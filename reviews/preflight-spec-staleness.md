# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/preflight-spec-staleness/spec.md",
  "plan_digest": "sha256:4be92c9f220cff6e03206483328acc81fd3de1894913f7c9c4df48caa0d0a87b",
  "implementation_session": "",
  "chunks": [
    {
      "id": "SS-C1",
      "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
      "tip": "c585f9d08d9f7d320d5e56dfc88c69a09f43733c",
      "plan_digest": "sha256:e74aa34bfd73177aca13de8907cef226127cc11275da45fa868160f4cd059831",
      "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
      "acceptance_rows": [
        "SS1",
        "SS2",
        "SS3",
        "SS4",
        "SS5",
        "SS6",
        "SS7",
        "SS8",
        "SS9",
        "SS10",
        "SS11",
        "SS12",
        "SS35",
        "SS36",
        "SS37",
        "SS38",
        "SS39",
        "SS40",
        "SS51",
        "SS52",
        "SS56"
      ],
      "verification": [
        {
          "id": "ft375_t01-ss-01-focused",
          "performer": "claude:ft375_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b4132eb97e1d6dfbd35824b91b6771d7fc431640",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01",
            "digest": "sha256:94faca3df1c1fed9a935eee458d498bee04472b49d8239a646b5700ad6e2aca9",
            "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  2,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  3,bit,swap,failed,2,TestDriftIdentity,5,1,0,null,yes\n  4,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  5,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestDriftObjects|TestDriftIdentity|TestAuditMovement,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,4503\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,615,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:151: drifted members = (1)\",25\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,534,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:151: drifted members = (1)\",25\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,451,5\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:172: SS9 spec drift = (0) [], want the audit red alone and \\\"specs/example/spec.md|specs/example/spec.md|modified|c55b53562e122867fb07f345a99cca80bb6936cb|d52507c627fa40f185341a5acc139938da2076de\\\":\",16\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:151: drifted members = (1)\",24\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,543,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:151: drifted members = (1)\",27\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,533,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:151: drifted members = (1)\",26\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestDriftObjects|TestDriftIdentity|TestAuditMovement'",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-01-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t01",
              "digest": "sha256:94faca3df1c1fed9a935eee458d498bee04472b49d8239a646b5700ad6e2aca9",
              "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  2,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  3,bit,swap,failed,2,TestDriftIdentity,5,1,0,null,yes\n  4,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  5,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestDriftObjects|TestDriftIdentity|TestAuditMovement,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,4503\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,615,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:151: drifted members = (1)\",25\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,534,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:151: drifted members = (1)\",25\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,451,5\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:172: SS9 spec drift = (0) [], want the audit red alone and \\\"specs/example/spec.md|specs/example/spec.md|modified|c55b53562e122867fb07f345a99cca80bb6936cb|d52507c627fa40f185341a5acc139938da2076de\\\":\",16\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:151: drifted members = (1)\",24\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,543,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:151: drifted members = (1)\",27\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,533,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:151: drifted members = (1)\",26\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t01-ss-01-check-1",
          "performer": "claude:ft375_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b4132eb97e1d6dfbd35824b91b6771d7fc431640",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01",
            "digest": "sha256:54e3743a38e15c7e1b81447a3612565b5b0addd244eb378e7eb8943bed8c486f",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,39306,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23808,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,41261,270\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,1277,12\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t01-ss-01-check-2",
          "performer": "claude:ft375_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b4132eb97e1d6dfbd35824b91b6771d7fc431640",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01",
            "digest": "sha256:b901c7e7442b0c3816e52dd52004d59989468752cbfe8da62e8517061bcdf8e0",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14510,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-2",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t01-ss-01-check-3",
          "performer": "claude:ft375_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b4132eb97e1d6dfbd35824b91b6771d7fc431640",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01",
            "digest": "sha256:e6328747dd5cc33a956cc78ab8d31437301618d4b24ddb5b55a8f40f8c8ab0b9",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3981,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-3",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft375_t01-ss-01-check-4",
          "performer": "claude:ft375_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b4132eb97e1d6dfbd35824b91b6771d7fc431640",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01",
            "digest": "sha256:bcd102b88170c169644262a505c91acfb98a13dad884d6b859bf94e2c3d23505",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14595,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t01-ss-01-check-5",
          "performer": "claude:ft375_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b4132eb97e1d6dfbd35824b91b6771d7fc431640",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01",
            "digest": "sha256:b184e547bc57048a3a9c9fbfb4f5fc5d40966032c81232c084b6868d30768a66",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,34532,1118\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ss-01-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1-ss-01-focused",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0e41de513c4f3d907df0be76f81b2de273dd76c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:6908e8536fcb63ce4c8d3afdb8357f10eacc8f19048aeeeb0df9f37225034e18",
            "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  2,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  3,bit,swap,failed,2,TestDriftIdentity,5,1,0,null,yes\n  4,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  5,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestDriftObjects|TestDriftIdentity|TestAuditMovement,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,4503\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,616,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,590,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,453,5\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:174: SS9 spec drift = (0) [], want the audit red alone and \\\"specs/example/spec.md|specs/example/spec.md|modified|c55b53562e122867fb07f345a99cca80bb6936cb|d52507c627fa40f185341a5acc139938da2076de\\\":\",16\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",24\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,587,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",28\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,605,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",27\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestDriftObjects|TestDriftIdentity|TestAuditMovement'",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-01-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t01_rep1",
              "digest": "sha256:6908e8536fcb63ce4c8d3afdb8357f10eacc8f19048aeeeb0df9f37225034e18",
              "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  2,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  3,bit,swap,failed,2,TestDriftIdentity,5,1,0,null,yes\n  4,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  5,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestDriftObjects|TestDriftIdentity|TestAuditMovement,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,4503\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,616,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,590,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,453,5\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:174: SS9 spec drift = (0) [], want the audit red alone and \\\"specs/example/spec.md|specs/example/spec.md|modified|c55b53562e122867fb07f345a99cca80bb6936cb|d52507c627fa40f185341a5acc139938da2076de\\\":\",16\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",24\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,587,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",28\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,605,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",27\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t01_rep1-ss-01-check-1",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0e41de513c4f3d907df0be76f81b2de273dd76c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:bd5c76e8cf48261c503694492a260a6318bfc7d685adb4dad81ceae6fe4850f3",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,39244,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23982,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,41047,270\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,1394,13\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1-ss-01-check-2",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0e41de513c4f3d907df0be76f81b2de273dd76c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:df0126cd60ef4b933b5ed536c8254322e1f0ed328943b7f17a432496f0a4d2a0",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15201,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-2",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1-ss-01-check-3",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0e41de513c4f3d907df0be76f81b2de273dd76c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:0304c806436e3331b751bfdf1028fc96e8e780f90ed19e529247f9c55edd1cb7",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3906,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-3",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1-ss-01-check-4",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0e41de513c4f3d907df0be76f81b2de273dd76c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:fb8ff631426dfee6a1c0d7944eb4c183b47dc9ff4e377d854e1d9c9da60ad751",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14925,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1-ss-01-check-5",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0e41de513c4f3d907df0be76f81b2de273dd76c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:1406148d6668449a0d23f7e0f1773fa276ab6908e04eda92f776456221da262d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,34808,1118\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ss-01-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1-ss-01-check-6",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0e41de513c4f3d907df0be76f81b2de273dd76c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:af9441ffb2fd7f646464fb44fcd77e235016e38e7d0636dd061529925d8415e3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,1273,172\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-6",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1b-ss-01-focused",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "410a4337ba0bc62d31ab5f62a7c55710837caaaf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:ee30298d638e7966a9473d23921e681ea8f45b49eb84fafdd591b3a0075ede22",
            "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  2,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  3,bit,swap,failed,2,TestDriftIdentity,5,1,0,null,yes\n  4,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  5,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestDriftObjects|TestDriftIdentity|TestAuditMovement,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,4503\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,592,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,583,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,450,5\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:174: SS9 spec drift = (0) [], want the audit red alone and \\\"specs/example/spec.md|specs/example/spec.md|modified|c55b53562e122867fb07f345a99cca80bb6936cb|d52507c627fa40f185341a5acc139938da2076de\\\":\",16\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",24\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,604,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",28\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,584,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",27\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestDriftObjects|TestDriftIdentity|TestAuditMovement'",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-01-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t01_rep1",
              "digest": "sha256:ee30298d638e7966a9473d23921e681ea8f45b49eb84fafdd591b3a0075ede22",
              "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  2,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  3,bit,swap,failed,2,TestDriftIdentity,5,1,0,null,yes\n  4,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\n  5,bit,swap,failed,1,TestDriftObjects,5,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestDriftObjects|TestDriftIdentity|TestAuditMovement,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,4503\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,592,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,583,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,450,5\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:174: SS9 spec drift = (0) [], want the audit red alone and \\\"specs/example/spec.md|specs/example/spec.md|modified|c55b53562e122867fb07f345a99cca80bb6936cb|d52507c627fa40f185341a5acc139938da2076de\\\":\",16\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",24\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,604,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",28\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,584,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftObjects,\"drift_test.go:153: drifted members = (1)\",27\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t01_rep1b-ss-01-check-1",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "410a4337ba0bc62d31ab5f62a7c55710837caaaf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:3db7c888de6344d104dd95b398cd7fb7ce8e59038184ee9bc6f0e1a7dbbe04dc",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,39537,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23713,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,41408,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,1443,13\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1b-ss-01-check-2",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "410a4337ba0bc62d31ab5f62a7c55710837caaaf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:00f40c4915118ac757ea85a8e3ef619a1d7e4ed5edf4700c4a8b469f38ce8343",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14458,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-2",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1b-ss-01-check-3",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "410a4337ba0bc62d31ab5f62a7c55710837caaaf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:7e789b6d111943e77ca0cf9d4bd4496dc735deb64d2f15cc06e776cc1892ee26",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3967,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-3",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1b-ss-01-check-4",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "410a4337ba0bc62d31ab5f62a7c55710837caaaf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:67cda0136049d2b7997ad3c45749ec4158b083f39603d16752db9ccb8410be91",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14737,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1b-ss-01-check-5",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "410a4337ba0bc62d31ab5f62a7c55710837caaaf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:e9e522ebc669842e2738bd6e4c4716fe192d8dda612932839a01995e5233bd6e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,34425,1118\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ss-01-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep1b-ss-01-check-6",
          "performer": "claude:ft375_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "410a4337ba0bc62d31ab5f62a7c55710837caaaf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep1",
            "digest": "sha256:785c5facb67a7bd70cbd8b5fd7147e3adb1458d38c6e752a82a573759b528285",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,1294,172\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-6",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep2-ss-01-focused",
          "performer": "claude:ft375_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep2",
            "digest": "sha256:d5bbccfe8eab0cbc5b4d21ef11c1d9e0d1f296af19c8b80bf4962a5569d98aba",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDriftIdentity,5,1,0,null,yes\n  2,bit,swap,failed,1,TestDriftIdentity,5,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestDriftObjects|TestDriftIdentity|TestAuditMovement,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2122\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,605,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:217: dropped fence drift = [\\\"specs/example/spec.md|specs/example/spec.md|modified|99b3d1ec3a9b920fbc394f341ec73e805909cbc6|c55b53562e122867fb07f345a99cca80bb6936cb\\\"], want [\\\"outside/fenced/member.txt|outside/fenced/member.txt|… (470 bytes)\",19\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,546,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:206: dropped claim drift = [\\\"specs/example/spec.md|specs/example/spec.md|modified|e3b972d9c6fde23b8a492121e09650552f419d48|4982493dba3b3a59aa78ce7f2d2c09526c65642f\\\"], want [\\\"outside/claim_test.go|outside/claim_test.go|modified… (589 bytes)\",19\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestDriftObjects|TestDriftIdentity|TestAuditMovement'",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-01-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t01_rep2",
              "digest": "sha256:d5bbccfe8eab0cbc5b4d21ef11c1d9e0d1f296af19c8b80bf4962a5569d98aba",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDriftIdentity,5,1,0,null,yes\n  2,bit,swap,failed,1,TestDriftIdentity,5,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestDriftObjects|TestDriftIdentity|TestAuditMovement,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2122\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,605,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:217: dropped fence drift = [\\\"specs/example/spec.md|specs/example/spec.md|modified|99b3d1ec3a9b920fbc394f341ec73e805909cbc6|c55b53562e122867fb07f345a99cca80bb6936cb\\\"], want [\\\"outside/fenced/member.txt|outside/fenced/member.txt|… (470 bytes)\",19\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,546,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:206: dropped claim drift = [\\\"specs/example/spec.md|specs/example/spec.md|modified|e3b972d9c6fde23b8a492121e09650552f419d48|4982493dba3b3a59aa78ce7f2d2c09526c65642f\\\"], want [\\\"outside/claim_test.go|outside/claim_test.go|modified… (589 bytes)\",19\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t01_rep2-ss-01-check-1",
          "performer": "claude:ft375_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep2",
            "digest": "sha256:c217b010954733a4266fdea20c999d67d820a47b95fccfd95958b451458edfc0",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,39463,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,24175,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,41401,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,1470,13\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep2-ss-01-check-2",
          "performer": "claude:ft375_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep2",
            "digest": "sha256:a39edce2839bc4f16672ef5ee945204e8fe04227bcaa61d4efda2da832dff42f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14430,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-2",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep2-ss-01-check-3",
          "performer": "claude:ft375_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep2",
            "digest": "sha256:36082a3ab4b55d9f88a4f419b3f91c1cb2e96d383a83bafd2bc734ea6c04f5e9",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3897,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-3",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep2-ss-01-check-4",
          "performer": "claude:ft375_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep2",
            "digest": "sha256:933dd03b0ca121b6a613323a86e5df5e72003a7ca55e67f71ce5d0afbaa76e99",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14413,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep2-ss-01-check-5",
          "performer": "claude:ft375_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep2",
            "digest": "sha256:e2b03b2c8b9f2f915d921204227acb484b2e07e809d243394a25454f96be694a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,34311,1118\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ss-01-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep2-ss-01-check-6",
          "performer": "claude:ft375_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep2",
            "digest": "sha256:074a13dc9f59985376e4a6f80b80014c274514f668de81719ef19f51808692fc",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,1269,172\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-6",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep3-ss-01-focused",
          "performer": "claude:ft375_t01_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep3",
            "digest": "sha256:2378a87fbf4eed6716f9076d78537925b5aa65441e6d01fb4d984dc5857d419b",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestDriftIdentity,5,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestDriftObjects|TestDriftIdentity|TestAuditMovement,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1061\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,632,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:221: dropped fence drift = [\\\"outside/fenced/member.txt|outside/fenced/member.txt|modified|5626abf0f72e58d7a153368ba57db4c673c0e171|f719efd430d52bcfc8566a43b2eb655688d38871\\\" \\\"specs/example/spec.md|specs/example/spec.md|modified… (744 bytes)\",20\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestDriftObjects|TestDriftIdentity|TestAuditMovement'",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-01-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t01_rep3",
              "digest": "sha256:2378a87fbf4eed6716f9076d78537925b5aa65441e6d01fb4d984dc5857d419b",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestDriftIdentity,5,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestDriftObjects|TestDriftIdentity|TestAuditMovement,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1061\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,632,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestDriftIdentity,\"drift_test.go:221: dropped fence drift = [\\\"outside/fenced/member.txt|outside/fenced/member.txt|modified|5626abf0f72e58d7a153368ba57db4c673c0e171|f719efd430d52bcfc8566a43b2eb655688d38871\\\" \\\"specs/example/spec.md|specs/example/spec.md|modified… (744 bytes)\",20\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t01_rep3-ss-01-check-1",
          "performer": "claude:ft375_t01_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep3",
            "digest": "sha256:fd9d84db030f3ac62320115532db777c8a8cc3aea9958ee59386d3745ce259dc",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,38904,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23830,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,41068,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,1457,13\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep3-ss-01-check-2",
          "performer": "claude:ft375_t01_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep3",
            "digest": "sha256:df823087f65204d8a25dc47f51c4344b4e113acba42102253b85d07721ede41c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14407,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-2",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep3-ss-01-check-3",
          "performer": "claude:ft375_t01_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep3",
            "digest": "sha256:64971fdcd754dfaaf3ebb78f1ee4caff44b556154c31cc55d79b4660cc5e3bde",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3984,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-3",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep3-ss-01-check-4",
          "performer": "claude:ft375_t01_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep3",
            "digest": "sha256:7a18506e30b7509815f4b18d316809cc913a5bedf8db54310762998f1f2e036f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14726,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep3-ss-01-check-5",
          "performer": "claude:ft375_t01_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep3",
            "digest": "sha256:d12eed35f2893732f3bafeef92f4197d8c45d603b51b2f2e5e33194ffaf5b537",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,34517,1118\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ss-01-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t01_rep3-ss-01-check-6",
          "performer": "claude:ft375_t01_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t01_rep3",
            "digest": "sha256:956999cfe866093f106db36d436625d1ff6fb7f8b9f595c3005f1a43e0741dce",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,1299,172\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-01-check-6",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "ft375_ss-c1_standards_r1",
          "performer": "claude:ft375_ss-c1_standards_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b4132eb97e1d6dfbd35824b91b6771d7fc431640",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_standards_r1",
            "digest": "sha256:edafc057380d0b5c83dec62be90986602a21b8c627ee95530b58e21cfbfa410a",
            "excerpt": "Standards: 3 findings; worst: S1\n\nEvidence sha256:af29…e6 binds (`--check-current`: current=true, dirty_record=false, head e70242684). Delta read once via `git diff 81b10e749..7a1535fdf` (non-specs part in full, plus stat). I fetched no evidence consumer rows: the source greps answered every question I had.\n\n- **S1** `internal/preflight/staleness/collect.go:379-387` (`subjectPath`, called at :98). Rule: AGENTS.md \"Knowledge duplication is a defect… a parser and its count\". The manifest says the subject grammar belongs to each premise's collector (`ownership/manifest.go:74-75`). Evidence: `subjectPath` re-derives that grammar and reads only a top-level `\"path\"`. Several subject kinds put their path elsewhere, so their paths silently never join audit scope:\n  - nested `owner` (`ownership/helpers.go:228`, `rules.go:29`, `members.go:50-52`)\n  - `caller`/`target`/`owner` (`symbols.go:214`, `:315`)\n  - `directory` (`sizing.go:150`)\n  - `pattern`/`scope` (`sizing.go:84`)\n  This is also a concrete gap against spec line 118 (\"claim source files\") and SS12 (\"cannot hide a changed premise owner\"). Static only, because the charge forbids probes. Disposition: auto-fix (one ownership-owned subject-path projection; may need an in-scope fence plan commit). Confidence 7.\n- **S2** `internal/preflight/staleness/collect.go:393-415` (`members`). Rule: \"one source per fact… a parser\". Evidence: it hand-parses `ls-tree -z` records (`meta\\tpath`, 3 fields, `IndexEntry.IsRegularFile`). `ReadTreeFile` already parses that same grammar at `internal/git/tree.go:169-171`. The gate has a third private full-listing parser (`internal/gate/tree_snapshot.go:131`, `parseTreeSnapshot`). Disposition: auto-fix (one exported `git` tree-listing parser; fence expansion via plan commit if `internal/git` is unfenced). Confidence 6.\n- **S3** `internal/preflight/staleness/drift_test.go:620-624`. Rule: craft-comments, \"no argument for its own correctness… The commit or spec owns the red record. Do not retain a mutation transcript or test result in a code comment.\" Evidence: the header says \"A dropped member, a collapsed kind, or a renamed column turns a case red.\" That is an unrecorded red claim addressed to the reviewer. Disposition: auto-fix (drop the sentence; record the red where the AGENTS.md test-expectation exception requires it). Confidence 6.\n\nThe `staleness.Collect(root string, …)` signature breaks no standard: it matches `ownership.Collect` and spec line 79.\n\nAdvice\n- `gather_inputs.go:67`: the audit form reuses the hint \"review source has uncommitted changes\" for a build audit, which misleads.\n- `gather.go:166` reuses `ownershipFailure` for an audit failure; a separate name reads more clearly.\n- `gather.go:155` parses the manifest from `content` a second time (the first parse is in `closure.go:117`). Passing the parsed value is cheaper.\n- The test's `table` (`drift_test.go:660-674`) adds a local TOON block splitter. If whole-response decoding fails, a shared `preflighttest` helper beats a per-package one.\n- `staleness.Error` repeats `ownership.Error` (`ownership.go:159`). The spec mandates it, so it is no finding.\n\nRead: `git diff` stat and non-specs delta; `internal/git/tree.go:100-210`; `internal/gate/tree_snapshot.go:60-140`; `internal/preflight/ownership/manifest.go` (whole), `premises.go:280-323`, `members.go:1-90`, `source_test.go:20-55`; `specs/preflight-spec-staleness/spec.md:60-127`, plus rg hits at :289, :406, :554; `.agents/skills/bench-craft-comments/SKILL.md`; `.agents/skills/bench-craft-review/references/finding-discipline.md`; rg sweeps of `internal/` for ls-tree, rename, subject-decoding and helper duplicates. Evidence cursors fetched: none.\n\nNot verified:\n- Whether a red record exists for the test's independent expectations (`absent`, table headers) under the AGENTS.md exception.\n- S1 by an executable run.\n- `internal/prose` relevance: the delta has no Markdown check.\n- `bench-craft-review/SKILL.md` itself (I read only its finding-discipline reference) and `projects/benchkit.md`.\n- The spec/ticket hunks in the range.\n"
          },
          "axis": "Standards",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "7a1535fdf85b124cd79d9106747269acadba764a",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "ft375_ss-c1_spec_r1",
          "performer": "claude:ft375_ss-c1_spec_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b4132eb97e1d6dfbd35824b91b6771d7fc431640",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_spec_r1",
            "digest": "sha256:418137b3282ec987df2938b4e2f6e588e758a3e6464fbb273c64311387ec5d67",
            "excerpt": "Spec: 1 finding; worst: P1\n\nThe evidence binds: `--check-current` returned current=true at e70242684, dirty_record=false. Delta read: 81b10e749..7a1535fd (8 commits).\n\n**P1. The chunk's spec edit left eight planning-metadata seam locations that no longer resolve.**\n- Requirement (spec.md:178): \"An omitted citation, unrecognized table, partial inventory, manual_claims entry, or absent metadata makes Complete false.\" Spec.md:621 also says: \"This authored inventory is declared complete…\"\n- What changed: commit 4ac8c6f3c rewrote the SS1–SS12, SS35–SS40, SS51, SS52 and SS56 coverage cells into backticked form. It then moved each `planned-TestDrift*` and `planned-TestAudit*` seam location from its coverage row (line 391/395) to line 636, the JSON line itself. That line holds the old needle more than once.\n- Evidence: I ran the read-only `bench preflight build preflight-spec-staleness --audit-spec --base 81b10e74 --source-tip e7024268`. It exits 1 and prints 8 rows of the form `metadata,\"seam planned-TestDriftObjects: …spec.md line 636 does not hold its needle … once\"`. Those 8 rows cover TestDriftObjects, TestDriftIdentity, TestAuditPrecedence, TestAuditMovement, TestAuditTrust, TestAuditGrammar, TestAuditSafety and TestAuditBound.\n- Before the chunk: at 81b10e749 the TestDriftObjects seam pointed at line 391 with a needle that matched that coverage row. I checked only that one seam's text; I did not run the audit at the base.\n- Impact: the declared-complete inventory now carries diagnostics that SS-C2 and SS-C3 cannot clear. Whether they refuse ticket 02's charge is unknown.\n- Fix: point each seam (and each acceptance decision `source`) at its coverage-map row with the backticked needle.\n- Disposition: auto-fix. Confidence: 7.\n\n**Rows**\n- SS1–SS6 and SS11 met (drift.go:534–592, collect.go:393–417, TestDriftObjects).\n- SS10 met as SS-C1 scopes it. The empty table prints, but complete stays false. Ticket 01:33 and spec.md:294 decide this, so it is a no-op.\n- SS7–SS9 and SS12 met (TestDriftIdentity, collect.go:288–355).\n- SS35–SS38 met (gather_inputs.go:235 dirty check runs before gather; the existing movement, tip-current and binary-seal rows still apply).\n- SS39–SS40 met (operations.go registry row with `required: FlagBase, FlagTip`).\n- SS51 met (collect.go:411 refuses a member that is not a regular file).\n- SS52 met (KindAuditSpec is unbounded, so the root spill applies; TestAuditBound).\n- SS56 met (the audit does not touch the prepared pack).\n- SS-C1 never claims omission: `unsupported` always adds cited-symbols and cited-tests diagnostics (collect.go:280, :310).\n- `Collect(root, …)`: no behavior change. It reads only committed objects at source.Base and source.Tip, inside the movement-checked attempt, and it follows the ownership.Collect convention (ownership.go:169). Not a finding.\n- Scope creep: none. Every production path is inside ticket 01's `Writes:` line.\n\n**Advice**\n- The tip manifest comes from the checkout bytes (gather.go:155). The tip spec bytes come from git (collect.go:293). Taking the tip manifest from `ParseManifest(current)` would keep both inputs on committed bytes. This only differs when tip ≠ HEAD, which tip-current already turns red.\n- drift.go pairs a deleted path with an added path by blob ID alone. Two unrelated files with identical bytes would print as a rename.\n- The audit's dirty refusal reuses the hint \"review source has uncommitted changes\" in build mode.\n- My first audit run refused \"source not clean\". The probes on the shared tree probably caused this; the retry was clean.\n\nRead: spec.md (all of it), tickets/01, bench-craft-review SKILL.md (grep) and references/finding-discipline.md, the chunk diff, command.go:100–214, gather.go:30–110, ownership.go:15–39, spec/fences.go FenceTokens, bench coverage output, git log of the range.\n\nNot verified: test execution and probes (the Coverage axis owns them); whether the P1 diagnostics refuse a later charge; the audit at base..base.\n"
          },
          "axis": "Spec",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "7a1535fdf85b124cd79d9106747269acadba764a",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "ft375_ss-c1_coverage_r1",
          "performer": "claude:ft375_ss-c1_coverage_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b4132eb97e1d6dfbd35824b91b6771d7fc431640",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_coverage_r1",
            "digest": "sha256:84380848f63cac4bf3d19026f762fff7e6266a141ff14c19a9d4364e98c6e3dd",
            "excerpt": "Coverage: 4 findings; worst: C1\n\nEvidence check: `--check-current` bound (current=true, head e70242684, dirty=false).\n\n- **C1**: the baseline spec has a fence entry or claim source that the tip spec drops, and a member under it changes. Break: if the collector reads the tip spec's fences and manifest only, the scope loses that member and its drift row disappears. Spec 111 says \"at both identities\"; the code is collect.go:297. Needed test: a TestDriftIdentity case where the tip removes a fence and changes its member, and the member row is expected. Probe: silent. Disposition: auto-fix (test). Confidence 8.\n- **C2**: a special file (link) exists only at the baseline. Break: if the baseline members refusal is ignored, the \"special before source\" refusal is lost. SS51 says \"before/after\", but TestAuditSafety (drift_test.go:898) places the link only at the tip. The code is collect.go:298-301. Needed test: a base-side link case in TestAuditSafety. Probe: silent. Disposition: auto-fix. Confidence 8.\n- **C3**: a member path holds `\"`, `\\`, a tab, or a non-ASCII byte. Break: if the listing uses newline records instead of `-z` (collect.go:394,399), git quotes the path and the fence never matches it. That gives a false empty drift. The requirement is Edge inventory spec:463 (\"retains spaces and metacharacters literally\"). Needed test: such a member in TestDriftObjects. Probe: silent. Disposition: auto-fix. Confidence 7.\n- **C4**: the tip differs from the checkout HEAD, and the two spec manifests differ. Break: gather.go:110 passes `content`, which is the checkout spec through resolveSpecInput (gather_inputs.go:90-92), as the \"tip spec\" manifest. The Collect doc (collect.go:284) says the manifest is the tip spec's, and spec:127 says \"Worktree content cannot replace committed source bytes\". The scope then gets claim sources from HEAD, not from the tip. The red tip-current row from SS37 masks this, so no false green occurs. No row covers it. Probe: none. Disposition: ask-user (the likely fix is `ParseManifest(current)` in Collect). Confidence 5.\n\n## Advice\n- The tickets-folder scope entry (collect.go:331) is unguarded: probe silent. SS31 in SS-C2 (ticket 04) owns that row, so confirm its test changes a ticket file.\n- No test pairs `--audit-spec` with `--ticket`, `--max-store-bytes`, or `--cursor`. I did not verify how the grammar rejects them.\n\n## Probes run\nAll probes used `--package ./internal/preflight/... --run 'TestDrift|TestAudit'`, and every one restored the file (restored=yes).\n1. Baseline fences and manifest dropped from scope: silent.\n2. Tickets entry swapped for the spec path: silent.\n3. Baseline members error ignored (occurrence 3): silent.\n4. `-z` dropped and records split on newline: silent.\n5. Two earlier attempts were invalid (build-failed, unused var/import); the retries replace them.\n\nThe tree is clean: `git status --short` prints nothing.\n\n## Read\n- The diff 81b10e749..7a1535fdf (internal/ and cmd/)\n- spec.md lines 100-139 and 284-295 (chunk table) and 380-470 (coverage map, edges)\n- gather.go:33-112\n- gather_inputs.go:90-93\n\n## Not verified\n- The skill files: I did not re-read bench-craft-review or bench-craft-tdd.\n- The grammar rejection of the extra flags (see Advice).\n- Gitlink (mode 160000) behavior.\n- The control-byte operand refusal.\n- Whether tip-current is red for every pair with tip ≠ HEAD.\n"
          },
          "axis": "Coverage",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "7a1535fdf85b124cd79d9106747269acadba764a",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": []
        },
        {
          "id": "ft375_ss-c1_standards_r2",
          "performer": "claude:ft375_ss-c1_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "410a4337ba0bc62d31ab5f62a7c55710837caaaf",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_standards_r2",
            "digest": "sha256:36ba8a6703b9342989ec7fd6a667fcfcead88633d597e6b25e59fa3a1f43e946",
            "excerpt": "Standards confirm: S1 closed, S2 closed, S3 closed; new: N1, N2\n\nEvidence bound: `--check-current` showed current=true at 040a2953 with no dirty record.\n\n- **S1 closed.** `ownership.SubjectPaths` (internal/preflight/ownership/ownership.go:217) reads paths through each collector's own decoder: decodeCensusSubject, ruleSubject, ownerSubject, decodeImportSubject, decodeAnchorSubject, decodeLinesSubject/Glob/Headroom, testSubject and reusedOwnerSubject. The local `subjectPath` is deleted from staleness/collect.go, and a subject that fails to decode becomes a diagnostic.\n- **S2 closed.** `git.ListTree` (internal/git/tree.go:151) now serves ReadTreeFile, TreeWithoutFile, the census `trackedAt` (sizing.go:291) and staleness `members` (collect.go:130).\n- **S3 closed.** The drift_test.go header no longer claims a red result.\n- **C4 code.** The fact now has one source. `Collect` parses both inventories from the committed spec bytes (collect.go:33), `gather` passes no manifest, and the doc comments are timeless.\n\nNew:\n- **N1** internal/git/tree.go:153. The comment says \"This function is the one reader of that grammar.\" That is false: other code still hand-parses full `ls-tree -z` records (internal/commitment/repository/light_path.go:121, internal/consumers/hunks.go:217, internal/roadmap/tree_validation.go:55). Rule: \"What survives merge is only what's true of the code as it stands\" (craft-comments). Disposition: concrete comment defect. Narrow the claim or delete it, or leave migrating the other callers for later. Confidence 7.\n- **N2** internal/preflight/ownership/source_test.go:585. TestSubjectPaths has a hand-written table of expected paths for each kind, and I found no recorded red for it in the commits, the review record or the spec. AGENTS.md: \"An independently authored test expectation counts as duplicated implementation knowledge, unless … someone must record and demonstrate that red.\" Disposition: either record a red (for example, drop the writer_probe path from `censusSubject.paths`) or get a reviewer ruling. Confidence 5.\n\nAdvice\n- The `subjectPaths` map lists premise kinds a second time next to `sizingPremises`, `symbolPremises` and the `premises` switch. The kinds match today, and a missing kind shows up as a diagnostic, not a silent gap.\n\nRead: AGENTS.md (in context); .agents/skills/bench-craft-comments/SKILL.md; the repair diff 7a1535fd..decb17ff; ownership/helpers.go:170-217; members.go:35-74; rg sweeps for `ls-tree`; commit messages in that range.\n\nNot verified: whether the tree_test.go named in plan commit 3b4220c exists (it is not in the diff; this is the Coverage axis); behaviour at run time (no tests run).\n"
          },
          "axis": "Standards",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "decb17ff468daefc64c1b0c64c4b19d89f480568",
          "finding_ids": [
            "N1",
            "N2"
          ],
          "supersedes": [
            "ft375_ss-c1_standards_r1"
          ]
        },
        {
          "id": "ft375_ss-c1_spec_r2",
          "performer": "claude:ft375_ss-c1_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "410a4337ba0bc62d31ab5f62a7c55710837caaaf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_spec_r2",
            "digest": "sha256:dd45024cf9b4b54bb984b98e02bc1f0b8e700c1a3d2467834694b33372a00f3a",
            "excerpt": "Spec confirm: P1 closed, C4 closed; new: none\n\nEvidence check: `bench preflight evidence ... --check-current` binds (`current=true`, `dirty_record=false`, pair 81b10e74..040a2953).\n\nP1: closed. At the given pair, `--audit-spec` prints only two `audit_diagnostics`: `cited-symbols` and `cited-tests`, each with \"no ... collector grades this audit\". It prints no `metadata` diagnostic. I parsed spec.md line 647 (the planning JSON) at the worktree tip and checked every location against its line:\n- All 8 SS-C1 seams (TestDriftObjects through TestAuditBound, lines 396–447) hold their needle exactly once.\n- So do all 21 SS-C1 decisions (SS1–12, SS35–40, SS51, SS52, SS56), at both `location` and `source`.\n- No other decision failed.\n\nC4: closed.\n- Spec line 79 gives the signature `Collect(root string, source ownership.Source, sourceFacts ownership.Snapshot) (AuditFacts, error)`. It matches internal/preflight/staleness/collect.go:24.\n- Line 80 (\"Collect parses both planning inventories from the committed spec bytes\") is met. collect.go:25-33 reads both sides through `git.ReadTreeFile` at Base and Tip and passes each through `ParseManifest`. gather.go:110 no longer passes checkout `content`.\n- Line 128 (\"Worktree content cannot replace committed source bytes\") is met. TestAuditCommittedInventory (inventory_test.go:93) checks it.\n- The new diagnostic for an unread premise subject (collect.go:838-841 in the delta, through `ownership.SubjectPaths` in ownership.go) stays within line 216 (\"An unsupported input is diagnostic and keeps manual review required\"). Decide sets `Complete = len(Diagnostics)==0` and `ManualRequired = !Complete` (staleness/decision.go:19,21). `fence-exclusion` is a valid manifest kind (manifest.go:147) that has no collector, so it lands as a diagnostic. It is not dropped.\n- No regression to SS1–12, SS35–40, SS51, SS52, or SS56. The drift_test.go changes only add cases: a quoted path for SS1–6, a nested reused-owner subject and a dropped-claim case for SS12, and a link refusal on both sides for SS51. No assertion was removed.\n\nAdvice:\n- The `subjectPaths` map (ownership.go) is a second kind-to-decoder registry next to the `premises` switch (helpers.go:182-209). This is a Standards-axis concern, not a spec defect.\n- SubjectPaths runs on premises that are not current, so they also produce diagnostics. That is conservative and allowed by the spec.\n- `git.ListTree` has no direct test in internal/git. tree_test.go is fenced but missing from the delta.\n\nRead: the repair delta (all files except reviews/), spec.md lines 74–143 and 216, collect.go:18-65, decision.go (rg), helpers.go:174-219, manifest.go:140-148, the audit-spec spill.\n\nNot verified: test execution; the reviews/ record; the SS-C1 needles at the base commit.\n"
          },
          "axis": "Spec",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "decb17ff468daefc64c1b0c64c4b19d89f480568",
          "finding_ids": [],
          "supersedes": [
            "ft375_ss-c1_spec_r1"
          ]
        },
        {
          "id": "ft375_ss-c1_coverage_r2",
          "performer": "claude:ft375_ss-c1_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "410a4337ba0bc62d31ab5f62a7c55710837caaaf",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_coverage_r2",
            "digest": "sha256:19bac5bf2fd587484e83b84f0ce0b5d5e2471b26124d29088cf0bd55a05def42",
            "excerpt": "Coverage confirm: C1 open (fence half); C2, C3, C4, X1 closed; new: none\n\nEvidence check: `--check-current` binds (current=true, dirty_record=false, source_tip 040a2953). I read the repair delta 7a1535fd..decb17ff once. `git status --short` is empty at exit.\n\n**C1, still open (fence half)**\n- Input: in `scope` at /home/mgibs/.bench/worktrees/bench-2826441890/b8e8080d78971d37a88341856b9407bb-bb87c9fc56cc0af9f6e84908b2d9a3a6/internal/preflight/staleness/collect.go:69, I changed the loop to `for _, content := range specs[1:] {`, which drops the baseline fence tokens.\n- Break: no test failed. That held for the narrow selection and for the full `./internal/preflight/...` plus `./internal/git` run.\n- Why: the new dropped-claim case in drift_test.go (around line 196–206) drops only the inventory claims. No staleness fixture changes a fence, so a baseline-only fence entry has no witness.\n- Spec line 112 requires every fence entry \"at both identities\", so this misses a binding requirement and is a blocking coverage gap. The inventory half is closed.\n- Disposition: needs a fresh repair session for the ticket that owns drift_test.go. Confidence 8.\n\n**Probes run** (every cell restored=yes)\n1. C1, baseline fence loop `specs[1:]`: silent.\n2. C1, baseline manifest replaced by the tip manifest: bit (TestDriftIdentity:206).\n3. C2, symlink exempted in `members`: bit (TestAuditSafety, \"after\" case).\n4. C2, special-file check made side-gated (`strict` flag, false for Base): bit (TestAuditSafety:343, \"before\" case).\n5. C3, ListTree splits at the last tab: bit (TestDriftObjects:122).\n6. C3, ListTree rewrites non-ASCII path bytes: bit (TestDriftObjects:153).\n7. C4, tip manifest parsed from the checkout file through `os.ReadFile`: bit (TestAuditCommittedInventory:47).\n8. X1, `censusSubject.paths` drops the writer-probe path: bit (TestSubjectPaths, posture).\n9. X1, `censusSubject.paths` drops the listed members: bit (TestSubjectPaths, retirement).\n\nA first attempt at probe 9 used `--omit`. It did not compile, so it counts as invalid and is not a result.\n\n**Advice**\nThe tip fence tokens come from committed bytes, but no test shows that. A checkout spec with a different fence would not turn TestAuditCommittedInventory red. C4's charge names only claim sources, so this is optional.\n\n**Not verified**\n- Whether the author's C1 probe dropped fences and inventory together, which would hide the fence gap.\n- The gate and the system suite.\n- Packages outside `./internal/preflight/...` and `./internal/git`.\n"
          },
          "axis": "Coverage",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "decb17ff468daefc64c1b0c64c4b19d89f480568",
          "finding_ids": [
            "C1"
          ],
          "supersedes": [
            "ft375_ss-c1_coverage_r1"
          ]
        },
        {
          "id": "ft375_ss-c1_standards_r3",
          "performer": "claude:ft375_ss-c1_standards_r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_standards_r3",
            "digest": "sha256:0b59dd7c6e983efccab9e0c1c2fa6c9606a8077ed73b9362dcba2749818c1690",
            "excerpt": "Standards confirm 2: N1 closed; new: none\n\nEvidence sha256:4d183986… binds to assignment bb87c9fc… at source tip a0eb0685 (current=true, dirty_record=false).\n\nAdvice: N1 is closed. The false sentence is gone from internal/git/tree.go:151-153. The comment that remains is true of the code: lines 159-170 split on NUL, refuse an unterminated record, and cut at the first tab. The new TestDriftIdentity case (internal/preflight/staleness/drift_test.go:209-218) builds its expected rows through the existing `modified` closure (line 188). It adds no second copy of the drift-row grammar. Its comment states the why in timeless present and carries no provenance. Line 172 still spells that row grammar out by hand, separate from the closure. That line is outside this delta, so I am not grading it here.\n\nNot verified: whether the new case goes red when the dropped-fence handling is removed (I ran no probe), and whether a red was recorded for it. I did not list the other ls-tree parsers.\n"
          },
          "axis": "Standards",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "c56e4a7f92e37107440b0e3fd4c092f9ed22188c",
          "finding_ids": [],
          "supersedes": [
            "ft375_ss-c1_standards_r2"
          ]
        },
        {
          "id": "ft375_ss-c1_spec_r3",
          "performer": "claude:ft375_ss-c1_spec_r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_spec_r3",
            "digest": "sha256:20b60d597c87bc6919a4e6b6812a094f743b0ab2fea787d3f43fc97df806ff92",
            "excerpt": "Spec confirm 2: pass; new: none\n\nThe evidence check binds: current=true, dirty_record=false, source_tip a0eb0685. The only commit after c56e4a7f is the record commit a0eb0685.\n\n- \"Enumerate every exact fence entry and every tracked member of a fence prefix at both identities.\" (spec.md:112): met. The new case in drift_test.go:209-218 commits a baseline fence of `outside/fenced/` through `SpecBody(slug, extra fence line)` (preflighttest/fixture.go:121). The tip drops that fence and changes `outside/fenced/member.txt`. The case asserts that the drift rows are exactly the spec and that member. Disposition: closes C1 (fence half). Confidence 8.\n- No SS-C1 row regressed (SS1–12, SS35–40, SS51, SS52, SS56): the delta only adds that test block and trims one comment at internal/git/tree.go:153. It removes no assertion and changes no production logic. Confidence 8.\n\nAdvice: none.\n\nNot verified: no tests or probes were run. I did not confirm a recorded red for the baseline-fence loop mutation (`specs[1:]`). I did not read the a0eb0685 record contents.\n"
          },
          "axis": "Spec",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "c56e4a7f92e37107440b0e3fd4c092f9ed22188c",
          "finding_ids": [],
          "supersedes": [
            "ft375_ss-c1_spec_r2"
          ]
        },
        {
          "id": "ft375_ss-c1_coverage_r3",
          "performer": "claude:ft375_ss-c1_coverage_r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_coverage_r3",
            "digest": "sha256:0ff19e25e9df8e8ba1eee59d791a0b3e3c200b374956cd65006ac5c9ded42e43",
            "excerpt": "Coverage confirm 2: C1 open; new: C1b\n\nThe evidence binds: `--check-current` returned current=true and dirty_record=false at source_tip a0eb0685. After the probe, `git status --short` is empty.\n\n**Probe**\n- Site: internal/preflight/staleness/collect.go:71. On the baseline pass only (`len(entries) == 2`), the mutation deletes every exact fence token, meaning any token that does not end in `/`. Prefix tokens and the tip tokens stay.\n- Selection: `--package ./internal/preflight/... --run 'TestDrift|TestAudit'`.\n- Result: silent, cause passed, 0 failures, restored yes.\n\n**Finding**\n- C1b. Input: an exact (non-prefix) fence entry that only the baseline spec declares. Break: the new drift_test.go:209-219 case covers only a dropped prefix (`outside/fenced/`), so nothing catches a dropped exact entry. This misses spec line 112 (\"every exact fence entry … at both identities\"). Disposition: blocking. It goes to a fresh repair session for ticket 01. Confidence 7.\n\n**Not verified**\n- Whether the tip side's exact entries have a witness.\n- The full `./internal/preflight/...` run without the `--run` filter.\n- The gate.\n"
          },
          "axis": "Coverage",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "c56e4a7f92e37107440b0e3fd4c092f9ed22188c",
          "finding_ids": [
            "C1b"
          ],
          "supersedes": [
            "ft375_ss-c1_coverage_r2"
          ]
        },
        {
          "id": "ft375_ss-c1_standards_r4",
          "performer": "claude:ft375_ss-c1_standards_r4",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_standards_r4",
            "digest": "sha256:28e2817af5d5b7a61c6c169976dc9df6f9500140c7bd0a56ffff06508088bea1",
            "excerpt": "Standards confirm 3: pass; new: none\n\nEvidence sha256:bcd54a0b… binds. It reports current=true, dirty_record=false, source_tip 80d77348.\n\nFindings: none.\n- internal/preflight/staleness/drift_test.go:209-210. The comment rule in bench-craft-comments allows this comment. It states timelessly why only the baseline fence names the files. It has no narration and no provenance. Disposition: none. Confidence: high.\n- internal/preflight/staleness/drift_test.go:211-222. The one-source rule allows this code. The expectation `fenced` comes from the same fixture constants (`spec`, `member`, `exact`) through `modified`. It does not copy any implementation knowledge, and it makes no new fixture harness. Disposition: none. Confidence: high.\n\nNot verified: I ran no tests and no probes. The red behavior of the new exact-fence case is not shown. I did not read the drift implementation, so I cannot confirm that the expected row order matches `sorted`.\n"
          },
          "axis": "Standards",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "c585f9d08d9f7d320d5e56dfc88c69a09f43733c",
          "finding_ids": [],
          "supersedes": [
            "ft375_ss-c1_standards_r3"
          ]
        },
        {
          "id": "ft375_ss-c1_spec_r4",
          "performer": "claude:ft375_ss-c1_spec_r4",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_spec_r4",
            "digest": "sha256:2ac76542744d48e103490a991b5ecae66b4cf03818699dfc9baa1a15e6b3e9ff",
            "excerpt": "Spec confirm 3: pass; new: none\n\nThe evidence binds: current=true, dirty_record=false, source_tip 80d77348.\n\n- spec.md:112 \"Enumerate every exact fence entry and every tracked member of a fence prefix at both identities.\" Met. The new case at internal/preflight/staleness/drift_test.go:209-228 commits a baseline fence with the prefix `outside/fenced/` and the exact entry `outside/exact.txt`, using SpecBody's variadic extra fence lines (preflighttest/fixture.go:121). The tip drops both entries and changes both files. The case asserts that the drift rows are exactly the spec, the member, and the exact file. This closes C1b. Confidence 8.\n- No production change. The delta touches only drift_test.go, reviews/preflight-spec-staleness.md, and the spec's completion-plan JSON. The test edit only widens the existing fence case, and no assertion in SS1–12, SS35–40, SS51, SS52, or SS56 was removed. Confidence 8.\n\nNot verified: I ran no tests or probes and saw no recorded red for the exact-token mutation. I did not check commits c585f9d..80d77348.\n"
          },
          "axis": "Spec",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "c585f9d08d9f7d320d5e56dfc88c69a09f43733c",
          "finding_ids": [],
          "supersedes": [
            "ft375_ss-c1_spec_r3"
          ]
        },
        {
          "id": "ft375_ss-c1_coverage_r4",
          "performer": "claude:ft375_ss-c1_coverage_r4",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38a4ca7901fc67f8466ddd176dbeecf8cc084972",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c1_coverage_r4",
            "digest": "sha256:4178bfcb78946ab3df4e1547b15335a89f259f8ab91b4ac726c073fd84217929",
            "excerpt": "Coverage confirm 3: C1b closed; new: none\n\nThe evidence check binds (current=true, dirty_record=false, source_tip 80d77348). After the probe, `git status --short` is empty.\n\n**Probe:** bit, restored yes.\n- Site: internal/preflight/staleness/collect.go:136, in `members`. The mutation honours an exact match only for the spec path (`&& (member.Path != entry || entry == entries[0])`).\n- Run: `--package ./internal/preflight/... --run 'TestDrift|TestAudit'`. TestDriftIdentity failed.\n\n**Findings:** none blocking.\n- Advice A1. Input: the probe above. Break: the first failure is the claim-drift assertion at drift_test.go:193, not the new exact-fence case at drift_test.go:209-226. The probe therefore does not show that the new case bites on its own. With the two earlier bites, I judge C1b closed. Disposition: optional advice under the bounded repair policy. Confidence 6.\n\n**Not verified:** a probe aimed only at exact fence tokens through `tickets.WritesPath`, the full package run, and the gate.\n"
          },
          "axis": "Coverage",
          "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
          "tip": "c585f9d08d9f7d320d5e56dfc88c69a09f43733c",
          "finding_ids": [],
          "supersedes": [
            "ft375_ss-c1_coverage_r3"
          ]
        }
      ]
    },
    {
      "id": "SS-C2",
      "base": "c585f9d08d9f7d320d5e56dfc88c69a09f43733c",
      "tip": "f791668126e10cb8ed959664b58c0ac0bc7cf953",
      "plan_digest": "sha256:4be92c9f220cff6e03206483328acc81fd3de1894913f7c9c4df48caa0d0a87b",
      "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
      "acceptance_rows": [
        "SS13",
        "SS14",
        "SS15",
        "SS16",
        "SS17",
        "SS18",
        "SS19",
        "SS20",
        "SS21",
        "SS22",
        "SS23",
        "SS24",
        "SS25",
        "SS26",
        "SS27",
        "SS28",
        "SS29",
        "SS30",
        "SS31",
        "SS32",
        "SS33",
        "SS34",
        "SS49",
        "SS50",
        "SS57",
        "SS58"
      ],
      "verification": [
        {
          "id": "ft375_t04-ss-04-focused",
          "performer": "claude:ft375_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04",
            "digest": "sha256:6ae96e18621d793ac8a37ff45662a4baaefb2771d3f6cfb5dae824b9bef19c45",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestAuditCompleteness/SS50_present_empty,17,1,0,null,yes\n  2,bit,swap,failed,1,TestAuditCompleteness/SS33_partial,17,1,0,null,yes\n  3,bit,swap,failed,1,TestAuditCompleteness/SS57_SS58_pending_native,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestCoverageIdentity|TestAuditCompleteness,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,3173\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,836,17\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestAuditCompleteness/SS50_present_empty,\"completeness_test.go:215: diagnostics = [\\\"metadata|row SS28 seam omits internal/example/kept_test.go (TestKept) from seams\\\" \\\"metadata|row SS28 seam omits internal/example/kept_test.go (TestOther) from seams\\\"], want the prefixes [\\\"metadata|r… (438 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,884,17\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestAuditCompleteness/SS33_partial,\"completeness_test.go:215: diagnostics = [], want the prefixes [\\\"bench-planning-facts|the planning inventory is not complete: metadata: the premise inventory is partial\\\"]\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,847,17\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestAuditCompleteness/SS57_SS58_pending_native,\"completeness_test.go:215: diagnostics = [], want the prefixes [\\\"native-SS47|native-future manual claim\\\" \\\"native-SS48|native-future manual claim\\\" \\\"native-SS53|native-future manual claim\\\" \\\"native-SS54|native-future manual claim\\\" \\\"native-SS55|… (268 bytes)\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestCoverageIdentity|TestAuditCompleteness'",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-04-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t04",
              "digest": "sha256:6ae96e18621d793ac8a37ff45662a4baaefb2771d3f6cfb5dae824b9bef19c45",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestAuditCompleteness/SS50_present_empty,17,1,0,null,yes\n  2,bit,swap,failed,1,TestAuditCompleteness/SS33_partial,17,1,0,null,yes\n  3,bit,swap,failed,1,TestAuditCompleteness/SS57_SS58_pending_native,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestCoverageIdentity|TestAuditCompleteness,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,3173\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,836,17\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestAuditCompleteness/SS50_present_empty,\"completeness_test.go:215: diagnostics = [\\\"metadata|row SS28 seam omits internal/example/kept_test.go (TestKept) from seams\\\" \\\"metadata|row SS28 seam omits internal/example/kept_test.go (TestOther) from seams\\\"], want the prefixes [\\\"metadata|r… (438 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,884,17\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestAuditCompleteness/SS33_partial,\"completeness_test.go:215: diagnostics = [], want the prefixes [\\\"bench-planning-facts|the planning inventory is not complete: metadata: the premise inventory is partial\\\"]\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,847,17\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestAuditCompleteness/SS57_SS58_pending_native,\"completeness_test.go:215: diagnostics = [], want the prefixes [\\\"native-SS47|native-future manual claim\\\" \\\"native-SS48|native-future manual claim\\\" \\\"native-SS53|native-future manual claim\\\" \\\"native-SS54|native-future manual claim\\\" \\\"native-SS55|… (268 bytes)\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t04-ss-04-check-1",
          "performer": "claude:ft375_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04",
            "digest": "sha256:cc37f23a135fc9c9955cd6ad53b31a15976885f76c9f382e9e55e5db2d152e47",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,38869,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23670,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,40862,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,3185,52\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t04-ss-04-check-2",
          "performer": "claude:ft375_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04",
            "digest": "sha256:0582558b204b51d9e78d1b565a5965c4f519069c2fbaf345250065ee9f2b40ad",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,pass,1976,145\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-2",
          "command": "bench test --package ./internal/coverage",
          "exit_code": 0
        },
        {
          "id": "ft375_t04-ss-04-check-3",
          "performer": "claude:ft375_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04",
            "digest": "sha256:5b0b71dc6dc1a15b4cb5cb39c1c55017bcf50eda975df4c354428d09c92781ba",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4263,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-3",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft375_t04-ss-04-check-4",
          "performer": "claude:ft375_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04",
            "digest": "sha256:301696dc3cd894eb5e64634e001f59b6de0e7f205ada8cc506c3f22b616cd2a0",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14733,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t04-ss-04-check-5",
          "performer": "claude:ft375_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04",
            "digest": "sha256:917eaa4e5fad5a51eb664bbabf460c7b494c0d70e8ce8f5a1cccf9d54ec38331",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14440,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-5",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t02-ss-02-focused",
          "performer": "claude:ft375_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02",
            "digest": "sha256:f57e76ad630ab3841dc6665650bd7d45a0f116370b0446b823288e0a307be52c",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,3,TestSymbolIdentity,11,1,0,null,yes\n  2,bit,omit,failed,1,TestSymbolIdentity/SS16_rule,11,1,0,null,yes\n  3,bit,omit,failed,2,TestSymbolIdentity,11,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestSymbolIdentity,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,4239\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,280,11\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity,\"identity_test.go:108: symbol drift rows = 4, want 6: map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/example/sym.go|absent|symbol example.Gone absent|func Gone()|absent:true internal/example/… (425 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS14_change,\"identity_test.go:99: drift lacks \\\"internal/example/sym.go|internal/example/sym.go|symbol example.Change signature|func Change(a int) int|func Change(a string) int\\\"; symbol rows map[absent|internal/example/sym.go|symbol example.Fresh absent|… (549 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS14_receiver,\"identity_test.go:99: drift lacks \\\"internal/example/sym.go|internal/example/sym.go|symbol example.Origin.Method signature|func (Origin) Method()|func (*Origin) Method()\\\"; symbol rows map[absent|internal/example/sym.go|symbol example.Fresh ab… (554 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,280,11\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS16_rule,\"identity_test.go:99: drift lacks \\\"internal/example/sym.go|internal/example/sym.go|symbol example.Rule kind|const|var Rule = 1\\\"; symbol rows map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/exa… (790 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,290,11\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity,\"identity_test.go:108: symbol drift rows = 5, want 6: map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/example/sym.go|absent|symbol example.Gone absent|func Gone()|absent:true internal/example/… (593 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS15_move,\"identity_test.go:99: drift lacks \\\"internal/example/sym.go|internal/example/moved.go|symbol example.Move moved|func Move()|func Move()\\\"; symbol rows map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true inte… (688 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-focused",
          "command": "bench test --package ./internal/preflight/... --run TestSymbolIdentity",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-02-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t02",
              "digest": "sha256:f57e76ad630ab3841dc6665650bd7d45a0f116370b0446b823288e0a307be52c",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,3,TestSymbolIdentity,11,1,0,null,yes\n  2,bit,omit,failed,1,TestSymbolIdentity/SS16_rule,11,1,0,null,yes\n  3,bit,omit,failed,2,TestSymbolIdentity,11,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestSymbolIdentity,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,4239\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,280,11\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity,\"identity_test.go:108: symbol drift rows = 4, want 6: map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/example/sym.go|absent|symbol example.Gone absent|func Gone()|absent:true internal/example/… (425 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS14_change,\"identity_test.go:99: drift lacks \\\"internal/example/sym.go|internal/example/sym.go|symbol example.Change signature|func Change(a int) int|func Change(a string) int\\\"; symbol rows map[absent|internal/example/sym.go|symbol example.Fresh absent|… (549 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS14_receiver,\"identity_test.go:99: drift lacks \\\"internal/example/sym.go|internal/example/sym.go|symbol example.Origin.Method signature|func (Origin) Method()|func (*Origin) Method()\\\"; symbol rows map[absent|internal/example/sym.go|symbol example.Fresh ab… (554 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,280,11\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS16_rule,\"identity_test.go:99: drift lacks \\\"internal/example/sym.go|internal/example/sym.go|symbol example.Rule kind|const|var Rule = 1\\\"; symbol rows map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/exa… (790 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,290,11\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity,\"identity_test.go:108: symbol drift rows = 5, want 6: map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/example/sym.go|absent|symbol example.Gone absent|func Gone()|absent:true internal/example/… (593 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS15_move,\"identity_test.go:99: drift lacks \\\"internal/example/sym.go|internal/example/moved.go|symbol example.Move moved|func Move()|func Move()\\\"; symbol rows map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true inte… (688 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t02-ss-02-check-1",
          "performer": "claude:ft375_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02",
            "digest": "sha256:21adb2577fe814b45dc2669906f0822e49b1cc731a9867ed0aea1cc1c0c4885f",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,38956,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23544,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,40523,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,3161,52\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t02-ss-02-check-2",
          "performer": "claude:ft375_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02",
            "digest": "sha256:b81b9c0cb45cd6430c7fa11dd1b27bf2cd09372a60a700b46f37d81ead91df4f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,pass,1809,66\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-2",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0
        },
        {
          "id": "ft375_t02-ss-02-check-3",
          "performer": "claude:ft375_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02",
            "digest": "sha256:2ecaa512be88e585d2d6867ee760cf91c2cfaec30542472817ecbdd18bd3e08e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14462,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t02-ss-02-check-4",
          "performer": "claude:ft375_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02",
            "digest": "sha256:b14439cda4ff428acef8d39536509aa9108d52afdeebc5828a7ff483ff3d9a7f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14135,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t03-ss-03-focused",
          "performer": "claude:ft375_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03",
            "digest": "sha256:84b74b9af50028d448df983718ff876edf87c93b9c8b534d53356703c09ecb5d",
            "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestTestIdentity,11,1,0,null,yes\n  2,bit,swap,failed,1,TestTestIdentity,11,1,0,null,yes\n  3,bit,swap,failed,1,TestTestIdentity/SS20_gone,11,1,0,null,yes\n  4,bit,swap,failed,2,TestTestIdentity,11,1,0,null,yes\n  5,bit,swap,failed,2,TestTestIdentity,11,1,0,null,yes\n  6,bit,swap,failed,1,TestTestIdentity/SS27_tagged,11,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestTestIdentity,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,7820\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,199,11\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:203: test drift rows = 2, want 5: map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59:true int… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS21_removed,\"identity_test.go:194: drift lacks \\\"internal/example/removed_test.go|internal/example/removed_test.go|test TestRemoved function-absent|7a1c971cb5216bc18a3548824c46a740c6c0e22f|d8fc4bab7558b887d077b760389d8da02647f2a5\\\"; test rows map[internal… (530 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS25_nosub,\"identity_test.go:194: drift lacks \\\"internal/example/kept_test.go|internal/example/kept_test.go|test TestSub/case two function-absent|8b7fd3111d935d68eaec647261034b4ba7727c50|8b7fd3111d935d68eaec647261034b4ba7727c50\\\"; test rows map[internal/… (529 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS27_untagged,\"identity_test.go:194: drift lacks \\\"internal/example/untagged_test.go|internal/example/untagged_test.go|test TestUntagged function-absent|3f3b86a034e2ffa9aa8221047f7eb6cb08eac515|d8fc4bab7558b887d077b760389d8da02647f2a5\\\"; test rows map[inter… (533 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,179,11\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:203: test drift rows = 6, want 5: map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59:true int… (1020 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,181,11\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS20_gone,\"identity_test.go:194: drift lacks \\\"internal/example/gone_test.go|absent|test TestGone file-absent|46ba373cb6b4a058a49ea7d864dbb2b8cb9b8ba3|absent\\\"; test rows map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly… (1020 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,189,11\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:203: test drift rows = 4, want 5: map[internal/example/gone_test.go|absent|test TestGone file-absent|46ba373cb6b4a058a49ea7d864dbb2b8cb9b8ba3|absent:true internal/example/kept_test.go|internal/example/kept_test.go|test Test… (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS23_early,\"identity_test.go:194: drift lacks \\\"internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59\\\"; test rows map[internal/exa… (904 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,184,11\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:203: test drift rows = 6, want 5: map[internal/example/computed_test.go|internal/example/computed_test.go|test TestComputed/x function-absent|c3bbec34e0c6c5ab6b87591cc3855bdf1141ec1e|c3bbec34e0c6c5ab6b87591cc3855bdf1141ec1e… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS26_computed,\"identity_test.go:197: diagnostic of computed = \\\"\\\", want prefix \\\"cited test TestComputed/x is unsupported\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  6\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,182,11\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS27_tagged,\"identity_test.go:197: diagnostic of tagged = \\\"\\\", want prefix \\\"cited test TestTagged is unsupported-context\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-03-focused",
          "command": "bench test --package ./internal/preflight/... --run TestTestIdentity",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-03-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t03",
              "digest": "sha256:84b74b9af50028d448df983718ff876edf87c93b9c8b534d53356703c09ecb5d",
              "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestTestIdentity,11,1,0,null,yes\n  2,bit,swap,failed,1,TestTestIdentity,11,1,0,null,yes\n  3,bit,swap,failed,1,TestTestIdentity/SS20_gone,11,1,0,null,yes\n  4,bit,swap,failed,2,TestTestIdentity,11,1,0,null,yes\n  5,bit,swap,failed,2,TestTestIdentity,11,1,0,null,yes\n  6,bit,swap,failed,1,TestTestIdentity/SS27_tagged,11,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestTestIdentity,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,7820\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,199,11\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:203: test drift rows = 2, want 5: map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59:true int… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS21_removed,\"identity_test.go:194: drift lacks \\\"internal/example/removed_test.go|internal/example/removed_test.go|test TestRemoved function-absent|7a1c971cb5216bc18a3548824c46a740c6c0e22f|d8fc4bab7558b887d077b760389d8da02647f2a5\\\"; test rows map[internal… (530 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS25_nosub,\"identity_test.go:194: drift lacks \\\"internal/example/kept_test.go|internal/example/kept_test.go|test TestSub/case two function-absent|8b7fd3111d935d68eaec647261034b4ba7727c50|8b7fd3111d935d68eaec647261034b4ba7727c50\\\"; test rows map[internal/… (529 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS27_untagged,\"identity_test.go:194: drift lacks \\\"internal/example/untagged_test.go|internal/example/untagged_test.go|test TestUntagged function-absent|3f3b86a034e2ffa9aa8221047f7eb6cb08eac515|d8fc4bab7558b887d077b760389d8da02647f2a5\\\"; test rows map[inter… (533 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,179,11\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:203: test drift rows = 6, want 5: map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59:true int… (1020 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,181,11\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS20_gone,\"identity_test.go:194: drift lacks \\\"internal/example/gone_test.go|absent|test TestGone file-absent|46ba373cb6b4a058a49ea7d864dbb2b8cb9b8ba3|absent\\\"; test rows map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly… (1020 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,189,11\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:203: test drift rows = 4, want 5: map[internal/example/gone_test.go|absent|test TestGone file-absent|46ba373cb6b4a058a49ea7d864dbb2b8cb9b8ba3|absent:true internal/example/kept_test.go|internal/example/kept_test.go|test Test… (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS23_early,\"identity_test.go:194: drift lacks \\\"internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59\\\"; test rows map[internal/exa… (904 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,184,11\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:203: test drift rows = 6, want 5: map[internal/example/computed_test.go|internal/example/computed_test.go|test TestComputed/x function-absent|c3bbec34e0c6c5ab6b87591cc3855bdf1141ec1e|c3bbec34e0c6c5ab6b87591cc3855bdf1141ec1e… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS26_computed,\"identity_test.go:197: diagnostic of computed = \\\"\\\", want prefix \\\"cited test TestComputed/x is unsupported\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  6\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,182,11\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS27_tagged,\"identity_test.go:197: diagnostic of tagged = \\\"\\\", want prefix \\\"cited test TestTagged is unsupported-context\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t03-ss-03-check-1",
          "performer": "claude:ft375_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03",
            "digest": "sha256:864828cf15400415752304d674e7e5331005722ee0a8719f4947b8370cfa310d",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,39121,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23590,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,40836,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,3157,52\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-03-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t03-ss-03-check-2",
          "performer": "claude:ft375_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03",
            "digest": "sha256:c09fa265683ec84fa9070fafe1679097c125ce6da6319edfaa9f25620696904c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,pass,1651,145\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-03-check-2",
          "command": "bench test --package ./internal/coverage",
          "exit_code": 0
        },
        {
          "id": "ft375_t03-ss-03-check-3",
          "performer": "claude:ft375_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03",
            "digest": "sha256:f7c91009848439419f3c051f86f1f273cd2f715879407450db0e8a27a8e613db",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14607,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-03-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t03-ss-03-check-4",
          "performer": "claude:ft375_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03",
            "digest": "sha256:c124b4456a968dd6931e524b975cda4b19476ec9fba1253abc7f29f113b1c5ce",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14496,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-03-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t03-ss-03-check-5",
          "performer": "claude:ft375_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03",
            "digest": "sha256:15e0410cba15a6e24e79f5d44e6484993adae009088c58a7b7391c7daf4cf8cc",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,33893,1118\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ss-03-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t04_rep1-ss-04-focused",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:7eadccba54607c444ddc0146bc750929a5e77ef3d29b71f008309a86d7065d4a",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCoverageIdentity,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestCoverageIdentity|TestAuditCompleteness,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,920\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,863,17\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestCoverageIdentity,\"completeness_test.go:89: oversized baseline ticket = (1), want the ticket source refusal and no staleness result:\",19\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestCoverageIdentity|TestAuditCompleteness'",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-04-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t04_rep1",
              "digest": "sha256:7eadccba54607c444ddc0146bc750929a5e77ef3d29b71f008309a86d7065d4a",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCoverageIdentity,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestCoverageIdentity|TestAuditCompleteness,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,920\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,863,17\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestCoverageIdentity,\"completeness_test.go:89: oversized baseline ticket = (1), want the ticket source refusal and no staleness result:\",19\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t04_rep1-ss-04-check-1",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:7889b5a6be37d53277b33625fb0f80dc9bef9f5a058a083b192ee2fbbca3c9c2",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,37668,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,22754,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,39615,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,3938,56\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t04_rep1-ss-04-check-2",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:4bab2b4f882f93bdc174f6d7e63253ce77ee5b6cd28d6912c3216c37d6107bd1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,pass,1666,145\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-2",
          "command": "bench test --package ./internal/coverage",
          "exit_code": 0
        },
        {
          "id": "ft375_t04_rep1-ss-04-check-3",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:5c4b60190638055dbdc8257d2173012ee6ab70c3c11a9fa34b1fb9d24fccbe1e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3795,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-3",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft375_t04_rep1-ss-04-check-4",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:abac8366892845f9a16cc291d3e7dcd05234fcb8d7c37ab378c3234edfa3e035",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13841,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t04_rep1-ss-04-check-5",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:46f10cce6a3973ef32d48e1f3c1ebbeb031b7b26857ed940c6c719192da78de1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,13552,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-5",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t02_rep1-ss-02-focused",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:39d858808496703c198777ca02fce4c06b0053771a25a90762bcaf77ea187d7e",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestSymbolIdentity/SS17_bare,14,1,0,null,yes\n  2,bit,omit,failed,1,TestSymbolIdentity/SS17_bare,14,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestSymbolIdentity,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1786\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,698,14\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS17_bare,\"identity_test.go:128: premise bare = \\\"unknown\\\", want \\\"green\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,676,14\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS17_bare,\"identity_test.go:134: diagnostic of bare = \\\"\\\", want prefix \\\"cited symbol Dup is ambiguous\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-focused",
          "command": "bench test --package ./internal/preflight/... --run TestSymbolIdentity",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-02-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t02_rep1",
              "digest": "sha256:39d858808496703c198777ca02fce4c06b0053771a25a90762bcaf77ea187d7e",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestSymbolIdentity/SS17_bare,14,1,0,null,yes\n  2,bit,omit,failed,1,TestSymbolIdentity/SS17_bare,14,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestSymbolIdentity,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1786\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,698,14\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS17_bare,\"identity_test.go:128: premise bare = \\\"unknown\\\", want \\\"green\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,676,14\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS17_bare,\"identity_test.go:134: diagnostic of bare = \\\"\\\", want prefix \\\"cited symbol Dup is ambiguous\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t02_rep1-ss-02-check-1",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:eb5b1ed0cb615b5c009e6e26022850807bf19f4034709e5246eebfad5574780a",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,37525,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,22848,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,39556,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,3983,56\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t02_rep1-ss-02-check-2",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:49c9687a7c514859dd9dd4e49e4e44eec6266f2bfa91d51e90d6e68cadc10f05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,pass,1752,66\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-2",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0
        },
        {
          "id": "ft375_t02_rep1-ss-02-check-3",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:8932364bcb15ab4a5326ca5c00c74074925b3e24b1f9e4fa4a56b6845791987c",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13958,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t02_rep1-ss-02-check-4",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:e6718df9616017763f6f744fa662b76e54420821e2c157148b1e80951c0f7187",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14014,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t02_rep1-ss-02-check-5",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:c61268d8a563dd811bbb54e5f6963cfa04e76ffe545d848eedefde4ba5b35795",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,33344,1118\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ss-02-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t02_rep1b-ss-02-focused",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:f290ab2c4c9242b8971acb2e2752b674160bca085c297882ff4642e95228d815",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,3,TestSymbolIdentity,14,1,0,null,yes\n  2,bit,omit,failed,1,TestSymbolIdentity/SS16_rule,14,1,0,null,yes\n  3,bit,omit,failed,2,TestSymbolIdentity,14,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestSymbolIdentity,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,3666\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,797,14\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity,\"identity_test.go:140: symbol drift rows = 4, want 6: map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/example/sym.go|absent|symbol example.Gone absent|func Gone()|absent:true internal/example/… (425 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS14_change,\"identity_test.go:128: premise change = \\\"green\\\", want \\\"unknown\\\"\",2\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS14_receiver,\"identity_test.go:128: premise receiver = \\\"green\\\", want \\\"unknown\\\"\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,768,14\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS16_rule,\"identity_test.go:131: drift lacks \\\"internal/example/sym.go|internal/example/sym.go|symbol example.Rule kind|const|var Rule = 1\\\"; symbol rows map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/ex… (791 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,778,14\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity,\"identity_test.go:140: symbol drift rows = 5, want 6: map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/example/sym.go|absent|symbol example.Gone absent|func Gone()|absent:true internal/example/… (593 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS15_move,\"identity_test.go:128: premise move = \\\"green\\\", want \\\"unknown\\\"\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-focused",
          "command": "bench test --package ./internal/preflight/... --run TestSymbolIdentity",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-02-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t02_rep1",
              "digest": "sha256:f290ab2c4c9242b8971acb2e2752b674160bca085c297882ff4642e95228d815",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,3,TestSymbolIdentity,14,1,0,null,yes\n  2,bit,omit,failed,1,TestSymbolIdentity/SS16_rule,14,1,0,null,yes\n  3,bit,omit,failed,2,TestSymbolIdentity,14,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestSymbolIdentity,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,3666\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,797,14\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity,\"identity_test.go:140: symbol drift rows = 4, want 6: map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/example/sym.go|absent|symbol example.Gone absent|func Gone()|absent:true internal/example/… (425 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS14_change,\"identity_test.go:128: premise change = \\\"green\\\", want \\\"unknown\\\"\",2\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS14_receiver,\"identity_test.go:128: premise receiver = \\\"green\\\", want \\\"unknown\\\"\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,768,14\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS16_rule,\"identity_test.go:131: drift lacks \\\"internal/example/sym.go|internal/example/sym.go|symbol example.Rule kind|const|var Rule = 1\\\"; symbol rows map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/ex… (791 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,778,14\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity,\"identity_test.go:140: symbol drift rows = 5, want 6: map[absent|internal/example/sym.go|symbol example.Fresh absent|absent|func Fresh():true internal/example/sym.go|absent|symbol example.Gone absent|func Gone()|absent:true internal/example/… (593 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestSymbolIdentity/SS15_move,\"identity_test.go:128: premise move = \\\"green\\\", want \\\"unknown\\\"\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t02_rep1b-ss-02-check-1",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:1369a992d743c059c48857e2a2069100f48c31f84dfd8d66485d58e96d81ee53",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,45266,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,26915,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,47262,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,4124,56\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t02_rep1b-ss-02-check-2",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:8e89216d5063f2ed094629292986c0ab4dbf0ab9e90bbc77c321126ca6b7556a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,pass,1991,66\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-2",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0
        },
        {
          "id": "ft375_t02_rep1b-ss-02-check-3",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:1d400bf1d4c4ff526de2b610d1a9182f8d8a1b359a6dabb81a65d5f9de02e798",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18414,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t02_rep1b-ss-02-check-4",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:0fe2628bd3462f2721e9be09599155894b0f462605caf857b7efd068c03f0c55",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18172,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-02-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t02_rep1b-ss-02-check-5",
          "performer": "claude:ft375_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t02_rep1",
            "digest": "sha256:f18c49dbbf5616e993794e35a01a5fd4d1285911118ced71096e117ed3f384ec",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,41677,1118\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ss-02-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t03_rep1-ss-03-focused",
          "performer": "claude:ft375_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03_rep1",
            "digest": "sha256:0541a8242d04dda68c9a6f1ed9d35e1917c8f3b539db4d1991415b1acc298c65",
            "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestTestIdentity,12,1,0,null,yes\n  2,bit,swap,failed,1,TestTestIdentity,12,1,0,null,yes\n  3,bit,swap,failed,1,TestTestIdentity/SS20_gone,12,1,0,null,yes\n  4,bit,swap,failed,2,TestTestIdentity,12,1,0,null,yes\n  5,bit,swap,failed,3,TestTestIdentity,12,1,0,null,yes\n  6,bit,swap,failed,1,TestTestIdentity/SS27_tagged,12,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestTestIdentity,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,8022\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,334,12\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:253: test drift rows = 2, want 5: map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59:true int… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS21_removed,\"identity_test.go:244: drift lacks \\\"internal/example/removed_test.go|internal/example/removed_test.go|test TestRemoved function-absent|7a1c971cb5216bc18a3548824c46a740c6c0e22f|d8fc4bab7558b887d077b760389d8da02647f2a5\\\"; test rows map[internal… (530 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS25_nosub,\"identity_test.go:244: drift lacks \\\"internal/example/kept_test.go|internal/example/kept_test.go|test TestSub/case two function-absent|8b7fd3111d935d68eaec647261034b4ba7727c50|8b7fd3111d935d68eaec647261034b4ba7727c50\\\"; test rows map[internal/… (529 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS27_untagged,\"identity_test.go:244: drift lacks \\\"internal/example/untagged_test.go|internal/example/untagged_test.go|test TestUntagged function-absent|3f3b86a034e2ffa9aa8221047f7eb6cb08eac515|d8fc4bab7558b887d077b760389d8da02647f2a5\\\"; test rows map[inter… (533 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,345,12\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:253: test drift rows = 6, want 5: map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59:true int… (1020 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,323,12\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS20_gone,\"identity_test.go:244: drift lacks \\\"internal/example/gone_test.go|absent|test TestGone file-absent|46ba373cb6b4a058a49ea7d864dbb2b8cb9b8ba3|absent\\\"; test rows map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly… (1020 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,344,12\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:253: test drift rows = 4, want 5: map[internal/example/gone_test.go|absent|test TestGone file-absent|46ba373cb6b4a058a49ea7d864dbb2b8cb9b8ba3|absent:true internal/example/kept_test.go|internal/example/kept_test.go|test Test… (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS23_early,\"identity_test.go:244: drift lacks \\\"internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59\\\"; test rows map[internal/exa… (904 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,347,12\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:253: test drift rows = 7, want 5: map[internal/example/computed_test.go|internal/example/computed_test.go|test TestComputed/x function-absent|c3bbec34e0c6c5ab6b87591cc3855bdf1141ec1e|c3bbec34e0c6c5ab6b87591cc3855bdf1141ec1e… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS26_computed,\"identity_test.go:247: diagnostic of computed = \\\"\\\", want prefix \\\"cited test TestComputed/x is unsupported\\\"\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS26_deep,\"identity_test.go:247: diagnostic of deep = \\\"\\\", want prefix \\\"cited test TestSub/case one/inner is unsupported\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  6\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,347,12\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS27_tagged,\"identity_test.go:247: diagnostic of tagged = \\\"\\\", want prefix \\\"cited test TestTagged is unsupported-context\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-03-focused",
          "command": "bench test --package ./internal/preflight/... --run TestTestIdentity",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-03-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t03_rep1",
              "digest": "sha256:0541a8242d04dda68c9a6f1ed9d35e1917c8f3b539db4d1991415b1acc298c65",
              "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestTestIdentity,12,1,0,null,yes\n  2,bit,swap,failed,1,TestTestIdentity,12,1,0,null,yes\n  3,bit,swap,failed,1,TestTestIdentity/SS20_gone,12,1,0,null,yes\n  4,bit,swap,failed,2,TestTestIdentity,12,1,0,null,yes\n  5,bit,swap,failed,3,TestTestIdentity,12,1,0,null,yes\n  6,bit,swap,failed,1,TestTestIdentity/SS27_tagged,12,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestTestIdentity,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,8022\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,334,12\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:253: test drift rows = 2, want 5: map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59:true int… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS21_removed,\"identity_test.go:244: drift lacks \\\"internal/example/removed_test.go|internal/example/removed_test.go|test TestRemoved function-absent|7a1c971cb5216bc18a3548824c46a740c6c0e22f|d8fc4bab7558b887d077b760389d8da02647f2a5\\\"; test rows map[internal… (530 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS25_nosub,\"identity_test.go:244: drift lacks \\\"internal/example/kept_test.go|internal/example/kept_test.go|test TestSub/case two function-absent|8b7fd3111d935d68eaec647261034b4ba7727c50|8b7fd3111d935d68eaec647261034b4ba7727c50\\\"; test rows map[internal/… (529 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS27_untagged,\"identity_test.go:244: drift lacks \\\"internal/example/untagged_test.go|internal/example/untagged_test.go|test TestUntagged function-absent|3f3b86a034e2ffa9aa8221047f7eb6cb08eac515|d8fc4bab7558b887d077b760389d8da02647f2a5\\\"; test rows map[inter… (533 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,345,12\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:253: test drift rows = 6, want 5: map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59:true int… (1020 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,323,12\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS20_gone,\"identity_test.go:244: drift lacks \\\"internal/example/gone_test.go|absent|test TestGone file-absent|46ba373cb6b4a058a49ea7d864dbb2b8cb9b8ba3|absent\\\"; test rows map[internal/example/future_test.go|internal/example/future_test.go|test TestEarly… (1020 bytes)\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,344,12\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:253: test drift rows = 4, want 5: map[internal/example/gone_test.go|absent|test TestGone file-absent|46ba373cb6b4a058a49ea7d864dbb2b8cb9b8ba3|absent:true internal/example/kept_test.go|internal/example/kept_test.go|test Test… (731 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS23_early,\"identity_test.go:244: drift lacks \\\"internal/example/future_test.go|internal/example/future_test.go|test TestEarly already-shipped|97bdb2f7bdd8fcaf99b6af6ba4e33ef73a633747|3efbc754a2dfd0a7fddb7764c3caabee382b8d59\\\"; test rows map[internal/exa… (904 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,347,12\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity,\"identity_test.go:253: test drift rows = 7, want 5: map[internal/example/computed_test.go|internal/example/computed_test.go|test TestComputed/x function-absent|c3bbec34e0c6c5ab6b87591cc3855bdf1141ec1e|c3bbec34e0c6c5ab6b87591cc3855bdf1141ec1e… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS26_computed,\"identity_test.go:247: diagnostic of computed = \\\"\\\", want prefix \\\"cited test TestComputed/x is unsupported\\\"\",1\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS26_deep,\"identity_test.go:247: diagnostic of deep = \\\"\\\", want prefix \\\"cited test TestSub/case one/inner is unsupported\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  6\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,347,12\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestTestIdentity/SS27_tagged,\"identity_test.go:247: diagnostic of tagged = \\\"\\\", want prefix \\\"cited test TestTagged is unsupported-context\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t03_rep1-ss-03-check-1",
          "performer": "claude:ft375_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03_rep1",
            "digest": "sha256:786da8395fc724207d9ef43a43fd13e41ffe341704d901d13edef10c63a90610",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,52221,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,31187,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,53159,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,4683,56\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-03-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t03_rep1-ss-03-check-2",
          "performer": "claude:ft375_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03_rep1",
            "digest": "sha256:e62000c367683c8230f368f66f6ec4855a27e30dafd71658ab7e21691dfef81e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,pass,2748,145\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-03-check-2",
          "command": "bench test --package ./internal/coverage",
          "exit_code": 0
        },
        {
          "id": "ft375_t03_rep1-ss-03-check-3",
          "performer": "claude:ft375_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03_rep1",
            "digest": "sha256:c7ce33b2bb205d58eee6802f024ee9433325c26980af8b26e89974574fdb724b",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18264,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-03-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t03_rep1-ss-03-check-4",
          "performer": "claude:ft375_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03_rep1",
            "digest": "sha256:d26aeea25c0384c41025f0ca43b70d7e64a3ff14999779d36238c20f52df6d7e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18626,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-03-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft375_t03_rep1-ss-03-check-5",
          "performer": "claude:ft375_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t03_rep1",
            "digest": "sha256:f9cf5cabb0618dd56944a1f9f5a7137e085ec9643c44f25ce01ff43beb26bea4",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,42509,1118\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ss-03-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t04_rep1b-ss-04-focused",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:9fcce2c47be4198fcf5748739d289086024c4dd209394b61133e43c8e2b24f1b",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestAuditCompleteness/SS34_omitted_citation,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestCoverageIdentity|TestAuditCompleteness,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1330\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,1136,17\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestAuditCompleteness/SS34_omitted_citation,\"completeness_test.go:224: diagnostics = [], want the prefixes [\\\"metadata|row SS28 seam omits internal/example/kept_test.go (TestOther)\\\"]\",2\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestAuditCompleteness/SS50_present_empty,\"completeness_test.go:224: diagnostics = [\\\"metadata|row SS28 seam omits internal/example/kept_test.go (TestKept) from seams\\\" \\\"bench-planning-facts|the planning inventory declares no claim, so it cannot certify its own scope\\\"], want the prefi… (454 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestCoverageIdentity|TestAuditCompleteness'",
          "exit_code": 0,
          "probe": {
            "mutation": "ss-04-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft375_t04_rep1",
              "digest": "sha256:9fcce2c47be4198fcf5748739d289086024c4dd209394b61133e43c8e2b24f1b",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestAuditCompleteness/SS34_omitted_citation,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestCoverageIdentity|TestAuditCompleteness,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1330\nmutant_report[1]{mutant}:\n  1\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,fail,1136,17\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestAuditCompleteness/SS34_omitted_citation,\"completeness_test.go:224: diagnostics = [], want the prefixes [\\\"metadata|row SS28 seam omits internal/example/kept_test.go (TestOther)\\\"]\",2\n  github.com/gibbonmi/bench/internal/preflight/staleness,TestAuditCompleteness/SS50_present_empty,\"completeness_test.go:224: diagnostics = [\\\"metadata|row SS28 seam omits internal/example/kept_test.go (TestKept) from seams\\\" \\\"bench-planning-facts|the planning inventory declares no claim, so it cannot certify its own scope\\\"], want the prefi… (454 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft375_t04_rep1b-ss-04-check-1",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:dbf92298d4ac470aad7da7fccb0df80d97e576ed932148663717945601e7808a",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,51059,416\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,30029,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,53145,271\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,5017,56\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft375_t04_rep1b-ss-04-check-2",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:6d8d34b93a6414566918f91cc8be4f8279b8d7acb7a98d658fa9ae48faec13a9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,pass,2629,145\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-2",
          "command": "bench test --package ./internal/coverage",
          "exit_code": 0
        },
        {
          "id": "ft375_t04_rep1b-ss-04-check-3",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:865b179beb72d22ff5ee5a9f1e491dffaae4de130e62d1ff2dd6600d54b64a51",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5217,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-3",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft375_t04_rep1b-ss-04-check-4",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:02947b2d1c7ac1417cb22078a71c891eeafaed2db408fcc9cdd3e990f841e71b",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17740,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft375_t04_rep1b-ss-04-check-5",
          "performer": "claude:ft375_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "19c1928ba06437a64fb1ff4b5a8bd3150eb7ac53",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft375_t04_rep1",
            "digest": "sha256:c411e3da3166db940910b5d51cf0f90ef6de750c2d6a1596d359a21be3d141dd",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18026,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ss-04-check-5",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "ft375_ss-c2_standards_r1",
          "performer": "claude:ft375_ss-c2_standards_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c2_standards_r1",
            "digest": "sha256:89ac6d9dbbb5db60f8699e35559d5ee9c0cfa2212a7b2e01ca7bd5f49a8821a7",
            "excerpt": "Standards: 3 findings; worst: S1\n\nThe evidence check bound: sha256:9209cead… is current for assignment bb87c9fc…, base c585f9d0, source tip 83b3d693, dirty_record false. I read the delta c585f9d0..758d350d with specs and reviews excluded. Paths below are under the worktree /home/mgibs/.bench/worktrees/bench-2826441890/b8e8080d78971d37a88341856b9407bb-bb87c9fc56cc0af9f6e84908b2d9a3a6.\n\n- **S1** (item 1): internal/preflight/ownership/ownership.go:375-378 against internal/consumers/resolve.go:198-216. Rule: \"Two derivations of the same fact must collapse into one source.\" `SourceDeclaration` spells the declaration-kind words (\"func\", \"type\", \"const\", \"var\", \"method\") a second time, beside `objKind`. `Owner.Mismatch` (ownership.go:201) compares the baseline kind from `SourceDeclaration` with the tip kind from `match.Kind`, which comes from `objKind`. If one word changes in `objKind`, every cited symbol reads as a false `kind` discrepancy at the baseline. The syntax-side reading is necessary, but the second word list is not. Owner: ticket 02. Fix: one kind-word source that both sides use, probably exported from consumers. internal/consumers is outside 02's Writes, so the fix needs a fence-expansion plan commit. Confidence 7.\n- **S2** (item 5): internal/preflight/preflighttest/fixture.go:192 against internal/preflight/ownership/charge_fixture_test.go:166-182. Rule: fixture harnesses \"remain single-sourced\". The new `CommitChargeSource` repeats `chargeRepo.seed` step for step, with the same slug \"example\" and the same foo.go body. The steps are: checkout main, write files, SeedAdmission, commit, `checkout -B feature`, commit, `ChargeArgs`. The delta added a second copy and did not move `seed` onto it. Owner: ticket 04. Fix: `seed` calls `CommitChargeSource`; charge_fixture_test.go needs a fence-expansion plan commit. Confidence 7.\n- **S3** (item 2): internal/preflight/staleness/identity_test.go:18 against internal/preflight/ownership/premises_fixture_test.go:265. Same rule. `symbolPremise` spells the symbol-premise item shape again, including the rule that expected kind and signature repeat the owner. `seamItem` (staleness/inventory_test.go:25) also matches ownership/charge_fixture_test.go:25. A shared seam already exists: inventorytest, which completeness_test.go already uses. Owner: ticket 02 for the symbol builder; the seam builder came from 03 and 04. A wrong shape fails loudly when it decodes, so the drift risk is low. Disposition: repair through inventorytest, or the reviewer can veto it as incidental fixture text. Confidence 5.\n\nItems with no finding:\n- Item 3: the delta changes only the decoder that the anchor-kind entry calls (ownership.go:319-322). The map itself predates the delta, so this is advice only.\n- Item 4: no finding. Ticket 03:14 says \"Keep a single grammar owner; the coverage citation check calls the same seam\". `checkCitation` (coverage/citations.go:218) and staleness/tests.go:64 both call `ResolveTest`. The delta did not touch `declaresTestFunction` (evidence.go:165) or `rowSeamCitations` (premises.go:23).\n- Item 6: no finding, cohesion only. I found no binding size or placement rule that the move breaks.\n\n**Advice**\n- `anchorKind` is registered in `symbolPremises` (symbols.go) but now sits in helpers.go, next to the census plumbing and `topLevel`. It belongs in its own file.\n- ownership/evidence.go:165 `declaresTestFunction` is still a test-function grammar outside `ResolveTest`. This predates the delta; a later drain could fold it in.\n- `subjectPaths` is still a kind registry beside `symbolPremises` and `memberCensuses`. This came from SS-C1.\n- The string \"unsupported-context\" appears twice: in staleness `TestContext` (types.go:105) and in ownership `DiscrepancyContext` (ownership.go:169). The two types are separate, but the output word is the same, so consider one source.\n\n**Read:** the delta, the evidence check, consumers/resolve.go:190-229, ownership/evidence.go:150-199, ownership/premises.go:1-70, ownership/ownership.go:275-324, the two ownership fixture files (excerpts), tickets 02 and 03, and the AGENTS.md and BENCH.md rules from context.\n\n**Not verified:**\n- I did not re-read projects/benchkit.md, bench-craft-review SKILL.md, finding-discipline.md, or bench-craft-comments SKILL.md this session, so the comment-standard grading is not done.\n- I did not review ticket 04's full text.\n- I did not check consumer rows.\n- I did not check whether `completeness.go`'s claim count repeats an existing emptiness check in the manifest.\n"
          },
          "axis": "Standards",
          "base": "c585f9d08d9f7d320d5e56dfc88c69a09f43733c",
          "tip": "bc54d2f77717f0f4e14b14f73acb359879770db4",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "ft375_ss-c2_spec_r1",
          "performer": "claude:ft375_ss-c2_spec_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c2_spec_r1",
            "digest": "sha256:2bac6c5c598bd1a1be606d43a5ca094e6bfd7446de28d33444776fbbdddfda8d",
            "excerpt": "Spec: 1 finding; worst: P1\n\nEvidence binding: `bench preflight evidence sha256:9209cead… --check-current` bound the evidence (current=true, pair c585f9d0..83b3d693).\n\n- **P1, scope creep into ownership policy (author item 1).** Spec line 326 says \"Consume the landed ownership projections without modifying their policy.\" Ticket 04 line 155 says \"Preserve ordinary charge and after-ticket applicability without changing ownership policy.\" The change is in `internal/preflight/ownership/symbols.go:101`, inside the shared `identify` that `exact` now calls: `case name == owner.Symbol && len(all) > 1` returns ambiguous. Traced input: an owner with Package `…/internal/example` and bare Symbol `Dup`, where a second `Dup` exists in `…/sub`. Before this chunk, the package filter left one match and the premise passed. Now the premise faults. Every `exact` caller inherits the change (`rules.go:45,58,156` and `symbols.go:298`), so charge readiness for symbol, rule, and bound premises changes too. The package key is required (`metadata.go:138-141`), so the selection was not arbitrary. The audit could class ambiguity in staleness alone. No live spec breaks today: only this spec has planning facts, and all of its symbols are qualified. Owner: ticket 02. Disposition: ask-user. Confidence 7.\n\n**Row verdicts**\n- SS13–SS19: met (`identity_test.go` TestSymbolIdentity cases, plus `staleness/identity.go`).\n- SS20–SS27: met (`staleness/tests.go`, `coverage.ResolveTest`). The coverage check keeps its prior behavior.\n- SS28–SS31: met (TestCoverageIdentity, `coverage.go`, Decide).\n- SS32–SS34, SS49, SS50, SS57, SS58: met (TestAuditCompleteness, `completeness.go`).\n\n**Other author-flagged items**\n- **Item 2:** within spec. Line 144 says \"non-Go language … stay diagnostic.\" The ownership charge still refuses in that case; only the reason text changed.\n- **Item 3:** within spec. Lines 148/156 define test metadata as seams. Reused-test premises stay graded by the census, and a non-green one stays open through `completeness.go:864-868`.\n- **Item 4:** the change to line 134 is non-behavioral and follows the tree. It is fine for veto. Ticket 04 line 11 now matches the moved coverage-map seams (`completeness_test.go`). The plan expansion is consistent: ticket 02's Writes now includes the ownership files, and the fence and scope now hold `premises_fixture_test.go`.\n- **Item 5:** by reasoning, consistent with line 165. The tip metadata flips seams from planned to existing, so the baseline lacks them, which gives cited-tests red. The spec changed, which gives spec-drift red. The five native claims make the audit complete=false and manual_required=true. Reds exit 1 per line 225. I did not confirm this by a run (see Not verified).\n\n**Advice**\n- SS49 reads \"refuses complete audit collection.\" The build implements it as an exit-0 diagnostic with complete=false, not as an exit-1 refusal. That is a defensible reading, but the reviewer may want to confirm it.\n- A missing subtest is labeled `function-absent` (`tests.go:745-746`), although the parent function exists. The finding is still red; only the label is imprecise.\n\n**Read:** spec.md (all 650 lines); tickets 02, 03, 04 at 758d350d; the full chunk diff (`git diff c585f9d0..758d350d`, spilled to `/home/mgibs/.bench/responses/bench-2826441890/primary/1792300359564396741-5f9f5a3d103a2c18.out`); `ownership/helpers.go:180-205`; `ownership/evidence.go:30-84`; `ownership/metadata.go:136-145`; the SKILL and finding-discipline docs. `bench coverage` returned mapped.\n\n**Not verified:**\n- I could not run the `--audit-spec` audit, which refused both times. At pin 758d350d it refused with \"source not clean\" because HEAD is 83b3d693. At pin 83b3d693 it refused because the worktree became dirty: ` M internal/preflight/staleness/tests.go`. I did not write that file. Another writer is active in this review venue, so the coordinator should check it before accepting any axis result.\n- The plan-expansion `bench learning` record.\n- The behavior of the bare-symbol inventories in `specs/shared-test-fixtures/inventory.json`.\n"
          },
          "axis": "Spec",
          "base": "c585f9d08d9f7d320d5e56dfc88c69a09f43733c",
          "tip": "bc54d2f77717f0f4e14b14f73acb359879770db4",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "ft375_ss-c2_coverage_r1",
          "performer": "claude:ft375_ss-c2_coverage_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6987147fc849eda9d61f1577d1293a7f3d0e8f3c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft375_ss-c2_coverage_r1",
            "digest": "sha256:0d7878e462f043329546547cf3019ad02e4a810d2203591733ca09c7531263e1",
            "excerpt": "Coverage: 4 findings; worst: C1\n\nThe evidence binding check passed (current=true, dirty_record=false). From 758d350d to HEAD 83b3d693, only reviews/preflight-spec-staleness.md changes, so the source I graded is the frozen chunk. Worktree: /home/mgibs/.bench/worktrees/bench-2826441890/b8e8080d78971d37a88341856b9407bb-bb87c9fc56cc0af9f6e84908b2d9a3a6\n\n- **C1. Unreadable cited-test object.** The refusal at internal/preflight/staleness/tests.go:60-61 can be swallowed and every test stays green. That refusal is required behavior, per spec.md:127 (\"Unreadable before or after objects refuse collection.\") and the Collect contract at collect.go:19-21. Expected break: the audit grades an unreadable file as an empty source, so an existing seam reads as file-absent red, a planned seam reads as valid, and nothing refuses. Missing test: a refusal case in internal/preflight/staleness/drift_test.go beside TestAuditSafety. It can use a committed test file over the 2 MiB ControlRecordLimit (internal/bounds/bounds.go:55) and should assert \"cited test source refused\" with no staleness table. The same gap holds for StageSymbol (identity.go:57-59) and StageTicket (coverage.go:77-79): there is no test reference to any of the three stages. Probe P1 was silent. Owner: ticket 03 (tests.go); 02 and 04 for the siblings. Disposition: repair the tests. Confidence 8.\n- **C2. Deeper subtest citation (`TestSub/case one/x`).** ResolveTest at internal/coverage/citations.go:63 should return unjudged, which makes cited-tests diagnostic and Complete false. Turning off the `deeper` branch leaves the audit green and complete, which breaks SS26 and spec.md:466 (\"No unsupported state can yield complete true\"). Missing test: a TestTestIdentity case with a three-segment citation whose first segment is a declared subtest, expecting an \"unsupported\" diagnostic. Probe P2 was silent. Owner: ticket 03. Disposition: repair the test. Confidence 8.\n- **C3. Unchanged cited method.** SourceDeclaration has a baseline method-kind branch at internal/preflight/ownership/ownership.go:344-346. Turning it into \"func\" leaves the tests green. Expected break: an unchanged cited method gets a false red `kind` drift row. TestSymbolIdentity's only method premise changes at the tip, so the baseline comparison never runs for it. Missing test: an unchanged method premise in TestSymbolIdentity that adds no row (SS14/SS18 count assertion at identity_test.go:107). Probe P3 was silent. Owner: ticket 02. Disposition: repair the test. Confidence 7.\n- **C4. Baseline owner that does not parse as Go, while the tip resolves.** The unread case at internal/preflight/staleness/identity.go:35-36 should keep it diagnostic as unread (spec.md:144 non-Go stays diagnostic). Under the mutant it becomes a red `absent` row. The fixture's prose premise is unread at the tip, so it never reaches this baseline branch. Missing test: a premise whose base object has a Go syntax error and whose tip is valid. Probe P4 was silent. Owner: ticket 02. Disposition: repair the test. Confidence 5, because the red outcome also blocks completion and only the classification and exit code differ.\n\nProbes run:\n- **P1.** tests.go: StageTest refusal → `return \"\", \"\", nil`, package ./internal/preflight/.... Silent, restored yes.\n- **P2.** citations.go: `(!literal || deeper)` → `(!literal || deeper && false)`, coverage and preflight packages. Silent, restored yes. A first attempt was invalid (build failed) and was also restored.\n- **P3.** ownership.go: `kind = \"method\"` → `kind = \"func\"`. Silent, restored yes.\n- **P4.** identity.go: `case unread != \"\":` → `&& false`. Silent, restored yes.\n- **P5.** ownership.go: baseline `top.tok == token.CONST` → `false`. Bit TestSymbolIdentity/SS16_rule, restored yes.\n\n`git status --short` is empty. I ran every probe as a direct swap; no mutants JSON was written.\n\nAdvice: an unchanged const or var premise in TestSymbolIdentity would also guard against false red rows (optional; no row binds it). The order of the planned+unjudged branches in testAt (tests.go:63-68) has no witness. That is optional, because SS26 only requires the result to stay diagnostic.\n\nNot verified:\n- Whether a fixture over 2 MiB reaches ReadControlBlob before an earlier member or bounds refusal. If it does not, the C1 fixture needs another way to make the read fail.\n- Reachability of the StageSymbol and StageTicket refusals; I did not probe them.\n- Full preflight suite outside ./internal/preflight/... and ./internal/coverage/....\n- SS28–SS34, SS49, SS50, SS57 and SS58 beyond reading the diff. The listed prior probes already cover their main branches.\n"
          },
          "axis": "Coverage",
          "base": "c585f9d08d9f7d320d5e56dfc88c69a09f43733c",
          "tip": "bc54d2f77717f0f4e14b14f73acb359879770db4",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
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
      "from": "sha256:78e919fc76dd6b2ec59b518316b4b9687704d3a9976adbc7fa67671b320d82e1",
      "to": "sha256:59027695274b16599f6e50120d1bf38e3646525ba2f82a6fbb43bb0e7cc7afdb",
      "chunk_ids": {
        "SS-C1": [
          "SS-C1"
        ]
      }
    },
    {
      "from": "sha256:59027695274b16599f6e50120d1bf38e3646525ba2f82a6fbb43bb0e7cc7afdb",
      "to": "sha256:acbe78295e38ec098e24c56c60e2280c11bc01a9a30a5c1f2223ae14859fbd98",
      "chunk_ids": {
        "SS-C1": [
          "SS-C1"
        ]
      }
    },
    {
      "from": "sha256:acbe78295e38ec098e24c56c60e2280c11bc01a9a30a5c1f2223ae14859fbd98",
      "to": "sha256:212429719a60ae5886b1d9e625e5c99fb4bffba99d92bf66a48ea84ea45e78a7",
      "chunk_ids": {
        "SS-C1": [
          "SS-C1"
        ]
      }
    },
    {
      "from": "sha256:212429719a60ae5886b1d9e625e5c99fb4bffba99d92bf66a48ea84ea45e78a7",
      "to": "sha256:e74aa34bfd73177aca13de8907cef226127cc11275da45fa868160f4cd059831",
      "chunk_ids": {
        "SS-C1": [
          "SS-C1"
        ]
      }
    },
    {
      "from": "sha256:e74aa34bfd73177aca13de8907cef226127cc11275da45fa868160f4cd059831",
      "to": "sha256:9073893f67ccb984152307ed0a54ac13ca668f4ad8e750a0c534e524a1dc7ed5",
      "chunk_ids": {
        "SS-C1": [
          "SS-C1"
        ]
      }
    },
    {
      "from": "sha256:9073893f67ccb984152307ed0a54ac13ca668f4ad8e750a0c534e524a1dc7ed5",
      "to": "sha256:4be92c9f220cff6e03206483328acc81fd3de1894913f7c9c4df48caa0d0a87b",
      "chunk_ids": {
        "SS-C1": [
          "SS-C1"
        ],
        "SS-C2": [
          "SS-C2"
        ]
      }
    }
  ]
}
```

## SS-C1 review round 1

The frozen pair is 81b10e74..7a1535fd. Every finding falls on ticket 01, so one repair session takes all of them. Repair cycle 1 of 2 opens.

## Standards

Count: 3. Worst: S1.

- S1 (auto-fix, confidence 7): subjectPath in staleness/collect.go derives the premise subject grammar again. It reads only a top-level path. The nested owner, caller, target, directory, and pattern paths never join the audit scope. Rules: one source per fact, spec line 118, and SS12.
- S2 (auto-fix, confidence 6): staleness/collect.go members parses the ls-tree -z records a second time. internal/git/tree.go ReadTreeFile owns that grammar (AGENTS.md one source per fact).
- S3 (auto-fix, confidence 6): the drift_test.go header claims a red result in a comment (craft-comments).

## Spec

Count: 1. Worst: P1.

- P1 (auto-fix, confidence 7): eight SS-C1 seam locations in the planning metadata do not hold their needle. The seam rewrite kept the old row text (spec line 178). The audit prints eight metadata diagnostics.

## Coverage

Count: 4. Worst: C1.

- C1 (auto-fix, confidence 8): no test drops a baseline fence entry at the tip while its member changes. The probe was silent (spec line 111).
- C2 (auto-fix, confidence 8): no test places a special file only at the baseline; the probe was silent (SS51).
- C3 (auto-fix, confidence 7): no test names a member path with a quote, a backslash, a tab, or a non-ASCII byte. A probe without -z was silent.
- C4 (auto-fix, confidence 5): the tip manifest comes from the checkout bytes, not from the committed tip (spec line 127). The coordinator accepted it as a defect under --auto-approve.

## SS-C1 confirming round 1

Repair cycle 1 closed S1, S2, S3, P1, C2, C3, C4, and the coordinator finding X1. A consensus of a Fable/high and a Sonnet/high consultant decided C4: Collect drops its manifest parameter.

The coordinator probe for X1 bit at 040a2953. It swapped the headroom projection to return no path, and TestSubjectPaths failed with restored yes. That run is the recorded red for N2.

Repair cycle 2 of 2 opens for ticket 01. It takes C1 and N1 and runs on opus/medium by reviewer direction.

## Standards

Count: 2. Worst: N1.

- N1 (auto-fix, confidence 7): the ListTree comment in internal/git/tree.go says that it is the one reader of the ls-tree grammar. Three other parsers exist (craft-comments).
- N2 (no-op, confidence 5): TestSubjectPaths had no recorded red. The coordinator probe above records it.

## Spec

Count: 0. Worst: none.

## Coverage

Count: 1. Worst: C1.

- C1 (auto-fix, confidence 8): the fence half stays open. A probe that drops the baseline fence tokens from the scope was silent (spec line 112).

## SS-C1 confirming round 2

Repair cycle 2 closed N1 and the prefix part of C1. The coordinator probe inserted a skip of the baseline fence tokens, and TestDriftIdentity failed with restored yes.

Repair cycle 3 opens past the two-cycle allowance under --auto-approve. It takes C1b only, as a test-only repair on opus/medium.

## Standards

Count: 0. Worst: none.

## Spec

Count: 0. Worst: none.

## Coverage

Count: 1. Worst: C1b.

- C1b (auto-fix, confidence 7): no test drops an exact fence entry that only the baseline declares. A probe that removes the baseline exact tokens was silent (spec line 112).

## SS-C1 confirming round 3

Repair cycle 3 closed C1b. The coordinator probe made both passes read the tip fence tokens, and TestDriftIdentity failed with restored yes. All three axes pass, and no SS-C1 finding stays open. SS-C1 used three repair cycles, one past the allowance under --auto-approve.

## Standards

Count: 0. Worst: none.

## Spec

Count: 0. Worst: none.

## Coverage

Count: 0. Worst: none.

Advice: the Coverage probe in members failed first at the claim-drift assertion, not at the new exact-fence case.

## SS-C2 review round 1

The frozen pair is c585f9d0..758d350d. Repair cycle 1 of 2 opens with one session for each affected ticket, in the order 02, 03, 04, on opus/medium.

A consensus of a Fable/high and a Sonnet/high consultant decided P1: restore the prior exact behavior. The consensus also kept the SS49 diagnostic reading. The consultants split twice on S3, and the reviewer then ruled S3 advice only.

## Standards

Count: 3. Worst: S1.

- S1 (auto-fix, confidence 7, ticket 02): SourceDeclaration spells the declaration-kind words again beside consumers objKind (one source per fact).
- S2 (auto-fix, confidence 7, ticket 04): preflighttest CommitChargeSource repeats the ownership chargeRepo seed steps (fixture harnesses stay single-sourced).
- S3 (no-op, confidence 5): the staleness test item builders restate the manifest item shapes. The reviewer ruled it advice only.

## Spec

Count: 1. Worst: P1.

- P1 (auto-fix, confidence 7, ticket 02): the shared identify makes the ownership exact check fault an ambiguous bare symbol. That changes ownership policy (spec line 326). Restore exact and classify ambiguity only in the audit.

## Coverage

Count: 4. Worst: C1.

- C1 (auto-fix, confidence 8, tickets 02, 03, 04): no test reaches the unreadable-object refusal of the test, symbol, or ticket stage. A probe at the test stage was silent (spec line 127).
- C2 (auto-fix, confidence 8, ticket 03): no test cites a subtest deeper than one segment, so the unjudged branch has no witness (SS26).
- C3 (auto-fix, confidence 7, ticket 02): no test holds an unchanged cited method, so the baseline method kind has no witness.
- C4 (auto-fix, confidence 5, ticket 02): no test holds a baseline owner that does not parse while the tip resolves. The baseline unread branch has no witness.
