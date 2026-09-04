package gitguard

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
)

// TestClassifyRealCheckerResolvedComposition composes Classify with the real Checker over
// temp repositories. It pins the resolved-answer corner of the polarity matrix that the
// constant refYes/refNo fakes in verdict_test.go cannot reach. That corner is real ref
// resolution, real branch existence, and a real push destination, not a hand-set boolean.
// The checkout rows run in the one-commit fixture; each push row sets its own repository
// state. The probes resolve in the process working directory. So this test t.Chdir's into
// each fixture repo and must not run parallel.
func TestClassifyRealCheckerResolvedComposition(t *testing.T) {
	root := gc1Repo(t)
	t.Chdir(root)
	checker := realChecker()

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

	for _, c := range pushCompositionCases() {
		t.Run(c.name, func(t *testing.T) {
			t.Chdir(c.setup(t))
			if got := Classify(c.cmd, checker); got != c.want {
				t.Errorf("Classify(%q) = %q, want %q", c.cmd, got, c.want)
			}
		})
	}
}

// pushCompositionCases holds the push rows of the resolved composition: each one sets a
// real repository state, then reads the verdict the real facts compose. The rows cover
// the config escape hatch (PG20), the bare rule under each push.default mode (PG21 to
// PG27), a named remote with no refspec (PG29), a directory outside a repository (PG30),
// and the literal HEAD refspec on a detached head (PG40).
func pushCompositionCases() []struct {
	name  string
	setup func(t *testing.T) string
	cmd   string
	want  string
} {
	const (
		toDefault  = "git push to the default branch"
		unresolved = "git push with an unresolved destination"
	)
	onMain := func(t *testing.T) string {
		root := gcTopicRepo(t)
		runGit(t, root, "checkout", "-q", "main")
		return root
	}
	detached := func(t *testing.T) string {
		root := gcTopicRepo(t)
		runGit(t, root, "checkout", "-q", "--detach", "HEAD")
		return root
	}
	mode := func(name string) func(t *testing.T) string {
		return func(t *testing.T) string {
			root := gcTopicRepo(t)
			runGit(t, root, "config", "push.default", name)
			return root
		}
	}
	return []struct {
		name  string
		setup func(t *testing.T) string
		cmd   string
		want  string
	}{
		{"allowProtectedPush does not lift the default-branch rule", func(t *testing.T) string {
			root := onMain(t)
			runGit(t, root, "config", "bench.allowProtectedPush", "true")
			return root
		}, "git push origin main", toDefault},
		{"bare push on a topic branch with push.default unset allows", gcTopicRepo, "git push", ""},
		{"bare push on a topic branch under current allows", mode("current"), "git push", ""},
		{"bare push on the default branch blocks", onMain, "git push", toDefault},
		{"bare push under upstream reads the upstream branch", mode("upstream"), "git push", toDefault},
		{"bare push on a detached head is unresolved", detached, "git push", unresolved},
		{"bare push under matching is unresolved", mode("matching"), "git push", unresolved},
		{"bare push under nothing is unresolved", mode("nothing"), "git push", unresolved},
		{"named remote with no refspec obeys the bare rule", onMain, "git push origin", toDefault},
		{"push outside a repository is unresolved", func(t *testing.T) string { return t.TempDir() }, "git push", unresolved},
		{"HEAD refspec on a detached head is unresolved", detached, "git push origin HEAD", unresolved},
	}
}

// TestClassifyRealCheckerTimeoutComposition composes Classify with the real Checker
// against a PATH-front stub `git` that sleeps past the 2s refCheckTimeout bound
// (internal/git.refCheckTimeout). It pins that both probes' opposite fail-safe
// defaults land on "block" under composition. RefResolves times out to false (an
// unresolvable-looking ref blocks checkout), and BranchExists times out to true (an
// undeterminable branch is presumed present, blocking forced creation). The bare-push row
// composes the same hung git through the destination probes, which carry no bound of
// their own: each call returns empty output, and the destination fact reports no branch,
// so a hung git still denies. The test tolerates the wall-clock cost of the sleeping
// stub. It uses t.Chdir and no t.Parallel, for the same process-cwd reason as the
// resolved-composition test above.
func TestClassifyRealCheckerTimeoutComposition(t *testing.T) {
	stubDir := t.TempDir()
	stubGit(t, stubDir)
	t.Setenv("PATH", stubDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Chdir(t.TempDir())
	checker := realChecker()

	cases := []struct {
		name string
		cmd  string
		want string
	}{
		{"hung ref probe blocks checkout (ref presumed unresolvable)", "git checkout not-a-real-ref", "git checkout path"},
		{"hung branch probe blocks forced creation (branch presumed present)", "git checkout -B some-branch", "git checkout path"},
		{"hung bare-push probe leaves the destination unresolved", "git push", "git push with an unresolved destination"},
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

// realChecker composes the Checker the guard subcommand wires, over the process working
// directory. It mirrors cmd/bench's wiring, which is the composition under test here: a
// fact left nil would report no name and turn every push row into the unresolved label.
func realChecker() Checker {
	return Checker{
		RefResolves:     git.RefResolves,
		BranchExists:    git.BranchExists,
		DefaultBranch:   func() (string, bool) { return git.ResolvedDefault(".") },
		CheckedOut:      func() (string, bool) { return git.CheckedOutName(".") },
		BareDestination: func() (string, bool) { return git.BarePushDestination(".") },
	}
}

// gcTopicRepo builds a temp repo whose checked-out branch is "topic" and whose topic
// branch tracks origin/main, so git.ResolvedDefault resolves "main" and the bare-push
// destination has both a checked-out answer and an upstream answer. The remote-tracking
// ref is written by hand, so the fixture needs no second repository and no network.
func gcTopicRepo(t *testing.T) string {
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
	runGit(t, root, "checkout", "-q", "-b", "topic")
	runGit(t, root, "remote", "add", "origin", "https://example.invalid/r.git")
	runGit(t, root, "update-ref", "refs/remotes/origin/main", "HEAD")
	runGit(t, root, "config", "branch.topic.remote", "origin")
	runGit(t, root, "config", "branch.topic.merge", "refs/heads/main")
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
