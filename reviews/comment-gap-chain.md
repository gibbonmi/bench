# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/comment-gap-chain/spec.md",
  "plan_digest": "sha256:d30f2a0b57067b74d9fa2c187538c772cadc87ed52b44b722869b4fe67279ba4",
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
      "reviews": [
        {
          "id": "cg-c1-standards-r1",
          "performer": "ft439-review-standards-c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "3a5a29097ed8fb1b61e79f1959bef281bdc80193",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft439-review-standards-c1",
            "digest": "sha256:777cf8bd11b1a33c617d9be76f0b56cd38aa5cddc16c24410c2ecc3b9b7d1921",
            "excerpt": "Standards axis, CG-C1, frozen pair beb3b302..890370c8, sonnet/medium. Preflight bind current=true.\nFindings (initial round, all severity \u003c= 2, none a correctness defect; the coordinator routes each to optional advice under the bounded repair policy):\n- S1 sev 2: the chain base case of TestReviewCheckpointKeepsStrictEvidence (internal/gate/review_checkpoint_commits_test.go:245-249) repeats five setup calls of chainedOverCommentGap (107-121). Axis disposition: no-op; incidental text, a shared form would need a callback.\n- S2 sev 1: spec.md:285 rewrites the commitment.json citations in the Won't-handle line. Axis disposition: ask-user or no-op as citation upkeep.\n- S3 sev 2: test doc comments cite spec row IDs (\"covers CG01 to CG03\" at review_checkpoint_commits_test.go:124, \"covers CG04 to CG09\" at fold_test.go:159); craft-comments treats an identifier as provenance. Axis disposition: auto-fix.\nClean: CG12 single call site in digestFollows; refusal clause identical in code, test, and spec; doc comments of digestFollows, follows, foldCoversBase are timeless and start with their names.\nImplementation command contributed: no.\n"
          },
          "axis": "Standards",
          "base": "beb3b30219e9665db63b83c6d385a7a0ccb74553",
          "tip": "890370c84f8a404e38bff0d77873f9b6b119368d",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cg-c1-spec-r1",
          "performer": "ft439-review-spec-c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "3a5a29097ed8fb1b61e79f1959bef281bdc80193",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft439-review-spec-c1",
            "digest": "sha256:95b04c0fb60de4d3f20f04b52eaf7f669a07c9e2d6b7d2f6a9371fc0d75a25ab",
            "excerpt": "Spec axis, CG-C1, frozen pair beb3b302..890370c8, sonnet/medium. Preflight bind current=true.\nFindings: P1 sev 1: spec.md:286 rewrites the FT380 citation to commitment.json:966-967 and 2143-2149; the new lines are real. Axis disposition: no-op.\nPer-row: CG01-CG03 hold (TestReviewCheckpointChainAcceptsACommentGap, chunk and complete exit 0, CG02 compares Base, Tip, SourceDigest, Reviews, Verification). CG04-CG09 hold (TestChainAcceptsAProvedCommentGap: chunk-then-comment, fold-then-comment, statement, directive, markdown, off-line two-parent tip; literal expectations). CG10, CG11 hold. CG12 holds (one call site fold.go:145 in digestFollows). CG13-CG17 hold with no edit.\nFour reader sites read digestFollows; both ancestry conjuncts kept; chain order unchanged.\nAdvice: the chain-base case writes correctCommentGap and then overwrites the same file.\nImplementation command contributed: no.\n"
          },
          "axis": "Spec",
          "base": "beb3b30219e9665db63b83c6d385a7a0ccb74553",
          "tip": "890370c84f8a404e38bff0d77873f9b6b119368d",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cg-c1-coverage-r1",
          "performer": "ft439-review-coverage-c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "3a5a29097ed8fb1b61e79f1959bef281bdc80193",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft439-review-coverage-c1",
            "digest": "sha256:d24d1f30553e5855fe2bbdc017de69275575bc2d0e5ad060d4f59a98e2bcc761",
            "excerpt": "Coverage axis, CG-C1, frozen pair beb3b302..890370c8, sonnet/medium. Preflight bind current=true.\nFindings: C1 sev 1: the \"covered == later ||\" clause of digestFollows (internal/reviewrecord/fold.go:149) is redundant, because commentgap.Prove returns nil for equal trees (commentgap.go:27-29). Probe --omit \"covered == later || \" over reviewrecord Chain|Fold|ReviewRecordSource|CommentOnlyGap: silent, restored. Axis disposition: ask-user or no-op.\nNo row gap: CG04-CG09 cases present; CG09 bites only the ancestry conjunct; CG14 error path unchanged; four call sites read digestFollows.\nFinal git status: empty. Implementation command contributed: no.\n"
          },
          "axis": "Coverage",
          "base": "beb3b30219e9665db63b83c6d385a7a0ccb74553",
          "tip": "890370c84f8a404e38bff0d77873f9b6b119368d",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "CG-C2",
      "base": "890370c84f8a404e38bff0d77873f9b6b119368d",
      "tip": "dd48bf924f651ae4d5a695b59d20391bba07b74c",
      "plan_digest": "sha256:d30f2a0b57067b74d9fa2c187538c772cadc87ed52b44b722869b4fe67279ba4",
      "source_digest": "c2dded1aef96d333ce550a709c25b5f4b59f4f03",
      "acceptance_rows": [
        "CG18",
        "CG19",
        "CG20",
        "CG21",
        "CG22",
        "CG23",
        "CG24",
        "CG25",
        "CG26"
      ],
      "verification": [
        {
          "id": "cg-c2-anchors",
          "performer": "ft439-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c2dded1aef96d333ce550a709c25b5f4b59f4f03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft439-author-02",
            "digest": "sha256:87736c99b0259d54150566f3c4eee5fa1a4f65bf8b50937e3d2a2070ff252447",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/registry_chunk_chain.go,omit,failed,1,TestChunkChainAnchors,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,^TestChunkChainAnchors$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,504\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,76,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestChunkChainAnchors,\"registry_chunk_chain_test.go:51: tree broken at \\\"A comment-only correction before the green checkpoint stays evidence-only and still takes a refreeze.\\\" in .agents/skills/bench-craft-line/references/bounded-repair-policy.md = [], want \\\"chunk… (325 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "dd48bf924f651ae4d5a695b59d20391bba07b74c",
          "requirement": "anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0,
          "probe": {
            "mutation": "Delete the Forbid row with the needle 'takes a refreeze' from internal/anchors/registry_chunk_chain.go, and record that TestChunkChainAnchors turns red on its mirror row.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft439-author-02",
              "digest": "sha256:87736c99b0259d54150566f3c4eee5fa1a4f65bf8b50937e3d2a2070ff252447",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/registry_chunk_chain.go,omit,failed,1,TestChunkChainAnchors,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,^TestChunkChainAnchors$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,504\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,76,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestChunkChainAnchors,\"registry_chunk_chain_test.go:51: tree broken at \\\"A comment-only correction before the green checkpoint stays evidence-only and still takes a refreeze.\\\" in .agents/skills/bench-craft-line/references/bounded-repair-policy.md = [], want \\\"chunk… (325 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "cg-c2-guidance",
          "performer": "ft439-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c2dded1aef96d333ce550a709c25b5f4b59f4f03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft439-author-02",
            "digest": "sha256:05d1677421f9bd41c6adaf874286864736d2801a3075e91b5c01775d8c09bda3",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.agents/commands/bench-implement-spec.md,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,447\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3841,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: chunk chain: the implement phase restored the refreeze of a comment-only correction\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "dd48bf924f651ae4d5a695b59d20391bba07b74c",
          "requirement": "guidance",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore the sentence 'Each source change except a record commit or a recorded fold moves the chunk tip and takes a refreeze of that chunk.' in the Land section of the implement phase, and record that the docs-currency-workflow check turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft439-author-02",
              "digest": "sha256:05d1677421f9bd41c6adaf874286864736d2801a3075e91b5c01775d8c09bda3",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.agents/commands/bench-implement-spec.md,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,447\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3841,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: chunk chain: the implement phase restored the refreeze of a comment-only correction\",1\nskips[0]{package,test,reason}:\n"
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
  },
  "amendments": [
    {
      "from": "sha256:aa15a1af150c53a83df89c155d1d30f25af15bd1fddf360feba92d60b9f278f0",
      "to": "sha256:d30f2a0b57067b74d9fa2c187538c772cadc87ed52b44b722869b4fe67279ba4",
      "chunk_ids": {
        "CG-C1": [
          "CG-C1"
        ]
      }
    }
  ]
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

## CG-C1 round 1 review

The three axes ran on sonnet / medium over the frozen pair `beb3b302..890370c8`. Each axis passes with zero finding IDs. The repair count of CG-C1 is 0 of 2.

## Standards

Findings: 0. Worst issue: none.

## Spec

Findings: 0. Worst issue: none.

## Coverage

Findings: 0. Worst issue: none.

## CG-C1 advice

Each item is an initial-round concern at severity 2 or lower and is not a correctness defect. Thus the bounded repair policy routes each item to optional advice, and no item takes a repair cycle.

- Standards S1 (severity 2): the `chain base` case of `TestReviewCheckpointKeepsStrictEvidence` repeats five setup calls of `chainedOverCommentGap`. A shared form needs a callback, so the repeated text stays.
- Standards S3 (severity 2): the doc comments of the two new tests cite spec row IDs (`internal/gate/review_checkpoint_commits_test.go:124`, `internal/reviewrecord/fold_test.go:159`). `craft-comments` treats an identifier as provenance. The next edit of those lines removes the row IDs.
- Standards S2 and Spec P1 (severity 1): the spec delta corrects the FT380 criterion citation. The staleness pass made that correction in its enabling plan commit, so no repair target remains.
- Coverage C1 (severity 1): the `covered == later` clause of `digestFollows` duplicates the equal-tree answer of `commentgap.Prove`, and a probe that omits it is silent. The spec decision names both forms, and the clause skips a tree diff for equal digests.
- Spec advice: the `chain base` case writes `correctCommentGap` and then overwrites the same file.
