## Outcome

The build cut the FT411 review tests that a spec-stage probe proved subsumed, and each cut kept a biting probe on its kept test.
Chunk FT411-C1 removed five adoption tests, merged three boundary tests into one, and cut seven preview rows.
Chunk FT411-C2 removed eight worktree tests, merged three census landings into one, and cut two resume rows.
The worktree count pin went from 880 to 872, and no production file changed.
Landing `d96982a30` published `Status: implemented` on a green gate.

## Cost

- landing: d96982a30791d29a93975acadf870e9d36a9ece8, range d96982a30791d29a93975acadf870e9d36a9ece8^1..d96982a30791d29a93975acadf870e9d36a9ece8^2
- wall: 68 min
- commits: 26 total, 19 bookkeeping, 7 change
- graph: 6 tickets, chain 2, 2 chunks, 1 serial checkpoints
- dispatch: 6 authors, 1 repairs
- trace: commit 26 spans in 241256 ms, lane 0 spans in 0 ms, landing 1 spans in 338845 ms
- wall per ticket: 11.5 min against 60 min
- bookkeeping per ticket: 3.2 against 7.1

Provider usage and charges are unknown.
The build did not time the cut packages before and after the cuts, so the gate saving is not measured.

## Gate-stage timings

- landing: commit d96982a30791d29a93975acadf870e9d36a9ece8, trace 2d121bbe0c41d4e6f9da86b948d4fbf9
- gofmt: 189 ms
- vet: 1627 ms
- test: 230476 ms
- race: 7072 ms
- system: 52518 ms
- shellcheck: 533 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 01 | TC2 to TC7 and TC43 bit on the kept tests | verified | 9 | held | opus / medium / implementer |
| ticket 01 | TC8: the merged test keeps the process check | claimed | 8 | held | opus / medium / implementer |
| ticket 01 | TC9: the LTE cells cite the kept tests | claimed | 8 | held | opus / medium / implementer |
| ticket 02 | TC1 and TC10 to TC19 bit on the kept tests | verified | 9 | held | opus / medium / implementer |
| ticket 02 | TC45 bit on the decision table | verified | 8 | held | opus / medium / implementer |
| ticket 03 | TC20 to TC24, TC41, and TC42 bit | verified | 9 | held | opus / medium / implementer |
| ticket 03 | TC44: the four kept rows remain | claimed | 8 | held | opus / medium / implementer |
| ticket 04 | TC25, TC26, TC29 to TC33, and TC35 bit | verified | 9 | held | opus / high / implementer |
| ticket 05 | TC27, TC28, and TC34 bit | verified | 9 | held | opus / high / implementer |
| ticket 06 | TC35 to TC40 bit | verified | 9 | held | opus / high / implementer |
| repair r02 | both C1 findings are fixed | claimed | unknown | abstained | opus / medium / repair |
| orchestrator | one repair round for the build | claimed | 6 | held | opus / medium / orchestrator |
| C1 Standards | the nil check is loosened | claimed | 6 | held | sonnet / high / reviewer |
| C1 Standards | the SD4 comment is stale | claimed | 5 | held | sonnet / high / reviewer |
| C1 Spec | the nil check is loosened | claimed | 4 | held | sonnet / high / reviewer |
| C1 Coverage | the nil check is loosened | claimed | 4 | held | sonnet / high / reviewer |
| C1 Coverage | the TC3 why-clause overstates the kept test | claimed | 3 | refuted | sonnet / high / reviewer |
| C2 Standards | the face-list comment holds provenance | claimed | 4 | refuted | sonnet / high / reviewer |
| C2 Standards | the fixture loop needs a guard clause | claimed | 5 | refuted | sonnet / high / reviewer |
| C2 Standards | the freshness doc comment omits two IDs | claimed | 4 | refuted | sonnet / high / reviewer |
| C2 Standards | a reflowed comment wraps short | claimed | 3 | refuted | sonnet / high / reviewer |
| C2 Coverage | the slug-and-path test name overstates its row | claimed | 6 | refuted | sonnet / high / reviewer |
| C2 Coverage | the skipped PL30 cell stays declared | claimed | 5 | refuted | sonnet / high / reviewer |

Brier mean: 0.130 over 22 pairs; 1 abstention.

Each author kept to its planned ticket and chunk boundaries, and no author needed a fence widen.
Every plan probe bit with a byte-exact restore, and no probe needed a replacement.

## Coordinator catches

- The ticket 02 author loosened a nil check so that a moved success row fit a refusal table, and flagged it. All three axes found it, and one repair restored the exact check.
- The orchestrator examined an 86-line spec deletion in an assignment commit. `bench record assignment` rewrites the plan fence as one compact JSON line, so the content stayed intact.
- The review preflight refused the frozen chunk tip as the source tip, because it derives the source tip from the record commit.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-cut-subsumed-repair-tests.md | 0 | none |
| 02-cut-subsumed-adopt-tests.md | 1 | spec-row |
| 03-cut-preview-rows-and-guard-kept-tests.md | 0 | none |
| 04-cut-census-effects-and-freshness-tests.md | 0 | none |
| 05-cut-resume-rows.md | 0 | none |
| 06-cut-refusal-route-tests.md | 0 | none |

The ticket 02 round is `spec-row`, because the ticket put a success row in a refusal table that asserts nil slices.

## Agent-experience improvements

### Bench CLI

- Let one `bench record excerpt` call stage every row of a mutant-set response, because five authors made one call for each row.
  Feeds: new
- Document the mutant-set JSON schema in `bench probe --help`, because three authors listed the scratchpad to learn it (census `ls` 3, `cat` 2, owned gap).
  Feeds: new
- State in `bench probe --help` that one `--mutants` call covers one subject file, because the ticket 01 author learned it from a refusal.
  Feeds: new
- Let `bench record ticket` carry a requirement whose test file the repair delta does not change, because the C1 repair reran three unreached probes.
  Feeds: new
- Add a `bench record amendment` form that writes the version 2 execution block, because the orchestrator edited the plan fence with a script (census `python3` 1, owned gap).
  Feeds: new
- Census account: the landing counted 15 raw calls. The 7 `rg` calls are sanctioned reads, and the 2 `cd` calls are unknown, so these need no change.
  Feeds: none

### Skills

- Name the probe grouping by subject file in the author charge template of `/bench-implement-spec`, so each author plans one probe call per subject file.
  Feeds: new
- Say in `/bench-review-implementation` step 1 that the review preflight takes the record commit as `--source-tip`, while the axes read the frozen chunk tip.
  Feeds: new

### Process

- Put a moved success row in its own assertion, not in a refusal table, when a ticket moves rows into a kept table.
  Feeds: none
- Time the cut test packages at the base and at the landed tip on a quiet host, because this build cannot report its gate saving.
  Feeds: new