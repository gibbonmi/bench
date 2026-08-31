package worktree

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
)

// buildRecorder is the `build` join a row drives instead of a real compile. It records
// every call and writes body at the output path, so a row can grade the arguments, the
// artifact, or the failure without a Go toolchain.
type buildRecorder struct {
	calls  [][2]string
	body   string
	result error
}

func (r *buildRecorder) join(_ context.Context, worktree, output string) error {
	r.calls = append(r.calls, [2]string{worktree, output})
	if r.result != nil {
		return r.result
	}
	if r.body != "" {
		if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
			return err
		}
		return os.WriteFile(output, []byte(r.body), 0o755)
	}
	return nil
}

// buildJoins is the seam set a row drives, with the recorder standing in for the build.
func buildJoins(recorder *buildRecorder) joins {
	j := defaultJoins()
	j.build = recorder.join
	return j
}

// plantBuildScript writes the stub the production join executes. The stub writes marker
// to its second argument, which is the output path the verb passes, so a row reads the
// script's own bytes back out of `dist/bench`.
func plantBuildScript(t *testing.T, worktree, marker string) {
	t.Helper()
	mustMkdirAll(t, filepath.Join(worktree, "scripts"), 0o755)
	body := "#!/usr/bin/env bash\nset -eu\nmkdir -p \"$(dirname \"$2\")\"\nprintf '%s' '" + marker + "' > \"$2\"\n"
	mustWrite(t, filepath.Join(worktree, "scripts", "go-build.sh"), []byte(body), 0o755)
}

// builtExecutable returns what the build left at the worktree's `dist/bench`.
func builtExecutable(t *testing.T, worktree string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(worktree, "dist", "bench"))
	mustNoError(t, err)
	return string(data)
}

// relabelAssignment gives one assignment a label the create grammar refuses, so a row can
// grade how the next line renders a hostile label. The Bench lock carries a digest of the
// label, so the rewrite relocks the checkout; without that step the creation bundle
// refuses the target before the build runs.
func relabelAssignment(t *testing.T, root string, assignment intent.Assignment, label string) intent.Assignment {
	t.Helper()
	assignment.Label = label
	mustNoError(t, intent.PutAssignment(root, assignment))
	mustNoError(t, unlockWorktree(root, assignment.Worktree))
	mustNoError(t, lockWorktree(root, assignment.Worktree, lockReason(assignment)))
	return assignment
}

// pathWithoutGo returns a PATH that holds git and no Go toolchain. Git stays reachable
// because the target resolves through the ledger and the checkout before the build seam
// is reached at all.
func pathWithoutGo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	git, err := exec.LookPath("git")
	mustNoError(t, err)
	mustNoError(t, os.Symlink(git, filepath.Join(dir, "git")))
	return dir
}

// WF1: the verb hands the build seam the worktree it resolved and that worktree's own
// `dist/bench`, so no caller has to name either path.
func TestBuildCallsTheJoinWithTheWorktreeAndOutput(t *testing.T) {
	t.Parallel()
	root, creation, _ := newOwnedAssignment(t, "build-join-arguments")
	recorder := &buildRecorder{}
	var stdout, stderr bytes.Buffer
	code := buildWith(buildJoins(recorder), root, Home(), []string{creation.Assignment.Label}, &stdout, &stderr)
	requireTest(t, code == 0, "build exit = %d, stderr %q", code, stderr.String())
	want := [2]string{creation.Assignment.Worktree, filepath.Join(creation.Assignment.Worktree, "dist", "bench")}
	requireTest(t, len(recorder.calls) == 1, "build join calls = %d, want 1", len(recorder.calls))
	requireTest(t, recorder.calls[0] == want, "build join call = %v, want %v", recorder.calls[0], want)
}

