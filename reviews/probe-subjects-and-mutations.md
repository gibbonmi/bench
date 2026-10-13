# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/probe-subjects-and-mutations/spec.md",
  "plan_digest": "sha256:be6c53c260f9f2d5de9df46ee418c9be074fa8fa1932cdf443e747e2b1c446ce",
  "implementation_session": "",
  "chunks": [
    {
      "id": "PS-A",
      "base": "8209612990d250568855d50b2735518e31d55223",
      "tip": "9de43b5280ba3648855f4401cac15e9f3dc6ecc3",
      "plan_digest": "sha256:be6c53c260f9f2d5de9df46ee418c9be074fa8fa1932cdf443e747e2b1c446ce",
      "source_digest": "235b86dafbb7b9a3d8a73e4d59e99bdf353cb7a0",
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
