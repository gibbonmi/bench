package durablefiletest

import (
	"errors"
	"io/fs"
	"os"
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
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Fatalf("%s = %q, want %q", path, data, want)
	}
}
