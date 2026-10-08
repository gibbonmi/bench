# Review outcomes

## Repair state

Chunk TF-C1 used repair cycle 1 of 2. The cycle repairs S1, S2, S3, and C1 through a fresh repair session for ticket 01. C1 is the one hardening cycle of TF-C1. No reviewer extension exists.

## Standards

TF-C1 round 1: 4 findings, worst S1.

- S1 (auto-fix, confidence 7): `internal/testreport/root_conformance_test.go:59-67` adds a third copy of the `selectRunBinary` override harness beside `installOwnedSelection` and `check_test.go:121-129`.
- S2 (auto-fix, confidence 6): `root_conformance_test.go:22` and `:101` and `check_test.go:226` spell the kit path fact three times, and `liveKitRoot` owns it.
- S3 (auto-fix, confidence 7): `check_test.go:201` keeps an independent help-list expectation that omits `root-conformance`, so an omission from help stays green.
- S4 (no-op, confidence 4): `testBenchSource` is a one-line forwarder to `gate.KitRoot`. The ticket requires that forwarder.

## Spec

TF-C1 round 1: 0 findings. TF1 to TF5 and the `testBenchSource` line hold.

Advice: the TF2 test runs the `--full` form of the command. The bare form keeps the same exit and diagnostic.

## Coverage

TF-C1 round 1: 1 finding, worst C1.

- C1 (auto-fix, confidence 8): `internal/gate/phases.go:151` has no test for the kit clause of `RootConformanceRuns`. A probe that drops the clause stays green.

Advice: the empty-kit clause has no caller that passes an empty kit.

```bench-review-record
{
  "version": 2,
  "spec": "specs/ticket-checkpoint-floor/spec.md",
  "plan_digest": "sha256:d9204cda44ff80f9e8917ef012b9baac98fccc7252cfedaa57d44cd6e32d6080",
  "implementation_session": "",
  "chunks": [
    {
      "id": "TF-C1",
      "base": "60da9ca95d42fda9f625fda4cb2ea96280bff4c2",
      "tip": "c0629ddc8ab1537d213840594febbb77605c013b",
      "plan_digest": "sha256:d9204cda44ff80f9e8917ef012b9baac98fccc7252cfedaa57d44cd6e32d6080",
      "source_digest": "844ebbe4af8f14822c400f2796c02feccf1e6874",
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
  }
}
```
