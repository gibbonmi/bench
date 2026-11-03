# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/prepublication-preflight/spec.md",
  "plan_digest": "sha256:351505691273d6505a9a5fcf36349186e657baf3e6769643711555ead6885e4d",
  "implementation_session": "",
  "chunks": [
    {
      "id": "C1",
      "base": "3d0598ba1b589d392eb8d11aa0c658e27e809e4c",
      "tip": "411b02763e57c9f901b47697ab5271e03b4d45fd",
      "plan_digest": "sha256:351505691273d6505a9a5fcf36349186e657baf3e6769643711555ead6885e4d",
      "source_digest": "52ea596345ba18cf4db25e522c2d2a8acf93aa5f",
      "acceptance_rows": [
        "PP01",
        "PP02",
        "PP03",
        "PP04",
        "PP05",
        "PP06",
        "PP07",
        "PP08",
        "PP09",
        "PP32",
        "PP33",
        "PP34"
      ],
      "verification": [
        {
          "id": "c1-01-candidate-tests",
          "performer": "ppf-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "52ea596345ba18cf4db25e522c2d2a8acf93aa5f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ppf-author-01",
            "digest": "sha256:3241a28cd297394534a83bcb28401c9755c441fc7989d1eae245f821bb03c8e3",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/plan.go,swap,failed,1,TestBuildVerdictAtGradesTheCandidate/PP02_completion_plan_at_the_candidate,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestBuildVerdictAtGradesTheCandidate$,passed,5\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,397\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,333,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestBuildVerdictAtGradesTheCandidate/PP02_completion_plan_at_the_candidate,\"candidate_test.go:36: plan candidate = (0), want 1 with completion-plan red at the candidate:\",4\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "411b02763e57c9f901b47697ab5271e03b4d45fd",
          "requirement": "candidate-tests",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0,
          "probe": {
            "mutation": "Read the completion plan at HEAD in the candidate grade and record the PP02 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ppf-author-01",
              "digest": "sha256:3241a28cd297394534a83bcb28401c9755c441fc7989d1eae245f821bb03c8e3",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/plan.go,swap,failed,1,TestBuildVerdictAtGradesTheCandidate/PP02_completion_plan_at_the_candidate,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestBuildVerdictAtGradesTheCandidate$,passed,5\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,397\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,333,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestBuildVerdictAtGradesTheCandidate/PP02_completion_plan_at_the_candidate,\"candidate_test.go:36: plan candidate = (0), want 1 with completion-plan red at the candidate:\",4\nskips[0]{package,test,reason}:\n"
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
