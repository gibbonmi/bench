# Review outcomes

## GS1 initial review pickup

The initial review read the frozen pair `e4db0d2e..eee73331`. Each axis ran on sonnet at high effort, by reviewer direction. The axes found 16 findings. After the collapse of S2 into C1, 6 repair targets remain. Ticket 1 owns F1 and S3, and ticket 2 owns F2, C1, C2, and C5.

Each ticket gets its own fresh repair session on opus at medium effort, by reviewer direction. The chunk has used 0 of its 2 repair cycles, and the next repair is cycle 1.

The reviewer deferred each ask-user finding to two consultants, opus at medium effort and sonnet at medium effort. A fable delegate at medium effort decided the one stalemate. The orchestrator committed the resulting spec corrections in one plan commit before the repair.

## GS1 repair cycle 1

Repair cycle 1 closed all 6 repair targets, so the chunk has used 1 of its 2 repair cycles. The ticket 1 repair `0cf5c883` closed F1 and S3, and the ticket 2 repair `cdeaf570` closed F2, C1, C2, and C5. Both repairs change test files only. The chunk is refrozen at `e4db0d2e..cdeaf570`, and both repair sessions recorded their verifications at that tip. The orchestrator's swap probe of the nested walk bit. A confirming round of all three axes reads the repair delta, because the plan bytes in the delta block a coordinator confirmation.

## GS1 confirming round 1

The confirming round read the repair delta `eee73331..cdeaf570` on sonnet at high effort. It closed F1, F2, S1, S2, S3, C1, C2, and C5. Standards and Spec passed, and Coverage found one blocker.

- N1, auto-fix, confidence 8. The `WalkDir` error branch of the ownership check has no witness, and a swap of its diagnostic read silent. Both consultants chose a repair. A silent mutant of a branch that the repair added fails the witness rule of `/bench-implement-spec`. Ticket 2 owns it, and repair cycle 2 adds one witness for an unreadable nested directory.
- Spec N1, no-op, confidence 5. The fable tie-break of the initial round decided that "real import bindings" bans text matching, so `spec.md` line 135 stays.
- Spec N2 and Coverage N2-N4, no-op. They name hardening that the spec permits, a symlinked directory that the repair did not widen, and pre-repair behavior.

## GS1 repair cycle 2

Repair cycle 2 closed N1, so the chunk has used 2 of its 2 repair cycles. The ticket 2 repair `2461ab1f` adds the witness `TestGuardScanLoopOwnerReportsAnUnreadableSubdir`, and the repair changes a test file only. The chunk is refrozen at `e4db0d2e..2461ab1f`, and both tickets recorded their verifications at that tip. The orchestrator's insertion probe, which drops the walk error before the switch, bit. A confirming round of all three axes reads the repair delta `cdeaf570..2461ab1f`.

## Standards

The axis found 4 findings. The worst issue is F2, because the test respells the diagnostic text that the check spells.

- F1, auto-fix, confidence 6. The `familySGolden` comment in `internal/guards/scan_loop_test.go` narrates its provenance and argues its own independence. Keep only the fact that it is the reference capture of scan family S.
- F2, auto-fix, confidence 5. `internal/conformance/guard_scan_owner_test.go` respells the suffixes " is missing or unreadable" and "cannot parse " that the check spells. Share each piece as one constant.
- F3, no-op, confidence 4. The five construct names stay an independent expectation, because a list derived from the check would stay green when a case is dropped. The ticket 2 author recorded a natural red before the production edit: `TestGuardScanLoopOwnerBitesOnPrivateCopy` failed one `t.Errorf` for each of the five construct kinds. Both consultants chose this disposition.
- F4, no-op, confidence 6. The GS05 and GS28 read-count witnesses are necessary, and the select ban limits their risk.

## Spec

The axis found 5 findings. The worst issue is S1, because the Edge inventory listed a family S member that the golden omits.

- S1, no-op, confidence 6. No Command cancel event lands between two candidates, so GS05 owns that member through a direct Scan. Both consultants chose a spec correction, and the Edge inventory now states this split.
- S2, auto-fix, confidence 5. This finding names the same defect as C1, and the C1 repair closes it.
- S3, auto-fix, confidence 4. The inspection-cancelled family S members depend on a 10ms wall-clock window. Make the cancel moment independent of a short timer, and keep the golden bytes unchanged.
- S4, no-op, confidence 3. GS11 observes its contract.
- S5, no-op, confidence 3. GS23 closes the second-path failure of GS22.

## Coverage

The axis found 8 findings. The worst issue is C1, because a pipeline in a subdirectory of `internal/guards` escapes the check, and the spec says "beneath".

- C1, auto-fix, confidence 7. Walk every production Go file beneath `internal/guards`, skip `testdata`, and add a nested-directory witness.
- C2, auto-fix, confidence 6. Skip or name a special file whose name ends in `.go` before the parser opens it, as the hostile-input checklist of `projects/benchkit.md` requires. Add a witness.
- C3, no-op, confidence 5. Both consultants chose a spec correction. The check grades `make` with a literal channel type, and a Won't handle line names the named channel type.
- C4, no-op, confidence 5. The consultants split. The fable tie-break chose a Won't handle line for goroutine starts with no go statement and no check change. The phrase "real import bindings" bans text matching, so the check meets it.
- C5, auto-fix, confidence 5. Only `fs.ErrNotExist` takes the optional-surface exit. Another `os.Stat` error gets a named diagnostic. Add a witness.
- C6, C7, and C8, no-op, confidence 6. GS23 covers the GS22 second path, and the preservation rows are legitimate regression guards. The axis demonstrated the GS21 invalid-arguments red with a probe at `internal/guards/guards.go`.

