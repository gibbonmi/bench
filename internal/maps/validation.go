package maps

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ValidateDecisionMap validates one map at its repository-relative path.
func ValidateDecisionMap(root, path string, compiled bool, content []byte) (DecisionMap, []Diagnostic) {
	m, parsed := ParseDecisionMap(content)
	diagnostics := make([]Diagnostic, 0, len(parsed))
	for _, diagnostic := range parsed {
		diagnostics = append(diagnostics, Diagnostic{Message: path + ": " + diagnostic.Message})
	}
	for _, diagnostic := range graphDiagnostics(m) {
		diagnostics = append(diagnostics, Diagnostic{Message: path + ": " + diagnostic.Message})
	}
	for _, diagnostic := range readinessDiagnostics(root, m, compiled) {
		diagnostics = append(diagnostics, Diagnostic{Message: path + ": " + diagnostic.Message})
	}
	return m, diagnostics
}

func graphDiagnostics(m DecisionMap) []Diagnostic {
	var diagnostics []Diagnostic
	byID := make(map[string]DecisionTicket)
	for _, ticket := range m.Tickets {
		if original, exists := byID[ticket.ID]; exists {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("duplicate ID #%s: %s conflicts with %s", ticket.ID, ticket.Title, original.Title)})
			continue
		}
		byID[ticket.ID] = ticket
	}
	walk := GraphWalk{}
	for _, ticket := range m.Tickets {
		walk.Names = append(walk.Names, ticket.ID)
		walk.Edges = append(walk.Edges, blockers(ticket.BlockedBy))
	}
	walk.Fault = func(fault GraphFault, node int, target string) string {
		ticket := m.Tickets[node]
		switch fault {
		case FaultDuplicateEdge:
			return fmt.Sprintf("ticket #%s: %s duplicate blocker #%s", ticket.ID, ticket.Title, target)
		case FaultDanglingEdge:
			return fmt.Sprintf("ticket #%s: %s dangling blocker #%s", ticket.ID, ticket.Title, target)
		case FaultSelfEdge:
			return fmt.Sprintf("ticket #%s: %s self-edge #%s -> #%s", ticket.ID, ticket.Title, ticket.ID, target)
		default:
			blocked := byID[target]
			return fmt.Sprintf("cycle edge ticket #%s: %s -> ticket #%s: %s (#%s -> #%s)", ticket.ID, ticket.Title, target, blocked.Title, ticket.ID, target)
		}
	}
	walk.Edge = func(node int, target string) string {
		ticket, blocker := m.Tickets[node], byID[target]
		if resolved(ticket) && !resolved(blocker) {
			return fmt.Sprintf("resolved ticket #%s: %s depends on unresolved #%s: %s", ticket.ID, ticket.Title, target, blocker.Title)
		}
		return ""
	}
	for _, message := range walk.Diagnostics() {
		diagnostics = append(diagnostics, Diagnostic{Message: message})
	}
	return diagnostics
}

// blockers names the decision-map dependency list inside the shared field grammar.
func blockers(value string) []string {
	return FieldList(value, "none", "#")
}

func resolved(ticket DecisionTicket) bool {
	return ticketAnswerState(ticket) == "resolved"
}

func ticketAnswerState(ticket DecisionTicket) string {
	answer := strings.TrimSpace(ticket.Answer)
	switch {
	case strings.HasPrefix(answer, "— (deferred"), strings.Contains(answer, "GRILL DEFERRED"):
		return "deferred"
	case answer == "", strings.HasPrefix(answer, "— (open"):
		return "frontier"
	default:
		return "resolved"
	}
}

func readinessDiagnostics(root string, m DecisionMap, compiled bool) []Diagnostic {
	var diagnostics []Diagnostic
	if compiled && m.Status != "ready" {
		diagnostics = append(diagnostics, Diagnostic{Message: "compiled map must be ready"})
	}
	if m.Status == "ready" {
		for _, ticket := range m.Tickets {
			if !resolved(ticket) {
				diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ready map has unresolved ticket #%s: %s", ticket.ID, ticket.Title)})
			}
		}
		if strings.TrimSpace(m.Fog) != "" {
			diagnostics = append(diagnostics, Diagnostic{Message: "ready map has non-empty Not yet specified"})
		}
	}
	for _, section := range []struct{ name, body string }{
		{"Not yet specified", m.Fog}, {"Spec-writer discretion", m.Discretion}, {"Out of scope", m.OutOfScope},
	} {
		if strings.TrimSpace(section.body) != "" && !markdownBullets(section.body) {
			diagnostics = append(diagnostics, Diagnostic{Message: section.name + " must be a Markdown bullet list"})
		}
	}
	diagnostics = append(diagnostics, sourceDiagnostics(root, m.Sources)...)
	return diagnostics
}

