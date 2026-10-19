package ownership_test

import (
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// ticketsRepo seeds the conformant build repository whose main branch carries tickets as
// its whole tickets folder, and whose feature branch edits internal/example/foo.go. A nil
// map removes the folder, and an empty map leaves it present and empty. It answers the
// ordinary build arguments over that tip.
func ticketsRepo(t *testing.T, tickets map[string]string) []string {
	t.Helper()
	root, _ := preflighttest.SeedConformant(t)
	preflighttest.RunGit(t, "checkout", "-q", "main")
	dir := "specs/" + slug + "/tickets"
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	for name, body := range tickets {
		preflighttest.MustWriteFile(t, dir+"/"+name, body)
	}
	commitmenttest.SeedAdmission(t, root, specPath)
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", "tickets")
	preflighttest.RunGit(t, "checkout", "-q", "-B", "feature")
	preflighttest.MustWriteFile(t, "internal/example/foo.go", "package example\n\nconst c1 = 1\n")
	preflighttest.RunGit(t, "add", "internal/example/foo.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")
	if tickets != nil && len(tickets) == 0 {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return []string{"build", slug, "--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", preflighttest.RunGit(t, "rev-parse", "HEAD")}
}

// TestUncitedRows drives the ordinary build over each tickets folder state.
func TestUncitedRows(t *testing.T) {
	t.Run("OC61 OC62 uncited row", func(t *testing.T) {
		build := ticketsRepo(t, map[string]string{"one.md": preflighttest.TicketDoc("One", "PF1") + "\nPF2 waits for a later ticket.\n"})
		verdict, code := preflight.Command(build)
		if code != 1 || !strings.Contains(verdict, "rows-owned,red") || !strings.Contains(verdict, "cited by no ticket file: PF2") {
			t.Fatalf("ordinary build = (%d), want the rows-owned red for PF2:\n%s", code, verdict)
		}
	})
	t.Run("OC64 present empty directory", func(t *testing.T) {
		if verdict, _ := preflight.Command(ticketsRepo(t, map[string]string{})); !strings.Contains(verdict, "rows-owned,red") {
			t.Fatalf("ordinary build over an empty directory lacks the rows-owned red:\n%s", verdict)
		}
	})
	t.Run("OC65 absent directory", func(t *testing.T) {
		if verdict, _ := preflight.Command(ticketsRepo(t, nil)); strings.Contains(verdict, "rows-owned") {
			t.Fatalf("ordinary build over an absent directory grades rows-owned:\n%s", verdict)
		}
	})
}

// TestExplicitAuthority grades that the ticket Writes: line is the ownership fence: a path
// that a leftover spec fence section omits still prepares the charge.
func TestExplicitAuthority(t *testing.T) {
	t.Run("OC86 ticket path outside the spec fence", func(t *testing.T) {
		tr := coveredTree()
		tr.Fence = []string{}
		for _, write := range tr.Writes {
			tr.Fence = append(tr.Fence, "- `"+write+"`")
		}
		tr.Writes = append(tr.Writes, "internal/other/extra.go")
		tr.Files["internal/other/extra.go"] = "package other\n"
		root := tr.Seed(t)
		if out, code := preflight.Command(preflighttest.ChargeArgs(t, root, slug)); code != 0 {
			t.Fatalf("charge for a ticket path outside the spec fence = (%d), want 0:\n%s", code, out)
		}
	})
}
