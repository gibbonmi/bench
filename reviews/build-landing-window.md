# Review outcomes

BW-A round 1 reviewed the pair `dcbbfb80..ced1dedc`.
Repair 1 at `2aa2bf49` closed all seven findings, and the confirming round passed on all three axes.
BW-A used one of its two repair cycles.

BW-B round 1 reviewed the pair `2aa2bf49..6c90a3c6`.
Repair 1 at `4057ca20` closed all five findings, and the confirming round passed on all three axes.
BW-B used one of its two repair cycles.

BW-C round 1 reviewed the pair `4057ca20..db1bdc5d`, and repair 1 closed S1, S2, and C1-C4 at `6568c4c2`.
The BW-C confirming round found S3 and P1, so BW-C repair 2 follows.
A fable/high consultation chose the owner for BW-C S1 and the content binding for BW-C C1.

## Standards

Findings: 3. Worst issue: S1, a copied gate-barrier fixture harness.

- S1, auto-fix, confidence 9: `barrierLandingFixture` in `internal/worktree/land_lease_test.go:47` copies the FIFO gate barrier of `TestCommitmentGateRace` (`internal/worktree/land_identity_test.go:241-264`). Fold that test onto the one harness.
- S2, auto-fix, confidence 8: `systemReleaseGateBarrier` in `internal/systemtest/owner_land_race_test.go:182` repeats the body of `releaseArtifactBarrier` (`internal/systemtest/owner_artifact_recovery_test.go:147`). Call the existing helper.
- S3, auto-fix, confidence 8: six new comment sentences are over 25 words. They are in `lease.go:17`, `lease_test.go:26`, `build_records_test.go:18` and `:53`, `faces_land.go:3`, and `land_lease_test.go:104`.

BW-B findings: 3. Worst issue: BW-B S1, a hand-spelled shell quote in a test expectation.

- BW-B S1, auto-fix, confidence 7: `foldOfMain` in `internal/worktree/land_moved_destination_test.go:18` hand-spells the quote of `main`. Derive it through `sanitize.ShellQuote`, as `merge_route_test.go:188` does.
- BW-B S2, auto-fix, confidence 6: `internal/worktree/land_refusal.go:291` strips `refs/heads/` that `land.go:194` added. Pass the branch that the caller already holds.
- BW-B S3, auto-fix, confidence 5: `destinationUnchanged` at `internal/landing/gitexec.go:145` drops a `ResolveCommit` error without a comment. State that the publish compare-and-swap still refuses.

BW-C findings: 2. Worst issue: BW-C S1, a default-branch ancestry check with two owners that disagree on a failed query.

- BW-C S1, auto-fix, confidence 7: `FoldParents` (`internal/reviewrecord/fold.go:29-30`) repeats the ancestry check of `internal/worktree/merge.go:354-372`. Move it to one `internal/git` owner that refuses on a failed query. Ticket 3 repairs it.
- BW-C S2, auto-fix, confidence 6: `MergedTestRequirement` (`internal/reviewrecord/check.go:91-92`) hand-spells the package command that `internal/preflight/check_floor.go:197-203` owns. Give the command one owner. Ticket 4 repairs it.
- BW-C S3, auto-fix, confidence 8: `foldCoversBase` (`internal/reviewrecord/fold.go:115`) and `coverage.go:46` and `:85` still call the ancestry query through `benchgit.OK`. Route them through `git.IsAncestor`. Ticket 3 repairs it in repair 2.

## Spec

Findings: 1. Worst issue: P1, a false comparison clause in the Edge inventory.

- P1, auto-fix, confidence 9: the Edge inventory says that the lease lock opens without following a link "as the gate lock opens its file". The gate lock uses a plain open, so remove that clause.

BW-B findings: 0. The Spec axis accepted the three author flags. BW23 uses an injected typed refusal at the `j.landReviewed` hook, because no real fixture reaches that state.

BW-C findings: 0. The transitive second-fold reading reconciles the spec lines on the covered digest and the fold `From`.

- BW-C P1, auto-fix, confidence 9: the BW-C chunk row omitted row BW98. The orchestrator fixed it in a plan commit before repair 2.

## Coverage

Findings: 3. Worst issue: C1, a BW7 test that does not detect a second transaction.

- C1, auto-fix, confidence 8: the BW7 test passes when `PublishAdmitted` drops the lease in a second transaction. Make the competing acquire take the intent lock first, and show it red under that mutation.
- C2, auto-fix, confidence 7: no test makes the ledger write fail after the ref update. Add a test that asserts nil and a lease record that the next acquire takes over.
- C3, auto-fix, confidence 7: no test puts a symbolic link at the lease lock path. Add a refusal case beside BW15 that expects `landing lease unavailable:`.

BW-B findings: 2. Worst issue: BW-B C1, a widened cause match that no test detects.

- BW-B C1, auto-fix, confidence 7: a swap of the cause match at `internal/worktree/land_refusal.go:284` to any non-empty cause stays green. Add a row that pins a non-proof gate cause to `land-infrastructure`.
- BW-B C2, auto-fix, confidence 5: no test reaches the branch where the ancestry query fails. Add a row that pins its handback.

BW-C findings: 4. Worst issue: BW-C C1, a fold resolution path that a later commit can change past the fence.

- BW-C C1, auto-fix, confidence 6: `unauthorizedPaths` skips a resolution path by name alone. Bind the skip to the bytes and mode at the fold commit, and add row BW98. Ticket 4 repairs it.
- BW-C C2, auto-fix, confidence 7: no test gives a chunk base with the fold digest that does not descend from the fold commit. Ticket 3 adds it.
- BW-C C3, auto-fix, confidence 7: no test records a fold whose first parent has only a comment-only gap. Ticket 3 adds it.
- BW-C C4, auto-fix, confidence 7: no test records two folds in one gap. Ticket 3 adds it.

