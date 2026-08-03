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
	requireAll := func(rel string, needles ...string) {
		for _, needle := range needles {
			require(rel, needle)
		}
	}
	requireAll(".agents/commands/bench-write-spec.md", "acceptance coverage map", "why it catches the failure", "red signal", "bench-craft-spec", "seam diagram", "tests attach here", "edge inventory", "Won't handle", "hostile-input checklist", "Superseded by", "spec-retire:", "Status: staged", "new session on the mid tier", "mostly not", "runs at the mid tier", "Every draft gets the pass")
	requireAll(".agents/skills/bench-craft-spec/SKILL.md", "why it catches the failure", "re-run idempotency", "separate capability", "Slicing a build for delegates")
	requireAll(".agents/skills/bench-craft-tdd/SKILL.md", "bench-craft-spec", "acceptance row", "not TDD-able", "call count", "row schema and the red-signal definition are", "floor, not the ceiling")
	requireAll(".agents/skills/bench-craft-review/SKILL.md", "bench-craft-spec", "an edge nobody decided")
	requireAll(".agents/skills/bench-craft-delegate/SKILL.md", "Slicing a build for delegates", "a claim, not a result", "bench worktree create --request", "git stash", "releases the worktree it cut")
	requireAll(".agents/skills/bench-craft-seams/SKILL.md", "failure modes", "structure.budgets")
	requireAll(".agents/commands/bench-implement-spec.md", "coverage table", "already covered", "turning red-to-green", "bench coverage <spec>", "When the build stops short", "the coordinator that cut it owns its retirement", "Status: implemented", "reviews/<spec-slug>.md", "names and deletes the file", "bench commit -m")
	requireAll(".agents/commands/bench-review-implementation.md", "acceptance coverage map", "mapped behavior", "bench diff --full", "bench diff --full --commit", "## Coverage", "Coverage axis", "craft-review", "craft-delegate", "reviews/<spec-slug>.md", "same session that writes it", "actionable findings", "writes no artifact", "section per axis", "finding count, its worst issue", "doc citation its axis supplied")
	requireAll(".agents/commands/bench-final-check.md", ".bench/gate.sh", "BENCH_GATE", "ship-tier verification has not run", "craft-gate", "bench commit -m", "retained exact green evidence")
	requireAll(".agents/commands/bench-setup-repo.md", "hostile-input checklist", "craft-gate")
	requireAll(".agents/commands/bench-debug.md", "diff-filter=D", "through the accused command")
	requireAll(".agents/commands/bench-what-next.md", "Reconcile first", "through the accused command", "empties to zero", "verdict in the batch diff", "one uncommitted batch diff", "commit on green", "## Recommended sequence")
	requireAll(".agents/commands/bench-assess.md", "verify the previous assessment's backlog landed", "read-only area sweeps on the mid tier", "synthesize adversarially on the top tier", "replaces its predecessor", "/bench-what-next")
	requireAll(".bench/BENCH.md", "bench worktree release", "bench worktree clean", "bench worktree recovery", "sole gate", "terminal final-check never repays", "reauthors promotion")
	requireAll("projects/benchkit.md", "hostile-input checklist", "Spec falsification pass", "shared-build-cache opt-in")

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
	// forbid strips comments and spacing so a wrapped or commented contradiction still fires.
	forbid := func(rel, needle, diag string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if strings.Contains(collapseSpace(stripHTMLComments(readIfExists(path))), needle) {
			diags = append(diags, diag)
		}
	}
	// scopedSection resolves one H2 section, diagnosing absent, missing, or duplicate owners
	// exactly once, ahead of the per-anchor matching. It returns the section body
	// collapsed and lowercased: whole-file requireCollapsed neither strips
	// comments nor sees headings, so a sentence commented out or pasted into the
	// wrong section would satisfy it; case folding also catches sentence-initial recasing.
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
	requireCollapsed(".agents/commands/bench-implement-spec.md", "Every spec-backed run assigns genuine write work to at least one write subagent before the first implementation edit", ".agents/commands/bench-implement-spec.md dropped the mandatory spec-backed write-delegation-before-first-edit contract")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "independent vertical slices fan out to separate parallel subagents within the harness's concurrency limit; dependent slices run sequentially; a spec that lands as one atomic diff is delegated whole to one worktree-isolated write subagent", ".agents/commands/bench-implement-spec.md dropped a delegation routing shape (independent-parallel, dependent-sequential, or atomic-whole)")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "A read-only helper (research, review, planning, search) does not satisfy the write requirement.", ".agents/commands/bench-implement-spec.md dropped the read-only exclusion from the write requirement")

	// One owner per workflow agreement: each repaired agreement pins the owner's
	// full statement with a require and the reintroduced contradiction with a
	// forbid, so the same fact cannot be stated twice and drift apart again.

	// Shaping is situational; README names both authorized routes without
	// copying the command's full entry contract.
	forbid("README.md", "Every spec has a decision map behind it", "README.md reintroduced mandatory decision maps; shaping is situational")
	requireCollapsed("README.md", "Decision maps are situational", "README.md dropped the situational decision-map vocabulary")

	// Delegation — craft-delegate owns the capability-aware policy in full;
	// /bench-implement-spec points at it and states no inline threshold of its own.
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md", "admitted by the lighter-path threshold", ".agents/skills/bench-craft-delegate/SKILL.md dropped the lighter-path inline allowance from the delegation policy")
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md", "stops before editing and emits one executable resume handoff to a subagent-capable harness — the repository path, the working branch or worktree, the spec or change name, the destination harness, and that harness's exact invocation", ".agents/skills/bench-craft-delegate/SKILL.md dropped the no-write-subagent stop-and-handoff rule")
	forbid(".agents/commands/bench-implement-spec.md", "the sole inline exception", ".agents/commands/bench-implement-spec.md restates an inline threshold of its own; craft-delegate owns the capability-aware delegation policy")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "states no inline threshold of its own", ".agents/commands/bench-implement-spec.md dropped the pointer to craft-delegate's capability-aware delegation policy")

	// Landing — promotion alone owns reviewed spec-backed landing; final-check
	// only reports its retained terminal evidence, while ordinary work keeps the
	// established gate-then-commit route.
	forbid(".agents/commands/bench-implement-spec.md", `bench commit -m "<msg>" --spec`, ".agents/commands/bench-implement-spec.md reclaims the spec-backed landing from spec build promote")
	requireCollapsed(".agents/commands/bench-final-check.md", "`bench spec build promote` is the sole spec-backed gate, commit, and `Status: implemented` author.", ".agents/commands/bench-final-check.md gives spec-backed landing authority to a second command")
	requireCollapsed(".agents/commands/bench-final-check.md", "A terminal promoted run gets no second gate or landing mutation", ".agents/commands/bench-final-check.md dropped the no-second-gate terminal report rule")
	requireCollapsed(".agents/commands/bench-final-check.md", "Light-path and ordinary non-lifecycle work retain the gate-then-commit path", ".agents/commands/bench-final-check.md dropped ordinary final-check landing behavior")
	forbid(".agents/commands/bench-final-check.md", "Use `bench commit -m \"<msg>\" --spec", ".agents/commands/bench-final-check.md reintroduced bench commit --spec as a spec-build landing author")
	forbid(".agents/commands/bench-final-check.md", "still performed via `bench spec implemented <slug>`", ".agents/commands/bench-final-check.md reintroduced bench spec implemented as a second status author")

	// Red observation — /bench-debug commits the repro only in the project's
	// expected-failure form; a red-tree commit has no sanctioned path.
	forbid(".agents/commands/bench-debug.md", "commit that test before launching the shift", ".agents/commands/bench-debug.md reintroduces the red repro commit before the shift; the repro is committed in the project's expected-failure form so the tree stays green")
	requireCollapsed(".agents/commands/bench-debug.md", "committed in the project's expected-failure form", ".agents/commands/bench-debug.md dropped the expected-failure quarantine form for committing the repro")
	requireCollapsed(".agents/commands/bench-debug.md", "quarantine marker naming the bug", ".agents/commands/bench-debug.md dropped the quarantine marker naming the bug")
	requireCollapsed(".agents/commands/bench-debug.md", "keeps the repro out of the shift and runs it by hand", ".agents/commands/bench-debug.md dropped the no-expected-failure-form fallback")
	requireCollapsed(".agents/skills/bench-craft-seams/SKILL.md", "check both the file-length budget and the directory's file-count headroom", ".agents/skills/bench-craft-seams/SKILL.md dropped the structure split-vs-grant headroom rule")
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md", "pins every file-tool path to that root", ".agents/skills/bench-craft-delegate/SKILL.md dropped the shared-worktree file-tool path pin")
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md", "names a commit-specific sentinel", ".agents/skills/bench-craft-delegate/SKILL.md dropped the fix-pass snapshot sentinel precondition")
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
	requireCollapsed(".agents/commands/bench-write-spec.md", "stale-command-reference sweep remains fail-closed across staged specs", ".agents/commands/bench-write-spec.md dropped the staged-spec fail-closed command sweep posture")
	requireCollapsed(".agents/commands/bench-what-next.md", "use `bench spec history <slug>` for the shipped-row check", ".agents/commands/bench-what-next.md dropped the bench spec history shipped-row check")
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
	requireCollapsed(".bench/BENCH.md", "ship only in the Bench kit repository", ".bench/BENCH.md does not state that the kit-maintenance surfaces ship only in the kit repository")
	requireCollapsed(".bench/BENCH.md", "a linked repo upgrades with `bench upgrade`", ".bench/BENCH.md does not name bench upgrade as the consumer's route onto a newer kit")
	requireCollapsed(".bench/BENCH-reference.md", "`.agents/skills/bench-craft-synthesis/SKILL.md` (kit-only)", ".bench/BENCH-reference.md skills index does not mark the kit-only craft-synthesis row")

	requireCollapsed(".bench/BENCH.md", "Parked ideas land in `capture/IDEAS.md`", ".bench/BENCH.md Capture section does not name capture/IDEAS.md as the capture sink")
	requireCollapsed(".bench/BENCH.md", "append the dated line (`- YYYY-MM-DD <text>`) to `capture/IDEAS.md`", ".bench/BENCH.md Capture section lost the no-PATH fallback append to capture/IDEAS.md")

	// Shared-rule placement — checkSharedRuleSingleSource's marker list owns each
	// sentence's presence and non-duplication; these anchors add only placement,
	// pinning the rule inside the section that owns it.
	if body, ok := scopedSection(".bench/BENCH.md", "Workflow"); ok {
		requireInSection(body, fixDontParkMarker, ".bench/BENCH.md Workflow section dropped the fix-don't-park rule; a mid-work defect fix belongs in the active workflow, not the backlog")
	}
	if body, ok := scopedSection(".bench/BENCH.md", "How to talk to me"); ok {
		requireInSection(body, sourceWarrantMarker, ".bench/BENCH.md How to talk to me section dropped the outside-source warrant rule; a claim resting on a source outside the tree names what was and was not read")
	}
	forbid(".bench/BENCH.md", "thorough", ".bench/BENCH.md phrases the outside-source warrant rule as thoroughness; the rule asks for disclosure of what went unread, which the reviewer can check — thoroughness nobody can")
	requireCollapsed(".agents/commands/bench-write-spec.md", "promote-then-delete commit removes the spec's `ROADMAP.md` row", ".agents/commands/bench-write-spec.md does not remove the spec's ROADMAP.md row in the promote-then-delete commit (row presence is status)")
	requireCollapsed(".agents/commands/bench-shape-idea.md", "never pause for permission or a re-prompt", ".agents/commands/bench-shape-idea.md dropped the resume-mode grill continuation rule; a running grill carries into newly-unblocked tickets without pausing for a re-prompt")
	requireCollapsed(".agents/commands/bench-review-implementation.md", "For an active spec build, submit the bounded receipt with `bench spec build review <slug> --evidence <receipt>`.", ".agents/commands/bench-review-implementation.md dropped exact-candidate review receipt submission")
	requireCollapsed(".agents/commands/bench-review-implementation.md", "Accepted findings become ownership-fenced repair tickets and return to `/bench-implement-spec`; a clean or risk-accepted review proceeds to `bench spec build promote <slug>`.", ".agents/commands/bench-review-implementation.md routes active spec-build findings or a clean review outside the lifecycle")
	forbid(".agents/commands/bench-review-implementation.md", "a clean or risk-accepted review goes to `/bench-final-check`", ".agents/commands/bench-review-implementation.md routes an active reviewed candidate to final-check before promotion")
	implementSpec := strings.ToLower(collapseSpace(stripHTMLComments(readIfExists(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md")))))
	implementSpec = strings.ReplaceAll(implementSpec, "`", "")
	if strings.Contains(implementSpec, "for an accepted repair finding, the coordinator may instead write the repair directly to the working branch before promote") {
		diags = append(diags, "bench-implement-spec permits an accepted repair to bypass provisional assignment and write directly to the working branch")
	}
	if strings.Contains(implementSpec, "not parallelizable") {
		diags = append(diags, "bench-implement-spec permits a generic unused-slot reason outside the closed set")
	}

	// Ticket guidance is convention-only, so these anchors pin the load-bearing
	// workflow clauses without inventing a parser for ticket files.
	requireCollapsed(".agents/commands/bench-implement-spec.md", "Charge `craft-tickets` before the first implementation edit", ".agents/commands/bench-implement-spec.md dropped the craft-tickets breakdown charge")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "derive ticket files from the spec's stories and seams", ".agents/commands/bench-implement-spec.md dropped the ticket breakdown derivation from the spec's stories and seams")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "`specs/<slug>/tickets/`", ".agents/commands/bench-implement-spec.md dropped the ticket breakdown destination")
	requireCollapsed(".agents/commands/bench-implement-spec.md", "under the session's existing approval surface", ".agents/commands/bench-implement-spec.md dropped the ticket breakdown approval surface")

	requireCollapsed(".bench/BENCH.md", "one independently-green ticket", ".bench/BENCH.md dropped the light-path independently-green-ticket observable")
	requireCollapsed(".bench/BENCH.md", "crosses no declared seam", ".bench/BENCH.md dropped the light-path declared-seam observable")

	requireCollapsed(".agents/skills/bench-craft-line/SKILL.md", "| Orchestration | mid + medium |", ".agents/skills/bench-craft-line/SKILL.md dropped the orchestration mid/medium stage default")
	requireCollapsed(".agents/skills/bench-craft-line/SKILL.md", "| Ticket implementation | cheap + low |", ".agents/skills/bench-craft-line/SKILL.md dropped the ticket implementation cheap/low stage default")
	requireCollapsed(".agents/skills/bench-craft-line/SKILL.md", "| Review (axis or falsification) | mid + high |", ".agents/skills/bench-craft-line/SKILL.md dropped the review mid/high stage default")

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

	requireCollapsed(".agents/skills/bench-craft-tickets/SKILL.md", "`craft-spec` owns the spec-time **who-writes-where** fence",
		".agents/skills/bench-craft-tickets/SKILL.md dropped the craft-spec ownership-fence cross-pointer")
	requireCollapsed(".agents/skills/bench-craft-spec/SKILL.md", "`craft-tickets` owns the build-time **what-lands-green-next** unit",
		".agents/skills/bench-craft-spec/SKILL.md dropped the craft-tickets build-time-unit cross-pointer")

	// The taught ticket shape is what the specbuild parser accepts, so each field
	// anchors to the template block that a cold author copies rather than to the
	// file: an appendix restating a field would satisfy a whole-file needle while
	// leaving the copied block unparseable. The prohibition anchors to the section
	// owning the gate cadence for the same reason.
	const ticketsSkill = ".agents/skills/bench-craft-tickets/SKILL.md"
	ticketTemplateRequires := []struct{ needle, diag string }{
		{"- [ ] [AB1] <observable behavioral criterion>", "dropped the labeled single-line acceptance row from the ticket template"},
		{"- [ ] [AB2] <observable behavioral criterion>", "dropped the second labeled acceptance row from the ticket template"},
		{"Ownership fence: `<path prefix>`, `<path prefix>`", "dropped the one-line backticked ownership fence from the ticket template"},
		{"Assumptions: <clause>; <clause>", "dropped the semicolon-separated assumptions line from the ticket template"},
		{"Blocked by: <sibling ticket file basenames, or none>", "dropped the basename-keyed blocked-by line from the ticket template"},
		{"| criterion | mutation | owner | operation sequence | |---|---|---|---| | <ID> |", "dropped the red-mutations table from the ticket template"},
	}
	if body, ok := scopedSection(ticketsSkill, "Write one file per ticket"); ok {
		for _, a := range ticketTemplateRequires {
			requireInSection(body, a.needle, ticketsSkill+" "+a.diag)
		}
	}
	// The breakdown procedure and the classification it invokes are two sections,
	// so each rule anchors to the one a coordinator reads it from: sizing and
	// sequencing rules belong beside the expand–migrate–contract list, while the
	// branch, the concurrency method, the one-line ceiling, and the cadence rules
	// belong inside the numbered method that applies them.
	breakdownRequires := []struct{ needle, diag string }{
		{"A wide refactor takes the expand–migrate–contract sequence instead of ordinary grouping",
			"dropped the blast-radius classification branch from the breakdown method's first step"},
		{"Concurrent eligibility is fence disjointness: two tickets run at once only when their ownership fences share no path.",
			"dropped fence disjointness as the mechanical concurrent-eligibility check beside the independently-green rule"},
		{"Name every real blocker by sibling ticket file basename.",
			"dropped the basename-keyed blocker naming from the breakdown method's third step"},
		{"A one-line change pays at most one shared test-harness line: below that ceiling it takes no fresh worktree, no fresh delegate, and no full gate by default.",
			"dropped the one-line test-harness ceiling beside the independently-green rule"},
		{"names which command authors gate evidence — `bench gate`, the canonical producing entry — and which phase consumes it",
			"dropped the evidence-authorship rule from the ticket cadence paragraph; a cadence-changing ticket names the producing command and the consuming phase"},
		{"The ticket carries behavioral acceptance checkboxes, not a project-gate checkbox: the green landing commit is the one source for that verdict.",
			"dropped the gate-checkbox prohibition from the ticket cadence paragraph"},
	}
	if body, ok := scopedSection(ticketsSkill, "Draft the breakdown"); ok {
		for _, a := range breakdownRequires {
			requireInSection(body, a.needle, ticketsSkill+" "+a.diag)
		}
		forbidInSection(body, "by sibling ticket title",
			ticketsSkill+" names blockers by ticket title in the breakdown method; a title dies at the next retitle, and the basename is what `--ticket` already names")
	}
	classifyRequires := []struct{ needle, diag string }{
		{"each batch sized by exactly one ownership fence",
			"dropped the one-ownership-fence sizing rule for migrate batches"},
		{"The contract ticket's `Blocked by:` names every migration ticket basename",
			"dropped the rule that the contract ticket's Blocked by names every migration ticket"},
	}
	if body, ok := scopedSection(ticketsSkill, "Classify before slicing"); ok {
		for _, a := range classifyRequires {
			requireInSection(body, a.needle, ticketsSkill+" "+a.diag)
		}
	}
	// The contracts-discovery rules anchor to the step that runs them — between
	// the drafted breakdown and the first ticket file — rather than to the file:
	// the same sentences written under the template would satisfy a whole-file
	// needle while leaving the slicing-time step a coordinator actually reads
	// silent about what crosses each fence.
	contractsRequires := []struct{ needle, diag string }{
		{"names four facts: its type, its membership or domain rule, its ordering, and its absence semantics",
			"dropped the four facts every fence-crossing value names in the contracts-discovery step"},
		{"asserted against the real producer and the whole enumerated family",
			"dropped the real-producer-and-enumerated-family assertion target from the consumer-ticket contract row"},
		{"When neither side can assert an invariant alone, add a junction ticket that can.",
			"dropped the junction-creation half of the junction rule from the contracts-discovery step"},
		{"a junction row discovered more than one ticket downstream moves a narrower copy of the row to the junction",
			"dropped the downstream-copy half of the junction rule from the contracts-discovery step"},
		{"from the tree after earlier tickets land — never from the spec's account of the base",
			"dropped the re-derive-claims-from-the-tree rule from the contracts-discovery step"},
	}
	if body, ok := scopedSection(ticketsSkill, "Discover the contracts before writing files"); ok {
		for _, a := range contractsRequires {
			requireInSection(body, a.needle, ticketsSkill+" "+a.diag)
		}
	}

	// The charge duties reach a low-context delegate only from the section it is
	// written in, so each needle is scoped to the section that owns the duty: the
	// self-probe, the probe-site rule, the template's probe-kind vocabulary, and
	// registry tracing belong to the charge a coordinator writes; backup isolation
	// belongs to the worktree rules the delegate runs under. A whole-file needle
	// would be satisfied by the same sentence written anywhere in the skill.
	const delegateSkill = ".agents/skills/bench-craft-delegate/SKILL.md"
	chargeRequires := []struct{ needle, diag string }{
		{"require the delegate to apply it to its own finished work, report the observed result, and add the missing row when the mutation comes back silently green",
			"dropped the delegate self-probe duty from the charge"},
		{"differs in site from every probe the delegate ran",
			"lets the coordinator probe repeat a site the delegate already probed"},
		{"report the observed result and the mutation's kind (omission or swap)",
			"dropped the omission/swap probe-kind vocabulary from the charge template"},
		{"names every registry the family already appears in, traced from one existing sibling",
			"dropped the registry-tracing duty from a family-extending charge"},
	}
	if body, ok := scopedSection(delegateSkill, "The charge"); ok {
		for _, a := range chargeRequires {
			requireInSection(body, a.needle, delegateSkill+" "+a.diag)
		}
	}
	if body, ok := scopedSection(delegateSkill, "Isolation"); ok {
		requireInSection(body,
			"the copy lives inside the delegate's own worktree under a unique name, and every restore names exact files, never a glob",
			delegateSkill+" dropped worktree-local backup isolation or admitted a glob restore")
	}

	// The spec-side halves anchor to the section a spec author reads each from:
	// the contracts pointer to the section that draws fences, the lifecycle class
	// to the run a walk enumerates, and the profile entry to the checklist that
	// walk consults. A whole-file needle would accept the same sentence parked
	// anywhere in the file, leaving the section that runs the rule silent.
	for _, scoped := range []struct{ rel, section, needle, diag string }{
		{".agents/skills/bench-craft-spec/SKILL.md", "Slicing a build for delegates",
			"Each fence carries value contracts across it, and `craft-tickets` owns naming them in `Discover the contracts before writing files`; this section points at that step by name rather than restating what it requires.",
			".agents/skills/bench-craft-spec/SKILL.md dropped the contracts-discovery pointer from the slicing section"},
		{".agents/skills/bench-craft-spec/SKILL.md", "The edge inventory",
			"re-run idempotency, process-boundary lifecycle, hostile environment",
			".agents/skills/bench-craft-spec/SKILL.md dropped the process-boundary lifecycle class from the canonical edge-class run"},
		{"projects/benchkit.md", "Hostile-input checklist (shell CLI)",
			"state serialized by one process and reloaded by a fresh one",
			"projects/benchkit.md dropped the process-boundary lifecycle entry from the hostile-input checklist"},
	} {
		if body, ok := scopedSection(scoped.rel, scoped.section); ok {
			requireInSection(body, scoped.needle, scoped.diag)
		}
	}

	for _, anchor := range []struct{ rel, needle, diag string }{
		{".agents/commands/bench-implement-spec.md", "`start` → `assign` → `checkpoint` → `integrate` → `review` → `promote`; `status` inspects the run and `abandon` plans or applies cleanup.", "bench-implement-spec dropped or reordered the eight-operation spec-build lifecycle"},
		{".agents/commands/bench-implement-spec.md", "Re-derive the complete ready frontier and the harness's live capacity before dispatch. Assign every ownership-safe ticket up to the smaller of frontier size and available capacity.", "bench-implement-spec dropped initial frontier capacity dispatch"},
		{".agents/commands/bench-implement-spec.md", "Refill the ownership-safe frontier after every integration or assignment release while another delegate remains active.", "bench-implement-spec replaced continuous frontier refill with drain-then-refill cadence"},
		{".agents/commands/bench-implement-spec.md", "For every unused harness slot, record exactly one reason: dependency, overlapping ownership fence, unavailable harness capacity, or measured resource constraint.", "bench-implement-spec dropped the closed unused-slot reason set"},
		{".agents/commands/bench-implement-spec.md", "Submit focused delegate evidence plus the coordinator-owned, different-kind probe through `checkpoint`.", "bench-implement-spec dropped focused evidence or the coordinator-owned different-kind probe"},
		{".agents/commands/bench-implement-spec.md", "Review the exact candidate composition before `promote`.", "bench-implement-spec moved composed review after promotion"},
		{".agents/commands/bench-implement-spec.md", "Accepted findings become new ownership-fenced repair tickets and re-enter `assign`, `checkpoint`, and `integrate` before a fresh composed review.", "bench-implement-spec routes an accepted repair outside the provisional lifecycle"},
		{".agents/commands/bench-implement-spec.md", "When the branch tip moves, `promote` is the operation that recomposes the run onto the new tip, and recomposition discards the review.", "bench-implement-spec dropped moved-tip recomposition through promote or its review discard"},
		{".agents/commands/bench-implement-spec.md", "The repair round is therefore repair → `promote` → `review` → `assign` … `integrate` → `review` → `promote`.", "bench-implement-spec dropped the ordered moved-tip repair round"},
		{".agents/commands/bench-review-implementation.md", "For an active spec build, read `bench spec build status <slug> --full` and bind the review inputs to the exact candidate subject and recorded run base it reports. Confirm that subject is unchanged immediately before receipt submission; a changed candidate invalidates the review rather than letting a delta review authorize a new composition.", "bench-review-implementation dropped exact-candidate review input binding"},
		{".agents/commands/bench-implement-spec.md", "After a green promotion, run `/bench-final-check` only for its terminal retained-evidence report and implementation retro.", "bench-implement-spec orphaned final-check's terminal report and implementation retro after promotion"},
		{".agents/commands/bench-implement-spec.md", "Do not run `bench commit` for a provisional spec-build ticket", "bench-implement-spec restored a per-ticket whole-project gate to provisional spec builds"},
		{".agents/skills/bench-craft-delegate/SKILL.md", "The coordinator probe's mutation kind differs from the delegate author's mutation kind.", "craft-delegate allows the coordinator probe to repeat the author's mutation kind"},
		{".agents/skills/bench-craft-delegate/SKILL.md", "The lifecycle checkpoints, integrates, and releases the assignment; the coordinator does not run a generic release.", "craft-delegate dropped the spec-build checkpoint-integrate-release exception"},
		{".agents/skills/bench-craft-delegate/SKILL.md", "A provisional checkpoint is not project-green evidence and cannot satisfy a done-claim.", "craft-delegate lets provisional evidence claim project green"},
		{".bench/BENCH.md", "`bench spec build start|assign|checkpoint|integrate|review|status|promote|abandon`", ".bench/BENCH.md dropped an operation from the spec-build inventory"},
		{".bench/BENCH.md", "Provisional cadence is exclusive to reviewed spec-backed builds; light-path work, `bench shift`, and ordinary `bench commit` remain commit-on-green.", ".bench/BENCH.md broadened provisional cadence beyond reviewed spec-backed builds"},
		{".bench/BENCH-reference.md", "| `start` | create or resume the subject-bound run from exact-tip whole-tree green, including a narrow verdict whose inherited evidence still covers every skip |", "BENCH-reference misroutes spec build start"},
		{".bench/BENCH-reference.md", "| `assign` | lease one ownership-fenced ticket worktree |", "BENCH-reference misroutes spec build assign"},
		{".bench/BENCH-reference.md", "| `checkpoint` | validate focused evidence and bind a provisional commit |", "BENCH-reference misroutes spec build checkpoint"},
		{".bench/BENCH-reference.md", "| `integrate` | compare-and-swap one verified checkpoint into the candidate |", "BENCH-reference misroutes spec build integrate"},
		{".bench/BENCH-reference.md", "| `review` | bind three-axis evidence to the exact candidate |", "BENCH-reference misroutes spec build review"},
		{".bench/BENCH-reference.md", "| `status` | inspect durable state and retained evidence |", "BENCH-reference misroutes spec build status"},
		{".bench/BENCH-reference.md", "| `promote` | gate and publish the exact reviewed composition;", "BENCH-reference misroutes spec build promote"},
		{".bench/BENCH-reference.md", "a moved tip recomposes through `promote`, discarding the review |", "BENCH-reference dropped promote's moved-tip recomposition from the lifecycle lookup"},
		{".bench/BENCH-reference.md", "| `abandon` | plan or apply recoverable cleanup |", "BENCH-reference misroutes spec build abandon"},
		{"bin/bench.sh", "bench spec build assign <slug> --ticket <ticket> --request <id>", "bench help dropped or malformed spec build assign grammar"},
		{"bin/bench.sh", "bench spec build start <slug>", "bench help dropped spec build start grammar"},
		{"bin/bench.sh", "bench spec build checkpoint <slug> --assignment <id> --evidence <receipt>", "bench help dropped or malformed spec build checkpoint grammar"},
		{"bin/bench.sh", "bench spec build integrate <slug> --assignment <id>", "bench help dropped or malformed spec build integrate grammar"},
		{"bin/bench.sh", "bench spec build review <slug> --evidence <receipt>", "bench help dropped or malformed spec build review grammar"},
		{"bin/bench.sh", "bench spec build status <slug> [--full]", "bench help dropped spec build status grammar"},
		{"bin/bench.sh", "bench spec build promote <slug>", "bench help dropped spec build promote grammar"},
		{"bin/bench.sh", "bench spec build abandon <slug> [--apply <fingerprint>]", "bench help dropped abandon plan/apply grammar"},
		{"projects/benchkit.md", "gpt-5.6-sol / high", "benchkit profile replaced the approved spec-build guidance line"},
		{"projects/benchkit.md", "Both dogfood traces use the public porcelain", "benchkit profile dropped the public-porcelain dogfood traces"},
		{"projects/benchkit.md", "Run `bench structure` before and after the guidance cut", "benchkit profile dropped the spec-build structure preflight"},
		{"CHANGELOG.md", "Light-path changes, `bench shift`, and ordinary `bench commit` keep commit-on-green cadence.", "CHANGELOG dropped the unchanged-path control for provisional spec builds"},
	} {
		requireCollapsed(anchor.rel, anchor.needle, anchor.diag)
	}
	requireCollapsed(".agents/commands/bench-implement-spec.md", "A repair pass integrates the findings accepted for its round, performs one fresh composed review, and stops at the next promotion result; another semantic review round opens only when the composition changes or the reviewer requests one.",
		".agents/commands/bench-implement-spec.md dropped the terminal repair-pass bound")
	for _, rawRoute := range []string{"Create the checkpoint with `git commit`", "advance the candidate with `git update-ref`", "create the assignment with `git worktree`", "replay the patch with `git cherry-pick`"} {
		forbid(".agents/commands/bench-implement-spec.md", rawRoute, "bench-implement-spec synthesizes lifecycle Git plumbing outside the eight public operations")
	}

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

	requireCollapsed(".agents/commands/bench-final-check.md", "only after `bench spec build status <slug> --full` reports a terminal promoted run and its retained exact green evidence",
		".agents/commands/bench-final-check.md dropped the terminal promoted implementation-retro placement")
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
	requireCollapsed(".agents/commands/bench-final-check.md", "rewrite `capture/retros/<spec-slug>.md` in full",
		".agents/commands/bench-final-check.md dropped whole-file implementation-retro replacement")
	requireCollapsed(".agents/commands/bench-final-check.md", "Do not run another gate or commit just to capture the retro",
		".agents/commands/bench-final-check.md added an implementation-retro gate outside the normal cadence")
	requireCollapsed(".bench/BENCH.md", "/bench-final-check` writes `capture/retros/<spec-slug>.md`",
		".bench/BENCH.md dropped the implementation-retro capture owner")

	requireCollapsed(".agents/commands/bench-what-next.md", "The snapshot's `retros` bodies are the only retro evidence this run reads",
		".agents/commands/bench-what-next.md dropped the roadmap-context retro drain source")
	requireCollapsed(".agents/commands/bench-what-next.md", "merge into an existing roadmap row, a new roadmap row, a learning-or-rule disposition, or an explicit dismissal",
		".agents/commands/bench-what-next.md dropped an implementation-retro recommendation disposition")
	requireCollapsed(".agents/commands/bench-what-next.md", "remove every drained `capture/retros/*.md` file in the same reviewer-approved batch",
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
		{"Concrete defects — bugs, spec misses, missing coverage — return through repair assignment, checkpoint, and integration without stopping; contestable design and judgment findings are flagged in the exit report for reviewer veto, not applied",
			"dropped a finding-disposition half (provisional repair cadence, or flag-don't-apply judgment)"},
		{"Submit the delegate's bounded receipt through `bench spec build review` only while its candidate subject still matches the exact reviewed composition",
			"dropped exact-subject review submission from the full orchestration"},
		{"Then run **final-check inline** only to report the retained terminal evidence and capture the retro",
			"dropped terminal final-check reporting after promotion"},
		{"bounded by this command's terminal repair-pass bound and routed through `craft-delegate`'s repair allowance",
			"dropped the repair-bound pointers; this mode adds no second version of either rule"},
		{"At every phase boundary the run writes the phase reached into `capture/session-handoff.md`'s State section, then refreshes the pin block with `bench handoff --next <command>`",
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
		{"Every phase claim in the exit report cites the record that proves it: the review delegate's invocation, lifecycle status and retained evidence, and the `capture/session-handoff.md` boundary rewrite",
			"dropped the record-citation rule for phase claims in the exit report"},
		{"ask a second, separate question: whether to add a cross-harness falsification pass over the diff before promotion — the Codex CLI at the top binding, charged to refute the claim that the spec was implemented rather than to grade it against the three axes",
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
		{"stops on a concrete defect", "stops the run on a concrete defect; concrete defects return through the provisional repair cadence"},
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
