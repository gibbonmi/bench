# Retro: go-build-cache-footprint

## Outcome

The spec landed at `5df9c8c3` on 2026-08-27 from source pair `408c160e` to
`8cd1e045`. The first review pair was `caeb19fb` to `95163e15`. Bench now
owns one Go build cache at `$HOME/.cache/bench/go-build`, and every Bench Go
argv carries `-trimpath`. Holders share a record lock. The verbs `bench cache`
and `bench cache clean` exist, and a green gate prints one `go-build-cache:`
line.
The measured fresh-path gate growth fell from 198,839,178 bytes and 2,621
files to 258,482 bytes and 197 files.

## Gate-stage timings

| stage | elapsed |
|---|---|
| landing gate: gofmt | 90 ms |
| landing gate: vet | 708 ms |
| landing gate: test | 45,370 ms |
| landing gate: race | 2,215 ms |
| landing gate: system | 9,874 ms |
| landing gate: shellcheck | 444 ms |
| ticket lanes (gofmt, prose, vet, build), 17 runs | about 30 s each |
| whole gates on the source, 6 runs | about 75 s each |

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized; no charge took a spec slice. Eleven ticket and
repair charges ran on Opus at low or medium effort. Eight landed first-pass on
behavior. Ticket 02 took one comment-placement fix. Ticket 04 edited one file
outside the fence.

Ticket 06 stopped correctly on a spec contradiction (C12). The eight review
axes and scoped checks each returned cited findings with no uncited claim.
The merge delegate composed three seams into a rewritten package. It reported
the one-parent commit instead of a raw fallback.

## Coordinator catches

- The whole gate under `-trimpath` reddened `internal/bounds`, which story 15
  did not name; the fence widened by one test file.
- Ticket 04's delegate edited `bin/bench.sh` outside the fence; the coordinator
  reverted it before the preflight.
- Ticket 05's `holdWait` literal restated a bounds policy value, and its
  helper-process skip counted as an environment skip in the kit gate.
- The C12 seam named a consumer-repo journey that never reaches the report
  tail; the reviewer moved the row.
- The `bench commit` after `git merge` recorded one parent, so the base was
  not an ancestor until a raw commit recorded the second.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01 test roots | 0 | none |
| 02 cache path | 1 | delegate-error |
| 03 trimpath argvs | 1 | spec-row |
| 04 bench cache verb | 1 | delegate-error |
| 05 lock and clean | 2 | tree-drift, tree-drift |
| 06 gate line and log | 1 | spec-row |
| 07 measurement | 0 | none |
| 08 revert | 0 | none |
| 09 bind holders | 0 | none |
| 10 single-source refusals | 0 | none |
| 11 reporter fallback | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Let `bench commit` record `MERGE_HEAD` as a second parent, or make the
  landing refusal name the raw commit that records the merge.
  Feeds: new
- Add a Bench projection of a checkout's changed-path list, so a commit path
  list never comes from `git status`. The census entry
  `go-build-cache-footprint census 372` in `capture/learnings.md` names the
  heads.
  Feeds: new
- Give `bench gate` a root argument, so a measurement on a scratch checkout
  needs no `env -C`.
  Feeds: none

### Skills

- `bench-craft-spec`: a spec that changes a compile flag lists the reds from
  one hand run of the whole gate under that flag, not a grep for one API.
  Feeds: new
- `bench-craft-delegate`: a charge names the fence as a refusal boundary above
  any mirror-the-sibling instruction.
  Feeds: none
- `bench-craft-tdd`: a re-exec helper test returns silently outside its role
  env; it never skips, because the kit gate reds an environment skip.
  Feeds: none

### Process

- Record every mid-build reviewer decision in the checkpoint ticket at the
  time it is taken, so a later review axis finds it.
  Feeds: none
- Promote the Bench build cache decisions to an ADR: one owned directory, a
  shared lock, no eviction inside a gate, and the 10 GiB bound.
  Feeds: new
- When two sessions share one machine, the second landing composes `main`
  through the refusal's prescribed merge. Run the scoped review on the merge
  delta before the retry.
  Feeds: none
