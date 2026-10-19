# Review outcomes

## CR1 review pickup

The CR1 review ran on the frozen pair `36b1fafc..1e8846bb`. Each axis ran on opus at high effort, by the reviewer's direction. The chunk has used 0 of its 2 repair cycles.

The raw count is 14 findings, and the repair-target count is 10. S2 and P2 name one fix, and P1 and C1 name one fix. S6, S7, and S8 have no repair target. Ticket 01 owns S1 and S2. Ticket 02 owns S3, S4, S5, P1, C2, and C3. The orchestrator owns P3 as a plan commit.

### Standards

Count: 8. Worst issue: S2.

- S1, auto-fix, confidence 9. `internal/preflight/ownership/premises_test.go:94`, `:162`, `:239`, and `:287` say the tests run through the proposal, but the bodies call `ownership.Collect`. Correct the four comments.
- S2, auto-fix, confidence 8. `CHANGELOG.md:25` says that a build charge refuses until each required current premise has proof, and the heading at `:22` names the removed closure. Correct the `[Unreleased]` entry.
- S3, auto-fix, confidence 7. The `extra` parameter of `seedChargeRepo` in `internal/preflight/ownership/charge_fixture_test.go:152` receives nil from every caller. Remove it.
- S4, auto-fix, confidence 6. `census.declaration` in `internal/preflight/ownership/helpers.go:62` returns a line that no caller reads. Remove the return value and correct its comment.
- S5, auto-fix, confidence 6. `CurrentComplete`, `RequiredCurrent`, and `PendingNative` in `internal/preflight/ownership/ownership.go:58-63` have no production reader after the `charge-premises` row went. Remove each field that no code reads.
- S6, no-op, confidence 6. The approved spec (story 47 and CR91) directs the deletion of the known-issues row. The reviewer can reopen the row as a new issue.
- S7, no-op, confidence 4. A rename of `internal/tickets/registry_data.go` is outside the acceptance.
- S8, no-op, confidence 3. The shadowed `source` in `internal/preflight/gather.go:105` is legacy code that the diff did not introduce.

### Spec

Count: 3. Worst issue: P2.

- P1, auto-fix, confidence 6. The CR06 assertion in `internal/preflight/charge_test.go:193` cannot fail, because a prepared charge prints no check rows. Add a `Decide` row-list test under the charge operations. C1 names the same fix.
- P2, auto-fix, confidence 8. S2 names the same fix.
- P3, auto-fix, confidence 6. Spec line 105 says the grader goes whole. Correct it to match ticket 02.

### Coverage

Count: 3. Worst issue: C1.

- C1, auto-fix, confidence 9. A kept green `charge-premises` row under `KindPrepareEvidence` or `KindCurrentEvidence` passes every test. Pin the exact row list of both charge operations in a `Decide` test.
- C2, auto-fix, confidence 6. `tickets.ParseWrite` admits `/`, `/ (new)`, `./`, and `.`. Refuse an empty path, and add those entries to `TestParseWriteRefusesAnUncleanPath`.
- C3, auto-fix, confidence 7. A blacklist body of `tickets.RepositoryPath` passes the clean-path tests. Add `internal/./x.go` as a refused entry and `a..b/x.go` as an accepted entry.

