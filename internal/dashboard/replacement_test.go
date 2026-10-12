package dashboard

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
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
	target := filepath.Join(t.TempDir(), "bench-dashboard.html")
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
	target := filepath.Join(t.TempDir(), "bench-dashboard.html")
	const prior = "<!DOCTYPE html>prior page"
	if err := os.WriteFile(target, []byte(prior), 0o600); err != nil {
		t.Fatal(err)
	}
	cause := errors.New("injected dashboard rename failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

	err := atomicWrite(rec.Replacer(), target, replacementPage())

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != prior {
		t.Fatalf("page = %q, want the prior page %q", got, prior)
	}
}
