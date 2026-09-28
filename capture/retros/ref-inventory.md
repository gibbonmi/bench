## Outcome

FT199 landed as `ab886f5b` on 2026-09-27, 13 hours after the spec staged at `72a749a3`. The build ran six tickets in four chunks with one fresh Opus/high author per ticket. `bench worktree clean --discard-branch --unclaimed` now classifies every unclaimed ref as landed, subsumed, or unique, and the bulk sweep removes only the first two classes. A unique ref discards by target with a dated discarded ref written first, and `bench spec retire` lists the superseded candidates and the unique count. The coverage map holds 100 rows; 20 of them joined during the build by reviewer decision or plan expansion. The review record holds 49 findings over 12 rounds, and the landing gate is green.

## Gate-stage timings

- landing: commit ab886f5b5ee04867110e79e4e1879a920a50434f, trace c441794301c0104e1c097b41bd0f8e0b
- gofmt: 154 ms
- vet: 1237 ms
- test: 133198 ms
- race: 3022 ms
- system: 43466 ms
- shellcheck: 548 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| RI-C1a author | RI57 glossary row delivered without a test red | claimed | 9 | held | claude-opus-5-5 / high / implementer |
| RI-C1a author | RI82 glossary row delivered without a test red | claimed | 9 | held | claude-opus-5-5 / high / implementer |
| RI-C1a Standards | R1 test expectation without a recorded red | finding | 7 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a Standards | R2 comment register | finding | 5 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a Standards | R3 fixture harness duplicated | finding | 5 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a Standards | R4 duplicated owner-letter rule | finding | 4 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a Standards | R5 judgment call on naming | finding | 3 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a Standards | R6 no-op | finding | 3 | abstained | claude-opus-5-5 / high / reviewer |
| RI-C1a Spec | R7 chunk table omission | finding | 8 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a Spec | R8 cost sentence | finding | 4 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a Coverage | R9 symref in a Bench namespace deleted by the sweep | finding | 9 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a Coverage | R10 content-landed recorded holder | finding | 9 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a Coverage | R11 class-only stale fingerprint | finding | 8 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a Coverage | R12 record states as holders | finding | 8 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a round 2 | R13 landed-holder reasons | finding | 7 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a round 2 | R14 owner-letter helper | finding | 4 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a round 2 | R15 comment wording | finding | 4 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a round 2 | R16 no-op | finding | 3 | abstained | claude-opus-5-5 / high / reviewer |
| RI-C1a round 2 | R17 dangling symref decision | finding | 5 | held | claude-opus-5-5 / high / reviewer |
| RI-C1a round 2 | R18 chunk table cell | finding | 4 | held | claude-opus-5-5 / high / reviewer |
| RI-C1b Standards | R19 status route shape | finding | 6 | held | claude-opus-5-5 / high / reviewer |
| RI-C1b Standards | R20 plan command spelling duplicated | finding | 6 | held | claude-opus-5-5 / high / reviewer |
| RI-C1b Standards | R21 faulted detail | finding | 7 | held | claude-opus-5-5 / high / reviewer |
| RI-C1b Standards | R22 comment register | finding | 5 | held | claude-opus-5-5 / high / reviewer |
| RI-C1b Standards | R23 judgment call | finding | 4 | held | claude-opus-5-5 / high / reviewer |
| RI-C1b Coverage | R24 symref counted twice | finding | 7 | held | claude-opus-5-5 / high / reviewer |
| RI-C1b Coverage | R25 planner failure route | finding | 6 | held | claude-opus-5-5 / high / reviewer |
| RI-C1b Coverage | R26 blob-tip Git-state fallback | finding | 8 | held | claude-opus-5-5 / high / reviewer |
| RI-C2a Standards | R27 guard comment | finding | 7 | held | claude-opus-5-5 / high / reviewer |
| RI-C2a Standards | R28 moved-ref harness duplicated | finding | 5 | held | claude-opus-5-5 / high / reviewer |
| RI-C2a Standards | R29 comment does not parse | finding | 5 | held | claude-opus-5-5 / high / reviewer |
| RI-C2a Standards | R30 fixture re-derives the path | finding | 3 | held | claude-opus-5-5 / high / reviewer |
| RI-C2a Standards | R31 no-op | finding | 4 | abstained | claude-opus-5-5 / high / reviewer |
| RI-C2a Coverage | R32 sweep follows a discarded symref | finding | 9 | held | claude-opus-5-5 / high / reviewer |
| RI-C2a Coverage | R33 suffixed date segment | finding | 7 | held | claude-opus-5-5 / high / reviewer |
| RI-C2a Coverage | R34 UTC date for a local instant | finding | 6 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b Standards | R35 stale unique detail trimmed twice | finding | 9 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b Standards | R36 two selector derivations | finding | 6 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b Standards | R37 second argv parse | finding | 6 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b Standards | R38 expectations without recorded reds | finding | 5 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b Standards | R39 no-op | finding | 4 | abstained | claude-opus-5-5 / high / reviewer |
| RI-C2b Spec | R40 trailing `--` hides the listing | finding | 6 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b Coverage | R41 symref at the planned path | finding | 8 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b Coverage | R42 refusal guard uncovered | finding | 8 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b Coverage | R43 fault guard untested | finding | 7 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b Coverage | R44 requalify never observed | finding | 6 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b Coverage | R45 read-write window | finding | 5 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b round 2 | R46 missing learning entry | finding | 7 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b round 2 | R47 step count in the seam notes | finding | 6 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b round 2 | R48 symref planted in the write window | finding | 8 | held | claude-opus-5-5 / high / reviewer |
| RI-C2b round 3 | R49 symref planted before the delete | finding | 8 | held | claude-opus-5-5 / high / reviewer |

