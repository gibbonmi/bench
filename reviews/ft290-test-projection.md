# Review outcomes

## TP-C3 review pickup

The TP-C3 review ran on the frozen pair `21da7e2e..23c328db`. Each axis ran on sonnet at high effort. The chunk has used 1 of its 2 repair cycles. The repair sessions `claude:ft290_t8_r1` and `claude:ft290_t9_r1` repaired the six targets at `4f78c7b3` and `f04f5dd6`, and the chunk tip is now `f04f5dd6`.

The confirming round passed on Standards and Spec. The Coverage axis found C4 in the repair delta, so the chunk starts its second and last repair cycle.

- C4, auto-fix, confidence 7. `rootRecords` in `internal/canary/root.go` can return a real inventory error as an empty answer, and no test reaches that path. Story 24 says a wrong inventory never reads as empty. Ticket 8 owns the fix: add a test with one fixture name in two families that expects exit 1 and the diagnostic.

The repair session `claude:ft290_t8_r2` repaired C4 at `8159a136`, and the chunk tip is now `8159a136`. The chunk has used 2 of its 2 repair cycles. The C4 swap now bites in `internal/testreport`, which serves the user-facing output, and it stays silent in `internal/canary`.

The second confirming round passed on Standards and Spec. The Coverage axis found C5.

- C5, auto-fix, confidence 8. `FixturePins` also calls `rootRecords`, and the build preflight reads its error. The C4 swap stays silent in `internal/canary` and `internal/preflight`, so a real inventory error there reads as an empty pin map. Ticket 8 owns the fix: add a canary test that expects `FixturePins` to return the error for one fixture name in two families.

The reviewer stated that `--auto-approve` applies to repair rounds. Under that explicit statement, the orchestrator extends the TP-C3 allowance by one repair cycle for C5 only. The extension adds one repair round to the implementation retro.

The raw count is 6 findings, and the repair-target count is 6. Ticket 8 owns S1, C1, and the `fixtures_face_test.go` half of S2. Ticket 9 owns P1, C2, C3, and the `checks_face_test.go` half of S2. P1 and C2 name one fix.

The orchestrator decided the three `ask-user` findings under the reviewer's auto-approval for spec, ticket, and repair expansions. Each decision stays open to reviewer veto.

### Standards

Count: 2. Worst issue: S1.

- S1, ask-user, accepted, confidence 5. The face adds a fifth join of the `tests/canary` root and a second rule that maps the no-fixtures error to an empty answer. Put one accessor in a new `internal/canary/root.go`, and let every caller use it. The ticket 8 fence grows by that file.
- S2, auto-fix, confidence 4. The face test comments name omissions whose reds the record does not hold. Record a probe red for each named omission, or narrow each comment to the recorded reds.

### Spec

Count: 1. Worst issue: P1.

- P1, ask-user, accepted, confidence 4. `--checks` exits 1 on an invalid canary inventory, and no row grades it. Keep the exit 1, because story 24 says a wrong inventory never reads as empty. The new row TP56 grades it.

### Coverage

Count: 3. Worst issue: C1.

- C1, auto-fix, confidence 8. The sort test plants two fixtures of one family, so a missing sort or a sort by name can stay green. Plant enough fixtures across two families that name order and path order differ.
- C2, auto-fix, confidence 7. No test runs `--checks` over an invalid inventory. Add the TP56 test.
- C3, ask-user, accepted, confidence 4. No test refuses `--checks` with `--changed`, `--package`, `--base`, or `--run`. Add those refusals to the grammar refusal test.

## TP-C2 review pickup

The TP-C2 review ran on the frozen pair `71582c05..5fd1de93`. Each axis ran on sonnet at high effort. The chunk has used 1 of its 2 repair cycles. The repair sessions `claude:ft290_t6_r1` and `claude:ft290_t7_r1` repaired the four targets at `f4a5a9c3` and `21da7e2e`, and the chunk tip is now `21da7e2e`. The ticket 6 repair took one fence expansion to `cmd/bench/test_command.go`. The confirming round of all three axes passed with zero findings, and each earlier silent mutation now bites.

The confirming axes gave advice with no finding ID. `testHelpSuffix` has no comment on its trimmed prefix. No test pins the `usage: ` prefix of the help text. The repair author's excerpt holds the post-repair red of the usage literal.

The raw count is 5 findings, and the repair-target count is 4, because S1 and C1 name one fix. Ticket 6 owns S1, S2, and C1. Ticket 7 owns S3 and C2.

The orchestrator decided the two `ask-user` findings under the reviewer's auto-approval for spec and ticket expansions. Each decision stays open to reviewer veto.

### Standards

Count: 3. Worst issue: S1.

- S1, auto-fix, confidence 7. The `Cmd` and `Help` fields in `command.go` hold the same usage line, typed twice. Keep one usage constant, and derive `Help` from it.
- S2, ask-user, accepted, confidence 5. The `Suffix` of the `test` help row types the same grammar a third time. Derive it from the usage constant, so tickets 8 and 9 edit one owner. The rendered bytes do not change.
- S3, auto-fix, confidence 5. A comment in `unknown_check_test.go` cites the spec as provenance. State what the rows grade.

### Spec

Count: 0. Worst issue: none. All nine TP-C2 rows are closed.

### Coverage

Count: 2. Worst issue: C1.

