package conformance

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	structuredPhaseDeclaration = "- **Structured Bench phase conversation:**"
	structuredPhaseUnavailable = ".bench/BENCH.md cannot verify the structured Bench phase contract because shared rules are missing or empty"
	// Shared by checkSharedRuleSingleSource's marker list and the section-scoped
	// placement anchors, so one string pins each rule's presence, non-duplication,
	// and placement.
	fixDontParkMarker   = "lands in the active workflow as its own commit"
	sourceWarrantMarker = "names what you read and what you did not"
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
	require(".agents/skills/bench-craft-spec/SKILL.md", "Slicing a build for delegates")
	require(".agents/skills/bench-craft-delegate/SKILL.md", "Slicing a build for delegates")
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
	require(".agents/commands/bench-final-check.md", "ship-tier verification has not run")
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
	require(".agents/skills/bench-craft-delegate/SKILL.md", "releases the worktree it cut")
	require(".agents/commands/bench-implement-spec.md", "When the build stops short")
	require(".agents/commands/bench-implement-spec.md", "the coordinator that cut it owns its retirement")
	require(".bench/BENCH.md", "bench worktree release")
	require(".bench/BENCH.md", "bench worktree clean")
	require(".bench/BENCH.md", "bench worktree recovery")
	require(".agents/commands/bench-write-spec.md", "Superseded by")
	require(".agents/commands/bench-write-spec.md", "spec-retire:")
	require(".agents/commands/bench-write-spec.md", "Status: staged")
	require(".agents/commands/bench-write-spec.md", "new session on the mid tier")
	require(".agents/commands/bench-write-spec.md", "mostly not")
	require(".agents/commands/bench-write-spec.md", "runs at the mid tier")
	require(".agents/commands/bench-write-spec.md", "Every draft gets the pass")
	require("projects/benchkit.md", "Spec falsification pass")
	require("projects/benchkit.md", "shared-build-cache opt-in")
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
	// scopedSection resolves the one H2 section a family of scoped anchors owns,
	// diagnosing an absent file, a missing section, and a duplicated heading each
	// exactly once, ahead of the per-anchor matching. It returns the section body
	// collapsed and lowercased: whole-file requireCollapsed neither strips
	// comments nor sees headings, so a sentence commented out or pasted into the
	// wrong section would satisfy it, and case is folded so a sentence-initial
	// recasing cannot slip past a needle.
	scopedSection := func(rel, section string) (string, bool) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if !exists(path) {
			diags = append(diags, "section-scoped anchor file missing: "+rel)
			return "", false
		}
		body, count := markdownH2Sections(stripHTMLComments(readIfExists(path)), section)
		if count == 0 {
			diags = append(diags, fmt.Sprintf("%s is missing the %q section that owns a scoped anchor", rel, section))
			return "", false
		}
		if count > 1 {
			diags = append(diags, fmt.Sprintf("%s carries %d %q sections; a scoped anchor needs exactly one owning section", rel, count, section))
			return "", false
		}
		return strings.ToLower(collapseSpace(body)), true
	}
	// requireInSection is the placement half of a section-owned anchor and
	// forbidInSection the must-not half, banning the contradiction from the
	// section that owns the fact so a contradicting sentence cannot sit beside
	// the fact it undoes. Both take a body scopedSection already resolved, so a
	// broken section reports once rather than once per anchor.
	requireInSection := func(body, needle, diag string) {
		if !strings.Contains(body, strings.ToLower(needle)) {
			diags = append(diags, diag)
		}
	}
	forbidInSection := func(body, needle, diag string) {
		if strings.Contains(body, strings.ToLower(needle)) {
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

	// Shaping is situational; README names both authorized routes without
	// copying the command's full entry contract.
	forbid("README.md", "Every spec has a decision map behind it",
		"README.md reintroduced mandatory decision maps; shaping is situational")
	requireCollapsed("README.md", "Decision maps are situational",
		"README.md dropped the situational decision-map vocabulary")

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
		{"review repair", "Repairs beyond the allowance under Delegate or inline are routed as Verifying the done-claim directs"},
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

	// Shared-rule placement — checkSharedRuleSingleSource's marker list owns each
	// sentence's presence and non-duplication; these anchors add only placement,
	// pinning the rule inside the section that owns it.
	if body, ok := scopedSection(".bench/BENCH.md", "Workflow"); ok {
		requireInSection(body, fixDontParkMarker,
			".bench/BENCH.md Workflow section dropped the fix-don't-park rule; a mid-work defect fix belongs in the active workflow, not the backlog")
	}
	if body, ok := scopedSection(".bench/BENCH.md", "How to talk to me"); ok {
		requireInSection(body, sourceWarrantMarker,
			".bench/BENCH.md How to talk to me section dropped the outside-source warrant rule; a claim resting on a source outside the tree names what was and was not read")
	}
	forbid(".bench/BENCH.md", "thorough",
		".bench/BENCH.md phrases the outside-source warrant rule as thoroughness; the rule asks for disclosure of what went unread, which the reviewer can check — thoroughness nobody can")
	requireCollapsed(".agents/commands/bench-write-spec.md", "promote-then-delete commit removes the spec's `ROADMAP.md` row",
		".agents/commands/bench-write-spec.md does not remove the spec's ROADMAP.md row in the promote-then-delete commit (row presence is status)")
	requireCollapsed(".agents/commands/bench-shape-idea.md", "never pause for permission or a re-prompt",
		".agents/commands/bench-shape-idea.md dropped the resume-mode grill continuation rule; a running grill carries into newly-unblocked tickets without pausing for a re-prompt")
	requireCollapsed(".agents/commands/bench-review-implementation.md", "This phase makes no fixes and runs no gate: findings that need work go to `/bench-implement-spec`, which owns the fix pass, the pickup file's resolution, and the terminal repair-pass bound; a clean or risk-accepted review goes to `/bench-final-check`, which owns the oracle run.",
		".agents/commands/bench-review-implementation.md dropped the hand-off-don't-repair rule; the fix pass belongs to implement-spec and the oracle run to final-check")

	// Ticket guidance is convention-only, so these anchors pin the load-bearing
	// workflow clauses without inventing a parser for ticket files.
	requireCollapsed(".agents/commands/bench-implement-spec.md", "Charge `craft-tickets` before the first implementation edit",
		".agents/commands/bench-implement-spec.md dropped the craft-tickets breakdown charge")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "derive ticket files from the spec's stories and seams",
		".agents/commands/bench-implement-spec.md dropped the ticket breakdown derivation from the spec's stories and seams")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "`specs/<slug>/tickets/`",
		".agents/commands/bench-implement-spec.md dropped the ticket breakdown destination")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "under the session's existing approval surface",
		".agents/commands/bench-implement-spec.md dropped the ticket breakdown approval surface")

	requireCollapsed(".bench/BENCH.md", "one independently-green ticket",
		".bench/BENCH.md dropped the light-path independently-green-ticket observable")
	requireCollapsed(".bench/BENCH.md", "crosses no declared seam",
		".bench/BENCH.md dropped the light-path declared-seam observable")

	requireCollapsed(".agents/skills/bench-craft-line/SKILL.md", "| Orchestration | mid + medium |",
		".agents/skills/bench-craft-line/SKILL.md dropped the orchestration mid/medium stage default")
	requireCollapsed(".agents/skills/bench-craft-line/SKILL.md", "| Ticket implementation | cheap + low |",
		".agents/skills/bench-craft-line/SKILL.md dropped the ticket implementation cheap/low stage default")
	requireCollapsed(".agents/skills/bench-craft-line/SKILL.md", "| Review (axis or falsification) | mid + high |",
		".agents/skills/bench-craft-line/SKILL.md dropped the review mid/high stage default")

	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "smallest independently-green",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the smallest-independently-green ticket contract")
	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "one write-delegate charge",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the one-write-delegate-charge ticket contract")
	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "Blocked by:",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the ticket template Blocked by heading")
	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "## What to build",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the ticket template What to build heading")
	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "## Acceptance",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the ticket template Acceptance heading")
	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "- [ ] <Observable behavioral criterion>",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the ticket template behavioral acceptance checkbox")

	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "`craft-spec` owns the spec-time **who-writes-where** fence",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the craft-spec ownership-fence cross-pointer")
	requireCollapsed(".agents/skills/bench-craft-spec/SKILL.md", "`craft-tickets` owns the build-time **what-lands-green-next** unit",
		".agents/skills/bench-craft-spec/SKILL.md dropped the craft-tickets build-time-unit cross-pointer")

	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "do not run a standalone full gate",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the no-standalone-full-gate ticket cadence")
	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "`bench commit` is the only per-ticket full-project-gate boundary",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the bench-commit-only ticket gate boundary")
	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "If it goes red, repair from that output and retry",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the ticket red repair-and-retry cadence")
	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "normal green path runs one full gate",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the ticket one-full-gate green cadence")
	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "/bench-final-check` remains the final full gate over the composed feature",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the composed-feature final gate")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "no standalone full gate before landing",
		".agents/commands/bench-implement-spec.md dropped the no-standalone-full-gate ticket cadence")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "`bench commit` is the only per-ticket full-project-gate boundary",
		".agents/commands/bench-implement-spec.md dropped the bench-commit-only ticket gate boundary")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "If it goes red, repair from that output and retry",
		".agents/commands/bench-implement-spec.md dropped the ticket red repair-and-retry cadence")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "normal green path is one full gate",
		".agents/commands/bench-implement-spec.md dropped the ticket one-full-gate green cadence")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "/bench-final-check` still runs the final full gate over the composed feature",
		".agents/commands/bench-implement-spec.md dropped the composed-feature final gate")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "A repair pass integrates the findings accepted for its round and stops at its repair ticket's green landing; another semantic review round opens only when that gate fails or the reviewer requests one.",
		".agents/commands/bench-implement-spec.md dropped the terminal repair-pass bound")

	requireCollapsed(".agents/commands/bench-write-spec.md", "Top-level `decisions/` holds pre-spec working maps",
		".agents/commands/bench-write-spec.md dropped the top-level pre-spec working-map posture")
	requireCollapsed(".agents/commands/bench-shape-idea.md", "Run `bench maps` before declaring readiness",
		".agents/commands/bench-shape-idea.md dropped the bench-maps readiness check")
	requireCollapsed(".agents/commands/bench-write-spec.md", "move (do not copy) the source map and any map-owned assets from top-level `decisions/` into `specs/<slug>/decisions/`",
		".agents/commands/bench-write-spec.md dropped the compile-time decision-map move")
	requireCollapsed(".agents/commands/bench-write-spec.md", "update every reference to the moved paths",
		".agents/commands/bench-write-spec.md dropped the moved decision-map reference update")
	requireCollapsed(".agents/commands/bench-write-spec.md", "re-run reads the already-compiled spec-local map; it never recreates a top-level copy",
		".agents/commands/bench-write-spec.md dropped the compiled decision-map re-run lifecycle")
	requireCollapsed(".agents/commands/bench-write-spec.md", "Whole-folder retirement removes the compiled maps and map-owned assets",
		".agents/commands/bench-write-spec.md dropped the whole-folder compiled decision-map retirement")

	requireCollapsed(".agents/commands/bench-final-check.md", "only after both the spec's landing gate and landing commit are green",
		".agents/commands/bench-final-check.md dropped the after-green implementation-retro placement")
	for _, section := range []string{
		"## Outcome",
		"## Gate-stage timings",
		"## Ticket-versus-spec-slice and delegate performance",
		"## Coordinator catches",
		"## Agent-experience improvements",
		"### Bench CLI",
		"### Skills",
		"### Process",
	} {
		requireCollapsed(".agents/commands/bench-final-check.md", section,
			fmt.Sprintf(".agents/commands/bench-final-check.md dropped the required implementation-retro evidence section: %s", section))
	}
	requireCollapsed(".agents/commands/bench-final-check.md", "rewrite `.bench/retros/<spec-slug>.md` in full",
		".agents/commands/bench-final-check.md dropped whole-file implementation-retro replacement")
	requireCollapsed(".agents/commands/bench-final-check.md", "Do not run another gate or commit just to capture the retro",
		".agents/commands/bench-final-check.md added an implementation-retro gate outside the normal cadence")
	requireCollapsed(".bench/BENCH.md", "/bench-final-check` writes `.bench/retros/<spec-slug>.md`",
		".bench/BENCH.md dropped the implementation-retro capture owner")

	requireCollapsed(".agents/commands/bench-what-next.md", "The snapshot's `retros` bodies are the only retro evidence this run reads",
		".agents/commands/bench-what-next.md dropped the roadmap-context retro drain source")
	requireCollapsed(".agents/commands/bench-what-next.md", "merge into an existing roadmap row, a new roadmap row, a learning-or-rule disposition, or an explicit dismissal",
		".agents/commands/bench-what-next.md dropped an implementation-retro recommendation disposition")
	requireCollapsed(".agents/commands/bench-what-next.md", "remove every drained `.bench/retros/*.md` file in the same reviewer-approved batch",
		".agents/commands/bench-what-next.md dropped the delete-all implementation-retro drain rule")
	requireCollapsed(".bench/BENCH.md", "/bench-what-next` owns their reviewed drain",
		".bench/BENCH.md dropped the implementation-retro drain owner")

	// Verify hooks — the Roles verify rule's point-of-use pointers, one per phase
	// command. Whole-file anchors by design: the three files hang their hooks on
	// three different structures, so the gate proves presence and the review axis
	// grades placement.
	requireCollapsed(".agents/commands/bench-shape-idea.md",
		"Before asking me about a fact, look it up in the tree",
		".agents/commands/bench-shape-idea.md dropped the look-it-up-before-asking hook; the shaping grill spends reviewer answers on decisions, and looks facts up in the tree")
	requireCollapsed(".agents/commands/bench-implement-spec.md",
		"Verify a claim against the tree as it stands, not against memory; a claim over a whole set is verified by enumerating the set, never by extending one measured member.",
		".agents/commands/bench-implement-spec.md dropped the verify-against-the-tree hook; the build verifies claims against the tree, not memory, and a set claim by enumerating the set")
	requireCollapsed(".agents/commands/bench-review-implementation.md",
		"A finding cites what its axis read now, not what it recalls; the bar for a universal claim — cite the enumeration or name itself a sample — is that citation standard's.",
		".agents/commands/bench-review-implementation.md dropped the cite-what-you-read hook; a review finding cites what its axis read now, and a universal claim answers to craft-review's citation standard")

	// --full — the orchestrated-run contract lives in one bounded section of
	// bench-implement-spec.md and nowhere else. Every anchor is scoped to that
	// section so it means placement, not mere presence, and each fact with a
	// cheap contradicting reading is a require/forbid pair.
	const fullRunFile = ".agents/commands/bench-implement-spec.md"
	const fullRunSection = "The `--full` run"
	fullRunRequires := []struct{ needle, diag string }{
		{"`--full` is an opt-in flag: plain `/bench-implement-spec` keeps today's implement-only semantics",
			"dropped the opt-in entry contract; plain invocation keeps implement-only semantics"},
		{"`/bench-review-implementation`, `/bench-final-check`, and `/bench-debug` stay standalone commands for strict phased use and mid-run resumption",
			"dropped the standalone-phases guarantee for strict phased use and mid-run resumption"},
		{"`--full` invoked with no spec argument, or with one naming a path that does not exist, refuses at the entry contract and says which of the two it was",
			"dropped the missing-or-unknown-spec refusal; a full run never builds a guessed target"},
		{"spawn one fresh-context delegate charged with the standalone `/bench-review-implementation` contract and given the spec and the diff and nothing else",
			"dropped the fresh-context review delegate and its spec-and-diff-only inputs"},
		{"Inline self-review is closed, not deprioritized — the context that produced the code carries the assumptions that produced its bugs — and a context-inheriting delegate is the same failure wearing a delegate's name",
			"dropped the closure of inline self-review and context inheritance"},
		{"The delegate's done-claim answers to invariant 1 in `.bench/BENCH.md` and `craft-delegate`'s verification discipline",
			"dropped the done-claim verification pointer to invariant 1 and craft-delegate"},
		{"When the review delegate cannot run or returns nothing, the run stops and reports at that boundary rather than proceeding to final-check with review unrun",
			"dropped the stop-at-the-boundary rule for an unavailable or empty-handed review delegate"},
		{"Concrete defects — bugs, spec misses, missing coverage — are fixed and re-gated without stopping; contestable design and judgment findings are flagged in the exit report for reviewer veto, not applied",
			"dropped a finding-disposition half (fix-and-re-gate concrete, or flag-don't-apply judgment)"},
		{"bounded by this command's terminal repair-pass bound and routed through `craft-delegate`'s repair allowance",
			"dropped the repair-bound pointers; this mode adds no second version of either rule"},
		{"At every phase boundary the run writes the phase reached into `session-handoff.md`'s State section, then refreshes the pin block with `bench handoff --next <command>`",
			"dropped the phase-boundary State write and pin-block refresh"},
		{"A re-invoked `--full` resumes from the phase the handoff names instead of re-implementing from the top; a stale handoff is arbitrated by `AGENTS.md`'s tree-wins rule",
			"dropped the resume-don't-restart rule for a re-invoked run"},
		{"pause and ask the reviewer as a structured decision list with a recommendation for this run, offering three routes: continue at the mid binding; escalate to the top binding in this harness; escalate to the top binding via the Codex CLI",
			"dropped the pause-and-ask escalation menu with its three fixed routes"},
		{"a route is omitted only when this harness cannot invoke it at all, and the omission is stated rather than silent",
			"dropped the stated-omission rule for an uninvokable escalation route"},
		{"Never escalate without asking. A harness with no structured-prompt surface asks the same question as a plain numbered list",
			"dropped the never-escalate-without-asking rule and its numbered-list fallback"},
		{"The run implements the spec's stories and nothing else: work noticed outside them — an adjacent refactor, an unrelated improvement, a story the spec chose not to take — is recorded for the reviewer rather than built",
			"dropped the scope fence; out-of-story work is recorded, never built"},
		{"one disposition per row of the spec's acceptance coverage map — implemented, deferred, or won't-handle — named row by row against `bench coverage <spec>`'s enumeration",
			"dropped the per-row coverage accounting against bench coverage's enumeration"},
		{"When the spec carries no coverage map, the report says so and accounts for the user stories instead",
			"dropped the no-coverage-map fallback to user-story accounting"},
		{"Every phase claim in the exit report cites the record that proves it: the review delegate's invocation, the commit shas the phase landed, the `session-handoff.md` boundary rewrite",
			"dropped the record-citation rule for phase claims in the exit report"},
		{"ask a second, separate question: whether to add a cross-harness falsification pass over the diff before final-check — the Codex CLI at the top binding, charged to refute the claim that the spec was implemented rather than to grade it against the three axes",
			"dropped the falsification-pass offer with its refutation charge"},
		{"its own question, not a fourth route in the escalation menu, and it never runs standing: absent the trigger it is not offered",
			"dropped the falsification pass's separate-question and never-standing bounds"},
	}
	fullRunBody, fullRunOK := scopedSection(fullRunFile, fullRunSection)
	if fullRunOK {
		for _, a := range fullRunRequires {
			requireInSection(fullRunBody, a.needle,
				fullRunFile+" `--full` section "+a.diag)
		}
	}
	fullRunForbids := []struct{ needle, diag string }{
		{"by default", "phrases `--full` as default-on; the mode is opt-in and plain invocation keeps implement-only semantics"},
		{"infers the spec", "reintroduces a spec-inferring fallback; a missing or unknown spec argument refuses and says which"},
		{"review the diff inline", "reintroduces the inline self-review fallback; the review runs in one fresh-context delegate"},
		{"stops on a concrete defect", "stops the run on a concrete defect; concrete defects are fixed and re-gated in-run"},
		{"applies a judgment finding", "applies a judgment finding; design and judgment findings are flagged for reviewer veto"},
		{"obviously needed", "reintroduces escalate-if-obviously-needed; escalation always pauses and asks the reviewer"},
		{"while the file is open", "licenses opportunistic improvement while the file is open; the scope fence records it instead"},
		{"fully implemented", "reports completion in aggregate; the exit report gives one disposition per coverage-map row"},
		{"without its record", "reports a phase complete without its record; every phase claim cites the record that proves it"},
		{"on every run", "makes the falsification pass standing; it is offered only on the size trigger"},
	}
	if fullRunOK {
		for _, a := range fullRunForbids {
			forbidInSection(fullRunBody, a.needle,
				fullRunFile+" `--full` section "+a.diag)
		}
	}

	diags = append(diags, checkSpecAuthorizationContract(root)...)
	for _, anchor := range []struct{ rel, needle, diag string }{
		{".agents/commands/bench-write-spec.md", "Record exactly one `Decision source:` line", ".agents/commands/bench-write-spec.md dropped the exactly-one Decision source contract"},
		{".agents/commands/bench-write-spec.md", "re-read and re-verify every structured `## Sources` entry", ".agents/commands/bench-write-spec.md dropped the map Sources re-verification or single-manifest contract"},
		{".agents/commands/bench-write-spec.md", "without copying a research manifest into the spec", ".agents/commands/bench-write-spec.md dropped the map Sources re-verification or single-manifest contract"},
		{".agents/commands/bench-write-spec.md", "Ask at most two late clarification questions, one at a time, each with a recommended answer; route a dependency tree or multi-session fog to `$bench-shape-idea`", ".agents/commands/bench-write-spec.md dropped the bounded late-uncertainty route"},
		{".agents/commands/bench-shape-idea.md", "Decision maps are situational, not mandatory", ".agents/commands/bench-shape-idea.md dropped the situational-map contract"},
		{".agents/commands/bench-shape-idea.md", "Run `bench maps --template` for the canonical paste-ready schema", ".agents/commands/bench-shape-idea.md dropped the canonical map-template pointer"},
		{".agents/commands/bench-shape-idea.md", "A decision ticket is one unresolved reviewer choice", ".agents/commands/bench-shape-idea.md dropped decision-ticket vocabulary"},
		{".agents/commands/bench-shape-idea.md", "Shaping owns reviewer decisions, constraints, exclusions, research objects, rejected alternatives, and bounded discretion", ".agents/commands/bench-shape-idea.md dropped shaping ownership of reviewer decisions, constraints, exclusions, research objects, rejected alternatives, or bounded discretion"},
		{".agents/commands/bench-write-spec.md", "Spec authoring owns engineering seams, deep-versus-thin design, tests, acceptance coverage, hostile-input attachment, and gate attachment", ".agents/commands/bench-write-spec.md dropped spec ownership of engineering seams, tests, coverage, hostile inputs, or gate attachment"},
	} {
		requireCollapsed(anchor.rel, anchor.needle, anchor.diag)
	}
	forbid(".agents/commands/bench-shape-idea.md", "## Handoff",
		".agents/commands/bench-shape-idea.md reintroduced decision-map Handoff ownership")
	forbid(".agents/commands/bench-write-spec.md", "## Handoff",
		".agents/commands/bench-write-spec.md reintroduced decision-map Handoff consumption")
	forbid(".agents/commands/bench-write-spec.md", "map's Handoff",
		".agents/commands/bench-write-spec.md reintroduced decision-map Handoff consumption")
	requireCollapsed(".agents/skills/bench-craft-grill/SKILL.md", "Grill is a decision-ticket type",
		".agents/skills/bench-craft-grill/SKILL.md dropped decision-ticket vocabulary")
	requireCollapsed(".agents/skills/bench-craft-spec/SKILL.md", "authorized decision source",
		".agents/skills/bench-craft-spec/SKILL.md reintroduced Handoff-owned acceptance behavior")
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md", "named decision source",
		".agents/skills/bench-craft-delegate/SKILL.md reintroduced Handoff-owned delegate context")
	requireCollapsed("CONTEXT.md", "**decision ticket**",
		"CONTEXT.md dropped the decision-ticket versus implementation-ticket distinction")
	requireCollapsed("README.md", "independently-green implementation tickets",
		"README.md dropped implementation-ticket vocabulary")
	requireCollapsed("projects/benchkit.md", "derives engineering seams, tests, coverage, hostile-input attachment, and gate attachment",
		"projects/benchkit.md dropped spec-authoring phase ownership")
	requireCollapsed("CHANGELOG.md", "separated shaping decision tickets from independently-green implementation tickets",
		"CHANGELOG.md dropped the decision-ticket phase-ownership change")

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

// markdownH2Section returns the first "## title" section body, or "" when the
// heading is absent. A caller that must reject a duplicated heading — one whose
// second body could contradict the first unseen — uses markdownH2Sections for
// the occurrence count.
func markdownH2Section(text, title string) string {
	body, _ := markdownH2Sections(text, title)
	return body
}

// markdownH2Sections returns the first "## title" section body and how many
// times the heading occurs in text.
func markdownH2Sections(text, title string) (string, int) {
	lines := strings.Split(text, "\n")
	heading := "## " + title
	count := 0
	start := -1
	end := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == heading {
			count++
			if count == 1 {
				start = i + 1
			}
			continue
		}
		if start >= 0 && end < 0 && strings.HasPrefix(line, "## ") {
			end = i
		}
	}
	if count == 0 {
		return "", 0
	}
	if end < 0 {
		end = len(lines)
	}
	return strings.Join(lines[start:end], "\n"), count
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
