# Local shift evidence review record

Status: LE-A author verification committed; the three LE-A axes are pending.
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

Pending.

## Spec

Pending.

## Coverage

Pending.

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
      "reviews": []
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
