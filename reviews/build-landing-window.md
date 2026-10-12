# Review outcomes

BW-A round 1 reviewed the pair `dcbbfb80..ced1dedc`.
The repair allowance of BW-A is two cycles, and the build has used zero of them.

## Standards

Findings: 3. Worst issue: S1, a copied gate-barrier fixture harness.

- S1, auto-fix, confidence 9: `barrierLandingFixture` in `internal/worktree/land_lease_test.go:47` copies the FIFO gate barrier of `TestCommitmentGateRace` (`internal/worktree/land_identity_test.go:241-264`). Fold that test onto the one harness.
- S2, auto-fix, confidence 8: `systemReleaseGateBarrier` in `internal/systemtest/owner_land_race_test.go:182` repeats the body of `releaseArtifactBarrier` (`internal/systemtest/owner_artifact_recovery_test.go:147`). Call the existing helper.
- S3, auto-fix, confidence 8: six new comment sentences are over 25 words. They are in `lease.go:17`, `lease_test.go:26`, `build_records_test.go:18` and `:53`, `faces_land.go:3`, and `land_lease_test.go:104`.

## Spec

Findings: 1. Worst issue: P1, a false comparison clause in the Edge inventory.

- P1, auto-fix, confidence 9: the Edge inventory says that the lease lock opens without following a link "as the gate lock opens its file". The gate lock uses a plain open, so remove that clause.

## Coverage

Findings: 3. Worst issue: C1, a BW7 test that does not detect a second transaction.

- C1, auto-fix, confidence 8: the BW7 test passes when `PublishAdmitted` drops the lease in a second transaction. Make the competing acquire take the intent lock first, and show it red under that mutation.
- C2, auto-fix, confidence 7: no test makes the ledger write fail after the ref update. Add a test that asserts nil and a lease record that the next acquire takes over.
- C3, auto-fix, confidence 7: no test puts a symbolic link at the lease lock path. Add a refusal case beside BW15 that expects `landing lease unavailable:`.

