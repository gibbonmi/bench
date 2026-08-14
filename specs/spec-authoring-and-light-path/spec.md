# spec-authoring-and-light-path

Status: staged

Decision source: reviewer-confirmed conversation, 2026-08-14 (light-path shape, slicing move, two-loop authorship verification, --reviewer grammar, drain-time implement-now, and the seam-pause drop each closed by explicit reviewer answer).

Verification log: spec 5 + tickets 3 iteration(s) to accept — reading the enforcement surfaces (fixture EXPECTs, MUTATE strings, bespoke count checks) by content before drafting rows, instead of by name, would plausibly have made it 2 + 1; the residual round was a ticket-sequencing cycle.

## Problem

The workflow is heavier than the risk it manages in three places. An
under-threshold change still pays a `craft-tickets` charge, a breakdown-approval
pause, a write-delegate, and a worktree. Ticket slicing happens in
`/bench-implement-spec`, so the session that holds the spec's decision context
hands off before the tickets exist. And spec authoring is routed to a fresh
mid-tier session, discarding the shaping session's context, while the
falsification pass is a single fire-and-forget iteration with no way to choose
the reviewer model or effort — and cross-harness reviewer invocations burn turns
rediscovering CLI syntax.

## Solution

Lighten the light path to write-one-ticket-then-implement-inline. Move ticket
slicing into `/bench-write-spec` so one approval table covers spec and tickets.
Let the main session author both, verified by two uncapped mid-tier review
loops — spec first, breakdown second — with a recorded per-loop iteration count
and a faster-approve note, overridable with `--reviewer <tier-or-model>
[effort]`. Pin the exact cross-family CLI recipes and the
own-family-uses-native-agents rule. Let a `/bench-what-next` drain implement
light-path-eligible items on the spot through a verified write-delegate instead
of parking them on the roadmap. Repair the four stale anchor rows found while
reading the enforcement surface.

## User stories

1. As a reviewer, when a change meets the light-path observables (one
   independently-green ticket, crosses no declared seam), the worker writes the
   one ticket file and implements it inline in the same session — no breakdown
   approval pause, no write-delegate, no worktree, no stop-for-seam-confirmation
   before the first TDD test — then gates and commits on green. `craft-tickets`
   gains the matching carve-out so its reviewer-approved-breakdown and
   one-write-delegate-charge rules are explicitly scoped to spec-backed builds.
   Line: fable / high. Kit guidance prose takes the `craft-line` leverage override.
2. As a reviewer, `/bench-write-spec` authors `specs/<slug>/tickets/` after its
   spec loop accepts and one approval table covers spec and breakdown;
   `/bench-implement-spec` consumes existing tickets at entry and routes a
   spec-backed run with no tickets back to `/bench-write-spec` instead of
   slicing there. Line: fable / high. Leverage override.
3. As a reviewer, the session holding the decision source authors the spec at
   whatever tier it runs, and verification runs as two loops: loop 1 — a
   read-only same-family mid-tier delegate at high effort reviews the spec
   with the falsification questions, author-fix/re-review uncapped until no
   blocking findings — then the session slices the tickets from the verified
   spec, then loop 2 — the same delegate shape grades the breakdown only
   against `craft-tickets` rules, uncapped until clean. My sign-off stays the
   hard stop, and the spec records `Verification log: spec <n> + tickets <m>
   iteration(s) to accept — <note>` naming what would have reached accept in
   one pass. When either count exceeds 1, the author also appends a
   `capture/learnings.md` entry — which stage missed, what the delegate
   caught, why, and the proposed rule or skill change — so the existing
   `/bench-what-next` drain turns repeat misses into process edits.
   Line: fable / high. Leverage override.
