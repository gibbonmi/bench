// Candidate check tests for the gate checkpoint: a declared build's checkpoint refuses before
// any ledger read or write when the running executable is not the installed candidate.
package gate_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// declaredBody is the spec body of a declared build. The field is one that no ledger type
// holds.
const declaredBody = recordtest.Body + "Persisted fields: intent-ledger:commitment.fixture\n"

// declaredBuild is a recorded first chunk in the checkout of the assignment of label, which
// binds the build. The plan holds as many chunks as the chunks argument states, and its spec
// declares a persisted field. The primary checkout declares Bench build inputs, and its
// broker manifest binds an installed candidate that is not this test binary. It answers the
// build fixture and the primary checkout.
func declaredBuild(t *testing.T, label string, chunks int) (*recordtest.Fixture, string) {
	t.Helper()
	t.Setenv("BENCH_HOME", t.TempDir())
	var primary string
	f := gate.AttachedCheckpointFixture(t, func(tb testing.TB, root string, _ int) *recordtest.Fixture {
		primary = root
		return linkedPrepare(tb, root, chunks, declaredBody)
	})
	build := *f
	build.Root = commitmenttest.BindBuild(t, f.Root, label, recordtest.Spec)
	commitmenttest.Write(t, primary, freshness.BuildInputsManifest, "cmd\n")
	commitmenttest.Write(t, primary, "dist/bench", "installed candidate\n")
	commitmenttest.BindManifest(t, primary, freshness.PublishedExecutable(primary))
	return &build, primary
}

// installRunning binds this test binary in the broker manifest of the primary checkout, as
// the operator's install binds the running executable.
func installRunning(t *testing.T, primary string) {
	t.Helper()
	running, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	commitmenttest.BindManifest(t, primary, running)
}

// openSecondDeclaredBuild approves a second outcome in parallel with the build of the
// fixture, with a spec that declares a persisted field, and starts it in the checkout of a
// new assignment. It answers the second outcome.
func openSecondDeclaredBuild(t *testing.T, primary string) string {
	t.Helper()
	const outcome, path = "second", "specs/second/spec.md"
	commitmenttest.Write(t, primary, path, declaredBody)
	commitmenttest.EditPolicy(t, primary, func(policy *commitment.Policy) {
		source := commitment.SourceBinding{ID: "second-spec", Path: path, Identity: commitment.Identity([]byte(declaredBody))}
		policy.Milestones[0].Outcomes = append(policy.Milestones[0].Outcomes, commitment.Outcome{ID: outcome,
			Criteria: []commitment.Criterion{{ID: "second.done", Text: "The second delivery is accepted."}}, Deliverables: []commitment.DeliveryBinding{{Source: source}}})
		policy.ParallelGrants = []commitment.ParallelGrant{{Outcomes: []string{commitmenttest.DeliveryOutcome, outcome}}}
	})
	commitmenttest.Commit(t, primary, "approve a second declared build")
	checkout := commitmenttest.Assignment(t, primary, "gate-overlap-second")
	if err := (commitrepo.Store{Root: checkout}).Start(outcome, "gate-overlap-second", path); err != nil {
		t.Fatal(err)
	}
	return outcome
}

// overlapFixture is the gate face walk's overlap fixture: the build of the fixture runs the
// installed candidate, and a second declared build is open. The reviewer ends the second
// build, which the fixture carries out by blocking its outcome.
func overlapFixture() gateFaceFixture {
	const label = "gate-overlap"
	var primary, second string
	return gateFaceFixture{
		face:     "declared-builds-overlap",
		args:     chunkArgs,
		contains: routetest.ReviewerMarker + "bench commitment block --outcome ",
		build: func(t *testing.T) gateSet {
			var f *recordtest.Fixture
			f, primary = declaredBuild(t, label, 2)
			installRunning(t, primary)
			second = openSecondDeclaredBuild(t, primary)
			return gateSet{f: f, label: label}
		},
		carry: map[int]func(*testing.T, gateSet){0: func(t *testing.T, _ gateSet) {
			if err := (commitrepo.Store{Root: primary}).Block(second, "one primary executable serves one candidate"); err != nil {
				t.Fatal(err)
			}
		}},
	}
}

