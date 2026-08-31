// Package tickets owns the ticket-file grammar. It is a pure decision domain:
// immutable bytes in, a parsed ticket and ordered diagnostics out. It reads no
// file, and it drives the shared field scan and graph walk that internal/maps
// exports, so the two schemas cannot drift.
package tickets

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/maps"
)

// Ticket is the parsed form of one ticket file.
type Ticket struct {
	// Name is the file basename, which is the identity a blocker edge resolves against.
	Name  string
	Title string
	// Blockers holds the declared blocker basenames in declared order. The none sentinel yields none.
	Blockers []string
	// Writes holds the declared ownership entries in declared order.
	Writes []string
	// Covers holds the declared row-ID tokens in declared order.
	Covers []string
}

const (
	fieldTitle     = "title"
	fieldBlockedBy = "Blocked by"
	fieldWrites    = "Writes"
	fieldCovers    = "Covers"
	sectionBuild   = "What to build"
	sectionAccept  = "Acceptance"
	sectionCharge  = "Delegate charge"
	noneSentinel   = "none"
	fenceMarker    = "```"
)

var ticketFields = []maps.FieldSpec{
	{Name: sectionBuild, Syntax: "## What to build", Heading: true},
	{Name: sectionAccept, Syntax: "## Acceptance", Heading: true},
	{Name: sectionCharge, Syntax: "## Delegate charge", Heading: true},
	{Name: fieldBlockedBy, Syntax: "Blocked by:"},
	{Name: fieldWrites, Syntax: "Writes:"},
	{Name: fieldCovers, Syntax: "Covers:"},
	{Name: fieldTitle, Syntax: "# "},
}

// requiredFields names every field a ticket must carry, in diagnostic order.
var requiredFields = []string{fieldTitle, fieldBlockedBy, fieldWrites, fieldCovers, sectionBuild, sectionAccept}

var coversToken = regexp.MustCompile(`^([A-Za-z]+)([1-9][0-9]*)$`)

func fieldScan() maps.FieldScan {
	return maps.FieldScan{
		Table: ticketFields,
		Duplicate: func(spec maps.FieldSpec, _ string) string {
			if spec.Heading {
				return "duplicate " + spec.Name + " section"
			}
			return "duplicate " + spec.Name
		},
	}
}

// missingMessage names one absent required field.
func missingMessage(name string) string {
	for _, spec := range ticketFields {
		if spec.Name == name && spec.Heading {
			return "missing " + name + " section"
		}
	}
	return "missing " + name
}

// ParseTicket grades one ticket file. Name is the file basename, siblings holds
// every basename in the same tickets folder, and tag is the spec tag. An empty
// tag skips the citation tag rule, which is the tickets-only posture.
// The diagnostics arrive in one stable order: duplicate fields, an unterminated
// fence, absent required fields, the citation grammar, then the declared edges.
func ParseTicket(name string, content []byte, siblings []string, tag string) (Ticket, []string) {
	ticket := Ticket{Name: name}
	scanned, diagnostics := fieldScan().Scan(content)

	values := make(map[string]string, len(ticketFields))
	present := make(map[string]bool, len(ticketFields))
	fences := 0
	for _, entry := range scanned {
		if strings.HasPrefix(entry.Text, fenceMarker) {
			fences++
		}
		if entry.Fenced || entry.Field == "" || entry.Diagnostic != "" {
			continue
		}
		present[entry.Field] = true
		if _, seen := values[entry.Field]; !seen {
			values[entry.Field] = entry.Value
		}
	}
	if fences%2 == 1 {
		diagnostics = append(diagnostics, "unterminated fence")
	}

	ticket.Title = values[fieldTitle]
	ticket.Blockers = list(values[fieldBlockedBy])
	ticket.Writes = list(values[fieldWrites])
	ticket.Covers = list(values[fieldCovers])

	for _, required := range requiredFields {
		if !present[required] {
			diagnostics = append(diagnostics, missingMessage(required))
			continue
		}
		if !isHeading(required) && values[required] == "" {
			diagnostics = append(diagnostics, missingMessage(required))
		}
	}
	diagnostics = append(diagnostics, coversDiagnostics(ticket, tag)...)
	return ticket, append(diagnostics, edgeDiagnostics(ticket, siblings)...)
}

