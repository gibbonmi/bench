package gocache

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// LockFile is the name of the cache lock inside the Bench build cache. A holder locks it
// shared for a run's span, and `bench cache clean` locks it exclusively for the removal.
// Go's own clean removes the two-hex subdirectories alone, so the file survives a clean.
const LockFile = "bench.lock"

// holdWait bounds how long a holder waits for a clean to finish. A clean is one
// `go clean -cache` call, so a run that starts beside one waits milliseconds. The wait is
// bounded rather than open, because a stuck clean must not hold a gate forever.
const holdWait = 5 * time.Second

// holdPoll is the interval between the retries inside holdWait.
const holdPoll = 20 * time.Millisecond

// RecordLock is the whole-file fcntl lock request of the given type. Every Bench record
// lock covers the whole file from byte zero, so the shape is stated once here rather than
// re-derived at each lock site.
func RecordLock(typ int16) syscall.Flock_t {
	return syscall.Flock_t{Type: typ, Whence: int16(io.SeekStart), Start: 0, Len: 0}
}

// heldLock is one open descriptor on a lock path and the count of this process's holders
// that share it. The count exists because closing any descriptor for a file drops every
// record lock this process holds on that file, so two in-process holders must share one
// descriptor and close it once.
type heldLock struct {
	file  *os.File
	count int
}

// held tracks this process's shared cache locks by lock path. POSIX record locks are owned
// per process: a second request on the same file replaces the first, and F_GETLK never
// reports the caller's own lock. So the in-process answer cannot come from the kernel, and
// this map is that answer — the same posture the gate execution lock keeps.
var held = struct {
	sync.Mutex
	paths map[string]*heldLock
}{paths: map[string]*heldLock{}}

// Holder is one acquired shared cache lock. Release returns the hold; the descriptor stays
// open while any holder in this process still needs it.
type Holder struct {
	path     string
	released bool
}

// Hold takes the shared cache lock for the Bench build cache that env derives, and keeps
// its descriptor open until Release. It creates the directory and the lock file when they
// are absent, because a first run must lock rather than refuse. A gate run, a lane run, and
// a `bench test` run each call it, so one rule covers every holder and two of them hold the
// lock at the same time without either waiting.
func Hold(env []string) (*Holder, error) {
	dir, err := Dir(env)
	if err != nil {
		return nil, err
	}
	return HoldDir(dir)
}

// HoldDir is Hold over an explicit directory, which is the seam the lock tests drive.
func HoldDir(dir string) (*Holder, error) {
	path := filepath.Join(dir, LockFile)
	held.Lock()
	defer held.Unlock()
	if entry := held.paths[path]; entry != nil {
		entry.count++
		return &Holder{path: path}, nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	if err := acquireShared(file); err != nil {
		file.Close()
		return nil, err
	}
	held.paths[path] = &heldLock{file: file, count: 1}
	return &Holder{path: path}, nil
}

// acquireShared requests the read lock without waiting, then retries until holdWait runs
// out. Only an exclusive holder — a running clean — refuses the request, so the retry
// answers the one contention that exists and reports the refusal rather than blocking.
func acquireShared(file *os.File) error {
	deadline := time.Now().Add(holdWait)
	for {
		lock := RecordLock(syscall.F_RDLCK)
		err := syscall.FcntlFlock(file.Fd(), syscall.F_SETLK, &lock)
		if err == nil {
			return nil
		}
		if !contended(err) || !time.Now().Before(deadline) {
			return err
		}
		time.Sleep(holdPoll)
	}
}

// Release returns this hold. The descriptor closes, and the record lock with it, when the
// last holder in this process releases. A second Release is a no-op, so a deferred Release
// beside an explicit one cannot drop another holder's lock.
func (h *Holder) Release() error {
	if h == nil || h.released {
		return nil
	}
	h.released = true
	held.Lock()
	defer held.Unlock()
	entry := held.paths[h.path]
	if entry == nil {
		return nil
	}
	entry.count--
	if entry.count > 0 {
		return nil
	}
	delete(held.paths, h.path)
	lock := RecordLock(syscall.F_UNLCK)
	err := syscall.FcntlFlock(entry.file.Fd(), syscall.F_SETLK, &lock)
	if closeErr := entry.file.Close(); err == nil {
		err = closeErr
	}
	return err
}

// localHolder reports whether this process already holds the shared lock on path. A clean
// asks first, because the kernel answers F_GETLK about other processes alone and would let
// a clean inside a holder's own process replace that holder's lock.
func localHolder(path string) bool {
	held.Lock()
	defer held.Unlock()
	return held.paths[path] != nil
}

// contended reports whether err is the kernel's "another process holds a conflicting lock"
// answer. Linux returns EAGAIN and POSIX allows EACCES, so both spellings count.
func contended(err error) bool {
	return errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EACCES)
}
