// history.go implements `bench spec history <slug>`. It compiles the hand-run recovery
// incantation `/bench-debug` and `/bench-write-spec` both name: `git log
// --grep=spec-retire: <slug>` plus `git log --diff-filter=D -- specs/<slug>.md`. It
// merges the two into one compiled, deduped, AXI-conformant query. It never writes. It
// is a read-only sibling of Flip and retireCommand, and shares their slugOf/specArg
// argument convention.
package spec

import (
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// historyLogFormat is shared by both git log queries below: full hash, abbreviated
// hash, committer date in strict ISO-8601, and the subject. A lexical sort on ISO-8601
// is a chronological sort, and its first 10 bytes are the displayed YYYY-MM-DD. The
// format is NUL-delimited, so a comma or quote in the subject arrives raw for the
// caller to TOON-escape a single layer downstream — the same discipline
// internal/diff's parseLogFormat uses.
const historyLogFormat = "--format=%H%x00%h%x00%cI%x00%s"

// historyEntry is one merged, deduped row: a commit tagged retire or delete. The row is
// keyed by its full hash for dedupe. It carries the full ISO-8601 timestamp for the
// sort; the rendered date is this field's first 10 bytes.
type historyEntry struct {
	full, short, iso, kind, subject string
}

// HistoryFact is one typed retirement/deletion record.
type HistoryFact struct {
	Slug, Hash, Date, Kind, Subject string
}

// History returns the merged history facts for slug, newest first.
func History(slug string) ([]HistoryFact, error) {
	retireRaw, err := historyRetireLog(slug)
	if err != nil {
		return nil, err
	}
	deleteRaw, err := historyDeleteLog(slug)
	if err != nil {
		return nil, err
	}
	var retires []historyEntry
	for _, e := range parseHistoryLog(retireRaw, "retire") {
		if retireTokenMatches(e.subject, slug) {
			retires = append(retires, e)
		}
	}
	merged := mergeHistory(retires, parseHistoryLog(deleteRaw, "delete"))
	out := make([]HistoryFact, len(merged))
	for i, e := range merged {
		date := e.iso
		if len(date) > 10 {
			date = date[:10]
		}
		out[i] = HistoryFact{Slug: slug, Hash: e.short, Date: date, Kind: e.kind, Subject: e.subject}
	}
	return out, nil
}

// parseHistoryLog turns one `git log <historyLogFormat>` invocation's raw output into
// rows tagged with kind. It is the one parser for both queries below, so the
// NUL-framing and 4-field shape has a single source.
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

// mergeHistory merges lists that parseHistoryLog already tagged with their own kind. It
// dedupes by full commit hash and keeps the first-seen kind. Callers pass the retire
// list first, so a commit matching both queries keeps the `retire` tag rather than a
// duplicate or an overwrite. The common case is a `bench spec retire` commit that both
// deletes the file and carries the message. The function sorts the result newest-first
// by the full ISO-8601 timestamp. A lexical sort on ISO-8601 is exactly a chronological
// sort, so it needs no time parsing.
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
	// Sort newest first with a simple insertion-sort-stable selection pass. The list
	// stays small — retirement history, not a full repo log — so this needs no import
	// beyond what is already here.
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].iso > out[j-1].iso; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// historyRetireLog runs the retire-message query: every commit whose message literally
// contains `spec-retire: <slug>`. It passes --fixed-strings, so a slug carrying regex
// metacharacters like `.` or `*` matches as literal text, not a pattern. --grep is a
// contains match, so this is a coarse candidate filter only. A slug that is a prefix of
// another slug — `dash` vs `dashboard` — matches both messages. retireTokenMatches
// applies the exact cut on the parsed subject.
func historyRetireLog(slug string) ([]byte, error) {
	return git.Raw("log", "--fixed-strings", "--grep=spec-retire: "+slug, historyLogFormat)
}

// retireTokenMatches reports whether a commit subject carries the retire token for
// exactly this slug: `spec-retire: <slug>` extending to the end of the subject. Slugs
// may contain spaces, so no word boundary can terminate one. End-of-subject is the only
// sound terminator, and it is also the commit-subject shape `bench spec retire`
// instructs.
func retireTokenMatches(subject, slug string) bool {
	return strings.HasSuffix(subject, "spec-retire: "+slug)
}

// historyDeleteLog runs the file-deletion query: every commit that deleted both
// historical specs/<slug>.md and folder specs/<slug>/spec.md. The pathspecs carry two
// magic words. `literal` matches a slug containing `*`, `?`, or `[` as an exact path
// rather than a glob. `top` resolves the path from the repository root regardless of
// the process cwd, because a subdirectory invocation must see the same history a
// repo-root invocation does.
func historyDeleteLog(slug string) ([]byte, error) {
	flat := ":(literal,top)specs/" + slug + ".md"
	folder := ":(literal,top)specs/" + slug + "/spec.md"
	return git.Raw("log", "--diff-filter=D", historyLogFormat, "--", flat, folder)
}

// historyCommand runs `bench spec history <slug>`. It merges the two queries above,
// dedupes, and renders `history[N]{hash,date,kind,subject}:` — the definitive empty
// state when neither query matches. It reports structured stdout errors on a git
// failure or an unrepresentable TOON cell, such as a control byte in a commit subject.
// Exit codes are honest: 0 for ok, including the empty state, 1 for a git or render
// failure, and 2 for a usage error.
func historyCommand(rest []string) (string, int) {
	arg, out, code, ok := specArg("bench spec history", "usage: bench spec history <spec.md | slug>\n", rest)
	if !ok {
		return out, code
	}
	if _, err := git.Root(); err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	slug := slugOf(arg)

	facts, err := History(slug)
	if err != nil {
		return toon.Errorf("git history derivation failed", err.Error()) + "\n", 1
	}
	rows := make([][]string, len(facts))
	for i, e := range facts {
		rows[i] = []string{e.Hash, e.Date, e.Kind, e.Subject}
	}
	tbl, err := toon.Table("history", []string{"hash", "date", "kind", "subject"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return tbl, 0
}
