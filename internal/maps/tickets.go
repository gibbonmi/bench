package maps

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// ticketsDirName is the folder a split map keeps its decision ticket files in.
// The folder sits under the index's own topic folder, one level below the
// candidate scan, so the scan never reads a ticket file as a map.
const ticketsDirName = "tickets"

var ticketBasename = regexp.MustCompile(`^[1-9][0-9]*$`)

// gistLine is one Decisions so far entry. The link is relative to the map file,
// and the parser resolves the entry by the ticket-number segment alone.
var gistLine = regexp.MustCompile(`^- \[[^\]]*\]\([^)]*/` + ticketsDirName + `/([1-9][0-9]*)\.md\):\s*\S`)

// indexHeading reports the index section one line opens, or an empty name.
func (s decisionMapSchema) indexHeading(line string) string {
	for _, section := range s.indexSections {
		if line == section.syntax {
			return section.heading
		}
	}
	return ""
}

// ticketFileScan drives the shared field scan over one decision ticket file.
// The whole file is one scope, so no field declares itself scoped and the
// duplicate diagnostic names the ticket the file is.
func (s decisionMapSchema) ticketFileScan(id string) FieldScan {
	var table []FieldSpec
	for _, name := range []string{"title", "Blocked by", "Type", "Question", "Answer"} {
		f := s.field(name)
		table = append(table, FieldSpec{Name: f.name, Syntax: f.syntax, Heading: f.heading})
	}
	return FieldScan{
		Table:     table,
		Duplicate: func(spec FieldSpec, _ string) string { return fmt.Sprintf("ticket #%s: duplicate %s", id, spec.Name) },
	}
}

// ticketDiagnostics grades one decision ticket's fields. Both file shapes share
// it, so an inline ticket and a ticket file answer with one message set.
func ticketDiagnostics(ticket DecisionTicket, answerSeen bool) []Diagnostic {
	var diagnostics []Diagnostic
	report := func(message string) {
		diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: %s", ticket.ID, message)})
	}
	if ticket.BlockedBy == "" {
		report("missing Blocked by")
	} else if !blockersField.MatchString(ticket.BlockedBy) {
		report("malformed Blocked by")
	}
	if ticket.Type == "" {
		report("missing Type")
	} else if !canonicalDecisionMapSchema.hasType(ticket.Type) {
		report(fmt.Sprintf("unsupported Type %q", ticket.Type))
	}
	if ticket.Question == "" {
		report("missing Question")
	}
	if !answerSeen {
		report("missing Answer")
	}
	return diagnostics
}

// parseDecisionTicket reads one ticket file. The basename supplies the id, and
// the `# ` line supplies the title.
func parseDecisionTicket(id string, content []byte) (DecisionTicket, []Diagnostic) {
	ticket := DecisionTicket{ID: id}
	scanned, _ := canonicalDecisionMapSchema.ticketFileScan(id).Scan(content)
	var diagnostics []Diagnostic
	answerSeen, section := false, ""
	for _, entry := range scanned {
		if entry.Fenced {
			continue
		}
		if entry.Diagnostic != "" {
			diagnostics = append(diagnostics, Diagnostic{Message: entry.Diagnostic})
			continue
		}
		switch entry.Field {
		case "title":
			ticket.Title = entry.Value
		case "Blocked by":
			ticket.BlockedBy = entry.Value
		case "Type":
			ticket.Type = entry.Value
		case "Question":
			section = "Question"
		case "Answer":
			answerSeen, section = true, "Answer"
		default:
			line := strings.TrimSpace(entry.Text)
			if line == "" {
				break
			}
			if section == "Question" {
				ticket.Question = appendSectionLine(ticket.Question, line)
			} else if section == "Answer" {
				ticket.Answer = appendSectionLine(ticket.Answer, line)
			}
		}
	}
	if ticket.Title == "" {
		diagnostics = append(diagnostics, Diagnostic{Message: "ticket #" + id + ": missing title"})
	}
	return ticket, append(diagnostics, ticketDiagnostics(ticket, answerSeen)...)
}

// ticketFolder is where the ticket files of the map at path live.
func ticketFolder(root, path string) string {
	return filepath.Join(root, filepath.FromSlash(strings.TrimSuffix(path, ".md")), ticketsDirName)
}

