// Spec-less landing tests: publication, gate refusal, base resolution, and resume without a spec.
package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
)

// WL1, WL2, WL5, and WL10: the spec-less landing publishes and releases like a spec
// landing, keeps both reviewed parents and the marker, transitions no spec, and prints
// the same record.
func TestLandCommandSpecLessLandsPublishesAndReleases(t *testing.T) {
	t.Parallel()
	request := "spec-less-land"
	f := specLessLandingFixture(t, request)
	specsBefore := gitOutput(t, f.creation.Path, "rev-parse", f.tip+":specs")
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base, f.tip, f.creation.Path)...))
	published := gitOutput(t, f.root, "rev-parse", "main")
	tree := gitOutput(t, f.root, "rev-parse", published+"^{tree}")
	want := wantEffects("not-applicable") + "landed{source_base=" + f.base + ",source_tip=" + f.tip + ",destination_base=" + f.base + ",published_commit=" + published + ",tree=" + tree + ",worktree=released,census=0}\n"
	if r.exit != 0 || r.stdout != want {
		t.Fatalf("spec-less land = (%d, %q, %q), want (0, %q)", r.exit, r.stdout, r.stderr, want)
	}
	parents := strings.Fields(gitOutput(t, f.root, "rev-list", "--parents", "-n", "1", published))
	if len(parents) != 3 || parents[1] != f.base || parents[2] != f.tip {
		t.Fatalf("published parents = %q, want destination %s and source %s", parents, f.base, f.tip)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got := gitOutput(t, f.root, "rev-parse", published+":specs"); got != specsBefore {
		t.Fatalf("published specs tree = %s, want the composed %s", got, specsBefore)
	}
	if got := gitOutput(t, f.root, "show", published+":specs/x/spec.md"); strings.Contains(got, "Status: implemented") {
		t.Fatalf("spec-less landing transitioned a spec: %q", got)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
	if _, err := os.Stat(f.creation.Path); !os.IsNotExist(err) {
		t.Fatalf("spec-less landing retained the worktree: %v", err)
	}
}

// WL3 and WL22: the gate still owns the spec-less path, and its refusal exits 1 with
// nothing published.
func TestLandCommandSpecLessGateRefusalPublishesNothing(t *testing.T) {
	t.Parallel()
	request := "spec-less-gate-red"
	f := specLessLandingFixture(t, request)
	mustWrite(t, filepath.Join(f.root, ".bench", "gate-prospective.sh"), []byte("#!/bin/sh\nset -eu\nruntime=$1\nprintf g >> '"+f.tally+"'\nexit 1\n"), 0o755)
	gitRun(t, f.root, "add", ".bench/gate-prospective.sh")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "red prospective gate")
	base := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")

	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.HasPrefix(r.stdout, "refused{detail=prospective authorization refused") {
		t.Fatalf("spec-less gate refusal = (%d, %q, %q), want exit 1 and an authorization refusal", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != base {
		t.Fatalf("gate refusal published main=%s, want %s", got, base)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
}

// WL6: the identity proofs still run without a spec, and the refusal names both sides.
func TestLandCommandSpecLessRefusesSourceTipMismatch(t *testing.T) {
	t.Parallel()
	request := "spec-less-tip-mismatch"
	f := specLessLandingFixture(t, request)
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base, f.base, f.creation.Path)...))
	// The route re-points the caller's own command at the tip the worktree holds, and a
	// spec-less landing re-runs spec-less.
	want := "refused{detail=worktree source tip mismatch,observed=" + f.base + ",wanted=" + f.tip +
		",next=" + landingRerun(request, f.base, f.tip, "", f.creation.Path, f.creation.Assignment.ID) + "}\n"
	if r.exit != 1 || r.stdout != want || len(r.stderr) != 0 {
		t.Fatalf("spec-less tip mismatch = (%d, %q, %q), want (1, %q, empty)", r.exit, r.stdout, r.stderr, want)
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("identity refusal ran the gate: %v", err)
	}
}

// WL23: the range proof still runs without a spec, and it refuses before the gate.
func TestLandCommandSpecLessRefusesNonAncestorBaseBeforeTheGate(t *testing.T) {
	t.Parallel()
	request := "spec-less-nonancestor-base"
	f := specLessLandingFixture(t, request)
	commitInWorktree(t, f.root, "destination-only", "destination\n", "destination movement")
	unrelated := gitOutput(t, f.root, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, unrelated, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "reviewed source range is invalid") || !strings.Contains(r.stdout, "not an ancestor") {
		t.Fatalf("spec-less non-ancestor base = (%d, %q, %q), want an invalid-range refusal", r.exit, r.stdout, r.stderr)
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("range refusal ran the gate: %v", err)
	}
}

