# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/subsumed-test-cuts/spec.md",
  "plan_digest": "sha256:e1564e017bb221e293307e9ca521767f8831bc4731280ac8f2de1aeef747bfca",
  "implementation_session": "",
  "chunks": [
    {
      "id": "FT411-C1",
      "base": "931b7ea78a972b7f45af2e7bfcfd7ddfc40c3a77",
      "tip": "c299acc0b86168670e4e75b012e6b8bff6a1e78e",
      "plan_digest": "sha256:e1564e017bb221e293307e9ca521767f8831bc4731280ac8f2de1aeef747bfca",
      "source_digest": "13c164f90efc3783e7e9f05a4090bce87d0f4159",
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
    }
  ]
}
```
