// Tests that the craft-tickets skill example parses clean under the live grammar.
package tickets

import (
	"os"
	"strings"
	"testing"
)

const (
	skillPath    = "../../.agents/skills/bench-craft-tickets/SKILL.md"
	exampleBegin = "<!-- ticket-example:begin -->"
	exampleEnd   = "<!-- ticket-example:end -->"
	exampleName  = "render-cancelled-jobs-in-status.md"
	exampleTag   = "CJ"
)

// TestCraftTicketsExampleParsesClean feeds the craft-tickets skill's own marked
// example to the live parser. The example is the starting point a cold author
// copies, so a diagnostic here reds the first ticket that author writes.
func TestCraftTicketsExampleParsesClean(t *testing.T) {
	content, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read skill: %v", err)
	}
	example := markedExample(t, string(content))
	siblings := []string{exampleName, "parse-cancelled-job-records.md"}
	ticket, diagnostics := ParseTicket(exampleName, []byte(example), siblings, exampleTag)
	if len(diagnostics) != 0 {
		t.Fatalf("example diagnostics = %v, want none", diagnostics)
	}
	if ticket.Title == "" {
		t.Fatalf("example parsed no title")
	}
}

// markedExample returns the ticket body between the example markers, with the
// enclosing markdown fence removed so the parser reads live grammar.
func markedExample(t *testing.T, content string) string {
	t.Helper()
	_, after, found := strings.Cut(content, exampleBegin)
	if !found {
		t.Fatalf("skill file has no %s marker", exampleBegin)
	}
	block, _, found := strings.Cut(after, exampleEnd)
	if !found {
		t.Fatalf("skill file has no %s marker", exampleEnd)
	}
	var body []string
	for _, line := range strings.Split(strings.TrimSpace(block), "\n") {
		if strings.HasPrefix(line, "```") {
			continue
		}
		body = append(body, line)
	}
	if len(body) == 0 {
		t.Fatalf("marked example is empty")
	}
	return strings.Join(body, "\n") + "\n"
}