func markdownBullets(body string) bool {
	seenBullet := false
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			seenBullet = true
			continue
		}
		if !seenBullet || (len(line) == len(strings.TrimLeft(line, " \t"))) {
			return false
		}
	}
	return seenBullet
}

func sourceDiagnostics(root, body string) []Diagnostic {
	var diagnostics []Diagnostic
	var record []string
	validate := func(lines []string) {
		if len(lines) == 0 {
			return
		}
		locator := strings.TrimPrefix(lines[0], "- ")
		isPath := strings.HasPrefix(locator, "Path: ")
		isURL := strings.HasPrefix(locator, "URL: ")
		if !isPath && !isURL {
			diagnostics = append(diagnostics, Diagnostic{Message: "Sources entry must start with Path or URL"})
			return
		}
		value := sourceLocator(strings.TrimPrefix(strings.TrimPrefix(locator, "Path: "), "URL: "))
		expected := []string{"Supports", "Drift"}
		seen := make(map[string]bool, len(expected))
		next := 0
		for _, line := range lines[1:] {
			name, fieldValue, hasSeparator := strings.Cut(line, ":")
			if !hasSeparator || (name != "Supports" && name != "Drift") {
				diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("Sources %s unexpected field %s", value, name)})
				continue
			}
			if seen[name] {
				diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("Sources %s duplicate field %s", value, name)})
				continue
			}
			if next >= len(expected) || name != expected[next] {
				expectedName := "no further fields"
				if next < len(expected) {
					expectedName = expected[next]
				}
				diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("Sources %s field %s is out of order; expected %s", value, name, expectedName)})
				continue
			}
			if strings.TrimSpace(fieldValue) == "" {
				diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("Sources %s field %s must be non-empty", value, name)})
				continue
			}
			seen[name] = true
			next++
		}
		if !seen["Supports"] {
			diagnostics = append(diagnostics, Diagnostic{Message: "Sources " + value + " missing Supports"})
		}
		if !seen["Drift"] {
			diagnostics = append(diagnostics, Diagnostic{Message: "Sources " + value + " missing Drift"})
		}
		if isPath {
			if reason := validateSourcePath(root, value); reason != "" {
				diagnostics = append(diagnostics, Diagnostic{Message: "Sources Path " + value + " " + reason})
			}
		} else if !validSourceURL(value) {
			diagnostics = append(diagnostics, Diagnostic{Message: "Sources URL " + value + " must be an absolute HTTP(S) URL with a host"})
		}
	}
	for _, line := range nonEmptyLines(body) {
		if strings.HasPrefix(line, "- ") {
			validate(record)
			record = []string{line}
		} else if len(record) == 0 {
			diagnostics = append(diagnostics, Diagnostic{Message: "Sources entry must start with a Markdown bullet"})
		} else {
			record = append(record, line)
		}
	}
	validate(record)
	return diagnostics
}

func nonEmptyLines(body string) []string {
	var lines []string
	for _, line := range strings.Split(body, "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	return lines
}

func sourceLocator(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "`") || strings.HasSuffix(value, "`") {
		if len(value) < 2 || !strings.HasPrefix(value, "`") || !strings.HasSuffix(value, "`") {
			return ""
		}
		return value[1 : len(value)-1]
	}
	return value
}

func validateSourcePath(root, source string) string {
	if source == "" || filepath.IsAbs(source) {
		return "must be a non-empty repository-relative path"
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return "root cannot be resolved"
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "root cannot be resolved"
	}
	path := filepath.Join(root, source)
	if relative, err := filepath.Rel(root, path); err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "escapes the repository"
	}
	info, err := os.Lstat(path)
	if err != nil {
		return "must name an existing regular file"
	}
	if info.Mode()&os.ModeSymlink != 0 {
		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			return "must name an existing regular file"
		}
	}
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "escapes the repository"
	}
	info, err = os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return "must name an existing regular file"
	}
	return ""
}

func validSourceURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}