```bench-review-record
{
  "version": 2,
  "spec": "specs/build-landing-window/spec.md",
  "plan_digest": "sha256:b02696cfc82f5e064c45f5a220cc8788f75b5d8615aec9daf12beb59d9b6815e",
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
        },
        {
          "id": "bwa-standards-r2",
          "performer": "claude:blw_bwa_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:blw_bwa_standards_r2",
            "digest": "sha256:221a9d4bace14f3d5d8feb20d5852a53ddd7741d3c6c83ce6a533ade13e82733",
            "excerpt": "Axis: Standards (confirming) — 0 findings\nS1 closed | land_identity_test.go:238 uses barrierLandingFixture; barrier only at land_lease_test.go:47\nS2 closed | systemReleaseGateBarrier removed; releaseArtifactBarrier sole source (owner_artifact_recovery_test.go:147)\nS3 closed | all delta comment sentences <=22 words\nC1 helper: second copy of intent lock taking (intent.go:249); not blocking per publication_lock_test.go precedent\nAdvice: an intent test helper that holds the lock would give both tests one source\n"
          },
          "axis": "Standards",
          "base": "dcbbfb805419eef05174ea379b99547cfc8889da",
          "tip": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
          "finding_ids": [],
          "supersedes": [
            "bwa-standards-r1"
          ]
        },
        {
          "id": "bwa-spec-r2",
          "performer": "claude:blw_bwa_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:blw_bwa_spec_r2",
            "digest": "sha256:069c7397a0c25a98ed90ac32ec2538501c00ca594e50eb334fc87021501082b4",
            "excerpt": "Spec confirming round, BW-A, delta ced1dedc..2aa2bf49.\nP1 closed: spec.md:606 now reads \"...is opened without following.\" with the gate-lock comparison removed.\nProduction delta is comment-only (lease.go:17-18, faces_land.go:3-4); Landing lease and Records sections unchanged.\nTestPublishKeepsTheLeaseWhenTheLedgerWriteFails grades spec.md:219: nil publish, stale record kept, takeover clears it.\nTestLeaseLockUnavailable dangling-link and live-link rows grade spec.md:606 via the O_NOFOLLOW refusal.\nNew findings: 0. Blocking: none.\n"
          },
          "axis": "Spec",
          "base": "dcbbfb805419eef05174ea379b99547cfc8889da",
          "tip": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
          "finding_ids": [],
          "supersedes": [
            "bwa-spec-r1"
          ]
        },
        {
          "id": "bwa-coverage-r2",
          "performer": "claude:blw_bwa_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "5ef9843142d8ce9ad3693414c3b91b6c32e33141",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:blw_bwa_coverage_r2",
            "digest": "sha256:abe570796f62b2ef721716ee3dea9abd2633fcfdc312c2a6eb9049552ea61fc7",
            "excerpt": "Coverage confirming round, BW-A, delta ced1dedc..2aa2bf49: 0 findings; C1, C2, C3 closed.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/publication.go,swap,failed,1,yes\n  TestPublishReleasesTheLeaseInItsDecision,\"lease_test.go:118: at the end of the publish decision: leases [{Destination:refs/heads/main Assignment:5e1cf428…}], lease lock lock is held; want no lease and a free lock\"\n  bit,internal/commitment/repository/publication.go,swap,failed,1,yes\n  TestPublishKeepsTheLeaseWhenTheLedgerWriteFails,\"lease_test.go:140: publication whose ledger write fails = replace intent ledger: rename …\"\n  bit,internal/gocache/lock.go,swap,failed,2,yes\n  TestLeaseLockUnavailable/dangling-link + /live-link,\"land_lease_test.go:228: … refused{detail=composition must not start …\"\nAdvice: the C1 helper spins hot until bounds.TestDeadline.\n"
          },
          "axis": "Coverage",
          "base": "dcbbfb805419eef05174ea379b99547cfc8889da",
          "tip": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
          "finding_ids": [],
          "supersedes": [
            "bwa-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "BW-B",
      "base": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
      "tip": "4057ca206c8859f1f6f9634345426bebc8d77954",
      "plan_digest": "sha256:440c1b0b55e357d00d6e66cc8c8a42b2ceb210a19ec0c64bf37e7427ad156d72",
      "source_digest": "ca6bd4b8092b1010171b96b3b14eecbe92d48e51",
      "acceptance_rows": [
        "BW17",
        "BW18",
        "BW19",
        "BW20",
        "BW21",
        "BW22",
        "BW23",
        "BW24"
      ],
      "verification": [
        {
          "id": "t2-land-focused-v1",
          "performer": "claude:blw_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0ffb66e4290a68a222c8d14c98cec94fcae5da90",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2/t2-land-focused.txt",
            "digest": "sha256:4616340cebaf4fd2b1dab3531f386ac451fd2a5f3a72795392effc68275ee070",
            "excerpt": "bench test --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'\n# through bench worktree exec \"FT342-FT253\" at 6c90a3c648d6d283b0ede6d7dd3dee0bcbf24347 (clean); exit 0; elapsed 3.6 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1679,20\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-land-focused",
          "command": "bench test --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the ancestry condition that maps a completion proof cause to land-destination-moved, so that a completion proof cause over a moved destination routes to land-handback. TestCompletionOverAMovedDestinationNamesTheFold turns red.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2/t2-land-focused.txt",
              "digest": "sha256:4616340cebaf4fd2b1dab3531f386ac451fd2a5f3a72795392effc68275ee070",
              "excerpt": "bench test --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'\n# through bench worktree exec \"FT342-FT253\" at 6c90a3c648d6d283b0ede6d7dd3dee0bcbf24347 (clean); exit 0; elapsed 3.6 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1679,20\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2-authorization-focused-v1",
          "performer": "claude:blw_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0ffb66e4290a68a222c8d14c98cec94fcae5da90",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2/t2-authorization-focused.txt",
            "digest": "sha256:a1fb93800970ce12d83a4923c0f276c58e31c4067e3d22e250f8b854eb100438",
            "excerpt": "bench test --package ./internal/gate/authorization --run 'TestCompletionProofCauseTravelsAsAType'\n# through bench worktree exec \"FT342-FT253\" at 6c90a3c648d6d283b0ede6d7dd3dee0bcbf24347 (clean); exit 0; elapsed 2.0 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,66,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-authorization-focused",
          "command": "bench test --package ./internal/gate/authorization --run 'TestCompletionProofCauseTravelsAsAType'",
          "exit_code": 0
        },
        {
          "id": "t2-floor-v1",
          "performer": "claude:blw_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0ffb66e4290a68a222c8d14c98cec94fcae5da90",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2/t2-floor.txt",
            "digest": "sha256:a68618894c41f0d6ab1babf069ac42e54d0b34350b519a102a3cae0fdec69d75",
            "excerpt": "bench test --package ./internal/gate/... --package ./internal/landing/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/commit --package ./cmd/bench --package ./internal/conformance\n# through bench worktree exec \"FT342-FT253\" at 6c90a3c648d6d283b0ede6d7dd3dee0bcbf24347 (clean); exit 0; elapsed 108.9 s\npackages[12]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26554,333\n  github.com/gibbonmi/bench/internal/commit,pass,13666,117\n  github.com/gibbonmi/bench/internal/conformance,pass,80467,1690\n  github.com/gibbonmi/bench/internal/gate,pass,31389,430\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,269,14\n  github.com/gibbonmi/bench/internal/gate/greenmarker,pass,54,4\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,339,36\n  github.com/gibbonmi/bench/internal/landing,pass,10122,221\n  github.com/gibbonmi/bench/internal/landing/published,no-tests,0,0\n  github.com/gibbonmi/bench/internal/landing/settlepolicy,pass,3,20\n  github.com/gibbonmi/bench/internal/refusalroute,pass,5,17\n  github.com/gibbonmi/bench/internal/worktree,pass,76955,1415\nfailures[0]{package,test,line,lines}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-floor",
          "command": "bench test --package ./internal/gate/... --package ./internal/landing/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/commit --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t2-system-race-v1",
          "performer": "claude:blw_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0ffb66e4290a68a222c8d14c98cec94fcae5da90",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2/t2-system-race.txt",
            "digest": "sha256:1f8acdf583d258f788a43ea201812b1300c79d14a8956db49e8286ea592d434c",
            "excerpt": "bench test --check system --run TestWorktreeLandPublicRaceAndRerun\n# through bench worktree exec \"FT342-FT253\" at 6c90a3c648d6d283b0ede6d7dd3dee0bcbf24347 (clean, after bench worktree build FT342-FT253); exit 0; elapsed 4.1 s\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,1807,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-system-race",
          "command": "bench test --check system --run TestWorktreeLandPublicRaceAndRerun",
          "exit_code": 0
        },
        {
          "id": "t2-root-conformance-v1",
          "performer": "claude:blw_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0ffb66e4290a68a222c8d14c98cec94fcae5da90",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2/t2-root-conformance.txt",
            "digest": "sha256:7519c04d5cf7119afb30485692f151773c306f9bbc036aa0abb32690c6434e5f",
            "excerpt": "bench test --check root-conformance\n# through bench worktree exec \"FT342-FT253\" at 6c90a3c648d6d283b0ede6d7dd3dee0bcbf24347 (clean); exit 0; elapsed 23.1 s\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18288,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-land-focused-v2",
          "performer": "claude:blw_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca6bd4b8092b1010171b96b3b14eecbe92d48e51",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2rep1/t2-land-focused.txt",
            "digest": "sha256:8ed2b63652e45e77867c5f90c21fca0ddb4b1e6cd1145627dc89e86fc669f5d6",
            "excerpt": "bench test --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'\n# through bench worktree exec \"FT342-FT253\" at 4057ca206c8859f1f6f9634345426bebc8d77954 (clean); exit 0; elapsed 3 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1527,23\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n# plan probe:\nbench probe internal/worktree/land_refusal.go --swap 'ancestryErr != nil || holds' --with 'ancestryErr != nil || !holds' --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'\n# through bench worktree exec \"FT342-FT253\" at 4057ca206c8859f1f6f9634345426bebc8d77954 (clean); exit 0; elapsed 11 s; plan probe of t2-land-focused\ntree[1]{target,head,dirty}:\n  FT342-FT253,4057ca206c8859f1f6f9634345426bebc8d77954,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/land_refusal.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture,passed,23\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold,\"land_moved_destination_test.go:43: face land-destination-moved = (1, \\\"refused{detail=prospective authorization refused: infrastructure (gate subject unavailable: completion composition changes moved.txt; review the destination delta),next=r… (950 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-destination-moved,\"identity_component_test.go:238: face land-destination-moved = (1, \\\"refused{detail=prospective authorization refused: infrastructure (gate subject unavailable: completion composition changes moved.txt; review the destination delta),next=revi… (964 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/current-source,\"land_moved_destination_test.go:132: current-source = (1, \\\"refused{detail=landing destination moved past the reviewed source: prospective authorization refused: infrastructure (completion proof fault),next=bench worktree merge --from 'main' … (1020 bytes)\",2\nskips[0]{package,test,reason}:\n# probe exit 0 (reread without a wrapper), restored=yes\n# exit 0 confirmed by a second run at the same head (only reviews/build-landing-window.md dirty, from the chunk record); outcome pass, failures 0\n"
          },
          "requirement": "t2-land-focused",
          "command": "bench test --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the ancestry condition that maps a completion proof cause to land-destination-moved, so that a completion proof cause over a moved destination routes to land-handback. TestCompletionOverAMovedDestinationNamesTheFold turns red.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2rep1/t2-land-focused.txt",
              "digest": "sha256:8ed2b63652e45e77867c5f90c21fca0ddb4b1e6cd1145627dc89e86fc669f5d6",
              "excerpt": "bench test --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'\n# through bench worktree exec \"FT342-FT253\" at 4057ca206c8859f1f6f9634345426bebc8d77954 (clean); exit 0; elapsed 3 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1527,23\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n# plan probe:\nbench probe internal/worktree/land_refusal.go --swap 'ancestryErr != nil || holds' --with 'ancestryErr != nil || !holds' --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'\n# through bench worktree exec \"FT342-FT253\" at 4057ca206c8859f1f6f9634345426bebc8d77954 (clean); exit 0; elapsed 11 s; plan probe of t2-land-focused\ntree[1]{target,head,dirty}:\n  FT342-FT253,4057ca206c8859f1f6f9634345426bebc8d77954,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/land_refusal.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture,passed,23\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold,\"land_moved_destination_test.go:43: face land-destination-moved = (1, \\\"refused{detail=prospective authorization refused: infrastructure (gate subject unavailable: completion composition changes moved.txt; review the destination delta),next=r… (950 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-destination-moved,\"identity_component_test.go:238: face land-destination-moved = (1, \\\"refused{detail=prospective authorization refused: infrastructure (gate subject unavailable: completion composition changes moved.txt; review the destination delta),next=revi… (964 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/current-source,\"land_moved_destination_test.go:132: current-source = (1, \\\"refused{detail=landing destination moved past the reviewed source: prospective authorization refused: infrastructure (completion proof fault),next=bench worktree merge --from 'main' … (1020 bytes)\",2\nskips[0]{package,test,reason}:\n# probe exit 0 (reread without a wrapper), restored=yes\n# exit 0 confirmed by a second run at the same head (only reviews/build-landing-window.md dirty, from the chunk record); outcome pass, failures 0\n"
            }
          }
        },
        {
          "id": "t2-authorization-focused-v2",
          "performer": "claude:blw_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca6bd4b8092b1010171b96b3b14eecbe92d48e51",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2rep1/t2-authorization-focused.txt",
            "digest": "sha256:71f590a408d5a43c4c09956e52e6ea0fa91b8b91497578f27efe0d8a024463a4",
            "excerpt": "bench test --package ./internal/gate/authorization --run 'TestCompletionProofCauseTravelsAsAType'\n# through bench worktree exec \"FT342-FT253\" at 4057ca206c8859f1f6f9634345426bebc8d77954 (clean); exit 0; elapsed 1 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,67,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n# exit 0 confirmed by a second run at the same head (only reviews/build-landing-window.md dirty, from the chunk record); outcome pass, failures 0\n"
          },
          "requirement": "t2-authorization-focused",
          "command": "bench test --package ./internal/gate/authorization --run 'TestCompletionProofCauseTravelsAsAType'",
          "exit_code": 0
        },
        {
          "id": "t2-floor-v2",
          "performer": "claude:blw_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca6bd4b8092b1010171b96b3b14eecbe92d48e51",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2rep1/t2-floor.txt",
            "digest": "sha256:a2c57ea64d441bfe6dc02586b70b98f21fef76ce866b85776997eab563404082",
            "excerpt": "bench test --package ./internal/gate/... --package ./internal/landing/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/commit --package ./cmd/bench --package ./internal/conformance\n# through bench worktree exec \"FT342-FT253\" at 4057ca206c8859f1f6f9634345426bebc8d77954 (clean); exit 0; elapsed 104 s\npackages[12]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,24164,333\n  github.com/gibbonmi/bench/internal/commit,pass,12875,117\n  github.com/gibbonmi/bench/internal/conformance,pass,79818,1690\n  github.com/gibbonmi/bench/internal/gate,pass,31248,430\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,254,14\n  github.com/gibbonmi/bench/internal/gate/greenmarker,pass,50,4\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,320,36\n  github.com/gibbonmi/bench/internal/landing,pass,9849,221\n  github.com/gibbonmi/bench/internal/landing/published,no-tests,0,0\n  github.com/gibbonmi/bench/internal/landing/settlepolicy,pass,2,20\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,17\n  github.com/gibbonmi/bench/internal/worktree,pass,76068,1418\nfailures[0]{package,test,line,lines}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n# exit 0 confirmed by a second run at the same head (only reviews/build-landing-window.md dirty, from the chunk record); outcome pass, failures 0\n"
          },
          "requirement": "t2-floor",
          "command": "bench test --package ./internal/gate/... --package ./internal/landing/... --package ./internal/worktree --package ./internal/refusalroute --package ./internal/commit --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t2-system-race-v2",
          "performer": "claude:blw_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca6bd4b8092b1010171b96b3b14eecbe92d48e51",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2rep1/t2-system-race.txt",
            "digest": "sha256:f25f573c742cdc59fa86beddde6c2f4f61258c84f2aa9e57293a305e27c8f75e",
            "excerpt": "bench test --check system --run TestWorktreeLandPublicRaceAndRerun\n# through bench worktree exec \"FT342-FT253\" at 4057ca206c8859f1f6f9634345426bebc8d77954 (clean), after bench worktree build FT342-FT253; exit 0; elapsed 4 s\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,1811,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n# exit 0 confirmed by a second run at the same head (only reviews/build-landing-window.md dirty, from the chunk record); outcome pass, failures 0\n"
          },
          "requirement": "t2-system-race",
          "command": "bench test --check system --run TestWorktreeLandPublicRaceAndRerun",
          "exit_code": 0
        },
        {
          "id": "t2-root-conformance-v2",
          "performer": "claude:blw_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca6bd4b8092b1010171b96b3b14eecbe92d48e51",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2rep1/t2-root-conformance.txt",
            "digest": "sha256:d8f953cb773080c8c8dff734a7a0d3544b39a6bbc53778b26c51568c2011d053",
            "excerpt": "bench test --check root-conformance\n# through bench worktree exec \"FT342-FT253\" at 4057ca206c8859f1f6f9634345426bebc8d77954 (clean); exit 0; elapsed 20 s\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17914,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n# exit 0 confirmed by a second run at the same head (only reviews/build-landing-window.md dirty, from the chunk record); outcome pass, failures 0\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-land-focused-v3",
          "performer": "claude:blw_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca6bd4b8092b1010171b96b3b14eecbe92d48e51",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2rep1/t2-land-focused-v3.txt",
            "digest": "sha256:7e67b3345165f1d453b4b4bba7ad9018ecfb580af815be2c8699c9d8f76d2998",
            "excerpt": "bench probe internal/worktree/land_refusal.go --swap 'ancestryErr != nil || holds' --with 'ancestryErr != nil || !holds' --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'\n# through bench worktree exec \"FT342-FT253\" at 0b99f67d3f8e7daa65e4e4d6bc07dc799c8d58a8 (clean); bench probe exit 0; mutated focused run exit 1 (probe_execution.focused_exit_code); restored=yes\ntree[1]{target,head,dirty}:\n  FT342-FT253,0b99f67d3f8e7daa65e4e4d6bc07dc799c8d58a8,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/land_refusal.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture,passed,23\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1275\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,1656,23\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold,\"land_moved_destination_test.go:43: face land-destination-moved = (1, \\\"refused{detail=prospective authorization refused: infrastructure (gate subject unavailable: completion composition changes moved.txt; review the destination delta),next=r… (951 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-destination-moved,\"identity_component_test.go:238: face land-destination-moved = (1, \\\"refused{detail=prospective authorization refused: infrastructure (gate subject unavailable: completion composition changes moved.txt; review the destination delta),next=revi… (963 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/current-source,\"land_moved_destination_test.go:132: current-source = (1, \\\"refused{detail=landing destination moved past the reviewed source: prospective authorization refused: infrastructure (completion proof fault),next=bench worktree merge --from 'main' … (1020 bytes)\",2\nskips[0]{package,test,reason}:\n\nbench test --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'\n# through bench worktree exec \"FT342-FT253\" at 0b99f67d3f8e7daa65e4e4d6bc07dc799c8d58a8 (clean); exit 0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1700,23\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-land-focused",
          "command": "bench test --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the ancestry condition that maps a completion proof cause to land-destination-moved, so that a completion proof cause over a moved destination routes to land-handback. TestCompletionOverAMovedDestinationNamesTheFold turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t2rep1/t2-land-focused-v3.txt",
              "digest": "sha256:7e67b3345165f1d453b4b4bba7ad9018ecfb580af815be2c8699c9d8f76d2998",
              "excerpt": "bench probe internal/worktree/land_refusal.go --swap 'ancestryErr != nil || holds' --with 'ancestryErr != nil || !holds' --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'\n# through bench worktree exec \"FT342-FT253\" at 0b99f67d3f8e7daa65e4e4d6bc07dc799c8d58a8 (clean); bench probe exit 0; mutated focused run exit 1 (probe_execution.focused_exit_code); restored=yes\ntree[1]{target,head,dirty}:\n  FT342-FT253,0b99f67d3f8e7daa65e4e4d6bc07dc799c8d58a8,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/land_refusal.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture,passed,23\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1275\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,1656,23\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCompletionOverAMovedDestinationNamesTheFold,\"land_moved_destination_test.go:43: face land-destination-moved = (1, \\\"refused{detail=prospective authorization refused: infrastructure (gate subject unavailable: completion composition changes moved.txt; review the destination delta),next=r… (951 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestLandingRefusalRegistryHasAProducingFixture/land-destination-moved,\"identity_component_test.go:238: face land-destination-moved = (1, \\\"refused{detail=prospective authorization refused: infrastructure (gate subject unavailable: completion composition changes moved.txt; review the destination delta),next=revi… (963 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestProofFaultOverACurrentSourceHandsBack/current-source,\"land_moved_destination_test.go:132: current-source = (1, \\\"refused{detail=landing destination moved past the reviewed source: prospective authorization refused: infrastructure (completion proof fault),next=bench worktree merge --from 'main' … (1020 bytes)\",2\nskips[0]{package,test,reason}:\n\nbench test --package ./internal/worktree --run 'TestCompletionOverAMovedDestinationNamesTheFold|TestCompareAndSwapRefusalNamesTheFold|TestCheckoutChangeWithoutARefMoveKeepsTheHandback|TestProofFaultOverACurrentSourceHandsBack|TestLandingRefusalRegistryHasAProducingFixture'\n# through bench worktree exec \"FT342-FT253\" at 0b99f67d3f8e7daa65e4e4d6bc07dc799c8d58a8 (clean); exit 0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1700,23\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "bwb-standards-r1",
          "performer": "claude:blw_bwb_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0ffb66e4290a68a222c8d14c98cec94fcae5da90",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:blw_bwb_standards",
            "digest": "sha256:b69268ccadc2a368691a152564f3f160132e62c8da66056c056f13f1afdbd7bd",
            "excerpt": "Standards, BW-B, 2aa2bf49..6c90a3c6: 3 findings.\nS1 auto-fix blocking: land_moved_destination_test.go:18 foldOfMain hand-spells ShellQuote(\"main\"); precedent merge_route_test.go:188 derives it; no recorded red.\nS2 auto-fix judgment: land_refusal.go:291 strips refs/heads/ that land.go:194 added; caller holds branch.\nS3 auto-fix: gitexec.go:145 ignores ResolveCommit error without stating the CAS backstop.\nFlags: (a) pre-existing literal pattern, advice; (b) causeJoiner single-sources producer and parser, sound; (c) refuted, see S1.\nAdvice: export gate-handback from refusalroute; Face/Cause dual naming; \"moved past\" wording; rerun built twice in test; system test literal quoting.\n"
          },
          "axis": "Standards",
          "base": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
          "tip": "6c90a3c648d6d283b0ede6d7dd3dee0bcbf24347",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "bwb-spec-r1",
          "performer": "claude:blw_bwb_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0ffb66e4290a68a222c8d14c98cec94fcae5da90",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:blw_bwb_spec",
            "digest": "sha256:bf13b3e0f4211ccaeead5f57e0031a7041a59a5402a3328c5d3d287a55f5d9cb",
            "excerpt": "Axis: Spec — chunk BW-B, base 2aa2bf49, tip 6c90a3c6\nFindings: 0\nBW17-BW24: each row built at the cited seam; tests present, not run by this axis.\nThree refusal sites match spec :390-392; handbacks match :394-395 and R10.\nCause travels as authorization.Cause; no reason-text routing (land_refusal.go:284).\nFlag (a): sentence \": \" cause follows the tree convention at land_refusal.go:151-153.\nFlag (b): BW23 injects the cause at j.landReviewed; ancestry runs on real git.\nFlag (c): ancestry-error handback is unspecified; deferred to Coverage.\nAdvice: record BW23's injected-fixture decision in the review record.\n"
          },
          "axis": "Spec",
          "base": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
          "tip": "6c90a3c648d6d283b0ede6d7dd3dee0bcbf24347",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "bwb-coverage-r1",
          "performer": "claude:blw_bwb_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0ffb66e4290a68a222c8d14c98cec94fcae5da90",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:blw_bwb_coverage",
            "digest": "sha256:8659c930165af3b91b511df8713de14f129066b1425e695ed6591403f1e7eccd",
            "excerpt": "Axis: Coverage — 2 findings, 0 blocking.\nC1 | non-proof gate face at completion landing | widened Cause match routes it to handback/moved | no landing row pins non-proof Cause -> land-infrastructure | auto-fix | no | 7\nC2 | ancestry query fails | swapped branch routes to land-destination-moved | no row/test for the unreadable-ancestry branch | auto-fix | no | 5\nP1 land_refusal.go:284 swap Cause==CompletionProof -> Cause!=\"\" : silent (1415 tests, ./internal/worktree), restored\nP2 land_refusal.go:285 swap ancestryErr!=nil||holds -> ancestryErr==nil&&holds : silent (374 tests), restored\nAdvice: proof fault over a lagging source always folds; TipMovedError fires on ref rewind.\ngit status --short: empty.\n"
          },
          "axis": "Coverage",
          "base": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
          "tip": "6c90a3c648d6d283b0ede6d7dd3dee0bcbf24347",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "bwb-standards-r2",
          "performer": "claude:blw_bwb_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca6bd4b8092b1010171b96b3b14eecbe92d48e51",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:blw_bwb_standards_r2",
            "digest": "sha256:366d1a013d055d3813a41945cdff6c4f9c9f95a721031741a65a59f68f31c2ec",
            "excerpt": "Standards (confirming), BW-B, 6c90a3c6..4057ca20: 0 findings.\nS1 closed: land_moved_destination_test.go:19 derives via sanitize.ShellQuote(\"main\").\nS2 closed: land.go:215 passes branch; land_refusal.go TrimPrefix removed.\nS3 closed: gitexec.go comment states the publish CAS backstop.\nComments: all added sentences <=17 words, present tense.\nAdvice: non-proof-cause literal mirrors gate face name; test name vs infrastructure row.\n"
          },
          "axis": "Standards",
          "base": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
          "tip": "4057ca206c8859f1f6f9634345426bebc8d77954",
          "finding_ids": [],
          "supersedes": [
            "bwb-standards-r1"
          ]
        },
        {
          "id": "bwb-spec-r2",
          "performer": "claude:blw_bwb_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca6bd4b8092b1010171b96b3b14eecbe92d48e51",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:blw_bwb_spec_r2",
            "digest": "sha256:31766cef5f4ff91f25b107823af59012d5d76126d2c5c59d78e16962c1cfdb80",
            "excerpt": "Spec confirming round, BW-B, delta 6c90a3c6..4057ca20. Findings: 0.\nFace, three-step route, three refusal sites, two land-handback cases, and the typed cause are unchanged (spec :383-400).\nThe fold fact now comes from the caller's branch (land_refusal.go:292, land.go:215); same --from value.\nNon-proof row expects land-infrastructure: matches spec :392 and :23.\nUnreadable-ancestry row expects land-handback: the spec does not decide this case; no contradiction.\nThe axis ran its binding from the primary checkout; the orchestrator reran --check-current through the worktree: current=true at 38430278.\nAdvice: a spec line for the unreadable-ancestry handback would record the intent.\n"
          },
          "axis": "Spec",
          "base": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
          "tip": "4057ca206c8859f1f6f9634345426bebc8d77954",
          "finding_ids": [],
          "supersedes": [
            "bwb-spec-r1"
          ]
        },
        {
          "id": "bwb-coverage-r2",
          "performer": "claude:blw_bwb_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca6bd4b8092b1010171b96b3b14eecbe92d48e51",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:blw_bwb_coverage_r2",
            "digest": "sha256:da4a5e3d1592c4d7e0e257c252d46dc87d1a9f422d34feb816215151c1a7d05e",
            "excerpt": "Coverage confirming round, BW-B, delta 6c90a3c6..4057ca20: 0 findings; C1 and C2 closed.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/land_refusal.go,swap,failed,1,yes\n  TestProofFaultOverACurrentSourceHandsBack/non-proof-cause, land_moved_destination_test.go:132 (Cause != \"\" swap)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/land_refusal.go,swap,failed,1,yes\n  TestProofFaultOverACurrentSourceHandsBack/unreadable-ancestry, land_moved_destination_test.go:132 (ancestry branch swap)\nThe unreadable-ancestry row proves that IsAncestor returned an error; no flake source found.\nAdvice: put the object restore in t.Cleanup right after the hide.\n"
          },
          "axis": "Coverage",
          "base": "2aa2bf491f2b9009ba46acd15a48fe1301ded44f",
          "tip": "4057ca206c8859f1f6f9634345426bebc8d77954",
          "finding_ids": [],
          "supersedes": [
            "bwb-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "BW-C",
      "base": "4057ca206c8859f1f6f9634345426bebc8d77954",
      "tip": "901d91524635f47b5c646f02c990db0817c6f1b5",
      "plan_digest": "sha256:b02696cfc82f5e064c45f5a220cc8788f75b5d8615aec9daf12beb59d9b6815e",
      "source_digest": "141da2f0e03766e1dd8147aec1dd793a0d0c69dd",
      "acceptance_rows": [
        "BW25",
        "BW26",
        "BW27",
        "BW28",
        "BW29",
        "BW30",
        "BW31",
        "BW32",
        "BW33",
        "BW34",
        "BW35",
        "BW36",
        "BW37",
        "BW38",
        "BW41",
        "BW39",
        "BW40",
        "BW42",
        "BW43",
        "BW44",
        "BW45",
        "BW93",
        "BW98"
      ],
      "verification": [
        {
          "id": "t4-review-focused-v1",
          "performer": "claude:blw_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4/probe-bw42.log",
            "digest": "sha256:ff0ebed20f4f251a59607c831bc56dece49455cff8de9facb74bbe4e2ab95428",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/check.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestFoldResolutionNeedsReview|TestFoldDerivedFinalRequirements,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,385\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,309,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestFoldDerivedFinalRequirements,\"fold_test.go:107: completion without a merged-test result = <nil>, want the verification form for merged-test:./pkg/a\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-review-focused",
          "command": "bench test --package ./internal/reviewrecord --run 'TestFoldResolutionNeedsReview|TestFoldDerivedFinalRequirements'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the merged-test derived requirements that checkCompletion adds for each fold package, so that the complete checkpoint passes without a merged-test result. TestFoldDerivedFinalRequirements turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4/probe-bw42.log",
              "digest": "sha256:ff0ebed20f4f251a59607c831bc56dece49455cff8de9facb74bbe4e2ab95428",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/check.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestFoldResolutionNeedsReview|TestFoldDerivedFinalRequirements,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,385\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,309,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestFoldDerivedFinalRequirements,\"fold_test.go:107: completion without a merged-test result = <nil>, want the verification form for merged-test:./pkg/a\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t4-preflight-focused-v1",
          "performer": "claude:blw_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4/t4-preflight-focused.txt",
            "digest": "sha256:5cac59e4917efec85a6efa289ae2d2990bddfd43fe32f469f1baea6f9ed384eb",
            "excerpt": "bench worktree exec \"FT342-FT253\" -- bench test --package ./internal/preflight --run 'TestPathsAuthorizedAdmitsARecordedResolutionPath'\nexit 0; elapsed about 3 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,163,1\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-preflight-focused",
          "command": "bench test --package ./internal/preflight --run 'TestPathsAuthorizedAdmitsARecordedResolutionPath'",
          "exit_code": 0
        },
        {
          "id": "t4-route-focused-v1",
          "performer": "claude:blw_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4/t4-route-focused.txt",
            "digest": "sha256:c219a5f23314e0dc947972dd976daa7a6349ebefe57c074cedaa86ae286d5527",
            "excerpt": "bench worktree exec \"FT342-FT253\" -- bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'\nexit 0; elapsed about 5 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1307,22\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-route-focused",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'",
          "exit_code": 0
        },
        {
          "id": "t4-help-focused-v1",
          "performer": "claude:blw_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4/t4-help-focused.txt",
            "digest": "sha256:a9da0f5c8428947fbd99b15940ab0a8ad2491acf671fb830318d09ea21156493",
            "excerpt": "bench worktree exec \"FT342-FT253\" -- bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'\nexit 0; elapsed about 3 s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,4,1\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-help-focused",
          "command": "bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'",
          "exit_code": 0
        },
        {
          "id": "t4-floor-v1",
          "performer": "claude:blw_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4/t4-floor.txt",
            "digest": "sha256:dc747eadb2ad50c22bc117ae7dd988ad73449bf42fde69074fe72111b192aa4f",
            "excerpt": "bench worktree exec \"FT342-FT253\" -- bench test --package ./internal/reviewrecord/... --package ./internal/preflight/... --package ./internal/gate/... --package ./internal/worktree --package ./internal/refusalroute --package ./cmd/bench --package ./internal/conformance\nexit 0; elapsed about 122 s\npackages[15]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,24505,333\n  github.com/gibbonmi/bench/internal/conformance,pass,78424,1690\n  github.com/gibbonmi/bench/internal/gate,pass,27842,430\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,254,14\n  github.com/gibbonmi/bench/internal/gate/greenmarker,pass,46,4\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,302,36\n  github.com/gibbonmi/bench/internal/preflight,pass,44655,401\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,24890,229\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/refusalroute,pass,5,17\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,2563,89\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,10214,179\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/worktree,pass,75661,1418\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-floor",
          "command": "bench test --package ./internal/reviewrecord/... --package ./internal/preflight/... --package ./internal/gate/... --package ./internal/worktree --package ./internal/refusalroute --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t4-root-conformance-v1",
          "performer": "claude:blw_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4/t4-root-conformance.txt",
            "digest": "sha256:7a0ccedcaee8ba94ea94b89ceee6f83dc43ab565bc556306550be92f1be96502",
            "excerpt": "bench worktree exec \"FT342-FT253\" -- bench test --check root-conformance\nexit 0; elapsed about 24 s\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18632,1\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3-record-focused-v1",
          "performer": "claude:blw_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3/tip/t3-record-focused.txt",
            "digest": "sha256:e37dc140df81a8f7a91570af92329574fc6970f6aea138f3664f05b086abca78",
            "excerpt": "bench worktree exec \"FT342-FT253\" -- bench test --package ./internal/reviewrecord/recordcmd --run 'TestRecordFold'\nexit 0\ntree[1]{target,head,dirty}:\n  FT342-FT253,3dd6beda004daef6ed1d83b5e218e1f54b340585,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,886,11\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-record-focused",
          "command": "bench test --package ./internal/reviewrecord/recordcmd --run 'TestRecordFold'",
          "exit_code": 0
        },
        {
          "id": "t3-chain-focused-v1",
          "performer": "claude:blw_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3/tip/probe-plan.answer",
            "digest": "sha256:5034158a6a2f2f57bcfe4edc7530a20fd8e07a524fa9bdfaa42bcac04ff3a5d2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/coverage.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestFoldChainAcceptsARecordedFoldBetweenChunks|TestFoldCoversACompletionSourcePastTheLastChunk,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,537\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,273,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestFoldChainAcceptsARecordedFoldBetweenChunks,\"fold_test.go:41: chunk 2 on the record commit after a recorded fold: chunk 2: stale review chain gap: expected base 47bbb00e36f44f5fa384e01c8393576ae3866cea, the chunk 1 tip, got 921b7e9e4f866eb4ad055bfb0a900745ee3cbadb; plan commits land b… (457 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-chain-focused",
          "command": "bench test --package ./internal/reviewrecord --run 'TestFoldChainAcceptsARecordedFoldBetweenChunks|TestFoldCoversACompletionSourcePastTheLastChunk'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the fold-coverage branch for a later chunk base in checkSource, so that the record commit after a recorded fold reads as a chain gap. TestFoldChainAcceptsARecordedFoldBetweenChunks turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3/tip/probe-plan.answer",
              "digest": "sha256:5034158a6a2f2f57bcfe4edc7530a20fd8e07a524fa9bdfaa42bcac04ff3a5d2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/coverage.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestFoldChainAcceptsARecordedFoldBetweenChunks|TestFoldCoversACompletionSourcePastTheLastChunk,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,537\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,273,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestFoldChainAcceptsARecordedFoldBetweenChunks,\"fold_test.go:41: chunk 2 on the record commit after a recorded fold: chunk 2: stale review chain gap: expected base 47bbb00e36f44f5fa384e01c8393576ae3866cea, the chunk 1 tip, got 921b7e9e4f866eb4ad055bfb0a900745ee3cbadb; plan commits land b… (457 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3-gate-focused-v1",
          "performer": "claude:blw_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3/tip/t3-gate-focused.txt",
            "digest": "sha256:ef2ac87a1308ae9472e3a5eecd18541ca271dad30dafe159b5301056cfe26e79",
            "excerpt": "bench worktree exec \"FT342-FT253\" -- bench test --package ./internal/gate --run 'TestReviewCheckpointChainGapNamesTheExpectedBase'\nexit 0\ntree[1]{target,head,dirty}:\n  FT342-FT253,3dd6beda004daef6ed1d83b5e218e1f54b340585,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,91,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-gate-focused",
          "command": "bench test --package ./internal/gate --run 'TestReviewCheckpointChainGapNamesTheExpectedBase'",
          "exit_code": 0
        },
        {
          "id": "t3-floor-v1",
          "performer": "claude:blw_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3/tip/t3-floor.txt",
            "digest": "sha256:e32e47ef8e099dd17b93432d785d9f8b461cc0a88b0b103d546d8bd932c925f6",
            "excerpt": "bench worktree exec \"FT342-FT253\" -- bench test --package ./internal/reviewrecord/... --package ./internal/testreport --package ./internal/gate/... --package ./internal/landing/... --package ./internal/worktree --package ./cmd/bench --package ./internal/conformance\nexit 0\ntree[1]{target,head,dirty}:\n  FT342-FT253,3dd6beda004daef6ed1d83b5e218e1f54b340585,false\npackages[14]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,23938,333\n  github.com/gibbonmi/bench/internal/conformance,pass,77751,1690\n  github.com/gibbonmi/bench/internal/gate,pass,29443,430\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,230,14\n  github.com/gibbonmi/bench/internal/gate/greenmarker,pass,40,4\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,281,36\n  github.com/gibbonmi/bench/internal/landing,pass,9482,221\n  github.com/gibbonmi/bench/internal/landing/published,no-tests,0,0\n  github.com/gibbonmi/bench/internal/landing/settlepolicy,pass,2,20\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,2196,89\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,8874,179\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,89284,225\n  github.com/gibbonmi/bench/internal/worktree,pass,79901,1418\nfailures[0]{package,test,line,lines}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-floor",
          "command": "bench test --package ./internal/reviewrecord/... --package ./internal/testreport --package ./internal/gate/... --package ./internal/landing/... --package ./internal/worktree --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t3-root-conformance-v1",
          "performer": "claude:blw_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3/tip/t3-root-conformance.txt",
            "digest": "sha256:28ddc43311f0de2800a23cf7bf23f5ca485baf742de1280a9f9e031d6080a5e5",
            "excerpt": "bench worktree exec \"FT342-FT253\" -- bench test --check root-conformance\nexit 0\ntree[1]{target,head,dirty}:\n  FT342-FT253,3dd6beda004daef6ed1d83b5e218e1f54b340585,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16149,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t4-review-focused-v2",
          "performer": "claude:blw_t4_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4rep1/t4-review-focused.txt",
            "digest": "sha256:b47920c56b1b425e25a2cd794ec814dfd2602b63ea111dafbc44036f3267e818",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/check.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestFoldResolutionNeedsReview|TestFoldDerivedFinalRequirements,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,385\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,287,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestFoldDerivedFinalRequirements,\"fold_test.go:130: completion without a merged-test result = <nil>, want the verification form for merged-test:./pkg/a\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-review-focused",
          "command": "bench test --package ./internal/reviewrecord --run 'TestFoldResolutionNeedsReview|TestFoldDerivedFinalRequirements'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the merged-test derived requirements that checkCompletion adds for each fold package, so that the complete checkpoint passes without a merged-test result. TestFoldDerivedFinalRequirements turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4rep1/t4-review-focused.txt",
              "digest": "sha256:b47920c56b1b425e25a2cd794ec814dfd2602b63ea111dafbc44036f3267e818",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/check.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestFoldResolutionNeedsReview|TestFoldDerivedFinalRequirements,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,385\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,287,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestFoldDerivedFinalRequirements,\"fold_test.go:130: completion without a merged-test result = <nil>, want the verification form for merged-test:./pkg/a\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t4-preflight-focused-v2",
          "performer": "claude:blw_t4_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4rep1/t4-preflight-focused.txt",
            "digest": "sha256:742fff6f75a511fe2908e725f0cb6aa1be27f9d02515398a38d11caa24639a1d",
            "excerpt": "$ bench test --package ./internal/preflight --run 'TestPathsAuthorizedAdmitsARecordedResolutionPath'\ntree[1]{target,head,dirty}:\n  FT342-FT253,6568c4c20f4cde69168fac0ef0c52633d5831950,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,480,7\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit: 0\n"
          },
          "requirement": "t4-preflight-focused",
          "command": "bench test --package ./internal/preflight --run 'TestPathsAuthorizedAdmitsARecordedResolutionPath'",
          "exit_code": 0
        },
        {
          "id": "t4-route-focused-v2",
          "performer": "claude:blw_t4_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4rep1/t4-route-focused.txt",
            "digest": "sha256:cb209cd3b7d83fea84d8858c0d1b3e127a4e8e8d657562ec1f6e20cf32b0ad4a",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  FT342-FT253,6568c4c20f4cde69168fac0ef0c52633d5831950,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1306,22\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit: 0\n"
          },
          "requirement": "t4-route-focused",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'",
          "exit_code": 0
        },
        {
          "id": "t4-help-focused-v2",
          "performer": "claude:blw_t4_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4rep1/t4-help-focused.txt",
            "digest": "sha256:5fc47b1d4b960237c488090417d700e0f5441644c84ca3a813b5a6d9a1d46ad4",
            "excerpt": "$ bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'\ntree[1]{target,head,dirty}:\n  FT342-FT253,6568c4c20f4cde69168fac0ef0c52633d5831950,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,5,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit: 0\n"
          },
          "requirement": "t4-help-focused",
          "command": "bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'",
          "exit_code": 0
        },
        {
          "id": "t4-floor-v2",
          "performer": "claude:blw_t4_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4rep1/t4-floor.txt",
            "digest": "sha256:73422c34a5b057f81ad56118a9d1f3c21611a6602ad8374fc2912888cc2fab89",
            "excerpt": "$ bench test --package ./internal/reviewrecord/... --package ./internal/preflight/... --package ./internal/gate/... --package ./internal/worktree --package ./internal/refusalroute --package ./cmd/bench --package ./internal/conformance --package ./internal/git\ntree[1]{target,head,dirty}:\n  FT342-FT253,6568c4c20f4cde69168fac0ef0c52633d5831950,true\npackages[16]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,23358,333\n  github.com/gibbonmi/bench/internal/conformance,pass,73864,1690\n  github.com/gibbonmi/bench/internal/gate,pass,27588,430\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,243,14\n  github.com/gibbonmi/bench/internal/gate/greenmarker,pass,45,4\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,287,36\n  github.com/gibbonmi/bench/internal/git,pass,1737,172\n  github.com/gibbonmi/bench/internal/preflight,pass,43731,407\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,24302,229\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/refusalroute,pass,6,17\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,2755,91\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,9714,181\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/worktree,pass,72587,1418\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit: 0\n"
          },
          "requirement": "t4-floor",
          "command": "bench test --package ./internal/reviewrecord/... --package ./internal/preflight/... --package ./internal/gate/... --package ./internal/worktree --package ./internal/refusalroute --package ./cmd/bench --package ./internal/conformance --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "t4-root-conformance-v2",
          "performer": "claude:blw_t4_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t4rep1/t4-root-conformance.txt",
            "digest": "sha256:b1544149311e7ae02ce3fc7196f3562bd807b2576e1030e5a854b0b8ce830c65",
            "excerpt": "$ bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  FT342-FT253,6568c4c20f4cde69168fac0ef0c52633d5831950,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17178,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit: 0\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3-record-focused-v2",
          "performer": "claude:blw_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3rep1/tip/t3-record-focused.txt",
            "digest": "sha256:9c8148f1994947e0ed34eb8f4c98ecc35eeebf097cbd66f6d5193ff8293ffe02",
            "excerpt": "$ bench worktree exec \"FT342-FT253\" -- bench test --package ./internal/reviewrecord/recordcmd --run 'TestRecordFold'\nexit=0 wall=2757ms\ntree[1]{target,head,dirty}:\n  FT342-FT253,3cdff1788acc402bb7635c92e3907029f0fa8d72,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,1090,13\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-record-focused",
          "command": "bench test --package ./internal/reviewrecord/recordcmd --run 'TestRecordFold'",
          "exit_code": 0
        },
        {
          "id": "t3-chain-focused-v2",
          "performer": "claude:blw_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3rep1/tip/t3-chain-focused.txt",
            "digest": "sha256:e89156e455cd6839cab67ddb75dd721e9160e889a409db748cdb3d3b7879394e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/coverage.go,omit,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestFoldChainAcceptsARecordedFoldBetweenChunks|TestFoldCoversACompletionSourcePastTheLastChunk,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,910\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,534,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestFoldChainAcceptsARecordedFoldBetweenChunks/1_folds,\"fold_test.go:50: chunk 2 on the record commit after 1 recorded folds: chunk 2: stale review chain gap: expected base 49b38266871386554e56d5cb78bf04e710961144, the chunk 1 tip, got bca42d2ed97191bd3f50ce0492d117c917b6c2bc; plan commits land … (458 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestFoldChainAcceptsARecordedFoldBetweenChunks/2_folds,\"fold_test.go:50: chunk 2 on the record commit after 2 recorded folds: chunk 2: stale review chain gap: expected base 49b38266871386554e56d5cb78bf04e710961144, the chunk 1 tip, got f6b9623e63dac2457b7fc52abf4bb90bd2faf9f8; plan commits land … (458 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-chain-focused",
          "command": "bench test --package ./internal/reviewrecord --run 'TestFoldChainAcceptsARecordedFoldBetweenChunks|TestFoldCoversACompletionSourcePastTheLastChunk'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the fold-coverage branch for a later chunk base in checkSource, so that the record commit after a recorded fold reads as a chain gap. TestFoldChainAcceptsARecordedFoldBetweenChunks turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3rep1/tip/t3-chain-focused.txt",
              "digest": "sha256:e89156e455cd6839cab67ddb75dd721e9160e889a409db748cdb3d3b7879394e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/coverage.go,omit,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestFoldChainAcceptsARecordedFoldBetweenChunks|TestFoldCoversACompletionSourcePastTheLastChunk,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,910\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,534,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestFoldChainAcceptsARecordedFoldBetweenChunks/1_folds,\"fold_test.go:50: chunk 2 on the record commit after 1 recorded folds: chunk 2: stale review chain gap: expected base 49b38266871386554e56d5cb78bf04e710961144, the chunk 1 tip, got bca42d2ed97191bd3f50ce0492d117c917b6c2bc; plan commits land … (458 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestFoldChainAcceptsARecordedFoldBetweenChunks/2_folds,\"fold_test.go:50: chunk 2 on the record commit after 2 recorded folds: chunk 2: stale review chain gap: expected base 49b38266871386554e56d5cb78bf04e710961144, the chunk 1 tip, got f6b9623e63dac2457b7fc52abf4bb90bd2faf9f8; plan commits land … (458 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3-gate-focused-v2",
          "performer": "claude:blw_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3rep1/tip/t3-gate-focused.txt",
            "digest": "sha256:89ca5d6bf7c33a12c9301866026c930da98eb60bfec33e286c8eefae99a4470c",
            "excerpt": "$ bench worktree exec \"FT342-FT253\" -- bench test --package ./internal/gate --run 'TestReviewCheckpointChainGapNamesTheExpectedBase'\nexit=0 wall=1835ms\ntree[1]{target,head,dirty}:\n  FT342-FT253,3cdff1788acc402bb7635c92e3907029f0fa8d72,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,97,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-gate-focused",
          "command": "bench test --package ./internal/gate --run 'TestReviewCheckpointChainGapNamesTheExpectedBase'",
          "exit_code": 0
        },
        {
          "id": "t3-floor-v2",
          "performer": "claude:blw_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3rep1/tip/t3-floor.txt",
            "digest": "sha256:cfce18eeec7626fd164ed43dd2f8ea0e959a91bf9884338f77c5d9153821e9c6",
            "excerpt": "$ bench worktree exec \"FT342-FT253\" -- bench test --package ./internal/reviewrecord/... --package ./internal/testreport --package ./internal/gate/... --package ./internal/landing/... --package ./internal/worktree --package ./cmd/bench --package ./internal/conformance --package ./internal/git\nexit=0 wall=113417ms\ntree[1]{target,head,dirty}:\n  FT342-FT253,3cdff1788acc402bb7635c92e3907029f0fa8d72,false\npackages[15]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,25032,333\n  github.com/gibbonmi/bench/internal/conformance,pass,80469,1690\n  github.com/gibbonmi/bench/internal/gate,pass,30247,430\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,234,14\n  github.com/gibbonmi/bench/internal/gate/greenmarker,pass,43,4\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,277,36\n  github.com/gibbonmi/bench/internal/git,pass,1742,172\n  github.com/gibbonmi/bench/internal/landing,pass,9906,221\n  github.com/gibbonmi/bench/internal/landing/published,no-tests,0,0\n  github.com/gibbonmi/bench/internal/landing/settlepolicy,pass,2,20\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,2633,91\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,9238,181\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,92005,225\n  github.com/gibbonmi/bench/internal/worktree,pass,83882,1418\nfailures[0]{package,test,line,lines}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-floor",
          "command": "bench test --package ./internal/reviewrecord/... --package ./internal/testreport --package ./internal/gate/... --package ./internal/landing/... --package ./internal/worktree --package ./cmd/bench --package ./internal/conformance --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "t3-root-conformance-v2",
          "performer": "claude:blw_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "/tmp/claude-1000/-home-mgibs-workspace-bench/145e33b1-639f-44f5-b13c-f23a9a43ac0e/scratchpad/blw/t3rep1/tip/t3-root-conformance.txt",
            "digest": "sha256:527d3d074a487cb94bbb25dc492f3697e384144acf5675472711e00eb8680b20",
            "excerpt": "$ bench worktree exec \"FT342-FT253\" -- bench test --check root-conformance\nexit=0 wall=18659ms\ntree[1]{target,head,dirty}:\n  FT342-FT253,3cdff1788acc402bb7635c92e3907029f0fa8d72,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16621,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "bwc-standards-r1",
          "performer": "claude:blw_bwc_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:blw_bwc_standards",
            "digest": "sha256:1f8418eae8779f29c2562d6f07ec95c7454f91638e39956d68d69056de73457e",
            "excerpt": "Axis: Standards — 2 findings, worst S1.\nS1 fold.go:29-30 repeats the default-branch ancestry predicate of\n  merge.go:354-372 (one-source rule); the copies drift on query failure.\n  ask-user, blocking, conf 7.\nS2 check.go:91-92 hand-spells the floor package command that\n  check_floor.go:198-203 owns. ask-user, blocking, conf 6.\nReviewer direction met: record writes use write->replace; no t.Parallel.\nAdvice: swallowed git errors in follows/FoldParents; repairValues omits Fold;\n  FoldParents export unused outside package.\n"
          },
          "axis": "Standards",
          "base": "4057ca206c8859f1f6f9634345426bebc8d77954",
          "tip": "db1bdc5d31f26b60c029566175a428aca4563266",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "bwc-spec-r1",
          "performer": "claude:blw_bwc_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:blw_bwc_spec",
            "digest": "sha256:6980645db841a24cb7881956cd64004021b66e67df0aaf6c5533665cd8cf5a92",
            "excerpt": "Axis: Spec — chunk BW-C, pair 4057ca20..db1bdc5d — 0 findings, pass.\nBW25-BW45 and BW93 are built and graded by their cited tests.\nFold refusals match spec text; each refuses before the record write.\nChain gap and stale source repair through form `fold`; texts match.\nmerged-test:<pkg> derives once in finalRequirements; paths-authorized skips Resolution.\nmerge-fold-red has four steps and no FT342.\nAdvice: coveringFold reads \"From covers the predecessor's digest\" transitively (F2 untested);\nplaceholder <full-destination-commit> differs from the spec's <from-commit>; the Compose\nconflict-tree error path has no direct test.\n"
          },
          "axis": "Spec",
          "base": "4057ca206c8859f1f6f9634345426bebc8d77954",
          "tip": "db1bdc5d31f26b60c029566175a428aca4563266",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "bwc-coverage-r1",
          "performer": "claude:blw_bwc_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "476ab806de9a78412dbe22cdc072481e915ec1f0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:blw_bwc_coverage",
            "digest": "sha256:36628ff195a66358c76b003b971797b254a749bfd1a0e98aa51d304696f15f72",
            "excerpt": "Axis: Coverage — 4 findings, none blocking.\nC1 ask-user: paths_authorized.go:28 exempts a fold resolution path by name; a later build write to it passes the fence (DefaultTipPaths binds content).\nC2 auto-fix: foldCoversBase ancestry (fold.go:111) untested — probe silent.\nC3 auto-fix: comment-only first parent (fold.go:137, spec:334) untested — probe silent.\nC4 auto-fix: second fold in one gap (fold.go:92) untested — probe silent.\nProbes P1-P3 restored=yes; worktree clean.\nAdvice: resolver not excluded from fold review; backslash/newline path bytes untested.\n"
          },
          "axis": "Coverage",
          "base": "4057ca206c8859f1f6f9634345426bebc8d77954",
          "tip": "db1bdc5d31f26b60c029566175a428aca4563266",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": []
        },
        {
          "id": "bwc-standards-r2",
          "performer": "claude:blw_bwc_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:blw_bwc_standards_r2",
            "digest": "sha256:3c30f0df1564ebe792dfc4647ec1b477ff80783ab44d88b91fa5758dd35adc95",
            "excerpt": "Standards confirming round, BW-C, delta db1bdc5d..6568c4c2: 1 finding.\nS1 closed: git.IsAncestor (internal/git/git.go:219) sole owner;\n  DefaultContains (default_branch.go:39) used by fold.go:30, merge.go:359.\nS2 closed: reviewrecord.PackageTestCommand (check.go:93) used by\n  check.go:100 and check_floor.go:202.\nS3: fold.go:115 foldCoversBase, and coverage.go:46 and :85, call merge-base\n  --is-ancestor through benchgit.OK; a failed query reads as no. One-source rule.\n  The orchestrator grades it blocking under the bounded repair policy.\n"
          },
          "axis": "Standards",
          "base": "4057ca206c8859f1f6f9634345426bebc8d77954",
          "tip": "6568c4c20f4cde69168fac0ef0c52633d5831950",
          "finding_ids": [
            "S3"
          ],
          "supersedes": [
            "bwc-standards-r1"
          ]
        },
        {
          "id": "bwc-spec-r2",
          "performer": "claude:blw_bwc_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:blw_bwc_spec_r2",
            "digest": "sha256:dcd371ccfb7ffefb3a82c1390f3d502928e7d8c968d45e76066981e3521915dc",
            "excerpt": "Spec confirming round, BW-C, delta db1bdc5d..6568c4c2: 1 finding.\nP1 spec.md:427: the BW-C chunk row names BW25-BW45 and BW93, but plan commit\n  3032523a added BW98 to ticket 4 Covers. BENCH.md plan-expansion rule. Blocking.\nRows BW25-BW45 and BW93 keep their behavior; the fold refusal texts are unchanged.\nC1 implements the amended R8 and BW98, with the newest listing fold and the dirty guard.\nA failed default-branch query now refuses with its cause, before any write.\n"
          },
          "axis": "Spec",
          "base": "4057ca206c8859f1f6f9634345426bebc8d77954",
          "tip": "6568c4c20f4cde69168fac0ef0c52633d5831950",
          "finding_ids": [
            "P1"
          ],
          "supersedes": [
            "bwc-spec-r1"
          ]
        },
        {
          "id": "bwc-coverage-r2",
          "performer": "claude:blw_bwc_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b0928e81bb076546ec03ff2f77a55d2d7e3c6479",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:blw_bwc_coverage_r2",
            "digest": "sha256:3403c45fd471629b0b452897f160616b5b7e7f5a21aa649c17a64429090f8921",
            "excerpt": "Coverage confirming round, BW-C, delta db1bdc5d..6568c4c2: 0 findings; C1-C4 closed.\nC1 gather_inputs.go name-only swap: bit, 3 failed, restored yes.\nC2 foldCoversBase ancestry removed: bit, 2 failed, restored yes.\nC3 follows strict digest: bit, 1 failed, restored yes.\nC4 coveringFold stops advancing: bit, 1 failed, restored yes.\nS1 bypass DefaultContains true on failure, ./internal/... 9744 tests: bit, 1 failed, restored yes.\nDirty rule and newest-fold selection probes: bit.\n"
          },
          "axis": "Coverage",
          "base": "4057ca206c8859f1f6f9634345426bebc8d77954",
          "tip": "6568c4c20f4cde69168fac0ef0c52633d5831950",
          "finding_ids": [],
          "supersedes": [
            "bwc-coverage-r1"
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
      "from": "sha256:e5fa3d7f7092f75a9f8bef63b4fd638d15330f29195691d7945b1fecd38042a3",
      "to": "sha256:9c7408fd49da2cbcb5c1e08c68770ed29e667560eecc2524aff5518d31fcce3a",
      "chunk_ids": {
        "BW-A": [
          "BW-A"
        ]
      }
    },
    {
      "from": "sha256:9c7408fd49da2cbcb5c1e08c68770ed29e667560eecc2524aff5518d31fcce3a",
      "to": "sha256:52082aafa798935d8dba9b9521a8e470557b9bf38383be1e474ccf3b513f1007",
      "chunk_ids": {
        "BW-A": [
          "BW-A"
        ]
      }
    },
    {
      "from": "sha256:52082aafa798935d8dba9b9521a8e470557b9bf38383be1e474ccf3b513f1007",
      "to": "sha256:440c1b0b55e357d00d6e66cc8c8a42b2ceb210a19ec0c64bf37e7427ad156d72",
      "chunk_ids": {
        "BW-A": [
          "BW-A"
        ],
        "BW-B": [
          "BW-B"
        ]
      }
    },
    {
      "from": "sha256:440c1b0b55e357d00d6e66cc8c8a42b2ceb210a19ec0c64bf37e7427ad156d72",
      "to": "sha256:cbbe23b76e3faffd89b75fcdac95b2bd39c75c6bd6a427fc8cb3ade1f8ed32ce",
      "chunk_ids": {
        "BW-A": [
          "BW-A"
        ],
        "BW-B": [
          "BW-B"
        ]
      }
    },
    {
      "from": "sha256:cbbe23b76e3faffd89b75fcdac95b2bd39c75c6bd6a427fc8cb3ade1f8ed32ce",
      "to": "sha256:b48b0493be97dee45f23a112128b66238d7c62246a8f6d5c50f0d04f36db7781",
      "chunk_ids": {
        "BW-A": [
          "BW-A"
        ],
        "BW-B": [
          "BW-B"
        ],
        "BW-C": [
          "BW-C"
        ]
      }
    },
    {
      "from": "sha256:b48b0493be97dee45f23a112128b66238d7c62246a8f6d5c50f0d04f36db7781",
      "to": "sha256:b02696cfc82f5e064c45f5a220cc8788f75b5d8615aec9daf12beb59d9b6815e",
      "chunk_ids": {
        "BW-A": [
          "BW-A"
        ],
        "BW-B": [
          "BW-B"
        ],
        "BW-C": [
          "BW-C"
        ]
      }
    }
  ]
}
```
