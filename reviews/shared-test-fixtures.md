# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/shared-test-fixtures/spec.md",
  "plan_digest": "sha256:1a14401f8264a19fe71ffd7de3b0ee308d915f2540baaa80cdb2e8d2ac381d8c",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GF-C1",
      "base": "4138ef1aa0c219c40396a9b799eac8769207b1f8",
      "tip": "6002ec0cfad73163762ef34044426b4cfb5b3298",
      "plan_digest": "sha256:15c922e92fbdf0b1d287370928c19669972477750af95e1a5a1e000f2e67a7b2",
      "source_digest": "57762fa392b4a1a6977422e6139578f2372ad26a",
      "acceptance_rows": [
        "GF1",
        "GF2",
        "GF3",
        "GF4",
        "GF5",
        "GF6",
        "GF7",
        "GF8",
        "GF9",
        "GF11",
        "GF12",
        "GF13",
        "GF14",
        "GF15",
        "GF20",
        "GF32",
        "GF33",
        "GF34",
        "GF35",
        "GF36",
        "GF37"
      ],
      "verification": [
        {
          "id": "t01-probe-01-run",
          "performer": "claude:gf-t01-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "46b8be48be138daa2cfc6e7648e239a442c6d0fd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t01-author-20261029",
            "digest": "sha256:50a97210eb09d66a7d42cb38d2c169b7485a6246191d1201684770c1f759102f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/gittest.go,swap,failed,1,TestRunFailureReturn,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestRunFailureReturn,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,371\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,8,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestRunFailureReturn,\"commands_test.go:134: Run(absent) = \\\"fatal: Needed a single revision\\\", <nil>; want empty output and the wrapped process error\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "63e498366d28674bb6a390f903391e8019bcea5f",
          "requirement": "probe-01-run",
          "command": "bench probe internal/gittest/gittest.go --swap 'return \"\", fmt.Errorf(\"git %s: %w: %s\", strings.Join(args, \" \"), err, out)' --with 'return strings.TrimSpace(string(out)), nil' --package ./internal/gittest --run TestRunFailureReturn",
          "exit_code": 0,
          "probe": {
            "mutation": "Return a failed command as success. The goroutine then receives no error, so the returned-failure test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t01-author-20261029",
              "digest": "sha256:50a97210eb09d66a7d42cb38d2c169b7485a6246191d1201684770c1f759102f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/gittest.go,swap,failed,1,TestRunFailureReturn,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestRunFailureReturn,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,371\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,8,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestRunFailureReturn,\"commands_test.go:134: Run(absent) = \\\"fatal: Needed a single revision\\\", <nil>; want empty output and the wrapped process error\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01-probe-01-identity",
          "performer": "claude:gf-t01-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "46b8be48be138daa2cfc6e7648e239a442c6d0fd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t01-author-20261029",
            "digest": "sha256:747f858e2d905d9ac0d82003474513091eba3207bebfc722d59eee5804dabeec",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFixtureIdentity,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestFixtureIdentity,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,337\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,21,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestFixtureIdentity,\"commands_test.go:301: git commit --quiet -m fixture: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "63e498366d28674bb6a390f903391e8019bcea5f",
          "requirement": "probe-01-identity",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/gittest --run TestFixtureIdentity",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t01-author-20261029",
              "digest": "sha256:747f858e2d905d9ac0d82003474513091eba3207bebfc722d59eee5804dabeec",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFixtureIdentity,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestFixtureIdentity,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,337\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,21,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestFixtureIdentity,\"commands_test.go:301: git commit --quiet -m fixture: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01r1-probe-01-run",
          "performer": "claude:gf-t01-repair1-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "57762fa392b4a1a6977422e6139578f2372ad26a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t01-repair1-20261029",
            "digest": "sha256:1980561e3d7c9e9775a239ee6e9bd7691934c89134720298a50b137bff98a22b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/gittest.go,swap,failed,1,TestRunFailureReturn,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestRunFailureReturn,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,372\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,16,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestRunFailureReturn,\"commands_test.go:134: Run(absent) = \\\"fatal: Needed a single revision\\\", <nil>; want empty output and the wrapped process error\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "6002ec0cfad73163762ef34044426b4cfb5b3298",
          "requirement": "probe-01-run",
          "command": "bench probe internal/gittest/gittest.go --swap 'return \"\", fmt.Errorf(\"git %s: %w: %s\", strings.Join(args, \" \"), err, out)' --with 'return strings.TrimSpace(string(out)), nil' --package ./internal/gittest --run TestRunFailureReturn",
          "exit_code": 0,
          "probe": {
            "mutation": "Return a failed command as success. The goroutine then receives no error, so the returned-failure test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t01-repair1-20261029",
              "digest": "sha256:1980561e3d7c9e9775a239ee6e9bd7691934c89134720298a50b137bff98a22b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/gittest.go,swap,failed,1,TestRunFailureReturn,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestRunFailureReturn,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,372\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,16,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestRunFailureReturn,\"commands_test.go:134: Run(absent) = \\\"fatal: Needed a single revision\\\", <nil>; want empty output and the wrapped process error\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01r1-probe-01-identity",
          "performer": "claude:gf-t01-repair1-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "57762fa392b4a1a6977422e6139578f2372ad26a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t01-repair1-20261029",
            "digest": "sha256:3addf57d8b02203f8746d8783645421bcbd0c934e0335b4a21ce21adf8c9073c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFixtureIdentity,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestFixtureIdentity,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,337\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,22,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestFixtureIdentity,\"commands_test.go:308: git commit --quiet -m fixture: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "6002ec0cfad73163762ef34044426b4cfb5b3298",
          "requirement": "probe-01-identity",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/gittest --run TestFixtureIdentity",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t01-repair1-20261029",
              "digest": "sha256:3addf57d8b02203f8746d8783645421bcbd0c934e0335b4a21ce21adf8c9073c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFixtureIdentity,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestFixtureIdentity,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,337\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,22,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestFixtureIdentity,\"commands_test.go:308: git commit --quiet -m fixture: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gf-c1-standards-r1",
          "performer": "claude:gf-c1-standards-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "46b8be48be138daa2cfc6e7648e239a442c6d0fd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c1-standards-r1",
            "digest": "sha256:c190c618ad15a7041e43007c9e7f7ead6bc8a61ee3a7c0107416b87708894aa0",
            "excerpt": "Standards axis, chunk GF-C1 (claude:gf-c1-standards-r1). Result: no findings. Identity single source at internal/gittest/identity.go:7-8; command at commands.go:14 is the one git process; testrepo snapshot exports removed; no dependency added. Advice: shared failure formatter (commands.go:34, gittest.go:321); init-plus-IdentityConfig loop repeat (snapshot.go:52, gittest.go:177); commit -qm through MustRun where Commit exists; working_tree.go holds only the package doc.\n"
          },
          "axis": "Standards",
          "base": "4138ef1aa0c219c40396a9b799eac8769207b1f8",
          "tip": "63e498366d28674bb6a390f903391e8019bcea5f",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "gf-c1-spec-r1",
          "performer": "claude:gf-c1-spec-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "46b8be48be138daa2cfc6e7648e239a442c6d0fd",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c1-spec-r1",
            "digest": "sha256:69943a751284e83f3ac04f5f3692e4f3d4f3eacec1245e3df6b202a1d4371408",
            "excerpt": "Spec axis, chunk GF-C1 (claude:gf-c1-spec-r1). Result: no binding finding; rows GF1-GF9, GF11-GF15, GF20, GF32-GF37 each have an implementation and an observing test. F1 commands.go:91-101 OutputBytes runs its own Output call beside Run; behavior matches GF3; no-op, confidence 3. F2 internal/testrepo/working_tree.go:1-3 holds only the package doc; no-op, confidence 2. Advice: recordcmd/ticket_test.go:121 restates bench@local for a later ticket.\n"
          },
          "axis": "Spec",
          "base": "4138ef1aa0c219c40396a9b799eac8769207b1f8",
          "tip": "63e498366d28674bb6a390f903391e8019bcea5f",
          "finding_ids": [
            "F1",
            "F2"
          ],
          "supersedes": []
        },
        {
          "id": "gf-c1-coverage-r1",
          "performer": "claude:gf-c1-coverage-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "46b8be48be138daa2cfc6e7648e239a442c6d0fd",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c1-coverage-r1",
            "digest": "sha256:44c3443d5bbb2a807230e88f67e5dd61567395ac0641362b8dba9cbe7e252234",
            "excerpt": "Coverage axis, chunk GF-C1 (claude:gf-c1-coverage-r1). C1 commands.go:28-35 no test makes OutputBytes fail on nonzero exit or missing Git; auto-fix, confidence 8. C2 identity.go:20-22 nothing grades IdentityArgs content; auto-fix, confidence 7. C3 gittest_test.go:12-21 and snapshot.go:44-47 no empty tracked file for the empty-versus-absent edge; auto-fix, confidence 7. C4 snapshot_test.go:40-50 GF15 grades only the commit failure; auto-fix or no-op, confidence 5. Advice: no newline-in-path case in TestCommandArgv.\n"
          },
          "axis": "Coverage",
          "base": "4138ef1aa0c219c40396a9b799eac8769207b1f8",
          "tip": "63e498366d28674bb6a390f903391e8019bcea5f",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": []
        }
      ],
      "confirmations": [
        {
          "id": "gf-c1-confirm-r1",
          "performer": "claude:gf-orchestrator-20261029",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "high",
          "source_digest": "57762fa392b4a1a6977422e6139578f2372ad26a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:gf-orchestrator-20261029",
            "digest": "sha256:854aca0f6659235236fe9e4f05dc2c66ef86052d3997572adcec8dd638d13cb7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFixtureIdentity,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestFixtureIdentity,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,443\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,43,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestFixtureIdentity,\"commands_test.go:315: IdentityArgs identity = \\\"bench test <other@example.invalid>|bench test <other@example.invalid>\\\"; want \\\"bench test <bench@example.invalid>|bench test <bench@example.invalid>\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "base": "4138ef1aa0c219c40396a9b799eac8769207b1f8",
          "reviewed_tip": "63e498366d28674bb6a390f903391e8019bcea5f",
          "repaired_tip": "6002ec0cfad73163762ef34044426b4cfb5b3298",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
            "F1",
            "F2"
          ],
          "supersedes": [],
          "probe": {
            "mutation": "swap the IdentityArgs email option to other@example.invalid",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:gf-orchestrator-20261029",
              "digest": "sha256:854aca0f6659235236fe9e4f05dc2c66ef86052d3997572adcec8dd638d13cb7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFixtureIdentity,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gittest,TestFixtureIdentity,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,443\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,fail,43,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gittest,TestFixtureIdentity,\"commands_test.go:315: IdentityArgs identity = \\\"bench test <other@example.invalid>|bench test <other@example.invalid>\\\"; want \\\"bench test <bench@example.invalid>|bench test <bench@example.invalid>\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ]
    },
    {
      "id": "GF-C2",
      "base": "6002ec0cfad73163762ef34044426b4cfb5b3298",
      "tip": "874d872c35a274eeeca55f7831508b46fb634b2b",
      "plan_digest": "sha256:ddaff991fc0933258b524e9d31c8ee54634c54b909d173458ecbf743ad44bb67",
      "source_digest": "93ef6b5139b4f395e2afdd49f2d053ba28ab3bf7",
      "acceptance_rows": [
        "GF16",
        "GF29",
        "GF30"
      ],
      "verification": [
        {
          "id": "t02-probe-02",
          "performer": "claude:gf-t02-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b83c32ffed193740ace21dd04a9c5657fd3cc09c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gf-t02-author-20261029",
            "digest": "sha256:eba9fb0edf5afadb180933cea13b41e4888329ae96351380619e51b6f26da098",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestReadStagedIndexReportsModesAndOmitsDeletions,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestReadStagedIndexReportsModesAndOmitsDeletions,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,352\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,15,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestReadStagedIndexReportsModesAndOmitsDeletions,\"staged_test.go:71: git commit --quiet -m base: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "0b272ed6f12cb704c69cc9c2026fc9d608680766",
          "requirement": "probe-02",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/git --run TestReadStagedIndexReportsModesAndOmitsDeletions",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gf-t02-author-20261029",
              "digest": "sha256:eba9fb0edf5afadb180933cea13b41e4888329ae96351380619e51b6f26da098",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestReadStagedIndexReportsModesAndOmitsDeletions,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestReadStagedIndexReportsModesAndOmitsDeletions,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,352\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,15,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestReadStagedIndexReportsModesAndOmitsDeletions,\"staged_test.go:71: git commit --quiet -m base: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t03-probe-03",
          "performer": "claude:gf-t03-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c755650ecc36f64f986c63224be4cb3aa3804231",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gf-t03-author-20261029",
            "digest": "sha256:c2194f7f886af2ad60e8f3308af534f25b0543c9648feb09513c6ab528b3d2b0",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestRetireDeletesTheFolderAndExitsZero,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,TestRetireDeletesTheFolderAndExitsZero,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,344\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,9,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestRetireDeletesTheFolderAndExitsZero,\"retire_test.go:176: git commit --quiet -m base: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d225e88cf9c2091efc2ebfaa62a7b29987642005",
          "requirement": "probe-03",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/spec --run TestRetireDeletesTheFolderAndExitsZero",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gf-t03-author-20261029",
              "digest": "sha256:c2194f7f886af2ad60e8f3308af534f25b0543c9648feb09513c6ab528b3d2b0",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestRetireDeletesTheFolderAndExitsZero,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,TestRetireDeletesTheFolderAndExitsZero,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,344\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,9,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestRetireDeletesTheFolderAndExitsZero,\"retire_test.go:176: git commit --quiet -m base: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t04-probe-04",
          "performer": "claude:gf-t04-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "80d2541f2df65975277aff4ba750e4c88108df35",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gf-t04-author-20261029",
            "digest": "sha256:2cf5e79aba42358fa99f72216ff5fe9269c00f45fd4940fd6f8d7773a7016057",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFlowCountsDetailFilesOverTheWindow,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmapflow,TestFlowCountsDetailFilesOverTheWindow,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,357\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmapflow,fail,13,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmapflow,TestFlowCountsDetailFilesOverTheWindow,\"flow_test.go:142: git commit --quiet -m base: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "a94afe984816d50c98dbd1a67ae4935fe863721e",
          "requirement": "probe-04",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/roadmapflow --run TestFlowCountsDetailFilesOverTheWindow",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gf-t04-author-20261029",
              "digest": "sha256:2cf5e79aba42358fa99f72216ff5fe9269c00f45fd4940fd6f8d7773a7016057",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFlowCountsDetailFilesOverTheWindow,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmapflow,TestFlowCountsDetailFilesOverTheWindow,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,357\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmapflow,fail,13,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmapflow,TestFlowCountsDetailFilesOverTheWindow,\"flow_test.go:142: git commit --quiet -m base: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t04r1-probe-04",
          "performer": "claude:gf-t04-repair1-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "93ef6b5139b4f395e2afdd49f2d053ba28ab3bf7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gf-t04-repair1-20261029",
            "digest": "sha256:1a4ae0c42f62b52df921b0c1d98cb26a84d5e1edd5c91d636824f8ac8cc9fc68",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFlowCountsDetailFilesOverTheWindow,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmapflow,TestFlowCountsDetailFilesOverTheWindow,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,357\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmapflow,fail,12,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmapflow,TestFlowCountsDetailFilesOverTheWindow,\"flow_test.go:142: git commit --quiet -m base: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "874d872c35a274eeeca55f7831508b46fb634b2b",
          "requirement": "probe-04",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/roadmapflow --run TestFlowCountsDetailFilesOverTheWindow",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gf-t04-repair1-20261029",
              "digest": "sha256:1a4ae0c42f62b52df921b0c1d98cb26a84d5e1edd5c91d636824f8ac8cc9fc68",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestFlowCountsDetailFilesOverTheWindow,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmapflow,TestFlowCountsDetailFilesOverTheWindow,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,357\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmapflow,fail,12,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmapflow,TestFlowCountsDetailFilesOverTheWindow,\"flow_test.go:142: git commit --quiet -m base: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gf-c2-standards-r1",
          "performer": "claude:gf-c2-standards-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "80d2541f2df65975277aff4ba750e4c88108df35",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c2-standards-r1",
            "digest": "sha256:60d752bf0add889a7a238b14fa65168ee9f83b9f8ce5264cb16b69b0779926c5",
            "excerpt": "Standards axis, chunk GF-C2 (claude:gf-c2-standards-r1). S1 internal/structure/structure_test.go:349-359 private run renames OutputBytes with its own TrimRight; auto-fix, confidence 7. S2 private dated exec runners in internal/diff/compatibility_test.go:150, internal/maps/freshness_test.go:53-60, internal/roadmap/retro_cost_test.go:32-48, and costRepo.git runs plain add, rev-parse, checkout; ask-user, confidence 5. S3 the IdentityConfig apply loop repeats in structure, maps, status fixtures; ask-user, confidence 4. S4 status_producible_test.go gitDir parameter rename; auto-fix, confidence 5. Advice: an IdentityArgs commit helper; a comment on the pinned compatibility identity.\n"
          },
          "axis": "Standards",
          "base": "6002ec0cfad73163762ef34044426b4cfb5b3298",
          "tip": "a94afe984816d50c98dbd1a67ae4935fe863721e",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4"
          ],
          "supersedes": []
        },
        {
          "id": "gf-c2-spec-r1",
          "performer": "claude:gf-c2-spec-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "80d2541f2df65975277aff4ba750e4c88108df35",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c2-spec-r1",
            "digest": "sha256:587d61ebc34aadb591058a597cf1eaf839fae6bd0bf01d599d69b6ee5c471c50",
            "excerpt": "Spec axis, chunk GF-C2 (claude:gf-c2-spec-r1). P1 (axis id F1) internal/status/status_signals_test.go:710-717 literal reformat and the roadmap_test.go import move are outside the fixture reach; auto-fix, confidence 6. P2 (axis id F2) internal/maps/freshness_test.go:56 commitAt is a timestamp-only private runner; siblings use t.Setenv plus gittest; ask-user, confidence 4. Kept raw sites classified correctly: diff command_test.go:69,78, matrix_test.go:277 and rawNoIndexPatch (expected exits), compatibility_test.go:125-137 (pinned identity and dates), state_scan_test.go:208 (stdin), retro_cost_test.go:40 (distinct author and committer times), structure run (output shape). Advice: shared identity apply helper; Commit in place of MustRun commit.\n"
          },
          "axis": "Spec",
          "base": "6002ec0cfad73163762ef34044426b4cfb5b3298",
          "tip": "a94afe984816d50c98dbd1a67ae4935fe863721e",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "gf-c2-coverage-r1",
          "performer": "claude:gf-c2-coverage-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "80d2541f2df65975277aff4ba750e4c88108df35",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c2-coverage-r1",
            "digest": "sha256:70742b54b7d716362edc3a0d2b797f5b16aa8e7536bedf1c8fd35fc2279fe4c4",
            "excerpt": "Coverage axis, chunk GF-C2 (claude:gf-c2-coverage-r1). V1 (axis id F1) internal/status/status_signals_test.go:264-267 literal reformat and the roadmap_test.go import move have no behavioral effect; no-op or revert, confidence 7. No changed subject found in the stdout-to-combined readers, the raw byte reads, the structure runner, the specialized probes, the stub-git tests, the identity change, or the empty -C root.\n"
          },
          "axis": "Coverage",
          "base": "6002ec0cfad73163762ef34044426b4cfb5b3298",
          "tip": "a94afe984816d50c98dbd1a67ae4935fe863721e",
          "finding_ids": [
            "V1"
          ],
          "supersedes": []
        }
      ],
      "confirmations": [
        {
          "id": "gf-c2-confirm-r1",
          "performer": "claude:gf-orchestrator-20261029",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "high",
          "source_digest": "93ef6b5139b4f395e2afdd49f2d053ba28ab3bf7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:gf-orchestrator-20261029",
            "digest": "sha256:8089e9e1a7032797161571ec9b3e14342ce881a5f0d6b0eec537acc4a256dcd5",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/maps/freshness_test.go,insert,failed,1,TestCommandProjectsTheStaleRowOnlyWhenTheAssetIsCommittedAfterTheIndex/asset_after_index,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/maps,TestCommandProjectsTheStaleRowOnlyWhenTheAssetIsCommittedAfterTheIndex,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,584\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/maps,fail,34,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/maps,TestCommandProjectsTheStaleRowOnlyWhenTheAssetIsCommittedAfterTheIndex/asset_after_index,\"freshness_test.go:79: Command(nil) = (exit 0, \\\"maps[1]{map,title,type,state,blockers,path}:\\\\\\\\n  alpha,Which parser owns the map?,Research,frontier,\\\\\\\\\\\"\\\\\\\\\\\",decisions/alpha/tickets/1.md\\\\\\\\nhelp[1]{cmd,why}:\\\\\\\\n  /bench-shape-idea,\\\\\\\\\\\"shape alpha: Which… (587 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "base": "6002ec0cfad73163762ef34044426b4cfb5b3298",
          "reviewed_tip": "a94afe984816d50c98dbd1a67ae4935fe863721e",
          "repaired_tip": "874d872c35a274eeeca55f7831508b46fb634b2b",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "P1",
            "P2",
            "V1"
          ],
          "supersedes": [],
          "probe": {
            "mutation": "insert a fixed committer date before the shared commit in commitAt",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:gf-orchestrator-20261029",
              "digest": "sha256:8089e9e1a7032797161571ec9b3e14342ce881a5f0d6b0eec537acc4a256dcd5",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/maps/freshness_test.go,insert,failed,1,TestCommandProjectsTheStaleRowOnlyWhenTheAssetIsCommittedAfterTheIndex/asset_after_index,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/maps,TestCommandProjectsTheStaleRowOnlyWhenTheAssetIsCommittedAfterTheIndex,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,584\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/maps,fail,34,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/maps,TestCommandProjectsTheStaleRowOnlyWhenTheAssetIsCommittedAfterTheIndex/asset_after_index,\"freshness_test.go:79: Command(nil) = (exit 0, \\\"maps[1]{map,title,type,state,blockers,path}:\\\\\\\\n  alpha,Which parser owns the map?,Research,frontier,\\\\\\\\\\\"\\\\\\\\\\\",decisions/alpha/tickets/1.md\\\\\\\\nhelp[1]{cmd,why}:\\\\\\\\n  /bench-shape-idea,\\\\\\\\\\\"shape alpha: Which… (587 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ]
    },
    {
      "id": "GF-C5A",
      "base": "874d872c35a274eeeca55f7831508b46fb634b2b",
      "tip": "cc197d6b3a004ed2c2b776573744f5b17c0cdd70",
      "plan_digest": "sha256:f53e416b2780300479d1e14e50b208cc60fe0103991607ebcf4320c09c6ed09f",
      "source_digest": "f9f1c8ab2431d8d681bd4016d5b14ccd862836ae",
      "acceptance_rows": [
        "GF16",
        "GF19",
        "GF29",
        "GF30"
      ],
      "verification": [
        {
          "id": "t05a-probe-05a",
          "performer": "claude:gf-t05a-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9ee4017b1ec4643f4a22d7c8f5d4ae6209635031",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gf-t05a-author-20261029",
            "digest": "sha256:883c475dd6916bb0495acf5b60e97ceaabf71df5d206c58a882ff959ff772a31",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestShellWrapperKitMatchesKitForRoot,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestShellWrapperKitMatchesKitForRoot,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,354\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,12,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot,\"wrapper_kit_test.go:111: git commit --quiet -m seed --allow-empty: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2d44df8b24f2998cc6f35feaf697deafcb172f86",
          "requirement": "probe-05a",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./cmd/bench --run TestShellWrapperKitMatchesKitForRoot",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gf-t05a-author-20261029",
              "digest": "sha256:883c475dd6916bb0495acf5b60e97ceaabf71df5d206c58a882ff959ff772a31",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestShellWrapperKitMatchesKitForRoot,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestShellWrapperKitMatchesKitForRoot,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,354\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,12,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot,\"wrapper_kit_test.go:111: git commit --quiet -m seed --allow-empty: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05a2-probe-05a",
          "performer": "claude:gf-t05a-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a5bf6ae22a750abdc790b354a808e777fa66e6f7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gf-t05a-author-20261029",
            "digest": "sha256:c71509ec0327b9f7c53732d17e3e0de82954789796d6988db0f35125f6ecf852",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestShellWrapperKitMatchesKitForRoot,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestShellWrapperKitMatchesKitForRoot,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,354\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,19,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot,\"wrapper_kit_test.go:111: git commit --quiet -m seed --allow-empty: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "6b30e5797d5056f5442ae6db3d0270e66a283af0",
          "requirement": "probe-05a",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./cmd/bench --run TestShellWrapperKitMatchesKitForRoot",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gf-t05a-author-20261029",
              "digest": "sha256:c71509ec0327b9f7c53732d17e3e0de82954789796d6988db0f35125f6ecf852",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestShellWrapperKitMatchesKitForRoot,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestShellWrapperKitMatchesKitForRoot,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,354\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,19,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot,\"wrapper_kit_test.go:111: git commit --quiet -m seed --allow-empty: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05b-probe-05b",
          "performer": "claude:gf-t05b-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "49e4928b94431a2a9d8478c16bdb95150f8f0541",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gf-t05b-author-20261029",
            "digest": "sha256:cf13d590b743c6bf278af9155aec17b935f6c32e992a526e427c5d502dd3374c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestCommitmentLightPathCommit,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commit,TestCommitmentLightPathCommit,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,349\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,fail,10,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commit,TestCommitmentLightPathCommit,\"commitment_test.go:219: git commit --quiet -m bootstrap: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "a8a6108eec28c50c398d57a68d5a8d4525d509f6",
          "requirement": "probe-05b",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/commit --run TestCommitmentLightPathCommit",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gf-t05b-author-20261029",
              "digest": "sha256:cf13d590b743c6bf278af9155aec17b935f6c32e992a526e427c5d502dd3374c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestCommitmentLightPathCommit,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commit,TestCommitmentLightPathCommit,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,349\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,fail,10,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commit,TestCommitmentLightPathCommit,\"commitment_test.go:219: git commit --quiet -m bootstrap: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05ar1-probe-05a",
          "performer": "claude:gf-t05a-repair1-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "380210fe111f693e2ebef2d299e2925d24e78607",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gf-t05a-repair1-20261029",
            "digest": "sha256:883c475dd6916bb0495acf5b60e97ceaabf71df5d206c58a882ff959ff772a31",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestShellWrapperKitMatchesKitForRoot,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestShellWrapperKitMatchesKitForRoot,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,354\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,12,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot,\"wrapper_kit_test.go:111: git commit --quiet -m seed --allow-empty: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "b1ea4021967b511f5d799eb1f422a700e634eff2",
          "requirement": "probe-05a",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./cmd/bench --run TestShellWrapperKitMatchesKitForRoot",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gf-t05a-repair1-20261029",
              "digest": "sha256:883c475dd6916bb0495acf5b60e97ceaabf71df5d206c58a882ff959ff772a31",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestShellWrapperKitMatchesKitForRoot,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestShellWrapperKitMatchesKitForRoot,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,354\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,12,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot,\"wrapper_kit_test.go:111: git commit --quiet -m seed --allow-empty: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05br1-probe-05b",
          "performer": "claude:gf-t05b-repair1-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f9f1c8ab2431d8d681bd4016d5b14ccd862836ae",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gf-t05b-repair1-20261029",
            "digest": "sha256:95c0038374c625bd0470c5851e7b918377736ff20229691851a3f9334e910a71",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestCommitmentLightPathCommit,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commit,TestCommitmentLightPathCommit,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,349\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,fail,13,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commit,TestCommitmentLightPathCommit,\"commitment_test.go:219: git commit --quiet -m bootstrap: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "cc197d6b3a004ed2c2b776573744f5b17c0cdd70",
          "requirement": "probe-05b",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/commit --run TestCommitmentLightPathCommit",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gf-t05b-repair1-20261029",
              "digest": "sha256:95c0038374c625bd0470c5851e7b918377736ff20229691851a3f9334e910a71",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestCommitmentLightPathCommit,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commit,TestCommitmentLightPathCommit,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,349\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,fail,13,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commit,TestCommitmentLightPathCommit,\"commitment_test.go:219: git commit --quiet -m bootstrap: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gf-c5a-standards-r1",
          "performer": "claude:gf-c5a-standards-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "49e4928b94431a2a9d8478c16bdb95150f8f0541",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c5a-standards-r1",
            "digest": "sha256:96a8b440583982d72fe7e9e624ec0d16cfde203a9268d67e5016e5bdfca0ee0a",
            "excerpt": "Standards axis, chunk GF-C5A (claude:gf-c5a-standards-r1). H1 (axis id F1) internal/adopt/link_hook_test.go:190 hookTestRepo is a one-line alias of gittest.Repo with 11 callers; auto-fix, confidence 6. Clean: runAdoptGit, runHookGit, outcomeGit, runGitIn removed; no literal identity in cmd/bench, internal/adopt, internal/commit, internal/gate tests; merge_parent_test.go:53 expected conflict kept. Advice: an initialize-in-directory gittest seam for gitInitRepo and initializeLandingRepo.\n"
          },
          "axis": "Standards",
          "base": "874d872c35a274eeeca55f7831508b46fb634b2b",
          "tip": "a8a6108eec28c50c398d57a68d5a8d4525d509f6",
          "finding_ids": [
            "H1"
          ],
          "supersedes": []
        },
        {
          "id": "gf-c5a-spec-r1",
          "performer": "claude:gf-c5a-spec-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "49e4928b94431a2a9d8478c16bdb95150f8f0541",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c5a-spec-r1",
            "digest": "sha256:050ccc7dcaba51d7d74f2dbf713aa5b342458340263df1e55e4305b1caac7778",
            "excerpt": "Spec axis, chunk GF-C5A (claude:gf-c5a-spec-r1). R1 (axis id S1) cmd/bench/command_registry_test.go:399 runAXIGit and preflight_version_test.go:146 trimmedAXIGit are generic fatal wrappers over production git.Raw that stay unmigrated; callers in command_registry_test.go, help_inventory_test.go, selected_queries_test.go, spec_retire_listing_test.go, preflight_version_test.go, tree_scope_test.go; auto-fix, confidence 7. R2 (axis id S2) cmd/bench/tree_scope_test.go:83-87 treeRowRepo commits a dated commit through runAXIGit; auto-fix, confidence 7. Conforming: internal/commit, internal/gate, internal/adopt classifications and retained sites. Advice: name the wrapper_kit_test.go worktree-add site in the pickup.\n"
          },
          "axis": "Spec",
          "base": "874d872c35a274eeeca55f7831508b46fb634b2b",
          "tip": "a8a6108eec28c50c398d57a68d5a8d4525d509f6",
          "finding_ids": [
            "R1",
            "R2"
          ],
          "supersedes": []
        },
        {
          "id": "gf-c5a-coverage-r1",
          "performer": "claude:gf-c5a-coverage-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "49e4928b94431a2a9d8478c16bdb95150f8f0541",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c5a-coverage-r1",
            "digest": "sha256:6fb7ae438ba4c76b55619d7ffa451216e04343d48cc15fe4b317696086c3cb42",
            "excerpt": "Coverage axis, chunk GF-C5A (claude:gf-c5a-coverage-r1). Findings: none. Stdout-only readers moved to combined output parse silent-stderr queries; headPaths kept raw bytes; merge --no-commit and the PATH stubs kept; no vacuous pass. Advice: the GF29 captures record only whether an identity is set.\n"
          },
          "axis": "Coverage",
          "base": "874d872c35a274eeeca55f7831508b46fb634b2b",
          "tip": "a8a6108eec28c50c398d57a68d5a8d4525d509f6",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "gf-c5a-standards-r2",
          "performer": "claude:gf-c5a-standards-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "380210fe111f693e2ebef2d299e2925d24e78607",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c5a-standards-r2",
            "digest": "sha256:8f0a4030641e58136269c0223a10ef176decb0cb4b7fbb56dd3507a5f0fef9d6",
            "excerpt": "Standards confirming round 2, chunk GF-C5A (claude:gf-c5a-standards-r2). Findings: none. R1 closed: runAXIGit and trimmedAXIGit removed, callers on gittest, branchTips raw through OutputBytes. R2 closed: treeRowRepo dates through t.Setenv, commit through gittest.Commit. The repair delta adds no runner, identity copy, or rename-only wrapper.\n"
          },
          "axis": "Standards",
          "base": "874d872c35a274eeeca55f7831508b46fb634b2b",
          "tip": "b1ea4021967b511f5d799eb1f422a700e634eff2",
          "finding_ids": [],
          "supersedes": [
            "gf-c5a-standards-r1"
          ]
        },
        {
          "id": "gf-c5a-spec-r2",
          "performer": "claude:gf-c5a-spec-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "380210fe111f693e2ebef2d299e2925d24e78607",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c5a-spec-r2",
            "digest": "sha256:18b701ec42663240a8ee36acc353b3997d67e48a93f78e48762d1df2a152ea8b",
            "excerpt": "Spec confirming round 2, chunk GF-C5A (claude:gf-c5a-spec-r2). Findings: none. R1 and R2 closed under Migration boundaries; every migrated call is fixture setup; the widen in 7708c6ff2 stays inside cmd/bench and changes no acceptance row.\n"
          },
          "axis": "Spec",
          "base": "874d872c35a274eeeca55f7831508b46fb634b2b",
          "tip": "b1ea4021967b511f5d799eb1f422a700e634eff2",
          "finding_ids": [],
          "supersedes": [
            "gf-c5a-spec-r1"
          ]
        },
        {
          "id": "gf-c5a-coverage-r2",
          "performer": "claude:gf-c5a-coverage-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "380210fe111f693e2ebef2d299e2925d24e78607",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c5a-coverage-r2",
            "digest": "sha256:763c8045d9a113ad0201c238c03ced10b864a88d00910cc919c277a5198230c6",
            "excerpt": "Coverage confirming round 2, chunk GF-C5A (claude:gf-c5a-coverage-r2). Findings: none. branchTips stays raw; the trimmed reads were trimmed before; no vacuous pass; the repair witness and the coordinator probe grade the inputs.\n"
          },
          "axis": "Coverage",
          "base": "874d872c35a274eeeca55f7831508b46fb634b2b",
          "tip": "b1ea4021967b511f5d799eb1f422a700e634eff2",
          "finding_ids": [],
          "supersedes": [
            "gf-c5a-coverage-r1"
          ]
        }
      ],
      "confirmations": [
        {
          "id": "gf-c5a-confirm-r2",
          "performer": "claude:gf-orchestrator-20261029",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "high",
          "source_digest": "f9f1c8ab2431d8d681bd4016d5b14ccd862836ae",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:gf-orchestrator-20261029",
            "digest": "sha256:e2d371f99ebb06814972eb562d5e2717817d7fbb9c6eb859b452986995c3a815",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact_test.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,ordinary-build-census,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,470\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,503,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: branch-native census forbids repository at internal/gate/prospectiveartifact/prospectiveartifact_test.go:306\",1\nskips[0]{package,test,reason}:\n"
          },
          "base": "874d872c35a274eeeca55f7831508b46fb634b2b",
          "reviewed_tip": "b1ea4021967b511f5d799eb1f422a700e634eff2",
          "repaired_tip": "cc197d6b3a004ed2c2b776573744f5b17c0cdd70",
          "finding_ids": [],
          "supersedes": [],
          "probe": {
            "mutation": "swap the raw worktree-list read back to gittest.Output; the census check turns red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:gf-orchestrator-20261029",
              "digest": "sha256:e2d371f99ebb06814972eb562d5e2717817d7fbb9c6eb859b452986995c3a815",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact_test.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,ordinary-build-census,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,470\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,503,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: branch-native census forbids repository at internal/gate/prospectiveartifact/prospectiveartifact_test.go:306\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ]
    },
    {
      "id": "GF-C5C",
      "base": "cc197d6b3a004ed2c2b776573744f5b17c0cdd70",
      "tip": "fd1234c51eb80f59ff191f3b069ee03f9992f765",
      "plan_digest": "sha256:8a878c0f6b72f4b296defc18ecec53c329e2fe7e2b0a98f6c4fd60e8ff9b011f",
      "source_digest": "a440eb60a9d37dee59b931229c6fa4b7682561a6",
      "acceptance_rows": [
        "GF16",
        "GF19",
        "GF29",
        "GF30"
      ],
      "verification": [
        {
          "id": "t05c-probe-05c",
          "performer": "claude:gf-t05c-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "795c28978221db70ba68d483183c76e2a2f1a963",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t05c-author-20261029",
            "digest": "sha256:75b6e74c120707ce46d90fddd0b366c1eb7d0fd1f407f5c231e7f57e2dead463",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,6,TestEvidenceRequiredSourceStates/CE41_absent,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,TestEvidenceRequiredSourceStates,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1309\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,47,10\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE41_absent,\"charge_evidence_test.go:37: git commit -q -m c0: exit status 128: Author identity unknown\",8\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE42_empty,\"charge_evidence_test.go:37: git commit -q -m c0: exit status 128: Author identity unknown\",8\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE43_symlink,\"charge_evidence_test.go:37: git commit -q -m c0: exit status 128: Author identity unknown\",8\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE47_directory,\"charge_evidence_test.go:37: git commit -q -m c0: exit status 128: Author identity unknown\",8\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE48_control_byte,\"charge_evidence_test.go:37: git commit -q -m c0: exit status 128: Author identity unknown\",8\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE48_invalid_UTF-8,\"charge_evidence_test.go:47: git commit -q -m c0: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ebc37cf4ce0a71cd74127eccd2348738e9049871",
          "requirement": "probe-05c",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/preflight --run TestEvidenceRequiredSourceStates",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t05c-author-20261029",
              "digest": "sha256:75b6e74c120707ce46d90fddd0b366c1eb7d0fd1f407f5c231e7f57e2dead463",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,6,TestEvidenceRequiredSourceStates/CE41_absent,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,TestEvidenceRequiredSourceStates,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1309\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,47,10\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE41_absent,\"charge_evidence_test.go:37: git commit -q -m c0: exit status 128: Author identity unknown\",8\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE42_empty,\"charge_evidence_test.go:37: git commit -q -m c0: exit status 128: Author identity unknown\",8\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE43_symlink,\"charge_evidence_test.go:37: git commit -q -m c0: exit status 128: Author identity unknown\",8\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE47_directory,\"charge_evidence_test.go:37: git commit -q -m c0: exit status 128: Author identity unknown\",8\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE48_control_byte,\"charge_evidence_test.go:37: git commit -q -m c0: exit status 128: Author identity unknown\",8\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceRequiredSourceStates/CE48_invalid_UTF-8,\"charge_evidence_test.go:47: git commit -q -m c0: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05d-probe-05d",
          "performer": "claude:gf-t05d-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c0ec29fb6b4017eaef40ca4ff2ed17d8ea5e569b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t05d-author-20261029",
            "digest": "sha256:f7c9ff6f3bd6c0bbc5c193e535135a2d6a8a8a875721abd3ea347fe836ba21b3",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestSnapshotProofLifecycle,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,TestSnapshotProofLifecycle,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,379\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,6,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestSnapshotProofLifecycle,\"intent_test.go:143: git -c user.name= -c user.email=bench@example.invalid commit -qm init: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "6a7ba7e9344ea58e69499d92f2ce5dad6823425d",
          "requirement": "probe-05d",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/intent --run TestSnapshotProofLifecycle",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t05d-author-20261029",
              "digest": "sha256:f7c9ff6f3bd6c0bbc5c193e535135a2d6a8a8a875721abd3ea347fe836ba21b3",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestSnapshotProofLifecycle,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,TestSnapshotProofLifecycle,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,379\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,6,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestSnapshotProofLifecycle,\"intent_test.go:143: git -c user.name= -c user.email=bench@example.invalid commit -qm init: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05d2-probe-05d",
          "performer": "claude:gf-t05d-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6db3e77342d0220a9794838bccf4f165124a6e73",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t05d-author-20261029",
            "digest": "sha256:f7c9ff6f3bd6c0bbc5c193e535135a2d6a8a8a875721abd3ea347fe836ba21b3",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestSnapshotProofLifecycle,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,TestSnapshotProofLifecycle,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,379\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,6,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestSnapshotProofLifecycle,\"intent_test.go:143: git -c user.name= -c user.email=bench@example.invalid commit -qm init: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "dfe7edd94d504e9a7503f3f7810ff29cace1564c",
          "requirement": "probe-05d",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/intent --run TestSnapshotProofLifecycle",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t05d-author-20261029",
              "digest": "sha256:f7c9ff6f3bd6c0bbc5c193e535135a2d6a8a8a875721abd3ea347fe836ba21b3",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestSnapshotProofLifecycle,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,TestSnapshotProofLifecycle,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,379\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,6,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestSnapshotProofLifecycle,\"intent_test.go:143: git -c user.name= -c user.email=bench@example.invalid commit -qm init: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05dr1-probe-05d",
          "performer": "claude:gf-t05d-repair1-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a440eb60a9d37dee59b931229c6fa4b7682561a6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t05d-repair1-20261029",
            "digest": "sha256:ba0f94fdb1b9997016da7fe0fe2040c8364ae4bbf014d250b69c22ca05b7eae9",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestSnapshotProofLifecycle,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,TestSnapshotProofLifecycle,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,379\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,5,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestSnapshotProofLifecycle,\"intent_test.go:143: git -c user.name= -c user.email=bench@example.invalid commit -qm init: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "fd1234c51eb80f59ff191f3b069ee03f9992f765",
          "requirement": "probe-05d",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/intent --run TestSnapshotProofLifecycle",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t05d-repair1-20261029",
              "digest": "sha256:ba0f94fdb1b9997016da7fe0fe2040c8364ae4bbf014d250b69c22ca05b7eae9",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestSnapshotProofLifecycle,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,TestSnapshotProofLifecycle,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,379\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,5,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestSnapshotProofLifecycle,\"intent_test.go:143: git -c user.name= -c user.email=bench@example.invalid commit -qm init: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gf-c5c-standards-r1",
          "performer": "claude:gf-c5c-standards-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6db3e77342d0220a9794838bccf4f165124a6e73",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c5c-standards-r1",
            "digest": "sha256:f36d0cfe1ecabf4f025fa298aefe29319fe6ecdf9aea4d932e716457bedc13b9",
            "excerpt": "Standards axis, chunk GF-C5C (claude:gf-c5c-standards-r1). K1 (axis id F1) internal/shift/fault_test.go:72-78 faultFixtureCore re-derives gittest.RepoOnBranch(t, \"main\"), and internal/conformance/docs_workflow_checks_test.go:203-206 uses Repo plus the identity loop before a commit; auto-fix, confidence 6. K2 (axis id F2) internal/landing/pending_merge_fixture_test.go:62-70 pendingMerge keeps a private exec merge and a stale comment, although gittest.Run returns the failure with the combined output; auto-fix, confidence 5. Advice: baseForNested is a lazy element; the nested gitlink setup repeats in two files.\n"
          },
          "axis": "Standards",
          "base": "cc197d6b3a004ed2c2b776573744f5b17c0cdd70",
          "tip": "dfe7edd94d504e9a7503f3f7810ff29cace1564c",
          "finding_ids": [
            "K1",
            "K2"
          ],
          "supersedes": []
        },
        {
          "id": "gf-c5c-spec-r1",
          "performer": "claude:gf-c5c-spec-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6db3e77342d0220a9794838bccf4f165124a6e73",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c5c-spec-r1",
            "digest": "sha256:6213300b4c0dd86d01c27a82ab677ccff01c2432f8ee4e1fddcda27011f32a18",
            "excerpt": "Spec axis, chunk GF-C5C (claude:gf-c5c-spec-r1). Findings: none. Retained sites conform: reviewfiles.go fast-import, pendingMerge, kit_run_test.go child environment, git_policy_test.go production queries, conformance check bodies, literal source data, the landing_reviewed_test.go budget edit. Every private runner is gone; no default identity copy remains in the fence. Advice: TempDir plus init -q repeats gittest.Repo in four tests.\n"
          },
          "axis": "Spec",
          "base": "cc197d6b3a004ed2c2b776573744f5b17c0cdd70",
          "tip": "dfe7edd94d504e9a7503f3f7810ff29cace1564c",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "gf-c5c-coverage-r1",
          "performer": "claude:gf-c5c-coverage-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6db3e77342d0220a9794838bccf4f165124a6e73",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c5c-coverage-r1",
            "digest": "sha256:1965f1dd870ee24a50aa8a2a6ba9bb850ddc8c0b0e233277105c58fa0fa8533b",
            "excerpt": "Coverage axis, chunk GF-C5C (claude:gf-c5c-coverage-r1). Findings: none. Combined-to-stdout readers read stdout-only data; trimmed readers keep their subject; RawPatch stays raw; the fast-import, kit-run, and pendingMerge probes are intact; tempGitRepoWithLines stays error-returning; no vacuous pass.\n"
          },
          "axis": "Coverage",
          "base": "cc197d6b3a004ed2c2b776573744f5b17c0cdd70",
          "tip": "dfe7edd94d504e9a7503f3f7810ff29cace1564c",
          "finding_ids": [],
          "supersedes": []
        }
      ],
      "confirmations": [
        {
          "id": "gf-c5c-confirm-r1",
          "performer": "claude:gf-orchestrator-20261029",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "high",
          "source_digest": "a440eb60a9d37dee59b931229c6fa4b7682561a6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:gf-orchestrator-20261029",
            "digest": "sha256:90c6346b50581d1afcfaf41a2ce4333f8d6193ca1e8a3364a1139922eaacae34",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shift/fault_test.go,swap,failed,1,TestLoopRetainsAndLocksDirtyWorktree,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/shift,TestLoopRetainsAndLocksDirtyWorktree,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,343\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shift,fail,8,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shift,TestLoopRetainsAndLocksDirtyWorktree,\"fault_test.go:214: git commit --quiet -m init: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
          },
          "base": "cc197d6b3a004ed2c2b776573744f5b17c0cdd70",
          "reviewed_tip": "dfe7edd94d504e9a7503f3f7810ff29cace1564c",
          "repaired_tip": "fd1234c51eb80f59ff191f3b069ee03f9992f765",
          "finding_ids": [
            "K1",
            "K2"
          ],
          "supersedes": [],
          "probe": {
            "mutation": "swap gittest.RepoOnBranch for gittest.Repo in faultFixtureCore",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:gf-orchestrator-20261029",
              "digest": "sha256:90c6346b50581d1afcfaf41a2ce4333f8d6193ca1e8a3364a1139922eaacae34",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shift/fault_test.go,swap,failed,1,TestLoopRetainsAndLocksDirtyWorktree,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/shift,TestLoopRetainsAndLocksDirtyWorktree,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,343\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shift,fail,8,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shift,TestLoopRetainsAndLocksDirtyWorktree,\"fault_test.go:214: git commit --quiet -m init: exit status 128: Author identity unknown\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ]
    },
    {
      "id": "GF-C6",
      "base": "fd1234c51eb80f59ff191f3b069ee03f9992f765",
      "tip": "8f4bb23108b5fe050c3e7fbd9dd90a07c23c46f6",
      "plan_digest": "sha256:b1937fdaf5392277821c129deb56670966c6c5dd95400bf063340a34f67096c3",
      "source_digest": "08dc59cabcc9e1259bbe4259f6cd167e32d93744",
      "acceptance_rows": [
        "GF17",
        "GF29",
        "GF30"
      ],
      "verification": [
        {
          "id": "t06-probe-06",
          "performer": "claude:gf-t06-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f0602aea95b14fe0e71590b6ec8b5e20850ac3f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t06-author-20261029",
            "digest": "sha256:1cecba3ffe778abee0aa02cc2a00494c7ee932f9ed4aeab3fc6f798c1f6a284d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestWindowAdmitsTheBuildsOwnLanding,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestWindowAdmitsTheBuildsOwnLanding,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,325\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,11,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestWindowAdmitsTheBuildsOwnLanding,\"land_window_test.go:102: git commit -q -m base: exit status 128\",9\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "960823e88f9f78234a8c195d96200b1134b0f688",
          "requirement": "probe-06",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/worktree --run TestWindowAdmitsTheBuildsOwnLanding",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t06-author-20261029",
              "digest": "sha256:1cecba3ffe778abee0aa02cc2a00494c7ee932f9ed4aeab3fc6f798c1f6a284d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestWindowAdmitsTheBuildsOwnLanding,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestWindowAdmitsTheBuildsOwnLanding,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,325\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,11,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestWindowAdmitsTheBuildsOwnLanding,\"land_window_test.go:102: git commit -q -m base: exit status 128\",9\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t07-probe-07",
          "performer": "claude:gf-t07-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "08dc59cabcc9e1259bbe4259f6cd167e32d93744",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t07-author-20261029",
            "digest": "sha256:7c39e5d932b4cbfd7ea0b828704bd565594bfa879ce6b9ad642a243a995151d6",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestMergeStartsABackgroundGrade,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestMergeStartsABackgroundGrade,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,325\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,11,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestMergeStartsABackgroundGrade,\"merge_background_test.go:49: git commit -q -m base: exit status 128\",9\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8f4bb23108b5fe050c3e7fbd9dd90a07c23c46f6",
          "requirement": "probe-07",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --package ./internal/worktree --run TestMergeStartsABackgroundGrade",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t07-author-20261029",
              "digest": "sha256:7c39e5d932b4cbfd7ea0b828704bd565594bfa879ce6b9ad642a243a995151d6",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestMergeStartsABackgroundGrade,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestMergeStartsABackgroundGrade,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,325\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,11,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestMergeStartsABackgroundGrade,\"merge_background_test.go:49: git commit -q -m base: exit status 128\",9\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gf-c6-standards-r1",
          "performer": "claude:gf-c6-standards-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "08dc59cabcc9e1259bbe4259f6cd167e32d93744",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c6-standards-r1",
            "digest": "sha256:5826dcd8dad6c6c5fd8157d2e5b207c24fdd62f228b5f13ed862d67893eeb774",
            "excerpt": "Standards axis, chunk GF-C6 (claude:gf-c6-standards-r1). No defect. Two items the axis itself dispositioned no-op: reset_refusal_test.go:132 raw commit headers feed a hash-prefix collision search (retained specialized site); bare commit calls in merge_test.go, clean_landed_apply_test.go, clean_set_apply_test.go read the same identity through IdentityConfig. withIdentity is the one consumer of IdentityArgs; comments match the code. Advice: ticket 06 wording names one harness file where the census admits three.\n"
          },
          "axis": "Standards",
          "base": "fd1234c51eb80f59ff191f3b069ee03f9992f765",
          "tip": "8f4bb23108b5fe050c3e7fbd9dd90a07c23c46f6",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "gf-c6-spec-r1",
          "performer": "claude:gf-c6-spec-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "08dc59cabcc9e1259bbe4259f6cd167e32d93744",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c6-spec-r1",
            "digest": "sha256:388c8fda60f162dda733065bf27a139c379554d1629c3e7dd4ff03511daf245d",
            "excerpt": "Spec axis, chunk GF-C6 (claude:gf-c6-spec-r1). Findings: none. Every migrated child starts through gitRun, gitOutput, or descendant; no private default identity remains in the fences except the specialized reset_refusal_test.go:132 site and production clean.go; the landAssignment cherry-pick -x change stays inside the approved behavior; leasedRepo needs no repository-config identity. Advice: leasedRepo could use journeyRepo.\n"
          },
          "axis": "Spec",
          "base": "fd1234c51eb80f59ff191f3b069ee03f9992f765",
          "tip": "8f4bb23108b5fe050c3e7fbd9dd90a07c23c46f6",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "gf-c6-coverage-r1",
          "performer": "claude:gf-c6-coverage-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "08dc59cabcc9e1259bbe4259f6cd167e32d93744",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c6-coverage-r1",
            "digest": "sha256:4057fd211d438742b272a94c095f689122892dbf404443a43dc1bba161597444",
            "excerpt": "Coverage axis, chunk GF-C6 (claude:gf-c6-coverage-r1). The plan probes probe-06 and probe-07 turn red at a repository-config identity that predates the chunk, so they do not grade the migration; the axis's own probe (withIdentity body reduced to return args) turned TestReleaseOwnerRestoresCleanAndUnleases red at lifecycle_test.go:94 and restored; proposed ask-user or no-op, confidence 5. No GF17 bypass found; the census guards it; no-op, confidence 4. Clean hunts: no other identical-object replay, leasedRepo has no later product commit, no vacuous pass.\n"
          },
          "axis": "Coverage",
          "base": "fd1234c51eb80f59ff191f3b069ee03f9992f765",
          "tip": "8f4bb23108b5fe050c3e7fbd9dd90a07c23c46f6",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "GF-C8",
      "base": "8f4bb23108b5fe050c3e7fbd9dd90a07c23c46f6",
      "tip": "64eda701bbfce69ad016a30a5890597ccd1962b5",
      "plan_digest": "sha256:5b4a8983e799c7e785028a7d2173239eb8a005409c5672be0154d929501628b0",
      "source_digest": "e5e6818dae5997facd755edf1905684f6cf6a047",
      "acceptance_rows": [
        "GF18",
        "GF29",
        "GF30"
      ],
      "verification": [
        {
          "id": "t08-probe-08",
          "performer": "claude:gf-t08-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e5e6818dae5997facd755edf1905684f6cf6a047",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:gf-t08-author-20261029",
            "digest": "sha256:3b665835d13c42259e0c65da1d6e141b0d4d77b9de6dea0bbd396b744d338e25",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestWorktreeReauthorizeJourney,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,^TestWorktreeReauthorizeJourney$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,614\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,262,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestWorktreeReauthorizeJourney,\"owner_selection_test.go:280: git [\\\"-c\\\" \\\"user.name=\\\" \\\"-c\\\" \\\"user.email=bench@example.invalid\\\" \\\"commit\\\" \\\"-qm\\\" \\\"base\\\"] = (128, \\\"Author identity unknown\\\\\\\\n\\\\\\\\n*** Please tell me who you are.\\\\\\\\n\\\\\\\\nRun\\\\\\\\n\\\\\\\\n  git config --global user.email \\\\\\\\\\\"you@example.… (468 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "64eda701bbfce69ad016a30a5890597ccd1962b5",
          "requirement": "probe-08",
          "command": "bench probe internal/gittest/identity.go --swap '\"bench test\"' --with '\"\"' --check system --run '^TestWorktreeReauthorizeJourney$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Empty the canonical default name. Each commit through the shared identity then fails, so the named migrated test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:gf-t08-author-20261029",
              "digest": "sha256:3b665835d13c42259e0c65da1d6e141b0d4d77b9de6dea0bbd396b744d338e25",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gittest/identity.go,swap,failed,1,TestWorktreeReauthorizeJourney,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,^TestWorktreeReauthorizeJourney$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,614\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,262,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestWorktreeReauthorizeJourney,\"owner_selection_test.go:280: git [\\\"-c\\\" \\\"user.name=\\\" \\\"-c\\\" \\\"user.email=bench@example.invalid\\\" \\\"commit\\\" \\\"-qm\\\" \\\"base\\\"] = (128, \\\"Author identity unknown\\\\\\\\n\\\\\\\\n*** Please tell me who you are.\\\\\\\\n\\\\\\\\nRun\\\\\\\\n\\\\\\\\n  git config --global user.email \\\\\\\\\\\"you@example.… (468 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gf-c8-standards-r1",
          "performer": "claude:gf-c8-standards-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e5e6818dae5997facd755edf1905684f6cf6a047",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c8-standards-r1",
            "digest": "sha256:d0ad9fefbbafe8a1fb8d03faa5477e29d45dda0feecbc93a091df462fb5eafc2",
            "excerpt": "Standards axis, chunk GF-C8 (claude:gf-c8-standards-r1). Findings: none. No private identity remains in internal/systemtest; systemCommitArgs and systemIdentity read only gittest.IdentityArgs and IdentityConfig; neither starts a process; comments match.\n"
          },
          "axis": "Standards",
          "base": "8f4bb23108b5fe050c3e7fbd9dd90a07c23c46f6",
          "tip": "64eda701bbfce69ad016a30a5890597ccd1962b5",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "gf-c8-spec-r1",
          "performer": "claude:gf-c8-spec-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e5e6818dae5997facd755edf1905684f6cf6a047",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c8-spec-r1",
            "digest": "sha256:a233ceca9365c276891a35a2a617722962a583a6ace4702af99f7c805f458532",
            "excerpt": "Spec axis, chunk GF-C8 (claude:gf-c8-spec-r1). Findings: none. Every migrated Git child starts through owner.runAt; the repository budget is untouched; config-form and per-command-form sites keep their form; no test grades the former Route Test identity. Advice: process_lifetime_cache_test.go commits through gittest outside owner.runAt, a pre-existing site outside the fence.\n"
          },
          "axis": "Spec",
          "base": "8f4bb23108b5fe050c3e7fbd9dd90a07c23c46f6",
          "tip": "64eda701bbfce69ad016a30a5890597ccd1962b5",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "gf-c8-coverage-r1",
          "performer": "claude:gf-c8-coverage-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e5e6818dae5997facd755edf1905684f6cf6a047",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c8-coverage-r1",
            "digest": "sha256:fb1b84e1af77fed2dbc33fc216d0cb0869b04d8cd49b5acd3854a95a5925f112",
            "excerpt": "Coverage axis, chunk GF-C8 (claude:gf-c8-coverage-r1). Findings: none. No replay depends on identity equality; no form switch; no vacuous pass; the starts ledger is unchanged. Advice: childEnvironment inherits the host Git configuration, so a host global identity would mask an identity migration.\n"
          },
          "axis": "Coverage",
          "base": "8f4bb23108b5fe050c3e7fbd9dd90a07c23c46f6",
          "tip": "64eda701bbfce69ad016a30a5890597ccd1962b5",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "GF-C9",
      "base": "64eda701bbfce69ad016a30a5890597ccd1962b5",
      "tip": "2c20d8e0a724dd1716d0cdaa059e01e8c1212e78",
      "plan_digest": "sha256:1a14401f8264a19fe71ffd7de3b0ee308d915f2540baaa80cdb2e8d2ac381d8c",
      "source_digest": "0b057c9ea23f58bea3a292fbb06c08a890126d3d",
      "acceptance_rows": [
        "GF21",
        "GF22",
        "GF23",
        "GF30"
      ],
      "verification": [
        {
          "id": "t09-probe-09",
          "performer": "claude:gf-t09-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0b057c9ea23f58bea3a292fbb06c08a890126d3d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gf-t09-author-20261029",
            "digest": "sha256:e71ed2927acb4235500b1b3f47c9be2ccc39e44d58f18209620948d9c9a69a93",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/capability/unavailable.go,swap,failed,1,TestUnavailableReason,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/capability,TestUnavailableReason,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,596\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/capability,fail,2,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/capability,TestUnavailableReason,\"unavailable_test.go:52: Unavailable reason = \\\"symlink unavailable: symlink a b: operation not permitted\\\", want \\\"symlink capability unavailable: symlink a b: operation not permitted\\\"\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/capability,TestUnavailableReason/skips,\"unavailable_test.go:20: symlink unavailable: symlink a b: operation not permitted\"\n"
          },
          "source_commit": "2c20d8e0a724dd1716d0cdaa059e01e8c1212e78",
          "requirement": "probe-09",
          "command": "bench probe internal/capability/unavailable.go --swap '\" capability unavailable: \"' --with '\" unavailable: \"' --package ./internal/capability --run TestUnavailableReason",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the capability word from the derived prefix. The reason then differs from the contract, so the reason test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gf-t09-author-20261029",
              "digest": "sha256:e71ed2927acb4235500b1b3f47c9be2ccc39e44d58f18209620948d9c9a69a93",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/capability/unavailable.go,swap,failed,1,TestUnavailableReason,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/capability,TestUnavailableReason,passed,2\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,596\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/capability,fail,2,2\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/capability,TestUnavailableReason,\"unavailable_test.go:52: Unavailable reason = \\\"symlink unavailable: symlink a b: operation not permitted\\\", want \\\"symlink capability unavailable: symlink a b: operation not permitted\\\"\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/capability,TestUnavailableReason/skips,\"unavailable_test.go:20: symlink unavailable: symlink a b: operation not permitted\"\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gf-c9-standards-r1",
          "performer": "claude:gf-c9-standards-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "0b057c9ea23f58bea3a292fbb06c08a890126d3d",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c9-standards-r1",
            "digest": "sha256:0edd74cdc21abf9e08f640041829178ea7e99fc06063cf6731569d75f5d3eba9",
            "excerpt": "Standards axis, chunk GF-C9 (claude:gf-c9-standards-r1). N1 (axis id F1) internal/capability/unavailable_test.go:13-28 runUnavailable is a third copy of the subtest-over-file-transport harness inlined at capability_test.go:114-133 and :159-178; proposed no-op, confidence 3. Clean: one source of the prefix at unavailable.go:11; no FIFO or symlink caller in the fence writes its own prefix; the independent reason literal has a recorded omission red; no comment drift.\n"
          },
          "axis": "Standards",
          "base": "64eda701bbfce69ad016a30a5890597ccd1962b5",
          "tip": "2c20d8e0a724dd1716d0cdaa059e01e8c1212e78",
          "finding_ids": [
            "N1"
          ],
          "supersedes": []
        },
        {
          "id": "gf-c9-spec-r1",
          "performer": "claude:gf-c9-spec-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "0b057c9ea23f58bea3a292fbb06c08a890126d3d",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c9-spec-r1",
            "digest": "sha256:36e5d7a6cf92988170494b19885f561e4d369569e24c5ce086f03865b9688a0a",
            "excerpt": "Spec axis, chunk GF-C9 (claude:gf-c9-spec-r1). N2 (axis id F1) unavailable_test.go:119-126 TestUnavailableReason exercises only Symlink, so a hard-coded class token in the prefix survives; GF21 binds derivation from the class; auto-fix, confidence 6. N3 (axis id F2) the GF18 plan commit (spec GF18 row, GF-C13 row, ticket 13) sits in this range outside ticket 09; ask-user to confirm the plan expansion was recorded, confidence 4. GF22, GF23, the contract, the six migrated classes, and the Privilege calls conform.\n"
          },
          "axis": "Spec",
          "base": "64eda701bbfce69ad016a30a5890597ccd1962b5",
          "tip": "2c20d8e0a724dd1716d0cdaa059e01e8c1212e78",
          "finding_ids": [
            "N2",
            "N3"
          ],
          "supersedes": []
        },
        {
          "id": "gf-c9-coverage-r1",
          "performer": "claude:gf-c9-coverage-r1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "0b057c9ea23f58bea3a292fbb06c08a890126d3d",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gf-c9-coverage-r1",
            "digest": "sha256:9e605d9bfb8ea1767f2a9c1f3f22594c4065687544939ea8183495a4e4b21740",
            "excerpt": "Coverage axis, chunk GF-C9 (claude:gf-c9-coverage-r1). N4 (axis id C9-COV-1) unavailable_test.go:47-53 a hard-coded class token in the prefix survives because only Symlink grades the reason; auto-fix, confidence 9. N5 (C9-COV-2) detail TrimSpace or ToLower mutants survive; no-op, confidence 5. N6 (C9-COV-3) no parallel transport case; Unavailable adds no write path; no-op, confidence 6.\n"
          },
          "axis": "Coverage",
          "base": "64eda701bbfce69ad016a30a5890597ccd1962b5",
          "tip": "2c20d8e0a724dd1716d0cdaa059e01e8c1212e78",
          "finding_ids": [
            "N4",
            "N5",
            "N6"
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
      "from": "sha256:dd48d59465866a1c9b1f35f8f9eecee4e98f8e465d8643764328edbea174193f",
      "to": "sha256:15c922e92fbdf0b1d287370928c19669972477750af95e1a5a1e000f2e67a7b2",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ]
      }
    },
    {
      "from": "sha256:15c922e92fbdf0b1d287370928c19669972477750af95e1a5a1e000f2e67a7b2",
      "to": "sha256:faf692b54b3699aa9bfd216ed49ad566258e0330b90102580b97de862be1d5fd",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ]
      }
    },
    {
      "from": "sha256:faf692b54b3699aa9bfd216ed49ad566258e0330b90102580b97de862be1d5fd",
      "to": "sha256:811bcf440b2d88ead423909e96ba6728e92ed947bc317f7d06f969e9e21d2ddf",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ]
      }
    },
    {
      "from": "sha256:811bcf440b2d88ead423909e96ba6728e92ed947bc317f7d06f969e9e21d2ddf",
      "to": "sha256:8a52cd552bd41e3d755db776a1c77f5371bd4f1a65e94ee521b32ac05709ba00",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ]
      }
    },
    {
      "from": "sha256:8a52cd552bd41e3d755db776a1c77f5371bd4f1a65e94ee521b32ac05709ba00",
      "to": "sha256:ddaff991fc0933258b524e9d31c8ee54634c54b909d173458ecbf743ad44bb67",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ]
      }
    },
    {
      "from": "sha256:ddaff991fc0933258b524e9d31c8ee54634c54b909d173458ecbf743ad44bb67",
      "to": "sha256:d9897e338a549f2aaeab8a0e4f21046027df8b54f57e5f0feaecb189a07dfa9c",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ]
      }
    },
    {
      "from": "sha256:d9897e338a549f2aaeab8a0e4f21046027df8b54f57e5f0feaecb189a07dfa9c",
      "to": "sha256:142b2f7c9288c50f1957cc9921b80ad2805b2eb2a0b142c44c6fc666f556d580",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ]
      }
    },
    {
      "from": "sha256:142b2f7c9288c50f1957cc9921b80ad2805b2eb2a0b142c44c6fc666f556d580",
      "to": "sha256:bd1eedb145049f9d18d247761194b53a00a3046a9942d49a69ad123970e12d0e",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ]
      }
    },
    {
      "from": "sha256:bd1eedb145049f9d18d247761194b53a00a3046a9942d49a69ad123970e12d0e",
      "to": "sha256:9124d52d6f85b93eb8c64e15419a727788dffa7fc80ebcedadd00cd795641e22",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ]
      }
    },
    {
      "from": "sha256:9124d52d6f85b93eb8c64e15419a727788dffa7fc80ebcedadd00cd795641e22",
      "to": "sha256:f53e416b2780300479d1e14e50b208cc60fe0103991607ebcf4320c09c6ed09f",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ]
      }
    },
    {
      "from": "sha256:f53e416b2780300479d1e14e50b208cc60fe0103991607ebcf4320c09c6ed09f",
      "to": "sha256:04a91e4d096c18351fae0e348f695e53ba82310e2510eec52a90a71c4fc94d33",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ]
      }
    },
    {
      "from": "sha256:04a91e4d096c18351fae0e348f695e53ba82310e2510eec52a90a71c4fc94d33",
      "to": "sha256:a3967d35127d4fcdfff4d1e62cbd897f146f1fb1a246dbf3a51191668671db63",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ],
        "GF-C5C": [
          "GF-C5C"
        ]
      }
    },
    {
      "from": "sha256:a3967d35127d4fcdfff4d1e62cbd897f146f1fb1a246dbf3a51191668671db63",
      "to": "sha256:95314d621569659d44dade204ac1f040da6f2cf762d7391f2647bf386bf87f2a",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ],
        "GF-C5C": [
          "GF-C5C"
        ]
      }
    },
    {
      "from": "sha256:95314d621569659d44dade204ac1f040da6f2cf762d7391f2647bf386bf87f2a",
      "to": "sha256:8a878c0f6b72f4b296defc18ecec53c329e2fe7e2b0a98f6c4fd60e8ff9b011f",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ],
        "GF-C5C": [
          "GF-C5C"
        ]
      }
    },
    {
      "from": "sha256:8a878c0f6b72f4b296defc18ecec53c329e2fe7e2b0a98f6c4fd60e8ff9b011f",
      "to": "sha256:d8576e9f92c8e64e0cf2f9d1954f7d3461fd4629442582ebb6c8f51f0133fd54",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ],
        "GF-C5C": [
          "GF-C5C"
        ]
      }
    },
    {
      "from": "sha256:d8576e9f92c8e64e0cf2f9d1954f7d3461fd4629442582ebb6c8f51f0133fd54",
      "to": "sha256:b1937fdaf5392277821c129deb56670966c6c5dd95400bf063340a34f67096c3",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ],
        "GF-C5C": [
          "GF-C5C"
        ],
        "GF-C6": [
          "GF-C6"
        ]
      }
    },
    {
      "from": "sha256:b1937fdaf5392277821c129deb56670966c6c5dd95400bf063340a34f67096c3",
      "to": "sha256:5b4a8983e799c7e785028a7d2173239eb8a005409c5672be0154d929501628b0",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ],
        "GF-C5C": [
          "GF-C5C"
        ],
        "GF-C6": [
          "GF-C6"
        ]
      }
    },
    {
      "from": "sha256:5b4a8983e799c7e785028a7d2173239eb8a005409c5672be0154d929501628b0",
      "to": "sha256:1a14401f8264a19fe71ffd7de3b0ee308d915f2540baaa80cdb2e8d2ac381d8c",
      "chunk_ids": {
        "GF-C1": [
          "GF-C1"
        ],
        "GF-C2": [
          "GF-C2"
        ],
        "GF-C5A": [
          "GF-C5A"
        ],
        "GF-C5C": [
          "GF-C5C"
        ],
        "GF-C6": [
          "GF-C6"
        ],
        "GF-C8": [
          "GF-C8"
        ]
      }
    }
  ]
}
```

## GF-C1 review pickup

The repair allowance of GF-C1 is two cycles. The first repair cycle starts after this record.

### Standards

Findings: 0. Worst issue: none.

Advice: share the failure formatter and the identity configuration loop. Use `Commit` in place of `commit -qm` through `MustRun`. Rename the file that holds only the package doc.

### Spec

Findings: 2. Worst issue: F1.

- F1, `internal/gittest/commands.go:91`: `OutputBytes` runs its own output call beside `Run`. The behavior matches GF3. Disposition: no-op. Confidence: 3.
- F2, `internal/testrepo/working_tree.go:1`: the file holds only the package doc, and the `Writes:` line names it. Disposition: no-op. Confidence: 2.

### Coverage

Findings: 4. Worst issue: C1.

- C1, `internal/gittest/commands.go:28`: no test makes `OutputBytes` fail on a nonzero exit or a missing Git. The edge inventory requires fatal wrappers to fail. Disposition: auto-fix. Confidence: 8.
- C2, `internal/gittest/identity.go:20`: no test grades the content of `IdentityArgs`, which later migrations use. Disposition: auto-fix. Confidence: 7.
- C3, `internal/gittest/gittest_test.go:12`: the snapshot fixture holds no empty tracked file, and the edge inventory names GF11 for the empty-versus-absent edge. Disposition: auto-fix. Confidence: 7.
- C4, `internal/gittest/snapshot_test.go:40`: GF15 grades only the commit failure. GF15 binds one failed command, and an earlier swallowed step still ends in the returned commit error. Disposition: no-op. Confidence: 5.

Repair targets: 3 (C1, C2, C3), all on ticket 01, in one combined repair charge.

### GF-C1 repair cycle 1

The ticket 01 repair session committed `6002ec0cf`, a test-only repair of C1, C2, and C3. It refroze GF-C1 at that tip.
Its witness reds each came from one `bench probe` of production code, and each restore passed:

- C1: a swap of the `OutputBytes` fatal call for a log call turned `TestOutputBytesFailure` red.
- C2: an omission of the name option from `IdentityArgs` turned `TestFixtureIdentity` red.
- C3: an insertion that skips zero-length files in the snapshot turned `TestKitCopyPreservesTheVisibleWorkingTree` red.

The coordinator probe swapped the email option of `IdentityArgs`, a different kind at a different site. It turned `TestFixtureIdentity` red, and the restore passed.
The repair consumed one of the two repair cycles of GF-C1.

## GF-C2 review pickup

The repair allowance of GF-C2 is two cycles. The first repair cycle starts after this record.

### Standards

Findings: 4. Worst issue: S1.

- S1, `internal/structure/structure_test.go:349`: the private `run` keeps a stdout-only output shape that callers outside the fence read. It is not a pure rename of a `gittest` helper, so the rename rule does not apply. Disposition: no-op. Confidence: 7.
- S2, `internal/roadmap/retro_cost_test.go:32`: the private runner also performs plain `add`, `rev-parse`, and `checkout` steps. One exception does not exempt unrelated setup, so these steps move to `gittest`. Disposition: auto-fix. Confidence: 5.
- S3, `internal/maps/freshness_test.go:32`: the identity configuration loop repeats in three fixtures. The identity values have one owner, and the loop is incidental composition. Disposition: no-op. Confidence: 4.
- S4, `internal/status/status_producible_test.go:56`: the parameter rename keeps the closure on one line inside the structure budget. Disposition: no-op. Confidence: 5.

### Spec

Findings: 2. Worst issue: P2.

- P1, `internal/status/status_signals_test.go:710`: the literal reformat keeps the file inside its structure budget, and the import move is `goimports` output. Disposition: no-op. Confidence: 6.
- P2, `internal/maps/freshness_test.go:56`: `commitAt` is a timestamp-only private runner. Sibling fixtures set the dates with `t.Setenv` and commit through `gittest`. Disposition: auto-fix. Confidence: 4.

### Coverage

Findings: 1. Worst issue: V1.

- V1, `internal/status/status_signals_test.go:264`: the same layout edit as P1 has no behavioral effect. Disposition: no-op. Confidence: 7.

Repair targets: 2 (S2, P2), both on ticket 04, in one combined repair charge.

### Retained specialized and intentional sites

- `internal/diff/command_test.go` and `internal/diff/matrix_test.go` (`rawNoIndexPatch` and the `diff --check` probe): an expected nonzero exit.
- `internal/diff/compatibility_test.go` (`seedCompatibilityRepo`): an intentional subject, because pinned commit ids grade its identity and dates.
- `internal/handoff/state_scan_test.go` (`writeBlob`): custom standard input.
- `internal/roadmap/retro_cost_test.go` (`gitAuthored`): deliberate author and committer times. The fixture grades the committer times, and it keeps the ungraded author offset as an input.
- `internal/structure/structure_test.go` (`run`): a stdout-only output shape over `gittest.OutputBytes`.

### GF-C2 repair cycle 1

The ticket 04 repair session committed `874d872c3`, a test-only repair of S2 and P2. It refroze GF-C2 at that tip.
Its witnesses each came from one `bench probe`, and each restore passed:

- P2: an omission of the committer date in `commitAt` turned the freshness test red.
- S2: a swap of the retro-cost committer stamp turned `TestRetroCostSection` red.
- A swap of the retro-cost author offset stayed silent. Thus the classification above names that offset as an ungraded input.

The coordinator probe inserted a fixed committer date before the shared commit in `commitAt`. It turned the freshness test red, and the restore passed.
The repair consumed one of the two repair cycles of GF-C2.

## GF-C5A review pickup

The repair allowance of GF-C5A is two cycles. The first repair cycle starts after this record.

### Standards

Findings: 1. Worst issue: H1.

- H1, `internal/adopt/link_hook_test.go:190`: `hookTestRepo` is a one-line alias of `gittest.Repo`. The spec's rename rule names only the execution helpers, and the alias holds no execution or identity policy. Disposition: no-op. Confidence: 6.

### Spec

Findings: 2. Worst issue: R1.

- R1, `cmd/bench/command_registry_test.go:399`: `runAXIGit` and `trimmedAXIGit` build fixture repositories through the production adapter. The spec forbids that use, so both wrappers leave and their callers move to `gittest`. Disposition: auto-fix. Confidence: 7.
- R2, `cmd/bench/tree_scope_test.go:83`: `treeRowRepo` commits its dated commit through `runAXIGit`. A dated commit uses `t.Setenv` and then a `gittest` commit. Disposition: auto-fix. Confidence: 7.

### Coverage

Findings: 0. Worst issue: none.

Repair targets: 2 (R1, R2), both on ticket 05a, in one combined repair charge.

### Retained specialized and intentional sites

- `internal/commit/merge_parent_test.go` (`startPendingMerge`): an expected merge conflict with a nonzero exit.
- `internal/adopt/link_hook_test.go` (the recorder and the failing stub tests): a stub Git on the path and an expected failure.
- `cmd/bench/wrapper_kit_test.go` (`newWrapperKitFixture`): a linked worktree that the wrapper resolution grades. The plan probe of ticket 13 plants a sibling runner beside this site.

### GF-C5A repair cycle 1

The ticket 05a repair session committed `b1ea40219`, a test-only repair of R1 and R2. It refroze GF-C5A at that tip.
It also migrated four more `runAXIGit` sites in the fenced `spec_retire_listing_test.go`, so that both wrappers could leave.
Its witness swapped the `treeRowRepo` head query for a tree query, and `TestTreeRowLeadsTreeResponse` turned red. The restore passed.
The date pin of `treeRowRepo` is an ungraded input.

The coordinator probe omitted the staging step of `newAXIEnvelopeRepo`. `TestSelectedWorktreeCommandRoutes` turned red, and the restore passed.
The repair consumed one of the two repair cycles of GF-C5A.

### GF-C5A repair cycle 2

The confirming round passed on all three axes. The chunk checkpoint then went red on the branch-native census at `internal/gate/prospectiveartifact/prospectiveartifact_test.go:306`.
The census reads a `gittest.Output` call with a `worktree` literal as a repository constructor in a gate test. The line was a raw read of the worktree list.

The ticket 05b repair session committed `cc197d6b3`, which reads that text through `gittest.OutputBytes`. It refroze GF-C5A at that tip.
Its witness ran the census on the unedited tree (red) and after the edit (green).

The coordinator probe swapped the line back. The census check turned red, and the restore passed.
The repair consumed the second and last repair cycle of GF-C5A.

## GF-C5C review pickup

The repair allowance of GF-C5C is two cycles. The first repair cycle starts after this record.

### Standards

Findings: 2. Worst issue: K1.

- K1, `internal/shift/fault_test.go:72`: `faultFixtureCore` repeats the body of `gittest.RepoOnBranch(t, "main")`. The fixture moves to that builder. The docs-workflow fixture keeps the ambient default branch, which is the closed disposition of the maps and status fixtures. Disposition: auto-fix for the shift fixture, no-op for the docs-workflow fixture. Confidence: 6.
- K2, `internal/landing/pending_merge_fixture_test.go:62`: `pendingMerge` keeps a private merge process, and its comment says the shared owner cannot tolerate the failure. `gittest.Run` returns that failure with the combined output. Disposition: auto-fix. Confidence: 5.

### Spec

Findings: 0. Worst issue: none.

### Coverage

Findings: 0. Worst issue: none.

Repair targets: 2 (K1, K2), both on ticket 05d, in one combined repair charge.

### Retained specialized and intentional sites

- `internal/preflight/preflighttest/reviewfiles.go` (`SeedReviewPair`): the fast-import input stream, the pinned committer, and the pinned times.
- `internal/env/kit_run_test.go` (`checkKitRunIsolation`): an explicit child environment.
- `internal/env/git_policy_test.go` (`TestGitPolicyHasOneSource`, `gitPolicyFindings`): production-subject queries.
- `internal/conformance` check bodies (`trackedPaths`, `checkExecutableGitModes`, `checkGateInputPathsNotIgnored`, `spillProbeOutput`): reads of the graded root.
- `internal/preflight/preflighttest/fixture.go` (`RunGit`): an exported facade over `gittest` with callers in other packages.

### GF-C5C repair cycle 1

The ticket 05d repair session committed `fd1234c51`, a test-only repair of K1 and K2, and refroze GF-C5C at that tip.
The session returned no report, so the coordinator verified the repair diff against the charge and the record commit `63f33948d`.

The coordinator probe swapped `gittest.RepoOnBranch` for `gittest.Repo` in `faultFixtureCore`. `TestLoopRetainsAndLocksDirtyWorktree` turned red on a missing author identity, and the restore passed.
The repair consumed one of the two repair cycles of GF-C5C.

## GF-C6 review pickup

All three axes passed. The axes dispositioned each item they raised as no-op, so the chunk takes no repair cycle.

### Standards

Findings: 0. Worst issue: none.

- The raw commit headers at `internal/worktree/reset_refusal_test.go:132` feed a hash-prefix collision search. They stay a specialized site, and ticket 13 must list them as an exact exception.
- The bare `commit` calls in the worktree tests read the same identity constants through `gittest.IdentityConfig`.

### Spec

Findings: 0. Worst issue: none.

### Coverage

Findings: 0. Worst issue: none.

The plan probes `probe-06` and `probe-07` turn red at a repository-config identity that predates this chunk, so they do not grade the migration.
The Coverage axis reduced `withIdentity` to `return args`. `TestReleaseOwnerRestoresCleanAndUnleases` turned red at the `leasedRepo` commit, and the restore passed. That probe is the discriminating evidence of the identity migration.

### Author witnesses outside the plan probes

- Ticket 06: a swap of a migrated `gitRun` call for an ordinary process turned `TestSerialJourneyHarnessCensus` red, so the census grades the `descendant` route.
- Ticket 07: a swap of `commitInWorktree` for a direct `gittest` commit turned `TestSerialJourneyHarnessCensus` red. An omission of `-x` in `landAssignment` turned `TestDirtyAssignmentWithoutOwnCommitIsNotLanded` red.

### Orchestrator decision for veto

Under the shared identity, the cherry-pick in `landAssignment` can rebuild the branch commit byte for byte within one second. The landed fixture then loses its separate landing commit, and eleven tests fail.
The orchestrator made the landing commit explicit with `cherry-pick -x` in `internal/worktree/landed_test.go`. The patch stays the same, and the product reads no part of the added trailer. The rejected option kept the old private identity, which hid the timing dependence.

### Retained specialized sites

- `internal/worktree/reset_refusal_test.go` (`TestResetRefusesAnAmbiguousCheckpoint`): a raw commit body with fixed input bytes.

## GF-C8 review pickup

All three axes passed with no finding, so the chunk takes no repair cycle.

### Standards

Findings: 0. Worst issue: none.

### Spec

Findings: 0. Worst issue: none.

### Coverage

Findings: 0. Worst issue: none.

### Author witnesses outside the plan probe

- A reduction of `systemCommitArgs` to no identity turned `TestWorktreeReauthorizeJourney` red.
- A loop of `systemIdentity` over no commands turned `TestStatusRouteExecutesUnclaimedCleanup` red.
- A swap of the `owner.runAt` call in `systemGit` for a direct process stayed silent across the selected system tests and their censuses.

### Consultant decision on GF18

No gate refuses a direct Git process in `internal/systemtest`, because the ordinary-build census skips that directory. Spec row GF18 claimed that such a bypass fails an owner or census assertion. The gap predates this build.

Two consultants (opus at high effort and sonnet at high effort) both chose to move the GF18 bite into ticket 13. Ticket 13 replaces the census's skip of `internal/systemtest` with an exact list of allowed sites, by file and enclosing function, through the existing architecture scan. The silent `systemGit` bypass mutant is then its recorded red.

The pre-existing `gittest` commits in `process_lifetime_cache_test.go` and `process_lifetime_preservation_test.go` take an explicit classification in that list.
Until GF-C13, the GF18 route stays checked by reading, as this chunk's axes did.

### Retained specialized sites

None.

## GF-C9 review pickup

The repair allowance of GF-C9 is two cycles. The first repair cycle starts after this record.

### Standards

Findings: 1. Worst issue: N1.

- N1, `internal/capability/unavailable_test.go:13`: `runUnavailable` is a third copy of the file-transport skip harness that `capability_test.go` inlines twice. AGENTS.md names a pasted fixture harness as a defect, and both files are in the fence. The two older copies reuse the one helper. Disposition: auto-fix. Confidence: 3.

### Spec

Findings: 2. Worst issue: N2.

- N2, `internal/capability/unavailable_test.go:119`: the reason test grades only the symlink class, so a hard-coded class token in the prefix survives. GF21 binds the derivation of the prefix from the class. Disposition: auto-fix. Confidence: 6.
- N3, `specs/shared-test-fixtures/spec.md`: the GF18 plan commit sits in this range. It came from the GF-C8 consultant decision, and the plan-expansion learning was captured before the commit. Disposition: no-op. Confidence: 4.

### Coverage

Findings: 3. Worst issue: N4.

- N4, `internal/capability/unavailable_test.go:47`: the same gap as N2. Disposition: auto-fix. Confidence: 9.
- N5, `internal/capability/unavailable_test.go:47`: a trim or case mutant of the detail survives. No row grades the detail's spacing or case. Disposition: no-op. Confidence: 5.
- N6, `specs/shared-test-fixtures/spec.md:302`: no test drives two parallel callers into one skip log. `Unavailable` adds no write path, so the existing single-write transport contract holds. Disposition: no-op. Confidence: 6.

Repair targets: 2 (N1, and N2 with N4), both on ticket 09, in one combined repair charge.
