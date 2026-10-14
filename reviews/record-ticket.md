# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/record-ticket/spec.md",
  "plan_digest": "sha256:e8251dd1728182f946187f20d4eaba706b5aef43e7b9466d453dc0a7f24489e1",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RT-RETAIN",
      "base": "53e60dc58f0432fdadf2b31b78d4e7a25f85570d",
      "tip": "9d989e126b4cbeac507b12ea3137250fd9c6b294",
      "plan_digest": "sha256:e8251dd1728182f946187f20d4eaba706b5aef43e7b9466d453dc0a7f24489e1",
      "source_digest": "e3239b89055da6012a766088982793175b60df59",
      "acceptance_rows": [
        "RT1",
        "RT2",
        "RT3",
        "RT4",
        "RT5",
        "RT6",
        "RT7",
        "RT8",
        "RT65"
      ],
      "verification": [
        {
          "id": "ft405_t01-retention",
          "performer": "claude:ft405_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft405_t01",
            "digest": "sha256:478e609c4e4b03228b63522691ab5613be1bce9db5810171f2c3779a005b2cde",
            "excerpt": "packages[11]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18816,345\n  github.com/gibbonmi/bench/internal/conformance,pass,81955,1694\n  github.com/gibbonmi/bench/internal/gate,pass,47840,489\n  github.com/gibbonmi/bench/internal/git,pass,1417,172\n  github.com/gibbonmi/bench/internal/handoff,pass,2960,56\n  github.com/gibbonmi/bench/internal/responsebound,pass,266,61\n  github.com/gibbonmi/bench/internal/responsebound/responseboundtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/roadmap,pass,1905,176\n  github.com/gibbonmi/bench/internal/shift,pass,6375,116\n  github.com/gibbonmi/bench/internal/status,pass,8259,142\n  github.com/gibbonmi/bench/internal/worktree,pass,77044,1433\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "retention",
          "command": "bench test --package ./cmd/bench --package ./internal/responsebound/... --package ./internal/git --package ./internal/gate --package ./internal/handoff --package ./internal/roadmap --package ./internal/shift --package ./internal/status --package ./internal/worktree --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft405_t01-retain-root",
          "performer": "claude:ft405_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft405_t01",
            "digest": "sha256:2a33e4364ef1af5ff2ab111c6dbf515049beae83a45ab8986775641c82cbd5b1",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft405-build,9d989e126b4cbeac507b12ea3137250fd9c6b294,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15332,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "retain-root",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft405_t01-nested-exec",
          "performer": "claude:ft405_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft405_t01",
            "digest": "sha256:d01d525aa63ccb3d71293d82464fc70e3ad0a218a4800ce4fd159c3776a9262d",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft405-build,9d989e126b4cbeac507b12ea3137250fd9c6b294,true\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,2,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,539,2\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "nested-exec",
          "command": "bench test --check system --run TestExecNested",
          "exit_code": 0
        },
        {
          "id": "ft405_t01-record-skip",
          "performer": "claude:ft405_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft405_t01",
            "digest": "sha256:b048b8cab438469381e9018b42eb30add7d60968249659c6b570e8c9a518fad2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,cmd/bench/retain.go,omit,failed,1,TestRetainSkipsRecordWorktreeAndRetiringCalls/record,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,all,passed,345\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,527\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,14822,345\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestRetainSkipsRecordWorktreeAndRetiringCalls/record,\"retain_test.go:102: retained after [\\\"record\\\" \\\"chunk\\\"] = {Verb:record Args:[chunk] Exit:0 Head:726829397659fa99a2294cdda885d1cf0a718bd0 Dirty:[] Output:[108 105 110 101 32 48 49 10 108 105 110 101 32 48 50 10]}, want the test response {Verb:… (352 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "record-skip",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0,
          "probe": {
            "mutation": "drop the record-family skip from retention",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft405_t01",
              "digest": "sha256:b048b8cab438469381e9018b42eb30add7d60968249659c6b570e8c9a518fad2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,cmd/bench/retain.go,omit,failed,1,TestRetainSkipsRecordWorktreeAndRetiringCalls/record,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,all,passed,345\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,527\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,14822,345\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestRetainSkipsRecordWorktreeAndRetiringCalls/record,\"retain_test.go:102: retained after [\\\"record\\\" \\\"chunk\\\"] = {Verb:record Args:[chunk] Exit:0 Head:726829397659fa99a2294cdda885d1cf0a718bd0 Dirty:[] Output:[108 105 110 101 32 48 49 10 108 105 110 101 32 48 50 10]}, want the test response {Verb:… (352 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ft405_rv_retain-standards",
          "performer": "claude:ft405_rv_retain_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft405_rv_retain_standards",
            "digest": "sha256:c0ec1af7b003186005dd6f637f8478f5e1276bf145ffaee1753e44955e659a43",
            "excerpt": "## Standards\n\nI found 4 findings: 1 hard and 3 judgment. The worst is S1, where the response-store scope layout is now built in two places.\n\n- **S1, hard.** `internal/responsebound/retain.go:40` builds `levels := []string{storeDirName, poolkey.Key(root), scope}` again. `openSpill` already builds the same layout at `internal/responsebound/store.go:46-51`, through `location`, which already computes the key. Rule: AGENTS.md \"one source per fact\". The two copies already differ: `openSpill` runs `filepath.Abs(home)` and `sanitize.LineSafe`, and `retainedPath` runs neither. Disposition: auto-fix, by extracting one scope-path helper in `store.go`, which is inside the ticket's Writes. Confidence 6.\n- **S2, judgment.** `internal/shift/shift.go:53` still runs its own `status --porcelain -z --no-renames` read and passes the result to `ParsePorcelainZ`, beside the new `git.PorcelainStatus` at `internal/git/status.go:262`. Rule: AGENTS.md \"one source per fact\". Spec line 150 claims that only \"`LandedState` and `Facts` already run that read\", so moving shift's read is outside the acceptance. Shift also ignores the read error on purpose (audit #11). Disposition: ask-user. Confidence 6.\n- **S3, judgment.** `internal/responsebound/owner.go:78-81, 215-217` adds `whole` and `replayed` (smell: Mutable Data). They store again what the old phase and `retained` already held: Finish copies each write into a second buffer instead of reading `retained` before the reset, as the ticket says (ticket line 16). `replayed` also has no field comment, unlike the fields around it. Disposition: auto-fix. Confidence 4.\n- **S4, judgment.** The doc comment at `internal/responsebound/retain.go:389-390` says \"the complete output of the finished response o and the call that call answers\". The repeated word \"call\" is ambiguous (craft-comments register, ASD-STE100 clarity). Disposition: auto-fix. Confidence 5.\n\n### Advice\n- `cmd/bench/retain.go:14` skips commands by name instead of reading a fact from the registry. Ticket line 43 chose a name comparison for line headroom, so I did not file it (a finding would be a no-op). The comment at lines 52-55 justifies the whole `worktree` family by `exec` alone, but the skip also covers `land`, `merge`, and `reset`.\n- Moved code keeps its error handling: `LandedState` still wraps with `fmt.Errorf(\"git status %s: %w\")` (`internal/git/status.go:66-68`), and `Facts` still returns the error unwrapped, as before.\n\n### What I examined\nThe full diff `53e60dc5..9d989e12`, the ticket, spec lines 143-170 and the write list, `owner.go`, `store.go`, `retire.go`, `size.go`, `status.go`, `shift.go:30-55`, `land_identity.go:330-354`, and the registry structs. A search of production Go for `status --porcelain` and `ParsePorcelainZ(` callers. The comments and the test fixtures. The preflight evidence check reported `current=true`. No files written; no tests or probes run.\n\n### Did the implementation command cause any issue?\nNot the command itself; the spec caused S2. A seam that moves a shared read should cite the full search over its callers.\n"
          },
          "axis": "Standards",
          "base": "53e60dc58f0432fdadf2b31b78d4e7a25f85570d",
          "tip": "9d989e126b4cbeac507b12ea3137250fd9c6b294",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4"
          ],
          "supersedes": []
        },
        {
          "id": "ft405_rv_retain-spec",
          "performer": "claude:ft405_rv_retain_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft405_rv_retain_spec",
            "digest": "sha256:f46fa028e38b19f13479a95ee65efcb84d62ca69f335fb98f2eddf2c38ec2ad7",
            "excerpt": "## Spec\n\nI found 2 findings, both low severity. The worst one is P1: RT7 and RT65 are tested at the `Owner.Retain` seam instead of through the dispatcher that the testing decision names. I judge that deviation acceptable. Rows RT1–RT8 and RT65 are otherwise fully met, and each seam cell names a real test that turns red under the row's named mutation.\n\n- **P1** Spec line 311: \"The retention tests drive one bounded call through the command dispatcher in a fixture pool worktree, then read the retained response.\" Where: `internal/responsebound/retain_test.go:467`, `TestRetainClearsAnIncompleteResponse`. Gap: RT7 and RT65 call `Owner.Retain` directly; the dispatcher cannot reach a failed spill without the test-only hooks `SetCreate` and `FaultAfter` (`internal/responsebound/export_test.go:24,28`). The named mutations do fail at this seam. Contradicts the spec text but not its behavior; follow the tree and flag for reviewer veto. Disposition: no-op. Confidence 7.\n- **P2** Spec line 167: \"A response with no complete bytes clears the retained response.\" Where: `internal/responsebound/retain.go:405-407`. A failed `HEAD` or porcelain read also clears the retained response (subtest `retain_test.go:484`). Small unstated behavior, defensible. Disposition: no-op. Confidence 5.\n\nPassed checks: the four skip classes (`cmd/bench/retain.go:56-58`, `retain.go:394-396`); complete bytes only in the `retaining` and `spilling` phases (`owner.go:211-253`); `HEAD` and dirty paths read only after `completeOutput` succeeds (`retain.go:401-405`); `PorcelainStatus` (`status.go:270`) is the one shared reader that `LandedState` and `Facts` call; one `durablefile.Replace` whose error the dispatcher discards (`census_output.go:29`); the file sits under the assignment scope so `Drop` removes it, and `retained.json` does not match the spill pattern.\n\n## Optional advice\n- `retainedPath` calls `poolkey.Key` twice, so a primary-checkout call pays one `git rev-parse --git-common-dir` (spec line 153 reads loosely).\n- The literals `\"record\"` and `\"worktree\"` in `retainsResponse` repeat registry names (`cmd/bench/main.go:113,132`).\n\n## Examined\nPreflight evidence current at 6c34bbc4; the diff 53e60dc5..9d989e12; `bench coverage`; spec stories 1–10, the Retained response decision, testing decisions, coverage map, edge inventory; ticket 01; `owner.go`, `store.go`, `size.go`, `poolkey.go`, `git.go`; the `worktree exec` leaf; the system-test helpers. No tests or probes run. The implementation command contributed to neither finding.\n"
          },
          "axis": "Spec",
          "base": "53e60dc58f0432fdadf2b31b78d4e7a25f85570d",
          "tip": "9d989e126b4cbeac507b12ea3137250fd9c6b294",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ft405_rv_retain-coverage",
          "performer": "claude:ft405_rv_retain_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft405_rv_retain_coverage",
            "digest": "sha256:79f8290c917c3d1317df54161726ffb6dbe4af2a7c4f882e80bc25c5d48ff461",
            "excerpt": "## Coverage\n\nFindings: 2. The worst gap is that no test runs two complete retaining calls back to back, so nothing proves the second response replaces the first. A write-once mutant passes every test. That defeats story 1 (\"the complete last response\") and the edge-inventory line \"the last retention write wins\" (spec.md:418).\n\n- **C1.** State: a retained file from an earlier complete call exists when a second complete `bench test` call runs. Expected break: the store keeps the first response, so a later `bench record excerpt` stages an old run. Each RT1, RT4, and RT8 write lands on an absent file, and the RT7/RT65 test clears rather than replaces. Test that should exist: an RT1 case in `cmd/bench/retain_test.go` that runs `test --a` (exit 3) and then `test --b` (exit 0), and asserts the second call is retained. Disposition: auto-fix. Confidence 8.\n- **C2.** Input: a response within the bound that interleaves stdout and stderr writes. Expected break: only stdout reaches the retained Output while a spilled response keeps both streams, so spec.md:147 (\"complete response bytes in arrival order\") fails. Every retention test writes stdout only (`exitingHandler` at `cmd/bench/response_bound_test.go:80`, `StdoutLines` in `internal/responsebound/retain_test.go`). Test that should exist: an RT1 within-bound case whose handler interleaves the two streams, with an exact assertion on Output order. Disposition: auto-fix. Confidence 7.\n\n## Probes (independent of the author's probes; git status clean afterwards)\n1. Write-once retention: in `internal/responsebound/retain.go`, insert `if _, e := os.Stat(path); e == nil { return nil };` before `return durablefile.Replace(...)`; packages `./cmd/bench` and `./internal/responsebound/...`. Verdict: silent (406 tests ran, all passed). restored: yes.\n2. Stdout-only replay: in `internal/responsebound/owner.go:215`, wrap the `replayed` append in `if w.stream == o.stdout`; same packages. Verdict: silent. restored: yes.\n\n## Optional advice\n- `commit` is a bounded verb, so retention reads HEAD after a commit; no test tells a read after the verb returns from a read before it (spec.md:148).\n- No test makes a call from outside any repository. The only aliases map to `help`.\n\n## Examined\nThe frozen diff 53e60dc5..9d989e12; spec stories 1–10, spec.md:143–170, the edge inventory and Won't handle lines; the ticket, coverage map, and consumer list; `owner.go` write, spill, and `Finish`; the registry dispatch and aliases; the benchkit hostile-input checklist; store states absent, symlinked scope, and existing file followed by a clear. No sign that the implementation command caused either gap.\n"
          },
          "axis": "Coverage",
          "base": "53e60dc58f0432fdadf2b31b78d4e7a25f85570d",
          "tip": "9d989e126b4cbeac507b12ea3137250fd9c6b294",
          "finding_ids": [
            "C1",
            "C2"
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

## Repair state

RT-RETAIN: repair cycle 1 of 2 is in progress for ticket 01. No reviewer extension exists.

## Standards

RT-RETAIN: 4 findings. The worst issue is S1.

- S1 (auto-fix, confidence 6): `internal/responsebound/retain.go:40` builds the response-store scope path a second time beside `location` in `internal/responsebound/store.go:46-51`.
- S2 (auto-fix, confidence 6): `internal/shift/shift.go:53` keeps its own porcelain status read beside `git.PorcelainStatus`. The orchestrator accepted it under the auto-approve direction, because `shift.go` is on the ticket's `Writes:` line. The shift call keeps its ignored error.
- S3 (auto-fix, confidence 4): `internal/responsebound/owner.go:78-81` and `owner.go:215-217` add the `whole` and `replayed` fields, which store again what the phase and the retained writes hold.
- S4 (auto-fix, confidence 5): the doc comment at `internal/responsebound/retain.go:389-390` repeats the word "call" ambiguously.

## Spec

RT-RETAIN: 0 actionable findings. P1 and P2 are no-op.

- P1 (no-op, confidence 7): RT7 and RT65 run at the `Owner.Retain` seam, not through the dispatcher. The behavior holds, and the deviation stays open to reviewer veto.
- P2 (no-op, confidence 5): a failed `HEAD` or porcelain read also clears the retained response.

## Coverage

RT-RETAIN: 2 findings. The worst issue is C1.

- C1 (auto-fix, confidence 8): no test runs two complete retaining calls in sequence, so a write-once mutant stays silent.
- C2 (auto-fix, confidence 7): no test interleaves stdout and stderr in a response within the bound, so a stdout-only mutant stays silent.
