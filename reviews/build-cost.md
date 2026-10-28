# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/build-cost/spec.md",
  "plan_digest": "sha256:afed42110e5cf94f89620cf212044f2461429b735a54964f193ce5cadf70276d",
  "implementation_session": "",
  "chunks": [
    {
      "id": "BC-C1",
      "base": "77ec3655a4601a57165b33d58c4328bdfdf72d34",
      "tip": "f0d811335ac529969264109ac3439c54bc649f1c",
      "plan_digest": "sha256:afed42110e5cf94f89620cf212044f2461429b735a54964f193ce5cadf70276d",
      "source_digest": "09e8efc193f31a48ceba77ac89b1b2bf23665fa8",
      "acceptance_rows": [
        "BC24",
        "BC25",
        "BC26",
        "BC27",
        "BC28",
        "BC29",
        "BC30",
        "BC31",
        "BC32",
        "BC23",
        "BC33",
        "BC34",
        "BC35",
        "BC36",
        "BC37"
      ],
      "verification": [
        {
          "id": "t01-probe-01",
          "performer": "claude:bc-t01-author-20261028",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "09e8efc193f31a48ceba77ac89b1b2bf23665fa8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:bc-t01-author-20261028",
            "digest": "sha256:b872e89b309b3b13235935f0ef056788bfc8e3603c20742a479365a4ff4bd8ef",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/plan.go,swap,failed,3,TestPlanProbeOnlyChunks/absent_chunk_list,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestPlanProbeOnlyChunks,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,765\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,101,6\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/absent_chunk_list,\"plan_probe_test.go:64: plan grade refused the probe-only plan: chunk 2: missing required verification inventory\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/empty_chunk_list,\"plan_probe_test.go:64: plan grade refused the probe-only plan: chunk 2: missing required verification inventory\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/missing_probe_entry,\"plan_probe_test.go:74: chunk 2: missing required verification inventory\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f0d811335ac529969264109ac3439c54bc649f1c",
          "requirement": "probe-01",
          "command": "bench probe internal/reviewrecord/plan.go --swap 'if len(items) == 0 && owned == nil {' --with 'if len(items) == 0 {' --package ./internal/reviewrecord --run TestPlanProbeOnlyChunks",
          "exit_code": 0,
          "probe": {
            "mutation": "Refuse an empty chunk list again. The probe-only plan with an empty chunk list must then refuse, so the plan test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:bc-t01-author-20261028",
              "digest": "sha256:b872e89b309b3b13235935f0ef056788bfc8e3603c20742a479365a4ff4bd8ef",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/plan.go,swap,failed,3,TestPlanProbeOnlyChunks/absent_chunk_list,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestPlanProbeOnlyChunks,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,765\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,101,6\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/absent_chunk_list,\"plan_probe_test.go:64: plan grade refused the probe-only plan: chunk 2: missing required verification inventory\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/empty_chunk_list,\"plan_probe_test.go:64: plan grade refused the probe-only plan: chunk 2: missing required verification inventory\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestPlanProbeOnlyChunks/missing_probe_entry,\"plan_probe_test.go:74: chunk 2: missing required verification inventory\",1\nskips[0]{package,test,reason}:\n"
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
