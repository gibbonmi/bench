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
// behind a fixed `diff_body:` marker line — the raw output of `git diff
// <base>...HEAD` (three-dot: changes on the HEAD side since merge-base) passed
// through verbatim, undecorated by TOON so a hunk header or `+`/`-` line survives
// unmangled. Bare `bench diff` is byte-for-byte unaffected by the flag's existence.
package diff

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

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

func changedFiles(base string) ([][]string, error) {
	raw, err := git.Raw("diff", "--name-status", "--no-renames", "-z", base+"...HEAD")
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

func commitLog(base string) ([][]string, error) {
	raw, err := git.Raw("log", "--format=%h%x00%s", base+"..HEAD")
	if err != nil {
		return nil, err
	}
	return parseLogFormat(raw), nil
}

const fullHelp = `usage: bench diff [--full]
  --full appends, after the files table, a log[N]{sha,subject} TOON table (git
  log <base>..HEAD) and, last, a verbatim diff_body: block (git diff
  <base>...HEAD) — the raw diff is passed through unescaped, not TOON-encoded.
  A commit subject carrying a control byte makes --full refuse: it exits 1
  with the unrepresentable-TOON-cell error instead of rendering a mangled log
  row.
`

// Command implements `bench diff`.
func Command(args []string) (string, int) {
	full := false
	switch {
	case len(args) == 0:
	case args[0] == "-h" || args[0] == "--help":
		return fullHelp, 0
	case len(args) == 1 && args[0] == "--full":
		full = true
	case len(args) == 2 && args[0] == "--full":
		// A recognized leading flag followed by a second argument: name the second
		// argument as the offender, not the flag that parsed fine.
		return toon.Usage("bench diff", args[1]) + "\n", 2
	default:
		return toon.Usage("bench diff", args[0]) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	base, method := resolveBase()
	if base == "" {
		def := git.DefaultBranch(root)
		mb, err := git.Output("merge-base", def, "HEAD")
		if err != nil {
			return toon.Errorf("cannot resolve a review base",
				"no merge-base with '"+def+"'; record one with: git config branch.<name>.benchBase <sha>") + "\n", 1
		}
		base = mb
		if method == "" {
			method = "merge-base"
		}
	}

	branch, _ := git.Output("symbolic-ref", "--quiet", "--short", "HEAD")
	branchLabel := branch
	if branchLabel == "" {
		branchLabel = "(detached)"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "branch: %s\n", branchLabel)
	fmt.Fprintf(&b, "base: %s\n", base)
	fmt.Fprintf(&b, "method: %s\n", method)
	files, err := changedFiles(base)
	if err != nil {
		return toon.Errorf("git diff --name-status failed", err.Error()) + "\n", 1
	}
	tbl, err := toon.Table("files", []string{"status", "path"}, files)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	b.WriteString(tbl)
	if full {
		logRows, err := commitLog(base)
		if err != nil {
			return toon.Errorf("git log failed", err.Error()) + "\n", 1
		}
		logTbl, err := toon.Table("log", []string{"sha", "subject"}, logRows)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		b.WriteString(logTbl)
		b.WriteString("diff_body:\n")
		body, err := git.Raw("diff", base+"...HEAD")
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
