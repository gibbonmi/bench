package probe

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/toon"
)

// startReadHook runs between the start read and the subject lock, and lockOwnerHook runs
// between the lock file's open and its identity check. Both are no-ops outside the tests,
// which change the subject or the lock file inside those windows.
var (
	startReadHook = func() {}
	lockOwnerHook = func(path string) {}
)

// lockAttempts bounds the acquire, so a lock file that keeps changing refuses rather than
// holding the probe in a retry loop.
const lockAttempts = 3

// errLockReplaced is the reason of a lock file that another file replaced at each attempt.
var errLockReplaced = errors.New("the lock file was replaced while the lock was taken")

// subjectLock is the held lock of one subject: the open lock file and the path it was
// opened at.
type subjectLock struct {
	file *os.File
	path string
}

// subjectLockPath answers the lock file of the subject at the absolute path. The name is the
// path's SHA-256, so no operand byte reaches it. The directory is a sibling of the preserve
// directory, so a lock release never removes a kept copy.
func subjectLockPath(path string) string {
	sum := sha256.Sum256([]byte(path))
	return filepath.Join(benchhome.Dir(), "probe-locks", hex.EncodeToString(sum[:]))
}

// lockSubject takes the subject's lock without waiting, because a waiting probe would hide a
// concurrent author. The lock is a record lock, so it dies with its process, and a lock file
// that a killed probe left grants the lock. An attempt whose open answers ENOENT, or whose
// held descriptor no longer names the file at the path, starts again: a holder's release
// removes the file that a waiter may already have opened.
func lockSubject(subject subject) (subjectLock, string) {
	startReadHook()
	path := subjectLockPath(subject.path)
	var last error
	for range lockAttempts {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return subjectLock{}, lockUnavailable(err)
		}
		// os.OpenFile opens every descriptor close-on-exec. A record lock needs write access.
		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW, 0o600)
		if errors.Is(err, fs.ErrNotExist) {
			last = err
			continue
		}
		if err != nil {
			return subjectLock{}, lockUnavailable(err)
		}
		request := gocache.RecordLock(syscall.F_WRLCK)
		if err := syscall.FcntlFlock(file.Fd(), syscall.F_SETLK, &request); err != nil {
			_ = file.Close()
			return subjectLock{}, toon.Errorf("probe subject in use", "another probe holds "+subject.display) + "\n"
		}
		lockOwnerHook(path)
		if namesFile(path, file) {
			return subjectLock{file: file, path: path}, ""
		}
		_ = file.Close()
		last = errLockReplaced
	}
	return subjectLock{}, lockUnavailable(last)
}

// namesFile reports whether path still names the file that the descriptor holds. It reads
// the held identity through the descriptor, so it opens no second descriptor.
func namesFile(path string, file *os.File) bool {
	held, err := file.Stat()
	if err != nil {
		return false
	}
	named, err := os.Lstat(path)
	return err == nil && os.SameFile(held, named)
}

// release removes the lock file while the lock is still held, so a waiter that opened that
// file fails its identity check. It then closes the descriptor, which drops the lock, and
// removes the lock directory. That removal fails while another probe's lock file is in it,
// and the failure is the intended answer, so it is ignored.
func (l subjectLock) release() {
	_ = os.Remove(l.path)
	_ = l.file.Close()
	_ = os.Remove(filepath.Dir(l.path))
}

func lockUnavailable(err error) string {
	return toon.Errorf("probe subject lock unavailable", err.Error()) + "\n"
}

// rereadSubject reads the subject again under the lock and compares it with the start read.
// Another probe can restore its mutant between the start read and the lock, so a probe that
// skipped this read would preserve that mutant as its own start bytes.
func rereadSubject(root, operand string, start subject) string {
	current, line := resolveSubject(root, operand)
	if line != "" {
		return line
	}
	if !bytes.Equal(current.start, start.start) || current.mode != start.mode {
		return toon.Errorf("probe subject changed", start.display+" changed before the subject lock; run the probe again") + "\n"
	}
	return ""
}
