# Review outcomes

## TP-C1a review pickup

The TP-C1a review ran on the frozen pair `16970a22..068ab04a`. Each axis ran on sonnet at high effort, by the reviewer's direction. The chunk has used 1 of its 2 repair cycles. The repair sessions `claude:ft290_t2_r1` and `claude:ft290_t3_r1` repaired the six targets at `3a97977c` and `d49b0697`, and the chunk tip is now `d49b0697`. The confirming round of all three axes passed with zero findings, and each earlier silent mutation now bites.

Advice from the confirming Coverage axis, with no finding ID: no test pins a nonzero `Ran` on the build-failed outcome, and no row requires one.

The raw count is 7 findings, and the repair-target count is 6, because S1 has no repair target. Ticket 2 owns C3. Ticket 3 owns S2, S3, C1, C2, and C4, because its test file holds the second helper and most of the tags.

### Standards

Count: 3. Worst issue: S2.

- S1, no-op, confidence 6. `named_check.go` keeps the `request.run` pass-through arm. The spec section that ticket 3 implements says "When `--run` is present and no test ran, the current `go test reported no test runs` refusal wins". The arm passes the `runGoTest` refusal through and derives no second copy of it. TP21 of ticket 6 reaches it.
- S2, auto-fix, confidence 5. `tests_run_test.go` and `check_row_test.go` each define the same canned-events `Command` helper, and `TestPackageRunWithNoTestKeepsExitZero` repeats it inline. Keep one helper.
- S3, auto-fix, confidence 5. The new test comments carry `(Coverage row TPn.)` tags, which are provenance under `craft-comments`. Delete the new tags.

### Spec

Count: 0. Worst issue: none. All ten TP-C1a rows are closed.

### Coverage

Count: 4. Worst issue: C1.