// WL24: the landed record carries the resolved review base, never the flag's spelling.
func TestLandCommandSpecLessLandedSourceBaseIsTheResolvedBase(t *testing.T) {
	t.Parallel()
	request := "spec-less-resolved-base"
	f := specLessLandingFixture(t, request)
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base[:12], f.tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "landed{source_base="+f.base+",") || strings.Contains(r.stdout, f.base[:12]+",") {
		t.Fatalf("abbreviated spec-less base = (%d, %q, %q), want the resolved base in the record", r.exit, r.stdout, r.stderr)
	}
}

// WL7: a spec-less landing interrupted after publication resumes without a spec.
func TestResumeLandCommandSpecLessCompletesAnInterruptedLanding(t *testing.T) {
	t.Parallel()
	request := "spec-less-resume"
	f := markerLandingFixture(t, request, false)
	r, published := interruptLandingAtMarker(t, f, specLessLandArgs(request, f.base, f.tip, f.creation.Path)...)
	if strings.Contains(r.stdout, "--spec") {
		t.Fatalf("spec-less resume instruction named a spec: %q", r.stdout)
	}
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, f.creation.Path}
	r = runVerb(t, verbLand, f.call(args...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") || len(r.stderr) != 0 {
		t.Fatalf("spec-less resume = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("spec-less resume reran the gate: tally=%q error=%v", got, err)
	}
}

// WL25: a resume without a spec completes a published spec-backed landing's marker and
// release, and publishes nothing a second time.
func TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding(t *testing.T) {
	t.Parallel()
	request := "spec-backed-spec-less-resume"
	f := markerLandingFixture(t, request, true)
	_, published := interruptLandingAtMarker(t, f, landArgs(request, f.base, f.tip, f.creation.Path)...)
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, f.creation.Path}
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") || len(r.stderr) != 0 {
		t.Fatalf("spec-less resume of a spec landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != published {
		t.Fatalf("resume republished: main=%s, want %s", got, published)
	}
	if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("resume reran the gate: tally=%q error=%v", got, err)
	}
}

// The edge under WL1: the flag is optional, but an empty value stays a usage error.
func TestLandCommandRefusesAnEmptySpecValue(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"first run", []string{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "", "-m", "m", "path"}},
		{"resume", []string{"--resume", "p", "--request", "r", "--base", "b", "--source-tip", "s", "--spec", "", "path"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := runVerb(t, verbLand, verbCall{home: Home(), args: tc.args})
			if r.exit != 2 || len(r.stdout) != 0 || !strings.Contains(r.stderr, `--spec ""`) {
				t.Fatalf("empty --spec = (%d, %q, %q), want (2, empty, a usage line naming the empty value)", r.exit, r.stdout, r.stderr)
			}
		})
	}
}

// planningLandingFixture is one planning assignment with no delivery binding. Its gate
// passes and tallies each run. An adopted fixture commits the protected A then B
// commitment; otherwise the repository has no policy.
func planningLandingFixture(t *testing.T, request string, adopted bool) landingFixture {
	t.Helper()
	root := newWorktreeRepo(t)
	tally, count := landingGateTally(t, root)
	landingGateFixture(t).MustWrite(t, root, "set -eu\n"+count, "set -eu\nruntime=$1\n"+count)
	if adopted {
		commitmenttest.SeedProtected(t, root)
	}
	commitmenttest.Commit(t, root, "planning base")
	home := filepath.Join(t.TempDir(), "bench-home")
	creation := mustCreate(t, root, home, request, "planning")
	return landingFixture{ownedAssignment: ownedAssignment{repoHome: repoHome{root, home}, creation: creation}, base: gitOutput(t, root, "rev-parse", "HEAD"), tally: tally}
}

// benchCommit runs the built commit verb in the planning worktree for paths.
func benchCommit(t *testing.T, worktree string, paths ...string) (int, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmd := descendant(t, testRunBinary(t), append([]string{"commit", "-m", "planning commit", "--"}, paths...)...)
	cmd.Dir, cmd.Stdout, cmd.Stderr = worktree, &stdout, &stderr
	return exitCode(cmd.Run()), stdout.String() + stderr.String()
}

