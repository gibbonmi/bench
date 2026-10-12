# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/build-landing-window/spec.md",
  "plan_digest": "sha256:e5fa3d7f7092f75a9f8bef63b4fd638d15330f29195691d7945b1fecd38042a3",
  "implementation_session": "",
  "chunks": [
    {
      "id": "BW-A",
      "base": "dcbbfb805419eef05174ea379b99547cfc8889da",
      "tip": "ced1dedc83e74057e0af3499d3b3ccdca9930319",
      "plan_digest": "sha256:e5fa3d7f7092f75a9f8bef63b4fd638d15330f29195691d7945b1fecd38042a3",
      "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
      "acceptance_rows": [
        "BW1",
        "BW2",
        "BW3",
        "BW4",
        "BW5",
        "BW6",
        "BW7",
        "BW8",
        "BW9",
        "BW10",
        "BW11",
        "BW12",
        "BW13",
        "BW14",
        "BW15",
        "BW16"
      ],
      "verification": [
        {
          "id": "t1-land-focused-v1",
          "performer": "claude:blw_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1/v-t1-land-focused.txt",
            "digest": "sha256:63c43d629e1c698817ed540e552543d2cd4794c19d3b61be2355968acb00bfad",
            "excerpt": "bench test --package ./internal/worktree --run 'TestLandHoldsTheDestinationLease|TestSecondLandingMeetsTheLease|TestLandReleasesTheLeaseOnRefusal|TestLandResumeTakesNoLease|TestLeaseLockUnavailable|TestLandingRefusalRegistryHasAProducingFixture'\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 6 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1973,20\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nplan probe: see probe-t1-land-focused.txt (verdict bit, restored yes)\n"
          },
          "requirement": "t1-land-focused",
          "command": "bench test --package ./internal/worktree --run 'TestLandHoldsTheDestinationLease|TestSecondLandingMeetsTheLease|TestLandReleasesTheLeaseOnRefusal|TestLandResumeTakesNoLease|TestLeaseLockUnavailable|TestLandingRefusalRegistryHasAProducingFixture'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the lease acquire call that the landing makes after its admission proofs and before landReviewed, so that a second landing reaches composition and the gate lock. TestSecondLandingMeetsTheLease turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1/v-t1-land-focused.txt",
              "digest": "sha256:63c43d629e1c698817ed540e552543d2cd4794c19d3b61be2355968acb00bfad",
              "excerpt": "bench test --package ./internal/worktree --run 'TestLandHoldsTheDestinationLease|TestSecondLandingMeetsTheLease|TestLandReleasesTheLeaseOnRefusal|TestLandResumeTakesNoLease|TestLeaseLockUnavailable|TestLandingRefusalRegistryHasAProducingFixture'\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 6 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1973,20\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nplan probe: see probe-t1-land-focused.txt (verdict bit, restored yes)\n"
            }
          }
        },
        {
          "id": "t1-lease-focused-v1",
          "performer": "claude:blw_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1/v-t1-lease-focused.txt",
            "digest": "sha256:d87c62087763a6f9b4a463f44a0d5266494efa4c9869016460ea926ae55ba08c",
            "excerpt": "bench test --package ./internal/commitment/repository --run 'TestPublishReleasesTheLeaseInItsDecision|TestLeaseTakesOverADeadOwner'\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 2 s\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,108,2\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t1-lease-focused",
          "command": "bench test --package ./internal/commitment/repository --run 'TestPublishReleasesTheLeaseInItsDecision|TestLeaseTakesOverADeadOwner'",
          "exit_code": 0
        },
        {
          "id": "t1-records-focused-v1",
          "performer": "claude:blw_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1/v-t1-records-focused.txt",
            "digest": "sha256:2ed3fbbf239305962c40ee9a03ae4d9120aea628990d60d3aee8f1d775cabb28",
            "excerpt": "bench test --package ./internal/intent --run 'TestValidateCommitmentRefusesDuplicateBuildRecords|TestTolerantReadKeepsTheBuildRecords|TestIntentReExportsEveryMovedLedgerName'\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 4 s\n  github.com/gibbonmi/bench/internal/intent,pass,10,9\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t1-records-focused",
          "command": "bench test --package ./internal/intent --run 'TestValidateCommitmentRefusesDuplicateBuildRecords|TestTolerantReadKeepsTheBuildRecords|TestIntentReExportsEveryMovedLedgerName'",
          "exit_code": 0
        },
        {
          "id": "t1-inventory-focused-v1",
          "performer": "claude:blw_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1/v-t1-inventory-focused.txt",
            "digest": "sha256:ee840d4ca4fd7bfb1d080781e87cb4956c03b5219626b43487131513fd8d1fb8",
            "excerpt": "bench test --package ./internal/conformance --run 'TestCommitmentRecordInventoryNamesTheBuildRecords'\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 6 s\n  github.com/gibbonmi/bench/internal/conformance,pass,12,1\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t1-inventory-focused",
          "command": "bench test --package ./internal/conformance --run 'TestCommitmentRecordInventoryNamesTheBuildRecords'",
          "exit_code": 0
        },
        {
          "id": "t1-floor-v1",
          "performer": "claude:blw_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1/v-t1-floor.txt",
            "digest": "sha256:2942764d0eb85418fcab370548a4d9bf5c8a77a15cdad92fe90b428bd3bd5bfe",
            "excerpt": "bench test --package ./internal/intent/... --package ./internal/commitment/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/gate/... --package ./internal/gocache --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 105 s (final tree)\n  github.com/gibbonmi/bench/internal/anchors,pass,3028,179\n  github.com/gibbonmi/bench/internal/commitment,pass,6327,175\n  github.com/gibbonmi/bench/internal/commitment/commitcmd,pass,831,19\n  github.com/gibbonmi/bench/internal/commitment/commitmenttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,5903,109\n  github.com/gibbonmi/bench/internal/conformance,pass,83571,1690\n  github.com/gibbonmi/bench/internal/gate,pass,32979,430\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,272,11\n  github.com/gibbonmi/bench/internal/gate/greenmarker,pass,69,4\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,345,36\n  github.com/gibbonmi/bench/internal/gocache,pass,216,42\n  github.com/gibbonmi/bench/internal/intent,pass,1212,64\n  github.com/gibbonmi/bench/internal/intent/admissionpolicy,pass,3,28\n  github.com/gibbonmi/bench/internal/intent/ledger,pass,3,1\n  github.com/gibbonmi/bench/internal/refusalroute,pass,5,17\n  github.com/gibbonmi/bench/internal/worktree,pass,77615,1406\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-floor",
          "command": "bench test --package ./internal/intent/... --package ./internal/commitment/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/gate/... --package ./internal/gocache --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-data-handling-v1",
          "performer": "claude:blw_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1/v-t1-data-handling.txt",
            "digest": "sha256:aaad682a24ff9d501293414f1a5563691db05268906749794f33f791bee1c621",
            "excerpt": "bench test --check data-handling-derivation\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 2 s (after bench worktree build FT342-FT253)\n  data-handling-derivation,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,10,1\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t1-data-handling",
          "command": "bench test --check data-handling-derivation",
          "exit_code": 0
        },
        {
          "id": "t1-system-lease-v1",
          "performer": "claude:blw_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1/v-t1-system-lease.txt",
            "digest": "sha256:d209209b981eec9eb346d620a129087b936ccf3adb0d94c87cfcbad7b491c58d",
            "excerpt": "bench test --check system --run TestTwoLandingProcessesShareOneLease\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 5 s\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\n  github.com/gibbonmi/bench/internal/systemtest,pass,1215,1\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t1-system-lease",
          "command": "bench test --check system --run TestTwoLandingProcessesShareOneLease",
          "exit_code": 0
        },
        {
          "id": "t1-root-conformance-v1",
          "performer": "claude:blw_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1/v-t1-root-conformance.txt",
            "digest": "sha256:6fd13d0bfc568c4da202f5afe096db305acfbe11c90d6f8cdf2b45b083733400",
            "excerpt": "bench test --check root-conformance\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 24 s (after bench worktree build FT342-FT253)\n  root-conformance,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,18755,1\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
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
