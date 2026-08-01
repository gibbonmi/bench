---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /bench-write-spec, or for a change the lighter-path threshold in .bench/BENCH.md lets you build without a spec. Declares its line and uses TDD at seams. An opt-in `--full <spec>` run carries the spec through implementation, exact-candidate review, promotion, and terminal reporting — resuming from capture/session-handoff.md.
---

# /bench-implement-spec — do the work at the seams

## Entry orientation

This is the implementation phase. It starts from an approved spec or a tiny change
with an obvious seam, declares the line, works vertical slices at the agreed seams,
and uses the acceptance coverage map to keep the build target fixed.

## Exit handoff

Close by reporting the implemented stories, the acceptance coverage status for
each row, and the durable spec-build state. A reviewed spec-backed build closes
through `bench spec build promote`, then `/bench-final-check` reports its retained
terminal evidence and captures the retro; a light-path build keeps the existing
`/bench-review-implementation` then `/bench-final-check` route. A build that stops
short exits through "When the build stops short" below and recommends its one
durable next action.

Implement the spec at the seams it names. If there's no spec, the change must fall
under the lighter-path threshold `.bench/BENCH.md`'s "Right-size the process"
paragraph sets; if it doesn't, stop and run `/bench-write-spec` first.

## First derive the tickets

Charge `craft-tickets` before the first implementation edit. For a spec-backed
build, derive ticket files from the spec's stories and seams, write them under
`specs/<slug>/tickets/`, and present the breakdown as the build plan under the
session's existing approval surface. For the light path, apply the same
procedure to the unspecced change and write its one ticket.

For a reviewed spec-backed build, run `bench spec build start <slug>` before the
first assignment. The public lifecycle is `start` → `assign` → `checkpoint` →
`integrate` → `review` → `promote`; `status` inspects the run and `abandon`
plans or applies cleanup. These eight operations are the complete mutation
surface: the harness supplies readiness and evidence but does not reproduce the
lifecycle in Git commands.

Re-derive the complete ready frontier and the harness's live capacity before
dispatch. Assign every ownership-safe ticket up to the smaller of frontier size
and available capacity. For every unused harness slot, record exactly one
reason: dependency, overlapping ownership fence, unavailable harness capacity,
or measured resource constraint. Refill the ownership-safe frontier after every
integration or assignment release while another delegate remains active. A
one-shot fan-out followed by waiting for the frontier to drain is not this
cadence.

Light-path work does not enter `bench spec build`. It retains one independently
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
  is an independently-attributed frontier ticket with its own fresh charge.
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
  and the promoted composition that closes them names and deletes the file, so
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
  a different mutation kind, then assembles the bounded checkpoint receipt
  outside the assignment worktree. The coordinator probe is independent evidence,
  not a replay or transcription of the delegate's pass.
- Submit focused delegate evidence plus the coordinator-owned, different-kind
  probe through `checkpoint`.
- Submit that receipt through `bench spec build checkpoint`, then use
  `bench spec build integrate`; the lifecycle owns checkpoint commits,
  compare-and-swap integration, and assignment release. Do not run `bench commit`
  for a provisional spec-build ticket, and do not release an integrated assignment
  through generic worktree porcelain.
- The harness never runs `git commit`, `git update-ref`, `git worktree`, or
  patch-replay plumbing for this lifecycle. A public operation refusal is a
  lifecycle result to diagnose or re-enter, never permission to synthesize it.
- During structure housekeeping, apply `craft-seams`' split-or-grant rule.
- One small change at a time, repo stays green — invariant 4 in `.bench/BENCH.md`.
- Every delegation during the build carries its own line and, when the spec has
  a coverage map, its stories' coverage rows in the charge — every time; the
  rules are `craft-delegate`'s (model half: `craft-line`). Every write-delegate's
  charge also names `craft-comments` by path
  (`.agents/skills/bench-craft-comments/SKILL.md`) — comment discipline rides
  the charge, not the delegate's priors.
- Verify each returned provisional claim against its charged rows and report the
  round in one line: checkpointed, or the missed case and the re-charged repair.
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

