# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/candidate-binary-install/spec.md",
  "plan_digest": "sha256:4333de2593d7167b64eba95841fcbe4868d8e26e8bc5936f9005879a34ae4738",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CB-C1",
      "base": "97a979263a299429d187efcf6d1845ef93199ea5",
      "tip": "2465ab5cb31f0f579fded76b9e6fc1143767aa5d",
      "plan_digest": "sha256:4333de2593d7167b64eba95841fcbe4868d8e26e8bc5936f9005879a34ae4738",
      "source_digest": "a6d0a3153fa7dfbb6584fb8a50c2d93f2ce2b249",
      "acceptance_rows": [
        "CB1",
        "CB2",
        "CB3",
        "CB4",
        "CB5",
        "CB6",
        "CB7"
      ],
      "verification": [
        {
          "id": "t1-spec",
          "performer": "claude:ft419-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6d0a3153fa7dfbb6584fb8a50c2d93f2ce2b249",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t1-author-20261020",
            "digest": "sha256:9742f64358bf066e4dcd99f8ec732447ae19aa17b161f9ea4d40541ba91189b9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,pass,1772,165\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "spec",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0
        },
        {
          "id": "t1-root-conformance",
          "performer": "claude:ft419-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6d0a3153fa7dfbb6584fb8a50c2d93f2ce2b249",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t1-author-20261020",
            "digest": "sha256:9fb00d1b107ce78d210fb700bf164b7578ce4581f850e2adc5c750f70ca1f138",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13253,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-t1-fence-probe",
          "performer": "claude:ft419-t1-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6d0a3153fa7dfbb6584fb8a50c2d93f2ce2b249",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t1-author-20261020",
            "digest": "sha256:bd457de95d5916595b30e301afed12d8c7774d67f9a6833f26065b714248fd58",
            "excerpt": "mutants[8]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,5,TestPersistedFieldsHeader/one_line,165,1,0,null,yes\n  2,bit,swap,failed,2,TestPersistedFieldsHeader/a_fenced_line,165,1,0,null,yes\n  3,bit,omit,failed,2,TestPersistedFieldsHeader/a_fenced_line,165,1,0,null,yes\n  4,bit,swap,failed,1,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,,\",165,1,0,null,yes\n  5,bit,swap,failed,1,\"TestPersistedFieldsHeader/refuses_gate-evidence:x\",165,1,0,null,yes\n  6,bit,swap,failed,1,\"TestPersistedFieldsHeader/refuses_intent-ledger:commitment..gaps\",165,1,0,null,yes\n  7,bit,swap,failed,1,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,_review-record:folds,_intent-ledger:a\",165,1,0,null,yes\n  8,bit,swap,failed,4,\"TestPersistedFieldsHeader/refuses_gate-evidence:x\",165,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/spec,all,passed,9\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,6223\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1693,165\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/one_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{}, <nil>; want []spec.PersistedField{spec.PersistedField{Record:\\\"intent-ledger\\\", Path:[]string{\\\"commitment\\\", \\\"gaps\\\"}}, spec.PersistedField{Record:\\\"review-record\\\", Path:[]string{\\\"… (254 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_gate-evidence:x\",\"persisted_test.go:36: PersistedFields(\\\"gate-evidence:x\\\") error = <nil>; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,,\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a,,\\\") error = <nil>; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,_review-record:folds,_intent-ledger:a\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a, review-record:folds, intent-ledger:a\\\") error = <nil>; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:commitment..gaps\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:commitment..gaps\\\") error = <nil>; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1783,165\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/a_fenced_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{spec.PersistedField{Record:\\\"intent-ledger\\\", Path:[]string(nil)}}, <nil>; want []spec.PersistedField(nil), nil\",1\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/no_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{spec.PersistedField{Record:\\\"intent-ledger\\\", Path:[]string(nil)}}, <nil>; want []spec.PersistedField(nil), nil\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1732,165\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/a_fenced_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{spec.PersistedField{Record:\\\"review-record\\\", Path:[]string{\\\"folds\\\"}}}, <nil>; want []spec.PersistedField(nil), nil\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestRetireNextLineNamesTheBoardRemainder/FC6:_fenced_Roadmap_line_is_not_read\",\"retire_test.go:159: next line = \\\"next: promote durable content, remove the ROADMAP row FT7 and its roadmap/FT7.md detail file, commit as `spec-retire: s`\\\\\\\\n\\\", want \\\"next: promote durable content, remove the ROADMAP row FT<n> and its roadmap/… (291 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1677,165\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,,\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a,,\\\") error = <nil>; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1693,165\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_gate-evidence:x\",\"persisted_test.go:36: PersistedFields(\\\"gate-evidence:x\\\") error = <nil>; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  6\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1662,165\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:commitment..gaps\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:commitment..gaps\\\") error = <nil>; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  7\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1625,165\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,_review-record:folds,_intent-ledger:a\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a, review-record:folds, intent-ledger:a\\\") error = <nil>; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  8\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1722,165\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_gate-evidence:x\",\"persisted_test.go:36: PersistedFields(\\\"gate-evidence:x\\\") error = persisted-field-bad: unknown record \\\"gate-evidence\\\" in \\\"gate-evidence:x\\\"; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,,\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a,,\\\") error = persisted-field-bad: empty item; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,_review-record:folds,_intent-ledger:a\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a, review-record:folds, intent-ledger:a\\\") error = persisted-field-bad: duplicate item \\\"intent-ledger:a\\\"; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:commitment..gaps\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:commitment..gaps\\\") error = persisted-field-bad: empty path segment in \\\"intent-ledger:commitment..gaps\\\"; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-fence-probe",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the code-fence state in the Persisted fields header scan; TestPersistedFieldsHeader must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t1-author-20261020",
              "digest": "sha256:bd457de95d5916595b30e301afed12d8c7774d67f9a6833f26065b714248fd58",
              "excerpt": "mutants[8]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,5,TestPersistedFieldsHeader/one_line,165,1,0,null,yes\n  2,bit,swap,failed,2,TestPersistedFieldsHeader/a_fenced_line,165,1,0,null,yes\n  3,bit,omit,failed,2,TestPersistedFieldsHeader/a_fenced_line,165,1,0,null,yes\n  4,bit,swap,failed,1,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,,\",165,1,0,null,yes\n  5,bit,swap,failed,1,\"TestPersistedFieldsHeader/refuses_gate-evidence:x\",165,1,0,null,yes\n  6,bit,swap,failed,1,\"TestPersistedFieldsHeader/refuses_intent-ledger:commitment..gaps\",165,1,0,null,yes\n  7,bit,swap,failed,1,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,_review-record:folds,_intent-ledger:a\",165,1,0,null,yes\n  8,bit,swap,failed,4,\"TestPersistedFieldsHeader/refuses_gate-evidence:x\",165,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/spec,all,passed,9\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,6223\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1693,165\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/one_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{}, <nil>; want []spec.PersistedField{spec.PersistedField{Record:\\\"intent-ledger\\\", Path:[]string{\\\"commitment\\\", \\\"gaps\\\"}}, spec.PersistedField{Record:\\\"review-record\\\", Path:[]string{\\\"… (254 bytes)\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_gate-evidence:x\",\"persisted_test.go:36: PersistedFields(\\\"gate-evidence:x\\\") error = <nil>; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,,\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a,,\\\") error = <nil>; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,_review-record:folds,_intent-ledger:a\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a, review-record:folds, intent-ledger:a\\\") error = <nil>; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:commitment..gaps\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:commitment..gaps\\\") error = <nil>; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1783,165\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/a_fenced_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{spec.PersistedField{Record:\\\"intent-ledger\\\", Path:[]string(nil)}}, <nil>; want []spec.PersistedField(nil), nil\",1\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/no_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{spec.PersistedField{Record:\\\"intent-ledger\\\", Path:[]string(nil)}}, <nil>; want []spec.PersistedField(nil), nil\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1732,165\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/a_fenced_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{spec.PersistedField{Record:\\\"review-record\\\", Path:[]string{\\\"folds\\\"}}}, <nil>; want []spec.PersistedField(nil), nil\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestRetireNextLineNamesTheBoardRemainder/FC6:_fenced_Roadmap_line_is_not_read\",\"retire_test.go:159: next line = \\\"next: promote durable content, remove the ROADMAP row FT7 and its roadmap/FT7.md detail file, commit as `spec-retire: s`\\\\\\\\n\\\", want \\\"next: promote durable content, remove the ROADMAP row FT<n> and its roadmap/… (291 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1677,165\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,,\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a,,\\\") error = <nil>; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  5\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1693,165\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_gate-evidence:x\",\"persisted_test.go:36: PersistedFields(\\\"gate-evidence:x\\\") error = <nil>; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  6\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1662,165\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:commitment..gaps\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:commitment..gaps\\\") error = <nil>; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  7\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1625,165\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,_review-record:folds,_intent-ledger:a\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a, review-record:folds, intent-ledger:a\\\") error = <nil>; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  8\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1722,165\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_gate-evidence:x\",\"persisted_test.go:36: PersistedFields(\\\"gate-evidence:x\\\") error = persisted-field-bad: unknown record \\\"gate-evidence\\\" in \\\"gate-evidence:x\\\"; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,,\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a,,\\\") error = persisted-field-bad: empty item; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:a,_review-record:folds,_intent-ledger:a\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:a, review-record:folds, intent-ledger:a\\\") error = persisted-field-bad: duplicate item \\\"intent-ledger:a\\\"; want a persisted-fields-invalid error\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestPersistedFieldsHeader/refuses_intent-ledger:commitment..gaps\",\"persisted_test.go:36: PersistedFields(\\\"intent-ledger:commitment..gaps\\\") error = persisted-field-bad: empty path segment in \\\"intent-ledger:commitment..gaps\\\"; want a persisted-fields-invalid error\",1\nskips[0]{package,test,reason}:\n"
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