```bench-review-record
{
  "version": 2,
  "spec": "specs/build-landing-window/spec.md",
  "plan_digest": "sha256:9c7408fd49da2cbcb5c1e08c68770ed29e667560eecc2524aff5518d31fcce3a",
  "implementation_session": "",
  "chunks": [
    {
      "id": "BW-A",
      "base": "dcbbfb805419eef05174ea379b99547cfc8889da",
      "tip": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
      "plan_digest": "sha256:9c7408fd49da2cbcb5c1e08c68770ed29e667560eecc2524aff5518d31fcce3a",
      "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
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
        },
        {
          "id": "t1-land-focused-v2",
          "performer": "claude:blw_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1rep1/v-t1-land-focused.txt",
            "digest": "sha256:d4aa3d17f192b302eeb05f1863e858718835a475b7b4e7525eeed0c7129ddcc6",
            "excerpt": "bench test --package ./internal/worktree --run 'TestLandHoldsTheDestinationLease|TestSecondLandingMeetsTheLease|TestLandReleasesTheLeaseOnRefusal|TestLandResumeTakesNoLease|TestLeaseLockUnavailable|TestLandingRefusalRegistryHasAProducingFixture'\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed 6 s (repair tree before commit; its content is the committed 2aa2bf49, with no edit between)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,2122,23\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nplan probe: see probe-t1-land-focused.txt (verdict bit, 6 failed tests including TestSecondLandingMeetsTheLease, restored yes)\n"
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
              "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1rep1/v-t1-land-focused.txt",
              "digest": "sha256:d4aa3d17f192b302eeb05f1863e858718835a475b7b4e7525eeed0c7129ddcc6",
              "excerpt": "bench test --package ./internal/worktree --run 'TestLandHoldsTheDestinationLease|TestSecondLandingMeetsTheLease|TestLandReleasesTheLeaseOnRefusal|TestLandResumeTakesNoLease|TestLeaseLockUnavailable|TestLandingRefusalRegistryHasAProducingFixture'\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed 6 s (repair tree before commit; its content is the committed 2aa2bf49, with no edit between)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,2122,23\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nplan probe: see probe-t1-land-focused.txt (verdict bit, 6 failed tests including TestSecondLandingMeetsTheLease, restored yes)\n"
            }
          }
        },
        {
          "id": "t1-lease-focused-v2",
          "performer": "claude:blw_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1rep1/v-t1-lease-focused.txt",
            "digest": "sha256:ff93d48e6bf1a31043518e70c2687ab33860d5dc872e230e608cbc43b9686338",
            "excerpt": "bench test --package ./internal/commitment/repository --run 'TestPublishReleasesTheLeaseInItsDecision|TestLeaseTakesOverADeadOwner'\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed 2 s (repair tree before commit; its content is the committed 2aa2bf49, with no edit between)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,114,2\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nC1 probe (second transaction): see probe-c1-second-transaction.txt (bit, restored yes)\nC2 test TestPublishKeepsTheLeaseWhenTheLedgerWriteFails runs under t1-floor (repository pass, 110 tests); C2 probe: see probe-c2-published-guard.txt (bit, restored yes)\n"
          },
          "requirement": "t1-lease-focused",
          "command": "bench test --package ./internal/commitment/repository --run 'TestPublishReleasesTheLeaseInItsDecision|TestLeaseTakesOverADeadOwner'",
          "exit_code": 0
        },
        {
          "id": "t1-records-focused-v2",
          "performer": "claude:blw_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1rep1/v-t1-records-focused.txt",
            "digest": "sha256:830f7aee721dfb104f6790befb246f803e14627544cc6b8a5f8a824e8d74a7ec",
            "excerpt": "bench test --package ./internal/intent --run 'TestValidateCommitmentRefusesDuplicateBuildRecords|TestTolerantReadKeepsTheBuildRecords|TestIntentReExportsEveryMovedLedgerName'\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed 2 s (repair tree before commit; its content is the committed 2aa2bf49, with no edit between)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,pass,10,9\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-records-focused",
          "command": "bench test --package ./internal/intent --run 'TestValidateCommitmentRefusesDuplicateBuildRecords|TestTolerantReadKeepsTheBuildRecords|TestIntentReExportsEveryMovedLedgerName'",
          "exit_code": 0
        },
        {
          "id": "t1-inventory-focused-v2",
          "performer": "claude:blw_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1rep1/v-t1-inventory-focused.txt",
            "digest": "sha256:c9cc7e048cd8141b3e9516d062f5aaba8918b8ef71e55a1252c91ae9ef9069bb",
            "excerpt": "bench test --package ./internal/conformance --run 'TestCommitmentRecordInventoryNamesTheBuildRecords'\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed 2 s (repair tree before commit; its content is the committed 2aa2bf49, with no edit between)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,10,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-inventory-focused",
          "command": "bench test --package ./internal/conformance --run 'TestCommitmentRecordInventoryNamesTheBuildRecords'",
          "exit_code": 0
        },
        {
          "id": "t1-floor-v2",
          "performer": "claude:blw_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1rep1/v-t1-floor.txt",
            "digest": "sha256:ca5d5e99fa8fde73ed10de552247ef050bd596088e339e47afe910085b26a828",
            "excerpt": "bench test --package ./internal/intent/... --package ./internal/commitment/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/gate/... --package ./internal/gocache --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed 105 s (repair tree before commit; its content is the committed 2aa2bf49, with no edit between)\npackages[17]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,24270,333\n  github.com/gibbonmi/bench/internal/anchors,pass,2305,179\n  github.com/gibbonmi/bench/internal/commitment,pass,5782,175\n  github.com/gibbonmi/bench/internal/commitment/commitcmd,pass,823,19\n  github.com/gibbonmi/bench/internal/commitment/commitmenttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,5547,110\n  github.com/gibbonmi/bench/internal/conformance,pass,85797,1690\n  github.com/gibbonmi/bench/internal/gate,pass,32483,430\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,189,11\n  github.com/gibbonmi/bench/internal/gate/greenmarker,pass,47,4\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,312,36\n  github.com/gibbonmi/bench/internal/gocache,pass,215,42\n  github.com/gibbonmi/bench/internal/intent,pass,1149,64\n  github.com/gibbonmi/bench/internal/intent/admissionpolicy,pass,3,28\n  github.com/gibbonmi/bench/internal/intent/ledger,pass,2,1\n  github.com/gibbonmi/bench/internal/refusalroute,pass,5,17\n  github.com/gibbonmi/bench/internal/worktree,pass,78449,1409\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-floor",
          "command": "bench test --package ./internal/intent/... --package ./internal/commitment/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/gate/... --package ./internal/gocache --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-data-handling-v2",
          "performer": "claude:blw_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1rep1/v-t1-data-handling.txt",
            "digest": "sha256:0fd72940a5008f19b4c9ca6bbf2082ff89a538ea39eb89343cb9ddce09206fcc",
            "excerpt": "bench test --check data-handling-derivation\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 3 s (at 2aa2bf49 after bench worktree build FT342-FT253)\ncheck[1]{name,kind,tests_run,subjects}:\n  data-handling-derivation,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,10,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-data-handling",
          "command": "bench test --check data-handling-derivation",
          "exit_code": 0
        },
        {
          "id": "t1-system-lease-v2",
          "performer": "claude:blw_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1rep1/v-t1-system-lease.txt",
            "digest": "sha256:b02a0a5e0f9e3356d2cb55e947d982809502f5c6e9732e1bcc9d13ee2c234fc8",
            "excerpt": "bench test --check system --run TestTwoLandingProcessesShareOneLease\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 6 s (at 2aa2bf49 after bench worktree build FT342-FT253)\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,1321,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nalso (S2 touched this file): bench test --check system --run TestWorktreeLandPublicRaceAndRerun; exit 0\n  github.com/gibbonmi/bench/internal/systemtest,pass,1871,1\n"
          },
          "requirement": "t1-system-lease",
          "command": "bench test --check system --run TestTwoLandingProcessesShareOneLease",
          "exit_code": 0
        },
        {
          "id": "t1-root-conformance-v2",
          "performer": "claude:blw_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t1rep1/v-t1-root-conformance.txt",
            "digest": "sha256:43eff0f801c23243d4f568132b64645fe07dee4bbeb4c740e72fbfb6e70f3a3d",
            "excerpt": "bench test --check root-conformance\nrun through: bench worktree exec \"FT342-FT253\" -- <command>; exit 0, elapsed about 25 s (at 2aa2bf49 after bench worktree build FT342-FT253)\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,19304,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "bwa-standards-r1",
          "performer": "claude:blw_bwa_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:blw_bwa_standards",
            "digest": "sha256:2a579305eb5624c0889539019bb431135e837a1faca11a659be7a252747ae400",
            "excerpt": "Axis: Standards — 3 findings (worst: S1)\nS1 | internal/worktree/land_lease_test.go:47 | fixture harness duplicated from land_identity_test.go:241-264 | auto-fix (Writes expansion) | blocking | 9\nS2 | internal/systemtest/owner_land_race_test.go:182 | systemReleaseGateBarrier duplicates releaseArtifactBarrier | auto-fix | blocking | 8\nS3 | six new comment sentences over 25 words (STE) | auto-fix, evidence-only | blocking | 8\nFlag (b) refuted: Composed(FactLeaseOwner) follows faces_reset.go:12 precedent.\nAdvice: Composed doc wording; HoldExclusive/lockForClean overlap; leaseLanding parameter clump.\n"
          },
          "axis": "Standards",
          "base": "dcbbfb805419eef05174ea379b99547cfc8889da",
          "tip": "ced1dedc83e74057e0af3499d3b3ccdca9930319",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "bwa-spec-r1",
          "performer": "claude:blw_bwa_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:blw_bwa_spec",
            "digest": "sha256:c324857b586d09742b97eefc88f66bf0aa483182614d465d42e0fa4071785b8a",
            "excerpt": "Axis: Spec — BW-A @ ced1dedc (evidence sha256:694e08cf…)\nFindings: 1 (non-blocking)\nP1 spec.md:606 \"opened without following, as the gate lock opens its file\" — the lease\n   uses O_NOFOLLOW (gocache/lock.go:108); the gate lock does not (gate/run_transaction.go:70).\n   Non-behavioral spec contradiction; drop the comparison clause; flag for veto. conf 9.\nRows: BW1-BW16 built and graded.\nAuthor flag (b): PublishAdmitted nil after ref move matches spec:219; no finding.\nAdvice:\n- Test the ledger-write-failure branch at publication.go:47.\n- Pin the land-infrastructure face in TestLeaseLockUnavailable.\n- Gap uniqueness keyed on outcome only (matches spec:183).\n"
          },
          "axis": "Spec",
          "base": "dcbbfb805419eef05174ea379b99547cfc8889da",
          "tip": "ced1dedc83e74057e0af3499d3b3ccdca9930319",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "bwa-coverage-r1",
          "performer": "claude:blw_bwa_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9c5287c4e471d17997ad54c4176662f29030292",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:blw_bwa_coverage",
            "digest": "sha256:d57e0a830cbff268ae9cbd5a48c0f1c57b222caa88091fbcb19891516caff7ec",
            "excerpt": "Coverage BW-A (tip ced1dedc): 3 findings, 1 blocking.\nC1 BW7 test non-discriminating: lease drop moved to a second Transact inside\n   PublishAdmitted -> probe silent (7 tests), restored. The 10ms O_EXCL poll lets\n   the immediate second transaction win. Blocking.\nC2 Post-publish ledger-write failure: guard `published && err == nil` -> probe\n   silent (109 tests), restored. New failure mode with no row. Non-blocking.\nC3 Symlink at lease lock: O_NOFOLLOW removed -> silent in gocache and worktree\n   lease tests, restored. Edge inventory claims it. Non-blocking.\nAdvice: destination read precedes the lease; the kernel-contended path is proven\nonly by system BW10.\n"
          },
          "axis": "Coverage",
          "base": "dcbbfb805419eef05174ea379b99547cfc8889da",
          "tip": "ced1dedc83e74057e0af3499d3b3ccdca9930319",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
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
      "from": "sha256:e5fa3d7f7092f75a9f8bef63b4fd638d15330f29195691d7945b1fecd38042a3",
      "to": "sha256:9c7408fd49da2cbcb5c1e08c68770ed29e667560eecc2524aff5518d31fcce3a",
      "chunk_ids": {
        "BW-A": [
          "BW-A"
        ]
      }
    }
  ]
}
```
