# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/gate-reuse/spec.md",
  "plan_digest": "sha256:99a53aa7b6eaee91fb9b3ac47a50d59b64e1313ddd3eb348932b26da0b60be7c",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GR-C1",
      "base": "21b5918db91ee0badcfd4fc26678f3d180ae9e6a",
      "tip": "71fa9105aa396dfa355bf42397ce0d1ae7614e93",
      "plan_digest": "sha256:745dd33a081ad71b9ee0cbd09774d8de2bd19290864dcdce214a7796bf5ccb4d",
      "source_digest": "7b40c63dd81311f1bf4472c3892376209959c4b7",
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
        },
        {
          "id": "gr-t01-r1-gate",
          "performer": "claude:gr-t01-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7b40c63dd81311f1bf4472c3892376209959c4b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t01-repair1-20261026",
            "digest": "sha256:503ca3bc5a04572da7701a45cb8f9bfe785df376936578b977bd03bc2025ae20",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,268269,678\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "71fa9105aa396dfa355bf42397ce0d1ae7614e93",
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "gr-t01-r1-authorization",
          "performer": "claude:gr-t01-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7b40c63dd81311f1bf4472c3892376209959c4b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t01-repair1-20261026",
            "digest": "sha256:1ba70f9d5f990a2adda5906adfa04de81c8a56acb865aed6b855284bed9a3c99",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,2242,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "71fa9105aa396dfa355bf42397ce0d1ae7614e93",
          "requirement": "authorization",
          "command": "bench test --package ./internal/gate/authorization",
          "exit_code": 0
        },
        {
          "id": "gr-t01-r1-root-conformance",
          "performer": "claude:gr-t01-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7b40c63dd81311f1bf4472c3892376209959c4b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t01-repair1-20261026",
            "digest": "sha256:7897c36750190974480d032dae42552a1d5965070050e86ecf6549fded1c68f7",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17069,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "71fa9105aa396dfa355bf42397ce0d1ae7614e93",
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "gr-t01-r1-posture-probe",
          "performer": "claude:gr-t01-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7b40c63dd81311f1bf4472c3892376209959c4b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t01-repair1-20261026",
            "digest": "sha256:3dc3383c6f1c893fde53e00f857b3ce04d8501dfaabf278f74685e38c5c1b9c5",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,swap,failed,1,TestOracleIdentityCarriesTheCapabilityPosture,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestOracleIdentityCarriesTheCapabilityPosture,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,409\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,465,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityCarriesTheCapabilityPosture,\"oracle_identity_test.go:77: identity = 76576f19f866b136d7148f462f57f61e820686efbd024d390b9b83f17ff74819 under both postures, want two identities\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "71fa9105aa396dfa355bf42397ce0d1ae7614e93",
          "requirement": "posture-probe",
          "command": "bench probe internal/gate/subject.go --swap 'strconv.FormatBool(strict())' --with 'strconv.FormatBool(false)' --package ./internal/gate --run TestOracleIdentityCarriesTheCapabilityPosture --expect TestOracleIdentityCarriesTheCapabilityPosture",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the posture from the grade identity. The two postures then give one identity, so the posture row fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gr-t01-repair1-20261026",
              "digest": "sha256:3dc3383c6f1c893fde53e00f857b3ce04d8501dfaabf278f74685e38c5c1b9c5",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,swap,failed,1,TestOracleIdentityCarriesTheCapabilityPosture,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestOracleIdentityCarriesTheCapabilityPosture,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,409\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,465,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityCarriesTheCapabilityPosture,\"oracle_identity_test.go:77: identity = 76576f19f866b136d7148f462f57f61e820686efbd024d390b9b83f17ff74819 under both postures, want two identities\",1\nskips[0]{package,test,reason}:\n"
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
        },
        {
          "id": "gr-c1-standards-2",
          "performer": "claude:gr-c1-standards2-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7b40c63dd81311f1bf4472c3892376209959c4b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c1-standards2-20261026",
            "digest": "sha256:241c88996ee86014f529575d94fff80f3a9f56d4c8773dc2f030cfd0c9d676f4",
            "excerpt": "Standards confirming round, GR-C1, repair delta 8b1fd9cb1..71fa9105a, record 38f86f2f2. sonnet/high.\nS1 confirmed: confinedPath calls relativeInside (internal/gate/subject.go:324); the only inside-root predicate in subject.go is relativeInside (446-449).\ngradePath (subject.go:436-445): one derivation shared by framePath and hashLinkChain; no duplication.\nframeBaselineSchedule move confirmed unchanged (prospective.go:75-88); subject.go 512 of its 513-line budget.\nNew findings: none. Raw count 0.\n"
          },
          "axis": "Standards",
          "base": "21b5918db91ee0badcfd4fc26678f3d180ae9e6a",
          "tip": "71fa9105aa396dfa355bf42397ce0d1ae7614e93",
          "finding_ids": [],
          "supersedes": [
            "gr-c1-standards-1"
          ]
        },
        {
          "id": "gr-c1-spec-2",
          "performer": "claude:gr-c1-spec2-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7b40c63dd81311f1bf4472c3892376209959c4b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c1-spec2-20261026",
            "digest": "sha256:a9c3e29875993c2cf429e55eb6dac561126bb4bc7b8dc8583a800645b1ecd111",
            "excerpt": "Spec confirming round, GR-C1, repair delta 8b1fd9cb1..71fa9105a. sonnet/high.\nP1 confirmed: subject.go:411-420 frames the raw target to the oracle at the same position (oracle bytes unchanged, ticket 01 line 27) and the gradePath spelling to the grade identity for an absolute target.\nP2 confirmed: subject.go:143-148 opens the subject when the running executable digest is empty, after s.Oracle is set; unreadableBaselineRunner stays oracle-only.\nNo acceptance row changes behavior. New findings: none. Raw count 0.\n"
          },
          "axis": "Spec",
          "base": "21b5918db91ee0badcfd4fc26678f3d180ae9e6a",
          "tip": "71fa9105aa396dfa355bf42397ce0d1ae7614e93",
          "finding_ids": [],
          "supersedes": [
            "gr-c1-spec-1"
          ]
        },
        {
          "id": "gr-c1-coverage-2",
          "performer": "claude:gr-c1-coverage2-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7b40c63dd81311f1bf4472c3892376209959c4b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c1-coverage2-20261026",
            "digest": "sha256:076fb12d0da848bcf5a54488c4acf8d5df7370815048d5a39ff37064cec2782c",
            "excerpt": "Coverage confirming round, GR-C1, repair delta 8b1fd9cb1..71fa9105a. sonnet/high. bench test --run of the three repair tests passed at 38f86f2f2.\nC1 confirmed: grade_identity_test.go:111-119 fails an always-runtime-root or always-ProspectiveGateSh mutant.\nP1/C3 confirmed: linked case at grade_identity_test.go:20-44 fails under the old raw-target frame.\nP2 confirmed: oracle_identity_test.go:92-96 fails under the old closed subject.\nNew findings: K1 low auto-fix, the repair witness reds were not recorded; closed by the orchestrator's record prose \"GR-C1 repair witness reds\". K2 low no-op, the oracle raw-target half is not isolated; the Spec confirming axis confirmed unchanged oracle bytes. K3 low no-op, the production \"\" default is not unit-tested; GR58 in GR-C3 covers the real executable chain.\nRaw count 3 (0 blocking).\n"
          },
          "axis": "Coverage",
          "base": "21b5918db91ee0badcfd4fc26678f3d180ae9e6a",
          "tip": "71fa9105aa396dfa355bf42397ce0d1ae7614e93",
          "finding_ids": [],
          "supersedes": [
            "gr-c1-coverage-1"
          ]
        }
      ]
    },
    {
      "id": "GR-C2",
      "base": "a5f9c772c7b5fb4812ab5b42c1c1611e84bbec4f",
      "tip": "f40e49568deeaff93e72b607d8f94724ff1a5bc0",
      "plan_digest": "sha256:99a53aa7b6eaee91fb9b3ac47a50d59b64e1313ddd3eb348932b26da0b60be7c",
      "source_digest": "d0eb2c5d8d0ca8a2643cc6628bf8fda02e53638e",
      "acceptance_rows": [
        "GR12",
        "GR13",
        "GR10",
        "GR11",
        "GR14",
        "GR15",
        "GR16",
        "GR17",
        "GR18",
        "GR19",
        "GR20",
        "GR21",
        "GR65",
        "GR53",
        "GR57",
        "GR22",
        "GR47",
        "GR23",
        "GR24",
        "GR44"
      ],
      "verification": [
        {
          "id": "gr-t02--gate",
          "performer": "claude:gr-t02-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bfc08f5174612cf8526a98e7229d1df3e8d101ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-author-20261026",
            "digest": "sha256:60e0c899b57ffb8f1846874b00cfe46c18483b8d42c8d725ac09ee8e20bd13e4",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,206034,685\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c702eb07d1869efa03cd602ad305776e80744aa8",
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "gr-t02--authorization",
          "performer": "claude:gr-t02-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bfc08f5174612cf8526a98e7229d1df3e8d101ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-author-20261026",
            "digest": "sha256:49573b7279b5c59235b5880f93be6f720312d66482e206339c9c21b282edbf4b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,1604,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c702eb07d1869efa03cd602ad305776e80744aa8",
          "requirement": "authorization",
          "command": "bench test --package ./internal/gate/authorization",
          "exit_code": 0
        },
        {
          "id": "gr-t02--worktree",
          "performer": "claude:gr-t02-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bfc08f5174612cf8526a98e7229d1df3e8d101ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-author-20261026",
            "digest": "sha256:3acf2f5ace81248507d7489332395ab3bcdb50b44d16c280cbd6302dbebbf25b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,130317,1516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c702eb07d1869efa03cd602ad305776e80744aa8",
          "requirement": "worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "gr-t02--root-conformance",
          "performer": "claude:gr-t02-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bfc08f5174612cf8526a98e7229d1df3e8d101ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-author-20261026",
            "digest": "sha256:ea2ddbe6c26422e4a833dd7558f2fcfd397889aab1fe1becc258c40ff67eb711",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12817,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c702eb07d1869efa03cd602ad305776e80744aa8",
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "gr-t02--reuse-probe",
          "performer": "claude:gr-t02-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bfc08f5174612cf8526a98e7229d1df3e8d101ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-author-20261026",
            "digest": "sha256:4b532e01fbe33f970ba53fa5c86efa7f6810103a3916d9471806e36e14006a41",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/engine.go,swap,failed,1,TestGradeKeyCrossesCheckouts,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGradeKeyCrossesCheckouts,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,318\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,975,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGradeKeyCrossesCheckouts,\"grade_key_test.go:68: primary stdout = \\\"\\\", want the reuse line alone\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c702eb07d1869efa03cd602ad305776e80744aa8",
          "requirement": "reuse-probe",
          "command": "bench probe internal/gate/engine.go --swap 'inspection.ReusableGreen, inspection.Reason = true, \"\"' --with 'inspection.ReusableGreen, inspection.Reason = false, \"probe\"' --package ./internal/gate --run TestGradeKeyCrossesCheckouts --expect TestGradeKeyCrossesCheckouts",
          "exit_code": 0,
          "probe": {
            "mutation": "Make the store answer no green. The primary run then runs the gate again, so the cross-checkout row fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gr-t02-author-20261026",
              "digest": "sha256:4b532e01fbe33f970ba53fa5c86efa7f6810103a3916d9471806e36e14006a41",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/engine.go,swap,failed,1,TestGradeKeyCrossesCheckouts,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGradeKeyCrossesCheckouts,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,318\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,975,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGradeKeyCrossesCheckouts,\"grade_key_test.go:68: primary stdout = \\\"\\\", want the reuse line alone\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "gr-t02--commands",
          "performer": "claude:gr-t02-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bfc08f5174612cf8526a98e7229d1df3e8d101ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-author-20261026",
            "digest": "sha256:cf21d2944e21601cf3ae1f0abd7bc0027e519b059fbbc13f6eb2f224cbc22c87",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26985,406\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c702eb07d1869efa03cd602ad305776e80744aa8",
          "requirement": "commands",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "gr-t03-gate-03",
          "performer": "claude:gr-t03-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "54e7b686fd638aa7404aabce8412dbbeb84dda34",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t03-author-20261026",
            "digest": "sha256:b64c03df75b04257615a7a41b1909f14eda7510e009a3de9eba81c79e79a2506",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,210119,689\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d3b24e188662787ff8ab2f1bed5525834c89f0db",
          "requirement": "gate-03",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "gr-t03-root-conformance-03",
          "performer": "claude:gr-t03-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "54e7b686fd638aa7404aabce8412dbbeb84dda34",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t03-author-20261026",
            "digest": "sha256:aaad372a4e782bc8eed3eca3e9c0d78e17167211d21583263b38c5635276713a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13633,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d3b24e188662787ff8ab2f1bed5525834c89f0db",
          "requirement": "root-conformance-03",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "gr-t02-r1-gate",
          "performer": "claude:gr-t02-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "284e8acfb4e04c41dbad217431941138e992b782",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-repair1-20261026",
            "digest": "sha256:f558320f3f81b9b2a44674e1942536a111547431f616c59347af7875a050dc67",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,209695,697\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f1fe7d8bc736ef1d209f7a98df721c1edfe1a144",
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "gr-t02-r1-authorization",
          "performer": "claude:gr-t02-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "284e8acfb4e04c41dbad217431941138e992b782",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-repair1-20261026",
            "digest": "sha256:4ba30417de3cbe41537df65a7792a2d9a9c97c7123386d87484772033f82351f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,1554,18\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f1fe7d8bc736ef1d209f7a98df721c1edfe1a144",
          "requirement": "authorization",
          "command": "bench test --package ./internal/gate/authorization",
          "exit_code": 0
        },
        {
          "id": "gr-t02-r1-worktree",
          "performer": "claude:gr-t02-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "284e8acfb4e04c41dbad217431941138e992b782",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-repair1-20261026",
            "digest": "sha256:509b7fb53cc6ccca7d52ef8b0ace58778f0ba892ac6695aef49f3156d2e964e3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,129935,1516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f1fe7d8bc736ef1d209f7a98df721c1edfe1a144",
          "requirement": "worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "gr-t02-r1-root-conformance",
          "performer": "claude:gr-t02-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "284e8acfb4e04c41dbad217431941138e992b782",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-repair1-20261026",
            "digest": "sha256:0cc9f839c2e53d9b44f41e9687a73282c461332766a034bb129d44fe352507ba",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12485,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f1fe7d8bc736ef1d209f7a98df721c1edfe1a144",
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "gr-t02-r1-reuse-probe",
          "performer": "claude:gr-t02-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "284e8acfb4e04c41dbad217431941138e992b782",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-repair1-20261026",
            "digest": "sha256:091f8db70fb57fba5ff624079dd0b3bfb80e41b808daa13c08c7cea09fb8acaa",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/engine.go,swap,failed,1,TestGradeKeyCrossesCheckouts,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGradeKeyCrossesCheckouts,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,318\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,964,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGradeKeyCrossesCheckouts,\"grade_key_test.go:68: primary stdout = \\\"\\\", want the reuse line alone\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f1fe7d8bc736ef1d209f7a98df721c1edfe1a144",
          "requirement": "reuse-probe",
          "command": "bench probe internal/gate/engine.go --swap 'inspection.ReusableGreen, inspection.Reason = true, \"\"' --with 'inspection.ReusableGreen, inspection.Reason = false, \"probe\"' --package ./internal/gate --run TestGradeKeyCrossesCheckouts --expect TestGradeKeyCrossesCheckouts",
          "exit_code": 0,
          "probe": {
            "mutation": "Make the store answer no green. The primary run then runs the gate again, so the cross-checkout row fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gr-t02-repair1-20261026",
              "digest": "sha256:091f8db70fb57fba5ff624079dd0b3bfb80e41b808daa13c08c7cea09fb8acaa",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/engine.go,swap,failed,1,TestGradeKeyCrossesCheckouts,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGradeKeyCrossesCheckouts,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,318\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,964,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGradeKeyCrossesCheckouts,\"grade_key_test.go:68: primary stdout = \\\"\\\", want the reuse line alone\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "gr-t02-r1-commands",
          "performer": "claude:gr-t02-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "284e8acfb4e04c41dbad217431941138e992b782",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t02-repair1-20261026",
            "digest": "sha256:431bcfb5e27ae7a63ac8b1854377cf3d77aef78fdca99c01fe28e1aa0160fc61",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,25192,406\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f1fe7d8bc736ef1d209f7a98df721c1edfe1a144",
          "requirement": "commands",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "gr-t03-r1-gate-03",
          "performer": "claude:gr-t03-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d0eb2c5d8d0ca8a2643cc6628bf8fda02e53638e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t03-repair1-20261026",
            "digest": "sha256:7eda3da9e79d601202f727743d66c1de5a22c974732dd6801f3a6ab07f7cc42c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,208528,700\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f40e49568deeaff93e72b607d8f94724ff1a5bc0",
          "requirement": "gate-03",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "gr-t03-r1-root-conformance-03",
          "performer": "claude:gr-t03-repair1-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d0eb2c5d8d0ca8a2643cc6628bf8fda02e53638e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gr-t03-repair1-20261026",
            "digest": "sha256:63fd92107f3f6ef5f92add872b9c4dd67bb49854ed4028ae245ed4b2a4b02789",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13501,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f40e49568deeaff93e72b607d8f94724ff1a5bc0",
          "requirement": "root-conformance-03",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "gr-c2-standards-1",
          "performer": "claude:gr-c2-standards-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "54e7b686fd638aa7404aabce8412dbbeb84dda34",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c2-standards-20261026",
            "digest": "sha256:69ce759f3481950186ffc73388ce8d4140e0754da42362259fff2e147612771d",
            "excerpt": "Standards axis, GR-C2, pair a5f9c772c..1c7f5ce2f, check-current current=true. sonnet/high.\nS1 advisory: internal/gate/engine.go:149-169 loadEvidence restates the record-class knowledge and the shared schema/tree check that verdict.go:393 validateRecordBytes owns.\nS2 advisory auto-fix: internal/gate/run_durability.go:106 retainGreen hand-builds the ready record that readyVerdict (verdict.go:119-122) now owns.\nS3 advisory: run_durability.go:531-541 and 546-550 reuseWritesNothing and recordReusedGreen ask the same \"cache describes this green\" question.\nS4 advisory no-op: reuseLine literal in grade_key_test.go:25 and land_journey_test.go:82.\nRaw 4 (0 blocking).\nDisposition (orchestrator): S1 accepted, ticket 02 repair (one shared check helper; the evidence store stays its own record class). S2 and S3 accepted, ticket 03 repair. S4 no-op: an independent expectation that a constant cannot share across packages.\n"
          },
          "axis": "Standards",
          "base": "a5f9c772c7b5fb4812ab5b42c1c1611e84bbec4f",
          "tip": "d3b24e188662787ff8ab2f1bed5525834c89f0db",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "gr-c2-spec-1",
          "performer": "claude:gr-c2-spec-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "54e7b686fd638aa7404aabce8412dbbeb84dda34",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c2-spec-20261026",
            "digest": "sha256:b8cc9fd0f71902dcb310d205543e55092cba6860494fab632fd4e6118a8388bb",
            "excerpt": "Spec axis, GR-C2, pair a5f9c772c..1c7f5ce2f, check-current current=true. sonnet/high. Rows GR10-GR24, GR44, GR47, GR53, GR57, GR65 map to a change or a named test.\nP1 advisory: run_durability.go:531-541 GR44 covers a bare lock holder only; a live holder with a Pending cache still refuses.\nP2 advisory: no recorded witness red for GR19, GR21, GR53, GR57, GR65 (existing, unchanged tests).\nP3 advisory: specs/landing-test-efficiency/spec.md:167 and internal/gate/prospective_owner_test.go:3 still say \"baseline runner identity\".\nP4 advisory: the ticket 02 widenings cite --auto-approve.\nP5 advisory: a destination green written by another executable is not inherited; a red is labeled candidate-owned.\nItem (a): LTE40 and the journey test are not weakened.\nRaw 5 (0 blocking).\nDisposition (orchestrator): P3 accepted, ticket 02 repair (the LTE ticket 09 text stays as that ticket's own record). P1 no change, bench learning. P2 no-op: no new independent expectation. P4 no-op: bench spec widen carried --decision under the reviewer's --auto-approve. P5 accepted as a consequence of story 7, recorded in the same learning.\n"
          },
          "axis": "Spec",
          "base": "a5f9c772c7b5fb4812ab5b42c1c1611e84bbec4f",
          "tip": "d3b24e188662787ff8ab2f1bed5525834c89f0db",
          "finding_ids": [
            "P3"
          ],
          "supersedes": []
        },
        {
          "id": "gr-c2-coverage-1",
          "performer": "claude:gr-c2-coverage-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "54e7b686fd638aa7404aabce8412dbbeb84dda34",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c2-coverage-20261026",
            "digest": "sha256:79240df5102195c4643798ae36e3b59dc8d16b192b788adaa7f20cf85884fe1b",
            "excerpt": "Coverage axis, GR-C2, pair a5f9c772c..1c7f5ce2f, check-current current=true. sonnet/high. The GR-C2 tests and four neighbours passed under bench test --run.\nC1 advisory: internal/gate/engine.go:167-177 loadEvidence has no malformed-record test (missing grade, legacy oracle, duplicate key, wrong mode, schema or tree).\nC2 advisory: run_transaction.go:94-102 a lock taken after the reuseWritesNothing check still refuses.\nC3 advisory: run_durability.go:86-91 no test has an expired or drifted same-tree green cache beside fresh evidence.\nC4 advisory: authorization.go:228 no test pins Grade in evidenceToken.\nC5 advisory no-op: land_journey_test.go:80-84 no longer exercises Authorize Green on the identical tree.\nC6 advisory: no recorded witness for GR19, GR21, GR53, GR57, GR65.\nRaw 6 (0 blocking).\nDisposition (orchestrator): C1 and C4 accepted, ticket 02 repair. C3 accepted, ticket 03 repair. C2 no change, bench learning (merged with Spec P1). C5 no-op: the GR-C3 landing rows cover the identical-tree authorization. C6 no-op: existing, unchanged tests.\n"
          },
          "axis": "Coverage",
          "base": "a5f9c772c7b5fb4812ab5b42c1c1611e84bbec4f",
          "tip": "d3b24e188662787ff8ab2f1bed5525834c89f0db",
          "finding_ids": [
            "C1",
            "C3",
            "C4"
          ],
          "supersedes": []
        },
        {
          "id": "gr-c2-standards-2",
          "performer": "claude:gr-c2-standards2-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d0eb2c5d8d0ca8a2643cc6628bf8fda02e53638e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c2-standards2-20261026",
            "digest": "sha256:b7ac5f545e1f08b30177b5113db8ed4a28c8361d5fca99ccdd6185bf1f66fe9f",
            "excerpt": "Standards confirming round, GR-C2, repair deltas 25bd01294..f1fe7d8bc (ticket 02) and 20df42e77..f40e49568 (ticket 03). sonnet/high. Evidence check-current current=true at 9c9c014a5.\nS1 confirmed: validateStoreRecordCommon (internal/gate/verdict.go:409-414) is the one schema-and-tree check; validateRecordBytes (verdict.go:393) and loadEvidence (engine.go:175) call it.\nS2 confirmed: retainGreen (run_durability.go:112-113) builds from readyVerdict.\nS3 confirmed: cacheDescribesFreshGreen (run_durability.go:95-97) answers both the pre-lock and the under-lock question.\nNew findings: none. A pre-existing hard-coded schema outside the delta went to bench learning.\n"
          },
          "axis": "Standards",
          "base": "a5f9c772c7b5fb4812ab5b42c1c1611e84bbec4f",
          "tip": "f40e49568deeaff93e72b607d8f94724ff1a5bc0",
          "finding_ids": [],
          "supersedes": [
            "gr-c2-standards-1"
          ]
        },
        {
          "id": "gr-c2-spec-2",
          "performer": "claude:gr-c2-spec2-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d0eb2c5d8d0ca8a2643cc6628bf8fda02e53638e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c2-spec2-20261026",
            "digest": "sha256:f293b6757d50d5ad5e49db10c1871fdb479915c794cb3847167cd754959ecb8f",
            "excerpt": "Spec confirming round, GR-C2, repair deltas 25bd01294..f1fe7d8bc and 20df42e77..f40e49568. sonnet/high.\nP3 confirmed: specs/landing-test-efficiency/spec.md:167 reads \"exact tree, run binary member, and freshness window\" with the line's other claims kept; its tickets are untouched; internal/gate/prospective_owner_test.go:1-3 states the grade key.\nNo GR-C2 acceptance row changes behavior. retainGreen writes an equal record (Grade set, Oracle empty).\nNew findings: none.\nThe axis ran check-current from the primary checkout, which refused; the orchestrator ran bench preflight evidence sha256:8544c861 --check-current in FT314-gate-reuse and it returned current=true at 9c9c014a5.\n"
          },
          "axis": "Spec",
          "base": "a5f9c772c7b5fb4812ab5b42c1c1611e84bbec4f",
          "tip": "f40e49568deeaff93e72b607d8f94724ff1a5bc0",
          "finding_ids": [],
          "supersedes": [
            "gr-c2-spec-1"
          ]
        },
        {
          "id": "gr-c2-coverage-2",
          "performer": "claude:gr-c2-coverage2-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d0eb2c5d8d0ca8a2643cc6628bf8fda02e53638e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gr-c2-coverage2-20261026",
            "digest": "sha256:4a82487b2e39e4750d5aa45436d7a72677436015441e81dbd775b4848b93c10d",
            "excerpt": "Coverage confirming round, GR-C2, repair deltas 25bd01294..f1fe7d8bc and 20df42e77..f40e49568. sonnet/high. The three tests pass under bench test --run at 9c9c014a5.\nC1 confirmed: TestMalformedEvidenceAnswersNoLookup (internal/gate/grade_key_test.go:168, assertion :198) fails when the shared check or the field check is dropped; recorded reds support it.\nC3 confirmed: TestReusedGreenReplacesAStaleGreenCacheOfItsTree (reused_green_test.go:52, assertion :79) fails on both rows under the recorded probe.\nC4 confirmed: TestEvidenceTokenBindsTheGradeKey (authorization_test.go:81, assertion :86-87) fails when Grade leaves the token.\nNew findings: none.\nThe axis ran check-current from the primary checkout, which refused; the orchestrator ran it in FT314-gate-reuse and it returned current=true at 9c9c014a5.\n"
          },
          "axis": "Coverage",
          "base": "a5f9c772c7b5fb4812ab5b42c1c1611e84bbec4f",
          "tip": "f40e49568deeaff93e72b607d8f94724ff1a5bc0",
          "finding_ids": [],
          "supersedes": [
            "gr-c2-coverage-1"
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
      "from": "sha256:b0f6edda9075e930683fb916045172e340e0ada400ebacc98510788749872dce",
      "to": "sha256:745dd33a081ad71b9ee0cbd09774d8de2bd19290864dcdce214a7796bf5ccb4d",
      "chunk_ids": {
        "GR-C1": [
          "GR-C1"
        ]
      }
    },
    {
      "from": "sha256:745dd33a081ad71b9ee0cbd09774d8de2bd19290864dcdce214a7796bf5ccb4d",
      "to": "sha256:3b9a6ca1ff2998dcbf24b991500bcc26ad4419bed38ebfdb7cd7ad81a9187f70",
      "chunk_ids": {
        "GR-C1": [
          "GR-C1"
        ]
      }
    },
    {
      "from": "sha256:3b9a6ca1ff2998dcbf24b991500bcc26ad4419bed38ebfdb7cd7ad81a9187f70",
      "to": "sha256:36f992a233ba0bcc938f4bb62a282818e29c760fec821b62412e53d16d56ac9e",
      "chunk_ids": {
        "GR-C1": [
          "GR-C1"
        ],
        "GR-C2": [
          "GR-C2"
        ]
      }
    },
    {
      "from": "sha256:36f992a233ba0bcc938f4bb62a282818e29c760fec821b62412e53d16d56ac9e",
      "to": "sha256:99a53aa7b6eaee91fb9b3ac47a50d59b64e1313ddd3eb348932b26da0b60be7c",
      "chunk_ids": {
        "GR-C1": [
          "GR-C1"
        ],
        "GR-C2": [
          "GR-C2"
        ]
      }
    }
  ],
  "folds": [
    {
      "id": "F1",
      "after": "GR-C1",
      "from": "9777f212d75559847ddc797441539ca9ea82d693",
      "commit": "a5f9c772c7b5fb4812ab5b42c1c1611e84bbec4f",
      "incoming": "59897559fc09fc158277b4219064d20062eaa30b",
      "source_digest": "ca3e7b301f1a455a1c2a903c739db24e4d5db4be",
      "resolution": [],
      "resolution_digest": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "test_packages": [],
      "reviews": []
    },
    {
      "id": "F2",
      "after": "GR-C2",
      "from": "1aa7dcdc97da575ce95104f8e535d4115e7b06ac",
      "commit": "1c9f15870f7afca3b123065ebac24b2898c8dfb2",
      "incoming": "bfe85e9e9d34cc3a08deef457c3b829db3c997c8",
      "source_digest": "88a2a177ac21fa0a45c7769cfe2953b20cbdb8a5",
      "resolution": [],
      "resolution_digest": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "test_packages": [],
      "reviews": []
    }
  ]
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

