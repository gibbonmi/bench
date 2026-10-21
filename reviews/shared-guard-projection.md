# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/shared-guard-projection/spec.md",
  "plan_digest": "sha256:26860796c78f035de7c2a1e29d67091b9d1200c6a543fe521f3f02df7be53ea7",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GP1",
      "base": "f86203658f5190beac7a34ce17e18473d20fe426",
      "tip": "93e90d3ee55ea6adb6d713389077326e688fec12",
      "plan_digest": "sha256:ed266643ad636c39c00201d70cd717be17e818a9770d0a1ee83d3b5f690b8df7",
      "source_digest": "09cf33d9e84ed49744b212f5792a0c2ea4a4bfbd",
      "acceptance_rows": [
        "GP01",
        "GP02",
        "GP03",
        "GP05",
        "GP06",
        "GP07",
        "GP08",
        "GP09",
        "GP10",
        "GP72",
        "GP74"
      ],
      "verification": [
        {
          "id": "t1-git-context-omission",
          "performer": "claude:sgp-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9176090363f3e82572626823207983480c37cc54",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
            "digest": "sha256:e85efeaaca93794962a7490c9b5c6803a02647316225ec2097caa0d07f973a3d",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGitProjectionAuthority/child_push,4,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gitguard,^TestGitProjectionAuthority$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,442\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,fail,2,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gitguard,TestGitProjectionAuthority/child_push,\"projection_contract_test.go:22: Classify(\\\"bench worktree exec X -- git push origin topic\\\") = \\\"\\\", want \\\"git push with an unresolved destination\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-git-context-omission",
          "command": "bench test --package ./internal/gitguard --run '^TestGitProjectionAuthority$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit only DifferentCheckout propagation at the Git consumer, keeping WorktreeExec and its other facts; child push must red the existing unresolved-destination assertion.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
              "digest": "sha256:e85efeaaca93794962a7490c9b5c6803a02647316225ec2097caa0d07f973a3d",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGitProjectionAuthority/child_push,4,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gitguard,^TestGitProjectionAuthority$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,442\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,fail,2,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gitguard,TestGitProjectionAuthority/child_push,\"projection_contract_test.go:22: Classify(\\\"bench worktree exec X -- git push origin topic\\\") = \\\"\\\", want \\\"git push with an unresolved destination\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-wrapper-projection-omission",
          "performer": "claude:sgp-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9176090363f3e82572626823207983480c37cc54",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
            "digest": "sha256:675dc7d93bf474a216c9ef8f9c765403f9ee9e208c6019b642960d2c7bfbe6bc",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestShellChildren/zsh,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/shellcommand,^TestShellChildren$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,387\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,2,8\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestShellChildren/zsh,\"projection_test.go:47: ShellChildren([\\\"zsh\\\" \\\"-c\\\" \\\"git status\\\"], 0) = [], want [{2 [git status]}]\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-wrapper-projection-omission",
          "command": "bench test --package ./internal/shellcommand --run '^TestShellChildren$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit a supported shell basename or c-flag projection while preserving the API; the independent ordered child cases must red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
              "digest": "sha256:675dc7d93bf474a216c9ef8f9c765403f9ee9e208c6019b642960d2c7bfbe6bc",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestShellChildren/zsh,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/shellcommand,^TestShellChildren$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,387\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,2,8\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestShellChildren/zsh,\"projection_test.go:47: ShellChildren([\\\"zsh\\\" \\\"-c\\\" \\\"git status\\\"], 0) = [], want [{2 [git status]}]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-t1-projection",
          "performer": "claude:sgp-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9176090363f3e82572626823207983480c37cc54",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
            "digest": "sha256:906e3d92e5bf32166bc67f9b597ff520a4e077fd876ef22eec6119207aad9fe9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,pass,2,42\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-projection",
          "command": "bench test --package ./internal/shellcommand",
          "exit_code": 0
        },
        {
          "id": "t1-t1-git-policy",
          "performer": "claude:sgp-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9176090363f3e82572626823207983480c37cc54",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
            "digest": "sha256:db110083248fe95dd8266bae0b6dd13cd553ee59e84ac69c04aa0c8a27be4f4e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,pass,15171,237\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-git-policy",
          "command": "bench test --package ./internal/gitguard",
          "exit_code": 0
        },
        {
          "id": "t1-t1-root-conformance",
          "performer": "claude:sgp-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9176090363f3e82572626823207983480c37cc54",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
            "digest": "sha256:bab2fadca06a406d5fe9fceaca92860518ea32d289891e14ad3b0ba7e0a98140",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12689,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-t1-git-context-omission",
          "performer": "claude:sgp-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9176090363f3e82572626823207983480c37cc54",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
            "digest": "sha256:e85efeaaca93794962a7490c9b5c6803a02647316225ec2097caa0d07f973a3d",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGitProjectionAuthority/child_push,4,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gitguard,^TestGitProjectionAuthority$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,442\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,fail,2,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gitguard,TestGitProjectionAuthority/child_push,\"projection_contract_test.go:22: Classify(\\\"bench worktree exec X -- git push origin topic\\\") = \\\"\\\", want \\\"git push with an unresolved destination\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-git-context-omission",
          "command": "bench test --package ./internal/gitguard --run '^TestGitProjectionAuthority$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit only DifferentCheckout propagation at the Git consumer, keeping WorktreeExec and its other facts; child push must red the existing unresolved-destination assertion.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
              "digest": "sha256:e85efeaaca93794962a7490c9b5c6803a02647316225ec2097caa0d07f973a3d",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGitProjectionAuthority/child_push,4,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gitguard,^TestGitProjectionAuthority$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,442\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,fail,2,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gitguard,TestGitProjectionAuthority/child_push,\"projection_contract_test.go:22: Classify(\\\"bench worktree exec X -- git push origin topic\\\") = \\\"\\\", want \\\"git push with an unresolved destination\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-t1-wrapper-projection-omission",
          "performer": "claude:sgp-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9176090363f3e82572626823207983480c37cc54",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
            "digest": "sha256:675dc7d93bf474a216c9ef8f9c765403f9ee9e208c6019b642960d2c7bfbe6bc",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestShellChildren/zsh,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/shellcommand,^TestShellChildren$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,387\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,2,8\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestShellChildren/zsh,\"projection_test.go:47: ShellChildren([\\\"zsh\\\" \\\"-c\\\" \\\"git status\\\"], 0) = [], want [{2 [git status]}]\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-wrapper-projection-omission",
          "command": "bench test --package ./internal/shellcommand --run '^TestShellChildren$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit a supported shell basename or c-flag projection while preserving the API; the independent ordered child cases must red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t1-author-20261021",
              "digest": "sha256:675dc7d93bf474a216c9ef8f9c765403f9ee9e208c6019b642960d2c7bfbe6bc",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestShellChildren/zsh,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/shellcommand,^TestShellChildren$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,387\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,2,8\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestShellChildren/zsh,\"projection_test.go:47: ShellChildren([\\\"zsh\\\" \\\"-c\\\" \\\"git status\\\"], 0) = [], want [{2 [git status]}]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r1-t1-projection",
          "performer": "claude:sgp-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09cf33d9e84ed49744b212f5792a0c2ea4a4bfbd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-repair1-20261021",
            "digest": "sha256:12fab91f8bbfbe21da8fc58a08044d2b6dc98719f843ae8c6437466e6ec2c26a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,pass,2,43\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-projection",
          "command": "bench test --package ./internal/shellcommand",
          "exit_code": 0
        },
        {
          "id": "r1-t1-git-policy",
          "performer": "claude:sgp-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09cf33d9e84ed49744b212f5792a0c2ea4a4bfbd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-repair1-20261021",
            "digest": "sha256:55e7668a2a17c8b0c2eac162922b51bd3e00d019e32f21b078d7f82ccb9d4511",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,pass,15165,237\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-git-policy",
          "command": "bench test --package ./internal/gitguard",
          "exit_code": 0
        },
        {
          "id": "r1-t1-root-conformance",
          "performer": "claude:sgp-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09cf33d9e84ed49744b212f5792a0c2ea4a4bfbd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-repair1-20261021",
            "digest": "sha256:bb9cc91367ffdf4ed92d92136b1a5e24a67f4134a93b46c9463b9e1e26a6d980",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12826,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "r1-t1-git-context-omission",
          "performer": "claude:sgp-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09cf33d9e84ed49744b212f5792a0c2ea4a4bfbd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-repair1-20261021",
            "digest": "sha256:e85efeaaca93794962a7490c9b5c6803a02647316225ec2097caa0d07f973a3d",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGitProjectionAuthority/child_push,4,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gitguard,^TestGitProjectionAuthority$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,442\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,fail,2,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gitguard,TestGitProjectionAuthority/child_push,\"projection_contract_test.go:22: Classify(\\\"bench worktree exec X -- git push origin topic\\\") = \\\"\\\", want \\\"git push with an unresolved destination\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-git-context-omission",
          "command": "bench test --package ./internal/gitguard --run '^TestGitProjectionAuthority$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit only DifferentCheckout propagation at the Git consumer, keeping WorktreeExec and its other facts; child push must red the existing unresolved-destination assertion.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t1-repair1-20261021",
              "digest": "sha256:e85efeaaca93794962a7490c9b5c6803a02647316225ec2097caa0d07f973a3d",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGitProjectionAuthority/child_push,4,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gitguard,^TestGitProjectionAuthority$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,442\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,fail,2,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gitguard,TestGitProjectionAuthority/child_push,\"projection_contract_test.go:22: Classify(\\\"bench worktree exec X -- git push origin topic\\\") = \\\"\\\", want \\\"git push with an unresolved destination\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r1-t1-wrapper-projection-omission",
          "performer": "claude:sgp-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09cf33d9e84ed49744b212f5792a0c2ea4a4bfbd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t1-repair1-20261021",
            "digest": "sha256:d2867cdc02e532debb96d5818c727a9e492a2d1b111d5c46cef4800b619b5048",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestShellChildren/zsh,9,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/shellcommand,^TestShellChildren$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,387\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,2,9\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestShellChildren/zsh,\"projection_test.go:49: ShellChildren([\\\"zsh\\\" \\\"-c\\\" \\\"git status\\\"], 0) = [], want [{2 [git status]}]\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-wrapper-projection-omission",
          "command": "bench test --package ./internal/shellcommand --run '^TestShellChildren$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit a supported shell basename or c-flag projection while preserving the API; the independent ordered child cases must red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t1-repair1-20261021",
              "digest": "sha256:d2867cdc02e532debb96d5818c727a9e492a2d1b111d5c46cef4800b619b5048",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestShellChildren/zsh,9,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/shellcommand,^TestShellChildren$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,387\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,2,9\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestShellChildren/zsh,\"projection_test.go:49: ShellChildren([\\\"zsh\\\" \\\"-c\\\" \\\"git status\\\"], 0) = [], want [{2 [git status]}]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gp1-r1-standards",
          "performer": "claude:sgp-review-standards-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9176090363f3e82572626823207983480c37cc54",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-standards-r1",
            "digest": "sha256:252191d6a3410a14261398c5304a9b2d263060a9ff0b5369f4b77fc3daf63444",
            "excerpt": "STANDARDS axis, GP1 (base f86203658, code tip a4053ed70, evidence tip f29f599e2). 2 findings. Worst: F2, a redirected-push test row that repeats an existing row and has no recorded red.\n\nFindings\n- F1 | auto-fix | internal/gitguard/projection_contract_test.go:11-13 | The doc comment of TestGitProjectionAuthority is one 29-word sentence. Binding: craft-comments (SKILL.md:12) requires ASD-STE100 prose, and your charge sets a limit of 25 words per sentence. Fix: split it into two sentences. This is a comment-only change, so it needs no repair cycle. | confidence 8\n- F2 | ask-user | internal/gitguard/projection_contract_test.go:22 (the \"redirected push\" row) | The row is the same tuple (cmd, Checker, want) as verdict_test.go:96. Binding: AGENTS.md \"one source per fact\", test-expectation exception. That exception needs a necessary independence plus a recorded and demonstrated red. The only recorded probe is t1-git-context-omission (reviews/shared-guard-projection.md:31-56, on f29f599e2). It mutates DifferentCheckout, so it bites only child_push (:42). Nothing records a red for the redirected row. Omitting redirect handling in the untouched redirectsRepository would also redden verdict_test.go:96-98, so this row's independence is not shown to be necessary. Counter-weight: spec GP10 (spec.md:371 in the diff) and ticket acceptance line 33 name \"child and redirected push\" under this seam. The reviewer must choose between two fixes. Either record a redirect-omission red for this row, or accept it as a spec-bound row and amend GP10. | confidence 5\n\nNamed question: child push versus verdict_test.go:239\n- Decision: the child-push row meets the exception, so it is not a defect.\n- Condition 2 is met. The recorded probe (reviews/shared-guard-projection.md:42, :55) shows \"bit\" on TestGitProjectionAuthority/child_push with restore \"pass\".\n- Condition 1 is partly met. The omission also reddens verdict_test.go:239, so independence is not strictly necessary for a red gate. But ticket acceptance line 33 names TestGitProjectionAuthority as the test that must go red. The spec binds the duplicate, so I treat it as a bound witness and not as drift.\n- The redirected-push row does not get the same treatment (see F2).\n\nChecked with no finding\n- No second wrapper or exec-child grammar remains in internal/gitguard. In the diff, `wrappers`, `wrapperCFlagRe`, the `regexp` import and `worktreeExecChild` are removed. They now exist once, in internal/shellcommand/projection.go:10-12 and :29-37. benchExecutables stays in scan.go:78 as the caller-owned Bench-identity test, which the ticket requires.\n- Behaviour parity of the scan.go migration, from traced old and new code: the first `-c` child matches the old first-flag-then-break behaviour; an empty exec child yields \"\" in both versions; the case order is unchanged.\n- stdlib rule: projection.go uses slices.Contains and slices.Index. No hand-written loop that a standard-library function replaces.\n- Error handling: ShellChildren and WorktreeExec check bounds (projection.go:24, :50) and return no child on a missing operand or separator.\n- Comments: every other new or edited sentence is 25 words or fewer. No narrated history and no FT tags.\n- Naming fits the surrounding code.\n\nAdvice (no finding ID)\n- internal/shellcommand/projection_test.go:288: the loop variable `child` shadows the package-level type `child` (projection_test.go:228). It compiles, but a rename would read more clearly.\n- internal/shellcommand/projection.go:12: `wrapperCFlagRe` has no comment. It was also uncommented in its old home, so this is not a regression.\n- specs/shared-guard-projection/spec.md was edited in the diff but is not on ticket 1's `Writes:` line. The edit only turns `planned` markers into real test citations, which map-discipline.md:126 requires when the covering ticket commits.\n\nCommands (all via bench worktree exec shared-guard-projection-build --): bench preflight evidence sha256:572e5979... --check-current returned current=true, dirty_record=false; git diff f86203658 a4053ed70 (--stat and full); git show of the ticket, scan.go, verdict_test.go, and the review record; git grep over a4053ed70 and f29f599e2; git log f86203658..a4053ed70 --stat. No tests, no probe, no file written.\n"
          },
          "axis": "Standards",
          "base": "f86203658f5190beac7a34ce17e18473d20fe426",
          "tip": "a4053ed70c625935cd64c2d4fe5943f99348d88c",
          "finding_ids": [
            "F1",
            "F2"
          ],
          "supersedes": []
        },
        {
          "id": "gp1-r1-spec",
          "performer": "claude:sgp-review-spec-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9176090363f3e82572626823207983480c37cc54",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-spec-r1",
            "digest": "sha256:7a15290d0cf3d4bc8819690271abf9bb201cd6abb455fcaee7e83d329f720261",
            "excerpt": "SPEC axis, chunk GP1 (base f86203658, code tip a4053ed70, record tip f29f599e): no findings.\n\nFindings: none. The reviewer needs no routing decision for this axis.\n\nCount: 0 findings. Worst issue: none.\n\nWhat I checked, all by re-reading source at the tip:\n- API shapes. WrapperChild, ShellChildren, ExecProjection and WorktreeExec match the spec's \"Projection API and caller policy\" signatures exactly.\n  - ShellChildren: basenames sh/bash/zsh via slices.Contains, regex ^-[A-Za-z]*c[A-Za-z]*$, argv order, no recursion. A negative or too-large index gives nil. A flag in the last position gives no child, and each position is read independently (GP02/GP03, family W).\n  - WorktreeExec: adjacent worktree and exec, first later -- via slices.Index, child returned unchanged with no join or reparse. HasSeparator is true even for an empty child. DifferentCheckout is len(Child) > 0. It trusts the caller's Bench identity, as the spec says.\n- Git caller policy (scan.go:44-71). The keyword trim stays in scan. case benchExecutables[base] is the Bench basename test, and it runs before WorktreeExec. The default branch takes only children[0], so the first matching shell child is consumed even when it finds no Git. Both branches recurse with allowWrapper=false. exec.DifferentCheckout is passed as elsewhere. scan and scanWords signatures are unchanged. Edge cases bash -c, bash -c -c 'git status', an exec with no --, and an empty child after -- give the same result as before.\n- Migration table scan.go rows. wrappers, wrapperCFlagRe, worktreeExecChild and the regexp import are gone from gitguard. shellcommand.go is untouched. gitguard holds 11 of 12 allowed files. The new files are 61 and 129 lines.\n- Coverage rows GP01-GP03, GP05-GP10, GP72, GP74 each cite a real test; GP10 red-bite recorded at reviews/shared-guard-projection.md:42.\n- Spec cell rewrites: limited to nine cells (GP02, GP03, GP05-GP08, GP10, GP72, GP74); each cites a test that lands in this commit, as .agents/commands/bench-implement-spec.md:52 requires. .bench/commitment.json is untouched.\n- Review record duplicate and doubled-prefix rows: bench record ticket builds the evidence ID as prefix + \"-\" + requirement.ID (internal/reviewrecord/recordcmd/ticket.go:143), so t1-t1- is mechanical. checkVerification matches on Requirement and the last match wins (internal/reviewrecord/check.go:46-48). Evidence-ID uniqueness is enforced by unclaimed (change.go:144); the two rows have distinct IDs. Neither breaks a record contract or the chunk checkpoint (confidence 7; parse.go validation not read in full).\n\nAdvice (optional, no finding ID):\n- Plan requirement IDs that begin with t1- produce cosmetic t1-t1- evidence IDs.\n- The duplicate probe rows are redundant evidence. No binding requirement forces their removal.\n\nCommands run, all via bench worktree exec shared-guard-projection-build --: bench preflight evidence sha256:572e5979...88 --check-current returned current=true, dirty_record=false; git diff --stat, git diff, git show of the spec, ticket, scan.go and the review record; rg over source. No tests, no bench probe, no files written.\n"
          },
          "axis": "Spec",
          "base": "f86203658f5190beac7a34ce17e18473d20fe426",
          "tip": "a4053ed70c625935cd64c2d4fe5943f99348d88c",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "gp1-r1-coverage",
          "performer": "claude:sgp-review-coverage-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9176090363f3e82572626823207983480c37cc54",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-coverage-r1",
            "digest": "sha256:2bd2e2f6b498a1f40370884d0282067300188200f1ff3e27fa8ba293ca62efe2",
            "excerpt": "COVERAGE axis, chunk GP1 (base f86203658, code tip a4053ed70, evidence tip f29f599e). Findings: 2, both auto-fix, both from the one hardening cycle the policy permits. The worst is C2. No acceptance row fails.\n\nFindings (one per line)\n- C1 | auto-fix | internal/shellcommand/projection.go:50 (the words[executable+2] != \"exec\" check) and projection_test.go:117 (TestWorktreeExecContext) | Binding: spec.md:96 \"recognizes adjacent worktree and exec words\"; GP06 \"Header recognition is confused with runnable child discovery\" | Confidence 6. Probe 1 swapped the check for !slices.Contains(words[executable+2:], \"exec\"). The verdict was silent: shellcommand and gitguard both passed (42 and 237 tests). Break input: bench worktree list exec -- git reset --hard. The mutant projects a child and Git refuses; the spec says there is no exec head. No test has worktree adjacent but exec non-adjacent with a later --. Fix: add a negative case to TestWorktreeExecContext (or a sibling) that expects ExecProjection{}.\n- C2 | auto-fix | internal/shellcommand/projection.go:30 and projection_test.go:27 (TestShellChildren) | Binding: spec.md:91 \"It performs no recursion\"; W4 edge row \"a nested wrapper ... neither recurses twice\"; GP02 | Confidence 5. Probe 2 changed the append so it also adds the children of the child's own commands (one extra level). The verdict was silent: both packages passed. Git reads only children[0], and GP09 runs through Git, so the inner child appended after the outer one is never seen. No ShellChildren case has a wrapper inside a -c operand, for example bash -c \"sh -c 'git push'\"; the expected result is one child. Fix: add that case to TestShellChildren. Counter-view: the spec assigns nested depth to GP09, not GP02; treat as optional advice if you read it that way.\n\nNamed gaps, checked and cleared\n- GP05: TestWorktreeExec has \"two words\" and --env A=1 before --; the child comes back unchanged.\n- GP06: TestWorktreeExecMissingDelimiter compares the whole struct; a -- before exec is covered by the no-projection case in TestWorktreeExecContext.\n- W family: /bin/bash -lc, bash -c -c 'git status', a missing operand, a trailing -lc, unsupported dash, and both bad indexes are covered; the zsh omission is bitten.\n- GP74: covered; an every-child loop reds (held).\n- GP01 and GP09: pre-existing characterization rows.\n\nAdvice (no ID)\n- The GP10 \"redirected push\" and GP74 \"first reset child\" cases pass on old and new code; their red comes only from the held mutants. The redirect branch is unchanged code.\n- bench worktree exec -- git push (no target) has no test; old and new code agree on it.\n\nCommands and results\n1. bench preflight evidence sha256:572e5979... --check-current: current=true, dirty_record=false.\n2. git diff f86203658 a4053ed70 (and --stat): read.\n3. bench probe --help: usage only.\n4. Probe 1 (adjacency of exec): silent, cause passed, restored=yes.\n5. Probe 2 (ShellChildren one-level recursion): silent, cause passed, restored=yes.\nProbes used: 2 of 2.\n"
          },
          "axis": "Coverage",
          "base": "f86203658f5190beac7a34ce17e18473d20fe426",
          "tip": "a4053ed70c625935cd64c2d4fe5943f99348d88c",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "gp1-r2-standards",
          "performer": "claude:sgp-review-standards-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "09cf33d9e84ed49744b212f5792a0c2ea4a4bfbd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-standards-r2",
            "digest": "sha256:83a767caa89254a57ce8d9afb40fcf2649157c2c796a72221a6e7df9d157ce37",
            "excerpt": "Standards, GP1 round 1 confirmation: F1 and F2 are both closed, and the repair delta adds no new Standards finding.\n\nF1: closed.\n- internal/gitguard/projection_contract_test.go:5-7 (at 93e90d3ee) now holds three sentences of 10, 12 and 7 words.\n- Each sentence is under the 25-word limit, and the comment still states the same contract.\n\nF2: closed by the recorded red, with the row kept.\n- The raw answer file /tmp/claude-1000/-home-mgibs-workspace-bench/f5571912-814f-48c2-9e29-363a829710e2/scratchpad/sgp-t1-r1-redirect.out shows `bit,internal/gitguard/verdict_push.go,swap,failed,2,TestGitProjectionAuthority/child_push,yes`.\n- Its failures list includes TestGitProjectionAuthority/redirected_push at projection_contract_test.go:22: Classify(\"git -C /other push origin topic\") = \"\", want \"git push with an unresolved destination\".\n- The mutated line exists at 93e90d3ee:internal/gitguard/verdict_push.go:19, `if viaXargs || redirected {`.\n- That meets the AGENTS.md test-expectation exception for the row at projection_contract_test.go:16. The omission turns this row red, and the red is recorded and the file restored (restored=yes).\n\nNew findings: none.\n- Repair-delta hunks against a4053ed70, with an unchanged production tree: the comment rewrite in the gitguard test file; the \"nested wrapper stays one child\" case at projection_test.go:43-44; the TestWorktreeExecContext loop that adds {\"bench\",\"worktree\",\"list\",\"exec\",\"--\",\"git\",\"reset\",\"--hard\"}.\n- One source per fact: the new cases test the shellcommand owner directly, with inputs and expectations unlike the gitguard rows.\n- The new cases carry no comments. C1 and C2 are hardening-cycle checks; no blocking condition on them.\n\nOptional advice (no ID): unchanged from round 1 (the child loop variable shadows the type child; wrapperCFlagRe has no comment).\n\nCommands run, all through bench worktree exec shared-guard-projection-build --: bench preflight evidence sha256:5590e889... --check-current returned current=true, dirty_record=false, source_tip ca726066b; git diff a4053ed70 93e90d3ee; git show 93e90d3ee:internal/gitguard/projection_contract_test.go; git grep -n \"viaXargs || redirected\" 93e90d3ee -- internal/gitguard. No test, no probe, no file written.\n"
          },
          "axis": "Standards",
          "base": "f86203658f5190beac7a34ce17e18473d20fe426",
          "tip": "93e90d3ee55ea6adb6d713389077326e688fec12",
          "finding_ids": [],
          "supersedes": [
            "gp1-r1-standards"
          ]
        },
        {
          "id": "gp1-r2-spec",
          "performer": "claude:sgp-review-spec-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "09cf33d9e84ed49744b212f5792a0c2ea4a4bfbd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-spec-r2",
            "digest": "sha256:aa06cd61f9a7e316d333de3abe1ccba1694145970d1db3ffd5f12c4e20053089",
            "excerpt": "SPEC axis, GP1 confirming round 1: PASS. No findings. Count 0; worst issue none. Nothing was run beyond reads (no tests, no probe).\n\nChecks:\n- Repair delta is test-only. Per git diff --stat a4053ed70 93e90d3ee, the Go changes are internal/gitguard/projection_contract_test.go (comment split only) and internal/shellcommand/projection_test.go (+2 cases). The other files in the range are the spec.md plan line (repair assignment record) and the review record. projection.go and scan.go did not change, so the Projection API signatures and the Git caller policy are unchanged from the initial pass.\n- Acceptance rows GP01-03, GP05-10, GP72, GP74 still cite the same tests. No assertion was removed or weakened.\n- C1 (projection_test.go:127) matches spec.md:96 \"recognizes adjacent worktree and exec words\". For {bench,worktree,list,exec,--,git,reset,--hard}, words[executable+2] is \"list\", which fails the != \"exec\" check at projection.go:50, so the result is ExecProjection{}.\n- C2 (projection_test.go:43-44) matches spec.md:91 \"It performs no recursion\". ShellChildren reads only the outer argv, so it yields one child at ArgumentIndex 2: [sh -c \"git push\"].\n- Refrozen record at ca726066 adds only reviews/shared-guard-projection.md over 93e90d3: chunk GP1 base f86203658, tip 93e90d3ee, acceptance rows unchanged.\n- The r1-t1-* entries are the five GP1 plan verification IDs with the r1 prefix, all with performer claude:sgp-t1-repair1-20261021, matching the plan's repair assignment, and source_digest equal to the chunk source_digest.\n- Test counts are consistent with the delta: shellcommand 42 to 43 tests; gitguard stays at 237.\n\nAdvice (no ID):\n- The r1 probes replay the original mutants; the red for the new C1 and C2 cases is the repair's own probe set and the Coverage axis evidence. Not binding.\n\nCommands (all via bench worktree exec shared-guard-projection-build --): bench preflight evidence sha256:5590e889... --check-current returned current=true, dirty_record=false; git diff a4053ed70 93e90d3ee; git diff --stat 93e90d3ee ca726066b; git log --oneline f86203658..ca726066b; git show of projection.go, projection_test.go, shellcommand.go, spec.md, and the review record; rg over spec.md.\n"
          },
          "axis": "Spec",
          "base": "f86203658f5190beac7a34ce17e18473d20fe426",
          "tip": "93e90d3ee55ea6adb6d713389077326e688fec12",
          "finding_ids": [],
          "supersedes": [
            "gp1-r1-spec"
          ]
        },
        {
          "id": "gp1-r2-coverage",
          "performer": "claude:sgp-review-coverage-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "09cf33d9e84ed49744b212f5792a0c2ea4a4bfbd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-coverage-r2",
            "digest": "sha256:e982f871a23aba1ef701488f44759e9fdedbeea444b58dca9f3ce13ef14eb310",
            "excerpt": "COVERAGE axis, GP1 round 1 confirmation. Both initial findings are closed, and I found no new finding. I used no bench probe, because the saved answer establishes the bite.\n\nC1: closed.\n- The repair adds {\"bench\",\"worktree\",\"list\",\"exec\",\"--\",\"git\",\"reset\",\"--hard\"} to the no-projection loop of TestWorktreeExecContext (internal/shellcommand/projection_test.go).\n- Repair mutant 1 in sgp-t1-r1-mutants.json swaps words[executable+2] != \"exec\" for !slices.Contains(words[executable+2:], \"exec\"), the same mutation as my probe 1.\n- The saved answer reports verdict bit, baseline passed, restored yes. The only failing test is TestWorktreeExecContext: projection_test.go:130: WorktreeExec([\"bench\" \"worktree\" \"list\" \"exec\" \"--\" \"git\" \"reset\" \"--hard\"]) = {... Child:[git reset --hard] DifferentCheckout:true}, want no projection.\n\nC2: closed.\n- The repair adds the case \"nested wrapper stays one child\": bash -c \"sh -c 'git push'\" expects the single child {2, [sh -c git push]}.\n- Repair mutant 2 appends ShellChildren(..., 0) of each child's own commands, the same one-level recursion as my probe 2.\n- The saved answer reports verdict bit, baseline passed, restored yes. The only failing test is TestShellChildren/nested_wrapper_stays_one_child: projection_test.go:49: ... = [{2 [sh -c git push]} {2 [git push]}], want [{2 [sh -c git push]}].\n\nRepair delta (git diff a4053ed70 93e90d3ee): two test files; the gitguard change is a comment reflow only. projection.go is unchanged, so the mutants apply to the code tip.\n\nNew findings: none.\n\nOptional advice (no ID): the GP10 redirected-push and GP74 first-reset-child cases pass on old and new code; bench worktree exec -- git push with no target has no test.\n\nCommands:\n1. bench worktree exec shared-guard-projection-build -- bench preflight evidence sha256:5590e889... --check-current: head ca726066, current true, dirty_record false.\n2. bench worktree exec shared-guard-projection-build -- git diff a4053ed70 93e90d3ee: read.\n3. Read sgp-t1-r1-cases.out, sgp-t1-r1-mutants.json, and gp1-coverage-r1.txt in the scratchpad.\nProbes used: 0 of 1.\n"
          },
          "axis": "Coverage",
          "base": "f86203658f5190beac7a34ce17e18473d20fe426",
          "tip": "93e90d3ee55ea6adb6d713389077326e688fec12",
          "finding_ids": [],
          "supersedes": [
            "gp1-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "GP2",
      "base": "6183e1e43df713b5dbbe2c58dc86df7106b78507",
      "tip": "743b19190a6ab51a6c7a1d65f4071659223ff53a",
      "plan_digest": "sha256:26860796c78f035de7c2a1e29d67091b9d1200c6a543fe521f3f02df7be53ea7",
      "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
      "acceptance_rows": [
        "GP14",
        "GP15",
        "GP16",
        "GP17",
        "GP18",
        "GP19",
        "GP20",
        "GP21",
        "GP43",
        "GP47",
        "GP49",
        "GP51",
        "GP54",
        "GP57",
        "GP60",
        "GP22",
        "GP23",
        "GP24",
        "GP25",
        "GP26",
        "GP27",
        "GP28",
        "GP29",
        "GP48",
        "GP50",
        "GP52",
        "GP53",
        "GP55",
        "GP56",
        "GP58",
        "GP59",
        "GP61",
        "GP04",
        "GP11",
        "GP12",
        "GP13",
        "GP62",
        "GP63",
        "GP64",
        "GP65",
        "GP66",
        "GP68",
        "GP69",
        "GP70",
        "GP73"
      ],
      "verification": [
        {
          "id": "a4-t4-bench-policy",
          "performer": "claude:sgp-t4-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-author-20261021",
            "digest": "sha256:c50559ae4406329f5073644cded4766043f047a347ee1a4b8cfd57f483049e4f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,pass,3,28\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-bench-policy",
          "command": "bench test --package ./internal/benchguard",
          "exit_code": 0
        },
        {
          "id": "a4-t4-actual-core-entries",
          "performer": "claude:sgp-t4-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-author-20261021",
            "digest": "sha256:03ff74cfe3d60cb87e4c53097977d561b91e7ecd1b64c8d4dd44b5b0dffce32a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14754,397\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-actual-core-entries",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "a4-t4-existing-census",
          "performer": "claude:sgp-t4-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-author-20261021",
            "digest": "sha256:0aa10bed9bff908c5b324f87708b29b357b2b8c5e55f5301a3ae92aa376d297d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/census,pass,81,67\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-existing-census",
          "command": "bench test --package ./internal/census",
          "exit_code": 0
        },
        {
          "id": "a4-t4-conformance-closure",
          "performer": "claude:sgp-t4-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-author-20261021",
            "digest": "sha256:ac3e07bc1eedd8c730af7ac9de88b6d3d93a8f2e8850a2cc6372ca4612cbc9de",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,34450,1127\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "a4-t4-root-conformance",
          "performer": "claude:sgp-t4-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-author-20261021",
            "digest": "sha256:b9fa65f4101f8331d7176ab4a220884241099c01a29e8fd69771ca391b70fa88",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12758,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "a4-t4-reviewer-recursion-omission",
          "performer": "claude:sgp-t4-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-author-20261021",
            "digest": "sha256:6b57909d6023f7cae652e6e9f0b74a5ff097d6cc7b51fa7146e6521d3ff68ac5",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestReviewerVerbKeepsRecursiveReach,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/benchguard,^TestReviewerVerbKeepsRecursiveReach$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,443\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/benchguard,TestReviewerVerbKeepsRecursiveReach,\"projection_contract_test.go:89: ReviewerVerb(\\\"bash -c \\\\\\\\\\\"sh -c 'bench commitment override --assignment x --reason y'\\\\\\\\\\\"\\\") = false, want true\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-reviewer-recursion-omission",
          "command": "bench test --package ./internal/benchguard --run '^TestReviewerVerbKeepsRecursiveReach$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace only the reviewer walk's recursion into a wrapper child stream with a one-level verb check; the nested-wrapper case must red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t4-author-20261021",
              "digest": "sha256:6b57909d6023f7cae652e6e9f0b74a5ff097d6cc7b51fa7146e6521d3ff68ac5",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestReviewerVerbKeepsRecursiveReach,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/benchguard,^TestReviewerVerbKeepsRecursiveReach$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,443\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/benchguard,TestReviewerVerbKeepsRecursiveReach,\"projection_contract_test.go:89: ReviewerVerb(\\\"bash -c \\\\\\\\\\\"sh -c 'bench commitment override --assignment x --reason y'\\\\\\\\\\\"\\\") = false, want true\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "a2-t2-owner-values",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:46fe003d563183afbaa192d908d9f74879898f8dede42631dface30dad7193cc",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,pass,3,50\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-owner-values",
          "command": "bench test --package ./internal/shellcommand",
          "exit_code": 0
        },
        {
          "id": "a2-t2-git-policy",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:5a056919785168a84836a91bb6a577f57dec8d3e587d1273884f7246594ca512",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,pass,15168,250\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-git-policy",
          "command": "bench test --package ./internal/gitguard",
          "exit_code": 0
        },
        {
          "id": "a2-t2-actual-core-entries",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:8bb70b500a8ab8547e1953956ff14462061843716d951b988ad0af418132a454",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14797,397\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-actual-core-entries",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "a2-t2-conformance-closure",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:31c0d4b2a77086e55cd126df44078107dbea379eb803df3890065b7ac6bad78a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,34427,1127\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "a2-t2-root-conformance",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:9774a81291b5039e5970065a471356dc486b981acc3ab9400a0166298c18d1f2",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12586,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "a2-t2-canary-bites",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:0064ceda55eee0b6aa37d7fb13856d15fcd3249cb7233221b5c7edaf48cb7076",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3303,288\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-canary-bites",
          "command": "bench test --package ./internal/conformance --run '^TestEveryRetainedFixtureBitesThroughRegisteredOwner$'",
          "exit_code": 0
        },
        {
          "id": "a2-t2-git-decoder-form-omission",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:7af52e2a559e284c1af06eab3dc33d274971c913b7fe2b23075bd98f3dc29343",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestGuardGitEnvelopePolicy/duplicate_command_null,16,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./cmd/bench,^TestGuardGitEnvelopePolicy$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1214\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,5,16\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/duplicate_command_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\",\\\\\\\\\\\"command\\\\\\\\\\\":null}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/duplicate_input_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\",\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"},\\\\\\\\\\\"tool_input\\\\\\\\\\\":null,\\\\\\\\\\\"cwd\\\\\\\\\\\":\\\\\\\\\\\"/repo\\\\\\\\\\\"}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/nested_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"COMMAND\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/root_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"TOOL_INPUT\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-git-decoder-form-omission",
          "command": "bench test --package ./cmd/bench --run '^TestGuardGitEnvelopePolicy$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace only the Git caller CommandStruct selection with InputMap; ordered folded/null/early-error cases must red through the actual guardGit entry.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
              "digest": "sha256:7af52e2a559e284c1af06eab3dc33d274971c913b7fe2b23075bd98f3dc29343",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestGuardGitEnvelopePolicy/duplicate_command_null,16,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./cmd/bench,^TestGuardGitEnvelopePolicy$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1214\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,5,16\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/duplicate_command_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\",\\\\\\\\\\\"command\\\\\\\\\\\":null}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/duplicate_input_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\",\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"},\\\\\\\\\\\"tool_input\\\\\\\\\\\":null,\\\\\\\\\\\"cwd\\\\\\\\\\\":\\\\\\\\\\\"/repo\\\\\\\\\\\"}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/nested_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"COMMAND\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/root_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"TOOL_INPUT\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "a3-t3-bench-policy",
          "performer": "claude:sgp-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-author-20261021",
            "digest": "sha256:c50559ae4406329f5073644cded4766043f047a347ee1a4b8cfd57f483049e4f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,pass,3,28\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-bench-policy",
          "command": "bench test --package ./internal/benchguard",
          "exit_code": 0
        },
        {
          "id": "a3-t3-write-policy",
          "performer": "claude:sgp-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-author-20261021",
            "digest": "sha256:57683df3fee739f11d15b03e4e6a378a152f0652546b9ea6d577ae89ef2f6c62",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/writeguard,pass,2,28\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-write-policy",
          "command": "bench test --package ./internal/writeguard",
          "exit_code": 0
        },
        {
          "id": "a3-t3-actual-core-entries",
          "performer": "claude:sgp-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-author-20261021",
            "digest": "sha256:59bd8398d193fa0cbaf074df46385ef065d06d6a65033471356c66ec16ed68d7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14443,397\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-actual-core-entries",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "a3-t3-conformance-closure",
          "performer": "claude:sgp-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-author-20261021",
            "digest": "sha256:bd371764337ccaa65bd509ef2205e808f73f99aa19ba37acb7f54b50c9103343",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,33939,1127\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "a3-t3-root-conformance",
          "performer": "claude:sgp-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-author-20261021",
            "digest": "sha256:24888b4ec36d49500e51fa4a5d921b6a1f44cb3d8a83ba6c6948403efdc856bd",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12766,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "a3-t3-bench-decoder-form-omission",
          "performer": "claude:sgp-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-author-20261021",
            "digest": "sha256:dfe10eca47e8431b3abcf9c3a969a82ccadef43c0fb964fd9f54e0f441ec4815",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestBenchEnvelopeValue/NUL_command,10,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/benchguard,^TestBenchEnvelope(RequiredField|Value)$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1283\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,fail,2,10\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/NUL_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"git\\\\\\\\u0000status\\\"}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field has control byte\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/duplicate_command_null,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"git reset --hard\\\",\\\"command\\\":null}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/empty_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"\\\"}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/null_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":null}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-bench-decoder-form-omission",
          "command": "bench test --package ./internal/benchguard --run '^TestBenchEnvelope(RequiredField|Value)$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace only the Bench caller InputMap selection with CommandStruct; exact-map and map-null cases must red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t3-author-20261021",
              "digest": "sha256:dfe10eca47e8431b3abcf9c3a969a82ccadef43c0fb964fd9f54e0f441ec4815",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestBenchEnvelopeValue/NUL_command,10,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/benchguard,^TestBenchEnvelope(RequiredField|Value)$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1283\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,fail,2,10\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/NUL_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"git\\\\\\\\u0000status\\\"}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field has control byte\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/duplicate_command_null,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"git reset --hard\\\",\\\"command\\\":null}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/empty_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"\\\"}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/null_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":null}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "a3-t3-write-decoder-form-omission",
          "performer": "claude:sgp-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-author-20261021",
            "digest": "sha256:3e1a7e6a59d6e5de5baba8766f60daa62dc8bdcf52c10979bab0021bf6aa12d5",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGuardWriteEnvelopePolicy/duplicate_cwd_null,13,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./cmd/bench,^TestGuardWriteEnvelopePolicy$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,836\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,12,13\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestGuardWriteEnvelopePolicy/duplicate_cwd_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"cwd\\\\\\\\\\\":\\\\\\\\\\\"/dev/shm/2CLYPW/t/TestGuardWriteEnvelopePolicy4043833738/001\\\\\\\\\\\",\\\\\\\\\\\"cwd\\\\\\\\\\\":null,\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardWriteEnvelopePolicy/write_root_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"CWD\\\\\\\\\\\":\\\\\\\\\\\"/dev/shm/2CLYPW/t/TestGuardWriteEnvelopePolicy4043833738/001\\\\\\\\\\\",\\\\\\\\\\\"TOOL_INPUT\\\\\\\\\\\":{\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-write-decoder-form-omission",
          "command": "bench test --package ./cmd/bench --run '^TestGuardWriteEnvelopePolicy$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace only the Write caller CwdInputMap selection with InputMap; actual tracked-relative-primary cwd cases must red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t3-author-20261021",
              "digest": "sha256:3e1a7e6a59d6e5de5baba8766f60daa62dc8bdcf52c10979bab0021bf6aa12d5",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGuardWriteEnvelopePolicy/duplicate_cwd_null,13,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./cmd/bench,^TestGuardWriteEnvelopePolicy$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,836\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,12,13\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestGuardWriteEnvelopePolicy/duplicate_cwd_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"cwd\\\\\\\\\\\":\\\\\\\\\\\"/dev/shm/2CLYPW/t/TestGuardWriteEnvelopePolicy4043833738/001\\\\\\\\\\\",\\\\\\\\\\\"cwd\\\\\\\\\\\":null,\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardWriteEnvelopePolicy/write_root_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"CWD\\\\\\\\\\\":\\\\\\\\\\\"/dev/shm/2CLYPW/t/TestGuardWriteEnvelopePolicy4043833738/001\\\\\\\\\\\",\\\\\\\\\\\"TOOL_INPUT\\\\\\\\\\\":{\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r4-t4-bench-policy",
          "performer": "claude:sgp-t4-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-repair1-20261021",
            "digest": "sha256:6f8b76a7ebe3208f3d41f7dc1e4096f29fe1fe039be9540e050f2612db51c7f0",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,pass,2,29\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-bench-policy",
          "command": "bench test --package ./internal/benchguard",
          "exit_code": 0
        },
        {
          "id": "r4-t4-actual-core-entries",
          "performer": "claude:sgp-t4-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-repair1-20261021",
            "digest": "sha256:c72d7f3ca38b9c29f156ad9b4b6b84e550d86232350195509d49ed934c9b0cfa",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14563,398\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-actual-core-entries",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "r4-t4-existing-census",
          "performer": "claude:sgp-t4-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-repair1-20261021",
            "digest": "sha256:81391e947200cbcc3c4031158f320141971f237efac5bed44f4f178add3a1e5e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/census,pass,131,67\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-existing-census",
          "command": "bench test --package ./internal/census",
          "exit_code": 0
        },
        {
          "id": "r4-t4-conformance-closure",
          "performer": "claude:sgp-t4-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-repair1-20261021",
            "digest": "sha256:15498087a58a9040077b8e581fa57ff69146252380d2556a7a824e5c6f50df51",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,35047,1127\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "r4-t4-root-conformance",
          "performer": "claude:sgp-t4-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-repair1-20261021",
            "digest": "sha256:6decfb18a91ed775c62887c0094b9f38733122ce18f4e9a3db372dc1fff3dcbf",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12843,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "r4-t4-reviewer-recursion-omission",
          "performer": "claude:sgp-t4-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t4-repair1-20261021",
            "digest": "sha256:af1d4b2d6c7cdf77a7673315f822c9e33b2de1c0d929012dc7f2f6cf0bbf24b5",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestReviewerVerbKeepsRecursiveReach,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/benchguard,^TestReviewerVerbKeepsRecursiveReach$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,443\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/benchguard,TestReviewerVerbKeepsRecursiveReach,\"projection_contract_test.go:95: ReviewerVerb(\\\"bash -c \\\\\\\\\\\"sh -c 'bench commitment override --assignment x --reason y'\\\\\\\\\\\"\\\") = false, want true\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-reviewer-recursion-omission",
          "command": "bench test --package ./internal/benchguard --run '^TestReviewerVerbKeepsRecursiveReach$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace only the reviewer walk's recursion into a wrapper child stream with a one-level verb check; the nested-wrapper case must red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t4-repair1-20261021",
              "digest": "sha256:af1d4b2d6c7cdf77a7673315f822c9e33b2de1c0d929012dc7f2f6cf0bbf24b5",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestReviewerVerbKeepsRecursiveReach,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/benchguard,^TestReviewerVerbKeepsRecursiveReach$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,443\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/benchguard,TestReviewerVerbKeepsRecursiveReach,\"projection_contract_test.go:95: ReviewerVerb(\\\"bash -c \\\\\\\\\\\"sh -c 'bench commitment override --assignment x --reason y'\\\\\\\\\\\"\\\") = false, want true\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r3-t3-bench-policy",
          "performer": "claude:sgp-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-repair1-20261021",
            "digest": "sha256:e7b33ca54138d54370163ca516484aa2c2dedc3f2ee67a6f96b8072099bac3c1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,pass,3,29\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-bench-policy",
          "command": "bench test --package ./internal/benchguard",
          "exit_code": 0
        },
        {
          "id": "r3-t3-write-policy",
          "performer": "claude:sgp-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-repair1-20261021",
            "digest": "sha256:57683df3fee739f11d15b03e4e6a378a152f0652546b9ea6d577ae89ef2f6c62",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/writeguard,pass,2,28\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-write-policy",
          "command": "bench test --package ./internal/writeguard",
          "exit_code": 0
        },
        {
          "id": "r3-t3-actual-core-entries",
          "performer": "claude:sgp-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-repair1-20261021",
            "digest": "sha256:5b7f897e73575974992d16746b95ddba6a4caad17e11294c3f774c58c3f25d23",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14524,398\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-actual-core-entries",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "r3-t3-conformance-closure",
          "performer": "claude:sgp-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-repair1-20261021",
            "digest": "sha256:a0c06abc05aa09f9bfa596be3d61a815000462c3c72a9189584b9cc2a4fae040",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,33679,1127\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "r3-t3-root-conformance",
          "performer": "claude:sgp-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-repair1-20261021",
            "digest": "sha256:ab17319514f1f0cc18cc234f1d00a1becfe9d3557a63a6968619e5852382edf8",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12800,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "r3-t3-bench-decoder-form-omission",
          "performer": "claude:sgp-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-repair1-20261021",
            "digest": "sha256:dfe10eca47e8431b3abcf9c3a969a82ccadef43c0fb964fd9f54e0f441ec4815",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestBenchEnvelopeValue/NUL_command,10,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/benchguard,^TestBenchEnvelope(RequiredField|Value)$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1283\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,fail,2,10\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/NUL_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"git\\\\\\\\u0000status\\\"}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field has control byte\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/duplicate_command_null,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"git reset --hard\\\",\\\"command\\\":null}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/empty_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"\\\"}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/null_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":null}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-bench-decoder-form-omission",
          "command": "bench test --package ./internal/benchguard --run '^TestBenchEnvelope(RequiredField|Value)$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace only the Bench caller InputMap selection with CommandStruct; exact-map and map-null cases must red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t3-repair1-20261021",
              "digest": "sha256:dfe10eca47e8431b3abcf9c3a969a82ccadef43c0fb964fd9f54e0f441ec4815",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestBenchEnvelopeValue/NUL_command,10,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/benchguard,^TestBenchEnvelope(RequiredField|Value)$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1283\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,fail,2,10\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/NUL_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"git\\\\\\\\u0000status\\\"}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field has control byte\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/duplicate_command_null,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"git reset --hard\\\",\\\"command\\\":null}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/empty_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":\\\"\\\"}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestBenchEnvelopeValue/null_command,\"projection_contract_test.go:21: CommandFromEnvelope({\\\"tool_input\\\":{\\\"command\\\":null}}) = (\\\"\\\", \\\"command field missing\\\"), want (\\\"\\\", \\\"command field empty\\\")\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r3-t3-write-decoder-form-omission",
          "performer": "claude:sgp-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t3-repair1-20261021",
            "digest": "sha256:22b9a9ac7e125f05d92c753ecf0e53b21dd960a568997359bae65346393bd120",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGuardWriteEnvelopePolicy/duplicate_cwd_null,13,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./cmd/bench,^TestGuardWriteEnvelopePolicy$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,836\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,11,13\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestGuardWriteEnvelopePolicy/duplicate_cwd_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"cwd\\\\\\\\\\\":\\\\\\\\\\\"/dev/shm/MNWKB4/t/TestGuardWriteEnvelopePolicy2807368414/001\\\\\\\\\\\",\\\\\\\\\\\"cwd\\\\\\\\\\\":null,\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardWriteEnvelopePolicy/write_root_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"CWD\\\\\\\\\\\":\\\\\\\\\\\"/dev/shm/MNWKB4/t/TestGuardWriteEnvelopePolicy2807368414/001\\\\\\\\\\\",\\\\\\\\\\\"TOOL_INPUT\\\\\\\\\\\":{\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-write-decoder-form-omission",
          "command": "bench test --package ./cmd/bench --run '^TestGuardWriteEnvelopePolicy$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace only the Write caller CwdInputMap selection with InputMap; actual tracked-relative-primary cwd cases must red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t3-repair1-20261021",
              "digest": "sha256:22b9a9ac7e125f05d92c753ecf0e53b21dd960a568997359bae65346393bd120",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGuardWriteEnvelopePolicy/duplicate_cwd_null,13,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./cmd/bench,^TestGuardWriteEnvelopePolicy$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,836\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,11,13\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestGuardWriteEnvelopePolicy/duplicate_cwd_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"cwd\\\\\\\\\\\":\\\\\\\\\\\"/dev/shm/MNWKB4/t/TestGuardWriteEnvelopePolicy2807368414/001\\\\\\\\\\\",\\\\\\\\\\\"cwd\\\\\\\\\\\":null,\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardWriteEnvelopePolicy/write_root_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"CWD\\\\\\\\\\\":\\\\\\\\\\\"/dev/shm/MNWKB4/t/TestGuardWriteEnvelopePolicy2807368414/001\\\\\\\\\\\",\\\\\\\\\\\"TOOL_INPUT\\\\\\\\\\\":{\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "b2-t2-owner-values",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:46fe003d563183afbaa192d908d9f74879898f8dede42631dface30dad7193cc",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,pass,3,50\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-owner-values",
          "command": "bench test --package ./internal/shellcommand",
          "exit_code": 0
        },
        {
          "id": "b2-t2-git-policy",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:5a056919785168a84836a91bb6a577f57dec8d3e587d1273884f7246594ca512",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,pass,15168,250\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-git-policy",
          "command": "bench test --package ./internal/gitguard",
          "exit_code": 0
        },
        {
          "id": "b2-t2-actual-core-entries",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:400bdea2f428256a670ef153b2098409947e646c7d71fe7b5dc2f695debbb0bc",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14438,398\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-actual-core-entries",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "b2-t2-conformance-closure",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:079a843e0b9975dd4f9dab9355b14bfd06e01f1c075c5f3b8baf2a3e8b6ebbc9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,34129,1127\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "b2-t2-root-conformance",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:99474a3f1c32ecb4c5be3bd5a1047c4b1ee24ca9b69c74705dcfa183c2f060e6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12659,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "b2-t2-canary-bites",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:72729a06f09db984178ad66cd7563152a2b91ed028b76fb65ad6867032fe4276",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3274,288\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-canary-bites",
          "command": "bench test --package ./internal/conformance --run '^TestEveryRetainedFixtureBitesThroughRegisteredOwner$'",
          "exit_code": 0
        },
        {
          "id": "b2-t2-git-decoder-form-omission",
          "performer": "claude:sgp-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
            "digest": "sha256:5d343b21e2b84e4cc00dcec2f516c46d79960b94312e77314f7cd959d01578aa",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,5,TestGuardGitEnvelopePolicy/duplicate_command_null,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./cmd/bench,^TestGuardGitEnvelopePolicy$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1453\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,5,17\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/duplicate_command_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\",\\\\\\\\\\\"command\\\\\\\\\\\":null}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/duplicate_input_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\",\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"},\\\\\\\\\\\"tool_input\\\\\\\\\\\":null,\\\\\\\\\\\"cwd\\\\\\\\\\\":\\\\\\\\\\\"/repo\\\\\\\\\\\"}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/nested_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"COMMAND\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/root_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"TOOL_INPUT\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/unrelated_cwd_type,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"cwd\\\\\\\\\\\":7,\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-git-decoder-form-omission",
          "command": "bench test --package ./cmd/bench --run '^TestGuardGitEnvelopePolicy$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace only the Git caller CommandStruct selection with InputMap; ordered folded/null/early-error cases must red through the actual guardGit entry.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:sgp-t2-author-20261021",
              "digest": "sha256:5d343b21e2b84e4cc00dcec2f516c46d79960b94312e77314f7cd959d01578aa",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,5,TestGuardGitEnvelopePolicy/duplicate_command_null,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./cmd/bench,^TestGuardGitEnvelopePolicy$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1453\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,5,17\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/duplicate_command_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\",\\\\\\\\\\\"command\\\\\\\\\\\":null}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/duplicate_input_null,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\",\\\\\\\\\\\"file_path\\\\\\\\\\\":\\\\\\\\\\\"README.md\\\\\\\\\\\"},\\\\\\\\\\\"tool_input\\\\\\\\\\\":null,\\\\\\\\\\\"cwd\\\\\\\\\\\":\\\\\\\\\\\"/repo\\\\\\\\\\\"}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/nested_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"COMMAND\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/root_key_casing,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"TOOL_INPUT\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\n  github.com/gibbonmi/bench/cmd/bench,TestGuardGitEnvelopePolicy/unrelated_cwd_type,\"guard_projection_test.go:37: entry(\\\"{\\\\\\\\\\\"cwd\\\\\\\\\\\":7,\\\\\\\\\\\"tool_input\\\\\\\\\\\":{\\\\\\\\\\\"command\\\\\\\\\\\":\\\\\\\\\\\"git reset --hard\\\\\\\\\\\"}}\\\") exit = 0, want 2\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gp2-r1-standards",
          "performer": "claude:sgp-review-standards-gp2-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-standards-gp2-r1",
            "digest": "sha256:e17fc41882f924b05b7d3735262a9987a927e0e10c8dfabc39b700c140185ff1",
            "excerpt": "STANDARDS axis, chunk GP2 (base 6183e1e43, code tip 813efeb55, evidence tip 9f83ef35). The axis has 1 finding, low confidence; the worst issue is the unrecorded \"respelled error text turns this row red\" claim in two new tests.\n\nFindings\n- F1 | auto-fix | ticket 3 | confidence 5\n  - Where: internal/writeguard/envelope_contract_test.go:29-31 (the \"duplicate input null\" row compares err.Error() to \"file_path field missing\"). The claim is repeated at internal/benchguard/projection_contract_test.go:12-13.\n  - Binding requirement: AGENTS.md says an independently authored expectation counts as duplicated implementation knowledge unless the independence is necessary for a named mutation to turn the gate red and that red is recorded and demonstrated.\n  - Evidence: the review record at 9f83ef35 holds only two ticket 3 mutation probes (the Bench InputMap to CommandStruct swap and the Write CwdInputMap to InputMap swap). No respelled-text mutant appears.\n  - Benchguard: the recorded swap does redden TestBenchEnvelopeValue through the expected-text mismatch, so the benchguard rows arguably meet the exception.\n  - Writeguard: the recorded Write probe reddens TestGuardWriteEnvelopePolicy, not this row. The text comparison is not shown to be necessary.\n  - Repair: record the respell red, or drop the text compare in the writeguard row and reword the comment.\n  - The six error-text rows in benchguard (projection_contract_test.go:30-53) pin texts that acceptance rows GP22-GP24 say to keep, so no separate finding.\n\nAnswers to the named questions\n1. Duplicated E0 expectations in gitguard_test.go TestCommandFromEnvelope (lines 10-29) and main_test.go TestGuardGitBlockAllow (lines 524-545): not a defect this chunk must fix. Neither test is in the diff, both pass unchanged, no ticket's Writes: covers them, and main_test.go is 719 lines. Optional advice.\n2. Spec caller-sweep lines (spec.md:630-638): no prose correction is owed. Line 630 ties the sweep to f07508b2add49469d686614c70f575531e8ea445, the pre-migration subject. The retirement is graded by ticket 4's acceptance, which the delta satisfies.\n\nGraded and clean\n- One source per fact: gitguard, benchguard and writeguard have no json.Unmarshal on an envelope; benchguard has no isWrapper, isExecHead, wrapperFlag, \"--\", \"worktree\" or \"exec\" literals; reviewer.go reads children only through WorktreeExec and ShellChildren; GP69 holds (runsBenchUnderExec reads the child only through WorktreeExec(...).Child).\n- Behavior preservation in moved code: CommandFromEnvelope and PathFromEnvelope keep their error order and texts; the reviewer walk returns false for an exec head with no --; judge and spanInvokesBench search every matching -c child one level deep.\n- Shared harness: one checkGuardEnvelope in cmd/bench/guard_projection_test.go (lines 30-55).\n- Pool and reviewer-verb expectations come from benchguard.PoolReferenceMessage and ReviewerVerbMessage.\n- Stdlib reimplementation: none. Comments: no narrated history; longest sentences 22-23 words.\n- MUTATE.json re-anchor: old anchor matches the gitguard.go import block at 813efeb55; EXPECT unchanged.\n- Line limits: benchguard.go 363, reviewer.go 60, guard_projection_test.go 148.\n- Fence: every changed path is in its ticket's Writes:.\n\nAdvice (no finding ID, no disposition)\n- Run bench learning to record that TestCommandFromEnvelope and TestGuardGitBlockAllow repeat the E0 rows now owned by GP20 and GP21.\n- Envelope carries all three fields behind one type, so Command() on an InputMap envelope silently returns \"\". The spec fixes this shape.\n- Ticket 3 and 4 Writes: lines mark shared test files \"(new)\" although an earlier ticket creates them.\n\nCount: 1 finding (F1, ticket 3).\n\nCommands, all through bench worktree exec shared-guard-projection-build --: bench preflight evidence sha256:31f00ef8... --check-current returned current=true, dirty_record=false; git diff --stat, git diff, git show, git log; rg for retired names and error texts. No tests, no bench probe.\n"
          },
          "axis": "Standards",
          "base": "6183e1e43df713b5dbbe2c58dc86df7106b78507",
          "tip": "813efeb55dbbf65cdd3777178012391c4e241039",
          "finding_ids": [
            "F1"
          ],
          "supersedes": []
        },
        {
          "id": "gp2-r1-spec",
          "performer": "claude:sgp-review-spec-gp2-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-spec-gp2-r1",
            "digest": "sha256:2d3a8c2b4377367e0e534b5b2e912701969d238235a9133df650233cb0b18410",
            "excerpt": "SPEC axis, chunk GP2 (tickets 2-4), base 6183e1e43 to code tip 813efeb55.\n\nFindings: 0 (no finding survived refutation). Worst issue: none.\n\nAdvice (optional, no finding ID):\n- Commit 813efeb55 is titled \"migrate Bench projection consumers onto ShellChildren and WorktreeExec\"; tickets 2 and 3 use the \"FT366 shared-guard-projection tN:\" prefix. No rule binds the prefix.\n\nWhat the diff matches in the spec:\n- Envelope API: each form decodes through one json.Unmarshal into the exact anonymous target and tags the form table lists. Root keys fold in all three forms; the nested command key folds only in CommandStruct. InputString reads exact map keys; missing gives present=false; a type error gives present=true with the error; a null raw value gives \"\" with present=true and no error. The owner chooses no exit, resolves no path, and validates no field. The extra unknown-form error branch is harmless.\n- Caller policy: Git is string-only and returns \"\" on any decode error. Bench keeps its signature and error texts. Write decodes the outer cwd before file_path and path handling, so a bad cwd errors even for an absolute path (envelope_contract_test.go:13-18). Git and Bench targets contain no cwd field.\n- Decoder matrix, row by row: Git actual entry guard_projection_test.go:57-77 (GP47, 49, 51, 54, 57); Bench value reader projection_contract_test.go:27-55 (GP48, 50, 52, 55, 58); Write reader envelope_contract_test.go:27-33 (GP53) and :17 (GP61); Write actual entry guard_projection_test.go:98-121 (GP56, 59, GP29); owner merge rows envelope_test.go:59-79 (GP60, GP16); owner rows GP14-19 envelope_test.go:26-132; Git and Bench E0 guard_projection_test.go:57-94 (GP21, GP25); Git reader rows gitguard/projection_contract_test.go:28-55 (GP20).\n- Bench policy: judge searches every matching -c child one level and calls WorktreeExec(...).Head only after isBench; spanInvokesBench uses ShellChildren; ReviewerVerb recurses at every level through ShellChildren and WorktreeExec with isBench first; an exec head with no -- gives a nil Child; cmd/bench/guards.go is not in the diff, so the entry order is unchanged (GP13, 65, 66 at guard_projection_test.go:132-148).\n- runsBenchUnderExec keeps its Bench identity test and reads the exec child only through shellcommand.WorktreeExec; RR12 untouched.\n- Preserved closure: MUTATE.json re-anchored on the shellcommand import block with no change to the diagnostic; type Checker struct byte-identical; hook shims not in the diff.\n- Retired copies: git grep of isWrapper, isExecHead, wrapperFlag, and encoding/json in guard production files at 813efeb55 returns no hit.\n- Structure and exclusions: benchguard.go 363 lines; cmd/bench/main.go, internal/lines/lines.go, cmd/bench/guards.go, projection.go, scan.go and .bench/hooks unchanged.\n- Spec cell rewrites: each rewritten cell sits in the commit of the ticket that covers the row, and each cited test lands in that commit.\n- Diff scope: every changed file is in the GP2 Writes: scope; no edit outside acceptance reach.\n\nCommands run (all through bench worktree exec shared-guard-projection-build --): bench preflight evidence sha256:31f00ef8...e2c6 --check-current returned current=true, dirty_record=false; git diff --stat and git diff 6183e1e43 813efeb55; git log; git show -U0 --word-diff=porcelain on the three code commits; git show of the spec, tickets, and changed files; git grep for test names, retired symbols, and imports. No tests.\n"
          },
          "axis": "Spec",
          "base": "6183e1e43df713b5dbbe2c58dc86df7106b78507",
          "tip": "813efeb55dbbf65cdd3777178012391c4e241039",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "gp2-r1-coverage",
          "performer": "claude:sgp-review-coverage-gp2-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6818978d77749d4c28f43936b62fe37236394b19",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-coverage-gp2-r1",
            "digest": "sha256:96993c8ee1f4663af47a377bc0ff9bfc0dd04e156a659f0a588c662f8c130989",
            "excerpt": "COVERAGE, chunk GP2 (base 6183e1e43, code tip 813efeb55, evidence tip 9f83ef35). Findings: 3, all test-only auto-fix. Worst: C1.\n\nFindings (disposition | cite | binding requirement | ticket | confidence)\n- C1 | auto-fix | internal/benchguard/benchguard.go:138 (invokes(child.Stream, resolver, false) in judge) and :305 (same call in spanInvokesBench, reached by Classify through judgeExec at :175); tests internal/benchguard/projection_contract_test.go:66 and benchguard_test.go:72-84 | spec \"Git and Bench Classify limit traversal to one level across shell and worktree child forms\", edge W4 \"neither recurses twice\", and GP04's why-clause \"A universal first-child or recursive policy changes the Bench verdict\" | ticket 4 | 7\n  - Probe 1 swapped :138 false for true. It read SILENT: benchguard 28 tests and cmd/bench 397 tests passed, restored.\n  - Derived, not run: bash -c \"sh -c 'bench status'\" && echo hi is unblocked on the tip and blocked under the :138 mutant. bench worktree exec L -- true; bash -c \"sh -c 'bench status'\" is allowed on the tip and refused under the coordinator's :305 mutant.\n  - The test gap predates this diff, but ticket 4 rewrote both lines, and the :305 flag is Classify-reachable.\n  - Fix: add a two-level negative to TestGuardWrapperSelection (an outer-syntax case and an exec-later-segment case) and to TestInvokesBenchWalksOneWrapperLevel.\n- C2 | auto-fix | internal/benchguard/benchguard.go:26 (DecodeEnvelope(data, shellcommand.InputMap)) and internal/benchguard/projection_contract_test.go:27-55 | spec \"Neither Git nor Bench begins to validate an unrelated cwd field\"; GP18 \"Eager sibling validation changes Git and Bench posture\" | ticket 3 | 6\n  - GP18 is witnessed only at the owner (TestEnvelopeIndependentFields). No caller test passes a wrong-typed cwd.\n  - Probe 2 swapped Bench's InputMap for CwdInputMap. It read SILENT (benchguard and cmd/bench passed, restored).\n  - Fix: add a {\"cwd\":7,\"tool_input\":{\"command\":\"git reset --hard\"}} row to TestBenchEnvelopeType (command readable, no error) and to TestGuardGitEnvelopePolicy (exit 2).\n- C3 | auto-fix | internal/benchguard/benchguard.go:127 (outer && shellcommand.WorktreeExec(...).Head); tests benchguard_test.go TestClassifySpanScopedFollowOns | spec GP12 \"outer exec follow-on and heredoc exceptions retain verdicts\"; judge doc \"the exec exception covers a direct head only, so it never applies inside that string\" | ticket 4 | 5\n  - Probe 3 dropped outer &&. It read SILENT (both packages passed, restored). No test puts a bench worktree exec follow-on inside a wrapper string.\n  - Derived, not run: bash -c 'bench worktree exec L -- cp a b; cp b a' is blocked on the tip; without outer it reaches judgeExec and is allowed. The gap predates the diff; ticket 4 rewrote the line.\n  - Fix: add that command as a refuse row.\n\nNamed gaps cleared\n- Decoder matrix: every row has a witness at its stated seam (Git entry rows through guardGit; Write entry rows through guardFileWrite against a real primary checkout with a tracked README.md; Bench and Write contract rows; GP60 in TestEnvelopeInputShape for all three forms).\n- GP17, GP19: covered for all three forms. GP18: owner level only (see C2).\n- GP04: first-child and last-child mutants bite cases 1 and 2; the depth half is C1.\n- GP63: nested exec bites through the plan probe set.\n- GP68: RR12 bites loss of exec-child detection; rg '\"--\"|\"worktree\"|\"exec\"' over refusal_route_guard_test.go, reviewer.go and benchguard.go found no match.\n- GP43: relies on the recorded t2-canary-bites verification.\n\nRows that pass on old and new code with no demonstrated red: GP11, GP12 (except C3), GP26, GP27 (the coordinator NUL probe bit), GP64, GP70, GP13, GP65, GP66, GP68.\n\nAdvice (no ID)\n1. The GP63 no-separator case is weak: namesVerb would return false for any exec head. It is the spec's named input, so keep it.\n2. runsBenchUnderExec keeps its Bench identity test with no negative witness; predates the diff.\n3. A first cwd:null has no entry-level case.\n\nPolicy: all three findings are checks beyond the approved acceptance rows. They fit GP2's single hardening cycle, and no blocking condition holds.\n\nCommands run, all through bench worktree exec shared-guard-projection-build --: bench preflight evidence sha256:31f00ef8... --check-current: current=true, dirty_record=false; git diff, git show, git log, git grep, rg (read-only); 3 probes (3 of 3 used), each silent and restored yes.\n"
          },
          "axis": "Coverage",
          "base": "6183e1e43df713b5dbbe2c58dc86df7106b78507",
          "tip": "813efeb55dbbf65cdd3777178012391c4e241039",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "gp2-r2-standards",
          "performer": "claude:sgp-review-standards-gp2-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-standards-gp2-r2",
            "digest": "sha256:0cf85d52156f10447e6ba98ef225f0224ac4bf4a3deb9d41e5bad1b47a36ae63",
            "excerpt": "STANDARDS confirming round 1, chunk GP2 (base 6183e1e43, code tip 743b19190, evidence tip 2d9f64794).\n\nEvidence check: bench preflight evidence sha256:8b3925a2... --check-current returned current=true, dirty_record=false, source_tip 2d9f64794.\n\nF1: CLOSED. The coordinator's respell probe answer (/tmp/claude-1000/-home-mgibs-workspace-bench/f5571912-814f-48c2-9e29-363a829710e2/scratchpad/sgp-gp2-f1-respell.out) shows verdict bit, restored yes, first_failure TestWriteEnvelopeCwdContract/duplicate_input_null. Its failure text is envelope_contract_test.go:31: PathFromEnvelope(...) = (\"\", file_path missing), want the file_path field missing error. That is the named mutation turning the independent expectation red, which is the AGENTS.md exception. The delta does not touch envelope_contract_test.go (29-31).\n\nNew finding (repair delta):\n- F2 | auto-fix | internal/benchguard/projection_contract_test.go:60-61 | confidence 8\n  - The added comment is one 33-word descriptive sentence: \"The nested rows read one level only, in the outer judge and in the exec later-segment scan, and the exec exception covers a direct head only, so an exec inside a wrapper refuses.\"\n  - Binding requirement: craft-comments SKILL.md:12 requires ASD-STE100 prose; .agents/skills/bench-craft-spec/references/ste-prose.md:17-18 caps a descriptive sentence at 25 words. Same rule as GP1 round 1 F1.\n  - A failed documented mandatory standard stays blocking even without an automated check (bounded-repair-policy.md:31), so it is not demoted to advice.\n  - Fix: split it into two sentences of 25 words or fewer each. Change no code. The correction is comment-only, so it is evidence-only per bounded-repair-policy.md:47-48.\n\nChecked, no defect:\n- New \"unrelated cwd type\" rows at cmd/bench/guard_projection_test.go:31 and internal/benchguard/projection_contract_test.go:44 test different layers, so there is no duplicated knowledge.\n- New Classify rows at projection_contract_test.go:70-72 are consistent with the table and add no comment.\n\nOptional advice: none.\n\nCommands run, all through bench worktree exec shared-guard-projection-build --: bench preflight evidence ... --check-current; git diff 813efeb55 743b19190 and --stat; git show 743b19190: for envelope_contract_test.go and projection_contract_test.go; git log -S; rg on the review file and ste-prose.md. No tests, no probes, no files written.\n"
          },
          "axis": "Standards",
          "base": "6183e1e43df713b5dbbe2c58dc86df7106b78507",
          "tip": "743b19190a6ab51a6c7a1d65f4071659223ff53a",
          "finding_ids": [
            "F2"
          ],
          "supersedes": [
            "gp2-r1-standards"
          ]
        },
        {
          "id": "gp2-r2-spec",
          "performer": "claude:sgp-review-spec-gp2-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-spec-gp2-r2",
            "digest": "sha256:890f449c84805ebdfc7ab16cea9a96bb643671e87821cd47f8ed1b869c1875e0",
            "excerpt": "SPEC axis, GP2 confirming round 1 (base 6183e1e43, code tip 743b19190, evidence tip 2d9f64794): PASS, no findings.\n\nEvidence current: bench preflight evidence sha256:8b3925a2... --check-current returned current=true, dirty_record=false, head 2d9f64794.\n\nRows traced against the tip's code:\n- C2, Bench (internal/benchguard/projection_contract_test.go:44): Bench decodes with InputMap, whose target has no cwd field (envelope.go:44-51). The cwd 7 is ignored, so the command is readable and err is \"\". A swap to CwdInputMap would error on cwd. Matches spec.md:193.\n- C2, Git (cmd/bench/guard_projection_test.go:76): Git decodes with CommandStruct, which also has no cwd field (envelope.go:35-43). The command survives, so guardGit exits 2.\n- C1 row A, bash -c \"sh -c 'bench status'\" && echo hi (projection_contract_test.go:38): the outer judge calls invokes(child.Stream, resolver, false) (benchguard.go:138), which returns false for the nested sh, so the verdict is allowed.\n- C1 row B, bench worktree exec L -- true; bash -c \"sh -c 'bench status'\" (projection_contract_test.go:39): judgeExec scans later segments with spanInvokesBench(..., true) (benchguard.go:175); the nested sh is not Bench at one level, so the verdict is allowed. Both rows match \"Git and Bench Classify limit traversal to one level\" (spec lines 105 and 473, W4).\n- C3, bash -c 'bench worktree exec L -- cp a b; cp b a' (projection_contract_test.go:40, ;, blocked): the outer flag blocks the exec exception inside the child. The base code has the same outer && isExecHead guard (6183e1e43:benchguard.go:134), so the row pins retained behavior (GP12).\n- The delta only adds rows and comment text, so every GP2 acceptance row is intact.\n\nRefrozen record (reviews/shared-guard-projection.md at 2d9f64794): GP2 base 6183e1e4, tip 743b1919, source_digest 070c2f92, plan_digest 26860796. All 20 plan verification ids appear at source 070c2f92 with prefixes b2-, r3-, r4-. The r3 and r4 performers match the plan's repair assignments. The repaired files fit the Writes: lines.\n\nAdvice (no ID):\n- No recorded plan probe shows red for the new Bench cwd row or the C1 and C3 rows; the repair probe sets hold that evidence.\n- The b2- rows are attributed to claude:sgp-t2-author-20261021, which has no repair assignment; that session is the ticket 2 author.\n- The new comment at projection_contract_test.go lines 58-60 has a sentence of about 33 words. This is a Standards matter.\n\nCommands run, all through bench worktree exec shared-guard-projection-build --: bench preflight evidence ... --check-current, git diff 813efeb55 743b19190, git log, git show, rg. No tests and no probes.\n"
          },
          "axis": "Spec",
          "base": "6183e1e43df713b5dbbe2c58dc86df7106b78507",
          "tip": "743b19190a6ab51a6c7a1d65f4071659223ff53a",
          "finding_ids": [],
          "supersedes": [
            "gp2-r1-spec"
          ]
        },
        {
          "id": "gp2-r2-coverage",
          "performer": "claude:sgp-review-coverage-gp2-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "070c2f9217dbd3ebda6c591e60ee3d2513cbbbe8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:sgp-review-coverage-gp2-r2",
            "digest": "sha256:c0fed787c2281241797c952e93313540fb2d5264b103a325237912a647036c70",
            "excerpt": "COVERAGE confirming round 1, GP2: C1, C2 and C3 are all closed. I found no new finding and ran no probe (0 of 1 used).\n\nCommands:\n- bench preflight evidence sha256:8b3925a2... --check-current returned current=true, dirty_record=false, source_tip=2d9f64794, base=6183e1e43.\n- Read the saved probe answers and mutant sets in /tmp/claude-1000/-home-mgibs-workspace-bench/f5571912-814f-48c2-9e29-363a829710e2/scratchpad/.\n- git diff 813efeb55 743b19190 for the two test files.\n- git diff --stat 1c01218bf 743b19190 shows only the two test files and one spec.md line changed; production code is identical to what the repair probes ran against.\n\nPer finding:\n- C1 (two-level wrapper depth): closed. The :138 and :305 false-to-true swaps in sgp-t4-r1-mutants.json match benchguard.go:138 and :305. sgp-t4-r1-cases.out shows mutant 1 (:138) bit on bash -c \"sh -c 'bench status'\" && echo hi (projection_contract_test.go:76, got (true,\"&&\"), want (false,\"\")). Mutant 2 (:305) bit on bench worktree exec L -- true; bash -c \"sh -c 'bench status'\". Both restored=yes.\n- C2 (Bench caller cwd independence): closed. The InputMap-to-CwdInputMap swap at benchguard.go:26 bit TestBenchEnvelopeType/unrelated_cwd_type (sgp-t3-r1-bench.out), restored=yes.\n- C3 (outer-only exec exception): closed. Mutant 3 (outer && dropped at benchguard.go:127) bit bash -c 'bench worktree exec L -- cp a b; cp b a', which got (false,\"\") where it wants (true,\";\"). restored=yes.\n\nGit-side C2 red route: adequate. gitguard.go:113 decodes with CommandStruct, which ignores cwd; envelope.go offers no form that decodes both a nested command and cwd, so a form swap cannot red the row. sgp-t3-r1-git.out shows an extra CwdInputMap decode mutant and a raw-bytes \"cwd\":7 mutant both failed TestGuardGitEnvelopePolicy/unrelated_cwd_type (exit 0, want 2) and restored.\n\nNew findings: none.\n\nOptional advice (no ID): the hardening-cycle Git row cannot be reddened by a pure form swap; this is advice only.\n"
          },
          "axis": "Coverage",
          "base": "6183e1e43df713b5dbbe2c58dc86df7106b78507",
          "tip": "743b19190a6ab51a6c7a1d65f4071659223ff53a",
          "finding_ids": [],
          "supersedes": [
            "gp2-r1-coverage"
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
      "from": "sha256:b81505d2f7e258a72e0ec8ccd6ff68968ed27183cf6dbfa9251140eced2210bf",
      "to": "sha256:ed266643ad636c39c00201d70cd717be17e818a9770d0a1ee83d3b5f690b8df7",
      "chunk_ids": {
        "GP1": [
          "GP1"
        ]
      }
    },
    {
      "from": "sha256:ed266643ad636c39c00201d70cd717be17e818a9770d0a1ee83d3b5f690b8df7",
      "to": "sha256:867fc6719dcec05d1ffed542b7afc332df124b4192ba607bc4f467f87f77230d",
      "chunk_ids": {
        "GP1": [
          "GP1"
        ]
      }
    },
    {
      "from": "sha256:867fc6719dcec05d1ffed542b7afc332df124b4192ba607bc4f467f87f77230d",
      "to": "sha256:26860796c78f035de7c2a1e29d67091b9d1200c6a543fe521f3f02df7be53ea7",
      "chunk_ids": {
        "GP1": [
          "GP1"
        ],
        "GP2": [
          "GP2"
        ]
      }
    }
  ]
}
```
