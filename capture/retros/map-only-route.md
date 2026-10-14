## Outcome

The `map-only-route` build landed at `38224aac` with a green gate. One ticket in one chunk delivered the reviewer-selected map-only route. It added a new `craft-spec` reference, one pointer sentence in `/bench-write-spec`, and 22 `RequireInSection` anchors with an independent test. It also updated the project profile, the field guide, and the changelog. The full route stays the default.

The build ran as `--full --delegate` with an author limit of 1. By reviewer direction, `gemini-4-argon-high` wrote the ticket through the `agy` CLI, and Opus/low wrote both repairs. Opus/high ran nine review axes over three rounds, and a Fable/medium session orchestrated.

## Gate-stage timings

- landing: commit 38224aac1e719dd973daed40fd4d0751550de269, trace 22382852f84fa2b7a4713a31e42d1966
- gofmt: 147 ms
- vet: 1429 ms
- test: 175806 ms
- race: 3460 ms
- system: 57586 ms
- shellcheck: 564 ms

## Ticket-versus-spec-slice and delegate performance

The author stayed inside its fence, touched 7 of 27 `Writes:` entries, and named a reason for each untouched one. It showed the one required red first and passed every floor check on its first attempt. The session took about 40 minutes of wall time, and about 30 of them went to the floor runs and the eight record calls.

Review round 1 found 7 repair targets. Three were prose accuracy slips in the field guide and the project profile. The others were a duplicated path literal in the test, a passive comment, and a sentence that copied a false claim from the spec. The orchestrator's cycle 1 expansion added a stop that the map validator already enforces, and the confirming round caught it.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| line declaration | expected repair rounds 1 | claimed | 6 | refuted (2 cycles) | fable / medium / orchestrator |
| ticket 1 return | acceptance rows 1-4 pass their named checks | verified | 10 | held | gemini-4-argon-high / high / implementer |
| ticket 1 return | acceptance rows 5-10 (readers, exclusions, reference content) | claimed | 10 | held | gemini-4-argon-high / high / implementer |
| repair 1 return | P1, P2/S1, S2, S4, S5, S6, C1 closed | verified | 7-9 | held | opus / low / repair |
| repair 2 return | S8 and P3/C2 closed | verified | 9 | held | opus / low / repair |

Brier mean 0.038 over 20 pairs, 0 abstentions.

## Coordinator catches

- The author's done-claim held under an independent swap probe on the pointer, which the live `root-conformance` check caught.
- The cycle 2 author recorded its eight verifications as hand-written summaries with a placeholder probe command. The orchestrator sent it back to re-record from the raw spill files.
- Every charge said `--probe-exit-code 0` for a bit probe; `bench record completion` refused, because that cell is the mutated run's exit code. The author re-recorded the entry with 1.
- The first cycle 1 plan commit failed the prose lane on two sentences over 25 words. The first `bench learning` call failed the same bound three times.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-add-the-map-only-route.md | 2 | delegate-error, spec-row |

## Agent-experience improvements

### Bench CLI

- Collapse the per-ticket record ceremony into one `bench record verification` call that imports every floor excerpt of a plan chunk. Chain the freeze, the record commit, and the post-ticket preflight into one step.
  Feeds: new
- State in the `bench record verification` help that `--probe-exit-code` is the mutated run's exit code, and refuse 0 with a message that names the `probe_execution` row.
  Feeds: none
- Let `bench record amendment` answer a clear no-op before the first chunk entry exists, so an author charge needs no retry rule.
  Feeds: none

### Skills

- Make the `craft-delegate` charge template name the `probe_execution` cell that fills `--probe-exit-code`, and require a raw native excerpt for each record.
  Feeds: none
- Add to `craft-review` that a finding which asks for a new stop or refusal cites the validator's current diagnostics for that input first.
  Feeds: none

### Process

- Make the orchestrator read the validator before it takes a pre-approved spec expansion that adds a stop. The MR34 expansion on this build restated a diagnostic that `bench maps` already gives.
  Feeds: none
- Compare the two gemini guidance trials in the scorecard before any routing change. Both passed every check first-pass, and both took one prose repair round.
  Feeds: none