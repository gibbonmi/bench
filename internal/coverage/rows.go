package coverage

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/markdown"
	specref "github.com/gibbonmi/bench/internal/spec"
)

const escSentinel = "\x1c" // stands in for an escaped `\|` across the pipe split

var (
	storyNumRe = regexp.MustCompile(`^[0-9]+\. `)
	mapEndRe   = regexp.MustCompile(`^#{2,} `)
	// notCoveredRe is the one grammar for a story-level coverage exception. A story no
	// row references must either gain a row or carry this line with a reason.
	notCoveredRe = regexp.MustCompile(`^Not covered: story ([0-9]+) — (.*)$`)
)

// historicalMarker is the literal opt-out comment. It may appear anywhere in a spec.
// It exempts the spec from the coverage-map requirement, a no-map state, and from
// row validation, a mapped-but-historical state.
const historicalMarker = "<!-- coverage-map: historical -->"

// The cell names a coverage-map header can carry. A schema is a list of these in
// cell order, and a header line is their join. Each name is spelled once, so the
// header match and the violation messages read it from the same place.
const (
	fieldRow      = "row"
	fieldStory    = "story"
	fieldBehavior = "behavior"
	fieldSeam     = "seam"
	fieldWhy      = "why it catches the failure"
)

// schema is the one descriptor for an accepted header: its field names, in cell
// order. Row width, the name a violation message quotes, and every cell offset all
// derive from that list. Each check reads its cell by name, so no check can address
// the wrong column when a second header joins the set. To add a header, add a
// descriptor here instead of editing the checks.
type schema struct {
	fields []string
}

// schemas are the accepted headers, tried in order. The first is also the projection
// fallback for a map whose header matched none of them. The non-opt-in header stays
// first, because projection() falls back to schemas[0], and that fallback's output is
// pinned.
var schemas = []schema{
	{fields: []string{fieldStory, fieldBehavior, fieldSeam, fieldWhy}},
	{fields: []string{fieldRow, fieldStory, fieldBehavior, fieldSeam, fieldWhy}},
}

// header is the lowercased header line this schema answers to: the join of its own
// field names. The match and the descriptor cannot drift apart.
func (s schema) header() string { return strings.Join(s.fields, "|") }

// known reports whether a header resolved to a descriptor at all.
func (s schema) known() bool { return len(s.fields) > 0 }

// width is the cell count this schema's data rows must have.
func (s schema) width() int { return len(s.fields) }

// optIn reports whether this schema carries per-row IDs, which is to say a row field.
func (s schema) optIn() bool { return s.index(fieldRow) >= 0 }

// index is the cell offset of a named field, or -1 when this schema has no such field.
func (s schema) index(name string) int {
	for i, f := range s.fields {
		if f == name {
			return i
		}
	}
	return -1
}

// cell reads one row's named field. A field this schema does not carry, and a field
// past a short row's last cell, both read as empty. That matches a present-but-blank
// cell, which the empty-cell and row-ID checks already treat as a fault.
func (s schema) cell(r dataRow, name string) string {
	i := s.index(name)
	if i < 0 || i >= len(r.cells) {
		return ""
	}
	return r.cells[i]
}

// schemaFor resolves a lowercased header line to its descriptor. A header matching
// none yields the zero schema, which reports !known().
func schemaFor(header string) schema {
	for _, s := range schemas {
		if s.header() == header {
			return s
		}
	}
	return schema{}
}

type dataRow struct {
	cells []string      // one map row's cells, trimmed, exactly as many as were written
	line  int           // the 1-based spec line the row sits on
	span  markdown.Span // the row's original bytes, without the line terminator
}

// storyLine is one numbered story line of the User stories section, in authored order. A
// repeated number keeps each occurrence, so a selection can refuse it as ambiguous.
type storyLine struct {
	number int
	span   markdown.Span // the story's original bytes, without the line terminator
}

// parsed is the result of one scan of a spec. It records whether the scan saw the
// map header, the historical opt-out, and whether it reached a header line at all.
// It also records the schema that header resolved to, the declared story numbers,
// and the data rows.
type parsed struct {
	seen       bool
	historical bool
	gotHeader  bool
	sch        schema // the descriptor the header matched; zero when none did
	storyNums  map[int]bool
	stories    []storyLine
	notCovered map[int]string
	dataRows   []dataRow
	// mapSpan is the original bytes of the map: its heading line through the byte before the
	// heading that ends it, or through the end of the content.
	mapSpan markdown.Span
	// fault is the block reader's fault, which leaves the rest of the content unread.
	fault *markdown.Fault
}

// projection is the descriptor Rows reads cells through. A header matching no
// descriptor has none of its own, so it projects through schemas[0], the four-cell
// reduced schema. In that schema, story, behavior, and seam stay at offsets 0, 1,
// and 2. Check refuses such a map before any other cell read.
func (p parsed) projection() schema {
	if p.sch.known() {
		return p.sch
	}
	return schemas[0]
}

