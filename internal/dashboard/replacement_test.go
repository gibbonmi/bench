package dashboard

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// replacementPage is a real rendered page with board content.
func replacementPage() string {
	s := baseSnapshot()
	s.Ideas = []string{"Render the dashboard through the shared replacement."}
	return Render(s)
}

// D26: a directory-sync failure after the rename fails the page write, and the published
// page holds the complete rendered HTML at mode 0600.
func TestDashboardReplacement(t *testing.T) {
	t.Parallel()
	target := filepath.Join(t.TempDir(), pageName)
	page := replacementPage()
	cause := errors.New("injected dashboard directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

	err := atomicWrite(rec.Replacer(), target, page)

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	durablefiletest.RequireFile(t, target, 0o600, page)
}

// D36: a rename failure refuses the page write and leaves the prior page.
func TestDashboardPreRenameFailure(t *testing.T) {
	t.Parallel()
	target := filepath.Join(t.TempDir(), pageName)
	const prior = "<!DOCTYPE html>prior page"
	if err := os.WriteFile(target, []byte(prior), 0o600); err != nil {
		t.Fatal(err)
	}
	cause := errors.New("injected dashboard rename failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

	err := atomicWrite(rec.Replacer(), target, replacementPage())

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	durablefiletest.RequireBytes(t, target, prior)
}

// A directory at the page path makes the native replacement fail, and Command reports
// that failure as its write refusal.
func TestDashboardCommandReportsReplacementFailure(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	gitDir, err := git.AdminDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(gitDir, pageName), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	r, c := Command(nil)

	if c != 1 || !strings.Contains(r, "cannot write bench-dashboard.html") {
		t.Fatalf("Command(nil) = %q, %d, want exit 1 and %q", r, c, "cannot write bench-dashboard.html")
	}
}
