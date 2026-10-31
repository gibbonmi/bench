# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/pinned-row-edits/spec.md",
  "plan_digest": "sha256:7f79cd015f4f0d8108f19c6a266efb37c99dc32558e1f1744c4e4d703f69336e",
  "implementation_session": "",
  "chunks": [
    {
      "id": "PE-C1",
      "base": "fba398122a3165ab735ed509aa85b05136ea716b",
      "tip": "5fb78f311a1458299ddb06dd07105a52e28e026e",
      "plan_digest": "sha256:7f79cd015f4f0d8108f19c6a266efb37c99dc32558e1f1744c4e4d703f69336e",
      "source_digest": "3e2273e3108a461ca260be080473ed19a4529566",
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
