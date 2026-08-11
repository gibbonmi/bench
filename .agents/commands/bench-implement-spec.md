---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /bench-write-spec, or for a change the lighter-path threshold in .bench/BENCH.md lets you build without a spec. Declares its line and uses TDD at seams. An opt-in `--full <spec>` run carries the spec through implementation, review, the final green landing, and terminal reporting — resuming from capture/session-handoff.md.
---

# /bench-implement-spec — do the work at the seams

## Entry orientation

This is the implementation phase. It starts from an approved spec or a tiny change
with an obvious seam, declares the line, works vertical slices at the agreed seams,
and uses the acceptance coverage map to keep the build target fixed.
A spec-backed run opens with `bench preflight build <slug>` before the first
edit; a red preflight stops the phase.

## Exit handoff

Close by reporting the implemented stories, the acceptance coverage status for
each row, and the landed commits. A reviewed spec-backed build closes through
`/bench-review-implementation` over the composed diff, then the final green
landing commit and `/bench-final-check`; a light-path build keeps the existing
`/bench-review-implementation` then `/bench-final-check` route. A build that stops
short exits through "When the build stops short" below and recommends its one
durable next action.

Implement the spec at the seams it names. If there's no spec, the change must fall
under the lighter-path threshold `.bench/BENCH.md`'s "Right-size the process"
paragraph sets; if it doesn't, stop and run `/bench-write-spec` first.

Pre-build research scales to the spec's staleness and vagueness, never to its
story count. A recently staged spec with named seams gets **one** read-only
delegate charged narrowly: confirm the named seams still hold in the tree,
and return the fixture-and-seam inventory (`craft-delegate` names it) for
each fence the build will charge. Drift found there routes back to
`/bench-write-spec` with the finding quoted — it is never compensated for
with a wider build-time research fan-out; a spec with no seams to confirm is
underspecified, which is the same route.

## First derive the tickets

Charge `craft-tickets` before the first implementation edit. For a spec-backed
build, derive ticket files from the spec's stories and seams, write them under
`specs/<slug>/tickets/`, and present the breakdown as the build plan under the
session's existing approval surface. For the light path, apply the same
procedure to the unspecced change and write its one ticket.

With the ticket files written and nothing assigned, run the breakdown review
`craft-tickets` owns: one fresh read-only delegate grades the ticket breakdown
before any assignment; a harness that cannot spawn one runs the pass inline and
flags it in the build plan. Its findings are reslices, repaired before the
first ticket is assigned rather than through a repair round after code lands.

Fence drift takes the same route and no ticket is assigned without a complete handoff ledger.
Send seam or fence drift through `$bench-write-spec`, obtain approval for the
repaired spec, then require `craft-tickets`' completed breakdown review before
assigning any ticket; `craft-tickets` owns the ledger schema.

Complete `craft-tickets`' contract discovery before assignment. Nothing
mechanical validates ticket shape or fences for you: an incomplete ticket
surfaces as a composed red several tickets later, so repair it before it is
assigned.

Tickets land serially, commit-on-green. Work respects `Blocked by:` order — a
ticket starts only after every blocker has landed — and each ticket's verified
diff lands on the branch with path-scoped `bench commit` before the next
ticket's landing.

Light-path work takes the same landing: one independently
green ticket, focused checks, and path-scoped `bench commit` as its atomic
full-project-gate landing.

## Open with the line

Declare the line before touching code — the declaration template, the decision
table that picks the row, and the escalation ladder are all `craft-line`'s.

## Route the venue

Pick the venue before the first implementation edit; it is part of the phase
contract, not a preference:

- Every spec-backed run assigns genuine write work to at least one write subagent
  before the first implementation edit — the reviewer never has to ask. How each
  delegate is charged, isolated, and verified is `craft-delegate`'s; its model and
  effort are `craft-line`'s.
- Route by slice shape: independent vertical slices fan out to separate parallel
  subagents within the harness's concurrency limit; dependent slices run
  sequentially; a spec that lands as one atomic diff is delegated whole to one
  worktree-isolated write subagent, and the invoking session still verifies the
  returned done-claim through the `craft-delegate` contract. Each routed slice
  is an independently-attributed ticket with its own fresh charge.
