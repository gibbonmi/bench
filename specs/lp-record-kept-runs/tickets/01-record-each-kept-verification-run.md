# Record each kept verification run

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md
Covers: none

## What to build

`/bench-implement-spec` tells a ticket author to run each verification requirement of its ticket once, at the ticket commit, and to keep the output file of that run. The author passes that file to `bench record verification --excerpt`. Thus a chunk freeze or a record step runs no requirement a second time.

The gate-reuse GR-C1 freeze ran the `./internal/gate` package again, because the author did not keep the runs that it made at its ticket commit. One run of that package takes about 70 seconds.

## Acceptance

- [ ] The Build and Land sections of the command state the rule, and the command stays inside its line budget. The plan's probe stays a verification requirement, so the author still runs it once at the commit.
- [ ] `bench gate` passes at the landing.
