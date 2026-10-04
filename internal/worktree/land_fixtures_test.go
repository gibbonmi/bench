// Shared test fixtures for the landing command families: repository fixtures, argument builders, and cross-family helpers.
package worktree

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testrepo"
)

func landingGateFixture(t *testing.T, environment ...string) *testrepo.GateFixture {
	t.Helper()
	return testrepo.NewGateFixture(t.TempDir(), environment...)
}

// publicLandingFixture mints one private Bench home and carries it in the value, so the
// fixture binds no process environment and the test it serves stays parallel-eligible. The
// value's call builds a verb call at that home. processHomeCall lands at the process home
// instead.
func publicLandingFixture(t *testing.T, request, ignored, declaration string) landingFixture {
	t.Helper()
	return publicLandingFixtureAtHome(t, request, ignored, declaration, filepath.Join(t.TempDir(), "bench-home"))
}

func publicLandingFixtureAtHome(t *testing.T, request, ignored, declaration, home string) landingFixture {
	t.Helper()
	return landingFixtureAtHome(t, request, ignored, declaration, home, true)
}

// specLessLandingFixture is the public landing fixture whose gate grades the landed
// source alone. A spec-less landing publishes no transition, so a gate that demanded
// `Status: implemented` would refuse every spec-less composition for the wrong reason.
func specLessLandingFixture(t *testing.T, request string) landingFixture {
	t.Helper()
	return landingFixtureAtHome(t, request, "", "", filepath.Join(t.TempDir(), "bench-home"), false)
}

// foldedLandingFixture builds a foldedLanding from the spec-less landing fixture.
func foldedLandingFixture(t *testing.T, request string) foldedLanding {
	t.Helper()
	f := specLessLandingFixture(t, request)
	mustWrite(t, filepath.Join(f.root, "advance.txt"), []byte("destination advance\n"), 0o644)
	gitRun(t, f.root, "add", "advance.txt")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "destination advance")
	f.base = gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local",
		"merge", "-q", "--no-ff", "-m", "fold the destination", f.base)
	fold := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, f.creation.Path, "owned.txt", "reviewed bytes after the fold\n", "work after the fold")
	f.tip = gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	return foldedLanding{landingFixture: f, fold: fold}
}

func landingFixtureAtHome(t *testing.T, request, ignored, declaration, home string, gradeSpec bool) landingFixture {
	t.Helper()
	return landingFixtureWithGateStep(t, request, ignored, declaration, home, gradeSpec, nil)
}

