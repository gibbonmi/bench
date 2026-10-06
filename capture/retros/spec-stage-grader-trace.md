## Outcome

The build delivers FT376 in one chunk and one ticket, and it landed at `5a2dc13e` with a green gate. The `Before the map locks` checklist in `map-discipline.md` now traces pin operators, entry reads, and derived expectations to their graders. It also maps consolidated rules, checks quantified obligations at each site and in each ticket, and traces workflow-step writes. A new bullet makes the caller sweep run `bench consumers`, unexported functions included. Nine anchors and nine canaries hold the rules, and all 21 acceptance rows reconcile as covered.

The ticket had a fresh Opus author at high effort, and one fresh Opus/high session ran the repair. By reviewer direction, Fable/high ran all six review axes. GT-C1 used one repair cycle of two.

## Gate-stage timings

- landing: commit 5a2dc13e6cc6c27caeaadbb01464283a0a1706aa, trace f358d3c6c5e2e60df80ebd6062dfab94
- gofmt: 127 ms
- vet: 1330 ms
- test: 252095 ms
- race: 3105 ms
- system: 79157 ms
- shellcheck: 541 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 1 | GT1-GT19 red then green, seven named probes bit | verified | 9 | held | claude-opus-5-5 / high / implementation |
| ticket 1 repair | GT20, GT21 red then green, nine named probes bit | verified | 9 | held | claude-opus-5-5 / high / repair |
| GT-C1 Standards S1 | a changelog sentence has 32 words | verified | 8 | held | claude-fable-5-1 / high / review |
| GT-C1 Standards S2 | a test comment has a four-noun cluster | verified | 4 | held | claude-fable-5-1 / high / review |
| GT-C1 Standards S3 | the record shows only the N1 probe | verified | 3 | missed | claude-fable-5-1 / high / review |
| GT-C1 Spec P1 | the chunk table omits the conformance row | verified | 5 | held | claude-fable-5-1 / high / review |
| GT-C1 Coverage C1 | N3 and N4 have no canary | verified | 7 | held | claude-fable-5-1 / high / review |
| GT-C1 confirming Coverage C2 | a needle without its leading "No" stays green | verified | 3 | missed | claude-fable-5-1 / high / review |
| GT-C1 confirming Coverage C3 | the N4 probe has no record | verified | 4 | missed | claude-fable-5-1 / high / review |

Brier mean: 0.122 over 9 pairs. Abstentions: 0.

## Coordinator catches

- The spec-stage handoff called the build ready, but the staged spec was not on `main` and FT376 had no deliverable binding. The coordinator landed the spec, planned and landed the binding, and then started the commitment.
- A commitment change landed on `main` during the build. After that, every commit on the build branch refused. The coordinator kept the dirty record in a reset envelope, merged `main` before the first review, and moved the chunk base.
- The author recorded a probe exit code of 0. The checkpoint refuses that value, so the coordinator sent the author back to append a superseding entry.
- The merge gate failed once on `TestResumeLandCommandRepeatsTheCensusCount` with an evidence-absent refusal. The same merge passed on retry.
- The Spec axis found that the coordinator's own plan amendment omitted the new conformance row from the chunk table.
- The coordinator refuted S3 and C3 against the author returns, which list every named probe with its bite and its restore.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-trace-spec-rows-to-graders.md | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Make `bench probe` print the exit code of the probed test run, because `bench record verification` requires that value.
  Feeds: new
- Act on the `FT376-build` landing census entry of 17 raw calls, which `bench learning` records with its verb heads.
  Feeds: new
- Make `bench commit` admit a record-only branch commit when only the commitment policy on `main` changed.
  Feeds: new
- Make `bench test --full` show every failing assertion line of a test, so one red run can show each rule of a multi-rule harness test.
  Feeds: new

### Skills

- Make the spec probe list name every anchor needle with its expected verdict.
  Feeds: new

### Process

- Close a spec stage only after the staged spec and its deliverable binding land on `main`, and name the binding state in the handoff.
  Feeds: new
- For a Markdown path that the prose gate excludes, tell the author to apply the sentence bound by hand.
  Feeds: new