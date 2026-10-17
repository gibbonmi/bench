# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/anchor-admission-prune/spec.md",
  "plan_digest": "sha256:294bb59a470b025b0faa1f465f8561a826966e0879b0b69714709866831db182",
  "implementation_session": "",
  "chunks": [
    {
      "id": "AP-C2",
      "base": "ca0e6ec76eac0a2ef65d22c1c5405b948bb15dcf",
      "tip": "1e169d060fc785115f5b6f42eff6a83761bae863",
      "plan_digest": "sha256:294bb59a470b025b0faa1f465f8561a826966e0879b0b69714709866831db182",
      "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
      "acceptance_rows": [
        "AP02",
        "AP03",
        "AP04",
        "AP10",
        "AP25",
        "AP26",
        "AP27",
        "AP31",
        "AP46"
      ],
      "verification": [
        {
          "id": "t1-t1-anchors",
          "performer": "claude:aap_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t01",
            "digest": "sha256:bc1cf1f92ded0de4ca5c3c36229d3c25cdecbb3a885206be420b6f556112370d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,1723,181\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t1-t1-docs-currency-workflow",
          "performer": "claude:aap_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t01",
            "digest": "sha256:3edeee72ada88b9370cf9fdbb32754a73daff4a95a2e70e45ab119c6b780b173",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4600,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-docs-currency-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t1-t1-conformance",
          "performer": "claude:aap_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t01",
            "digest": "sha256:56e103fd930c1fdf24313d2c68701c5fcebf4aa6138e998f454516c68cdfef31",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,51490,1561\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-t1-cmd-bench",
          "performer": "claude:aap_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t01",
            "digest": "sha256:f61428016e9012f55ecea1e9090d037635d28457632dbc4c78e2dc1a6fc59366",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14734,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1-t1-root-conformance",
          "performer": "claude:aap_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t01",
            "digest": "sha256:fa8c8775eb396a0cee658a7df678d7e7f43fd85231509c8fc6e09314f6c8009d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14785,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
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
