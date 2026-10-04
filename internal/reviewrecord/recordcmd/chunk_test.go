package recordcmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/toon"
)

// linked is a delegated fixture in a linked worktree with count planned chunks.
func linked(t *testing.T, count int) *recordtest.Fixture {
	t.Helper()
	f, _ := recordtest.NewLinked(t, count, recordtest.Delegate)
	return f
}

// recorded is a fixture with count planned chunks whose chunk 1 entry is written.
func recorded(t *testing.T, count int, prepare ...func(*recordtest.Fixture)) *recordtest.Fixture {
	t.Helper()
	f := linked(t, count)
	for _, apply := range prepare {
		apply(f)
	}
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	return f
}

func slug(t *testing.T) string {
	t.Helper()
	slug, err := rr.Slug(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	return slug
}

func recordFile(t *testing.T, f *recordtest.Fixture) string {
	t.Helper()
	path, err := rr.RecordPath(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(f.Root, filepath.FromSlash(path))
}

// advance commits one source change and returns the commit before it and the commit it made.
func advance(f *recordtest.Fixture, label string) (string, string) {
	base := f.Tip()
	f.Write("source.txt", label+"\n")
	f.Commit(label)
	return base, f.Tip()
}

func chunkArgs(t *testing.T, id, base, tip string) []string {
	return []string{"chunk", slug(t), "--chunk", id, "--base", base, "--tip", tip}
}

// record runs the chunk form and requires exit 0.
func record(t *testing.T, f *recordtest.Fixture, id, base, tip string) string {
	t.Helper()
	out, code := recordcmd.Command(f.Root, chunkArgs(t, id, base, tip))
	if code != 0 {
		t.Fatalf("chunk %s exit = %d, want 0; output=%q", id, code, out)
	}
	return out
}

func read(t *testing.T, f *recordtest.Fixture) rr.Record {
	t.Helper()
	record, err := rr.Read(f.Root, recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func bytesOf(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func onlyChunk(t *testing.T, f *recordtest.Fixture) rr.Chunk {
	t.Helper()
	record := read(t, f)
	if len(record.Chunks) != 1 {
		t.Fatalf("record chunks = %d, want 1", len(record.Chunks))
	}
	return record.Chunks[0]
}

func TestRecordChunkWritesFullCommitIDs(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base[:12], "HEAD")
	chunk := onlyChunk(t, f)
	if chunk.Base != base || chunk.Tip != tip {
		t.Fatalf("chunk pair = %s..%s, want the full IDs %s..%s", chunk.Base, chunk.Tip, base, tip)
	}
}

func TestRecordChunkSourceDigestExcludesTheRecord(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	f.Commit("record chunk 1")
	record(t, f, "1", base, "HEAD")
	want, err := rr.SourceDigest(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	chunk := onlyChunk(t, f)
	if chunk.SourceDigest != want || chunk.SourceDigest == f.Tree() {
		t.Fatalf("source digest = %s, want %s, which differs from the tip tree %s", chunk.SourceDigest, want, f.Tree())
	}
}

func TestRecordChunkPlanDigestReadsThePlan(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	plan, err := rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if chunk := onlyChunk(t, f); chunk.PlanDigest != plan.Digest {
		t.Fatalf("plan digest = %s, want %s", chunk.PlanDigest, plan.Digest)
	}
}

func TestRecordChunkCopiesAcceptanceRows(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	if rows := onlyChunk(t, f).AcceptanceRows; !reflect.DeepEqual(rows, []string{"E1"}) {
		t.Fatalf("acceptance rows = %q, want [E1]", rows)
	}
}

func TestRecordChunkCreatesAVersionTwoRecord(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	got := read(t, f)
	plan, err := rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != 2 || got.Spec != recordtest.Spec || got.PlanDigest != plan.Digest || got.ImplementationSession != "" || len(got.Chunks) != 1 {
		t.Fatalf("record = %+v, want version 2, %s, plan digest %s, no session, and one chunk", got, recordtest.Spec, plan.Digest)
	}
}

func TestRecordChunkCreatesInAnEmptyDirectory(t *testing.T) {
	f := linked(t, 1)
	if err := os.Mkdir(filepath.Dir(recordFile(t, f)), 0o755); err != nil {
		t.Fatal(err)
	}
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	if _, err := os.Lstat(recordFile(t, f)); err != nil {
		t.Fatalf("record file: %v", err)
	}
}

func TestRecordChunkRefusesAVersionOnePlan(t *testing.T) {
	f, _ := recordtest.NewLinked(t, 1)
	base, tip := advance(f, "chunk 1")
	out, code := recordcmd.Command(f.Root, chunkArgs(t, "1", base, tip))
	if code != 1 || !strings.Contains(out, "version 2") {
		t.Fatalf("version 1 plan = exit %d, output %q; want exit 1 naming version 2", code, out)
	}
	if _, err := os.Lstat(recordFile(t, f)); !os.IsNotExist(err) {
		t.Fatalf("record file after refusal: %v, want absent", err)
	}
}

func TestRecordChunkAppendsAfterTheRecordedChunk(t *testing.T) {
	f := linked(t, 2)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	base, tip = advance(f, "chunk 2")
	record(t, f, "2", base, tip)
	chunks := read(t, f).Chunks
	if len(chunks) != 2 || chunks[0].ID != "1" || chunks[1].ID != "2" {
		t.Fatalf("chunks = %+v, want chunk 1 then chunk 2", chunks)
	}
}

func TestRecordChunkUpdateKeepsResults(t *testing.T) {
	f := linked(t, 1)
	base := f.Tip()
	f.RecordChunk(base)
	f.Save()
	_, tip := advance(f, "repair chunk 1")
	record(t, f, "1", base, tip)
	chunk := onlyChunk(t, f)
	if chunk.Tip != tip || len(chunk.Verification) != 2 || len(chunk.Reviews) != 3 {
		t.Fatalf("chunk = tip %s, %d verification and %d review results; want tip %s with 2 and 3", chunk.Tip, len(chunk.Verification), len(chunk.Reviews), tip)
	}
}

func TestRecordChunkRepeatIsByteIdentical(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	first := bytesOf(t, recordFile(t, f))
	record(t, f, "1", base, tip)
	if second := bytesOf(t, recordFile(t, f)); !bytes.Equal(first, second) {
		t.Fatalf("second run changed the record:\n%s\nwant:\n%s", second, first)
	}
}

func TestRecordChunkRefusesAnUnplannedChunk(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	before := bytesOf(t, recordFile(t, f))
	out, code := recordcmd.Command(f.Root, chunkArgs(t, "9", base, tip))
	if code != 1 || !strings.Contains(out, "chunk 9") {
		t.Fatalf("unplanned chunk = exit %d, output %q; want exit 1 naming chunk 9", code, out)
	}
	if after := bytesOf(t, recordFile(t, f)); !bytes.Equal(before, after) {
		t.Fatalf("refusal changed the record")
	}
}

func TestRecordChunkRefusesAnUnknownRevision(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	before := bytesOf(t, recordFile(t, f))
	out, code := recordcmd.Command(f.Root, chunkArgs(t, "1", base, "no-such-rev"))
	if code != 1 {
		t.Fatalf("unknown revision = exit %d, output %q; want exit 1", code, out)
	}
	if after := bytesOf(t, recordFile(t, f)); !bytes.Equal(before, after) {
		t.Fatalf("refusal changed the record")
	}
}

// wantRow is the one output table of a chunk form, derived from the written entry.
func wantRow(t *testing.T, action string, chunk rr.Chunk) string {
	t.Helper()
	table, err := toon.Table("chunk", []string{"id", "action", "base", "tip", "source_digest", "plan_digest", "rows"},
		[][]string{{chunk.ID, action, chunk.Base, chunk.Tip, chunk.SourceDigest, chunk.PlanDigest, strconv.Itoa(len(chunk.AcceptanceRows))}})
	if err != nil {
		t.Fatal(err)
	}
	return table
}

func TestRecordChunkReportsCreated(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	out := record(t, f, "1", base, tip)
	if want := wantRow(t, "created", onlyChunk(t, f)); out != want {
		t.Fatalf("output = %q, want %q", out, want)
	}
}

func TestRecordChunkReportsAdded(t *testing.T) {
	f := linked(t, 2)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	base, tip = advance(f, "chunk 2")
	out := record(t, f, "2", base, tip)
	if want := wantRow(t, "added", read(t, f).Chunks[1]); out != want {
		t.Fatalf("output = %q, want %q", out, want)
	}
}

func TestRecordChunkReportsUpdated(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	_, tip = advance(f, "repair chunk 1")
	out := record(t, f, "1", base, tip)
	if want := wantRow(t, "updated", onlyChunk(t, f)); out != want {
		t.Fatalf("output = %q, want %q", out, want)
	}
}

func TestRecordKeepsTheProseAroundTheFence(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	const before, after = "# Review outcomes\n\nThe pickup prose before the fence.\n\n", "\n## Pickup\n\nThe prose after the fence.\n"
	path := recordFile(t, f)
	if err := os.WriteFile(path, []byte(before+string(bytesOf(t, path))[len("# Review outcomes\n\n"):]+after), 0o644); err != nil {
		t.Fatal(err)
	}
	_, tip = advance(f, "repair chunk 1")
	record(t, f, "1", base, tip)
	if got := string(bytesOf(t, path)); !strings.HasPrefix(got, before) || !strings.HasSuffix(got, "```\n"+after) {
		t.Fatalf("record = %q, want the prose %q before and %q after the fence", got, before, after)
	}
}

func TestRecordChangesOnlyTheRecordPath(t *testing.T) {
	f := linked(t, 1)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	path, err := rr.RecordPath(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if status := f.Git("status", "--porcelain", "--untracked-files=all"); status != "?? "+path {
		t.Fatalf("git status = %q, want only %q", status, path)
	}
}
