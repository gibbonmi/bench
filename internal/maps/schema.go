package maps

import (
	"fmt"
	"regexp"
	"strings"
)

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
	indexSections       []terminalSection
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
	// A map index carries these two sections. They are graded by indexDiagnostics
	// against the ticket files, so they are not terminal sections.
	indexSections: []terminalSection{
		{heading: "Notes", syntax: "## Notes"},
		{heading: "Decisions so far", syntax: "## Decisions so far"},
	},
	unsupportedHeadings: []string{"## Handoff"},
}

// decisionMapSourcesExample teaches the Sources record grammar in the rendered
// template. The locator is a URL, because validateSourcePath resolves a Path locator
// against the repository root, where no placeholder file exists. Supports and Drift
// stay on two physical lines, because a Sources record keeps each field on one line.
const decisionMapSourcesExample = "\n- URL: https://example.invalid/decision-source\n" +
	"  Supports: <the decision this source supports>\n" +
	"  Drift: <the change that makes this source stale>\n"

// decisionMapAssetRule states, in the rendered template, where a map-owned asset stays.
// The candidate scanner lists only the direct children of a decisions directory, so an
// asset in this nested directory is never read as a decision map. The sentence renders
// without Markdown quoting, because the anchor registry pins the bare path.
const decisionMapAssetRule = "A map-owned asset stays in the map's assets folder, decisions/<topic>/assets/."

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
	Notes       string
	Decisions   string
	// IndexSections holds the index sections the document opened. A section can
	// be present and empty, which the body alone cannot report.
	IndexSections map[string]bool
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

// fieldScan drives the shared field scan over one map index. The index owns no
// ticket field, so the scoped half of the field table stays with ticketFileScan.
func (s decisionMapSchema) fieldScan() FieldScan {
	var table []FieldSpec
	for _, f := range s.fields {
		if f.scoped {
			continue
		}
		table = append(table, FieldSpec{Name: f.name, Syntax: f.syntax, Heading: f.heading})
	}
	for _, terminal := range append(append([]terminalSection{}, s.terminalSections...), s.indexSections...) {
		table = append(table, FieldSpec{Name: terminal.heading, Syntax: terminal.syntax, Heading: true})
	}
	return FieldScan{
		Table: table,
		Duplicate: func(spec FieldSpec, _ string) string {
			if spec.Name == "title" || spec.Name == "Status" {
				return "duplicate " + spec.Name
			}
			return "duplicate " + spec.Name + " section"
		},
	}
}

// parseDecisionMap parses one map index. Every decision lives in a ticket file
// beside the index, so an inline decision-ticket heading is refused here and the
// ticket-count rule belongs to the caller that read the folder. The topic names
// the index's own folder, which the refusal points at.
func parseDecisionMap(topic string, content []byte) (DecisionMap, []Diagnostic) {
	m := DecisionMap{IndexSections: map[string]bool{}}
	var diagnostics []Diagnostic
	seenTerminal := map[string]bool{}
	section := ""

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
				diagnostics = append(diagnostics, Diagnostic{Message: "unsupported Handoff section"})
				section = "Handoff"
				unsupported = true
				break
			}
		}
		if unsupported {
			continue
		}
		if match := inlineTicketHeading.FindStringSubmatch(line); match != nil {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("inline ticket #%s: move it to %s/%s/%s.md", match[1], topic, ticketsDirName, match[1])})
			section = "inline ticket"
			continue
		}
		if heading := canonicalDecisionMapSchema.indexHeading(line); heading != "" {
			if duplicate {
				report()
			}
			m.IndexSections[heading] = true
			section = heading
			continue
		}
		if heading := canonicalDecisionMapSchema.terminalHeading(line); heading != "" {
			if duplicate {
				report()
			}
			seenTerminal[heading] = true
			section = heading
			continue
		}
		if entry.Field == "Destination" {
			if duplicate {
				report()
			}
			section = "Destination"
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
		case "Notes":
			m.Notes = appendSectionLine(m.Notes, line)
		case "Decisions so far":
			m.Decisions = appendSectionLine(m.Decisions, line)
		}
	}
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

// DecisionMapTemplate renders the canonical map-index Markdown skeleton. The
// decisions themselves live in the ticket files DecisionTicketTemplate renders.
func DecisionMapTemplate() string {
	var b strings.Builder
	b.WriteString(canonicalDecisionMapSchema.field("title").syntax)
	b.WriteString("<decision map title>\n\n")
	b.WriteString(canonicalDecisionMapSchema.field("Status").syntax)
	b.WriteString(canonicalDecisionMapSchema.statuses[0])
	b.WriteString("\n\n")
	b.WriteString(canonicalDecisionMapSchema.field("Destination").syntax)
	b.WriteString("\n\n<what this map decides>\n")
	for _, index := range canonicalDecisionMapSchema.indexSections {
		body := "- [<decision question>](" + templateTicketsSegment + "/1.md): <gist>"
		if index.heading == "Notes" {
			body = decisionMapAssetRule
		}
		b.WriteString("\n" + index.syntax + "\n\n" + body + "\n")
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
