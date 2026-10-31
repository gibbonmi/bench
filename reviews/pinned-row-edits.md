# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/pinned-row-edits/spec.md",
  "plan_digest": "sha256:403f9717792add2a8526f4b0564e0f52260a4a6f2e6f07ca167773e76b05023e",
  "implementation_session": "",
  "chunks": [
    {
      "id": "PE-C1",
      "base": "fba398122a3165ab735ed509aa85b05136ea716b",
      "tip": "42537669f1fe01765198242087ecdf08fa4a9380",
      "plan_digest": "sha256:403f9717792add2a8526f4b0564e0f52260a4a6f2e6f07ca167773e76b05023e",
      "source_digest": "3e8a8c55c547fac104283d73a89d2d25fed9a883",
      "acceptance_rows": [
        "PE01",
        "PE02",
        "PE03",
        "PE04",
        "PE05",
        "PE06",
        "PE07",
        "PE08",
        "PE09",
        "PE10",
        "PE11",
        "PE12",
        "PE13",
        "PE14",
        "PE15",
        "PE16",
        "PE17",
        "PE18",
        "PE19",
        "PE20",
        "PE21",
        "PE22",
        "PE23",
        "PE24",
        "PE25",
        "PE26",
        "PE27",
        "PE28",
        "PE29"
      ],
      "verification": [
        {
          "id": "pe-c1-t1-plan-approve",
          "performer": "ft437-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e2273e3108a461ca260be080473ed19a4529566",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft437-author-01",
            "digest": "sha256:3afc32b64ae7555f3f32b0f2491631ead5a9ca2381ff3e36878a9cec434138e9",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,16,TestPinnedRowApproveRereadsTheWorkingFile/changed-after-plan,40,1,0,null,yes\n  2,bit,swap,failed,4,TestPinnedRowPlanRefusesAHostileWorkingFile/empty,26,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment,PinnedRow|HeadBinding,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "5fb78f311a1458299ddb06dd07105a52e28e026e",
          "requirement": "plan-approve",
          "command": "bench test --package ./internal/commitment --run 'PinnedRow|HeadBinding'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make readTree return the default-branch tip for a changed outcome source, and record that TestPinnedRowPlanReadsTheWorkingFile turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft437-author-01",
              "digest": "sha256:3afc32b64ae7555f3f32b0f2491631ead5a9ca2381ff3e36878a9cec434138e9",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,16,TestPinnedRowApproveRereadsTheWorkingFile/changed-after-plan,40,1,0,null,yes\n  2,bit,swap,failed,4,TestPinnedRowPlanRefusesAHostileWorkingFile/empty,26,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment,PinnedRow|HeadBinding,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "pe-c1-t1-hostile-read",
          "performer": "ft437-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e2273e3108a461ca260be080473ed19a4529566",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft437-author-01",
            "digest": "sha256:3afc32b64ae7555f3f32b0f2491631ead5a9ca2381ff3e36878a9cec434138e9",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,16,TestPinnedRowApproveRereadsTheWorkingFile/changed-after-plan,40,1,0,null,yes\n  2,bit,swap,failed,4,TestPinnedRowPlanRefusesAHostileWorkingFile/empty,26,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment,PinnedRow|HeadBinding,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "5fb78f311a1458299ddb06dd07105a52e28e026e",
          "requirement": "hostile-read",
          "command": "bench test --package ./internal/commitment --run 'PinnedRowPlanRefusesAHostileWorkingFile'",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace the guarded working read with a plain os.ReadFile, and record that TestPinnedRowPlanRefusesAHostileWorkingFile turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft437-author-01",
              "digest": "sha256:3afc32b64ae7555f3f32b0f2491631ead5a9ca2381ff3e36878a9cec434138e9",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,16,TestPinnedRowApproveRereadsTheWorkingFile/changed-after-plan,40,1,0,null,yes\n  2,bit,swap,failed,4,TestPinnedRowPlanRefusesAHostileWorkingFile/empty,26,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment,PinnedRow|HeadBinding,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "pe-c1-t1-commit",
          "performer": "ft437-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e2273e3108a461ca260be080473ed19a4529566",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft437-author-01",
            "digest": "sha256:069d1389b99bf8ac3d83962085c2a114dc057287ff95a274ccb620e885d085cf",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCommitmentApprovedRowEdit/policy-only,31,1,0,null,yes\n  2,bit,swap,failed,4,TestCommitmentApprovedRowEdit/approved,31,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commit ./internal/commitment/repository,CommitmentApprovedRowEdit|CommitmentProtectedScope|AdmitPublication|AuthorizeCandidate,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "5fb78f311a1458299ddb06dd07105a52e28e026e",
          "requirement": "commit",
          "command": "bench test --package ./internal/commit --run 'CommitmentApprovedRowEdit|CommitmentProtectedScope'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make approvedTransition skip the source read of a changed outcome source, and record that TestCommitmentApprovedRowEdit turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft437-author-01",
              "digest": "sha256:069d1389b99bf8ac3d83962085c2a114dc057287ff95a274ccb620e885d085cf",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCommitmentApprovedRowEdit/policy-only,31,1,0,null,yes\n  2,bit,swap,failed,4,TestCommitmentApprovedRowEdit/approved,31,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commit ./internal/commitment/repository,CommitmentApprovedRowEdit|CommitmentProtectedScope|AdmitPublication|AuthorizeCandidate,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "pe-c1-t1-landing",
          "performer": "ft437-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e2273e3108a461ca260be080473ed19a4529566",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft437-author-01",
            "digest": "sha256:069d1389b99bf8ac3d83962085c2a114dc057287ff95a274ccb620e885d085cf",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCommitmentApprovedRowEdit/policy-only,31,1,0,null,yes\n  2,bit,swap,failed,4,TestCommitmentApprovedRowEdit/approved,31,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commit ./internal/commitment/repository,CommitmentApprovedRowEdit|CommitmentProtectedScope|AdmitPublication|AuthorizeCandidate,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "5fb78f311a1458299ddb06dd07105a52e28e026e",
          "requirement": "landing",
          "command": "bench test --package ./internal/commitment/repository --run 'AdmitPublication|AuthorizeCandidate'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make approvedTransition read the default-branch tip for a changed outcome source, and record that TestAdmitPublicationGradesAChangedRowAtTheComposedTree turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft437-author-01",
              "digest": "sha256:069d1389b99bf8ac3d83962085c2a114dc057287ff95a274ccb620e885d085cf",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCommitmentApprovedRowEdit/policy-only,31,1,0,null,yes\n  2,bit,swap,failed,4,TestCommitmentApprovedRowEdit/approved,31,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commit ./internal/commitment/repository,CommitmentApprovedRowEdit|CommitmentProtectedScope|AdmitPublication|AuthorizeCandidate,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "pe-c1-t1-retro-walk",
          "performer": "ft437-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e2273e3108a461ca260be080473ed19a4529566",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft437-author-01",
            "digest": "sha256:aab373af132de802efcb7126ceb762eefc157dd743de9af7b4d8b2336ac423f8",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/retro.go,swap,failed,2,TestRetroRefusesSymlinkedDestinationComponents/missing,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,Retro,passed,43\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,810\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,559,43\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestRetroRefusesSymlinkedDestinationComponents/missing,\"retro_test.go:139: retro through missing symlink = \\\"error: cannot write capture/retros/linked.md — mkdirat capture/retros: statat capture/retros: path escapes from parent\\\\\\\\n\\\"/1, want a symbolic-link refusal\",1\n  github.com/gibbonmi/bench/internal/roadmap,TestRetroRefusesSymlinkedDestinationComponents/outside,\"retro_test.go:139: retro through outside symlink = \\\"error: cannot write capture/retros/linked.md — mkdirat capture/retros: statat capture/retros: path escapes from parent\\\\\\\\n\\\"/1, want a symbolic-link refusal\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5fb78f311a1458299ddb06dd07105a52e28e026e",
          "requirement": "retro-walk",
          "command": "bench test --package ./internal/roadmap --run 'Retro'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the call of the moved component walk from the retro write, and record that TestRetroRefusesSymlinkedDestinationComponents turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft437-author-01",
              "digest": "sha256:aab373af132de802efcb7126ceb762eefc157dd743de9af7b4d8b2336ac423f8",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/retro.go,swap,failed,2,TestRetroRefusesSymlinkedDestinationComponents/missing,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,Retro,passed,43\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,810\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,559,43\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestRetroRefusesSymlinkedDestinationComponents/missing,\"retro_test.go:139: retro through missing symlink = \\\"error: cannot write capture/retros/linked.md — mkdirat capture/retros: statat capture/retros: path escapes from parent\\\\\\\\n\\\"/1, want a symbolic-link refusal\",1\n  github.com/gibbonmi/bench/internal/roadmap,TestRetroRefusesSymlinkedDestinationComponents/outside,\"retro_test.go:139: retro through outside symlink = \\\"error: cannot write capture/retros/linked.md — mkdirat capture/retros: statat capture/retros: path escapes from parent\\\\\\\\n\\\"/1, want a symbolic-link refusal\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "pe-c1-t1-rep1-plan-approve",
          "performer": "ft437-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e8a8c55c547fac104283d73a89d2d25fed9a883",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft437-author-01-rep1",
            "digest": "sha256:df8cd178e77402ba40265a6232541b9c7850a542c2ca4a8c27e3902b51f63e78",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestPinnedRowPlanReadsTheWorkingFile/linked-parent,4,1,0,null,yes\n  2,bit,swap,failed,4,TestPinnedRowPlanRefusesAHostileWorkingFile/empty,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment,^TestPinnedRowPlanReadsTheWorkingFile$|^TestPinnedRowPlanRefusesAHostileWorkingFile$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "42537669f1fe01765198242087ecdf08fa4a9380",
          "requirement": "plan-approve",
          "command": "bench test --package ./internal/commitment --run 'PinnedRow|HeadBinding'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make readTree return the default-branch tip for a changed outcome source, and record that TestPinnedRowPlanReadsTheWorkingFile turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft437-author-01-rep1",
              "digest": "sha256:df8cd178e77402ba40265a6232541b9c7850a542c2ca4a8c27e3902b51f63e78",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestPinnedRowPlanReadsTheWorkingFile/linked-parent,4,1,0,null,yes\n  2,bit,swap,failed,4,TestPinnedRowPlanRefusesAHostileWorkingFile/empty,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment,^TestPinnedRowPlanReadsTheWorkingFile$|^TestPinnedRowPlanRefusesAHostileWorkingFile$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "pe-c1-t1-rep1-hostile-read",
          "performer": "ft437-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e8a8c55c547fac104283d73a89d2d25fed9a883",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft437-author-01-rep1",
            "digest": "sha256:df8cd178e77402ba40265a6232541b9c7850a542c2ca4a8c27e3902b51f63e78",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestPinnedRowPlanReadsTheWorkingFile/linked-parent,4,1,0,null,yes\n  2,bit,swap,failed,4,TestPinnedRowPlanRefusesAHostileWorkingFile/empty,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment,^TestPinnedRowPlanReadsTheWorkingFile$|^TestPinnedRowPlanRefusesAHostileWorkingFile$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "42537669f1fe01765198242087ecdf08fa4a9380",
          "requirement": "hostile-read",
          "command": "bench test --package ./internal/commitment --run 'PinnedRowPlanRefusesAHostileWorkingFile'",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace the guarded working read with a plain os.ReadFile, and record that TestPinnedRowPlanRefusesAHostileWorkingFile turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft437-author-01-rep1",
              "digest": "sha256:df8cd178e77402ba40265a6232541b9c7850a542c2ca4a8c27e3902b51f63e78",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestPinnedRowPlanReadsTheWorkingFile/linked-parent,4,1,0,null,yes\n  2,bit,swap,failed,4,TestPinnedRowPlanRefusesAHostileWorkingFile/empty,8,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commitment,^TestPinnedRowPlanReadsTheWorkingFile$|^TestPinnedRowPlanRefusesAHostileWorkingFile$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "pe-c1-t1-rep1-commit",
          "performer": "ft437-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e8a8c55c547fac104283d73a89d2d25fed9a883",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft437-author-01-rep1",
            "digest": "sha256:f5ec3ee91818a4217110703dc320494cc101e90d04974df2f48f53745b41be96",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCommitmentApprovedRowEdit/policy-only,4,1,0,null,yes\n  2,bit,swap,failed,2,TestAdmitPublicationGradesAChangedRowAtTheComposedTree/approved,3,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commit ./internal/commitment/repository,^TestCommitmentApprovedRowEdit$|^TestAdmitPublicationGradesAChangedRowAtTheComposedTree$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "42537669f1fe01765198242087ecdf08fa4a9380",
          "requirement": "commit",
          "command": "bench test --package ./internal/commit --run 'CommitmentApprovedRowEdit|CommitmentProtectedScope'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make approvedTransition skip the source read of a changed outcome source, and record that TestCommitmentApprovedRowEdit turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft437-author-01-rep1",
              "digest": "sha256:f5ec3ee91818a4217110703dc320494cc101e90d04974df2f48f53745b41be96",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCommitmentApprovedRowEdit/policy-only,4,1,0,null,yes\n  2,bit,swap,failed,2,TestAdmitPublicationGradesAChangedRowAtTheComposedTree/approved,3,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commit ./internal/commitment/repository,^TestCommitmentApprovedRowEdit$|^TestAdmitPublicationGradesAChangedRowAtTheComposedTree$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "pe-c1-t1-rep1-landing",
          "performer": "ft437-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e8a8c55c547fac104283d73a89d2d25fed9a883",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft437-author-01-rep1",
            "digest": "sha256:f5ec3ee91818a4217110703dc320494cc101e90d04974df2f48f53745b41be96",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCommitmentApprovedRowEdit/policy-only,4,1,0,null,yes\n  2,bit,swap,failed,2,TestAdmitPublicationGradesAChangedRowAtTheComposedTree/approved,3,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commit ./internal/commitment/repository,^TestCommitmentApprovedRowEdit$|^TestAdmitPublicationGradesAChangedRowAtTheComposedTree$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "42537669f1fe01765198242087ecdf08fa4a9380",
          "requirement": "landing",
          "command": "bench test --package ./internal/commitment/repository --run 'AdmitPublication|AuthorizeCandidate'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make approvedTransition read the default-branch tip for a changed outcome source, and record that TestAdmitPublicationGradesAChangedRowAtTheComposedTree turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft437-author-01-rep1",
              "digest": "sha256:f5ec3ee91818a4217110703dc320494cc101e90d04974df2f48f53745b41be96",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCommitmentApprovedRowEdit/policy-only,4,1,0,null,yes\n  2,bit,swap,failed,2,TestAdmitPublicationGradesAChangedRowAtTheComposedTree/approved,3,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/commit ./internal/commitment/repository,^TestCommitmentApprovedRowEdit$|^TestAdmitPublicationGradesAChangedRowAtTheComposedTree$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "pe-c1-t1-rep1-retro-walk",
          "performer": "ft437-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e8a8c55c547fac104283d73a89d2d25fed9a883",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft437-author-01-rep1",
            "digest": "sha256:063b9909fb29e762a87e098a51c0682ac4e937847af1bd2ee4a557059d50fd98",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/retro.go,swap,failed,2,TestRetroRefusesSymlinkedDestinationComponents/missing,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestRetroRefusesSymlinkedDestinationComponents$,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,808\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,27,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestRetroRefusesSymlinkedDestinationComponents/missing,\"retro_test.go:139: retro through missing symlink = \\\"error: cannot write capture/retros/linked.md — mkdirat capture/retros: statat capture/retros: path escapes from parent\\\\\\\\n\\\"/1, want a symbolic-link refusal\",1\n  github.com/gibbonmi/bench/internal/roadmap,TestRetroRefusesSymlinkedDestinationComponents/outside,\"retro_test.go:139: retro through outside symlink = \\\"error: cannot write capture/retros/linked.md — mkdirat capture/retros: statat capture/retros: path escapes from parent\\\\\\\\n\\\"/1, want a symbolic-link refusal\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "42537669f1fe01765198242087ecdf08fa4a9380",
          "requirement": "retro-walk",
          "command": "bench test --package ./internal/roadmap --run 'Retro'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the call of the moved component walk from the retro write, and record that TestRetroRefusesSymlinkedDestinationComponents turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft437-author-01-rep1",
              "digest": "sha256:063b9909fb29e762a87e098a51c0682ac4e937847af1bd2ee4a557059d50fd98",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/retro.go,swap,failed,2,TestRetroRefusesSymlinkedDestinationComponents/missing,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestRetroRefusesSymlinkedDestinationComponents$,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,808\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,27,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestRetroRefusesSymlinkedDestinationComponents/missing,\"retro_test.go:139: retro through missing symlink = \\\"error: cannot write capture/retros/linked.md — mkdirat capture/retros: statat capture/retros: path escapes from parent\\\\\\\\n\\\"/1, want a symbolic-link refusal\",1\n  github.com/gibbonmi/bench/internal/roadmap,TestRetroRefusesSymlinkedDestinationComponents/outside,\"retro_test.go:139: retro through outside symlink = \\\"error: cannot write capture/retros/linked.md — mkdirat capture/retros: statat capture/retros: path escapes from parent\\\\\\\\n\\\"/1, want a symbolic-link refusal\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "pe-c1-standards-1",
          "performer": "ft437-review-standards-01",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e2273e3108a461ca260be080473ed19a4529566",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft437-review-standards-01",
            "digest": "sha256:67208fd953f14370a81a4fc92eef22182228c4c45164641de924cc056af6f566",
            "excerpt": "Standards axis, PE-C1, opus/medium. 4 findings.\nS1 (7, auto-fix): internal/commitment/repository/head_binding_test.go:168-177 pastes the pinnedRows fixture of internal/commitment/pinned_row_test.go:22-31; AGENTS.md one source per fact (fixture harness pasted N times). Move it into commitmenttest.\nS2 (6, auto-fix): the \"commitment source %q changed from X to Y\" expectation is built three times (pinned_row_test.go:45-48, commit/commitment_test.go:239, repository/head_binding_test.go:165); head_binding_test.go already holds changedStage at :59-62.\nS3 (5, advice): readTree(..., checkout bool) flag argument with bare literals (sources.go:148, candidate.go:211).\nS4 (3, advice): two link walks in internal/bounds/classify.go:151-167 and :170-185.\n"
          },
          "axis": "Standards",
          "base": "fba398122a3165ab735ed509aa85b05136ea716b",
          "tip": "5fb78f311a1458299ddb06dd07105a52e28e026e",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "pe-c1-spec-1",
          "performer": "ft437-review-spec-01",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e2273e3108a461ca260be080473ed19a4529566",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft437-review-spec-01",
            "digest": "sha256:fa17c2c48d96cc4f01ea2cbb56361b243ba85002e2a50a253e099fccc898314b",
            "excerpt": "Spec axis, PE-C1, opus/medium. 3 findings. All 29 rows map to code and tests; read-location table matches readTree; PE26 and PE29 hold.\nP1 (9, auto-fix): PE21, PE22, PE25 why-clauses are false: under a transition protectedCandidate re-reads the candidate sources and refuses with the same changed-from text, wrapped as \"candidate changes protected commitment\". The author's negative assertion (commit/commitment_test.go:120-122) is within the acceptance.\nP2 (7, auto-fix): the PE25 landing subtest checks only strings.Contains; a skipped landing read stays green. Mirror the negative assertion.\nP3 (6, auto-fix): spec \"No change reaches ... protectedCandidate\" is false: candidate.go edits its signature for sourceRead.\n"
          },
          "axis": "Spec",
          "base": "fba398122a3165ab735ed509aa85b05136ea716b",
          "tip": "5fb78f311a1458299ddb06dd07105a52e28e026e",
          "finding_ids": [
            "P1",
            "P2",
            "P3"
          ],
          "supersedes": []
        },
        {
          "id": "pe-c1-coverage-1",
          "performer": "ft437-review-coverage-01",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3e2273e3108a461ca260be080473ed19a4529566",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft437-review-coverage-01",
            "digest": "sha256:98273e2a2c0000a3cdc13d4a5ce2e23e451f200678fc33f6c87983ba27bcfdd5",
            "excerpt": "Coverage axis, PE-C1, opus/medium. 3 findings, none blocking. Hostile checklist fully decided; identity parity pinned by PE15/PE16; refusals write nothing. No test or probe run.\nC1 (4, auto-fix): internal/commitment/delivery.go:207 ChangedSource: removing the ID half or the Path half survives; no test moves an id to a new path or binds a new id at a current path.\nC2 (7 survive, low impact, advice): sources.go:66-72 refusal detail checked only by prefix.\nC3 (6, advice): classify.go:159 non-NotExist Lstat error branch untested; plan still refuses.\n"
          },
          "axis": "Coverage",
          "base": "fba398122a3165ab735ed509aa85b05136ea716b",
          "tip": "5fb78f311a1458299ddb06dd07105a52e28e026e",
          "finding_ids": [
            "C1"
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
      "from": "sha256:7f79cd015f4f0d8108f19c6a266efb37c99dc32558e1f1744c4e4d703f69336e",
      "to": "sha256:403f9717792add2a8526f4b0564e0f52260a4a6f2e6f07ca167773e76b05023e",
      "chunk_ids": {
        "PE-C1": [
          "PE-C1"
        ]
      }
    }
  ]
}
```

## PE-C1 review round 1

The three axes ran on opus at medium effort over `fba398122..5fb78f311`. The repair allowance of PE-C1 is two cycles, and this round consumes none. Cycle 1 repairs the accepted findings below.

## Standards

Findings: 2 accepted. Worst issue: S1.

- S1, auto-fix, confidence 7: `internal/commitment/repository/head_binding_test.go:168` pastes the `pinnedRows` fixture of `internal/commitment/pinned_row_test.go:22`. The AGENTS.md one-source rule names a fixture harness pasted N times. Move the fixture into `commitmenttest` and call it from both packages.
- S2, auto-fix, confidence 6: three tests build the `changed from` expectation each in their own way. The repository test file already holds the `changedStage` builder. Build the expectation in one place.

## Spec

Findings: 3 accepted. Worst issue: P1.

- P1, auto-fix, confidence 9: the catch clauses of PE21, PE22, and PE25 are false. Under a transition, `protectedCandidate` refuses the same input with the same `changed from` text and wraps it as `candidate changes protected commitment`. The clauses must name the face check as the catch.
- P2, auto-fix, confidence 7: the PE25 landing subtest checks only the substring. Add the negative assertion of the face, as the commit rows do.
- P3, auto-fix, confidence 6: the sentence "No change reaches ... `protectedCandidate`" is false, because the change edits its signature. State that only the signature changes.

## Coverage

Findings: 1 accepted. Worst issue: C1.

- C1, auto-fix, confidence 4: no test makes a mutation of either half of the id-and-path match of `ChangedSource` fail (`internal/commitment/delivery.go:207`). Add one subtest that moves an id to a new path and one that binds a new id at a current path. Each subtest must make plan keep its old read location. This is the one hardening cycle of PE-C1.

## Advice

- Standards S3: `readTree` takes a bare bool flag, and `false` needs a comment at one caller. A named mode would read better.
- Standards S4: `internal/bounds` holds two link walks with the same body. The spec decides to keep the absolute walk and the root-relative walk separate.
- Coverage C2: the hostile-file refusals match only the `refused:` prefix.
- Coverage C3: no test reaches the `Lstat` error branch that is not a missing file, and the plan still refuses.
