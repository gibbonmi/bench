# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/comment-gap-chain/spec.md",
  "plan_digest": "sha256:aa15a1af150c53a83df89c155d1d30f25af15bd1fddf360feba92d60b9f278f0",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CG-C1",
      "base": "beb3b30219e9665db63b83c6d385a7a0ccb74553",
      "tip": "890370c84f8a404e38bff0d77873f9b6b119368d",
      "plan_digest": "sha256:aa15a1af150c53a83df89c155d1d30f25af15bd1fddf360feba92d60b9f278f0",
      "source_digest": "3a5a29097ed8fb1b61e79f1959bef281bdc80193",
      "acceptance_rows": [
        "CG01",
        "CG02",
        "CG03",
        "CG04",
        "CG05",
        "CG06",
        "CG07",
        "CG08",
        "CG09",
        "CG10",
        "CG11",
        "CG12",
        "CG13",
        "CG14",
        "CG15",
        "CG16",
        "CG17"
      ],
      "verification": [
        {
          "id": "cg-c1-chain-gap",
          "performer": "ft439-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3a5a29097ed8fb1b61e79f1959bef281bdc80193",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft439-author-01",
            "digest": "sha256:e6b4225b9aadecacb977d650ff7477e63586febd438889daf35551856642cc18",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/fold.go,swap,failed,1,TestChainAcceptsAProvedCommentGap/fold_then_comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord/...,Chain|Fold|ReviewRecordSource|CommentOnlyGap,passed,35\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,756\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,2411,17\n  github.com/gibbonmi/bench/internal/reviewrecord/carry,pass,3,0\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,1708,18\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestChainAcceptsAProvedCommentGap/fold_then_comment,\"fold_test.go:117: chunk 2 after a proved comment gap: chunk 2: stale review chain gap: expected base 96b0f7468a3461725501c1b0c286a5852b9a7b0d, the chunk 1 tip, got 3ac76fb646ff37df51d3ef2b9998ca03f64dd564; plan commits land before the ticke… (490 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "890370c84f8a404e38bff0d77873f9b6b119368d",
          "requirement": "chain-gap",
          "command": "bench test --package ./internal/reviewrecord/... --run 'Chain|Fold|ReviewRecordSource|CommentOnlyGap'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make foldCoversBase compare the fold digest with the base digest by equality, and record that the fold-then-comment case of TestChainAcceptsAProvedCommentGap turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft439-author-01",
              "digest": "sha256:e6b4225b9aadecacb977d650ff7477e63586febd438889daf35551856642cc18",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/fold.go,swap,failed,1,TestChainAcceptsAProvedCommentGap/fold_then_comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord/...,Chain|Fold|ReviewRecordSource|CommentOnlyGap,passed,35\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,756\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,2411,17\n  github.com/gibbonmi/bench/internal/reviewrecord/carry,pass,3,0\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,1708,18\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestChainAcceptsAProvedCommentGap/fold_then_comment,\"fold_test.go:117: chunk 2 after a proved comment gap: chunk 2: stale review chain gap: expected base 96b0f7468a3461725501c1b0c286a5852b9a7b0d, the chunk 1 tip, got 3ac76fb646ff37df51d3ef2b9998ca03f64dd564; plan commits land before the ticke… (490 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "cg-c1-checkpoint",
          "performer": "ft439-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3a5a29097ed8fb1b61e79f1959bef281bdc80193",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft439-author-01",
            "digest": "sha256:d00c2c49807af841abc1c5c938971ae8cdcc9b887c6d01596b594050965b39b7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/coverage.go,swap,failed,2,TestReviewCheckpointChainAcceptsACommentGap/chunk,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,ReviewCheckpoint|DeltaCheckpointRunsWholeForANonTestDelta,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,882\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,10658,89\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointChainAcceptsACommentGap/chunk,\"review_checkpoint_commits_test.go:140: chunk checkpoint over a comment gap = 1  gate subject unavailable: completion evidence: chunk 2: stale review chain gap: expected base 4e98cd0486d13e0cc2d75a9ab7c5820bba8b9fbe, the chunk 1 tip, got 4da… (563 bytes)\",3\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointChainAcceptsACommentGap/complete,\"review_checkpoint_commits_test.go:140: complete checkpoint over a comment gap = 1  gate subject unavailable: completion evidence: chunk 2: stale review chain gap: expected base fc6b3ecb73bc264f6cfb678072e4e335f0af6e5e, the chunk 1 tip, got … (566 bytes)\",3\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "890370c84f8a404e38bff0d77873f9b6b119368d",
          "requirement": "checkpoint",
          "command": "bench test --package ./internal/gate --run 'ReviewCheckpoint|DeltaCheckpointRunsWholeForANonTestDelta'",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore the equality comparison of the base digest with the previous chunk digest in checkSource, and record that TestReviewCheckpointChainAcceptsACommentGap turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft439-author-01",
              "digest": "sha256:d00c2c49807af841abc1c5c938971ae8cdcc9b887c6d01596b594050965b39b7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/coverage.go,swap,failed,2,TestReviewCheckpointChainAcceptsACommentGap/chunk,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,ReviewCheckpoint|DeltaCheckpointRunsWholeForANonTestDelta,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,882\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,10658,89\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointChainAcceptsACommentGap/chunk,\"review_checkpoint_commits_test.go:140: chunk checkpoint over a comment gap = 1  gate subject unavailable: completion evidence: chunk 2: stale review chain gap: expected base 4e98cd0486d13e0cc2d75a9ab7c5820bba8b9fbe, the chunk 1 tip, got 4da… (563 bytes)\",3\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointChainAcceptsACommentGap/complete,\"review_checkpoint_commits_test.go:140: complete checkpoint over a comment gap = 1  gate subject unavailable: completion evidence: chunk 2: stale review chain gap: expected base fc6b3ecb73bc264f6cfb678072e4e335f0af6e5e, the chunk 1 tip, got … (566 bytes)\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": []
    }
  ],
  "completion": {
    "state": "pending",
    "source_digest": "",
    "performer": "",
    "reconciliation": {},
    "verification": []
  }
}
```

## CG-C1 ticket 01 author witness reds

The ticket 01 author `ft439-author-01` (opus / medium) wrote the new and changed tests first. It ran them in one `bench test` call against the unchanged production code at `fb429dbea`, and these cases went red.
The plan probes `chain-gap` and `checkpoint` answered `bit` after the ticket commit `890370c84`, and the record holds them.

| subject | condition | failing witness | row |
| --- | --- | --- | --- |
| `internal/reviewrecord/coverage.go` | equality branch of the chain | `TestChainAcceptsAProvedCommentGap/chunk_then_comment` | CG04 |
| `internal/reviewrecord/fold.go` | equality in `foldCoversBase` | `TestChainAcceptsAProvedCommentGap/fold_then_comment` | CG05 |
| `internal/reviewrecord/coverage.go` | equality branch of the chain | `TestReviewCheckpointChainAcceptsACommentGap/chunk` | CG01, CG02 |
| `internal/reviewrecord/coverage.go` | equality branch of the chain | `TestReviewCheckpointChainAcceptsACommentGap/complete` | CG03 |
| `internal/reviewrecord/coverage.go` | old refusal clause | `TestReviewCheckpointChainGapNamesTheExpectedBase` | CG11 |

The statement, directive, Markdown, and off-line refusal cases (CG06 to CG09) passed against the old code, because the old chain already refused them.
