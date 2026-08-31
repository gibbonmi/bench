// Package coverage ports `bench coverage`, the acceptance-coverage-map parser the
// gate's docs layer and the review phase both consume. Extraction mode emits the
// spec's state and rows as TOON. `--check` validates the map: a canonical header,
// rows as wide as that header declares, non-empty cells, story references against
// the exact declared story set, and a historical opt-out. It requires a map at all
// unless the spec is marked historical. The header itself fixes which columns it
// carries, never the cell count.
//
// A map may opt into per-row IDs by leading the header with a `row` column. An
// opted-in map's IDs are grammar-checked and spec-local unique, and ParseSpec exports
// them to other packages. The validation phrasings are load-bearing, because
// downstream consumers match them by substring, so this is the one validator for the
// convention.
package coverage

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
	specref "github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand.
// Arity, flag recognition, `--`, and help all come from there rather than a local
// switch. MinArgs stays 0, so the absent-spec case keeps its own named message below
// rather than the generic missing-positional one.
var grammar = usage.Grammar{
	Cmd:     "bench coverage",
	Help:    "usage: bench coverage [--check] <spec.md | slug>",
	Flags:   []usage.Flag{{Name: "--check"}},
	MaxArgs: 1,
}

const escSentinel = "\x1c" // stands in for an escaped `\|` across the pipe split

// storyPartPattern is the one grammar for a story-reference part: a number, or a
// range joined by an en-dash or hyphen, capturing its endpoints.
const storyPartPattern = `([0-9]+)(?:[ \t]*(?:–|-)[ \t]*([0-9]+))?`

