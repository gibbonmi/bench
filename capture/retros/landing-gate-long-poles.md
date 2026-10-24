## Outcome

FT402 landed as `f4c6107080dd286f2d0acdfc3df294978b386f19` from source tip `a4c724413` on base `7f585da44`. The build ran 16 tickets in 10 chunks as a `--full --auto-approve` run with an author limit of 1. The reviewer was away and deferred every reviewer-owned decision to an Opus/high and Sonnet/high consultant pair, with Fable/high as the tie-breaker.

The build shares the probe and testreport fixtures once for each package run. It adds a git run value that caches location answers, and it carries the call directory and the child environment through each verb run. The system suite now runs in parallel in one process behind a census with a closed serial table. The serial table fell from 105 rows to 2 owner-repository rows.

Every chunk passed its checkpoint. The final verification passed 9 of 9 requirements, and the completion reconciled all 38 coverage rows. The retained red proofs are 30 plan probes, and each one bit with its restore passed.

## Gate-stage timings

- landing: commit f4c6107080dd286f2d0acdfc3df294978b386f19, trace 1ea550d21b8de326b0bd75f5714a6b6a
- gofmt: 182 ms
- vet: 1420 ms
- test: 168867 ms
- race: 2248 ms
- system: 35104 ms
- shellcheck: 547 ms

The baseline gates on `main` took 296 to 308 s, with a test phase of 173 to 174 s and a system phase of 100 to 108 s. The LGP-C4B checkpoint gate took about 214 s with a system phase of 39.2 s, and the landing gate took about 208 s. The build saves about 90 to 100 s for each landing, nearly all of it in the system phase. The test phase did not change outside the spread, because `internal/probe` fell from 88 to 90 s to 33 s but other packages still bound the phase.

## Ticket-versus-spec-slice and delegate performance

Each ticket fit its slice, and no author crossed a chunk boundary. Two first authors stopped before a commit: ticket 05 met a wrong predicted red count, and ticket 08 met a decision that the single-read census refuses. Fresh authors landed both after a spec correction. Ticket 14's author found and fixed a real leak: a crash test left phase processes that re-created deleted `/dev/shm` temp roots. The LGP-C4B authors moved 103 serial tests into the parallel set, and no moved test stayed serial.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 14 | each test in the eight files passes in parallel in two runs | verified | 9 | true | opus / xhigh / implementer |
| ticket 15 | LGP34 times recorded before and after | verified | 9 | true | opus / xhigh / implementer |
| ticket 15 | each test in the five files passes in parallel in two runs | verified | 9 | true | opus / xhigh / implementer |
| ticket 15 | the injected-port registry still finds TestEvidenceInterruptedPublication | verified | 9 | true | opus / xhigh / implementer |
| ticket 15 | charge_evidence_test.go stays at or under 400 lines | verified | 10 | true | opus / xhigh / implementer |
| ticket 15 | the probe reds the census | verified | 10 | true | opus / xhigh / implementer |
| ticket 16 | a third owner row, a non-reaching subject, and a missing subject red by name | verified | 9 | true | opus / xhigh / implementer |
| ticket 16 | a pending-migration row reds by name | verified | 9 | true | opus / xhigh / implementer |
| ticket 16 | the system argv stays unchanged | verified | 10 | true | opus / xhigh / implementer |
| ticket 16 | the pending-migration probe bites | verified | 10 | true | opus / xhigh / implementer |
| LGP-C4B Standards | S1 stale serial-table comment | finding | 6 | true | opus / medium / reviewer |
| LGP-C4B Standards | S2 duplicated pid-file wait | finding | 5 | true | opus / medium / reviewer |
| LGP-C4B Standards | S3 incomplete awaitSpan comment | finding | 4 | true | opus / medium / reviewer |
| LGP-C4B Coverage | C1 whitespace or near-miss serial reason passes | finding | 8 | true | opus / medium / reviewer |
| LGP-C4B Coverage | C2 non-literal TMPDIR binding passes | finding | 6 | true | opus / medium / reviewer |
| LGP-C4B Coverage | C3 phase kill only on the success path | finding | 6 | true | opus / medium / reviewer |
| LGP-C4B r1 Standards | S5 the no-reason rule has two sources | finding | 5 | true | opus / medium / reviewer |
| LGP-C4B r1 Standards | S6 comment restates enforced rules | finding | 3 | true | opus / medium / reviewer |

