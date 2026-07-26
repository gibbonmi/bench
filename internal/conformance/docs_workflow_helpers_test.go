package conformance

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	structuredPhaseDeclaration = "- **Structured Bench phase conversation:**"
	structuredPhaseUnavailable = ".bench/BENCH.md cannot verify the structured Bench phase contract because shared rules are missing or empty"
)

func checkWorkflowAnchors(root string) []string {
	var diags []string
	require := func(rel, needle string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if !exists(path) {
			diags = append(diags, "acceptance coverage anchor file missing: "+rel)
			return
		}
		if !strings.Contains(readIfExists(path), needle) {
			diags = append(diags, fmt.Sprintf("%s missing acceptance coverage anchor: %s", rel, needle))
		}
	}

	require(".agents/commands/bench-write-spec.md", "acceptance coverage map")
	require(".agents/commands/bench-write-spec.md", "why it catches the failure")
	require(".agents/commands/bench-write-spec.md", "red signal")
	require(".agents/skills/bench-craft-spec/SKILL.md", "why it catches the failure")
	require(".agents/skills/bench-craft-spec/SKILL.md", "re-run idempotency")
	require(".agents/skills/bench-craft-spec/SKILL.md", "separate capability")
	require(".agents/commands/bench-write-spec.md", "bench-craft-spec")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "bench-craft-spec")
	require(".agents/skills/bench-craft-review/SKILL.md", "bench-craft-spec")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "acceptance row")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "not TDD-able")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "call count")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "row schema and the red-signal definition are")
	require(".agents/commands/bench-implement-spec.md", "coverage table")
	require(".agents/commands/bench-implement-spec.md", "already covered")
	require(".agents/commands/bench-implement-spec.md", "turning red-to-green")
	require(".agents/commands/bench-implement-spec.md", "bench coverage <spec>")
	require(".agents/commands/bench-review-implementation.md", "acceptance coverage map")
	require(".agents/commands/bench-review-implementation.md", "mapped behavior")
	require(".agents/commands/bench-review-implementation.md", "bench diff --full")
	require(".agents/commands/bench-review-implementation.md", "bench diff --full --commit")
	require(".agents/commands/bench-final-check.md", ".bench/gate.sh")
	require(".agents/commands/bench-final-check.md", "BENCH_GATE")
	require(".agents/commands/bench-write-spec.md", "seam diagram")
	require(".agents/commands/bench-write-spec.md", "tests attach here")
	require(".agents/commands/bench-write-spec.md", "edge inventory")
	require(".agents/commands/bench-write-spec.md", "Won't handle")
	require(".agents/commands/bench-write-spec.md", "hostile-input checklist")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "floor, not the ceiling")
	require(".agents/skills/bench-craft-seams/SKILL.md", "failure modes")
	require(".agents/skills/bench-craft-seams/SKILL.md", "structure.budgets")
	require(".agents/commands/bench-review-implementation.md", "## Coverage")
	require(".agents/commands/bench-review-implementation.md", "Coverage axis")
	require(".agents/commands/bench-setup-repo.md", "hostile-input checklist")
	require("projects/benchkit.md", "hostile-input checklist")
	require(".agents/commands/bench-setup-repo.md", "craft-gate")
	require(".agents/commands/bench-final-check.md", "craft-gate")
	require(".agents/commands/bench-review-implementation.md", "craft-review")
	require(".agents/skills/bench-craft-review/SKILL.md", "an edge nobody decided")
	require(".agents/commands/bench-review-implementation.md", "craft-delegate")
	require(".agents/skills/bench-craft-delegate/SKILL.md", "a claim, not a result")
	require(".agents/skills/bench-craft-delegate/SKILL.md", "bench worktree create --request")
	require(".agents/skills/bench-craft-delegate/SKILL.md", "git stash")
	require(".agents/commands/bench-implement-spec.md", "When the build stops short")
	require(".agents/commands/bench-write-spec.md", "Superseded by")
	require(".agents/commands/bench-shape-idea.md", "## Handoff")
	require(".agents/commands/bench-shape-idea.md", "Hostile-input owner")
	require(".agents/commands/bench-shape-idea.md", "Dependency order")
	require(".agents/commands/bench-shape-idea.md", "n/a \u2014")
	require(".agents/commands/bench-write-spec.md", "map's Handoff")
	require(".agents/commands/bench-write-spec.md", "spec-retire:")
	require(".agents/commands/bench-write-spec.md", "Status: staged")
	require(".agents/commands/bench-write-spec.md", "new session on the mid tier")
	require(".agents/commands/bench-write-spec.md", "no map backs the draft")
	require(".agents/commands/bench-write-spec.md", "written in the same session as the draft")
	require(".agents/commands/bench-write-spec.md", "mostly not observed reds")
	require(".agents/commands/bench-write-spec.md", "runs at the mid tier")
	require(".agents/commands/bench-write-spec.md", "Every draft gets the pass")
	require("projects/benchkit.md", "Spec falsification pass")
	require(".agents/commands/bench-implement-spec.md", "Status: implemented")
	require(".agents/commands/bench-debug.md", "diff-filter=D")
	require(".agents/commands/bench-review-implementation.md", "reviews/<spec-slug>.md")
	require(".agents/commands/bench-review-implementation.md", "same session that writes it")
	require(".agents/commands/bench-implement-spec.md", "reviews/<spec-slug>.md")
	require(".agents/commands/bench-implement-spec.md", "names and deletes the file")
	require(".agents/commands/bench-final-check.md", "not outlive the decision it captured")
	require(".agents/commands/bench-implement-spec.md", "bench commit -m")
	require(".agents/commands/bench-final-check.md", "bench commit -m")
	require(".agents/commands/bench-review-implementation.md", "actionable findings")
	require(".agents/commands/bench-review-implementation.md", "writes no artifact")
	require(".agents/commands/bench-review-implementation.md", "same green fix commit")
	require(".agents/commands/bench-review-implementation.md", "section per axis")
	require(".agents/commands/bench-review-implementation.md", "finding count, its worst issue")
	require(".agents/commands/bench-review-implementation.md", "doc citation its axis supplied")
	require(".agents/commands/bench-what-next.md", "Reconcile first")
	require(".agents/commands/bench-what-next.md", "through the accused command")
	require(".agents/commands/bench-debug.md", "through the accused command")
	require(".agents/commands/bench-what-next.md", "empties to zero")
	require(".agents/commands/bench-what-next.md", "verdict in the batch diff")
	require(".agents/commands/bench-what-next.md", "one uncommitted batch diff")
	require(".agents/commands/bench-what-next.md", "commit on green")
	require(".agents/commands/bench-what-next.md", "## Recommended sequence")

	require(".agents/commands/bench-assess.md", "verify the previous assessment's backlog landed")
	require(".agents/commands/bench-assess.md", "read-only area sweeps on the mid tier")
	require(".agents/commands/bench-assess.md", "synthesize adversarially on the top tier")
	require(".agents/commands/bench-assess.md", "replaces its predecessor")
	require(".agents/commands/bench-assess.md", "/bench-what-next")

	requireCollapsed := func(rel, needle, diag string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if !exists(path) {
			diags = append(diags, "acceptance coverage anchor file missing: "+rel)
			return
		}
		if !strings.Contains(collapseSpace(readIfExists(path)), needle) {
			diags = append(diags, diag)
		}
	}
	// forbid is the must-not-contain half of an anchor pair: a workflow agreement
	// with one owner stays repaired only while the contradicting sentence cannot
	// return to a non-owner document. Comments are stripped and whitespace
	// collapsed so a wrapped or commented-out reintroduction still fires.
	forbid := func(rel, needle, diag string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if strings.Contains(collapseSpace(stripHTMLComments(readIfExists(path))), needle) {
			diags = append(diags, diag)
		}
	}
	diags = append(diags, checkStructuredPhaseContract(readIfExists(filepath.Join(root, ".bench", "BENCH.md")))...)

	requireCollapsed(".agents/commands/bench-implement-spec.md", "apply `craft-seams`' split-or-grant rule",
		".agents/commands/bench-implement-spec.md dropped the craft-seams split-or-grant pointer")
	requireCollapsed(".agents/commands/bench-implement-spec.md",
		"Every spec-backed run assigns genuine write work to at least one write subagent before the first implementation edit",
		".agents/commands/bench-implement-spec.md dropped the mandatory spec-backed write-delegation-before-first-edit contract")
	requireCollapsed(".agents/commands/bench-implement-spec.md",
		"independent vertical slices fan out to separate parallel subagents within the harness's concurrency limit; dependent slices run sequentially; a spec that lands as one atomic diff is delegated whole to one worktree-isolated write subagent",
		".agents/commands/bench-implement-spec.md dropped a delegation routing shape (independent-parallel, dependent-sequential, or atomic-whole)")
	requireCollapsed(".agents/commands/bench-implement-spec.md",
		"A read-only helper (research, review, planning, search) does not satisfy the write requirement.",
		".agents/commands/bench-implement-spec.md dropped the read-only exclusion from the write requirement")

	// One owner per workflow agreement: each repaired agreement pins the owner's
	// full statement with a require and the reintroduced contradiction with a
	// forbid, so the same fact cannot be stated twice and drift apart again.

	// Shaping — /bench-write-spec's entry contract owns the map requirement;
	// README points at the inline-map route instead of offering a skip.
	forbid("README.md", "Skip `/bench-shape-idea`",
		"README.md reintroduces the shaping skip route; every spec has a map behind it and /bench-write-spec's entry contract owns the inline-map recording path")
	requireCollapsed("README.md", "`/bench-write-spec`'s entry contract records the map inline",
		"README.md dropped the inline-map route pointer; /bench-write-spec's entry contract owns the shaping requirement")

	// Delegation — craft-delegate owns the capability-aware policy in full;
	// /bench-implement-spec points at it and states no inline threshold of its own.
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md",
		"admitted by the lighter-path threshold",
		".agents/skills/bench-craft-delegate/SKILL.md dropped the lighter-path inline allowance from the delegation policy")
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md",
		"stops before editing and emits one executable resume handoff to a subagent-capable harness — the repository path, the working branch or worktree, the spec or change name, the destination harness, and that harness's exact invocation",
		".agents/skills/bench-craft-delegate/SKILL.md dropped the no-write-subagent stop-and-handoff rule")
	forbid(".agents/commands/bench-implement-spec.md", "the sole inline exception",
		".agents/commands/bench-implement-spec.md restates an inline threshold of its own; craft-delegate owns the capability-aware delegation policy")
	requireCollapsed(".agents/commands/bench-implement-spec.md",
		"states no inline threshold of its own",
		".agents/commands/bench-implement-spec.md dropped the pointer to craft-delegate's capability-aware delegation policy")

	// Landing — /bench-final-check owns the landing commit and the status
	// transition; /bench-implement-spec ends at its last green build commit.
	forbid(".agents/commands/bench-implement-spec.md", `bench commit -m "<msg>" --spec`,
		".agents/commands/bench-implement-spec.md reclaims the landing --spec commit; /bench-final-check owns the landing commit and the Status: implemented transition")
	requireCollapsed(".agents/commands/bench-implement-spec.md",
		"ends at its last green build commit",
		".agents/commands/bench-implement-spec.md dropped the hand-off that ends at the last green build commit")
	requireCollapsed(".agents/commands/bench-final-check.md",
		"owns the landing commit and the spec's `Status: implemented` transition",
		".agents/commands/bench-final-check.md dropped the landing-commit and status-transition ownership")
	requireCollapsed(".agents/commands/bench-final-check.md",
		"nothing left to commit is reported green",
		".agents/commands/bench-final-check.md dropped the honest no-op for a branch with nothing to commit")
	requireCollapsed(".agents/commands/bench-final-check.md",
		"still performed via `bench spec implemented <slug>`",
		".agents/commands/bench-final-check.md dropped the bench spec implemented route for the status flip")

	// Red observation — /bench-debug commits the repro only in the project's
	// expected-failure form; a red-tree commit has no sanctioned path.
	forbid(".agents/commands/bench-debug.md", "commit that test before launching the shift",
		".agents/commands/bench-debug.md reintroduces the red repro commit before the shift; the repro is committed in the project's expected-failure form so the tree stays green")
	requireCollapsed(".agents/commands/bench-debug.md",
		"committed in the project's expected-failure form",
		".agents/commands/bench-debug.md dropped the expected-failure quarantine form for committing the repro")
	requireCollapsed(".agents/commands/bench-debug.md",
		"quarantine marker naming the bug",
		".agents/commands/bench-debug.md dropped the quarantine marker naming the bug")
	requireCollapsed(".agents/commands/bench-debug.md",
		"keeps the repro out of the shift and runs it by hand",
		".agents/commands/bench-debug.md dropped the no-expected-failure-form fallback")
	requireCollapsed(".agents/skills/bench-craft-seams/SKILL.md", "check both the file-length budget and the directory's file-count headroom",
		".agents/skills/bench-craft-seams/SKILL.md dropped the structure split-vs-grant headroom rule")
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md", "pins every file-tool path to that root",
		".agents/skills/bench-craft-delegate/SKILL.md dropped the shared-worktree file-tool path pin")
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md", "names a commit-specific sentinel",
		".agents/skills/bench-craft-delegate/SKILL.md dropped the fix-pass snapshot sentinel precondition")
	delegationPolicyAnchors := []struct {
		caseName string
		anchor   string
	}{
		{"insertion", "one source-line insertion"},
		{"deletion", "one source-line deletion"},
		{"replacement", "one source-line replacement"},
		{"replacement accounting", "A replacement counts as one correction"},
		{"request scope", "does not reset when work is split into tasks, slices, delegates, or verification rounds"},
		{"review repair", "Repairs beyond the allowance under Delegate or inline are re-charged to a write-delegate"},
		{"owning checkout", "coordinator verifies the repair in the checkout that owns the diff"},
	}
	for _, policyCase := range delegationPolicyAnchors {
		requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md", policyCase.anchor,
			fmt.Sprintf(".agents/skills/bench-craft-delegate/SKILL.md dropped the %s delegation-policy case", policyCase.caseName))
	}
	requireCollapsed(".agents/commands/bench-write-spec.md", "stale-command-reference sweep remains fail-closed across staged specs",
		".agents/commands/bench-write-spec.md dropped the staged-spec fail-closed command sweep posture")
	requireCollapsed(".agents/commands/bench-what-next.md", "use `bench spec history <slug>` for the shipped-row check",
		".agents/commands/bench-what-next.md dropped the bench spec history shipped-row check")
	whatNext := readIfExists(filepath.Join(root, ".agents", "commands", "bench-what-next.md"))
	if whatNext != "" && (strings.Count(whatNext, "bench roadmap --context") != 1 ||
		!strings.Contains(collapseSpace(whatNext), "If the query fails, stop the phase") ||
		!strings.Contains(collapseSpace(whatNext), "manual evidence reconstruction")) {
		diags = append(diags, "bench-what-next dropped the roadmap context query")
	}
	// Kit-only surfaces — the shipped guide and its generated skills index name only
	// what a consumer actually has: the guide says once that the maintenance surfaces
	// ship in the kit repository alone and points consumers at `bench upgrade`, and the
	// index marks the same rows the payload allowlist withholds.
	requireCollapsed(".bench/BENCH.md", "ship only in the Bench kit repository",
		".bench/BENCH.md does not state that the kit-maintenance surfaces ship only in the kit repository")
	requireCollapsed(".bench/BENCH.md", "a linked repo upgrades with `bench upgrade`",
		".bench/BENCH.md does not name bench upgrade as the consumer's route onto a newer kit")
	requireCollapsed(".bench/BENCH-reference.md", "`.agents/skills/bench-craft-synthesis/SKILL.md` (kit-only)",
		".bench/BENCH-reference.md skills index does not mark the kit-only craft-synthesis row")

	requireCollapsed(".bench/BENCH.md", "Parked ideas land in `IDEAS.md`",
		".bench/BENCH.md Capture section does not name IDEAS.md as the capture sink")
	requireCollapsed(".bench/BENCH.md", "append the dated line (`- YYYY-MM-DD <text>`) to `IDEAS.md`",
		".bench/BENCH.md Capture section lost the no-PATH fallback append to IDEAS.md")
	requireCollapsed(".agents/commands/bench-write-spec.md", "promote-then-delete commit removes the spec's `ROADMAP.md` row",
		".agents/commands/bench-write-spec.md does not remove the spec's ROADMAP.md row in the promote-then-delete commit (row presence is status)")
	requireCollapsed(".agents/commands/bench-shape-idea.md", "never pause for permission or a re-prompt",
		".agents/commands/bench-shape-idea.md dropped the resume-mode grill continuation rule; a running grill carries into newly-unblocked tickets without pausing for a re-prompt")
	requireCollapsed(".agents/commands/bench-review-implementation.md", "Integrate the findings accepted for this round, run focused checks for the changed behavior, then run one final gate and stop. Open another semantic review round only when that gate fails or the reviewer requests one.",
		".agents/commands/bench-review-implementation.md dropped the terminal repair-pass bound")

	shapeIdeaPath := filepath.Join(root, ".agents", "commands", "bench-shape-idea.md")
	shapeIdeaActive := collapseSpace(stripHTMLComments(readIfExists(shapeIdeaPath)))
	if exists(shapeIdeaPath) && strings.Contains(shapeIdeaActive, "straight to `/bench-write-spec`") {
		diags = append(diags, ".agents/commands/bench-shape-idea.md reintroduces the skip-to-spec bypass fragment; every idea must yield a map with a Handoff before a spec")
	}
	writeSpecPath := filepath.Join(root, ".agents", "commands", "bench-write-spec.md")
	writeSpecActive := collapseSpace(stripHTMLComments(readIfExists(writeSpecPath)))
	if exists(writeSpecPath) {
		if !strings.Contains(writeSpecActive, "refuses to run without") {
			diags = append(diags, ".agents/commands/bench-write-spec.md dropped the map-required entry contract (refuses to run without a complete map)")
		}
		reviewerClosedFastPathAnchors := []string{
			"Default spec authoring starts in a fresh mid-tier session",
			"sole same-session exception",
			"every load-bearing fork has already been put to the reviewer and closed in the current session",
			"write those decisions directly into a new decision map with a complete Handoff",
			"continue from that file rather than unwritten grill memory",
			"compile the spec without routing through `/bench-shape-idea`",
		}
		for _, anchor := range reviewerClosedFastPathAnchors {
			if !strings.Contains(writeSpecActive, anchor) {
				diags = append(diags, ".agents/commands/bench-write-spec.md dropped the active reviewer-closed-forks same-session fast path")
				break
			}
		}
		if !strings.Contains(writeSpecActive, "if any fork remains open, run `/bench-shape-idea` and keep the normal map gate") {
			diags = append(diags, ".agents/commands/bench-write-spec.md dropped the explicit open-fork fallback to /bench-shape-idea")
		}
	}
	if exists(shapeIdeaPath) && !strings.Contains(shapeIdeaActive, "`/bench-write-spec`'s entry contract owns the narrow recording path for decisions already closed with the reviewer in the current session") {
		diags = append(diags, ".agents/commands/bench-shape-idea.md dropped the pointer to /bench-write-spec's entry contract for reviewer-closed forks")
	}

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
	section := markdownH2Section(stripHTMLComments(sharedRules), "How to talk to me")
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
	declared := structuredPhaseClauseNames(block)
	if len(declared) == 0 {
		return []string{".bench/BENCH.md structured Bench phase contract declares no named clauses"}
	}
	bodies, bodyDiags := structuredPhaseClauseBodies(block)
	diags := append([]string(nil), bodyDiags...)
	seen := map[string]bool{}
	for _, name := range declared {
		if seen[name] {
			diags = append(diags, fmt.Sprintf(".bench/BENCH.md structured Bench phase contract declares clause %q more than once", name))
			continue
		}
		seen[name] = true
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

func markdownH2Section(text, title string) string {
	lines := strings.Split(text, "\n")
	heading := "## " + title
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == heading {
			start = i + 1
			continue
		}
		if start >= 0 && strings.HasPrefix(line, "## ") {
			return strings.Join(lines[start:i], "\n")
		}
	}
	if start < 0 {
		return ""
	}
	return strings.Join(lines[start:], "\n")
}

func stripHTMLComments(text string) string {
	for {
		start := strings.Index(text, "<!--")
		if start < 0 {
			return text
		}
		end := strings.Index(text[start+4:], "-->")
		if end < 0 {
			return text[:start]
		}
		text = text[:start] + text[start+4+end+3:]
	}
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

func collapseSpace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
