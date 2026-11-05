# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/production-test-seam-policy/spec.md",
  "plan_digest": "sha256:5c08a62494fc1fdd742618a16f4d9ccab37c50326fb6f594f7fef3f23db32a4b",
  "implementation_session": "",
  "chunks": [
    {
      "id": "TS1",
      "base": "77de12e23e7fc3f4745ac7132e1ceeae02ce396c",
      "tip": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
      "plan_digest": "sha256:5c08a62494fc1fdd742618a16f4d9ccab37c50326fb6f594f7fef3f23db32a4b",
      "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
      "acceptance_rows": [
        "PS01",
        "PS02",
        "PS03",
        "PS12",
        "PS13",
        "PS14",
        "PS41",
        "PS42",
        "PS43",
        "PS50",
        "PS51",
        "PS52",
        "PS53",
        "PS54",
        "PS55",
        "PS61"
      ],
      "verification": [
        {
          "id": "ts1-01-wrapper-bypass",
          "performer": "ft343-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-author-01",
            "digest": "sha256:5ff4c883590de2b57efe65fadff342ba3c996a22b74d4c1e1136f5892cbb5329",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/injected_ports_test.go,swap,failed,2,TestEveryRetainedFixtureBitesThroughRegisteredOwner,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestEveryRetainedFixtureBitesThroughRegisteredOwner$,passed,289\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,771\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,4108,289\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner,\"fixture_bite_test.go:33: completed fixture proofs do not equal the current producer: got 287, want 288\",1\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner/unregistered-port,\"fixture_bite_test.go:27: unregistered-port did not bite through owner injected-port-registry; want \\\"injected port has no registry row: internal/gitguard.Auditor is injected at a parameter of auditVerdict in internal/gitguard/gitguard.go\\\" in… (253 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "requirement": "wrapper-bypass",
          "command": "bench test --package ./internal/conformance --run '^TestEveryRetainedFixtureBitesThroughRegisteredOwner$'",
          "exit_code": 0,
          "probe": {
            "mutation": "make checkInjectedPortRegistry return without calling the relocated owner; require the injected-ports canary bite to fail with its exact expected diagnostic",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft343-author-01",
              "digest": "sha256:5ff4c883590de2b57efe65fadff342ba3c996a22b74d4c1e1136f5892cbb5329",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/injected_ports_test.go,swap,failed,2,TestEveryRetainedFixtureBitesThroughRegisteredOwner,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestEveryRetainedFixtureBitesThroughRegisteredOwner$,passed,289\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,771\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,4108,289\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner,\"fixture_bite_test.go:33: completed fixture proofs do not equal the current producer: got 287, want 288\",1\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner/unregistered-port,\"fixture_bite_test.go:27: unregistered-port did not bite through owner injected-port-registry; want \\\"injected port has no registry row: internal/gitguard.Auditor is injected at a parameter of auditVerdict in internal/gitguard/gitguard.go\\\" in… (253 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
