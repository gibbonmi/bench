# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/cleanup-member-traversal/spec.md",
  "plan_digest": "sha256:8917268dbf6107a6d8cda61954f441bb91d3430e791577ba9c7ee427bc99837c",
  "implementation_session": "",
  "chunks": [
    {
      "id": "C1",
      "base": "76c8d256cd35d5c74fba8d9273d2bfcb51f89a71",
      "tip": "94eb8696d87a613cebca4cefdcffea5b37389019",
      "plan_digest": "sha256:8917268dbf6107a6d8cda61954f441bb91d3430e791577ba9c7ee427bc99837c",
      "source_digest": "23e60de8115a289143ee457b2b5ebe265a9e5923",
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
        },
        {
          "id": "C1-T01-R1-explicit-preflight",
          "performer": "ft363-repair-01-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "23e60de8115a289143ee457b2b5ebe265a9e5923",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-repair-01-r1",
            "digest": "sha256:4e37adc266d9a56817a2bf680a14ddfacd2ef11cafa4e670d7ac245759a82d17",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,100,1,0,null,yes\n  2,bit,swap,failed,7,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,100,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,100,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "94eb8696d87a613cebca4cefdcffea5b37389019",
          "requirement": "explicit-preflight",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the explicit all-member preflight before the recorded visits and record the R06 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-repair-01-r1",
              "digest": "sha256:4e37adc266d9a56817a2bf680a14ddfacd2ef11cafa4e670d7ac245759a82d17",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,100,1,0,null,yes\n  2,bit,swap,failed,7,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,100,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,100,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "C1-T01-R1-later-requalify",
          "performer": "ft363-repair-01-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "23e60de8115a289143ee457b2b5ebe265a9e5923",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-repair-01-r1",
            "digest": "sha256:4e37adc266d9a56817a2bf680a14ddfacd2ef11cafa4e670d7ac245759a82d17",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,100,1,0,null,yes\n  2,bit,swap,failed,7,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,100,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,100,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "94eb8696d87a613cebca4cefdcffea5b37389019",
          "requirement": "later-requalify",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the later requalification before a reached recorded visit and record the R09 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-repair-01-r1",
              "digest": "sha256:4e37adc266d9a56817a2bf680a14ddfacd2ef11cafa4e670d7ac245759a82d17",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,100,1,0,null,yes\n  2,bit,swap,failed,7,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,100,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,100,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "C1-T01-R1-locked-requalify",
          "performer": "ft363-repair-01-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "23e60de8115a289143ee457b2b5ebe265a9e5923",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-repair-01-r1",
            "digest": "sha256:d7cc79fa093eaeb4c7925c2c70396f722fbb282d2200f4e83354c4d886bf8795",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/resume.go,swap,failed,5,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,100\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1510\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,5644,100\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,\"clean_set_outcomes_test.go:322: late-drift apply = (0, \\\"worktree_cleanup[2]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/LIQV7U/t/TestCleanSetApplyTimeStaleRefusalexplicit_late_drift_through_th596009054/001/.bench… (902 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetLateDrift,\"clean_set_apply_test.go:204: apply = (\u003cnil>, mutated=true), want the under-lock recheck to refuse the raced member\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalAuthorityPreservation,\"cleanup_members_modes_test.go:150: apply error = \u003cnil>, want the locked recheck to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/locked-drift,\"cleanup_members_baseline_test.go:119: explicit/locked-drift differs from its frozen record at line 15: got \\\"  error none\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",54\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/landed/locked-drift,\"cleanup_members_baseline_test.go:119: landed/locked-drift differs from its frozen record at line 15: got \\\"  error cleanup receipt has invalid identity\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",53\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "94eb8696d87a613cebca4cefdcffea5b37389019",
          "requirement": "locked-requalify",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the locked qualification check of the recorded cleanup transaction and record the R10 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-repair-01-r1",
              "digest": "sha256:d7cc79fa093eaeb4c7925c2c70396f722fbb282d2200f4e83354c4d886bf8795",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/resume.go,swap,failed,5,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,100\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1510\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,5644,100\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,\"clean_set_outcomes_test.go:322: late-drift apply = (0, \\\"worktree_cleanup[2]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/LIQV7U/t/TestCleanSetApplyTimeStaleRefusalexplicit_late_drift_through_th596009054/001/.bench… (902 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetLateDrift,\"clean_set_apply_test.go:204: apply = (\u003cnil>, mutated=true), want the under-lock recheck to refuse the raced member\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalAuthorityPreservation,\"cleanup_members_modes_test.go:150: apply error = \u003cnil>, want the locked recheck to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/locked-drift,\"cleanup_members_baseline_test.go:119: explicit/locked-drift differs from its frozen record at line 15: got \\\"  error none\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",54\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/landed/locked-drift,\"cleanup_members_baseline_test.go:119: landed/locked-drift differs from its frozen record at line 15: got \\\"  error cleanup receipt has invalid identity\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",53\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "C1-T01-R1-stop-at-failure",
          "performer": "ft363-repair-01-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "23e60de8115a289143ee457b2b5ebe265a9e5923",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-repair-01-r1",
            "digest": "sha256:4c5fb9d0d78cd4e110dda43d6efb7a951ed0557aa661ff19887a60acc99a1e7b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/cleanup_members.go,swap,failed,13,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,100\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,3198\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,5566,100\nfailures[13]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,\"clean_set_outcomes_test.go:322: late-drift apply = (0, \\\"worktree_cleanup[2]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/6DPH3O/t/TestCleanSetApplyTimeStaleRefusalexplicit_late_drift_through_th3012840869/001/.benc… (926 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetLateDrift,\"clean_set_apply_test.go:204: apply = (\u003cnil>, mutated=true), want the under-lock recheck to refuse the raced member\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,\"clean_set_refusal_test.go:80: apply error = \u003cnil>, want the member's own requalify to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetMemberDriftAfterPreflight/explicit_member_drifts_while_still_removable,\"clean_set_refusal_test.go:135: apply error = \u003cnil>, want the member's own requalify to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetPartialApply,\"clean_set_outcomes_test.go:45: apply error = \u003cnil>, want the injected transaction fault\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetUnstartedOutcomes,\"clean_set_outcomes_test.go:79: partial apply = (0, \\\"worktree_cleanup[3]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/6DPH3O/t/TestCleanSetUnstartedOutcomes2585395255/001/.bench-home/worktrees/001-959000473/b5de0a0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalAuthorityPreservation,\"cleanup_members_modes_test.go:150: apply error = \u003cnil>, want the locked recheck to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/first-failure,\"cleanup_members_baseline_test.go:119: explicit/first-failure differs from its frozen record at line 18: got \\\"  error none\\\", frozen \\\"  error baseline transaction fault\\\"\",65\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/later-failure,\"cleanup_members_baseline_test.go:119: explicit/later-failure differs from its frozen record at line 18: got \\\"  error none\\\", frozen \\\"  error baseline transaction fault\\\"\",65\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/locked-drift,\"cleanup_members_baseline_test.go:119: explicit/locked-drift differs from its frozen record at line 15: got \\\"  error none\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",53\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalCompletedPrefix,\"cleanup_members_modes_test.go:81: apply error = \u003cnil>, want the injected fault\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalStopsAtFailedVisit/first,\"cleanup_members_test.go:97: error = \u003cnil>, want the failed visit's own error\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalStopsAtFailedVisit/later,\"cleanup_members_test.go:97: error = \u003cnil>, want the failed visit's own error\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "94eb8696d87a613cebca4cefdcffea5b37389019",
          "requirement": "stop-at-failure",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "exit_code": 0,
          "probe": {
            "mutation": "Continue the shared traversal after a failed visit and record the R14 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-repair-01-r1",
              "digest": "sha256:4c5fb9d0d78cd4e110dda43d6efb7a951ed0557aa661ff19887a60acc99a1e7b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/cleanup_members.go,swap,failed,13,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,100\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,3198\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,5566,100\nfailures[13]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetApplyTimeStaleRefusal/explicit_late_drift_through_the_command,\"clean_set_outcomes_test.go:322: late-drift apply = (0, \\\"worktree_cleanup[2]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/6DPH3O/t/TestCleanSetApplyTimeStaleRefusalexplicit_late_drift_through_th3012840869/001/.benc… (926 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetLateDrift,\"clean_set_apply_test.go:204: apply = (\u003cnil>, mutated=true), want the under-lock recheck to refuse the raced member\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,\"clean_set_refusal_test.go:80: apply error = \u003cnil>, want the member's own requalify to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetMemberDriftAfterPreflight/explicit_member_drifts_while_still_removable,\"clean_set_refusal_test.go:135: apply error = \u003cnil>, want the member's own requalify to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetPartialApply,\"clean_set_outcomes_test.go:45: apply error = \u003cnil>, want the injected transaction fault\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanSetUnstartedOutcomes,\"clean_set_outcomes_test.go:79: partial apply = (0, \\\"worktree_cleanup[3]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /dev/shm/6DPH3O/t/TestCleanSetUnstartedOutcomes2585395255/001/.bench-home/worktrees/001-959000473/b5de0a0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalAuthorityPreservation,\"cleanup_members_modes_test.go:150: apply error = \u003cnil>, want the locked recheck to refuse\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/first-failure,\"cleanup_members_baseline_test.go:119: explicit/first-failure differs from its frozen record at line 18: got \\\"  error none\\\", frozen \\\"  error baseline transaction fault\\\"\",65\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/later-failure,\"cleanup_members_baseline_test.go:119: explicit/later-failure differs from its frozen record at line 18: got \\\"  error none\\\", frozen \\\"  error baseline transaction fault\\\"\",65\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalBaseline/explicit/locked-drift,\"cleanup_members_baseline_test.go:119: explicit/locked-drift differs from its frozen record at line 15: got \\\"  error none\\\", frozen \\\"  error cleanup fingerprint is stale\\\"\",53\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalCompletedPrefix,\"cleanup_members_modes_test.go:81: apply error = \u003cnil>, want the injected fault\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalStopsAtFailedVisit/first,\"cleanup_members_test.go:97: error = \u003cnil>, want the failed visit's own error\",1\n  github.com/gibbonmi/bench/internal/worktree,TestCleanupTraversalStopsAtFailedVisit/later,\"cleanup_members_test.go:97: error = \u003cnil>, want the failed visit's own error\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "C1-T01-R1-recorded-suffix",
          "performer": "ft363-repair-01-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "23e60de8115a289143ee457b2b5ebe265a9e5923",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-repair-01-r1",
            "digest": "sha256:4e37adc266d9a56817a2bf680a14ddfacd2ef11cafa4e670d7ac245759a82d17",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,100,1,0,null,yes\n  2,bit,swap,failed,7,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,100,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,100,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "94eb8696d87a613cebca4cefdcffea5b37389019",
          "requirement": "recorded-suffix",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the unreached recorded members from the explicit stop projection and record the R15 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft363-repair-01-r1",
              "digest": "sha256:4e37adc266d9a56817a2bf680a14ddfacd2ef11cafa4e670d7ac245759a82d17",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,2,TestCleanSetPreflightAllRows/explicit,100,1,0,null,yes\n  2,bit,swap,failed,7,TestCleanSetMemberDriftAfterPreflight/drifts_to_retained,100,1,0,null,yes\n  3,bit,swap,failed,5,TestCleanSetPartialApply,100,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "C1-R1-STD",
          "performer": "ft363-review-c1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f83eb0463ad351fd34d053f3eaade732acdb79b1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-review-c1",
            "digest": "sha256:2b8c0eb1e1d37da894c55a97d11e4b02c02a713b03ab3f01d143d1817731fd72",
            "excerpt": "## Standards\nFindings: 0. Worst issue: none.\nThe shared owner at internal/worktree/cleanup_members.go:97-125 handles only iteration, pass-through, accumulation and the stop. It derives no fingerprint, identity or classification, and it takes no run. The explicit adapter supplies the predicate, the visit, the retained row and the stop projection; its visit closure captures a.run (internal/worktree/clean_set.go:333-353). The test-count pin went from 874 to 882 (internal/worktree/parallel_census_test.go:926), matching 8 new top-level tests.\n"
          },
          "axis": "Standards",
          "base": "76c8d256cd35d5c74fba8d9273d2bfcb51f89a71",
          "tip": "5e0110372b6adf94156ef2a64ac42935e19632a5",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "C1-R1-SPEC",
          "performer": "ft363-review-c1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f83eb0463ad351fd34d053f3eaade732acdb79b1",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft363-review-c1",
            "digest": "sha256:c312b7b875a0f49780143f7d639dfcf58c4d7f77bb1f41492042536a37a322c0",
            "excerpt": "## Spec\nFindings: 1. Worst issue: C1-SPEC-1.\nC1-SPEC-1 | sev 4 | conf 8 | ask-user | internal/worktree/cleanup_members_baseline_record_test.go:131-134, internal/worktree/cleanup_members_baseline_test.go:146-153\nRequirement: spec.md:325-327 \"Verify that each fingerprint binds the same inputs within its own run\" and \"A changed digest relation, input omission, ... must fail the comparison.\" Defect: digests are named by role and nothing derives a digest from its rendered inputs. Probe: removed options.Full from fingerprintExplicitSet (clean_set.go:239), ran CleanupTraversalBaseline: silent.\nRow audit: R01, R02, R04, R05, R08, R11, R12, R14, R15, R19, R24, R28 witnessed in explicit mode; R27 explicit only; R03 held by the unchanged successor; R06, R09, R10 rest on the recorded plan probes; R26 baseline part has 31 cases, ticket 00 is test-only and ticket 01 left the records unchanged; R13 partial (C1-COV-1).\n"
          },
          "axis": "Spec",
          "base": "76c8d256cd35d5c74fba8d9273d2bfcb51f89a71",
          "tip": "5e0110372b6adf94156ef2a64ac42935e19632a5",
          "finding_ids": [
            "C1-SPEC-1"
          ],
          "supersedes": []
        },
        {
          "id": "C1-R1-COV",
          "performer": "ft363-review-c1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f83eb0463ad351fd34d053f3eaade732acdb79b1",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft363-review-c1",
            "digest": "sha256:ac970814bd12a17e0f32c39c5cab91395d0d79618380bfbd1f519f11b2bda6d3",
            "excerpt": "## Coverage\nFindings: 1. Worst issue: C1-COV-1.\nC1-COV-1 | sev 4 | conf 7 | auto-fix | internal/worktree/clean_set.go:341; witness internal/worktree/clean_set_refusal_test.go:89-113\nInput: in explicit mode, a member that drifts after preflight but stays removable, at its in-loop requalify. Gap: the only R13 witness drives applyLandedSet. Probe: inserted `if current.plan.Action.Removes() { return current.plan, err }` before line 341, ran CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval: 99 tests passed, silent. Fix: add an explicit subtest to TestCleanSetMemberDriftAfterPreflight.\nProbes: explicit visits predicate forced true: bit (7 failures). Stop projection members[failed+1:] to members[failed+2:]: bit (TestCleanSetPartialApply). Final git status: clean.\n"
          },
          "axis": "Coverage",
          "base": "76c8d256cd35d5c74fba8d9273d2bfcb51f89a71",
          "tip": "5e0110372b6adf94156ef2a64ac42935e19632a5",
          "finding_ids": [
            "C1-COV-1"
          ],
          "supersedes": []
        },
        {
          "id": "C1-R2-STD",
          "performer": "ft363-review-c1-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "23e60de8115a289143ee457b2b5ebe265a9e5923",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-review-c1-r2",
            "digest": "sha256:b8373c36c7367562dc404a81f0ccceb5a4b4290cec2cfaad15029ea9b5e59271",
            "excerpt": "## Standards\nOutcome: closed. Findings: 0.\nThe test no longer holds its own copy of the fingerprint inputs; it renders explicitFingerprintInputs (cleanup_members_baseline_test.go:148-161), which fingerprintExplicitSet hashes (clean_set.go:173, 239-263). The frozen records did not change. requalifyExplicitRow moved byte-identical to clean_set_apply.go:229-241; call sites unchanged.\n"
          },
          "axis": "Standards",
          "base": "76c8d256cd35d5c74fba8d9273d2bfcb51f89a71",
          "tip": "94eb8696d87a613cebca4cefdcffea5b37389019",
          "finding_ids": [],
          "supersedes": [
            "C1-R1-STD"
          ]
        },
        {
          "id": "C1-R2-SPEC",
          "performer": "ft363-review-c1-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "23e60de8115a289143ee457b2b5ebe265a9e5923",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-review-c1-r2",
            "digest": "sha256:144e5060ceaaa6e39fa14c07162340880c91cb335e4452feb414f424daf17249",
            "excerpt": "## Spec (C1-SPEC-1)\nOutcome: closed. Findings: 0.\nOmitting an input now fails the baseline (cleanup_members_baseline_test.go:118, spec.md:327). Each plan digest must equal fingerprintExplicitSet of the rendered inputs in the same run (cleanup_members_baseline_test.go:149-151, spec.md:325). Probe: removed options.Full at clean_set.go:240: bit, 15 failures, first TestCleanupTraversalBaseline/explicit/all-success, restored yes.\n"
          },
          "axis": "Spec",
          "base": "76c8d256cd35d5c74fba8d9273d2bfcb51f89a71",
          "tip": "94eb8696d87a613cebca4cefdcffea5b37389019",
          "finding_ids": [],
          "supersedes": [
            "C1-R1-SPEC"
          ]
        },
        {
          "id": "C1-R2-COV",
          "performer": "ft363-review-c1-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "23e60de8115a289143ee457b2b5ebe265a9e5923",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft363-review-c1-r2",
            "digest": "sha256:ff907a100723b64b530737df2d8661e7eb8298d6358d792584cf837133850120",
            "excerpt": "## Coverage (C1-COV-1)\nOutcome: closed. Findings: 0.\nThe new subtest (clean_set_refusal_test.go:116-143) drifts the second member at its own in-loop requalify (clean_set.go:337-340) and asserts errStaleFingerprint, a not-attempted row, and that the member survives. Probe: inserted `if current.plan.Action.Removes() { err = nil }` before clean_set.go:341: bit, 1 failure, explicit_member_drifts_while_still_removable at clean_set_refusal_test.go:135, restored yes. Final git status: clean.\n"
          },
          "axis": "Coverage",
          "base": "76c8d256cd35d5c74fba8d9273d2bfcb51f89a71",
          "tip": "94eb8696d87a613cebca4cefdcffea5b37389019",
          "finding_ids": [],
          "supersedes": [
            "C1-R1-COV"
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
    },
    {
      "from": "sha256:3f9255958885eae9a6109137977b57408b8d5287016e33fd22b9d4d306fa75a5",
      "to": "sha256:8917268dbf6107a6d8cda61954f441bb91d3430e791577ba9c7ee427bc99837c",
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

## C1 review round 1

One unified opus / medium session reviewed the frozen pair `76c8d256c..5e0110372`. Repair cycles consumed: 0 of 2 before this round's repair.

## Standards

Findings: 0. Worst issue: none.

## Spec

Findings: 1. Worst issue: C1-SPEC-1.

- C1-SPEC-1, severity 4, confidence 8. The baseline does not fail when a fingerprint input is omitted (`internal/worktree/cleanup_members_baseline_test.go:146`). The test holds a hand copy of the input list of `fingerprintExplicitSet`. Disposition: auto-fix, by orchestrator decision under `--auto-approve`, flagged for reviewer veto. The repair moves the explicit input list into one production function that the fingerprint hashes, and the baseline renders that list from production.

## Coverage

Findings: 1. Worst issue: C1-COV-1.

- C1-COV-1, severity 4, confidence 7. Explicit mode has no witness for drift that stays removable at the in-loop requalification (`internal/worktree/clean_set.go:341`). Disposition: auto-fix. The repair adds an explicit subtest to `TestCleanSetMemberDriftAfterPreflight`.

## Advice

- With `BENCH_CAPTURE_CLEANUP_BASELINE=1`, the baseline test rewrites the records and still passes. A failed test after a capture write is safer.
- Cases that name members by label sort their rows, so row order is checked only against the planned order of the current run.
- The commits section collapses member tips to one rank name, so the member-to-commit relation is not compared.

Both findings route to one repair session on ticket 01, which holds every path of both fixes.

## C1 repair 1 witness reds

The repair session ft363-repair-01-r1 closed both findings at repair commit `94eb8696d`. Repair cycles consumed: 1 of 2. Each mutation was restored, and the package then passed.

| Finding | Row | Mutation | Red test |
|---|---|---|---|
| C1-COV-1 | R13 | The requalification returns the current plan when it still removes. | `TestCleanSetMemberDriftAfterPreflight/explicit_member_drifts_while_still_removable` |
| C1-SPEC-1 | R26 | `options.Full` is removed from `explicitFingerprintInputs`. | `TestCleanupTraversalBaseline`, 15 failures |

The repair moved `requalifyExplicitRow` unchanged from `clean_set.go` to `clean_set_apply.go`, so that `clean_set.go` stays within the 400-line bound.

## C1 confirming round

A fresh unified opus / medium session read the repair delta `5e0110372..94eb8696d`. It closed C1-SPEC-1 and C1-COV-1, and it found no new finding. Each axis replayed its finding's probe, and each probe bit. Repair cycles consumed for C1: 1 of 2.
