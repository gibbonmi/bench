# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/git-query-cache/spec.md",
  "plan_digest": "sha256:7d3a811db66c80945a137d70bfb8b0c7995202a3a5bf521a5e9bca38fb700356",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GQ-C1",
      "base": "f979bae6807f82e9bc4a259b78076c3170710597",
      "tip": "de8410c563e7ea8e2fd2d7a1b2bebac97fd646f1",
      "plan_digest": "sha256:7d3a811db66c80945a137d70bfb8b0c7995202a3a5bf521a5e9bca38fb700356",
      "source_digest": "ab9506e1be46c9e86a6302b435dd4a45219d6a7e",
      "acceptance_rows": [
        "GQ1",
        "GQ2",
        "GQ3",
        "GQ4",
        "GQ5",
        "GQ6",
        "GQ7",
        "GQ8",
        "GQ9",
        "GQ10",
        "GQ13",
        "GQ14",
        "GQ15",
        "GQ16",
        "GQ17",
        "GQ18",
        "GQ19",
        "GQ20",
        "GQ21",
        "GQ22",
        "GQ23",
        "GQ24",
        "GQ25",
        "GQ26",
        "GQ43",
        "GQ44",
        "GQ45",
        "GQ47"
      ],
      "verification": [
        {
          "id": "t01-probe-01-hit",
          "performer": "claude:gq-t01-author-20261028",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ab9506e1be46c9e86a6302b435dd4a45219d6a7e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gq-t01-author-20261028",
            "digest": "sha256:f8f661e0cfed631a1c568308ddfda0ddde5bc474b41dc8ecf4e4c59b31d25213",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,7,TestRunObjectReads/absent_path,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReads,passed,11\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1424\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,67,11\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/absent_path,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/ancestor,\"run_object_test.go:110: starts of \\\"--is-ancestor 8103df5ae7d83f834a8d236b36d62eaee5a4ead6\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/blob,\"run_object_test.go:110: starts of \\\"cat-file blob ce013625030ba8dba906f756967f9e9ca394464a\\\" = 4, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/commit,\"run_object_test.go:110: starts of \\\"rev-parse --verify --quiet 8103df5ae7d83f834a8d236b36d62eaee5a4ead6\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/non-ancestor,\"run_object_test.go:110: starts of \\\"--is-ancestor 9324b9253cfba87b42cd827ede089247cbcfad3f\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_file,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_path,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "de8410c563e7ea8e2fd2d7a1b2bebac97fd646f1",
          "requirement": "probe-01-hit",
          "command": "bench probe internal/git/run.go --swap 'if raw, ok := r.objects[key]; ok {' --with 'if raw, ok := r.objects[key]; false && ok {' --package ./internal/git --run TestRunObjectReads",
          "exit_code": 0,
          "probe": {
            "mutation": "Never answer an object read from the store. Each second read then starts git, so the object-read test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gq-t01-author-20261028",
              "digest": "sha256:f8f661e0cfed631a1c568308ddfda0ddde5bc474b41dc8ecf4e4c59b31d25213",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,7,TestRunObjectReads/absent_path,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReads,passed,11\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1424\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,67,11\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/absent_path,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/ancestor,\"run_object_test.go:110: starts of \\\"--is-ancestor 8103df5ae7d83f834a8d236b36d62eaee5a4ead6\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/blob,\"run_object_test.go:110: starts of \\\"cat-file blob ce013625030ba8dba906f756967f9e9ca394464a\\\" = 4, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/commit,\"run_object_test.go:110: starts of \\\"rev-parse --verify --quiet 8103df5ae7d83f834a8d236b36d62eaee5a4ead6\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/non-ancestor,\"run_object_test.go:110: starts of \\\"--is-ancestor 9324b9253cfba87b42cd827ede089247cbcfad3f\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_file,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_path,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01-probe-01-live",
          "performer": "claude:gq-t01-author-20261028",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ab9506e1be46c9e86a6302b435dd4a45219d6a7e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gq-t01-author-20261028",
            "digest": "sha256:fed66f682741bdb6213eefbd1c8c4b1420aec6c0cb6dd49364a3ce1eca64313b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,1,TestRunObjectReadsStayLive,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReadsStayLive,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,510\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,31,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReadsStayLive,\"run_object_test.go:166: reads after the commit and the update-ref = \\\"hello\\\\\\\\n|100644 blob 5b69858846f588e863c6d82fbcf630f2e0adb550\\\\\\\\ta.txt|18c9d8cb4d427298e0689ab1cb3c2ca91d0d6ff8|true\\\", want \\\"next\\\\\\\\n|100644 blob 5b69858846f588e863c6d82fbcf630… (433 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "de8410c563e7ea8e2fd2d7a1b2bebac97fd646f1",
          "requirement": "probe-01-live",
          "command": "bench probe internal/git/run.go --swap 'if !FullObjectID(object) {' --with 'if false {' --package ./internal/git --run TestRunObjectReadsStayLive",
          "exit_code": 0,
          "probe": {
            "mutation": "Cache a blob read whose operand is a ref. The second read of the moved ref then returns the old bytes, so the live-read test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gq-t01-author-20261028",
              "digest": "sha256:fed66f682741bdb6213eefbd1c8c4b1420aec6c0cb6dd49364a3ce1eca64313b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,1,TestRunObjectReadsStayLive,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReadsStayLive,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,510\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,31,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReadsStayLive,\"run_object_test.go:166: reads after the commit and the update-ref = \\\"hello\\\\\\\\n|100644 blob 5b69858846f588e863c6d82fbcf630f2e0adb550\\\\\\\\ta.txt|18c9d8cb4d427298e0689ab1cb3c2ca91d0d6ff8|true\\\", want \\\"next\\\\\\\\n|100644 blob 5b69858846f588e863c6d82fbcf630… (433 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