## GR-C1 repair witness reds

The ticket 01 repair session reported these reds in its `Witness return:` field for repair commit 71fa9105a:

- P2: TestOracleIdentityOpensWithoutTheRunningExecutable ran against the code before the repair and failed with "subject closed = true, reason "", want open with "running executable unavailable"".
- P1 and C3: the linked case of TestGradeIdentityLeavesOutTheCheckoutPath failed on the old code with "want one grade identity across checkouts" for `linked = true`.
- C1: the new expectation of TestGradeIdentityFollowsTheRunBinarySource passes on the code before the repair, so two `bench probe --run` mutants demonstrated its red. The mutant that changes `return kitRoot(storageRoot)` to `return runtimeRoot` in `internal/gate/gate.go` made the test fail. The mutant that changes `sourceKind := res.Kind` to `sourceKind := ProspectiveGateSh` in `internal/gate/subject.go` also made it fail. Both probes restored their subject.

The orchestrator wrote these reds after it prepared the confirming charge. A `bench learning` entry records that slip.

## GR-C2 ticket 02 author witness reds

The ticket 02 author probed `internal/gate` at its ticket work. Each mutant below made its named test fail, and each probe restored its subject.

- The engine mutant `frame(h, plan.Identity)` to `frame(h, plan.Oracle)` restores the old key that the checkout path binds. It failed TestGradeKeyCrossesCheckouts (GR10, GR11) at `grade_key_test.go:68`.
- The same mutant failed TestProspectiveEvidenceKeysToTheTreeAndRunBinary (GR13) at `prospective_owner_test.go:221`.
- The same mutant failed TestGradeKeyBindsThePosture (GR14) at `grade_key_test.go:80`.
- The same mutant failed TestCompleteCheckpointReusesAnUnpurposedGreenOfItsTree (GR16) at `grade_key_test.go:103`.
- The same mutant failed TestGradeKeyRedRemovesAnotherCheckoutsGreen (GR18) at `grade_key_test.go:138`.
- An omitted `frame(h, plan.Tree)` failed TestProspectiveEvidenceKeysToTheTreeAndRunBinary (GR12) at `prospective_owner_test.go:216`.
- An omitted `frameBaselineSchedule(id, identityRoot)` failed TestGradeKeyBindsTheBaselineSchedule (GR15) at `grade_key_test.go:89`.
- An expiry test changed to `if false {` failed TestGradeKeyKeepsTheSourceWindow (GR17) at `grade_key_test.go:122`.
- A dropped grade comparison failed TestGradeKeyRefusesAMisnamedRecord (GR20) at `grade_key_test.go:163`.
- After the journey-test fix at c702eb07d, the mutant `if !evidence.ReusableGreen {` to `if true {` in `engine.go` failed TestLandCommandPublicRealGitJourney/clean at `land_journey_test.go:83`.

