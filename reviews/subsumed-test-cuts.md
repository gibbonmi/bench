# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/subsumed-test-cuts/spec.md",
  "plan_digest": "sha256:4576078b7e43a8211933e7577e092a7faf6f5a6880d23a56bdca340592b59c5f",
  "implementation_session": "",
  "chunks": [
    {
      "id": "FT411-C1",
      "base": "931b7ea78a972b7f45af2e7bfcfd7ddfc40c3a77",
      "tip": "3ea94e85180f56a97c537aadde2e074377bba1e1",
      "plan_digest": "sha256:4576078b7e43a8211933e7577e092a7faf6f5a6880d23a56bdca340592b59c5f",
      "source_digest": "889f110fe1cf439429acfc0fdbabe38d22884033",
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
