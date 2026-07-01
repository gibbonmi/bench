# TDD acceptance coverage — make feature builds prove the promised behavior

> **BOOTSTRAPPED (2026-06-30).** Comparison against Matt Pocock's engineering skills
> found that Bench already has the bug-fix discipline in `/bench-debug`; the gap is
> the feature-build path. It says TDD happens at pre-agreed seams, but it does not
> force a behavior-to-test inventory that proves every promised story has a
> red-capable signal before implementation proceeds.

## Grounding

- Bench's TDD skill already has the core loop: red, green, refactor; public seam
  tests; no internal mocks; gate as oracle.
- `/bench-write-spec` already decides seams and asks for testing decisions, but its
  template does not require a per-story acceptance-test map.
- `/bench-implement-spec` already says vertical slices, not all-tests-first, but it
  does not require the agent to close each story against a named test signal.
- `/bench-debug` already carries the hard-bug loop: tight repro, red-capable signal,
  minimise, regression test, and "no correct seam is a finding." Do not duplicate
  that whole bug workflow into feature TDD.
- The external pattern worth folding is narrow: list observable behaviors before
  code, test one behavior at a time, and require each test to be red-capable for
  the behavior it claims to protect.

## #1: Which artifact owns the behavior-to-test map?

Type: Grill

### Question
Should the acceptance coverage map live in the spec's Testing decisions, in
`craft-tdd`, in `/bench-implement-spec`, or across all three? The layer matters: if
the spec does not own it, the agent invents tests mid-loop; if only the spec owns it,
the implementation phase may not mechanically close every behavior.

### Answer
Resolved (grill, 2026-06-30). The behavior-to-test map is owned across all
three artifacts, with the spec as the source of truth:

- `/bench-write-spec` writes the acceptance coverage map in Testing decisions, so
  the target is external before implementation starts.
- `craft-tdd` defines what makes a valid mapped row and red test: observable
  behavior at the chosen seam, not internals or agent-invented coverage.
- `/bench-implement-spec` must close each mapped behavior while building, so a
  green local test is not enough if mapped story breadth remains uncovered.

This keeps correctness external to the agent while making implementation
mechanically accountable for every mapped behavior.

## #2: What is the smallest useful acceptance coverage row?

Blocked by: #1
Type: Grill

### Question
Define the minimum row shape. Candidate fields: user story, observable behavior,
seam, red-capable signal, and why the signal catches the failure. Too little leaves
room for vague coverage; too much turns every spec into ceremony.

### Answer
— (open)

## #3: What does "red-capable" mean for new feature behavior?

Blocked by: #2
Type: Grill

### Question
For bug fixes, red-capable means the loop catches the user's exact symptom. For
feature work, define the equivalent: should a test prove the missing capability
fails before code, prove a documented edge case fails, or merely assert the
intended post-build behavior?

### Answer
— (open)

## #4: How does `/bench-implement-spec` prove it covered the map?

Blocked by: #1, #2, #3
Type: Grill

### Question
Should implementation require a per-slice closeout that names the row just made
green, a final coverage table before the gate, or only rely on `/bench-review-implementation`
to catch omissions? The goal is to prevent "green on a few tests, unfinished on the
story breadth" without adding a noisy checklist to every tiny change.

### Answer
— (open)

## #5: What quality bar prevents shallow or internal tests?

Blocked by: #2
Type: Research

### Question
Decide whether the existing `craft-seams` and `craft-tdd` wording is enough, or
whether the new map needs explicit rejection criteria: internal mocks, private
methods, call counts, database peeking instead of interface behavior, or tests that
would pass while the user-visible behavior is broken.

### Answer
— (open)

## #6: What proves this kit change actually works?

Blocked by: #1, #2, #3, #4, #5
Type: Research

### Question
Because this changes agent behavior rather than runtime code, define the dogfood
proof before building: a small spec that uses the coverage map, a shift that closes
each row, and a review/gate pass that would have caught at least one omitted
behavior. The proof must fit Bench's anti-sediment bar: real gap filled, no broad
process bloat.

### Answer
— (open)