- C1, auto-fix, confidence 6. No test pins the usage text, and a swap that drops the system run form stayed green. Add a test that pins the form in a usage refusal and in the help text.
- C2, ask-user, accepted, confidence 5. No test reaches the branch where the running executable cannot be named. The spec now lists `executable: unknown` as a flagged addition. Add a test for that branch.

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
  "plan_digest": "sha256:ec1780551fcfb36e85fe605183ad34a3c1b45f058f5c0c2e6eda382c605655a5",
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
    },
    {
      "id": "TP-C2",
      "base": "71582c05f4028c1ed6a643f9747895251aed31a3",
      "tip": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
      "plan_digest": "sha256:90908c6e7bac9a90b47002f5c1f5c78d32305faf06f2717a13a33b682a7389cd",
      "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
      "acceptance_rows": [
        "TP19",
        "TP20",
        "TP21",
        "TP22",
        "TP23",
        "TP24",
        "TP25",
        "TP26",
        "TP54"
      ],
      "verification": [
        {
          "id": "t6-testreport-v1",
          "performer": "claude:ft290_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0c638fcc90b2ee1f59c92e3d99a1a27ee44ad04f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t6",
            "digest": "sha256:a1845c9a0031a70227c6c9759c6052e3a2a68c0b1563ed51c68ee452b6411012",
            "excerpt": "tree 5fd1de93ac74d8d76d3562590fdec6f1c5e87466\n$ bench test --package ./internal/testreport  (exit 0)\n  github.com/gibbonmi/bench/internal/testreport,pass,37270,167\nfailures[0]\n$ bench probe internal/testreport/selection_facts.go --swap 'case r.focused.check == proseCheckName:' --with 'case r.focused.check == proseCheckName || r.focused.check == gate.SystemPhaseName:' --package ./internal/testreport --run '^TestSystemRequestRunFact$'\n  bit,internal/testreport/selection_facts.go,swap,failed,1,yes  (mutated run exit 1)\n  TestSystemRequestRunFact: Run() = \"all\", want \"^TestX$\"\n"
          },
          "requirement": "t6-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at Request.Run(): return AllTests for the system request with a pattern, as the joined prose and system case did. TestSystemRequestRunFact must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t6",
              "digest": "sha256:a1845c9a0031a70227c6c9759c6052e3a2a68c0b1563ed51c68ee452b6411012",
              "excerpt": "tree 5fd1de93ac74d8d76d3562590fdec6f1c5e87466\n$ bench test --package ./internal/testreport  (exit 0)\n  github.com/gibbonmi/bench/internal/testreport,pass,37270,167\nfailures[0]\n$ bench probe internal/testreport/selection_facts.go --swap 'case r.focused.check == proseCheckName:' --with 'case r.focused.check == proseCheckName || r.focused.check == gate.SystemPhaseName:' --package ./internal/testreport --run '^TestSystemRequestRunFact$'\n  bit,internal/testreport/selection_facts.go,swap,failed,1,yes  (mutated run exit 1)\n  TestSystemRequestRunFact: Run() = \"all\", want \"^TestX$\"\n"
            }
          }
        },
        {
          "id": "t6-cmd-v1",
          "performer": "claude:ft290_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0c638fcc90b2ee1f59c92e3d99a1a27ee44ad04f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t6",
            "digest": "sha256:d51df7627413955022228544fc26ec9cc0811062dc7b43ccb0b54a52b8763fd8",
            "excerpt": "tree 5fd1de93ac74d8d76d3562590fdec6f1c5e87466\n$ bench test --package ./cmd/bench  (exit 0)\n  github.com/gibbonmi/bench/cmd/bench,pass,17891,330\nfailures[0]\n"
          },
          "requirement": "t6-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t6-axi-query-registry-v1",
          "performer": "claude:ft290_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0c638fcc90b2ee1f59c92e3d99a1a27ee44ad04f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t6",
            "digest": "sha256:c64060a5492d29a7fd76ce7845d39a3227ea48e82b18bab5dc246c2372316d37",
            "excerpt": "tree 5fd1de93ac74d8d76d3562590fdec6f1c5e87466\n$ ./dist/bench test --check axi-query-registry  (exit 0)\n  axi-query-registry,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,6,1\n"
          },
          "requirement": "t6-axi-query-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t6-subcommand-routing-v1",
          "performer": "claude:ft290_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0c638fcc90b2ee1f59c92e3d99a1a27ee44ad04f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t6",
            "digest": "sha256:ae9812f38eea6254569d2f35e905585e8122e3815ae51bda0eafbe47d300c067",
            "excerpt": "tree 5fd1de93ac74d8d76d3562590fdec6f1c5e87466\n$ ./dist/bench test --check subcommand-routing  (exit 0)\n  subcommand-routing,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,25,1\n"
          },
          "requirement": "t6-subcommand-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t6-package-core-guard-v1",
          "performer": "claude:ft290_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0c638fcc90b2ee1f59c92e3d99a1a27ee44ad04f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t6",
            "digest": "sha256:f7a5a55bd944b398a52aad870e3ac7326ad46bf65b9ff8cd158b0a4c649933ce",
            "excerpt": "tree 5fd1de93ac74d8d76d3562590fdec6f1c5e87466\n$ ./dist/bench test --check package-core-guard  (exit 0)\n  package-core-guard,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,2363,1\n"
          },
          "requirement": "t6-package-core-guard",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t7-testreport-v1",
          "performer": "claude:ft290_t7",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0c638fcc90b2ee1f59c92e3d99a1a27ee44ad04f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t7",
            "digest": "sha256:eff700094503bbe1a9f27c69e1e2da040809703cd9d0130a46aef807f1f7198c",
            "excerpt": "HEAD 5fd1de93ac74d8d76d3562590fdec6f1c5e87466 (only reviews/ft290-test-projection.md modified, outside the package)\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\nexit 0\n  github.com/gibbonmi/bench/internal/testreport,pass,36461,167\nfailures[0] skips[0]\n\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/testreport/named_check.go --swap 'seal = sources' --with 'seal = \"unsealed\"; _ = sources' --package ./internal/testreport --run '^TestUnknownCheckNamesSealSources$'\nprobe exit 0\n  bit,internal/testreport/named_check.go,swap,failed,1,yes\nmutated run: internal/testreport,fail,250,1 (exit 1)\n  TestUnknownCheckNamesSealSources: \"...seal: unsealed...\"\n"
          },
          "requirement": "t7-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the seal line of the unknown-check refusal: replace the source digest from freshness.SealDigests with the literal unsealed. TestUnknownCheckNamesSealSources must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t7",
              "digest": "sha256:eff700094503bbe1a9f27c69e1e2da040809703cd9d0130a46aef807f1f7198c",
              "excerpt": "HEAD 5fd1de93ac74d8d76d3562590fdec6f1c5e87466 (only reviews/ft290-test-projection.md modified, outside the package)\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\nexit 0\n  github.com/gibbonmi/bench/internal/testreport,pass,36461,167\nfailures[0] skips[0]\n\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/testreport/named_check.go --swap 'seal = sources' --with 'seal = \"unsealed\"; _ = sources' --package ./internal/testreport --run '^TestUnknownCheckNamesSealSources$'\nprobe exit 0\n  bit,internal/testreport/named_check.go,swap,failed,1,yes\nmutated run: internal/testreport,fail,250,1 (exit 1)\n  TestUnknownCheckNamesSealSources: \"...seal: unsealed...\"\n"
            }
          }
        },
        {
          "id": "t7-probe-v1",
          "performer": "claude:ft290_t7",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0c638fcc90b2ee1f59c92e3d99a1a27ee44ad04f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t7",
            "digest": "sha256:b744bdf9c9a15e5d2cfd73e0160097a00344d332e17abaa14bc64fb6a739c2e3",
            "excerpt": "HEAD 5fd1de93ac74d8d76d3562590fdec6f1c5e87466 (only reviews/ft290-test-projection.md modified)\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/probe\nexit 0\n  github.com/gibbonmi/bench/internal/probe,pass,20300,100\nfailures[0] skips[0]\n"
          },
          "requirement": "t7-probe",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        },
        {
          "id": "t6-testreport-r1",
          "performer": "claude:ft290_t6_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t6_r1",
            "digest": "sha256:38e8be03b39546028523ad02da6caf59361dc3029a55515793b9462140079d1b",
            "excerpt": "HEAD 21da7e2e89fe107910def9543c3af0149fe9f2e0\nbench test --package ./internal/testreport\n  github.com/gibbonmi/bench/internal/testreport,pass,34030,169; failures[0]; exit 0\nbench probe internal/testreport/selection_facts.go --swap 'case r.focused.check == proseCheckName:' --with 'case r.focused.check == proseCheckName || r.focused.check == gate.SystemPhaseName:' --package ./internal/testreport --run '^TestSystemRequestRunFact$'\n  bit, failed_tests 1, restored yes; TestSystemRequestRunFact: Run() = \"all\", want \"^TestX$\"; mutated run exit 1\n"
          },
          "requirement": "t6-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at Request.Run(): return AllTests for the system request with a pattern, as the joined prose and system case did. TestSystemRequestRunFact must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t6_r1",
              "digest": "sha256:38e8be03b39546028523ad02da6caf59361dc3029a55515793b9462140079d1b",
              "excerpt": "HEAD 21da7e2e89fe107910def9543c3af0149fe9f2e0\nbench test --package ./internal/testreport\n  github.com/gibbonmi/bench/internal/testreport,pass,34030,169; failures[0]; exit 0\nbench probe internal/testreport/selection_facts.go --swap 'case r.focused.check == proseCheckName:' --with 'case r.focused.check == proseCheckName || r.focused.check == gate.SystemPhaseName:' --package ./internal/testreport --run '^TestSystemRequestRunFact$'\n  bit, failed_tests 1, restored yes; TestSystemRequestRunFact: Run() = \"all\", want \"^TestX$\"; mutated run exit 1\n"
            }
          }
        },
        {
          "id": "t6-cmd-r1",
          "performer": "claude:ft290_t6_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t6_r1",
            "digest": "sha256:739660dced819d76d2238bd6996cfa09cdf7fea9f7f2efa02650118b1eb6dbd2",
            "excerpt": "HEAD 21da7e2e89fe107910def9543c3af0149fe9f2e0\nbench test --package ./cmd/bench\n  github.com/gibbonmi/bench/cmd/bench,pass,16046,330; failures[0]; exit 0\n"
          },
          "requirement": "t6-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t6-axi-query-registry-r1",
          "performer": "claude:ft290_t6_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t6_r1",
            "digest": "sha256:ecd77a15b8222a5e7293bc75348c6e82733d6570c57671e806105b3575cbb3b6",
            "excerpt": "HEAD 21da7e2e89fe107910def9543c3af0149fe9f2e0\n./dist/bench test --check axi-query-registry\n  check axi-query-registry,conformance,tests_run 1; conformance pass 6ms; failures[0]; exit 0\n"
          },
          "requirement": "t6-axi-query-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t6-subcommand-routing-r1",
          "performer": "claude:ft290_t6_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t6_r1",
            "digest": "sha256:bc61cb8351a9ea392f751ba25f99c1c3c33bfa3b39da02ae70b9294d564897cb",
            "excerpt": "HEAD 21da7e2e89fe107910def9543c3af0149fe9f2e0\n./dist/bench test --check subcommand-routing\n  check subcommand-routing,conformance,tests_run 1; conformance pass 24ms; failures[0]; exit 0\n"
          },
          "requirement": "t6-subcommand-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t6-package-core-guard-r1",
          "performer": "claude:ft290_t6_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t6_r1",
            "digest": "sha256:5795a7a5eca4fb8db3fd4b70d54fc59ad8f47ed98b6c712527f58aceef0e71a8",
            "excerpt": "HEAD 21da7e2e89fe107910def9543c3af0149fe9f2e0\n./dist/bench test --check package-core-guard\n  check package-core-guard,conformance,tests_run 1; conformance pass 2347ms; failures[0]; exit 0\n"
          },
          "requirement": "t6-package-core-guard",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t7-testreport-r1",
          "performer": "claude:ft290_t7_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t7_r1",
            "digest": "sha256:cb5e07bcf9e06128ec40e419d9383b7602c89de4a76593f884e1ca287c54d178",
            "excerpt": "HEAD 21da7e2e89fe107910def9543c3af0149fe9f2e0 (only reviews/ft290-test-projection.md modified, other sessions)\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\ngithub.com/gibbonmi/bench/internal/testreport,pass,34234,169\nfailures[0] skips[0]\n\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/testreport/named_check.go --swap 'seal = sources' --with 'seal = \"unsealed\"; _ = sources' --package ./internal/testreport --run '^TestUnknownCheckNamesSealSources$'\nprobe: bit,internal/testreport/named_check.go,swap,failed,1,yes\npackage: internal/testreport,fail,168,1\nfailure: TestUnknownCheckNamesSealSources unknown_check_test.go:85 \"...\\nseal: unsealed\\n...\"\n"
          },
          "requirement": "t7-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the seal line of the unknown-check refusal: replace the source digest from freshness.SealDigests with the literal unsealed. TestUnknownCheckNamesSealSources must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t7_r1",
              "digest": "sha256:cb5e07bcf9e06128ec40e419d9383b7602c89de4a76593f884e1ca287c54d178",
              "excerpt": "HEAD 21da7e2e89fe107910def9543c3af0149fe9f2e0 (only reviews/ft290-test-projection.md modified, other sessions)\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/testreport\ngithub.com/gibbonmi/bench/internal/testreport,pass,34234,169\nfailures[0] skips[0]\n\n$ bench worktree exec \"ft290-test-projection\" -- bench probe internal/testreport/named_check.go --swap 'seal = sources' --with 'seal = \"unsealed\"; _ = sources' --package ./internal/testreport --run '^TestUnknownCheckNamesSealSources$'\nprobe: bit,internal/testreport/named_check.go,swap,failed,1,yes\npackage: internal/testreport,fail,168,1\nfailure: TestUnknownCheckNamesSealSources unknown_check_test.go:85 \"...\\nseal: unsealed\\n...\"\n"
            }
          }
        },
        {
          "id": "t7-probe-r1",
          "performer": "claude:ft290_t7_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t7_r1",
            "digest": "sha256:e4c50c15c914b7fa65ea5023c6a383f2fc252cadbbd88904bdd3b93dc2863350",
            "excerpt": "HEAD 21da7e2e89fe107910def9543c3af0149fe9f2e0 (only reviews/ft290-test-projection.md modified, other sessions)\n$ bench worktree exec \"ft290-test-projection\" -- bench test --package ./internal/probe\ngithub.com/gibbonmi/bench/internal/probe,pass,17785,100\nfailures[0] skips[0]\n"
          },
          "requirement": "t7-probe",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "c2-standards-r1",
          "performer": "claude:ft290_c2_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "0c638fcc90b2ee1f59c92e3d99a1a27ee44ad04f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c2_standards",
            "digest": "sha256:d2c8d4143d439d1d445a87dfdb5b635fddf76ef5d46ea7632aea9e3d5a2e1f43",
            "excerpt": "axis=Standards findings=3 worst=S1 (TP-C2, 71582c05..5fd1de93)\nS1 | auto-fix | conf 7 | internal/testreport/command.go:24-25 | Cmd and Help hold the same usage line typed twice; no test pins either field's text, so dropping the new form from one passes. AGENTS.md one source per fact. | Declare one usage-line const; Help = \"usage: \" + const.\nS2 | ask-user | conf 5 | cmd/bench/main.go:106 vs internal/testreport/command.go:24 | The test help row Suffix retypes the same grammar minus the \"bench test\" prefix; nothing derives or cross-checks it. cmd/bench already imports testreport. | Reviewer decides whether Suffix derives from the exported usage const.\nS3 | auto-fix | conf 5 | internal/testreport/unknown_check_test.go:255-258 | \"the spec fixes the line order...\" is provenance under craft-comments. | State what the rows grade and drop \"the spec fixes\".\nAdvice: the system-name tests answer different questions; the seal read has one caller; the \"-run\", pattern pairs feed one argv builder; \"running executable\" is used correctly; the two cmd/bench pins overlap by spec mandate.\nExamined: --check-current, full chunk diff, tickets 6-7, spec 240-275 and 515-535, CONTEXT.md, craft-comments, smell baseline, freshness.SealDigests, usage.Grammar. No tests run.\n"
          },
          "axis": "Standards",
          "base": "71582c05f4028c1ed6a643f9747895251aed31a3",
          "tip": "5fd1de93ac74d8d76d3562590fdec6f1c5e87466",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "c2-spec-r1",
          "performer": "claude:ft290_c2_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "0c638fcc90b2ee1f59c92e3d99a1a27ee44ad04f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c2_spec",
            "digest": "sha256:1341a0455668aaa93fa7321c70b2bb873fa9e5a6045c63e85f45ce74e7b9e8e0",
            "excerpt": "axis=Spec findings=0 worst=none (TP-C2, 71582c05..5fd1de93)\nAll nine rows closed: TP19 TestSystemRunPatternReachesGoArgv (-run after the suite operands); TP20 TestSystemRequestRunFact via Prepare (prose keeps AllTests); TP21 TestSystemRunPatternNoMatchRefusalWins (refusal before the zero-rule branch); TP22 TestNamedCheckRefusalMatrix; TP23 TestProseRefusesRunPattern; TP24 TestUnknownNamedCheckReportsOperandAndInventory in unknown_check_test.go; TP25, TP26, TP54 in unknown_check_test.go.\nRefusal order, the seal sources value or unsealed, and sanitize.Controls on the name and the path match \"The refusal identity\". The grammar equals the spec text minus the ticket 8 and 9 forms; --in appears only in the rendered pins.\nFlag 1: executable: unknown on an os.Executable failure contradicts no spec line; no row binds it.\nFlag 2: the TP24 seam edit names the file that holds the test.\nAdvice: list the executable: unknown literal under Flagged additions; the TP20 case sits in system_run_test.go, not beside AllTests in selection_facts_test.go (non-behavioral).\nExamined: coverage projection, the chunk diff, spec sections, tickets 6-7, command.go, named_check.go, selection_facts.go, freshness.go, sanitize.Controls. No tests run.\n"
          },
          "axis": "Spec",
          "base": "71582c05f4028c1ed6a643f9747895251aed31a3",
          "tip": "5fd1de93ac74d8d76d3562590fdec6f1c5e87466",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c2-coverage-r1",
          "performer": "claude:ft290_c2_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "0c638fcc90b2ee1f59c92e3d99a1a27ee44ad04f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c2_coverage",
            "digest": "sha256:a81810f58230e3dd9723a54ec7a5f78cf9723c7c77b24bb9cef163d2978a5026",
            "excerpt": "axis=Coverage findings=2 worst=C1 (TP-C2, 71582c05..5fd1de93)\nC1 | auto-fix | conf 6 | internal/testreport/command.go:24-25; spec \"The grammar\" 248-251 | No test pins the Cmd or Help text: dropping the system run form from Cmd probed silent in ./internal/testreport and ./cmd/bench, and from Help silent in ./internal/testreport. Only the Suffix is pinned. | Add a testreport test whose independent expectation holds the system run form in a usage refusal and in the help text.\nC2 | ask-user | conf 5 | internal/testreport/named_check.go:31; spec \"The refusal identity\" 181-186 | The os.Executable failure branch prints executable: unknown; no test reaches it (swapping \"unknown\" for \"\" probed silent), and the spec names no fallback. | Reviewer decides the fallback word; then add a test with a failing runningExecutable asserting exit 2 and the whole refusal.\nBit (restored yes): --run allowed for line-routing; --run allowed for prose; the no-match pass-through omitted; the != system guard in Request.Run omitted; ^pattern$ rewrite; sources swapped for the executable digest; sanitize.Controls dropped on the path.\nLive: bench test ./cmd/bench passed (330 tests); ./dist/bench test --check system --run '^TestNoSuchSystemTest$' exit 1 with the no-runs refusal and no zero-rule title.\nAdvice: an invalid regex exits 1 with Go's text; the seal-sources test cannot tell a read value from a recomputed digest; --in primary uses the older child.\nExamined: chunk diff, spec rows TP19-TP26 and TP54, edge inventory, grammar, tickets 6-7, freshness, sanitize, gate.SystemSuite; git status clean.\n"
          },
          "axis": "Coverage",
          "base": "71582c05f4028c1ed6a643f9747895251aed31a3",
          "tip": "5fd1de93ac74d8d76d3562590fdec6f1c5e87466",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "c2-standards-c1",
          "performer": "claude:ft290_c2_standards_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c2_standards_c1",
            "digest": "sha256:d5c3bfc7e197b9c90d114cd763d6f5ff77f523530e4de720a077341c24785a44",
            "excerpt": "axis=Standards findings=0 worst=none (confirming round, repair delta 5fd1de93..21da7e2e)\nS1 confirmed: command.go:23 holds one exported Usage; Cmd and Help derive from it.\nS2 confirmed: the test help row reads testHelpSuffix = strings.TrimPrefix(testreport.Usage, \"bench test\"); the removed literal matches, so rendered bytes are unchanged.\nS3 confirmed: the unknown_check_test.go header states what the rows grade, with no spec provenance.\nDuplication: the usage_text_test.go literal is an independent expectation under the AGENTS.md exception; its post-repair red is in the repair author's excerpt, which this axis did not read (conf 3, below the bar). checkInventory reads production namedChecks; answerRunningExecutable is the one swap owner.\nAdvice: testHelpSuffix has no comment explaining the trimmed prefix.\nExamined: repair delta via worktree exec, --check-current, named_check.go, rg sweeps, the review record, spec flagged additions. No tests run.\n"
          },
          "axis": "Standards",
          "base": "71582c05f4028c1ed6a643f9747895251aed31a3",
          "tip": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "finding_ids": [],
          "supersedes": [
            "c2-standards-r1"
          ]
        },
        {
          "id": "c2-spec-c1",
          "performer": "claude:ft290_c2_spec_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c2_spec_c1",
            "digest": "sha256:fa8f09bcf12840725351973b1bbe8613e4cafc63501fabf07636d54d2098f189",
            "excerpt": "axis=Spec findings=0 worst=none (confirming round, repair delta 5fd1de93..21da7e2e)\nRendered bytes unchanged: Usage holds the old string; testHelpSuffix equals the old Suffix; Help equals the old Help; unknownCheck is outside the delta.\nAll nine TP-C2 rows keep their named tests; the unknown_check_test.go change is a helper extraction with no weaker assertion.\nUsage equals spec.md:251 minus the ticket 8 and 9 forms; --in appears only in the rendered pin.\nTestUnknownCheckNamesUnknownExecutable matches the Flagged additions line and \"The refusal identity\".\nThe fence expansion appears in ticket 6 Writes and the spec fence list; the usage constant is a flagged addition.\nAdvice: tickets 8 and 9 still say add the form to the grammar and the help row; one Usage edit now feeds both.\nExamined: --check-current at f3783290, repair delta, spec sections, tickets 6-9, named_check.go, command_registry_test.go:715-720. No tests run.\n"
          },
          "axis": "Spec",
          "base": "71582c05f4028c1ed6a643f9747895251aed31a3",
          "tip": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "finding_ids": [],
          "supersedes": [
            "c2-spec-r1"
          ]
        },
        {
          "id": "c2-coverage-c1",
          "performer": "claude:ft290_c2_coverage_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "31efa65fe4d819a5f626b2f25f6fbad216d83e3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c2_coverage_c1",
            "digest": "sha256:b55eb07249374b30ebb7a43ade1551f718324f28f9ca48da4f4d4e38ffb88c69",
            "excerpt": "axis=Coverage findings=0 worst=none (confirming round, repair delta 5fd1de93..21da7e2e)\nReplays, no --run, all restored=yes:\nC1 drop the system form from Usage: bit in ./internal/testreport (TestUsageTextNamesSystemRunForm) and in ./cmd/bench (TestHelpInventoryIsComplete, TestTestHelpNamesOnlyRunnableFocusedForms).\nC1b Help-only drop: bit (TestUsageTextNamesSystemRunForm, --help case).\nC2 \"unknown\" swapped for \"\": bit (TestUnknownCheckNamesUnknownExecutable).\nNew bypass: a testHelpSuffix prefix mismatch bit both cmd/bench pins. A double space after \"usage:\" in Help was silent; cosmetic, no row binds it.\nExamined: --check-current at f3783290, repair delta, command.go, named_check.go, the cmd/bench pins; git status clean after every probe.\n"
          },
          "axis": "Coverage",
          "base": "71582c05f4028c1ed6a643f9747895251aed31a3",
          "tip": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "finding_ids": [],
          "supersedes": [
            "c2-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "TP-C3",
      "base": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
      "tip": "8159a136023072cf42f4265f4b74cd932e539ac7",
      "plan_digest": "sha256:ec1780551fcfb36e85fe605183ad34a3c1b45f058f5c0c2e6eda382c605655a5",
      "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
      "acceptance_rows": [
        "TP27",
        "TP28",
        "TP29",
        "TP30",
        "TP33",
        "TP40",
        "TP52",
        "TP53",
        "TP31",
        "TP32",
        "TP34",
        "TP35",
        "TP36",
        "TP37",
        "TP38",
        "TP39",
        "TP50",
        "TP56"
      ],
      "verification": [
        {
          "id": "t8-testreport-v1",
          "performer": "claude:ft290_t8",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8",
            "digest": "sha256:7d8737696abed3d4a4fb92e59122ab603d8c8f27b32f80eceaa9708f68ba870a",
            "excerpt": "tip 23c328db1eec38cc3031c0b9dd7a4157d547c74f\n$ bench test --package ./internal/testreport\n  github.com/gibbonmi/bench/internal/testreport,pass,34292,186  failures[0] skips[0] exit 0\n$ bench probe internal/testreport/fixtures.go --swap \"if fixture.Check == check {\" --with \"if fixture.Family == check {\" --package ./internal/testreport --run ^TestFixturesFaceHonorsCheckMarker$\n  bit,internal/testreport/fixtures.go,swap,failed,1,yes  (mutated run: testreport fail,5,1; exit 1)\n"
          },
          "requirement": "t8-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the fixtures owner filter: compare the fixture family name in place of the owner that canary.Fixtures resolves. TestFixturesFaceHonorsCheckMarker must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t8",
              "digest": "sha256:7d8737696abed3d4a4fb92e59122ab603d8c8f27b32f80eceaa9708f68ba870a",
              "excerpt": "tip 23c328db1eec38cc3031c0b9dd7a4157d547c74f\n$ bench test --package ./internal/testreport\n  github.com/gibbonmi/bench/internal/testreport,pass,34292,186  failures[0] skips[0] exit 0\n$ bench probe internal/testreport/fixtures.go --swap \"if fixture.Check == check {\" --with \"if fixture.Family == check {\" --package ./internal/testreport --run ^TestFixturesFaceHonorsCheckMarker$\n  bit,internal/testreport/fixtures.go,swap,failed,1,yes  (mutated run: testreport fail,5,1; exit 1)\n"
            }
          }
        },
        {
          "id": "t8-canary-v1",
          "performer": "claude:ft290_t8",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8",
            "digest": "sha256:f2dfa4bf49d0db84e87355e639c00ba154e84cbc3869dd158cdc5e301bfbc3d4",
            "excerpt": "tip 23c328db1eec38cc3031c0b9dd7a4157d547c74f\n$ bench test --package ./internal/canary\n  github.com/gibbonmi/bench/internal/canary,pass,16,26  failures[0] exit 0\n"
          },
          "requirement": "t8-canary",
          "command": "bench test --package ./internal/canary",
          "exit_code": 0
        },
        {
          "id": "t8-cmd-v1",
          "performer": "claude:ft290_t8",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8",
            "digest": "sha256:bf83bcdee76b0f24ef8f37e311f9b497ffc1e9da44a68842fcf5d0d720807e5c",
            "excerpt": "tip 23c328db1eec38cc3031c0b9dd7a4157d547c74f\n$ bench test --package ./cmd/bench\n  github.com/gibbonmi/bench/cmd/bench,pass,14784,330  failures[0] exit 0\n"
          },
          "requirement": "t8-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t8-ordinary-build-census-v1",
          "performer": "claude:ft290_t8",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8",
            "digest": "sha256:55cb8b9afa7b4f2770aa6918ba5b0db56debb3d56e3aeccebff0d0f62b7a3c15",
            "excerpt": "tip 23c328db1eec38cc3031c0b9dd7a4157d547c74f\n$ ./dist/bench test --check ordinary-build-census\n  check: ordinary-build-census,conformance,1,0\n  internal/conformance,pass,310,1  exit 0\n"
          },
          "requirement": "t8-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t8-axi-query-registry-v1",
          "performer": "claude:ft290_t8",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8",
            "digest": "sha256:dc6a37578a48a1acb63bf4363720baac2acd51da90dbdf2496d1ebb57fc0d095",
            "excerpt": "tip 23c328db1eec38cc3031c0b9dd7a4157d547c74f\n$ ./dist/bench test --check axi-query-registry\n  check: axi-query-registry,conformance,1,0\n  internal/conformance,pass,5,1  exit 0\n"
          },
          "requirement": "t8-axi-query-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t8-subcommand-routing-v1",
          "performer": "claude:ft290_t8",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8",
            "digest": "sha256:49545ca1d4cf962618aee045b549aa41c71914a418586aba732c6714d2c7d13d",
            "excerpt": "tip 23c328db1eec38cc3031c0b9dd7a4157d547c74f\n$ ./dist/bench test --check subcommand-routing\n  check: subcommand-routing,conformance,1,0\n  internal/conformance,pass,20,1  exit 0\n"
          },
          "requirement": "t8-subcommand-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t8-package-core-guard-v1",
          "performer": "claude:ft290_t8",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8",
            "digest": "sha256:b7795db06425342069bc1ad1667e60d261f66224760c88a6cebf5e982ac56991",
            "excerpt": "tip 23c328db1eec38cc3031c0b9dd7a4157d547c74f\n$ ./dist/bench test --check package-core-guard\n  check: package-core-guard,conformance,1,0\n  internal/conformance,pass,2379,1  exit 0\n"
          },
          "requirement": "t8-package-core-guard",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t9-testreport-v1",
          "performer": "claude:ft290_t9",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9",
            "digest": "sha256:316a20c98a608a03f7e169f94f4d47f6ba23d7bb4454095d8d17ab34403fe45c",
            "excerpt": "@23c328db (only reviews/ft290-test-projection.md modified)\n$ bench test --package ./internal/testreport\n  github.com/gibbonmi/bench/internal/testreport,pass,34741,186\nexit 0\n$ bench probe internal/testreport/fixtures.go --swap 'namedCheckKind(check), len(families)}' --with 'namedCheckKind(check), len(owned)}' --package ./internal/testreport --run '^TestChecksFaceCountsFamilies$'\n  bit,internal/testreport/fixtures.go,swap,failed,1,yes\n  TestChecksFaceCountsFamilies: package-core-guard families = \"3\"; want 2 (mutated run exit 1)\n"
          },
          "requirement": "t9-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the families cell: count the owned fixtures in place of the distinct non-empty family names. TestChecksFaceCountsFamilies must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t9",
              "digest": "sha256:316a20c98a608a03f7e169f94f4d47f6ba23d7bb4454095d8d17ab34403fe45c",
              "excerpt": "@23c328db (only reviews/ft290-test-projection.md modified)\n$ bench test --package ./internal/testreport\n  github.com/gibbonmi/bench/internal/testreport,pass,34741,186\nexit 0\n$ bench probe internal/testreport/fixtures.go --swap 'namedCheckKind(check), len(families)}' --with 'namedCheckKind(check), len(owned)}' --package ./internal/testreport --run '^TestChecksFaceCountsFamilies$'\n  bit,internal/testreport/fixtures.go,swap,failed,1,yes\n  TestChecksFaceCountsFamilies: package-core-guard families = \"3\"; want 2 (mutated run exit 1)\n"
            }
          }
        },
        {
          "id": "t9-cmd-v1",
          "performer": "claude:ft290_t9",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9",
            "digest": "sha256:16b60758a48580892e0b7d48eec7cc4048045ea2a636033a552b7e60c1c8d51b",
            "excerpt": "@23c328db\n$ bench test --package ./cmd/bench\n  github.com/gibbonmi/bench/cmd/bench,pass,14384,330\nexit 0\n"
          },
          "requirement": "t9-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t9-axi-query-registry-v1",
          "performer": "claude:ft290_t9",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9",
            "digest": "sha256:b2a04e42f3c28aa1a6b99eccab3c98129cc0b2c01e9836ca060fed483b2fe762",
            "excerpt": "@23c328db, dist/bench rebuilt there\n$ ./dist/bench test --check axi-query-registry\n  axi-query-registry,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,7,1\nexit 0\n"
          },
          "requirement": "t9-axi-query-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t9-subcommand-routing-v1",
          "performer": "claude:ft290_t9",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9",
            "digest": "sha256:1c3a36a08809fcc4c85c42a0a29da9991a4635a0f9dc022708f25f8a0143a860",
            "excerpt": "@23c328db, dist/bench rebuilt there\n$ ./dist/bench test --check subcommand-routing\n  subcommand-routing,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,21,1\nexit 0\n"
          },
          "requirement": "t9-subcommand-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t9-package-core-guard-v1",
          "performer": "claude:ft290_t9",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9",
            "digest": "sha256:0299cc0f766501595a9dd88015cbdbce2e3f30dc5e88d23570bfc83d97d67b18",
            "excerpt": "@23c328db, dist/bench rebuilt there\n$ ./dist/bench test --check package-core-guard\n  package-core-guard,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,2385,1\nexit 0\n"
          },
          "requirement": "t9-package-core-guard",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t8-testreport-r1",
          "performer": "claude:ft290_t8_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r1",
            "digest": "sha256:5eaf7bdec3f807d1aab3040f0e942ead26d1d40573ff6e6d7b28b894b5543169",
            "excerpt": "$ bench test --package ./internal/testreport\ntree: ft290-test-projection,f04f5dd6fad737947218a5332874f40a43bc4bcc\ngithub.com/gibbonmi/bench/internal/testreport,pass,38309,187\nfailures[0] skips[0]\nprobe: bench probe internal/testreport/fixtures.go --swap \"fixture.Check == check\" --with \"fixture.Family == check\" --package ./internal/testreport --run '^TestFixturesFaceHonorsCheckMarker$'\nprobe[1]: bit,internal/testreport/fixtures.go,swap,failed,1,restored yes\nmutated run: testreport,fail,6,1 TestFixturesFaceHonorsCheckMarker fixtures_face_test.go:85\n"
          },
          "requirement": "t8-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the fixtures owner filter: compare the fixture family name in place of the owner that canary.Fixtures resolves. TestFixturesFaceHonorsCheckMarker must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t8_r1",
              "digest": "sha256:5eaf7bdec3f807d1aab3040f0e942ead26d1d40573ff6e6d7b28b894b5543169",
              "excerpt": "$ bench test --package ./internal/testreport\ntree: ft290-test-projection,f04f5dd6fad737947218a5332874f40a43bc4bcc\ngithub.com/gibbonmi/bench/internal/testreport,pass,38309,187\nfailures[0] skips[0]\nprobe: bench probe internal/testreport/fixtures.go --swap \"fixture.Check == check\" --with \"fixture.Family == check\" --package ./internal/testreport --run '^TestFixturesFaceHonorsCheckMarker$'\nprobe[1]: bit,internal/testreport/fixtures.go,swap,failed,1,restored yes\nmutated run: testreport,fail,6,1 TestFixturesFaceHonorsCheckMarker fixtures_face_test.go:85\n"
            }
          }
        },
        {
          "id": "t8-canary-r1",
          "performer": "claude:ft290_t8_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r1",
            "digest": "sha256:af41e2f9bbdd5a18159c38cfee1595076797a758d511416eb6c482e2e339eb9f",
            "excerpt": "$ bench test --package ./internal/canary\ntree: ft290-test-projection,f04f5dd6fad737947218a5332874f40a43bc4bcc\ngithub.com/gibbonmi/bench/internal/canary,pass,18,26\nfailures[0] skips[0]\n"
          },
          "requirement": "t8-canary",
          "command": "bench test --package ./internal/canary",
          "exit_code": 0
        },
        {
          "id": "t8-cmd-r1",
          "performer": "claude:ft290_t8_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r1",
            "digest": "sha256:a08b07167129ecd79130e85e8fd9ea738fb8324245a8a5c1b104a26d5743ff60",
            "excerpt": "$ bench test --package ./cmd/bench\ntree: ft290-test-projection,f04f5dd6fad737947218a5332874f40a43bc4bcc\ngithub.com/gibbonmi/bench/cmd/bench,pass,19007,330\nfailures[0] skips[0]\n"
          },
          "requirement": "t8-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t8-ordinary-build-census-r1",
          "performer": "claude:ft290_t8_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r1",
            "digest": "sha256:272fdc99e060addec438d38a5ab9ebe57fcbbb6f7043d000341877e3d63510ba",
            "excerpt": "$ ./dist/bench test --check ordinary-build-census\ntree: ft290-test-projection,f04f5dd6fad737947218a5332874f40a43bc4bcc\ncheck: ordinary-build-census,conformance,tests_run 1\ninternal/conformance,pass,359,1\nfailures[0] skips[0]\n"
          },
          "requirement": "t8-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t8-axi-query-registry-r1",
          "performer": "claude:ft290_t8_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r1",
            "digest": "sha256:7bf051ee42e9e9794851c93307516fd530086d3b377c8242561d5b620c5eedc1",
            "excerpt": "$ ./dist/bench test --check axi-query-registry\ntree: ft290-test-projection,f04f5dd6fad737947218a5332874f40a43bc4bcc\ncheck: axi-query-registry,conformance,tests_run 1\ninternal/conformance,pass,7,1\nfailures[0] skips[0]\n"
          },
          "requirement": "t8-axi-query-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t8-subcommand-routing-r1",
          "performer": "claude:ft290_t8_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r1",
            "digest": "sha256:0d7944540b7546f954ec8425e9ee25a78e9abc05fe8792710ddfda915cb853bc",
            "excerpt": "$ ./dist/bench test --check subcommand-routing\ntree: ft290-test-projection,f04f5dd6fad737947218a5332874f40a43bc4bcc\ncheck: subcommand-routing,conformance,tests_run 1\ninternal/conformance,pass,25,1\nfailures[0] skips[0]\n"
          },
          "requirement": "t8-subcommand-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t8-package-core-guard-r1",
          "performer": "claude:ft290_t8_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r1",
            "digest": "sha256:5b165647ad90a85ca60815e16b98ecd3b77c915064b21865c5a80fd2b55c80a4",
            "excerpt": "$ ./dist/bench test --check package-core-guard\ntree: ft290-test-projection,f04f5dd6fad737947218a5332874f40a43bc4bcc\ncheck: package-core-guard,conformance,tests_run 1\ninternal/conformance,pass,2810,1\nfailures[0] skips[0]\n"
          },
          "requirement": "t8-package-core-guard",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t9-testreport-r1",
          "performer": "claude:ft290_t9_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9_r1",
            "digest": "sha256:98c1ee5475b13419809f4476d860c11235877ee5de79b2463dbad264146a5c58",
            "excerpt": "HEAD f04f5dd6 (only reviews/ft290-test-projection.md modified, not ours)\nbench test --package ./internal/testreport\ngithub.com/gibbonmi/bench/internal/testreport,pass,37086,187 ; failures[0] skips[0] ; exit 0\nbench probe internal/testreport/fixtures.go --swap 'namedCheckKind(check), len(families)}' --with 'namedCheckKind(check), len(owned)}' --package ./internal/testreport --run '^TestChecksFaceCountsFamilies$'\nprobe: bit,internal/testreport/fixtures.go,swap,failed,1,yes ; mutated run fail (exit 1): package-core-guard families = \"3\"; want 2\n"
          },
          "requirement": "t9-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the families cell: count the owned fixtures in place of the distinct non-empty family names. TestChecksFaceCountsFamilies must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t9_r1",
              "digest": "sha256:98c1ee5475b13419809f4476d860c11235877ee5de79b2463dbad264146a5c58",
              "excerpt": "HEAD f04f5dd6 (only reviews/ft290-test-projection.md modified, not ours)\nbench test --package ./internal/testreport\ngithub.com/gibbonmi/bench/internal/testreport,pass,37086,187 ; failures[0] skips[0] ; exit 0\nbench probe internal/testreport/fixtures.go --swap 'namedCheckKind(check), len(families)}' --with 'namedCheckKind(check), len(owned)}' --package ./internal/testreport --run '^TestChecksFaceCountsFamilies$'\nprobe: bit,internal/testreport/fixtures.go,swap,failed,1,yes ; mutated run fail (exit 1): package-core-guard families = \"3\"; want 2\n"
            }
          }
        },
        {
          "id": "t9-cmd-r1",
          "performer": "claude:ft290_t9_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9_r1",
            "digest": "sha256:669fa395111079c88eaeba34c38b1a13bb017d8cc617cd30ea96a6243809bf7d",
            "excerpt": "HEAD f04f5dd6\nbench test --package ./cmd/bench\ngithub.com/gibbonmi/bench/cmd/bench,pass,23744,330 ; failures[0] skips[0] ; exit 0\n"
          },
          "requirement": "t9-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t9-axi-query-registry-r1",
          "performer": "claude:ft290_t9_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9_r1",
            "digest": "sha256:76cc94db453e8522d8bf2db1ace27d8cf0aeefddd80cd7564e2222955a48f062",
            "excerpt": "HEAD f04f5dd6 ; ./dist/bench test --check axi-query-registry\ncheck: axi-query-registry,conformance,1,0 ; internal/conformance,pass,9,1 ; failures[0] ; exit 0\n"
          },
          "requirement": "t9-axi-query-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t9-subcommand-routing-r1",
          "performer": "claude:ft290_t9_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9_r1",
            "digest": "sha256:36ec217d2b9979e7a438674b9163bb3db60c9aa3c24d4b24752d51c73a94083f",
            "excerpt": "HEAD f04f5dd6 ; ./dist/bench test --check subcommand-routing\ncheck: subcommand-routing,conformance,1,0 ; internal/conformance,pass,22,1 ; failures[0] ; exit 0\n"
          },
          "requirement": "t9-subcommand-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t9-package-core-guard-r1",
          "performer": "claude:ft290_t9_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9_r1",
            "digest": "sha256:07ae0e960f2e101f521e8ba7f51292cc7dd5b3a88fc4af5d545bf11574f9811f",
            "excerpt": "HEAD f04f5dd6 ; ./dist/bench test --check package-core-guard\ncheck: package-core-guard,conformance,1,0 ; internal/conformance,pass,2917,1 ; failures[0] ; exit 0\n"
          },
          "requirement": "t9-package-core-guard",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t8-testreport-r2",
          "performer": "claude:ft290_t8_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r2",
            "digest": "sha256:0010e8474dbb7ddcb5054c805349dfa0599b2b4e3fd29f8e088d38f505c88b8c",
            "excerpt": "bench test --package ./internal/testreport\ntree: ft290-test-projection,8159a136023072cf42f4265f4b74cd932e539ac7\ngithub.com/gibbonmi/bench/internal/testreport,pass,34808,188\nfailures[0]\nprobe: bench probe internal/testreport/fixtures.go --swap \"fixture.Check == check\" --with \"fixture.Family == check\" --package ./internal/testreport --run ^TestFixturesFaceHonorsCheckMarker$\nprobe: bit,failed,failed_tests=1,restored=yes (mutated run fail, TestFixturesFaceHonorsCheckMarker)\n"
          },
          "requirement": "t8-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the fixtures owner filter: compare the fixture family name in place of the owner that canary.Fixtures resolves. TestFixturesFaceHonorsCheckMarker must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t8_r2",
              "digest": "sha256:0010e8474dbb7ddcb5054c805349dfa0599b2b4e3fd29f8e088d38f505c88b8c",
              "excerpt": "bench test --package ./internal/testreport\ntree: ft290-test-projection,8159a136023072cf42f4265f4b74cd932e539ac7\ngithub.com/gibbonmi/bench/internal/testreport,pass,34808,188\nfailures[0]\nprobe: bench probe internal/testreport/fixtures.go --swap \"fixture.Check == check\" --with \"fixture.Family == check\" --package ./internal/testreport --run ^TestFixturesFaceHonorsCheckMarker$\nprobe: bit,failed,failed_tests=1,restored=yes (mutated run fail, TestFixturesFaceHonorsCheckMarker)\n"
            }
          }
        },
        {
          "id": "t8-canary-r2",
          "performer": "claude:ft290_t8_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r2",
            "digest": "sha256:a9a172c769968ce663d593f38a5ae3d249e6fcd3032ff9b4cf202ce32a6e0ffb",
            "excerpt": "bench test --package ./internal/canary\ntree: ft290-test-projection,8159a136023072cf42f4265f4b74cd932e539ac7\ngithub.com/gibbonmi/bench/internal/canary,pass,16,26\nfailures[0]\n"
          },
          "requirement": "t8-canary",
          "command": "bench test --package ./internal/canary",
          "exit_code": 0
        },
        {
          "id": "t8-cmd-r2",
          "performer": "claude:ft290_t8_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r2",
            "digest": "sha256:8ed1a5075057d7b07678bcad217d436dd8150d033545205303e34fde16f10e72",
            "excerpt": "bench test --package ./cmd/bench\ntree: ft290-test-projection,8159a136023072cf42f4265f4b74cd932e539ac7\ngithub.com/gibbonmi/bench/cmd/bench,pass,14423,330\nfailures[0]\n"
          },
          "requirement": "t8-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t8-ordinary-build-census-r2",
          "performer": "claude:ft290_t8_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r2",
            "digest": "sha256:48a70f4a8f420986478deb3e4e680b96dd2ed1146fe1a410d9647dd2f09bdde8",
            "excerpt": "./dist/bench test --check ordinary-build-census\ntree: ft290-test-projection,8159a136023072cf42f4265f4b74cd932e539ac7\ncheck: ordinary-build-census,conformance,tests_run=1\ngithub.com/gibbonmi/bench/internal/conformance,pass,304,1\nfailures[0]\n"
          },
          "requirement": "t8-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t8-axi-query-registry-r2",
          "performer": "claude:ft290_t8_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r2",
            "digest": "sha256:190ec76d4a961c999063c104ba6ec64921b9814dfba9d008ef53770bd8ed3bdc",
            "excerpt": "./dist/bench test --check axi-query-registry\ntree: ft290-test-projection,8159a136023072cf42f4265f4b74cd932e539ac7\ncheck: axi-query-registry,conformance,tests_run=1\ngithub.com/gibbonmi/bench/internal/conformance,pass,6,1\nfailures[0]\n"
          },
          "requirement": "t8-axi-query-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t8-subcommand-routing-r2",
          "performer": "claude:ft290_t8_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r2",
            "digest": "sha256:0684fc31da02bddb38a89cb214fc2faddbcc2e00ba29a94de799a0de752c48e9",
            "excerpt": "./dist/bench test --check subcommand-routing\ntree: ft290-test-projection,8159a136023072cf42f4265f4b74cd932e539ac7\ncheck: subcommand-routing,conformance,tests_run=1\ngithub.com/gibbonmi/bench/internal/conformance,pass,21,1\nfailures[0]\n"
          },
          "requirement": "t8-subcommand-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t8-package-core-guard-r2",
          "performer": "claude:ft290_t8_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t8_r2",
            "digest": "sha256:42a338a2e9f6c4b6a74e4d32539e7f9b679cec342ee1f8270fbd4839a2eaec76",
            "excerpt": "./dist/bench test --check package-core-guard\ntree: ft290-test-projection,8159a136023072cf42f4265f4b74cd932e539ac7\ncheck: package-core-guard,conformance,tests_run=1\ngithub.com/gibbonmi/bench/internal/conformance,pass,2412,1\nfailures[0]\n"
          },
          "requirement": "t8-package-core-guard",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t9-testreport-r2",
          "performer": "claude:ft290_t9_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9_r1",
            "digest": "sha256:b64cb3accb60f0f8668eafd6fd357d0b9ba8056925ba7924c68eb6660d468569",
            "excerpt": "HEAD 8159a136 (only reviews/ft290-test-projection.md modified, not ours)\nbench test --package ./internal/testreport\ngithub.com/gibbonmi/bench/internal/testreport,pass,34801,188 ; failures[0] skips[0] ; exit 0\nbench probe internal/testreport/fixtures.go --swap 'namedCheckKind(check), len(families)}' --with 'namedCheckKind(check), len(owned)}' --package ./internal/testreport --run '^TestChecksFaceCountsFamilies$'\nprobe: bit,internal/testreport/fixtures.go,swap,failed,1,yes ; mutated run fail (exit 1): package-core-guard families = \"3\"; want 2\n"
          },
          "requirement": "t9-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the families cell: count the owned fixtures in place of the distinct non-empty family names. TestChecksFaceCountsFamilies must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft290_t9_r1",
              "digest": "sha256:b64cb3accb60f0f8668eafd6fd357d0b9ba8056925ba7924c68eb6660d468569",
              "excerpt": "HEAD 8159a136 (only reviews/ft290-test-projection.md modified, not ours)\nbench test --package ./internal/testreport\ngithub.com/gibbonmi/bench/internal/testreport,pass,34801,188 ; failures[0] skips[0] ; exit 0\nbench probe internal/testreport/fixtures.go --swap 'namedCheckKind(check), len(families)}' --with 'namedCheckKind(check), len(owned)}' --package ./internal/testreport --run '^TestChecksFaceCountsFamilies$'\nprobe: bit,internal/testreport/fixtures.go,swap,failed,1,yes ; mutated run fail (exit 1): package-core-guard families = \"3\"; want 2\n"
            }
          }
        },
        {
          "id": "t9-cmd-r2",
          "performer": "claude:ft290_t9_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9_r1",
            "digest": "sha256:7e0433847f06e0b0246745ea6faa2d23fcc0eca2bdac495997532aebe9cf8a71",
            "excerpt": "HEAD 8159a136\nbench test --package ./cmd/bench\ngithub.com/gibbonmi/bench/cmd/bench,pass,14444,330 ; failures[0] skips[0] ; exit 0\n"
          },
          "requirement": "t9-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t9-axi-query-registry-r2",
          "performer": "claude:ft290_t9_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9_r1",
            "digest": "sha256:5cfe30ec46a5bb0792ca0584d147cdad3865335dc677a0a1d19e79365aeec901",
            "excerpt": "HEAD 8159a136 ; ./dist/bench test --check axi-query-registry\ncheck: axi-query-registry,conformance,1,0 ; internal/conformance,pass,6,1 ; failures[0] ; exit 0\n"
          },
          "requirement": "t9-axi-query-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t9-subcommand-routing-r2",
          "performer": "claude:ft290_t9_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9_r1",
            "digest": "sha256:32bf8bbe39163dbe5a413093c2a7ac0472f07d4909a50dbb5ecd91a52794b204",
            "excerpt": "HEAD 8159a136 ; ./dist/bench test --check subcommand-routing\ncheck: subcommand-routing,conformance,1,0 ; internal/conformance,pass,21,1 ; failures[0] ; exit 0\n"
          },
          "requirement": "t9-subcommand-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t9-package-core-guard-r2",
          "performer": "claude:ft290_t9_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_t9_r1",
            "digest": "sha256:92a886ac02fad77b17b39fd765a9c1f134cdd218225520cc2b4df3a41b16cc40",
            "excerpt": "HEAD 8159a136 ; ./dist/bench test --check package-core-guard\ncheck: package-core-guard,conformance,1,0 ; internal/conformance,pass,2411,1 ; failures[0] ; exit 0\n"
          },
          "requirement": "t9-package-core-guard",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "c3-standards-r1",
          "performer": "claude:ft290_c3_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c3_standards",
            "digest": "sha256:b92ed1f766af96ac16ebe3d7e5979012c22e0e51aea38a63c1944ae5cfc14f6a",
            "excerpt": "axis=Standards findings=2 worst=S1 (TP-C3, 21da7e2e..23c328db)\nS1 | ask-user | conf 5 | internal/testreport/fixtures.go:19-20 vs internal/canary/inventory.go:46,56,135,160-161 | checkFixtures adds a fifth production join of root/tests/canary, the first outside internal/canary, and canaryInventory re-derives FixturePins' ErrNoFixtures-means-empty rule. AGENTS.md one source per fact. An accessor in a new internal/canary file fits without growing inventory.go. | Repair ticket 8 after a plan commit that adds the new canary file to its Writes, or the reviewer accepts the duplicate.\nS2 | auto-fix | conf 4 | fixtures_face_test.go:5-6; checks_face_test.go:8-10,123 | Header comments name omissions (absolute path, lost row, lost check, fixed kind, Usage edit) whose reds are not recorded; only the owner-filter and families-count probes are. | Record one probe red per named omission, or narrow the comments to the demonstrated reds.\nClean: one ownedFixtures filter; --checks reads namedChecks(); kind from namedCheckKind; exclusivity in parseFocusedRequest; TP34 parses --help, not the producer; no new dependency; no comment provenance.\nAdvice: decision.go:44 repeats the ErrNoFixtures text (outside the fence); TestInventoryFacesStartNoChild is near redundant; writeProseCheckFile name misleads.\nExamined: --check-current at ada702e2, the chunk diff, inventory.go, decision.go, named_check.go, command.go, spec 100-262 and 584-618, tickets 8-9, review record rows. No tests run.\n"
          },
          "axis": "Standards",
          "base": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "tip": "23c328db1eec38cc3031c0b9dd7a4157d547c74f",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c3-spec-r1",
          "performer": "claude:ft290_c3_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c3_spec",
            "digest": "sha256:ef4cb3bb92a6cb22f7b027472ea05c8034539528ef550eed0673a6793956f289",
            "excerpt": "axis=Spec findings=1 worst=P1 (TP-C3, 21da7e2e..23c328db)\nP1 | ask-user | conf 4 | internal/testreport/fixtures.go:345-347; spec \"The inventory faces\" 188-214, TP33, story 25 | --checks exits 1 with \"canary inventory invalid\" on a bad CHECK marker; no row grades it and story 25 says --checks runs at exit 0, so the spec reads two ways. | Reviewer decides: record exit 1 as intended with a test through ticket 9, or make --checks tolerate the error.\nAll 17 rows closed: TP27-TP30, TP33, TP52, TP53, TP40 in fixtures_face_test.go; TP31 via inventoryFace and TestInventoryFacesStartNoChild; TP32, TP34-TP38 in checks_face_test.go; TP39 TestInventoryGrammarRefusals; TP50 via the cmd/bench pins.\nUsage equals spec line 251; bench help renders it with --in after the first bench test only. ownedFixtures is the one owner; checksInventory reads namedChecks() and namedCheckKind. The sentinel, FixturePins errors.Is, the untouched Select literal, and the unchanged inventory.go line count match \"Fence disposition\".\nNotes: OutcomeNoTestRun on a listing contradicts no spec line, and no probe path reaches a face. Ticket 9's refactor of ticket 8 helpers is clean.\nAdvice: discoverFixtures maps any ReadDir failure to ErrNoFixtures (pre-existing); the inventory diagnostic prints raw paths; the two-row sort test catches a missing sort only half the time.\nExamined: the chunk diff, spec, tickets 8-9, command.go, named_check.go, outcome.go, fixtures.go, inventory.go, toon.go, probe consumers; ran ./dist/bench help, test --help, one --fixtures listing.\n"
          },
          "axis": "Spec",
          "base": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "tip": "23c328db1eec38cc3031c0b9dd7a4157d547c74f",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "c3-coverage-r1",
          "performer": "claude:ft290_c3_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7ba315fbd51fb0e5d41094f4e0546f51d9aee1d7",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c3_coverage",
            "digest": "sha256:6bf8916b589f26a9a4b4b76c74d07b34e1469aca8179bece9ccd93b95fc65961",
            "excerpt": "axis=Coverage findings=3 worst=C1 (TP-C3, 21da7e2e..23c328db)\nC1 | auto-fix | conf 8 | fixtures_face_test.go:41-52; spec TP27 \"sorted by path\" | Two same-family fixtures over a map: dropping the sort bit once and was silent once; swapping the sort key from path to name was silent. | Plant six or more fixtures across two families with one CHECK owner where name and path order differ.\nC2 | auto-fix | conf 7 | fixtures.go:343-347; spec 201-202, story 24 | No test runs --checks over an invalid inventory; swapping the refusal for found = nil was silent. | Add a --checks row over an unknown-CHECK fixture expecting exit 1 and \"names unknown check\".\nC3 | ask-user | conf 4 | command.go:72,75; TP39 | An allowlist bypass that lets --checks take --changed, --package, --base, or --run was silent. | Add those --checks refusals to TestInventoryGrammarRefusals.\nLive and correct: CHECK whitespace trimmed; empty CHECK exits 1; duplicate names exit 1; symlinks skipped; tab and newline escaped; U+007F raw; families counts.\nOpen (pre-existing): any ReadDir error maps to ErrNoFixtures, so a file or unreadable tests/canary reads as empty.\nExamined: --check-current at ada702e2, the chunk diff, spec, inventory.go, command.go, fixtures.go, named_check.go; bench test ./internal/canary passed; live --checks (45 rows) and --fixtures (27 rows); every probe restored; git status clean.\n"
          },
          "axis": "Coverage",
          "base": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "tip": "23c328db1eec38cc3031c0b9dd7a4157d547c74f",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "c3-standards-c1",
          "performer": "claude:ft290_c3_standards_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c3_standards_c1",
            "digest": "sha256:13dbe9bb981b7b2a8c7131df7926cd14be1407dcd367ae26b67d44122d69d076",
            "excerpt": "axis=Standards findings=0 worst=none (confirming round, repair delta 23c328db..f04f5dd6)\nS1 confirmed: root.go Dir is the only production tests/canary join; rootRecords is the one ErrNoFixtures-means-empty rule (callers RootFixtures and FixturePins); fixturesOf is the one record-to-Fixture body and Fixtures wraps it; Inventory, FixturePins, UnboundConformanceFamilies, and both faces call the accessor; inventory.go shrank 415 to 400; decision.go untouched.\nS2 confirmed: each omission named in the fixtures_face_test.go and checks_face_test.go headers has a recorded probe red in the t8r1 and t9r1 excerpts.\nComments: the new doc comments state behavior with no provenance.\nAdvice: decision.go:44 repeats the empty message (outside the fence); some tests outside the delta still join tests/canary by hand; the Fixtures doc comment could name ErrNoFixtures.\nExamined: --check-current, repair delta, inventory.go, the two face test headers, rg sweeps, spec flagged additions, review pickup, probe excerpts. No tests run.\n"
          },
          "axis": "Standards",
          "base": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "tip": "f04f5dd6fad737947218a5332874f40a43bc4bcc",
          "finding_ids": [],
          "supersedes": [
            "c3-standards-r1"
          ]
        },
        {
          "id": "c3-spec-c1",
          "performer": "claude:ft290_c3_spec_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c3_spec_c1",
            "digest": "sha256:5677e2c45d16d701050a042b36efeeaaddcbfc3a5787a5dbe8b4a367786d0ff0",
            "excerpt": "axis=Spec findings=0 worst=none (confirming round, repair delta 23c328db..f04f5dd6)\nNo rendered-output change: Dir equals the old join; Fixtures routes through the moved loop and still returns ErrNoFixtures; FixturePins still returns an empty map and skips CHECK validation; RootFixtures returns an empty map where canaryInventory returned nil.\nAll 18 TP-C3 rows have named tests, TP56 TestChecksFaceRefusesInvalidInventory included; no assertion is weaker.\nTP56 matches story 24 and shares the --fixtures refusal path. The four new --checks refusals match spec line 260. Every written path is inside the fence.\nAdvice: the TP27 row text still describes two fixtures a and b; the repaired test is a stricter ten-row exact match.\nExamined: --check-current at 8d294b25, the repair delta, the spec and ticket diff, inventory.go, root.go, command.go, fixtures_face_test.go, spec rows TP27-TP56 and the fence list. No tests run.\n"
          },
          "axis": "Spec",
          "base": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "tip": "f04f5dd6fad737947218a5332874f40a43bc4bcc",
          "finding_ids": [],
          "supersedes": [
            "c3-spec-r1"
          ]
        },
        {
          "id": "c3-coverage-c1",
          "performer": "claude:ft290_c3_coverage_c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "dc1ee0f29f592caa2cdc7b2333880b1c463e7a81",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c3_coverage_c1",
            "digest": "sha256:d15a7838d7b6f4d63f0d9630943ee5eb967e4b27b4addb08dfc8161e63d2b4ea",
            "excerpt": "axis=Coverage findings=1 worst=C1 (confirming round, repair delta 23c328db..f04f5dd6)\nC1 | auto-fix | conf 7 | internal/canary/root.go:94-100 (rootRecords); story 24 | Swapping `return records, err` for `return records, nil` in rootRecords was silent in ./internal/canary and ./internal/testreport: a real discoverFixtures error (a fixture name in two families, an unreadable family dir, a marker read error) reads as an empty inventory. The refusal tests plant only an unknown CHECK, which fails in fixturesOf. | Add a test that plants one fixture name in two families and expects exit 1 with the \"multiple families\" diagnostic, or a canary-level RootFixtures error test.\nReplays bit (restored yes): the no-op sort comparator (twice), the name sort, the --checks refusal swap, and the --changed, --package, and --run allowlist swaps; --base is an equivalent mutant. New: a wrong Dir join bit in both packages; a skipped fixtureCheck error bit in testreport.\nExamined: --check-current, the repair delta, fixtures.go, command.go, root.go, inventory.go, fixtures_face_test.go, spec story lines; bench test ./internal/canary passed; git status clean.\n"
          },
          "axis": "Coverage",
          "base": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "tip": "f04f5dd6fad737947218a5332874f40a43bc4bcc",
          "finding_ids": [
            "C4"
          ],
          "supersedes": [
            "c3-coverage-r1"
          ]
        },
        {
          "id": "c3-standards-c2",
          "performer": "claude:ft290_c3_standards_c2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c3_standards_c2",
            "digest": "sha256:317da44aa535ee5cd2cd2f1b1bdbd14b733ac0934d4ad86b983648d1c697a1df",
            "excerpt": "axis=Standards findings=0 worst=none (second confirming round, repair delta f04f5dd6..8159a136)\nThe delta adds only internal/testreport/fixtures_duplicate_test.go. It reuses plantFixture and inventoryFace, so no harness is duplicated.\nThe expected fragment \"twin\" appears in multiple families repeats the inline fmt.Errorf at inventory.go:244, which has no shared constant; the C4 swap red is recorded in ft290/t8r2/c4-red.txt.\nThe test comment states the why; naming and structure match the neighboring face tests; no dependency or production code.\nAdvice: a red for a wording-only change of the diagnostic is not recorded.\nExamined: --check-current, the delta, fixtures_face_test.go:1-70, inventory.go:244, rg for \"multiple families\", the c4-red excerpt. No tests run.\n"
          },
          "axis": "Standards",
          "base": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "tip": "8159a136023072cf42f4265f4b74cd932e539ac7",
          "finding_ids": [],
          "supersedes": [
            "c3-standards-c1"
          ]
        },
        {
          "id": "c3-spec-c2",
          "performer": "claude:ft290_c3_spec_c2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft290_c3_spec_c2",
            "digest": "sha256:75075ca39ac26501aa97dc7603c30c07a54ab410b3ebbdad628622db0757f294",
            "excerpt": "axis=Spec findings=0 worst=none (second confirming round, repair delta f04f5dd6..8159a136)\nThe delta is one new test file; no production file or existing test changes.\nThe test plants one fixture name in two families and requires exit 1 with the inventory.go:244 diagnostic from both --checks and --fixtures, which matches story 24 and \"The inventory faces\". rootRecords returns that error, because only ErrNoFixtures maps to empty.\nNo TP-C3 row lost its test; the file sits inside ticket 8's internal/testreport/ fence.\nBinding note: this axis did not rerun --check-current; the delta read used the frozen pair.\nExamined: the delta, inventory.go 78-260, root.go, fixtures.go, the face tests, toon.Errorf, spec line 83, ticket 8. No tests run.\n"
          },
          "axis": "Spec",
          "base": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "tip": "8159a136023072cf42f4265f4b74cd932e539ac7",
          "finding_ids": [],
          "supersedes": [
            "c3-spec-c1"
          ]
        },
        {
          "id": "c3-coverage-c2",
          "performer": "claude:ft290_c3_coverage_c2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "eabae1284bc39ceae2cacb4ff14effae6a95daec",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft290_c3_coverage_c2",
            "digest": "sha256:b863aeef4e75a259232a6350187da790291bf53bef40bba2a7812970213b9892",
            "excerpt": "axis=Coverage findings=1 worst=C5 (second confirming round, repair delta f04f5dd6..8159a136)\nC5 | auto-fix | conf 8 | internal/canary/root.go:26 (rootRecords); callers inventory.go:146 (FixturePins), internal/preflight/gather.go:297-299 | Swapping `return records, err` for `return records, nil` is silent in ./internal/canary (26 tests) and ./internal/preflight (353 tests): FixturePins would turn a real discoverFixtures error into an empty pin map, and the preflight would skip its \"fixture inventory not readable\" failure. No test asserts that path. | Add an internal/canary test that plants one fixture name in two families and asserts FixturePins returns the \"appears in multiple families\" error; re-probe the swap in ./internal/canary.\nC4 closed for the testreport faces: the swap bit in ./internal/testreport (TestInventoryFacesRefuseDuplicateFixtureName); a wrapped-error-as-empty bypass also bit there.\nCommand note: a repair of a shared helper should replay its mutation across every production caller.\nExamined: --check-current, the delta, root.go, inventory.go 70-300, gather.go 280-330, rg for callers and tests; four probes, all restored yes; git status clean.\n"
          },
          "axis": "Coverage",
          "base": "21da7e2e89fe107910def9543c3af0149fe9f2e0",
          "tip": "8159a136023072cf42f4265f4b74cd932e539ac7",
          "finding_ids": [
            "C5"
          ],
          "supersedes": [
            "c3-coverage-c1"
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
    },
    {
      "from": "sha256:0a578218a0d886a9acb001ee1ccf0168a84abd0b6589b1d3cd4de21a1337e123",
      "to": "sha256:ff60b6deddfa6e8eb884659bb0537e1ad180758e8bfb59dc2ad6e61907227d4c",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ],
        "TP-C1b": [
          "TP-C1b"
        ]
      }
    },
    {
      "from": "sha256:ff60b6deddfa6e8eb884659bb0537e1ad180758e8bfb59dc2ad6e61907227d4c",
      "to": "sha256:acc488347939b21c1b27b5c453cf5b81846267e1c8dd014737d8c5c7a200cc64",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ],
        "TP-C1b": [
          "TP-C1b"
        ],
        "TP-C2": [
          "TP-C2"
        ]
      }
    },
    {
      "from": "sha256:acc488347939b21c1b27b5c453cf5b81846267e1c8dd014737d8c5c7a200cc64",
      "to": "sha256:faab16952fedf136c9dd876f74ed853635eb338a234640dc5903d182f48f866d",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ],
        "TP-C1b": [
          "TP-C1b"
        ],
        "TP-C2": [
          "TP-C2"
        ]
      }
    },
    {
      "from": "sha256:faab16952fedf136c9dd876f74ed853635eb338a234640dc5903d182f48f866d",
      "to": "sha256:7eb209e7e947acc32004d2af598d3543421671facbc0e26ff0763c03db3e5f18",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ],
        "TP-C1b": [
          "TP-C1b"
        ],
        "TP-C2": [
          "TP-C2"
        ]
      }
    },
    {
      "from": "sha256:7eb209e7e947acc32004d2af598d3543421671facbc0e26ff0763c03db3e5f18",
      "to": "sha256:90908c6e7bac9a90b47002f5c1f5c78d32305faf06f2717a13a33b682a7389cd",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ],
        "TP-C1b": [
          "TP-C1b"
        ],
        "TP-C2": [
          "TP-C2"
        ]
      }
    },
    {
      "from": "sha256:90908c6e7bac9a90b47002f5c1f5c78d32305faf06f2717a13a33b682a7389cd",
      "to": "sha256:b88f0a385e878e95fc42dee1716deaa73c232bdf8e43d03326931ccdbcb8cf7d",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ],
        "TP-C1b": [
          "TP-C1b"
        ],
        "TP-C2": [
          "TP-C2"
        ]
      }
    },
    {
      "from": "sha256:b88f0a385e878e95fc42dee1716deaa73c232bdf8e43d03326931ccdbcb8cf7d",
      "to": "sha256:0cce2a8e053c994e4fadb3529c75d1deb4a1eb72722d3050e9669dee30444287",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ],
        "TP-C1b": [
          "TP-C1b"
        ],
        "TP-C2": [
          "TP-C2"
        ],
        "TP-C3": [
          "TP-C3"
        ]
      }
    },
    {
      "from": "sha256:0cce2a8e053c994e4fadb3529c75d1deb4a1eb72722d3050e9669dee30444287",
      "to": "sha256:ec1780551fcfb36e85fe605183ad34a3c1b45f058f5c0c2e6eda382c605655a5",
      "chunk_ids": {
        "TP-C1a": [
          "TP-C1a"
        ],
        "TP-C1b": [
          "TP-C1b"
        ],
        "TP-C2": [
          "TP-C2"
        ],
        "TP-C3": [
          "TP-C3"
        ]
      }
    }
  ]
}
```
