package maps

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// DecisionMapCandidate identifies a directly discovered map and its ownership.
type DecisionMapCandidate struct {
	Path     string
	Compiled bool
}

func isDirectoryDoc(name string) bool {
	return strings.EqualFold(strings.TrimSuffix(name, filepath.Ext(name)), "README")
}

// discoverDirectoryCandidates is the sole direct-child candidate policy for active
// and compiled decision-map directories. Callers choose which directories to compose.
func discoverDirectoryCandidates(root, dir string, compiled bool) ([]DecisionMapCandidate, bounds.FileState, string) {
	classified := bounds.ClassifyDir(dir)
	if classified.State != bounds.StateParsed {
		return nil, classified.State, classified.Reason
	}
	var candidates []DecisionMapCandidate
	for _, entry := range classified.Entries {
		name := entry.Name()
		if entry.IsDir() || strings.HasPrefix(name, ".") || !strings.HasSuffix(name, ".md") || isDirectoryDoc(name) {
			continue
		}
		path, err := filepath.Rel(root, filepath.Join(dir, name))
		if err != nil {
			return nil, bounds.StateUnreadable, err.Error()
		}
		candidates = append(candidates, DecisionMapCandidate{Path: filepath.ToSlash(path), Compiled: compiled})
	}
	return candidates, bounds.StateParsed, ""
}

// DiscoverDecisionMapCandidates finds active and compiled direct Markdown candidates.
func DiscoverDecisionMapCandidates(root string) ([]DecisionMapCandidate, error) {
	candidates, diagnostics := discoverDecisionMapCandidates(root)
	if len(diagnostics) > 0 {
		return nil, fmt.Errorf("%s", diagnostics[0])
	}
	return candidates, nil
}

// discoverDecisionMapCandidates is the sole active-and-compiled candidate traversal.
// Callers receive every readable candidate plus any independent directory diagnostics.
func discoverDecisionMapCandidates(root string) ([]DecisionMapCandidate, []string) {
	var candidates []DecisionMapCandidate
	var diagnostics []string
	appendDirectory := func(dir string, compiled bool) {
		discovered, state, reason := discoverDirectoryCandidates(root, dir, compiled)
		if state == bounds.StateAbsent {
			return
		}
		if state != bounds.StateParsed && state != bounds.StateEmpty {
			rel, err := filepath.Rel(root, dir)
			if err != nil {
				rel = dir
			}
			diagnostics = append(diagnostics, fmt.Sprintf("%s: %s: %s", filepath.ToSlash(rel), state, reason))
			return
		}
		candidates = append(candidates, discovered...)
	}
	appendDirectory(filepath.Join(root, DecisionsDir), false)
	specs := filepath.Join(root, "specs")
	classified := bounds.ClassifyDir(specs)
	if classified.State == bounds.StateAbsent {
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Path < candidates[j].Path })
		return candidates, diagnostics
	}
	if classified.State != bounds.StateParsed && classified.State != bounds.StateEmpty {
		diagnostics = append(diagnostics, fmt.Sprintf("specs: %s: %s", classified.State, classified.Reason))
		return candidates, diagnostics
	}
	for _, spec := range classified.Entries {
		if !spec.IsDir() || strings.HasPrefix(spec.Name(), ".") {
			continue
		}
		appendDirectory(filepath.Join(specs, spec.Name(), DecisionsDir), true)
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Path < candidates[j].Path })
	return candidates, diagnostics
}

type field struct {
	name    string
	syntax  string
	heading bool
	scoped  bool
}

type terminalSection struct {
	heading string
	syntax  string
}

type decisionMapSchema struct {
	statuses            []string
	types               []string
	fields              []field
	terminalSections    []terminalSection
	ticketHeading       string
	unsupportedHeadings []string
}

