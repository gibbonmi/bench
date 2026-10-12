package repairpilot

import (
	"errors"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

// D29: the pilot document store returns the injected directory-sync failure of its
// replacement after the rename, and the document holds the encoded replacement at the
// pilot mode. The replace-failure row of TestRepairPilotStorage holds the pre-rename refusal.
func TestPilotReplacement(t *testing.T) {
	options, replacement, _ := replacementFixture(t)
	cause := errors.New("injected pilot directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

	err := (Store{Path: documentPath(options), Files: FileOps{Replacer: rec.Replacer()}}).Replace(replacement)

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	durablefiletest.RequireFile(t, documentPath(options), 0o600, string(documentBytes(t, replacement)))
}
