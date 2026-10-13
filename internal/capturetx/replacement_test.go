package capturetx

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

func captureManifest(state string) manifest {
	return manifest{Schema: 1, ID: "d-0123456789ab", State: state, Sources: []manifestSource{{
		Name: "ideas", Path: "/repo/capture/IDEAS.md", Blob: "source-0.bin", Bytes: 6, Digest: "abc",
	}}}
}

const sealedManifest = `{
  "schema": 1,
  "id": "d-0123456789ab",
  "state": "sealed",
  "sources": [
    {
      "name": "ideas",
      "path": "/repo/capture/IDEAS.md",
      "blob": "source-0.bin",
      "bytes": 6,
      "digest": "abc"
    }
  ]
}
`

// D20: the capture manifest writer creates its missing transaction directory and returns
// the injected directory-sync failure of its replacement after the rename, and the
// manifest holds the encoded manifest at mode 0644.
func TestCaptureReplacement(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), directoryName)
	cause := errors.New("injected capture directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

	err := writeManifestWith(rec.Replacer(), dir, captureManifest("sealed"))

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	durablefiletest.RequireFile(t, filepath.Join(dir, manifestName), 0o644, sealedManifest)
}

// D36: a rename failure of the capture manifest replacement returns an unpublished
// failure that wraps its cause and leaves the former manifest.
func TestCapturePreRenameFailure(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := writeManifestDir(dir, captureManifest("preparing")); err != nil {
		t.Fatal(err)
	}
	former, err := os.ReadFile(filepath.Join(dir, manifestName))
	if err != nil {
		t.Fatal(err)
	}
	cause := errors.New("injected capture rename failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

	err = writeManifestWith(rec.Replacer(), dir, captureManifest("sealed"))

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	durablefiletest.RequireBytes(t, filepath.Join(dir, manifestName), string(former))
}

// D30: a manifest write whose transaction directory cannot be created refuses before the
// replacement starts, and the blocking file keeps its bytes.
func TestCaptureDirectoryRefusal(t *testing.T) {
	t.Parallel()
	const blocked = "blocking file"
	parent := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(parent, []byte(blocked), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(parent, directoryName)
	rec := &durablefiletest.Recorder{}

	err := writeManifestWith(rec.Replacer(), dir, captureManifest("sealed"))

	durablefiletest.RequireRefusal(t, err, rec.Trace)
	durablefiletest.RequireBytes(t, parent, blocked)
}
