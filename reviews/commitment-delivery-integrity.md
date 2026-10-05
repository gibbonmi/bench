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