// landingFixtureWithGateStep is landingFixtureAtHome with one more prospective gate line.
// A non-nil step receives the gate fixture, so the step declares each command it runs,
// and the destination root. Its line runs after the tally.
func landingFixtureWithGateStep(t *testing.T, request, ignored, declaration, home string, gradeSpec bool, step func(*testrepo.GateFixture, string) string) landingFixture {
	t.Helper()
	gateSpec, prospectiveSpec := "", ""
	f := landingGateFixture(t)
	if gradeSpec {
		gateSpec = "IFS= read -r status < specs/x/spec.md\n[ \"$status\" = \"Status: implemented\" ]\n"
		// The two scripts assert the same status by different means, so a landing that ran
		// the wrong one goes red. Both means stay POSIX, because a bare runner has no rg.
		prospectiveSpec = f.Command("grep") + " -q '^Status: implemented$' specs/x/spec.md\n"
	}
	root := newWorktreeRepo(t)
	common := gitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	tally := filepath.Join(common, "bench-land-gate-tally")
	// The tally path is a literal in both gate scripts, so the fixture binds no name
	// into the process environment and declares an empty gate environment. A gate that
	// read the path from an exported name would make every caller serial.
	count := "printf g >> '" + tally + "'\n"
	extra := ""
	if step != nil {
		extra = step(f, root)
	}
	f.MustWrite(t, root, "set -eu\n"+gateSpec+"[ -f owned.txt ]\n"+count, "set -eu\nruntime=$1\n"+prospectiveSpec+"[ -f owned.txt ]\n"+count+extra)
	if declaration != "" {
		mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte("{\"schema\":1,\"paths\":[\""+declaration+"\"]}\n"), 0o644)
	}
	ignore := ""
	if ignored != "" {
		ignore = strings.Split(ignored, "/")[0] + "/\n"
		mustWrite(t, filepath.Join(root, ".gitignore"), []byte(ignore), 0o644)
	}
	specBody := "# x\n\nStatus: staged\n\n## User stories\n1. Land source.\n\n### Acceptance coverage map\n| row | story | behavior | seam | why it catches the failure |\n|---|---|---|---|---|\n| E1 | 1 | lands | command | catches failure |\n\n## Ownership fences\n\n- `owned.txt`\n- `reviews/x.md`\n- `" + siblingReviewPath + "`\n"
	prepared := recordtest.Prepare(t, root, 1, "specs/x/spec.md", specBody)
	// The ticket writes the fence less its review pickup, so a landing that grades
	// fence-writes reads a union-exact spec.
	prepared.SetTicketWrites("owned.txt (new), " + siblingReviewPath + " (new)")
	commitmenttest.SeedAdmission(t, root, "specs/x/spec.md")
	prepared.Commit("approve fixture delivery")
	base := prepared.Tip()
	creation := mustCreate(t, root, home, request, "public landing")
	commitmenttest.Admit(t, creation.Path, request, "specs/x/spec.md")
	commitInWorktree(t, creation.Path, "owned.txt", "reviewed bytes\n", "reviewed source")
	refreshLandingEvidence(t, creation.Path, base)
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	if ignored != "" {
		mustMkdirAll(t, filepath.Dir(filepath.Join(creation.Path, filepath.FromSlash(ignored))), 0o755)
		mustWrite(t, filepath.Join(creation.Path, filepath.FromSlash(ignored)), []byte("residue\n"), 0o600)
	}
	return landingFixture{ownedAssignment: ownedAssignment{repoHome: repoHome{root, home}, creation: creation}, base: base, tip: tip, tally: tally}
}

// blockLandingReconcile plants a nested repository in the destination at root, so the
// real residue guard fails the landing's reconcile step after the publication. The
// returned repair removes that repository, so a resume then reconciles.
func blockLandingReconcile(t *testing.T, root string) (repair func()) {
	t.Helper()
	nested := plantNestedRepository(t, root)
	return func() { mustRemove(t, nested) }
}

// postPublicationFault is one landing step that runs after the publication, with the
// fixture composition that makes that step fail. name is the step name that the
// interrupted landing reports in its worktree cell. build receives the joins that the
// landing would run under. It returns the fixture, the joins for the interrupted
// landing, and the repair that lets a resume finish the step.
type postPublicationFault struct {
	name  string
	build func(t *testing.T, request string, j joins) (f landingFixture, broken joins, repair func())
}

// postPublicationFaults holds one row for each step after the publication. The marker
// and reconcile faults are real destination states. The release fault breaks its seam
// in the copy of the joins that build returns, so the caller's joins stay whole.
var postPublicationFaults = []postPublicationFault{
	{name: "marker", build: func(t *testing.T, request string, j joins) (landingFixture, joins, func()) {
		return markerLandingFixture(t, request, true), j, func() {}
	}},
	{name: "reconcile", build: func(t *testing.T, request string, j joins) (landingFixture, joins, func()) {
		f := publicLandingFixture(t, request, "", "")
		return f, j, blockLandingReconcile(t, f.root)
	}},
	{name: "release", build: func(t *testing.T, request string, j joins) (landingFixture, joins, func()) {
		j.releaseLandingAssignment = func(joins, ambient, string, []string, io.Writer, io.Writer) int { return 1 }
		return publicLandingFixture(t, request, "", ""), j, func() {}
	}},
}

