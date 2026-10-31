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
	"github.com/gibbonmi/bench/internal/tickets"
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

// storyPartPattern is the one grammar for a story-reference part: a number, or a
// range joined by an en-dash or hyphen, capturing its endpoints.
const storyPartPattern = `([0-9]+)(?:[ \t]*(?:–|-)[ \t]*([0-9]+))?`

var (
	edgeRe = regexp.MustCompile(`^[Ee][Dd][Gg][Ee]`)
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

// rowIDRe anchors the shared row-ID grammar to a whole cell, spec-local unique. Its
// one submatch is the row's alphabetic tag, so a caller reads the tag from the same
// match that decides whether the cell is well-formed.
var rowIDRe = regexp.MustCompile(`^` + tickets.RowIDPattern + `$`)

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
		story, parts, ok := storyReference(s.cell(r, fieldStory))
		if !ok {
			v = append(v, fmt.Sprintf("coverage map row %d has an unrecognized story reference '%s'", rn, story))
			continue
		}
		fanOut := 0
		for _, part := range parts {
			start, end := part.start, part.end
			if !part.ranged {
				v = append(v, p.checkStoryMember(rn, start)...)
				referenced[start] = true
				fanOut++
				continue
			}
			if end < start {
				v = append(v, fmt.Sprintf("coverage map row %d has a story range with end before start '%s-%s'", rn, part.startText, part.endText))
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
