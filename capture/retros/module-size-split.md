# Retro: module-size-split

## Outcome

The spec landed as `Status: implemented` at commit `74ef1125`. The final phase amended row R13's scope (the reviewer widened it by one file, 2026-08-24) and split `internal/lines/lines_test.go` into four files as ticket 13. `bench structure` fell from 71 issues at the spec's base to 55 at the landing, and R13's bound of at most 55 holds. All 20 coverage rows passed, and the three review axes returned zero findings on the finale.

## Gate-stage timings

The landing gate at `74ef1125`: gofmt 0.1 s, vet 1.3 s, test 65.7 s, race 4.6 s, system 10.1 s, shellcheck 0.4 s. The test phase dominates; `internal/worktree` alone takes 56 s of it.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized; no delegate received a spec slice. The finale's one Opus/low charge (44k tokens) landed first-pass with a body-line multiset proof and three attributed mutation probes of its own. This extends the batch-2 result: eleven Opus/low charges, eleven first-pass diffs. The spec amendment was authored inline by the orchestrator, with one Sonnet/medium review round.

## Coordinator catches

The coordinator caught nothing in the delegate's finale work; every claim verified. The material catch of the phase was pre-build: ticket 13's premise was stale (its target already held 209 lines). That shortfall stopped batch 2 and forced the scope amendment. The orchestrator's own misses: three prose-bound red gates on capture prose, one unrecoverable worktree request token, and a primary-checkout retire that needed a stash.

## Repair attribution

| ticket | repair rounds | causes |
| --- | --- | --- |
| 01 | 0 | none |
| 02 | 0 | none |
| 03 | 0 | none |
| 04 | 0 | none |
| 05 | 0 | none |
| 06 | 0 | none |
| 07 | 0 | none |
| 08 | 0 | none |
| 09 | 0 | none |
| 10 | 1 | spec-row |
| 11 | 0 | none |
| 12 | 0 | none |
| 13 | 1 | spec-row |

Ticket 10's round: a census check pins `owner_test.go` by path, and the spec never named it. Ticket 13's round is the stale-premise amendment: the spec's line count for the target predated the tree.

## Agent-experience improvements

### Bench CLI

- Make `bench preflight build` compare each ticket's stated line count against the tree, so a stale premise reds before a charge spends tokens.
  Feeds: new
- Make `bench worktree create` echo the request token back in its output, so a landing never depends on the caller's own recording.
  Feeds: new
- Make `bench spec retire` refuse the primary checkout the same way `bench commit` does, so the retire never dirties a checkout that cannot commit.
  Feeds: new

### Skills

- Add the prose-bound check to the capture-writing step in the final-check and handoff guidance. Three landings in a row paid red gates to over-bound capture sentences.
  Feeds: new

### Process

- Record the worktree request token in a scratch file in the same command that mints it, before the create call runs.
  Feeds: none
