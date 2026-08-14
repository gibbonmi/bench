// Package diff owns the coherent Git review snapshot. A live response resolves its
// branch range once, renders every bounded and complete section from that attempt,
// and compares each patch-observable identity before emitting any bytes. A named
// commit remains an immutable first-parent view and omits unrelated checkout facts.
package diff

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
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

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, repeated flags, and help all come from there rather
// than a local switch.
// Help is fullHelp without its trailing newline, because the caller appends one.
var grammar = usage.Grammar{
	Cmd:   "bench diff",
	Help:  strings.TrimSuffix(fullHelp, "\n"),
	Flags: []usage.Flag{{Name: "--full"}, {Name: "--commit", HasValue: true}, {Name: "--base", HasValue: true, NoEmptyValue: true}},
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
// straight through to `git diff --name-status --no-renames -z`: the resolved base
// alone for the branch-relative path, so Git includes index and tracked worktree
// changes, or two bare refs ("base", "head") for the commit-relative path — `git diff` treats
// two positional refs the same as an explicit two-dot range, which is the exact
// two-commit diff `--commit` needs.
func changedFiles(rangeArgs ...string) ([][]string, error) {
	args := append([]string{"diff", "--name-status", "--no-renames", "-z"}, rangeArgs...)
	raw, err := git.Raw(args...)
	if err != nil {
		return nil, err
	}
	return parseNameStatusZ(raw), nil
}

// parseLogFormat turns `git log --format=%h%x00%s` output into sha/subject rows. Each
// commit is one line (a subject is, by definition, the first line of the commit
// message and carries no embedded newline); the NUL between sha and subject is a
// delimiter git itself never puts in either field, so a comma or quote in the
// subject arrives raw for the caller to TOON-escape a single layer downstream — the
// same NUL-framing discipline parseNameStatusZ uses for paths.
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

// commitLog renders the log table for a `git log` range expression — always a
// literal two-dot range string ("base..HEAD" or "base..head"), since `git log`'s
// two-dot form (unlike `git diff`'s) has a distinct meaning from two bare refs.
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

const fullHelp = `usage: bench diff [--full] [--commit <sha>] [--base <commit>]
  Live mode reports one movement-checked revision, aggregate, files, checkout,
  and whitespace snapshot. --full also appends the commit log and verbatim
  tracked patch, then path-sorted raw Git patches for untracked regular files.
  --commit <sha> reports the immutable first-parent view of one resolved commit
  and omits unrelated live-checkout facts. Bounded results advertise the exact
  --full invocation; complete and clean results advertise an empty help table.
`

// diffRange is the resolved review range every rendering section shares: the base
// and method preamble lines, and the argument shapes changedFiles/commitLog/diffBody
// need for either the branch-relative (resolved base through worktree) range or the
// commit-relative (exact first-parent) range.
type diffRange struct {
	base      string
	head      string
	method    string
	filesArgs []string
	logRange  string
	bodyArgs  []string
}

// SourceRange is the frozen, inclusive-ancestor source identity shared by explicit
// diff and preflight. Tip is captured with Base so callers never pair a path set
// from one source revision with the identity of another.
type SourceRange struct {
	Base, Tip      string
	CommittedPaths []string
}

// MovementSnapshot is the typed checkout state one movement-checked read exposes to
// its caller. The root and facts stay coupled so derived source ranges and paths
// cannot accidentally inspect the process working directory.
type MovementSnapshot struct {
	root  string
	Facts git.DiffFacts
}

// MovementResult is the single outcome of a root-bound attempt. Drift remains
// distinct from a read failure so callers retry only repository movement.
type MovementResult struct {
	DriftKind, DriftHint string
	Kind, Hint           string
}

// MovementChecked captures one root-bound checkout identity, runs read against its
// typed facts, then verifies the identity again. A non-empty drift asks the caller to
// retry; kind and hint preserve an ordinary read failure without reclassifying it.
func MovementChecked(root string, read func(MovementSnapshot) (kind, hint string)) MovementResult {
	facts, err := git.AllFilesFacts(root)
	if err != nil {
		return MovementResult{Kind: "checkout facts failed", Hint: err.Error()}
	}
	before, err := capturedIdentity(root, facts)
	if err != nil {
		return MovementResult{Kind: "snapshot identity failed", Hint: err.Error()}
	}
	if kind, hint := read(MovementSnapshot{root: root, Facts: facts}); kind != "" {
		return MovementResult{Kind: kind, Hint: hint}
	}
	snapshotAfterRead()
	after, err := recapturedIdentity(root, facts)
	if err != nil {
		return MovementResult{Kind: "snapshot identity failed", Hint: err.Error()}
	}
	if drift := before.drifted(after); drift != "" {
		return MovementResult{DriftKind: drift, DriftHint: "the " + drift + " changed while reading; retry the exact invocation"}
	}
	return MovementResult{}
}

// MovementCheckedRetry is the movement-retry policy every drift-sensitive read shares:
// one retry when the only thing that failed was repository movement, then the second
// attempt's answer stands — including its terminal drift hint, which no caller spells
// for itself. A read failure is never retried, so a caller cannot turn a broken read
// into a drift refusal.
func MovementCheckedRetry(root string, read func(MovementSnapshot) (kind, hint string)) MovementResult {
	var result MovementResult
	for attempt := 0; attempt < 2; attempt++ {
		result = MovementChecked(root, read)
		if result.Kind != "" || result.DriftKind == "" {
			return result
		}
	}
	return result
}

// ResolveSourceRange resolves base against the snapshot's captured source tip.
func (snapshot MovementSnapshot) ResolveSourceRange(base string) (SourceRange, string, string) {
	return ResolveSourceRange(snapshot.root, base, snapshot.Facts.Head)
}

// SourceSnapshotPaths returns source paths derived from the same captured checkout.
func (snapshot MovementSnapshot) SourceSnapshotPaths(source SourceRange) ([]string, error) {
	return sourceSnapshotPaths(snapshot.root, source, snapshot.Facts)
}

// ResolveSourceRange resolves an explicit source base and its exact source tip.
func ResolveSourceRange(root, base, tip string) (SourceRange, string, string) {
	b, err := git.Output("-C", root, "rev-parse", "--verify", base+"^{commit}")
	if err != nil {
		return SourceRange{}, "cannot resolve --base", "'" + base + "' does not name a commit reachable in this repository"
	}
	if !git.OK("-C", root, "merge-base", "--is-ancestor", b, tip) {
		return SourceRange{}, "--base is not an ancestor", "'" + b + "' is not an ancestor of '" + tip + "'"
	}
	paths, err := changedFilesAt(root, b, tip)
	if err != nil {
		return SourceRange{}, "committed source paths not readable", err.Error()
	}
	committed := make([]string, 0, len(paths))
	for _, path := range paths {
		committed = append(committed, path[1])
	}
	return SourceRange{Base: b, Tip: tip, CommittedPaths: committed}, "", ""
}

// resolveCommitRange builds the diffRange for `--commit <sha>`: base is <sha>'s
// resolved first parent. The sha is verified before anything renders — an
// unresolvable sha and a root commit's missing parent are each their own
// structured error (kind, hint), never a leaked git failure.
func resolveCommitRange(root, commitArg string) (dr diffRange, errKind, errHint string) {
	headSha, err := git.Output("-C", root, "rev-parse", "--verify", commitArg+"^{commit}")
	if err != nil {
		return diffRange{}, "cannot resolve --commit",
			"'" + commitArg + "' does not name a commit reachable in this repository"
	}
	baseSha, err := git.Output("-C", root, "rev-parse", "--verify", commitArg+"^")
	if err != nil {
		return diffRange{}, "--commit has no parent",
			"'" + commitArg + "' is a root commit — there is no first parent to diff against"
	}
	return diffRange{
		base:      baseSha,
		head:      headSha,
		method:    "commit " + headSha,
		filesArgs: []string{baseSha, headSha},
		logRange:  baseSha + ".." + headSha,
		bodyArgs:  []string{baseSha, headSha},
	}, "", ""
}

// resolveBranchRange builds the diffRange for bare `bench diff`/`--full`: the
// recorded-key base when it names a reachable ancestor, else merge-base with the
// default branch — byte-identical to the pre-`--commit` behavior.
func resolveBranchRange(root string) (dr diffRange, errKind, errHint string) {
	base, method, errKind, errHint := ResolveReviewBase(root)
	if errKind != "" {
		return diffRange{}, errKind, errHint
	}
	return diffRange{
		base:      base,
		method:    method,
		filesArgs: []string{base},
		logRange:  base + "..HEAD",
		bodyArgs:  []string{base},
	}, "", ""
}

// Command implements `bench diff`.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, full := parsed.Flags["--full"]
	commitArg, hasCommit := parsed.Flags["--commit"]
	baseArg, hasBase := parsed.Flags["--base"]
	if hasCommit && hasBase {
		return toon.Usage(grammar.Cmd, "--base") + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	invocation := append([]string{"diff"}, args...)
	out, drift, driftHint, errKind, errHint := renderAttempt(root, hasCommit, commitArg, hasBase, baseArg, full)
	if errKind != "" {
		return toon.Errorf(errKind, errHint) + "\n", 1
	}
	if drift == "" {
		return out, 0
	}
	help, err := axi.RenderHelp([]axi.Action{axi.RetryDiff(invocation)})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return toon.Errorf("snapshot drift", driftHint) + "\n" + help, 1
}

func renderAttempt(root string, hasCommit bool, commitArg string, hasBase bool, baseArg string, full bool) (out, drift, driftHint, errKind, errHint string) {
	var dr diffRange
	if hasCommit {
		dr, errKind, errHint = resolveCommitRange(root, commitArg)
		if errKind != "" {
			return "", "", "", errKind, errHint
		}
		out, errKind, errHint := renderCommit(root, dr, full)
		return out, "", "", errKind, errHint
	}
	result := MovementCheckedRetry(root, func(snapshot MovementSnapshot) (kind, hint string) {
		if hasBase {
			dr, kind, hint = resolveExplicitRange(root, baseArg, snapshot.Facts.Head)
		} else {
			dr, kind, hint = resolveBranchRangeFromFacts(root, snapshot.Facts)
		}
		if kind != "" {
			return kind, hint
		}
		out, kind, hint = renderLive(root, dr, snapshot.Facts, full)
		return kind, hint
	})
	return out, result.DriftKind, result.DriftHint, result.Kind, result.Hint
}

func resolveExplicitRange(root, base, head string) (diffRange, string, string) {
	source, kind, hint := ResolveSourceRange(root, base, head)
	if kind != "" {
		return diffRange{}, kind, hint
	}
	return diffRange{base: source.Base, head: source.Tip, method: "explicit", filesArgs: []string{source.Base}, logRange: source.Base + ".." + source.Tip, bodyArgs: []string{source.Base}}, "", ""
}

// SourceSnapshotPaths returns the complete explicit-source inventory: committed,
// index, tracked-worktree, and untracked paths. The committed half belongs to SourceRange.
func SourceSnapshotPaths(root string, source SourceRange) ([]string, error) {
	facts, err := git.AllFilesFacts(root)
	if err != nil {
		return nil, err
	}
	return sourceSnapshotPaths(root, source, facts)
}

func sourceSnapshotPaths(root string, source SourceRange, facts git.DiffFacts) ([]string, error) {
	paths := append([]string(nil), source.CommittedPaths...)
	tracked, err := ChangedFilePathsAt(root, source.Base)
	if err != nil {
		return nil, err
	}
	paths = append(paths, tracked...)
	for _, entry := range facts.Changes {
		if entry.Status == "??" {
			paths = append(paths, entry.Path)
		}
	}
	sort.Strings(paths)
	unique := paths[:0]
	for _, path := range paths {
		if len(unique) == 0 || unique[len(unique)-1] != path {
			unique = append(unique, path)
		}
	}
	return unique, nil
}

func renderLive(root string, dr diffRange, facts git.DiffFacts, full bool) (string, string, string) {
	tracked, err := changedFilesAt(root, dr.filesArgs...)
	if err != nil {
		return "", "git diff --name-status failed", err.Error()
	}
	files := make([][]string, 0, len(tracked)+len(facts.Changes))
	for _, row := range tracked {
		files = append(files, []string{row[0], row[1], ""})
	}
	for _, entry := range facts.Changes {
		if entry.Status == "??" {
			files = append(files, []string{"?", entry.Path, fileKind(root, entry.Path)})
		}
	}
	commits, err := commitLogAt(root, dr.logRange)
	if err != nil {
		return "", "git log failed", err.Error()
	}
	insertions, deletions, err := shortstat(root, dr.filesArgs...)
	if err != nil {
		return "", "git diff --shortstat failed", err.Error()
	}
	clean, offenses, err := whitespace(root, dr.filesArgs...)
	if err != nil {
		return "", "git diff --check failed", err.Error()
	}
	staged, unstaged, untracked := statusCounts(facts.Changes)
	branch := facts.Branch
	if branch == "HEAD" {
		branch = "(detached)"
	}
	revision, err := toon.Table("revision", []string{"branch", "default", "ahead", "behind", "base", "method", "head"}, [][]string{{branch, facts.DefaultBranch, strconv.Itoa(facts.Ahead), strconv.Itoa(facts.Behind), dr.base, dr.method, facts.Head}})
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	aggregate, err := toon.Table("aggregate", []string{"commits", "files", "insertions", "deletions", "staged", "unstaged", "untracked"}, [][]string{{strconv.Itoa(len(commits)), strconv.Itoa(len(files)), strconv.Itoa(insertions), strconv.Itoa(deletions), strconv.Itoa(staged), strconv.Itoa(unstaged), strconv.Itoa(untracked)}})
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	fileTable, err := toon.Table("files", []string{"status", "path", "kind"}, files)
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	checkout, err := toon.Table("checkout", []string{"index", "worktree", "path"}, checkoutRows(facts.Changes))
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	white, err := toon.Table("whitespace", []string{"clean", "offenses"}, [][]string{{strconv.FormatBool(clean), strconv.Itoa(offenses)}})
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	var b strings.Builder
	b.WriteString(revision)
	b.WriteString(aggregate)
	b.WriteString(fileTable)
	b.WriteString(checkout)
	b.WriteString(white)
	if full {
		logTable, err := toon.Table("log", []string{"sha", "subject"}, commits)
		if err != nil {
			return "", "unrepresentable TOON cell", err.Error()
		}
		body, err := diffBodyAt(root, dr.bodyArgs...)
		if err != nil {
			return "", "git diff failed", err.Error()
		}
		untrackedBody, err := untrackedPatches(root, facts)
		if err != nil {
			return "", "git diff --no-index failed", err.Error()
		}
		b.WriteString(logTable)
		b.WriteString("diff_body:\n")
		b.Write(body)
		b.Write(untrackedBody)
	}
	actions := []axi.Action(nil)
	if !full && len(files) > 0 {
		if dr.method == "explicit" {
			actions = append(actions, axi.InspectFullBase(dr.base))
		} else {
			actions = append(actions, axi.InspectFull(""))
		}
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	b.WriteString(help)
	return b.String(), "", ""
}

func renderCommit(root string, dr diffRange, full bool) (string, string, string) {
	files, err := changedFilesAt(root, dr.filesArgs...)
	if err != nil {
		return "", "git diff --name-status failed", err.Error()
	}
	commits, err := commitLogAt(root, dr.logRange)
	if err != nil {
		return "", "git log failed", err.Error()
	}
	insertions, deletions, err := shortstat(root, dr.filesArgs...)
	if err != nil {
		return "", "git diff --shortstat failed", err.Error()
	}
	revision, err := toon.Table("revision", []string{"commit", "base", "method"}, [][]string{{dr.head, dr.base, dr.method}})
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	aggregate, err := toon.Table("aggregate", []string{"commits", "files", "insertions", "deletions"}, [][]string{{strconv.Itoa(len(commits)), strconv.Itoa(len(files)), strconv.Itoa(insertions), strconv.Itoa(deletions)}})
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	for i := range files {
		files[i] = append(files[i], "")
	}
	fileTable, err := toon.Table("files", []string{"status", "path", "kind"}, files)
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	var b strings.Builder
	b.WriteString(revision)
	b.WriteString(aggregate)
	b.WriteString(fileTable)
	if full {
		logTable, err := toon.Table("log", []string{"sha", "subject"}, commits)
		if err != nil {
			return "", "unrepresentable TOON cell", err.Error()
		}
		body, err := diffBodyAt(root, dr.bodyArgs...)
		if err != nil {
			return "", "git diff failed", err.Error()
		}
		b.WriteString(logTable)
		b.WriteString("diff_body:\n")
		b.Write(body)
	}
	actions := []axi.Action(nil)
	if !full && len(files) > 0 {
		actions = append(actions, axi.InspectFull(dr.head))
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return "", "unrepresentable TOON cell", err.Error()
	}
	b.WriteString(help)
	return b.String(), "", ""
}

func resolveBranchRangeFromFacts(root string, facts git.DiffFacts) (dr diffRange, errKind, errHint string) {
	base, method := "", ""
	if facts.RecordedBase != "" {
		switch {
		case !git.OK("-C", root, "cat-file", "-e", facts.RecordedBase+"^{commit}"):
			method = "merge-base (recorded sha unreachable)"
		case !git.OK("-C", root, "merge-base", "--is-ancestor", facts.RecordedBase, facts.Head):
			method = "merge-base (recorded sha not an ancestor)"
		default:
			base, method = facts.RecordedBase, "recorded"
		}
	}
	if base == "" {
		if !facts.DefaultResolved {
			return diffRange{}, "cannot resolve a review base", "this repository has no resolvable default branch; record one with: git config branch.<name>.benchBase <sha>"
		}
		var err error
		base, err = git.Output("-C", root, "merge-base", facts.DefaultTip, facts.Head)
		if err != nil {
			return diffRange{}, "cannot resolve a review base", "no merge-base with '" + facts.DefaultBranch + "'; record one with: git config branch.<name>.benchBase <sha>"
		}
		if method == "" {
			method = "merge-base"
		}
	}
	return diffRange{
		base:      base,
		head:      facts.Head,
		method:    method,
		filesArgs: []string{base},
		logRange:  base + ".." + facts.Head,
		bodyArgs:  []string{base},
	}, "", ""
}

// ResolveReviewBase is the single source of the resolved review base for
// bench diff and its consumers: the recorded `branch.<name>.benchBase` key when
// it names a reachable ancestor of HEAD, else merge-base with the resolved
// default branch, with method carrying which path answered — `recorded`,
// `merge-base`, or one of the loud fallback labels when a recorded key is
// present but unusable. A non-empty errKind/errHint with an empty base is the
// only absence shape; base is never empty on a nil error.
func ResolveReviewBase(root string) (base, method, errKind, errHint string) {
	base, method = resolveBase()
	if base != "" {
		return base, method, "", ""
	}
	def, ok := git.ResolvedDefault(root)
	if !ok {
		return "", "", "cannot resolve a review base",
			"this repository has no resolvable default branch; record one with: git config branch.<name>.benchBase <sha>"
	}
	mb, err := git.Output("merge-base", def, "HEAD")
	if err != nil {
		return "", "", "cannot resolve a review base",
			"no merge-base with '" + def + "'; record one with: git config branch.<name>.benchBase <sha>"
	}
	if method == "" {
		method = "merge-base"
	}
	return mb, method, "", ""
}

// ChangedFilePaths is the single source of the changed-path set for bench diff and
// its consumers: the path half of changedFiles' rows for a resolved review base,
// carrying bench diff's exact committed+index+tracked-worktree semantics — `git diff
// --name-status --no-renames -z <base>` — so a consumer never re-derives that git
// invocation or the NUL-pair parse.
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

// resolveBase returns the recorded-key base and `recorded` when the key names a
// reachable ancestor, or ("", <loud fallback method>) when the key is present but
// unreachable/divergent, or ("","") when there is no usable recorded key.
func resolveBase() (base, method string) {
	branch, _ := git.Output("symbolic-ref", "--quiet", "--short", "HEAD")
	if branch == "" {
		return "", ""
	}
	key, _ := git.Output("config", "branch."+branch+".benchBase")
	if key == "" {
		return "", ""
	}
	switch {
	case !git.OK("cat-file", "-e", key+"^{commit}"):
		return "", "merge-base (recorded sha unreachable)"
	case !git.OK("merge-base", "--is-ancestor", key, "HEAD"):
		return "", "merge-base (recorded sha not an ancestor)"
	default:
		return key, "recorded"
	}
}
