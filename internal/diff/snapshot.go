// Git-shell parsing, snapshot identity, and drift detection for package diff.
package diff

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// snapshotAfterRead is a test seam for a deterministic mutation between the two
// identity captures. Production leaves it as a no-op.
var snapshotAfterRead = func() {}

// SetSnapshotAfterReadForTest installs the deterministic movement seam and returns
// its restoration function. Production callers leave the seam as its no-op default.
func SetSnapshotAfterReadForTest(after func()) func() {
	previous := snapshotAfterRead
	if after == nil {
		after = func() {}
	}
	snapshotAfterRead = after
	return func() { snapshotAfterRead = previous }
}

// parseNameStatusZ turns `git diff --name-status --no-renames -z` output into
// status/path rows. The NUL framing carries each path raw (never git C-quoted), so a
// path with a space, a non-ASCII byte, or a quote arrives intact and is TOON-escaped a
// single layer downstream.
func parseNameStatusZ(raw []byte) [][]string {
	parts := strings.Split(string(raw), "\x00")
	var rows [][]string
	for i := 0; i+1 < len(parts); i += 2 {
		st, p := parts[i], parts[i+1]
		if st == "" && p == "" {
			continue
		}
		rows = append(rows, []string{st, p})
	}
	return rows
}

// changedFiles renders the files table for a `git diff` range. rangeArgs is passed
// straight through to `git diff --name-status --no-renames -z`: the resolved base alone
// for the branch-relative path, so Git includes index and tracked worktree changes, or
// two bare refs ("base", "head") for the commit-relative path. `git diff` treats two
// positional refs the same as an explicit two-dot range, which is the exact two-commit
// diff `--commit` needs.
func changedFiles(rangeArgs ...string) ([][]string, error) {
	args := append([]string{"diff", "--name-status", "--no-renames", "-z"}, rangeArgs...)
	raw, err := git.Raw(args...)
	if err != nil {
		return nil, err
	}
	return parseNameStatusZ(raw), nil
}

// parseLogFormat turns `git log --format=%h%x00%s` output into sha/subject rows. Each
// commit is one line, because a subject is by definition the first line of the commit
// message and carries no embedded newline. The NUL between sha and subject is a
// delimiter git itself never puts in either field, so a comma or quote in the subject
// arrives raw for the caller to TOON-escape a single layer downstream, the same
// NUL-framing discipline parseNameStatusZ uses for paths.
func parseLogFormat(raw []byte) [][]string {
	s := strings.TrimRight(string(raw), "\n")
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	rows := make([][]string, 0, len(lines))
	for _, line := range lines {
		sha, subject, ok := strings.Cut(line, "\x00")
		if !ok {
			continue
		}
		rows = append(rows, []string{sha, subject})
	}
	return rows
}

// commitLog renders the log table for a `git log` range expression: always a literal
// two-dot range string ("base..HEAD" or "base..head"). `git log`'s two-dot form, unlike
// `git diff`'s, has a distinct meaning from two bare refs.
func commitLog(rangeExpr string) ([][]string, error) {
	raw, err := git.Raw("log", "--format=%h%x00%s", rangeExpr)
	if err != nil {
		return nil, err
	}
	return parseLogFormat(raw), nil
}

// diffBody renders the raw `git diff` body. Same rangeArgs contract as changedFiles:
// the resolved base for the branch-relative path, or two bare refs for
// the commit-relative path.
func diffBody(rangeArgs ...string) ([]byte, error) {
	args := append([]string{"diff"}, rangeArgs...)
	return git.Raw(args...)
}

func changedFilesAt(root string, rangeArgs ...string) ([][]string, error) {
	args := append([]string{"-C", root, "diff", "--name-status", "--no-renames", "-z"}, rangeArgs...)
	raw, err := git.Raw(args...)
	if err != nil {
		return nil, err
	}
	return parseNameStatusZ(raw), nil
}

func diffBodyAt(root string, rangeArgs ...string) ([]byte, error) {
	args := append([]string{"-C", root, "diff"}, rangeArgs...)
	return git.Raw(args...)
}

func commitLogAt(root, rangeExpr string) ([][]string, error) {
	raw, err := git.Raw("-C", root, "log", "--format=%h%x00%s", rangeExpr)
	if err != nil {
		return nil, err
	}
	return parseLogFormat(raw), nil
}

