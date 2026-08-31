# Retro — otel-seam-record (FT274)

## Outcome

FT274 landed. The published commit is 2c4d263b on main, over the reviewed pair base effc457b to tip 95423686. The landing carries 13 tickets, one build repair, five review-repair commits, and two fence amendments. The coverage map validates with 28 rows. The record writes OTLP-JSON spans for the gate, phase, lane, commit, landing, worktree, and hook seams.

## Gate-stage timings

The landing gate ran green: gofmt 91 ms, vet 865 ms, test 52986 ms, race 5032 ms, system 18774 ms, shellcheck 472 ms.

## Ticket-versus-spec-slice and delegate performance

Every write charge was ticket-sized; no charge received a spec slice. Twelve of thirteen ticket charges landed first-pass on behavior. Each delegate ran a mutation probe that bit, and each reported its unsure calls. The one repair charge closed nine review targets in four commits, first-pass.

## Coordinator catches

- The ticket 12 delegate edited `cmd/bench/command_registry.go` outside its Writes list; the review preflight confirmed the fence gap.
- The landing refused seven repair files that no fence row named; the fence gained them before the re-land.
- The prose lane refused four sentences and one paragraph in the pickup artifact before its commit.
- The coordinator probes bit at distinct sites: the commit path count, the exec subject, the hook seam prefix, and the FIFO refusal.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| add-the-pinned-otel-dependencies | 0 | none |
| move-the-bench-home-read | 0 | none |
| add-the-encoder-and-the-attribute-set | 0 | none |
| write-the-keyed-record-file | 1 | spec-row |
| write-the-start-and-end-lines | 0 | none |
| record-the-gate-run-and-its-phases | 1 | spec-row |
| keep-the-started-phase-line-after-a-kill | 0 | none |
| check-that-each-registered-seam-starts-a-span | 0 | none |
| record-the-lane-span-with-its-failing-check | 0 | none |
| record-the-commit-and-landing-spans | 1 | ticket-slicing |
| record-a-span-per-worktree-verb | 1 | delegate-error |
| record-a-span-per-hook-plumbing-verb | 1 | ticket-slicing |
| declare-the-harness-measure-cells | 0 | none |

## Agent-experience improvements

### Bench CLI

- Charge templates name the worktree path once, because this landing's census learning records 38 raw calls with cd, ls, and W= heads.
  Feeds: none
- The fold through `bench worktree merge` runs the fence check from `bench preflight`, so an out-of-fence file surfaces at the fold, not at the landing.
  Feeds: new

### Skills

- The `bench-craft-delegate` charge shape includes the spec fence list, so a delegate reports an out-of-fence write before it edits.
  Feeds: none

### Process

- The coordinator runs `bench preflight build` after each fold, so a fence gap attributes to its own ticket.
  Feeds: none
