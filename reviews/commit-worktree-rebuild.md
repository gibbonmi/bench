# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/commit-worktree-rebuild/spec.md",
  "plan_digest": "sha256:43bcc719cba9b1bf8d9efe50be6ffc20fb51fe361d8084a238deb7e5996fbb1e",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CW-C1",
      "base": "8fcad8341e760744d06ce632ea0a978684025f81",
      "tip": "856141b4132375c0b870c48ea8c89d7f7314d809",
      "plan_digest": "sha256:43bcc719cba9b1bf8d9efe50be6ffc20fb51fe361d8084a238deb7e5996fbb1e",
      "source_digest": "4de49d6a6a65bc49342369d3735bf52383d138e9",
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
        "CW49",
        "CW51"
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
        },
        {
          "id": "cw-c1-rep1-owners",
          "performer": "ft406-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "82340a1a7b9eace8db09de391de358c27f7ad25d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep1",
            "digest": "sha256:b383f7321fdc4d7d37295e93417fe2f76ee5ce7cc9954d8cefe2a5208b0c138c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,DiffChangesBuildInputs|Refresh,passed,15\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,435\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,1342,15\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,\"decision_test.go:135: Refresh(nil build that publishes nothing) = (true, \u003cnil>, \u003cnil>), want no rebuild, build error \u003cnil>, and a seal refusal\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8fc874450b0ac60b0fa2f7b1555883d0bc74abe6",
          "requirement": "owners",
          "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs|Refresh'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make Refresh return rebuilt true after a nil build error without the second Verify read, and record that TestRefreshSecondReadDecides turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep1",
              "digest": "sha256:b383f7321fdc4d7d37295e93417fe2f76ee5ce7cc9954d8cefe2a5208b0c138c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,DiffChangesBuildInputs|Refresh,passed,15\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,435\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,1342,15\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,\"decision_test.go:135: Refresh(nil build that publishes nothing) = (true, \u003cnil>, \u003cnil>), want no rebuild, build error \u003cnil>, and a seal refusal\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "cw-c1-rep1-deleted-input",
          "performer": "ft406-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "82340a1a7b9eace8db09de391de358c27f7ad25d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep1",
            "digest": "sha256:c74183da5db99eca27ef4ed8241607c3a464cf165820b0f365a5be124ff11cab",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_buildinputs.go,swap,failed,3,TestDiffChangesBuildInputs/deleted_input,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,8\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,737\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,753,8\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestDiffChangesBuildInputs/deleted_input,\"decision_test.go:81: DiffChangesBuildInputs(deleted input) = false, want true\",1\n  github.com/gibbonmi/bench/internal/freshness,TestDiffChangesBuildInputs/deleted_non-ASCII_input,\"decision_test.go:81: DiffChangesBuildInputs(deleted non-ASCII input) = false, want true\",1\n  github.com/gibbonmi/bench/internal/freshness,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,\"decision_test.go:81: DiffChangesBuildInputs(input renamed to a test file) = false, want true\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8fc874450b0ac60b0fa2f7b1555883d0bc74abe6",
          "requirement": "deleted-input",
          "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the deleted-file test from DiffChangesBuildInputs so that it reads only the post-commit input set, and record that TestDiffChangesBuildInputs turns red on the deleted-input case.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep1",
              "digest": "sha256:c74183da5db99eca27ef4ed8241607c3a464cf165820b0f365a5be124ff11cab",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_buildinputs.go,swap,failed,3,TestDiffChangesBuildInputs/deleted_input,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,8\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,737\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,753,8\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestDiffChangesBuildInputs/deleted_input,\"decision_test.go:81: DiffChangesBuildInputs(deleted input) = false, want true\",1\n  github.com/gibbonmi/bench/internal/freshness,TestDiffChangesBuildInputs/deleted_non-ASCII_input,\"decision_test.go:81: DiffChangesBuildInputs(deleted non-ASCII input) = false, want true\",1\n  github.com/gibbonmi/bench/internal/freshness,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,\"decision_test.go:81: DiffChangesBuildInputs(input renamed to a test file) = false, want true\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "cw-c1-rep1-renamed-input",
          "performer": "ft406-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "82340a1a7b9eace8db09de391de358c27f7ad25d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep1",
            "digest": "sha256:17914165d06e58d0ab50af0c65e77635f64c81edc775cf3daef9873a0e7cac9d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_buildinputs.go,swap,failed,1,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,8\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,377\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,822,8\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,\"decision_test.go:81: DiffChangesBuildInputs(input renamed to a test file) = false, want true\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8fc874450b0ac60b0fa2f7b1555883d0bc74abe6",
          "requirement": "renamed-input",
          "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove --no-renames from the diff that DiffChangesBuildInputs reads, and record that TestDiffChangesBuildInputs turns red on the renamed-input case.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep1",
              "digest": "sha256:17914165d06e58d0ab50af0c65e77635f64c81edc775cf3daef9873a0e7cac9d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_buildinputs.go,swap,failed,1,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,8\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,377\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,822,8\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,\"decision_test.go:81: DiffChangesBuildInputs(input renamed to a test file) = false, want true\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "cw-c1-rep1-landing",
          "performer": "ft406-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "82340a1a7b9eace8db09de391de358c27f7ad25d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep1",
            "digest": "sha256:544cc0f299db6828ac6e615c6a517b0c19cf9755d2edcfe3b2e4c70bc3ec858d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_buildinputs.go,swap,failed,2,TestLandCommandReportsInstallStepForABrokerChangingDiff,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,LandCommandReportsInstallStepForABrokerChangingDiff|LandReportsInstallStepForADeletedBrokerInput|LandRefreshesTheBrokerAfterPublication|LandReportsAFailedRefresh|LandSkipsAFreshBroker|RefreshHoldsAnInstalledCandidate,passed,13\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,909\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,2895,13\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:91: landing printed no broker install notice: \\\"landing source{review_base=add5401c323abcefff2591fb4e9342ea42dc0a7a,assignment_start=52e479bf21b0c5b4bb248314d6ae540d03c1b716}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (810 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsInstallStepForADeletedBrokerInput,\"land_broker_notice_test.go:71: landing printed no broker install notice: \\\"landing source{review_base=f0371836265b35e06cd44f55cba20432a23d55af,assignment_start=ccfefc9e773f7a60ca2975a19f5603d434b44535}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (770 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8fc874450b0ac60b0fa2f7b1555883d0bc74abe6",
          "requirement": "landing",
          "command": "bench test --package ./internal/worktree --run 'LandCommandReportsInstallStepForABrokerChangingDiff|LandReportsInstallStepForADeletedBrokerInput|LandRefreshesTheBrokerAfterPublication|LandReportsAFailedRefresh|LandSkipsAFreshBroker|RefreshHoldsAnInstalledCandidate'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make DiffChangesBuildInputs return false for every diff, and record that TestLandCommandReportsInstallStepForABrokerChangingDiff turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep1",
              "digest": "sha256:544cc0f299db6828ac6e615c6a517b0c19cf9755d2edcfe3b2e4c70bc3ec858d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_buildinputs.go,swap,failed,2,TestLandCommandReportsInstallStepForABrokerChangingDiff,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,LandCommandReportsInstallStepForABrokerChangingDiff|LandReportsInstallStepForADeletedBrokerInput|LandRefreshesTheBrokerAfterPublication|LandReportsAFailedRefresh|LandSkipsAFreshBroker|RefreshHoldsAnInstalledCandidate,passed,13\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,909\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,2895,13\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:91: landing printed no broker install notice: \\\"landing source{review_base=add5401c323abcefff2591fb4e9342ea42dc0a7a,assignment_start=52e479bf21b0c5b4bb248314d6ae540d03c1b716}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (810 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsInstallStepForADeletedBrokerInput,\"land_broker_notice_test.go:71: landing printed no broker install notice: \\\"landing source{review_base=f0371836265b35e06cd44f55cba20432a23d55af,assignment_start=ccfefc9e773f7a60ca2975a19f5603d434b44535}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (770 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "cw-c1-rep2-owners",
          "performer": "ft406-author-01-rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4de49d6a6a65bc49342369d3735bf52383d138e9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep2",
            "digest": "sha256:e2735729d25e06358faac34ecd1ed429e29a6f4b43ecdadac60098568a3f5eac",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,DiffChangesBuildInputs|Refresh,passed,15\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,435\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,1112,15\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,\"decision_test.go:135: Refresh(nil build that publishes nothing) = (true, \u003cnil>, \u003cnil>), want no rebuild, build error \u003cnil>, and a seal refusal\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "856141b4132375c0b870c48ea8c89d7f7314d809",
          "requirement": "owners",
          "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs|Refresh'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make Refresh return rebuilt true after a nil build error without the second Verify read, and record that TestRefreshSecondReadDecides turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep2",
              "digest": "sha256:e2735729d25e06358faac34ecd1ed429e29a6f4b43ecdadac60098568a3f5eac",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,DiffChangesBuildInputs|Refresh,passed,15\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,435\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,1112,15\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestRefreshSecondReadDecides/nil_build_that_publishes_nothing,\"decision_test.go:135: Refresh(nil build that publishes nothing) = (true, \u003cnil>, \u003cnil>), want no rebuild, build error \u003cnil>, and a seal refusal\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "cw-c1-rep2-deleted-input",
          "performer": "ft406-author-01-rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4de49d6a6a65bc49342369d3735bf52383d138e9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep2",
            "digest": "sha256:bac258d2145b8ee8e554b913fd24baa5d00436b5e9ae1c44d4491156e64d7789",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestDiffChangesBuildInputs/deleted_input,8,1,0,null,yes\n  2,bit,omit,failed,1,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,8,1,0,null,yes\n  3,bit,omit,failed,4,TestDiffChangesBuildInputs/build_script,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "856141b4132375c0b870c48ea8c89d7f7314d809",
          "requirement": "deleted-input",
          "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the deleted-file test from DiffChangesBuildInputs so that it reads only the post-commit input set, and record that TestDiffChangesBuildInputs turns red on the deleted-input case.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep2",
              "digest": "sha256:bac258d2145b8ee8e554b913fd24baa5d00436b5e9ae1c44d4491156e64d7789",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestDiffChangesBuildInputs/deleted_input,8,1,0,null,yes\n  2,bit,omit,failed,1,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,8,1,0,null,yes\n  3,bit,omit,failed,4,TestDiffChangesBuildInputs/build_script,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "cw-c1-rep2-renamed-input",
          "performer": "ft406-author-01-rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4de49d6a6a65bc49342369d3735bf52383d138e9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep2",
            "digest": "sha256:bac258d2145b8ee8e554b913fd24baa5d00436b5e9ae1c44d4491156e64d7789",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestDiffChangesBuildInputs/deleted_input,8,1,0,null,yes\n  2,bit,omit,failed,1,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,8,1,0,null,yes\n  3,bit,omit,failed,4,TestDiffChangesBuildInputs/build_script,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "856141b4132375c0b870c48ea8c89d7f7314d809",
          "requirement": "renamed-input",
          "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove --no-renames from the diff that DiffChangesBuildInputs reads, and record that TestDiffChangesBuildInputs turns red on the renamed-input case.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep2",
              "digest": "sha256:bac258d2145b8ee8e554b913fd24baa5d00436b5e9ae1c44d4491156e64d7789",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestDiffChangesBuildInputs/deleted_input,8,1,0,null,yes\n  2,bit,omit,failed,1,TestDiffChangesBuildInputs/input_renamed_to_a_test_file,8,1,0,null,yes\n  3,bit,omit,failed,4,TestDiffChangesBuildInputs/build_script,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/freshness,DiffChangesBuildInputs,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "cw-c1-rep2-landing",
          "performer": "ft406-author-01-rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4de49d6a6a65bc49342369d3735bf52383d138e9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep2",
            "digest": "sha256:812b3f22002cee854ea1297148b353757bbd729f7ed12706ef2dfb9399be6b0f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_buildinputs.go,swap,failed,2,TestLandCommandReportsInstallStepForABrokerChangingDiff,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,LandCommandReportsInstallStepForABrokerChangingDiff|LandReportsInstallStepForADeletedBrokerInput|LandRefreshesTheBrokerAfterPublication|LandReportsAFailedRefresh|LandSkipsAFreshBroker|RefreshHoldsAnInstalledCandidate,passed,13\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,909\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,3046,13\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:91: landing printed no broker install notice: \\\"landing source{review_base=286aa424b9cd5147de98781d01b397c3fdadb611,assignment_start=7c696fa6ad848a48ac73d984901da0b48975572a}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (810 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsInstallStepForADeletedBrokerInput,\"land_broker_notice_test.go:71: landing printed no broker install notice: \\\"landing source{review_base=70bfe9b6515706cf55678944e47aeedbc88fa8d7,assignment_start=47bc657c4619f4aec153ae51cfd31f918fd7b81f}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (770 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "856141b4132375c0b870c48ea8c89d7f7314d809",
          "requirement": "landing",
          "command": "bench test --package ./internal/worktree --run 'LandCommandReportsInstallStepForABrokerChangingDiff|LandReportsInstallStepForADeletedBrokerInput|LandRefreshesTheBrokerAfterPublication|LandReportsAFailedRefresh|LandSkipsAFreshBroker|RefreshHoldsAnInstalledCandidate'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make DiffChangesBuildInputs return false for every diff, and record that TestLandCommandReportsInstallStepForABrokerChangingDiff turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft406-author-01-rep2",
              "digest": "sha256:812b3f22002cee854ea1297148b353757bbd729f7ed12706ef2dfb9399be6b0f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_buildinputs.go,swap,failed,2,TestLandCommandReportsInstallStepForABrokerChangingDiff,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,LandCommandReportsInstallStepForABrokerChangingDiff|LandReportsInstallStepForADeletedBrokerInput|LandRefreshesTheBrokerAfterPublication|LandReportsAFailedRefresh|LandSkipsAFreshBroker|RefreshHoldsAnInstalledCandidate,passed,13\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,909\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,3046,13\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:91: landing printed no broker install notice: \\\"landing source{review_base=286aa424b9cd5147de98781d01b397c3fdadb611,assignment_start=7c696fa6ad848a48ac73d984901da0b48975572a}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (810 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsInstallStepForADeletedBrokerInput,\"land_broker_notice_test.go:71: landing printed no broker install notice: \\\"landing source{review_base=70bfe9b6515706cf55678944e47aeedbc88fa8d7,assignment_start=47bc657c4619f4aec153ae51cfd31f918fd7b81f}\\\\\\\\nlanding refresh failed: bench binary \\\\\\\\… (770 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "cw-c1-r1-standards",
          "performer": "ft406-review-c1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "5d93cb2bfbb1812bf823fb3c46c8381a8efa8fa3",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft406-review-c1",
            "digest": "sha256:f4cb68d25d72080e731af4b897430315dd46a3f17b2480abe8e3bbc3627474b0",
            "excerpt": "Standards axis (unified session), CW-C1, frozen pair 8fcad834..2a5c53e3. Count: 1 finding. Worst: S1.\nS1, confidence 3, internal/worktree/land_effects.go:137-138: the sentence \"Refresh settles completion from the tree, so a resume needs no journal\" repeats the completion rationale that freshness.Refresh's doc comment owns. no-op.\nNo hard violation: buildInputs returns the inputs and the package directories from one go list call; relativeBuildInputs feeds BuildInputs and the diff owner; no second decoder; no new dependency.\nCommand contribution: none.\n"
          },
          "axis": "Standards",
          "base": "8fcad8341e760744d06ce632ea0a978684025f81",
          "tip": "2a5c53e3558ed3e39a1720d242a17008a4f1f3ae",
          "finding_ids": [
            "cw-c1-std-1"
          ],
          "supersedes": []
        },
        {
          "id": "cw-c1-r1-spec",
          "performer": "ft406-review-c1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "5d93cb2bfbb1812bf823fb3c46c8381a8efa8fa3",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft406-review-c1",
            "digest": "sha256:9bd59e1632e06ad5e6698fcc6147d97132a65065f39a8110197c95f2abe57e1d",
            "excerpt": "Spec axis (unified session), CW-C1, frozen pair 8fcad834..2a5c53e3. Count: 1 finding. Worst: P1.\nP1, confidence 9, internal/worktree/land_effects_test.go:183: CW46 is unclosed. TestLandReportsAFailedRefresh uses refreshJoins(nil), whose build returns nil, so no test grades the text \"landing refresh build failed: \". auto-fix: a failing-build assertion on that literal.\nEvery other CW-C1 row holds; CW10 passes (rg name-only finds only a comment at land.go:182).\nCommand contribution: the ticket acceptance line bundles six kept rows under \"keep or gain the answers\"; /bench-write-spec should require a recorded red for each kept-rule row mapped to an existing test.\n"
          },
          "axis": "Spec",
          "base": "8fcad8341e760744d06ce632ea0a978684025f81",
          "tip": "2a5c53e3558ed3e39a1720d242a17008a4f1f3ae",
          "finding_ids": [
            "cw-c1-spec-1"
          ],
          "supersedes": []
        },
        {
          "id": "cw-c1-r1-coverage",
          "performer": "ft406-review-c1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "5d93cb2bfbb1812bf823fb3c46c8381a8efa8fa3",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft406-review-c1",
            "digest": "sha256:0a9d7a7f13f32dd2213ad9ca0acc08471d38b0edd59703af35310138f8c4307e",
            "excerpt": "Coverage axis (unified session), CW-C1, frozen pair 8fcad834..2a5c53e3. Count: 2 findings. Worst: C1.\nC1, confidence 9, internal/worktree/land_effects.go:156: the CW46 gap; a swap of the \"landing refresh build failed: \" text stays silent (14 landing tests pass, restored yes). auto-fix, same repair as P1.\nC2, confidence 6, internal/freshness/freshness_buildinputs.go:281,285-291: --name-status is parsed without -z, so Git C-quotes a non-ASCII deleted path such as cmd/bench/é.go, the .go suffix test fails, and the owner answers false. No row or Won't handle covers it. ask-user (spec.md:137 fixes the command).\nProbes: build-failure text swap silent; _test.go exclusion omit bit; Refresh drops buildErr bit; status D to Z bit; first seal read disabled bit. Every restored cell yes.\nCommand contribution: none.\n"
          },
          "axis": "Coverage",
          "base": "8fcad8341e760744d06ce632ea0a978684025f81",
          "tip": "2a5c53e3558ed3e39a1720d242a17008a4f1f3ae",
          "finding_ids": [
            "cw-c1-cov-1",
            "cw-c1-cov-2"
          ],
          "supersedes": []
        },
        {
          "id": "cw-c1-r2-standards",
          "performer": "ft406-review-c1-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4de49d6a6a65bc49342369d3735bf52383d138e9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft406-review-c1-r2",
            "digest": "sha256:e0428dfcb7d3e1427f097645dd139918a88a24669adb9facd01ab24106691251",
            "excerpt": "Standards axis (unified confirming session), CW-C1, repair delta 2a5c53e3..856141b4. Count: 0 findings.\nO1 closed: git.ParseNameStatusZ at internal/git/status.go:280-295 is the one owner; internal/diff/snapshot.go:43,94 and internal/freshness/freshness_buildinputs.go:136 call it; the moved body is byte-identical; git.Raw errors still answer false (freshness_buildinputs.go:132-135).\nAdvice: internal/freshness/decision_test.go:41 doc comment breaks early (rewrap).\nCommand contribution: none.\n"
          },
          "axis": "Standards",
          "base": "8fcad8341e760744d06ce632ea0a978684025f81",
          "tip": "856141b4132375c0b870c48ea8c89d7f7314d809",
          "finding_ids": [],
          "supersedes": [
            "cw-c1-r1-standards"
          ]
        },
        {
          "id": "cw-c1-r2-spec",
          "performer": "ft406-review-c1-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4de49d6a6a65bc49342369d3735bf52383d138e9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft406-review-c1-r2",
            "digest": "sha256:b9ebda5820eb608d862d101861d2e35845996dbf5b66e18222814c6f98d67672",
            "excerpt": "Spec axis (unified confirming session), CW-C1, repair delta 2a5c53e3..856141b4. Count: 0 findings.\nP1 closed: land_effects_test.go:184-188 fails the build with \"compile exploded\" and checks failed and the exact text \"landing refresh build failed: compile exploded\\n\" (land_effects.go:156, spec row CW46).\nC2/CW51 closed: spec.md:137-138 and row CW51 name git diff -z --no-renames --name-status, matching freshness_buildinputs.go:132; subtest \"deleted non-ASCII input\".\nEvery delta path is on ticket 01's widened Writes line; nothing outside the acceptance rows.\nCommand contribution: none.\n"
          },
          "axis": "Spec",
          "base": "8fcad8341e760744d06ce632ea0a978684025f81",
          "tip": "856141b4132375c0b870c48ea8c89d7f7314d809",
          "finding_ids": [],
          "supersedes": [
            "cw-c1-r1-spec"
          ]
        },
        {
          "id": "cw-c1-r2-coverage",
          "performer": "ft406-review-c1-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4de49d6a6a65bc49342369d3735bf52383d138e9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft406-review-c1-r2",
            "digest": "sha256:3b5a031b1d2c85c3de518ffa807ea064a7015c2f12d606d8546193fc529cfcfd",
            "excerpt": "Coverage axis (unified confirming session), CW-C1, repair delta 2a5c53e3..856141b4. Count: 0 findings.\nC1 closed: probe land_effects.go swap \"if buildErr != nil {\" -> \"if buildErr == nil {\": bit, TestLandReportsAFailedRefresh, restored yes.\nC2 closed: probe freshness_buildinputs.go guard that skips a non-ASCII name in the deletion test: bit, only TestDiffChangesBuildInputs/deleted_non-ASCII_input, restored yes.\nNo new hostile-input gap: NUL framing keeps tab, newline, and quote bytes; an odd trailing segment drops.\nCommand contribution: none.\n"
          },
          "axis": "Coverage",
          "base": "8fcad8341e760744d06ce632ea0a978684025f81",
          "tip": "856141b4132375c0b870c48ea8c89d7f7314d809",
          "finding_ids": [],
          "supersedes": [
            "cw-c1-r1-coverage"
          ]
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
      "from": "sha256:caecfd37289f7819cd30ab7da133e5019f36c7a65df862ec6a16d535b9272523",
      "to": "sha256:086b266147d63e0ac6dbbf355ab726e1f706ead33525d9afba8501b37d83da54",
      "chunk_ids": {
        "CW-C1": [
          "CW-C1"
        ]
      }
    },
    {
      "from": "sha256:086b266147d63e0ac6dbbf355ab726e1f706ead33525d9afba8501b37d83da54",
      "to": "sha256:43bcc719cba9b1bf8d9efe50be6ffc20fb51fe361d8084a238deb7e5996fbb1e",
      "chunk_ids": {
        "CW-C1": [
          "CW-C1"
        ]
      }
    }
  ]
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

