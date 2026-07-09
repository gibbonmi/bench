// Package git is the one source of git subprocess invocation for the AXI query
// commands. Every ported parser shells out to git through here — git stays the
// source of repository truth (root, config, merge-base, diff), exactly as the shell
// commands did, so there is one place the invocation form lives.
package git

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
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

// PorcelainEntry is one record of `git status --porcelain -z --no-renames`: the two
// status characters (XY) and the path. With --no-renames a record is always
// `XY <path>`, so the path begins at byte 3.
type PorcelainEntry struct {
	Status string // the two XY status characters
	Path   string // the record's path, verbatim
}

// ParsePorcelainZ splits `git status --porcelain -z --no-renames` output into entries.
// The -z framing is NUL-delimited and never C-quotes, so a path with spaces, glob
// characters, or a literal newline survives whole — the one source of that framing
// knowledge for every caller (the shift staging diff and the commit block-check).
func ParsePorcelainZ(raw []byte) []PorcelainEntry {
	var entries []PorcelainEntry
	for record := range bytes.SplitSeq(raw, []byte{0}) {
		if len(record) <= 3 {
			continue // trailing empty after the final NUL, or a malformed short record
		}
		entries = append(entries, PorcelainEntry{Status: string(record[:2]), Path: string(record[3:])})
	}
	return entries
}

// GateCacheFile is the filename of the cached gate verdict, written under the git
// directory (never the worktree, so it is never a diff or commit candidate). It is the
// one source gate.Record (the writer) and status.GateVerdict (the reader) both compose
// with their own git-dir resolution.
const GateCacheFile = "bench-last-gate"

// DefaultBranch is the repository's default branch: origin/HEAD's short name with the
// `origin/` prefix stripped, falling back to "main" when the ref is unset (no remote
// HEAD) or empty. The one source both `diff` and `status` read, so the two surfaces
// agree on what "default branch" is.
func DefaultBranch(root string) string {
	out, err := Output("-C", root, "symbolic-ref", "--quiet", "--short", "refs/remotes/origin/HEAD")
	if err != nil || out == "" {
		return "main"
	}
	return strings.TrimPrefix(out, "origin/")
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

// TreeHash returns the content hash of tracked-plus-untracked-unignored files under
// root, computed through a THROWAWAY index so the real index is never touched — this
// is the gate verdict cache key. It returns the literal "none" on any failure or an
// empty result. The temp index lives outside the repo so it can't join the tree it
// hashes; `git add -A` respects .gitignore, which is the intended scope.
func TreeHash(root string) string {
	dir, err := os.MkdirTemp("", "bench-tree")
	if err != nil {
		return "none"
	}
	defer os.RemoveAll(dir)
	idx := filepath.Join(dir, "index")

	// Seed the throwaway index from HEAD, falling back to an empty tree in a repo
	// with no commits yet, then stage everything on disk and write the tree.
	if !idxOK(root, idx, "read-tree", "HEAD") {
		if !idxOK(root, idx, "read-tree", "--empty") {
			return "none"
		}
	}
	if !idxOK(root, idx, "add", "-A") {
		return "none"
	}
	hash, err := idxOutput(root, idx, "write-tree")
	if err != nil || hash == "" {
		return "none"
	}
	return hash
}

// ChangedPathsBetweenTrees reports the root-relative paths whose content differs between
// two tree objects. It shells out to `git diff --name-only <from> <to>` so the compared
// trees stay the same source of truth `bench status` already uses. Any invalid tree,
// missing object, or diff failure returns ok=false so callers can fail closed.
func ChangedPathsBetweenTrees(root, fromTree, toTree string) ([]string, bool) {
	if fromTree == "" || toTree == "" || fromTree == "none" || toTree == "none" {
		return nil, false
	}
	out, err := Output("-C", root, "diff", "--name-only", fromTree, toTree)
	if err != nil {
		return nil, false
	}
	if out == "" {
		return []string{}, true
	}
	return strings.Split(out, "\n"), true
}

// idxCommand builds a `git -C root <args>` command whose index is the throwaway idx
// file rather than the repository's own — the shared invocation form for TreeHash.
func idxCommand(root, idx string, args ...string) *exec.Cmd {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
	return cmd
}

// idxOK reports whether the throwaway-index git command exited zero.
func idxOK(root, idx string, args ...string) bool {
	return idxCommand(root, idx, args...).Run() == nil
}

// idxOutput runs the throwaway-index git command and returns stdout with a single
// trailing newline trimmed.
func idxOutput(root, idx string, args ...string) (string, error) {
	var out bytes.Buffer
	cmd := idxCommand(root, idx, args...)
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimRight(out.String(), "\n"), err
}
