## Outcome

The preflight-spec-staleness build (FT375) landed at 5371131b256d701b85df54de62361541de9d4c74 from the reviewed pair 81b10e74..fcb5de85, and the landing gate was green. `bench preflight build <slug> --audit-spec --base --source-tip` now reports the exact drift, the cited-symbol and cited-test states, each unsupported claim, and one staleness decision. `/bench-implement-spec` omits the staleness pass only when that decision states `complete` true and `manual_required` false.

The build ran five tickets in three chunks as a `--full --auto-approve` run with an author limit of 1. The chunks used 3, 1, and 2 repair cycles. SS-C1 and SS-C3 went past one cycle under `--auto-approve`. Two fresh Claude subagents ran the native rows on approved fixtures, before and after each guidance repair.

## Gate-stage timings

- landing: commit 5371131b256d701b85df54de62361541de9d4c74, trace 585b757414d2836756db2b09dd07e73c
- gofmt: 194 ms
- vet: 1776 ms
- test: 238168 ms
- race: 9306 ms
- system: 70420 ms
- shellcheck: 649 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 01 author | the audit source drift met SS1 to SS20 inside its fence | held after three repairs | unknown | abstain | claude-opus-5-5 / high / implementer |
| tickets 02 to 05 authors | each ticket met its rows inside its fence | held, with one repair each for tickets 02 to 04 and two for ticket 05 | unknown | abstain | claude-opus-5-5 / medium / implementer |
| repair sessions | each repair closed every accepted finding of its round | held; SS-C1 needed a third repair for one new Coverage finding | unknown | abstain | claude-opus-5-5 / high, medium / repair |
| native sessions | the clean fixture omits the pass and the diagnostic fixture charges it | held in both runs | unknown | abstain | claude-opus-5-5 / medium / verification |

No author or repair return carried a claims table with integer confidences, so the build holds no labeled pair. Brier mean: unknown over 0 pairs; 4 abstentions.

## Coordinator catches

- The entry preflight was red on `writes-resolve`, because landed work had deleted 31 canaries. The staleness pass took the red into the enabling plan commit.
- The orchestrator's needle-relocation script matched needles by substring and moved 29 SS-C1 seam locations to wrong rows. SS-C1 review round 1 caught it as P1, and a plan commit pointed each location back.
- Coordinator probes found three silent mutants: the `Complete` term, the headroom projection, and the red-over-diagnostic precedence. Each probe added a witness in a repair.
- Two focused entries for tickets 02 and 04 recorded a probe other than the named plan probe. Both were recorded again with the named probe.
- `bench record ticket` could not stage a native excerpt, so the native rows and each ticket entry were recorded one at a time.
- The orchestrator refroze SS-C2 at record commits. SS-C3 then opened on the last source commit, and the completion entry refused a review chain gap. A refreeze on the recorded SS-C2 tip cost one more full review round.
- The orchestrator first recorded the SS-C3 repair verifications under its own identity. The completion entry refused them, and the repair session recorded them again as the owed author.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-audit-source-drift.md | 3 | one-source, delegate-error, delegate-error |
| 02-audit-symbol-identities.md | 1 | one-source |
| 03-audit-test-citations.md | 1 | delegate-error |
| 04-audit-coverage-and-completeness.md | 1 | one-source |
| 05-consume-exact-manual-omission.md | 2 | one-source, one-source |

## Agent-experience improvements

### Bench CLI

- The landing census counted 27 raw calls: rg 14, python3 6, for 2, ls 2, and one each of awk, cd, and cp. Their Bench forms are the Grep tool for rg, none for python3, for, awk, and cp, `bench worktree path` for ls, and `bench worktree exec` for cd; add a `bench plan show <slug>` read of the completion plan, because each python3 call parsed that JSON line.
  Feeds: new
- Let `bench record ticket` take an excerpt file for a native requirement, so one call records a chunk that holds native rows.
  Feeds: new
- Make `bench record verification` refuse a performer that is not the owed author of the plan, so the refusal comes at record time and not at the completion entry.
  Feeds: new
- Make `bench record chunk` refuse a tip that only adds record commits after the last source commit of the chunk.
  Feeds: new
- Make `bench gate --checkpoint` for a chunk check the chain to its predecessor, so a chain gap shows before the next review round.
  Feeds: new

### Skills

- State in the implement-spec phase that a chunk tip is its last source commit, and that the successor base descends from the recorded predecessor tip.
  Feeds: new
- State in the repair charge that the repair session records its own ticket verifications with `bench record verification` under its session identity.
  Feeds: new

### Process

- Make each relocation or bulk-edit script of the orchestrator match whole lines and print its target list before it applies.
  Feeds: none
- Treat a missing-file check as the only proof of absence in the binary-seal collector, because the build made that row exempt from the omission conjunction.
  Feeds: new
- Require a claims table with integer confidences in every author and repair return, so a build holds labeled pairs.
  Feeds: none