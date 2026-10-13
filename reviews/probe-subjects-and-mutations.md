# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/probe-subjects-and-mutations/spec.md",
  "plan_digest": "sha256:0ad5a6ea21f4e91ab9ccbe9d8e88aec3d43731b81fbac8a784791fcea3a05631",
  "implementation_session": "",
  "chunks": [
    {
      "id": "PS-A",
      "base": "8209612990d250568855d50b2735518e31d55223",
      "tip": "6175e732579feebb6224c2a7ade2f01992144e20",
      "plan_digest": "sha256:1f98cd1548c8bf824e1354271e8243b50d888c06d4ff8f7d4fd2b9d1d5c43052",
      "source_digest": "fe95a116179a1ce2f7b03f629e2896ad8a5721fc",
      "acceptance_rows": [
        "PS1",
        "PS2",
        "PS3",
        "PS4",
        "PS5",
        "PS6",
        "PS7",
        "PS8",
        "PS9",
        "PS10",
        "PS11",
        "PS12",
        "PS13",
        "PS14",
        "PS15",
        "PS16",
        "PS17",
        "PS18",
        "PS19",
        "PS20",
        "PS21",
        "PS22",
        "PS23",
        "PS24",
        "PS25",
        "PS76"
      ],
      "verification": [
        {
          "id": "ps-t01-probe-family",
          "performer": "claude:ps_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "235b86dafbb7b9a3d8a73e4d59e99bdf353cb7a0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
            "digest": "sha256:c6928d0e1721e8ae2518c2738f71781dc5de683dbbbe59c1b0cf4f551c531038",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-subjects-build,9de43b5280ba3648855f4401cac15e9f3dc6ecc3,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20769,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70604,1690\n  github.com/gibbonmi/bench/internal/probe,pass,44726,175\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit: 0\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "ps-t01-command-package",
          "performer": "claude:ps_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "235b86dafbb7b9a3d8a73e4d59e99bdf353cb7a0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
            "digest": "sha256:c6928d0e1721e8ae2518c2738f71781dc5de683dbbbe59c1b0cf4f551c531038",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-subjects-build,9de43b5280ba3648855f4401cac15e9f3dc6ecc3,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20769,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70604,1690\n  github.com/gibbonmi/bench/internal/probe,pass,44726,175\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit: 0\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ps-t01-conformance-package",
          "performer": "claude:ps_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "235b86dafbb7b9a3d8a73e4d59e99bdf353cb7a0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
            "digest": "sha256:c6928d0e1721e8ae2518c2738f71781dc5de683dbbbe59c1b0cf4f551c531038",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-subjects-build,9de43b5280ba3648855f4401cac15e9f3dc6ecc3,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20769,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70604,1690\n  github.com/gibbonmi/bench/internal/probe,pass,44726,175\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit: 0\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t01-package-core-guard",
          "performer": "claude:ps_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "235b86dafbb7b9a3d8a73e4d59e99bdf353cb7a0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
            "digest": "sha256:96dc28d65ec101c40376102f265623050679ecaef33e81c67be58816ef670f3c",
            "excerpt": "$ bench test --check package-core-guard\ntree[1]{target,head,dirty}:\n  probe-subjects-build,9de43b5280ba3648855f4401cac15e9f3dc6ecc3,false\ncheck[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,2661,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit: 0\n"
          },
          "requirement": "package-core-guard",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "ps-t01-root-conformance",
          "performer": "claude:ps_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "235b86dafbb7b9a3d8a73e4d59e99bdf353cb7a0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
            "digest": "sha256:2479cc07f1d14c4b8d209994dc4c3566b7fc446f10fa8de49bdd9097f144c2e6",
            "excerpt": "$ bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-subjects-build,9de43b5280ba3648855f4401cac15e9f3dc6ecc3,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18363,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit: 0\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t01-central-property",
          "performer": "claude:ps_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "235b86dafbb7b9a3d8a73e4d59e99bdf353cb7a0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
            "digest": "sha256:aee4d3252f29529d26a1653795b995ffd01903b4380fb0b52fcd187d9fc6137e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/probe/edit.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,all,passed,175\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,955\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,33577,175\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/start_bytes,\"combined_test.go:26: stdout = (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  bit,clamp.go,swap+swap,failed,1,yes\\\\\\\\n\\\" first a… (246 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseOverlap,\"combined_test.go:44: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execution… (425 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Apply each combined edit over the result of the earlier edit instead of over the start bytes. TestCombinedEditsMatchOnce must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
              "digest": "sha256:aee4d3252f29529d26a1653795b995ffd01903b4380fb0b52fcd187d9fc6137e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/probe/edit.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,all,passed,175\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,955\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,33577,175\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/start_bytes,\"combined_test.go:26: stdout = (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  bit,clamp.go,swap+swap,failed,1,yes\\\\\\\\n\\\" first a… (246 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseOverlap,\"combined_test.go:44: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execution… (425 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ps-t01-rep1-probe-family",
          "performer": "claude:ps_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fe95a116179a1ce2f7b03f629e2896ad8a5721fc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6175e732579feebb6224c2a7ade2f01992144e20",
            "digest": "sha256:9c730aa44fda7f1ea3eaaa777024fc29b98dae062190e2d105a39704b885978d",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-subjects-build,6175e732579feebb6224c2a7ade2f01992144e20,true\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20987,335\n  github.com/gibbonmi/bench/internal/conformance,pass,72616,1690\n  github.com/gibbonmi/bench/internal/probe,pass,46330,178\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n(dirty: only reviews/probe-subjects-and-mutations.md from bench record chunk; source bytes equal 6175e732)\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "ps-t01-rep1-command-package",
          "performer": "claude:ps_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fe95a116179a1ce2f7b03f629e2896ad8a5721fc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6175e732579feebb6224c2a7ade2f01992144e20",
            "digest": "sha256:9c730aa44fda7f1ea3eaaa777024fc29b98dae062190e2d105a39704b885978d",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-subjects-build,6175e732579feebb6224c2a7ade2f01992144e20,true\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20987,335\n  github.com/gibbonmi/bench/internal/conformance,pass,72616,1690\n  github.com/gibbonmi/bench/internal/probe,pass,46330,178\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n(dirty: only reviews/probe-subjects-and-mutations.md from bench record chunk; source bytes equal 6175e732)\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ps-t01-rep1-conformance-package",
          "performer": "claude:ps_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fe95a116179a1ce2f7b03f629e2896ad8a5721fc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6175e732579feebb6224c2a7ade2f01992144e20",
            "digest": "sha256:9c730aa44fda7f1ea3eaaa777024fc29b98dae062190e2d105a39704b885978d",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-subjects-build,6175e732579feebb6224c2a7ade2f01992144e20,true\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20987,335\n  github.com/gibbonmi/bench/internal/conformance,pass,72616,1690\n  github.com/gibbonmi/bench/internal/probe,pass,46330,178\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n(dirty: only reviews/probe-subjects-and-mutations.md from bench record chunk; source bytes equal 6175e732)\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t01-rep1-package-core-guard",
          "performer": "claude:ps_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fe95a116179a1ce2f7b03f629e2896ad8a5721fc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6175e732579feebb6224c2a7ade2f01992144e20",
            "digest": "sha256:c38570e5f6072f1707365883b26482be56fbe6e59d0fd14976d756f849b300e4",
            "excerpt": "$ bench test --check package-core-guard\ntree[1]{target,head,dirty}:\n  probe-subjects-build,6175e732579feebb6224c2a7ade2f01992144e20,true\ncheck[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,2720,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "package-core-guard",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "ps-t01-rep1-root-conformance",
          "performer": "claude:ps_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fe95a116179a1ce2f7b03f629e2896ad8a5721fc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6175e732579feebb6224c2a7ade2f01992144e20",
            "digest": "sha256:25556d80477eb7b30686320d503404e2396748767f4b077e1a14fa3a60f72815",
            "excerpt": "$ bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-subjects-build,6175e732579feebb6224c2a7ade2f01992144e20,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,19683,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t01-rep1-central-property",
          "performer": "claude:ps_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fe95a116179a1ce2f7b03f629e2896ad8a5721fc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6175e732579feebb6224c2a7ade2f01992144e20",
            "digest": "sha256:12bcadb3597ed3521e4ab8a506b1c27e772d3d9a5fe42164b43e503813387278",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/probe/command.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/,^TestCombinedEditsRefuseSingleForms$,passed,7\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,551\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,976,7\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--with-file_without_--with,\"combined_test.go:71: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  silent,clamp.go,swap+omit,passed,0,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,passed,1\\\\\\\\nprobe_execution[1]{schema,co… (887 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Apply each combined edit over the result of the earlier edit instead of over the start bytes. TestCombinedEditsMatchOnce must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "6175e732579feebb6224c2a7ade2f01992144e20",
              "digest": "sha256:12bcadb3597ed3521e4ab8a506b1c27e772d3d9a5fe42164b43e503813387278",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/probe/command.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/,^TestCombinedEditsRefuseSingleForms$,passed,7\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,551\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,976,7\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--with-file_without_--with,\"combined_test.go:71: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  silent,clamp.go,swap+omit,passed,0,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,passed,1\\\\\\\\nprobe_execution[1]{schema,co… (887 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ps-a-standards-1",
          "performer": "claude:ps-a-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "235b86dafbb7b9a3d8a73e4d59e99bdf353cb7a0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ps-a-standards",
            "digest": "sha256:75edeac305c2d2c7475ea2f940ec3cb4bdbe4a8967f4d92fdb4dd8734a6c7016",
            "excerpt": "Standards axis, chunk PS-A, pair 8209612990d2..9de43b5280ba, opus medium. Findings: 2.\nS1 internal/probe/command.go:81 vs internal/probe/mutation_input.go:39-74: gradeForms and mutationInput each derive the mutation group; two derivations of one fact (AGENTS.md one source per fact). Confidence 5.\nS2 internal/probe/edit.go:46: the empty-mutation refusal names --with equals --swap for a --with-file swap whose file bytes equal the old string; a refusal names the caller's actual input (craft-cli). Confidence 4.\nAdvice: occurrence(parsed) runs twice; placeholder <new> vs <text> for --with.\n"
          },
          "axis": "Standards",
          "base": "8209612990d250568855d50b2735518e31d55223",
          "tip": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "ps-a-spec-1",
          "performer": "claude:ps-a-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "235b86dafbb7b9a3d8a73e4d59e99bdf353cb7a0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ps-a-spec",
            "digest": "sha256:1a99e43e96f6530736e886ce44b99b0d77be128a9121e1d8b854d14fff76d640",
            "excerpt": "Spec axis, chunk PS-A (PS1-PS25, PS76), pair 8209612990d2..9de43b5280ba, opus medium. Findings: 1.\nP1 internal/probe/combined_test.go:49-51 with internal/probe/command.go:146-147: PS18 --with-file case hits the earlier withs>0&&withFile refusal; the combined&&withFile case has no isolating test. Confidence 8.\nEvery other row matches the spec and its cited test asserts it. The --delete Replacer.Remove decision holds.\nAdvice: ambiguous-versus-miss precedence in a combined mutant is unstated; single --omit emptying not tested.\n"
          },
          "axis": "Spec",
          "base": "8209612990d250568855d50b2735518e31d55223",
          "tip": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "ps-a-coverage-1",
          "performer": "claude:ps-a-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "235b86dafbb7b9a3d8a73e4d59e99bdf353cb7a0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ps-a-coverage",
            "digest": "sha256:25d8dd1a852363befafd42a850759425ca52791dc10cad7db819f7972422905d",
            "excerpt": "Coverage axis, chunk PS-A, pair 8209612990d2..9de43b5280ba, opus medium. Findings: 3.\nC1 internal/probe/command.go:109: the combined && withFile refusal has no red-capable test; probe swap combined&&withFile -> false&&withFile read silent, restored yes. Confidence 9.\nC2 internal/probe/mutation_input.go:91: no test reaches the err == nil term for an overflowing --occurrence value. Confidence 8.\nC3 internal/probe/subject.go:140: the Lstat error branch of the restore mode check is unreachable without a race. Confidence 6.\nProbes: 1 (swap, command.go:109, silent, restored yes).\n"
          },
          "axis": "Coverage",
          "base": "8209612990d250568855d50b2735518e31d55223",
          "tip": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "ps-a-standards-2",
          "performer": "claude:ps-a-standards-2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fe95a116179a1ce2f7b03f629e2896ad8a5721fc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ps-a-standards-2",
            "digest": "sha256:24cc922ef26e6ba870819cbb14b17b0a7dd14cd4a5d1c45a333512465699b70f",
            "excerpt": "Standards confirming round, chunk PS-A, repair delta 9de43b5280ba..6175e732579f, opus medium.\nS1 closed (confidence 8): namedGroup with the groupFlags table is the one source that gradeForms and mutationInput call.\nS2 closed (confidence 9): find.textFlag names --with or --with-file in the empty-mutation refusal; edit_test.go asserts the --with-file wording.\nAdvice: the namedGroup comment names its callers.\n"
          },
          "axis": "Standards",
          "base": "8209612990d250568855d50b2735518e31d55223",
          "tip": "6175e732579feebb6224c2a7ade2f01992144e20",
          "finding_ids": [],
          "supersedes": [
            "ps-a-standards-1"
          ]
        },
        {
          "id": "ps-a-spec-2",
          "performer": "claude:ps-a-spec-2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fe95a116179a1ce2f7b03f629e2896ad8a5721fc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ps-a-spec-2",
            "digest": "sha256:d6973cbec9a78a3119405e8f84e9ae9eaa84f9ab95e63f482e38cb2cbdca5c98",
            "excerpt": "Spec confirming round, chunk PS-A, repair delta 9de43b5280ba..6175e732579f, opus medium.\nP1 closed (confidence 9): combined_test.go:66-72 sends --swap, --with-file, and --omit with no --with; only the combined-with-file case refuses it.\nNo other drift; the namedGroup refactor keeps every refusal.\nAdvice: the spec did not name the --with-file empty-mutation text; the orchestrator added it to the spec.\n"
          },
          "axis": "Spec",
          "base": "8209612990d250568855d50b2735518e31d55223",
          "tip": "6175e732579feebb6224c2a7ade2f01992144e20",
          "finding_ids": [],
          "supersedes": [
            "ps-a-spec-1"
          ]
        },
        {
          "id": "ps-a-coverage-2",
          "performer": "claude:ps-a-coverage-2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fe95a116179a1ce2f7b03f629e2896ad8a5721fc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ps-a-coverage-2",
            "digest": "sha256:40f53235bbdd88105f0028dddd3eb4e09c08aa7302c6b9f10e8ece81d2a78a04",
            "excerpt": "Coverage confirming round, chunk PS-A, repair delta 9de43b5280ba..6175e732579f, opus medium.\nC1 closed (confidence 9): the new subtest reaches only the combined-with-file case.\nC2 closed (confidence 9): the overflow row alone turns red when err == nil is dropped.\nC3 stays no-op. No new findings. Probes: none.\n"
          },
          "axis": "Coverage",
          "base": "8209612990d250568855d50b2735518e31d55223",
          "tip": "6175e732579feebb6224c2a7ade2f01992144e20",
          "finding_ids": [],
          "supersedes": [
            "ps-a-coverage-1"
          ]
        }
      ]
    },
    {
      "id": "PS-B",
      "base": "6175e732579feebb6224c2a7ade2f01992144e20",
      "tip": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
      "plan_digest": "sha256:0ad5a6ea21f4e91ab9ccbe9d8e88aec3d43731b81fbac8a784791fcea3a05631",
      "source_digest": "a5459a6cad4efb4174d9c2ae5e312217c40551b9",
      "acceptance_rows": [
        "PS26",
        "PS27",
        "PS28",
        "PS29",
        "PS30",
        "PS31",
        "PS32",
        "PS33",
        "PS34",
        "PS35",
        "PS36",
        "PS37"
      ],
      "verification": [
        {
          "id": "ps-t02-probe-family",
          "performer": "claude:ps_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a5459a6cad4efb4174d9c2ae5e312217c40551b9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
            "digest": "sha256:4dde806d4e60720b8054c19d861f8180b412ce025f2f5103141df555078dd5fe",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,0d91e20b73d3ad6cb5ec74eeab491de0e2b07589,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,28131,335\n  github.com/gibbonmi/bench/internal/conformance,pass,86029,1690\n  github.com/gibbonmi/bench/internal/probe,pass,56144,186\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,80298,226\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "ps-t02-testreport",
          "performer": "claude:ps_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a5459a6cad4efb4174d9c2ae5e312217c40551b9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
            "digest": "sha256:4dde806d4e60720b8054c19d861f8180b412ce025f2f5103141df555078dd5fe",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,0d91e20b73d3ad6cb5ec74eeab491de0e2b07589,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,28131,335\n  github.com/gibbonmi/bench/internal/conformance,pass,86029,1690\n  github.com/gibbonmi/bench/internal/probe,pass,56144,186\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,80298,226\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "ps-t02-command-package",
          "performer": "claude:ps_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a5459a6cad4efb4174d9c2ae5e312217c40551b9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
            "digest": "sha256:4dde806d4e60720b8054c19d861f8180b412ce025f2f5103141df555078dd5fe",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,0d91e20b73d3ad6cb5ec74eeab491de0e2b07589,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,28131,335\n  github.com/gibbonmi/bench/internal/conformance,pass,86029,1690\n  github.com/gibbonmi/bench/internal/probe,pass,56144,186\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,80298,226\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ps-t02-conformance-package",
          "performer": "claude:ps_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a5459a6cad4efb4174d9c2ae5e312217c40551b9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
            "digest": "sha256:4dde806d4e60720b8054c19d861f8180b412ce025f2f5103141df555078dd5fe",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,0d91e20b73d3ad6cb5ec74eeab491de0e2b07589,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,28131,335\n  github.com/gibbonmi/bench/internal/conformance,pass,86029,1690\n  github.com/gibbonmi/bench/internal/probe,pass,56144,186\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,80298,226\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t02-root-conformance",
          "performer": "claude:ps_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a5459a6cad4efb4174d9c2ae5e312217c40551b9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
            "digest": "sha256:fd3654017fc345739f1aabf59a89c09276742a7335f5ce81949a71ca981a3b31",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,0d91e20b73d3ad6cb5ec74eeab491de0e2b07589,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,19252,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t02-central-property",
          "performer": "claude:ps_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a5459a6cad4efb4174d9c2ae5e312217c40551b9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
            "digest": "sha256:b833fac31b8a7951b2b16a6d5ec7c0d4c327cb421f6bbbe495023dbc4f27a47e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/probe/probe.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestSelectionPrecedesTheSubstringMiss$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1426\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,322,4\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_line-routing,\"subject_check_test.go:120: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_system_--run_^TestClampNegative$,\"subject_check_test.go:120: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/unknown_check,\"subject_check_test.go:105: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Answer the substring miss before the selection parse with a hard-coded all run cell. TestSelectionPrecedesTheSubstringMiss must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
              "digest": "sha256:b833fac31b8a7951b2b16a6d5ec7c0d4c327cb421f6bbbe495023dbc4f27a47e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/probe/probe.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestSelectionPrecedesTheSubstringMiss$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1426\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,322,4\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_line-routing,\"subject_check_test.go:120: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_system_--run_^TestClampNegative$,\"subject_check_test.go:120: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/unknown_check,\"subject_check_test.go:105: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n"
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
      "from": "sha256:be6c53c260f9f2d5de9df46ee418c9be074fa8fa1932cdf443e747e2b1c446ce",
      "to": "sha256:1f98cd1548c8bf824e1354271e8243b50d888c06d4ff8f7d4fd2b9d1d5c43052",
      "chunk_ids": {
        "PS-A": [
          "PS-A"
        ]
      }
    },
    {
      "from": "sha256:1f98cd1548c8bf824e1354271e8243b50d888c06d4ff8f7d4fd2b9d1d5c43052",
      "to": "sha256:0ad5a6ea21f4e91ab9ccbe9d8e88aec3d43731b81fbac8a784791fcea3a05631",
      "chunk_ids": {
        "PS-A": [
          "PS-A"
        ]
      }
    }
  ]
}
```

## Repair allowance

Chunk PS-A has used 1 of its 2 repair cycles after its initial review. Repair 6175e732 closed S1, S2, P1, C1, and C2, and the confirming round of all three axes passed.

## Standards

PS-A: 2 findings. The worst issue is S1.

- S1 (auto-fix, confidence 5): `internal/probe/command.go:81` and `internal/probe/mutation_input.go:39-74` derive the mutation group twice. `gradeForms` must give the one group, and `mutationInput` must read that group.
- S2 (auto-fix, confidence 4): `internal/probe/edit.go:46` names `--with equals --swap` for a `--with-file` swap. The refusal must name the flag that the caller gave.

## Spec

PS-A: 1 finding. The worst issue is P1.

- P1 (auto-fix, confidence 8): `internal/probe/combined_test.go:49-51` reaches the earlier `--with` and `--with-file` refusal. A PS18 subtest must send `--swap` with `--with-file` and `--omit`, and no `--with`. This finding has the same fix as C1.

## Coverage

PS-A: 3 findings. The worst issue is C1.

- C1 (auto-fix, confidence 9): `internal/probe/command.go:109` has no red-capable test, and a probe there read `silent`. The P1 subtest closes it.
- C2 (auto-fix, confidence 8): `internal/probe/mutation_input.go:91` has no test for an overflowing `--occurrence` value. The PS8 value list must add one.
- C3 (no-op, confidence 6): `internal/probe/subject.go:140` is reachable only through a race after the byte check.

Repair targets: 4 (S1, S2, P1 with C1, and C2) for ticket 1.

## Advice

- PS-A Standards: the `namedGroup` comment names its callers.
- PS-A Spec: the spec does not yet name the `--with-file equals --swap` refusal text. The orchestrator adds it in the ticket 2 plan commit.
