// Package learnings ports `bench learnings`: the open journal headings of
// capture/learnings.md as a `learnings[N]{date,title}:` TOON table. The heading parser
// is the single source `bench learnings` and (through the Go binary) `bench status`
// both read, so the open-count and the row listing count by one rule.
package learnings

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:  "bench learnings",
	Help: "usage: bench learnings",
}

// Entry is one typed journal item. Body preserves the source text below the
// heading so roadmap context and the human projections share this parser.
type Entry struct {
	Date, Title, State, Body string
	Line                     int
}

// Malformed is one heading that started a `## ` entry but did not parse as a dated
// one. Line is the 1-based source line, carried so a surface can name where to look.
type Malformed struct {
	Reason, Raw string
	Line        int
}

// Entries parses every well-formed dated journal heading in document order.
// Rows is deliberately a projection of this typed representation.
func Entries(content []byte) []Entry {
	entries, _ := Parse(content)
	return entries
}

// Parse returns typed entries plus every malformed heading fragment.
func Parse(content []byte) ([]Entry, []Malformed) {
	lines := strings.Split(string(content), "\n")
	var out []Entry
	var malformed []Malformed
	for i := 0; i < len(lines); {
		line := strings.TrimSuffix(lines[i], "\r")
		if !isDatedHeading(line) {
			switch {
			case strings.HasPrefix(line, "## "):
				// The shipped scaffold's own worked example (internal/adopt seeds every fresh
				// repo with it under "Format per entry:") is documentation, not a broken record,
				// so an unedited template never counts as malformed.
				if !isTemplatePlaceholder(line) {
					malformed = append(malformed, Malformed{Reason: "malformed learning heading", Raw: line, Line: i + 1})
				}
			case isLostDatedLine(line):
				malformed = append(malformed, Malformed{Reason: lostDatedLineReason, Raw: line, Line: i + 1})
			}
			i++
			continue
		}
		date, title, state := parseHeading(line)
		start := i + 1
		i = start
		// A dated line inside the body is collected rather than appended as it is found,
		// so a heading that is itself malformed keeps its record ahead of the body's and
		// every reason stays in ascending source-line order.
		var bodyLost []Malformed
		for i < len(lines) && !strings.HasPrefix(strings.TrimSuffix(lines[i], "\r"), "## ") {
			if bodyLine := strings.TrimSuffix(lines[i], "\r"); isLostDatedLine(bodyLine) {
				bodyLost = append(bodyLost, Malformed{Reason: lostDatedLineReason, Raw: bodyLine, Line: i + 1})
			}
			i++
		}
		body := strings.Join(lines[start:i], "\n")
		body = strings.Trim(body, "\n")
		if state != "open" {
			malformed = append(malformed, Malformed{Reason: "dated learning heading must end with [open]", Raw: line, Line: start})
		} else {
			out = append(out, Entry{Date: date, Title: title, State: state, Body: body, Line: start})
		}
		malformed = append(malformed, bodyLost...)
	}
	return out, malformed
}

// Rows parses the open headings of a learnings journal into date/title rows. A
// trailing CR (CRLF journal) is stripped first; the title is the heading minus its
// `## <date>` prefix, the following separator run (spaces, an ASCII hyphen, or an
// em-dash), and the trailing `[open]…`, with surrounding whitespace removed.
func Rows(content []byte) [][]string { return openRows(Entries(content)) }

// openRows projects the open-state entries to date/title rows — the one source Rows
// and Command both narrow to, so the listing and the open-count agree by construction.
func openRows(entries []Entry) [][]string {
	var rows [][]string
	for _, e := range entries {
		if e.State != "open" {
			continue
		}
		rows = append(rows, []string{e.Date, e.Title})
	}
	return rows
}

func parseHeading(line string) (date, title, state string) {
	rest := strings.TrimPrefix(line, "## ")
	date = rest
	if i := strings.IndexFunc(rest, isSpace); i >= 0 {
		date = rest[:i]
	}
	title = stripLeadingSeparators(strings.TrimPrefix(rest, date))
	if i := strings.LastIndex(title, "["); i >= 0 && strings.HasSuffix(title, "]") {
		state = title[i+1 : len(title)-1]
		title = title[:i]
	}
	title = strings.TrimRightFunc(title, isSpace)
	return date, title, state
}

func isDatedHeading(line string) bool {
	return strings.HasPrefix(line, "## ") && opensWithDate(line[len("## "):])
}

// lostDatedLineReason is the reason a dated line that is not a well-formed heading
// carries. It names the writer's repair — use a heading — rather than the parser's
// disappointment, and it is distinct from the two `## ` reasons so a reader can tell
// "you used the wrong marker" from "your heading is broken".
const lostDatedLineReason = "dated learning entry is not a heading"

