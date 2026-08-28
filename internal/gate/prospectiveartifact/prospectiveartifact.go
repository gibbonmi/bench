// Package prospectiveartifact owns temporary prospective gate artifacts.
package prospectiveartifact

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

const (
	bundlePrefix    = "bench-prospective-artifact-"
	ownerRecordName = "owner.json"
	checkoutName    = "checkout"
)

type ownerRecord struct {
	Schema    int    `json:"schema"`
	OwnerPID  int    `json:"owner_pid"`
	CommonDir string `json:"common_dir"`
}

// Owner holds one prospective checkout bundle.
type Owner struct {
	root       string
	repository string
	remove     func(string, string) error
	closeOnce  sync.Once
	closeErr   error
}

// Root returns the bundle root.
func (o *Owner) Root() string { return o.root }

// Checkout returns the fixed checkout child path.
func (o *Owner) Checkout() string { return filepath.Join(o.root, checkoutName) }

// Factory supplies the operating-system seams for a bundle owner.
type Factory struct {
	TempRoot string
	Probe    func(int) error
	Remove   func(string, string) error
}

// Open creates a bundle owner for repository.
func Open(repository string) (*Owner, error) {
	return Factory{}.Open(repository)
}

// Open recovers dead bundles before it publishes a new owner record.
func (f Factory) Open(repository string) (*Owner, error) {
	common, err := canonicalCommonDir(repository)
	if err != nil {
		return nil, err
	}
	base := f.TempRoot
	if base == "" {
		base = os.TempDir()
	}
	remove := f.Remove
	if remove == nil {
		remove = removeWorktree
	}
	if err := sweep(base, repository, common, f.probe(), remove); err != nil {
		return nil, err
	}
	root, err := os.MkdirTemp(base, bundlePrefix)
	if err != nil {
		return nil, fmt.Errorf("create prospective artifact bundle: %w", err)
	}
	o := &Owner{root: root, repository: repository, remove: remove}
	record := ownerRecord{Schema: 1, OwnerPID: os.Getpid(), CommonDir: common}
	if err := writeRecord(root, record); err != nil {
		_ = os.RemoveAll(root)
		return nil, err
	}
	return o, nil
}

// Materialize creates the private checkout and replaces its tree with tree.
func (o *Owner) Materialize(tree string) error {
	if output, err := exec.Command("git", "-C", o.repository, "worktree", "add", "--quiet", "--detach", o.Checkout(), "HEAD").CombinedOutput(); err != nil {
		return fmt.Errorf("create prospective checkout: %s", strings.TrimSpace(string(output)))
	}
	if output, err := exec.Command("git", "-C", o.Checkout(), "read-tree", "--reset", "-u", tree).CombinedOutput(); err != nil {
		return fmt.Errorf("materialize prospective tree: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

// Close removes the Git registration before it removes the owned bundle root.
func (o *Owner) Close() error {
	if o == nil {
		return nil
	}
	o.closeOnce.Do(func() { o.closeErr = removeBundle(o.repository, o.root, o.Checkout(), o.remove) })
	return o.closeErr
}

func (f Factory) probe() func(int) error {
	if f.Probe != nil {
		return f.Probe
	}
	return func(pid int) error { return syscall.Kill(pid, 0) }
}

func canonicalCommonDir(repository string) (string, error) {
	common, err := benchgit.CommonDir(repository)
	if err != nil {
		return "", fmt.Errorf("resolve prospective repository: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(common)
	if err != nil {
		return "", fmt.Errorf("resolve prospective repository: %w", err)
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.IsDir() {
		return "", errors.New("resolve prospective repository: common directory is unavailable")
	}
	return filepath.Clean(resolved), nil
}

func sweep(base, repository, common string, probe func(int) error, remove func(string, string) error) error {
	entries, err := os.ReadDir(base)
	if err != nil {
		return fmt.Errorf("scan prospective artifact bundles: %w", err)
	}
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), bundlePrefix) {
			continue
		}
		root := filepath.Join(base, entry.Name())
		info, err := os.Lstat(root)
		if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			continue
		}
		record, ok := readRecord(root)
		if !ok || record.CommonDir != common {
			continue
		}
		if err := probe(record.OwnerPID); err == nil || errors.Is(err, syscall.EPERM) {
			continue
		} else if !errors.Is(err, syscall.ESRCH) {
			continue
		}
		if err := removeBundle(repository, root, filepath.Join(root, checkoutName), remove); err != nil {
			return fmt.Errorf("recover prospective artifact bundle: %w", err)
		}
	}
	return nil
}

func readRecord(root string) (ownerRecord, bool) {
	path := filepath.Join(root, ownerRecordName)
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		return ownerRecord{}, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ownerRecord{}, false
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	var record ownerRecord
	if decoder.Decode(&record) != nil || record.Schema != 1 || record.OwnerPID <= 0 || record.CommonDir == "" {
		return ownerRecord{}, false
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return ownerRecord{}, false
	}
	return record, true
}

func writeRecord(root string, record ownerRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(root, ".owner-")
	if err != nil {
		return err
	}
	name := temporary.Name()
	defer func() { _ = temporary.Close(); _ = os.Remove(name) }()
	if err := temporary.Chmod(0o600); err != nil {
		return err
	}
	if _, err := temporary.Write(append(data, '\n')); err != nil {
		return err
	}
	if err := temporary.Sync(); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(name, filepath.Join(root, ownerRecordName)); err != nil {
		return err
	}
	directory, err := os.Open(root)
	if err != nil {
		return err
	}
	defer directory.Close()
	return directory.Sync()
}

func removeBundle(repository, root, checkout string, remove func(string, string) error) error {
	worktrees, err := benchgit.Worktrees(repository)
	if err != nil {
		return err
	}
	for _, worktree := range worktrees {
		if filepath.Clean(worktree.Path) == filepath.Clean(checkout) {
			if err := remove(repository, checkout); err != nil {
				return err
			}
			break
		}
	}
	return os.RemoveAll(root)
}

func removeWorktree(repository, checkout string) error {
	output, err := exec.Command("git", "-C", repository, "worktree", "remove", "--force", checkout).CombinedOutput()
	if err != nil {
		return fmt.Errorf("remove prospective checkout: %s", strings.TrimSpace(string(output)))
	}
	return nil
}
