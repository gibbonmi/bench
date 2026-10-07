# Review outcomes

## TP-C1b review pickup

The TP-C1b review ran on the frozen pair `d49b0697..5a1c3c3b`. Each axis ran on sonnet at high effort. The chunk has used 1 of its 2 repair cycles. The repair sessions `claude:ft290_t4_r1` and `claude:ft290_t5_r1` repaired the three targets at `fd6c3e07` and `71582c05`, and the chunk tip is now `71582c05`. The confirming round of all three axes passed with zero findings, and each earlier silent mutation now bites in its owning package.

The confirming Coverage axis gave advice with no finding ID. The zero-subject prose test matches by prefix, so a suffix that only `--full` adds stays green. No row binds that case.

The raw count is 3 findings, and the repair-target count is 3. Ticket 4 owns C1 and C2. Ticket 5 owns S1, because ticket 5 changed the unit of `Outcome.FailedTests`.

### Standards

Count: 1. Worst issue: S1.

- S1, auto-fix, confidence 5. The `OutcomeFailed` doc comment in `outcome.go` still names a failing test row, but a `--full` row is now one diagnostic line. The edited `Outcome` comment is also one long line. Reword the comment and wrap it.

### Spec

Count: 0. Worst issue: none. All twelve TP-C1b rows are closed. No spec line contradicts the `OutcomeFailed` kind of a prose grader refusal.

### Coverage

Count: 2. Worst issue: C1.

