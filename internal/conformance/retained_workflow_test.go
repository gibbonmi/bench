package conformance

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/lines"
)

func TestRetainedWorkflow(t *testing.T) {
	h := NewHarness(t)
	if diags := checkRetainedWorkflow(h.Root); len(diags) != 0 {
		t.Fatalf("retained workflow is incomplete:\n%s", strings.Join(diags, "\n"))
	}
	runAnchorBites(t, anchorsWithDiagnosticPrefix("retained workflow: "), func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}

func checkRetainedWorkflow(root string) []string {
	var diags []string
	binding := lines.ParseBinding([]byte(readIfExists(filepath.Join(root, ".bench", "lines.env"))))
	if got, want := binding.Cell("codex", "mid"), "gpt-5.6-sol"; got != want {
		diags = append(diags, fmt.Sprintf("Codex mid binding = %q, want %q", got, want))
	}
	if got, want := binding.Cell("codex", "top"), "gpt-6-astra"; got != want {
		diags = append(diags, fmt.Sprintf("Codex top binding = %q, want %q", got, want))
	}
	if got, want := binding.Cell("codex", "cheap"), "gpt-5.6-terra"; got != want {
		diags = append(diags, fmt.Sprintf("Codex cheap binding = %q, want %q", got, want))
	}
	wantDiagnostics := []string{
		"retained workflow: craft-spec dropped the one implementation-line template",
		"retained workflow: craft-spec dropped the implementation-line reason factors",
		"retained workflow: craft-spec dropped the harder-chunk annotation",
		"retained workflow: craft-line dropped the conditional review stage default",
		"retained workflow: craft-line dropped the Codex mid-to-top review branch",
		"retained workflow: craft-line dropped the default mid review branch",
		"retained workflow: review phase dropped the conditional review-line reference",
		"retained workflow: project profile dropped the conditional review-line reference",
		"retained workflow: craft-line dropped the user-directed model-switch boundary",
		"retained workflow: operating guide dropped retained implementation authorship",
		"retained workflow: operating guide dropped chunk outcome and review checkpoint",
		"retained workflow: operating guide dropped serial ticket checkpoint",
		"retained workflow: operating guide allows advancement before three-axis chunk review",
		"retained workflow: operating guide dropped final acceptance reconciliation",
		"retained workflow: operating guide restored routine final full review",
		"retained workflow: operating guide dropped diagnostic-only consultation boundary",
		"retained workflow: implementation phase dropped retained policy reference",
		"retained workflow: drain restored delegated batch authorship",
		"retained workflow: field guide dropped retained chunk review",
		"retained workflow: craft-spec dropped the canonical plan-expansion owner",
		"retained workflow: operating guide dropped changed-chunk identity mapping",
		"retained workflow: operating guide dropped plan-update timing or preserved guarantees",
		"retained workflow: operating guide dropped expansion learning and drain ownership",
		"retained workflow: implementation phase dropped plan-expansion timing",
		"retained workflow: craft-gate dropped approved expansion timing",
		"retained workflow: craft-tickets restored Writes as an approval boundary",
		"retained workflow: delegation discipline restored Writes as a refusal boundary",
		"retained workflow: operating guide dropped the delegated opt-in entry",
		"retained workflow: operating guide dropped the delegated prerequisite-checkpoint wait",
		"retained workflow: operating guide dropped the green predecessor dispatch rule",
		"retained workflow: operating guide dropped the pending-or-red predecessor stop",
		"retained workflow: operating guide dropped delegated ticket-author repair ownership",
		"retained workflow: craft-line dropped the delegated tier authorization",
		"retained workflow: craft-line dropped the unbound delegated model stop",
		"retained workflow: craft-line dropped the active-writer author limit",
		"retained workflow: craft-line dropped the delegated mid review route",
		"retained workflow: craft-delegate dropped the delegated per-ticket author",
		"retained workflow: delegation discipline dropped the no-progress transfer trigger",
		"retained workflow: delegation discipline dropped the terminal-failure transfer trigger",
		"retained workflow: delegation discipline dropped the cap-exhaustion transfer trigger",
		"retained workflow: delegation discipline dropped the lost-session transfer trigger",
		"retained workflow: delegation discipline dropped the confirmed writer termination",
		"retained workflow: implementation phase dropped the delegated entry refusals",
		"retained workflow: implementation phase dropped the delegated dispatch declaration",
		"retained workflow: implementation phase dropped delegated resumption contents",
		"retained workflow: review phase dropped the integrated chunk-tip review fence",
		"retained workflow: review phase dropped the delegated axis exclusions",
		"retained workflow: final check dropped the delegated account reconciliation",
		"retained workflow: final check dropped the no-paid-comparison boundary",
		"retained workflow: final check dropped the orchestrator final verification",
		"retained workflow: craft-tickets dropped the delegated serial ticket checkpoint",
	}
	wantPredicates := map[string]struct {
		file    string
		section string
		needle  string
	}{
		"retained workflow: craft-spec dropped the implementation-line reason factors": {
			file:    ".agents/skills/bench-craft-spec/SKILL.md",
			section: "Template",
			needle:  "Implementation-line reason: <hardest material chunk, spec precision, seam uncertainty, and test strength>.",
		},
		"retained workflow: craft-line dropped the Codex mid-to-top review branch": {
			file:   ".agents/skills/bench-craft-line/SKILL.md",
			needle: "A Codex mid implementation sends each axis to the Codex top binding.",
		},
		"retained workflow: craft-line dropped the default mid review branch": {
			file:   ".agents/skills/bench-craft-line/SKILL.md",
			needle: "Every other implementation sends each axis to the invoking harness's mid binding.",
		},
		"retained workflow: craft-line dropped the user-directed model-switch boundary": {
			file:   ".agents/skills/bench-craft-line/SKILL.md",
			needle: "A different implementation model or session requires user direction. The author can adjust effort in the retained session and reports the change.",
		},
		"retained workflow: operating guide dropped retained implementation authorship": {
			file:   ".bench/BENCH.md",
			needle: "The retained implementation session writes production changes, tests, probes, and repairs.",
		},
		"retained workflow: operating guide dropped chunk outcome and review checkpoint": {
			file:   ".bench/BENCH.md",
			needle: "An implementation chunk is one coherent behavior outcome with acceptance rows, tests, and a review checkpoint.",
		},
		"retained workflow: operating guide dropped serial ticket checkpoint": {
			file:   ".bench/BENCH.md",
			needle: "A ticket remains a serial green commit checkpoint.",
		},
		"retained workflow: operating guide allows advancement before three-axis chunk review": {
			file:   ".bench/BENCH.md",
			needle: "After a chunk's ticket commits, freeze its delta and run Standards, Spec, and Coverage against the whole approved spec before starting its successor.",
		},
		"retained workflow: operating guide dropped final acceptance reconciliation": {
			file:   ".bench/BENCH.md",
			needle: "After the last chunk, the retained author reconciles overall acceptance and integration.",
		},
		"retained workflow: operating guide restored routine final full review": {
			file:   ".bench/BENCH.md",
			needle: "Repeat delegated review only for a later delta or a cross-chunk concern that invalidates prior evidence.",
		},
		"retained workflow: operating guide dropped diagnostic-only consultation boundary": {
			file:   ".bench/BENCH.md",
			needle: "A brief read-only diagnostic consultation can inspect evidence, but it receives no implementation or repair assignment.",
		},
		"retained workflow: implementation phase dropped retained policy reference": {
			file:   ".agents/commands/bench-implement-spec.md",
			needle: "Follow `.bench/BENCH.md`'s retained implementation and chunk-review policy.",
		},
		"retained workflow: drain restored delegated batch authorship": {
			file:   ".agents/commands/bench-drain.md",
			needle: "If tracked changes remain, the retained drain session authors the complete tracked batch.",
		},
		"retained workflow: field guide dropped retained chunk review": {
			file:   "docs/field-guide.html",
			needle: "One implementation session retains authorship through the approved ticket graph and its chunk reviews.",
		},
		"retained workflow: craft-spec dropped the canonical plan-expansion owner": {
			file:    ".agents/skills/bench-craft-spec/SKILL.md",
			section: "Slicing a build for delegates",
			needle:  "During a build, `.bench/BENCH.md` owns approved in-scope plan expansion.",
		},
		"retained workflow: operating guide dropped changed-chunk identity mapping": {
			file:   ".bench/BENCH.md",
			needle: "When chunk boundaries change, record old-to-new stable IDs.",
		},
		"retained workflow: operating guide dropped plan-update timing or preserved guarantees": {
			file:   ".bench/BENCH.md",
			needle: "Before using a plan expansion, update the affected spec and tickets; preserve acceptance coverage, dependencies, review checkpoints, existing checks, pass criteria, and required behavior.",
		},
		"retained workflow: operating guide dropped expansion learning and drain ownership": {
			file:   ".bench/BENCH.md",
			needle: "Record every plan or gate expansion with `bench learning`, including what changed, why, and verification; `/bench-drain` owns its later disposition.",
		},
		"retained workflow: implementation phase dropped plan-expansion timing": {
			file:   ".agents/commands/bench-implement-spec.md",
			needle: "When evidence requires an in-scope plan, `Writes:`, or gate expansion, apply `.bench/BENCH.md`'s approved plan-expansion policy before using it.",
		},
		"retained workflow: craft-gate dropped approved expansion timing": {
			file:   ".agents/skills/bench-craft-gate/SKILL.md",
			needle: "An in-scope gate addition during implementation follows `.bench/BENCH.md`'s approved plan-expansion policy before the author uses it.",
		},
		"retained workflow: craft-tickets restored Writes as an approval boundary": {
			file:   ".agents/skills/bench-craft-tickets/SKILL.md",
			needle: "`Writes:` predicts the touched paths; `.bench/BENCH.md` owns how the retained author updates that expectation before an approved in-scope expansion is used.",
		},
		"retained workflow: delegation discipline restored Writes as a refusal boundary": {
			file:    ".agents/skills/bench-craft-delegate/references/delegation-discipline.md",
			section: "In the charge",
			needle:  "A user-directed write delegate treats `Writes:` as an expectation.",
		},
		"retained workflow: operating guide dropped the delegated opt-in entry": {
			file:   ".bench/BENCH.md",
			needle: "`--delegate` applies only to an approved `$bench-implement-spec --full <spec>` run with an approved ticket graph.",
		},
		"retained workflow: operating guide dropped the delegated prerequisite-checkpoint wait": {
			file:   ".bench/BENCH.md",
			needle: "A ticket in a dependent chunk waits for every prerequisite chunk checkpoint.",
		},
		"retained workflow: operating guide dropped the green predecessor dispatch rule": {
			file:   ".bench/BENCH.md",
			needle: "A same-chunk successor ticket starts after its predecessor ticket commits green.",
		},
		"retained workflow: operating guide dropped the pending-or-red predecessor stop": {
			file:   ".bench/BENCH.md",
			needle: "A pending or red predecessor commit stops that successor dispatch.",
		},
		"retained workflow: review phase dropped the integrated chunk-tip review fence": {
			file:   ".agents/commands/bench-review-implementation.md",
			needle: "A delegated chunk review starts after every ticket of the chunk reaches the integrated chunk tip.",
		},
		"retained workflow: craft-line dropped the delegated mid review route": {
			file:   ".agents/skills/bench-craft-line/SKILL.md",
			needle: "Every delegated review axis uses the invoking harness's configured mid binding at high effort.",
		},
		"retained workflow: delegation discipline dropped the confirmed writer termination": {
			file:    ".agents/skills/bench-craft-delegate/references/delegation-discipline.md",
			section: "Delegated author transfer",
			needle:  "Every author transfer waits for confirmed termination of the old writer.",
		},
	}
	family := anchorsWithDiagnosticPrefix("retained workflow: ")
	for _, want := range wantDiagnostics {
		index := slices.IndexFunc(family, func(anchor anchors.Anchor) bool { return anchor.Diagnostic == want })
		if index < 0 {
			diags = append(diags, "retained-workflow anchor is absent: "+want)
			continue
		}
		if predicate, ok := wantPredicates[want]; ok {
			anchor := family[index]
			if anchor.File != predicate.file || anchor.Section != predicate.section || anchor.Needle != predicate.needle {
				diags = append(diags, fmt.Sprintf("retained-workflow predicate drifted for %q: %#v", want, anchor))
			}
		}
	}
	for file, forbidden := range map[string][]string{
		".agents/commands/bench-implement-spec.md": {
			"Every spec-backed run assigns genuine write work to a write subagent",
		},
		".agents/commands/bench-write-spec.md": {
			"A fresh session\nbuilds after ticket approval",
			"fresh mid-tier build session",
		},
		".agents/commands/bench-drain.md": {
			"one later write delegate authors the complete tracked batch",
			"one later write delegate to draft the complete tracked pass",
		},
		".agents/commands/bench-review-implementation.md": {
			"Initial review blocks on the full",
			"repair-scoped re-review",
		},
		".agents/skills/bench-craft-tickets/SKILL.md": {
			"one fresh write-delegate charge",
		},
		".agents/skills/bench-craft-delegate/SKILL.md": {
			"All other code authorship\nruns as a write-delegation",
			"Otherwise a fresh charge in an isolated worktree carries the finding",
		},
		".agents/skills/bench-craft-spec/SKILL.md": {
			"A build may not edit its own spec's acceptance rows",
		},
		".agents/skills/bench-craft-delegate/references/delegation-discipline.md": {
			"takes a\n  fence extension in a continuation",
		},
		"docs/field-guide.html": {
			"their build starts in a fresh session",
			"per fresh write-delegate context",
		},
		"projects/benchkit.md": {
			"After ticket approval, a fresh\n  mid-tier session starts the build.",
		},
	} {
		text := readIfExists(filepath.Join(root, filepath.FromSlash(file)))
		for _, phrase := range forbidden {
			if strings.Contains(text, phrase) {
				diags = append(diags, fmt.Sprintf("live workflow reader %s retains stale implementation instruction %q", file, phrase))
			}
		}
	}
	return diags
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
