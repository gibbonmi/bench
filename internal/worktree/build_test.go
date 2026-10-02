package worktree

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
)

// buildRecorder is the `buildSubject` join a row drives instead of a real compile. It records
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
	j.buildSubject = recorder.join
	return j
}

// plantBuildScript writes the stub the production join executes. The stub writes marker
// to its last argument, which is the output path the verb passes, so a row reads the
// script's own bytes back out of `dist/bench`. Reading the last argument rather than the
// second keeps the stub honest about the real script's grammar, which accepts options
// ahead of the two positionals. The stub also records every argument, one per line,
// beside the output, so a row can grade the options the production join passed.
func plantBuildScript(t *testing.T, worktree, marker string) {
	t.Helper()
	mustMkdirAll(t, filepath.Join(worktree, "scripts"), 0o755)
	body := "#!/usr/bin/env bash\nset -eu\nout=\"${!#}\"\nmkdir -p \"$(dirname \"$out\")\"\nprintf '%s\\n' \"$@\" > \"$out.argv\"\nprintf '%s' '" + marker + "' > \"$out\"\n"
	mustWrite(t, filepath.Join(worktree, "scripts", "go-build.sh"), []byte(body), 0o755)
}

// builtExecutable returns what the build left at the worktree's `dist/bench`.
func builtExecutable(t *testing.T, worktree string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(worktree, "dist", "bench"))
	mustNoError(t, err)
	return string(data)
}

// buildScriptArgv returns the arguments the planted stub received, in order.
func buildScriptArgv(t *testing.T, worktree string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(worktree, "dist", "bench.argv"))
	mustNoError(t, err)
	return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
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
	f := newOwnedAssignment(t, "build-join-arguments")
	recorder := &buildRecorder{}
	r := runVerb(t, verbBuild, repoHome{f.root, Home()}.callWith(buildJoins(recorder), f.creation.Assignment.Label))
	requireTest(t, r.exit == 0, "build exit = %d, stderr %q", r.exit, r.stderr)
	want := [2]string{f.creation.Assignment.Worktree, filepath.Join(f.creation.Assignment.Worktree, "dist", "bench")}
	requireTest(t, len(recorder.calls) == 1, "build join calls = %d, want 1", len(recorder.calls))
	requireTest(t, recorder.calls[0] == want, "build join call = %v, want %v", recorder.calls[0], want)
}

// WF2: the production join reaches the worktree's own build script, so the executable and
// its seal come from the sanctioned producer rather than from a bare `go build`.
func TestBuildRunsTheWorktreeBuildScript(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "build-runs-script")
	plantBuildScript(t, f.creation.Path, "script-authored")
	r := runVerb(t, verbBuild, repoHome{f.root, Home()}.callWith(defaultJoins(), f.creation.Assignment.Label))
	requireTest(t, r.exit == 0, "build exit = %d, stderr %q", r.exit, r.stderr)
	requireTest(t, builtExecutable(t, f.creation.Path) == "script-authored", "dist/bench = %q, want the script's bytes", builtExecutable(t, f.creation.Path))
}

// FT327: the verb builds the worktree's own published executable, so it names no manifest
// directory. The script then publishes the broker manifest beside the wrapper in `bin/`,
// where the doctor row and the landing read it, and not beside `dist/bench`.
func TestBuildLeavesTheManifestDirectoryToTheScript(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "build-manifest-default")
	plantBuildScript(t, f.creation.Path, "subject-build")
	r := runVerb(t, verbBuild, repoHome{f.root, Home()}.callWith(defaultJoins(), f.creation.Assignment.Label))
	requireTest(t, r.exit == 0, "build exit = %d, stderr %q", r.exit, r.stderr)
	argv := buildScriptArgv(t, f.creation.Path)
	requireTest(t, !slices.Contains(argv, "--manifest-dir"), "the build script received %q, want no --manifest-dir so the manifest lands beside the wrapper", argv)
}