Brier mean 0.107 over 18 pairs; 2 abstentions (Standards S4 parked with no label, and ticket 16's LGP35 and LGP36 rows owned by the orchestrator). The table holds only the LGP-C4B returns, because the context compaction removed the earlier chunks' claim tables.

## Coordinator catches

- The coordinator refreshed four stale census citations after ticket 15 moved `TestEvidenceStatelessProcesses`, in a plan commit before the ticket 16 dispatch.
- The coordinator closed the LGP34 record gap of LGP-C4B with an evidence-only Spec excerpt that pairs each ticket's before and after times.
- The coordinator routed S4 (rename the system-suite journey-repository wrapper) and C2 (a wider process-state rule) to the consultants. They agreed on no repair, and both became `bench idea` entries.
- The coordinator committed the LGP-C4B repair records before four unrepaired tickets restaged, so the chunk had to refreeze at a record-only commit, and all six tickets restaged.
- The coordinator first landed with a guessed request token and a reauthorize that unbound the FT402 delivery. Reauthorizing with the token of `bench commitment start` restored the binding. The scorecard already held this rule.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-share-the-probe-clamp-golden.md | 1 | one-source |
| 02-build-the-testreport-selection-once.md | 1 | one-source |
| 03-add-the-git-run-value.md | 1 | one-source |
| 04-carry-the-run-value-and-read-owner-markers-once.md | 0 | none |
| 05-read-the-explicit-plan-locations-once.md | 0 | none |
| 06-carry-the-call-directory.md | 2 | edge-found, one-source |
| 07-move-the-working-directory-tests.md | 0 | none |
| 08-carry-the-child-environment.md | 1 | edge-found |
| 09-move-the-exec-child-environment-tests.md | 1 | one-source |
| 10-add-the-system-parallel-census.md | 3 | edge-found, edge-found, edge-found |
| 11-run-the-owner-and-recovery-tests-in-parallel.md | 0 | none |
| 12-run-the-compatibility-and-adoption-tests-in-parallel.md | 0 | none |
| 13-run-the-session-start-and-guard-tests-in-parallel.md | 0 | none |
| 14-run-the-land-and-telemetry-tests-in-parallel.md | 1 | edge-found |
| 15-run-the-exec-and-charge-tests-in-parallel.md | 0 | none |
| 16-close-the-system-serial-table.md | 2 | edge-found, one-source |

## Agent-experience improvements

### Bench CLI

- Add a plan verb that lists and appends the verification rows of a completion plan, because 25 raw `python3` calls read or edited the plan JSON (owned gap).
  Feeds: new
- Add a review-record reader that returns one chunk entry or one verification by ID, because authors and axes used raw `rg` with offsets to find entries in a 5000-line record (owned gap).
  Feeds: new
- Count raw `rg` calls on Bench-owned files apart from source searches in the census, because 164 sanctioned `rg` calls hide the owned-gap reads of the plan and the review record.
  Feeds: new
- Route a line-count check to `bench structure` headroom for a named file, because 6 raw `wc` calls measured line caps that Bench already grades (missed form).
  Feeds: FT307
- Make `bench record ticket` say that only the chunk's last ticket refreezes, when a non-last ticket passes `--tip` at a newer commit.
  Feeds: new
- Make the `bench probe --omit` ambiguity refusal list the line of each match, and quote the first census report line in the probe failure row.
  Feeds: FT168
- Spill the full failure body of a multi-line `bench test` failure, because the failures row shows only its first line.
  Feeds: FT415
- Accept `--with-file` on each repeated `--swap` of `bench probe`, so a multi-line insertion needs no shell quoting.
  Feeds: new
- The 15 raw `cat`, `ls`, `sed`, `awk`, `cd`, and `for` calls have no recorded purpose (unknown), so this retro proposes no change for them.
  Feeds: none

### Skills

- State in the `bench-implement-spec` Land section that after a refreeze of a multi-ticket chunk, the record commit waits until every ticket restages at the new tip.
  Feeds: new
- Name `bench record verification --final` and `bench record completion` in the `bench-implement-spec` Land section, with the source that each takes, because the landing refusal was the only route to them.
  Feeds: new

### Process

- After a context compaction, reread the scorecard's Current decisions before the landing, because this build broke the request-token rule that the scorecard already held.
  Feeds: none
- Write the commitment request token into the handoff state file, so a resumed orchestrator lands with the bound token.
  Feeds: new