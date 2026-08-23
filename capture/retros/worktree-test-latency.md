# Retro — worktree-test-latency

## Outcome

The spec landed at `eab92796` on 2026-08-23 (reviewed pair `4a8aa16a` →
`8226a2df`). The worktree suite now selects one Bench executable and resolves
effect inputs at one boundary. It owns landing, lifecycle, and reclaim
decisions in three pure packages, and runs serial proof-backed journeys
through one harness. The package span fell from a 125.790-second baseline median to
56.898 seconds. The evidence lives in
`specs/worktree-test-latency/evidence/demand-reduction.md`.

## Gate-stage timings

The landing gate ran six phases: gofmt, vet, test, race, system, and
shellcheck, all green. Phase-level wall time is not printed; the package rows
are the retained measurement. `internal/worktree` ran 57.769 seconds in the
landing gate, and the three policy packages ran under 10 milliseconds each.
`internal/systemtest` ran 8.157 seconds, and the race phase ran about 3
seconds. Across the build's twelve gate runs, the worktree package span held
between 55.2 and 58.0 seconds.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized on one retained integration worktree; no charge
received a whole spec slice. Eight write charges ran at `fable` / medium:
seven tickets and one review repair. Six of eight landed first-pass at the
diff level. Ticket 04 returned once for a conformance red, and the ticket 07
and repair charges each returned once for a prose bound. Every delegate
self-probe bit, and every coordinator probe at a distinct kind and site also
bit. The three review axes and the scoped re-review ran at `sonnet` / high
with zero re-runs.

## Coordinator catches

- No delegate done-claim was false; all eight independent probes went red by
  name and restored green.
- The coordinator caught its own lost worktree request token and recovered it
  through `bench worktree reauthorize` after a digest brute-force failed.
- The coordinator caught two of its own CLI-contract violations (appended
  pipelines on `bench commit`) and captured them in `capture/learnings.md`.
- The landing refusal on `reviews/worktree-test-latency.md` exposed a fence
  gap; the fix was a spec fence amendment on the finding cadence.

## Repair attribution

| ticket | repair rounds | causes |
|---|---|---|
| 01-select-one-test-run-binary | 0 | none |
| 02-add-explicit-effect-inputs | 0 | none |
| 03-extract-landing-policy | 0 | none |
| 04-extract-lifecycle-policy | 1 | spec-row |
| 05-extract-reclaim-policy | 0 | none |
| 06-contract-to-serial-journeys | 0 | none |
| 07-record-reduced-demand | 1 | delegate-error |
| review repair (precedence cases) | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

- Make `bench worktree create` record its request token in the assignment
  record, so a later landing can recover it without a reauthorize round.
  Feeds: new

### Skills

- Make `bench-craft-spec` require that a code-moving spec names every
  conformance-pinned consumer of the moved symbols. Ticket 04's only red
  came from an unnamed pinned consumer.
  Feeds: new

### Process

- Batch capture-file commits (ideas, learnings, handoff) into one gate run,
  because each solo capture commit pays a full five-minute gate.
  Feeds: none
- Fence the review pickup path `reviews/<slug>.md` in the spec's ownership
  fences at spec-author time, so a landing never refuses on the review
  phase's own artifact.
  Feeds: new
