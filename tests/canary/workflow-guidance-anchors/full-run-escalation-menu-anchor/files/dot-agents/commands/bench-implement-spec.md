---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /bench-write-spec, or for a change the lighter-path threshold in .bench/BENCH.md lets you build without a spec. Declares its line, uses TDD at seams, ends on a green gate. An opt-in `--full <spec>` run carries the spec on to push-ready — implement, fresh-delegate review, final check — resuming from session-handoff.md.
---

# /bench-implement-spec — do the work at the seams

## Entry orientation

This is the implementation phase. It starts from an approved spec or a tiny change
with an obvious seam, declares the line, works vertical slices at the agreed seams,
and uses the acceptance coverage map to keep the build target fixed.

## Exit handoff

Close by reporting the implemented stories, the acceptance coverage status for
each row, the gate result, and any semantic review findings already fixed. The
recommended next command is `/bench-review-implementation` when semantic review has
not run yet; otherwise it is `/bench-final-check`. A build that stops short exits
through "When the build stops short" below instead, and recommends its route.

Implement the spec at the seams it names. If there's no spec, the change must fall
under the lighter-path threshold `.bench/BENCH.md`'s "Right-size the process"
paragraph sets; if it doesn't, stop and run `/bench-write-spec` first.

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
  returned done-claim through the `craft-delegate` contract.
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
- Run typecheck and the relevant single test file frequently as you go. Run the
  full gate once at the end.
- During structure housekeeping, apply `craft-seams`' split-or-grant rule.
- One small change at a time, repo stays green — invariant 4 in `.bench/BENCH.md`.
- Every delegation during the build carries its own line and, when the spec has
  a coverage map, its stories' coverage rows in the charge — every time; the
  rules are `craft-delegate`'s (model half: `craft-line`). Every write-delegate's
  charge also names `craft-comments` by path
  (`.agents/skills/bench-craft-comments/SKILL.md`) — comment discipline rides
  the charge, not the delegate's priors.
- Verify each returned done-claim against its charged rows and report the round
  in one line: accepted, or the missed case and the re-charged repair.
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
2. **Keep what's real:** committed green work stays committed; uncommitted work
   is described and left in the worktree. Nothing gets squash-finished to fake
   completion.
3. **Route by cause, and recommend one:**
   - wrong tier (the model ground, the gate disagreed) → re-declare one tier up
     per the `craft-line` ladder and resume this phase;
   - wrong spec (a story is unbuildable as written) → back to `/bench-write-spec`
     with the finding quoted;
   - wrong scope (the spec is bigger than one build) → propose the split; the
     reviewer decides.
4. **Leave the worktree owned:** the coordinator that cut it owns its retirement.
   A route that drops the work releases the worktree now (`bench worktree
   release`, the creating request's id); a route that carries uncommitted work
   forward names the worktree path in the report, and the session resuming that
   route retires it with path-addressed `bench worktree clean` once the work
   lands.

## Close on green

- The build is done when `bench gate` is green, and only then — invariant 1 in
  `.bench/BENCH.md` (the gate is the oracle).
- A green gate proves what the tests observe. Before handing back, drive the
  changed path once end-to-end — invoke the real command, endpoint, or call the
  diff changes and read its output. A mismatch here is a defect to fix or
  surface, never a footnote.
- Land each green build step with path-scoped `bench commit -m "<msg>"`, naming
  the files that step touched. It gates and commits the named paths atomically,
  and it enforces the commit discipline so you don't have to restate it. When it
  refuses over an unexplained working-tree file, surface the file to the
  reviewer; don't commit or revert it on your own.
- This phase ends at its last green build commit and hands off: the landing
  commit and the `Status: implemented` transition are `/bench-final-check`'s.
  Never run `bench commit --spec` and never edit a spec's status line from this
  phase.
- Before the final gate, emit the coverage table for every acceptance row —
  `bench coverage <spec>` produces it and `bench coverage --check <spec>` validates
  the map; don't hand-assemble it. Classify each row `green`, `already covered`,
  or `not TDD-able`. If any mapped behavior is missing, partial, or unclassified,
  the build is not ready for the gate.
- Once the gate is green, run `/bench-review-implementation` — the semantic three-axis
  pass (Standards + Spec + Coverage) that catches what the gate can't: right thing
  built the wrong way, wrong thing built cleanly, or breaking inputs nothing
  exercises. Read its findings, fix what matters, re-run the gate.
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

- **Implement inline**, under everything above this section.
- **Review in a delegate:** spawn one fresh-context delegate charged with the
  standalone `/bench-review-implementation` contract and given the spec and
  the diff and nothing else. Inline self-review is closed, not deprioritized —
  the context that produced the code carries the assumptions that produced its
  bugs — and a context-inheriting delegate is the same failure wearing a
  delegate's name. Verify the delegate's done-claim against the gate and
  `git status` before acting on it — invariant 1 in `.bench/BENCH.md` and
  `craft-delegate`'s verification discipline, not restated here. When the
  review delegate cannot run or returns nothing, the run stops and reports at
  that boundary rather than proceeding to final-check with review unrun.
- **Final-check inline** per `/bench-final-check`, and **debug inline** with
  `/bench-debug` when an issue needs deep analysis.

**Finding disposition.** Concrete defects — bugs, spec misses, missing
coverage — are fixed and re-gated without stopping; contestable design and
judgment findings are flagged in the exit report for reviewer veto, not
applied. The repair pass is bounded by `/bench-review-implementation`'s
terminal repair-pass rule and routed through `craft-delegate`'s repair
allowance; this mode adds no second version of either.

**The scope fence.** The run implements the spec's stories and nothing else:
work noticed outside them — an adjacent refactor, an unrelated improvement, a
story the spec chose not to take — is recorded for the reviewer rather than
built. This is the counterweight to the fix-don't-park route, whose
park-versus-fix test is `.bench/BENCH.md`'s Workflow section's: a defect in
the run's own path is fixed in-run; everything wider is recorded.

**Phase boundaries persist.** At every phase boundary the run writes the phase
reached into `session-handoff.md`'s State section, then refreshes the pin
block with `bench handoff --next <command>`. The phase reached is the one fact
this mode adds; every other field is the existing handoff contract's. A
re-invoked `--full` resumes from the phase the handoff names instead of
re-implementing from the top, and where the handoff and the tree disagree the
tree wins.

**Escalation asks first.** Those three are the fixed menu — a route is
omitted only when this harness cannot invoke it at all, and the omission is
stated rather than silent. Never escalate without asking. A harness with no
structured-prompt surface asks the same question as a plain numbered list.
The Codex route is `codex exec`: set the model, the reasoning effort, the
working directory, and a non-interactive approval posture — `codex exec
--help` is authoritative for the flag spellings, so none are pinned here.

**The falsification offer.** On that same size trigger, ask a second, separate
question: whether to add a cross-harness falsification pass over the diff
before final-check — the Codex CLI at the top binding, charged to refute the
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
that proves it: the review delegate's invocation, the commit shas the phase
landed, the `session-handoff.md` boundary rewrite.
