# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/complete-checkpoint-closure/spec.md",
  "plan_digest": "sha256:0cd15d2e700fa2603ac87fe1df65b5eb74a88ecc1bd5ea045a287b45c763e326",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CC1",
      "base": "ea3e7867a6019c2ce11edbf16bae033aaba81a2f",
      "tip": "c1c76ff057669cb9407cbcb1b02ada1000ce926c",
      "plan_digest": "sha256:0cd15d2e700fa2603ac87fe1df65b5eb74a88ecc1bd5ea045a287b45c763e326",
      "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
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
        },
        {
          "id": "v-t1-gate-package-r1",
          "performer": "claude:ft392_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t1",
            "digest": "sha256:c39625233142325663b784d473391c1a5b6005524105e44256cbdf28985a2c64",
            "excerpt": "tree[1]{target,head,dirty}:\n  complete-checkpoint-closure,c1c76ff057669cb9407cbcb1b02ada1000ce926c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,25125,398\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-gate-package",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "v-t2-gate-package-r1",
          "performer": "claude:ft392_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2_repair1",
            "digest": "sha256:fb09845e594e164ff862b6410e1ed14d749446da47a500cb4ab5f2a327dd85ab",
            "excerpt": "tip c1c76ff057669cb9407cbcb1b02ada1000ce926c\nbench test --package ./internal/gate\ngithub.com/gibbonmi/bench/internal/gate,pass,25074,398\n"
          },
          "requirement": "t2-gate-package",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "v-t2-complete-checkpoint-r1",
          "performer": "claude:ft392_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2_repair1",
            "digest": "sha256:e722108e8e58524480640ac513a8e008e548f07f2cee96091037ea664912f09b",
            "excerpt": "tip c1c76ff057669cb9407cbcb1b02ada1000ce926c\nbench test --package ./internal/gate --run 'TestCompleteCheckpoint|TestChunkCheckpointGradesTheCheckoutTree|TestReviewCheckpoint|TestCommitmentExactTransform'\ngithub.com/gibbonmi/bench/internal/gate,pass,6949,67\n"
          },
          "requirement": "t2-complete-checkpoint",
          "command": "bench test --package ./internal/gate --run 'TestCompleteCheckpoint|TestChunkCheckpointGradesTheCheckoutTree|TestReviewCheckpoint|TestCommitmentExactTransform'",
          "exit_code": 0
        },
        {
          "id": "v-t2-public-route-r1",
          "performer": "claude:ft392_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2_repair1",
            "digest": "sha256:78deeae3e873489c9731109960a79a5ea51ef313ec02bc129852aef8bdeb2c3f",
            "excerpt": "tip c1c76ff057669cb9407cbcb1b02ada1000ce926c\nbench test --package ./cmd/bench --run 'TestGateCheckpointRoute'\ngithub.com/gibbonmi/bench/cmd/bench,pass,411,1\n"
          },
          "requirement": "t2-public-route",
          "command": "bench test --package ./cmd/bench --run 'TestGateCheckpointRoute'",
          "exit_code": 0
        },
        {
          "id": "v-t2-landing-journey-r1",
          "performer": "claude:ft392_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2_repair1",
            "digest": "sha256:f552ac44438c8be54ba4b3fd4863e7fb208f7445124f0202ab0e22fd58dd0eaf",
            "excerpt": "tip c1c76ff057669cb9407cbcb1b02ada1000ce926c\nbench test --package ./internal/worktree --run 'TestLandCommandPublicRealGitJourney'\ngithub.com/gibbonmi/bench/internal/worktree,pass,4723,7\n"
          },
          "requirement": "t2-landing-journey",
          "command": "bench test --package ./internal/worktree --run 'TestLandCommandPublicRealGitJourney'",
          "exit_code": 0
        },
        {
          "id": "v-t2-route-omission-proof-r1",
          "performer": "claude:ft392_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2_repair1",
            "digest": "sha256:621bb84dd354e7db0c32e3e6c002f70fce7d42b53c7631363cb7df0c005a5660",
            "excerpt": "tip c1c76ff057669cb9407cbcb1b02ada1000ce926c\nbench test --package ./internal/gate --run 'TestCompleteCheckpoint'\ngithub.com/gibbonmi/bench/internal/gate,pass,2457,12\nprobe: internal/gate/gate.go swap 'if checkpoint.Complete {' -> 'if false {' (route --complete through the ordinary checkout evaluation)\nprobe[1]: bit,internal/gate/gate.go,swap,failed,9,yes\nmutated run: github.com/gibbonmi/bench/internal/gate,fail,1738,12\nfailed: TestCompleteCheckpointGradesTheClosedTree complete_checkpoint_test.go:105 (exit 0, want exit 9 and the citation), plus 8 other TestCompleteCheckpoint cases\nprobe exit code 1 is the failing go test exit, not a probe output field\n"
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
              "ref": "claude-agent:ft392_t2_repair1",
              "digest": "sha256:621bb84dd354e7db0c32e3e6c002f70fce7d42b53c7631363cb7df0c005a5660",
              "excerpt": "tip c1c76ff057669cb9407cbcb1b02ada1000ce926c\nbench test --package ./internal/gate --run 'TestCompleteCheckpoint'\ngithub.com/gibbonmi/bench/internal/gate,pass,2457,12\nprobe: internal/gate/gate.go swap 'if checkpoint.Complete {' -> 'if false {' (route --complete through the ordinary checkout evaluation)\nprobe[1]: bit,internal/gate/gate.go,swap,failed,9,yes\nmutated run: github.com/gibbonmi/bench/internal/gate,fail,1738,12\nfailed: TestCompleteCheckpointGradesTheClosedTree complete_checkpoint_test.go:105 (exit 0, want exit 9 and the citation), plus 8 other TestCompleteCheckpoint cases\nprobe exit code 1 is the failing go test exit, not a probe output field\n"
            }
          }
        },
        {
          "id": "v-t2-clean-checkout-proof-r1",
          "performer": "claude:ft392_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2_repair1",
            "digest": "sha256:b9db3ddeaa45e2c3e547300ca87f12089aee32a7ba8e0a818f9f798bd0e566bb",
            "excerpt": "tip c1c76ff057669cb9407cbcb1b02ada1000ce926c\nbench test --package ./internal/gate --run 'TestCompleteCheckpointRefusesADirtyCheckout'\ngithub.com/gibbonmi/bench/internal/gate,pass,411,6\nprobe: internal/gate/complete_checkpoint.go swap 'if working != sourceTree {' -> 'if working == \"\" && working != sourceTree {' (remove the clean-checkout refusal)\nprobe[1]: bit,internal/gate/complete_checkpoint.go,swap,failed,4,yes\nfailed: tracked_edit, untracked_file (exit 0); uncommitted_record, stale_evidence ('completion is incomplete or stale' instead of the clean-checkout refusal)\nrecord-exemption probe not rerun at this tip\nprobe exit code 1 is the failing go test exit, not a probe output field\n"
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
              "ref": "claude-agent:ft392_t2_repair1",
              "digest": "sha256:b9db3ddeaa45e2c3e547300ca87f12089aee32a7ba8e0a818f9f798bd0e566bb",
              "excerpt": "tip c1c76ff057669cb9407cbcb1b02ada1000ce926c\nbench test --package ./internal/gate --run 'TestCompleteCheckpointRefusesADirtyCheckout'\ngithub.com/gibbonmi/bench/internal/gate,pass,411,6\nprobe: internal/gate/complete_checkpoint.go swap 'if working != sourceTree {' -> 'if working == \"\" && working != sourceTree {' (remove the clean-checkout refusal)\nprobe[1]: bit,internal/gate/complete_checkpoint.go,swap,failed,4,yes\nfailed: tracked_edit, untracked_file (exit 0); uncommitted_record, stale_evidence ('completion is incomplete or stale' instead of the clean-checkout refusal)\nrecord-exemption probe not rerun at this tip\nprobe exit code 1 is the failing go test exit, not a probe output field\n"
            }
          }
        },
        {
          "id": "v-t2-witness-proof-r1",
          "performer": "claude:ft392_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_t2_repair1",
            "digest": "sha256:4956323053c55f3ea8d7cc6bc0e186db9a8567631cc3e5904317e2dcd352b0c8",
            "excerpt": "tip c1c76ff057669cb9407cbcb1b02ada1000ce926c\nbench test --package ./internal/gate --run 'TestReviewCheckpointReuse|TestCompleteCheckpointEvidenceNamesThePublishedTree'\ngithub.com/gibbonmi/bench/internal/gate,pass,1100,2\nprobe: internal/gate/run_outcomes_test.go one swap: fixture witnesses checkout-relative and outcomeWitness.path returns the checkout root\nprobe[1]: bit,internal/gate/run_outcomes_test.go,swap,failed,2,yes\nmutated run: github.com/gibbonmi/bench/internal/gate,fail,705,2\nfailed: TestCompleteCheckpointEvidenceNamesThePublishedTree complete_checkpoint_test.go:136 (.gate-record-during absent); TestReviewCheckpointReuse review_checkpoint_test.go:101 (2 runs)\nprobe exit code 1 is the failing go test exit, not a probe output field\n"
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
              "ref": "claude-agent:ft392_t2_repair1",
              "digest": "sha256:4956323053c55f3ea8d7cc6bc0e186db9a8567631cc3e5904317e2dcd352b0c8",
              "excerpt": "tip c1c76ff057669cb9407cbcb1b02ada1000ce926c\nbench test --package ./internal/gate --run 'TestReviewCheckpointReuse|TestCompleteCheckpointEvidenceNamesThePublishedTree'\ngithub.com/gibbonmi/bench/internal/gate,pass,1100,2\nprobe: internal/gate/run_outcomes_test.go one swap: fixture witnesses checkout-relative and outcomeWitness.path returns the checkout root\nprobe[1]: bit,internal/gate/run_outcomes_test.go,swap,failed,2,yes\nmutated run: github.com/gibbonmi/bench/internal/gate,fail,705,2\nfailed: TestCompleteCheckpointEvidenceNamesThePublishedTree complete_checkpoint_test.go:136 (.gate-record-during absent); TestReviewCheckpointReuse review_checkpoint_test.go:101 (2 runs)\nprobe exit code 1 is the failing go test exit, not a probe output field\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "r-cc1-standards",
          "performer": "claude:ft392_review_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft392_review_standards",
            "digest": "sha256:a1d042a67ed8cda45aea52613f1ad9823bb21c0a34d0abc35a334b672a37a553",
            "excerpt": "Standards axis, chunk CC1, evidence sha256:793a47b2 current=true at 5f13265d. 3 findings; worst S1.\nS1 complete_checkpoint_test.go:198 respells the cleanCheckoutRefusal literal from complete_checkpoint.go:12; no recorded red needs the independence. auto-fix, confidence 6.\nS2 complete_checkpoint.go:17-18 comment claims the checkpoint and the landing cannot grade two different trees; landing.go:241 transforms the composition tree. auto-fix, confidence 6.\nS3 published.Tree plus WithCompletion pairing is hand-built at landing.go:241/249 and complete_checkpoint.go:35/39. ask-user, confidence 4.\n"
          },
          "axis": "Standards",
          "base": "ea3e7867a6019c2ce11edbf16bae033aaba81a2f",
          "tip": "1ee00ddfe063f207281768d2f5757b923ca95102",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "r-cc1-spec",
          "performer": "claude:ft392_review_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_review_spec",
            "digest": "sha256:798ed75f19dc6e665eac1a12abd9aa1bf36edfe54128acf8f913b8f72ae1f4c1",
            "excerpt": "Spec axis, chunk CC1, evidence sha256:793a47b2 current=true at 5f13265d. 0 findings.\nCC01-CC32 met. Clean refusal at complete_checkpoint.go:32 precedes compose (:35), grade (:39), and the stale-evidence refusal (checkpoint.go:145). gate.go:204 keeps notifyGateSignals and mode; engine.go:46 keeps ExecuteTree arguments; BENCH-reference.md:400 carries the CC28 sentence.\n"
          },
          "axis": "Spec",
          "base": "ea3e7867a6019c2ce11edbf16bae033aaba81a2f",
          "tip": "1ee00ddfe063f207281768d2f5757b923ca95102",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "r-cc1-coverage",
          "performer": "claude:ft392_review_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "584eac658ba71f284cb17b3683d2b03d08f6bf33",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft392_review_coverage",
            "digest": "sha256:a1712ffc9ae9a244c13a200cca879ae184f9b4c32305292c1c39df65e30a82aa",
            "excerpt": "Coverage axis, chunk CC1, evidence sha256:793a47b2 current=true at 5f13265d. 1 finding; worst C1.\nC1 complete_checkpoint.go:63-69 clean check: no test pins that an ignored file passes (spec line 94). Probe refusing on ignored files: silent on ./internal/gate (397 tests) and on cmd/bench TestGateCheckpointRoute, restored yes. auto-fix, confidence 7.\nCC04 fixture probe red at complete_checkpoint_test.go:100; CC08 red at :173; CC10 red at :190; CC31 red at :145; all restored yes.\n"
          },
          "axis": "Coverage",
          "base": "ea3e7867a6019c2ce11edbf16bae033aaba81a2f",
          "tip": "1ee00ddfe063f207281768d2f5757b923ca95102",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "r-cc1-standards-confirm",
          "performer": "claude:ft392_confirm_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_confirm_standards",
            "digest": "sha256:50b02d5b550de22c86dd19956c0ca75f9028a589e185b7362e85c8c055d6b894",
            "excerpt": "Standards confirming round, chunk CC1, evidence sha256:e53fc44a current=true at ee5f1724; repair delta 00c12fef..c1c76ff0. 0 findings.\nS1 confirmed: complete_checkpoint_test.go:241-242 asserts on cleanCheckoutRefusal (complete_checkpoint.go:12). S2 confirmed: complete_checkpoint.go:16-18 states only what the code does. C1 confirmed: complete_checkpoint_test.go:203-209.\nAdvice (no ID): the absence check at complete_checkpoint_test.go:241 spells \"completion is incomplete or stale\", which reviewrecord/check.go:62 owns inline.\n"
          },
          "axis": "Standards",
          "base": "ea3e7867a6019c2ce11edbf16bae033aaba81a2f",
          "tip": "c1c76ff057669cb9407cbcb1b02ada1000ce926c",
          "finding_ids": [],
          "supersedes": [
            "r-cc1-standards"
          ]
        },
        {
          "id": "r-cc1-spec-confirm",
          "performer": "claude:ft392_confirm_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_confirm_spec",
            "digest": "sha256:abc0290f304195e840afa131e47dcd6b435a578f819c3fd1258da07f1a47112d",
            "excerpt": "Spec confirming round, chunk CC1, evidence sha256:e53fc44a current=true at ee5f1724; repair delta 00c12fef..c1c76ff0. 0 findings.\nS1, S2, C1 confirmed. The repair changes only comment lines 16-18 of complete_checkpoint.go; the ignored-file case matches spec.md:94; CC01-CC32 do not regress.\n"
          },
          "axis": "Spec",
          "base": "ea3e7867a6019c2ce11edbf16bae033aaba81a2f",
          "tip": "c1c76ff057669cb9407cbcb1b02ada1000ce926c",
          "finding_ids": [],
          "supersedes": [
            "r-cc1-spec"
          ]
        },
        {
          "id": "r-cc1-coverage-confirm",
          "performer": "claude:ft392_confirm_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft392_confirm_coverage",
            "digest": "sha256:400a8c073c6968b21fe1423587f48e9fe5e829246f2e5265904c8904fd9b2e7a",
            "excerpt": "Coverage confirming round, chunk CC1, evidence sha256:e53fc44a current=true at ee5f1724; repair delta 00c12fef..c1c76ff0. 0 findings.\nC1 probe (refuse on ignored files): bit, ignored_file failed at complete_checkpoint_test.go:207, restored yes.\nS1 probe (never refuse): bit, 4 failures at complete_checkpoint_test.go:242, restored yes. git status clean at finish.\n"
          },
          "axis": "Coverage",
          "base": "ea3e7867a6019c2ce11edbf16bae033aaba81a2f",
          "tip": "c1c76ff057669cb9407cbcb1b02ada1000ce926c",
          "finding_ids": [],
          "supersedes": [
            "r-cc1-coverage"
          ]
        }
      ]
    }
  ],
  "completion": {
    "state": "completed",
    "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
    "performer": "claude:ft392-orchestrator-20261007",
    "reconciliation": {
      "CC01": "covered",
      "CC02": "covered",
      "CC03": "covered",
      "CC04": "covered",
      "CC05": "covered",
      "CC06": "covered",
      "CC07": "covered",
      "CC08": "covered",
      "CC09": "covered",
      "CC10": "covered",
      "CC11": "covered",
      "CC12": "covered",
      "CC13": "covered",
      "CC14": "covered",
      "CC15": "covered",
      "CC16": "covered",
      "CC17": "covered",
      "CC18": "covered",
      "CC19": "covered",
      "CC20": "covered",
      "CC21": "covered",
      "CC22": "covered",
      "CC23": "covered",
      "CC24": "covered",
      "CC25": "covered",
      "CC26": "covered",
      "CC27": "covered",
      "CC28": "covered",
      "CC29": "covered",
      "CC30": "covered",
      "CC31": "covered",
      "CC32": "covered"
    },
    "verification": [
      {
        "id": "v-final-ordinary-integration",
        "performer": "claude:ft392-orchestrator-20261007",
        "role": "integration-verification",
        "model": "opus",
        "effort": "medium",
        "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude-session:ft392-orchestrator-20261007",
          "digest": "sha256:77e15cc720d0c43a9db43ea1c4af35f4b0a92a61bac98109123f28e4421a6b2b",
          "excerpt": "bench test --package ./... at d6d3a6df: exit 0\ntree[1]{target,head,dirty}:\n  complete-checkpoint-closure,d6d3a6df89a64efd3eb66fc43eac76d9e09c355b,false\npackages[122]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench,pass,2,22\npackages[122] pass; capability skips only (fifo socket, privilege device)\n"
        },
        "requirement": "ordinary-integration",
        "command": "bench test --package ./...",
        "exit_code": 0
      },
      {
        "id": "v-final-coverage",
        "performer": "claude:ft392-orchestrator-20261007",
        "role": "integration-verification",
        "model": "opus",
        "effort": "medium",
        "source_digest": "e72ea3f615d2c12f8f810c7e5e5159cdf3a0a396",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude-session:ft392-orchestrator-20261007",
          "digest": "sha256:ffcd6bc8e819206611422f6453941435c4fa5c679f8a3f627553f4de2be885dd",
          "excerpt": "bench coverage --check complete-checkpoint-closure at d6d3a6df: exit 0\nok: coverage map valid — 32 row(s)\nuncited: 25 row(s) with no seam-cell citation — CC01, CC02, CC03, CC04, CC05, CC06, CC07, CC08, CC09, CC10, CC11, CC12, CC13, CC14, CC15, CC16, CC17, CC20, CC21, CC24, CC28, CC29, CC30, CC31, CC32\n"
        },
        "requirement": "coverage",
        "command": "bench coverage --check complete-checkpoint-closure",
        "exit_code": 0
      }
    ]
  },
  "amendments": [
    {
      "from": "sha256:252e69174ce3424e55e73d397fdecf3a4a7df3ca029e597de1edea0e8d5661ac",
      "to": "sha256:0cd15d2e700fa2603ac87fe1df65b5eb74a88ecc1bd5ea045a287b45c763e326",
      "chunk_ids": {
        "CC1": [
          "CC1"
        ]
      }
    }
  ]
}
```

## Standards

Chunk CC1 has 3 findings. The worst finding is S1.

- S1 (auto-fix, confidence 6): `internal/gate/complete_checkpoint_test.go:198` spells the clean-checkout refusal text again. `internal/gate/complete_checkpoint.go:12` owns that text as `cleanCheckoutRefusal`. No recorded red needs an independent spelling, so the test must use the constant.
- S2 (auto-fix, confidence 6): the comment at `internal/gate/complete_checkpoint.go:17-18` says that the checkpoint and the landing cannot grade two different trees. The landing transforms its composition tree at `internal/landing/landing.go:241`, so the claim is false when the destination moves. Remove the claim.
- S3 (no-op, confidence 4): the pair of `published.Tree` and `WithCompletion` is built at `internal/landing/landing.go:241-249` and at `internal/gate/complete_checkpoint.go:35-39`. The closed reviewer decision keeps the landing unchanged, and `published.Tree` stays the one closure derivation. So no repair target remains in this spec.

## Spec

Chunk CC1 has 0 findings. Rows CC01 to CC32 are met.

## Coverage

Chunk CC1 has 1 finding. The worst finding is C1.

- C1 (auto-fix, confidence 7): no test makes sure that an ignored file passes the clean-checkout check at `internal/gate/complete_checkpoint.go:63-69`. The spec decides that the working-tree hash excludes ignored files. A probe that refuses on each ignored file stayed silent on `./internal/gate` and on `TestGateCheckpointRoute`. Add an ignored-file case to `TestCompleteCheckpointRefusesADirtyCheckout` that expects exit 0.

## Repair state

Chunk CC1 used 1 of its 2 repair cycles and its 1 hardening cycle. The repair session `claude:ft392_t2_repair1` fixed S1, S2, and C1 in one cycle. The confirming round of all three axes read the repair delta and found 0 findings. Each confirming result supersedes the first-round result of its axis.

The record carries the repair plan amendment before the chunk update. An earlier record order put the chunk update first, and the chunk checkpoint refused it. The rebuilt record holds the same entries and excerpt digests.

## Advice

The absence check in `internal/gate/complete_checkpoint_test.go` spells the stale-evidence text, which `internal/reviewrecord/check.go` owns as an inline error. Two older tests use the same literal. This advice has no finding ID and goes to the drain.
