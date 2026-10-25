# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/no-go-schedule/spec.md",
  "plan_digest": "sha256:f7c4b9dfcb60734c5667d0a2306b537d6995e9fb0298efa86eb1c64756040342",
  "implementation_session": "",
  "chunks": [
    {
      "id": "NG-C1",
      "base": "d30ccffeac02152febd38db38dfccfdcf0f6cf7f",
      "tip": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
      "plan_digest": "sha256:f7c4b9dfcb60734c5667d0a2306b537d6995e9fb0298efa86eb1c64756040342",
      "source_digest": "eb1341d34a4433afe8b6c7856211f62971757177",
      "acceptance_rows": [
        "NG1",
        "NG2",
        "NG3",
        "NG4"
      ],
      "verification": [
        {
          "id": "t01-ticket-1-check-1",
          "performer": "claude:ngs-t01-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eb1341d34a4433afe8b6c7856211f62971757177",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ngs-t01-author-20261025",
            "digest": "sha256:07819a85f9898b12010b6603dda28329ecb6590c0040295487d977b1caadfcdd",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,64283,590\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
          "requirement": "ticket-1-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-2",
          "performer": "claude:ngs-t01-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eb1341d34a4433afe8b6c7856211f62971757177",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ngs-t01-author-20261025",
            "digest": "sha256:2ef9de0bcc32b863ebe4a896ec81f3352c1de3e921219a23a50607e527bcfdd9",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13386,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
          "requirement": "ticket-1-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-probe",
          "performer": "claude:ngs-t01-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eb1341d34a4433afe8b6c7856211f62971757177",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ngs-t01-author-20261025",
            "digest": "sha256:ab2ca06d9e31e8c328a36359d2d3aa1090b1e66f65dc731d8a7d5d8e3a032f55",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,swap,failed,1,TestOracleIdentityIgnoresTheTree,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,590\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,460\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,64459,590\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityIgnoresTheTree,\"oracle_identity_test.go:39: identity = 6abcf6ab77289023b8ec2e4521c3a7353fb34c08363560c135f24f9dd58e6ffe and de05407ab4ef765295e74e42747e1ed13e79aa5b360ab8ae66df730fe0c2f58d, want one identity across trees\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
          "requirement": "ticket-1-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "In buildSubjectForTree, frame the tree into the identity digest too. TestOracleIdentityIgnoresTheTree must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ngs-t01-author-20261025",
              "digest": "sha256:ab2ca06d9e31e8c328a36359d2d3aa1090b1e66f65dc731d8a7d5d8e3a032f55",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,swap,failed,1,TestOracleIdentityIgnoresTheTree,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,590\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,460\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,64459,590\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityIgnoresTheTree,\"oracle_identity_test.go:39: identity = 6abcf6ab77289023b8ec2e4521c3a7353fb34c08363560c135f24f9dd58e6ffe and de05407ab4ef765295e74e42747e1ed13e79aa5b360ab8ae66df730fe0c2f58d, want one identity across trees\",1\nskips[0]{package,test,reason}:\n"
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
