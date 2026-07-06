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

	describe := runPrePushDescribe(t, f)
	describe.RequireExit(0)
	describe.RequireContains(describe.Stdout, ".bench drift from bench gate pin")

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

func runPrePush(t *testing.T, f contract.Fixture, stdin string) contract.Probe {
	t.Helper()
	return contract.RunAtWithInput(t, f, f.Root, nil, stdin, "bash", prePushPath(t, f))
}

func runPrePushDescribe(t *testing.T, f contract.Fixture) contract.Probe {
	t.Helper()
	return contract.RunAtWithInput(t, f, f.Root, nil, "", "bash", prePushPath(t, f), "--describe")
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