- C1, auto-fix, confidence 9. No test runs a red prose tree with `--full`, and a swap that drops the findings under `--full` stayed green. Add `--full` cases to the red, refusal, and zero-subject prose tests.
- C2, auto-fix, confidence 8. The subject test uses only a file exclusion row, and a swap that ignores directory rows stayed green. This repository's exclusion file uses directory rows. Add a directory row and a file under it.

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
  "plan_digest": "sha256:0a578218a0d886a9acb001ee1ccf0168a84abd0b6589b1d3cd4de21a1337e123",
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
        },
        {
          "id": "t1-testreport-v3",
          "performer": "claude:ft290_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t1",
            "digest": "sha256:b1f072c8ed444d64df88435c5bffb19d7254414741a6a9df963cdfe4c086fdf3",
            "excerpt": "source: 72c5cbc4e94dd5ad6ece7616431f73b74daf42f2 (chunk TP-C1a), tree clean\n$ bench test --package ./internal/testreport\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,pass,35016\nfailures[0] skips[0]\n$ bench probe internal/testreport/named_check.go --omit ' + \"\\n\" + namedCheckInventory()' --package ./internal/testreport --run '^TestUnknownNamedCheckReportsOperandAndInventory$'\nbench probe exit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/named_check.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,^TestUnknownNamedCheckReportsOperandAndInventory$,passed,1\npackages[1]: github.com/gibbonmi/bench/internal/testreport,fail,3\nfailure: TestUnknownNamedCheckReportsOperandAndInventory check_test.go:217 (refusal lost the checks inventory)\nmutated run exit code: the probe output prints no exit-code cell for the mutated run.\nRecorded probe exit code 1, inferred from cause=failed, failed_tests=1, and package status fail.\n"
          },
          "requirement": "t1-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Omission at the moved unknown-name branch of the named-check owner: omit the namedCheckInventory call in the unknown-check refusal with bench probe --omit. TestUnknownNamedCheckReportsOperandAndInventory must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t1",
              "digest": "sha256:b1f072c8ed444d64df88435c5bffb19d7254414741a6a9df963cdfe4c086fdf3",
              "excerpt": "source: 72c5cbc4e94dd5ad6ece7616431f73b74daf42f2 (chunk TP-C1a), tree clean\n$ bench test --package ./internal/testreport\nexit 0\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,pass,35016\nfailures[0] skips[0]\n$ bench probe internal/testreport/named_check.go --omit ' + \"\\n\" + namedCheckInventory()' --package ./internal/testreport --run '^TestUnknownNamedCheckReportsOperandAndInventory$'\nbench probe exit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/named_check.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,^TestUnknownNamedCheckReportsOperandAndInventory$,passed,1\npackages[1]: github.com/gibbonmi/bench/internal/testreport,fail,3\nfailure: TestUnknownNamedCheckReportsOperandAndInventory check_test.go:217 (refusal lost the checks inventory)\nmutated run exit code: the probe output prints no exit-code cell for the mutated run.\nRecorded probe exit code 1, inferred from cause=failed, failed_tests=1, and package status fail.\n"
            }
          }
        },
        {
          "id": "t2-testreport-r2",
          "performer": "claude:ft290_t2_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t2_r1",
            "digest": "sha256:93af44ce3e2170196ce75606fefae53d2e0c7ffe027ac54cc67ef33e95678815",
            "excerpt": "HEAD 72c5cbc4e94dd5ad6ece7616431f73b74daf42f2\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/testreport/testreport.go --swap 'len(r.ranTests[pkg])' --with '0' --package ./internal/testreport --run '^TestPackagesRowCountsRunEvents$'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/testreport.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,fail,58\nfailures[1]: TestPackagesRowCountsRunEvents got \"canned,pass,250,0\", want \"canned,pass,250,2\"\nMutated focused run exit code: 1, inferred from verdict bit, cause failed, failed_tests 1, and package status fail (bench probe prints no exit cell for the mutated run).\n\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,pass,31944\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\nexit 0\n"
          },
          "requirement": "t2-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the packages-row producer: replace the tests_run count of distinct run events with the constant 0. TestPackagesRowCountsRunEvents must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t2_r1",
              "digest": "sha256:93af44ce3e2170196ce75606fefae53d2e0c7ffe027ac54cc67ef33e95678815",
              "excerpt": "HEAD 72c5cbc4e94dd5ad6ece7616431f73b74daf42f2\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/testreport/testreport.go --swap 'len(r.ranTests[pkg])' --with '0' --package ./internal/testreport --run '^TestPackagesRowCountsRunEvents$'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/testreport.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,fail,58\nfailures[1]: TestPackagesRowCountsRunEvents got \"canned,pass,250,0\", want \"canned,pass,250,2\"\nMutated focused run exit code: 1, inferred from verdict bit, cause failed, failed_tests 1, and package status fail (bench probe prints no exit cell for the mutated run).\n\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,pass,31944\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\nexit 0\n"
            }
          }
        },
        {
          "id": "t3-testreport-r2",
          "performer": "claude:ft290_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "369008241fa2cf8fa5c91120bc1c6cfc7c5fcfeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t3_r1",
            "digest": "sha256:63ad4b44b04ce59435d56501d60ed11f174cca83b8f44301559774707c04d99f",
            "excerpt": "HEAD 72c5cbc4e94dd5ad6ece7616431f73b74daf42f2\n$ bench probe internal/testreport/named_check.go --omit '<OutcomeNoTestRun zero-rule branch>' --package ./internal/testreport --run '^TestNamedCheckRanNothingExitsOne$'\n  probe: bit,internal/testreport/named_check.go,omit,failed,1,yes\n  mutated run: internal/testreport,fail,42 ms; failures[1] TestNamedCheckRanNothingExitsOne check_row_test.go:87\n  mutated-run exit 1: inferred from verdict bit, cause failed, failed_tests 1, package status fail (bench probe prints no exit cell)\n$ bench test --package ./internal/testreport\n  internal/testreport,pass,32417 ms; failures[0]; skips[0]; exit 0\n"
          },
          "requirement": "t3-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Omission at the zero rule: omit the OutcomeNoTestRun branch that prints the named check ran nothing title and exits 1. TestNamedCheckRanNothingExitsOne must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t3_r1",
              "digest": "sha256:63ad4b44b04ce59435d56501d60ed11f174cca83b8f44301559774707c04d99f",
              "excerpt": "HEAD 72c5cbc4e94dd5ad6ece7616431f73b74daf42f2\n$ bench probe internal/testreport/named_check.go --omit '<OutcomeNoTestRun zero-rule branch>' --package ./internal/testreport --run '^TestNamedCheckRanNothingExitsOne$'\n  probe: bit,internal/testreport/named_check.go,omit,failed,1,yes\n  mutated run: internal/testreport,fail,42 ms; failures[1] TestNamedCheckRanNothingExitsOne check_row_test.go:87\n  mutated-run exit 1: inferred from verdict bit, cause failed, failed_tests 1, package status fail (bench probe prints no exit cell)\n$ bench test --package ./internal/testreport\n  internal/testreport,pass,32417 ms; failures[0]; skips[0]; exit 0\n"
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
    },
    {
      "id": "TP-C1b",
      "base": "d49b069704efe603a84203bdeffa9614c4802c37",
      "tip": "71582c05f4028c1ed6a643f9747895251aed31a3",
      "plan_digest": "sha256:0a578218a0d886a9acb001ee1ccf0168a84abd0b6589b1d3cd4de21a1337e123",
      "source_digest": "6d30a8f09034e3ce40838dd97f09d8f380ddd8cc",
      "acceptance_rows": [
        "TP7",
        "TP8",
        "TP9",
        "TP10",
        "TP11",
        "TP12",
        "TP13",
        "TP14",
        "TP15",
        "TP16",
        "TP17",
        "TP49"
      ],
      "verification": [
        {
          "id": "t4-testreport-v1",
          "performer": "claude:ft290_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9dda6bd8e26b08d40413e6ff3299191f6bb9f9ee",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t4",
            "digest": "sha256:3fc87ffd1de421b63d19edd28f29306982a0c42f434f6cad30e95f04334cbbc7",
            "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\ntree: ft290-test-projection,5a1c3c3b8e1ffbbd04ff12c85e8d0dca56a6734a\nexit 0\n  github.com/gibbonmi/bench/internal/testreport,pass,30928,158\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t4-prose-v1",
          "performer": "claude:ft290_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9dda6bd8e26b08d40413e6ff3299191f6bb9f9ee",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t4",
            "digest": "sha256:36b7180cd6a3daab35153be458ab405dae97a847a133c70c2066b7b0a3b793b9",
            "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/prose\ntree: ft290-test-projection,5a1c3c3b8e1ffbbd04ff12c85e8d0dca56a6734a\nexit 0\n  github.com/gibbonmi/bench/internal/prose,pass,346,79\nfailures[0]{package,test,line,lines}:\n\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/prose/walk.go --omit $'\\t\\tif g.ex.excluded(rel) {\\n\\t\\t\\tcontinue\\n\\t\\t}\\n' --package ./internal/prose --run '^TestGradeReportsGradedSubjects$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/prose/walk.go,omit,failed,1,yes\n  github.com/gibbonmi/bench/internal/prose,fail,3,1\n  TestGradeReportsGradedSubjects: GradeTree() = {Subjects:[keep.md skip.md] Findings:[]}\n"
          },
          "requirement": "t4-prose",
          "command": "bench test --package ./internal/prose",
          "exit_code": 0,
          "probe": {
            "mutation": "Omission at the subject answer in internal/prose/walk.go: omit the exclusion-set test, so that an excluded file counts as a subject. TestGradeReportsGradedSubjects must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t4",
              "digest": "sha256:36b7180cd6a3daab35153be458ab405dae97a847a133c70c2066b7b0a3b793b9",
              "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/prose\ntree: ft290-test-projection,5a1c3c3b8e1ffbbd04ff12c85e8d0dca56a6734a\nexit 0\n  github.com/gibbonmi/bench/internal/prose,pass,346,79\nfailures[0]{package,test,line,lines}:\n\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/prose/walk.go --omit $'\\t\\tif g.ex.excluded(rel) {\\n\\t\\t\\tcontinue\\n\\t\\t}\\n' --package ./internal/prose --run '^TestGradeReportsGradedSubjects$'\nexit 0\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/prose/walk.go,omit,failed,1,yes\n  github.com/gibbonmi/bench/internal/prose,fail,3,1\n  TestGradeReportsGradedSubjects: GradeTree() = {Subjects:[keep.md skip.md] Findings:[]}\n"
            }
          }
        },
        {
          "id": "t4-prose-mechanics-v1",
          "performer": "claude:ft290_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9dda6bd8e26b08d40413e6ff3299191f6bb9f9ee",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t4",
            "digest": "sha256:7c53cd950aebecd80e75ba6ce86c57d233de13316e88e9f430942a8f28979e94",
            "excerpt": "$ bench worktree exec \"ft290-test-projection\" -- ./dist/bench test --check prose-mechanics\ntree: ft290-test-projection,5a1c3c3b8e1ffbbd04ff12c85e8d0dca56a6734a\nexit 0\ncheck[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,148,1\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-prose-mechanics",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        },
        {
          "id": "t5-testreport-v1",
          "performer": "claude:ft290_t5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9dda6bd8e26b08d40413e6ff3299191f6bb9f9ee",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t5",
            "digest": "sha256:1b97d55a5e0adcc7a2886f9aa168d8d4104003556c0e9f39c6c527e8bef95bcf",
            "excerpt": "tip 5a1c3c3b8e1ffbbd04ff12c85e8d0dca56a6734a\nbench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\nexit 0, wall 30412 ms\n  github.com/gibbonmi/bench/internal/testreport,pass,30412,158\nfailures[0]{package,test,line,lines}:\n\nprobe: bench probe internal/testreport/outcome.go --swap \"failed++\" --with \"failed += len(f.rows(true))\" --package ./internal/testreport --run '^TestFullFailedTestsCountsTests$'\n  bit,internal/testreport/outcome.go,swap,failed,1,yes\nmutated run: testreport,fail,45,1 (exit 1); TestFullFailedTestsCountsTests: FailedTests:3, want 1\n"
          },
          "requirement": "t5-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the failed-test count: count the --full failures rows in place of the distinct failed tests for Outcome.FailedTests. TestFullFailedTestsCountsTests must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t5",
              "digest": "sha256:1b97d55a5e0adcc7a2886f9aa168d8d4104003556c0e9f39c6c527e8bef95bcf",
              "excerpt": "tip 5a1c3c3b8e1ffbbd04ff12c85e8d0dca56a6734a\nbench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\nexit 0, wall 30412 ms\n  github.com/gibbonmi/bench/internal/testreport,pass,30412,158\nfailures[0]{package,test,line,lines}:\n\nprobe: bench probe internal/testreport/outcome.go --swap \"failed++\" --with \"failed += len(f.rows(true))\" --package ./internal/testreport --run '^TestFullFailedTestsCountsTests$'\n  bit,internal/testreport/outcome.go,swap,failed,1,yes\nmutated run: testreport,fail,45,1 (exit 1); TestFullFailedTestsCountsTests: FailedTests:3, want 1\n"
            }
          }
        },
        {
          "id": "t5-probe-v1",
          "performer": "claude:ft290_t5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9dda6bd8e26b08d40413e6ff3299191f6bb9f9ee",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t5",
            "digest": "sha256:07cbe6c738b82de4c95fb3dc733a650fe62933ea9c2a64d90f54474cae933744",
            "excerpt": "tip 5a1c3c3b8e1ffbbd04ff12c85e8d0dca56a6734a\nbench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/probe\nexit 0, wall 16241 ms\n  github.com/gibbonmi/bench/internal/probe,pass,16241,100\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t5-probe",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        },
        {
          "id": "t4-testreport-r1",
          "performer": "claude:ft290_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6d30a8f09034e3ce40838dd97f09d8f380ddd8cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t4_r1",
            "digest": "sha256:36d06ca8cc4d0a47b4370cd8224d7b590e639c5b938b8f4a2938abc1c03894ee",
            "excerpt": "HEAD 71582c05f4028c1ed6a643f9747895251aed31a3\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\n  github.com/gibbonmi/bench/internal/testreport,pass,32447,158\nfailures[0] skips[0]\n"
          },
          "requirement": "t4-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t4-prose-r1",
          "performer": "claude:ft290_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6d30a8f09034e3ce40838dd97f09d8f380ddd8cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t4_r1",
            "digest": "sha256:ec772514bcd17336da01d9321a60b0fa8a7c82794672d3412742af08b51270cd",
            "excerpt": "HEAD 71582c05f4028c1ed6a643f9747895251aed31a3\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/prose\n  github.com/gibbonmi/bench/internal/prose,pass,350,79\nfailures[0] skips[0]\n$ bench probe internal/prose/walk.go --omit '<if g.ex.excluded(rel) { continue } block>' --package ./internal/prose --run '^TestGradeReportsGradedSubjects$'\nprobe: bit,internal/prose/walk.go,omit,failed,1,yes\n  mutated run: github.com/gibbonmi/bench/internal/prose,fail,4,1 (exit 1)\n  TestGradeReportsGradedSubjects walk_test.go:302: GradeTree() = {Subjects:[drafts/skip.md keep.md skip.md] Findings:[]}, want the one subject keep.md and no finding\n"
          },
          "requirement": "t4-prose",
          "command": "bench test --package ./internal/prose",
          "exit_code": 0,
          "probe": {
            "mutation": "Omission at the subject answer in internal/prose/walk.go: omit the exclusion-set test, so that an excluded file counts as a subject. TestGradeReportsGradedSubjects must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t4_r1",
              "digest": "sha256:ec772514bcd17336da01d9321a60b0fa8a7c82794672d3412742af08b51270cd",
              "excerpt": "HEAD 71582c05f4028c1ed6a643f9747895251aed31a3\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/prose\n  github.com/gibbonmi/bench/internal/prose,pass,350,79\nfailures[0] skips[0]\n$ bench probe internal/prose/walk.go --omit '<if g.ex.excluded(rel) { continue } block>' --package ./internal/prose --run '^TestGradeReportsGradedSubjects$'\nprobe: bit,internal/prose/walk.go,omit,failed,1,yes\n  mutated run: github.com/gibbonmi/bench/internal/prose,fail,4,1 (exit 1)\n  TestGradeReportsGradedSubjects walk_test.go:302: GradeTree() = {Subjects:[drafts/skip.md keep.md skip.md] Findings:[]}, want the one subject keep.md and no finding\n"
            }
          }
        },
        {
          "id": "t4-prose-mechanics-r1",
          "performer": "claude:ft290_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6d30a8f09034e3ce40838dd97f09d8f380ddd8cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t4_r1",
            "digest": "sha256:af40ffd58130abe63c49b0d1f81562b1cbf96591359d24b38b25edd335663bd3",
            "excerpt": "HEAD 71582c05f4028c1ed6a643f9747895251aed31a3 (dist/bench rebuilt at this tip by coordinator)\n$ bench worktree exec \"ft290-test-projection\" -- ./dist/bench test --check prose-mechanics\ncheck: prose-mechanics,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,146,1\nfailures[0] skips[0]\n"
          },
          "requirement": "t4-prose-mechanics",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        },
        {
          "id": "t5-testreport-r1",
          "performer": "claude:ft290_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6d30a8f09034e3ce40838dd97f09d8f380ddd8cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t5_r1",
            "digest": "sha256:6acee6a7737c939fb3245eb008f400c89a7572a25da250940eb2e140cc761b7a",
            "excerpt": "HEAD 71582c05f4028c1ed6a643f9747895251aed31a3\n$ bench test --package ./internal/testreport\n  github.com/gibbonmi/bench/internal/testreport,pass,34901,158\nexit 0\n$ bench probe internal/testreport/outcome.go --swap \"failed++\" --with \"failed += len(f.rows(true))\" --package ./internal/testreport --run '^TestFullFailedTestsCountsTests$'\n  bit,internal/testreport/outcome.go,swap,failed,1,yes\n  mutated run: testreport,fail,44,1; TestFullFailedTestsCountsTests FailedTests:3 want 1 (exit 1)\n"
          },
          "requirement": "t5-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the failed-test count: count the --full failures rows in place of the distinct failed tests for Outcome.FailedTests. TestFullFailedTestsCountsTests must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t5_r1",
              "digest": "sha256:6acee6a7737c939fb3245eb008f400c89a7572a25da250940eb2e140cc761b7a",
              "excerpt": "HEAD 71582c05f4028c1ed6a643f9747895251aed31a3\n$ bench test --package ./internal/testreport\n  github.com/gibbonmi/bench/internal/testreport,pass,34901,158\nexit 0\n$ bench probe internal/testreport/outcome.go --swap \"failed++\" --with \"failed += len(f.rows(true))\" --package ./internal/testreport --run '^TestFullFailedTestsCountsTests$'\n  bit,internal/testreport/outcome.go,swap,failed,1,yes\n  mutated run: testreport,fail,44,1; TestFullFailedTestsCountsTests FailedTests:3 want 1 (exit 1)\n"
            }
          }
        },
        {
          "id": "t5-probe-r1",
          "performer": "claude:ft290_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6d30a8f09034e3ce40838dd97f09d8f380ddd8cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t5_r1",
            "digest": "sha256:a91d1f80c1cd8d726070326edf99c0f6c3ac01973960f8cf5e2a0bdb2cad712a",
            "excerpt": "HEAD 71582c05f4028c1ed6a643f9747895251aed31a3\n$ bench test --package ./internal/probe\n  github.com/gibbonmi/bench/internal/probe,pass,18170,100\nexit 0\n"
          },
          "requirement": "t5-probe",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "c1b-standards-r1",
          "performer": "claude:ft290_c1b_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9dda6bd8e26b08d40413e6ff3299191f6bb9f9ee",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c1b_standards",
            "digest": "sha256:629df5879c1b98b2d29e345bd6035629552371d989803acf73b465286665f5de",
            "excerpt": "axis=Standards findings=1 worst=S1 (TP-C1b, d49b0697..5a1c3c3b)\nS1 | auto-fix | conf 5 | internal/testreport/outcome.go:12, 24-25 | craft-comments Aging: the OutcomeFailed doc comment still says \"at least one failing test row\", but a --full row is now one diagnostic line and FailedTests counts distinct tests; the edited Outcome comment is one over-long line. | Reword line 12 to \"at least one failed test\" and re-wrap lines 24-27.\nRefuted duplication candidates: GradeTree reuses collect and the one exclusion owner; ranNothingTitle is the single title owner; report.failed() feeds both the rows and outcome(); no new canned-events helper; TP10/TP11 read the owner's text from prose.Grade with preconditions and exact whole-output matches.\nAdvice: a Grader method returning findings and graded could avoid a double exclusion call; failed()/failures() names are close.\nExamined: full diff, tickets 4-5, walk.go, subject.go, named_check.go, outcome.go, testreport.go, test helpers, AGENTS.md, BENCH.md, benchkit profile, craft-comments. No tests run.\n"
          },
          "axis": "Standards",
          "base": "d49b069704efe603a84203bdeffa9614c4802c37",
          "tip": "5a1c3c3b8e1ffbbd04ff12c85e8d0dca56a6734a",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "c1b-spec-r1",
          "performer": "claude:ft290_c1b_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9dda6bd8e26b08d40413e6ff3299191f6bb9f9ee",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c1b_spec",
            "digest": "sha256:8a6b12dfc322fcaf778a97dda5c65b59b11371cfb0292045b75b799d9cf71f79",
            "excerpt": "axis=Spec findings=0 worst=none (TP-C1b, d49b0697..5a1c3c3b)\nAll twelve rows closed: TP7 TestProseGreenPrintsOnlyCheckRow, TP8 TestProseFullListsSubjects (walk order conflicts with sorted order), TP9 TestProseZeroSubjectsExitsOne, TP10 TestProseRedKeepsFindingsAfterCheckRow, TP11 TestProseGraderRefusalPrintsCheckRow, TP12 TestGradeReportsGradedSubjects, TP13-TP17 exact-match failures tests, TP49 TestFullFailedTestsCountsTests.\nprose.Grade = GradeTree(root).Findings keeps its zero-subject pass; its conformance caller is untouched; one walk and one exclusion test.\nTestFullFailureDiagnostics keeps order, no-ANSI intent, default preview, and FailedTests 3. No scope creep.\nOrchestrator flag: a grader refusal keeping OutcomeFailed contradicts no spec line; probe refuses prose, so no reader sees the kind.\nAdvice: no direct test of the tracked-only subject list through GradeTree; collect's own tests cover the walk.\nExamined: spec, tickets 4-5, the chunk diff, walk.go, full_failure_test.go, subject.go. No tests run.\n"
          },
          "axis": "Spec",
          "base": "d49b069704efe603a84203bdeffa9614c4802c37",
          "tip": "5a1c3c3b8e1ffbbd04ff12c85e8d0dca56a6734a",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1b-coverage-r1",
          "performer": "claude:ft290_c1b_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9dda6bd8e26b08d40413e6ff3299191f6bb9f9ee",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c1b_coverage",
            "digest": "sha256:e194843cc24cc9b51e31e1921f3b63ce33abe1b1fd746c3df5fb5c8f0dabff2b",
            "excerpt": "axis=Coverage findings=2 worst=C1 (TP-C1b, d49b0697..5a1c3c3b)\nC1 | auto-fix | conf 9 | internal/testreport/named_check.go:305; TP10, story 8 | No test runs a red prose tree with --full. Swapping the findings branch to `len(grade.Findings) > 0 && !full` probed silent: a red --check prose --full run drops every finding and exits 0. TP9 and TP11 are default-only too. | Add --full cases to the red, refusal, and zero-subject prose tests.\nC2 | auto-fix | conf 8 | internal/prose/walk.go:60, walk_test.go:292; TP12 | TestGradeReportsGradedSubjects uses only a file-row exclusion. Swapping g.ex.excluded(rel) for g.ex.files[rel] probed silent in prose and testreport; this repo's .bench/prose-exclusions uses directory rows. | Add a directory-prefix exclusion row and a file under it to the fixture.\nCovered: the blank-line skip in decode and the zero-subject early return both bit.\nAdvice: a control byte in a .md path makes toon.Table refuse only under --full (unexercised); GradeTree evaluates the exclusion twice.\nExamined: chunk diff, spec rows TP7-TP17 and TP49, edge inventory, walk.go, subject.go, exclusions.go, testreport.go, named_check.go; bench test ./internal/probe passed; ./dist/bench --check prose (prose,prose,0,365) and prose-mechanics passed; 4 probes, all restored yes; git status clean.\n"
          },
          "axis": "Coverage",
          "base": "d49b069704efe603a84203bdeffa9614c4802c37",
          "tip": "5a1c3c3b8e1ffbbd04ff12c85e8d0dca56a6734a",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "c1b-standards-c1",
          "performer": "claude:ft290_c1b_standards_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6d30a8f09034e3ce40838dd97f09d8f380ddd8cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c1b_standards_c1",
            "digest": "sha256:62bc13f20068effecaa40c64f0228068257045b779c5e486f63188edf6693c34",
            "excerpt": "axis=Standards findings=0 worst=none (confirming round, repair delta 5a1c3c3b..71582c05)\nS1 confirmed: outcome.go:12 names failed tests; the Outcome comment is wrapped and states FailedTests counts distinct failed tests; rg over internal/testreport comments finds none that treats a failures row as a failed test.\nNo duplicated knowledge: proseFailureModes is one source for the two mode argument lists; want is built once per test; the directory exclusion fixture row adds no second parser.\nVenue note: the one git diff of the repair delta ran from the primary checkout, which shares the object store; file reads used the worktree.\nExamined: --check-current at f616ed37, repair delta, outcome.go, named_check.go:88-117, prose/exclusions.go. No tests run.\n"
          },
          "axis": "Standards",
          "base": "d49b069704efe603a84203bdeffa9614c4802c37",
          "tip": "71582c05f4028c1ed6a643f9747895251aed31a3",
          "finding_ids": [],
          "supersedes": [
            "c1b-standards-r1"
          ]
        },
        {
          "id": "c1b-spec-c1",
          "performer": "claude:ft290_c1b_spec_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6d30a8f09034e3ce40838dd97f09d8f380ddd8cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c1b_spec_c1",
            "digest": "sha256:4757498635fa503f7c6591462237595900e8396cb85be282f40fc2fc0fcb28a9",
            "excerpt": "axis=Spec findings=0 worst=none (confirming round, repair delta 5a1c3c3b..71582c05)\nNo production behavior change: two test files plus comment-only text in outcome.go.\nAll twelve TP-C1b rows keep their named tests; TP9, TP10, TP11 now loop over both modes with unchanged assertions; TP12 adds a directory exclusion row.\nrunProseCheck returns the findings and zero-subject branches before it reads full, so the --full cases match spec \"The prose result\" (lines 149-152); the directory row matches exclusions.go:16 and spec lines 142-143.\nExamined: --check-current, repair delta, spec sections, named_check.go:70-131, exclusions.go, rg for test names. No tests run.\n"
          },
          "axis": "Spec",
          "base": "d49b069704efe603a84203bdeffa9614c4802c37",
          "tip": "71582c05f4028c1ed6a643f9747895251aed31a3",
          "finding_ids": [],
          "supersedes": [
            "c1b-spec-r1"
          ]
        },
        {
          "id": "c1b-coverage-c1",
          "performer": "claude:ft290_c1b_coverage_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6d30a8f09034e3ce40838dd97f09d8f380ddd8cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c1b_coverage_c1",
            "digest": "sha256:f76ede4d6e043f327012e1704c1d1065065a52a8d7c0af880ad0047790e58158",
            "excerpt": "axis=Coverage findings=0 worst=none (confirming round, repair delta 5a1c3c3b..71582c05)\nC1 replay (&& !full at named_check.go:99, --package ./internal/testreport): bit, 2 failures from the --full mode loop.\nC2 replay (g.ex.files[rel] at walk.go:60): bit in ./internal/prose (TestGradeReportsGradedSubjects); silent in ./internal/testreport, which owns no directory-row row (TP12 belongs to internal/prose).\nNew bypass: a --full-only suffix after the findings output bit; a --full-only suffix after the zero-subject error was silent because TP9 matches by prefix, and no row binds it.\nAdvice: make the zero-subject expectation an exact-bytes match in both modes.\nExamined: --check-current at f616ed37, repair delta, named_check.go, walk.go, prose_check_test.go, spec rows TP8-TP12; 5 probes, all restored yes; git status clean.\n"
          },
          "axis": "Coverage",
          "base": "d49b069704efe603a84203bdeffa9614c4802c37",
          "tip": "71582c05f4028c1ed6a643f9747895251aed31a3",
          "finding_ids": [],
          "supersedes": [
            "c1b-coverage-r1"
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
      "from": "sha256:37523705f05ff5539429fd7c190772de817b553748203d244e9ad3900e2c0a57",
      "to": "sha256:e7275892d011a725f88fd6f4f231059449f393bf693fbd5723403f2c7ded0c8e",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ]
      }
    },
    {
      "from": "sha256:e7275892d011a725f88fd6f4f231059449f393bf693fbd5723403f2c7ded0c8e",
      "to": "sha256:1a858540618fdb0d00d294ecff2d46ad3e367de73fba5e3107d2943004030bac",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ]
      }
    },
    {
      "from": "sha256:1a858540618fdb0d00d294ecff2d46ad3e367de73fba5e3107d2943004030bac",
      "to": "sha256:0a578218a0d886a9acb001ee1ccf0168a84abd0b6589b1d3cd4de21a1337e123",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ],
        "TP-C1b": [
          "TP-C1b"
        ]
      }
    }
  ]
}
```
