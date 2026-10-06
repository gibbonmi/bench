# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/ft290-test-projection/spec.md",
  "plan_digest": "sha256:37523705f05ff5539429fd7c190772de817b553748203d244e9ad3900e2c0a57",
  "implementation_session": "",
  "chunks": [
    {
      "id": "TP-C1a",
      "base": "16970a226a4e36f2d3d1eee5b5e262f0c76de548",
      "tip": "068ab04a456b87e56655e0a500525322919efcc1",
      "plan_digest": "sha256:37523705f05ff5539429fd7c190772de817b553748203d244e9ad3900e2c0a57",
      "source_digest": "ce41587b1bc16e9ce6820db7d98c466c03b7e510",
      "acceptance_rows": [
        "TP1",
        "TP2",
        "TP47",
        "TP3",
        "TP4",
        "TP5",
        "TP6",
        "TP18",
        "TP51",
        "TP55"
      ],
      "verification": [
        {
          "id": "t1-testreport-v1",
          "performer": "claude:ft290_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ce41587b1bc16e9ce6820db7d98c466c03b7e510",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t1",
            "digest": "sha256:f5eb11da73dca1bb2fd217f95cc0a2cfbc18edfb190aa6b222f9c2582500b00c",
            "excerpt": "source: 068ab04a456b87e56655e0a500525322919efcc1 (chunk TP-C1a)\n$ bench test --package ./internal/testreport\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,pass,32148\nfailures[0] skips[0]\n$ bench probe internal/testreport/named_check.go --omit ' + \"\\n\" + namedCheckInventory()' --package ./internal/testreport --run '^TestUnknownNamedCheckReportsOperandAndInventory$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/named_check.go,omit,failed,1,yes\nfailure: TestUnknownNamedCheckReportsOperandAndInventory check_test.go:217 (refusal lost the checks inventory)\n"
          },
          "requirement": "t1-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Omission at the moved unknown-name branch of the named-check owner: omit the namedCheckInventory call in the unknown-check refusal with bench probe --omit. TestUnknownNamedCheckReportsOperandAndInventory must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t1",
              "digest": "sha256:f5eb11da73dca1bb2fd217f95cc0a2cfbc18edfb190aa6b222f9c2582500b00c",
              "excerpt": "source: 068ab04a456b87e56655e0a500525322919efcc1 (chunk TP-C1a)\n$ bench test --package ./internal/testreport\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,pass,32148\nfailures[0] skips[0]\n$ bench probe internal/testreport/named_check.go --omit ' + \"\\n\" + namedCheckInventory()' --package ./internal/testreport --run '^TestUnknownNamedCheckReportsOperandAndInventory$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/named_check.go,omit,failed,1,yes\nfailure: TestUnknownNamedCheckReportsOperandAndInventory check_test.go:217 (refusal lost the checks inventory)\n"
            }
          }
        },
        {
          "id": "t2-testreport-v1",
          "performer": "claude:ft290_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ce41587b1bc16e9ce6820db7d98c466c03b7e510",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t2",
            "digest": "sha256:b53482d0f011d45abb9d4f8ab223f56c24706a5270c9dfdb9d1e9cc6c534c2e3",
            "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\nhead a1a0edc8cf4277005373022c16d9f577263e9cee dirty=false; exit 0\n  github.com/gibbonmi/bench/internal/testreport,pass,35376\nfailures[0]{package,test,line}:\n$ bench probe internal/testreport/testreport.go --swap 'len(r.ranTests[pkg])' --with '0' --package ./internal/testreport --run '^TestPackagesRowCountsRunEvents$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/testreport.go,swap,failed,1,yes\n"
          },
          "requirement": "t2-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the packages-row producer: replace the tests_run count of distinct run events with the constant 0. TestPackagesRowCountsRunEvents must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t2",
              "digest": "sha256:b53482d0f011d45abb9d4f8ab223f56c24706a5270c9dfdb9d1e9cc6c534c2e3",
              "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\nhead a1a0edc8cf4277005373022c16d9f577263e9cee dirty=false; exit 0\n  github.com/gibbonmi/bench/internal/testreport,pass,35376\nfailures[0]{package,test,line}:\n$ bench probe internal/testreport/testreport.go --swap 'len(r.ranTests[pkg])' --with '0' --package ./internal/testreport --run '^TestPackagesRowCountsRunEvents$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/testreport.go,swap,failed,1,yes\n"
            }
          }
        },
        {
          "id": "t2-probe-v1",
          "performer": "claude:ft290_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ce41587b1bc16e9ce6820db7d98c466c03b7e510",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t2",
            "digest": "sha256:faff6342e6d14688826d4a17fe45371e91c4349e597e9e36bec922f2cf2be1a2",
            "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/probe\nhead a1a0edc8cf4277005373022c16d9f577263e9cee dirty=false; exit 0\n  github.com/gibbonmi/bench/internal/probe,pass,17869\nfailures[0]{package,test,line}:\n"
          },
          "requirement": "t2-probe",
          "command": "bench test --package ./internal/probe",
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
