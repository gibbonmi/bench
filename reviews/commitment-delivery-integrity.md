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

## FD-C1 repair state

Repair cycle 1 of 2 is consumed, and it closes C1 to C4. A fresh ticket 1
repair session added test rows in commit `4e1fb87f`. The confirming round of
all three axes at that tip found no finding. Probes confirmed C2, C3, and C4,
and the Spec axis confirmed C1 by reading.

Advice: the `slices.Equal` clause in `retainedObligationFree` is redundant,
because both bindings on that path list no obligation. The `changed` fixture
in `authority_test.go` repeats the `legacy` builder expression.

## FD-C2 pickup

The FD-C2 review covers the frozen pair `4e1fb87f..90697fd0`. The Coverage
findings start repair cycle 1 of 2. This cycle is the one hardening cycle of
the chunk, because rows FD12 to FD22 and FD25 prove.

### FD-C2 Standards

Findings: 0. Worst issue: none.

Advice: the roadmap-row source rule appears inline in
`TestCommitmentPlanSourcesAreCanonical`, beside `bindObligationFree`, and in
`authority_internal_test.go`.

### FD-C2 Spec

Findings: 0. Worst issue: none. Every chunk row is met, and the FD-C1 rows
still hold.

### FD-C2 Coverage

Findings: 2. Worst issue: C1.

- C1 (`auto-fix`, confidence 5): `internal/commitment/authority.go:111`. No
  `boundSources` test has two sources with an equal ID and path and a
  different identity. A probe that drops the identity tiebreak was silent.
- C3 (`auto-fix`, confidence 5): `internal/commitment/authority.go:110`. No
  test tells byte order from a case-folded compare. A probe that folds the ID
  case was silent.

The reviewer decided on 2026-10-05 to leave the whole-binding comparison at
`internal/commitment/authority.go:74` unpinned. A test for it would freeze the
re-pin behavior that the spec puts out of scope, so the idea is parked as
intake.

