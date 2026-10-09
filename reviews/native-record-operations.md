# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/native-record-operations/spec.md",
  "plan_digest": "sha256:db85f3d092897f31ec6d6d341cb0a8c11a741dcc7d6587808e14a5945a113680",
  "implementation_session": "",
  "chunks": [
    {
      "id": "NR-C1",
      "base": "b925d1c1fd675b0cc299543ff3ba4a60a4ddf978",
      "tip": "e2068642f17fb744698626a66a26c6e67a7ce157",
      "plan_digest": "sha256:db85f3d092897f31ec6d6d341cb0a8c11a741dcc7d6587808e14a5945a113680",
      "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
      "acceptance_rows": [
        "NR1",
        "NR2",
        "NR3",
        "NR4",
        "NR5",
        "NR6",
        "NR7",
        "NR8",
        "NR9",
        "NR10",
        "NR11",
        "NR12",
        "NR13",
        "NR14",
        "NR36"
      ],
      "verification": [
        {
          "id": "nr-c1-t1-record-domain",
          "performer": "claude:nro_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:nro_t1_r1 bench test --package ./internal/reviewrecord at e2068642",
            "digest": "sha256:916eac24401b8e0a62ed9659a086a13890f563cd9272eae7cc425effe05754a5",
            "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,1561,73\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "record-domain",
          "command": "bench test --package ./internal/reviewrecord",
          "exit_code": 0
        },
        {
          "id": "nr-c1-t1-record-forms",
          "performer": "claude:nro_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:nro_t1_r1 bench test --package ./internal/reviewrecord/recordcmd at e2068642",
            "digest": "sha256:57a5c3b8fab2d3aab65b8b62c52900f04d2dc85433867973d92a3b747647a16a",
            "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,8425,141\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "record-forms",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "nr-c1-t1-command-package",
          "performer": "claude:nro_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:nro_t1_r1 bench test --package ./cmd/bench at e2068642",
            "digest": "sha256:b76268cf4a22f00b8849db64aa2b6da14c9bec3f4b8dbea2adaacdbfbbd4f597",
            "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14069,334\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "nr-c1-t1-native-assignment-journey",
          "performer": "claude:nro_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:nro_t1_r1 bench test --package ./cmd/bench --run ^TestNativeRecordCLIJourney$ at e2068642",
            "digest": "sha256:940103462384c4cb643692618ed71408d44dd74e0f90bd0c0fff5c0418bf5068",
            "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,102,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "native-assignment-journey",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'",
          "exit_code": 0
        },
        {
          "id": "nr-c1-t1-guidance",
          "performer": "claude:nro_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:nro_t1_r1 bench test --check docs-currency-workflow at e2068642",
            "digest": "sha256:b1ebf2ac23789885598f2d39bf218e3f0b680e82c9526a429165a3ad168cb316",
            "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5409,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "guidance",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "nr-c1-t1-root-conformance",
          "performer": "claude:nro_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:nro_t1_r1 bench test --check root-conformance at e2068642",
            "digest": "sha256:e189f99b61a482d519541fe00f91c1d7c65abedb332665c1cd73e0500793ba90",
            "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16434,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "nr-c1-t1-anchors",
          "performer": "claude:nro_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:nro_t1_r1 bench test --package ./internal/anchors at e2068642",
            "digest": "sha256:40617b6a7b5b5bdcad40e8a1649c0de84d6875e2b44bc8fc6d0b924bcd4a8012",
            "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,1167,178\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "nr-c1-t1-conformance-package",
          "performer": "claude:nro_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:nro_t1_r1 bench test --package ./internal/conformance at e2068642",
            "digest": "sha256:d1b5d74f8ae14f7f15a2c7974ed5442e371fa873acb2748a94526891fbfbbbfe",
            "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,54757,1677\nfailures[0]{package,test,line,lines}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceProseBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem: listen unix /tmp/X6N65R/t/TestGuidanceProseBudgetRefusesNonRegularSubjectssocket175670898/001/.agents/skills/bench-craft-linked/SKILL.md: bind: invalid argument\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestSkillDescriptionBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem: listen unix /tmp/X6N65R/t/TestSkillDescriptionBudgetRefusesNonRegularSubjectssocket3610578920/001/.agents/skills/bench-craft-planted/SKILL.md: bind: invalid argument\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "nr-c1-t1-record-fixtures",
          "performer": "claude:nro_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:nro_t1_r1 bench test --package ./internal/reviewrecord/recordtest at e2068642",
            "digest": "sha256:b4b87f026cdf6a76b409392e7edca3a7d2c1a98af24f3c56b80dfc3eb4871a28",
            "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "record-fixtures",
          "command": "bench test --package ./internal/reviewrecord/recordtest",
          "exit_code": 0
        },
        {
          "id": "nr-c1-t1-whole-plan-omission",
          "performer": "claude:nro_t1_r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f8a65e2cc006d71c605967d925149a2cf2af06d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:nro_t1_r1 focused green and bench probe omit of the checkPlan call in internal/reviewrecord/assignment.go at e2068642",
            "digest": "sha256:3ebb042d1fa4e81f5cdac13ab554c28dbed1212145ded28294073e26634901ce",
            "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,197,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\ntree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/assignment.go,omit,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord/recordcmd,^TestAssignmentWholePlanValidation$,passed,6\nspilled{lines=12,bytes=1960,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/1494f9bd9d4036914f3897f88c9afdf6/1791536500750183733-abc69c5ebad2bc45.out}\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentWholePlanValidation/cyclic_dependency,\"assignment_test.go:152: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",9b7fe78f3825c7b491a43ebc7e8a7fe73c634c2f\\\\\\\\n\\\"; want exit 1 naming \\\"cycle\\\"\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentWholePlanValidation/invalid_non-target_ticket,\"assignment_test.go:142: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",9c8d026b8030b98cbef1c0cb38dd7c85577f0763\\\\\\\\n\\\"; want exit 1 naming \\\"2.md\\\"\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentWholePlanValidation/missing_dependency,\"assignment_test.go:147: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",71f10da2b26a75a85efa0f6ead7b2729617efe25\\\\\\\\n\\\"; want exit 1 naming \\\"9.md\\\"\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentWholePlanValidation/session_of_another_ticket,\"assignment_test.go:158: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  2.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",ff0b343216296251da9c2de7fc49d44b8251ccbd\\\\\\\\n\\\"; want exit 1 naming \\\"already authors ticket… (246 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "whole-plan-omission",
          "command": "bench test --package ./internal/reviewrecord/recordcmd --run '^TestAssignmentWholePlanValidation$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit only the complete candidate-validation call in RecordAssignment, retaining the validator used by ReadPlan; an unrelated invalid ticket must red TestAssignmentWholePlanValidation.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:nro_t1_r1 focused green and bench probe omit of the checkPlan call in internal/reviewrecord/assignment.go at e2068642",
              "digest": "sha256:3ebb042d1fa4e81f5cdac13ab554c28dbed1212145ded28294073e26634901ce",
              "excerpt": "tree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,197,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\ntree[1]{target,head,dirty}:\n  nro-build,e2068642f17fb744698626a66a26c6e67a7ce157,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/assignment.go,omit,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord/recordcmd,^TestAssignmentWholePlanValidation$,passed,6\nspilled{lines=12,bytes=1960,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/1494f9bd9d4036914f3897f88c9afdf6/1791536500750183733-abc69c5ebad2bc45.out}\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentWholePlanValidation/cyclic_dependency,\"assignment_test.go:152: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",9b7fe78f3825c7b491a43ebc7e8a7fe73c634c2f\\\\\\\\n\\\"; want exit 1 naming \\\"cycle\\\"\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentWholePlanValidation/invalid_non-target_ticket,\"assignment_test.go:142: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",9c8d026b8030b98cbef1c0cb38dd7c85577f0763\\\\\\\\n\\\"; want exit 1 naming \\\"2.md\\\"\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentWholePlanValidation/missing_dependency,\"assignment_test.go:147: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  1.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",71f10da2b26a75a85efa0f6ead7b2729617efe25\\\\\\\\n\\\"; want exit 1 naming \\\"9.md\\\"\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestAssignmentWholePlanValidation/session_of_another_ticket,\"assignment_test.go:158: assignment = exit 0, output \\\"assignment[1]{ticket,session,occurrence,predecessor,source}:\\\\\\\\n  2.md,fixture-author-1.md,\\\\\\\\\\\"1\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\",ff0b343216296251da9c2de7fc49d44b8251ccbd\\\\\\\\n\\\"; want exit 1 naming \\\"already authors ticket… (246 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
