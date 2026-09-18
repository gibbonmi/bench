package chargeevidence

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"syscall"
)

// lock is one held whole-file advisory lock.
type lock struct{ file *os.File }

func (l *lock) release() {
	if l != nil && l.file != nil {
		_ = syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		_ = l.file.Close()
		l.file = nil
	}
}

// acquire takes one store lock. A writer creates an absent lock file. Only preparation
// creates a lock file, so a reader of a store without one refuses the absent artifact and
// changes nothing. A caller that passes LOCK_NB refuses an active holder instead of waiting.
func acquire(dir *os.Root, name string, how int, create bool) (*lock, error) {
	flags := os.O_RDWR | syscall.O_NOFOLLOW | syscall.O_NONBLOCK
	if create {
		flags |= os.O_CREATE
	}
	info, err := dir.Lstat(name)
	if err == nil && !info.Mode().IsRegular() {
		return nil, refuse(RefuseUnsafe, "store lock %s is %s", name, kindOf(info))
	}
	if errors.Is(err, fs.ErrNotExist) && !create {
		return nil, refuse(RefuseAbsent, "the evidence store holds no artifacts")
	}
	file, err := dir.OpenFile(name, flags, 0o600)
	if err != nil {
		return nil, refuse(RefuseUnsafe, "store lock %s is not openable: %v", name, err)
	}
	if info, err := file.Stat(); err != nil || !info.Mode().IsRegular() {
		file.Close()
		return nil, refuse(RefuseUnsafe, "store lock %s is not a regular file", name)
	}
	if err := syscall.Flock(int(file.Fd()), how); err != nil {
		file.Close()
		// A waiting caller that cannot take its lock met a storage failure. A caller that
		// asked not to wait met an active holder, which is the refusal cleanup reports.
		if how&syscall.LOCK_NB != 0 {
			return nil, refuse(RefuseBusy, "store lock %s is held by an active reader or writer", name)
		}
		return nil, refuse(RefuseStorage, "store lock %s is not acquirable: %v", name, err)
	}
	return &lock{file: file}, nil
}

// readRegular reads one store object after proving that the opened file is the regular file
// the directory names. It opens without blocking and without following a link.
func readRegular(dir *os.Root, name string) ([]byte, error) {
	file, info, err := openRegular(dir, name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, refuse(RefuseStorage, "store object %s is unreadable: %v", name, err)
	}
	if err := unchanged(dir, name, file, info); err != nil {
		return nil, err
	}
	return data, nil
}

func openRegular(dir *os.Root, name string) (*os.File, os.FileInfo, error) {
	before, err := dir.Lstat(name)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil, refuse(RefuseAbsent, "the evidence store holds no such artifact")
	}
	if err != nil {
		return nil, nil, refuse(RefuseUnsafe, "store object is not inspectable: %v", err)
	}
	if !before.Mode().IsRegular() {
		return nil, nil, refuse(RefuseUnsafe, "store object is %s, not a regular file", kindOf(before))
	}
	file, err := dir.OpenFile(name, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, nil, refuse(RefuseUnsafe, "store object is not openable: %v", err)
	}
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(before, opened) {
		file.Close()
		return nil, nil, refuse(RefuseReplaced, "store object changed while it was opened")
	}
	return file, opened, nil
}

// unchanged proves that the directory still names the file that was read, with the same
// length and modification time.
func unchanged(dir *os.Root, name string, file *os.File, opened os.FileInfo) error {
	after, err := file.Stat()
	current, lerr := dir.Lstat(name)
	if err != nil || lerr != nil || !os.SameFile(opened, current) || after.Size() != opened.Size() || !after.ModTime().Equal(opened.ModTime()) {
		return refuse(RefuseReplaced, "store object changed while it was read")
	}
	return nil
}

func kindOf(info os.FileInfo) string {
	switch mode := info.Mode(); {
	case mode&fs.ModeSymlink != 0:
		return "a symlink"
	case mode&fs.ModeNamedPipe != 0:
		return "a FIFO"
	case mode&fs.ModeSocket != 0:
		return "a socket"
	case mode&fs.ModeDevice != 0:
		return "a device"
	case mode.IsDir():
		return "a directory"
	}
	return "an irregular file"
}
