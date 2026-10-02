package recordcmd_test

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/toon"
)

// amendArgs is the amendment argv at HEAD. Each value follows as one --map flag, in order.
func amendArgs(t *testing.T, maps ...string) []string {
	t.Helper()
	args := formArgs(t, "amendment", map[string]string{"--source": "HEAD"})
	for _, value := range maps {
		args = append(args, "--map", value)
	}
	return args
}

// amend runs the amendment form and requires exit 0.
func amend(t *testing.T, f *recordtest.Fixture, maps ...string) string {
	t.Helper()
	return succeed(t, f, amendArgs(t, maps...))
}

// replan commits a plan change that keeps every chunk ID: the first final requirement
// runs command.
func replan(f *recordtest.Fixture, command string) {
	f.Plan.FinalVerification[0].Command = command
	f.RewritePlan()
}

// rename commits a plan whose chunks carry ids, in order.
func rename(f *recordtest.Fixture, ids ...string) {
	for i, id := range ids {
		f.Plan.Chunks[i].ID = id
	}
	f.RewritePlan()
}

// onlyAmendment is the one amendment in the record of f.
func onlyAmendment(t *testing.T, f *recordtest.Fixture) rr.Amendment {
	t.Helper()
	amendments := read(t, f).Amendments
	if len(amendments) != 1 {
		t.Fatalf("record amendments = %d, want 1", len(amendments))
	}
	return amendments[0]
}

func TestRecordAmendmentComputesBothDigests(t *testing.T) {
	f := recorded(t, 1)
	from := read(t, f).PlanDigest
	replan(f, "go test -count=1 ./...")
	amend(t, f)
	plan, err := rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if got := onlyAmendment(t, f); got.From != from || got.To != plan.Digest || got.From == got.To {
		t.Fatalf("amendment = %s to %s, want %s to the plan digest at HEAD %s", got.From, got.To, from, plan.Digest)
	}
}

func TestRecordAmendmentMovesThePlanDigest(t *testing.T) {
	f := recorded(t, 1)
	replan(f, "go test -count=1 ./...")
	amend(t, f)
	if got := read(t, f); got.PlanDigest != onlyAmendment(t, f).To {
		t.Fatalf("record plan digest = %s, want the amendment's to %s", got.PlanDigest, onlyAmendment(t, f).To)
	}
}

func TestRecordAmendmentMapsEachChunkToItself(t *testing.T) {
	f := recorded(t, 2)
	base, tip := advance(f, "chunk 2")
	record(t, f, "2", base, tip)
	replan(f, "go test -count=1 ./...")
	amend(t, f)
	if got := onlyAmendment(t, f).ChunkIDs; !reflect.DeepEqual(got, map[string][]string{"1": {"1"}, "2": {"2"}}) {
		t.Fatalf("chunk mapping = %q, want each recorded chunk mapped to itself", got)
	}
}

func TestRecordAmendmentChainsFromTheLastDigest(t *testing.T) {
	f := recorded(t, 2)
	rename(f, "1a", "1b")
	amend(t, f, "1=1a,1b")
	replan(f, "go test -count=1 ./...")
	amend(t, f)
	amendments := read(t, f).Amendments
	if len(amendments) != 2 || amendments[1].From != amendments[0].To || amendments[1].To != read(t, f).PlanDigest {
		t.Fatalf("amendments = %+v, want the second to chain from the first's to", amendments)
	}
	if got := amendments[1].ChunkIDs; !reflect.DeepEqual(got, map[string][]string{"1a": {"1a"}, "1b": {"1b"}}) {
		t.Fatalf("second chunk mapping = %q, want chunk 1 mapped through the first amendment to 1a and 1b", got)
	}
}

func TestRecordAmendmentWritesAMappedSplit(t *testing.T) {
	f := recorded(t, 2)
	rename(f, "1a", "1b")
	amend(t, f, "1=1a,1b")
	if got := onlyAmendment(t, f).ChunkIDs; !reflect.DeepEqual(got, map[string][]string{"1": {"1a", "1b"}}) {
		t.Fatalf("chunk mapping = %q, want chunk 1 split into 1a and 1b", got)
	}
}

func TestRecordAmendmentRefusesAnUnchangedPlan(t *testing.T) {
	f := recorded(t, 1)
	refuseArgs(t, f, "unchanged", amendArgs(t))
}

func TestRecordAmendmentRefusesAnUnmappedChunk(t *testing.T) {
	f := recorded(t, 2)
	rename(f, "1a", "2")
	refuseArgs(t, f, "chunk 1 ", amendArgs(t))
}