// list splits one comma-separated field value. An absent or empty value, like
// the none sentinel, declares no entry.
func list(value string) []string {
	if value == "" {
		return nil
	}
	return maps.FieldList(value, noneSentinel, "")
}

func isHeading(name string) bool {
	for _, spec := range ticketFields {
		if spec.Name == name {
			return spec.Heading
		}
	}
	return false
}

// coversDiagnostics grades the citation grammar of one ticket.
func coversDiagnostics(ticket Ticket, tag string) []string {
	var diagnostics []string
	seen := make(map[string]bool, len(ticket.Covers))
	for _, token := range ticket.Covers {
		if seen[token] {
			diagnostics = append(diagnostics, fmt.Sprintf("duplicate Covers token %s", token))
			continue
		}
		seen[token] = true
		match := coversToken.FindStringSubmatch(token)
		if match == nil {
			diagnostics = append(diagnostics, fmt.Sprintf("malformed Covers token %q", token))
			continue
		}
		if tag != "" && match[1] != tag {
			diagnostics = append(diagnostics, fmt.Sprintf("foreign Covers token %s: spec tag is %s", token, tag))
		}
	}
	return diagnostics
}

// edgeDiagnostics grades one ticket's declared blocker edges against its
// siblings. The sibling nodes declare no edge, so the walk reports the
// duplicate, dangling, and self classes here; Cycles owns the cycle class.
func edgeDiagnostics(ticket Ticket, siblings []string) []string {
	walk := maps.GraphWalk{Names: []string{ticket.Name}, Edges: [][]string{ticket.Blockers}}
	for _, sibling := range sortedOthers(siblings, ticket.Name) {
		walk.Names = append(walk.Names, sibling)
		walk.Edges = append(walk.Edges, nil)
	}
	walk.Fault = func(fault maps.GraphFault, _ int, target string) string {
		return fmt.Sprintf("%s: %s blocker %s", ticket.Name, faultWord(fault), target)
	}
	return walk.Diagnostics()
}

func faultWord(fault maps.GraphFault) string {
	switch fault {
	case maps.FaultDuplicateEdge:
		return "duplicate"
	case maps.FaultDanglingEdge:
		return "dangling"
	default:
		return "self"
	}
}

// Cycles reports one edge of each blocker cycle over a parsed ticket set. Every
// duplicate, dangling, and self edge is dropped first, because ParseTicket
// already named it and only a resolved forward edge can close a cycle.
func Cycles(parsed []Ticket) []string {
	known := make(map[string]bool, len(parsed))
	for _, ticket := range parsed {
		known[ticket.Name] = true
	}
	walk := maps.GraphWalk{Names: make([]string, 0, len(parsed)), Edges: make([][]string, 0, len(parsed))}
	for _, ticket := range parsed {
		seen := make(map[string]bool, len(ticket.Blockers))
		var edges []string
		for _, blocker := range ticket.Blockers {
			if seen[blocker] || !known[blocker] || blocker == ticket.Name {
				continue
			}
			seen[blocker] = true
			edges = append(edges, blocker)
		}
		walk.Names = append(walk.Names, ticket.Name)
		walk.Edges = append(walk.Edges, edges)
	}
	walk.Fault = func(_ maps.GraphFault, node int, target string) string {
		return fmt.Sprintf("cycle edge %s -> %s", walk.Names[node], target)
	}
	return walk.Diagnostics()
}

func sortedOthers(names []string, exclude string) []string {
	others := make([]string, 0, len(names))
	for _, name := range names {
		if name != exclude {
			others = append(others, name)
		}
	}
	sort.Strings(others)
	return others
}
