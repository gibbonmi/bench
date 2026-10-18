package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

// retainedWorkflowPrefix opens the diagnostic of every retained-workflow registry row.
const retainedWorkflowPrefix = "retained workflow: "

func TestRetainedWorkflow(t *testing.T) {
	h := NewHarness(t)
	if diags := checkRetainedWorkflow(h.Root); len(diags) != 0 {
		t.Fatalf("retained workflow is incomplete:\n%s", strings.Join(diags, "\n"))
	}
	runAnchorBites(t, anchorsWithDiagnosticPrefix(retainedWorkflowPrefix), func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}

// TestFreshTicketAuthors grades the fresh ticket author rows. Each row bites alone on a
// synthetic tree, and the shipped guidance satisfies every row. A Forbid row bites when its
// retired sentence returns.
func TestFreshTicketAuthors(t *testing.T) {
	family := anchorsWithDiagnosticPrefix("fresh ticket author: ")
	if len(family) == 0 {
		t.Fatal("no fresh ticket author row is registered")
	}
	runAnchorBites(t, family, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
	requireConformantLiveTree(t, family)
}

// checkRetainedWorkflow reads the guidance at root for every retained-workflow row. The
// registry is the one list: the check takes each row's group and diagnostic from the rows
// whose diagnostic opens with the family prefix, and it holds no expectation of its own.
func checkRetainedWorkflow(root string) []string {
	var groups []anchors.Group
	for _, anchor := range anchorsWithDiagnosticPrefix(retainedWorkflowPrefix) {
		if !slices.Contains(groups, anchor.Group) {
			groups = append(groups, anchor.Group)
		}
	}
	var diags []string
	for _, group := range groups {
		for _, diag := range anchors.EvaluateGroup(root, group) {
			if strings.HasPrefix(diag, retainedWorkflowPrefix) {
				diags = append(diags, diag)
			}
		}
	}
	return diags
}

// TestRetainedWorkflowReadsRegistry proves that the retained-workflow check keeps no second
// list. No string literal of this file equals a registry diagnostic or needle. On a tree that
// contradicts every family row, the check reports each family row once, and Entries() gives
// that count, so a registry edit changes no expectation here.
func TestRetainedWorkflowReadsRegistry(t *testing.T) {
	t.Run("no registry literal", func(t *testing.T) {
		registered := map[string]bool{}
		for _, anchor := range anchors.Entries() {
			registered[anchor.Diagnostic] = true
			registered[anchor.Needle] = true
		}
		file, err := parser.ParseFile(token.NewFileSet(), "retained_workflow_test.go", nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(node ast.Node) bool {
			literal, ok := node.(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				return true
			}
			if value, err := strconv.Unquote(literal.Value); err == nil && value != "" && registered[value] {
				t.Errorf("retained-workflow check holds a literal that the registry also holds: %q", value)
			}
			return true
		})
	})
	t.Run("count from registry", func(t *testing.T) {
		family := anchorsWithDiagnosticPrefix(retainedWorkflowPrefix)
		if len(family) == 0 {
			t.Fatal("no retained-workflow row is registered")
		}
		root := t.TempDir()
		bodies := map[string]*strings.Builder{}
		var files []string
		for _, anchor := range family {
			body, ok := bodies[anchor.File]
			if !ok {
				body = &strings.Builder{}
				body.WriteString("# Fixture\n\n" + plantedAnchorContradiction + "\n\n")
				bodies[anchor.File] = body
				files = append(files, anchor.File)
			}
			if anchor.Section != "" && !strings.Contains(body.String(), "\n## "+anchor.Section+"\n") {
				body.WriteString("## " + anchor.Section + "\n\n" + plantedAnchorContradiction + "\n\n")
			}
		}
		for _, anchor := range family {
			if anchor.Kind == anchors.Forbid {
				bodies[anchor.File].WriteString(anchor.Needle + "\n\n")
			}
		}
		for _, file := range files {
			path := filepath.Join(root, filepath.FromSlash(file))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(bodies[file].String()), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		var want []string
		for _, anchor := range family {
			want = append(want, anchor.Diagnostic)
		}
		got := checkRetainedWorkflow(root)
		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Fatalf("retained-workflow check on a contradicted tree reported %d diagnostics, want the %d registry rows:\ngot:\n%s\nwant:\n%s", len(got), len(want), strings.Join(got, "\n"), strings.Join(want, "\n"))
		}
	})
}

func checkImplementationChunkTable(spec string) []string {
	const heading = "## Implementation chunks"
	start := strings.Index(spec, heading)
	if start < 0 {
		return []string{"retained workflow: implementation chunk table is absent"}
	}
	section := spec[start+len(heading):]
	if end := strings.Index(section, "\n## "); end >= 0 {
		section = section[:end]
	}
	lines := strings.Split(section, "\n")
	header := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "| chunk / ticket | blocked by | delivered outcome | acceptance rows | tests | harder chunk |" {
			header = i
			break
		}
	}
	if header < 0 {
		return []string{"retained workflow: implementation chunk table dropped outcome, acceptance rows, or tests columns"}
	}
	var diags []string
	rows := 0
	for _, line := range lines[header+2:] {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			break
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		if len(cells) != 6 {
			diags = append(diags, "retained workflow: implementation chunk row has the wrong shape")
			continue
		}
		rows++
		for index, label := range map[int]string{0: "stable ID", 2: "outcome", 3: "acceptance rows", 4: "tests"} {
			if strings.TrimSpace(cells[index]) == "" {
				diags = append(diags, fmt.Sprintf("retained workflow: implementation chunk %d dropped %s", rows, label))
			}
		}
	}
	if rows == 0 {
		diags = append(diags, "retained workflow: implementation chunk table has no planned chunks")
	}
	return diags
}

func TestRetainedWorkflowChunkTableBites(t *testing.T) {
	valid := "## Implementation chunks\n\n| chunk / ticket | blocked by | delivered outcome | acceptance rows | tests | harder chunk |\n| --- | --- | --- | --- | --- | --- |\n| 1.md | none | retained behavior | W1 | `TestRetainedWorkflow` | no |\n\n## Testing decisions\n"
	if diags := checkImplementationChunkTable(valid); len(diags) != 0 {
		t.Fatalf("valid implementation chunk table failed: %v", diags)
	}
	for _, tc := range []struct {
		name       string
		from       string
		to         string
		diagnostic string
	}{
		{name: "outcome", from: "retained behavior", to: "", diagnostic: "dropped outcome"},
		{name: "acceptance rows", from: "W1", to: "", diagnostic: "dropped acceptance rows"},
		{name: "tests", from: "`TestRetainedWorkflow`", to: "", diagnostic: "dropped tests"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mutated := strings.Replace(valid, tc.from, tc.to, 1)
			if !containsDiagnostic(checkImplementationChunkTable(mutated), tc.diagnostic) {
				t.Fatalf("chunk-table mutation did not bite: %s", tc.diagnostic)
			}
		})
	}
}
