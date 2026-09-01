# Retro: spec-authoring-discipline

## Outcome

The landing commit `d312714f` published `specs/spec-authoring-discipline/spec.md` as
`Status: implemented` and closed FT220, FT257, and FT278. The reviewed pair was
base `34889eef` and source tip `7bc1803d`. Six tickets and one repair ticket landed
through the retained integration source `sad-integration`. The diff touched 89
files with 853 insertions and 43 deletions, and it added 25 canary fixtures.

## Gate-stage timings

The landing gate ran once over the composed pair. Each earlier fold into the
integration source paid one prospective gate of the same shape.

| phase | elapsed |
| --- | --- |
| gofmt | 94 ms |
| vet | 876 ms |
| test | 52.1 s |
| race | 2.3 s |
| system | 21.2 s |
| shellcheck | 461 ms |

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized. Seven write charges ran at Opus medium, and six
landed first-pass on behavior. The reader-sweep charge stopped on an out-of-fence
blocker and cited the two pinning tests. It then finished first-pass under one
coordinator direction. Every delegate returned a red-before
and green-after log per row and a self-probe that bit. Four review charges ran at
Opus medium: three axes and one repair-scoped re-review. The axes returned 12 raw
findings that collapsed to one repair target.

## Coordinator catches

- A coordinator probe on a wrapped sentence matched nothing, and the gate stayed
  green. A `cmp` against the copy aside showed no byte changed. The corrected
  probe went red at both readers.
- The conformance package alone cannot show a live-tree anchor red. The gate entry
  test with `BENCH_CONFORMANCE_ROOT` set to the worktree can.
- The reader-sweep delegate first reported the pinning tests as a fence block. The
  pinned bytes end at `rule.`, so a second sentence after them keeps both pins.
- The review found a bare "map" on a reflowed line that no delegate self-audit
  caught.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| add-a-sources-example-to-the-map-template | 0 | none |
| require-a-row-for-each-in-scope-promise | 0 | none |
| split-the-wrapped-field-refusal | 0 | none |
| name-the-reader-sweep-and-move-the-ship-test | 1 | spec-row |
| add-the-decision-map-authoring-steps | 0 | none |
| define-the-map-glossary-terms | 0 | none |
| repair-bare-map-on-moved-lines | 0 | none |

## Agent-experience improvements

### Bench CLI

- Give `bench worktree exec` a probe form that copies a file aside, runs the
  command, and restores the file. The census entry `sad-integration landed with
  census 11` then records no shell-variable head and no `rm`.
  Feeds: new
- Let `bench gate-prose` refuse a file as the root operand with a usage message
  and print the offending line on a fail verdict.
  Feeds: new
- Make the fixture materializer collapse whitespace the way the anchor evaluator
  does, or name the wrap in its refusal.
  Feeds: new
- Add a non-gate handle that evaluates the anchor registry against a live tree.
  Feeds: new
- Make `bench test` accept the `--check <name>` flag its usage line advertises.
  Feeds: new

### Skills

- Add to `craft-spec`'s reader sweep an `rg` over `tests/` and `internal/conformance`
  for the literal bytes of every sentence a build moves.
  Feeds: new

### Process

- Require the coordinator to run `cmp` against the copy aside before it reads any
  probe verdict.
  Feeds: none
- Run a live-tree anchor probe through the gate entry test with the worktree as
  the conformance root, never through the fixture-bite test alone.
  Feeds: none
