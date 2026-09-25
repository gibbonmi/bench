package responsebound

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// The spill store sits under the Bench home at responses/<repo-key>/<scope>/. The repo
// key is the census key of the repository, and a process outside any repository uses
// noRepository. The scope is the assignment id of the worktree the process runs in, and
// primaryScope otherwise. A primary scope keeps the newest primaryRetained spills.
const (
	storeDirName    = "responses"
	noRepository    = "none"
	primaryScope    = "primary"
	primaryRetained = 64
	privateDir      = 0o700
	privateFile     = 0o600
	spillSuffix     = ".out"
	nameRandBytes   = 8
)

var (
	errUnsafePath = errors.New("spill path is not line-safe")
	errNotOwnDir  = errors.New("not a real directory")
)

// openSpill creates one new spill file and answers it with its absolute path. A path
// that is not line-safe, a store directory that is not a real directory, and a failed
// create all return an error, and the caller takes the create-failure route. A new
// spill in a primary scope then prunes that scope.
func openSpill(root func() string, retiring bool, create func(string) (io.WriteCloser, error)) (io.WriteCloser, string, error) {
	home, err := filepath.Abs(benchhome.Dir())
	if err != nil {
		return nil, "", err
	}
	key, scope := location(home, root(), retiring)
	path := filepath.Join(home, storeDirName, key, scope, spillName())
	if !sanitize.LineSafe(path) {
		return nil, "", errUnsafePath
	}
	if err := privateDirs(home, storeDirName, key, scope); err != nil {
		return nil, "", err
	}
	file, err := create(path)
	if err != nil {
		return nil, "", err
	}
	if scope == primaryScope {
		prunePrimary(filepath.Dir(path))
	}
	return file, path, nil
}

// location answers the repo key and the scope of a process whose repository root is
// root, or the empty string outside a repository. A retiring verb takes the primary
// scope wherever it runs.
func location(home, root string, retiring bool) (string, string) {
	if root == "" {
		return noRepository, primaryScope
	}
	key := poolkey.Key(root)
	if retiring {
		return key, primaryScope
	}
	segment, err := filepath.Rel(filepath.Join(poolkey.Pools(home), key), root)
	if err != nil {
		return key, primaryScope
	}
	if id, ok := poolkey.SplitAssignmentSegment(segment); ok {
		return key, id
	}
	return key, primaryScope
}

// privateDirs creates each store directory below home with mode 0700. It never follows a
// symlink: a component that already exists must be a real directory.
func privateDirs(home string, components ...string) error {
	if err := os.MkdirAll(home, privateDir); err != nil {
		return err
	}
	dir := home
	for _, component := range components {
		dir = filepath.Join(dir, component)
		if err := os.Mkdir(dir, privateDir); err != nil && !errors.Is(err, fs.ErrExist) {
			return err
		}
		if err := realDir(dir); err != nil {
			return err
		}
	}
	return nil
}

// realDir refuses a path that is not a real directory, so a symlink is never followed.
func realDir(dir string) error {
	info, err := os.Lstat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return &fs.PathError{Op: "mkdir", Path: dir, Err: errNotOwnDir}
	}
	return nil
}

// spillName generates a file name, so no operand or output byte forms the path. The time
// prefix orders the names by creation.
func spillName() string {
	random := make([]byte, nameRandBytes)
	_, _ = rand.Read(random)
	return fmt.Sprintf("%d-%x%s", time.Now().UnixNano(), random, spillSuffix)
}

// exclusiveCreate creates a new spill file at mode 0600. The exclusive create refuses an
// existing name, and a symlink at that name with it.
func exclusiveCreate(path string) (io.WriteCloser, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, privateFile)
}

// failureReason renders err for a spill-failed line. A path error names its operation
// and its cause but not its path, and the text is line-safe.
func failureReason(err error) string {
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return sanitize.Controls(pathErr.Op + ": " + pathErr.Err.Error())
	}
	return sanitize.Controls(err.Error())
}