// candidateFixture is the gate face walk's candidate fixture: a declared build whose primary
// broker manifest binds an installed candidate that is not this test binary. The operator's
// install binds the running executable, which the fixture carries out by rebinding the
// manifest to this test binary.
func candidateFixture() gateFaceFixture {
	const label = "gate-candidate"
	var primary string
	return gateFaceFixture{
		face:     "candidate-not-installed",
		args:     chunkArgs,
		contains: "bench worktree build --install " + sanitize.ShellQuote(label) + "; then bench gate --in " + sanitize.ShellQuote(label) + " ",
		build: func(t *testing.T) gateSet {
			var f *recordtest.Fixture
			f, primary = declaredBuild(t, label, 2)
			return gateSet{f: f, label: label}
		},
		carry: map[int]func(*testing.T, gateSet){0: func(t *testing.T, _ gateSet) { installRunning(t, primary) }},
	}
}

// refuseUninstalled runs the checkpoint that args select in the checkout of f and requires a
// candidate-not-installed refusal with one route and no oracle output. It answers the
// diagnostic and the printed route.
func refuseUninstalled(t *testing.T, f *recordtest.Fixture, args []string) (string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := gate.RunCommand(append([]string{f.Root}, args...), &stdout, &stderr)
	next, printed := routetest.Next(stderr.String())
	if code == 0 || stdout.Len() != 0 || !printed || !strings.Contains(stderr.String(), "candidate-not-installed") {
		t.Fatalf("checkpoint = (%d, %q, %q), want a candidate-not-installed refusal with one route and no oracle output", code, stdout.String(), stderr.String())
	}
	return stderr.String(), next
}

// TestCheckpointPersistedFields is CB29, CB30, CB32, CB33, CB36, and CB57. The installed
// candidate is not this test binary, so the binary stands in for a version-cache executable
// that runs after an install: it knows every field and is still not the installed bytes. A
// malformed declaration hands back with its own error.
func TestCheckpointPersistedFields(t *testing.T) {
	t.Run("chunk", func(t *testing.T) {
		// The record lacks a review, so a check after the review-record read prints that
		// refusal and its route first.
		f, _ := declaredBuild(t, "build", 2)
		f.Record.Chunks[0].Reviews = f.Record.Chunks[0].Reviews[:2]
		f.Save()
		commitmenttest.RequireLedgerKept(t, f.Root, func() { refuseUninstalled(t, f, chunkArgs) })
	})
	t.Run("complete", func(t *testing.T) {
		f, _ := declaredBuild(t, "build", 1)
		gate.RetainCompletion(f)
		commitmenttest.RequireLedgerKept(t, f.Root, func() { refuseUninstalled(t, f, completeArgs) })
	})
	t.Run("unreadable ledger", func(t *testing.T) {
		f, _ := declaredBuild(t, "build", 2)
		commitmenttest.PlantCommitmentField(t, f.Root, "fixture")
		_, cause := intent.Read(f.Root)
		diagnostic, next := refuseUninstalled(t, f, chunkArgs)
		if cause == nil || strings.Contains(diagnostic, cause.Error()) || !strings.HasPrefix(next, "from the primary checkout, run bench worktree build --install <target>; then ") {
			t.Fatalf("checkpoint over the decode error %v = %q, want the refusal and the install slot, not the decode error", cause, diagnostic)
		}
	})
	t.Run("author worktree", func(t *testing.T) {
		f, _ := declaredBuild(t, "first", 2)
		author := *f
		author.Root = commitmenttest.BindBuild(t, f.Root, "second", recordtest.Spec)
		_, next := refuseUninstalled(t, &author, chunkArgs)
		if want := "from the primary checkout, run bench worktree build --install 'first'; then bench gate --in 'second' --checkpoint 'specs/example/spec.md' --chunk '1'"; next != want {
			t.Fatalf("route from the author worktree:\n got %q\nwant %q", next, want)
		}
	})
	t.Run("malformed declaration", func(t *testing.T) {
		f, _ := declaredBuild(t, "build", 2)
		spec := string(readFile(t, f.Root, recordtest.Spec))
		f.Write(recordtest.Spec, strings.Replace(spec, "intent-ledger:commitment.fixture", "gate-evidence:fixture", 1))
		f.Commit("break the declaration")
		var stdout, stderr bytes.Buffer
		code := gate.RunCommand(append([]string{f.Root}, chunkArgs...), &stdout, &stderr)
		next, _ := routetest.Next(stderr.String())
		if code == 0 || !strings.Contains(stderr.String(), "persisted-fields-invalid") || !strings.HasPrefix(next, routetest.ReviewerMarker) {
			t.Fatalf("checkpoint = (%d, %q), want the persisted-fields-invalid error and a reviewer route", code, stderr.String())
		}
	})
}
