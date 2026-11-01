# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/rebuild-remedy-owner/spec.md",
  "plan_digest": "sha256:5af4ab8f4f17fc2df240a7bc3eeff75e28ca925d9eba968241d049adfe46a288",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RB-C1",
      "base": "7d02ba15a2dc0a589a3352eefb3e11f3948ac89f",
      "tip": "40fbf22e237dd5c6604e7aa45f3468eaa09005b6",
      "plan_digest": "sha256:5af4ab8f4f17fc2df240a7bc3eeff75e28ca925d9eba968241d049adfe46a288",
      "source_digest": "fb5825339f0f0b051f977f6fd92a3b4d29daf22a",
      "acceptance_rows": [
        "RB01",
        "RB02",
        "RB03",
        "RB04",
        "RB05",
        "RB06",
        "RB07",
        "RB08",
        "RB09",
        "RB10",
        "RB11",
        "RB12",
        "RB17",
        "RB18",
        "RB19",
        "RB20",
        "RB21",
        "RB22",
        "RB26"
      ],
      "verification": [
        {
          "id": "rb-c1-t01-owner-kind",
          "performer": "ft344-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fb5825339f0f0b051f977f6fd92a3b4d29daf22a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft344-author-01",
            "digest": "sha256:758a6a08763a1fb4b8b7728f076832dfb584190d7e272b221c83e19848cf6aca",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,5,TestRebuildActionReadsTheProcessHome,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/benchguard ./internal/freshness,RebuildRemedy|RebuildAction|RefusalCarries|RefusalRebuildAction,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1956\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,fail,109,5\n  github.com/gibbonmi/bench/internal/freshness,pass,5,1\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/benchguard,TestRebuildActionReadsTheProcessHome,\"rebuild_remedy_test.go:93: RB08: RebuildAction(\\\"/dev/shm/AB6P3F/t/TestRebuildActionReadsTheProcessHome2811353708/001/worktrees/k/x\\\") = \\\"cd '/dev/shm/AB6P3F/t/TestRebuildActionReadsTheProcessHome2811353708/001/worktrees/k/x' && bash scripts/… (547 bytes)\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestRebuildRemedyHostileRoot,\"rebuild_remedy_test.go:77: RB07: \\\"cd '/dev/shm/AB6P3F/t/TestRebuildRemedyHostileRoot3429196343/001/worktrees/k/it'\\\\\\\\\\\\\\\\''s a *' && bash scripts/go-build.sh '/dev/shm/AB6P3F/t/TestRebuildRemedyHostileRoot3429196343/001/worktrees/k/it'\\\\\\\\\\\\\\\\''s a *'… (369 bytes)\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestRebuildRemedyKind,\"rebuild_remedy_test.go:43: RB01 a pool root: RebuildIn(\\\"/dev/shm/AB6P3F/t/TestRebuildRemedyKind2747663244/001\\\", \\\"/dev/shm/AB6P3F/t/TestRebuildRemedyKind2747663244/001/worktrees/k/x\\\") = {Action:cd '/dev/shm/AB6P3F/t/TestRebuildRemedyKind2747… (584 bytes)\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestRebuildRemedyPassesTheHooks,\"rebuild_remedy_test.go:59: RB03: PoolReference(\\\"cd '/dev/shm/AB6P3F/t/TestRebuildRemedyPassesTheHooks4091076315/001/worktrees/k/x' && bash scripts/go-build.sh '/dev/shm/AB6P3F/t/TestRebuildRemedyPassesTheHooks4091076315/001/worktrees/k/x' '… (437 bytes)\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestRefusalCarriesThePoolRemedy,\"rebuild_remedy_test.go:106: RB09: Error() = \\\"bench binary \\\\\\\\\\\"/dev/shm/AB6P3F/t/TestRefusalCarriesThePoolRemedy1402083413/001/worktrees/k/x/dist/bench\\\\\\\\\\\" is untrusted: fixture; rebuild with cd '/dev/shm/AB6P3F/t/TestRefusalCarriesThePoolRemedy… (508 bytes)\",2\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "40fbf22e237dd5c6604e7aa45f3468eaa09005b6",
          "requirement": "owner-kind",
          "command": "bench test --package ./internal/benchguard --package ./internal/freshness --run 'RebuildRemedy|RebuildAction|RefusalCarries|RefusalRebuildAction'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make RebuildIn return CheckoutRebuildAction for every root, and record that TestRebuildRemedyPassesTheHooks turns red on the pool-root case.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft344-author-01",
              "digest": "sha256:758a6a08763a1fb4b8b7728f076832dfb584190d7e272b221c83e19848cf6aca",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,5,TestRebuildActionReadsTheProcessHome,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/benchguard ./internal/freshness,RebuildRemedy|RebuildAction|RefusalCarries|RefusalRebuildAction,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1956\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,fail,109,5\n  github.com/gibbonmi/bench/internal/freshness,pass,5,1\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/benchguard,TestRebuildActionReadsTheProcessHome,\"rebuild_remedy_test.go:93: RB08: RebuildAction(\\\"/dev/shm/AB6P3F/t/TestRebuildActionReadsTheProcessHome2811353708/001/worktrees/k/x\\\") = \\\"cd '/dev/shm/AB6P3F/t/TestRebuildActionReadsTheProcessHome2811353708/001/worktrees/k/x' && bash scripts/… (547 bytes)\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestRebuildRemedyHostileRoot,\"rebuild_remedy_test.go:77: RB07: \\\"cd '/dev/shm/AB6P3F/t/TestRebuildRemedyHostileRoot3429196343/001/worktrees/k/it'\\\\\\\\\\\\\\\\''s a *' && bash scripts/go-build.sh '/dev/shm/AB6P3F/t/TestRebuildRemedyHostileRoot3429196343/001/worktrees/k/it'\\\\\\\\\\\\\\\\''s a *'… (369 bytes)\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestRebuildRemedyKind,\"rebuild_remedy_test.go:43: RB01 a pool root: RebuildIn(\\\"/dev/shm/AB6P3F/t/TestRebuildRemedyKind2747663244/001\\\", \\\"/dev/shm/AB6P3F/t/TestRebuildRemedyKind2747663244/001/worktrees/k/x\\\") = {Action:cd '/dev/shm/AB6P3F/t/TestRebuildRemedyKind2747… (584 bytes)\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestRebuildRemedyPassesTheHooks,\"rebuild_remedy_test.go:59: RB03: PoolReference(\\\"cd '/dev/shm/AB6P3F/t/TestRebuildRemedyPassesTheHooks4091076315/001/worktrees/k/x' && bash scripts/go-build.sh '/dev/shm/AB6P3F/t/TestRebuildRemedyPassesTheHooks4091076315/001/worktrees/k/x' '… (437 bytes)\",1\n  github.com/gibbonmi/bench/internal/benchguard,TestRefusalCarriesThePoolRemedy,\"rebuild_remedy_test.go:106: RB09: Error() = \\\"bench binary \\\\\\\\\\\"/dev/shm/AB6P3F/t/TestRefusalCarriesThePoolRemedy1402083413/001/worktrees/k/x/dist/bench\\\\\\\\\\\" is untrusted: fixture; rebuild with cd '/dev/shm/AB6P3F/t/TestRefusalCarriesThePoolRemedy… (508 bytes)\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "rb-c1-t01-command-entry",
          "performer": "ft344-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fb5825339f0f0b051f977f6fd92a3b4d29daf22a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft344-author-01",
            "digest": "sha256:6ed9ef7858ee60d0833cba7c0d6ccb98ed6e5c8f8fd4de3befb208aa67198b7e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestCommandsBriefGradesTheRunningExecutable/a_mismatched_seal_in_a_pool_checkout_names_the_worktree_build,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,CommandsBriefGradesTheRunningExecutable,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,576\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,3178,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestCommandsBriefGradesTheRunningExecutable/a_mismatched_seal_in_a_pool_checkout_names_the_worktree_build,\"commands_brief_test.go:70: commands --brief refusal = \\\"error: bench binary does not match this checkout's sources — cd '/dev/shm/JZGRFX/t/TestCommandsBriefGradesTheRunningExecutablea_mismatched_seal_in183059740/001/home/worktrees/k/x' && ba… (676 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "40fbf22e237dd5c6604e7aa45f3468eaa09005b6",
          "requirement": "command-entry",
          "command": "bench test --package ./cmd/bench --run 'CommandsBriefGradesTheRunningExecutable'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make RebuildAction read an empty home in place of benchhome.Dir(), and record that the pool-root subtest of TestCommandsBriefGradesTheRunningExecutable turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft344-author-01",
              "digest": "sha256:6ed9ef7858ee60d0833cba7c0d6ccb98ed6e5c8f8fd4de3befb208aa67198b7e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestCommandsBriefGradesTheRunningExecutable/a_mismatched_seal_in_a_pool_checkout_names_the_worktree_build,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,CommandsBriefGradesTheRunningExecutable,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,576\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,3178,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestCommandsBriefGradesTheRunningExecutable/a_mismatched_seal_in_a_pool_checkout_names_the_worktree_build,\"commands_brief_test.go:70: commands --brief refusal = \\\"error: bench binary does not match this checkout's sources — cd '/dev/shm/JZGRFX/t/TestCommandsBriefGradesTheRunningExecutablea_mismatched_seal_in183059740/001/home/worktrees/k/x' && ba… (676 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "rb-c1-t01-shell-copy",
          "performer": "ft344-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fb5825339f0f0b051f977f6fd92a3b4d29daf22a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft344-author-01",
            "digest": "sha256:8a34579e3b7b5afddc4b80d5d24570124b39b3280ed6f8a50b4bc55641fd2008",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.bench/lib/resolve-bench.sh,insert,failed,1,TestBenchFollowOnHookDegradedRim,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,TestBenchFollowOnHookDegradedRim,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,582\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,770,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestBenchFollowOnHookDegradedRim,\"bench_follow_on_test.go:334: RB17 stale core with a pool kit = (2, \\\"\\\", \\\"BLOCKED: the Bench binary is stale, so no Bench verb can answer. Rebuild it: bench worktree build '/dev/shm/URYD6H/t/TestBenchFollowOnHookDegradedRim3548653876/007/work… (567 bytes)\",2\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "40fbf22e237dd5c6604e7aa45f3468eaa09005b6",
          "requirement": "shell-copy",
          "command": "bench test --check system",
          "exit_code": 0,
          "probe": {
            "mutation": "Add a pool-path case to bench_rebuild_action that prints bench worktree build, and record that TestBenchFollowOnHookDegradedRim turns red on the pool-kit case.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft344-author-01",
              "digest": "sha256:8a34579e3b7b5afddc4b80d5d24570124b39b3280ed6f8a50b4bc55641fd2008",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.bench/lib/resolve-bench.sh,insert,failed,1,TestBenchFollowOnHookDegradedRim,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,TestBenchFollowOnHookDegradedRim,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,582\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,770,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestBenchFollowOnHookDegradedRim,\"bench_follow_on_test.go:334: RB17 stale core with a pool kit = (2, \\\"\\\", \\\"BLOCKED: the Bench binary is stale, so no Bench verb can answer. Rebuild it: bench worktree build '/dev/shm/URYD6H/t/TestBenchFollowOnHookDegradedRim3548653876/007/work… (567 bytes)\",2\nskips[0]{package,test,reason}:\n"
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
