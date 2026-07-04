// Package git is the one source of git subprocess invocation for the AXI query
// commands. Every ported parser shells out to git through here — git stays the
// source of repository truth (root, config, merge-base, diff), exactly as the shell
// commands did, so there is one place the invocation form lives.
package git

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Root returns the working tree's top-level directory, or an error when the cwd is
// not inside a git repository (the `not in a git repository` posture of every command).
func Root() (string, error) {
	return Output("rev-parse", "--show-toplevel")
}

// DefaultBranch is the repository's default branch: origin/HEAD's short name with the
// `origin/` prefix stripped, falling back to "main" when the ref is unset (no remote
// HEAD) or empty. The one source both `diff` and `status` read — and the Go mirror of
// bench.sh's default_branch — so the three surfaces agree on what "default branch" is.
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
