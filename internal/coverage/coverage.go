// Package coverage ports `bench coverage`: the acceptance-coverage-map parser the
// gate's docs layer and the review phase both consume. Extraction mode emits the
// spec's state and rows as TOON; `--check` validates the map (canonical header,
// five-cell rows, non-empty cells, story references against the exact declared
// story set, historical opt-out) and requires a map at all unless the spec is
// marked historical. A map may opt into per-row IDs by leading the header with a
// `row` column; an opted-in map's IDs are grammar-checked, spec-local unique, and
// exported to other packages via ParseSpec. The validation phrasings are
// load-bearing — downstream consumers match them by substring — so this is the
// one validator for the convention.
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

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
// MinArgs stays 0 so the absent-spec case keeps its own named message below rather than
// the generic missing-positional one.
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
	// notCoveredRe is the one grammar for a story-level coverage exception: a story no
	// row references must either gain a row or carry this line with a reason.
	notCoveredRe = regexp.MustCompile(`^Not covered: story ([0-9]+) — (.*)$`)
	// clauseRe finds a `;` outside backticks: the cheap tell of a behavior cell that
	// states more than one predicate, which no single test can go red on.
	clauseRe = regexp.MustCompile("`[^`]*`|;")
	parenRe  = regexp.MustCompile(`[ \t]*\(.*\)$`)
	// storyPartRe matches one comma-separated part of a story reference; storyRefRe
	// matches the whole comma list. Both compose storyPartPattern, so they cannot
	// disagree about what a part looks like: every trimmed part storyRefRe accepts is
	// one storyPartRe matches, which is why the submatch below needs no nil check.
	storyPartRe = regexp.MustCompile(`^` + storyPartPattern + `$`)
	storyRefRe  = regexp.MustCompile(`^` + storyPartPattern + `([ \t]*,[ \t]*` + storyPartPattern + `)*$`)
)

// historicalMarker is the literal opt-out comment: present anywhere in a spec, it
// exempts the spec from the coverage-map requirement (a no-map state) and from
// row validation (a mapped-but-historical state).
const historicalMarker = "<!-- coverage-map: historical -->"

var fieldNames5 = [5]string{"story", "behavior", "seam", "red signal", "why it catches the failure"}
var fieldNames6 = [6]string{"row", "story", "behavior", "seam", "red signal", "why it catches the failure"}

// rowIDPattern is the one row-ID grammar the map's leading `row` cell answers to:
// an uppercase tag plus a number.
const rowIDPattern = `[A-Z]+[0-9]+`

// rowIDRe anchors rowIDPattern to a whole cell, spec-local unique.
var rowIDRe = regexp.MustCompile(`^` + rowIDPattern + `$`)

type dataRow struct {
	ncells int
	all    [6]string // c[1..6] (c[1..5] when not opted in), trimmed; empty beyond ncells
}

// parsed is the result of one scan of a spec: whether the map header was seen, the
// historical opt-out, whether the parsed header matched the canonical one, whether
// that header opted into row IDs, the declared story numbers, and the data rows.
type parsed struct {
	seen       bool
	historical bool
	gotHeader  bool
	headerOK   bool
	optIn      bool
	storyNums  map[int]bool
	notCovered map[int]string
	dataRows   []dataRow
}

// width reports the cell count an opted-in map's rows must have (6) versus a legacy
// map's (5); fields and storyOffset shift in lockstep with the same choice.
func (p parsed) width() int {
	if p.optIn {
		return 6
	}
	return 5
}

// fields names the cells at width(), in order, for empty-cell and count violations.
func (p parsed) fields() []string {
	if p.optIn {
		return fieldNames6[:]
	}
	return fieldNames5[:]
}

// storyOffset is the index of the story cell: 0 for a legacy row, 1 when a leading
// row-ID cell has shifted every other field one place right.
func (p parsed) storyOffset() int {
	if p.optIn {
		return 1
	}
	return 0
}

