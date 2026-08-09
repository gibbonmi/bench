package conformance

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

const (
	structuredPhaseDeclaration = "- **Structured Bench phase conversation:**"
	structuredPhaseUnavailable = ".bench/BENCH.md cannot verify the structured Bench phase contract because shared rules are missing or empty"
)

// structuredPhaseClauseOrder is the contract itself, in declaration order: .bench/BENCH.md
// must name exactly these clauses, once each. Reading the required set out of the document
// instead lets a self-consistent deletion or addition define its own passing contract.
var structuredPhaseClauseOrder = []string{"progress", "exit", "omission", "cohesion"}

func structuredPhaseClauseRequired(name string) bool {
	return slices.Contains(structuredPhaseClauseOrder, name)
}

func checkWorkflowAnchors(root string) []string {
	diags := anchors.EvaluateGroup(root, anchors.BeforeStructured)
	diags = append(diags, checkStructuredPhaseContract(readIfExists(filepath.Join(root, ".bench", "BENCH.md")))...)
	diags = append(diags, anchors.EvaluateGroup(root, anchors.AfterStructured)...)
	whatNext := readIfExists(filepath.Join(root, ".agents", "commands", "bench-what-next.md"))
	if whatNext != "" && (strings.Count(whatNext, "bench roadmap --context") != 1 ||
		!strings.Contains(collapseSpace(whatNext), "If the query fails, stop the phase") ||
		!strings.Contains(collapseSpace(whatNext), "manual evidence reconstruction")) {
		diags = append(diags, "bench-what-next dropped the roadmap context query")
	}
	diags = append(diags, anchors.EvaluateGroup(root, anchors.AfterRoadmapContext)...)
	implementSpec := strings.ToLower(collapseSpace(stripHTMLComments(readIfExists(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md")))))
	implementSpec = strings.ReplaceAll(implementSpec, "`", "")
	if strings.Contains(implementSpec, "for an accepted repair finding, the coordinator may instead write the repair directly to the working branch before promote") {
		diags = append(diags, "bench-implement-spec permits an accepted repair to bypass provisional assignment and write directly to the working branch")
	}
	if strings.Contains(implementSpec, "not parallelizable") {
		diags = append(diags, "bench-implement-spec permits a generic unused-slot reason outside the closed set")
	}
	diags = append(diags, anchors.EvaluateGroup(root, anchors.AfterImplementSpec)...)
	diags = append(diags, checkSpecAuthorizationContract(root)...)
	diags = append(diags, anchors.EvaluateGroup(root, anchors.AfterSpecAuthorization)...)
	readme := readIfExists(filepath.Join(root, "README.md"))
	if readme != "" {
		if !strings.Contains(readme, "session-start.sh") {
			diags = append(diags, "README layout omits .bench/hooks/session-start.sh")
		}
		if !strings.Contains(readme, "bench.sh") {
			diags = append(diags, "README layout omits the real bin/bench.sh filename")
		}
		if !strings.Contains(readme, "benchkit.md") {
			diags = append(diags, "README layout omits projects/benchkit.md")
		}
		if strings.Contains(readme, "\u2502   \u2514\u2500\u2500 bench                 #") {
			diags = append(diags, "README layout still names bin/bench instead of bin/bench.sh")
		}
	}

	if text := readIfExists(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md")); text != "" && !strings.Contains(text, "craft-line") {
		diags = append(diags, "bench-implement-spec.md does not reference craft-line")
	}
	if text := readIfExists(filepath.Join(root, ".agents", "commands", "bench-write-spec.md")); text != "" {
		if !strings.Contains(text, "craft-line") {
			diags = append(diags, "bench-write-spec.md does not reference craft-line")
		}
		if !strings.Contains(text, "model and effort") {
			diags = append(diags, "bench-write-spec.md does not mandate per-story model and effort")
		}
		if strings.Count(text, "`Bootstrap authority before execution` rule") != 2 {
			diags = append(diags, "bench-write-spec.md does not apply craft-spec's named bootstrap-authority rule during edge walking and falsification")
		}
	}
	if text := readIfExists(filepath.Join(root, ".agents", "skills", "bench-craft-spec", "SKILL.md")); text != "" {
		if !strings.Contains(text, "## Bootstrap authority before execution") {
			diags = append(diags, ".agents/skills/bench-craft-spec/SKILL.md dropped the bootstrap-authority pre-execution trace")
		} else if !strings.Contains(text, "before launching the next executable") || strings.Contains(text, "after launching the next executable") {
			diags = append(diags, ".agents/skills/bench-craft-spec/SKILL.md validates a bootstrap authority after launch")
		}
	}
	if text := readIfExists(filepath.Join(root, ".bench", "BENCH-reference.md")); text != "" && !strings.Contains(text, "BENCH_MODEL") {
		diags = append(diags, "BENCH-reference.md adapter contract does not document BENCH_MODEL")
	}
	return diags
}

func checkStructuredPhaseContract(sharedRules string) []string {
	if strings.TrimSpace(sharedRules) == "" {
		return []string{structuredPhaseUnavailable}
	}
	section := anchors.MarkdownH2Section(stripHTMLComments(sharedRules), "How to talk to me")
	if section == "" {
		return []string{".bench/BENCH.md dropped the structured Bench phase contract from the active How to talk to me section"}
	}

	lines := strings.Split(section, "\n")
	start := -1
	for i, line := range lines {
		if strings.HasPrefix(line, structuredPhaseDeclaration) {
			start = i
			break
		}
	}
	if start < 0 {
		return []string{".bench/BENCH.md dropped the structured Bench phase contract declaration from the active How to talk to me section"}
	}
	declarationBody := strings.TrimSpace(strings.TrimPrefix(lines[start], structuredPhaseDeclaration))
	if declarationBody == "" || structuredPhaseClauseIsNegated(declarationBody) {
		return []string{".bench/BENCH.md negated or emptied the structured Bench phase contract declaration in the active How to talk to me section"}
	}

	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "- ") {
			end = i
			break
		}
	}
	block := lines[start:end]
	bodies, bodyDiags := structuredPhaseClauseBodies(block)
	diags := append([]string(nil), bodyDiags...)
	seen := map[string]bool{}
	var declaredOrder []string
	for _, name := range structuredPhaseClauseNames(block) {
		switch {
		case seen[name]:
			diags = append(diags, fmt.Sprintf(".bench/BENCH.md structured Bench phase contract declares clause %q more than once", name))
		case !structuredPhaseClauseRequired(name):
			seen[name] = true
			diags = append(diags, fmt.Sprintf(".bench/BENCH.md structured Bench phase contract declares unknown clause %q", name))
		default:
			seen[name] = true
			declaredOrder = append(declaredOrder, name)
		}
	}
	// declaredOrder holds the first occurrence of each required name, so order is
	// compared only once every required name appears; a missing name skips the
	// comparison, while duplicate and unknown diagnostics can accompany it.
	if len(declaredOrder) == len(structuredPhaseClauseOrder) {
		for i, name := range declaredOrder {
			if name != structuredPhaseClauseOrder[i] {
				diags = append(diags, fmt.Sprintf(
					".bench/BENCH.md structured Bench phase contract declares clause %q out of contract order; expected %q at position %d",
					name, structuredPhaseClauseOrder[i], i+1))
				break
			}
		}
	}
	for _, name := range structuredPhaseClauseOrder {
		if !seen[name] {
			diags = append(diags, fmt.Sprintf(".bench/BENCH.md structured Bench phase contract does not declare the %s clause", name))
			continue
		}
		body := strings.TrimSpace(bodies[name])
		if body == "" || structuredPhaseClauseIsNegated(body) {
			diags = append(diags, fmt.Sprintf(".bench/BENCH.md dropped the structured Bench phase %s clause", name))
		}
	}
	for name := range bodies {
		if !seen[name] {
			diags = append(diags, fmt.Sprintf(".bench/BENCH.md structured Bench phase clause %q is not named by its contract declaration", name))
		}
	}
	return diags
}

func structuredPhaseClauseNames(block []string) []string {
	var names []string
	for _, line := range block[1:] {
		if strings.HasPrefix(line, "  - **") {
			break
		}
		for {
			start := strings.Index(line, "`")
			if start < 0 {
				break
			}
			line = line[start+1:]
			end := strings.Index(line, "`")
			if end < 0 {
				break
			}
			name := strings.ToLower(strings.TrimSpace(line[:end]))
			if name != "" {
				names = append(names, name)
			}
			line = line[end+1:]
		}
	}
	return names
}

func structuredPhaseClauseBodies(block []string) (map[string]string, []string) {
	bodies := map[string]string{}
	var diags []string
	current := ""
	for _, line := range block[1:] {
		if strings.HasPrefix(line, "  - **") {
			rest := strings.TrimPrefix(line, "  - **")
			labelEnd := strings.Index(rest, ":**")
			if labelEnd < 0 {
				current = ""
				continue
			}
			current = strings.ToLower(strings.TrimSpace(rest[:labelEnd]))
			if _, exists := bodies[current]; exists {
				diags = append(diags, fmt.Sprintf(".bench/BENCH.md structured Bench phase clause %q appears more than once", current))
			}
			bodies[current] = strings.TrimSpace(rest[labelEnd+3:])
			continue
		}
		if current != "" && strings.HasPrefix(line, "    ") {
			bodies[current] = strings.TrimSpace(bodies[current] + " " + strings.TrimSpace(line))
		}
	}
	return bodies, diags
}

func structuredPhaseClauseIsNegated(body string) bool {
	lower := strings.ToLower(strings.TrimSpace(body))
	for _, prefix := range []string{"do not ", "don't ", "never ", "not ", "no "} {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

var (
	markdownH2Sections = anchors.MarkdownH2Sections
	collapseSpace      = anchors.CollapseSpace
	stripHTMLComments  = anchors.StripHTMLComments
)

// TestStructuredPhaseContractPinsTheFixedClauseSet enumerates the four required clause
// names independently of the parser. A contract read out of the document stays green for
// any self-consistent deletion, duplication, or addition, so the expectation has to name
// the set the guide is required to carry.
func TestStructuredPhaseContractPinsTheFixedClauseSet(t *testing.T) {
	bodyLines := map[string]string{
		"progress": "  - **Progress:** Use compact bold **Status:** and **Next:** labels.",
		"exit":     "  - **Exit:** A phase exit leads with `## Result`.",
		"omission": "  - **Omission:** Omit empty progress groups and exit sections.",
		"cohesion": "  - **Cohesion:** Keep related sentences together.",
		"cadence":  "  - **Cadence:** Check in once per iteration.",
	}
	guide := func(declared, bodied []string) string {
		quoted := make([]string, 0, len(declared))
		for _, name := range declared {
			quoted = append(quoted, "`"+name+"`")
		}
		lines := []string{
			"# Bench Operating Guide",
			"",
			"## How to talk to me",
			"",
			structuredPhaseDeclaration + " Apply the named clauses",
			"  " + strings.Join(quoted, ", ") + " proportionally.",
		}
		for _, name := range bodied {
			lines = append(lines, bodyLines[name])
		}
		return strings.Join(lines, "\n") + "\n"
	}
	complete := []string{"progress", "exit", "omission", "cohesion"}
	without := func(dropped string) []string {
		kept := make([]string, 0, len(complete)-1)
		for _, name := range complete {
			if name != dropped {
				kept = append(kept, name)
			}
		}
		return kept
	}

	if diags := checkStructuredPhaseContract(guide(complete, complete)); len(diags) != 0 {
		t.Fatalf("the complete clause set failed its own contract:\n%s", strings.Join(diags, "\n"))
	}

	cases := []struct {
		name     string
		declared []string
		bodied   []string
		want     []string
	}{
		{"missing progress", without("progress"), without("progress"), []string{`does not declare the progress clause`}},
		{"missing exit", without("exit"), without("exit"), []string{`does not declare the exit clause`}},
		{"missing omission", without("omission"), without("omission"), []string{`does not declare the omission clause`}},
		{"missing cohesion", without("cohesion"), without("cohesion"), []string{`does not declare the cohesion clause`}},
		{
			"no declared names",
			nil,
			nil,
			[]string{
				`does not declare the progress clause`,
				`does not declare the exit clause`,
				`does not declare the omission clause`,
				`does not declare the cohesion clause`,
			},
		},
		{
			"duplicate name",
			[]string{"progress", "exit", "exit", "omission", "cohesion"},
			complete,
			[]string{`declares clause "exit" more than once`},
		},
		{
			"unknown name",
			[]string{"progress", "exit", "omission", "cohesion", "cadence"},
			[]string{"progress", "exit", "omission", "cohesion", "cadence"},
			[]string{`declares unknown clause "cadence"`},
		},
		{
			// The rename keeps the declared count at four, so a membership-only contract
			// stays green on it; both halves of the rename have to be attributed.
			"renamed progress to cadence",
			[]string{"cadence", "exit", "omission", "cohesion"},
			[]string{"cadence", "exit", "omission", "cohesion"},
			[]string{`does not declare the progress clause`, `declares unknown clause "cadence"`},
		},
		{
			"reordered declaration",
			[]string{"exit", "progress", "omission", "cohesion"},
			complete,
			[]string{`declares clause "exit" out of contract order; expected "progress" at position 1`},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			diags := checkStructuredPhaseContract(guide(testCase.declared, testCase.bodied))
			for _, want := range testCase.want {
				if !containsDiagnostic(diags, want) {
					t.Fatalf("clause set stayed valid; want %q in diagnostics:\n%s", want, strings.Join(diags, "\n"))
				}
			}
		})
	}
}