- A read-only helper (research, review, planning, search) does not satisfy the
  write requirement. When any change may remain inline, and what a harness that
  cannot spawn a write subagent does instead of editing, are `craft-delegate`'s
  capability-aware delegation policy — this phase states no inline threshold of
  its own. For the resume handoff that policy requires, this phase's invocation
  is `$bench-implement-spec` in Codex, `/bench-implement-spec` in Claude Code,
  or the canonical `.agents/commands/bench-implement-spec.md` file in another
  subagent-capable AGENTS harness.

## Then build

- Before the first slice, check for `reviews/<spec-slug>.md`. When it exists, its
  findings are part of this build's target: resolve them alongside the stories,
  and the green fix commit that closes them names and deletes the file, so
  resolved findings never resurface as pickup work.
- Work the user stories in vertical slices, not all-tests-first horizontal ones.
- Use TDD only at the pre-agreed seams; its bounds are `craft-tdd`'s.
- If the spec has an acceptance coverage map, seed your harness's native task
  list (Codex plan, Claude todos), if it has one, from `bench coverage <spec>` —
  one task per coverage row — so the reviewer can watch build progress. Each
  vertical slice names the coverage row it is turning red-to-green before
  editing that slice, and marks that row's task in-progress as it does. Rows
  marked `already covered` or `not TDD-able` keep their recorded reason; don't
  silently upgrade them into TDD coverage.
- Each write delegate runs the ticket's focused checks and records its own
  mutation probe. The coordinator independently probes the returned tree through
  a different mutation kind — independent evidence, not a replay or
  transcription of the delegate's pass — then lands the verified diff with
  path-scoped `bench commit -m "<msg>" <path>...`, which gates and commits
  atomically, and releases the worktree.
- During structure housekeeping, apply `craft-seams`' split-or-grant rule.
- One small change at a time, repo stays green — invariant 4 in `.bench/BENCH.md`.
- Every delegation during the build carries its own line and, when the spec has
  a coverage map, its stories' coverage rows in the charge — every time; the
  rules are `craft-delegate`'s (model half: `craft-line`). Every write-delegate's
  charge also names `craft-comments` by path
  (`.agents/skills/bench-craft-comments/SKILL.md`) — comment discipline rides
  the charge, not the delegate's priors.
- Verify each returned done-claim against its charged rows and report the
  round in one line: landed, or the missed case and the re-charged repair.
- Verify a claim against the tree as it stands, not against memory; a claim
  over a whole set is verified by enumerating the set, never by extending one
  measured member.
- Route code authorship through `craft-delegate`: the coordinator scopes, routes,
  and verifies; a write-delegate authors the code. That skill owns worktree
  isolation and verification discipline. Use `craft-line` to choose each
  delegate's model and effort.
- The venue comes from "Route the venue" above — spec-backed work is never
  inlined here. When an atomic build collapses the stories into one delegate
  charge, flag the collapsed per-story lines in the exit report.
- For broad renames or reference refactors, dry-run the file scope before editing,
  then verify old stems in every form: `/name`, `$name`, bare basenames in
  inventories, and `dir/name` path forms. Separator slashes inside prose are not
  command invocations.

## When a delegate is blocked outside its fence

A write delegate whose deterministic repro fails on a defect outside its
ticket's ownership fence stops implementation edits at that boundary. It does
not run `/bench-debug` — that phase is the reviewer's to invoke, never an
agent's. The delegate returns a bounded blocked report instead of a done-claim:
the exact repro command with its red exit and output digest, the failing surface
it observed, its assignment ID, recorded base, and in-fence dirty paths. It
never edits outside its fence and never absorbs the defect into its own diff — a
fix smuggled through the adapter's fence is a review defect even when it makes
the ticket green.

The coordinator relays that report and stops, recommending the reviewer run the
debug phase in this harness's own form against the named failing surface. That
reviewer-run pass produces the debug receipt — the blocked report's evidence plus
the confirmed cause, the paths the repair must own, and whether the ticket can
proceed once the repair lands. Nothing below runs until the receipt exists.

