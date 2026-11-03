# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/cleanup-member-traversal/spec.md",
  "plan_digest": "sha256:3f9255958885eae9a6109137977b57408b8d5287016e33fd22b9d4d306fa75a5",
  "implementation_session": "",
  "chunks": [
    {
      "id": "C1",
      "base": "76c8d256cd35d5c74fba8d9273d2bfcb51f89a71",
      "tip": "5e0110372b6adf94156ef2a64ac42935e19632a5",
      "plan_digest": "sha256:3f9255958885eae9a6109137977b57408b8d5287016e33fd22b9d4d306fa75a5",
      "source_digest": "f83eb0463ad351fd34d053f3eaade732acdb79b1",
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
        },
        {
          "id": "C1-T01-explicit-preflight",
          "performer": "ft363-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f83eb0463ad351fd34d053f3eaade732acdb79b1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-author-01",
            "digest": "sha256:7904ee497271b27c230dd9e6585279ab4442914ae23c43100c61a5598809342d",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,99,1,0,null,yes\n  2,bit,swap,failed,6,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,99,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,99,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "5e0110372b6adf94156ef2a64ac42935e19632a5",
          "requirement": "explicit-preflight",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the explicit all-member preflight before the recorded visits and record the R06 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-author-01",
              "digest": "sha256:7904ee497271b27c230dd9e6585279ab4442914ae23c43100c61a5598809342d",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,99,1,0,null,yes\n  2,bit,swap,failed,6,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,99,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,99,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "C1-T01-later-requalify",
          "performer": "ft363-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f83eb0463ad351fd34d053f3eaade732acdb79b1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-author-01",
            "digest": "sha256:7904ee497271b27c230dd9e6585279ab4442914ae23c43100c61a5598809342d",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,99,1,0,null,yes\n  2,bit,swap,failed,6,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,99,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,99,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "5e0110372b6adf94156ef2a64ac42935e19632a5",
          "requirement": "later-requalify",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the later requalification before a reached recorded visit and record the R09 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-author-01",
              "digest": "sha256:7904ee497271b27c230dd9e6585279ab4442914ae23c43100c61a5598809342d",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,99,1,0,null,yes\n  2,bit,swap,failed,6,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,99,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,99,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "C1-T01-locked-requalify",
          "performer": "ft363-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f83eb0463ad351fd34d053f3eaade732acdb79b1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-author-01",
            "digest": "sha256:7c192a48811fefd3562e8dca807f3c4a6a20a986f960d41dca7ea616e249a1c2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/resume.go,swap,failed,5,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,99\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1509\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,5577,99\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,\"clean_set_outcomes_test.go:322: late-drift apply = (0, \\\"worktree_cleanup[2]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/KVVQSZ/t/TestCleanSetApplyTimeStaleRefusalexplicit_late_drift_through_th1971209900/001/.benc… (904 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetLateDrift,\"clean_set_apply_test.go:204: apply = (\u003cnil>, mutated=true), want the under-lock recheck to refuse the raced member\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalAuthorityPreservation,\"cleanup_members_modes_test.go:150: apply error = \u003cnil>, want the locked recheck to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/locked-drift,\"cleanup_members_baseline_test.go:119: explicit/locked-drift differs from its frozen record at line 15: got \\\"  error none\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",54\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/landed/locked-drift,\"cleanup_members_baseline_test.go:119: landed/locked-drift differs from its frozen record at line 15: got \\\"  error cleanup receipt has invalid identity\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",53\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5e0110372b6adf94156ef2a64ac42935e19632a5",
          "requirement": "locked-requalify",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the locked qualification check of the recorded cleanup transaction and record the R10 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-author-01",
              "digest": "sha256:7c192a48811fefd3562e8dca807f3c4a6a20a986f960d41dca7ea616e249a1c2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/resume.go,swap,failed,5,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,99\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1509\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,5577,99\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,\"clean_set_outcomes_test.go:322: late-drift apply = (0, \\\"worktree_cleanup[2]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/KVVQSZ/t/TestCleanSetApplyTimeStaleRefusalexplicit_late_drift_through_th1971209900/001/.benc… (904 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetLateDrift,\"clean_set_apply_test.go:204: apply = (\u003cnil>, mutated=true), want the under-lock recheck to refuse the raced member\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalAuthorityPreservation,\"cleanup_members_modes_test.go:150: apply error = \u003cnil>, want the locked recheck to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/locked-drift,\"cleanup_members_baseline_test.go:119: explicit/locked-drift differs from its frozen record at line 15: got \\\"  error none\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",54\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/landed/locked-drift,\"cleanup_members_baseline_test.go:119: landed/locked-drift differs from its frozen record at line 15: got \\\"  error cleanup receipt has invalid identity\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",53\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "C1-T01-stop-at-failure",
          "performer": "ft363-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f83eb0463ad351fd34d053f3eaade732acdb79b1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-author-01",
            "digest": "sha256:6df584a8fda30d3e65694d3ff826ba518809579683c723303d6fcc6198f19373",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/cleanup_members.go,swap,failed,12,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,99\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2970\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,5481,99\nfailures[12]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,\"clean_set_outcomes_test.go:322: late-drift apply = (0, \\\"worktree_cleanup[2]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/EYV6UO/t/TestCleanSetApplyTimeStaleRefusalexplicit_late_drift_through_th182124007/001/.bench… (924 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetLateDrift,\"clean_set_apply_test.go:204: apply = (\u003cnil>, mutated=true), want the under-lock recheck to refuse the raced member\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,\"clean_set_refusal_test.go:80: apply error = \u003cnil>, want the member's own requalify to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetPartialApply,\"clean_set_outcomes_test.go:45: apply error = \u003cnil>, want the injected transaction fault\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetUnstartedOutcomes,\"clean_set_outcomes_test.go:79: partial apply = (0, \\\"worktree_cleanup[3]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/EYV6UO/t/TestCleanSetUnstartedOutcomes2527818133/001/.bench-home/worktrees/001-1852413856/95bb73… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalAuthorityPreservation,\"cleanup_members_modes_test.go:150: apply error = \u003cnil>, want the locked recheck to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/first-failure,\"cleanup_members_baseline_test.go:119: explicit/first-failure differs from its frozen record at line 18: got \\\"  error none\\\", frozen \\\"  error baseline transaction fault\\\"\",65\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/later-failure,\"cleanup_members_baseline_test.go:119: explicit/later-failure differs from its frozen record at line 18: got \\\"  error none\\\", frozen \\\"  error baseline transaction fault\\\"\",65\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/locked-drift,\"cleanup_members_baseline_test.go:119: explicit/locked-drift differs from its frozen record at line 15: got \\\"  error none\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",53\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalCompletedPrefix,\"cleanup_members_modes_test.go:81: apply error = \u003cnil>, want the injected fault\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalStopsAtFailedVisit/first,\"cleanup_members_test.go:97: error = \u003cnil>, want the failed visit's own error\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalStopsAtFailedVisit/later,\"cleanup_members_test.go:97: error = \u003cnil>, want the failed visit's own error\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5e0110372b6adf94156ef2a64ac42935e19632a5",
          "requirement": "stop-at-failure",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "exit_code": 0,
          "probe": {
            "mutation": "Continue the shared traversal after a failed visit and record the R14 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-author-01",
              "digest": "sha256:6df584a8fda30d3e65694d3ff826ba518809579683c723303d6fcc6198f19373",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/cleanup_members.go,swap,failed,12,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,99\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2970\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,5481,99\nfailures[12]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,\"clean_set_outcomes_test.go:322: late-drift apply = (0, \\\"worktree_cleanup[2]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/EYV6UO/t/TestCleanSetApplyTimeStaleRefusalexplicit_late_drift_through_th182124007/001/.bench… (924 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetLateDrift,\"clean_set_apply_test.go:204: apply = (\u003cnil>, mutated=true), want the under-lock recheck to refuse the raced member\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,\"clean_set_refusal_test.go:80: apply error = \u003cnil>, want the member's own requalify to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetPartialApply,\"clean_set_outcomes_test.go:45: apply error = \u003cnil>, want the injected transaction fault\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetUnstartedOutcomes,\"clean_set_outcomes_test.go:79: partial apply = (0, \\\"worktree_cleanup[3]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/EYV6UO/t/TestCleanSetUnstartedOutcomes2527818133/001/.bench-home/worktrees/001-1852413856/95bb73… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalAuthorityPreservation,\"cleanup_members_modes_test.go:150: apply error = \u003cnil>, want the locked recheck to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/first-failure,\"cleanup_members_baseline_test.go:119: explicit/first-failure differs from its frozen record at line 18: got \\\"  error none\\\", frozen \\\"  error baseline transaction fault\\\"\",65\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/later-failure,\"cleanup_members_baseline_test.go:119: explicit/later-failure differs from its frozen record at line 18: got \\\"  error none\\\", frozen \\\"  error baseline transaction fault\\\"\",65\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/locked-drift,\"cleanup_members_baseline_test.go:119: explicit/locked-drift differs from its frozen record at line 15: got \\\"  error none\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",53\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalCompletedPrefix,\"cleanup_members_modes_test.go:81: apply error = \u003cnil>, want the injected fault\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalStopsAtFailedVisit/first,\"cleanup_members_test.go:97: error = \u003cnil>, want the failed visit's own error\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalStopsAtFailedVisit/later,\"cleanup_members_test.go:97: error = \u003cnil>, want the failed visit's own error\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "C1-T01-recorded-suffix",
          "performer": "ft363-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f83eb0463ad351fd34d053f3eaade732acdb79b1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-author-01",
            "digest": "sha256:7904ee497271b27c230dd9e6585279ab4442914ae23c43100c61a5598809342d",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,99,1,0,null,yes\n  2,bit,swap,failed,6,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,99,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,99,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "5e0110372b6adf94156ef2a64ac42935e19632a5",
          "requirement": "recorded-suffix",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the unreached recorded members from the explicit stop projection and record the R15 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-author-01",
              "digest": "sha256:7904ee497271b27c230dd9e6585279ab4442914ae23c43100c61a5598809342d",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,99,1,0,null,yes\n  2,bit,swap,failed,6,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,99,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,99,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
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
    },
    {
      "from": "sha256:6326ba026c86a6e50db4f8cfb3e74eec9c5d576e5d39b782ab0072a0fbb3fd77",
      "to": "sha256:3f9255958885eae9a6109137977b57408b8d5287016e33fd22b9d4d306fa75a5",
      "chunk_ids": {
        "C1": [
          "C1"
        ]
      }
    }
  ]
}
```

## C1 author witness reds

The ticket 01 author reported these development reds outside the plan probes. Each mutation was restored, and the package then passed.

| Row | Mutation | Red test |
|---|---|---|
| R01, R05 | The visit predicate is forced false. | `TestCleanupTraversalEmptyModes/explicit/retained-only`, `TestCleanupTraversalRetainedVisits` |
| R02 | Each outcome is prepended, not appended. | `TestCleanupTraversalOrder` |
| R19 | The stop projection has no completed prefix. | `TestCleanupTraversalCompletedPrefix` |
| R27 | The `StepMemberRequalify` hit is omitted. | `TestCleanupTraversalModeFaultBoundaries/explicit/all-success` |
| R28 | The locked fingerprint check is disabled. | `TestCleanupTraversalAuthorityPreservation` |

The author reported one gap: the only cited R13 witness, the subtest for drift that stays removable, drives landed mode only.
