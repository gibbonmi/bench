## Outcome

The 15-ticket spec landed at `eaa2201c` from reviewed source `61f7ca36`. Seven chunk checkpoints and the final landing gate passed.
The completion record preserves five final checks and reconciles 86 coverage rows. WS45 and WS81 stay review-owned, and the other 84 rows cite tests.

The joins value now keeps 15 fields. Real fixtures replace 12 test seams, and one ambient value carries the home, clock, kit root, and writer.
A static census refuses reads below an entry. Production owns table names, and not-called tests prove that the verb used the joins route.

The comparison is one descriptive run. The dollar cells are API-equivalent estimates, not provider charges, and author costs stay separate from review costs.

| arm | author shape | author USD | review USD | total USD | wall time | blind rank |
|---|---|---:|---:|---:|---:|---:|
| A | Opus full slices | 6.57 | 2.77 | 9.34 | 50.5 min | 3 |
| B | Opus three-ticket slices | 8.96 | 3.14 | 12.10 | 69.0 min | 2 |
| C | Sonnet ticket slices | 6.04 | 6.86 | 12.90 | 71.2 min | 1 |

The reviewer selected arm B. The result does not change a model default because the run has no comparable repeat.
All four helper worktrees released after their gate logs moved to the primary capture.

The final-close dogfood wrote the missing completion record against the original source. The record preserved five checks and reconciled all 86 coverage rows.
The kit fix landed green at `6281fc3f`. Focused checks passed, and an independent mutation changed covered evidence to pending.
Its census recorded two raw Python calls. The fix gate passed all required phases and skipped ShellCheck by capability.

## Gate-stage timings

- landing: commit eaa2201cea5fd35cbedb79bf91a5b41a9d5edad1, trace f83af7c35844a26a13f5da6c465db6cb
- gofmt: 118 ms
- vet: 1310 ms
- test: 190451 ms
- race: 3096 ms
- system: 46969 ms
- shellcheck: 52 ms (skipped)

The gate reported eight capability skips: four FIFO skips and four privilege skips.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 1 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 2 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 3 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 4 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 5 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 6 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 7 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 8 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 9 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 10 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 11 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 12 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 13 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 14 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |
| ticket 15 | its planned acceptance rows hold | verified | unknown | held | opus / high / implementer |

Brier mean: unknown over 0 labeled pairs. Abstentions: 15, because the retained ticket records preserve no integer author confidence.

Each ticket stayed in its planned chunk. Tickets 8 and 12 needed two and three author attempts, and each other initial author used one attempt.
Arm B had two lane failures on the ticket 12 file budget. The author reduced the file within the three-attempt cap.
Reviews returned one repair round for tickets 1, 2, 3, 5, 6, 7, 8, 14, and 15. Ticket 12 needed two rounds.

## Coordinator catches

- The independent review found one-source defects in tickets 1, 6, 8, 12, and 14. Each repair removed the second derivation.
- Coverage review found silent mutations in tickets 1, 2, 3, 5, 7, 12, and 15. Each repair added a biting check.
- The coordinator probed every accepted repair at another site or kind. Each retained coordinator probe turned red.
- The completion dogfood found that the final evidence had no supported writer. The kit repair added that writer after the original build finished.
- Three later gates passed their test phases but found generated package artifacts. The system npm ran `prepare` despite `--ignore-scripts`.
- The corrected Node 25 and npm 11 path created no artifacts. This host defect belongs to tooling, not to an author or the spec.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-read-the-kit-value-once-in-gate.md | 1 | one-source |
| 10-fault-the-merge-reconcile-with-a-stale-index-lock.md | 0 | none |
| 11-lift-each-read-below-a-census-entry.md | 0 | none |
| 12-refuse-a-read-below-a-census-entry.md | 2 | one-source, delegate-error |
| 13-make-the-serial-ceiling-exact.md | 0 | none |
| 14-name-each-table-in-production.md | 1 | one-source |
| 15-require-the-joins-route-of-a-not-called-stub.md | 1 | delegate-error |
| 2-carry-an-ambient-value-below-each-verb-entry.md | 1 | delegate-error |
| 3-pass-the-kit-value-to-merge-and-land.md | 1 | delegate-error |
| 4-interrupt-the-landing-marker-with-a-gate-script.md | 0 | none |
| 5-fault-the-landing-follow-on-steps-with-real-fixtures.md | 1 | delegate-error |
| 6-land-the-stubbed-landing-tests-for-real.md | 1 | one-source |
| 7-fault-the-cleanup-reads-with-real-fixtures.md | 1 | delegate-error |
| 8-fault-the-reauthorize-unlock-with-a-denied-admin-directory.md | 1 | one-source |
| 9-fault-the-reset-move-with-real-fixtures.md | 0 | none |

## Agent-experience improvements

### Bench CLI

- Replace the `sr-integration census 30 raw calls` flow with one bounded command that reports command heads and output volume.
  Feeds: new
- Make package dry-run use the declared Node toolchain, so an older npm cannot create gate artifacts through `prepare`.
  Feeds: new

### Skills

- Require final-check to dogfood the completion writer before landing, so a missing evidence path fails before the source release.
  Feeds: none
- Require comparison reviews to stop at their declared test-run cap, so the retained cost and quality measures stay comparable.
  Feeds: none

### Process

- Preserve each helper gate log before worktree release, so the comparison keeps exact green evidence after cleanup.
  Feeds: none
- Keep the arm result descriptive until a comparable repeat supports a routing decision.
  Feeds: none