// WF3: success names the assignment and the absolute executable, so the reader never has
// to derive either from the target they typed.
func TestBuildPrintsTheTable(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "build-prints-table")
	r := runVerb(t, verbBuild, repoHome{f.root, Home()}.callWith(buildJoins(&buildRecorder{}), f.creation.Assignment.Label))
	requireTest(t, r.exit == 0, "build exit = %d, stderr %q", r.exit, r.stderr)
	want, err := toon.Table("worktree_build", []string{"worktree", "executable"}, [][]string{{f.creation.Assignment.ID, filepath.Join(f.creation.Assignment.Worktree, "dist", "bench")}})
	mustNoError(t, err)
	requireTest(t, strings.HasPrefix(r.stdout, want), "build printed %q, want the table %q", r.stdout, want)
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
			f := newOwnedAssignment(t, "build-next-"+strings.ReplaceAll(row.name, " ", "-"))
			assignment := relabelAssignment(t, f.root, f.creation.Assignment, row.label)
			r := runVerb(t, verbBuild, repoHome{f.root, Home()}.callWith(buildJoins(&buildRecorder{}), assignment.ID))
			requireTest(t, r.exit == 0, "build exit = %d, stderr %q", r.exit, r.stderr)
			want := "next[1]:\n  bench worktree exec " + row.address(assignment) + " -- ./dist/bench <verb>\n"
			requireTest(t, strings.HasSuffix(r.stdout, want), "build printed %q, want it to end with %q", r.stdout, want)
		})
	}
}

// WF5: a build failure prints the builder's own sentence and then the tree it ran in, so
// the reader never needs a raw path lookup to act on the failure.
func TestBuildFailureNamesTheWorktree(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "build-failure-names")
	recorder := &buildRecorder{result: errors.New("build script exited 1")}
	r := runVerb(t, verbBuild, repoHome{f.root, Home()}.callWith(buildJoins(recorder), f.creation.Assignment.Label))
	requireTest(t, r.exit == 1, "build exit = %d, want 1", r.exit)
	want := "bench worktree build: build script exited 1\nworktree: " + f.creation.Assignment.Worktree + "\n"
	requireTest(t, r.stderr == want, "build printed %q, want %q", r.stderr, want)
}

// WF6: an absent Go toolchain is refused by the builder's own sentence, so the refusal
// names the tool rather than an exec failure the reader has to decode.
func TestBuildRefusesWithoutGoOnPath(t *testing.T) {
	f := newOwnedAssignment(t, "build-without-go")
	plantBuildScript(t, f.creation.Path, "never-runs")
	bindEnv(t, "PATH", pathWithoutGo(t))
	r := runVerb(t, verbBuild, repoHome{f.root, Home()}.callWith(defaultJoins(), f.creation.Assignment.Label))
	requireTest(t, r.exit == 1, "build exit = %d, want 1", r.exit)
	requireTest(t, strings.Contains(r.stderr, "Go is absent from PATH"), "build printed %q, want the builder's Go sentence", r.stderr)
	requireTest(t, strings.HasSuffix(r.stderr, "worktree: "+f.creation.Assignment.Worktree+"\n"), "build printed %q, want it to end with the worktree line", r.stderr)
}

// WF7: an interrupted build exits 130 and still leaves the tree's path, so a signal reads
// apart from a broken build.
func TestBuildCancelExitsOneHundredThirty(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "build-cancelled")
	recorder := &buildRecorder{result: context.Canceled}
	r := runVerb(t, verbBuild, repoHome{f.root, Home()}.callWith(buildJoins(recorder), f.creation.Assignment.Label))
	requireTest(t, r.exit == 130, "build exit = %d, want 130", r.exit)
	requireTest(t, strings.HasSuffix(r.stderr, "worktree: "+f.creation.Assignment.Worktree+"\n"), "build printed %q, want it to end with the worktree line", r.stderr)
}

// WF9: a rebuild replaces the executable in place, so an edit-and-rebuild loop needs no
// clean step between the two runs.
func TestBuildReplacesAPriorExecutable(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "build-replaces-prior")
	for _, marker := range []string{"first-build", "second-build"} {
		plantBuildScript(t, f.creation.Path, marker)
		r := runVerb(t, verbBuild, repoHome{f.root, Home()}.callWith(defaultJoins(), f.creation.Assignment.Label))
		requireTest(t, r.exit == 0, "build for %s exit = %d, stderr %q", marker, r.exit, r.stderr)
	}
	requireTest(t, builtExecutable(t, f.creation.Path) == "second-build", "dist/bench = %q, want the second build's bytes", builtExecutable(t, f.creation.Path))
}

// WF11: the build leaves residue under `dist/` alone, so the landing's residue rule stays
// green after a worktree has been built.
func TestBuildWritesOnlyUnderDist(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "build-writes-under-dist")
	plantBuildScript(t, f.creation.Path, "declared-output")
	gitRun(t, f.creation.Path, "add", "scripts/go-build.sh")
	gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "add the build script")
	r := runVerb(t, verbBuild, repoHome{f.root, Home()}.callWith(defaultJoins(), f.creation.Assignment.Label))
	requireTest(t, r.exit == 0, "build exit = %d, stderr %q", r.exit, r.stderr)
	listing, err := descendant(t, "git", "-C", f.creation.Path, "status", "--porcelain", "--untracked-files=all").Output()
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
