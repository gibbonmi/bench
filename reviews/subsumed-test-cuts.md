# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/subsumed-test-cuts/spec.md",
  "plan_digest": "sha256:231c4c21fcd42b76e165ed5bf163f33aa4482652664ac995293048a25c1501b1",
  "implementation_session": "",
  "chunks": [
    {
      "id": "FT411-C1",
      "base": "931b7ea78a972b7f45af2e7bfcfd7ddfc40c3a77",
      "tip": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
      "plan_digest": "sha256:0d0f369cac0d4a134ad41fa09c025f63069706c7ace97e6b2b9fd35cead1e997",
      "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
      "acceptance_rows": [
        "TC2",
        "TC3",
        "TC4",
        "TC5",
        "TC6",
        "TC7",
        "TC8",
        "TC9",
        "TC43",
        "TC1",
        "TC10",
        "TC11",
        "TC12",
        "TC13",
        "TC14",
        "TC15",
        "TC16",
        "TC17",
        "TC18",
        "TC19",
        "TC45",
        "TC20",
        "TC21",
        "TC22",
        "TC23",
        "TC24",
        "TC41",
        "TC42",
        "TC44"
      ],
      "verification": [
        {
          "id": "t01-tc2",
          "performer": "ft411-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "889f110fe1cf439429acfc0fdbabe38d22884033",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-01",
            "digest": "sha256:0cf6ef3048c78e1fa49026640bea1f468c92cf58000497d0c5bedbe778c70be3",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompatibilityUndo,1,1,0,null,yes\n  2,bit,swap,failed,1,TestCompatibilityUndo,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityUndo$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "3ea94e85180f56a97c537aadde2e074377bba1e1",
          "requirement": "tc2",
          "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityUndo$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/transaction/publication.go: swap `remove = os.Remove` for `remove = func(string) error { return nil }`; expect TestCompatibilityUndo",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-01",
              "digest": "sha256:0cf6ef3048c78e1fa49026640bea1f468c92cf58000497d0c5bedbe778c70be3",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompatibilityUndo,1,1,0,null,yes\n  2,bit,swap,failed,1,TestCompatibilityUndo,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityUndo$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t01-tc3",
          "performer": "ft411-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "889f110fe1cf439429acfc0fdbabe38d22884033",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-01",
            "digest": "sha256:0cf6ef3048c78e1fa49026640bea1f468c92cf58000497d0c5bedbe778c70be3",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompatibilityUndo,1,1,0,null,yes\n  2,bit,swap,failed,1,TestCompatibilityUndo,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityUndo$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "3ea94e85180f56a97c537aadde2e074377bba1e1",
          "requirement": "tc3",
          "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityUndo$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/transaction/publication.go: swap `_, err = f.Write(value.Data)` for a write of a copy whose first byte has bit 0x20 flipped; expect TestCompatibilityUndo",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-01",
              "digest": "sha256:0cf6ef3048c78e1fa49026640bea1f468c92cf58000497d0c5bedbe778c70be3",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompatibilityUndo,1,1,0,null,yes\n  2,bit,swap,failed,1,TestCompatibilityUndo,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityUndo$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t01-tc4",
          "performer": "ft411-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "889f110fe1cf439429acfc0fdbabe38d22884033",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-01",
            "digest": "sha256:3c171df95ce6ca03500eaf57727409b68b4db7133597a051c757a7706dbbf8d9",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  2,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  3,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityRecoveryBoundaries$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "3ea94e85180f56a97c537aadde2e074377bba1e1",
          "requirement": "tc4",
          "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityRecoveryBoundaries$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/compatibility/recovery.go: swap `{\"security-policy\", \"decision\",` for `{\"security-policy\", \"info\",`; expect TestCompatibilityRecoveryBoundaries",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-01",
              "digest": "sha256:3c171df95ce6ca03500eaf57727409b68b4db7133597a051c757a7706dbbf8d9",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  2,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  3,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityRecoveryBoundaries$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t01-tc5",
          "performer": "ft411-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "889f110fe1cf439429acfc0fdbabe38d22884033",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-01",
            "digest": "sha256:3c171df95ce6ca03500eaf57727409b68b4db7133597a051c757a7706dbbf8d9",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  2,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  3,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityRecoveryBoundaries$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "3ea94e85180f56a97c537aadde2e074377bba1e1",
          "requirement": "tc5",
          "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityRecoveryBoundaries$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/compatibility/recovery.go: swap `{\"process-interruption\", \"decision\",` for `{\"process-interruption\", \"info\",`; expect TestCompatibilityRecoveryBoundaries",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-01",
              "digest": "sha256:3c171df95ce6ca03500eaf57727409b68b4db7133597a051c757a7706dbbf8d9",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  2,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  3,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityRecoveryBoundaries$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t01-tc6",
          "performer": "ft411-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "889f110fe1cf439429acfc0fdbabe38d22884033",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-01",
            "digest": "sha256:3c171df95ce6ca03500eaf57727409b68b4db7133597a051c757a7706dbbf8d9",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  2,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  3,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityRecoveryBoundaries$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "3ea94e85180f56a97c537aadde2e074377bba1e1",
          "requirement": "tc6",
          "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityRecoveryBoundaries$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/compatibility/recovery.go: swap `{\"private-runtime\", \"decision\",` for `{\"private-runtime\", \"info\",`; expect TestCompatibilityRecoveryBoundaries",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-01",
              "digest": "sha256:3c171df95ce6ca03500eaf57727409b68b4db7133597a051c757a7706dbbf8d9",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  2,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  3,bit,swap,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityRecoveryBoundaries$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t01-tc7",
          "performer": "ft411-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "889f110fe1cf439429acfc0fdbabe38d22884033",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-01",
            "digest": "sha256:c13905323c56a8badcb11a87f506b6e9129c95b452e1a3d2a3ff6eef1bdc4a3f",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  2,bit,insert,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityRecoveryBoundaries$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "3ea94e85180f56a97c537aadde2e074377bba1e1",
          "requirement": "tc7",
          "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityRecoveryBoundaries$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/compatibility.go: insert a write of `x` to `$CODEX_HOME/config.toml` before `namespace := filepath.Join(home, \"compatibility-repairs\")`; expect TestCompatibilityRecoveryBoundaries",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-01",
              "digest": "sha256:c13905323c56a8badcb11a87f506b6e9129c95b452e1a3d2a3ff6eef1bdc4a3f",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  2,bit,insert,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityRecoveryBoundaries$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t01-tc43",
          "performer": "ft411-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "889f110fe1cf439429acfc0fdbabe38d22884033",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-01",
            "digest": "sha256:c13905323c56a8badcb11a87f506b6e9129c95b452e1a3d2a3ff6eef1bdc4a3f",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  2,bit,insert,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityRecoveryBoundaries$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "3ea94e85180f56a97c537aadde2e074377bba1e1",
          "requirement": "tc43",
          "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityRecoveryBoundaries$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/compatibility.go: insert a write of `x` to `$CODEX_HOME/trust.json` before `namespace := filepath.Join(home, \"compatibility-repairs\")`; expect TestCompatibilityRecoveryBoundaries",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-01",
              "digest": "sha256:c13905323c56a8badcb11a87f506b6e9129c95b452e1a3d2a3ff6eef1bdc4a3f",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\n  2,bit,insert,failed,1,TestCompatibilityRecoveryBoundaries,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt/repairtest,^TestCompatibilityRecoveryBoundaries$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc1",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:a58008107f8d27c0a4a54bab749485bb5f4f4cebc676d7748f8979ec6e4e3baa",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDoctorParsesHelpAndUnknownFlags,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestDoctorParsesHelpAndUnknownFlags$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc1",
          "command": "bench test --package ./internal/adopt --run '^TestDoctorParsesHelpAndUnknownFlags$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/compatibility.go: swap `args[0] == \"--compat\"` for `args[0] != \"--compat\"`; expect TestDoctorParsesHelpAndUnknownFlags",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:a58008107f8d27c0a4a54bab749485bb5f4f4cebc676d7748f8979ec6e4e3baa",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDoctorParsesHelpAndUnknownFlags,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestDoctorParsesHelpAndUnknownFlags$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc10",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc10",
          "command": "bench test --package ./internal/adopt --run '^TestSetupGateRejectsIgnoredDeclaredInputs$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/setup.go: omit `BENCH_CONFORMANCE_CONSUMER_ONLY=1 `; expect TestSetupGateRejectsIgnoredDeclaredInputs",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc11",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc11",
          "command": "bench test --package ./internal/adopt --run '^TestSetupSeedsGateInputsInZeroSignalRepo$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/setup.go: omit `\"BENCH_RUN_BINARY\", `; expect TestSetupSeedsGateInputsInZeroSignalRepo",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc12",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc12",
          "command": "bench test --package ./internal/adopt --run '^TestSetupTwiceLeavesGateInputsIdentical$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/setup.go: swap `content: scaffoldGateInputs()` for `content: scaffoldGateInputs() + map[bool]string{false: \" \"}[facts.zeroSignal]`; expect TestSetupTwiceLeavesGateInputsIdentical",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc13",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
            "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc13",
          "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: swap `; exists {` for `; exists && false {`; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
              "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc14",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
            "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc14",
          "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: omit ` || asset.Fingerprint == \"\"`; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
              "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc15",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
            "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc15",
          "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: insert an `empty` refusal for an input with no current and no desired asset before `decisions := make([]LifecycleDecision, 0, len(current)+len(desired))`; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
              "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc16",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
            "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc16",
          "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: swap `case lifecycleAdd, lifecycleChange, lifecycleRemove:` for `case lifecycleAdd, lifecycleChange:`; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
              "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc17",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
            "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc17",
          "command": "bench test --package ./internal/adopt --run '^TestLifecycleDecisionOrder$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: swap `return decisions[i].Path < decisions[j].Path` for `return false && decisions[i].Path < decisions[j].Path`; expect TestLifecycleDecisionOrder",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
              "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc45",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
            "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc45",
          "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: swap `if !validAssetPath(asset.Path) || asset.Fingerprint == \"\" {` for `if asset.Fingerprint == \"\" {`; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:c1546a17bc3587dac957c35dd734e3b498b6897c4404539964ccfb7b0ec6aa08",
              "excerpt": "mutants[6]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,19,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,19,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/legacy_empty,19,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,19,1,0,null,yes\n  5,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\n  6,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$|^TestLifecycleDecisionOrder$,passed,7\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t02-tc19",
          "performer": "ft411-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d21a19f8c8ac45e6dbcc126b99631692454a52ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-02",
            "digest": "sha256:b19b267bb7ba9fa70b70f055b858390c9a808a98e6bd10c6a3686047903938ff",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestPrePushHookIgnoresALegacyPin,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestPrePushHookIgnoresALegacyPin$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "8e98a19c5aa647a9cbc11fbbb81a417f860dac2e",
          "requirement": "tc19",
          "command": "bench test --package ./internal/adopt --run '^TestPrePushHookIgnoresALegacyPin$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/prepush.sh: insert `echo \"topic push\" >&2` before `exit 0`; expect TestPrePushHookIgnoresALegacyPin",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-02",
              "digest": "sha256:b19b267bb7ba9fa70b70f055b858390c9a808a98e6bd10c6a3686047903938ff",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestPrePushHookIgnoresALegacyPin,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestPrePushHookIgnoresALegacyPin$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t03-tc20",
          "performer": "ft411-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-03",
            "digest": "sha256:66212da928d5ece8b5818788844655efac7d98b3b61bec1bbee9d396b30a3f1e",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestLifecycleUpgradePreviewCompatibility/plan_kinds/adapter_unowned,17,1,0,null,yes\n  2,bit,swap,failed,6,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\n  3,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/dropped_row,17,1,0,null,yes\n  4,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "requirement": "tc20",
          "command": "bench test --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/upgrade.go: swap `if fingerprint == \"\" {` for `if false {`; expect TestLifecycleUpgradePreviewCompatibility",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-03",
              "digest": "sha256:66212da928d5ece8b5818788844655efac7d98b3b61bec1bbee9d396b30a3f1e",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestLifecycleUpgradePreviewCompatibility/plan_kinds/adapter_unowned,17,1,0,null,yes\n  2,bit,swap,failed,6,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\n  3,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/dropped_row,17,1,0,null,yes\n  4,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t03-tc21",
          "performer": "ft411-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-03",
            "digest": "sha256:66212da928d5ece8b5818788844655efac7d98b3b61bec1bbee9d396b30a3f1e",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestLifecycleUpgradePreviewCompatibility/plan_kinds/adapter_unowned,17,1,0,null,yes\n  2,bit,swap,failed,6,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\n  3,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/dropped_row,17,1,0,null,yes\n  4,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "requirement": "tc21",
          "command": "bench test --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/upgrade.go: swap `fingerprint := old` for `fingerprint := entry.kind + \":new\" + old[:0]`; expect TestLifecycleUpgradePreviewCompatibility",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-03",
              "digest": "sha256:66212da928d5ece8b5818788844655efac7d98b3b61bec1bbee9d396b30a3f1e",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestLifecycleUpgradePreviewCompatibility/plan_kinds/adapter_unowned,17,1,0,null,yes\n  2,bit,swap,failed,6,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\n  3,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/dropped_row,17,1,0,null,yes\n  4,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t03-tc22",
          "performer": "ft411-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-03",
            "digest": "sha256:66212da928d5ece8b5818788844655efac7d98b3b61bec1bbee9d396b30a3f1e",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestLifecycleUpgradePreviewCompatibility/plan_kinds/adapter_unowned,17,1,0,null,yes\n  2,bit,swap,failed,6,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\n  3,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/dropped_row,17,1,0,null,yes\n  4,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "requirement": "tc22",
          "command": "bench test --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/upgrade.go: swap `removed++` for `added++`; expect TestLifecycleUpgradePreviewCompatibility",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-03",
              "digest": "sha256:66212da928d5ece8b5818788844655efac7d98b3b61bec1bbee9d396b30a3f1e",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestLifecycleUpgradePreviewCompatibility/plan_kinds/adapter_unowned,17,1,0,null,yes\n  2,bit,swap,failed,6,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\n  3,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/dropped_row,17,1,0,null,yes\n  4,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t03-tc23",
          "performer": "ft411-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-03",
            "digest": "sha256:66212da928d5ece8b5818788844655efac7d98b3b61bec1bbee9d396b30a3f1e",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestLifecycleUpgradePreviewCompatibility/plan_kinds/adapter_unowned,17,1,0,null,yes\n  2,bit,swap,failed,6,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\n  3,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/dropped_row,17,1,0,null,yes\n  4,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "requirement": "tc23",
          "command": "bench test --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/upgrade.go: swap `Preserve: []string{\"CLAUDE.md\"}` for `Preserve: nil`; expect TestLifecycleUpgradePreviewCompatibility",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-03",
              "digest": "sha256:66212da928d5ece8b5818788844655efac7d98b3b61bec1bbee9d396b30a3f1e",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestLifecycleUpgradePreviewCompatibility/plan_kinds/adapter_unowned,17,1,0,null,yes\n  2,bit,swap,failed,6,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\n  3,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/dropped_row,17,1,0,null,yes\n  4,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/plan_kinds/CLAUDE.md_row,17,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t03-tc24",
          "performer": "ft411-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-03",
            "digest": "sha256:0e5518b92eaf2b9127e8dc1e8e7f1a33861fcb998a8cf585feb09618b8969cf8",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/hook_refresh,2,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$/^hook_refresh$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "requirement": "tc24",
          "command": "bench test --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$/^hook_refresh$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/upgrade.go: swap `direction != 0 && prePushRefreshEligible(root)` for `false && direction != 0 && prePushRefreshEligible(root)`; expect TestLifecycleUpgradePreviewCompatibility",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-03",
              "digest": "sha256:0e5518b92eaf2b9127e8dc1e8e7f1a33861fcb998a8cf585feb09618b8969cf8",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLifecycleUpgradePreviewCompatibility/hook_refresh,2,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$/^hook_refresh$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t03-tc41",
          "performer": "ft411-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-03",
            "digest": "sha256:66035fcd07c4af9e278f0afe649b10120d3419e9a4b4f1437b1393e25dae55be",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDetectedProjectGateRejectsIgnoredDeclaredInput,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  check,system,^TestDetectedProjectGateRejectsIgnoredDeclaredInput$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "requirement": "tc41",
          "command": "bench test --check system --run '^TestDetectedProjectGateRejectsIgnoredDeclaredInput$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/setup.go: swap `+ consumerGateHygieneCheck() + facts.gateCommand` for `+ facts.gateCommand`; expect TestDetectedProjectGateRejectsIgnoredDeclaredInput",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-03",
              "digest": "sha256:66035fcd07c4af9e278f0afe649b10120d3419e9a4b4f1437b1393e25dae55be",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDetectedProjectGateRejectsIgnoredDeclaredInput,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  check,system,^TestDetectedProjectGateRejectsIgnoredDeclaredInput$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t03-tc42",
          "performer": "ft411-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-author-03",
            "digest": "sha256:7faadc4e526bd8fd24bd74af549476233b6daf5be4f8cecda50d6f1da68d02c2",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLifecycleLinkPartitions/dangling_leaf_collision,2,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleLinkPartitions$/^dangling_leaf_collision$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "requirement": "tc42",
          "command": "bench test --package ./internal/adopt --run '^TestLifecycleLinkPartitions$/^dangling_leaf_collision$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/link_stage.go: swap `destInfo, err := os.Lstat(dest)` for `destInfo, err := os.Stat(dest)`; expect TestLifecycleLinkPartitions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-author-03",
              "digest": "sha256:7faadc4e526bd8fd24bd74af549476233b6daf5be4f8cecda50d6f1da68d02c2",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLifecycleLinkPartitions/dangling_leaf_collision,2,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleLinkPartitions$/^dangling_leaf_collision$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc1",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:a58008107f8d27c0a4a54bab749485bb5f4f4cebc676d7748f8979ec6e4e3baa",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDoctorParsesHelpAndUnknownFlags,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestDoctorParsesHelpAndUnknownFlags$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc1",
          "command": "bench test --package ./internal/adopt --run '^TestDoctorParsesHelpAndUnknownFlags$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/compatibility.go: swap `args[0] == \"--compat\"` for `args[0] != \"--compat\"`; expect TestDoctorParsesHelpAndUnknownFlags",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:a58008107f8d27c0a4a54bab749485bb5f4f4cebc676d7748f8979ec6e4e3baa",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDoctorParsesHelpAndUnknownFlags,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestDoctorParsesHelpAndUnknownFlags$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc10",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc10",
          "command": "bench test --package ./internal/adopt --run '^TestSetupGateRejectsIgnoredDeclaredInputs$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/setup.go: omit `BENCH_CONFORMANCE_CONSUMER_ONLY=1 `; expect TestSetupGateRejectsIgnoredDeclaredInputs",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc11",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc11",
          "command": "bench test --package ./internal/adopt --run '^TestSetupSeedsGateInputsInZeroSignalRepo$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/setup.go: omit `\"BENCH_RUN_BINARY\", `; expect TestSetupSeedsGateInputsInZeroSignalRepo",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc12",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc12",
          "command": "bench test --package ./internal/adopt --run '^TestSetupTwiceLeavesGateInputsIdentical$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/setup.go: swap `content: scaffoldGateInputs()` for `content: scaffoldGateInputs() + map[bool]string{false: \" \"}[facts.zeroSignal]`; expect TestSetupTwiceLeavesGateInputsIdentical",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:73a58af80747b1d2ce1e741d60b0992091b74202d2fdca49e5c820fbb9247077",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestSetupGateRejectsIgnoredDeclaredInputs/undeclared,4,1,0,null,yes\n  2,bit,omit,failed,1,TestSetupSeedsGateInputsInZeroSignalRepo,1,1,0,null,yes\n  3,bit,swap,failed,1,TestSetupTwiceLeavesGateInputsIdentical,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestSetupGateRejectsIgnoredDeclaredInputs$|^TestSetupSeedsGateInputsInZeroSignalRepo$|^TestSetupTwiceLeavesGateInputsIdentical$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc13",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:04781eb931d71761d236a82414702933a8e514aa8056cd8bee315b7cfb38565d",
            "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,18,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,18,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/present_empty_execution,18,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,18,1,0,null,yes\n  5,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc13",
          "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: swap `; exists {` for `; exists && false {`; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:04781eb931d71761d236a82414702933a8e514aa8056cd8bee315b7cfb38565d",
              "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,18,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,18,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/present_empty_execution,18,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,18,1,0,null,yes\n  5,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc14",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:04781eb931d71761d236a82414702933a8e514aa8056cd8bee315b7cfb38565d",
            "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,18,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,18,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/present_empty_execution,18,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,18,1,0,null,yes\n  5,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc14",
          "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: omit ` || asset.Fingerprint == \"\"`; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:04781eb931d71761d236a82414702933a8e514aa8056cd8bee315b7cfb38565d",
              "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,18,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,18,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/present_empty_execution,18,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,18,1,0,null,yes\n  5,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc15",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:04781eb931d71761d236a82414702933a8e514aa8056cd8bee315b7cfb38565d",
            "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,18,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,18,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/present_empty_execution,18,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,18,1,0,null,yes\n  5,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc15",
          "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: insert an `empty` refusal for an input with no current and no desired asset before `decisions := make([]LifecycleDecision, 0, len(current)+len(desired))`; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:04781eb931d71761d236a82414702933a8e514aa8056cd8bee315b7cfb38565d",
              "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,18,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,18,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/present_empty_execution,18,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,18,1,0,null,yes\n  5,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc16",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:04781eb931d71761d236a82414702933a8e514aa8056cd8bee315b7cfb38565d",
            "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,18,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,18,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/present_empty_execution,18,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,18,1,0,null,yes\n  5,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc16",
          "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: swap `case lifecycleAdd, lifecycleChange, lifecycleRemove:` for `case lifecycleAdd, lifecycleChange:`; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:04781eb931d71761d236a82414702933a8e514aa8056cd8bee315b7cfb38565d",
              "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,18,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,18,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/present_empty_execution,18,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,18,1,0,null,yes\n  5,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc17",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:c3ab7614591ad3f12a68b53216c0b6a4973f777d63c182e9e82b8ee398b2ac0d",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleDecisionOrder$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc17",
          "command": "bench test --package ./internal/adopt --run '^TestLifecycleDecisionOrder$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: swap `return decisions[i].Path < decisions[j].Path` for `return false && decisions[i].Path < decisions[j].Path`; expect TestLifecycleDecisionOrder",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:c3ab7614591ad3f12a68b53216c0b6a4973f777d63c182e9e82b8ee398b2ac0d",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestLifecycleDecisionOrder/execution,4,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestLifecycleDecisionOrder$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc45",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:04781eb931d71761d236a82414702933a8e514aa8056cd8bee315b7cfb38565d",
            "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,18,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,18,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/present_empty_execution,18,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,18,1,0,null,yes\n  5,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc45",
          "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/decision.go: swap `if !validAssetPath(asset.Path) || asset.Fingerprint == \"\" {` for `if asset.Fingerprint == \"\" {`; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:04781eb931d71761d236a82414702933a8e514aa8056cd8bee315b7cfb38565d",
              "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestCompleteLifecycleDecisions/legacy_duplicate,18,1,0,null,yes\n  2,bit,omit,failed,1,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,18,1,0,null,yes\n  3,bit,insert,failed,1,TestCompleteLifecycleDecisions/present_empty_execution,18,1,0,null,yes\n  4,bit,swap,failed,3,TestCompleteLifecycleDecisions/legacy,18,1,0,null,yes\n  5,bit,swap,failed,2,TestCompleteLifecycleDecisions/current_before_desired,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "r02-tc19",
          "performer": "ft411-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
            "digest": "sha256:b19b267bb7ba9fa70b70f055b858390c9a808a98e6bd10c6a3686047903938ff",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestPrePushHookIgnoresALegacyPin,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestPrePushHookIgnoresALegacyPin$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "requirement": "tc19",
          "command": "bench test --package ./internal/adopt --run '^TestPrePushHookIgnoresALegacyPin$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/prepush.sh: insert `echo \"topic push\" >&2` before `exit 0`; expect TestPrePushHookIgnoresALegacyPin",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft411-repair-02",
              "digest": "sha256:b19b267bb7ba9fa70b70f055b858390c9a808a98e6bd10c6a3686047903938ff",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,1,TestPrePushHookIgnoresALegacyPin,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/adopt,^TestPrePushHookIgnoresALegacyPin$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c1-standards",
          "performer": "ft411-review-c1-standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:c1-standards",
            "digest": "sha256:30968bfc40656ac0eb491a3bb661a0e43cf36bf6c98d669a4479454a31bdc836",
            "excerpt": "Standards axis FT411-C1 (931b7ea78..c299acc0b), sonnet/high. 2 advice, 0 blocking.\nF1 advice auto-fix conf 6 internal/adopt/lifecycle_decision_test.go:197,200: refusal-table check loosened from nil to len(...) != 0 for the legacy empty row; restore the nil check on refusal rows.\nF2 advice auto-fix conf 5 internal/adopt/setup_test.go:76: doc comment \"pins SD4.\" is stale after the T1-9 move adds the exact seeded-bytes check.\nChecked clean: merged boundary assertions, helpers, census key, fence.\n"
          },
          "axis": "Standards",
          "base": "931b7ea78a972b7f45af2e7bfcfd7ddfc40c3a77",
          "tip": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "finding_ids": [
            "c1-f1-nil-check",
            "c1-f2-sd4-comment"
          ],
          "supersedes": []
        },
        {
          "id": "c1-spec",
          "performer": "ft411-review-c1-spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:c1-spec",
            "digest": "sha256:d6a0ee12139c36a42ccc63e812341bb342d1b8252012c65b26f096e4b0ec2e8b",
            "excerpt": "Spec axis FT411-C1 (931b7ea78..c299acc0b), sonnet/high. 1 advice, 0 blocking.\nF1 advice conf 4 internal/adopt/lifecycle_decision_test.go:~200: shared refusal assertion loosened from nil to len(...) != 0 for the legacy empty row; keep the nil check where want != \"\".\nTC8 pass, TC44 pass, T1-5 bytes and modes pass, T1-9 and T1-10 moves pass, LTE cells pass, planned cells pass, all 25 recorded probes bit.\n"
          },
          "axis": "Spec",
          "base": "931b7ea78a972b7f45af2e7bfcfd7ddfc40c3a77",
          "tip": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "finding_ids": [
            "c1-f1-nil-check"
          ],
          "supersedes": []
        },
        {
          "id": "c1-coverage",
          "performer": "ft411-review-c1-coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:c1-coverage",
            "digest": "sha256:1c6494bd8305b2330a8e5d61b45dc1f3c37bad9e0f52c7c4dbff6cec1251cfd6",
            "excerpt": "Coverage axis FT411-C1 (931b7ea78..c299acc0b), sonnet/high. 2 advice, 0 blocking. No probe run.\nF1 advice ask-user conf 4 internal/adopt/lifecycle_decision_test.go:200: nil check weakened to admit the legacy empty success row; keep nil for refusal rows and give legacy empty its own assertion.\nF3 advice no-op conf 3 spec TC3 why-clause: kept TestCompatibilityUndo bit on the byte flip through undo divergence, not a repaired-bytes assertion; probe tc3 bit.\nChecked: TC45 escape row held, preview cuts held by plan kinds, topic push held, merge predicates kept, no dead helper. git status clean.\n"
          },
          "axis": "Coverage",
          "base": "931b7ea78a972b7f45af2e7bfcfd7ddfc40c3a77",
          "tip": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "finding_ids": [
            "c1-f1-nil-check"
          ],
          "supersedes": []
        }
      ],
      "confirmations": [
        {
          "id": "c1-confirm-1",
          "performer": "session_01JK5CPYv9wmkqVE69kw7W9b",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "medium",
          "source_digest": "97df05fc3a2c6e7a91ec99118964b973f19c71fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-code:orchestrator:c1-confirm-1",
            "digest": "sha256:7eb9ed070376884fd957654693db1209c2fd286f8b8d61da0901413440f4be14",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/adopt/decision.go,swap,failed,13,TestCompleteLifecycleDecisions/current_before_desired,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,18\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,4822\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,18\nfailures[13]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/current_before_desired,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed asset \\\\\\\\\\\"../c\\\\\\\\\\\" is invalid\\\"}, want the refusal \\\"managed asset \\\\\\\\\\\"../c\\\\\\\\\\\" is invalid\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/desired_before_preserve,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed asset \\\\\\\\\\\"../d\\\\\\\\\\\" is invalid\\\"}, want the refusal \\\"managed asset \\\\\\\\\\\"../d\\\\\\\\\\\" is invalid\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/dropped_escape,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"refused manifest path: ../escape\\\"}, want the refusal \\\"refused manifest path: ../escape\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/dropped_special,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"refused non-regular dropped asset: dir\\\"}, want the refusal \\\"refused non-regular dropped as… (249 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/dropped_unreadable,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"dropped asset \\\\\\\\\\\"unread\\\\\\\\\\\" has no fingerprint\\\"}, want the refusal \\\"dropped asset \\\\\\\\\\\"unread\\\\\\\\\\\" … (259 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/duplicate,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed asset \\\\\\\\\\\"twice\\\\\\\\\\\" is duplicated\\\"}, want the refusal \\\"managed asset \\\\\\\\\\\"twice\\\\\\\\\\\" is dupl… (247 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/foreign_adapter,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"conflict: .claude/commands/x.md has a symlink parent directory\\\"}, want the refusal \\\"confli… (297 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/invalid_preserve,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed path \\\\\\\\\\\"../x\\\\\\\\\\\" is invalid\\\"}, want the refusal \\\"managed path \\\\\\\\\\\"../x\\\\\\\\\\\" is invalid\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/legacy_duplicate,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed asset \\\\\\\\\\\"x\\\\\\\\\\\" is duplicated\\\"}, want the refusal \\\"managed asset \\\\\\\\\\\"x\\\\\\\\\\\" is duplicated\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed asset \\\\\\\\\\\"x\\\\\\\\\\\" is invalid\\\"}, want the refusal \\\"managed asset \\\\\\\\\\\"x\\\\\\\\\\\" is invalid\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/seed_symlink_parent,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"conflict: linked/seed has a symlink parent directory\\\"}, want the refusal \\\"conflict: linked… (277 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/symlink_parent_write,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"conflict: linked/asset has a symlink parent directory\\\"}, want the refusal \\\"conflict: linke… (279 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/unknown_action,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"lifecycle action \\\\\\\\\\\"relink\\\\\\\\\\\" is invalid\\\"}, want the refusal \\\"lifecycle action \\\\\\\\\\\"relink\\\\\\\\\\\" is… (249 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "base": "931b7ea78a972b7f45af2e7bfcfd7ddfc40c3a77",
          "reviewed_tip": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
          "repaired_tip": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
          "finding_ids": [
            "c1-f1-nil-check",
            "c1-f2-sd4-comment"
          ],
          "supersedes": [],
          "probe": {
            "mutation": "internal/adopt/decision.go: swap 'return LifecyclePlan{Refusal: refusal}' for a refusal return with empty non-nil Decisions and Operations; expect TestCompleteLifecycleDecisions",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-code:orchestrator:c1-confirm-1",
              "digest": "sha256:7eb9ed070376884fd957654693db1209c2fd286f8b8d61da0901413440f4be14",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/adopt/decision.go,swap,failed,13,TestCompleteLifecycleDecisions/current_before_desired,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,18\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,4822\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,18\nfailures[13]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/current_before_desired,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed asset \\\\\\\\\\\"../c\\\\\\\\\\\" is invalid\\\"}, want the refusal \\\"managed asset \\\\\\\\\\\"../c\\\\\\\\\\\" is invalid\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/desired_before_preserve,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed asset \\\\\\\\\\\"../d\\\\\\\\\\\" is invalid\\\"}, want the refusal \\\"managed asset \\\\\\\\\\\"../d\\\\\\\\\\\" is invalid\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/dropped_escape,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"refused manifest path: ../escape\\\"}, want the refusal \\\"refused manifest path: ../escape\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/dropped_special,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"refused non-regular dropped asset: dir\\\"}, want the refusal \\\"refused non-regular dropped as… (249 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/dropped_unreadable,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"dropped asset \\\\\\\\\\\"unread\\\\\\\\\\\" has no fingerprint\\\"}, want the refusal \\\"dropped asset \\\\\\\\\\\"unread\\\\\\\\\\\" … (259 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/duplicate,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed asset \\\\\\\\\\\"twice\\\\\\\\\\\" is duplicated\\\"}, want the refusal \\\"managed asset \\\\\\\\\\\"twice\\\\\\\\\\\" is dupl… (247 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/foreign_adapter,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"conflict: .claude/commands/x.md has a symlink parent directory\\\"}, want the refusal \\\"confli… (297 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/invalid_preserve,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed path \\\\\\\\\\\"../x\\\\\\\\\\\" is invalid\\\"}, want the refusal \\\"managed path \\\\\\\\\\\"../x\\\\\\\\\\\" is invalid\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/legacy_duplicate,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed asset \\\\\\\\\\\"x\\\\\\\\\\\" is duplicated\\\"}, want the refusal \\\"managed asset \\\\\\\\\\\"x\\\\\\\\\\\" is duplicated\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/legacy_missing_fingerprint,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"managed asset \\\\\\\\\\\"x\\\\\\\\\\\" is invalid\\\"}, want the refusal \\\"managed asset \\\\\\\\\\\"x\\\\\\\\\\\" is invalid\\\"\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/seed_symlink_parent,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"conflict: linked/seed has a symlink parent directory\\\"}, want the refusal \\\"conflict: linked… (277 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/symlink_parent_write,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"conflict: linked/asset has a symlink parent directory\\\"}, want the refusal \\\"conflict: linke… (279 bytes)\",1\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/unknown_action,\"lifecycle_decision_test.go:203: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation{}, Refusal:\\\"lifecycle action \\\\\\\\\\\"relink\\\\\\\\\\\" is invalid\\\"}, want the refusal \\\"lifecycle action \\\\\\\\\\\"relink\\\\\\\\\\\" is… (249 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ]
    },
    {
      "id": "FT411-C2",
      "base": "824802ac05e7eb4f070b5bdf86506172bec49eb3",
      "tip": "d466c480239e31d45e37c79a2240f6145afb7b19",
      "plan_digest": "sha256:231c4c21fcd42b76e165ed5bf163f33aa4482652664ac995293048a25c1501b1",
      "source_digest": "e886fe61608d9babaa528e30c0c603aad6cbb2f8",
      "acceptance_rows": [
        "TC25",
        "TC26",
        "TC29",
        "TC30",
        "TC31",
        "TC32",
        "TC33",
        "TC27",
        "TC28",
        "TC34",
        "TC35",
        "TC36",
        "TC37",
        "TC38",
        "TC39",
        "TC40"
      ],
      "verification": [
        {
          "id": "t04-tc25",
          "performer": "ft411-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e886fe61608d9babaa528e30c0c603aad6cbb2f8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft411-author-04",
            "digest": "sha256:085cb72c8b924ab0a040f5e6be994c7d80738660e40db790c45d81f5b69d428a",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandCommandPrintsNoHeadsLineWithoutRecords,1,1,0,null,yes\n  2,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\n  3,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandCommandPrintsNoHeadsLineWithoutRecords$|^TestLandCommandStatesTheCensusCountAndDropsTheRecords$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "d466c480239e31d45e37c79a2240f6145afb7b19",
          "requirement": "tc25",
          "command": "bench test --package ./internal/worktree --run '^TestLandCommandPrintsNoHeadsLineWithoutRecords$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/worktree/land.go: swap `return counts[assignment]` for `return counts[assignment] + 1`; expect TestLandCommandPrintsNoHeadsLineWithoutRecords",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft411-author-04",
              "digest": "sha256:085cb72c8b924ab0a040f5e6be994c7d80738660e40db790c45d81f5b69d428a",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandCommandPrintsNoHeadsLineWithoutRecords,1,1,0,null,yes\n  2,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\n  3,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandCommandPrintsNoHeadsLineWithoutRecords$|^TestLandCommandStatesTheCensusCountAndDropsTheRecords$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t04-tc26",
          "performer": "ft411-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e886fe61608d9babaa528e30c0c603aad6cbb2f8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft411-author-04",
            "digest": "sha256:b947daf0b0c4350b20188752c6855d119d8ca4e42f9ce5e89d2372d0321d0ea6",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandSkipsTheRefreshWithoutBuildInputs,1,1,0,null,yes\n  2,bit,swap,failed,1,TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandSkipsTheRefreshWithoutBuildInputs$|^TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "d466c480239e31d45e37c79a2240f6145afb7b19",
          "requirement": "tc26",
          "command": "bench test --package ./internal/worktree --run '^TestLandSkipsTheRefreshWithoutBuildInputs$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/worktree/land_effects.go: swap `return landedComplete(stdout, result, active, records)` for a call that keeps the exit code, prints one trailer line after the landed record, and returns the code; expect TestLandSkipsTheRefreshWithoutBuildInputs",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft411-author-04",
              "digest": "sha256:b947daf0b0c4350b20188752c6855d119d8ca4e42f9ce5e89d2372d0321d0ea6",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandSkipsTheRefreshWithoutBuildInputs,1,1,0,null,yes\n  2,bit,swap,failed,1,TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandSkipsTheRefreshWithoutBuildInputs$|^TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t04-tc29",
          "performer": "ft411-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e886fe61608d9babaa528e30c0c603aad6cbb2f8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft411-author-04",
            "digest": "sha256:b947daf0b0c4350b20188752c6855d119d8ca4e42f9ce5e89d2372d0321d0ea6",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandSkipsTheRefreshWithoutBuildInputs,1,1,0,null,yes\n  2,bit,swap,failed,1,TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandSkipsTheRefreshWithoutBuildInputs$|^TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "d466c480239e31d45e37c79a2240f6145afb7b19",
          "requirement": "tc29",
          "command": "bench test --package ./internal/worktree --run '^TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/worktree/land_effects.go: swap `if j.releaseLandingAssignment(j, a, root, []string{\"--request\", request, path}, io.Discard, &diagnostic) == 0 {` for `if true {`; expect TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft411-author-04",
              "digest": "sha256:b947daf0b0c4350b20188752c6855d119d8ca4e42f9ce5e89d2372d0321d0ea6",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandSkipsTheRefreshWithoutBuildInputs,1,1,0,null,yes\n  2,bit,swap,failed,1,TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandSkipsTheRefreshWithoutBuildInputs$|^TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion$,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t04-tc30",
          "performer": "ft411-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e886fe61608d9babaa528e30c0c603aad6cbb2f8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft411-author-04",
            "digest": "sha256:085cb72c8b924ab0a040f5e6be994c7d80738660e40db790c45d81f5b69d428a",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandCommandPrintsNoHeadsLineWithoutRecords,1,1,0,null,yes\n  2,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\n  3,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandCommandPrintsNoHeadsLineWithoutRecords$|^TestLandCommandStatesTheCensusCountAndDropsTheRecords$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "d466c480239e31d45e37c79a2240f6145afb7b19",
          "requirement": "tc30",
          "command": "bench test --package ./internal/worktree --run '^TestLandCommandStatesTheCensusCountAndDropsTheRecords$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/worktree/land.go: omit `fmt.Fprintf(stderr, \"census heads{%s}\\n\", breakdown)`; expect TestLandCommandStatesTheCensusCountAndDropsTheRecords",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft411-author-04",
              "digest": "sha256:085cb72c8b924ab0a040f5e6be994c7d80738660e40db790c45d81f5b69d428a",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandCommandPrintsNoHeadsLineWithoutRecords,1,1,0,null,yes\n  2,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\n  3,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandCommandPrintsNoHeadsLineWithoutRecords$|^TestLandCommandStatesTheCensusCountAndDropsTheRecords$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t04-tc31",
          "performer": "ft411-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e886fe61608d9babaa528e30c0c603aad6cbb2f8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft411-author-04",
            "digest": "sha256:085cb72c8b924ab0a040f5e6be994c7d80738660e40db790c45d81f5b69d428a",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandCommandPrintsNoHeadsLineWithoutRecords,1,1,0,null,yes\n  2,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\n  3,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandCommandPrintsNoHeadsLineWithoutRecords$|^TestLandCommandStatesTheCensusCountAndDropsTheRecords$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "d466c480239e31d45e37c79a2240f6145afb7b19",
          "requirement": "tc31",
          "command": "bench test --package ./internal/worktree --run '^TestLandCommandStatesTheCensusCountAndDropsTheRecords$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/worktree/land.go: omit `fmt.Fprintf(stderr, \"census output{%s}\\n\", breakdown)`; expect TestLandCommandStatesTheCensusCountAndDropsTheRecords",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft411-author-04",
              "digest": "sha256:085cb72c8b924ab0a040f5e6be994c7d80738660e40db790c45d81f5b69d428a",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandCommandPrintsNoHeadsLineWithoutRecords,1,1,0,null,yes\n  2,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\n  3,bit,omit,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandCommandPrintsNoHeadsLineWithoutRecords$|^TestLandCommandStatesTheCensusCountAndDropsTheRecords$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t04-tc32",
          "performer": "ft411-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e886fe61608d9babaa528e30c0c603aad6cbb2f8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft411-author-04",
            "digest": "sha256:8633b423c8c84f6a485fa122601a2f27a6f79dddabf8c66cccd551dedbaddbab",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandCommandStatesTheCensusCountAndDropsTheRecords$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "d466c480239e31d45e37c79a2240f6145afb7b19",
          "requirement": "tc32",
          "command": "bench test --package ./internal/worktree --run '^TestLandCommandStatesTheCensusCountAndDropsTheRecords$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/worktree/lifecycle.go: swap `_ = census.Drop(a.home, root, id)` for `_ = census.Drop`; expect TestLandCommandStatesTheCensusCountAndDropsTheRecords",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft411-author-04",
              "digest": "sha256:8633b423c8c84f6a485fa122601a2f27a6f79dddabf8c66cccd551dedbaddbab",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandCommandStatesTheCensusCountAndDropsTheRecords$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "t04-tc33",
          "performer": "ft411-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e886fe61608d9babaa528e30c0c603aad6cbb2f8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft411-author-04",
            "digest": "sha256:8633b423c8c84f6a485fa122601a2f27a6f79dddabf8c66cccd551dedbaddbab",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandCommandStatesTheCensusCountAndDropsTheRecords$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "d466c480239e31d45e37c79a2240f6145afb7b19",
          "requirement": "tc33",
          "command": "bench test --package ./internal/worktree --run '^TestLandCommandStatesTheCensusCountAndDropsTheRecords$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/census/output.go: swap `assignment+outputSuffix` for `assignment` in the RecordOutput write; expect TestLandCommandStatesTheCensusCountAndDropsTheRecords",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft411-author-04",
              "digest": "sha256:8633b423c8c84f6a485fa122601a2f27a6f79dddabf8c66cccd551dedbaddbab",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestLandCommandStatesTheCensusCountAndDropsTheRecords,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,^TestLandCommandStatesTheCensusCountAndDropsTheRecords$,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
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
      "from": "sha256:4576078b7e43a8211933e7577e092a7faf6f5a6880d23a56bdca340592b59c5f",
      "to": "sha256:68132c138c14f1f05b79f12799c0bd19247a27c7ae9d3b7c8652c19e23b0850e",
      "chunk_ids": {
        "FT411-C1": [
          "FT411-C1"
        ]
      }
    },
    {
      "from": "sha256:68132c138c14f1f05b79f12799c0bd19247a27c7ae9d3b7c8652c19e23b0850e",
      "to": "sha256:e1564e017bb221e293307e9ca521767f8831bc4731280ac8f2de1aeef747bfca",
      "chunk_ids": {
        "FT411-C1": [
          "FT411-C1"
        ]
      }
    },
    {
      "from": "sha256:e1564e017bb221e293307e9ca521767f8831bc4731280ac8f2de1aeef747bfca",
      "to": "sha256:0d0f369cac0d4a134ad41fa09c025f63069706c7ace97e6b2b9fd35cead1e997",
      "chunk_ids": {
        "FT411-C1": [
          "FT411-C1"
        ]
      }
    },
    {
      "from": "sha256:0d0f369cac0d4a134ad41fa09c025f63069706c7ace97e6b2b9fd35cead1e997",
      "to": "sha256:231c4c21fcd42b76e165ed5bf163f33aa4482652664ac995293048a25c1501b1",
      "chunk_ids": {
        "FT411-C1": [
          "FT411-C1"
        ]
      }
    }
  ]
}
```

## FT411-C1 review findings

Frozen pair: 931b7ea78 to c299acc0b. Each axis ran on sonnet at high effort.
The raw count is 5 findings, and the repair-target count is 2.

### Standards

The axis reports 2 findings, and the worst is F1.

- F1 `c1-f1-nil-check`, auto-fix, confidence 6: `internal/adopt/lifecycle_decision_test.go:200` loosens the refusal-row check from a nil check to a length check for the `legacy empty` row. The repair keeps the nil check on the refusal rows and gives `legacy empty` its own assertion.
- F2 `c1-f2-sd4-comment`, auto-fix, confidence 5: the doc comment at `internal/adopt/setup_test.go:76` names only SD4, but the test now also holds the exact seeded bytes. The repair states both facts.

### Spec

The axis reports 1 finding, and the worst is F1.

- F1 `c1-f1-nil-check`, auto-fix, confidence 4: the same finding as Standards F1.

### Coverage

The axis reports 2 findings, and the worst is F1.

- F1 `c1-f1-nil-check`, auto-fix, confidence 4: the same finding as Standards F1.
- F3, no-op, confidence 3: the TC3 why-clause overstates what the kept test asserts. Probe tc3 bit on `TestCompatibilityUndo`, so no repair target remains.
