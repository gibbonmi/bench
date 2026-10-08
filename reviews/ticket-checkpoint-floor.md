# Review outcomes

## Repair state

Chunk TF-C1 used repair cycle 1 of 2 for S1, S2, S3, and C1. C1 is the one hardening cycle of TF-C1.

Round 2 confirmed those folds and found S5. Chunk TF-C1 uses repair cycle 2 of 2 for S5 through a fresh repair session for ticket 01. No reviewer extension exists. Under `--auto-approve`, the reviewer approved a fence expansion in advance. Round 3 confirmed the fold of S5 with 0 findings on every axis, so TF-C1 closes with both repair cycles used and no open finding.

Chunk TF-C2 uses repair cycle 1 of 2 for S1, S2, C1, C2, C3, C4, and C6 through a fresh repair session for ticket 02. C1, C2, and C6 are the one hardening cycle of TF-C2. S1 widens the ticket 02 fence by `internal/gate/phases.go` under the `--auto-approve` fence approval.

## Standards

TF-C1 round 1: 4 findings, worst S1.

- S1 (auto-fix, confidence 7): `internal/testreport/root_conformance_test.go:59-67` adds a third copy of the `selectRunBinary` override harness beside `installOwnedSelection` and `check_test.go:121-129`.
- S2 (auto-fix, confidence 6): `root_conformance_test.go:22` and `:101` and `check_test.go:226` spell the kit path fact three times, and `liveKitRoot` owns it.
- S3 (auto-fix, confidence 7): `check_test.go:201` keeps an independent help-list expectation that omits `root-conformance`, so an omission from help stays green.
- S4 (no-op, confidence 4): `testBenchSource` is a one-line forwarder to `gate.KitRoot`. The ticket requires that forwarder.

TF-C1 round 2: 1 finding, worst S5. The folds of S1, S2, and S3 hold.

- S5 (auto-fix, confidence 5): `internal/testreport/root_conformance_test.go:47-54` repeats the factory body of `installCannedSelection` in `outcome_test.go:79-88`. Only an observer flag differs.
- S6 (no-op, confidence 4): the help-list expectation in `check_test.go` has a recorded red. The repair session dropped `root-conformance` from `namedChecks` through `bench probe`, and the test failed at `check_test.go:194`.

TF-C1 round 3: 0 findings. The fold of S5 holds.

TF-C2 round 1: 3 findings, worst S1.

- S1 (auto-fix, confidence 6): `internal/testreport/floor.go:32` repeats the gate's rule that a root holds a regular `go.mod` file (`internal/gate/phases.go:115`).
- S2 (auto-fix, confidence 5): `floor.go:59-69` and `selection.go:231-241` each spell the rule that maps a path to its direct package.
- S3 (no-op, confidence 3): the floor does not clean an entry path. A `Writes:` entry is already a clean repository path.

## Spec

TF-C1 round 1: 0 findings. TF1 to TF5 and the `testBenchSource` line hold.

Advice: the TF2 test runs the `--full` form of the command. The bare form keeps the same exit and diagnostic.

TF-C1 round 2: 0 findings. TF1 to TF5 hold.

TF-C1 round 3: 0 findings. TF1 to TF5 hold.

TF-C2 round 1: 0 blocking findings. Every TF-C2 row holds.

Advice: P1, `go.work` and `go.work.sum` also add `./...` because the floor reuses the `--changed` metadata rule. The final report flags this widening for reviewer veto. P2, a `go.mod` entry leaves the embed flag false, and only the whole-suite command satisfies its element. P3, the author return holds the probe reds of TF21 and TF22.

## Coverage

TF-C1 round 1: 1 finding, worst C1.

- C1 (auto-fix, confidence 8): `internal/gate/phases.go:151` has no test for the kit clause of `RootConformanceRuns`. A probe that drops the clause stays green.

Advice: the empty-kit clause has no caller that passes an empty kit.

TF-C1 round 2: 0 findings. The fold of C1 holds.

TF-C1 round 3: 0 findings. The callers of `installCannedSelection` keep their behavior.

TF-C2 round 1: 5 findings, worst C1.

- C1 (auto-fix, confidence 8): `floor.go:40-43` has no test for the package-index error. A probe that drops the error stays green.
- C2 (auto-fix, confidence 7): `floor.go:36-39` has no floor test for the loader error.
- C3 (auto-fix, confidence 9): TF12 tests `go.mod` only, and story 11 also names `go.sum`.
- C4 (auto-fix, confidence 7): no test sets the embed flag through a directory entry, which story 15 requires.
- C6 (auto-fix, confidence 8): the refusal of a package outside the repository in `selection.go:54` has no test.

Advice: C5, a `go.mod` that is a directory or a link. The S1 repair makes the floor read the gate's rule, so the gate decides it.

