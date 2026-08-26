---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /bench-write-spec, or for a change the lighter-path threshold in .bench/BENCH.md lets you build without a spec. Declares its line and uses TDD at seams. An opt-in `--full <spec>` run carries the spec through implementation, review, the final green landing, and terminal reporting — resuming from capture/session-handoff.md.
---

# /bench-implement-spec — do the work at the seams

## Entry orientation

This is the implementation phase. It starts from an approved spec, or from a change
small enough for the lighter-path threshold. It declares the line and works vertical
slices at the pre-agreed seams. If there is no spec, the change must fall under that
threshold in `.bench/BENCH.md`'s "Right-size the process" paragraph. If it does not
fall under that threshold, route to `/bench-write-spec` first.

If a spec-backed run has no `specs/<slug>/tickets/` directory or that directory contains no
ticket files, return to `/bench-write-spec`; ticket slicing and approval belong there.

## Exit handoff

Close by reporting the implemented stories, each acceptance row's coverage status, and
the landed commits. A reviewed spec-backed build closes through `/bench-review-implementation`
and `bench worktree land`, then `/bench-final-check`. A light-path build goes straight
to `/bench-final-check`. A build that stops short exits through "When the build stops
short" below, which recommends its one durable next action.

## Declare the line, validate the tickets, route the venue

Declare the line before you touch code. `craft-line` owns the template, the tier decision,
and the escalation ladder. After the write-spec-phase approval,
create or retain one integration worktree at the reviewed graph commit; that
commit is its frozen review base. Enter it and run `bench preflight build <slug>`
before the first edit; a red preflight stops the phase.

Every spec-backed run assigns genuine write work to a write subagent before
implementation; `craft-delegate` owns isolation and verification, and `craft-line`
owns its model and effort. A read-only helper does not satisfy this requirement.
`craft-delegate` governs incapable harnesses.

## Build

Work the user stories in vertical slices at the agreed seams. Use TDD only where
`craft-tdd` marks the seam. When the spec carries a coverage map, seed the
harness's native task list from `bench coverage <spec>`, whose
`rows[N]{story,behavior,seam}` projection is one task per row. Charge the whole frontier in parallel unless `craft-tickets` names a conflict. Each write
delegate runs its ticket's focused checks and its own mutation probe. The
coordinator independently verifies the returned diff before it lands.

## Land

Tickets commit green serially in `Blocked by:` order on the retained integration
source. Then run `/bench-review-implementation` over its explicit frozen base and
source tip. Accepted findings commit there as repair tickets and move the tip.
Then run `/bench-review-implementation` in repair-scoped mode with the accepted
repair predicates and the prior reviewed tip; that command owns the scope and
convergence rule. From the clean destination, hand the accepted pair to
`bench worktree land`; it composes, gates, publishes `Status: implemented`, and
releases the source. Then run `/bench-final-check`.

## When the build stops short

Report the state: what is done, what remains, the coverage table, and what
consumed the cap. Keep landed work landed. Then route by cause:

- Wrong tier: escalate one `craft-line` rung and resume.
- Wrong spec: route to `/bench-write-spec` with the finding quoted.
- Wrong scope: propose the split for the reviewer to decide.

## `--full <spec>`

This command orchestrates implementation, a delegate-run review, the final landing,
and `/bench-final-check`. It refuses on a missing or unnamed spec. At every phase
boundary it writes the phase reached into `capture/session-handoff.md` and refreshes
the pin with `bench handoff`. A diff large enough to risk missed bugs pauses to ask
the reviewer before escalating tier or adding a cross-harness falsification pass.
Both are offers; the command never applies them silently.
