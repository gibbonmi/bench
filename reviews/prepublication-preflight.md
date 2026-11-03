# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/prepublication-preflight/spec.md",
  "plan_digest": "sha256:eb3e51dca0f37f18117f06ce2c8ec17848d0ac921cc1b5598a9d43b0830ca3c6",
  "implementation_session": "",
  "chunks": [
    {
      "id": "C1",
      "base": "3d0598ba1b589d392eb8d11aa0c658e27e809e4c",
      "tip": "816c402f3e370e4cadf906470e5b5a62b9d72349",
      "plan_digest": "sha256:439822b5cbe0fa3d7522f3a589cf67fd2fd93e440fb9f9beed6a74d913859390",
      "source_digest": "bcd174d44e66961aa744235d1d943c2208ab4fc1",
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
        },
        {
          "id": "c1-01-r1-candidate-tests",
          "performer": "ppf-repair-01-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bcd174d44e66961aa744235d1d943c2208ab4fc1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ppf-repair-01-r1",
            "digest": "sha256:802cf6107e6fa7d8d7c4bf3a7bc29c98f0f7e5a92e8352fac1020da365a2b1d0",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/plan.go,swap,failed,1,TestBuildVerdictAtGradesTheCandidate/PP02_completion_plan_at_the_candidate,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestBuildVerdictAtGradesTheCandidate$,passed,5\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,397\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,327,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestBuildVerdictAtGradesTheCandidate/PP02_completion_plan_at_the_candidate,\"candidate_test.go:37: plan candidate = (0), want 1 with completion-plan red at the candidate:\",4\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "816c402f3e370e4cadf906470e5b5a62b9d72349",
          "requirement": "candidate-tests",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0,
          "probe": {
            "mutation": "Read the completion plan at HEAD in the candidate grade and record the PP02 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ppf-repair-01-r1",
              "digest": "sha256:802cf6107e6fa7d8d7c4bf3a7bc29c98f0f7e5a92e8352fac1020da365a2b1d0",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/plan.go,swap,failed,1,TestBuildVerdictAtGradesTheCandidate/PP02_completion_plan_at_the_candidate,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestBuildVerdictAtGradesTheCandidate$,passed,5\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,397\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,327,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestBuildVerdictAtGradesTheCandidate/PP02_completion_plan_at_the_candidate,\"candidate_test.go:37: plan candidate = (0), want 1 with completion-plan red at the candidate:\",4\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c1-standards-r1",
          "performer": "ppf-review-c1-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "52ea596345ba18cf4db25e522c2d2a8acf93aa5f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ppf-review-c1-standards",
            "digest": "sha256:fa6e6b385a0715ea3c24305a886266f5d012c5d880bce11964643a041de14b85",
            "excerpt": "Standards axis, chunk C1 (base 3d0598ba, tip 411b0276). Findings: 4. Worst: F1.\nF1 internal/git/status.go:63-65: WorktreeDirtyAt holds two dirty rules chosen by the spelling \"HEAD\"; porcelain counts a staged-then-reverted change dirty, the tree rule counts it clean; BuildVerdictAt passes a resolved SHA, so PP04 byte equality depends on index state. Violates ticket 01 \"dirty rule keeps one source\". medium, ask-user, confidence 7.\nF2 internal/git/status.go:71-101 repeats the TreeHash steps of internal/git/tree.go:18-41 with a different seed. low, auto-fix (widen tree.go), confidence 6.\nF3 \"HEAD\" literal rule twice: candidate.go:123 headRef and status.go:64. low, auto-fix, goes away with F1, confidence 5.\nF4 spec.md seam cells rewritten outside ticket Writes. low, no-op if convention holds, confidence 4.\nImplement command contributed to F1: the ticket did not define dirty against a non-HEAD commit.\n"
          },
          "axis": "Standards",
          "base": "3d0598ba1b589d392eb8d11aa0c658e27e809e4c",
          "tip": "411b02763e57c9f901b47697ab5271e03b4d45fd",
          "finding_ids": [
            "C1-STD-F1",
            "C1-STD-F2",
            "C1-STD-F3",
            "C1-STD-F4"
          ],
          "supersedes": []
        },
        {
          "id": "c1-spec-r1",
          "performer": "ppf-review-c1-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "52ea596345ba18cf4db25e522c2d2a8acf93aa5f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ppf-review-c1-spec",
            "digest": "sha256:f4791a2f9586b8e269504afd0a46787b47e4b8b743b8379fc91c3c6c4a778c3a",
            "excerpt": "Spec axis, chunk C1 (base 3d0598ba, tip 411b0276). All 12 C1 rows audited; every listed candidate read substitutes; one renderer; commitment row kept; public --source-tip keeps HEAD. Findings: 2. Worst: S1.\nS1 internal/preflight/candidate.go:141 with command.go:163-164: on snapshot drift the candidate grade prints a rerun of bench preflight build --source-tip \u003ccandidate>, which reds tip-current before publication; the right retry is the caller's commit call. low, ask-user or route to C3, confidence 6.\nS2 internal/git/status.go:78-102 WorktreeDirtyAt divergence meets spec.md:107; differs only for a staged change reverted in the worktree. no-op, confidence 6.\nImplement command: no contribution.\n"
          },
          "axis": "Spec",
          "base": "3d0598ba1b589d392eb8d11aa0c658e27e809e4c",
          "tip": "411b02763e57c9f901b47697ab5271e03b4d45fd",
          "finding_ids": [
            "C1-SPEC-S1",
            "C1-SPEC-S2"
          ],
          "supersedes": []
        },
        {
          "id": "c1-coverage-r1",
          "performer": "ppf-review-c1-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "52ea596345ba18cf4db25e522c2d2a8acf93aa5f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ppf-review-c1-coverage",
            "digest": "sha256:8a2ba1961f9677a8cccb2ed8ebdf991ca2e3b4d6b27ae680d42b220c7dfa3c21",
            "excerpt": "Coverage axis, chunk C1 (base 3d0598ba, tip 411b0276). Findings: 2. Worst: C1.\nC1 internal/preflight/gather.go:115 withFoldResolution(root, head, ...) -> WorktreeDirtyAt (gather_inputs.go:131): no fixture has a fold entry; a regression to headRef stays silent. medium, auto-fix, confidence 8. Probe: silent,internal/preflight/gather.go,swap,passed,0,null,yes (swap withFoldResolution(root, head, -> withFoldResolution(root, headRef, ; ./internal/preflight; 398 tests; baseline passed).\nC2 internal/preflight/gather.go:131 buildDeletedPaths: no candidate fixture has a ticket Writes naming a deleted path, so PP04's why-clause is not met for build deletions. low, ask-user, confidence 6.\nAuthor gaps: WorktreeDirtyAt error returns no-op; plan mutant vs PP04 no-op.\nImplement command contributed: spec quantified-obligation line lists only the varied reads.\n"
          },
          "axis": "Coverage",
          "base": "3d0598ba1b589d392eb8d11aa0c658e27e809e4c",
          "tip": "411b02763e57c9f901b47697ab5271e03b4d45fd",
          "finding_ids": [
            "C1-COV-F1",
            "C1-COV-F2"
          ],
          "supersedes": []
        },
        {
          "id": "c1-standards-r2",
          "performer": "ppf-confirm-c1-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bcd174d44e66961aa744235d1d943c2208ab4fc1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ppf-confirm-c1-standards",
            "digest": "sha256:dd7ebe90df8055f8cc30e17136d4c3ff5be5bb3fcee37ae483037855d12d37b7",
            "excerpt": "Standards confirming round, chunk C1 (repair delta 7dfe3ebf..816c402f, tip 4f526e63). C1-STD-F2 closed: one checkoutTree helper (internal/git/tree.go:30) serves TreeHash (tree.go:18) and WorktreeDirtyAt (status.go:71); seed order HEAD then --empty kept; \"none\" kept; indexStep keeps the git %s: %w: %s error. New findings: none. Not verified: no tests run.\n"
          },
          "axis": "Standards",
          "base": "3d0598ba1b589d392eb8d11aa0c658e27e809e4c",
          "tip": "816c402f3e370e4cadf906470e5b5a62b9d72349",
          "finding_ids": [],
          "supersedes": [
            "c1-standards-r1"
          ]
        },
        {
          "id": "c1-spec-r2",
          "performer": "ppf-confirm-c1-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bcd174d44e66961aa744235d1d943c2208ab4fc1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ppf-confirm-c1-spec",
            "digest": "sha256:356cb0283dbd2ec53abf23877fb849a2ff6fbf38bc17234e2ae5296b7626a6eb",
            "excerpt": "Spec confirming round, chunk C1 (tip 4f526e63). C1-SPEC-S1 closed: BuildVerdictAt passes nil args (candidate.go:25); verdictCommand skips snapshotDriftRefusal on nil args (command.go:161-167); public path unchanged (command.go:99). C1-STD-F1 closed: the spec sentence after spec.md:107 matches status.go:64-75. New advisory (low, non-blocking): \"differ only for a staged change that the checkout reverts\" overclaims; assume-unchanged, skip-worktree, and tracked-but-ignored files also differ.\n"
          },
          "axis": "Spec",
          "base": "3d0598ba1b589d392eb8d11aa0c658e27e809e4c",
          "tip": "816c402f3e370e4cadf906470e5b5a62b9d72349",
          "finding_ids": [],
          "supersedes": [
            "c1-spec-r1"
          ]
        },
        {
          "id": "c1-coverage-r2",
          "performer": "ppf-confirm-c1-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bcd174d44e66961aa744235d1d943c2208ab4fc1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ppf-confirm-c1-coverage",
            "digest": "sha256:ce2da1569a240cc55e8d571f3c8ed2a8064eacbb6e3afb14cd7c70efd2fa22b9",
            "excerpt": "Coverage confirming round, chunk C1 (tip 4f526e63). C1-COV-F1 closed (gather.go:115; candidate_test.go:209; seedFoldCandidate candidate_fixture_test.go:145). C1-COV-F2 closed (gather.go:131; candidate_test.go:210; seedBuildDeletionCandidate candidate_fixture_test.go:157). Probe set over ./internal/preflight, exit 0, baseline passed:\n1,bit,swap,failed,1,TestBuildVerdictAtMatchesThePublishedVerdict/PP04_fold_resolution,8,1,0,null,yes\n2,bit,swap,failed,1,TestBuildVerdictAtMatchesThePublishedVerdict/PP04_build_deletion,8,1,0,null,yes\nNew gaps: none.\n"
          },
          "axis": "Coverage",
          "base": "3d0598ba1b589d392eb8d11aa0c658e27e809e4c",
          "tip": "816c402f3e370e4cadf906470e5b5a62b9d72349",
          "finding_ids": [],
          "supersedes": [
            "c1-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "C2",
      "base": "816c402f3e370e4cadf906470e5b5a62b9d72349",
      "tip": "7827cb0bec9a148bd39a3a254dca7052dce2b332",
      "plan_digest": "sha256:eb3e51dca0f37f18117f06ce2c8ec17848d0ac921cc1b5598a9d43b0830ca3c6",
      "source_digest": "61408cb094d9dcb5e936c0e3ea1d33ef3f938581",
      "acceptance_rows": [
        "PP10",
        "PP11",
        "PP12",
        "PP13",
        "PP14",
        "PP15",
        "PP16"
      ],
      "verification": [
        {
          "id": "c2-02-check-seam",
          "performer": "ppf-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "61408cb094d9dcb5e936c0e3ea1d33ef3f938581",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ppf-author-02",
            "digest": "sha256:f143aa77912fd3340f09b8fd4c0844ced286c82e9f5d0d5d91659bfae0e52eb0",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/landing/landing.go,swap,failed,1,TestLandRunsTheCheckBeforePublication,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRunsTheCheckBeforePublication$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,416\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,fail,28,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRunsTheCheckBeforePublication,\"check_test.go:39: destination inside the check = 24e370f26c78a4f3118999ce965cbc07cf0149eb, want the expected base 44f3953a32ecff70b28d871779d02a569b6cf367\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "7827cb0bec9a148bd39a3a254dca7052dce2b332",
          "requirement": "check-seam",
          "command": "bench test --package ./internal/landing --package ./internal/commit --package ./internal/refusalroute",
          "exit_code": 0,
          "probe": {
            "mutation": "Call the check after the ref update and record the PP10 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ppf-author-02",
              "digest": "sha256:f143aa77912fd3340f09b8fd4c0844ced286c82e9f5d0d5d91659bfae0e52eb0",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/landing/landing.go,swap,failed,1,TestLandRunsTheCheckBeforePublication,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRunsTheCheckBeforePublication$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,416\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,fail,28,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRunsTheCheckBeforePublication,\"check_test.go:39: destination inside the check = 24e370f26c78a4f3118999ce965cbc07cf0149eb, want the expected base 44f3953a32ecff70b28d871779d02a569b6cf367\",1\nskips[0]{package,test,reason}:\n"
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
      "from": "sha256:351505691273d6505a9a5fcf36349186e657baf3e6769643711555ead6885e4d",
      "to": "sha256:439822b5cbe0fa3d7522f3a589cf67fd2fd93e440fb9f9beed6a74d913859390",
      "chunk_ids": {
        "C1": [
          "C1"
        ]
      }
    },
    {
      "from": "sha256:439822b5cbe0fa3d7522f3a589cf67fd2fd93e440fb9f9beed6a74d913859390",
      "to": "sha256:eb3e51dca0f37f18117f06ce2c8ec17848d0ac921cc1b5598a9d43b0830ca3c6",
      "chunk_ids": {
        "C1": [
          "C1"
        ]
      }
    }
  ]
}
```

## C1 round 1 pickup

The frozen pair is 3d0598ba1..411b02763. The axes found 8 raw findings. After the orchestrator collapsed the findings that name the same fix, 4 repair targets remain, and all of them are on ticket 01.

## Standards

Count: 4. Worst issue: C1-STD-F1.

- C1-STD-F1 (confidence 7): `internal/git/status.go:63-65` `WorktreeDirtyAt` uses the porcelain rule at `HEAD` and the tree rule at another commit. Disposition: no-op on the code. The spec states the tree test before publication (`spec.md:107`), and the Spec axis graded it as met. The two rules differ only for a staged change that the checkout reverts. The orchestrator records this divergence in the spec prose for reviewer veto.
- C1-STD-F2 (confidence 6): `internal/git/status.go:71-101` repeats the checkout-tree steps of `TreeHash` in `internal/git/tree.go:18-41`. Disposition: auto-fix. Share one checkout-tree helper with a seed argument. The plan commit widens `internal/git/tree.go`.
- C1-STD-F3 (confidence 5): the `"HEAD"` literal appears in `internal/preflight/candidate.go:123` and `internal/git/status.go:64`. Disposition: no-op. The literal is the Git name of the live head in two packages, not a duplicated rule.
- C1-STD-F4 (confidence 4): the ticket commit rewrote the spec's seam cells. Disposition: no-op. `/bench-implement-spec` requires the ticket commit to rewrite each covered `planned` seam cell.

## Spec

Count: 2. Worst issue: C1-SPEC-S1.

- C1-SPEC-S1 (confidence 6): `internal/preflight/candidate.go:141` with `internal/preflight/command.go:163-164`. On a snapshot drift, the candidate grade prints a rerun of `bench preflight build --source-tip <candidate>`, and that command reds `tip-current` before publication. Disposition: auto-fix. The candidate grade does not print the public preflight rerun on a drift.
- C1-SPEC-S2 (confidence 6): the `WorktreeDirtyAt` divergence meets `spec.md:107`. Disposition: no-op.

## Coverage

Count: 2. Worst issue: C1-COV-F1.

- C1-COV-F1 (confidence 8): `internal/preflight/gather.go:115` passes the candidate to `withFoldResolution`, and no fixture holds a fold entry. A swap to `headRef` probed `silent`. Disposition: auto-fix. Add a PP04 fixture whose review record folds a changed path.
- C1-COV-F2 (confidence 6): `internal/preflight/gather.go:131` `buildDeletedPaths` has no candidate fixture where a ticket `Writes:` line names a deleted path. Disposition: auto-fix under the run's advance approval. Add that PP04 fixture.