```bench-review-record
{
  "version": 2,
  "spec": "specs/commitment-delivery-integrity/spec.md",
  "plan_digest": "sha256:de78ccd7ca2d1e0e46fb3800c0176a9ac1f589f2786aee72b6432691f2ca8dae",
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
        },
        {
          "id": "t2-commitment-v2",
          "performer": "claude:fd_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t2",
            "digest": "sha256:a99fdc5fdce181ff4fbd6d6dc3753e6537a4fdc55d0d2bfa3123442139e8b3f5",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,6313\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n$ bench worktree exec FT390 -- bench probe internal/commitment/delivery.go --swap 'return policy, nil, fmt.Errorf(\"deliverable %q of outcome %q names no obligation; list the sources that it satisfies with bench commitment plan --input <file>\", path, outcome.ID)' --with 'return policy, nil, nil' --package ./internal/commitment --run 'TestCommitmentDeliverRowless$'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/delivery.go,swap,failed,1,yes\nmutated focused test run: exit 1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,3\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentDeliverRowless,\"delivery_test.go:96: obligation-free Deliver = <nil>, want the names-no-obligation refusal for an outcome with sources\"\n"
          },
          "requirement": "t2-commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0,
          "probe": {
            "mutation": "In Deliver, return the policy unchanged with no error for an obligation-free binding of an outcome with sources. TestCommitmentDeliverRowless must fail and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:fd_t2",
              "digest": "sha256:a99fdc5fdce181ff4fbd6d6dc3753e6537a4fdc55d0d2bfa3123442139e8b3f5",
              "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,6313\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n$ bench worktree exec FT390 -- bench probe internal/commitment/delivery.go --swap 'return policy, nil, fmt.Errorf(\"deliverable %q of outcome %q names no obligation; list the sources that it satisfies with bench commitment plan --input <file>\", path, outcome.ID)' --with 'return policy, nil, nil' --package ./internal/commitment --run 'TestCommitmentDeliverRowless$'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/delivery.go,swap,failed,1,yes\nmutated focused test run: exit 1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,3\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentDeliverRowless,\"delivery_test.go:96: obligation-free Deliver = <nil>, want the names-no-obligation refusal for an outcome with sources\"\n"
            }
          }
        },
        {
          "id": "t2-commitment-repository-v2",
          "performer": "claude:fd_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t2",
            "digest": "sha256:f18dbdce893c9d21a3844544b2235b3ca27bbdc671e42e832f08814e643845ca",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment/repository\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,1751\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-commitment-repository",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t2-worktree-v2",
          "performer": "claude:fd_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t2",
            "digest": "sha256:e7335b433646a5d9bcb594b2b1d9b06a2638ecfae4aab4344b234feac9e08384",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/worktree\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,67799\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"unix sockets unavailable\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable\"\n"
          },
          "requirement": "t2-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t2-conformance-v2",
          "performer": "claude:fd_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t2",
            "digest": "sha256:eab60a53d3643a7d50fba8f7bf62e379018d7bebc3983cd40fcc80a4e61efa22",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/conformance\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,43431\nfailures[0]{package,test,line}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceProseBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestSkillDescriptionBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem\"\n"
          },
          "requirement": "t2-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t2-bench-v2",
          "performer": "claude:fd_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t2",
            "digest": "sha256:1ab3b053bf4f472244ab15ced58f2f94422d1dad6dba3839d0d77626d8dc0ad4",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./cmd/bench\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,13639\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
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
        },
        {
          "id": "fdc1-standards-r2",
          "performer": "claude:fdc1_standards_r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fdc1_standards_r2",
            "digest": "sha256:b3c86ad3951401ac94671140454680cbb52188789039ea5ef81d88a27b22e11e",
            "excerpt": "Standards axis, FD-C1 confirming round (repair delta 33465a09..4e1fb87f), sonnet high, evidence sha256:f0bab143 check-current=true.\nFindings: 0.\nAdvice (no finding ID): authority_test.go:97-98 the changed fixture repeats the legacy builder expression before it changes one identity; pre-existing; cosmetic.\nRefuted: bindObligationFree is the one builder and obligationFreeMilestone wraps it; comments state current behavior; the message literal is a pre-existing independent expectation.\n"
          },
          "axis": "Standards",
          "base": "9febee8f284bfed4f6de714717f4364adede738e",
          "tip": "4e1fb87f9d5378c895e4219e74e339bb87cb1d46",
          "finding_ids": [],
          "supersedes": [
            "fdc1-standards-r1"
          ]
        },
        {
          "id": "fdc1-spec-r2",
          "performer": "claude:fdc1_spec_r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fdc1_spec_r2",
            "digest": "sha256:51d871bcfc6c0202cb447601cd28851699dc1fd60f88513b28befe448f29daca",
            "excerpt": "Spec axis, FD-C1 confirming round (repair delta 33465a09..4e1fb87f), sonnet high, evidence sha256:f0bab143 check-current=true.\nFindings: 0.\nFolds C1-C4 confirmed by code reading (authority_test.go:110-112 against authority.go:108-129).\nFD1, FD3, FD4, FD23, FD24, FD26 still met; FD9 (parse_test.go:133) accepts the moved fixture path; FD2 and FD5 do not use the helper.\n"
          },
          "axis": "Spec",
          "base": "9febee8f284bfed4f6de714717f4364adede738e",
          "tip": "4e1fb87f9d5378c895e4219e74e339bb87cb1d46",
          "finding_ids": [],
          "supersedes": [
            "fdc1-spec-r1"
          ]
        },
        {
          "id": "fdc1-coverage-r2",
          "performer": "claude:fdc1_coverage_r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "579e52ff4d3dadd051e6a9dce337c58646b90739",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fdc1_coverage_r2",
            "digest": "sha256:ff27c423ab15151fc40af372fc58e060f7964d4627b29b4daaa425f09cd52f84",
            "excerpt": "Coverage axis, FD-C1 confirming round (repair delta 33465a09..4e1fb87f), sonnet high, evidence sha256:f0bab143 check-current=true.\np1 authority.go:129 omit obligationFree(outcome, kept): bit, restored yes (C2 confirmed).\np2 authority.go:109 swap Deliverables[:1]: bit, restored yes (C3 confirmed).\np3 authority.go:129 omit slices.Equal: silent, restored yes; equivalent mutant because both bindings are obligation-free on that path; refuted, no finding.\nCoordinator probe authority.go:108 swap Outcomes[:1]: bit, 2 failed (legacy_binding_moved_to_another_outcome, second_outcome), restored yes (C4 confirmed).\nC1: repair session probe t3 bit; Spec axis confirmed by reading.\nFindings: 0.\nAdvice (no finding ID): the slices.Equal clause at authority.go:129 is redundant production code.\n"
          },
          "axis": "Coverage",
          "base": "9febee8f284bfed4f6de714717f4364adede738e",
          "tip": "4e1fb87f9d5378c895e4219e74e339bb87cb1d46",
          "finding_ids": [],
          "supersedes": [
            "fdc1-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "FD-C2",
      "base": "4e1fb87f9d5378c895e4219e74e339bb87cb1d46",
      "tip": "f85f69225daee9c88cbfd45e18584750ccc7fcb4",
      "plan_digest": "sha256:de78ccd7ca2d1e0e46fb3800c0176a9ac1f589f2786aee72b6432691f2ca8dae",
      "source_digest": "4a344fe82f87536a8f9d0fdefd8fe132f1e51b93",
      "acceptance_rows": [
        "FD12",
        "FD13",
        "FD14",
        "FD15",
        "FD16",
        "FD17",
        "FD18",
        "FD19",
        "FD25",
        "FD20",
        "FD21",
        "FD22"
      ],
      "verification": [
        {
          "id": "t3-commitment-v1",
          "performer": "claude:fd_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "93c18dfe15b94618c9651a12fd9761cebd4f4d6a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t3",
            "digest": "sha256:6545425bb51cc4261c65ce9d5be60291b5bb589063f26e02d70f71db7c14a3eb",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit 0\ntree: FT390,f7aae5ad0910cb6f946d3a1ae6276e1adb23e871,false\npackages: github.com/gibbonmi/bench/internal/commitment,pass,5847\nfailures[0] skips[0]\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --swap \"open, _ := Unsettled(*current)\" --with \"open := policySources(*current)\" --package ./internal/commitment --run TestCommitmentPlanSourcesOmitDroppedDeliverable\nprobe{verdict,subject,mutation,cause,failed_tests,restored}: bit,internal/commitment/authority.go,swap,failed,1,yes\npackages: github.com/gibbonmi/bench/internal/commitment,fail,3\nfailures: TestCommitmentPlanSourcesOmitDroppedDeliverable,\"authority_test.go:140: plan sources = [{ID:FT1 ...} {ID:spec ...}] ...\"\n"
          },
          "requirement": "t3-commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0,
          "probe": {
            "mutation": "In BuildPlan, add each unsettled deliverable of the current policy to the plan sources again. TestCommitmentPlanSourcesOmitDroppedDeliverable must fail and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:fd_t3",
              "digest": "sha256:6545425bb51cc4261c65ce9d5be60291b5bb589063f26e02d70f71db7c14a3eb",
              "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit 0\ntree: FT390,f7aae5ad0910cb6f946d3a1ae6276e1adb23e871,false\npackages: github.com/gibbonmi/bench/internal/commitment,pass,5847\nfailures[0] skips[0]\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --swap \"open, _ := Unsettled(*current)\" --with \"open := policySources(*current)\" --package ./internal/commitment --run TestCommitmentPlanSourcesOmitDroppedDeliverable\nprobe{verdict,subject,mutation,cause,failed_tests,restored}: bit,internal/commitment/authority.go,swap,failed,1,yes\npackages: github.com/gibbonmi/bench/internal/commitment,fail,3\nfailures: TestCommitmentPlanSourcesOmitDroppedDeliverable,\"authority_test.go:140: plan sources = [{ID:FT1 ...} {ID:spec ...}] ...\"\n"
            }
          }
        },
        {
          "id": "t3-commitment-repository-v1",
          "performer": "claude:fd_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "93c18dfe15b94618c9651a12fd9761cebd4f4d6a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t3",
            "digest": "sha256:223a70e42a62f05c6d1db1ab5b9ff74787abc6cb4d2aa46d84efdccb0621b850",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment/repository\nexit 0\ntree: FT390,f7aae5ad0910cb6f946d3a1ae6276e1adb23e871,false\npackages: github.com/gibbonmi/bench/internal/commitment/repository,pass,2233\nfailures[0] skips[0]\n"
          },
          "requirement": "t3-commitment-repository",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t3-conformance-v1",
          "performer": "claude:fd_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "93c18dfe15b94618c9651a12fd9761cebd4f4d6a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t3",
            "digest": "sha256:63dbe36ac81da408a56536eb394e261c068cb8eabb9aaa758e3ad57c4660a837",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/conformance\nexit 0\ntree: FT390,f7aae5ad0910cb6f946d3a1ae6276e1adb23e871,false\npackages: github.com/gibbonmi/bench/internal/conformance,pass,39910\nfailures[0]\nskips[3]: TestGuidanceProseBudgetRefusesNonRegularSubjects/socket (capability: fifo), TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device (capability: privilege), TestSkillDescriptionBudgetRefusesNonRegularSubjects/socket (capability: fifo)\n"
          },
          "requirement": "t3-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t4-commitment-v1",
          "performer": "claude:fd_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "93c18dfe15b94618c9651a12fd9761cebd4f4d6a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t4",
            "digest": "sha256:899442f177fc89585baa3c39e5a020bb4a38188e7ec35a56e7d92f6aaf1fad32",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit 0\ntree[1]{target,head,dirty}:\n  FT390,9d224f92399108cbdf8bd1a06906f64cd02fda85,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,6360\nfailures[0]{package,test,line}:\n\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --swap \"slices.SortFunc(sorted, func(a, b SourceBinding) int {\" --with \"_ = (func(a, b SourceBinding) int {\" --package ./internal/commitment --run TestBoundSourcesIgnoresInputOrder\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/authority.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,4\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestBoundSourcesIgnoresInputOrder,\"authority_internal_test.go:15: boundSources lists = [{FT9 ...} {FT1 ...}] ...\"\nnote: bench probe --omit of the sort statement is refused (build fails: \"cmp\" imported and not used), so the probe runs as a swap that leaves the copy unsorted.\n"
          },
          "requirement": "t4-commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0,
          "probe": {
            "mutation": "In boundSources, return the copy without the sort. TestBoundSourcesIgnoresInputOrder must fail and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:fd_t4",
              "digest": "sha256:899442f177fc89585baa3c39e5a020bb4a38188e7ec35a56e7d92f6aaf1fad32",
              "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit 0\ntree[1]{target,head,dirty}:\n  FT390,9d224f92399108cbdf8bd1a06906f64cd02fda85,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,6360\nfailures[0]{package,test,line}:\n\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --swap \"slices.SortFunc(sorted, func(a, b SourceBinding) int {\" --with \"_ = (func(a, b SourceBinding) int {\" --package ./internal/commitment --run TestBoundSourcesIgnoresInputOrder\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/authority.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,4\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestBoundSourcesIgnoresInputOrder,\"authority_internal_test.go:15: boundSources lists = [{FT9 ...} {FT1 ...}] ...\"\nnote: bench probe --omit of the sort statement is refused (build fails: \"cmp\" imported and not used), so the probe runs as a swap that leaves the copy unsorted.\n"
            }
          }
        },
        {
          "id": "t4-commitment-repository-v1",
          "performer": "claude:fd_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "93c18dfe15b94618c9651a12fd9761cebd4f4d6a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t4",
            "digest": "sha256:ff3de982962c716bce3687c14b4603fe27eb7fff26ce966b56875b4200ab3b65",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment/repository\nexit 0\ntree[1]{target,head,dirty}:\n  FT390,9d224f92399108cbdf8bd1a06906f64cd02fda85,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,2322\nfailures[0]{package,test,line}:\n"
          },
          "requirement": "t4-commitment-repository",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t4-conformance-v1",
          "performer": "claude:fd_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "93c18dfe15b94618c9651a12fd9761cebd4f4d6a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t4",
            "digest": "sha256:6c177728cc18e6111e6eeba99201af182b5c61ae15cd73287ccbaa6aa3172c19",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/conformance\nexit 0\ntree[1]{target,head,dirty}:\n  FT390,9d224f92399108cbdf8bd1a06906f64cd02fda85,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,41717\nfailures[0]{package,test,line}:\nskips[3]: two unix-socket subjects and one character-device subject (environment capability skips)\n"
          },
          "requirement": "t4-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t4-commitment-v2",
          "performer": "claude:fd_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4a344fe82f87536a8f9d0fdefd8fe132f1e51b93",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t4_r1",
            "digest": "sha256:52136ca6e3fd450ae1ddc95fd5cb9034339bd1b0e788d6e3c3c880dd349ce7a6",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,bb73045c213a6dfe97314b245a6b52954c465a03,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,5598\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --swap 'slices.SortFunc(sorted, func(a, b SourceBinding) int {' --with '_ = (func(a, b SourceBinding) int {' --package ./internal/commitment --run TestBoundSourcesIgnoresInputOrder\nexit: 0 (bench probe verb)\ntree[1]{target,head,dirty}:\n  FT390,bb73045c213a6dfe97314b245a6b52954c465a03,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/authority.go,swap,failed,1,yes\nselection: package ./internal/commitment, run TestBoundSourcesIgnoresInputOrder, baseline passed\nmutated focused test run exit: 1 (cause failed, package status fail)\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestBoundSourcesIgnoresInputOrder,\"authority_internal_test.go:20: boundSources lists = [{a roadmap/a.md ...} {FT9 roadmap/FT9.md ...} {FT9 ... (1020 bytes)\"\n\nRepair coverage probes (run at the f85f6922 test content before commit; --package ./internal/commitment --run TestBoundSourcesIgnoresInputOrder; baseline passed):\nC1 --swap 'cmp.Compare(a.Identity, b.Identity)' --with 'cmp.Compare(0, 0)'\n  bit,internal/commitment/authority.go,swap,failed,1,yes\n  failed: TestBoundSourcesIgnoresInputOrder (authority_internal_test.go:20, list starts FT1, FT9 sha256:a18b..., so input order leaks through the tie)\nC3 --swap 'cmp.Compare(a.ID, b.ID)' --with 'cmp.Compare(strings.ToLower(a.ID), strings.ToLower(b.ID))'\n  bit,internal/commitment/authority.go,swap,failed,1,yes\n  failed: TestBoundSourcesIgnoresInputOrder (authority_internal_test.go:20, list starts with a before FT1)\n"
          },
          "requirement": "t4-commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0,
          "probe": {
            "mutation": "In boundSources, return the copy without the sort. TestBoundSourcesIgnoresInputOrder must fail and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:fd_t4_r1",
              "digest": "sha256:52136ca6e3fd450ae1ddc95fd5cb9034339bd1b0e788d6e3c3c880dd349ce7a6",
              "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,bb73045c213a6dfe97314b245a6b52954c465a03,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,5598\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec FT390 -- bench probe internal/commitment/authority.go --swap 'slices.SortFunc(sorted, func(a, b SourceBinding) int {' --with '_ = (func(a, b SourceBinding) int {' --package ./internal/commitment --run TestBoundSourcesIgnoresInputOrder\nexit: 0 (bench probe verb)\ntree[1]{target,head,dirty}:\n  FT390,bb73045c213a6dfe97314b245a6b52954c465a03,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/authority.go,swap,failed,1,yes\nselection: package ./internal/commitment, run TestBoundSourcesIgnoresInputOrder, baseline passed\nmutated focused test run exit: 1 (cause failed, package status fail)\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestBoundSourcesIgnoresInputOrder,\"authority_internal_test.go:20: boundSources lists = [{a roadmap/a.md ...} {FT9 roadmap/FT9.md ...} {FT9 ... (1020 bytes)\"\n\nRepair coverage probes (run at the f85f6922 test content before commit; --package ./internal/commitment --run TestBoundSourcesIgnoresInputOrder; baseline passed):\nC1 --swap 'cmp.Compare(a.Identity, b.Identity)' --with 'cmp.Compare(0, 0)'\n  bit,internal/commitment/authority.go,swap,failed,1,yes\n  failed: TestBoundSourcesIgnoresInputOrder (authority_internal_test.go:20, list starts FT1, FT9 sha256:a18b..., so input order leaks through the tie)\nC3 --swap 'cmp.Compare(a.ID, b.ID)' --with 'cmp.Compare(strings.ToLower(a.ID), strings.ToLower(b.ID))'\n  bit,internal/commitment/authority.go,swap,failed,1,yes\n  failed: TestBoundSourcesIgnoresInputOrder (authority_internal_test.go:20, list starts with a before FT1)\n"
            }
          }
        },
        {
          "id": "t4-commitment-repository-v2",
          "performer": "claude:fd_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4a344fe82f87536a8f9d0fdefd8fe132f1e51b93",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t4_r1",
            "digest": "sha256:b418c8dfca26e1b3ed2e08a3e4dbcc42ed80268f63ab8648444ebf3032c44f45",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/commitment/repository\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,bb73045c213a6dfe97314b245a6b52954c465a03,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,2118\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-commitment-repository",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t4-conformance-v2",
          "performer": "claude:fd_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4a344fe82f87536a8f9d0fdefd8fe132f1e51b93",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fd_t4_r1",
            "digest": "sha256:4040eebce9e605e3f8cc7f673ad68796f6759a5d86a5fad0dfa6402b1b2397c6",
            "excerpt": "$ bench worktree exec FT390 -- bench test --package ./internal/conformance\nexit: 0\ntree[1]{target,head,dirty}:\n  FT390,bb73045c213a6dfe97314b245a6b52954c465a03,false\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,39594\nfailures[0]{package,test,line}:\nskips[3]{package,test,reason}:\n  TestGuidanceProseBudgetRefusesNonRegularSubjects/socket (capability: unix sockets unavailable on this filesystem)\n  TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device (capability: cannot create a character device)\n  TestSkillDescriptionBudgetRefusesNonRegularSubjects/socket (capability: unix sockets unavailable on this filesystem)\n"
          },
          "requirement": "t4-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "fdc2-standards-r1",
          "performer": "claude:fdc2_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "93c18dfe15b94618c9651a12fd9761cebd4f4d6a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fdc2_standards",
            "digest": "sha256:71c3aceec3ede408e9b4416838e0dd26c7180b24a3673d4f729240279cdd92b6",
            "excerpt": "Standards axis, FD-C2 (4e1fb87f..90697fd0), sonnet high, evidence sha256:7423d84e check-current=true.\nFindings: 0.\nAdvice (no finding ID): the roadmap-row SourceBinding rule (path roadmap/<row>.md, identity of the row) is spelled inline at authority_test.go:119 beside bindObligationFree (authority_test.go:85) and in authority_internal_test.go:80-81; one line; the internal test cannot share across the import cycle.\nRefuted: second source encoding or order (only boundSources, authority.go:105-118); second settled rule (BuildPlan and policySources call Unsettled); RemoveTickets pairs with WriteTickets; FD12/FD14/FD18 share one table body; comments are timeless.\n"
          },
          "axis": "Standards",
          "base": "4e1fb87f9d5378c895e4219e74e339bb87cb1d46",
          "tip": "90697fd0dfc23439ae1f6204b0f08b692f7d3930",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "fdc2-spec-r1",
          "performer": "claude:fdc2_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "93c18dfe15b94618c9651a12fd9761cebd4f4d6a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:fdc2_spec",
            "digest": "sha256:ee78c7f1fc620bb9731c821c935ffc3fa6f5858f0912d61d110a0180bf1e03a2",
            "excerpt": "Spec axis, FD-C2 (4e1fb87f..90697fd0), sonnet high, evidence sha256:7423d84e check-current=true.\nFindings: 0.\nRows FD12-FD22 and FD25 met (TestCommitmentPlanSurvivesRemovedDeliverable, TestCommitmentPlanSourcesOmitDroppedDeliverable, TestCommitmentRemovalBindsRemovedSource, TestCommitmentPlanAfterDelivery, TestCommitmentApprovalBindsKeptDeliverable, TestBoundSourcesIgnoresInputOrder, TestCommitmentPlanSourcesAreCanonical).\nFD-C1 rows still hold by reading. Closed decisions honored: only Unsettled(*current) sources join; one boundSources owns order and encoding; no repository owner edit.\n"
          },
          "axis": "Spec",
          "base": "4e1fb87f9d5378c895e4219e74e339bb87cb1d46",
          "tip": "90697fd0dfc23439ae1f6204b0f08b692f7d3930",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "fdc2-coverage-r1",
          "performer": "claude:fdc2_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "93c18dfe15b94618c9651a12fd9761cebd4f4d6a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:fdc2_coverage",
            "digest": "sha256:374c7bc40b6f2b3825f0c06729b172e7853ebcfec6ae0e96b9bd5316ce61319a",
            "excerpt": "Coverage axis, FD-C2 (4e1fb87f..90697fd0), sonnet high, evidence sha256:7423d84e check-current=true, all probes restored.\nP1 authority.go:111 drop the identity tiebreak: silent. P2 drop the path tiebreak: bit. P3/P4 authority.go:74 match by ID only: silent in commitment and repository. P5 authority.go:110 case-folded ID compare: silent.\nC1 authority.go:111 no boundSources test with equal ID and path and different identity; auto-fix; confidence 5.\nC3 authority.go:110 byte order is not distinguished from a case-folded compare; auto-fix; confidence 5.\nReviewer decision 2026-10-05 (no finding ID): the whole-binding comparison at authority.go:74 (P3/P4) stays unpinned, because pinning it freezes the out-of-scope re-pin behavior; parked as intake.\nCount 2 (C1, C3). Worst: C1.\n"
          },
          "axis": "Coverage",
          "base": "4e1fb87f9d5378c895e4219e74e339bb87cb1d46",
          "tip": "90697fd0dfc23439ae1f6204b0f08b692f7d3930",
          "finding_ids": [
            "C1",
            "C3"
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
    },
    {
      "from": "sha256:37236d7f76612ed55db6b4a27b59899919a497c041af08415e14869512b9d73d",
      "to": "sha256:aad803245152d604f292375130e8ccf0aae170848859f48fa1b44d4e36300ae8",
      "chunk_ids": {
        "FD-C1": [
          "FD-C1"
        ]
      }
    },
    {
      "from": "sha256:aad803245152d604f292375130e8ccf0aae170848859f48fa1b44d4e36300ae8",
      "to": "sha256:de78ccd7ca2d1e0e46fb3800c0176a9ac1f589f2786aee72b6432691f2ca8dae",
      "chunk_ids": {
        "FD-C1": [
          "FD-C1"
        ],
        "FD-C2": [
          "FD-C2"
        ]
      }
    }
  ]
}
```
