package commit

// The step rows call refreshWorktreeBuild over a committed kit tree with an injected build,
// and the commit rows run the command over a kit variant of the launch fixture.

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/treetarget/kittest"
)

// rebuiltLine is the whole stderr of a step that rebuilt the worktree build.
const rebuiltLine = "worktree build: rebuilt dist/bench\n"

// secondInput is a second Go file of the kit's cmd/bench package, so a commit can delete a
// sealed build input.
const secondInput = "cmd/bench/second.go"

// The commits that a step row makes after the kit commit: a change of the listed build
// input, a change of a Markdown file alone, and the deletion of secondInput.
const inputCommit, markdownCommit, deletionCommit = "input", "markdown", "deletion"

// writeKit writes the kit tree with secondInput. The tree ignores dist/, as the kit does, so
// a fixture commit tracks no build output.
func writeKit(t *testing.T, root string) {
	t.Helper()
	kittest.WriteTree(t, root)
	mustWrite(t, filepath.Join(root, filepath.FromSlash(secondInput)), "package main\n", 0o644)
	mustWrite(t, filepath.Join(root, ".gitignore"), "/dist/\n", 0o644)
}

// staleKit writes the kit tree with a published build that went stale before the next
// commit, as a build does after a Go edit.
func staleKit(t *testing.T, root string) {
	t.Helper()
	writeKit(t, root)
	publishBuild(t, root, freshness.PublishedExecutable(root))
	kittest.EditBuildInput(t, root)
}

// publishBuild publishes a sealed script of the tree at root at executable.
func publishBuild(t *testing.T, root, executable string) {
	t.Helper()
	staged := filepath.Join(t.TempDir(), "bench")
	mustWrite(t, staged, "#!/bin/sh\nexit 0\n", 0o755)
	if err := freshness.Publish(root, staged, executable, filepath.Dir(executable), "fixture"); err != nil {
		t.Fatal(err)
	}
}

// kitRepo is a repository at a path with a space whose last commit holds the kit tree.
func kitRepo(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "kit root")
	mustMkdirAll(t, root)
	initializeLandingRepo(t, root, 0)
	writeKit(t, root)
	commitmenttest.Commit(t, root, "kit tree")
	return root
}

// commitChange makes and commits one change of kind in root.
func commitChange(t *testing.T, root, kind string) {
	t.Helper()
	switch kind {
	case inputCommit:
		kittest.EditBuildInput(t, root)
	case markdownCommit:
		mustWrite(t, filepath.Join(root, "README.md"), "# Notes\n", 0o644)
	case deletionCommit:
		if err := os.Remove(filepath.Join(root, filepath.FromSlash(secondInput))); err != nil {
			t.Fatal(err)
		}
	}
	commitmenttest.Commit(t, root, kind)
}

// fakeBuild is an injected build. Each call records the HEAD of its root, publishes a
// sealed build at its output when publish is set, and answers err.
type fakeBuild struct {
	t       *testing.T
	publish bool
	err     error
	heads   []string
}

func (b *fakeBuild) build(_ context.Context, root, executable string) error {
	b.heads = append(b.heads, gittest.Output(b.t, root, "rev-parse", "HEAD"))
	if b.publish {
		publishBuild(b.t, root, executable)
	}
	return b.err
}

// runStep runs the step over the HEAD of root with build under a private home. It returns
// that home and the stderr of the step.
func runStep(t *testing.T, root string, build func(context.Context, string, string) error) (home, stderr string) {
	t.Helper()
	home = t.TempDir()
	var out strings.Builder
	refreshWorktreeBuild(root, gittest.Output(t, root, "rev-parse", "HEAD"), home, build, &out)
	return home, out.String()
}

// rebuildRow is one step row: whether a build is published after the kit commit, the
// commits that follow, and the count of builds that the step runs. A row that builds once
// prints the rebuilt line and leaves a build that verifies; a row that builds nothing
// prints nothing.
type rebuildRow struct {
	published bool
	commits   []string
	builds    int
}

