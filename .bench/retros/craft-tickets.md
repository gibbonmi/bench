## Outcome

The `craft-tickets` spec landed as 14 green implementation commits after its
staging commit. The final ticket committed 13 paths as `4fae717` on
`bench/assign/cd165a0b561885248bf91321f402e874/199e9f9018300b9cc9d5740df5b751fa`,
including the spec transition to `Status: implemented`. The topic branch is
awaiting reviewer merge into `main`; ship-tier verification has not run.

## Gate-stage timings

The final ticket's explicit dev-tier gate took approximately 150 seconds
wall-clock. Its slowest reported package stages were artifact contracts
141.446s, worktree tests 59.276s, runtime contracts 52.496s, gate tests 50.030s,
surface contracts 40.212s, conformance suite 27.664s, publication contracts
20.034s, race tests 8.125s, and prep-release surface contracts 6.776s. Build,
gofmt, vet, test, race, conformance-suite, conformance, contract, shellcheck,
and canary all finished green; the canary reported every fixture bit. The run
reported 270 environment capability skips and no capability skips. The landing
command reused that fresh exact-tree verdict, so it did not run the gate again.

## Ticket-versus-spec-slice and delegate performance

The spec produced 14 independently green implementation commits, and the final
ticket stayed isolated in one assignment worktree. This session had no retained
delegate timing or token evidence and no comparable spec-slice delegate charge,
so a delegate-performance comparison cannot be substantiated from the available
record. The visible operational difference was that ticket isolation preserved
the final diff, but the unfinished assignment was easy to overlook from the
primary checkout.

## Coordinator catches

The coordinator initially gated `main` instead of the active assignment
worktree. `bench worktree list` exposed the mismatch: the assignment was at the
same committed HEAD but held the final ticket as uncommitted edits. After the
correct worktree gate went green, the first `bench commit` attempt refused
because two nested mutated-command fixture files were outside the named path
set. Inspecting the complete fixture trees and naming those files made the
atomic commit succeed without broad staging.

## Agent-experience improvements

### Bench CLI

Make final-check subject selection expose an active dirty assignment before a
primary-checkout gate starts. The current aggregate dashboard signals dirty
work but does not identify the authoritative final-check tree early enough; a
subject-oriented diagnostic would prevent a full gate run against the wrong
checkout.

### Skills

Add an entry check to final-check that resolves active assignment worktrees and
states the exact oracle subject before invoking the gate. The missing check
cost one redundant full gate and made the first green verdict irrelevant to the
ticket being landed.

### Process

Require the implementation handoff to name any still-open assignment and
whether its branch is committed. This keeps the phase boundary explicit when
the primary checkout is clean for implementation code but the final ticket
still exists only in an isolated worktree.
