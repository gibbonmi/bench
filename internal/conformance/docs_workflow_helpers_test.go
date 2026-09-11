package conformance

import (
	"fmt"
	"os"
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

// structuredPhaseClauseOrder is the contract itself, in declaration order.
// .bench/BENCH.md must name exactly these clauses, once each. Reading the required set
// out of the document instead would let a self-consistent deletion or addition define its
// own passing contract.
var structuredPhaseClauseOrder = []string{"progress", "exit", "omission", "cohesion"}

func structuredPhaseClauseRequired(name string) bool {
	return slices.Contains(structuredPhaseClauseOrder, name)
}

func checkWorkflowAnchors(root string) []string {
	diags := anchors.EvaluateGroup(root, anchors.BeforeStructured)
	diags = append(diags, checkStructuredPhaseContract(readIfExists(filepath.Join(root, ".bench", "BENCH.md")))...)
	diags = append(diags, anchors.EvaluateGroup(root, anchors.AfterStructured)...)
	whatNext := readIfExists(filepath.Join(root, ".agents", "commands", "bench-drain.md"))
	diags = append(diags, checkRoadmapContextQuery(whatNext)...)
	diags = append(diags, anchors.EvaluateGroup(root, anchors.AfterRoadmapContext)...)
	diags = append(diags, anchors.EvaluateGroup(root, anchors.AfterImplementSpec)...)
	implementSpec := readIfExists(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md"))
	reviewImplementation := readIfExists(filepath.Join(root, ".agents", "commands", "bench-review-implementation.md"))
	diags = append(diags, checkReviewConvergenceContract(implementSpec, reviewImplementation)...)
	diags = append(diags, checkIntegrationSourceWorkflowCurrency(root)...)
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

	if text := implementSpec; text != "" && !strings.Contains(text, "craft-line") {
		diags = append(diags, "bench-implement-spec.md does not reference craft-line")
	}
	if text := readIfExists(filepath.Join(root, ".agents", "commands", "bench-write-spec.md")); text != "" {
		if !strings.Contains(text, "craft-line") {
			diags = append(diags, "bench-write-spec.md does not reference craft-line")
		}
		if !strings.Contains(text, "implementation-line recommendation contract") {
			diags = append(diags, "bench-write-spec.md does not reference craft-spec's implementation-line recommendation contract")
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

func checkReviewConvergenceContract(implementSpec, reviewImplementation string) []string {
	if implementSpec == "" || reviewImplementation == "" {
		return nil
	}
	normalize := func(text string) string {
		text = strings.ToLower(collapseSpace(stripHTMLComments(text)))
		return strings.NewReplacer("`", "", "<code>", "", "</code>", "").Replace(text)
	}
	implementSpec = normalize(implementSpec)
	reviewImplementation = normalize(reviewImplementation)

	reviewRequirements := []string{
		"each planned chunk takes one review across standards, spec, and coverage",
		"the axes read the whole approved spec and focus on the frozen chunk-base..chunk-tip delta",
		"current repair coverage closes those predicates",
		"repeat delegated review only for a later semantic delta or a cross-chunk concern that invalidates prior evidence",
		"the successor chunk starts only after findings and repair coverage close",
		"after the last chunk, the retained author reconciles overall acceptance and integration before landing",
	}
	for _, requirement := range reviewRequirements {
		if !strings.Contains(reviewImplementation, requirement) {
			return []string{"bench-review-implementation dropped the chunk-review convergence contract: " + requirement}
		}
	}
	for _, requirement := range []string{
		"after the last ticket in a chunk, freeze the chunk delta and run the three review axes before advancing",
		"accepted findings return to the retained author",
		"after the last chunk, reconcile every acceptance row and the integrated behavior",
		"repeat delegated review only when a later delta or cross-chunk concern invalidates prior evidence",
	} {
		if !strings.Contains(implementSpec, requirement) {
			return []string{"bench-implement-spec dropped the retained chunk-review loop: " + requirement}
		}
	}
	return nil
}

func TestReviewConvergenceContractCurrentDocs(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	implementSpec := readIfExists(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md"))
	reviewImplementation := readIfExists(filepath.Join(root, ".agents", "commands", "bench-review-implementation.md"))
	if diags := checkReviewConvergenceContract(implementSpec, reviewImplementation); len(diags) != 0 {
		t.Fatalf("review convergence contract is incomplete:\n%s", strings.Join(diags, "\n"))
	}

	normalizedReview := strings.ToLower(collapseSpace(reviewImplementation))
	mutated := strings.Replace(
		normalizedReview,
		"the successor chunk starts only after findings and repair coverage close",
		"the successor chunk starts before findings and repair coverage close",
		1,
	)
	if mutated == normalizedReview {
		t.Fatal("scope-widening mutation did not apply")
	}
	if !containsDiagnostic(checkReviewConvergenceContract(implementSpec, mutated), "chunk-review convergence contract") {
		t.Fatal("advancing before chunk findings close did not bite")
	}
}

func checkRoadmapContextQuery(whatNext string) []string {
	if whatNext == "" {
		return nil
	}
	indexCall := "`bench roadmap --context`"
	rowCall := "`bench roadmap --context --row <ids>`"
	if strings.Count(whatNext, indexCall) != 1 || strings.Count(whatNext, rowCall) != 1 ||
		strings.Count(whatNext, "bench roadmap --context") != 2 ||
		!strings.Contains(collapseSpace(whatNext), "If the query fails, stop the phase") ||
		!strings.Contains(collapseSpace(whatNext), "manual evidence reconstruction") {
		return []string{"bench-drain dropped the roadmap context query"}
	}
	return nil
}

func TestRoadmapContextQueryCheckBites(t *testing.T) {
	const guidance = "Use `bench roadmap --context` once, then fetch `bench roadmap --context --row <ids>`. If the query fails, stop the phase; manual evidence reconstruction is not a fallback."
	for _, mutation := range []struct{ name, old, replacement string }{
		{"invocation set", "`bench roadmap --context --row <ids>`", "`bench roadmap --context --full`"},
		{"query failure", "If the query fails, stop the phase", "If the query fails, continue"},
		{"manual reconstruction", "manual evidence reconstruction", "ad hoc reconstruction"},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			mutated := strings.Replace(guidance, mutation.old, mutation.replacement, 1)
			if !containsDiagnostic(checkRoadmapContextQuery(mutated), "bench-drain dropped the roadmap context query") {
				t.Fatal("mutation did not bite")
			}
		})
	}
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
	// declaredOrder holds the first occurrence of each required name. The check compares
	// order only after every required name appears. A missing name skips the order
	// comparison, but duplicate and unknown diagnostics can still accompany it.
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
// any self-consistent deletion, duplication, or addition. The expectation must name the
// set the guide is required to carry.
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
			// The rename keeps the declared count at four. A membership-only contract stays green
			// on it. The test must attribute both halves of the rename.
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

// plantedAnchorContradiction replaces a file's conformant text. It matches no anchor
// needle, so it contradicts a Require row and satisfies a Forbid row.
const plantedAnchorContradiction = "planted contradictory clause"

func anchorsWithDiagnosticPrefix(prefix string) []anchors.Anchor {
	var found []anchors.Anchor
	for _, anchor := range anchors.Entries() {
		if strings.HasPrefix(anchor.Diagnostic, prefix) {
			found = append(found, anchor)
		}
	}
	return found
}

// runAnchorBites proves that each anchor in a family bites alone. A subtest writes the
// conformant file, requires silence, plants the contradiction, and requires the
// diagnostic. A Forbid row swaps the pair, because its needle is the contradiction.
// subject names each subtest. Two anchors that share one subject hide each other in the
// run output, so a repeated subject fails the family.
func runAnchorBites(t *testing.T, family []anchors.Anchor, subject func(anchors.Anchor) string) {
	t.Helper()
	seen := map[string]bool{}
	for _, anchor := range family {
		name := subject(anchor)
		if seen[name] {
			t.Fatalf("anchor bite subject repeated: %s", name)
		}
		seen[name] = true
		t.Run(name, func(t *testing.T) {
			conformant, contradictory := anchor.Needle, plantedAnchorContradiction
			if anchor.Kind == anchors.Forbid {
				conformant, contradictory = contradictory, anchor.Needle
			}
			root := t.TempDir()
			path := filepath.Join(root, filepath.FromSlash(anchor.File))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			write := func(text string) {
				t.Helper()
				if anchor.Section != "" {
					text = "# Fixture\n\n## " + anchor.Section + "\n\n" + text
				}
				if err := os.WriteFile(path, []byte(text+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			write(conformant)
			if diags := checkWorkflowAnchors(root); containsDiagnostic(diags, anchor.Diagnostic) {
				t.Fatalf("anchor is red while its file conforms: %s", anchor.Diagnostic)
			}
			write(contradictory)
			if diags := checkWorkflowAnchors(root); !containsDiagnostic(diags, anchor.Diagnostic) {
				t.Fatalf("the contradictory file did not bite with %q", anchor.Diagnostic)
			}
		})
	}
}

const integrationSourceDiagnosticPrefix = "workflow integration source: "

func integrationSourceWorkflowAnchors() []anchors.Anchor {
	return anchorsWithDiagnosticPrefix(integrationSourceDiagnosticPrefix)
}

func checkIntegrationSourceWorkflowCurrency(root string) []string {
	var diags []string
	for _, anchor := range integrationSourceWorkflowAnchors() {
		if anchor.File == "internal/usage/worktree.go" {
			continue
		}
		text := strings.ToLower(anchors.CollapseSpace(anchors.StripHTMLComments(
			readIfExists(filepath.Join(root, filepath.FromSlash(anchor.File))),
		)))
		text = strings.NewReplacer("`", "", "<code>", "", "</code>", "").Replace(text)
		for _, stale := range []string{
			"benchbase",
			"sole landing path",
			"path-scoped bench commit is the only",
			"branch's recorded pre-shift base",
		} {
			if strings.Contains(text, stale) {
				diags = append(diags, fmt.Sprintf("%s retains stale scalar or sole-path workflow claim %q", anchor.File, stale))
			}
		}
	}
	return diags
}

func TestIntegrationSourceWorkflowAnchorsBiteIndependently(t *testing.T) {
	workflowAnchors := integrationSourceWorkflowAnchors()
	if got, want := len(workflowAnchors), 10; got != want {
		t.Fatalf("integration-source workflow anchor count = %d, want %d", got, want)
	}
	runAnchorBites(t, workflowAnchors, func(anchor anchors.Anchor) string { return anchor.File })
}