var canonicalDecisionMapSchema = decisionMapSchema{
	statuses: []string{"shaping", "ready"},
	types:    []string{"Research", "Prototype", "Grill", "Task"},
	fields: []field{
		{name: "title", syntax: "# "},
		{name: "Status", syntax: "Status: "},
		{name: "Destination", syntax: "## Destination", heading: true},
		{name: "Blocked by", syntax: "Blocked by: ", scoped: true},
		{name: "Type", syntax: "Type: ", scoped: true},
		{name: "Question", syntax: "### Question", heading: true, scoped: true},
		{name: "Answer", syntax: "### Answer", heading: true, scoped: true},
	},
	terminalSections: []terminalSection{
		{heading: "Not yet specified", syntax: "## Not yet specified"},
		{heading: "Spec-writer discretion", syntax: "## Spec-writer discretion"},
		{heading: "Out of scope", syntax: "## Out of scope"},
		{heading: "Sources", syntax: "## Sources"},
	},
	ticketHeading:       "## #",
	unsupportedHeadings: []string{"## Handoff"},
}

// resolvedBlockedRule states, in the rendered template, the graph rule the walk in
// graphDiagnostics already enforces. The template states the rule; it adds no check.
const resolvedBlockedRule = "A resolved decision ticket cannot stay blocked by an unresolved ticket."

// decisionMapSourcesExample teaches the Sources record grammar in the rendered
// template. The locator is a URL, because validateSourcePath resolves a Path locator
// against the repository root, where no placeholder file exists. Supports and Drift
// stay on two physical lines, because a Sources record keeps each field on one line.
const decisionMapSourcesExample = "\n- URL: https://example.invalid/decision-source\n" +
	"  Supports: <the decision this source supports>\n" +
	"  Drift: <the change that makes this source stale>\n"

// decisionMapAssetRule states, in the rendered template, where a map-owned asset stays.
// The candidate scanner lists only the direct children of a decisions directory, so an
// asset in this nested directory is never read as a decision map. The template states the
// convention; it adds no check.
const decisionMapAssetRule = "A map-owned asset stays in `decisions/assets/`."

var decisionHeading = regexp.MustCompile(`^([1-9][0-9]*):\s+(.+?)\s*$`)
var blockersField = regexp.MustCompile(`^(none|#[1-9][0-9]*(, #[1-9][0-9]*)*)$`)

// Diagnostic describes one structural decision-map problem.
type Diagnostic struct{ Message string }

// DecisionTicket is one parsed decision ticket.
type DecisionTicket struct {
	ID        string
	Title     string
	BlockedBy string
	Type      string
	Question  string
	Answer    string
}

// DecisionMap is the parsed, schema-owned form of a decision map.
type DecisionMap struct {
	Title       string
	Status      string
	Destination string
	Tickets     []DecisionTicket
	Fog         string
	Discretion  string
	OutOfScope  string
	Sources     string
}

func (s decisionMapSchema) hasStatus(status string) bool {
	for _, candidate := range s.statuses {
		if status == candidate {
			return true
		}
	}
	return false
}

func (s decisionMapSchema) hasType(typ string) bool {
	for _, candidate := range s.types {
		if typ == candidate {
			return true
		}
	}
	return false
}

func (s decisionMapSchema) field(name string) field {
	for _, field := range s.fields {
		if field.name == name {
			return field
		}
	}
	panic("unknown decision-map schema field: " + name)
}

func (s decisionMapSchema) terminalHeading(line string) string {
	for _, section := range s.terminalSections {
		if line == section.syntax {
			return section.heading
		}
	}
	return ""
}

func (s decisionMapSchema) ticket(line string) (id, title string, ok bool) {
	match := decisionHeading.FindStringSubmatch(strings.TrimPrefix(line, s.ticketHeading))
	if !strings.HasPrefix(line, s.ticketHeading) || match == nil {
		return "", "", false
	}
	return match[1], match[2], true
}

