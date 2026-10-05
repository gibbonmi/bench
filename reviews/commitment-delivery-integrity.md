# Review outcomes

## FD-C1 pickup

The FD-C1 review covers the frozen pair `9febee8f..29eabd30`. The Coverage
findings start repair cycle 1 of 2. This cycle is the one hardening cycle of
the chunk, because rows FD1 to FD11, FD23, FD24, and FD26 prove.

## Standards

Findings: 0. Worst issue: none.

Advice: `obligationFreeRoute` in the landing fixture reuses the seed of
`deliveryRoutes[1]` by position.

## Spec

Findings: 0. Worst issue: none. Every closed decision held, and every chunk
row is met.

## Coverage

Findings: 4. Worst issue: C2.

- C1 (`auto-fix`, confidence 8): `internal/commitment/authority.go:127`. No
  test moves a kept legacy binding to another outcome. Probe t3 was silent.
- C2 (`auto-fix`, confidence 9): `internal/commitment/authority.go:129`. No
  test keeps a binding of a rowless outcome that gains sources. Probe t4 was
  silent.
- C3 (`auto-fix`, confidence 8): `internal/commitment/authority.go:109`. No
  test puts a new binding after a kept binding in one outcome. Probe t5 was
  silent.
- C4 (`auto-fix`, confidence 8): `internal/commitment/authority.go:108`. No
  test puts the binding on a later outcome of a milestone. Probe t6 was
  silent.

