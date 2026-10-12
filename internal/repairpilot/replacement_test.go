package repairpilot

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

// D29: the pilot document store returns the injected directory-sync failure of its
// replacement after the rename, and the document holds the encoded replacement. The custody
// and replace-failure rows of TestRepairPilotStorage hold the mode and the pre-rename refusal.
func TestPilotReplacement(t *testing.T) {
	options, replacement, _ := replacementFixture(t)
	cause := errors.New("injected pilot directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

	err := (Store{Path: documentPath(options), Files: FileOps{Replacer: rec.Replacer()}}).Replace(replacement)

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	durablefiletest.RequireBytes(t, documentPath(options), string(documentBytes(t, replacement)))
}

// D30: activation refuses a pilot directory that it cannot create before any replacement
// starts, and no pilot document exists.
func TestPilotDirectoryRefusal(t *testing.T) {
	options := testOptions(t.TempDir())
	cause := errors.New("injected pilot directory refusal")
	rec := &durablefiletest.Recorder{}
	options.Files = FileOps{
		MkdirAll: func(string, os.FileMode) error { return cause },
		Replacer: rec.Replacer(),
	}

	out, code := Command(options, []string{"activate"})

	if code != 1 || !strings.Contains(out, cause.Error()) {
		t.Fatalf("activate = output %q, exit %d; want the injected refusal at exit 1", out, code)
	}
	if len(rec.Trace) != 0 {
		t.Fatalf("replacement trace = %v, want no replacement stage", rec.Trace)
	}
	durablefiletest.RequireAbsent(t, documentPath(options))
}