- C1, auto-fix, confidence 6. No test pins the `check` row on an `OutcomeFailed` named-check result. Add a test over the canned failing set.
- C2, auto-fix, confidence 6. Only streams with one run event pin the `tests_run` cell of the `check` row. Add a stream with two distinct run events.
- C3, auto-fix, confidence 4. No test pins the `e.Test != ""` guard on run events, which the zero rule depends on. Add a package-level run event to a zero-rule stream.
- C4, auto-fix, confidence 4. No test pins the interrupt pass-through arm or a failure with no run event. Add both tests.

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
      "tip": "d49b069704efe603a84203bdeffa9614c4802c37",
      "plan_digest": "sha256:e7275892d011a725f88fd6f4f231059449f393bf693fbd5723403f2c7ded0c8e",
      "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
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
        },
        {
          "id": "t3-testreport-v1",
          "performer": "claude:ft290_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ce41587b1bc16e9ce6820db7d98c466c03b7e510",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t3",
            "digest": "sha256:b3860d001760ee40ed4f0058e8ee5b25e8c557c6288fd88224ff84a0e761f0c4",
            "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\ntree: ft290-test-projection,44d7ba54c598dcea0a1093cff6527896a31868b4,false\nexit 0\n  github.com/gibbonmi/bench/internal/testreport,pass,34523\nfailures[0]{package,test,line}:\n\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/testreport/named_check.go --omit $'if outcome.Kind == OutcomeNoTestRun {\\n\\t\\treturn outcome, row + toon.Errorf(\"named check ran nothing\", \"no test emitted a run event\") + \"\\\\n\", 1\\n\\t}' --package ./internal/testreport --run '^TestNamedCheckRanNothingExitsOne$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/named_check.go,omit,failed,1,yes\n"
          },
          "requirement": "t3-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Omission at the zero rule: omit the OutcomeNoTestRun branch that prints the named check ran nothing title and exits 1. TestNamedCheckRanNothingExitsOne must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t3",
              "digest": "sha256:b3860d001760ee40ed4f0058e8ee5b25e8c557c6288fd88224ff84a0e761f0c4",
              "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\ntree: ft290-test-projection,44d7ba54c598dcea0a1093cff6527896a31868b4,false\nexit 0\n  github.com/gibbonmi/bench/internal/testreport,pass,34523\nfailures[0]{package,test,line}:\n\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/testreport/named_check.go --omit $'if outcome.Kind == OutcomeNoTestRun {\\n\\t\\treturn outcome, row + toon.Errorf(\"named check ran nothing\", \"no test emitted a run event\") + \"\\\\n\", 1\\n\\t}' --package ./internal/testreport --run '^TestNamedCheckRanNothingExitsOne$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/named_check.go,omit,failed,1,yes\n"
            }
          }
        },
        {
          "id": "t3-probe-v1",
          "performer": "claude:ft290_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ce41587b1bc16e9ce6820db7d98c466c03b7e510",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t3",
            "digest": "sha256:d52f2bb7549f53ef86f2af436e5c403c73ef46ef433deffaee73d14116ba977b",
            "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/probe\ntree: ft290-test-projection,44d7ba54c598dcea0a1093cff6527896a31868b4,false\nexit 0\n  github.com/gibbonmi/bench/internal/probe,pass,16591\nfailures[0]{package,test,line}:\n"
          },
          "requirement": "t3-probe",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        },
        {
          "id": "t3-ordinary-build-census-v1",
          "performer": "claude:ft290_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ce41587b1bc16e9ce6820db7d98c466c03b7e510",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t3",
            "digest": "sha256:643625bb9af0a53371a74b1fc341fd60773b2a2da07bd3f2d13beaab641e43d2",
            "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- bench test --check ordinary-build-census\ntree: ft290-test-projection,44d7ba54c598dcea0a1093cff6527896a31868b4,false\nexit 0\n  github.com/gibbonmi/bench/internal/conformance,pass,276\nfailures[0]{package,test,line}:\n"
          },
          "requirement": "t3-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t2-testreport-r1",
          "performer": "claude:ft290_t2_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t2_r1",
            "digest": "sha256:bece0e23c5426cca695c3697314cf46846834e4a1299d3977dcf29818966113a",
            "excerpt": "HEAD d49b069704efe603a84203bdeffa9614c4802c37\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,pass,32025\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\nexit 0\n\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/testreport/testreport.go --swap 'len(r.ranTests[pkg])' --with '0' --package ./internal/testreport --run '^TestPackagesRowCountsRunEvents$'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/testreport.go,swap,failed,1,yes\nfailures[1]: TestPackagesRowCountsRunEvents got \"canned,pass,250,0\", want \"canned,pass,250,2\"\nexit 0\n"
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
              "ref": "claude-agent:ft290_t2_r1",
              "digest": "sha256:bece0e23c5426cca695c3697314cf46846834e4a1299d3977dcf29818966113a",
              "excerpt": "HEAD d49b069704efe603a84203bdeffa9614c4802c37\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,pass,32025\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\nexit 0\n\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/testreport/testreport.go --swap 'len(r.ranTests[pkg])' --with '0' --package ./internal/testreport --run '^TestPackagesRowCountsRunEvents$'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/testreport.go,swap,failed,1,yes\nfailures[1]: TestPackagesRowCountsRunEvents got \"canned,pass,250,0\", want \"canned,pass,250,2\"\nexit 0\n"
            }
          }
        },
        {
          "id": "t2-probe-r1",
          "performer": "claude:ft290_t2_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t2_r1",
            "digest": "sha256:cfffc5a7ac8d0e3b9755400431c057bb167829c740a384fb5717f6b1e7138911",
            "excerpt": "HEAD d49b069704efe603a84203bdeffa9614c4802c37\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/probe\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/probe,pass,14897\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\nexit 0\n"
          },
          "requirement": "t2-probe",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        },
        {
          "id": "t3-testreport-r1",
          "performer": "claude:ft290_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t3_r1",
            "digest": "sha256:e0d9420fecd55d07a18864422c9af76f74f08cb2c5082fdd77369643ed6eb4ed",
            "excerpt": "HEAD d49b069704efe603a84203bdeffa9614c4802c37\n$ bench test --package ./internal/testreport\n  internal/testreport,pass,32646 ms; failures[0]; skips[0]; exit 0\n$ bench probe internal/testreport/named_check.go --omit '<OutcomeNoTestRun zero-rule branch>' --package ./internal/testreport --run '^TestNamedCheckRanNothingExitsOne$'\n  bit,internal/testreport/named_check.go,omit,failed,1,yes\n"
          },
          "requirement": "t3-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Omission at the zero rule: omit the OutcomeNoTestRun branch that prints the named check ran nothing title and exits 1. TestNamedCheckRanNothingExitsOne must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t3_r1",
              "digest": "sha256:e0d9420fecd55d07a18864422c9af76f74f08cb2c5082fdd77369643ed6eb4ed",
              "excerpt": "HEAD d49b069704efe603a84203bdeffa9614c4802c37\n$ bench test --package ./internal/testreport\n  internal/testreport,pass,32646 ms; failures[0]; skips[0]; exit 0\n$ bench probe internal/testreport/named_check.go --omit '<OutcomeNoTestRun zero-rule branch>' --package ./internal/testreport --run '^TestNamedCheckRanNothingExitsOne$'\n  bit,internal/testreport/named_check.go,omit,failed,1,yes\n"
            }
          }
        },
        {
          "id": "t3-probe-r1",
          "performer": "claude:ft290_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t3_r1",
            "digest": "sha256:ee27b86a94af51c77ad31d0e5db4aa51a365c614496147bf011a33e6bbf0695a",
            "excerpt": "HEAD d49b069704efe603a84203bdeffa9614c4802c37\n$ bench test --package ./internal/probe\n  internal/probe,pass,16206 ms; failures[0]; skips[0]; exit 0\n"
          },
          "requirement": "t3-probe",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        },
        {
          "id": "t3-ordinary-build-census-r1",
          "performer": "claude:ft290_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t3_r1",
            "digest": "sha256:74402d7e68ce7395cbaa1ed1c9118797331e609f99eeee476f7f567ddd1da76e",
            "excerpt": "HEAD d49b069704efe603a84203bdeffa9614c4802c37\n$ bench test --check ordinary-build-census\n  internal/conformance,pass,282 ms; failures[0]; skips[0]; exit 0\n"
          },
          "requirement": "t3-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t1-testreport-v2",
          "performer": "claude:ft290_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t1",
            "digest": "sha256:5edc1f9c770198167e4579f3514aff0e78f849a4dc235bbbafed48e73551febb",
            "excerpt": "source: d49b069704efe603a84203bdeffa9614c4802c37 (chunk TP-C1a, after review repair)\n$ bench test --package ./internal/testreport\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,pass,33197\nfailures[0] skips[0]\n$ bench probe internal/testreport/named_check.go --omit ' + \"\\n\" + namedCheckInventory()' --package ./internal/testreport --run '^TestUnknownNamedCheckReportsOperandAndInventory$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/named_check.go,omit,failed,1,yes\nfailure: TestUnknownNamedCheckReportsOperandAndInventory check_test.go:217 (refusal lost the checks inventory)\n"
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
              "digest": "sha256:5edc1f9c770198167e4579f3514aff0e78f849a4dc235bbbafed48e73551febb",
              "excerpt": "source: d49b069704efe603a84203bdeffa9614c4802c37 (chunk TP-C1a, after review repair)\n$ bench test --package ./internal/testreport\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,pass,33197\nfailures[0] skips[0]\n$ bench probe internal/testreport/named_check.go --omit ' + \"\\n\" + namedCheckInventory()' --package ./internal/testreport --run '^TestUnknownNamedCheckReportsOperandAndInventory$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/named_check.go,omit,failed,1,yes\nfailure: TestUnknownNamedCheckReportsOperandAndInventory check_test.go:217 (refusal lost the checks inventory)\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c1a-standards-r1",
          "performer": "claude:ft290_c1a_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "ce41587b1bc16e9ce6820db7d98c466c03b7e510",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c1a_standards",
            "digest": "sha256:2699f31af56b947b58f7c9ffda61fc79a260bdd8d7339143e09018c9db74d5f0",
            "excerpt": "axis=Standards findings=3 worst=S1\nS1 | auto-fix | conf 6 | internal/testreport/named_check.go:57,61 | `request.run != \"\" && outcome.Kind == OutcomeNoTestRun` in runNamedCheck is dead: parseFocusedRequest refuses --check with --run (command.go:71), and runGoTest owns the --run no-test refusal (command.go:251). One source per fact; Speculative Generality. | Delete the clause and the \"--run refusal\" comment phrase.\nS2 | auto-fix | conf 5 | tests_run_test.go:9-14 vs check_row_test.go:13-18,67-69,83-92 | Tickets 2 and 3 each wrote the same canned-events Command helper; TestPackageRunWithNoTestKeepsExitZero repeats the install pair inline. Fixture harness pasted N times. | Keep one helper taking root, args, events, exit.\nS3 | auto-fix | conf 5 | check_row_test.go:21,76; tests_run_test.go:18,52 | New \"(Coverage row TPn.)\" comment tags are provenance (craft-comments \"Treat an identifier as provenance\"). | Delete the parentheticals, or close as no-op if the convention is kept.\nAdvice: selection_facts.go:44 could derive the probe-target answer from namedCheckKind; header literal expectations have no recorded header-mutation red; ranTests is the single run-count source; namedCheckKind is the single kind owner.\nExamined: --check-current, git diff 16970a22..068ab04a, command.go, outcome.go, named_check.go, craft-comments, craft-review; no tests run.\n"
          },
          "axis": "Standards",
          "base": "16970a226a4e36f2d3d1eee5b5e262f0c76de548",
          "tip": "068ab04a456b87e56655e0a500525322919efcc1",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "c1a-spec-r1",
          "performer": "claude:ft290_c1a_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "ce41587b1bc16e9ce6820db7d98c466c03b7e510",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c1a_spec",
            "digest": "sha256:d37d1866926088809f6dff9c2e94df508774a295045713c2b460a696dbb33d7e",
            "excerpt": "axis=Spec findings=0 worst=none\nAll ten TP-C1a rows closed by static read: TP1 TestPackagesRowCountsRunEvents, TP2 TestPackagesRowNoTestsCountsZero, TP3 TestNamedCheckPrintsCheckRowFirst, TP4 TestSystemCheckRowKind, TP5 TestNamedCheckRanNothingExitsOne, TP51 TestNamedCheckBuildFailureWinsOverZeroRule, TP6 TestPackageRunWithNoTestKeepsExitZero, TP55 TestChangedRunWithNoTestKeepsExitZero, TP47 TestPackageFormHasNoSelectedBy, TP18 TestNamedCheckRunsOnlyRegisteredDevScope (in place).\nTitle, check-row-first, OutcomeNoTestRun keying, no row before a verdict, and the --run refusal pass-through all match the spec. No scope creep; posture-change table fully applied.\nAdvice: no row grades the Refused/Interrupted arms of runNamedCheck; the --run pass-through arm is unreachable until ticket 6 (TP21), justified by the spec sentence; the plan edit is orchestrator record keeping.\nExamined: --check-current, bench coverage, full diff, spec sections, tickets 1-3, outcome.go, command.go, test helpers. No tests run.\n"
          },
          "axis": "Spec",
          "base": "16970a226a4e36f2d3d1eee5b5e262f0c76de548",
          "tip": "068ab04a456b87e56655e0a500525322919efcc1",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1a-coverage-r1",
          "performer": "claude:ft290_c1a_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "ce41587b1bc16e9ce6820db7d98c466c03b7e510",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c1a_coverage",
            "digest": "sha256:b6a8dda42b078c807b965e8d4cc7773043281ffe63b2b38f926dbb18be20430c",
            "excerpt": "axis=Coverage findings=4 worst=C1\nC1 | auto-fix | conf 6 | named_check.go:60-61; spec \"first block of each named-check result that reached a verdict\" | No test pins the check row on an OutcomeFailed result; adding OutcomeFailed to the passthrough arm probed silent. | Add a named-check test over the canned failing set asserting the check header and a tests_run 1 row as prefix, exit 1.\nC2 | auto-fix | conf 6 | named_check.go:62 | Only single-run-event streams pin the check-row tests_run cell; min(outcome.Ran, 1) probed silent. | Carry two distinct run events in TP3 or TP5 and assert line-routing,conformance,2,0.\nC3 | auto-fix | conf 4 | testreport.go:72 | The e.Test != \"\" guard on run events is unpinned; a package-level run event would defeat the zero rule; dropping the guard probed silent. | Add a package-level run event to the TP5 stream; assert exit 1, tests_run 0.\nC4 | auto-fix | conf 4 | named_check.go:60 | (a) Interrupt passthrough unpinned (removal probed silent); (b) a fail-with-no-run-event stream is unpinned (Ran==0 && !BuildFailed && !Passed swap probed silent). | Add an installSignallingGo named-check interrupt test asserting no check[ prefix; add a fail-without-run stream asserting failures, not the zero-rule title.\nCovered: TP55 and TP6 bit under a zero rule on every form; a per-package cap bit outcome tests.\nAdvice: an all-skipped named check prints tests_run 1 at exit 0 (matches the spec); the --run arm belongs to ticket 6 (TP21).\nExamined: delta 16970a22..068ab04a, spec rows and edge sections, tickets 1-3, outcome.go, testreport.go, command.go, environment.go, consumers; 9 bench probe calls, all restored yes; git status clean.\n"
          },
          "axis": "Coverage",
          "base": "16970a226a4e36f2d3d1eee5b5e262f0c76de548",
          "tip": "068ab04a456b87e56655e0a500525322919efcc1",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": []
        },
        {
          "id": "c1a-standards-c1",
          "performer": "claude:ft290_c1a_standards_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c1a_standards_c1",
            "digest": "sha256:3618ad0088e61de3329274d3bf4f8ae69d8690845abab232b6e5ecf72378d575",
            "excerpt": "axis=Standards findings=0 worst=none (confirming round, repair delta 068ab04a..d49b0697)\nS2 confirmed: commandOverEvents (tests_run_test.go:8-14) is the one canned-events helper; no namedCheckOverEvents remains; the interrupt test uses installSignallingGo, a different stub; twoTestsRunEvents is one shared stream.\nS3 confirmed: rg \"Coverage row TP\" internal has no hits; new comments state current behavior.\nS1 no-op confirmed; no binding rule contests it.\nAdvice: repeated want literal and argv across check-row tests is incidental test text.\nExamined: repair delta, outcome_test.go 100-229, craft-comments, rg for helper and tag patterns, --check-current at f05830c6. No tests run.\n"
          },
          "axis": "Standards",
          "base": "16970a226a4e36f2d3d1eee5b5e262f0c76de548",
          "tip": "d49b069704efe603a84203bdeffa9614c4802c37",
          "finding_ids": [],
          "supersedes": [
            "c1a-standards-r1"
          ]
        },
        {
          "id": "c1a-spec-c1",
          "performer": "claude:ft290_c1a_spec_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c1a_spec_c1",
            "digest": "sha256:91f4ec2e2f4768c382535bef9aae98021df48a2cdc17505de59719207b2721c4",
            "excerpt": "axis=Spec findings=0 worst=none (confirming round, repair delta 068ab04a..d49b0697)\nRepair is test-only: only check_row_test.go and tests_run_test.go change under internal/.\nAll ten TP-C1a rows keep a named test with unchanged assertions: TP1, TP2, TP47 in tests_run_test.go; TP3, TP4, TP5, TP6, TP51, TP55 in check_row_test.go; TP18 in check_test.go (ordinary-build-census,conformance,1,0).\nNew tests C1, C2, C3, C4a, C4b agree with spec \"The check row and the zero rule\" (spec.md:114-137).\nExamined: --check-current at f05830c6, repair delta, spec rows, outcome_test.go 195-264, production files by targeted rg. No tests run.\n"
          },
          "axis": "Spec",
          "base": "16970a226a4e36f2d3d1eee5b5e262f0c76de548",
          "tip": "d49b069704efe603a84203bdeffa9614c4802c37",
          "finding_ids": [],
          "supersedes": [
            "c1a-spec-r1"
          ]
        },
        {
          "id": "c1a-coverage-c1",
          "performer": "claude:ft290_c1a_coverage_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c1a_coverage_c1",
            "digest": "sha256:62ffab055d65ec508a525bdcc049d80ba2a124a7c714300fe935d736331acd4a",
            "excerpt": "axis=Coverage findings=0 worst=none (confirming round, repair delta 068ab04a..d49b0697)\nAll five earlier silent mutations now bite under --package ./internal/testreport, restored=yes each:\nC1 OutcomeFailed pass-through -> bit (TestNamedCheckFailurePrintsCheckRowFirst)\nC2 min(outcome.Ran, 1) -> bit (TestNamedCheckRowCountsEachRunTest)\nC3 drop e.Test != \"\" -> bit (TestPackagesRowIgnoresUnnamedRunEvent)\nC4a omit interrupt clause -> bit (TestNamedCheckInterruptKeepsItsBytes)\nC4b Ran==0 && !BuildFailed && !Passed -> bit (TestNamedCheckFailureWithNoRunEventIsNotRanNothing)\nNew-bypass probes on outcome.go: cross-package sum and Failed-arm Ran both bit; the shared helper hides no per-test difference.\nAdvice: Ran on the OutcomeBuildFailed arm (outcome.go:96) is unpinned with a nonzero value; no row requires it.\nExamined: --check-current at f05830c6, repair delta, named_check.go, outcome.go, testreport.go; 8 probes; git status clean before and after.\n"
          },
          "axis": "Coverage",
          "base": "16970a226a4e36f2d3d1eee5b5e262f0c76de548",
          "tip": "d49b069704efe603a84203bdeffa9614c4802c37",
          "finding_ids": [],
          "supersedes": [
            "c1a-coverage-r1"
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
  }
}
```
