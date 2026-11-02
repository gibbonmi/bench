## Outcome

FT406 landed at `75ae8da7` on a green landing gate. A lane-passed `bench commit` in a kit worktree now rebuilds a stale `dist/bench` before its background grade. It builds an absent `dist/bench` only when the diff names a build input. `freshness.DiffChangesBuildInputs` and `freshness.Refresh` own the diff test and the refresh rule for the landing and the commit. `git.ParseNameStatusZ` owns the NUL name-status parse. The stale-build known issue retired, and an anchor forbids its return.

The build was a `--full --delegate --auto-approve` run with `--reviewer opus medium`. The plan set `execution.review_mode` to `unified`, so one session ran the three axes of each round. CW-C1 used both repair cycles, and CW-C2 used one.

## Cost

- landing: 75ae8da754cbfdf7ef9caa29b2bceedbd0ad0eb1, range 75ae8da754cbfdf7ef9caa29b2bceedbd0ad0eb1^1..75ae8da754cbfdf7ef9caa29b2bceedbd0ad0eb1^2
- wall: 120 min
- commits: 25 total, 19 bookkeeping, 6 change
- graph: 2 tickets, chain 2, 2 chunks, 2 serial checkpoints
- dispatch: 2 authors, 3 repairs
- trace: commit 26 spans in 241322 ms, lane 0 spans in 0 ms, landing 1 spans in 307814 ms
- wall per ticket: 60.0 min against 60 min
- bookkeeping per ticket: 9.5 against 7.1, over budget
- delegates: about 1.08M subagent tokens over 11 sessions; the ticket 02 author used 242k, and the four unified review rounds used 298k together.
- provider usage and charges: unknown.

## Gate-stage timings

- landing: commit 75ae8da754cbfdf7ef9caa29b2bceedbd0ad0eb1, trace b7ef9c466b0c9dcfd2009772159b913e
- gofmt: 203 ms
- vet: 1715 ms
- test: 226585 ms
- race: 3662 ms
- system: 49695 ms
- shellcheck: 547 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| CW-C1 ticket 01 | CW01 diff test names a build input | verified | 9 | held | opus / medium / implementer |
| CW-C1 ticket 01 | CW46 kept test grades the build-failure text | claimed | 5 | refuted | opus / medium / implementer |
| CW-C1 ticket 01 | CW44 landing reports a deleted input | verified | 9 | held | opus / medium / implementer |
| CW-C2 ticket 02 | CW13 stale build rebuilds once | verified | 9 | held | opus / medium / implementer |
| CW-C2 ticket 02 | CW38 anchor turns red on the retired row | verified | 4 | refuted | opus / medium / implementer |
| CW-C2 ticket 02 | CW45 declared build still rebuilds | verified | 9 | held | opus / medium / implementer |
| CW-C1 repair 1 | CW46 closed | verified | 9 | held | opus / medium / repair |
| CW-C1 repair 1 | CW51 closed | verified | 9 | held | opus / medium / repair |
| CW-C1 repair 2 | O1 one parser owner | verified | unknown | abstained | opus / medium / repair |
| CW-C2 repair 1 | CW52 closed | verified | unknown | abstained | opus / medium / repair |
| CW-C1 review 1 | C2 quoted path misses a deleted input | claimed | 6 | held | opus / medium / reviewer |
| CW-C1 review 1 | S1 comment repeats the owner | claimed | 3 | refuted | opus / medium / reviewer |
| CW-C2 review 1 | C1 dangling build has no row | claimed | 7 | held | opus / medium / reviewer |
| build | one repair round for each chunk | claimed | 6 | refuted | opus / medium / orchestrator |

Brier mean 0.098 over 12 pairs; 2 abstained.

## Coordinator catches

- `bench consumers` found two `buildInputs` callers in `internal/freshness/freshness.go`, and ticket 01 took that widen before dispatch.
- A coordinator probe at the stack-landing caller stayed silent over 44 tests. No row covers that caller, so a learning holds it.
- The CW-C1 repair added a second NUL name-status parser beside `parseNameStatusZ`. The orchestrator routed it as O1 to repair cycle 2.
- The `known-issues` plan probe read a synthetic-tree test and stayed silent. A coordinator probe showed that `docs-currency-workflow` bites, and a plan commit pointed the requirement at it.
- The CW-C2 confirming round found that the spec says a dangling `dist/bench` rebuilds, but the publication refuses a link. The build follows the tree, and the wording waits for reviewer veto.
- The orchestrator did not read this scorecard before the line declaration. It recorded a no-op finding with a finding ID, so the Standards axis read `fail`, and it repeated a plan probe that the decisions already warn against.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-shared-freshness-owners.md | 2 | spec-row, one-source |
| 02-commit-rebuilds-worktree-build.md | 1 | edge-found |

## Agent-experience improvements

### Bench CLI

- Let a spill file hold the full native response, so that a Read of a spill never spills again; the scaffold, a probe, and a test run each took two reads here.
  Feeds: new
- Let `bench record assignment` accept a staged spec edit in the same commit, so that a repair round takes one plan commit and not two.
  Feeds: new
- Start no background grade for a commit whose paths are only Markdown records, because one such grade sealed the cache and refused the next lane commit.
  Feeds: new
- Let `bench probe --mutants` stage each row under a `--requirement` name, so that an author runs no paired `bench record excerpt` calls.
  Feeds: new
- Census: 4 raw calls. Two `cat` calls wrote mutant-set and scratch files, an owned gap that an inline mutant flag removes. One `git rev-parse` compared HEAD with `main`, a missed form because each Bench response prints the head. One `rg` was a sanctioned read.
  Feeds: none

### Skills

- Let the staleness pass omit its delegates when the reviewed graph commit is the source tip and the audit drift list is empty; two slices spent 57k tokens to return current.
  Feeds: new
- Require in `/bench-write-spec` one recorded red for each kept-rule row that maps to an existing test, because CW46 named a test that never graded its text.
  Feeds: new
- Require in `/bench-write-spec` that each handled edge of the edge inventory names its row, because the dangling `dist/bench` edge had no row.
  Feeds: new

### Process

- Keep `execution.review_mode` `unified` for a spec of two or three tickets; four unified rounds found nine real findings for about 298k tokens.
  Feeds: none
- Retire known-issues row 13 in the next drain, because three freezes after an assignment record passed with no amendment record.
  Feeds: new
