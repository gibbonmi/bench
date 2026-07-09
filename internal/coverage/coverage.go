// Package coverage ports `bench coverage`: the acceptance-coverage-map parser the
// gate's docs layer and the review phase both consume. Extraction mode emits the
// spec's state and rows as TOON; `--check` validates the map (canonical header,
// five-cell rows, non-empty cells, in-range story references, historical opt-out).
// The validation phrasings are load-bearing — downstream consumers match them by
// substring — so this is the one validator for the convention.
package coverage

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	specref "github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/toon"
)

const escSentinel = "\x1c" // stands in for an escaped `\|` across the pipe split

var (
	storyNumRe = regexp.MustCompile(`^[0-9]+\. `)
	mapEndRe   = regexp.MustCompile(`^#{2,} `)
	edgeRe     = regexp.MustCompile(`^[Ee][Dd][Gg][Ee]`)
	nonDigitRe = regexp.MustCompile(`[^0-9]+`)
	parenRe    = regexp.MustCompile(`[ \t]*\(.*\)$`)
	// A story reference: a number, a range (en-dash or hyphen), or a comma list of them.
	storyRefRe = regexp.MustCompile(`^[0-9]+([ \t]*(–|-)[ \t]*[0-9]+)?([ \t]*,[ \t]*[0-9]+([ \t]*(–|-)[ \t]*[0-9]+)?)*$`)
)

var fieldNames = [5]string{"story", "behavior", "seam", "red signal", "why it catches the failure"}

type dataRow struct {
	ncells int
	all    [5]string // c[1..5], trimmed; empty beyond ncells
}

// parsed is the result of one scan of a spec: whether the map header was seen, the
// historical opt-out, whether the parsed header matched the canonical one, the highest
// story number, and the data rows.
type parsed struct {
	seen       bool
	historical bool
	gotHeader  bool
	headerOK   bool
	maxStory   int
	dataRows   []dataRow
}

func parse(content []byte) parsed {
	var p parsed
	inStories, inMap := false, false
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSuffix(raw, "\r")
		if line == "<!-- coverage-map: historical -->" {
			p.historical = true
		}
		if line == "## User stories" {
			inStories = true
			continue
		}
		if inStories && strings.HasPrefix(line, "## ") {
			inStories = false
		}
		if inStories && storyNumRe.MatchString(line) {
			num := line
			if i := strings.IndexByte(num, '.'); i >= 0 {
				num = num[:i]
			}
			if n, err := strconv.Atoi(num); err == nil && n > p.maxStory {
				p.maxStory = n
			}
		}
		if line == "### Acceptance coverage map" && !p.seen {
			p.seen = true
			inMap = true
			continue
		}
		if inMap && mapEndRe.MatchString(line) {
			inMap = false
		}
		if inMap {
			p.processMapLine(line)
		}
	}
	return p
}

func (p *parsed) processMapLine(raw string) {
	line := trimSpaceTab(raw)
	if !strings.HasPrefix(line, "|") {
		return
	}
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	line = strings.ReplaceAll(line, `\|`, escSentinel)
	cells := strings.Split(line, "|")
	allSep := true
	for i, c := range cells {
		c = strings.ReplaceAll(c, escSentinel, `\|`)
		c = trimSpaceTab(c)
		cells[i] = c
		if !isDashes(c) {
			allSep = false
		}
	}
	if !p.gotHeader {
		p.gotHeader = true
		p.headerOK = strings.ToLower(strings.Join(cells, "|")) == "story|behavior|seam|red signal|why it catches the failure"
		return
	}
	if allSep {
		return
	}
	dr := dataRow{ncells: len(cells)}
	for i := 0; i < 5 && i < len(cells); i++ {
		dr.all[i] = cells[i]
	}
	p.dataRows = append(p.dataRows, dr)
}

// State classifies the spec: no-map (no coverage header), historical (opted out), or
// mapped.
func State(p parsed) string {
	switch {
	case !p.seen:
		return "no-map"
	case p.historical:
		return "historical"
	default:
		return "mapped"
	}
}