// cannedGreenShape is what a bounded green run prints: the phase table, the skip count,
// and the verdict. A landing fixture prints these bytes itself, so the journey grades the
// relay rather than the engine that would have produced them.
const cannedGreenShape = "phases[1]{phase,verdict,elapsed_ms}:\n  build,green,12\ncapability-skips: 0 (capability=0 environment=0)\ngate: green\n"

// commitCannedShapeGate replaces the fixture's prospective gate with one that prints a
// canned shape and settles green. The landing runs the prospective script, so that is the
// one script the relay carries.
func commitCannedShapeGate(t *testing.T, root, shape string) (base string) {
	t.Helper()
	script := "#!/bin/sh\nprintf '%s' " + sanitize.ShellQuote(shape) + "\nexit 0\n"
	mustWrite(t, filepath.Join(root, ".bench", "gate-prospective.sh"), []byte(script), 0o755)
	gitRun(t, root, "add", ".bench/gate-prospective.sh")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "canned bounded gate")
	return gitOutput(t, root, "rev-parse", "HEAD")
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}

// commitLandingBuildInputs commits the Go build-input manifest whose presence puts a
// landing in the dev context, where the command has to prove its own executable. It is
// committed rather than dropped in place so the destination stays clean. The freshness
// refusal is then the only thing a landing can fail on.
func commitLandingBuildInputs(t *testing.T, root string) {
	t.Helper()
	manifest := filepath.Join(root, filepath.FromSlash(freshness.BuildInputsManifest))
	mustMkdirAll(t, filepath.Dir(manifest), 0o755)
	mustWrite(t, manifest, []byte(freshness.BuildInputLine("build_script", "scripts/go-build.sh")), 0o644)
	gitRun(t, root, "add", freshness.BuildInputsManifest)
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "declare go build inputs")
}

func landArgs(request, base, tip, path string) []string {
	return []string{"--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land", path}
}

func specLessLandArgs(request, base, tip, path string) []string {
	return []string{"--request", request, "--base", base, "--source-tip", tip, "-m", "land", path}
}

func stageLandSpec(t *testing.T, root, source string) {
	t.Helper()
	path := filepath.Join(root, "specs", "x", "spec.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, path, []byte("Status: staged\n"), 0o644)
	gitRun(t, root, "add", "specs/x/spec.md")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "stage spec")
	gitRun(t, source, "rebase", "main")
}

// ticketsOnlyLandingFixture is the spec-less landing fixture with a tickets-only
// `specs/t/` folder committed at the review base and carried into the source. A
// light-path change has exactly this shape: tickets, no spec.md.
func ticketsOnlyLandingFixture(t *testing.T, request string) landingFixture {
	t.Helper()
	return ticketsOnlyFolderFixture(t, specLessLandingFixture(t, request))
}

// ticketsOnlyFolderFixture adds the tickets-only folder to a spec-less landing fixture.
func ticketsOnlyFolderFixture(t *testing.T, f landingFixture) landingFixture {
	t.Helper()
	mustMkdirAll(t, filepath.Join(f.root, "specs", "t", "tickets"), 0o755)
	mustWrite(t, filepath.Join(f.root, "specs", "t", "tickets", "one.md"), []byte("Light path ticket.\n"), 0o644)
	gitRun(t, f.root, "add", "specs/t")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "tickets-only folder")
	f.base = gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.root, "update-ref", "refs/bench/green/main", f.base)
	gitRun(t, f.creation.Path, "rebase", "main")
	f.tip = gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	return f
}

func ticketsOnlyLandArgs(request, base, tip, slug, path string) []string {
	return append([]string{"--spec", slug}, specLessLandArgs(request, base, tip, path)...)
}

// withFenceEntry adds one fence entry to a spec the landing fixture prepared.
func withFenceEntry(spec []byte, entry string) []byte { return recordtest.WithFenceEntry(spec, entry) }
