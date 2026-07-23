package gate

import (
	"io"
	"os"
	"sync"
	"syscall"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

type gateEngine interface {
	Now() time.Time
	BuildSubject(string) (subject, error)
	PostRunSubject(string) (subject, error)
	GitDir(string) (string, error)
	OpenLock(string) (gateFile, error)
	Acquire(gateFile) error
	Unlock(gateFile) error
	CreateTemp(string, string) (gateFile, error)
	Rename(string, string) error
	OpenDir(string) (gateFile, error)
	WriteFile(string, []byte, os.FileMode) error
	Remove(string) error
}

type productionGateEngine struct{}

var executionLockOwners = struct {
	sync.Mutex
	paths map[string]bool
}{paths: map[string]bool{}}

func recordLock(typ int16) syscall.Flock_t {
	return syscall.Flock_t{Type: typ, Whence: int16(io.SeekStart), Start: 0, Len: 0}
}

func (productionGateEngine) Now() time.Time                              { return time.Now().UTC() }
func (productionGateEngine) BuildSubject(root string) (subject, error)   { return buildSubject(root) }
func (productionGateEngine) PostRunSubject(root string) (subject, error) { return buildSubject(root) }
func (productionGateEngine) GitDir(root string) (string, error) {
	return benchgit.Output("-C", root, "rev-parse", "--absolute-git-dir")
}
func (productionGateEngine) OpenLock(path string) (gateFile, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
}
func (productionGateEngine) Acquire(f gateFile) error {
	executionLockOwners.Lock()
	defer executionLockOwners.Unlock()
	if executionLockOwners.paths[f.Name()] {
		return syscall.EAGAIN
	}
	lock := recordLock(syscall.F_WRLCK)
	if err := syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, &lock); err != nil {
		return err
	}
	executionLockOwners.paths[f.Name()] = true
	return nil
}
func (productionGateEngine) Unlock(f gateFile) error {
	executionLockOwners.Lock()
	defer executionLockOwners.Unlock()
	lock := recordLock(syscall.F_UNLCK)
	err := syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, &lock)
	delete(executionLockOwners.paths, f.Name())
	return err
}
func (productionGateEngine) CreateTemp(dir, pattern string) (gateFile, error) {
	return os.CreateTemp(dir, pattern)
}
func (productionGateEngine) Rename(oldpath, newpath string) error  { return os.Rename(oldpath, newpath) }
func (productionGateEngine) OpenDir(path string) (gateFile, error) { return os.Open(path) }
func (productionGateEngine) WriteFile(path string, data []byte, mode os.FileMode) error {
	return os.WriteFile(path, data, mode)
}
func (productionGateEngine) Remove(path string) error { return os.Remove(path) }