`craft-tickets` is the sole repair-reslicing owner. The coordinator charges
`craft-tickets` with the validated receipt and applies its guidance. The receipt
already fixes the cause, required fence, and proceed condition, so it takes one
or more repair tickets, never a small spec. Each ticket adds its own ownership
fence, at least one acceptance row, and each row's red mutation.
One repair ticket remains the common case. This command owns only the resulting
tickets' orchestration.

The coordinator authors every repair ticket file; the implementing delegates
never author their own acceptance criteria. Having them draft criteria for
coordinator approval adds no independent source, only a round trip.

With the receipt and coordinator-authored tickets in hand, the repair takes
the ordinary serial cadence: write every repair ticket under
`specs/<slug>/tickets/`, land that planning batch with path-scoped
`bench commit`, then implement and land each repair ticket commit-on-green in
`Blocked by:` order. Resume the blocked ticket only after the whole repair has
landed and the receipt's proceed condition is satisfied; its delegate rebases
onto the repaired tree and completes through the ordinary landing.

## When the build stops short

A build that exhausts its iteration cap or ends with stories unmet exits through a
defined route — never a silent grind, never an abandoned worktree:

1. **Report state:** stories done vs. remaining, the coverage table as it
   stands, the gate verdict, and what consumed the cap.
2. **Keep what's real:** landed tickets stay landed; unfinished work remains in
   its owned worktree. Nothing gets squash-finished to fake completion.
3. **Route by cause, and recommend one:**
   - wrong tier (the model ground, the gate disagreed) → re-declare one tier up
     per the `craft-line` ladder and resume this phase;
   - wrong spec (a story is unbuildable as written) → back to `/bench-write-spec`
     with the finding quoted;
   - wrong scope (the spec is bigger than one build) → propose the split; the
     reviewer decides.
4. **Leave the worktree owned:** every cut worktree keeps its owner — the
   coordinator that cut it owns its retirement, through `bench worktree release`
   or `bench worktree clean`.

## Close on green

- Once every implementation ticket short of the final landing has landed, run
  `/bench-review-implementation` over the composed diff — the semantic
  three-axis pass (Standards + Spec + Coverage) that catches what the gate can't:
  right thing built the wrong way, wrong thing built cleanly, or breaking inputs
  nothing exercises. Review before the final landing, not after.
- Accepted findings become new ownership-fenced repair tickets and land through
  the same serial commit-on-green cadence before a fresh review of the changed
  diff.
- **The terminal repair-pass bound.** A repair pass lands the findings accepted
  for its round, performs one fresh review, and stops at that result; another
  semantic review round opens only when the diff changes or the reviewer
  requests one.
- Before the final landing, emit the coverage table for every acceptance row —
  `bench coverage <spec>` produces it and `bench coverage --check <spec>` validates
  the map; don't hand-assemble it. Classify each row `green`, `already covered`,
  or `not TDD-able`. If any mapped behavior is missing, partial, or unclassified,
  the build is not ready to land.
- The final green landing commit is path-scoped
  `bench commit --spec <slug> -m "<msg>" <path>...`: it runs the gate, commits
  only on green, and is the sole author of the spec's `Status: implemented`
  flip. Every earlier ticket lands the same way without `--spec`; ordinary
  `bench commit` keeps its gate-then-commit contract everywhere.
- Then summarize what changed in plain language and hand back. I own the merge;
  propose it, don't perform it.

For UI work, if your project has an interaction-layer skill and a screenshot loop,
they're part of the gate alongside the `craft-design-system` skill — a green test suite is
necessary but not sufficient for UI.

## The `--full` run

`--full` is an opt-in flag: plain `/bench-implement-spec` keeps today's
implement-only semantics, and `/bench-review-implementation`,
`/bench-final-check`, and `/bench-debug` stay standalone commands for strict
phased use and mid-run resumption. `--full` invoked with no spec argument, or
with one naming a path that does not exist, refuses at the entry contract and
says which of the two it was — a run that guessed its target would build the
wrong spec with every phase green.

