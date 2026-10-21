# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/guard-scan-loop/spec.md",
  "plan_digest": "sha256:08e457bec47256e4119881f481cd37ea8d054b3070bd8602c6cd7b5c61443dce",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GS1",
      "base": "e4db0d2e66ffae4045c9c29c11c4696462442bc4",
      "tip": "eee733317545f94b7df5efcc56f23e1e93437e26",
      "plan_digest": "sha256:08e457bec47256e4119881f481cd37ea8d054b3070bd8602c6cd7b5c61443dce",
      "source_digest": "d360efb6c316c4523f5c5d342d3d3756b00856ce",
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
