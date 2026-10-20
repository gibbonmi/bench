// Candidate check tests for the gate checkpoint: a declared build's checkpoint refuses before
// any ledger read or write when the running executable is not the installed candidate.
package gate_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
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

// declaredBuild is the recorded first chunk of a plan of chunks chunks whose spec declares a
// persisted field, in the checkout of the assignment of label, which binds the build. The
// primary checkout declares Bench build inputs, and its broker manifest binds an installed
// candidate that is not this test binary. It answers the build fixture and the primary
// checkout.
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
	install(t, primary, freshness.PublishedExecutable(primary))
	return &build, primary
}

// install binds executable in the broker manifest of the primary checkout, as the install's
// doctor repair does.
func install(t *testing.T, primary, executable string) {
	t.Helper()
	dir := filepath.Join(primary, "bin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, _, err := brokermanifest.Write(dir, executable, "0.0.0"); err != nil {
		t.Fatal(err)
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
		carry: map[int]func(*testing.T, gateSet){0: func(t *testing.T, _ gateSet) {
			running, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			install(t, primary, running)
		}},
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
