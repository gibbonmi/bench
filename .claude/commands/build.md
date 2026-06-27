---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /spec, or for a change small enough that the seams are obvious. Declares its line, uses TDD at seams, ends on a green gate.
---

# /build — do the work at the seams

Implement the spec at the seams it names. If there's no spec, the change must be
small enough that the seam is obvious; if it isn't, stop and run `/spec` first.

## Open with the line

State it in one line before touching code:

> Line: <model> / <effort> / ~<token cap>. <one clause why.>

Match the line to the work. Plumbing at a known seam → cheap model, low effort.
The genuinely uncertain seam → top model, high effort. Don't escalate silently;
if you hit the cap, stop and report.

## Then build

- Work the user stories in vertical slices, not all-tests-first horizontal ones.
- Use **TDD only at the pre-agreed seams** (see `tdd-at-seams`). Elsewhere, write
  the code and let the gate catch regressions. TDD everything is the cost trap.
- Run typecheck and the relevant single test file frequently as you go. Run the
  full gate once at the end.
- Smallest diffs that advance a story. Read before you write. Compose existing
  seams before inventing new ones.

## Close on green

- The build is done when `bench gate` is green — not before, and not because the
  diff looks right. If the gate is red, the build continues or stops with an
  explanation; it never declares done on red.
- Do not weaken a test or a check to reach green. If a check is wrong, surface it.
- Once the gate is green, run `/review` — the semantic two-axis pass (Standards +
  Spec) that catches what the gate can't: right thing built the wrong way, or wrong
  thing built cleanly. Read its findings, fix what matters, re-run the gate.
- Then summarize what changed in plain language and hand back. I own the merge;
  propose it, don't perform it.

For UI work, if your project has an interaction-layer skill and a screenshot loop,
they're part of the gate alongside the `design-system` skill — a green test suite is
necessary but not sufficient for UI.
