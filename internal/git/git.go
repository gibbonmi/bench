// Package git is the one source of git subprocess invocation for the AXI query
// commands. Every ported parser shells out to git through here. Git stays the source of
// repository truth — root, config, merge-base, diff — exactly as the shell commands did,
// so the invocation form lives in one place.
package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// refCheckTimeout is the hook-scoped fail-safe for destructive-git classification,
// unlike the policy-owned worktree discovery bound below. It bounds the ref and branch
// existence probes the destructive-git guard runs per classification — internal/gitguard's
// checkout and forced-creation verdicts. A hung git must never stall a PreToolUse Bash
// hook, so the code bounds each probe at two seconds and resolves it to its caller's
// fail-safe default.
const refCheckTimeout = 2 * time.Second

var worktreeListTimeout = bounds.WorktreeListTimeout

// SetWorktreeListTimeoutForTest installs a test-only discovery bound and restores it.
func SetWorktreeListTimeoutForTest(limit time.Duration) func() {
	previous := worktreeListTimeout
	worktreeListTimeout = limit
	return func() { worktreeListTimeout = previous }
}

// RefResolves reports whether arg names a commit-ish that resolves in the process
// working directory. That directory is the agent's cwd, where the guarded Bash command
// runs, not a fixed root: the check runs `git rev-parse --verify --quiet
// <arg>^{commit}`. On a timeout, or a failure to run git at all, it returns false. The
// function treats an undeterminable target as unresolvable, so a checkout of it fails
// closed and blocks; the fail-toward-blocking default is deliberate.
func RefResolves(arg string) bool {
	exitZero, ran := refCheck(arg + "^{commit}")
	if !ran {
		return false
	}
	return exitZero
}

// BranchExists reports whether refs/heads/<name> exists in the cwd repo. On a timeout or
// a failure to run git, it returns TRUE, the opposite default from RefResolves. Its only
// caller — forced branch or switch creation — must block when it cannot rule out that the
// force would clobber an existing branch.
func BranchExists(name string) bool {
	exitZero, ran := refCheck("refs/heads/" + name)
	if !ran {
		return true
	}
	return exitZero
}

// refCheck runs `git rev-parse --verify --quiet <ref>` under the time bound and reports
// (git exited zero, git ran to a verdict). ran is false when the probe timed out or git
// could not run — the undeterminable branch each caller resolves to its own fail-safe
// default.
func refCheck(ref string) (exitZero, ran bool) {
	ctx, cancel := context.WithTimeout(context.Background(), refCheckTimeout)
	defer cancel()
	err := exec.CommandContext(ctx, "git", "rev-parse", "--verify", "--quiet", ref).Run()
	if ctx.Err() == context.DeadlineExceeded {
		return false, false
	}
	if err == nil {
		return true, true
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return false, true
	}
	return false, false
}

// Root returns the working tree's top-level directory. It returns an error when the cwd
// is not inside a git repository — the `not in a git repository` posture of every
// command.
func Root() (string, error) {
	return Output("rev-parse", "--show-toplevel")
}

// RootAt returns the top-level directory for the repository containing dir. Harness
// events carry an explicit cwd, so their routing must not depend on the hook process's
// ambient working directory.
func RootAt(dir string) (string, error) {
	return Output("-C", dir, "rev-parse", "--show-toplevel")
}

// PorcelainEntry is one record of `git status --porcelain -z --no-renames`: the two
// status characters (XY) and the path. With --no-renames a record is always
// `XY <path>`, so the path begins at byte 3.
type PorcelainEntry struct {
	Status string // the two XY status characters
	Path   string // the record's path, verbatim
}

// RepoFacts is the typed local repository state used by read-only query owners.
// DefaultResolved makes the default-branch cells readable. When it is false,
// DefaultBranch is empty and Ahead/Behind are zero, because there is no branch to
// measure against — an unknown, not a measurement.
type RepoFacts struct {
	Branch, DefaultBranch string
	DefaultResolved       bool
	Dirty                 bool
	Ahead, Behind         int
	Changes               []PorcelainEntry
}

// DiffFacts is the additive facts path for bench diff. It leaves Facts unchanged for
// existing consumers. It expands untracked directories into the individual entries a
// coherent patch can actually show.
type DiffFacts struct {
	RepoFacts
	Head, DefaultTip, RecordedBase string
	Porcelain                      []byte
}

// LandedStateFact is the offline git verdict used by status. DirtyPaths describes the named
// checkout; commit and branch counts describe the repository.
type LandedStateFact struct {
	DirtyPaths      int
	UnpushedCommits int
	UniqueBranches  int
}

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

