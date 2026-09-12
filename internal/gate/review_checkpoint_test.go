package gate

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

func checkpointFixture(t *testing.T) *recordtest.Fixture {
	return attachedCheckpointFixture(t, recordtest.Attach)
}

// attachedCheckpointFixture builds one recorded chunk from either record form.
// The version 1 and delegated checkpoints share this one build sequence, so a
// change to it cannot drift between them.
func attachedCheckpointFixture(t *testing.T, attach func(testing.TB, string, int) *recordtest.Fixture) *recordtest.Fixture {
	t.Helper()
	f := attach(t, outcomeFixture(t), 1)
	f.AddChunk()
	f.Save()
	f.Commit("retain review evidence")
	return f
}

func runCheckpoint(t *testing.T, f *recordtest.Fixture) (int, string) {
	t.Helper()
	var out, err bytes.Buffer
	code := RunCommand([]string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", "1"}, &out, &err)
	return code, out.String() + err.String()
}

func TestReviewCheckpoint(t *testing.T) {
	cases := []struct {
		name, reason string
		change       func(*recordtest.Fixture)
	}{
		{"pending axis", "pending", func(f *recordtest.Fixture) { f.Record.Chunks[0].Reviews[0].State = "pending" }},
		{"failed transport", "failed", func(f *recordtest.Fixture) { f.Record.Chunks[0].Reviews[0].State = "failed" }},
		{"skipped axis", "skipped", func(f *recordtest.Fixture) { f.Record.Chunks[0].Reviews[0].State = "skipped" }},
		{"missing axis", "missing Coverage", func(f *recordtest.Fixture) { f.Record.Chunks[0].Reviews = f.Record.Chunks[0].Reviews[:2] }},
		{"partial verification", "verification tests", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification = f.Record.Chunks[0].Verification[1:] }},
		{"partial additional verification", "verification additional", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification = f.Record.Chunks[0].Verification[:1] }},
		{"missing verification", "verification tests", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification = nil }},
		{"pending verification", "verification tests", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification[0].State = "pending" }},
		{"failed verification", "verification tests", func(f *recordtest.Fixture) { one := 1; f.Record.Chunks[0].Verification[0].ExitCode = &one }},
		{"stale verification", "verification tests", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification[0].SourceDigest = f.Record.Chunks[0].Base }},
		{"missing verifier", "terminal source", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification[0].Performer = "" }},
		{"wrong verifier", "verification tests", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification[0].Performer = "other" }},
		{"missing probe", "probe", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification[0].Probe = nil }},
		{"silent probe", "probe", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification[0].Probe.Outcome = "silent" }},
		{"failed restore", "restore", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification[0].Probe.Restore = "failed" }},
		{"review as verification", "author verification", func(f *recordtest.Fixture) { f.Record.Chunks[0].Verification[0].Role = "independent-review" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := checkpointFixture(t)
			tc.change(f)
			f.Save()
			code, out := runCheckpoint(t, f)
			if code == 0 || !strings.Contains(out, tc.reason) {
				t.Fatalf("checkpoint accepted invalid evidence or lost reason: exit %d: %s", code, out)
			}
			if _, err := os.Stat(filepath.Join(f.Root, ".gate-run-count")); !os.IsNotExist(err) {
				t.Fatalf("oracle ran before evidence refusal: %v", err)
			}
		})
	}
}

func TestReviewCheckpointReuse(t *testing.T) {
	f := checkpointFixture(t)
	f.Complete()
	f.Save()
	f.Commit("retain complete evidence")
	if result := Execute(context.Background(), f.Root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("ordinary: %+v", result)
	}
	if code, out := runCheckpoint(t, f); code != 0 {
		t.Fatalf("checkpoint: %d %s", code, out)
	}
	if got := outcomeRuns(t, f.Root); got != 2 {
		t.Fatalf("ordinary green satisfied stronger checkpoint: %d runs", got)
	}
	if code, out := runCheckpoint(t, f); code != 0 || !strings.Contains(out, "reused") {
		t.Fatalf("checkpoint reuse: %d %s", code, out)
	}
	var out, diagnostic bytes.Buffer
	if code := RunCommand([]string{f.Root, "--checkpoint", f.Record.Spec, "--complete"}, &out, &diagnostic); code != 0 {
		t.Fatalf("complete checkpoint: %d %s", code, diagnostic.String())
	}
	if got := outcomeRuns(t, f.Root); got != 3 {
		t.Fatalf("chunk green satisfied complete purpose: %d runs", got)
	}
	f.Record.Chunks[0].Reviews = f.Record.Chunks[0].Reviews[:2]
	f.Save()
	if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "missing Coverage") {
		t.Fatalf("record mutation reused green: %d %s", code, out)
	}
}

