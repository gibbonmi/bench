# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/no-go-schedule/spec.md",
  "plan_digest": "sha256:11ed1fc05e800de44f46a007e13f3d2a5f224fab1ec958ed43cdcaf5fd8708f8",
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
    }
  ]
}
```

## Repair state

NG-C1 consumed repair cycle 1 of 2 after its initial review (round r1). No reviewer extension applies.

## Standards

NG-C1 round r1: 0 findings. Worst issue: none.

## Spec

NG-C1 round r1: 0 findings. Worst issue: none.

## Coverage

NG-C1 round r1: 1 finding. Worst issue: C1, no test varies an identity input other than the tree, the posture, and the tool bytes.

- C1 | auto-fix | confidence 6 | `internal/gate/oracle_identity_test.go` | A writer swap at `internal/gate/subject.go:77-105` drops an identity input and leaves NG1-NG3 green. Ticket 01 line 9 gives the identity every oracle frame except the tree.

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
