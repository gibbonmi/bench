# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/shared-guard-projection/spec.md",
  "plan_digest": "sha256:b81505d2f7e258a72e0ec8ccd6ff68968ed27183cf6dbfa9251140eced2210bf",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GP1",
      "base": "f86203658f5190beac7a34ce17e18473d20fe426",
      "tip": "a4053ed70c625935cd64c2d4fe5943f99348d88c",
      "plan_digest": "sha256:b81505d2f7e258a72e0ec8ccd6ff68968ed27183cf6dbfa9251140eced2210bf",
      "source_digest": "9176090363f3e82572626823207983480c37cc54",
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