## CW-C1 review round 1

One unified session (`execution.review_mode` is `unified`) ran the three axes separately on opus at medium effort over `8fcad834..2a5c53e3`.
The orchestrator accepts P1, C1, and C2 under `--auto-approve`, and dismisses S1.
P1 and C1 name one repair target, and C2 names a second, so the round has two repair targets.
Both targets are on paths of ticket 01, so one repair session takes both. C2 changes production code, so CW-C1 takes a confirming round after the repair.
The spec amendment for C2 adds `-z` to the diff of `DiffChangesBuildInputs` and adds row CW51 to ticket 01.

## Standards

Findings: 1, dismissed. Worst issue: S1.

- S1 (`cw-c1-std-1`), no-op, confidence 3: `internal/worktree/land_effects.go:137-138` states that `freshness.Refresh` settles completion. The sentence names the owner and is a short pointer, not a second copy of the rationale.

## Spec

Findings: 1 accepted. Worst issue: P1.

- P1 (`cw-c1-spec-1`), auto-fix, confidence 9: CW46 is unclosed. `TestLandReportsAFailedRefresh` (`internal/worktree/land_effects_test.go:183`) uses a build that returns nil, so no test grades the text `landing refresh build failed: `. Add a failing-build assertion on that text.

## Coverage

Findings: 2 accepted. Worst issue: C1.

