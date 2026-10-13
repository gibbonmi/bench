package durablefiletest

import (
	"errors"
	"io/fs"
	"os"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// PostRenameStages names each stage that fails after the rename, so a caller test can
// fail every one of them and expect a published failure.
var PostRenameStages = []string{"directory-open", "directory-sync", "directory-close"}

// RequireFailure fails t unless err is a *durablefile.Error that failed at stage, whose
// Published field equals published, and whose chain holds cause.
func RequireFailure(t testing.TB, err error, stage string, published bool, cause error) {
	t.Helper()
	var replaceErr *durablefile.Error
	if !errors.As(err, &replaceErr) || replaceErr.Stage != stage || replaceErr.Published != published || !errors.Is(err, cause) {
		t.Fatalf("err = %v, want a *durablefile.Error at stage %s with Published=%t that wraps %v", err, stage, published, cause)
	}
}

// RequireRefusal fails t unless err is the caller's own refusal, not a *durablefile.Error,
// and trace holds no replacement step.
func RequireRefusal[T any](t testing.TB, err error, trace []T) {
	t.Helper()
	var replaceErr *durablefile.Error
	if err == nil || errors.As(err, &replaceErr) {
		t.Fatalf("err = %v, want the caller's own refusal, not a replacement failure", err)
	}
	RequireNoReplacement(t, trace)
}

// RequireNoReplacement fails t unless trace holds no replacement step. A caller whose
// entry reports its refusal as output, not as an error, checks only this half.
func RequireNoReplacement[T any](t testing.TB, trace []T) {
	t.Helper()
	if len(trace) != 0 {
		t.Fatalf("replacement trace = %v, want no replacement step", trace)
	}
}

// RequireFile fails t unless path names a regular file, not a symbolic link, whose
// permission bits equal mode and whose bytes equal want.
func RequireFile(t testing.TB, path string, mode fs.FileMode, want string) {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != mode {
		t.Fatalf("%s mode = %v, want a regular file at %v", path, info.Mode(), mode)
	}
	RequireBytes(t, path, want)
}

// RequireAbsent fails t unless nothing exists at path, not even a symbolic link. A path
// below a regular file is absent too, because a refusal test can block a parent that way.
func RequireAbsent(t testing.TB, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) && !errors.Is(err, syscall.ENOTDIR) {
		t.Fatalf("%s stat err = %v, want it absent", path, err)
	}
}

// RequireBytes fails t unless the bytes of the file at path equal want.
func RequireBytes(t testing.TB, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Fatalf("%s = %q, want %q", path, data, want)
	}
}