```bench-review-record
{
  "version": 2,
  "spec": "specs/closure-removal/spec.md",
  "plan_digest": "sha256:28c270ec1eed58e47fed9261ef9ff6d80486340f37c637856ac33111d7dcf7ec",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CR1",
      "base": "36b1fafc45426299b27b68f26fe627cadf000639",
      "tip": "1e8846bb094761993e922fd085648ab3de33344f",
      "plan_digest": "sha256:28c270ec1eed58e47fed9261ef9ff6d80486340f37c637856ac33111d7dcf7ec",
      "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
      "acceptance_rows": [
        "CR07",
        "CR08",
        "CR09",
        "CR10",
        "CR80",
        "CR81",
        "CR82",
        "CR83",
        "CR84",
        "CR85",
        "CR91",
        "CR93",
        "CR95",
        "CR01",
        "CR02",
        "CR03",
        "CR04",
        "CR05",
        "CR06",
        "CR11",
        "CR12",
        "CR13",
        "CR14",
        "CR15",
        "CR16",
        "CR17",
        "CR18",
        "CR19",
        "CR20",
        "CR21",
        "CR22",
        "CR89",
        "CR96"
      ],
      "verification": [
        {
          "id": "t2-t2-preflight-packages",
          "performer": "claude:ft293_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2",
            "digest": "sha256:1a574dd382e565a393855ea7c2db0f6f36294e9d0877937cd28f83370768d5c2",
            "excerpt": "packages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,24676,369\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,17720,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,9950,175\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,4973,69\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd",
          "exit_code": 0
        },
        {
          "id": "t2-t2-tickets-help-anchors",
          "performer": "claude:ft293_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2",
            "digest": "sha256:b9671242311550f46ab9094b2bc823718c88436008df1e1ea811240718665186",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/tickets/writes.go,omit,failed,1,TestParseWriteRefusesAnUncleanPath,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/tickets ./cmd/bench ./internal/anchors ./internal/canary ./internal/conformance,all,passed,1661\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,825\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15323,346\n  github.com/gibbonmi/bench/internal/anchors,pass,114,123\n  github.com/gibbonmi/bench/internal/canary,pass,7,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33471,1116\n  github.com/gibbonmi/bench/internal/tickets,fail,3,44\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/tickets,TestParseWriteRefusesAnUncleanPath,\"writes_test.go:50: ParseWrite(\\\"internal/x/../../cmd/y.go\\\") = {Path:internal/x/../../cmd/y.go Kind:0}, want an unclean-path refusal\",7\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-tickets-help-anchors",
          "command": "bench test --package ./internal/tickets --package ./cmd/bench --package ./internal/anchors --package ./internal/canary --package ./internal/conformance",
          "exit_code": 0,
          "probe": {
            "mutation": "Omission: remove the clean-path refusal from tickets.ParseWrite; TestParseWriteRefusesAnUncleanPath must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t2",
              "digest": "sha256:b9671242311550f46ab9094b2bc823718c88436008df1e1ea811240718665186",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/tickets/writes.go,omit,failed,1,TestParseWriteRefusesAnUncleanPath,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/tickets ./cmd/bench ./internal/anchors ./internal/canary ./internal/conformance,all,passed,1661\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,825\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15323,346\n  github.com/gibbonmi/bench/internal/anchors,pass,114,123\n  github.com/gibbonmi/bench/internal/canary,pass,7,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33471,1116\n  github.com/gibbonmi/bench/internal/tickets,fail,3,44\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/tickets,TestParseWriteRefusesAnUncleanPath,\"writes_test.go:50: ParseWrite(\\\"internal/x/../../cmd/y.go\\\") = {Path:internal/x/../../cmd/y.go Kind:0}, want an unclean-path refusal\",7\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t2-t2-root-conformance",
          "performer": "claude:ft293_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2",
            "digest": "sha256:cc99b42edfe38defb14f1f7e4624f572c34ef7eb1c468af37bbce9e03bea8dbe",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14606,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-t2-guidance-anchors",
          "performer": "claude:ft293_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2",
            "digest": "sha256:040fbe641b69c8364627390825e4760a73714fbc566f22ff602ea578485f2178",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4285,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t2-t2-guidance-budgets",
          "performer": "claude:ft293_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2",
            "digest": "sha256:7833b79d87da4342e711cb9d6c0cdd75b6c3afdfca28d197838c202bd9080c37",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t1-t1-preflight-packages",
          "performer": "claude:ft293_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1",
            "digest": "sha256:0278573194a4b093c118baa828234926bd62138cf537699ceaad8b061ff5aa37",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/evidencecmd/operations.go,swap+swap,failed,1,TestPreflightBuildRefusesRemovedForms,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight ./internal/preflight/ownership ./internal/preflight/ownership/ownershiptest ./internal/preflight/staleness ./internal/preflight/evidencecmd,all,passed,845\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,727\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,24629,369\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,17636,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,9642,175\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,4744,69\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestPreflightBuildRefusesRemovedForms,\"command_build_test.go:132: [build s --propose-writes --ticket a.md --base B --source-tip T] = (1), want exit 2 naming \\\"unknown argument: --propose-writes\\\":\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd",
          "exit_code": 0,
          "probe": {
            "mutation": "Omission at the removed-form refusal: restore the --propose-writes flag in the preflight build grammar so it parses again; TestPreflightBuildRefusesRemovedForms must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t1",
              "digest": "sha256:0278573194a4b093c118baa828234926bd62138cf537699ceaad8b061ff5aa37",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/evidencecmd/operations.go,swap+swap,failed,1,TestPreflightBuildRefusesRemovedForms,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight ./internal/preflight/ownership ./internal/preflight/ownership/ownershiptest ./internal/preflight/staleness ./internal/preflight/evidencecmd,all,passed,845\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,727\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,24629,369\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,17636,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,9642,175\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,4744,69\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestPreflightBuildRefusesRemovedForms,\"command_build_test.go:132: [build s --propose-writes --ticket a.md --base B --source-tip T] = (1), want exit 2 naming \\\"unknown argument: --propose-writes\\\":\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-t1-tickets-help-anchors",
          "performer": "claude:ft293_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1",
            "digest": "sha256:d7d62425aa2cddc18241dc327d07f6064ec4851eea77863b9ce87e2bfa77222b",
            "excerpt": "packages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14971,346\n  github.com/gibbonmi/bench/internal/anchors,pass,108,123\n  github.com/gibbonmi/bench/internal/canary,pass,7,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33111,1116\n  github.com/gibbonmi/bench/internal/tickets,pass,3,44\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-tickets-help-anchors",
          "command": "bench test --package ./internal/tickets --package ./cmd/bench --package ./internal/anchors --package ./internal/canary --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-t1-root-conformance",
          "performer": "claude:ft293_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1",
            "digest": "sha256:dc3213e6e7d411ee03d0a3c6a6fcdf718457d17ef3ae00a2dff0f03e72d0014f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13429,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-t1-guidance-anchors",
          "performer": "claude:ft293_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1",
            "digest": "sha256:12bb8449605e8dbc52ea5502bdc2f73d3718bac20a019e700c1b1f6fd13399d4",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4028,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t1-t1-guidance-budgets",
          "performer": "claude:ft293_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1",
            "digest": "sha256:7833b79d87da4342e711cb9d6c0cdd75b6c3afdfca28d197838c202bd9080c37",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "cr1-standards",
          "performer": "claude:ft293_cr1_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1_standards",
            "digest": "sha256:2f1ea0bafbf98dab4b32874d8128ef585cc76dff94976bf37533720fef6a5c31",
            "excerpt": "Standards axis, FT293 CR1 (base 36b1fafc4, tip 1e8846bb0). Evidence check-current passed.\nS1 auto-fix 9: premises_test.go:94,:162,:239,:287 doc comments still say \"through the actual proposal\"; bodies call premises(t, root) over ownership.Collect (craft-comments: never leave a comment describing the code that was).\nS2 auto-fix 8: CHANGELOG.md:25 still says \"A build charge refuses until its declared-complete inventory and each required current premise have proof.\"; heading :22 names a removed feature.\nS3 auto-fix 7: charge_fixture_test.go:152 seedChargeRepo extra map param; every caller passes nil (Speculative Generality).\nS4 auto-fix 6: ownership/helpers.go:62 census.declaration returns an unread line int; comment :59-61 promises it.\nS5 ask-user 6: ownership.go:58-63 CurrentComplete, RequiredCurrent, PendingNative computed in premises.go:301-318, read by no production code.\nS6 ask-user 6: docs/known-issues.md:6 retire rule; the deleted anchor-registry row's issue is still true.\nS7 no-op 4: tickets/registry_data.go name holds no registry.\nS8 no-op 3: gather.go:105 source shadows the parameter.\nCount 8 (2 hard, 6 judgment). Worst: S2.\n"
          },
          "axis": "Standards",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "1e8846bb094761993e922fd085648ab3de33344f",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "S5",
            "S6",
            "S7",
            "S8"
          ],
          "supersedes": []
        },
        {
          "id": "cr1-spec",
          "performer": "claude:ft293_cr1_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1_spec",
            "digest": "sha256:b8a791ae6d3bc523a8099956a9de860f8ce8368096cc9ef9b8fbb6a4b0d07605",
            "excerpt": "Spec axis, FT293 CR1 (base 36b1fafc4, chunk tip 1e8846bb0). Evidence check-current ran current.\nP1 auto-fix 6: CR06 \"The build charge verdict holds no charge-premises row\"; charge_test.go:193 asserts absence on prepared output that prints no check rows, so it cannot fail; no Decide test sets Operation to KindPrepareEvidence.\nP2 auto-fix 8: story 48; CHANGELOG.md:25 still states the charge-premises refusal.\nP3 auto-fix 6: spec.md:105 \"loses its binding-table grader, because that grader only keeps the closure's table valid\" is stale after plan commit 5fb3be66a.\nAXI keep is inside scope (check_floor.go:79-81 reads tickets.CommandPackage; story 11).\nRows: CR06 partial (P1); CR01-CR05, CR07-CR22, CR80-CR85, CR89, CR91, CR93, CR95, CR96 closed; CR15 closed by rg sweep; CR76 holds.\nCount 3. Worst: P2.\n"
          },
          "axis": "Spec",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "1e8846bb094761993e922fd085648ab3de33344f",
          "finding_ids": [
            "P1",
            "P2",
            "P3"
          ],
          "supersedes": []
        },
        {
          "id": "cr1-coverage",
          "performer": "claude:ft293_cr1_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b5f65579b446a4e6440a2680948bce3bcf63552",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1_coverage",
            "digest": "sha256:3fb3db6129505bfb7817c3cb0e6521f44b603dc233a8e6d008b1c2a7c9f91ac0",
            "excerpt": "Coverage axis, FT293 CR1. check-current current=true, tree clean before and after probes.\nC1 auto-fix 9: a green charge-premises row kept under PrepareEvidence/CurrentEvidence passes TestChargePreparesWithUnprovedPremise; probe 1 and 1b (insert green row in decision.go) silent, restored yes. Need a Decide row-list test under PrepareEvidence and CurrentEvidence.\nC2 auto-fix 6: tickets.ParseWrite accepts \"/\", \"/ (new)\", \"./\", \".\"; probe 2 bit, restored yes. Add \"/\" and \"/ (new)\" to TestParseWriteRefusesAnUncleanPath and refuse an empty path in ParseWrite.\nC3 auto-fix 7: a blacklist RepositoryPath body passes TestParseWrite*, TestWritesEntryCover, TestCommandBuildRefusesAnUncleanWritesEntry; probe 3 silent, restored yes. Add internal/./x.go refused and a..b/x.go accepted.\nPremise migration: no kept collector lost its only red-capable test. Command-package list: same 14 packages; grader breadth unchanged.\nCount 3. Worst: C1.\n"
          },
          "axis": "Coverage",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "1e8846bb094761993e922fd085648ab3de33344f",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
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
