# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/commitment-delivery-integrity/spec.md",
  "plan_digest": "sha256:e1af045505879acfbe050e31acc668e35464b452e9a3f71f2b38add2733381ef",
  "implementation_session": "",
  "chunks": [
    {
      "id": "FD-C1",
      "base": "9febee8f284bfed4f6de714717f4364adede738e",
      "tip": "29eabd3090ad23b7e6dc359a1922029df50cab79",
      "plan_digest": "sha256:e1af045505879acfbe050e31acc668e35464b452e9a3f71f2b38add2733381ef",
      "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
      "acceptance_rows": [
        "FD1",
        "FD2",
        "FD3",
        "FD4",
        "FD5",
        "FD9",
        "FD23",
        "FD24",
        "FD26",
        "FD6",
        "FD7",
        "FD8",
        "FD10",
        "FD11"
      ],
      "verification": [
        {
          "id": "t1-commitment-v1",
          "performer": "claude:fd_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t1",
            "digest": "sha256:e4b3b2c2bdbb56cbcdb24dbb609e2a57475407d2d11389434ec9c3b06da8cd74",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,29eabd3090ad23b7e6dc359a1922029df50cab79,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,6727\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --omit '<refuseObligationFree call and its return>' --package ./internal/commitment --run TestCommitmentPlanRefusesObligationFreeBinding\nexit: 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/authority.go,omit,failed,3,yes\nfailures[3]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/changed_legacy_binding,\"authority_test.go:104: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/inactive_milestone,\"authority_test.go:104: BuildPlan() = <nil>, want the refusal of outcome \\\"B\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/new_binding,\"authority_test.go:104: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n"
          },
          "requirement": "t1-commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0,
          "probe": {
            "mutation": "In BuildPlan, skip the refusal of a new obligation-free binding. TestCommitmentPlanRefusesObligationFreeBinding must fail and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:fd_t1",
              "digest": "sha256:e4b3b2c2bdbb56cbcdb24dbb609e2a57475407d2d11389434ec9c3b06da8cd74",
              "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,29eabd3090ad23b7e6dc359a1922029df50cab79,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,6727\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --omit '<refuseObligationFree call and its return>' --package ./internal/commitment --run TestCommitmentPlanRefusesObligationFreeBinding\nexit: 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/authority.go,omit,failed,3,yes\nfailures[3]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/changed_legacy_binding,\"authority_test.go:104: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/inactive_milestone,\"authority_test.go:104: BuildPlan() = <nil>, want the refusal of outcome \\\"B\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/new_binding,\"authority_test.go:104: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n"
            }
          }
        },
        {
          "id": "t1-conformance-v1",
          "performer": "claude:fd_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t1",
            "digest": "sha256:a36d8bf55636a03e2ceaa973224d5e200ab0ce0af8edfe0c699552cab6eece8d",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/conformance\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,29eabd3090ad23b7e6dc359a1922029df50cab79,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38257\nfailures[0]{package,test,line}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceProseBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestSkillDescriptionBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem\"\n"
          },
          "requirement": "t1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-bench-v1",
          "performer": "claude:fd_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t1",
            "digest": "sha256:6830587c1ff5cd04da23b13c62943da06a71dcf4947dfd00ae6cf28a36cb7f5a",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./cmd/bench\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,29eabd3090ad23b7e6dc359a1922029df50cab79,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,12690\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t2-commitment-v1",
          "performer": "claude:fd_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t2",
            "digest": "sha256:d4ff6d7729d892fb16fdcf9a3c932a10ccb8ac1b67768fb5714b6a9460324565",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,5334\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n$ bench worktree exec FT390 -- bench probe internal/commitment/delivery.go --swap 'return policy, nil, fmt.Errorf(\"deliverable %q of outcome %q names no obligation; list the sources that it satisfies with bench commitment plan --input <file>\", path, outcome.ID)' --with 'return policy, nil, nil' --package ./internal/commitment --run 'TestCommitmentDeliverRowless$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/delivery.go,swap,failed,1,yes\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentDeliverRowless,\"delivery_test.go:96: obligation-free Deliver = <nil>, want the names-no-obligation refusal for an outcome with sources\"\n"
          },
          "requirement": "t2-commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0,
          "probe": {
            "mutation": "In Deliver, return the policy unchanged with no error for an obligation-free binding of an outcome with sources. TestCommitmentDeliverRowless must fail and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:fd_t2",
              "digest": "sha256:d4ff6d7729d892fb16fdcf9a3c932a10ccb8ac1b67768fb5714b6a9460324565",
              "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,5334\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n$ bench worktree exec FT390 -- bench probe internal/commitment/delivery.go --swap 'return policy, nil, fmt.Errorf(\"deliverable %q of outcome %q names no obligation; list the sources that it satisfies with bench commitment plan --input <file>\", path, outcome.ID)' --with 'return policy, nil, nil' --package ./internal/commitment --run 'TestCommitmentDeliverRowless$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/delivery.go,swap,failed,1,yes\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentDeliverRowless,\"delivery_test.go:96: obligation-free Deliver = <nil>, want the names-no-obligation refusal for an outcome with sources\"\n"
            }
          }
        },
        {
          "id": "t2-commitment-repository-v1",
          "performer": "claude:fd_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t2",
            "digest": "sha256:31ee2945b0a885ade8f4ef94885206823d33bb2d5a088058a3aaf0ef4437a72e",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment/repository\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,1753\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-commitment-repository",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t2-worktree-v1",
          "performer": "claude:fd_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t2",
            "digest": "sha256:5bf3be5fe6170fb6e0784a192515a74b731cb51d73ebb1b79d21d3528780a7ad",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/worktree\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,64619\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"capability: fifo: unix sockets unavailable\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable\"\n"
          },
          "requirement": "t2-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t2-conformance-v1",
          "performer": "claude:fd_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t2",
            "digest": "sha256:bd91d661fdd22fbff5c42797bcf6628c0261b772aadd519ab223903ca4ae005f",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/conformance\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38130\nfailures[0]{package,test,line}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceProseBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestSkillDescriptionBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem\"\n"
          },
          "requirement": "t2-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t2-bench-v1",
          "performer": "claude:fd_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t2",
            "digest": "sha256:bedeef4dce9aa1b40299bbfc5a3eb3293c2e4bf8699aa84cc16c84e3db2e7a9a",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./cmd/bench\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,13016\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
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
