# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/cleanup-member-traversal/spec.md",
  "plan_digest": "sha256:6326ba026c86a6e50db4f8cfb3e74eec9c5d576e5d39b782ab0072a0fbb3fd77",
  "implementation_session": "",
  "chunks": [
    {
      "id": "C1",
      "base": "76c8d256cd35d5c74fba8d9273d2bfcb51f89a71",
      "tip": "74e7512e3d459d25f184f1bf2cef1a6d547748a6",
      "plan_digest": "sha256:589aa35576d3e85b161c9894002d688f6dfec0850944362a1a8aaf21c17a0e1d",
      "source_digest": "bcf5a3429fd19b30b7c757a360d84adb7d409f61",
      "acceptance_rows": [
        "R01",
        "R02",
        "R03",
        "R04",
        "R05",
        "R06",
        "R08",
        "R09",
        "R10",
        "R11",
        "R12",
        "R13",
        "R14",
        "R15",
        "R19",
        "R24",
        "R27",
        "R28"
      ],
      "verification": [
        {
          "id": "C1-T00-baseline-capture",
          "performer": "ft363-author-00",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bcf5a3429fd19b30b7c757a360d84adb7d409f61",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-author-00",
            "digest": "sha256:b028b4908e6a66b54905ff12c08946e28cd530db54da1839bb84eeca90acace2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/clean_set.go,swap,failed,2,TestCleanupTraversalBaseline/explicit/first-failure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,CleanupTraversalBaseline,passed,32\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,788\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,2118,32\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/first-failure,\"cleanup_members_baseline_test.go:119: explicit/first-failure differs from its frozen record at line 20: got \\\"  rows 1 nil=false order=prefix 1 of 3\\\", frozen \\\"  rows 3 nil=false order=prefix 3 of 3\\\"\",63\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/later-failure,\"cleanup_members_baseline_test.go:119: explicit/later-failure differs from its frozen record at line 20: got \\\"  rows 2 nil=false order=prefix 2 of 3\\\", frozen \\\"  rows 3 nil=false order=prefix 3 of 3\\\"\",64\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "74e7512e3d459d25f184f1bf2cef1a6d547748a6",
          "requirement": "baseline-capture",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversalBaseline'",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the unreached recorded suffix rows from the explicit stopped result and record the R26 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-author-00",
              "digest": "sha256:b028b4908e6a66b54905ff12c08946e28cd530db54da1839bb84eeca90acace2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/clean_set.go,swap,failed,2,TestCleanupTraversalBaseline/explicit/first-failure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,CleanupTraversalBaseline,passed,32\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,788\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,2118,32\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/first-failure,\"cleanup_members_baseline_test.go:119: explicit/first-failure differs from its frozen record at line 20: got \\\"  rows 1 nil=false order=prefix 1 of 3\\\", frozen \\\"  rows 3 nil=false order=prefix 3 of 3\\\"\",63\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/later-failure,\"cleanup_members_baseline_test.go:119: explicit/later-failure differs from its frozen record at line 20: got \\\"  rows 2 nil=false order=prefix 2 of 3\\\", frozen \\\"  rows 3 nil=false order=prefix 3 of 3\\\"\",64\nskips[0]{package,test,reason}:\n"
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
  },
  "amendments": [
    {
      "from": "sha256:589aa35576d3e85b161c9894002d688f6dfec0850944362a1a8aaf21c17a0e1d",
      "to": "sha256:6326ba026c86a6e50db4f8cfb3e74eec9c5d576e5d39b782ab0072a0fbb3fd77",
      "chunk_ids": {
        "C1": [
          "C1"
        ]
      }
    }
  ]
}
```
