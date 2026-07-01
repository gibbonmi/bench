# TDD acceptance coverage — make feature builds prove the promised behavior

> **RESOLVED (2026-07-01).** All tickets (#1-#6) are decided. The build path is
> clear: `/bench-write-spec` should write the implementation spec for adding
> acceptance coverage maps to the feature-build workflow, and that spec must
> dogfood the five-field map decided here.

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
Resolved (grill, 2026-07-01). The smallest useful acceptance coverage row has
five fields:

- `story` — the user story this row protects.
- `behavior` — the observable behavior the caller/user should see.
- `seam` — where the test attaches.
- `red signal` — the command or test that must fail before the behavior exists.
- `why it catches the failure` — the short explanation that proves this signal
  would catch a real miss instead of merely exercising code.

This is the minimum shape that ties story breadth to caller-visible behavior,
places the test at a chosen seam, and forces the agent to explain why the signal is
protective rather than decorative.

## #3: What does "red-capable" mean for new feature behavior?

Blocked by: #2
Type: Grill

### Question
For bug fixes, red-capable means the loop catches the user's exact symptom. For
feature work, define the equivalent: should a test prove the missing capability
fails before code, prove a documented edge case fails, or merely assert the
intended post-build behavior?

### Answer
Resolved (grill, 2026-07-01). For feature work, a red signal means the
named command or test has already been run before implementation and fails because
the mapped behavior is absent or wrong at the chosen seam.

If the signal already passes, the row is already covered or is testing the wrong
thing; do not count it as new TDD coverage. If the signal cannot be run before
implementation, the spec must say why and mark the row as not TDD-able rather
than pretending it has a red signal.

## #4: How does `/bench-implement-spec` prove it covered the map?

Blocked by: #1, #2, #3
Type: Grill

### Question
Should implementation require a per-slice closeout that names the row just made
green, a final coverage table before the gate, or only rely on `/bench-review-implementation`
to catch omissions? The goal is to prevent "green on a few tests, unfinished on the
story breadth" without adding a noisy checklist to every tiny change.

### Answer
Resolved (grill, 2026-07-01). `/bench-implement-spec` proves coverage in
two lightweight ways:

- During the build, each vertical slice names the acceptance coverage row it is
  turning red-to-green.
- Before the final gate, the agent emits a compact coverage table showing every
  row as `green`, `already covered`, or `not TDD-able`, with the recorded reason
  for any row that could not start from a red signal.

`/bench-review-implementation` still audits omissions, but implementation cannot
outsource basic coverage accounting to review. The build phase must show that
every mapped behavior was closed or explicitly classified before it claims green.

## #5: What quality bar prevents shallow or internal tests?

Blocked by: #2
Type: Research

### Question
Decide whether the existing `craft-seams` and `craft-tdd` wording is enough, or
whether the new map needs explicit rejection criteria: internal mocks, private
methods, call counts, database peeking instead of interface behavior, or tests that
would pass while the user-visible behavior is broken.

### Answer
Resolved (research, 2026-07-01). The existing `craft-seams` and
`craft-tdd` principles are correct but not sufficient by themselves for the new
acceptance coverage map. The quality bar should be folded into the map/test
validity wording, not left implicit across multiple skills.

Valid rows and tests must be rejected when they:

- attach below the chosen seam or test private/internal behavior;
- mock internal collaborators or assert call counts/order instead of observable
  behavior;
- verify by peeking around the interface instead of through the chosen seam;
- use a signal that would still pass while the mapped user-visible behavior is
  broken;
- describe implementation shape rather than caller-visible behavior.

Mocks remain allowed at real system boundaries such as time, randomness, network,
filesystem, or external APIs when the seam requires it. The point is not "no
mocks"; it is "no internal test doubles that let an agent satisfy the test while
missing the behavior."

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
Resolved (research, 2026-07-01). This kit change is proven only when it is
dogfooded through the workflow it changes. Adoption requires three proofs:

1. **Self-use in the implementation spec.** The spec for this change must include
   an acceptance coverage map using the decided five-field row shape. Rows should
   cover the spec template, `craft-tdd` row/test validity, and
   `/bench-implement-spec` closeout behavior.
2. **Shift closeout against the map.** The implementation pass must name each row
   as it turns red-to-green and end with the compact coverage table from #4 before
   the gate. Any row that cannot start red is classified as `already covered` or
   `not TDD-able` with a reason.
3. **Real dogfood run after the edits.** With the changed kit in the working tree,
   run a small real Bench planning/build path on this repo or another linked repo:
   `/bench-write-spec` for a tiny behavior change, then a build/shift that closes
   the coverage map, then `/bench-review-implementation`, then `bench gate`.

A static wording check can support the proof, but cannot replace it. If the real
dogfood run cannot be completed, the synthesis is incomplete and the change stays
proposed rather than adopted.
