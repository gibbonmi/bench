package publication

import (
	"errors"
	"os"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

// replacementRecord returns one publication record with result.
func replacementRecord(result string) Record {
	return Record{
		ReleaseIndexSHA256: "abc",
		Path:               "first",
		Profile:            "public",
		Transitions:        []Transition{},
		Provenance:         []Provenance{{Package: "@bench/cli", SHA256: "def", Kind: "wrapper"}},
		Result:             result,
	}
}

// D24: a record save creates its missing record directory and returns the injected
// directory-sync failure of its replacement after the rename, and the record holds the
// encoded record at mode 0644.
func TestPublicationReplacement(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	cause := errors.New("injected publication directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

	err := saveRecordWith(rec.Replacer(), root, replacementRecord("staged"))

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	durablefiletest.RequireFile(t, RecordPath(root), 0o644, `{
  "schema_version": 1,
  "release_index_sha256": "abc",
  "path": "first",
  "profile": "public",
  "transitions": [],
  "provenance": [
    {
      "package": "@bench/cli",
      "sha256": "def",
      "kind": "wrapper"
    }
  ],
  "result": "staged"
}
`)
}

// D36: a rename failure of a record save returns an unpublished failure that wraps its
// cause and leaves the former record.
func TestPublicationPreRenameFailure(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := SaveRecord(root, replacementRecord("staged")); err != nil {
		t.Fatal(err)
	}
	former, err := os.ReadFile(RecordPath(root))
	if err != nil {
		t.Fatal(err)
	}
	cause := errors.New("injected publication rename failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

	err = saveRecordWith(rec.Replacer(), root, replacementRecord("published"))

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	durablefiletest.RequireFile(t, RecordPath(root), 0o644, string(former))
}
