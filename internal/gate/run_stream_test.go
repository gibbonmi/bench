package gate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// openTestStreamFile opens one stream file the way a run's own open does, so a test
// that needs a retained stream takes the real file rather than a stand-in writer.
func openTestStreamFile(t *testing.T) *os.File {
	t.Helper()
	path := filepath.Join(t.TempDir(), gateLogStreamName(seededRun(0)))
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_APPEND, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = file.Close() })
	return file
}

func readStreamFile(t *testing.T, file *os.File) string {
	t.Helper()
	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// BG24 at the seam this ticket owns: a line reaches the stream file when it arrives,
// not when the phase settles. The read below happens with the phase still running and
// its writers still open, which is the state a killed run leaves behind. A buffer
// flushed at settle answers nothing here, and such a run keeps nothing.
func TestPhaseStreamWritesEachLineWhenItArrivesNotWhenThePhaseSettles(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	file := openTestStreamFile(t)
	streams.retain(file)

	out, _, closeWriters := streams.open(Phase{Name: "test"})
	if _, err := io.WriteString(out, "first finding\nsecond finding\n"); err != nil {
		t.Fatal(err)
	}

	want := "[test] first finding\n[test] second finding\n"
	if got := readStreamFile(t, file); got != want {
		t.Fatalf("stream file mid-phase = %q, want %q", got, want)
	}
	closeWriters()
}

// Every line of a completed run reaches the file, in arrival order, under the phase that
// wrote it. The two streams of one phase interleave, and a second phase's lines file
// under their own name, so one file reads as one run.
func TestGateRunStreamHoldsEveryPhaseLineInArrivalOrder(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	file := openTestStreamFile(t)
	streams.retain(file)

	out, errOut, closeWriters := streams.open(Phase{Name: "vet"})
	io.WriteString(out, "vet said one\n")
	io.WriteString(errOut, "vet said two\n")
	io.WriteString(out, "vet said three")
	closeWriters()
	gofmtOut, _, closeGofmt := streams.open(Phase{Name: "gofmt"})
	io.WriteString(gofmtOut, "gofmt said one\n")
	closeGofmt()

	want := "[vet] vet said one\n[vet] vet said two\n[vet] vet said three\n[gofmt] gofmt said one\n"
	if got := readStreamFile(t, file); got != want {
		t.Fatalf("stream file = %q, want %q", got, want)
	}
}

// BG35: a run whose stream file opens names it once on stderr, beside the progress log
// line the run already prints. The path is the file the report's more-row sends a reader
// to, so the run must name the file it actually opened.
func TestGateRunNamesTheStreamFileOnceOnStderr(t *testing.T) {
	root := newLoggingPruneRoot(t)
	var stderr bytes.Buffer

	ctx, finish := beginGateRunLog(context.Background(), root, &stderr, "dev")
	defer finish(Result{})

	stream := gateRunStreamFile(ctx)
	if stream == nil {
		t.Fatal("the run retained no stream file, so this asserts nothing")
	}
	if got := strings.Count(stderr.String(), "gate: stream "); got != 1 {
		t.Fatalf("stderr named the stream %d times, want 1:\n%s", got, stderr.String())
	}
	if want := "gate: stream " + stream.Name() + "\n"; !strings.Contains(stderr.String(), want) {
		t.Errorf("stderr = %q, want it to carry %q", stderr.String(), want)
	}
	// A run that opened both files names both, and names each exactly once.
	if got := strings.Count(stderr.String(), "gate: progress log "); got != 1 {
		t.Errorf("stderr named the progress log %d times, want 1:\n%s", got, stderr.String())
	}
}

// The edge inventory's symlinked .logs: a run writes no stream through a link, because
// MkdirAll follows one whose target is a directory and the run would write outside the
// tree it names. The report then says the stream is unavailable rather than naming a
// file that holds nothing.
func TestGateRunWritesNoStreamThroughASymlinkedLogDir(t *testing.T) {
	root, target := newSymlinkedLogRoot(t)
	var stderr bytes.Buffer

	ctx, finish := beginGateRunLog(context.Background(), root, &stderr, "dev")
	defer finish(Result{})

	if stream := gateRunStreamFile(ctx); stream != nil {
		t.Fatalf("the run opened %s through a symlinked .logs", stream.Name())
	}
	entries, err := os.ReadDir(target)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), gateLogStreamSuffix) {
			t.Errorf("a stream file landed through the link: %s", entry.Name())
		}
	}
	if strings.Contains(stderr.String(), "gate: stream ") {
		t.Errorf("stderr named a stream the run never opened: %q", stderr.String())
	}
	streams := newPhaseStreams(io.Discard)
	streams.retain(gateRunStreamFile(ctx))
	if got := streams.path(); got != "" {
		t.Errorf("streams.path() = %q, want the stream reported unavailable", got)
	}
}