## When the build stops short

A build that exhausts its iteration cap or ends with stories unmet exits through a
defined route — never a silent grind, never an abandoned worktree:

1. **Report state:** stories done vs. remaining, the coverage table as it
   stands, the gate verdict, and what consumed the cap.
2. **Keep what's real:** retained checkpoints and candidate integrations stay
   provisional; uncheckpointed work remains in its owned assignment. Nothing
   gets squash-finished to fake completion.
3. **Route by cause, and recommend one:**
   - wrong tier (the model ground, the gate disagreed) → re-declare one tier up
     per the `craft-line` ladder and resume this phase;
   - wrong spec (a story is unbuildable as written) → back to `/bench-write-spec`
     with the finding quoted;
   - wrong scope (the spec is bigger than one build) → propose the split; the
     reviewer decides.
4. **Leave the worktree owned:** a spec-backed run remains in durable lifecycle
   state; if the reviewer drops it, use `bench spec build abandon` and its
   fingerprinted apply. A light-path coordinator that cut a generic worktree
   still owns its retirement: the coordinator that cut it owns its retirement.

## Close on green

- Provisional checkpoints are evidence, not green commits. The exact composed
  candidate receives a fresh Standards, Spec, and Coverage review before any
  promotion attempt; submit its bounded receipt through `bench spec build review`.
- Review the exact candidate composition before `promote`.
- Accepted findings become new ownership-fenced repair tickets and re-enter
  `assign`, `checkpoint`, and `integrate` before a fresh composed review. A repair
  never writes directly to the working branch and a delta-only review cannot
  authorize the changed composition.
- `bench spec build promote` is the sole spec-backed whole-project-gate boundary.
  It constructs the prospective `Status: implemented` tree, runs the gate for
  that exact subject, and publishes the authorized squash only on green. A red
  leaves the run provisional and routes candidate-owned defects back to repair.
- After a green promotion, run `/bench-final-check` only for its terminal
  retained-evidence report and implementation retro. It does not run another
  gate, commit the spec, or author `Status: implemented`.
- Light-path work still lands its ticket with path-scoped
  `bench commit -m "<msg>"`; ordinary `bench commit` retains its gate-then-commit
  contract and never creates provisional evidence.
- Before the composed feature's final gate, emit the coverage table for every
  acceptance row —
  `bench coverage <spec>` produces it and `bench coverage --check <spec>` validates
  the map; don't hand-assemble it. Classify each row `green`, `already covered`,
  or `not TDD-able`. If any mapped behavior is missing, partial, or unclassified,
  the build is not ready for the gate.
- Once every ticket is integrated and released, run
  `/bench-review-implementation` over the exact candidate — the semantic
  three-axis pass (Standards + Spec + Coverage) that catches what the gate can't:
  right thing built the wrong way, wrong thing built cleanly, or breaking inputs
  nothing exercises. Convert its result to the review receipt the public
  lifecycle validates.
- **The terminal repair-pass bound.** A repair pass integrates the findings
  accepted for its round, performs one fresh composed review, and stops at the
  next promotion result; another semantic review round opens only when the
  composition changes or the reviewer requests one.
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
- Submit the delegate's bounded receipt through `bench spec build review` only
  while its candidate subject still matches the exact reviewed composition.
- **Promote inline** through `bench spec build promote`; its exact prospective
  subject is the final composed gate boundary. Then run **final-check inline**
  only to report the retained terminal evidence and capture the retro. Use
  **debug inline** with `/bench-debug` when an issue needs deep analysis.

**Finding disposition.** Concrete defects — bugs, spec misses, missing
coverage — return through repair assignment, checkpoint, and integration without
stopping; contestable design and
judgment findings are flagged in the exit report for reviewer veto, not
applied. Re-gating follows the repair-ticket cadence: focused evidence,
checkpoint, integration, fresh composed review, then promote. The repair pass is
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
before promotion — the Codex CLI at the top binding, charged to refute the
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
that proves it: the review delegate's invocation, lifecycle status and retained
evidence, and the `capture/session-handoff.md` boundary rewrite.
