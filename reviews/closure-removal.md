# Review outcomes

## CR4 third confirming-round pickup

The final reconciliation found that the landing deletes `roadmap/FT293.md`, which `decisions/architecture-planning.md` cites as a source. A fresh ticket 06 session moved that source to the delegate guidance, and the third confirming round read that delta. The chunk used 3 repair cycles, and the reviewer's `--auto-approve` covers the third.

The round passed on all three axes, and no finding stays open. The widening of the decision map cites the advance approval, and the orchestrator flags it for reviewer veto.

## CR4 second confirming-round pickup

The second confirming round read the repair delta `093b9a75..b0d84d87`. Each axis ran on opus at high effort. The chunk used 2 of its 2 repair cycles.

The round closed S5, S8, and S9, and no finding stays open. The gate still reds a budget overrun, the restored forbidden plan-expansion sentence, and an over-long sentence.

Optional advice:

- The operating guide keeps one gerund used as a noun.
- `/bench-write-spec` and the slicing checks keep "fence expansion" outside this delta.
- The spec's solution and story 18 keep "the directory of each signed-off entry".

## CR4 confirming-round pickup

The confirming round read the repair delta `688d194f..093b9a75`. Each axis ran on opus at high effort. The chunk has used 1 of its 2 repair cycles, and the next repair is cycle 2.

The round closed S1 to S4, S6, S7, P1, and C1. Three findings stay open, and ticket 06 owns all three.

- S5, auto-fix, confidence 6. The guidance still uses several names for one widening and for one reviewer decision. Use one name for each.
- S8, auto-fix, confidence 6. The operating guide says "the directory of each signed-off `Writes:` entry", which reads as every entry. The scope is the directory of any one signed-off entry. Correct the scope phrase, and make the repair clause conditional on that scope.
- S9, auto-fix, confidence 4. The lead-in of the decision list in the operating guide does not match its last item. Make the lead-in fit every item.

Optional advice: `/bench-implement-spec` uses a gerund as a noun in one sentence. The spec decision on where `/bench-implement-spec` names the verb differs from the tree after the S3 repair; the orchestrator flags it for reviewer veto. A code span over part of a forbid needle passes a case-folded forbid by matcher design, and the orchestrator parks it with `bench learning`.

## CR4 review pickup

The CR4 review ran on the frozen pair `ffee6a38..688d194f`. Each axis ran on opus at high effort. The chunk has used 0 of its 2 repair cycles.

The raw count is 11 findings, and the repair-target count is 9, because P2 and C2 have no repair target. Ticket 06 owns every repair target.

### Standards

Count: 7. Worst issue: S1.

- S1, auto-fix, confidence 7. The craft-review skill and the review command each define the fence as the ticket `Writes:` lines and drop the implicit entries. Point both sentences to the definition in the craft-spec skill.
- S2, auto-fix, confidence 6. A sentence in `/bench-implement-spec` has 23 words and joins an instruction and a description. Split it.
- S3, auto-fix, confidence 6. `/bench-implement-spec` restates the widen scope rule that the operating guide owns. Point to the operating guide.
- S4, auto-fix, confidence 5. The slicing checks name `/bench-implement-spec` as the widening owner, and the other files name the plan-expansion policy. Name one owner.
- S5, auto-fix, confidence 6. The guidance uses several terms for one widening and for one scope. Use one term for each.
- S6, auto-fix, confidence 5. The delegate rule uses vague nouns for a path and for the `Writes:` line. Name the path and the line.
- S7, auto-fix, confidence 3. The plan-expansion paragraph holds a four-item inline list and a garden-path sentence. Use a vertical list.

### Spec

Count: 2. Worst issue: P1.

- P1, auto-fix, confidence 6. The operating guide states only the outside half of the scope rule after the deletion. State that an in-scope widening, a repair's included, takes a plan commit and no stop.
- P2, no-op, confidence 4. The spec's story 50 exclusion and its upper-bound decision disagree. The build followed the decision.

### Coverage

Count: 2. Worst issue: C1.

- C1, auto-fix, confidence 8. The CR78 forbid anchor keeps case and emphasis, so a bold or recapitalized sentence passes. Use the case-folded emphasis kind, and add a bold case.
- C2, no-op, confidence 6. The removed-form forbid rows name only the files that the spec rows name.

## CR3 second confirming-round pickup

The second confirming round read the repair delta `9d26cb3d..ffee6a38`. Each axis ran on opus at high effort. The chunk used 2 of its 2 repair cycles, and the second cycle was its hardening cycle.

The round closed C3, C4, and P1, and no finding stays open. The CR75 section still equals the ticket union plus the pickup.

Optional advice, because each item grades only a check from the hardening cycle or is a non-blocking judgment call:

- `gatherTicketsWithPolicy` is a one-line pass-through that keeps the walk's doc.
- No test pins checkout and tracked parity for `a/b/z.md` beside `a-b/z.md`.
- No test pins the representable refusal of a tracked non-`.md` entry.

The Spec axis wrote and deleted two temporary files outside the repository against its read-only charge. The tree did not change.

## CR3 confirming-round pickup

The confirming round read the repair delta `82309770..9d26cb3d`. Each axis ran on opus at high effort. The chunk has used 1 of its 2 repair cycles, and the next repair is cycle 2.

The round closed S1 to S4 and C1, and the CR75 section still equals the ticket union plus the pickup. Three findings stay open, and ticket 04 owns all three.

- C3, auto-fix, confidence 8. No test pins that a failed read of the tracked tickets refuses the landing. A fallback to the checkout tickets survives every test, and a committed symbolic link in the tickets folder reaches that path. Add a landing case for that link.
- C4 and P1, auto-fix, confidence 6. `bench preflight review` reads the checkout tickets, so an ignored ticket-shaped file can turn its coverage and `paths-authorized` verdicts green. The review preflight grades a clean committed source, so it reads the tracked tickets of the source tip, as the landing does. The orchestrator corrects the spec sentence and flags this call for reviewer veto.
- C2, auto-fix, confidence 7. No test pins that the tracked ticket walk keeps the copy that `tickets.Enumerate` keeps when a nested duplicate sorts first. Add a parity case.

## CR3 review pickup

The CR3 review ran on the frozen pair `a883604f..82309770`. Each axis ran on opus at high effort. The chunk has used 0 of its 2 repair cycles.

The raw count is 10 findings, and the repair-target count is 7. S2 and C2 share the S1 fix, and S5, P1, and P2 have no CR3 repair target. Ticket 04 owns S1, S2, C1, C2, and C3. Ticket 05 owns S3 and S4. Ticket 06 owns P2.

### Standards

Count: 5. Worst issue: S1.

- S1, auto-fix, confidence 7. `internal/preflight/staleness/collect.go` repeats the ticket read and parse of `coversAt` in `internal/preflight/staleness/coverage.go`. Give the ticket set of a commit one source that both audit faces read.
- S2, auto-fix, confidence 4. The audit's ticket set does not apply the duplicate-basename rule of `tickets.Enumerate`. The S1 fix closes it.
- S3, auto-fix, confidence 5. The craft-spec skill says "plus the implicit entries" and does not name them. Name `capture`, the spec folder, and the review pickup.
- S4, auto-fix, confidence 4. The craft-spec skill says "Each fence" after the build leaves one derived fence. Correct the term.
- S5, no-op, confidence 3. A rename of `fence_writes.go` is outside the acceptance.

### Spec

Count: 2. Worst issue: P2.

- P1, no-op, confidence 3. The `fence` block lists the ticket union, as the spec decision states.
- P2, no-op for CR3, confidence 8. The craft-review skill still says that the ownership fence lives in craft-spec. Ticket 06 owns that sentence.

The section of this spec covers the whole ticket union plus `reviews/closure-removal.md`, with no gap (CR75).

### Coverage

Count: 3. Worst issue: C1.

