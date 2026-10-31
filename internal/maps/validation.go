package maps

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/git"
)

// ValidateDecisionMap validates one map at its repository-relative path.
func ValidateDecisionMap(root, path string, compiled bool, content []byte) (DecisionMap, []Diagnostic) {
	tickets, ticketFolderDiagnostics := splitTickets(root, path)
	m, parsed := parseDecisionMap(mapTopic(path), content)
	m.Tickets = tickets
	parsed = append(parsed, ticketFolderDiagnostics...)
	if len(tickets) == 0 {
		parsed = append(parsed, Diagnostic{Message: "missing decision ticket"})
	}
	parsed = append(parsed, indexDiagnostics(path, m, tickets)...)
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

// mapTopic names the map at a repository-relative path by its basename without ".md".
func mapTopic(path string) string { return strings.TrimSuffix(filepath.Base(path), ".md") }

func graphDiagnostics(m DecisionMap) []Diagnostic {
	var diagnostics []Diagnostic
	// A ticket id is its file's basename, so the directory read cannot deliver two
	// tickets with one id and no duplicate-id rule is reachable.
	byID := make(map[string]DecisionTicket)
	for _, ticket := range m.Tickets {
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

// SourceResolver answers whether a cleaned repository-relative slash path names a regular
// file. It returns the empty string when it does, and its own reason text when it does not.
type SourceResolver func(path string) string

// SourceDiagnostics grades the Sources section of the map at the repository-relative path
// with resolve, and prefixes each diagnostic with that path. It reads no other section, so
// the map's own structure stays the gate's grade.
func SourceDiagnostics(path string, content []byte, resolve SourceResolver) []Diagnostic {
	m, _ := parseDecisionMap(mapTopic(path), content)
	var diagnostics []Diagnostic
	for _, diagnostic := range sourceRecordDiagnostics(m.Sources, resolve) {
		diagnostics = append(diagnostics, Diagnostic{Message: path + ": " + diagnostic.Message})
	}
	return diagnostics
}

// sourceDiagnostics grades one Sources body against the checkout at root.
func sourceDiagnostics(root, body string) []Diagnostic {
	return sourceRecordDiagnostics(body, workingTreeResolver(root))
}

// sourceRecordDiagnostics is the one Sources record grammar: the record split, the field
// order, the locator rule, the URL rule, and the lexical path rule. resolve answers each
// Path locator that the lexical rule admits.
func sourceRecordDiagnostics(body string, resolve SourceResolver) []Diagnostic {
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
		// A wrapped continuation that holds a colon cuts into a name with a space, so the
		// space guard routes it to the one-physical-line message before the unknown-field
		// branch reads it as a field.
		oneLine := func(line string) Diagnostic {
			return Diagnostic{Message: fmt.Sprintf("Sources %s line %q has no field name; write each Sources record field on one physical line", value, line)}
		}
		for _, line := range lines[1:] {
			name, fieldValue, hasSeparator := strings.Cut(line, ":")
			if !hasSeparator {
				diagnostics = append(diagnostics, oneLine(line))
				continue
			}
			if strings.Contains(name, " ") {
				diagnostics = append(diagnostics, oneLine(line))
				continue
			}
			if !slices.Contains(expected, name) {
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
			if reason := sourcePathReason(value, resolve); reason != "" {
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

// The Path reasons that the lexical rule and the resolvers share: no regular file at the
// path, and a path that leads outside the repository.
const (
	notRegularFile    = "must name an existing regular file"
	escapesRepository = "escapes the repository"
)

// validateSourcePath grades one Path locator against the checkout at root.
func validateSourcePath(root, source string) string {
	return sourcePathReason(source, workingTreeResolver(root))
}

// sourcePathReason runs the lexical path rule and then resolve on the cleaned slash path.
func sourcePathReason(source string, resolve SourceResolver) string {
	clean, reason := lexicalSourcePath(source)
	if reason != "" {
		return reason
	}
	return resolve(clean)
}

// lexicalSourcePath refuses an empty, absolute, or escaping locator before any resolver
// reads it, and gives the cleaned slash path, so ./x and a/../x both reach a resolver as x.
func lexicalSourcePath(source string) (string, string) {
	if source == "" || filepath.IsAbs(source) {
		return "", "must be a non-empty repository-relative path"
	}
	clean := path.Clean(filepath.ToSlash(source))
	if clean == ".." || strings.HasPrefix(clean, "../") {
		return "", escapesRepository
	}
	return clean, ""
}

// workingTreeResolver resolves a cleaned path against the checkout at root. It follows a
// link, and a link that leads outside the repository escapes it.
func workingTreeResolver(root string) SourceResolver {
	return func(source string) string {
		root, err := canonicalpath.Resolve(root)
		if err != nil {
			return "root cannot be resolved"
		}
		path := filepath.Join(root, filepath.FromSlash(source))
		info, err := os.Lstat(path)
		if err != nil {
			return notRegularFile
		}
		if info.Mode()&os.ModeSymlink != 0 {
			path, err = filepath.EvalSymlinks(path)
			if err != nil {
				return notRegularFile
			}
		}
		relative, err := filepath.Rel(root, path)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return escapesRepository
		}
		info, err = os.Stat(path)
		if err != nil || !info.Mode().IsRegular() {
			return notRegularFile
		}
		return ""
	}
}

// PublishedTreeResolver resolves a cleaned path against its tree entry at commit in root. It
// admits a regular-file entry, mode 100644 or 100755, that the delivery closure of
// deliverable keeps; deleted holds each path that the closure deletes. A missing, link,
// directory, or submodule entry is no regular file.
func PublishedTreeResolver(root, commit, deliverable string, deleted []string) SourceResolver {
	return func(source string) string {
		entries, err := git.ListTree(root, commit, "--", ":(literal)"+source)
		if err != nil || !slices.ContainsFunc(entries, func(entry git.TreeEntry) bool { return entry.Path == source && entry.IsRegularFile() }) {
			return notRegularFile
		}
		if slices.Contains(deleted, source) {
			return "is deleted by the delivery closure of " + deliverable
		}
		return ""
	}
}

func validSourceURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}
