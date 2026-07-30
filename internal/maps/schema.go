package maps

import (
	"fmt"
	"os"
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
	var candidates []DecisionMapCandidate
	appendDirectory := func(dir string, compiled bool) error {
		discovered, state, reason := discoverDirectoryCandidates(root, dir, compiled)
		if state == bounds.StateAbsent {
			return nil
		}
		if state.Failed() {
			return fmt.Errorf("read %s: %s", dir, reason)
		}
		candidates = append(candidates, discovered...)
		return nil
	}
	if err := appendDirectory(filepath.Join(root, DecisionsDir), false); err != nil {
		return nil, err
	}
	specs, err := os.ReadDir(filepath.Join(root, "specs"))
	if os.IsNotExist(err) {
		return candidates, nil
	}
	if err != nil {
		return nil, err
	}
	for _, spec := range specs {
		if !spec.IsDir() || strings.HasPrefix(spec.Name(), ".") {
			continue
		}
		if err := appendDirectory(filepath.Join(root, "specs", spec.Name(), DecisionsDir), true); err != nil {
			return nil, err
		}
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Path < candidates[j].Path })
	return candidates, nil
}

type field struct {
	name    string
	syntax  string
	heading bool
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
		{name: "Blocked by", syntax: "Blocked by: "},
		{name: "Type", syntax: "Type: "},
		{name: "Question", syntax: "### Question", heading: true},
		{name: "Answer", syntax: "### Answer", heading: true},
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

// ParseDecisionMap parses a decision map according to the canonical schema.
func ParseDecisionMap(content []byte) (DecisionMap, []Diagnostic) {
	var m DecisionMap
	var diagnostics []Diagnostic
	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	seenTerminal := map[string]bool{}
	var current *DecisionTicket
	var ticketFields map[string]bool
	section := ""
	inFence := false
	seenTitle, seenStatus, seenDestination := false, false, false

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
		if !ticketFields["Answer"] {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: missing Answer", current.ID)})
		}
		m.Tickets = append(m.Tickets, *current)
		current = nil
		ticketFields = nil
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if strings.HasPrefix(line, canonicalDecisionMapSchema.field("title").syntax) {
			if seenTitle {
				diagnostics = append(diagnostics, Diagnostic{Message: "duplicate title"})
			} else {
				m.Title = strings.TrimSpace(strings.TrimPrefix(line, canonicalDecisionMapSchema.field("title").syntax))
				seenTitle = true
			}
			continue
		}
		if strings.HasPrefix(line, canonicalDecisionMapSchema.field("Status").syntax) {
			if seenStatus {
				diagnostics = append(diagnostics, Diagnostic{Message: "duplicate Status"})
			} else {
				m.Status = strings.TrimSpace(strings.TrimPrefix(line, canonicalDecisionMapSchema.field("Status").syntax))
				seenStatus = true
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
			if seenTerminal[heading] {
				diagnostics = append(diagnostics, Diagnostic{Message: "duplicate " + heading + " section"})
			}
			seenTerminal[heading] = true
			section = heading
			continue
		}
		if line == canonicalDecisionMapSchema.field("Destination").syntax {
			finishTicket()
			if seenDestination {
				diagnostics = append(diagnostics, Diagnostic{Message: "duplicate Destination section"})
			}
			seenDestination = true
			section = "Destination"
			continue
		}
		if id, title, ok := canonicalDecisionMapSchema.ticket(line); ok {
			finishTicket()
			current = &DecisionTicket{ID: id, Title: title}
			ticketFields = map[string]bool{}
			section = "ticket"
			continue
		}
		if current != nil {
			setField := func(name, value string, set func(string)) {
				if ticketFields[name] {
					diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: duplicate %s", current.ID, name)})
					return
				}
				ticketFields[name] = true
				set(value)
			}
			switch {
			case strings.HasPrefix(line, canonicalDecisionMapSchema.field("Blocked by").syntax):
				setField("Blocked by", strings.TrimSpace(strings.TrimPrefix(line, canonicalDecisionMapSchema.field("Blocked by").syntax)), func(value string) { current.BlockedBy = value })
			case strings.HasPrefix(line, canonicalDecisionMapSchema.field("Type").syntax):
				setField("Type", strings.TrimSpace(strings.TrimPrefix(line, canonicalDecisionMapSchema.field("Type").syntax)), func(value string) { current.Type = value })
			case line == canonicalDecisionMapSchema.field("Question").syntax:
				if ticketFields["Question"] {
					diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: duplicate Question", current.ID)})
				}
				ticketFields["Question"] = true
				section = "Question"
			case line == canonicalDecisionMapSchema.field("Answer").syntax:
				if ticketFields["Answer"] {
					diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: duplicate Answer", current.ID)})
				}
				ticketFields["Answer"] = true
				section = "Answer"
			case strings.TrimSpace(line) != "" && section == "Question":
				current.Question = appendSectionLine(current.Question, strings.TrimSpace(line))
			case strings.TrimSpace(line) != "" && section == "Answer":
				current.Answer = appendSectionLine(current.Answer, strings.TrimSpace(line))
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
	}
	return b.String()
}
