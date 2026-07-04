---
name: craft-tdd
description: How to apply test-driven development without the cost blowup or the over-fit-and-stop failure. Use whenever writing tests first, someone says "red-green-refactor" or "write the test first", or building at a seam that /bench-write-spec marked for TDD. Reach for this before starting any TDD pass to bound where and how it applies.
index: writing tests first
---

# TDD at seams

TDD is powerful where the spec is known and the seam is right, and a liability
where the agent invents its own target and over-fits to it. This skill keeps it
on the productive side.

## Apply it only at pre-agreed seams

TDD every line and you pay a large cost tax for ceremony and watch the loop do the
minimum to pass its own tests, then stop. So:

- TDD **only** at the seams `/bench-write-spec` named. At those seams the test target is
  external — I chose the seam and the behavior — so passing the test means
  matching my spec, not the agent's guess.
- Off the marked seams, write the code and let the gate catch regressions —
  remembering the gate only catches what some test observes. If off-seam code
  carries behavior no seam can observe, that is a seam-set defect: surface it,
  don't skip it.
- The spec's user stories are the breadth **floor, not the ceiling**. At a marked
  seam, enumerate the behavior's failure modes by walking the edge classes —
  the canonical list lives in `/bench-write-spec`'s edge-inventory step — and
  propose them as coverage rows for the
  reviewer to veto; never silently skip an edge the spec forgot. The over-fit guard constrains what *correct* means (the reviewer chose
  the seam and the semantics), never which inputs get exercised. If a story isn't
  covered, that's a gap to fix, not a stop.

**Catch yourself inventing a test target mid-loop? Stop — that is the exact
failure this skill exists to prevent.** The seam and the semantics come from
the spec, never from the loop.

## The cycle

1. **Red** — write one failing test at the seam, asserting observable behavior
   through the interface. One test, not a test file: batch-writing every case
   upfront and implementing against the pile skips the loop that makes TDD work.
   One logical assertion; test what the caller cares about,
   not how it's done. Never mock an internal collaborator — if a refactor that
   keeps behavior breaks the test, the test was bound to an internal, so move it
   out to a real seam. Run it; confirm it fails for the right reason — an
   assertion failure, not a compile error. A compile-error red only proves the
   symbol is missing; it says nothing about whether the assertion can catch a
   wrong implementation. The
   expected value comes from the spec or an independent computation — never from
   running the implementation and pasting back what it returned. Vacuity check:
   an assertion that would also pass against a no-op implementation asserts
   nothing.
2. **Green** — the smallest change that passes it. "Minimal to pass" is a local
   rule for this step, not a license to stop short of the spec's breadth.
3. **Refactor** — clean up with the test green. A good test survives this; if your
   refactor breaks it, the test was attached to an internal — re-place the seam.
   Look for: duplication (extract), long methods (private helpers, tests stay on
   the interface), shallow modules (combine or deepen — see `craft-seams`), feature envy
   (move logic to where the data lives), primitive obsession (introduce value
   objects), and existing code the new code just revealed as a problem.

Inside a `bench shift` iteration, complete red→green within the same iteration: a
red gate rolls the worktree back and deletes uncommitted work, so a test left
failing when the iteration ends is destroyed, not carried to the next one.

## Acceptance rows

When `/bench-write-spec` includes an acceptance coverage map, treat each
acceptance row as the unit of TDD coverage. A valid row has `story`, `behavior`,
`seam`, `red signal`, and `why it catches the failure`.

- `behavior` is caller-visible. It is not a data shape, private method, or
  implementation step.
- `seam` is where callers cross the interface and where the test attaches.
- `red signal` is a command or test run immediately before implementing that
  row's slice, failing because the mapped behavior is absent or wrong. Rows go
  red one at a time as each slice starts — never batched into an upfront all-red
  test file. If the signal already passes, classify the row as `already covered`;
  if it cannot run before implementation, classify it as `not TDD-able` with the
  reason.
- `why it catches the failure` explains why this signal would fail when the mapped
  user-visible behavior is broken.

Reject rows and tests that attach below the chosen seam, test private behavior,
mock an internal collaborator, use an internal test double to satisfy the test
without the behavior, assert call count or call order, peek around the interface,
or describe implementation shape instead of observable behavior. Reject
nondeterminism too: a test whose verdict depends on wall-clock time, iteration
order, or a live network is a broken oracle. Mocks are fine at
real system boundaries — time, randomness, network, filesystem, external APIs —
when the seam requires them. But a boundary stub scripted to return exactly the
success shape hollows the test out — the test passes while the real integration
is never exercised. Keep stubs honest: realistic shapes, and failure behaviors
drawn from the edge inventory, not only the happy reply.

```
repo = IssueRepo(fake_http); repo.close(41); assert repo.get(41).status == "closed"
```
Good — the fake stands at a real system boundary (the injected network client)
and the assertion reads observable state back through the interface.

```
repo.close(41); mock_cache.invalidate.assert_called_once()
```
Bad — the mock pins an internal collaborator and asserts call count, so a
refactor that keeps behavior kills the test.

## The oracle is the gate, not you

A green test you wrote is not proof of done. Done is `bench gate` green across the
whole spec. Never edit, relax, or delete a test to reach green — if a test is
wrong, stop and say so. The agent's own assertions are never the completion
signal; the gate is.

## Declare the line

A TDD pass is a multi-cycle stage — declare model/effort/cap first. High effort
for the uncertain seam; low for mechanical ones. Don't grind past the cap.