func fileKind(root, path string) string {
	info, err := os.Lstat(filepath.Join(root, path))
	if err != nil {
		return "dangling-symlink"
	}
	mode := info.Mode()
	if mode.IsRegular() {
		return ""
	}
	if mode&os.ModeSymlink != 0 {
		if _, err := os.Stat(filepath.Join(root, path)); err != nil {
			return "dangling-symlink"
		}
		return "symlink"
	}
	if mode&os.ModeNamedPipe != 0 {
		return "fifo"
	}
	if mode&os.ModeSocket != 0 {
		return "socket"
	}
	if mode&os.ModeDevice != 0 {
		return "device"
	}
	return "special"
}

func checkoutRows(changes []git.PorcelainEntry) [][]string {
	rows := make([][]string, 0, len(changes))
	for _, change := range changes {
		if len(change.Status) != 2 {
			continue
		}
		x, y := string(change.Status[0]), string(change.Status[1])
		if x == " " {
			x = "-"
		}
		if y == " " {
			y = "-"
		}
		rows = append(rows, []string{x, y, change.Path})
	}
	return rows
}

func statusCounts(changes []git.PorcelainEntry) (staged, unstaged, untracked int) {
	for _, change := range changes {
		if change.Status == "??" {
			untracked++
			continue
		}
		if len(change.Status) == 2 && change.Status[0] != ' ' {
			staged++
		}
		if len(change.Status) == 2 && change.Status[1] != ' ' {
			unstaged++
		}
	}
	return
}

func shortstat(root string, args ...string) (insertions, deletions int, err error) {
	cmd := append([]string{"-C", root, "diff", "--shortstat"}, args...)
	out, err := git.Output(cmd...)
	if err != nil || out == "" {
		return 0, 0, err
	}
	for _, field := range strings.Split(out, ",") {
		words := strings.Fields(field)
		if len(words) < 2 {
			continue
		}
		n, parseErr := strconv.Atoi(words[0])
		if parseErr != nil {
			continue
		}
		if strings.HasPrefix(words[1], "insertion") {
			insertions = n
		}
		if strings.HasPrefix(words[1], "deletion") {
			deletions = n
		}
	}
	return insertions, deletions, nil
}

func whitespace(root string, args ...string) (clean bool, offenses int, err error) {
	cmd := append([]string{"-C", root, "diff", "--check"}, args...)
	out, runErr := git.Raw(cmd...)
	if runErr == nil {
		return true, 0, nil
	}
	if len(out) == 0 {
		return false, 0, runErr
	}
	return false, len(strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")), nil
}

type snapshotIdentity struct {
	head, defaultTip, recordedBase, index, porcelain string
	dirtyContent, dirtyMode, dirtyGitlink            string
}

func digest(bytes []byte) string { return fmt.Sprintf("%x", sha256.Sum256(bytes)) }

func capturedIdentity(root string, facts git.DiffFacts) (snapshotIdentity, error) {
	identity := snapshotIdentity{
		head:         facts.Head,
		defaultTip:   facts.DefaultTip,
		recordedBase: facts.RecordedBase,
		porcelain:    digest(facts.Porcelain),
	}
	var err error
	identity.index, err = indexIdentity(root)
	if err != nil {
		return snapshotIdentity{}, err
	}
	identity.dirtyContent, identity.dirtyMode, identity.dirtyGitlink, err = pathIdentities(root, facts.Changes)
	return identity, err
}

func recapturedIdentity(root string, facts git.DiffFacts) (snapshotIdentity, error) {
	after := facts
	var err error
	after.Head, err = git.Output("-C", root, "rev-parse", "HEAD")
	if err != nil {
		return snapshotIdentity{}, err
	}
	if facts.DefaultResolved {
		after.DefaultTip, err = git.Output("-C", root, "rev-parse", facts.DefaultBranch)
		if err != nil {
			return snapshotIdentity{}, err
		}
	}
	if facts.Branch != "HEAD" {
		after.RecordedBase, _ = git.Output("-C", root, "config", "branch."+facts.Branch+".benchBase")
	}
	after.Porcelain, after.Changes, err = git.AllFilesStatus(root)
	if err != nil {
		return snapshotIdentity{}, err
	}
	return capturedIdentity(root, after)
}

func indexIdentity(root string) (string, error) {
	if indexPath, err := git.Output("-C", root, "rev-parse", "--git-path", "index"); err == nil {
		if !filepath.IsAbs(indexPath) {
			indexPath = filepath.Join(root, indexPath)
		}
		data, readErr := os.ReadFile(indexPath)
		if readErr != nil {
			return "", readErr
		}
		return digest(data), nil
	} else {
		return "", err
	}
}

func pathIdentities(root string, changes []git.PorcelainEntry) (contentDigest, modeDigest, gitlinkDigest string, err error) {
	contents := make([]string, 0, len(changes))
	modes := make([]string, 0, len(changes))
	gitlinks := make([]string, 0, len(changes))
	for _, entry := range changes {
		path := filepath.Join(root, entry.Path)
		info, err := os.Lstat(path)
		if err != nil {
			contents = append(contents, entry.Path+":missing")
			modes = append(modes, entry.Path+":missing")
			continue
		}
		modes = append(modes, entry.Path+":"+info.Mode().String())
		if head, ok := gitlinkHead(root, entry.Path); ok {
			gitlinks = append(gitlinks, entry.Path+":"+head)
			continue
		}
		value := ""
		switch {
		case info.Mode().IsRegular():
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return "", "", "", readErr
			}
			value = digest(data)
		case info.Mode()&os.ModeSymlink != 0:
			target, readErr := os.Readlink(path)
			if readErr != nil {
				return "", "", "", readErr
			}
			value = target
		default:
			value = fmt.Sprintf("%d:%d", info.Size(), info.ModTime().UnixNano())
		}
		contents = append(contents, entry.Path+":"+value)
	}
	sort.Strings(contents)
	sort.Strings(modes)
	sort.Strings(gitlinks)
	return digest([]byte(strings.Join(contents, "\x00"))), digest([]byte(strings.Join(modes, "\x00"))), digest([]byte(strings.Join(gitlinks, "\x00"))), nil
}

