# Local shift evidence review record

Status: the LE-A confirming round returned 4 findings; repair cycle 2 is pending.
Spec: specs/local-shift-evidence/spec.md
Assignment: 8854df6a652ec4400d952339b55940b6
Author: claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN
Line: opus (claude-opus-5-5) / medium / uncapped
Review line: opus / high / one iteration for each axis
Post-review repair cycles consumed: LE-A 1 of 2
Expected repair rounds: 2
Confidence: 5

## LE-A author verification

The first frozen pair was `ba9b8621..fff1e5ab`. After repair cycle 1, the chunk tip is `33b65a3e`, and the four planned checks passed there too.

| Row | Test | Probe |
|---|---|---|
| LE11 | `TestAHandedOffChildJoinsItsParentSpan` | Omit the traceparent inject: bit. |
| LE12 | `TestTheGateSpellsNoHandoffVariable` | Restore a handoff literal in `telemetry.go`: bit. |
| LE1 | `TestEncodeWritesTheRecordSchema` | Pre-edit red. |
| LE2 | `TestEncodeWritesTheSetVersion` | Pre-edit red. |
| LE3 | `TestOtelGateRecordNamesTheStampedVersion` | Omit `SetVersion` (system suite, copy aside): bit. |
| LE4 | `TestBeginWritesNoEnvironmentResource` | Write the SDK resource back: bit in the full package run (helper process). |
| LE5 | `TestEncodeResourceHoldsExactlyTheBenchKeys` | Write the SDK resource back: bit. |
| LE96 | `TestEncodeWithoutAVersionWritesNoVersionKey` | Pre-edit red. |
| LE6 | `TestEncodeDropsAnUndeclaredAttribute` | Turn the filter off: bit. |
| LE7 | `TestEncodeKeepsEveryDeclaredAttribute` | Drop `bench.record` in the filter: bit. |
| LE8 | `TestReadSelectedReportsAnUnknownSchemaAsMalformed` | Pre-edit red. |
| LE9 | `TestReadSpansReturnsNoSpanOfAnUnknownSchema` | Pre-edit red. |
| LE10 | `TestReadSpansReadsALegacyLine` | Require the schema key: bit. |
| LE13 | `TestAnAppendPastTheLimitSealsTheLiveSegment` | Pre-edit red. |
| LE92 | `TestTwoRotationsSealConsecutiveSequences` | Pre-edit red. |
| LE94 | `TestARotationSkipsAPlantedSequenceName` | Ignore the present names: bit. |
| LE14 | `TestThePruneKeepsTheRetainedCount` | Drop the prune: bit. |
| LE93 | `TestThePruneRemovesTheLowestSequenceWhateverItsTime` | Drop the prune: bit. |
| LE15 | `TestReadSpansReadsTheSealedSegmentsInOrder` | Read the live segment only: bit. |
| LE16 | `TestNewestLandingReadsStagesFromASealedSegment` | Read the live segment only: bit. |
| LE17 | `TestReadSelectedReadsASealedSegment` | Read the live segment only: bit. |
| LE18 | `TestTwoWritersAcrossRotationsLeaveOnlyWholeLines` | Truncate instead of rename: bit. |
| LE19 | `TestAHeldRotationLockSkipsTheRotation` | Rotate without the lock: bit. |
| LE20 | `TestReadSpansRefusesAFIFOSegment` | Omit the segment grade loop: bit. |
| LE21 | `TestReadSpansRefusesASymlinkedSegment` | Omit the segment grade loop: silent, because the no-follow open also refuses the link. |

## LE-A author notes

- LE4 writes its record in a helper process, because the SDK builds its default resource once for each process.
- LE94: the reviewer decided that the sealed-name listing is the existence check, so the writer has no free-name step.
- LE21 has two defenses: the grade and the no-follow open. Each defense alone keeps the row green.
- An unstamped `dev` build hands no version to the record, so its lines carry no `service.version` key. The spec names a stamped version only.
- `bounds.RecordSegmentLimit` is spelled `1 << 24`. The spelling `16 << 20` collides by text with an unrelated value in `internal/harnesstranscript/read.go`, and the bounds-policy check then reds. A learning records this call for the reviewer.

