---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /bench-write-spec, or for a change the lighter-path threshold in .bench/BENCH.md lets you build without a spec. Declares its line and uses TDD at seams. An opt-in `--full <spec>` run carries the spec through implementation, review, the final green landing, and terminal reporting — resuming from capture/session-handoff.md.
---

# /bench-implement-spec — do the work at the seams

## Entry

A spec-backed run opens with `bench preflight build <slug>` before the first
edit; a red preflight stops the phase. If there's no spec, the change must
fall under the lighter-path threshold in `.bench/BENCH.md`'s "Right-size the
process" paragraph, or route to `/bench-write-spec` first.

## Declare the line, derive the tickets, route the venue

Declare the line before touching code — `craft-line` owns the template, the
tier decision, and the escalation ladder. Charge `craft-tickets` to break the
spec's stories and seams (or the one light-path change) into tickets under
`specs/<slug>/tickets/`, then run its reviewer-approved breakdown before any
ticket is assigned — the AFK batch-approval carve-out in `.bench/BENCH.md` is
the only no-round-trip route. Every spec-backed run assigns genuine
write work to at least one write subagent before the first implementation
edit — `craft-delegate` owns isolation and verification, `craft-line` owns
its model and effort. A read-only helper never satisfies the write
requirement; when the harness cannot spawn a write subagent,
`craft-delegate`'s capability-aware policy governs what runs inline instead.

## Build

Work the user stories in vertical slices at the agreed seams, TDD only where
`craft-tdd` marks the seam. When the spec carries a coverage map, seed the
harness's native task list from `bench coverage <spec>`, one task per row.
Each write delegate runs its ticket's focused checks and its own mutation
probe; the coordinator independently verifies the returned diff before
landing it.

## Land

Tickets land serially in `Blocked by:` order, each through path-scoped
`bench commit`, gate-then-commit atomically. Once every ticket short of the
final landing is in, run `/bench-review-implementation` over the composed
diff; accepted findings land as repair tickets through the same cadence. The
final green landing commit is `bench commit --spec <slug>` — it is the sole
author of the spec's `Status: implemented` flip. Then run `/bench-final-check`.

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