// WF2: the production join reaches the worktree's own build script, so the executable and
// its seal come from the sanctioned producer rather than from a bare `go build`.
func TestBuildRunsTheWorktreeBuildScript(t *testing.T) {
	t.Parallel()
	root, creation, _ := newOwnedAssignment(t, "build-runs-script")
	plantBuildScript(t, creation.Path, "script-authored")
	var stdout, stderr bytes.Buffer
	code := buildWith(defaultJoins(), root, Home(), []string{creation.Assignment.Label}, &stdout, &stderr)
	requireTest(t, code == 0, "build exit = %d, stderr %q", code, stderr.String())
	requireTest(t, builtExecutable(t, creation.Path) == "script-authored", "dist/bench = %q, want the script's bytes", builtExecutable(t, creation.Path))
}

// WF3: success names the assignment and the absolute executable, so the reader never has
// to derive either from the target they typed.
func TestBuildPrintsTheTable(t *testing.T) {
	t.Parallel()
	root, creation, _ := newOwnedAssignment(t, "build-prints-table")
	var stdout, stderr bytes.Buffer
	code := buildWith(buildJoins(&buildRecorder{}), root, Home(), []string{creation.Assignment.Label}, &stdout, &stderr)
	requireTest(t, code == 0, "build exit = %d, stderr %q", code, stderr.String())
	want, err := toon.Table("worktree_build", []string{"worktree", "executable"}, [][]string{{creation.Assignment.ID, filepath.Join(creation.Assignment.Worktree, "dist", "bench")}})
	mustNoError(t, err)
	requireTest(t, strings.HasPrefix(stdout.String(), want), "build printed %q, want the table %q", stdout.String(), want)
}

// WF4: the next line is paste-safe for every label. A line-safe label is shell-quoted, and
// a label holding a control byte gives way to the assignment id, which never holds one.
func TestBuildNamesTheExecFormForTheLabel(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name    string
		label   string
		address func(intent.Assignment) string
	}{
		{name: "quoted and globbed", label: "it's a*b", address: func(intent.Assignment) string { return `'it'\''s a*b'` }},
		{name: "newline", label: "two\nlines", address: func(a intent.Assignment) string { return a.ID }},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			root, creation, _ := newOwnedAssignment(t, "build-next-"+strings.ReplaceAll(row.name, " ", "-"))
			assignment := relabelAssignment(t, root, creation.Assignment, row.label)
			var stdout, stderr bytes.Buffer
			code := buildWith(buildJoins(&buildRecorder{}), root, Home(), []string{assignment.ID}, &stdout, &stderr)
			requireTest(t, code == 0, "build exit = %d, stderr %q", code, stderr.String())
			want := "next[1]:\n  bench worktree exec " + row.address(assignment) + " -- ./dist/bench <verb>\n"
			requireTest(t, strings.HasSuffix(stdout.String(), want), "build printed %q, want it to end with %q", stdout.String(), want)
		})
	}
}

// WF5: a build failure prints the builder's own sentence and then the tree it ran in, so
// the reader never needs a raw path lookup to act on the failure.
func TestBuildFailureNamesTheWorktree(t *testing.T) {
	t.Parallel()
	root, creation, _ := newOwnedAssignment(t, "build-failure-names")
	recorder := &buildRecorder{result: errors.New("build script exited 1")}
	var stdout, stderr bytes.Buffer
	code := buildWith(buildJoins(recorder), root, Home(), []string{creation.Assignment.Label}, &stdout, &stderr)
	requireTest(t, code == 1, "build exit = %d, want 1", code)
	want := "bench worktree build: build script exited 1\nworktree: " + creation.Assignment.Worktree + "\n"
	requireTest(t, stderr.String() == want, "build printed %q, want %q", stderr.String(), want)
}

