# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/landing-gate-long-poles/spec.md",
  "plan_digest": "sha256:400cd4647d7f162a60c2c9c5d31dde7bdfd1f3789c3b6a3c878532c65ddfcdc3",
  "implementation_session": "",
  "chunks": [
    {
      "id": "LGP-C1",
      "base": "7f585da4489289b34df9cf3beeef2c59de8edb6d",
      "tip": "5cd7849154fcfed5f255cc99a6812d2f38d8cc7a",
      "plan_digest": "sha256:400cd4647d7f162a60c2c9c5d31dde7bdfd1f3789c3b6a3c878532c65ddfcdc3",
      "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
      "acceptance_rows": [
        "LGP1",
        "LGP2",
        "LGP3",
        "LGP4",
        "LGP5",
        "LGP37",
        "LGP34",
        "LGP6",
        "LGP7",
        "LGP8"
      ],
      "verification": [
        {
          "id": "t02-testreport",
          "performer": "claude:lgp-t02-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t02-author-20261021",
            "digest": "sha256:f015edfa1829b30a1818f304f1a6b4a00a2701d4a35b3692b21aef36bfe5693a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,75544,268\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t02-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t02-root-conformance",
          "performer": "claude:lgp-t02-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t02-author-20261021",
            "digest": "sha256:4309f069daf41f55a21f16178c1abb6d7b9009611058670e57199e76c7c202f8",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16954,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t02-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t02-selection-once",
          "performer": "claude:lgp-t02-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t02-author-20261021",
            "digest": "sha256:776ae2809d4365f98718502518bf2a7dc71bf53cd2299510826d29849e6fc691",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/root_conformance_test.go,swap,failed,1,TestOwnedSelectionBuildsOnce,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,TestOwnedSelectionBuildsOnce,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,513\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,4,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestOwnedSelectionBuildsOnce,\"selection_once_test.go:86: two kit copies ran the builder 2 times and selected map[/dev/shm/3H3UVT/t/TestOwnedSelectionBuildsOnce4096956678/001/bench-run-1400777246/bench:true /dev/shm/3H3UVT/t/TestOwnedSelectionBuildsOnce4096956678/001/ben… (303 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t02-selection-once",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "build the selection in each installOwnedSelection call",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t02-author-20261021",
              "digest": "sha256:776ae2809d4365f98718502518bf2a7dc71bf53cd2299510826d29849e6fc691",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/root_conformance_test.go,swap,failed,1,TestOwnedSelectionBuildsOnce,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,TestOwnedSelectionBuildsOnce,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,513\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,4,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestOwnedSelectionBuildsOnce,\"selection_once_test.go:86: two kit copies ran the builder 2 times and selected map[/dev/shm/3H3UVT/t/TestOwnedSelectionBuildsOnce4096956678/001/bench-run-1400777246/bench:true /dev/shm/3H3UVT/t/TestOwnedSelectionBuildsOnce4096956678/001/ben… (303 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01-pkg-cmd-bench",
          "performer": "claude:lgp-t01-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-author-20261021",
            "digest": "sha256:1fe7df6c0d551b0f1994b8d857b89df5722fed83e2e3b53504f3035e6648279a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19993,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t01-pkg-internal-probe-probetest",
          "performer": "claude:lgp-t01-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-author-20261021",
            "digest": "sha256:9a918b669075b6c5924526ca06ccff7e47c580b13650fd72d550a1ba14cd0209",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-pkg-internal-probe-probetest",
          "command": "bench test --package ./internal/probe/probetest",
          "exit_code": 0
        },
        {
          "id": "t01-pkg-internal-reviewrecord-recordcmd",
          "performer": "claude:lgp-t01-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-author-20261021",
            "digest": "sha256:ef5b941e89a226fd5f576269943735696f4c4b7b1028793449e59cd4f305ff8c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,23531,261\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-pkg-internal-reviewrecord-recordcmd",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "t01-probe",
          "performer": "claude:lgp-t01-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-author-20261021",
            "digest": "sha256:c57a81c9d550d506aad315e927de9e247875ee1279fb99e581cb9733cf8eaf29",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,pass,52821,318\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-probe",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        },
        {
          "id": "t01-root-conformance",
          "performer": "claude:lgp-t01-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-author-20261021",
            "digest": "sha256:3210efe233907d82217b7f63b87a7060eab7f8e0db26d48566e9ea4a7362b9b2",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16433,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t01-golden-use",
          "performer": "claude:lgp-t01-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-author-20261021",
            "digest": "sha256:3ada88949838057d04a4cefea35ffbc16a2524363dcf6afbc190e58ba49a6ddf",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestClampGoldenBuildsOnce,3,1,0,null,yes\n  2,bit,swap,failed,1,TestClampCopiesAreIsolated,3,1,0,null,yes\n  3,bit,insert,failed,1,TestClampCopyKeepsAFreshSeal,3,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/probe,^TestClamp(GoldenBuildsOnce|CopiesAreIsolated|CopyKeepsAFreshSeal)$,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1405\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,624,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampGoldenBuildsOnce,\"golden_test.go:118: the fixture publisher ran 2 times, want 1\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,207,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopiesAreIsolated,\"golden_test.go:125: subject bytes = \\\"package probefixture\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,243,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopyKeepsAFreshSeal,\"golden_test.go:132: the copy's run binary refuses its own root: selected Bench executable \\\"/dev/shm/AVHYDY/t/TestClampCopyKeepsAFreshSeal131029836/001/bin/bench\\\" does not match source \\\"/dev/shm/AVHYDY/t/TestClampCopyKeepsAFreshSeal131029836… (662 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-golden-use",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0,
          "probe": {
            "mutation": "build a fresh repository in clampFixture",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t01-author-20261021",
              "digest": "sha256:3ada88949838057d04a4cefea35ffbc16a2524363dcf6afbc190e58ba49a6ddf",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestClampGoldenBuildsOnce,3,1,0,null,yes\n  2,bit,swap,failed,1,TestClampCopiesAreIsolated,3,1,0,null,yes\n  3,bit,insert,failed,1,TestClampCopyKeepsAFreshSeal,3,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/probe,^TestClamp(GoldenBuildsOnce|CopiesAreIsolated|CopyKeepsAFreshSeal)$,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1405\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,624,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampGoldenBuildsOnce,\"golden_test.go:118: the fixture publisher ran 2 times, want 1\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,207,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopiesAreIsolated,\"golden_test.go:125: subject bytes = \\\"package probefixture\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,243,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopyKeepsAFreshSeal,\"golden_test.go:132: the copy's run binary refuses its own root: selected Bench executable \\\"/dev/shm/AVHYDY/t/TestClampCopyKeepsAFreshSeal131029836/001/bin/bench\\\" does not match source \\\"/dev/shm/AVHYDY/t/TestClampCopyKeepsAFreshSeal131029836… (662 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01-copy-isolation",
          "performer": "claude:lgp-t01-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-author-20261021",
            "digest": "sha256:3ada88949838057d04a4cefea35ffbc16a2524363dcf6afbc190e58ba49a6ddf",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestClampGoldenBuildsOnce,3,1,0,null,yes\n  2,bit,swap,failed,1,TestClampCopiesAreIsolated,3,1,0,null,yes\n  3,bit,insert,failed,1,TestClampCopyKeepsAFreshSeal,3,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/probe,^TestClamp(GoldenBuildsOnce|CopiesAreIsolated|CopyKeepsAFreshSeal)$,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1405\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,624,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampGoldenBuildsOnce,\"golden_test.go:118: the fixture publisher ran 2 times, want 1\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,207,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopiesAreIsolated,\"golden_test.go:125: subject bytes = \\\"package probefixture\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,243,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopyKeepsAFreshSeal,\"golden_test.go:132: the copy's run binary refuses its own root: selected Bench executable \\\"/dev/shm/AVHYDY/t/TestClampCopyKeepsAFreshSeal131029836/001/bin/bench\\\" does not match source \\\"/dev/shm/AVHYDY/t/TestClampCopyKeepsAFreshSeal131029836… (662 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-copy-isolation",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0,
          "probe": {
            "mutation": "return the golden root instead of a private copy",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t01-author-20261021",
              "digest": "sha256:3ada88949838057d04a4cefea35ffbc16a2524363dcf6afbc190e58ba49a6ddf",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestClampGoldenBuildsOnce,3,1,0,null,yes\n  2,bit,swap,failed,1,TestClampCopiesAreIsolated,3,1,0,null,yes\n  3,bit,insert,failed,1,TestClampCopyKeepsAFreshSeal,3,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/probe,^TestClamp(GoldenBuildsOnce|CopiesAreIsolated|CopyKeepsAFreshSeal)$,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1405\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,624,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampGoldenBuildsOnce,\"golden_test.go:118: the fixture publisher ran 2 times, want 1\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,207,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopiesAreIsolated,\"golden_test.go:125: subject bytes = \\\"package probefixture\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,243,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopyKeepsAFreshSeal,\"golden_test.go:132: the copy's run binary refuses its own root: selected Bench executable \\\"/dev/shm/AVHYDY/t/TestClampCopyKeepsAFreshSeal131029836/001/bin/bench\\\" does not match source \\\"/dev/shm/AVHYDY/t/TestClampCopyKeepsAFreshSeal131029836… (662 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
