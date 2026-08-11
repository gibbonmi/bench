package preflight

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// --- fixture harness: mirrors internal/coverage's TestCommand and internal/diff's
// review_base_test.go patterns — real git commands in a throwaway repo, exact
// output/exit assertions against the public Command entry point only. ---

func runGit(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func mustWriteFile(t *testing.T, path, body string) {
	t.Helper()
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q): %v", dir, err)
		}
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

// specBody renders a bootstrap-conformant spec for slug: staged status, a valid
// opted-in coverage map declaring row PF1 and PF2, and one backticked fence entry
// authorizing internal/<slug>/.
func specBody(slug string, extraFenceLines ...string) string {
	var b strings.Builder
	b.WriteString("# " + slug + "\n\nStatus: staged\n\n")
	b.WriteString("## User stories\n1. As a, I want b, so c.\n\n")
	b.WriteString("### Acceptance coverage map\n")
	b.WriteString("| row | story | behavior | seam | red signal | why it catches the failure |\n")
	b.WriteString("|---|---|---|---|---|---|\n")
	b.WriteString("| PF1 | 1 | does x | cli seam | cmd fails | catches z |\n")
	b.WriteString("| PF2 | 1 | does y | cli seam | cmd fails | catches w |\n")
	b.WriteString("\n## Ownership fences\n\n")
	b.WriteString("- `internal/" + slug + "/` (implementation)\n")
	for _, line := range extraFenceLines {
		b.WriteString(line + "\n")
	}
	return b.String()
}

// initRepo starts a fresh git repo at t.TempDir(), chdir'd into it, with an author
// identity configured — the common prefix of every seed function below.
func initRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Chdir(root)
	runGit(t, "init", "-q", "-b", "main")
	runGit(t, "config", "user.email", "t@example.com")
	runGit(t, "config", "user.name", "t")
	return root
}

// seedConformant builds the PF1 tracer fixture: a base commit on main carrying a
// bootstrap-conformant spec and a ticket citing both declared rows, then a feature
// branch with one authorized change — the tree every check should answer green over.
func seedConformant(t *testing.T) (root, slug string) {
	t.Helper()
	slug = "example"
	root = initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", "Ticket citing PF1 and PF2.\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	runGit(t, "add", "internal/"+slug+"/foo.go")
	runGit(t, "commit", "-q", "-m", "c1")
	return root, slug
}

// seedBuildFresh builds the build-mode fresh fixture: a base commit on main carrying
// a bootstrap-conformant spec with no tickets/ directory at all, then a feature
// branch with one authorized change — the tree B1 answers not-applicable ticket rows
// over.
func seedBuildFresh(t *testing.T) (root, slug string) {
	t.Helper()
	slug = "example"
	root = initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	runGit(t, "add", "internal/"+slug+"/foo.go")
	runGit(t, "commit", "-q", "-m", "c1")
	return root, slug
}
