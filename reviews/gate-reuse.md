# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/gate-reuse/spec.md",
  "plan_digest": "sha256:b0f6edda9075e930683fb916045172e340e0ada400ebacc98510788749872dce",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GR-C1",
      "base": "21b5918db91ee0badcfd4fc26678f3d180ae9e6a",
      "tip": "5600fe2cf34166e06e8206026fcedffa21462139",
      "plan_digest": "sha256:b0f6edda9075e930683fb916045172e340e0ada400ebacc98510788749872dce",
      "source_digest": "7be3b9eaf7234c5dd70e8c2ec68e6c4162a7ddd8",
      "acceptance_rows": [
        "GR1",
        "GR2",
        "GR3",
        "GR4",
        "GR5",
        "GR6",
        "GR7",
        "GR8",
        "GR9",
        "GR45",
        "GR54",
        "GR55"
      ],
      "verification": [
        {
          "id": "gr-t01-gate",
          "performer": "claude:gr-t01-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7be3b9eaf7234c5dd70e8c2ec68e6c4162a7ddd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t01-author-20261026",
            "digest": "sha256:9cd2daa636dd9277aa842eb8fd1620d31ca7698f024440fdec833b83376ef02c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,203539,677\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5600fe2cf34166e06e8206026fcedffa21462139",
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "gr-t01-authorization",
          "performer": "claude:gr-t01-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7be3b9eaf7234c5dd70e8c2ec68e6c4162a7ddd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t01-author-20261026",
            "digest": "sha256:59dfd0c38a847154e209b2cbf2ab8475d0d64bef3116c87e3d79360b77f3150d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,1668,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5600fe2cf34166e06e8206026fcedffa21462139",
          "requirement": "authorization",
          "command": "bench test --package ./internal/gate/authorization",
          "exit_code": 0
        },
        {
          "id": "gr-t01-root-conformance",
          "performer": "claude:gr-t01-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7be3b9eaf7234c5dd70e8c2ec68e6c4162a7ddd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t01-author-20261026",
            "digest": "sha256:fc12785052414e2bce2ba2c2cf3ebe2d31d400e5562cccfb121f088a1717cd8b",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,11949,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5600fe2cf34166e06e8206026fcedffa21462139",
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "gr-t01-posture-probe",
          "performer": "claude:gr-t01-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7be3b9eaf7234c5dd70e8c2ec68e6c4162a7ddd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t01-author-20261026",
            "digest": "sha256:655eeaddc40d2efe6ebc7a5867d6b073c94973ff0b50426375e680862fbcfc88",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,swap,failed,1,TestOracleIdentityCarriesTheCapabilityPosture,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestOracleIdentityCarriesTheCapabilityPosture,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,409\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,310,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityCarriesTheCapabilityPosture,\"oracle_identity_test.go:77: identity = 1f10bbc6c837344b29ca9d4cd6d702980744d24d7b26a56dfa59727516238eb4 under both postures, want two identities\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5600fe2cf34166e06e8206026fcedffa21462139",
          "requirement": "posture-probe",
          "command": "bench probe internal/gate/subject.go --swap 'strconv.FormatBool(strict())' --with 'strconv.FormatBool(false)' --package ./internal/gate --run TestOracleIdentityCarriesTheCapabilityPosture --expect TestOracleIdentityCarriesTheCapabilityPosture",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the posture from the grade identity. The two postures then give one identity, so the posture row fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gr-t01-author-20261026",
              "digest": "sha256:655eeaddc40d2efe6ebc7a5867d6b073c94973ff0b50426375e680862fbcfc88",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,swap,failed,1,TestOracleIdentityCarriesTheCapabilityPosture,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestOracleIdentityCarriesTheCapabilityPosture,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,409\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,310,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityCarriesTheCapabilityPosture,\"oracle_identity_test.go:77: identity = 1f10bbc6c837344b29ca9d4cd6d702980744d24d7b26a56dfa59727516238eb4 under both postures, want two identities\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gr-c1-standards-1",
          "performer": "claude:gr-c1-standards-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7be3b9eaf7234c5dd70e8c2ec68e6c4162a7ddd8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c1-standards-20261026",
            "digest": "sha256:38b10f6d0ed8677ed5bca2188f67bf1bd36b8dec82aa19f71c66e3030a25c3d3",
            "excerpt": "Standards axis, GR-C1, frozen pair 21b5918db..be7ab2aba, evidence sha256:2b9eee23 check-current current=true dirty=false. sonnet/high.\n\nS1 advisory auto-fix: internal/gate/subject.go:433 new relativeInside duplicates the inside-root predicate that confinedPath keeps inline at subject.go:321. confinedPath should call the helper.\nS2 advisory ask-user: internal/gate/subject.go:462 runningExecutableDigest returns the bare literal \"unreadable\"; unreadableBaselineRunner = \"unreadable\" exists at internal/gate/prospective.go:50. Two sentinels, same spelling, no shared name.\n\nClean: comment prose; oracle frame order kept; identityPath removed with no caller; frameBaselineSchedule reuses readManifest; one runBinarySource helper; gate.go 376 and run_transaction.go 370 lines; internal import only.\nRaw 2 (0 blocking). Repair targets 2. Worst: S1.\nDisposition (orchestrator): S1 accepted, repair. S2 no-op after tiebreak: unreadableBaselineRunner records a different fact; the executable member fails closed under Q1.\n"
          },
          "axis": "Standards",
          "base": "21b5918db91ee0badcfd4fc26678f3d180ae9e6a",
          "tip": "5600fe2cf34166e06e8206026fcedffa21462139",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "gr-c1-spec-1",
          "performer": "claude:gr-c1-spec-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7be3b9eaf7234c5dd70e8c2ec68e6c4162a7ddd8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c1-spec-20261026",
            "digest": "sha256:f6b5e80cbdc34dbddfc534f03df5591544960ee860bf7cadbc1aeb077df492de",
            "excerpt": "Spec axis, GR-C1, frozen pair 21b5918..be7ab2a, evidence check-current current=true dirty=false. sonnet/high.\n\nP1 advisory ask-user: internal/gate/subject.go:408 hashLinkChain frames the raw os.Readlink target into both digests, so an absolute in-checkout link puts the checkout path into the grade identity (story 2).\nP2 advisory ask-user: internal/gate/subject.go:461-463 and prospective.go:45-52 an unreadable running executable returns a constant, so two unreadable executables share one identity member (story 7).\n\nRows GR1-GR9, GR45, GR54, GR55 audited and closed. Plan re-cut 5600fe2cf keeps the posture probe intent; recorded probe bit. No scope creep.\nRaw 2 (0 blocking). Repair targets 2. Worst: P1.\nDisposition (orchestrator, after opus/sonnet consultation and fable tiebreak): P1 accepted, repair (raw target to the oracle, relative spelling to the identity through one helper shared with framePath, plus a link case). P2 accepted, repair (the running executable member fails closed through s.open).\n"
          },
          "axis": "Spec",
          "base": "21b5918db91ee0badcfd4fc26678f3d180ae9e6a",
          "tip": "5600fe2cf34166e06e8206026fcedffa21462139",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "gr-c1-coverage-1",
          "performer": "claude:gr-c1-coverage-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7be3b9eaf7234c5dd70e8c2ec68e6c4162a7ddd8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c1-coverage-20261026",
            "digest": "sha256:30afe366d6e6c74e7e9a03e8d5698e20d336e23d9de6629468ccde7ad7c80def",
            "excerpt": "Coverage axis, GR-C1, frozen pair 21b5918d..be7ab2ab, evidence check-current current=true. sonnet/high. bench test over the GR-C1 tests passed (20 tests).\n\nC1 advisory ask-user: subject.go:93-100,142 and gate.go:308-315 the fallback arm of runBinarySource (no self-built binary, runtime root differs from identity root) is never exercised.\nC2 advisory ask-user: subject.go:93-100 sourceKind follows the prospective script's presence, not the evaluation mode.\nC3 advisory ask-user: subject.go:408 an absolute in-checkout link target enters the grade identity raw.\nC4 advisory ask-user: subject.go:231-243 the runningExecutableDigest seam replaces the whole chain; the real chain has no test.\nC5 advisory ask-user: reviews/gate-reuse.md records one witness red only.\nC6 advisory ask-user: oracle_identity_test.go:120-123 policy leaves the grade identity; no row asks for it.\nRaw 6 (0 blocking). Worst: C1.\nDisposition (orchestrator): C1 accepted, repair (one test). C2 no change; bench learning flagged for reviewer veto (fable tiebreak). C3 merged with Spec P1. C4 no-op: GR58 in GR-C3 tests the real chain. C5 closed by record commit e7a51cc4f. C6 no change: decision #3 replaces the policy version (both consultants agree).\n"
          },
          "axis": "Coverage",
          "base": "21b5918db91ee0badcfd4fc26678f3d180ae9e6a",
          "tip": "5600fe2cf34166e06e8206026fcedffa21462139",
          "finding_ids": [
            "C1"
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
  }
}
```

## GR-C1 author witness reds

The ticket 01 author ran the new tests against the unchanged production code at 21b5918db. Six tests failed: TestGradeIdentityLeavesOutTheCheckoutPath (GR1, GR2), TestGradeIdentityFollowsTheLoaderRead (GR5), TestGradeIdentityGivesADefectOneValue (GR6), TestReaderPlanAdmitsATableOfASiblingCheckout (GR9), TestGradeIdentityFollowsTheRunBinarySource (GR54), and TestGradeIdentityGivesAnUndeclaredSourceOneValue (GR55).

The author then probed `internal/gate/subject.go` at the ticket commit, and each mutant below made its expected test fail:

- An absent schedule framed like an empty one failed TestGradeIdentityCarriesTheBaselineSchedule (GR3).
- Dropped schedule bytes failed TestGradeIdentityCarriesTheBaselineSchedule (GR4).
- `grade = rel` swapped with `_, grade = rel, path` failed TestGradeIdentityKeepsTheOrdinaryAndProspectiveSubjectsEqual (GR45).
- A dropped running-executable frame failed TestOracleIdentityAndOracleShareTheClosure (GR8).
- A run-binary member replaced with `resolutionName(sourceKind)` failed TestOracleIdentityAndOracleShareTheClosure (GR8).

Two first cuts of these mutants did not compile, and the author cut them again with the same intent. The planned posture probe is in the fenced record as `gr-t01-posture-probe`.
