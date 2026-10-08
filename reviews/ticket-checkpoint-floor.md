# Review outcomes

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
