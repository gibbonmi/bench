package recordcmd_test

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/toon"
)

func readyCompletion(t *testing.T) (*recordtest.Fixture, string) {
	t.Helper()
	f := recorded(t, 1)
	verify(t, f)
	verify(t, f, unprobed)
	for _, axis := range rr.Axes() {
		review(t, f, axis, "review-"+axis)
	}
	source := f.Tip()
	verify(t, f, final(source))
	verify(t, f, final(source), map[string]string{"--requirement": "integration", "--id": "integration-1"})
	return f, source
}

func TestRecordCompletionWritesDerivedEntry(t *testing.T) {
	f, source := readyCompletion(t)
	out, code := recordcmd.Command(f.Root, []string{"completion", slug(t), "--source", source})
	if code != 0 {
		t.Fatalf("completion = exit %d, output %q; want exit 0", code, out)
	}
	got := read(t, f).Completion
	if got.State != "completed" || got.SourceDigest != onlyChunk(t, f).SourceDigest || got.Performer != recordtest.Orchestrator {
		t.Fatalf("completion = %s at %s by %s, want completed at the chunk source by %s", got.State, got.SourceDigest, got.Performer, recordtest.Orchestrator)
	}
	if want := map[string]string{"E1": "covered"}; !reflect.DeepEqual(got.Reconciliation, want) {
		t.Fatalf("reconciliation = %v, want %v", got.Reconciliation, want)
	}
	if len(got.Verification) != len(f.Plan.FinalVerification) {
		t.Fatalf("final verification count = %d, want %d preserved", len(got.Verification), len(f.Plan.FinalVerification))
	}
	want, err := toon.Table("completion", []string{"state", "source_digest", "performer", "rows", "verification"},
		[][]string{{"completed", got.SourceDigest, recordtest.Orchestrator, "1", "2"}})
	if err != nil {
		t.Fatal(err)
	}
	if out != want {
		t.Fatalf("output = %q, want %q", out, want)
	}
}

func TestRecordCompletionRefusesIncompleteEvidenceWithoutWrite(t *testing.T) {
	one := 1
	tests := []struct {
		name, want string
		mutate     func(*rr.Record)
	}{
		{"author verification", "missing verification additional", func(record *rr.Record) {
			record.Chunks[0].Verification = record.Chunks[0].Verification[:1]
		}},
		{"review", "missing Coverage", func(record *rr.Record) {
			record.Chunks[0].Reviews = record.Chunks[0].Reviews[:2]
		}},
		{"final verification", "verification acceptance is incomplete, failed, or stale", func(record *rr.Record) {
			record.Completion.Verification[0].Outcome = "fail"
			record.Completion.Verification[0].ExitCode = &one
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, source := readyCompletion(t)
			record := read(t, f)
			test.mutate(&record)
			f.ReplaceRecord(record)
			path := recordFile(t, f)
			before := bytesOf(t, path)
			out, code := recordcmd.Command(f.Root, []string{"completion", slug(t), "--source", source})
			if code != 1 || !strings.Contains(out, test.want) {
				t.Fatalf("completion = exit %d, output %q; want exit 1 containing %q", code, out, test.want)
			}
			if after := bytesOf(t, path); !bytes.Equal(after, before) {
				t.Fatalf("record changed after refusal\nbefore: %s\nafter: %s", before, after)
			}
		})
	}
}

func TestRecordCompletionPreservesProseAndReplaysByteForByte(t *testing.T) {
	f, source := readyCompletion(t)
	path := recordFile(t, f)
	document := strings.Replace(string(bytesOf(t, path)), "# Review outcomes\n\n", "# Review outcomes\n\nKeep this prose.\n\n", 1) + "After the record.\n"
	if err := os.WriteFile(path, []byte(document), 0o644); err != nil {
		t.Fatal(err)
	}
	args := []string{"completion", slug(t), "--source", source}
	if out, code := recordcmd.Command(f.Root, args); code != 0 {
		t.Fatalf("first completion = exit %d, output %q; want exit 0", code, out)
	}
	first := bytesOf(t, path)
	if !bytes.Contains(first, []byte("Keep this prose.")) || !bytes.HasSuffix(first, []byte("After the record.\n")) {
		t.Fatalf("completion did not preserve surrounding prose: %s", first)
	}
	if out, code := recordcmd.Command(f.Root, args); code != 0 {
		t.Fatalf("replayed completion = exit %d, output %q; want exit 0", code, out)
	}
	if second := bytesOf(t, path); !bytes.Equal(second, first) {
		t.Fatalf("replayed record changed bytes\nfirst: %s\nsecond: %s", first, second)
	}
}
