# Local shift evidence review record

Status: LE-A review round 1 returned 10 findings; repairs are pending.
Spec: specs/local-shift-evidence/spec.md
Assignment: 8854df6a652ec4400d952339b55940b6
Author: claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN
Line: opus (claude-opus-5-5) / medium / uncapped
Review line: opus / high / one iteration for each axis
Post-review repair cycles consumed: LE-A 0 of 2
Expected repair rounds: 2
Confidence: 5

## LE-A author verification

The frozen pair is `ba9b8621..fff1e5ab`. The four planned checks passed at the tip.

| Row | Test | Probe |
|---|---|---|
| LE11 | `TestAHandedOffChildJoinsItsParentSpan` | Omit the traceparent inject: bit. |
| LE12 | `TestTheGateSpellsNoHandoffVariable` | Restore a handoff literal in `telemetry.go`: bit. |
| LE1 | `TestEncodeWritesTheRecordSchema` | Pre-edit red. |
| LE2 | `TestEncodeWritesTheSetVersion` | Pre-edit red. |
| LE3 | `TestOtelGateRecordNamesTheStampedVersion` | Omit `SetVersion` (system suite, copy aside): bit. |
| LE4 | `TestBeginWritesNoEnvironmentResource` | Write the SDK resource back: bit when the test runs alone. |
| LE5 | `TestEncodeResourceHoldsExactlyTheBenchKeys` | Write the SDK resource back: bit. |
| LE96 | `TestEncodeWithoutAVersionWritesNoVersionKey` | Pre-edit red. |
| LE6 | `TestEncodeDropsAnUndeclaredAttribute` | Turn the filter off: bit. |
| LE7 | `TestEncodeKeepsEveryDeclaredAttribute` | Drop `bench.record` in the filter: bit. |
| LE8 | `TestReadSelectedReportsAnUnknownSchemaAsMalformed` | Pre-edit red. |
| LE9 | `TestReadSpansReturnsNoSpanOfAnUnknownSchema` | Pre-edit red. |
| LE10 | `TestReadSpansReadsALegacyLine` | Require the schema key: bit. |
| LE13 | `TestAnAppendPastTheLimitSealsTheLiveSegment` | Pre-edit red. |
| LE92 | `TestTwoRotationsSealConsecutiveSequences` | Pre-edit red. |
| LE94 | `TestARotationSkipsAPlantedSequenceName` | Ignore the present names: bit. Remove only the free-name loop: silent. |
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

- LE4 shares a process-wide cache with the SDK default resource. In a full package run, an earlier test builds that resource first, so LE4 bites only alone. LE5 bites the same mutation in the full run.
- LE94: the sealed-name listing already counts a planted file, so the free-name loop guards only a race with a writer outside the lock. The spec requires the loop, and the build keeps it.
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

```bench-review-record
{
  "version": 1,
  "spec": "specs/local-shift-evidence/spec.md",
  "plan_digest": "sha256:7fb022a3db20bc5ce8289757ac685876565f9812804aa23984d52761468ce096",
  "implementation_session": "claude-code:session_01PzPVd5kMaFqKt7bLjSjtgN",
  "chunks": [
    {
      "id": "LE-A",
      "base": "ba9b86216536fa93f9f51410c31aea8e6f61725b",
      "tip": "fff1e5abedad3f270e0385e46fc743a09c2c67ef",
      "plan_digest": "sha256:7fb022a3db20bc5ce8289757ac685876565f9812804aa23984d52761468ce096",
      "source_digest": "0bef00fe9c58c2f02b2041309d34a355398aba28",
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
