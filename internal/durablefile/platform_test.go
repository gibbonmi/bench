package durablefile_test

import (
	"bytes"
	"os"
	"runtime"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

// TestNativeReplacement qualifies the operating-system operations on the filesystem of the
// test's temporary directory: every stage, including both synchronizations, succeeds. The
// run log names the host and the directory, and the evidence of the run names that
// directory's filesystem.
func TestNativeReplacement(t *testing.T) {
	path, payload := destination(t, []byte("old")), []byte("native payload\n")
	t.Logf("host %s/%s, destination %s", runtime.GOOS, runtime.GOARCH, path)
	rec := &durablefiletest.Recorder{}
	if err := rec.Replacer().Replace(path, payload, 0o644); err != nil {
		t.Fatalf("native replacement through each stage: %v", err)
	}
	if want := slices.Concat(durablefiletest.BeforeRename, durablefiletest.AfterRename); !slices.Equal(rec.Trace, want) {
		t.Fatalf("native trace = %q, want %q", rec.Trace, want)
	}
	if err := durablefile.Replace(path, payload[:6], 0o600); err != nil {
		t.Fatalf("native Replace: %v", err)
	}
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := readBack(t, path); !bytes.Equal(got, payload[:6]) || info.Mode().Perm() != 0o600 || !info.Mode().IsRegular() {
		t.Fatalf("destination = %q at %v, want %q in a regular file at 0600", got, info.Mode(), payload[:6])
	}
}
