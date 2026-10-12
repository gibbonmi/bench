package brokermanifest

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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
	// The staging name ends in a random suffix with no dot, so everything up to its last dot
	// must equal the prefix exactly; a shorter prefix would still match as a prefix.
	if len(rec.Temps) != 1 {
		t.Fatalf("temps = %v, want one staging file", rec.Temps)
	}
	if temp := filepath.Base(rec.Temps[0]); temp[:strings.LastIndex(temp, ".")+1] != TemporaryPrefix {
		t.Fatalf("staging file %q, want it named %q plus a random suffix", temp, TemporaryPrefix)
	}
}

// D36: a rename failure refuses the manifest write and leaves no manifest.
func TestBrokerPreRenameFailure(t *testing.T) {
	t.Parallel()
	dir, broker := replacementBroker(t)
	cause := errors.New("injected broker rename failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

	_, _, err := write(rec.Replacer(), dir, broker, "1.2.3")

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	durablefiletest.RequireAbsent(t, filepath.Join(dir, Name))
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
	durablefiletest.RequireAbsent(t, filepath.Join(dir, Name))
}

// The public Write returns the refusal of its native replacement path.
func TestBrokerWriteReturnsRefusal(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	if _, _, err := Write(dir, dir, "1.2.3"); err == nil {
		t.Fatal("Write err = nil, want the refusal of a broker that is not a regular file")
	}
	durablefiletest.RequireAbsent(t, filepath.Join(dir, Name))
}
