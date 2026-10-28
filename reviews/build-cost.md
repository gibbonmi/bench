# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/build-cost/spec.md",
  "plan_digest": "sha256:09980ff61f50b93eee0fb166dc5464a9c75860fe31e90086b5f02f2ca01b8a0a",
  "implementation_session": "",
  "chunks": [
    {
      "id": "BC-C1",
      "base": "77ec3655a4601a57165b33d58c4328bdfdf72d34",
      "tip": "8224ff46132acdcdd36e4c86da54a2f49283c1d3",
      "plan_digest": "sha256:09980ff61f50b93eee0fb166dc5464a9c75860fe31e90086b5f02f2ca01b8a0a",
      "source_digest": "8558e1a9088a828e6abf0f6e8f00574a2ecaa55c",
      "acceptance_rows": [
        "BC24",
        "BC25",
        "BC26",
        "BC27",
        "BC28",
        "BC29",
        "BC30",
        "BC31",
        "BC32",
        "BC23",
        "BC33",
        "BC34",
        "BC35",
        "BC36",
        "BC37"
      ],
      "verification": [
        {
          "id": "t01-probe-01",
          "performer": "claude:bc-t01-author-20261028",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "09e8efc193f31a48ceba77ac89b1b2bf23665fa8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:bc-t01-author-20261028",
            "digest": "sha256:b872e89b309b3b13235935f0ef056788bfc8e3603c20742a479365a4ff4bd8ef",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/plan.go,swap,failed,3,TestPlanProbeOnlyChunks/absent_chunk_list,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestPlanProbeOnlyChunks,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,765\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,101,6\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/absent_chunk_list,\"plan_probe_test.go:64: plan grade refused the probe-only plan: chunk 2: missing required verification inventory\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/empty_chunk_list,\"plan_probe_test.go:64: plan grade refused the probe-only plan: chunk 2: missing required verification inventory\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/missing_probe_entry,\"plan_probe_test.go:74: chunk 2: missing required verification inventory\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f0d811335ac529969264109ac3439c54bc649f1c",
          "requirement": "probe-01",
          "command": "bench probe internal/reviewrecord/plan.go --swap 'if len(items) == 0 && owned == nil {' --with 'if len(items) == 0 {' --package ./internal/reviewrecord --run TestPlanProbeOnlyChunks",
          "exit_code": 0,
          "probe": {
            "mutation": "Refuse an empty chunk list again. The probe-only plan with an empty chunk list must then refuse, so the plan test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:bc-t01-author-20261028",
              "digest": "sha256:b872e89b309b3b13235935f0ef056788bfc8e3603c20742a479365a4ff4bd8ef",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/plan.go,swap,failed,3,TestPlanProbeOnlyChunks/absent_chunk_list,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestPlanProbeOnlyChunks,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,765\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,101,6\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/absent_chunk_list,\"plan_probe_test.go:64: plan grade refused the probe-only plan: chunk 2: missing required verification inventory\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/empty_chunk_list,\"plan_probe_test.go:64: plan grade refused the probe-only plan: chunk 2: missing required verification inventory\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/missing_probe_entry,\"plan_probe_test.go:74: chunk 2: missing required verification inventory\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t02-probe-02",
          "performer": "claude:bc-t02-author-20261028",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "1973ab354a72c1d6c865d7ecedda279f379d1d8a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:bc-t02-author-20261028",
            "digest": "sha256:698e50a3fa8e1c863d7ddafbfd6aa9aa024d8f24f20572abc18503d7abdc259e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/plan.go,swap,failed,1,TestPlanProbeOnlyChunks/plain_chunk_requirement,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestPlanProbeOnlyChunks,passed,7\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,392\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,158,7\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/plain_chunk_requirement,\"plan_probe_test.go:62: plan grade = <nil>, want a refusal naming \\\"chunk 2: verification plain has no probe\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "94c5e6fb21bcb29adc29b58701eca5da8c96eb01",
          "requirement": "probe-02",
          "command": "bench probe internal/reviewrecord/plan.go --swap 'if owned != nil && item.Probe == \"\" {' --with 'if false {' --package ./internal/reviewrecord --run TestPlanProbeOnlyChunks",
          "exit_code": 0,
          "probe": {
            "mutation": "Disable the plain chunk requirement refusal. The plan with a plain chunk requirement must then pass, so the plan test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:bc-t02-author-20261028",
              "digest": "sha256:698e50a3fa8e1c863d7ddafbfd6aa9aa024d8f24f20572abc18503d7abdc259e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/plan.go,swap,failed,1,TestPlanProbeOnlyChunks/plain_chunk_requirement,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestPlanProbeOnlyChunks,passed,7\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,392\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,158,7\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/plain_chunk_requirement,\"plan_probe_test.go:62: plan grade = <nil>, want a refusal naming \\\"chunk 2: verification plain has no probe\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "bc-c1-standards-1",
          "performer": "claude:bc-c1-standards-20261028",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8558e1a9088a828e6abf0f6e8f00574a2ecaa55c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:bc-c1-standards-20261028",
            "digest": "sha256:8ea31afe3c91ffcc5e6f579b85dd9f02f5e41478d34252196b113c242e3755ce",
            "excerpt": "axis: Standards\nchunk: BC-C1  pair: 77ec3655a..8224ff461\nfindings: 5\nworst: Removing check-floor leaves the command-package registry with no production reader, and a conformance grade still maintains it (F2).\n- F1 [auto-fix] [confidence 8] .agents/commands/bench-implement-spec.md:25 — \"The author runs the floor checks and witnesses of several packages in one `bench test` call\" still names the retired check floor; nothing defines \"floor checks\" now — reword to \"The author runs the witnesses of several packages in one `bench test` call that repeats `--package`\".\n- F2 [ask-user] [confidence 7] internal/tickets/registry_data.go:8-28 with internal/conformance/ticket_grammar_test.go:21-25,167-170 — CommandPackage has no production caller after check_floor.go left; a conformance grade keeps a dead registry current — the reviewer decides whether to retire commandPackages, CommandPackage, and gradeCommandPackages; internal/tickets/ is outside every signed-off Writes: directory.\n- F3 [auto-fix] [confidence 7] internal/testreport/selection_facts.go:74-78; internal/testreport/floor.go:11,16-22 — Request.Filtered has no caller, Floor.Checks and Floor.Embed are written and never read in production, and the Floor doc still says \"one ticket's check floor\" — delete them with their floor_test cases and reword the doc.\n- F4 [auto-fix] [confidence 6] internal/reviewrecord/recordcmd/stage_test.go:186 — the rewritten test keeps the legacy RC48 tag, whose row asserted the opposite subject — drop the RC48 tag from the edited comment.\n- F5 [auto-fix] [confidence 4] .agents/skills/bench-craft-spec/references/completion-plan.md:51 — one 117-column line in a file wrapped at about 80 columns — rewrap lines 47-52.\nreads: craft-review; AGENTS.md; BENCH.md; craft-comments; evidence metadata and consumers (touched=false rows walked); the delta; tickets 01-03; the BC-C1 plan; selection_facts.go; floor.go; fold.go:55-85; tickets/registry_data.go; ticket_grammar_test.go:160-235; rg sweeps.\nnot-verified: non-Go readers of commandPackages; projects/benchkit.md.\ncommand-contribution: yes, partly — when a ticket deletes production code, the author charge should run bench consumers on each symbol the deleted code called and report each that keeps no production reader (F2, F3).\n"
          },
          "axis": "Standards",
          "base": "77ec3655a4601a57165b33d58c4328bdfdf72d34",
          "tip": "8224ff46132acdcdd36e4c86da54a2f49283c1d3",
          "finding_ids": [
            "F1",
            "F3",
            "F4",
            "F5"
          ],
          "supersedes": []
        },
        {
          "id": "bc-c1-spec-1",
          "performer": "claude:bc-c1-spec-20261028",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8558e1a9088a828e6abf0f6e8f00574a2ecaa55c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:bc-c1-spec-20261028",
            "digest": "sha256:cd6bb8fd2608e01ad9c839c6c2552e8155bb45e2c2aba6d9e7ddb44e12960125",
            "excerpt": "axis: Spec\nchunk: BC-C1  pair: 77ec3655a..8224ff461\nfindings: 2\nworst: bench-implement-spec.md still tells the author to run \"the floor checks\", a term that no longer has a definition after the floor retired.\n- F1 [auto-fix] [confidence 6] .agents/commands/bench-implement-spec.md:25 — spec.md:173 retires the check-floor guidance, but line 25 still says \"The author runs the floor checks and witnesses of several packages in one `bench test` call\"; ticket 03's \"three other sentences\" missed it — change \"the floor checks\" to \"the development checks\".\n- F2 [no-op] [confidence 5] specs/build-cost/spec.md:187 — self-hosting step 5 says the orchestrator installs the worktree build, but the install applies only to a declared build and this spec declares no persisted field; the not-applicable call is correct — reword the sentence in a later plan commit.\nrows: BC23-BC30 met; BC31 review-owned met; BC32 met (placement accepted); BC33 met; BC34 review-owned met (bench-implement-spec.md:48); BC35 met; BC36 review-owned met (completion-plan.md:38); BC37 met.\nreads: craft-review; evidence metadata; the chunk diff; spec.md; tickets 01-03; bench coverage; per-commit spec word-diff; ADR 0012; BENCH-reference.md:398-409; reviewrecord/check.go:40-66; ownership/evidence.go:206-245; rg sweeps.\nnot-verified: no tests run; whether the charges ran the new binary through --in; the code was graded at 8224ff4.\ncommand-contribution: no command change is necessary\n"
          },
          "axis": "Spec",
          "base": "77ec3655a4601a57165b33d58c4328bdfdf72d34",
          "tip": "8224ff46132acdcdd36e4c86da54a2f49283c1d3",
          "finding_ids": [
            "F1",
            "F2"
          ],
          "supersedes": []
        },
        {
          "id": "bc-c1-coverage-1",
          "performer": "claude:bc-c1-coverage-20261028",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8558e1a9088a828e6abf0f6e8f00574a2ecaa55c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:bc-c1-coverage-20261028",
            "digest": "sha256:da5d188bf46e959ece73c5b6685f92587e803623e2b00c3eaefb5e838554d85b",
            "excerpt": "axis: Coverage\nchunk: BC-C1  pair: 77ec3655a..8224ff461\nfindings: 5\nworst: Ticket 01 split the empty-inventory refusal of requirementsValid by list. The final-list half is now the only half that no test checks.\n- F1 [auto-fix] [confidence 7] internal/reviewrecord/plan.go:180 — a plan with final_verification [] or absent; a wrong version that drops the owned == nil refusal or tests items == nil accepts it, and every test stays green — add \"empty final list\" and \"absent final list\" cases to TestPlanProbeOnlyChunks wanting \"final verification: missing required verification inventory\".\n- F2 [auto-fix] [confidence 7] internal/reviewrecord/recordcmd/ticket_probe_test.go:19 — the want row is built by hand from raw SHAs and copies the TOON quoting rule; a SHA in the quoted class turns BC29 red with no behavior change — build the want row through toon.Table as recordcmd/chunk_test.go:301 does.\n- F3 [auto-fix] [confidence 4] internal/reviewrecord/ticket.go:48 — a version 1 plan whose chunk lists no requirement is now legal; no test refuses the ticket form on it — add recordtest.Pooled(t, 1, recordtest.Owing(nil)) wanting \"the ticket form reads a version 2 plan\".\n- F4 [ask-user] [confidence 3] internal/reviewrecord/ticket.go:109 — no test refreezes a chunk for a ticket that owes nothing; low risk, freezeTip has no owed-count branch — reviewer decides on a refreeze case.\n- F5 [no-op] [confidence 5] internal/preflight/ownership/premises_test.go:273 — OC87 now asserts reviewrecord/check.go:60, not evidence.go:229; the unreachable-branch learning holds it; the fixture comment at premises_fixture_test.go:229 is now false (Standards item).\ntests-run: --package ./internal/reviewrecord --package ./internal/reviewrecord/recordcmd --run 'TestPlanProbeOnlyChunks|TestTicket' pass (7 + 44); --package ./internal/preflight --package ./internal/anchors --run 'TestDecideHasNoCheckFloorRow|TestCommandBuildRedsAPlainChunkRequirement|TestBuildCostAnchors' pass (2 + 1); two usage refusals ran no test.\nreads: evidence check-current; craft-review Coverage charge; BC23-BC37 rows; the code delta; plan.go, ticket.go, recordcmd/stage.go, recordtest/carry.go; evidence.go; benchkit.md checklist.\nnot-verified: no probe ran, so F1 and F3 rest on grep and code reading; reachability of evidence.go:229.\ncommand-contribution: no command change is necessary\n"
          },
          "axis": "Coverage",
          "base": "77ec3655a4601a57165b33d58c4328bdfdf72d34",
          "tip": "8224ff46132acdcdd36e4c86da54a2f49283c1d3",
          "finding_ids": [
            "F1",
            "F2",
            "F3",
            "F5"
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
  },
  "amendments": [
    {
      "from": "sha256:afed42110e5cf94f89620cf212044f2461429b735a54964f193ce5cadf70276d",
      "to": "sha256:efd81b765659c0c100c9c772db9e2519769687f6535d71d70ec3384b47c17814",
      "chunk_ids": {
        "BC-C1": [
          "BC-C1"
        ]
      }
    },
    {
      "from": "sha256:efd81b765659c0c100c9c772db9e2519769687f6535d71d70ec3384b47c17814",
      "to": "sha256:09980ff61f50b93eee0fb166dc5464a9c75860fe31e90086b5f02f2ca01b8a0a",
      "chunk_ids": {
        "BC-C1": [
          "BC-C1"
        ]
      }
    }
  ]
}
```

## BC-C1 review, round 1 (pair 77ec3655a..8224ff461)

Raw findings: Standards 5, Spec 2, Coverage 5. Repair targets after the merge of the duplicates: 9 on three tickets, 1 spec correction, 1 parked, 1 declined.

## Standards

Findings: 5. Worst: the command-package registry keeps no production reader after the check-floor row retired (F2).

- F1 auto-fix, confidence 8: `.agents/commands/bench-implement-spec.md:25` still names "the floor checks". Ticket 03 repair; one fix with Spec F1.
- F2 ask-user, confidence 7: `internal/tickets/registry_data.go:8-28` and `internal/conformance/ticket_grammar_test.go:21-25,167-170` keep a registry that no production code reads. Parked to `bench learning`, because the retirement removes a conformance grade outside every signed-off directory.
- F3 auto-fix, confidence 7: `internal/testreport/selection_facts.go:74-78` and `internal/testreport/floor.go:11,16-22` keep `Request.Filtered`, `Floor.Checks`, and `Floor.Embed` with no production reader. Ticket 01 repair.
- F4 auto-fix, confidence 6: `internal/reviewrecord/recordcmd/stage_test.go:186` keeps the RC48 tag on a test whose subject changed. Ticket 02 repair.
- F5 auto-fix, confidence 4: `.agents/skills/bench-craft-spec/references/completion-plan.md:51` holds one over-long line. Ticket 03 repair.

## Spec

Findings: 2. Worst: the build guidance still names the retired floor (F1).

- F1 auto-fix, confidence 6: `.agents/commands/bench-implement-spec.md:25`; the fix reads "the development checks". Ticket 03 repair.
- F2 no-op, confidence 5: `specs/build-cost/spec.md:187` says that the orchestrator installs the build. The orchestrator corrects the sentence in the repair plan commit.

## Coverage

Findings: 5. Worst: no test refuses an empty or absent final list (F1).

- F1 auto-fix, confidence 7: `internal/reviewrecord/plan.go:180`; add the empty and absent final-list cases. Ticket 01 repair.
- F2 auto-fix, confidence 7: `internal/reviewrecord/recordcmd/ticket_probe_test.go:19`; build the `want` row through `toon.Table`. Ticket 01 repair.
- F3 auto-fix, confidence 4: `internal/reviewrecord/ticket.go:48`; add a version 1 plan with an empty chunk that the ticket form refuses. Ticket 01 repair.
- F4 ask-user, confidence 3: a refreeze case for a ticket that owes nothing. Declined, because `freezeTip` and `changeChunk` have no branch on the owed count.
- F5 no-op, confidence 5: `internal/preflight/ownership/premises_fixture_test.go:229` holds a comment that is now false. Ticket 02 repair.