// BG26: a run whose stream file is refused leaves the table bounded and the more-row
// saying the stream is unavailable. The refusal is partial by construction here. The
// record's own open runs first and succeeds, and only the stream is refused, which is the
// state a run reaches when it keeps its verdict and loses the whole its table points at.
// A refused stream costs the reader that whole, never the bound, and the run names no
// file it did not open.
func TestARefusedStreamLeavesTheTableBoundedAndTheStreamUnavailable(t *testing.T) {
	root, _ := newSymlinkedLogRoot(t)

	var runErr bytes.Buffer
	ctx, finish := beginGateRunLog(context.Background(), root, &runErr, "dev")
	if strings.Contains(runErr.String(), "progress logging unavailable") {
		t.Fatalf("the run opened no record, so the refusal here is not the partial one: %q", runErr.String())
	}
	if stream := gateRunStreamFile(ctx); stream != nil {
		t.Fatalf("the run opened %s, so this asserts nothing about a refused stream", stream.Name())
	}
	if strings.Contains(runErr.String(), "gate: stream ") {
		t.Errorf("stderr named a stream the run never opened: %q", runErr.String())
	}
	streams := newPhaseStreams(io.Discard)
	streams.retain(gateRunStreamFile(ctx))
	var raw strings.Builder
	for i := 1; i <= 50; i++ {
		fmt.Fprintf(&raw, "finding %d\n", i)
	}
	writePhaseStream(t, streams, "vet", raw.String(), "")
	var stdout, stderr bytes.Buffer
	aggregateAndReport([]phaseResult{{Name: "vet", Argv: []string{"go", "vet", "./..."}, Code: 1}}, false, streams, &stdout, &stderr)

	rows := rowsForPhase(t, stdout.String(), "vet")
	if len(rows) != failureRowCap+1 {
		t.Fatalf("fifty findings yielded %d rows, want %d", len(rows), failureRowCap+1)
	}
	if want := "+30 more lines (stream unavailable)"; rows[failureRowCap] != want {
		t.Errorf("more-row = %q, want %q", rows[failureRowCap], want)
	}

	// The record the run did open still holds this run's own start and finish, so the
	// refusal above cost the run its stream alone.
	finish(Result{})
	record := gateRunLogFrom(t, ctx).file.Name()
	data, err := os.ReadFile(record)
	if err != nil {
		t.Fatalf("the record the run opened is unreadable: %v", err)
	}
	for _, event := range []string{`"event":"gate.start"`, `"event":"gate.finish"`} {
		if !strings.Contains(string(data), event) {
			t.Errorf("%s = %q, want it to carry %s", record, data, event)
		}
	}
}