// refusePlanningLanding lands the planning source spec-less and requires a commitment
// refusal that names want before the gate, with the destination and source unchanged.
func refusePlanningLanding(t *testing.T, f landingFixture, request, want string) {
	t.Helper()
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "refused{detail=commitment: ") || !strings.Contains(r.stdout, want) {
		t.Fatalf("planning landing = (%d, %q, %q), want a commitment refusal naming %q", r.exit, r.stdout, r.stderr, want)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != f.base {
		t.Fatalf("refused planning landing published main=%s, want %s", got, f.base)
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("refused planning landing ran the gate: %v", err)
	}
	if got := gitOutput(t, f.creation.Path, "rev-parse", "HEAD"); got != tip {
		t.Fatalf("refused planning landing moved the source to %s, want %s", got, tip)
	}
}

// Before adoption, a planning assignment commits and lands a decision map and a
// staged spec. Adoption cannot require already-admitted planning work.
func TestCommitmentPlanningBootstrap(t *testing.T) {
	t.Parallel()
	f := planningLandingFixture(t, "planning", false)
	mustMkdirAll(t, filepath.Join(f.creation.Path, "decisions"), 0o755)
	mustMkdirAll(t, filepath.Join(f.creation.Path, "specs", "y"), 0o755)
	mustWrite(t, filepath.Join(f.creation.Path, "decisions", "delivery-map.md"), []byte("# Delivery map\n\nDecide the delivery order.\n"), 0o644)
	mustWrite(t, filepath.Join(f.creation.Path, "specs", "y", "spec.md"), []byte("# y\n\nStatus: staged\n"), 0o644)
	if code, out := benchCommit(t, f.creation.Path, "decisions/delivery-map.md", "specs/y/spec.md"); code != 0 {
		t.Fatalf("planning commit = (%d, %q), want exit 0", code, out)
	}
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.call(specLessLandArgs("planning", f.base, tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released") {
		t.Fatalf("planning landing = (%d, %q, %q), want a published release", r.exit, r.stdout, r.stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	if got := gitOutput(t, f.root, "show", published+":specs/y/spec.md"); got != "# y\n\nStatus: staged" {
		t.Fatalf("published staged spec = %q", got)
	}
	if got := gitOutput(t, f.root, "show", published+":decisions/delivery-map.md"); !strings.Contains(got, "Decide the delivery order.") {
		t.Fatalf("published decision map = %q", got)
	}
}

// Before adoption, a planning assignment cannot publish a production file through
// either the commit or a raw-Git source that it then lands.
func TestCommitmentPlanningFence(t *testing.T) {
	t.Parallel()
	f := planningLandingFixture(t, "planning", false)
	before := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	mustWrite(t, filepath.Join(f.creation.Path, "tool.go"), []byte("package tool\n"), 0o644)
	if code, out := benchCommit(t, f.creation.Path, "tool.go"); code != 1 || !strings.Contains(out, "commitment adoption required") {
		t.Fatalf("planning production commit = (%d, %q), want an adoption refusal", code, out)
	}
	if got := gitOutput(t, f.creation.Path, "rev-parse", "HEAD"); got != before {
		t.Fatalf("refused planning commit moved HEAD to %s", got)
	}
	commitInWorktree(t, f.creation.Path, "tool.go", "package tool\n", "raw production commit")
	refusePlanningLanding(t, f, "planning", "commitment adoption required")
}

// A planning landing that renames active A's roadmap row refuses publication.
func TestCommitmentProtectedRename(t *testing.T) {
	t.Parallel()
	f := planningLandingFixture(t, "planning", true)
	board := filepath.Join(f.creation.Path, "ROADMAP.md")
	body, err := os.ReadFile(board)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, board, bytes.Replace(body, []byte("**FT1 — A**"), []byte("**FT9 — A**"), 1), 0o644)
	mustWrite(t, filepath.Join(f.creation.Path, "roadmap", "FT9.md"), []byte("**FT9 — A**\n\nKeep the A obligation.\n"), 0o644)
	mustRemove(t, filepath.Join(f.creation.Path, "roadmap", "FT1.md"))
	gitRun(t, f.creation.Path, "add", "-A")
	gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "rename the protected row")
	refusePlanningLanding(t, f, "planning", "candidate changes protected commitment")
}

// A planning landing that recommends unrelated C first refuses publication.
func TestCommitmentProtectedSequence(t *testing.T) {
	t.Parallel()
	f := planningLandingFixture(t, "planning", true)
	body, err := os.ReadFile(filepath.Join(f.creation.Path, "ROADMAP.md"))
	if err != nil {
		t.Fatal(err)
	}
	commitInWorktree(t, f.creation.Path, "ROADMAP.md", strings.Replace(string(body), "1. A\n2. B", "1. C\n2. A\n3. B", 1), "recommend unrelated work first")
	refusePlanningLanding(t, f, "planning", "protected recommended sequence changed")
}
