---
name: tdd-at-seams
description: How to apply test-driven development without the cost blowup or the over-fit-and-stop failure. Use whenever writing tests first, doing red-green-refactor, or building at a seam that /spec marked for TDD. Reach for this before starting any TDD pass to bound where and how it applies.
---

# TDD at seams

TDD is powerful where the spec is known and the seam is right, and a liability
where the agent invents its own target and over-fits to it. This skill keeps it
on the productive side.

## Apply it only at pre-agreed seams

TDD every line and you pay a large cost tax for ceremony and watch the loop do the
minimum to pass its own tests, then stop. So:

- TDD **only** at the seams `/spec` named. At those seams the test target is
  external — I chose the seam and the behavior — so passing the test means
  matching my spec, not the agent's guess.
- Off the marked seams, write the code and let the gate catch regressions.
- The breadth comes from the spec's user stories, not from what tests the agent
  could imagine. If a story isn't covered, that's a gap to fix, not a stop.

## The cycle

1. **Red** — write one failing test at the seam, asserting observable behavior
   through the interface. One logical assertion; test what the caller cares about,
   not how it's done. Never mock an internal collaborator — if a refactor that
   keeps behavior breaks the test, the test was bound to an internal, so move it
   out to a real seam. Run it; confirm it fails for the right reason.
2. **Green** — the smallest change that passes it. "Minimal to pass" is a local
   rule for this step, not a license to stop short of the spec's breadth.
3. **Refactor** — clean up with the test green. A good test survives this; if your
   refactor breaks it, the test was attached to an internal — re-place the seam.
   Look for: duplication (extract), long methods (private helpers, tests stay on
   the interface), shallow modules (combine or deepen — see `seams`), feature envy
   (move logic to where the data lives), primitive obsession (introduce value
   objects), and existing code the new code just revealed as a problem.

## The oracle is the gate, not you

A green test you wrote is not proof of done. Done is `bench gate` green across the
whole spec. Never edit, relax, or delete a test to reach green — if a test is
wrong, stop and say so. The agent's own assertions are never the completion
signal; the gate is.

## Declare the line

A TDD pass is a multi-cycle stage — declare model/effort/cap first. High effort
for the uncertain seam; low for mechanical ones. Don't grind past the cap.
