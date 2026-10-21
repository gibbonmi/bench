# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/shared-guard-projection/spec.md",
  "plan_digest": "sha256:ed266643ad636c39c00201d70cd717be17e818a9770d0a1ee83d3b5f690b8df7",
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
    }
  ]
}
```