4. As a reviewer, `--reviewer <tier-or-model> [effort]` overrides the
   verification delegate: a tier resolves through the invoking harness's own
   `lines.env` column, a model id must already be bound in `lines.env` (an
   unbound id is refused), an own-family id runs through the harness's native
   agent surface (never that family's CLI), and a cross-family id runs through
   the pinned CLI recipe so no turns are spent on syntax; a harness with no
   native subagent surface falls back to its own family's CLI.
   Line: fable / high. Leverage override.
5. As a reviewer, when a `/bench-what-next` drain verdicts an item that meets
   the light-path observables, I can choose "implement now" instead of a
   `ROADMAP.md` row: the session writes the one ticket file, spawns a
   write-delegate charged with it, verifies the returned diff against the
   ticket's acceptance rows and the gate, and lands it as its own commit on
   green — recorded as a second named exception to the drain's
   one-batch-commit rule, beside the existing two-spec-retire exception. Items
   needing a reviewer decision, a new seam, or spec-level design still
   graduate to the roadmap. Line: fable / high. Leverage override.
6. As a maintainer, the four stale anchor rows match the tree again: the
   benchkit hostile-input needle matches its real casing, the review-preflight
   needle matches the current entry sentence, the dead `shared-build-cache
   opt-in` row is retired, and the `new session on the mid tier` row is
   retired with story 3's replacement landing in its own ticket. Bundled into
   this spec by reviewer decision under the Fix-don't-park rule rather than
   shipped alone. Line: opus / medium. The profile caches only the effort for
   conformance work; opus is this harness's mid binding, so the model is
   derived, not cached.

## Implementation decisions

- **Light path** (`.bench/BENCH.md` right-size table, `craft-tdd`,
  `craft-tickets`, `craft-delegate`): the route becomes "write the one ticket
  file (`craft-tickets` owns the template), then implement it inline in this
  session" — the table remains the standing approval; both threshold
  observables and the `Right-size the process` marker are retained verbatim.
  `craft-tdd`'s light-path bullet replaces the stop-for-seam-confirmation
  with: the ticket file names the test seam and the reviewer vetoes post-hoc.
  `craft-tickets` gains a light-path carve-out — the right-size table is the
  ticket's standing approval and it implements inline — so its "only route
  onto the frontier" and "one fresh write-delegate charge" rules bind
  spec-backed builds only (displacing lines; the file sits at its 100-line
  budget). `craft-delegate`'s lighter-path inline allowance wording stays the
  allowance it already is.
- **Slicing move** (`bench-write-spec.md`, `bench-implement-spec.md`,
  `craft-tickets`): write-spec gains a slicing step that runs after
  the spec's verification accepts — charge `craft-tickets`, write
  `specs/<slug>/tickets/`, and carry the numbered title/`Blocked by:`/outcome
  list into the approval table. Landing order: the slicing step lands first
  wording its trigger against the existing falsification review (true of the
  file at that commit); the verification-loops ticket re-points that trigger
  to loop 1 when it installs the loops — no commit ships a reference to a
  step the file does not contain. BENCH.md's workflow list gains the same
  fact in one clause: step 2 reads "lock stories, seams, and gate
  expectations, and slice the tickets". In implement-spec, remove only the two slicing
  clauses (the `Charge \`craft-tickets\` ...` breakdown charge and its
  reviewer-approved-breakdown / AFK-carve-out continuation — one grammatical
  sentence in the live file); the three
  anchored delegation sentences stay verbatim — the write-subagent assignment,
  the read-only-helper disqualifier, and the `craft-delegate`
  incapable-harness clause — while the integration-worktree sentence (which
  has no registry row today) is reworded only as far as its now-orphaned
  "After approval" referent, naming the write-spec-phase approval, and gains
  a first-time Require row and fixture whose diagnostic must not carry the
  `workflow integration source: ` prefix. Its entry validates tickets
  exist (preflight's present-tickets path already handles them — no production
  Go change) and routes a ticketless spec-backed run to `/bench-write-spec`.
  `craft-tickets`' trigger text moves from build entry to spec authoring
  (its `index:` frontmatter stays, so `.bench/BENCH-reference.md` needs no
  edit); template, breakdown, and frontier rules for spec-backed builds
  unchanged. write-spec's Exit handoff sentence is rewritten — sign-off, then
  a fresh mid-tier build session on one retained integration source, with
  slicing no longer routed through `/bench-implement-spec` — landing as a
  Forbid+Require pair on its anchored post-slicing-handoff needle, its
  backing fixture updated or created. The Forbid's diagnostic must not begin
  with the `workflow integration source: ` prefix — a bespoke helper counts
  that family at an exact size and fatals on a second row per file, and its
  test file stays unfenced.
- **Authorship + two verification loops** (`bench-write-spec.md`,
  `projects/benchkit.md` Lines): "Who runs this phase" says the session holding
  the decision source authors; step 9 becomes two verification loops
  run by a read-only delegate, default same-family mid at high effort via the
  native agent surface: loop 1 takes the spec alone with the existing
  falsification questions, and only after it accepts does the slicing step
  run; loop 2 then takes the breakdown alone against `craft-tickets` rules.
  Each loop repeats author-fix/re-review until no blocking findings. The
  uncapped loops are a recorded reviewer exception to invariant 2's iteration
  cap; each round is still reported in one line. On close the author writes
  `Verification log: spec <n> + tickets <m> iteration(s) to accept — <note>`
  into the spec (template gains the line under `Decision source:`); a count
  above 1 in either loop also appends one `capture/learnings.md` entry — which
  stage missed, what was caught, why, the proposed rule change — so feedback
  drains through the existing reviewed path with no new ledger surface. Human
  sign-off stays the hard stop. The rewrite keeps the `Bootstrap authority
  before execution` rule at exactly its current two occurrences (a bespoke
  check counts them). benchkit's Lines rows for spec authoring, the
  falsification pass, and the ticket-breakdown review pass are rewritten to
  this state; the harness×tier binding table is untouched, so `line-routing`
  stays green. `bench-shape-idea.md`'s exit recommendation drops the
  fresh-mid-tier routing and recommends `/bench-write-spec` from the session
  holding the ready decision source, as a Forbid+Require pair with its
  fixture updated or created.
- **--reviewer + recipes** (`bench-write-spec.md`, `craft-delegate` +
  new `references/` file): grammar `--reviewer <tier-or-model> [effort]`
  (`--reviewer mid xhigh` under Codex resolves `BENCH_CODEX_MID` at xhigh).
  Refusal of unbound ids matches the check-agent-line hook. Recipes live in
  `.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md`
  (unbudgeted; the skill-dir symlink already exposes it to Claude Code):
  `claude -p --model <id> --effort <level> "<charge>"` and
  `codex exec --sandbox read-only -m <id> -c model_reasoning_effort=<level>
  "<charge>"`, verified against the installed CLIs at build time, plus the
  no-native-surface fallback rule. Own-family → native agent surface, never
  the CLI (flagged for veto: the fallback keeps Codex able to run an
  own-family reviewer). craft-delegate's own-family rule and pointer displace
  equivalent lines (the file sits at its 120-line budget).
- **Drain-time light path** (`bench-what-next.md`, `.bench/BENCH.md` capture
  wording): the drain verdict set gains "implement now" for items meeting the
  light-path observables — reviewer chooses it per item during the reviewed
  drain, the session writes the one ticket file, spawns a write-delegate
  (`craft-delegate` isolation and `craft-line` routing apply), and verifies
  the returned diff against the ticket's acceptance rows and the gate. The
  implemented item lands as its own commit on green, named in
  bench-what-next.md as the second exception to the drain's one-batch-commit
  rule (beside two-spec-retire); a second Require row pins the exception with
  a **new** backing fixture, the existing `one uncommitted batch diff` needle
  staying as-is — no fixture backs that rule today. BENCH.md's "graduate into ROADMAP.md only
  through a reviewed `/bench-what-next` drain" stays true — an implemented item
  closes instead of graduating. Note the deliberate asymmetry: the interactive
  light path implements inline; the drain-time route delegates because the main
  session is coordinating the drain and owns verification.
- **Anchor discipline**: every changed anchored clause updates its
  `internal/anchors/registry_data.go` row and its
  `tests/canary/workflow-guidance-anchors/` fixture in the same ticket, and
  **every replaced clause lands as a pair** — a Forbid on the retired text
  plus a Require on the successor — so the additive cheat (new clause beside
  the old) reds. A third enforcement surface exists:
  `TestWorkflowCadenceAnchorsRejectDeletionAndSwap`
  (`internal/conformance/fixture_bite_test.go`) hardcodes byte-exact
  substrings, line wraps included, from bench-write-spec.md, craft-tickets,
  craft-delegate, benchkit.md, and craft-spec (not an edit target here),
  failing unless each occurs exactly once.
  Any edit that reflows a pinned substring updates that test's string in the
  same ticket; edits near pinned strings otherwise keep them byte-identical
  — the test file is fenced for exactly this. Existing fixtures this spec
  rewrites: `craft-tdd-light-path-seam-gate`,
  `write-spec-review-made-conditional`, `write-spec-review-tier-escalated`,
  `benchkit-spec-ownership`. The other implement-spec, ticket, and write-spec
  fixtures (`implement-spec-mandatory-delegation-anchor`,
  `implement-spec-inline-exception`, `write-spec-handoff-anchor`,
  `ticket-stage-routing-anchor` — a craft-line fixture — and
  `ticket-light-path-anchor`, whose variant pins the retained `crosses no
  declared seam` observable) pin clauses this spec does not change and must
  stay green untouched. New clauses each get a Require row and a **new**
  fixture — no successors exist for them: main-session authorship (replacing
  the retired `new session on the mid tier` row), inline light path,
  craft-tickets carve-out, slicing-in-write-spec, implement-spec entry
  validation, the reworded worktree sentence, verification-log line,
  learnings hook, --reviewer grammar, recipes, own-family rule, implement-now
  verdict, the drain commit exception, and the shape-idea exit
  recommendation. Separately, the three delegation anchors and the two
  repaired stale needles keep their existing Require rows — no duplicate rows
  — and gain first-time fixtures. The build re-derives the complete affected-fixture set by
  sweeping every changed phrase across `tests/canary/` before landing. Needles
  at risk of being satisfied by quoted examples elsewhere in the same file use
  `RequireInSection`. Multi-clause route sentences are pinned by one long
  needle covering the whole sentence, so dropping any clause fails the
  substring.
- **Fixture census stays independently green per ticket**: every ticket that
  adds or removes a top-level canary fixture updates the checked-in expected
  binding count in `internal/canary/inventory_test.go` in the same commit. The
  independently authored literal remains the omission sentinel; it is not
  derived from the producer under test.
- **Prose budgets bind the edits**: `.bench/BENCH.md` ≤ 180 (now 177),
  `bench-implement-spec.md` ≤ 60 (deletion frees room), `craft-tickets` ≤ 100
  (at limit) and `craft-delegate` ≤ 120 (at limit) — additions displace
  lines; `craft-tdd` has headroom (117/120); `bench-write-spec.md`,
  `bench-what-next.md`, `benchkit.md`, and the references file are unbudgeted.
- **Drift surfaces**: `docs/field-guide.html` passages restating "fresh
  mid-tier session" and slicing-at-implement are updated in the same change
  (under-threshold, no deferral); the edit near the retained
  frontier-ticket sentence must leave that anchored sentence byte-identical.
  No gate observable exists for the field-guide passages themselves — the
  update is review-checked prose, recorded in the edge inventory.
  `docs/reporesident-distillation.md`'s light-path trigger sentence stays
  accurate and is left alone.

## Testing decisions

- A good test here is the existing fixture-bite contract: materialize the
  broken variant → the owner check must emit the fixture's EXPECT diagnostic;
  restore the live file → that diagnostic must vanish. Both halves are
  exercised by `go test ./internal/conformance -run
  TestEveryRetainedFixtureBitesThroughRegisteredOwner`, which the gate's test
  phase runs.
- Seams receiving tests: the anchors registry + canary fixture family (prior
  art: every fixture under `tests/canary/workflow-guidance-anchors/`). No new
  seams; no production Go behavior changes.
- Gate seam observing the feature: `bench gate` — conformance checks
  `docs-currency-workflow` (fixture bites), `guidance-prose-budgets` (line
  budgets), `skills-index-command-adapters` (index/adapters), plus
  `bench coverage --check` across staged specs.

### Seam diagram

    trigger: gate test phase (go test ./internal/conformance)
        │
        ▼
    guidance file + registry row  ──▶  [ anchors.EvaluateGroup via check owner ]  ──▶  diagnostics
    fixture broken variant        ──▶  [ fixture bite: materialize → red,      ]
                                       [ restore live → that red vanishes      ]
                      ◀ tests attach here: fixture EXPECT written first, so the
                        restore-half stays red until the live clause lands

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| WF1 | 1 | BENCH.md light-path route reads write-ticket-then-implement-inline, no approval pause/delegate/worktree, observables and marker retained; the old route cell text is retired | fixture bite (new light-path-inline fixture; `ticket-light-path-anchor` stays green untouched) with a Forbid on the old route text | observed red: restore-half red until the live route lands | dropped observable or surviving old-route text leaves a diagnostic standing after restore |
| WF2 | 1 | craft-tdd light-path bullet: ticket names the seam, reviewer vetoes post-hoc, no live stop; the stop-for-confirmation clause is retired | fixture bite (craft-tdd-light-path-seam-gate rewritten) pairing Forbid on the retired stop with Require on the successor | observed red: restore-half red until the bullet is replaced | the additive cheat — new bullet beside the old stop — trips the Forbid |
| WF3 | 2 | bench-write-spec.md charges `craft-tickets` after the spec's verification (trigger worded at landing against the current falsification review, whose verdict is advisory; re-pointed to loop 1 by the loops ticket), names `specs/<slug>/tickets/`, folds the breakdown into the approval table, and its Exit handoff no longer routes slicing through `/bench-implement-spec` while keeping the retained-integration-source build recommendation; BENCH.md workflow step 2 carries the slice clause and field-guide's slicing-at-implement passage is corrected (both read-observable, no fixture) | fixture bite (new slicing-in-write-spec fixture; registry rows moved from the implement file; the post-slicing-handoff needle rewritten as Forbid+Require) | observed red: restore-half red until the slicing step and rewritten handoff land | write-spec without the post-verification slicing step, or still routing slicing to implement, fails a needle |
| WF4 | 2 | bench-implement-spec.md drops only the two slicing clauses, keeps the three anchored delegation sentences verbatim, rewords only the worktree sentence's "After approval" referent to the write-spec-phase approval, validates tickets at entry, and routes a ticketless spec-backed run to `/bench-write-spec` | fixture bite (Forbid on the old charge + Require on the route-back; the three delegation anchors and the reworded worktree sentence gain first-time fixtures — the worktree sentence also gains its first row, diagnostic outside the `workflow integration source: ` prefix) | observed red: restore-half red for the route-back and reworded-worktree needles; the three delegation fixtures red on the bite half, their clauses being already live | a retained charge trips the Forbid; a missing route-back trips the Require; over-deleting any of the four sentences reds a new fixture |
| WF5 | 3 | bench-write-spec.md says the session holding the decision source authors spec and tickets; the fresh-mid-tier-authors wording is retired, and field-guide's fresh-mid-tier authoring passages are corrected (read-observable, no fixture) | fixture bite (new authorship fixture; the `new session on the mid tier` row is retired and a Forbid pins out the old wording) | observed red: restore-half red until the replacement lands | the additive cheat — authorship clause beside the old routing — trips the Forbid |
| WF6 | 3 | two-loop verification clause: read-only same-family mid at high via native agent surface; loop 1 on the spec before slicing, loop 2 on the breakdown; each uncapped until no blocking findings; sign-off stays the hard stop; the one-iteration wording is retired | fixture bite (write-spec-review-made-conditional / -tier-escalated rewritten) pairing Forbid on the retired one-iteration text with Require on the loop clause | observed red: restore-half red until the two-loop clause replaces the old text | a combined, capped, sign-off-skipping, or additive wording fails a needle |
| WF7 | 3 | template and process require `Verification log: spec <n> + tickets <m> iteration(s) to accept — <note>` written into the spec at close | fixture bite (new verification-log fixture, RequireInSection on Template) | observed red: restore-half red until the template line lands | a spec process without the two-count log line fails the section-scoped needle |
| WF8 | 4 | `--reviewer <tier-or-model> [effort]` grammar: tier resolves through the invoking harness's own column, unbound ids refused | fixture bite (new reviewer-override fixture, long needle over the grammar sentence) | observed red: restore-half red until the grammar lands | dropping the flag, the tier-resolution clause, or the refusal breaks the long-needle substring |
| WF9 | 4 | recipes file carries both exact CLI commands and the no-native-surface fallback; craft-delegate pins own-family→native and points at the file | fixture bite (new `files/`-style recipes fixture on the references file + new craft-delegate pointer fixture) | observed red: restore-half red until file content and pointer land | a dropped command, fallback, pointer, or rule fails its needle |
| WF10 | 3 | benchkit.md Lines rows read main-session authorship, the two loops, and slicing-at-spec; `Spec falsification pass` needle retained; binding table untouched | fixture bite (benchkit-spec-ownership rewritten) + line-routing check | observed red: bite-half red by construction — the rewritten fixture's mutation string is re-derived from the new Lines text, which does not exist until the rows land | stale Lines rows contradict the shipped process for every cold session |
| WF11 | 6 | repaired needles (`Hostile-input checklist`, current preflight-review sentence) hold on the live tree and are fixture-enforced | fixture bite (two new fixtures) | observed red: restore-half red while the needle mismatches the tree | the case-mismatched and reworded needles are exactly what restore-half detects |
| WF12 | 6 | dead rows removed: `shared-build-cache opt-in` retired here, `new session on the mid tier` retired by WF5's ticket | root-conformance sweep (`TestRootConformance` with `BENCH_CONFORMANCE_ROOT` set — prep-release's run) | observed red: the sweep reports all four stale diagnostics today; the story-6 diagnostics must vanish after repair, with the ~8 pre-existing reds — unrelated in cause, though two sit on bench-implement-spec.md's missing section headings — recorded as inherited baseline | the sweep is the live-tree observable the ordinary gate lacks for unfixtured rows; before/after output lands in ticket evidence |
| WF13 | 1–5 | every budgeted file stays within its prose budget after the edits | guidance-prose-budgets check | already covered (gate conformance check) | an over-budget edit reds the gate without new tests |
| WF14 | 5 | bench-what-next.md offers the implement-now verdict for light-path-eligible drained items, its route sentence naming ticket file, write-delegate, and main-session verification against the ticket's acceptance rows | fixture bite (new implement-now fixture, one long needle over the route sentence) | observed red: restore-half red until the verdict route lands | dropping any clause of the route breaks the long-needle substring |
| WF15 | 3 | write-spec requires a `capture/learnings.md` entry when either loop's count exceeds 1 | fixture bite (new learnings-hook fixture) | observed red: restore-half red until the learnings clause lands | a process without the learnings hook fails the needle, so misses never reach the drain |
| WF16 | 1 | craft-tickets scopes reviewer-approved breakdown and one-write-delegate-charge to spec-backed builds and carves out the light-path ticket with the table as its standing approval | fixture bite (new carve-out fixture; the retained breakdown fixture stays green) | observed red: restore-half red until the carve-out lands | without the carve-out the skill contradicts BENCH.md's route and the contradiction is anchored red |
| WF17 | 5 | the implemented drain item lands as its own commit on green, named as the second exception to the drain's one-batch-commit rule | fixture bite (new batch-commit-exception fixture — no fixture backs the rule today) | observed red: restore-half red until the named exception lands | a drain folding the implementation into the batch diff or leaving it uncommitted violates the retained one-batch rule or the new exception needle |
| WF18 | 3 | bench-shape-idea.md's exit recommends `/bench-write-spec` from the session holding the ready decision source; the fresh-mid-tier routing is retired | fixture bite (shape-idea exit fixture updated or created) pairing Forbid on the retired routing with Require on the successor | observed red: restore-half red until the recommendation is replaced | the additive cheat or a surviving fresh-mid-tier recommendation trips a needle |
| WF19 | 1 | the lighten-light-path ticket's fixture additions update the independently authored canary binding count in the same commit | `TestRunReportsAcceptedInventoryBindings` | observed red: the ticket's two new fixtures make `go test ./internal/canary` report 185 while the expected count remains 183 | an omitted same-commit census update leaves the ticket gate red |
| WF20 | 2 | the move-slicing-into-write-spec ticket's fixture additions update the independently authored canary binding count in the same commit | `TestRunReportsAcceptedInventoryBindings` | add the ticket's top-level fixtures while retaining the previous expected count; `go test ./internal/canary` fails | an omitted same-commit census update leaves the ticket gate red |
| WF21 | 4 | the cross-harness-reviewer-recipes ticket's fixture additions update the independently authored canary binding count in the same commit | `TestRunReportsAcceptedInventoryBindings` | add the ticket's top-level fixtures while retaining the previous expected count; `go test ./internal/canary` fails | an omitted same-commit census update leaves the ticket gate red |
| WF22 | 5 | the drain-time-light-path ticket's fixture additions update the independently authored canary binding count in the same commit | `TestRunReportsAcceptedInventoryBindings` | add the ticket's top-level fixtures while retaining the previous expected count; `go test ./internal/canary` fails | an omitted same-commit census update leaves the ticket gate red |
| WF23 | 6 | the repair-stale-anchors ticket's fixture additions update the independently authored canary binding count in the same commit | `TestRunReportsAcceptedInventoryBindings` | add the ticket's top-level fixtures while retaining the previous expected count; `go test ./internal/canary` fails | an omitted same-commit census update leaves the ticket gate red |
| WF24 | 3 | the main-session-authorship ticket's fixture additions update the independently authored canary binding count in the same commit | `TestRunReportsAcceptedInventoryBindings` | add the ticket's top-level fixtures while retaining the previous expected count; `go test ./internal/canary` fails | an omitted same-commit census update leaves the ticket gate red |
| WF25 | 3 | the verification-loops ticket's fixture additions update the independently authored canary binding count in the same commit | `TestRunReportsAcceptedInventoryBindings` | add the ticket's top-level fixtures while retaining the previous expected count; `go test ./internal/canary` fails | an omitted same-commit census update leaves the ticket gate red |
| WF26 | 4 | the reviewer-override-flag ticket's fixture addition updates the independently authored canary binding count in the same commit | `TestRunReportsAcceptedInventoryBindings` | add the ticket's top-level fixture while retaining the previous expected count; `go test ./internal/canary` fails | an omitted same-commit census update leaves the ticket gate red |

