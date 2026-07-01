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
not run yet; otherwise it is `/bench-final-check`.

Implement the spec at the seams it names. If there's no spec, the change must be
small enough that the seam is obvious; if it isn't, stop and run `/bench-write-spec` first.

## Open with the line

State it in one line before touching code:

> Line: <model> / <effort> / ~<token cap>. <one clause why.>

Match the line to the work. Plumbing at a known seam → cheap model, low effort.
The genuinely uncertain seam → top model, high effort. Don't escalate silently;
if you hit the cap, stop and report.

## Then build

- Work the user stories in vertical slices, not all-tests-first horizontal ones.
- Use **TDD only at the pre-agreed seams** (see `craft-tdd`). Elsewhere, write
  the code and let the gate catch regressions — the gate only catches what some
  test observes, so behavior no seam can observe is a seam-set defect to surface,
  not to skip. TDD everything is the cost trap.
- If the spec has an acceptance coverage map, each vertical slice names the
  coverage row it is turning red-to-green before editing that slice. Rows marked
  `already covered` or `not TDD-able` keep their recorded reason; don't silently
  upgrade them into TDD coverage.
- Run typecheck and the relevant single test file frequently as you go. Run the
  full gate once at the end.
- Smallest diffs that advance a story. Read before you write. Compose existing
  seams before inventing new ones.
- For broad renames or reference refactors, dry-run the file scope before editing,
  then verify old stems in every form: `/name`, `$name`, bare basenames in
  inventories, and `dir/name` path forms. Separator slashes inside prose are not
  command invocations.

## Close on green

- The build is done when `bench gate` is green — not before, and not because the
  diff looks right. If the gate is red, the build continues or stops with an
  explanation; it never declares done on red.
- Do not weaken a test or a check to reach green. If a check is wrong, surface it.
- Before the final gate, emit a compact coverage table for every acceptance row:
  `green`, `already covered`, or `not TDD-able`. If any mapped behavior is missing,
  partial, or unclassified, the build is not ready for the gate.
- Once the gate is green, run `/bench-review-implementation` — the semantic two-axis pass (Standards +
  Spec) that catches what the gate can't: right thing built the wrong way, or wrong
  thing built cleanly. Read its findings, fix what matters, re-run the gate.
- Then summarize what changed in plain language and hand back. I own the merge;
  propose it, don't perform it.

For UI work, if your project has an interaction-layer skill and a screenshot loop,
they're part of the gate alongside the `craft-design-system` skill — a green test suite is
necessary but not sufficient for UI.
