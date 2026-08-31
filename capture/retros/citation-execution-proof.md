# Retro — citation-execution-proof (FT133)

## Outcome

FT133 landed. The published commit is 5b5daa51 on main, over the reviewed pair base 61d3e8a0 to tip b827d51f. The landing carries 7 tickets, one hermeticity repair, one review-repair commit, and one fence amendment. It also closes the pickup. The coverage map validates with 27 rows.

A citation now proves an executed test. A mention and a stale subtest red. A mixed-tag map refuses. Every allowlist source stats, and `--check` names the uncited rows.

## Gate-stage timings

The landing gate ran green: gofmt 95 ms, vet 838 ms, test 50730 ms, race 4534 ms, system 17707 ms, shellcheck 488 ms.

## Ticket-versus-spec-slice and delegate performance

Every write charge was ticket-sized; no charge received a spec slice. Seven of seven ticket charges landed first-pass on behavior, and the repair charge closed eight review targets in one pass. Each delegate ran a mutation probe that bit, and the coordinator probe at a distinct site and kind also bit on every accepted diff.

## Coordinator catches

- The composed merge gate found a non-hermetic test. The system-set fixture read the ambient `BENCH_KIT`. The coordinator attributed the red to the grade ticket. One inserted line repaired it.
- The canary delegate reported its own out-of-Writes edit to `internal/conformance/registry_test.go`; the fence gained the file before the landing.
- The installed `bench` binary predates the diff, so a `--check` sanity run could not show the new uncited report. The Spec axis saw the report through the built code. The evidence for a new CLI behavior comes from the built tree, not the installed binary.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| derive-executed-tag-census | 0 | none |
| refuse-mixed-tag-row-ids | 0 | none |
| stat-allowlist-sources | 0 | none |
| grade-citation-execution | 1 | delegate-error |
| tighten-citation-grammar | 0 | none |
| report-uncited-rows | 0 | none |
| prove-diagnostics-with-canaries | 1 | ticket-slicing |

## Agent-experience improvements

### Bench CLI

- The census learning "citation-execution-proof census: 193 raw calls" records the heads. Variable-prefix path assignments count near 150, plus cat, cd, ls, and sed. A path-once affordance for worktree file edits removes most of them.
  Feeds: new
- The follow-on hook refuses a non-Bench segment that precedes a `bench worktree exec` on the same line. An edit plus its test run then costs two calls. Five delegates hit it.
  Feeds: new
- `bench worktree exec` carries no timeout, so a wedged child (the FIFO probe) holds the run until the harness kill and the red is unreadable.
  Feeds: new

### Skills

- The `bench-craft-tickets` Writes guidance names every registry an enumerated family joins, because the canary ticket omitted `internal/conformance/registry_test.go` and forced a fence amendment.
  Feeds: none

### Process

- A charge for a test that needs the system tag set pins `BENCH_KIT` in the fixture. The gate exports an ambient kit, and that ambient kit flips the census under composition.
  Feeds: none