### Edge inventory

- Quoted grammar token (a needle satisfied by example prose elsewhere in the
  file) — covered by WF7's `RequireInSection` and the same treatment for any
  ambiguous new needle (anchor-discipline build rule).
- Additive-cheat replacement (new clause appended beside the retired one) —
  covered by the Forbid+Require pairs in WF1, WF2, WF4, WF5, WF6.
- Over-deletion (removing more of an anchored paragraph than the spec names) —
  covered by WF4's new fixtures: the three delegation anchors are unfixtured
  today and the worktree sentence has no row at all, so the ticket creates
  all four to make over-deletion of any sentence observable in the ordinary
  gate.
- Absent vs empty file — already covered: the registry emits
  `anchor file missing` and the fixture harness refuses an empty EXPECT.
- Case sensitivity of Require matching — covered by WF11 (the defect class
  being repaired).
- Budget boundary at exactly the limit — covered by WF13; craft-delegate and
  craft-tickets sit at their limits, so any net addition reds.
- Non-ASCII whitespace in edited markdown — **Won't handle**: `CollapseSpace`
  uses `strings.Fields` (Unicode-aware), so an NBSP cannot unanchor a needle in
  this matcher; no new exposure is added.
- Special files / symlinks in touched paths — **Won't handle**: the new
  references file is a regular tracked file inside an already-symlinked skill
  directory; existing kit-compliance and budget refusal posture is unchanged.