// Rows returns story/seam/red_signal per data row for a mapped spec; nil otherwise.
func Rows(p parsed) [][]string {
	if State(p) != "mapped" {
		return nil
	}
	var rows [][]string
	for _, r := range p.dataRows {
		rows = append(rows, []string{r.all[0], r.all[2], r.all[3]})
	}
	return rows
}

// Check returns one violation message per problem for a mapped spec; a historical or
// no-map spec is silent (nil). The phrasings are matched by substring downstream.
func Check(p parsed) []string {
	if State(p) != "mapped" {
		return nil
	}
	if !p.headerOK {
		return []string{"coverage map missing the canonical header"}
	}
	if len(p.dataRows) == 0 {
		return []string{"coverage map has no data rows"}
	}
	var v []string
	for idx, r := range p.dataRows {
		rn := idx + 1
		if r.ncells != 5 {
			v = append(v, fmt.Sprintf("coverage map row %d has %d cells (want 5)", rn, r.ncells))
			continue
		}
		for i := 0; i < 5; i++ {
			if r.all[i] == "" {
				v = append(v, fmt.Sprintf("coverage map row %d has an empty '%s' cell", rn, fieldNames[i]))
			}
		}
		story := parenRe.ReplaceAllString(r.all[0], "")
		if story == "" || edgeRe.MatchString(story) {
			continue
		}
		if storyRefRe.MatchString(story) {
			for _, tok := range strings.Fields(nonDigitRe.ReplaceAllString(story, " ")) {
				if n, err := strconv.Atoi(tok); err == nil && n > p.maxStory {
					v = append(v, fmt.Sprintf("coverage map row %d references story %s but the spec numbers only %d", rn, tok, p.maxStory))
				}
			}
		} else {
			v = append(v, fmt.Sprintf("coverage map row %d has an unrecognized story reference '%s'", rn, story))
		}
	}
	return v
}

func trimSpaceTab(s string) string { return strings.Trim(s, " \t") }

func isDashes(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != '-' {
			return false
		}
	}
	return true
}

// Command implements `bench coverage [--check] <spec.md | slug>`.
func Command(args []string) (string, int) {
	check := false
	spec := ""
	for _, a := range args {
		switch {
		case a == "--check":
			check = true
		case a == "-h" || a == "--help":
			return "usage: bench coverage [--check] <spec.md | slug>\n", 0
		case strings.HasPrefix(a, "-"):
			return toon.Usage("bench coverage", a) + "\n", 2
		default:
			if spec != "" {
				return toon.Usage("bench coverage", a) + "\n", 2
			}
			spec = a
		}
	}
	if spec == "" {
		return "usage: bench coverage (missing argument: <spec.md> is required)\n", 2
	}
	content, resolved, tried, ok, err := specref.Resolve("", spec)
	if err != nil {
		return toon.Errorf("spec not readable: "+err.Error(), "fix the file's permissions or pass another spec") + "\n", 1
	}
	if !ok {
		return toon.Errorf("spec not found: "+strings.Join(tried, ", "), "pass a path to a spec markdown file") + "\n", 1
	}
	spec = resolved
	p := parse(content)
	if check {
		violations := Check(p)
		if len(violations) == 0 {
			// A mapped spec with zero violations gets a definitive pass line — silent
			// success is indistinguishable from a check that produced no output by
			// accident. A historical or no-map spec has nothing to validate, so it
			// keeps the AXI empty-state posture (silent, exit 0).
			if State(p) == "mapped" {
				return fmt.Sprintf("ok: coverage map valid — %d row(s)\n", len(p.dataRows)), 0
			}
			return "", 0
		}
		var b strings.Builder
		for _, viol := range violations {
			fmt.Fprintln(&b, toon.Errorf(spec+" "+viol, "fix the map or mark it <!-- coverage-map: historical -->"))
		}
		return b.String(), 1
	}
	var b strings.Builder
	fmt.Fprintf(&b, "spec: %s\n", spec)
	fmt.Fprintf(&b, "state: %s\n", State(p))
	tbl, err := toon.Table("rows", []string{"story", "seam", "red_signal"}, Rows(p))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	b.WriteString(tbl)
	return b.String(), 0
}