```bench-review-record
{
  "version": 2,
  "spec": "specs/guard-scan-loop/spec.md",
  "plan_digest": "sha256:552f0b3d19b40a0481174b21ae276df18c64cfbc7270fc30f94b0b66660045cd",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GS1",
      "base": "e4db0d2e66ffae4045c9c29c11c4696462442bc4",
      "tip": "2461ab1f5d246ddd59307d74cecc45bb55dbdb3e",
      "plan_digest": "sha256:552f0b3d19b40a0481174b21ae276df18c64cfbc7270fc30f94b0b66660045cd",
      "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
      "acceptance_rows": [
        "GS01",
        "GS02",
        "GS03",
        "GS04",
        "GS05",
        "GS06",
        "GS07",
        "GS08",
        "GS09",
        "GS10",
        "GS11",
        "GS12",
        "GS13",
        "GS14",
        "GS15",
        "GS16",
        "GS17",
        "GS18",
        "GS19",
        "GS20",
        "GS21",
        "GS22",
        "GS27",
        "GS28",
        "GS29",
        "GS30",
        "GS23",
        "GS24",
        "GS25",
        "GS26"
      ],
      "verification": [
        {
          "id": "t2-t2-registered-owner",
          "performer": "claude:gsl-t2-author-medium-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-author-medium-20261020",
            "digest": "sha256:34a4930feacf842a35b593436ffc61282fe1fdfa2ac7aa9dfa21744e72407599",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,2588,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-registered-owner",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t2-t2-conformance-closure",
          "performer": "claude:gsl-t2-author-medium-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-author-medium-20261020",
            "digest": "sha256:8ac87902ecbc908a031ddf854e5865f2e9518c94fe7fb6d71be796a5a198e60d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,35750,1123\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t2-t2-root-conformance",
          "performer": "claude:gsl-t2-author-medium-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-author-medium-20261020",
            "digest": "sha256:c683a0ffa1cc00549b052f5bf3f515f5a4d48145e7e54f92b8f3c56ca9e656d9",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13883,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-t2-copy-check-omission",
          "performer": "claude:gsl-t2-author-medium-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-author-medium-20261020",
            "digest": "sha256:c3aa62b8498fafacdcba6c335e96c9080d0e807b1a0780fd2bdb2b717db987aa",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGuardScanLoopOwnerBitesOnPrivateCopy,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/conformance,^TestGuardScanLoopOwnerBitesOnPrivateCopy$,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,400\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuardScanLoopOwnerBitesOnPrivateCopy,\"guard_scan_owner_test.go:155: want a go statement diagnostic on the private pipeline copy, got []\",5\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-copy-check-omission",
          "command": "bench test --package ./internal/conformance --run '^TestGuardScanLoopOwnerBitesOnPrivateCopy$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit only the registered checkGuardScanLoopOwner call while retaining its function and tests; an unused old worker pipeline must stop producing the expected registered diagnostic and red this witness.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gsl-t2-author-medium-20261020",
              "digest": "sha256:c3aa62b8498fafacdcba6c335e96c9080d0e807b1a0780fd2bdb2b717db987aa",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGuardScanLoopOwnerBitesOnPrivateCopy,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/conformance,^TestGuardScanLoopOwnerBitesOnPrivateCopy$,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,400\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuardScanLoopOwnerBitesOnPrivateCopy,\"guard_scan_owner_test.go:155: want a go statement diagnostic on the private pipeline copy, got []\",5\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-t1-scan-and-command",
          "performer": "claude:gsl-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
            "digest": "sha256:e192d277a2cff5d91e7bdbafad8753907bac9e42dfb9833a04aa282fd4b23c61",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,pass,397,54\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-scan-and-command",
          "command": "bench test --package ./internal/guards",
          "exit_code": 0
        },
        {
          "id": "t1-t1-bounds-policy",
          "performer": "claude:gsl-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
            "digest": "sha256:2845b2438ef099916edff839bb2460d39867f89b12bf8bc45a5555774a56ca99",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  bounds-policy,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,235,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-bounds-policy",
          "command": "bench test --check bounds-policy",
          "exit_code": 0
        },
        {
          "id": "t1-t1-conformance-closure",
          "performer": "claude:gsl-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
            "digest": "sha256:599dd883b29246a5a68725f92a1e91a076a4831b98ff30c23626eeaa71c35be5",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,36384,1123\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-t1-root-conformance",
          "performer": "claude:gsl-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
            "digest": "sha256:4259ae1a40e1fb5f82ba6f20496a07a1954310daa14e504827523bcb809a28db",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13663,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-t1-command-closure",
          "performer": "claude:gsl-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
            "digest": "sha256:8be28703f308532d7438d2ad13a3e8e23896912bcdd33021b6b9aa6823eebba1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15223,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-command-closure",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1-t1-enumeration-return-omission",
          "performer": "claude:gsl-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
            "digest": "sha256:27cabe066002fefad126b797e9fda0d5ed74748bd5d51b5ddd0b795f68f78691",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopEnumerationReturnsAfterCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopEnumerationReturnsAfterCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,489\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopEnumerationReturnsAfterCancellation,\"scan_loop_test.go:335: Scan = guards.ScanResult{Rows:[][]string(nil), Status:\\\"incomplete\\\", Inspected:\\\"0\\\", Total:\\\"2\\\", Omitted:\\\"2\\\", Reason:\\\"timeout\\\"} after operations [\\\"enumerate\\\"]\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-enumeration-return-omission",
          "command": "bench test --package ./internal/guards --run '^TestScanLoopEnumerationReturnsAfterCancellation$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Gate the post-enumeration checkpoint on a non-nil enumeration error only; a nil-error candidate slice returned after cancellation must report known counts and red this witness.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
              "digest": "sha256:27cabe066002fefad126b797e9fda0d5ed74748bd5d51b5ddd0b795f68f78691",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopEnumerationReturnsAfterCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopEnumerationReturnsAfterCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,489\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopEnumerationReturnsAfterCancellation,\"scan_loop_test.go:335: Scan = guards.ScanResult{Rows:[][]string(nil), Status:\\\"incomplete\\\", Inspected:\\\"0\\\", Total:\\\"2\\\", Omitted:\\\"2\\\", Reason:\\\"timeout\\\"} after operations [\\\"enumerate\\\"]\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-t1-differential-golden-bite",
          "performer": "claude:gsl-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
            "digest": "sha256:67e1a85a72ad1628ed4ba08054fdeda21e238df1f9ba871b98b5a1494c0dc4c7",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,8,TestScanLoopDifferentialCommandFamilies/complete,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopDifferentialCommandFamilies$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1833\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,302,19\nfailures[8]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/complete,\"scan_loop_test.go:263: member complete differs from testdata/scan-family-s.golden\",43\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/fifo-header,\"scan_loop_test.go:263: member fifo-header differs from testdata/scan-family-s.golden\",37\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-empty,\"scan_loop_test.go:263: member inspection-cancelled-empty differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-rows,\"scan_loop_test.go:263: member inspection-cancelled-rows differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/static-headers,\"scan_loop_test.go:263: member static-headers differs from testdata/scan-family-s.golden\",49\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/symlink-and-missing-headers,\"scan_loop_test.go:263: member symlink-and-missing-headers differs from testdata/scan-family-s.golden\",47\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-claude-and-codex,\"scan_loop_test.go:263: member wiring-claude-and-codex differs from testdata/scan-family-s.golden\",41\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-empty-and-invalid,\"scan_loop_test.go:263: member wiring-empty-and-invalid differs from testdata/scan-family-s.golden\",37\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-differential-golden-bite",
          "command": "bench test --package ./internal/guards --run '^TestScanLoopDifferentialCommandFamilies$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Reverse the candidate inspection order in the candidate loop; the independently captured scan-family-s.golden expectation must red this witness on the complete member.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
              "digest": "sha256:67e1a85a72ad1628ed4ba08054fdeda21e238df1f9ba871b98b5a1494c0dc4c7",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,8,TestScanLoopDifferentialCommandFamilies/complete,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopDifferentialCommandFamilies$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1833\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,302,19\nfailures[8]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/complete,\"scan_loop_test.go:263: member complete differs from testdata/scan-family-s.golden\",43\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/fifo-header,\"scan_loop_test.go:263: member fifo-header differs from testdata/scan-family-s.golden\",37\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-empty,\"scan_loop_test.go:263: member inspection-cancelled-empty differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-rows,\"scan_loop_test.go:263: member inspection-cancelled-rows differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/static-headers,\"scan_loop_test.go:263: member static-headers differs from testdata/scan-family-s.golden\",49\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/symlink-and-missing-headers,\"scan_loop_test.go:263: member symlink-and-missing-headers differs from testdata/scan-family-s.golden\",47\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-claude-and-codex,\"scan_loop_test.go:263: member wiring-claude-and-codex differs from testdata/scan-family-s.golden\",41\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-empty-and-invalid,\"scan_loop_test.go:263: member wiring-empty-and-invalid differs from testdata/scan-family-s.golden\",37\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-t1-post-inspection-omission",
          "performer": "claude:gsl-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
            "digest": "sha256:39519a6de7f2e254f566175e4285294c2020d7454b7e8a5dd8c78858f492d852",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopPostInspectionCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopPostInspectionCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,575\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,4,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopPostInspectionCancellation,\"scan_loop_test.go:352: Scan = guards.ScanResult{Rows:[][]string{[]string{\\\"a\\\", \\\"boundary of a\\\", \\\"denies of a\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}, []string{\\\"b\\\", \\\"boundary of b\\\", \\\"denies of b\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}}, Status:\\\"complete\\\", Insp… (344 bytes)\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-post-inspection-omission",
          "command": "bench test --package ./internal/guards --run '^TestScanLoopPostInspectionCancellation$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit only the post-inspection context checkpoint; a returned interrupted candidate must publish rows and red the completed-prefix contract.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gsl-t1-author-20261020",
              "digest": "sha256:39519a6de7f2e254f566175e4285294c2020d7454b7e8a5dd8c78858f492d852",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopPostInspectionCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopPostInspectionCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,575\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,4,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopPostInspectionCancellation,\"scan_loop_test.go:352: Scan = guards.ScanResult{Rows:[][]string{[]string{\\\"a\\\", \\\"boundary of a\\\", \\\"denies of a\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}, []string{\\\"b\\\", \\\"boundary of b\\\", \\\"denies of b\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}}, Status:\\\"complete\\\", Insp… (344 bytes)\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2r1-t2-registered-owner",
          "performer": "claude:gsl-t2-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-repair1-20261021",
            "digest": "sha256:ea83a744fff4141b6a3eaaf205dcddff7b652c091eeebd13f9ad1b8b27856d83",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,2562,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-registered-owner",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-conformance-closure",
          "performer": "claude:gsl-t2-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-repair1-20261021",
            "digest": "sha256:3f950cf3e844ba2af18692171d158a2670c2e1e0f2d886b3c0a1304d5a299ed7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,35820,1126\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-root-conformance",
          "performer": "claude:gsl-t2-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-repair1-20261021",
            "digest": "sha256:c683a0ffa1cc00549b052f5bf3f515f5a4d48145e7e54f92b8f3c56ca9e656d9",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13883,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-copy-check-omission",
          "performer": "claude:gsl-t2-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-repair1-20261021",
            "digest": "sha256:3f226d7e5510b67dbb40090d2ec3ce9a53b6592d9e9668f5f04a28c4080ed0b4",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGuardScanLoopOwnerBitesOnPrivateCopy,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/conformance,^TestGuardScanLoopOwnerBitesOnPrivateCopy$,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,400\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuardScanLoopOwnerBitesOnPrivateCopy,\"guard_scan_owner_test.go:182: want a go statement diagnostic on the private pipeline copy, got []\",5\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-copy-check-omission",
          "command": "bench test --package ./internal/conformance --run '^TestGuardScanLoopOwnerBitesOnPrivateCopy$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit only the registered checkGuardScanLoopOwner call while retaining its function and tests; an unused old worker pipeline must stop producing the expected registered diagnostic and red this witness.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gsl-t2-repair1-20261021",
              "digest": "sha256:3f226d7e5510b67dbb40090d2ec3ce9a53b6592d9e9668f5f04a28c4080ed0b4",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGuardScanLoopOwnerBitesOnPrivateCopy,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/conformance,^TestGuardScanLoopOwnerBitesOnPrivateCopy$,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,400\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuardScanLoopOwnerBitesOnPrivateCopy,\"guard_scan_owner_test.go:182: want a go statement diagnostic on the private pipeline copy, got []\",5\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r1-t1-scan-and-command",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:03e346c42bd1fca8781a6c5027762641e854e1c6a383c1387ab9d4f1701c5730",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,pass,6373,54\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-scan-and-command",
          "command": "bench test --package ./internal/guards",
          "exit_code": 0
        },
        {
          "id": "t1r1-t1-bounds-policy",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:6c4d187786206f737507d0d7533acf313a933fea1ca854e9ba27ed4880117098",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  bounds-policy,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,230,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-bounds-policy",
          "command": "bench test --check bounds-policy",
          "exit_code": 0
        },
        {
          "id": "t1r1-t1-conformance-closure",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:0892a661d2d377c48cf414a6d3a1561954cdce436f1aa12ddc4bf9fc9441bb66",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,36006,1126\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1r1-t1-root-conformance",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:93c858bb6f7f5b46b7685ff7fae236ebfdc4cc0bbdc22fa36e88e6d3292cb9b8",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13606,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1r1-t1-command-closure",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:00560787cd12f48454a92e2e37187b18b2456fd3a7420d9c1c429f855ba268d9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14728,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-command-closure",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1r1-t1-enumeration-return-omission",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:27cabe066002fefad126b797e9fda0d5ed74748bd5d51b5ddd0b795f68f78691",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopEnumerationReturnsAfterCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopEnumerationReturnsAfterCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,489\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopEnumerationReturnsAfterCancellation,\"scan_loop_test.go:335: Scan = guards.ScanResult{Rows:[][]string(nil), Status:\\\"incomplete\\\", Inspected:\\\"0\\\", Total:\\\"2\\\", Omitted:\\\"2\\\", Reason:\\\"timeout\\\"} after operations [\\\"enumerate\\\"]\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-enumeration-return-omission",
          "command": "bench test --package ./internal/guards --run '^TestScanLoopEnumerationReturnsAfterCancellation$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Gate the post-enumeration checkpoint on a non-nil enumeration error only; a nil-error candidate slice returned after cancellation must report known counts and red this witness.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
              "digest": "sha256:27cabe066002fefad126b797e9fda0d5ed74748bd5d51b5ddd0b795f68f78691",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopEnumerationReturnsAfterCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopEnumerationReturnsAfterCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,489\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopEnumerationReturnsAfterCancellation,\"scan_loop_test.go:335: Scan = guards.ScanResult{Rows:[][]string(nil), Status:\\\"incomplete\\\", Inspected:\\\"0\\\", Total:\\\"2\\\", Omitted:\\\"2\\\", Reason:\\\"timeout\\\"} after operations [\\\"enumerate\\\"]\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r1-t1-differential-golden-bite",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:93343ca645108092fb18f65fed2c9a829d0fb56d6bf1ec56ba17a131379be6cd",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,8,TestScanLoopDifferentialCommandFamilies/complete,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopDifferentialCommandFamilies$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1834\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,6239,19\nfailures[8]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/complete,\"scan_loop_test.go:263: member complete differs from testdata/scan-family-s.golden\",43\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/fifo-header,\"scan_loop_test.go:263: member fifo-header differs from testdata/scan-family-s.golden\",37\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-empty,\"scan_loop_test.go:263: member inspection-cancelled-empty differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-rows,\"scan_loop_test.go:263: member inspection-cancelled-rows differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/static-headers,\"scan_loop_test.go:263: member static-headers differs from testdata/scan-family-s.golden\",49\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/symlink-and-missing-headers,\"scan_loop_test.go:263: member symlink-and-missing-headers differs from testdata/scan-family-s.golden\",47\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-claude-and-codex,\"scan_loop_test.go:263: member wiring-claude-and-codex differs from testdata/scan-family-s.golden\",41\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-empty-and-invalid,\"scan_loop_test.go:263: member wiring-empty-and-invalid differs from testdata/scan-family-s.golden\",37\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-differential-golden-bite",
          "command": "bench test --package ./internal/guards --run '^TestScanLoopDifferentialCommandFamilies$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Reverse the candidate inspection order in the candidate loop; the independently captured scan-family-s.golden expectation must red this witness on the complete member.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
              "digest": "sha256:93343ca645108092fb18f65fed2c9a829d0fb56d6bf1ec56ba17a131379be6cd",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,8,TestScanLoopDifferentialCommandFamilies/complete,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopDifferentialCommandFamilies$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1834\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,6239,19\nfailures[8]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/complete,\"scan_loop_test.go:263: member complete differs from testdata/scan-family-s.golden\",43\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/fifo-header,\"scan_loop_test.go:263: member fifo-header differs from testdata/scan-family-s.golden\",37\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-empty,\"scan_loop_test.go:263: member inspection-cancelled-empty differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-rows,\"scan_loop_test.go:263: member inspection-cancelled-rows differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/static-headers,\"scan_loop_test.go:263: member static-headers differs from testdata/scan-family-s.golden\",49\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/symlink-and-missing-headers,\"scan_loop_test.go:263: member symlink-and-missing-headers differs from testdata/scan-family-s.golden\",47\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-claude-and-codex,\"scan_loop_test.go:263: member wiring-claude-and-codex differs from testdata/scan-family-s.golden\",41\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-empty-and-invalid,\"scan_loop_test.go:263: member wiring-empty-and-invalid differs from testdata/scan-family-s.golden\",37\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r1-t1-post-inspection-omission",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:39519a6de7f2e254f566175e4285294c2020d7454b7e8a5dd8c78858f492d852",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopPostInspectionCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopPostInspectionCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,575\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,4,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopPostInspectionCancellation,\"scan_loop_test.go:352: Scan = guards.ScanResult{Rows:[][]string{[]string{\\\"a\\\", \\\"boundary of a\\\", \\\"denies of a\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}, []string{\\\"b\\\", \\\"boundary of b\\\", \\\"denies of b\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}}, Status:\\\"complete\\\", Insp… (344 bytes)\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-post-inspection-omission",
          "command": "bench test --package ./internal/guards --run '^TestScanLoopPostInspectionCancellation$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit only the post-inspection context checkpoint; a returned interrupted candidate must publish rows and red the completed-prefix contract.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
              "digest": "sha256:39519a6de7f2e254f566175e4285294c2020d7454b7e8a5dd8c78858f492d852",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopPostInspectionCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopPostInspectionCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,575\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,4,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopPostInspectionCancellation,\"scan_loop_test.go:352: Scan = guards.ScanResult{Rows:[][]string{[]string{\\\"a\\\", \\\"boundary of a\\\", \\\"denies of a\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}, []string{\\\"b\\\", \\\"boundary of b\\\", \\\"denies of b\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}}, Status:\\\"complete\\\", Insp… (344 bytes)\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2r2-t2-registered-owner",
          "performer": "claude:gsl-t2-repair2-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-repair2-20261021",
            "digest": "sha256:11486b9fd4572df5d059e2a9c84fd5c0604eaaabcefcedf32be86e5da34ab94e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,2567,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-registered-owner",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t2r2-t2-conformance-closure",
          "performer": "claude:gsl-t2-repair2-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-repair2-20261021",
            "digest": "sha256:ccdb457c5e2638eb980d098586797e6dbc3d662f0eb949b352c50607ca92ff51",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,35503,1127\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t2r2-t2-root-conformance",
          "performer": "claude:gsl-t2-repair2-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-repair2-20261021",
            "digest": "sha256:025ae3c01947cdf49ccb8776e3f5fdc6c746aae74c4d6f4ad67192cda694ddd2",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13269,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2r2-t2-copy-check-omission",
          "performer": "claude:gsl-t2-repair2-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t2-repair2-20261021",
            "digest": "sha256:3f226d7e5510b67dbb40090d2ec3ce9a53b6592d9e9668f5f04a28c4080ed0b4",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGuardScanLoopOwnerBitesOnPrivateCopy,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/conformance,^TestGuardScanLoopOwnerBitesOnPrivateCopy$,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,400\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuardScanLoopOwnerBitesOnPrivateCopy,\"guard_scan_owner_test.go:182: want a go statement diagnostic on the private pipeline copy, got []\",5\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-copy-check-omission",
          "command": "bench test --package ./internal/conformance --run '^TestGuardScanLoopOwnerBitesOnPrivateCopy$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit only the registered checkGuardScanLoopOwner call while retaining its function and tests; an unused old worker pipeline must stop producing the expected registered diagnostic and red this witness.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gsl-t2-repair2-20261021",
              "digest": "sha256:3f226d7e5510b67dbb40090d2ec3ce9a53b6592d9e9668f5f04a28c4080ed0b4",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGuardScanLoopOwnerBitesOnPrivateCopy,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/conformance,^TestGuardScanLoopOwnerBitesOnPrivateCopy$,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,400\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuardScanLoopOwnerBitesOnPrivateCopy,\"guard_scan_owner_test.go:182: want a go statement diagnostic on the private pipeline copy, got []\",5\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r2-t1-scan-and-command",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:b23fda31b41a1bb13c96b55d6272e92f07ee10f001695f4a797fc53e1b33612c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,pass,6353,54\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-scan-and-command",
          "command": "bench test --package ./internal/guards",
          "exit_code": 0
        },
        {
          "id": "t1r2-t1-bounds-policy",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:858bc235c25cf5a7055131d9c9924765c9d039285fb3a28843bf1c0b1ed8b71a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  bounds-policy,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,229,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-bounds-policy",
          "command": "bench test --check bounds-policy",
          "exit_code": 0
        },
        {
          "id": "t1r2-t1-conformance-closure",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:5e2fb05c71c1fd728ee99f1a07aa74e7830d773d7efc9c9af2915d087489beec",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,35593,1127\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-conformance-closure",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1r2-t1-root-conformance",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:7435777686970920ec1e0f80feddf180c668c4cc281bb2485b49342884f0aa6d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13438,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1r2-t1-command-closure",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:66d6bbc2aa2b5d1862cd2e4100e13cfb6e61da98a62cfbdbfa194bc5e670ed95",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14787,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-command-closure",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1r2-t1-enumeration-return-omission",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:27cabe066002fefad126b797e9fda0d5ed74748bd5d51b5ddd0b795f68f78691",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopEnumerationReturnsAfterCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopEnumerationReturnsAfterCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,489\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopEnumerationReturnsAfterCancellation,\"scan_loop_test.go:335: Scan = guards.ScanResult{Rows:[][]string(nil), Status:\\\"incomplete\\\", Inspected:\\\"0\\\", Total:\\\"2\\\", Omitted:\\\"2\\\", Reason:\\\"timeout\\\"} after operations [\\\"enumerate\\\"]\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-enumeration-return-omission",
          "command": "bench test --package ./internal/guards --run '^TestScanLoopEnumerationReturnsAfterCancellation$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Gate the post-enumeration checkpoint on a non-nil enumeration error only; a nil-error candidate slice returned after cancellation must report known counts and red this witness.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
              "digest": "sha256:27cabe066002fefad126b797e9fda0d5ed74748bd5d51b5ddd0b795f68f78691",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopEnumerationReturnsAfterCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopEnumerationReturnsAfterCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,489\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopEnumerationReturnsAfterCancellation,\"scan_loop_test.go:335: Scan = guards.ScanResult{Rows:[][]string(nil), Status:\\\"incomplete\\\", Inspected:\\\"0\\\", Total:\\\"2\\\", Omitted:\\\"2\\\", Reason:\\\"timeout\\\"} after operations [\\\"enumerate\\\"]\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r2-t1-differential-golden-bite",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:5840b00f59052a1f4d03f98b35940a7b10df774e374f23a3bbf64c69082b74c4",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,8,TestScanLoopDifferentialCommandFamilies/complete,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopDifferentialCommandFamilies$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1834\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,6262,19\nfailures[8]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/complete,\"scan_loop_test.go:263: member complete differs from testdata/scan-family-s.golden\",43\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/fifo-header,\"scan_loop_test.go:263: member fifo-header differs from testdata/scan-family-s.golden\",37\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-empty,\"scan_loop_test.go:263: member inspection-cancelled-empty differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-rows,\"scan_loop_test.go:263: member inspection-cancelled-rows differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/static-headers,\"scan_loop_test.go:263: member static-headers differs from testdata/scan-family-s.golden\",49\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/symlink-and-missing-headers,\"scan_loop_test.go:263: member symlink-and-missing-headers differs from testdata/scan-family-s.golden\",47\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-claude-and-codex,\"scan_loop_test.go:263: member wiring-claude-and-codex differs from testdata/scan-family-s.golden\",41\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-empty-and-invalid,\"scan_loop_test.go:263: member wiring-empty-and-invalid differs from testdata/scan-family-s.golden\",37\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-differential-golden-bite",
          "command": "bench test --package ./internal/guards --run '^TestScanLoopDifferentialCommandFamilies$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Reverse the candidate inspection order in the candidate loop; the independently captured scan-family-s.golden expectation must red this witness on the complete member.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
              "digest": "sha256:5840b00f59052a1f4d03f98b35940a7b10df774e374f23a3bbf64c69082b74c4",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,8,TestScanLoopDifferentialCommandFamilies/complete,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopDifferentialCommandFamilies$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1834\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,6262,19\nfailures[8]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/complete,\"scan_loop_test.go:263: member complete differs from testdata/scan-family-s.golden\",43\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/fifo-header,\"scan_loop_test.go:263: member fifo-header differs from testdata/scan-family-s.golden\",37\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-empty,\"scan_loop_test.go:263: member inspection-cancelled-empty differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/inspection-cancelled-rows,\"scan_loop_test.go:263: member inspection-cancelled-rows differs from testdata/scan-family-s.golden\",29\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/static-headers,\"scan_loop_test.go:263: member static-headers differs from testdata/scan-family-s.golden\",49\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/symlink-and-missing-headers,\"scan_loop_test.go:263: member symlink-and-missing-headers differs from testdata/scan-family-s.golden\",47\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-claude-and-codex,\"scan_loop_test.go:263: member wiring-claude-and-codex differs from testdata/scan-family-s.golden\",41\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopDifferentialCommandFamilies/wiring-empty-and-invalid,\"scan_loop_test.go:263: member wiring-empty-and-invalid differs from testdata/scan-family-s.golden\",37\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r2-t1-post-inspection-omission",
          "performer": "claude:gsl-t1-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "364c4867dce31d62207fb8a1be27a9ac4774b2fb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
            "digest": "sha256:300bdbf931805075606e97f7e8adcb7f8d8f21377766c1c2acc8002781594119",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopPostInspectionCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopPostInspectionCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,575\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopPostInspectionCancellation,\"scan_loop_test.go:352: Scan = guards.ScanResult{Rows:[][]string{[]string{\\\"a\\\", \\\"boundary of a\\\", \\\"denies of a\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}, []string{\\\"b\\\", \\\"boundary of b\\\", \\\"denies of b\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}}, Status:\\\"complete\\\", Insp… (344 bytes)\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-post-inspection-omission",
          "command": "bench test --package ./internal/guards --run '^TestScanLoopPostInspectionCancellation$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit only the post-inspection context checkpoint; a returned interrupted candidate must publish rows and red the completed-prefix contract.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gsl-t1-repair1-20261021",
              "digest": "sha256:300bdbf931805075606e97f7e8adcb7f8d8f21377766c1c2acc8002781594119",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestScanLoopPostInspectionCancellation,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/guards,^TestScanLoopPostInspectionCancellation$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,575\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanLoopPostInspectionCancellation,\"scan_loop_test.go:352: Scan = guards.ScanResult{Rows:[][]string{[]string{\\\"a\\\", \\\"boundary of a\\\", \\\"denies of a\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}, []string{\\\"b\\\", \\\"boundary of b\\\", \\\"denies of b\\\", \\\"\\\", \\\"\\\", \\\"current\\\", \\\"claude\\\"}}, Status:\\\"complete\\\", Insp… (344 bytes)\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gs1-r1-standards",
          "performer": "claude:gsl-review-standards-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gsl-review-standards-r1",
            "digest": "sha256:02983951314c50d6d8732be6e296eb1b7f61d05e9065788b6be8ff2cad32f9e4",
            "excerpt": "STANDARDS axis, GS1 (base e4db0d2e6, code tip eee733317). check-current: current=true.\nF1 | auto-fix | internal/guards/scan_loop_test.go:24-26 | conf 6. The familySGolden comment narrates provenance and argues its own independence; craft-comments forbids narration. Keep only \"reference capture of scan family S, reviewed apart from the candidate\".\nF2 | auto-fix | internal/conformance/guard_scan_owner_test.go:180, :191 vs :47, :54 | conf 5. The test respells the diagnostic suffixes \" is missing or unreadable\" and \"cannot parse \" that the check spells; only the prefix is shared. AGENTS.md one source per fact. Share the two pieces as constants.\nF3 | ask-user | internal/conformance/guard_scan_owner_test.go:153 vs :69-87 | conf 4. The test lists the five construct names that scanPipelineConstruct returns, with no recorded per-construct red. Record a mutation, or derive the list from the check.\nF4 | no-op | internal/guards/scan_loop_test.go:340, :372-374 | conf 6. GS05 and GS28 read-count witnesses are necessary and not brittle beyond the spec; the select ban limits the risk.\nCount 4 (2 auto-fix, 1 ask-user, 1 no-op). Worst: F2. Structure limits all met; the move is verbatim; no edit outside the acceptance's reach; 6 touched=false consumers read, none stale.\nAdvice: a shared production-.go listing helper does not exist; make(namedChanType) is not detected; the in-process fifo-header member could hang on a regression; magic numbers 5 and 7 could derive from len(names).\n"
          },
          "axis": "Standards",
          "base": "e4db0d2e66ffae4045c9c29c11c4696462442bc4",
          "tip": "eee733317545f94b7df5efcc56f23e1e93437e26",
          "finding_ids": [
            "F1",
            "F2",
            "F3"
          ],
          "supersedes": []
        },
        {
          "id": "gs1-r1-spec",
          "performer": "claude:gsl-review-spec-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gsl-review-spec-r1",
            "digest": "sha256:3ed5204e2bf3f39b8fba09989188649698fe07f563dbf38a3ae22395ef11d5e6",
            "excerpt": "SPEC axis, GS1 (base e4db0d2e6, code tip eee733317). check-current: current=true.\nS1 | ask-user | spec.md:314, :293, :271; scan_loop_test.go:129-130 | conf 6. Family S omits the \"between candidates\" member; Command's timeout cannot place it and the old select cannot take a read-count cancel. GS05 covers it by a direct Scan. The spec records no waiver.\nS2 | ask-user | guard_scan_owner_test.go:39-45 vs spec.md:130, :373 | conf 5. The ownership check reads only the top level of internal/guards; a nested production package gets no diagnostic. Recursive walk that skips testdata, or a spec amendment to \"directly in\".\nS3 | auto-fix | scan_loop_test.go:30, :142-143 | conf 4. The inspection-cancelled family S members depend on a 10ms wall-clock window; a CPU stall changes the tuple. spec.md:304, :326 forbid timing-defined expectations.\nS4 | no-op | scan_loop_test.go:379-391 vs spec.md:94 | conf 3. GS11 uses a 20ms hold, not a release barrier, but it observes the contract.\nS5 | no-op | scan_loop_test.go:393-399 | conf 3. GS22 observes its stated behavior; GS23 backs it.\nCount 5 (2 ask-user, 1 auto-fix, 2 no-op). Worst: S1. All 30 rows observed by their cited tests except as S1 notes.\nAdvice: checkpointContext tests pin the exact count of ctx.Err reads; the check does not catch `for range ch` or `new(chan T)`, which the spec does not require.\n"
          },
          "axis": "Spec",
          "base": "e4db0d2e66ffae4045c9c29c11c4696462442bc4",
          "tip": "eee733317545f94b7df5efcc56f23e1e93437e26",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "gs1-r1-coverage",
          "performer": "claude:gsl-review-coverage-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gsl-review-coverage-r1",
            "digest": "sha256:d9fe84149f7458b2b6864d06f5f1fc480fc6065b56eabd0824b11e3a1dedda47",
            "excerpt": "COVERAGE axis, GS1 (base e4db0d2e6, code tip eee733317). check-current: current=true.\nC1 | auto-fix | internal/conformance/guard_scan_owner_test.go:41-45 vs spec.md:130, :373 | conf 7. The check reads only direct children of internal/guards; a nested pipeline gets no diagnostic.\nC2 | auto-fix | guard_scan_owner_test.go:52; projects/benchkit.md:226 | conf 6. A FIFO or special file named *.go reaches parser.ParseFile and can hang; the hostile-input checklist requires rejection before read.\nC3 | ask-user | guard_scan_owner_test.go:83-84; spec.md:131 | conf 5. make(namedChanType) is not reported as channel construction; a fix needs go/types.\nC4 | ask-user | guard_scan_owner_test.go:70-90; spec.md:240, :131 | conf 5. sync.WaitGroup.Go, context.AfterFunc, and time.AfterFunc start goroutines without a go statement; no decision covers them.\nC5 | auto-fix | guard_scan_owner_test.go:36; projects/benchkit.md:216-223 | conf 5. Every os.Stat error reads as absent; only ErrNotExist should take the optional-surface exit.\nC6 | no-op | scan_loop_test.go:393-399 | conf 6. GS23 closes the second-path failure that GS22 names.\nC7 | no-op | scan_loop_test.go:324-328, :365, :373, :379 | conf 6. Preservation rows without a demonstrated red are legitimate regression guards.\nC8 | no-op | scan_loop_test.go:331, :339-342 | conf 6. GS03, GS05, GS07 reds derived by reading.\nGS21 invalid-arguments red demonstrated: probe swap at guards.go:274 read bit, first failure TestScanLoopDifferentialCommandFamilies/invalid-arguments, restored yes.\nGolden review: f5c12ccbc changed no guards.go; golden bytes unchanged to eee733317; faithful capture of the old pipeline; no member tailored to the new code.\nCount 8 (5 actionable, 3 no-op). Worst: C1.\nAdvice: 10ms cancelTick and 20ms operationHold could flake on a loaded runner; the live-tree test is the only catch for a renamed internal/guards.\n"
          },
          "axis": "Coverage",
          "base": "e4db0d2e66ffae4045c9c29c11c4696462442bc4",
          "tip": "eee733317545f94b7df5efcc56f23e1e93437e26",
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
          "id": "gs1-r2-standards",
          "performer": "claude:gsl-review-standards-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gsl-review-standards-r2",
            "digest": "sha256:88b4c01c46d99a61273803083cd84139ad54efacb8076dc2980ac8918bf9a3c9",
            "excerpt": "STANDARDS confirming round, GS1, repair delta eee73331..cdeaf570. check-current: current=true.\nF1: closed. scan_loop_test.go:24 states only that familySGolden is the reference capture of scan family S.\nF2: closed. guard_scan_owner_test.go:33-41 holds one constant for each diagnostic piece; the check (:55, :63, :68, :75, :81) and the tests (:207, :218, :262, :285) share them; no respelled literal remains in the file.\nNew findings: 0.\nAdvice: the cancelTick comment at scan_loop_test.go:27 is slightly wide; each blocked-* member now waits about 3s in total; guard_scan_owner_test.go:242 and :180 respell construct labels, confidence 3, not a block.\n"
          },
          "axis": "Standards",
          "base": "e4db0d2e66ffae4045c9c29c11c4696462442bc4",
          "tip": "cdeaf570831bad3feee048c7bc74922e570ee3a9",
          "finding_ids": [],
          "supersedes": [
            "gs1-r1-standards"
          ]
        },
        {
          "id": "gs1-r2-spec",
          "performer": "claude:gsl-review-spec-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gsl-review-spec-r2",
            "digest": "sha256:3e063b8db993770e18051f7f796e1b072ca96e78baa49a1b0ed1d6808b9a39ec",
            "excerpt": "SPEC confirming round, GS1, repair delta eee73331..cdeaf570 and spec corrections 2854ddb4..06216f4c. check-current: current=true.\nS1: closed. spec.md:314 gives the between-candidates member to GS05 by direct Scan; familyS (scan_loop_test.go:129-144) omits it; GS27 (spec.md:293) agrees.\nS2 (C1): closed. checkGuardScanLoopOwner walks every production file beneath internal/guards and skips testdata (guard_scan_owner_test.go:58-73), as spec.md:130 says.\nS3: closed. The inspection-cancelled members use inspectionWindow (scan_loop_test.go:30, :142-143); scan-family-s.golden is byte-identical across the repair.\nSpec corrections match the repaired check (spec.md:131, :338-339, :375).\nN1 | ask-user, orchestrator disposition no-op | spec.md:135 | conf 5. \"AST nodes and real import bindings\" versus the AST-only check. The fable tie-break of the initial round decided that this sentence bans text matching, and the line stays unchanged.\nN2 | no-op | guard_scan_owner_test.go:50-73 | conf 6. The stat, unreadable-subdirectory, and non-regular-file diagnostics harden the check in the spirit of spec.md:138.\nCount 2 new, 0 blocking, 0 open named findings.\nAdvice: the blocking-operation member still uses the 10ms cancelTick (scan_loop_test.go:144), outside S3.\n"
          },
          "axis": "Spec",
          "base": "e4db0d2e66ffae4045c9c29c11c4696462442bc4",
          "tip": "cdeaf570831bad3feee048c7bc74922e570ee3a9",
          "finding_ids": [],
          "supersedes": [
            "gs1-r1-spec"
          ]
        },
        {
          "id": "gs1-r2-coverage",
          "performer": "claude:gsl-review-coverage-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "748b0101d016df166d50a31f68f4df1f475c8edc",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gsl-review-coverage-r2",
            "digest": "sha256:0e46addf896bf14ccf0e566e2faf03191970c0cef94b774998ba3d0b8f6de01e",
            "excerpt": "COVERAGE confirming round, GS1, repair delta eee73331..cdeaf570. check-current: current=true.\nC1: closed. guard_scan_owner_test.go:77-92 walks beneath internal/guards and skips testdata at any depth; TestGuardScanLoopOwnerGradesNestedProductionFiles demands only the nested diagnostic.\nC2: closed. A non-regular .go entry is named before any open (:86); TestGuardScanLoopOwnerRefusesASpecialGoFile ran with a real FIFO and a bounded deadline.\nC5: closed. Only fs.ErrNotExist takes the nil exit (:62); TestGuardScanLoopOwnerReportsAnUnstatableOwnerDir passed.\nN1 | ask-user | guard_scan_owner_test.go:81-82 | conf 8. The WalkDir error branch has no witness: a swap of its diagnostic append for `_ = rel` read silent, and all 10 TestGuardScanLoopOwner* tests stayed green.\nN2 | no-op | guard_scan_owner_test.go:77-92 | conf 4. WalkDir does not follow a symlinked directory; the repair did not widen this.\nN3 | no-op | guard_scan_owner_test.go:62 | conf 3. A regular file named internal/guards still reads as absent, as before the repair.\nN4 | no-op | reviews/guard-scan-loop.md | conf 3. The record fence holds no red-to-green entry for the new C2 and C5 witnesses.\nCount 4 (1 ask-user, 3 no-op). Worst: N1.\nCommands: bench test --package ./internal/conformance --run '^TestGuardScanLoopOwner' passed 10; one bench probe swap at :82 read silent, restored yes.\n"
          },
          "axis": "Coverage",
          "base": "e4db0d2e66ffae4045c9c29c11c4696462442bc4",
          "tip": "cdeaf570831bad3feee048c7bc74922e570ee3a9",
          "finding_ids": [
            "N1"
          ],
          "supersedes": [
            "gs1-r1-coverage"
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
      "from": "sha256:08e457bec47256e4119881f481cd37ea8d054b3070bd8602c6cd7b5c61443dce",
      "to": "sha256:accf0d321caa04266e1927a4e654e5a04c4eec6af9019f85ac0fda07ae04341f",
      "chunk_ids": {
        "GS1": [
          "GS1"
        ]
      }
    },
    {
      "from": "sha256:accf0d321caa04266e1927a4e654e5a04c4eec6af9019f85ac0fda07ae04341f",
      "to": "sha256:552f0b3d19b40a0481174b21ae276df18c64cfbc7270fc30f94b0b66660045cd",
      "chunk_ids": {
        "GS1": [
          "GS1"
        ]
      }
    }
  ]
}
```
