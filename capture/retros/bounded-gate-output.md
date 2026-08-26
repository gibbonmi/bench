# Retro: bounded-gate-output (FT185)

## Outcome

The gate engine now buffers each phase's stream instead of relaying it. A
green run prints a bounded `phases[N]{phase,verdict,elapsed_ms}` table, one
`capability-skips` line, then `gate: green`. A red run prints a bounded
`failures[N]{phase,line}` table, then `gate: red`. The complete stream is
retained at `.logs/gate-<run>.out`, beside the progress log, under the same
twenty-run retention. `bench commit` and `bench worktree land` relay the
same engine output unchanged.

Landed at `253d01ec` on `main`, published by `bench worktree land` from the
retained integration source at base `81e3415f`. The spec retired at
`5df2f5ca`.

Two review passes ran. The first was an initial three-axis pass over the
seven build tickets. The second was a repair-scoped re-review over the
four accepted repairs. Both passes are clean. The repair-scoped pass found
only stale-wording follow-ups in already-landed ticket files, since fixed.

## Gate-stage timings

From the landing's own gate run (six phases, all green):

| phase | elapsed |
| --- | --- |
| gofmt | 89ms |
| vet | 1.66s |
| test | 47.8s |
| race | 4.83s |
| system | 10.3s |
| shellcheck | 499ms |

## Ticket-versus-spec-slice and delegate performance

The write-spec phase (a prior session) sliced the spec into seven
ticket-sized charges (01–07). Each charge carried an explicit `Writes:`
fence and its own acceptance rows. All seven landed first-pass on
behavior; no delegate needed a second charge to fix its own diff. Tickets
04 and 05 were authored as separate delegates from the same base (03's
tip), per the spec's own parallel-authoring note. They then landed
serially, 04 before 05, to avoid a `report_test.go` collision the spec had
already flagged. Tickets 06 and 07 ran genuinely in parallel on disjoint
files (tests and docs), with no contention.

The four review repairs (A–D) were charged smaller than a full ticket
each, scoped to one review finding cluster per delegate. They shared files
only where unavoidable: `report.go`, across repair A and one trailing
comment fix. All four landed first-pass on behavior. Two ran in parallel
per round (A with C, then B with D) on disjoint files. Both rounds landed
clean.

## Coordinator catches

- Every one of eleven delegate charges (seven tickets, four repairs) got an
  independent coordinator mutation probe. Each probe ran at a distinct
  site and kind from the delegate's own probe. Every probe bit correctly.
  No delegate's self-probe was masking a gap; that is itself evidence the
  delegate probes were real.
- After amending the spec's acceptance coverage map with three new rows
  (BG37–BG39) for the accepted repairs, `bench preflight review` failed
  `rows-owned`: the new rows cited no ticket file. This is a real process
  gap. A review-approved repair changes the spec but not a ticket file,
  unless the coordinator adds one. Fixed by writing ticket 08, which names
  the four repairs retroactively.
- The repair-scoped re-review caught four non-behavioral spec
  contradictions. Two superseded ticket-file acceptance bullets (01, 02)
  still described the pre-repair `WARNING: DATA RACE` fallback behavior.
  One story's rationale clause (17) named a race report as the fallback's
  example, after the repair had moved races out of that fallback. One
  row's wording (BG26) still said "unwritable" after its fixture moved
  from a chmod to a symlink. All four were wording-only. The coordinator
  fixed them directly rather than delegating, since each was a
  single-paragraph text correction in an already-committed spec artifact.

## Repair attribution

| ticket | repair rounds | cause |
| --- | --- | --- |
| 01 | 1 | spec-row |
| 02 | 1 | spec-row |
| 03 | 1 | ticket-slicing |
| 04 | 0 | none |
| 05 | 1 | ticket-slicing |
| 06 | 1 | ticket-slicing |
| 07 | 0 | none |
| 08 (review repairs) | 0 | none |

Ticket 01's round: the classifier's fallback-only "no row" behavior did not
match how `go test -race` actually interleaves a `WARNING: DATA RACE` line
with a `--- FAIL:` block. The acceptance row (BG32) tested a fixture no
real race report produces. Ticket 02's round: the acceptance row BG01 said
"and nothing else," but the engine printed the capability-skips totals
line on every run, not only green. Tickets 03, 05, and 06's rounds each declared a seam: the `capability`
reservation, the `.out`-retention wiring, and BG26's partial-failure
ordering. The owning ticket's own test did not actually exercise it.

## Agent-experience improvements

### Bench CLI

- `bench preflight review` refuses when an amended coverage map's row
  cites no ticket file. Nothing prompts a coordinator to check this before
  amending a spec mid-repair. A short note in the spec-amendment path of
  `/bench-review-implementation`'s exit handoff would have caught this a
  step earlier.
  Feeds: FT185
- Census: this landing recorded `census=133` on the main publish and
  `census=5` on the spec-retire housekeeping publish, 138 raw calls total.
  The coordinator's best-recollection breakdown (`git diff`/`status`/`log`,
  `cp`, `sed -i`, `rg`) is captured in `capture/learnings.md`. No bench
  verb covers a coordinator's per-file diff review or a mutation probe's
  edit-run-revert cycle.
  Feeds: none

### Skills

- `bench-craft-review`'s repair-scoped mode worked cleanly once invoked.
  Nothing in `/bench-implement-spec --full` or
  `/bench-review-implementation` names the exact moment a coordinator
  should write the retroactive repair ticket file. Naming it explicitly,
  after repairs land and before the repair-scoped re-review, would remove
  a guess.
  Feeds: FT185

### Process

- The four accepted review repairs shared two files pairwise. `report.go`
  crossed the BG01/StartErr repair and the trailing comment fix.
  `testlines` crossed the classifier repair and its own tests. Scoping
  each repair delegate to disjoint files where possible, and sequencing
  the rest, held up under this landing's size. No repair delegate
  collided with another's concurrent edit.
  Feeds: none
