## Outcome

Six closed or shipped maps leave the active inventory. The integrity check
rejects a map that exists both at top level and under a spec. The drain audit
covers closed maps without a spec. The gate-budget decision survives in an ADR.

The original `bench maps` repro passed after the cleanup. The copied-map
regression failed before the fix and passed after it. A live copied-map probe
made `bench test --check decision-map-integrity` fail with the ownership diagnostic.
The same command passed after removal of the temporary copy.

## Gate-stage timings

The worktree gate passed all six phases. Gofmt took 119 ms, vet took
1250 ms, and test took 130595 ms. Race took 3728 ms, system took
50774 ms, and shellcheck took 626 ms. The run reported eight capability skips
and zero environment skips.

## Ticket-versus-spec-slice and delegate performance

This debug repair has no spec or tickets. The user requested continuation in
the same session. No delegate served in this run.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| map retirement | Six obsolete maps leave the projection | observed | unknown | unknown | gpt-6-astra / high / author |

## Coordinator catches

The reference audit found a sixth obsolete map: gate-concurrency.
Its only residual question belonged to the closed gate-budget proposal.
The retired files were copied aside before removal.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| spec-less debug repair | 1 | tree/tooling: stale map ownership and incomplete closure audit |

## Agent-experience improvements

### Bench CLI

The integrity check names both conflicting paths. Disabling its condition made
the regression fail; `bench probe` reported `bit` and restored the source.
Feeds: none

### Skills

The drain audit includes maps closed without a spec. This repair implements
the improvement; no open improvement remains.
Feeds: none

### Process

Semantic closure still requires evidence from the tree and recorded decisions.
The structural check catches copied maps; it does not infer completion from prose.
Feeds: none