Brier mean 0.163 over 47 labeled pairs, with 4 abstentions. Every retained finding held through its confirming round, so the mean measures under-confidence, not error.

## Coordinator catches

- The RI-C1a author reported two anchor rows as verified with no test red; the record labels them claimed.
- The ticket 5 author stopped on two fence gaps that the plan did not name, and a plan commit expanded the fence before it resumed.
- The RI-C2b round 2 Spec axis found that a fence expansion had no `bench learning` entry; the orchestrator had recorded it only in prose.
- The ticket 4 repair 2 session stopped on a material acceptance shortfall. Git 2.43 replaces a dangling symref under every write form, so the RI102 wording could not hold. The orchestrator confirmed the cause with a throwaway loop and a scratch repository before the reviewer restated the row.
- The Codex Astra consultation on that decision produced no output in over three hours; the orchestrator stopped it and the reviewer decided in conversation.
- Every accepted repair took an independent coordinator probe of a different kind and site, and each bit.
- The first landing refused because the `main` composition changed two guidance files outside every fence; the composition then took its own review round.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-classify-unclaimed-refs.md | 2 | spec-row, one-source |
| 2-route-status-to-the-plan.md | 2 | spec-row, spec-row |
| 3-sweep-discarded-refs.md | 1 | spec-row |
| 4-discard-a-unique-ref-by-target.md | 3 | one-source, spec-row, spec-row |
| 5-list-retire-candidates.md | 1 | one-source |
| 6-repair-glossary-shift-namespace.md | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Give the review record verbs a digest command for a commit and an amendment appender. The landing census counted 136 raw calls: rg 93, wc 18, python3 5, sed 5, gofmt 3, ls 3, and seven single calls. The two verbs remove the python3 heads and most of the rg heads.
  Feeds: new
- Accept the `sha256:` prefix in the evidence cursor, because two authors lost a call to `invalid-cursor` before they dropped it.
  Feeds: none
- Give the review preflight charge a form for a Markdown-only landing composition that needs no fixture closure. The composition round ran without an evidence id.
  Feeds: new
- Name the `binary-seal` rebuild in the build preflight's own next cell after a Go change, as one author noted.
  Feeds: none

### Skills

- State in `craft-delegate` that a window with a write and a delete gets a planted symref at both refs in the first Coverage round. RI102 and RI103 arrived one round apart.
  Feeds: FT199
- State in `craft-line` that a consultation with no output past its expected window is stopped and rerouted, not awaited.
  Feeds: none

### Process

- Cite the last chunk's seams and run the `main` composition before the final reconciliation. Each late plan commit after the chunk tip forced author reruns and a record move.
  Feeds: none
- Record every fence or plan expansion with `bench learning` in the same step as its commit, before the dispatch it enables.
  Feeds: none
- Keep the RI-C2b symref race findings as one hostile-input family in the profile. A future transaction spec then plants a symref at every ref it touches from its first round.
  Feeds: FT199