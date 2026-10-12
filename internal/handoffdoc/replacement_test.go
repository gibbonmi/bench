package handoffdoc

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

// replacementSection is the section that each replacement test writes.
var replacementSection = Section{
	Key:    "9f2ab77",
	Fields: []Field{{LabelRequestToken, "`9f2ab77`"}},
	Next:   "`bench gate`",
	State:  "The phase is live.",
}

// D21: a section write returns the injected directory-sync failure of its replacement
// after the rename, and the document holds the rendered document at mode 0644.
func TestHandoffReplacement(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "session-handoff.md")
	cause := errors.New("injected handoff directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

	err := writeSection(rec.Replacer(), path, replacementSection)

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	durablefiletest.RequireFile(t, path, 0o644, "## main\n\nNext command:\n\n### State\n\n"+
		"## request 9f2ab77\n\nRequest token: `9f2ab77`\nNext command: `bench gate`\n\n### State\n\nThe phase is live.\n\n")
}

// D36: a rename failure of a section write returns an unpublished failure that wraps its
// cause and leaves the prior document bytes.
func TestHandoffPreRenameFailure(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "session-handoff.md")
	if err := EnsureMain(path); err != nil {
		t.Fatal(err)
	}
	prior, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	cause := errors.New("injected handoff rename failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

	err = writeSection(rec.Replacer(), path, replacementSection)

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	durablefiletest.RequireFile(t, path, 0o644, string(prior))
}
