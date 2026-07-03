// Package learnings ports `bench learnings`: the open journal headings of
// .bench/learnings.md as a `learnings[N]{date,title}:` TOON table. The heading parser
// is the single source `bench learnings` and (through the Go binary) `bench status`
// both read, so the open-count and the row listing count by one rule.
package learnings

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// openRe selects an open heading: `## <ISO date> … [open]`. The same shape the shell
// LEARNINGS_OPEN_RE matched; a template example without a real date never matches.
var openRe = regexp.MustCompile(`^## [0-9]{4}-[0-9]{2}-[0-9]{2}.*\[open\]`)

// Rows parses the open headings of a learnings journal into date/title rows. A
// trailing CR (CRLF journal) is stripped first; the title is the heading minus its
// `## <date>` prefix, the following separator run (spaces, an ASCII hyphen, or an
// em-dash), and the trailing `[open]…`, with surrounding whitespace removed.
func Rows(content []byte) [][]string {
	var rows [][]string
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSuffix(line, "\r")
		if !openRe.MatchString(line) {
			continue
		}
		date, title := parseHeading(line)
		rows = append(rows, []string{date, title})
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
	switch {
	case len(args) == 0:
	case args[0] == "-h" || args[0] == "--help":
		return "usage: bench learnings\n", 0
	default:
		return toon.Usage("bench learnings", args[0]) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	content, _ := os.ReadFile(filepath.Join(root, ".bench", "learnings.md"))
	return toon.Table("learnings", []string{"date", "title"}, Rows(content)), 0
}
