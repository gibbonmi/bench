package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExplicitBaseReviewOwnsSourceRangeNotDestinationHandoff(t *testing.T) {
	root, slug := seedConformant(t)
	base := runGit(t, "rev-parse", "main")
	tip := runGit(t, "rev-parse", "feature")
	runGit(t, "checkout", "-q", "-b", "destination")
	mustWriteFile(t, "capture/session-handoff.md", "destination only\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "destination handoff")
	runGit(t, "checkout", "-q", "feature")
	facts, boot := Gather(root, "review", slug, base)
	if boot != nil {
		t.Fatalf("Gather = %s: %s", boot.Kind, boot.Hint)
	}
	if strings.Join(facts.ChangedPaths, ",") != "internal/example/foo.go" {
		t.Fatalf("source paths = %v, want only source-authored path", facts.ChangedPaths)
	}
	configBefore, err := os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	out, code := Command([]string{"review", slug, "--base", base})
	if code != 0 || !strings.Contains(out, "diff-nonempty,green") || !strings.Contains(out, "source[1]{base,tip}") || !strings.Contains(out, base) || !strings.Contains(out, tip) {
		t.Fatalf("explicit review = (%d):\n%s", code, out)
	}
	configAfter, err := os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	if string(configAfter) != string(configBefore) {
		t.Fatal("accepted explicit preflight changed Git config bytes")
	}
	mustWriteFile(t, "internal/example/staged.go", "package example\n")
	runGit(t, "add", "internal/example/staged.go")
	mustWriteFile(t, "internal/example/foo.go", "package example\n// worktree\n")
	mustWriteFile(t, "internal/example/untracked.go", "package example\n")
	facts, boot = Gather(root, "build", slug, base)
	if boot != nil {
		t.Fatalf("build Gather = %s: %s", boot.Kind, boot.Hint)
	}
	for _, want := range []string{"internal/example/foo.go", "internal/example/staged.go", "internal/example/untracked.go"} {
		if !containsPath(facts.ChangedPaths, want) {
			t.Fatalf("complete source snapshot = %v, missing %s", facts.ChangedPaths, want)
		}
	}
}

func containsPath(paths []string, want string) bool {
	for _, path := range paths {
		if path == want {
			return true
		}
	}
	return false
}

func TestExplicitBasePreflightRefusesDirtyReview(t *testing.T) {
	_, slug := seedConformant(t)
	base := runGit(t, "rev-parse", "main")
	mustWriteFile(t, "dirty.txt", "dirty\n")
	out, code := Command([]string{"review", slug, "--base", base})
	if code != 1 || !strings.Contains(out, "source not clean") {
		t.Fatalf("dirty source = (%d):\n%s", code, out)
	}
}
