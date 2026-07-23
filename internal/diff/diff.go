// Package diff ports `bench diff`: the single source of review-base truth for the
// review phase. The recorded `branch.<name>.benchBase` key wins when it names a
// reachable ancestor; otherwise merge-base with the default branch. The preamble
// names which method resolved so a review agent can see a fallback happen, and a
// recorded sha that is unreachable or not an ancestor falls back loudly rather than
// silently diffing from the wrong base. Changed files follow as a
// `files[N]{status,path}:` table with paths escaped exactly once.
//
// The `--full` flag appends the rest of the base-relative picture a review agent
// otherwise hand-runs as two extra git calls: a `log[N]{sha,subject}:` TOON table
// from `git log <base>..HEAD` (two-dot: commits on HEAD since base), and — last,
// behind a fixed `diff_body:` marker line — the raw output of `git diff <base>`
// (committed, index, and tracked worktree changes since the resolved base) passed
// through verbatim, undecorated by TOON so a hunk header or `+`/`-` line survives
// unmangled. Bare `bench diff` is byte-for-byte unaffected by the flag's existence.
//
// `--commit <sha>` bounds the same three sections to one already-landed commit
// instead of the current branch: base becomes the commit's own first parent
// (`git rev-parse --verify <sha>^`), so the range is an exact two-commit diff, not
// a merge-base-relative one — for a merge commit this is everything the merge
// brought in. benchBase/merge-base resolution is skipped entirely; `method: commit
// <sha>` is how a reader sees that the override applies. The sha is verified before
// any section renders: an unresolvable sha or a root commit's missing parent is
// its own structured error, never a leaked git failure.
package diff

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, repeated flags, and help all come from there rather
// than a local switch.
// Help is fullHelp without its trailing newline, because the caller appends one.
var grammar = usage.Grammar{
	Cmd:   "bench diff",
	Help:  strings.TrimSuffix(fullHelp, "\n"),
	Flags: []usage.Flag{{Name: "--full"}, {Name: "--commit", HasValue: true}},
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

const fullHelp = `usage: bench diff [--full] [--commit <sha>]
  --full appends, after the files table, a log[N]{sha,subject} TOON table (git
  log <base>..HEAD) and, last, a verbatim diff_body: block (git diff
  <base>, including tracked index and worktree changes) — the raw diff is passed
  through unescaped, not TOON-encoded. The log remains commit-only.
  A commit subject carrying a control byte makes --full refuse: it exits 1
  with the unrepresentable-TOON-cell error instead of rendering a mangled log
  row.
  --commit <sha> bounds every section to one already-landed commit instead of
  the current branch: base becomes <sha>'s first parent, skipping
  benchBase/merge-base resolution entirely — the fallback for reviewing work
  after it has already merged, when the branch-relative diff is empty. A root
  commit (no parent) or an unresolvable <sha> exits 1 with a structured
  error; a bare --commit (no value) or a repeated --commit exits 2.
`

// diffRange is the resolved review range every rendering section shares: the base
// and method preamble lines, and the argument shapes changedFiles/commitLog/diffBody
// need for either the branch-relative (resolved base through worktree) range or the
// commit-relative (exact first-parent) range.
type diffRange struct {
	base      string
	method    string
	filesArgs []string
	logRange  string
	bodyArgs  []string
}

// resolveCommitRange builds the diffRange for `--commit <sha>`: base is <sha>'s
// resolved first parent. The sha is verified before anything renders — an
// unresolvable sha and a root commit's missing parent are each their own
// structured error (kind, hint), never a leaked git failure.
func resolveCommitRange(commitArg string) (dr diffRange, errKind, errHint string) {
	headSha, err := git.Output("rev-parse", "--verify", commitArg+"^{commit}")
	if err != nil {
		return diffRange{}, "cannot resolve --commit",
			"'" + commitArg + "' does not name a commit reachable in this repository"
	}
	baseSha, err := git.Output("rev-parse", "--verify", commitArg+"^")
	if err != nil {
		return diffRange{}, "--commit has no parent",
			"'" + commitArg + "' is a root commit — there is no first parent to diff against"
	}
	return diffRange{
		base:      baseSha,
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
	base, method := resolveBase()
	if base == "" {
		def := git.DefaultBranch(root)
		mb, err := git.Output("merge-base", def, "HEAD")
		if err != nil {
			return diffRange{}, "cannot resolve a review base",
				"no merge-base with '" + def + "'; record one with: git config branch.<name>.benchBase <sha>"
		}
		base = mb
		if method == "" {
			method = "merge-base"
		}
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
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	var dr diffRange
	var errKind, errHint string
	if hasCommit {
		dr, errKind, errHint = resolveCommitRange(commitArg)
	} else {
		dr, errKind, errHint = resolveBranchRange(root)
	}
	if errKind != "" {
		return toon.Errorf(errKind, errHint) + "\n", 1
	}

	branch, _ := git.Output("symbolic-ref", "--quiet", "--short", "HEAD")
	branchLabel := branch
	if branchLabel == "" {
		branchLabel = "(detached)"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "branch: %s\n", branchLabel)
	fmt.Fprintf(&b, "base: %s\n", dr.base)
	fmt.Fprintf(&b, "method: %s\n", dr.method)
	files, err := changedFiles(dr.filesArgs...)
	if err != nil {
		return toon.Errorf("git diff --name-status failed", err.Error()) + "\n", 1
	}
	tbl, err := toon.Table("files", []string{"status", "path"}, files)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	b.WriteString(tbl)
	if full {
		logRows, err := commitLog(dr.logRange)
		if err != nil {
			return toon.Errorf("git log failed", err.Error()) + "\n", 1
		}
		logTbl, err := toon.Table("log", []string{"sha", "subject"}, logRows)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		b.WriteString(logTbl)
		b.WriteString("diff_body:\n")
		body, err := diffBody(dr.bodyArgs...)
		if err != nil {
			return toon.Errorf("git diff failed", err.Error()) + "\n", 1
		}
		b.Write(body)
	}
	return b.String(), 0
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
