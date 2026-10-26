# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/no-go-schedule/spec.md",
  "plan_digest": "sha256:fe8eec05cb96401f72af5a70e983c0f599e7a6c35d83366b8c15758f6f77e8be",
  "implementation_session": "",
  "chunks": [
    {
      "id": "NG-C1",
      "base": "d30ccffeac02152febd38db38dfccfdcf0f6cf7f",
      "tip": "9768ad919f94d609cd4e32853833d0f0fb31443f",
      "plan_digest": "sha256:11ed1fc05e800de44f46a007e13f3d2a5f224fab1ec958ed43cdcaf5fd8708f8",
      "source_digest": "6f1577e3ac031bfcffccbb20ad32254da0ef9ee0",
      "acceptance_rows": [
        "NG1",
        "NG2",
        "NG3",
        "NG4"
      ],
      "verification": [
        {
          "id": "t01-ticket-1-check-1",
          "performer": "claude:ngs-t01-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eb1341d34a4433afe8b6c7856211f62971757177",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ngs-t01-author-20261025",
            "digest": "sha256:07819a85f9898b12010b6603dda28329ecb6590c0040295487d977b1caadfcdd",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,64283,590\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
          "requirement": "ticket-1-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-2",
          "performer": "claude:ngs-t01-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eb1341d34a4433afe8b6c7856211f62971757177",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ngs-t01-author-20261025",
            "digest": "sha256:2ef9de0bcc32b863ebe4a896ec81f3352c1de3e921219a23a50607e527bcfdd9",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13386,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
          "requirement": "ticket-1-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-probe",
          "performer": "claude:ngs-t01-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eb1341d34a4433afe8b6c7856211f62971757177",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ngs-t01-author-20261025",
            "digest": "sha256:ab2ca06d9e31e8c328a36359d2d3aa1090b1e66f65dc731d8a7d5d8e3a032f55",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,swap,failed,1,TestOracleIdentityIgnoresTheTree,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,590\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,460\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,64459,590\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityIgnoresTheTree,\"oracle_identity_test.go:39: identity = 6abcf6ab77289023b8ec2e4521c3a7353fb34c08363560c135f24f9dd58e6ffe and de05407ab4ef765295e74e42747e1ed13e79aa5b360ab8ae66df730fe0c2f58d, want one identity across trees\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
          "requirement": "ticket-1-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "In buildSubjectForTree, frame the tree into the identity digest too. TestOracleIdentityIgnoresTheTree must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ngs-t01-author-20261025",
              "digest": "sha256:ab2ca06d9e31e8c328a36359d2d3aa1090b1e66f65dc731d8a7d5d8e3a032f55",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,swap,failed,1,TestOracleIdentityIgnoresTheTree,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,590\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,460\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,64459,590\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityIgnoresTheTree,\"oracle_identity_test.go:39: identity = 6abcf6ab77289023b8ec2e4521c3a7353fb34c08363560c135f24f9dd58e6ffe and de05407ab4ef765295e74e42747e1ed13e79aa5b360ab8ae66df730fe0c2f58d, want one identity across trees\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01r1-ticket-1-check-1",
          "performer": "claude:ngs-t01-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6f1577e3ac031bfcffccbb20ad32254da0ef9ee0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ngs-t01-repair1-20261025",
            "digest": "sha256:6bb9160ad6be6e32514bdf8cd65af87dc2edccd82139ddfc1ec271322a1c419c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,64297,597\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "9768ad919f94d609cd4e32853833d0f0fb31443f",
          "requirement": "ticket-1-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-2",
          "performer": "claude:ngs-t01-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6f1577e3ac031bfcffccbb20ad32254da0ef9ee0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ngs-t01-repair1-20261025",
            "digest": "sha256:dd5b8de655f090720506716af7f59f24126200947607809996d21854824b01df",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13407,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "9768ad919f94d609cd4e32853833d0f0fb31443f",
          "requirement": "ticket-1-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-probe",
          "performer": "claude:ngs-t01-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6f1577e3ac031bfcffccbb20ad32254da0ef9ee0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ngs-t01-repair1-20261025",
            "digest": "sha256:9b789dd1a1b6a8120cc78c1797ac36eb6beb6fd2090aa8349fdc221510c265ee",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,swap,failed,1,TestOracleIdentityIgnoresTheTree,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,597\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,460\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,64493,597\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityIgnoresTheTree,\"oracle_identity_test.go:53: identity = 3cd56d1c2b363d8a542f300a8719450e12a2920bf7819568796187d3f539913b and 1ff28b742cc2d4565d5da72e8efc7b6097ac6e7c5bf08f31b3c67e0ef6e954e2, want one identity across trees\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "9768ad919f94d609cd4e32853833d0f0fb31443f",
          "requirement": "ticket-1-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "In buildSubjectForTree, frame the tree into the identity digest too. TestOracleIdentityIgnoresTheTree must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ngs-t01-repair1-20261025",
              "digest": "sha256:9b789dd1a1b6a8120cc78c1797ac36eb6beb6fd2090aa8349fdc221510c265ee",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,swap,failed,1,TestOracleIdentityIgnoresTheTree,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,597\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,460\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,64493,597\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityIgnoresTheTree,\"oracle_identity_test.go:53: identity = 3cd56d1c2b363d8a542f300a8719450e12a2920bf7819568796187d3f539913b and 1ff28b742cc2d4565d5da72e8efc7b6097ac6e7c5bf08f31b3c67e0ef6e954e2, want one identity across trees\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c1r1-standards",
          "performer": "claude:ngs-c1-standards-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "eb1341d34a4433afe8b6c7856211f62971757177",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c1-standards-20261025",
            "digest": "sha256:0c62646a2acec3f28d69d0064cbe87b2933e2b40e2e461a974852847b92ae43d",
            "excerpt": "Axis: Standards, NG-C1, ticket 01. Session claude:ngs-c1-standards-20261025. Frozen pair d30ccffe..8c68fe1f.\nCount: 0 findings. Worst issue: none.\nRe-derived: oracle frame order unchanged (subject.go:73-81; tree to h only, all other frames through w); one pass feeds both digests; posture read from strict() and framed into id only (subject.go:131-132); subject.go 456/513 lines; three callers of buildSubjectForTree, field additive; tests map one-to-one to NG1-NG3.\nAdvice (no ID): (1) oracle_identity_test.go:37,52,71 open comments with row IDs NG1-NG3; dropping them is cleaner. (2) verdict.go:~88 comment names a landing consumer that later tickets add.\nImplementation command contribution: none.\n"
          },
          "axis": "Standards",
          "base": "d30ccffeac02152febd38db38dfccfdcf0f6cf7f",
          "tip": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1r1-spec",
          "performer": "claude:ngs-c1-spec-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "eb1341d34a4433afe8b6c7856211f62971757177",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c1-spec-20261025",
            "digest": "sha256:b6c3fd85f7baa327774b376978e46e04b2c7fc4a709926bda2b25f51af700c5e",
            "excerpt": "Axis: Spec, NG-C1, ticket 01, rows NG1-NG4. Session claude:ngs-c1-spec-20261025. Frozen pair d30ccffe..8c68fe1f.\nCount: 0 findings. Worst issue: none.\nNG4 audited in buildSubjectForTree (subject.go:73-133): h receives the same seven leading frames in the same order; the tree goes to h only; id alone gets requireCapabilitiesEnv and strict() after s.Oracle. NG1-NG3 tests match their rows. Plan commits c0d1266d6, f80bbf194, bcdb41cb7: citation corrections and the v2 split hold; union of checks equal; probes recorded by a bench learning entry; no material acceptance change.\nAdvice: none.\n"
          },
          "axis": "Spec",
          "base": "d30ccffeac02152febd38db38dfccfdcf0f6cf7f",
          "tip": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1r1-coverage",
          "performer": "claude:ngs-c1-coverage-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "eb1341d34a4433afe8b6c7856211f62971757177",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c1-coverage-20261025",
            "digest": "sha256:31cad880f835500ad9de5be807488cd733a9a2871ef53402d662e10619175d67",
            "excerpt": "Axis: Coverage, NG-C1, ticket 01, rows NG1-NG4. Session claude:ngs-c1-coverage-20261025. Frozen pair d30ccffe..8c68fe1f.\nCount: 1 finding. Worst issue: C1.\nC1 severity 4, confidence 6, auto-fix, ticket 01 (internal/gate/oracle_identity_test.go): no test varies a non-tree, non-tool input of the identity digest. A w->h swap at subject.go:77, :78, :80-82, :89, :99-101, or :104-105 drops an input from the identity and leaves NG1-NG3 green. Proposed: one table-driven test that changes each input (policy, identity root, BENCH_GATE, PATH, a declared environment value, a declared path's bytes, the manifest text) and asserts that both digests move. Implementation command contributed: yes; the ticket lacks a row for \"any non-tree oracle input changes the identity\".\nAdvice (no ID): sameSubject skips Identity (sound today); NG2 survives framing raw os.Getenv; NG4 has no golden value (review-owned); open subjects get an Identity; spec cites subject.go:73 and :123, now stale.\n"
          },
          "axis": "Coverage",
          "base": "d30ccffeac02152febd38db38dfccfdcf0f6cf7f",
          "tip": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        }
      ],
      "confirmations": [
        {
          "id": "c1r1-confirmation",
          "performer": "claude:ngs-orchestrator-20261025",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6f1577e3ac031bfcffccbb20ad32254da0ef9ee0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:ngs-orchestrator-20261025:c1-probe",
            "digest": "sha256:cf37a82458551c87994c502007bdad3e88e1d7dde63f67f2864272d3bde224e8",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,insert,failed,7,TestOracleIdentityAndOracleShareTheClosure/PATH,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestOracleIdentity,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1908\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,49,10\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/PATH,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across PATH, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/baseline_runner,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across baseline runner, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/declared_path_bytes,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across declared path bytes, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/environment_value,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across environment value, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/identity_root,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across identity root, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/policy,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across policy, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/tool_bytes,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across tool bytes, want two identities\",1\nskips[0]{package,test,reason}:\n"
          },
          "base": "d30ccffeac02152febd38db38dfccfdcf0f6cf7f",
          "reviewed_tip": "8c68fe1fcc1111a64f6c7f5f5fcbe9cbf0e4bd54",
          "repaired_tip": "9768ad919f94d609cd4e32853833d0f0fb31443f",
          "finding_ids": [
            "C1"
          ],
          "supersedes": [],
          "probe": {
            "mutation": "insert id.Reset() before the identity posture frames in internal/gate/subject.go",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:ngs-orchestrator-20261025:c1-probe",
              "digest": "sha256:cf37a82458551c87994c502007bdad3e88e1d7dde63f67f2864272d3bde224e8",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/subject.go,insert,failed,7,TestOracleIdentityAndOracleShareTheClosure/PATH,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestOracleIdentity,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1908\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,49,10\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/PATH,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across PATH, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/baseline_runner,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across baseline runner, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/declared_path_bytes,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across declared path bytes, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/environment_value,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across environment value, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/identity_root,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across identity root, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/policy,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across policy, want two identities\",1\n  github.com/gibbonmi/bench/internal/gate,TestOracleIdentityAndOracleShareTheClosure/tool_bytes,\"oracle_identity_test.go:146: identity = b936ce0d10832c24f6c7ac516e9fb524d5cb8279265c987f40f4fb5009ba7f79 across tool bytes, want two identities\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ]
    },
    {
      "id": "NG-C2",
      "base": "9768ad919f94d609cd4e32853833d0f0fb31443f",
      "tip": "5aaad0ad36af9a7ee95845845585a233c6713621",
      "plan_digest": "sha256:fe8eec05cb96401f72af5a70e983c0f599e7a6c35d83366b8c15758f6f77e8be",
      "source_digest": "e440572cfc3420f5cb25f2084ba68ef7711bab1a",
      "acceptance_rows": [
        "NG5",
        "NG6",
        "NG7",
        "NG8",
        "NG9",
        "NG10",
        "NG11",
        "NG12",
        "NG72",
        "NG81",
        "NG13",
        "NG14",
        "NG15",
        "NG16",
        "NG17",
        "NG18",
        "NG19",
        "NG20",
        "NG21",
        "NG22",
        "NG23",
        "NG24",
        "NG25",
        "NG26",
        "NG27",
        "NG28",
        "NG71"
      ],
      "verification": [
        {
          "id": "t03-ticket-3-check-1",
          "performer": "claude:ngs-t03-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e440572cfc3420f5cb25f2084ba68ef7711bab1a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t03-author-20261025",
            "digest": "sha256:1a8ec55582141f730ae349ccc344b015da665b1f050fb9b64b492efa86e24689",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,69295,623\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5aaad0ad36af9a7ee95845845585a233c6713621",
          "requirement": "ticket-3-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t03-ticket-3-check-2",
          "performer": "claude:ngs-t03-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e440572cfc3420f5cb25f2084ba68ef7711bab1a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t03-author-20261025",
            "digest": "sha256:5a6bd5a5bfac80246edf34c70617c090655953a6d9fc2f3d36a031bf8663c8d2",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,1341,179\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5aaad0ad36af9a7ee95845845585a233c6713621",
          "requirement": "ticket-3-check-2",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "t03-ticket-3-check-3",
          "performer": "claude:ngs-t03-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e440572cfc3420f5cb25f2084ba68ef7711bab1a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t03-author-20261025",
            "digest": "sha256:c1de50ff3562ef0c41f78e82ce8e79f84a0ad2d2584195c5aa0c6e2cf16aaea4",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13382,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5aaad0ad36af9a7ee95845845585a233c6713621",
          "requirement": "ticket-3-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t03-ticket-3-probe",
          "performer": "claude:ngs-t03-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e440572cfc3420f5cb25f2084ba68ef7711bab1a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t03-author-20261025",
            "digest": "sha256:0235626aa4eef90bbb99d7b19ccc55f96b1ab8f1db9ec21a7a2a0ff052c7bc07",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/reader_table.go,omit,failed,1,TestReaderTableRequiresTheConformanceWitness,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,623\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,523\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,69326,623\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestReaderTableRequiresTheConformanceWitness,\"reader_table_test.go:167: derivation = map[.agents/:[internal/conformance] .bench/BENCH.md:[internal/conformance] .bench/commitment.json:[internal/conformance] CHANGELOG.md:[internal/conformance] ROADMAP.md:[internal/conformance] capture/:[… (454 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5aaad0ad36af9a7ee95845845585a233c6713621",
          "requirement": "ticket-3-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "In the reader-table derivation, remove the internal/conformance witness refusal. TestReaderTableRequiresTheConformanceWitness must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ngs-t03-author-20261025",
              "digest": "sha256:0235626aa4eef90bbb99d7b19ccc55f96b1ab8f1db9ec21a7a2a0ff052c7bc07",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/reader_table.go,omit,failed,1,TestReaderTableRequiresTheConformanceWitness,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,623\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,523\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,69326,623\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestReaderTableRequiresTheConformanceWitness,\"reader_table_test.go:167: derivation = map[.agents/:[internal/conformance] .bench/BENCH.md:[internal/conformance] .bench/commitment.json:[internal/conformance] CHANGELOG.md:[internal/conformance] ROADMAP.md:[internal/conformance] capture/:[… (454 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t02-ticket-2-check-1",
          "performer": "claude:ngs-t02-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c2b5ce7e36dfc13afe9c9cf8c97ea445863c9468",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t02-author-20261025",
            "digest": "sha256:5d17c0a3f9c2bd604601317d6fa9c6ab13582b83799d153c23f07ce46291a69c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,68080,606\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e2320d7138b93c708ded2bff0a6d85db598cfd19",
          "requirement": "ticket-2-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-2",
          "performer": "claude:ngs-t02-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c2b5ce7e36dfc13afe9c9cf8c97ea445863c9468",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t02-author-20261025",
            "digest": "sha256:1b8be29bebd5b2aaba7bd5ab1066f1665bdebbfe81d94603dda9e57214dc3a03",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4336,43\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e2320d7138b93c708ded2bff0a6d85db598cfd19",
          "requirement": "ticket-2-check-2",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-3",
          "performer": "claude:ngs-t02-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c2b5ce7e36dfc13afe9c9cf8c97ea445863c9468",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t02-author-20261025",
            "digest": "sha256:daded535e600af67852f854731c0f8e9448cd416757b516a2258acd2cdf92af2",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,379,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e2320d7138b93c708ded2bff0a6d85db598cfd19",
          "requirement": "ticket-2-check-3",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-4",
          "performer": "claude:ngs-t02-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c2b5ce7e36dfc13afe9c9cf8c97ea445863c9468",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t02-author-20261025",
            "digest": "sha256:542656d86d539e1e3a0e080a3d71d4862c9fc34b4137368b35355d17cbca4de6",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,37909,1143\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "source_commit": "e2320d7138b93c708ded2bff0a6d85db598cfd19",
          "requirement": "ticket-2-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-5",
          "performer": "claude:ngs-t02-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c2b5ce7e36dfc13afe9c9cf8c97ea445863c9468",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t02-author-20261025",
            "digest": "sha256:ebb21f36389050ff9c6e84ac652a29b65522276811fdbb0abe9a5ee3e6d94cc5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13273,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e2320d7138b93c708ded2bff0a6d85db598cfd19",
          "requirement": "ticket-2-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-probe",
          "performer": "claude:ngs-t02-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c2b5ce7e36dfc13afe9c9cf8c97ea445863c9468",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t02-author-20261025",
            "digest": "sha256:2794932d9f7167f75202834dc8ed3b75194241d020e27238950c34d0c0ea346d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/testlog_exec.go,swap,failed,2,TestGateGoExecAddsTheTestLog,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,606\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,586\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,68453,606\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateGoExecAddsTheTestLog,\"testlog_exec_test.go:47: probe arguments = [\\\"-a\\\" \\\"b c\\\"], want [\\\"-test.testlogfile=/dev/shm/VZFZMN/t/TestGateGoExecAddsTheTestLog424366856/001/pkg%2Fsub.log\\\" \\\"-a\\\" \\\"b c\\\"]\",1\n  github.com/gibbonmi/bench/internal/gate,TestWholeTestPhaseWritesOneLogPerPackage,\"testlog_exec_test.go:132: log of a = \\\"\\\", <nil>; want a Go test log\",4\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e2320d7138b93c708ded2bff0a6d85db598cfd19",
          "requirement": "ticket-2-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "In the gate-go exec shim, drop the -test.testlogfile operand. TestGateGoExecAddsTheTestLog must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ngs-t02-author-20261025",
              "digest": "sha256:2794932d9f7167f75202834dc8ed3b75194241d020e27238950c34d0c0ea346d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/testlog_exec.go,swap,failed,2,TestGateGoExecAddsTheTestLog,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,606\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,586\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,68453,606\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateGoExecAddsTheTestLog,\"testlog_exec_test.go:47: probe arguments = [\\\"-a\\\" \\\"b c\\\"], want [\\\"-test.testlogfile=/dev/shm/VZFZMN/t/TestGateGoExecAddsTheTestLog424366856/001/pkg%2Fsub.log\\\" \\\"-a\\\" \\\"b c\\\"]\",1\n  github.com/gibbonmi/bench/internal/gate,TestWholeTestPhaseWritesOneLogPerPackage,\"testlog_exec_test.go:132: log of a = \\\"\\\", <nil>; want a Go test log\",4\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c2r1-standards",
          "performer": "claude:ngs-c2-standards-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e440572cfc3420f5cb25f2084ba68ef7711bab1a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c2-standards-20261025",
            "digest": "sha256:a138d14875ed1233c5281e5612d447dc495f1d17455576af4e4f56abc77ba5df",
            "excerpt": "Axis: Standards, NG-C2, tickets 02 and 03. Session claude:ngs-c2-standards-20261025. Frozen pair 9768ad91..5aaad0ad.\nCount: 2 findings. Worst issue: S1.\nS1 severity 5, confidence 7, auto-fix, ticket 03 (internal/gate/reader_store.go:35-38,137-145): testRunFlags {-count,-timeout,-exec} restates the run flags that BaseTestArgv (phases.go:230) and traceExecOperand (testlog_exec.go:29) emit; ticket 03 line 17 forbids a restated list. Fix: the traced phase builder exposes the build flags as one value that BaseTestArgv and testPackages both read.\nS2 severity 4, confidence 8, auto-fix, ticket 03 (internal/git/tree.go:175-176): symlinkMode \"120000\" duplicates lane_select.go:110 (and settlepolicy.go:47). Fix: export the mode from internal/git and point lane_select.go at it.\nAdvice (no ID): three test-binary trampolines (processtest init, the shim init, Replace's probe) could share one processtest entry; optional.\n"
          },
          "axis": "Standards",
          "base": "9768ad919f94d609cd4e32853833d0f0fb31443f",
          "tip": "5aaad0ad36af9a7ee95845845585a233c6713621",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c2r1-spec",
          "performer": "claude:ngs-c2-spec-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e440572cfc3420f5cb25f2084ba68ef7711bab1a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c2-spec-20261025",
            "digest": "sha256:3c0c11edb7ed1a2169b8d4dc304259779b2d34aa709f5b4c27e2435d88451649",
            "excerpt": "Axis: Spec, NG-C2, tickets 02 and 03, rows NG5-NG28, NG71, NG72, NG81. Session claude:ngs-c2-spec-20261025. Frozen pair 9768ad91..5aaad0ad.\nCount: 0 findings. Worst issue: none.\nEach row's named test exists and drives its behavior. Shim, key, traced argv, quoting, env strip, derivation rule order and refusals, persistence after the green verdict only, log-directory removal, and the marker-keeping prune all match the spec.\nAdvice (no ID): testPackages runs go list with the ambient env, not gateEnv; a run-binary path with both quote marks has no cmd/go spelling (a red run, not a pass); branch_native_phases_test.go:104 derives its expectation from traceExecOperand; no negative test that other phases lack BENCH_GATE_TESTLOG; no direct unit test of TreeSymlinks; the root-package skip changes no table today.\n"
          },
          "axis": "Spec",
          "base": "9768ad919f94d609cd4e32853833d0f0fb31443f",
          "tip": "5aaad0ad36af9a7ee95845845585a233c6713621",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c2r1-coverage",
          "performer": "claude:ngs-c2-coverage-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e440572cfc3420f5cb25f2084ba68ef7711bab1a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c2-coverage-20261025",
            "digest": "sha256:8b7427ea547372232213e7505654bdec4e45ef940ecd1959c542423369a83e5a",
            "excerpt": "Axis: Coverage, NG-C2, tickets 02 and 03. Session claude:ngs-c2-coverage-20261025. Frozen pair 9768ad91..5aaad0ad.\nCount: 6 findings. Worst issue: C1.\nC1 sev 5, conf 7, auto-fix, ticket 03: no test pins that the table persists only after the green verdict persists (spec edge inventory; run_transaction.go:239-244).\nC2 sev 5, conf 5, ask-user, ticket 03: go list runs under the operator's environment (reader_store.go:146 via stepOutput) instead of the closed oracle environment, so ambient GOFLAGS/GOOS can change the package list; no test.\nC3 sev 4, conf 8, auto-fix, ticket 03: no test of an absolute chdir followed by a relative open (reader_table.go:93).\nC4 sev 3, conf 8, auto-fix, ticket 03: the empty-name refusal (reader_table.go:86-88) has no test.\nC5 sev 4, conf 5, ask-user, ticket 03: familyOf also claims a family's own directory (reader_families.go:42); not in the spec text; no test.\nC6 sev 4, conf 6, ask-user, ticket 03: the eleven-family declaration is not pinned member by member (reader_families.go:10-22).\nAdvice (no ID): persisted Children not asserted; link chains and cycles untested; prune with no marker and tied mtimes untested; withTestLog reads the ambient env in-process; testLogFile resolves links but the derivation compares against the unresolved root.\n"
          },
          "axis": "Coverage",
          "base": "9768ad919f94d609cd4e32853833d0f0fb31443f",
          "tip": "5aaad0ad36af9a7ee95845845585a233c6713621",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
            "C5",
            "C6"
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
      "from": "sha256:f7c4b9dfcb60734c5667d0a2306b537d6995e9fb0298efa86eb1c64756040342",
      "to": "sha256:11ed1fc05e800de44f46a007e13f3d2a5f224fab1ec958ed43cdcaf5fd8708f8",
      "chunk_ids": {
        "NG-C1": [
          "NG-C1"
        ]
      }
    },
    {
      "from": "sha256:11ed1fc05e800de44f46a007e13f3d2a5f224fab1ec958ed43cdcaf5fd8708f8",
      "to": "sha256:fe8eec05cb96401f72af5a70e983c0f599e7a6c35d83366b8c15758f6f77e8be",
      "chunk_ids": {
        "NG-C1": [
          "NG-C1"
        ]
      }
    }
  ]
}
```

## Repair state

NG-C1 consumed repair cycle 1 of 2 after its initial review (round r1). No reviewer extension applies.

NG-C2 consumed repair cycle 1 of 2 after its initial review (round r1). No reviewer extension applies.

## Standards

NG-C1 round r1: 0 findings. Worst issue: none.

NG-C2 round r1: 2 findings. Worst issue: S1, a restated list of the run flags of the traced test phase.

- S1 | auto-fix | confidence 7 | `internal/gate/reader_store.go:35-38` | `testRunFlags` restates the run flags that the traced phase builder emits. Ticket 03 line 17 forbids a restated list.
- S2 | auto-fix | confidence 8 | `internal/git/tree.go:175-176` | The symbolic-link mode gets a second production copy beside `internal/gate/lane_select.go:110`.

## Spec

NG-C1 round r1: 0 findings. Worst issue: none.

NG-C2 round r1: 0 findings. Worst issue: none.

## Coverage

NG-C1 round r1: 1 finding. Worst issue: C1, no test varies an identity input other than the tree, the posture, and the tool bytes.

- C1 | auto-fix | confidence 6 | `internal/gate/oracle_identity_test.go` | A writer swap at `internal/gate/subject.go:77-105` drops an identity input and leaves NG1-NG3 green. Ticket 01 line 9 gives the identity every oracle frame except the tree.

NG-C2 round r1: 6 findings. Worst issue: C1, no test pins that the table persists only after the green verdict persists.

- C1 | auto-fix | confidence 7 | `internal/gate/run_transaction.go:239-244` | No test pins the persistence order that the spec's edge inventory states.
- C2 | auto-fix | confidence 5 | `internal/gate/reader_store.go:146` | The package listing runs under the operator's environment, not the closed oracle environment. The axis asked for the user; the orchestrator routed the environment part as auto-fix.
- C3 | auto-fix | confidence 8 | `internal/gate/reader_table.go:93` | No test reads a relative path after an absolute `chdir`.
- C4 | auto-fix | confidence 8 | `internal/gate/reader_table.go:86-88` | No test feeds a log line with an empty name.
- C5 | no-op | confidence 5 | `internal/gate/reader_families.go:42` | The consultants kept the claim of a family's own directory, and the spec text now states it.
- C6 | no-op | confidence 6 | `internal/gate/reader_families.go:10-22` | A test that pins each family duplicates the declaration, which AGENTS.md forbids, and review owns that declaration.

## Repair witnesses

The NG-C1 repair commit `9768ad919` makes `TestOracleIdentityAndOracleShareTheClosure` a table of seven inputs. Each mutant below swapped the identity writer for the oracle writer in `internal/gate/subject.go`. Each one read `bit` with `restored` `yes`:

- The policy frame failed the `policy` case.
- The identity-root frame failed the `identity_root` case.
- The resolution loop failed the `PATH` case.
- The prospective preparation failed the `baseline_runner` case.
- The environment frames failed the `environment_value` case.
- A `hashRepoPath` call replaced with `error(nil)` failed the `declared_path_bytes` case.

The swap of the `path` frames stayed silent. The orchestrator accepts it as an equivalent mutant, because `manifestIdentity` already frames each declared path name into both digests. The coordinator probe inserted `id.Reset()` before the posture frames, and seven cases failed.

## Advice

- The row-ID prefixes in the comments of `internal/gate/oracle_identity_test.go` are provenance; dropping them is cleaner.
- `sameSubject` does not compare `Identity`; a later ticket that compares the field must add it there or check `Closed` first.
- NG2 would survive a frame of the raw `BENCH_REQUIRE_CAPABILITIES` value in place of `strict()`.
- The spec cites `internal/gate/subject.go:73` and `:123`, which moved in this delta.
- Three test binaries serve as helper programs through their own `init` trampolines; one shared `processtest` entry could replace them.
- The persisted child rows, link chains and cycles, and a prune with no green marker have no direct test.
- A run-binary path that holds both quote marks has no `cmd/go` spelling, so the run goes red.

## Consultant decisions

The ticket 03 author asked for three decisions. The consultants (opus / high, sonnet / high) kept the directory claim of `familyOf` and the prune by modification time. They split on the root test package. The tie-breaker (fable / medium) chose to amend: the own-directory rule does not apply to the root package, and row NG82 covers the change.