// fieldScan drives the shared field scan with the decision-map field table. A
// decision ticket opens one scope, and a section heading closes it.
func (s decisionMapSchema) fieldScan() FieldScan {
	table := make([]FieldSpec, 0, len(s.fields)+len(s.terminalSections))
	for _, f := range s.fields {
		table = append(table, FieldSpec{Name: f.name, Syntax: f.syntax, Heading: f.heading, Scoped: f.scoped})
	}
	for _, terminal := range s.terminalSections {
		table = append(table, FieldSpec{Name: terminal.heading, Syntax: terminal.syntax, Heading: true})
	}
	return FieldScan{
		Table: table,
		Scope: func(line string) (string, bool) {
			if id, _, ok := s.ticket(line); ok {
				return id, true
			}
			if line == s.field("Destination").syntax || s.terminalHeading(line) != "" {
				return "", true
			}
			for _, heading := range s.unsupportedHeadings {
				if line == heading {
					return "", true
				}
			}
			return "", false
		},
		Duplicate: func(spec FieldSpec, scope string) string {
			switch {
			case spec.Scoped:
				return fmt.Sprintf("ticket #%s: duplicate %s", scope, spec.Name)
			case spec.Name == "title", spec.Name == "Status":
				return "duplicate " + spec.Name
			default:
				return "duplicate " + spec.Name + " section"
			}
		},
	}
}

// ParseDecisionMap parses a decision map according to the canonical schema.
func ParseDecisionMap(content []byte) (DecisionMap, []Diagnostic) {
	var m DecisionMap
	var diagnostics []Diagnostic
	seenTerminal := map[string]bool{}
	var current *DecisionTicket
	answerSeen := false
	section := ""

	finishTicket := func() {
		if current == nil {
			return
		}
		if current.BlockedBy == "" {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: missing Blocked by", current.ID)})
		} else if !blockersField.MatchString(current.BlockedBy) {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: malformed Blocked by", current.ID)})
		}
		if current.Type == "" {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: missing Type", current.ID)})
		} else if !canonicalDecisionMapSchema.hasType(current.Type) {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: unsupported Type %q", current.ID, current.Type)})
		}
		if current.Question == "" {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: missing Question", current.ID)})
		}
		if !answerSeen {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: missing Answer", current.ID)})
		}
		m.Tickets = append(m.Tickets, *current)
		current = nil
	}

	scanned, _ := canonicalDecisionMapSchema.fieldScan().Scan(content)
	for _, entry := range scanned {
		if entry.Fenced {
			continue
		}
		line := entry.Text
		duplicate := entry.Diagnostic != ""
		report := func() { diagnostics = append(diagnostics, Diagnostic{Message: entry.Diagnostic}) }
		switch entry.Field {
		case "title":
			if duplicate {
				report()
			} else {
				m.Title = entry.Value
			}
			continue
		case "Status":
			if duplicate {
				report()
			} else {
				m.Status = entry.Value
			}
			continue
		}
		unsupported := false
		for _, heading := range canonicalDecisionMapSchema.unsupportedHeadings {
			if line == heading {
				finishTicket()
				diagnostics = append(diagnostics, Diagnostic{Message: "unsupported Handoff section"})
				section = "Handoff"
				unsupported = true
				break
			}
		}
		if unsupported {
			continue
		}
		if heading := canonicalDecisionMapSchema.terminalHeading(line); heading != "" {
			finishTicket()
			if duplicate {
				report()
			}
			seenTerminal[heading] = true
			section = heading
			continue
		}
		if entry.Field == "Destination" {
			finishTicket()
			if duplicate {
				report()
			}
			section = "Destination"
			continue
		}
		if id, title, ok := canonicalDecisionMapSchema.ticket(line); ok {
			finishTicket()
			current = &DecisionTicket{ID: id, Title: title}
			answerSeen = false
			section = "ticket"
			continue
		}
		if current != nil {
			switch entry.Field {
			case "Blocked by":
				if duplicate {
					report()
					break
				}
				current.BlockedBy = entry.Value
			case "Type":
				if duplicate {
					report()
					break
				}
				current.Type = entry.Value
			case "Question":
				if duplicate {
					report()
				}
				section = "Question"
			case "Answer":
				if duplicate {
					report()
				}
				answerSeen = true
				section = "Answer"
			default:
				if strings.TrimSpace(line) == "" {
					break
				}
				if section == "Question" {
					current.Question = appendSectionLine(current.Question, strings.TrimSpace(line))
				} else if section == "Answer" {
					current.Answer = appendSectionLine(current.Answer, strings.TrimSpace(line))
				}
			}
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		switch section {
		case "Destination":
			m.Destination = strings.TrimSpace(line)
		case "Not yet specified":
			m.Fog = appendSectionLine(m.Fog, line)
		case "Spec-writer discretion":
			m.Discretion = appendSectionLine(m.Discretion, line)
		case "Out of scope":
			m.OutOfScope = appendSectionLine(m.OutOfScope, line)
		case "Sources":
			m.Sources = appendSectionLine(m.Sources, line)
		}
	}
	finishTicket()
	if m.Title == "" {
		diagnostics = append(diagnostics, Diagnostic{Message: "missing title"})
	}
	if m.Status == "" {
		diagnostics = append(diagnostics, Diagnostic{Message: "missing Status"})
	} else if !canonicalDecisionMapSchema.hasStatus(m.Status) {
		diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("unsupported Status %q", m.Status)})
	}
	if m.Destination == "" {
		diagnostics = append(diagnostics, Diagnostic{Message: "missing Destination"})
	}
	if len(m.Tickets) == 0 {
		diagnostics = append(diagnostics, Diagnostic{Message: "missing decision ticket"})
	}
	for _, terminal := range canonicalDecisionMapSchema.terminalSections {
		if !seenTerminal[terminal.heading] {
			diagnostics = append(diagnostics, Diagnostic{Message: "missing " + terminal.heading + " section"})
		}
	}
	return m, diagnostics
}