func requireRebuild(t *testing.T, row rebuildRow) {
	t.Helper()
	root := kitRepo(t)
	if row.published {
		publishBuild(t, root, freshness.PublishedExecutable(root))
	}
	for _, kind := range row.commits {
		commitChange(t, root, kind)
	}
	build := &fakeBuild{t: t, publish: true}
	_, stderr := runStep(t, root, build.build)
	want := ""
	if row.builds == 1 {
		want = rebuiltLine
		if err := freshness.Verify(root, freshness.PublishedExecutable(root)); err != nil {
			t.Fatalf("after the rebuild Verify = %v, want nil", err)
		}
	}
	if len(build.heads) != row.builds || stderr != want {
		t.Fatalf("step = (%d builds, %q), want (%d builds, %q)", len(build.heads), stderr, row.builds, want)
	}
}

// TestWorktreeRebuildRepairsAStaleBuild is CW13 to CW15.
func TestWorktreeRebuildRepairsAStaleBuild(t *testing.T) {
	requireRebuild(t, rebuildRow{published: true, commits: []string{inputCommit}, builds: 1})
}

// TestWorktreeRebuildBuildsAnAbsentBuildForAnInputChange is CW16.
func TestWorktreeRebuildBuildsAnAbsentBuildForAnInputChange(t *testing.T) {
	requireRebuild(t, rebuildRow{commits: []string{inputCommit}, builds: 1})
}

// TestWorktreeRebuildBuildsAnAbsentBuildForADeletedInput is CW43.
func TestWorktreeRebuildBuildsAnAbsentBuildForADeletedInput(t *testing.T) {
	requireRebuild(t, rebuildRow{commits: []string{deletionCommit}, builds: 1})
}

// TestWorktreeRebuildSkipsAFreshBuild is CW17 and CW18.
func TestWorktreeRebuildSkipsAFreshBuild(t *testing.T) {
	requireRebuild(t, rebuildRow{published: true, commits: []string{markdownCommit}})
}

// TestWorktreeRebuildSkipsAMarkdownCommitWithoutABuild is CW19.
func TestWorktreeRebuildSkipsAMarkdownCommitWithoutABuild(t *testing.T) {
	requireRebuild(t, rebuildRow{commits: []string{markdownCommit}})
}

// TestWorktreeRebuildRepairsAStaleBuildAfterAMarkdownCommit is CW20.
func TestWorktreeRebuildRepairsAStaleBuildAfterAMarkdownCommit(t *testing.T) {
	requireRebuild(t, rebuildRow{published: true, commits: []string{inputCommit, markdownCommit}, builds: 1})
}

// TestWorktreeRebuildRepairsADeletedInput is CW21.
func TestWorktreeRebuildRepairsADeletedInput(t *testing.T) {
	requireRebuild(t, rebuildRow{published: true, commits: []string{deletionCommit}, builds: 1})
}

// TestWorktreeRebuildReportsAFailure is CW22 to CW24 and CW50. Each row builds a stale
// build once in a root whose path holds a space. A failure prints one line with the remedy
// once, and the second seal read decides whether the build failed.
func TestWorktreeRebuildReportsAFailure(t *testing.T) {
	boom := errors.New("boom")
	for _, row := range []struct {
		name  string
		build fakeBuild
		cause string
	}{
		{name: "build error", build: fakeBuild{err: boom}, cause: "boom"},
		{name: "control bytes", build: fakeBuild{err: errors.New("boom\n\x1b[31m")}, cause: `boom\n\u001b[31m`},
		{name: "seal refusal", build: fakeBuild{}, cause: "seal source digest does not match current build inputs"},
		{name: "published then failed", build: fakeBuild{publish: true, err: boom}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := kitRepo(t)
			publishBuild(t, root, freshness.PublishedExecutable(root))
			commitChange(t, root, inputCommit)
			build := row.build
			build.t = t
			home, stderr := runStep(t, root, build.build)
			want := rebuiltLine
			if row.cause != "" {
				cause := row.cause
				if row.build.err == nil {
					cause = "bench binary " + strconv.Quote(freshness.PublishedExecutable(root)) + " is untrusted: " + cause
				}
				want = "worktree build: failed: " + cause + "; rebuild with " + freshness.RebuildIn(home, root).Action + "\n"
			}
			if len(build.heads) != 1 || stderr != want {
				t.Fatalf("step = (%d builds, %q), want (1 build, %q)", len(build.heads), stderr, want)
			}
		})
	}
}

// pathWithoutGo returns a PATH that holds git and no Go toolchain.
func pathWithoutGo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if path, err := exec.LookPath("git"); err != nil || os.Symlink(path, filepath.Join(dir, "git")) != nil {
		t.Fatalf("link git into %s: %v", dir, err)
	}
	return dir
}

