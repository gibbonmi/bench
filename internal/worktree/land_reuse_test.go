// The landing and the merge reuse another run's whole green. Every verb here runs in this
// test process, so each run has one running executable, and the gate tally counts the whole
// runs of the fixture's prospective gate.
package worktree

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/gate"
)

// requireOneRun requires that the fixture gate behind tally ran exactly once.
func requireOneRun(t *testing.T, tally, when string) {
	t.Helper()
	requireTally(t, tally, "g", when)
}

// requireTally requires that the fixture gate tally reads want, one g for each run.
func requireTally(t *testing.T, tally, want, when string) {
	t.Helper()
	if got, err := os.ReadFile(tally); err != nil || string(got) != want {
		t.Fatalf("gate tally %s = %q, %v; want %q", when, got, err, want)
	}
}

// completeCheckpointIn runs the complete checkpoint of the fixture spec in source and
// requires it green.
func completeCheckpointIn(t *testing.T, source string) {
	t.Helper()
	var out, diagnostic bytes.Buffer
	if code := gate.RunCommand([]string{source, "--checkpoint", closureSpec, "--complete"}, &out, &diagnostic); code != 0 {
		t.Fatalf("complete checkpoint = %d; stdout=%q stderr=%q", code, out.String(), diagnostic.String())
	}
}

// GR25, GR26: in a fixture with no lane, the merge of the destination grades the merge tip,
// and a spec-less landing of that tip reuses the green and advances the green marker.
func TestLandReusesTheMergeGreen(t *testing.T) {
	t.Parallel()
	request := "land-reuse-merge"
	f := specLessLandingFixture(t, request)
	base := advanceDestination(t, f.root)
	merged := runVerb(t, verbMerge, mergeSetOver(t, f.root, f.home).merge("--from", "main", f.creation.Assignment.ID))
	if merged.exit != 0 {
		t.Fatalf("merge = (%d, %q, %q)", merged.exit, merged.stdout, merged.stderr)
	}
	requireOneRun(t, f.tally, "after the merge")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")

	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 0 {
		t.Fatalf("landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	requireOneRun(t, f.tally, "after the landing of the merge tip")
	if published, marker := gitOutput(t, f.root, "rev-parse", "main"), projectGreenMarker(t, f.root); marker != published {
		t.Fatalf("project-green marker = %q, want the published commit %s", marker, published)
	}
}

// GR66: in a fixture with no lane, a merge whose composed tree another checkout graded green
// adds no whole-gate run.
func TestMergeReusesAGreenOfTheComposedTree(t *testing.T) {
	t.Parallel()
	f := specLessLandingFixture(t, "merge-reuse")
	advanceDestination(t, f.root)
	tree := gitOutput(t, f.root, "merge-tree", "--write-tree", f.tip, "main")
	var stdout, stderr bytes.Buffer
	if result := gate.ExecuteTree(context.Background(), f.root, tree, &stdout, &stderr); result.ActionExit != 0 {
		t.Fatalf("destination run of the composed tree = %+v, stderr=%q", result, stderr.String())
	}

	merged := runVerb(t, verbMerge, mergeSetOver(t, f.root, f.home).merge("--from", "main", f.creation.Assignment.ID))
	if merged.exit != 0 || !strings.Contains(merged.stdout, ",tree="+tree) {
		t.Fatalf("merge = (%d, %q, %q), want the composed tree %s", merged.exit, merged.stdout, merged.stderr, tree)
	}
	requireOneRun(t, f.tally, "after the merge of a graded tree")
}

// GR27: a green complete checkpoint in the source worktree answers the spec landing.
func TestLandReusesTheCompleteCheckpointGreen(t *testing.T) {
	t.Parallel()
	request := "land-reuse-checkpoint"
	f := publicLandingFixture(t, request, "", "")
	completeCheckpointIn(t, f.creation.Path)

	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 0 {
		t.Fatalf("landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	requireOneRun(t, f.tally, "after the landing of the checkpoint tree")
}

// GR28: the run binary member follows the source root of the run binary, so a source that
// changes a Bench build input still reuses its own complete checkpoint green. The destination
// is a Go module that builds its run binary, so each prospective run builds from the graded
// tree. That member reads no kit, so the test keeps the process kit and runs in parallel.
// The build script fails, so the refresh effect after the publication reports failed.
func TestLandReusesACheckpointGreenOfAGoChange(t *testing.T) {
	t.Parallel()
	request := "land-reuse-build-input"
	f := seededLandingFixture(t, request, "", "", t.TempDir(), true, nil, closureAdmission, runBinaryMain)
	commitLandingBuildInputs(t, f.root)
	mustWrite(t, filepath.Join(f.root, "go.mod"), []byte("module fixture\n\ngo 1.21\n"), 0o644)
	for _, path := range []string{landingBuildScript, runBinaryMain} {
		mustMkdirAll(t, filepath.Dir(filepath.Join(f.root, path)), 0o755)
	}
	mustWrite(t, filepath.Join(f.root, landingBuildScript), []byte("#!/bin/sh\nexit 1\n"), 0o755)
	commitInWorktree(t, f.root, runBinaryMain, "package main\n\nfunc main() {}\n", "declare the run binary source")
	gitRun(t, f.root, "add", "go.mod", landingBuildScript)
	gitRun(t, f.root, withIdentity("commit", "-qm", "declare the Go module")...)
	base := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	commitInWorktree(t, f.creation.Path, runBinaryMain, "package main\n\nfunc main() { println() }\n", "change a build input")
	refreshLandingEvidence(t, f.creation.Path, base)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	bindRunningBroker(t, f.root)
	completeCheckpointIn(t, f.creation.Path)

	r := runVerb(t, verbLand, f.call(landArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 3 || !strings.Contains(r.stdout, wantEffects("failed")) || gitOutput(t, f.root, "rev-parse", "main") == base {
		t.Fatalf("landing = (%d, %q, %q), want a publication with a failed refresh", r.exit, r.stdout, r.stderr)
	}
	requireOneRun(t, f.tally, "after the landing of a build-input change")
}

// bindRunningBroker binds this test binary in the broker manifest of the primary checkout at
// root, so a complete checkpoint of a repository with Bench build inputs keeps its own
// process. The checkout excludes the manifest, so it stays clean.
func bindRunningBroker(t *testing.T, root string) {
	t.Helper()
	commitmenttest.BindRunning(t, root)
	exclude := mustAdminPath(t, root, "info/exclude")
	excluded, err := os.ReadFile(exclude)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	manifest, err := filepath.Rel(root, brokermanifest.Path(root))
	mustNoError(t, err)
	mustMkdirAll(t, filepath.Dir(exclude), 0o755)
	mustWrite(t, exclude, append(excluded, "/"+filepath.ToSlash(manifest)+"\n"...), 0o644)
}

// runBinaryMain is the command source of the fixture's run binary, a Bench build input.
const runBinaryMain = "cmd/bench/main.go"