func parse(content []byte) parsed {
	var p parsed
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
		switch strings.ToLower(strings.Join(cells, "|")) {
		case "story|behavior|seam|red signal|why it catches the failure":
			p.headerOK = true
		case "row|story|behavior|seam|red signal|why it catches the failure":
			p.headerOK = true
			p.optIn = true
		}
		return
	}
	if allSep {
		return
	}
	dr := dataRow{ncells: len(cells)}
	for i := 0; i < 6 && i < len(cells); i++ {
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
// The row-ID cell of an opted-in map, if any, is not part of this schema.
func Rows(p parsed) [][]string {
	if State(p) != "mapped" {
		return nil
	}
	off := p.storyOffset()
	var rows [][]string
	for _, r := range p.dataRows {
		rows = append(rows, []string{r.all[off], r.all[off+2], r.all[off+3]})
	}
	return rows
}

// rowIDs returns the leading row-ID cell of each data row, in map order, for an
// opted-in mapped spec; nil for a legacy map or a spec that is not mapped. It backs
// ParseSpec, the package's one exported entry point for callers outside the
// package, which cannot construct a parsed value themselves.
func rowIDs(p parsed) []string {
	if State(p) != "mapped" || !p.optIn {
		return nil
	}
	ids := make([]string, len(p.dataRows))
	for i, r := range p.dataRows {
		ids[i] = r.all[0]
	}
	return ids
}

// Check returns one violation message per problem. A historical spec (mapped or
// not) is silent (nil); an unmapped spec with no historical marker is a violation
// in itself — an unmapped spec cannot pass the gate's docs layer by having
// nothing to validate. The phrasings are matched by substring downstream.
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
	if !p.headerOK {
		return []string{"coverage map missing the canonical header"}
	}
	if len(p.dataRows) == 0 {
		return []string{"coverage map has no data rows"}
	}
	width, fields, storyOff := p.width(), p.fields(), p.storyOffset()
	seenIDs := make(map[string]int) // row id -> first row number that used it
	referenced := make(map[int]bool)
	var v []string
	for idx, r := range p.dataRows {
		rn := idx + 1
		if r.ncells != width {
			v = append(v, fmt.Sprintf("coverage map row %d has %d cells (want %d)", rn, r.ncells, width))
			continue
		}
		for i := 0; i < width; i++ {
			if r.all[i] == "" {
				v = append(v, fmt.Sprintf("coverage map row %d has an empty '%s' cell", rn, fields[i]))
			}
		}
		if p.optIn && r.all[0] != "" {
			id := r.all[0]
			if !rowIDRe.MatchString(id) {
				v = append(v, fmt.Sprintf("coverage map row %d has a malformed row id '%s'", rn, id))
			} else if first, dup := seenIDs[id]; dup {
				v = append(v, fmt.Sprintf("coverage map row %d has a duplicate row id '%s' (first used at row %d)", rn, id, first))
			} else {
				seenIDs[id] = rn
			}
		}
		if strings.Contains(stripBackticks(r.all[storyOff+1]), ";") {
			v = append(v, fmt.Sprintf("coverage map row %d behavior states more than one predicate (';' outside backticks); split the row", rn))
		}
		story := parenRe.ReplaceAllString(r.all[storyOff], "")
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
			// A range stands for every number it spans, not just its endpoints: a spec
			// declaring 1, 2, 4 has both ends of `2-4` but no story 3 for the row to
			// cover. Only the first gap is reported — it names the row's fault, and a
			// wide range would otherwise bury the rest of the output.
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
// spec's exact declared set. Each failure mode names itself: 0 is never a valid
// story number, a spec that declares no stories at all says so plainly, and any
// other non-member is reported against the declared set rather than a maximum —
// a spec that skips a number (1, 2, 4) makes "numbers only N" false for a row
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

// ParseSpec reads and parses the coverage map at path, the package's one exported
// entry point for callers that cannot construct the unexported parsed type
// themselves: internal/preflight reads a spec's opt-in verdict, ordered row IDs, and
// Check violations through here rather than re-deriving map structure. ids is nil
// when the map is not opted into row IDs (legacy or absent).
func ParseSpec(path string) (optIn bool, ids []string, violations []string, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, nil, nil, err
	}
	p := parse(content)
	return p.optIn, rowIDs(p), Check(p), nil
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
		violations := Check(p)
		if len(violations) == 0 {
			// A pass is always a definitive one-line result, never silence, so "nothing
			// to validate" is never mistaken for "nothing printed by accident". Check is
			// silent for two states — a mapped spec with no violations, and a historical
			// marker (mapped or not) — and a third silent state added there needs its own
			// pass line here.
			if State(p) == "mapped" {
				return checkOKLine(fmt.Sprintf("coverage map valid — %d row(s)", len(p.dataRows))), 0
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
	tbl, err := toon.Table("rows", []string{"story", "seam", "red_signal"}, Rows(p))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	b.WriteString(tbl)
	violations := Check(p)
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
