// Shared test fixtures for the landing command families: repository fixtures, argument builders, and cross-family helpers.
package worktree

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/landing"
)

func publicLandingFixture(t *testing.T, request, ignored, declaration string) (string, Creation, string, string, string) {
	t.Helper()
	return publicLandingFixtureAtHome(t, request, ignored, declaration, filepath.Join(t.TempDir(), "bench-home"))
}

func publicLandingFixtureAtHome(t *testing.T, request, ignored, declaration, home string) (string, Creation, string, string, string) {
	t.Helper()
	return landingFixtureAtHome(t, request, ignored, declaration, home, true)
}

// specLessLandingFixture is the public landing fixture whose gate grades the landed
// source alone. A spec-less landing publishes no transition, so a gate that demanded
// `Status: implemented` would refuse every spec-less composition for the wrong reason.
func specLessLandingFixture(t *testing.T, request string) (string, Creation, string, string, string) {
	t.Helper()
	return landingFixtureAtHome(t, request, "", "", filepath.Join(t.TempDir(), "bench-home"), false)
}

func landingFixtureAtHome(t *testing.T, request, ignored, declaration, home string, gradeSpec bool) (string, Creation, string, string, string) {
	t.Helper()
	gateSpec, prospectiveSpec := "", ""
	if gradeSpec {
		gateSpec = "IFS= read -r status < specs/x/spec.md\n[ \"$status\" = \"Status: implemented\" ]\n"
		prospectiveSpec = "rg -q '^Status: implemented$' specs/x/spec.md\n"
	}
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", home)
	common := gitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	tally := filepath.Join(common, "bench-land-gate-tally")
	bindEnv(t, "LAND_GATE_TALLY", tally)
	mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate.sh"), []byte("#!/bin/sh\nset -eu\n"+gateSpec+"[ -f owned.txt ]\nprintf g >> \"$LAND_GATE_TALLY\"\n"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate-prospective.sh"), []byte("#!/bin/sh\nset -eu\nruntime=$1\n"+prospectiveSpec+"[ -f owned.txt ]\nprintf g >> \"$LAND_GATE_TALLY\"\n"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate-inputs.json"), []byte("{\"schema\":1,\"closure\":\"local\",\"environment\":[\"LAND_GATE_TALLY\"],\"paths\":[],\"tools\":[]}\n"), 0o644)
	if declaration != "" {
		mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte("{\"schema\":1,\"paths\":[\""+declaration+"\"]}\n"), 0o644)
	}
	ignore := ""
	if ignored != "" {
		ignore = strings.Split(ignored, "/")[0] + "/\n"
		mustWrite(t, filepath.Join(root, ".gitignore"), []byte(ignore), 0o644)
	}
	specBody := "# x\n\nStatus: staged\n\n## User stories\n1. Land source.\n\n### Acceptance coverage map\n| row | story | behavior | seam | why it catches the failure |\n|---|---|---|---|---|\n| LX1 | 1 | lands | command | catches failure |\n\n## Ownership fences\n\n- `owned.txt`\n"
	mustMkdirAll(t, filepath.Join(root, "specs", "x", "tickets"), 0o755)
	mustWrite(t, filepath.Join(root, "specs", "x", "spec.md"), []byte(specBody), 0o644)
	mustWrite(t, filepath.Join(root, "specs", "x", "tickets", "one.md"), []byte("Ticket covers LX1.\n"), 0o644)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "landing base")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	creation := mustCreate(t, root, request, "public landing")
	commitInWorktree(t, creation.Path, "owned.txt", "reviewed bytes\n", "reviewed source")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	if ignored != "" {
		mustMkdirAll(t, filepath.Dir(filepath.Join(creation.Path, filepath.FromSlash(ignored))), 0o755)
		mustWrite(t, filepath.Join(creation.Path, filepath.FromSlash(ignored)), []byte("residue\n"), 0o600)
	}
	return root, creation, base, tip, tally
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
func commitLandingBuildInputs(t *testing.T, root, body string) {
	t.Helper()
	mustMkdirAll(t, filepath.Join(root, "scripts"), 0o755)
	mustWrite(t, filepath.Join(root, "scripts", "go-build.inputs"), []byte(body), 0o644)
	gitRun(t, root, "add", "scripts/go-build.inputs")
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

func stubLandJoins(t *testing.T, base, tip string) func() {
	t.Helper()
	oldLand, oldMarker, oldReconcile, oldRelease, oldAuthorize := landReviewed, advanceLandingMarker, reconcileLanding, releaseLandingAssignment, authorizeLandingSource
	landReviewed = func(context.Context, landing.ReviewedRequest) (landing.ReviewedResult, error) {
		return landing.ReviewedResult{SourceBase: base, SourceTip: tip, DestinationBase: base, Commit: strings.Repeat("a", 40), Tree: strings.Repeat("b", 40)}, nil
	}
	advanceLandingMarker = func(context.Context, string, string, string, string) error { return nil }
	reconcileLanding = func(string, string, string, string) error { return nil }
	releaseLandingAssignment = ReleaseCommand
	authorizeLandingSource = func(string, string, string) (diff.SourceRange, error) {
		return diff.SourceRange{Base: base, Tip: tip}, nil
	}
	return func() {
		landReviewed, advanceLandingMarker, reconcileLanding, releaseLandingAssignment, authorizeLandingSource = oldLand, oldMarker, oldReconcile, oldRelease, oldAuthorize
	}
}

// ticketsOnlyLandingFixture is the spec-less landing fixture with a tickets-only
// `specs/t/` folder committed at the review base and carried into the source. A
// light-path change has exactly this shape: tickets, no spec.md.
func ticketsOnlyLandingFixture(t *testing.T, request string) (string, Creation, string, string, string) {
	t.Helper()
	root, creation, _, _, tally := specLessLandingFixture(t, request)
	mustMkdirAll(t, filepath.Join(root, "specs", "t", "tickets"), 0o755)
	mustWrite(t, filepath.Join(root, "specs", "t", "tickets", "one.md"), []byte("Light path ticket.\n"), 0o644)
	gitRun(t, root, "add", "specs/t")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "tickets-only folder")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, root, "update-ref", "refs/bench/green/main", base)
	gitRun(t, creation.Path, "rebase", "main")
	return root, creation, base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), tally
}

func ticketsOnlyLandArgs(request, base, tip, slug, path string) []string {
	return append([]string{"--spec", slug}, specLessLandArgs(request, base, tip, path)...)
}
