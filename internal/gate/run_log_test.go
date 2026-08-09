package gate

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestGateRunLogRecordsMajorEvents(t *testing.T) {
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	writeGateTestFile(t, root, ".gitignore", ".logs/\n", 0o644)
	var stderr bytes.Buffer
	ctx, finish := beginGateRunLog(context.Background(), root, &stderr, "fresh")
	exit := 0
	logGateEvent(ctx, gateLogRecord{Event: "phase.start", Phase: "conformance"})
	logGateEvent(ctx, gateLogRecord{Event: "phase.finish", Phase: "conformance", Exit: &exit, ElapsedMS: 12})
	finish(Result{})

	entries, err := os.ReadDir(filepath.Join(root, ".logs"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("gate log entries = %v, %v, want one", entries, err)
	}
	file, err := os.Open(filepath.Join(root, ".logs", entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var records []gateLogRecord
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record gateLogRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			t.Fatalf("decode gate log: %v", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	want := []string{"gate.start", "phase.start", "phase.finish", "gate.finish"}
	if len(records) != len(want) {
		t.Fatalf("gate log records = %#v, want events %v", records, want)
	}
	for i, event := range want {
		if records[i].Schema != gateLogSchema || records[i].Run == "" || records[i].Time == "" || records[i].Event != event {
			t.Fatalf("record %d = %#v, want schema/run/time and event %q", i, records[i], event)
		}
	}
	if records[2].Exit == nil || *records[2].Exit != 0 || records[2].ElapsedMS != 12 {
		t.Fatalf("phase finish = %#v, want exit 0 and 12ms", records[2])
	}
	if !bytes.Contains(stderr.Bytes(), []byte(filepath.Join(root, ".logs"))) {
		t.Fatalf("progress log announcement = %q, want log path", stderr.String())
	}
}

func TestGateRunLogDoesNotWriteIntoAnUnignoredSubject(t *testing.T) {
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	ctx, finish := beginGateRunLog(context.Background(), root, io.Discard, "ordinary")
	logGateEvent(ctx, gateLogRecord{Event: "phase.start", Phase: "test"})
	finish(Result{})
	if _, err := os.Stat(filepath.Join(root, ".logs")); !os.IsNotExist(err) {
		t.Fatalf("unignored log directory = %v, want absent", err)
	}
}

func TestInheritedGateRunLogAppendsPhaseEventsAndSuppressesInnerCanary(t *testing.T) {
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	writeGateTestFile(t, root, ".gitignore", ".logs/\n", 0o644)
	ctx, finish := beginGateRunLog(context.Background(), root, io.Discard, "fresh")
	log := ctx.Value(gateRunLogKey{}).(*gateRunLog)
	t.Setenv(gateLogPathEnv, log.file.Name())
	t.Setenv(gateLogRootEnv, log.root)
	t.Setenv(gateLogRunEnv, log.run)

	inherited, closeInherited := inheritGateRunLog(context.Background(), io.Discard)
	logGateEvent(inherited, gateLogRecord{Event: "phase.start", Phase: "test"})
	closeInherited()

	t.Setenv("BENCH_CANARY_INNER", "1")
	inner, closeInner := inheritGateRunLog(context.Background(), io.Discard)
	logGateEvent(inner, gateLogRecord{Event: "phase.start", Phase: "inner"})
	closeInner()
	finish(Result{})

	data, err := os.ReadFile(log.file.Name())
	if err != nil {
		t.Fatal(err)
	}
	if got := bytes.Count(data, []byte(`"event":"phase.start"`)); got != 1 {
		t.Fatalf("phase-start records = %d, want one inherited outer event; log=%s", got, data)
	}
	if bytes.Contains(data, []byte(`"phase":"inner"`)) {
		t.Fatalf("inner canary wrote to outer log: %s", data)
	}
}
