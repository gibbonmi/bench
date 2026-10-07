# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/complete-checkpoint-closure/spec.md",
  "plan_digest": "sha256:252e69174ce3424e55e73d397fdecf3a4a7df3ca029e597de1edea0e8d5661ac",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CC1",
      "base": "ea3e7867a6019c2ce11edbf16bae033aaba81a2f",
      "tip": "1ee00ddfe063f207281768d2f5757b923ca95102",
      "plan_digest": "sha256:252e69174ce3424e55e73d397fdecf3a4a7df3ca029e597de1edea0e8d5661ac",
      "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
      "acceptance_rows": [
        "CC01",
        "CC02",
        "CC03",
        "CC04",
        "CC05",
        "CC06",
        "CC07",
        "CC08",
        "CC09",
        "CC10",
        "CC11",
        "CC12",
        "CC13",
        "CC14",
        "CC15",
        "CC16",
        "CC17",
        "CC18",
        "CC19",
        "CC20",
        "CC21",
        "CC22",
        "CC23",
        "CC24",
        "CC25",
        "CC26",
        "CC27",
        "CC28",
        "CC29",
        "CC30",
        "CC31",
        "CC32"
      ],
      "verification": [
        {
          "id": "v-t1-gate-package",
          "performer": "claude:ft392_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t1",
            "digest": "sha256:76d6033109f7ce2cfe9c077b29a64d3883dbbcdfababc84dad276cfec49335bb",
            "excerpt": "tree[1]{target,head,dirty}:\n  complete-checkpoint-closure,1ee00ddfe063f207281768d2f5757b923ca95102,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,24983,397\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-gate-package",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "v-t2-gate-package",
          "performer": "claude:ft392_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2",
            "digest": "sha256:7ac161b97b6746f35429d62302b610fa68e8ce5d49f27537e233949e68034d64",
            "excerpt": "tip 1ee00ddfe063f207281768d2f5757b923ca95102\nbench test --package ./internal/gate\ngithub.com/gibbonmi/bench/internal/gate,pass,25123,397\nfailures[0] skips[0]\n"
          },
          "requirement": "t2-gate-package",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "v-t2-complete-checkpoint",
          "performer": "claude:ft392_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2",
            "digest": "sha256:3aa7b3ff06a609fcbba322183cd9f9a01dc9719b770f181bce32ca33748a3c5e",
            "excerpt": "tip 1ee00ddfe063f207281768d2f5757b923ca95102\nbench test --package ./internal/gate --run 'TestCompleteCheckpoint|TestChunkCheckpointGradesTheCheckoutTree|TestReviewCheckpoint|TestCommitmentExactTransform'\ngithub.com/gibbonmi/bench/internal/gate,pass,6813,66\nfailures[0] skips[0]\n"
          },
          "requirement": "t2-complete-checkpoint",
          "command": "bench test --package ./internal/gate --run 'TestCompleteCheckpoint|TestChunkCheckpointGradesTheCheckoutTree|TestReviewCheckpoint|TestCommitmentExactTransform'",
          "exit_code": 0
        },
        {
          "id": "v-t2-public-route",
          "performer": "claude:ft392_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2",
            "digest": "sha256:21faef66b2ddfc9f2fcfc42c1f39da6c5bfc62da574e17bed7a73581f665d5b4",
            "excerpt": "tip 1ee00ddfe063f207281768d2f5757b923ca95102\nbench test --package ./cmd/bench --run 'TestGateCheckpointRoute'\ngithub.com/gibbonmi/bench/cmd/bench,pass,409,1\nfailures[0] skips[0]\n"
          },
          "requirement": "t2-public-route",
          "command": "bench test --package ./cmd/bench --run 'TestGateCheckpointRoute'",
          "exit_code": 0
        },
        {
          "id": "v-t2-landing-journey",
          "performer": "claude:ft392_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2",
            "digest": "sha256:a8b4bbe0bf2236cd90f7d61a77862631918b6b85ee4d42b8fe8963f05b72f99e",
            "excerpt": "tip 1ee00ddfe063f207281768d2f5757b923ca95102\nbench test --package ./internal/worktree --run 'TestLandCommandPublicRealGitJourney'\ngithub.com/gibbonmi/bench/internal/worktree,pass,4735,7\nfailures[0] skips[0]\n"
          },
          "requirement": "t2-landing-journey",
          "command": "bench test --package ./internal/worktree --run 'TestLandCommandPublicRealGitJourney'",
          "exit_code": 0
        },
        {
          "id": "v-t2-route-omission-proof",
          "performer": "claude:ft392_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2",
            "digest": "sha256:f6f57e06c851c403de97a07d001ea9b92a494de96aa363927075521866dfea42",
            "excerpt": "tip 1ee00ddfe063f207281768d2f5757b923ca95102\nbench test --package ./internal/gate --run 'TestCompleteCheckpoint'\ngithub.com/gibbonmi/bench/internal/gate,pass,2406,11\nprobe: internal/gate/gate.go swap 'if checkpoint.Complete {' -> 'if false {' (route --complete through the ordinary checkout evaluation)\nprobe[1]: bit,internal/gate/gate.go,swap,failed,9,yes\nfailed: TestCompleteCheckpointGradesTheClosedTree complete_checkpoint_test.go:105 complete checkpoint on a source whose delivery closure is red = (0, \"\"), want exit 9\n"
          },
          "requirement": "t2-route-omission-proof",
          "command": "bench test --package ./internal/gate --run 'TestCompleteCheckpoint'",
          "exit_code": 0,
          "probe": {
            "mutation": "Route --complete through the ordinary checkout evaluation. TestCompleteCheckpointGradesTheClosedTree must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft392_t2",
              "digest": "sha256:f6f57e06c851c403de97a07d001ea9b92a494de96aa363927075521866dfea42",
              "excerpt": "tip 1ee00ddfe063f207281768d2f5757b923ca95102\nbench test --package ./internal/gate --run 'TestCompleteCheckpoint'\ngithub.com/gibbonmi/bench/internal/gate,pass,2406,11\nprobe: internal/gate/gate.go swap 'if checkpoint.Complete {' -> 'if false {' (route --complete through the ordinary checkout evaluation)\nprobe[1]: bit,internal/gate/gate.go,swap,failed,9,yes\nfailed: TestCompleteCheckpointGradesTheClosedTree complete_checkpoint_test.go:105 complete checkpoint on a source whose delivery closure is red = (0, \"\"), want exit 9\n"
            }
          }
        },
        {
          "id": "v-t2-clean-checkout-proof",
          "performer": "claude:ft392_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2",
            "digest": "sha256:1c96f5e360d39e9e15259768486a1dd5912d6d4b3e58b3ed6e506ff9eafc8887",
            "excerpt": "tip 1ee00ddfe063f207281768d2f5757b923ca95102\nbench test --package ./internal/gate --run 'TestCompleteCheckpointRefusesADirtyCheckout'\ngithub.com/gibbonmi/bench/internal/gate,pass,223,5\nprobe 1 (recorded): complete_checkpoint.go swap 'if working != sourceTree {' -> 'if working == \"\" && working != sourceTree {' (remove the clean-checkout refusal)\nprobe[1]: bit,internal/gate/complete_checkpoint.go,swap,failed,4,yes\nfailed: tracked_edit, untracked_file (exit 0 and oracle ran), uncommitted_record, stale_evidence (stale-evidence reason instead of the clean-checkout refusal)\nprobe 2 (exempt the review record from the refusal): bit,internal/gate/complete_checkpoint.go,swap,failed,1,yes\nfailed: uncommitted_record (reached 'completion is incomplete or stale')\n"
          },
          "requirement": "t2-clean-checkout-proof",
          "command": "bench test --package ./internal/gate --run 'TestCompleteCheckpointRefusesADirtyCheckout'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the clean-checkout refusal, or exempt the review record from it. Each named refusal assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft392_t2",
              "digest": "sha256:1c96f5e360d39e9e15259768486a1dd5912d6d4b3e58b3ed6e506ff9eafc8887",
              "excerpt": "tip 1ee00ddfe063f207281768d2f5757b923ca95102\nbench test --package ./internal/gate --run 'TestCompleteCheckpointRefusesADirtyCheckout'\ngithub.com/gibbonmi/bench/internal/gate,pass,223,5\nprobe 1 (recorded): complete_checkpoint.go swap 'if working != sourceTree {' -> 'if working == \"\" && working != sourceTree {' (remove the clean-checkout refusal)\nprobe[1]: bit,internal/gate/complete_checkpoint.go,swap,failed,4,yes\nfailed: tracked_edit, untracked_file (exit 0 and oracle ran), uncommitted_record, stale_evidence (stale-evidence reason instead of the clean-checkout refusal)\nprobe 2 (exempt the review record from the refusal): bit,internal/gate/complete_checkpoint.go,swap,failed,1,yes\nfailed: uncommitted_record (reached 'completion is incomplete or stale')\n"
            }
          }
        },
        {
          "id": "v-t2-witness-proof",
          "performer": "claude:ft392_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2",
            "digest": "sha256:436b59043cf5e06ebd8757784d787f07b812a515cc4013d4fc7d2ba5181abc50",
            "excerpt": "tip 1ee00ddfe063f207281768d2f5757b923ca95102\nbench test --package ./internal/gate --run 'TestReviewCheckpointReuse|TestCompleteCheckpointEvidenceNamesThePublishedTree'\ngithub.com/gibbonmi/bench/internal/gate,pass,1130,2\nprobe: internal/gate/run_outcomes_test.go one swap: fixture witnesses checkout-relative and outcomeWitness.path returns the checkout root\nprobe[1]: bit,internal/gate/run_outcomes_test.go,swap,failed,2,yes\nfailed: TestCompleteCheckpointEvidenceNamesThePublishedTree complete_checkpoint_test.go:136 (.gate-record-during absent); TestReviewCheckpointReuse review_checkpoint_test.go:101 (2 runs)\n"
          },
          "requirement": "t2-witness-proof",
          "command": "bench test --package ./internal/gate --run 'TestReviewCheckpointReuse|TestCompleteCheckpointEvidenceNamesThePublishedTree'",
          "exit_code": 0,
          "probe": {
            "mutation": "Point the witness helper and the fixture witnesses at the checkout. The CC25 and CC32 assertions must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft392_t2",
              "digest": "sha256:436b59043cf5e06ebd8757784d787f07b812a515cc4013d4fc7d2ba5181abc50",
              "excerpt": "tip 1ee00ddfe063f207281768d2f5757b923ca95102\nbench test --package ./internal/gate --run 'TestReviewCheckpointReuse|TestCompleteCheckpointEvidenceNamesThePublishedTree'\ngithub.com/gibbonmi/bench/internal/gate,pass,1130,2\nprobe: internal/gate/run_outcomes_test.go one swap: fixture witnesses checkout-relative and outcomeWitness.path returns the checkout root\nprobe[1]: bit,internal/gate/run_outcomes_test.go,swap,failed,2,yes\nfailed: TestCompleteCheckpointEvidenceNamesThePublishedTree complete_checkpoint_test.go:136 (.gate-record-during absent); TestReviewCheckpointReuse review_checkpoint_test.go:101 (2 runs)\n"
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