// isLostDatedLine reports whether line leads with a date but is not a heading, which
// is how a writer loses an entry: appended as a bullet, a quote, or plain text, it
// parses to nothing today. A line already starting `## ` is excluded outright, because
// the two heading reasons own it and a second record would double-report it.
//
// The prefix walk strips a run — possibly empty, so a date flush at column one is
// still reached — of runes that are each either whitespace or one of the markdown
// markers a writer reaches for. unicode.IsSpace is the exact predicate rather than the
// ASCII isSpace below: this serves hand-edited markdown, where a pasted U+00A0 or
// U+3000 must not re-open the silent drop, while the zero-width U+200B and U+FEFF are
// not White_Space and stay non-separators a reader can see.
func isLostDatedLine(line string) bool {
	if strings.HasPrefix(line, "## ") {
		return false
	}
	return opensWithDate(strings.TrimLeftFunc(line, func(r rune) bool {
		return unicode.IsSpace(r) || strings.ContainsRune("-*+>#", r)
	}))
}

// opensWithDate reports whether s begins with a `YYYY-MM-DD` digit shape. It is the one
// definition of the journal's date grammar, shared by the heading rule and the lost-line
// rule so the two cannot drift apart, and it is deliberately shape-only: a calendar parse
// here would judge `2026-88-88` differently from the heading rule that already accepts it.
func opensWithDate(s string) bool {
	const shape = "2006-01-02"
	if len(s) < len(shape) {
		return false
	}
	for i, b := range []byte(s[:len(shape)]) {
		if i == 4 || i == 7 {
			if b != '-' {
				return false
			}
			continue
		}
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// isTemplatePlaceholder reports whether line is the scaffold's own worked example, whose
// date field is the literal `<date>` token rather than a date.
func isTemplatePlaceholder(line string) bool {
	rest := strings.TrimPrefix(line, "## ")
	date := rest
	if i := strings.IndexFunc(rest, isSpace); i >= 0 {
		date = rest[:i]
	}
	return date == "<date>"
}

// hasAnyHeading reports whether content attempts the journal shape at all — any line
// beginning `## `, dated, malformed, or the template placeholder alike. Its absence is
// what unsupported-schema means for this document: bytes that never attempt a heading,
// as distinct from a heading attempt that failed to parse.
func hasAnyHeading(content []byte) bool {
	for _, raw := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(strings.TrimSuffix(raw, "\r"), "## ") {
			return true
		}
	}
	return false
}

func hasJournalSchema(content []byte) bool {
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSuffix(raw, "\r")
		if line == "" {
			continue
		}
		return line == JournalSchemaHeading
	}
	return false
}

func stripLeadingSeparators(s string) string {
	for {
		switch {
		case strings.HasPrefix(s, " "):
			s = s[1:]
		case strings.HasPrefix(s, "-"):
			s = s[1:]
		case strings.HasPrefix(s, "—"):
			s = s[len("—"):]
		default:
			return s
		}
	}
}

// isSpace adapts toon.IsSpace (the one source of the AXI whitespace class) to a rune
// predicate for IndexFunc/TrimRightFunc; a multibyte rune is >= 0x80 and never a space.
func isSpace(r rune) bool { return r < 0x80 && toon.IsSpace(byte(r)) }

// JournalSchemaHeading identifies a zero-entry learnings journal.
const JournalSchemaHeading = "# Learnings — usage journal"

// JournalPath is the repo-relative journal. It is exported because the name is one
// fact with three readers — this command, the roadmap drain that counts its open
// headings, and the status row that names it when the read fails — and a literal
// repeated at each of them is how the three drift apart.
const JournalPath = "capture/learnings.md"

// Command implements `bench learnings`. Unknown argument → usage on stdout, exit 2;
// outside a repo → structured error on stdout, exit 1. Absence is the only
// authoritative empty state: a missing journal renders the empty table at exit 0. Any
// other non-parsed classifier state — empty, malformed bytes, unreadable, wrong-type —
// or a document that never attempts a `## ` heading at all (unsupported-schema) exits 1
// with a structured `error:` line naming the state, so a read failure can never print
// as an empty table. A parsed document with malformed headings among its well-formed
// entries renders every well-formed row plus one row per malformed heading and still
// exits 1 (story 9): a broken entry is surfaced, never silently dropped.
func Command(args []string) (string, int) {
	if _, line, code := usage.Parse(grammar, args); line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	c := bounds.Classify(filepath.Join(root, JournalPath), bounds.ControlRecordLimit)
	switch c.State {
	case bounds.StateAbsent:
		out, err := toon.Table("learnings", []string{"date", "title"}, nil)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		help, err := axi.RenderHelp(nil)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		return out + help, 0
	case bounds.StateParsed:
		// falls through to the structural parse below
	default:
		return toon.RecordError(JournalPath, c.State, c.Reason) + "\n", 1
	}
	if !hasAnyHeading(c.Data) && !hasJournalSchema(c.Data) {
		return toon.RecordError(JournalPath, bounds.StateUnsupportedSchema, "no dated heading found") + "\n", 1
	}
	entries, malformed := Parse(c.Data)
	rows := openRows(entries)
	for _, m := range malformed {
		rows = append(rows, []string{fmt.Sprintf("line %d", m.Line), m.Reason})
	}
	out, err := toon.Table("learnings", []string{"date", "title"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	actions := make([]axi.Action, 0, len(entries))
	for _, entry := range entries {
		actions = append(actions, axi.HarnessPhase("/bench-drain", "verdict "+entry.Date+": "+entry.Title))
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	out += help
	if len(malformed) > 0 {
		return out, 1
	}
	return out, 0
}