```bench-review-record
{
  "version": 2,
  "spec": "specs/ticket-checkpoint-floor/spec.md",
  "plan_digest": "sha256:802ce2a779ccf3a4e748614f8b6837557f1f56e1ebe0232880dedbf0dba23a21",
  "implementation_session": "",
  "chunks": [
    {
      "id": "TF-C1",
      "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
      "tip": "776cac6e43bb1a8686bfe1399181c511bc044a74",
      "plan_digest": "sha256:63741a116956d2ebb6084503ed645e4a874351d06671cd4699639b095d8dc998",
      "source_digest": "e32356e41a7f6c51a086dfbd734eee8d0acae0ba",
      "acceptance_rows": [
        "TF1",
        "TF2",
        "TF3",
        "TF4",
        "TF5"
      ],
      "verification": [
        {
          "id": "t1-pkg-internal-testreport",
          "performer": "claude:ft349_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "844ebbe4af8f14822c400f2796c02feccf1e6874",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1/pkg-internal-testreport",
            "digest": "sha256:4a3f4c827a35e57f6b2ecfe9c26a54cbb966e330435920cdea73795f8b036ed5",
            "excerpt": "command: bench test --package ./internal/testreport\nexit=0 wall=81s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,74298,205\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "pkg-internal-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t1-pkg-internal-gate",
          "performer": "claude:ft349_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "844ebbe4af8f14822c400f2796c02feccf1e6874",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1/pkg-internal-gate",
            "digest": "sha256:587d1076c52217fa5f96a5c15c9d993ae32379554fd7354ab1ce8607faf30c86",
            "excerpt": "command: bench test --package ./internal/gate\nexit=0 wall=35s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,30422,420\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "pkg-internal-gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t1-pkg-cmd-bench",
          "performer": "claude:ft349_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "844ebbe4af8f14822c400f2796c02feccf1e6874",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1/pkg-cmd-bench",
            "digest": "sha256:f9d884cf7d04b180c4eaa7c3c347c1d4a475d0ea4cc451ce3c8c061d6fecd73d",
            "excerpt": "command: bench test --package ./cmd/bench\nexit=0 wall=25s\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20402,333\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1-root-conformance",
          "performer": "claude:ft349_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "844ebbe4af8f14822c400f2796c02feccf1e6874",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1/root-conformance",
            "digest": "sha256:c48974ab50513be561cd3afeea481679fbd0df16f25647e60d5df32167063a4f",
            "excerpt": "command: bench test --check root-conformance\nexit=0 wall=21s\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16447,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-root-scope-probe",
          "performer": "claude:ft349_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "844ebbe4af8f14822c400f2796c02feccf1e6874",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1/root-scope-probe",
            "digest": "sha256:aad7ba9d35076aac6e8459476d31bea598b0a64b185bb1ea61115db6677bc364",
            "excerpt": "command: bench probe internal/testreport/environment.go --swap (empty scope return) --with (return check) --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'\nexit=0 wall=5s\ntree[1]{target,head,dirty}:\n  ft349-build,a177c00d1c2e49dbc282442f201ee222921d6a5b,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/environment.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,^TestRootConformanceCheckRunsUnscoped$,passed,1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,65,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestRootConformanceCheckRunsUnscoped,\"root_conformance_test.go:95: root-conformance child lacks the line \\\"BENCH_CONFORMANCE_SCOPE=\\\":\",84\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-scope-probe",
          "command": "bench test --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the empty conformance scope for the check name in the root-conformance environment. TestRootConformanceCheckRunsUnscoped turns red.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft349_t1/root-scope-probe",
              "digest": "sha256:aad7ba9d35076aac6e8459476d31bea598b0a64b185bb1ea61115db6677bc364",
              "excerpt": "command: bench probe internal/testreport/environment.go --swap (empty scope return) --with (return check) --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'\nexit=0 wall=5s\ntree[1]{target,head,dirty}:\n  ft349-build,a177c00d1c2e49dbc282442f201ee222921d6a5b,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/environment.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,^TestRootConformanceCheckRunsUnscoped$,passed,1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,65,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestRootConformanceCheckRunsUnscoped,\"root_conformance_test.go:95: root-conformance child lacks the line \\\"BENCH_CONFORMANCE_SCOPE=\\\":\",84\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r1-pkg-internal-testreport",
          "performer": "claude:ft349_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca91c85ea8501c3bc21f49e8bdbeb7e68c9bcc6b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair1/pkg-internal-testreport",
            "digest": "sha256:4f69cb1cfd0605e95e8c916fed42b067155f02da90d0782ee1e8b33eda26df7a",
            "excerpt": "cmd: bench test --package ./internal/testreport\nexit: 0\nhead: a044f1155a78e044297f8993b201e0ca17f4e9ea dirty=false\n  github.com/gibbonmi/bench/internal/testreport,pass,58219,207\nfailures[0]\n"
          },
          "requirement": "pkg-internal-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t1r1-pkg-internal-gate",
          "performer": "claude:ft349_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca91c85ea8501c3bc21f49e8bdbeb7e68c9bcc6b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair1/pkg-internal-gate",
            "digest": "sha256:907bc23886639003c89b92c007341aeec56b0de1252f55c3a11840b89e6b7f47",
            "excerpt": "cmd: bench test --package ./internal/gate\nexit: 0\nhead: a044f1155a78e044297f8993b201e0ca17f4e9ea dirty=false\n  github.com/gibbonmi/bench/internal/gate,pass,18522,420\nfailures[0]\n"
          },
          "requirement": "pkg-internal-gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t1r1-pkg-cmd-bench",
          "performer": "claude:ft349_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca91c85ea8501c3bc21f49e8bdbeb7e68c9bcc6b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair1/pkg-cmd-bench",
            "digest": "sha256:0674b69871ea1e0e8272133f5314df775ca8cdf0fdf4f5557c55d13c61d6f30f",
            "excerpt": "cmd: bench test --package ./cmd/bench\nexit: 0\nhead: a044f1155a78e044297f8993b201e0ca17f4e9ea dirty=false\n  github.com/gibbonmi/bench/cmd/bench,pass,13359,333\nfailures[0]\n"
          },
          "requirement": "pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1r1-root-conformance",
          "performer": "claude:ft349_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca91c85ea8501c3bc21f49e8bdbeb7e68c9bcc6b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair1/root-conformance",
            "digest": "sha256:93d5a7156bb7c077eefcbf5005038a12f363c7f671cbfd57b14439f5ee482f68",
            "excerpt": "cmd: bench test --check root-conformance\nexit: 0\nhead: a044f1155a78e044297f8993b201e0ca17f4e9ea dirty=false\n  root-conformance,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,15699,1\nfailures[0]\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1r1-root-scope-probe",
          "performer": "claude:ft349_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ca91c85ea8501c3bc21f49e8bdbeb7e68c9bcc6b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair1/root-scope-probe",
            "digest": "sha256:8ef4578c090183a067e61298c47f585fd204222b5c540f0a1de52fb97addfdcd",
            "excerpt": "green: bench test --package ./internal/testreport --run ^TestRootConformanceCheckRunsUnscoped$ -> exit 0, pass,51,1\nprobe: bench probe internal/testreport/environment.go --swap return \"\" --with return check (root-conformance branch of conformanceScope)\n  bit,internal/testreport/environment.go,swap,failed,1,yes\n  TestRootConformanceCheckRunsUnscoped: root_conformance_test.go:85: child lacks the line \"BENCH_CONFORMANCE_SCOPE=\"\nhead: a044f1155a78e044297f8993b201e0ca17f4e9ea dirty=false\nprobe exit: 0 (rerun at a044f115)\n"
          },
          "requirement": "root-scope-probe",
          "command": "bench test --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the empty conformance scope for the check name in the root-conformance environment. TestRootConformanceCheckRunsUnscoped turns red.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft349_t1_repair1/root-scope-probe",
              "digest": "sha256:8ef4578c090183a067e61298c47f585fd204222b5c540f0a1de52fb97addfdcd",
              "excerpt": "green: bench test --package ./internal/testreport --run ^TestRootConformanceCheckRunsUnscoped$ -> exit 0, pass,51,1\nprobe: bench probe internal/testreport/environment.go --swap return \"\" --with return check (root-conformance branch of conformanceScope)\n  bit,internal/testreport/environment.go,swap,failed,1,yes\n  TestRootConformanceCheckRunsUnscoped: root_conformance_test.go:85: child lacks the line \"BENCH_CONFORMANCE_SCOPE=\"\nhead: a044f1155a78e044297f8993b201e0ca17f4e9ea dirty=false\nprobe exit: 0 (rerun at a044f115)\n"
            }
          }
        },
        {
          "id": "t1r2-pkg-internal-testreport",
          "performer": "claude:ft349_t1_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e32356e41a7f6c51a086dfbd734eee8d0acae0ba",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair2/pkg-internal-testreport",
            "digest": "sha256:11e5672ec966c3c8144143c90f2baa1ff609d6202737b69a41c57308f760fef6",
            "excerpt": "command: bench worktree exec ft349-build -- bench test --package ./internal/testreport\nhead: 776cac6e43bb1a8686bfe1399181c511bc044a74 (clean)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,63952,207\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit=0 wall=71.3s\n"
          },
          "requirement": "pkg-internal-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t1r2-pkg-internal-gate",
          "performer": "claude:ft349_t1_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e32356e41a7f6c51a086dfbd734eee8d0acae0ba",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair2/pkg-internal-gate",
            "digest": "sha256:ade346f379d4d01deaa8ed1bd047bb988b1b05539af2df8940936050747f371b",
            "excerpt": "command: bench worktree exec ft349-build -- bench test --package ./internal/gate\nhead: 776cac6e43bb1a8686bfe1399181c511bc044a74 (clean)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,20413,420\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit=0 wall=24.6s\n"
          },
          "requirement": "pkg-internal-gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t1r2-pkg-cmd-bench",
          "performer": "claude:ft349_t1_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e32356e41a7f6c51a086dfbd734eee8d0acae0ba",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair2/pkg-cmd-bench",
            "digest": "sha256:bb7021a9ac712288b7ab0f202370a830f336bba02d7fdf272b2add3b2ecb7811",
            "excerpt": "command: bench worktree exec ft349-build -- bench test --package ./cmd/bench\nhead: 776cac6e43bb1a8686bfe1399181c511bc044a74 (clean)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14653,333\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit=0 wall=17.7s\n"
          },
          "requirement": "pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1r2-root-conformance",
          "performer": "claude:ft349_t1_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e32356e41a7f6c51a086dfbd734eee8d0acae0ba",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair2/root-conformance",
            "digest": "sha256:2fe7be65bf512a0c490f9bd248226ff37272be68df4c115635a98deff3aca773",
            "excerpt": "command: bench worktree exec ft349-build -- bench test --check root-conformance\nhead: 776cac6e43bb1a8686bfe1399181c511bc044a74 (clean)\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16429,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit=0 wall=18.5s\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1r2-root-scope-probe",
          "performer": "claude:ft349_t1_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e32356e41a7f6c51a086dfbd734eee8d0acae0ba",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair2/root-scope-probe",
            "digest": "sha256:11943d512cd06aff513b868b4215c4db5141187c8b463b85be69b0865a132059",
            "excerpt": "command: bench worktree exec ft349-build -- bench probe internal/testreport/environment.go --swap '\tif check == rootConformanceCheckName {\\n\t\treturn \"\"' --with '\tif check == rootConformanceCheckName {\\n\t\treturn check' --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'\nhead: 776cac6e43bb1a8686bfe1399181c511bc044a74 (clean)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/environment.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestRootConformanceCheckRunsUnscoped,\"root_conformance_test.go:77: root-conformance child lacks the line \\\"BENCH_CONFORMANCE_SCOPE=\\\":\",84\nprobe exit=0 wall=3.2s\n"
          },
          "requirement": "root-scope-probe",
          "command": "bench test --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the empty conformance scope for the check name in the root-conformance environment. TestRootConformanceCheckRunsUnscoped turns red.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft349_t1_repair2/root-scope-probe",
              "digest": "sha256:11943d512cd06aff513b868b4215c4db5141187c8b463b85be69b0865a132059",
              "excerpt": "command: bench worktree exec ft349-build -- bench probe internal/testreport/environment.go --swap '\tif check == rootConformanceCheckName {\\n\t\treturn \"\"' --with '\tif check == rootConformanceCheckName {\\n\t\treturn check' --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'\nhead: 776cac6e43bb1a8686bfe1399181c511bc044a74 (clean)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/environment.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestRootConformanceCheckRunsUnscoped,\"root_conformance_test.go:77: root-conformance child lacks the line \\\"BENCH_CONFORMANCE_SCOPE=\\\":\",84\nprobe exit=0 wall=3.2s\n"
            }
          }
        },
        {
          "id": "t1r2-root-scope-probe-c1",
          "performer": "claude:ft349_t1_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e32356e41a7f6c51a086dfbd734eee8d0acae0ba",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t1_repair2/root-scope-probe",
            "digest": "sha256:11943d512cd06aff513b868b4215c4db5141187c8b463b85be69b0865a132059",
            "excerpt": "command: bench worktree exec ft349-build -- bench probe internal/testreport/environment.go --swap '\tif check == rootConformanceCheckName {\\n\t\treturn \"\"' --with '\tif check == rootConformanceCheckName {\\n\t\treturn check' --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'\nhead: 776cac6e43bb1a8686bfe1399181c511bc044a74 (clean)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/environment.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestRootConformanceCheckRunsUnscoped,\"root_conformance_test.go:77: root-conformance child lacks the line \\\"BENCH_CONFORMANCE_SCOPE=\\\":\",84\nprobe exit=0 wall=3.2s\n"
          },
          "requirement": "root-scope-probe",
          "command": "bench test --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the empty conformance scope for the check name in the root-conformance environment. TestRootConformanceCheckRunsUnscoped turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft349_t1_repair2/root-scope-probe",
              "digest": "sha256:11943d512cd06aff513b868b4215c4db5141187c8b463b85be69b0865a132059",
              "excerpt": "command: bench worktree exec ft349-build -- bench probe internal/testreport/environment.go --swap '\tif check == rootConformanceCheckName {\\n\t\treturn \"\"' --with '\tif check == rootConformanceCheckName {\\n\t\treturn check' --package ./internal/testreport --run '^TestRootConformanceCheckRunsUnscoped$'\nhead: 776cac6e43bb1a8686bfe1399181c511bc044a74 (clean)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/environment.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestRootConformanceCheckRunsUnscoped,\"root_conformance_test.go:77: root-conformance child lacks the line \\\"BENCH_CONFORMANCE_SCOPE=\\\":\",84\nprobe exit=0 wall=3.2s\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c1-standards-r1",
          "performer": "claude:ft349_c1_standards",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "844ebbe4af8f14822c400f2796c02feccf1e6874",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft349_c1_standards",
            "digest": "sha256:4172eb6ff1d9e87d4f3f71875df1440a27317a18d68e8cccf2fbc756924eb96b",
            "excerpt": "Standards axis TF-C1 (fable/high), diff 60da9ca9..c0629ddc. 4 advisory, 0 blocking (axis-stated); worst S1.\nS1 root_conformance_test.go:59-67 selectRunBinary override harness exists three times (here, installOwnedSelection :40-43, check_test.go:121-129); outcome_test.go:79-88 installCannedSelection owns the stub. auto-fix. conf 7.\nS2 root_conformance_test.go:22,:101 and check_test.go:226 spell filepath.Join(\"..\",\"..\") kit fact three times; liveKitRoot owns it. auto-fix. conf 6.\nS3 check_test.go:201 TestNamedCheckHelpListsEverySupportedCheck expectation omits root-conformance (and prose); omission from help stays green. auto-fix. conf 7.\nS4 command.go:331-336 testBenchSource is a one-line forwarder to gate.KitRoot. no-op. conf 4.\n"
          },
          "axis": "Standards",
          "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
          "tip": "c0629ddc8ab1537d213840594febbb77605c013b",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4"
          ],
          "supersedes": []
        },
        {
          "id": "c1-spec-r1",
          "performer": "claude:ft349_c1_spec",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "844ebbe4af8f14822c400f2796c02feccf1e6874",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_c1_spec",
            "digest": "sha256:d86f770d6b82fb2f2be561cab4a750a657ad43f2d74fad425f08b2b50f6828f3",
            "excerpt": "Spec axis TF-C1 (fable/high), diff 60da9ca9..c0629ddc. 1 advisory, 0 blocking; worst P1.\nP1 root_conformance_test.go:112 TF2 test runs --full --check root-conformance while spec names the bare command; behavior holds in bare form. no-op. conf 4.\nTF1-TF5 and the testBenchSource line satisfied; seam cells TF1-TF5 in citation form naming existing tests; bench coverage mapped, 71 rows.\n"
          },
          "axis": "Spec",
          "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
          "tip": "c0629ddc8ab1537d213840594febbb77605c013b",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1-coverage-r1",
          "performer": "claude:ft349_c1_coverage",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "844ebbe4af8f14822c400f2796c02feccf1e6874",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft349_c1_coverage",
            "digest": "sha256:646601a4a84087400a6791f57b15753be8877b649d2c0ac4adec3079789465ba",
            "excerpt": "Coverage axis TF-C1 (fable/high), diff 60da9ca9..c0629ddc. 2 advisory, 0 blocking; worst C1.\nC1 internal/gate/phases.go:151 kit clause of RootConformanceRuns untested: probe swap 'sameDirectory(root, kitRootAt(root, kit))' -> 'sameDirectory(root, root)' silent on ./internal/gate (420) and testreport TF1/TF3/TF5 (3), restored=yes. auto-fix: assert RootConformanceRuns(liveKitRoot(t), t.TempDir()) == false. conf 8.\nC2 phases.go:149-151 empty-kit clause unreachable from callers; probe silent, restored=yes. no-op. conf 7.\n"
          },
          "axis": "Coverage",
          "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
          "tip": "c0629ddc8ab1537d213840594febbb77605c013b",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "c1-standards-r2",
          "performer": "claude:ft349_c1_standards_r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "ca91c85ea8501c3bc21f49e8bdbeb7e68c9bcc6b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft349_c1_standards_r2",
            "digest": "sha256:02e021096666f73b419975df12ba0347142008113c70daf2f2a75300d12126e5",
            "excerpt": "Standards axis TF-C1 round 2 (sonnet/high), repair delta 2f4e9c6a..a044f115. Folds S1, S2, S3 held. 2 new low findings; worst S5.\nS5 root_conformance_test.go:47-54 runbinary.Factory literal repeats the body of installCannedSelection (outcome_test.go:79-88); only a selected=true flag differs. auto-fix. conf 5.\nS6 check_test.go:192 help expectation restates namedChecks(); no recorded red seen. no-op unless enforced. conf 4.\n"
          },
          "axis": "Standards",
          "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
          "tip": "a044f1155a78e044297f8993b201e0ca17f4e9ea",
          "finding_ids": [
            "S5"
          ],
          "supersedes": [
            "c1-standards-r1"
          ]
        },
        {
          "id": "c1-spec-r2",
          "performer": "claude:ft349_c1_spec_r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "ca91c85ea8501c3bc21f49e8bdbeb7e68c9bcc6b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_c1_spec_r2",
            "digest": "sha256:b957f862f3606226d797caaf0cec3bcf5063166c1079596a6586a87b455f2c53",
            "excerpt": "Spec axis TF-C1 round 2 (sonnet/high), repair delta 2f4e9c6a..a044f115. TF1-TF5 held; seam cells cite existing tests (root_conformance_test.go:71,90,109,131,138); bench coverage mapped, 71 rows; delta inside ticket 01 Writes. 0 findings.\n"
          },
          "axis": "Spec",
          "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
          "tip": "a044f1155a78e044297f8993b201e0ca17f4e9ea",
          "finding_ids": [],
          "supersedes": [
            "c1-spec-r1"
          ]
        },
        {
          "id": "c1-coverage-r2",
          "performer": "claude:ft349_c1_coverage_r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "ca91c85ea8501c3bc21f49e8bdbeb7e68c9bcc6b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_c1_coverage_r2",
            "digest": "sha256:5d1aab89c011cc0f4fca1c16d7141631745bb78b9c54c43c5943a62026b840b6",
            "excerpt": "Coverage axis TF-C1 round 2 (sonnet/high), repair delta 2f4e9c6a..a044f115. C1 held: kit-clause swap probe bit TestRootConformanceCheckRefusesAForeignRoot/foreign_kit at root_conformance_test.go:118, restored yes. S1 fold held: swaps in named_check.go (runGoTest dir) and environment.go (BENCH_KIT append) bit TestNamedCheckRunsFromKitAgainstLinkedConsumer at check_test.go:128 and :136, restored yes. 0 findings.\n"
          },
          "axis": "Coverage",
          "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
          "tip": "a044f1155a78e044297f8993b201e0ca17f4e9ea",
          "finding_ids": [],
          "supersedes": [
            "c1-coverage-r1"
          ]
        },
        {
          "id": "c1-standards-r3",
          "performer": "claude:ft349_c1_standards_r3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e32356e41a7f6c51a086dfbd734eee8d0acae0ba",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_c1_standards_r3",
            "digest": "sha256:13f98993a24ae36cef7d7d8049089ed83de0b76d86be8ac945b47aa658add3f9",
            "excerpt": "Standards axis TF-C1 round 3 (sonnet/high), repair delta e5c10e3b..776cac6e. S5 held: rootConformanceRun calls installCannedSelection (root_conformance_test.go:46); the one owner is outcome_test.go:81. No duplicated knowledge added; new doc comment obeys craft-comments. 0 findings.\n"
          },
          "axis": "Standards",
          "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
          "tip": "776cac6e43bb1a8686bfe1399181c511bc044a74",
          "finding_ids": [],
          "supersedes": [
            "c1-standards-r2"
          ]
        },
        {
          "id": "c1-spec-r3",
          "performer": "claude:ft349_c1_spec_r3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e32356e41a7f6c51a086dfbd734eee8d0acae0ba",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_c1_spec_r3",
            "digest": "sha256:5a7e098ba6c0fd803c1749c0d402062fba63f91dc47a73dd072838022cf594c5",
            "excerpt": "Spec axis TF-C1 round 3 (sonnet/high), repair delta e5c10e3b..776cac6e. TF1-TF5 held (root_conformance_test.go:63,82,101,123,130). Delta inside the widened ticket 01 Writes line; 7c308b1b adds outcome_test.go to the ticket and the spec fence consistently. bench coverage mapped, 71 rows. 0 findings.\n"
          },
          "axis": "Spec",
          "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
          "tip": "776cac6e43bb1a8686bfe1399181c511bc044a74",
          "finding_ids": [],
          "supersedes": [
            "c1-spec-r2"
          ]
        },
        {
          "id": "c1-coverage-r3",
          "performer": "claude:ft349_c1_coverage_r3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e32356e41a7f6c51a086dfbd734eee8d0acae0ba",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_c1_coverage_r3",
            "digest": "sha256:29834cd6acafab9a8ca3fb4c0c23d365e65ced2295a334d51c13d2d2841333db",
            "excerpt": "Coverage axis TF-C1 round 3 (sonnet/high), repair delta e5c10e3b..776cac6e. Fold held: six other installCannedSelection callers discard the *bool; factory unchanged apart from the observer. Probe swap causeChanged \"changed\" -> \"chnged\" in selection.go:40 bit TestChangedRowsCarrySelectedBy at selected_by_test.go:36, restored yes. bench test --package ./internal/testreport pass, 207 tests, about 66s. 0 findings.\n"
          },
          "axis": "Coverage",
          "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
          "tip": "776cac6e43bb1a8686bfe1399181c511bc044a74",
          "finding_ids": [],
          "supersedes": [
            "c1-coverage-r2"
          ]
        }
      ]
    },
    {
      "id": "TF-C2",
      "base": "07f1ee469e43e5d45a111cd5788df8c17f4d1cd2",
      "tip": "b9bf693e55636a0ef9331cf9226e75f14862fcda",
      "plan_digest": "sha256:802ce2a779ccf3a4e748614f8b6837557f1f56e1ebe0232880dedbf0dba23a21",
      "source_digest": "de79880a3a9ded9103cc4c9f40afcdc21c7d3c5e",
      "acceptance_rows": [
        "TF7",
        "TF8",
        "TF9",
        "TF10",
        "TF11",
        "TF12",
        "TF13",
        "TF14",
        "TF15",
        "TF19",
        "TF20",
        "TF21",
        "TF22",
        "TF71"
      ],
      "verification": [
        {
          "id": "t2-pkg-internal-testreport",
          "performer": "claude:ft349_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "de79880a3a9ded9103cc4c9f40afcdc21c7d3c5e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t2/pkg-internal-testreport",
            "digest": "sha256:515ef8e72024cb5eed230e49d10f0acf937e5ec2449b43ca6d6640b124722992",
            "excerpt": "command: bench test --package ./internal/testreport\nexit: 0\nwall: 59.98s\n  github.com/gibbonmi/bench/internal/testreport,pass,58249,219\nfailures[0]\nskips[0]\n"
          },
          "requirement": "pkg-internal-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t2-pkg-internal-tickets",
          "performer": "claude:ft349_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "de79880a3a9ded9103cc4c9f40afcdc21c7d3c5e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t2/pkg-internal-tickets",
            "digest": "sha256:0d1d5a25130216e675546efb942ede987c88a76617d06bf81f19e362b7f258da",
            "excerpt": "command: bench test --package ./internal/tickets\nexit: 0\nwall: 1.50s\n  github.com/gibbonmi/bench/internal/tickets,pass,3,46\nfailures[0]\nskips[0]\n"
          },
          "requirement": "pkg-internal-tickets",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "t2-pkg-cmd-bench",
          "performer": "claude:ft349_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "de79880a3a9ded9103cc4c9f40afcdc21c7d3c5e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t2/pkg-cmd-bench",
            "digest": "sha256:7b3cd28d39377c7cd75c8f7660f9d766b0fe1b1b76322e5d81a4aeb1fa4d9353",
            "excerpt": "command: bench test --package ./cmd/bench\nexit: 0\nwall: 18.91s\n  github.com/gibbonmi/bench/cmd/bench,pass,16208,333\nfailures[0]\nskips[0]\n"
          },
          "requirement": "pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t2-root-conformance",
          "performer": "claude:ft349_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "de79880a3a9ded9103cc4c9f40afcdc21c7d3c5e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t2/root-conformance",
            "digest": "sha256:8eca419ef980f53ef09f00ba0017c546ad106019d38986bc6ecf101e8e462d18",
            "excerpt": "command: bench test --check root-conformance\nexit: 0\nwall: 19.10s\n  root-conformance,conformance,1,0\n  github.com/gibbonmi/bench/internal/conformance,pass,16615,1\nfailures[0]\nskips[0]\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-floor-probe",
          "performer": "claude:ft349_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "de79880a3a9ded9103cc4c9f40afcdc21c7d3c5e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_t2/floor-probe",
            "digest": "sha256:7e55771aa0ef2274e5a9ff619525732cd396e221d485b08525c4ba792a6c8974",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft349-build,773775a31b8d61c08b192b2e6871b939d59408bd,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/floor.go,omit,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,^TestWritesFloor,passed,11\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,88,11\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestWritesFloorCoversADirectoryEntry,\"floor_test.go:75: floor = {Packages:[] Embed:false}, want packages [\\\"./nest\\\" \\\"./nest/inner\\\"] and embed false\",1\n  github.com/gibbonmi/bench/internal/testreport,TestWritesFloorKeepsTheSegmentBoundary,\"floor_test.go:82: floor = {Packages:[] Embed:false}, want packages [\\\"./internal/git\\\"] and embed false\",1\nskips[0]{package,test,reason}:\ncommand: bench probe internal/testreport/floor.go --omit <directory-entry loop> --package ./internal/testreport --run ^TestWritesFloor\nbench probe exit: 0\nwall: 3.56s\nmutated test run: status fail, 2 failed tests (exit 1 inferred; the probe prints no child exit code)\nbaseline: passed, 11 ran\nrestored: yes\n"
          },
          "requirement": "floor-probe",
          "command": "bench test --package ./internal/testreport --run '^TestWritesFloor'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the directory-entry rule from the floor function. TestWritesFloorCoversADirectoryEntry turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft349_t2/floor-probe",
              "digest": "sha256:7e55771aa0ef2274e5a9ff619525732cd396e221d485b08525c4ba792a6c8974",
              "excerpt": "tree[1]{target,head,dirty}:\n  ft349-build,773775a31b8d61c08b192b2e6871b939d59408bd,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/testreport/floor.go,omit,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,^TestWritesFloor,passed,11\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,88,11\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestWritesFloorCoversADirectoryEntry,\"floor_test.go:75: floor = {Packages:[] Embed:false}, want packages [\\\"./nest\\\" \\\"./nest/inner\\\"] and embed false\",1\n  github.com/gibbonmi/bench/internal/testreport,TestWritesFloorKeepsTheSegmentBoundary,\"floor_test.go:82: floor = {Packages:[] Embed:false}, want packages [\\\"./internal/git\\\"] and embed false\",1\nskips[0]{package,test,reason}:\ncommand: bench probe internal/testreport/floor.go --omit <directory-entry loop> --package ./internal/testreport --run ^TestWritesFloor\nbench probe exit: 0\nwall: 3.56s\nmutated test run: status fail, 2 failed tests (exit 1 inferred; the probe prints no child exit code)\nbaseline: passed, 11 ran\nrestored: yes\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c2-standards-r1",
          "performer": "claude:ft349_c2_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "de79880a3a9ded9103cc4c9f40afcdc21c7d3c5e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft349_c2_standards",
            "digest": "sha256:f1bda59133e9b0beb5e57555680794f79081ac2ea191d34ada0cb692820a762b",
            "excerpt": "Standards axis TF-C2 round 1 (opus/high), ticket diff 773775a3..b9bf693e. 3 findings, worst S1.\nS1 blocking floor.go:32 own go.mod regular-file check repeats the gate's unexported isRegularFile rule (internal/gate/phases.go:115, :343). auto-fix: export one gate predicate. conf 6.\nS2 advisory floor.go:59-69 and selection.go:231-241 each spell the path-to-direct-package rule; only the maps are shared. auto-fix: one packageIndex lookup. conf 5.\nS3 advisory floor does not Clean entries as selection.go:133 does. no-op. conf 3.\nisGoMetadata reuse for go.work: no finding.\n"
          },
          "axis": "Standards",
          "base": "07f1ee469e43e5d45a111cd5788df8c17f4d1cd2",
          "tip": "b9bf693e55636a0ef9331cf9226e75f14862fcda",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c2-spec-r1",
          "performer": "claude:ft349_c2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "de79880a3a9ded9103cc4c9f40afcdc21c7d3c5e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft349_c2_spec",
            "digest": "sha256:1f55009a14a23d503a9181f9313c7ad4d8070a4ea0700128b8b46d2699391385",
            "excerpt": "Spec axis TF-C2 round 1 (opus/high), ticket diff 773775a3..b9bf693e. All 14 rows and Request.Changed() implemented. 3 advisory findings, 0 blocking.\nP1 floor.go:55 isGoMetadata also maps go.work and go.work.sum to ./...; spec.md:134 names go.mod and go.sum. accept; flag for reviewer veto. no-op. conf 8.\nP2 floor.go:54-57 a go.mod entry leaves Embed false. no-op. conf 3.\nP3 no recorded probe red for TF21 and TF22 in the record. ask-user. conf 4.\n"
          },
          "axis": "Spec",
          "base": "07f1ee469e43e5d45a111cd5788df8c17f4d1cd2",
          "tip": "b9bf693e55636a0ef9331cf9226e75f14862fcda",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c2-coverage-r1",
          "performer": "claude:ft349_c2_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "de79880a3a9ded9103cc4c9f40afcdc21c7d3c5e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft349_c2_coverage",
            "digest": "sha256:39e2d35353756cd520fdb39b5e53e88aa12532256f2665ac4acb85ea3e09ce4f",
            "excerpt": "Coverage axis TF-C2 round 1 (opus/high), ticket diff 773775a3..b9bf693e. 6 probes, all silent, all restored=yes. 6 advisory findings; worst C1.\nC1 floor.go:40-43 newPackageIndex error swallowed -> silent. auto-fix: floor test with a package Dir outside root. conf 8.\nC2 floor.go:36-39 loader error swallowed -> silent. no-op if TF45 fails through WritesFloor, else auto-fix. conf 7.\nC3 floor.go:55 isGoMetadata -> entryPath == \"go.mod\" silent; TF12 tests go.mod only. auto-fix: add go.sum. conf 9.\nC4 floor.go:70-73 directory loop add -> selected[...]=true silent; embed flag for a directory entry untested. auto-fix. conf 7.\nC5 floor.go:32 IsRegular swap silent; go.mod directory or symlink undecided. ask-user. conf 6.\nC6 selection.go:54 outside-repository refusal untested (predates this change). auto-fix with C1. conf 8.\n"
          },
          "axis": "Coverage",
          "base": "07f1ee469e43e5d45a111cd5788df8c17f4d1cd2",
          "tip": "b9bf693e55636a0ef9331cf9226e75f14862fcda",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
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
      "from": "sha256:d9204cda44ff80f9e8917ef012b9baac98fccc7252cfedaa57d44cd6e32d6080",
      "to": "sha256:5ab185e24118847b76b5943bca905f1d204398378c7054a0d09b2307b1b6c981",
      "chunk_ids": {
        "TF-C1": [
          "TF-C1"
        ]
      }
    },
    {
      "from": "sha256:5ab185e24118847b76b5943bca905f1d204398378c7054a0d09b2307b1b6c981",
      "to": "sha256:63741a116956d2ebb6084503ed645e4a874351d06671cd4699639b095d8dc998",
      "chunk_ids": {
        "TF-C1": [
          "TF-C1"
        ]
      }
    },
    {
      "from": "sha256:63741a116956d2ebb6084503ed645e4a874351d06671cd4699639b095d8dc998",
      "to": "sha256:35975f8a331862c515848b125af2ab8e9e06ca5abe4ff7682fe97bba3f133356",
      "chunk_ids": {
        "TF-C1": [
          "TF-C1"
        ]
      }
    },
    {
      "from": "sha256:35975f8a331862c515848b125af2ab8e9e06ca5abe4ff7682fe97bba3f133356",
      "to": "sha256:802ce2a779ccf3a4e748614f8b6837557f1f56e1ebe0232880dedbf0dba23a21",
      "chunk_ids": {
        "TF-C1": [
          "TF-C1"
        ]
      }
    }
  ]
}
```
