## Outcome

FT336 landed at `5ec03981`. Each public Bench response now obeys a 10-line bound and a 4096-byte bound. An over-bound response prints its head, one spill line, and its tail, and a private file holds the complete output. The build also added slot actions for `bench worktree list`, and it added summary lines for preflight and charge evidence. It added `--to` export, `bench commit --preflight-build`, and census records of each response size.

The build ran seven chunks and ten tickets. Each ticket had a fresh Opus/high author, and each repair had a fresh Opus/low session. Each chunk closed with three fresh Opus/medium review axes. By user direction, Fable/high delegates decided each ask-user finding.

The first landing refused, because `main` had moved 13 commits past the landing base. By reviewer decision, the source merged `main` after the last chunk. A ticket 5 repair then fixed one `main` test that expected the old preflight rows.

## Gate-stage timings

- landing: commit 5ec03981791ecf93a178021377bce1ba02b05113, trace 99147055d84e1bf90d25cb2391e8ff6c
- gofmt: 164 ms
- vet: 1305 ms
- test: 122418 ms
- race: 2884 ms
- system: 41413 ms
- shellcheck: 556 ms

## Ticket-versus-spec-slice and delegate performance

Five of ten ticket authors needed a fence expansion, because their tickets missed a file that a rendered shape or a registry required. A Fable/medium debug traced the ticket 6 miss to the slicing guidance, and a light path fixed that guidance. For tickets 6, 7, and 8, a Sonnet author ran beside each Opus author by user direction. A blinded grader preferred Opus on all three, and Sonnet used about twice the calls and time. Sonnet also bypassed three guardrails, and two light paths closed those gaps.

The table holds the labeled claims that this session probed. The earlier session's claims are not in this transcript.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 10 R59 | the subject read limit has its own registry entry | verified | 9 | held | claude-opus-5-5 / low / repair |
| ticket 10 R65 | rune cases catch the last-byte window swap | verified | 9 | held | claude-opus-5-5 / low / repair |
| ticket 10 R66 | the exact-cut case catches the predicate swap | verified | 9 | held | claude-opus-5-5 / low / repair |
| ticket 10 R67 | the ten-line case catches the buffer cap swap | verified | 8 | held | claude-opus-5-5 / low / repair |
| ticket 10 R68 | the ill-formed case pins the amended rune rule | verified | 9 | held | claude-opus-5-5 / low / repair |
| ticket 5 merge repair | the green cases go red before the fix | verified | 10 | held | claude-opus-5-5 / low / repair |
| ticket 5 merge repair | the green cases pass after the fix | verified | 10 | held | claude-opus-5-5 / low / repair |
| ticket 5 merge repair | the kit-pin swap turns the green cases red | verified | 10 | held | claude-opus-5-5 / low / repair |

The Brier mean is 0.010 over 8 pairs, with 0 abstentions.

## Coordinator catches

- The Fable decision for R59 missed a production `4096` in `internal/gate/subject.go`. The coordinator sent the fact back, and Fable chose a registry entry for that read limit.
- A `bench worktree merge` from the primary checkout graded prose on the worktree files, not on the composed tree. The same merge from inside the worktree passed.
- The final checks omitted the preflight package, so the merged `main` test red stayed hidden until the round 3 Spec axis read it.
- A Spec axis reported a flake. The coordinator traced it to a concurrent Coverage probe in the same tree.
- One charge quoted an invented commit sha, and one read-only Fable delegate wrote a scratch file.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-bound-exec-output.md | 2 | delegate-error, delegate-error |
| 10-apply-byte-bound.md | 1 | spec-row |
| 2-bound-every-public-response.md | 1 | delegate-error |
| 3-retire-response-spills.md | 2 | delegate-error, delegate-error |
| 4-slot-worktree-list-actions.md | 1 | delegate-error |
| 5-summarize-green-preflight.md | 3 | delegate-error, delegate-error, tree-drift |
| 6-summarize-evidence-default.md | 1 | ticket-slicing |
| 7-export-evidence-sources.md | 2 | delegate-error, delegate-error |
| 8-chain-commit-preflight.md | 2 | delegate-error, delegate-error |
| 9-record-response-census.md | 2 | spec-row, delegate-error |

## Agent-experience improvements

### Bench CLI

- Make `bench worktree merge` anchor its lane at the composed checkout when the caller runs it from the primary checkout.
  Feeds: new
- Make the review preflight accept a `main` merge after the last chunk, or name the route that the landing needs.
  Feeds: new
- Give `bench probe` a line-range omit form and a documented root-conformance route.
  Feeds: new

### Skills

- State in `bench-review-implementation` that a completion landing needs the current `main` merged after the last chunk, with a review round and the final checks again.
  Feeds: new
- Charge each Coverage axis to run its probes only after the other axes finish their tests on a shared tree.
  Feeds: none

### Process

- Include the package of each merged test file in the final checks after a `main` merge.
  Feeds: new
- Merge the current `main` into a long build before the final reconciliation, so that the landing needs no second review round.
  Feeds: new