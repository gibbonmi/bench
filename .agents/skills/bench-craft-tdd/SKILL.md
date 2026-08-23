---
name: craft-tdd
description: How to apply test-driven development without the cost blowup or the over-fit-and-stop failure. Use whenever writing tests first, someone says "red-green-refactor" or "write the test first", or building at a seam that /bench-write-spec marked for TDD.
index: writing tests first
---

# TDD at seams

TDD is powerful where the spec is known and the seam is right. It is a liability
where the agent invents its own target and over-fits to it. This skill keeps it
on the productive side.

## Apply it only at pre-agreed seams

TDD every line, and you pay a cost tax for ceremony while the loop does the
minimum to pass its own tests, then stops. So:

- TDD **only** at reviewer-confirmed seams. Spec-backed work consumes the seams
  `/bench-write-spec` named: spec sign-off already confirmed them, so it takes the
  signed-off seam without a second reviewer gate. At those seams the test target is
  external — I chose it — so a passing test matches my spec, not the agent's
  guess.
- Light-path work names the test seam in its ticket file and starts without a live
  confirmation stop; the reviewer can veto the seam post-hoc. The right-size
  table's standing approval gives the named seam external-target authority.
- Off the marked seams, write the code and let the gate catch regressions —
  remembering the gate only catches what some test observes. If off-seam code carries
  behavior no seam can observe, that is a seam-set defect: surface it, don't skip it.
- The spec's user stories are the breadth **floor, not the ceiling**. At a
  marked seam, enumerate the behavior's failure modes by walking the classes
  below, and propose them as coverage rows for the reviewer to veto. Never
  silently skip an edge the spec forgot. The over-fit guard constrains what
  *correct* means (the reviewer chose the seam and the semantics), never which
  inputs get exercised. If a story isn't covered, that's a gap to fix, not a stop.

Walk these classes: error path, empty/absent input, boundary values, malformed
input, interrupted or partial state, re-run idempotency, process-boundary
lifecycle, hostile environment, plus the project profile's hostile-input
checklist. A control resolving a class must exercise the **new** surface.

**Catch yourself inventing a test target mid-loop? Stop — that is the exact
failure this skill exists to prevent.** The seam and the semantics come from
the spec, never from the loop.

## The cycle

1. **Red** — write one failing test at the seam, asserting observable behavior through
   the interface. Write one test, not a test file: batch-writing every case upfront and
   implementing against the pile skips the loop that makes TDD work. Keep one logical
   assertion, and test what the caller cares about, not how it is done.
   - Never mock an internal collaborator. If a refactor that keeps behavior breaks the
     test, the test was bound to an internal, so move it out to a real seam.
   - Run the test, and confirm it fails for the right reason: an assertion failure,
     not a compile error. A compile-error red only proves the symbol is missing; it
     says nothing about whether the assertion can catch a wrong implementation.
   - In a compiled language, stub the minimal declarations with no behavior so the one
     test compiles, then confirm the behavioral red. Missing symbols are pressure to
     scaffold, never license to batch the file.
   - The expected value comes from the spec or an independent computation, written as a
     literal. Never run the implementation and paste back what it returned. Never recompute
     it in the test with the implementation's own algorithm: a test that mirrors the code
     passes by construction.
   - Vacuity check: an assertion that would also pass against a no-op implementation asserts
     nothing. `references/tests.md` shows what a good test looks like: the four properties,
     with a worked good/bad pair.
2. **Green** — the smallest change that passes it. "Minimal to pass" is a local
   rule for this step, not a license to stop short of the spec's breadth.
3. **Refactor** — clean up with the test green. A good test survives this; if your
   refactor breaks it, the test was attached to an internal — re-place the seam.
   Hunt `craft-review`'s smell baseline plus shallow modules (combine or
   deepen — see `craft-seams`), and existing code the new code just revealed
   as a problem. Tests stay on the interface while helpers move beneath it.

Inside a `bench shift` iteration, complete red→green within the same iteration.
A red gate rolls the worktree back and deletes uncommitted work. The rollback destroys a test
left failing when the iteration ends; it does not carry to the next one.

## Acceptance rows

When `/bench-write-spec` includes an acceptance coverage map, treat each
acceptance row as the unit of TDD coverage. The row schema is `bench-craft-spec`'s
and the red-signal classification is this skill's. Don't restate the map's
fields here, and classify a row only from the run that observes it:

- Rows go red one at a time as each slice starts — never batched into an
  upfront all-red test file. The row's red signal runs immediately before
  that slice's implementation and fails because the mapped behavior is absent or
  wrong.
- If the signal already passes, classify the row as `already covered`; if it
  cannot run before implementation, classify it as `not TDD-able` with the
  reason.

Reject a row or test that:
- attaches below the chosen seam
- tests private behavior
- mocks an internal collaborator
- uses an internal double to satisfy the test without the behavior
- asserts call count or call order
- peeks around the interface
- describes implementation shape instead of observable behavior
- shows nondeterminism: a test whose verdict depends on wall-clock time, iteration order, or a live network is a broken oracle

Mocks are fine at real system boundaries — time, randomness, network,
filesystem, external APIs — when the seam requires them. But a boundary stub
scripted to return exactly the success shape hollows the test out — the test
passes while nothing exercises the real integration. Keep stubs honest:
realistic shapes, and failure behaviors drawn from the edge inventory, not
only the happy reply. The mock-or-not rules and the honest-stub good/bad pair
are in `references/mocking.md`.

## The oracle is the gate, not you

A green test you wrote is not proof of done — and neither is a runner's summary line. A
run whose only matched test skipped still prints `ok` (Go does), so read the test
output, not the summary. Done is `bench gate` green across the whole spec. Never edit,
relax, or delete a test to reach green — if a test is wrong, stop and say so. The
agent's own assertions are never the completion signal; the gate is.

## Declare the line

A TDD pass is a multi-cycle stage — declare model/effort/cap first. High effort
for the uncertain seam; low for mechanical ones. Don't grind past the cap.