func parse(content []byte) parsed {
	var p parsed
	inMap := false
	doc := markdown.Read(string(content))
	p.fault = doc.Fault
	for _, line := range doc.Lines {
		text := line.Raw
		if text == historicalMarker {
			p.historical = true
		}
		if line.Class != markdown.Body {
			continue
		}
		raw := markdown.Span{Start: line.Offset, End: line.Offset + len(text)}
		inStories := line.Section == specref.StoriesSection
		if inStories && line.Heading {
			continue
		}
		if inStories && storyNumRe.MatchString(text) {
			num := text
			if i := strings.IndexByte(num, '.'); i >= 0 {
				num = num[:i]
			}
			if n, err := strconv.Atoi(num); err == nil {
				if p.storyNums == nil {
					p.storyNums = make(map[int]bool)
				}
				p.storyNums[n] = true
				p.stories = append(p.stories, storyLine{number: n, span: raw})
			}
		}
		if m := notCoveredRe.FindStringSubmatch(text); m != nil {
			if n, err := strconv.Atoi(m[1]); err == nil {
				if p.notCovered == nil {
					p.notCovered = make(map[int]string)
				}
				p.notCovered[n] = trimSpaceTab(m[2])
			}
		}
		if text == "### Acceptance coverage map" && !p.seen {
			p.seen = true
			inMap = true
			p.mapSpan = markdown.Span{Start: line.Offset, End: len(content)}
			continue
		}
		if inMap && mapEndRe.MatchString(text) {
			inMap, p.mapSpan.End = false, line.Offset
		}
		if inMap {
			p.processMapLine(text, line.Number, raw)
		}
	}
	return p
}

func (p *parsed) processMapLine(raw string, lineNum int, span markdown.Span) {
	cells, ok := pipeCells(raw)
	if !ok {
		return
	}
	if !p.gotHeader {
		p.gotHeader = true
		p.sch = schemaFor(strings.ToLower(strings.Join(cells, "|")))
		return
	}
	if separatorRow(cells) {
		return
	}
	p.dataRows = append(p.dataRows, dataRow{cells: cells, line: lineNum, span: span})
}

// pipeCells is the one pipe-table cell splitter: it splits one Markdown table line into
// its trimmed cells, and an escaped `\|` stays inside its cell. A line that does not open
// with a pipe is no table line, so ok is false.
func pipeCells(raw string) (cells []string, ok bool) {
	line := trimSpaceTab(raw)
	if !strings.HasPrefix(line, "|") {
		return nil, false
	}
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	line = strings.ReplaceAll(line, `\|`, escSentinel)
	cells = strings.Split(line, "|")
	for i, c := range cells {
		cells[i] = trimSpaceTab(strings.ReplaceAll(c, escSentinel, `\|`))
	}
	return cells, true
}

// separatorRow reports whether every cell of a table line is dashes only, which is the
// header separator.
func separatorRow(cells []string) bool {
	for _, c := range cells {
		if !isDashes(c) {
			return false
		}
	}
	return true
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

// Rows returns the story, behavior, and seam cells of each data row for a mapped
// spec, or nil otherwise. The row's schema resolves the cells by name. A map that
// opts into row IDs projects the same three fields as a non-opt-in one, and the
// leading row-ID cell is not part of the projection. Every accepted header carries
// these three fields — behavior is the cell that names what to build — so a caller
// reads one row shape whichever schema the spec uses.
func Rows(p parsed) [][]string {
	if State(p) != "mapped" {
		return nil
	}
	s := p.projection()
	var rows [][]string
	for _, r := range p.dataRows {
		rows = append(rows, []string{s.cell(r, fieldStory), s.cell(r, fieldBehavior), s.cell(r, fieldSeam)})
	}
	return rows
}

// rowIDs returns the leading row-ID cell of each data row, in map order, for an
// opted-in mapped spec. It returns nil for a non-opt-in map or a spec that is not
// mapped. It backs ParseSpec, the package's one exported entry point for callers
// outside the package. Those callers cannot construct a parsed value themselves.
func rowIDs(p parsed) []string {
	if State(p) != "mapped" || !p.sch.optIn() {
		return nil
	}
	ids := make([]string, len(p.dataRows))
	for i, r := range p.dataRows {
		ids[i] = p.sch.cell(r, fieldRow)
	}
	return ids
}

// MapRow is one data row of a spec's coverage map: its row ID when the map opts into row
// IDs, its story, behavior, and seam cells, and the 1-based spec line it sits on.
type MapRow struct {
	ID, Story, Behavior, Seam string
	Line                      int
}

// MapRows projects the coverage map in content into typed rows, in map order. A spec that
// is not mapped has none. It reads the one parse that Rows and rowIDs read, so a caller
// outside the package never splits the pipe table again.
func MapRows(content []byte) []MapRow {
	p := parse(content)
	if State(p) != "mapped" {
		return nil
	}
	s := p.projection()
	rows := make([]MapRow, len(p.dataRows))
	for i, r := range p.dataRows {
		rows[i] = MapRow{ID: s.cell(r, fieldRow), Story: s.cell(r, fieldStory), Behavior: s.cell(r, fieldBehavior), Seam: s.cell(r, fieldSeam), Line: r.line}
	}
	return rows
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