```bench-review-record
{
  "version": 2,
  "spec": "specs/commitment-delivery-integrity/spec.md",
  "plan_digest": "sha256:37236d7f76612ed55db6b4a27b59899919a497c041af08415e14869512b9d73d",
  "implementation_session": "",
  "chunks": [
    {
      "id": "FD-C1",
      "base": "9febee8f284bfed4f6de714717f4364adede738e",
      "tip": "4e1fb87f9d5378c895e4219e74e339bb87cb1d46",
      "plan_digest": "sha256:37236d7f76612ed55db6b4a27b59899919a497c041af08415e14869512b9d73d",
      "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
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
        },
        {
          "id": "t1-commitment-v2",
          "performer": "claude:fd_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t1_r1",
            "digest": "sha256:c192b6825b1fe0dcf9e9083b1798ebfe0c7d46dafd122d871c1bccbc190d28b6",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,a7414079fc441354e4587d4910ecc710be4935cb,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,6196\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --swap 'if obligationFree(outcome, binding) && !retainedObligationFree(current, outcome.ID, binding) {' --with 'if false && obligationFree(outcome, binding) && !retainedObligationFree(current, outcome.ID, binding) {' --package ./internal/commitment --run TestCommitmentPlanRefusesObligationFreeBinding\nexit: 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/authority.go,swap,failed,7,yes\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/new_binding,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/new_binding_after_a_kept_binding,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"draft\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/rowless_binding_gains_a_row,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/second_outcome,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"B\\\" binding \\\"spec\\\"\"\n\nRepair coverage probes (run at 4e1fb87f content before commit; --run TestCommitmentPlanRefusesObligationFreeBinding):\nC1 --swap 'if outcome.ID == id {' --with 'if outcome.ID == id || len(outcome.Deliverables) > 0 {'\n  bit,internal/commitment/authority.go,swap,failed,2,yes\n  failed: legacy_binding_moved_to_another_outcome, inactive_milestone\nC2 --omit ' && obligationFree(outcome, kept)'\n  bit,internal/commitment/authority.go,omit,failed,1,yes\n  failed: rowless_binding_gains_a_row\nC3 --swap 'range outcome.Deliverables {' --with 'range outcome.Deliverables[:min(1, len(outcome.Deliverables))] {'\n  bit,internal/commitment/authority.go,swap,failed,1,yes\n  failed: new_binding_after_a_kept_binding\nC4 --swap 'range milestone.Outcomes {' --with 'range milestone.Outcomes[:min(1, len(milestone.Outcomes))] {'\n  bit,internal/commitment/authority.go,swap,failed,2,yes\n  failed: legacy_binding_moved_to_another_outcome, second_outcome\n"
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
              "ref": "claude-agent:fd_t1_r1",
              "digest": "sha256:c192b6825b1fe0dcf9e9083b1798ebfe0c7d46dafd122d871c1bccbc190d28b6",
              "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,a7414079fc441354e4587d4910ecc710be4935cb,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,6196\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --swap 'if obligationFree(outcome, binding) && !retainedObligationFree(current, outcome.ID, binding) {' --with 'if false && obligationFree(outcome, binding) && !retainedObligationFree(current, outcome.ID, binding) {' --package ./internal/commitment --run TestCommitmentPlanRefusesObligationFreeBinding\nexit: 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/authority.go,swap,failed,7,yes\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/new_binding,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/new_binding_after_a_kept_binding,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"draft\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/rowless_binding_gains_a_row,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/second_outcome,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"B\\\" binding \\\"spec\\\"\"\n\nRepair coverage probes (run at 4e1fb87f content before commit; --run TestCommitmentPlanRefusesObligationFreeBinding):\nC1 --swap 'if outcome.ID == id {' --with 'if outcome.ID == id || len(outcome.Deliverables) > 0 {'\n  bit,internal/commitment/authority.go,swap,failed,2,yes\n  failed: legacy_binding_moved_to_another_outcome, inactive_milestone\nC2 --omit ' && obligationFree(outcome, kept)'\n  bit,internal/commitment/authority.go,omit,failed,1,yes\n  failed: rowless_binding_gains_a_row\nC3 --swap 'range outcome.Deliverables {' --with 'range outcome.Deliverables[:min(1, len(outcome.Deliverables))] {'\n  bit,internal/commitment/authority.go,swap,failed,1,yes\n  failed: new_binding_after_a_kept_binding\nC4 --swap 'range milestone.Outcomes {' --with 'range milestone.Outcomes[:min(1, len(milestone.Outcomes))] {'\n  bit,internal/commitment/authority.go,swap,failed,2,yes\n  failed: legacy_binding_moved_to_another_outcome, second_outcome\n"
            }
          }
        },
        {
          "id": "t1-conformance-v2",
          "performer": "claude:fd_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t1_r1",
            "digest": "sha256:441c5ebf39d293c93048c5c07a0738621fd1736bb10eb64f46572525317742b9",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/conformance\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,a7414079fc441354e4587d4910ecc710be4935cb,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38391\nfailures[0]{package,test,line}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceProseBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestSkillDescriptionBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem\"\n"
          },
          "requirement": "t1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-bench-v2",
          "performer": "claude:fd_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t1_r1",
            "digest": "sha256:cd5fdaca7e874300fc1cddd31301e0a0cf6075604b39aded102fd5412b702d62",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./cmd/bench\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,a7414079fc441354e4587d4910ecc710be4935cb,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,12831\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1-commitment-v3",
          "performer": "claude:fd_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t1_r1",
            "digest": "sha256:08027838a5c8ed1d9f61ffd52aedab6a7a0927d2867d457f54c2d54fba96be38",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,a7414079fc441354e4587d4910ecc710be4935cb,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,6196\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --swap 'if obligationFree(outcome, binding) && !retainedObligationFree(current, outcome.ID, binding) {' --with 'if false && obligationFree(outcome, binding) && !retainedObligationFree(current, outcome.ID, binding) {' --package ./internal/commitment --run TestCommitmentPlanRefusesObligationFreeBinding\nexit: 0 (bench probe verb)\ntree[1]{target,head,dirty}:\n  FT390,ad660e7e13a4a510785a652905ff21df3e6d5f6c,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/authority.go,swap,failed,7,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding,passed,10\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,3\nmutated focused test run exit: 1 (cause failed, package status fail; the failing go test run exits 1)\nfailures[7]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/changed_legacy_binding,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/inactive_milestone,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"B\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/legacy_binding_moved_to_another_outcome,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"B\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/new_binding,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/new_binding_after_a_kept_binding,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"draft\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/rowless_binding_gains_a_row,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/second_outcome,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"B\\\" binding \\\"spec\\\"\"\n\nRepair coverage probes (run at 4e1fb87f content before commit; --run TestCommitmentPlanRefusesObligationFreeBinding):\nC1 --swap 'if outcome.ID == id {' --with 'if outcome.ID == id || len(outcome.Deliverables) > 0 {'\n  bit,internal/commitment/authority.go,swap,failed,2,yes\n  failed: legacy_binding_moved_to_another_outcome, inactive_milestone\nC2 --omit ' && obligationFree(outcome, kept)'\n  bit,internal/commitment/authority.go,omit,failed,1,yes\n  failed: rowless_binding_gains_a_row\nC3 --swap 'range outcome.Deliverables {' --with 'range outcome.Deliverables[:min(1, len(outcome.Deliverables))] {'\n  bit,internal/commitment/authority.go,swap,failed,1,yes\n  failed: new_binding_after_a_kept_binding\nC4 --swap 'range milestone.Outcomes {' --with 'range milestone.Outcomes[:min(1, len(milestone.Outcomes))] {'\n  bit,internal/commitment/authority.go,swap,failed,2,yes\n  failed: legacy_binding_moved_to_another_outcome, second_outcome\n"
          },
          "requirement": "t1-commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0,
          "probe": {
            "mutation": "In BuildPlan, skip the refusal of a new obligation-free binding. TestCommitmentPlanRefusesObligationFreeBinding must fail and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:fd_t1_r1",
              "digest": "sha256:08027838a5c8ed1d9f61ffd52aedab6a7a0927d2867d457f54c2d54fba96be38",
              "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,a7414079fc441354e4587d4910ecc710be4935cb,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,6196\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --swap 'if obligationFree(outcome, binding) && !retainedObligationFree(current, outcome.ID, binding) {' --with 'if false && obligationFree(outcome, binding) && !retainedObligationFree(current, outcome.ID, binding) {' --package ./internal/commitment --run TestCommitmentPlanRefusesObligationFreeBinding\nexit: 0 (bench probe verb)\ntree[1]{target,head,dirty}:\n  FT390,ad660e7e13a4a510785a652905ff21df3e6d5f6c,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/authority.go,swap,failed,7,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding,passed,10\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,3\nmutated focused test run exit: 1 (cause failed, package status fail; the failing go test run exits 1)\nfailures[7]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/changed_legacy_binding,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/inactive_milestone,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"B\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/legacy_binding_moved_to_another_outcome,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"B\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/new_binding,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/new_binding_after_a_kept_binding,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"draft\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/rowless_binding_gains_a_row,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"A\\\" binding \\\"spec\\\"\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPlanRefusesObligationFreeBinding/second_outcome,\"authority_test.go:123: BuildPlan() = <nil>, want the refusal of outcome \\\"B\\\" binding \\\"spec\\\"\"\n\nRepair coverage probes (run at 4e1fb87f content before commit; --run TestCommitmentPlanRefusesObligationFreeBinding):\nC1 --swap 'if outcome.ID == id {' --with 'if outcome.ID == id || len(outcome.Deliverables) > 0 {'\n  bit,internal/commitment/authority.go,swap,failed,2,yes\n  failed: legacy_binding_moved_to_another_outcome, inactive_milestone\nC2 --omit ' && obligationFree(outcome, kept)'\n  bit,internal/commitment/authority.go,omit,failed,1,yes\n  failed: rowless_binding_gains_a_row\nC3 --swap 'range outcome.Deliverables {' --with 'range outcome.Deliverables[:min(1, len(outcome.Deliverables))] {'\n  bit,internal/commitment/authority.go,swap,failed,1,yes\n  failed: new_binding_after_a_kept_binding\nC4 --swap 'range milestone.Outcomes {' --with 'range milestone.Outcomes[:min(1, len(milestone.Outcomes))] {'\n  bit,internal/commitment/authority.go,swap,failed,2,yes\n  failed: legacy_binding_moved_to_another_outcome, second_outcome\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "fdc1-standards-r1",
          "performer": "claude:fdc1_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fdc1_standards",
            "digest": "sha256:8ffda479b549219eadc3ac41664c8b63d30a2db171b338f647f2ec9af82df82b",
            "excerpt": "Standards axis, FD-C1 (9febee8f..29eabd30), sonnet high, evidence sha256:e06ff36b check-current=true.\nS1 commitment_landing_fixture_test.go:62 obligationFreeRoute reuses deliveryRoutes[1].seed by position; no-op; confidence 3.\nS2 authority.go:129 retainedObligationFree compares the two DeliveryBinding fields; refuted (model.go:45-48 has two fields); no-op; confidence 2.\nS3 authority_test.go:79, delivery_test.go:15 two fixture builders with different identity conventions; refuted (different shapes, shared helper exists); no-op; confidence 2.\nTwo refusal texts: not a defect; one predicate obligationFree (delivery.go:41) serves Deliver and BuildPlan; texts differ by spec lines 87 and 89; plan remedy clause harmless; no-op; confidence 7.\nCount 3, all no-op. Worst: S1.\nAdvice (no finding ID): S1 commitment_landing_fixture_test.go:62 positional seed reuse deliveryRoutes[1].seed.\n"
          },
          "axis": "Standards",
          "base": "9febee8f284bfed4f6de714717f4364adede738e",
          "tip": "29eabd3090ad23b7e6dc359a1922029df50cab79",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "fdc1-spec-r1",
          "performer": "claude:fdc1_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fdc1_spec",
            "digest": "sha256:46ade649f197ca38e7617a5f3700e6a464c1fb9da74d529449dcab93b5fdbe73",
            "excerpt": "Spec axis, FD-C1 (9febee8f..29eabd30), sonnet high, evidence sha256:e06ff36b check-current=true.\nP1 spec.md:227, ticket 2 Writes line: the census pin fence expansion (plan commit bf4654cc) and the pre-expansion charge identity are not a behavioral deviation; no-op; confidence 7.\nClosed decisions held: one predicate obligationFree (delivery.go:41) called by Deliver and BuildPlan; Validate untouched; BuildPlan refuses new or changed bindings in every milestone; Deliver checks the active milestone; no repository or landing caller edit.\nRows FD1-FD11, FD23, FD24, FD26: met.\nCount 1, no-op. Worst: P1 (record hygiene).\n"
          },
          "axis": "Spec",
          "base": "9febee8f284bfed4f6de714717f4364adede738e",
          "tip": "29eabd3090ad23b7e6dc359a1922029df50cab79",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "fdc1-coverage-r1",
          "performer": "claude:fdc1_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "876937bb5fdcad91e0f23a047f59964ca3c3fb61",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:fdc1_coverage",
            "digest": "sha256:293da5b6d0a12d9fbc8a80fbd3b87daa0ab70f6aafe8a58dd949c48795c6e1ea",
            "excerpt": "Coverage axis, FD-C1 (9febee8f..29eabd30), sonnet high, evidence sha256:e06ff36b check-current=true, tree clean after probes.\nC1 authority.go:127 retention across outcomes untested; probe t3 swap (match any outcome with deliverables) silent, restored yes; auto-fix; confidence 8.\nC2 authority.go:129 \"already obligation-free in current\" clause untested (rowless outcome gains sources, binding kept); probe t4 omit silent, restored yes; auto-fix; confidence 9.\nC3 authority.go:109 only the first deliverable of an outcome is exercised; probe t5 swap Deliverables[:1] silent, restored yes; auto-fix; confidence 8.\nC4 authority.go:108 only the first outcome of a milestone is exercised; probe t6 swap Outcomes[:1] silent, restored yes; auto-fix; confidence 8.\nRefuted: obligation list change on a kept binding; outcome moved between milestones; landing order (commitment_light_landing_test.go:347-361); Store.Plan receipt; inactive-only Deliver.\nCount 4. Worst: C2.\n"
          },
          "axis": "Coverage",
          "base": "9febee8f284bfed4f6de714717f4364adede738e",
          "tip": "29eabd3090ad23b7e6dc359a1922029df50cab79",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": []
        }
      ]
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
      "from": "sha256:e1af045505879acfbe050e31acc668e35464b452e9a3f71f2b38add2733381ef",
      "to": "sha256:37236d7f76612ed55db6b4a27b59899919a497c041af08415e14869512b9d73d",
      "chunk_ids": {
        "FD-C1": [
          "FD-C1"
        ]
      }
    }
  ]
}
```
