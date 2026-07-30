package maps

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
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
	for _, ticket := range m.Tickets {
		seen := make(map[string]bool)
		for _, blocker := range blockers(ticket.BlockedBy) {
			if seen[blocker] {
				diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: %s duplicate blocker #%s", ticket.ID, ticket.Title, blocker)})
				continue
			}
			seen[blocker] = true
			blockerTicket, exists := byID[blocker]
			if !exists {
				diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: %s dangling blocker #%s", ticket.ID, ticket.Title, blocker)})
				continue
			}
			if blocker == ticket.ID {
				diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("ticket #%s: %s self-edge #%s -> #%s", ticket.ID, ticket.Title, ticket.ID, blocker)})
			}
			if resolved(ticket) && !resolved(blockerTicket) {
				diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("resolved ticket #%s: %s depends on unresolved #%s: %s", ticket.ID, ticket.Title, blocker, blockerTicket.Title)})
			}
		}
	}
	diagnostics = append(diagnostics, cycleDiagnostics(byID)...)
	return diagnostics
}

func cycleDiagnostics(tickets map[string]DecisionTicket) []Diagnostic {
	var diagnostics []Diagnostic
	visiting, visited := make(map[string]bool), make(map[string]bool)
	var visit func(string, string)
	visit = func(id, from string) {
		if visiting[id] {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("cycle edge #%s -> #%s", from, id)})
			return
		}
		if visited[id] {
			return
		}
		visiting[id] = true
		for _, blocker := range blockers(tickets[id].BlockedBy) {
			if _, exists := tickets[blocker]; exists && blocker != id {
				visit(blocker, id)
			}
		}
		delete(visiting, id)
		visited[id] = true
	}
	ids := make([]string, 0, len(tickets))
	for id := range tickets {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		visit(id, "")
	}
	return diagnostics
}

func blockers(value string) []string {
	if value == "none" {
		return nil
	}
	parts := strings.Split(value, ", ")
	for i := range parts {
		parts[i] = strings.TrimPrefix(parts[i], "#")
	}
	return parts
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
		supports, drift := false, false
		for _, line := range lines[1:] {
			if strings.HasPrefix(line, "Supports: ") && strings.TrimSpace(strings.TrimPrefix(line, "Supports: ")) != "" {
				supports = true
			}
			if strings.HasPrefix(line, "Drift: ") && strings.TrimSpace(strings.TrimPrefix(line, "Drift: ")) != "" {
				drift = true
			}
		}
		if !supports {
			diagnostics = append(diagnostics, Diagnostic{Message: "Sources " + value + " missing Supports"})
		}
		if !drift {
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
		if strings.HasPrefix(line, "- Path: ") || strings.HasPrefix(line, "- URL: ") {
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
