## Outcome

FT349 landed as `d0b62266` through `bench worktree land`, and the whole-project gate was green. The build preflight now grades each ticket's check floor in the `check-floor` row. `bench test --check root-conformance` is the one command for the unscoped root conformance pass. `bench record chunk` and `bench record completion` refuse a `planned` seam cell. The guidance states the floor, the cell rewrite, the marker lifetime, and the root conformance command. 

The source pair was `9b99dc0c` to `c324a603`. All 74 acceptance rows reconciled, and the complete checkpoint was green.

The build ran five chunks with five fresh Opus/high authors and seven repair sessions. TF-C1 and TF-C2 used both repair cycles, TF-C3 used one, TF-C4 used two, and TF-C5 used none. Review rounds ran on Fable/high, then Opus/high, then Fable/medium by reviewer direction. Confirming rounds ran on Sonnet/high and then Fable/medium. Provider usage and charges remain unknown.

## Gate-stage timings

- landing: commit d0b62266587865469dff55076c9e3d6432046652, trace 764f77322a0efc3e4d48383f51313479
- gofmt: 164 ms
- vet: 1414 ms
- test: 294432 ms
- race: 2998 ms
- system: 80303 ms
- shellcheck: 544 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 01 TF1 | the root-conformance check runs unscoped | verified | 9 | held | Opus / high / implementer |
| ticket 02 TF20 | the floor and the changed selection agree | verified | 8 | held | Opus / high / implementer |
| ticket 03 TF32 | a compound command satisfies nothing | verified | 8 | held | Opus / high / implementer |
| ticket 04 TF67 | the chunk freeze reads the tip | verified | 8 | held | Opus / high / implementer |
| ticket 05 TF60 | the floor sentence is anchored and bites | verified | 9 | held | Opus / high / implementer |

Brier mean 0.028 over 5 pairs; 0 abstained.

Each author worked inside its planned ticket and chunk. Ticket 03 stopped twice before its first edit to ask for a fence. The spec named an unexported owner and omitted the baselines that pin the row list. Ticket 04 reported four copies outside its fence before review. The orchestrator widened the fence, and the same author folded them before the freeze. Every author's named probe bit.

## Coordinator catches

- The orchestrator ran `bench record chunk` before `bench record amendment` after plan commits, so the TF-C1 checkpoint refused a stale plan digest. The record was replayed in order on a new commit.
- A delegate recorded the probe command exit code where the checkpoint reads the mutated run exit code.
- A stale worktree binary turned `binary-seal` red after a ticket commit.
- Three coordinator probes were invalid or equivalent: two mutants left an import unused, and one swapped a commit for the same commit. Each was replaced before acceptance.
- A plan-text correction after the TF-C3 freeze staled the chunk source, so the chunk was frozen again.
- The ticket 03 author found that `bench test --checks` and a fixture listing satisfied a floor element; row TF72 closed that false pass.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-name-the-root-conformance-check.md | 2 | one-source, one-source |
| 02-derive-the-ticket-check-floor.md | 2 | one-source, one-source |
| 03-grade-the-plan-against-the-floor.md | 1 | spec-row |
| 04-refuse-a-planned-cell-at-the-freeze.md | 2 | delegate-error, delegate-error |
| 05-state-the-floor-in-the-guidance.md | 0 | none |

## Agent-experience improvements

### Bench CLI

- Make `bench record chunk` refuse a plan digest that differs from the record digest, and name `bench record amendment` in the refusal.
  Feeds: FT318
- Make `bench probe` report the mutated test run exit code, name an unused import in an invalid verdict, and lock one subject file.
  Feeds: FT168
- Add one `bench record ticket` verb that records the amendment, freezes the chunk, and records each verification from an excerpt folder.
  Feeds: new
- Make `bench commit` rebuild the worktree binary when the commit changes a Go build input.
  Feeds: new
- Print the whole-call wall time in `bench test` and `bench probe`, and keep failures rows ahead of a spill cut; the census entry `ft349-build census: 15 raw calls` holds the detail.
  Feeds: new

### Skills

- Make `craft-delegate` require each author charge to name the existing owners, seams, and fixture helpers that the ticket composes.
  Feeds: new
- Make `craft-spec` confirm that each owner a spec names is exported, and list the fixtures that pin a row list.
  Feeds: new
- Make `craft-line` name Fable/medium as the review line for a ticketed build after one held-out comparison confirms it.
  Feeds: new

### Process

- Give the chunk freeze and the author verification records to the ticket author in the implement-spec guidance.
  Feeds: new
- Apply a plan-text correction before the chunk freeze, or freeze the chunk again before any record at the old source.
  Feeds: new
- Widen a fence before the freeze when an active author reports a copy outside it, so the author folds the copy before review.
  Feeds: new