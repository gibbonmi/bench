package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// Worktree is one registered checkout reported by git. Worktrees is the sole
// parser for worktree-list porcelain so every consumer agrees on path framing,
// branch identity, detached state, and locks.
type Worktree struct {
	Path       string
	Branch     string
	BranchRef  string
	Detached   bool
	Locked     bool
	LockReason string
}

// BenchLeaseFilename names the private lifecycle record shared by git and pool lifecycle.
const BenchLeaseFilename = "bench-lease"

const investigateGitFailureAction = "investigate the git failure"

// ResolutionError reports a common-directory resolution that cannot be trusted.
type ResolutionError struct {
	Path   string
	Err    error
	Action string
}

// WorktreeFailure exposes the recovery action for a worktree discovery failure.
type WorktreeFailure interface {
	error
	WorktreeAction() string
}

// WorktreeAdminError identifies an admin entry whose filesystem shape is not
// safe for git's worktree enumeration. The code only lstats the entry.
type WorktreeAdminError struct {
	Path   string
	Shape  string
	Action string
}

func (e *WorktreeAdminError) Error() string {
	return fmt.Sprintf("worktree admin entry %q has shape %s; %s", e.Path, e.Shape, e.Action)
}

func (e *WorktreeAdminError) WorktreeAction() string { return e.Action }

// WorktreeScanError reports a failure while inspecting a worktree admin entry.
type WorktreeScanError struct {
	Path   string
	Err    error
	Action string
}

func (e *WorktreeScanError) Error() string {
	return fmt.Sprintf("cannot inspect worktree admin entry %q: %v; %s", e.Path, e.Err, e.Action)
}

func (e *WorktreeScanError) Unwrap() error { return e.Err }

func (e *WorktreeScanError) WorktreeAction() string { return e.Action }

// ScanWorktreeAdmin refuses malformed entries before git can open them. Every
// direct entry must be a regular file or directory.
// Git 2.43.0 blocks open-for-read on FIFO admin files, which is why this check remains a
// preflight until git bounds those reads itself.
func ScanWorktreeAdmin(commonDir string) error {
	base := filepath.Join(commonDir, "worktrees")
	info, err := os.Lstat(base)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return worktreeScanError("worktrees", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return worktreeAdminError("worktrees", "symlink")
	}
	if !info.IsDir() {
		return nil
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		return worktreeScanError("worktrees", err)
	}
	for _, id := range entries {
		idPath := filepath.Join(base, id.Name())
		idInfo, err := os.Lstat(idPath)
		if err != nil {
			return worktreeScanError(filepath.Join("worktrees", id.Name()), err)
		}
		if idInfo.Mode().IsRegular() {
			continue
		}
		if !idInfo.IsDir() || idInfo.Mode()&os.ModeSymlink != 0 {
			return worktreeAdminError(filepath.Join("worktrees", id.Name()), fileShape(idInfo.Mode()))
		}
		children, err := os.ReadDir(idPath)
		if err != nil {
			return worktreeScanError(filepath.Join("worktrees", id.Name()), err)
		}
		for _, child := range children {
			childPath := filepath.Join(idPath, child.Name())
			childInfo, err := os.Lstat(childPath)
			if err != nil {
				return worktreeScanError(filepath.Join("worktrees", id.Name(), child.Name()), err)
			}
			if childInfo.Mode().IsRegular() || childInfo.IsDir() {
				continue
			}
			return worktreeAdminError(filepath.Join("worktrees", id.Name(), child.Name()), fileShape(childInfo.Mode()))
		}
	}
	return nil
}

func worktreeAdminError(path, shape string) error {
	return &WorktreeAdminError{Path: path, Shape: shape, Action: "inspect and remove it"}
}

func worktreeScanError(path string, err error) error {
	return &WorktreeScanError{Path: path, Err: err, Action: investigateGitFailureAction}
}

func fileShape(mode os.FileMode) string {
	switch {
	case mode&os.ModeSymlink != 0:
		return "symlink"
	case mode&os.ModeNamedPipe != 0:
		return "fifo"
	case mode&os.ModeSocket != 0:
		return "socket"
	case mode&os.ModeDevice != 0:
		return "device"
	default:
		return "non-regular"
	}
}