// TestWorktreeRebuildNamesAnAbsentGo is CW25. The package build, with no Go on PATH, fails
// with a line that names the absent Go.
func TestWorktreeRebuildNamesAnAbsentGo(t *testing.T) {
	root := kitRepo(t)
	publishBuild(t, root, freshness.PublishedExecutable(root))
	commitChange(t, root, inputCommit)
	t.Setenv("PATH", pathWithoutGo(t))
	_, stderr := runStep(t, root, runbinary.BuildSubject)
	if !strings.HasPrefix(stderr, "worktree build: failed: ") || !strings.Contains(stderr, "Go is absent from PATH") || strings.Count(stderr, "\n") != 1 {
		t.Fatalf("step stderr = %q, want one failure line that names the absent Go", stderr)
	}
}

// TestWorktreeRebuildSkipsALinkedRepository is CW26. A repository with no build-input
// manifest takes no build, even with a dist/bench that no seal authenticates.
func TestWorktreeRebuildSkipsALinkedRepository(t *testing.T) {
	root := filepath.Join(t.TempDir(), "linked")
	mustMkdirAll(t, filepath.Join(root, "dist"))
	initializeLandingRepo(t, root, 0)
	mustWrite(t, freshness.PublishedExecutable(root), "#!/bin/sh\nexit 0\n", 0o755)
	build := &fakeBuild{t: t, publish: true}
	if _, stderr := runStep(t, root, build.build); len(build.heads) != 0 || stderr != "" {
		t.Fatalf("step = (%d builds, %q), want no build and no output", len(build.heads), stderr)
	}
}

// TestWorktreeRebuildIgnoresAHeldInstall is CW45. The ledger holds an open declared build,
// and the worktree's own broker manifest binds its stale build, as an install leaves it.
// The step still rebuilds that build.
func TestWorktreeRebuildIgnoresAHeldInstall(t *testing.T) {
	root, _ := landingRepo(t, 0, staleKit)
	primary, err := git.PrimaryCheckout(root)
	if err != nil {
		t.Fatal(err)
	}
	const declared = "specs/declared/spec.md"
	commitmenttest.Write(t, primary, declared, "# Declared\n\nStatus: staged\nPersisted fields: intent-ledger:commitment.fixture\n")
	commitmenttest.Commit(t, primary, "declare a persisted field")
	identity, err := commitrepo.SourceIdentity(root, "main", declared)
	if err != nil {
		t.Fatal(err)
	}
	// The start already claimed the bound outcome, so the rebind keeps that claim.
	commitmenttest.EditCommitment(t, root, func(state *intent.CommitmentState) {
		state.Bindings[0].Deliverable, state.Bindings[0].Identity = declared, identity
		claim := intent.OutcomeClaim{Milestone: state.Bindings[0].Milestone, Outcome: state.Bindings[0].Outcome}
		if !slices.Contains(state.Claims, claim) {
			state.Claims = append(state.Claims, claim)
		}
	})
	executable := freshness.PublishedExecutable(root)
	commitmenttest.BindManifest(t, root, executable)
	if _, open, err := (commitrepo.Store{Root: root}).DeclaredBuild(); err != nil || !open {
		t.Fatalf("declared build = (%t, %v), want an open declared build", open, err)
	}
	fields, err := brokermanifest.Read(brokermanifest.Path(root))
	digest, digestErr := brokermanifest.Digest(executable)
	if err != nil || digestErr != nil || fields[brokermanifest.DigestField] != digest {
		t.Fatalf("manifest = (%v, %v), digest = (%q, %v), want the manifest to bind %s", fields, err, digest, digestErr, executable)
	}
	build := &fakeBuild{t: t, publish: true}
	if _, stderr := runStep(t, root, build.build); len(build.heads) != 1 || stderr != rebuiltLine {
		t.Fatalf("step = (%d builds, %q), want (1 build, %q)", len(build.heads), stderr, rebuiltLine)
	}
}

// kitLaunch is the launch fixture over a stale kit tree, marked as the kit, whose commits
// run build as the worktree build.
func kitLaunch(t *testing.T, check string, build *fakeBuild) launchFixture {
	t.Helper()
	f := newLaunchFixtureWith(t, check, staleKit)
	f.asKit(t)
	swapWorktreeBuild(t, build)
	return f
}

