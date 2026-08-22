package gitguard

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
)

// TestClassifyRealCheckerResolvedComposition composes Classify with the real
// Checker{git.RefResolves, git.BranchExists} over a temp repo with one commit and one
// branch. It pins the resolved-answer corner of the polarity matrix that the constant
// refYes/refNo fakes in verdict_test.go cannot reach. That corner is real ref
// resolution and real branch existence, not a hand-set boolean. The probes resolve in
// the process working directory. So this test t.Chdir's into the fixture repo and
// must not run parallel.
func TestClassifyRealCheckerResolvedComposition(t *testing.T) {
	root := gc1Repo(t)
	t.Chdir(root)
	checker := Checker{RefResolves: git.RefResolves, BranchExists: git.BranchExists}

	cases := []struct {
		name string
		cmd  string
		want string
	}{
		{"unresolvable ref blocks checkout", "git checkout not-a-real-ref", "git checkout path"},
		{"resolvable ref permits checkout", "git checkout main", ""},
		{"forced creation clobbering an existing branch blocks", "git checkout -B feature", "git checkout path"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify(c.cmd, checker); got != c.want {
				t.Errorf("Classify(%q) = %q, want %q", c.cmd, got, c.want)
			}
		})
	}
}

// TestClassifyRealCheckerTimeoutComposition composes Classify with the real Checker
// against a PATH-front stub `git` that sleeps past the 2s refCheckTimeout bound
// (internal/git.refCheckTimeout). It pins that both probes' opposite fail-safe
// defaults land on "block" under composition. RefResolves times out to false (an
// unresolvable-looking ref blocks checkout), and BranchExists times out to true (an
// undeterminable branch is presumed present, blocking forced creation). The test
// tolerates the wall-clock cost of two ~2s probe timeouts. It uses t.Chdir and no
// t.Parallel, for the same process-cwd reason as the resolved-composition test above.
func TestClassifyRealCheckerTimeoutComposition(t *testing.T) {
	stubDir := t.TempDir()
	stubGit(t, stubDir)
	t.Setenv("PATH", stubDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Chdir(t.TempDir())
	checker := Checker{RefResolves: git.RefResolves, BranchExists: git.BranchExists}

	cases := []struct {
		name string
		cmd  string
		want string
	}{
		{"hung ref probe blocks checkout (ref presumed unresolvable)", "git checkout not-a-real-ref", "git checkout path"},
		{"hung branch probe blocks forced creation (branch presumed present)", "git checkout -B some-branch", "git checkout path"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify(c.cmd, checker); got != c.want {
				t.Errorf("Classify(%q) = %q, want %q", c.cmd, got, c.want)
			}
		})
	}
}

// gc1Repo builds a temp repo with one commit on "main" and one branch "feature",
// exercising the exec-git-in-test idiom.
func gc1Repo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init", "-q", "-b", "main")
	runGit(t, root, "config", "user.email", "test@example.invalid")
	runGit(t, root, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", ".")
	runGit(t, root, "commit", "-qm", "initial")
	runGit(t, root, "branch", "feature")
	return root
}

func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// stubGit writes a PATH-front `git` shell script that sleeps past refCheckTimeout
// (2s) before doing anything else, standing in for a hung git process.
func stubGit(t *testing.T, dir string) {
	t.Helper()
	script := "#!/bin/sh\nsleep 3\n"
	path := filepath.Join(dir, "git")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
}