- C1 (`cw-c1-cov-1`), auto-fix, confidence 9: the CW46 gap at `internal/worktree/land_effects.go:156`. A swap of the build-failure text stays silent. The repair of P1 closes it.
- C2 (`cw-c1-cov-2`), auto-fix, confidence 6: `internal/freshness/freshness_buildinputs.go:281,285-291` parses `--name-status` without `-z`. Git quotes a non-ASCII deleted path, so the `.go` suffix test misses it. Add `-z`, parse the NUL-separated fields, and add the CW51 case to `TestDiffChangesBuildInputs`.

## CW-C1 repair cycle 1 verification

The repair session `ft406-author-01-rep1` (opus / medium) committed `8fc874450` and refroze CW-C1 at `667622a78`.
CW46 went red when the build-failure text was swapped, and CW51 went red when `-z` was removed. Each restored cell reads `yes`.
The `deleted-input` plan probe ran as `if false && status == "D"`, because a literal removal of the block does not compile. The orchestrator accepts that form, because it removes the deleted-file test as the plan states.

The orchestrator found one defect in the repair delta. `DiffChangesBuildInputs` now parses NUL-separated `--name-status` pairs itself, and `parseNameStatusZ` in `internal/diff/snapshot.go` already owns that parse. `AGENTS.md` grades duplicated knowledge as a defect.

- O1 (`cw-c1-orch-1`), auto-fix: move the parser to one exported owner in `internal/git/status.go`, beside the porcelain parser, and make both callers read it. Repair cycle 2 of 2 takes this finding. The ticket 01 widen cites `--auto-approve`.

## CW-C1 confirming round

A fresh unified session ran the three axes on opus at medium effort over the repair delta `2a5c53e3..856141b4`.
It closed P1, C1, C2, and O1, and each axis passes. CW-C1 used both repair cycles of its allowance.
Its probes used new sites: `buildErr == nil` in `refreshBroker` bit `TestLandReportsAFailedRefresh`, and a non-ASCII skip in the deletion test bit only the CW51 case. Each restored cell reads `yes`.

Advice: the doc comment at `internal/freshness/decision_test.go:41` breaks early and can take a rewrap.
