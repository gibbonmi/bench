---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /bench-write-spec, or for a change the lighter-path threshold in .bench/BENCH.md lets you build without a spec. Declares its line and uses TDD at seams. An opt-in `--full <spec>` run carries the spec through implementation, review, the final green landing, and terminal reporting — resuming from capture/session-handoff.md.
---

# /bench-implement-spec — do the work at the seams

## Entry

If there's no spec, the change must fall under the lighter-path threshold in
`.bench/BENCH.md`'s "Right-size the process" paragraph, or route to
`/bench-write-spec` first.

If a spec-backed run has no `specs/<slug>/tickets/` directory or that directory contains no
ticket files, return to `/bench-write-spec`; ticket slicing and approval belong there.

## Declare the line, validate the tickets, route the venue

Declare the line before touching code — `craft-line` owns the template, the
tier decision, and the escalation ladder. After the write-spec-phase approval,
create or retain one integration worktree at the reviewed graph commit; that
commit is its frozen review base. Enter it and run `bench preflight build <slug>`
before the first edit; a red preflight stops the phase.
Every spec-backed run assigns genuine write work to a write subagent before
implementation: `craft-delegate`
owns isolation and verification, `craft-line` its model and effort. A read-only
helper does not satisfy this; `craft-delegate` governs incapable harnesses.

## Build

Work the user stories in vertical slices at the agreed seams, TDD only where
`craft-tdd` marks the seam. When the spec carries a coverage map, seed the
harness's native task list from `bench coverage <spec>`, one task per row.
Each write delegate runs its ticket's focused checks and its own mutation
probe; the coordinator independently verifies the returned diff before
landing it.

## Land

Tickets commit green serially in `Blocked by:` order on the retained integration
source. Then run `/bench-review-implementation` over its explicit frozen base
and source tip. Accepted findings commit there as repair tickets and move the
tip, so review runs again. From the clean destination, hand the accepted pair to
`bench worktree land`; it composes, gates, publishes `Status: implemented`, and
releases the source. Then run `/bench-final-check`.

## When the build stops short

Report state (done vs. remaining, the coverage table, what consumed the
cap), keep landed work landed, and route by cause: wrong tier → escalate one
`craft-line` rung and resume; wrong spec → `/bench-write-spec` with the
finding quoted; wrong scope → propose the split for the reviewer to decide.

## `--full <spec>`

Orchestrates implementation, a delegate-run review, the final landing, and
`/bench-final-check`; refuses on a missing or unnamed spec. At every phase
boundary it writes the phase reached into `capture/session-handoff.md` and
refreshes the pin with `bench handoff`. A diff large enough to risk missed
bugs pauses to ask the reviewer before escalating tier or adding a
cross-harness falsification pass — both are offers, never silently applied.
