## Outcome

The build landed `native-record-operations` on `main` at 576441482d5ad5ac5c9acaf73463f701bb4b77bd from source tip 90799f4a. It delivered the native `bench record assignment` form and the stale-plan guard at the chunk writer with typed repair routes. It also delivered complete probe result import with the `probe_execution` row and `--probe-result`. Three chunks closed with one repair cycle each, and all 45 acceptance rows hold citations. The build also retired the FT318 and FT168 rows of `docs/known-issues.md`. The run id is `nro-build-20261009`.

Three orchestrator calls stay open to a reviewer veto. The completion ordering note and the copied-flag sources show in the root record help rows. An existing test pins one grammar line per form help. Every biting `bench probe` answer now spills, because the `probe_execution` row crosses the 10-line bound. The test accessor `Find` answers only a first reserved frame that is `spilled`, and no test pins that on a two-frame response.

## Gate-stage timings

- landing: commit 576441482d5ad5ac5c9acaf73463f701bb4b77bd, trace 4d2cb0cb2b4959b00d3de767adf9ffb7
- gofmt: 168 ms
- vet: 1411 ms
- test: 350170 ms
- race: 4360 ms
- system: 183114 ms
- shellcheck: 587 ms

The whole-project gate ran six times in this build. Those runs were three chunk checkpoints, one repeated NR-C3 checkpoint after a late `main` fold, and the two merge gates. Each run took 6 to 10 minutes, so the gate cost about 50 minutes of the 5-hour wall time. The orchestrated part after the reviewer restart took 4 hours and 40 minutes.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| NR-C1 ticket 1 | NR4 whole-plan validation verified | verified | 9 | held | opus / high / implementer |
| NR-C1 Standards | S1 preimage loop is duplicated | auto-fix | 6 | refuted | fable / medium / reviewer |
| NR-C1 Standards | S2 ticket-path rule has two sources | ask-user | 6 | refuted | fable / medium / reviewer |
| NR-C1 Spec | P1 second preimage check met by inspection | no-op | 7 | held | fable / medium / reviewer |
| NR-C1 Coverage | C1 same-session replacement untested | auto-fix | 7 | held | fable / medium / reviewer |
| NR-C1 Coverage | C2 non-positive author limit untested | auto-fix | 4 | held | fable / medium / reviewer |
| NR-C2 Standards | S1 form names spelled twice | auto-fix | 8 | held | fable / medium / reviewer |
| NR-C2 Spec | P2 completion help lacks the note | partial | 6 | refuted | fable / medium / reviewer |
| NR-C2 Spec | P3 slug placeholder never filled | partial | 4 | held | fable / medium / reviewer |
| NR-C2 Coverage | C1 substituted values unobserved | open | 7 | held | fable / medium / reviewer |
| NR-C2 Coverage | C2 amendment route unobserved | open | 6 | held | fable / medium / reviewer |
| NR-C2 Coverage | C4 unplanned chunk misrouted | open | 4 | held | fable / medium / reviewer |
| NR-C2 Coverage | C5 absent-record chunk value omitted | open | 4 | held | fable / medium / reviewer |
| NR-C3 Standards | S1 size test pastes the checkout harness | ask-user | 6 | held | fable / medium / reviewer |
| NR-C3 Standards | S2 spill-path loop repeated | auto-fix | 5 | held | fable / medium / reviewer |
| NR-C3 Spec | P1 copied-flag sources only in the root row | ask-user | 6 | refuted | fable / medium / reviewer |
| NR-C3 Spec | P2 help qualifier drifts from the table | auto-fix | 3 | held | fable / medium / reviewer |
| NR-C3 Spec | P3 guidance addition not authorized | ask-user | 4 | held | fable / medium / reviewer |
| NR-C3 Coverage | C1 noncanonical envelope unobserved | gap | 9 | held | fable / medium / reviewer |
| NR-C3 round 2 Coverage | C2 Find semantic untested | finding | 7 | refuted | fable / medium / reviewer |
| NR-C2 orchestrator | Source value in stale-plan repair unobserved | candidate | abstained | refuted | fable / high / orchestrator |

Brier mean 0.256 over 20 pairs, 1 abstained.

Each ticket author stayed inside its fence or stopped at it. Ticket 1 resumed from a paused predecessor's preserved diff and committed in 17 minutes. Ticket 2 stopped at a fence gap after a 5-minute survey, then committed 25 minutes after the fence widened. Ticket 3 stopped at a test import cycle that two files outside its fence caused, then committed 19 minutes after the fence widened. Every named probe bit on the first run, and every repair session closed its targets in one commit.

## Coordinator catches

- The ticket 2 and ticket 3 fences each omitted files that a closed decision or an import cycle required. The orchestrator widened them in plan commits before the authors edited.
- The NR-C2 review found the form-name literals in the two files the first fence expansion did not cover, so a second expansion joined the repair.
- The orchestrator's own Coverage candidate for the stale-plan repair source was refuted by the axis, which found the domain test that observes it.
- A repair session hit the verification refusal at the repair tip before the chunk was refrozen, which cost one round trip.
- A late `main` move after the last checkpoint forced a refreeze, a re-record, a third review round, and a second checkpoint.
- Three parallel review axes with one probing axis left the tree clean in all nine rounds.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-record-native-assignments.md | 1 | check-gap |
| 2-enforce-native-record-order.md | 1 | one-source |
| 3-import-complete-probe-results.md | 1 | one-source |

## Agent-experience improvements

### Bench CLI

- Let `bench record verification` take its excerpt from the last response of the worktree. An unspilled response leaves no file, and an author copied 13 excerpts by hand.
  Feeds: new
- Let `bench test` accept a repeated `--package` flag, because an author ran `--changed` to cover several consumer packages.
  Feeds: new
- Let `bench probe` print the omitted or swapped text in its result row and rerun the focused green after the restore. Each probe now needs a second `bench test` call.
  Feeds: FT168
- Let `bench record completion` accept a review that covers only the fold delta past the reviewed chunk tip. A late `main` move cost a refreeze and a full extra round.
  Feeds: new
- Let the build preflight flag a symbol that a closed decision moves out of a file the ticket fence omits. Two fence gaps cost 20 minutes of author survey each.
  Feeds: FT293
- The landing census recorded 9 raw calls: python3 4, rg 2, cat 1, git log 1, ls 1. A `bench assessment stage` verb that appends a timing event would own the python3 calls.
  Feeds: new

### Skills

- Name the chunk refreeze order in `craft-delegate`'s repair-charge template, because a repair session returned on the verification refusal before the orchestrator refroze the chunk.
  Feeds: none
- State in `craft-review` that a confirming round runs its three axes in parallel with one probing axis. Four such rounds kept the tree clean and cut each round to 2 to 5 minutes.
  Feeds: none

### Process

- Fold `main` inside the last chunk immediately before its confirming round, and compare `main` again right before the checkpoint. A move after the checkpoint costs about 25 minutes.
  Feeds: FT342
- Count the lines of the largest ordinary answer against the response bound before a spec adds a table to a bounded response. The `probe_execution` row now spills every biting probe answer.
  Feeds: none
- Run `bench gate-prose` on the review record before every record commit, because three prose refusals each cost a commit round trip.
  Feeds: none