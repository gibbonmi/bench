---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /bench-write-spec, or for a change small enough that the seams are obvious. Declares its line, uses TDD at seams, ends on a green gate.
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

Implement the spec at the seams it names. If there's no spec, the change must be
small enough that the seam is obvious; if it isn't, stop and run `/bench-write-spec` first.

## Open with the line

Declare the line before touching code — the declaration template, the decision
table that picks the row, and the escalation ladder are all `craft-line`'s.

## Then build

- Work the user stories in vertical slices, not all-tests-first horizontal ones.
- Use TDD only at the pre-agreed seams; its bounds are `craft-tdd`'s.
- If the spec has an acceptance coverage map, each vertical slice names the
  coverage row it is turning red-to-green before editing that slice. Rows marked
  `already covered` or `not TDD-able` keep their recorded reason; don't silently
  upgrade them into TDD coverage.
- Run typecheck and the relevant single test file frequently as you go. Run the
  full gate once at the end.
- One small change at a time, repo stays green — invariant 4 in `.bench/BENCH.md`.
- Every delegation during the build carries its own line — the rules are
  `craft-delegate`'s (model half: `craft-line`).
- For broad renames or reference refactors, dry-run the file scope before editing,
  then verify old stems in every form: `/name`, `$name`, bare basenames in
  inventories, and `dir/name` path forms. Separator slashes inside prose are not
  command invocations.

## When the build stops short

A build that hits its token cap or ends with stories unmet exits through a
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

## Close on green

- The build is done when `bench gate` is green, and only then — invariant 1 in
  `.bench/BENCH.md` (the gate is the oracle).
- A green gate proves what the tests observe. Before handing back, drive the
  changed path once end-to-end — invoke the real command, endpoint, or call the
  diff changes and read its output. A mismatch here is a defect to fix or
  surface, never a footnote.
- When a commit happens in this phase, stage the files the build actually touched,
  explicitly — never a blind `git add -A`. An unexplained working-tree file blocks
  the commit: surface it to the reviewer; don't commit or revert it on your own.
- Advance the spec's lifecycle in the green-gate commit that finishes the build:
  flip its `Status: staged` line to `Status: implemented`, and stage that one-line
  edit with the build's files. `implemented` then honestly means *built, gate-green,
  awaiting review/merge* — the state `bench status`'s retirement signal keys on once
  the spec reaches the default branch. Never write a line-start `Status: implemented`
  into any other `specs/*.md`, or that detector fires on a spec that is not done.
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
