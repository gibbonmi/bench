package handoffdoc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// lockDeadline bounds the wait for the document lock, and lockPoll is the retry
// interval. Both mirror the intent ledger's acquire loop, which is the repository's
// existing answer for how long a caller waits on a small local file. The ledger's
// own values are unexported literals inside its acquire function, so this is a
// mirror rather than a shared constant.
const (
	lockDeadline = 2 * time.Second
	lockPoll     = 10 * time.Millisecond
)

// DocumentPath is the session handoff document's path, relative to the checkout that
// owns it. The leaf owns the grammar, so it owns the name too: a consumer that spelled
// the path itself could watch a file the writer never touches.
const DocumentPath = "capture/session-handoff.md"

// LockPath is the lock file that guards one document. It sits beside the document
// rather than inside it, so a reader that opens the document never contends.
func LockPath(path string) string { return path + ".lock" }

// Read parses the document at path without taking the lock. An absent or empty
// file reads as a fresh document, so a first run in a repo needs no scaffold step.
// Callers that intend to write go through Update instead.
func Read(path string) (*Document, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return New(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read handoff document: %w", err)
	}
	if len(content) == 0 {
		return New(), nil
	}
	return Parse(path, content)
}

// Update runs one read-modify-write under an exclusive lock, then replaces the
// file by rename. The lock is what makes two live phases safe: temp-and-rename
// alone makes each write atomic, but two writers that both read the old document
// still land one section apiece and the later rename drops the other's.
func Update(path string, mutate func(*Document) error) error {
	release, err := acquire(path)
	if err != nil {
		return err
	}
	defer release()

	doc, err := Read(path)
	if err != nil {
		return err
	}
	if err := mutate(doc); err != nil {
		return err
	}
	return replace(path, doc.Render())
}

// WriteSection rewrites one section and leaves every other section's bytes as the
// file had them.
func WriteSection(path string, section Section) error {
	return Update(path, func(doc *Document) error {
		doc.Put(section)
		doc.EnsureMain()
		return nil
	})
}

// RemoveSection drops the section under key and leaves main behind. It is the
// retirement path's one call, so a landing, a release, and a clean all remove a
// section the same way.
func RemoveSection(path, key string) error {
	return Update(path, func(doc *Document) error {
		doc.Remove(key)
		doc.EnsureMain()
		return nil
	})
}

// EnsureMain writes main into the document when it holds none.
func EnsureMain(path string) error {
	return Update(path, func(doc *Document) error {
		doc.EnsureMain()
		return nil
	})
}

// acquire takes the exclusive flock on the lock file and returns its release. The
// lock is polled rather than blocking, so a holder that never releases refuses the
// caller at the deadline instead of hanging a phase close forever.
//
// The release unlinks the lock file, so the document is left alone on disk and a
// landing sees no undeclared residue. Unlinking makes the flock reclaimable rather
// than authoritative: a waiter can be holding the removed inode by the time it wins
// the lock, and a third writer that arrives after the unlink locks a fresh file at
// the same path. So the winner compares the open file against the path and reopens
// when they differ. A writer killed mid-hold leaves the file behind, and the next
// writer takes it as it stands.
func acquire(path string) (func(), error) {
	lock := LockPath(path)
	if err := os.MkdirAll(filepath.Dir(lock), 0o755); err != nil {
		return nil, fmt.Errorf("open handoff lock %s: %w", lock, err)
	}
	deadline := time.Now().Add(lockDeadline)
	for {
		file, err := os.OpenFile(lock, os.O_CREATE|os.O_RDWR, 0o600)
		if err != nil {
			return nil, fmt.Errorf("open handoff lock %s: %w", lock, err)
		}
		held, err := flockUntil(file, deadline)
		if err != nil {
			_ = file.Close()
			return nil, fmt.Errorf("lock handoff document at %s: %w", lock, err)
		}
		if !held {
			_ = file.Close()
			return nil, stillHeld(lock)
		}
		if !isFileAt(file, lock) {
			_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
			_ = file.Close()
			if time.Now().After(deadline) {
				return nil, stillHeld(lock)
			}
			continue
		}
		return func() {
			_ = os.Remove(lock)
			_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
			_ = file.Close()
		}, nil
	}
}

// stillHeld is the deadline refusal. It names the lock path, because deleting that
// file is the operator's remedy when no writer is left to release it.
func stillHeld(lock string) error {
	return fmt.Errorf("handoff lock %s is still held after %s; wait for the other phase close, or delete the lock file if no writer holds it", lock, lockDeadline)
}

// flockUntil polls for the exclusive flock until it wins or the deadline passes. A
// false with no error is the deadline; an error is the lock call itself failing.
func flockUntil(file *os.File, deadline time.Time) (bool, error) {
	for {
		err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return true, nil
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) {
			return false, err
		}
		if time.Now().After(deadline) {
			return false, nil
		}
		time.Sleep(lockPoll)
	}
}

// isFileAt reports whether the open file is still the file the path names. An absent
// path and a replaced one both answer false, and both mean the same thing: the lock
// this caller holds no longer guards the document.
func isFileAt(file *os.File, path string) bool {
	opened, err := file.Stat()
	if err != nil {
		return false
	}
	current, err := os.Stat(path)
	if err != nil {
		return false
	}
	return os.SameFile(opened, current)
}

// replace writes the rendered document to a sibling temp file, syncs it, and
// renames it over the document. A writer killed between the write and the rename
// leaves the prior document intact.
func replace(path, content string) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("write handoff document: %w", err)
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.WriteString(content); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write handoff document: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync handoff document: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close handoff document: %w", err)
	}
	if err := os.Chmod(name, 0o644); err != nil {
		return fmt.Errorf("write handoff document: %w", err)
	}
	if err := os.Rename(name, path); err != nil {
		return fmt.Errorf("replace handoff document: %w", err)
	}
	if handle, err := os.Open(dir); err == nil {
		_ = handle.Sync()
		_ = handle.Close()
	}
	return nil
}
