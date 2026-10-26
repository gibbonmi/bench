# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/gate-reuse/spec.md",
  "plan_digest": "sha256:745dd33a081ad71b9ee0cbd09774d8de2bd19290864dcdce214a7796bf5ccb4d",
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
