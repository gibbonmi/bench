# Align bench with the Pocock workflow shape

Status: ready

## Destination

Bench keeps its oracles (gate authority, three-axis review, worktree isolation,
capture/drain, handoff) but inverts its doctrine-to-process ratio to match the
Pocock workflow: lifecycle machinery removed or gutted, completeness
enumeration moved upstream into a domain-modeling discipline, reviews
re-derive from primary sources, and the kit prose shrinks to doctrine leaves.
Done when the reviewed decision set below is resolved and ready for one or
more specs.

## #1: What replaces the provisional spec-build lifecycle?

Blocked by: none
Type: Grill

### Question

Remove `bench spec build` (assign/checkpoint/integrate/review/promote,
receipts, refresh, fences-as-enforcement, ~11.7k lines in
`internal/specbuild`) entirely in favor of serial commit-on-green tickets on a
real branch — or gut it to a thin status tracker? What survives: anything?

### Answer

Full removal (reviewer, 2026-08-11). Nothing survives — no thin
status-tracker wrapper. Serial commit-on-green per ticket on a feature branch;
parallel delegates only for disjoint tickets that merge to the branch
immediately on green. `bench spec` keeps only `staged`/`implemented` status;
`bench commit`'s gate-then-commit contract is the sole landing path.

## #2: Disposition of the seven staged lifecycle specs and in-flight repair tickets

Blocked by: #1
Type: Grill

### Question

`axi-spec-build-complete` (+24 repair tickets), `axi-coherent-diff`,
`axi-query-disclosure`, `axi-compatibility-oracle`, `checkpoint-scoped-review`,
`ticket-bundle-refusal`, `single-build-serial-gate` all extend the machinery
#1 may remove. Retire, rescope, or build which?

### Answer

Retire four, park three (reviewer, 2026-08-11). Retired: `axi-spec-build-complete`
(with its 24 repair tickets and the uncommitted ticket edits),
`checkpoint-scoped-review`, `ticket-bundle-refusal` (its no-author-asserted-
keep-together rule survives as a sentence in the craft-tickets rewrite), and
`axi-compatibility-oracle` (baseline pinned to a pre-removal subject; re-derive
fresh after the surface stabilizes if still wanted). Parked to the roadmap for
re-ranking after the reshape: `axi-coherent-diff`, `axi-query-disclosure`
(both need rescoping off spec-build surfaces and the retired `axi.Action`
landing route), `single-build-serial-gate` (independent; promote references
rot). The reshape is the only active program.

## #3: Which parsed ticket fields survive?

Blocked by: #1
Type: Grill

### Question

Closure tokens, Contracts, Integration surfaces, covers annotations, and
red-mutation tables are consistency validators. Cut to Pocock shape (what to
build / acceptance / blocked-by) — or keep any field as machine-parsed?

### Answer

