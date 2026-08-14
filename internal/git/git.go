// Package git is the one source of git subprocess invocation for the AXI query
// commands. Every ported parser shells out to git through here — git stays the
// source of repository truth (root, config, merge-base, diff), exactly as the shell
// commands did, so there is one place the invocation form lives.
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
)

// refCheckTimeout bounds the ref/branch existence probes the destructive-git guard
// runs per classification (internal/gitguard's checkout and forced-creation verdicts):
// a hung git must never stall a PreToolUse Bash hook, so each probe is bounded at two
// seconds and then resolves to its caller's fail-safe default.
const refCheckTimeout = 2 * time.Second

// RefResolves reports whether arg names a commit-ish that resolves in the process
// working directory — the agent's cwd, where the guarded Bash command would run, not a
// fixed root (`git rev-parse --verify --quiet <arg>^{commit}`). On a timeout or a
// failure to run git at all it returns false: an undeterminable target is treated as
// unresolvable, so a checkout of it fails closed (blocks) — the fail-toward-blocking
// default is deliberate, not incidental.
func RefResolves(arg string) bool {
	exitZero, ran := refCheck(arg + "^{commit}")
	if !ran {
		return false
	}
	return exitZero
}

// BranchExists reports whether refs/heads/<name> exists in the cwd repo. On a timeout
// or a failure to run git it returns TRUE — the opposite default from RefResolves —
// because its only caller (forced branch/switch creation) must block when it cannot
// rule out that the force would clobber an existing branch.
func BranchExists(name string) bool {
	exitZero, ran := refCheck("refs/heads/" + name)
	if !ran {
		return true
	}
	return exitZero
}

// refCheck runs `git rev-parse --verify --quiet <ref>` under the time bound and reports
// (git exited zero, git ran to a verdict). ran is false when the probe timed out or
// git could not be executed — the undeterminable branch each caller resolves to its own
// fail-safe default.
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

// Root returns the working tree's top-level directory, or an error when the cwd is
// not inside a git repository (the `not in a git repository` posture of every command).
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
// DefaultResolved is what makes the default-branch cells readable: when it is false,
// DefaultBranch is empty and Ahead/Behind are zero because there is no branch to measure
// against — an unknown, not a measurement.
type RepoFacts struct {
	Branch, DefaultBranch string
	DefaultResolved       bool
	Dirty                 bool
	Ahead, Behind         int
	Changes               []PorcelainEntry
}

// DiffFacts is the additive facts path for bench diff. It intentionally leaves Facts
// unchanged for existing consumers while expanding untracked directories into the
// individual entries a coherent patch can actually show.
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

// Worktrees returns every registered checkout using NUL-framed porcelain. The
// framing is required because a valid worktree path may contain a newline.
func Worktrees(root string) ([]Worktree, error) {
	raw, err := Raw("-C", root, "worktree", "list", "--porcelain", "-z")
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

// PruneLandedBranches removes local non-default branches whose work is already
// contained in the default branch and which are neither checked out in a registered
// worktree nor named by a caller-protected lifecycle record. The exact old OID makes
// each deletion fail closed if the branch moves after classification.
func PruneLandedBranches(root string, protectedBranches []string) (int, error) {
	def, ok := ResolvedDefault(root)
	if !ok {
		return 0, errors.New("git repository has no resolvable default branch")
	}
	worktrees, err := Worktrees(root)
	if err != nil {
		return 0, fmt.Errorf("git worktree list: %w", err)
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

// LandedState derives checkout-local dirtiness and repository-wide commit and branch facts.
// excludedDirtyPaths are omitted only from the named checkout's dirty count.
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
// that case: a repository with no commits still has a named branch, and losing the whole
// snapshot over a missing commit is the worse answer. It is exported because the probe
// chain, not the phrasing built from it, is what every caller shares — a caller that wants
// detachment reported as "no branch" tests the returned literal rather than running the
// two git queries a second time.
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
		// Divergence is derived from `rev-list <default>...HEAD`, which errors against a
		// branch that does not exist; the caller gets the unresolved state and the rest of
		// the snapshot instead of a failed read of the whole thing.
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

// AllFilesStatus returns Git's raw all-files porcelain plus untracked special
// entries that Git omits from that stream. The supplemental walk only classifies
// non-regular, non-directory, non-symlink nodes and asks Git whether each is tracked
// or ignored; it never opens the node.
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
// The -z framing is NUL-delimited and never C-quotes, so a path with spaces, glob
// characters, or a literal newline survives whole — the one source of that framing
// knowledge for every caller (the shift staging diff and the commit block-check).
func ParsePorcelainZ(raw []byte) []PorcelainEntry {
	entries, _ := ParsePorcelainZStrict(raw)
	return entries
}

// ParsePorcelainZStrict parses and validates NUL-framed porcelain-v1 records.
// Rename and copy records carry a second, path-only NUL record; it is returned
// with an empty Status so callers can preserve the framing while filtering.
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
// reverse-applying the branch's cumulative diff to def's tree — the last is what
// proves a squash-landing, where the branch's commits were composed into one commit
// no other proof can see. Merge-only content cannot be proven by git cherry and is
// deliberately kept; the merge check runs before either content proof so that
// conservatism holds for the fourth proof too.
//
// A true verdict is the sole authority every cleanup path deletes a branch on, so
// ambiguity never rounds up: whatever the reverse-apply proof cannot generate, apply
// cleanly, or represent byte-for-byte (reverseAppliesToDefault names the forms) is
// reported not landed. The cost of refusing is an orphaned branch; the cost of a
// wrong true is lost work with nothing standing behind it.
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

// Output runs `git <args>` and returns stdout with a single trailing newline trimmed;
// err is non-nil on a nonzero exit. Used for single-value reads (root, a config key,
// a resolved sha) where the trailing newline is noise.
func Output(args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimRight(out.String(), "\n"), err
}

// OK reports whether `git <args>` exits zero, discarding all output — the test form
// (cat-file -e, merge-base --is-ancestor) where only the exit code matters.
func OK(args ...string) bool {
	return exec.Command("git", args...).Run() == nil
}

// Raw runs `git <args>` and returns stdout verbatim (no trimming) with the exit
// status; used for `diff -z` output whose NUL framing and any trailing bytes are load-bearing.
func Raw(args ...string) ([]byte, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}
