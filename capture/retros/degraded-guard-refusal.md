## Outcome

The FT366 degraded-guard build landed as `bd7beca66989d7e56a4b515775f67b8682c44cb7` with a green whole-project gate. The degraded Git hook now refuses an unclosed quote, a trailing escape, and a malformed wrapper child with a parse-repair diagnostic. The core's unlexed refusal now carries only its repair advice. All 42 acceptance rows reconciled, including row GR42, which the review added.

The build ran one ticket in one chunk on Claude Code as a `--full --auto-approve` run. The reviewer directed the Opus/high author line, because the spec declared a Codex model. The chunk used 2 repair cycles: one test-only repair session and one orchestrator spec correction.

## Gate-stage timings

- landing: commit bd7beca66989d7e56a4b515775f67b8682c44cb7, trace 8d1d7955f7e343af9a4f8a1779ada6f1
- gofmt: 170 ms
- vet: 1411 ms
- test: 164480 ms
- race: 2602 ms
- system: 58045 ms
- shellcheck: 525 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 01 | GR01-GR06, GR10-GR12, GR26, GR27, GR30, GR32, GR33 red then green | verified | 9 | held | opus / high / implementer |
| ticket 01 | GR07-GR09, GR17, GR18, GR29 red then green | verified | 8 | held | opus / high / implementer |
| ticket 01 | GR13-GR16, GR19-GR25, GR28, GR31, GR34-GR41 green with preservation reds | verified | 9 | held | opus / high / implementer |
| repair 1 | S1, S2, S3, P2, and C1 closed | verified | 9 | held | opus / high / repair |
| repair 1 | S4 closed | verified | 8 | held | opus / high / repair |

Brier mean 0.022 over 5 pairs, 0 abstained.

The single ticket matched its chunk slice, and the author met every original row in one pass. The review still found seven repair targets. Four of them were one-source defects, because the narrowed `Writes:` line kept the shared system-test helpers out of reach.

## Coordinator catches

- The staleness pass found ten fence paths that only the deleted ownership closure had pinned. It also found the `check-floor` red and five drifted citations.
- The coordinator's independent probe swapped the parse-failure return at the `invokes_git` entry, and GR04 turned red.
- `bench record confirmation` refused the test-only path, because the repair range also changed the spec. The build then needed a full confirming round.
- The coordinator recorded the Coverage advice C3 as a finding ID, and the issuing axis had to reaffirm it as advice.
- The coordinator's own spec expansion left a stale path count (P3). Its correction moved the chunk tip and made every GR1 verification and review result stale.
- The coordinator reran the GR1 verifications under its own session, but the checkpoint owes them to the latest ticket assignment, so the repair session reran them.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-refuse-malformed-guard-commands.md | 2 | one-source, other |

## Agent-experience improvements

### Bench CLI

- Census account: the landing counted 18 raw calls, with rg 16 (Bench form: none), cp 1 (none), and python3 1 (none); add a `bench spec` form that amends a version 1 plan to a version 2 delegate plan, which removes the python3 call.
  Feeds: new
- Make `bench record ticket` skip an `--id-prefix` that the requirement id already carries, so that it writes no `t1-t1-` identifiers.
  Feeds: new
- Make `bench consumers` accept the full import-path symbol form, or name the accepted form in its refusal.
  Feeds: new
- Let `bench record confirmation` accept a repair range whose only non-test change is the coverage-map seam cell of a row that the repair adds.
  Feeds: new

### Skills

- Make `craft-line` resolve a spec `Line:` through the tier binding of each adopted harness, so that a build in another harness needs no author-line stop.
  Feeds: new
- Make the staleness pass charge name each rule removal that landed since the reviewed base, so that its slices search for paths that the removed rule shaped.
  Feeds: new

### Process

- Read the provider scorecard before the line declaration, because this build repeated three lessons that its current decisions already state.
  Feeds: none
- Batch every spec correction into the repair's enabling plan commit, so that no correction follows the chunk freeze.
  Feeds: none
- Sweep the staged specs for closure-derived fence paths once, instead of one build-entry staleness pass at a time.
  Feeds: new