The planned reuse probe bit at `grade_key_test.go:68` at c702eb07d.

## GR-C2 ticket 03 author witness reds

The ticket 03 author probed `internal/gate` with `bench probe --run`, and each probe restored its subject.

- A cache write changed to `return nil` failed TestReusedGreenRecordsTheGateCache (GR22) at `reused_green_test.go:24`. It also failed TestReusedGreenReplacesARedCacheOfItsTree (GR47) at `reused_green_test.go:44`.
- A write rule that compares only `CachedTree` failed TestReusedGreenReplacesARedCacheOfItsTree (GR47) at `reused_green_test.go:44`.
- An ignored write error failed TestReusedGreenCacheWriteFailureRefuses (GR24) at `durable_replacement_test.go:307`.
- A `lockHeld` result forced to false failed TestReusedGreenBesideAHeldLockWritesNothing (GR44) at `reused_green_test.go:56`.
- Both cache checks changed to `if false` failed TestGateRunGreenReusesRecordedVerdict (GR23) at `run_outcomes_test.go:71`. The two checks mask each other, so this joint probe is the witness for both.

The ticket 02 verification entries carry the ids `gr-t02--<requirement>`, because `bench record ticket` adds its own hyphen to the prefix.

## GR-C2 ticket 02 repair witness reds

The ticket 02 repair session reported these reds for repair commit f1fe7d8bc, and each probe restored its subject.

- C1: an omitted `validateStoreRecordCommon(*record) != nil ||` in `engine.go` failed TestMalformedEvidenceAnswersNoLookup at `grade_key_test.go:198` for the wrong-schema and wrong-tree rows.
- C1: an omitted `requireObjectFields(read.data, evidenceFields) != nil ||` failed the same test at `grade_key_test.go:198` for the legacy-oracle and missing-grade rows.
- C4: an omitted `inspection.Grade` term in `evidenceToken` failed TestEvidenceTokenBindsTheGradeKey at `authorization_test.go:87`.
- The planned reuse probe bit at `grade_key_test.go:68`.

## GR-C2 ticket 03 repair witness reds

The ticket 03 repair session reported this red for repair commit f40e49568. The probe restored its subject.

- C3: a `cacheDescribesFreshGreen` body that also accepts any same-tree green cache skips the write. It failed both rows of TestReusedGreenReplacesAStaleGreenCacheOfItsTree at `reused_green_test.go:79`, which are the expired row and the drifted-oracle row.
- S2 and S3 change no behavior, and the whole `internal/gate` package passes at f40e49568.
