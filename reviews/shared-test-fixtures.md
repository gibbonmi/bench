# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/shared-test-fixtures/spec.md",
  "plan_digest": "sha256:dd48d59465866a1c9b1f35f8f9eecee4e98f8e465d8643764328edbea174193f",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GF-C1",
      "base": "4138ef1aa0c219c40396a9b799eac8769207b1f8",
      "tip": "63e498366d28674bb6a390f903391e8019bcea5f",
      "plan_digest": "sha256:dd48d59465866a1c9b1f35f8f9eecee4e98f8e465d8643764328edbea174193f",
      "source_digest": "46b8be48be138daa2cfc6e7648e239a442c6d0fd",
      "acceptance_rows": [
        "GF1",
        "GF2",
        "GF3",
        "GF4",
        "GF5",
        "GF6",
        "GF7",
        "GF8",
        "GF9",
        "GF11",
        "GF12",
        "GF13",
        "GF14",
        "GF15",
        "GF20",
        "GF32",
        "GF33",
        "GF34",
        "GF35",
        "GF36",
        "GF37"
      ],
      "verification": [
        {
          "id": "t01-probe-01-run",
          "performer": "claude:gf-t01-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "46b8be48be138daa2cfc6e7648e239a442c6d0fd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t01-author-20261029",
            "digest": "sha256:50a97210eb09d66a7d42cb38d2c169b7485a6246191d1201684770c1f759102f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/gittest.go,swap,failed,1,TestRunFailureReturn,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestRunFailureReturn,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,371\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,8,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestRunFailureReturn,\"commands_test.go:134: Run(absent) = \\\"fatal: Needed a single revision\\\", <nil>; want empty output and the wrapped process error\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "63e498366d28674bb6a390f903391e8019bcea5f",
          "requirement": "probe-01-run",
          "command": "bench probe internal/gittest/gittest.go --swap 'return \"\", fmt.Errorf(\"git %s: %w: %s\", strings.Join(args, \" \"), err, out)' --with 'return strings.TrimSpace(string(out)), nil' --package ./internal/gittest --run TestRunFailureReturn",
          "exit_code": 0,
          "probe": {
            "mutation": "Return a failed command as success. The goroutine then receives no error, so the returned-failure test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t01-author-20261029",
              "digest": "sha256:50a97210eb09d66a7d42cb38d2c169b7485a6246191d1201684770c1f759102f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/gittest.go,swap,failed,1,TestRunFailureReturn,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestRunFailureReturn,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,371\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,8,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestRunFailureReturn,\"commands_test.go:134: Run(absent) = \\\"fatal: Needed a single revision\\\", <nil>; want empty output and the wrapped process error\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01-probe-01-identity",
          "performer": "claude:gf-t01-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "46b8be48be138daa2cfc6e7648e239a442c6d0fd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t01-author-20261029",
            "digest": "sha256:747f858e2d905d9ac0d82003474513091eba3207bebfc722d59eee5804dabeec",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFixtureIdentity,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestFixtureIdentity,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,337\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,21,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestFixtureIdentity,\"commands_test.go:301: git commit --quiet -m fixture: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "63e498366d28674bb6a390f903391e8019bcea5f",
          "requirement": "probe-01-identity",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/gittest --run TestFixtureIdentity",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t01-author-20261029",
              "digest": "sha256:747f858e2d905d9ac0d82003474513091eba3207bebfc722d59eee5804dabeec",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFixtureIdentity,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestFixtureIdentity,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,337\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,21,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestFixtureIdentity,\"commands_test.go:301: git commit --quiet -m fixture: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
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
