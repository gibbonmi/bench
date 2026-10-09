package transaction

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/canonicalpath"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

// Lease excludes competing writers for its canonical destinations.
type Lease struct {
	files []*os.File
	keys  map[string]bool
}

// Lock acquires the complete destination set before an adoption writer mutates it.
func Lock(destinations []string) (*Lease, error) {
	// A process-local temporary directory or Bench home would split independent writers.
	shared, err := filepath.EvalSymlinks("/tmp")
	if err != nil {
		return nil, err
	}
	return lockIn(filepath.Join(shared, "bench-adoption-locks-"+strconv.Itoa(os.Getuid())), destinations)
}

// lockIn acquires destinations through the lock files in dir.
func lockIn(dir string, destinations []string) (*Lease, error) {
	if err := os.Mkdir(dir, 0o700); err != nil && !os.IsExist(err) {
		return nil, err
	}
	info, err := os.Lstat(dir)
	if err != nil {
		return nil, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !info.IsDir() || info.Mode().Perm() != 0o700 || !ok || stat.Uid != uint32(os.Getuid()) {
		return nil, fmt.Errorf("invalid adoption lock directory: %s", dir)
	}
	keys := make(map[string]bool, len(destinations))
	for _, dest := range destinations {
		key, err := destinationKey(dest)
		if err != nil {
			return nil, err
		}
		keys[key] = true
	}
	ordered := make([]string, 0, len(keys))
	for key := range keys {
		ordered = append(ordered, key)
	}
	sort.Strings(ordered)
	lease := &Lease{keys: keys}
	for _, key := range ordered {
		f, err := acquire(filepath.Join(dir, key))
		if err != nil {
			return nil, errors.Join(fmt.Errorf("adoption destination is unavailable; retry after the competing writer completes: %w", err), lease.Close())
		}
		lease.files = append(lease.files, f)
	}
	return lease, nil
}

// acquire holds the lock file at path. Only a holder removes a lock file, so a lock taken
// on an inode that path no longer names belongs to a released holder and is retried.
func acquire(path string) (*os.File, error) {
	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0o600)
		if err != nil {
			return nil, err
		}
		held, err := f.Stat()
		if err != nil || !held.Mode().IsRegular() || held.Mode().Perm() != 0o600 {
			err = fmt.Errorf("invalid adoption lock file")
		}
		if err == nil {
			err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		}
		if err == nil {
			current, statErr := os.Lstat(path)
			if statErr == nil && os.SameFile(held, current) {
				return f, nil
			}
			if statErr != nil && !os.IsNotExist(statErr) {
				err = statErr
			}
		}
		_ = f.Close()
		if err != nil {
			return nil, err
		}
	}
}

func canonicalDestination(path string) (string, error) {
	directory, name := filepath.Split(path)
	if name == "" {
		return "", fmt.Errorf("managed destination has no file name")
	}
	if directory == "" {
		directory = "."
	}
	parent, err := resolveDirectory(directory)
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, name), nil
}

func resolveDirectory(path string) (string, error) {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("not a directory: %s", path)
		}
		return canonicalpath.Resolve(path)
	}
	if !os.IsNotExist(err) {
		return "", err
	}
	if info, statErr := os.Lstat(path); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
		return "", err
	}
	parent, name := filepath.Split(strings.TrimRight(path, string(filepath.Separator)))
	if name == "" || name == "." || name == ".." {
		return "", err
	}
	if parent == "" {
		parent = "."
	}
	resolved, err := resolveDirectory(parent)
	if err != nil {
		return "", err
	}
	return filepath.Join(resolved, name), nil
}

// Close releases this writer's destinations. It removes each lock file while it still holds
// the lock, so the directory does not grow and acquire can detect the released inode.
func (l *Lease) Close() error {
	var failures []error
	for i := len(l.files) - 1; i >= 0; i-- {
		failures = append(failures, os.Remove(l.files[i].Name()), l.files[i].Close())
	}
	l.files = nil
	l.keys = nil
	return errors.Join(failures...)
}

func destinationKey(path string) (string, error) {
	canonical, err := canonicalDestination(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(canonical))), nil
}

func (l *Lease) covers(paths []string) error {
	for _, path := range paths {
		key, err := destinationKey(path)
		if err != nil {
			return err
		}
		if !l.keys[key] {
			return fmt.Errorf("adoption plan changed outside its held destinations: %s", path)
		}
	}
	return nil
}
