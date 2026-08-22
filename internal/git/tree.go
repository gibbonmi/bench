package git

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TreeHash returns the content hash of tracked and untracked, unignored files under root.
// It computes the hash through a THROWAWAY index, so the real index stays untouched; the
// gate uses this hash as its verdict cache key. It returns the literal "none" on any
// failure or an empty result. The temp index lives outside the repo, so it cannot join
// the tree it hashes. `git add -A` respects .gitignore, which is the intended scope.
func TreeHash(root string) string {
	dir, err := os.MkdirTemp("", "bench-tree")
	if err != nil {
		return "none"
	}
	defer os.RemoveAll(dir)
	idx := filepath.Join(dir, "index")

	// Seed the throwaway index from HEAD. In a repo with no commits yet, fall back to an
	// empty tree. Then stage everything on disk and write the tree.
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

// reverseAppliesToDefault reports whether branch's cumulative diff against its merge base
// with def applies in reverse to def's tree. A true result proves the branch's content is
// already present in def, however it landed, which is what survives a squash. The apply
// runs against a throwaway index seeded from def, the TreeHash idiom, so it touches
// neither the working tree nor the real index.
//
// A true verdict authorizes branch deletion, so every step refuses rather than guesses.
// A missing merge base, a diff that fails to generate, or an apply that fails for any
// reason all refuse the branch. The function refuses a submodule pointer outright. A
// patch carries only the subproject's sha, not its content, so a clean apply would not
// prove the work is present. The apply itself stays byte- and mode-exact — full-index
// binary patches, no rename detection, and whitespace leniency off. Loosening any of
// these trades an orphaned branch for silently destroyed work, the wrong direction.
func reverseAppliesToDefault(root, branch, def string) bool {
	base, err := Output("-C", root, "merge-base", def, branch)
	if err != nil || base == "" {
		return false
	}
	changes, err := Output("-C", root, "diff", "--raw", "--no-renames", base, branch)
	if err != nil {
		return false
	}
	if changes == "" {
		// The branch tree equals its merge base's, an ancestor of def: nothing to prove.
		return true
	}
	for _, line := range strings.Split(changes, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0] == ":160000" || fields[1] == "160000" {
			return false
		}
		// The function refuses a mode change on a surviving entry — a chmod, or a
		// file/symlink typechange. git apply treats a preimage mode mismatch as a warning,
		// not a failure, so a clean apply would not prove the mode landed. Adds and deletes
		// keep one side at 000000, and the apply itself verifies their modes.
		if fields[0] != ":000000" && fields[1] != "000000" && fields[0][1:] != fields[1] {
			return false
		}
	}
	patch, err := Raw("-C", root, "diff", "--binary", "--no-renames", "--full-index", base, branch)
	if err != nil || len(patch) == 0 {
		return false
	}
	dir, err := os.MkdirTemp("", "bench-landed")
	if err != nil {
		return false
	}
	defer os.RemoveAll(dir)
	idx := filepath.Join(dir, "index")
	if !idxOK(root, idx, "read-tree", def) {
		return false
	}
	check := idxCommand(root, idx, "apply", "--cached", "--check", "--reverse", "--no-ignore-whitespace", "--whitespace=nowarn")
	check.Stdin = bytes.NewReader(patch)
	return check.Run() == nil
}

// idxCommand builds a `git -C root <args>` command whose index is the throwaway idx
// file rather than the repository's own. This is the shared invocation form for
// TreeHash.
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
