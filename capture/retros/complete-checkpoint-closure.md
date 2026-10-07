## Outcome

FT392 landed at `f84a9809` from source `ea3e7867..923f9941`, and the landing gate was green. `bench gate --checkpoint <spec> --complete` now grades the tree that the landing publishes. It composes that tree with `published.Tree` from the committed tip and grades it through the prospective run under the completion obligation. A dirty checkout refuses before the oracle runs. Two tickets in one chunk covered 32 acceptance rows, and the delivery closed `roadmap/FT392.md`. The complete checkpoint of this build ran on the new route and was green.

## Gate-stage timings

- landing: commit f84a9809e6f2c9a3f013989b42f701c71e7a4f8a, trace 20bbbb01fa208f5d852a6e69d490cef8
- gofmt: 140 ms
- vet: 1362 ms
- test: 270219 ms
- race: 3082 ms
- system: 80654 ms
- shellcheck: 520 ms

## Ticket-versus-spec-slice and delegate performance

Two fresh Opus/high authors each committed their ticket first-pass inside the planned fence, and each named probe bit with an exact restore. The ticket 02 author needed none of the 36 approved closure paths. One fresh Opus/high repair session closed S1, S2, and C1 in one cycle. Opus/high ran all six review axes, and the confirming round found 0 findings.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| line declaration | 1 expected repair round | held | 0.6 | expectation | Opus / medium / orchestrator |
| CC1 Standards | S1 test respells cleanCheckoutRefusal | held | 0.6 | finding | Opus / high / reviewer |
| CC1 Standards | S2 comment claims one graded tree | held | 0.6 | finding | Opus / high / reviewer |
| CC1 Standards | S3 published.Tree pairing built twice | refuted | 0.4 | finding | Opus / high / reviewer |
| CC1 Coverage | C1 ignored-file pass untested | held | 0.7 | finding | Opus / high / reviewer |

Brier mean: 0.146 over 5 pairs; 0 abstained.

## Coordinator catches

- `bench commitment start` refused twice: the primary checkout owns no assignment, and the policy bound no deliverable to FT392. The reviewer approved a deliverable plan, which landed at `ea3e7867`.
- The CC1 chunk checkpoint refused with a stale plan digest, because the coordinator recorded the chunk update before the repair plan amendment. The record was rebuilt in the correct order from the committed state and from the retained excerpts.
- The ticket 02 author saw no red for CC04, CC08, CC10, and CC31. The Coverage axis then turned each of the four rows red under a probe.
- The authors recorded a probe exit code of 1 from the failing test run, because `bench probe` prints no mutated-run exit code.
- The coordinator did not read the Claude scorecard before it selected the line. Three of its recorded decisions would have prevented the two refusals and the late census read.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-move-fixture-witnesses-to-the-common-directory.md | 0 | none |
| 02-grade-the-published-tree-at-the-complete-checkpoint.md | 1 | one-source |

## Agent-experience improvements

### Bench CLI

- Make `bench record chunk` refuse a plan digest that the record cannot map, and name `bench record amendment` as the next command.
  Feeds: new
- Make the deliverable refusal of `bench commitment start` name `bench commitment plan` with the deliverable binding as its next command.
  Feeds: new
- Make `bench probe` print the exit code of the mutated run, so that `bench record verification` takes it from the probe output.
  Feeds: new
- Add a `bench record` form that appends a version 2 author assignment, so that no orchestrator edits the plan fence with a script.
  Feeds: new
- The census learning for this landing records 27 raw calls: rg 20, sed 4, python3 2, ls 1.
  Feeds: none

### Skills

- Make `/bench-write-spec` bind the staged spec as the outcome deliverable in the same spec-stage landing.
  Feeds: new
- Make `/bench-implement-spec` state the record order after a plan commit: the amendment, then the chunk, then the verification.
  Feeds: new

### Process

- Read the provider scorecard before the line declaration, as `craft-line` states, because its decisions name the record order and the census read.
  Feeds: none
- Let the orchestrator close the staleness pass without a dispatch when the drift touches no path that the spec or a ticket names.
  Feeds: new