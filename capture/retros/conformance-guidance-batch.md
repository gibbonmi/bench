# Conformance and guidance batch

## Outcome

The batch starts at 14:24:06 UTC and has a 15:24:06 UTC deadline on 2026-09-22.
FT326 lands at c8f33c9d979ab0260e16fa46540dd2bc6cadce52.
FT297 lands at 34f16d1cd29a990e04778bf6f4e4985cf5c4d95a.
Each passes the serialized full gate and independent Standards, Spec, and Coverage review.
The sentence-boundary proposal remains excluded.

FT219 and FT102 have committed implementations and three passing native review axes.
Their six fresh-reader scenarios pass, but their landing gates remain pending.
The batch reserves its last gate slot for this capture commit.
No source or focused test runner changes a tree during its full gate.

FT322 at 9339641557369f2afc20ccb545314cee6995ba0c and FT323 at 333bc7ce38821638c5d352a65ac896d2642f943e are existing ancestors.
Their roadmap rows are retired without rebuilding either item.
Their repair counts remain zero of two and one of two.

FT290 is not admitted because its test-determinism prerequisite has only staged documentation.
The prerequisite's collision note claims a named-check owner move that current source does not contain.
Correct that note, implement its eight tickets and four checkpoints, and land it before refreshing FT290.
FT108 needs the FT89 guidance owner first; FT299 lacks a decided executable rehearsal route.
Decision rows and missing-reproduction items remain pending.

## Gate-stage timings

| surface | gofmt ms | vet ms | test ms | race ms | system ms | shellcheck |
| --- | --- | --- | --- | --- | --- | --- |
| FT326 landing | 148 | 1360 | 159938 | 3743 | 44843 | skipped, 35 ms |
| FT297 landing | 132 | 1138 | 166103 | 7079 | 44401 | skipped, 31 ms |
| Guidance review composition | 126 | 1175 | 171655 | 4180 | 46511 | skipped, 61 ms |

The FT297 timing record is the landing log at `.logs/gate-20260922T151346.150780016Z-3279304.jsonl`.
The final capture gate writes its own log after this record commits.
Each completed gate reports eight capability skips and no environment skips.
The skips comprise four FIFO cases and four privilege cases.
The guidance composition gate is review evidence, not landing evidence.

## Ticket-versus-spec-slice and delegate performance

| ticket | scope | result | retained author | review |
| --- | --- | --- | --- | --- |
| FT326 | distinctness only | landed, one repair | Astra/ultra, ft326_author | three native Sol/high passes |
| FT297 | path and Go flag contracts | landed, one repair | Astra/ultra, ft297_author | three native Sol/high passes |
| FT219 | existing ready-map refresh | reviewed, landing pending | Astra/ultra, coordinator | three native Sol/high passes |
| FT102 | synthesis owner checks | reviewed, landing pending | Astra/ultra, coordinator | three native Sol/high passes |

| surface | claim | status | confidence | label | model / effort / role |
| --- | --- | --- | --- | --- | --- |
| FT326 | one expected repair | verified | 7 | 1 | Astra / ultra / implementer |
| FT297 | one expected repair | verified | 7 | 1 | Astra / ultra / implementer |
| FT219 | zero expected repairs | claimed | 8 | pending | Astra / ultra / implementer |
| FT102 | zero expected repairs | claimed | 8 | pending | Astra / ultra / implementer |
| Fresh adoption | six scenarios match the owner text | claimed | unknown | pending | Terra / low / reader |

The batch repair-forecast Brier mean is 0.09 over two labeled pairs, with zero abstentions.
The three native review axes have no calibrated outcome labels; their Brier mean remains unknown.
The fresh reader reports qualitative confidence only, so its numeric calibration remains unknown.
Run tokens, costs, and comparative fast-mode latency remain unknown.
The normalized assessment is `batch-conformance-20260922-142406` in the local assessment store.

## Coordinator catches

The current source proves that FT290's prerequisite has not landed.
The readiness pass reduces FT102 to its existing synthesis owner and avoids a second policy copy.
The FT326 census uses 14 current registry slices instead of the stale roadmap count.
A failed native spawn is retained as orchestration evidence; distinct review sessions are reused by axis.

The primary checkout has five unrelated paths preserved in a verified external recovery copy.
Automatic approval review first rejects temporary clearing under the preservation instruction.
The user then approves temporary clearing and exact restoration.
Restoration and the final hash check follow the last landing and are recorded in the local handoff.

## Repair attribution

| ticket | rounds | causes |
| --- | --- | --- |
| FT326 | 1 of 2 | delegate: expectation proof, repeated diagnosis, cross-family probe, Unicode boundary coverage |
| FT297 | 1 of 2 | delegate: repeated import parser and missing special-file refusal |
| FT219 | 0 of 2 | none; landing pending |
| FT102 | 0 of 2 | none; landing pending |

FT326 preserves one refuted request for all kind pairs because exact-kind identity already owns that rule.
Its first needle-omission mutation fails compilation and supplies no passing evidence.
The corrected mutation compiles, turns the named tests red, and is restored.
FT297's four FIFO refusal probes turn red under the classifier bypass and are restored.
No repair or hardening allowance resets when the budget restarts.
The unrelated jev repair assignment remains at five consumed cycles, with no sixth cycle authorized.

## Agent-experience improvements

### Bench CLI

- Make the price of a review-venue merge explicit because this composition runs a full gate.
  Feeds: FT314

### Skills

- Keep source identity and exact terminal excerpts in each native review record.
  Feeds: FT318

### Process

- Reserve a full gate for capture before accepting another landing into a bounded batch.
  Feeds: none