The run carries the named spec from build to push-ready by orchestrating the
standalone phases; their contracts stay theirs:

- **Implement through ticket delegates**, under everything above this section.
- **Review in a delegate:** spawn one fresh-context delegate charged with the
  standalone `/bench-review-implementation` contract and given the spec and
  the diff and nothing else. Inline self-review is closed, not deprioritized —
  the context that produced the code carries the assumptions that produced its
  bugs — and a context-inheriting delegate is the same failure wearing a
  delegate's name. The delegate's done-claim answers to invariant 1 in
  `.bench/BENCH.md` and `craft-delegate`'s verification discipline. When the
  review delegate cannot run or returns nothing, the run stops and reports at
  that boundary rather than proceeding to final-check with review unrun.
- **Land inline:** the final green landing commit — path-scoped
  `bench commit --spec <slug>` — is the composed gate boundary. Then run
  **final-check inline** to report the landed evidence and capture the retro. When an
  issue needs deep analysis, the run does not debug inline: it stops at that
  boundary, reports the failing surface, and recommends the reviewer invoke
  `/bench-debug` themselves in this harness's form. That phase is expensive and
  reviewer-invoked by design — no agent charges it.

**Finding disposition.** Concrete defects — bugs, spec misses, missing
coverage — return through repair tickets landed commit-on-green without
stopping; contestable design and
judgment findings are flagged in the exit report for reviewer veto, not
applied. Re-gating follows the repair-ticket cadence: focused evidence, a
green landing, then a fresh review of the changed diff. The repair pass is
bounded by this command's terminal repair-pass bound and routed through
`craft-delegate`'s repair allowance; this mode adds no second version of
either.

**The scope fence.** The run implements the spec's stories and nothing else:
work noticed outside them — an adjacent refactor, an unrelated improvement, a
story the spec chose not to take — is recorded for the reviewer rather than
built. The fence bounds the run's diff to the spec's stories and seams; it
counterweights the fix-don't-park route, and whether a discovered defect is
fixed or parked stays `.bench/BENCH.md`'s Workflow section's test.

**Phase boundaries persist.** At every phase boundary the run writes the phase
reached into `capture/session-handoff.md`'s State section, then refreshes the pin
block with `bench handoff --next <command>`. The phase reached is the one fact
this mode adds; every other field is the existing handoff contract's. A
re-invoked `--full` resumes from the phase the handoff names instead of
re-implementing from the top; a stale handoff is arbitrated by `AGENTS.md`'s
tree-wins rule.

**Escalation asks first.** When the diff is large enough that the mid binding
could miss important bugs, pause and ask the reviewer as a structured decision
list with a recommendation for this run, offering three routes: continue at
the mid binding; escalate to the top binding in this harness; escalate to the
top binding via the Codex CLI. Those three are the fixed menu — a route is
omitted only when this harness cannot invoke it at all, and the omission is
stated rather than silent. Never escalate without asking. A harness with no
structured-prompt surface asks the same question as a plain numbered list.
The Codex route is `codex exec`: set the model, the reasoning effort, the
working directory, and a non-interactive approval posture — `codex exec
--help` is authoritative for the flag spellings, so none are pinned here.

**The falsification offer.** On that same size trigger, ask a second, separate
question: whether to add a cross-harness falsification pass over the diff
before the final landing — the Codex CLI at the top binding, charged to refute the
claim that the spec was implemented rather than to grade it against the three
axes. It is its own question, not a fourth route in the escalation menu, and
it never runs standing: absent the trigger it is not offered. When this
harness cannot invoke Codex the pass is omitted and the omission stated, the
same posture the escalation menu takes.

**The exit report accounts and cites.** The report gives one disposition per
row of the spec's acceptance coverage map — implemented, deferred, or
won't-handle — named row by row against `bench coverage <spec>`'s enumeration.
When the spec carries no coverage map, the report says so and accounts for the
user stories instead. Every phase claim in the exit report cites the record
that proves it: the review delegate's invocation, the landed commits with
their gate verdicts, and the `capture/session-handoff.md` boundary rewrite.
