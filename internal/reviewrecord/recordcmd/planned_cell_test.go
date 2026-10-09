package recordcmd_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/coverage"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// The seam cells of the fixture maps: one cites a test, and one plans a test.
const (
	cited   = "`internal/x/x_test.go` (`TestX`)"
	planned = "planned TestX in internal/x/x_test.go"
)

// withMap writes an opted-in coverage map into the fixture spec, with one row for each seam
// cell: E1 holds the first cell, E2 the second, and so on. The caller commits the change.
func withMap(seams ...string) func(*recordtest.Fixture) {
	return func(f *recordtest.Fixture) {
		f.T.Helper()
		var b strings.Builder
		b.WriteString("\n### Acceptance coverage map\n\n| row | story | behavior | seam | why it catches the failure |\n| --- | --- | --- | --- | --- |\n")
		for i, seam := range seams {
			fmt.Fprintf(&b, "| E%d | 1 | b | %s | w |\n", i+1, seam)
		}
		f.Write(recordtest.Spec, specText(f)+b.String())
	}
}

// specText is the fixture spec in the working copy.
func specText(f *recordtest.Fixture) string {
	f.T.Helper()
	data, err := os.ReadFile(filepath.Join(f.Root, filepath.FromSlash(recordtest.Spec)))
	if err != nil {
		f.T.Fatal(err)
	}
	return string(data)
}

// oneChunk puts every ticket of the plan into its first chunk, so that chunk covers two rows.
func oneChunk(plan *rr.Plan) {
	for _, chunk := range plan.Chunks[1:] {
		plan.Chunks[0].Tickets = append(plan.Chunks[0].Tickets, chunk.Tickets...)
	}
	plan.Chunks = plan.Chunks[:1]
}

// freeze commits the map that seams gives into a new chunk 1 fixture with count chunks, and
// runs the chunk form for chunk 1 over that commit.
func freeze(t *testing.T, count int, seams ...string) (*recordtest.Fixture, string, int) {
	t.Helper()
	f := linked(t, count)
	withMap(seams...)(f)
	base, tip := advance(f, "chunk 1")
	out, code := recordcmd.Command(f.Root, chunkArgs(t, "1", base, tip))
	return f, out, code
}

func TestRecordChunkRefusesAPlannedCell(t *testing.T) {
	f, out, code := freeze(t, 1, planned)
	if code != 1 || !strings.Contains(out, "bench record chunk refused") {
		t.Fatalf("planned chunk row = exit %d, output %q; want exit 1 with the chunk refusal", code, out)
	}
	if _, err := os.Lstat(recordFile(t, f)); !os.IsNotExist(err) {
		t.Fatalf("record file after refusal: %v, want absent", err)
	}
}

func TestRecordChunkNamesEachPlannedRow(t *testing.T) {
	f, _ := recordtest.NewLinked(t, 2, oneChunk, recordtest.Delegate)
	withMap(planned, planned)(f)
	base, tip := advance(f, "chunk 1")
	out, code := recordcmd.Command(f.Root, chunkArgs(t, "1", base, tip))
	if code != 1 || !strings.Contains(out, "E1") || !strings.Contains(out, "E2") || !strings.Contains(out, coverage.CitationForm) {
		t.Fatalf("two planned rows = exit %d, output %q; want exit 1 naming E1, E2, and %s", code, out, coverage.CitationForm)
	}
}

func TestRecordChunkRefusalWritesNothing(t *testing.T) {
	f := linked(t, 2)
	withMap(cited, planned)(f)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	before := bytesOf(t, recordFile(t, f))
	base, tip = advance(f, "chunk 2")
	if out, code := recordcmd.Command(f.Root, chunkArgs(t, "2", base, tip)); code != 1 {
		t.Fatalf("planned chunk 2 row = exit %d, output %q; want exit 1", code, out)
	}
	if after := bytesOf(t, recordFile(t, f)); !bytes.Equal(before, after) {
		t.Fatalf("refusal changed the record:\n%s\nwant:\n%s", after, before)
	}
}

func TestRecordChunkAcceptsACitedRow(t *testing.T) {
	if _, out, code := freeze(t, 1, cited); code != 0 {
		t.Fatalf("cited chunk row = exit %d, output %q; want exit 0", code, out)
	}
}

func TestRecordChunkAcceptsAReviewOwnedRow(t *testing.T) {
	if _, out, code := freeze(t, 1, "review-owned: the planned canary stays green"); code != 0 {
		t.Fatalf("review-owned chunk row = exit %d, output %q; want exit 0", code, out)
	}
}

func TestRecordChunkIgnoresAnotherChunksRow(t *testing.T) {
	if _, out, code := freeze(t, 2, cited, planned); code != 0 {
		t.Fatalf("planned row of chunk 2 = exit %d, output %q; want chunk 1 to record", code, out)
	}
}

func TestRecordChunkReadsTheTip(t *testing.T) {
	f := linked(t, 1)
	withMap(cited)(f)
	base, tip := advance(f, "chunk 1")
	f.Write(recordtest.Spec, strings.Replace(specText(f), cited, planned, 1))
	record(t, f, "1", base, tip)
}

// completed writes complete chunk and final evidence through the fixture for a spec whose
// map holds seams, then runs the completion form at the chunk tip. The fixture writes the
// chunk entry itself, so a planned cell that the chunk form refuses still reaches completion.
func completed(t *testing.T, seams ...string) (string, int) {
	t.Helper()
	f := linked(t, 1)
	withMap(seams...)(f)
	base, source := advance(f, "chunk 1")
	f.Reload()
	f.Record.PlanDigest = f.Plan.Digest
	f.RecordChunk(base)
	f.Complete()
	f.Save()
	return recordcmd.Command(f.Root, []string{"completion", slug(t), "--source", source})
}

func TestRecordCompletionRefusesAPlannedCell(t *testing.T) {
	if out, code := completed(t, planned); code != 1 || !strings.Contains(out, "E1") {
		t.Fatalf("planned plan row = exit %d, output %q; want exit 1 naming E1", code, out)
	}
}

func TestRecordCompletionAcceptsCitedRows(t *testing.T) {
	if out, code := completed(t, cited); code != 0 {
		t.Fatalf("cited plan rows = exit %d, output %q; want exit 0", code, out)
	}
}