// A run's stream file is never opened over a path some other run already took.
// openGateStreamFile opens with O_EXCL, so a taken name is a refusal rather than an
// append, and two runs sharing one instant cannot interleave their lines in one file.
func TestGateRunOpensNoStreamOverATakenPath(t *testing.T) {
	root := newPruneRoot(t)
	run := seededRun(0)
	if err := os.WriteFile(gateLogStreamPath(root, run), []byte("another run's lines\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if file := openGateStreamFile(root, run); file != nil {
		file.Close()
		t.Fatalf("the run opened %s over a path already taken", file.Name())
	}
	data, err := os.ReadFile(gateLogStreamPath(root, run))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "another run's lines\n" {
		t.Errorf("the taken file = %q, want it untouched", data)
	}
}

// newSymlinkedLogRoot builds a root whose .logs is a symlink to a directory elsewhere.
// MkdirAll and the record's own open both follow such a link, so the run keeps its
// record; only the stream's Lstat refuses. This is the one fixture that reaches a run
// with a record and no stream.
func newSymlinkedLogRoot(t *testing.T) (root, target string) {
	t.Helper()
	root, target = t.TempDir(), t.TempDir()
	if err := os.Symlink(target, gateLogDir(root)); err != nil {
		t.Fatal(err)
	}
	stubGateLogPathIgnored(t)
	return root, target
}

// gateRunLogFrom answers the run log a context carries, which is where a test asks what
// the run actually opened.
func gateRunLogFrom(t *testing.T, ctx context.Context) *gateRunLog {
	t.Helper()
	log, _ := ctx.Value(gateRunLogKey{}).(*gateRunLog)
	if log == nil {
		t.Fatal("the context carries no run log")
	}
	return log
}

// The child process that runs the phases appends to the parent's file. The parent hands
// the path over in the environment, exactly as it hands over the record's, so one run
// leaves one stream however many processes wrote it.
func TestGateRunStreamInheritsTheParentsFile(t *testing.T) {
	root := newLoggingPruneRoot(t)
	ctx, finish := beginGateRunLog(context.Background(), root, io.Discard, "dev")
	parent := gateRunStreamFile(ctx)
	if parent == nil {
		t.Fatal("the parent retained no stream file, so this asserts nothing")
	}
	setGateLogEnv(t, withGateRunLogEnv(ctx, nil))

	childCtx, closeChild := inheritGateRunLog(context.Background(), io.Discard)
	child := gateRunStreamFile(childCtx)
	if child == nil {
		t.Fatal("the child inherited no stream file")
	}
	if child.Name() != parent.Name() {
		t.Fatalf("child stream = %s, want the parent's %s", child.Name(), parent.Name())
	}

	streams := newPhaseStreams(io.Discard)
	streams.retain(child)
	writePhaseStream(t, streams, "vet", "child finding\n", "")
	closeChild()
	finish(Result{})

	if got := readStreamFile(t, parent); !strings.Contains(got, "[vet] child finding\n") {
		t.Errorf("the parent's file = %q, want the child's line in it", got)
	}
}

// A parent that retained no stream still hands its child a run to write records into.
// Only the record decides whether a run is inherited, so a refused stream costs the
// child its file and nothing else.
func TestGateRunInheritsTheRecordWhenTheParentRetainedNoStream(t *testing.T) {
	root := newLoggingPruneRoot(t)
	run := seededRun(0)
	seedRecords(t, root, 0)
	setGateLogEnv(t, []string{
		gateLogPathEnv + "=" + gateLogRecordPath(root, run),
		gateLogRootEnv + "=" + root,
		gateLogRunEnv + "=" + run,
		gateLogStreamEnv + "=",
	})

	childCtx, closeChild := inheritGateRunLog(context.Background(), io.Discard)
	defer closeChild()

	log, _ := childCtx.Value(gateRunLogKey{}).(*gateRunLog)
	if log == nil {
		t.Fatal("a parent without a stream left the child without a run log")
	}
	if log.streamPath() != "" {
		t.Errorf("the child opened %q, want no stream", log.streamPath())
	}
}

// setGateLogEnv puts one composed environment into this test's process, which is where
// the inheriting child reads it from.
func setGateLogEnv(t *testing.T, env []string) {
	t.Helper()
	for _, name := range gateLogEnv {
		t.Setenv(name, "")
	}
	for _, item := range env {
		name, value, _ := strings.Cut(item, "=")
		t.Setenv(name, value)
	}
}
