// Package learnings ports `bench learnings`: the open journal headings of
// .bench/learnings.md as a `learnings[N]{date,title}:` TOON table. The heading parser
// is the single source `bench learnings` and (through the Go binary) `bench status`
// both read, so the open-count and the row listing count by one rule.
package learnings

import (
	"os"
	"path/filepath"
	"strings"

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
}

type Malformed struct{ Reason, Raw string }

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
		if !strings.HasPrefix(line, "## ") || len(line) < 13 || line[7] != '-' || line[10] != '-' {
			if strings.HasPrefix(line, "## ") {
				malformed = append(malformed, Malformed{"malformed learning heading", line})
			}
			i++
			continue
		}
		date, title := parseHeading(line)
		state := ""
		if j := strings.LastIndex(line, "["); j >= 0 && strings.HasSuffix(line, "]") {
			state = line[j+1 : len(line)-1]
		}
		start := i + 1
		i = start
		for i < len(lines) && !strings.HasPrefix(strings.TrimSuffix(lines[i], "\r"), "## ") {
			i++
		}
		body := strings.Join(lines[start:i], "\n")
		body = strings.Trim(body, "\n")
		out = append(out, Entry{Date: date, Title: title, State: state, Body: body})
	}
	return out, malformed
}

// Rows parses the open headings of a learnings journal into date/title rows. A
// trailing CR (CRLF journal) is stripped first; the title is the heading minus its
// `## <date>` prefix, the following separator run (spaces, an ASCII hyphen, or an
// em-dash), and the trailing `[open]…`, with surrounding whitespace removed.
func Rows(content []byte) [][]string {
	var rows [][]string
	for _, e := range Entries(content) {
		if e.State != "open" {
			continue
		}
		rows = append(rows, []string{e.Date, e.Title})
	}
	return rows
}

func parseHeading(line string) (date, title string) {
	rest := strings.TrimPrefix(line, "## ")
	date = rest
	if i := strings.IndexFunc(rest, isSpace); i >= 0 {
		date = rest[:i]
	}
	title = stripLeadingSeparators(strings.TrimPrefix(rest, date))
	// Strip from the last `[open]` (the shell `%[open]*` shortest-suffix removal),
	// then trim trailing whitespace so no padding rides into the field.
	if i := strings.LastIndex(title, "[open]"); i >= 0 {
		title = title[:i]
	}
	title = strings.TrimRightFunc(title, isSpace)
	return date, title
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

// Command implements `bench learnings`. Unknown argument → usage on stdout, exit 2;
// outside a repo → structured error on stdout, exit 1; otherwise the TOON table, exit 0.
func Command(args []string) (string, int) {
	if _, line, code := usage.Parse(grammar, args); line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	content, _ := os.ReadFile(filepath.Join(root, ".bench", "learnings.md"))
	out, err := toon.Table("learnings", []string{"date", "title"}, Rows(content))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 0
}
