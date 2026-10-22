# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/landing-gate-long-poles/spec.md",
  "plan_digest": "sha256:b5e476e1aadba9767c08baeae66130b308c5e2b5302df2a0170bf35907fad091",
  "implementation_session": "",
  "chunks": [
    {
      "id": "LGP-C1",
      "base": "7f585da4489289b34df9cf3beeef2c59de8edb6d",
      "tip": "be0aca0c526c177be15b360eeba46c7b57927ee8",
      "plan_digest": "sha256:e15cb1da923ae005902078362e7c2e00b733194f22922820f90fa00d060865c7",
      "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
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
        },
        {
          "id": "t02-r1-t02-testreport",
          "performer": "claude:lgp-t02-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t02-repair1-20261021",
            "digest": "sha256:7ec2f0055ccad2deeda90fdfb9c437e3b5b2c9b271a2443a93ca139a5427573b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,76757,268\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t02-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t02-r1-t02-root-conformance",
          "performer": "claude:lgp-t02-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t02-repair1-20261021",
            "digest": "sha256:9e2666505401191a80a3839ce28696d22e78aec4b12e38f2edf25e1e7956a23b",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16859,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t02-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t02-r1-t02-selection-once",
          "performer": "claude:lgp-t02-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t02-repair1-20261021",
            "digest": "sha256:2e5c96cfef5f398f38ffa62a0926840b382f08fbffa5288a9e66bf047a8d63c2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/root_conformance_test.go,swap,failed,1,TestOwnedSelectionBuildsOnce,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,TestOwnedSelectionBuildsOnce,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,513\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,4,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestOwnedSelectionBuildsOnce,\"selection_once_test.go:50: two kit copies ran the builder 2 times and selected map[/dev/shm/A6TSH4/t/TestOwnedSelectionBuildsOnce991126804/001/bench-run-2349916703/bench:true /dev/shm/A6TSH4/t/TestOwnedSelectionBuildsOnce991126804/001/bench… (301 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:bench-writer-high:lgp-t02-repair1-20261021",
              "digest": "sha256:2e5c96cfef5f398f38ffa62a0926840b382f08fbffa5288a9e66bf047a8d63c2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/root_conformance_test.go,swap,failed,1,TestOwnedSelectionBuildsOnce,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,TestOwnedSelectionBuildsOnce,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,513\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,4,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestOwnedSelectionBuildsOnce,\"selection_once_test.go:50: two kit copies ran the builder 2 times and selected map[/dev/shm/A6TSH4/t/TestOwnedSelectionBuildsOnce991126804/001/bench-run-2349916703/bench:true /dev/shm/A6TSH4/t/TestOwnedSelectionBuildsOnce991126804/001/bench… (301 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t02-r1-t02-runbinary",
          "performer": "claude:lgp-t02-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t02-repair1-20261021",
            "digest": "sha256:8892eae49d08bc12e78ac8f4ae3e5ef105bcd3736967a57daa172ab523e63c8c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,pass,26094,35\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t02-runbinary",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "t02-r1-t02-worktree",
          "performer": "claude:lgp-t02-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t02-repair1-20261021",
            "digest": "sha256:e87c5f730c8513f625a1b4d11d6953c581c7aa3ecfa36ba21c28170812186a5e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,89518,1491\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t02-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t02-r1-t02-select-inherit",
          "performer": "claude:lgp-t02-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t02-repair1-20261021",
            "digest": "sha256:6986e1c6980ff3ba6514be6e6e0df2631ce02612acfb5b88430082a8110dc8fd",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/runbinary/runbinary.go,swap,failed,5,TestSharedInheritedSelectionWinsWithZeroBuilds,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,^TestShared,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1929\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,4,6\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestSharedInheritedSelectionWinsWithZeroBuilds,\"shared_test.go:42: first caller selected \\\"/dev/shm/D3FY2T/t/TestSharedInheritedSelectionWinsWithZeroBuilds2728413877/002/bench-run-781537048/bench\\\", <nil>; want the inherited \\\"/dev/shm/D3FY2T/t/TestSharedInheritedSelectionWinsWithZeroBuilds… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/runbinary,TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild/empty,\"shared_test.go:65: the present empty value \\\"\\\" selected \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuildempty1064602961/001/bench-run-462498608/bench\\\", want a refusal\",1\n  github.com/gibbonmi/bench/internal/runbinary,TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild/missing,\"shared_test.go:65: the present missing value \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild2462681282/001/absent\\\" selected \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuildmissing392224… (292 bytes)\",1\n  github.com/gibbonmi/bench/internal/runbinary,TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild/no_source_root,\"shared_test.go:65: the present no source root value \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild2462681282/005/bench\\\" selected \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuildno_sour… (297 bytes)\",1\n  github.com/gibbonmi/bench/internal/runbinary,TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild/relative,\"shared_test.go:65: the present relative value \\\"dist/bench\\\" selected \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuildrelativ343471130/001/bench-run-3137205700/bench\\\", want a refusal\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t02-select-inherit",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0,
          "probe": {
            "mutation": "ignore a present inherited value in Factory.Select; the inherited-selection test must red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t02-repair1-20261021",
              "digest": "sha256:6986e1c6980ff3ba6514be6e6e0df2631ce02612acfb5b88430082a8110dc8fd",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/runbinary/runbinary.go,swap,failed,5,TestSharedInheritedSelectionWinsWithZeroBuilds,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,^TestShared,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1929\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,4,6\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestSharedInheritedSelectionWinsWithZeroBuilds,\"shared_test.go:42: first caller selected \\\"/dev/shm/D3FY2T/t/TestSharedInheritedSelectionWinsWithZeroBuilds2728413877/002/bench-run-781537048/bench\\\", <nil>; want the inherited \\\"/dev/shm/D3FY2T/t/TestSharedInheritedSelectionWinsWithZeroBuilds… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/runbinary,TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild/empty,\"shared_test.go:65: the present empty value \\\"\\\" selected \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuildempty1064602961/001/bench-run-462498608/bench\\\", want a refusal\",1\n  github.com/gibbonmi/bench/internal/runbinary,TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild/missing,\"shared_test.go:65: the present missing value \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild2462681282/001/absent\\\" selected \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuildmissing392224… (292 bytes)\",1\n  github.com/gibbonmi/bench/internal/runbinary,TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild/no_source_root,\"shared_test.go:65: the present no source root value \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild2462681282/005/bench\\\" selected \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuildno_sour… (297 bytes)\",1\n  github.com/gibbonmi/bench/internal/runbinary,TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild/relative,\"shared_test.go:65: the present relative value \\\"dist/bench\\\" selected \\\"/dev/shm/D3FY2T/t/TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuildrelativ343471130/001/bench-run-3137205700/bench\\\", want a refusal\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t02-r1-t02-pkg-cmd-bench",
          "performer": "claude:lgp-t02-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t02-repair1-20261021",
            "digest": "sha256:a768910fb2f2e886d1e348a41f1cf48c2dcdcc7571dec50b18ee996dae6a1608",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19686,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t02-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t01-r1-t01-pkg-cmd-bench",
          "performer": "claude:lgp-t01-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
            "digest": "sha256:34a87052d7fc6dd43a7dcaaa276237a7d8aed343ad6f876453f000029d51c639",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19949,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t01-r1-t01-pkg-internal-probe-probetest",
          "performer": "claude:lgp-t01-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
            "digest": "sha256:9a918b669075b6c5924526ca06ccff7e47c580b13650fd72d550a1ba14cd0209",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-pkg-internal-probe-probetest",
          "command": "bench test --package ./internal/probe/probetest",
          "exit_code": 0
        },
        {
          "id": "t01-r1-t01-pkg-internal-reviewrecord-recordcmd",
          "performer": "claude:lgp-t01-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
            "digest": "sha256:79606df07eca0042e0d32119f7e9be9719fcaae6726b310713ac48491208a747",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,23815,261\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-pkg-internal-reviewrecord-recordcmd",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "t01-r1-t01-probe",
          "performer": "claude:lgp-t01-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
            "digest": "sha256:192dcc445f192daed11794c9ebbfeb77004419c4352e4b44bd617afeaefde869",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,pass,54364,318\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-probe",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        },
        {
          "id": "t01-r1-t01-root-conformance",
          "performer": "claude:lgp-t01-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
            "digest": "sha256:d1061c0be3f85ef5dcd1be068dca23f7d225232876339019c1bc9aa14798123b",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16785,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t01-r1-t01-golden-use",
          "performer": "claude:lgp-t01-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
            "digest": "sha256:062a51f3143ce191923a762d43a44c2c142712dd3f5db5995394044c61d4a43f",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestClampGoldenBuildsOnce,318,1,0,null,yes\n  2,bit,swap,failed,196,TestClampCopiesAreIsolated,286,1,0,null,yes\n  3,bit,insert,failed,110,TestClampCopyKeepsAFreshSeal,318,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/probe,all,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,79195\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,76146,318\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampGoldenBuildsOnce,\"golden_test.go:119: the fixture publisher ran 55 times, want 1\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,32669,286\nfailures[196]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopiesAreIsolated,\"golden_test.go:126: subject bytes = \\\"package probefixture\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMakeOneMutant,\"combined_test.go:17: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^Tes… (545 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/ambiguous,\"combined_test.go:34: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,\\\\\\\\\\\"\\\\\\\\\\\"… (461 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/empty_pair,\"combined_test.go:39: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/start_bytes,\"combined_test.go:29: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap,refused,0,null,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegative$,pa… (589 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseOverlap,\"combined_test.go:47: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,\\\\\\\\\\\"\\\\\\\\\\\"… (444 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--delete,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--insert-before,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--omit-file,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--unwrap,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--with-file,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--with-file_without_--with,\"combined_test.go:73: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteRemovesTheSubjectForTheRun,\"delete_test.go:41: /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture… (354 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_and_restore,\"delete_test.go:103: /dev/shm/AVONEO/t/TestDeleteReportsARestoreFailurefailed_removal_and_restore4049586931/002/probe/repo-567104589/1792634926202973261-2144253/clamp.go = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Cla… (448 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_restored,\"delete_test.go:108: chmod /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_name_prefix,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_passing_mutated_run,\"verdict_test.go:60: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/another_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/prose_check,\"verdict_test.go:67: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/the_named_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/baseline_refusal,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/build_failure,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/prose_bite,\"verdict_test.go:101: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/silent_run,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/two_failed_tests,\"verdict_test.go:77: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/delete_with_with,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/delete_with_with-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/insert-before_without_text,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_combined,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_delete,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_omit-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/swap_without_its_with,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/with-file_without_swap,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/with_and_with-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/ambiguous_anchor,\"edit_test.go:50: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with,\"edit_test.go:38: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,insert,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with-file,\"edit_test.go:45: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,insert,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_control_byte_in_the_first_build_line,\"verdict_test.go:137: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_run-binary_build_script_failure_at_the_baseline,\"verdict_test.go:152: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/an_unused_import,\"verdict_test.go:119: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/inherited,\"set_test.go:110: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/owned,\"set_test.go:110: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_restore-failed,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_silent,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/two_bites,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryCap,\"set_input_test.go:290: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryMatchesTheSingleForm,\"set_test.go:31: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusalNamesTheEntry,\"set_input_test.go:310: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusesProseExpect,\"set_input_test.go:203: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusesSelectionTokens,\"set_input_test.go:181: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryTakesACombinedMutant,\"set_test.go:69: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetExpectIsPerEntry,\"set_test.go:59: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetGradesEachEntry,\"set_test.go:21: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissIsPerEntry,\"set_test.go:78: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetNamesABuildFailure,\"set_test.go:357: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetOfMissesStartsNoRun,\"set_test.go:87: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetOfProseBuildsNothing,\"set_test.go:265: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetPrintsEachReport,\"set_test.go:335: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/inherited,\"set_test.go:250: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/owned,\"set_test.go:250: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRestoresAfterEachMutant,\"set_test.go:155: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReusesTheBuildOutsideTheClosure,\"set_test.go:236: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRowsKeepTheirSchema,\"set_test.go:284: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:131: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetSharesOneBaseline,\"set_test.go:96: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/directory-sync,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/failed_restore,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/rename,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAPreserveFailure,\"set_test.go:222: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnARestoreFailure,\"set_test.go:172: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnASignalDuringARestore,\"baseline_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,substring-miss,0,null,0,null,1,null,untouched\\\\\\\\n  2,bit,swap,f… (962 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnUnreadableGate,\"lock_test.go:298: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (363 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsUnderAGateRun,\"lock_test.go:286: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (356 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/held_by_another_probe,\"lock_test.go:266: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/tried_during_mutant_2,\"lock_test.go:273: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesRepeatedPackages,\"set_test.go:366: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsExcludesTopLevelMutation,\"set_input_test.go:212: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsFileRefusesHostilePaths,\"set_input_test.go:229: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsFileRefusesMalformedJSON,\"set_input_test.go:259: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/out_of_range,\"edit_test.go:66: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation occurrence out of range — the old string matches 2 times, --occurrence is 3\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/second_match,\"edit_test.go:60: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_-1,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_0,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_99999999999999999999,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_x,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmissionThatEmptiesTheSubjectRefuses/combined,\"delete_test.go:138: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empties the subject — use --delete to remove clamp.go for the run\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmissionThatEmptiesTheSubjectRefuses/omit-file,\"delete_test.go:138: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empties the subject — use --delete to remove clamp.go for the run\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesAnUnreadableInputBeforeRunning/--omit-file,\"omit_file_test.go:77: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesAnUnreadableInputBeforeRunning--omit-file2737955549/005/mut… (270 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesAnUnreadableInputBeforeRunning/--with-file,\"omit_file_test.go:77: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesAnUnreadableInputBeforeRunning--with-file43956786/005/mutat… (268 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_absent,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_absent1411704031/005/… (269 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_directory,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_direct4008031769/005/… (274 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_empty,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_empty2723731046/005/m… (267 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_oversized,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_oversi322350743/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_special,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_specia1647278768/005/… (277 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_symlink,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_symlin2171921339/005/… (272 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_absent,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_absent4234421728/005/… (269 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_directory,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_direct1967178204/005/… (274 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_empty,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_empty18947531/005/mut… (265 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_oversized,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_oversi120057229/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_special,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_specia3710825161/005/… (277 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_symlink,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_symlin698694105/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesZeroAndMultipleMatchesSeparately/multiple_matches,\"omit_file_test.go:97: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesZeroAndMultipleMatchesSeparately/zero_matches,\"omit_file_test.go:90: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,null,untouched\\\\\\\\ns… (446 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileUsesExactBytesAndReportsOmit,\"omit_file_test.go:20: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,omit,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first … (247 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeBitesWhenTheFocusedTestFails,\"probe_test.go:195: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeBuildsTheMutantRunBinary,\"subject_check_test.go:109: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeCarriesTheFocusedRunTables,\"probe_test.go:247: exit = 1, want 0\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheRestoreFailureSuffix,\"result_test.go:68: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheTestsThatRan,\"probe_test.go:278: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"selection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegative$,passed,2\\\\\\\\n\\\" and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeFailedRestore,\"replacement_test.go:172: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenNoTestRuns,\"probe_test.go:236: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,no-test-run,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenTheMutationDoesNotCompile,\"probe_test.go:225: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,build-failed,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsSilentWhenTheFocusedTestPasses,\"probe_test.go:215: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  silent,clamp.go,swap,passed,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationPreparationFailure,\"replacement_test.go:205: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoveryFailure,\"replacement_test.go:100: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-close,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-open,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-sync,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/focused_run,\"replacement_test.go:143: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/mutation_failure,\"replacement_test.go:143: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCheckSelection,\"outcome_test.go:236: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the check selection row \\\"selection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,line-routing,^TestRootConformance$,passed,1\\\\\\\\n\\\" and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesThePackageSelection,\"outcome_test.go:247: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection… (334 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeOmitsTheMatchOnce,\"probe_test.go:207: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,omit,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbePackagesShareOneBuildPerRun,\"package_test.go:59: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first an… (245 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbePreservesTheSubjectBeforeTheRun,\"outcome_test.go:27: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_directory,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_special_file,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_symlink,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/absent,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/outside_the_repository,\"refusal_test.go:127: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnAmbiguousMutation/no_match,\"refusal_test.go:61: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nsel… (444 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnAmbiguousMutation/several_matches,\"refusal_test.go:61: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnEmptyLaterPackage,\"refusal_test.go:176: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnEmptyMutation,\"refusal_test.go:71: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empty — --with equals --swap\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUngradedProseSubject/matched,\"subject_check_test.go:65: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUngradedProseSubject/substring_miss,\"subject_check_test.go:65: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUnknownCheckBeforeAnyWrite,\"refusal_test.go:184: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesUnderALiveGateRun,\"refusal_test.go:90: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: gate execution in progress — wait for the gate run to finish before you mutate the tree\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenPreservationFails,\"refusal_test.go:209: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenTheMutationWriteFails,\"refusal_test.go:223: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:31: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection[… (346 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReplacement,\"replacement_test.go:49: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsAReadBackMismatch,\"outcome_test.go:69: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  restore-failed,clamp.go,swap,failed,1,TestClampNegative,no\\\\\\\\n… (255 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsARestoreFailure,\"outcome_test.go:51: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  restore-failed,clamp.go,swap,failed,1,TestClampNegative,no\\\\\\\\n… (255 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/bit,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/red_baseline,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/substring_miss,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResolvesTheSubjectFromTheWorkingDirectory,\"probe_test.go:294: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultEnvelopeRefusals,\"result_test.go:133: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultSchemaVersion,\"result_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeSubjectRefusalSkipsReplacement,\"replacement_test.go:185: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheProseCheck,\"subject_check_test.go:43: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/kit,\"subject_check_test.go:79: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/not_the_kit,\"subject_check_test.go:95: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUnwrapBitesAndRestores,\"unwrap_test.go:6: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_omit-file,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_unwrap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_delete,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/no_mutation,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/no_selection,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_unwrap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/run_with_a_check_other_than_system,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/two_selections,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/unwrap_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/with_only_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeWritesNoRecord,\"outcome_test.go:166: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantRequiresTheCompleteAnswer,\"result_test.go:208: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantSharesTheReaderRules,\"result_test.go:232: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadResultRefusesAMutantSet,\"result_test.go:199: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_line-routing,\"subject_check_test.go:136: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_system_--run_^TestClampNegative$,\"subject_check_test.go:136: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/unknown_check,\"subject_check_test.go:125: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockFollowsTheGateLock,\"lock_test.go:59: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegativ… (399 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockIgnoresAFileWithoutAHolder,\"lock_test.go:131: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockKeepsItsFilesPrivate,\"lock_test.go:250: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReadsOnlyContentionAsInUse,\"lock_test.go:184: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_on_every_attempt,\"lock_test.go:162: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_once,\"lock_test.go:148: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesALinkAtTheLockPath,\"lock_test.go:234: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesASecondProbe,\"lock_test.go:44: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegativ… (399 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesAnUnusableDirectory,\"lock_test.go:175: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReleasesOnInterrupt,\"lock_test.go:115: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_on_every_attempt,\"lock_test.go:214: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_once,\"lock_test.go:200: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapEqualToOldRefusesEmpty,\"edit_test.go:28: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empty — --with-file equals --swap\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapsMultilineBytes,\"edit_test.go:18: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestZeroMatchesStayASubstringMiss/combined,\"edit_test.go:118: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\n… (447 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestZeroMatchesStayASubstringMiss/occurrence,\"edit_test.go:118: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselec… (442 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,153196,318\nfailures[110]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopyKeepsAFreshSeal,\"golden_test.go:133: the copy's run binary refuses its own root: selected Bench executable \\\"/dev/shm/3G6OMR/t/TestClampCopyKeepsAFreshSeal4239075508/001/bin/bench\\\" does not match source \\\"/dev/shm/3G6OMR/t/TestClampCopyKeepsAFreshSeal42390755… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMakeOneMutant,\"combined_test.go:17: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap+omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^T… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/start_bytes,\"combined_test.go:29: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestCl… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteRemovesTheSubjectForTheRun,\"delete_test.go:40: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_and_restore,\"delete_test.go:95: exit = 1, want 2\",8\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_restored,\"delete_test.go:111: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/restore_proves_the_mode,\"delete_test.go:83: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_name_prefix,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_passing_mutated_run,\"verdict_test.go:63: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/another_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/the_named_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/build_failure,\"verdict_test.go:97: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/silent_run,\"verdict_test.go:97: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/two_failed_tests,\"verdict_test.go:80: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with,\"edit_test.go:38: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,insert,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with-file,\"edit_test.go:45: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,insert,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_control_byte_in_the_first_build_line,\"verdict_test.go:148: stdout = \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/an_unused_import,\"verdict_test.go:128: stdout = \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/inherited,\"set_test.go:116: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/owned,\"set_test.go:116: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_restore-failed,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_silent,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/two_bites,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryMatchesTheSingleForm,\"set_test.go:44: single form of entry 1 = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryTakesACombinedMutant,\"set_test.go:73: stdout = (\\\"mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap+swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\nselection[1]{… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetExpectIsPerEntry,\"set_test.go:64: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetGradesEachEntry,\"set_test.go:24: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissIsPerEntry,\"set_test.go:81: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,omit,substring-miss,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,bas… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/baseline,\"baseline_test.go:248: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/interrupt,\"baseline_test.go:239: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/restore_failure,\"baseline_test.go:256: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetNamesABuildFailure,\"set_test.go:361: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetPrintsEachReport,\"set_test.go:347: stdout = (\\\"mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,omit,… (1020 bytes)\",4\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/inherited,\"set_test.go:257: stdout = \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/owned,\"set_test.go:257: stdout = \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRestoresAfterEachMutant,\"set_test.go:161: exit = 1, want 0\",9\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReusesTheBuildOutsideTheClosure,\"set_test.go:240: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRowsKeepTheirSchema,\"set_test.go:300: selection row = \\\"  package,./,^TestClampNegative$,refused,0\\\", want \\\"  package,./,^TestClampNegative$,passed,0\\\"\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:147: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetSharesOneBaseline,\"set_test.go:100: exit = 1, want 0\",10\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/directory-sync,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/failed_restore,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/rename,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAPreserveFailure,\"set_test.go:228: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnARestoreFailure,\"set_test.go:177: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnASignalDuringARestore,\"baseline_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnInterruptedBuild,\"baseline_test.go:222: no process created /dev/shm/3G6OMR/t/TestMutantSetStopsOnAnInterruptedBuild3031501550/004/parked within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnUnreadableGate,\"lock_test.go:298: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnInterrupt,\"baseline_test.go:201: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsUnderAGateRun,\"lock_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/tried_during_mutant_2,\"lock_test.go:273: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesRepeatedPackages,\"set_test.go:374: exit = 1, want 0\",9\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/second_match,\"edit_test.go:60: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegati… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileUsesExactBytesAndReportsOmit,\"omit_file_test.go:20: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeBitesWhenTheFocusedTestFails,\"probe_test.go:195: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeBuildsTheMutantRunBinary,\"subject_check_test.go:113: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,system,^Tes… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeCarriesTheFocusedRunTables,\"probe_test.go:247: exit = 1, want 0\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheRestoreFailureSuffix,\"result_test.go:74: restore failure = (1, {Verdict:invalid Restored:false CommandExit:1 FocusedExit:<nil>}), want command 2, focused 1, unrestored\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheTestsThatRan,\"probe_test.go:278: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeFailedRestore,\"replacement_test.go:178: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestCla… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeForwardsFullToTheFocusedRun/baseline,\"baseline_test.go:160: --full did not reach the run:\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeForwardsFullToTheFocusedRun/mutated_run,\"baseline_test.go:160: --full did not reach the run:\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenNoTestRuns,\"probe_test.go:236: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenTheMutationDoesNotCompile,\"probe_test.go:225: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0\\\\\\\\… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsSilentWhenTheFocusedTestPasses,\"probe_test.go:215: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoveryFailure,\"replacement_test.go:110: exit = 1, want 2\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-close,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-open,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-sync,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/focused_run,\"replacement_test.go:151: exit = 1, want 2\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/mutation_failure,\"replacement_test.go:151: exit = 1, want 2\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCheckSelection,\"outcome_test.go:236: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,line-routing,^Tes… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCopyWhenTheRenderAndTheRestoreFail,\"outcome_test.go:106: stdout = (\\\"error: unrepresentable TOON cell — toon: unsupported control character U+0007 in string\\\\\\\\n\\\", 1), want \\\"error: unrepresentable TOON cell — \\\" first and 2\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesThePackageSelection,\"outcome_test.go:247: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeOmitsTheMatchOnce,\"probe_test.go:207: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbePackagesShareOneBuildPerRun,\"package_test.go:59: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./clamp ./other,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbePreservesTheSubjectBeforeTheRun,\"outcome_test.go:27: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesARedBaseline,\"baseline_test.go:103: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenPreservationFails,\"refusal_test.go:213: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenTheMutationWriteFails,\"refusal_test.go:230: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRendersASubjectWithASpace,\"outcome_test.go:142: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,my pkg/clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^Test… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:31: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./clamp ./other,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReplacement,\"replacement_test.go:71: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClamp… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineRefusal,\"baseline_test.go:137: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineWithoutAVerdict/baseline-build-failed,\"baseline_test.go:127: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineWithoutAVerdict/baseline-no-test-run,\"baseline_test.go:127: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsAReadBackMismatch,\"outcome_test.go:69: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsARestoreFailure,\"outcome_test.go:51: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/bit,\"result_test.go:53: probe = (1, {Verdict:invalid Restored:false CommandExit:1 FocusedExit:<nil>}), want command exit 0, verdict bit, restored true\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeResolvesTheSubjectFromTheWorkingDirectory,\"probe_test.go:294: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresAFileWithoutTrailingNewline,\"outcome_test.go:152: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,notrailing.txt,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestC… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresOnInterrupt/baseline,\"baseline_test.go:186: Go test start 1 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresOnInterrupt/mutated_run,\"baseline_test.go:186: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultEnvelopeRefusals,\"result_test.go:133: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultSchemaVersion,\"result_test.go:94: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRunsTheBaselineFirst,\"baseline_test.go:70: exit = 1, want 0\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/kit,\"subject_check_test.go:86: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,system,^Test… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeUnwrapBitesAndRestores,\"unwrap_test.go:23: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,unwrap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeWritesNoRecord,\"outcome_test.go:166: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantRequiresTheCompleteAnswer,\"result_test.go:208: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantSharesTheReaderRules,\"result_test.go:232: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestReadResultRefusesAMutantSet,\"result_test.go:199: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockAllowsAnotherSubject,\"lock_test.go:107: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,other.txt,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockIgnoresAFileWithoutAHolder,\"lock_test.go:131: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockKeepsItsFilesPrivate,\"lock_test.go:250: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_once,\"lock_test.go:148: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReleasesOnInterrupt,\"lock_test.go:115: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_once,\"lock_test.go:200: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapsMultilineBytes,\"edit_test.go:18: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegati… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
              "digest": "sha256:062a51f3143ce191923a762d43a44c2c142712dd3f5db5995394044c61d4a43f",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestClampGoldenBuildsOnce,318,1,0,null,yes\n  2,bit,swap,failed,196,TestClampCopiesAreIsolated,286,1,0,null,yes\n  3,bit,insert,failed,110,TestClampCopyKeepsAFreshSeal,318,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/probe,all,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,79195\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,76146,318\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampGoldenBuildsOnce,\"golden_test.go:119: the fixture publisher ran 55 times, want 1\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,32669,286\nfailures[196]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopiesAreIsolated,\"golden_test.go:126: subject bytes = \\\"package probefixture\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMakeOneMutant,\"combined_test.go:17: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^Tes… (545 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/ambiguous,\"combined_test.go:34: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,\\\\\\\\\\\"\\\\\\\\\\\"… (461 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/empty_pair,\"combined_test.go:39: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/start_bytes,\"combined_test.go:29: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap,refused,0,null,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegative$,pa… (589 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseOverlap,\"combined_test.go:47: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,\\\\\\\\\\\"\\\\\\\\\\\"… (444 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--delete,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--insert-before,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--omit-file,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--unwrap,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--with-file,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--with-file_without_--with,\"combined_test.go:73: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteRemovesTheSubjectForTheRun,\"delete_test.go:41: /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture… (354 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_and_restore,\"delete_test.go:103: /dev/shm/AVONEO/t/TestDeleteReportsARestoreFailurefailed_removal_and_restore4049586931/002/probe/repo-567104589/1792634926202973261-2144253/clamp.go = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Cla… (448 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_restored,\"delete_test.go:108: chmod /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_name_prefix,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_passing_mutated_run,\"verdict_test.go:60: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/another_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/prose_check,\"verdict_test.go:67: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/the_named_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/baseline_refusal,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/build_failure,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/prose_bite,\"verdict_test.go:101: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/silent_run,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/two_failed_tests,\"verdict_test.go:77: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/delete_with_with,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/delete_with_with-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/insert-before_without_text,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_combined,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_delete,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_omit-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/swap_without_its_with,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/with-file_without_swap,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/with_and_with-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/ambiguous_anchor,\"edit_test.go:50: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with,\"edit_test.go:38: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,insert,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with-file,\"edit_test.go:45: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,insert,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_control_byte_in_the_first_build_line,\"verdict_test.go:137: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_run-binary_build_script_failure_at_the_baseline,\"verdict_test.go:152: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/an_unused_import,\"verdict_test.go:119: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/inherited,\"set_test.go:110: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/owned,\"set_test.go:110: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_restore-failed,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_silent,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/two_bites,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryCap,\"set_input_test.go:290: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryMatchesTheSingleForm,\"set_test.go:31: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusalNamesTheEntry,\"set_input_test.go:310: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusesProseExpect,\"set_input_test.go:203: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusesSelectionTokens,\"set_input_test.go:181: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryTakesACombinedMutant,\"set_test.go:69: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetExpectIsPerEntry,\"set_test.go:59: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetGradesEachEntry,\"set_test.go:21: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissIsPerEntry,\"set_test.go:78: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetNamesABuildFailure,\"set_test.go:357: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetOfMissesStartsNoRun,\"set_test.go:87: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetOfProseBuildsNothing,\"set_test.go:265: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetPrintsEachReport,\"set_test.go:335: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/inherited,\"set_test.go:250: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/owned,\"set_test.go:250: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRestoresAfterEachMutant,\"set_test.go:155: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReusesTheBuildOutsideTheClosure,\"set_test.go:236: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRowsKeepTheirSchema,\"set_test.go:284: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:131: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetSharesOneBaseline,\"set_test.go:96: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/directory-sync,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/failed_restore,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/rename,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAPreserveFailure,\"set_test.go:222: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnARestoreFailure,\"set_test.go:172: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnASignalDuringARestore,\"baseline_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,substring-miss,0,null,0,null,1,null,untouched\\\\\\\\n  2,bit,swap,f… (962 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnUnreadableGate,\"lock_test.go:298: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (363 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsUnderAGateRun,\"lock_test.go:286: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (356 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/held_by_another_probe,\"lock_test.go:266: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/tried_during_mutant_2,\"lock_test.go:273: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesRepeatedPackages,\"set_test.go:366: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsExcludesTopLevelMutation,\"set_input_test.go:212: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsFileRefusesHostilePaths,\"set_input_test.go:229: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsFileRefusesMalformedJSON,\"set_input_test.go:259: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/out_of_range,\"edit_test.go:66: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation occurrence out of range — the old string matches 2 times, --occurrence is 3\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/second_match,\"edit_test.go:60: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_-1,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_0,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_99999999999999999999,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_x,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmissionThatEmptiesTheSubjectRefuses/combined,\"delete_test.go:138: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empties the subject — use --delete to remove clamp.go for the run\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmissionThatEmptiesTheSubjectRefuses/omit-file,\"delete_test.go:138: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empties the subject — use --delete to remove clamp.go for the run\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesAnUnreadableInputBeforeRunning/--omit-file,\"omit_file_test.go:77: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesAnUnreadableInputBeforeRunning--omit-file2737955549/005/mut… (270 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesAnUnreadableInputBeforeRunning/--with-file,\"omit_file_test.go:77: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesAnUnreadableInputBeforeRunning--with-file43956786/005/mutat… (268 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_absent,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_absent1411704031/005/… (269 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_directory,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_direct4008031769/005/… (274 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_empty,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_empty2723731046/005/m… (267 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_oversized,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_oversi322350743/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_special,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_specia1647278768/005/… (277 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_symlink,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_symlin2171921339/005/… (272 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_absent,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_absent4234421728/005/… (269 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_directory,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_direct1967178204/005/… (274 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_empty,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_empty18947531/005/mut… (265 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_oversized,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_oversi120057229/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_special,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_specia3710825161/005/… (277 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_symlink,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_symlin698694105/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesZeroAndMultipleMatchesSeparately/multiple_matches,\"omit_file_test.go:97: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesZeroAndMultipleMatchesSeparately/zero_matches,\"omit_file_test.go:90: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,null,untouched\\\\\\\\ns… (446 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileUsesExactBytesAndReportsOmit,\"omit_file_test.go:20: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,omit,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first … (247 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeBitesWhenTheFocusedTestFails,\"probe_test.go:195: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeBuildsTheMutantRunBinary,\"subject_check_test.go:109: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeCarriesTheFocusedRunTables,\"probe_test.go:247: exit = 1, want 0\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheRestoreFailureSuffix,\"result_test.go:68: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheTestsThatRan,\"probe_test.go:278: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"selection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegative$,passed,2\\\\\\\\n\\\" and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeFailedRestore,\"replacement_test.go:172: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenNoTestRuns,\"probe_test.go:236: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,no-test-run,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenTheMutationDoesNotCompile,\"probe_test.go:225: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,build-failed,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsSilentWhenTheFocusedTestPasses,\"probe_test.go:215: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  silent,clamp.go,swap,passed,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationPreparationFailure,\"replacement_test.go:205: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoveryFailure,\"replacement_test.go:100: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-close,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-open,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-sync,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/focused_run,\"replacement_test.go:143: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/mutation_failure,\"replacement_test.go:143: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCheckSelection,\"outcome_test.go:236: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the check selection row \\\"selection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,line-routing,^TestRootConformance$,passed,1\\\\\\\\n\\\" and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesThePackageSelection,\"outcome_test.go:247: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection… (334 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeOmitsTheMatchOnce,\"probe_test.go:207: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,omit,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbePackagesShareOneBuildPerRun,\"package_test.go:59: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first an… (245 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbePreservesTheSubjectBeforeTheRun,\"outcome_test.go:27: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_directory,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_special_file,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_symlink,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/absent,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/outside_the_repository,\"refusal_test.go:127: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnAmbiguousMutation/no_match,\"refusal_test.go:61: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nsel… (444 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnAmbiguousMutation/several_matches,\"refusal_test.go:61: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnEmptyLaterPackage,\"refusal_test.go:176: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnEmptyMutation,\"refusal_test.go:71: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empty — --with equals --swap\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUngradedProseSubject/matched,\"subject_check_test.go:65: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUngradedProseSubject/substring_miss,\"subject_check_test.go:65: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUnknownCheckBeforeAnyWrite,\"refusal_test.go:184: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesUnderALiveGateRun,\"refusal_test.go:90: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: gate execution in progress — wait for the gate run to finish before you mutate the tree\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenPreservationFails,\"refusal_test.go:209: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenTheMutationWriteFails,\"refusal_test.go:223: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:31: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection[… (346 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReplacement,\"replacement_test.go:49: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsAReadBackMismatch,\"outcome_test.go:69: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  restore-failed,clamp.go,swap,failed,1,TestClampNegative,no\\\\\\\\n… (255 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsARestoreFailure,\"outcome_test.go:51: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  restore-failed,clamp.go,swap,failed,1,TestClampNegative,no\\\\\\\\n… (255 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/bit,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/red_baseline,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/substring_miss,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResolvesTheSubjectFromTheWorkingDirectory,\"probe_test.go:294: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultEnvelopeRefusals,\"result_test.go:133: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultSchemaVersion,\"result_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeSubjectRefusalSkipsReplacement,\"replacement_test.go:185: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheProseCheck,\"subject_check_test.go:43: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/kit,\"subject_check_test.go:79: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/not_the_kit,\"subject_check_test.go:95: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUnwrapBitesAndRestores,\"unwrap_test.go:6: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_omit-file,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_unwrap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_delete,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/no_mutation,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/no_selection,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_unwrap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/run_with_a_check_other_than_system,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/two_selections,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/unwrap_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/with_only_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeWritesNoRecord,\"outcome_test.go:166: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantRequiresTheCompleteAnswer,\"result_test.go:208: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantSharesTheReaderRules,\"result_test.go:232: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadResultRefusesAMutantSet,\"result_test.go:199: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_line-routing,\"subject_check_test.go:136: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_system_--run_^TestClampNegative$,\"subject_check_test.go:136: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/unknown_check,\"subject_check_test.go:125: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockFollowsTheGateLock,\"lock_test.go:59: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegativ… (399 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockIgnoresAFileWithoutAHolder,\"lock_test.go:131: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockKeepsItsFilesPrivate,\"lock_test.go:250: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReadsOnlyContentionAsInUse,\"lock_test.go:184: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_on_every_attempt,\"lock_test.go:162: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_once,\"lock_test.go:148: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesALinkAtTheLockPath,\"lock_test.go:234: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesASecondProbe,\"lock_test.go:44: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegativ… (399 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesAnUnusableDirectory,\"lock_test.go:175: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReleasesOnInterrupt,\"lock_test.go:115: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_on_every_attempt,\"lock_test.go:214: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_once,\"lock_test.go:200: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapEqualToOldRefusesEmpty,\"edit_test.go:28: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empty — --with-file equals --swap\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapsMultilineBytes,\"edit_test.go:18: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestZeroMatchesStayASubstringMiss/combined,\"edit_test.go:118: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\n… (447 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestZeroMatchesStayASubstringMiss/occurrence,\"edit_test.go:118: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselec… (442 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,153196,318\nfailures[110]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopyKeepsAFreshSeal,\"golden_test.go:133: the copy's run binary refuses its own root: selected Bench executable \\\"/dev/shm/3G6OMR/t/TestClampCopyKeepsAFreshSeal4239075508/001/bin/bench\\\" does not match source \\\"/dev/shm/3G6OMR/t/TestClampCopyKeepsAFreshSeal42390755… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMakeOneMutant,\"combined_test.go:17: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap+omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^T… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/start_bytes,\"combined_test.go:29: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestCl… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteRemovesTheSubjectForTheRun,\"delete_test.go:40: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_and_restore,\"delete_test.go:95: exit = 1, want 2\",8\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_restored,\"delete_test.go:111: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/restore_proves_the_mode,\"delete_test.go:83: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_name_prefix,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_passing_mutated_run,\"verdict_test.go:63: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/another_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/the_named_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/build_failure,\"verdict_test.go:97: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/silent_run,\"verdict_test.go:97: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/two_failed_tests,\"verdict_test.go:80: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with,\"edit_test.go:38: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,insert,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with-file,\"edit_test.go:45: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,insert,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_control_byte_in_the_first_build_line,\"verdict_test.go:148: stdout = \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/an_unused_import,\"verdict_test.go:128: stdout = \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/inherited,\"set_test.go:116: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/owned,\"set_test.go:116: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_restore-failed,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_silent,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/two_bites,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryMatchesTheSingleForm,\"set_test.go:44: single form of entry 1 = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryTakesACombinedMutant,\"set_test.go:73: stdout = (\\\"mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap+swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\nselection[1]{… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetExpectIsPerEntry,\"set_test.go:64: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetGradesEachEntry,\"set_test.go:24: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissIsPerEntry,\"set_test.go:81: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,omit,substring-miss,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,bas… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/baseline,\"baseline_test.go:248: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/interrupt,\"baseline_test.go:239: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/restore_failure,\"baseline_test.go:256: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetNamesABuildFailure,\"set_test.go:361: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetPrintsEachReport,\"set_test.go:347: stdout = (\\\"mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,omit,… (1020 bytes)\",4\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/inherited,\"set_test.go:257: stdout = \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/owned,\"set_test.go:257: stdout = \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRestoresAfterEachMutant,\"set_test.go:161: exit = 1, want 0\",9\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReusesTheBuildOutsideTheClosure,\"set_test.go:240: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRowsKeepTheirSchema,\"set_test.go:300: selection row = \\\"  package,./,^TestClampNegative$,refused,0\\\", want \\\"  package,./,^TestClampNegative$,passed,0\\\"\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:147: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetSharesOneBaseline,\"set_test.go:100: exit = 1, want 0\",10\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/directory-sync,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/failed_restore,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/rename,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAPreserveFailure,\"set_test.go:228: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnARestoreFailure,\"set_test.go:177: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnASignalDuringARestore,\"baseline_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnInterruptedBuild,\"baseline_test.go:222: no process created /dev/shm/3G6OMR/t/TestMutantSetStopsOnAnInterruptedBuild3031501550/004/parked within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnUnreadableGate,\"lock_test.go:298: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnInterrupt,\"baseline_test.go:201: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsUnderAGateRun,\"lock_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/tried_during_mutant_2,\"lock_test.go:273: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesRepeatedPackages,\"set_test.go:374: exit = 1, want 0\",9\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/second_match,\"edit_test.go:60: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegati… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileUsesExactBytesAndReportsOmit,\"omit_file_test.go:20: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeBitesWhenTheFocusedTestFails,\"probe_test.go:195: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeBuildsTheMutantRunBinary,\"subject_check_test.go:113: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,system,^Tes… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeCarriesTheFocusedRunTables,\"probe_test.go:247: exit = 1, want 0\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheRestoreFailureSuffix,\"result_test.go:74: restore failure = (1, {Verdict:invalid Restored:false CommandExit:1 FocusedExit:<nil>}), want command 2, focused 1, unrestored\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheTestsThatRan,\"probe_test.go:278: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeFailedRestore,\"replacement_test.go:178: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestCla… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeForwardsFullToTheFocusedRun/baseline,\"baseline_test.go:160: --full did not reach the run:\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeForwardsFullToTheFocusedRun/mutated_run,\"baseline_test.go:160: --full did not reach the run:\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenNoTestRuns,\"probe_test.go:236: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenTheMutationDoesNotCompile,\"probe_test.go:225: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0\\\\\\\\… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsSilentWhenTheFocusedTestPasses,\"probe_test.go:215: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoveryFailure,\"replacement_test.go:110: exit = 1, want 2\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-close,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-open,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-sync,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/focused_run,\"replacement_test.go:151: exit = 1, want 2\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/mutation_failure,\"replacement_test.go:151: exit = 1, want 2\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCheckSelection,\"outcome_test.go:236: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,line-routing,^Tes… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCopyWhenTheRenderAndTheRestoreFail,\"outcome_test.go:106: stdout = (\\\"error: unrepresentable TOON cell — toon: unsupported control character U+0007 in string\\\\\\\\n\\\", 1), want \\\"error: unrepresentable TOON cell — \\\" first and 2\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesThePackageSelection,\"outcome_test.go:247: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeOmitsTheMatchOnce,\"probe_test.go:207: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbePackagesShareOneBuildPerRun,\"package_test.go:59: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./clamp ./other,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbePreservesTheSubjectBeforeTheRun,\"outcome_test.go:27: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesARedBaseline,\"baseline_test.go:103: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenPreservationFails,\"refusal_test.go:213: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenTheMutationWriteFails,\"refusal_test.go:230: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRendersASubjectWithASpace,\"outcome_test.go:142: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,my pkg/clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^Test… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:31: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./clamp ./other,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReplacement,\"replacement_test.go:71: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClamp… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineRefusal,\"baseline_test.go:137: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineWithoutAVerdict/baseline-build-failed,\"baseline_test.go:127: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineWithoutAVerdict/baseline-no-test-run,\"baseline_test.go:127: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsAReadBackMismatch,\"outcome_test.go:69: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsARestoreFailure,\"outcome_test.go:51: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/bit,\"result_test.go:53: probe = (1, {Verdict:invalid Restored:false CommandExit:1 FocusedExit:<nil>}), want command exit 0, verdict bit, restored true\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeResolvesTheSubjectFromTheWorkingDirectory,\"probe_test.go:294: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresAFileWithoutTrailingNewline,\"outcome_test.go:152: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,notrailing.txt,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestC… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresOnInterrupt/baseline,\"baseline_test.go:186: Go test start 1 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresOnInterrupt/mutated_run,\"baseline_test.go:186: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultEnvelopeRefusals,\"result_test.go:133: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultSchemaVersion,\"result_test.go:94: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRunsTheBaselineFirst,\"baseline_test.go:70: exit = 1, want 0\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/kit,\"subject_check_test.go:86: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,system,^Test… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeUnwrapBitesAndRestores,\"unwrap_test.go:23: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,unwrap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeWritesNoRecord,\"outcome_test.go:166: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantRequiresTheCompleteAnswer,\"result_test.go:208: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantSharesTheReaderRules,\"result_test.go:232: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestReadResultRefusesAMutantSet,\"result_test.go:199: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockAllowsAnotherSubject,\"lock_test.go:107: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,other.txt,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockIgnoresAFileWithoutAHolder,\"lock_test.go:131: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockKeepsItsFilesPrivate,\"lock_test.go:250: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_once,\"lock_test.go:148: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReleasesOnInterrupt,\"lock_test.go:115: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_once,\"lock_test.go:200: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapsMultilineBytes,\"edit_test.go:18: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegati… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01-r1-t01-copy-isolation",
          "performer": "claude:lgp-t01-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
            "digest": "sha256:062a51f3143ce191923a762d43a44c2c142712dd3f5db5995394044c61d4a43f",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestClampGoldenBuildsOnce,318,1,0,null,yes\n  2,bit,swap,failed,196,TestClampCopiesAreIsolated,286,1,0,null,yes\n  3,bit,insert,failed,110,TestClampCopyKeepsAFreshSeal,318,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/probe,all,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,79195\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,76146,318\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampGoldenBuildsOnce,\"golden_test.go:119: the fixture publisher ran 55 times, want 1\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,32669,286\nfailures[196]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopiesAreIsolated,\"golden_test.go:126: subject bytes = \\\"package probefixture\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMakeOneMutant,\"combined_test.go:17: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^Tes… (545 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/ambiguous,\"combined_test.go:34: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,\\\\\\\\\\\"\\\\\\\\\\\"… (461 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/empty_pair,\"combined_test.go:39: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/start_bytes,\"combined_test.go:29: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap,refused,0,null,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegative$,pa… (589 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseOverlap,\"combined_test.go:47: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,\\\\\\\\\\\"\\\\\\\\\\\"… (444 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--delete,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--insert-before,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--omit-file,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--unwrap,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--with-file,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--with-file_without_--with,\"combined_test.go:73: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteRemovesTheSubjectForTheRun,\"delete_test.go:41: /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture… (354 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_and_restore,\"delete_test.go:103: /dev/shm/AVONEO/t/TestDeleteReportsARestoreFailurefailed_removal_and_restore4049586931/002/probe/repo-567104589/1792634926202973261-2144253/clamp.go = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Cla… (448 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_restored,\"delete_test.go:108: chmod /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_name_prefix,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_passing_mutated_run,\"verdict_test.go:60: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/another_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/prose_check,\"verdict_test.go:67: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/the_named_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/baseline_refusal,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/build_failure,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/prose_bite,\"verdict_test.go:101: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/silent_run,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/two_failed_tests,\"verdict_test.go:77: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/delete_with_with,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/delete_with_with-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/insert-before_without_text,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_combined,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_delete,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_omit-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/swap_without_its_with,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/with-file_without_swap,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/with_and_with-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/ambiguous_anchor,\"edit_test.go:50: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with,\"edit_test.go:38: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,insert,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with-file,\"edit_test.go:45: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,insert,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_control_byte_in_the_first_build_line,\"verdict_test.go:137: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_run-binary_build_script_failure_at_the_baseline,\"verdict_test.go:152: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/an_unused_import,\"verdict_test.go:119: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/inherited,\"set_test.go:110: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/owned,\"set_test.go:110: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_restore-failed,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_silent,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/two_bites,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryCap,\"set_input_test.go:290: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryMatchesTheSingleForm,\"set_test.go:31: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusalNamesTheEntry,\"set_input_test.go:310: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusesProseExpect,\"set_input_test.go:203: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusesSelectionTokens,\"set_input_test.go:181: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryTakesACombinedMutant,\"set_test.go:69: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetExpectIsPerEntry,\"set_test.go:59: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetGradesEachEntry,\"set_test.go:21: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissIsPerEntry,\"set_test.go:78: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetNamesABuildFailure,\"set_test.go:357: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetOfMissesStartsNoRun,\"set_test.go:87: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetOfProseBuildsNothing,\"set_test.go:265: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetPrintsEachReport,\"set_test.go:335: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/inherited,\"set_test.go:250: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/owned,\"set_test.go:250: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRestoresAfterEachMutant,\"set_test.go:155: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReusesTheBuildOutsideTheClosure,\"set_test.go:236: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRowsKeepTheirSchema,\"set_test.go:284: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:131: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetSharesOneBaseline,\"set_test.go:96: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/directory-sync,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/failed_restore,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/rename,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAPreserveFailure,\"set_test.go:222: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnARestoreFailure,\"set_test.go:172: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnASignalDuringARestore,\"baseline_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,substring-miss,0,null,0,null,1,null,untouched\\\\\\\\n  2,bit,swap,f… (962 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnUnreadableGate,\"lock_test.go:298: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (363 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsUnderAGateRun,\"lock_test.go:286: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (356 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/held_by_another_probe,\"lock_test.go:266: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/tried_during_mutant_2,\"lock_test.go:273: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesRepeatedPackages,\"set_test.go:366: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsExcludesTopLevelMutation,\"set_input_test.go:212: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsFileRefusesHostilePaths,\"set_input_test.go:229: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsFileRefusesMalformedJSON,\"set_input_test.go:259: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/out_of_range,\"edit_test.go:66: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation occurrence out of range — the old string matches 2 times, --occurrence is 3\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/second_match,\"edit_test.go:60: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_-1,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_0,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_99999999999999999999,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_x,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmissionThatEmptiesTheSubjectRefuses/combined,\"delete_test.go:138: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empties the subject — use --delete to remove clamp.go for the run\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmissionThatEmptiesTheSubjectRefuses/omit-file,\"delete_test.go:138: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empties the subject — use --delete to remove clamp.go for the run\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesAnUnreadableInputBeforeRunning/--omit-file,\"omit_file_test.go:77: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesAnUnreadableInputBeforeRunning--omit-file2737955549/005/mut… (270 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesAnUnreadableInputBeforeRunning/--with-file,\"omit_file_test.go:77: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesAnUnreadableInputBeforeRunning--with-file43956786/005/mutat… (268 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_absent,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_absent1411704031/005/… (269 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_directory,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_direct4008031769/005/… (274 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_empty,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_empty2723731046/005/m… (267 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_oversized,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_oversi322350743/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_special,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_specia1647278768/005/… (277 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_symlink,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_symlin2171921339/005/… (272 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_absent,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_absent4234421728/005/… (269 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_directory,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_direct1967178204/005/… (274 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_empty,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_empty18947531/005/mut… (265 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_oversized,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_oversi120057229/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_special,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_specia3710825161/005/… (277 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_symlink,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_symlin698694105/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesZeroAndMultipleMatchesSeparately/multiple_matches,\"omit_file_test.go:97: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesZeroAndMultipleMatchesSeparately/zero_matches,\"omit_file_test.go:90: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,null,untouched\\\\\\\\ns… (446 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileUsesExactBytesAndReportsOmit,\"omit_file_test.go:20: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,omit,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first … (247 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeBitesWhenTheFocusedTestFails,\"probe_test.go:195: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeBuildsTheMutantRunBinary,\"subject_check_test.go:109: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeCarriesTheFocusedRunTables,\"probe_test.go:247: exit = 1, want 0\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheRestoreFailureSuffix,\"result_test.go:68: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheTestsThatRan,\"probe_test.go:278: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"selection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegative$,passed,2\\\\\\\\n\\\" and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeFailedRestore,\"replacement_test.go:172: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenNoTestRuns,\"probe_test.go:236: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,no-test-run,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenTheMutationDoesNotCompile,\"probe_test.go:225: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,build-failed,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsSilentWhenTheFocusedTestPasses,\"probe_test.go:215: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  silent,clamp.go,swap,passed,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationPreparationFailure,\"replacement_test.go:205: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoveryFailure,\"replacement_test.go:100: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-close,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-open,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-sync,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/focused_run,\"replacement_test.go:143: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/mutation_failure,\"replacement_test.go:143: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCheckSelection,\"outcome_test.go:236: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the check selection row \\\"selection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,line-routing,^TestRootConformance$,passed,1\\\\\\\\n\\\" and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesThePackageSelection,\"outcome_test.go:247: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection… (334 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeOmitsTheMatchOnce,\"probe_test.go:207: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,omit,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbePackagesShareOneBuildPerRun,\"package_test.go:59: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first an… (245 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbePreservesTheSubjectBeforeTheRun,\"outcome_test.go:27: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_directory,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_special_file,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_symlink,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/absent,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/outside_the_repository,\"refusal_test.go:127: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnAmbiguousMutation/no_match,\"refusal_test.go:61: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nsel… (444 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnAmbiguousMutation/several_matches,\"refusal_test.go:61: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnEmptyLaterPackage,\"refusal_test.go:176: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnEmptyMutation,\"refusal_test.go:71: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empty — --with equals --swap\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUngradedProseSubject/matched,\"subject_check_test.go:65: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUngradedProseSubject/substring_miss,\"subject_check_test.go:65: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUnknownCheckBeforeAnyWrite,\"refusal_test.go:184: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesUnderALiveGateRun,\"refusal_test.go:90: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: gate execution in progress — wait for the gate run to finish before you mutate the tree\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenPreservationFails,\"refusal_test.go:209: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenTheMutationWriteFails,\"refusal_test.go:223: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:31: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection[… (346 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReplacement,\"replacement_test.go:49: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsAReadBackMismatch,\"outcome_test.go:69: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  restore-failed,clamp.go,swap,failed,1,TestClampNegative,no\\\\\\\\n… (255 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsARestoreFailure,\"outcome_test.go:51: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  restore-failed,clamp.go,swap,failed,1,TestClampNegative,no\\\\\\\\n… (255 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/bit,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/red_baseline,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/substring_miss,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResolvesTheSubjectFromTheWorkingDirectory,\"probe_test.go:294: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultEnvelopeRefusals,\"result_test.go:133: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultSchemaVersion,\"result_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeSubjectRefusalSkipsReplacement,\"replacement_test.go:185: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheProseCheck,\"subject_check_test.go:43: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/kit,\"subject_check_test.go:79: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/not_the_kit,\"subject_check_test.go:95: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUnwrapBitesAndRestores,\"unwrap_test.go:6: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_omit-file,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_unwrap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_delete,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/no_mutation,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/no_selection,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_unwrap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/run_with_a_check_other_than_system,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/two_selections,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/unwrap_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/with_only_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeWritesNoRecord,\"outcome_test.go:166: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantRequiresTheCompleteAnswer,\"result_test.go:208: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantSharesTheReaderRules,\"result_test.go:232: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadResultRefusesAMutantSet,\"result_test.go:199: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_line-routing,\"subject_check_test.go:136: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_system_--run_^TestClampNegative$,\"subject_check_test.go:136: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/unknown_check,\"subject_check_test.go:125: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockFollowsTheGateLock,\"lock_test.go:59: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegativ… (399 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockIgnoresAFileWithoutAHolder,\"lock_test.go:131: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockKeepsItsFilesPrivate,\"lock_test.go:250: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReadsOnlyContentionAsInUse,\"lock_test.go:184: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_on_every_attempt,\"lock_test.go:162: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_once,\"lock_test.go:148: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesALinkAtTheLockPath,\"lock_test.go:234: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesASecondProbe,\"lock_test.go:44: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegativ… (399 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesAnUnusableDirectory,\"lock_test.go:175: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReleasesOnInterrupt,\"lock_test.go:115: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_on_every_attempt,\"lock_test.go:214: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_once,\"lock_test.go:200: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapEqualToOldRefusesEmpty,\"edit_test.go:28: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empty — --with-file equals --swap\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapsMultilineBytes,\"edit_test.go:18: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestZeroMatchesStayASubstringMiss/combined,\"edit_test.go:118: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\n… (447 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestZeroMatchesStayASubstringMiss/occurrence,\"edit_test.go:118: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselec… (442 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,153196,318\nfailures[110]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopyKeepsAFreshSeal,\"golden_test.go:133: the copy's run binary refuses its own root: selected Bench executable \\\"/dev/shm/3G6OMR/t/TestClampCopyKeepsAFreshSeal4239075508/001/bin/bench\\\" does not match source \\\"/dev/shm/3G6OMR/t/TestClampCopyKeepsAFreshSeal42390755… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMakeOneMutant,\"combined_test.go:17: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap+omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^T… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/start_bytes,\"combined_test.go:29: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestCl… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteRemovesTheSubjectForTheRun,\"delete_test.go:40: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_and_restore,\"delete_test.go:95: exit = 1, want 2\",8\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_restored,\"delete_test.go:111: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/restore_proves_the_mode,\"delete_test.go:83: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_name_prefix,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_passing_mutated_run,\"verdict_test.go:63: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/another_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/the_named_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/build_failure,\"verdict_test.go:97: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/silent_run,\"verdict_test.go:97: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/two_failed_tests,\"verdict_test.go:80: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with,\"edit_test.go:38: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,insert,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with-file,\"edit_test.go:45: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,insert,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_control_byte_in_the_first_build_line,\"verdict_test.go:148: stdout = \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/an_unused_import,\"verdict_test.go:128: stdout = \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/inherited,\"set_test.go:116: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/owned,\"set_test.go:116: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_restore-failed,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_silent,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/two_bites,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryMatchesTheSingleForm,\"set_test.go:44: single form of entry 1 = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryTakesACombinedMutant,\"set_test.go:73: stdout = (\\\"mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap+swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\nselection[1]{… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetExpectIsPerEntry,\"set_test.go:64: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetGradesEachEntry,\"set_test.go:24: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissIsPerEntry,\"set_test.go:81: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,omit,substring-miss,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,bas… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/baseline,\"baseline_test.go:248: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/interrupt,\"baseline_test.go:239: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/restore_failure,\"baseline_test.go:256: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetNamesABuildFailure,\"set_test.go:361: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetPrintsEachReport,\"set_test.go:347: stdout = (\\\"mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,omit,… (1020 bytes)\",4\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/inherited,\"set_test.go:257: stdout = \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/owned,\"set_test.go:257: stdout = \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRestoresAfterEachMutant,\"set_test.go:161: exit = 1, want 0\",9\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReusesTheBuildOutsideTheClosure,\"set_test.go:240: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRowsKeepTheirSchema,\"set_test.go:300: selection row = \\\"  package,./,^TestClampNegative$,refused,0\\\", want \\\"  package,./,^TestClampNegative$,passed,0\\\"\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:147: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetSharesOneBaseline,\"set_test.go:100: exit = 1, want 0\",10\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/directory-sync,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/failed_restore,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/rename,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAPreserveFailure,\"set_test.go:228: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnARestoreFailure,\"set_test.go:177: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnASignalDuringARestore,\"baseline_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnInterruptedBuild,\"baseline_test.go:222: no process created /dev/shm/3G6OMR/t/TestMutantSetStopsOnAnInterruptedBuild3031501550/004/parked within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnUnreadableGate,\"lock_test.go:298: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnInterrupt,\"baseline_test.go:201: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsUnderAGateRun,\"lock_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/tried_during_mutant_2,\"lock_test.go:273: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesRepeatedPackages,\"set_test.go:374: exit = 1, want 0\",9\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/second_match,\"edit_test.go:60: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegati… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileUsesExactBytesAndReportsOmit,\"omit_file_test.go:20: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeBitesWhenTheFocusedTestFails,\"probe_test.go:195: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeBuildsTheMutantRunBinary,\"subject_check_test.go:113: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,system,^Tes… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeCarriesTheFocusedRunTables,\"probe_test.go:247: exit = 1, want 0\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheRestoreFailureSuffix,\"result_test.go:74: restore failure = (1, {Verdict:invalid Restored:false CommandExit:1 FocusedExit:<nil>}), want command 2, focused 1, unrestored\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheTestsThatRan,\"probe_test.go:278: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeFailedRestore,\"replacement_test.go:178: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestCla… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeForwardsFullToTheFocusedRun/baseline,\"baseline_test.go:160: --full did not reach the run:\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeForwardsFullToTheFocusedRun/mutated_run,\"baseline_test.go:160: --full did not reach the run:\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenNoTestRuns,\"probe_test.go:236: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenTheMutationDoesNotCompile,\"probe_test.go:225: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0\\\\\\\\… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsSilentWhenTheFocusedTestPasses,\"probe_test.go:215: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoveryFailure,\"replacement_test.go:110: exit = 1, want 2\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-close,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-open,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-sync,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/focused_run,\"replacement_test.go:151: exit = 1, want 2\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/mutation_failure,\"replacement_test.go:151: exit = 1, want 2\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCheckSelection,\"outcome_test.go:236: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,line-routing,^Tes… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCopyWhenTheRenderAndTheRestoreFail,\"outcome_test.go:106: stdout = (\\\"error: unrepresentable TOON cell — toon: unsupported control character U+0007 in string\\\\\\\\n\\\", 1), want \\\"error: unrepresentable TOON cell — \\\" first and 2\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesThePackageSelection,\"outcome_test.go:247: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeOmitsTheMatchOnce,\"probe_test.go:207: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbePackagesShareOneBuildPerRun,\"package_test.go:59: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./clamp ./other,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbePreservesTheSubjectBeforeTheRun,\"outcome_test.go:27: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesARedBaseline,\"baseline_test.go:103: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenPreservationFails,\"refusal_test.go:213: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenTheMutationWriteFails,\"refusal_test.go:230: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRendersASubjectWithASpace,\"outcome_test.go:142: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,my pkg/clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^Test… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:31: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./clamp ./other,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReplacement,\"replacement_test.go:71: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClamp… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineRefusal,\"baseline_test.go:137: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineWithoutAVerdict/baseline-build-failed,\"baseline_test.go:127: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineWithoutAVerdict/baseline-no-test-run,\"baseline_test.go:127: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsAReadBackMismatch,\"outcome_test.go:69: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsARestoreFailure,\"outcome_test.go:51: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/bit,\"result_test.go:53: probe = (1, {Verdict:invalid Restored:false CommandExit:1 FocusedExit:<nil>}), want command exit 0, verdict bit, restored true\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeResolvesTheSubjectFromTheWorkingDirectory,\"probe_test.go:294: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresAFileWithoutTrailingNewline,\"outcome_test.go:152: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,notrailing.txt,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestC… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresOnInterrupt/baseline,\"baseline_test.go:186: Go test start 1 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresOnInterrupt/mutated_run,\"baseline_test.go:186: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultEnvelopeRefusals,\"result_test.go:133: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultSchemaVersion,\"result_test.go:94: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRunsTheBaselineFirst,\"baseline_test.go:70: exit = 1, want 0\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/kit,\"subject_check_test.go:86: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,system,^Test… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeUnwrapBitesAndRestores,\"unwrap_test.go:23: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,unwrap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeWritesNoRecord,\"outcome_test.go:166: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantRequiresTheCompleteAnswer,\"result_test.go:208: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantSharesTheReaderRules,\"result_test.go:232: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestReadResultRefusesAMutantSet,\"result_test.go:199: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockAllowsAnotherSubject,\"lock_test.go:107: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,other.txt,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockIgnoresAFileWithoutAHolder,\"lock_test.go:131: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockKeepsItsFilesPrivate,\"lock_test.go:250: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_once,\"lock_test.go:148: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReleasesOnInterrupt,\"lock_test.go:115: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_once,\"lock_test.go:200: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapsMultilineBytes,\"edit_test.go:18: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegati… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
              "digest": "sha256:062a51f3143ce191923a762d43a44c2c142712dd3f5db5995394044c61d4a43f",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestClampGoldenBuildsOnce,318,1,0,null,yes\n  2,bit,swap,failed,196,TestClampCopiesAreIsolated,286,1,0,null,yes\n  3,bit,insert,failed,110,TestClampCopyKeepsAFreshSeal,318,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/probe,all,passed,1\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,79195\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,76146,318\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampGoldenBuildsOnce,\"golden_test.go:119: the fixture publisher ran 55 times, want 1\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,32669,286\nfailures[196]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopiesAreIsolated,\"golden_test.go:126: subject bytes = \\\"package probefixture\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMakeOneMutant,\"combined_test.go:17: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^Tes… (545 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/ambiguous,\"combined_test.go:34: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,\\\\\\\\\\\"\\\\\\\\\\\"… (461 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/empty_pair,\"combined_test.go:39: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/start_bytes,\"combined_test.go:29: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap,refused,0,null,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegative$,pa… (589 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseOverlap,\"combined_test.go:47: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,\\\\\\\\\\\"\\\\\\\\\\\"… (444 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--delete,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--insert-before,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--omit-file,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--unwrap,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--with-file,\"combined_test.go:64: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsRefuseSingleForms/--with-file_without_--with,\"combined_test.go:73: subject bytes = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\n… (314 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteRemovesTheSubjectForTheRun,\"delete_test.go:41: /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n > 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\", want \\\"package probefixture… (354 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_and_restore,\"delete_test.go:103: /dev/shm/AVONEO/t/TestDeleteReportsARestoreFailurefailed_removal_and_restore4049586931/002/probe/repo-567104589/1792634926202973261-2144253/clamp.go = \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Cla… (448 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_restored,\"delete_test.go:108: chmod /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_name_prefix,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_passing_mutated_run,\"verdict_test.go:60: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/another_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/prose_check,\"verdict_test.go:67: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/the_named_test,\"verdict_test.go:52: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/baseline_refusal,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/build_failure,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/prose_bite,\"verdict_test.go:101: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/silent_run,\"verdict_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/two_failed_tests,\"verdict_test.go:77: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/delete_with_with,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/delete_with_with-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/insert-before_without_text,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_combined,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_delete,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/occurrence_with_omit-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/swap_without_its_with,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/with-file_without_swap,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestFormMatrixRefusesPartnerFaults/with_and_with-file,\"edit_test.go:97: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/ambiguous_anchor,\"edit_test.go:50: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with,\"edit_test.go:38: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,insert,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with-file,\"edit_test.go:45: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,insert,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_control_byte_in_the_first_build_line,\"verdict_test.go:137: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_run-binary_build_script_failure_at_the_baseline,\"verdict_test.go:152: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/an_unused_import,\"verdict_test.go:119: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/inherited,\"set_test.go:110: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/owned,\"set_test.go:110: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_restore-failed,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_silent,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/two_bites,\"set_test.go:320: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryCap,\"set_input_test.go:290: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryMatchesTheSingleForm,\"set_test.go:31: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusalNamesTheEntry,\"set_input_test.go:310: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusesProseExpect,\"set_input_test.go:203: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryRefusesSelectionTokens,\"set_input_test.go:181: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryTakesACombinedMutant,\"set_test.go:69: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetExpectIsPerEntry,\"set_test.go:59: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetGradesEachEntry,\"set_test.go:21: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissIsPerEntry,\"set_test.go:78: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetNamesABuildFailure,\"set_test.go:357: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetOfMissesStartsNoRun,\"set_test.go:87: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetOfProseBuildsNothing,\"set_test.go:265: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetPrintsEachReport,\"set_test.go:335: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/inherited,\"set_test.go:250: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/owned,\"set_test.go:250: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRestoresAfterEachMutant,\"set_test.go:155: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReusesTheBuildOutsideTheClosure,\"set_test.go:236: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRowsKeepTheirSchema,\"set_test.go:284: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:131: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetSharesOneBaseline,\"set_test.go:96: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/directory-sync,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/failed_restore,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/rename,\"set_test.go:202: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAPreserveFailure,\"set_test.go:222: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnARestoreFailure,\"set_test.go:172: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnASignalDuringARestore,\"baseline_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,substring-miss,0,null,0,null,1,null,untouched\\\\\\\\n  2,bit,swap,f… (962 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnUnreadableGate,\"lock_test.go:298: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (363 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsUnderAGateRun,\"lock_test.go:286: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (356 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/held_by_another_probe,\"lock_test.go:266: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/tried_during_mutant_2,\"lock_test.go:273: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,bit,swap,f… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesRepeatedPackages,\"set_test.go:366: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsExcludesTopLevelMutation,\"set_input_test.go:212: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsFileRefusesHostilePaths,\"set_input_test.go:229: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantsFileRefusesMalformedJSON,\"set_input_test.go:259: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/out_of_range,\"edit_test.go:66: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation occurrence out of range — the old string matches 2 times, --occurrence is 3\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/second_match,\"edit_test.go:60: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_-1,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_0,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_99999999999999999999,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/usage_x,\"edit_test.go:72: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmissionThatEmptiesTheSubjectRefuses/combined,\"delete_test.go:138: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empties the subject — use --delete to remove clamp.go for the run\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmissionThatEmptiesTheSubjectRefuses/omit-file,\"delete_test.go:138: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empties the subject — use --delete to remove clamp.go for the run\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesAnUnreadableInputBeforeRunning/--omit-file,\"omit_file_test.go:77: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesAnUnreadableInputBeforeRunning--omit-file2737955549/005/mut… (270 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesAnUnreadableInputBeforeRunning/--with-file,\"omit_file_test.go:77: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesAnUnreadableInputBeforeRunning--with-file43956786/005/mutat… (268 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_absent,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_absent1411704031/005/… (269 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_directory,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_direct4008031769/005/… (274 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_empty,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_empty2723731046/005/m… (267 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_oversized,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_oversi322350743/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_special,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_specia1647278768/005/… (277 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--omit-file_symlink,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--omit-file_symlin2171921339/005/… (272 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_absent,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_absent4234421728/005/… (269 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_directory,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_direct1967178204/005/… (274 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_empty,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_empty18947531/005/mut… (265 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_oversized,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_oversi120057229/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_special,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_specia3710825161/005/… (277 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesInvalidInputsBeforeRunning/--with-file_symlink,\"omit_file_test.go:60: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation file unavailable — /dev/shm/AVONEO/t/TestOmitFileRefusesInvalidInputsBeforeRunning--with-file_symlin698694105/005/m… (271 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesZeroAndMultipleMatchesSeparately/multiple_matches,\"omit_file_test.go:97: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileRefusesZeroAndMultipleMatchesSeparately/zero_matches,\"omit_file_test.go:90: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,substring-miss,0,null,untouched\\\\\\\\ns… (446 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileUsesExactBytesAndReportsOmit,\"omit_file_test.go:20: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,omit,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first … (247 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeBitesWhenTheFocusedTestFails,\"probe_test.go:195: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeBuildsTheMutantRunBinary,\"subject_check_test.go:109: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeCarriesTheFocusedRunTables,\"probe_test.go:247: exit = 1, want 0\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheRestoreFailureSuffix,\"result_test.go:68: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheTestsThatRan,\"probe_test.go:278: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"selection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegative$,passed,2\\\\\\\\n\\\" and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeFailedRestore,\"replacement_test.go:172: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenNoTestRuns,\"probe_test.go:236: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,no-test-run,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenTheMutationDoesNotCompile,\"probe_test.go:225: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,build-failed,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsSilentWhenTheFocusedTestPasses,\"probe_test.go:215: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  silent,clamp.go,swap,passed,0,null,yes\\\\\\\\n\\\" first and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationPreparationFailure,\"replacement_test.go:205: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoveryFailure,\"replacement_test.go:100: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-close,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-open,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-sync,\"replacement_test.go:81: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/focused_run,\"replacement_test.go:143: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/mutation_failure,\"replacement_test.go:143: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCheckSelection,\"outcome_test.go:236: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the check selection row \\\"selection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,line-routing,^TestRootConformance$,passed,1\\\\\\\\n\\\" and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesThePackageSelection,\"outcome_test.go:247: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection… (334 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeOmitsTheMatchOnce,\"probe_test.go:207: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,omit,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbePackagesShareOneBuildPerRun,\"package_test.go:59: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first an… (245 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbePreservesTheSubjectBeforeTheRun,\"outcome_test.go:27: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_directory,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_special_file,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/a_symlink,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/absent,\"refusal_test.go:122: open /dev/shm/AVONEO/t/probe-golden-3804286733/repo/clamp.go: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAHostileSubject/outside_the_repository,\"refusal_test.go:127: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnAmbiguousMutation/no_match,\"refusal_test.go:61: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nsel… (444 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnAmbiguousMutation/several_matches,\"refusal_test.go:61: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnEmptyLaterPackage,\"refusal_test.go:176: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnEmptyMutation,\"refusal_test.go:71: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empty — --with equals --swap\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUngradedProseSubject/matched,\"subject_check_test.go:65: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUngradedProseSubject/substring_miss,\"subject_check_test.go:65: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesAnUnknownCheckBeforeAnyWrite,\"refusal_test.go:184: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesUnderALiveGateRun,\"refusal_test.go:90: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: gate execution in progress — wait for the gate run to finish before you mutate the tree\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenPreservationFails,\"refusal_test.go:209: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenTheMutationWriteFails,\"refusal_test.go:223: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:31: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection[… (346 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReplacement,\"replacement_test.go:49: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsAReadBackMismatch,\"outcome_test.go:69: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  restore-failed,clamp.go,swap,failed,1,TestClampNegative,no\\\\\\\\n… (255 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsARestoreFailure,\"outcome_test.go:51: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  restore-failed,clamp.go,swap,failed,1,TestClampNegative,no\\\\\\\\n… (255 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/bit,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/red_baseline,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/substring_miss,\"result_test.go:48: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResolvesTheSubjectFromTheWorkingDirectory,\"probe_test.go:294: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and… (244 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultEnvelopeRefusals,\"result_test.go:133: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultSchemaVersion,\"result_test.go:94: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeSubjectRefusalSkipsReplacement,\"replacement_test.go:185: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheProseCheck,\"subject_check_test.go:43: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/kit,\"subject_check_test.go:79: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/not_the_kit,\"subject_check_test.go:95: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUnwrapBitesAndRestores,\"unwrap_test.go:6: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_omit-file,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/delete_with_unwrap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_delete,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/insert-before_with_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/no_mutation,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/no_selection,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_omit,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit-file_with_unwrap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/omit_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/run_with_a_check_other_than_system,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/two_selections,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/unwrap_with_replacement,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeUsageNamesOneSelectionAndOneMutation/with_only_swap,\"refusal_test.go:166: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeWritesNoRecord,\"outcome_test.go:166: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantRequiresTheCompleteAnswer,\"result_test.go:208: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantSharesTheReaderRules,\"result_test.go:232: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadResultRefusesAMutantSet,\"result_test.go:199: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_line-routing,\"subject_check_test.go:136: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/--check_system_--run_^TestClampNegative$,\"subject_check_test.go:136: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSelectionPrecedesTheSubstringMiss/unknown_check,\"subject_check_test.go:125: CopyFS repo/pipe: invalid argument\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockFollowsTheGateLock,\"lock_test.go:59: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegativ… (399 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockIgnoresAFileWithoutAHolder,\"lock_test.go:131: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockKeepsItsFilesPrivate,\"lock_test.go:250: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReadsOnlyContentionAsInUse,\"lock_test.go:184: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_on_every_attempt,\"lock_test.go:162: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_once,\"lock_test.go:148: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesALinkAtTheLockPath,\"lock_test.go:234: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesASecondProbe,\"lock_test.go:44: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegativ… (399 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRefusesAnUnusableDirectory,\"lock_test.go:175: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReleasesOnInterrupt,\"lock_test.go:115: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_on_every_attempt,\"lock_test.go:214: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want the lock refusal and 1\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_once,\"lock_test.go:200: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and … (243 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapEqualToOldRefusesEmpty,\"edit_test.go:28: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"error: probe mutation empty — --with-file equals --swap\\\\\\\\n\\\", 1)\",1\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapsMultilineBytes,\"edit_test.go:18: stdout = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\n\\\" first and 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestZeroMatchesStayASubstringMiss/combined,\"edit_test.go:118: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+omit,substring-miss,0,null,untouched\\\\\\\\n… (447 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestZeroMatchesStayASubstringMiss/occurrence,\"edit_test.go:118: refusal = (\\\"error: probe subject unavailable — clamp.go is absent\\\\\\\\n\\\", 1), want (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,substring-miss,0,null,untouched\\\\\\\\nselec… (442 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,153196,318\nfailures[110]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestClampCopyKeepsAFreshSeal,\"golden_test.go:133: the copy's run binary refuses its own root: selected Bench executable \\\"/dev/shm/3G6OMR/t/TestClampCopyKeepsAFreshSeal4239075508/001/bin/bench\\\" does not match source \\\"/dev/shm/3G6OMR/t/TestClampCopyKeepsAFreshSeal42390755… (668 bytes)\",1\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMakeOneMutant,\"combined_test.go:17: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap+omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^T… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestCombinedEditsMatchOnce/start_bytes,\"combined_test.go:29: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap+swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestCl… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteRemovesTheSubjectForTheRun,\"delete_test.go:40: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_and_restore,\"delete_test.go:95: exit = 1, want 2\",8\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/failed_removal_restored,\"delete_test.go:111: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestDeleteReportsARestoreFailure/restore_proves_the_mode,\"delete_test.go:83: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,delete,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_name_prefix,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_passing_mutated_run,\"verdict_test.go:63: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/a_subtest_of_the_named_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/another_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/the_named_test,\"verdict_test.go:55: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/build_failure,\"verdict_test.go:97: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/silent_run,\"verdict_test.go:97: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestFirstFailureNamesTheFirstFailedTest/two_failed_tests,\"verdict_test.go:80: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with,\"edit_test.go:38: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,insert,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestInsertBeforePutsTextBeforeTheAnchor/with-file,\"edit_test.go:45: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,insert,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/a_control_byte_in_the_first_build_line,\"verdict_test.go:148: stdout = \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantBuildFailureNamesItself/an_unused_import,\"verdict_test.go:128: stdout = \\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/inherited,\"set_test.go:116: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetBaselineEndsTheSet/owned,\"set_test.go:116: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_restore-failed,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/bit_and_silent,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetCommandExit/two_bites,\"set_test.go:327: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryMatchesTheSingleForm,\"set_test.go:44: single form of entry 1 = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetEntryTakesACombinedMutant,\"set_test.go:73: stdout = (\\\"mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap+swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\nselection[1]{… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetExpectIsPerEntry,\"set_test.go:64: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetGradesEachEntry,\"set_test.go:24: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissIsPerEntry,\"set_test.go:81: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,omit,substring-miss,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,bas… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/baseline,\"baseline_test.go:248: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/interrupt,\"baseline_test.go:239: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetMissKeepsItsCause/restore_failure,\"baseline_test.go:256: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetNamesABuildFailure,\"set_test.go:361: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetPrintsEachReport,\"set_test.go:347: stdout = (\\\"mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,omit,… (1020 bytes)\",4\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/inherited,\"set_test.go:257: stdout = \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReportsTheBuildCount/owned,\"set_test.go:257: stdout = \\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,b… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRestoresAfterEachMutant,\"set_test.go:161: exit = 1, want 0\",9\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetReusesTheBuildOutsideTheClosure,\"set_test.go:240: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRowsKeepTheirSchema,\"set_test.go:300: selection row = \\\"  package,./,^TestClampNegative$,refused,0\\\", want \\\"  package,./,^TestClampNegative$,passed,0\\\"\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:147: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetSharesOneBaseline,\"set_test.go:100: exit = 1, want 0\",10\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/directory-sync,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/failed_restore,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAMutationFailure/rename,\"set_test.go:207: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAPreserveFailure,\"set_test.go:228: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnARestoreFailure,\"set_test.go:177: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnASignalDuringARestore,\"baseline_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnInterruptedBuild,\"baseline_test.go:222: no process created /dev/shm/3G6OMR/t/TestMutantSetStopsOnAnInterruptedBuild3031501550/004/parked within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnAnUnreadableGate,\"lock_test.go:298: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsOnInterrupt,\"baseline_test.go:201: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetStopsUnderAGateRun,\"lock_test.go:286: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesOneSubjectLock/tried_during_mutant_2,\"lock_test.go:273: stdout = (\\\"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,swap… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetTakesRepeatedPackages,\"set_test.go:374: exit = 1, want 0\",9\n  github.com/gibbonmi/bench/internal/probe,TestOccurrenceSelectsOneMatch/second_match,\"edit_test.go:60: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegati… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestOmitFileUsesExactBytesAndReportsOmit,\"omit_file_test.go:20: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeBitesWhenTheFocusedTestFails,\"probe_test.go:195: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeBuildsTheMutantRunBinary,\"subject_check_test.go:113: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,system,^Tes… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeCarriesTheFocusedRunTables,\"probe_test.go:247: exit = 1, want 0\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheRestoreFailureSuffix,\"result_test.go:74: restore failure = (1, {Verdict:invalid Restored:false CommandExit:1 FocusedExit:<nil>}), want command 2, focused 1, unrestored\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeCountsTheTestsThatRan,\"probe_test.go:278: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeFailedRestore,\"replacement_test.go:178: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestCla… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeForwardsFullToTheFocusedRun/baseline,\"baseline_test.go:160: --full did not reach the run:\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeForwardsFullToTheFocusedRun/mutated_run,\"baseline_test.go:160: --full did not reach the run:\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenNoTestRuns,\"probe_test.go:236: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsInvalidWhenTheMutationDoesNotCompile,\"probe_test.go:225: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused,0\\\\\\\\… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeIsSilentWhenTheFocusedTestPasses,\"probe_test.go:215: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoveryFailure,\"replacement_test.go:110: exit = 1, want 2\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-close,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-open,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeMutationRecoverySuccess/directory-sync,\"replacement_test.go:90: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/focused_run,\"replacement_test.go:151: exit = 1, want 2\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesACopyTheEncoderRefuses/mutation_failure,\"replacement_test.go:151: exit = 1, want 2\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCheckSelection,\"outcome_test.go:236: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,line-routing,^Tes… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesTheCopyWhenTheRenderAndTheRestoreFail,\"outcome_test.go:106: stdout = (\\\"error: unrepresentable TOON cell — toon: unsupported control character U+0007 in string\\\\\\\\n\\\", 1), want \\\"error: unrepresentable TOON cell — \\\" first and 2\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeNamesThePackageSelection,\"outcome_test.go:247: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeOmitsTheMatchOnce,\"probe_test.go:207: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,omit,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbePackagesShareOneBuildPerRun,\"package_test.go:59: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./clamp ./other,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbePreservesTheSubjectBeforeTheRun,\"outcome_test.go:27: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesARedBaseline,\"baseline_test.go:103: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenPreservationFails,\"refusal_test.go:213: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRefusesWhenTheMutationWriteFails,\"refusal_test.go:230: refusal = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRendersASubjectWithASpace,\"outcome_test.go:142: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,my pkg/clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^Test… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:31: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./clamp ./other,… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReplacement,\"replacement_test.go:71: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClamp… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineRefusal,\"baseline_test.go:137: Go test starts = 0, want only the baseline's one\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineWithoutAVerdict/baseline-build-failed,\"baseline_test.go:127: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,all,refused… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsABaselineWithoutAVerdict/baseline-no-test-run,\"baseline_test.go:127: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampN… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsAReadBackMismatch,\"outcome_test.go:69: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsARestoreFailure,\"outcome_test.go:51: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNeg… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeReportsNormalizedExecutionExit/bit,\"result_test.go:53: probe = (1, {Verdict:invalid Restored:false CommandExit:1 FocusedExit:<nil>}), want command exit 0, verdict bit, restored true\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeResolvesTheSubjectFromTheWorkingDirectory,\"probe_test.go:294: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresAFileWithoutTrailingNewline,\"outcome_test.go:152: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,notrailing.txt,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestC… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresOnInterrupt/baseline,\"baseline_test.go:186: Go test start 1 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeRestoresOnInterrupt/mutated_run,\"baseline_test.go:186: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultEnvelopeRefusals,\"result_test.go:133: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeResultSchemaVersion,\"result_test.go:94: probe = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeRunsTheBaselineFirst,\"baseline_test.go:70: exit = 1, want 0\",8\n  github.com/gibbonmi/bench/internal/probe,TestProbeTakesTheSystemCheck/kit,\"subject_check_test.go:86: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  check,system,^Test… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeUnwrapBitesAndRestores,\"unwrap_test.go:23: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,unwrap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNe… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestProbeWritesNoRecord,\"outcome_test.go:166: exit = 1, want 0\",1\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantRequiresTheCompleteAnswer,\"result_test.go:208: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestReadMutantSharesTheReaderRules,\"result_test.go:232: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestReadResultRefusesAMutantSet,\"result_test.go:199: stdout = (\\\"mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\\\\\\\\n  1,invalid,swap,baseline-refused,0,null,0,null,1,null,untouched\\\\\\\\n  2,invalid,om… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockAllowsAnotherSubject,\"lock_test.go:107: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,other.txt,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNega… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockIgnoresAFileWithoutAHolder,\"lock_test.go:131: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockKeepsItsFilesPrivate,\"lock_test.go:250: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRechecksTheIdentity/replaced_once,\"lock_test.go:148: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockReleasesOnInterrupt,\"lock_test.go:115: Go test start 2 did not happen within 23s\",1\n  github.com/gibbonmi/bench/internal/probe,TestSubjectLockRetriesAnAbsentLockDirectory/removed_once,\"lock_test.go:200: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegat… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/probe,TestWithFileSwapsMultilineBytes,\"edit_test.go:18: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  invalid,clamp.go,swap,baseline-refused,0,null,untouched\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./,^TestClampNegati… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01-r1-t01-testrepo",
          "performer": "claude:lgp-t01-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
            "digest": "sha256:83c2e6ea94a547e3103f445ba129cdc60e07c38b32fae55b5cfbaa52b7821eb1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testrepo,pass,39,14\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-testrepo",
          "command": "bench test --package ./internal/testrepo",
          "exit_code": 0
        },
        {
          "id": "t01-r1-t01-adopt",
          "performer": "claude:lgp-t01-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
            "digest": "sha256:1710becd2c39d588d1dace867adb7053a40d2db8ba4e9a2e57420a27a5b93e65",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,pass,24162,244\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-adopt",
          "command": "bench test --package ./internal/adopt",
          "exit_code": 0
        },
        {
          "id": "t01-r1-t01-copy-modes",
          "performer": "claude:lgp-t01-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
            "digest": "sha256:5d8acd17490c3d7b1fcc10c91188d0eadfc6fc4ef39510f0728ac06eb4209cff",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testrepo/copy_tree.go,swap,failed,1,TestCopyTreeKeepsModes,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testrepo,all,passed,14\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,316\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testrepo,fail,40,14\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testrepo,TestCopyTreeKeepsModes,\"copy_tree_test.go:54: a-private mode = -rw-r--r--, want -rw-------\",4\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t01-copy-modes",
          "command": "bench test --package ./internal/testrepo",
          "exit_code": 0,
          "probe": {
            "mutation": "drop the mode restore of the shared tree copy; the copy-mode test must red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t01-repair1-20261021",
              "digest": "sha256:5d8acd17490c3d7b1fcc10c91188d0eadfc6fc4ef39510f0728ac06eb4209cff",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testrepo/copy_tree.go,swap,failed,1,TestCopyTreeKeepsModes,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testrepo,all,passed,14\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,316\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testrepo,fail,40,14\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testrepo,TestCopyTreeKeepsModes,\"copy_tree_test.go:54: a-private mode = -rw-r--r--, want -rw-------\",4\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c1-standards-r0",
          "performer": "claude:lgp-c1-standards-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c1-standards-20261021",
            "digest": "sha256:5c1740c213710a93174fc305147dce9b1a53c799fc217340f50f299f5a6a4c79",
            "excerpt": "Standards axis, LGP-C1 (7f585da44..5cd784915): 4 findings, 3 accepted. Worst: S1. S1 (ask-user, confidence 8): ownedRunBinary (internal/testreport/selection_once_test.go:12-52) repeats the worktree testRunSelector (internal/worktree/test_run_test.go:16-90) and is a third, weaker derivation of the runbinary.Factory.ReuseOrOwn inherit-or-own policy (it skips requireSourceRoot). Consultants split (sonnet: accept; opus: extract); the fable tie-breaker chose one Factory.Select(ctx, sourceRoot, raw, present) method and one exported once-holder inside internal/runbinary. Disposition: auto-fix, ticket 02 repair. S2 (ask-user, confidence 6): the probe golden repeats the adopt golden harness and drops its mode-restoring walk. Consultant consensus: extract only the mode-preserving tree copy into one internal/testrepo helper that adopt and probe share. Disposition: auto-fix, ticket 01 repair. S3 (auto-fix, confidence 8): the fixture doc comment at internal/probe/probe_test.go:15-16 still names probetest.Repository. Disposition: auto-fix, ticket 01 repair. S4 (no-op, confidence 7): the counters are an atomic or the fake builder, not a mutex-guarded log; the count requirement holds and no standard fails.\n"
          },
          "axis": "Standards",
          "base": "7f585da4489289b34df9cf3beeef2c59de8edb6d",
          "tip": "5cd7849154fcfed5f255cc99a6812d2f38d8cc7a",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "c1-spec-r0",
          "performer": "claude:lgp-c1-spec-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c1-spec-20261021",
            "digest": "sha256:116692ba6e841ab8e083001f86d0eff9169ec530dd6a33fb76525c3c21da4c3c",
            "excerpt": "Spec axis, LGP-C1 (7f585da44..5cd784915): 4 findings, 3 accepted. Worst: P1. Rows LGP1-5, LGP7, LGP8, LGP37 met; LGP6 and LGP34 partial. P1 (ask-user, confidence 6): story 8 asks for a real owned build, but the inherited BENCH_RUN_BINARY wins in the gate, so no testreport test builds. Consultant consensus: amend story 8 to name the surviving real-build test TestBuildLeavesTheWrapperManifestUntouched (internal/runbinary/build_script_test.go:26-57). Disposition: spec prose amendment. P2 (auto-fix, confidence 7): LGP34 before-and-after --full times were only in commit messages. Coordinator disposition: evidence-only; this entry records them: internal/probe 74925 ms before, 44304 ms after (dc17fbba8); internal/testreport 79679 ms before, 77267 ms after (5cd784915), inside the recorded spread. P3 (ask-user, confidence 6): os.CopyFS does not keep modes, against spec line 118. Disposition: closed by the S2 repair (mode-preserving copy). P4 (ask-user, confidence 5): LGP6's why-it-catches names a bypass that TestOwnedSelectionBuildsOnce does not see. Consultant consensus: amend the clause to \"A helper that builds on each installOwnedSelection call invokes the builder a second time.\" Disposition: spec prose amendment.\n"
          },
          "axis": "Spec",
          "base": "7f585da4489289b34df9cf3beeef2c59de8edb6d",
          "tip": "5cd7849154fcfed5f255cc99a6812d2f38d8cc7a",
          "finding_ids": [
            "P1",
            "P3",
            "P4"
          ],
          "supersedes": []
        },
        {
          "id": "c1-coverage-r0",
          "performer": "claude:lgp-c1-coverage-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a209794fd9556f6ed4a4da99debb782345e2bbda",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c1-coverage-20261021",
            "digest": "sha256:fb86bfc5852ad18de1a6cfb54c8be370fa8396d1eea503be815a2f1b2f058a9b",
            "excerpt": "Coverage axis, LGP-C1 (7f585da44..5cd784915): 1 finding, accepted. Worst: C1. C1 (auto-fix, confidence 8): the inherited branch of ownedRunBinary (valid inherited wins; present invalid refuses, spec line 126) has no test; probe swap of the inherited test for a never-true one read silent, restored yes. Disposition: auto-fix as the chunk's one hardening cycle, in the ticket 02 repair (the Factory.Select extraction of S1 carries the branch tests). Other probes: copy shares the golden executable (silent, restored yes; advice: seal holds no path); TestMain cleanup removed (silent, restored yes; advice). Refuted: hard-link sharing, re-entry helpers building a golden, inherited selection removal at Close. File modes: advice, joined to S2/P3.\n"
          },
          "axis": "Coverage",
          "base": "7f585da4489289b34df9cf3beeef2c59de8edb6d",
          "tip": "5cd7849154fcfed5f255cc99a6812d2f38d8cc7a",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "c1-standards-r1",
          "performer": "claude:lgp-c1-standards-r1-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c1-standards-r1-20261021",
            "digest": "sha256:e45b1222c788e2b028071fa303e3c1ec5c92bc83ae076ddfbc6c63fc3c3dd3f9",
            "excerpt": "Standards axis, LGP-C1 confirming round (repair delta 5cd784915..be0aca0c5): S1, S2, S3 closed, 0 new findings. S1 closed (9): Factory.Select holds the rule (internal/runbinary/runbinary.go:77-85), ReuseOrOwn calls it (:70-73), and runbinary.Shared (shared.go:12) serves the testreport and worktree selectors, each with its own capture. S2 closed (9): testrepo.CopyTree (internal/testrepo/copy_tree.go:15) holds the mode-preserving copy; both goldens call it. S3 closed (10): probe_test.go:15-16 names clampCopy. Advice only: a comment that Shared serves test runs; the worktree refusal table overlaps shared_test.go:204-227; the CopyTree destination comment overstates the limit.\n"
          },
          "axis": "Standards",
          "base": "7f585da4489289b34df9cf3beeef2c59de8edb6d",
          "tip": "be0aca0c526c177be15b360eeba46c7b57927ee8",
          "finding_ids": [],
          "supersedes": [
            "c1-standards-r0"
          ]
        },
        {
          "id": "c1-spec-r1",
          "performer": "claude:lgp-c1-spec-r1-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c1-spec-r1-20261021",
            "digest": "sha256:5eecc1c407a1406cc82ce8e72e85be5b178a36848fa3b4f9765baae49344e485",
            "excerpt": "Spec axis, LGP-C1 confirming round (repair delta 5cd784915..be0aca0c5): P1, P3, P4 closed. P1 closed (9): story 8 names TestBuildLeavesTheWrapperManifestUntouched, which runs the real Build on a kit copy (internal/runbinary/build_script_test.go:26-57). P3 closed (9): testrepo.CopyTree restores permission bits, and clampCopy uses it. P4 closed (9): the LGP6 row names the per-call build, and TestOwnedSelectionBuildsOnce proves it. ReuseOrOwn keeps its behavior (runbinary.go:70-85). P5 (auto-fix, confidence 8): spec lines 121, 126, and 128 cite runbinary.go and test_run_test.go lines that the repair shifted. Coordinator disposition: a prose-only spec correction in the next plan commit, before the ticket 03 dispatch, as the Land rule routes it; it meets no blocking condition of the bounded repair policy, so this round carries no open finding.\n"
          },
          "axis": "Spec",
          "base": "7f585da4489289b34df9cf3beeef2c59de8edb6d",
          "tip": "be0aca0c526c177be15b360eeba46c7b57927ee8",
          "finding_ids": [],
          "supersedes": [
            "c1-spec-r0"
          ]
        },
        {
          "id": "c1-coverage-r1",
          "performer": "claude:lgp-c1-coverage-r1-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccf2c687b713b6d6c0336145de2ed8fb1290cae5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c1-coverage-r1-20261021",
            "digest": "sha256:ceb33cd0d94dd3fac80b9cd91c9c227662fbdcc0ee56a0d085d86a3a4faab74b",
            "excerpt": "Coverage axis, LGP-C1 confirming round (repair delta 5cd784915..be0aca0c5): C1 closed (9), 0 new findings. internal/runbinary/shared_test.go:188-227 pins the valid-inherited and present-invalid branches through runbinary.Shared; worktree SB1-SB3 and the parallel journey test pass. Independent probes, all restored yes: Shared.Path argument s.Present to s.Present && s.Raw != \"\" bit (TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild/empty); CopyTree symlink skip removed bit (TestCopyTreeKeepsModes); CopyTree directory entries skipped bit (TestCopyTreeKeepsModes, dir mode). Advice: no test covers a source directory without the search bit.\n"
          },
          "axis": "Coverage",
          "base": "7f585da4489289b34df9cf3beeef2c59de8edb6d",
          "tip": "be0aca0c526c177be15b360eeba46c7b57927ee8",
          "finding_ids": [],
          "supersedes": [
            "c1-coverage-r0"
          ]
        }
      ]
    },
    {
      "id": "LGP-C2A",
      "base": "be0aca0c526c177be15b360eeba46c7b57927ee8",
      "tip": "4c4d1df002f9b734aa5db8a99032dabf669c7318",
      "plan_digest": "sha256:504f3662dc285101a10b12a57d0d3171265e21af7bf269af10240c19c8d2b397",
      "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
      "acceptance_rows": [
        "LGP9",
        "LGP10",
        "LGP11",
        "LGP12",
        "LGP13",
        "LGP14",
        "LGP15",
        "LGP34"
      ],
      "verification": [
        {
          "id": "t03-git",
          "performer": "claude:lgp-t03-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "7c41bdcf5553577385591dc3ef047af10544dd78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-author-20261021",
            "digest": "sha256:0f54ffd4eca704b2ebe4319d82723d6acbb5792f567d33c377de825338dbad78",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,1619,177\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-git",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "t03-poolkey",
          "performer": "claude:lgp-t03-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "7c41bdcf5553577385591dc3ef047af10544dd78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-author-20261021",
            "digest": "sha256:a24b67b6f05cbc52f5ff78d0590eb0daaf0b8f2634d6e44edb63e6fde41f0387",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/poolkey,pass,20,10\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-poolkey",
          "command": "bench test --package ./internal/poolkey",
          "exit_code": 0
        },
        {
          "id": "t03-root-conformance",
          "performer": "claude:lgp-t03-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "7c41bdcf5553577385591dc3ef047af10544dd78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-author-20261021",
            "digest": "sha256:074e9e751fd3bb2713950d6dc646462192b9789e806051eee3af4b805b0bba98",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16695,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t03-revalidation",
          "performer": "claude:lgp-t03-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "7c41bdcf5553577385591dc3ef047af10544dd78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-author-20261021",
            "digest": "sha256:53c0e8bc782e849741d39bdd77931a6afbfcb3939e9488598a8a17d54c005c74",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,1,TestRunRevalidatesACachedCommonDir,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,all,passed,177\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,399\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1590,177\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunRevalidatesACachedCommonDir,\"run_cache_test.go:92: second read = <nil>, want the symlink refusal of \\\"/dev/shm/QPTAG5/t/TestRunRevalidatesACachedCommonDir37012447/001/.git\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-revalidation",
          "command": "bench test --package ./internal/git",
          "exit_code": 0,
          "probe": {
            "mutation": "skip the validation on a cached common directory",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:lgp-t03-author-20261021",
              "digest": "sha256:53c0e8bc782e849741d39bdd77931a6afbfcb3939e9488598a8a17d54c005c74",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,1,TestRunRevalidatesACachedCommonDir,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,all,passed,177\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,399\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1590,177\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunRevalidatesACachedCommonDir,\"run_cache_test.go:92: second read = <nil>, want the symlink refusal of \\\"/dev/shm/QPTAG5/t/TestRunRevalidatesACachedCommonDir37012447/001/.git\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t03-r1-t03-git",
          "performer": "claude:lgp-t03-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021",
            "digest": "sha256:4699607975d155670eebfa347ec4b6ddbf0f350afc1a256f83201bbed94968b2",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,1610,179\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-git",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "t03-r1-t03-poolkey",
          "performer": "claude:lgp-t03-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021",
            "digest": "sha256:53ea3916008c2325461b3864e9c1e9c0a714d83afd492527ff4631581f105975",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/poolkey,pass,17,10\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-poolkey",
          "command": "bench test --package ./internal/poolkey",
          "exit_code": 0
        },
        {
          "id": "t03-r1-t03-root-conformance",
          "performer": "claude:lgp-t03-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021",
            "digest": "sha256:f9f9ee075f37ce2ea3136667c61d220588fdee357ec86a668219cb48659ae56c",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16655,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t03-r1-t03-revalidation",
          "performer": "claude:lgp-t03-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021",
            "digest": "sha256:b1d49d4a1c935a063930ffa6901a2fe3fc41e8e11a5c775c8d99ffa656438863",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestRunRevalidatesACachedCommonDir,179,1,0,null,yes\n  2,bit,swap,failed,1,TestRunKeysTheAdminDirApartFromTheCommonDir,179,1,0,null,yes\n  3,bit,swap,failed,1,TestRunStoresNoFailedAnswer,179,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/git,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1318\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1611,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunRevalidatesACachedCommonDir,\"run_cache_test.go:97: second read = <nil>, want the symlink refusal of \\\"/dev/shm/U3AKEW/t/TestRunRevalidatesACachedCommonDir2522912558/001/.git\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1628,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunKeysTheAdminDirApartFromTheCommonDir,\"run_cache_test.go:73: run.AdminDir(/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/002/linked) = \\\"/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/001/.git\\\", <nil>, want an answer apart from th… (341 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1614,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunStoresNoFailedAnswer,\"run_cache_test.go:161: common-directory starts = 1, want 2\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-revalidation",
          "command": "bench test --package ./internal/git",
          "exit_code": 0,
          "probe": {
            "mutation": "skip the validation on a cached common directory",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021",
              "digest": "sha256:b1d49d4a1c935a063930ffa6901a2fe3fc41e8e11a5c775c8d99ffa656438863",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestRunRevalidatesACachedCommonDir,179,1,0,null,yes\n  2,bit,swap,failed,1,TestRunKeysTheAdminDirApartFromTheCommonDir,179,1,0,null,yes\n  3,bit,swap,failed,1,TestRunStoresNoFailedAnswer,179,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/git,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1318\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1611,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunRevalidatesACachedCommonDir,\"run_cache_test.go:97: second read = <nil>, want the symlink refusal of \\\"/dev/shm/U3AKEW/t/TestRunRevalidatesACachedCommonDir2522912558/001/.git\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1628,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunKeysTheAdminDirApartFromTheCommonDir,\"run_cache_test.go:73: run.AdminDir(/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/002/linked) = \\\"/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/001/.git\\\", <nil>, want an answer apart from th… (341 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1614,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunStoresNoFailedAnswer,\"run_cache_test.go:161: common-directory starts = 1, want 2\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t03-r1-t03-gittest",
          "performer": "claude:lgp-t03-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021",
            "digest": "sha256:676ea7b4e7ab18ce01c6bc1d498d25cf0438b297351bc890a60ad90915efcd42",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gittest,pass,35,2\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-gittest",
          "command": "bench test --package ./internal/gittest",
          "exit_code": 0
        },
        {
          "id": "t03-r1-t03-admin-key",
          "performer": "claude:lgp-t03-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021",
            "digest": "sha256:b1d49d4a1c935a063930ffa6901a2fe3fc41e8e11a5c775c8d99ffa656438863",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestRunRevalidatesACachedCommonDir,179,1,0,null,yes\n  2,bit,swap,failed,1,TestRunKeysTheAdminDirApartFromTheCommonDir,179,1,0,null,yes\n  3,bit,swap,failed,1,TestRunStoresNoFailedAnswer,179,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/git,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1318\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1611,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunRevalidatesACachedCommonDir,\"run_cache_test.go:97: second read = <nil>, want the symlink refusal of \\\"/dev/shm/U3AKEW/t/TestRunRevalidatesACachedCommonDir2522912558/001/.git\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1628,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunKeysTheAdminDirApartFromTheCommonDir,\"run_cache_test.go:73: run.AdminDir(/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/002/linked) = \\\"/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/001/.git\\\", <nil>, want an answer apart from th… (341 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1614,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunStoresNoFailedAnswer,\"run_cache_test.go:161: common-directory starts = 1, want 2\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-admin-key",
          "command": "bench test --package ./internal/git",
          "exit_code": 0,
          "probe": {
            "mutation": "key Run.AdminDir under the common-directory location; the linked-worktree test must red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021",
              "digest": "sha256:b1d49d4a1c935a063930ffa6901a2fe3fc41e8e11a5c775c8d99ffa656438863",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestRunRevalidatesACachedCommonDir,179,1,0,null,yes\n  2,bit,swap,failed,1,TestRunKeysTheAdminDirApartFromTheCommonDir,179,1,0,null,yes\n  3,bit,swap,failed,1,TestRunStoresNoFailedAnswer,179,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/git,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1318\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1611,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunRevalidatesACachedCommonDir,\"run_cache_test.go:97: second read = <nil>, want the symlink refusal of \\\"/dev/shm/U3AKEW/t/TestRunRevalidatesACachedCommonDir2522912558/001/.git\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1628,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunKeysTheAdminDirApartFromTheCommonDir,\"run_cache_test.go:73: run.AdminDir(/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/002/linked) = \\\"/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/001/.git\\\", <nil>, want an answer apart from th… (341 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1614,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunStoresNoFailedAnswer,\"run_cache_test.go:161: common-directory starts = 1, want 2\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t03-r1-t03-failure-memo",
          "performer": "claude:lgp-t03-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021",
            "digest": "sha256:b1d49d4a1c935a063930ffa6901a2fe3fc41e8e11a5c775c8d99ffa656438863",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestRunRevalidatesACachedCommonDir,179,1,0,null,yes\n  2,bit,swap,failed,1,TestRunKeysTheAdminDirApartFromTheCommonDir,179,1,0,null,yes\n  3,bit,swap,failed,1,TestRunStoresNoFailedAnswer,179,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/git,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1318\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1611,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunRevalidatesACachedCommonDir,\"run_cache_test.go:97: second read = <nil>, want the symlink refusal of \\\"/dev/shm/U3AKEW/t/TestRunRevalidatesACachedCommonDir2522912558/001/.git\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1628,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunKeysTheAdminDirApartFromTheCommonDir,\"run_cache_test.go:73: run.AdminDir(/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/002/linked) = \\\"/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/001/.git\\\", <nil>, want an answer apart from th… (341 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1614,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunStoresNoFailedAnswer,\"run_cache_test.go:161: common-directory starts = 1, want 2\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t03-failure-memo",
          "command": "bench test --package ./internal/git",
          "exit_code": 0,
          "probe": {
            "mutation": "store a failed answer in Run.answer; the failed-query test must red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021",
              "digest": "sha256:b1d49d4a1c935a063930ffa6901a2fe3fc41e8e11a5c775c8d99ffa656438863",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestRunRevalidatesACachedCommonDir,179,1,0,null,yes\n  2,bit,swap,failed,1,TestRunKeysTheAdminDirApartFromTheCommonDir,179,1,0,null,yes\n  3,bit,swap,failed,1,TestRunStoresNoFailedAnswer,179,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/git,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1318\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1611,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunRevalidatesACachedCommonDir,\"run_cache_test.go:97: second read = <nil>, want the symlink refusal of \\\"/dev/shm/U3AKEW/t/TestRunRevalidatesACachedCommonDir2522912558/001/.git\\\"\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1628,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunKeysTheAdminDirApartFromTheCommonDir,\"run_cache_test.go:73: run.AdminDir(/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/002/linked) = \\\"/dev/shm/7NW27O/t/TestRunKeysTheAdminDirApartFromTheCommonDir2495499539/001/.git\\\", <nil>, want an answer apart from th… (341 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,1614,179\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunStoresNoFailedAnswer,\"run_cache_test.go:161: common-directory starts = 1, want 2\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c2a-standards-r0",
          "performer": "claude:lgp-c2a-standards-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c41bdcf5553577385591dc3ef047af10544dd78",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2a-standards-20261021",
            "digest": "sha256:061c61e961fc0a1c23fdac533dafda47684c5960922faeaa2209d83e40f8fd2f",
            "excerpt": "Standards axis, LGP-C2A (be0aca0c5..de63c6090): 1 finding, accepted. Worst: S1. S1 (auto-fix, confidence 4): the stub-git start-counting reader is derived twice, in gitStarts and loggedStub (internal/git/run_cache_test.go:15-29) and inline (internal/poolkey/run_test.go:17, :26-30); AGENTS.md keeps fixture harnesses single-sourced. Coordinator disposition: auto-fix in the ticket 03 repair; one start counter beside StubGitDir in internal/gittest, which tickets 04 and 05 also need. Clean: the query-then-validate step lives in Run.directory, the trim in locationAnswer, and the package forms build new(Run). Advice: the Run doc comment repeats the location list; a test file spells --git-common-dir, which the owner check permits in tests.\n"
          },
          "axis": "Standards",
          "base": "be0aca0c526c177be15b360eeba46c7b57927ee8",
          "tip": "de63c6090aa9328c5c0f131eba4498765dc69846",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "c2a-spec-r0",
          "performer": "claude:lgp-c2a-spec-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c41bdcf5553577385591dc3ef047af10544dd78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2a-spec-20261021",
            "digest": "sha256:49e0d1685c58262f840b32b40d2003ce6bd21a4030416d577c3fadab558e6de8",
            "excerpt": "Spec axis, LGP-C2A (be0aca0c5..de63c6090): rows LGP9-LGP15 met; LGP34 partial. P1 (auto-fix, confidence 8): the record held no before time for LGP34. Coordinator disposition: evidence-only; this entry records the ticket 03 author's bench test --full --package ./internal/git times: 1564 ms with 172 tests before, 1582 ms with 177 tests after (de63c6090), no change within the spread. LGP12 holds: Run exposes no ref reader, and the test interleaves the live ResolveCommit with run reads. Advice: LGP12 cannot catch a future run-level ref method; the Run comment states the rule.\n"
          },
          "axis": "Spec",
          "base": "be0aca0c526c177be15b360eeba46c7b57927ee8",
          "tip": "de63c6090aa9328c5c0f131eba4498765dc69846",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c2a-coverage-r0",
          "performer": "claude:lgp-c2a-coverage-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c41bdcf5553577385591dc3ef047af10544dd78",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2a-coverage-20261021",
            "digest": "sha256:9a2b62924bfe1266caf73dc22bb7db247d584d136f3cd9af5197e449880e7c71",
            "excerpt": "Coverage axis, LGP-C2A (be0aca0c5..de63c6090): 2 findings, accepted. Worst: C1. C1 (auto-fix, confidence 7): Run.AdminDir keyed under the common-directory location passes, because LGP10 runs in a primary checkout where both answers are equal; probe swap runQuery{adminDirLocation, root} for runQuery{commonDirLocation, root} read silent, restored yes. Ticket 04 reads both in one run. C2 (auto-fix, confidence 6): a failed first query cached in Run.answer passes, although run.go:77-78 says a failed start stores nothing; probe storing raw unconditionally read silent, restored yes. Disposition: both in the ticket 03 repair, as the one hardening cycle. Refuted: root spellings keyed apart by design, revalidation on every read, package forms uncached, Canonical delegates to CanonicalIn. Advice: a nil *Run panics; a relative root is keyed as text.\n"
          },
          "axis": "Coverage",
          "base": "be0aca0c526c177be15b360eeba46c7b57927ee8",
          "tip": "de63c6090aa9328c5c0f131eba4498765dc69846",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "c2a-standards-r1",
          "performer": "claude:lgp-c2a-standards-r1-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2a-standards-r1-20261021",
            "digest": "sha256:08acc06cc556e74aa8af04b6a05d48eaf090ed962b8e62b5aacc212281c944f5",
            "excerpt": "Standards axis, LGP-C2A confirming round (repair delta de63c6090..4c4d1df00): S1 closed (9), 0 new findings. One gittest.StubGitStarts(t, logPath, flag) at internal/gittest/gittest.go:140-155 serves internal/git/run_cache_test.go and internal/poolkey/run_test.go:24; the local gitStarts and the inline count are gone. Advice: the two new tests carry no coverage-row tag; the StubGitStarts doc names StubGitDir.\n"
          },
          "axis": "Standards",
          "base": "be0aca0c526c177be15b360eeba46c7b57927ee8",
          "tip": "4c4d1df002f9b734aa5db8a99032dabf669c7318",
          "finding_ids": [],
          "supersedes": [
            "c2a-standards-r0"
          ]
        },
        {
          "id": "c2a-spec-r1",
          "performer": "claude:lgp-c2a-spec-r1-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2a-spec-r1-20261021",
            "digest": "sha256:25db2b1e8b9276f0109ad48f46d2476c01122c0b2b921564461b572687b3f1d8",
            "excerpt": "Spec axis, LGP-C2A confirming round (repair delta de63c6090..4c4d1df00): LGP9-LGP15 met at 4c4d1df00, 0 new findings. The delta touches only internal/git/run_cache_test.go, internal/poolkey/run_test.go, and internal/gittest/gittest.go, so no production file of internal/git or internal/poolkey changes; each row's test keeps its name and its assertions. Advice: StubGitStarts counts lines, which matches one stub line per start; the poolkey failure message no longer prints the log.\n"
          },
          "axis": "Spec",
          "base": "be0aca0c526c177be15b360eeba46c7b57927ee8",
          "tip": "4c4d1df002f9b734aa5db8a99032dabf669c7318",
          "finding_ids": [],
          "supersedes": [
            "c2a-spec-r0"
          ]
        },
        {
          "id": "c2a-coverage-r1",
          "performer": "claude:lgp-c2a-coverage-r1-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4a16f6b13947814e1769e848c0eb30d91ff85547",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2a-coverage-r1-20261021",
            "digest": "sha256:1763df2a438eff32395f805c935ddd2a5138a6d1938d218273fe9dbca746fab6",
            "excerpt": "Coverage axis, LGP-C2A confirming round (repair delta de63c6090..4c4d1df00): C1 closed (9), C2 closed (8), 0 new findings. Independent probes, all restored yes: RootAt keyed as the admin directory bit (TestRunCachesTopLevelAndAdminDir); StubGitStarts counting every log line bit (TestRunCachesTheCommonDir and two more); a failed answer stored only for the admin and top-level locations read silent. Advice: answer has no per-location branch and t03-failure-memo bites, so a failed AdminDir case would only harden the test.\n"
          },
          "axis": "Coverage",
          "base": "be0aca0c526c177be15b360eeba46c7b57927ee8",
          "tip": "4c4d1df002f9b734aa5db8a99032dabf669c7318",
          "finding_ids": [],
          "supersedes": [
            "c2a-coverage-r0"
          ]
        }
      ]
    },
    {
      "id": "LGP-C2B",
      "base": "4c4d1df002f9b734aa5db8a99032dabf669c7318",
      "tip": "bb8a207e3b56c973ee1cd6fe043e4bc3167e6351",
      "plan_digest": "sha256:f97a2dd92aa8f25d612c082e08c62d69946a1bc987a3ae263cae84b7ac1b56b3",
      "source_digest": "fe4abb4dc1ffa973a0494d808d083b5bd8cd9d56",
      "acceptance_rows": [
        "LGP17",
        "LGP21",
        "LGP22",
        "LGP34"
      ],
      "verification": [
        {
          "id": "t04-pkg-cmd-bench",
          "performer": "claude:lgp-t04-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "fe4abb4dc1ffa973a0494d808d083b5bd8cd9d56",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t04-author-20261021",
            "digest": "sha256:c82afa71b7a56f8746febc38aade5b2f87feb21b940e08fe49ec99242ad27f48",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20456,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t04-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t04-worktree",
          "performer": "claude:lgp-t04-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "fe4abb4dc1ffa973a0494d808d083b5bd8cd9d56",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t04-author-20261021",
            "digest": "sha256:27a03bb313a7ed9ea56629a9326bab3b3ecd81eccd4b54035b924ce784e8d5b4",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,87545,1492\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t04-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t04-root-conformance",
          "performer": "claude:lgp-t04-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "fe4abb4dc1ffa973a0494d808d083b5bd8cd9d56",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t04-author-20261021",
            "digest": "sha256:4f2d445178d2799cf0aae733d65f3b90946ddbbc4a3012d0045428fcc331691a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16784,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t04-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t04-marker-run-value",
          "performer": "claude:lgp-t04-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "fe4abb4dc1ffa973a0494d808d083b5bd8cd9d56",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t04-author-20261021",
            "digest": "sha256:eadee60d96f1c7d95a8b846915d1da20c0521ff51bb777150fd4bcc691fd63c8",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/owner_marker.go,swap+swap+swap+swap,failed,1,TestOwnerMarkerReadsTheCommonDirOnce,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1492\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,368\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,88765,1492\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestOwnerMarkerReadsTheCommonDirOnce,\"journey_test.go:257: one owner-marker validation started the common-directory query 3 times, want 1\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t04-marker-run-value",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "skip the run value in validateOwnerMarker",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:lgp-t04-author-20261021",
              "digest": "sha256:eadee60d96f1c7d95a8b846915d1da20c0521ff51bb777150fd4bcc691fd63c8",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/owner_marker.go,swap+swap+swap+swap,failed,1,TestOwnerMarkerReadsTheCommonDirOnce,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1492\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,368\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,88765,1492\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestOwnerMarkerReadsTheCommonDirOnce,\"journey_test.go:257: one owner-marker validation started the common-directory query 3 times, want 1\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c2b-standards-r0",
          "performer": "claude:lgp-c2b-standards-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fe4abb4dc1ffa973a0494d808d083b5bd8cd9d56",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2b-standards-20261021",
            "digest": "sha256:76b0509cfdcfbbcbfaafb0780a4b6ce7086ea18adbfab4291749940c465262bc",
            "excerpt": "Standards axis, LGP-C2B (4c4d1df00..bb8a207e3): 0 findings. The canonical-root rule has one source, poolkey.CanonicalIn; canonicalRoot and canonicalRootIn are one-line delegates (internal/worktree/worktree.go:415, :418). validateOwnerMarker moved to owner_marker.go with only the run reads changed, and no second body remains; the thin form keeps subshell.go:225 and two test callers. Every touched comment matches the code. Advice: &git.Run{} against new(git.Run) is incidental text.\n"
          },
          "axis": "Standards",
          "base": "4c4d1df002f9b734aa5db8a99032dabf669c7318",
          "tip": "bb8a207e3b56c973ee1cd6fe043e4bc3167e6351",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c2b-spec-r0",
          "performer": "claude:lgp-c2b-spec-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fe4abb4dc1ffa973a0494d808d083b5bd8cd9d56",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2b-spec-20261021",
            "digest": "sha256:ca40dc581f2a548fcaecd761f10ad4b180befc7368929c7cd22eab78f7824769",
            "excerpt": "Spec axis, LGP-C2B (4c4d1df00..bb8a207e3): 0 findings. LGP17 met: validateOwnerMarkerIn reads canonicalRootIn, run.Worktrees, and run.CommonDir under one key; TestOwnerMarkerReadsTheCommonDirOnce (journey_test.go:247) asserts one start. Carrier met: ambient.run, newAmbient builds &git.Run{}, callAmbient calls newAmbient, no joins field holds a run, the three ownerMarkerRefusal callers build a run on the call line, and reauthorizeParsed takes a.run. Two-runs item met inside the LGP17 test. LGP21 and LGP22 met: 805 + 1 = 806 tests, 45 - 0 + 1 = 46 serial. LGP34: coordinator records the ticket 04 author's bench test --full --package ./internal/worktree times, 88859 ms (1491 tests) before and 86662 ms (1492 tests) after, inside the spread. AdminDir through the run is within spec line 136.\n"
          },
          "axis": "Spec",
          "base": "4c4d1df002f9b734aa5db8a99032dabf669c7318",
          "tip": "bb8a207e3b56c973ee1cd6fe043e4bc3167e6351",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c2b-coverage-r0",
          "performer": "claude:lgp-c2b-coverage-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fe4abb4dc1ffa973a0494d808d083b5bd8cd9d56",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2b-coverage-20261021",
            "digest": "sha256:106e033302437fec73c8c4b6f6f614aa09472cab4e419db2a24cdc4ad9ab5bdd",
            "excerpt": "Coverage axis, LGP-C2B (4c4d1df00..bb8a207e3): 0 findings. Probes, all restored yes: canonicalRootIn given a fresh run bit (LGP17, 2 starts); newAmbient run set to nil bit (panic and LGP17); an extra git.CommonDir read inserted in owner_marker.go bit (LGP17); reauthorize a.run swapped for a fresh run read silent, an equivalent mutant today. Refuted: shared runs across verb runs, joins reuse, a stale cache after a registration write, the single-read census. Advice: nothing pins reauthorize's a.run until a second reader joins the verb.\n"
          },
          "axis": "Coverage",
          "base": "4c4d1df002f9b734aa5db8a99032dabf669c7318",
          "tip": "bb8a207e3b56c973ee1cd6fe043e4bc3167e6351",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "LGP-C2C",
      "base": "bb8a207e3b56c973ee1cd6fe043e4bc3167e6351",
      "tip": "cd9889779c97c9d50b07158bce399e6b40612737",
      "plan_digest": "sha256:c29cdfbe2ba7223fba8e656138b5210901464cb9f4b555854aabd86698ec4368",
      "source_digest": "1bc5392d81da0bc27398f21fc9886a61f8bc59bf",
      "acceptance_rows": [
        "LGP16",
        "LGP21",
        "LGP22",
        "LGP34"
      ],
      "verification": [
        {
          "id": "t05-pkg-cmd-bench",
          "performer": "claude:lgp-t05-author2-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "1bc5392d81da0bc27398f21fc9886a61f8bc59bf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t05-author2-20261021",
            "digest": "sha256:8d01b84995fc9afeb28081c4b6b0ad96783be317c225f19157dc6742a88af7fd",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15832,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t05-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t05-worktree",
          "performer": "claude:lgp-t05-author2-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "1bc5392d81da0bc27398f21fc9886a61f8bc59bf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t05-author2-20261021",
            "digest": "sha256:d01dc353a141b2d6388caea7c184b370944d8f0f3dd59ad456492ec44904640e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,67191,1493\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t05-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t05-root-conformance",
          "performer": "claude:lgp-t05-author2-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "1bc5392d81da0bc27398f21fc9886a61f8bc59bf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t05-author2-20261021",
            "digest": "sha256:d4192c07261c0623c4496c10f09e57b8b1cba09596ba42b075cb7fc2993c236e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18894,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t05-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t05-plan-run-value",
          "performer": "claude:lgp-t05-author2-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "1bc5392d81da0bc27398f21fc9886a61f8bc59bf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:lgp-t05-author2-20261021",
            "digest": "sha256:3fd74654de1bfb9ec4676265fbebcc585d2a0ae08fb49aa4fe0733f8af415729",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/subshell.go,swap,failed,1,TestExplicitPlanReadsEachLocationOnce,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1493\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,359\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,65754,1493\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestExplicitPlanReadsEachLocationOnce,\"journey_test.go:279: one explicit plan started the common-directory query 4 times, want 3\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t05-plan-run-value",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "read the common directory again in planExplicitWith",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:lgp-t05-author2-20261021",
              "digest": "sha256:3fd74654de1bfb9ec4676265fbebcc585d2a0ae08fb49aa4fe0733f8af415729",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/subshell.go,swap,failed,1,TestExplicitPlanReadsEachLocationOnce,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1493\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,359\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,65754,1493\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestExplicitPlanReadsEachLocationOnce,\"journey_test.go:279: one explicit plan started the common-directory query 4 times, want 3\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c2c-standards-r0",
          "performer": "claude:lgp-c2c-standards-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1bc5392d81da0bc27398f21fc9886a61f8bc59bf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2c-standards-20261021",
            "digest": "sha256:06a23879df9b063129b2ece3e25d8ce4cda3f742d8947fca77f43bc926014848",
            "excerpt": "Standards axis, LGP-C2C (bb8a207e3..cd9889779): 0 findings. The thin validateOwnerMarker form has only test callers now (journey_test.go, admin_readers_test.go, eligibility_test.go), but ticket 04 keeps it as the graded LGP17 seam, and it delegates in one line, so no fact has two sources: no-op. The five &git.Run{} call-line constructions follow ticket 05's direction and ticket 04's constructor shape. The three setup lines that TestExplicitPlanReadsEachLocationOnce repeats are incidental helper calls. Every touched comment matches the code, and classifier.go and clean_landed.go gain no line. Advice: subshell.go:200 still reads git.AdminDir(target) directly while the nested validation reads run.AdminDir; ticket 05 leaves it out of scope.\n"
          },
          "axis": "Standards",
          "base": "bb8a207e3b56c973ee1cd6fe043e4bc3167e6351",
          "tip": "cd9889779c97c9d50b07158bce399e6b40612737",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c2c-spec-r0",
          "performer": "claude:lgp-c2c-spec-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1bc5392d81da0bc27398f21fc9886a61f8bc59bf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2c-spec-20261021",
            "digest": "sha256:c06321c9a877a72852da55f79b4ca70d35dd94de1a3b70942d009694a7e7b0ac",
            "excerpt": "Spec axis, LGP-C2C (bb8a207e3..cd9889779): 0 findings. LGP16 met: planExplicitWith (subshell.go:160) reads canonicalRootIn, run.Worktrees, run.CommonDir, and validateOwnerMarkerIn through the run value; the intent.Address ledger reads and the default-branch reads stay live; worktree.go:347 and resume.go:86 pass a.run, and the five other callers pass &git.Run{}; TestExplicitPlanReadsEachLocationOnce asserts 3. LGP21 and LGP22 met: 806 + 1 = 807 tests, 46 - 0 + 1 = 47 serial. LGP34: coordinator records the ticket 05 author's bench test --full --package ./internal/worktree times, 67314 ms (1492 tests) before and 66960 ms (1493 tests) after, inside the spread. Plan commit 83bfa89b7 (red prediction eight to six) is accurate and non-behavioral. Advice: planExplicitSet could pass a.run to its rows; spec line 135 could adopt ticket 05's wording for callers without an ambient.\n"
          },
          "axis": "Spec",
          "base": "bb8a207e3b56c973ee1cd6fe043e4bc3167e6351",
          "tip": "cd9889779c97c9d50b07158bce399e6b40612737",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c2c-coverage-r0",
          "performer": "claude:lgp-c2c-coverage-20261021",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1bc5392d81da0bc27398f21fc9886a61f8bc59bf",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c2c-coverage-20261021",
            "digest": "sha256:5ae48f71f3a921d3137fad05a6c6d8d60620bc8a8db570173cfcc6c92020cf4c",
            "excerpt": "Coverage axis, LGP-C2C (bb8a207e3..cd9889779): 0 findings. Probes, all restored yes: the nested validateOwnerMarkerIn swapped for the fresh-run thin form at subshell.go:225 bit (4 starts, want 3); canonicalRootIn swapped for canonicalRoot at subshell.go:161 bit (4 starts); resume.go:86 a.run swapped for &git.Run{} read silent across the whole package, an efficiency-only mutant that no acceptance row covers. Refuted: a stale shared a.run answer (the run caches only location directories, revalidates the common directory, and the fingerprint reads a live AdminDir); a root spelling cache miss (spec accepts one extra start); the single-read census. Advice: a verb-level start count could pin the ambient wiring; one plan reads the admin directory twice (subshell.go:200 and owner_marker.go:54), within ticket scope.\n"
          },
          "axis": "Coverage",
          "base": "bb8a207e3b56c973ee1cd6fe043e4bc3167e6351",
          "tip": "cd9889779c97c9d50b07158bce399e6b40612737",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "LGP-C3A",
      "base": "cd9889779c97c9d50b07158bce399e6b40612737",
      "tip": "14d72a647f88c503680baf8e6d4253f587fdceba",
      "plan_digest": "sha256:8c97d52b575a5fb918a786e9098901522c03d40ce6f294f8d2eac848e89eb4d6",
      "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
      "acceptance_rows": [
        "LGP18",
        "LGP22",
        "LGP34"
      ],
      "verification": [
        {
          "id": "t06-pkg-cmd-bench",
          "performer": "claude:lgp-t06-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0bca031bcf55878ac2b0a1202fda122b4e3e33ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-author-20261022",
            "digest": "sha256:bc96a75b1919dd5e23324d25d9f6f00d4c963fe3335123164f58cd35382f6352",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16125,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t06-worktree",
          "performer": "claude:lgp-t06-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0bca031bcf55878ac2b0a1202fda122b4e3e33ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-author-20261022",
            "digest": "sha256:03a3c40e9befd13a6fae5cc962b89fa81914eee4f2125539b8f2b253b34b9fde",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,63972,1495\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t06-canonicalpath",
          "performer": "claude:lgp-t06-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0bca031bcf55878ac2b0a1202fda122b4e3e33ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-author-20261022",
            "digest": "sha256:0e215200df57dd24d2c7ac061bd57bc68af495eec0b2a2af6162c3fa80075653",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/canonicalpath,pass,2,12\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-canonicalpath",
          "command": "bench test --package ./internal/canonicalpath",
          "exit_code": 0
        },
        {
          "id": "t06-root-conformance",
          "performer": "claude:lgp-t06-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0bca031bcf55878ac2b0a1202fda122b4e3e33ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-author-20261022",
            "digest": "sha256:7d404c7716871dce52a97c9a460a9b6c9df0e29fe000f37691df6ec1f93d3356",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12966,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06-call-directory",
          "performer": "claude:lgp-t06-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0bca031bcf55878ac2b0a1202fda122b4e3e33ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-author-20261022",
            "digest": "sha256:a34096f62f91421804b521dbb6d1ea96919b4e8c15bd945e2fb705f2abeceb87",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/subshell.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,65439,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:65: clean \\\"./fb9f22b1767d97ffe4b107bade374bec-f6904e3db200e0e8a830bc70c736fc83\\\" = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktrees/bench-2826441890… (735 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-call-directory",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "resolve a relative operand against the process directory",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t06-author-20261022",
              "digest": "sha256:a34096f62f91421804b521dbb6d1ea96919b4e8c15bd945e2fb705f2abeceb87",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/subshell.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,65439,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:65: clean \\\"./fb9f22b1767d97ffe4b107bade374bec-f6904e3db200e0e8a830bc70c736fc83\\\" = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktrees/bench-2826441890… (735 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06-r1-t06-pkg-cmd-bench",
          "performer": "claude:lgp-t06-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ffc425ae18affb5db9e65ee556d33ffa92016da0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair1-20261022",
            "digest": "sha256:f2ab88fef3cffece4ccf029213cbd538f5b868292bc2281b45d9e84ebdbbba2e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15679,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t06-r1-t06-worktree",
          "performer": "claude:lgp-t06-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ffc425ae18affb5db9e65ee556d33ffa92016da0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair1-20261022",
            "digest": "sha256:5c035c3676f0d5d0e091b219efa96b8834eda65d7c43b39a749a8926387ee0e7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,64976,1495\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t06-r1-t06-canonicalpath",
          "performer": "claude:lgp-t06-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ffc425ae18affb5db9e65ee556d33ffa92016da0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair1-20261022",
            "digest": "sha256:0e215200df57dd24d2c7ac061bd57bc68af495eec0b2a2af6162c3fa80075653",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/canonicalpath,pass,2,12\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-canonicalpath",
          "command": "bench test --package ./internal/canonicalpath",
          "exit_code": 0
        },
        {
          "id": "t06-r1-t06-root-conformance",
          "performer": "claude:lgp-t06-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ffc425ae18affb5db9e65ee556d33ffa92016da0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair1-20261022",
            "digest": "sha256:ad513d08b16fdc330c9d529ac4cbddf9101a3dd5031347e6cb88ae333af4f344",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12893,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06-r1-t06-call-directory",
          "performer": "claude:lgp-t06-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ffc425ae18affb5db9e65ee556d33ffa92016da0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair1-20261022",
            "digest": "sha256:40b09c802041841ee2b652548e852c70791c8bff163ad788fb2ea37bdc73be8f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/subshell.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,65118,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:66: clean \\\"./24da2fe974cae11318a60ff690fbf9e4-1a77b2ea472a099952dfcd1242ae85fd\\\" = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktrees/bench-2826441890… (734 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-call-directory",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "resolve a relative operand against the process directory",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t06-repair1-20261022",
              "digest": "sha256:40b09c802041841ee2b652548e852c70791c8bff163ad788fb2ea37bdc73be8f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/subshell.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,65118,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:66: clean \\\"./24da2fe974cae11318a60ff690fbf9e4-1a77b2ea472a099952dfcd1242ae85fd\\\" = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktrees/bench-2826441890… (734 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06-r1-t06-routed-sites",
          "performer": "claude:lgp-t06-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ffc425ae18affb5db9e65ee556d33ffa92016da0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair1-20261022",
            "digest": "sha256:69e4bb8f83dc02dca7dd1bcb6493c7017f4d6936e180b30ba014b9bab62001c8",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/land.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,534\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,65112,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:92: land \\\"./3f90e89ff54ba15e5908cb25fd007ef8-0f3c7ada52b05855faad71b2f4049a7c\\\" = 1 \\\"refused{detail=request token matches no assignment,next=later proofs in this group did not run; bench worktree list}\\\\\\\\n\\\" \\\"\\\", want … (421 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-routed-sites",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "resolve the land, land resume, and reauthorize operands against the process directory; the LGP18 test must red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t06-repair1-20261022",
              "digest": "sha256:69e4bb8f83dc02dca7dd1bcb6493c7017f4d6936e180b30ba014b9bab62001c8",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/land.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,534\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,65112,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:92: land \\\"./3f90e89ff54ba15e5908cb25fd007ef8-0f3c7ada52b05855faad71b2f4049a7c\\\" = 1 \\\"refused{detail=request token matches no assignment,next=later proofs in this group did not run; bench worktree list}\\\\\\\\n\\\" \\\"\\\", want … (421 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06-r1-t06-clean-failed-read",
          "performer": "claude:lgp-t06-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ffc425ae18affb5db9e65ee556d33ffa92016da0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair1-20261022",
            "digest": "sha256:e68bf1a033f7c84e8b0cba48211b63bddb94c9c530aad88ba55526b43de18ac1",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/worktree.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,64927,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:69: clean \\\"./09c537a47719ecaaf3656b72fc29ae84-4945877479d0f429ae517103bbb3ffbe\\\" without a directory = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktr… (600 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-clean-failed-read",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "drop the carried directory-read failure on the clean path; the clean failed-read test must red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t06-repair1-20261022",
              "digest": "sha256:e68bf1a033f7c84e8b0cba48211b63bddb94c9c530aad88ba55526b43de18ac1",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/worktree.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,64927,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:69: clean \\\"./09c537a47719ecaaf3656b72fc29ae84-4945877479d0f429ae517103bbb3ffbe\\\" without a directory = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktr… (600 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06-r2-t06-pkg-cmd-bench",
          "performer": "claude:lgp-t06-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
            "digest": "sha256:46a710075188b975b07085a7cdee100a8c6404a822d8af4d290697aea6cf8c51",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14385,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t06-r2-t06-worktree",
          "performer": "claude:lgp-t06-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
            "digest": "sha256:6ca84999b5353ef86302520a7ab281800b5e532f888fdd4dd008040fbdb8b38d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,65342,1495\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t06-r2-t06-canonicalpath",
          "performer": "claude:lgp-t06-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
            "digest": "sha256:0e215200df57dd24d2c7ac061bd57bc68af495eec0b2a2af6162c3fa80075653",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/canonicalpath,pass,2,12\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-canonicalpath",
          "command": "bench test --package ./internal/canonicalpath",
          "exit_code": 0
        },
        {
          "id": "t06-r2-t06-root-conformance",
          "performer": "claude:lgp-t06-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
            "digest": "sha256:7d404c7716871dce52a97c9a460a9b6c9df0e29fe000f37691df6ec1f93d3356",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12966,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06-r2-t06-call-directory",
          "performer": "claude:lgp-t06-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
            "digest": "sha256:cffd0cbc63c5c7936d9c8eefbfb5971c86c503343c7daf7885b03a7d8b6aaaa2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/subshell.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,64786,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:67: clean \\\"./58b7e6d014664dbe3bd56fdd7c2a2460-dd316dc7cc5e359a9d25b53c33956f81\\\" = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktrees/bench-2826441890… (735 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-call-directory",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "resolve a relative operand against the process directory",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
              "digest": "sha256:cffd0cbc63c5c7936d9c8eefbfb5971c86c503343c7daf7885b03a7d8b6aaaa2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/subshell.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,64786,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:67: clean \\\"./58b7e6d014664dbe3bd56fdd7c2a2460-dd316dc7cc5e359a9d25b53c33956f81\\\" = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktrees/bench-2826441890… (735 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06-r2-t06-routed-sites",
          "performer": "claude:lgp-t06-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
            "digest": "sha256:d7c8331d799eeb25cf3d0cb3a03309c9516e96f1f414f931678e4856b082908c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/land.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,534\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,64509,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:102: land \\\"./ac76ad7076009a48e3fb3505e2ebacd9-f84f0a2cb68a7115d35d873896123fd7\\\" = 1 \\\"refused{detail=request token matches no assignment,next=later proofs in this group did not run; bench worktree list}\\\\\\\\n\\\" \\\"\\\", want… (422 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-routed-sites",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "resolve the land, land resume, and reauthorize operands against the process directory; the LGP18 test must red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
              "digest": "sha256:d7c8331d799eeb25cf3d0cb3a03309c9516e96f1f414f931678e4856b082908c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/land.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,534\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,64509,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:102: land \\\"./ac76ad7076009a48e3fb3505e2ebacd9-f84f0a2cb68a7115d35d873896123fd7\\\" = 1 \\\"refused{detail=request token matches no assignment,next=later proofs in this group did not run; bench worktree list}\\\\\\\\n\\\" \\\"\\\", want… (422 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06-r2-t06-clean-failed-read",
          "performer": "claude:lgp-t06-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
            "digest": "sha256:7b59a51d2dd3cd34a905a69747980bc2da2182320976b6865b438d71b6934e76",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/worktree.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,65008,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:70: clean \\\"./a2df9bbec01acf8c307234c8b2f94dd6-3cec1143a949eb059a87a786392af5c8\\\" without a directory = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktr… (600 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-clean-failed-read",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "drop the carried directory-read failure on the clean path; the clean failed-read test must red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
              "digest": "sha256:7b59a51d2dd3cd34a905a69747980bc2da2182320976b6865b438d71b6934e76",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/worktree.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,65008,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:70: clean \\\"./a2df9bbec01acf8c307234c8b2f94dd6-3cec1143a949eb059a87a786392af5c8\\\" without a directory = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktr… (600 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06-r2-t06-apply-call-directory",
          "performer": "claude:lgp-t06-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
            "digest": "sha256:725108894d653dcd4227915397066100d3dd07164b4827ad35072135178aa79b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/worktree.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,64807,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:76: clean \\\"./37faf5a99377809ecafbc3e92feb79b5-551255bf72d5662cfa1acb68849ac0b8\\\" --apply = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktrees/bench-28… (749 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t06-apply-call-directory",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "gate the clean call-directory join to the plan posture; the clean --apply case of the LGP18 test must red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t06-repair2-20261022",
              "digest": "sha256:725108894d653dcd4227915397066100d3dd07164b4827ad35072135178aa79b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/worktree.go,swap,failed,1,TestVerbCallDirectoryResolvesOperands,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1495\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,531\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,64807,1495\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestVerbCallDirectoryResolvesOperands,\"identifier_operand_test.go:76: clean \\\"./37faf5a99377809ecafbc3e92feb79b5-551255bf72d5662cfa1acb68849ac0b8\\\" --apply = 1 \\\"worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\\\\\\\\n  /home/mgibs/.bench/worktrees/bench-28… (749 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c3a-standards-r0",
          "performer": "claude:lgp-c3a-standards-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0bca031bcf55878ac2b0a1202fda122b4e3e33ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3a-standards-20261022",
            "digest": "sha256:7d0813a4b1a341201f7f9c367f5adbb1f9f0354b61d7aba7623c746abe414c07",
            "excerpt": "Standards axis, LGP-C3A (cd9889779..b31f7de1e): 2 findings. S1 (accepted, repair): under a failed call-directory read, clean passes the relative operand down as typed, and planExplicitWith -> resolveOperand -> canonicalpath.Resolve reads the process directory again below the boundary (worktree.go:279-283, subshell.go:121-124); the callPath comment claims no lower code does that, and effects.go says lower owners never read the current directory. S2 (accepted, auto-fix): the boundary's `~` test (worktree.go:282) spells the home-operand rule a second time beside expandHomeTarget (path.go:277-278). Advice: accept Absolute plus Resolve as the directory form; rename reauthorize's closure parameter to drop the run alias.\n"
          },
          "axis": "Standards",
          "base": "cd9889779c97c9d50b07158bce399e6b40612737",
          "tip": "b31f7de1ef3bafe11ae5909bbfdf1b9f1ad539e5",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c3a-spec-r0",
          "performer": "claude:lgp-c3a-spec-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0bca031bcf55878ac2b0a1202fda122b4e3e33ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3a-spec-20261022",
            "digest": "sha256:32f7bca83a1856b87807545613f76f0ae72b68c0b635592d82930e41e9de45de",
            "excerpt": "Spec axis, LGP-C3A (cd9889779..b31f7de1e): 1 finding. P1 (accepted, repair, same defect as S1): spec line 156, \"A failed read is carried forward and surfaces only when a relative operand needs it\", fails on clean, where callPath's error is dropped (worktree.go:282-284) and the plan reads the process directory again; repair by rendering the existing refusal at the boundary through one named error shared with resolveOperand. LGP18 met (TestVerbCallDirectoryResolvesOperands; t06-call-directory bit). LGP22 met: 807 + 2 = 809 tests, ceiling 47. LGP34: coordinator records the ticket 06 author's bench test --full --package ./internal/worktree times, 65266 ms (1493 tests) before and 65576 ms (1495 tests) after, inside the spread. Absolute operand under a failed read: met. Author decisions 1, 3, 4, and 5 accepted; decision 1 is a non-behavioral form change flagged for veto.\n"
          },
          "axis": "Spec",
          "base": "cd9889779c97c9d50b07158bce399e6b40612737",
          "tip": "b31f7de1ef3bafe11ae5909bbfdf1b9f1ad539e5",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "c3a-coverage-r0",
          "performer": "claude:lgp-c3a-coverage-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0bca031bcf55878ac2b0a1202fda122b4e3e33ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3a-coverage-20261022",
            "digest": "sha256:6339d31f795a6023fb7e643b7a89c1cf1c17f4718cda7d19bbdb34d79fd1700c",
            "excerpt": "Coverage axis, LGP-C3A (cd9889779..b31f7de1e): 2 findings. C1 (accepted, repair): the LGP18 test exercises only clean and release; canonicalPathIn swapped back to canonicalPath read silent at land.go:137, reauthorize.go:70, and land_resume.go:36 (restored yes), so three routed sites have no test. C2 (accepted, folded into the S1/P1 repair): no test pins clean under a failed directory value. Probes that bit (restored yes): the `~` guard forced true; canonicalpath.Absolute swapped for filepath.Join; an added filepath.Abs read in callPath (single-read census); clean's boundary callPath omitted. Advice: the clean error path at worktree.go:407 affects only the displayed Target.\n"
          },
          "axis": "Coverage",
          "base": "cd9889779c97c9d50b07158bce399e6b40612737",
          "tip": "b31f7de1ef3bafe11ae5909bbfdf1b9f1ad539e5",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "c3a-standards-r1",
          "performer": "claude:lgp-c3a-standards-r1-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ffc425ae18affb5db9e65ee556d33ffa92016da0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3a-standards-r1-20261022",
            "digest": "sha256:6f9b1068ec6eba6f88d41b6e39e9ee4ecfabb91a3c8a6c92fcd7f1684207725a",
            "excerpt": "Standards axis r1, LGP-C3A (cd9889779..4f2541aea): S1, S2 closed (homeOperand is the one home-operand predicate; clean refuses at the boundary with errTargetNotCanonical). 2 new findings, both accepted for repair: S3, the callPath doc (subshell.go:121-125) says each caller passes the path no further, but clean passes the unresolved path on as the refusal row's Target; S4, the operand refusal plan unresolvedPlan(x, ReasonUncertain, ...) is built at worktree.go:286 and subshell.go:189, so the reason-code pairing has two sources. Advice: a failed read gives different refusal texts per verb family; only clean's is new.\n"
          },
          "axis": "Standards",
          "base": "cd9889779c97c9d50b07158bce399e6b40612737",
          "tip": "4f2541aea85b0ed2c45ad572997ba8c3be999aee",
          "finding_ids": [
            "S3",
            "S4"
          ],
          "supersedes": [
            "c3a-standards-r0"
          ]
        },
        {
          "id": "c3a-spec-r1",
          "performer": "claude:lgp-c3a-spec-r1-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ffc425ae18affb5db9e65ee556d33ffa92016da0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3a-spec-r1-20261022",
            "digest": "sha256:ef562372580c1e22a588a05d7aa1b471745747c83903b0c013148f56bc7676c7",
            "excerpt": "Spec axis r1, LGP-C3A (cd9889779..4f2541aea): 0 findings; S1, S2, P1, C1, and C2 closed. Spec line 156 now holds: callPath errors only for a relative operand, and clean refuses at the boundary with the plan's own text. The --apply change from an error row to a retain row keeps ticket 06 line 25's refusal text and exit 1; no spec line binds the row action. LGP18 met (t06-call-directory bit); LGP22 met, 807 + 2 = 809, ceiling 47; LGP34: coordinator records the repair tip's bench test --package ./internal/worktree time, 64976 ms (1495 tests), against 65266 ms before and 65576 ms after the ticket, inside the spread. Every acceptance item met.\n"
          },
          "axis": "Spec",
          "base": "cd9889779c97c9d50b07158bce399e6b40612737",
          "tip": "4f2541aea85b0ed2c45ad572997ba8c3be999aee",
          "finding_ids": [],
          "supersedes": [
            "c3a-spec-r0"
          ]
        },
        {
          "id": "c3a-coverage-r1",
          "performer": "claude:lgp-c3a-coverage-r1-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ffc425ae18affb5db9e65ee556d33ffa92016da0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3a-coverage-r1-20261022",
            "digest": "sha256:5255b8419c21dc76dfb6c690e33ca20476e395bf8cfb1cd329ef1147ea3baeb9",
            "excerpt": "Coverage axis r1, LGP-C3A (cd9889779..4f2541aea): C1 and C2 closed; land_resume.go:36 and reauthorize.go:70 canonicalPath swaps bit, and a callPath bypass that drops a.dirErr bit (restored yes). 1 new finding, accepted for repair: C3, gating the clean boundary to the plan posture (`fingerprint == \"\"`) read silent across the whole package, so clean --apply sends a relative operand to planExplicitWith and cleanupIdentity unpinned, against ticket 06 line 23 and spec line 156; the LGP18 test needs two --apply cases (good directory names the creation path; failed read exits 1 with a retain row and the not-canonical text).\n"
          },
          "axis": "Coverage",
          "base": "cd9889779c97c9d50b07158bce399e6b40612737",
          "tip": "4f2541aea85b0ed2c45ad572997ba8c3be999aee",
          "finding_ids": [
            "C3"
          ],
          "supersedes": [
            "c3a-coverage-r0"
          ]
        },
        {
          "id": "c3a-standards-r2",
          "performer": "claude:lgp-c3a-standards-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3a-standards-r2-20261022",
            "digest": "sha256:950a5127c4ad7793c10c420b310c61680c8fc93914560ca6c5eecd8161d8cbc3",
            "excerpt": "Standards axis r2, LGP-C3A (cd9889779..14d72a647): 0 findings; S1 to S4 closed. The callPath doc (subshell.go:121-126) matches its two callers. One operandRefusalPlan constructor (subshell.go:169) serves the plan (subshell.go:197) and the clean boundary (worktree.go:286), and errTargetNotCanonical is the one refusal value that the tests compare by symbol. homeOperand (path.go:273) is the one home-operand predicate. Advice: renaming reauthorize's closure parameter would drop the run alias. A multi-target clean resolves through clean_set.go, where targetPath refuses a relative path (Coverage r0 and r1 confirmed it), so the call-directory join does not apply there.\n"
          },
          "axis": "Standards",
          "base": "cd9889779c97c9d50b07158bce399e6b40612737",
          "tip": "14d72a647f88c503680baf8e6d4253f587fdceba",
          "finding_ids": [],
          "supersedes": [
            "c3a-standards-r1"
          ]
        },
        {
          "id": "c3a-spec-r2",
          "performer": "claude:lgp-c3a-spec-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3a-spec-r2-20261022",
            "digest": "sha256:89720478c28745be9d40cc1f7c81efd1c1a5d9299f42846f8181448bbbed09b3",
            "excerpt": "Spec axis r2, LGP-C3A (cd9889779..14d72a647): 0 findings; S1 to S4, P1, and C1 to C3 closed. Spec line 156 holds on clean, release, land, land resume, and reauthorize; operandRefusalPlan gives every operand refusal one reason code and the errTargetNotCanonical text, as ticket 06 line 25 requires. LGP18 met: t06-call-directory, t06-routed-sites, t06-clean-failed-read, and t06-apply-call-directory all bit at the r2 tip. LGP22 met: 807 + 2 = 809, ceiling 47. LGP34: coordinator records the second repair tip's bench test --package ./internal/worktree time, 65342 ms (1495 tests), against the ticket's 65266 ms before and 65576 ms after, inside the spread. Every acceptance item met.\n"
          },
          "axis": "Spec",
          "base": "cd9889779c97c9d50b07158bce399e6b40612737",
          "tip": "14d72a647f88c503680baf8e6d4253f587fdceba",
          "finding_ids": [],
          "supersedes": [
            "c3a-spec-r1"
          ]
        },
        {
          "id": "c3a-coverage-r2",
          "performer": "claude:lgp-c3a-coverage-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "66dedb119854ada4ff0ac6b1d1c945fa48b9bc60",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3a-coverage-r2-20261022",
            "digest": "sha256:5d738473083a4bb083dc7e3740b90075c89383a58cbdbd19b888442280a6675d",
            "excerpt": "Coverage axis r2, LGP-C3A (cd9889779..14d72a647): 0 findings; S3, S4, and C3 closed, and the r0 findings stay closed. Probes, restored yes: the boundary refusal dropped under --apply only bit (identifier_operand_test.go:79, the failed-read --apply case); the boundary retain row turned into an error row under --apply only bit (:79); an inline unresolvedPlan with a divergent reason code at subshell.go:197 read silent across the whole package, an equivalent mutant because no output renders ReasonCode. Refuted: multi-target clean (targetPath refuses a relative path), cleanupIdentity callers (canonical paths), the canonicalPathIn fallback (only a ~ operand), spaces and glob characters (text join). Advice: one full-package baseline run failed TestLandCommandCarriesTheBaselineScheduleRootIntoTheProspectiveGate with an infrastructure refusal, and it passed alone and on the next run.\n"
          },
          "axis": "Coverage",
          "base": "cd9889779c97c9d50b07158bce399e6b40612737",
          "tip": "14d72a647f88c503680baf8e6d4253f587fdceba",
          "finding_ids": [],
          "supersedes": [
            "c3a-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "LGP-C3B",
      "base": "14d72a647f88c503680baf8e6d4253f587fdceba",
      "tip": "b2f333d0b73d6058a4f881d15649de14854d698f",
      "plan_digest": "sha256:672b86bd3295ae362b43842629b00cdbd2cbbd02d14bbded0388a400f2fa24b0",
      "source_digest": "7f455bc4e131f44b0458e65bda0c4d801b68846c",
      "acceptance_rows": [
        "LGP20",
        "LGP21",
        "LGP22",
        "LGP23",
        "LGP34"
      ],
      "verification": [
        {
          "id": "t07-pkg-cmd-bench",
          "performer": "claude:lgp-t07-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f455bc4e131f44b0458e65bda0c4d801b68846c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t07-author-20261022",
            "digest": "sha256:f2274b426fb24e355358152f3b8c440a484ea486bfcc647c5818e8cf300c1fcf",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14732,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t07-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t07-worktree",
          "performer": "claude:lgp-t07-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f455bc4e131f44b0458e65bda0c4d801b68846c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t07-author-20261022",
            "digest": "sha256:badf6487cb13219e77df2d0f1c240550497dd074bb9cc691aa6232d91ca197c7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,63074,1496\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t07-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t07-root-conformance",
          "performer": "claude:lgp-t07-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f455bc4e131f44b0458e65bda0c4d801b68846c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t07-author-20261022",
            "digest": "sha256:705f36ba10d270f4943234364bbe670da21c4d4ef51752f4fe23fb30a8748c61",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12995,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t07-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t07-ceiling",
          "performer": "claude:lgp-t07-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f455bc4e131f44b0458e65bda0c4d801b68846c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t07-author-20261022",
            "digest": "sha256:d1c074f30abc07b6078967cea59282d0b72242f0fe8528bebd902dc74c279266",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/parallel_census_test.go,swap,failed,1,TestSerialSetStaysBelowTheCeiling,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1496\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,404\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,63068,1496\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestSerialSetStaysBelowTheCeiling,\"parallel_census_test.go:1052: the package holds 31 serial tests, below the ceiling of 47: lower worktreeSerialCeiling to 31 in this change\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t07-ceiling",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "leave the ceiling at its ticket-06 value",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t07-author-20261022",
              "digest": "sha256:d1c074f30abc07b6078967cea59282d0b72242f0fe8528bebd902dc74c279266",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/parallel_census_test.go,swap,failed,1,TestSerialSetStaysBelowTheCeiling,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1496\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,404\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,63068,1496\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestSerialSetStaysBelowTheCeiling,\"parallel_census_test.go:1052: the package holds 31 serial tests, below the ceiling of 47: lower worktreeSerialCeiling to 31 in this change\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c3b-standards-r0",
          "performer": "claude:lgp-c3b-standards-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7f455bc4e131f44b0458e65bda0c4d801b68846c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3b-standards-20261022",
            "digest": "sha256:37ec89f330dae17f47211791c5fc29df4c0d24f67488e79f7fe9b1da3a8d21a1",
            "excerpt": "Standards axis, LGP-C3B (14d72a647..b2f333d0b): 0 findings. Every former chdir site that calls a joins-form verb uses repoHome.callFromRoot, and only land_resume_refusal_test.go:346 keeps chdir. TestRemovedDirectoryOperandStaysSerial reuses censusFacts, so no second census walk exists, and its \"chdir\" literal names the serial helper that reachedEffects reports. The four-class comment in parallel_census_test.go lists the four current classes. The delta touches only internal/worktree test files. Advice: resume_reconcile_test.go:354-355 sets joins on a callFromRoot call by hand, as mergeSet.merge does.\n"
          },
          "axis": "Standards",
          "base": "14d72a647f88c503680baf8e6d4253f587fdceba",
          "tip": "b2f333d0b73d6058a4f881d15649de14854d698f",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c3b-spec-r0",
          "performer": "claude:lgp-c3b-spec-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7f455bc4e131f44b0458e65bda0c4d801b68846c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3b-spec-20261022",
            "digest": "sha256:2d2a6ac2a193d75855688ecd7e3e834a22123b8d6c4605d37df94afc5efbf3b5",
            "excerpt": "Spec axis, LGP-C3B (14d72a647..b2f333d0b): 0 blocking findings. LGP20 met: the 16 moved tests call t.Parallel and drop chdir, and the live census passes. LGP21 met: 45 + 1 + 1 - 16 + 0 = 31 serial. LGP22 met: 805 + 1 + 1 + 2 + 1 = 810. LGP23 met: TestRemovedDirectoryOperandStaysSerial pins the chdir serial reason, and TestLandingHandsBackAPathWithNoCanonicalForm still grades the real filepath.Abs failure (land_resume_refusal_test.go:354). Ticket line 20 holds: path, list, show, and exec build no ambient and targetPath reads no directory for an identifier or absolute operand. P1 and P2 (evidence placement) close here: the coordinator records the ticket 07 author's bench test --full --package ./internal/worktree times, 65067 ms (1495 tests) before and 62560 ms (1496 tests) after, inside the spread, and the moved and stayed lists with reasons and the no-directory proof are in ticket commit b2f333d0b's body.\n"
          },
          "axis": "Spec",
          "base": "14d72a647f88c503680baf8e6d4253f587fdceba",
          "tip": "b2f333d0b73d6058a4f881d15649de14854d698f",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c3b-coverage-r0",
          "performer": "claude:lgp-c3b-coverage-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7f455bc4e131f44b0458e65bda0c4d801b68846c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3b-coverage-20261022",
            "digest": "sha256:e1828406f39cd1e73f2b62fedadf66f45bea2f4cbeecdac2422a6cf773104806",
            "excerpt": "Coverage axis, LGP-C3B (14d72a647..b2f333d0b): 0 binding findings. Probes, all restored yes: callFromRoot returning a failed directory read read silent across the 11 joins-form tests, because each operand is absolute or an identifier, so the helper's effect is the joins-form route and not a directory read (advice, no repair: parallel safety is the graded subject and holds); removing mustRemove(t, gone) while chdir stays bit TestLandingHandsBackAPathWithNoCanonicalForm, so that serial test grades the real filepath.Abs failure itself; omitting one t.Parallel in request_token_test.go bit TestParallelCensusOnTheLiveTree. Every git child in internal/worktree and internal/git on a moved test's path carries -C or Dir. Advice: the callFromRoot comment could say the test does not read the directory value.\n"
          },
          "axis": "Coverage",
          "base": "14d72a647f88c503680baf8e6d4253f587fdceba",
          "tip": "b2f333d0b73d6058a4f881d15649de14854d698f",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "LGP-C3C",
      "base": "b2f333d0b73d6058a4f881d15649de14854d698f",
      "tip": "8ac213e021c3e86955957c6c92073dcf7bf968b8",
      "plan_digest": "sha256:d2657bc9aba2dcd702db7ef57791f341c08425eec7593e80663637f9fda45ee9",
      "source_digest": "450319442053d40f7641de2e84dfa03330ef5c26",
      "acceptance_rows": [
        "LGP19",
        "LGP22",
        "LGP24",
        "LGP34"
      ],
      "verification": [
        {
          "id": "t08-pkg-cmd-bench",
          "performer": "claude:lgp-t08-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "af5ffed93be6a7c13b2d61d29823624d986a8bb4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t08-author2-20261022",
            "digest": "sha256:a76d58400a51fb547cd4a74b788ac0d2e6c345ef83bee5f02491a198f31aba86",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14614,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t08-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t08-worktree",
          "performer": "claude:lgp-t08-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "af5ffed93be6a7c13b2d61d29823624d986a8bb4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t08-author2-20261022",
            "digest": "sha256:b10a474b69eb007dbfb5d44172e07935d66a414db7b22542d79b00267ea5779e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,63227,1498\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t08-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t08-root-conformance",
          "performer": "claude:lgp-t08-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "af5ffed93be6a7c13b2d61d29823624d986a8bb4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t08-author2-20261022",
            "digest": "sha256:ee3ed1465ce50861f2875926b3116869a1831fc8b497d928a7242afc36b00c9b",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12885,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t08-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t08-environ-owner",
          "performer": "claude:lgp-t08-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "af5ffed93be6a7c13b2d61d29823624d986a8bb4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t08-author2-20261022",
            "digest": "sha256:8e150dd15bbda61e66980af90e3b7c01ee7c213a29701999fbf1ad83bee69c22",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/exec.go,swap+swap,failed,3,TestEffectBoundaryCensus,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1498\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,716\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,62731,1498\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestEffectBoundaryCensus,\"effect_census_test.go:45: exec.go:201: ambient read os.Environ( below the effect boundary; resolve it in effects.go and pass the value down\",1\n  github.com/gibbonmi/bench/internal/worktree,TestExecChildTakesTheCallEnvironment,\"exec_call_env_test.go:30: the child environment lacks the call value:\",78\n  github.com/gibbonmi/bench/internal/worktree,TestSingleReadCensusOnTheLiveTree,\"single_read_census_cases_test.go:200: the single-read census reports 1 reads:\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t08-environ-owner",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "restore the os.Environ read in execEnv",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t08-author2-20261022",
              "digest": "sha256:8e150dd15bbda61e66980af90e3b7c01ee7c213a29701999fbf1ad83bee69c22",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/exec.go,swap+swap,failed,3,TestEffectBoundaryCensus,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1498\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,716\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,62731,1498\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestEffectBoundaryCensus,\"effect_census_test.go:45: exec.go:201: ambient read os.Environ( below the effect boundary; resolve it in effects.go and pass the value down\",1\n  github.com/gibbonmi/bench/internal/worktree,TestExecChildTakesTheCallEnvironment,\"exec_call_env_test.go:30: the child environment lacks the call value:\",78\n  github.com/gibbonmi/bench/internal/worktree,TestSingleReadCensusOnTheLiveTree,\"single_read_census_cases_test.go:200: the single-read census reports 1 reads:\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t08-r1-t08-pkg-cmd-bench",
          "performer": "claude:lgp-t08-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "450319442053d40f7641de2e84dfa03330ef5c26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t08-repair1-20261022",
            "digest": "sha256:59eb6871bb9bb24bc2ff292fe69b78622119e46e3c20fcb1b91e636afdccc0d1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14993,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t08-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t08-r1-t08-worktree",
          "performer": "claude:lgp-t08-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "450319442053d40f7641de2e84dfa03330ef5c26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t08-repair1-20261022",
            "digest": "sha256:7ad985bd240388535a641418f129a36ed5abc54e07433df1907695d0d5c76549",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,63270,1498\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t08-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t08-r1-t08-root-conformance",
          "performer": "claude:lgp-t08-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "450319442053d40f7641de2e84dfa03330ef5c26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t08-repair1-20261022",
            "digest": "sha256:b95d5216eeecfa98db39eff94b51b6b4d558bbf98b262743421fb688524f2bba",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12629,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t08-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t08-r1-t08-environ-owner",
          "performer": "claude:lgp-t08-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "450319442053d40f7641de2e84dfa03330ef5c26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t08-repair1-20261022",
            "digest": "sha256:c0a1025a19ab8875e328a7c6e0b3b7d8c5dd2fc1b1af9cbf429ff181fc248b25",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/exec.go,swap+swap,failed,3,TestEffectBoundaryCensus,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1498\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,716\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,62958,1498\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestEffectBoundaryCensus,\"effect_census_test.go:45: exec.go:202: ambient read os.Environ( below the effect boundary; resolve it in effects.go and pass the value down\",1\n  github.com/gibbonmi/bench/internal/worktree,TestExecChildTakesTheCallEnvironment,\"exec_call_env_test.go:30: the child environment lacks the call value:\",78\n  github.com/gibbonmi/bench/internal/worktree,TestSingleReadCensusOnTheLiveTree,\"single_read_census_cases_test.go:200: the single-read census reports 1 reads:\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t08-environ-owner",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "restore the os.Environ read in execEnv",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t08-repair1-20261022",
              "digest": "sha256:c0a1025a19ab8875e328a7c6e0b3b7d8c5dd2fc1b1af9cbf429ff181fc248b25",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/exec.go,swap+swap,failed,3,TestEffectBoundaryCensus,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1498\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,716\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,62958,1498\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestEffectBoundaryCensus,\"effect_census_test.go:45: exec.go:202: ambient read os.Environ( below the effect boundary; resolve it in effects.go and pass the value down\",1\n  github.com/gibbonmi/bench/internal/worktree,TestExecChildTakesTheCallEnvironment,\"exec_call_env_test.go:30: the child environment lacks the call value:\",78\n  github.com/gibbonmi/bench/internal/worktree,TestSingleReadCensusOnTheLiveTree,\"single_read_census_cases_test.go:200: the single-read census reports 1 reads:\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t08-r1-t08-recovery-environ",
          "performer": "claude:lgp-t08-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "450319442053d40f7641de2e84dfa03330ef5c26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t08-repair1-20261022",
            "digest": "sha256:41309ff78137cf4849d809adb15a8afb936ef983f283c2e683e980c13feac172",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/clean.go,swap,failed,1,TestResetApplyPreservesTheLayersAndMovesTheCheckout,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1498\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,385\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,63046,1498\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestResetApplyPreservesTheLayersAndMovesTheCheckout,\"reset_apply_test.go:212: envelope committer date = 1792653168, want the call environment's 1000000000\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t08-recovery-environ",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "drop the threaded environ in gitInput; the recovery-chain environment test must red",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t08-repair1-20261022",
              "digest": "sha256:41309ff78137cf4849d809adb15a8afb936ef983f283c2e683e980c13feac172",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/clean.go,swap,failed,1,TestResetApplyPreservesTheLayersAndMovesTheCheckout,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,all,passed,1498\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,385\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,63046,1498\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestResetApplyPreservesTheLayersAndMovesTheCheckout,\"reset_apply_test.go:212: envelope committer date = 1792653168, want the call environment's 1000000000\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c3c-standards-r0",
          "performer": "claude:lgp-c3c-standards-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "af5ffed93be6a7c13b2d61d29823624d986a8bb4",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3c-standards-20261022",
            "digest": "sha256:97388da1495f8bc8abe027397f84c15af8ad0934deb35325bacdfc55131776e7",
            "excerpt": "Standards axis, LGP-C3C (b2f333d0b..8c3e7f18f): 2 findings. S1 (closed in this round, no code change): TestVerbCallRefusesAnEnvironmentValue checks an independently written refusal string with no recorded red; the Coverage axis demonstrated that red (probe P4, the checkVerbCall environ refusal disabled, bit, restored yes), which meets AGENTS.md's record-and-demonstrate condition. S2 (accepted, repair): the ambient doc (effects.go:40-41) says environ is the base environment of each child the verb starts, and the verbCall doc (verb_runner_test.go:49) repeats it, but git.Output children in captureLayers (layers.go:373) inherit the process environment; narrow the doc to the children that take the value and drop the restatement. Advice: ExecCommandResolving builds a whole newAmbient for .environ while RunTreeChild calls baseEnviron.\n"
          },
          "axis": "Standards",
          "base": "b2f333d0b73d6058a4f881d15649de14854d698f",
          "tip": "8c3e7f18f0a9f2c275de73170032f6a015aa6f74",
          "finding_ids": [
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c3c-spec-r0",
          "performer": "claude:lgp-c3c-spec-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "af5ffed93be6a7c13b2d61d29823624d986a8bb4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3c-spec-20261022",
            "digest": "sha256:36995cd56e3de31304dfcfa0980f1c3cb64a6b00a41a2229c48fa93692a1a35f",
            "excerpt": "Spec axis, LGP-C3C (b2f333d0b..8c3e7f18f): 0 open findings. LGP19 met (TestExecChildTakesTheCallEnvironment; t08-environ-owner bit). LGP22 met: 810 + 2 = 812, ceiling 31. LGP24 met: effect_census_test.go adds os.Environ, and the only production os.Environ read is effects.go baseEnviron. The four sites and both chains match the amended ticket; subshellAt and RunTreeChild keep their signatures; exec_test.go has no net line gain; verb_runner_check_test.go has 394 lines. P1 (record placement) closes here: the coordinator records the ticket 08 author's bench test --full --package ./internal/worktree times from commit 8c3e7f18f's body, 63309 ms (1496 tests) before at 2e8dd2851 and 63324 ms (1498 tests) after, inside the spread. Advice: installAttributed puts environ after target.\n"
          },
          "axis": "Spec",
          "base": "b2f333d0b73d6058a4f881d15649de14854d698f",
          "tip": "8c3e7f18f0a9f2c275de73170032f6a015aa6f74",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c3c-coverage-r0",
          "performer": "claude:lgp-c3c-coverage-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "af5ffed93be6a7c13b2d61d29823624d986a8bb4",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3c-coverage-20261022",
            "digest": "sha256:a07e9ba324c129ac0b9bf87f7db193c84a245c6cc693090bb6be339c0f7a3f92",
            "excerpt": "Coverage axis, LGP-C3C (b2f333d0b..8c3e7f18f): 2 findings. C1 (accepted as the chunk's hardening cycle): gitInput dropping the threaded environ (slices.Concat(extraEnv)) read silent across the whole package (1498 tests, restored yes), so no test pins the reset or cleanup recovery chain to the call's base environment. C2 (closed here): probe P4, the checkVerbCall environ refusal disabled, bit TestVerbCallRefusesAnEnvironmentValue, restored yes, which records that expectation's red. Other probes, restored yes: nil environ to repairPrimary on restore and on install bit; repairPrimary reading baseEnviron bit the single-read census; baseEnviron returning nil bit; os.Environ in gitInput bit both censuses; syscall.Environ in subshellAt and a method-value os.Environ read bit only through a behavior test and the single-read census. Advice: the effect census regex misses syscall.Environ and a method value; about 20 git children with nil Env lie outside the spec's four sites.\n"
          },
          "axis": "Coverage",
          "base": "b2f333d0b73d6058a4f881d15649de14854d698f",
          "tip": "8c3e7f18f0a9f2c275de73170032f6a015aa6f74",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "c3c-standards-r1",
          "performer": "claude:lgp-c3c-standards-r1-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "450319442053d40f7641de2e84dfa03330ef5c26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3c-standards-r1-20261022",
            "digest": "sha256:09ea69600df248dd25c1045cfc2cc806a6ba4a6d3b81816ae4acd2a3db92e3cc",
            "excerpt": "Standards axis r1, LGP-C3C (b2f333d0b..8ac213e02): S2 closed; the ambient doc names exactly the children that take the base environment, the verbCall doc points to it, and the baseEnviron doc names its one chain. 1 raised item, dispositioned no-op: subshellAt takes both an environ parameter and the ambient, but ticket 08 line 25 keeps that parameter for its test callers, and finding-discipline makes a finding against a ticket's explicit keep decision a no-op. Advice: ExecCommandResolving builds a whole newAmbient for .environ; the reset_apply_test.go comment on commitTree's date restates the code.\n"
          },
          "axis": "Standards",
          "base": "b2f333d0b73d6058a4f881d15649de14854d698f",
          "tip": "8ac213e021c3e86955957c6c92073dcf7bf968b8",
          "finding_ids": [],
          "supersedes": [
            "c3c-standards-r0"
          ]
        },
        {
          "id": "c3c-spec-r1",
          "performer": "claude:lgp-c3c-spec-r1-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "450319442053d40f7641de2e84dfa03330ef5c26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3c-spec-r1-20261022",
            "digest": "sha256:caba5d1fc3db9e63e5a24d57b111d00009726b0b4decaaf491ba68c2ebdcdc51",
            "excerpt": "Spec axis r1, LGP-C3C (b2f333d0b..8ac213e02): 0 findings; S2 and C1 closed. LGP19 met; LGP24 met (the only production os.Environ read is effects.go baseEnviron); LGP22 met, 810 + 2 = 812, ceiling 31; LGP34: coordinator records 63309 ms (1496 tests) before the ticket and 63270 ms (1498 tests) at the repair tip, inside the spread. The four sites and both chains match the amended ticket; RunTreeChild and subshellAt keep their signatures; exec_test.go 421 lines, verb_runner_check_test.go 394. Advice: ticket 08's \"count rises by one\" undercounts the required refusal test (+2).\n"
          },
          "axis": "Spec",
          "base": "b2f333d0b73d6058a4f881d15649de14854d698f",
          "tip": "8ac213e021c3e86955957c6c92073dcf7bf968b8",
          "finding_ids": [],
          "supersedes": [
            "c3c-spec-r0"
          ]
        },
        {
          "id": "c3c-coverage-r1",
          "performer": "claude:lgp-c3c-coverage-r1-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "450319442053d40f7641de2e84dfa03330ef5c26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3c-coverage-r1-20261022",
            "digest": "sha256:32f74822d09c306b9b775efca6d488119556c74d71759fd877e85ade92969dbd",
            "excerpt": "Coverage axis r1, LGP-C3C (b2f333d0b..8ac213e02): 0 findings; S2 and C1 closed. Probes, restored yes: applyReset receiving nil at reset.go bit TestResetApplyPreservesTheLayersAndMovesTheCheckout (committer date); captureLayers receiving nil in reset_envelope.go bit the same test; recoverAssignmentWithFault receiving nil at lifecycle.go:523 read silent across the package; installChildEnv receiving nil at build.go:43 read silent under the focused build selection. Advice: the cleanup recovery chain and the doctor child have no pin on the threaded environment; no acceptance row requires one, and the chunk's hardening cycle is spent.\n"
          },
          "axis": "Coverage",
          "base": "b2f333d0b73d6058a4f881d15649de14854d698f",
          "tip": "8ac213e021c3e86955957c6c92073dcf7bf968b8",
          "finding_ids": [],
          "supersedes": [
            "c3c-coverage-r0"
          ]
        }
      ]
    },
    {
      "id": "LGP-C3D",
      "base": "8ac213e021c3e86955957c6c92073dcf7bf968b8",
      "tip": "fd7ac338eea9394fa3cf878e9ea495cff9441874",
      "plan_digest": "sha256:b5e476e1aadba9767c08baeae66130b308c5e2b5302df2a0170bf35907fad091",
      "source_digest": "7a9817df854f82559365cabf9cd949f1f8632fc0",
      "acceptance_rows": [
        "LGP20",
        "LGP21",
        "LGP22",
        "LGP34"
      ],
      "verification": [
        {
          "id": "t09-pkg-cmd-bench",
          "performer": "claude:lgp-t09-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a9817df854f82559365cabf9cd949f1f8632fc0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t09-author-20261022",
            "digest": "sha256:b19c1970bee32d338f8f59652bf21a28f9ea6433d0827a8e18e690c536f3f9ac",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14414,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t09-pkg-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t09-worktree",
          "performer": "claude:lgp-t09-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a9817df854f82559365cabf9cd949f1f8632fc0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t09-author-20261022",
            "digest": "sha256:e021ca430642989dee90dc452498f5d0a1e5993d8c0b93085c60cded3c2b1072",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,62647,1498\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t09-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "t09-root-conformance",
          "performer": "claude:lgp-t09-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a9817df854f82559365cabf9cd949f1f8632fc0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t09-author-20261022",
            "digest": "sha256:b5ea876ac0ab4c9b4ad419bb4aaf803b51ab8bb8ebb872336cc1c96286e52d65",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12795,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t09-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t09-exec-bind",
          "performer": "claude:lgp-t09-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a9817df854f82559365cabf9cd949f1f8632fc0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:lgp-t09-author-20261022",
            "digest": "sha256:fb94119dbd296a80a46e4068d094178a694ea993a7ea5430dab0f297577e1f44",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/exec_call_env_test.go,insert,failed,1,TestParallelCensusOnTheLiveTree,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,^TestParallelCensusOnTheLiveTree$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,327\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,90,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestParallelCensusOnTheLiveTree,\"parallel_census_test.go:1029: the parallel census reports 10 breaks:\",11\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t09-exec-bind",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0,
          "probe": {
            "mutation": "bind the caller environment again in childOutputAtHome",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:lgp-t09-author-20261022",
              "digest": "sha256:fb94119dbd296a80a46e4068d094178a694ea993a7ea5430dab0f297577e1f44",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/exec_call_env_test.go,insert,failed,1,TestParallelCensusOnTheLiveTree,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,^TestParallelCensusOnTheLiveTree$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,327\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,90,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestParallelCensusOnTheLiveTree,\"parallel_census_test.go:1029: the parallel census reports 10 breaks:\",11\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c3d-standards-r0",
          "performer": "claude:lgp-c3d-standards-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7a9817df854f82559365cabf9cd949f1f8632fc0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3d-standards-20261022",
            "digest": "sha256:cf81dac1fb9ce20350b84fe7a40409bbc1540087c4adc39a65b34b42e74f7fe3",
            "excerpt": "Standards axis, LGP-C3D (8ac213e02..fd7ac338e): 1 finding. S1 (accepted, auto-fix): the census comment (parallel_census_test.go:1035-1036) says an exec child test takes its caller environment from the call and stays out of the set, which repeats the fact that the helper comment (exec_call_env_test.go:21) owns and names something outside the three serial classes it lists; delete the census sentence. Clean: ten t.Parallel additions match the ten TestExecChild tests; the census comment names three classes; the ceiling is 21; no edit outside the ticket's tests. Advice: the helper chain is split across two files; \"PATH=\" + os.Getenv(\"PATH\") appears in two call bases.\n"
          },
          "axis": "Standards",
          "base": "8ac213e021c3e86955957c6c92073dcf7bf968b8",
          "tip": "fd7ac338eea9394fa3cf878e9ea495cff9441874",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "c3d-spec-r0",
          "performer": "claude:lgp-c3d-spec-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7a9817df854f82559365cabf9cd949f1f8632fc0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3d-spec-20261022",
            "digest": "sha256:69d5e056d841e2ef8f7943f432d9f965ab3bf4e8c087a16688a309f73ab43e8a",
            "excerpt": "Spec axis, LGP-C3D (8ac213e02..fd7ac338e): 1 finding. P1 (accepted, repair): the recorded t09-exec-bind run used --run \"^TestParallelCensusOnTheLiveTree$\" while the plan names bench test --package ./internal/worktree, against the rule \"Run the plan's named probe exactly\"; the plan was not amended, and the record shows the plan command beside a narrowed selection. The Coverage axis then ran the plan's exact probe on the whole package and it bit cleanly, so the repair re-records the plan's run and the plan stays. LGP20 met (ten tests parallel, no bindEnv left, the mutant census reports 10 breaks); LGP21 met, 31 - 10 + 0 = 21; LGP22 met, 812 unchanged; LGP34: coordinator records the ticket 09 author's bench test --full --package ./internal/worktree times, 63189 ms (1498 tests) before and 63078 ms (1498 tests) after, inside the spread. Ticket lines 13, 15, and 17 met (exec_test.go 408 lines).\n"
          },
          "axis": "Spec",
          "base": "8ac213e021c3e86955957c6c92073dcf7bf968b8",
          "tip": "fd7ac338eea9394fa3cf878e9ea495cff9441874",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "c3d-coverage-r0",
          "performer": "claude:lgp-c3d-coverage-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7a9817df854f82559365cabf9cd949f1f8632fc0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:lgp-c3d-coverage-20261022",
            "digest": "sha256:7c15780bbac5f954e6c82b3c5cfce33d9eccf26eb603cfc74e793a2e6d9d496f",
            "excerpt": "Coverage axis, LGP-C3D (8ac213e02..fd7ac338e): 0 findings. Probes, all restored yes: the planned t09-exec-bind mutation (bindEnv in childOutputAtHome) on the whole package as the plan states bit with 3 failures (a t.Setenv panic in a moved test, the census with 10 breaks, the ceiling at 31 against 21) and left no temp directory; execEnv reading baseEnviron instead of the call environ bit two moved tests; removing one t.Parallel bit the census; dropping BENCH_EXEC_TEST_CARRIED from the caller slice bit. Silent and equivalent: withHome stripping a dummy name (exec.Cmd keeps the last duplicate). Advice: dropping BENCH_KIT from the caller slice reads silent, so TestExecChildDropsWrapperRouting could first assert its fixture carries the value; this predates the ticket.\n"
          },
          "axis": "Coverage",
          "base": "8ac213e021c3e86955957c6c92073dcf7bf968b8",
          "tip": "fd7ac338eea9394fa3cf878e9ea495cff9441874",
          "finding_ids": [],
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
      "from": "sha256:400cd4647d7f162a60c2c9c5d31dde7bdfd1f3789c3b6a3c878532c65ddfcdc3",
      "to": "sha256:e15cb1da923ae005902078362e7c2e00b733194f22922820f90fa00d060865c7",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ]
      }
    },
    {
      "from": "sha256:e15cb1da923ae005902078362e7c2e00b733194f22922820f90fa00d060865c7",
      "to": "sha256:6fc699b363a5ca38032b7bd24f0eb9e055a3519c6924419b580011ee584ebdd6",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ]
      }
    },
    {
      "from": "sha256:6fc699b363a5ca38032b7bd24f0eb9e055a3519c6924419b580011ee584ebdd6",
      "to": "sha256:504f3662dc285101a10b12a57d0d3171265e21af7bf269af10240c19c8d2b397",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ],
        "LGP-C2A": [
          "LGP-C2A"
        ]
      }
    },
    {
      "from": "sha256:504f3662dc285101a10b12a57d0d3171265e21af7bf269af10240c19c8d2b397",
      "to": "sha256:f97a2dd92aa8f25d612c082e08c62d69946a1bc987a3ae263cae84b7ac1b56b3",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ],
        "LGP-C2A": [
          "LGP-C2A"
        ]
      }
    },
    {
      "from": "sha256:f97a2dd92aa8f25d612c082e08c62d69946a1bc987a3ae263cae84b7ac1b56b3",
      "to": "sha256:c29cdfbe2ba7223fba8e656138b5210901464cb9f4b555854aabd86698ec4368",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ],
        "LGP-C2A": [
          "LGP-C2A"
        ],
        "LGP-C2B": [
          "LGP-C2B"
        ]
      }
    },
    {
      "from": "sha256:c29cdfbe2ba7223fba8e656138b5210901464cb9f4b555854aabd86698ec4368",
      "to": "sha256:4ac131444b87c7a14de29125ae9b270da620935b6adaed6c5b60471f364eced9",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ],
        "LGP-C2A": [
          "LGP-C2A"
        ],
        "LGP-C2B": [
          "LGP-C2B"
        ],
        "LGP-C2C": [
          "LGP-C2C"
        ]
      }
    },
    {
      "from": "sha256:4ac131444b87c7a14de29125ae9b270da620935b6adaed6c5b60471f364eced9",
      "to": "sha256:91612ff61ecb2fc10a4be3e1df5548c24d844e8ebc268576f6c3adf48a13b10d",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ],
        "LGP-C2A": [
          "LGP-C2A"
        ],
        "LGP-C2B": [
          "LGP-C2B"
        ],
        "LGP-C2C": [
          "LGP-C2C"
        ],
        "LGP-C3A": [
          "LGP-C3A"
        ]
      }
    },
    {
      "from": "sha256:91612ff61ecb2fc10a4be3e1df5548c24d844e8ebc268576f6c3adf48a13b10d",
      "to": "sha256:8c97d52b575a5fb918a786e9098901522c03d40ce6f294f8d2eac848e89eb4d6",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ],
        "LGP-C2A": [
          "LGP-C2A"
        ],
        "LGP-C2B": [
          "LGP-C2B"
        ],
        "LGP-C2C": [
          "LGP-C2C"
        ],
        "LGP-C3A": [
          "LGP-C3A"
        ]
      }
    },
    {
      "from": "sha256:8c97d52b575a5fb918a786e9098901522c03d40ce6f294f8d2eac848e89eb4d6",
      "to": "sha256:672b86bd3295ae362b43842629b00cdbd2cbbd02d14bbded0388a400f2fa24b0",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ],
        "LGP-C2A": [
          "LGP-C2A"
        ],
        "LGP-C2B": [
          "LGP-C2B"
        ],
        "LGP-C2C": [
          "LGP-C2C"
        ],
        "LGP-C3A": [
          "LGP-C3A"
        ]
      }
    },
    {
      "from": "sha256:672b86bd3295ae362b43842629b00cdbd2cbbd02d14bbded0388a400f2fa24b0",
      "to": "sha256:76f8f9e1e0d3dc9d5b45af7c1f854109efe64adf1f66fe0e82cb4229f31b9702",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ],
        "LGP-C2A": [
          "LGP-C2A"
        ],
        "LGP-C2B": [
          "LGP-C2B"
        ],
        "LGP-C2C": [
          "LGP-C2C"
        ],
        "LGP-C3A": [
          "LGP-C3A"
        ],
        "LGP-C3B": [
          "LGP-C3B"
        ]
      }
    },
    {
      "from": "sha256:76f8f9e1e0d3dc9d5b45af7c1f854109efe64adf1f66fe0e82cb4229f31b9702",
      "to": "sha256:d2657bc9aba2dcd702db7ef57791f341c08425eec7593e80663637f9fda45ee9",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ],
        "LGP-C2A": [
          "LGP-C2A"
        ],
        "LGP-C2B": [
          "LGP-C2B"
        ],
        "LGP-C2C": [
          "LGP-C2C"
        ],
        "LGP-C3A": [
          "LGP-C3A"
        ],
        "LGP-C3B": [
          "LGP-C3B"
        ],
        "LGP-C3C": [
          "LGP-C3C"
        ]
      }
    },
    {
      "from": "sha256:d2657bc9aba2dcd702db7ef57791f341c08425eec7593e80663637f9fda45ee9",
      "to": "sha256:b5e476e1aadba9767c08baeae66130b308c5e2b5302df2a0170bf35907fad091",
      "chunk_ids": {
        "LGP-C1": [
          "LGP-C1"
        ],
        "LGP-C2A": [
          "LGP-C2A"
        ],
        "LGP-C2B": [
          "LGP-C2B"
        ],
        "LGP-C2C": [
          "LGP-C2C"
        ],
        "LGP-C3A": [
          "LGP-C3A"
        ],
        "LGP-C3B": [
          "LGP-C3B"
        ],
        "LGP-C3C": [
          "LGP-C3C"
        ]
      }
    }
  ]
}
```

## Pickup: LGP-C1 initial review

The frozen pair is `7f585da44..5cd784915`. Repair cycle 1 of 2 starts here, and C1 uses the one hardening cycle of the chunk.

### Standards

Count: 4 findings, 3 accepted. Worst issue: S1.

- S1 (auto-fix, confidence 8): `ownedRunBinary` (`internal/testreport/selection_once_test.go:12-52`) repeats the worktree `testRunSelector` (`internal/worktree/test_run_test.go:16-90`). It is a third, weaker form of the inherit-or-own rule of `runbinary.Factory.ReuseOrOwn`. Route: ticket 02 repair. The fable tie-breaker chose one `Factory.Select` method and one exported once-holder in `internal/runbinary`.
- S2 (auto-fix, confidence 6): the probe golden repeats the adopt golden harness and drops its mode walk (`internal/probe/golden_test.go:99`, `internal/adopt/golden_test.go:88-104`). Route: ticket 01 repair. One mode-preserving tree copy in `internal/testrepo` serves adopt and probe.
- S3 (auto-fix, confidence 8): the `fixture` doc comment (`internal/probe/probe_test.go:15-16`) still names `probetest.Repository`. Route: ticket 01 repair.

### Spec

Count: 4 findings, 3 accepted. Worst issue: P1.

- P1 (auto-fix, confidence 6): story 8 asks for a real owned build, but the inherited selection wins in the gate. Route: a spec prose amendment of story 8, by consultant consensus.
- P3 (auto-fix, confidence 6): `os.CopyFS` does not keep the modes that spec line 118 promises. Route: closed by the S2 repair.
- P4 (auto-fix, confidence 5): the LGP6 catch clause names a bypass that its test does not see. Route: a spec prose amendment of the LGP6 row, by consultant consensus.

P2 is closed in the record: the `c1-spec-r0` excerpt holds the LGP34 timings.

### Coverage

Count: 1 finding, accepted. Worst issue: C1.

- C1 (auto-fix, confidence 8): no test covers the inherited branch of the testreport selection, so a probe that disables it reads silent. Route: ticket 02 repair, as the one hardening cycle.

### LGP-C1 state

Repair cycle 1 closed every accepted finding at the tip `be0aca0c5`, and the confirming round passed on all three axes. The chunk used 1 of its 2 repair cycles and its one hardening cycle.
The confirming Spec axis raised P5: spec lines 121, 126, and 128 cite lines that the repair shifted. The next plan commit corrects them before the ticket 03 dispatch.

## Pickup: LGP-C2A initial review

The frozen pair is `be0aca0c5..de63c6090`. Repair cycle 1 of 2 starts here, and C1 and C2 use the one hardening cycle of the chunk. Every accepted finding falls on ticket 03, so one combined repair charge carries them.

### Standards

Count: 1 finding, accepted. Worst issue: S1.

- S1 (auto-fix, confidence 4): the stub-git start counter is derived twice (`internal/git/run_cache_test.go:15-29`, `internal/poolkey/run_test.go:17`, `:26-30`). Route: ticket 03 repair; one start counter beside `StubGitDir` in `internal/gittest`.

### Spec

Count: 1 finding, closed in the record. Worst issue: P1.

P1 is closed in the record: the `c2a-spec-r0` excerpt holds the LGP34 timings.

### Coverage

Count: 2 findings, accepted. Worst issue: C1.

- C1 (auto-fix, confidence 7): an admin-directory read keyed as the common directory passes, because LGP10 runs in a primary checkout (`internal/git/run.go:49`). Route: ticket 03 repair, with a linked-worktree test.
- C2 (auto-fix, confidence 6): a cached failed query passes, although `internal/git/run.go:77-78` says a failed start stores nothing. Route: ticket 03 repair, with a failed-query test.

### LGP-C2A state

Repair cycle 1 closed S1, C1, and C2 at the tip `4c4d1df00`, and the confirming round passed on all three axes. The chunk used 1 of its 2 repair cycles and its one hardening cycle.
`gittest.StubGitStarts` is the one start counter for tickets 04 and 05.