## Standards

Findings: 4. Worst issue: LEA-S2, a second parser of the sealed name in the test.

- LEA-S1 (auto-fix, confidence 6): `resourceOf`, `spanAttributeKeys`, and `schemaLine` each copy the encode-and-parse steps that `encodeFixture` owns (`internal/otelrecord/encode_test.go`, `internal/otelrecord/reader_test.go`). Rule: `AGENTS.md`, one source per fact, fixture harnesses.
- LEA-S2 (auto-fix, confidence 6): `sealedNames` in `internal/otelrecord/retention_test.go` parses sealed names with its own `traces-` prefix and accepts names that `parseSealedName` rejects. Rule: one source per fact, and the spec's one sealed-name helper.
- LEA-S3 (auto-fix, confidence 7): two new test comments tell history (`processor_test.go` LE4 comment, `encode_test.go` `sdkResourceSpan` comment). Rule: `craft-comments` present tense, invariant 3.
- LEA-S4 (auto-fix, confidence 5): the comment above `version` in `cmd/bench/main.go` still spells `"dev"`, which `unstampedVersion` now owns. Rule: `craft-comments`, one source owns a fact.

## Spec

Findings: 3. Worst issue: LEA-P1, where the LE94 test cannot fail on the mutation its row names.

- LEA-P1 (ask-user, confidence 8): LE94 says a rotation with no existence check reds the test. The sealed-name listing already counts the planted file, so the free-name loop is unreachable under the lock, and its removal stays green. LEA-C2 names the same defect.
- LEA-P2 (auto-fix, confidence 6): LE4 can fail only when it runs alone, because the SDK caches its default resource once per process. The row must fail in the chunk's own package run.
- LEA-P3 (auto-fix, confidence 8): the 25 LE-A rows still carry a planned seam. The spec says that the build adds each citation when its test lands.

## Coverage

Findings: 3. Worst issue: LEA-C1, where a reader can return an incomplete record with no error.

- LEA-C1 (ask-user, confidence 7): a rotation between the sealed-name listing and the live open seals a segment that the reader never lists. The reader then misses those spans and reports no error. No row and no Won't handle line decides this interleaving.
- LEA-C2 (ask-user, confidence 9): merged with LEA-P1.
- LEA-C3 (auto-fix, confidence 7): `ReadSelected` numbers lines across every segment, and the assessment collection attaches each problem to the live path. A problem in a sealed segment then names a wrong line of `traces.jsonl`.

Raw findings: 10. Repair targets after the merge of LEA-P1 and LEA-C2: 9.

## LE-A repair cycle 1

The reviewer closed LEA-P1 and LEA-C2 with the removal of the free-name step and a new LE94 reason. The reviewer closed LEA-C1 with a Won't handle line. The ticket `repair-le-a-review-round-1.md` covers the amended row.

| Finding | Repair | Evidence |
|---|---|---|
| LEA-S1 | `encodeInto` owns the encode and parse steps. | The record tests pass. |
| LEA-S2 | `sealedNames` reads through `sealedSequences`. | The record tests pass. |
| LEA-S3 | The two comments state the present behavior. | Review of the text. |
| LEA-S4 | The `main.go` comment names `unstampedVersion`. | Review of the text. |
| LEA-P1, LEA-C2 | The writer has no free-name step, and LE94 grades a rotation that ignores the present names. | That probe bit. |
| LEA-P2 | LE4 writes its record in a helper process. | The SDK-resource probe bit LE4 in the full package run. |
| LEA-P3 | Each LE-A row cites its test. | `bench coverage --check` reports a valid map. |
| LEA-C1 | A Won't handle line records the reader window. | Reviewer decision. |
| LEA-C3 | `ReadSelected` names the sealed segment and counts lines per segment. | `TestReadSelectedNamesTheSegmentOfAProblem`; the cross-segment probe bit. |

## LE-A confirming round

The three axes graded the pair `ba9b8621..721250ae` with the source digest `869256de`. Every round 1 finding is closed, except that LEA-C3 is only half closed (LEA-C4).

