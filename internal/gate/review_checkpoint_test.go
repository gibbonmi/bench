package gate

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

func checkpointFixture(t *testing.T) *recordtest.Fixture {
	return attachedCheckpointFixture(t, recordtest.Attach)
}

// attachedCheckpointFixture builds one recorded chunk from either record form.
// The version 1 and delegated checkpoints share this one build sequence, so a
// change to it cannot drift between them.
func attachedCheckpointFixture(t *testing.T, attach func(testing.TB, string, int) *recordtest.Fixture, prepare ...func(*recordtest.Fixture)) *recordtest.Fixture {
	t.Helper()
	f := attach(t, outcomeFixture(t), 1)
	for _, apply := range prepare {
		apply(f)
	}
	f.AddChunk()
	f.Save()
	f.Commit("retain review evidence")
	return f
}

// The checkpoint route walk runs in the external test package, because the steps it runs
// verbatim reach packages that import this one. These names hand it the checkpoint
// fixtures and the one run seam that no public entry reaches.
var (
	AttachedCheckpointFixture = attachedCheckpointFixture
	RetainCompletion          = retainCompletion
)

// RunCheckpointArmed runs the checkpoint that args select at root, and calls arm once the run
// holds its execution lock: after it accepts its subject and before it validates that
// subject. It answers the exit.
func RunCheckpointArmed(t *testing.T, root string, args []string, arm func(), stderr io.Writer) int {
	t.Helper()
	_, mode, checkpoint, err := parseGateArgs(args, false)
	if err != nil {
		t.Fatal(err)
	}
	armed := func(ctx context.Context) (context.Context, func()) { arm(); return ctx, func() {} }
	return executeAfterAcquire(asVerb(WithCheckpoint(context.Background(), checkpoint)), root, io.Discard, stderr, armed, mode).ActionExit
}

// RunCompletionTree grades tree as the gate verb grades a complete checkpoint of spec at tip,
// so a fixture can grade a tree that is not the transform of that source. It answers the exit.
func RunCompletionTree(root, spec, tip, tree string, stderr io.Writer) int {
	ctx := asVerb(WithCompletion(context.Background(), spec, tip))
	return executeTreeWithOwner(ctx, root, tree, io.Discard, stderr, nil, reuseFreshGreen, nil).ActionExit
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
			if _, err := os.Stat(runCountWitness.path(t, f.Root)); !os.IsNotExist(err) {
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

// TestReviewCheckpointRefusalRoute: the checkpoint cannot show what the completion
// evidence lacks, and `bench preflight review <slug>` reports the plan row and the record
// state. So that refusal names the read, and only the read: the fixed write-access help
// row names no cause of this refusal. An accepted checkpoint carries no route.
func TestReviewCheckpointRefusalRoute(t *testing.T) {
	const route = "next=bench preflight review 'example'"
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
			if strings.Contains(out, "help[1]{cmd,why}") {
				t.Fatalf("refusal printed a help row beside its route: %s", out)
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

// TestCheckpointRouteRendersHostileFacts drives the funnel over a slug, a spec path, and a
// chunk id that hold a space or a byte that is not line-safe. A value with a space prints
// shell-quoted, and a value that is not line-safe prints its slot, so no raw value reaches
// the route line.
func TestCheckpointRouteRendersHostileFacts(t *testing.T) {
	const spaced, unsafe = "specs/a b/spec.md", "specs/a\u0085b/spec.md"
	evidence, capture := evidenceError{errors.New("evidence")}, errors.New("capture")
	quote := sanitize.ShellQuote
	for _, tc := range []struct {
		name       string
		cause      error
		checkpoint Checkpoint
		want       string
	}{
		{"spaced slug", evidence, Checkpoint{Spec: spaced, Chunk: "1"}, "bench preflight review " + quote("a b")},
		{"unsafe slug", evidence, Checkpoint{Spec: unsafe, Chunk: "1"}, "bench preflight review <" + refusalroute.FactSlug + ">"},
		{"spaced spec path and chunk", capture, Checkpoint{Spec: spaced, Chunk: "c 1"}, " --checkpoint " + quote(spaced) + " --chunk " + quote("c 1")},
		{"unsafe spec path and chunk", capture, Checkpoint{Spec: unsafe, Chunk: "c\u00851"}, " --checkpoint <spec-path> --chunk <id>"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := outcomeFixture(t)
			var stderr bytes.Buffer
			ctx := asVerb(WithCheckpoint(context.Background(), tc.checkpoint))
			executeSubjectWithRunBinary(ctx, root, root, io.Discard, &stderr, nil, reuseFreshGreen, failedAcceptEvaluation{err: tc.cause}, nil, "")
			if next, one := routetest.Next(stderr.String()); !one || !strings.HasSuffix(next, tc.want) || !sanitize.LineSafe(next) {
				t.Fatalf("stderr = %q, want one line-safe next= route that ends with %q", stderr.String(), tc.want)
			}
		})
	}
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
