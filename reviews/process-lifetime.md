# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/process-lifetime/spec.md",
  "plan_digest": "sha256:af4cfdbe03e0fa38cd38cd92c2f3781a1f7bcedc6bc1446394bf266282ef64b4",
  "implementation_session": "",
  "chunks": [
    {
      "id": "PL-C1",
      "base": "4115afad73961d526b3df1015ee85ea4a71c84b3",
      "tip": "36c2fa86aac94a56792529dbdae08ec1e6652701",
      "plan_digest": "sha256:c36d67f25b6ac6ce1b621f6b196cb4fb4e063b2830bdb5ba4bb7405337a33930",
      "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
      "acceptance_rows": [
        "PL2",
        "PL3",
        "PL4",
        "PL5",
        "PL6",
        "PL7",
        "PL8",
        "PL9",
        "PL10",
        "PL11",
        "PL12",
        "PL13",
        "PL14",
        "PL15",
        "PL16",
        "PL17",
        "PL18",
        "PL19",
        "PL28",
        "PL66",
        "PL67",
        "PL68",
        "PL69",
        "PL124",
        "PL125"
      ],
      "verification": [
        {
          "id": "t01-ticket-1-check-1",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:d02b1d252f095bc1a68b2d08d13b27fb79abeca529e5e2ed98cfc2b4a45ffcd4",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,3550,33\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-1",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-2",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:f2e5cf66d7a4b234db06f8c197763ccaced30893be4af9b50849b26e8d16d237",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,pass,3952,63\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-2",
          "command": "bench test --package ./internal/bounds",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-3",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:f216035271ebec65897d39691c51517379597ea2c3bf3e11f1f0941dc3b0d8a7",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14354,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-4",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:e9485480b2bce62efbd4b0cc4d3577f87b08871076072de8596c5c7996e6d872",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17301,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-5",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:8fd2377ba3e639c452f16dcfe3b2e222146a6dd5323617eb6f9a2e14da707df9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,102,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-5",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-6",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:f780a8547c4c5b282feef2e0c4dc3c6f416d1e85c72b0afb46d69fb0b7f205d6",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,36256,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-1-check-6",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-7",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:c9eb8ae6de1c44a7dea8bf2e89c1df75f6274721739ba23191590976376f6f3e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-7",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-8",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:0570173198b71d68ffdbfaa4facbe957045d752e3f35aaf27496ec6af277bf05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-8",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-9",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:0aaabee70dfcea4b354fdc4790e3db36125abe22a576723f2b6660764c8834f7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-9",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-10",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:085563a55cd74cec959ff54381673a618521b021277e421cb1b052c3c05d4e05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,2,4\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-10",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-check-11",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:933b0f0474c12b388587dac721f21cf6738a3f059a3c401dc1a13c47b5acfe41",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-11",
          "command": "bench test --package ./internal/subprocess/processtest",
          "exit_code": 0
        },
        {
          "id": "t01-ticket-1-omission",
          "performer": "claude:plt-t01-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
            "digest": "sha256:88e034e44aa5d204af0d569c80a6676d44f0f4430d05903af238f5cf292768ee",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,3,TestLifetimeBlockedStart,33,1,0,null,yes\n  2,bit,insert,failed,2,TestLifetimeExitAndIncomplete,28,1,0,null,yes\n  3,bit,swap,failed,7,TestLifetimeBlockedInput,33,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2535\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,64484,33\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedStart,\"lifetime_test.go:54: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeLateStart,\"lifetime_test.go:68: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeStartingCancel,\"lifetime_test.go:42: timeout: waited 20.6s for the lifetime to return\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,2065,28\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeExitAndIncomplete,\"lifetime_test.go:123: cleanup failure replaced the observed exit: {Started:true StartErr:<nil> PID:870109 PGID:870109 Waited:true Code:0 Signal:signal 0 WaitErr:waitid: no child processes}\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeSingleWaiter,\"lifetime_test.go:97: one launch created 2 waiters\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,145173,33\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedInput,\"lifetime_test.go:286: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedWriter,\"lifetime_test.go:174: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeDelayedCapture,\"lifetime_test.go:134: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeExitAndIncomplete,\"lifetime_test.go:121: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeInheritedDescriptor,\"lifetime_test.go:218: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeSingleFinalWindow,\"lifetime_test.go:252: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeWorkerOwnership,\"lifetime_test.go:265: timeout: waited 20.6s for the lifetime to return\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-omission",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit remembered cancellation for a held Start, start a second waiter, or wait unconditionally on a blocked consumer. The corresponding real-child or controlled-start assertion must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t01-author-20261022",
              "digest": "sha256:88e034e44aa5d204af0d569c80a6676d44f0f4430d05903af238f5cf292768ee",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,3,TestLifetimeBlockedStart,33,1,0,null,yes\n  2,bit,insert,failed,2,TestLifetimeExitAndIncomplete,28,1,0,null,yes\n  3,bit,swap,failed,7,TestLifetimeBlockedInput,33,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2535\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,64484,33\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedStart,\"lifetime_test.go:54: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeLateStart,\"lifetime_test.go:68: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeStartingCancel,\"lifetime_test.go:42: timeout: waited 20.6s for the lifetime to return\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,2065,28\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeExitAndIncomplete,\"lifetime_test.go:123: cleanup failure replaced the observed exit: {Started:true StartErr:<nil> PID:870109 PGID:870109 Waited:true Code:0 Signal:signal 0 WaitErr:waitid: no child processes}\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeSingleWaiter,\"lifetime_test.go:97: one launch created 2 waiters\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,145173,33\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedInput,\"lifetime_test.go:286: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedWriter,\"lifetime_test.go:174: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeDelayedCapture,\"lifetime_test.go:134: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeExitAndIncomplete,\"lifetime_test.go:121: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeInheritedDescriptor,\"lifetime_test.go:218: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeSingleFinalWindow,\"lifetime_test.go:252: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeWorkerOwnership,\"lifetime_test.go:265: timeout: waited 20.6s for the lifetime to return\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01r1-ticket-1-check-1",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:4109c48dc61bc8c333e8379cee677dacaa4d9a50765062512e3c9af35839d252",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,3599,33\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-1",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-2",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:99e664978059e7d12ad0b26b10fb892199d2e46f51e274348cd56be954ce3a24",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,pass,3940,63\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-2",
          "command": "bench test --package ./internal/bounds",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-3",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:8c71de58f6fd086d5d2aa72e1f99e8f5f460aa49bc839d2c8effd7bbd565e370",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14094,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-4",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:6098d83a3b12179138af3798b1df2d63bc9b8d6234422d2a9f1e818fcd2d9c05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16832,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-5",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:87e3ed0085cedfbe08518008f388b05d62c87ec20bcafb3ed30b1d729b8bc928",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,106,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-5",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-6",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:eb356f895f2a180bdc8e787badf910b528041d98baf39b59d60a1359b45d3071",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,35392,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-1-check-6",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-7",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:c9eb8ae6de1c44a7dea8bf2e89c1df75f6274721739ba23191590976376f6f3e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-7",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-8",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:0570173198b71d68ffdbfaa4facbe957045d752e3f35aaf27496ec6af277bf05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-8",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-9",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:0aaabee70dfcea4b354fdc4790e3db36125abe22a576723f2b6660764c8834f7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-9",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-10",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:085563a55cd74cec959ff54381673a618521b021277e421cb1b052c3c05d4e05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,2,4\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-10",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-check-11",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:933b0f0474c12b388587dac721f21cf6738a3f059a3c401dc1a13c47b5acfe41",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-check-11",
          "command": "bench test --package ./internal/subprocess/processtest",
          "exit_code": 0
        },
        {
          "id": "t01r1-ticket-1-omission",
          "performer": "claude:plt-t01-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
            "digest": "sha256:52b5e1edce3736ee3a969ff8c074ddc20fc07208a13e48499808778bd10afe41",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,3,TestLifetimeBlockedStart,33,1,0,null,yes\n  2,bit,insert,failed,3,TestLifetimeExitAndIncomplete,28,1,0,null,yes\n  3,bit,swap,failed,7,TestLifetimeBlockedInput,33,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2869\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,64571,33\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedStart,\"lifetime_test.go:54: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeLateStart,\"lifetime_test.go:70: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeStartingCancel,\"lifetime_test.go:42: timeout: waited 20.6s for the lifetime to return\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,2157,28\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeExitAndIncomplete,\"lifetime_test.go:125: cleanup failure replaced the observed exit: {Started:true StartErr:<nil> PID:971270 PGID:971270 Waited:true Code:0 Signal:signal 0 WaitErr:waitid: no child processes}\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeResistantGroup,\"lifetime_test.go:194: TERM-resistant group after grace: process={Started:true StartErr:<nil> PID:971291 PGID:971291 Waited:true Code:0 Signal:signal 0 WaitErr:waitid: no child processes} cleanup={Status:complete Obligations:[] Group:{State:… (274 bytes)\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeSingleWaiter,\"lifetime_test.go:99: one launch created 2 waiters\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,145211,33\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedInput,\"lifetime_test.go:319: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedWriter,\"lifetime_test.go:180: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeDelayedCapture,\"lifetime_test.go:136: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeExitAndIncomplete,\"lifetime_test.go:123: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeInheritedDescriptor,\"lifetime_test.go:224: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeSingleFinalWindow,\"lifetime_test.go:269: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeWorkerOwnership,\"lifetime_test.go:298: timeout: waited 20.6s for the lifetime to return\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-1-omission",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit remembered cancellation for a held Start, start a second waiter, or wait unconditionally on a blocked consumer. The corresponding real-child or controlled-start assertion must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t01-repair1-20261022",
              "digest": "sha256:52b5e1edce3736ee3a969ff8c074ddc20fc07208a13e48499808778bd10afe41",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,3,TestLifetimeBlockedStart,33,1,0,null,yes\n  2,bit,insert,failed,3,TestLifetimeExitAndIncomplete,28,1,0,null,yes\n  3,bit,swap,failed,7,TestLifetimeBlockedInput,33,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2869\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,64571,33\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedStart,\"lifetime_test.go:54: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeLateStart,\"lifetime_test.go:70: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeStartingCancel,\"lifetime_test.go:42: timeout: waited 20.6s for the lifetime to return\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,2157,28\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeExitAndIncomplete,\"lifetime_test.go:125: cleanup failure replaced the observed exit: {Started:true StartErr:<nil> PID:971270 PGID:971270 Waited:true Code:0 Signal:signal 0 WaitErr:waitid: no child processes}\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeResistantGroup,\"lifetime_test.go:194: TERM-resistant group after grace: process={Started:true StartErr:<nil> PID:971291 PGID:971291 Waited:true Code:0 Signal:signal 0 WaitErr:waitid: no child processes} cleanup={Status:complete Obligations:[] Group:{State:… (274 bytes)\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeSingleWaiter,\"lifetime_test.go:99: one launch created 2 waiters\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,145211,33\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedInput,\"lifetime_test.go:319: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeBlockedWriter,\"lifetime_test.go:180: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeDelayedCapture,\"lifetime_test.go:136: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeExitAndIncomplete,\"lifetime_test.go:123: timeout: waited 20.6s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeInheritedDescriptor,\"lifetime_test.go:224: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeSingleFinalWindow,\"lifetime_test.go:269: timeout: waited 20.45s for the lifetime to return\",1\n  github.com/gibbonmi/bench/internal/subprocess,TestLifetimeWorkerOwnership,\"lifetime_test.go:298: timeout: waited 20.6s for the lifetime to return\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c1-standards-r1",
          "performer": "claude:plt-c1-standards-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c1-standards-20261022",
            "digest": "sha256:35f1e24bc06ecc85cf6d3088a4408eb686a153fcc64460f09176815941fcd9cf",
            "excerpt": "## Standards axis, PL-C1 (base 4115afad, tip 7769a5ea)\n\nEvidence check: --check-current returned current=true, dirty_record=false.\n\nFindings\nS1 | auto-fix | 6 | internal/subprocess/control.go:248-249 | AGENTS.md \"one source per fact... two derivations of the same fact must collapse into one source\"; protection/observe.go:44 doc \"Package protection owns the observations...\" | controller.signal reclassifies ESRCH as GroupAbsent on its own, which duplicates the ESRCH-only-proves-absence rule that protection.observeGroup (observe.go:79-93) owns.\nS2 | auto-fix | 6 | internal/subprocess/lifetime_test.go:246-250 | bench-craft-tdd SKILL.md:101 \"a test whose verdict depends on wall-clock time... is a broken oracle\" | TestLifetimeSingleFinalWindow relies on a raw 50ms pause so that the owner consumes waitDone before fire. On a loaded host, run's select can take final.c first. The snapshot then carries ObligationWaiter, and requireIncomplete(..., ObligationStdout) goes red with no code defect. Also, alive(leader) is true for a zombie, so the loop proves only reaping, not that the owner consumed the waiter's result.\n\nAdvice (no IDs, not counted)\n- bounds.ProcessPoll is an exported mutable package var. A function or const-backed constructor would prevent reassignment.\n- Process.PGID always equals PID (outcome.go snapshot). Keep it only if the spec needs the separate field.\n- ProcessShutdownGrace has no production consumer yet. The ticket mandates it for later callers, so it is not a defect.\n\nCount and worst issue: 2 findings. The worst is S2, a scheduling-dependent verdict in a PL66 test.\n"
          },
          "axis": "Standards",
          "base": "4115afad73961d526b3df1015ee85ea4a71c84b3",
          "tip": "2e56083604aa6706a47923d03e4b533887c780b2",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c1-spec-r1",
          "performer": "claude:plt-c1-spec-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c1-spec-20261022",
            "digest": "sha256:f9ab11da18102b7da16868239c9a613a19da21c4512b9d28bc9ca0b2d24edd42",
            "excerpt": "Spec axis, PL-C1 (claude:plt-c1-spec-20261022)\n\nThe evidence check reported current=true at tip 7769a5ea. I found no Spec findings.\n\nFindings: None.\n\nI traced every row ticket 01 covers (PL2-PL19, PL28, PL66-PL69, PL124-PL125) from the coverage map to its cited test. Each test exercises the stated behavior at the stated seam, and each row's why-it-catches clause would fail under the named mutation. internal/subprocess/lifetime_test.go lines 20-291 cover every row except PL15-18 and PL28/PL67; internal/subprocess/protection/observe_test.go lines 26-56 cover PL15-18; internal/bounds/lifetime_test.go lines 18-73 cover PL28 and PL67. The controller (internal/subprocess/control.go) matches spec lines 127-139 and 193-204.\n\nThe five accepted interface calls: none is behavioral (window values keep the registry durations; Registration and Protection belong to ticket 02; the signal through Cancel matches spec line 123; every RunOutput caller passes a bytes.Buffer read after return; processtest is a test fixture only).\n\nNo edit falls outside the acceptance's reach.\n\nAdvice\n- internal/subprocess/streams.go:57 clears cmd.Stdin, so bounds.run now drops a stdin that a caller sets on the command. None of the bounds.Run/RunOutput callers sets Stdin, so this is a latent trap only.\n- TestLifetimeResistantGroup does not prove that TERM came before KILL. The PL12 witness requires only missing escalation.\n- bounds.ProcessPoll puts a 10ms literal outside the const registry.\n\nCount: 0 findings. Worst issue: none.\n"
          },
          "axis": "Spec",
          "base": "4115afad73961d526b3df1015ee85ea4a71c84b3",
          "tip": "2e56083604aa6706a47923d03e4b533887c780b2",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1-coverage-r1",
          "performer": "claude:plt-c1-coverage-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dd9860f250f48ac55476ff340f5c1c0d7de8b498",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c1-coverage-20261022",
            "digest": "sha256:5c35f23bc9b1b3ea5dd9719eb2fd8d188b7d6751447c3f5caff87d94c4bd281c",
            "excerpt": "Coverage axis, PL-C1. Preflight evidence check: current=true, dirty_record=false. All probes read restored=yes, and none was invalid.\n\nFindings\nC1 | auto-fix | 7 | internal/subprocess/control.go:177 | probe silent, restored yes; spec.md:136 \"If that deadline already expired, the owner forces termination without another grace.\" | TestLifetimeLateStart (PL125) uses no grace and SIGKILL, so it cannot tell the c.forced KILL from the remembered signal.\nC2 | ask-user | 5 | internal/bounds/bounds.go:257 | probe silent, restored yes; ticket 01 line 40 \"A raw child status cannot certify required cleanup or publication.\" | When a captured stream fails or is force-closed, bounds should not return ProcessComplete, but no test fails when that check is removed.\nC3 | auto-fix | 5 | internal/subprocess/outcome.go:128 | probe silent, restored yes; spec.md:146 Cleanup \"Not requested, complete, or incomplete\". | A forced closure through the OutputDelay window, or a drain, can report CleanupNotRequested and no test notices.\nC4 | no-op | 4 | internal/subprocess/control.go:81 | probe silent, restored yes; spec.md:159 \"preserves caller platform attributes\". | No C1 test bites this; PL89 (a later ticket, builder Pdeathsig) is the planned biting row.\n\nAdvice\n- Nothing tests the Unbounded branch in ProcessWindow (bounds.go:199); no C1 caller passes Unbounded.\n- Drain mode (control.go:451) has no C1 test; PL31 and PL34 cover it later.\n- Nothing tests StreamFailed, Cleanup.SignalErr, the stdin EPIPE swallow, the pgid<=1 guard, or a second Start and a Wait before Start.\n- TestLifetimeSingleFinalWindow orders its events with a fixed 50 ms pause.\n\nTotals: 4 findings. Worst is C1.\n\nProbes: bounds.go:257 stream-status check made false, silent, yes; bounds.go:199 Unbounded check made false, silent, yes; outcome.go:128 made case c.stopping, silent, yes; control.go:177 if c.forced made false, silent, yes; control.go:81 SysProcAttr == nil made true, silent, yes.\n"
          },
          "axis": "Coverage",
          "base": "4115afad73961d526b3df1015ee85ea4a71c84b3",
          "tip": "2e56083604aa6706a47923d03e4b533887c780b2",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": []
        },
        {
          "id": "c1-standards-r2",
          "performer": "claude:plt-c1-standards-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c1-standards-r2-20261022",
            "digest": "sha256:1e76da24b1ba71874d89f3a2063a537f22b06edc962f9b6332b8f3e15c8a9094",
            "excerpt": "PL-C1 Standards confirming round r2: both r1 findings are closed, and the repair adds no new finding. The evidence check returned current=true at tip 746689e57, with dirty_record=false.\n\nS1: closed. protection.Classify at internal/subprocess/protection/observe.go:47-51 is now the only place that decides ESRCH means absence. observeGroup delegates to it (observe.go:37). controller.signal at internal/subprocess/control.go:246-253 stores the result only when Classify says the group is absent, and keeps every other non-nil error as signalErr. The doc for Cleanup in outcome.go now states that same rule, so it does not restate ESRCH.\n\nS2: closed. The 50 ms pause and the alive(leader) loop are gone. The test now syncs on the polled channel (internal/subprocess/lifetime_test.go:72-79, 88-92); the time.After there only bounds a failure. observe returns early while waitPending is true (control.go:258), so Poll runs only from onWait (control.go:195-196), after the owner took the waiter result. The joinGroup zombie keeps the group present (lifetime_test.go:279). requireIncomplete(t, o, ObligationGroup, ObligationStdout) at lifetime_test.go:274 leaves out the waiter obligation, so an unconsumed waiter result turns the test red.\n\nNew findings: none. Every added or changed comment sentence is 25 words or fewer and describes current code.\n\nFinding count: 0 open, 2 closed.\n\nAdvice\n- starts is written on the owner's goroutine and read on the test goroutine (lifetime_test.go:97); it predates the repair and is ordered only through waitWithin's channel.\n"
          },
          "axis": "Standards",
          "base": "4115afad73961d526b3df1015ee85ea4a71c84b3",
          "tip": "36c2fa86aac94a56792529dbdae08ec1e6652701",
          "finding_ids": [],
          "supersedes": [
            "c1-standards-r1"
          ]
        },
        {
          "id": "c1-spec-r2",
          "performer": "claude:plt-c1-spec-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c1-spec-r2-20261022",
            "digest": "sha256:4a5626cc9b04136a08071ac4e105c4bda2edd3351e48157e52f45ef675af719a",
            "excerpt": "Spec axis, PL-C1 round r2 (claude:plt-c1-spec-r2-20261022)\n\nThe evidence check returned current=true and dirty_record=false at tip 746689e57.\n\nConfirmations\n- PL10 TestLifetimeForcedClose (internal/subprocess/lifetime_test.go:184-196): uses the OutputDelay path. It still asserts StreamClosed with data \"early\", and the new assertion requires CleanupComplete, which outcome.go:128 (c.closed) supplies. This matches spec.md:146.\n- PL125 TestLifetimeLateStart (lifetime_test.go:62-76): a TERM-resistant child, with SIGTERM under a real grace; the Start is released after grace and final expire. control.go:177-180 (c.forced, then SIGKILL) is the only path that reaps the child inside final. It still asserts one waiter. This matches spec.md:136 and closes r1 C1.\n- PL66 TestLifetimeSingleFinalWindow (lifetime_test.go:241-271): the 50 ms pause is gone; polled fires only through observe() after the waiter result (control.go:258); starts == 1 still enforces one shared window.\n- PL15-PL18: observe_test.go:25-56 still drives ObserveGroup and observeGroup, which delegate to Classify (observe.go:42). Test names are unchanged, so no citation drifted.\n- Classify (observe.go:47-56): nil means present, ESRCH means absent, every other error means uncertain and keeps its Err. This matches spec.md:200-202. control.go:246-254 keeps the signal error as evidence unless the answer proves absence (spec.md:203).\n- No edit falls outside the acceptance's reach.\n\nFindings: None.\n\nAdvice\n- No test exercises how controller.signal grades EPERM into SignalErr; r1 advice already noted this.\n\nCount: 0 findings.\n"
          },
          "axis": "Spec",
          "base": "4115afad73961d526b3df1015ee85ea4a71c84b3",
          "tip": "36c2fa86aac94a56792529dbdae08ec1e6652701",
          "finding_ids": [],
          "supersedes": [
            "c1-spec-r1"
          ]
        },
        {
          "id": "c1-coverage-r2",
          "performer": "claude:plt-c1-coverage-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6a961a5af5525e0724759d154bcc7d615c32ed9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c1-coverage-r2-20261022",
            "digest": "sha256:0ae9f6aec46917726250b8f06af1e1c0246094cb364c7100a32688c3d754478b",
            "excerpt": "Coverage axis, PL-C1 confirming round r2 (claude:plt-c1-coverage-r2-20261022). The evidence check returned current=true and dirty_record=false at tip 746689e5. Every probe reports restored=yes.\n\nC1: closed. My probe omitted the KILL in the forced branch of onStart (internal/subprocess/control.go:177-179). TestLifetimeLateStart bit at lifetime_test.go:72 and timed out on completion. With the earlier c.forced->false and KILL->remembered-signal probes, three different mutations of spec.md:136 now bite.\n\nC3: closed. My probe changed the case at internal/subprocess/outcome.go:128 so it no longer includes c.closed. TestLifetimeForcedClose bit at lifetime_test.go:173 with \"Status:not-requested, want complete\" (spec.md:146).\n\nC2: no-op reaffirmed. Under processPolicy, the loop at internal/bounds/bounds.go:256-260 cannot be reached without a prior cancellation: the policy sets no OutputDelay, Normal is Preserve, and it has no Drain (bounds.go:213-217), so closeOwned (control.go:122,125) runs only on final-window expiry; beginFinal runs only from force() or the Drain branch of onWait (control.go:186-190); StreamFailed (streams.go:123) needs a read error from bytes.Buffer.ReadFrom on the owner end of the pipe (streams.go:95-96); a cancellation returns early at bounds.go:250-253. The branch is defense-in-depth with no binding requirement that a test bites it.\n\nC4: no-op unchanged.\n\nProbes\n- control.go onStart | omit c.signal(syscall.SIGKILL) (occurrence 1) | bit | restored yes\n- outcome.go:128 | swap to case c.stopping || c.needGroup: | bit | restored yes\n- control.go signal | omit c.group = answer | silent across ./internal/subprocess/... | restored yes\n\nNew findings: none. The silent probe hit the refactored Classify branch in signal; the mutant is effectively equivalent because ESRCH comes only after reaping, and onWait->observe() then records the absence (control.go:196-197,258-260). No spec row binds the early record.\n\nAdvice\n- A direct witness for the ESRCH absence that signal records would have to separate that record from the observe() that follows.\n\nFinding count: 0 new; C1 and C3 closed, C2 and C4 no-op.\n"
          },
          "axis": "Coverage",
          "base": "4115afad73961d526b3df1015ee85ea4a71c84b3",
          "tip": "36c2fa86aac94a56792529dbdae08ec1e6652701",
          "finding_ids": [],
          "supersedes": [
            "c1-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "PL-C2",
      "base": "36c2fa86aac94a56792529dbdae08ec1e6652701",
      "tip": "3e9708ec537419d91b9740dd5238fe13f02e874a",
      "plan_digest": "sha256:dbb8b18b9518678f70a2d2a437f7ae9e355507e4dbb28f1b411810cc37504e8d",
      "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
      "acceptance_rows": [
        "PL20",
        "PL21",
        "PL22",
        "PL23",
        "PL24",
        "PL25",
        "PL26",
        "PL27",
        "PL39",
        "PL40",
        "PL41",
        "PL42",
        "PL43",
        "PL44",
        "PL45",
        "PL46",
        "PL47",
        "PL48",
        "PL49",
        "PL50",
        "PL51",
        "PL52",
        "PL53",
        "PL54",
        "PL70",
        "PL71",
        "PL72",
        "PL73",
        "PL74",
        "PL104",
        "PL105",
        "PL106",
        "PL107",
        "PL108",
        "PL109",
        "PL110",
        "PL111",
        "PL112",
        "PL113",
        "PL120"
      ],
      "verification": [
        {
          "id": "t02-ticket-2-check-1",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:bef02d78ea325bfa2622bc04fe27b12b0229fdd32ffc3c1a4c121e021b778c46",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,42,41\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-1",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-2",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:266b18f9ccd057e4ec9c2eecc0e9fceff9f13cc2ed46c912de51789a90f80532",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,pass,3949,63\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-2",
          "command": "bench test --package ./internal/bounds",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-3",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:b37918dd66be754e50593fcc0e87c3bbfb6018dff4e26937d877adb7305cf0b6",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,88,44\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-3",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-4",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:d09c5bc2b5c59d39dd1d17f67a4105937cce52d0e85bdc3d7b4db04fd200751d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,206,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,43365,206\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-4",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-5",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:828b5c5b9ad104e0faf8bddba6ad3e9583c9aeb3a0eeff40bbd572e0ed04bb30",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12942,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-6",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:f98b4745f79e6cd636f671d7bdc9a8afd712d5765463a4e6d955ee41303f7386",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14521,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-7",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:6c926d40204c36382fa02b8a8ae2306b1d1776984fc5cb1f3252f3fbed8d8f2e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,99,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-7",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-8",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:486da63162330289a823861a0907aa14ebba88c49aaa3d179ad64ed56fa608aa",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,36101,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-2-check-8",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-9",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:3b2214bf5e1da03908018b148ff1aa2ec5eb37b4316a9fe4131f64303f8fa2b9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,3915,34\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-9",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-10",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:c9eb8ae6de1c44a7dea8bf2e89c1df75f6274721739ba23191590976376f6f3e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-10",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-11",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:0570173198b71d68ffdbfaa4facbe957045d752e3f35aaf27496ec6af277bf05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-11",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-12",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:0aaabee70dfcea4b354fdc4790e3db36125abe22a576723f2b6660764c8834f7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-12",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-check-13",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:933b0f0474c12b388587dac721f21cf6738a3f059a3c401dc1a13c47b5acfe41",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-13",
          "command": "bench test --package ./internal/subprocess/processtest",
          "exit_code": 0
        },
        {
          "id": "t02-ticket-2-omission",
          "performer": "claude:plt-t02-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
            "digest": "sha256:cc9ee03c435bda8c49f209cec77a66639e866ae0bced16ee6716d613d6d8f9a3",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/bounds/bounds.go,swap,failed,2,TestProcessLifetimeNestedOverlap,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,all,passed,206\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,654\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,206,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,40719,206\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestProcessLifetimeNestedOverlap,\"process_lifetime_test.go:128: provider blocked used the cache with protection \\\"absent\\\\\\\\n\\\" (<nil>), want registered\",1\n  github.com/gibbonmi/bench/internal/systemtest,TestProcessLifetimeNestedRegistration,\"process_lifetime_test.go:102: provider only used the cache with protection \\\"absent\\\\\\\\n\\\" (<nil>), want registered\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-omission",
          "command": "bench test --check system",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop inherited descriptors from the real bounds launch. The nested-provider witness must fail through bench test --check system.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:plt-t02-author2-20261022",
              "digest": "sha256:cc9ee03c435bda8c49f209cec77a66639e866ae0bced16ee6716d613d6d8f9a3",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/bounds/bounds.go,swap,failed,2,TestProcessLifetimeNestedOverlap,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,all,passed,206\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,654\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,206,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,40719,206\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestProcessLifetimeNestedOverlap,\"process_lifetime_test.go:128: provider blocked used the cache with protection \\\"absent\\\\\\\\n\\\" (<nil>), want registered\",1\n  github.com/gibbonmi/bench/internal/systemtest,TestProcessLifetimeNestedRegistration,\"process_lifetime_test.go:102: provider only used the cache with protection \\\"absent\\\\\\\\n\\\" (<nil>), want registered\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t02r1-ticket-2-check-1",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:fd317faa6d52e0a8bb44826e77b74a2704909aba4e95300c994e6a79afd20343",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,41,42\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-1",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-2",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:393b5a759a15604b40d8b993e8538625cc4091922de72fce99b94e576ac39f7b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,pass,3919,64\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-2",
          "command": "bench test --package ./internal/bounds",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-3",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:bca20b2d1d91f66ab3274612cae251e696ee383069be9b8a4ca3ad90e4d4bed7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,97,44\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-3",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-4",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:bd5b16d6d0489dc3db1460bddafb9482438f9c17b8bbe3fccae38c32530fb564",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,206,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,48062,206\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-4",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-5",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:2a3556a753967e927a4608d8547d7a0c57e1f6397991acd49df7b7447ad08bca",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12734,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-6",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:0f59ae512e6f83597a085450e315aa94cf3ea4212bf449dbef2e00dfa94f7f47",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14421,403\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-7",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:8fd2377ba3e639c452f16dcfe3b2e222146a6dd5323617eb6f9a2e14da707df9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,102,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-7",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-8",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:823b4df983f1e9db29cd1cb74871cffd0f478c91475f1e68fd4fdee4e98cb107",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,35735,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-2-check-8",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-9",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:5146800f07abbbf2426af66f089e623741dbe35f7c773bcbcab002f6194cdaeb",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,3923,34\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-9",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-10",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:c9eb8ae6de1c44a7dea8bf2e89c1df75f6274721739ba23191590976376f6f3e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-10",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-11",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:0570173198b71d68ffdbfaa4facbe957045d752e3f35aaf27496ec6af277bf05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-11",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-12",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:0aaabee70dfcea4b354fdc4790e3db36125abe22a576723f2b6660764c8834f7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-12",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-check-13",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:933b0f0474c12b388587dac721f21cf6738a3f059a3c401dc1a13c47b5acfe41",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-check-13",
          "command": "bench test --package ./internal/subprocess/processtest",
          "exit_code": 0
        },
        {
          "id": "t02r1-ticket-2-omission",
          "performer": "claude:plt-t02-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
            "digest": "sha256:9862003a0d778e59f3c1f7f72b483944fc48e4abfa8179d1df34b82936549610",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/bounds/bounds.go,swap,failed,2,TestProcessLifetimeNestedOverlap,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,all,passed,206\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,654\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,206,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,44856,206\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestProcessLifetimeNestedOverlap,\"process_lifetime_test.go:128: provider blocked used the cache with protection \\\"absent\\\\\\\\n\\\" (<nil>), want registered\",1\n  github.com/gibbonmi/bench/internal/systemtest,TestProcessLifetimeNestedRegistration,\"process_lifetime_test.go:102: provider only used the cache with protection \\\"absent\\\\\\\\n\\\" (<nil>), want registered\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-2-omission",
          "command": "bench test --check system",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop inherited descriptors from the real bounds launch. The nested-provider witness must fail through bench test --check system.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:plt-t02-repair1-20261022",
              "digest": "sha256:9862003a0d778e59f3c1f7f72b483944fc48e4abfa8179d1df34b82936549610",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/bounds/bounds.go,swap,failed,2,TestProcessLifetimeNestedOverlap,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,all,passed,206\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,654\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,206,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,44856,206\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestProcessLifetimeNestedOverlap,\"process_lifetime_test.go:128: provider blocked used the cache with protection \\\"absent\\\\\\\\n\\\" (<nil>), want registered\",1\n  github.com/gibbonmi/bench/internal/systemtest,TestProcessLifetimeNestedRegistration,\"process_lifetime_test.go:102: provider only used the cache with protection \\\"absent\\\\\\\\n\\\" (<nil>), want registered\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c2-standards-r1",
          "performer": "claude:plt-c2-standards-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c2-standards-20261022",
            "digest": "sha256:6cffe68cf4cb8f50648bc337489535ce212a6af5adeaf35f9a40b65da4c0a8c6",
            "excerpt": "Standards axis, PL-C2 (FT362-plt; frozen pair 36c2fa86 to 58e4b351, record e4b782c4b)\n\nThe evidence check passed: current=true, dirty_record=false.\n\nFindings\nS1 | no-op | 4 | internal/subprocess/protection/fixture_test.go:138 and internal/systemtest/process_lifetime_test.go:141 | AGENTS.md, fixture harnesses \"remain single-sourced\" | The diff adds two bounded-wait helpers, within[T] and exitWithin, each running a call in a goroutine and selecting against bounds.TestDeadline, the same shape as waitWithin in internal/subprocess/lifetime_helpers_test.go:59. The deadline policy stays single-sourced in bounds, so the repeated part is the select scaffold, close to incidental text. No repair cycle recommended.\n\nRejected candidates: Environ vs capability.WithoutEnvironment (spec.md:246 closes the leaf's imports; capability is the test-skip package); the strict reader (readRecord uses DecodeExactDocument; jsonfile has no no-follow reader); group-absence waits (WaitAbsent uses protection.ObserveGroup; WaitPID now shares waitFor); status literals (ObligationGroup aliases protection.ObligationGroup); the Install global (spec.md:251-254 and the ticket require it); test seams (SetReplacer and SetProbe in export_test.go only; Lifetime.kill follows the start/wait field precedent); comments (no narration or provenance; present tense).\n\nTotals: 1 finding. The worst is S1, a judgment-grade fixture duplication.\n\nAdvice\n- registration.go Started formats r.state in its refusal (\"not protected\"), but the cause is a closed or non-pending registration; naming r.record.Identity would describe it better.\n- If a third production env-filter copy appears, it belongs in a leaf both owners can import.\n\nImplementation command: no command change.\n"
          },
          "axis": "Standards",
          "base": "36c2fa86aac94a56792529dbdae08ec1e6652701",
          "tip": "58e4b351cdb9250de98cbcb8e6e3b3e7df355514",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "c2-spec-r1",
          "performer": "claude:plt-c2-spec-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c2-spec-20261022",
            "digest": "sha256:186abc89d3fdd3175cfed7b2006f16a152fd65358cd9f831fd8396e49dc0e6fe",
            "excerpt": "Spec axis, PL-C2. Evidence current=true. Two findings; neither blocks the production contract. Worst is P1.\n\nFindings\nP1 | auto-fix | medium | internal/subprocess/protection/registration_test.go:204-206 and record_test.go:225-229 | PL71, PL110 (spec 661, 700) | On error, OpenResource returns Descriptor{}, so requireRefused(t, p, d) refuses only because the path is empty. It never launches against the visible-but-unsynced resource.json (PL71) or the control-byte descriptor (PL110); Prepare's own check() refusal of an unsafe descriptor is never exercised.\nP2 | auto-fix | low | internal/subprocess/protection/record_test.go:59-61 | PL39 (spec 629) | The foreign-binding test changes the store and the generation but never the Binding field; a validator or matches that ignored the binding would pass this row.\n\nAdvice\n- TestRegistrationPartialPublish's loop over the first store's records has no count check; noStartUser guards this elsewhere.\n- In TestRecoveryMissingCompletion, the outer Outcome.Protection is Safe while a nested user is unresolved; release authority stays with SealAndObserve, and later consumers must not treat Outcome.Protection as release proof (spec 400-402).\n- After Started fails, the registration closes, so pending records stay indefinitely; conservative and consistent with spec 353-357.\n\nContracts verified in code: Prepare validates the whole vector and admission before its first mutation and releases locks before Start; sealing refuses new users while Started and Finish still update; rollback publishes no-start and returns uncertain on failure; a Started failure forces KILL through the in-memory handle; grade treats pending as protected and never infers no-start, requiring group ESRCH then complete workers or actor ESRCH; recovery probes with signal 0 only; ParseVector and Environ handle BENCH_PROCESS_RESOURCES and bounds.run refuses a malformed vector; schema fields, 0700/0600 modes, and the strict reader match spec 256-287. Outside Writes:, the ticket commit changes only spec.md coverage citations, which spec 584 expects.\n"
          },
          "axis": "Spec",
          "base": "36c2fa86aac94a56792529dbdae08ec1e6652701",
          "tip": "58e4b351cdb9250de98cbcb8e6e3b3e7df355514",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "c2-coverage-r1",
          "performer": "claude:plt-c2-coverage-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9624eb4bfdce87ebee84ef2abb9b05167cb5d621",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c2-coverage-20261022",
            "digest": "sha256:67eee62d7259f339d4339a5c996063af3f82105e843a69c7805e0e7cafbf758f",
            "excerpt": "Coverage axis, PL-C2. Evidence current=true, dirty_record=false. Four findings, all silent probes on behavior the spec states; worst is C3.\n\nFindings\nC1 | auto-fix | high | internal/gocache/protection.go:136 | probe silent, restored yes; spec \"Descriptors cannot choose arbitrary files, symlinks...\" | The canonical-path check has no biting test. The link case in TestProtectionRefusesAForeignCache is still refused by cacheIdentity's Lstat on the final symlink; a symlink parent component would pass without the check.\nC2 | auto-fix | medium | internal/gocache/protection.go:142 | probe silent, restored yes; spec \"The validator checks the resource kind's recognized store location and existing owner identity\" | The lock-file check has no biting test; the unlocked case keeps the original binding, so the binding check refuses it first.\nC3 | auto-fix | high | internal/subprocess/protection/store.go:43-44 | probe silent across ./internal/subprocess/..., ./internal/bounds, ./internal/gocache, restored yes; spec \"Malformed, missing, foreign, and sealed descriptors prevent launch\" and \"each resource owner's identity validator\" | No test covers a kind with no validator row or a nil Protocol; with the mutant a kind with no row skips validation.\nC4 | auto-fix | medium | internal/bounds/bounds.go:226-229 | probe silent on ./internal/bounds and ./internal/subprocess/..., restored yes; spec \"Malformed ... descriptors prevent launch\" | No bounds test sets a malformed BENCH_PROCESS_RESOURCES; ParseVector rows cover only the parser, not the refusal in run.\n\nAdvice\n- bounds.run does not copy Outcome.Protection into ProcessResult; no spec line asks for it.\n- Environ removes the entry for an empty vector; production never calls it that way.\n- The h.dir == \"\" arm of Holder.Protect is unreachable from HoldDir.\n- The errIdentityPending snapshot arm (outcome.go:373-374) was not probed.\n\nProbes: gocache/protection.go lock check made false (silent, yes); gocache/protection.go canonical check made false (silent, yes); protection/store.go no-validator refusal to return nil (silent, yes); bounds/bounds.go malformed check made false (silent, yes).\n"
          },
          "axis": "Coverage",
          "base": "36c2fa86aac94a56792529dbdae08ec1e6652701",
          "tip": "58e4b351cdb9250de98cbcb8e6e3b3e7df355514",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": []
        },
        {
          "id": "c2-standards-r2",
          "performer": "claude:plt-c2-standards-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c2-standards-r2-20261022",
            "digest": "sha256:ca62ffec458be645cf067c4914f60fd9ef22d33c393d6bc5d39e7a2a411d44a2",
            "excerpt": "Standards axis, PL-C2 confirming round r2. Pass. The repair delta adds no new Standards defect. Evidence current=true, dirty_record=false, tip 699f1b4e4.\n\nOne source per fact: the new cases reuse the single-sourced fixtures newResource, newProtocol, requireRefused (internal/subprocess/protection/fixture_test.go:159), resourceRecord, and plant (record_test.go:19, :30). No harness is pasted again. The gocache \"parent link\" case extends the existing table at protection_test.go:71.\n\nComments: all six new comment sentences are 25 words or fewer and describe current code; the longest, record_test.go:92-93, is 22 words.\n\nr1 S1 stays a no-op; this delta does not reopen it.\n\nAdvice: none.\n"
          },
          "axis": "Standards",
          "base": "36c2fa86aac94a56792529dbdae08ec1e6652701",
          "tip": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "finding_ids": [],
          "supersedes": [
            "c2-standards-r1"
          ]
        },
        {
          "id": "c2-spec-r2",
          "performer": "claude:plt-c2-spec-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c2-spec-r2-20261022",
            "digest": "sha256:9ce56e6e2b85979eb3e49edbb05a389224f337c1bb779c33a8cb2aba8eef97f9",
            "excerpt": "Spec axis, PL-C2 round r2. Evidence current=true, dirty_record=false at tip 699f1b4e4. No new findings.\n\n- P1, PL110 half: closed. TestRecordUnsafePath launches the full control-byte descriptor through requireRefused, which checks the marker stays absent, Started is false, and StartErr is set (fixture_test.go:159-165); the launch error must equal the parse check's own error. This matches spec.md:700 and spec.md:750.\n- P1, PL71 half: reaffirmed as no-op. Row PL71 refuses only the first launch (spec.md:661); spec.md:317-318 requires the launch refused and the visible record kept for observation and repair, and asks for no lasting refusal. The OpenResource error is the witness for that row.\n- P2: closed. The test copies the second resource's Binding into the first resource's published record and requires the launch refused (record_test.go:84-86). This covers spec.md:629.\n\nThe C1 to C4 test-only repairs stay within PL-C2's rows; no defect.\n\nAdvice: P2 does not check which error refused the launch; matching the binding-mismatch error would tighten it.\n"
          },
          "axis": "Spec",
          "base": "36c2fa86aac94a56792529dbdae08ec1e6652701",
          "tip": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "finding_ids": [],
          "supersedes": [
            "c2-spec-r1"
          ]
        },
        {
          "id": "c2-coverage-r2",
          "performer": "claude:plt-c2-coverage-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f1894600687d86e6e43827a9509cf8f69a71e77c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c2-coverage-r2-20261022",
            "digest": "sha256:fee7e544d3023e2426352d2313e0568467a000a42bbed78ff0e4ecbb30ef6726",
            "excerpt": "Coverage axis, PL-C2 round r2. Evidence current=true, tip 699f1b4e, not dirty. No probe ran: none of the changed tests looks vacuous on reading. All four r1 findings are closed.\n\n- C1 closed. The \"parent link\" case keeps the store, binding, and lock, and changes only a parent component, so the canonical check at internal/gocache/protection.go:38 is the only refusal. Test: protection_test.go:68-71.\n- C2 closed. unlocked is a sibling directory under the canonical parent whose store and binding match it, so only the lock check at protection.go:44 refuses it. Test: protection_test.go:64-65.\n- C3 closed. newResource makes a real cache generation; with a nil table, an empty table, or a worktree-only table, owner (internal/subprocess/protection/store.go:57-62) refuses first through union (registration.go:97). Test: record_test.go:76-82. The nil *Protocol subcase may fail the mutant by a later nil-receiver panic rather than by admission; unverified and low risk.\n- C4 closed. The malformed guard is at internal/bounds/bounds.go:226-228, and the assertion (internal/bounds/lifetime_test.go:56-61) requires ProcessStart with a non-nil error and no started process, matching the orchestrator's status-swap probe.\n\nNo new findings.\n"
          },
          "axis": "Coverage",
          "base": "36c2fa86aac94a56792529dbdae08ec1e6652701",
          "tip": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "finding_ids": [],
          "supersedes": [
            "c2-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "PL-C3",
      "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
      "tip": "32aeb2db509ee2977b351b7cbfe0e48276a27993",
      "plan_digest": "sha256:af4cfdbe03e0fa38cd38cd92c2f3781a1f7bcedc6bc1446394bf266282ef64b4",
      "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
      "acceptance_rows": [
        "PL29",
        "PL30",
        "PL31",
        "PL34",
        "PL35",
        "PL36",
        "PL126",
        "PL127",
        "PL37",
        "PL38",
        "PL55",
        "PL56",
        "PL57",
        "PL76",
        "PL77",
        "PL78",
        "PL32",
        "PL33",
        "PL58",
        "PL59",
        "PL60",
        "PL75",
        "PL121",
        "PL123",
        "PL25",
        "PL128",
        "PL129"
      ],
      "verification": [
        {
          "id": "t07-ticket-7-check-1",
          "performer": "claude:plt-t07-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t07-author2-20261022",
            "digest": "sha256:00006daa80597321114748758918032580d928cdbe9f2c81726814f6e2007ebd",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,425,42\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-1",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "t07-ticket-7-check-2",
          "performer": "claude:plt-t07-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t07-author2-20261022",
            "digest": "sha256:27261a3df0cf1d97062be6cc6a8b039540536f7e09a4856940fbfea169668fba",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,28577,515\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t07-ticket-7-check-3",
          "performer": "claude:plt-t07-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t07-author2-20261022",
            "digest": "sha256:ccebed2f3a230c597634baea098957a4976309ab1b7ad29724dfd85555feb19c",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,209,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,63404,209\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t07-ticket-7-check-4",
          "performer": "claude:plt-t07-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t07-author2-20261022",
            "digest": "sha256:1a445f7173162db60a23ddbac53708946e5a89200de4c2f3095410d14158ef65",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5214,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-4",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t07-ticket-7-check-5",
          "performer": "claude:plt-t07-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t07-author2-20261022",
            "digest": "sha256:6fecc44ab6e9bb6df6475fec11dd2719c1d082e0adb183df4e2d7c2df0b0a38d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17478,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t07-ticket-7-check-6",
          "performer": "claude:plt-t07-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t07-author2-20261022",
            "digest": "sha256:b1445511accf0940ef3b9bf015193becd9711048b78409296c47727bdffd0748",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20765,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t07-ticket-7-check-7",
          "performer": "claude:plt-t07-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t07-author2-20261022",
            "digest": "sha256:628bb382f16e95ead91f5fa0616465690ff16b64216ddf6dde44561c0a7e1c22",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,124,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-7",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t07-ticket-7-check-8",
          "performer": "claude:plt-t07-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t07-author2-20261022",
            "digest": "sha256:1ed2a6f8d47a672b94c0a0d42db992d29b1724a59fe71c67800cfde3af9bb432",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,47740,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-7-check-8",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t07-ticket-7-omission",
          "performer": "claude:plt-t07-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t07-author2-20261022",
            "digest": "sha256:5a8a3af57732b80db33059b8b6aca92c0c48b914a2b8c4bd3856c98b4bcff488",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,TestDirectClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate/prospectiveartifact,all,passed,42\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,353\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,436,42\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestDirectClose,\"lifetime_test.go:81: Close with a live group user = <nil>, want the retention\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-omission",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore owner-death-only sweep, ignore a nested registered group, or discard lane selection release failure. The real Close, second-process sweep, or green-record witness must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t07-author2-20261022",
              "digest": "sha256:5a8a3af57732b80db33059b8b6aca92c0c48b914a2b8c4bd3856c98b4bcff488",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,TestDirectClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate/prospectiveartifact,all,passed,42\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,353\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,436,42\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestDirectClose,\"lifetime_test.go:81: Close with a live group user = <nil>, want the retention\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t03-ticket-3-check-1",
          "performer": "claude:plt-t03-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t03-author-20261022",
            "digest": "sha256:ad5416d74b2c68deffae0d7837490a9c69706c19ead32b76373bb9a13629d418",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,28492,515\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t03-ticket-3-check-2",
          "performer": "claude:plt-t03-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t03-author-20261022",
            "digest": "sha256:8488bc2fad1dd2bba66f7b06c5fa3290a358bbaa9144ed259adeef5a4aa9fc4a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17030,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t03-ticket-3-check-3",
          "performer": "claude:plt-t03-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t03-author-20261022",
            "digest": "sha256:dd40fda7721b09501f857629a2271895f76afb8e21151e9f7c9acedd97a458c3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,47509,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-3-check-3",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t03-ticket-3-omission",
          "performer": "claude:plt-t03-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t03-author-20261022",
            "digest": "sha256:5111973dc478a4cf81b2e3015298f8efee8889d5e953e15c397c7a5fd2bb506c",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGateInterruptOpensWithINTAndExits130,515,1,0,null,yes\n  2,bit,swap,failed,1,TestGateTimeoutOpensWithTERMAndExits124,515,1,0,null,yes\n  3,bit,swap,failed,1,TestGateNormalDrainRemovesReadyDescendant,515,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1163\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28542,515\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateInterruptOpensWithINTAndExits130,\"lifetime_test.go:56: interrupted phase opener = \\\"TERM\\\", exit = 130, interrupted = true, want INT, 130, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29053,515\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTimeoutOpensWithTERMAndExits124,\"lifetime_test.go:64: timed-out phase opener = \\\"INT\\\", exit = 130, interrupted = true, want TERM, 124, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28992,515\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateNormalDrainRemovesReadyDescendant,\"lifetime_test.go:78: descendant 681589 survived normal completion\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-omission",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace INT with TERM, map timeout to interruption, or omit normal group drain. The actual gate signal and descendant-survival fixtures must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t03-author-20261022",
              "digest": "sha256:5111973dc478a4cf81b2e3015298f8efee8889d5e953e15c397c7a5fd2bb506c",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGateInterruptOpensWithINTAndExits130,515,1,0,null,yes\n  2,bit,swap,failed,1,TestGateTimeoutOpensWithTERMAndExits124,515,1,0,null,yes\n  3,bit,swap,failed,1,TestGateNormalDrainRemovesReadyDescendant,515,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1163\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28542,515\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateInterruptOpensWithINTAndExits130,\"lifetime_test.go:56: interrupted phase opener = \\\"TERM\\\", exit = 130, interrupted = true, want INT, 130, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29053,515\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTimeoutOpensWithTERMAndExits124,\"lifetime_test.go:64: timed-out phase opener = \\\"INT\\\", exit = 130, interrupted = true, want TERM, 124, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,28992,515\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateNormalDrainRemovesReadyDescendant,\"lifetime_test.go:78: descendant 681589 survived normal completion\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t04-ticket-4-check-1",
          "performer": "claude:plt-t04-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t04-author-20261022",
            "digest": "sha256:1c30ba7395a4326f5284fa716d1ae20e8a0aa8a3f72aa7f8ce2dc1cdb0d343f3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,pass,28790,37\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-1",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "t04-ticket-4-check-2",
          "performer": "claude:plt-t04-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t04-author-20261022",
            "digest": "sha256:751b1450ababdf98962188df918b97233d2c9e9416959c10753e6c3d84bbf8cd",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16978,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t04-ticket-4-check-3",
          "performer": "claude:plt-t04-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t04-author-20261022",
            "digest": "sha256:523b474cc972cfb312737f2b58ddebcf03a9cfb4b01f83d77eafd0b08b72cf63",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20117,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-3",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t04-ticket-4-check-4",
          "performer": "claude:plt-t04-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t04-author-20261022",
            "digest": "sha256:3abbc40e0fe367b0d329934b494c738b9d89658c0bf77cee78f532da4781f58b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,48168,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-4-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t04-ticket-4-omission",
          "performer": "claude:plt-t04-author-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t04-author-20261022",
            "digest": "sha256:872e1025c5df963a9882af1831f64f94ed07e8e10f6b34e952eb81509b6b6ba9",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuilderNormalDrain,37,1,0,null,yes\n  2,bit,swap,failed,1,TestBuilderCancel,37,1,0,null,yes\n  3,bit,swap,failed,1,TestBinaryClose,37,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/runbinary,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1196\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,28266,37\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderNormalDrain,\"lifetime_test.go:26: build private Bench executable: builder protection protected: <nil>\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,26490,37\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderCancel,\"lifetime_test.go:90: descendant signal log = \\\"\\\" (open /dev/shm/OGLPFA/t/TestBuilderCancel2760237569/002/signals: no such file or directory), want TERM before KILL\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,28785,37\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:64: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-omission",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove normal builder drain or delete the selection before all registered groups are absent. The live descendant or executable-retention assertion must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t04-author-20261022",
              "digest": "sha256:872e1025c5df963a9882af1831f64f94ed07e8e10f6b34e952eb81509b6b6ba9",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuilderNormalDrain,37,1,0,null,yes\n  2,bit,swap,failed,1,TestBuilderCancel,37,1,0,null,yes\n  3,bit,swap,failed,1,TestBinaryClose,37,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/runbinary,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1196\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,28266,37\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderNormalDrain,\"lifetime_test.go:26: build private Bench executable: builder protection protected: <nil>\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,26490,37\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderCancel,\"lifetime_test.go:90: descendant signal log = \\\"\\\" (open /dev/shm/OGLPFA/t/TestBuilderCancel2760237569/002/signals: no such file or directory), want TERM before KILL\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,28785,37\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:64: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t04r1-ticket-4-check-1",
          "performer": "claude:plt-t04-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair1-20261022",
            "digest": "sha256:39e9471e87f43f66d277fb47450d9ddc0eb3f572ccda8a2e654263b38fed79fa",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,pass,28530,37\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-1",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "t04r1-ticket-4-check-2",
          "performer": "claude:plt-t04-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair1-20261022",
            "digest": "sha256:09b239b513fe127cf26652ca41299c9ea27094d1097d877f9099f800ce4100ca",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17044,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t04r1-ticket-4-check-3",
          "performer": "claude:plt-t04-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair1-20261022",
            "digest": "sha256:092b6b0691a5ab7ef0966281494cb1befd8fa2e75bffa53f7135ea9c711a8383",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20610,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-3",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t04r1-ticket-4-check-4",
          "performer": "claude:plt-t04-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair1-20261022",
            "digest": "sha256:f63f9516afe7c6bae78e3d849ed8f343a78d780a7ba265baa4fc25c1c655ae4b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,47887,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-4-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t04r1-ticket-4-omission",
          "performer": "claude:plt-t04-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair1-20261022",
            "digest": "sha256:e219e223a023a6d44e0d8dc3f5de5a284fef74936dec6f48f9531b59a9674a48",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/runbinary/runbinary.go,swap,failed,1,TestBuilderNormalDrain,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,TestBuilderNormalDrain,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,8,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderNormalDrain,\"lifetime_test.go:26: build private Bench executable: builder protection protected: <nil>\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-omission",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove normal builder drain or delete the selection before all registered groups are absent. The live descendant or executable-retention assertion must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t04-repair1-20261022",
              "digest": "sha256:e219e223a023a6d44e0d8dc3f5de5a284fef74936dec6f48f9531b59a9674a48",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/runbinary/runbinary.go,swap,failed,1,TestBuilderNormalDrain,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,TestBuilderNormalDrain,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,8,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderNormalDrain,\"lifetime_test.go:26: build private Bench executable: builder protection protected: <nil>\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05-ticket-5-check-1",
          "performer": "claude:plt-t05-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t05-author2-20261022",
            "digest": "sha256:1532b4f82feb14ef7a3b7223f744893f58bc3eae265bd4c02d890a2df210ec67",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,pass,1433,49\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-1",
          "command": "bench test --package ./internal/env",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-check-2",
          "performer": "claude:plt-t05-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t05-author2-20261022",
            "digest": "sha256:30b1d104b9dc04bcf7d93a5f080e0c43bb8e38d3f6ec78812a68e1ed9aab06ff",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,29749,515\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-check-3",
          "performer": "claude:plt-t05-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t05-author2-20261022",
            "digest": "sha256:367d37c06d26bed348d053612b019903dde190674cf87ba09cf7d2b18f0e1048",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17410,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-check-4",
          "performer": "claude:plt-t05-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t05-author2-20261022",
            "digest": "sha256:7797bfc984683d505fa438de37b9cdc9dfbb581d6e2341a89020b4d68edabcb7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20633,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-check-5",
          "performer": "claude:plt-t05-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t05-author2-20261022",
            "digest": "sha256:9059a6d391f72a6d7df9f698b724ef943c7a468436230ae689d376bc6ae5f09c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,47771,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-5-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-check-6",
          "performer": "claude:plt-t05-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t05-author2-20261022",
            "digest": "sha256:648dd273f66dd3e29890e071b413b9b2674d8993054fe0b3e26f37b1c30fc825",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4042,35\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-6",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-check-7",
          "performer": "claude:plt-t05-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t05-author2-20261022",
            "digest": "sha256:0b906dd0d6e7bd2515471d74bd0841bd9db5c2bc3916cfacb1804cfd60386074",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,51,42\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-7",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t05-ticket-5-omission",
          "performer": "claude:plt-t05-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t05-author2-20261022",
            "digest": "sha256:75f94af770142b8a8d9e249b807e0ba3aad367a8eec2eb3dba59705deff2255a",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestKitRunClose,49,1,0,null,yes\n  2,bit,swap,failed,1,TestKitRunClose,49,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/env,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,884\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,1422,49\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:40: user outcome = {Process:{Started:true StartErr:<nil> PID:816207 PGID:816207 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[]}… (506 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,1325,49\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:44: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-omission",
          "command": "bench test --package ./internal/env",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the private-run descriptor or discard KitTestRun.Close failure. The child-owned file survives while a successful phase result must disappear.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t05-author2-20261022",
              "digest": "sha256:75f94af770142b8a8d9e249b807e0ba3aad367a8eec2eb3dba59705deff2255a",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestKitRunClose,49,1,0,null,yes\n  2,bit,swap,failed,1,TestKitRunClose,49,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/env,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,884\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,1422,49\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:40: user outcome = {Process:{Started:true StartErr:<nil> PID:816207 PGID:816207 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[]}… (506 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,1325,49\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:44: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05r1-ticket-5-check-1",
          "performer": "claude:plt-t05-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair1-20261022",
            "digest": "sha256:dddeba42f87f96326b1cca5d8441c427551f9b590fd098e1e6a4492f85e7f669",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,pass,1403,49\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-1",
          "command": "bench test --package ./internal/env",
          "exit_code": 0
        },
        {
          "id": "t05r1-ticket-5-check-2",
          "performer": "claude:plt-t05-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair1-20261022",
            "digest": "sha256:5baf36425f84a50f3750e52e09d309fcd468535b5620da440acc17a8b123d87c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,30335,515\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t05r1-ticket-5-check-3",
          "performer": "claude:plt-t05-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair1-20261022",
            "digest": "sha256:06b64c6f29c6916f14b49ed99935ced4ab98fa25ce7671dd7094bc8a096d8f0a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17359,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t05r1-ticket-5-check-4",
          "performer": "claude:plt-t05-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair1-20261022",
            "digest": "sha256:d975a963b678242b85a89c1de75082fd558ba7328d55fc8256ecb2cb180d8c02",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20748,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t05r1-ticket-5-check-5",
          "performer": "claude:plt-t05-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair1-20261022",
            "digest": "sha256:38727b079e10841d2f4145500bd7bc7db218b6a19fe390c7e9da75ae980c09c3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,48551,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-5-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t05r1-ticket-5-check-6",
          "performer": "claude:plt-t05-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair1-20261022",
            "digest": "sha256:48c9c7b75a0884dd2dd15a78a109146e3d11b6c91b78aac0349466261a31d2ad",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4050,35\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-6",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t05r1-ticket-5-check-7",
          "performer": "claude:plt-t05-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair1-20261022",
            "digest": "sha256:13a4d54a11d7d19660826ae994c330667f0fbfe0f064ec6fdf7877c8ee06bb03",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,52,42\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-7",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t05r1-ticket-5-omission",
          "performer": "claude:plt-t05-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair1-20261022",
            "digest": "sha256:5f921cdb05d38aebfdb726d4256182f02f3140149e4d19656c207720ae9595f6",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/env/kit_run.go,swap,failed,1,TestKitRunClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/env,all,passed,49\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,334\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,1392,49\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:44: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-omission",
          "command": "bench test --package ./internal/env",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the private-run descriptor or discard KitTestRun.Close failure. The child-owned file survives while a successful phase result must disappear.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t05-repair1-20261022",
              "digest": "sha256:5f921cdb05d38aebfdb726d4256182f02f3140149e4d19656c207720ae9595f6",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/env/kit_run.go,swap,failed,1,TestKitRunClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/env,all,passed,49\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,334\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,1392,49\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:44: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06-ticket-6-check-1",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:28467ef7bc0a307cdcf6efddbc55a4b16c1dff00b42aff06972f4b1b4fa62cf7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,264,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-1",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-check-2",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:d9dee3a10c3ba69c296c4150be3850165abf510e62ebaae4d569df913b0d20ce",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,29741,515\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-check-3",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:db337d7c274e76982f2f2502a404db0334e4f2418ca61c2688f1dfae69e27e5e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,209,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,61940,209\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-check-4",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:da14e7fe7ef726a5398bfc8505342a805a37e42706b7a82ffb35ea1a95d53708",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17728,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-check-5",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:28bd8580c7c9160b372a6635433db7b9eda8a67b643f624c419b558ef2c5542c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,48732,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-6-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t06-ticket-6-omission",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:8925e5f875e2255a6eeed2a7f86b5f6890ac88dd5aab9080d12b340b7316ef34",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  2,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  3,bit,omit,failed,2,TestCleanOnAnAbsentDirectory,48,1,0,null,yes\n  4,bit,swap,failed,1,TestProtectionEmptyCache,48,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gocache,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2365\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,349,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:164: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/V7UDAT/t/TestProtectionUnreadableProtection815885429/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protecti… (250 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,349,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:164: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/CYRS3I/t/TestProtectionUnreadableProtection1708784181/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,266,48\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanOnAnAbsentDirectory,\"clean_test.go:175: clean = exit 1, want 0; out=\\\"error: cache lock unavailable — open /dev/shm/QITPAB/t/TestCleanOnAnAbsentDirectory1706497144/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionAbsentCache,\"protection_test.go:174: clean of an absent cache = (\\\"error: cache lock unavailable — open /dev/shm/QITPAB/t/TestProtectionAbsentCache542391112/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,273,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionEmptyCache,\"protection_test.go:197: clean of an empty protected cache = (\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-omission",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0,
          "probe": {
            "mutation": "Use only the dropped owner lock, treat unreadable records as empty, or publish green before terminal persistence. The second CLI cleaner and real gate record witnesses must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
              "digest": "sha256:8925e5f875e2255a6eeed2a7f86b5f6890ac88dd5aab9080d12b340b7316ef34",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  2,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  3,bit,omit,failed,2,TestCleanOnAnAbsentDirectory,48,1,0,null,yes\n  4,bit,swap,failed,1,TestProtectionEmptyCache,48,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gocache,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2365\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,349,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:164: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/V7UDAT/t/TestProtectionUnreadableProtection815885429/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protecti… (250 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,349,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:164: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/CYRS3I/t/TestProtectionUnreadableProtection1708784181/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,266,48\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanOnAnAbsentDirectory,\"clean_test.go:175: clean = exit 1, want 0; out=\\\"error: cache lock unavailable — open /dev/shm/QITPAB/t/TestCleanOnAnAbsentDirectory1706497144/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionAbsentCache,\"protection_test.go:174: clean of an absent cache = (\\\"error: cache lock unavailable — open /dev/shm/QITPAB/t/TestProtectionAbsentCache542391112/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,273,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionEmptyCache,\"protection_test.go:197: clean of an empty protected cache = (\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t03r1-ticket-3-check-1",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:21ebd3dc7fa3b092befec9c69bec60e1ed3eac4aeed324522373f42440c96b01",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,24557,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t03r1-ticket-3-check-2",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:7bac8d4542ef0652e9692c488ca08c458665985a8ef3a1ff0b22027f4cc70a6a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15734,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t03r1-ticket-3-check-3",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:31ba018e0b5d31c8dbe8b85ab8547e00d021bdfa6b03674a10f8953df92a813d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,40174,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-3-check-3",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t03r1-ticket-3-omission",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:ce16f385b6ddd503914026215d987ab7b574c3f5f04da9ed3c6a6a94be4bc95b",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGateInterruptOpensWithINTAndExits130,516,1,0,null,yes\n  2,bit,swap,failed,1,TestGateTimeoutOpensWithTERMAndExits124,516,1,0,null,yes\n  3,bit,swap,failed,1,TestGateNormalDrainRemovesReadyDescendant,516,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1164\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,24985,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateInterruptOpensWithINTAndExits130,\"lifetime_test.go:56: interrupted phase opener = \\\"TERM\\\", exit = 130, interrupted = true, want INT, 130, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,24737,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTimeoutOpensWithTERMAndExits124,\"lifetime_test.go:64: timed-out phase opener = \\\"INT\\\", exit = 130, interrupted = true, want TERM, 124, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,25231,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateNormalDrainRemovesReadyDescendant,\"lifetime_test.go:78: descendant 1721618 survived normal completion\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-omission",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace INT with TERM, map timeout to interruption, or omit normal group drain. The actual gate signal and descendant-survival fixtures must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
              "digest": "sha256:ce16f385b6ddd503914026215d987ab7b574c3f5f04da9ed3c6a6a94be4bc95b",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGateInterruptOpensWithINTAndExits130,516,1,0,null,yes\n  2,bit,swap,failed,1,TestGateTimeoutOpensWithTERMAndExits124,516,1,0,null,yes\n  3,bit,swap,failed,1,TestGateNormalDrainRemovesReadyDescendant,516,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1164\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,24985,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateInterruptOpensWithINTAndExits130,\"lifetime_test.go:56: interrupted phase opener = \\\"TERM\\\", exit = 130, interrupted = true, want INT, 130, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,24737,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTimeoutOpensWithTERMAndExits124,\"lifetime_test.go:64: timed-out phase opener = \\\"INT\\\", exit = 130, interrupted = true, want TERM, 124, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,25231,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateNormalDrainRemovesReadyDescendant,\"lifetime_test.go:78: descendant 1721618 survived normal completion\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t04r2-ticket-4-check-1",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:45ae80a426353e5944344727a08e3174e6248c4e9c63fafd05f6cea42904990c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,pass,30810,38\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-1",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "t04r2-ticket-4-check-2",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:28cf37f6d993daace599a71420683d2d95104d18f14ce009613d7f7c33f1cc29",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15732,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t04r2-ticket-4-check-3",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:6f123cc132c9899519f8b26b083239920d6e5083bd3c654e18da491cc4ea19ae",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20464,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-3",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t04r2-ticket-4-check-4",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:bd48c5081fc0a8f72a62efd6ddd7fd8e6e32c3751efd74875a5796a751069442",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,42699,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-4-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t04r2-ticket-4-check-5",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:95600bd866b591c59badf047102ffd3786b703750dc7de06f9448c23e977ef96",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,3981,36\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,48,46\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-5",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0
        },
        {
          "id": "t04r2-ticket-4-check-6",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:6e353c966b752e5879221cb9f89e6885ac32aa2d16d3a88e849962ff0078137f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,261,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-6",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t04r2-ticket-4-omission",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:ab8d9f119580275e8725a0f7689ea2f2a4a932450aa200cf52298061231ea493",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuilderNormalDrain,38,1,0,null,yes\n  2,bit,swap,failed,1,TestBinaryClose,38,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/runbinary,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,760\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,30209,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderNormalDrain,\"lifetime_test.go:26: build private Bench executable: builder process protection protected: <nil>\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,25732,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:98: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-omission",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove normal builder drain or delete the selection before all registered groups are absent. The live descendant or executable-retention assertion must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
              "digest": "sha256:ab8d9f119580275e8725a0f7689ea2f2a4a932450aa200cf52298061231ea493",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuilderNormalDrain,38,1,0,null,yes\n  2,bit,swap,failed,1,TestBinaryClose,38,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/runbinary,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,760\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,30209,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderNormalDrain,\"lifetime_test.go:26: build private Bench executable: builder process protection protected: <nil>\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,25732,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:98: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05r2-ticket-5-check-1",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:5d81744413f68d3264d071c28c48b13fa5cd2edcf901e180a62ea676a4d40b51",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,pass,860,50\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-1",
          "command": "bench test --package ./internal/env",
          "exit_code": 0
        },
        {
          "id": "t05r2-ticket-5-check-2",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:0326aaaef0b1dd2263beabbe73cf1ced16895c7b6c41c1b5a040f7589b00503e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,31815,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t05r2-ticket-5-check-3",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:28cf37f6d993daace599a71420683d2d95104d18f14ce009613d7f7c33f1cc29",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15732,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t05r2-ticket-5-check-4",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:6f123cc132c9899519f8b26b083239920d6e5083bd3c654e18da491cc4ea19ae",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20464,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t05r2-ticket-5-check-5",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:bd48c5081fc0a8f72a62efd6ddd7fd8e6e32c3751efd74875a5796a751069442",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,42699,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-5-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t05r2-ticket-5-check-6",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:38dad289400c473aaa89af357a936b9f3349226220887dcb743da1defd404cf1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4002,36\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-6",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t05r2-ticket-5-check-7",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:6cc6fa26414506739590f9b02fd6eeae97acfcb415275fce93a616c27d4db940",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,47,46\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-7",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t05r2-ticket-5-omission",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:d876aa857b760a9414f248f89a53560dfac464d4f33b22ce11a5a29ce20e0ac7",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestKitRunClose,50,1,0,null,yes\n  2,bit,swap,failed,1,TestKitRunClose,50,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/env,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,882\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,838,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:61: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,879,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:57: user outcome = {Process:{Started:true StartErr:<nil> PID:2321714 PGID:2321714 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[… (508 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-omission",
          "command": "bench test --package ./internal/env",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the private-run descriptor or discard KitTestRun.Close failure. The child-owned file survives while a successful phase result must disappear.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
              "digest": "sha256:d876aa857b760a9414f248f89a53560dfac464d4f33b22ce11a5a29ce20e0ac7",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestKitRunClose,50,1,0,null,yes\n  2,bit,swap,failed,1,TestKitRunClose,50,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/env,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,882\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,838,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:61: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,879,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:57: user outcome = {Process:{Started:true StartErr:<nil> PID:2321714 PGID:2321714 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[… (508 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t07r1-ticket-7-check-1",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:44290b676111e7054994d609d4eff710df16f566c0653f55ff3306b824e3e310",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,369,43\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-1",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "t07r1-ticket-7-check-2",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:0326aaaef0b1dd2263beabbe73cf1ced16895c7b6c41c1b5a040f7589b00503e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,31815,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t07r1-ticket-7-check-3",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:0a5a04c348eec8610aae50050c86689377bceb69da156d26f87951c39e064a31",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,209,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,49100,209\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t07r1-ticket-7-check-4",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:eac5d90d292ed08437b9f70effd5b95d477b498f8a5c2d4f316b46c05d151508",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4002,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-4",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t07r1-ticket-7-check-5",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:28cf37f6d993daace599a71420683d2d95104d18f14ce009613d7f7c33f1cc29",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15732,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t07r1-ticket-7-check-6",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:6f123cc132c9899519f8b26b083239920d6e5083bd3c654e18da491cc4ea19ae",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20464,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t07r1-ticket-7-check-7",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:145338acf7e0629553ec303461cb7ca1eae0e2aa6c98866d1e81c6be6d2ebf39",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,109,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-7",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t07r1-ticket-7-check-8",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:bd48c5081fc0a8f72a62efd6ddd7fd8e6e32c3751efd74875a5796a751069442",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,42699,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-7-check-8",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t07r1-ticket-7-omission",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:7925463e01abeadb7dbdd0aef890ce12f5818e378432907a55b4714a5a3349fd",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,2,TestDirectClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate/prospectiveartifact,all,passed,43\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,556\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,450,43\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestDirectClose,\"lifetime_test.go:98: Close with a live group user = <nil>, want the retention\",1\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOpenRefusesRecoveryWhenRegistrationRemovalFails,\"lifetime_test.go:176: first recovery with a registration failure = nil, want refusal\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-omission",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore owner-death-only sweep, ignore a nested registered group, or discard lane selection release failure. The real Close, second-process sweep, or green-record witness must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
              "digest": "sha256:7925463e01abeadb7dbdd0aef890ce12f5818e378432907a55b4714a5a3349fd",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,2,TestDirectClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate/prospectiveartifact,all,passed,43\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,556\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,450,43\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestDirectClose,\"lifetime_test.go:98: Close with a live group user = <nil>, want the retention\",1\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOpenRefusesRecoveryWhenRegistrationRemovalFails,\"lifetime_test.go:176: first recovery with a registration failure = nil, want refusal\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t07r1b-ticket-7-check-1",
          "performer": "claude:plt-t07-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair1-20261022",
            "digest": "sha256:ee7056d2db0aec3196762cb4f70948f8050ad304c7c0ea9aed84e7fc18a3469b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,359,43\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-1",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "t07r1b-ticket-7-check-2",
          "performer": "claude:plt-t07-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair1-20261022",
            "digest": "sha256:e96c7438df8acfe89a54a102b08be7869d9e578ebc6e67bcb73ea6ac8ec4c74b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,24041,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t07r1b-ticket-7-check-3",
          "performer": "claude:plt-t07-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair1-20261022",
            "digest": "sha256:09ea4714aa735a040b1b901151df904a74c0a9a9c3d89efd8f75014b70094173",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,209,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,47166,209\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t07r1b-ticket-7-check-4",
          "performer": "claude:plt-t07-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair1-20261022",
            "digest": "sha256:2f383b092f4cc40e68638582d5828906e9203d466f6ed3f30a12def84b250dfb",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3847,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-4",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t07r1b-ticket-7-check-5",
          "performer": "claude:plt-t07-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair1-20261022",
            "digest": "sha256:99d32ed2f607baebcfedaa9ab438626afb4827966df8cd35c4b9ce385f6350c7",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14180,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t07r1b-ticket-7-check-6",
          "performer": "claude:plt-t07-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair1-20261022",
            "digest": "sha256:731a76bbfdc5e9219bcf5c91d7e04afd40b7ff6e99d38743db859868c2d94f04",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16655,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t07r1b-ticket-7-check-7",
          "performer": "claude:plt-t07-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair1-20261022",
            "digest": "sha256:0ae9c6b1a47802a7628f3da0e12bcf89199cfe9ec49ca31460612e16a345d6b7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,112,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-7",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t07r1b-ticket-7-check-8",
          "performer": "claude:plt-t07-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair1-20261022",
            "digest": "sha256:c97cc8466117a7e26c25442ac023f5ce9b757e69548693380bc3d430284e72cf",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38958,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-7-check-8",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t07r1b-ticket-7-omission",
          "performer": "claude:plt-t07-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair1-20261022",
            "digest": "sha256:57a43bc5a126f1ae4265c516c0bb723817d9db0f3817dfdd7b69a99ce1b9f91d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,TestDirectClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate/prospectiveartifact,all,passed,43\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,353\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,322,43\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestDirectClose,\"lifetime_test.go:98: Close with a live group user = <nil>, want the retention\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-omission",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore owner-death-only sweep, ignore a nested registered group, or discard lane selection release failure. The real Close, second-process sweep, or green-record witness must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t07-repair1-20261022",
              "digest": "sha256:57a43bc5a126f1ae4265c516c0bb723817d9db0f3817dfdd7b69a99ce1b9f91d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,TestDirectClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate/prospectiveartifact,all,passed,43\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,353\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,322,43\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestDirectClose,\"lifetime_test.go:98: Close with a live group user = <nil>, want the retention\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t04r3-ticket-4-check-1",
          "performer": "claude:plt-t04-repair3-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair3-20261022",
            "digest": "sha256:c0e24989b88fd0831610c061199d6e0a1246fb0088fab68443227b023f0e8a87",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,pass,24553,38\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-1",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "t04r3-ticket-4-check-2",
          "performer": "claude:plt-t04-repair3-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair3-20261022",
            "digest": "sha256:8b827b2e2410de975c9506a2fa7e53f7230f3c0ac9e2007d06c706b81cd52043",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14352,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t04r3-ticket-4-check-3",
          "performer": "claude:plt-t04-repair3-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair3-20261022",
            "digest": "sha256:d2b9b518451b3e9dc42a491808b3b0a70000b211853cea7e9d72bb9150a027ad",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17370,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-3",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t04r3-ticket-4-check-4",
          "performer": "claude:plt-t04-repair3-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair3-20261022",
            "digest": "sha256:81569562027a28fb47e1748303b4f8786f4122e9d52cd4a07d8b71187c5e0769",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,39772,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-4-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t04r3-ticket-4-check-5",
          "performer": "claude:plt-t04-repair3-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair3-20261022",
            "digest": "sha256:b7ee5204a5b56a537948bd6f8328cd399dd112766b3f5c56dc3086fef20ef300",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,3951,36\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,46,46\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-5",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0
        },
        {
          "id": "t04r3-ticket-4-check-6",
          "performer": "claude:plt-t04-repair3-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair3-20261022",
            "digest": "sha256:728760258e5068b73f942eb4e79ae493de5c0c4628e6adff7f3fa97747451253",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,169,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-6",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t04r3-ticket-4-omission",
          "performer": "claude:plt-t04-repair3-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair3-20261022",
            "digest": "sha256:9f230da23b0e8f312580e8f2d1cddb10fcf56618be977041e972b8158c5f1f90",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuilderNormalDrain,38,1,0,null,yes\n  2,bit,omit,failed,1,TestBinaryClose,38,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/runbinary,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,760\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,23805,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderNormalDrain,\"lifetime_test.go:26: build private Bench executable: builder process protection protected: <nil>\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,23938,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:98: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-omission",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove normal builder drain or delete the selection before all registered groups are absent. The live descendant or executable-retention assertion must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t04-repair3-20261022",
              "digest": "sha256:9f230da23b0e8f312580e8f2d1cddb10fcf56618be977041e972b8158c5f1f90",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestBuilderNormalDrain,38,1,0,null,yes\n  2,bit,omit,failed,1,TestBinaryClose,38,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/runbinary,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,760\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,23805,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderNormalDrain,\"lifetime_test.go:26: build private Bench executable: builder process protection protected: <nil>\",2\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,23938,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:98: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05r2b-ticket-5-check-1",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:6903c46a76b4b106418f8b5bcd242af729585d24374aa8b87774cea270b1e630",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,pass,764,50\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-1",
          "command": "bench test --package ./internal/env",
          "exit_code": 0
        },
        {
          "id": "t05r2b-ticket-5-check-2",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:95295742818c2acdc2a03631ff8bd6e8e0bff075cc8808956249c7e12e7c7f8b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,24571,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t05r2b-ticket-5-check-3",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:cd2ac4e1b76139e19d42b7ce319961c3cf57e64e14f5518df9171a4a6c98cbed",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17397,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t05r2b-ticket-5-check-4",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:b075108359df5a08d303276bddff5d739601b92f40c927f0b34da54deda31a3d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,21783,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t05r2b-ticket-5-check-5",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:8eb3a671b2e871966fe13900ff138ce94fabd01a8c677efd88f7595506e298a4",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,48864,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-5-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t05r2b-ticket-5-check-6",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:4bdccf35709571fb9ccdf885934cdcdfb1ef1b133b0380e6c0bdb605768546ff",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4101,36\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-6",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t05r2b-ticket-5-check-7",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:6c9485e44b75122fa5c105a2fc817891433384e29461e0f2ae4a5b5b2193b4f3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,52,46\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-7",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t05r2b-ticket-5-omission",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:25e86a75268986bfbb8b15b8ae85a3c47516085a0d2b82d980ddaebb3ab57afc",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/env/kit_run.go,swap,failed,1,TestKitRunClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/env,all,passed,50\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,490\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,2051,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:57: user outcome = {Process:{Started:true StartErr:<nil> PID:3750727 PGID:3750727 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[… (508 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-omission",
          "command": "bench test --package ./internal/env",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the private-run descriptor or discard KitTestRun.Close failure. The child-owned file survives while a successful phase result must disappear.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
              "digest": "sha256:25e86a75268986bfbb8b15b8ae85a3c47516085a0d2b82d980ddaebb3ab57afc",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/env/kit_run.go,swap,failed,1,TestKitRunClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/env,all,passed,50\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,490\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,2051,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:57: user outcome = {Process:{Started:true StartErr:<nil> PID:3750727 PGID:3750727 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[… (508 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06b-ticket-6-check-1",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:f81648f0f20b386d263bace33f62028486042f717f243fe3d91a5cd49d273ee1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,152,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-1",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t06b-ticket-6-check-2",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:945c2876af88cba35058005c9a4558336f4cccadd2e05619667815ac2943d97c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,32159,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t06b-ticket-6-check-3",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:f1d4767120fc62a311bc1d02efd43dce29a41e907ae0177ca6ea71e92d0da18f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,209,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,60101,209\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t06b-ticket-6-check-4",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:de692295b38165917500ad3034a3ed377a881ecfc4a69bc45dbe0595cf0f25e9",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17648,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06b-ticket-6-check-5",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:f2c7aa7a1bdf157d709b9ee63ed30dd4e6dfb3f1171074753ff63e94b27c31e0",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,56476,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-6-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t06b-ticket-6-omission",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:704ce3051082c1705bd89ea2df769e45ded1650ffd495ae71c508da735801b83",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  2,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  3,bit,omit,failed,2,TestCleanOnAnAbsentDirectory,48,1,0,null,yes\n  4,bit,swap,failed,1,TestProtectionEmptyCache,48,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gocache,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2365\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,243,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:165: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/BOP4LP/t/TestProtectionUnreadableProtection3430648816/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,213,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:165: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/ZWKOSF/t/TestProtectionUnreadableProtection1177041840/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,146,48\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanOnAnAbsentDirectory,\"clean_test.go:175: clean = exit 1, want 0; out=\\\"error: cache lock unavailable — open /dev/shm/ZPNUJ5/t/TestCleanOnAnAbsentDirectory792211444/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionAbsentCache,\"protection_test.go:175: clean of an absent cache = (\\\"error: cache lock unavailable — open /dev/shm/ZPNUJ5/t/TestProtectionAbsentCache1696668615/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,153,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionEmptyCache,\"protection_test.go:198: clean of an empty protected cache = (\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-omission",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0,
          "probe": {
            "mutation": "Use only the dropped owner lock, treat unreadable records as empty, or publish green before terminal persistence. The second CLI cleaner and real gate record witnesses must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
              "digest": "sha256:704ce3051082c1705bd89ea2df769e45ded1650ffd495ae71c508da735801b83",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  2,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  3,bit,omit,failed,2,TestCleanOnAnAbsentDirectory,48,1,0,null,yes\n  4,bit,swap,failed,1,TestProtectionEmptyCache,48,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gocache,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2365\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,243,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:165: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/BOP4LP/t/TestProtectionUnreadableProtection3430648816/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,213,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:165: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/ZWKOSF/t/TestProtectionUnreadableProtection1177041840/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,146,48\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanOnAnAbsentDirectory,\"clean_test.go:175: clean = exit 1, want 0; out=\\\"error: cache lock unavailable — open /dev/shm/ZPNUJ5/t/TestCleanOnAnAbsentDirectory792211444/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionAbsentCache,\"protection_test.go:175: clean of an absent cache = (\\\"error: cache lock unavailable — open /dev/shm/ZPNUJ5/t/TestProtectionAbsentCache1696668615/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,153,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionEmptyCache,\"protection_test.go:198: clean of an empty protected cache = (\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t14-ticket-14-check-1",
          "performer": "claude:plt-t14-author-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t14-author-20261023",
            "digest": "sha256:0e72fe3549ada50bbe8543ed1ec5db0b25f6a5514c2015be9c4854338bda1a45",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4078,38\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,62,54\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-14-check-1",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0
        },
        {
          "id": "t14-ticket-14-check-2",
          "performer": "claude:plt-t14-author-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t14-author-20261023",
            "digest": "sha256:4b1e7754afc02007e2b96848cae0eb8fe9de0c08119f6574e20eb2ad22031ba6",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,24689,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-2",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t14-ticket-14-check-3",
          "performer": "claude:plt-t14-author-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t14-author-20261023",
            "digest": "sha256:7137835fc89c1c934e64c0c5780dcf27f6e09f78a6bd5c80ab477a653ad56567",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,122,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-3",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t14-ticket-14-check-4",
          "performer": "claude:plt-t14-author-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t14-author-20261023",
            "digest": "sha256:bec982b817470ebd90f6bb54813410511fe1eb96b2da5ec539afbb371c7b7f85",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,46156,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-14-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t14-ticket-14-check-5",
          "performer": "claude:plt-t14-author-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t14-author-20261023",
            "digest": "sha256:6bc4359381c8c7b7c6ef050eafac26b7a4153efa89470aa489095ae995cd425d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16701,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t14-ticket-14-check-6",
          "performer": "claude:plt-t14-author-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t14-author-20261023",
            "digest": "sha256:2d2f3e0971a22c706da1ce5a6a813da420f7744317216cd7c13c0e5500eef3a7",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,65916,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-6",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t14-ticket-14-check-7",
          "performer": "claude:plt-t14-author-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t14-author-20261023",
            "digest": "sha256:1b06a693f346cb301c419fdbc2206ae3bff010159bd69939e9ec35aac61491ea",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,31738,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-7",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t14-ticket-14-omission",
          "performer": "claude:plt-t14-author-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t14-author-20261023",
            "digest": "sha256:d78b2796fce1afab7a2be45a3e31359b419052e203a975a5426b9ff26b6928ea",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/startgate.go,swap,failed,1,TestStartGateActorDeathBeforeOpen,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/...,all,passed,92\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,856\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,4055,38\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,60,54\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestStartGateActorDeathBeforeOpen,\"startgate_test.go:72: the program ran after its actor died before the open signal: <nil>\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-14-omission",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Execute the program at end of file, or grade a pending user with an absent actor as protected. The ActorDeathBeforeOpen or PendingActorAbsent witness must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t14-author-20261023",
              "digest": "sha256:d78b2796fce1afab7a2be45a3e31359b419052e203a975a5426b9ff26b6928ea",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/startgate.go,swap,failed,1,TestStartGateActorDeathBeforeOpen,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/...,all,passed,92\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,856\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,4055,38\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,60,54\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestStartGateActorDeathBeforeOpen,\"startgate_test.go:72: the program ran after its actor died before the open signal: <nil>\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
            }
          }
        },
        {
          "id": "t03r1b-ticket-3-check-1",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:8a8124a916ca3b5889c1e65be58b554dbde44aeb2f8f0996439491fe8b9991b4",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,31624,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t03r1b-ticket-3-check-2",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:2e9e206f8b374cf0ef9032fbb22c0b9237d0e0155ba2ba47bb0ead6bdb968d7e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,19666,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t03r1b-ticket-3-check-3",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:a9d8680aae6666533b2406a36b5132c35035398c1cae84c86d5366693a804458",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,47490,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-3-check-3",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t03r1b-ticket-3-omission",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:f9cb6968d4b846bf8f81b2f74b2f7b455a814ed98e4b21c62f95cacd71c66d1e",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGateInterruptOpensWithINTAndExits130,516,1,0,null,yes\n  2,bit,swap,failed,1,TestGateTimeoutOpensWithTERMAndExits124,516,1,0,null,yes\n  3,bit,swap,failed,1,TestGateNormalDrainRemovesReadyDescendant,516,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1164\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,32212,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateInterruptOpensWithINTAndExits130,\"lifetime_test.go:56: interrupted phase opener = \\\"TERM\\\", exit = 130, interrupted = true, want INT, 130, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,32329,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTimeoutOpensWithTERMAndExits124,\"lifetime_test.go:64: timed-out phase opener = \\\"INT\\\", exit = 130, interrupted = true, want TERM, 124, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,32338,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateNormalDrainRemovesReadyDescendant,\"lifetime_test.go:78: descendant 3605668 survived normal completion\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-omission",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace INT with TERM, map timeout to interruption, or omit normal group drain. The actual gate signal and descendant-survival fixtures must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
              "digest": "sha256:f9cb6968d4b846bf8f81b2f74b2f7b455a814ed98e4b21c62f95cacd71c66d1e",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGateInterruptOpensWithINTAndExits130,516,1,0,null,yes\n  2,bit,swap,failed,1,TestGateTimeoutOpensWithTERMAndExits124,516,1,0,null,yes\n  3,bit,swap,failed,1,TestGateNormalDrainRemovesReadyDescendant,516,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1164\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,32212,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateInterruptOpensWithINTAndExits130,\"lifetime_test.go:56: interrupted phase opener = \\\"TERM\\\", exit = 130, interrupted = true, want INT, 130, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,32329,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTimeoutOpensWithTERMAndExits124,\"lifetime_test.go:64: timed-out phase opener = \\\"INT\\\", exit = 130, interrupted = true, want TERM, 124, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,32338,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateNormalDrainRemovesReadyDescendant,\"lifetime_test.go:78: descendant 3605668 survived normal completion\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t04r4-ticket-4-check-1",
          "performer": "claude:plt-t04-repair4-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
            "digest": "sha256:4dbb641468e18f6d568feaf7ead13e6b731de4bc760b1cf5988a2ad4667ac351",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,pass,29043,38\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-1",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "t04r4-ticket-4-check-2",
          "performer": "claude:plt-t04-repair4-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
            "digest": "sha256:7619f5ae9daf98b635cd2a129d503653a17f4cef8664895c77f169a0b09c5e09",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17340,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t04r4-ticket-4-check-3",
          "performer": "claude:plt-t04-repair4-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
            "digest": "sha256:8e1af329969f2adfec49091c71c76f54bef12854797c994af8d2742199fb18e3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,25498,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-3",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t04r4-ticket-4-check-4",
          "performer": "claude:plt-t04-repair4-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
            "digest": "sha256:55e6c952d45e398b71ecdfca50a4e7d5897b3cb298e61b5cdffd06e1e8e5e51e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,46878,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-4-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t04r4-ticket-4-check-5",
          "performer": "claude:plt-t04-repair4-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
            "digest": "sha256:a6caf0891309047c7bebb6b215e8c9d081be6a49538c4df99b718cad9ffeb030",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4084,38\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,62,54\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-4-check-5",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0
        },
        {
          "id": "t04r4-ticket-4-check-6",
          "performer": "claude:plt-t04-repair4-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
            "digest": "sha256:f669b5cded00ad10cc66df7258efc3533af1b60551134e7bbbc055477651630e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,267,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-6",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t04r4-ticket-4-check-7",
          "performer": "claude:plt-t04-repair4-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
            "digest": "sha256:feac0721a75f57b67d28ecf831dcf1d0f0d99274d28d9dd02f71510b91439c7c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,pass,1310,50\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-7",
          "command": "bench test --package ./internal/env",
          "exit_code": 0
        },
        {
          "id": "t04r4-ticket-4-check-8",
          "performer": "claude:plt-t04-repair4-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
            "digest": "sha256:f0353162df1dc346c0d829af6c0ae52ccb537b0f7f8ecf36f27ab4ceed64f074",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,435,43\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-8",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "t04r4-ticket-4-check-9",
          "performer": "claude:plt-t04-repair4-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
            "digest": "sha256:428d823d7210d57715991bdead6736147efd309dc5fd817613b22ba029cae251",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,32160,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-9",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t04r4-ticket-4-omission",
          "performer": "claude:plt-t04-repair4-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
            "digest": "sha256:27cbb7efd75239299b17d0a678be63f069c286a312bc730713e03ece95cff7fc",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/runbinary/runbinary.go,swap,failed,1,TestBinaryClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,TestBinaryClose$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,344\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,8,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:98: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-omission",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove normal builder drain or delete the selection before all registered groups are absent. The live descendant or executable-retention assertion must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t04-repair4-20261023",
              "digest": "sha256:27cbb7efd75239299b17d0a678be63f069c286a312bc730713e03ece95cff7fc",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/runbinary/runbinary.go,swap,failed,1,TestBinaryClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,TestBinaryClose$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,344\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,8,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:98: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05r2c-ticket-5-check-1",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:feab895936e70039687dc37bcfae400641b070ef4124ca6549150e484cd8d45a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,pass,1371,50\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-1",
          "command": "bench test --package ./internal/env",
          "exit_code": 0
        },
        {
          "id": "t05r2c-ticket-5-check-2",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:6199ed56f155c5e3d422e0e43e14e218b14ec748bccc3cd491913de45177b19c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,32169,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t05r2c-ticket-5-check-3",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:094d26a3416780afe1ca1afe151b30f1f4ab27ce47d5eeafe22e9a68dda1881f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16738,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t05r2c-ticket-5-check-4",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:795045602850925bbe0e72e4f37f5898d1780b9fd953bcbc9a73785501581e9b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,25044,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t05r2c-ticket-5-check-5",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:0ebe4fc857b60380f08609a4094080d9f126ea293b615d8db8f4914a4e46e9f7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,46548,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-5-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t05r2c-ticket-5-check-6",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:9c392b259b84ddebbb4d5d3c4aadfb15fb497bdec7eabc4760cda1e8cda11e8d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4056,38\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-6",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t05r2c-ticket-5-check-7",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:89c2ea40a590acc49d1034126547b6fd6ff91966495bdf67a998fcda1c389c3a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,61,54\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-5-check-7",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t05r2c-ticket-5-omission",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:35565121b39f3ab5cb6e8b9432ccb2c257ef77827e8d0c90ea687ce2592d7d35",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/env/kit_run.go,swap,failed,1,TestKitRunClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/env,all,passed,50\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,490\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,1437,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:57: user outcome = {Process:{Started:true StartErr:<nil> PID:3728293 PGID:3728293 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[… (508 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-omission",
          "command": "bench test --package ./internal/env",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the private-run descriptor or discard KitTestRun.Close failure. The child-owned file survives while a successful phase result must disappear.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
              "digest": "sha256:35565121b39f3ab5cb6e8b9432ccb2c257ef77827e8d0c90ea687ce2592d7d35",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/env/kit_run.go,swap,failed,1,TestKitRunClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/env,all,passed,50\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,490\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,1437,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:57: user outcome = {Process:{Started:true StartErr:<nil> PID:3728293 PGID:3728293 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[… (508 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06c-ticket-6-check-1",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:dbf3cbe45e9c11fa661c4a7f0a6db6c4a2646e5088b860eba3e9e6f6c1ec8fc1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,260,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-1",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t06c-ticket-6-check-2",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:94ce6911cfadd82e8f263410d0abec936dd9c11919ac8127729dcaac28c72266",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,32875,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t06c-ticket-6-check-3",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:82fb3920daaea250c69fc7d5ad713ee4d0ffd5bb88d40cc6337c8f71612280f1",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,62114,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t06c-ticket-6-check-4",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:7e894027ffa513888ae2d6efabffcab2e010b170bc228fd18b4a4d504765c212",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17787,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06c-ticket-6-check-5",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:72b8e95d54179c941085b8aa6f0547be671a2b1de274e5ce261b3e4e60e5c48e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,47567,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-6-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t06c-ticket-6-omission",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:61f08e35dbaf23bb46cb3461a36690050ef32b2f4697b09b1685a2d091263611",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  2,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  3,bit,omit,failed,2,TestCleanOnAnAbsentDirectory,48,1,0,null,yes\n  4,bit,swap,failed,1,TestProtectionEmptyCache,48,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gocache,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2366\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,334,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/CNHVSR/t/TestProtectionUnreadableProtection3321946623/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,324,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/QTMVDU/t/TestProtectionUnreadableProtection4214816891/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,240,48\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanOnAnAbsentDirectory,\"clean_test.go:175: clean = exit 1, want 0; out=\\\"error: cache lock unavailable — open /dev/shm/AJFLBH/t/TestCleanOnAnAbsentDirectory2392132198/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionAbsentCache,\"protection_test.go:177: clean of an absent cache = (\\\"error: cache lock unavailable — open /dev/shm/AJFLBH/t/TestProtectionAbsentCache1139134954/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,269,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionEmptyCache,\"protection_test.go:200: clean of an empty protected cache = (\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-omission",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0,
          "probe": {
            "mutation": "Use only the dropped owner lock, treat unreadable records as empty, or publish green before terminal persistence. The second CLI cleaner and real gate record witnesses must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
              "digest": "sha256:61f08e35dbaf23bb46cb3461a36690050ef32b2f4697b09b1685a2d091263611",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  2,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  3,bit,omit,failed,2,TestCleanOnAnAbsentDirectory,48,1,0,null,yes\n  4,bit,swap,failed,1,TestProtectionEmptyCache,48,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gocache,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2366\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,334,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/CNHVSR/t/TestProtectionUnreadableProtection3321946623/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,324,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/QTMVDU/t/TestProtectionUnreadableProtection4214816891/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,240,48\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanOnAnAbsentDirectory,\"clean_test.go:175: clean = exit 1, want 0; out=\\\"error: cache lock unavailable — open /dev/shm/AJFLBH/t/TestCleanOnAnAbsentDirectory2392132198/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionAbsentCache,\"protection_test.go:177: clean of an absent cache = (\\\"error: cache lock unavailable — open /dev/shm/AJFLBH/t/TestProtectionAbsentCache1139134954/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,269,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionEmptyCache,\"protection_test.go:200: clean of an empty protected cache = (\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t07r2-ticket-7-check-1",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:abe194412b6d08ed71b99126dff554ebd075db0f230b253f5499ea2d3469eb7f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,438,43\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-1",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "t07r2-ticket-7-check-2",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:f84796d8d7ab2630014e3a21afd9491f617155899ff75e97bb1c1dae38c137e1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,33510,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t07r2-ticket-7-check-3",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:6a412ef230b30bdaea0030a0cc2e8d31259fad9e56dd1d98ffb6d1483e0c65e1",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,58438,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t07r2-ticket-7-check-4",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:0a6b5a174e9aed7766e3751bb38cfb3074f79ef656492416d88f16fffc4835df",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5069,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-4",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t07r2-ticket-7-check-5",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:3e0e074b666185f8546cd5ac87ee9445766e2c36eb299c771c308c37494794e5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17523,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t07r2-ticket-7-check-6",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:0e8a9d8a659e4dc624a0fe474d33c5d7e61941699f6c930cbb4b667c637c63be",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26160,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t07r2-ticket-7-check-7",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:e908087131657cb3e3a2a71345c9e013d71ded28a37a477bde48c0192e089964",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,126,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-7",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t07r2-ticket-7-check-8",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:323724090e8b32c5f0f10b98281f13f6e5d0c48315f103f2d8bae8b44e435677",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,47723,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-7-check-8",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t07r2-ticket-7-omission",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:db4a1d741f47078030fbb750684402eb898bec40f365ac5d005fa08b7f8cdb57",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/lane.go,swap,failed,1,TestLaneSelectionCloseFailure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,516\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,325\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,33321,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestLaneSelectionCloseFailure,\"lifetime_resources_test.go:264: lane with a nested user = <nil>; output:\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-omission",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore owner-death-only sweep, ignore a nested registered group, or discard lane selection release failure. The real Close, second-process sweep, or green-record witness must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
              "digest": "sha256:db4a1d741f47078030fbb750684402eb898bec40f365ac5d005fa08b7f8cdb57",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/lane.go,swap,failed,1,TestLaneSelectionCloseFailure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,516\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,325\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,33321,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestLaneSelectionCloseFailure,\"lifetime_resources_test.go:264: lane with a nested user = <nil>; output:\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t14r1-ticket-14-check-1",
          "performer": "claude:plt-t14-repair1-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair1-20261023",
            "digest": "sha256:137c3d1ab2adeb2e038f69d9cb3251d95a79753e8213e8826a5233ccef7fb1dd",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4169,39\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,56,59\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-14-check-1",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0
        },
        {
          "id": "t14r1-ticket-14-check-2",
          "performer": "claude:plt-t14-repair1-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair1-20261023",
            "digest": "sha256:4ef524edf7478cb257e3535b6f4c8e16325484d4e205d3533eaf9c01e9d19412",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26392,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-2",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t14r1-ticket-14-check-3",
          "performer": "claude:plt-t14-repair1-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair1-20261023",
            "digest": "sha256:0ae9c6b1a47802a7628f3da0e12bcf89199cfe9ec49ca31460612e16a345d6b7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,112,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-3",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t14r1-ticket-14-check-4",
          "performer": "claude:plt-t14-repair1-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair1-20261023",
            "digest": "sha256:7affd51c69ac90b0437cfb248f705d60413cd9cb1a701b9b43b1e065e250f3a1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,42378,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-14-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t14r1-ticket-14-check-5",
          "performer": "claude:plt-t14-repair1-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair1-20261023",
            "digest": "sha256:b241e790807c918534a7e630b92f3eb8b416b66d0fb0535b72c0965217844f66",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14946,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t14r1-ticket-14-check-6",
          "performer": "claude:plt-t14-repair1-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair1-20261023",
            "digest": "sha256:3b27a5700c6b9703831adcebcf73180409bce1dd685de2ef0311ddd59a15813c",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,49098,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-6",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t14r1-ticket-14-check-7",
          "performer": "claude:plt-t14-repair1-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair1-20261023",
            "digest": "sha256:40724d0e62a30a0877a08fe86bf8266f916a02198ad11b8e351677e88550d014",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,28483,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-7",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t14r1-ticket-14-omission",
          "performer": "claude:plt-t14-repair1-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair1-20261023",
            "digest": "sha256:180a091806af52f9fb3d64cf84e787fd84a0e5aeef21c1aa2472d2e04aff9d0e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/startgate.go,swap,failed,1,TestStartGateActorDeathBeforeOpen,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/...,all,passed,98\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,856\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,4104,39\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,58,59\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestStartGateActorDeathBeforeOpen,\"startgate_test.go:74: the program ran after its actor died before the open signal: <nil>\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-14-omission",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Execute the program at end of file, or grade a pending user with an absent actor as protected. The ActorDeathBeforeOpen or PendingActorAbsent witness must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t14-repair1-20261023",
              "digest": "sha256:180a091806af52f9fb3d64cf84e787fd84a0e5aeef21c1aa2472d2e04aff9d0e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/startgate.go,swap,failed,1,TestStartGateActorDeathBeforeOpen,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/...,all,passed,98\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,856\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,4104,39\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,58,59\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestStartGateActorDeathBeforeOpen,\"startgate_test.go:74: the program ran after its actor died before the open signal: <nil>\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
            }
          }
        },
        {
          "id": "t04r5-ticket-4-check-1",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:f4b349666edb180586ad7d2ddfb532708e5f0b89e8f0b9d0411b81aab3488e76",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,pass,25179,38\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-1",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "t04r5-ticket-4-check-2",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:4d00f0b0a4bb7f0f19940f0080c606891c230ba614a5236a5711f8185643021c",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15018,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t04r5-ticket-4-check-3",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:541ab9d68509dfbff351b042c4864e85a053229ef29a0672d880b41a6d31ccde",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,21261,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-3",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t04r5-ticket-4-check-4",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:5308afb63cf570123e9fb46bffaa07b0f119f5face59e05aee78bbda19094192",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,41549,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-4-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t04r5-ticket-4-check-5",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:7aeff2eb482fdf3b73196083d92f7188e88722a44dffebbaa708fbe66ea242cb",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4125,39\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,56,59\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-4-check-5",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0
        },
        {
          "id": "t04r5-ticket-4-check-6",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:77490f9b9d81d96e9ce62dd51ad93e8d596e24d3b67a3a4224ef42101032f42a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,178,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-6",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t04r5-ticket-4-check-7",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:0ba901ce4a8ba0037cc754088bf70daa43bc4037730c67d49fa4923dda375388",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,pass,816,50\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-7",
          "command": "bench test --package ./internal/env",
          "exit_code": 0
        },
        {
          "id": "t04r5-ticket-4-check-8",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:29a42106ce3ffc2b2f1d6cba1e62dd40daf747e5223b189f30d1c0a437bc06a7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,835,43\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-8",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "t04r5-ticket-4-check-9",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:3c4c5c37f3337cc2b7cb65b4e60f8c85818744a1a2eeba2aaf8f2e19ecb7500c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,28268,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-9",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t04r5-ticket-4-check-10",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:dc9d57e76e68753a21b12cc9bfc9fc7af1e2ee78cb72f9422e69fd4f9983d5e6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,52783,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-10",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t04r5-ticket-4-omission",
          "performer": "claude:plt-t04-repair5-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
            "digest": "sha256:c207e72ba35bc266201b6978c4849ea367e2576380dfe252d7905b93c5764ef5",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/runbinary/runbinary.go,swap,failed,1,TestBinaryClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,all,passed,38\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,349\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,25714,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:98: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-omission",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove normal builder drain or delete the selection before all registered groups are absent. The live descendant or executable-retention assertion must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t04-repair5-20261023",
              "digest": "sha256:c207e72ba35bc266201b6978c4849ea367e2576380dfe252d7905b93c5764ef5",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/runbinary/runbinary.go,swap,failed,1,TestBinaryClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,all,passed,38\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,349\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,25714,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:98: Close with a live registered group = <nil>, want a protected retention\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t03r1c-ticket-3-check-1",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:0da55d144c6cf364f1645b102bc046eab2e606f6e819c0caf5dcd52e7597becd",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,28677,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t03r1c-ticket-3-check-2",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:96ed7caab78e52b13676fbd821712d046c5a17e063189bf46fd5cbbd252c3e19",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15508,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t03r1c-ticket-3-check-3",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:df6abdefe5ddc551400fb7d0cd184d5dd6b28820964bf3e21ad683049cdf0db9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,44590,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-3-check-3",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t03r1c-ticket-3-omission",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:297b254d2b359472e298eabebb8efddcc04798b45419b895c0e08dcb15120ad4",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGateInterruptOpensWithINTAndExits130,516,1,0,null,yes\n  2,bit,swap,failed,1,TestGateTimeoutOpensWithTERMAndExits124,516,1,0,null,yes\n  3,bit,swap,failed,1,TestGateNormalDrainRemovesReadyDescendant,516,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1163\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,33640,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateInterruptOpensWithINTAndExits130,\"lifetime_test.go:56: interrupted phase opener = \\\"TERM\\\", exit = 130, interrupted = true, want INT, 130, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,33848,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTimeoutOpensWithTERMAndExits124,\"lifetime_test.go:64: timed-out phase opener = \\\"INT\\\", exit = 130, interrupted = true, want TERM, 124, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,32871,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateNormalDrainRemovesReadyDescendant,\"lifetime_test.go:78: descendant 523678 survived normal completion\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-omission",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace INT with TERM, map timeout to interruption, or omit normal group drain. The actual gate signal and descendant-survival fixtures must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
              "digest": "sha256:297b254d2b359472e298eabebb8efddcc04798b45419b895c0e08dcb15120ad4",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGateInterruptOpensWithINTAndExits130,516,1,0,null,yes\n  2,bit,swap,failed,1,TestGateTimeoutOpensWithTERMAndExits124,516,1,0,null,yes\n  3,bit,swap,failed,1,TestGateNormalDrainRemovesReadyDescendant,516,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1163\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,33640,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateInterruptOpensWithINTAndExits130,\"lifetime_test.go:56: interrupted phase opener = \\\"TERM\\\", exit = 130, interrupted = true, want INT, 130, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,33848,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTimeoutOpensWithTERMAndExits124,\"lifetime_test.go:64: timed-out phase opener = \\\"INT\\\", exit = 130, interrupted = true, want TERM, 124, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,32871,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateNormalDrainRemovesReadyDescendant,\"lifetime_test.go:78: descendant 523678 survived normal completion\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05r2d-ticket-5-check-1",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:fff76940a30801b7045a7c26c605ce7dbfd8b60634d828df409dc5abb010b700",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,pass,1485,50\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-1",
          "command": "bench test --package ./internal/env",
          "exit_code": 0
        },
        {
          "id": "t05r2d-ticket-5-check-2",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:b04237dfb26730ee2647ad498bb5fa7b5a80faa2ff3ee62f4e3116cd19e96cc6",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,33555,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t05r2d-ticket-5-check-3",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:08cd0be6f56473dde558cb62ed5a64c51912c7cd09cbb46dd162e0cf579477b4",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17401,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t05r2d-ticket-5-check-4",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:47725f3ceccbd04e47538e677bf7586f3603a20120ff55cf5ac3291ae31d45c2",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26301,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t05r2d-ticket-5-check-5",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:a341f27ee346518f330d4d0ce2c3040bf0cfc7da19c79f7ee48e586f6c31cc29",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,47423,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-5-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t05r2d-ticket-5-check-6",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:1bba328730feef18a45505d28e5c9455debbb442347089db66cab82436fc8a57",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4161,39\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-6",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t05r2d-ticket-5-check-7",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:6b8388998b5143690cdf1b75ab4f1fd045e262ca97656632658d288383f526c3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,58,59\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-5-check-7",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t05r2d-ticket-5-omission",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:d72df514f3faa1d06d7cfef2a71aa68ca1b6e896d4b4ec2bf946490f82fcda76",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/env/kit_run.go,swap,failed,1,TestKitRunClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/env,all,passed,50\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,490\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,1440,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:57: user outcome = {Process:{Started:true StartErr:<nil> PID:581443 PGID:581443 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[]}… (506 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-omission",
          "command": "bench test --package ./internal/env",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the private-run descriptor or discard KitTestRun.Close failure. The child-owned file survives while a successful phase result must disappear.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
              "digest": "sha256:d72df514f3faa1d06d7cfef2a71aa68ca1b6e896d4b4ec2bf946490f82fcda76",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/env/kit_run.go,swap,failed,1,TestKitRunClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/env,all,passed,50\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,490\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,1440,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:57: user outcome = {Process:{Started:true StartErr:<nil> PID:581443 PGID:581443 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[]}… (506 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06d-ticket-6-check-1",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:97b33321aefdba09ad3b7478545d36d8ba9408c190b2301817bedd1bf843d112",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,273,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-1",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t06d-ticket-6-check-2",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:98af5af75c38f7d3779bfc7fb7332db5cc4c7ec5a60292a3a0805278b29f0b45",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,32905,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t06d-ticket-6-check-3",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:97e9856bc566736aa778439b2c214d1002287febd928f1ecd5d8be1847588560",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,57777,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t06d-ticket-6-check-4",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:8327131e780d9f2091491f7a8b06977644082c7815d1f2db4c2a44db4cdac30e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17406,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06d-ticket-6-check-5",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:5a13950ff7378f2d74ef1f6ca10e280398917fd47d65db5df79b837dd06e3a8c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,47328,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-6-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t06d-ticket-6-omission",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:6c929a99c0792351193752564c5cf8cbd50f6b8d61fb820b65c6c5e5488b4ff7",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  2,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  3,bit,omit,failed,2,TestCleanOnAnAbsentDirectory,48,1,0,null,yes\n  4,bit,swap,failed,3,TestCleanRemovesTheShardsAndKeepsTheRest,48,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gocache,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2819\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,388,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/KCJALZ/t/TestProtectionUnreadableProtection630943995/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protecti… (250 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,398,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/JOOYBW/t/TestProtectionUnreadableProtection4276480310/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,311,48\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanOnAnAbsentDirectory,\"clean_test.go:175: clean = exit 1, want 0; out=\\\"error: cache lock unavailable — open /dev/shm/YH22VG/t/TestCleanOnAnAbsentDirectory1148633245/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionAbsentCache,\"protection_test.go:177: clean of an absent cache = (\\\"error: cache lock unavailable — open /dev/shm/YH22VG/t/TestProtectionAbsentCache996967024/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,261,48\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanRemovesTheShardsAndKeepsTheRest,\"clean_test.go:143: clean = exit 1, want 0; out=\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestCleanReportsThePreRemovalMeasurement,\"clean_test.go:159: clean = exit 1, want 0; out=\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionEmptyCache,\"protection_test.go:200: clean of an empty protected cache = (\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-omission",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0,
          "probe": {
            "mutation": "Use only the dropped owner lock, treat unreadable records as empty, or publish green before terminal persistence. The second CLI cleaner and real gate record witnesses must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
              "digest": "sha256:6c929a99c0792351193752564c5cf8cbd50f6b8d61fb820b65c6c5e5488b4ff7",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  2,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  3,bit,omit,failed,2,TestCleanOnAnAbsentDirectory,48,1,0,null,yes\n  4,bit,swap,failed,3,TestCleanRemovesTheShardsAndKeepsTheRest,48,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gocache,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2819\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,388,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/KCJALZ/t/TestProtectionUnreadableProtection630943995/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protecti… (250 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,398,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/JOOYBW/t/TestProtectionUnreadableProtection4276480310/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,311,48\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanOnAnAbsentDirectory,\"clean_test.go:175: clean = exit 1, want 0; out=\\\"error: cache lock unavailable — open /dev/shm/YH22VG/t/TestCleanOnAnAbsentDirectory1148633245/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionAbsentCache,\"protection_test.go:177: clean of an absent cache = (\\\"error: cache lock unavailable — open /dev/shm/YH22VG/t/TestProtectionAbsentCache996967024/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,261,48\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanRemovesTheShardsAndKeepsTheRest,\"clean_test.go:143: clean = exit 1, want 0; out=\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestCleanReportsThePreRemovalMeasurement,\"clean_test.go:159: clean = exit 1, want 0; out=\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionEmptyCache,\"protection_test.go:200: clean of an empty protected cache = (\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t07r2b-ticket-7-check-1",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:913c0ea58223620a96772c8e0399469a01ed982a11a83fca1f513677f20ec0e8",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,971,43\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-1",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "t07r2b-ticket-7-check-2",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:2ec2803f01cd9a397264edd9b9b6644f4e5a5c7852aadf7c0238595abf259123",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,33894,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t07r2b-ticket-7-check-3",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:b659dedcf084c098a232919ac781f2b698e294b6a112095041f0a2ecc927ce1f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,60408,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t07r2b-ticket-7-check-4",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:a3928179a0930f912a1fd6c42a14c0b41eb9df51a16a62587bc72229f6b2d31f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5145,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-4",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t07r2b-ticket-7-check-5",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:e9f3536142d79521912a3eefc990466dc7aad6b4496d71a4442ba28a06d671f4",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17365,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t07r2b-ticket-7-check-6",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:d4ce7be3aa76038b584cf2d0125f0300355bfd559db4f8d4b286957370dc1cf0",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,25897,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t07r2b-ticket-7-check-7",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:86cf992eeb857beb65f47d74bfd15a9cc042c6aa665c971b252cd3531491c9d6",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,127,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-7",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t07r2b-ticket-7-check-8",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:297505d74f7336bae4ba90226832a7ca1c95d2da0bad949c997324cf4b4f85a9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,46689,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-7-check-8",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t07r2b-ticket-7-omission",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:1e95572891b1686910d9f807716ddf1cccdea929dd61bcff08a731fcab323dbf",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/lane.go,swap,failed,1,TestLaneSelectionCloseFailure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,516\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,325\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,33581,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestLaneSelectionCloseFailure,\"lifetime_resources_test.go:264: lane with a nested user = <nil>; output:\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-omission",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore owner-death-only sweep, ignore a nested registered group, or discard lane selection release failure. The real Close, second-process sweep, or green-record witness must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
              "digest": "sha256:1e95572891b1686910d9f807716ddf1cccdea929dd61bcff08a731fcab323dbf",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/lane.go,swap,failed,1,TestLaneSelectionCloseFailure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,516\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,325\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,33581,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestLaneSelectionCloseFailure,\"lifetime_resources_test.go:264: lane with a nested user = <nil>; output:\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t14r2-ticket-14-check-1",
          "performer": "claude:plt-t14-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair2-20261023",
            "digest": "sha256:7d06cc342ba3acbb27c0ddb0e426376d73d553999c367c0f84344c9fdded5fb7",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4066,39\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,52,59\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-14-check-1",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0
        },
        {
          "id": "t14r2-ticket-14-check-2",
          "performer": "claude:plt-t14-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair2-20261023",
            "digest": "sha256:5a6816f577d90111e4f6cffd4bccdcf72d2efa200e6bba170e56577a76f32f92",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19984,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-2",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t14r2-ticket-14-check-3",
          "performer": "claude:plt-t14-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair2-20261023",
            "digest": "sha256:54965290673920ffe0d877794aab2c5c674453ba64f4fd6793b2c6ca587f6b7e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,107,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-3",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t14r2-ticket-14-check-4",
          "performer": "claude:plt-t14-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair2-20261023",
            "digest": "sha256:ac6ed2391f030914936bec1ddad5dcbbc0dc661324ad63fc13676e796ad2cc97",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38428,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-14-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t14r2-ticket-14-check-5",
          "performer": "claude:plt-t14-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair2-20261023",
            "digest": "sha256:4d95354dd2e08857b052e315be12d235fcb0cf76b3d197c1fc9f45b03e0a93a7",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13668,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t14r2-ticket-14-check-6",
          "performer": "claude:plt-t14-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair2-20261023",
            "digest": "sha256:d2bc763d79f5bed60e5911791bdd70b47330c6200cf33fed2b6030802cfbdc64",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,43787,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-6",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t14r2-ticket-14-check-7",
          "performer": "claude:plt-t14-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair2-20261023",
            "digest": "sha256:0c324e1b1a29b3c1aaf25d26e67519f525197c03e6c7053ba5bdabe9fa8fcef8",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,27046,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-14-check-7",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t14r2-ticket-14-omission",
          "performer": "claude:plt-t14-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t14-repair2-20261023",
            "digest": "sha256:8a3243ace72851863f0fe50d81c5a062280bad4448e4567e450fa9aef233bebc",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/startgate.go,swap,failed,1,TestStartGateActorDeathBeforeOpen,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/...,all,passed,98\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,856\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,4070,39\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,54,59\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestStartGateActorDeathBeforeOpen,\"startgate_test.go:74: the program ran after its actor died before the open signal: <nil>\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-14-omission",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Execute the program at end of file, or grade a pending user with an absent actor as protected. The ActorDeathBeforeOpen or PendingActorAbsent witness must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t14-repair2-20261023",
              "digest": "sha256:8a3243ace72851863f0fe50d81c5a062280bad4448e4567e450fa9aef233bebc",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/startgate.go,swap,failed,1,TestStartGateActorDeathBeforeOpen,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/...,all,passed,98\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,856\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,4070,39\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,54,59\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestStartGateActorDeathBeforeOpen,\"startgate_test.go:74: the program ran after its actor died before the open signal: <nil>\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
            }
          }
        },
        {
          "id": "t04r6-ticket-4-check-1",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:8e308239ae4ae86afdef27b031641d497696f074d8aab62ae87771237862faee",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,pass,23760,38\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-1",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "t04r6-ticket-4-check-2",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:b0dc7889f0e663651a865b5e428a1c02bdf28143c2a19fa8422a18935f30ae79",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13907,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t04r6-ticket-4-check-3",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:787d2728dc78fceb9dfb3d70417fb4ec6d1100c5390f42da42cd578e2ec5d335",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19916,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-3",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t04r6-ticket-4-check-4",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:6ec946b5aed6508afedd6539685fc167a0904a72d6956572bc677f8c04bad328",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38894,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-4-check-4",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t04r6-ticket-4-check-5",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:0288fb7a42df0ef84c3e4202439a8b3534972559caabab3db755b7502b4225d0",
            "excerpt": "packages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4065,39\n  github.com/gibbonmi/bench/internal/subprocess/processtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,53,59\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-4-check-5",
          "command": "bench test --package ./internal/subprocess/...",
          "exit_code": 0
        },
        {
          "id": "t04r6-ticket-4-check-6",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:50ba979e7573fda8d44c4febbef025188b1b1fe9cd246b850b4fdbbe4d5d7b71",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,193,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-6",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t04r6-ticket-4-check-7",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:7fd4988f1ebdfd570c42f74d4e38148e1a097b889f514292c34bbb211c502a67",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,pass,1002,50\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-7",
          "command": "bench test --package ./internal/env",
          "exit_code": 0
        },
        {
          "id": "t04r6-ticket-4-check-8",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:74ba5ad2c27723831f40c640fd4082fe1e5565e104f128014d68aca21ba7bea1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,823,43\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-8",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "t04r6-ticket-4-check-9",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:9cfcb4cf6af32934c72f7377469bbde6f3e8924077fa3e5fd96775c7c1322524",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,26944,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-9",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t04r6-ticket-4-check-10",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:613d1388a04fc6108d0392c62db1fa1b68eebfd2893e0593a1c960b90084569a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,48041,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-check-10",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t04r6-ticket-4-omission",
          "performer": "claude:plt-t04-repair6-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
            "digest": "sha256:ebf09beb30064c7b7c2904327b5a11aa1f28ffea96ad1e30bb84af27d00856cb",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/runbinary/runbinary.go,swap,failed,2,TestBinaryClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,all,passed,38\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,494\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,24547,38\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:98: Close with a live registered group = <nil>, want a protected retention\",1\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderNormalDrain,\"lifetime_test.go:41: protection store after Close = <nil>, want absent\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-4-omission",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove normal builder drain or delete the selection before all registered groups are absent. The live descendant or executable-retention assertion must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t04-repair6-20261023",
              "digest": "sha256:ebf09beb30064c7b7c2904327b5a11aa1f28ffea96ad1e30bb84af27d00856cb",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/runbinary/runbinary.go,swap,failed,2,TestBinaryClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,all,passed,38\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,494\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,24547,38\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBinaryClose,\"lifetime_resources_test.go:98: Close with a live registered group = <nil>, want a protected retention\",1\n  github.com/gibbonmi/bench/internal/runbinary,TestBuilderNormalDrain,\"lifetime_test.go:41: protection store after Close = <nil>, want absent\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t03r1d-ticket-3-check-1",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:cec249550afb55ad3a0b32092ba909a4f6d54827cf6fbda89a2c505244969c06",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,27116,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-check-1",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t03r1d-ticket-3-check-2",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:ff52871fe011661cc1913bd8aa7c266bdd55260c7c6c6e1f85c6ca272dbdf717",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14717,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-check-2",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t03r1d-ticket-3-check-3",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:2e41d90d8c87ba1502b68c34bffac69d7aefc8f74fcc2a243a4dafb8937abc43",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,42917,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-3-check-3",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t03r1d-ticket-3-omission",
          "performer": "claude:plt-t03-repair1-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
            "digest": "sha256:f26c75a89e221e716e218d7dd258f6cd982ec6200d4151d705d16034025397ae",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGateInterruptOpensWithINTAndExits130,516,1,0,null,yes\n  2,bit,swap,failed,1,TestGateTimeoutOpensWithTERMAndExits124,516,1,0,null,yes\n  3,bit,swap,failed,1,TestGateNormalDrainRemovesReadyDescendant,516,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1164\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29718,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateInterruptOpensWithINTAndExits130,\"lifetime_test.go:56: interrupted phase opener = \\\"TERM\\\", exit = 130, interrupted = true, want INT, 130, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29275,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTimeoutOpensWithTERMAndExits124,\"lifetime_test.go:64: timed-out phase opener = \\\"INT\\\", exit = 130, interrupted = true, want TERM, 124, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,27085,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateNormalDrainRemovesReadyDescendant,\"lifetime_test.go:78: descendant 1628271 survived normal completion\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-3-omission",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Replace INT with TERM, map timeout to interruption, or omit normal group drain. The actual gate signal and descendant-survival fixtures must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t03-repair1-20261022",
              "digest": "sha256:f26c75a89e221e716e218d7dd258f6cd982ec6200d4151d705d16034025397ae",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGateInterruptOpensWithINTAndExits130,516,1,0,null,yes\n  2,bit,swap,failed,1,TestGateTimeoutOpensWithTERMAndExits124,516,1,0,null,yes\n  3,bit,swap,failed,1,TestGateNormalDrainRemovesReadyDescendant,516,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gate,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1164\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29718,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateInterruptOpensWithINTAndExits130,\"lifetime_test.go:56: interrupted phase opener = \\\"TERM\\\", exit = 130, interrupted = true, want INT, 130, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,29275,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTimeoutOpensWithTERMAndExits124,\"lifetime_test.go:64: timed-out phase opener = \\\"INT\\\", exit = 130, interrupted = true, want TERM, 124, true\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,27085,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateNormalDrainRemovesReadyDescendant,\"lifetime_test.go:78: descendant 1628271 survived normal completion\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t05r2e-ticket-5-check-1",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:75a250758a33ee23aab5ecb14ad295b1775878fdaf9c3a8ea240715d93555faa",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,pass,739,50\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-1",
          "command": "bench test --package ./internal/env",
          "exit_code": 0
        },
        {
          "id": "t05r2e-ticket-5-check-2",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:2c60f2038c174e72d5f3143d872c26966a8e8143b7bea86bc49dcb53e225a883",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,27744,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t05r2e-ticket-5-check-3",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:e4f12d0a1319c2bf6300777ae53c5cec72e32014866e7d6c810ac8f8f0aaafdd",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14681,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t05r2e-ticket-5-check-4",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:9c4bddfa0b68508661224cb4887c4865283b7140c003c04e8056e5780ad44a5d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19935,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t05r2e-ticket-5-check-5",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:c1db44ae2e88f321dd096a990d168f3d033c54a99b70fa2caff80261d9e3c3c7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,39245,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-5-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t05r2e-ticket-5-check-6",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:6ff553f88fc18258865d329636416511e1eec22942f9cae02f022b8ee13f5b24",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,pass,4070,39\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-check-6",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0
        },
        {
          "id": "t05r2e-ticket-5-check-7",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:9e8007cbba23a0abe5815bc00f9b82cc02aa27086d0b2d834b0398a5d1898894",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,pass,52,59\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "requirement": "ticket-5-check-7",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0
        },
        {
          "id": "t05r2e-ticket-5-omission",
          "performer": "claude:plt-t05-repair2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
            "digest": "sha256:10aa9b70a45e2f4da1bddc93947320ffb6968f5adb43c575e8e5dc5effb3b593",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/env/kit_run.go,swap,failed,1,TestKitRunClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/env,all,passed,50\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,489\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,822,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:59: user outcome = {Process:{Started:true StartErr:<nil> PID:1690194 PGID:1690194 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[… (508 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-5-omission",
          "command": "bench test --package ./internal/env",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the private-run descriptor or discard KitTestRun.Close failure. The child-owned file survives while a successful phase result must disappear.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t05-repair2-20261022",
              "digest": "sha256:10aa9b70a45e2f4da1bddc93947320ffb6968f5adb43c575e8e5dc5effb3b593",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/env/kit_run.go,swap,failed,1,TestKitRunClose,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/env,all,passed,50\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,489\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/env,fail,822,50\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/env,TestKitRunClose,\"kit_run_lifetime_test.go:59: user outcome = {Process:{Started:true StartErr:<nil> PID:1690194 PGID:1690194 Waited:true Code:0 Signal:signal 0 WaitErr:<nil>} Cancellation:{Cause:0 Signal:signal 0} Streams:{Stdin:{Status:none Err:<nil> Data:[… (508 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t06e-ticket-6-check-1",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:950e25e81e26c93dfc7eedf52f1e81fdb0b8061d3285ddad50099d41e1717ac8",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,pass,149,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-1",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0
        },
        {
          "id": "t06e-ticket-6-check-2",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:75f1694a9983e94f0542521f02d088392d694545cd98c046af5319471dd38266",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,27471,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t06e-ticket-6-check-3",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:987c239c950115e076ffb64e004af622101235671e778cb29c691bd48d702dab",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,52434,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t06e-ticket-6-check-4",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:99ed4907eb932389b4fa170f7a6a97a3ff9cd84f391f809f8bfb98a2c4de6b9a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13898,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t06e-ticket-6-check-5",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:43f6fb7707a2c2487f239a5af481dd9023da63a64ff2dd2f7445fd531e01e93e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38782,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-6-check-5",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t06e-ticket-6-omission",
          "performer": "claude:plt-t06-author2-20261022",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
            "digest": "sha256:656c62a780a57e2ab3d7b21ec1a197410bde070c4565b1a05d7f16e3c52d8b97",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  2,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  3,bit,omit,failed,2,TestCleanOnAnAbsentDirectory,48,1,0,null,yes\n  4,bit,swap,failed,3,TestCleanRemovesTheShardsAndKeepsTheRest,48,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gocache,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2819\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,194,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/OWUYSE/t/TestProtectionUnreadableProtection3527858357/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,262,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/SSBWIE/t/TestProtectionUnreadableProtection2451021478/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,150,48\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanOnAnAbsentDirectory,\"clean_test.go:175: clean = exit 1, want 0; out=\\\"error: cache lock unavailable — open /dev/shm/6GKYJI/t/TestCleanOnAnAbsentDirectory143597574/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionAbsentCache,\"protection_test.go:177: clean of an absent cache = (\\\"error: cache lock unavailable — open /dev/shm/6GKYJI/t/TestProtectionAbsentCache4265937827/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,195,48\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanRemovesTheShardsAndKeepsTheRest,\"clean_test.go:143: clean = exit 1, want 0; out=\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestCleanReportsThePreRemovalMeasurement,\"clean_test.go:159: clean = exit 1, want 0; out=\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionEmptyCache,\"protection_test.go:200: clean of an empty protected cache = (\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-6-omission",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0,
          "probe": {
            "mutation": "Use only the dropped owner lock, treat unreadable records as empty, or publish green before terminal persistence. The second CLI cleaner and real gate record witnesses must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:plt-t06-author2-20261022",
              "digest": "sha256:656c62a780a57e2ab3d7b21ec1a197410bde070c4565b1a05d7f16e3c52d8b97",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  2,bit,swap,failed,1,TestProtectionUnreadableProtection,48,1,0,null,yes\n  3,bit,omit,failed,2,TestCleanOnAnAbsentDirectory,48,1,0,null,yes\n  4,bit,swap,failed,3,TestCleanRemovesTheShardsAndKeepsTheRest,48,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/gocache,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2819\nmutant_report[1]{mutant}:\n  1\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,194,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/OWUYSE/t/TestProtectionUnreadableProtection3527858357/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,262,48\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionUnreadableProtection,\"protection_test.go:167: clean over unreadable protection = (\\\"go_build_cache_clean[1]{dir,bytes_removed,files_removed}:\\\\\\\\n  /dev/shm/SSBWIE/t/TestProtectionUnreadableProtection2451021478/001/.cache/bench/go-build,26,2\\\\\\\\n\\\", 0), want the protect… (251 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,150,48\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanOnAnAbsentDirectory,\"clean_test.go:175: clean = exit 1, want 0; out=\\\"error: cache lock unavailable — open /dev/shm/6GKYJI/t/TestCleanOnAnAbsentDirectory143597574/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionAbsentCache,\"protection_test.go:177: clean of an absent cache = (\\\"error: cache lock unavailable — open /dev/shm/6GKYJI/t/TestProtectionAbsentCache4265937827/001/.cache/bench/go-build/bench.lock: no such file or directory\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,195,48\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestCleanRemovesTheShardsAndKeepsTheRest,\"clean_test.go:143: clean = exit 1, want 0; out=\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestCleanReportsThePreRemovalMeasurement,\"clean_test.go:159: clean = exit 1, want 0; out=\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionEmptyCache,\"protection_test.go:200: clean of an empty protected cache = (\\\"error: cache protected — protection safe for users []; retry after they finish\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t07r2c-ticket-7-check-1",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:d1fa9cba0770afc5ca424a26e55f01dca26bfcb4ad7c5d950174e0e71c8037f5",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,837,43\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-1",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "t07r2c-ticket-7-check-2",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:7cdf3775c12dc6bbfc9e79dfe65ca55180ecfeb828c3f984b279ec64c8660ca2",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,27099,516\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-2",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t07r2c-ticket-7-check-3",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:39f9947d1fbc0559fdc7e4e9d4ed866b3cda03229b03b73c56066b78eee58b36",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,210,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,46697,210\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-3",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t07r2c-ticket-7-check-4",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:02b43fbf5f9c66fdc136e5faa93ae827ac2d7636f3c03fd94745f7647c963ccc",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3887,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-4",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t07r2c-ticket-7-check-5",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:f34fd2e74b3a2e753c34f08aaeaabae719f34f274c57b0d136f99ee67a16d805",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13802,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-5",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t07r2c-ticket-7-check-6",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:2f7e559e00507949caf6f2bfd482b8eb8706f7c4c9686b79c8d89bcab0e6bde5",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19586,404\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t07r2c-ticket-7-check-7",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:7638d01f8690b961ff108a5ea0e50330f7ae8f0d023da3de5a8f49be1848f563",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,103,123\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-check-7",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t07r2c-ticket-7-check-8",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:f08eefcf7837a7a15f8f939eaaae9ba09f12d43eeab37c346752e5074950d4a8",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38854,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "ticket-7-check-8",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t07r2c-ticket-7-omission",
          "performer": "claude:plt-t07-repair2-20261023",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3359764140326f05a33dc1b7f56e4cdec1d2b394",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
            "digest": "sha256:d466669ac1d76d5b2ff2cad90d556cf42c4efa1820d962016fbad1c4fd899a62",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/lane.go,swap,failed,1,TestLaneSelectionCloseFailure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,516\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,325\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,27365,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestLaneSelectionCloseFailure,\"lifetime_resources_test.go:264: lane with a nested user = <nil>; output:\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "ticket-7-omission",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore owner-death-only sweep, ignore a nested registered group, or discard lane selection release failure. The real Close, second-process sweep, or green-record witness must fail.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:plt-t07-repair2-20261023",
              "digest": "sha256:d466669ac1d76d5b2ff2cad90d556cf42c4efa1820d962016fbad1c4fd899a62",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/lane.go,swap,failed,1,TestLaneSelectionCloseFailure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,516\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,325\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,27365,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestLaneSelectionCloseFailure,\"lifetime_resources_test.go:264: lane with a nested user = <nil>; output:\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c3-standards-r1",
          "performer": "claude:plt-c3-standards-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-standards-20261022",
            "digest": "sha256:70bded6ceb7d792a3caa0afa460e2a1498e95e832e204988b87c3144b57a3da3",
            "excerpt": "Standards axis, PL-C3. Evidence current=true. Three findings; worst is S3, a builder zero exit with a failed capture stream reports success. No test or probe run.\n\nFindings\nS1 | ask-user | 8 | internal/env/protection.go:31, internal/runbinary/protection.go:35, internal/gate/prospectiveartifact/protection.go:15 and :69 | AGENTS.md one source per fact | The directory identity adapter is pasted three times (protect; ValidateProtection with canonical path, fixed process-users store, binding check; dev:ino identity at runbinary/protection.go:66), copying internal/gocache/protection.go:68 and :99. The copies differ only in Kind, an optional name prefix and error nouns. One source belongs in internal/subprocess/protection. Owners: 04, 05, 07.\nS2 | ask-user | 7 | internal/env/kit_run.go:140, internal/runbinary/runbinary.go:244; Resources() at kit_run.go ~125, runbinary.go ~221, prospectiveartifact.go ~284 | AGENTS.md one source per fact | \"Seal and observe, keep the directory unless Safe, else remove\" appears twice with the same error shape; the descriptor-to-slice accessor appears three times. Fold into the S1 source. Owners: 04, 05.\nS3 | ask-user | 7 | internal/gate/runner.go:~96-110 (uncertified), internal/runbinary/runbinary.go:378-381, internal/bounds/bounds.go:262-267 | AGENTS.md one source per fact; spec.md:68 story 6; ticket 04 \"A raw child status cannot certify required cleanup or publication\" | The normal-exit certification rule exists in three forms: the gate checks obligations, streams and protection; the builder checks cleanup and protection but not stream status; bounds.run checks streams but not protection. Defect: a builder zero exit with a failed captured stream returns nil. Fix: one predicate on subprocess.Outcome that all callers use. Owners: 03, 04.\n\nAdvice\n- gateProcessPolicy and builderPolicy differ only in Grace (incidental).\n- The two TestMain functions that unset protection.Env are incidental; processtest.Main could take the step.\n- The 10 ms literal in releaseArtifactBarrier repeats bounds.ProcessPoll.\n- PLnn IDs in test doc comments follow tree convention.\n\nImplementation command: name directory identity validation and normal-exit certification as single sources in the ticket (\"reuse, do not copy\").\n"
          },
          "axis": "Standards",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "f50919c538584b625cf1184bd16ab8a5e9ce5c39",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "c3-spec-r1",
          "performer": "claude:plt-c3-spec-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-spec-20261022",
            "digest": "sha256:17b941d9bfd7db94c83e2d209c2469513924328577c195b874f7a71238c89aad",
            "excerpt": "Spec axis, PL-C3. Evidence current=true, dirty_record=false. One finding; worst is P1.\n\nFindings\nP1 | auto-fix | 7 | internal/gate/run_transaction.go:194, internal/gate/lane.go:140, internal/gate/phases.go:178 | spec.md:90 story 28 \"unresolved children retain the executable they can still use\"; spec.md:413 private binary row; ticket 04 \"Expose its Descriptor through the existing selection composition path\" | Only the selection's own builder registers on Selection.Resources() (runbinary.go:296). The oracle and lane checks register on cache and bundle only; withRunBinary passes BENCH_RUN_BINARY but no descriptor. A nested phase that runs $BENCH_RUN_BINARY survives a zero oracle exit; bundle.CloseAfter(selection.Close) finds the selection generation Safe and os.RemoveAll deletes the executable (runbinary.go:250) while the bundle stays retained. PL32's test checks only the checkout; PL36 registers a user by hand. Owner: ticket 07 (its Writes holds run_transaction.go, lane.go, phases.go).\n\nContracts confirmed: gate certification and retainedFailure (PL32, PL78); KitTestRun.Close seals before removal (PL38); cache clean refuses unresolved protection under the exclusive lock; rule (a) guard at control.go with Inherited() parsed before launch and Command.Run installing the table; composition table has four rows; ADR 0019 states the decided state; PL33, PL38, PL78, PL123 tests use real helpers at the stated seams.\n\nAdvice\n- releaseUnprotected indexes disposition.Resources[0]; safe only because Recover's read-error path fills the slice.\n\nImplementation command: ticket 04 asked for descriptor exposure but gave it no gate write path, and its PL36 witness passes with a hand-registered user, so the gap fell between tickets 04 and 07.\n"
          },
          "axis": "Spec",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "f50919c538584b625cf1184bd16ab8a5e9ce5c39",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "c3-coverage-r1",
          "performer": "claude:plt-c3-coverage-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "17efbc34bfae7a14a844e11242cd9273461cbd66",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-coverage-20261022",
            "digest": "sha256:ad02304d36e6b180ef70985e23df09eba41f5de0b6c411629b7f61234f8ccd95",
            "excerpt": "Coverage axis, PL-C3. Evidence current=true at source tip 88f81d3, tree clean. Three findings, all silent probes on behavior the spec states; worst is C1.\n\nFindings\nC1 | auto-fix | high | internal/runbinary/protection.go:49 | probe silent, restored yes; spec 388 \"existing owner identity\", 391 \"Malformed, missing, foreign, and sealed descriptors prevent launch\" | No test hands a binary descriptor whose Binding does not match the selection's dev:ino; the only call is a positive check (lifetime_resources_test.go:47). Owner: ticket 04.\nC2 | auto-fix | high | internal/env/protection.go:42 | probe silent, restored yes; spec 388, 391 | No negative test for a foreign kit-test binding. Owner: ticket 05.\nC3 | auto-fix | high | internal/gate/prospectiveartifact/protection.go:56 | probe silent, restored yes; spec 388, 391 | No negative test for a foreign bundle binding. Owner: ticket 07.\ngocache has the foreign-descriptor negative test the other three adapters lack (internal/gocache/protection_test.go:77).\n\nAdvice\n- The three gaps share one shape; one foreign-binding check per owner closes them, as the chunk's single hardening cycle.\n- outcome.go:136 Uncertain versus rule (a), the retainedFailure branch in run_transaction.go, and the adapters' store and prefix checks were not probed.\n\nProbes: each swapped `if binding != d.Binding {` for `if false && binding != d.Binding {`; every mutant compiled. runbinary/protection.go on ./internal/runbinary, ./cmd/bench, ./internal/subprocess/... (silent, yes); prospectiveartifact/protection.go on ./internal/gate/prospectiveartifact (silent, yes); env/protection.go on ./internal/env (silent, yes). System suite not run.\n"
          },
          "axis": "Coverage",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "f50919c538584b625cf1184bd16ab8a5e9ce5c39",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "c3-standards-r2",
          "performer": "claude:plt-c3-standards-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-standards-r2-20261022",
            "digest": "sha256:315cd2a1c74e44da6be7cdcf3b76a731eef19bb90f551266d60433037fd081a5",
            "excerpt": "Standards axis, PL-C3 confirming round r2. Evidence current=true at 83198984e. All r1 findings closed; no new finding.\n\nS1 closed: one adapter, protection.Directory (internal/subprocess/protection/directory.go:24-94), composed by env/protection.go:9, runbinary/protection.go:12, gate/prospectiveartifact/protection.go:9, gocache/protection.go:57; no owner keeps its own identity copy.\nS2 closed: one release rule, Generation.Release (directory.go:~113-143), and one accessor, Generation.Resources(); callers env/kit_run.go:135,:127, runbinary/runbinary.go:243,:227, prospectiveartifact.go:316,:143. gocache/lock.go:171 retire seals without removal and is a different rule.\nS3 closed: one predicate, Outcome.Uncertified() (subprocess/outcome.go:26), used by gate/runner.go:84 and runbinary/runbinary.go:373; the builder now checks streams. bounds.go:264 excluded by the PL-C1 C2 decision.\n\nAdvice\n- The three owner validator pins repeat one test body (env, prospectiveartifact, runbinary); each restates the directory.go:80 binding message format.\n- The Generation.Release parameter shares the method's name in its doc.\n- Observation-failure messages in runbinary and env now read \"<name> stays retained: <err>\"; the directory is still kept.\n"
          },
          "axis": "Standards",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "72c097e5b349fb51f69d6a8bdb159a54cd45ba61",
          "finding_ids": [],
          "supersedes": [
            "c3-standards-r1"
          ]
        },
        {
          "id": "c3-spec-r2",
          "performer": "claude:plt-c3-spec-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-spec-r2-20261022",
            "digest": "sha256:cd65754ccc90f2d58ccf197efcb9863c05863df339449fe55a98e09bfab4165e",
            "excerpt": "Spec axis, PL-C3 confirming round r2. Evidence current=true, dirty_record=false at 83198984e. P1 closed; no new finding.\n\nP1 closed: internal/gate/phases.go:176 withRunBinary adds selection.Resources() to every phase, and phases.go:318 appends run.Resources() instead of overwriting; internal/gate/run_transaction.go:195 registers the oracle on cache, bundle, and selection; internal/gate/lane.go:140-142 and :296 send lane checks through withRunBinary. TestLaneSelectedExecutableUser (internal/gate/lifetime_resources_test.go:~41) and PL32's bench-run-*/bench glob (internal/systemtest/process_lifetime_test.go:~399) witness it. Of three production receivers of BENCH_RUN_BINARY, phases.go:182 and run_transaction.go:179 are covered; internal/testreport/environment.go:89/101 is outside PL-C3.\n\nContracts kept: directory.go:60-83 refuses a non-canonical path, wrong store, missing prefix or lock, missing directory, and foreign binding (spec.md:387-391); the sealed refusal stays in the protocol. Generation.Release (directory.go:119-144) keeps the safe, retained, and failed-release split that KitTestRun.Close, Selection.release, and releaseBundle relied on. outcome.go:22-39 messages match the removed gate uncertified text.\n\nAdvice\n- The focused-test child receives BENCH_RUN_BINARY but registers no selection descriptor (internal/testreport); this belongs to ticket 10 in PL-C4.\n- runBuildScript now also refuses an incomplete stream with \"builder process\" messages.\n"
          },
          "axis": "Spec",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "72c097e5b349fb51f69d6a8bdb159a54cd45ba61",
          "finding_ids": [],
          "supersedes": [
            "c3-spec-r1"
          ]
        },
        {
          "id": "c3-coverage-r2",
          "performer": "claude:plt-c3-coverage-r2-20261022",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-coverage-r2-20261022",
            "digest": "sha256:68eaba3869eca945dfb7f69423bbbba6b7558b6e594f44d265900df3c3fadd06",
            "excerpt": "Coverage axis, PL-C3 confirming round r2. Evidence current=true, dirty=false at 83198984e. C1, C2, C3 closed; no new finding.\n\nProbes (focused selections; both restored yes):\n1. internal/subprocess/protection/directory.go:79, binding comparison made never to refuse; --run Foreign across protection, runbinary, env, gate/prospectiveartifact, gocache: bit, 7 failures in all 5 packages: TestDirectoryRefusesAForeignBinding/{plain,prefix,lock}, TestSelectionValidatorRefusesAForeignBinding (closes C1), TestKitRunValidatorRefusesAForeignBinding (closes C2), TestBundleValidatorRefusesAForeignBinding (closes C3), TestProtectionRefusesAForeignCache.\n2. internal/subprocess/outcome.go:175, Uncertified stream check disabled; --run \"Uncertified|Stream|Lifetime\" across subprocess, runbinary, gate: bit, TestUncertifiedRefusesEachUncertainPart (outcome_test.go:37).\n\nAdvice\n- Only the subprocess unit test bites the stream part; the gate's use of Uncertified() was probed by a repair session. The P1 registration in phases.go was not probed this round.\n"
          },
          "axis": "Coverage",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "72c097e5b349fb51f69d6a8bdb159a54cd45ba61",
          "finding_ids": [],
          "supersedes": [
            "c3-coverage-r1"
          ]
        },
        {
          "id": "c3-standards-r3",
          "performer": "claude:plt-c3-standards-r3-20261023",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-standards-r3-20261023",
            "digest": "sha256:26bef761ebcec2deb599ae30e8c0679fda0d1309132ce681cb160f793cabcd03",
            "excerpt": "Standards axis, PL-C3 round r3 (store and start-gate amendment). Evidence current=true, dirty_record=false at d9e29a54a. Two findings, no blocker; worst is S1.\n\nFindings\nS1 | auto-fix | 5 | internal/env/kit_run.go:136-156, internal/runbinary/runbinary.go:164-171, internal/gate/prospectiveartifact/prospectiveartifact.go:317-323 | AGENTS.md one source per fact; Generation.Discard doc (protection/directory.go ~127) | Each owner carries the same order: Release, remove the resource, Discard. Nothing enforces it, and an owner that omits Discard leaks a store in the tmpfs root. One Generation helper that takes the removal makes it one source. Owner 04 (prospective site 07).\nS2 | auto-fix | 5 | internal/systemtest/process_lifetime_cache_test.go:47 | AGENTS.md one source per fact | The bare protection.Env override repeats the drop that systemEnvironmentFrom now owns (owner_environment_test.go:21). Owner 06.\n\nAdvice\n- A Bench binary that inherits BENCH_SUBPROCESS_START_GATE exits 1 with no message; the key is single-sourced and removed before exec.\n- Off Linux, darwin is coherent; other systems refuse at bootID before any gate starts.\n- Directory.Store creates directories as a side effect of Validate.\n- root_linux.go checks /dev/shm on every call, so the root can change mid-boot.\n- The \"fork/exec\" literal repeats the os.StartProcess Op string.\n\nChecked and clean: one root derivation with no environment input; gocache.ProtectionStore derives through Describe; error meanings kept; comment register.\n"
          },
          "axis": "Standards",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "e4ca40405fad4b55cf0c7209d9cdd028d751e70d",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": [
            "c3-standards-r2"
          ]
        },
        {
          "id": "c3-spec-r3",
          "performer": "claude:plt-c3-spec-r3-20261023",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-spec-r3-20261023",
            "digest": "sha256:a4e566132f29924609dbb87302145a8203125db5eee59e4aed0cfae3511a7441",
            "excerpt": "Spec axis, PL-C3 round r3 (store and start-gate amendment). Evidence current=true, dirty_record=false at d9e29a54a. Five findings; worst is P1.\n\nFindings\nP1 | auto-fix | 7 | internal/gate/prospectiveartifact/prospectiveartifact.go:241-249, internal/subprocess/protection/store.go:175, internal/subprocess/protection/recovery.go:19-21 | spec \"An absent store names no user in the current boot\" | An absent store makes lock hit ErrNotExist and SealAndObserve return Uncertain, so a sweep keeps such a bundle on every run. Only gocache/clean.go:88 treats an absent store as no users. Owner 07 (shared seam in protection: 04).\nP2 | auto-fix | 8 | prospectiveartifact.go:119-121 | \"Safe removal deletes the store after it deletes the resource\" | A Publish failure after Open removes the root but never discards the store. Owner 07.\nP3 | auto-fix | 6 | internal/env/kit_run.go:151 | same | RemoveAll of a run deletes resources a child created inside it without their stores being deleted; the likeliest cause of /dev/shm growth from 135 to 240 stores. Owner 05.\nP4 | auto-fix | 6 | runbinary.go:135, kit_run.go:57-60, prospectiveartifact.go:113-114, store.go:78-113 | same | OpenResource can fail after it creates the store; the zero generation makes Discard a no-op, so the store stays. Owners 04/05/07.\nP5 | veto | 8 | internal/subprocess/protection/directory.go:118-119 | \"a SHA-256 digest\" | The store name uses the first 16 bytes of the digest; directory_test.go:73 repeats the truncation. Owner 14.\n\nAdvice\n- Cache clean keeps the cache directory, so one store per cache path is consistent; tests that create caches under t.TempDir without SetStoreBase leave one store per test.\n- Selection, kit-run and bundle Close and sweep call Discard after a successful removal.\n- The start gate, PL25 grading, schema 1, the gate-free path, and the identityErr kill path match the spec; PL126 replaces storeBase, so the /dev/shm versus /tmp choice and darwin root are untested; the PL127 owner case skips unless root.\n- Every package init linked into the binary runs inside the gate before exec; not audited.\n"
          },
          "axis": "Spec",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "e4ca40405fad4b55cf0c7209d9cdd028d751e70d",
          "finding_ids": [
            "P1",
            "P2",
            "P3",
            "P4",
            "P5"
          ],
          "supersedes": [
            "c3-spec-r2"
          ]
        },
        {
          "id": "c3-coverage-r3",
          "performer": "claude:plt-c3-coverage-r3-20261023",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afcd6e340d9a811b7919bb28616a09602bbfc48a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-coverage-r3-20261023",
            "digest": "sha256:0e5e08c9857cc477338a69e4e5b69ff9dbd7ea6639142be4d321d243190bcb63",
            "excerpt": "Coverage axis, PL-C3 round r3 (store and start-gate amendment). Evidence current=true, clean at d9e29a54a. Four findings; worst is C1-C3, a store-leak path no test catches. All 6 probes restored.\n\nFindings\nC1 | auto-fix | 9 | internal/runbinary/runbinary.go:249 | silent, restored; spec \"Safe removal deletes the store after it deletes the resource\" | Discard swapped for return nil passes 38 tests. Owner 04.\nC2 | auto-fix | 9 | internal/env/kit_run.go:154 | silent, restored; same | Skipping Discard passes 50 tests. Owner 05.\nC3 | auto-fix | 9 | internal/gate/prospectiveartifact/prospectiveartifact.go:323 | silent, restored; same | Discard swapped for return nil passes 43 tests. Owner 07.\nC4 | auto-fix | 8 | internal/subprocess/control.go:185 | silent, restored; spec \"The actor sends that signal only after Started succeeds in every registered store\" | Releasing the gate after a failed Started passes 38 tests; no test injects a Started failure. Owner 14.\n\nProbes: runbinary.go Discard→return nil (silent, yes); kit_run.go Discard→error(nil) (silent, yes); prospectiveartifact.go Discard→return nil (silent, yes); directory.go store comparison disabled (bit TestStoreRoot, TestProtectionRefusesAForeignCache; yes); control.go Started check→false (silent, yes); control.go case c.reg != nil→case true (bit 15 tests by timeout; yes).\n\nAdvice\n- gocache clean never calls Discard; the cache directory survives a clean, so the store rule may not apply there.\n- The no-registration path bites through hangs, not a dedicated assertion.\n"
          },
          "axis": "Coverage",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "e4ca40405fad4b55cf0c7209d9cdd028d751e70d",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": [
            "c3-coverage-r2"
          ]
        },
        {
          "id": "c3-standards-r4",
          "performer": "claude:plt-c3-standards-r4-20261023",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-standards-r4-20261023",
            "digest": "sha256:8a22976875505055c535b3af0d001821dfe9e5c7f9920a485b2543181e36780e",
            "excerpt": "Standards axis, PL-C3 confirming round r4. Evidence current=true at b3eaea5f0. S1 and S2 closed; the r3 repair delta stays inside each ticket's Writes: line. Three new low-severity findings, no blocker.\n\nFindings\nN1 | auto-fix | 2 | internal/subprocess/protection/directory.go:217 | AGENTS.md one source per fact | discardNested writes the Lstat/ErrNotExist absence check inline; the same delta adds absent(path) at recovery.go:121 for that check, with a store-specific parameter name and comment. Call absent(r.ResourcePath) and make the helper path-neutral. Owner 04.\nN2 | auto-fix | 2 | internal/subprocess/startgate_test.go:146-148 | error handling kept | sealStore discards the filepath.WalkDir error and lacks t.Helper(); a failed walk leaves a partly sealed store and a misleading assertion. Owner 14.\nN3 | auto-fix | 1 | internal/gocache/clean.go:82 | comment register | The edited doc comment runs past the file's usual 100-column wrap. Owner 04.\n\nClosed\nS1: protection.Generation.Remove (directory.go:182) is the one safe removal; owners call it at env/kit_run.go:135, runbinary/runbinary.go:242, prospectiveartifact.go:120 and :317; no production caller of Release or Discard remains outside Remove.\nS2: process_lifetime_cache_test.go no longer adds protection.Env; systemEnvironmentFrom (owner_environment_test.go:16, :20-21) is the single drop.\n\nAdvice\n- Generation.Release and Discard stay exported, so the compiler does not enforce Remove; unexporting them is a follow-up candidate.\n- gocache/clean.go:89 keeps its own ErrNotExist fallback; its comment could say it now catches only a store with no generation.\n"
          },
          "axis": "Standards",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "3e0eaaad6e0d5d5ed9d83e0be33556b3e7e07220",
          "finding_ids": [
            "N1",
            "N2",
            "N3"
          ],
          "supersedes": [
            "c3-standards-r3"
          ]
        },
        {
          "id": "c3-spec-r4",
          "performer": "claude:plt-c3-spec-r4-20261023",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-spec-r4-20261023",
            "digest": "sha256:e6fcea1433f8c07233ec228e930821cf1e808f7eda3c21f22f7c62aff82f640d",
            "excerpt": "Spec axis, PL-C3 confirming round r4. Evidence current=true at b3eaea5f0 (code tip 3e0eaaad6, clean). The r3 repair closes P1-P5; no new spec deviation in e4ca40405..3e0eaaad6.\n\nFindings\nNo new finding.\n\nClosed\nP1: recovery.go:23-25 and :74-76 grade an absent store before lock through unused (:127-134), Safe with no users and no store created; store.go:182-194 lock never creates the store; the prospective sweep reaches it through prospectiveartifact.go:246 and :317; TestRecoveryAbsentStore (lifecycle_test.go:261). Spec: \"An absent store names no user in the current boot.\"\nP2: prospectiveartifact.go:119-121 sends a failed Publish through Generation.Remove, which deletes the bundle and its store (directory.go:182-196).\nP3: kit_run.go:134-136 closes through Generation.Remove; discardNested (directory.go:201-222) deletes each store of a deleted nested resource; TestGenerationRemove (lifecycle_test.go:215).\nP4: OpenResource (store.go:79-131) deletes a store it created when open fails; all three owners open through Directory.Open (directory.go:65-74); TestOpenFailureRemovesItsOwnStore (lifecycle_test.go:284).\nP5: directory.go:60-61 matches spec line 268 (first 16 bytes of a SHA-256 digest, in hexadecimal).\n\nAdvice\n- SealAndObserve's absent-store check (recovery.go:23) runs outside the protocol lock; a comment could state the owner-only creation assumption.\n- discardNested scans the store root without the lock and relies on the outer generation being safe; its doc comment could say so.\n- clean.go:89-90 can run release twice if unused's release error wraps fs.ErrNotExist; minor, not a spec deviation.\n"
          },
          "axis": "Spec",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "3e0eaaad6e0d5d5ed9d83e0be33556b3e7e07220",
          "finding_ids": [],
          "supersedes": [
            "c3-spec-r3"
          ]
        },
        {
          "id": "c3-coverage-r4",
          "performer": "claude:plt-c3-coverage-r4-20261023",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c5866e7b844ae4aa5a10238969fc23ba7c557e8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:plt-c3-coverage-r4-20261023",
            "digest": "sha256:ee7a5c34e5537751adeb525c2ad616460bd5cd2f2d2ad9427ba3081e00fb7249",
            "excerpt": "Coverage axis, PL-C3 confirming round r4. Evidence current=true at b3eaea5f0. C1-C4 closed at their current sites; the discard moved into one seam, protection.Generation.Remove (directory.go:192), and every probe there bit and restored. One residual gap: no owner package test observes the store discard.\n\nFindings\nR4-C1 | veto | 4 | internal/runbinary/runbinary.go:242, internal/env/kit_run.go:135, internal/gate/prospectiveartifact/prospectiveartifact.go:120 and :317 | inferred from passing owner packages, no owner-site probe; spec \"Safe removal deletes the store after it deletes the resource\" | Both directory.go mutants bit only TestGenerationRemove; runbinary (38), env (50) and prospectiveartifact (43) passed with the store present. An owner that went back to Release plus RemoveAll would pass. Owners 04, 05, 07.\n\nProbes:\n- directory.go:192 own-store Discard dropped: bit TestGenerationRemove; restored yes.\n- directory.go:192 discardNested dropped: bit TestGenerationRemove; restored yes.\n- control.go:186 gate released on failed Started (single form): bit TestStartGateStartedFailure; restored yes.\n\nClosed\nC1, C2, C3: runbinary.go:242, kit_run.go:135, prospectiveartifact.go:120 and :317 go through Generation.Remove, and its discard site bites; residual is R4-C1.\nC4: TestStartGateStartedFailure bites. Spec: \"The actor sends that signal only after Started succeeds in every registered store.\"\n\nAdvice\n- One store-absence assertion per owner close test closes R4-C1.\n"
          },
          "axis": "Coverage",
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "tip": "3e0eaaad6e0d5d5ed9d83e0be33556b3e7e07220",
          "finding_ids": [
            "R4-C1"
          ],
          "supersedes": [
            "c3-coverage-r3"
          ]
        }
      ],
      "confirmations": [
        {
          "id": "c3-confirm-r2",
          "performer": "claude:plt-orchestrator-20261022",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "high",
          "source_digest": "173af91dfd496d875cb634e90e7054df6c6486aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:plt-orchestrator-20261022",
            "digest": "sha256:930efda2f3c8ab3aeb2d519726614dbba234381a6da76dcd3f2ea549c6721ee5",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/phases.go,swap,failed,1,TestLaneSelectedExecutableUser,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,516\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,334\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,25878,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestLaneSelectedExecutableUser,\"lifetime_resources_test.go:253: retained bench-run-*/bench = [], <nil>; want one\",1\nskips[0]{package,test,reason}:\n"
          },
          "base": "3e9708ec537419d91b9740dd5238fe13f02e874a",
          "reviewed_tip": "72c097e5b349fb51f69d6a8bdb159a54cd45ba61",
          "repaired_tip": "75c0cb5371dde7a9be08dab3760ca07f821ffc23",
          "finding_ids": [],
          "supersedes": [],
          "probe": {
            "mutation": "withRunBinary drops the selection's descriptor from every phase's resources",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:plt-orchestrator-20261022",
              "digest": "sha256:930efda2f3c8ab3aeb2d519726614dbba234381a6da76dcd3f2ea549c6721ee5",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/phases.go,swap,failed,1,TestLaneSelectedExecutableUser,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,516\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,334\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,25878,516\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestLaneSelectedExecutableUser,\"lifetime_resources_test.go:253: retained bench-run-*/bench = [], <nil>; want one\",1\nskips[0]{package,test,reason}:\n"
            }
          }
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
      "from": "sha256:a3d8f9f0d05b81006d7e1a91c4898a0c28fee10543b663198fcf80e08429465d",
      "to": "sha256:c36d67f25b6ac6ce1b621f6b196cb4fb4e063b2830bdb5ba4bb7405337a33930",
      "chunk_ids": {
        "PL-C1": [
          "PL-C1"
        ]
      }
    },
    {
      "from": "sha256:c36d67f25b6ac6ce1b621f6b196cb4fb4e063b2830bdb5ba4bb7405337a33930",
      "to": "sha256:c50d5467850cec958ff45993f032691dee138de7662bca7b6d67d6e589266777",
      "chunk_ids": {
        "PL-C1": [
          "PL-C1"
        ]
      }
    },
    {
      "from": "sha256:c50d5467850cec958ff45993f032691dee138de7662bca7b6d67d6e589266777",
      "to": "sha256:dbb8b18b9518678f70a2d2a437f7ae9e355507e4dbb28f1b411810cc37504e8d",
      "chunk_ids": {
        "PL-C1": [
          "PL-C1"
        ],
        "PL-C2": [
          "PL-C2"
        ]
      }
    },
    {
      "from": "sha256:dbb8b18b9518678f70a2d2a437f7ae9e355507e4dbb28f1b411810cc37504e8d",
      "to": "sha256:eba58dd14aa0003e578bafade100c714a3dd1c554ce3508af878811fb8342bef",
      "chunk_ids": {
        "PL-C1": [
          "PL-C1"
        ],
        "PL-C2": [
          "PL-C2"
        ]
      }
    },
    {
      "from": "sha256:eba58dd14aa0003e578bafade100c714a3dd1c554ce3508af878811fb8342bef",
      "to": "sha256:cce28f6ec806cd14a2584f82658ea24c135b62aa57a76382469b7b815c01569a",
      "chunk_ids": {
        "PL-C1": [
          "PL-C1"
        ],
        "PL-C2": [
          "PL-C2"
        ],
        "PL-C3": [
          "PL-C3"
        ]
      }
    },
    {
      "from": "sha256:cce28f6ec806cd14a2584f82658ea24c135b62aa57a76382469b7b815c01569a",
      "to": "sha256:1d5c9d89d3d5bf071cbc0c8b773d92c28686b143851d1ead9f1a84d261f9bfd6",
      "chunk_ids": {
        "PL-C1": [
          "PL-C1"
        ],
        "PL-C2": [
          "PL-C2"
        ],
        "PL-C3": [
          "PL-C3"
        ]
      }
    },
    {
      "from": "sha256:1d5c9d89d3d5bf071cbc0c8b773d92c28686b143851d1ead9f1a84d261f9bfd6",
      "to": "sha256:3bf4de5bb086b76069553ee314c6fef0e5bddf003656cb21a2463f275ad0ac88",
      "chunk_ids": {
        "PL-C1": [
          "PL-C1"
        ],
        "PL-C2": [
          "PL-C2"
        ],
        "PL-C3": [
          "PL-C3"
        ]
      }
    },
    {
      "from": "sha256:3bf4de5bb086b76069553ee314c6fef0e5bddf003656cb21a2463f275ad0ac88",
      "to": "sha256:caa570aec8f6c35f7fea13e9d45f205ccac5e0d483c9580efbc060af4fbca693",
      "chunk_ids": {
        "PL-C1": [
          "PL-C1"
        ],
        "PL-C2": [
          "PL-C2"
        ],
        "PL-C3": [
          "PL-C3"
        ]
      }
    },
    {
      "from": "sha256:caa570aec8f6c35f7fea13e9d45f205ccac5e0d483c9580efbc060af4fbca693",
      "to": "sha256:af4cfdbe03e0fa38cd38cd92c2f3781a1f7bcedc6bc1446394bf266282ef64b4",
      "chunk_ids": {
        "PL-C1": [
          "PL-C1"
        ],
        "PL-C2": [
          "PL-C2"
        ],
        "PL-C3": [
          "PL-C3"
        ]
      }
    }
  ]
}
```

## Repair state

PL-C1 consumed repair cycle 1 of 2 after its initial review (round r1). No reviewer extension applies.
Repair 1 (commit 36c2fa86a) closed S1, S2, C1, and C3. The C2 disposition moved to no-op by consultant tie-break, as a record-only correction.
The confirming round r2 passed on all three axes with no new finding, and the Coverage axis reaffirmed the C2 no-op.

PL-C2 consumed repair cycle 1 of 2 after its initial review (round r1). No reviewer extension applies. The reviewer directed minimal probing during this round.
Repair 1 (commit 3e9708ec5, test-only) closed P2, C1 to C4, and the PL110 half of P1. The Opus and Sonnet consultants agreed that the PL71 half of P1 is a no-op.
The confirming round r2 passed on all three axes with no new finding, and the Spec axis reaffirmed the PL71 no-op.

PL-C3 took two pre-review repairs before its record commit: ticket 04 (commit 9f83d4328, test-only) and ticket 05 (commit f50919c53, test-only). These repairs fixed reds that the floor checks found at the frozen tip; they consume no review repair cycle.
PL-C3 round r1 found 7 findings. The orchestrator set every finding to auto-fix under the reviewer's advance approval.

Repair cycle 1 of 2 ran one fresh repair session per affected ticket, and the repairs closed all 7 findings.
The repair commits are 918c78c13 and f44d12d4b (ticket 04), 148f5f7d9 (ticket 05), 5f675b171 (ticket 07), and 72c097e5b (ticket 03).
The confirming round r2 passed on all three axes with no new finding.

The PL-C3 checkpoint gate then went red: launches serialized on durable writes under every inherited store's lock.
The reviewer approved an amendment on 2026-10-23: tmpfs process-user stores, no inherited vector in system tests, and a start gate in new ticket 14.
Round r3 reviewed that amendment and found 11 findings. Repair cycle 2 of 2 ran one repair session for ticket 04 (commit 65a501d96) and one for ticket 14 (commit 3e0eaaad6).

The confirming round r4 closed all 11 r3 findings. The Spec axis passed, and the Standards and Coverage axes found 4 new low-severity findings.
Repair cycle 3 runs past the allowance under the reviewer's advance approval: one repair session for ticket 04 and one for ticket 14.

## Standards

PL-C1 round r1: 2 findings. Worst issue: S2, a scheduling-dependent verdict in a PL66 test.

- S1 | auto-fix | confidence 6 | `internal/subprocess/control.go:248-249` | `controller.signal` reclassifies ESRCH as absence, which duplicates the rule that `protection.observeGroup` owns (AGENTS.md one source per fact).
- S2 | auto-fix | confidence 6 | `internal/subprocess/lifetime_test.go:246-250` | `TestLifetimeSingleFinalWindow` orders its events with a raw 50 ms pause (craft-tdd broken-oracle rule).

PL-C2 round r1: 1 finding. Worst issue: S1, a fixture scaffold that repeats incidental text.

- S1 | no-op | confidence 4 | `internal/subprocess/protection/fixture_test.go:138` | Two bounded-wait helpers repeat the select scaffold of `waitWithin`; the deadline policy keeps one source in `bounds`.

PL-C3 round r1: 3 findings. Worst issue: S3, a builder zero exit with a failed capture stream reports success.

- S1 | auto-fix | confidence 8 | `internal/env/protection.go:31`, `internal/runbinary/protection.go:35`, `internal/gate/prospectiveartifact/protection.go:15` | The directory identity adapter has three copies (AGENTS.md one source per fact). The copies differ only in kind, name prefix, and error nouns. One source goes in `internal/subprocess/protection`.
- S2 | auto-fix | confidence 7 | `internal/env/kit_run.go:140`, `internal/runbinary/runbinary.go:244` | The seal-then-remove-unless-safe rule has two copies, and the descriptor accessor has three. They fold into the S1 source.
- S3 | auto-fix | confidence 7 | `internal/gate/runner.go:96-110`, `internal/runbinary/runbinary.go:378-381` | The normal-exit certification rule has two forms, and the builder form omits the stream check (spec.md:68 story 6). One predicate on `subprocess.Outcome` serves the gate and the builder. `bounds` keeps its own composition, as the PL-C1 C2 decision states.

PL-C3 round r3: 2 findings. Worst issue: S1, the release, removal, and store-deletion order has three copies.

- S1 | auto-fix | confidence 5 | `internal/env/kit_run.go:136-156`, `internal/runbinary/runbinary.go:164-171`, `internal/gate/prospectiveartifact/prospectiveartifact.go:317-323` | Each owner repeats release, removal, and `Discard` (AGENTS.md one source per fact). One `Generation` helper takes the removal. Owner: ticket 04.
- S2 | auto-fix | confidence 5 | `internal/systemtest/process_lifetime_cache_test.go:47` | The bare `protection.Env` override repeats the drop that `systemEnvironmentFrom` owns. Owner: ticket 04, widened.

PL-C3 round r4: 3 findings. Worst issue: N1, `discardNested` repeats the absence check of the `absent` helper.

- N1 | auto-fix | confidence 2 | `internal/subprocess/protection/directory.go:217` | `discardNested` writes the absence check inline, and the `absent` helper at `recovery.go:121` owns it (AGENTS.md one source per fact). The helper takes a path-neutral name and comment. Owner: ticket 04.
- N2 | auto-fix | confidence 2 | `internal/subprocess/startgate_test.go:146-148` | `sealStore` discards the walk error and does not call `t.Helper()`. Owner: ticket 14.
- N3 | auto-fix | confidence 1 | `internal/gocache/clean.go:82` | The edited doc comment runs past the file's wrap width. Owner: ticket 04.

## Spec

PL-C1 round r1: 0 findings. Worst issue: none.

PL-C2 round r1: 2 findings. Worst issue: P1, a launch-refusal assertion that refuses only an empty descriptor.

- P1 | auto-fix | confidence 6 | `internal/subprocess/protection/registration_test.go:204-206` and `record_test.go:225-229` | The PL71 and PL110 tests never launch against the unsynced record or the unsafe descriptor. The PL110 half is repaired. The PL71 half is a no-op, because row PL71 states the first launch only, and the `OpenResource` error witnesses it. A lasting refusal would add a record state that no row approves (spec.md:317-318, :661).
- P2 | auto-fix | confidence 5 | `internal/subprocess/protection/record_test.go:59-61` | The PL39 foreign-binding test never changes the `Binding` field.

PL-C3 round r1: 1 finding. Worst issue: P1, the selected executable is deleted while a nested user that inherits it survives.

- P1 | auto-fix | confidence 7 | `internal/gate/run_transaction.go:194`, `internal/gate/lane.go:140`, `internal/gate/phases.go:178` | The oracle and the lane checks never register on the selection's descriptor (spec.md:90 story 28, spec.md:413). The selection then reads safe. Its executable is removed while a nested phase can still run it. Owner: ticket 07.

PL-C3 round r3: 5 findings. Worst issue: P1, a bundle whose store is absent stays retained forever.

- P1 | auto-fix | confidence 7 | `internal/subprocess/protection/store.go:175`, `recovery.go:19-21` | An absent store returns uncertain, but the spec says an absent store names no user in the current boot. Owner: ticket 04.
- P2 | auto-fix | confidence 8 | `internal/gate/prospectiveartifact/prospectiveartifact.go:119-121` | A Publish failure removes the bundle but keeps its store. Owner: ticket 04.
- P3 | auto-fix | confidence 6 | `internal/env/kit_run.go:151` | Run removal deletes resources that children created inside it and keeps their stores. Owner: ticket 04.
- P4 | auto-fix | confidence 6 | `internal/subprocess/protection/store.go:78-113` | A failed OpenResource can keep the store it created. Owner: ticket 04.
- P5 | auto-fix | confidence 8 | `internal/subprocess/protection/directory.go:118-119` | The store name uses the first 16 bytes of the digest. A prose-only spec correction states that length.

PL-C3 round r4: 0 findings. Worst issue: none.

## Coverage

PL-C1 round r1: 4 findings. Worst issue: C1, the expired-deadline KILL of a late Start has no biting test.

- C1 | auto-fix | confidence 7 | `internal/subprocess/control.go:177` | No test bites the forced KILL without another grace for a late Start after an expired deadline (spec.md:136).
- C2 | no-op | confidence 5 | `internal/bounds/bounds.go:257` | No test fails when bounds certifies completion over a failed or force-closed stream (ticket 01 line 40). The repair found that no real child reaches the loop under the production policy. The consultants split, and the fable tie-break chose no-op. The loop is the bounds status projection over the separate stream facts of `subprocess.Outcome` (spec.md:38; ticket 01:25). Each later caller composes those facts under its own contract (spec.md:169-183; ticket 09:16). The branch stays as defense-in-depth with no test.
- C3 | auto-fix | confidence 5 | `internal/subprocess/outcome.go:128` | No test fails when a forced closure reports the cleanup as not requested (spec.md:146).
- C4 | no-op | confidence 4 | `internal/subprocess/control.go:81` | PL89 in PL-C5 is the planned witness for the platform-attribute rule.

PL-C2 round r1: 4 findings. Worst issue: C3, a kind with no validator row gets admission with no biting test.

- C1 | auto-fix | confidence 8 | `internal/gocache/protection.go:136` | No test bites the canonical-path check; a symlink parent component would pass.
- C2 | auto-fix | confidence 6 | `internal/gocache/protection.go:142` | No test bites the lock-file check; the binding check refuses the test case first.
- C3 | auto-fix | confidence 8 | `internal/subprocess/protection/store.go:43-44` | No test bites the refusal of a kind with no validator row.
- C4 | auto-fix | confidence 6 | `internal/bounds/bounds.go:226-229` | No bounds test sets a malformed `BENCH_PROCESS_RESOURCES`; only the parser is tested.

PL-C3 round r1: 3 findings. Worst issue: C1, a foreign binary descriptor passes validation with no biting test.

- C1 | auto-fix | confidence 8 | `internal/runbinary/protection.go:49` | No test bites the binding comparison for a binary descriptor (spec.md:388, :391). Owner: ticket 04.
- C2 | auto-fix | confidence 8 | `internal/env/protection.go:42` | No test bites the binding comparison for a kit-test descriptor. Owner: ticket 05.
- C3 | auto-fix | confidence 8 | `internal/gate/prospectiveartifact/protection.go:56` | No test bites the binding comparison for a bundle descriptor. Owner: ticket 07.

PL-C3 round r3: 4 findings. Worst issue: C1 to C3, no test bites a skipped store deletion.

- C1 | auto-fix | confidence 9 | `internal/runbinary/runbinary.go:249` | Skipping `Discard` passes every test. Owner: ticket 04.
- C2 | auto-fix | confidence 9 | `internal/env/kit_run.go:154` | Skipping `Discard` passes every test. Owner: ticket 04.
- C3 | auto-fix | confidence 9 | `internal/gate/prospectiveartifact/prospectiveartifact.go:323` | Skipping `Discard` passes every test. Owner: ticket 04.
- C4 | auto-fix | confidence 8 | `internal/subprocess/control.go:185` | No test injects a Started failure, so opening the gate after it passes. Owner: ticket 14.

PL-C3 round r4: 1 finding. Worst issue: R4-C1, no owner test sees its store deleted after a safe close.

- R4-C1 | auto-fix | confidence 4 | `internal/runbinary/runbinary.go:242`, `internal/env/kit_run.go:135`, `internal/gate/prospectiveartifact/prospectiveartifact.go:120` | Only `TestGenerationRemove` bites a skipped store deletion. An owner that went back to release and removal without `Generation.Remove` passes every owner test. One store-absence assertion goes in each owner close test. The axis marked it veto; the orchestrator set it to auto-fix under the reviewer's advance approval. Owner: ticket 04, widened.

## Advice

- The bounds runner does not forward a caller's `cmd.Stdin`; no current caller sets it.
- `bounds.ProcessPoll` is an exported mutable variable.
- `TestLifetimeResistantGroup` does not prove that TERM came before KILL.
- The `Unbounded` window branch, drain mode, `StreamFailed`, and `Cleanup.SignalErr` have no PL-C1 test; later rows cover drain mode.
- PL-C2: the `requireRefused` call on the empty descriptor in `TestRegistrationResourceSyncFailure` proves nothing and can go.
- PL-C2: a later consumer must not treat `Outcome.Protection` as release proof; release authority stays with `SealAndObserve`.
- PL-C3: the 10 ms poll literal in `releaseArtifactBarrier` repeats `bounds.ProcessPoll`.
- PL-C3: `releaseUnprotected` indexes `disposition.Resources[0]`, which is safe only because the read-error path of `Recover` fills the slice.
- PL-C3: the `Uncertain` state at `outcome.go:136` and the store and prefix checks of the adapters were not probed.
- PL-C3 r2: the focused-test child receives `BENCH_RUN_BINARY` but registers no selection descriptor in `internal/testreport`. Ticket 10 in PL-C4 owns that path.
- PL-C3 r2: the three owner validator pins repeat one test body; each pin has a recorded red, so the test-expectation exception holds.
- PL-C3 r3: a Bench binary that inherits the start-gate key exits 1 with no message.
- PL-C3 r3: `Directory.Store` creates directories as a side effect of `Validate`.
- PL-C3 r3: the root check reads `/dev/shm` on every call, so the root can change during one boot.
- PL-C3 r3: PL126 replaces the store base, so the `/dev/shm` and `/tmp` choice and the darwin root have no test.
- PL-C3 r3: every package `init` in the binary runs inside the start gate before exec.
- PL-C3 r4: `Generation.Release` and `Generation.Discard` stay exported, so the compiler does not force an owner through `Generation.Remove`.
- PL-C3 r4: the absent-store check in `SealAndObserve` and the root scan in `discardNested` run outside the protocol lock. Each relies on the owner alone to create a store.
- PL-C3 r4: `gocache/clean.go:89-90` can run release twice when the release error of `unused` wraps `fs.ErrNotExist`.
