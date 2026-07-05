package contract

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAXIFailClosedContracts(t *testing.T) {
	t.Parallel()
	skipIfSubjectBenchMissing(t)
	runParallel(t, "block-dangerous-git analyzer-missing fail-closed contract", testAXIBlockDangerousGitAnalyzerMissing)
	runParallel(t, "block-dangerous-git binary-missing fail-closed (git-shaped)", testAXIBlockDangerousGitBinaryMissing)
	runParallel(t, "block-dangerous-git core-errored fail-closed contract", testAXIBlockDangerousGitCoreErrored)
}

func testAXIBlockDangerousGitAnalyzerMissing(t *testing.T) {
	f := NewFixture(t)
	hook := filepath.Join(SubjectRoot(t), ".bench", "hooks", "block-dangerous-git.sh")
	pathEnv := map[string]string{"PATH": "/usr/bin:/bin"}

	describe := runFixtureCommand(t, f, f.Root, pathEnv, "", bashPath(t), hook, "--describe")
	describe.RequireExit(0)
	requireAXILine(t, describe.Stdout, "denies: manifest unavailable (analyzer missing)")

	enforce := runFixtureCommand(t, f, f.Root, pathEnv, `{"tool_input":{"command":"git push"}}`, bashPath(t), hook)
	enforce.RequireExit(2)
	enforce.RequireContains(enforce.Stderr, "BLOCKED")
}

func testAXIBlockDangerousGitBinaryMissing(t *testing.T) {
	f := NewFixture(t)
	hook := filepath.Join(SubjectRoot(t), ".bench", "hooks", "block-dangerous-git.sh")
	stubBin := filepath.Join(f.Root, "stubbin")
	if err := os.MkdirAll(stubBin, 0o755); err != nil {
		t.Fatalf("create stubbin: %v", err)
	}
	writeExecutableAt(t, stubBin, "bench", "#!/usr/bin/env bash\nexit 127\n")
	pathEnv := map[string]string{"PATH": stubBin + ":/usr/bin:/bin"}

	gitShaped := runFixtureCommand(t, f, f.Root, pathEnv, `{"tool_input":{"command":"git push"}}`, bashPath(t), hook)
	gitShaped.RequireExit(2)
	gitShaped.RequireContains(gitShaped.Stderr, "BLOCKED")

	nonGit := runFixtureCommand(t, f, f.Root, pathEnv, `{"tool_input":{"command":"ls -la"}}`, bashPath(t), hook)
	nonGit.RequireExit(0)
}

func testAXIBlockDangerousGitCoreErrored(t *testing.T) {
	f := NewFixture(t)
	hook := filepath.Join(SubjectRoot(t), ".bench", "hooks", "block-dangerous-git.sh")
	stubBin := filepath.Join(f.Root, "stubbin")
	if err := os.MkdirAll(stubBin, 0o755); err != nil {
		t.Fatalf("create stubbin: %v", err)
	}
	writeExecutableAt(t, stubBin, "bench", "#!/usr/bin/env bash\nexit 3\n")

	enforce := runFixtureCommand(t, f, f.Root, map[string]string{"PATH": stubBin + ":/usr/bin:/bin"}, `{"tool_input":{"command":"git push"}}`, bashPath(t), hook)
	enforce.RequireExit(2)
	enforce.RequireContains(enforce.Stderr, "BLOCKED")
}
