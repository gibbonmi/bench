# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/probe-subjects-and-mutations/spec.md",
  "plan_digest": "sha256:95540fdff85b79045b84441bbcf357fda6bcb980917f14f3acf6ea387f7ac7b9",
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
      "tip": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
      "plan_digest": "sha256:61cc4c138e8f1feee50a4ed368358a8048d738fd5c1e2e7b71fbadd2f7e4fb89",
      "source_digest": "f161a3a5e243a1a11d35a5cac635988d5c050c55",
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
        "PS37",
        "PS77"
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
        },
        {
          "id": "ps-t02-rep1-probe-family",
          "performer": "claude:ps_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a71101bcf99eeb2dfb29654c0a18b6ad3c774006",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d4f03f030be17feb46f5e64e709663f9d8979abb",
            "digest": "sha256:d62270910b8bf94bedaa87d85e51843e80a48d91ce4d2e9a2e70f3fff0f03b79",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,d4f03f030be17feb46f5e64e709663f9d8979abb,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,pass,67894,190\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "ps-t02-rep1-testreport",
          "performer": "claude:ps_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a71101bcf99eeb2dfb29654c0a18b6ad3c774006",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d4f03f030be17feb46f5e64e709663f9d8979abb",
            "digest": "sha256:bed032faf90ef3345188cabd073f1dc380ae9c6e8e06be37ac16744a3016284c",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,d4f03f030be17feb46f5e64e709663f9d8979abb,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,88234,226\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "ps-t02-rep1-command-package",
          "performer": "claude:ps_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a71101bcf99eeb2dfb29654c0a18b6ad3c774006",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d4f03f030be17feb46f5e64e709663f9d8979abb",
            "digest": "sha256:f21a5e2dce4818214c06d6448a258c11fc05d930e1c4f8aa8159188b6c3a80b1",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,d4f03f030be17feb46f5e64e709663f9d8979abb,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,35003,335\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ps-t02-rep1-conformance-package",
          "performer": "claude:ps_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a71101bcf99eeb2dfb29654c0a18b6ad3c774006",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d4f03f030be17feb46f5e64e709663f9d8979abb",
            "digest": "sha256:0590cb9d0338299373108d4dc37a0bdb2034b32b8a43a837f16be4329483479c",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,d4f03f030be17feb46f5e64e709663f9d8979abb,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,92746,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t02-rep1-root-conformance",
          "performer": "claude:ps_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a71101bcf99eeb2dfb29654c0a18b6ad3c774006",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d4f03f030be17feb46f5e64e709663f9d8979abb",
            "digest": "sha256:037b4bc3012695440262fbf727c4b6095c2f8b82c23aa20cc1c49f66293a6a30",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,d4f03f030be17feb46f5e64e709663f9d8979abb,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,19494,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t02-rep1-central-property",
          "performer": "claude:ps_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a71101bcf99eeb2dfb29654c0a18b6ad3c774006",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d4f03f030be17feb46f5e64e709663f9d8979abb",
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
              "ref": "d4f03f030be17feb46f5e64e709663f9d8979abb",
              "digest": "sha256:b833fac31b8a7951b2b16a6d5ec7c0d4c327cb421f6bbbe495023dbc4f27a47e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/probe/probe.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestSelectionPrecedesTheSubstringMiss$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1426\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,322,4\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_line-routing,\"subject_check_test.go:120: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_system_--run_^TestClampNegative$,\"subject_check_test.go:120: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/unknown_check,\"subject_check_test.go:105: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ps-t02-rep2-probe-family",
          "performer": "claude:ps_t02_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f161a3a5e243a1a11d35a5cac635988d5c050c55",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
            "digest": "sha256:2e7f96eed4e8a082cc05b992148d858fe1e41d241d75ed38169adc2cbf64865a",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,145d5439330c0a5bc7232b11d4be5aefb1dd47c5,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,pass,48136,190\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "ps-t02-rep2-testreport",
          "performer": "claude:ps_t02_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f161a3a5e243a1a11d35a5cac635988d5c050c55",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
            "digest": "sha256:3be13badf5a2b478ee088541c6af93b8756b363befb23ba372e9965da80d5a60",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,145d5439330c0a5bc7232b11d4be5aefb1dd47c5,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,72359,226\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "ps-t02-rep2-command-package",
          "performer": "claude:ps_t02_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f161a3a5e243a1a11d35a5cac635988d5c050c55",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
            "digest": "sha256:9a25a0804f8e95c90f1894a396ca3e0889108fa692c71cf438fa7fee6758e2de",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,145d5439330c0a5bc7232b11d4be5aefb1dd47c5,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,21897,335\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ps-t02-rep2-conformance-package",
          "performer": "claude:ps_t02_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f161a3a5e243a1a11d35a5cac635988d5c050c55",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
            "digest": "sha256:cb1e08f46aa9e954a7c87f319b272a344e5bed273140f33b18cf581584ff365e",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,145d5439330c0a5bc7232b11d4be5aefb1dd47c5,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,74727,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t02-rep2-root-conformance",
          "performer": "claude:ps_t02_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f161a3a5e243a1a11d35a5cac635988d5c050c55",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
            "digest": "sha256:7a805cc1442d45cc82730a9aff717efab6aa0923be06cbd9ea639e4e31f912ea",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,145d5439330c0a5bc7232b11d4be5aefb1dd47c5,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16677,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t02-rep2-central-property",
          "performer": "claude:ps_t02_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f161a3a5e243a1a11d35a5cac635988d5c050c55",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
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
              "ref": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
              "digest": "sha256:b833fac31b8a7951b2b16a6d5ec7c0d4c327cb421f6bbbe495023dbc4f27a47e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/probe/probe.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestSelectionPrecedesTheSubstringMiss$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1426\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,322,4\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_line-routing,\"subject_check_test.go:120: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (688 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_system_--run_^TestClampNegative$,\"subject_check_test.go:120: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (680 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/unknown_check,\"subject_check_test.go:105: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,\\\\\\\\\\\"\\\\\\\\\\\",all,\\\\\\\\\\\"\\\\\\\\\\\",0\\\\\\\\nprobe_execut… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ps-b-standards-1",
          "performer": "claude:ps-b-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a5459a6cad4efb4174d9c2ae5e312217c40551b9",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ps-b-standards",
            "digest": "sha256:80479e9d2c91d2a125a15eb3f6ac9c35a958105e0f7c13285db1ef9582ff804c",
            "excerpt": "Standards axis, chunk PS-B, pair 6175e732579f..0d91e20b73d3, opus medium. Findings: 2.\nS1 internal/probe/command.go:151 repeats the only-system-takes---run rule of internal/testreport/command.go:89 (AGENTS.md one source per fact); fix with one exported testreport predicate that both call. Confidence 7.\nS2 internal/testreport/selection_facts.go:95 calls every --check name a conformance check, but prose and system are not conformance checks (craft-cli: help states current behavior). Confidence 6.\nproseCheck deletion: no-op, confidence 9.\n"
          },
          "axis": "Standards",
          "base": "6175e732579feebb6224c2a7ade2f01992144e20",
          "tip": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "ps-b-spec-1",
          "performer": "claude:ps-b-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a5459a6cad4efb4174d9c2ae5e312217c40551b9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ps-b-spec",
            "digest": "sha256:b394a3a3501be6ba253f45b383fb856f0ae57d8cd7a27637c65bcc1eedcd7b4e",
            "excerpt": "Spec axis, chunk PS-B (PS26-PS37), pair 6175e732579f..0d91e20b73d3, opus medium. Findings: 0.\nEvery row matches the spec, and each cited test asserts it. The doc-comment refusal order, the prose exclusion fixture, and the retired route texts are confirmed.\n"
          },
          "axis": "Spec",
          "base": "6175e732579feebb6224c2a7ade2f01992144e20",
          "tip": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ps-b-coverage-1",
          "performer": "claude:ps-b-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a5459a6cad4efb4174d9c2ae5e312217c40551b9",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ps-b-coverage",
            "digest": "sha256:f9bd364ef22445de7e908d3940aca490117c44372556191f73b1258b13aa7ed0",
            "excerpt": "Coverage axis, chunk PS-B, pair 6175e732579f..0d91e20b73d3, opus medium. Findings: 2.\nC1 internal/probe/probe.go:79-94: a --check prose probe over a subject that the prose grade never reads answers a false silent; no row covers it. Confidence 7.\nC2 internal/probe/command.go:149: no test pins that --check system without --run is accepted. Confidence 5.\nProbes: none.\n"
          },
          "axis": "Coverage",
          "base": "6175e732579feebb6224c2a7ade2f01992144e20",
          "tip": "0d91e20b73d3ad6cb5ec74eeab491de0e2b07589",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "ps-b-standards-2",
          "performer": "claude:ps-b-standards-2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a71101bcf99eeb2dfb29654c0a18b6ad3c774006",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ps-b-standards-2",
            "digest": "sha256:35c2bbce59104d51831578cfb74cb616a11de37f03b02cc8723770a4bf29f404",
            "excerpt": "Standards confirming round, chunk PS-B, repair delta 0d91e20b73d3..d4f03f030be1, opus medium.\nS1 closed (9): testreport.CheckTakesRun owns the --run rule for three callers.\nS2 closed (9): the inventory note no longer names a conformance class.\nS4 open (auto-fix, 7): internal/probe/probe.go:107 hint 'grades only .md files' advertises the suffix rule of internal/prose/walk.go:27 (AGENTS.md enforcement and advertisement).\nS3 new (auto-fix, 5): the CheckTakesRun comment names its callers and misses Request.Run; the GradesProse comment describes prose.IsSubjectName, which it never consults.\n"
          },
          "axis": "Standards",
          "base": "6175e732579feebb6224c2a7ade2f01992144e20",
          "tip": "d4f03f030be17feb46f5e64e709663f9d8979abb",
          "finding_ids": [
            "S3",
            "S4"
          ],
          "supersedes": [
            "ps-b-standards-1"
          ]
        },
        {
          "id": "ps-b-spec-2",
          "performer": "claude:ps-b-spec-2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a71101bcf99eeb2dfb29654c0a18b6ad3c774006",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ps-b-spec-2",
            "digest": "sha256:191b7e517033f12d987d2de5e3c11a45413b4531debded569b731f322ff70d89",
            "excerpt": "Spec confirming round, chunk PS-B, repair delta 0d91e20b73d3..d4f03f030be1, opus medium.\nPS77 holds (9): order, exit 1, exact text, and the cited test assert it. PS26-PS37 still hold (8).\nRuling: the refusal-order summary sentence omits the not-graded step, a non-behavioral contradiction for the orchestrator to fix in the spec.\n"
          },
          "axis": "Spec",
          "base": "6175e732579feebb6224c2a7ade2f01992144e20",
          "tip": "d4f03f030be17feb46f5e64e709663f9d8979abb",
          "finding_ids": [],
          "supersedes": [
            "ps-b-spec-1"
          ]
        },
        {
          "id": "ps-b-coverage-2",
          "performer": "claude:ps-b-coverage-2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a71101bcf99eeb2dfb29654c0a18b6ad3c774006",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ps-b-coverage-2",
            "digest": "sha256:55ecfa92eefcd5e221027b0cce8b43d6b4af7a980eb3731f26b559a8d7c9a610",
            "excerpt": "Coverage confirming round, chunk PS-B, repair delta 0d91e20b73d3..d4f03f030be1, opus medium.\nC1 closed (8): TestProbeRefusesAnUngradedProseSubject reds without the guard in both subtests.\nC2 closed (7): the bare --check system row reds under a swap that requires --run.\nC3 (no-op): non-system --run refusal of bench test is pinned at internal/testreport/check_test.go:151 and :176 and usage_text_test.go:16. Probes: none.\n"
          },
          "axis": "Coverage",
          "base": "6175e732579feebb6224c2a7ade2f01992144e20",
          "tip": "d4f03f030be17feb46f5e64e709663f9d8979abb",
          "finding_ids": [],
          "supersedes": [
            "ps-b-coverage-1"
          ]
        },
        {
          "id": "ps-b-standards-3",
          "performer": "claude:ps-b-standards-3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f161a3a5e243a1a11d35a5cac635988d5c050c55",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ps-b-standards-3",
            "digest": "sha256:74b42820e23c4fffd2e64eca71eeb735feb4470946d7a7e81e8071cdd9b485a9",
            "excerpt": "Standards confirming round 2, chunk PS-B, repair delta d4f03f030be1..145d5439330c, opus medium.\nS4 closed (9): the hint at internal/probe/probe.go:107 reads <subject> is not a prose subject.\nS3 closed (9): the GradesProse and CheckTakesRun comments state the code without the prose rule or a caller list. No new blocking findings.\n"
          },
          "axis": "Standards",
          "base": "6175e732579feebb6224c2a7ade2f01992144e20",
          "tip": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
          "finding_ids": [],
          "supersedes": [
            "ps-b-standards-2"
          ]
        },
        {
          "id": "ps-b-spec-3",
          "performer": "claude:ps-b-spec-3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f161a3a5e243a1a11d35a5cac635988d5c050c55",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ps-b-spec-3",
            "digest": "sha256:ccef94fd91bb3c4d736d776c6128c929fb53973c7c551eeb875ffc44521b187b",
            "excerpt": "Spec confirming round 2, chunk PS-B, repair delta d4f03f030be1..145d5439330c, opus medium.\nHolds (9): the code prints the spec's not-graded refusal exactly, and the test pins it. PS77 and PS26-PS37 hold, and the refusal-order sentence matches the code.\n"
          },
          "axis": "Spec",
          "base": "6175e732579feebb6224c2a7ade2f01992144e20",
          "tip": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
          "finding_ids": [],
          "supersedes": [
            "ps-b-spec-2"
          ]
        },
        {
          "id": "ps-b-coverage-3",
          "performer": "claude:ps-b-coverage-3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f161a3a5e243a1a11d35a5cac635988d5c050c55",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ps-b-coverage-3",
            "digest": "sha256:34ffeda9383c9b72e4fc23fe73de34956e748fa0fc663bcdae82cc717a251216",
            "excerpt": "Coverage confirming round 2, chunk PS-B, repair delta d4f03f030be1..145d5439330c, opus medium.\nClosed: subject_check_test.go:51 pins the changed refusal text through an exact compare; the delta adds no production branch. Probes: none.\n"
          },
          "axis": "Coverage",
          "base": "6175e732579feebb6224c2a7ade2f01992144e20",
          "tip": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
          "finding_ids": [],
          "supersedes": [
            "ps-b-coverage-2"
          ]
        }
      ]
    },
    {
      "id": "PS-C",
      "base": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
      "tip": "7a55152eebda40993195f4279847529b1dc982c0",
      "plan_digest": "sha256:95540fdff85b79045b84441bbcf357fda6bcb980917f14f3acf6ea387f7ac7b9",
      "source_digest": "330ba04d2a41cf4c5c2663d5656760ea5a0d2c8c",
      "acceptance_rows": [
        "PS38",
        "PS39",
        "PS40",
        "PS41",
        "PS42",
        "PS43",
        "PS54",
        "PS55",
        "PS56",
        "PS57",
        "PS44",
        "PS45",
        "PS46",
        "PS47",
        "PS48",
        "PS49",
        "PS50",
        "PS51",
        "PS52",
        "PS53"
      ],
      "verification": [
        {
          "id": "ps-t04-t4-probe-family",
          "performer": "claude:ps_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:0962a943faf0d62f270ffee37c7c7132d811e3093e60e2f0e3b41f014aae314d",
            "excerpt": "tree: probe-subjects-build at ab67d6d580a238ccaf654ce2fcf73911ac4d381f (one PS-C floor call)\npackages[N]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,pass,53000,206\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "ps-t04-t4-testreport",
          "performer": "claude:ps_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:8b04b375d998c89d51ab7604e87097e829acab6c29095e517edab1077a59abd8",
            "excerpt": "tree: probe-subjects-build at ab67d6d580a238ccaf654ce2fcf73911ac4d381f (one PS-C floor call)\npackages[N]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,63425,232\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "ps-t04-t4-runbinary",
          "performer": "claude:ps_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:d4bdba307b63c6fcd7630ff9b142a3250d78643890143966c57dbf97f62774b7",
            "excerpt": "tree: probe-subjects-build at ab67d6d580a238ccaf654ce2fcf73911ac4d381f (one PS-C floor call)\npackages[N]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,pass,24713,29\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-runbinary",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "ps-t04-t4-record-forms",
          "performer": "claude:ps_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:f6f5ba2f22b0e2947a38ec63f4d9a4a3dd4171e5ba06749661374cfef3411286",
            "excerpt": "tree: probe-subjects-build at ab67d6d580a238ccaf654ce2fcf73911ac4d381f (one PS-C floor call)\npackages[N]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,12959,184\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-record-forms",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "ps-t04-t4-gate-readers",
          "performer": "claude:ps_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:cbb717b4dfcef150ad577bacb583b5ceb6a519337c9751c52e51774eeab432d6",
            "excerpt": "tree: probe-subjects-build at ab67d6d580a238ccaf654ce2fcf73911ac4d381f (one PS-C floor call)\npackages[N]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,38362,455\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-gate-readers",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "ps-t04-t4-command-package",
          "performer": "claude:ps_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:9ddb10da156990b4013234d9605dd676fda0243a3bfc39dcc504ec3c47e7b590",
            "excerpt": "tree: probe-subjects-build at ab67d6d580a238ccaf654ce2fcf73911ac4d381f (one PS-C floor call)\npackages[N]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18811,335\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ps-t04-t4-conformance-package",
          "performer": "claude:ps_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:a2cefb8427d3ff6905262ff87eaf6db5bba0ae83849f65b0d558f6b193401f3d",
            "excerpt": "tree: probe-subjects-build at ab67d6d580a238ccaf654ce2fcf73911ac4d381f (one PS-C floor call)\npackages[N]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,60394,1690\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "t4-conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t04-t4-root-conformance",
          "performer": "claude:ps_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:fe5a928f3bd22fae4db778c13727c5db6154581fc57a839b445b57562382dc3a",
            "excerpt": "tree: probe-subjects-build at ab67d6d580a238ccaf654ce2fcf73911ac4d381f\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15077,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t04-t4-central-property",
          "performer": "claude:ps_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:8c0740bcb885dffe8e517a4ac4e666bcaa383f205965d60fd361cc002a4a99b7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/outcome.go,swap,failed,3,TestExecuteClassifiesTheOutcome/refused,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,all,passed,232\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,973\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,55125,232\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestExecuteClassifiesTheOutcome/refused,\"outcome_test.go:73: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",1\n  github.com/gibbonmi/bench/internal/testreport,TestOnlyAScriptFailureIsBuildFailed/a_cancelled_build,\"outcome_test.go:161: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",2\n  github.com/gibbonmi/bench/internal/testreport,TestOnlyAScriptFailureIsBuildFailed/no_Go_on_PATH,\"outcome_test.go:168: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-central-property",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Map every builder error to build-failed instead of only a BuildError. TestOnlyAScriptFailureIsBuildFailed must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
              "digest": "sha256:8c0740bcb885dffe8e517a4ac4e666bcaa383f205965d60fd361cc002a4a99b7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/outcome.go,swap,failed,3,TestExecuteClassifiesTheOutcome/refused,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,all,passed,232\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,973\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,55125,232\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestExecuteClassifiesTheOutcome/refused,\"outcome_test.go:73: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",1\n  github.com/gibbonmi/bench/internal/testreport,TestOnlyAScriptFailureIsBuildFailed/a_cancelled_build,\"outcome_test.go:161: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",2\n  github.com/gibbonmi/bench/internal/testreport,TestOnlyAScriptFailureIsBuildFailed/no_Go_on_PATH,\"outcome_test.go:168: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ps-t03-t3-probe-family",
          "performer": "claude:ps_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:41465ac6920f2ee24c697123b6de1ee1532b2df4c521c8ff3a369222d3d83050",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,57e6bc12ce18e67cd531ebcd526e8fff92946832,false\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16972,335\n  github.com/gibbonmi/bench/internal/conformance,pass,58717,1690\n  github.com/gibbonmi/bench/internal/gate,pass,38569,455\n  github.com/gibbonmi/bench/internal/probe,pass,53751,206\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,13956,184\n  github.com/gibbonmi/bench/internal/runbinary,pass,24486,29\n  github.com/gibbonmi/bench/internal/testreport,pass,60460,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "ps-t03-t3-testreport",
          "performer": "claude:ps_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:41465ac6920f2ee24c697123b6de1ee1532b2df4c521c8ff3a369222d3d83050",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,57e6bc12ce18e67cd531ebcd526e8fff92946832,false\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16972,335\n  github.com/gibbonmi/bench/internal/conformance,pass,58717,1690\n  github.com/gibbonmi/bench/internal/gate,pass,38569,455\n  github.com/gibbonmi/bench/internal/probe,pass,53751,206\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,13956,184\n  github.com/gibbonmi/bench/internal/runbinary,pass,24486,29\n  github.com/gibbonmi/bench/internal/testreport,pass,60460,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "ps-t03-t3-runbinary",
          "performer": "claude:ps_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:41465ac6920f2ee24c697123b6de1ee1532b2df4c521c8ff3a369222d3d83050",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,57e6bc12ce18e67cd531ebcd526e8fff92946832,false\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16972,335\n  github.com/gibbonmi/bench/internal/conformance,pass,58717,1690\n  github.com/gibbonmi/bench/internal/gate,pass,38569,455\n  github.com/gibbonmi/bench/internal/probe,pass,53751,206\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,13956,184\n  github.com/gibbonmi/bench/internal/runbinary,pass,24486,29\n  github.com/gibbonmi/bench/internal/testreport,pass,60460,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-runbinary",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "ps-t03-t3-record-forms",
          "performer": "claude:ps_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:41465ac6920f2ee24c697123b6de1ee1532b2df4c521c8ff3a369222d3d83050",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,57e6bc12ce18e67cd531ebcd526e8fff92946832,false\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16972,335\n  github.com/gibbonmi/bench/internal/conformance,pass,58717,1690\n  github.com/gibbonmi/bench/internal/gate,pass,38569,455\n  github.com/gibbonmi/bench/internal/probe,pass,53751,206\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,13956,184\n  github.com/gibbonmi/bench/internal/runbinary,pass,24486,29\n  github.com/gibbonmi/bench/internal/testreport,pass,60460,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-record-forms",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "ps-t03-t3-gate-readers",
          "performer": "claude:ps_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:41465ac6920f2ee24c697123b6de1ee1532b2df4c521c8ff3a369222d3d83050",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,57e6bc12ce18e67cd531ebcd526e8fff92946832,false\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16972,335\n  github.com/gibbonmi/bench/internal/conformance,pass,58717,1690\n  github.com/gibbonmi/bench/internal/gate,pass,38569,455\n  github.com/gibbonmi/bench/internal/probe,pass,53751,206\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,13956,184\n  github.com/gibbonmi/bench/internal/runbinary,pass,24486,29\n  github.com/gibbonmi/bench/internal/testreport,pass,60460,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-gate-readers",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "ps-t03-t3-command-package",
          "performer": "claude:ps_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:41465ac6920f2ee24c697123b6de1ee1532b2df4c521c8ff3a369222d3d83050",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,57e6bc12ce18e67cd531ebcd526e8fff92946832,false\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16972,335\n  github.com/gibbonmi/bench/internal/conformance,pass,58717,1690\n  github.com/gibbonmi/bench/internal/gate,pass,38569,455\n  github.com/gibbonmi/bench/internal/probe,pass,53751,206\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,13956,184\n  github.com/gibbonmi/bench/internal/runbinary,pass,24486,29\n  github.com/gibbonmi/bench/internal/testreport,pass,60460,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ps-t03-t3-conformance-package",
          "performer": "claude:ps_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:41465ac6920f2ee24c697123b6de1ee1532b2df4c521c8ff3a369222d3d83050",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,57e6bc12ce18e67cd531ebcd526e8fff92946832,false\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16972,335\n  github.com/gibbonmi/bench/internal/conformance,pass,58717,1690\n  github.com/gibbonmi/bench/internal/gate,pass,38569,455\n  github.com/gibbonmi/bench/internal/probe,pass,53751,206\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,13956,184\n  github.com/gibbonmi/bench/internal/runbinary,pass,24486,29\n  github.com/gibbonmi/bench/internal/testreport,pass,60460,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t03-t3-root-conformance",
          "performer": "claude:ps_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:a531d706eb82484a1eb02210c84830efa1eae3f63272f9ef78a8ed5820fdb311",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,57e6bc12ce18e67cd531ebcd526e8fff92946832,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14897,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t03-t3-central-property",
          "performer": "claude:ps_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
            "digest": "sha256:0a423e60a0983c3c736f6018ca43f81b0722c541f90898b2f6e4aabdbff405ba",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/verdict.go,swap,failed,1,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestExpectRequiresTheNamedWitness$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,615\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,1200,6\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,\"verdict_test.go:57: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,expect-missed,1,TestClampNegative/zero,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^Tes… (737 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Match --expect by exact name only and drop the parent-and-slash match. TestExpectRequiresTheNamedWitness must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
              "digest": "sha256:0a423e60a0983c3c736f6018ca43f81b0722c541f90898b2f6e4aabdbff405ba",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/verdict.go,swap,failed,1,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestExpectRequiresTheNamedWitness$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,615\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,1200,6\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,\"verdict_test.go:57: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,expect-missed,1,TestClampNegative/zero,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^Tes… (737 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ps-t04-rep1-t4-probe-family",
          "performer": "claude:ps_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "330ba04d2a41cf4c5c2663d5656760ea5a0d2c8c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "7a55152eebda40993195f4279847529b1dc982c0",
            "digest": "sha256:8d003a3a9f44b87bf71e4c3110eb64075af1102b95e77896d67a1ee9560a3e0a",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,f6bbaada788285cbbb3f67dafab7502b5402c599,true\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17454,335\n  github.com/gibbonmi/bench/internal/conformance,pass,59540,1690\n  github.com/gibbonmi/bench/internal/gate,pass,39483,455\nspilled{lines=12,bytes=821,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/e7c800827d5fd2eed7e68741a4ec1720/1791864791403580194-87ee99c06a41de15.out}\n  github.com/gibbonmi/bench/internal/runbinary,pass,25602,29\n  github.com/gibbonmi/bench/internal/testreport,pass,62356,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "ps-t04-rep1-t4-testreport",
          "performer": "claude:ps_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "330ba04d2a41cf4c5c2663d5656760ea5a0d2c8c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "7a55152eebda40993195f4279847529b1dc982c0",
            "digest": "sha256:8d003a3a9f44b87bf71e4c3110eb64075af1102b95e77896d67a1ee9560a3e0a",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,f6bbaada788285cbbb3f67dafab7502b5402c599,true\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17454,335\n  github.com/gibbonmi/bench/internal/conformance,pass,59540,1690\n  github.com/gibbonmi/bench/internal/gate,pass,39483,455\nspilled{lines=12,bytes=821,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/e7c800827d5fd2eed7e68741a4ec1720/1791864791403580194-87ee99c06a41de15.out}\n  github.com/gibbonmi/bench/internal/runbinary,pass,25602,29\n  github.com/gibbonmi/bench/internal/testreport,pass,62356,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "ps-t04-rep1-t4-runbinary",
          "performer": "claude:ps_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "330ba04d2a41cf4c5c2663d5656760ea5a0d2c8c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "7a55152eebda40993195f4279847529b1dc982c0",
            "digest": "sha256:8d003a3a9f44b87bf71e4c3110eb64075af1102b95e77896d67a1ee9560a3e0a",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,f6bbaada788285cbbb3f67dafab7502b5402c599,true\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17454,335\n  github.com/gibbonmi/bench/internal/conformance,pass,59540,1690\n  github.com/gibbonmi/bench/internal/gate,pass,39483,455\nspilled{lines=12,bytes=821,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/e7c800827d5fd2eed7e68741a4ec1720/1791864791403580194-87ee99c06a41de15.out}\n  github.com/gibbonmi/bench/internal/runbinary,pass,25602,29\n  github.com/gibbonmi/bench/internal/testreport,pass,62356,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-runbinary",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "ps-t04-rep1-t4-record-forms",
          "performer": "claude:ps_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "330ba04d2a41cf4c5c2663d5656760ea5a0d2c8c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "7a55152eebda40993195f4279847529b1dc982c0",
            "digest": "sha256:8d003a3a9f44b87bf71e4c3110eb64075af1102b95e77896d67a1ee9560a3e0a",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,f6bbaada788285cbbb3f67dafab7502b5402c599,true\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17454,335\n  github.com/gibbonmi/bench/internal/conformance,pass,59540,1690\n  github.com/gibbonmi/bench/internal/gate,pass,39483,455\nspilled{lines=12,bytes=821,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/e7c800827d5fd2eed7e68741a4ec1720/1791864791403580194-87ee99c06a41de15.out}\n  github.com/gibbonmi/bench/internal/runbinary,pass,25602,29\n  github.com/gibbonmi/bench/internal/testreport,pass,62356,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-record-forms",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "ps-t04-rep1-t4-gate-readers",
          "performer": "claude:ps_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "330ba04d2a41cf4c5c2663d5656760ea5a0d2c8c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "7a55152eebda40993195f4279847529b1dc982c0",
            "digest": "sha256:8d003a3a9f44b87bf71e4c3110eb64075af1102b95e77896d67a1ee9560a3e0a",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,f6bbaada788285cbbb3f67dafab7502b5402c599,true\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17454,335\n  github.com/gibbonmi/bench/internal/conformance,pass,59540,1690\n  github.com/gibbonmi/bench/internal/gate,pass,39483,455\nspilled{lines=12,bytes=821,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/e7c800827d5fd2eed7e68741a4ec1720/1791864791403580194-87ee99c06a41de15.out}\n  github.com/gibbonmi/bench/internal/runbinary,pass,25602,29\n  github.com/gibbonmi/bench/internal/testreport,pass,62356,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-gate-readers",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "ps-t04-rep1-t4-command-package",
          "performer": "claude:ps_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "330ba04d2a41cf4c5c2663d5656760ea5a0d2c8c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "7a55152eebda40993195f4279847529b1dc982c0",
            "digest": "sha256:8d003a3a9f44b87bf71e4c3110eb64075af1102b95e77896d67a1ee9560a3e0a",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,f6bbaada788285cbbb3f67dafab7502b5402c599,true\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17454,335\n  github.com/gibbonmi/bench/internal/conformance,pass,59540,1690\n  github.com/gibbonmi/bench/internal/gate,pass,39483,455\nspilled{lines=12,bytes=821,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/e7c800827d5fd2eed7e68741a4ec1720/1791864791403580194-87ee99c06a41de15.out}\n  github.com/gibbonmi/bench/internal/runbinary,pass,25602,29\n  github.com/gibbonmi/bench/internal/testreport,pass,62356,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ps-t04-rep1-t4-conformance-package",
          "performer": "claude:ps_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "330ba04d2a41cf4c5c2663d5656760ea5a0d2c8c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "7a55152eebda40993195f4279847529b1dc982c0",
            "digest": "sha256:8d003a3a9f44b87bf71e4c3110eb64075af1102b95e77896d67a1ee9560a3e0a",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,f6bbaada788285cbbb3f67dafab7502b5402c599,true\npackages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17454,335\n  github.com/gibbonmi/bench/internal/conformance,pass,59540,1690\n  github.com/gibbonmi/bench/internal/gate,pass,39483,455\nspilled{lines=12,bytes=821,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/e7c800827d5fd2eed7e68741a4ec1720/1791864791403580194-87ee99c06a41de15.out}\n  github.com/gibbonmi/bench/internal/runbinary,pass,25602,29\n  github.com/gibbonmi/bench/internal/testreport,pass,62356,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t04-rep1-t4-root-conformance",
          "performer": "claude:ps_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "330ba04d2a41cf4c5c2663d5656760ea5a0d2c8c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "7a55152eebda40993195f4279847529b1dc982c0",
            "digest": "sha256:8beb59e2e86120562e1c4ccf9bd58a81c439b51c9b4fbe28b8fac224e68e7397",
            "excerpt": "tree[1]{target,head,dirty}:\n  probe-subjects-build,f6bbaada788285cbbb3f67dafab7502b5402c599,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15098,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ps-t04-rep1-t4-central-property",
          "performer": "claude:ps_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "330ba04d2a41cf4c5c2663d5656760ea5a0d2c8c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "7a55152eebda40993195f4279847529b1dc982c0",
            "digest": "sha256:8c0740bcb885dffe8e517a4ac4e666bcaa383f205965d60fd361cc002a4a99b7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/outcome.go,swap,failed,3,TestExecuteClassifiesTheOutcome/refused,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,all,passed,232\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,973\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,55125,232\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestExecuteClassifiesTheOutcome/refused,\"outcome_test.go:73: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",1\n  github.com/gibbonmi/bench/internal/testreport,TestOnlyAScriptFailureIsBuildFailed/a_cancelled_build,\"outcome_test.go:161: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",2\n  github.com/gibbonmi/bench/internal/testreport,TestOnlyAScriptFailureIsBuildFailed/no_Go_on_PATH,\"outcome_test.go:168: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-central-property",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Map every builder error to build-failed instead of only a BuildError. TestOnlyAScriptFailureIsBuildFailed must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "7a55152eebda40993195f4279847529b1dc982c0",
              "digest": "sha256:8c0740bcb885dffe8e517a4ac4e666bcaa383f205965d60fd361cc002a4a99b7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/outcome.go,swap,failed,3,TestExecuteClassifiesTheOutcome/refused,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,all,passed,232\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,973\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,55125,232\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestExecuteClassifiesTheOutcome/refused,\"outcome_test.go:73: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",1\n  github.com/gibbonmi/bench/internal/testreport,TestOnlyAScriptFailureIsBuildFailed/a_cancelled_build,\"outcome_test.go:161: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",2\n  github.com/gibbonmi/bench/internal/testreport,TestOnlyAScriptFailureIsBuildFailed/no_Go_on_PATH,\"outcome_test.go:168: outcome = {Kind:build-failed FailedTests:0 FailedNames:[] Ran:0 BuildLine:}, want {Kind:refused FailedTests:0 FailedNames:[] Ran:0 BuildLine:}\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ps-c-standards-1",
          "performer": "claude:ps-c-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ps-c-standards",
            "digest": "sha256:3874fbea981c6bfa36aeeb3f4db61ba307b50a0ee65093814a8407d00fe6d4c0",
            "excerpt": "Standards axis, chunk PS-C, pair 145d5439330c..ab67d6d580a2, opus medium. Findings: 2.\nS1 internal/testreport/outcome_test.go:323: sameOutcome lists every Outcome field by hand, so a new field goes unchecked in silence (AGENTS.md one source per fact). Confidence 6.\nS2 internal/probe/verdict_test.go:134-138: the build-script fixture install repeats buildEachRun in subject_check_fixture_test.go:30-35 (a fixture harness pasted twice). Confidence 5.\nConsumers of Factory.Own and probe.ReadResult: no defect.\n"
          },
          "axis": "Standards",
          "base": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
          "tip": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "ps-c-spec-1",
          "performer": "claude:ps-c-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ps-c-spec",
            "digest": "sha256:5ad3075f084ec110d6c5e09dfba2ea16e52ba1d9cdda9e273b1eb92cbee7150f",
            "excerpt": "Spec axis, chunk PS-C (PS38-PS57), pair 145d5439330c..ab67d6d580a2, opus medium. Findings: 0.\nAll 20 rows hold. Ruling: a red baseline names its first failure, consistent with the spec. Factory.Own readers print unchanged bytes; TestNativeRecordCLIJourney passes the v2 envelope without an edit.\nAdvice: A1 add a sentence that a red baseline names its own first failure; A2 no row for a mutated-run build-script failure table; A3 mutant table survives a restore failure.\n"
          },
          "axis": "Spec",
          "base": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
          "tip": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ps-c-coverage-1",
          "performer": "claude:ps-c-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f38889bea0cbab80653b4b790a85a3c55552d0ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ps-c-coverage",
            "digest": "sha256:4618b89c55c462f0015b9782bcf721f02115ecd3e08873af929c744728dca3a1",
            "excerpt": "Coverage axis, chunk PS-C, pair 145d5439330c..ab67d6d580a2, opus medium. Findings: 2.\nC1 internal/probe/verdict.go:45: no test runs --expect over a passing mutated run, so dropping the failed-kind guard stays green. Confidence 9.\nC2 internal/probe/verdict.go:49: no test puts a control byte in the first build line, so dropping sanitize.Preview stays green. Confidence 8.\nProbes: none.\n"
          },
          "axis": "Coverage",
          "base": "145d5439330c0a5bc7232b11d4be5aefb1dd47c5",
          "tip": "ab67d6d580a238ccaf654ce2fcf73911ac4d381f",
          "finding_ids": [
            "C1",
            "C2"
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
    },
    {
      "from": "sha256:0ad5a6ea21f4e91ab9ccbe9d8e88aec3d43731b81fbac8a784791fcea3a05631",
      "to": "sha256:6124b1fe814c7d1d1ddd8422233ef9ffed6cec824407d7140135e64c06f02c57",
      "chunk_ids": {
        "PS-A": [
          "PS-A"
        ],
        "PS-B": [
          "PS-B"
        ]
      }
    },
    {
      "from": "sha256:6124b1fe814c7d1d1ddd8422233ef9ffed6cec824407d7140135e64c06f02c57",
      "to": "sha256:61cc4c138e8f1feee50a4ed368358a8048d738fd5c1e2e7b71fbadd2f7e4fb89",
      "chunk_ids": {
        "PS-A": [
          "PS-A"
        ],
        "PS-B": [
          "PS-B"
        ]
      }
    },
    {
      "from": "sha256:61cc4c138e8f1feee50a4ed368358a8048d738fd5c1e2e7b71fbadd2f7e4fb89",
      "to": "sha256:eda3f737cedfbbb29511d8b04197c36b17a332d75a896cb01b29bc4bf2acf8e4",
      "chunk_ids": {
        "PS-A": [
          "PS-A"
        ],
        "PS-B": [
          "PS-B"
        ]
      }
    },
    {
      "from": "sha256:eda3f737cedfbbb29511d8b04197c36b17a332d75a896cb01b29bc4bf2acf8e4",
      "to": "sha256:95540fdff85b79045b84441bbcf357fda6bcb980917f14f3acf6ea387f7ac7b9",
      "chunk_ids": {
        "PS-A": [
          "PS-A"
        ],
        "PS-B": [
          "PS-B"
        ],
        "PS-C": [
          "PS-C"
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

## PS-B pickup

Chunk PS-B has used 2 of its 2 repair cycles after its initial review. Repair d4f03f03 closed S1, S2, C1, and C2. Repair 145d5439 closed S3 and S4, and the second confirming round of all three axes passed.

### Standards

PS-B: 2 findings. The worst issue is S1.

- S1 (auto-fix, confidence 7): `internal/probe/command.go:151` repeats the rule of `internal/testreport/command.go:89` that only the system suite takes `--run`. One exported `testreport` predicate must own the rule, and both graders must call it.
- S2 (auto-fix, confidence 6): `internal/testreport/selection_facts.go:95` calls each `--check` name a conformance check. The note must name a check from the `bench test --help` inventory without that class.

### Spec

PS-B: 0 findings.

### Coverage

PS-B: 2 findings. The worst issue is C1.

- C1 (auto-fix after a plan expansion, confidence 7): `internal/probe/probe.go:79-94` answers `silent` for a `--check prose` probe of a subject that the prose grade never reads. The probe must refuse a subject that `prose.IsSubjectName` rejects.
- C2 (auto-fix, confidence 5): `internal/probe/command.go:149` has no row that accepts `--check system` without `--run`.

Repair targets: 4 (S1, S2, C1, and C2) for ticket 2.

### PS-B confirming round

- S3 (auto-fix, confidence 5): the `CheckTakesRun` comment in `internal/testreport/selection_facts.go` names its callers and misses `Request.Run`. The `GradesProse` comment describes `prose.IsSubjectName`, which that method never consults.
- S4 (auto-fix, confidence 7): the not-graded hint at `internal/probe/probe.go:107` states the `.md` rule that `prose.IsSubjectName` enforces. The hint must say that the subject is not a prose subject, without the rule.
- C3 (no-op): `internal/testreport/check_test.go:151` and `:176` already pin the non-system `--run` refusal of `bench test`.

## PS-C pickup

Chunk PS-C has used 0 of its 2 repair cycles after its initial review.

### Standards

PS-C: 2 findings. The worst issue is S1.

- S1 (auto-fix, ticket 3, confidence 6): `sameOutcome` in `internal/testreport/outcome_test.go` lists every `Outcome` field by hand, so a new field goes unchecked. The helper must compare the whole value.
- S2 (auto-fix, ticket 4, confidence 5): `internal/probe/verdict_test.go:134-138` repeats the build-script install of `buildEachRun`. One fixture helper must own it.

### Spec

PS-C: 0 findings. The orchestrator adds the advice A1 sentence on a red baseline to the Verdict paragraph.

### Coverage

PS-C: 2 findings. The worst issue is C1.

- C1 (auto-fix, ticket 3, confidence 9): no test runs `--expect` over a passing mutated run, so the failed-kind guard at `internal/probe/verdict.go:45` has no red-capable test.
- C2 (auto-fix, ticket 4, confidence 8): no test puts a control byte in the first build line, so `sanitize.Preview` at `internal/probe/verdict.go:49` has no red-capable test.

Repair targets: 4. Ticket 3 takes S1 and C1, and ticket 4 takes S2 and C2.