Only `Blocked by:` stays parsed (reviewer, 2026-08-11) — the one field with a
real consumer (frontier computation). Closure, Contracts, Integration
surfaces, covers annotations, and the red-mutation table are deleted as
schema; a contract that matters becomes a What-to-build sentence plus an
acceptance checkbox, re-derived by review from the tree (#6). Ownership fence
demotes from enforcement whitelist to an advisory Writes: note the coordinator
uses for parallel-disjointness judgment, never for refusal or review
narrowing. Ticket shape: title, Blocked by:, What to build, Acceptance.

## #4: Adopt a domain-modeling discipline upstream of the spec?

Blocked by: none
Type: Grill

### Question

Add a domain-modeling skill plus living CONTEXT.md glossary (challenge terms,
edge-case scenarios at concept boundaries, code-vs-claim cross-referencing,
inline updates) as the completeness-enumeration home — and where does it fire:
inside shaping, its own phase, or ambient?

### Answer

Companion skill, no hard gate (reviewer, 2026-08-11). New `bench-craft-domain`
skill mirroring Pocock's domain-modeling leaf: glossary challenge, canonical
terms with Avoid lists, edge-case scenarios at concept boundaries,
code-vs-claim cross-referencing, inline CONTEXT.md updates, glossary-only
content. Charged by craft-grill, /bench-shape-idea, and /bench-write-spec;
everything else reads CONTEXT.md ambiently. Enumeration duty (input
partitions, hostile classes, family membership) lives here as domain
scenarios walked with the reviewer during shaping. craft-adr stays, aligned
to the paragraph-ADR standard. No refusal gate on spec writing.

## #5: Which doctrine leaves does the kit adopt?

Blocked by: none
Type: Grill

### Question

Candidates: tests.md/mocking.md-style example leaves under craft-tdd;
DEEPENING dependency categories and DESIGN-IT-TWICE parallel radical designs
under craft-seams; a /prototype skill; frontier-rounds grilling replacing pure
one-question-at-a-time in craft-grill. Which land, and which stay out?

### Answer

All four (reviewer, 2026-08-11): craft-tdd gains tests.md and mocking.md
example leaves referenced by path; craft-seams gains DEEPENING.md
(dependency-category table → test strategy) and DESIGN-IT-TWICE.md (parallel
radically-different interface designs for spec-phase exploration); a new
/prototype skill backs the existing Prototype ticket type (throwaway, trivial
to run, no persistence, surface the state, capture verdict then discard);
craft-grill moves to frontier-rounds grilling — whole settled-prerequisite
frontier per numbered round, each question with a recommended answer, serial
only when answers genuinely reshape the next question.

## #6: Review re-derivation mandate

Blocked by: #1
Type: Grill

### Question

Reviews currently validate declared artifacts. Mandate independent
re-derivation from primary sources (full fence set from the spec, input
partition from the producer, ancestry from main) — and spend review
parallelism on independent derivations rather than fenced concurrent writing?

### Answer

Adopted (reviewer, 2026-08-11). Re-derive-then-compare becomes the review
skill's core rule: three axes stay as parallel fresh-context delegates;
Coverage derives the complete input partition from the actual producer and
the authorized write set from the spec, then diffs both against the
candidate; Spec quotes spec lines against behavior, never ticket claims;
Standards keeps the Fowler baseline. A finding cites what it re-derived; a
review that only confirmed declarations is incomplete by definition.

## #7: The executable preflight — contents and home

Blocked by: #1
Type: Grill

### Question

Mechanical reality checks: base is current main; every written path
spec-authorized; every spec row owned by a ticket; producer-derived membership
equals inventory; rev-parse + non-empty diff before review. Which land, and
where — the gate, a phase-entry check, or both?

### Answer

All five checks, as one `bench preflight` command run at phase entry — not in
the gate (reviewer, 2026-08-11). The gate stays the done-oracle over the
tree; preflight is the start-oracle over artifacts-vs-reality. Build entry
and review entry invoke it; a red preflight stops the phase.

## #8: HITL gates versus delegate passes

Blocked by: none
Type: Grill

### Question

Restore reviewer-in-the-loop gates — seam confirmation before TDD, a
ticket-breakdown quiz iterated with the reviewer — replacing or augmenting the
current delegate breakdown review? Which boundaries get the human?

### Answer

Both gates, replacing the delegate breakdown review (reviewer, 2026-08-11).
(a) No TDD test at an unconfirmed seam — spec sign-off already covers
spec-backed work, so this bites only light-path work. (b) Ticket breakdown
presented as a numbered list (title / blocked-by / what it delivers) and
iterated with the reviewer until approved. The batch-approval AFK carve-out
in BENCH.md stays.

## #9: The kit prose diet

Blocked by: #1, #3
Type: Grill

### Question

Target line budgets after the removals: craft-tickets 411→~100,
bench-implement-spec 389→~50, plus the always-loaded BENCH.md/AGENTS.md
surface. What ratio and budget does the reviewer set, and what is the audit
rule that keeps it?

### Answer

Hard budgets in the profile (reviewer, 2026-08-11): craft-tickets ≤ 100
lines, bench-implement-spec ≤ 60, every other skill ≤ 120 with doctrine split
into on-demand reference leaves, always-loaded BENCH.md ≤ 150. The gate's
conformance sweep gains a per-file line-budget check; a breach is red, and
raising a budget is a reviewer edit to the profile.

## #10: Fate of craft-line and delegation charging

Blocked by: #1
Type: Grill

### Question

Does tier/line governance and the delegation discipline survive the lifecycle
removal, and in what form?

### Answer

Keep both, slimmed (reviewer, 2026-08-11). The line declaration stays as cost
governance; craft-delegate shrinks to its durable core — fresh context per
ticket, worktree isolation for writes, verify the done-claim against tree and
gate — shedding the receipt/charge ceremony. The `no-mistakes` upstream
(kunchenguid) is not this ticket's concern: per the kit's recorded assessment
its transferable part is the fail-closed default and the finding-action
vocabulary (no-op / auto-fix / ask-user), which land in #6's finding
disposition and #7's fail-closed preflight posture, not in line governance.

## #11: CLI shrink scope beyond spec build

Blocked by: #1
Type: Grill

### Question

Which `bench` subcommands shrink or go with the machinery?

### Answer

Remove what serviced the lifecycle: recovery-ref machinery, `spec build
reclaim`, provisional-ref plumbing (reviewer, 2026-08-11). Keep `worktree`
(create/path/exec/release/clean), `gate`, `commit`, `status`, `guards`, and
the capture/roadmap surfaces. Exact inventory is spec-writer discretion
within that rule.

## #12: Migration path for linked repos

Blocked by: #2
Type: Grill

### Question

How do linked repos consume the removal?

### Answer

No backwards-compatibility work at all (reviewer, 2026-08-11): removed
subcommands simply become unknown commands; no shims, no compatibility
messaging beyond a changelog line. Ordinary `bench upgrade` syncs kit files;
a linked repo mid-spec-build lands or discards provisional work before
upgrading, noted in the changelog.

## #13: How does the reshape split into specs?

Blocked by: #1, #5, #7
Type: Grill

### Question

The shaped scope holds several independently useful behaviors: lifecycle
removal (code + CLI + ticket schema + workflow prose), doctrine adoption (new
skills, leaves, review rewrite, HITL gates, prose diet), and the preflight
command. One bundled spec, or a split — and in what order?

### Answer

Three specs plus immediate housekeeping (reviewer, 2026-08-11). Housekeeping
now, no spec: retire the four specs, park the three roadmap rows, discard the
uncommitted repair-ticket edits. Then Spec A — lifecycle removal
(#1, #3, #11, #12): delete `internal/specbuild`, shrink the CLI, reduce the
ticket schema, cut lifecycle references from workflow prose. Then Spec B —
`bench preflight` (#7). Then Spec C — doctrine adoption (#4, #5, #6, #8, #9,
#10) as one kit-prose batch.

## Not yet specified

## Spec-writer discretion

- Exact CLI subcommand inventory within #11's keep/remove rule.

## Out of scope

- Gate authority (invariant 1) — not weakened by any ticket here.
- Dropping the Coverage axis to match Pocock's two-axis review.
- Tracker-backed collaborative maps (closed in skills-assessment.md).

## Sources

- Path: `skills-assessment.md`
  Supports: #5 — prior v1.1 comparison; adoptions 1–10 already applied.
  Drift: assessed upstream v1.1; re-check the upstream ref before Spec C.
- Path: `capture/learnings.md`
  Supports: #1, #2, #7 — lifecycle failure evidence, with the retired run's 24 repair tickets.
  Drift: the Pocock skills repo sits outside this tree at `~/workspace/skills` (read 2026-08-11 for the nested-skill evidence behind #1, #4, #5, #6, #8); re-read before Spec C.