// LocalBranches is the sole local-head enumeration and parsing owner.
func LocalBranches(root string) ([]string, error) {
	out, err := Output("-C", root, "for-each-ref", "--format=%(refname:short)", "refs/heads/")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return []string{}, nil
	}
	return strings.Split(out, "\n"), nil
}

// DeleteBranchExact removes one full branch ref only while it still has the
// caller-proven OID.
func DeleteBranchExact(root, ref, oid string) error {
	out, err := exec.Command("git", "-C", root, "update-ref", "-d", ref, oid).CombinedOutput()
	if err != nil {
		return fmt.Errorf("update exact branch ref: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// PruneLandedBranches removes local non-default branches whose work the default branch
// already contains. It skips a branch checked out in a registered worktree or named by a
// caller-protected lifecycle record. The exact old OID makes each deletion fail closed if
// the branch moves after classification.
func PruneLandedBranches(root string, protectedBranches []string) (int, error) {
	def, ok := ResolvedDefault(root)
	if !ok {
		return 0, errors.New("git repository has no resolvable default branch")
	}
	worktrees, err := Worktrees(root)
	if err != nil {
		return 0, fmt.Errorf("worktree discovery failed: %w", err)
	}
	checkedOut := map[string]bool{}
	for _, worktree := range worktrees {
		if worktree.Branch != "" {
			checkedOut[worktree.Branch] = true
		}
	}
	for _, branch := range protectedBranches {
		checkedOut[strings.TrimPrefix(branch, "refs/heads/")] = true
	}
	branches, err := LocalBranches(root)
	if err != nil {
		return 0, fmt.Errorf("git local branches: %w", err)
	}
	pruned := 0
	for _, branch := range branches {
		if branch == def || checkedOut[branch] {
			continue
		}
		landed, _, err := LandedInDefault(root, branch, def)
		if err != nil {
			return pruned, fmt.Errorf("git landedness %s: %w", branch, err)
		}
		if !landed {
			continue
		}
		oid, err := Output("-C", root, "rev-parse", "--verify", "refs/heads/"+branch+"^{commit}")
		if err != nil {
			return pruned, fmt.Errorf("git branch identity %s: %w", branch, err)
		}
		if err := DeleteBranchExact(root, "refs/heads/"+branch, oid); err != nil {
			return pruned, fmt.Errorf("delete landed branch %s: %w", branch, err)
		}
		pruned++
	}
	return pruned, nil
}

// LandedState derives checkout-local dirtiness and repository-wide commit and branch
// facts. The dirty count omits excludedDirtyPaths, but only for the named checkout.
func LandedState(root string, excludedDirtyPaths ...string) (LandedStateFact, error) {
	excluded := make(map[string]bool, len(excludedDirtyPaths))
	for _, path := range excludedDirtyPaths {
		excluded[path] = true
	}
	dirty := map[string]bool{}
	raw, err := Raw("-C", root, "status", "--porcelain=v1", "-z", "--no-renames")
	if err != nil {
		return LandedStateFact{}, fmt.Errorf("git status %s: %w", root, err)
	}
	for _, entry := range ParsePorcelainZ(raw) {
		if excluded[entry.Path] {
			continue
		}
		dirty[entry.Path] = true
	}
	def, ok := ResolvedDefault(root)
	if !ok {
		return LandedStateFact{}, errors.New("git repository has no resolvable default branch")
	}
	branches, err := LocalBranches(root)
	if err != nil {
		return LandedStateFact{}, fmt.Errorf("git local branches: %w", err)
	}
	commits := map[string]bool{}
	unique := map[string]bool{}
	for _, branch := range branches {
		if branch != def {
			landed, _, err := LandedInDefault(root, branch, def)
			if err != nil {
				return LandedStateFact{}, fmt.Errorf("git landedness %s: %w", branch, err)
			}
			if !landed {
				unique[branch] = true
			}
		}
		upstream, err := Output("-C", root, "for-each-ref", "--format=%(upstream:short)", "refs/heads/"+branch)
		if err != nil {
			return LandedStateFact{}, fmt.Errorf("git upstream %s: %w", branch, err)
		}
		if upstream == "" {
			continue
		}
		ahead, err := Output("-C", root, "rev-list", upstream+".."+branch)
		if err != nil {
			return LandedStateFact{}, fmt.Errorf("git ahead %s: %w", branch, err)
		}
		for _, commit := range strings.Split(ahead, "\n") {
			if commit != "" {
				commits[commit] = true
			}
		}
	}
	return LandedStateFact{DirtyPaths: len(dirty), UnpushedCommits: len(commits), UniqueBranches: len(unique)}, nil
}

// CheckedOutBranch names the branch HEAD points at, or the literal "HEAD" when detached.
// `rev-parse --abbrev-ref` fails outright on an unborn branch, so the symbolic ref settles
// that case. A repository with no commits still has a named branch, and losing the whole
// snapshot over a missing commit is the worse answer. The function is exported because
// every caller shares the probe chain, not the phrasing built from it. A caller that wants
// detachment reported as "no branch" tests the returned literal instead of running the two
// git queries again.
func CheckedOutBranch(root string) (string, error) {
	if name, err := Output("-C", root, "rev-parse", "--abbrev-ref", "HEAD"); err == nil && name != "" {
		return name, nil
	}
	return Output("-C", root, "symbolic-ref", "--quiet", "--short", "HEAD")
}

// Facts derives repository state without mutating the worktree or index.
func Facts(root string) (RepoFacts, error) {
	branch, err := CheckedOutBranch(root)
	if err != nil {
		return RepoFacts{}, err
	}
	raw, err := Raw("-C", root, "status", "--porcelain=v1", "-z", "--no-renames")
	if err != nil {
		return RepoFacts{}, err
	}
	f := RepoFacts{Branch: branch, Changes: ParsePorcelainZ(raw)}
	f.Dirty = len(f.Changes) > 0
	def, ok := ResolvedDefault(root)
	if !ok {
		// The code derives divergence from `rev-list <default>...HEAD`, which errors against
		// a branch that does not exist. The caller gets the unresolved state and the rest
		// of the snapshot instead of a failed read of the whole thing.
		return f, nil
	}
	f.DefaultBranch, f.DefaultResolved = def, true
	counts, err := Output("-C", root, "rev-list", "--left-right", "--count", f.DefaultBranch+"...HEAD")
	if err != nil {
		return RepoFacts{}, fmt.Errorf("git rev-list: %w", err)
	}
	if _, err := fmt.Sscanf(counts, "%d\t%d", &f.Behind, &f.Ahead); err != nil {
		return RepoFacts{}, fmt.Errorf("parse git divergence: %w", err)
	}
	return f, nil
}

// AllFilesFacts derives the diff-specific status facts with Git's all-files
// untracked policy. Existing Facts callers retain their collapsed-directory output.
func AllFilesFacts(root string) (DiffFacts, error) {
	branch, err := CheckedOutBranch(root)
	if err != nil {
		return DiffFacts{}, err
	}
	head, _ := Output("-C", root, "rev-parse", "HEAD")
	recordedBase := ""
	if branch != "HEAD" {
		recordedBase, _ = Output("-C", root, "config", "branch."+branch+".benchBase")
	}
	raw, changes, err := AllFilesStatus(root)
	if err != nil {
		return DiffFacts{}, err
	}
	f := DiffFacts{
		RepoFacts:    RepoFacts{Branch: branch, Changes: changes},
		Head:         head,
		RecordedBase: recordedBase,
		Porcelain:    raw,
	}
	f.Dirty = len(f.Changes) > 0
	def, ok := ResolvedDefault(root)
	if !ok {
		return f, nil
	}
	defaultTip, err := Output("-C", root, "rev-parse", def)
	if err != nil {
		return DiffFacts{}, fmt.Errorf("git default tip: %w", err)
	}
	f.DefaultBranch, f.DefaultTip, f.DefaultResolved = def, defaultTip, true
	counts, err := Output("-C", root, "rev-list", "--left-right", "--count", defaultTip+"..."+head)
	if err != nil {
		return DiffFacts{}, fmt.Errorf("git rev-list: %w", err)
	}
	if _, err := fmt.Sscanf(counts, "%d\t%d", &f.Behind, &f.Ahead); err != nil {
		return DiffFacts{}, fmt.Errorf("parse git divergence: %w", err)
	}
	return f, nil
}

// AllFilesStatus returns Git's raw all-files porcelain plus untracked special entries
// that Git omits from that stream. The supplemental walk only classifies non-regular,
// non-directory, non-symlink nodes. It asks Git whether each is tracked or ignored, and
// it never opens the node.
func AllFilesStatus(root string) ([]byte, []PorcelainEntry, error) {
	raw, err := Raw("-C", root, "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all")
	if err != nil {
		return nil, nil, err
	}
	changes := ParsePorcelainZ(raw)
	seen := make(map[string]bool, len(changes))
	for _, change := range changes {
		seen[change.Path] = true
	}
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Name() == ".git" {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() || entry.Type().IsRegular() || entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if seen[rel] || OK("-C", root, "ls-files", "--error-unmatch", "--", rel) || OK("-C", root, "check-ignore", "-q", "--", rel) {
			return nil
		}
		seen[rel] = true
		changes = append(changes, PorcelainEntry{Status: "??", Path: rel})
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	return raw, changes, nil
}

// ParsePorcelainZ splits `git status --porcelain -z --no-renames` output into entries.
// The -z framing is NUL-delimited and never C-quotes. A path with spaces, glob
// characters, or a literal newline survives whole. This function is the one source of
// that framing knowledge for every caller — the shift staging diff and the commit
// block-check.
func ParsePorcelainZ(raw []byte) []PorcelainEntry {
	entries, _ := ParsePorcelainZStrict(raw)
	return entries
}

// ParsePorcelainZStrict parses and validates NUL-framed porcelain-v1 records. A rename
// or copy record carries a second, path-only NUL record. The parser returns that record
// with an empty Status, so callers can preserve the framing while filtering.
func ParsePorcelainZStrict(raw []byte) ([]PorcelainEntry, error) {
	var entries []PorcelainEntry
	for offset := 0; offset < len(raw); {
		end := bytes.IndexByte(raw[offset:], 0)
		if end < 0 {
			return nil, errors.New("malformed checkout status")
		}
		record := raw[offset : offset+end]
		offset += end + 1
		if len(record) < 4 || record[2] != ' ' {
			return nil, errors.New("malformed checkout status")
		}
		entries = append(entries, PorcelainEntry{Status: string(record[:2]), Path: string(record[3:])})
		status := record[:2]
		if status[0] == 'R' || status[0] == 'C' || status[1] == 'R' || status[1] == 'C' {
			end = bytes.IndexByte(raw[offset:], 0)
			if end < 1 {
				return nil, errors.New("malformed checkout status")
			}
			entries = append(entries, PorcelainEntry{Path: string(raw[offset : offset+end])})
			offset += end + 1
		}
	}
	return entries, nil
}

// GateCacheFile is the filename of the cached gate verdict, written under the git
// directory (never the worktree, so it is never a diff or commit candidate). It is the
// one gate owner resolves and composes with its absolute Git-directory path.
const GateCacheFile = "bench-last-gate"

// LandedInDefault proves a local branch landed by ancestry, patch containment, or
// reverse-applying the branch's cumulative diff to def's tree. The reverse-apply proof
// alone proves a squash-landing, where the branch's commits compose into one commit no
// other proof can see. git cherry cannot prove merge-only content, so the function keeps
// the merge check. It runs that check before either content proof, which keeps
// conservatism for the fourth proof too.
//
// A true verdict is the sole authority every cleanup path deletes a branch on, so
// ambiguity never rounds up. The function reports not landed whenever the reverse-apply
// proof cannot generate, apply cleanly, or represent the diff byte-for-byte —
// reverseAppliesToDefault names the exact forms. The cost of refusing is an orphaned
// branch. The cost of a wrong true is lost work with nothing standing behind it.
func LandedInDefault(root, branch, def string) (landed, byContent bool, err error) {
	ancestor := exec.Command("git", "-C", root, "merge-base", "--is-ancestor", branch, def)
	if err := ancestor.Run(); err == nil {
		return true, false, nil
	} else if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
		return false, false, fmt.Errorf("merge-base %s against %s: %w", branch, def, err)
	}
	merges, err := Output("-C", root, "rev-list", "--merges", "--max-count=1", def+".."+branch)
	if err != nil {
		return false, false, fmt.Errorf("list merges on %s: %w", branch, err)
	}
	if merges != "" {
		return false, false, nil
	}
	out, err := Output("-C", root, "cherry", def, branch)
	if err != nil {
		return false, false, fmt.Errorf("compare patches on %s: %w", branch, err)
	}
	for _, line := range strings.Split(out, "\n") {
		if line != "" && !strings.HasPrefix(line, "-") {
			if reverseAppliesToDefault(root, branch, def) {
				return true, true, nil
			}
			return false, false, nil
		}
	}
	return true, true, nil
}

// Output runs `git <args>` and returns stdout with a single trailing newline trimmed.
// err is non-nil on a nonzero exit. Callers use it for single-value reads — a root, a
// config key, a resolved sha — where the trailing newline is noise.
func Output(args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimRight(out.String(), "\n"), err
}

// OK reports whether `git <args>` exits zero. It discards all output. This is the test
// form — cat-file -e, merge-base --is-ancestor — where only the exit code matters.
func OK(args ...string) bool {
	return exec.Command("git", args...).Run() == nil
}

// Raw runs `git <args>` and returns stdout verbatim, with no trimming, along with the
// exit status. Callers use it for `diff -z` output, whose NUL framing and trailing bytes
// are load-bearing.
func Raw(args ...string) ([]byte, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}
