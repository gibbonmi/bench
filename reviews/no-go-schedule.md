# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/no-go-schedule/spec.md",
  "plan_digest": "sha256:d1077d01a3ec7991d40fabd7a7b56e387cb8b270c52105ccd278f5b050ffea11",
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
      "tip": "8885e8d8ab5717f8da8b3e6eba35d3a176af0dab",
      "plan_digest": "sha256:68ef1b1ef031160d9eec83650abcea4fd38721acdb6654476a69446626e82429",
      "source_digest": "b5a4125d38d9546033c0264cf0a91e491a8032d5",
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
        "NG71",
        "NG82"
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
        },
        {
          "id": "t03r1-ticket-3-check-1",
          "performer": "claude:ngs-t03-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b5a4125d38d9546033c0264cf0a91e491a8032d5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t03-repair1-20261025",
            "digest": "sha256:ce1be5f00ad0fd2146a46924b09ab4b6a58c67b71430fb9b45fa5c731f4ab8d5",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,70161,626\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8885e8d8ab5717f8da8b3e6eba35d3a176af0dab",
          "requirement": "ticket-3-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t03r1-ticket-3-check-2",
          "performer": "claude:ngs-t03-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b5a4125d38d9546033c0264cf0a91e491a8032d5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t03-repair1-20261025",
            "digest": "sha256:1cd04fea799c67744b14b77464813667326efb072bce803b2c86157230296942",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,1351,179\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8885e8d8ab5717f8da8b3e6eba35d3a176af0dab",
          "requirement": "ticket-3-check-2",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "t03r1-ticket-3-check-3",
          "performer": "claude:ngs-t03-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b5a4125d38d9546033c0264cf0a91e491a8032d5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t03-repair1-20261025",
            "digest": "sha256:751d53a21543a06741096484acef7ba69c933ea122e0c0ca666b1fc2ae66fd4b",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13502,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8885e8d8ab5717f8da8b3e6eba35d3a176af0dab",
          "requirement": "ticket-3-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t03r1-ticket-3-probe",
          "performer": "claude:ngs-t03-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b5a4125d38d9546033c0264cf0a91e491a8032d5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t03-repair1-20261025",
            "digest": "sha256:ed1a90352fd82dc06906ae0d0adf33797773a2605cf8ebce433506f090037dc4",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/reader_table.go,omit,failed,1,TestReaderTableRequiresTheConformanceWitness,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,626\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,523\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,70516,626\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestReaderTableRequiresTheConformanceWitness,\"reader_table_test.go:182: derivation = map[.agents/:[internal/conformance] .bench/BENCH.md:[internal/conformance] .bench/commitment.json:[internal/conformance] CHANGELOG.md:[internal/conformance] ROADMAP.md:[internal/conformance] capture/:[… (454 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8885e8d8ab5717f8da8b3e6eba35d3a176af0dab",
          "requirement": "ticket-3-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "In the reader-table derivation, remove the internal/conformance witness refusal. TestReaderTableRequiresTheConformanceWitness must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ngs-t03-repair1-20261025",
              "digest": "sha256:ed1a90352fd82dc06906ae0d0adf33797773a2605cf8ebce433506f090037dc4",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/reader_table.go,omit,failed,1,TestReaderTableRequiresTheConformanceWitness,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,626\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,523\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,70516,626\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestReaderTableRequiresTheConformanceWitness,\"reader_table_test.go:182: derivation = map[.agents/:[internal/conformance] .bench/BENCH.md:[internal/conformance] .bench/commitment.json:[internal/conformance] CHANGELOG.md:[internal/conformance] ROADMAP.md:[internal/conformance] capture/:[… (454 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
        },
        {
          "id": "c2r2-standards",
          "performer": "claude:ngs-c2r2-standards-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "b5a4125d38d9546033c0264cf0a91e491a8032d5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c2r2-standards-20261025",
            "digest": "sha256:266ad3751a8e04bff4bb9fa76d74d5147e12e1770ca57f451f5f827d64c47020",
            "excerpt": "Axis: Standards, NG-C2 confirming round r2. Session claude:ngs-c2r2-standards-20261025. Repair delta 5aaad0ad..8885e8d8.\nCount: 0 findings. Worst issue: none.\nS1 confirmed: testBuildFlags (phases.go:220-222) is the one source, read by BaseTestArgv and testPackages; testRunFlags is gone. S2 confirmed: internal/git exports SymlinkMode and GitlinkMode; lane_select.go:303 reads them. C2 fold: a direct exec.Command in gate production code is allowed (precedent gate_go.go:254,277, runner.go:335; census, process-lifetime-closure, and bounds-policy do not grade a synchronous cmd.Output()).\nAdvice (no ID): an env parameter on stepOutput would remove the re-spelled call shape; go-build-vcs cannot grade a variadic go list argv; older copies of the link modes sit outside the delta.\n"
          },
          "axis": "Standards",
          "base": "9768ad919f94d609cd4e32853833d0f0fb31443f",
          "tip": "8885e8d8ab5717f8da8b3e6eba35d3a176af0dab",
          "finding_ids": [],
          "supersedes": [
            "c2r1-standards"
          ]
        },
        {
          "id": "c2r2-spec",
          "performer": "claude:ngs-c2r2-spec-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "b5a4125d38d9546033c0264cf0a91e491a8032d5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c2r2-spec-20261025",
            "digest": "sha256:3d6dcd41acecda7d97e615a327fd72761dc5af514fa78c525739ccddaa1c37d9",
            "excerpt": "Axis: Spec, NG-C2 confirming round r2. Session claude:ngs-c2r2-spec-20261025. Repair delta 5aaad0ad..8885e8d8 and plan amendment 4d4d7950..621c9f2f.\nCount: 0 findings. Worst issue: none.\nNG82 holds (reader_table.go:101; TestReaderTableCountsTheRootPackageReads), NG16 still holds, the citation cell is valid. Spec line 113 matches familyOf. C2: the listing runs under mergeEnv(plan.Env, phase.Env), the traced run's environment. S1: every phase argv stays byte-identical.\nAdvice (no ID): the directory claim of familyOf has no unit test; testBuildFlags is a mutable package slice, copied at both uses.\n"
          },
          "axis": "Spec",
          "base": "9768ad919f94d609cd4e32853833d0f0fb31443f",
          "tip": "8885e8d8ab5717f8da8b3e6eba35d3a176af0dab",
          "finding_ids": [],
          "supersedes": [
            "c2r1-spec"
          ]
        },
        {
          "id": "c2r2-coverage",
          "performer": "claude:ngs-c2r2-coverage-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "b5a4125d38d9546033c0264cf0a91e491a8032d5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c2r2-coverage-20261025",
            "digest": "sha256:89835154387a708a3649f303f78d2ea34f3d1c3de29433a4f2dcacb5708ac441",
            "excerpt": "Axis: Coverage, NG-C2 confirming round r2. Session claude:ngs-c2r2-coverage-20261025. Repair delta 5aaad0ad..8885e8d8.\nCount: 0 findings. Worst issue: none.\nC1, C2, C3, C4, and NG82: each new test is red-capable against the defect it names (traced: verdict-before-table order at run_transaction.go:239-243; ambient GOFLAGS tag lists an unlogged package; absolute chdir through logPath's IsAbs branch; empty-name refusal at reader_table.go:87-89; root-package skip at reader_table.go:101).\nAdvice (no ID): the lane gitlink branch (lane_select.go:303) and the tree.go:89 gitlink guard have no biting test; both predate the chunk.\n"
          },
          "axis": "Coverage",
          "base": "9768ad919f94d609cd4e32853833d0f0fb31443f",
          "tip": "8885e8d8ab5717f8da8b3e6eba35d3a176af0dab",
          "finding_ids": [],
          "supersedes": [
            "c2r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "NG-C3",
      "base": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
      "tip": "7204223529757f414366245ec4f8db01de7a3f65",
      "plan_digest": "sha256:44a38079f9734bd6a1f106fe16ba748911098f8d5a19e6f5cbd42d7c7d55e6e4",
      "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
      "acceptance_rows": [
        "NG29",
        "NG30",
        "NG31",
        "NG32",
        "NG33",
        "NG34",
        "NG35",
        "NG36",
        "NG79",
        "NG80",
        "NG37",
        "NG38",
        "NG39",
        "NG40",
        "NG41",
        "NG42",
        "NG43",
        "NG44",
        "NG45",
        "NG46",
        "NG47",
        "NG48",
        "NG49",
        "NG50",
        "NG51",
        "NG52",
        "NG53",
        "NG54",
        "NG73",
        "NG74",
        "NG75",
        "NG76",
        "NG77",
        "NG55",
        "NG56",
        "NG57",
        "NG58",
        "NG59",
        "NG60",
        "NG83",
        "NG84"
      ],
      "verification": [
        {
          "id": "t06-ticket-6-check-1",
          "performer": "claude:ngs-t06-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d55a1dd09e3a53f63ffcf69525cbb542990431a8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-author-20261025",
            "digest": "sha256:7ed68123f8f752e4ce636f25ea72a3354c51e027da30cab644c2eefaf7815b46",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,219,16\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ff617775e2ca93ad218f92d8fad7395a2cacad38",
          "requirement": "ticket-6-check-1",
          "command": "bench test --package ./internal/gate/authorization",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-check-2",
          "performer": "claude:ngs-t06-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d55a1dd09e3a53f63ffcf69525cbb542990431a8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-author-20261025",
            "digest": "sha256:cac37e7a7961d4ca7bfb9dee397a4b90ae41f10a7dc8bd0448f2f3789e9fbaf3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,9655,222\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "source_commit": "ff617775e2ca93ad218f92d8fad7395a2cacad38",
          "requirement": "ticket-6-check-2",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-check-3",
          "performer": "claude:ngs-t06-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d55a1dd09e3a53f63ffcf69525cbb542990431a8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-author-20261025",
            "digest": "sha256:4dd45db048d97fd7c3d214a30b3a4a4b0adbbe9cf6d2ba5b4bcc89c19fe1736c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,73772,1515\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ff617775e2ca93ad218f92d8fad7395a2cacad38",
          "requirement": "ticket-6-check-3",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-check-4",
          "performer": "claude:ngs-t06-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d55a1dd09e3a53f63ffcf69525cbb542990431a8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-author-20261025",
            "digest": "sha256:ac6b7ad585e26bf0c4a31e2818a77ddfd72ee54a806b5b50f93afc023560977a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree/landingpolicy,pass,3,46\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ff617775e2ca93ad218f92d8fad7395a2cacad38",
          "requirement": "ticket-6-check-4",
          "command": "bench test --package ./internal/worktree/landingpolicy",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-check-5",
          "performer": "claude:ngs-t06-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d55a1dd09e3a53f63ffcf69525cbb542990431a8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-author-20261025",
            "digest": "sha256:0f4465dcbd109106c9e85ed6b80e9bddc47e155b30de3a6a4cdbb1c559e29370",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,21444,406\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ff617775e2ca93ad218f92d8fad7395a2cacad38",
          "requirement": "ticket-6-check-5",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-check-6",
          "performer": "claude:ngs-t06-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d55a1dd09e3a53f63ffcf69525cbb542990431a8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-author-20261025",
            "digest": "sha256:e7cf5608364d973d3bb73c54e30b65967b8285fe25156f99ab029669bd90c8f3",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13651,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ff617775e2ca93ad218f92d8fad7395a2cacad38",
          "requirement": "ticket-6-check-6",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-probe",
          "performer": "claude:ngs-t06-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d55a1dd09e3a53f63ffcf69525cbb542990431a8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-author-20261025",
            "digest": "sha256:739f5a10a0d5d671b44724d4738bbb6b9b01d28b781acd92995415917531aac0",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/land_marker.go,swap,failed,1,TestReaderLandingLeavesTheGreenMarker,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestReaderLandingLeavesTheGreenMarker,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,450\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,381,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestReaderLandingLeavesTheGreenMarker,\"land_reader_test.go:38: refs/bench/green/main = dc07699aa03d0d9ad9ada6fcffdbaa21b39f92be after a reader green publication, want the prior marker 638ef50fc7f1b1086dd4966e7f774ae3086f49c9\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ff617775e2ca93ad218f92d8fad7395a2cacad38",
          "requirement": "ticket-6-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "In the reviewed landing, advance the green marker after a reader green publication. TestReaderLandingLeavesTheGreenMarker must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ngs-t06-author-20261025",
              "digest": "sha256:739f5a10a0d5d671b44724d4738bbb6b9b01d28b781acd92995415917531aac0",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/land_marker.go,swap,failed,1,TestReaderLandingLeavesTheGreenMarker,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestReaderLandingLeavesTheGreenMarker,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,450\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,381,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestReaderLandingLeavesTheGreenMarker,\"land_reader_test.go:38: refs/bench/green/main = dc07699aa03d0d9ad9ada6fcffdbaa21b39f92be after a reader green publication, want the prior marker 638ef50fc7f1b1086dd4966e7f774ae3086f49c9\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t04-ticket-4-check-1",
          "performer": "claude:ngs-t04-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f2625783b7ba638fc6f8ce51fca224e880b1385",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t04-author-20261025",
            "digest": "sha256:c76074f9125f4321e00cbd7eddd4b186bb7209ec8bd6de6aa55ee46f4d18180f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,78877,637\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c1c527bf7fec4cfb89961c98806d204a9d86cf36",
          "requirement": "ticket-4-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t04-ticket-4-check-2",
          "performer": "claude:ngs-t04-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f2625783b7ba638fc6f8ce51fca224e880b1385",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t04-author-20261025",
            "digest": "sha256:3c1b2268e0e9e5e4239eda85c4e7bfb9dc4c51c889b5e4ae246d084eb8cb91d2",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,1327,179\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c1c527bf7fec4cfb89961c98806d204a9d86cf36",
          "requirement": "ticket-4-check-2",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "t04-ticket-4-check-3",
          "performer": "claude:ngs-t04-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f2625783b7ba638fc6f8ce51fca224e880b1385",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t04-author-20261025",
            "digest": "sha256:3c1aad2c0d3b1caa931b7015f3952241a47b836005c2438acfbfb6428348351f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15361,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c1c527bf7fec4cfb89961c98806d204a9d86cf36",
          "requirement": "ticket-4-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t04-ticket-4-probe",
          "performer": "claude:ngs-t04-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f2625783b7ba638fc6f8ce51fca224e880b1385",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t04-author-20261025",
            "digest": "sha256:bf195e83b82b121b506fae6e90a1108388a21164339430756dbf17f60f9320cb",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/delta_checkpoint.go,swap,failed,1,TestReaderRedWritesNoRecord,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,^TestReaderRedWritesNoRecord$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,333\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,68,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestReaderRedWritesNoRecord,\"reader_set_test.go:130: ReaderRecorded after a red reader run = true, <nil>; want false\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c1c527bf7fec4cfb89961c98806d204a9d86cf36",
          "requirement": "ticket-4-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "In the reader run, write the reader record before the status check. TestReaderRedWritesNoRecord must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ngs-t04-author-20261025",
              "digest": "sha256:bf195e83b82b121b506fae6e90a1108388a21164339430756dbf17f60f9320cb",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/delta_checkpoint.go,swap,failed,1,TestReaderRedWritesNoRecord,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,^TestReaderRedWritesNoRecord$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,333\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,68,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestReaderRedWritesNoRecord,\"reader_set_test.go:130: ReaderRecorded after a red reader run = true, <nil>; want false\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05-ticket-5-check-1",
          "performer": "claude:ngs-t05-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "06f68fe87b09c4d66e6b87c83fb520eb3744b041",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t05-author-20261025",
            "digest": "sha256:dd67a1fc27b0f12262870d01fb65183436f0c36b0e40c0fedf624315089fa541",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,79529,660\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d2837921053f74ff6f4aaec695b7a9dfe83b8abc",
          "requirement": "ticket-5-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-check-2",
          "performer": "claude:ngs-t05-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "06f68fe87b09c4d66e6b87c83fb520eb3744b041",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t05-author-20261025",
            "digest": "sha256:f7c26dc447966949bfa1d11a4292e84660f649fe244ec257679edc824255f099",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,1345,180\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d2837921053f74ff6f4aaec695b7a9dfe83b8abc",
          "requirement": "ticket-5-check-2",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-check-3",
          "performer": "claude:ngs-t05-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "06f68fe87b09c4d66e6b87c83fb520eb3744b041",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t05-author-20261025",
            "digest": "sha256:3d2438dc9aeb28146b475021acdf1b4213ae2fab77c7cd009c32673dcb1e3b2d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/packagesurface,pass,3,14\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d2837921053f74ff6f4aaec695b7a9dfe83b8abc",
          "requirement": "ticket-5-check-3",
          "command": "bench test --package ./internal/packagesurface",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-check-4",
          "performer": "claude:ngs-t05-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "06f68fe87b09c4d66e6b87c83fb520eb3744b041",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t05-author-20261025",
            "digest": "sha256:eb25afe99f76ad4f39c0be285268416618722ac8a81e5af7bba2c0c1a8abb1b0",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13933,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d2837921053f74ff6f4aaec695b7a9dfe83b8abc",
          "requirement": "ticket-5-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-probe",
          "performer": "claude:ngs-t05-author-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "06f68fe87b09c4d66e6b87c83fb520eb3744b041",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t05-author-20261025",
            "digest": "sha256:4b7deebf4d9c67346619f6c943936cfcd0e1a53dec189bb08318d045d63d7914",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/reader_plan.go,swap,failed,1,TestReaderPlanRunsWholeAfterFiveReaderLandings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestReaderPlanRunsWholeAfterFiveReaderLandings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,369\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,71,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestReaderPlanRunsWholeAfterFiveReaderLandings,\"reader_plan_test.go:326: plan = true \\\"\\\"; want the whole gate with the refusal \\\"reader landing bound\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d2837921053f74ff6f4aaec695b7a9dfe83b8abc",
          "requirement": "ticket-5-probe",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "In the reader plan, change the bound of five first-parent commits to six. TestReaderPlanRunsWholeAfterFiveReaderLandings must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ngs-t05-author-20261025",
              "digest": "sha256:4b7deebf4d9c67346619f6c943936cfcd0e1a53dec189bb08318d045d63d7914",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/reader_plan.go,swap,failed,1,TestReaderPlanRunsWholeAfterFiveReaderLandings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestReaderPlanRunsWholeAfterFiveReaderLandings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,369\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,71,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestReaderPlanRunsWholeAfterFiveReaderLandings,\"reader_plan_test.go:326: plan = true \\\"\\\"; want the whole gate with the refusal \\\"reader landing bound\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06r1-ticket-6-check-1",
          "performer": "claude:ngs-t06-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b78f2b2b4025079aff7990d9bde89be284363",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair1-20261025",
            "digest": "sha256:c3794b9543f4f8492c3c7d9a6e72a9e0434c8af95a3246a9ad016459dbd7debc",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,224,16\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "92c79f1af4edea5a92f895027c16f242a7ffc81c",
          "requirement": "ticket-6-check-1",
          "command": "bench test --package ./internal/gate/authorization",
          "exit_code": 0
        },
        {
          "id": "t06r1-ticket-6-check-2",
          "performer": "claude:ngs-t06-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b78f2b2b4025079aff7990d9bde89be284363",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair1-20261025",
            "digest": "sha256:a1acb68d0504fcb87906c42926671c10d6bbf5ba0ef723947ff579a5ea6f1e25",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,9521,222\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "source_commit": "92c79f1af4edea5a92f895027c16f242a7ffc81c",
          "requirement": "ticket-6-check-2",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "t06r1-ticket-6-check-3",
          "performer": "claude:ngs-t06-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b78f2b2b4025079aff7990d9bde89be284363",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair1-20261025",
            "digest": "sha256:d1e849a4fe1c14b352da2a66b3318638e9e10c05a51180a137721a65f12f25c1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,73845,1516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "92c79f1af4edea5a92f895027c16f242a7ffc81c",
          "requirement": "ticket-6-check-3",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t06r1-ticket-6-check-4",
          "performer": "claude:ngs-t06-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b78f2b2b4025079aff7990d9bde89be284363",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair1-20261025",
            "digest": "sha256:b234765d58cc2fdc4b3575946c0965d9cd542b38c2698c49f7018cb6da16a38a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree/landingpolicy,pass,2,46\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "92c79f1af4edea5a92f895027c16f242a7ffc81c",
          "requirement": "ticket-6-check-4",
          "command": "bench test --package ./internal/worktree/landingpolicy",
          "exit_code": 0
        },
        {
          "id": "t06r1-ticket-6-check-5",
          "performer": "claude:ngs-t06-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b78f2b2b4025079aff7990d9bde89be284363",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair1-20261025",
            "digest": "sha256:972b61629cb8fc51fac898e607196d3b2bfb78d94a12ac707fafbfc88253196e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,22100,406\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "92c79f1af4edea5a92f895027c16f242a7ffc81c",
          "requirement": "ticket-6-check-5",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t06r1-ticket-6-check-6",
          "performer": "claude:ngs-t06-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b78f2b2b4025079aff7990d9bde89be284363",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair1-20261025",
            "digest": "sha256:af041d3d102cb7b0b92d1ef6494245aac0b6acc2fd32cccffc1ed4a713b57153",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13631,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "92c79f1af4edea5a92f895027c16f242a7ffc81c",
          "requirement": "ticket-6-check-6",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06r1-ticket-6-probe",
          "performer": "claude:ngs-t06-repair1-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b78f2b2b4025079aff7990d9bde89be284363",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair1-20261025",
            "digest": "sha256:ef20120b8a4cf3743fad863c9842d6816591d8319d7d552d7e03308b700ebc92",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/land_marker.go,swap,failed,1,TestReaderLandingLeavesTheGreenMarker,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,^TestReaderLandingLeavesTheGreenMarker$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,450\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,381,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestReaderLandingLeavesTheGreenMarker,\"land_reader_test.go:38: refs/bench/green/main = e1bce28bee4c218919950fd2de46a8913e27933c after a reader green publication, want the prior marker f8cd2b4f8a7fe957fcbd04d07983a5c71f24c24f\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "92c79f1af4edea5a92f895027c16f242a7ffc81c",
          "requirement": "ticket-6-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "In the reviewed landing, advance the green marker after a reader green publication. TestReaderLandingLeavesTheGreenMarker must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ngs-t06-repair1-20261025",
              "digest": "sha256:ef20120b8a4cf3743fad863c9842d6816591d8319d7d552d7e03308b700ebc92",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/land_marker.go,swap,failed,1,TestReaderLandingLeavesTheGreenMarker,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,^TestReaderLandingLeavesTheGreenMarker$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,450\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,381,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestReaderLandingLeavesTheGreenMarker,\"land_reader_test.go:38: refs/bench/green/main = e1bce28bee4c218919950fd2de46a8913e27933c after a reader green publication, want the prior marker f8cd2b4f8a7fe957fcbd04d07983a5c71f24c24f\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06r2-ticket-6-check-1",
          "performer": "claude:ngs-t06-repair2-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair2-20261025",
            "digest": "sha256:fabb72fe6984076ba8d922441b5c91ffc57c4f075caedc37ed1e275d20a4fddc",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,227,16\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "7204223529757f414366245ec4f8db01de7a3f65",
          "requirement": "ticket-6-check-1",
          "command": "bench test --package ./internal/gate/authorization",
          "exit_code": 0
        },
        {
          "id": "t06r2-ticket-6-check-2",
          "performer": "claude:ngs-t06-repair2-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair2-20261025",
            "digest": "sha256:4a20dd9ec48de5bdd897f9417339522fa0c684b486bd983e69d6378b5edb1d5e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,9669,222\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "source_commit": "7204223529757f414366245ec4f8db01de7a3f65",
          "requirement": "ticket-6-check-2",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "t06r2-ticket-6-check-3",
          "performer": "claude:ngs-t06-repair2-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair2-20261025",
            "digest": "sha256:33b5584b657353cd0e374107df42b7ff278d8ae10461ad11ed46a2c6c44ab205",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,74445,1516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "7204223529757f414366245ec4f8db01de7a3f65",
          "requirement": "ticket-6-check-3",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t06r2-ticket-6-check-4",
          "performer": "claude:ngs-t06-repair2-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair2-20261025",
            "digest": "sha256:ac6b7ad585e26bf0c4a31e2818a77ddfd72ee54a806b5b50f93afc023560977a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree/landingpolicy,pass,3,46\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "7204223529757f414366245ec4f8db01de7a3f65",
          "requirement": "ticket-6-check-4",
          "command": "bench test --package ./internal/worktree/landingpolicy",
          "exit_code": 0
        },
        {
          "id": "t06r2-ticket-6-check-5",
          "performer": "claude:ngs-t06-repair2-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair2-20261025",
            "digest": "sha256:9809fc66de8a5d8f46143d722d4d4d8a8a7d458572eac0a0385d8981d038141c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,23018,406\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "7204223529757f414366245ec4f8db01de7a3f65",
          "requirement": "ticket-6-check-5",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t06r2-ticket-6-check-6",
          "performer": "claude:ngs-t06-repair2-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair2-20261025",
            "digest": "sha256:89891d49ed704b64fce194d40d2b4e32b5479988cda8d60cfb863a81c1c4002c",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13526,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "7204223529757f414366245ec4f8db01de7a3f65",
          "requirement": "ticket-6-check-6",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06r2-ticket-6-probe",
          "performer": "claude:ngs-t06-repair2-20261025",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t06-repair2-20261025",
            "digest": "sha256:4f879d6f2484000a9b07c4a9e0586a93b1a3d96c289d29b25332baa540ebcd35",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/land_marker.go,swap,failed,1,TestReaderLandingLeavesTheGreenMarker,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestReaderLandingLeavesTheGreenMarker,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,450\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,381,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestReaderLandingLeavesTheGreenMarker,\"land_reader_test.go:38: refs/bench/green/main = 9ad609e9fbc4f7e81ae4cd68971184976662bb61 after a reader green publication, want the prior marker be7b63ea4c7d81745a58e57148bf90aeb367ec2b\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "7204223529757f414366245ec4f8db01de7a3f65",
          "requirement": "ticket-6-probe",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "In the reviewed landing, advance the green marker after a reader green publication. TestReaderLandingLeavesTheGreenMarker must fail, then pass after the restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ngs-t06-repair2-20261025",
              "digest": "sha256:4f879d6f2484000a9b07c4a9e0586a93b1a3d96c289d29b25332baa540ebcd35",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/land_marker.go,swap,failed,1,TestReaderLandingLeavesTheGreenMarker,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestReaderLandingLeavesTheGreenMarker,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,450\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,381,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestReaderLandingLeavesTheGreenMarker,\"land_reader_test.go:38: refs/bench/green/main = 9ad609e9fbc4f7e81ae4cd68971184976662bb61 after a reader green publication, want the prior marker be7b63ea4c7d81745a58e57148bf90aeb367ec2b\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c3r1-standards",
          "performer": "claude:ngs-c3-standards-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d55a1dd09e3a53f63ffcf69525cbb542990431a8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c3-standards-20261025",
            "digest": "sha256:a4bed550e5473051b93f15a06d4f6d4ba411a5c7f0db838c409f0c0f7b232ec0",
            "excerpt": "Axis: Standards, NG-C3, tickets 04, 05, 06. Session claude:ngs-c3-standards-20261025. Frozen pair f9e20cf3..ff617775.\nCount: 3 findings. Worst issue: S2.\nS1 sev 5, conf 9, auto-fix: the kit-run predicate is derived twice (reader_plan.go:77 and reader_store.go:55 beginReaderTrace). Extract one predicate.\nS2 sev 6, conf 9, ask-user: the per-file embed parse is duplicated (embed_directives.go:47-67 directivePatterns and assets.go:98-118 EmbedTargets) and already diverges on backticks and the all: prefix. EmbedTargets should call directivePatterns.\nS3 sev 5, conf 9, auto-fix: the link-mode check on both sides of a change repeats changeClasses (reader_plan.go:117-123 and lane_select.go:302-306). Extract one predicate.\nrunReduced (ticket 04) is a faithful extraction: no finding.\nAdvice (no ID): two mechanisms set Schedule; commitTree repeats an older rev-parse idiom; the refusal relocation is pure.\n"
          },
          "axis": "Standards",
          "base": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
          "tip": "ff617775e2ca93ad218f92d8fad7395a2cacad38",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "c3r1-spec",
          "performer": "claude:ngs-c3-spec-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d55a1dd09e3a53f63ffcf69525cbb542990431a8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c3-spec-20261025",
            "digest": "sha256:29a8ad990d11aa28a917488067c79d5aa36d13947e65662a160a4a3d6b6a45b4",
            "excerpt": "Axis: Spec, NG-C3, tickets 04, 05, 06, rows NG29-NG59, NG73-NG77, NG79, NG80. Session claude:ngs-c3-spec-20261025. Frozen pair f9e20cf3..ff617775.\nCount: 3 findings (axis IDs S1-S3, recorded as P1-P3). Worst issue: P1.\nP1 sev 5, conf 6, ask-user, ticket 06: attribute (authorization.go:33-41) maps a reader-set timeout (124) or cancel (130) to candidate from the exit codes alone; the whole path gives infrastructure. Define red as a failed phase table.\nP2 sev 5, conf 5, ask-user, ticket 06: publishedReaderRecorded (land_resume.go:664-671) treats a read error as no record, so a resumed reader landing can advance the green marker (spec lines 38-39 forbid it). Fail closed.\nP3 sev 5, conf 6, auto-fix, ticket 05: decide reads changes through ComposedChanges/TreeChanges without --ignore-submodules=none (tree.go:285-288); an ambient diff.ignoreSubmodules hides a gitlink change from case 6. Use TreeChangesIncludingSubmodules.\nAdvice (no ID): the duplicated embed parse; a destination-tree read failure is labelled no green marker; the resumeMarkerFacts wiring of the reader fact has no biting test; the reader-record store has no prune.\n"
          },
          "axis": "Spec",
          "base": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
          "tip": "ff617775e2ca93ad218f92d8fad7395a2cacad38",
          "finding_ids": [
            "P1",
            "P2",
            "P3"
          ],
          "supersedes": []
        },
        {
          "id": "c3r1-coverage",
          "performer": "claude:ngs-c3-coverage-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d55a1dd09e3a53f63ffcf69525cbb542990431a8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c3-coverage-20261025",
            "digest": "sha256:8edd27df22d12771f2cb19a62c67ed27a8d1393ee556e2a9587af5a8aa42ad2c",
            "excerpt": "Axis: Coverage, NG-C3, tickets 04, 05, 06. Session claude:ngs-c3-coverage-20261025. Frozen pair f9e20cf3..ff617775.\nCount: 7 findings. Worst issue: C5.\nC5 sev 6, conf 8, auto-fix, ticket 06: the resume wiring of PublishedReaderRecorded (land_resume.go:296-314) has no biting test; NG59 presets the fact.\nC1 sev 5, conf 6, ask-user, tickets 05/06: an open subject can land on reader green; the plan never reads plan.Closed and the table persists on any whole green.\nC2 sev 5, conf 7, auto-fix, ticket 06: a reader timeout or cancel maps to candidate (same as Spec P1).\nC3 sev 5, conf 8, auto-fix, ticket 05: the IsAncestor clause of case 3 (reader_plan.go:86) has no test.\nC4 sev 5, conf 8, auto-fix, ticket 05: the chain-gap test places its only gap at index 0.\nC6 sev 5, conf 6, auto-fix, ticket 05: TreeChanges hides a gitlink under ambient config (same as Spec P3).\nC7 sev 5, conf 7, auto-fix, ticket 05: strings.Fields splits a quoted embed pattern with a space; the loop copies EmbedTargets (same root as Standards S2).\nAdvice (no ID): no test swaps cases 10/11 or 5/6; the plan chain fixture is linear; the reader-record store has no bound.\n"
          },
          "axis": "Coverage",
          "base": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
          "tip": "ff617775e2ca93ad218f92d8fad7395a2cacad38",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
            "C5",
            "C6",
            "C7"
          ],
          "supersedes": []
        },
        {
          "id": "c3r2-standards",
          "performer": "claude:ngs-c3r2-standards-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "039b78f2b2b4025079aff7990d9bde89be284363",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c3r2-standards-20261025",
            "digest": "sha256:ce9acf2f98c19bd875968700d9afdea4c0dceb906afd955fe524d00244f34f44",
            "excerpt": "Axis: Standards, NG-C3 confirming round r2. Session claude:ngs-c3r2-standards-20261025. Repair delta ff617775..92c79f1a.\nCount: 1 finding. Worst issue: S4.\nS1, S2, S3 confirmed closed: kitRun (reader_store.go:54-56) is the one kit-run predicate; EmbedTargets calls directivePatterns and keeps its build-source scope; linkChange (lane_select.go:71) serves changeClasses and the reader plan.\nS4 sev 5, conf 6, auto-fix: phaseTableRed = 1 (authorization.go:38-39) restates the gate's red exit (report.go:82). The gate should export one named constant used at both sites.\nAdvice (no ID): the worktree fixture hand-writes the reader-record JSON; reader_plan.go:94 re-derives the destination tree.\n"
          },
          "axis": "Standards",
          "base": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
          "tip": "92c79f1af4edea5a92f895027c16f242a7ffc81c",
          "finding_ids": [
            "S4"
          ],
          "supersedes": [
            "c3r1-standards"
          ]
        },
        {
          "id": "c3r2-spec",
          "performer": "claude:ngs-c3r2-spec-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "039b78f2b2b4025079aff7990d9bde89be284363",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c3r2-spec-20261025",
            "digest": "sha256:58a91c0f6851037afa6ff8f95eec2a5fe9da9c83b28b8d8649c275770d41b943",
            "excerpt": "Axis: Spec, NG-C3 confirming round r2. Session claude:ngs-c3r2-spec-20261025. Repair delta ff617775..92c79f1a and plan amendment c9e9e8f5..8175d208.\nCount: 0 findings. Worst issue: none.\nP1, P2, P3, NG83, NG84, and the N1 order (derive before the release, persist after the verdict) are satisfied.\nAdvice (no ID): the resume refusal branch has no verb-level test; an Invalid reader record reads as no record (matches the spec's \"valid\"); the red exit literal is duplicated; NG83's why-cell names only the plan half; the lane's ComposedChanges still uses TreeChanges.\n"
          },
          "axis": "Spec",
          "base": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
          "tip": "92c79f1af4edea5a92f895027c16f242a7ffc81c",
          "finding_ids": [],
          "supersedes": [
            "c3r1-spec"
          ]
        },
        {
          "id": "c3r2-coverage",
          "performer": "claude:ngs-c3r2-coverage-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "039b78f2b2b4025079aff7990d9bde89be284363",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c3r2-coverage-20261025",
            "digest": "sha256:0e53c610ef23313cb81f84f51b0bb25ca7d852e191f57c0eaf68c712e88cd2ac",
            "excerpt": "Axis: Coverage, NG-C3 confirming round r2. Session claude:ngs-c3r2-coverage-20261025. Repair delta ff617775..92c79f1a.\nCount: 0 findings. Worst issue: none.\nEach new test is red-capable against its finding: C1/NG83 (plan and store), C3, C4, C5/P2, C6/P3, C7, P1/C2, NG84.\nAdvice (no ID): the verb-level resume refusal is unpinned (severity 4); EmbedTargets now strips all: and backquotes with no direct test.\n"
          },
          "axis": "Coverage",
          "base": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
          "tip": "92c79f1af4edea5a92f895027c16f242a7ffc81c",
          "finding_ids": [],
          "supersedes": [
            "c3r1-coverage"
          ]
        },
        {
          "id": "c3r3-standards",
          "performer": "claude:ngs-c3r3-standards-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c3r3-standards-20261025",
            "digest": "sha256:532b493775fe0fd77005478b5c35a8adf99b48424dc25fb0692a48499dc8401e",
            "excerpt": "Axis: Standards, NG-C3 confirming round r3. Session claude:ngs-c3r3-standards-20261025. Delta 92c79f1a..72042235.\nCount: 0 findings. S4 closed: gate.PhaseTableRedExit (report.go:37-38) is the one source, returned at report.go:85 and compared at authorization.go:176; phaseTableRed is gone; behavior unchanged.\n"
          },
          "axis": "Standards",
          "base": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
          "tip": "7204223529757f414366245ec4f8db01de7a3f65",
          "finding_ids": [],
          "supersedes": [
            "c3r2-standards"
          ]
        },
        {
          "id": "c3r3-spec",
          "performer": "claude:ngs-c3r3-spec-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c3r3-spec-20261025",
            "digest": "sha256:e032f5d25b82e91b417ec687b067c5be8a6b30c4b5d34d0061695451251bc3f6",
            "excerpt": "Axis: Spec, NG-C3 confirming round r3. Session claude:ngs-c3r3-spec-20261025. Delta 92c79f1a..72042235.\nCount: 0 findings. The rename changes no behavior of NG55 or NG56: a red reader set still maps to candidate, and 124 and 130 still map to infrastructure.\n"
          },
          "axis": "Spec",
          "base": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
          "tip": "7204223529757f414366245ec4f8db01de7a3f65",
          "finding_ids": [],
          "supersedes": [
            "c3r2-spec"
          ]
        },
        {
          "id": "c3r3-coverage",
          "performer": "claude:ngs-c3r3-coverage-20261025",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "19e419008a71273091c26b576cdc720e90b4cb0f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ngs-c3r3-coverage-20261025",
            "digest": "sha256:5ca20586ca129ec2f3751332b116b7b4f5600b9e5f6406770f1790a975e856f7",
            "excerpt": "Axis: Coverage, NG-C3 confirming round r3. Session claude:ngs-c3r3-coverage-20261025. Delta 92c79f1a..72042235.\nCount: 0 findings. reader_kind_test.go keeps the literal exit map {1: Candidate, 124: Infrastructure, 130: Infrastructure}, so a change of PhaseTableRedExit turns it red.\n"
          },
          "axis": "Coverage",
          "base": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
          "tip": "7204223529757f414366245ec4f8db01de7a3f65",
          "finding_ids": [],
          "supersedes": [
            "c3r2-coverage"
          ]
        }
      ]
    },
    {
      "id": "NG-C6",
      "base": "7204223529757f414366245ec4f8db01de7a3f65",
      "tip": "a977ad8b7436b68ab74ac76930746ae89280d409",
      "plan_digest": "sha256:d1077d01a3ec7991d40fabd7a7b56e387cb8b270c52105ccd278f5b050ffea11",
      "source_digest": "a66ef5ea86b51abd99ba800a940f49b82a165098",
      "acceptance_rows": [
        "NG68",
        "NG69",
        "NG70"
      ],
      "verification": [
        {
          "id": "t09-ticket-9-check-1",
          "performer": "claude:ngs-t09-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a66ef5ea86b51abd99ba800a940f49b82a165098",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t09-author-20261026",
            "digest": "sha256:ed928f34f1288fb0d387ed1b0e9cc58578ce113e40764b8aa83b62678e90436c",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,7305,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "a977ad8b7436b68ab74ac76930746ae89280d409",
          "requirement": "ticket-9-check-1",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t09-ticket-9-check-2",
          "performer": "claude:ngs-t09-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a66ef5ea86b51abd99ba800a940f49b82a165098",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t09-author-20261026",
            "digest": "sha256:2db154f2b4366e26a6effddda9090149ac1cf97e63d38a0076bcfbef6cb27b77",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,120,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "a977ad8b7436b68ab74ac76930746ae89280d409",
          "requirement": "ticket-9-check-2",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        },
        {
          "id": "t09-ticket-9-check-3",
          "performer": "claude:ngs-t09-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a66ef5ea86b51abd99ba800a940f49b82a165098",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ngs-t09-author-20261026",
            "digest": "sha256:c311dfbca31e591521ee7ab0cf3252a196cff97c42c0a49b388cb50a25b074fa",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18035,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "a977ad8b7436b68ab74ac76930746ae89280d409",
          "requirement": "ticket-9-check-3",
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
    },
    {
      "from": "sha256:fe8eec05cb96401f72af5a70e983c0f599e7a6c35d83366b8c15758f6f77e8be",
      "to": "sha256:68ef1b1ef031160d9eec83650abcea4fd38721acdb6654476a69446626e82429",
      "chunk_ids": {
        "NG-C1": [
          "NG-C1"
        ],
        "NG-C2": [
          "NG-C2"
        ]
      }
    },
    {
      "from": "sha256:68ef1b1ef031160d9eec83650abcea4fd38721acdb6654476a69446626e82429",
      "to": "sha256:ba62935cfa82c71017eacd4dd743bfaee0bcc4d3fbce6f11d4b3112cc8695d27",
      "chunk_ids": {
        "NG-C1": [
          "NG-C1"
        ],
        "NG-C2": [
          "NG-C2"
        ]
      }
    },
    {
      "from": "sha256:ba62935cfa82c71017eacd4dd743bfaee0bcc4d3fbce6f11d4b3112cc8695d27",
      "to": "sha256:6a41bd1bd06ebac944ca1dcf933eb7f3ec560f7ac75908afb40859eaae1674b7",
      "chunk_ids": {
        "NG-C1": [
          "NG-C1"
        ],
        "NG-C2": [
          "NG-C2"
        ],
        "NG-C3": [
          "NG-C3"
        ]
      }
    },
    {
      "from": "sha256:6a41bd1bd06ebac944ca1dcf933eb7f3ec560f7ac75908afb40859eaae1674b7",
      "to": "sha256:44a38079f9734bd6a1f106fe16ba748911098f8d5a19e6f5cbd42d7c7d55e6e4",
      "chunk_ids": {
        "NG-C1": [
          "NG-C1"
        ],
        "NG-C2": [
          "NG-C2"
        ],
        "NG-C3": [
          "NG-C3"
        ]
      }
    },
    {
      "from": "sha256:44a38079f9734bd6a1f106fe16ba748911098f8d5a19e6f5cbd42d7c7d55e6e4",
      "to": "sha256:d1077d01a3ec7991d40fabd7a7b56e387cb8b270c52105ccd278f5b050ffea11",
      "chunk_ids": {
        "NG-C1": [
          "NG-C1"
        ],
        "NG-C2": [
          "NG-C2"
        ],
        "NG-C3": [
          "NG-C3"
        ]
      }
    }
  ],
  "folds": [
    {
      "id": "F1",
      "after": "NG-C2",
      "from": "05bd4025a5545c73e22e5171e48823665aa0f473",
      "commit": "f9e20cf3569b172cb80e40d33f094f1ed50a32ba",
      "incoming": "f862eafc9f24ab18a6b75d2be37f1da686fb9ed3",
      "source_digest": "f7dcc62f827f4ec189f3e98de47be2116cfdf2ea",
      "resolution": [],
      "resolution_digest": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "test_packages": [],
      "reviews": []
    }
  ]
}
```

## Repair state

NG-C1 consumed repair cycle 1 of 2 after its initial review (round r1). No reviewer extension applies.

NG-C2 consumed repair cycle 1 of 2 after its initial review (round r1). No reviewer extension applies.

NG-C3 consumed repair cycle 1 of 2 after its initial review (round r1), and cycle 2 of 2 for the round r2 finding S4. One combined repair session on ticket 06 carries every accepted finding, because the findings share files across tickets 03, 05, and 06. No reviewer extension applies.

## Standards

NG-C1 round r1: 0 findings. Worst issue: none.

NG-C2 round r1: 2 findings. Worst issue: S1, a restated list of the run flags of the traced test phase.

- S1 | auto-fix | confidence 7 | `internal/gate/reader_store.go:35-38` | `testRunFlags` restates the run flags that the traced phase builder emits. Ticket 03 line 17 forbids a restated list.
- S2 | auto-fix | confidence 8 | `internal/git/tree.go:175-176` | The symbolic-link mode gets a second production copy beside `internal/gate/lane_select.go:110`.

NG-C3 round r1: 3 findings. Worst issue: S2, a duplicated embed parse that already diverges.

- S1 | auto-fix | confidence 9 | `internal/gate/reader_plan.go:77` | The kit-run predicate repeats `beginReaderTrace` in `internal/gate/reader_store.go:55`.
- S2 | auto-fix | confidence 9 | `internal/packagesurface/embed_directives.go:47-67` | The embed parse repeats `EmbedTargets` and differs on backticks and the `all:` prefix. The axis asked for the user; `EmbedTargets` keeps its file scope.
- S3 | auto-fix | confidence 9 | `internal/gate/reader_plan.go:117-123` | The link-mode check repeats `changeClasses` in `internal/gate/lane_select.go:302-306`.

NG-C3 round r2 (confirming): S1, S2, and S3 are closed. 1 new finding.

- S4 | auto-fix | confidence 6 | `internal/gate/authorization/authorization.go:38-39` | `phaseTableRed = 1` restates the red exit of `internal/gate/report.go:82`.

NG-C3 round r3 (confirming): all three axes pass with 0 findings. S4 is closed by `PhaseTableRedExit`.

## Spec

NG-C1 round r1: 0 findings. Worst issue: none.

NG-C2 round r1: 0 findings. Worst issue: none.

NG-C2 round r2 (confirming): all three axes pass with 0 findings. Each repair test is red-capable against the finding it closes.

NG-C3 round r1: 3 findings. Worst issue: P1, a reader-set timeout that attributes as a candidate red.

- P1 | auto-fix | confidence 6 | `internal/gate/authorization/authorization.go:33-41` | A reader-set timeout or cancel maps to `candidate`. The axis asked for the user; a timeout and a cancel now map to `infrastructure`.
- P2 | auto-fix | confidence 5 | `internal/worktree/land_resume.go:664-671` | An unreadable reader record at resume reads as no record, so the marker can advance. The axis asked for the user; the resume now refuses.
- P3 | auto-fix | confidence 6 | `internal/gate/reader_plan.go:93` | The change list omits submodules, so ambient configuration can hide a gitlink change from case 6.

NG-C3 round r2 (confirming): 0 findings. P1, P2, P3, NG83, NG84, and the N1 order hold.

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

NG-C3 round r1: 7 findings. Worst issue: C5, the resume wiring of the reader fact has no biting test.

- C1 | auto-fix | confidence 6 | `internal/gate/reader_plan.go:76-152` | An open subject can land on `reader green`, and a table persists on an open whole green. The axis asked for the user; the plan now refuses an open subject, and rows NG83 and NG84 cover the change.
- C2 | auto-fix | confidence 7 | `internal/gate/authorization/authorization.go:39` | A timeout or a cancel maps to `candidate`. Spec finding P1 holds the same repair.
- C3 | auto-fix | confidence 8 | `internal/gate/reader_plan.go:86` | No test reaches the ancestor clause of case 3.
- C4 | auto-fix | confidence 8 | `internal/gate/reader_plan_test.go:341` | The chain-gap test places its only gap at the first commit.
- C5 | auto-fix | confidence 8 | `internal/worktree/land_resume.go:296-314` | No test bites the resume wiring of the reader fact.
- C6 | auto-fix | confidence 6 | `internal/gate/reader_plan.go:93` | Ambient configuration can hide a gitlink change. Spec finding P3 holds the same repair.
- C7 | auto-fix | confidence 7 | `internal/packagesurface/embed_directives.go:55` | A quoted embed pattern with a space splits into two words. Standards finding S2 holds the shared parser.

NG-C3 round r2 (confirming): 0 findings. Each new test is red-capable against the finding it closes.

## Measurement (NG60)

The ticket 06 author landed two `ROADMAP.md` changes in a disposable clone of `ff617775e`, with a separate Bench home. The first landing ran the whole gate in 275 seconds, green, with the case-3 refusal. The second landing also ran the whole gate, in 315 seconds, with the refusal `no reader table`. Both runs logged that `go list` failed in a deleted prospective checkout.

The orchestrator files this result as finding N1 on ticket 03's code. `settleRun` releases the prospective checkout before `trace.persist` lists the test packages, so no landing stores a reader table. The NG60 premise stays unproven until the repair lands and the measurement runs again.

The orchestrator ran the measurement again in a disposable clone of `e4cd9398d`, which holds both NG-C3 repairs. Each landing added one comment below the title of `ROADMAP.md`. The first landing ran the whole gate in 354 seconds, green. The second landing printed `gate: green (reader set)` in 173 seconds. The reader set took less than half the time of the whole gate, so the premise holds and NG60 is met.

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
