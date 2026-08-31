package maps

import (
	"sort"
	"strconv"
	"strings"
)

// FieldSpec declares one line-oriented field of a Markdown grammar.
type FieldSpec struct {
	Name string
	// Syntax is the prefix that introduces the field, or the whole line when Heading is true.
	Syntax  string
	Heading bool
	// Scoped counts a repeat inside the open scope, not across the document.
	Scoped bool
}

// FieldLine is one source line resolved against a field table.
type FieldLine struct {
	// Text is the line without its newline and with CRLF already normalized.
	Text string
	// Fenced marks a code-fence marker or a line inside a fence. Such a line matches no field.
	Fenced bool
	// Field names the matched field, and stays empty when the line matches none.
	Field string
	Value string
	// Scope is the open scope identity at this line.
	Scope string
	// Diagnostic holds the duplicate-field message when this occurrence repeats an earlier one.
	Diagnostic string
}

// FieldScan reads a line-oriented field grammar out of immutable bytes. One scan
// serves every schema that declares its fields as prefixes or exact headings.
type FieldScan struct {
	// Table is matched in order; the first match wins.
	Table []FieldSpec
	// Scope reports a scope change: the opened scope identity, or an empty
	// identity when the line closes the open scope. Two scopes with one
	// identity stay distinct for duplicate accounting.
	Scope func(line string) (identity string, changed bool)
	// Duplicate formats the diagnostic for a repeated field. Scope is the open scope identity.
	Duplicate func(spec FieldSpec, scope string) string
}

// Scan resolves every line against the table and returns the duplicate-field
// diagnostics in line order. A fenced line never parses as a field.
func (s FieldScan) Scan(content []byte) ([]FieldLine, []string) {
	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	scanned := make([]FieldLine, 0, len(lines))
	var diagnostics []string
	seen := make(map[string]bool)
	scope, scopeKey := "", ""
	scopes := 0
	inFence := false
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
			scanned = append(scanned, FieldLine{Text: line, Fenced: true})
			continue
		}
		if inFence {
			scanned = append(scanned, FieldLine{Text: line, Fenced: true})
			continue
		}
		if s.Scope != nil {
			if identity, changed := s.Scope(line); changed {
				scopes++
				scope, scopeKey = identity, strconv.Itoa(scopes)+"\x00"+identity
				if identity == "" {
					scopeKey = ""
				}
			}
		}
		entry := FieldLine{Text: line, Scope: scope}
		for _, spec := range s.Table {
			if spec.Heading && line != spec.Syntax {
				continue
			}
			if !spec.Heading && !strings.HasPrefix(line, spec.Syntax) {
				continue
			}
			entry.Field = spec.Name
			entry.Value = strings.TrimSpace(strings.TrimPrefix(line, spec.Syntax))
			if spec.Scoped && scopeKey == "" {
				break
			}
			key := "\x00" + spec.Name
			if spec.Scoped {
				key = scopeKey + key
			}
			if seen[key] && s.Duplicate != nil {
				entry.Diagnostic = s.Duplicate(spec, scope)
				diagnostics = append(diagnostics, entry.Diagnostic)
			}
			seen[key] = true
			break
		}
		scanned = append(scanned, entry)
	}
	return scanned, diagnostics
}

// FieldList splits one comma-separated field value into its entries. The none
// sentinel yields no entry, and itemPrefix is trimmed from each entry.
func FieldList(value, none, itemPrefix string) []string {
	if value == none {
		return nil
	}
	entries := strings.Split(value, ", ")
	for i := range entries {
		entries[i] = strings.TrimPrefix(entries[i], itemPrefix)
	}
	return entries
}

// GraphFault names one class of declared-edge fault.
type GraphFault string

// The declared-edge fault classes one walk reports.
const (
	FaultDuplicateEdge GraphFault = "duplicate"
	FaultDanglingEdge  GraphFault = "dangling"
	FaultSelfEdge      GraphFault = "self"
	FaultCycleEdge     GraphFault = "cycle"
)

// GraphWalk grades the declared edges of a named-node dependency graph. The
// walk carries no schema vocabulary; the caller's formatter supplies every word.
type GraphWalk struct {
	// Names holds one identity per node in declaration order. A repeated identity resolves to its first node.
	Names []string
	// Edges holds the declared target identities of each node, in declared order.
	Edges [][]string
	// Fault formats one fault. The node index addresses Names, and target names the declared edge target.
	Fault func(fault GraphFault, node int, target string) string
	// Edge grades one resolved, non-duplicate edge. An empty return adds no diagnostic.
	Edge func(node int, target string) string
}

// Diagnostics reports every duplicate, dangling, and self edge in declaration
// order, then one edge of each cycle.
func (w GraphWalk) Diagnostics() []string {
	first := make(map[string]int, len(w.Names))
	for i, name := range w.Names {
		if _, exists := first[name]; !exists {
			first[name] = i
		}
	}
	var diagnostics []string
	for node := range w.Names {
		seen := make(map[string]bool)
		for _, target := range w.edges(node) {
			if seen[target] {
				diagnostics = append(diagnostics, w.Fault(FaultDuplicateEdge, node, target))
				continue
			}
			seen[target] = true
			if _, exists := first[target]; !exists {
				diagnostics = append(diagnostics, w.Fault(FaultDanglingEdge, node, target))
				continue
			}
			if target == w.Names[node] {
				diagnostics = append(diagnostics, w.Fault(FaultSelfEdge, node, target))
			}
			if w.Edge == nil {
				continue
			}
			if message := w.Edge(node, target); message != "" {
				diagnostics = append(diagnostics, message)
			}
		}
	}
	return append(diagnostics, w.cycleDiagnostics(first)...)
}

func (w GraphWalk) edges(node int) []string {
	if node < 0 || node >= len(w.Edges) {
		return nil
	}
	return w.Edges[node]
}

func (w GraphWalk) cycleDiagnostics(first map[string]int) []string {
	var diagnostics []string
	visiting, visited := make(map[string]bool), make(map[string]bool)
	var visit func(name, from string)
	visit = func(name, from string) {
		if visiting[name] {
			diagnostics = append(diagnostics, w.Fault(FaultCycleEdge, first[from], name))
			return
		}
		if visited[name] {
			return
		}
		visiting[name] = true
		for _, target := range w.edges(first[name]) {
			if _, exists := first[target]; exists && target != name {
				visit(target, name)
			}
		}
		delete(visiting, name)
		visited[name] = true
	}
	names := make([]string, 0, len(first))
	for name := range first {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		visit(name, "")
	}
	return diagnostics
}
