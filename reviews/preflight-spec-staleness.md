# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/preflight-spec-staleness/spec.md",
  "plan_digest": "sha256:78e919fc76dd6b2ec59b518316b4b9687704d3a9976adbc7fa67671b320d82e1",
  "implementation_session": "",
  "chunks": [
    {
      "id": "SS-C1",
      "base": "81b10e749ad6e3b655a3f030d06050e2ef6e0dde",
      "tip": "7a1535fdf85b124cd79d9106747269acadba764a",
      "plan_digest": "sha256:78e919fc76dd6b2ec59b518316b4b9687704d3a9976adbc7fa67671b320d82e1",
      "source_digest": "b4132eb97e1d6dfbd35824b91b6771d7fc431640",
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