- Standards, 2 findings. Worst issue: LEA-S5.
  - LEA-S5 (auto-fix, confidence 8): the `TestARotationSkipsAPlantedSequenceName` comment still names the removed free-name check. Rule: `craft-comments`, update a comment with its code.
  - LEA-S6 (auto-fix, confidence 7): the `TestReadSelectedNamesTheSegmentOfAProblem` comment cites a review finding ID. Rule: `craft-comments`, no provenance.
- Spec, 0 findings.
- Coverage, 2 findings. Worst issue: LEA-C4.
  - LEA-C4 (auto-fix, confidence 9): no test fails when the per-segment line reset is removed. A live-segment problem after a sealed segment can then name a wrong line.
  - LEA-C5 (ask-user, confidence 6): a planted name at the maximum sequence makes the next sequence wrap to 0. Each later rotation then renames over sealed segment 0. The reviewer decided on 2026-09-23 that the writer refuses that rotation.

Raw findings: 4. Repair targets: 4.

```bench-review-record
{
  "version": 1,
  "spec": "specs/local-shift-evidence/spec.md",
  "plan_digest": "sha256:596aee0bd583c82d52e1eb70145f7f41cd945d378da153ad5f7d14a0660ea1fc",
  "implementation_session": "claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN",
  "chunks": [
    {
      "id": "LE-A",
      "base": "ba9b86216536fa93f9f51410c31aea8e6f61725b",
      "tip": "33b65a3e073da681cf8e45c1ce02233e493dd4f9",
      "plan_digest": "sha256:596aee0bd583c82d52e1eb70145f7f41cd945d378da153ad5f7d14a0660ea1fc",
      "source_digest": "869256de0fde5b87aa31099a0910430f83f1bf71",
      "acceptance_rows": [
        "LE1",
        "LE2",
        "LE3",
        "LE4",
        "LE5",
        "LE96",
        "LE6",
        "LE7",
        "LE8",
        "LE9",
        "LE10",
        "LE11",
        "LE12",
        "LE13",
        "LE92",
        "LE94",
        "LE14",
        "LE93",
        "LE15",
        "LE16",
        "LE17",
        "LE18",
        "LE19",
        "LE20",
        "LE21"
      ],
      "verification": [
        {
          "id": "LE-A-verify-otelrecord",
          "performer": "claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN",
          "role": "author-verification",
          "model": "claude-opus-5-5",
          "effort": "medium",
          "source_digest": "0bef00fe9c58c2f02b2041309d34a355398aba28",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "bench test --package ./internal/otelrecord at fff1e5ab",
            "digest": "sha256:516fc6f027a5b7af046693f400afa77e2519d2be14fcee04fc3c48010adde55b",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/otelrecord,pass,217\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "otelrecord",
          "command": "bench test --package ./internal/otelrecord",
          "exit_code": 0
        },
        {
          "id": "LE-A-verify-gate",
          "performer": "claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN",
          "role": "author-verification",
          "model": "claude-opus-5-5",
          "effort": "medium",
          "source_digest": "0bef00fe9c58c2f02b2041309d34a355398aba28",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "bench test --package ./internal/gate at fff1e5ab",
            "digest": "sha256:0c5e996e1ea23d8ec58399660699508323f97ebd3a5eccd213cd6a93965120a7",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,10402\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "LE-A-verify-cmd",
          "performer": "claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN",
          "role": "author-verification",
          "model": "claude-opus-5-5",
          "effort": "medium",
          "source_digest": "0bef00fe9c58c2f02b2041309d34a355398aba28",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "bench test --package ./cmd/bench at fff1e5ab",
            "digest": "sha256:33d9a3351ecdc2441e2531c891d99e42b3ae6f50c9ed2db4f7f5cc21c49d5a15",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,11353\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "LE-A-verify-system",
          "performer": "claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN",
          "role": "author-verification",
          "model": "claude-opus-5-5",
          "effort": "medium",
          "source_digest": "0bef00fe9c58c2f02b2041309d34a355398aba28",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "bench test --check system at fff1e5ab",
            "digest": "sha256:366d7196552a2e5b1e75cc2c8d83be04e4e1277565863702e0f97dade2d203c7",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,39281\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "LE-A-verify-otelrecord-2",
          "performer": "claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN",
          "role": "author-verification",
          "model": "claude-opus-5-5",
          "effort": "medium",
          "source_digest": "869256de0fde5b87aa31099a0910430f83f1bf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "bench test --package ./internal/otelrecord at 33b65a3e",
            "digest": "sha256:bdf956482cfad7d4a811880d747b1d9be07a6c274d312bbd92c507bd1578568d",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/otelrecord,pass,196\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "otelrecord",
          "command": "bench test --package ./internal/otelrecord",
          "exit_code": 0
        },
        {
          "id": "LE-A-verify-gate-2",
          "performer": "claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN",
          "role": "author-verification",
          "model": "claude-opus-5-5",
          "effort": "medium",
          "source_digest": "869256de0fde5b87aa31099a0910430f83f1bf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "bench test --package ./internal/gate at 33b65a3e",
            "digest": "sha256:13206398342dcb25443171220ade3bbccda4a8b60f97966be85a5315a765cd9b",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,9569\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "LE-A-verify-cmd-2",
          "performer": "claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN",
          "role": "author-verification",
          "model": "claude-opus-5-5",
          "effort": "medium",
          "source_digest": "869256de0fde5b87aa31099a0910430f83f1bf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "bench test --package ./cmd/bench at 33b65a3e",
            "digest": "sha256:67e58bc8519d40a16d1d087fc0fba64bd224246a009c0ef8714d176b02d5bf21",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,9956\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "LE-A-verify-system-2",
          "performer": "claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN",
          "role": "author-verification",
          "model": "claude-opus-5-5",
          "effort": "medium",
          "source_digest": "869256de0fde5b87aa31099a0910430f83f1bf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "bench test --check system at 33b65a3e",
            "digest": "sha256:58668e2606831bae902b3c79c1dd9befd350df464aa2e1ea71543bb80611796d",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,39064\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "system",
          "command": "bench test --check system",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "LE-A-review-standards-1",
          "performer": "claude-code:subagent:a29c984de1edb6a27",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0bef00fe9c58c2f02b2041309d34a355398aba28",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-code subagent claude-code:subagent:a29c984de1edb6a27, evidence sha256:51031371421e116a8ed1af81494f5695fdbfbf470fb37216de836e1d28bcf08e",
            "digest": "sha256:ffce5a2806bb8bc98ee6164ae4497ba82cef96b36c41f3096469e74a0ff7bf7a",
            "excerpt": "Standards axis, LE-A. Outcome: 4 findings (3 hard, 1 judgment call). Worst: LEA-S2, because the test's copy of the sealed-name parser accepts names the production parser rejects."
          },
          "axis": "Standards",
          "base": "ba9b86216536fa93f9f51410c31aea8e6f61725b",
          "tip": "817915b45963f432d5373466257d9544de6ada11",
          "finding_ids": [
            "LEA-S1",
            "LEA-S2",
            "LEA-S3",
            "LEA-S4"
          ],
          "supersedes": []
        },
        {
          "id": "LE-A-review-spec-1",
          "performer": "claude-code:subagent:a265ef2193aedd2c0",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0bef00fe9c58c2f02b2041309d34a355398aba28",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-code subagent claude-code:subagent:a265ef2193aedd2c0, evidence sha256:51031371421e116a8ed1af81494f5695fdbfbf470fb37216de836e1d28bcf08e",
            "digest": "sha256:e68bc8f3a4f715885f19111ac63f1a4cd67d8602ebc1059cbf1c1a46ce5f9c7e",
            "excerpt": "Spec axis, chunk LE-A. outcome: findings (3). Worst issue: LEA-P1. The LE94 test cannot fail when the spec-required free-name check is deleted."
          },
          "axis": "Spec",
          "base": "ba9b86216536fa93f9f51410c31aea8e6f61725b",
          "tip": "817915b45963f432d5373466257d9544de6ada11",
          "finding_ids": [
            "LEA-P1",
            "LEA-P2",
            "LEA-P3"
          ],
          "supersedes": []
        },
        {
          "id": "LE-A-review-coverage-1",
          "performer": "claude-code:subagent:abf61b440adc50922",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0bef00fe9c58c2f02b2041309d34a355398aba28",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-code subagent claude-code:subagent:abf61b440adc50922, evidence sha256:51031371421e116a8ed1af81494f5695fdbfbf470fb37216de836e1d28bcf08e",
            "digest": "sha256:aeca561cbec9897a443faaaacc1755b7a5867145745dd2629d30b5263307da85",
            "excerpt": "outcome: findings (3). Findings: 3. Worst issue: LEA-C1, where a reader can return an incomplete record with no error after a rotation."
          },
          "axis": "Coverage",
          "base": "ba9b86216536fa93f9f51410c31aea8e6f61725b",
          "tip": "817915b45963f432d5373466257d9544de6ada11",
          "finding_ids": [
            "LEA-C1",
            "LEA-C2",
            "LEA-C3"
          ],
          "supersedes": []
        },
        {
          "id": "LE-A-review-standards-2",
          "performer": "claude-code:subagent:a29c984de1edb6a27",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "869256de0fde5b87aa31099a0910430f83f1bf71",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-code subagent claude-code:subagent:a29c984de1edb6a27, evidence sha256:04fdd1a78e6d13f5964aa2b9d5dad0c2fca08d0f7264f7098cbe531c7e7c5c1c",
            "digest": "sha256:63fa5aee82e732343a1c89abeadf851520ab94d5fe7cd1451aefc920339da93b",
            "excerpt": "Standards confirming round, LE-A. Outcome: 2 findings. All four earlier folds are closed. Worst: LEA-S5, a test comment that describes a mutation class the reviewer removed."
          },
          "axis": "Standards",
          "base": "ba9b86216536fa93f9f51410c31aea8e6f61725b",
          "tip": "721250aec7be579afb0870975a9a30c57a0c4c6b",
          "finding_ids": [
            "LEA-S5",
            "LEA-S6"
          ],
          "supersedes": [
            "LE-A-review-standards-1"
          ]
        },
        {
          "id": "LE-A-review-spec-2",
          "performer": "claude-code:subagent:a265ef2193aedd2c0",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "869256de0fde5b87aa31099a0910430f83f1bf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-code subagent claude-code:subagent:a265ef2193aedd2c0, evidence sha256:04fdd1a78e6d13f5964aa2b9d5dad0c2fca08d0f7264f7098cbe531c7e7c5c1c",
            "digest": "sha256:5381fddf40018a8915f1e922307d94d3d0760d5d11b571081e3a1ad2a21f622f",
            "excerpt": "outcome: pass. The repair delta introduces no blocking Spec finding and leaves none open. Count: 0. Worst issue: none."
          },
          "axis": "Spec",
          "base": "ba9b86216536fa93f9f51410c31aea8e6f61725b",
          "tip": "721250aec7be579afb0870975a9a30c57a0c4c6b",
          "finding_ids": [],
          "supersedes": [
            "LE-A-review-spec-1"
          ]
        },
        {
          "id": "LE-A-review-coverage-2",
          "performer": "claude-code:subagent:abf61b440adc50922",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "869256de0fde5b87aa31099a0910430f83f1bf71",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-code subagent claude-code:subagent:abf61b440adc50922, evidence sha256:04fdd1a78e6d13f5964aa2b9d5dad0c2fca08d0f7264f7098cbe531c7e7c5c1c",
            "digest": "sha256:30ab276a2b4e09d6b839965a2431097f73118a49ae032ec860334130927dc2b2",
            "excerpt": "outcome: findings (2). Findings: 2. Worst issue: LEA-C4. The per-segment line count has no test that can fail, so LEA-C3 is only half closed."
          },
          "axis": "Coverage",
          "base": "ba9b86216536fa93f9f51410c31aea8e6f61725b",
          "tip": "721250aec7be579afb0870975a9a30c57a0c4c6b",
          "finding_ids": [
            "LEA-C4",
            "LEA-C5"
          ],
          "supersedes": [
            "LE-A-review-coverage-1"
          ]
        }
      ]
    }
  ],
  "completion": {
    "state": "",
    "source_digest": "",
    "performer": "",
    "reconciliation": {},
    "verification": []
  }
}
```
