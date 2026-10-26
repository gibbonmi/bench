# Stage each verification run when it ends

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md
Covers: none

## What to build

`/bench-implement-spec` tells a ticket author to commit the ticket on a green lane, and then to run each verification requirement once at that commit. The author stages each run with `bench record excerpt` before its next Bench call, because that verb stages only the last retained response. Thus a chunk freeze or a record step runs no requirement a second time.

The gate-reuse GR-C1 freeze ran each requirement again, because the author did not stage the runs that it made at its ticket commit. One run of the `./internal/gate` package takes about 70 seconds.

## Acceptance

- [ ] The Build and Land sections of the command state the rule, and the command stays inside its line budget. The plan's probe stays a verification requirement, so the author still runs it once at the commit.
- [ ] `bench gate` passes at the landing.
