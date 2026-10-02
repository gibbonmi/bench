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
	identity os.FileInfo
	Kind     string      `json:"kind"`
	Mode     os.FileMode `json:"mode,omitempty"`
	Link     string      `json:"link,omitempty"`
	Digest   string      `json:"digest,omitempty"`
	Data     []byte      `json:"-"`
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
		return image{Kind: "symlink", Link: target, identity: info}, err
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
	return image{Kind: "file", Mode: current.Mode(), Digest: digestBytes(data), Data: data, identity: current}, nil
}

func same(a, b image) bool {
	return a.Kind == b.Kind && a.Mode == b.Mode && a.Link == b.Link && a.Digest == b.Digest
}

func sameObserved(a, b image) bool {
	if !same(a, b) {
		return false
	}
	if a.identity == nil || b.identity == nil {
		return a.identity == nil && b.identity == nil
	}
	return os.SameFile(a.identity, b.identity)
}

// SyncDirectory reports both persistence and close failures for a directory.
func SyncDirectory(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return errors.Join(f.Sync(), f.Close())
}

func publish(path string, value image, persist func(string) error, rename func(string, string) error) error {
	parent := filepath.Dir(path)
	if err := ensureDirectory(parent, 0o755, persist); err != nil {
		return err
	}
	if value.Kind == "absent" {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return persist(parent)
	}
	f, err := os.CreateTemp(parent, ".bench-adopt-")
	if err != nil {
		return err
	}
	temp := f.Name()
	defer os.Remove(temp)
	if value.Kind == "symlink" {
		if err := f.Close(); err != nil {
			return err
		}
		if err := os.Remove(temp); err != nil {
			return err
		}
		if err := os.Symlink(value.Link, temp); err != nil {
			return err
		}
	} else if value.Kind == "file" {
		_, err = f.Write(value.Data)
		if err == nil {
			err = f.Chmod(value.Mode)
		}
		if err == nil {
			err = f.Sync()
		}
		if err = errors.Join(err, f.Close()); err != nil {
			return err
		}
	} else {
		_ = f.Close()
		return fmt.Errorf("unknown image kind %q", value.Kind)
	}
	if err := rename(temp, path); err != nil {
		return err
	}
	return persist(parent)
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
