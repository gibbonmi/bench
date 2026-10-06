package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

type image struct {
	Identity *fileIdentity `json:"identity,omitempty"`
	Kind     string        `json:"kind"`
	Mode     os.FileMode   `json:"mode,omitempty"`
	Link     string        `json:"link,omitempty"`
	Digest   string        `json:"digest,omitempty"`
	Data     []byte        `json:"-"`
}

func inspect(path string) (image, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return image{Kind: "absent"}, nil
	}
	if err != nil {
		return image{}, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		return image{Kind: "symlink", Link: target, Identity: identityOf(info)}, err
	}
	if !info.Mode().IsRegular() {
		return image{}, fmt.Errorf("not a regular managed asset: %s", path)
	}
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return image{}, err
	}
	current, statErr := f.Stat()
	if statErr != nil || !current.Mode().IsRegular() || !os.SameFile(info, current) {
		_ = f.Close()
		return image{}, fmt.Errorf("managed asset changed during read: %s", path)
	}
	data, readErr := io.ReadAll(f)
	closeErr := f.Close()
	if err := errors.Join(readErr, closeErr); err != nil {
		return image{}, err
	}
	return image{Kind: "file", Mode: current.Mode(), Digest: digestBytes(data), Data: data, Identity: identityOf(current)}, nil
}

func same(a, b image) bool {
	return a.Kind == b.Kind && a.Mode == b.Mode && a.Link == b.Link && a.Digest == b.Digest
}

type fileIdentity struct {
	Device uint64 `json:"device"`
	Inode  uint64 `json:"inode"`
}

func identityOf(info os.FileInfo) *fileIdentity {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil
	}
	return &fileIdentity{Device: uint64(stat.Dev), Inode: uint64(stat.Ino)}
}

func sameObserved(a, b image) bool {
	if !same(a, b) {
		return false
	}
	if a.Kind == "absent" {
		return true
	}
	return a.Identity != nil && b.Identity != nil && *a.Identity == *b.Identity
}

// SyncDirectory reports both persistence and close failures for a directory.
func SyncDirectory(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return errors.Join(f.Sync(), f.Close())
}

func digestBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func ensureDirectory(path string, mode os.FileMode, persist func(string) error) error {
	info, err := os.Lstat(path)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("not a directory: %s", path)
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	parent := filepath.Dir(path)
	if parent == path {
		return err
	}
	if err := ensureDirectory(parent, mode, persist); err != nil {
		return err
	}
	if err := os.Mkdir(path, mode); err != nil {
		if !os.IsExist(err) {
			return err
		}
		return ensureDirectory(path, mode, persist)
	}
	return persist(parent)
}
