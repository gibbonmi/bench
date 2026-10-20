# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/candidate-binary-install/spec.md",
  "plan_digest": "sha256:a043f0448b7a5ecaf95aae6c6be7cfda8f21961bcd4df1281a437248f8789f86",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CB-C1",
      "base": "97a979263a299429d187efcf6d1845ef93199ea5",
      "tip": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
      "plan_digest": "sha256:90738f167f2b1ec19be414edfd1c6feb7cfeee86402486076859080b943cc922",
      "source_digest": "358f56d009d7b99e814b812c1a1c0af953cb870c",
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
        },
        {
          "id": "t1r2-spec",
          "performer": "claude:ft419-t1-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "358f56d009d7b99e814b812c1a1c0af953cb870c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t1-repair2-20261020",
            "digest": "sha256:ef7b223a8420e83f43f2301865ae5e77724aa3d5735ccacff486adde71448c01",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,pass,1726,175\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "spec",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0
        },
        {
          "id": "t1r2-root-conformance",
          "performer": "claude:ft419-t1-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "358f56d009d7b99e814b812c1a1c0af953cb870c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t1-repair2-20261020",
            "digest": "sha256:0d70282d8f49254a3c35d85c71ec5c0f695da923ec71c7b38c388629cbcfafb3",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13577,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1r2-t1-fence-probe",
          "performer": "claude:ft419-t1-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "358f56d009d7b99e814b812c1a1c0af953cb870c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t1-repair2-20261020",
            "digest": "sha256:7bd08470952aeb549083574e2bdb3777edc154acc62ac0ec69f05397121e2b33",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/spec.go,omit,failed,1,TestPersistedFieldsHeader/a_fenced_line,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,TestPersistedFieldsHeader,passed,18\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,437\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,2,18\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/a_fenced_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{spec.PersistedField{Record:\\\"review-record\\\", Path:[]string{\\\"folds\\\"}}}, <nil>; want []spec.PersistedField(nil), nil\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:bench-writer-high:ft419-t1-repair2-20261020",
              "digest": "sha256:7bd08470952aeb549083574e2bdb3777edc154acc62ac0ec69f05397121e2b33",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/spec.go,omit,failed,1,TestPersistedFieldsHeader/a_fenced_line,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,TestPersistedFieldsHeader,passed,18\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,437\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,2,18\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestPersistedFieldsHeader/a_fenced_line,\"persisted_test.go:22: PersistedFields = []spec.PersistedField{spec.PersistedField{Record:\\\"review-record\\\", Path:[]string{\\\"folds\\\"}}}, <nil>; want []spec.PersistedField(nil), nil\",1\nskips[0]{package,test,reason}:\n"
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
        },
        {
          "id": "cbc1-standards-r2",
          "performer": "claude:ft419-cbc1-standards-r2-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "358f56d009d7b99e814b812c1a1c0af953cb870c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc1-standards-r2-20261020",
            "digest": "sha256:e4869591c019a91a3616235d101efe6503d905ba141d71caf3d03c231072f325",
            "excerpt": "## Standards (CB-C1 repair cycle 2, delta 057ee735b..704a040fd, opus/medium)\nNo finding. The plan assignment record follows the repair1 shape.\nAdvice (no finding ID; the bounded repair policy makes a finding that grades only the N1 hardening check optional advice): the U+00A0 row at internal/spec/persisted_test.go:37 uses an invisible literal; an escape would show it. The three whitespace rows share one subtest base name.\n"
          },
          "axis": "Standards",
          "base": "97a979263a299429d187efcf6d1845ef93199ea5",
          "tip": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
          "finding_ids": [],
          "supersedes": [
            "cbc1-standards-confirm"
          ]
        },
        {
          "id": "cbc1-spec-r2",
          "performer": "claude:ft419-cbc1-spec-r2-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "358f56d009d7b99e814b812c1a1c0af953cb870c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc1-spec-r2-20261020",
            "digest": "sha256:c706ad5a52703dcfcf12e4d6a949ce5378142c02a440a39bcdaf8ebf46854887",
            "excerpt": "## Spec (CB-C1 repair cycle 2, delta 057ee735b..704a040fd, opus/medium)\nNo finding. The tab and U+00A0 rows serve CB66 (spec.md:124, :325). The plan JSON change only appends the repair2 assignment; with it removed, both plans canonicalize to one digest.\n"
          },
          "axis": "Spec",
          "base": "97a979263a299429d187efcf6d1845ef93199ea5",
          "tip": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
          "finding_ids": [],
          "supersedes": [
            "cbc1-spec-confirm"
          ]
        },
        {
          "id": "cbc1-coverage-r2",
          "performer": "claude:ft419-cbc1-coverage-r2-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "358f56d009d7b99e814b812c1a1c0af953cb870c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc1-coverage-r2-20261020",
            "digest": "sha256:16ba8bd1d36cc3aa68328665d6eb592bfe470da1f41271b5395f90658985333e",
            "excerpt": "## Coverage (CB-C1 repair cycle 2, delta 057ee735b..704a040fd, opus/medium)\nN1 closed. Swap to a space-only test bit 2 (tab, U+00A0); bypass to ASCII-only whitespace bit 1 (U+00A0); restored=yes each. No new finding.\n"
          },
          "axis": "Coverage",
          "base": "97a979263a299429d187efcf6d1845ef93199ea5",
          "tip": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
          "finding_ids": [],
          "supersedes": [
            "cbc1-coverage-confirm"
          ]
        }
      ]
    },
    {
      "id": "CB-C2",
      "base": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
      "tip": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
      "plan_digest": "sha256:8451fc1cd35152b165098f00691638091485dab2e16dd21032046eb32e8d79a6",
      "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
      "acceptance_rows": [
        "CB13",
        "CB14",
        "CB15",
        "CB16",
        "CB17",
        "CB18",
        "CB19",
        "CB20",
        "CB21",
        "CB23",
        "CB54",
        "CB55",
        "CB56",
        "CB63",
        "CB24",
        "CB25",
        "CB50",
        "CB51",
        "CB52",
        "CB53",
        "CB59",
        "CB62",
        "CB68"
      ],
      "verification": [
        {
          "id": "t3-t3-worktree",
          "performer": "claude:ft419-t3-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-author-20261020",
            "digest": "sha256:14da1a2e5709b47aa0a0aed8ef6ecd5b4a187ec253bcaa9927a644762c9a79f4",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,85256,1466\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t3-t3-system",
          "performer": "claude:ft419-t3-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-author-20261020",
            "digest": "sha256:857e05d8eecccfaadd3e52fc7f81a9c8904eee76f469411fc592f5c54203461e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,3,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,42077,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-system",
          "command": "bench test --check system --run TestCandidateInstall",
          "exit_code": 0
        },
        {
          "id": "t3-t3-internal-commitment-repository",
          "performer": "claude:ft419-t3-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-author-20261020",
            "digest": "sha256:b3d802b2e87e34037bbd3d7c0113f3847103d4f85745fe5dd4e76361cf982738",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,7744,165\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-internal-commitment-repository",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t3-t3-prose",
          "performer": "claude:ft419-t3-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-author-20261020",
            "digest": "sha256:f643359e24c69d79a36cb582133d6c7a99ad25c36ca16a8362bbeca50336cef3",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  prose,prose,0,173\n"
          },
          "requirement": "t3-prose",
          "command": "bench test --check prose",
          "exit_code": 0
        },
        {
          "id": "t3-t3-cmd-bench",
          "performer": "claude:ft419-t3-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-author-20261020",
            "digest": "sha256:db4e6a530e161a2d7609473ebf37709ab2464b3bad896df17f8a95ec34e39502",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19271,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t3-t3-root-conformance",
          "performer": "claude:ft419-t3-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-author-20261020",
            "digest": "sha256:59de2a18c1df2ac5a210303a5244f6c86e96d6e254c6e0870d07bd1096ee5399",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17799,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3-t3-open-build-probe",
          "performer": "claude:ft419-t3-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-author-20261020",
            "digest": "sha256:cf2730b67924e50f2a5b18ceee40a6da36166a38c3fe1f015e211a9e5b18dbff",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,1466,1,0,null,yes\n  2,bit,swap,failed,1,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,1466,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1023\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,85162,1466\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_undeclared_open_build_rebuilds,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,85223,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-open-build-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "ignore the open declared build in the held-install predicate; TestRefreshHoldsAnInstalledCandidate must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t3-author-20261020",
              "digest": "sha256:cf2730b67924e50f2a5b18ceee40a6da36166a38c3fe1f015e211a9e5b18dbff",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,1466,1,0,null,yes\n  2,bit,swap,failed,1,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,1466,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1023\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,85162,1466\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_undeclared_open_build_rebuilds,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,85223,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3-t3-digest-probe",
          "performer": "claude:ft419-t3-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-author-20261020",
            "digest": "sha256:cf2730b67924e50f2a5b18ceee40a6da36166a38c3fe1f015e211a9e5b18dbff",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,1466,1,0,null,yes\n  2,bit,swap,failed,1,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,1466,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1023\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,85162,1466\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_undeclared_open_build_rebuilds,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,85223,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-digest-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "ignore the digest equality in the held-install predicate; TestRefreshHoldsAnInstalledCandidate must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t3-author-20261020",
              "digest": "sha256:cf2730b67924e50f2a5b18ceee40a6da36166a38c3fe1f015e211a9e5b18dbff",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,1466,1,0,null,yes\n  2,bit,swap,failed,1,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,1466,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1023\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,85162,1466\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_undeclared_open_build_rebuilds,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,85223,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,\"land_effects_test.go:289: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2-t2-worktree",
          "performer": "claude:ft419-t2-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
            "digest": "sha256:589a3dc3f964d561b4aab577b00c7d3a8574be7f6035c062ca54cf19ce8277b6",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,82897,1466\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t2-t2-system",
          "performer": "claude:ft419-t2-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
            "digest": "sha256:730ff2d2fcd779f115302ff82a5b5d4bf648bd380989b4f5a48d1a734781f2f4",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,3,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,42396,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-system",
          "command": "bench test --check system --run TestCandidateInstall",
          "exit_code": 0
        },
        {
          "id": "t2-t2-cmd-bench",
          "performer": "claude:ft419-t2-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
            "digest": "sha256:c0891eb78e2c86d4c9dab801b8fe83b11f7bb98e522f727ff1ba73345420ea67",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17793,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t2-t2-internal-conformance",
          "performer": "claude:ft419-t2-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
            "digest": "sha256:a5f57a08abd4972d66770978e08879293f29a8a15e5fa069058d1b91916679b2",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38694,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-internal-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t2-t2-internal-intent",
          "performer": "claude:ft419-t2-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
            "digest": "sha256:02f2d273f30704c8e814d28b6562fd2c28580ec1e8bb388a50a0390dcac730c1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,pass,1138,78\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-internal-intent",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0
        },
        {
          "id": "t2-t2-internal-usage",
          "performer": "claude:ft419-t2-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
            "digest": "sha256:9dc3a2106ba5d02c9c72884fc4cdb6f759767003b2ad7ad20a784047d4059575",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/usage,pass,2,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-internal-usage",
          "command": "bench test --package ./internal/usage",
          "exit_code": 0
        },
        {
          "id": "t2-t2-internal-commitment-repository",
          "performer": "claude:ft419-t2-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
            "digest": "sha256:b6dfce86fa6a86405681c0c2e65e8e357b0db828d3aefb8f2f2ce3d5fe0aa80c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6658,165\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-internal-commitment-repository",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t2-t2-root-conformance",
          "performer": "claude:ft419-t2-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
            "digest": "sha256:c502f142c4c2557a573b7845144bd17c9aec29b43a7896328394d8928446f45e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16109,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-t2-path-probe",
          "performer": "claude:ft419-t2-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
            "digest": "sha256:ffb95014830fb07ce9d821b40d9744854227ab634fb0446858fd07fa4d110305",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,1466,1,0,null,yes\n  2,bit,swap,failed,1,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,1466,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1158\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,74339,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,\"build_install_test.go:224: install = (0, \\\"\\\"), want exit 1 naming \\\"the running executable /dev/shm/RXDLIO/t/TestBuildInstallrefuses_before_any_writea_worktree_executable_3587036432/001/.bench-home/worktrees/001-4254242932/37cf61bf432478769c9… (328 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,82552,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,\"build_install_test.go:239: install = (1, \\\"bench worktree build: the intent ledger does not read: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench status\\\\\\\\n\\\"), want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-path-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the primary-executable path refusal of the install; TestBuildInstall must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
              "digest": "sha256:ffb95014830fb07ce9d821b40d9744854227ab634fb0446858fd07fa4d110305",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,1466,1,0,null,yes\n  2,bit,swap,failed,1,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,1466,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1158\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,74339,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,\"build_install_test.go:224: install = (0, \\\"\\\"), want exit 1 naming \\\"the running executable /dev/shm/RXDLIO/t/TestBuildInstallrefuses_before_any_writea_worktree_executable_3587036432/001/.bench-home/worktrees/001-4254242932/37cf61bf432478769c9… (328 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,82552,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,\"build_install_test.go:239: install = (1, \\\"bench worktree build: the intent ledger does not read: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench status\\\\\\\\n\\\"), want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2-t2-tolerant-probe",
          "performer": "claude:ft419-t2-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
            "digest": "sha256:ffb95014830fb07ce9d821b40d9744854227ab634fb0446858fd07fa4d110305",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,1466,1,0,null,yes\n  2,bit,swap,failed,1,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,1466,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1158\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,74339,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,\"build_install_test.go:224: install = (0, \\\"\\\"), want exit 1 naming \\\"the running executable /dev/shm/RXDLIO/t/TestBuildInstallrefuses_before_any_writea_worktree_executable_3587036432/001/.bench-home/worktrees/001-4254242932/37cf61bf432478769c9… (328 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,82552,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,\"build_install_test.go:239: install = (1, \\\"bench worktree build: the intent ledger does not read: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench status\\\\\\\\n\\\"), want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-tolerant-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "replace the tolerant bound-assignment lookup of the install with the strict ledger read; TestBuildInstall must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t2-author-20261020",
              "digest": "sha256:ffb95014830fb07ce9d821b40d9744854227ab634fb0446858fd07fa4d110305",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,1466,1,0,null,yes\n  2,bit,swap,failed,1,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,1466,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1158\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,74339,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,\"build_install_test.go:224: install = (0, \\\"\\\"), want exit 1 naming \\\"the running executable /dev/shm/RXDLIO/t/TestBuildInstallrefuses_before_any_writea_worktree_executable_3587036432/001/.bench-home/worktrees/001-4254242932/37cf61bf432478769c9… (328 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,82552,1466\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,\"build_install_test.go:239: install = (1, \\\"bench worktree build: the intent ledger does not read: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench status\\\\\\\\n\\\"), want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3r1-t3-worktree",
          "performer": "claude:ft419-t3-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-repair1-20261020",
            "digest": "sha256:6847e2eaa720153c9b5143e21c2361f1589b59a12477db15120611bcef27060b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,76613,1470\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t3r1-t3-system",
          "performer": "claude:ft419-t3-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-repair1-20261020",
            "digest": "sha256:f2c8826cd3d28033db49fd0d7b9ca22fa4c17bc99c816ee469a1eff553e055b6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,3,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,38312,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-system",
          "command": "bench test --check system --run TestCandidateInstall",
          "exit_code": 0
        },
        {
          "id": "t3r1-t3-internal-commitment-repository",
          "performer": "claude:ft419-t3-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-repair1-20261020",
            "digest": "sha256:19b1eb353eb8a3e7ca66cae6690245ae0dde3bb74d996cd3685c7f5a33ea7091",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6770,165\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-internal-commitment-repository",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t3r1-t3-prose",
          "performer": "claude:ft419-t3-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-repair1-20261020",
            "digest": "sha256:f643359e24c69d79a36cb582133d6c7a99ad25c36ca16a8362bbeca50336cef3",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  prose,prose,0,173\n"
          },
          "requirement": "t3-prose",
          "command": "bench test --check prose",
          "exit_code": 0
        },
        {
          "id": "t3r1-t3-cmd-bench",
          "performer": "claude:ft419-t3-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-repair1-20261020",
            "digest": "sha256:3f7a4be9f39d765d8009b2d31e782c891178e77169d3931dac6dd3bcc1cfd899",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18702,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t3r1-t3-root-conformance",
          "performer": "claude:ft419-t3-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-repair1-20261020",
            "digest": "sha256:6467682ab74bdaeb46ca6dc54e3aebb499138207cecc8f898b9a581f81536aff",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14843,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3r1-t3-open-build-probe",
          "performer": "claude:ft419-t3-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-repair1-20261020",
            "digest": "sha256:8a94b3c1d316e5233d8e8c896420fad61f0a5447a499b5180db6cafd063b3cb0",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,1470,1,0,null,yes\n  2,bit,swap,failed,1,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,1470,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1023\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,74168,1470\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_undeclared_open_build_rebuilds,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72689,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-open-build-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "ignore the open declared build in the held-install predicate; TestRefreshHoldsAnInstalledCandidate must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t3-repair1-20261020",
              "digest": "sha256:8a94b3c1d316e5233d8e8c896420fad61f0a5447a499b5180db6cafd063b3cb0",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,1470,1,0,null,yes\n  2,bit,swap,failed,1,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,1470,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1023\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,74168,1470\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_undeclared_open_build_rebuilds,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72689,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3r1-t3-digest-probe",
          "performer": "claude:ft419-t3-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t3-repair1-20261020",
            "digest": "sha256:8a94b3c1d316e5233d8e8c896420fad61f0a5447a499b5180db6cafd063b3cb0",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,1470,1,0,null,yes\n  2,bit,swap,failed,1,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,1470,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1023\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,74168,1470\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_undeclared_open_build_rebuilds,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72689,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-digest-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "ignore the digest equality in the held-install predicate; TestRefreshHoldsAnInstalledCandidate must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t3-repair1-20261020",
              "digest": "sha256:8a94b3c1d316e5233d8e8c896420fad61f0a5447a499b5180db6cafd063b3cb0",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,1470,1,0,null,yes\n  2,bit,swap,failed,1,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,1470,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1023\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,74168,1470\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_invalid_declaration_counts_as_none,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/an_undeclared_open_build_rebuilds,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72689,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRefreshHoldsAnInstalledCandidate/a_declared_build_with_another_manifest_digest_rebuilds,\"land_effects_test.go:308: refresh = \\\"held\\\" (\\\"\\\"), want \\\"complete\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2r1-t2-worktree",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:8f1f2cce64338867f30309d487057cffa9315c7c0ef71395ea9541d1fb628a53",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,71084,1470\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-system",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:2790311de178cbc58e5c6a43479e563a2e49ad49aec2ca139f0a754ff15b147d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,3,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,36149,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-system",
          "command": "bench test --check system --run TestCandidateInstall",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-cmd-bench",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:706e8c00a6398c81f7e442c6cfb6375c1089785a91ea8bfb9cb5082e4f000fe4",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15411,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-internal-conformance",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:578b4165728588227013a5fe623f9cef22c3138405e875684906533a980c976f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,35563,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-internal-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-internal-intent",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:5645b53f606ecab435d8e9ad148b2673860678b07ca47b7dc18d2938afad1a79",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,pass,1109,78\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-internal-intent",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-internal-usage",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:9dc3a2106ba5d02c9c72884fc4cdb6f759767003b2ad7ad20a784047d4059575",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/usage,pass,2,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-internal-usage",
          "command": "bench test --package ./internal/usage",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-internal-commitment-repository",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:637bc142404f310b68f85f11903d9e6ad5459111f1145c0a711ddabd93f62cf4",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6056,165\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-internal-commitment-repository",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-commitmenttest",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:a6bd8772a7701bb2c44799879209bac9ed4990f39d308bf4f8adc700372e3ab1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/commitmenttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-commitmenttest",
          "command": "bench test --package ./internal/commitment/commitmenttest",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-root-conformance",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:b05b7c9971112734ed943ae6f5670638766207ca6ecc382d186176200e4f2165",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15602,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2r1-t2-path-probe",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:6d92ee7e39cf61d5ec44c11094dad7622d064ce7a1cb617119def154fdefd56c",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,1470,1,0,null,yes\n  2,bit,swap,failed,1,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,1470,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1158\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72007,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,\"build_install_test.go:222: install = (0, \\\"\\\"), want exit 1 naming \\\"the running executable /dev/shm/MMDZAH/t/TestBuildInstallrefuses_before_any_writea_worktree_executable_355482682/001/.bench-home/worktrees/001-2011273003/cd6e1aecbbe8b176dd5a… (327 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72795,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,\"build_install_test.go:237: install = (1, \\\"bench worktree build: the intent ledger does not read: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench status\\\\\\\\n\\\"), want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-path-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the primary-executable path refusal of the install; TestBuildInstall must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
              "digest": "sha256:6d92ee7e39cf61d5ec44c11094dad7622d064ce7a1cb617119def154fdefd56c",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,1470,1,0,null,yes\n  2,bit,swap,failed,1,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,1470,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1158\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72007,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,\"build_install_test.go:222: install = (0, \\\"\\\"), want exit 1 naming \\\"the running executable /dev/shm/MMDZAH/t/TestBuildInstallrefuses_before_any_writea_worktree_executable_355482682/001/.bench-home/worktrees/001-2011273003/cd6e1aecbbe8b176dd5a… (327 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72795,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,\"build_install_test.go:237: install = (1, \\\"bench worktree build: the intent ledger does not read: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench status\\\\\\\\n\\\"), want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2r1-t2-tolerant-probe",
          "performer": "claude:ft419-t2-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
            "digest": "sha256:6d92ee7e39cf61d5ec44c11094dad7622d064ce7a1cb617119def154fdefd56c",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,1470,1,0,null,yes\n  2,bit,swap,failed,1,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,1470,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1158\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72007,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,\"build_install_test.go:222: install = (0, \\\"\\\"), want exit 1 naming \\\"the running executable /dev/shm/MMDZAH/t/TestBuildInstallrefuses_before_any_writea_worktree_executable_355482682/001/.bench-home/worktrees/001-2011273003/cd6e1aecbbe8b176dd5a… (327 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72795,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,\"build_install_test.go:237: install = (1, \\\"bench worktree build: the intent ledger does not read: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench status\\\\\\\\n\\\"), want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-tolerant-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "replace the tolerant bound-assignment lookup of the install with the strict ledger read; TestBuildInstall must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t2-repair1-20261020",
              "digest": "sha256:6d92ee7e39cf61d5ec44c11094dad7622d064ce7a1cb617119def154fdefd56c",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,1470,1,0,null,yes\n  2,bit,swap,failed,1,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,1470,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1158\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72007,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/refuses_before_any_write/a_worktree_executable_runs,\"build_install_test.go:222: install = (0, \\\"\\\"), want exit 1 naming \\\"the running executable /dev/shm/MMDZAH/t/TestBuildInstallrefuses_before_any_writea_worktree_executable_355482682/001/.bench-home/worktrees/001-2011273003/cd6e1aecbbe8b176dd5a… (327 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,72795,1470\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildInstall/finds_the_binding_past_an_unknown_commitment_field,\"build_install_test.go:237: install = (1, \\\"bench worktree build: the intent ledger does not read: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench status\\\\\\\\n\\\"), want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "cbc2-standards",
          "performer": "claude:ft419-cbc2-standards-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc2-standards-20261020",
            "digest": "sha256:7ade91d3896956da41750e2d3e63b39cef993c60827b2d079978f74ecdff2a56",
            "excerpt": "## Standards (CB-C2, pair 704a040f..74648d83, opus/medium)\nS1 ask-user conf 8: raw ledger-field plant duplicated (build_install_test.go:114, candidate_install_test.go:142); one owner in internal/commitment/commitmenttest.\nS2 auto-fix conf 7: landable spec text repeated (candidate_install_test.go:60, owner_landing_fixture_test.go:58).\nS3 auto-fix conf 7: installAssignment (build.go:227) repeats resolveAssignmentIn (path.go:67) and drops missingTreeRefusal.\nS4 auto-fix conf 5: manifest path composed at build.go:260, :333, land_effects.go:135.\nS5 auto-fix conf 6: wrapper path spelled at build.go:161 and exec.go:204.\nS6 auto-fix conf 6: window_test.go:194-196 names a restore that does not exist.\nS7 auto-fix conf 5: window.go readOpenBuilds comment argues its own design.\nS8 auto-fix conf 7: candidate_install_test.go:49-50 comment does not parse.\nS9 no-op conf 8: HoldsOpenBuild has no production caller; ticket 3 asks for it.\nCount 9. Worst: S3.\n"
          },
          "axis": "Standards",
          "base": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
          "tip": "cd364bb6a77549e627ddbe9c406641bf12e8706c",
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
          "id": "cbc2-spec",
          "performer": "claude:ft419-cbc2-spec-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc2-spec-20261020",
            "digest": "sha256:1063d42fc3f9b99b8cded24572da1e2007c3d86819ddf161940dd090e902d11d",
            "excerpt": "## Spec (CB-C2, pair 704a040f..74648d83, opus/medium)\nF1 auto-fix conf 7: CB56 is partial at its seam; the TestBuildInstall case starts with no earlier manifest. Amend the seam cell to also cite TestPublishRestoresPriorPairWhenSealPromotionFails.\nEvery other CB-C2 row met. Refusal order, next line, help placement, fail-closed hold, invalid-as-none, CB59 wording graded no-op.\nCount 1. Worst: F1.\n"
          },
          "axis": "Spec",
          "base": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
          "tip": "cd364bb6a77549e627ddbe9c406641bf12e8706c",
          "finding_ids": [
            "F1"
          ],
          "supersedes": []
        },
        {
          "id": "cbc2-coverage",
          "performer": "claude:ft419-cbc2-coverage-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "60add18bbfab6ed169fea56e9ffeb226aaff3071",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc2-coverage-20261020",
            "digest": "sha256:8c26ddcb2f2d64d0faf7b5f265c321a272bd901cd6612db7398ce4701c540aec",
            "excerpt": "## Coverage (CB-C2, pair 704a040f..74648d83, opus/medium)\nF1 auto-fix conf 8: no row or test has two open builds; mutant \"first build only\" in window.go DeclaredBuild stays silent (257 tests), restored=yes.\nF2 auto-fix conf 6: canonical physical path compare untested; mutant plain string compare at build.go:212-222 silent, restored=yes.\nF3 auto-fix conf 6: CB56 never checks an earlier manifest; mutant removing the manifest on the publishCandidate error path silent, restored=yes.\nCount 3. Worst: F1.\n"
          },
          "axis": "Coverage",
          "base": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
          "tip": "cd364bb6a77549e627ddbe9c406641bf12e8706c",
          "finding_ids": [
            "F1",
            "F2",
            "F3"
          ],
          "supersedes": []
        },
        {
          "id": "cbc2-standards-confirm",
          "performer": "claude:ft419-cbc2-standards-confirm-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc2-standards-confirm-20261020",
            "digest": "sha256:593dcc99df3e8208b2cebd065cdbb9588c34da8433309d7ff9ee57c64fb4531c",
            "excerpt": "## Standards (CB-C2 confirming round, repair delta 74648d83..5ca3a61d, opus/medium)\nS1 to S8 closed: one ledger-field plant (commitmenttest/admission.go:195-224), landableSpec, resolveAmong with the missing-tree refusal, manifestPath, wrapperPath, three comment fixes.\nJudgment calls graded acceptable: ledger read before the target-control check; the install gains the branch-checked-out check.\nNew findings: 0.\n"
          },
          "axis": "Standards",
          "base": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
          "tip": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "finding_ids": [],
          "supersedes": [
            "cbc2-standards"
          ]
        },
        {
          "id": "cbc2-spec-confirm",
          "performer": "claude:ft419-cbc2-spec-confirm-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc2-spec-confirm-20261020",
            "digest": "sha256:ab3ce3ee49c10a1f5d162f377f5a0be03d9a09debf3498507937ce5b71411b76",
            "excerpt": "## Spec (CB-C2 confirming round, repair delta 74648d83..5ca3a61d, opus/medium)\nCB56 met (seam cites TestBuildInstall and TestPublishRestoresPriorPairWhenSealPromotionFails). CB68 met and cited. No behavior beyond S1-S8 and F1-F3; plan changes are only the widenings, CB68, t2-commitmenttest, and repair assignments.\nNew findings: 0.\n"
          },
          "axis": "Spec",
          "base": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
          "tip": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "finding_ids": [],
          "supersedes": [
            "cbc2-spec"
          ]
        },
        {
          "id": "cbc2-coverage-confirm",
          "performer": "claude:ft419-cbc2-coverage-confirm-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e443b2673f804c8cb95cecd333658eae36e8524c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc2-coverage-confirm-20261020",
            "digest": "sha256:8a5aabb02245f76d11c7036d084910e03e899d8d553a3dea84a2a27a9b69b93f",
            "excerpt": "## Coverage (CB-C2 confirming round, repair delta 74648d83..5ca3a61d, opus/medium)\nF1/CB68, F2, F3, and the S3 missing-tree case closed; each independent mutant bit, restored=yes.\nAdvice (no finding ID; the bounded repair policy makes a finding that grades only the CB68 hardening check optional advice): no row has the declared build first and an undeclared build second; a last-build-only mutant in DeclaredBuild stays silent. The production loop checks every build.\nThe resolveAmong lineSafe gap is inherited and goes to bench learning.\n"
          },
          "axis": "Coverage",
          "base": "704a040fda0c0c5a3f50798252ffb638aa978ab5",
          "tip": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "finding_ids": [],
          "supersedes": [
            "cbc2-coverage"
          ]
        }
      ]
    },
    {
      "id": "CB-C3",
      "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
      "tip": "cee87305b25de0e474628722a3936eed0c3e41a5",
      "plan_digest": "sha256:0c508811425128eb96494093745f4074bfedd49a4b0941685db08436845bf193",
      "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
      "acceptance_rows": [
        "CB26",
        "CB27",
        "CB28",
        "CB29",
        "CB30",
        "CB31",
        "CB32",
        "CB33",
        "CB34",
        "CB35",
        "CB36",
        "CB57",
        "CB60",
        "CB61",
        "CB65",
        "CB69"
      ],
      "verification": [
        {
          "id": "t5-owner",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:f371156dad6b5b2d5aa1cd3ad31e74cea95f10cf5c1f72a590ca8f645b490709",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6117,175\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t5-manifest",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:a8ae0a8515973657d26dbd7a8f80c05ff01ad65ae10d8183358ad0c79e71228f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/brokermanifest,pass,2,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "manifest",
          "command": "bench test --package ./internal/brokermanifest",
          "exit_code": 0
        },
        {
          "id": "t5-gate",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:e61df139370c867933ae5a852e067a9a4d961463fdeb424815b0718a6d57d86a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,27694,507\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t5-record",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:8bea61b11c131565fc800aef55bfc2a700b10ecef67eaeac7d304308f2266b8f",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,2749,110\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,16955,257\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "record",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0
        },
        {
          "id": "t5-routes",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:4d1e7f47f88e58feec6b457965f88016f4a7b0ca8be8f02cfc1e0111655e6989",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,5,21\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "routes",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "t5-cmd-bench",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:d1c3180a818e255aff7816d71b7e956bf8d3abf64c17c0612effd26d440abb4b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14280,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t5-internal-conformance",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:fa68b22cbffcc3653f6384275f1a5f7d3cc1c3bf572ed29a3050f04fda490bc5",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,33339,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "internal-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t5-root-conformance",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:d8ae6e7df63186103a0323f96a325e91712f38e1609ef8ed0f7b714957290c80",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15340,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t5-system",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:52e5827486a06163961dd024f3cebce64f92e5cc98fa8d9438ae926f54ec012a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,3,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,33037,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "system",
          "command": "bench test --check system --run TestCandidateInstall",
          "exit_code": 0
        },
        {
          "id": "t5-t5-digest-probe",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:3425fec3d0f788a6780b7bae4cc14bd0c98607b9ad0b15131c86f55b9774b6a2",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,175,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,175,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1057\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5846,175\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5925,175\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-digest-probe",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the digest comparison of the candidate check; TestCandidateCheck must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
              "digest": "sha256:3425fec3d0f788a6780b7bae4cc14bd0c98607b9ad0b15131c86f55b9774b6a2",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,175,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,175,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1057\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5846,175\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5925,175\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5-t5-locator-probe",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:3425fec3d0f788a6780b7bae4cc14bd0c98607b9ad0b15131c86f55b9774b6a2",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,175,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,175,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1057\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5846,175\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5925,175\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-locator-probe",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "read the broker manifest beside the caller root instead of the primary checkout; TestCandidateCheck must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
              "digest": "sha256:3425fec3d0f788a6780b7bae4cc14bd0c98607b9ad0b15131c86f55b9774b6a2",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,175,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,175,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1057\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5846,175\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5925,175\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:105: check = <nil>; want a candidate-not-installed refusal true\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5-t5-chunk-order-probe",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:c5ab6facc20cca3b6c82b2b92ee2a020d978a0e37fd5b042c8bbeb6bcc670e1e",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestCheckpointFacesFollowTheirRoutes/checkpoint-subject-unavailable,507,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1295\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28086,507\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointFacesFollowTheirRoutes/checkpoint-subject-unavailable,\"refusal_route_test.go:139: face checkpoint-subject-unavailable = (1, \\\"complete checkpoint source unavailable: exit status 128\\\\\\\\nnext=bench doctor; then bench gate --in 'gate-subject' --fresh --checkpoint 'specs/example/spec.md' --complete\\\\\\\\nr… (519 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:111: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (808 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:145: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-chunk-order-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "move the candidate check of the chunk checkpoint after its execute path; TestCheckpointPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
              "digest": "sha256:c5ab6facc20cca3b6c82b2b92ee2a020d978a0e37fd5b042c8bbeb6bcc670e1e",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestCheckpointFacesFollowTheirRoutes/checkpoint-subject-unavailable,507,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1295\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28086,507\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointFacesFollowTheirRoutes/checkpoint-subject-unavailable,\"refusal_route_test.go:139: face checkpoint-subject-unavailable = (1, \\\"complete checkpoint source unavailable: exit status 128\\\\\\\\nnext=bench doctor; then bench gate --in 'gate-subject' --fresh --checkpoint 'specs/example/spec.md' --complete\\\\\\\\nr… (519 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:111: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (808 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:145: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5-t5-complete-order-probe",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:73b1dc1f58309eec60b95bf22e6e7591d7219a6709f10cd45db66fcb09a3b02b",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCheckpointPersistedFields/complete,507,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,430\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,27876,507\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:116: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-complete-order-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the candidate check of the complete checkpoint; TestCheckpointPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
              "digest": "sha256:73b1dc1f58309eec60b95bf22e6e7591d7219a6709f10cd45db66fcb09a3b02b",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCheckpointPersistedFields/complete,507,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,430\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,27876,507\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:116: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5-t5-record-order-probe",
          "performer": "claude:ft419-t5-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
            "digest": "sha256:452d62cdb4bce970f0125f72073e912a37a3e2c59663b5a67810e9d3fda29112",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestAssignmentPersistedFields/malformed_declaration,257,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord/recordcmd,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1007\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,16170,257\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/malformed_declaration,\"assignment_persisted_test.go:62: assignment = (0, \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",f0bf8e48d7f9008edbeabff37ba94d568c5d5784\\\\\\\\n\\\"), want the persisted-fields-invalid error and… (249 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/uninstalled,\"assignment_persisted_test.go:48: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",6ad7669c5a7a2556d9cf97623edf4fccb21b9cde\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not… (251 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-record-order-probe",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the candidate check of the assignment entry; TestAssignmentPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-author-20261020",
              "digest": "sha256:452d62cdb4bce970f0125f72073e912a37a3e2c59663b5a67810e9d3fda29112",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestAssignmentPersistedFields/malformed_declaration,257,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord/recordcmd,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1007\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,16170,257\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/malformed_declaration,\"assignment_persisted_test.go:62: assignment = (0, \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",f0bf8e48d7f9008edbeabff37ba94d568c5d5784\\\\\\\\n\\\"), want the persisted-fields-invalid error and… (249 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/uninstalled,\"assignment_persisted_test.go:48: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",6ad7669c5a7a2556d9cf97623edf4fccb21b9cde\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not… (251 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r1-owner",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:322828ee1486e14b755fedaa600553ddba6f80ac12463bf93c42293ff7fed524",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6137,179\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t5r1-manifest",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:a8ae0a8515973657d26dbd7a8f80c05ff01ad65ae10d8183358ad0c79e71228f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/brokermanifest,pass,2,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "manifest",
          "command": "bench test --package ./internal/brokermanifest",
          "exit_code": 0
        },
        {
          "id": "t5r1-gate",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:ca9e2a5718e3b408a0917be6cc5b7787e8bf78de35f9c0cb359942b5102e4fc6",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,28413,508\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t5r1-record",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:306b3d71c6fb39c0674e01df3653ce028357bbb1b111ee28bef99b6b5034fac5",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,3061,110\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,17549,260\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "record",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0
        },
        {
          "id": "t5r1-routes",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:083201dbdeffa4beaddc7caec9539221723505998f67c6ee5a7c07652c31da24",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,22\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "routes",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "t5r1-cmd-bench",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:d3908f673ab524f4cab530d772efe91eca1a717749edd93bac318396b4ea9a7f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17795,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t5r1-internal-conformance",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:44671fa36f35c4657f2fabad186aa459b21dce82c43a029c078344ec83c9250b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,33293,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "internal-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t5r1-commitmenttest",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:a6bd8772a7701bb2c44799879209bac9ed4990f39d308bf4f8adc700372e3ab1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/commitmenttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "commitmenttest",
          "command": "bench test --package ./internal/commitment/commitmenttest",
          "exit_code": 0
        },
        {
          "id": "t5r1-root-conformance",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:dccb3fefd353b2244a55ba3faae56612de2ef6a8529e06fb0ff67a12ae3cfc24",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13924,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t5r1-system",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:b1d2b5e80382cbf2a9f38fdf9152347355aa001b45d6a8e82e073854153c4420",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,3,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,34754,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "system",
          "command": "bench test --check system --run TestCandidateInstall",
          "exit_code": 0
        },
        {
          "id": "t5r1-t5-digest-probe",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:b46ef4ff875d1954a63a0ba7f4eefecf875f2aab1aeceaa29e9298cdc3af67b2",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,179,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,179,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1159\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6298,179\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6186,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-digest-probe",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the digest comparison of the candidate check; TestCandidateCheck must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
              "digest": "sha256:b46ef4ff875d1954a63a0ba7f4eefecf875f2aab1aeceaa29e9298cdc3af67b2",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,179,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,179,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1159\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6298,179\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6186,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r1-t5-locator-probe",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:b46ef4ff875d1954a63a0ba7f4eefecf875f2aab1aeceaa29e9298cdc3af67b2",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,179,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,179,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1159\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6298,179\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6186,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-locator-probe",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "read the broker manifest beside the caller root instead of the primary checkout; TestCandidateCheck must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
              "digest": "sha256:b46ef4ff875d1954a63a0ba7f4eefecf875f2aab1aeceaa29e9298cdc3af67b2",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,179,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,179,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1159\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6298,179\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6186,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:142: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r1-t5-chunk-order-probe",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:130b68231d9e1648a455381350d96fbcef7152e79888675be4f4b7495d03c987",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,2,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1348\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28946,508\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (808 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28927,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-chunk-order-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "move the candidate check of the chunk checkpoint after its execute path; TestCheckpointPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
              "digest": "sha256:130b68231d9e1648a455381350d96fbcef7152e79888675be4f4b7495d03c987",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,2,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1348\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28946,508\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (808 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28927,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r1-t5-complete-order-probe",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:130b68231d9e1648a455381350d96fbcef7152e79888675be4f4b7495d03c987",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,2,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1348\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28946,508\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (808 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28927,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-complete-order-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the candidate check of the complete checkpoint; TestCheckpointPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
              "digest": "sha256:130b68231d9e1648a455381350d96fbcef7152e79888675be4f4b7495d03c987",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,2,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1348\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28946,508\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (808 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28927,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r1-t5-record-order-probe",
          "performer": "claude:ft419-t5-repair1-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
            "digest": "sha256:719e2c2409a301cde72e4cc0ce3b1c585e4bd45e87a2c79a1436f211a6a63d68",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,4,TestAssignmentPersistedFields/an_undeclared_working_tree,260,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord/recordcmd,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1714\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,17069,260\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/an_undeclared_working_tree,\"assignment_persisted_test.go:86: assignment = exit 1, output \\\"error: bench record assignment refused — working specs/example/spec.md differs from its source; commit or discard the change first\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not-installed… (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/hostile_slug_and_flag_value,\"assignment_persisted_test.go:78: assignment = \\\"error: bench record assignment refused — missing or nonregular tree file specs/it's here/tickets/1.md\\\\\\\\n\\\", want one next= route that names no gate\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/malformed_declaration,\"assignment_persisted_test.go:92: assignment = (0, \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",6baf3a023d69264fc733f6ae153afaa37b9cdb80\\\\\\\\n\\\"), want the persisted-fields-invalid error and… (249 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/uninstalled,\"assignment_persisted_test.go:63: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",e071ed1aefe3034f3b736bf0c770548f7a533c1c\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not… (251 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-record-order-probe",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the candidate check of the assignment entry; TestAssignmentPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair1-20261020",
              "digest": "sha256:719e2c2409a301cde72e4cc0ce3b1c585e4bd45e87a2c79a1436f211a6a63d68",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,4,TestAssignmentPersistedFields/an_undeclared_working_tree,260,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord/recordcmd,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1714\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,17069,260\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/an_undeclared_working_tree,\"assignment_persisted_test.go:86: assignment = exit 1, output \\\"error: bench record assignment refused — working specs/example/spec.md differs from its source; commit or discard the change first\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not-installed… (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/hostile_slug_and_flag_value,\"assignment_persisted_test.go:78: assignment = \\\"error: bench record assignment refused — missing or nonregular tree file specs/it's here/tickets/1.md\\\\\\\\n\\\", want one next= route that names no gate\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/malformed_declaration,\"assignment_persisted_test.go:92: assignment = (0, \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",6baf3a023d69264fc733f6ae153afaa37b9cdb80\\\\\\\\n\\\"), want the persisted-fields-invalid error and… (249 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/uninstalled,\"assignment_persisted_test.go:63: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",e071ed1aefe3034f3b736bf0c770548f7a533c1c\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not… (251 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r2-owner",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:114c90fb03c89d0d9fb8fe46eaaf56ecb859d70286f28b3a07484914dee2748f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6117,180\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t5r2-manifest",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:a8ae0a8515973657d26dbd7a8f80c05ff01ad65ae10d8183358ad0c79e71228f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/brokermanifest,pass,2,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "manifest",
          "command": "bench test --package ./internal/brokermanifest",
          "exit_code": 0
        },
        {
          "id": "t5r2-gate",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:3c38022a3da09d39cd40747e699bca874f30f1aec85593a5b8cdbeff7605e4ed",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,29420,508\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t5r2-record",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:2d1baa929bbdbd2c517201c141143c5ac7335647112518d9b94e0059d6291d86",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,2759,110\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,16879,260\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "record",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0
        },
        {
          "id": "t5r2-routes",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:083201dbdeffa4beaddc7caec9539221723505998f67c6ee5a7c07652c31da24",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,22\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "routes",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "t5r2-cmd-bench",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:d65477a2a6ac41254b4e7984f9c69544b018b355ac6144ca5e4ffc429193e189",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16405,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t5r2-internal-conformance",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:a373e2bf24f109143baf7a6048d1230d635fb4b095d354ca9e8f0a3f7236aa38",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,33584,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "internal-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t5r2-commitmenttest",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:a6bd8772a7701bb2c44799879209bac9ed4990f39d308bf4f8adc700372e3ab1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/commitmenttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "commitmenttest",
          "command": "bench test --package ./internal/commitment/commitmenttest",
          "exit_code": 0
        },
        {
          "id": "t5r2-root-conformance",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:eaf203d656bbab62a31a0056dba44b9d0740d0c25aa6c164605b10d0627e2f04",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13756,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t5r2-system",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:49565c748b3e86cdee7a621bb677b8b9b568ea2cfe78296106bdc49a723994d2",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,3,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,34022,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "system",
          "command": "bench test --check system --run TestCandidateInstall",
          "exit_code": 0
        },
        {
          "id": "t5r2-t5-digest-probe",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:f867fcf0caeaec248a007fbc914b86feb15f4ed51b8d240ac3e80a537a1f8922",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,180,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,180,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1159\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6398,180\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6449,180\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-digest-probe",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the digest comparison of the candidate check; TestCandidateCheck must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
              "digest": "sha256:f867fcf0caeaec248a007fbc914b86feb15f4ed51b8d240ac3e80a537a1f8922",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,180,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,180,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1159\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6398,180\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6449,180\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r2-t5-locator-probe",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:f867fcf0caeaec248a007fbc914b86feb15f4ed51b8d240ac3e80a537a1f8922",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,180,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,180,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1159\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6398,180\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6449,180\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-locator-probe",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "read the broker manifest beside the caller root instead of the primary checkout; TestCandidateCheck must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
              "digest": "sha256:f867fcf0caeaec248a007fbc914b86feb15f4ed51b8d240ac3e80a537a1f8922",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,180,1,0,null,yes\n  2,bit,swap,failed,1,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,180,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1159\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6398,180\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6449,180\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:156: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r2-t5-chunk-order-probe",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:5387ea4fa33a80a95c76a8cd5c824ae1260913490275451ade2a64d14b4f404d",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,2,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1348\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28920,508\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (807 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29494,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-chunk-order-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "move the candidate check of the chunk checkpoint after its execute path; TestCheckpointPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
              "digest": "sha256:5387ea4fa33a80a95c76a8cd5c824ae1260913490275451ade2a64d14b4f404d",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,2,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1348\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28920,508\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (807 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29494,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r2-t5-complete-order-probe",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:5387ea4fa33a80a95c76a8cd5c824ae1260913490275451ade2a64d14b4f404d",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,2,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1348\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28920,508\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (807 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29494,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-complete-order-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the candidate check of the complete checkpoint; TestCheckpointPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
              "digest": "sha256:5387ea4fa33a80a95c76a8cd5c824ae1260913490275451ade2a64d14b4f404d",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,2,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1348\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28920,508\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (807 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29494,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r2-t5-record-order-probe",
          "performer": "claude:ft419-t5-repair2-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
            "digest": "sha256:11ea3beec40ef2675c6890d74ebe0cb5863c8b96f55f21d0c6f4025ce623c265",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,4,TestAssignmentPersistedFields/an_undeclared_working_tree,260,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord/recordcmd,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1714\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,16878,260\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/an_undeclared_working_tree,\"assignment_persisted_test.go:86: assignment = exit 1, output \\\"error: bench record assignment refused — working specs/example/spec.md differs from its source; commit or discard the change first\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not-installed… (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/hostile_slug_and_flag_value,\"assignment_persisted_test.go:78: assignment = \\\"error: bench record assignment refused — missing or nonregular tree file specs/it's here/tickets/1.md\\\\\\\\n\\\", want one next= route that names no gate\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/malformed_declaration,\"assignment_persisted_test.go:92: assignment = (0, \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",1086d0a65c6dc29700f8cc4d5b42289094e5fd89\\\\\\\\n\\\"), want the persisted-fields-invalid error and… (249 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/uninstalled,\"assignment_persisted_test.go:63: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",7f228a0a691db3f9b87031268d4c8010cf31d5bf\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not… (251 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-record-order-probe",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the candidate check of the assignment entry; TestAssignmentPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair2-20261020",
              "digest": "sha256:11ea3beec40ef2675c6890d74ebe0cb5863c8b96f55f21d0c6f4025ce623c265",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,4,TestAssignmentPersistedFields/an_undeclared_working_tree,260,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord/recordcmd,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1714\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,16878,260\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/an_undeclared_working_tree,\"assignment_persisted_test.go:86: assignment = exit 1, output \\\"error: bench record assignment refused — working specs/example/spec.md differs from its source; commit or discard the change first\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not-installed… (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/hostile_slug_and_flag_value,\"assignment_persisted_test.go:78: assignment = \\\"error: bench record assignment refused — missing or nonregular tree file specs/it's here/tickets/1.md\\\\\\\\n\\\", want one next= route that names no gate\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/malformed_declaration,\"assignment_persisted_test.go:92: assignment = (0, \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",1086d0a65c6dc29700f8cc4d5b42289094e5fd89\\\\\\\\n\\\"), want the persisted-fields-invalid error and… (249 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/uninstalled,\"assignment_persisted_test.go:63: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",7f228a0a691db3f9b87031268d4c8010cf31d5bf\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not… (251 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r3-owner",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:dca1a2aebb0e00513879d782e6d304169793f6b60045bc2737f91b18e7deb02a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6523,181\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t5r3-manifest",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:a8ae0a8515973657d26dbd7a8f80c05ff01ad65ae10d8183358ad0c79e71228f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/brokermanifest,pass,2,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "manifest",
          "command": "bench test --package ./internal/brokermanifest",
          "exit_code": 0
        },
        {
          "id": "t5r3-gate",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:de6835c7b6c5a0fec4fc86ea8dabe0211eda7bd92deeba77588fd7f75a1816a9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,29059,508\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t5r3-record",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:bfeca1bc3677c72fe9ec7140aa4124b3cc7d6cdd70e8ad8a0d9abef6de6dc3d2",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,2844,110\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,16599,260\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "record",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0
        },
        {
          "id": "t5r3-routes",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:083201dbdeffa4beaddc7caec9539221723505998f67c6ee5a7c07652c31da24",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,22\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "routes",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "t5r3-cmd-bench",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:1c7234afa57d129b5e500b80cfc283abd6b44f7246ffe199faf94fe58cbe1c2b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14112,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t5r3-internal-conformance",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:2aac7fdd645dcd9b7bf171ffed8b8e610fc17d1b5e2e232bedc1bb2e476faaa0",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,33333,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "internal-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t5r3-commitmenttest",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:a6bd8772a7701bb2c44799879209bac9ed4990f39d308bf4f8adc700372e3ab1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/commitmenttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "commitmenttest",
          "command": "bench test --package ./internal/commitment/commitmenttest",
          "exit_code": 0
        },
        {
          "id": "t5r3-root-conformance",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:44ad4cea96dc48f111d0aad2eb222a58cd06972bb754ef64a52a99c5bd1376e7",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13603,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t5r3-system",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:acabf28faaa4ac9e00f94f082c82d51094681cda370b5609f34cbc0d4a8888df",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,3,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,34162,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "system",
          "command": "bench test --check system --run TestCandidateInstall",
          "exit_code": 0
        },
        {
          "id": "t5r3-t5-digest-probe",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:9853bc91a7896b8021a5cb3c8696ba8668df261114067d5d97a7ebb91ebe104a",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,181,1,0,null,yes\n  2,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,181,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1655\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6317,181\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_wrapper_does_not_move_the_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6191,181\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_wrapper_does_not_move_the_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-digest-probe",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the digest comparison of the candidate check; TestCandidateCheck must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
              "digest": "sha256:9853bc91a7896b8021a5cb3c8696ba8668df261114067d5d97a7ebb91ebe104a",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,181,1,0,null,yes\n  2,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,181,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1655\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6317,181\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_wrapper_does_not_move_the_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6191,181\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_wrapper_does_not_move_the_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r3-t5-locator-probe",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:9853bc91a7896b8021a5cb3c8696ba8668df261114067d5d97a7ebb91ebe104a",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,181,1,0,null,yes\n  2,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,181,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1655\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6317,181\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_wrapper_does_not_move_the_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6191,181\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_wrapper_does_not_move_the_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-locator-probe",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "read the broker manifest beside the caller root instead of the primary checkout; TestCandidateCheck must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
              "digest": "sha256:9853bc91a7896b8021a5cb3c8696ba8668df261114067d5d97a7ebb91ebe104a",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,181,1,0,null,yes\n  2,bit,swap,failed,2,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,181,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment/repository,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1655\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6317,181\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_wrapper_does_not_move_the_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/another_digest_refuses,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,6191,181\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_root_reads_the_primary_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCandidateCheck/a_worktree_wrapper_does_not_move_the_manifest,\"install_check_test.go:167: check = <nil>; want the refusal face \\\"candidate-not-installed\\\", or a pass when it is empty\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r3-t5-chunk-order-probe",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:62527866a421f426045bd8e21891efdc1a18b9e5b972bbac9ce8d9995bcebc7c",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,3,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1700\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29347,508\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (809 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29178,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"candidate-not-installed: the running executable \\\\\\\\\\\"/dev/shm/27UK6Q/t/go-build2504943068/b001/gate.test\\\\\\\\\\\" has digest e139bc6138c40af4b722ea8058276ee4d0b3cb7b2f17b1ef1dc26982fd0ed8f3, and… (669 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-chunk-order-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "move the candidate check of the chunk checkpoint after its execute path; TestCheckpointPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
              "digest": "sha256:62527866a421f426045bd8e21891efdc1a18b9e5b972bbac9ce8d9995bcebc7c",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,3,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1700\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29347,508\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (809 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29178,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"candidate-not-installed: the running executable \\\\\\\\\\\"/dev/shm/27UK6Q/t/go-build2504943068/b001/gate.test\\\\\\\\\\\" has digest e139bc6138c40af4b722ea8058276ee4d0b3cb7b2f17b1ef1dc26982fd0ed8f3, and… (669 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r3-t5-complete-order-probe",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:62527866a421f426045bd8e21891efdc1a18b9e5b972bbac9ce8d9995bcebc7c",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,3,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1700\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29347,508\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (809 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29178,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"candidate-not-installed: the running executable \\\\\\\\\\\"/dev/shm/27UK6Q/t/go-build2504943068/b001/gate.test\\\\\\\\\\\" has digest e139bc6138c40af4b722ea8058276ee4d0b3cb7b2f17b1ef1dc26982fd0ed8f3, and… (669 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-complete-order-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the candidate check of the complete checkpoint; TestCheckpointPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
              "digest": "sha256:62527866a421f426045bd8e21891efdc1a18b9e5b972bbac9ce8d9995bcebc7c",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,3,TestCheckpointPersistedFields/chunk,508,1,0,null,yes\n  2,bit,swap,failed,1,TestCheckpointPersistedFields/complete,508,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1700\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29347,508\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/chunk,\"checkpoint_persisted_test.go:149: checkpoint = (1, \\\"\\\", \\\"gate subject unavailable: completion evidence: chunk 1: missing Coverage; record the native review result\\\\\\\\nnext=bench record 'review' --help\\\\\\\\ncandidate-not-installed: the running execu… (809 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"\\\"), want a candidate-not-installed refusal with one route and no oracle output\",1\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/malformed_declaration,\"checkpoint_persisted_test.go:183: checkpoint = (1, \\\"gate subject unavailable: completion evidence: stale plan digest; review and record the plan delta\\\\\\\\nnext=bench record 'amendment' --help\\\\\\\\npersisted-fields-invalid: unknown record \\\\\\\\\\\"gate-ev… (477 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29178,508\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestCheckpointPersistedFields/complete,\"checkpoint_persisted_test.go:154: checkpoint = (0, \\\"\\\", \\\"candidate-not-installed: the running executable \\\\\\\\\\\"/dev/shm/27UK6Q/t/go-build2504943068/b001/gate.test\\\\\\\\\\\" has digest e139bc6138c40af4b722ea8058276ee4d0b3cb7b2f17b1ef1dc26982fd0ed8f3, and… (669 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5r3-t5-record-order-probe",
          "performer": "claude:ft419-t5-repair3-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
            "digest": "sha256:e44e492b8c603e456d274b76373b901fe1875367448a2056417b9aa6e587d86e",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,4,TestAssignmentPersistedFields/an_undeclared_working_tree,260,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord/recordcmd,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1714\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,16329,260\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/an_undeclared_working_tree,\"assignment_persisted_test.go:86: assignment = exit 1, output \\\"error: bench record assignment refused — working specs/example/spec.md differs from its source; commit or discard the change first\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not-installed… (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/hostile_slug_and_flag_value,\"assignment_persisted_test.go:78: assignment = \\\"error: bench record assignment refused — missing or nonregular tree file specs/it's here/tickets/1.md\\\\\\\\n\\\", want one next= route that names no gate\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/malformed_declaration,\"assignment_persisted_test.go:92: assignment = (0, \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",1c17b0e2e6f8f2ed072e6ecf193b982351f450b7\\\\\\\\n\\\"), want the persisted-fields-invalid error and… (249 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/uninstalled,\"assignment_persisted_test.go:63: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",0ff1f2d7da8e8e520a41dbf4ec52fb80f8da72c3\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not… (251 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-record-order-probe",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the candidate check of the assignment entry; TestAssignmentPersistedFields must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t5-repair3-20261020",
              "digest": "sha256:e44e492b8c603e456d274b76373b901fe1875367448a2056417b9aa6e587d86e",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,4,TestAssignmentPersistedFields/an_undeclared_working_tree,260,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord/recordcmd,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1714\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,16329,260\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/an_undeclared_working_tree,\"assignment_persisted_test.go:86: assignment = exit 1, output \\\"error: bench record assignment refused — working specs/example/spec.md differs from its source; commit or discard the change first\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not-installed… (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/hostile_slug_and_flag_value,\"assignment_persisted_test.go:78: assignment = \\\"error: bench record assignment refused — missing or nonregular tree file specs/it's here/tickets/1.md\\\\\\\\n\\\", want one next= route that names no gate\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/malformed_declaration,\"assignment_persisted_test.go:92: assignment = (0, \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",1c17b0e2e6f8f2ed072e6ecf193b982351f450b7\\\\\\\\n\\\"), want the persisted-fields-invalid error and… (249 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentPersistedFields/uninstalled,\"assignment_persisted_test.go:63: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",0ff1f2d7da8e8e520a41dbf4ec52fb80f8da72c3\\\\\\\\n\\\"; want exit 1 naming \\\"candidate-not… (251 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "cbc3-standards",
          "performer": "claude:ft419-cbc3-standards-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-standards-20261020",
            "digest": "sha256:de44da6c7fa4c13a3d5cebdea61cec713d31a160f3f9239114bc36c57696852c",
            "excerpt": "## Standards (CB-C3, pair ec4144c3..7b8ca752, opus/medium)\nF1 auto-fix conf 9: install_check_test.go:95 spells the build-inputs path that freshness.BuildInputsManifest owns.\nF2 auto-fix conf 7: window.go:97 spells \"bin\" again; one exported brokermanifest path helper owns it.\nF3 ask-user conf 8: the manifest-binding fixture repeats in three test packages; one commitmenttest helper owns it (widening).\nF4 auto-fix conf 5: GateRerun, RecordRerun, and CommitRerun repeat one caller-rerun step; one constructor and renderer own it.\nF5 auto-fix conf 6: checkpoint_persisted_test.go:26 comment does not parse.\nR1-R4 no-op (intent import, DigestField, flag.term split, GateRerun step).\nCount 5. Worst: F3.\n"
          },
          "axis": "Standards",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "b1385cae54d4ef6e33317580b734784813b4b489",
          "finding_ids": [
            "F1",
            "F2",
            "F3",
            "F4",
            "F5"
          ],
          "supersedes": []
        },
        {
          "id": "cbc3-spec",
          "performer": "claude:ft419-cbc3-spec-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-spec-20261020",
            "digest": "sha256:ff7f8b8b0127ff79843531f12857ddf61983a0478480b933a05be0c3522696bd",
            "excerpt": "## Spec (CB-C3, pair ec4144c3..7b8ca752, opus/medium)\nF1 auto-fix conf 6: CheckEntry (window.go:169-182) reads and parses the declaration before the build-inputs test (window.go:199), so a linked repository refuses on a malformed or unreadable spec; the spec says the declaration stays inert there.\nEvery CB-C3 row met and cited. Count 1. Worst: F1.\n"
          },
          "axis": "Spec",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "b1385cae54d4ef6e33317580b734784813b4b489",
          "finding_ids": [
            "F1"
          ],
          "supersedes": []
        },
        {
          "id": "cbc3-coverage",
          "performer": "claude:ft419-cbc3-coverage-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fddc4550375bf3ad450a3b035949b00dbfbfc7e0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-coverage-20261020",
            "digest": "sha256:3db54bbbd566071da3e10cc84b5bb3e317301709e1dfbfecd05af12dbffd4dba",
            "excerpt": "## Coverage (CB-C3, pair ec4144c3..7b8ca752, opus/medium)\nF1 auto-fix conf 8: no test has a working-tree spec that differs from HEAD; mutant os.ReadFile at window.go:170 silent, restored=yes.\nF2 auto-fix conf 7: no CB61 case with a hostile slug or flag value; hand-quote mutants at command.go:761 and :750 silent, restored=yes.\nF3 ask-user conf 5: with two open declared builds the route names the first build's target, even for a caller in the second.\nCount 3. Worst: F1.\n"
          },
          "axis": "Coverage",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "b1385cae54d4ef6e33317580b734784813b4b489",
          "finding_ids": [
            "F1",
            "F2",
            "F3"
          ],
          "supersedes": []
        },
        {
          "id": "cbc3-standards-confirm",
          "performer": "claude:ft419-cbc3-standards-confirm-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-standards-confirm-20261020",
            "digest": "sha256:8f63027d353c4cce693b46793d20c4fc0c0ff188f4b482add7207a5b0c31cd5b",
            "excerpt": "## Standards (CB-C3 confirming round, repair delta 7b8ca752..dd995111, opus/medium)\nF1 to F5 closed. The declared-builds-overlap face matches its siblings, and its Reviewer authority follows the CommitmentDecision precedent.\nN1 auto-fix conf 7: internal/worktree/build.go:289,295 still derive the manifest path beside brokermanifest.Path; routed to the ticket 6 charge, whose Writes line holds build.go.\nN2 auto-fix conf 7: window.go:74 doc comment is stale after the build-inputs reorder.\nCount 2 new.\n"
          },
          "axis": "Standards",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "9df8de96038e542a88bc0366e003376bc9ce904f",
          "finding_ids": [
            "N1",
            "N2"
          ],
          "supersedes": [
            "cbc3-standards"
          ]
        },
        {
          "id": "cbc3-spec-confirm",
          "performer": "claude:ft419-cbc3-spec-confirm-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-spec-confirm-20261020",
            "digest": "sha256:e4f059bc98515b1f82b0c4691df659faddb0bc08a9f251733a050a5d6b86dbd4",
            "excerpt": "## Spec (CB-C3 confirming round, repair delta 7b8ca752..dd995111, opus/medium)\nSpec F1, CB69 and its paragraph, DeclaredBuild first-of-list, and repair scope met. The --in disposition under the existing Won't handle line is sound.\nN1 auto-fix conf 5: the overlap count is skipped on any declared-build read error, but the spec scopes the skip to an unreadable ledger; prose correction.\nN2 auto-fix conf 4: CB69 traces to story 44, which names the target, not the overlap refusal; prose correction adds a story.\nCount 2 new.\n"
          },
          "axis": "Spec",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "9df8de96038e542a88bc0366e003376bc9ce904f",
          "finding_ids": [
            "N1",
            "N2"
          ],
          "supersedes": [
            "cbc3-spec"
          ]
        },
        {
          "id": "cbc3-coverage-confirm",
          "performer": "claude:ft419-cbc3-coverage-confirm-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8bd178419cbbfddaa3622986f9f75001aca1692d",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-coverage-confirm-20261020",
            "digest": "sha256:08b003b73a0f026444e9756581630389f69f77c62acb38745e01498e91acf8a1",
            "excerpt": "## Coverage (CB-C3 confirming round, repair delta 7b8ca752..dd995111, opus/medium)\nF1, F2, F3/CB69, and Spec F1 closed; each independent mutant bit, restored=yes.\nN1 auto-fix conf 7: a wrong location in brokermanifest.Path stays silent over unit, adopt, cmd, and system runs, because the fixtures bind through the same Path. Add a case that binds through adopt.WriteBrokerManifest and expects the check to pass.\nCount 1 new.\n"
          },
          "axis": "Coverage",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "9df8de96038e542a88bc0366e003376bc9ce904f",
          "finding_ids": [
            "N1"
          ],
          "supersedes": [
            "cbc3-coverage"
          ]
        },
        {
          "id": "cbc3-standards-r2",
          "performer": "claude:ft419-cbc3-standards-r2-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-standards-r2-20261020",
            "digest": "sha256:89ee5abfbf24c919209e5fcc59b4859629dc2097ca1e568431f1e27e2598e5be",
            "excerpt": "## Standards (CB-C3 repair cycle 2, delta dd995111..7972fc64, opus/medium)\nN2 closed: window.go:74-75 now scopes the refusals to a repository that declares Bench build inputs.\nNo finding ID. S2 evidence: the new case's independent bin expectation has two demonstrated reds, the author's swap of \"bin\" to \"dist\" in brokermanifest.Path and the orchestrator's move of the writer dest in internal/adopt/broker.go; both bit with restored=yes.\nAdvice routed to the ticket 6 charge: the window.go:74 comment omits the unresolved-primary case.\n"
          },
          "axis": "Standards",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "d0c62cc4413621aff8e8124a65b9f72a069555d2",
          "finding_ids": [],
          "supersedes": [
            "cbc3-standards-confirm"
          ]
        },
        {
          "id": "cbc3-spec-r2",
          "performer": "claude:ft419-cbc3-spec-r2-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-spec-r2-20261020",
            "digest": "sha256:57e899ca1a2b7659368e0ca41f59d45c1d0d93a2f177ca9868511c6f01cd0214",
            "excerpt": "## Spec (CB-C3 repair cycle 2, delta dd995111..7972fc64, opus/medium)\nSpec N1, N2, the chunk test list, and the new case met; the diff adds no behavior.\nAdvice routed to the ticket 6 charge: window.go:117-118 still says a ledger that does not read lists no build, narrower than the code.\n"
          },
          "axis": "Spec",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "d0c62cc4413621aff8e8124a65b9f72a069555d2",
          "finding_ids": [],
          "supersedes": [
            "cbc3-spec-confirm"
          ]
        },
        {
          "id": "cbc3-coverage-r2",
          "performer": "claude:ft419-cbc3-coverage-r2-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d67adea4d1da24470d407189925f5d5ca7a5fec7",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-coverage-r2-20261020",
            "digest": "sha256:67fd7d50461b45ab855ea1c66fe48f926c95c691d6e27834e4e780951e599fcf",
            "excerpt": "## Coverage (CB-C3 repair cycle 2, delta dd995111..7972fc64, opus/medium)\nN1 closed: the production-writer case refuses on a moved reader or writer location.\nC1 auto-fix conf 7: a check that reads the manifest beside BENCH_WRAPPER stays silent (14 tests, restored=yes). No row sets BENCH_WRAPPER to a worktree bin/ with the worktree build as the executable. CB60's catch clause names that bypass.\nCount 1.\n"
          },
          "axis": "Coverage",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "d0c62cc4413621aff8e8124a65b9f72a069555d2",
          "finding_ids": [
            "C1"
          ],
          "supersedes": [
            "cbc3-coverage-confirm"
          ]
        },
        {
          "id": "cbc3-standards-r3",
          "performer": "claude:ft419-cbc3-standards-r3-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-standards-r3-20261020",
            "digest": "sha256:669b08ac22306c0451daee0d82a1340b12406109b1cc6574b579e81b9f8648d2",
            "excerpt": "## Standards (CB-C3 repair cycle 3, delta 7972fc64..e1f2dfdd, opus/medium)\nNo finding ID. Evidence for the new row's independent expectation: three demonstrated reds, each restored=yes. The author's insertion of a wrapper-relative read before brokermanifest.Read in window.go bit it during the debug loop. The orchestrator's swap of manifest := brokermanifest.Path(primary) with a wrapper-suffix read bit it. The orchestrator's swap of \"bin\" to \"dist\" in brokermanifest.Path bit it.\nwrapperAt has one source, and the comments meet craft-comments.\n"
          },
          "axis": "Standards",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "cee87305b25de0e474628722a3936eed0c3e41a5",
          "finding_ids": [],
          "supersedes": [
            "cbc3-standards-r2"
          ]
        },
        {
          "id": "cbc3-spec-r3",
          "performer": "claude:ft419-cbc3-spec-r3-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-spec-r3-20261020",
            "digest": "sha256:7dc24b66b3192b5acdc9915480d78474e7ab4ba5cbeb5f89e68d58e814713d3e",
            "excerpt": "## Spec (CB-C3 repair cycle 3, delta 7972fc64..e1f2dfdd, opus/medium)\nCB60 met, the clause \"never beside BENCH_WRAPPER\" met, the plan change only appends the cycle 3 assignment, and the delta adds no behavior. No finding.\n"
          },
          "axis": "Spec",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "cee87305b25de0e474628722a3936eed0c3e41a5",
          "finding_ids": [],
          "supersedes": [
            "cbc3-spec-r2"
          ]
        },
        {
          "id": "cbc3-coverage-r3",
          "performer": "claude:ft419-cbc3-coverage-r3-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "45ad78a20b7966f109b621de105c1034bde1fbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc3-coverage-r3-20261020",
            "digest": "sha256:77952ae35bdf8e1bca229fa1ad465a0205fedb8ad48ddca73f87c78887b71b51",
            "excerpt": "## Coverage (CB-C3 repair cycle 3, delta 7972fc64..e1f2dfdd, opus/medium)\nC1 closed: a read beside the caller root when BENCH_WRAPPER is set bit the new row, restored=yes.\nAdvice (no finding ID; the bounded repair policy makes a finding that grades only the extension's check optional advice, and the production code never reads beside BENCH_WRAPPER): a read beside BENCH_WRAPPER only when the root is the primary checkout stays silent (15 tests, restored=yes).\n"
          },
          "axis": "Coverage",
          "base": "ec4144c30be12b9b823b66ee38c29ee83c9d8fdb",
          "tip": "cee87305b25de0e474628722a3936eed0c3e41a5",
          "finding_ids": [],
          "supersedes": [
            "cbc3-coverage-r2"
          ]
        }
      ]
    },
    {
      "id": "CB-C4",
      "base": "cee87305b25de0e474628722a3936eed0c3e41a5",
      "tip": "6c74df34ab340683bfb1ed07ef1fbcdf4d4a2642",
      "plan_digest": "sha256:9aada2a977c8775e0d4132cd56d98e965741fef22865668ffad99bc3a577a390",
      "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
      "acceptance_rows": [
        "CB37",
        "CB38",
        "CB39",
        "CB40",
        "CB41",
        "CB42",
        "CB43",
        "CB49",
        "CB64"
      ],
      "verification": [
        {
          "id": "t6-worktree",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:ed9d46e8d3b7c6c1eddde94efcdffdca3ca88d2c1096911550414c6a2d20e4f1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,66240,1487\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t6-cmd-bench",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:3c6d7d8f313579eda152b91125a3ac72b97060b52272f0f2f5fe5398dad34431",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14562,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t6-internal-commitment-repository",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:1b4d4abce6156ce5fcff8bc0b3528f5c0dd445c94b54203becbebe4ecd845c64",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6328,181\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "internal-commitment-repository",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "t6-internal-conformance",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:7763ca5c76dac1952214180dd1b0e9f3b94681c7adf27eb06dbdabb68680934a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,33230,1116\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "internal-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t6-internal-intent",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:ffa2d26ccc236f5c2714a95c3bb902472f284f81dcfe9c7a101afb1283dfe6a5",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,pass,1097,79\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "internal-intent",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0
        },
        {
          "id": "t6-internal-usage",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:9dc3a2106ba5d02c9c72884fc4cdb6f759767003b2ad7ad20a784047d4059575",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/usage,pass,2,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "internal-usage",
          "command": "bench test --package ./internal/usage",
          "exit_code": 0
        },
        {
          "id": "t6-spec",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:cd06fa53765b00758e292a25b80060ed9155667c608251d5524d34cd506b7032",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,pass,1734,175\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "spec",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0
        },
        {
          "id": "t6-refusalroute",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:083201dbdeffa4beaddc7caec9539221723505998f67c6ee5a7c07652c31da24",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,22\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "refusalroute",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "t6-root-conformance",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:529a23dcabf3b1f960d7803e4d1c060e3bf55d323e0f944fce5717005d5a0399",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13792,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t6-t6-order-probe",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:77bb1f7b73b88637942dd679764d32bbeff387d82d389f01fbc5741e691163f4",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,3,TestBuildRestore/refuses_before_any_write/a_ledger_that_is_not_JSON,1487,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1387\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,69628,1487\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildRestore/refuses_before_any_write/a_ledger_that_is_not_JSON,\"build_restore_test.go:204: restore = (1, \\\"bench worktree build: read intent ledger: malformed JSON: decode persisted JSON: EOF\\\\\\\\nnext=bench worktree list\\\\\\\\n\\\"), want exit 1 naming \\\"next=bench worktree build --install <target>; then bench workt… (270 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestBuildRestore/refuses_before_any_write/a_ledger_that_the_strict_read_refuses,\"build_restore_test.go:204: restore = (1, \\\"bench worktree build: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench worktree list\\\\\\\\n\\\"), want exit 1 naming \\\"cannot read the intent ledger\\\"\",1\n  github.com/gibbonmi/bench/internal/worktree,TestBuildRestore/refuses_before_any_write/a_refused_ledger_whose_target_tree_is_gone,\"build_restore_test.go:204: restore = (1, \\\"bench worktree build: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench worktree list\\\\\\\\n\\\"), want exit 1 naming \\\"next=bench worktree build --instal… (298 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-order-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "run the open-build test before the strict ledger read in the restore; TestBuildRestore must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
              "digest": "sha256:77bb1f7b73b88637942dd679764d32bbeff387d82d389f01fbc5741e691163f4",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,3,TestBuildRestore/refuses_before_any_write/a_ledger_that_is_not_JSON,1487,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1387\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,69628,1487\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildRestore/refuses_before_any_write/a_ledger_that_is_not_JSON,\"build_restore_test.go:204: restore = (1, \\\"bench worktree build: read intent ledger: malformed JSON: decode persisted JSON: EOF\\\\\\\\nnext=bench worktree list\\\\\\\\n\\\"), want exit 1 naming \\\"next=bench worktree build --install <target>; then bench workt… (270 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestBuildRestore/refuses_before_any_write/a_ledger_that_the_strict_read_refuses,\"build_restore_test.go:204: restore = (1, \\\"bench worktree build: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench worktree list\\\\\\\\n\\\"), want exit 1 naming \\\"cannot read the intent ledger\\\"\",1\n  github.com/gibbonmi/bench/internal/worktree,TestBuildRestore/refuses_before_any_write/a_refused_ledger_whose_target_tree_is_gone,\"build_restore_test.go:204: restore = (1, \\\"bench worktree build: read intent ledger: malformed JSON: decode persisted JSON: json: unknown field \\\\\\\\\\\"fixture\\\\\\\\\\\"\\\\\\\\nnext=bench worktree list\\\\\\\\n\\\"), want exit 1 naming \\\"next=bench worktree build --instal… (298 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t6-t6-nested-probe",
          "performer": "claude:ft419-t6-author-20261020",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
            "digest": "sha256:58462a7c5ccf51307d4ff619771b8698e08984971eabded33735d52108b802b3",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestBuildRestore/a_declared_field_that_the_ledger_does_not_hold_needs_no_change,1487,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1361\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,67439,1487\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildRestore/a_declared_field_that_the_ledger_does_not_hold_needs_no_change,\"build_restore_test.go:155: ledger = {\\\"assignments\\\":[{\\\"branch\\\":\\\"refs/heads/bench/assign/d3f0fe6e8404202530aa19578188e1b9/82552e9fa0dfeac5b27ae75f1d7d1e0f\\\",\\\"created_at\\\":\\\"2026-10-20T22:06:26Z\\\",\\\"id\\\":\\\"82552e9fa0dfeac5b27ae75f1d7d1e0f\\\",\\\"label\\\":\\\"r… (774 bytes)\",3\n  github.com/gibbonmi/bench/internal/worktree,TestBuildRestore/clears_a_nested_field_from_each_element,\"build_restore_test.go:143: ledger holds no commitment object: map[assignments:[map[branch:refs/heads/bench/assign/fc17310c7126b1c788d4e6ec914bc4ad/599bff6ed9a0496bd0c4fba3dac7a5ec created_at:2026-10-20T22:06:24Z id:599bff6ed9a0496bd0c4fba3d… (734 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,\"TestBuildRestore/copies_the_ledger,_clears_the_field,_rebuilds,_and_repairs\",\"build_restore_test.go:122: ledger holds no commitment object: map[assignments:[map[branch:refs/heads/bench/assign/d67fd841972013eb987ab4aaac1217b3/d0b93df1d5bde2175079d0ca896edfce created_at:2026-10-20T22:06:26Z id:d0b93df1d5bde2175079d0ca8… (748 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-nested-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "delete only the top-level key in the raw-field removal; TestBuildRestore must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft419-t6-author-20261020",
              "digest": "sha256:58462a7c5ccf51307d4ff619771b8698e08984971eabded33735d52108b802b3",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestBuildRestore/a_declared_field_that_the_ledger_does_not_hold_needs_no_change,1487,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1361\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,67439,1487\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestBuildRestore/a_declared_field_that_the_ledger_does_not_hold_needs_no_change,\"build_restore_test.go:155: ledger = {\\\"assignments\\\":[{\\\"branch\\\":\\\"refs/heads/bench/assign/d3f0fe6e8404202530aa19578188e1b9/82552e9fa0dfeac5b27ae75f1d7d1e0f\\\",\\\"created_at\\\":\\\"2026-10-20T22:06:26Z\\\",\\\"id\\\":\\\"82552e9fa0dfeac5b27ae75f1d7d1e0f\\\",\\\"label\\\":\\\"r… (774 bytes)\",3\n  github.com/gibbonmi/bench/internal/worktree,TestBuildRestore/clears_a_nested_field_from_each_element,\"build_restore_test.go:143: ledger holds no commitment object: map[assignments:[map[branch:refs/heads/bench/assign/fc17310c7126b1c788d4e6ec914bc4ad/599bff6ed9a0496bd0c4fba3dac7a5ec created_at:2026-10-20T22:06:24Z id:599bff6ed9a0496bd0c4fba3d… (734 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,\"TestBuildRestore/copies_the_ledger,_clears_the_field,_rebuilds,_and_repairs\",\"build_restore_test.go:122: ledger holds no commitment object: map[assignments:[map[branch:refs/heads/bench/assign/d67fd841972013eb987ab4aaac1217b3/d0b93df1d5bde2175079d0ca896edfce created_at:2026-10-20T22:06:26Z id:d0b93df1d5bde2175079d0ca8… (748 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "cbc4-standards",
          "performer": "claude:ft419-cbc4-standards-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc4-standards-20261020",
            "digest": "sha256:9afeba2583741b6262c0096ba2816510c3bcfb17a30915b6134e2fc625bc20d2",
            "excerpt": "## Standards (CB-C4, pair cee87305..cca32bd6, opus/medium)\nF1 ask-user conf 6: BlockRoute (faces_gate.go:39-41) renders the shared blockOutcome step with no Reviewer authority, while the overlap face marks the same step Reviewer.\nF2 auto-fix conf 6: build_restore.go:77 builds the install-then-rerun route by hand with StepJoiner, which only the registry spells.\nF3 auto-fix conf 6: routeWord, targetSlot, slugSlot (build.go:362-375) copy the registry quoting rule and the usage placeholders.\nF4 auto-fix conf 6: window.go:119-121 comment does not match the located != nil branch.\nF5 auto-fix conf 4: window.go:146-147 comment keeps the narrowness the ticket corrected.\nCount 5. Worst: F1.\n"
          },
          "axis": "Standards",
          "base": "cee87305b25de0e474628722a3936eed0c3e41a5",
          "tip": "6c74df34ab340683bfb1ed07ef1fbcdf4d4a2642",
          "finding_ids": [
            "F1",
            "F2",
            "F3",
            "F4",
            "F5"
          ],
          "supersedes": []
        },
        {
          "id": "cbc4-spec",
          "performer": "claude:ft419-cbc4-spec-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc4-spec-20261020",
            "digest": "sha256:f57ccc5f3cbd36b64b90ae6c600356aa1649d3fba2edd41014e1f283755d47f7",
            "excerpt": "## Spec (CB-C4, pair cee87305..cca32bd6, opus/medium)\nNo finding. CB37-CB43, CB49, CB64 met; the copy folded into the removal and the decode-only check meet the spec's strict decode; refusal order, routes, rebuild owner, child environment, and carried items met.\nAdvice: spec line 198 attributes the copy to the verb; the ledger owner takes it under the lock.\n"
          },
          "axis": "Spec",
          "base": "cee87305b25de0e474628722a3936eed0c3e41a5",
          "tip": "6c74df34ab340683bfb1ed07ef1fbcdf4d4a2642",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cbc4-coverage",
          "performer": "claude:ft419-cbc4-coverage-20261020",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "30a32e8ce625a8375416a4670fa987de49a8198f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft419-cbc4-coverage-20261020",
            "digest": "sha256:38083a6256c886a0238acf1000db7ea7a6362f9ff1b50508d79b4e0a67bf6442",
            "excerpt": "## Coverage (CB-C4, pair cee87305..cca32bd6, opus/medium)\nF1 auto-fix conf 9: no test has an open build of an undeclared outcome; DeclaredBuild in place of HoldsOpenBuild stays silent, restored=yes.\nF2 auto-fix conf 9: no test declares several ledger paths; a short-circuit removal loop stays silent, restored=yes.\nF3 auto-fix conf 8: no test runs the removal while another holder has the ledger lock; a no-op lock stays silent, restored=yes.\nF4 auto-fix conf 8: no test declares a review-record path; a filter that accepts any record stays silent, restored=yes.\nF5 auto-fix conf 6: with no ledger file, RemoveFields fails on its read and the restore refuses with an endless rerun; CB49's rule says it should complete.\nCount 5. Worst: F1.\n"
          },
          "axis": "Coverage",
          "base": "cee87305b25de0e474628722a3936eed0c3e41a5",
          "tip": "6c74df34ab340683bfb1ed07ef1fbcdf4d4a2642",
          "finding_ids": [
            "F1",
            "F2",
            "F3",
            "F4",
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
    },
    {
      "from": "sha256:11740dde07e548048c06e174cbae4b13dbf2f9c97ab28592da7ffce6ed524e6e",
      "to": "sha256:90738f167f2b1ec19be414edfd1c6feb7cfeee86402486076859080b943cc922",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ]
      }
    },
    {
      "from": "sha256:90738f167f2b1ec19be414edfd1c6feb7cfeee86402486076859080b943cc922",
      "to": "sha256:06abc74fade0bc0761e3e4b1870f37989e0514d809a045a87af2e8bad9f3280d",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ]
      }
    },
    {
      "from": "sha256:06abc74fade0bc0761e3e4b1870f37989e0514d809a045a87af2e8bad9f3280d",
      "to": "sha256:6f620eb19f2ea53404fa00df9f10568f30e5ca3ba5527b63f8505233223eb38a",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ]
      }
    },
    {
      "from": "sha256:6f620eb19f2ea53404fa00df9f10568f30e5ca3ba5527b63f8505233223eb38a",
      "to": "sha256:d2bb9d5cac795e9b265eee84ebbd2f679ac9d05edf353343efe2b6c43646f811",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ]
      }
    },
    {
      "from": "sha256:d2bb9d5cac795e9b265eee84ebbd2f679ac9d05edf353343efe2b6c43646f811",
      "to": "sha256:8451fc1cd35152b165098f00691638091485dab2e16dd21032046eb32e8d79a6",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ]
      }
    },
    {
      "from": "sha256:8451fc1cd35152b165098f00691638091485dab2e16dd21032046eb32e8d79a6",
      "to": "sha256:0cdd9a879f1b8bed408506623917fe20b0c17bacdca3f2da193684e477f8f07d",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ]
      }
    },
    {
      "from": "sha256:0cdd9a879f1b8bed408506623917fe20b0c17bacdca3f2da193684e477f8f07d",
      "to": "sha256:2bfce79ff100ea966faba7279d981e406958c557ab6c1343c422074803a82e13",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ],
        "CB-C3": [
          "CB-C3"
        ]
      }
    },
    {
      "from": "sha256:2bfce79ff100ea966faba7279d981e406958c557ab6c1343c422074803a82e13",
      "to": "sha256:3ce3a54891bf8f86892e152a2582ea4b9edcee1cb93ad581350fa513e330870f",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ],
        "CB-C3": [
          "CB-C3"
        ]
      }
    },
    {
      "from": "sha256:3ce3a54891bf8f86892e152a2582ea4b9edcee1cb93ad581350fa513e330870f",
      "to": "sha256:0fc5c6f43017b69a38f2d778727c2e9ed18ce158d88918252ae1ebf85f92d610",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ],
        "CB-C3": [
          "CB-C3"
        ]
      }
    },
    {
      "from": "sha256:0fc5c6f43017b69a38f2d778727c2e9ed18ce158d88918252ae1ebf85f92d610",
      "to": "sha256:2938a66ff05b5a6ab30332163ab39e2477b3f2c39a0989f089687b0f0ad93bb5",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ],
        "CB-C3": [
          "CB-C3"
        ]
      }
    },
    {
      "from": "sha256:2938a66ff05b5a6ab30332163ab39e2477b3f2c39a0989f089687b0f0ad93bb5",
      "to": "sha256:0c508811425128eb96494093745f4074bfedd49a4b0941685db08436845bf193",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ],
        "CB-C3": [
          "CB-C3"
        ]
      }
    },
    {
      "from": "sha256:0c508811425128eb96494093745f4074bfedd49a4b0941685db08436845bf193",
      "to": "sha256:a89a46a8cd89fac2ba14e0b4e4a22376282c6b7fd0fb1a2df12ff2fc13a967fc",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ],
        "CB-C3": [
          "CB-C3"
        ]
      }
    },
    {
      "from": "sha256:a89a46a8cd89fac2ba14e0b4e4a22376282c6b7fd0fb1a2df12ff2fc13a967fc",
      "to": "sha256:3ede70a4f3cbb8792467937b82b5f3f091ad3e9dcc0a18fc953f66a4521adeae",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ],
        "CB-C3": [
          "CB-C3"
        ]
      }
    },
    {
      "from": "sha256:3ede70a4f3cbb8792467937b82b5f3f091ad3e9dcc0a18fc953f66a4521adeae",
      "to": "sha256:9aada2a977c8775e0d4132cd56d98e965741fef22865668ffad99bc3a577a390",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ],
        "CB-C3": [
          "CB-C3"
        ]
      }
    },
    {
      "from": "sha256:9aada2a977c8775e0d4132cd56d98e965741fef22865668ffad99bc3a577a390",
      "to": "sha256:a043f0448b7a5ecaf95aae6c6be7cfda8f21961bcd4df1281a437248f8789f86",
      "chunk_ids": {
        "CB-C1": [
          "CB-C1"
        ],
        "CB-C2": [
          "CB-C2"
        ],
        "CB-C3": [
          "CB-C3"
        ],
        "CB-C4": [
          "CB-C4"
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

CB-C2: 9 findings. Worst: S3. Repair owners: ticket 2 takes S1, S3, S4, and S5. Ticket 3 takes S2, S6, S7, and S8.

- S1 (auto-fix under the run's --auto-approve widening, confidence 8): the raw ledger-field plant has two copies. One helper in `internal/commitment/commitmenttest` owns it.
- S2 (auto-fix, confidence 7): the landable spec text repeats in `internal/systemtest/candidate_install_test.go`. One `landableSpec` helper owns it.
- S3 (auto-fix, confidence 7): `installAssignment` repeats `resolveAssignmentIn` and drops the missing-tree refusal. One shared resolver body owns it.
- S4 (auto-fix, confidence 5): three sites compose the broker manifest path. One `manifestPath` helper owns it.
- S5 (auto-fix, confidence 6): `build.go` and `exec.go` each spell the wrapper path. One `wrapperPath` helper owns it.
- S6, S7, S8 (auto-fix, confidence 5 to 7): three comments name a missing restore, argue their own design, or do not parse.
- S9 (no-op, confidence 8): ticket 3 asks for `HoldsOpenBuild`, and ticket 6 consumes it.

CB-C3: 5 findings. Worst: F3. The ticket 5 repair takes all five.

- F1 (auto-fix, confidence 9): a test spells the build-inputs path that `freshness.BuildInputsManifest` owns.
- F2 (auto-fix, confidence 7): `window.go` spells the `bin` directory again. One exported path helper in `brokermanifest` owns it.
- F3 (auto-fix under the run's --auto-approve widening, confidence 8): the manifest-binding fixture repeats in three test packages. One `commitmenttest` helper owns it.
- F4 (auto-fix, confidence 5): three caller-rerun steps repeat one shape. One constructor and one renderer own it.
- F5 (auto-fix, confidence 6): a test comment in `checkpoint_persisted_test.go` does not parse.

CB-C4: 5 findings. Worst: F1. The ticket 6 repair takes all five.

- F1 (auto-fix after decision, confidence 6): `BlockRoute` renders the shared `bench commitment block` step with no Reviewer authority. The CB-C3 review set that step to Reviewer, so the restore route takes Reviewer too.
- F2 (auto-fix, confidence 6): the restore builds its install-then-rerun route by hand. The registry owns that route shape.
- F3 (auto-fix, confidence 6): `routeWord`, `targetSlot`, and `slugSlot` copy the registry quoting rule and the usage placeholders.
- F4 and F5 (auto-fix, confidence 6 and 4): two `window.go` comments still do not match the code.

## Spec

CB-C1: 0 findings. Worst: none.

CB-C2: 1 finding. Worst: F1.

- F1 (auto-fix, confidence 7): CB56 is partial at its seam, because its case starts with no earlier manifest. The ticket 2 repair cites the owner test in the CB56 cell.

CB-C3: 1 finding. Worst: F1.

- F1 (auto-fix, confidence 6): `CheckEntry` reads and parses the declaration before its build-inputs test. A linked repository then refuses on a malformed or unreadable spec, but the spec keeps the declaration inert there.

CB-C4: 0 findings. Advice: a prose correction says that the ledger owner takes the copy under the lock.

## Coverage

CB-C1: 4 findings. Worst: C1.

- C1 (auto-fix after consultation, confidence 8): a path segment with whitespace or a colon parses (`internal/spec/spec.go:126-134`). The opus high and sonnet high consultants agreed on a refusal with `persisted-fields-invalid`; new row CB66.
- C2 (auto-fix, confidence 7): CB6 bites only an interior `..`; add an empty path, a leading dot, and a trailing dot (`internal/spec/persisted_test.go:36`).
- C3 (auto-fix, confidence 7): CB5 bites only `gate-evidence`; add a second unknown record name (`internal/spec/persisted_test.go:35`).
- C4 (auto-fix after consultation, confidence 7): a second unfenced `Persisted fields:` line replaces the first (`internal/spec/spec.go:115-119`). The consultants agreed on a refusal of a second line with `persisted-fields-invalid`; new row CB67.

Confirming round (repair delta 90c9f2f75..e479d3763): S1, S2, C1, C2, C3, and C4 are closed.

- N1 (auto-fix, confidence 7): a space-only mutant at `internal/spec/spec.go:150` stays silent, because the CB66 refusal list holds only a plain space. Add a tab case and a U+00A0 case (`internal/spec/persisted_test.go:27-40`).
- N2 (no-op, confidence 4): an identical second line is untested, but the code refuses every second line, so the CB67 harm cannot occur.

Repair state: CB-C1 repair cycle 1 of 2 closed S1, S2, and C1 to C4. Cycle 2 of 2 closed N1 with a test-only change. The repair-assignment plan commit put the spec in the delta, so the three axes reviewed it again, and each passed. The `--auto-approve` option of the run approves the CB66 and CB67 expansion.

Advice (CB-C1 cycle 2, no finding ID): the U+00A0 row at `internal/spec/persisted_test.go:37` uses an invisible literal. The three whitespace rows share one subtest base name. The bounded repair policy makes these optional advice, because they grade only the N1 hardening check.

CB-C2: 3 findings. Worst: F1. Ticket 3 takes F1, and ticket 2 takes F2 and F3.

- F1 (auto-fix, confidence 8): no test has two open builds. New row CB68 pins the hold when only the second build declares.
- F2 (auto-fix, confidence 6): no test reaches the running executable through a symlinked alias of the primary path.
- F3 (auto-fix, confidence 6): no test seeds an earlier manifest before a failed publication.

Repair state: CB-C2 repair cycle 1 of 2 closed every accepted finding. The confirming round passed on all three axes.

Advice (CB-C2 confirming round, no finding ID): no row has the declared build first and an undeclared build second. The bounded repair policy makes this optional advice, because it grades only the CB68 hardening check.

CB-C3: 3 findings. Worst: F1. The ticket 5 repair takes all three.

- F1 (auto-fix, confidence 8): no test has a working-tree spec that differs from HEAD, so a working-tree read stays silent.
- F2 (auto-fix, confidence 7): no CB61 case has a hostile slug or flag value, so a hand-quoted rerun stays silent.
- F3 (auto-fix after consultation, confidence 5): with two open declared builds, the route names the first build's target. The opus high and sonnet high consultants agreed on a refusal with its own face. New row CB69 pins it.

Repair state: CB-C3 repair cycle 1 of 2 closed all nine findings. The confirming round accepted four new items, and cycle 2 of 2 takes them:

- Coverage N1 (auto-fix, confidence 7): no test binds the manifest through a production writer. A wrong location in `brokermanifest.Path` stays silent. The ticket 5 repair adds that case.
- Standards N2 (auto-fix, confidence 7): a `window.go` doc comment is stale after the reorder. The ticket 5 repair fixes it.
- Spec N1 and N2 (auto-fix, confidence 5 and 4): two prose-only spec corrections go in the cycle 2 plan commit.
- Standards N1 (auto-fix, confidence 7): `internal/worktree/build.go` still derives the manifest path. The ticket 6 charge takes it, because its Writes line holds that file.

Cycle 2 of 2 closed Coverage N1, Standards N2, and Spec N1 and N2. Its round raised one blocking item:

- C1 (auto-fix, confidence 7): a check that reads the manifest beside `BENCH_WRAPPER` passes every test. Row CB60 names that bypass, so CB60 is only partly met.

Repair extension: the bounded allowance is spent. The run's --auto-approve option is the reviewer's advance approval of repair rounds past the allowance. It permits one test-only cycle 3 for C1, with one `TestCandidateCheck` row. The extension adds a repair round to the implementation retro.

Advice for the ticket 6 charge: two `window.go` comments are narrower than the code.

Cycle 3 closed C1. The round passed on all three axes. Three demonstrated reds, each with a restored subject, record the new row's independent expectation.

Advice (CB-C3 cycle 3, no finding ID): a read beside `BENCH_WRAPPER` only when the root is the primary checkout stays silent. The production code never reads beside `BENCH_WRAPPER`. The bounded repair policy makes this optional advice, because it grades only the extension's check.

CB-C4: 5 findings. Worst: F1. The ticket 6 repair takes all five.

- F1 (auto-fix, confidence 9): no test has an open build of an undeclared outcome. A swap to the declared-build answer stays silent.
- F2 (auto-fix, confidence 9): no test declares several ledger paths, so a short-circuit removal loop stays silent.
- F3 (auto-fix, confidence 8): no test runs the removal while another holder has the ledger lock.
- F4 (auto-fix, confidence 8): no test declares a review-record path, so a filter that accepts any record stays silent.
- F5 (auto-fix after bench-debug, confidence 6): with no ledger file, the restore refuses with a rerun that cannot succeed. The rule of CB49 says it completes.

Repair state: CB-C4 repair cycle 1 of 2 runs one ticket 6 session.
