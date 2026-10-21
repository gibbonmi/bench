# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/test-cost-signal/spec.md",
  "plan_digest": "sha256:11f891ae58901a2217867f2bcc653246fc6ee06a85da20fee822a3d0ede0ec54",
  "implementation_session": "",
  "chunks": [
    {
      "id": "TC-C1",
      "base": "7ed6d10a5fa7d0fefde556c10487f050a9df17a7",
      "tip": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
      "plan_digest": "sha256:a030a302362c213d596f53bf8de12a525abc5ee5de38c4ee5d1129e3aaaf1407",
      "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
      "acceptance_rows": [
        "TC1",
        "TC2",
        "TC3",
        "TC4",
        "TC5",
        "TC6",
        "TC7",
        "TC8",
        "TC9",
        "TC10",
        "TC11",
        "TC38"
      ],
      "verification": [
        {
          "id": "t1-bounds",
          "performer": "claude:tcs-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t1-author-20261021",
            "digest": "sha256:910cf032e678f12fe51e2828b5329301b73250c880df536d7eceef113211121c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,pass,567,61\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-bounds",
          "command": "bench test --package ./internal/bounds",
          "exit_code": 0
        },
        {
          "id": "t1-gate",
          "performer": "claude:tcs-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t1-author-20261021",
            "digest": "sha256:55d8b97a112f8b66848a7e6b4468be2ce2796568a640ce635fad925d830f6c59",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,34547,508\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t1-preprelease",
          "performer": "claude:tcs-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t1-author-20261021",
            "digest": "sha256:02f29fb8d3a0669833ff057f5bbffb71cb6ba5e4c96e21805500aff7462bf050",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preprelease,pass,14,12\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-preprelease",
          "command": "bench test --package ./internal/preprelease",
          "exit_code": 0
        },
        {
          "id": "t1-testreport",
          "performer": "claude:tcs-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t1-author-20261021",
            "digest": "sha256:ed5ba40870867b72062e2a2b525db8b997ed34e5d929ded5ff1be4183fd2d149",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,63529,240\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t1-conformance",
          "performer": "claude:tcs-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t1-author-20261021",
            "digest": "sha256:6bda3566c96e68baae3ec037c446bb7eeeace8a01a56e20475d71e3997dd3a06",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,43280,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-bounds-policy",
          "performer": "claude:tcs-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t1-author-20261021",
            "digest": "sha256:bb5446e5dfe6337f12dda75c4d9af63b16700965da02da6526b114c8fd660a3e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  bounds-policy,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,247,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-bounds-policy",
          "command": "bench test --check bounds-policy",
          "exit_code": 0
        },
        {
          "id": "t1-root-conformance",
          "performer": "claude:tcs-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t1-author-20261021",
            "digest": "sha256:27d4e4eadd45bfaf3b76d38ac7819cc055855660c5f4e7251c713b95a3a3ff6f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12648,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-policy-pin",
          "performer": "claude:tcs-t1-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t1-author-20261021",
            "digest": "sha256:98beb82ed7202216555f4a8da42ac43994a0012b3358596bb1b44fd272d76cab",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/bounds/bounds.go,swap,failed,1,TestProductionPolicyValues,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/bounds,^TestProductionPolicyValues$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,374\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/bounds,TestProductionPolicyValues,\"bounds_test.go:87: duration policy changed: provider=10s environment_discovery=2s refresh=30s guard=5s gate=45m0s go_test=5m0s\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-policy-pin",
          "command": "bench test --package ./internal/bounds --run '^TestProductionPolicyValues$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap GoTestTimeout = 10 * time.Minute for GoTestTimeout = 5 * time.Minute in internal/bounds/bounds.go; TestProductionPolicyValues must red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:tcs-t1-author-20261021",
              "digest": "sha256:98beb82ed7202216555f4a8da42ac43994a0012b3358596bb1b44fd272d76cab",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/bounds/bounds.go,swap,failed,1,TestProductionPolicyValues,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/bounds,^TestProductionPolicyValues$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,374\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/bounds,TestProductionPolicyValues,\"bounds_test.go:87: duration policy changed: provider=10s environment_discovery=2s refresh=30s guard=5s gate=45m0s go_test=5m0s\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c1-standards-r0",
          "performer": "claude:tcs-c1-standards-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c1-standards-20261021",
            "digest": "sha256:01e90a7edeb1761bcd126472671bc4161cbaf94d06d0aefec34d561118ed2dca",
            "excerpt": "Standards axis, TC-C1 (ticket 1-timeout-owner): 0 findings. Worst issue: none. One source per fact holds: BaseTestArgv in internal/gate/phases.go is the only production reader of bounds.GoTestTimeout for the argv; each exact-argv test builds the token by its own expression; the only literal timeout is the TestProductionPolicyValues pin (bounds_test.go:86). Every edited path is in Writes:, and .bench/commitment.json is not in the delta. Unknown (not a finding): the recorded red of the 5-minute mutation was outside the read set; the coordinator confirmed it in the ticket commit message and the t1-policy-pin record entry.\n"
          },
          "axis": "Standards",
          "base": "7ed6d10a5fa7d0fefde556c10487f050a9df17a7",
          "tip": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1-spec-r0",
          "performer": "claude:tcs-c1-spec-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c1-spec-20261021",
            "digest": "sha256:a15fbf25ae2dba839de7104ed7ea514354ddee933a78c1e0bdb758dfbf90964e",
            "excerpt": "Spec axis, TC-C1: 1 finding, 0 actionable. P1 (ask-user, confidence 5): the recorded red of GoTestTimeout = 5 * time.Minute was not in the code delta (reviews/ excluded). Coordinator disposition: no-op; the ticket commit 6ba2c0406 message records the red 'duration policy changed: ... go_test=5m0s', and reviews/test-cost-signal.md holds the t1-policy-pin probe entry (bit, restored yes). Clean: the six BaseTestArgv callers inherit the token right after -count=1; acceptance rows 1-5 hold; TC11 projects/benchkit.md:345-347 shows -timeout=10m0s in the test, race, and system rows.\n"
          },
          "axis": "Spec",
          "base": "7ed6d10a5fa7d0fefde556c10487f050a9df17a7",
          "tip": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1-coverage-r0",
          "performer": "claude:tcs-c1-coverage-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "ccc1785aef4112d6c998d2d9849bffdd7fdc8d11",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c1-coverage-20261021",
            "digest": "sha256:31c67fd0b89e8d942077c59d50de4cdb82627ee0c21bc038ce084e50d395fb24",
            "excerpt": "Coverage axis, TC-C1 (TC1-TC10, TC38): 0 findings. Worst issue: none. Every BaseTestArgv caller is pinned: system TC3, core TC1, race TC4 census prefix, preprelease TC5, focused TC6, check children TC7-TC9. No expectation calls BaseTestArgv; no 10m0s literal in a Go test. Probe: swap in internal/preprelease/preprelease.go replacing gate.BaseTestArgv(kit, \"-tags\" with a literal argv without the token: bit at TestShipConformanceStepArgvPrefix, restored yes. Residual unknown: the GOFLAGS precedence claim in the BaseTestArgv comment is not graded by a row.\n"
          },
          "axis": "Coverage",
          "base": "7ed6d10a5fa7d0fefde556c10487f050a9df17a7",
          "tip": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "TC-C2",
      "base": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
      "tip": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
      "plan_digest": "sha256:073860a0a02894a62a0f2a8cde55044b870e1588ff27390dc6320ca13f44ecb7",
      "source_digest": "b604b56d6c000230b4a9b8238003fb8e1f8927ac",
      "acceptance_rows": [
        "TC12",
        "TC13",
        "TC14",
        "TC15",
        "TC16",
        "TC17",
        "TC18",
        "TC19",
        "TC20",
        "TC21",
        "TC22",
        "TC23",
        "TC24",
        "TC25",
        "TC26",
        "TC39",
        "TC40"
      ],
      "verification": [
        {
          "id": "t2-testreport",
          "performer": "claude:tcs-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bd4d0e07d53386e9678bb22ae282c164e7949229",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-author-20261021",
            "digest": "sha256:74bb0eea0f1560dcb27bb3f6ce33818b5697aa7809236eff78e2f89d66db9ba3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,59028,253\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nslowest[3]{package,test,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,TestRootConformanceCheckPrintsFailures,20600\n  github.com/gibbonmi/bench/internal/testreport,TestNamedCheckRunsOnlyRegisteredDevScope,15000\n  github.com/gibbonmi/bench/internal/testreport,TestCommandDrainsBuilderGroupOnCancelSignal,6050\n"
          },
          "requirement": "t2-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t2-root-conformance",
          "performer": "claude:tcs-t2-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "bd4d0e07d53386e9678bb22ae282c164e7949229",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-author-20261021",
            "digest": "sha256:7ba1849b212ae216bfa227e14565c96661ec21139577e8ac9f389b4e8adf639f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13537,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nslowest[1]{package,test,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,13520\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-r1-t2-testreport",
          "performer": "claude:tcs-t2-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "aaf7ffde7e3059c73ac481b2fb34e993a3a7b3f2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-repair1-20261021",
            "digest": "sha256:e0937ba38ba3ffd443ad3086836f805ac3236dd015660a1b685da47e7eaf15d2",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,66923,255\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nslowest[3]{package,test,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,TestRootConformanceCheckPrintsFailures,20760\n  github.com/gibbonmi/bench/internal/testreport,TestNamedCheckRunsOnlyRegisteredDevScope,14710\n  github.com/gibbonmi/bench/internal/testreport,TestCommandDrainsBuilderGroupOnCancelSignal,6050\n"
          },
          "requirement": "t2-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t2-r1-t2-root-conformance",
          "performer": "claude:tcs-t2-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "aaf7ffde7e3059c73ac481b2fb34e993a3a7b3f2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-repair1-20261021",
            "digest": "sha256:ace74449f4ab1f6b28090244670a1ec9cb3276953a5a62bafaf79787f01b8f59",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12425,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nslowest[1]{package,test,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,12410\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-r2-t2-testreport",
          "performer": "claude:tcs-t2-repair2-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3b90591aac1d7ac27cb5ac499be470956ed3fe68",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-repair2-20261021",
            "digest": "sha256:b048d359264f36759f6ec759fa51ff9f2455d2054911aaaea8c881b1f26181e3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,65479,255\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nslowest[3]{package,test,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,TestRootConformanceCheckPrintsFailures,21980\n  github.com/gibbonmi/bench/internal/testreport,TestNamedCheckRunsOnlyRegisteredDevScope,16350\n  github.com/gibbonmi/bench/internal/testreport,TestCommandDrainsBuilderGroupOnCancelSignal,6060\n"
          },
          "requirement": "t2-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t2-r2-t2-root-conformance",
          "performer": "claude:tcs-t2-repair2-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3b90591aac1d7ac27cb5ac499be470956ed3fe68",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-repair2-20261021",
            "digest": "sha256:c9e09c8b3480b6a52f51272bab15798b87acae7cb73dea23fde46dc96a7f06c5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12357,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nslowest[1]{package,test,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,12340\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-r3-t2-testreport",
          "performer": "claude:tcs-t2-repair3-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9be17cb51fd80fb5567656b3ebfd8a142f0213c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-repair3-20261021",
            "digest": "sha256:0bcc2c73df01db7aefe6e2730713175a0afc555f8e14ae37c5944f5e8ac70101",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,59107,255\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nslowest[3]{package,test,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/testreport,TestRootConformanceCheckPrintsFailures,20740\n  github.com/gibbonmi/bench/internal/testreport,TestNamedCheckRunsOnlyRegisteredDevScope,14780\n  github.com/gibbonmi/bench/internal/testreport,TestCommandDrainsBuilderGroupOnCancelSignal,6040\n"
          },
          "requirement": "t2-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t2-r3-t2-root-conformance",
          "performer": "claude:tcs-t2-repair3-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9be17cb51fd80fb5567656b3ebfd8a142f0213c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-repair3-20261021",
            "digest": "sha256:d21b31bf831a8c02c97a1586b1137e62bf1ebc001c6b0a4563901307a329db45",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12592,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nslowest[1]{package,test,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,12580\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-r4-t2-testreport",
          "performer": "claude:tcs-t2-repair4-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b604b56d6c000230b4a9b8238003fb8e1f8927ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-repair4-20261021",
            "digest": "sha256:0328a31cb5eb8b683119cafa7fbe6977639a46ea8359f1e14e00a1c468270b2b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,59260,257\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t2-r4-t2-root-conformance",
          "performer": "claude:tcs-t2-repair4-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b604b56d6c000230b4a9b8238003fb8e1f8927ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-repair4-20261021",
            "digest": "sha256:e3580680ba59b6e9b0597f4f9f4529c832cdf61b4e62a517d49aef5d59d3a085",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13463,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-r4-t2-record-journey",
          "performer": "claude:tcs-t2-repair4-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b604b56d6c000230b4a9b8238003fb8e1f8927ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-repair4-20261021",
            "digest": "sha256:f7b8773c18f5ac538d7f3b641767001315f55e7ae87ad9a1626325e5ba1f1e08",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,1068,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-record-journey",
          "command": "bench test --package ./cmd/bench --run '^TestNativeRecordCLIJourney$'",
          "exit_code": 0
        },
        {
          "id": "t2-r4-t2-adoption-smoke",
          "performer": "claude:tcs-t2-repair4-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b604b56d6c000230b4a9b8238003fb8e1f8927ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t2-repair4-20261021",
            "digest": "sha256:c37282e58159cc63e504c0ee1a240c52870f9385768f92c7f5826fa11c1f1fe3",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,12188,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-adoption-smoke",
          "command": "bench test --check system --run '^TestAdoptionSmokeJourney$'",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "c2-standards-r0",
          "performer": "claude:tcs-c2-standards-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "bd4d0e07d53386e9678bb22ae282c164e7949229",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-standards-20261021",
            "digest": "sha256:eb4269a0a6ec290facfc5c4fbdd9b8e3a94d5e141cde9b00de1c8b2a3ccc2ac7",
            "excerpt": "Standards axis, TC-C2: 3 findings, 0 accepted. S1 (auto-fix proposed, confidence 6): outcome_test.go:49-74 golden rows spell the slowest block inline while failure_lines_test.go defines cannedSlowestBlock. Coordinator disposition: no-op; the baseGoldens table spells each full output inline by its own convention (skips[0] repeats seven times there), so this is incidental golden text. S2 (no-op, confidence 6): renderEvents in test_cost_test.go repeats two short decode-then-render steps of testreport_test.go; honest incidental repetition. S3 (no-op, confidence 4): the spec coverage-cell rewrite is outside Writes:; the delegation rule requires the citation to land with its test. Clean: one milliseconds helper, one slowestPerPackage constant, comments state current code.\n"
          },
          "axis": "Standards",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "e214d5647273fc31eb3ac4884d70438c014125f4",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c2-spec-r0",
          "performer": "claude:tcs-c2-spec-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "bd4d0e07d53386e9678bb22ae282c164e7949229",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-spec-20261021",
            "digest": "sha256:c57e328ef7f49617156203830c99731398010ebea828f077d6ab5dc2f03cf907",
            "excerpt": "Spec axis, TC-C2: 0 findings. Rows TC12-TC26 each cite an existing test that asserts the row behavior; TC23 and TC24 sit in the baseGoldens full-equality table of TestCommandKeepsItsBaseOutput. Implementation decisions hold: one math.Round helper used by both cells, last terminal event overwrites finished, skipped, and elapsedMS, one row-limit constant, render order packages, failures, skips, slowest, root_conformance, empty table form slowest[0]{package,test,elapsed_ms}:. No edit outside the acceptance reach.\n"
          },
          "axis": "Spec",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "e214d5647273fc31eb3ac4884d70438c014125f4",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c2-coverage-r0",
          "performer": "claude:tcs-c2-coverage-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "bd4d0e07d53386e9678bb22ae282c164e7949229",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-coverage-20261021",
            "digest": "sha256:296d71c882f11e58af9bb9dde0beb7f43dccda0ed63f5e82004e1dd63066d6ad",
            "excerpt": "Coverage axis, TC-C2: 2 findings, both accepted (auto-fix). C1 (confidence 8): internal/testreport/test_cost.go milliseconds helper; the math.Round decision is unpinned, because every fixture value lands on a whole millisecond; coordinator probe K1 (truncation swap) read silent, restored yes. Fix: a fixture with a sub-millisecond fraction that truncation and rounding render differently. C2 (confidence 6): internal/testreport/testreport.go elapsedMS assignment; \"the last terminal event of a test wins\" is unpinned for elapsed; coordinator probe K2 (first nonzero terminal wins) read silent, restored yes. Fix: a fixture with two terminal events for one test whose later value renders. No other defect: TC12-TC26 cited tests assert their rows; edges TC21 and TC23 asserted.\n"
          },
          "axis": "Coverage",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "e214d5647273fc31eb3ac4884d70438c014125f4",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "c2-standards-r1",
          "performer": "claude:tcs-c2-standards-r1-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "aaf7ffde7e3059c73ac481b2fb34e993a3a7b3f2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-standards-r1-20261021",
            "digest": "sha256:c221a1515f53c78ef643846de5b269b804f5e18e97f9fac843703c9707d8a755",
            "excerpt": "Standards axis, TC-C2 confirming round 1 (repair delta e214d5647..14b320f6a): 0 findings. Both new tests reuse the existing wantSlowest and testTerminal helpers; each pins one fact (0.0026 s renders 3 ms; the last terminal event renders 500 ms); no pasted harness and no duplicated fact.\n"
          },
          "axis": "Standards",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "14b320f6a981e20b6b8b5e7f8a7d5603f2ef23e6",
          "finding_ids": [],
          "supersedes": [
            "c2-standards-r0"
          ]
        },
        {
          "id": "c2-spec-r1",
          "performer": "claude:tcs-c2-spec-r1-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "aaf7ffde7e3059c73ac481b2fb34e993a3a7b3f2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-spec-r1-20261021",
            "digest": "sha256:a7f6be4143023111fbf7fb5c4276242bcaedfdcdaa84e720ef4d2376574dba05",
            "excerpt": "Spec axis, TC-C2 confirming round 1: 0 findings. C1 closed by spec: TestSlowestRoundsToNearestMillisecond feeds 0.0026 s and expects canned,TestA,3, which matches the math.Round decision. C2 closed by spec: TestSlowestLastTerminalEventWins sends fail 0.3 then pass 0.5 and expects 500, which matches \"the last terminal event of a test wins\". The repair contradicts no row TC12-TC26.\n"
          },
          "axis": "Spec",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "14b320f6a981e20b6b8b5e7f8a7d5603f2ef23e6",
          "finding_ids": [],
          "supersedes": [
            "c2-spec-r0"
          ]
        },
        {
          "id": "c2-coverage-r1",
          "performer": "claude:tcs-c2-coverage-r1-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "aaf7ffde7e3059c73ac481b2fb34e993a3a7b3f2",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-coverage-r1-20261021",
            "digest": "sha256:72039b3f5ffdc408827ffd17eb188cc7c815395452112ba6545264a09fab5b7d",
            "excerpt": "Coverage axis, TC-C2 confirming round 1: 1 finding. C3 (auto-fix, confidence 9): internal/testreport/test_cost_test.go TestSlowestLastTerminalEventWins does not close C2, because its last terminal value 0.5 is also the maximum; probe at internal/testreport/testreport.go:101 swapping to a max-of-terminals assignment read silent, restored yes. Fix: make the last terminal event the smaller value. C1 closed for truncation (K1 bit). Advice, not a finding: a ceiling helper still passes 0.0026 -> 3.\n"
          },
          "axis": "Coverage",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "14b320f6a981e20b6b8b5e7f8a7d5603f2ef23e6",
          "finding_ids": [
            "C3"
          ],
          "supersedes": [
            "c2-coverage-r0"
          ]
        },
        {
          "id": "c2-standards-r2",
          "performer": "claude:tcs-c2-standards-r2-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "3b90591aac1d7ac27cb5ac499be470956ed3fe68",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-standards-r2-20261021",
            "digest": "sha256:c45d922576cba174736aa7596d3454d695087833dd52895a84101d922c716674",
            "excerpt": "Standards axis, TC-C2 confirming round 2 (cycle 2 delta 14b320f6a..49eb64960): 0 findings. Both tests reuse wantSlowest and testTerminal; no duplicated knowledge, parser, harness, or registry. The rounding test pairs 0.0026 -> 3 with 0.0024 -> 2; the last-terminal test sends the smaller value last.\n"
          },
          "axis": "Standards",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "49eb64960a455b6ca0443269c644081c07c22806",
          "finding_ids": [],
          "supersedes": [
            "c2-standards-r1"
          ]
        },
        {
          "id": "c2-spec-r2",
          "performer": "claude:tcs-c2-spec-r2-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "3b90591aac1d7ac27cb5ac499be470956ed3fe68",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-spec-r2-20261021",
            "digest": "sha256:62189f0bc02c1d09443586528301eb5ca905f852707342f95c2723293b5742c7",
            "excerpt": "Spec axis, TC-C2 confirming round 2: 0 findings. TestSlowestRoundsToNearestMillisecond (0.0026 -> 3, 0.0024 -> 2) matches math.Round to whole milliseconds; TestSlowestLastTerminalEventWins (pass 0.5 then fail 0.3, expects 300) matches the last terminal event rule, and fail is an admitted last action. No conflict with rows TC12-TC26.\n"
          },
          "axis": "Spec",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "49eb64960a455b6ca0443269c644081c07c22806",
          "finding_ids": [],
          "supersedes": [
            "c2-spec-r1"
          ]
        },
        {
          "id": "c2-coverage-r2",
          "performer": "claude:tcs-c2-coverage-r2-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "3b90591aac1d7ac27cb5ac499be470956ed3fe68",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-coverage-r2-20261021",
            "digest": "sha256:14dc7a804d17702246b8b916ca728a5f4e7f168161ad314f7f3edfa1a0827153",
            "excerpt": "Coverage axis, TC-C2 confirming round 2: 1 finding. C4 (auto-fix, confidence 9): internal/testreport/test_cost_test.go:60-66 TestSlowestLastTerminalEventWins does not kill min-wins, because its one two-event case sends the smaller value last; probe at internal/testreport/testreport.go:101 with a min-of-terminals assignment read silent, restored yes. Fix: add the other order (smaller first, larger last) beside the existing case. The rounding residual is closed: the test separates round, truncate, and ceiling.\n"
          },
          "axis": "Coverage",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "49eb64960a455b6ca0443269c644081c07c22806",
          "finding_ids": [
            "C4"
          ],
          "supersedes": [
            "c2-coverage-r1"
          ]
        },
        {
          "id": "c2-standards-r3",
          "performer": "claude:tcs-c2-standards-r3-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "9be17cb51fd80fb5567656b3ebfd8a142f0213c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-standards-r3-20261021",
            "digest": "sha256:5c329c99fdee1c8d87783b60f4277c23894689eea3202830649af47e6e001aa2",
            "excerpt": "Standards axis, TC-C2 confirming round 3 (cycle 3 delta 49eb64960..849d606f3): 0 findings. The second case of TestSlowestLastTerminalEventWins reuses wantSlowest and testTerminal; the doc comment states the current two cases and holds no history.\n"
          },
          "axis": "Standards",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "849d606f3d57dc5321d847ba8076ebe3620d90a8",
          "finding_ids": [],
          "supersedes": [
            "c2-standards-r2"
          ]
        },
        {
          "id": "c2-spec-r3",
          "performer": "claude:tcs-c2-spec-r3-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "9be17cb51fd80fb5567656b3ebfd8a142f0213c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-spec-r3-20261021",
            "digest": "sha256:2f8ca0380e232f793834a5eb0563763c7d0a03c18edbed14326978e136c17ed8",
            "excerpt": "Spec axis, TC-C2 confirming round 3: 0 findings. The new case (fail 0.2 then pass 0.6, expects canned,TestA,600) matches \"the last terminal event of a test wins\" and the candidate rule; with the first case (300), first-wins, max-wins, and min-wins each go red. No contradiction with rows TC12-TC26.\n"
          },
          "axis": "Spec",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "849d606f3d57dc5321d847ba8076ebe3620d90a8",
          "finding_ids": [],
          "supersedes": [
            "c2-spec-r2"
          ]
        },
        {
          "id": "c2-coverage-r3",
          "performer": "claude:tcs-c2-coverage-r3-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "9be17cb51fd80fb5567656b3ebfd8a142f0213c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-coverage-r3-20261021",
            "digest": "sha256:f2d52a3b9a3b621d95c20fb8d8b4b07e03a65826ab05e3540c49c5cea2167fd8",
            "excerpt": "Coverage axis, TC-C2 confirming round 3: 0 findings. C4 closed (confidence 9): the added case kills min-wins (200 against 600), and the first case still kills first-wins and max-wins; this matches the K2, K4, and K5 probes that bit at this tip. Advice: none.\n"
          },
          "axis": "Coverage",
          "base": "f2a1a693405ee8443c98d7c592d81068b7898919",
          "tip": "849d606f3d57dc5321d847ba8076ebe3620d90a8",
          "finding_ids": [],
          "supersedes": [
            "c2-coverage-r2"
          ]
        },
        {
          "id": "c2-standards-r4",
          "performer": "claude:tcs-c2-standards-r4-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "9be17cb51fd80fb5567656b3ebfd8a142f0213c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-standards-r4-20261021",
            "digest": "sha256:0599e10e7ef04cc0f46895289b5eb3abc90d5baaf004dec8e759681b308055c5",
            "excerpt": "Standards axis, TC-C2 at the moved base 6ba2c0406 (code tip 849d606f3, evidence tip 15bc2bd8d): PASS, 0 findings. The added range 6ba2c0406..f2a1a6934 changes only reviews/test-cost-signal.md (the TC-C1 record and review commits) and one line of the spec completion plan (the ticket 2 assignment entry). No code and no behavior change.\n"
          },
          "axis": "Standards",
          "base": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
          "tip": "849d606f3d57dc5321d847ba8076ebe3620d90a8",
          "finding_ids": [],
          "supersedes": [
            "c2-standards-r3"
          ]
        },
        {
          "id": "c2-spec-r4",
          "performer": "claude:tcs-c2-spec-r4-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "9be17cb51fd80fb5567656b3ebfd8a142f0213c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-spec-r4-20261021",
            "digest": "sha256:d922b47871d4049822db645af12ce394932d6fd448f4133661df155d89327404",
            "excerpt": "Spec axis, TC-C2 at the moved base: PASS, 0 findings. The one spec change in 6ba2c0406..f2a1a6934 is the 2-slowest-tests.md assignment entry at spec.md:161; acceptance rows, tickets, coverage cells, and chunks are unchanged. The refrozen TC-C2 chunk entry names base 6ba2c0406 and tip 849d606f3.\n"
          },
          "axis": "Spec",
          "base": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
          "tip": "849d606f3d57dc5321d847ba8076ebe3620d90a8",
          "finding_ids": [],
          "supersedes": [
            "c2-spec-r3"
          ]
        },
        {
          "id": "c2-coverage-r4",
          "performer": "claude:tcs-c2-coverage-r4-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "9be17cb51fd80fb5567656b3ebfd8a142f0213c2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-coverage-r4-20261021",
            "digest": "sha256:4be746671cc53e697b53696d3e0564d5d6a7a0c3afd2f60cc63922504604d4fd",
            "excerpt": "Coverage axis, TC-C2 at the moved base: PASS, 0 findings. The added range changes no Go file and no tested behavior, so TC12-TC26 coverage is unchanged. The t2-r3 verification entries carry source digest 9be17cb51, which matches the TC-C2 chunk entry and the r3 reviews.\n"
          },
          "axis": "Coverage",
          "base": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
          "tip": "849d606f3d57dc5321d847ba8076ebe3620d90a8",
          "finding_ids": [],
          "supersedes": [
            "c2-coverage-r3"
          ]
        },
        {
          "id": "c2-standards-r5",
          "performer": "claude:tcs-c2-standards-r5-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "b604b56d6c000230b4a9b8238003fb8e1f8927ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-standards-r5-20261021",
            "digest": "sha256:baffa06e9d10c4e279bea5fbab51f655362e5e850097d79a927810aaf650e2e8",
            "excerpt": "Standards axis, TC-C2 round 5 (the --full amendment 09b938caf..4a966b105 and repair 5b51fcb8e..1795d4bbc): 3 findings, none blocking. F1 (auto-fix proposed, confidence 6): test_cost_test.go:147 fullRun and failure_lines_test.go:38 both spell the --full argv. Coordinator disposition: no-op; the base failure_lines_test.go spells that argv inline five times, so it is incidental test text by tree convention. F2 (auto-fix, confidence 5): spec \"Rendered-shape readers\" cites line numbers that the repair moved; fold: name the tests in the next plan commit. F3 (ask-user, confidence 4): spec prose restates the 10, 5, and 12 line values that internal/bounds owns; fold: name the bound constants in the next plan commit. Its second half (TC40 does not name TestAdoptionSmokeJourney) is no-op: requirement t2-adoption-smoke grades it. Clean: one render path selected by full, one block order, a short why comment, TC23 and TC24 cite real tests.\n"
          },
          "axis": "Standards",
          "base": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
          "tip": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
          "finding_ids": [],
          "supersedes": [
            "c2-standards-r4"
          ]
        },
        {
          "id": "c2-spec-r5",
          "performer": "claude:tcs-c2-spec-r5-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "b604b56d6c000230b4a9b8238003fb8e1f8927ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-spec-r5-20261021",
            "digest": "sha256:3dcd48a6b4edf70a6e8c2bbe0a87672c3381d84a7874015b6b799fd2c4f97993",
            "excerpt": "Spec axis, TC-C2 round 5: 1 finding, 0 accepted. Code matches the amended spec: no slowest block on a plain run, the block order on a full run, the empty table on a full run. TC23, TC24, and TC39 cite tests that assert their behavior by full equality. No claim remains that the table renders on every run; ticket 3 and TC-C3 stay consistent. P1 (confidence 5): the TC40 red was not shown in the read set. Coordinator disposition: no-op; the TC-C2 checkpoint gate at 09b938caf failed TestNativeRecordCLIJourney at record_operations_test.go:183 with the spilled 12-line answer, and the cycle 4 author observed that red at the base.\n"
          },
          "axis": "Spec",
          "base": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
          "tip": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
          "finding_ids": [],
          "supersedes": [
            "c2-spec-r4"
          ]
        },
        {
          "id": "c2-coverage-r5",
          "performer": "claude:tcs-c2-coverage-r5-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "b604b56d6c000230b4a9b8238003fb8e1f8927ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c2-coverage-r5-20261021",
            "digest": "sha256:3531db93f49b9e8ee4f3bd1fe5ce1819958035af4024d9433f5a6d50e66586b6",
            "excerpt": "Coverage axis, TC-C2 round 5: 0 findings. Every slowest-row test renders --full (renderEvents calls render(true); TC23 and TC24 use fullRun), so none passes on a plain run for the wrong reason. The plain paths of failure_lines_test.go, tests_run_test.go, and outcome_test.go equal the base bytes. The bench probe --full readers use strings.Contains and pass. bench test --package ./internal/testreport --package ./internal/probe passed (257 and 315 tests). Probes already bit: result += slowestBlock omission, if full -> if true, and the coordinator request.full = false insertion in command.go.\n"
          },
          "axis": "Coverage",
          "base": "6ba2c0406751db07573d8c3929cfa4dff99ab355",
          "tip": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
          "finding_ids": [],
          "supersedes": [
            "c2-coverage-r4"
          ]
        }
      ]
    },
    {
      "id": "TC-C3",
      "base": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
      "tip": "63d97956227bbe55db78767bd7f72bcc8b9a84ac",
      "plan_digest": "sha256:11f891ae58901a2217867f2bcc653246fc6ee06a85da20fee822a3d0ede0ec54",
      "source_digest": "963401f284430475ee64af4ed1c2a5976c1444aa",
      "acceptance_rows": [
        "TC27",
        "TC28",
        "TC29",
        "TC30",
        "TC31",
        "TC32",
        "TC33",
        "TC34",
        "TC35",
        "TC36",
        "TC37"
      ],
      "verification": [
        {
          "id": "t3-testreport",
          "performer": "claude:tcs-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "886346e22913f236e8a7ec8c309b8d8a1d47bdb6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t3-author-20261021",
            "digest": "sha256:f66e299657465985e34f0fbee7b352d0dd9f476dec16718d9ae9a15225a08d75",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,78925,264\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t3-conformance",
          "performer": "claude:tcs-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "886346e22913f236e8a7ec8c309b8d8a1d47bdb6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t3-author-20261021",
            "digest": "sha256:63c964f4da8ec2ab647fe703478a347a6b80cda9b194325d5f0e02c1f60b7773",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,45809,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t3-bounds-policy",
          "performer": "claude:tcs-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "886346e22913f236e8a7ec8c309b8d8a1d47bdb6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t3-author-20261021",
            "digest": "sha256:c115431ec6c29f710d2e73035e91d0b8c86d23827dea08c6ef014642fae4fd49",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  bounds-policy,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,254,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-bounds-policy",
          "command": "bench test --check bounds-policy",
          "exit_code": 0
        },
        {
          "id": "t3-root-conformance",
          "performer": "claude:tcs-t3-author-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "886346e22913f236e8a7ec8c309b8d8a1d47bdb6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t3-author-20261021",
            "digest": "sha256:5e85121e7e75de9f57b6c55079e646019517425d2603c196c18fb1bd6ce3625d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16093,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3-r1-t3-testreport",
          "performer": "claude:tcs-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "963401f284430475ee64af4ed1c2a5976c1444aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t3-repair1-20261021",
            "digest": "sha256:5fb3938566fdb755040906f66c9357079c5b64006f84a89d9d118eec341ae9a4",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,pass,77867,265\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "t3-r1-t3-conformance",
          "performer": "claude:tcs-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "963401f284430475ee64af4ed1c2a5976c1444aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t3-repair1-20261021",
            "digest": "sha256:dacb23d7a4f464649008c12d53e2d27cd8206b261fa606ffdef603674674a27e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,44780,1139\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t3-r1-t3-bounds-policy",
          "performer": "claude:tcs-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "963401f284430475ee64af4ed1c2a5976c1444aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t3-repair1-20261021",
            "digest": "sha256:b4007c6ea12ab07ac0733fa70d50fc44280e6d164444df1a16a9416b98b62fad",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  bounds-policy,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,252,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-bounds-policy",
          "command": "bench test --check bounds-policy",
          "exit_code": 0
        },
        {
          "id": "t3-r1-t3-root-conformance",
          "performer": "claude:tcs-t3-repair1-20261021",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "963401f284430475ee64af4ed1c2a5976c1444aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:tcs-t3-repair1-20261021",
            "digest": "sha256:b1d257555360a5c37607eafb256291efcbbe04a9e2f904f935d5dd284f5d9466",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15838,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "c3-standards-r0",
          "performer": "claude:tcs-c3-standards-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "886346e22913f236e8a7ec8c309b8d8a1d47bdb6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c3-standards-20261021",
            "digest": "sha256:52910e8dd2bb88a775cbdaa4f17c6bfc102f06374ed99a8f4dd4758700b1ef6b",
            "excerpt": "Standards axis, TC-C3: 2 findings, 1 accepted. S1 (no-op proposed, confidence 3): TestTimeoutWarningKeepsPassingExit repeats the setup of TestTimeoutWarningsRenderFirst; coordinator disposition no-op, because the two tests cite different rows (TC32 and TC30). S2 (auto-fix, confidence 3): internal/testreport/timeout_warning_test.go:117-119 the halfTimeoutMS comment ends with a clause that argues for its own red record; trim that clause. Clean: one sorted package list, the bounds-policy entry for test_cost.go, timeless comments in test_cost.go and testreport.go, STE changelog sentences, plan-commit spec corrections.\n"
          },
          "axis": "Standards",
          "base": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
          "tip": "7a554e12e44555f4553631eebc974765463287bc",
          "finding_ids": [
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c3-spec-r0",
          "performer": "claude:tcs-c3-spec-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "886346e22913f236e8a7ec8c309b8d8a1d47bdb6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c3-spec-20261021",
            "digest": "sha256:bf2f5691355a95436fbd058a22f4b8ab56c0ae375fdad46ed8107a2606bb869d",
            "excerpt": "Spec axis, TC-C3: 1 finding, accepted. P1 (auto-fix, confidence 4): CHANGELOG.md:11 \"Before, the limit came from the environment\" is imprecise; the old limit was the implicit Go default of ten minutes, and a GOFLAGS value could override it. Rows TC27-TC37 pass: threshold >= with the divisor constant, one millisecond below renders nothing, package order (three packages make the order check stronger), the block first, the check row first under --check, exits 0 and 1, healthy run unchanged, changelog names the flag, the linked-repository effect, slowest of --full, and timeout_warnings, no literal timeout in any test except the TC38 pin, and the bounds-policy entry.\n"
          },
          "axis": "Spec",
          "base": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
          "tip": "7a554e12e44555f4553631eebc974765463287bc",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "c3-coverage-r0",
          "performer": "claude:tcs-c3-coverage-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "886346e22913f236e8a7ec8c309b8d8a1d47bdb6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c3-coverage-20261021",
            "digest": "sha256:dfb9e6ebb5c992774b643dc3724f550df732cd55bc4df31595fd55f765653515",
            "excerpt": "Coverage axis, TC-C3: 1 finding, accepted. C1 (auto-fix, confidence 9): internal/testreport/timeout_warning_test.go fixtures all sit exactly at the threshold, so a timeout_ms cell derived from the package time (elapsed * divisor) passes; coordinator probe K6 at internal/testreport/test_cost.go read silent, restored yes, and this axis reran it with the same verdict. Fix: a warned fixture above the threshold whose row shows its own elapsed_ms and the constant timeout_ms, and unequal elapsed values in the order test. TC27-TC33 and TC37 otherwise cite tests that assert their rows; TC36 holds.\n"
          },
          "axis": "Coverage",
          "base": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
          "tip": "7a554e12e44555f4553631eebc974765463287bc",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "c3-standards-r1",
          "performer": "claude:tcs-c3-standards-r1-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "963401f284430475ee64af4ed1c2a5976c1444aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c3-standards-r1-20261021",
            "digest": "sha256:d872e64f1e14cbb3192ca6e08aeedf7e64dce0f0866abb5d9225081db257ee55",
            "excerpt": "Standards axis, TC-C3 confirming round 1 (repair 7a554e12e..63d979562): 0 findings. S2 closed: the halfTimeoutMS comment states what and why in one line. The warnedPackage type and atThreshold helper give each package's elapsed time one source for the events and the expected rows. The changelog sentence is plain STE.\n"
          },
          "axis": "Standards",
          "base": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
          "tip": "63d97956227bbe55db78767bd7f72bcc8b9a84ac",
          "finding_ids": [],
          "supersedes": [
            "c3-standards-r0"
          ]
        },
        {
          "id": "c3-spec-r1",
          "performer": "claude:tcs-c3-spec-r1-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "963401f284430475ee64af4ed1c2a5976c1444aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c3-spec-r1-20261021",
            "digest": "sha256:7ab1a27b6e64a56664b4f9fdead4b566c4b17bb3eae6359e4a497c57945cda29",
            "excerpt": "Spec axis, TC-C3 confirming round 1: 0 findings. P1 closed: CHANGELOG.md:10 states that the limit was the implicit Go default of ten minutes, which a -timeout in GOFLAGS could change. The changed tests still assert TC27-TC33 as the spec states them; TestTimeoutWarningAboveThreshold is additive.\n"
          },
          "axis": "Spec",
          "base": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
          "tip": "63d97956227bbe55db78767bd7f72bcc8b9a84ac",
          "finding_ids": [],
          "supersedes": [
            "c3-spec-r0"
          ]
        },
        {
          "id": "c3-coverage-r1",
          "performer": "claude:tcs-c3-coverage-r1-20261021",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "medium",
          "source_digest": "963401f284430475ee64af4ed1c2a5976c1444aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:tcs-c3-coverage-r1-20261021",
            "digest": "sha256:8c0f6e684be5ba562bfe4502379efbc8825dd272fa8cb9ae000d26a77746cf01",
            "excerpt": "Coverage axis, TC-C3 confirming round 1: 0 findings. C1 closed: each warned row takes its elapsed cell from its own fixture value, TestTimeoutWarningAboveThreshold uses the threshold plus one, and the order test uses unequal values, so K6 and K7 bite. TC36 holds. Advice, not a finding: the test-side /2 in halfTimeoutMS is an independent copy of the production divisor; the author's divisor 2->3 mutant bit against it.\n"
          },
          "axis": "Coverage",
          "base": "1795d4bbc0fafb61c150c65264fabd843fc46d3c",
          "tip": "63d97956227bbe55db78767bd7f72bcc8b9a84ac",
          "finding_ids": [],
          "supersedes": [
            "c3-coverage-r0"
          ]
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
      "from": "sha256:a030a302362c213d596f53bf8de12a525abc5ee5de38c4ee5d1129e3aaaf1407",
      "to": "sha256:321f4c6febbf4d188806e04290462c6ad9c3af218c9a23a4c831f8fce926e34d",
      "chunk_ids": {
        "TC-C1": [
          "TC-C1"
        ]
      }
    },
    {
      "from": "sha256:321f4c6febbf4d188806e04290462c6ad9c3af218c9a23a4c831f8fce926e34d",
      "to": "sha256:3b77567c43768aebdd895a498a8cfe407310f037da223b9e1c3448e195a089b2",
      "chunk_ids": {
        "TC-C1": [
          "TC-C1"
        ],
        "TC-C2": [
          "TC-C2"
        ]
      }
    },
    {
      "from": "sha256:3b77567c43768aebdd895a498a8cfe407310f037da223b9e1c3448e195a089b2",
      "to": "sha256:5d79c8012ea0883ed2194af97f5c65fd7f5e75d471531192d2e9a485bc892ee0",
      "chunk_ids": {
        "TC-C1": [
          "TC-C1"
        ],
        "TC-C2": [
          "TC-C2"
        ]
      }
    },
    {
      "from": "sha256:5d79c8012ea0883ed2194af97f5c65fd7f5e75d471531192d2e9a485bc892ee0",
      "to": "sha256:5c96c02f9ab8638951e0718696430d4ff0b135fc44b59b2a5c51f53cb4ed1b75",
      "chunk_ids": {
        "TC-C1": [
          "TC-C1"
        ],
        "TC-C2": [
          "TC-C2"
        ]
      }
    },
    {
      "from": "sha256:5c96c02f9ab8638951e0718696430d4ff0b135fc44b59b2a5c51f53cb4ed1b75",
      "to": "sha256:10fc6009bec677b4a1c75ff2ea3342ce81eeab5262feb7b6d04db55342f4b834",
      "chunk_ids": {
        "TC-C1": [
          "TC-C1"
        ],
        "TC-C2": [
          "TC-C2"
        ]
      }
    },
    {
      "from": "sha256:10fc6009bec677b4a1c75ff2ea3342ce81eeab5262feb7b6d04db55342f4b834",
      "to": "sha256:073860a0a02894a62a0f2a8cde55044b870e1588ff27390dc6320ca13f44ecb7",
      "chunk_ids": {
        "TC-C1": [
          "TC-C1"
        ],
        "TC-C2": [
          "TC-C2"
        ]
      }
    },
    {
      "from": "sha256:073860a0a02894a62a0f2a8cde55044b870e1588ff27390dc6320ca13f44ecb7",
      "to": "sha256:7235841b50c687080fb3e991ceeb10d9137a671a9f6274b9a190b5a6038941df",
      "chunk_ids": {
        "TC-C1": [
          "TC-C1"
        ],
        "TC-C2": [
          "TC-C2"
        ]
      }
    },
    {
      "from": "sha256:7235841b50c687080fb3e991ceeb10d9137a671a9f6274b9a190b5a6038941df",
      "to": "sha256:11f891ae58901a2217867f2bcc653246fc6ee06a85da20fee822a3d0ede0ec54",
      "chunk_ids": {
        "TC-C1": [
          "TC-C1"
        ],
        "TC-C2": [
          "TC-C2"
        ],
        "TC-C3": [
          "TC-C3"
        ]
      }
    }
  ]
}
```
