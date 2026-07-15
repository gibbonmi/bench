package surface

import (
	"github.com/gibbonmi/bench/internal/contract"
	"path/filepath"
	"strings"
	"testing"
)

func testManagedPrePushGatePinning(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("pin base")
	base := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	writeGatePin(t, f, base)

	unchanged := runPrePush(t, f, refLine("refs/heads/topic", base, "refs/heads/topic"))
	unchanged.RequireExit(0)
	unchanged.RequireNotContains(unchanged.Stderr, "blocked:")

	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 1\n")
	f.CommitAll("drift gate")
	drifted := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)

	blocked := runPrePush(t, f, refLine("refs/heads/topic", drifted, "refs/heads/topic"))
	if blocked.ExitCode == 0 {
		t.Fatalf("pinned pre-push allowed a drifted .bench tree\nstderr:\n%s", blocked.Stderr)
	}
	blocked.RequireContains(blocked.Stderr, "bench gate pin")

	cleanThenDrift := runPrePush(t, f,
		refLine("refs/heads/clean", base, "refs/heads/clean")+
			refLine("refs/heads/drifted", drifted, "refs/heads/drifted"))
	if cleanThenDrift.ExitCode == 0 {
		t.Fatal("pinned pre-push allowed a multi-ref push with a later drifted ref")
	}
	driftThenClean := runPrePush(t, f,
		refLine("refs/heads/drifted", drifted, "refs/heads/drifted")+
			refLine("refs/heads/clean", base, "refs/heads/clean"))
	if driftThenClean.ExitCode == 0 {
		t.Fatal("pinned pre-push allowed a multi-ref push with an earlier drifted ref")
	}

	deleteRef := runPrePush(t, f, refLine("refs/heads/topic", zeroOID(), "refs/heads/topic"))
	deleteRef.RequireExit(0)

	defaultBranch := runPrePush(t, f, refLine("refs/heads/main", base, "refs/heads/main"))
	if defaultBranch.ExitCode == 0 {
		t.Fatal("pre-push no longer blocks direct pushes to main")
	}
	defaultBranch.RequireContains(defaultBranch.Stderr, "direct push to main")
}

func TestManagedPrePushGatePinningEdges(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "managed pre-push unpinned notice contract failed", testManagedPrePushUnpinnedNotice)
	contract.RunParallel(t, "managed pre-push missing .bench tree contract failed", testManagedPrePushMissingBenchTree)
}

func testManagedPrePushUnpinnedNotice(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("base")
	head := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)

	unpinned := runPrePush(t, f, refLine("refs/heads/topic", head, "refs/heads/topic"))
	unpinned.RequireExit(0)
	unpinned.RequireContains(unpinned.Stderr, "bench gate pin")

	contract.WriteFileAbs(t, filepath.Join(gitDirPath(t, f), "bench-gate-pin"), "\n")
	malformed := runPrePush(t, f, refLine("refs/heads/topic", head, "refs/heads/topic"))
	malformed.RequireExit(0)
	malformed.RequireContains(malformed.Stderr, "bench gate pin")
}

func testManagedPrePushMissingBenchTree(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("pin base")
	base := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	writeGatePin(t, f, base)

	contract.Remove(t, filepath.Join(f.Root, ".bench"))
	f.CommitAll("remove bench")
	noBench := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	blocked := runPrePush(t, f, refLine("refs/heads/topic", noBench, "refs/heads/topic"))
	if blocked.ExitCode == 0 {
		t.Fatal("pinned pre-push allowed a commit with no .bench tree")
	}
	blocked.RequireContains(blocked.Stderr, "bench gate pin")
}

// TestManagedPrePushLiveDefaultBranch covers the pre-push guard's default-branch
// resolution: a repo linked before its remote existed baked a fabricated default, so the
// hook resolves origin/HEAD live at push time and falls back to the baked token only when
// that is unresolvable.
func TestManagedPrePushLiveDefaultBranch(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "pre-push did not guard the live default branch", testManagedPrePushLiveDefault)
	contract.RunParallel(t, "pre-push offline baked-fallback guard failed", testManagedPrePushBakedFallback)
}