// WF6: an absent Go toolchain is refused by the builder's own sentence, so the refusal
// names the tool rather than an exec failure the reader has to decode.
func TestBuildRefusesWithoutGoOnPath(t *testing.T) {
	root, creation, _ := newOwnedAssignment(t, "build-without-go")
	plantBuildScript(t, creation.Path, "never-runs")
	bindEnv(t, "PATH", pathWithoutGo(t))
	var stdout, stderr bytes.Buffer
	code := buildWith(defaultJoins(), root, Home(), []string{creation.Assignment.Label}, &stdout, &stderr)
	requireTest(t, code == 1, "build exit = %d, want 1", code)
	requireTest(t, strings.Contains(stderr.String(), "Go is absent from PATH"), "build printed %q, want the builder's Go sentence", stderr.String())
	requireTest(t, strings.HasSuffix(stderr.String(), "worktree: "+creation.Assignment.Worktree+"\n"), "build printed %q, want it to end with the worktree line", stderr.String())
}

// WF7: an interrupted build exits 130 and still leaves the tree's path, so a signal reads
// apart from a broken build.
func TestBuildCancelExitsOneHundredThirty(t *testing.T) {
	t.Parallel()
	root, creation, _ := newOwnedAssignment(t, "build-cancelled")
	recorder := &buildRecorder{result: context.Canceled}
	var stdout, stderr bytes.Buffer
	code := buildWith(buildJoins(recorder), root, Home(), []string{creation.Assignment.Label}, &stdout, &stderr)
	requireTest(t, code == 130, "build exit = %d, want 130", code)
	requireTest(t, strings.HasSuffix(stderr.String(), "worktree: "+creation.Assignment.Worktree+"\n"), "build printed %q, want it to end with the worktree line", stderr.String())
}

// WF9: a rebuild replaces the executable in place, so an edit-and-rebuild loop needs no
// clean step between the two runs.
func TestBuildReplacesAPriorExecutable(t *testing.T) {
	t.Parallel()
	root, creation, _ := newOwnedAssignment(t, "build-replaces-prior")
	var stdout, stderr bytes.Buffer
	for _, marker := range []string{"first-build", "second-build"} {
		plantBuildScript(t, creation.Path, marker)
		stdout.Reset()
		stderr.Reset()
		code := buildWith(defaultJoins(), root, Home(), []string{creation.Assignment.Label}, &stdout, &stderr)
		requireTest(t, code == 0, "build for %s exit = %d, stderr %q", marker, code, stderr.String())
	}
	requireTest(t, builtExecutable(t, creation.Path) == "second-build", "dist/bench = %q, want the second build's bytes", builtExecutable(t, creation.Path))
}

// WF11: the build leaves residue under `dist/` alone, so the landing's residue rule stays
// green after a worktree has been built.
func TestBuildWritesOnlyUnderDist(t *testing.T) {
	t.Parallel()
	root, creation, _ := newOwnedAssignment(t, "build-writes-under-dist")
	plantBuildScript(t, creation.Path, "declared-output")
	gitRun(t, creation.Path, "add", "scripts/go-build.sh")
	gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "add the build script")
	var stdout, stderr bytes.Buffer
	code := buildWith(defaultJoins(), root, Home(), []string{creation.Assignment.Label}, &stdout, &stderr)
	requireTest(t, code == 0, "build exit = %d, stderr %q", code, stderr.String())
	listing, err := descendant(t, "git", "-C", creation.Path, "status", "--porcelain", "--untracked-files=all").Output()
	mustNoError(t, err)
	untracked := untrackedPaths(string(listing))
	requireTest(t, len(untracked) > 0, "git reported no untracked path, so the row grades nothing")
	for _, path := range untracked {
		requireTest(t, strings.HasPrefix(path, "dist/"), "the build left %q outside dist/; every untracked path is %v", path, untracked)
	}
}

// untrackedPaths returns the paths a porcelain listing reports as untracked.
func untrackedPaths(listing string) []string {
	var paths []string
	for _, line := range strings.Split(strings.TrimSuffix(listing, "\n"), "\n") {
		if rest, found := strings.CutPrefix(line, "?? "); found {
			paths = append(paths, rest)
		}
	}
	return paths
}
