## Outcome

The build landed at 35d5a8fa583aef00cb6278182d3a2b8f6ae0cbac from source dcbbfb80..474b64ea on destination base 19b84bbe, and the whole-project gate was green.
It adds a landing lease, an open-build window with chunk gaps, a recorded fold of `main`, a terminal-bound landing override, and the moved-destination face.
All 100 acceptance rows passed, and the complete checkpoint was green at 474b64ea.
Fold F1 merged the FT354-durable-replace landing after BW-E with no resolution path.
The installed binary could not read the new fold record, so the build landed on its own candidate binary by reviewer decision on 2026-10-12.
Under `--auto-approve`, the back-out route changed to `git revert -m 1 --no-commit`, and that change waits for the reviewer's veto.

## Gate-stage timings

- landing: commit 35d5a8fa583aef00cb6278182d3a2b8f6ae0cbac, trace 8aa0f60c3e20198ecaca5e12757bdf5b
- gofmt: 164 ms
- vet: 1601 ms
- test: 186097 ms
- race: 2573 ms
- system: 65977 ms
- shellcheck: 619 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 7 author return | No test touches the real ledger | held | 8 | review record and ledger keys after BW-E | opus / high / implementation |
| ticket 8 author return | Each sentence states the enforced route | failed | 7 | BW-F Spec P1 and the ADR 0014 wait | opus / high / implementation |
| ticket 7 repair return | No shared-ledger write | held | 9 | ledger keys after the BW-E checkpoint | opus / high / repair |

Brier mean: 0.18 over 3 pairs. Abstentions: 0.

Each of the eight tickets went to one fresh Opus/high author, and every author kept its writes inside the ticket fence or stopped at it.
Every chunk took at least one repair cycle, and BW-C and BW-F took both cycles.
Most first-round findings were one-source defects: a copied harness, a second derivation of a predicate, or a hand-spelled value.
The ticket 8 author wrote true guidance for every enforced route except the back-out command, which the spec itself had wrong.
Each repair session ran `/bench-debug` and returned its phases after the reviewer's direction.

## Coordinator catches

- The BW-D checkpoint on the worktree binary wrote a `gaps` record that the installed binary refused. The coordinator stopped, and the reviewer approved a hand removal of that key after a backup.
- The installed binary refused the landing on the new `folds` field. A repro through the installed `bench gate --checkpoint --complete` confirmed the cause before the candidate landing.
- `BENCH_KIT` set to the primary checkout dropped the system phase from the coverage census. A one-variable repro chose the landing environment.
- The Spec axis found that the spec's back-out command fails on a two-parent landing, and the coordinator corrected the spec under `--auto-approve`.
- Repair 1 of BW-F reported false wording outside its fence, and the coordinator expanded the fence once instead of a third review round.
- The check-floor row caught every plan expansion that omitted a package, before any dispatch.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-hold-the-landing-lease.md | 1 | one-source |
| 2-route-a-moved-destination.md | 1 | one-source |
| 3-record-a-fold.md | 2 | one-source, one-source |
| 4-review-a-fold-resolution.md | 1 | one-source |
| 5-hold-main-for-an-open-build.md | 1 | one-source |
| 6-open-and-close-chunk-gaps.md | 1 | spec-row |
| 7-grant-a-reviewer-override.md | 1 | delegate-error |
| 8-document-the-landing-window.md | 2 | one-source, ticket-slicing |

## Agent-experience improvements

### Bench CLI

- Bind the landing's gate subject through a binary that reads the landed record format, so that the installed shim lands a new record field.
  Feeds: new
- Make the coverage census select the kit-source phases for any kit-source tree, so that a system-test citation passes under any `BENCH_KIT`.
  Feeds: new
- Census: 60 raw calls, recorded in the learning entry "FT342-FT253 census: 60 raw calls".
  Feeds: new

### Skills

- Make the repair-charge template require a full read of `bench-debug` and a returned Debug phases section.
  Feeds: new
- Make a glossary ticket sweep code comments and help text for the retired wording, not only the guidance files.
  Feeds: new

### Process

- Name the pre-landing writer and the landing binary in the edge inventory of any spec that adds a ledger or review-record field.
  Feeds: new
- Try each guidance command that acts on a landing commit on a two-parent fixture before the spec stage closes.
  Feeds: new