// testManagedPrePushLiveDefault links a repo with no remote (the hook bakes the fabricated
// "main" default), then points a later origin/HEAD at master: the guard must block the live
// default (master) and no longer the stale baked token (main).
func testManagedPrePushLiveDefault(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("base")
	base := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)

	f.Git("symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/master")

	blocked := runPrePush(t, f, refLine("refs/heads/master", base, "refs/heads/master"))
	if blocked.ExitCode == 0 {
		t.Fatalf("pre-push allowed a direct push to the live default branch master\nstderr:\n%s", blocked.Stderr)
	}
	blocked.RequireContains(blocked.Stderr, "direct push to master")

	// Live resolution supersedes the baked token: main is no longer the default, so a
	// push to it is allowed.
	allowed := runPrePush(t, f, refLine("refs/heads/main", base, "refs/heads/main"))
	allowed.RequireExit(0)
	allowed.RequireNotContains(allowed.Stderr, "blocked:")
}

// testManagedPrePushBakedFallback is the offline-path regression guard: with no remote,
// origin/HEAD is unresolvable, so the hook must fall back to the baked default and still
// block a direct push to it. It passes before and after the live-resolution change.
func testManagedPrePushBakedFallback(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("base")
	base := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)

	blocked := runPrePush(t, f, refLine("refs/heads/main", base, "refs/heads/main"))
	if blocked.ExitCode == 0 {
		t.Fatalf("pre-push allowed a direct push to the baked default with no remote\nstderr:\n%s", blocked.Stderr)
	}
	blocked.RequireContains(blocked.Stderr, "direct push to main")
}

// TestManagedPrePushNewlineLessFinalLine covers the read loop's newline-tail admission:
// git always LF-terminates pre-push stdin, but a hand-crafted or non-LF-terminated final
// ref line must still be checked, not silently dropped by a while-read loop that discards
// $line at EOF.
func TestManagedPrePushNewlineLessFinalLine(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "pre-push dropped a newline-less final ref line", testManagedPrePushNewlineLessFinalLine)
}

func testManagedPrePushNewlineLessFinalLine(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("base")
	base := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)

	// The push has a single ref line with no trailing newline: a while-read loop that
	// discards $line when read returns non-zero at EOF drops this line entirely, so the
	// direct-push-to-main check never runs and the offending push is wrongly allowed.
	line := strings.TrimSuffix(refLine("refs/heads/main", base, "refs/heads/main"), "\n")
	blocked := runPrePush(t, f, line)
	if blocked.ExitCode == 0 {
		t.Fatalf("pre-push allowed a direct push to main via a newline-less final ref line\nstderr:\n%s", blocked.Stderr)
	}
	blocked.RequireContains(blocked.Stderr, "direct push to main")
}

func runPrePush(t *testing.T, f contract.Fixture, stdin string) contract.Probe {
	t.Helper()
	return contract.RunAtWithInput(t, f, f.Root, nil, stdin, "bash", prePushPath(t, f))
}

func prePushPath(t *testing.T, f contract.Fixture) string {
	t.Helper()
	hooks := strings.TrimSpace(f.Git("rev-parse", "--git-path", "hooks").Stdout)
	return filepath.Join(f.Root, filepath.FromSlash(hooks), "pre-push")
}

func writeGatePin(t *testing.T, f contract.Fixture, commit string) {
	t.Helper()
	tree := strings.TrimSpace(f.Git("rev-parse", commit+":.bench").Stdout)
	contract.WriteFileAbs(t, filepath.Join(gitDirPath(t, f), "bench-gate-pin"), tree+"\n"+commit+"\n2026-07-06T00:00:00Z\n")
}

func gitDirPath(t *testing.T, f contract.Fixture) string {
	t.Helper()
	return strings.TrimSpace(f.Git("rev-parse", "--absolute-git-dir").Stdout)
}

func refLine(localRef, localOID, remoteRef string) string {
	return localRef + " " + localOID + " " + remoteRef + " " + zeroOID() + "\n"
}

func zeroOID() string {
	return "0000000000000000000000000000000000000000"
}
