package brokermanifest

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

// replacementBroker writes a broker executable in a fresh directory and returns that
// directory and the symlink-resolved broker path the manifest binds.
func replacementBroker(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	broker := filepath.Join(dir, "bench")
	if err := os.WriteFile(broker, []byte("broker executable\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	resolved, err := filepath.EvalSymlinks(broker)
	if err != nil {
		t.Fatal(err)
	}
	return dir, resolved
}

// D25: a directory-sync failure after the rename fails the manifest write, because a
// manifest the land route trusts must be durable before Write succeeds.
func TestBrokerReplacement(t *testing.T) {
	t.Parallel()
	dir, broker := replacementBroker(t)
	cause := errors.New("injected broker directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

	_, _, err := write(rec.Replacer(), dir, broker, "1.2.3")

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	durablefiletest.RequireFile(t, filepath.Join(dir, Name), 0o644, "path\t"+broker+"\nversion\t1.2.3\nplatform\t"+Platform()+
		"\nsha256\tc153e59aa5f0b544d81e9534460e0199291ae0258c26d16040eedaf5d0f92330\n")
}

// D36: a rename failure refuses the manifest write and leaves no manifest.
func TestBrokerPreRenameFailure(t *testing.T) {
	t.Parallel()
	dir, broker := replacementBroker(t)
	cause := errors.New("injected broker rename failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

	_, _, err := write(rec.Replacer(), dir, broker, "1.2.3")

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	requireNoManifest(t, dir)
}

// D30: a broker that is not a regular file refuses the manifest with the digest error
// before the replacement starts.
func TestBrokerRefusalSkipsReplacement(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	rec := &durablefiletest.Recorder{}

	_, _, err := write(rec.Replacer(), dir, dir, "1.2.3")

	var replaceErr *durablefile.Error
	if err == nil || errors.As(err, &replaceErr) {
		t.Fatalf("err = %v, want the caller's own refusal, not a replacement failure", err)
	}
	if len(rec.Trace) != 0 {
		t.Fatalf("trace = %v, want no replacement stage", rec.Trace)
	}
	requireNoManifest(t, dir)
}

// requireNoManifest fails t unless dir holds no manifest.
func requireNoManifest(t *testing.T, dir string) {
	t.Helper()
	if _, err := os.Lstat(filepath.Join(dir, Name)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("manifest stat err = %v, want it absent", err)
	}
}