func TestRecordAmendmentRefusesAMapForAnUnrecordedChunk(t *testing.T) {
	f := recorded(t, 1)
	replan(f, "go test -count=1 ./...")
	refuseArgs(t, f, "chunk 9 ", amendArgs(t, "9=1"))
}

func TestRecordAmendmentRefusesAMapToAnUnplannedChunk(t *testing.T) {
	f := recorded(t, 1)
	replan(f, "go test -count=1 ./...")
	refuseArgs(t, f, "chunk 9 ", amendArgs(t, "1=9"))
}

func TestRecordAmendmentRefusesAnUnresolvableChain(t *testing.T) {
	f := recorded(t, 1)
	original := f.Plan.FinalVerification[0].Command
	replan(f, "go test -count=1 ./...")
	amend(t, f)
	replan(f, original)
	amend(t, f)
	replan(f, "go test -count=2 ./...")
	refuseArgs(t, f, "ambiguous", amendArgs(t))
}

func TestRecordAmendmentRefusesAControlCharacter(t *testing.T) {
	f := recorded(t, 1)
	replan(f, "go test -count=1 ./...")
	refuseArgs(t, f, "--map", amendArgs(t, "1=1\x1b"))
}

func TestRecordAmendmentNeedsARecord(t *testing.T) {
	f := linked(t, 1)
	replan(f, "go test -count=1 ./...")
	out, code := recordcmd.Command(f.Root, amendArgs(t))
	if code != 1 || !strings.Contains(out, "bench record chunk") {
		t.Fatalf("amendment without a record = exit %d, output %q; want exit 1 naming bench record chunk", code, out)
	}
	if _, err := os.Lstat(recordFile(t, f)); !os.IsNotExist(err) {
		t.Fatalf("record file after refusal: %v, want absent", err)
	}
}

func TestRecordAmendmentReportsItsRow(t *testing.T) {
	f := recorded(t, 2)
	rename(f, "1a", "1b")
	out := amend(t, f, "1=1a,1b")
	got := onlyAmendment(t, f)
	want, err := toon.Table("amendment", []string{"from", "to", "chunks"}, [][]string{{got.From, got.To, strconv.Itoa(len(got.ChunkIDs))}})
	if err != nil {
		t.Fatal(err)
	}
	if out != want || !strings.HasPrefix(out, "amendment[1]{from,to,chunks}:") || len(got.ChunkIDs) != 1 {
		t.Fatalf("output = %q, want %q", out, want)
	}
}

func TestRecordedAmendmentPassesTheCheckpoint(t *testing.T) {
	f := recorded(t, 2)
	verify(t, f)
	verify(t, f, unprobed)
	for _, axis := range rr.Axes() {
		review(t, f, axis, "review-1-"+axis)
	}
	f.Commit("record chunk 1")
	base := f.Tip()
	f.Plan.FinalVerification[0].Command = "go test -count=1 ./..."
	f.WritePlan()
	_, tip := advance(f, "chunk 2")
	f.Reload()
	if err := rr.Check(f.Root, f.Tree(), f.Tip(), recordtest.Spec, "2", false); err == nil || !strings.Contains(err.Error(), "stale plan digest") {
		t.Fatalf("checkpoint before the amendment = %v, want stale plan digest", err)
	}
	amend(t, f)
	record(t, f, "2", base, tip)
	chunk2 := map[string]string{"--chunk": "2", "--performer": recordtest.Author("2.md")}
	verify(t, f, chunk2, map[string]string{"--id": "tests-2"})
	verify(t, f, chunk2, unprobed, map[string]string{"--id": "additional-2"})
	for _, axis := range rr.Axes() {
		succeed(t, f, reviewArgs(t, axis, "review-2-"+axis, nil, map[string]string{"--chunk": "2"}))
	}
	f.Commit("record chunk 2")
	if err := rr.Check(f.Root, f.Tree(), f.Tip(), recordtest.Spec, "2", false); err != nil {
		t.Fatalf("checkpoint for chunk 2: %v", err)
	}
}

func TestRecordAmendmentGrammarRefusals(t *testing.T) {
	f := recorded(t, 1)
	replan(f, "go test -count=1 ./...")
	for name, maps := range map[string][]string{
		"no equals":    {"1"},
		"empty old":    {"=1"},
		"empty new":    {"1="},
		"empty new ID": {"1=1a,"},
		"repeated key": {"1=2", "1=3"},
	} {
		out, code := recordcmd.Command(f.Root, amendArgs(t, maps...))
		if code != 2 || !strings.HasPrefix(out, "usage: bench record amendment") {
			t.Errorf("%s = exit %d, output %q; want exit 2 with a usage: bench record amendment line", name, code, out)
		}
	}
}
