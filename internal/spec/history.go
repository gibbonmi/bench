// history.go implements `bench spec history <slug>`: the compiled port of the
// hand-run recovery incantation `/bench-debug` and `/bench-write-spec` both name —
// `git log --grep=spec-retire: <slug>` plus `git log --diff-filter=D -- specs/<slug>.md`
// — into one compiled, merged, deduped, AXI-conformant query. It never writes; it is a
// read-only sibling of Flip/retireCommand, sharing their slugOf/specArg argument
// convention.
package spec

import (
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// historyLogFormat is shared by both git log queries below: full hash, abbreviated
// hash, committer date in strict ISO-8601 (lexically sortable = chronologically
// sortable, and its first 10 bytes are the displayed YYYY-MM-DD), and the subject —
// NUL-delimited so a comma or quote in the subject arrives raw for the caller to
// TOON-escape a single layer downstream, the same discipline internal/diff's
// parseLogFormat uses.
const historyLogFormat = "--format=%H%x00%h%x00%cI%x00%s"

// historyEntry is one merged, deduped row: a commit tagged retire or delete, keyed by
// its full hash for dedupe and carrying the full ISO-8601 timestamp for the sort
// (the rendered date is this field's first 10 bytes).
type historyEntry struct {
	full, short, iso, kind, subject string
}

// parseHistoryLog turns one `git log <historyLogFormat>` invocation's raw output into
// rows tagged with kind — the one parser for both queries below, so the NUL-framing
// and 4-field shape has a single source.
func parseHistoryLog(raw []byte, kind string) []historyEntry {
	s := strings.TrimRight(string(raw), "\n")
	if s == "" {
		return nil
	}
	var out []historyEntry
	for _, line := range strings.Split(s, "\n") {
		fields := strings.SplitN(line, "\x00", 4)
		if len(fields) != 4 {
			continue
		}
		out = append(out, historyEntry{full: fields[0], short: fields[1], iso: fields[2], kind: kind, subject: fields[3]})
	}
	return out
}

// mergeHistory merges lists (each already tagged with its own kind by
// parseHistoryLog), dedupes by full commit hash keeping the first-seen kind — callers
// pass the retire list first, so a commit matching both queries (the common case: a
// `bench spec retire` commit both deletes the file and carries the message) keeps the
// `retire` tag rather than being overwritten or duplicated — and sorts the result
// newest-first by the full ISO-8601 timestamp (a lexical sort on ISO-8601 is exactly a
// chronological sort, so no time parsing is needed).
func mergeHistory(lists ...[]historyEntry) []historyEntry {
	seen := make(map[string]bool)
	var out []historyEntry
	for _, list := range lists {
		for _, e := range list {
			if seen[e.full] {
				continue
			}
			seen[e.full] = true
			out = append(out, e)
		}
	}
	// Newest first; insertion-sort-stable via a simple selection pass keeps the
	// dependency footprint at zero (the list is small — retirement history, not
	// full repo log) and needs no import beyond what's already here.
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].iso > out[j-1].iso; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// historyRetireLog runs the retire-message query: every commit whose message
// literally contains `spec-retire: <slug>` (--fixed-strings, so a slug carrying regex
// metacharacters like `.` or `*` matches as literal text, not a pattern). --grep is a
// contains match, so this is a coarse candidate filter only — a slug that is a prefix
// of another slug (`dash` vs `dashboard`) matches both messages; retireTokenMatches
// applies the exact cut on the parsed subject.
func historyRetireLog(slug string) ([]byte, error) {
	return git.Raw("log", "--fixed-strings", "--grep=spec-retire: "+slug, historyLogFormat)
}

// retireTokenMatches reports whether a commit subject carries the retire token for
// exactly this slug: `spec-retire: <slug>` extending to the end of the subject.
// Slugs may contain spaces, so no word boundary can terminate one — end-of-subject
// is the only sound terminator, which is also the commit-subject shape
// `bench spec retire` instructs.
func retireTokenMatches(subject, slug string) bool {
	return strings.HasSuffix(subject, "spec-retire: "+slug)
}

// historyDeleteLog runs the file-deletion query: every commit that deleted
// specs/<slug>.md. The pathspec carries two magic words: `literal` so a slug
// containing `*`/`?`/`[` is matched as an exact path rather than a glob, and `top` so
// the path resolves from the repository root regardless of the process cwd — a
// subdirectory invocation must see the same history a repo-root invocation does.
func historyDeleteLog(slug string) ([]byte, error) {
	pathspec := ":(literal,top)specs/" + slug + ".md"
	return git.Raw("log", "--diff-filter=D", historyLogFormat, "--", pathspec)
}

// historyCommand runs `bench spec history <slug>`: merges the two queries above,
// dedupes, and renders `history[N]{hash,date,kind,subject}:` — the definitive empty
// state when neither query matches, structured stdout errors on a git failure or an
// unrepresentable TOON cell (a control byte in a commit subject), and honest exit
// codes (0 ok including the empty state, 1 a git/render failure, 2 a usage error).
func historyCommand(rest []string) (string, int) {
	arg, out, code, ok := specArg("bench spec history", "usage: bench spec history <spec.md | slug>\n", rest)
	if !ok {
		return out, code
	}
	if _, err := git.Root(); err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	slug := slugOf(arg)

	retireRaw, err := historyRetireLog(slug)
	if err != nil {
		return toon.Errorf("git log --grep failed", err.Error()) + "\n", 1
	}
	deleteRaw, err := historyDeleteLog(slug)
	if err != nil {
		return toon.Errorf("git log --diff-filter=D failed", err.Error()) + "\n", 1
	}

	var retireEntries []historyEntry
	for _, e := range parseHistoryLog(retireRaw, "retire") {
		if retireTokenMatches(e.subject, slug) {
			retireEntries = append(retireEntries, e)
		}
	}
	deleteEntries := parseHistoryLog(deleteRaw, "delete")
	merged := mergeHistory(retireEntries, deleteEntries)

	rows := make([][]string, len(merged))
	for i, e := range merged {
		date := e.iso
		if len(date) > 10 {
			date = date[:10]
		}
		rows[i] = []string{e.short, date, e.kind, e.subject}
	}
	tbl, err := toon.Table("history", []string{"hash", "date", "kind", "subject"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return tbl, 0
}
