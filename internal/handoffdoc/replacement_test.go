package handoffdoc

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

var replacementSection = Section{
	Key:    "9f2ab77",
	Fields: []Field{{LabelRequestToken, "`9f2ab77`"}},
	Next:   "`bench gate`",
	State:  "The phase is live.",
}

// D21: a directory failure after the rename fails the section write.
func TestHandoffReplacement(t *testing.T) {
	t.Parallel()
	for _, stage := range durablefiletest.PostRenameStages {
		t.Run(stage, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "session-handoff.md")
			cause := errors.New("injected handoff " + stage + " failure")
			rec := &durablefiletest.Recorder{Fail: map[string]error{stage: cause}}

			err := writeSection(rec.Replacer(), path, replacementSection)

			durablefiletest.RequireFailure(t, err, stage, true, cause)
			durablefiletest.RequireFile(t, path, 0o644, "## main\n\nNext command:\n\n### State\n\n"+
				"## request 9f2ab77\n\nRequest token: `9f2ab77`\nNext command: `bench gate`\n\n### State\n\nThe phase is live.\n\n")
		})
	}
}

// D36: a failure before the rename refuses the section write and leaves the prior
// document.
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
	durablefiletest.RequireBytes(t, path, string(prior))
}

// D30: a refused document change returns its own error before the replacement starts,
// and no document exists.
func TestHandoffRefusalSkipsReplacement(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "session-handoff.md")
	rec := &durablefiletest.Recorder{}

	err := update(rec.Replacer(), path, func(*Document) error { return errors.New("injected handoff change refusal") })

	durablefiletest.RequireRefusal(t, err, rec.Trace)
	durablefiletest.RequireAbsent(t, path)
}