func TestReviewCheckpointLaterSource(t *testing.T) {
	for _, committed := range []bool{false, true} {
		t.Run(map[bool]string{false: "dirty source", true: "unreviewed repair"}[committed], func(t *testing.T) {
			f := checkpointFixture(t)
			f.Write("source.txt", "later repair\n")
			if committed {
				f.Commit("unreviewed repair")
			}
			if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "chunk 1: stale reviewed source") {
				t.Fatalf("uncovered delta accepted: %d %s", code, out)
			}
		})
	}
}

func TestReviewCheckpointCanonicalAxes(t *testing.T) {
	// E10 requires this independent expectation; omitting Coverage from Axes must fail it.
	for _, axis := range []string{"Standards", "Spec", "Coverage"} {
		t.Run(axis, func(t *testing.T) {
			f := checkpointFixture(t)
			reviews := f.Record.Chunks[0].Reviews[:0]
			for _, item := range f.Record.Chunks[0].Reviews {
				if item.Axis != axis {
					reviews = append(reviews, item)
				}
			}
			f.Record.Chunks[0].Reviews = reviews
			f.Save()
			if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "missing "+axis) {
				t.Fatalf("canonical axis omitted without refusal: %d %s", code, out)
			}
		})
	}
}

func TestReviewCheckpointOrdinaryWork(t *testing.T) {
	for _, state := range []string{"absent", "active", "pending", "repair", "sibling"} {
		t.Run(state, func(t *testing.T) {
			f := checkpointFixture(t)
			switch state {
			case "absent":
				if err := os.Remove(filepath.Join(f.Root, "reviews/example.md")); err != nil {
					t.Fatal(err)
				}
			case "active":
				f.Record.Chunks = nil
				f.Save()
			case "pending":
				f.Record.Chunks[0].Reviews[0].State = "pending"
				f.Save()
			case "repair":
				f.Write("source.txt", "repair in progress\n")
			case "sibling":
				f.Write("specs/sibling/spec.md", "# Sibling\n\nStatus: staged\n")
			}
			lane, err := RunLane(context.Background(), LaneRequest{Root: f.Root, Tree: benchgit.TreeHash(f.Root), Checks: []Phase{{Name: "ordinary", Argv: []string{"true"}}}})
			if err != nil || !lane.Passed() {
				t.Fatalf("ordinary WIP lane refused: %+v %v", lane, err)
			}
			result := Execute(context.Background(), f.Root, &bytes.Buffer{}, &bytes.Buffer{})
			if result.ActionExit != 0 {
				t.Fatalf("ordinary WIP refused: %+v", result)
			}
		})
	}
}

func TestReviewCheckpointMissingRecord(t *testing.T) {
	f := checkpointFixture(t)
	if err := os.Remove(filepath.Join(f.Root, "reviews/example.md")); err != nil {
		t.Fatal(err)
	}
	if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "retain a valid native result record") {
		t.Fatalf("missing legacy record passed checkpoint: %d %s", code, out)
	}
}

// TestReviewCheckpointRefusalRoute grades the one refusal a single read answers.
// The checkpoint cannot show what the completion evidence lacks, and `bench
// preflight review <slug>` reports the plan row and the record state. So that
// refusal names the read. Every other refusal, and an accepted checkpoint, carry
// no route, because a route that names no state-derived action is a false path.
func TestReviewCheckpointRefusalRoute(t *testing.T) {
	const route = "next=bench preflight review example"
	for _, tc := range []struct {
		name, reason string
		change       func(*recordtest.Fixture)
	}{
		{"missing record", "retain a valid native result record", func(f *recordtest.Fixture) {
			if err := os.Remove(filepath.Join(f.Root, "reviews/example.md")); err != nil {
				t.Fatal(err)
			}
		}},
		{"missing plan", "completion evidence", func(f *recordtest.Fixture) {
			f.Write(recordtest.Spec, "# Example\n\nStatus: staged\n")
			f.Commit("drop the completion plan")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := checkpointFixture(t)
			tc.change(f)
			code, out := runCheckpoint(t, f)
			if code == 0 || !strings.Contains(out, tc.reason) || !strings.Contains(out, route+"\n") {
				t.Fatalf("refusal lost its reason or its route: %d %s", code, out)
			}
		})
	}

	t.Run("accepted", func(t *testing.T) {
		f := checkpointFixture(t)
		code, out := runCheckpoint(t, f)
		if code != 0 || strings.Contains(out, "next=") {
			t.Fatalf("accepted checkpoint carried a route: %d %s", code, out)
		}
	})
}

func TestReviewCheckpointFindingAndReviewIdentity(t *testing.T) {
	for _, kind := range []string{"findings", "source", "pair"} {
		t.Run(kind, func(t *testing.T) {
			f := checkpointFixture(t)
			r := &f.Record.Chunks[0].Reviews[0]
			switch kind {
			case "findings":
				r.Outcome = "findings"
				r.FindingIDs = []string{"unresolved"}
			case "source":
				r.SourceDigest = f.Record.Chunks[0].Base
			case "pair":
				r.Tip = f.Record.Chunks[0].Base
			}
			f.Save()
			if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "Standards") {
				t.Fatalf("invalid axis accepted: %d %s", code, out)
			}
		})
	}
}