- Stale field-guide prose surviving the edit — **Won't handle** beyond
  review: the passages are unanchored narrative HTML, and anchoring a doc
  site line-by-line is the separate capability already priced out of scope.
- Ticketless spec-backed run reaching implement anyway — **Won't handle**
  beyond the anchored route-back prose: `bench preflight build` legitimately
  treats absent tickets as its fresh-build case, so no gate observable exists
  without a production Go change, which is out of scope by decision.
- Interrupt / re-run idempotency / process-boundary lifecycle — **Won't
  handle**: docs-and-data change only; no process writes state.
- OpenCode-family `--reviewer` id — **Won't handle**: the column is unbound by
  closed decision, so the unbound-id refusal (WF8) is the defined behavior.

## Ownership fences

- `.bench/BENCH.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-what-next.md`
- `.agents/commands/bench-shape-idea.md`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `.agents/skills/bench-craft-tdd/SKILL.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md`
- `projects/benchkit.md`
- `docs/field-guide.html`
- `internal/anchors/registry_data.go`
- `internal/canary/inventory_test.go`
- `internal/conformance/fixture_bite_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `specs/spec-authoring-and-light-path/`

## Out of scope

- Wiring the existing root-conformance anchor sweep into the ordinary dev gate
  (the sweep exists — `TestRootConformance` under `BENCH_CONFORMANCE_ROOT`,
  run by prep-release — but the ordinary gate's test phase does not set that
  root, so unfixtured rows rot between releases) — separate capability,
  ~3 edits, 2 gate runs; parked to `capture/IDEAS.md` for a reviewer decision.
- OpenCode adoption or binding its `lines.env` column — separate capability,
  ~3 edits, 1 gate run.
- Any `bin/bench.sh` or production Go behavior change — preflight already
  accepts a present `tickets/` directory; nothing to build.
- Regenerating or restructuring `docs/field-guide.html` beyond the passages the
  workflow change falsifies — separate capability, ~2 edits, 1 gate run.
