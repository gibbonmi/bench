# Review outcomes

## CR2 second confirming-round pickup

The second confirming round read the repair delta `d66a978f..a883604f`. Each axis ran on opus at high effort. The chunk used 2 of its 2 repair cycles, and the second cycle was its hardening cycle.

The round closed S5, S6, and C6, and no finding stays open.

Optional advice, because each item grades only a check from the hardening cycle and no item shows a defect:

- The `wrongType` test constant copies the producer's state text with no recorded red.
- No rerun witness holds a slug or a basename with `'` or `$`.
- No `ParseTicket` test pins the first `Writes:` value of a ticket with a duplicate line.

## CR2 confirming-round pickup

The confirming round read the repair delta `390eb67d..d66a978f`. Each axis ran on opus at high effort. The chunk has used 1 of its 2 repair cycles, and the next repair is cycle 2.

The round closed S1, S2, S4, P1, P2, and C1 to C5. Three findings stay open, and ticket 03 owns all three.

- S5, auto-fix, confidence 4. `ParseTicket` and `WritesLine` in `internal/tickets/tickets.go` each write the rule that selects the live `Writes:` line. Give that rule one helper that both callers use.
- S6, auto-fix, confidence 5. The Coverage entry `cr2c-coverage` retains the red of the rerun assertion. The refusal text `outsideScope` in `widen_test.go` still has no recorded red. Derive that expectation from the producer's own text in the same package.
- C6, auto-fix, confidence 4. No witness pins the quoting of a ticket basename or a slug that holds a space. Add that case to `TestSpecWidenRefusesAnOutOfScopeEntry`.

Optional advice: the dogfood record names the two entries without their `(new)` marker.

## CR2 review pickup

The CR2 review ran on the frozen pair `1ea15b35..390eb67d`. Each axis ran on opus at high effort. The chunk has used 0 of its 2 repair cycles.

The raw count is 12 findings, and the repair-target count is 10, because S3 and S5 have no repair target. Ticket 03 owns every repair target. The orchestrator widens ticket 03 for S1, adds the plan requirement `t3-widen-table` for S4, and corrects the spec sentence of P2.

### Standards

Count: 5. Worst issue: S2.

- S1, auto-fix, confidence 7. `internal/spec/widen.go` restates the `Writes:` field name and the `, ` separator that `internal/tickets` and `internal/maps` own. Export each fact from its owner, and read it in the verb.
- S2, auto-fix, confidence 7. The out-of-scope refusal joins the refused entries without quotes, so a `(new)` entry breaks the rerun, and the rerun drops the in-scope entries. Quote each entry, and name the whole request.
- S3, no-op, confidence 4. The git error text in the hint slot follows the tree convention.
- S4, auto-fix, confidence 3. The table columns and the refusal text in `widen_test.go` are independent expectations with no recorded red. The `t3-widen-table` probe records that red.
- S5, no-op, confidence 3. The placement of `checkoutRefusal` follows the structure budget of `spec.go`.

### Spec

Count: 2. Worst issue: P1.

- P1, auto-fix, confidence 6. The CR35 catch does not hold, because `tickets.Covers` covers no path for a root prefix. Make the CR35 case red under a root-parent mutant, and correct the row's cell.
- P2, auto-fix, confidence 5. The spec and the changelog say that `--decision` records the reviewer's admission, but the verb stores no reference. Correct both sentences.

### Coverage

Count: 5. Worst issue: C1.

- C1, auto-fix, confidence 8. A marked or slash-spelled entry on either side reports `added` and appends a duplicate path. Add those cases to `TestSpecWidenReportsACoveredEntry`.
- C2, auto-fix, confidence 7. No witness pins the slash refusal of `--ticket sub/a.md` or of the slug `a/b`. Add both cases.
- C3, auto-fix, confidence 6. A default tip with only a non-ticket file or only a symbolic-link ticket counts as signed off. Add both cases to the CR53 test.
- C4, auto-fix, confidence 6. A ticket with `Writes: none` refuses a legal widening. Replace the sentinel, and add the case.
- C5, auto-fix, confidence 5. A signed-off entry with pathspec magic makes the scope read fail for the whole slug. Read the tree with literal pathspecs.

## CR1 third confirming-round pickup

The third confirming round read the delta `bad08290..1ea15b35`, which holds the test-only repair and its assignment. Each axis ran on opus at high effort. The chunk used 3 repair cycles, and the reviewer's `--auto-approve` covers the third.

The round closed C5, and no finding stays open. S10, no-op, confidence 7: the Coverage entry `cr1c3-coverage` retains two bypass probes that red the copied reason text, so the record now holds that red. The Standards axis reaffirmed this evidence-only disposition.

## CR1 second confirming-round pickup

The second confirming round read the repair delta `25863380..bad08290`. Each axis ran on opus at high effort. The chunk has used 2 of its 2 repair cycles. The reviewer's `--auto-approve` covers one more cycle, and the build log records that use.

The round closed S5, S9, C1, and C4 at the rule. One finding stays open, and ticket 02 owns it.

- C5, auto-fix, confidence 6. The sizing collector refuses a headroom premise whose directory is `""` or `.`, but no test pins that refusal. The orchestrator accepted the refusal on two consultations and flagged it for reviewer veto. Add a headroom case for each directory in the ownership package. This repair is test-only, so a coordinator confirmation closes it.

## CR1 confirming-round pickup

The confirming round read the repair delta `1e8846bb..25863380`. Each axis ran on opus at high effort. The chunk has used 1 of its 2 repair cycles, and the next repair is cycle 2.

The round closed S1, S2, S3, S4, P1, P2, P3, C2, and C3. Four findings stay open, and ticket 02 owns all four. The orchestrator adds the plan requirement `t2-charge-rows` for S9.

