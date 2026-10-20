# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/merge-target-lane/spec.md",
  "plan_digest": "sha256:322c86ff89de964810f82061879d5aa054f5d0e695e4bab97c9f0286471f655c",
  "implementation_session": "",
  "chunks": [
    {
      "id": "C1",
      "base": "d9dd9650cd42c277ae733c522a881b322d50e9e2",
      "tip": "03bec7089b851acf9865fd84b62493624e3970b7",
      "plan_digest": "sha256:322c86ff89de964810f82061879d5aa054f5d0e695e4bab97c9f0286471f655c",
      "source_digest": "bb27090aa8a633522f3c1e05b70672a72890b5a4",
      "acceptance_rows": [
        "R01",
        "R02",
        "R03",
        "R04",
        "R05",
        "R06",
        "R07",
        "R08",
        "R09",
        "R17",
        "R18"
      ],
      "verification": [
        {
          "id": "t1-author-t1-kit-rule",
          "performer": "claude:ft335-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bb27090aa8a633522f3c1e05b70672a72890b5a4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft335-t01-author-20261019",
            "digest": "sha256:bacaeddeaa4f508aa7cfe048de65edc6672efefa65b7f6a2a7ee963f9ab606bd",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,26,1,0,null,yes\n  2,bit,swap,failed,2,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,26,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate ./cmd/bench,KitForRoot|ShellWrapper,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1964\nmutant_report[1]{mutant}:\n  1\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,1003,18\n  github.com/gibbonmi/bench/internal/gate,fail,45,8\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,\"wrapper_kit_test.go:117: gate.KitForRoot(\\\"/dev/shm/FJRVYN/t/TestShellWrapperKitMatchesKitForRoot323442746/001/adopted/.bench\\\", \\\"/dev/shm/FJRVYN/t/TestShellWrapperKitMatchesKitForRoot323442746/001/adopted\\\") = \\\"/dev/shm/FJRVYN/t/TestShellWrap… (394 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestKitForRootKeepsAnAdoptedSubdirectoryKit,\"kit_root_test.go:58: KitForRoot(\\\"/dev/shm/FJRVYN/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit1506399746/001/.bench\\\", \\\"/dev/shm/FJRVYN/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit1506399746/002/with space/linked\\\") = \\\"/dev/shm/FJRVYN/t/Test… (403 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,1005,18\n  github.com/gibbonmi/bench/internal/gate,fail,55,8\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,\"wrapper_kit_test.go:117: gate.KitForRoot(\\\"/dev/shm/IKRMWT/t/TestShellWrapperKitMatchesKitForRoot2959542728/001/adopted/.bench\\\", \\\"/dev/shm/IKRMWT/t/TestShellWrapperKitMatchesKitForRoot2959542728/001/adopted\\\") = \\\"/dev/shm/IKRMWT/t/TestShellWr… (398 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestKitForRootKeepsAnAdoptedSubdirectoryKit,\"kit_root_test.go:58: KitForRoot(\\\"/dev/shm/IKRMWT/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit2200884880/001/.bench\\\", \\\"/dev/shm/IKRMWT/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit2200884880/002/with space/linked\\\") = \\\"/dev/shm/IKRMWT/t/Test… (403 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-kit-rule",
          "command": "bench test --package ./internal/gate --run 'KitForRoot|LaneForCommit|LaneFor|KitValue|KitSourceCheckout'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the top-level condition from KitForRoot and record the R03 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft335-t01-author-20261019",
              "digest": "sha256:bacaeddeaa4f508aa7cfe048de65edc6672efefa65b7f6a2a7ee963f9ab606bd",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,26,1,0,null,yes\n  2,bit,swap,failed,2,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,26,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate ./cmd/bench,KitForRoot|ShellWrapper,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1964\nmutant_report[1]{mutant}:\n  1\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,1003,18\n  github.com/gibbonmi/bench/internal/gate,fail,45,8\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,\"wrapper_kit_test.go:117: gate.KitForRoot(\\\"/dev/shm/FJRVYN/t/TestShellWrapperKitMatchesKitForRoot323442746/001/adopted/.bench\\\", \\\"/dev/shm/FJRVYN/t/TestShellWrapperKitMatchesKitForRoot323442746/001/adopted\\\") = \\\"/dev/shm/FJRVYN/t/TestShellWrap… (394 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestKitForRootKeepsAnAdoptedSubdirectoryKit,\"kit_root_test.go:58: KitForRoot(\\\"/dev/shm/FJRVYN/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit1506399746/001/.bench\\\", \\\"/dev/shm/FJRVYN/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit1506399746/002/with space/linked\\\") = \\\"/dev/shm/FJRVYN/t/Test… (403 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,1005,18\n  github.com/gibbonmi/bench/internal/gate,fail,55,8\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,\"wrapper_kit_test.go:117: gate.KitForRoot(\\\"/dev/shm/IKRMWT/t/TestShellWrapperKitMatchesKitForRoot2959542728/001/adopted/.bench\\\", \\\"/dev/shm/IKRMWT/t/TestShellWrapperKitMatchesKitForRoot2959542728/001/adopted\\\") = \\\"/dev/shm/IKRMWT/t/TestShellWr… (398 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestKitForRootKeepsAnAdoptedSubdirectoryKit,\"kit_root_test.go:58: KitForRoot(\\\"/dev/shm/IKRMWT/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit2200884880/001/.bench\\\", \\\"/dev/shm/IKRMWT/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit2200884880/002/with space/linked\\\") = \\\"/dev/shm/IKRMWT/t/Test… (403 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-author-t1-wrapper-differential",
          "performer": "claude:ft335-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bb27090aa8a633522f3c1e05b70672a72890b5a4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft335-t01-author-20261019",
            "digest": "sha256:bacaeddeaa4f508aa7cfe048de65edc6672efefa65b7f6a2a7ee963f9ab606bd",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,26,1,0,null,yes\n  2,bit,swap,failed,2,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,26,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate ./cmd/bench,KitForRoot|ShellWrapper,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1964\nmutant_report[1]{mutant}:\n  1\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,1003,18\n  github.com/gibbonmi/bench/internal/gate,fail,45,8\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,\"wrapper_kit_test.go:117: gate.KitForRoot(\\\"/dev/shm/FJRVYN/t/TestShellWrapperKitMatchesKitForRoot323442746/001/adopted/.bench\\\", \\\"/dev/shm/FJRVYN/t/TestShellWrapperKitMatchesKitForRoot323442746/001/adopted\\\") = \\\"/dev/shm/FJRVYN/t/TestShellWrap… (394 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestKitForRootKeepsAnAdoptedSubdirectoryKit,\"kit_root_test.go:58: KitForRoot(\\\"/dev/shm/FJRVYN/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit1506399746/001/.bench\\\", \\\"/dev/shm/FJRVYN/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit1506399746/002/with space/linked\\\") = \\\"/dev/shm/FJRVYN/t/Test… (403 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,1005,18\n  github.com/gibbonmi/bench/internal/gate,fail,55,8\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,\"wrapper_kit_test.go:117: gate.KitForRoot(\\\"/dev/shm/IKRMWT/t/TestShellWrapperKitMatchesKitForRoot2959542728/001/adopted/.bench\\\", \\\"/dev/shm/IKRMWT/t/TestShellWrapperKitMatchesKitForRoot2959542728/001/adopted\\\") = \\\"/dev/shm/IKRMWT/t/TestShellWr… (398 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestKitForRootKeepsAnAdoptedSubdirectoryKit,\"kit_root_test.go:58: KitForRoot(\\\"/dev/shm/IKRMWT/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit2200884880/001/.bench\\\", \\\"/dev/shm/IKRMWT/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit2200884880/002/with space/linked\\\") = \\\"/dev/shm/IKRMWT/t/Test… (403 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-wrapper-differential",
          "command": "bench test --package ./cmd/bench --run 'ShellWrapper'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the top-level condition from KitForRoot and record the R09 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft335-t01-author-20261019",
              "digest": "sha256:bacaeddeaa4f508aa7cfe048de65edc6672efefa65b7f6a2a7ee963f9ab606bd",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,26,1,0,null,yes\n  2,bit,swap,failed,2,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,26,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate ./cmd/bench,KitForRoot|ShellWrapper,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1964\nmutant_report[1]{mutant}:\n  1\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,1003,18\n  github.com/gibbonmi/bench/internal/gate,fail,45,8\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,\"wrapper_kit_test.go:117: gate.KitForRoot(\\\"/dev/shm/FJRVYN/t/TestShellWrapperKitMatchesKitForRoot323442746/001/adopted/.bench\\\", \\\"/dev/shm/FJRVYN/t/TestShellWrapperKitMatchesKitForRoot323442746/001/adopted\\\") = \\\"/dev/shm/FJRVYN/t/TestShellWrap… (394 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestKitForRootKeepsAnAdoptedSubdirectoryKit,\"kit_root_test.go:58: KitForRoot(\\\"/dev/shm/FJRVYN/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit1506399746/001/.bench\\\", \\\"/dev/shm/FJRVYN/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit1506399746/002/with space/linked\\\") = \\\"/dev/shm/FJRVYN/t/Test… (403 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,1005,18\n  github.com/gibbonmi/bench/internal/gate,fail,55,8\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestShellWrapperKitMatchesKitForRoot/an_adopted_repo's_own_kit_subdirectory,\"wrapper_kit_test.go:117: gate.KitForRoot(\\\"/dev/shm/IKRMWT/t/TestShellWrapperKitMatchesKitForRoot2959542728/001/adopted/.bench\\\", \\\"/dev/shm/IKRMWT/t/TestShellWrapperKitMatchesKitForRoot2959542728/001/adopted\\\") = \\\"/dev/shm/IKRMWT/t/TestShellWr… (398 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestKitForRootKeepsAnAdoptedSubdirectoryKit,\"kit_root_test.go:58: KitForRoot(\\\"/dev/shm/IKRMWT/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit2200884880/001/.bench\\\", \\\"/dev/shm/IKRMWT/t/TestKitForRootKeepsAnAdoptedSubdirectoryKit2200884880/002/with space/linked\\\") = \\\"/dev/shm/IKRMWT/t/Test… (403 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-author-t1-plumbing-owner",
          "performer": "claude:ft335-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bb27090aa8a633522f3c1e05b70672a72890b5a4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft335-t01-author-20261019",
            "digest": "sha256:345201fc1e4cd3d6d8387c343e9f7684d9ce80142956facbc436414a3a33249d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  git-plumbing-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,122,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-plumbing-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "t1-author-t1-root-conformance",
          "performer": "claude:ft335-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bb27090aa8a633522f3c1e05b70672a72890b5a4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft335-t01-author-20261019",
            "digest": "sha256:18b9bb50f0ad7a9f08501c28844d81ccee4ff9c5a2992cfb9afa1b7852bc4630",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17980,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-author-t1-gate-package",
          "performer": "claude:ft335-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bb27090aa8a633522f3c1e05b70672a72890b5a4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft335-t01-author-20261019",
            "digest": "sha256:0d0a34a2cbbc927b6f85c1af58522241e397660926399f752d8737391559361c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,32461,498\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-gate-package",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t1-author-t1-cmd-package",
          "performer": "claude:ft335-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bb27090aa8a633522f3c1e05b70672a72890b5a4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft335-t01-author-20261019",
            "digest": "sha256:6c6c15606f9fb63d53fa92d9cdc05094d4ca5f2352e20e3e19866ae94c893922",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19747,352\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-cmd-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
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