func appendSectionLine(section, line string) string {
	if section == "" {
		return line
	}
	return section + "\n" + line
}

// DecisionMapTemplate renders the canonical decision-map Markdown skeleton.
func DecisionMapTemplate() string {
	var b strings.Builder
	b.WriteString(canonicalDecisionMapSchema.field("title").syntax)
	b.WriteString("<decision map title>\n\n")
	b.WriteString(canonicalDecisionMapSchema.field("Status").syntax)
	b.WriteString(canonicalDecisionMapSchema.statuses[0])
	b.WriteString("\n\n")
	b.WriteString(canonicalDecisionMapSchema.field("Destination").syntax)
	b.WriteString("\n\n<what this map decides>\n\n")
	b.WriteString(canonicalDecisionMapSchema.ticketHeading)
	b.WriteString("1: <decision question>\n\n")
	b.WriteString(canonicalDecisionMapSchema.field("Blocked by").syntax)
	b.WriteString("none\n")
	b.WriteString(canonicalDecisionMapSchema.field("Type").syntax)
	b.WriteString(canonicalDecisionMapSchema.types[0])
	b.WriteString("\n\n")
	b.WriteString(resolvedBlockedRule)
	b.WriteString("\n")
	b.WriteString(decisionMapAssetRule)
	b.WriteString("\n")
	for _, field := range []field{canonicalDecisionMapSchema.field("Question"), canonicalDecisionMapSchema.field("Answer")} {
		b.WriteString("\n")
		b.WriteString(field.syntax)
		b.WriteString("\n\n<")
		b.WriteString(strings.ToLower(field.name))
		b.WriteString(">\n")
	}
	for _, terminal := range canonicalDecisionMapSchema.terminalSections {
		b.WriteString("\n")
		b.WriteString(terminal.syntax)
		b.WriteString("\n")
		if terminal.heading == "Sources" {
			b.WriteString(decisionMapSourcesExample)
		}
	}
	return b.String()
}