var (
	storyNumRe = regexp.MustCompile(`^[0-9]+\. `)
	mapEndRe   = regexp.MustCompile(`^#{2,} `)
	edgeRe     = regexp.MustCompile(`^[Ee][Dd][Gg][Ee]`)
	// notCoveredRe is the one grammar for a story-level coverage exception. A story no
	// row references must either gain a row or carry this line with a reason.
	notCoveredRe = regexp.MustCompile(`^Not covered: story ([0-9]+) — (.*)$`)
	// clauseRe finds a `;` outside backticks. This is the cheap tell of a behavior
	// cell that states more than one predicate, which no single test can go red on.
	clauseRe = regexp.MustCompile("`[^`]*`|;")
	parenRe  = regexp.MustCompile(`[ \t]*\(.*\)$`)
	// storyPartRe matches one comma-separated part of a story reference. storyRefRe
	// matches the whole comma list. Both compose storyPartPattern, so they cannot
	// disagree about what a part looks like. Every trimmed part storyRefRe accepts is
	// one storyPartRe matches, which is why the submatch below needs no nil check.
	storyPartRe = regexp.MustCompile(`^` + storyPartPattern + `$`)
	storyRefRe  = regexp.MustCompile(`^` + storyPartPattern + `([ \t]*,[ \t]*` + storyPartPattern + `)*$`)
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

// rowIDPattern is the one row-ID grammar the map's leading `row` cell answers to:
// an uppercase tag plus a number. The group captures the tag.
const rowIDPattern = `([A-Z]+)[0-9]+`

// rowIDRe anchors rowIDPattern to a whole cell, spec-local unique. Its one
// submatch is the row's alphabetic tag, so a caller reads the tag from the same
// match that decides whether the cell is well-formed.
var rowIDRe = regexp.MustCompile(`^` + rowIDPattern + `$`)

type dataRow struct {
	cells []string // one map row's cells, trimmed, exactly as many as were written
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
	notCovered map[int]string
	dataRows   []dataRow
	// The spec's `## Ownership fences` tokens, read through internal/spec's one fence
	// parser. fencesDeclared separates a spec that opens the section from one that
	// never declares it; only the first is graded for the review pickup.
	fences         []string
	fencesDeclared bool
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
	p.fences, p.fencesDeclared = specref.FenceTokens(content)
	inStories, inMap := false, false
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSuffix(raw, "\r")
		if line == historicalMarker {
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
			if n, err := strconv.Atoi(num); err == nil {
				if p.storyNums == nil {
					p.storyNums = make(map[int]bool)
				}
				p.storyNums[n] = true
			}
		}
		if m := notCoveredRe.FindStringSubmatch(line); m != nil {
			if n, err := strconv.Atoi(m[1]); err == nil {
				if p.notCovered == nil {
					p.notCovered = make(map[int]string)
				}
				p.notCovered[n] = trimSpaceTab(m[2])
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
		p.sch = schemaFor(strings.ToLower(strings.Join(cells, "|")))
		return
	}
	if allSep {
		return
	}
	p.dataRows = append(p.dataRows, dataRow{cells: cells})
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

// Check returns one violation message per problem. A historical spec, mapped or not,
// stays silent and returns nil. An unmapped spec with no historical marker is a
// violation in itself, because it cannot pass the gate's docs layer by having
// nothing to validate. Downstream consumers match the phrasings by substring.
func Check(p parsed) []string {
	switch State(p) {
	case "no-map":
		if p.historical {
			return nil
		}
		return []string{"coverage map missing and spec is not marked historical"}
	case "historical":
		return nil
	}
	if !p.sch.known() {
		return []string{"coverage map missing the canonical header"}
	}
	if len(p.dataRows) == 0 {
		return []string{"coverage map has no data rows"}
	}
	s := p.sch
	seenIDs := make(map[string]int)   // row id -> first row number that used it
	seenTags := make(map[string]bool) // alphabetic row-id tags this map declares
	referenced := make(map[int]bool)
	var v []string
	for idx, r := range p.dataRows {
		rn := idx + 1
		if len(r.cells) != s.width() {
			v = append(v, fmt.Sprintf("coverage map row %d has %d cells (want %d)", rn, len(r.cells), s.width()))
			continue
		}
		for _, name := range s.fields {
			if s.cell(r, name) == "" {
				v = append(v, fmt.Sprintf("coverage map row %d has an empty '%s' cell", rn, name))
			}
		}
		if id := s.cell(r, fieldRow); id != "" {
			m := rowIDRe.FindStringSubmatch(id)
			if m == nil {
				v = append(v, fmt.Sprintf("coverage map row %d has a malformed row id '%s'", rn, id))
			} else if first, dup := seenIDs[id]; dup {
				v = append(v, fmt.Sprintf("coverage map row %d has a duplicate row id '%s' (first used at row %d)", rn, id, first))
			} else {
				seenIDs[id] = rn
			}
			// A malformed cell contributes no tag: its message already names the fault,
			// and a guessed tag would add a second, derived red for the same cell.
			if m != nil {
				seenTags[m[1]] = true
			}
		}
		if strings.Contains(stripBackticks(s.cell(r, fieldBehavior)), ";") {
			v = append(v, fmt.Sprintf("coverage map row %d behavior states more than one predicate (';' outside backticks); split the row", rn))
		}
		story := parenRe.ReplaceAllString(s.cell(r, fieldStory), "")
		if story == "" || edgeRe.MatchString(story) {
			continue
		}
		if !storyRefRe.MatchString(story) {
			v = append(v, fmt.Sprintf("coverage map row %d has an unrecognized story reference '%s'", rn, story))
			continue
		}
		fanOut := 0
		for _, part := range strings.Split(story, ",") {
			m := storyPartRe.FindStringSubmatch(trimSpaceTab(part))
			start, _ := strconv.Atoi(m[1])
			if m[2] == "" {
				v = append(v, p.checkStoryMember(rn, start)...)
				referenced[start] = true
				fanOut++
				continue
			}
			end, _ := strconv.Atoi(m[2])
			if end < start {
				v = append(v, fmt.Sprintf("coverage map row %d has a story range with end before start '%s-%s'", rn, m[1], m[2]))
				continue
			}
			// A range stands for every number it spans, not just its endpoints. A spec
			// declaring 1, 2, 4 has both ends of `2-4` but no story 3 for the row to
			// cover. The check reports only the first gap: it names the row's fault, and
			// a wide range would otherwise bury the rest of the output.
			for n := start; n <= end; n++ {
				referenced[n] = true
				fanOut++
				if msgs := p.checkStoryMember(rn, n); msgs != nil {
					v = append(v, msgs...)
					break
				}
			}
		}
		if fanOut > bounds.CoverageRowStories {
			v = append(v, fmt.Sprintf("coverage map row %d references %d stories (max %d); an outcome family is not one red-capable row", rn, fanOut, bounds.CoverageRowStories))
		}
	}
	// The preflight's membership check scopes to one tag, so a second tag hides its
	// rows from that check. The map declares one tag; the message names each tag it
	// found, sorted, so the reader sees which rows disagree.
	if len(seenTags) > 1 {
		tags := make([]string, 0, len(seenTags))
		for t := range seenTags {
			tags = append(tags, t)
		}
		sort.Strings(tags)
		v = append(v, fmt.Sprintf("coverage map row ids carry more than one tag (%s); a map declares one tag", strings.Join(tags, ", ")))
	}
	// A declared story no row references is a breadth-floor promise nothing checks:
	// it needs a row or an explicit, reasoned exception line.
	declared := make([]int, 0, len(p.storyNums))
	for n := range p.storyNums {
		declared = append(declared, n)
	}
	sort.Ints(declared)
	for _, n := range declared {
		if referenced[n] {
			continue
		}
		reason, ok := p.notCovered[n]
		switch {
		case !ok:
			v = append(v, fmt.Sprintf("coverage map leaves story %d unreferenced; add a row or a `Not covered: story %d — <reason>` line", n, n))
		case reason == "":
			v = append(v, fmt.Sprintf("story %d is marked Not covered without a reason", n))
		}
	}
	return v
}

// stripBackticks removes inline code spans so a `;` inside a command or literal does
// not read as a second predicate.
func stripBackticks(cell string) string {
	return clauseRe.ReplaceAllStringFunc(cell, func(m string) string {
		if m == ";" {
			return m
		}
		return ""
	})
}

// checkStoryMember validates one story number referenced by row rn against the
// spec's exact declared set. Each failure mode names itself. Story 0 is never a
// valid story number, and a spec that declares no stories at all says so plainly.
// Any other non-member reports against the declared set rather than a maximum. A
// spec that skips a number — 1, 2, 4 — makes "numbers only N" false for a row
// referencing 3.
func (p parsed) checkStoryMember(rn, n int) []string {
	if n == 0 {
		return []string{fmt.Sprintf("coverage map row %d references story 0, which is not a valid story number", rn)}
	}
	if p.storyNums[n] {
		return nil
	}
	if len(p.storyNums) == 0 {
		return []string{fmt.Sprintf("coverage map row %d references story %d, but the spec declares no stories", rn, n)}
	}
	return []string{fmt.Sprintf("coverage map row %d references story %d, which the spec does not declare (has: %s)", rn, n, p.declaredStoriesList())}
}

// declaredStoriesList renders the spec's declared story numbers in ascending order,
// as the "has: ..." clause of the non-member violation message.
func (p parsed) declaredStoriesList() string {
	nums := make([]int, 0, len(p.storyNums))
	for n := range p.storyNums {
		nums = append(nums, n)
	}
	sort.Ints(nums)
	strs := make([]string, len(nums))
	for i, n := range nums {
		strs[i] = strconv.Itoa(n)
	}
	return strings.Join(strs, ", ")
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

// checkOKLine is the one construction for a --check pass line: the `ok: ` AXI
// prefix over a message that names why the check passed, never silence.
func checkOKLine(msg string) string { return fmt.Sprintf("ok: %s\n", msg) }

// uncitedLine renders the informational report that rides beside the pass line: the
// count of mapped rows no citation backs, and their names. It is one line, so a green
// check stays a bounded read. A map with nothing to report renders nothing, which keeps
// the pass line the whole response in the common case.
//
// The report joins no violation list and changes no exit code. It tells a build which
// rows are not yet wired; it does not claim they are faults.
func uncitedLine(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return fmt.Sprintf("uncited: %d row(s) with no seam-cell citation — %s\n", len(names), strings.Join(names, ", "))
}

// ParseSpec reads and parses the coverage map at path. It is the package's one
// exported entry point for callers that cannot construct the unexported parsed type
// themselves: internal/preflight reads a spec's opt-in verdict, ordered row IDs, and
// Check violations through here rather than re-deriving map structure. ids is nil
// when the map is not opted into row IDs, whether non-opt-in or absent.
func ParseSpec(path string) (optIn bool, ids []string, violations []string, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, nil, nil, err
	}
	p := parse(content)
	return p.sch.optIn(), rowIDs(p), CheckFiles(p, path), nil
}

// Command implements `bench coverage [--check] <spec.md | slug>`.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, check := parsed.Flags["--check"]
	spec := ""
	if len(parsed.Positionals) == 1 {
		spec = parsed.Positionals[0]
	}
	if spec == "" {
		return toon.MissingArg("bench coverage", "<spec.md> is required") + "\n", 2
	}
	// The repo root anchors the bare-slug fallback, exactly as it does for the
	// internal/spec subcommands, so which spec a slug names never depends on where the
	// caller is standing. Rendering stays repo-relative for the same reason: an absolute
	// path would make the same spec print differently per checkout.
	base := specref.RepoBase()
	content, resolved, tried, ok, err := specref.Resolve(base, spec)
	if err != nil {
		return toon.Errorf("spec not readable: "+err.Error(), "fix the file's permissions or pass another spec") + "\n", 1
	}
	if !ok {
		return toon.Errorf("spec not found: "+strings.Join(tried, ", "), "pass a path to a spec markdown file") + "\n", 1
	}
	spec = filepath.ToSlash(specref.RelTo(base, resolved))
	p := parse(content)
	if check {
		violations := CheckFiles(p, resolved)
		if len(violations) == 0 {
			// A pass is always a definitive one-line result, never silence, so "nothing to
			// validate" is never mistaken for "nothing printed by accident." Check stays
			// silent for two states: a mapped spec with no violations, and a historical
			// marker, mapped or not. A third silent state added there needs its own pass
			// line here.
			if State(p) == "mapped" {
				return checkOKLine(fmt.Sprintf("coverage map valid — %d row(s)", len(p.dataRows))) + uncitedLine(uncitedRows(p)), 0
			}
			return checkOKLine("coverage map historical — validation skipped"), 0
		}
		var b strings.Builder
		for _, viol := range violations {
			fmt.Fprintln(&b, toon.Errorf(spec+" "+viol, "fix the map or mark it "+historicalMarker))
		}
		return b.String(), 1
	}
	var b strings.Builder
	fmt.Fprintf(&b, "spec: %s\n", spec)
	fmt.Fprintf(&b, "state: %s\n", State(p))
	tbl, err := toon.Table("rows", []string{"story", "behavior", "seam"}, Rows(p))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	b.WriteString(tbl)
	violations := CheckFiles(p, resolved)
	rows := Rows(p)
	actions := make([]axi.Action, 0, len(rows))
	if len(rows) > 0 && len(violations) > 0 {
		actions = append(actions, axi.ExecutableInvocation("retry after repairing coverage map", axi.KnownArgument("coverage"), axi.KnownArgument("--check"), axi.KnownArgument(spec)))
	} else if len(rows) > 0 {
		for _, row := range rows {
			actions = append(actions, axi.ExecutableInvocation("check coverage for stories "+row[0], axi.KnownArgument("coverage"), axi.KnownArgument("--check"), axi.KnownArgument(spec)))
		}
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	b.WriteString(help)
	return b.String(), 0
}
