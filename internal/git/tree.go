package git

import (
	"bytes"
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
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
		if len(fields) < 2 || fields[0] == ":"+GitlinkMode || fields[1] == GitlinkMode {
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
	check := IndexCommand(root, idx, "apply", "--cached", "--check", "--reverse", "--no-ignore-whitespace", "--whitespace=nowarn")
	check.Stdin = bytes.NewReader(patch)
	return check.Run() == nil
}

// IndexCommand builds a `git -C root <args>` command whose index is the throwaway idx
// file rather than the repository's own. It is the one owner of the index-file
// invocation form; the caller chooses how to run the command and shape its result.
func IndexCommand(root, idx string, args ...string) *exec.Cmd {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
	return cmd
}

// idxOK reports whether the throwaway-index git command exited zero.
func idxOK(root, idx string, args ...string) bool {
	return IndexCommand(root, idx, args...).Run() == nil
}

// idxOutput runs the throwaway-index git command and returns stdout with a single
// trailing newline trimmed.
func idxOutput(root, idx string, args ...string) (string, error) {
	var out bytes.Buffer
	cmd := IndexCommand(root, idx, args...)
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimRight(out.String(), "\n"), err
}

// TreeEntry is one record of `git ls-tree`: the entry's six-digit mode, its object type,
// its object ID, and its repository-relative path.
type TreeEntry struct {
	Mode, Type, Object, Path string
}

// IsRegularFile reports whether the entry holds regular-file content.
func (e TreeEntry) IsRegularFile() bool { return regularFileModes[e.Mode] }

// ListTree runs `git -C root ls-tree -z <args>` and reads its records. Each record is
// `<mode> <type> <object>\t<path>` and ends at a NUL, so the path is every byte past the
// first tab and Git quotes none of it.
func ListTree(root string, args ...string) ([]TreeEntry, error) {
	raw, err := Raw(listTreeArgs(root, args)...)
	if err != nil {
		return nil, err
	}
	return parseTreeEntries(raw)
}

// listTreeArgs builds the git argv of one ls-tree listing, which is also its run key.
func listTreeArgs(root string, args []string) []string {
	return append([]string{"-C", root, "ls-tree", "-z"}, args...)
}

// parseTreeEntries reads the NUL-terminated records of one ls-tree listing.
func parseTreeEntries(raw []byte) ([]TreeEntry, error) {
	records := strings.Split(string(raw), "\x00")
	if records[len(records)-1] != "" {
		return nil, fmt.Errorf("git: unterminated tree entry %q", records[len(records)-1])
	}
	entries := make([]TreeEntry, 0, len(records)-1)
	for _, record := range records[:len(records)-1] {
		meta, path, ok := strings.Cut(record, "\t")
		fields := strings.Fields(meta)
		if !ok || len(fields) != 3 {
			return nil, fmt.Errorf("git: unreadable tree entry %q", record)
		}
		entries = append(entries, TreeEntry{Mode: fields[0], Type: fields[1], Object: fields[2], Path: path})
	}
	return entries, nil
}

// The tree-entry modes of a symbolic link and of a gitlink.
const (
	SymlinkMode = "120000"
	GitlinkMode = "160000"
)

// TreeSymlinks answers the tracked symbolic links of tree: the repository-relative path of
// each link, mapped to the link text that its blob holds.
func TreeSymlinks(root, tree string) (map[string]string, error) {
	entries, err := ListTree(root, "-r", tree)
	if err != nil {
		return nil, err
	}
	links := map[string]string{}
	for _, entry := range entries {
		if entry.Mode != SymlinkMode {
			continue
		}
		target, err := ReadControlBlob(root, entry.Object)
		if err != nil {
			return nil, err
		}
		links[entry.Path] = string(target)
	}
	return links, nil
}

// TreeWithoutFile returns a tree with one exact file removed through a private index.
func TreeWithoutFile(root, tree, path string) (string, error) {
	listing, err := ListTree(root, tree, "--", path)
	if err != nil {
		return "", err
	}
	if len(listing) > 0 && listing[0].Mode == "040000" {
		return "", fmt.Errorf("tree exclusion names a directory: %s", path)
	}
	dir, err := os.MkdirTemp("", "bench-tree-exclusion-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	idx := filepath.Join(dir, "index")
	if !idxOK(root, idx, "read-tree", tree) || !idxOK(root, idx, "update-index", "--force-remove", "--", path) {
		return "", fmt.Errorf("cannot exclude exact file %s", path)
	}
	return idxOutput(root, idx, "write-tree")
}

// ReadTreeFile reads a regular immutable blob under the control-record bound. It reads
// through a nil run, so each call starts its listing and its blob read.
func ReadTreeFile(root, tree, path string) ([]byte, error) {
	return (*Run)(nil).ReadTreeFile(root, tree, path)
}

// ReadControlBlob bounds an immutable blob read for control-file consumers.
func ReadControlBlob(root, object string) ([]byte, error) {
	cmd := exec.Command("git", controlBlobArgs(root, object)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	read := bounds.Read(stdout, bounds.ControlRecordLimit)
	_ = stdout.Close()
	waitErr := cmd.Wait()
	if read.Err != nil {
		return nil, read.Err
	}
	if waitErr != nil {
		return nil, waitErr
	}
	return read.Data, nil
}

// controlBlobArgs builds the git argv of one blob read, which is also its run key.
func controlBlobArgs(root, object string) []string {
	return []string{"-C", root, "cat-file", "blob", object}
}

// treeChangeFields is the field count of one raw-diff entry's metadata: the two
// modes, the two object IDs, and the status letter.
const treeChangeFields = 5

// TreeChange is one entry of the raw diff between the base tree and the composed
// tree. Status is Git's own status letter. SrcMode and DstMode are the six-digit modes
// of the two sides, and one of them is `000000` when that side holds nothing. Path
// carries the file's own bytes, so a name with a space or a byte above ASCII survives.
type TreeChange struct {
	Status  string
	SrcMode string
	DstMode string
	Path    string
}

// TreeChanges lists raw changes between two Git trees in path order.
// Rename detection is off, so each path has one metadata frame and one path frame.
// NUL framing preserves every path byte without Git's quoted-name encoding.
func TreeChanges(root, from, to string) ([]TreeChange, error) {
	return treeChanges(root, from, to, false)
}

// TreeChangesIncludingSubmodules keeps changed gitlinks visible despite Git configuration.
func TreeChangesIncludingSubmodules(root, from, to string) ([]TreeChange, error) {
	return treeChanges(root, from, to, true)
}

func treeChanges(root, from, to string, includeSubmodules bool) ([]TreeChange, error) {
	args := []string{"-C", root, "diff", "--raw", "--no-renames", "-z"}
	if includeSubmodules {
		args = append(args, "--ignore-submodules=none")
	}
	raw, err := Raw(append(args, from, to)...)
	if err != nil {
		return nil, fmt.Errorf("git: tree change list unavailable: %w", err)
	}
	frames := strings.Split(string(raw), "\x00")
	var changes []TreeChange
	for i := 0; i+1 < len(frames); i += 2 {
		change, err := parseTreeChange(frames[i], frames[i+1])
		if err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}
	return changes, nil
}

// parseTreeChange reads one `--raw -z` entry, whose metadata frame is
// `:<srcmode> <dstmode> <srcsha> <dstsha> <status>` and whose path is the frame after it.
func parseTreeChange(meta, path string) (TreeChange, error) {
	fields := strings.Fields(strings.TrimPrefix(meta, ":"))
	if !strings.HasPrefix(meta, ":") || len(fields) != treeChangeFields || path == "" {
		return TreeChange{}, fmt.Errorf("git: unreadable tree change entry %q", meta)
	}
	return TreeChange{Status: fields[4], SrcMode: fields[0], DstMode: fields[1], Path: path}, nil
}