func (e *ResolutionError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("git common directory %q has shape %v; %s", e.Path, e.Err, e.Action)
	}
	return fmt.Sprintf("git common directory resolution failed (%s); %s", e.Err, e.Action)
}

func (e *ResolutionError) Unwrap() error { return e.Err }

func (e *ResolutionError) WorktreeAction() string { return e.Action }

func commonDirArgs(root string) []string {
	return []string{"-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir"}
}

// CommonDir resolves the repository's shared administrative directory.
func CommonDir(root string) (string, error) {
	return Output(commonDirArgs(root)...)
}

func validateCommonDir(common string) (string, error) {
	if common == "" {
		return "", &ResolutionError{Err: errors.New("rev-parse returned an empty path"), Action: investigateGitFailureAction}
	}
	info, statErr := os.Lstat(common)
	if statErr != nil {
		return "", &ResolutionError{Path: common, Err: fmt.Errorf("missing path: %w", statErr), Action: investigateGitFailureAction}
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", &ResolutionError{Path: common, Err: errors.New("symlink"), Action: investigateGitFailureAction}
	}
	if !info.IsDir() {
		return "", &ResolutionError{Path: common, Err: errors.New("non-directory"), Action: investigateGitFailureAction}
	}
	return common, nil
}

func boundedGit(args ...string) ([]byte, error) {
	var stdout bytes.Buffer
	result := bounds.RunOutput(context.Background(), worktreeListTimeout, exec.Command("git", args...), &stdout)
	if result.Status == bounds.ProcessComplete {
		return stdout.Bytes(), nil
	}
	invocation := "git " + strings.Join(args, " ")
	if result.Status == bounds.ProcessExit {
		failure := fmt.Errorf("%s: %w", invocation, result.Err)
		if diagnostic := bytes.TrimRight(result.Output, "\r\n"); len(diagnostic) > 0 {
			failure = fmt.Errorf("%w: %s", failure, diagnostic)
		}
		return stdout.Bytes(), failure
	}
	if result.Status == bounds.ProcessTimeout {
		return nil, &ResolutionError{Err: fmt.Errorf("%s timed out after %s", invocation, worktreeListTimeout), Action: investigateGitFailureAction}
	}
	return nil, &ResolutionError{Err: fmt.Errorf("%s failed to start or was canceled: %w", invocation, result.Err), Action: investigateGitFailureAction}
}

// Worktrees returns every registered checkout using NUL-framed porcelain. The framing
// matters because a valid worktree path may contain a newline.
func Worktrees(root string) ([]Worktree, error) {
	commonRaw, err := boundedGit(commonDirArgs(root)...)
	if err != nil {
		if _, ok := err.(*ResolutionError); ok {
			return nil, err
		}
		return nil, &ResolutionError{Err: err, Action: investigateGitFailureAction}
	}
	common, err := validateCommonDir(strings.TrimRight(string(commonRaw), "\n"))
	if err != nil {
		return nil, err
	}
	if err := ScanWorktreeAdmin(common); err != nil {
		return nil, err
	}
	raw, err := boundedGit("-C", root, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		var typed *ResolutionError
		if errors.As(err, &typed) {
			return nil, err
		}
		return nil, err
	}
	var worktrees []Worktree
	var current *Worktree
	for field := range bytes.SplitSeq(raw, []byte{0}) {
		if len(field) == 0 {
			current = nil
			continue
		}
		line := string(field)
		switch {
		case strings.HasPrefix(line, "worktree "):
			worktrees = append(worktrees, Worktree{Path: strings.TrimPrefix(line, "worktree ")})
			current = &worktrees[len(worktrees)-1]
		case current != nil && strings.HasPrefix(line, "branch refs/heads/"):
			current.BranchRef = strings.TrimPrefix(line, "branch ")
			current.Branch = strings.TrimPrefix(current.BranchRef, "refs/heads/")
		case current != nil && line == "detached":
			current.Detached = true
		case current != nil && (line == "locked" || strings.HasPrefix(line, "locked ")):
			current.Locked = true
			current.LockReason = strings.TrimPrefix(line, "locked ")
		}
	}
	return worktrees, nil
}
