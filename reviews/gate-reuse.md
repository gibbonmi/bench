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

## GR-C1 author witness reds

The ticket 01 author ran the new tests against the unchanged production code at 21b5918db. Six tests failed: TestGradeIdentityLeavesOutTheCheckoutPath (GR1, GR2), TestGradeIdentityFollowsTheLoaderRead (GR5), TestGradeIdentityGivesADefectOneValue (GR6), TestReaderPlanAdmitsATableOfASiblingCheckout (GR9), TestGradeIdentityFollowsTheRunBinarySource (GR54), and TestGradeIdentityGivesAnUndeclaredSourceOneValue (GR55).

The author then probed `internal/gate/subject.go` at the ticket commit, and each mutant below made its expected test fail:

- An absent schedule framed like an empty one failed TestGradeIdentityCarriesTheBaselineSchedule (GR3).
- Dropped schedule bytes failed TestGradeIdentityCarriesTheBaselineSchedule (GR4).
- `grade = rel` swapped with `_, grade = rel, path` failed TestGradeIdentityKeepsTheOrdinaryAndProspectiveSubjectsEqual (GR45).
- A dropped running-executable frame failed TestOracleIdentityAndOracleShareTheClosure (GR8).
- A run-binary member replaced with `resolutionName(sourceKind)` failed TestOracleIdentityAndOracleShareTheClosure (GR8).

Two first cuts of these mutants did not compile, and the author cut them again with the same intent. The planned posture probe is in the fenced record as `gr-t01-posture-probe`.