- C1, auto-fix, confidence 8. No test pins that a `(deleted)` entry authorizes its own deletion. Add a landing or build case.
- C2, auto-fix, confidence 7. The audit scope drops a ticket under `tickets/sub/`, but the gather authorizes it. The S1 fix closes it, and a nested case pins it.
- C3, auto-fix, confidence 6. The landing gathers its tickets from the filesystem, so a git-ignored ticket-shaped file can authorize paths. The landing fence must read only the tracked tickets at the source tip. The orchestrator flags this call for reviewer veto.

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
  "plan_digest": "sha256:40e3bcd3128efe554b1c3045d1e51f4107e10d7f6f8c9590d56d82e1cb1420f9",
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
      "tip": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
      "plan_digest": "sha256:6c7a65800371afdd0e0250f01b0d90b912a017eb591baca4f30a76d7407ee411",
      "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
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
        },
        {
          "id": "t5r1-t5-preflight-packages",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:f57cb4b7d16c05a8791c20aee6919751c25fe7095d6f1d0fb63d79c632579e9c",
            "excerpt": "packages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,pass,7217,101\n  github.com/gibbonmi/bench/internal/preflight,pass,28767,366\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,21172,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,11952,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5520,70\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/preflighttest --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd --package ./internal/diff",
          "exit_code": 0
        },
        {
          "id": "t5r1-t5-fence-readers",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:1bcea9a15d8856ecd4263266946ff637c7c7ce4c3d4b92f4f50c23073f8ac9d4",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/coverage/citations.go,insert,failed,1,TestCommandCheckIgnoresAFenceSection,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/coverage ./internal/spec ./internal/chargeevidence ./internal/refusalroute ./internal/reviewrecord/recordtest ./internal/anchors ./cmd/bench ./internal/conformance,all,passed,2025\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1333\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18317,346\n  github.com/gibbonmi/bench/internal/anchors,pass,125,123\n  github.com/gibbonmi/bench/internal/chargeevidence,pass,25,122\n  github.com/gibbonmi/bench/internal/conformance,pass,36664,1116\n  github.com/gibbonmi/bench/internal/coverage,fail,3306,141\n  github.com/gibbonmi/bench/internal/refusalroute,pass,7,20\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/spec,pass,2396,157\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/coverage,TestCommandCheckIgnoresAFenceSection,\"fence_section_test.go:13: Command = (1, \\\"error: specs/s/spec.md specs/s/spec.md coverage map spec declares ownership fences without the review pickup — fix the map or mark it <!-- coverage-map: historical -->\\\\\\\\n\\\"), want exit 0 for a fence se… (273 bytes)\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/chargeevidence,TestEvidenceStoreKinds/CE94_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
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
              "ref": "claude-agent:ft293_t5_r1",
              "digest": "sha256:1bcea9a15d8856ecd4263266946ff637c7c7ce4c3d4b92f4f50c23073f8ac9d4",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/coverage/citations.go,insert,failed,1,TestCommandCheckIgnoresAFenceSection,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/coverage ./internal/spec ./internal/chargeevidence ./internal/refusalroute ./internal/reviewrecord/recordtest ./internal/anchors ./cmd/bench ./internal/conformance,all,passed,2025\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1333\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18317,346\n  github.com/gibbonmi/bench/internal/anchors,pass,125,123\n  github.com/gibbonmi/bench/internal/chargeevidence,pass,25,122\n  github.com/gibbonmi/bench/internal/conformance,pass,36664,1116\n  github.com/gibbonmi/bench/internal/coverage,fail,3306,141\n  github.com/gibbonmi/bench/internal/refusalroute,pass,7,20\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/spec,pass,2396,157\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/coverage,TestCommandCheckIgnoresAFenceSection,\"fence_section_test.go:13: Command = (1, \\\"error: specs/s/spec.md specs/s/spec.md coverage map spec declares ownership fences without the review pickup — fix the map or mark it <!-- coverage-map: historical -->\\\\\\\\n\\\"), want exit 0 for a fence se… (273 bytes)\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/chargeevidence,TestEvidenceStoreKinds/CE94_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t5r1-t5-landing",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:2e3770b09beb590b14e37ad090f94ac702dc368c75278a12012604044e572d38",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,66609,1438\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-landing",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t5r1-t5-system",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:7804b2f72b7ec016142ccd0997cd2b2df442ceb1b6de2d923ba8db0e26485cbb",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,170,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,70616,170\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t5r1-t5-root-conformance",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:fab75ed76732489a1f43101e9b17deb4dac7d6cadf2e59c79260a1f0e6355d6e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14563,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t5r1-t5-guidance-anchors",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:584985bd8b4686c4729c34e787d5682f065d2b02761ed8b55a3007e88f73ade6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4934,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t5r1-t5-guidance-budgets",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:b682aaa6a2f3a10908927cd53ce11e0bf9cf598d196444453baf84591b6728d5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t4r1-t4-preflight-packages",
          "performer": "claude:ft293_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r1",
            "digest": "sha256:44d4fbf62e39d1ea321f5ff128e3e0520f8ce6023b7d9027d967ac8414130858",
            "excerpt": "packages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,27919,366\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,19730,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,11487,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5532,70\n  github.com/gibbonmi/bench/internal/tickets,pass,3,45\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/preflighttest --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "t4r1-t4-fence-readers",
          "performer": "claude:ft293_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r1",
            "digest": "sha256:174c2ce4cf60f3abf4f9c8412a864b0c549a8ffc024368b32430b04f68aeb6c0",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16937,346\n  github.com/gibbonmi/bench/internal/anchors,pass,107,123\n  github.com/gibbonmi/bench/internal/chargeevidence,pass,23,122\n  github.com/gibbonmi/bench/internal/conformance,pass,36737,1116\n  github.com/gibbonmi/bench/internal/coverage,pass,3298,141\n  github.com/gibbonmi/bench/internal/refusalroute,pass,5,20\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/spec,pass,2286,157\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/chargeevidence,TestEvidenceStoreKinds/CE94_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-fence-readers",
          "command": "bench test --package ./internal/coverage --package ./internal/spec --package ./internal/chargeevidence --package ./internal/refusalroute --package ./internal/reviewrecord/recordtest --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t4r1-t4-landing",
          "performer": "claude:ft293_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r1",
            "digest": "sha256:051a81c73636246ca82591c84cad170c5cc55c35b691131c7124f6efeba239d2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/fence_writes.go,swap,failed,198,TestCommitmentClosureAdmission/rowless-FT2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1433\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,71131\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,600037,1433\nfailures[198]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/rowless-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (711 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (722 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT3,\"commitment_landing_test.go:101: closing FT3 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/rowless,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/spec,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/rowless,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/spec,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentDeliveryClosure,\"commitment_landing_test.go:27: delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/blocker,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/default-branch-policy,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentInterruptedGate,\"commitment_landing_test.go:143: interrupted closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLandingRefusesObligationFreeDelivery,\"commitment_light_landing_test.go:327: obligation-free landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/partly-delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/beyond-scope,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (753 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/directory-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/listed-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/outside-directory,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (762 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/scope-without-deliverable,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (799 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/sibling-prefix,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (755 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/unlisted,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (752 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-then-tickets,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-with-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-without-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/tickets-then-spec,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (721 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPartialDelivery,\"commitment_landing_test.go:56: partial landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPersistenceBeforeGate,\"commitment_landing_test.go:173: torn-intent landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPublishLock,\"land_identity_test.go:307: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then b… (660 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentStaleLanding,\"land_identity_test.go:209: displaced landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCompareAndSwapRefusalNamesTheFold,\"land_moved_destination_test.go:66: moved destination = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold,\"land_moved_destination_test.go:43: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (733 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:132: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict-pending,\"refusal_route_test.go:132: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestDelegatedIntegrationJourney,\"delegated_integration_test.go:237: delegated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (652 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh,\"panic: test timed out after 10m0s\",135\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/installed,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (741 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/kit-checkout,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (747 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCleansTheFoldedSibling,\"land_effects_cleanup_test.go:83: folded-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAcceptsAbbreviatedIdentities,\"land_surface_test.go:46: abbreviated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthenticatesDigestShapedRequestToken,\"land_reauthorization_test.go:102: digest-shaped request land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (756 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesCaptureOutsideTheFence,\"land_surface_test.go:135: capture landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (676 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesFromTicketWrites,\"land_derived_fence_test.go:34: ticket-only landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (679 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/another_review_record,\"land_derived_fence_test.go:92: other review landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (742 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/pickup,\"land_derived_fence_test.go:78: pickup landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (679 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandCarriesTheBaselineScheduleRootIntoTheProspectiveGate,\"land_freshness_test.go:288: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (700 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandComposesCaptureOntoMovedDestination,\"land_surface_test.go:71: unreviewed capture composition = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictNextPointsThroughUnsafePath,\"land_surface_test.go:360: unsafe-path conflict next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesThePath,\"land_surface_test.go:268: conflict refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesTheSourceRepair,\"land_surface_test.go:301: conflict repair next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandDisclosesAUnionResolution,\"land_surface_test.go:104: unreviewed union = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (658 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedBase,\"land_reauthorization_test.go:175: abbreviated base = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"ecf122deb87d\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/39,\"land_reauthorization_test.go:163: abbreviated source tip \\\"7f5185f69bc23f06709f03111ecb19940f629e5\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or a… (752 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/4,\"land_reauthorization_test.go:163: abbreviated source tip \\\"4415\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line … (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/upper-12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"0B8CA7C170D0\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (729 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandFenceRefusalNamesThePath,\"land_surface_test.go:220: fence refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; t… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal,\"land_freshness_test.go:114: forged-primary landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (699 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIncompleteNextUsesAssignmentPointerForUnsafePath,\"land_resume_test.go:38: unsafe-path incomplete = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease,\"land_freshness_test.go:72: single-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandLeavesTheDestinationUnchangedAfterARedProspectiveGate,\"land_freshness_test.go:185: gate tally = \\\"\\\", open /dev/shm/JYYKGJ/t/TestLandCommandLeavesTheDestinationUnchangedAfterARedProspective3134160856/003/.git/bench-land-gate-tally: no such file or directory, want one prospective run\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout,\"land_broker_notice_test.go:125: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (720 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion,\"land_freshness_test.go:39: stable-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/marker,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (644 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/reconcile,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (650 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/release,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (646 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsNoHeadsLineWithoutRecords,\"land_census_test.go:91: empty census landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsTheCensusHeadBreakdown,\"land_census_test.go:76: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/absent,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (646 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/concurrently-moved,\"land_flags_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/present,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (647 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicConflictRepairRequiresNewReviewedTip,\"land_journey_test.go:294: conflict result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicLandsAnInRangeSpecAmendment,\"land_spec_amendment_test.go:47: amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicPreservesHistoricalRuntimeLogs,\"land_journey_test.go:139: runtime-log landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/clean,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (646 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/declared-output,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-and-unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (658 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log-empty-declaration,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicResumeCompletesPublishedReleaseWithoutRepublishing,\"land_resume_test.go:55: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublishedReleaseFailureExitsIncomplete,\"land_journey_test.go:267: published release exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReadsOnlyTrackedTickets,\"land_derived_fence_test.go:61: ignored-ticket landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (723 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalKeepsControlBearingPathInOneTableRow,\"land_release_refusal_test.go:70: control refusal = (1, \\\"refused{detail=landing destination is not clean,next=reviewer: commit the destination's uncommitted work, or discard it; then bench worktree land --request 'refusal-controls' --base '0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalListsCollidingPaths,\"land_release_refusal_test.go:52: collision refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusesPostGateUnknownIgnoredMutation,\"land_journey_test.go:175: post-gate ignored mutation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReleaseDiagnosticCannotForgeTerminalLines,\"land_flags_test.go:132: hostile release result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRemovesEveryTemporaryProspectiveArtifact/published,\"land_freshness_test.go:214: landing exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:54: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/marker,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/reconcile,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/release,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeOverlappingDestinationEdit,\"land_journey_test.go:237: last-moment overlapping edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (742 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeTrackedDestinationEdit,\"land_journey_test.go:199: last-moment tracked edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesTheCensusCountAndDropsTheRecords,\"land_census_test.go:45: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesZeroForAnAssignmentWithNoRecords,\"land_census_test.go:61: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandEffectsRowPrecedesTheLandedRecord,\"land_effects_test.go:238: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then be… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandFenceRefusalOverAFoldedSourceKeepsTheCallerBase,\"land_folded_base_test.go:34: fence refusal = \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (717 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandGradesASourceCommittedByALanePass,\"land_journey_test.go:381: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandHoldsTheDestinationLease,\"land_lease_test.go:90: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writ… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandLeavesAPriorLandedAssignment,\"land_effects_cleanup_test.go:114: narrowed landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshReadsTheSeal,\"land_effects_test.go:220: unsealed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshesTheBrokerAfterPublication,\"land_effects_test.go:160: landing with a stale broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRelaysTheBoundedGreenShapeBeforeTheLandedRecord,\"land_journey_test.go:410: land exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReleasesTheLeaseOnRefusal,\"land_lease_test.go:171: red landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedCleanup,\"land_effects_cleanup_test.go:188: faulted cleanup = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedRefresh,\"land_effects_test.go:181: failed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandResumeTakesNoLease,\"land_lease_test.go:190: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (634 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsAnUnprovenSibling,\"land_effects_cleanup_test.go:140: retaining landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/ordinary,\"land_empty_sibling_test.go:45: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (654 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/resumed,\"land_empty_sibling_test.go:35: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (664 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsAFreshBroker,\"land_effects_test.go:140: landing with a fresh broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsTheRefreshWithoutBuildInputs,\"land_effects_test.go:118: landing without build inputs = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict,\"refusal_route_follow_test.go:59: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-compositi… (617 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict-pending,\"refusal_route_follow_test.go:59: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-c… (641 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-collision,\"refusal_route_follow_test.go:54: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-not-clean,\"refusal_route_follow_test.go:72: destination-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (729 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-build-window,\"refusal_route_follow_test.go:72: land-build-window re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-destination-moved,\"refusal_route_follow_test.go:54: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (739 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-handback,\"refusal_route_follow_test.go:72: land-handback re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-incomplete,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (675 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-infrastructure,\"refusal_route_follow_test.go:69: printed step \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen\\\" is not one simple command with every slot filled\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-lease-held,\"refusal_route_follow_test.go:72: land-lease-held re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-red,\"refusal_route_follow_test.go:59: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base '570… (553 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-destination-residue,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-marker,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-clean,\"refusal_route_follow_test.go:72: source-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-fenced,\"refusal_route_follow_test.go:72: source-not-fenced re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (717 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-tip-mismatch,\"refusal_route_follow_test.go:72: source-tip-mismatch re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (722 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingPrintsOutputBreakdown,\"land_census_output_test.go:25: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRedRouteNamesTheRepair,\"refusal_route_test.go:147: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base 'a6586c291… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/destination-collision,\"identity_component_test.go:238: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (744 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-destination-moved,\"identity_component_test.go:238: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (746 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-incomplete,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (687 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-destination-residue,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (701 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-marker,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingReleaseFactAdapterTranslatesReleaseExit,\"land_facts_test.go:211: release-failure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRetiresAnAbsentLandedSibling,\"missing_tree_recovery_test.go:231: absent-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/dangling-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/directory,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (700 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/live-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (700 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/current-source,\"land_moved_destination_test.go:132: current-source = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/non-proof-cause,\"land_moved_destination_test.go:132: non-proof-cause = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/unreadable-ancestry,\"land_moved_destination_test.go:126: rename /dev/shm/JYYKGJ/t/TestProofFaultOverACurrentSourceHandsBackunreadable-ancestry608295967/003/.git/objects/64/8e9e78b9d7508662f7de72aa96ae7bae00d449.hidden /dev/shm/JYYKGJ/t/TestProofFaultOverACurren… (370 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeHandsBackADetachedDestination,\"land_resume_refusal_test.go:366: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (666 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/./specs/x/spec.md,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/x,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandCompletesAnInterruptedMarker,\"land_resume_test.go:230: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandExpandsAbbreviatedIdentities,\"land_reauthorization_test.go:184: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandFollowupFailureExitsIncomplete,\"land_resume_test.go:19: first incomplete exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-path,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (719 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-state,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (721 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/lock,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/owner-marker,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/registration,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/request,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (703 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-path,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (666 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-request,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (669 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicBindsPublishedLandingIdentity,\"land_resume_test.go:101: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicCompletesAnAmendedSourceLanding,\"land_spec_amendment_test.go:75: interrupted amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandReconcilesAnUnreconciledPublishedCheckout,\"land_resume_test.go:183: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/absent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/behind,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/divergent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesNonAncestorReviewBaseWithoutMutation,\"land_resume_refusal_test.go:143: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (675 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesWhenTerminalReceiptWasEvicted,\"land_resume_refusal_test.go:204: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRepeatsTheCensusCount,\"land_resume_test.go:285: incomplete landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (676 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandSourceRefusalNamesTheCallersResume,\"land_resume_refusal_test.go:119: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandUnknownRequestNamesReauthorizeRecovery,\"land_reauthorization_test.go:16: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding,\"land_specless_test.go:151: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCompletesTheUnfinishedEffects,\"land_effects_cleanup_test.go:218: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandRunsTheEffectsAfterRelease,\"land_effects_cleanup_test.go:270: failed-refresh landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeReadsEffectStateFromTheTree,\"land_effects_test.go:269: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict,\"refusal_route_test.go:187: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict-pending,\"refusal_route_test.go:187: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/destination-collision,\"refusal_route_test.go:185: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-destination-residue,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-marker,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestSecondLandingMeetsTheLease,\"land_lease_test.go:133: first landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnclaimedLandingCompositionErrorHandsBack,\"refusal_route_test.go:164: unclaimed composition error = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (736 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnsafePathRouteUsesThePlaceholder/incomplete_resume,\"land_release_refusal_test.go:204: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (558 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowAdmitsTheBuildsOwnLanding,\"land_window_test.go:119: landing of the bound spec from the sibling = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: l… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowClosesAtCompletion,\"land_window_test.go:130: landing of the build's completion = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/build-start,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/gap-close,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (686 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:ft293_t4_r1",
              "digest": "sha256:051a81c73636246ca82591c84cad170c5cc55c35b691131c7124f6efeba239d2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/fence_writes.go,swap,failed,198,TestCommitmentClosureAdmission/rowless-FT2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1433\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,71131\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,600037,1433\nfailures[198]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/rowless-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (711 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (722 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT3,\"commitment_landing_test.go:101: closing FT3 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/rowless,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/spec,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/rowless,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/spec,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentDeliveryClosure,\"commitment_landing_test.go:27: delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/blocker,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/default-branch-policy,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentInterruptedGate,\"commitment_landing_test.go:143: interrupted closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLandingRefusesObligationFreeDelivery,\"commitment_light_landing_test.go:327: obligation-free landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/partly-delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/beyond-scope,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (753 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/directory-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/listed-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/outside-directory,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (762 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/scope-without-deliverable,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (799 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/sibling-prefix,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (755 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/unlisted,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (752 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-then-tickets,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-with-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-without-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/tickets-then-spec,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (721 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPartialDelivery,\"commitment_landing_test.go:56: partial landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPersistenceBeforeGate,\"commitment_landing_test.go:173: torn-intent landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPublishLock,\"land_identity_test.go:307: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then b… (660 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentStaleLanding,\"land_identity_test.go:209: displaced landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCompareAndSwapRefusalNamesTheFold,\"land_moved_destination_test.go:66: moved destination = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold,\"land_moved_destination_test.go:43: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (733 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:132: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict-pending,\"refusal_route_test.go:132: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestDelegatedIntegrationJourney,\"delegated_integration_test.go:237: delegated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (652 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh,\"panic: test timed out after 10m0s\",135\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/installed,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (741 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/kit-checkout,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (747 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCleansTheFoldedSibling,\"land_effects_cleanup_test.go:83: folded-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAcceptsAbbreviatedIdentities,\"land_surface_test.go:46: abbreviated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthenticatesDigestShapedRequestToken,\"land_reauthorization_test.go:102: digest-shaped request land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (756 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesCaptureOutsideTheFence,\"land_surface_test.go:135: capture landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (676 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesFromTicketWrites,\"land_derived_fence_test.go:34: ticket-only landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (679 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/another_review_record,\"land_derived_fence_test.go:92: other review landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (742 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/pickup,\"land_derived_fence_test.go:78: pickup landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (679 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandCarriesTheBaselineScheduleRootIntoTheProspectiveGate,\"land_freshness_test.go:288: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (700 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandComposesCaptureOntoMovedDestination,\"land_surface_test.go:71: unreviewed capture composition = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictNextPointsThroughUnsafePath,\"land_surface_test.go:360: unsafe-path conflict next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesThePath,\"land_surface_test.go:268: conflict refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesTheSourceRepair,\"land_surface_test.go:301: conflict repair next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandDisclosesAUnionResolution,\"land_surface_test.go:104: unreviewed union = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (658 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedBase,\"land_reauthorization_test.go:175: abbreviated base = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"ecf122deb87d\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/39,\"land_reauthorization_test.go:163: abbreviated source tip \\\"7f5185f69bc23f06709f03111ecb19940f629e5\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or a… (752 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/4,\"land_reauthorization_test.go:163: abbreviated source tip \\\"4415\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line … (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/upper-12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"0B8CA7C170D0\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (729 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandFenceRefusalNamesThePath,\"land_surface_test.go:220: fence refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; t… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal,\"land_freshness_test.go:114: forged-primary landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (699 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIncompleteNextUsesAssignmentPointerForUnsafePath,\"land_resume_test.go:38: unsafe-path incomplete = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease,\"land_freshness_test.go:72: single-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandLeavesTheDestinationUnchangedAfterARedProspectiveGate,\"land_freshness_test.go:185: gate tally = \\\"\\\", open /dev/shm/JYYKGJ/t/TestLandCommandLeavesTheDestinationUnchangedAfterARedProspective3134160856/003/.git/bench-land-gate-tally: no such file or directory, want one prospective run\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout,\"land_broker_notice_test.go:125: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (720 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion,\"land_freshness_test.go:39: stable-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/marker,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (644 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/reconcile,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (650 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/release,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (646 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsNoHeadsLineWithoutRecords,\"land_census_test.go:91: empty census landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsTheCensusHeadBreakdown,\"land_census_test.go:76: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/absent,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (646 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/concurrently-moved,\"land_flags_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/present,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (647 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicConflictRepairRequiresNewReviewedTip,\"land_journey_test.go:294: conflict result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicLandsAnInRangeSpecAmendment,\"land_spec_amendment_test.go:47: amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicPreservesHistoricalRuntimeLogs,\"land_journey_test.go:139: runtime-log landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/clean,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (646 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/declared-output,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-and-unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (658 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log-empty-declaration,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicResumeCompletesPublishedReleaseWithoutRepublishing,\"land_resume_test.go:55: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublishedReleaseFailureExitsIncomplete,\"land_journey_test.go:267: published release exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReadsOnlyTrackedTickets,\"land_derived_fence_test.go:61: ignored-ticket landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (723 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalKeepsControlBearingPathInOneTableRow,\"land_release_refusal_test.go:70: control refusal = (1, \\\"refused{detail=landing destination is not clean,next=reviewer: commit the destination's uncommitted work, or discard it; then bench worktree land --request 'refusal-controls' --base '0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalListsCollidingPaths,\"land_release_refusal_test.go:52: collision refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusesPostGateUnknownIgnoredMutation,\"land_journey_test.go:175: post-gate ignored mutation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReleaseDiagnosticCannotForgeTerminalLines,\"land_flags_test.go:132: hostile release result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRemovesEveryTemporaryProspectiveArtifact/published,\"land_freshness_test.go:214: landing exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:54: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/marker,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/reconcile,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/release,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeOverlappingDestinationEdit,\"land_journey_test.go:237: last-moment overlapping edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (742 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeTrackedDestinationEdit,\"land_journey_test.go:199: last-moment tracked edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesTheCensusCountAndDropsTheRecords,\"land_census_test.go:45: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesZeroForAnAssignmentWithNoRecords,\"land_census_test.go:61: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandEffectsRowPrecedesTheLandedRecord,\"land_effects_test.go:238: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then be… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandFenceRefusalOverAFoldedSourceKeepsTheCallerBase,\"land_folded_base_test.go:34: fence refusal = \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (717 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandGradesASourceCommittedByALanePass,\"land_journey_test.go:381: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandHoldsTheDestinationLease,\"land_lease_test.go:90: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writ… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandLeavesAPriorLandedAssignment,\"land_effects_cleanup_test.go:114: narrowed landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshReadsTheSeal,\"land_effects_test.go:220: unsealed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshesTheBrokerAfterPublication,\"land_effects_test.go:160: landing with a stale broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRelaysTheBoundedGreenShapeBeforeTheLandedRecord,\"land_journey_test.go:410: land exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReleasesTheLeaseOnRefusal,\"land_lease_test.go:171: red landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedCleanup,\"land_effects_cleanup_test.go:188: faulted cleanup = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedRefresh,\"land_effects_test.go:181: failed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandResumeTakesNoLease,\"land_lease_test.go:190: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (634 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsAnUnprovenSibling,\"land_effects_cleanup_test.go:140: retaining landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/ordinary,\"land_empty_sibling_test.go:45: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (654 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/resumed,\"land_empty_sibling_test.go:35: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (664 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsAFreshBroker,\"land_effects_test.go:140: landing with a fresh broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsTheRefreshWithoutBuildInputs,\"land_effects_test.go:118: landing without build inputs = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict,\"refusal_route_follow_test.go:59: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-compositi… (617 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict-pending,\"refusal_route_follow_test.go:59: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-c… (641 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-collision,\"refusal_route_follow_test.go:54: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-not-clean,\"refusal_route_follow_test.go:72: destination-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (729 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-build-window,\"refusal_route_follow_test.go:72: land-build-window re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-destination-moved,\"refusal_route_follow_test.go:54: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (739 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-handback,\"refusal_route_follow_test.go:72: land-handback re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-incomplete,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (675 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-infrastructure,\"refusal_route_follow_test.go:69: printed step \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen\\\" is not one simple command with every slot filled\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-lease-held,\"refusal_route_follow_test.go:72: land-lease-held re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-red,\"refusal_route_follow_test.go:59: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base '570… (553 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-destination-residue,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-marker,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-clean,\"refusal_route_follow_test.go:72: source-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-fenced,\"refusal_route_follow_test.go:72: source-not-fenced re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (717 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-tip-mismatch,\"refusal_route_follow_test.go:72: source-tip-mismatch re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (722 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingPrintsOutputBreakdown,\"land_census_output_test.go:25: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRedRouteNamesTheRepair,\"refusal_route_test.go:147: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base 'a6586c291… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/destination-collision,\"identity_component_test.go:238: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (744 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-destination-moved,\"identity_component_test.go:238: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (746 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-incomplete,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (687 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-destination-residue,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (701 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-marker,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingReleaseFactAdapterTranslatesReleaseExit,\"land_facts_test.go:211: release-failure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRetiresAnAbsentLandedSibling,\"missing_tree_recovery_test.go:231: absent-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/dangling-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/directory,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (700 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/live-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (700 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/current-source,\"land_moved_destination_test.go:132: current-source = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/non-proof-cause,\"land_moved_destination_test.go:132: non-proof-cause = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/unreadable-ancestry,\"land_moved_destination_test.go:126: rename /dev/shm/JYYKGJ/t/TestProofFaultOverACurrentSourceHandsBackunreadable-ancestry608295967/003/.git/objects/64/8e9e78b9d7508662f7de72aa96ae7bae00d449.hidden /dev/shm/JYYKGJ/t/TestProofFaultOverACurren… (370 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeHandsBackADetachedDestination,\"land_resume_refusal_test.go:366: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (666 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/./specs/x/spec.md,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/x,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandCompletesAnInterruptedMarker,\"land_resume_test.go:230: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandExpandsAbbreviatedIdentities,\"land_reauthorization_test.go:184: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandFollowupFailureExitsIncomplete,\"land_resume_test.go:19: first incomplete exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-path,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (719 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-state,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (721 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/lock,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/owner-marker,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/registration,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/request,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (703 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-path,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (666 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-request,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (669 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicBindsPublishedLandingIdentity,\"land_resume_test.go:101: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicCompletesAnAmendedSourceLanding,\"land_spec_amendment_test.go:75: interrupted amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandReconcilesAnUnreconciledPublishedCheckout,\"land_resume_test.go:183: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/absent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/behind,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/divergent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesNonAncestorReviewBaseWithoutMutation,\"land_resume_refusal_test.go:143: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (675 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesWhenTerminalReceiptWasEvicted,\"land_resume_refusal_test.go:204: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRepeatsTheCensusCount,\"land_resume_test.go:285: incomplete landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (676 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandSourceRefusalNamesTheCallersResume,\"land_resume_refusal_test.go:119: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandUnknownRequestNamesReauthorizeRecovery,\"land_reauthorization_test.go:16: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding,\"land_specless_test.go:151: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCompletesTheUnfinishedEffects,\"land_effects_cleanup_test.go:218: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandRunsTheEffectsAfterRelease,\"land_effects_cleanup_test.go:270: failed-refresh landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeReadsEffectStateFromTheTree,\"land_effects_test.go:269: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict,\"refusal_route_test.go:187: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict-pending,\"refusal_route_test.go:187: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/destination-collision,\"refusal_route_test.go:185: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-destination-residue,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-marker,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestSecondLandingMeetsTheLease,\"land_lease_test.go:133: first landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnclaimedLandingCompositionErrorHandsBack,\"refusal_route_test.go:164: unclaimed composition error = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (736 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnsafePathRouteUsesThePlaceholder/incomplete_resume,\"land_release_refusal_test.go:204: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (558 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowAdmitsTheBuildsOwnLanding,\"land_window_test.go:119: landing of the bound spec from the sibling = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: l… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowClosesAtCompletion,\"land_window_test.go:130: landing of the build's completion = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/build-start,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/gap-close,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (686 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t4r1-t4-system",
          "performer": "claude:ft293_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r1",
            "digest": "sha256:8094e5ed95bda585d93e77729f248a61e6e11ff7730124a735c26eb464f51879",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,170,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,69903,170\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t4r1-t4-root-conformance",
          "performer": "claude:ft293_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r1",
            "digest": "sha256:287c247d507651e8a83755280e6c559e65191e4eed9b254dc53eb5407531c059",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15078,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t4r1-t4-guidance-anchors",
          "performer": "claude:ft293_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r1",
            "digest": "sha256:a3048ae5bec0d5959453beae03393398af5a4cb898dee6ba0a89a1ea1a7369e5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4364,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t4r1-t4-guidance-budgets",
          "performer": "claude:ft293_t4_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r1",
            "digest": "sha256:0119757da3d9ed1fe105aa731f79a17727150131165456bf02d9e83c2f1730af",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t5r2-t5-preflight-packages",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:62de9b396a63196cfe940e6d2fdaf94315f7201af75729ae6eada4c28a537320",
            "excerpt": "packages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,pass,6906,101\n  github.com/gibbonmi/bench/internal/preflight,pass,28818,369\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,21051,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,12264,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5420,70\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/preflighttest --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd --package ./internal/diff",
          "exit_code": 0
        },
        {
          "id": "t5r2-t5-fence-readers",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:8bf05848370be2d4f9dae3bcb48b2700b1011f8e62ebeb5ac57658461011a239",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/coverage/citations.go,insert,failed,1,TestCommandCheckIgnoresAFenceSection,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/coverage ./internal/spec ./internal/chargeevidence ./internal/refusalroute ./internal/reviewrecord/recordtest ./internal/anchors ./cmd/bench ./internal/conformance,all,passed,2025\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1333\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18227,346\n  github.com/gibbonmi/bench/internal/anchors,pass,121,123\n  github.com/gibbonmi/bench/internal/chargeevidence,pass,22,122\n  github.com/gibbonmi/bench/internal/conformance,pass,36015,1116\n  github.com/gibbonmi/bench/internal/coverage,fail,3520,141\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,20\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/spec,pass,2429,157\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/coverage,TestCommandCheckIgnoresAFenceSection,\"fence_section_test.go:13: Command = (1, \\\"error: specs/s/spec.md specs/s/spec.md coverage map spec declares ownership fences without the review pickup — fix the map or mark it <!-- coverage-map: historical -->\\\\\\\\n\\\"), want exit 0 for a fence se… (273 bytes)\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/chargeevidence,TestEvidenceStoreKinds/CE94_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
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
              "ref": "claude-agent:ft293_t5_r1",
              "digest": "sha256:8bf05848370be2d4f9dae3bcb48b2700b1011f8e62ebeb5ac57658461011a239",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/coverage/citations.go,insert,failed,1,TestCommandCheckIgnoresAFenceSection,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/coverage ./internal/spec ./internal/chargeevidence ./internal/refusalroute ./internal/reviewrecord/recordtest ./internal/anchors ./cmd/bench ./internal/conformance,all,passed,2025\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1333\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18227,346\n  github.com/gibbonmi/bench/internal/anchors,pass,121,123\n  github.com/gibbonmi/bench/internal/chargeevidence,pass,22,122\n  github.com/gibbonmi/bench/internal/conformance,pass,36015,1116\n  github.com/gibbonmi/bench/internal/coverage,fail,3520,141\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,20\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/spec,pass,2429,157\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/coverage,TestCommandCheckIgnoresAFenceSection,\"fence_section_test.go:13: Command = (1, \\\"error: specs/s/spec.md specs/s/spec.md coverage map spec declares ownership fences without the review pickup — fix the map or mark it <!-- coverage-map: historical -->\\\\\\\\n\\\"), want exit 0 for a fence se… (273 bytes)\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/chargeevidence,TestEvidenceStoreKinds/CE94_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        },
        {
          "id": "t5r2-t5-landing",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:76a5116a7644b86235490a5cf0d4fd27222d2c1bbb8b4a069fb39ec08475410d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,64977,1439\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-landing",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t5r2-t5-system",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:e8e32ba257fc697d8b2d7dc5c29439223e2889a5a47d5aa506652bb5fdb2ab54",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,170,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,69138,170\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t5r2-t5-root-conformance",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:011fdc3c0e120b3fad28f319acf86df77d1706aa4f48a236f75806d365578633",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14889,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t5r2-t5-guidance-anchors",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:7d9ad523080d605ae2503e8debbbd4f14c8837300bcfa1b10a928fd470095de7",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4350,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t5r2-t5-guidance-budgets",
          "performer": "claude:ft293_t5_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t5_r1",
            "digest": "sha256:0119757da3d9ed1fe105aa731f79a17727150131165456bf02d9e83c2f1730af",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t4r2-t4-preflight-packages",
          "performer": "claude:ft293_t4_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r2",
            "digest": "sha256:fd7f16e120d3ba8248db143ba0d18f0bbe0a1114a4b2c9671b04c3b75d157ca1",
            "excerpt": "packages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,28101,369\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,19164,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,10684,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5399,70\n  github.com/gibbonmi/bench/internal/tickets,pass,3,45\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-preflight-packages",
          "command": "bench test --package ./internal/preflight --package ./internal/preflight/preflighttest --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "t4r2-t4-fence-readers",
          "performer": "claude:ft293_t4_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r2",
            "digest": "sha256:2144f5f64c071188963938c0273e8df441c005c5aee2275cbe07c67dc37f2f65",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16511,346\n  github.com/gibbonmi/bench/internal/anchors,pass,106,123\n  github.com/gibbonmi/bench/internal/chargeevidence,pass,25,122\n  github.com/gibbonmi/bench/internal/conformance,pass,35642,1116\n  github.com/gibbonmi/bench/internal/coverage,pass,3208,141\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,20\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/spec,pass,2288,157\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/chargeevidence,TestEvidenceStoreKinds/CE94_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-fence-readers",
          "command": "bench test --package ./internal/coverage --package ./internal/spec --package ./internal/chargeevidence --package ./internal/refusalroute --package ./internal/reviewrecord/recordtest --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t4r2-t4-landing",
          "performer": "claude:ft293_t4_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r2",
            "digest": "sha256:f090cf34dd5f68b785a1feef9675178b49165f935040f97cda30db5cf8681358",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/fence_writes.go,swap,failed,198,TestCommitmentClosureAdmission/rowless-FT2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1434\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,71136\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,600012,1434\nfailures[198]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/rowless-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (711 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT3,\"commitment_landing_test.go:101: closing FT3 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/rowless,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/spec,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/rowless,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (702 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/spec,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentDeliveryClosure,\"commitment_landing_test.go:27: delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/blocker,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/default-branch-policy,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentInterruptedGate,\"commitment_landing_test.go:143: interrupted closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLandingRefusesObligationFreeDelivery,\"commitment_light_landing_test.go:327: obligation-free landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/partly-delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/beyond-scope,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (753 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/directory-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/listed-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/outside-directory,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (762 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/scope-without-deliverable,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (799 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/sibling-prefix,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (757 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/unlisted,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (753 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-then-tickets,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-with-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (711 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-without-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/tickets-then-spec,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (721 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPartialDelivery,\"commitment_landing_test.go:56: partial landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (654 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPersistenceBeforeGate,\"commitment_landing_test.go:173: torn-intent landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPublishLock,\"land_identity_test.go:307: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then b… (660 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentStaleLanding,\"land_identity_test.go:209: displaced landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (699 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCompareAndSwapRefusalNamesTheFold,\"land_moved_destination_test.go:66: moved destination = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold,\"land_moved_destination_test.go:43: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (732 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:132: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict-pending,\"refusal_route_test.go:132: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestDelegatedIntegrationJourney,\"delegated_integration_test.go:237: delegated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (650 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/installed,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (741 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/kit-checkout,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (747 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCleansTheFoldedSibling,\"land_effects_cleanup_test.go:83: folded-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAcceptsAbbreviatedIdentities,\"land_surface_test.go:46: abbreviated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthenticatesDigestShapedRequestToken,\"land_reauthorization_test.go:102: digest-shaped request land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (756 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesCaptureOutsideTheFence,\"land_surface_test.go:135: capture landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesFromTicketWrites,\"land_derived_fence_test.go:34: ticket-only landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/another_review_record,\"land_derived_fence_test.go:114: other review landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (743 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/pickup,\"land_derived_fence_test.go:100: pickup landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandCarriesTheBaselineScheduleRootIntoTheProspectiveGate,\"land_freshness_test.go:288: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (701 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandComposesCaptureOntoMovedDestination,\"land_surface_test.go:71: unreviewed capture composition = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (679 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictNextPointsThroughUnsafePath,\"land_surface_test.go:360: unsafe-path conflict next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesThePath,\"land_surface_test.go:268: conflict refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (689 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesTheSourceRepair,\"land_surface_test.go:301: conflict repair next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandDisclosesAUnionResolution,\"land_surface_test.go:104: unreviewed union = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedBase,\"land_reauthorization_test.go:175: abbreviated base = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"7b99bcc967d6\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (725 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/39,\"land_reauthorization_test.go:163: abbreviated source tip \\\"15af999e6b1d44efafbc326f20845761966cf6d\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or a… (752 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/4,\"land_reauthorization_test.go:163: abbreviated source tip \\\"67e1\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line … (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/upper-12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"82BD8B100CAF\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (729 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandFenceRefusalNamesThePath,\"land_surface_test.go:220: fence refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; t… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal,\"land_freshness_test.go:114: forged-primary landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (699 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIncompleteNextUsesAssignmentPointerForUnsafePath,\"land_resume_test.go:38: unsafe-path incomplete = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease,\"land_freshness_test.go:72: single-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandLeavesTheDestinationUnchangedAfterARedProspectiveGate,\"land_freshness_test.go:185: gate tally = \\\"\\\", open /dev/shm/HN6ZFU/t/TestLandCommandLeavesTheDestinationUnchangedAfterARedProspective3748901803/003/.git/bench-land-gate-tally: no such file or directory, want one prospective run\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout,\"land_broker_notice_test.go:125: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (719 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion,\"land_freshness_test.go:39: stable-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/marker,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (643 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/reconcile,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (650 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/release,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (646 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsNoHeadsLineWithoutRecords,\"land_census_test.go:91: empty census landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsTheCensusHeadBreakdown,\"land_census_test.go:76: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/absent,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (648 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/concurrently-moved,\"land_flags_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/present,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (650 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicConflictRepairRequiresNewReviewedTip,\"land_journey_test.go:294: conflict result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicLandsAnInRangeSpecAmendment,\"land_spec_amendment_test.go:47: amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicPreservesHistoricalRuntimeLogs,\"land_journey_test.go:139: runtime-log landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/clean,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (647 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/declared-output,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-and-unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (657 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log-empty-declaration,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicResumeCompletesPublishedReleaseWithoutRepublishing,\"land_resume_test.go:55: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublishedReleaseFailureExitsIncomplete,\"land_journey_test.go:267: published release exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReadsOnlyTrackedTickets,\"land_derived_fence_test.go:61: ignored-ticket landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (723 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalKeepsControlBearingPathInOneTableRow,\"land_release_refusal_test.go:70: control refusal = (1, \\\"refused{detail=landing destination is not clean,next=reviewer: commit the destination's uncommitted work, or discard it; then bench worktree land --request 'refusal-controls' --base '7… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalListsCollidingPaths,\"land_release_refusal_test.go:52: collision refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusesPostGateUnknownIgnoredMutation,\"land_journey_test.go:175: post-gate ignored mutation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (679 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReleaseDiagnosticCannotForgeTerminalLines,\"land_flags_test.go:132: hostile release result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRemovesEveryTemporaryProspectiveArtifact/published,\"land_freshness_test.go:214: landing exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:54: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/marker,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (682 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/reconcile,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/release,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeOverlappingDestinationEdit,\"land_journey_test.go:237: last-moment overlapping edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (742 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeTrackedDestinationEdit,\"land_journey_test.go:199: last-moment tracked edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesTheCensusCountAndDropsTheRecords,\"land_census_test.go:45: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (687 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesZeroForAnAssignmentWithNoRecords,\"land_census_test.go:61: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandEffectsRowPrecedesTheLandedRecord,\"land_effects_test.go:238: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then be… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandFenceRefusalOverAFoldedSourceKeepsTheCallerBase,\"land_folded_base_test.go:34: fence refusal = \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandGradesASourceCommittedByALanePass,\"land_journey_test.go:381: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandHoldsTheDestinationLease,\"land_lease_test.go:90: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writ… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandLeavesAPriorLandedAssignment,\"land_effects_cleanup_test.go:114: narrowed landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshReadsTheSeal,\"land_effects_test.go:220: unsealed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshesTheBrokerAfterPublication,\"land_effects_test.go:160: landing with a stale broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRelaysTheBoundedGreenShapeBeforeTheLandedRecord,\"land_journey_test.go:410: land exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReleasesTheLeaseOnRefusal,\"land_lease_test.go:171: red landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedCleanup,\"land_effects_cleanup_test.go:188: faulted cleanup = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedRefresh,\"land_effects_test.go:181: failed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandResumeTakesNoLease,\"land_lease_test.go:190: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (634 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsAnUnprovenSibling,\"land_effects_cleanup_test.go:140: retaining landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/ordinary,\"land_empty_sibling_test.go:45: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/resumed,\"land_empty_sibling_test.go:35: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (665 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsAFreshBroker,\"land_effects_test.go:140: landing with a fresh broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsTheRefreshWithoutBuildInputs,\"land_effects_test.go:118: landing without build inputs = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict,\"refusal_route_follow_test.go:59: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-compositi… (618 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict-pending,\"refusal_route_follow_test.go:59: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-c… (642 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-collision,\"refusal_route_follow_test.go:54: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-not-clean,\"refusal_route_follow_test.go:72: destination-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (730 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-build-window,\"refusal_route_follow_test.go:72: land-build-window re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-destination-moved,\"refusal_route_follow_test.go:54: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (739 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-handback,\"refusal_route_follow_test.go:72: land-handback re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-incomplete,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (676 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-infrastructure,\"refusal_route_follow_test.go:69: printed step \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen\\\" is not one simple command with every slot filled\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-lease-held,\"refusal_route_follow_test.go:72: land-lease-held re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-red,\"refusal_route_follow_test.go:59: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base '803… (554 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-destination-residue,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-marker,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-clean,\"refusal_route_follow_test.go:72: source-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-fenced,\"refusal_route_follow_test.go:72: source-not-fenced re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-tip-mismatch,\"refusal_route_follow_test.go:72: source-tip-mismatch re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingPrintsOutputBreakdown,\"land_census_output_test.go:25: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (682 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRedRouteNamesTheRepair,\"refusal_route_test.go:147: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base '9394ad9c9… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture,\"panic: test timed out after 10m0s\",135\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/destination-collision,\"identity_component_test.go:238: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (743 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-destination-moved,\"identity_component_test.go:238: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (745 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-incomplete,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-destination-residue,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (700 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-marker,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingReleaseFactAdapterTranslatesReleaseExit,\"land_facts_test.go:211: release-failure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRetiresAnAbsentLandedSibling,\"missing_tree_recovery_test.go:231: absent-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/dangling-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/directory,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/live-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/current-source,\"land_moved_destination_test.go:132: current-source = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/non-proof-cause,\"land_moved_destination_test.go:132: non-proof-cause = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/unreadable-ancestry,\"land_moved_destination_test.go:126: rename /dev/shm/HN6ZFU/t/TestProofFaultOverACurrentSourceHandsBackunreadable-ancestry1230209754/003/.git/objects/99/dfdeb13c8911060006192c63b97f3460ee5f33.hidden /dev/shm/HN6ZFU/t/TestProofFaultOverACurre… (372 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeHandsBackADetachedDestination,\"land_resume_refusal_test.go:366: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (666 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/./specs/x/spec.md,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/x,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandCompletesAnInterruptedMarker,\"land_resume_test.go:230: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandExpandsAbbreviatedIdentities,\"land_reauthorization_test.go:184: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandFollowupFailureExitsIncomplete,\"land_resume_test.go:19: first incomplete exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-path,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (720 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-state,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (722 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/lock,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/owner-marker,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/registration,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/request,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-path,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-request,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicBindsPublishedLandingIdentity,\"land_resume_test.go:101: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicCompletesAnAmendedSourceLanding,\"land_spec_amendment_test.go:75: interrupted amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandReconcilesAnUnreconciledPublishedCheckout,\"land_resume_test.go:183: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/absent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/behind,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/divergent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesNonAncestorReviewBaseWithoutMutation,\"land_resume_refusal_test.go:143: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (675 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesWhenTerminalReceiptWasEvicted,\"land_resume_refusal_test.go:204: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRepeatsTheCensusCount,\"land_resume_test.go:285: incomplete landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandSourceRefusalNamesTheCallersResume,\"land_resume_refusal_test.go:119: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandUnknownRequestNamesReauthorizeRecovery,\"land_reauthorization_test.go:16: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding,\"land_specless_test.go:151: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCompletesTheUnfinishedEffects,\"land_effects_cleanup_test.go:218: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (674 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandRunsTheEffectsAfterRelease,\"land_effects_cleanup_test.go:270: failed-refresh landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeReadsEffectStateFromTheTree,\"land_effects_test.go:269: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict,\"refusal_route_test.go:187: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict-pending,\"refusal_route_test.go:187: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/destination-collision,\"refusal_route_test.go:185: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (733 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-destination-residue,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-marker,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestSecondLandingMeetsTheLease,\"land_lease_test.go:133: first landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket… (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnclaimedLandingCompositionErrorHandsBack,\"refusal_route_test.go:164: unclaimed composition error = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (736 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnsafePathRouteUsesThePlaceholder/incomplete_resume,\"land_release_refusal_test.go:204: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (558 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowAdmitsTheBuildsOwnLanding,\"land_window_test.go:119: landing of the bound spec from the sibling = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: l… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowClosesAtCompletion,\"land_window_test.go:130: landing of the build's completion = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/build-start,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/gap-close,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (685 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:ft293_t4_r2",
              "digest": "sha256:f090cf34dd5f68b785a1feef9675178b49165f935040f97cda30db5cf8681358",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/fence_writes.go,swap,failed,198,TestCommitmentClosureAdmission/rowless-FT2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1434\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,71136\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,600012,1434\nfailures[198]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/rowless-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (711 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT2,\"commitment_landing_test.go:101: closing FT2 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureAdmission/spec-residual-FT3,\"commitment_landing_test.go:101: closing FT3 = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/rowless,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureGateRed/spec,\"commitment_landing_test.go:125: red closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/rowless,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (702 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentClosureResume/spec,\"commitment_landing_test.go:222: interrupted delivery = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentDeliveryClosure,\"commitment_landing_test.go:27: delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/blocker,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentGateRace/default-branch-policy,\"land_identity_test.go:247: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentInterruptedGate,\"commitment_landing_test.go:143: interrupted closure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLandingRefusesObligationFreeDelivery,\"commitment_light_landing_test.go:327: obligation-free landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyClosure/partly-delivered-scope,\"commitment_light_landing_test.go:59: legacy delivery landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/beyond-scope,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (753 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/directory-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/listed-scope,\"land_spec_amendment_test.go:150: listed continuation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/outside-directory,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (762 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/scope-without-deliverable,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (799 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/sibling-prefix,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (757 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentLegacyContinuation/unlisted,\"land_spec_amendment_test.go:155: continuation landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (753 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-then-tickets,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-with-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (711 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/spec-without-board,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentNoRoadmapOwner/tickets-then-spec,\"commitment_light_landing_test.go:280: rowless delivery of specs/x/spec.md = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wri… (721 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPartialDelivery,\"commitment_landing_test.go:56: partial landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (654 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPersistenceBeforeGate,\"commitment_landing_test.go:173: torn-intent landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentPublishLock,\"land_identity_test.go:307: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then b… (660 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCommitmentStaleLanding,\"land_identity_test.go:209: displaced landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (699 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCompareAndSwapRefusalNamesTheFold,\"land_moved_destination_test.go:66: moved destination = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold,\"land_moved_destination_test.go:43: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (732 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:132: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict-pending,\"refusal_route_test.go:132: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestDelegatedIntegrationJourney,\"delegated_integration_test.go:237: delegated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (650 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/installed,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (741 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandBrokerNoticeAfterACompleteRefresh/kit-checkout,\"land_broker_notice_test.go:98: refreshed broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line… (747 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCleansTheFoldedSibling,\"land_effects_cleanup_test.go:83: folded-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAcceptsAbbreviatedIdentities,\"land_surface_test.go:46: abbreviated landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthenticatesDigestShapedRequestToken,\"land_reauthorization_test.go:102: digest-shaped request land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (756 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesCaptureOutsideTheFence,\"land_surface_test.go:135: capture landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesFromTicketWrites,\"land_derived_fence_test.go:34: ticket-only landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/another_review_record,\"land_derived_fence_test.go:114: other review landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (743 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandAuthorizesTheReviewPickup/pickup,\"land_derived_fence_test.go:100: pickup landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandCarriesTheBaselineScheduleRootIntoTheProspectiveGate,\"land_freshness_test.go:288: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (701 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandComposesCaptureOntoMovedDestination,\"land_surface_test.go:71: unreviewed capture composition = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (679 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictNextPointsThroughUnsafePath,\"land_surface_test.go:360: unsafe-path conflict next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesThePath,\"land_surface_test.go:268: conflict refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (689 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandConflictRefusalNamesTheSourceRepair,\"land_surface_test.go:301: conflict repair next = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandDisclosesAUnionResolution,\"land_surface_test.go:104: unreviewed union = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedBase,\"land_reauthorization_test.go:175: abbreviated base = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (692 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"7b99bcc967d6\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (725 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/39,\"land_reauthorization_test.go:163: abbreviated source tip \\\"15af999e6b1d44efafbc326f20845761966cf6d\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or a… (752 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/4,\"land_reauthorization_test.go:163: abbreviated source tip \\\"67e1\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line … (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandExpandsAbbreviatedSourceTip/upper-12,\"land_reauthorization_test.go:163: abbreviated source tip \\\"82BD8B100CAF\\\" = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Write… (729 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandFenceRefusalNamesThePath,\"land_surface_test.go:220: fence refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; t… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal,\"land_freshness_test.go:114: forged-primary landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (699 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandIncompleteNextUsesAssignmentPointerForUnsafePath,\"land_resume_test.go:38: unsafe-path incomplete = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (691 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease,\"land_freshness_test.go:72: single-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandLeavesTheDestinationUnchangedAfterARedProspectiveGate,\"land_freshness_test.go:185: gate tally = \\\"\\\", open /dev/shm/HN6ZFU/t/TestLandCommandLeavesTheDestinationUnchangedAfterARedProspective3748901803/003/.git/bench-land-gate-tally: no such file or directory, want one prospective run\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout,\"land_broker_notice_test.go:125: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (719 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion,\"land_freshness_test.go:39: stable-owner landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/marker,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (643 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/reconcile,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (650 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPostCASTerminalTable/release,\"land_flags_test.go:115: land = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench w… (646 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsNoHeadsLineWithoutRecords,\"land_census_test.go:91: empty census landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPrintsTheCensusHeadBreakdown,\"land_census_test.go:76: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/absent,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (648 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/concurrently-moved,\"land_flags_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandProjectGreenOrderTable/present,\"land_flags_test.go:219: marker result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (650 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicConflictRepairRequiresNewReviewedTip,\"land_journey_test.go:294: conflict result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen;… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicLandsAnInRangeSpecAmendment,\"land_spec_amendment_test.go:47: amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicPreservesHistoricalRuntimeLogs,\"land_journey_test.go:139: runtime-log landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/clean,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (647 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/declared-output,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-and-unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (657 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/runtime-log-empty-declaration,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicRealGitJourney/unknown-ignored,\"land_journey_test.go:64: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublicResumeCompletesPublishedReleaseWithoutRepublishing,\"land_resume_test.go:55: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandPublishedReleaseFailureExitsIncomplete,\"land_journey_test.go:267: published release exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (695 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReadsOnlyTrackedTickets,\"land_derived_fence_test.go:61: ignored-ticket landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (723 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalKeepsControlBearingPathInOneTableRow,\"land_release_refusal_test.go:70: control refusal = (1, \\\"refused{detail=landing destination is not clean,next=reviewer: commit the destination's uncommitted work, or discard it; then bench worktree land --request 'refusal-controls' --base '7… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusalListsCollidingPaths,\"land_release_refusal_test.go:52: collision refusal = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRefusesPostGateUnknownIgnoredMutation,\"land_journey_test.go:175: post-gate ignored mutation = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (679 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReleaseDiagnosticCannotForgeTerminalLines,\"land_flags_test.go:132: hostile release result = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRemovesEveryTemporaryProspectiveArtifact/published,\"land_freshness_test.go:214: landing exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandReportsInstallStepForABrokerChangingDiff,\"land_broker_notice_test.go:54: broker-changing landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/marker,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (682 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/reconcile,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing/release,\"land_freshness_test.go:247: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeOverlappingDestinationEdit,\"land_journey_test.go:237: last-moment overlapping edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line w… (742 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandRetainsJustInTimeTrackedDestinationEdit,\"land_journey_test.go:199: last-moment tracked edit landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesTheCensusCountAndDropsTheRecords,\"land_census_test.go:45: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (687 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandCommandStatesZeroForAnAssignmentWithNoRecords,\"land_census_test.go:61: landed record = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; the… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandEffectsRowPrecedesTheLandedRecord,\"land_effects_test.go:238: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then be… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandFenceRefusalOverAFoldedSourceKeepsTheCallerBase,\"land_folded_base_test.go:34: fence refusal = \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandGradesASourceCommittedByALanePass,\"land_journey_test.go:381: land exit=1 stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandHoldsTheDestinationLease,\"land_lease_test.go:90: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writ… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandLeavesAPriorLandedAssignment,\"land_effects_cleanup_test.go:114: narrowed landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshReadsTheSeal,\"land_effects_test.go:220: unsealed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen… (673 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRefreshesTheBrokerAfterPublication,\"land_effects_test.go:160: landing with a stale broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (707 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRelaysTheBoundedGreenShapeBeforeTheLandedRecord,\"land_journey_test.go:410: land exit = 1, want 0; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReleasesTheLeaseOnRefusal,\"land_lease_test.go:171: red landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then … (656 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedCleanup,\"land_effects_cleanup_test.go:188: faulted cleanup = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spe… (685 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandReportsAFailedRefresh,\"land_effects_test.go:181: failed refresh = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandResumeTakesNoLease,\"land_lease_test.go:190: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (634 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsAnUnprovenSibling,\"land_effects_cleanup_test.go:140: retaining landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (694 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/ordinary,\"land_empty_sibling_test.go:45: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; th… (655 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandRetainsSiblingBornAtPublishedTip/resumed,\"land_empty_sibling_test.go:35: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (665 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsAFreshBroker,\"land_effects_test.go:140: landing with a fresh broker = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandSkipsTheRefreshWithoutBuildInputs,\"land_effects_test.go:118: landing without build inputs = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict,\"refusal_route_follow_test.go:59: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-compositi… (618 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/composition-conflict-pending,\"refusal_route_follow_test.go:59: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-c… (642 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-collision,\"refusal_route_follow_test.go:54: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (735 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/destination-not-clean,\"refusal_route_follow_test.go:72: destination-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wi… (730 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-build-window,\"refusal_route_follow_test.go:72: land-build-window re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-destination-moved,\"refusal_route_follow_test.go:54: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line wit… (739 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-handback,\"refusal_route_follow_test.go:72: land-handback re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-incomplete,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (676 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-infrastructure,\"refusal_route_follow_test.go:69: printed step \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen\\\" is not one simple command with every slot filled\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-lease-held,\"refusal_route_follow_test.go:72: land-lease-held re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with ben… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/land-red,\"refusal_route_follow_test.go:59: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base '803… (554 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-destination-residue,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/resume-marker,\"refusal_route_follow_test.go:54: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-clean,\"refusal_route_follow_test.go:72: source-not-clean re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-not-fenced,\"refusal_route_follow_test.go:72: source-not-fenced re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (718 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingFacesFollowTheirRoutes/source-tip-mismatch,\"refusal_route_follow_test.go:72: source-tip-mismatch re-run = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (724 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingPrintsOutputBreakdown,\"land_census_output_test.go:25: landing evidence = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (682 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRedRouteNamesTheRepair,\"refusal_route_test.go:147: land-red next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-land-red' --base '9394ad9c9… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture,\"panic: test timed out after 10m0s\",135\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/destination-collision,\"identity_component_test.go:238: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (743 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-destination-moved,\"identity_component_test.go:238: face land-destination-moved = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (745 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-incomplete,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-destination-residue,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (700 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/resume-marker,\"identity_component_test.go:238: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingReleaseFactAdapterTranslatesReleaseExit,\"land_facts_test.go:211: release-failure landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec … (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRetiresAnAbsentLandedSibling,\"missing_tree_recovery_test.go:231: absent-sibling landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/dangling-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/directory,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLeaseLockUnavailable/live-link,\"land_lease_test.go:228: landing over a blocked lease lock = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with b… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/current-source,\"land_moved_destination_test.go:132: current-source = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench sp… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/non-proof-cause,\"land_moved_destination_test.go:132: non-proof-cause = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/unreadable-ancestry,\"land_moved_destination_test.go:126: rename /dev/shm/HN6ZFU/t/TestProofFaultOverACurrentSourceHandsBackunreadable-ancestry1230209754/003/.git/objects/99/dfdeb13c8911060006192c63b97f3460ee5f33.hidden /dev/shm/HN6ZFU/t/TestProofFaultOverACurre… (372 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeHandsBackADetachedDestination,\"land_resume_refusal_test.go:366: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (666 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/./specs/x/spec.md,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (710 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandAcceptsSpecSlugAndPath/x,\"land_resume_test.go:214: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandCompletesAnInterruptedMarker,\"land_resume_test.go:230: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (681 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandExpandsAbbreviatedIdentities,\"land_reauthorization_test.go:184: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandFollowupFailureExitsIncomplete,\"land_resume_test.go:19: first incomplete exit = 1, want 3; stdout=\\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-path,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (720 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/assignment-state,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (722 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/lock,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (698 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/owner-marker,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/registration,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (712 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandNamesEachIdentityComponent/request,\"identity_component_test.go:175: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (704 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-path,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPreauthenticatesCompletedRequestAndPath/wrong-request,\"land_resume_test.go:261: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then ben… (671 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicBindsPublishedLandingIdentity,\"land_resume_test.go:101: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicCompletesAnAmendedSourceLanding,\"land_spec_amendment_test.go:75: interrupted amended landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with… (693 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (697 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (716 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (714 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL25/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/clean,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/declared_ignored_output,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/nested_repository,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (709 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/staged_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/tracked-worktree_changes,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/undeclared_ignored_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (715 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandPublicRefusesDestructiveDestinationState/PL30/untracked_file,\"land_resume_refusal_test.go:71: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench s… (706 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandReconcilesAnUnreconciledPublishedCheckout,\"land_resume_test.go:183: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wid… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/absent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (683 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/behind,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves/divergent,\"land_resume_refusal_test.go:180: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (686 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesNonAncestorReviewBaseWithoutMutation,\"land_resume_refusal_test.go:143: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (675 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRefusesWhenTerminalReceiptWasEvicted,\"land_resume_refusal_test.go:204: landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; … (667 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandRepeatsTheCensusCount,\"land_resume_test.go:285: incomplete landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wide… (677 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandSourceRefusalNamesTheCallersResume,\"land_resume_refusal_test.go:119: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (678 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandUnknownRequestNamesReauthorizeRecovery,\"land_reauthorization_test.go:16: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench … (713 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding,\"land_specless_test.go:151: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (708 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandCompletesTheUnfinishedEffects,\"land_effects_cleanup_test.go:218: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (674 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeLandRunsTheEffectsAfterRelease,\"land_effects_cleanup_test.go:270: failed-refresh landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with be… (705 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestResumeReadsEffectStateFromTheTree,\"land_effects_test.go:269: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec wi… (659 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict,\"refusal_route_test.go:187: composition-conflict next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composition-con… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/composition-conflict-pending,\"refusal_route_test.go:187: composition-conflict-pending next = \\\"take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen; then bench worktree land --request 'landing-face-composi… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/destination-collision,\"refusal_route_test.go:185: face destination-collision = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (733 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-destination-residue,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (696 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestReviewerLandFacesOpenWithTheMarker/resume-marker,\"refusal_route_test.go:185: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec w… (670 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestSecondLandingMeetsTheLease,\"land_lease_test.go:133: first landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket… (672 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnclaimedLandingCompositionErrorHandsBack,\"refusal_route_test.go:164: unclaimed composition error = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with benc… (736 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnsafePathRouteUsesThePlaceholder/incomplete_resume,\"land_release_refusal_test.go:204: interrupted landing = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench… (558 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowAdmitsTheBuildsOwnLanding,\"land_window_test.go:119: landing of the bound spec from the sibling = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: l… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowClosesAtCompletion,\"land_window_test.go:130: landing of the build's completion = (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with … (684 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/build-start,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (690 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestWindowRecheckAtPublish/gap-close,\"land_window_test.go:208: landing finished without reaching its gate barrier: (1, \\\"refused{detail=reviewed source range or ownership fence is invalid,next=take the refusal_paths entries out of the reviewed range, or add them to a ticket's Wr… (685 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t4r2-t4-system",
          "performer": "claude:ft293_t4_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r2",
            "digest": "sha256:e9ef38ae77d000dc49028742c34485fe5a02c9a67666c32f03e4617159b65619",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,170,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,69400,170\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t4r2-t4-root-conformance",
          "performer": "claude:ft293_t4_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r2",
            "digest": "sha256:655243b255cc9e54ce66e3d659a564f6a9577fea918194abfdd7edfcdf2426f2",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14860,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t4r2-t4-guidance-anchors",
          "performer": "claude:ft293_t4_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r2",
            "digest": "sha256:df300cc7cb9ba51d5179e645f0206f458f56819fa05ff3e1d9e3d5b11d9b4dd9",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4564,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t4r2-t4-guidance-budgets",
          "performer": "claude:ft293_t4_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t4_r2",
            "digest": "sha256:b682aaa6a2f3a10908927cd53ce11e0bf9cf598d196444453baf84591b6728d5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "cr3-standards",
          "performer": "claude:ft293_cr3_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr3_standards",
            "digest": "sha256:ee83ac9a58ef44a2414d0d6f8f0f028b7983f2eb95939fce4378d87520e5eec3",
            "excerpt": "Standards axis, FT293 CR3 (a883604f..82309770). check-current current=true, dirty_record=false.\nS1 auto-fix 7: staleness/collect.go:84-96 copies coversAt (staleness/coverage.go:73-82): ReadControlBlob, the StageTicket wrap, ParseTicket args, and the ticket-file test (collect.go:89, coverage.go:38).\nS2 ask-user 4: ticketWrites lists tickets from git objects without the duplicate-basename rule of tickets.Enumerate, which ticketFence reads (enumerate.go:57-59).\nS3 auto-fix 5: craft-spec SKILL.md:65 says \"plus the implicit entries\" without naming them.\nS4 auto-fix 4: craft-spec SKILL.md:70 \"Each fence\" is ambiguous after one derived fence.\nS5 no-op 3: fence_writes.go keeps the retired check's file name.\nImplicit-entry set one source; CR65 generated from MetadataBlocks; no dead code; budgets hold. Count 5. Worst: S1.\n"
          },
          "axis": "Standards",
          "base": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "tip": "8230977099a3702443b965c5af3f5f24b6859173",
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
          "id": "cr3-spec",
          "performer": "claude:ft293_cr3_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr3_spec",
            "digest": "sha256:7569a3b9cc91d6a16770021c89b3121b85ec122dc414db3d18c0a94e83740f32",
            "excerpt": "Spec axis, FT293 CR3 (a883604f..82309770). check-current current=true, dirty_record=false. bench coverage --check: 94 rows valid.\nP1 no-op 3: story 32 is looser than decision spec.md:135; the fence block lists the ticket union as decided.\nP2 no-op for CR3 8: craft-review SKILL.md:42 \"The ownership fence lives in craft-spec\" is stale; ticket 06 owns it.\nRows CR55-CR72, CR77, CR86, CR87, CR90 closed; CR65 byte-equal; CR71 FenceTokens absent.\nCR75: union 220 paths, section 221; the only extra is reviews/closure-removal.md; no gap; every changed path since 36b1fafc is covered.\nCR76: no canary file changed; anchor rows keep the literal shape. Story 38 holds.\nCount 2 non-blocking. Worst: P2.\n"
          },
          "axis": "Spec",
          "base": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "tip": "8230977099a3702443b965c5af3f5f24b6859173",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "cr3-coverage",
          "performer": "claude:ft293_cr3_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9bc8d7cd35aeeea1516f1ed9082f8776e601067",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr3_coverage",
            "digest": "sha256:8c849a39c54bdf32b08107c7f1730b7b846d6235bfb670150004837260bba244",
            "excerpt": "Coverage axis, FT293 CR3. check-current current=true at 96e49cc4; 5 probes, all restored=yes; tree clean.\nP1 required: pickup widened to reviews: bit (TestCommandBuildIgnoresAFenceSection, TestCommandReviewAuthorizesFromTicketWrites, TestLandCommandAuthorizesTheReviewPickup/another_review_record).\nP2 spec-folder entry becomes specs: bit (LS9). P3 plain string prefix in fenceAuthorizes: bit (TestCommandFencePrefixBoundary, TestDecidePathsAuthorized, LS10).\nP4 ticketFence drops (deleted) entries: silent. P5 audit ticketWrites drops nested tickets: silent.\nC1 fix 8: no test pins that a (deleted) entry authorizes its own deletion.\nC2 fix 7: the audit scope drops tickets under tickets/sub/ while the gather authorizes them.\nC3 decision 6: AuthorizeReviewedSource gathers tickets from the filesystem, scan does not skip dot files, and the landing dirty check leaves out ignored files, so an ignored ticket-shaped file can authorize paths; the audit reads committed objects.\nCount 3. Worst: C1.\n"
          },
          "axis": "Coverage",
          "base": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "tip": "8230977099a3702443b965c5af3f5f24b6859173",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "cr3c-standards",
          "performer": "claude:ft293_cr3c_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr3c_standards",
            "digest": "sha256:dcb8345b628badc17af1f9f701eee3d4b479a1f62ace541cdc5bd025449cda80",
            "excerpt": "Standards confirming round, FT293 CR3, repair delta 82309770..9d26cb3d. check-current current=true.\nS1 closed 9: staleness.TicketsAt (coverage.go:29-51) is the one reader of a commit's ticket set for the Writes: face, the Covers: face, and AuthorizeReviewedSource; no coversAt copy remains.\nS2 closed 8: uniqueNames (enumerate.go:118) is the one duplicate-basename rule for Enumerate and EnumerateTracked.\nS3 closed 6: craft-spec SKILL.md:65 names the three implicit entries; the only guidance list of them.\nS4 closed 8: SKILL.md:71 uses \"ownership fence\" and \"part\".\nticketsDir: not made worse 8. Outside the delta: internal/spec/widen.go:192-208 reads tracked tickets with its own rule.\nOpen 0.\n"
          },
          "axis": "Standards",
          "base": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "tip": "9d26cb3d4e5ee01d1e401d26d727fda2821e2230",
          "finding_ids": [],
          "supersedes": [
            "cr3-standards"
          ]
        },
        {
          "id": "cr3c-spec",
          "performer": "claude:ft293_cr3c_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr3c_spec",
            "digest": "sha256:68d7871ea50fbb6a773fade9fe7c8f76846c20baf27059659c094e7a2474d6d3",
            "excerpt": "Spec confirming round, FT293 CR3, repair delta 82309770..9d26cb3d. check-current current=true at 4e2fe40a. bench coverage --check: 94 rows valid.\nLanding and build readers match spec.md:131 held 9; CR3 rows closed, CR66/CR67 cover nested and duplicate tickets held 8; no unrequested behavior held 8; CR75 section equals the union plus the pickup held 9.\nP1 flag 6: spec.md:131 says bench preflight review reads that one set, but review preflight reads the worktree tickets (gather.go:310), so an ignored ticket-shaped file can turn its paths-authorized green while the landing refuses.\nCount 1. Worst: P1.\n"
          },
          "axis": "Spec",
          "base": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "tip": "9d26cb3d4e5ee01d1e401d26d727fda2821e2230",
          "finding_ids": [
            "P1"
          ],
          "supersedes": [
            "cr3-spec"
          ]
        },
        {
          "id": "cr3c-coverage",
          "performer": "claude:ft293_cr3c_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2a40470ea22be7ddef48a2fb3c0314e35422ece2",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr3c_coverage",
            "digest": "sha256:12927ccf72fed22755668d651dd37b6f922a8f316021ec24b42ff9120293628d",
            "excerpt": "Coverage confirming round, FT293 CR3. check-current current at 4e2fe40a; 7 probes, all restored=yes; tree clean.\nC3 open 8: gather.go:211-218 omit facts.Tickets = nil bit TestLandCommandReadsOnlyTrackedTickets; keep checkout tickets on a TicketsAt error silent; read TicketsAt at HEAD silent (equivalent in practice). A committed symlink in the tickets folder reaches the error path.\nC2/S2 closed for the stated cases 7: slices.Reverse(walk) bit; depth-first sort silent; no test pins parity when a nested duplicate sorts first.\nC1 closed 9: (deleted) authorizes only a still-held path bit TestCommandAuthorizesADeletedEntry.\nC4 new 6: bench preflight review reads checkout tickets, and the clean check passes ignored files, so an ignored ticket-shaped file can turn rows-membership, paths-authorized, and the charge fence block green; the landing re-checks only paths.\nOpen: C3, C4. Worst: C3.\n"
          },
          "axis": "Coverage",
          "base": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "tip": "9d26cb3d4e5ee01d1e401d26d727fda2821e2230",
          "finding_ids": [
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": [
            "cr3-coverage"
          ]
        },
        {
          "id": "cr3c2-standards",
          "performer": "claude:ft293_cr3c2_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr3c2_standards",
            "digest": "sha256:f6d24bb40af3d855df9843922544d5a34b3e55c3130ac73fbe11962216970361",
            "excerpt": "Standards second confirming round, FT293 CR3, repair delta 9d26cb3d..ffee6a38. check-current current=true, dirty_record=false.\nTracked-ticket reader one source (staleness.TicketEntriesAt; TicketsAt parses it); review mode reads it (gather.go:91-95); the AuthorizeReviewedSource override is removed. Duplicate rule one source (uniqueNames). Representable check one source (unrepresentable, enumerate.go:66). Refusal conversion once (gather.go:320-326). Comments current. Test-count pin 798.\nAdvice (judgment, non-blocking): gatherTicketsWithPolicy is a one-line pass-through (gather.go:281-283) that keeps the walk's doc.\nOpen 0.\n"
          },
          "axis": "Standards",
          "base": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "tip": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "finding_ids": [],
          "supersedes": [
            "cr3c-standards"
          ]
        },
        {
          "id": "cr3c2-spec",
          "performer": "claude:ft293_cr3c2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr3c2_spec",
            "digest": "sha256:6fdb7e7fcdb6471339ceb19ed9668aba6fd80240b3e85d294e9aac3ee2f5d94c",
            "excerpt": "Spec second confirming round, FT293 CR3, repair delta 9d26cb3d..ffee6a38. check-current current=true. bench coverage --check: 94 rows valid.\nspec.md:131 matches review mode, the landing, and build mode held 9; P1 closed. CR3 rows held 8; CR64 and CR77 read the same tracked tickets. No unrequested behavior held 8. CR75: 224 section entries equal the 223 ticket paths plus the pickup held 9.\nProcess note: the axis wrote and deleted two temporary files under /dev/shm against its read-only charge; the tree is unchanged.\nNew findings: none. Open 0.\n"
          },
          "axis": "Spec",
          "base": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "tip": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "finding_ids": [],
          "supersedes": [
            "cr3c-spec"
          ]
        },
        {
          "id": "cr3c2-coverage",
          "performer": "claude:ft293_cr3c2_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "682ada7998b7ee53beb77d62d59493820d93994f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr3c2_coverage",
            "digest": "sha256:12731283fb7d009a953a58632fa68d367c4a040ce3b9cb705574fb07e32e8a43",
            "excerpt": "Coverage second confirming round, FT293 CR3. check-current current at 99e3bfc9; every probe restored=yes; tree clean.\nC3 closed 7: failed tracked read becomes an empty set bit TestLandCommandRefusesAnUnreadableTicketSet; skip a symlinked member bit.\nC4/P1 closed 8: review reads checkout tickets for every row while only the landing re-reads bit TestCommandReviewReadsOnlyTrackedTickets and TestCommandReviewRefusesAnUnrepresentableTrackedTicket.\nAdvice (hardening-cycle checks, code correct): byte-order sort in EnumerateTracked silent (a/b/z.md versus a-b/z.md parity has no test); representable check limited to .md entries silent (no non-.md case).\nOpen 0.\n"
          },
          "axis": "Coverage",
          "base": "a883604feb1171d1ae1c43159f1e0b4377b13a22",
          "tip": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "finding_ids": [],
          "supersedes": [
            "cr3c-coverage"
          ]
        }
      ]
    },
    {
      "id": "CR4",
      "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
      "tip": "1b79d041f615d1e7e8f322c98245c0de6ad7def7",
      "plan_digest": "sha256:40e3bcd3128efe554b1c3045d1e51f4107e10d7f6f8c9590d56d82e1cb1420f9",
      "source_digest": "268b808950fbf7fe79c783f91d7a64f110e719c2",
      "acceptance_rows": [
        "CR75",
        "CR76",
        "CR78",
        "CR79",
        "CR88",
        "CR94"
      ],
      "verification": [
        {
          "id": "t6-t6-anchors-and-help",
          "performer": "claude:ft293_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d790f343b3965cfd4d7b497e659303ee4f78de12",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6",
            "digest": "sha256:d4b50b19971536b961aa43b79b9eb31b02a4d77e342cd02ac4b3221c65fbc204",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16077,346\n  github.com/gibbonmi/bench/internal/anchors,pass,120,123\n  github.com/gibbonmi/bench/internal/conformance,pass,35534,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t6-anchors-and-help",
          "command": "bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t6-t6-root-conformance",
          "performer": "claude:ft293_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d790f343b3965cfd4d7b497e659303ee4f78de12",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6",
            "digest": "sha256:760e7c97fbb57c430cff39f7018f7696a5dfd9f69cad47a89a1cd6fdb9a62b92",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14580,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t6-t6-guidance-anchors",
          "performer": "claude:ft293_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d790f343b3965cfd4d7b497e659303ee4f78de12",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6",
            "digest": "sha256:26da588d9dc794fe7952cf6a0ce637528c3fd7202b4f7563f49418ff545420a7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.agents/skills/bench-craft-delegate/references/delegation-discipline.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,438\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,4509,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: closure removal: the delegate guidance called the fence a refusal boundary\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "Insertion: restore the refusal-boundary sentence in delegation-discipline.md; the forbid anchor must red docs-currency-workflow, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t6",
              "digest": "sha256:26da588d9dc794fe7952cf6a0ce637528c3fd7202b4f7563f49418ff545420a7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.agents/skills/bench-craft-delegate/references/delegation-discipline.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,438\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,4509,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: closure removal: the delegate guidance called the fence a refusal boundary\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t6-t6-guidance-budgets",
          "performer": "claude:ft293_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d790f343b3965cfd4d7b497e659303ee4f78de12",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6",
            "digest": "sha256:b682aaa6a2f3a10908927cd53ce11e0bf9cf598d196444453baf84591b6728d5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t6-t6-prose",
          "performer": "claude:ft293_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d790f343b3965cfd4d7b497e659303ee4f78de12",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6",
            "digest": "sha256:03e28f3a13489358121dd031f36d36810717eec5bf538347fab4655af99549d0",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,107,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-prose",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        },
        {
          "id": "t6r1-t6-anchors-and-help",
          "performer": "claude:ft293_t6_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bc7e5487836f3615df79dfe59687034d95e81e70",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r1",
            "digest": "sha256:12c23a58d2c5c2dc912122c45aa8236ffb4c271ebeaab50560c09e7d707a58b3",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16158,346\n  github.com/gibbonmi/bench/internal/anchors,pass,119,123\n  github.com/gibbonmi/bench/internal/conformance,pass,35396,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t6-anchors-and-help",
          "command": "bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t6r1-t6-root-conformance",
          "performer": "claude:ft293_t6_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bc7e5487836f3615df79dfe59687034d95e81e70",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r1",
            "digest": "sha256:8688daa7db66468f5df6b5a37464c4fe5bddc4a9740a036c8acbc0e6ac7c0851",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14365,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t6r1-t6-guidance-anchors",
          "performer": "claude:ft293_t6_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bc7e5487836f3615df79dfe59687034d95e81e70",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r1",
            "digest": "sha256:165cf671dcbfdfb8d82055cbbacfbabc3c86137f5d0074ba30aa8b153a20bdee",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.agents/skills/bench-craft-delegate/references/delegation-discipline.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,438\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,4334,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: closure removal: the delegate guidance called the fence a refusal boundary\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "Insertion: restore the refusal-boundary sentence in delegation-discipline.md; the forbid anchor must red docs-currency-workflow, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t6_r1",
              "digest": "sha256:165cf671dcbfdfb8d82055cbbacfbabc3c86137f5d0074ba30aa8b153a20bdee",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.agents/skills/bench-craft-delegate/references/delegation-discipline.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,438\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,4334,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: closure removal: the delegate guidance called the fence a refusal boundary\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t6r1-t6-guidance-budgets",
          "performer": "claude:ft293_t6_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bc7e5487836f3615df79dfe59687034d95e81e70",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r1",
            "digest": "sha256:b682aaa6a2f3a10908927cd53ce11e0bf9cf598d196444453baf84591b6728d5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t6r1-t6-prose",
          "performer": "claude:ft293_t6_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bc7e5487836f3615df79dfe59687034d95e81e70",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r1",
            "digest": "sha256:261d8ef1e993aa6910690c2dd6c7e0e7be10292b5893fcc929f24b55eae6c63c",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,100,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-prose",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        },
        {
          "id": "t6r2-t6-anchors-and-help",
          "performer": "claude:ft293_t6_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r2",
            "digest": "sha256:b1b0051e1cb3ddc4137e4beed32a9c5a4901b7e6d3cba054e25b55c52d483956",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15964,346\n  github.com/gibbonmi/bench/internal/anchors,pass,122,123\n  github.com/gibbonmi/bench/internal/conformance,pass,35456,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t6-anchors-and-help",
          "command": "bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t6r2-t6-root-conformance",
          "performer": "claude:ft293_t6_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r2",
            "digest": "sha256:e630612e05335fc118e692e1281836b4d225228deb6c5fa6ef6b21f1e73bf765",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14368,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t6r2-t6-guidance-anchors",
          "performer": "claude:ft293_t6_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r2",
            "digest": "sha256:af0444fba471d4a4c392029d16d872ab12eb2cdd0809b4a3a6168e03046fef19",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.agents/skills/bench-craft-delegate/references/delegation-discipline.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,438\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5032,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: closure removal: the delegate guidance called the fence a refusal boundary\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "Insertion: restore the refusal-boundary sentence in delegation-discipline.md; the forbid anchor must red docs-currency-workflow, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t6_r2",
              "digest": "sha256:af0444fba471d4a4c392029d16d872ab12eb2cdd0809b4a3a6168e03046fef19",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.agents/skills/bench-craft-delegate/references/delegation-discipline.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,438\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5032,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: closure removal: the delegate guidance called the fence a refusal boundary\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t6r2-t6-guidance-budgets",
          "performer": "claude:ft293_t6_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r2",
            "digest": "sha256:b682aaa6a2f3a10908927cd53ce11e0bf9cf598d196444453baf84591b6728d5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t6r2-t6-prose",
          "performer": "claude:ft293_t6_r2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r2",
            "digest": "sha256:ca561e1ddc7357ba55193c1e409322ec143e69a81ec4fd02e21f18bf1574587f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,97,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-prose",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        },
        {
          "id": "t6r3-t6-anchors-and-help",
          "performer": "claude:ft293_t6_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "268b808950fbf7fe79c783f91d7a64f110e719c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r3",
            "digest": "sha256:a656795cf6cdb325194619bb94b6555ae04435f14a9315840097464f9c310572",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16025,346\n  github.com/gibbonmi/bench/internal/anchors,pass,112,123\n  github.com/gibbonmi/bench/internal/conformance,pass,35813,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t6-anchors-and-help",
          "command": "bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t6r3-t6-root-conformance",
          "performer": "claude:ft293_t6_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "268b808950fbf7fe79c783f91d7a64f110e719c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r3",
            "digest": "sha256:8b827b2e2410de975c9506a2fa7e53f7230f3c0ac9e2007d06c706b81cd52043",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14352,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t6r3-t6-guidance-anchors",
          "performer": "claude:ft293_t6_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "268b808950fbf7fe79c783f91d7a64f110e719c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r3",
            "digest": "sha256:f3ea219dc073368311bda7a647dc4d21af9e44bad9b3a6f696f22a006a38f219",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.agents/skills/bench-craft-delegate/references/delegation-discipline.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,438\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5191,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: closure removal: the delegate guidance called the fence a refusal boundary\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-guidance-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "Insertion: restore the refusal-boundary sentence in delegation-discipline.md; the forbid anchor must red docs-currency-workflow, and the restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t6_r3",
              "digest": "sha256:f3ea219dc073368311bda7a647dc4d21af9e44bad9b3a6f696f22a006a38f219",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.agents/skills/bench-craft-delegate/references/delegation-discipline.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,438\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5191,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: closure removal: the delegate guidance called the fence a refusal boundary\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t6r3-t6-guidance-budgets",
          "performer": "claude:ft293_t6_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "268b808950fbf7fe79c783f91d7a64f110e719c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r3",
            "digest": "sha256:b682aaa6a2f3a10908927cd53ce11e0bf9cf598d196444453baf84591b6728d5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-guidance-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "t6r3-t6-prose",
          "performer": "claude:ft293_t6_r3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "268b808950fbf7fe79c783f91d7a64f110e719c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t6_r3",
            "digest": "sha256:5189701af7f4c83164bcfbd5ea1620e69a85a8037aca4b212f55a020f5f15f6e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,102,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-prose",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "cr4-standards",
          "performer": "claude:ft293_cr4_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d790f343b3965cfd4d7b497e659303ee4f78de12",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4_standards",
            "digest": "sha256:9d2596d522b71c671a4592c15313ea911e47b47736ccde513a917931e58edaa1",
            "excerpt": "Standards axis, FT293 CR4 (ffee6a38..688d194f). check-current current at 12d8d89c.\nS1 repair 7: craft-review SKILL.md:42 \"The fence lives in the ticket Writes: lines\" and bench-review-implementation.md:102 \"the union of the ticket Writes: lines\" restate the fence without the implicit entries that craft-spec SKILL.md:65 owns.\nS2 repair 6: bench-implement-spec.md:54 sentence is 23 tokens and joins an instruction and a description.\nS3 judgment 6: bench-implement-spec.md:87 restates the widen scope rule that BENCH.md:124 owns.\nS4 judgment 5: slicing-checks.md:47 names /bench-implement-spec as the widening owner; other files name BENCH.md's plan-expansion policy.\nS5 judgment 6: \"Writes: expansion\" beside \"ticket fence expansions\"; \"widens the line\" beside \"fence amendment\"; \"directory\" beside \"signed-off prefixes\".\nS6 judgment 5: delegation-discipline.md:36 \"missing need\", \"a need outside Writes:\", and \"the line\" are vague.\nS7 judgment 3: BENCH.md:124 inline four-item list and a garden-path sentence.\nAnchor row clean; paragraphs within six sentences; budgets hold. Count 7. Worst: S1.\n"
          },
          "axis": "Standards",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "688d194fa53ca58397bc88636250175f01419041",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "S5",
            "S6",
            "S7"
          ],
          "supersedes": []
        },
        {
          "id": "cr4-spec",
          "performer": "claude:ft293_cr4_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d790f343b3965cfd4d7b497e659303ee4f78de12",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4_spec",
            "digest": "sha256:42ab733ed6d0d609b6a85450feb5dbf1b7bd9f304607765590697e5df7f90853",
            "excerpt": "Spec axis, FT293 CR4 (ffee6a38..688d194f). check-current current=true. bench coverage --check: 94 rows valid.\nP1 should-fix 6: BENCH.md:124 states only the outside half of the scope rule; the in-scope no-stop and its plan commit, a repair's included, rest on implication after the clause deletion; BENCH.md:122 shortfall sentence lost its balance.\nP2 advisory 4: spec.md:324 (story 50, section unchanged) conflicts with spec.md:161 (upper-bound sentence becomes an expectation sentence); the build followed :161.\nRows: CR75 pass (223 ticket paths plus the pickup, zero gap); CR76 pass 8; CR78 pass; CR79 pass; CR88 pass with P1; CR94 pass; stories 40-50 pass; commitment.json ticket half kept.\nProcess note: the axis wrote and deleted two temporary files under /dev/shm against its read-only charge; the tree is unchanged.\nCount 2. Worst: P1.\n"
          },
          "axis": "Spec",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "688d194fa53ca58397bc88636250175f01419041",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "cr4-coverage",
          "performer": "claude:ft293_cr4_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d790f343b3965cfd4d7b497e659303ee4f78de12",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4_coverage",
            "digest": "sha256:3399dbc6535e79df22bf7e6552537c7691ccb6dea382bf7ccaa9a2148f64b1e2",
            "excerpt": "Coverage axis, FT293 CR4 (ffee6a38..688d194f). check-current current=true at 12d8d89c; tree clean.\nProbes: line break inside the phrase bit; **refusal boundary** silent; \"Is A Refusal Boundary\" silent; --propose-writes step in bench-implement-spec.md silent; delegation line 92 old wording silent (accepted, rule wording); craft-review reflow bit (123 > 122). All restored=yes.\nC1 advisory 8: the CR78 row uses Forbid, which keeps case and emphasis; ForbidCaseFoldedEmphasis exists (match.go:18).\nC2 advisory 6: --propose-writes forbid rows cover three named files only.\nSweep: no removed form in .agents, README.md, AGENTS.md, CONTEXT.md; docs/research and docs/audits hits are dated evidence; commitment.json hits are record slugs.\nCount 2. Worst: C1.\n"
          },
          "axis": "Coverage",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "688d194fa53ca58397bc88636250175f01419041",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "cr4c-standards",
          "performer": "claude:ft293_cr4c_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "bc7e5487836f3615df79dfe59687034d95e81e70",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4c_standards",
            "digest": "sha256:827ae06ff7798cab4e26697c3c1e2baa34b5bd13ca81e000f340d688f9dc6e14",
            "excerpt": "Standards confirming round, FT293 CR4, repair delta 688d194f..093b9a75. check-current current=true at 942c3980.\nS1 closed 9; S2 closed 8; S3 closed 7; S4 closed 7; S6 closed 8; S7 closed 7.\nS5 open 6: one widening still has several names (BENCH.md:124 \"expand ticket Writes: expectations\" beside :126 \"widens the Writes: line\"; bench-implement-spec.md:93 \"ticket fence expansions\"; delegation-discipline.md:38 \"fence amendment\"; bench-implement-spec.md:87 \"gives to the reviewer\" beside BENCH.md \"requires my decision\").\nS8 auto-fix 6: BENCH.md:126 \"inside the directory of each signed-off Writes: entry\" reads as every entry's directory; the scope is the union (widen.go:38, :227); \"a repair's path does too\" reads as unconditional.\nS9 auto-fix 4: the BENCH.md:126 lead-in \"Each of these changes ... cites for a path:\" does not match list item :131.\nS10 advisory 3: bench-implement-spec.md:54 \"Before using an in-scope plan\" uses a gerund as a noun.\nC1 test side acceptable 6: the DG15 selection predicate appears in two test packages; a needle rename fails loudly.\nOpen: S5, S8, S9. Worst: S8.\n"
          },
          "axis": "Standards",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "093b9a758c7ab506db65dc6809d77c65143ebc25",
          "finding_ids": [
            "S5",
            "S8",
            "S9"
          ],
          "supersedes": [
            "cr4-standards"
          ]
        },
        {
          "id": "cr4c-spec",
          "performer": "claude:ft293_cr4c_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "bc7e5487836f3615df79dfe59687034d95e81e70",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4c_spec",
            "digest": "sha256:b3ef505b055946144e924d495aa203f808d2a388675bb13952201c01a0d30f26",
            "excerpt": "Spec confirming round, FT293 CR4, repair delta 688d194f..093b9a75. check-current current=true at 942c3980. bench coverage --check: 94 rows valid.\nP1 held 9: BENCH.md:126 names bench spec widen and the in-scope plan commit with no stop, a repair's included; :131 lists the outside path; --decision cited.\nCR79 held 9. CR75 held 9 (226 entries equal the union plus the pickup). CR78 held 8. CR94 held 8. No unrequested behavior 8.\nAdvice: spec.md:163 says /bench-implement-spec names the verb at the plan-expansion sentence; after the S3 repair it names the verb at the pause and points to BENCH.md (one route, one hop). Flagged for reviewer veto.\nP3 advisory 3: BENCH.md:126 \"each signed-off entry\" can read as every entry.\nOpen 0.\n"
          },
          "axis": "Spec",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "093b9a758c7ab506db65dc6809d77c65143ebc25",
          "finding_ids": [],
          "supersedes": [
            "cr4-spec"
          ]
        },
        {
          "id": "cr4c-coverage",
          "performer": "claude:ft293_cr4c_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "bc7e5487836f3615df79dfe59687034d95e81e70",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4c_coverage",
            "digest": "sha256:2d503c42ce58b62b7ccf2e75e85faa19ecd92fe4e063f85a7e7dc429c3600bbd",
            "excerpt": "Coverage confirming round, FT293 CR4. check-current current=true at 942c3980; tree clean.\nC1 closed 8: probes on delegation-discipline.md under docs-currency-workflow: is a *refusal* _boundary_ bit; IS A REFUSAL BOUNDARY bit; soft line break inside emphasis bit; whole phrase in a code span bit; code span over part of the phrase silent. All restored=yes.\nDG15 needle swapped in registry_debug_loop.go:16: all three DG15-selecting tests red at their lookup; restored.\nC2 advisory 6: a code span over part of a needle passes a case-folded forbid by matcher design (DG15 shares it); park via bench learning.\nOpen 0.\n"
          },
          "axis": "Coverage",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "093b9a758c7ab506db65dc6809d77c65143ebc25",
          "finding_ids": [],
          "supersedes": [
            "cr4-coverage"
          ]
        },
        {
          "id": "cr4c2-standards",
          "performer": "claude:ft293_cr4c2_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4c2_standards",
            "digest": "sha256:f2da5db5f437a3c0bd6e70888325f5dc1251c6b5bc44bcb4b3840c273a209ef5",
            "excerpt": "Standards second confirming round, FT293 CR4, repair delta 093b9a75..b0d84d87. check-current current=true, dirty_record=false.\nS5 closed 7: one name \"widen the Writes: line\" / \"Writes: widening\" (BENCH.md:124, :126; bench-implement-spec.md:54, :87, :93; delegation-discipline.md:36, :38); the decision reads \"requires my decision\" / \"requires the reviewer's decision\".\nS8 closed 7: BENCH.md:126 \"a signed-off Writes: entry\"; \"This rule holds for a repair too.\" keeps the condition.\nS9 closed 7: the lead-in fits items :128-:131 (21 words).\nAdvice: BENCH.md:124 \"Before using a plan expansion\" uses a gerund as a noun; bench-write-spec.md:71 and slicing-checks.md:36 keep \"fence expansion\" (outside the delta).\nOpen 0.\n"
          },
          "axis": "Standards",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "b0d84d879e53da93cf2cdac833efb395629fc316",
          "finding_ids": [],
          "supersedes": [
            "cr4c-standards"
          ]
        },
        {
          "id": "cr4c2-spec",
          "performer": "claude:ft293_cr4c2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4c2_spec",
            "digest": "sha256:0f9cdf65bd7f780c476758c456a010befc76c2cee411f179e37be3a8075793ad",
            "excerpt": "Spec second confirming round, FT293 CR4, repair delta 093b9a75..b0d84d87. check-current current at fb7d0ffe. bench coverage --check: 94 rows valid.\nCR79 held 9; CR88 held 9; scope wording matches spec.md:121 held 9; CR94 held 9 (net zero lines); CR75 held 7; no rule change beyond the three wordings 8.\nAdvice: spec.md:27 and story 18 keep \"the directory of each signed-off entry\".\nOpen 0.\n"
          },
          "axis": "Spec",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "b0d84d879e53da93cf2cdac833efb395629fc316",
          "finding_ids": [],
          "supersedes": [
            "cr4c-spec"
          ]
        },
        {
          "id": "cr4c2-coverage",
          "performer": "claude:ft293_cr4c2_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4c2_coverage",
            "digest": "sha256:7e3532b04bf0bf1810fd82b5adc89312f347d47d3fc0223d9f870f49e6ab8a1b",
            "excerpt": "Coverage second confirming round, FT293 CR4. check-current current at fb7d0ffec; tree clean.\nProbes: one added line in bench-implement-spec.md, guidance-prose-budgets: red (94 over 93), restored yes; line plus blank line: red (95), restored yes; restore the BENCH.md forbid needle (anchor row 32), docs-currency-workflow: red, restored yes; lengthen BENCH.md:126 sentence to 27 words, prose-mechanics: red, restored yes.\nNo new gap; no guard removed. Open 0.\n"
          },
          "axis": "Coverage",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "b0d84d879e53da93cf2cdac833efb395629fc316",
          "finding_ids": [],
          "supersedes": [
            "cr4c-coverage"
          ]
        },
        {
          "id": "cr4c3-standards",
          "performer": "claude:ft293_cr4c3_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "268b808950fbf7fe79c783f91d7a64f110e719c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4c3_standards",
            "digest": "sha256:afa38a3894b8e0fcb86a0ca5d3f798824ebdfb56b76f1cdc455cd6e5ce0ea734",
            "excerpt": "Standards third confirming round, FT293 CR4, delta b0d84d87..1b79d041 (decision-map source). check-current current.\nClosed 8: the new source delegation-discipline.md:36 states the Writes: expectation rule that the entry supports; the entry keeps the Path/Supports/Drift shape; no live citation of roadmap/FT293 remains outside the outcome pin and dated research.\nFindings: none.\n"
          },
          "axis": "Standards",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "1b79d041f615d1e7e8f322c98245c0de6ad7def7",
          "finding_ids": [],
          "supersedes": [
            "cr4c2-standards"
          ]
        },
        {
          "id": "cr4c3-spec",
          "performer": "claude:ft293_cr4c3_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "268b808950fbf7fe79c783f91d7a64f110e719c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4c3_spec",
            "digest": "sha256:91788c84c527fba6ddb4e1ed6f5867f7ff7334fd612140919f71165d09acbd39",
            "excerpt": "Spec third confirming round, FT293 CR4, delta b0d84d87..1b79d041. check-current current at 4cd12990. bench coverage --check: 94 rows valid.\nNo story, decision, or coverage row changed held 9; CR75 section equals the union plus the pickup (227 paths) held 9; the decision-map edit stays in scope held 8.\nAdvice: the widening of decisions/architecture-planning.md cites auto-approve:ft293-full-20261018; flagged for reviewer veto.\nProcess note: the axis wrote and deleted two temporary files under /dev/shm against its read-only charge; the tree is unchanged.\nOpen 0.\n"
          },
          "axis": "Spec",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "1b79d041f615d1e7e8f322c98245c0de6ad7def7",
          "finding_ids": [],
          "supersedes": [
            "cr4c2-spec"
          ]
        },
        {
          "id": "cr4c3-coverage",
          "performer": "claude:ft293_cr4c3_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "268b808950fbf7fe79c783f91d7a64f110e719c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_cr4c3_coverage",
            "digest": "sha256:b282cf71796f5b41336d586a5b124bea08e140b2aa850e78b4084ffe3f07eeb1",
            "excerpt": "Coverage third confirming round, FT293 CR4. check-current current at 4cd129908.\nProbe: decisions/architecture-planning.md source swapped to roadmap/FT293-missing.md, check decision-map-integrity: bit (TestRootConformance, \"Sources Path roadmap/FT293-missing.md must name an existing regular file\"); restored yes; tree clean.\nFindings: none.\n"
          },
          "axis": "Coverage",
          "base": "ffee6a38c248bb4299ed5cbc6ff634a0e4ee5de4",
          "tip": "1b79d041f615d1e7e8f322c98245c0de6ad7def7",
          "finding_ids": [],
          "supersedes": [
            "cr4c2-coverage"
          ]
        }
      ]
    }
  ],
  "completion": {
    "state": "completed",
    "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
    "performer": "claude:ft293-orchestrator-20261018",
    "reconciliation": {
      "CR01": "covered",
      "CR02": "covered",
      "CR03": "covered",
      "CR04": "covered",
      "CR05": "covered",
      "CR06": "covered",
      "CR07": "covered",
      "CR08": "covered",
      "CR09": "covered",
      "CR10": "covered",
      "CR11": "covered",
      "CR12": "covered",
      "CR13": "covered",
      "CR14": "covered",
      "CR15": "covered",
      "CR16": "covered",
      "CR17": "covered",
      "CR18": "covered",
      "CR19": "covered",
      "CR20": "covered",
      "CR21": "covered",
      "CR22": "covered",
      "CR23": "covered",
      "CR24": "covered",
      "CR25": "covered",
      "CR26": "covered",
      "CR27": "covered",
      "CR28": "covered",
      "CR29": "covered",
      "CR30": "covered",
      "CR31": "covered",
      "CR32": "covered",
      "CR33": "covered",
      "CR34": "covered",
      "CR35": "covered",
      "CR36": "covered",
      "CR37": "covered",
      "CR38": "covered",
      "CR39": "covered",
      "CR40": "covered",
      "CR41": "covered",
      "CR42": "covered",
      "CR43": "covered",
      "CR44": "covered",
      "CR45": "covered",
      "CR46": "covered",
      "CR47": "covered",
      "CR48": "covered",
      "CR49": "covered",
      "CR50": "covered",
      "CR51": "covered",
      "CR52": "covered",
      "CR53": "covered",
      "CR54": "covered",
      "CR55": "covered",
      "CR56": "covered",
      "CR57": "covered",
      "CR58": "covered",
      "CR59": "covered",
      "CR60": "covered",
      "CR61": "covered",
      "CR62": "covered",
      "CR63": "covered",
      "CR64": "covered",
      "CR65": "covered",
      "CR66": "covered",
      "CR67": "covered",
      "CR68": "covered",
      "CR69": "covered",
      "CR70": "covered",
      "CR71": "covered",
      "CR72": "covered",
      "CR75": "covered",
      "CR76": "covered",
      "CR77": "covered",
      "CR78": "covered",
      "CR79": "covered",
      "CR80": "covered",
      "CR81": "covered",
      "CR82": "covered",
      "CR83": "covered",
      "CR84": "covered",
      "CR85": "covered",
      "CR86": "covered",
      "CR87": "covered",
      "CR88": "covered",
      "CR89": "covered",
      "CR90": "covered",
      "CR91": "covered",
      "CR92": "covered",
      "CR93": "covered",
      "CR94": "covered",
      "CR95": "covered",
      "CR96": "covered"
    },
    "verification": [
      {
        "id": "final-preflight-packages",
        "performer": "claude:ft293-orchestrator-20261018",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude:ft293-orchestrator-20261018",
          "digest": "sha256:955603f5b412a1f970e8f29c5c3ef78defb1034d43bf1c37d3fd473972e6815a",
          "excerpt": "tree[1]{target,head,dirty}:\n  closure-removal,503d23fa3e67b4f6f4553af45db7841a81322d04,false\npackages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,pass,7059,101\n  github.com/gibbonmi/bench/internal/preflight,pass,29913,369\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,21224,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,11811,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5425,70\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
        },
        "requirement": "preflight-packages",
        "command": "bench test --package ./internal/preflight --package ./internal/preflight/preflighttest --package ./internal/preflight/ownership --package ./internal/preflight/ownership/ownershiptest --package ./internal/preflight/staleness --package ./internal/preflight/evidencecmd --package ./internal/diff",
        "exit_code": 0
      },
      {
        "id": "final-coverage",
        "performer": "claude:ft293-orchestrator-20261018",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude:ft293-orchestrator-20261018",
          "digest": "sha256:a116dbbcfb3a6b57040e63bead84b5d2e716185e4d9f74819295150e58d9e457",
          "excerpt": "tree[1]{target,head,dirty}:\n  closure-removal,503d23fa3e67b4f6f4553af45db7841a81322d04,false\nok: coverage map valid — 94 row(s)\n"
        },
        "requirement": "coverage",
        "command": "bench coverage --check specs/closure-removal/spec.md",
        "exit_code": 0
      },
      {
        "id": "final-fence-readers",
        "performer": "claude:ft293-orchestrator-20261018",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude:ft293-orchestrator-20261018",
          "digest": "sha256:97023a0b2e7f12cd52617ef1f4dd9c3cad7d0484de5eaad03f6d0584f064fede",
          "excerpt": "tree[1]{target,head,dirty}:\n  closure-removal,503d23fa3e67b4f6f4553af45db7841a81322d04,true\npackages[10]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17182,346\n  github.com/gibbonmi/bench/internal/anchors,pass,119,123\n  github.com/gibbonmi/bench/internal/canary,pass,8,32\nspilled{lines=15,bytes=1093,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/ade8f4556745da60d1dc6299d63291b4/1792400883741179911-b81c7d19aedd10d8.out}\n  github.com/gibbonmi/bench/internal/tickets,pass,3,45\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/chargeevidence,TestEvidenceStoreKinds/CE94_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
        },
        "requirement": "fence-readers",
        "command": "bench test --package ./internal/coverage --package ./internal/spec --package ./internal/chargeevidence --package ./internal/refusalroute --package ./internal/reviewrecord/recordtest --package ./internal/tickets --package ./internal/anchors --package ./internal/canary --package ./internal/conformance --package ./cmd/bench",
        "exit_code": 0
      },
      {
        "id": "final-landing",
        "performer": "claude:ft293-orchestrator-20261018",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude:ft293-orchestrator-20261018",
          "digest": "sha256:7378390cb4175f32a1d4988408a97a27c84a29b44f45a05f9dd15be0af395734",
          "excerpt": "tree[1]{target,head,dirty}:\n  closure-removal,503d23fa3e67b4f6f4553af45db7841a81322d04,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,66972,1439\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
        },
        "requirement": "landing",
        "command": "bench test --package ./internal/worktree",
        "exit_code": 0
      },
      {
        "id": "final-system",
        "performer": "claude:ft293-orchestrator-20261018",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude:ft293-orchestrator-20261018",
          "digest": "sha256:83d32b800b44f442f5602e1b1da8cd0959f64c28ed3449cbb23f13229c72eb09",
          "excerpt": "tree[1]{target,head,dirty}:\n  closure-removal,503d23fa3e67b4f6f4553af45db7841a81322d04,true\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,170,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,70231,170\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
        },
        "requirement": "system",
        "command": "bench test --check system",
        "exit_code": 0
      },
      {
        "id": "final-guidance-anchors",
        "performer": "claude:ft293-orchestrator-20261018",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude:ft293-orchestrator-20261018",
          "digest": "sha256:edc303f523201c8b02e442f952ba6bb7696a770d946393ff8e728efc33670473",
          "excerpt": "tree[1]{target,head,dirty}:\n  closure-removal,503d23fa3e67b4f6f4553af45db7841a81322d04,true\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4602,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
        },
        "requirement": "guidance-anchors",
        "command": "bench test --check docs-currency-workflow",
        "exit_code": 0
      },
      {
        "id": "final-guidance-budgets",
        "performer": "claude:ft293-orchestrator-20261018",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude:ft293-orchestrator-20261018",
          "digest": "sha256:79e2e0e53766f25e160661411d2645ca449b5db49807ba530c0c9a712071c580",
          "excerpt": "tree[1]{target,head,dirty}:\n  closure-removal,503d23fa3e67b4f6f4553af45db7841a81322d04,true\ncheck[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
        },
        "requirement": "guidance-budgets",
        "command": "bench test --check guidance-prose-budgets",
        "exit_code": 0
      },
      {
        "id": "final-prose",
        "performer": "claude:ft293-orchestrator-20261018",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude:ft293-orchestrator-20261018",
          "digest": "sha256:20cfd207b128ea8ebe7543a0b22333b8ab7a46ad7e7a981f691108a24c3c8ea3",
          "excerpt": "tree[1]{target,head,dirty}:\n  closure-removal,503d23fa3e67b4f6f4553af45db7841a81322d04,true\ncheck[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,100,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
        },
        "requirement": "prose",
        "command": "bench test --check prose-mechanics",
        "exit_code": 0
      },
      {
        "id": "final-root-conformance",
        "performer": "claude:ft293-orchestrator-20261018",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude:ft293-orchestrator-20261018",
          "digest": "sha256:c2124c2e8c8ce0e6273e36f83b72e6d5e690a8ad6034750a91c8ced5c4e0ff71",
          "excerpt": "tree[1]{target,head,dirty}:\n  closure-removal,503d23fa3e67b4f6f4553af45db7841a81322d04,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15529,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
        },
        "requirement": "root-conformance",
        "command": "bench test --check root-conformance",
        "exit_code": 0
      },
      {
        "id": "final-coordinator-section-reader",
        "performer": "claude:ft293-orchestrator-20261018",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "3c16fb6b9dd0fa1f5eaeb445ab00f41f892c058b",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude:ft293-orchestrator-20261018",
          "digest": "sha256:6c548146e7b5c7bef68e9a2635c56543683864e2ffcd94e426588666a5258b62",
          "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/paths_authorized.go,swap,failed,2,TestCommandBuildIgnoresAFenceSection,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,all,passed,369\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,878\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,25008,369\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestCommandBuildIgnoresAFenceSection,\"command_fence_test.go:129: build = (1) with paths-authorized {Check:paths-authorized Verdict:red Detail:not authorized by any ownership fence: reviews/other.md Next: exempt:false}, want the red \\\"not authorized by any ownership fence: intern… (276 bytes)\",6\n  github.com/gibbonmi/bench/internal/preflight,TestCommandReviewAuthorizesFromTicketWrites,\"command_fence_test.go:139: review = (1) with paths-authorized {Check:paths-authorized Verdict:red Detail:not authorized by any ownership fence: reviews/other.md Next: exempt:false}, want the red \\\"not authorized by any ownership fence: inter… (277 bytes)\",6\nskips[0]{package,test,reason}:\n"
        },
        "requirement": "coordinator-section-reader",
        "command": "bench test --package ./internal/preflight",
        "exit_code": 0,
        "probe": {
          "mutation": "append the leftover fence-section tokens to the authorizing entries; TestCommandBuildIgnoresAFenceSection must fail",
          "outcome": "bit",
          "exit_code": 1,
          "restore": "pass",
          "native_ref": {
            "ref": "claude:ft293-orchestrator-20261018",
            "digest": "sha256:6c548146e7b5c7bef68e9a2635c56543683864e2ffcd94e426588666a5258b62",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/paths_authorized.go,swap,failed,2,TestCommandBuildIgnoresAFenceSection,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,all,passed,369\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,878\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,25008,369\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestCommandBuildIgnoresAFenceSection,\"command_fence_test.go:129: build = (1) with paths-authorized {Check:paths-authorized Verdict:red Detail:not authorized by any ownership fence: reviews/other.md Next: exempt:false}, want the red \\\"not authorized by any ownership fence: intern… (276 bytes)\",6\n  github.com/gibbonmi/bench/internal/preflight,TestCommandReviewAuthorizesFromTicketWrites,\"command_fence_test.go:139: review = (1) with paths-authorized {Check:paths-authorized Verdict:red Detail:not authorized by any ownership fence: reviews/other.md Next: exempt:false}, want the red \\\"not authorized by any ownership fence: inter… (277 bytes)\",6\nskips[0]{package,test,reason}:\n"
          }
        }
      }
    ]
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
    },
    {
      "from": "sha256:b8aa3bef62333ca44c9dcab717464e45f80187bdce5c374ceb36ba739453eb66",
      "to": "sha256:966368af4476a968b87a85ea361a5a016d608290d3dcd83c99abe2ecf60e065c",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ],
        "CR2": [
          "CR2"
        ],
        "CR3": [
          "CR3"
        ]
      }
    },
    {
      "from": "sha256:966368af4476a968b87a85ea361a5a016d608290d3dcd83c99abe2ecf60e065c",
      "to": "sha256:6c7a65800371afdd0e0250f01b0d90b912a017eb591baca4f30a76d7407ee411",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ],
        "CR2": [
          "CR2"
        ],
        "CR3": [
          "CR3"
        ]
      }
    },
    {
      "from": "sha256:6c7a65800371afdd0e0250f01b0d90b912a017eb591baca4f30a76d7407ee411",
      "to": "sha256:4f5ce59424ab22825cec7bcbe7cdba51aa3a1b9492032e7425a380b332e9921b",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ],
        "CR2": [
          "CR2"
        ],
        "CR3": [
          "CR3"
        ]
      }
    },
    {
      "from": "sha256:4f5ce59424ab22825cec7bcbe7cdba51aa3a1b9492032e7425a380b332e9921b",
      "to": "sha256:7f97a4d2c2010851cefe2b2cab6f1c3d0a903c769322652906b9928fcf7d0045",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ],
        "CR2": [
          "CR2"
        ],
        "CR3": [
          "CR3"
        ],
        "CR4": [
          "CR4"
        ]
      }
    },
    {
      "from": "sha256:7f97a4d2c2010851cefe2b2cab6f1c3d0a903c769322652906b9928fcf7d0045",
      "to": "sha256:9cfd3285582a56265c7b6a6bf4e6e7c4e3480a2f5b899c5eaab9598d0fb28385",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ],
        "CR2": [
          "CR2"
        ],
        "CR3": [
          "CR3"
        ],
        "CR4": [
          "CR4"
        ]
      }
    },
    {
      "from": "sha256:9cfd3285582a56265c7b6a6bf4e6e7c4e3480a2f5b899c5eaab9598d0fb28385",
      "to": "sha256:40e3bcd3128efe554b1c3045d1e51f4107e10d7f6f8c9590d56d82e1cb1420f9",
      "chunk_ids": {
        "CR1": [
          "CR1"
        ],
        "CR2": [
          "CR2"
        ],
        "CR3": [
          "CR3"
        ],
        "CR4": [
          "CR4"
        ]
      }
    }
  ]
}
```