// swapWorktreeBuild makes build the worktree build of each commit until t ends.
func swapWorktreeBuild(t *testing.T, build *fakeBuild) {
	t.Helper()
	old := worktreeBuild
	t.Cleanup(func() { worktreeBuild = old })
	worktreeBuild = build.build
}

// TestCommitRebuildsBeforeItsGrade is CW27 to CW29 and CW36. The build runs once at the
// published commit, its line comes before the synchronous outcome line of the grade, stdout
// keeps the commit record last, and the commit still starts one grade.
func TestCommitRebuildsBeforeItsGrade(t *testing.T) {
	build := &fakeBuild{t: t, publish: true}
	f := kitLaunch(t, passingCheck, build)
	f.seedRecord(t)
	before := head(t, f.root)
	code, stdout, stderr := runCommand(t, f.root, "-m", "m", "tracked.txt")
	published := head(t, f.root)
	f.launches.Await(t)
	f.launches.End(t)
	rebuilt, outcome := strings.Index(stderr, rebuiltLine), strings.Index(stderr, "background grade: latest ")
	if code != 0 || rebuilt < 0 || outcome < rebuilt {
		t.Fatalf("commit = (%d, %q), want exit 0 and the rebuilt line before the background outcome", code, stderr)
	}
	if published == before || !slices.Equal(build.heads, []string{published}) {
		t.Fatalf("build heads = %q, want one build at the published %s", build.heads, published)
	}
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	if !strings.HasPrefix(lines[len(lines)-1], "committed{published_commit=") || strings.Contains(stdout, "worktree build") {
		t.Fatalf("stdout = %q, want the commit record last and no worktree build text", stdout)
	}
	if got := f.launches.Launched(t); len(got) != 1 {
		t.Fatalf("launches = %+v, want one grade", got)
	}
}

// TestCommitKeepsExitZeroOnAFailedRebuild is CW30.
func TestCommitKeepsExitZeroOnAFailedRebuild(t *testing.T) {
	build := &fakeBuild{t: t, err: errors.New("boom")}
	f := kitLaunch(t, passingCheck, build)
	if code, stdout, stderr := runCommand(t, f.root, "-m", "m", "tracked.txt"); code != 0 || len(build.heads) != 1 {
		t.Fatalf("commit = (%d, %q, %q) after %d builds, want exit 0 after one failed build", code, stdout, stderr, len(build.heads))
	}
}

// TestUnpublishedCommitRebuildsNothing is CW31. A lane refusal and a dry run publish
// nothing, so each runs no build over the stale build.
func TestUnpublishedCommitRebuildsNothing(t *testing.T) {
	for _, row := range []struct {
		name, check string
		args        []string
		code        int
	}{
		{name: "lane refusal", check: failingCheck, args: []string{"-m", "m", "tracked.txt"}, code: 1},
		{name: "dry run", check: passingCheck, args: []string{"--dry-run", "-m", "m", "tracked.txt"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			build := &fakeBuild{t: t, publish: true}
			f := kitLaunch(t, row.check, build)
			if code, stdout, stderr := runCommand(t, f.root, row.args...); code != row.code || len(build.heads) != 0 {
				t.Fatalf("commit = (%d, %q, %q) after %d builds, want exit %d and no build", code, stdout, stderr, len(build.heads), row.code)
			}
		})
	}
}

// TestUnreconciledCommitRebuildsNothing is CW32. A commit that exits 3 over a stale build
// runs no build.
func TestUnreconciledCommitRebuildsNothing(t *testing.T) {
	build := &fakeBuild{t: t, publish: true}
	swapWorktreeBuild(t, build)
	publishedUnreconciledIn(t, "named", staleKit)
	if len(build.heads) != 0 {
		t.Fatalf("build heads = %q, want no build", build.heads)
	}
}

// TestKitIgnoresTheRebuildOutputs is CW33. Each output of a rebuild is ignored at the module
// root, so a rebuild leaves git status clean.
func TestKitIgnoresTheRebuildOutputs(t *testing.T) {
	for _, path := range []string{"dist/bench", "dist/bench.seal", "bin/bench-broker.manifest"} {
		if _, err := gittest.Run(filepath.Join("..", ".."), "check-ignore", "-q", "--", path); err != nil {
			t.Errorf("git check-ignore %s = %v, want ignored", path, err)
		}
	}
}