- S5, auto-fix, confidence 5. `PendingNative` in `internal/preflight/ownership/ownership.go` has no production reader. Remove the field, and let the tests read `Blocker`.
- S9, auto-fix, confidence 5. `TestDecideChargeRows` holds an independent row list with no recorded red. The `t2-charge-rows` probe records that red.
- C1, auto-fix, confidence 8. `TestDecideChargeRows` starts from facts with no pinned source tip, but a real prepare charge always pins it. Pin the tip in the prepare case, and expect the `tip-current` row.
- C4, auto-fix, confidence 7. The empty-path exemption of `tickets.RepositoryPath` is dead, and `ParseWrite` holds a second root rule. Give the root rule one source in `RepositoryPath`, and pin it with a test.

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
  "plan_digest": "sha256:b8aa3bef62333ca44c9dcab717464e45f80187bdce5c374ceb36ba739453eb66",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CR1",
      "base": "36b1fafc45426299b27b68f26fe627cadf000639",
      "tip": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
      "plan_digest": "sha256:2ad84f74d4d6d1a7784989f1d50a9fd2a6980b0da96e261cc62a2701ef03d4d9",
      "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
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
        },
        {
          "id": "t2r1-t2-preflight-packages",
          "performer": "claude:ft293_t2_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r1",
            "digest": "sha256:747d9f10e918e7d53dd628d6394c0a8707eaea054f46097044825a856bca6014",
            "excerpt": "packages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,24630,370\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,18178,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,9981,175\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5021,69\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-tickets-help-anchors",
          "performer": "claude:ft293_t2_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r1",
            "digest": "sha256:4395fd6ffaf3e78ec53cd128edf5e272e07f799b4ea604210b7ab810bbec2996",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/tickets/writes.go,omit,failed,1,TestParseWriteRefusesAnUncleanPath,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/tickets ./cmd/bench ./internal/anchors ./internal/canary ./internal/conformance,all,passed,1661\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,825\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15448,346\n  github.com/gibbonmi/bench/internal/anchors,pass,106,123\n  github.com/gibbonmi/bench/internal/canary,pass,8,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33425,1116\n  github.com/gibbonmi/bench/internal/tickets,fail,5,44\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/tickets,TestParseWriteRefusesAnUncleanPath,\"writes_test.go:57: ParseWrite(\\\"internal/x/../../cmd/y.go\\\") = {Path:internal/x/../../cmd/y.go Kind:0}, want an unclean-path refusal\",9\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
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
              "ref": "claude-agent:ft293_t2_r1",
              "digest": "sha256:4395fd6ffaf3e78ec53cd128edf5e272e07f799b4ea604210b7ab810bbec2996",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/tickets/writes.go,omit,failed,1,TestParseWriteRefusesAnUncleanPath,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/tickets ./cmd/bench ./internal/anchors ./internal/canary ./internal/conformance,all,passed,1661\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,825\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15448,346\n  github.com/gibbonmi/bench/internal/anchors,pass,106,123\n  github.com/gibbonmi/bench/internal/canary,pass,8,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33425,1116\n  github.com/gibbonmi/bench/internal/tickets,fail,5,44\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/tickets,TestParseWriteRefusesAnUncleanPath,\"writes_test.go:57: ParseWrite(\\\"internal/x/../../cmd/y.go\\\") = {Path:internal/x/../../cmd/y.go Kind:0}, want an unclean-path refusal\",9\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t2r1-t2-root-conformance",
          "performer": "claude:ft293_t2_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r1",
            "digest": "sha256:af03711acdc366d225e7f4f64095de0cf4ac210320a5371ec7aa6c7bc1b1032b",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13511,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-guidance-anchors",
          "performer": "claude:ft293_t2_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r1",
            "digest": "sha256:73f2d493ae09bbc2cbff31756f37bd3dc1da66688d6a0856ee631e8a656e7d92",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4161,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-guidance-budgets",
          "performer": "claude:ft293_t2_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r1",
            "digest": "sha256:7833b79d87da4342e711cb9d6c0cdd75b6c3afdfca28d197838c202bd9080c37",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t1r1-t1-preflight-packages",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:5534efc5515324d47d75ea9da179203dc882471dc7dc3bf3e22b74d04a1d2efe",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/evidencecmd/operations.go,swap+swap,failed,2,TestPreflightBuildRefusesRemovedForms,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight ./internal/preflight/ownership ./internal/preflight/ownership/ownershiptest ./internal/preflight/staleness ./internal/preflight/evidencecmd,all,passed,846\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,981\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,25186,370\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,fail,18310,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,10373,175\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5107,69\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestPreflightBuildRefusesRemovedForms,\"command_build_test.go:132: [build s --propose-writes --ticket a.md --base B --source-tip T] = (1), want exit 2 naming \\\"unknown argument: --propose-writes\\\":\",2\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,TestEvidenceFinalGuard,\"bound_test.go:26: operation \\\"build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>\\\" declares Bounded=false, want true for its kind\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:ft293_t1_r1",
              "digest": "sha256:5534efc5515324d47d75ea9da179203dc882471dc7dc3bf3e22b74d04a1d2efe",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/evidencecmd/operations.go,swap+swap,failed,2,TestPreflightBuildRefusesRemovedForms,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight ./internal/preflight/ownership ./internal/preflight/ownership/ownershiptest ./internal/preflight/staleness ./internal/preflight/evidencecmd,all,passed,846\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,981\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,25186,370\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,fail,18310,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,10373,175\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5107,69\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestPreflightBuildRefusesRemovedForms,\"command_build_test.go:132: [build s --propose-writes --ticket a.md --base B --source-tip T] = (1), want exit 2 naming \\\"unknown argument: --propose-writes\\\":\",2\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,TestEvidenceFinalGuard,\"bound_test.go:26: operation \\\"build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>\\\" declares Bounded=false, want true for its kind\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r1-t1-tickets-help-anchors",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:9864b2bf52b62f3fe46f720669c3282c1ba7180c9866c8da5db103b5b2b71325",
            "excerpt": "packages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15511,346\n  github.com/gibbonmi/bench/internal/anchors,pass,111,123\n  github.com/gibbonmi/bench/internal/canary,pass,7,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33458,1116\n  github.com/gibbonmi/bench/internal/tickets,pass,3,44\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-tickets-help-anchors",
          "command": "bench test --package ./internal/tickets --package ./cmd/bench --package ./internal/anchors --package ./internal/canary --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1r1-t1-root-conformance",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:45967cf546f1d3cc7028894ab08197fa7a09568c3c32f604fc4af87d39bd8377",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13702,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1r1-t1-guidance-anchors",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:cece434e58a39d00eef9cf04f6e4eec85187a0065b9309498a886c3c801d6a69",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4082,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t1r1-t1-guidance-budgets",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:7833b79d87da4342e711cb9d6c0cdd75b6c3afdfca28d197838c202bd9080c37",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t2r2-t2-preflight-packages",
          "performer": "claude:ft293_t2_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r2",
            "digest": "sha256:cd5f2be318a392d368cfccbab7eafbaaefed010ec3e68b4e5cfc1e88a40059ba",
            "excerpt": "packages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,24803,370\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,17698,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,9842,175\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,4858,69\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd",
          "exit_code": 0
        },
        {
          "id": "t2r2-t2-tickets-help-anchors",
          "performer": "claude:ft293_t2_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r2",
            "digest": "sha256:ed791bcf2b0a6b68e905e4e69f68d242c637d4994ec3aefabbe681f75229047f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/tickets/writes.go,omit,failed,1,TestParseWriteRefusesAnUncleanPath,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/tickets ./cmd/bench ./internal/anchors ./internal/canary ./internal/conformance,all,passed,1662\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,826\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17520,346\n  github.com/gibbonmi/bench/internal/anchors,pass,112,123\n  github.com/gibbonmi/bench/internal/canary,pass,9,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33525,1116\n  github.com/gibbonmi/bench/internal/tickets,fail,3,45\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/tickets,TestParseWriteRefusesAnUncleanPath,\"writes_test.go:57: ParseWrite(\\\"internal/x/../../cmd/y.go\\\") = {Path:internal/x/../../cmd/y.go Kind:0}, want an unclean-path refusal\",13\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
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
              "ref": "claude-agent:ft293_t2_r2",
              "digest": "sha256:ed791bcf2b0a6b68e905e4e69f68d242c637d4994ec3aefabbe681f75229047f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/tickets/writes.go,omit,failed,1,TestParseWriteRefusesAnUncleanPath,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/tickets ./cmd/bench ./internal/anchors ./internal/canary ./internal/conformance,all,passed,1662\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,826\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17520,346\n  github.com/gibbonmi/bench/internal/anchors,pass,112,123\n  github.com/gibbonmi/bench/internal/canary,pass,9,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33525,1116\n  github.com/gibbonmi/bench/internal/tickets,fail,3,45\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/tickets,TestParseWriteRefusesAnUncleanPath,\"writes_test.go:57: ParseWrite(\\\"internal/x/../../cmd/y.go\\\") = {Path:internal/x/../../cmd/y.go Kind:0}, want an unclean-path refusal\",13\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t2r2-t2-root-conformance",
          "performer": "claude:ft293_t2_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r2",
            "digest": "sha256:962316a1c14d2adf6f2303c8178c7e1ba1cb53b226766cb451faee92f0e9f2d4",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14139,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2r2-t2-guidance-anchors",
          "performer": "claude:ft293_t2_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r2",
            "digest": "sha256:0fe2c9e54a8165dd9a90c5c9351886ddaf6d6dd19d9dee6155dc375717c08e69",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4303,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t2r2-t2-guidance-budgets",
          "performer": "claude:ft293_t2_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r2",
            "digest": "sha256:7833b79d87da4342e711cb9d6c0cdd75b6c3afdfca28d197838c202bd9080c37",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t2r2-t2-charge-rows",
          "performer": "claude:ft293_t2_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r2",
            "digest": "sha256:f3cc1ea27002e4f7cdacb5a5b69c7b520c9c161555e45ac6b7103cb78ee6dc9f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/decision.go,insert,failed,1,TestDecideChargeRows,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,all,passed,370\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,509\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,21785,370\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestDecideChargeRows,\"charge_cells_test.go:46: Decide under operation 2 = red false, checks [base-current tip-current paths-authorized tickets-parse completion-plan check-floor blockers-resolve writes-resolve charge-premises fence-writes kit-pin binary-seal rows… (481 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-charge-rows",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0,
          "probe": {
            "mutation": "Insertion in Decide: append a green charge-premises row under KindPrepareEvidence when the source tip is pinned; TestDecideChargeRows must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t2_r2",
              "digest": "sha256:f3cc1ea27002e4f7cdacb5a5b69c7b520c9c161555e45ac6b7103cb78ee6dc9f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/decision.go,insert,failed,1,TestDecideChargeRows,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,all,passed,370\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,509\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,21785,370\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestDecideChargeRows,\"charge_cells_test.go:46: Decide under operation 2 = red false, checks [base-current tip-current paths-authorized tickets-parse completion-plan check-floor blockers-resolve writes-resolve charge-premises fence-writes kit-pin binary-seal rows… (481 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r2-t1-preflight-packages",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:4584bf9d7cab6ba943f6f9ff22e8cdf797182a1a9f739298ef2b4df403aa72ec",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/evidencecmd/operations.go,swap+swap,failed,2,TestPreflightBuildRefusesRemovedForms,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight ./internal/preflight/ownership ./internal/preflight/ownership/ownershiptest ./internal/preflight/staleness ./internal/preflight/evidencecmd,all,passed,846\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,980\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,24443,370\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,fail,17410,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,9930,175\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5018,69\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestPreflightBuildRefusesRemovedForms,\"command_build_test.go:132: [build s --propose-writes --ticket a.md --base B --source-tip T] = (1), want exit 2 naming \\\"unknown argument: --propose-writes\\\":\",2\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,TestEvidenceFinalGuard,\"bound_test.go:26: operation \\\"build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>\\\" declares Bounded=false, want true for its kind\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:ft293_t1_r1",
              "digest": "sha256:4584bf9d7cab6ba943f6f9ff22e8cdf797182a1a9f739298ef2b4df403aa72ec",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/evidencecmd/operations.go,swap+swap,failed,2,TestPreflightBuildRefusesRemovedForms,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight ./internal/preflight/ownership ./internal/preflight/ownership/ownershiptest ./internal/preflight/staleness ./internal/preflight/evidencecmd,all,passed,846\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,980\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,24443,370\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,fail,17410,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,9930,175\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5018,69\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestPreflightBuildRefusesRemovedForms,\"command_build_test.go:132: [build s --propose-writes --ticket a.md --base B --source-tip T] = (1), want exit 2 naming \\\"unknown argument: --propose-writes\\\":\",2\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,TestEvidenceFinalGuard,\"bound_test.go:26: operation \\\"build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>\\\" declares Bounded=false, want true for its kind\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r2-t1-tickets-help-anchors",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:241da0f18d71121947bf93c97d293553e18c8e8fb8ed31037a959e5631aae375",
            "excerpt": "packages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14818,346\n  github.com/gibbonmi/bench/internal/anchors,pass,116,123\n  github.com/gibbonmi/bench/internal/canary,pass,7,32\n  github.com/gibbonmi/bench/internal/conformance,pass,32973,1116\n  github.com/gibbonmi/bench/internal/tickets,pass,3,45\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-tickets-help-anchors",
          "command": "bench test --package ./internal/tickets --package ./cmd/bench --package ./internal/anchors --package ./internal/canary --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1r2-t1-root-conformance",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:abce10c7b8eb038d9fc7cb4428bc7dc19e8084814a02f13950e3b0f642f49032",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13890,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1r2-t1-guidance-anchors",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:8b8f7f72e6e2e89dea01853e7f782a5467e93be3b263799092e29a0b12b9ee02",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4021,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t1r2-t1-guidance-budgets",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:b682aaa6a2f3a10908927cd53ce11e0bf9cf598d196444453baf84591b6728d5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t2r3-t2-preflight-packages",
          "performer": "claude:ft293_t2_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r3",
            "digest": "sha256:1587bdb61942becabfcb2aa2b3106720bbd044e4c09c5bc7949c8038e913facb",
            "excerpt": "packages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,25073,370\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,18299,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,10310,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,4848,69\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd",
          "exit_code": 0
        },
        {
          "id": "t2r3-t2-tickets-help-anchors",
          "performer": "claude:ft293_t2_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r3",
            "digest": "sha256:63efc3dbeaeb8ef5be876f981006a50c2b33e1d3abd453885cdf493ef42d3e31",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/tickets/writes.go,omit,failed,1,TestParseWriteRefusesAnUncleanPath,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/tickets ./cmd/bench ./internal/anchors ./internal/canary ./internal/conformance,all,passed,1662\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,826\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15392,346\n  github.com/gibbonmi/bench/internal/anchors,pass,106,123\n  github.com/gibbonmi/bench/internal/canary,pass,7,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33168,1116\n  github.com/gibbonmi/bench/internal/tickets,fail,3,45\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/tickets,TestParseWriteRefusesAnUncleanPath,\"writes_test.go:57: ParseWrite(\\\"internal/x/../../cmd/y.go\\\") = {Path:internal/x/../../cmd/y.go Kind:0}, want an unclean-path refusal\",13\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
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
              "ref": "claude-agent:ft293_t2_r3",
              "digest": "sha256:63efc3dbeaeb8ef5be876f981006a50c2b33e1d3abd453885cdf493ef42d3e31",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/tickets/writes.go,omit,failed,1,TestParseWriteRefusesAnUncleanPath,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/tickets ./cmd/bench ./internal/anchors ./internal/canary ./internal/conformance,all,passed,1662\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,826\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15392,346\n  github.com/gibbonmi/bench/internal/anchors,pass,106,123\n  github.com/gibbonmi/bench/internal/canary,pass,7,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33168,1116\n  github.com/gibbonmi/bench/internal/tickets,fail,3,45\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/tickets,TestParseWriteRefusesAnUncleanPath,\"writes_test.go:57: ParseWrite(\\\"internal/x/../../cmd/y.go\\\") = {Path:internal/x/../../cmd/y.go Kind:0}, want an unclean-path refusal\",13\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t2r3-t2-root-conformance",
          "performer": "claude:ft293_t2_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r3",
            "digest": "sha256:efab2fe689fc0828064d403d96ea2eea565802a2f7c5db4cfe6d564ade1ecea1",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13540,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2r3-t2-guidance-anchors",
          "performer": "claude:ft293_t2_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r3",
            "digest": "sha256:e198407bfa6e189eb0099932884bfd97a9be959065f4dc658650484640bacd4e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4154,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t2r3-t2-guidance-budgets",
          "performer": "claude:ft293_t2_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r3",
            "digest": "sha256:7833b79d87da4342e711cb9d6c0cdd75b6c3afdfca28d197838c202bd9080c37",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t2r3-t2-charge-rows",
          "performer": "claude:ft293_t2_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t2_r3",
            "digest": "sha256:ea6eb7cd9b8a25417f15f6a0265ff53cab2efd3c0dd5e203962e486ab8986394",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/decision.go,insert,failed,1,TestDecideChargeRows,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,all,passed,370\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,509\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,21764,370\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestDecideChargeRows,\"charge_cells_test.go:46: Decide under operation 2 = red false, checks [base-current tip-current charge-premises paths-authorized tickets-parse completion-plan check-floor blockers-resolve writes-resolve fence-writes kit-pin binary-seal rows… (481 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-charge-rows",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0,
          "probe": {
            "mutation": "Insertion in Decide: append a green charge-premises row under KindPrepareEvidence when the source tip is pinned; TestDecideChargeRows must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t2_r3",
              "digest": "sha256:ea6eb7cd9b8a25417f15f6a0265ff53cab2efd3c0dd5e203962e486ab8986394",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/decision.go,insert,failed,1,TestDecideChargeRows,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,all,passed,370\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,509\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,21764,370\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestDecideChargeRows,\"charge_cells_test.go:46: Decide under operation 2 = red false, checks [base-current tip-current charge-premises paths-authorized tickets-parse completion-plan check-floor blockers-resolve writes-resolve fence-writes kit-pin binary-seal rows… (481 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r3-t1-preflight-packages",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:75bf348061336eb0ceaa63c6b2abcc6b2ddc68c70b6cd2354c0571020e8708b2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/evidencecmd/operations.go,swap+swap,failed,2,TestPreflightBuildRefusesRemovedForms,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight ./internal/preflight/ownership ./internal/preflight/ownership/ownershiptest ./internal/preflight/staleness ./internal/preflight/evidencecmd,all,passed,848\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,981\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,24819,370\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,fail,18379,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,10228,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,4884,69\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestPreflightBuildRefusesRemovedForms,\"command_build_test.go:132: [build s --propose-writes --ticket a.md --base B --source-tip T] = (1), want exit 2 naming \\\"unknown argument: --propose-writes\\\":\",2\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,TestEvidenceFinalGuard,\"bound_test.go:26: operation \\\"build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>\\\" declares Bounded=false, want true for its kind\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:ft293_t1_r1",
              "digest": "sha256:75bf348061336eb0ceaa63c6b2abcc6b2ddc68c70b6cd2354c0571020e8708b2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/evidencecmd/operations.go,swap+swap,failed,2,TestPreflightBuildRefusesRemovedForms,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight ./internal/preflight/ownership ./internal/preflight/ownership/ownershiptest ./internal/preflight/staleness ./internal/preflight/evidencecmd,all,passed,848\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,981\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,24819,370\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,fail,18379,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,10228,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,4884,69\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestPreflightBuildRefusesRemovedForms,\"command_build_test.go:132: [build s --propose-writes --ticket a.md --base B --source-tip T] = (1), want exit 2 naming \\\"unknown argument: --propose-writes\\\":\",2\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,TestEvidenceFinalGuard,\"bound_test.go:26: operation \\\"build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>\\\" declares Bounded=false, want true for its kind\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r3-t1-tickets-help-anchors",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:fc937023a325129fd80960349bb49371e771b409a57644fbfdda135a5be98705",
            "excerpt": "packages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14863,346\n  github.com/gibbonmi/bench/internal/anchors,pass,111,123\n  github.com/gibbonmi/bench/internal/canary,pass,7,32\n  github.com/gibbonmi/bench/internal/conformance,pass,33098,1116\n  github.com/gibbonmi/bench/internal/tickets,pass,3,45\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-tickets-help-anchors",
          "command": "bench test --package ./internal/tickets --package ./cmd/bench --package ./internal/anchors --package ./internal/canary --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1r3-t1-root-conformance",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:41a675f765a1f899bb965ea229328b705aa3093014cbf23814fbd4b5e0ed065b",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13344,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1r3-t1-guidance-anchors",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
            "digest": "sha256:dfbf84a9d0c0baef21be43079d969552c9bfda9b43f97e3443d2d96cffeb39e0",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4644,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t1r3-t1-guidance-budgets",
          "performer": "claude:ft293_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t1_r1",
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
        },
        {
          "id": "cr1c-standards",
          "performer": "claude:ft293_cr1c_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1c_standards",
            "digest": "sha256:4be3f52a87ccf1b4cc2598c4c23c5da1eec81b659b2b3c92b4d75b790319e797",
            "excerpt": "Standards confirming round, FT293 CR1, repair delta 1e8846bb..25863380. check-current passed.\nS1 closed 9; S2 closed 9; S3 closed 9; S4 closed 9.\nS5 open 5: PendingNative kept, read only by premises_test.go:353,:364; the Blocker checks already prove a pending claim.\nS9 new 5: TestDecideChargeRows (charge_cells_test.go:40-45) hand-lists 13 check names with no recorded red (AGENTS.md independent-expectation exception).\nOpen: S5, S9. Worst: S5.\n"
          },
          "axis": "Standards",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "2586338063a1fd52e1b4ef10f8e1a7c2effe1131",
          "finding_ids": [
            "S5",
            "S9"
          ],
          "supersedes": [
            "cr1-standards"
          ]
        },
        {
          "id": "cr1c-spec",
          "performer": "claude:ft293_cr1c_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1c_spec",
            "digest": "sha256:6671ae6f19eff9c1e9cf00d64d0d93744dc281ca4b9e2387c38b80b5fa1f3ed0",
            "excerpt": "Spec confirming round, FT293 CR1, repair delta 1e8846bb..25863380. check-current passed.\nP1 closed 9: TestDecideChargeRows pins the exact charge row list under KindPrepareEvidence and KindCurrentEvidence; CR05 holds.\nP2 closed 8: CHANGELOG.md:22-25 matches the tree; CR93 holds.\nP3 closed 9: spec line 105 matches ticket 02 and ticket_grammar_test.go:24-25.\nCR19, CR20 hold 9. No new finding. Open: 0.\n"
          },
          "axis": "Spec",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "2586338063a1fd52e1b4ef10f8e1a7c2effe1131",
          "finding_ids": [],
          "supersedes": [
            "cr1-spec"
          ]
        },
        {
          "id": "cr1c-coverage",
          "performer": "claude:ft293_cr1c_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e9bf7a93e3580fe323fbfb7d0591344b8de789a9",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1c_coverage",
            "digest": "sha256:5b265735598b259a3701764018a9541c3ef7930a6c2761b4fe0eb992d80bfc99",
            "excerpt": "Coverage confirming round, FT293 CR1. check-current current=true; 8 probes, all restored=yes; tree clean.\nC1 open 8: a green charge-premises row gated on KindPrepareEvidence && PinnedSourceTip != \"\" is silent across ./internal/preflight/...; TestDecideChargeRows uses baseFacts() with no pinned tip and omits tip-current.\nC2 closed 9: guard swaps bit on \"/\" and \"./\".\nC3 closed 9: Clean-without-IsLocal and contains-\"..\" predicates both bit.\nC4 new 7: RepositoryPath empty-path exemption (writes.go:75-77) is dead; removing it keeps all tests green; the root rule has two sources.\nOpen: C1, C4. Worst: C1.\n"
          },
          "axis": "Coverage",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "2586338063a1fd52e1b4ef10f8e1a7c2effe1131",
          "finding_ids": [
            "C1",
            "C4"
          ],
          "supersedes": [
            "cr1-coverage"
          ]
        },
        {
          "id": "cr1c2-standards",
          "performer": "claude:ft293_cr1c2_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1c2_standards",
            "digest": "sha256:e65087045e9724a76a07fba6966ffe0a983ac6c26336912fae03ec18714f86d9",
            "excerpt": "Standards second confirming round, FT293 CR1, repair delta 25863380..bad08290. check-current passed.\nS5 closed 9: PendingNative gone (ownership.go:52-59, premises.go:277-279); tests read Blocker.\nS9 closed 9: t2-charge-rows plan row (spec.md:208); record entry t2r2-t2-charge-rows probe bit, restore pass.\nC4 closed 8: one root rule at writes.go:72-73, pinned by writes_test.go:67-75; no second root rule in ParseWrite or sizing.go.\nNew findings: none. Open 0.\n"
          },
          "axis": "Standards",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "bad08290250ed94ac3d1fc5c37c309db8e90cbd3",
          "finding_ids": [],
          "supersedes": [
            "cr1c-standards"
          ]
        },
        {
          "id": "cr1c2-spec",
          "performer": "claude:ft293_cr1c2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1c2_spec",
            "digest": "sha256:d69b67c5291c1355cb644bf7c0dc9f76b124045be64796152b6543a786fd6bf6",
            "excerpt": "Spec second confirming round, FT293 CR1, repair delta 2586338..bad08290. check-current current=true.\nCR06 held 9; CR11, CR12 held 8; CR19-CR22 held 9.\nRoot headroom refusal: no contradiction 8; spec line 111 says the premise collectors call the grammar's rule.\nNew findings: none. Open 0.\n"
          },
          "axis": "Spec",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "bad08290250ed94ac3d1fc5c37c309db8e90cbd3",
          "finding_ids": [],
          "supersedes": [
            "cr1c-spec"
          ]
        },
        {
          "id": "cr1c2-coverage",
          "performer": "claude:ft293_cr1c2_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "5dc40802a6e6296ad8747dd0de5618101e560a47",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1c2_coverage",
            "digest": "sha256:1fdb71f33f6d2068b65b30aec1f379cf9d6e17ca90a547a36115639c368179f7",
            "excerpt": "Coverage second confirming round, FT293 CR1. check-current current=true; 6 probes, all restored=yes; tree clean.\nC1 closed 9: unpinned KindCurrentEvidence green row and n/a row between base-current and tip-current both bit TestDecideChargeRows.\nC4 closed at the rule, open at one caller 7: sizing.go:156 guard insertion silent (ownership, 175 tests); the headroom caller has no \"\" or \".\" test.\nS5 closed 8: completeness.go:22 skip bit TestAuditCompleteness/SS57_SS58; premises.go:199 green bit TestChargeManualApplicability/OC89_OC93.\nC5 new 6: the \".\" headroom refusal changed with no test or spec row; add a headroom row for \".\" and \"\" in the ownership package.\nOpen: C5 (with the C4 caller gap). Worst: C5.\n"
          },
          "axis": "Coverage",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "bad08290250ed94ac3d1fc5c37c309db8e90cbd3",
          "finding_ids": [
            "C5"
          ],
          "supersedes": [
            "cr1c-coverage"
          ]
        },
        {
          "id": "cr1c3-spec",
          "performer": "claude:ft293_cr1c3_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1c3_spec",
            "digest": "sha256:2770708ac168ddd3a14a0d955f3e570186028e4e460f07e4037f8bb2cfa77ffd",
            "excerpt": "Spec third confirming round, FT293 CR1, delta bad08290..1ea15b35. check-current current=true.\nNo production code in the delta: held 10. spec.md change is only the claude:ft293_t2_r3 assignment: held 10.\nCR1 rows stay closed: held 8. New TestSizingPremises cases pin the decided refusal (spec line 111): held 9.\nNew findings: none. Open 0.\n"
          },
          "axis": "Spec",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "finding_ids": [],
          "supersedes": [
            "cr1c2-spec"
          ]
        },
        {
          "id": "cr1c3-coverage",
          "performer": "claude:ft293_cr1c3_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1c3_coverage",
            "digest": "sha256:7df789e8ca37cade09ebbd9ceff570b6962079072ba9828dfe2510621022fa89",
            "excerpt": "Coverage third confirming round, FT293 CR1, delta bad08290..1ea15b35. check-current current=true; tree clean.\nC5 closed 9. Probes, none repeating the author or coordinator mutation:\n1. writes.go:73 RepositoryPath admits \"\" (prepend p == \"\" ||): bit, OC36_headroom-empty-directory failed, restored yes.\n2. sizing.go:156 check RepositoryPath(strings.TrimLeft(dir+\"/x\", \"./\")): bit, OC36_headroom-dot-directory and OC36_headroom-empty-directory failed, restored yes.\n3. writes.go:73 path.Clean(p) == p swapped to path.Clean(p) != \"\": bit, TestParseWriteRefusesAnUncleanPath failed, restored yes.\n\"./\" and \"x/..\" are refused by the one clean-spelling rule that TestParseWriteRefusesAnUncleanPath pins.\nNew findings: 0.\n"
          },
          "axis": "Coverage",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "finding_ids": [],
          "supersedes": [
            "cr1c2-coverage"
          ]
        },
        {
          "id": "cr1c3-standards",
          "performer": "claude:ft293_cr1c3_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f1debd11ebddd7a617aa148b96f5074b6c2e3a2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr1c3_standards",
            "digest": "sha256:e2ddae1d3b4db0004ed077f9fa8ea2bafc5f0f8b3e52270e60491c6c034b8a55",
            "excerpt": "Standards third confirming round, FT293 CR1, delta bad08290..1ea15b35. check-current current=true.\nC5 closed 8: two cases through headroomItem (premises_test.go:125-126, :148-149) hit sizing.go:156-157; doc comment :93-95 accurate.\nspec.md change is only the claude:ft293_t2_r3 assignment (line 185).\nS10 raised 5 (copied reason text with no in-tree red), then reaffirmed closed 7 as evidence-only after record commit bafcdf413 retained the cr1c3-coverage bypass probes (reviews/closure-removal.md:1214-1225).\nOpen 0.\n"
          },
          "axis": "Standards",
          "base": "36b1fafc45426299b27b68f26fe627cadf000639",
          "tip": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "finding_ids": [],
          "supersedes": [
            "cr1c2-standards"
          ]
        }
      ]
    },
    {
      "id": "CR2",
      "base": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
      "tip": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
      "plan_digest": "sha256:9421a83c29076fbf94debbe657d4e004e894c1ba8c9d4a336c3af0cd6712f7bc",
      "source_digest": "59dd4fd54fddd3eb9c4c14dacc0846e3d11ea1bb",
      "acceptance_rows": [
        "CR23",
        "CR24",
        "CR25",
        "CR26",
        "CR27",
        "CR28",
        "CR29",
        "CR30",
        "CR31",
        "CR32",
        "CR33",
        "CR34",
        "CR35",
        "CR36",
        "CR37",
        "CR38",
        "CR39",
        "CR40",
        "CR41",
        "CR42",
        "CR43",
        "CR44",
        "CR45",
        "CR46",
        "CR47",
        "CR48",
        "CR49",
        "CR50",
        "CR51",
        "CR52",
        "CR53",
        "CR54",
        "CR92"
      ],
      "verification": [
        {
          "id": "t3-t3-spec-and-help",
          "performer": "claude:ft293_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6b9636ad28b76a8ace91461a3ef6b2542a9cf341",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3",
            "digest": "sha256:ae6f97512cff772142aa23904c823a805be22dd0f9d552ae0b690cf598829794",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,7,TestSpecWidenAdmitsADecidedEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec ./cmd/bench ./internal/anchors ./internal/conformance,all,passed,1737\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2375\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16348,346\n  github.com/gibbonmi/bench/internal/anchors,pass,116,123\n  github.com/gibbonmi/bench/internal/conformance,pass,35491,1116\n  github.com/gibbonmi/bench/internal/spec,fail,1941,152\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:147: widen = (0, \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,intern… (310 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_directory,\"widen_test.go:217: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_with_its_prefix_spelling,\"widen_test.go:217: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/git2/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_refuses_another_root_file,\"widen_test.go:217: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,README.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:188: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/new/other.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenRefusesAnOutOfScopeEntry,\"widen_test.go:140: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenWritesAllOrNothing,\"widen_test.go:246: code = 0, want 1; out = \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-spec-and-help",
          "command": "bench test --package ./internal/spec --package ./cmd/bench --package ./internal/anchors --package ./internal/conformance",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the scope check: make the scope check admit every entry; TestSpecWidenRefusesAnOutOfScopeEntry must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t3",
              "digest": "sha256:ae6f97512cff772142aa23904c823a805be22dd0f9d552ae0b690cf598829794",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,7,TestSpecWidenAdmitsADecidedEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec ./cmd/bench ./internal/anchors ./internal/conformance,all,passed,1737\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2375\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16348,346\n  github.com/gibbonmi/bench/internal/anchors,pass,116,123\n  github.com/gibbonmi/bench/internal/conformance,pass,35491,1116\n  github.com/gibbonmi/bench/internal/spec,fail,1941,152\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:147: widen = (0, \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,intern… (310 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_directory,\"widen_test.go:217: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_with_its_prefix_spelling,\"widen_test.go:217: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/git2/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_refuses_another_root_file,\"widen_test.go:217: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,README.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:188: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/new/other.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenRefusesAnOutOfScopeEntry,\"widen_test.go:140: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenWritesAllOrNothing,\"widen_test.go:246: code = 0, want 1; out = \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t3-t3-root-conformance",
          "performer": "claude:ft293_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6b9636ad28b76a8ace91461a3ef6b2542a9cf341",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3",
            "digest": "sha256:b47bdb9414b0d88ee7b235880e7f6a476f62deb76d88830272774fa099795062",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14896,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3-t3-axi-registry",
          "performer": "claude:ft293_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6b9636ad28b76a8ace91461a3ef6b2542a9cf341",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3",
            "digest": "sha256:9bb0abdb3205d48eb90e52b735058f8f1eb39c284c726525d8897e6fe96907e6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  axi-query-registry,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-axi-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t3-t3-command-routing",
          "performer": "claude:ft293_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6b9636ad28b76a8ace91461a3ef6b2542a9cf341",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3",
            "digest": "sha256:af39cc20f715525000234871bf92af20ae8b9a8714a0e5b065d4f13a01fc6f36",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  subcommand-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,26,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-command-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t3r1-t3-spec-and-help",
          "performer": "claude:ft293_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "362a6a32ef914aa627fafd665359931873ddcef6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r1",
            "digest": "sha256:5c1f850f6cc95cc05ef8831c4d8b10ca74d977e53d1968468e49c26946e4ba40",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,7,TestSpecWidenAdmitsADecidedEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec ./internal/tickets ./internal/maps ./cmd/bench ./internal/anchors ./internal/conformance,all,passed,1888\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2476\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19231,346\n  github.com/gibbonmi/bench/internal/anchors,pass,135,123\n  github.com/gibbonmi/bench/internal/conformance,pass,36613,1116\n  github.com/gibbonmi/bench/internal/maps,pass,129,97\n  github.com/gibbonmi/bench/internal/spec,fail,2152,161\n  github.com/gibbonmi/bench/internal/tickets,pass,3,45\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:150: widen = (0, \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,intern… (310 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_directory,\"widen_test.go:219: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_with_its_prefix_spelling,\"widen_test.go:219: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/git2/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:231: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,README.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:191: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/new/other.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenRefusesAnOutOfScopeEntry,\"widen_test.go:140: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenWritesAllOrNothing,\"widen_test.go:244: code = 0, want 1; out = \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-spec-and-help",
          "command": "bench test --package ./internal/spec --package ./internal/tickets --package ./internal/maps --package ./cmd/bench --package ./internal/anchors --package ./internal/conformance",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the scope check: make the scope check admit every entry; TestSpecWidenRefusesAnOutOfScopeEntry must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t3_r1",
              "digest": "sha256:5c1f850f6cc95cc05ef8831c4d8b10ca74d977e53d1968468e49c26946e4ba40",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,7,TestSpecWidenAdmitsADecidedEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec ./internal/tickets ./internal/maps ./cmd/bench ./internal/anchors ./internal/conformance,all,passed,1888\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2476\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19231,346\n  github.com/gibbonmi/bench/internal/anchors,pass,135,123\n  github.com/gibbonmi/bench/internal/conformance,pass,36613,1116\n  github.com/gibbonmi/bench/internal/maps,pass,129,97\n  github.com/gibbonmi/bench/internal/spec,fail,2152,161\n  github.com/gibbonmi/bench/internal/tickets,pass,3,45\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:150: widen = (0, \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,intern… (310 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_directory,\"widen_test.go:219: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_with_its_prefix_spelling,\"widen_test.go:219: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/git2/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:231: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,README.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:191: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/new/other.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenRefusesAnOutOfScopeEntry,\"widen_test.go:140: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenWritesAllOrNothing,\"widen_test.go:244: code = 0, want 1; out = \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t3r1-t3-widen-table",
          "performer": "claude:ft293_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "362a6a32ef914aa627fafd665359931873ddcef6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r1",
            "digest": "sha256:2f4816c1e5d18a0e1944649a76173ec92ba5dee1700c6135e372d4e071c1f276",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,6,TestSpecWidenAddsAnInScopeEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,all,passed,159\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2146\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1696,159\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAddsAnInScopeEntry,\"widen_test.go:130: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:150: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,decided\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,int… (313 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:235: root file widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsSignedOffEntriesLiterally,\"widen_signedoff_test.go:50: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:198: out = \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\", want \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReportsACoveredEntry,\"widen_line_test.go:12: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y.go,current\\\\\\\\n  specs/s/tickets/a.md,internal/status/t.go,current\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a… (323 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-widen-table",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap in internal/spec/widen.go: rename the widen table column state to status; TestSpecWidenAddsAnInScopeEntry must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t3_r1",
              "digest": "sha256:2f4816c1e5d18a0e1944649a76173ec92ba5dee1700c6135e372d4e071c1f276",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,6,TestSpecWidenAddsAnInScopeEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,all,passed,159\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2146\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1696,159\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAddsAnInScopeEntry,\"widen_test.go:130: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:150: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,decided\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,int… (313 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:235: root file widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsSignedOffEntriesLiterally,\"widen_signedoff_test.go:50: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:198: out = \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\", want \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReportsACoveredEntry,\"widen_line_test.go:12: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y.go,current\\\\\\\\n  specs/s/tickets/a.md,internal/status/t.go,current\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a… (323 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3r1-t3-root-conformance",
          "performer": "claude:ft293_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "362a6a32ef914aa627fafd665359931873ddcef6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r1",
            "digest": "sha256:05899aa6bb4ab4d11dc03b91963e89f6b3e76bae9a0935682f8ea564d50f588d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15261,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3r1-t3-axi-registry",
          "performer": "claude:ft293_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "362a6a32ef914aa627fafd665359931873ddcef6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r1",
            "digest": "sha256:945ed42bc8506a59c3a98e468a4380dc6dc7bfb668bcd12bbcf937a6c8909cd6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  axi-query-registry,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,11,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-axi-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t3r1-t3-command-routing",
          "performer": "claude:ft293_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "362a6a32ef914aa627fafd665359931873ddcef6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r1",
            "digest": "sha256:af39cc20f715525000234871bf92af20ae8b9a8714a0e5b065d4f13a01fc6f36",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  subcommand-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,26,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-command-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t3r2-t3-spec-and-help",
          "performer": "claude:ft293_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a835df0f41c0746f4ce88193f844a38cf52d5f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r1",
            "digest": "sha256:6f3d228318c3f7f3c8a4ebc3df48dd3cf9213cf0721e408db59572e43c2d43d2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,7,TestSpecWidenAdmitsADecidedEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec ./internal/tickets ./internal/maps ./cmd/bench ./internal/anchors ./internal/conformance,all,passed,1888\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2476\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16804,346\n  github.com/gibbonmi/bench/internal/anchors,pass,117,123\n  github.com/gibbonmi/bench/internal/conformance,pass,36384,1116\n  github.com/gibbonmi/bench/internal/maps,pass,110,97\n  github.com/gibbonmi/bench/internal/spec,fail,2152,161\n  github.com/gibbonmi/bench/internal/tickets,pass,3,45\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:150: widen = (0, \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,intern… (310 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_directory,\"widen_test.go:219: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_with_its_prefix_spelling,\"widen_test.go:219: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/git2/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:231: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,README.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:191: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/new/other.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenRefusesAnOutOfScopeEntry,\"widen_test.go:140: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenWritesAllOrNothing,\"widen_test.go:244: code = 0, want 1; out = \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-spec-and-help",
          "command": "bench test --package ./internal/spec --package ./internal/tickets --package ./internal/maps --package ./cmd/bench --package ./internal/anchors --package ./internal/conformance",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the scope check: make the scope check admit every entry; TestSpecWidenRefusesAnOutOfScopeEntry must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t3_r1",
              "digest": "sha256:6f3d228318c3f7f3c8a4ebc3df48dd3cf9213cf0721e408db59572e43c2d43d2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,7,TestSpecWidenAdmitsADecidedEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec ./internal/tickets ./internal/maps ./cmd/bench ./internal/anchors ./internal/conformance,all,passed,1888\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2476\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16804,346\n  github.com/gibbonmi/bench/internal/anchors,pass,117,123\n  github.com/gibbonmi/bench/internal/conformance,pass,36384,1116\n  github.com/gibbonmi/bench/internal/maps,pass,110,97\n  github.com/gibbonmi/bench/internal/spec,fail,2152,161\n  github.com/gibbonmi/bench/internal/tickets,pass,3,45\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:150: widen = (0, \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,intern… (310 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_directory,\"widen_test.go:219: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_with_its_prefix_spelling,\"widen_test.go:219: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/git2/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:231: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,README.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:191: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/new/other.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenRefusesAnOutOfScopeEntry,\"widen_test.go:140: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenWritesAllOrNothing,\"widen_test.go:244: code = 0, want 1; out = \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t3r2-t3-widen-table",
          "performer": "claude:ft293_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a835df0f41c0746f4ce88193f844a38cf52d5f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r1",
            "digest": "sha256:346dc7084b3f3bd380b585eb51d217e6c98f0a1682566527f42ee87e70c8a6d9",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,6,TestSpecWidenAddsAnInScopeEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,all,passed,159\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2146\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1682,159\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAddsAnInScopeEntry,\"widen_test.go:130: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:150: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,decided\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,int… (313 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:235: root file widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsSignedOffEntriesLiterally,\"widen_signedoff_test.go:50: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:198: out = \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\", want \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReportsACoveredEntry,\"widen_line_test.go:12: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y.go,current\\\\\\\\n  specs/s/tickets/a.md,internal/status/t.go,current\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a… (323 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-widen-table",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap in internal/spec/widen.go: rename the widen table column state to status; TestSpecWidenAddsAnInScopeEntry must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t3_r1",
              "digest": "sha256:346dc7084b3f3bd380b585eb51d217e6c98f0a1682566527f42ee87e70c8a6d9",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,6,TestSpecWidenAddsAnInScopeEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,all,passed,159\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2146\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1682,159\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAddsAnInScopeEntry,\"widen_test.go:130: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:150: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,decided\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,int… (313 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:235: root file widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsSignedOffEntriesLiterally,\"widen_signedoff_test.go:50: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:198: out = \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\", want \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReportsACoveredEntry,\"widen_line_test.go:12: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y.go,current\\\\\\\\n  specs/s/tickets/a.md,internal/status/t.go,current\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a… (323 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3r2-t3-root-conformance",
          "performer": "claude:ft293_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a835df0f41c0746f4ce88193f844a38cf52d5f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r1",
            "digest": "sha256:f402d09322b9044376e4662e293de9828e197a891e7131c4e3206d2434fd0c0f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14725,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3r2-t3-axi-registry",
          "performer": "claude:ft293_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a835df0f41c0746f4ce88193f844a38cf52d5f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r1",
            "digest": "sha256:945ed42bc8506a59c3a98e468a4380dc6dc7bfb668bcd12bbcf937a6c8909cd6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  axi-query-registry,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,11,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-axi-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t3r2-t3-command-routing",
          "performer": "claude:ft293_t3_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a835df0f41c0746f4ce88193f844a38cf52d5f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r1",
            "digest": "sha256:7869d8783dc8a1904b65ac15f6cd7de12ced4f90afeb9403faf7f18e534d6676",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  subcommand-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,27,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-command-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t3r3-t3-spec-and-help",
          "performer": "claude:ft293_t3_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "59dd4fd54fddd3eb9c4c14dacc0846e3d11ea1bb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r2",
            "digest": "sha256:d79d53b2dd87219520e835d79649c2fa353af56cd2b029fcf33b6e8fcf68a331",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,7,TestSpecWidenAdmitsADecidedEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec ./internal/tickets ./internal/maps ./cmd/bench ./internal/anchors ./internal/conformance,all,passed,1889\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2476\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19265,346\n  github.com/gibbonmi/bench/internal/anchors,pass,128,123\n  github.com/gibbonmi/bench/internal/conformance,pass,36598,1116\n  github.com/gibbonmi/bench/internal/maps,pass,130,97\n  github.com/gibbonmi/bench/internal/spec,fail,2126,162\n  github.com/gibbonmi/bench/internal/tickets,pass,4,45\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:153: widen = (0, \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,intern… (310 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_directory,\"widen_test.go:222: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_with_its_prefix_spelling,\"widen_test.go:222: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/git2/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:234: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,README.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:194: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/new/other.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenRefusesAnOutOfScopeEntry,\"widen_test.go:137: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenWritesAllOrNothing,\"widen_test.go:247: code = 0, want 1; out = \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-spec-and-help",
          "command": "bench test --package ./internal/spec --package ./internal/tickets --package ./internal/maps --package ./cmd/bench --package ./internal/anchors --package ./internal/conformance",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap at the scope check: make the scope check admit every entry; TestSpecWidenRefusesAnOutOfScopeEntry must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t3_r2",
              "digest": "sha256:d79d53b2dd87219520e835d79649c2fa353af56cd2b029fcf33b6e8fcf68a331",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,7,TestSpecWidenAdmitsADecidedEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec ./internal/tickets ./internal/maps ./cmd/bench ./internal/anchors ./internal/conformance,all,passed,1889\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2476\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19265,346\n  github.com/gibbonmi/bench/internal/anchors,pass,128,123\n  github.com/gibbonmi/bench/internal/conformance,pass,36598,1116\n  github.com/gibbonmi/bench/internal/maps,pass,130,97\n  github.com/gibbonmi/bench/internal/spec,fail,2126,162\n  github.com/gibbonmi/bench/internal/tickets,pass,4,45\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:153: widen = (0, \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,intern… (310 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_directory,\"widen_test.go:222: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_directory_entry_refuses_a_sibling_with_its_prefix_spelling,\"widen_test.go:222: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/git2/x.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:234: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,README.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:194: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/new/other.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenRefusesAnOutOfScopeEntry,\"widen_test.go:137: code = 0, want 1; out = \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenWritesAllOrNothing,\"widen_test.go:247: code = 0, want 1; out = \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,added\\\\\\\\n\\\"\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t3r3-t3-widen-table",
          "performer": "claude:ft293_t3_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "59dd4fd54fddd3eb9c4c14dacc0846e3d11ea1bb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r2",
            "digest": "sha256:601fb86a13849978a1b0bb7e4b7b9730d04a1a7b646e9ec32f71c64b62fd29e7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,6,TestSpecWidenAddsAnInScopeEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,all,passed,160\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2146\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1689,160\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAddsAnInScopeEntry,\"widen_test.go:127: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:153: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,decided\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,int… (313 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:238: root file widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsSignedOffEntriesLiterally,\"widen_signedoff_test.go:50: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:201: out = \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\", want \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReportsACoveredEntry,\"widen_line_test.go:12: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y.go,current\\\\\\\\n  specs/s/tickets/a.md,internal/status/t.go,current\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a… (323 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-widen-table",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap in internal/spec/widen.go: rename the widen table column state to status; TestSpecWidenAddsAnInScopeEntry must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t3_r2",
              "digest": "sha256:601fb86a13849978a1b0bb7e4b7b9730d04a1a7b646e9ec32f71c64b62fd29e7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/widen.go,swap,failed,6,TestSpecWidenAddsAnInScopeEntry,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,all,passed,160\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2146\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1689,160\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAddsAnInScopeEntry,\"widen_test.go:127: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y_test.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenAdmitsADecidedEntry,\"widen_test.go:153: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/w.go,added\\\\\\\\n  specs/s/tickets/a.md,internal/other/z.go,decided\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,int… (313 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenDirectoryPrefix/a_root_file_entry_is_its_own_prefix,\"widen_test.go:238: root file widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,CHANGELOG.md,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsSignedOffEntriesLiterally,\"widen_signedoff_test.go:50: widen = (0, \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"), want exit 0 and \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReadsTheSignedOffScope,\"widen_test.go:201: out = \\\"widen[1]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\", want \\\"widen[1]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/b.md,internal/x/z.go,added\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/spec,TestSpecWidenReportsACoveredEntry,\"widen_line_test.go:12: widen = (0, \\\"widen[2]{ticket,entry,status}:\\\\\\\\n  specs/s/tickets/a.md,internal/x/y.go,current\\\\\\\\n  specs/s/tickets/a.md,internal/status/t.go,current\\\\\\\\n\\\"), want exit 0 and \\\"widen[2]{ticket,entry,state}:\\\\\\\\n  specs/s/tickets/a… (323 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3r3-t3-root-conformance",
          "performer": "claude:ft293_t3_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "59dd4fd54fddd3eb9c4c14dacc0846e3d11ea1bb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r2",
            "digest": "sha256:0796e26a49ee9fef583295b5df771d151a2099f52b18676bdb5583c51c13328a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15012,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3r3-t3-axi-registry",
          "performer": "claude:ft293_t3_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "59dd4fd54fddd3eb9c4c14dacc0846e3d11ea1bb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r2",
            "digest": "sha256:9bb0abdb3205d48eb90e52b735058f8f1eb39c284c726525d8897e6fe96907e6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  axi-query-registry,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-axi-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t3r3-t3-command-routing",
          "performer": "claude:ft293_t3_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "59dd4fd54fddd3eb9c4c14dacc0846e3d11ea1bb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t3_r2",
            "digest": "sha256:af39cc20f715525000234871bf92af20ae8b9a8714a0e5b065d4f13a01fc6f36",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  subcommand-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,26,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-command-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "cr2-standards",
          "performer": "claude:ft293_cr2_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "6b9636ad28b76a8ace91461a3ef6b2542a9cf341",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr2_standards",
            "digest": "sha256:af93306804b531c7a120a343e9d8dd49e10d85fee343f9d4da41c49093726e79",
            "excerpt": "Standards axis, FT293 CR2 (base 1ea15b35, tip 390eb67d). check-current current=true, dirty_record=false.\nS1 auto-fix 7: widen.go:40 writesField, :43 writesSeparator, :247 FieldScan table restate the Writes: syntax (tickets/tickets.go:47) and the \", \" separator (maps/fields.go:136); one-source-per-fact. Route: plan expansion of ticket 03 Writes.\nS2 auto-fix 7: widen.go:121 joins refused entries with \" \" unquoted; a (new) entry breaks the rerun, and the rerun drops in-scope entries; axi.ShellQuote exists (axi/action.go:284).\nS3 no-op 4: widen.go:190,195,205,216 pass git err text into the hint slot.\nS4 no-op unless no recorded red 3: widen_test.go:74 repeats the column list; :24 repeats the refusal text.\nS5 no-op 3: checkoutRefusal lives in retire.go.\nRetire move preserves behavior; WidenUsage is one source. Count 5. Worst: S2.\n"
          },
          "axis": "Standards",
          "base": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "tip": "390eb67dc1eea75dee6d9c63b8a3d8c40ceed185",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "S5"
          ],
          "supersedes": []
        },
        {
          "id": "cr2-spec",
          "performer": "claude:ft293_cr2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "6b9636ad28b76a8ace91461a3ef6b2542a9cf341",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr2_spec",
            "digest": "sha256:ec18edab980ef08aff7eefc2c8dc8110289c5690166ef3a3dd18c1b42eb915b7",
            "excerpt": "Spec axis, FT293 CR2 (base 1ea15b35, tip 390eb67d). check-current current, not dirty.\nP1 non-blocking 6: CR35 \"A root parent rule admits every path\" does not hold; Covers(\".\", p) covers nothing, so the CR35 subcase stays green under a root-parent mutant (writes.go:89-91, widen_test.go:208).\nP2 non-behavioral contradiction 5: spec.md:27 and CHANGELOG.md:29 say --decision \"records the reviewer's admission\"; the verb stores or prints no reference (widen.go:56, :107; spec.md:337 won't-handle).\nRows: CR35 partial; CR23-CR34, CR36-CR54, CR92 closed. All \"The plan verb\" decisions implemented; no won't-handle scope creep.\nCount 2. Worst: P1.\n"
          },
          "axis": "Spec",
          "base": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "tip": "390eb67dc1eea75dee6d9c63b8a3d8c40ceed185",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "cr2-coverage",
          "performer": "claude:ft293_cr2_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "6b9636ad28b76a8ace91461a3ef6b2542a9cf341",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr2_coverage",
            "digest": "sha256:290cbe41161ed3025a88ea301cfac6eb6c4d45413c8f5f36328878bcd55b7640",
            "excerpt": "Coverage axis, FT293 CR2 (1ea15b35..390eb67d). check-current current=true, dirty=false; tree clean before and after.\nProbes (./internal/spec, 152 tests, each restored=yes): P1 widen.go:95 append(covered, entry) silent; P2 :153 omit slash check silent; P3 :200 omit Ext check silent; P4 :266 if false silent; P5 :101 p := entry silent; P6 :200 omit IsRegularFile silent.\nC1 fix 8: marked or slash-spelled entries on either side report added and append a duplicate (CR38).\nC2 fix 7: --ticket sub/a.md and slug a/b have no witness for the slash check.\nC3 fix 6: a tip with only a non-.md file or only a symlink ticket counts as signed off (CR53).\nC4 fix 6: a Writes: none ticket refuses a legal widening.\nC5 flag 5: a :(glob) signed-off entry makes git ls-tree exit 128 and refuses every widening for the slug.\nCount 5. Worst: C1.\n"
          },
          "axis": "Coverage",
          "base": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "tip": "390eb67dc1eea75dee6d9c63b8a3d8c40ceed185",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
            "C5"
          ],
          "supersedes": []
        },
        {
          "id": "cr2c-standards",
          "performer": "claude:ft293_cr2c_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a835df0f41c0746f4ce88193f844a38cf52d5f3",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr2c_standards",
            "digest": "sha256:3885273b924a89e117c310906b4cccfa2183dcece51a03a128dd6d44226213b4",
            "excerpt": "Standards confirming round, FT293 CR2, repair delta 390eb67d..d66a978f. check-current current=true, dirty_record=false.\nS1 closed 8: tickets.WritesSyntax (tickets.go:53), maps.FieldSeparator (fields.go:131); local constants deleted.\nS2 closed 9: rerun names every entry, shell-quoted (widen.go:113-118; widen_test.go:331-332).\nS4 closed 7: t3r2-t3-widen-table probe bit, restore pass (reviews/closure-removal.md:1565-1589).\nS5 new 4: the live Writes: line selection is written twice (tickets.go:107-117 ParseTicket, :150-153 WritesLine); one helper could serve both.\nS6 new 5: the rerun-string assertion (widen_test.go:332) and the outsideScope refusal text (:24) have no recorded red.\nOpen: S5, S6. Worst: S6.\n"
          },
          "axis": "Standards",
          "base": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "tip": "d66a978fbea4282cd30f11f09252802f587f17a1",
          "finding_ids": [
            "S5",
            "S6"
          ],
          "supersedes": [
            "cr2-standards"
          ]
        },
        {
          "id": "cr2c-spec",
          "performer": "claude:ft293_cr2c_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a835df0f41c0746f4ce88193f844a38cf52d5f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr2c_spec",
            "digest": "sha256:f62a4522680bb0c355a2cbad9b77da35963a3bec73d80fcd454df4ceb9eb648a",
            "excerpt": "Spec confirming round, FT293 CR2, repair delta 390eb67d..d66a978f. check-current current=true, dirty_record=false.\nP1 closed 9: the b.md widen of CHANGELOG.md (widen_test.go:233-236) reds a root-parent mutant; spec.md:263 states that catch.\nP2 closed 9: spec.md:27 and CHANGELOG.md:29 say --decision admits a path after the reviewer decides and stores no reference.\nCR23-CR54, CR92 held 8; moved cells CR38, CR40-CR42 (widen_line_test.go), CR53 (widen_signedoff_test.go) held 9; no unrequested behavior 8; dogfood record held 7.\nAdvice (not a finding): spec.md:660 names the dogfood entries without their (new) marker.\nOpen 0.\n"
          },
          "axis": "Spec",
          "base": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "tip": "d66a978fbea4282cd30f11f09252802f587f17a1",
          "finding_ids": [],
          "supersedes": [
            "cr2-spec"
          ]
        },
        {
          "id": "cr2c-coverage",
          "performer": "claude:ft293_cr2c_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a835df0f41c0746f4ce88193f844a38cf52d5f3",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr2c_coverage",
            "digest": "sha256:dfb7379285d7f3da397aed2ef16af8bc31080d28897a81d7ffbb34b923d2ec8a",
            "excerpt": "Coverage confirming round, FT293 CR2. check-current current at 28d95f2e; 11 probes over ./internal/spec, every restored=yes; tree clean.\nC1 closed 9: added-entry append(covered, entry) bit; ticket-side TrimSuffix only bit; request-side strips markers not slash bit.\nC2 closed 9: pathSegment -> RepositoryPath on basename bit; on slug bit.\nC3 closed 9: IsRegularFile -> Type != \"blob\" bit; signed = len(listing) > 0 bit.\nC4 closed 8: kept = \" \" -> \"\" bit.\nC5 closed 7: drop :(literal) in scopePrefix bit; on the folder call silent (folder always starts with specs/).\nS2 closed 8: rerun loops over refused (widen.go:114) bit TestSpecWidenRefusesAnOutOfScopeEntry at the rerun assertion.\nC6 new 4: ShellQuote(basename) -> basename silent; pathSegment admits a spaced basename or slug; add a witness.\nNested specs/<slug>/tickets/sub/x.md: not a gap; tickets enumeration counts nested tickets.\nOpen: C6. Worst: C6.\n"
          },
          "axis": "Coverage",
          "base": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "tip": "d66a978fbea4282cd30f11f09252802f587f17a1",
          "finding_ids": [
            "C6"
          ],
          "supersedes": [
            "cr2-coverage"
          ]
        },
        {
          "id": "cr2c2-standards",
          "performer": "claude:ft293_cr2c2_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "59dd4fd54fddd3eb9c4c14dacc0846e3d11ea1bb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr2c2_standards",
            "digest": "sha256:ff722d77fc47a28b5d8566f5ef13518d15d494ba1571f3f448cd8d9488635f1d",
            "excerpt": "Standards second confirming round, FT293 CR2, delta d66a978f..a883604f. check-current current=true, dirty_record=false.\nS5 closed 9: liveLines (tickets.go:153-159) is the one selection rule for ParseTicket and WritesLine; behavior unchanged.\nS6 closed 8: outsideScope lives in widen.go:149; tests read it; the rerun red is recorded in cr2c-coverage.\nAdvice (hardening-cycle check, no defect): the wrongType test constant copies the producer state text with no recorded red.\nOpen 0.\n"
          },
          "axis": "Standards",
          "base": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "tip": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "finding_ids": [],
          "supersedes": [
            "cr2c-standards"
          ]
        },
        {
          "id": "cr2c2-spec",
          "performer": "claude:ft293_cr2c2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "59dd4fd54fddd3eb9c4c14dacc0846e3d11ea1bb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr2c2_spec",
            "digest": "sha256:5f58490d25de4bb50723ae2b173cc747a975bf32548219b5e7eaeed343d0eb47",
            "excerpt": "Spec second confirming round, FT293 CR2, delta d66a978f..a883604f. check-current current=true.\nCR23-CR54 held 10; no ticket-grammar change held 9 (one live line per field); no unrequested behavior held 10; CR92 held 9.\nNew findings: none. Open 0.\n"
          },
          "axis": "Spec",
          "base": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "tip": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "finding_ids": [],
          "supersedes": [
            "cr2c-spec"
          ]
        },
        {
          "id": "cr2c2-coverage",
          "performer": "claude:ft293_cr2c2_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "59dd4fd54fddd3eb9c4c14dacc0846e3d11ea1bb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr2c2_coverage",
            "digest": "sha256:1594dab0e56c58c0b594c3a48168324545b05a20f61e79dc91a42f4ac7d93d69",
            "excerpt": "Coverage second confirming round, FT293 CR2. Evidence current at a14573820; probes one at a time, every restored=yes; tree clean.\nC6 closed 8: raw slug when spaced bit; plain double quotes on basename bit (widen_test.go:145).\nS5 closed 7: last Writes: line bit TestSpecWidenKeepsOtherBytes/a_duplicate_Writes:_line_stays; last line for other fields bit TestParseTicketFieldMechanics/a_second_Blocked_by_line_is_a_duplicate_field; drop !line.Fenced silent and equivalent.\nAdvice (hardening-cycle checks, no defect): a slug or basename with ' or $ has no rerun witness; no ParseTicket test pins the first Writes: value for a duplicate line.\ninternal/tickets 45 tests and internal/maps 97 tests pass. Open 0.\n"
          },
          "axis": "Coverage",
          "base": "1ea15b3595545d8a829e52b82f0753b64b4327f9",
          "tip": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "finding_ids": [],
          "supersedes": [
            "cr2c-coverage"
          ]
        }
      ]
    },
    {
      "id": "CR3",
      "base": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
      "tip": "8230977099a3702443b965c5af3f5f24b6859173",
      "plan_digest": "sha256:b8aa3bef62333ca44c9dcab717464e45f80187bdce5c374ceb36ba739453eb66",
      "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
      "acceptance_rows": [
        "CR55",
        "CR56",
        "CR57",
        "CR58",
        "CR59",
        "CR60",
        "CR61",
        "CR62",
        "CR63",
        "CR64",
        "CR65",
        "CR66",
        "CR67",
        "CR68",
        "CR69",
        "CR72",
        "CR77",
        "CR90",
        "CR70",
        "CR71",
        "CR86",
        "CR87"
      ],
      "verification": [
        {
          "id": "t5-t5-preflight-packages",
          "performer": "claude:ft293_t5",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5",
            "digest": "sha256:7a560b9e784fa8a0c102444697f16b1fa294a11718f4394b14f69b5e9a766aa3",
            "excerpt": "packages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,pass,6977,101\n  github.com/gibbonmi/bench/internal/preflight,pass,28469,365\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,20283,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,11676,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5548,70\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/preflighttest --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd --package ./internal/diff",
          "exit_code": 0
        },
        {
          "id": "t5-t5-fence-readers",
          "performer": "claude:ft293_t5",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5",
            "digest": "sha256:ced19b784642ef1842633e27db5b750c6335ca6fbcffe68f2d387c4b4fa07dd4",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/coverage/citations.go,insert,failed,1,TestCommandCheckIgnoresAFenceSection,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/coverage ./internal/spec ./internal/chargeevidence ./internal/refusalroute ./internal/reviewrecord/recordtest ./internal/anchors ./cmd/bench ./internal/conformance,all,passed,2025\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1333\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16529,346\n  github.com/gibbonmi/bench/internal/anchors,pass,112,123\n  github.com/gibbonmi/bench/internal/chargeevidence,pass,24,122\n  github.com/gibbonmi/bench/internal/conformance,pass,36727,1116\n  github.com/gibbonmi/bench/internal/coverage,fail,3213,141\n  github.com/gibbonmi/bench/internal/refusalroute,pass,5,20\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/spec,pass,2312,157\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/coverage,TestCommandCheckIgnoresAFenceSection,\"fence_section_test.go:13: Command = (1, \\\"error: specs/s/spec.md coverage map spec declares ownership fences without the review pickup — fix the map or mark it <!-- coverage-map: historical -->\\\\\\\\n\\\"), want exit 0 for a fence section without th… (257 bytes)\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/chargeevidence,TestEvidenceStoreKinds/CE94_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t5-fence-readers",
          "command": "bench test --package ./internal/coverage --package ./internal/spec --package ./internal/chargeevidence --package ./internal/refusalroute --package ./internal/reviewrecord/recordtest --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0,
          "probe": {
            "mutation": "Insertion at the coverage parser: restore the checkReviewPickup grading of a fence section that omits the review pickup; TestCommandCheckIgnoresAFenceSection must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t5",
              "digest": "sha256:ced19b784642ef1842633e27db5b750c6335ca6fbcffe68f2d387c4b4fa07dd4",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/coverage/citations.go,insert,failed,1,TestCommandCheckIgnoresAFenceSection,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/coverage ./internal/spec ./internal/chargeevidence ./internal/refusalroute ./internal/reviewrecord/recordtest ./internal/anchors ./cmd/bench ./internal/conformance,all,passed,2025\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1333\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16529,346\n  github.com/gibbonmi/bench/internal/anchors,pass,112,123\n  github.com/gibbonmi/bench/internal/chargeevidence,pass,24,122\n  github.com/gibbonmi/bench/internal/conformance,pass,36727,1116\n  github.com/gibbonmi/bench/internal/coverage,fail,3213,141\n  github.com/gibbonmi/bench/internal/refusalroute,pass,5,20\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/spec,pass,2312,157\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/coverage,TestCommandCheckIgnoresAFenceSection,\"fence_section_test.go:13: Command = (1, \\\"error: specs/s/spec.md coverage map spec declares ownership fences without the review pickup — fix the map or mark it <!-- coverage-map: historical -->\\\\\\\\n\\\"), want exit 0 for a fence section without th… (257 bytes)\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/chargeevidence,TestEvidenceStoreKinds/CE94_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t5-t5-landing",
          "performer": "claude:ft293_t5",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5",
            "digest": "sha256:44cce94048c7564bbc85f1a9ff26331fd6bd00c4f3ec9acea86f49752031311e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,65825,1437\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-landing",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t5-t5-system",
          "performer": "claude:ft293_t5",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5",
            "digest": "sha256:c8da41aaa2bb011e278b94c22fd3496606d7d8b8ee7a4b0e9210c744d94b611a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,170,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,69726,170\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t5-t5-root-conformance",
          "performer": "claude:ft293_t5",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5",
            "digest": "sha256:d7d5646183536e8406b949599bb413bb742acf08fa26fab0cc4b9c85b6b232db",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15699,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t5-t5-guidance-anchors",
          "performer": "claude:ft293_t5",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5",
            "digest": "sha256:a56176fd0e0206f80d291ba9f4e20068af41b670df57b559e09a6788fb18b198",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4543,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t5-t5-guidance-budgets",
          "performer": "claude:ft293_t5",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5",
            "digest": "sha256:b682aaa6a2f3a10908927cd53ce11e0bf9cf598d196444453baf84591b6728d5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t4-t4-preflight-packages",
          "performer": "claude:ft293_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4",
            "digest": "sha256:978878f40db95e872394d2134ca127157ec4d3469ff2b71a32678e5a31d7d8e5",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,27323,365\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,19710,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,11250,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5735,70\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/preflighttest --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd",
          "exit_code": 0
        },
        {
          "id": "t4-t4-fence-readers",
          "performer": "claude:ft293_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4",
            "digest": "sha256:24ca6f3c6564eed444bd09da4c7c5eb6ec23775e6782e194fffd09386c53a04e",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16745,346\n  github.com/gibbonmi/bench/internal/anchors,pass,109,123\n  github.com/gibbonmi/bench/internal/chargeevidence,pass,22,122\n  github.com/gibbonmi/bench/internal/conformance,pass,36749,1116\n  github.com/gibbonmi/bench/internal/coverage,pass,3295,141\n  github.com/gibbonmi/bench/internal/refusalroute,pass,5,20\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/spec,pass,2364,157\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/chargeevidence,TestEvidenceStoreKinds/CE94_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-fence-readers",
          "command": "bench test --package ./internal/coverage --package ./internal/spec --package ./internal/chargeevidence --package ./internal/refusalroute --package ./internal/reviewrecord/recordtest --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t4-t4-landing",
          "performer": "claude:ft293_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4",
            "digest": "sha256:e22b02d9afaa81688421b0b25f9a2f773db303c55e1bd9bd5888c9f7fef567b3",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/fence_writes.go,swap,failed,196,TestCommitmentClosureAdmission/rowless-FT2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1432\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,70657\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,600097,1432\nfailures[196]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/rowless-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT3,\"commitment_landing_test.go:101: closing FT3 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/rowless,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/spec,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/rowless,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/spec,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentDeliveryClosure,\"commitment_landing_test.go:27: delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/blocker,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/default-branch-policy,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentInterruptedGate,\"commitment_landing_test.go:143: interrupted closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLandingRefusesObligationFreeDelivery,\"commitment_light_landing_test.go:327: obligation-free landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/partly-delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/beyond-scope,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (751 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/directory-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/listed-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/outside-directory,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (763 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/scope-without-deliverable,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (799 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/sibling-prefix,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (757 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/unlisted,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (753 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-then-tickets,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-with-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-without-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/tickets-then-spec,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (721 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPartialDelivery,\"commitment_landing_test.go:56: partial landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPersistenceBeforeGate,\"commitment_landing_test.go:173: torn-intent landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPublishLock,\"land_identity_test.go:307: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then b… (661 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentStaleLanding,\"land_identity_test.go:209: displaced landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCompareAndSwapRefusalNamesTheFold,\"land_moved_destination_test.go:66: moved destination = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold,\"land_moved_destination_test.go:43: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (733 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:132: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict-pending,\"refusal_route_test.go:132: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestDelegatedIntegrationJourney,\"delegated_integration_test.go:237: delegated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (649 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/installed,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (741 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/kit-checkout,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (746 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCleansTheFoldedSibling,\"land_effects_cleanup_test.go:83: folded-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAcceptsAbbreviatedIdentities,\"land_surface_test.go:46: abbreviated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthenticatesDigestShapedRequestToken,\"land_reauthorization_test.go:102: digest-shaped request land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (756 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesCaptureOutsideTheFence,\"land_surface_test.go:135: capture landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (676 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesFromTicketWrites,\"land_derived_fence_test.go:33: ticket-only landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/another_review_record,\"land_derived_fence_test.go:67: other review landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (742 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/pickup,\"land_derived_fence_test.go:53: pickup landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandCarriesTheBaselineScheduleRootIntoTheProspectiveGate,\"land_freshness_test.go:288: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (700 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandComposesCaptureOntoMovedDestination,\"land_surface_test.go:71: unreviewed capture composition = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (682 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictNextPointsThroughUnsafePath,\"land_surface_test.go:360: unsafe-path conflict next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesThePath,\"land_surface_test.go:268: conflict refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesTheSourceRepair,\"land_surface_test.go:301: conflict repair next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandDisclosesAUnionResolution,\"land_surface_test.go:104: unreviewed union = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedBase,\"land_reauthorization_test.go:175: abbreviated base = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"f6270c2f3bbe\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (725 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/39,\"land_reauthorization_test.go:163: abbreviated source tip \\\"db8e618d29422b77bc0cdd19b811eefdd28bb43\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or a… (752 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/4,\"land_reauthorization_test.go:163: abbreviated source tip \\\"6173\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line … (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/upper-12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"19B54E59C4B2\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandFenceRefusalNamesThePath,\"land_surface_test.go:220: fence refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; t… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal,\"land_freshness_test.go:114: forged-primary landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (699 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIncompleteNextUsesAssignmentPointerForUnsafePath,\"land_resume_test.go:38: unsafe-path incomplete = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease,\"land_freshness_test.go:72: single-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandLeavesTheDestinationUnchangedAfterARedProspectiveGate,\"land_freshness_test.go:185: gate tally = \\\"\\\", open /dev/shm/QYZTUF/t/TestLandCommandLeavesTheDestinationUnchangedAfterARedProspective3665232664/003/.git/bench-land-gate-tally: no such file or directory, want one prospective run\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout,\"land_broker_notice_test.go:125: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (720 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion,\"land_freshness_test.go:39: stable-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/marker,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (644 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/reconcile,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (649 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/release,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (645 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsNoHeadsLineWithoutRecords,\"land_census_test.go:91: empty census landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsTheCensusHeadBreakdown,\"land_census_test.go:76: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/absent,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (648 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/concurrently-moved,\"land_flags_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/present,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (649 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicConflictRepairRequiresNewReviewedTip,\"land_journey_test.go:294: conflict result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicLandsAnInRangeSpecAmendment,\"land_spec_amendment_test.go:47: amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicPreservesHistoricalRuntimeLogs,\"land_journey_test.go:139: runtime-log landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/clean,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (646 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/declared-output,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-and-unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log-empty-declaration,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (664 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicResumeCompletesPublishedReleaseWithoutRepublishing,\"land_resume_test.go:55: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublishedReleaseFailureExitsIncomplete,\"land_journey_test.go:267: published release exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalKeepsControlBearingPathInOneTableRow,\"land_release_refusal_test.go:70: control refusal = (1, \\\"refused{detail=landing destination is not clean,next=reviewer: commit the destination's uncommitted work, or discard it; then bench worktree land --request 'refusal-controls' --base '5… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalListsCollidingPaths,\"land_release_refusal_test.go:52: collision refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusesPostGateUnknownIgnoredMutation,\"land_journey_test.go:175: post-gate ignored mutation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReleaseDiagnosticCannotForgeTerminalLines,\"land_flags_test.go:132: hostile release result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRemovesEveryTemporaryProspectiveArtifact/published,\"land_freshness_test.go:214: landing exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:54: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (703 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/marker,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (682 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/reconcile,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/release,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeOverlappingDestinationEdit,\"land_journey_test.go:237: last-moment overlapping edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (741 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeTrackedDestinationEdit,\"land_journey_test.go:199: last-moment tracked edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesTheCensusCountAndDropsTheRecords,\"land_census_test.go:45: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesZeroForAnAssignmentWithNoRecords,\"land_census_test.go:61: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandEffectsRowPrecedesTheLandedRecord,\"land_effects_test.go:238: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then be… (669 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandFenceRefusalOverAFoldedSourceKeepsTheCallerBase,\"land_folded_base_test.go:34: fence refusal = \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (717 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandGradesASourceCommittedByALanePass,\"land_journey_test.go:381: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandHoldsTheDestinationLease,\"land_lease_test.go:90: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writ… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandLeavesAPriorLandedAssignment,\"land_effects_cleanup_test.go:114: narrowed landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshReadsTheSeal,\"land_effects_test.go:220: unsealed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshesTheBrokerAfterPublication,\"land_effects_test.go:160: landing with a stale broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRelaysTheBoundedGreenShapeBeforeTheLandedRecord,\"land_journey_test.go:410: land exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReleasesTheLeaseOnRefusal,\"land_lease_test.go:171: red landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedCleanup,\"land_effects_cleanup_test.go:188: faulted cleanup = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedRefresh,\"land_effects_test.go:181: failed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandResumeTakesNoLease,\"land_lease_test.go:190: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (633 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsAnUnprovenSibling,\"land_effects_cleanup_test.go:140: retaining landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/ordinary,\"land_empty_sibling_test.go:45: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/resumed,\"land_empty_sibling_test.go:35: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (664 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsAFreshBroker,\"land_effects_test.go:140: landing with a fresh broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsTheRefreshWithoutBuildInputs,\"land_effects_test.go:118: landing without build inputs = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict,\"refusal_route_follow_test.go:59: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-compositi… (617 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict-pending,\"refusal_route_follow_test.go:59: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-c… (641 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-collision,\"refusal_route_follow_test.go:54: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (734 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-not-clean,\"refusal_route_follow_test.go:72: destination-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (730 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-build-window,\"refusal_route_follow_test.go:72: land-build-window re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-destination-moved,\"refusal_route_follow_test.go:54: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (737 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-handback,\"refusal_route_follow_test.go:72: land-handback re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-incomplete,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-infrastructure,\"refusal_route_follow_test.go:69: printed step \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen\\\" is not one simple command with every slot filled\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-lease-held,\"refusal_route_follow_test.go:72: land-lease-held re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-red,\"refusal_route_follow_test.go:59: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base '900… (553 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-destination-residue,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-marker,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-clean,\"refusal_route_follow_test.go:72: source-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-fenced,\"refusal_route_follow_test.go:72: source-not-fenced re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-tip-mismatch,\"refusal_route_follow_test.go:72: source-tip-mismatch re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingPrintsOutputBreakdown,\"land_census_output_test.go:25: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRedRouteNamesTheRepair,\"refusal_route_test.go:147: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base 'f1f6e08a4… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/destination-collision,\"identity_component_test.go:238: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (744 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-destination-moved,\"identity_component_test.go:238: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (745 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-incomplete,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (687 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-destination-residue,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (701 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-marker,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingReleaseFactAdapterTranslatesReleaseExit,\"land_facts_test.go:211: release-failure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRetiresAnAbsentLandedSibling,\"missing_tree_recovery_test.go:231: absent-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/dangling-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/directory,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/live-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/current-source,\"land_moved_destination_test.go:132: current-source = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/non-proof-cause,\"land_moved_destination_test.go:132: non-proof-cause = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/unreadable-ancestry,\"land_moved_destination_test.go:126: rename /dev/shm/QYZTUF/t/TestProofFaultOverACurrentSourceHandsBackunreadable-ancestry2520870877/003/.git/objects/a3/d6e9901b8628b0df10120d1cb1bebf5e7aecad.hidden /dev/shm/QYZTUF/t/TestProofFaultOverACurre… (372 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeHandsBackADetachedDestination,\"land_resume_refusal_test.go:366: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (666 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/./specs/x/spec.md,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/x,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandCompletesAnInterruptedMarker,\"land_resume_test.go:230: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (679 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandExpandsAbbreviatedIdentities,\"land_reauthorization_test.go:184: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandFollowupFailureExitsIncomplete,\"land_resume_test.go:19: first incomplete exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-path,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (719 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-state,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (722 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/lock,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/owner-marker,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/registration,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/request,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-path,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-request,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicBindsPublishedLandingIdentity,\"land_resume_test.go:101: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (669 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicCompletesAnAmendedSourceLanding,\"land_spec_amendment_test.go:75: interrupted amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandReconcilesAnUnreconciledPublishedCheckout,\"land_resume_test.go:183: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/absent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/behind,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/divergent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (687 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesNonAncestorReviewBaseWithoutMutation,\"land_resume_refusal_test.go:143: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (675 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesWhenTerminalReceiptWasEvicted,\"land_resume_refusal_test.go:204: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRepeatsTheCensusCount,\"land_resume_test.go:285: incomplete landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandSourceRefusalNamesTheCallersResume,\"land_resume_refusal_test.go:119: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandUnknownRequestNamesReauthorizeRecovery,\"land_reauthorization_test.go:16: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding,\"land_specless_test.go:151: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCompletesTheUnfinishedEffects,\"land_effects_cleanup_test.go:218: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandRunsTheEffectsAfterRelease,\"land_effects_cleanup_test.go:270: failed-refresh landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeReadsEffectStateFromTheTree,\"land_effects_test.go:269: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict,\"refusal_route_test.go:187: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict-pending,\"refusal_route_test.go:187: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/destination-collision,\"refusal_route_test.go:185: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-destination-residue,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-marker,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestSecondLandingMeetsTheLease,\"land_lease_test.go:133: first landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnclaimedLandingCompositionErrorHandsBack,\"refusal_route_test.go:164: unclaimed composition error = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (736 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnsafePathRouteUsesThePlaceholder/incomplete_resume,\"land_release_refusal_test.go:204: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (558 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowAdmitsTheBuildsOwnLanding,\"land_window_test.go:119: landing of the bound spec from the sibling = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: l… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowClosesAtCompletion,\"land_window_test.go:130: landing of the build's completion = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (682 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/build-start,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/gap-close,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (686 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-landing",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "Omission: drop the review pickup from the implicit entries; TestLandCommandAuthorizesTheReviewPickup must fail, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t4",
              "digest": "sha256:e22b02d9afaa81688421b0b25f9a2f773db303c55e1bd9bd5888c9f7fef567b3",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/fence_writes.go,swap,failed,196,TestCommitmentClosureAdmission/rowless-FT2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1432\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,70657\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,600097,1432\nfailures[196]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/rowless-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT3,\"commitment_landing_test.go:101: closing FT3 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/rowless,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/spec,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/rowless,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/spec,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentDeliveryClosure,\"commitment_landing_test.go:27: delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/blocker,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/default-branch-policy,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentInterruptedGate,\"commitment_landing_test.go:143: interrupted closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLandingRefusesObligationFreeDelivery,\"commitment_light_landing_test.go:327: obligation-free landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/partly-delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/beyond-scope,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (751 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/directory-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/listed-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/outside-directory,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (763 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/scope-without-deliverable,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (799 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/sibling-prefix,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (757 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/unlisted,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (753 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-then-tickets,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-with-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-without-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/tickets-then-spec,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (721 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPartialDelivery,\"commitment_landing_test.go:56: partial landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPersistenceBeforeGate,\"commitment_landing_test.go:173: torn-intent landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPublishLock,\"land_identity_test.go:307: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then b… (661 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentStaleLanding,\"land_identity_test.go:209: displaced landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCompareAndSwapRefusalNamesTheFold,\"land_moved_destination_test.go:66: moved destination = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold,\"land_moved_destination_test.go:43: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (733 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:132: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict-pending,\"refusal_route_test.go:132: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestDelegatedIntegrationJourney,\"delegated_integration_test.go:237: delegated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (649 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/installed,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (741 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/kit-checkout,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (746 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCleansTheFoldedSibling,\"land_effects_cleanup_test.go:83: folded-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAcceptsAbbreviatedIdentities,\"land_surface_test.go:46: abbreviated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthenticatesDigestShapedRequestToken,\"land_reauthorization_test.go:102: digest-shaped request land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (756 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesCaptureOutsideTheFence,\"land_surface_test.go:135: capture landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (676 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesFromTicketWrites,\"land_derived_fence_test.go:33: ticket-only landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/another_review_record,\"land_derived_fence_test.go:67: other review landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (742 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/pickup,\"land_derived_fence_test.go:53: pickup landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandCarriesTheBaselineScheduleRootIntoTheProspectiveGate,\"land_freshness_test.go:288: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (700 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandComposesCaptureOntoMovedDestination,\"land_surface_test.go:71: unreviewed capture composition = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (682 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictNextPointsThroughUnsafePath,\"land_surface_test.go:360: unsafe-path conflict next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesThePath,\"land_surface_test.go:268: conflict refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesTheSourceRepair,\"land_surface_test.go:301: conflict repair next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandDisclosesAUnionResolution,\"land_surface_test.go:104: unreviewed union = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedBase,\"land_reauthorization_test.go:175: abbreviated base = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"f6270c2f3bbe\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (725 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/39,\"land_reauthorization_test.go:163: abbreviated source tip \\\"db8e618d29422b77bc0cdd19b811eefdd28bb43\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or a… (752 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/4,\"land_reauthorization_test.go:163: abbreviated source tip \\\"6173\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line … (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/upper-12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"19B54E59C4B2\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandFenceRefusalNamesThePath,\"land_surface_test.go:220: fence refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; t… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal,\"land_freshness_test.go:114: forged-primary landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (699 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIncompleteNextUsesAssignmentPointerForUnsafePath,\"land_resume_test.go:38: unsafe-path incomplete = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease,\"land_freshness_test.go:72: single-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandLeavesTheDestinationUnchangedAfterARedProspectiveGate,\"land_freshness_test.go:185: gate tally = \\\"\\\", open /dev/shm/QYZTUF/t/TestLandCommandLeavesTheDestinationUnchangedAfterARedProspective3665232664/003/.git/bench-land-gate-tally: no such file or directory, want one prospective run\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout,\"land_broker_notice_test.go:125: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (720 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion,\"land_freshness_test.go:39: stable-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/marker,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (644 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/reconcile,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (649 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/release,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (645 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsNoHeadsLineWithoutRecords,\"land_census_test.go:91: empty census landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsTheCensusHeadBreakdown,\"land_census_test.go:76: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/absent,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (648 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/concurrently-moved,\"land_flags_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/present,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (649 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicConflictRepairRequiresNewReviewedTip,\"land_journey_test.go:294: conflict result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicLandsAnInRangeSpecAmendment,\"land_spec_amendment_test.go:47: amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicPreservesHistoricalRuntimeLogs,\"land_journey_test.go:139: runtime-log landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/clean,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (646 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/declared-output,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-and-unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log-empty-declaration,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (664 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicResumeCompletesPublishedReleaseWithoutRepublishing,\"land_resume_test.go:55: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublishedReleaseFailureExitsIncomplete,\"land_journey_test.go:267: published release exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalKeepsControlBearingPathInOneTableRow,\"land_release_refusal_test.go:70: control refusal = (1, \\\"refused{detail=landing destination is not clean,next=reviewer: commit the destination's uncommitted work, or discard it; then bench worktree land --request 'refusal-controls' --base '5… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalListsCollidingPaths,\"land_release_refusal_test.go:52: collision refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusesPostGateUnknownIgnoredMutation,\"land_journey_test.go:175: post-gate ignored mutation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReleaseDiagnosticCannotForgeTerminalLines,\"land_flags_test.go:132: hostile release result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRemovesEveryTemporaryProspectiveArtifact/published,\"land_freshness_test.go:214: landing exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:54: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (703 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/marker,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (682 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/reconcile,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/release,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeOverlappingDestinationEdit,\"land_journey_test.go:237: last-moment overlapping edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (741 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeTrackedDestinationEdit,\"land_journey_test.go:199: last-moment tracked edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesTheCensusCountAndDropsTheRecords,\"land_census_test.go:45: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesZeroForAnAssignmentWithNoRecords,\"land_census_test.go:61: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandEffectsRowPrecedesTheLandedRecord,\"land_effects_test.go:238: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then be… (669 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandFenceRefusalOverAFoldedSourceKeepsTheCallerBase,\"land_folded_base_test.go:34: fence refusal = \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (717 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandGradesASourceCommittedByALanePass,\"land_journey_test.go:381: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandHoldsTheDestinationLease,\"land_lease_test.go:90: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writ… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandLeavesAPriorLandedAssignment,\"land_effects_cleanup_test.go:114: narrowed landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshReadsTheSeal,\"land_effects_test.go:220: unsealed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshesTheBrokerAfterPublication,\"land_effects_test.go:160: landing with a stale broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRelaysTheBoundedGreenShapeBeforeTheLandedRecord,\"land_journey_test.go:410: land exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReleasesTheLeaseOnRefusal,\"land_lease_test.go:171: red landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedCleanup,\"land_effects_cleanup_test.go:188: faulted cleanup = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedRefresh,\"land_effects_test.go:181: failed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandResumeTakesNoLease,\"land_lease_test.go:190: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (633 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsAnUnprovenSibling,\"land_effects_cleanup_test.go:140: retaining landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/ordinary,\"land_empty_sibling_test.go:45: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/resumed,\"land_empty_sibling_test.go:35: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (664 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsAFreshBroker,\"land_effects_test.go:140: landing with a fresh broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsTheRefreshWithoutBuildInputs,\"land_effects_test.go:118: landing without build inputs = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict,\"refusal_route_follow_test.go:59: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-compositi… (617 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict-pending,\"refusal_route_follow_test.go:59: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-c… (641 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-collision,\"refusal_route_follow_test.go:54: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (734 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-not-clean,\"refusal_route_follow_test.go:72: destination-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (730 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-build-window,\"refusal_route_follow_test.go:72: land-build-window re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-destination-moved,\"refusal_route_follow_test.go:54: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (737 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-handback,\"refusal_route_follow_test.go:72: land-handback re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-incomplete,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-infrastructure,\"refusal_route_follow_test.go:69: printed step \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen\\\" is not one simple command with every slot filled\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-lease-held,\"refusal_route_follow_test.go:72: land-lease-held re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-red,\"refusal_route_follow_test.go:59: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base '900… (553 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-destination-residue,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-marker,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-clean,\"refusal_route_follow_test.go:72: source-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-fenced,\"refusal_route_follow_test.go:72: source-not-fenced re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-tip-mismatch,\"refusal_route_follow_test.go:72: source-tip-mismatch re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingPrintsOutputBreakdown,\"land_census_output_test.go:25: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRedRouteNamesTheRepair,\"refusal_route_test.go:147: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base 'f1f6e08a4… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/destination-collision,\"identity_component_test.go:238: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (744 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-destination-moved,\"identity_component_test.go:238: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (745 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-incomplete,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (687 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-destination-residue,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (701 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-marker,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingReleaseFactAdapterTranslatesReleaseExit,\"land_facts_test.go:211: release-failure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRetiresAnAbsentLandedSibling,\"missing_tree_recovery_test.go:231: absent-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/dangling-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/directory,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/live-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/current-source,\"land_moved_destination_test.go:132: current-source = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/non-proof-cause,\"land_moved_destination_test.go:132: non-proof-cause = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/unreadable-ancestry,\"land_moved_destination_test.go:126: rename /dev/shm/QYZTUF/t/TestProofFaultOverACurrentSourceHandsBackunreadable-ancestry2520870877/003/.git/objects/a3/d6e9901b8628b0df10120d1cb1bebf5e7aecad.hidden /dev/shm/QYZTUF/t/TestProofFaultOverACurre… (372 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeHandsBackADetachedDestination,\"land_resume_refusal_test.go:366: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (666 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/./specs/x/spec.md,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/x,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandCompletesAnInterruptedMarker,\"land_resume_test.go:230: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (679 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandExpandsAbbreviatedIdentities,\"land_reauthorization_test.go:184: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandFollowupFailureExitsIncomplete,\"land_resume_test.go:19: first incomplete exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-path,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (719 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-state,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (722 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/lock,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/owner-marker,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/registration,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/request,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-path,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-request,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicBindsPublishedLandingIdentity,\"land_resume_test.go:101: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (669 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicCompletesAnAmendedSourceLanding,\"land_spec_amendment_test.go:75: interrupted amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandReconcilesAnUnreconciledPublishedCheckout,\"land_resume_test.go:183: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/absent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/behind,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/divergent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (687 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesNonAncestorReviewBaseWithoutMutation,\"land_resume_refusal_test.go:143: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (675 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesWhenTerminalReceiptWasEvicted,\"land_resume_refusal_test.go:204: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRepeatsTheCensusCount,\"land_resume_test.go:285: incomplete landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandSourceRefusalNamesTheCallersResume,\"land_resume_refusal_test.go:119: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandUnknownRequestNamesReauthorizeRecovery,\"land_reauthorization_test.go:16: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding,\"land_specless_test.go:151: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCompletesTheUnfinishedEffects,\"land_effects_cleanup_test.go:218: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandRunsTheEffectsAfterRelease,\"land_effects_cleanup_test.go:270: failed-refresh landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeReadsEffectStateFromTheTree,\"land_effects_test.go:269: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict,\"refusal_route_test.go:187: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict-pending,\"refusal_route_test.go:187: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/destination-collision,\"refusal_route_test.go:185: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-destination-residue,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-marker,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestSecondLandingMeetsTheLease,\"land_lease_test.go:133: first landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnclaimedLandingCompositionErrorHandsBack,\"refusal_route_test.go:164: unclaimed composition error = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (736 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnsafePathRouteUsesThePlaceholder/incomplete_resume,\"land_release_refusal_test.go:204: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (558 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowAdmitsTheBuildsOwnLanding,\"land_window_test.go:119: landing of the bound spec from the sibling = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: l… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowClosesAtCompletion,\"land_window_test.go:130: landing of the build's completion = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (682 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/build-start,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/gap-close,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (686 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t4-t4-system",
          "performer": "claude:ft293_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4",
            "digest": "sha256:c847042345da42cf27e3ba96bc15180af335713fc3a39d153ccee70b7dcede04",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,170,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,70351,170\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t4-t4-root-conformance",
          "performer": "claude:ft293_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4",
            "digest": "sha256:4bd186c2726e80ec95c3f18fd330fc4ac10f70096a355a2d5773c9d3bb480ec6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14857,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t4-t4-guidance-anchors",
          "performer": "claude:ft293_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4",
            "digest": "sha256:98c89f91dc35788825ff3170ed84aba1d23a5837c33a82fe6e90ced5add461cc",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4557,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t4-t4-guidance-budgets",
          "performer": "claude:ft293_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4",
            "digest": "sha256:0119757da3d9ed1fe105aa731f79a17727150131165456bf02d9e83c2f1730af",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
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
  },
  "amendments": [
    {
      "from": "sha256:28c270ec1eed58e47fed9261ef9ff6d80486340f37c637856ac33111d7dcf7ec",
      "to": "sha256:de5b65bac4b9d2078432cc702595272db14ab766d6242737453782ad7a94ef81",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ]
      }
    },
    {
      "from": "sha256:de5b65bac4b9d2078432cc702595272db14ab766d6242737453782ad7a94ef81",
      "to": "sha256:a557d522d03dc16e5ca4556972eeb2565c2ac5c644e5d770046705ce338f6a6c",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ]
      }
    },
    {
      "from": "sha256:a557d522d03dc16e5ca4556972eeb2565c2ac5c644e5d770046705ce338f6a6c",
      "to": "sha256:2ad84f74d4d6d1a7784989f1d50a9fd2a6980b0da96e261cc62a2701ef03d4d9",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ]
      }
    },
    {
      "from": "sha256:2ad84f74d4d6d1a7784989f1d50a9fd2a6980b0da96e261cc62a2701ef03d4d9",
      "to": "sha256:ca0df6b8c3ab10b55553f689626fe14f3729c44811ee60b2ecaff4bb6224eba5",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ]
      }
    },
    {
      "from": "sha256:ca0df6b8c3ab10b55553f689626fe14f3729c44811ee60b2ecaff4bb6224eba5",
      "to": "sha256:d9e2bdb101dd34464c7b671c63729eeb0fc18e1f835431a23963c7ef8f1bc8a2",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ],
        "CR2": [
          "CR2"
        ]
      }
    },
    {
      "from": "sha256:d9e2bdb101dd34464c7b671c63729eeb0fc18e1f835431a23963c7ef8f1bc8a2",
      "to": "sha256:85f78506160321bb5f2f0b029dcb6ffc622d852931e866e67d0eec733a187e6a",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ],
        "CR2": [
          "CR2"
        ]
      }
    },
    {
      "from": "sha256:85f78506160321bb5f2f0b029dcb6ffc622d852931e866e67d0eec733a187e6a",
      "to": "sha256:9421a83c29076fbf94debbe657d4e004e894c1ba8c9d4a336c3af0cd6712f7bc",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ],
        "CR2": [
          "CR2"
        ]
      }
    },
    {
      "from": "sha256:9421a83c29076fbf94debbe657d4e004e894c1ba8c9d4a336c3af0cd6712f7bc",
      "to": "sha256:b8aa3bef62333ca44c9dcab717464e45f80187bdce5c374ceb36ba739453eb66",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ],
        "CR2": [
          "CR2"
        ]
      }
    }
  ]
}
```