// splitTickets reads the ticket files beside one map index. The present flag
// separates a split map from an inline one, which keeps today's rules.
func splitTickets(root, path string) (tickets []DecisionTicket, diagnostics []Diagnostic, present bool) {
	folder := ticketFolder(root, path)
	classified := bounds.ClassifyDir(folder)
	switch classified.State {
	case bounds.StateAbsent:
		return nil, nil, false
	case bounds.StateEmpty, bounds.StateParsed:
	default:
		return nil, []Diagnostic{{Message: fmt.Sprintf("%s: %s: %s", ticketsDirName, classified.State, classified.Reason)}}, true
	}
	for _, entry := range classified.Entries {
		name := entry.Name()
		if entry.IsDir() || strings.HasPrefix(name, ".") || !strings.HasSuffix(name, ".md") || isDirectoryDoc(name) {
			continue
		}
		id := strings.TrimSuffix(name, ".md")
		if !ticketBasename.MatchString(id) {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("%s/%s: ticket file name must be a decision ticket number", ticketsDirName, name)})
			continue
		}
		file := bounds.Classify(filepath.Join(folder, name), bounds.ControlRecordLimit)
		if file.State != bounds.StateParsed && file.State != bounds.StateEmpty {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("%s/%s: %s: %s", ticketsDirName, name, file.State, file.Reason)})
			continue
		}
		ticket, ticketDiagnostics := parseDecisionTicket(id, file.Data)
		tickets = append(tickets, ticket)
		diagnostics = append(diagnostics, ticketDiagnostics...)
	}
	// A directory read is lexical, so #10 would precede #2. Numeric order keeps the
	// diagnostics and the projected rows in the order the reader numbered them.
	sort.SliceStable(tickets, func(i, j int) bool { return ticketNumber(tickets[i].ID) < ticketNumber(tickets[j].ID) })
	return tickets, diagnostics, true
}

func ticketNumber(id string) int {
	number, _ := strconv.Atoi(id)
	return number
}

// indexDiagnostics grades the index of a split map: the two sections it must
// carry, and the gists that must track its resolved tickets.
func indexDiagnostics(m DecisionMap, tickets []DecisionTicket) []Diagnostic {
	var diagnostics []Diagnostic
	for _, section := range canonicalDecisionMapSchema.indexSections {
		if !m.IndexSections[section.heading] {
			diagnostics = append(diagnostics, Diagnostic{Message: "missing " + section.heading + " section"})
		}
	}
	byID := make(map[string]DecisionTicket, len(tickets))
	for _, ticket := range tickets {
		byID[ticket.ID] = ticket
	}
	gists := map[string]bool{}
	for _, line := range nonEmptyLines(m.Decisions) {
		match := gistLine.FindStringSubmatch(line)
		if match == nil {
			diagnostics = append(diagnostics, Diagnostic{Message: "Decisions so far line has no ticket link"})
			continue
		}
		if gists[match[1]] {
			continue
		}
		gists[match[1]] = true
		switch ticket, found := byID[match[1]]; {
		case !found:
			diagnostics = append(diagnostics, Diagnostic{Message: "Decisions so far links missing ticket #" + match[1]})
		case !resolved(ticket):
			diagnostics = append(diagnostics, Diagnostic{Message: "Decisions so far links unresolved ticket #" + match[1]})
		}
	}
	for _, ticket := range tickets {
		if resolved(ticket) && !gists[ticket.ID] {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: resolved without a gist in Decisions so far", ticket.ID)})
		}
	}
	return diagnostics
}

// orphanTicketFolders reports a tickets folder whose map index is absent. The
// candidate scan reads index files only, so without this read the folder and
// every ticket in it would be graded by nothing.
func orphanTicketFolders(root, dir string) []string {
	classified := bounds.ClassifyDir(dir)
	if classified.State != bounds.StateParsed {
		return nil
	}
	var diagnostics []string
	for _, entry := range classified.Entries {
		topic := entry.Name()
		if !entry.IsDir() || strings.HasPrefix(topic, ".") {
			continue
		}
		folder := filepath.Join(dir, topic, ticketsDirName)
		if info, err := os.Stat(folder); err != nil || !info.IsDir() {
			continue
		}
		if info, err := os.Stat(filepath.Join(dir, topic+".md")); err == nil && info.Mode().IsRegular() {
			continue
		}
		rel, err := filepath.Rel(root, folder)
		if err != nil {
			rel = folder
		}
		diagnostics = append(diagnostics, fmt.Sprintf("%s: no map index at %s.md", filepath.ToSlash(rel), topic))
	}
	return diagnostics
}
