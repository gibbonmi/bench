# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/preflight-spec-staleness/spec.md",
  "plan_digest": "sha256:212429719a60ae5886b1d9e625e5c99fb4bffba99d92bf66a48ea84ea45e78a7",
  "implementation_session": "",
  "chunks": [
    {
      "id": "SS-C1",
      "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
      "tip": "c56e4a7f92e37107440b0e3fd4c092f9ed22188c",
      "plan_digest": "sha256:212429719a60ae5886b1d9e625e5c99fb4bffba99d92bf66a48ea84ea45e78a7",
      "source_digest": "91d8b3f0e12284e9d598d3a6f4a8f224a91534be",
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
