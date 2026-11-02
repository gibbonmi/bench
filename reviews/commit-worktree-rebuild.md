# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/commit-worktree-rebuild/spec.md",
  "plan_digest": "sha256:caecfd37289f7819cd30ab7da133e5019f36c7a65df862ec6a16d535b9272523",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CW-C1",
      "base": "8fcad8341e760744d06ce632ea0a978684025f81",
      "tip": "2a5c53e3558ed3e39a1720d242a17008a4f1f3ae",
      "plan_digest": "sha256:caecfd37289f7819cd30ab7da133e5019f36c7a65df862ec6a16d535b9272523",
      "source_digest": "5d93cb2bfbb1812bf823fb3c46c8381a8efa8fa3",
      "acceptance_rows": [
        "CW01",
        "CW02",
        "CW03",
        "CW04",
        "CW05",
        "CW06",
        "CW07",
        "CW08",
        "CW09",
        "CW10",
        "CW11",
        "CW12",
        "CW41",
        "CW42",
        "CW44",
        "CW46",
        "CW47",
        "CW49"
      ],
      "verification": [
        {
          "id": "cw-c1-owners",
          "performer": "ft406-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "5d93cb2bfbb1812bf823fb3c46c8381a8efa8fa3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01",
            "digest": "sha256:e34faa480264a5c2e7b616d22063fb9b54dc78ef4702f8ea005726cc7858b891",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,DiffChangesBuildInputs|Refresh,passed,14\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,435\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,1137,14\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,\"decision_test.go:132: Refresh(nil build that publishes nothing) = (true, \u003cnil>, \u003cnil>), want no rebuild, build error \u003cnil>, and a seal refusal\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2a5c53e3558ed3e39a1720d242a17008a4f1f3ae",
          "requirement": "owners",
          "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs|Refresh'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make Refresh return rebuilt true after a nil build error without the second Verify read, and record that TestRefreshSecondReadDecides turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01",
              "digest": "sha256:e34faa480264a5c2e7b616d22063fb9b54dc78ef4702f8ea005726cc7858b891",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,DiffChangesBuildInputs|Refresh,passed,14\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,435\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,1137,14\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,\"decision_test.go:132: Refresh(nil build that publishes nothing) = (true, \u003cnil>, \u003cnil>), want no rebuild, build error \u003cnil>, and a seal refusal\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "cw-c1-deleted-input",
          "performer": "ft406-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "5d93cb2bfbb1812bf823fb3c46c8381a8efa8fa3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01",
            "digest": "sha256:e3a3a2164bd749478e96e0cfd632f48fc15b0005926ee121653a6cd33fa84d12",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestDiffChangesBuildInputs/deleted_input,7,1,0,null,yes\n  2,bit,omit,failed,1,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,7,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "2a5c53e3558ed3e39a1720d242a17008a4f1f3ae",
          "requirement": "deleted-input",
          "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the deleted-file test from DiffChangesBuildInputs so that it reads only the post-commit input set, and record that TestDiffChangesBuildInputs turns red on the deleted-input case.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01",
              "digest": "sha256:e3a3a2164bd749478e96e0cfd632f48fc15b0005926ee121653a6cd33fa84d12",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestDiffChangesBuildInputs/deleted_input,7,1,0,null,yes\n  2,bit,omit,failed,1,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,7,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "cw-c1-renamed-input",
          "performer": "ft406-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "5d93cb2bfbb1812bf823fb3c46c8381a8efa8fa3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01",
            "digest": "sha256:e3a3a2164bd749478e96e0cfd632f48fc15b0005926ee121653a6cd33fa84d12",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestDiffChangesBuildInputs/deleted_input,7,1,0,null,yes\n  2,bit,omit,failed,1,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,7,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "2a5c53e3558ed3e39a1720d242a17008a4f1f3ae",
          "requirement": "renamed-input",
          "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove --no-renames from the diff that DiffChangesBuildInputs reads, and record that TestDiffChangesBuildInputs turns red on the renamed-input case.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01",
              "digest": "sha256:e3a3a2164bd749478e96e0cfd632f48fc15b0005926ee121653a6cd33fa84d12",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestDiffChangesBuildInputs/deleted_input,7,1,0,null,yes\n  2,bit,omit,failed,1,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,7,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "cw-c1-landing",
          "performer": "ft406-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "5d93cb2bfbb1812bf823fb3c46c8381a8efa8fa3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01",
            "digest": "sha256:40d9b65770f977577714ef7a5c0bc4bd5e5016dcab52fe910b7edc4a0634f09e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_buildinputs.go,swap,failed,2,TestLandCommandReportsInstallStepForABrokerChangingDiff,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,LandCommandReportsInstallStepForABrokerChangingDiff|LandReportsInstallStepForADeletedBrokerInput|LandRefreshesTheBrokerAfterPublication|LandReportsAFailedRefresh|LandSkipsAFreshBroker|RefreshHoldsAnInstalledCandidate,passed,13\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,909\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,2777,13\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:91: landing printed no broker install notice: \\\"landing source{review_base=d7352d7e0b7d1a7d44a7b9458553c915b7a52f04,assignment_start=f39a468c7bc38e7c3ebf6a7859af4c324417be4f}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (810 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsInstallStepForADeletedBrokerInput,\"land_broker_notice_test.go:71: landing printed no broker install notice: \\\"landing source{review_base=caf842eccab40f007d8ffdb9c44678e133301b66,assignment_start=5c5f64848d2f3e42d3a84d12ed11cca2769cf5d5}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (775 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2a5c53e3558ed3e39a1720d242a17008a4f1f3ae",
          "requirement": "landing",
          "command": "bench test --package ./internal/worktree --run 'LandCommandReportsInstallStepForABrokerChangingDiff|LandReportsInstallStepForADeletedBrokerInput|LandRefreshesTheBrokerAfterPublication|LandReportsAFailedRefresh|LandSkipsAFreshBroker|RefreshHoldsAnInstalledCandidate'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make DiffChangesBuildInputs return false for every diff, and record that TestLandCommandReportsInstallStepForABrokerChangingDiff turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01",
              "digest": "sha256:40d9b65770f977577714ef7a5c0bc4bd5e5016dcab52fe910b7edc4a0634f09e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_buildinputs.go,swap,failed,2,TestLandCommandReportsInstallStepForABrokerChangingDiff,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,LandCommandReportsInstallStepForABrokerChangingDiff|LandReportsInstallStepForADeletedBrokerInput|LandRefreshesTheBrokerAfterPublication|LandReportsAFailedRefresh|LandSkipsAFreshBroker|RefreshHoldsAnInstalledCandidate,passed,13\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,909\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,2777,13\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:91: landing printed no broker install notice: \\\"landing source{review_base=d7352d7e0b7d1a7d44a7b9458553c915b7a52f04,assignment_start=f39a468c7bc38e7c3ebf6a7859af4c324417be4f}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (810 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsInstallStepForADeletedBrokerInput,\"land_broker_notice_test.go:71: landing printed no broker install notice: \\\"landing source{review_base=caf842eccab40f007d8ffdb9c44678e133301b66,assignment_start=5c5f64848d2f3e42d3a84d12ed11cca2769cf5d5}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (775 bytes)\",1\nskips[0]{package,test,reason}:\n"
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

## CW-C1 ticket 01 author witness reds

The ticket 01 author `ft406-author-01` (opus / medium) committed ticket 01 at `2a5c53e35` and froze CW-C1 at `56d53b85c`.
The CW44 test went red against the unchanged production code, and each other witness went red under a `bench probe` mutation after the ticket commit.
The plan probes `owners`, `deleted-input`, `renamed-input`, and `landing` answered `bit` with `restored` `yes`, and the record holds them.

| subject | mutation | failing witness | row |
| --- | --- | --- | --- |
| `internal/worktree/land.go` | unchanged production code | `TestLandReportsInstallStepForADeletedBrokerInput` | CW44 |
| `internal/freshness/freshness_buildinputs.go` | intersection removed | `TestDiffChangesBuildInputs/build_script` | CW01 |
| `internal/freshness/freshness_buildinputs.go` | every changed path counts | `TestDiffChangesBuildInputs/markdown_only` | CW02 |
| `internal/freshness/freshness_buildinputs.go` | manifest guard answers true | `TestDiffChangesBuildInputs/no_build-input_manifest` | CW03 |
| `internal/freshness/freshness_buildinputs.go` | test-file exclusion removed | `TestDiffChangesBuildInputs/deleted_test_file` | CW42 |
| `internal/freshness/freshness_verify.go` | first `Verify` read removed | `TestRefreshBuildsOnlyAStaleExecutable/fresh` | CW04 |
| `internal/freshness/freshness_verify.go` | build call removed | `TestRefreshBuildsOnlyAStaleExecutable/stale` and `/absent` | CW05, CW08 |
| `internal/freshness/freshness_verify.go` | `buildErr` dropped | `TestRefreshSecondReadDecides/failing_build` | CW07 |
| `internal/worktree/land_effects.go` | `err` branch removed | `TestLandReportsAFailedRefresh` | CW46 |

The author reports one gap: the mutation `if false && buildErr != nil` in `refreshBroker` stays green, because the build of `TestLandReportsAFailedRefresh` returns nil. No test grades the text `landing refresh build failed: ` that CW46 names.
The coordinator probe swapped the base of the diff test in the stack landing for the tip, and 44 stack tests stayed green. No coverage row of this spec covers the stack landing, so that gap goes to the learnings.