func gitlinkHead(root, path string) (string, bool) {
	entry, err := git.Output("-C", root, "ls-files", "--stage", "--", path)
	if err != nil || !strings.HasPrefix(entry, "160000 ") {
		return "", false
	}
	head, err := git.Output("-C", filepath.Join(root, path), "rev-parse", "HEAD")
	if err != nil {
		return "missing", true
	}
	return head, true
}

func (i snapshotIdentity) drifted(after snapshotIdentity) string {
	for _, field := range []struct{ name, before, after string }{
		{"HEAD", i.head, after.head}, {"default tip", i.defaultTip, after.defaultTip},
		{"recorded base", i.recordedBase, after.recordedBase}, {"index", i.index, after.index},
		{"porcelain", i.porcelain, after.porcelain},
		{"dirty content", i.dirtyContent, after.dirtyContent},
		{"dirty mode", i.dirtyMode, after.dirtyMode},
		{"dirty gitlink HEAD", i.dirtyGitlink, after.dirtyGitlink},
	} {
		if field.before != field.after {
			return field.name
		}
	}
	return ""
}

func untrackedPatches(root string, facts git.DiffFacts) ([]byte, error) {
	var paths []string
	for _, entry := range facts.Changes {
		if entry.Status == "??" && fileKind(root, entry.Path) == "" {
			paths = append(paths, entry.Path)
		}
	}
	sort.Strings(paths)
	var body []byte
	for _, path := range paths {
		patch, err := git.Raw("-C", root, "diff", "--no-index", "--", "/dev/null", path)
		if err != nil && len(patch) == 0 {
			return nil, err
		}
		body = append(body, patch...)
	}
	return body, nil
}

// ChangedFilePaths is the single source of the changed-path set for bench diff and
// its consumers: the path half of changedFiles' rows for a resolved review base,
// carrying bench diff's exact committed+index+tracked-worktree semantics through
// `git diff --name-status --no-renames -z <base>`, so a consumer never re-derives that
// git invocation or the NUL-pair parse.
func ChangedFilePaths(base string) ([]string, error) {
	return ChangedFilePathsAt("", base)
}

// ChangedFilePathsAt is ChangedFilePaths rooted at root when supplied.
func ChangedFilePathsAt(root, base string) ([]string, error) {
	rows, err := changedFilesAt(root, base)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(rows))
	for _, row := range rows {
		paths = append(paths, row[1])
	}
	return paths, nil
}
