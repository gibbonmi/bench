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

// The three subject nouns a ResolutionError can carry. Each reader names its own
// fact, so one error type still serves every WorktreeFailure type check.
const (
	subjectCommonDir = "git common directory"
	subjectAdminDir  = "checkout administration directory"
	subjectAdminPath = "checkout administration path"
)

// ResolutionError reports a git administration fact that cannot be trusted. Subject
// names which of the three facts failed; an empty Subject renders as the common
// directory, the field's original sole meaning.
type ResolutionError struct {
	Path    string
	Err     error
	Action  string
	Subject string
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
	subject := e.Subject
	if subject == "" {
		subject = subjectCommonDir
	}
	if e.Path != "" {
		return fmt.Sprintf("%s %q has shape %v; %s", subject, e.Path, e.Err, e.Action)
	}
	return fmt.Sprintf("%s resolution failed (%s); %s", subject, e.Err, e.Action)
}

func (e *ResolutionError) Unwrap() error { return e.Err }

func (e *ResolutionError) WorktreeAction() string { return e.Action }

func commonDirArgs(root string) []string {
	return []string{"-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir"}
}

func adminDirArgs(root string) []string {
	return []string{"-C", root, "rev-parse", "--path-format=absolute", "--git-dir"}
}

func adminPathArgs(root, name string) []string {
	return []string{"-C", root, "rev-parse", "--git-path", name}
}

// CommonDir resolves the repository's shared administrative directory, refusing
// an empty, missing, symlinked, or non-directory answer the way Worktrees does.
func CommonDir(root string) (string, error) {
	raw, err := boundedGit(subjectCommonDir, commonDirArgs(root)...)
	if err != nil {
		return "", typeGitFailure(err, subjectCommonDir)
	}
	return validateCommonDir(strings.TrimRight(string(raw), "\n"), subjectCommonDir)
}

// AdminDir resolves the checkout's own administration directory: the common
// directory for a primary checkout, or the linked worktree's private admin entry
// beneath it. It validates its answer the way CommonDir does.
func AdminDir(root string) (string, error) {
	raw, err := boundedGit(subjectAdminDir, adminDirArgs(root)...)
	if err != nil {
		return "", typeGitFailure(err, subjectAdminDir)
	}
	return validateCommonDir(strings.TrimRight(string(raw), "\n"), subjectAdminDir)
}

// AdminPath resolves the absolute path of a named file inside the checkout's
// administration directory. It refuses an empty answer, and it runs no existence
// check, because absence is the caller's fact.
//
// The query uses git's default path format, not --path-format=absolute, because
// the absolute format resolves an existing symlink and would answer the target
// rather than the named file: the default path format keeps a symlinked entry's
// own path. The default format prints a path relative to git's working
// directory for a primary checkout, which is root, so a relative answer joins
// onto root. Root itself is resolved with filepath.EvalSymlinks before the join,
// because git's own -C resolves a symlinked path component physically: a root
// that reaches the repository through a symlink followed by ".." must join
// onto the same physical directory git ran in, not onto a lexically cleaned root
// that never saw the symlink.
func AdminPath(root, name string) (string, error) {
	raw, err := boundedGit(subjectAdminPath, adminPathArgs(root, name)...)
	if err != nil {
		return "", typeGitFailure(err, subjectAdminPath)
	}
	answer := strings.TrimRight(string(raw), "\n")
	if answer == "" {
		return "", &ResolutionError{Err: errors.New("rev-parse returned an empty path"), Action: investigateGitFailureAction, Subject: subjectAdminPath}
	}
	if filepath.IsAbs(answer) {
		return answer, nil
	}
	resolvedRoot, resolveErr := resolveAdminRoot(root)
	if resolveErr != nil {
		return "", resolveErr
	}
	return filepath.Join(resolvedRoot, answer), nil
}

// resolveAdminRoot resolves root to the physical, absolute directory git ran in.
// filepath.EvalSymlinks runs first, on the possibly-relative and possibly-unclean
// root, so a symlinked path component resolves physically before any ".." that
// follows it is applied — the same order the OS uses when git's -C changes into
// root. filepath.Abs then covers the case EvalSymlinks leaves relative, which it
// does not for an absolute input.
func resolveAdminRoot(root string) (string, error) {
	resolved, evalErr := filepath.EvalSymlinks(root)
	if evalErr != nil {
		return "", &ResolutionError{Path: root, Err: fmt.Errorf("cannot resolve symlinks in the checkout root: %w", evalErr), Action: investigateGitFailureAction, Subject: subjectAdminPath}
	}
	if filepath.IsAbs(resolved) {
		return resolved, nil
	}
	absoluteRoot, absErr := filepath.Abs(resolved)
	if absErr != nil {
		return "", &ResolutionError{Path: root, Err: fmt.Errorf("cannot absolutize the checkout root: %w", absErr), Action: investigateGitFailureAction, Subject: subjectAdminPath}
	}
	return absoluteRoot, nil
}

func validateCommonDir(common, subject string) (string, error) {
	if common == "" {
		return "", &ResolutionError{Err: errors.New("rev-parse returned an empty path"), Action: investigateGitFailureAction, Subject: subject}
	}
	info, statErr := os.Lstat(common)
	if statErr != nil {
		return "", &ResolutionError{Path: common, Err: fmt.Errorf("missing path: %w", statErr), Action: investigateGitFailureAction, Subject: subject}
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", &ResolutionError{Path: common, Err: errors.New("symlink"), Action: investigateGitFailureAction, Subject: subject}
	}
	if !info.IsDir() {
		return "", &ResolutionError{Path: common, Err: errors.New("non-directory"), Action: investigateGitFailureAction, Subject: subject}
	}
	return common, nil
}

// typeGitFailure wraps a boundedGit exit failure in a ResolutionError carrying
// subject, unless it is already typed — the timeout and start-failure branches of
// boundedGit return an already-typed error, and re-wrapping it would lose that
// subject.
func typeGitFailure(err error, subject string) error {
	if typed, ok := err.(*ResolutionError); ok {
		return typed
	}
	return &ResolutionError{Err: err, Action: investigateGitFailureAction, Subject: subject}
}

func boundedGit(subject string, args ...string) ([]byte, error) {
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
		return nil, &ResolutionError{Err: fmt.Errorf("%s timed out after %s", invocation, worktreeListTimeout), Action: investigateGitFailureAction, Subject: subject}
	}
	return nil, &ResolutionError{Err: fmt.Errorf("%s failed to start or was canceled: %w", invocation, result.Err), Action: investigateGitFailureAction, Subject: subject}
}

// Worktrees returns every registered checkout using NUL-framed porcelain. The framing
// matters because a valid worktree path may contain a newline.
func Worktrees(root string) ([]Worktree, error) {
	commonRaw, err := boundedGit(subjectCommonDir, commonDirArgs(root)...)
	if err != nil {
		return nil, typeGitFailure(err, subjectCommonDir)
	}
	common, err := validateCommonDir(strings.TrimRight(string(commonRaw), "\n"), subjectCommonDir)
	if err != nil {
		return nil, err
	}
	if err := ScanWorktreeAdmin(common); err != nil {
		return nil, err
	}
	raw, err := boundedGit(subjectCommonDir, "-C", root, "worktree", "list", "--porcelain", "-z")
	if err != nil {
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
