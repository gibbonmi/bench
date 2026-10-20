# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/candidate-binary-install/spec.md",
  "plan_digest": "sha256:11740dde07e548048c06e174cbae4b13dbf2f9c97ab28592da7ffce6ed524e6e",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CB-C1",
      "base": "97a979263a299429d187efcf6d1845ef93199ea5",
      "tip": "057ee735bc394f4fe97e96f989972d476838f1e2",
      "plan_digest": "sha256:11740dde07e548048c06e174cbae4b13dbf2f9c97ab28592da7ffce6ed524e6e",
      "source_digest": "14b0fa583c14ebe9802594de774433d8e92afc3e",
      "acceptance_rows": [
        "CB1",
        "CB2",
        "CB3",
        "CB4",
        "CB5",
        "CB6",
        "CB7",
        "CB66",
        "CB67"
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
        },
        {
          "id": "t1r1-spec",
          "performer": "claude:ft419-t1-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "14b0fa583c14ebe9802594de774433d8e92afc3e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t1-repair1-20261020",
            "digest": "sha256:4de815e3526e13c69f24e75673de1bc924ebb452dbf25558d0150c7c0b90c5c0",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,pass,1722,173\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "spec",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0
        },
        {
          "id": "t1r1-root-conformance",
          "performer": "claude:ft419-t1-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "14b0fa583c14ebe9802594de774433d8e92afc3e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t1-repair1-20261020",
            "digest": "sha256:7fc08464c9b4f960cb5e280253f5150bafac692d5a9b3055b9732690d3b96996",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13992,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1r1-t1-fence-probe",
          "performer": "claude:ft419-t1-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "14b0fa583c14ebe9802594de774433d8e92afc3e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t1-repair1-20261020",
            "digest": "sha256:a4339bdc64a5f3c919b14586423a08a529d9e631cfbd8ea97aec2f6dd9efc893",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/spec.go,omit,failed,2,TestPersistedFieldsHeader/a_fenced_line,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,all,passed,173\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,829\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1672,173\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/a_fenced_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{spec.PersistedField{Record:\\\"review-record\\\", Path:[]string{\\\"folds\\\"}}}, <nil>; want []spec.PersistedField(nil), nil\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestRetireNextLineNamesTheBoardRemainder/FC6:_fenced_Roadmap_line_is_not_read\",\"retire_test.go:159: next line = \\\"next: promote durable content, remove the ROADMAP row FT7 and its roadmap/FT7.md detail file, commit as `spec-retire: s`\\\\\\\\n\\\", want \\\"next: promote durable content, remove the ROADMAP row FT<n> and its roadmap/… (291 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:bench-writer-high:ft419-t1-repair1-20261020",
              "digest": "sha256:a4339bdc64a5f3c919b14586423a08a529d9e631cfbd8ea97aec2f6dd9efc893",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/spec.go,omit,failed,2,TestPersistedFieldsHeader/a_fenced_line,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,all,passed,173\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,829\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,1672,173\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/a_fenced_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{spec.PersistedField{Record:\\\"review-record\\\", Path:[]string{\\\"folds\\\"}}}, <nil>; want []spec.PersistedField(nil), nil\",1\n  github.com/gibbonmi/bench/internal/spec,\"TestRetireNextLineNamesTheBoardRemainder/FC6:_fenced_Roadmap_line_is_not_read\",\"retire_test.go:159: next line = \\\"next: promote durable content, remove the ROADMAP row FT7 and its roadmap/FT7.md detail file, commit as `spec-retire: s`\\\\\\\\n\\\", want \\\"next: promote durable content, remove the ROADMAP row FT<n> and its roadmap/… (291 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "cbc1-standards",
          "performer": "claude:ft419-cbc1-standards-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a6d0a3153fa7dfbb6584fb8a50c2d93f2ce2b249",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc1-standards-20261020",
            "digest": "sha256:2d6366266cfa70c0e7000ce5ca96f7702b2a4fc73fa9a8faf41c5fdc25ddbd5e",
            "excerpt": "## Standards (CB-C1, pair 97a979263..3a5878479, opus/medium)\nS1 auto-fix conf 7: internal/spec/spec.go:103-104 calls ErrPersistedFieldsInvalid \"the refusal face\"; \"refusal face\" is a registered internal/refusalroute term and ticket 1 adds no face (bench-craft-comments: accurate terms).\nS2 auto-fix conf 5: internal/spec/spec.go:117-118 says the declaration reads \"through the header scan of Status: and Roadmap:\"; the shared seam is headerLines, which the AwaitsRetirement comment names.\nCount 2. Worst: S1. Advice kept separate (nil vs empty list assertion; last-wins unpinned; LiveSpecSlugs own fence scan).\n"
          },
          "axis": "Standards",
          "base": "97a979263a299429d187efcf6d1845ef93199ea5",
          "tip": "2465ab5cb31f0f579fded76b9e6fc1143767aa5d",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "cbc1-spec",
          "performer": "claude:ft419-cbc1-spec-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a6d0a3153fa7dfbb6584fb8a50c2d93f2ce2b249",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc1-spec-20261020",
            "digest": "sha256:ecac3b62582592a9ccfb6eaaf63f034bc1b7aa34aafd86f98bb9e47fb632e573",
            "excerpt": "## Spec (CB-C1, pair 97a979263..3a5878479, opus/medium)\nNo finding. CB1-CB7 met (spec.go:122-152, persisted_test.go); seam cells cited; API fits tickets 2, 3, 5, 6; no scope creep. Author choices last-wins and no segment character rule graded acceptable against spec.md:120,124,126.\nCount 0. Worst: none. Advice: untrimmed path segments; unrecorded plan-expansion learning.\n"
          },
          "axis": "Spec",
          "base": "97a979263a299429d187efcf6d1845ef93199ea5",
          "tip": "2465ab5cb31f0f579fded76b9e6fc1143767aa5d",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cbc1-coverage",
          "performer": "claude:ft419-cbc1-coverage-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a6d0a3153fa7dfbb6584fb8a50c2d93f2ce2b249",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc1-coverage-20261020",
            "digest": "sha256:296d0a43f2765f1265c55b539c4a17830603709b0f208cdfed40dc9fe3733548",
            "excerpt": "## Coverage (CB-C1, pair 97a979263..3a5878479, opus/medium)\nC1 ask-user conf 8: path segments with whitespace or a colon parse (intent-ledger: commitment.gaps, commitment. gaps, a b, a:b); restore would remove nothing; spec.go:126-134.\nC2 auto-fix conf 7: CB6 bites only interior \"..\"; mutant strings.Contains(path,\"..\") silent; add empty path, leading dot, trailing dot cases; persisted_test.go:36.\nC3 auto-fix conf 7: CB5 bites only gate-evidence; denylist mutant silent; add a second unknown name; persisted_test.go:35.\nC4 ask-user conf 7: a second unfenced or ~~~-fenced Persisted fields: line replaces the header declaration (last wins) untested and undecided; spec.go:115-119, headerLines spec.go:79.\nProbes: CB6 swap silent; CB5 denylist silent; CR strip omit silent (equivalent); first-wins silent; hostile insert bit; ~~~ insert bit; all restored=yes.\nCount 4. Worst: C1.\n"
          },
          "axis": "Coverage",
          "base": "97a979263a299429d187efcf6d1845ef93199ea5",
          "tip": "2465ab5cb31f0f579fded76b9e6fc1143767aa5d",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": []
        },
        {
          "id": "cbc1-standards-confirm",
          "performer": "claude:ft419-cbc1-standards-confirm-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "14b0fa583c14ebe9802594de774433d8e92afc3e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc1-standards-confirm-20261020",
            "digest": "sha256:2d1e018c787c5e88763877722444983308fb02bf8842e0fd2f12fb969afe32fb",
            "excerpt": "## Standards (CB-C1 confirming round, repair delta 90c9f2f75..e479d3763, opus/medium)\nS1 closed: internal/spec/spec.go:104 now names ErrPersistedFieldsInvalid as the sentinel that each PersistedFields refusal wraps.\nS2 closed: internal/spec/spec.go:117-121 now names headerLines, the shared fence-aware line view.\nNew findings: 0. Advice only (wording \"reviewed header declaration\", 2/10).\n"
          },
          "axis": "Standards",
          "base": "97a979263a299429d187efcf6d1845ef93199ea5",
          "tip": "057ee735bc394f4fe97e96f989972d476838f1e2",
          "finding_ids": [],
          "supersedes": [
            "cbc1-standards"
          ]
        },
        {
          "id": "cbc1-spec-confirm",
          "performer": "claude:ft419-cbc1-spec-confirm-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "14b0fa583c14ebe9802594de774433d8e92afc3e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc1-spec-confirm-20261020",
            "digest": "sha256:3402df8abfa91938761a7046d5140f440876c4e8888f35bd74d433bec7b99178",
            "excerpt": "## Spec (CB-C1 confirming round, repair delta 90c9f2f75..e479d3763, opus/medium)\nCB66 met: spec.md:325 and :124; spec.go:150-151; persisted_test.go refusal table.\nCB67 met: spec.md:326 and :124; spec.go:129-131.\nSeam cells cited; headerLines and metadata unchanged; no scope beyond S1, S2, C1-C4.\nNew findings: 0.\n"
          },
          "axis": "Spec",
          "base": "97a979263a299429d187efcf6d1845ef93199ea5",
          "tip": "057ee735bc394f4fe97e96f989972d476838f1e2",
          "finding_ids": [],
          "supersedes": [
            "cbc1-spec"
          ]
        },
        {
          "id": "cbc1-coverage-confirm",
          "performer": "claude:ft419-cbc1-coverage-confirm-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "14b0fa583c14ebe9802594de774433d8e92afc3e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc1-coverage-confirm-20261020",
            "digest": "sha256:7bfac9e33ba79b2cce1af88f09237b8682df1a1e5ab65aaf63d8ee4c4be76171",
            "excerpt": "## Coverage (CB-C1 confirming round, repair delta 90c9f2f75..e479d3763, opus/medium)\nC1/CB66, C2, C3, C4/CB67 closed; each named mutant bit, restored=yes.\nN1 auto-fix conf 7: mutant unicode.IsSpace(r) && r == ' ' at spec.go:150 stays silent; add a tab case and a U+00A0 case to the CB66 refusal list (persisted_test.go:27-40).\nN2 no-op conf 4: an identical second line is untested, but the code refuses every second line and the CB67 harm cannot occur.\n"
          },
          "axis": "Coverage",
          "base": "97a979263a299429d187efcf6d1845ef93199ea5",
          "tip": "057ee735bc394f4fe97e96f989972d476838f1e2",
          "finding_ids": [
            "N1"
          ],
          "supersedes": [
            "cbc1-coverage"
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
      "from": "sha256:4333de2593d7167b64eba95841fcbe4868d8e26e8bc5936f9005879a34ae4738",
      "to": "sha256:f61261f1ce4bc05e292c015b9a144c132d8610eed2d6316998c1b0492bfa7ad9",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ]
      }
    },
    {
      "from": "sha256:f61261f1ce4bc05e292c015b9a144c132d8610eed2d6316998c1b0492bfa7ad9",
      "to": "sha256:11740dde07e548048c06e174cbae4b13dbf2f9c97ab28592da7ffce6ed524e6e",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ]
      }
    }
  ]
}
```

## Standards

CB-C1: 2 findings. Worst: S1.

- S1 (auto-fix, confidence 7): `internal/spec/spec.go:103-104` calls `ErrPersistedFieldsInvalid` a refusal face; that term names a registered `internal/refusalroute` face, and ticket 1 adds none.
- S2 (auto-fix, confidence 5): `internal/spec/spec.go:117-118` names the `Status:`/`Roadmap:` scan as the shared seam; the shared seam is `headerLines`.

## Spec

CB-C1: 0 findings. Worst: none.

## Coverage

CB-C1: 4 findings. Worst: C1.

- C1 (auto-fix after consultation, confidence 8): a path segment with whitespace or a colon parses (`internal/spec/spec.go:126-134`). The opus high and sonnet high consultants agreed on a refusal with `persisted-fields-invalid`; new row CB66.
- C2 (auto-fix, confidence 7): CB6 bites only an interior `..`; add an empty path, a leading dot, and a trailing dot (`internal/spec/persisted_test.go:36`).
- C3 (auto-fix, confidence 7): CB5 bites only `gate-evidence`; add a second unknown record name (`internal/spec/persisted_test.go:35`).
- C4 (auto-fix after consultation, confidence 7): a second unfenced `Persisted fields:` line replaces the first (`internal/spec/spec.go:115-119`). The consultants agreed on a refusal of a second line with `persisted-fields-invalid`; new row CB67.

Confirming round (repair delta 90c9f2f75..e479d3763): S1, S2, C1, C2, C3, and C4 are closed.

- N1 (auto-fix, confidence 7): a space-only mutant at `internal/spec/spec.go:150` stays silent, because the CB66 refusal list holds only a plain space. Add a tab case and a U+00A0 case (`internal/spec/persisted_test.go:27-40`).
- N2 (no-op, confidence 4): an identical second line is untested, but the code refuses every second line, so the CB67 harm cannot occur.

Repair state: CB-C1 repair cycle 1 of 2 closed S1, S2, and C1 to C4. Cycle 2 of 2 is a test-only repair of N1, and it closes through a coordinator confirmation. The `--auto-approve` option of the run approves the CB66 and CB67 expansion.
