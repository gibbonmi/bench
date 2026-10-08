# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/refusal-route-registry/spec.md",
  "plan_digest": "sha256:09259ca04468131aff9125326c043078881ec4190cafa38628483388cb98a239",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RR-C1a",
      "base": "9c228393356ae35e4f940c2071d12e36d42ede7a",
      "tip": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
      "plan_digest": "sha256:10d938cdde11624c65c40c632beebb04b986d194da031b1c8dd8d470d65aaa05",
      "source_digest": "e516146765693ca135ae558a2f01d88b654f6a99",
      "acceptance_rows": [
        "RR04",
        "RR05",
        "RR06",
        "RR07",
        "RR13"
      ],
      "verification": [
        {
          "id": "v-t1-registry",
          "performer": "claude:ft393_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a27426ae564319fe1c7c3febe1c7b31b25c7bb5d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t1",
            "digest": "sha256:dce953d229559745058b1e2c37a956ac6f8b8cc6f8eedb8d402e331c0cd6a4cf",
            "excerpt": "$ bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,da295794f164bc5fd459a5ebc0a6ae5ac7e8b4ef,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,13\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t1-renderer-proof",
          "performer": "claude:ft393_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a27426ae564319fe1c7c3febe1c7b31b25c7bb5d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t1",
            "digest": "sha256:d3062b7047d945b75d3b9c9ec5b9b42e42dbcc0f7b22e543e7855a6ffaa561a8",
            "excerpt": "$ bench test --package ./internal/refusalroute --run 'TestRouteRendering'\ntree[1]{target,head,dirty}:\n  ft393-build,da295794f164bc5fd459a5ebc0a6ae5ac7e8b4ef,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,10\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/route.go --swap 'stepJoiner     = \"; then \"' --with 'stepJoiner     = \", \"' --package ./internal/refusalroute --run 'TestRouteRendering'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/route.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/refusalroute,TestRouteRendering,passed,10\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/a_composed_command_that_is_not_line-safe_prints_its_placeholder,\"route_test.go:75: rendered route:\",3\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/a_preface_comes_first_after_the_reviewer_marker,\"route_test.go:75: rendered route:\",3\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/an_agent_route_joins_its_steps_in_declared_order,\"route_test.go:75: rendered route:\",3\n"
          },
          "requirement": "t1-renderer-proof",
          "command": "bench test --package ./internal/refusalroute --run 'TestRouteRendering'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make the renderer join the route steps with `, ` instead of `; then `. TestRouteRendering must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t1",
              "digest": "sha256:d3062b7047d945b75d3b9c9ec5b9b42e42dbcc0f7b22e543e7855a6ffaa561a8",
              "excerpt": "$ bench test --package ./internal/refusalroute --run 'TestRouteRendering'\ntree[1]{target,head,dirty}:\n  ft393-build,da295794f164bc5fd459a5ebc0a6ae5ac7e8b4ef,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,10\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/route.go --swap 'stepJoiner     = \"; then \"' --with 'stepJoiner     = \", \"' --package ./internal/refusalroute --run 'TestRouteRendering'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/route.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/refusalroute,TestRouteRendering,passed,10\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/a_composed_command_that_is_not_line-safe_prints_its_placeholder,\"route_test.go:75: rendered route:\",3\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/a_preface_comes_first_after_the_reviewer_marker,\"route_test.go:75: rendered route:\",3\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/an_agent_route_joins_its_steps_in_declared_order,\"route_test.go:75: rendered route:\",3\n"
            }
          }
        },
        {
          "id": "v-t1-registry-r1",
          "performer": "claude:ft393_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e516146765693ca135ae558a2f01d88b654f6a99",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t1_repair1",
            "digest": "sha256:b09a7a220b5aae87ea8253e4e81362b3db06315223b586c44292362e760ad3bb",
            "excerpt": "$ bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,5e6272170c81293be97bea4fa5b551b1d5dd6455,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,15\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t1-renderer-proof-r1",
          "performer": "claude:ft393_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e516146765693ca135ae558a2f01d88b654f6a99",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t1_repair1",
            "digest": "sha256:ce0b93d7e4fe3c43803d167869abe4def8745a9cf0a32216c935a5e934180351",
            "excerpt": "$ bench test --package ./internal/refusalroute --run 'TestRouteRendering'\ntree[1]{target,head,dirty}:\n  ft393-build,5e6272170c81293be97bea4fa5b551b1d5dd6455,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,12\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/route.go --swap 'stepJoiner     = \"; then \"' --with 'stepJoiner     = \", \"' --package ./internal/refusalroute --run 'TestRouteRendering'\ntree[1]{target,head,dirty}:\n  ft393-build,5e6272170c81293be97bea4fa5b551b1d5dd6455,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/route.go,swap,failed,4,yes\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/a_composed_command_that_is_not_line-safe_prints_its_placeholder,\"route_test.go:86: rendered route:\",3\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/a_preface_comes_first_after_the_reviewer_marker,\"route_test.go:86: rendered route:\",3\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/an_absent_composed_command_prints_its_placeholder,\"route_test.go:86: rendered route:\",3\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/an_agent_route_joins_its_steps_in_declared_order,\"route_test.go:86: rendered route:\",3\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-renderer-proof",
          "command": "bench test --package ./internal/refusalroute --run 'TestRouteRendering'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make the renderer join the route steps with `, ` instead of `; then `. TestRouteRendering must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t1_repair1",
              "digest": "sha256:ce0b93d7e4fe3c43803d167869abe4def8745a9cf0a32216c935a5e934180351",
              "excerpt": "$ bench test --package ./internal/refusalroute --run 'TestRouteRendering'\ntree[1]{target,head,dirty}:\n  ft393-build,5e6272170c81293be97bea4fa5b551b1d5dd6455,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,12\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/route.go --swap 'stepJoiner     = \"; then \"' --with 'stepJoiner     = \", \"' --package ./internal/refusalroute --run 'TestRouteRendering'\ntree[1]{target,head,dirty}:\n  ft393-build,5e6272170c81293be97bea4fa5b551b1d5dd6455,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/route.go,swap,failed,4,yes\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/a_composed_command_that_is_not_line-safe_prints_its_placeholder,\"route_test.go:86: rendered route:\",3\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/a_preface_comes_first_after_the_reviewer_marker,\"route_test.go:86: rendered route:\",3\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/an_absent_composed_command_prints_its_placeholder,\"route_test.go:86: rendered route:\",3\n  github.com/gibbonmi/bench/internal/refusalroute,TestRouteRendering/an_agent_route_joins_its_steps_in_declared_order,\"route_test.go:86: rendered route:\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "r-c1a-standards",
          "performer": "claude:ft393_c1a_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a27426ae564319fe1c7c3febe1c7b31b25c7bb5d",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c1a_standards",
            "digest": "sha256:2831767803c2f52efaca917ac539da11d3dfd58f84de944017b8722bab9b7077",
            "excerpt": "C1a Standards axis (claude:ft393_c1a_standards): evidence current=true. 4 findings, 0 blocking.\nC1a-S1 advisory route.go Step.command field set but unread; ticket 03 guard check owns the reader. no-op. conf 6\nC1a-S2 advisory registry.go uniqueNames is called only by tests; follows the approved design. no-op. conf 4\nC1a-S3 advisory registry_test.go/route_test.go comments cite acceptance rows; repo precedent. no-op. conf 4\nC1a-S4 advisory registry.go newIn has no doc comment naming the injectable-inventory seam. auto-fix. conf 3\nWorst: C1a-S1. Implementation command contributed to no finding.\n"
          },
          "axis": "Standards",
          "base": "9c228393356ae35e4f940c2071d12e36d42ede7a",
          "tip": "da295794f164bc5fd459a5ebc0a6ae5ac7e8b4ef",
          "finding_ids": [
            "C1a-S1",
            "C1a-S2",
            "C1a-S3",
            "C1a-S4"
          ],
          "supersedes": []
        },
        {
          "id": "r-c1a-spec",
          "performer": "claude:ft393_c1a_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a27426ae564319fe1c7c3febe1c7b31b25c7bb5d",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c1a_spec",
            "digest": "sha256:058ddcb6555d40054c3a24b248e1e7f1a3f2d22144ad6ceb387c5ce6c2d2c46f",
            "excerpt": "C1a Spec axis (claude:ft393_c1a_spec): evidence current=true. RR04, RR05, RR06, RR07, RR13 met. 3 findings, 0 blocking.\nC1a-P1 advisory route.go Operator/Operators never read Values; the guard check sample fill is ticket 03's seam. no-op. conf 6\nC1a-P2 advisory registry.go unregistered fallback discards the observed sentence; matches spec RR07. no-op. conf 3\nC1a-P3 advisory route.go Facts.Preface joined raw with no LineSafe gate, unlike Fact and Composed. auto-fix. conf 3\nWorst: C1a-P1. Implementation command contributed to no finding.\n"
          },
          "axis": "Spec",
          "base": "9c228393356ae35e4f940c2071d12e36d42ede7a",
          "tip": "da295794f164bc5fd459a5ebc0a6ae5ac7e8b4ef",
          "finding_ids": [
            "C1a-P1",
            "C1a-P2",
            "C1a-P3"
          ],
          "supersedes": []
        },
        {
          "id": "r-c1a-coverage",
          "performer": "claude:ft393_c1a_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a27426ae564319fe1c7c3febe1c7b31b25c7bb5d",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c1a_coverage",
            "digest": "sha256:3cb016ac8142fd780fdb6e5422aaa83fa52ad7ded8cc79142fd250888c7db837",
            "excerpt": "C1a Coverage axis (claude:ft393_c1a_coverage): evidence current=true. 4 findings, 0 blocking.\nC1a-C1 advisory route.go wordComposed empty value untested; dropping value != \"\" survives. auto-fix. conf 7\nC1a-C2 advisory route.go Operator never supplied a Values entry in a test; ticket 03 owns the sample-fill seam. no-op. conf 5\nC1a-C3 advisory route.go Facts.Preface concatenated with no LineSafe check. auto-fix. conf 4\nC1a-C4 advisory registry.go empty-name and empty-sentence degenerate states untested and unspecified. no-op. conf 3\nWorst: C1a-C1. Implementation command contributed to no finding.\n"
          },
          "axis": "Coverage",
          "base": "9c228393356ae35e4f940c2071d12e36d42ede7a",
          "tip": "da295794f164bc5fd459a5ebc0a6ae5ac7e8b4ef",
          "finding_ids": [
            "C1a-C1",
            "C1a-C2",
            "C1a-C3",
            "C1a-C4"
          ],
          "supersedes": []
        },
        {
          "id": "r-c1a-r1-standards",
          "performer": "claude:ft393_c1a_r1_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e516146765693ca135ae558a2f01d88b654f6a99",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c1a_r1_standards",
            "digest": "sha256:fec0cd4287f193d87225b38580e4841107ef6323fba907160a2cec1774f96d27",
            "excerpt": "RR-C1a confirming round, Standards (claude:ft393_c1a_r1_standards): evidence current=true at 5537e2e7. Folds C1a-P3/C1a-C3, C1a-C1, C1a-S4 confirmed. New findings: zero.\nOptional advice (no id): the Fact case spells the same line-safe predicate that asWritten holds; prefaceSlot sits under the prefaceJoiner comment.\nImplementation command contribution: none.\n"
          },
          "axis": "Standards",
          "base": "9c228393356ae35e4f940c2071d12e36d42ede7a",
          "tip": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "finding_ids": [],
          "supersedes": [
            "r-c1a-standards"
          ]
        },
        {
          "id": "r-c1a-r1-spec",
          "performer": "claude:ft393_c1a_r1_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e516146765693ca135ae558a2f01d88b654f6a99",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c1a_r1_spec",
            "digest": "sha256:3b468f5d86e24293381b752848257c80106517e937769720ee0d754f2900b8cd",
            "excerpt": "RR-C1a confirming round, Spec (claude:ft393_c1a_r1_spec): evidence current=true at 5537e2e7. Folds C1a-P3/C1a-C3, C1a-C1, C1a-S4 confirmed. RR04-RR07 and RR13 unchanged and met; the preface gate obeys the Edge inventory. New findings: zero.\nImplementation command contribution: none.\n"
          },
          "axis": "Spec",
          "base": "9c228393356ae35e4f940c2071d12e36d42ede7a",
          "tip": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "finding_ids": [],
          "supersedes": [
            "r-c1a-spec"
          ]
        },
        {
          "id": "r-c1a-r1-coverage",
          "performer": "claude:ft393_c1a_r1_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e516146765693ca135ae558a2f01d88b654f6a99",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c1a_r1_coverage",
            "digest": "sha256:ade6d084520d6013ae88a1316b48e9bf8d102b0d70c1e2efe1e894cbd4023ccc",
            "excerpt": "RR-C1a confirming round, Coverage (claude:ft393_c1a_r1_coverage): evidence current=true at 5537e2e7. Folds C1a-P3/C1a-C3, C1a-C1, C1a-S4 confirmed; each new case fails on its targeted mutation. New findings: zero.\nOptional advice (no id): a whitespace-only preface renders as written.\nImplementation command contribution: none.\n"
          },
          "axis": "Coverage",
          "base": "9c228393356ae35e4f940c2071d12e36d42ede7a",
          "tip": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "finding_ids": [],
          "supersedes": [
            "r-c1a-coverage"
          ]
        }
      ]
    },
    {
      "id": "RR-C1b",
      "base": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
      "tip": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
      "plan_digest": "sha256:f0639cf27495f50d19843b401bf84eab059aa614e05f76f81c98dc485790924b",
      "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
      "acceptance_rows": [
        "RR01",
        "RR02",
        "RR03",
        "RR14",
        "RR16",
        "RR17",
        "RR18",
        "RR19",
        "RR20",
        "RR08",
        "RR09",
        "RR10",
        "RR11",
        "RR12",
        "RR15",
        "RR58",
        "RR60",
        "RR51",
        "RR52",
        "RR53",
        "RR54"
      ],
      "verification": [
        {
          "id": "v-t2-registry",
          "performer": "claude:ft393_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2",
            "digest": "sha256:cd703ce3ea2316d4824af2e088ed8e1e7c83e92e1f72c33ed50891d67442fbbe",
            "excerpt": "$ bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t2-worktree-package",
          "performer": "claude:ft393_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2",
            "digest": "sha256:3fe9622575e0b35898cb7d57c7ed1472eec1b5650539b627efd1cb188ee3496d",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,67211,1326\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"clean_landed_hostile_test.go:98: unix sockets unavailable: listen unix ... (280 bytes)\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable: listen unix ... (279 bytes)\"\n"
          },
          "requirement": "t2-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t2-landing-faces",
          "performer": "claude:ft393_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2",
            "digest": "sha256:45041dc847c19e104043832b43ca2f895653858076b0f46ca1299a70c3f68c86",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestConflictRepairIsAReviewerRoute|TestLandCommandReportsEveryRefusalInOnePreflight'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1178,16\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestConflictRepairIsAReviewerRoute|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "exit_code": 0
        },
        {
          "id": "v-t2-conflict-proof",
          "performer": "claude:ft393_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2",
            "digest": "sha256:d048929fc93e6109cd84fc61fc5f5ccbfa3a3517aad9bf7f11df381e825e3b53",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,239,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nProbe: Declare the `composition-conflict` face with the agent authority.\n$ bench probe internal/refusalroute/registry.go --swap $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Reviewer,' --with $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Agent,' --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,228,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:100: composition-conflict next = \\\"git -C '...' merge ... (1020 bytes)\",2\n"
          },
          "requirement": "t2-conflict-proof",
          "command": "bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'",
          "exit_code": 0,
          "probe": {
            "mutation": "Declare the `composition-conflict` face with the agent authority. TestConflictRepairIsAReviewerRoute must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t2",
              "digest": "sha256:d048929fc93e6109cd84fc61fc5f5ccbfa3a3517aad9bf7f11df381e825e3b53",
              "excerpt": "$ bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,239,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nProbe: Declare the `composition-conflict` face with the agent authority.\n$ bench probe internal/refusalroute/registry.go --swap $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Reviewer,' --with $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Agent,' --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,228,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:100: composition-conflict next = \\\"git -C '...' merge ... (1020 bytes)\",2\n"
            }
          }
        },
        {
          "id": "v-t3-registry",
          "performer": "claude:ft393_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3",
            "digest": "sha256:9160e5eba21851b24fde2ede477f2394042fc9c0397f02c47b5ba128df3db00e",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,3,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t3-guard-check",
          "performer": "claude:ft393_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3",
            "digest": "sha256:08b24dd6660403e8160e28bee8beb9e91f8953b936d3cf9417009d144a01c672",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-guard-check",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'",
          "exit_code": 0
        },
        {
          "id": "v-t3-landing-faces",
          "performer": "claude:ft393_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3",
            "digest": "sha256:fafb516d496d28d3763c069d5cf5ac5a7944027ecd5669b925d58112b07d1cd9",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,3684,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "exit_code": 0
        },
        {
          "id": "v-t3-guard-proof",
          "performer": "claude:ft393_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3",
            "digest": "sha256:b31d03fea8ca015a68d0b9161f8c02571ae790b1125ef7de4b5605fa56d40f5b",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/registry.go --swap 'TreeCommand(\"bench commit\", Fact(FactLabel), Text(\"-m\"), Operator(\"msg\"), Text(\"--\"), Operators(\"path\")),' --with 'Command(Text(\"git merge\"), Operator(\"commit\")),' --package ./internal/conformance --run TestAgentRoutesPassTheWiredGuards\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestAgentRoutesPassTheWiredGuards,\"refusal_route_guard_test.go:111: face source-not-clean step 1 \\\"git merge '/bench-home/worktrees/pool/repository/assignment'\\\": gitguard denies git merge\",1\n"
          },
          "requirement": "t3-guard-proof",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'",
          "exit_code": 0,
          "probe": {
            "mutation": "Declare the first route step of the `source-not-clean` face as the command `git merge <commit>`. TestAgentRoutesPassTheWiredGuards must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t3",
              "digest": "sha256:b31d03fea8ca015a68d0b9161f8c02571ae790b1125ef7de4b5605fa56d40f5b",
              "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/registry.go --swap 'TreeCommand(\"bench commit\", Fact(FactLabel), Text(\"-m\"), Operator(\"msg\"), Text(\"--\"), Operators(\"path\")),' --with 'Command(Text(\"git merge\"), Operator(\"commit\")),' --package ./internal/conformance --run TestAgentRoutesPassTheWiredGuards\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestAgentRoutesPassTheWiredGuards,\"refusal_route_guard_test.go:111: face source-not-clean step 1 \\\"git merge '/bench-home/worktrees/pool/repository/assignment'\\\": gitguard denies git merge\",1\n"
            }
          }
        },
        {
          "id": "v-t4-registry",
          "performer": "claude:ft393_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4",
            "digest": "sha256:c120ecc6e8ea7faca24759db727af69bb17873bb155a6f3565ebb754bc7a3e2a",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t4-guard-check",
          "performer": "claude:ft393_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4",
            "digest": "sha256:08b24dd6660403e8160e28bee8beb9e91f8953b936d3cf9417009d144a01c672",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,38,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-guard-check",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'",
          "exit_code": 0
        },
        {
          "id": "v-t4-landing-faces",
          "performer": "claude:ft393_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4",
            "digest": "sha256:9352434e0cb00e75fbdb07084e79faf245736c528df5798ea0b3453d845221d5",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,4033,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "exit_code": 0
        },
        {
          "id": "v-t4-worktree-package",
          "performer": "claude:ft393_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4",
            "digest": "sha256:5c6194723ea01893f62fa5f3a17bb9f7142f8512c65ccb47e5cc2d90121e6f3a",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,66395,1326\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"capability: fifo: unix sockets unavailable: listen unix ... (265 bytes)\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable: listen unix ... (278 bytes)\"\n"
          },
          "requirement": "t4-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t4-recovery-verb",
          "performer": "claude:ft393_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4",
            "digest": "sha256:08db695d3f22ecb4d866344071f6826db6801919412b8726c03479658b720a40",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,4,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-recovery-verb",
          "command": "bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'",
          "exit_code": 0
        },
        {
          "id": "v-t4-recovery-proof",
          "performer": "claude:ft393_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4",
            "digest": "sha256:f2c64fd5bd7efe30157bcddf5b9d677942ba634757fdb59bd56cdadf15dffbcb",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/command.go --swap 'for _, face := range inventory {' --with 'for _, face := range inventory[:len(inventory)-1] {' --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/command.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/refusalroute,TestRecoveryListsEveryFace,\"command_test.go:25: stdout = \\\"recovery[10]{verb,face,authority,route}:... (1020 bytes)\",2\n"
          },
          "requirement": "t4-recovery-proof",
          "command": "bench test --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make `bench recovery` skip the last registered face. TestRecoveryListsEveryFace must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t4",
              "digest": "sha256:f2c64fd5bd7efe30157bcddf5b9d677942ba634757fdb59bd56cdadf15dffbcb",
              "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/command.go --swap 'for _, face := range inventory {' --with 'for _, face := range inventory[:len(inventory)-1] {' --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\ntree[1]{target,head,dirty}:\n  ft393-build,389cb337d7f198bb2925785dbfa1ab6e74e4a805,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/command.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/refusalroute,TestRecoveryListsEveryFace,\"command_test.go:25: stdout = \\\"recovery[10]{verb,face,authority,route}:... (1020 bytes)\",2\n"
            }
          }
        },
        {
          "id": "v-t2-registry-r1",
          "performer": "claude:ft393_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2_repair1",
            "digest": "sha256:03dd473e4de705847ae5f49516da3ea0f51cc027cc584f1f425636be30350e3b",
            "excerpt": "$ bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t2-worktree-package-r1",
          "performer": "claude:ft393_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2_repair1",
            "digest": "sha256:69939afeecaa19f43529d5148627f4d3beb5ffcc6fd48a5af598bf1e95f85e6e",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,69845,1338\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"unix sockets unavailable\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable\"\n"
          },
          "requirement": "t2-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t2-landing-faces-r1",
          "performer": "claude:ft393_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2_repair1",
            "digest": "sha256:d6564935640631091062af6f2d3b0cae6ad308edcc53e3289523a43c54de0f30",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestConflictRepairIsAReviewerRoute|TestLandCommandReportsEveryRefusalInOnePreflight'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1128,16\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestConflictRepairIsAReviewerRoute|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "exit_code": 0
        },
        {
          "id": "v-t2-conflict-proof-r1",
          "performer": "claude:ft393_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2_repair1",
            "digest": "sha256:519db6d84a87c2c2a5b095ed5c3d45d16e30c721609eecb8d41f33b7ae83bcaf",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,239,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nNamed probe: Declare the `composition-conflict` face with the agent authority\n$ bench probe internal/refusalroute/registry.go --swap $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Reviewer,' --with $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Agent,' --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,239,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:121: composition-conflict next = \\\"git -C '...' merge ...\\\" (no reviewer: prefix)\",2\n"
          },
          "requirement": "t2-conflict-proof",
          "command": "bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'",
          "exit_code": 0,
          "probe": {
            "mutation": "Declare the `composition-conflict` face with the agent authority. TestConflictRepairIsAReviewerRoute must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t2_repair1",
              "digest": "sha256:519db6d84a87c2c2a5b095ed5c3d45d16e30c721609eecb8d41f33b7ae83bcaf",
              "excerpt": "$ bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,239,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nNamed probe: Declare the `composition-conflict` face with the agent authority\n$ bench probe internal/refusalroute/registry.go --swap $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Reviewer,' --with $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Agent,' --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,239,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:121: composition-conflict next = \\\"git -C '...' merge ...\\\" (no reviewer: prefix)\",2\n"
            }
          }
        },
        {
          "id": "v-t3-registry-r1",
          "performer": "claude:ft393_t3_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3_repair1",
            "digest": "sha256:03dd473e4de705847ae5f49516da3ea0f51cc027cc584f1f425636be30350e3b",
            "excerpt": "$ bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t3-guard-check-r1",
          "performer": "claude:ft393_t3_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3_repair1",
            "digest": "sha256:f32b0549704ec4ad8619c0712325d6f4d5b1544d3f1a8a19202f0d1a67ad1359",
            "excerpt": "$ bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,53,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-guard-check",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'",
          "exit_code": 0
        },
        {
          "id": "v-t3-landing-faces-r1",
          "performer": "claude:ft393_t3_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3_repair1",
            "digest": "sha256:52894168297aa7c408448774d339a058ac91e1b23dc6fe0ff79a31c9ef9740d0",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,4183,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "exit_code": 0
        },
        {
          "id": "v-t3-worktree-package-r1",
          "performer": "claude:ft393_t3_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3_repair1",
            "digest": "sha256:a89ea6219dfc95919dccdf0553f45a7f167e767b4ecc280e3843a99c807e5f19",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,68975,1338\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"clean_landed_hostile_test.go:98: unix sockets unavailable: listen unix ... (281 bytes)\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable: listen unix ... (278 bytes)\"\n"
          },
          "requirement": "t3-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t3-guard-proof-r1",
          "performer": "claude:ft393_t3_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3_repair1",
            "digest": "sha256:e67e477a5888baa815b8c32740d0e6d939b3474458e656f0a18e412b0b3e4c50",
            "excerpt": "$ bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/registry.go --swap $'TreeCommand(\"bench commit\", Fact(FactLabel), Text(\"-m\"), Operator(\"msg\"), Text(\"--\"), Operators(\"path\")),\\n\\t\\t\\treview,' --with $'Command(Text(\"git merge\"), Fact(FactConflictCommit)),\\n\\t\\t\\treview,' --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestAgentRoutesPassTheWiredGuards,\"refusal_route_guard_test.go:111: face source-not-clean step 1 \\\"git merge '/bench-home/worktrees/pool/repository/assignment'\\\": gitguard denies git merge\",1\n"
          },
          "requirement": "t3-guard-proof",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'",
          "exit_code": 0,
          "probe": {
            "mutation": "Declare the first route step of the `source-not-clean` face as the command `git merge <commit>`. TestAgentRoutesPassTheWiredGuards must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t3_repair1",
              "digest": "sha256:e67e477a5888baa815b8c32740d0e6d939b3474458e656f0a18e412b0b3e4c50",
              "excerpt": "$ bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/registry.go --swap $'TreeCommand(\"bench commit\", Fact(FactLabel), Text(\"-m\"), Operator(\"msg\"), Text(\"--\"), Operators(\"path\")),\\n\\t\\t\\treview,' --with $'Command(Text(\"git merge\"), Fact(FactConflictCommit)),\\n\\t\\t\\treview,' --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestAgentRoutesPassTheWiredGuards,\"refusal_route_guard_test.go:111: face source-not-clean step 1 \\\"git merge '/bench-home/worktrees/pool/repository/assignment'\\\": gitguard denies git merge\",1\n"
            }
          }
        },
        {
          "id": "v-t4-registry-r1",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:03dd473e4de705847ae5f49516da3ea0f51cc027cc584f1f425636be30350e3b",
            "excerpt": "$ bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t4-guard-check-r1",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:0555a29d83d10d58c8818ecc4161ad58ddff9d5e20bbf1b8117bf54cff1fdd24",
            "excerpt": "$ bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,54,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-guard-check",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'",
          "exit_code": 0
        },
        {
          "id": "v-t4-landing-faces-r1",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:587a11b38f6654357749668a349d125250c9e9b88cbfd1e2dd5c0d5502b54bab",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,4203,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "exit_code": 0
        },
        {
          "id": "v-t4-worktree-package-r1",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:27d1d52890f259021db90ab543093b6b788c9deb04dff36f19b988a0121f8e0d",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,68036,1338\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,unix sockets unavailable (capability skip)\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,unix sockets unavailable (capability skip)\n"
          },
          "requirement": "t4-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t4-recovery-verb-r1",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:973da3dcbea1396797b31d74f45c9133413969905161fd0cb6534839008f4c08",
            "excerpt": "$ bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,4,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-recovery-verb",
          "command": "bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'",
          "exit_code": 0
        },
        {
          "id": "v-t4-recovery-proof-r1",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:653516b1f7c3e3f45702066a1928084735bd3ac774be49d9d77f06373ec03026",
            "excerpt": "$ bench test --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/command.go --swap 'for _, face := range inventory {' --with 'for _, face := range inventory[:len(inventory)-1] {' --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/command.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/refusalroute,TestRecoveryListsEveryFace,\"command_test.go:25: stdout = \\\"recovery[10]{verb,face,authority,route}:...\"\n"
          },
          "requirement": "t4-recovery-proof",
          "command": "bench test --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make `bench recovery` skip the last registered face. TestRecoveryListsEveryFace must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t4_repair1",
              "digest": "sha256:653516b1f7c3e3f45702066a1928084735bd3ac774be49d9d77f06373ec03026",
              "excerpt": "$ bench test --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\ntree[1]{target,head,dirty}:\n  ft393-build,f4669147db2ddf388fc9fdae06ba52f428870537,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/command.go --swap 'for _, face := range inventory {' --with 'for _, face := range inventory[:len(inventory)-1] {' --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/command.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/refusalroute,TestRecoveryListsEveryFace,\"command_test.go:25: stdout = \\\"recovery[10]{verb,face,authority,route}:...\"\n"
            }
          }
        },
        {
          "id": "v-t2-registry-r2",
          "performer": "claude:ft393_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2_repair1",
            "digest": "sha256:28d043fefb26ec06148fd32940c836c86dd6859ee4d0a93b1111232c877b4bc1",
            "excerpt": "$ bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t2-worktree-package-r2",
          "performer": "claude:ft393_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2_repair1",
            "digest": "sha256:5edc757e50cba72adae8084d6a840c08d10eda7bedb841494da18e518c74ec02",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,71990,1338\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"unix sockets unavailable\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable\"\n"
          },
          "requirement": "t2-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t2-landing-faces-r2",
          "performer": "claude:ft393_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2_repair1",
            "digest": "sha256:964285a31ebc56e30dba02441c25d8921d2c99755665fb85a6c14582e6302156",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestConflictRepairIsAReviewerRoute|TestLandCommandReportsEveryRefusalInOnePreflight'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1219,16\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestConflictRepairIsAReviewerRoute|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "exit_code": 0
        },
        {
          "id": "v-t2-conflict-proof-r2",
          "performer": "claude:ft393_t2_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t2_repair1",
            "digest": "sha256:9a8a8051b3a4c1e60caaf83d4afdaca00385a1c590038598a3c7c24c872d471d",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,252,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nNamed probe: Declare the `composition-conflict` face with the agent authority\n$ bench probe internal/refusalroute/registry.go --swap $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Reviewer,' --with $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Agent,' --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,227,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:121: composition-conflict next = \\\"git -C '...' merge ...\\\" (no reviewer: prefix)\",2\n"
          },
          "requirement": "t2-conflict-proof",
          "command": "bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'",
          "exit_code": 0,
          "probe": {
            "mutation": "Declare the `composition-conflict` face with the agent authority. TestConflictRepairIsAReviewerRoute must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t2_repair1",
              "digest": "sha256:9a8a8051b3a4c1e60caaf83d4afdaca00385a1c590038598a3c7c24c872d471d",
              "excerpt": "$ bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,252,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nNamed probe: Declare the `composition-conflict` face with the agent authority\n$ bench probe internal/refusalroute/registry.go --swap $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Reviewer,' --with $'Name:      \"composition-conflict\",\\n\\t\\tAuthority: Agent,' --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,227,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/composition-conflict,\"refusal_route_test.go:121: composition-conflict next = \\\"git -C '...' merge ...\\\" (no reviewer: prefix)\",2\n"
            }
          }
        },
        {
          "id": "v-t3-registry-r2",
          "performer": "claude:ft393_t3_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3_repair2",
            "digest": "sha256:df2bf9948515f183e258dda292dd59da209b91972058bc52c769f4eb48e937f9",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t3-guard-check-r2",
          "performer": "claude:ft393_t3_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3_repair2",
            "digest": "sha256:0d5697edd868b5fb705a9d2423f4a306ddc2433440b3ef2a5e4d07ceabce7126",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,46,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-guard-check",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'",
          "exit_code": 0
        },
        {
          "id": "v-t3-landing-faces-r2",
          "performer": "claude:ft393_t3_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3_repair2",
            "digest": "sha256:8ea755784e4ef3112be131ce35ce296a974b345bb2d1a99ed24995f38a4b4577",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,4078,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "exit_code": 0
        },
        {
          "id": "v-t3-worktree-package-r2",
          "performer": "claude:ft393_t3_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3_repair2",
            "digest": "sha256:620d37c88cff1a143c2f0acd43d067945d886cd1f0f2bbd357fb8db34a3c8dd0",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,72958,1338\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"capability: fifo: unix sockets unavailable\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable\"\n"
          },
          "requirement": "t3-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t3-guard-proof-r2",
          "performer": "claude:ft393_t3_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t3_repair2",
            "digest": "sha256:4818123977490acef992eaa7b6f04587bdd9e332bcc88b9a99faaa4720821e23",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/registry.go --swap 'TreeCommand(\"bench commit\", Fact(FactLabel), Text(\"-m\"), Operator(\"msg\"), Text(\"--\"), Operators(\"path\")),' --with 'Command(Text(\"git merge\"), Operator(\"commit\")),' --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestAgentRoutesPassTheWiredGuards,\"refusal_route_guard_test.go:111: face source-not-clean step 1 \\\"git merge '/bench-home/worktrees/pool/repository/assignment'\\\": gitguard denies git merge\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-guard-proof",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'",
          "exit_code": 0,
          "probe": {
            "mutation": "Declare the first route step of the `source-not-clean` face as the command `git merge <commit>`. TestAgentRoutesPassTheWiredGuards must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t3_repair2",
              "digest": "sha256:4818123977490acef992eaa7b6f04587bdd9e332bcc88b9a99faaa4720821e23",
              "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/registry.go --swap 'TreeCommand(\"bench commit\", Fact(FactLabel), Text(\"-m\"), Operator(\"msg\"), Text(\"--\"), Operators(\"path\")),' --with 'Command(Text(\"git merge\"), Operator(\"commit\")),' --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/registry.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestAgentRoutesPassTheWiredGuards,\"refusal_route_guard_test.go:111: face source-not-clean step 1 \\\"git merge '/bench-home/worktrees/pool/repository/assignment'\\\": gitguard denies git merge\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "v-t4-registry-r2",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:28d043fefb26ec06148fd32940c836c86dd6859ee4d0a93b1111232c877b4bc1",
            "excerpt": "$ bench test --package ./internal/refusalroute\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-registry",
          "command": "bench test --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "v-t4-guard-check-r2",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:7edf1d97d59e37faa1df8797c801789af39e723a34059d1ccbda3de5a483549a",
            "excerpt": "$ bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,45,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-guard-check",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'",
          "exit_code": 0
        },
        {
          "id": "v-t4-landing-faces-r2",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:453554aadaf5176011ec2564e6c15ec6e68bfe640848b50dcfbc87876d84a8c6",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,4494,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "exit_code": 0
        },
        {
          "id": "v-t4-worktree-package-r2",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:7c6a460e06627bcd401276f3313c29b57e411d17806904108c82bbd000ec5032",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,72501,1338\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,unix sockets unavailable (capability skip)\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,unix sockets unavailable (capability skip)\n"
          },
          "requirement": "t4-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t4-recovery-verb-r2",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:07205f456576fbce3a8ebf66d14262e5145b89f5bc8a9dd51a5fab1d40b8d553",
            "excerpt": "$ bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,5,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-recovery-verb",
          "command": "bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'",
          "exit_code": 0
        },
        {
          "id": "v-t4-recovery-proof-r2",
          "performer": "claude:ft393_t4_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t4_repair1",
            "digest": "sha256:101ec67bea2f06e815360ada25d0038ccb2ffa3e79cee7f269a24c8b33c99c31",
            "excerpt": "$ bench test --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/command.go --swap 'for _, face := range inventory {' --with 'for _, face := range inventory[:len(inventory)-1] {' --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/command.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/refusalroute,TestRecoveryListsEveryFace,\"command_test.go:25: stdout = \\\"recovery[10]{verb,face,authority,route}:...\"\n"
          },
          "requirement": "t4-recovery-proof",
          "command": "bench test --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make `bench recovery` skip the last registered face. TestRecoveryListsEveryFace must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t4_repair1",
              "digest": "sha256:101ec67bea2f06e815360ada25d0038ccb2ffa3e79cee7f269a24c8b33c99c31",
              "excerpt": "$ bench test --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\ntree[1]{target,head,dirty}:\n  ft393-build,33ac2f98902ac4a51bfa951218f50f3fd268319d,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/refusalroute,pass,2,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/command.go --swap 'for _, face := range inventory {' --with 'for _, face := range inventory[:len(inventory)-1] {' --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/command.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/refusalroute,TestRecoveryListsEveryFace,\"command_test.go:25: stdout = \\\"recovery[10]{verb,face,authority,route}:...\"\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "r-c1b-standards",
          "performer": "claude:ft393_c1b_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c1b_standards",
            "digest": "sha256:57fa0059cd8466b9d318bb373a5b835cd9b835d79b96a988286809f6e33e8fe5",
            "excerpt": "RR-C1b Standards (claude:ft393_c1b_standards): evidence current=true. 6 findings, all advisory.\nC1b-S1 registry.go and refusal_route_test.go comments cite FT341 as provenance. auto-fix. conf 4\nC1b-S2 worktree tests re-spell the reviewer marker and the step joiner with no recorded red. ask-user. conf 5\nC1b-S3 refusal_route_test.go over-long comment line and the unclear name placeholds. auto-fix. conf 6\nC1b-S4 three ad-hoc parsers of the printed record layout in the worktree route tests. auto-fix. conf 4\nC1b-S5 landingFaceRefusal takes two adjacent same-type strings. auto-fix. conf 4\nC1b-S6 bench recovery takes the operational AXI exemption. ask-user. conf 3\nWorst: C1b-S2. Implementation command contributed to C1b-S6 and C1b-S2.\n"
          },
          "axis": "Standards",
          "base": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "tip": "389cb337d7f198bb2925785dbfa1ab6e74e4a805",
          "finding_ids": [
            "C1b-S1",
            "C1b-S2",
            "C1b-S3",
            "C1b-S4",
            "C1b-S5",
            "C1b-S6"
          ],
          "supersedes": []
        },
        {
          "id": "r-c1b-spec",
          "performer": "claude:ft393_c1b_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c1b_spec",
            "digest": "sha256:74ab85de428a59f0987d8535c50a02dbec7078e8d0538c9a89b3703159c1a8ac",
            "excerpt": "RR-C1b Spec (claude:ft393_c1b_spec): evidence current=true. 3 findings, all advisory.\nC1b-P1 five identity-component landing refusals print no next= route. ask-user. conf 6\nC1b-P2 non-conflict landReviewed errors at land.go print no route and no ticket owns the remainder. ask-user. conf 5\nC1b-P3 RR18 has no independent reviewer-marker assertion. auto-fix. conf 7\nNo-op notes: the AXI exemption and the helpCommand move.\nWorst: C1b-P1. Implementation command contributed to C1b-P1 and C1b-P2.\n"
          },
          "axis": "Spec",
          "base": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "tip": "389cb337d7f198bb2925785dbfa1ab6e74e4a805",
          "finding_ids": [
            "C1b-P1",
            "C1b-P2",
            "C1b-P3"
          ],
          "supersedes": []
        },
        {
          "id": "r-c1b-coverage",
          "performer": "claude:ft393_c1b_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a902a55e90d2758014cc49eaeb59e60f9ce7493f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c1b_coverage",
            "digest": "sha256:e26683cd9360529be5d43ea442d9e20ab2a46f5a0a8d54e4446522b692396b6b",
            "excerpt": "RR-C1b Coverage (claude:ft393_c1b_coverage): evidence current=true. 7 findings, 2 blocking.\nC1b-C1 blocking the unsafe-path bench worktree path step refuses for the same input class. ask-user. conf 8\nC1b-C2 blocking routeless identity-component refusals and land.go landReviewed errors. ask-user. conf 7\nC1b-C3 advisory the follow walk passes when a different refusal replaces the face. ask-user. conf 7\nC1b-C4 advisory the review step is a Command, which blocks the ticket 05 land-red fixture. ask-user. conf 6\nC1b-C5 advisory composed slots take the sample path in the guard check. no-op. conf 5\nC1b-C6 advisory the new handback sites have no next= assertion. auto-fix. conf 5\nC1b-C7 advisory no test dispatches bench recovery through cmd/bench. auto-fix. conf 5\nWorst: C1b-C1. Implementation command contributed to no finding.\n"
          },
          "axis": "Coverage",
          "base": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "tip": "389cb337d7f198bb2925785dbfa1ab6e74e4a805",
          "finding_ids": [
            "C1b-C1",
            "C1b-C2",
            "C1b-C3",
            "C1b-C4",
            "C1b-C5",
            "C1b-C6",
            "C1b-C7"
          ],
          "supersedes": []
        },
        {
          "id": "r-c1b-r1-standards",
          "performer": "claude:ft393_c1b_r1_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c1b_r1_standards",
            "digest": "sha256:9f64bbc77a1e3d1e5633113ba40b7eeaee1af7f679b0aacd2b60267e95163020",
            "excerpt": "RR-C1b confirming round, Standards (claude:ft393_c1b_r1_standards): evidence current=true. Every fold confirmed. 1 advisory finding.\nC1b-RS1 land_resume_refusal_test.go spells the resume rerun inline beside the new resumeRerunOf helper. auto-fix. conf 6\nOptional advice: the recovery header repeats in a cmd/bench test; the repaired-tip face set has a production and a fixture expectation.\nImplementation command contribution: none.\n"
          },
          "axis": "Standards",
          "base": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "tip": "f4669147db2ddf388fc9fdae06ba52f428870537",
          "finding_ids": [
            "C1b-RS1"
          ],
          "supersedes": [
            "r-c1b-standards"
          ]
        },
        {
          "id": "r-c1b-r1-spec",
          "performer": "claude:ft393_c1b_r1_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c1b_r1_spec",
            "digest": "sha256:a20b87ab0dc08656594a026d6b725a7648e5ed8f312b60455b820bcffc8fa4df",
            "excerpt": "RR-C1b confirming round, Spec (claude:ft393_c1b_r1_spec): evidence current=true. Every fold confirmed except the recorded probe of C1b-P3. Rows RR01-RR03, RR08-RR20, RR51-RR54, RR58, RR60 stay met. The amended rows match the code.\nC1b-RP1 blocking the record shows no red for TestReviewerLandFacesOpenWithTheMarker. auto-fix. conf 6\nOptional advice: the seam cells of RR15, RR18, RR19, RR58, RR60 name the old test file.\nImplementation command contribution: none.\n"
          },
          "axis": "Spec",
          "base": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "tip": "f4669147db2ddf388fc9fdae06ba52f428870537",
          "finding_ids": [
            "C1b-RP1"
          ],
          "supersedes": [
            "r-c1b-spec"
          ]
        },
        {
          "id": "r-c1b-r1-coverage",
          "performer": "claude:ft393_c1b_r1_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "bebc4c2d296d0cb8f97e1236eaeaab4c94900f3b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c1b_r1_coverage",
            "digest": "sha256:deb735e970f753da7e20d2f4b255ec96eae602d269099fed2c5b9a9f863bf1bc",
            "excerpt": "RR-C1b confirming round, Coverage (claude:ft393_c1b_r1_coverage): evidence current=true. Every fold confirmed; the follow walk can no longer pass on a replacement refusal; the flake fix holds. 1 advisory finding.\nC1b-RC1 advisory the record does not show the post-fix red of the reviewer-marker test. auto-fix. conf 5\nOptional advice: the RR18 seam cell; resume source-not-fenced keeps its agent route.\nImplementation command contribution: none.\n"
          },
          "axis": "Coverage",
          "base": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "tip": "f4669147db2ddf388fc9fdae06ba52f428870537",
          "finding_ids": [
            "C1b-RC1"
          ],
          "supersedes": [
            "r-c1b-coverage"
          ]
        },
        {
          "id": "r-c1b-r2-standards",
          "performer": "claude:ft393_c1b_r2_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c1b_r2_standards",
            "digest": "sha256:052dd58446e465418fc627f8eeab3f0e954559b7e6c349ea6f06958e24849055",
            "excerpt": "RR-C1b repair cycle 2 confirming round, Standards (claude:ft393_c1b_r2_standards): evidence current=true at d85e78e9. C1b-RS1 confirmed; C1b-RP1 and C1b-RC1 confirmed closed by the record correction. New findings: zero.\nImplementation command contribution: none.\nOptional advice: the resume form also appears in land_journey_test.go and an internal/systemtest file, both outside this delta.\n"
          },
          "axis": "Standards",
          "base": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "tip": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
          "finding_ids": [],
          "supersedes": [
            "r-c1b-r1-standards"
          ]
        },
        {
          "id": "r-c1b-r2-spec",
          "performer": "claude:ft393_c1b_r2_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c1b_r2_spec",
            "digest": "sha256:3b480e8b6544bb8a817e7a92594c9bdc9e01a51b9c0dbd995fa470b17902e5df",
            "excerpt": "RR-C1b repair cycle 2 confirming round, Spec (claude:ft393_c1b_r2_spec): evidence current=true at d85e78e9. C1b-RS1 confirmed; C1b-RP1 and C1b-RC1 confirmed closed by the record correction. New findings: zero.\nImplementation command contribution: none.\nSpec also confirmed the corrected seam cells of RR15, RR18, RR19, RR58, RR60.\n"
          },
          "axis": "Spec",
          "base": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "tip": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
          "finding_ids": [],
          "supersedes": [
            "r-c1b-r1-spec"
          ]
        },
        {
          "id": "r-c1b-r2-coverage",
          "performer": "claude:ft393_c1b_r2_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9ac008b53fa935a61f3420a225b1bd754995a7d1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c1b_r2_coverage",
            "digest": "sha256:72c9b1efb8f27056be2d109e7b0f9f6a72ae8756dc8e7f066f1a7ef38c90fef7",
            "excerpt": "RR-C1b repair cycle 2 confirming round, Coverage (claude:ft393_c1b_r2_coverage): evidence current=true at d85e78e9. C1b-RS1 confirmed; C1b-RP1 and C1b-RC1 confirmed closed by the record correction. New findings: zero.\nImplementation command contribution: none.\n"
          },
          "axis": "Coverage",
          "base": "5e6272170c81293be97bea4fa5b551b1d5dd6455",
          "tip": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
          "finding_ids": [],
          "supersedes": [
            "r-c1b-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "RR-C2",
      "base": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
      "tip": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
      "plan_digest": "sha256:b777d31d70d7bd260d746692d440ae30a28d3d90ba994f94d264d015d0f79a27",
      "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
      "acceptance_rows": [
        "RR27",
        "RR28",
        "RR30",
        "RR31",
        "RR67",
        "RR21",
        "RR22",
        "RR23",
        "RR24",
        "RR25",
        "RR26",
        "RR55",
        "RR56",
        "RR57",
        "RR59",
        "RR29"
      ],
      "verification": [
        {
          "id": "v-t5-merge-routes",
          "performer": "claude:ft393_t5b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t5b",
            "digest": "sha256:7ab8b523e202b8147cadba1019817970d4e6625f4ac50d9d19e79d44520d5cd2",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal|TestLandingRedRouteNamesTheRepair'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1373,27\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-merge-routes",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal|TestLandingRedRouteNamesTheRepair'",
          "exit_code": 0
        },
        {
          "id": "v-t5-worktree-package",
          "performer": "claude:ft393_t5b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t5b",
            "digest": "sha256:64ee435dd5a93068b1860fdec834884e02c5050978f06f2193f7d24d567f3443",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,76384,1368\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"capability: fifo: unix sockets unavailable: ...\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable: ...\"\n"
          },
          "requirement": "t5-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t5-landing-package",
          "performer": "claude:ft393_t5b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t5b",
            "digest": "sha256:4e3a543ffc68cfecc86ed7ac819c1e59f1dec71b4fe81905490e511ef5fda02a",
            "excerpt": "$ bench test --package ./internal/landing\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,8898,188\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t5-landing-package",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "v-t5-conflict-proof",
          "performer": "claude:ft393_t5b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t5b",
            "digest": "sha256:ba301f9b849b87172ef49ee9c6f3c0f78cb1db6a0b48948b405ba86e289071ec",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,236,4\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nprobe: Declare the `merge-conflict` face with the agent authority\n$ bench probe internal/refusalroute/faces_merge.go --swap $'\"merge-conflict\",\\n\\t\\tAuthority: Reviewer,' --with $'\"merge-conflict\",\\n\\t\\tAuthority: Agent,' --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_merge.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/merge-conflict,\"refusal_route_test.go:132: merge-conflict next = \\\"git -C '...' me… (1020 bytes)\",2\n"
          },
          "requirement": "t5-conflict-proof",
          "command": "bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'",
          "exit_code": 0,
          "probe": {
            "mutation": "Declare the `merge-conflict` face with the agent authority. TestConflictRepairIsAReviewerRoute must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t5b",
              "digest": "sha256:ba301f9b849b87172ef49ee9c6f3c0f78cb1db6a0b48948b405ba86e289071ec",
              "excerpt": "$ bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,236,4\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nprobe: Declare the `merge-conflict` face with the agent authority\n$ bench probe internal/refusalroute/faces_merge.go --swap $'\"merge-conflict\",\\n\\t\\tAuthority: Reviewer,' --with $'\"merge-conflict\",\\n\\t\\tAuthority: Agent,' --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_merge.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/merge-conflict,\"refusal_route_test.go:132: merge-conflict next = \\\"git -C '...' me… (1020 bytes)\",2\n"
            }
          }
        },
        {
          "id": "v-t6-merge-routes",
          "performer": "claude:ft393_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t6",
            "digest": "sha256:662fb97c57e7c93c2d889f229bfd6ba5f5d73543562b6de31a7507a86546a5d8",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestRedSourceFoldNamesAnExit|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1831,29\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-merge-routes",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestRedSourceFoldNamesAnExit|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal'",
          "exit_code": 0
        },
        {
          "id": "v-t6-worktree-package",
          "performer": "claude:ft393_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t6",
            "digest": "sha256:b14a96ddb4136c90a00d3773c1b9098fb68de4fe37edc6de535127d80935447c",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,76350,1368\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"capability: fifo: unix sockets unavailable\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable\"\n"
          },
          "requirement": "t6-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t6-landing-package",
          "performer": "claude:ft393_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t6",
            "digest": "sha256:0fffc4a37124b2f94c25f8692fe8f884a8d8be584949f8b932173ba024606a6b",
            "excerpt": "$ bench test --package ./internal/landing\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,8830,188\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t6-landing-package",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "v-t6-target-alone-proof",
          "performer": "claude:ft393_t6",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t6",
            "digest": "sha256:0436cefdc1c2544e6c0e597b73deacd5c74a29571665815b5adc814a4612d92f",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,669,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench probe internal/worktree/merge_refusal.go --swap 'if refused.Result.Kind != authorization.Candidate {' --with 'if refused.Result.Kind == authorization.LaneFail {' --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/merge_refusal.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRedSourceFoldNamesAnExit/inherited,\"merge_route_test.go:206: red-source fold = (1, refused{detail=prospective authorization refused: inherited ...,next=reviewer: the fold of 'main' adds…)\",1\n"
          },
          "requirement": "t6-target-alone-proof",
          "command": "bench test --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make an `inherited` fold take the `merge-fold-red` face without the target-alone grade. TestRedSourceFoldNamesAnExit must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t6",
              "digest": "sha256:0436cefdc1c2544e6c0e597b73deacd5c74a29571665815b5adc814a4612d92f",
              "excerpt": "$ bench test --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,669,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench probe internal/worktree/merge_refusal.go --swap 'if refused.Result.Kind != authorization.Candidate {' --with 'if refused.Result.Kind == authorization.LaneFail {' --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/merge_refusal.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRedSourceFoldNamesAnExit/inherited,\"merge_route_test.go:206: red-source fold = (1, refused{detail=prospective authorization refused: inherited ...,next=reviewer: the fold of 'main' adds…)\",1\n"
            }
          }
        },
        {
          "id": "v-t7-merge-routes",
          "performer": "claude:ft393_t7b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t7b",
            "digest": "sha256:3ef7aa0c18bbe6e17ad7358ec14d1c7e193a3f0ca49555872e6e12c918986bf6",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestRedSourceFoldNamesAnExit|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1722,29\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t7-merge-routes",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestRedSourceFoldNamesAnExit|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal'",
          "exit_code": 0
        },
        {
          "id": "v-t7-worktree-package",
          "performer": "claude:ft393_t7b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t7b",
            "digest": "sha256:e4128beab462af0ac973fda79472d497e3b2c0b4fed44d7cd10907422d82f2ca",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,75585,1368\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"capability: fifo: unix sockets unavailable: listen unix ... (266 bytes)\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable: listen unix ... (279 bytes)\"\n"
          },
          "requirement": "t7-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t7-landing-package",
          "performer": "claude:ft393_t7b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t7b",
            "digest": "sha256:83c0277bc415b2c99e04de6cf90ab1afa3fb8f7b8f595c256e8c93b207a8eee9",
            "excerpt": "$ bench test --package ./internal/landing\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,8344,188\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t7-landing-package",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "v-t7-walk-proof",
          "performer": "claude:ft393_t7b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t7b",
            "digest": "sha256:c4ff0ee511c3ddbfbcd0012a426704d97a1209f2ab6cf390fd31428286a79de8",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1628,15\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/faces_reset.go --swap \"var resetFaces = []Face{\" --with \"var resetFaces = []Face{{Verb: Reset, Name: \\\"reset-unproduced\\\", Sentence: \\\"unproduced reset\\\", Authority: Reviewer, Route: handback},\" --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_reset.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,1329,15\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestMergeFacesFollowTheirRoutes,\"merge_route_test.go:96: registry reset face \\\"reset-unproduced\\\" has no producing fixture\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t7-walk-proof",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'",
          "exit_code": 0,
          "probe": {
            "mutation": "Register one more reset face that no fixture produces. TestMergeFacesFollowTheirRoutes must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t7b",
              "digest": "sha256:c4ff0ee511c3ddbfbcd0012a426704d97a1209f2ab6cf390fd31428286a79de8",
              "excerpt": "$ bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1628,15\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/faces_reset.go --swap \"var resetFaces = []Face{\" --with \"var resetFaces = []Face{{Verb: Reset, Name: \\\"reset-unproduced\\\", Sentence: \\\"unproduced reset\\\", Authority: Reviewer, Route: handback},\" --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,b91c696ce74e984b6cdd93d6d8758b35b61be62c,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_reset.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,1329,15\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestMergeFacesFollowTheirRoutes,\"merge_route_test.go:96: registry reset face \\\"reset-unproduced\\\" has no producing fixture\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "v-t5-merge-routes-r1",
          "performer": "claude:ft393_t5_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t5_repair1",
            "digest": "sha256:3d91d8b0bd60d1bbbc4a5541ac0ff7cc81f651df04c99fada2c4ee944ec6d59b",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal|TestLandingRedRouteNamesTheRepair'\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1901,34\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t5-merge-routes",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal|TestLandingRedRouteNamesTheRepair'",
          "exit_code": 0
        },
        {
          "id": "v-t5-worktree-package-r1",
          "performer": "claude:ft393_t5_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t5_repair1",
            "digest": "sha256:8f716df8f98aff9c38c61fa1ce5fa7b3b373500494b27b08471b5eaec137f2d7",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,79598,1382\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"capability: fifo: unix sockets unavailable: listen unix ... (266 bytes)\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable: listen unix ... (278 bytes)\"\n"
          },
          "requirement": "t5-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t5-landing-package-r1",
          "performer": "claude:ft393_t5_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t5_repair1",
            "digest": "sha256:4a37a98498769995da68afe63cf1b9c96fa0f8283439465af1eef588537db4c4",
            "excerpt": "$ bench test --package ./internal/landing\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,8703,188\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t5-landing-package",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "v-t5-conflict-proof-r1",
          "performer": "claude:ft393_t5_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t5_repair1",
            "digest": "sha256:1db8774b018f41873d10953098bce1dbcc4e3d05be9691dfff18be36f0f0e051",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,218,4\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nprobe: Declare the `merge-conflict` face with the agent authority\n$ bench probe internal/refusalroute/faces_merge.go --swap $'\"merge-conflict\",\\n\\t\\tAuthority: Reviewer,' --with $'\"merge-conflict\",\\n\\t\\tAuthority: Agent,' --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_merge.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/merge-conflict,\"refusal_route_test.go:132: merge-conflict next = \\\"git -C '...' merg… (1020 bytes)\",2\n"
          },
          "requirement": "t5-conflict-proof",
          "command": "bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'",
          "exit_code": 0,
          "probe": {
            "mutation": "Declare the `merge-conflict` face with the agent authority. TestConflictRepairIsAReviewerRoute must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t5_repair1",
              "digest": "sha256:1db8774b018f41873d10953098bce1dbcc4e3d05be9691dfff18be36f0f0e051",
              "excerpt": "$ bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,218,4\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nprobe: Declare the `merge-conflict` face with the agent authority\n$ bench probe internal/refusalroute/faces_merge.go --swap $'\"merge-conflict\",\\n\\t\\tAuthority: Reviewer,' --with $'\"merge-conflict\",\\n\\t\\tAuthority: Agent,' --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_merge.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestConflictRepairIsAReviewerRoute/merge-conflict,\"refusal_route_test.go:132: merge-conflict next = \\\"git -C '...' merg… (1020 bytes)\",2\n"
            }
          }
        },
        {
          "id": "v-t6-merge-routes-r1",
          "performer": "claude:ft393_t6_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t6_repair1",
            "digest": "sha256:289d0abca5f979159a2050eb67a8baf317eeef6646dd6ce2eb8b78efc4471b73",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestRedSourceFoldNamesAnExit|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal'\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,2573,36\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t6-merge-routes",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestRedSourceFoldNamesAnExit|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal'",
          "exit_code": 0
        },
        {
          "id": "v-t6-worktree-package-r1",
          "performer": "claude:ft393_t6_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t6_repair1",
            "digest": "sha256:1259b2aa1540e5406939becf2002d474f2fc3c09548012694db0d6270e6c25a1",
            "excerpt": "$ bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,80561,1382\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"clean_landed_hostile_test.go:98: unix sockets unavailable ...\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable ...\"\n"
          },
          "requirement": "t6-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t6-landing-package-r1",
          "performer": "claude:ft393_t6_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t6_repair1",
            "digest": "sha256:0b9738bb5dd3e0b52455e04806a4428bd82de0c7b13e3049e601e531a02c64cf",
            "excerpt": "$ bench test --package ./internal/landing\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,8507,188\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t6-landing-package",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "v-t6-target-alone-proof-r1",
          "performer": "claude:ft393_t6_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t6_repair1",
            "digest": "sha256:29fc5f351125ad17032de2da3e295aec1a354ee3406c6941fdc10203cef7ac5e",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,587,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench probe internal/worktree/merge_refusal.go --swap 'if refused.Result.Kind != authorization.Candidate {' --with 'if refused.Result.Kind == authorization.LaneFail {' --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/merge_refusal.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRedSourceFoldNamesAnExit/inherited,\"merge_route_test.go:219: red-source fold = (1, \\\"refused{detail=prospective authorization refused: inherited ...,next=reviewer: the fold of 'main' adds…)\",1\n"
          },
          "requirement": "t6-target-alone-proof",
          "command": "bench test --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'",
          "exit_code": 0,
          "probe": {
            "mutation": "Make an `inherited` fold take the `merge-fold-red` face without the target-alone grade. TestRedSourceFoldNamesAnExit must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t6_repair1",
              "digest": "sha256:29fc5f351125ad17032de2da3e295aec1a354ee3406c6941fdc10203cef7ac5e",
              "excerpt": "$ bench test --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,587,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench probe internal/worktree/merge_refusal.go --swap 'if refused.Result.Kind != authorization.Candidate {' --with 'if refused.Result.Kind == authorization.LaneFail {' --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/merge_refusal.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRedSourceFoldNamesAnExit/inherited,\"merge_route_test.go:219: red-source fold = (1, \\\"refused{detail=prospective authorization refused: inherited ...,next=reviewer: the fold of 'main' adds…)\",1\n"
            }
          }
        },
        {
          "id": "v-t7-merge-routes-r1",
          "performer": "claude:ft393_t7_repair1b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t7_repair1b",
            "digest": "sha256:01948eff47a1636d8ac9434615c80f314a15899dbe0a8cd19c599e0d61e5ea3b",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestRedSourceFoldNamesAnExit|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal'\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,2095,36\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t7-merge-routes",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestRedSourceFoldNamesAnExit|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal'",
          "exit_code": 0
        },
        {
          "id": "v-t7-worktree-package-r1",
          "performer": "claude:ft393_t7_repair1b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t7_repair1b",
            "digest": "sha256:d1e371f0c1b3d69d5d4ee3661a202a5a6c6b9bd35fbd62a8a18f74ba2db0d43b",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,79453,1382\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"clean_landed_hostile_test.go:98: unix sockets unavailable (capability skip)\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable\"\n"
          },
          "requirement": "t7-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "v-t7-landing-package-r1",
          "performer": "claude:ft393_t7_repair1b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t7_repair1b",
            "digest": "sha256:1dbdb2c7d7c76cbd70043ca1c2d97b73ff5c137de75c69bcb1a660c5a26488de",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/landing\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,8994,188\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t7-landing-package",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "v-t7-walk-proof-r1",
          "performer": "claude:ft393_t7_repair1b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t7_repair1b",
            "digest": "sha256:190462bcf6316e6bc279d23fccab089cac8f865f10c88de6767703e5758f5432",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,2451,22\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_reset.go --swap 'Name:      \"reset-handback\",' --with 'Name: \"reset-unproduced\", Sentence: \"unproduced reset face\", Authority: Agent, Route: []Step{Command(Text(\"bench worktree list\"))}}, {Verb: Reset, Name: \"reset-handback\",' --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_reset.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestMergeFacesFollowTheirRoutes,\"merge_route_test.go:96: registry reset face \\\"reset-unproduced\\\" has no producing fixture\",1\n"
          },
          "requirement": "t7-walk-proof",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'",
          "exit_code": 0,
          "probe": {
            "mutation": "Register one more reset face that no fixture produces. TestMergeFacesFollowTheirRoutes must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t7_repair1b",
              "digest": "sha256:190462bcf6316e6bc279d23fccab089cac8f865f10c88de6767703e5758f5432",
              "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,fba4fe785c6c55f8ad7a2cc1522efa33730332cb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,2451,22\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_reset.go --swap 'Name:      \"reset-handback\",' --with 'Name: \"reset-unproduced\", Sentence: \"unproduced reset face\", Authority: Agent, Route: []Step{Command(Text(\"bench worktree list\"))}}, {Verb: Reset, Name: \"reset-handback\",' --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_reset.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestMergeFacesFollowTheirRoutes,\"merge_route_test.go:96: registry reset face \\\"reset-unproduced\\\" has no producing fixture\",1\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "r-c2-standards",
          "performer": "claude:ft393_c2_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c2_standards",
            "digest": "sha256:9c8522d85174530edb669eb014160b2fbffe266a4c33791a0ed02d4de85fe96a",
            "excerpt": "RR-C2 Standards (claude:ft393_c2_standards): evidence current=true. 6 advisory findings.\nC2-S1 the hand-merge fixture step is pasted three times. auto-fix. conf 7\nC2-S2 the mergeRedRefusal comment uses the parameter name as a verb. auto-fix. conf 6\nC2-S3 an unwrapped comment line and a missing blank line. auto-fix. conf 8\nC2-S4 planLandedAssignment does not state its landed precondition. auto-fix. conf 4\nC2-S5 releaseAssignment builds missingTreeRelease twice. no-op or ask-user. conf 3\nC2-S6 wholeGateMergeFixture rebuilds the mergeSet literal. auto-fix. conf 3\nWorst: C2-S1. Implementation command contributed to no finding.\n"
          },
          "axis": "Standards",
          "base": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
          "tip": "b91c696ce74e984b6cdd93d6d8758b35b61be62c",
          "finding_ids": [
            "C2-S1",
            "C2-S2",
            "C2-S3",
            "C2-S4",
            "C2-S5",
            "C2-S6"
          ],
          "supersedes": []
        },
        {
          "id": "r-c2-spec",
          "performer": "claude:ft393_c2_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c2_spec",
            "digest": "sha256:2bb008237491cc60bfc5164fab9aee471ca7c113f7e238677187cd9ec537f5cb",
            "excerpt": "RR-C2 Spec (claude:ft393_c2_spec): evidence current=true. Every mapped row is implemented and tested. 5 advisory findings.\nC2-P1 the selective-lane target grade is sound but untested. no-op. conf 6\nC2-P2 a non-line-safe missing-tree path prints the whole route as one placeholder. ask-user. conf 5\nC2-P3 create --from keeps the exec-form sibling route. ask-user. conf 4\nC2-P4 cleanLandedSiblings now releases an absent landed sibling, untested. ask-user. conf 4\nC2-P5 the landed route choice and the clean selector conditions may differ. no-op unless reproduced. conf 3\nWorst: C2-P2. Implementation command contributed to no finding.\n"
          },
          "axis": "Spec",
          "base": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
          "tip": "b91c696ce74e984b6cdd93d6d8758b35b61be62c",
          "finding_ids": [
            "C2-P1",
            "C2-P2",
            "C2-P3",
            "C2-P4",
            "C2-P5"
          ],
          "supersedes": []
        },
        {
          "id": "r-c2-coverage",
          "performer": "claude:ft393_c2_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "4eeb67dac8ffa9bc3b73ed77fcbd3bd28caec77e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c2_coverage",
            "digest": "sha256:0cf28b903a1ebd1e4110d68e76c919c7a1f442d17c4b3efb3e5d53a9336e7f9d",
            "excerpt": "RR-C2 Coverage (claude:ft393_c2_coverage): evidence current=true. 7 findings, 1 blocking.\nC2-1 blocking a missing tree whose registration was pruned still fails in releaseRegistration. ask-user. conf 6\nC2-2 the target-alone grade base has no selective-lane test. auto-fix. conf 6\nC2-3 the unlanded missing-tree cause is not followed from its printed route. auto-fix. conf 6\nC2-4 the merge and reset placeholder arms are untested. auto-fix. conf 6\nC2-5 two infrastructure arms of the target-alone grade are unproduced. auto-fix. conf 5\nC2-6 the LeaseUnknown arm of missingTreeLeaseRefusal is untested. auto-fix. conf 5\nC2-7 an unfaced merge refusal hands back to the reviewer. ask-user. conf 4\nWorst: C2-1. Implementation command contributed to no finding.\n"
          },
          "axis": "Coverage",
          "base": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
          "tip": "b91c696ce74e984b6cdd93d6d8758b35b61be62c",
          "finding_ids": [
            "C2-1",
            "C2-2",
            "C2-3",
            "C2-4",
            "C2-5",
            "C2-6",
            "C2-7"
          ],
          "supersedes": []
        },
        {
          "id": "r-c2-r1-standards",
          "performer": "claude:ft393_c2_r1_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c2_r1_standards",
            "digest": "sha256:e1052ce4ea27710dac8697c812e26371adf7441caea5fbcd6f9fc72920e7a78c",
            "excerpt": "RR-C2 confirming round, Standards (claude:ft393_c2_r1_standards): evidence current=true. All 15 folds confirmed. Zero findings.\nConcern 1: not flaky; bench sets maintenance.auto off (internal/env/kit_run.go:69); a missing loose object fails loudly.\nConcern 2: list.go:113-116 is a second derivation for a present tree with a detached or foreign HEAD; unchanged since the chunk base; advice only.\nConcern 3: obeys the Edge inventory intent; AtCheckout(line, \"\") omits the path step by design.\nConcern 4: producingFixtures keys on (face, cause); RR29 and ticket 07 line 29 still say \"exactly one\" per face, a wording gap.\nImplementation command contributed to no finding.\n"
          },
          "axis": "Standards",
          "base": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
          "tip": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "finding_ids": [],
          "supersedes": [
            "r-c2-standards"
          ]
        },
        {
          "id": "r-c2-r1-spec",
          "performer": "claude:ft393_c2_r1_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c2_r1_spec",
            "digest": "sha256:a35ce50355c3c9b58be3e64abf76d7e1459a5e4288ff3e1af809d5232380631a",
            "excerpt": "RR-C2 confirming round, Spec (claude:ft393_c2_r1_spec): evidence current=true. All folds confirmed. Zero findings.\nAll 16 changed files sit in a ticket 05-07 Writes line; the pin 761 to 766 matches the five new tests.\nConcern 2: same predicate through the planner; lease read differs only in an edge case; unchanged since the chunk base; no finding.\nConcern 3: obeys the Edge inventory intent; the path <id> step refuses for a missing tree.\nConcern 4: covered; producingFixtures enforces one fixture for each face and cause.\nConcern 1: not flaky; a packed object fails loudly.\nAdvice: RR29 and ticket 07 line 29 should say \"for each cause\"; the Edge inventory should record the missing-tree exception.\nImplementation command contributed to no finding.\n"
          },
          "axis": "Spec",
          "base": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
          "tip": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "finding_ids": [],
          "supersedes": [
            "r-c2-spec"
          ]
        },
        {
          "id": "r-c2-r1-coverage",
          "performer": "claude:ft393_c2_r1_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c2_r1_coverage",
            "digest": "sha256:58a452442f4c4211c4cdbd5e0fb46f18a9733e34ddfdd59127d126623a49010f",
            "excerpt": "RR-C2 confirming round, Coverage (claude:ft393_c2_r1_coverage): evidence current=true. All folds confirmed.\nC2-RC1 advisory reviews/refusal-route-registry.md:2273-2283: for C2-2, C2-3, C2-4, C2-5, C2-6 and C2-P4 the repro column records only the pre-fix silent probe and no post-fix red verdict. auto-fix (evidence only). conf 5\nConcern 1: not flaky; the tip commit is loose by construction; a packed object fails loudly with the wrong face.\nConcern 2: list.go:113-124 is outside the delta; a second derivation for present trees; advice.\nConcern 3: obeys the intent; AtCheckout with an empty id is the registry's own form.\nConcern 4: covered by the face and cause key.\nImplementation command contributed to no finding.\n"
          },
          "axis": "Coverage",
          "base": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
          "tip": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "finding_ids": [
            "C2-RC1"
          ],
          "supersedes": [
            "r-c2-coverage"
          ]
        },
        {
          "id": "r-c2-r2-coverage",
          "performer": "claude:ft393_c2_r1_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "862129435b89ab30c479c604b3f28f8ea272ecc6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c2_r1_coverage",
            "digest": "sha256:11603a43a25bb261ca7a99ac8c091ae57547d838d10f60b0c65268b4ec169df9",
            "excerpt": "RR-C2 confirming round, Coverage reaffirmation (claude:ft393_c2_r1_coverage) at record commit 6e16b78e.\nC2-RC1 resolved: the section \"RR-C2 repair cycle 1 red verdicts after the fix\" (about lines 2358-2373) gives a post-fix red verdict for C2-2, C2-3, C2-4, C2-5, C2-6, C2-P4 and C2-P5. Zero new findings.\nThe verdicts are recorded claims from the repair returns; the axis did not rerun them. The commit changes only the review record.\n"
          },
          "axis": "Coverage",
          "base": "33ac2f98902ac4a51bfa951218f50f3fd268319d",
          "tip": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "finding_ids": [],
          "supersedes": [
            "r-c2-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "RR-C3",
      "base": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
      "tip": "0f2825e1678e14f2dd2ad40222cd46ac7772c810",
      "plan_digest": "sha256:129cae8ef3792d9ae5a007ed37c6d3df9bb3b1712d5c4b452690b42fe42b87f6",
      "source_digest": "a512c4e4ab3ccea3aa89a010f0c9dc1c8db8cfb7",
      "acceptance_rows": [
        "RR32",
        "RR33",
        "RR34",
        "RR35",
        "RR36",
        "RR37",
        "RR38"
      ],
      "verification": [
        {
          "id": "v-t8-commit-package",
          "performer": "claude:ft393_t8",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "07c190e3138ef94b2670f6584dadfbc5aa70464e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t8",
            "digest": "sha256:2fd94881ccd1f156a0524247a608ecf0b0d68bed53987ee04155722da54b74ca",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit\ntree[1]{target,head,dirty}:\n  ft393-build,5f14753f5c9f20414b1da8c4e0a7e45074dd4093,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,8472,87\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t8-commit-package",
          "command": "bench test --package ./internal/commit",
          "exit_code": 0
        },
        {
          "id": "v-t8-commit-exit-three",
          "performer": "claude:ft393_t8",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "07c190e3138ef94b2670f6584dadfbc5aa70464e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t8",
            "digest": "sha256:714c339f55c93eaa0aaeb1bdbbf8749a870a353f09c7dd4f1c327b2df54c00fb",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'\ntree[1]{target,head,dirty}:\n  ft393-build,5f14753f5c9f20414b1da8c4e0a7e45074dd4093,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,429,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t8-commit-exit-three",
          "command": "bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'",
          "exit_code": 0
        },
        {
          "id": "v-t8-reset-route-proof",
          "performer": "claude:ft393_t8",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "07c190e3138ef94b2670f6584dadfbc5aa70464e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t8",
            "digest": "sha256:b8c02636e0859b006122ee54d589e4fe41771fe366905dceedc8e53bdf598b30",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'\ntree[1]{target,head,dirty}:\n  ft393-build,5f14753f5c9f20414b1da8c4e0a7e45074dd4093,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,352,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_commit.go --swap 'Command(Text(\"bench worktree reset --to\"), Fact(FactPublishedCommit), Fact(FactCheckout))' --with 'Command(Text(\"git restore --staged --worktree --source\"), Fact(FactPublishedCommit), Text(\"--\"), Fact(FactCheckout))' --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,2,yes\n"
          },
          "requirement": "t8-reset-route-proof",
          "command": "bench test --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'",
          "exit_code": 0,
          "probe": {
            "mutation": "Give the `commit-published-unreconciled` face the old `git restore` route. TestPublishedUnreconciledRouteIsTheResetPlan must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t8",
              "digest": "sha256:b8c02636e0859b006122ee54d589e4fe41771fe366905dceedc8e53bdf598b30",
              "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'\ntree[1]{target,head,dirty}:\n  ft393-build,5f14753f5c9f20414b1da8c4e0a7e45074dd4093,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,352,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_commit.go --swap 'Command(Text(\"bench worktree reset --to\"), Fact(FactPublishedCommit), Fact(FactCheckout))' --with 'Command(Text(\"git restore --staged --worktree --source\"), Fact(FactPublishedCommit), Text(\"--\"), Fact(FactCheckout))' --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,2,yes\n"
            }
          }
        },
        {
          "id": "v-t9-commit-package",
          "performer": "claude:ft393_t9b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "07c190e3138ef94b2670f6584dadfbc5aa70464e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t9b",
            "digest": "sha256:21b47171b879f7252aea8d579aeff73c9335e2b8376721e4c001e51df639a09d",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit\ntree[1]{target,head,dirty}:\n  ft393-build,5f14753f5c9f20414b1da8c4e0a7e45074dd4093,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,8360,87\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t9-commit-package",
          "command": "bench test --package ./internal/commit",
          "exit_code": 0
        },
        {
          "id": "v-t9-commit-exit-three",
          "performer": "claude:ft393_t9b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "07c190e3138ef94b2670f6584dadfbc5aa70464e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t9b",
            "digest": "sha256:175712a6d03166bdcc2b1be82a7f84a27d2319277c997625816db4050df00e56",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'\ntree[1]{target,head,dirty}:\n  ft393-build,5f14753f5c9f20414b1da8c4e0a7e45074dd4093,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,424,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t9-commit-exit-three",
          "command": "bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'",
          "exit_code": 0
        },
        {
          "id": "v-t9-red-route-proof",
          "performer": "claude:ft393_t9b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "07c190e3138ef94b2670f6584dadfbc5aa70464e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t9b",
            "digest": "sha256:def7c2322525a2d707da6c17600862b603277bb9187181f32529ba4384e79bb4",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,5f14753f5c9f20414b1da8c4e0a7e45074dd4093,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,1425,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nProbe: Drop the re-run step from the `commit-red` route\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_commit.go --swap 'run reports\")), commitRerun}' --with 'run reports\"))}' --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commit,TestCommitFacesFollowTheirRoutes/commit-red,\"refusal_route_test.go:131: face commit-red = (1, \\\"\\\", \\\"error: prospective authorization refused: inherited (the gate ran red on the composed tree and no green baseline attributes the red to this diff)\\\\\\\\nnext=repair each failure that the run r… (379 bytes)\",1\n"
          },
          "requirement": "t9-red-route-proof",
          "command": "bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the re-run step from the `commit-red` route. TestCommitFacesFollowTheirRoutes must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t9b",
              "digest": "sha256:def7c2322525a2d707da6c17600862b603277bb9187181f32529ba4384e79bb4",
              "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,5f14753f5c9f20414b1da8c4e0a7e45074dd4093,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,1425,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nProbe: Drop the re-run step from the `commit-red` route\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_commit.go --swap 'run reports\")), commitRerun}' --with 'run reports\"))}' --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commit,TestCommitFacesFollowTheirRoutes/commit-red,\"refusal_route_test.go:131: face commit-red = (1, \\\"\\\", \\\"error: prospective authorization refused: inherited (the gate ran red on the composed tree and no green baseline attributes the red to this diff)\\\\\\\\nnext=repair each failure that the run r… (379 bytes)\",1\n"
            }
          }
        },
        {
          "id": "v-t8-commit-package-r1",
          "performer": "claude:ft393_t8_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b3f091957045d14435a1efab24b54f1f2b92f97",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t8_repair1",
            "digest": "sha256:d7e6e974d915354d722d64de38bfb7dd875ff43f1166c966d714573d6f0ec9d2",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit\ntree[1]{target,head,dirty}:\n  ft393-build,0d0e0c1f2d9319b69ccd954eba83374fa7747327,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,9936,103\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t8-commit-package",
          "command": "bench test --package ./internal/commit",
          "exit_code": 0
        },
        {
          "id": "v-t8-commit-exit-three-r1",
          "performer": "claude:ft393_t8_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b3f091957045d14435a1efab24b54f1f2b92f97",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t8_repair1",
            "digest": "sha256:d6b33a5497f64636c3e59240f3bc856db07af22bb2e219f2df3eb13ae69443f2",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'\ntree[1]{target,head,dirty}:\n  ft393-build,0d0e0c1f2d9319b69ccd954eba83374fa7747327,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,666,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t8-commit-exit-three",
          "command": "bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'",
          "exit_code": 0
        },
        {
          "id": "v-t8-reset-route-proof-r1",
          "performer": "claude:ft393_t8_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b3f091957045d14435a1efab24b54f1f2b92f97",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t8_repair1",
            "digest": "sha256:f43a32061d6dbce1e374ce5bec54b0c65b999a2c13c54340656f3d0607ca59e6",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'\ntree[1]{target,head,dirty}:\n  ft393-build,0d0e0c1f2d9319b69ccd954eba83374fa7747327,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,356,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_commit.go --swap 'Route:     []Step{Command(Text(\"bench worktree reset --to\"), Fact(FactPublishedCommit), Fact(FactCheckout))},' --with 'Route:     []Step{Command(Text(\"git restore --source\"), Fact(FactPublishedCommit), Text(\"--staged --worktree -- .\"))},' --package ./internal/commit --run TestPublishedUnreconciledRouteIsTheResetPlan\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,2,yes\n"
          },
          "requirement": "t8-reset-route-proof",
          "command": "bench test --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'",
          "exit_code": 0,
          "probe": {
            "mutation": "Give the `commit-published-unreconciled` face the old `git restore` route. TestPublishedUnreconciledRouteIsTheResetPlan must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t8_repair1",
              "digest": "sha256:f43a32061d6dbce1e374ce5bec54b0c65b999a2c13c54340656f3d0607ca59e6",
              "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'\ntree[1]{target,head,dirty}:\n  ft393-build,0d0e0c1f2d9319b69ccd954eba83374fa7747327,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,356,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_commit.go --swap 'Route:     []Step{Command(Text(\"bench worktree reset --to\"), Fact(FactPublishedCommit), Fact(FactCheckout))},' --with 'Route:     []Step{Command(Text(\"git restore --source\"), Fact(FactPublishedCommit), Text(\"--staged --worktree -- .\"))},' --package ./internal/commit --run TestPublishedUnreconciledRouteIsTheResetPlan\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,2,yes\n"
            }
          }
        },
        {
          "id": "v-t9-commit-package-r1",
          "performer": "claude:ft393_t9_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b3f091957045d14435a1efab24b54f1f2b92f97",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t9_repair1",
            "digest": "sha256:7b780475ec595bad42f4494751a7fb5493d7fe54c09fa75cd57e810958735f93",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit\ntree[1]{target,head,dirty}:\n  ft393-build,0d0e0c1f2d9319b69ccd954eba83374fa7747327,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,10113,103\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t9-commit-package",
          "command": "bench test --package ./internal/commit",
          "exit_code": 0
        },
        {
          "id": "v-t9-commit-exit-three-r1",
          "performer": "claude:ft393_t9_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b3f091957045d14435a1efab24b54f1f2b92f97",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t9_repair1",
            "digest": "sha256:527a034f370be9d832fb3b0a4c9287d817679fd017b34a87e051f7057d250578",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'\ntree[1]{target,head,dirty}:\n  ft393-build,0d0e0c1f2d9319b69ccd954eba83374fa7747327,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,429,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t9-commit-exit-three",
          "command": "bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'",
          "exit_code": 0
        },
        {
          "id": "v-t9-red-route-proof-r1",
          "performer": "claude:ft393_t9_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9b3f091957045d14435a1efab24b54f1f2b92f97",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t9_repair1",
            "digest": "sha256:53fcb5d4df26984d046d40faa1786c4a1fdf91aaa9bd7da8aeffa5431dcb0304",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,0d0e0c1f2d9319b69ccd954eba83374fa7747327,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,2128,9\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_commit.go --swap 'Route:     []Step{Instruction(Text(\"repair each failure that the run reports\")), commitRerun},' --with 'Route:     []Step{Instruction(Text(\"repair each failure that the run reports\"))},' --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,3,yes\n"
          },
          "requirement": "t9-red-route-proof",
          "command": "bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the re-run step from the `commit-red` route. TestCommitFacesFollowTheirRoutes must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t9_repair1",
              "digest": "sha256:53fcb5d4df26984d046d40faa1786c4a1fdf91aaa9bd7da8aeffa5431dcb0304",
              "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,0d0e0c1f2d9319b69ccd954eba83374fa7747327,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,2128,9\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_commit.go --swap 'Route:     []Step{Instruction(Text(\"repair each failure that the run reports\")), commitRerun},' --with 'Route:     []Step{Instruction(Text(\"repair each failure that the run reports\"))},' --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,3,yes\n"
            }
          }
        },
        {
          "id": "v-t8-commit-package-r2",
          "performer": "claude:ft393_t8_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a512c4e4ab3ccea3aa89a010f0c9dc1c8db8cfb7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t8_repair1",
            "digest": "sha256:c2a0f2ede1fa0c80488c7d4e8231aa920202f14ccf74d27cf5d0e869f2c5984e",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit\ntree[1]{target,head,dirty}:\n  ft393-build,0f2825e1678e14f2dd2ad40222cd46ac7772c810,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,10352,104\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t8-commit-package",
          "command": "bench test --package ./internal/commit",
          "exit_code": 0
        },
        {
          "id": "v-t8-commit-exit-three-r2",
          "performer": "claude:ft393_t8_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a512c4e4ab3ccea3aa89a010f0c9dc1c8db8cfb7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t8_repair1",
            "digest": "sha256:c2cc59fb8370f97ed14b51c6eb585d5a162c55a7fb509c6d8289c735220eec46",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'\ntree[1]{target,head,dirty}:\n  ft393-build,0f2825e1678e14f2dd2ad40222cd46ac7772c810,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,435,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t8-commit-exit-three",
          "command": "bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'",
          "exit_code": 0
        },
        {
          "id": "v-t8-reset-route-proof-r2",
          "performer": "claude:ft393_t8_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a512c4e4ab3ccea3aa89a010f0c9dc1c8db8cfb7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t8_repair1",
            "digest": "sha256:5a300fc541084b96e427210e8393306f1c55f80f6162e58c48935a7304a441a0",
            "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'\ntree[1]{target,head,dirty}:\n  ft393-build,0f2825e1678e14f2dd2ad40222cd46ac7772c810,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,348,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_commit.go --swap 'Route:     []Step{Command(Text(\"bench worktree reset --to\"), Fact(FactPublishedCommit), Fact(FactCheckout))},' --with 'Route:     []Step{Command(Text(\"git restore --source\"), Fact(FactPublishedCommit), Text(\"--staged --worktree -- .\"))},' --package ./internal/commit --run TestPublishedUnreconciledRouteIsTheResetPlan\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,2,yes\n"
          },
          "requirement": "t8-reset-route-proof",
          "command": "bench test --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'",
          "exit_code": 0,
          "probe": {
            "mutation": "Give the `commit-published-unreconciled` face the old `git restore` route. TestPublishedUnreconciledRouteIsTheResetPlan must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t8_repair1",
              "digest": "sha256:5a300fc541084b96e427210e8393306f1c55f80f6162e58c48935a7304a441a0",
              "excerpt": "$ bench worktree exec ft393-build -- bench test --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'\ntree[1]{target,head,dirty}:\n  ft393-build,0f2825e1678e14f2dd2ad40222cd46ac7772c810,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,348,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n$ bench worktree exec ft393-build -- bench probe internal/refusalroute/faces_commit.go --swap 'Route:     []Step{Command(Text(\"bench worktree reset --to\"), Fact(FactPublishedCommit), Fact(FactCheckout))},' --with 'Route:     []Step{Command(Text(\"git restore --source\"), Fact(FactPublishedCommit), Text(\"--staged --worktree -- .\"))},' --package ./internal/commit --run TestPublishedUnreconciledRouteIsTheResetPlan\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,2,yes\n"
            }
          }
        },
        {
          "id": "v-t9-commit-package-r2",
          "performer": "claude:ft393_t9_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a512c4e4ab3ccea3aa89a010f0c9dc1c8db8cfb7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t9_repair2",
            "digest": "sha256:e19fc69c5a9ea36f34b6038c9c742c26c70a87316a217b2d78afde64924b65f9",
            "excerpt": "$ bench test --package ./internal/commit\ntree[1]{target,head,dirty}:\n  ft393-build,0f2825e1678e14f2dd2ad40222cd46ac7772c810,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,10096,104\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t9-commit-package",
          "command": "bench test --package ./internal/commit",
          "exit_code": 0
        },
        {
          "id": "v-t9-commit-exit-three-r2",
          "performer": "claude:ft393_t9_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a512c4e4ab3ccea3aa89a010f0c9dc1c8db8cfb7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t9_repair2",
            "digest": "sha256:c09616f7bf38be9059be9f8b1c52a8c00a1eedc0258992634e8f6e5b55d1121d",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'\ntree[1]{target,head,dirty}:\n  ft393-build,0f2825e1678e14f2dd2ad40222cd46ac7772c810,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,423,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t9-commit-exit-three",
          "command": "bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'",
          "exit_code": 0
        },
        {
          "id": "v-t9-red-route-proof-r2",
          "performer": "claude:ft393_t9_repair2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a512c4e4ab3ccea3aa89a010f0c9dc1c8db8cfb7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t9_repair2",
            "digest": "sha256:830511dbf28698ae7f9cb62a02f43cb52d453492802fba0cb07015ea36a52dd3",
            "excerpt": "$ bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,0f2825e1678e14f2dd2ad40222cd46ac7772c810,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,2452,10\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/faces_commit.go --swap 'Instruction(Text(\"repair each failure that the run reports\")), commitRerun}' --with 'Instruction(Text(\"repair each failure that the run reports\"))}' --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,3,yes\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commit,TestCommitFacesFollowTheirRoutes/commit-red/,...\n  github.com/gibbonmi/bench/internal/commit,TestCommitFacesFollowTheirRoutes/commit-red/dry_run,...\n  github.com/gibbonmi/bench/internal/commit,TestCommitFacesFollowTheirRoutes/commit-red/lane,...\n"
          },
          "requirement": "t9-red-route-proof",
          "command": "bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the re-run step from the `commit-red` route. TestCommitFacesFollowTheirRoutes must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t9_repair2",
              "digest": "sha256:830511dbf28698ae7f9cb62a02f43cb52d453492802fba0cb07015ea36a52dd3",
              "excerpt": "$ bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\ntree[1]{target,head,dirty}:\n  ft393-build,0f2825e1678e14f2dd2ad40222cd46ac7772c810,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,2452,10\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/refusalroute/faces_commit.go --swap 'Instruction(Text(\"repair each failure that the run reports\")), commitRerun}' --with 'Instruction(Text(\"repair each failure that the run reports\"))}' --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/refusalroute/faces_commit.go,swap,failed,3,yes\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commit,TestCommitFacesFollowTheirRoutes/commit-red/,...\n  github.com/gibbonmi/bench/internal/commit,TestCommitFacesFollowTheirRoutes/commit-red/dry_run,...\n  github.com/gibbonmi/bench/internal/commit,TestCommitFacesFollowTheirRoutes/commit-red/lane,...\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "r-c3-standards",
          "performer": "claude:ft393_c3_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "07c190e3138ef94b2670f6584dadfbc5aa70464e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c3_standards",
            "digest": "sha256:0d0fb09c53fecd11c3e0aafec8737441ec98032b4a64b50549088e82f2ec38ef",
            "excerpt": "RR-C3 Standards (claude:ft393_c3_standards): evidence current=true. 4 findings, 1 blocking.\nC3-S1 blocking internal/commit/refusal_route_test.go:99-221 repeats the worktree route-walk harness (marker split, step count, carry-or-run, simple-command check, doctor exemption); tickets 10-12 would add more copies. ask-user (new shared seam). conf 7\nC3-S2 advisory internal/worktree/land_rerun.go:43-50 landingRerunArg repeats the quote-or-placeholder rule that refusalroute.Arg exports. auto-fix with a plan commit. conf 7\nC3-S3 advisory usage.WorktreeCreate and faces_commit.go:26 derive the create command twice; the primary-checkout output names it twice. ask-user. conf 5\nC3-S4 advisory spec line 178 says next= is on the stderr refusal line; the code and RR37 print its own line. auto-fix the spec text. conf 6\nWorst: C3-S1. Implementation command contributed to C3-S1: the plan names no shared route-walk seam.\n"
          },
          "axis": "Standards",
          "base": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "tip": "5f14753f5c9f20414b1da8c4e0a7e45074dd4093",
          "finding_ids": [
            "C3-S1",
            "C3-S2",
            "C3-S3",
            "C3-S4"
          ],
          "supersedes": []
        },
        {
          "id": "r-c3-spec",
          "performer": "claude:ft393_c3_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "07c190e3138ef94b2670f6584dadfbc5aa70464e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c3_spec",
            "digest": "sha256:01cc07ad2af8eedbb580a5f6c2fd636b7bceec04b7abf97bf6016dbea827151f",
            "excerpt": "RR-C3 Spec (claude:ft393_c3_spec): evidence current=true. Rows RR32 to RR38 each map to a test in the delta. 5 findings, 1 blocking.\nC3-P1 blocking internal/commit/commit.go:116-200 landingFace sends every non-authorization landing error to the reviewer handback, including causes an agent clears alone (a missing named path, untracked content, a Go format error). ask-user. conf 6\nC3-P2 advisory commit.go:122-123 and usage/worktree.go:65 the primary-checkout output names the create route twice, from two sources. ask-user. conf 6\nC3-P3 advisory internal/commit/refusal_route_test.go:471-571 re-implements the worktree route-walk helpers. auto-fix. conf 7\nC3-P4 advisory internal/worktree/land_rerun.go:45-50 landingRerunArg copies the rule that refusalroute.Arg exports. auto-fix. conf 6\nC3-P5 advisory commit.go:235-247 the rerun has no hostile-argument test (dry-run, preflight-build, multi-line message, control byte, no-assignment label). auto-fix. conf 7\nItem (b) next= on its own stderr line: no-op. Worst: C3-P1. Implementation command contributed to C3-P3.\n"
          },
          "axis": "Spec",
          "base": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "tip": "5f14753f5c9f20414b1da8c4e0a7e45074dd4093",
          "finding_ids": [
            "C3-P1",
            "C3-P2",
            "C3-P3",
            "C3-P4",
            "C3-P5"
          ],
          "supersedes": []
        },
        {
          "id": "r-c3-coverage",
          "performer": "claude:ft393_c3_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "07c190e3138ef94b2670f6584dadfbc5aa70464e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c3_coverage",
            "digest": "sha256:b499b0b5054bc0ae1e299f716e526908af12020d1f1e67d3377503387d31eac1",
            "excerpt": "RR-C3 Coverage (claude:ft393_c3_coverage): evidence current=true. 5 findings, all advisory.\nC3-1 advisory internal/commit/commit.go:238-247 the rerun arms (--dry-run, --preflight-build, msg and path placeholders) have no fixture; all fixtures run -m m -- a.txt. auto-fix. conf 7\nC3-2 advisory commit.go:116-163 and landing/attribution.go agent-clearable causes route to the reviewer handback; no decision records it. ask-user. conf 5\nC3-3 advisory internal/commit/refusal_route_test.go:100-106 the walk keys on Faces(Commit) only; ticket 11 widens it. no-op. conf 6\nC3-4 advisory commit.go:219-228 no lane-red or dry-run red fixture reaches commit-red. auto-fix. conf 6\nC3-5 advisory commit.go:236 no fixture covers the no-assignment <label> placeholder. auto-fix. conf 5\nWorst: C3-2. Implementation command: no finding traced to it.\n"
          },
          "axis": "Coverage",
          "base": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "tip": "5f14753f5c9f20414b1da8c4e0a7e45074dd4093",
          "finding_ids": [
            "C3-1",
            "C3-2",
            "C3-3",
            "C3-4",
            "C3-5"
          ],
          "supersedes": []
        },
        {
          "id": "r-c3-r1-standards",
          "performer": "claude:ft393_c3_r1_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9b3f091957045d14435a1efab24b54f1f2b92f97",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c3_r1_standards",
            "digest": "sha256:7cf8e67362461228260234cc9b9c9a625c9b69f2b2614df4f7006f844001cf38",
            "excerpt": "RR-C3 confirming round, Standards (claude:ft393_c3_r1_standards): evidence current=true. C3-S1, C3-S2, C3-S4 confirmed. 2 findings, 1 blocking.\nC3-RS1 advisory internal/refusalroute/faces_commit.go:53-54: the malformed lane declaration stays in commit-handback under a comment rule (\"a change to the gate's own checks is the reviewer's\") that the spec Authority section does not state. ask-user: amend the spec or move the cause. conf 6\nC3-RS2 blocking internal/landing/gitexec.go:19 with internal/commit/commit.go:216-226: the compare-and-swap refusal says \"rerun the landing\" but routes to commit-handback; a rerun-only agent face is the clean fit. ask-user. conf 6\nConcern 3: not a gap. Concern 4: routetest reusable as it stands.\nAdvice: reviewerRoute forwarder; the carry closure repeats in four drivers; Words rejects \"<\"; mixed NamedPathError literal forms; table test placement.\nImplementation command contributed to no finding.\n"
          },
          "axis": "Standards",
          "base": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "tip": "0d0e0c1f2d9319b69ccd954eba83374fa7747327",
          "finding_ids": [
            "C3-RS1",
            "C3-RS2"
          ],
          "supersedes": [
            "r-c3-standards"
          ]
        },
        {
          "id": "r-c3-r1-spec",
          "performer": "claude:ft393_c3_r1_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9b3f091957045d14435a1efab24b54f1f2b92f97",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c3_r1_spec",
            "digest": "sha256:6091acdb11e8683291e88664a785c6ace53305e7f2521aac44ded570f37644b8",
            "excerpt": "RR-C3 confirming round, Spec (claude:ft393_c3_r1_spec): evidence current=true. All folds confirmed for named-path causes. 2 findings, 1 blocking.\nC3-RP1 blocking internal/landing/gitexec.go:19 with internal/commit/commit.go:228: the compare-and-swap refusal falls to the reviewer handback; the Authority rule makes it an agent cause; it has no typed signal. ask-user (face and type). conf 6\nC3-RP2 advisory spec.md:1060 and ticket 09 line 32: RR36 says exactly one producing fixture per face; commit-red now has three. auto-fix plan wording. conf 7\nAdvice: name gate-check edits in the Authority list; commit-named-path is absent from the face table; concern 3 not a gap; concern 4 routetest reusable.\nImplementation command contributed to no finding.\n"
          },
          "axis": "Spec",
          "base": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "tip": "0d0e0c1f2d9319b69ccd954eba83374fa7747327",
          "finding_ids": [
            "C3-RP1",
            "C3-RP2"
          ],
          "supersedes": [
            "r-c3-spec"
          ]
        },
        {
          "id": "r-c3-r1-coverage",
          "performer": "claude:ft393_c3_r1_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "9b3f091957045d14435a1efab24b54f1f2b92f97",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c3_r1_coverage",
            "digest": "sha256:29e10f8e695ad99d134109f47bc27d2ada5d01cb68de0ca34617711e6fb203dd",
            "excerpt": "RR-C3 confirming round, Coverage (claude:ft393_c3_r1_coverage): evidence current=true. All folds confirmed for named paths. 1 finding, blocking.\nC3-RC1 blocking internal/landing/gitexec.go:19 and internal/commit/commit.go:219-228: the compare-and-swap refusal tells the caller to rerun, but its plain error routes to commit-handback; the spec Authority rule makes a cause that a Bench verb clears an agent cause (precedent checkpoint-tip-moved). No fixture covers it. ask-user (face choice). conf 6\nConcern 1: no finding; the gate invariant supports the malformed lane in handback; the spec should name the cause. advisory.\nConcern 4: routetest is reusable as it stands.\nAdvice: the wrap \"at least one path is required\" is unreachable through the commit and no mutation turns it red; spec.md:218-221 does not list commit-named-path.\nImplementation command contributed to no finding.\n"
          },
          "axis": "Coverage",
          "base": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "tip": "0d0e0c1f2d9319b69ccd954eba83374fa7747327",
          "finding_ids": [
            "C3-RC1"
          ],
          "supersedes": [
            "r-c3-coverage"
          ]
        },
        {
          "id": "r-c3-r2-standards",
          "performer": "claude:ft393_c3_r2_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a512c4e4ab3ccea3aa89a010f0c9dc1c8db8cfb7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c3_r2_standards",
            "digest": "sha256:9266e5887f6d13e77adb7a99cd11b8bd4cbeed15dad3520e69ed7a23cb447559",
            "excerpt": "RR-C3 cycle 2 confirming round, Standards (claude:ft393_c3_r2_standards): evidence current=true. C3-RC1/C3-RS2/C3-RP1, C3-RS1, and C3-RP2 confirmed. Zero findings.\nConcern 1: the fixture is deterministic; the admin-dir marker stops a second move, and it lives under t.TempDir.\nConcern 2: land and merge face selection unchanged. Concern 3: the three commit.go edits change no behavior.\nAdvice: move tipMovedFixture beside the other fixtures; two typed errors do not need a shared abstraction yet.\nImplementation command contributed to no finding.\n"
          },
          "axis": "Standards",
          "base": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "tip": "0f2825e1678e14f2dd2ad40222cd46ac7772c810",
          "finding_ids": [],
          "supersedes": [
            "r-c3-r1-standards"
          ]
        },
        {
          "id": "r-c3-r2-spec",
          "performer": "claude:ft393_c3_r2_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a512c4e4ab3ccea3aa89a010f0c9dc1c8db8cfb7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c3_r2_spec",
            "digest": "sha256:653003881bc597efe4a9a2d6596b71ca388250fc0e1c90315cc68eaef3c06068",
            "excerpt": "RR-C3 cycle 2 confirming round, Spec (claude:ft393_c3_r2_spec): evidence current=true. C3-RC1/C3-RS2/C3-RP1, C3-RS1, and C3-RP2 confirmed. Zero findings.\nConcern 2: no change; land and merge select faces by refusalError, AuthorizationRefusal, or the sentence, and the Error() text and Unwrap target are unchanged.\nAdvice: TipMovedError and NamedPathError could share one file.\nImplementation command contributed to no finding.\n"
          },
          "axis": "Spec",
          "base": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "tip": "0f2825e1678e14f2dd2ad40222cd46ac7772c810",
          "finding_ids": [],
          "supersedes": [
            "r-c3-r1-spec"
          ]
        },
        {
          "id": "r-c3-r2-coverage",
          "performer": "claude:ft393_c3_r2_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a512c4e4ab3ccea3aa89a010f0c9dc1c8db8cfb7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_c3_r2_coverage",
            "digest": "sha256:fc47c7e4bb6167f5eeb37679489d214a1d13df6dd195265459e1ee606940472a",
            "excerpt": "RR-C3 cycle 2 confirming round, Coverage (claude:ft393_c3_r2_coverage): evidence current=true. C3-RC1/C3-RS2/C3-RP1, C3-RS1, and C3-RP2 confirmed. Zero findings.\nThe review record gives two red states for commit-tip-moved: the fixture printed the handback before the fix, and a probe that removed the typed wrap returned bit.\nConcern 1: the fixture is deterministic; a gate that ran too early turns the walk red, not falsely green.\nConcern 2: no face selection changed. Concern 3: no behavior change.\nImplementation command contributed to no finding.\n"
          },
          "axis": "Coverage",
          "base": "fba4fe785c6c55f8ad7a2cc1522efa33730332cb",
          "tip": "0f2825e1678e14f2dd2ad40222cd46ac7772c810",
          "finding_ids": [],
          "supersedes": [
            "r-c3-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "RR-C4",
      "base": "0f2825e1678e14f2dd2ad40222cd46ac7772c810",
      "tip": "819154c6b9007e7c9019e19287f228a51462e02c",
      "plan_digest": "sha256:09259ca04468131aff9125326c043078881ec4190cafa38628483388cb98a239",
      "source_digest": "3bdd72ca8defafd138a1edd564cddbdffca873b8",
      "acceptance_rows": [
        "RR39",
        "RR40",
        "RR41",
        "RR42",
        "RR43",
        "RR44",
        "RR66",
        "RR68"
      ],
      "verification": [
        {
          "id": "v-t10-gate-package",
          "performer": "claude:ft393_t10b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8b17934c6d79e6c796068437f724283b0a90bee6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t10b",
            "digest": "sha256:efcb5b09cb5a08f36008b122e778682e541eb46032afb6915abe02854ad88d65",
            "excerpt": "$ bench test --package ./internal/gate\ntree[1]{target,head,dirty}:\n  ft393-build,7879b23afb093bf754886394774b2193e4347968,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,26683,408\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t10-gate-package",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "v-t10-checkpoint-routes",
          "performer": "claude:ft393_t10b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8b17934c6d79e6c796068437f724283b0a90bee6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t10b",
            "digest": "sha256:8b9b4f040133e8abbbe3fc948b8376da55123e5a86d12139d314f64f1ce9d281",
            "excerpt": "$ bench test --package ./internal/gate --run 'TestCheckpointFacesFollowTheirRoutes|TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'\ntree[1]{target,head,dirty}:\n  ft393-build,7879b23afb093bf754886394774b2193e4347968,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,1664,12\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t10-checkpoint-routes",
          "command": "bench test --package ./internal/gate --run 'TestCheckpointFacesFollowTheirRoutes|TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'",
          "exit_code": 0
        },
        {
          "id": "v-t10-help-row-proof",
          "performer": "claude:ft393_t10b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8b17934c6d79e6c796068437f724283b0a90bee6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t10b",
            "digest": "sha256:6bc5397fbb60102a87b7587c7109331a39f832b95629d9a56be7b05d214cce1e",
            "excerpt": "$ bench test --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'\ntree[1]{target,head,dirty}:\n  ft393-build,7879b23afb093bf754886394774b2193e4347968,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,260,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nProbe: Print the fixed write-access `help[1]{cmd,why}` row on the checkpoint refusal again\n$ bench probe internal/gate/run_transaction.go --swap 'return refuse(ctx, storageRoot, stderr, mode, funnelFace(err),' --with 'fmt.Fprint(stdout, \"help[1]{cmd,why}:\\n  bench gate --fresh,retry after restoring repository write access\\n\"); return refuse(ctx, storageRoot, stderr, mode, funnelFace(err),' --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,3,yes\n"
          },
          "requirement": "t10-help-row-proof",
          "command": "bench test --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'",
          "exit_code": 0,
          "probe": {
            "mutation": "Print the fixed write-access `help[1]{cmd,why}` row on the checkpoint refusal again. TestReviewCheckpointRefusalRoute must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t10b",
              "digest": "sha256:6bc5397fbb60102a87b7587c7109331a39f832b95629d9a56be7b05d214cce1e",
              "excerpt": "$ bench test --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'\ntree[1]{target,head,dirty}:\n  ft393-build,7879b23afb093bf754886394774b2193e4347968,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,260,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\nProbe: Print the fixed write-access `help[1]{cmd,why}` row on the checkpoint refusal again\n$ bench probe internal/gate/run_transaction.go --swap 'return refuse(ctx, storageRoot, stderr, mode, funnelFace(err),' --with 'fmt.Fprint(stdout, \"help[1]{cmd,why}:\\n  bench gate --fresh,retry after restoring repository write access\\n\"); return refuse(ctx, storageRoot, stderr, mode, funnelFace(err),' --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,3,yes\n"
            }
          }
        },
        {
          "id": "v-t10-gate-package-r1",
          "performer": "claude:ft393_t10_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3bdd72ca8defafd138a1edd564cddbdffca873b8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t10_repair1",
            "digest": "sha256:8d6d67e5ae6903deb17bd9a7f146cd88b12e2090aae432ab675272c66bc3d022",
            "excerpt": "$ bench test --package ./internal/gate\ntree[1]{target,head,dirty}:\n  ft393-build,819154c6b9007e7c9019e19287f228a51462e02c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,26841,418\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t10-gate-package",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "v-t10-checkpoint-routes-r1",
          "performer": "claude:ft393_t10_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3bdd72ca8defafd138a1edd564cddbdffca873b8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t10_repair1",
            "digest": "sha256:678fb49e461521571813fd03f049c0fa1721d6e96e23ae85b2b74c84a2fc1f23",
            "excerpt": "$ bench test --package ./internal/gate --run 'TestCheckpointFacesFollowTheirRoutes|TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'\ntree[1]{target,head,dirty}:\n  ft393-build,819154c6b9007e7c9019e19287f228a51462e02c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,1594,12\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t10-checkpoint-routes",
          "command": "bench test --package ./internal/gate --run 'TestCheckpointFacesFollowTheirRoutes|TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'",
          "exit_code": 0
        },
        {
          "id": "v-t10-help-row-proof-r1",
          "performer": "claude:ft393_t10_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3bdd72ca8defafd138a1edd564cddbdffca873b8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft393_t10_repair1",
            "digest": "sha256:5b90b79db1e0fadbecb774c5c549daded7b2bd1e6ca633bc9aefc711b89c7375",
            "excerpt": "$ bench test --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'\ntree[1]{target,head,dirty}:\n  ft393-build,819154c6b9007e7c9019e19287f228a51462e02c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,259,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/gate/run_transaction.go --swap 'return refuse(ctx, storageRoot, stderr, mode, funnelFace(err),' --with 'fmt.Fprintln(stdout, \"help[1]{cmd,why}:\"); return refuse(ctx, storageRoot, stderr, mode, funnelFace(err),' --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,3,yes\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateRunRetainsSubjectConstructionCause,\"run_outcomes_test.go:201: stdout = \\\"help[1]{cmd,why}:\\\\n\\\", want no help row beside the route\",1\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointRefusalRoute/missing_plan,\"review_checkpoint_test.go:248: refusal printed a help row beside its route: help[1]{cmd,why}:\",3\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointRefusalRoute/missing_record,\"review_checkpoint_test.go:248: refusal printed a help row beside its route: help[1]{cmd,why}:\",3\n"
          },
          "requirement": "t10-help-row-proof",
          "command": "bench test --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'",
          "exit_code": 0,
          "probe": {
            "mutation": "Print the fixed write-access `help[1]{cmd,why}` row on the checkpoint refusal again. TestReviewCheckpointRefusalRoute must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft393_t10_repair1",
              "digest": "sha256:5b90b79db1e0fadbecb774c5c549daded7b2bd1e6ca633bc9aefc711b89c7375",
              "excerpt": "$ bench test --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'\ntree[1]{target,head,dirty}:\n  ft393-build,819154c6b9007e7c9019e19287f228a51462e02c,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,259,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/gate/run_transaction.go --swap 'return refuse(ctx, storageRoot, stderr, mode, funnelFace(err),' --with 'fmt.Fprintln(stdout, \"help[1]{cmd,why}:\"); return refuse(ctx, storageRoot, stderr, mode, funnelFace(err),' --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,3,yes\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateRunRetainsSubjectConstructionCause,\"run_outcomes_test.go:201: stdout = \\\"help[1]{cmd,why}:\\\\n\\\", want no help row beside the route\",1\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointRefusalRoute/missing_plan,\"review_checkpoint_test.go:248: refusal printed a help row beside its route: help[1]{cmd,why}:\",3\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointRefusalRoute/missing_record,\"review_checkpoint_test.go:248: refusal printed a help row beside its route: help[1]{cmd,why}:\",3\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "r-c4-standards",
          "performer": "claude:ft393_c4_standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "8b17934c6d79e6c796068437f724283b0a90bee6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c4_standards",
            "digest": "sha256:421220a21bcca72ee612f67991eebae4b73e6242254dd7ff8bfc104b63fd54d4",
            "excerpt": "RR-C4 Standards (claude:ft393_c4_standards): evidence current=true. 6 findings, all advisory.\nC4-S1 advisory internal/gate/checkpoint.go:149: the forced --fresh rerun rule is derived in gate code and in the face comment and spec row. ask-user. conf 6\nC4-S2 advisory internal/gate/checkpoint.go:103-105,203: completionProofError wraps capture faults that its doc excludes; they route to the reviewer handback. ask-user. conf 7\nC4-S3 advisory internal/gate/refusal_route_test.go:612-617, run_outcomes_test.go:195-196: three copies of the next= line reader beside routetest. auto-fix. conf 6\nC4-S4 advisory internal/gate/checkpoint.go:158-172: a rerun with no arguments renders <arguments>, a placeholder the agent cannot fill. ask-user. conf 4\nC4-S5 advisory specs/native-record-operations/spec.md:176 cites the deleted gate.routedRefusal. ask-user. conf 6\nC4-S6 advisory gate test comments carry RR row ids; the tree convention is mixed. ask-user. conf 5\nWorst: C4-S2. Implementation command contributed to C4-S1 and C4-S3.\n"
          },
          "axis": "Standards",
          "base": "0f2825e1678e14f2dd2ad40222cd46ac7772c810",
          "tip": "7879b23afb093bf754886394774b2193e4347968",
          "finding_ids": [
            "C4-S1",
            "C4-S2",
            "C4-S3",
            "C4-S4",
            "C4-S5",
            "C4-S6"
          ],
          "supersedes": []
        },
        {
          "id": "r-c4-spec",
          "performer": "claude:ft393_c4_spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "8b17934c6d79e6c796068437f724283b0a90bee6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c4_spec",
            "digest": "sha256:93a63a3dffa6833548a1e8e5711ff5f94bd535260cf2b3c4cfbd359ede8b7889",
            "excerpt": "RR-C4 Spec (claude:ft393_c4_spec): evidence current=true. Rows RR39-RR44, RR66, RR68 each map to code and an assertion. 4 findings, all advisory.\nC4-P1 advisory internal/gate/checkpoint.go:198-200 with completion.go:37-44: completionProofError routes capture and read faults to gate-handback; spec line 306 gives them checkpoint-subject-unavailable. ask-user (completion.go outside the fence). conf 7\nC4-P2 advisory internal/gate/engine.go:55: the \"prospective gate subject unavailable\" refusal prints no next=; spec line 306 and story 34 apply. ask-user (engine.go outside the fence). conf 6\nC4-P3 advisory internal/refusalroute/faces_gate.go:14: Composed(FactSlug) prints a hostile slug unquoted; the edge inventory requires quoted path facts. ask-user. conf 6\nC4-P4 advisory internal/gate/checkpoint.go:114-126: refuse prints a gate next= inside land and commit authorization runs; story 35 asks for one route per refusal; pre-existing. ask-user. conf 6\nAuthor items (a) and (b): no-op. Worst: C4-P1. Implementation command contributed to no finding.\n"
          },
          "axis": "Spec",
          "base": "0f2825e1678e14f2dd2ad40222cd46ac7772c810",
          "tip": "7879b23afb093bf754886394774b2193e4347968",
          "finding_ids": [
            "C4-P1",
            "C4-P2",
            "C4-P3",
            "C4-P4"
          ],
          "supersedes": []
        },
        {
          "id": "r-c4-coverage",
          "performer": "claude:ft393_c4_coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "8b17934c6d79e6c796068437f724283b0a90bee6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft393_c4_coverage",
            "digest": "sha256:fe519c28844995e2c22a9842fa1105df0f5cc4b8754cd82c265cfb295855f6db",
            "excerpt": "RR-C4 Coverage (claude:ft393_c4_coverage): evidence current=true. Rows RR39-RR44, RR66, RR68 each map to a test that asserts the row text. 6 findings, 1 blocking.\nC4-1 advisory internal/refusalroute/faces_gate.go: the slug renders unquoted through Composed(FactSlug); no test renders a hostile slug, spec path, or chunk id. auto-fix. conf 7\nC4-2 blocking internal/gate/checkpoint.go:142-153: a funnel refusal inside a land or commit authorization run prints a gate next= (for a tickets-only landing, bench gate --in <label> <arguments>) that the caller does not run, before the caller's own route. ask-user. conf 6\nC4-3 advisory internal/gate/checkpoint.go:201-203: completionProofError routes capture faults to the reviewer handback; no test injects one. ask-user. conf 6\nC4-4 advisory internal/gate/engine.go:53-56: prospective gate subject unavailable prints no route; spec line 1147 may cover it. ask-user. conf 5\nC4-5 advisory internal/gate/checkpoint.go:149: the mode == forceRun term has no test. auto-fix. conf 6\nC4-6 advisory internal/gate/refusal_route_test.go:544: the handback and tip-moved walks pass with no repair. no-op. conf 4\nWorst: C4-2. Implementation command contributed to C4-2, C4-3, C4-4: a ticket that adds a printed route to a shared funnel needs a consumer sweep of its callers.\n"
          },
          "axis": "Coverage",
          "base": "0f2825e1678e14f2dd2ad40222cd46ac7772c810",
          "tip": "7879b23afb093bf754886394774b2193e4347968",
          "finding_ids": [
            "C4-1",
            "C4-2",
            "C4-3",
            "C4-4",
            "C4-5",
            "C4-6"
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
      "from": "sha256:ce25426aa29127e407386b36ae00130d3db020f722da4e0e7639a45f75ae8a07",
      "to": "sha256:10d938cdde11624c65c40c632beebb04b986d194da031b1c8dd8d470d65aaa05",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ]
      }
    },
    {
      "from": "sha256:10d938cdde11624c65c40c632beebb04b986d194da031b1c8dd8d470d65aaa05",
      "to": "sha256:798ef1ecdef294b8be9eb9af70ca2280e069e137efdb28b02071046d4f0f5132",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ]
      }
    },
    {
      "from": "sha256:798ef1ecdef294b8be9eb9af70ca2280e069e137efdb28b02071046d4f0f5132",
      "to": "sha256:1958801be24bda9d85bc5db93b391183c898455fb52adc9154b785f6f4a8d984",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ],
        "RR-C1b": [
          "RR-C1b"
        ]
      }
    },
    {
      "from": "sha256:1958801be24bda9d85bc5db93b391183c898455fb52adc9154b785f6f4a8d984",
      "to": "sha256:f0639cf27495f50d19843b401bf84eab059aa614e05f76f81c98dc485790924b",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ],
        "RR-C1b": [
          "RR-C1b"
        ]
      }
    },
    {
      "from": "sha256:f0639cf27495f50d19843b401bf84eab059aa614e05f76f81c98dc485790924b",
      "to": "sha256:43985cc5492c46ba00fc0a2dfe29aaeda496478ea800c97bdfe4e22ba1b89720",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ],
        "RR-C1b": [
          "RR-C1b"
        ]
      }
    },
    {
      "from": "sha256:43985cc5492c46ba00fc0a2dfe29aaeda496478ea800c97bdfe4e22ba1b89720",
      "to": "sha256:b777d31d70d7bd260d746692d440ae30a28d3d90ba994f94d264d015d0f79a27",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ],
        "RR-C1b": [
          "RR-C1b"
        ],
        "RR-C2": [
          "RR-C2"
        ]
      }
    },
    {
      "from": "sha256:b777d31d70d7bd260d746692d440ae30a28d3d90ba994f94d264d015d0f79a27",
      "to": "sha256:4af8e55ed11edfc3ea0c36181c13f11dd5742108f4f38d57f0b15bdf580b752b",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ],
        "RR-C1b": [
          "RR-C1b"
        ],
        "RR-C2": [
          "RR-C2"
        ]
      }
    },
    {
      "from": "sha256:4af8e55ed11edfc3ea0c36181c13f11dd5742108f4f38d57f0b15bdf580b752b",
      "to": "sha256:e6df2b3ac2271259efedd93c96a79966b915548e2cd37747833ccc21905b1260",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ],
        "RR-C1b": [
          "RR-C1b"
        ],
        "RR-C2": [
          "RR-C2"
        ],
        "RR-C3": [
          "RR-C3"
        ]
      }
    },
    {
      "from": "sha256:e6df2b3ac2271259efedd93c96a79966b915548e2cd37747833ccc21905b1260",
      "to": "sha256:129cae8ef3792d9ae5a007ed37c6d3df9bb3b1712d5c4b452690b42fe42b87f6",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ],
        "RR-C1b": [
          "RR-C1b"
        ],
        "RR-C2": [
          "RR-C2"
        ],
        "RR-C3": [
          "RR-C3"
        ]
      }
    },
    {
      "from": "sha256:129cae8ef3792d9ae5a007ed37c6d3df9bb3b1712d5c4b452690b42fe42b87f6",
      "to": "sha256:3ae2907dd73b212ed47dc6c3e373191338817deaa79e8f1aa1456e865b45dae6",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ],
        "RR-C1b": [
          "RR-C1b"
        ],
        "RR-C2": [
          "RR-C2"
        ],
        "RR-C3": [
          "RR-C3"
        ]
      }
    },
    {
      "from": "sha256:3ae2907dd73b212ed47dc6c3e373191338817deaa79e8f1aa1456e865b45dae6",
      "to": "sha256:09259ca04468131aff9125326c043078881ec4190cafa38628483388cb98a239",
      "chunk_ids": {
        "RR-C1a": [
          "RR-C1a"
        ],
        "RR-C1b": [
          "RR-C1b"
        ],
        "RR-C2": [
          "RR-C2"
        ],
        "RR-C3": [
          "RR-C3"
        ],
        "RR-C4": [
          "RR-C4"
        ]
      }
    }
  ]
}
```

## RR-C1a pickup

The RR-C1a review returned 11 advisory findings and no blocking finding.
The findings collapse to 3 repair targets.

### Standards

Count: 4. Worst issue: C1a-S1, the unread `Step.command` field.

- C1a-S4, auto-fix, confidence 3: `internal/refusalroute/registry.go` `newIn` has no doc comment that names the injectable inventory seam.
- C1a-S1, no-op, confidence 6: ticket 03 owns the reader of `Step.command`.
- C1a-S2 and C1a-S3, no-op, confidence 4: each follows the approved design or the repository precedent.

### Spec

Count: 3. Worst issue: C1a-P1, the sample fill of an operator slot.

- C1a-P3, auto-fix, confidence 3: `internal/refusalroute/route.go` `Face.Render` joins `Facts.Preface` with no line-safe gate. C1a-C3 names the same fix.
- C1a-P1, no-op, confidence 6: ticket 03 owns the sample fill of the guard check.
- C1a-P2, no-op, confidence 3: the unregistered sentence obeys RR07.

### Coverage

Count: 4. Worst issue: C1a-C1, the untested empty `Composed` value.

- C1a-C1, auto-fix, confidence 7: `internal/refusalroute/route.go` has no test for an empty `Composed` value, so a mutation that drops the empty-value check survives.
- C1a-C3, auto-fix, confidence 4: the same fix as C1a-P3.
- C1a-C2, no-op, confidence 5: ticket 03 owns the sample fill of the guard check.
- C1a-C4, no-op, confidence 3: the spec does not state the empty-name edge.

### Repair state

RR-C1a consumed 1 of its 2 repair cycles.
The repair commit 5e627217 closes C1a-P3, C1a-C3, C1a-C1, and C1a-S4.

### Confirming round

The confirming round confirmed every fold on all three axes and returned no new finding.
The axes retained this optional advice:

- The `Fact` case and `asWritten` spell the same line-safe predicate.
- A preface of only whitespace renders as written.

## RR-C1b pickup

The RR-C1b review returned 16 findings, and Coverage marked two of them blocking.
The orchestrator also holds one known defect: `TestLandingFacesFollowTheirRoutes/source-not-clean` failed once and passed on a rerun at the same source.
The findings collapse to 9 repair targets across tickets 02, 03, and 04, and one ticket 05 expansion.

### Standards

Count: 6. Worst issue: C1b-S2, the re-spelled marker and joiner with no recorded red.

- C1b-S1, auto-fix, confidence 4: comments cite a feature identifier as provenance. Ticket 02 repair.
- C1b-S2, auto-fix, confidence 5: the marker gets a recorded red through the C1b-P3 assertions, and the joiner gets one source. Tickets 02 and 03 repair.
- C1b-S3, auto-fix, confidence 6: an over-long comment line and an unclear closure name in the landing route test. Ticket 02 repair.
- C1b-S4, auto-fix, confidence 4: three test parsers of the printed record layout collapse to one. Ticket 03 repair.
- C1b-S5, no-op, confidence 4: the parameter list is a preference with no binding requirement, so the axis keeps it as advice.
- C1b-S6, no-op, confidence 3: all three axes judge the operational AXI exemption defensible, and it has precedent.

### Spec

Count: 3. Worst issue: C1b-P1, the routeless identity-component refusals.

- C1b-P1, auto-fix, confidence 6: the five identity components without a recovery route hand back through `land-handback`. Ticket 02 repair.
- C1b-P2, auto-fix, confidence 5: ticket 05 expands, so each remaining `landReviewed` error hands back through `land-handback`.
- C1b-P3, auto-fix, confidence 7: each reviewer land face gets a literal marker assertion and a recorded probe. Ticket 02 repair.

### Coverage

Count: 7. Worst issue: C1b-C1, the unsafe-path lookup that refuses its own input class.

- C1b-C1, no-op, confidence 8: Bench builds each pool path from hex and decimal segments, so the lookup guard is defensive. A consultation verified the pool path shape, and the spec decided the form.
- C1b-C2, auto-fix, confidence 7: the same fix as C1b-P1 and C1b-P2.
- C1b-C3, auto-fix, confidence 7: the follow walk proves that the rerun finishes, and the source reruns print the repaired source tip. Ticket 03 repair.
- C1b-C4, auto-fix, confidence 6: the review step becomes an instruction, so the ticket 05 `land-red` fixture can follow it. Ticket 02 repair.
- C1b-C5, no-op, confidence 5: the spec states the sample-path limit of the guard check.
- C1b-C6, auto-fix, confidence 5: the newly faced handback sites get `next=` assertions. Ticket 02 repair.
- C1b-C7, auto-fix, confidence 5: a test dispatches `bench recovery` through the command registry. Ticket 04 repair.

### Known defect

- The follow walk fails at random on the `source-not-clean` step, probably because a new wrapper script starts while it is still open for write. Ticket 03 repair.

### RR-C1b repair state

RR-C1b consumed 1 of its 2 repair cycles.
The ticket 02 repair 5a0a4f78, the ticket 03 repair c359944c, and the ticket 04 repair f4669147 close every auto-fix finding and the known defect.
The ticket 03 repair also adds a review step and the repaired source tip to the `source-not-clean` route.
It also hands back a dirty source on the resume path.

## Diagnosed causes

The repair charges of RR-C1a and RR-C1b did not run the debug step before each fix.
The orchestrator records each diagnosed cause here from the repair returns, after the repairs.

| finding | diagnosed cause | repro |
|---|---|---|
| C1a-P3, C1a-C3 | `Face.Render` joined the preface with no line-safe gate. | A preface with a control byte rendered raw. |
| C1a-C1 | No test held an absent composed value. | A probe that drops the empty-value check survived. |
| C1b-P1, C1b-C2 | The handback fallback skipped every refusal that named a component, and only the request component sets a route. | A probe that removes the fallback failed 14 tests. |
| C1b-C4 | The review step was a command step, but a review is a phase and not a shell command. | The follow walk cannot run a command step whose head is not `bench`. |
| C1b-P3 | The walk took the expected marker from the registry authority. | A probe that gives `destination-not-clean` agent authority stayed green before the fix. After the fix, the ticket 02 repair session ran the same probe against `TestReviewerLandFacesOpenWithTheMarker`: the verdict was `bit`, the subtest `destination-not-clean` failed, and the restore passed. |
| C1b-C3 | The source reruns kept the caller's tip, and the commit in the route moves the tip. | A probe that restores the caller tip made two reruns hit the tip mismatch. |
| Known flake | A parallel fork kept the new wrapper script open for write, so its start failed with "text file busy". | A scratch program failed about 300 of 2000 starts. |
| C1b-C7 | No test ran `bench recovery` through the command dispatch. | A probe that points the row at another handler stayed green before the fix. |

## RR-C1b confirming round

The confirming round confirmed every fold and returned three findings.

- C1b-RP1 and C1b-RC1, evidence-only, confidence 6: the record did not show the red of the reviewer-marker test. The cause table now holds the probe result that the ticket 02 repair return reported.
- C1b-RS1, auto-fix, confidence 6: a resume refusal test spells the resume rerun inline, beside the `resumeRerunOf` helper. Ticket 03 repair, cycle 2 of 2.

The orchestrator also corrects the seam cells of RR15, RR18, RR19, RR58, and RR60 in the spec, because their tests moved files.

### RR-C1b repair cycle 2

RR-C1b consumed 2 of its 2 repair cycles, and no further repair cycle remains without a reviewer extension.
The ticket 03 repair session ran the debug step and recorded this cause before its fix.

| finding | diagnosed cause | repro |
|---|---|---|
| C1b-RS1 | The cycle 1 repair added `resumeRerunOf` but left one resume test with its own inline rerun. | A probe that changes the helper form was silent before the fix and bit after it. |

The repair commit 33ac2f98 closes C1b-RS1.

### RR-C1b repair cycle 2 confirming round

The confirming round confirmed C1b-RS1, C1b-RP1, and C1b-RC1 on all three axes and returned no new finding.
The axes retained this optional advice: the resume form also appears in two test files outside this delta.

## RR-C2 pickup

The RR-C2 review returned 18 findings, and Coverage marked one of them blocking.
The findings collapse to repair targets in tickets 05, 06, and 07.
Each repair session runs the debug step and states the diagnosed cause before its fix.

### Standards

Count: 6. Worst issue: C2-S1, the hand-merge fixture step pasted three times.

- C2-S1, auto-fix, confidence 7: one helper serves the hand-merge fixture step. Ticket 05 repair.
- C2-S2, auto-fix, confidence 6: the `mergeRedRefusal` comment states the grade without the parameter name as a verb. Ticket 06 repair.
- C2-S3, auto-fix, confidence 8: wrap the comment line and add the blank line. Ticket 07 repair.
- C2-S4, auto-fix, confidence 4: `planLandedAssignment` states its landed precondition. Ticket 07 repair.
- C2-S5, auto-fix, confidence 3: `releaseAssignment` builds the missing-tree release once. Ticket 07 repair.
- C2-S6, auto-fix, confidence 3: the whole-gate fixture reuses `mergeSetAt`. Ticket 06 repair.

### Spec

Count: 5. Worst issue: C2-P2, the whole-route placeholder.

- C2-P1, no-op, confidence 6: C2-2 adds the selective-lane test.
- C2-P2, auto-fix, confidence 5: a non-line-safe missing-tree path prints a placeholder for each value, as the Edge inventory states. Ticket 07 repair.
- C2-P3, no-op, confidence 4: `create` is outside the six write verbs, so a learning records its exec route.
- C2-P4, auto-fix, confidence 4: a test proves that a landing retires an absent landed sibling. Ticket 07 repair.
- C2-P5, auto-fix, confidence 3: the ticket 07 repair checks the landed route choice against the clean selector. It fixes a gap that a repro shows.

### Coverage

Count: 7. Worst issue: C2-1, the pruned registration.

- C2-1, auto-fix, confidence 6: a missing tree whose registration Git already pruned must clear through the printed route. Ticket 07 repair.
- C2-2, auto-fix, confidence 6: a selective-lane fixture grades the target-alone base. Ticket 06 repair.
- C2-3, auto-fix, confidence 6: the walk follows the unlanded missing-tree cause from its printed route. Ticket 07 repair.
- C2-4, auto-fix, confidence 6: tests hold the placeholder arms of the merge rerun (ticket 05 repair) and the reset rerun (ticket 07 repair).
- C2-5, auto-fix, confidence 5: fixtures produce the two infrastructure arms of the target-alone grade. Ticket 06 repair.
- C2-6, auto-fix, confidence 5: a test holds the unknown-lease arm of the missing-tree release. Ticket 07 repair.
- C2-7, no-op, confidence 4: an unfaced merge refusal hands back, which is the fail-closed default of the spec.

### RR-C2 repair cycle 1

RR-C2 consumed 1 of its 2 repair cycles. The repair commits are eadc1b23 (ticket 05), e686bbda (ticket 06), d0c8c495 and fba4fe78 (ticket 07).
The repair sessions ran the debug step and stated these causes before their fixes.

| finding | diagnosed cause | repro |
|---|---|---|
| C2-S1 | Three tests wrote the hand-merge step inline, and two of them repeated the same resolution text. | `rg -n '"merge", "--no-commit"' internal/worktree` showed three conflict-merge sites. |
| C2-2 | The lane fixtures used a constant check that reads no change list, so no check depended on the base of the target-alone grade. | A `bench probe` that graded against `PreviousTip` instead of `Incoming` returned `silent`. |
| C2-5 | No fixture made the target-alone grade return `Infrastructure`: neither a failed lane read in `mergeTargetGrade` nor an unreadable target tree in `GradeTarget`. | A `bench probe` that swapped each arm's `Infrastructure` literal returned `silent`. |
| C2-S2 | The comment used the parameter name `grade` as a verb. | A read of `merge_refusal.go`. |
| C2-S6 | `mergeSetAt` cannot build a set without a lane commit, so the whole-gate fixture repeated the set literal. | A read of `merge_caller_root_test.go`. |
| C2-4 (merge half) | No worktree test drove a merge refusal with a value that is not line-safe, so nothing held the placeholder arm of `mergeRerun`. | A `bench probe` that printed the raw values in `mergeRerun` returned `silent` before the fix. |
| C2-1 | `releaseRegistration` failed when Git had already pruned the registration: a missing pool or a missing entry returned an error, so the release-leftover plan never finished. | The walk's `landed pruned` and `unlanded pruned` causes, red before the edit. |
| C2-3 | No production defect: the walk had no unlanded missing-tree fixture. | A `bench probe` on `path.go` that forced the landed route. |
| C2-6 | No production defect: the `LeaseUnknown` arm had no test. | A `bench probe` that removed `case LeaseUnknown:`. |
| C2-P2 | `missingTreeRefusal` passed the whole route as one `FactRecovery` value, so the `Composed` rule printed one placeholder for the whole route. | `TestMissingTreeRoutePrintsAPlaceholderInEachUnsafeSlot`, red before the edit. |
| C2-4 (reset half) | No production defect: the reset rerun placeholder arm had no test. | A `bench probe` that printed the raw reset target. |
| C2-P4 | No production defect: the absent landed sibling arm had no test. | A `bench probe` that disabled the absent arm of `planLandedAssignment`. |
| C2-P5 | The refusal chose `clean --landed` from branch ancestry alone, but the selector also requires an active record and no live lease. | `TestMissingTreeRouteAgreesWithTheLandedSelector`, red before the edit. |
| C2-S3, C2-S4, C2-S5 | Standards findings, applied as listed. | A read of the named lines. |
| C2-P5 (list follow-up) | The `list` help row chose the missing-tree route from the branch-ancestry cell alone, so it named `clean --landed` where the refusal named the release. | `TestMissingTreeRouteAgreesWithTheLandedSelector`, red before the edit. |

The ticket 07 repair left the `list` help row on the old landed rule, outside its fence. The orchestrator expanded the fence, and a fresh session moved both surfaces to one landed decision.

### RR-C2 repair cycle 1 red verdicts after the fix

The confirming round asked for the red verdict after each fix (C2-RC1). The repair returns hold these verdicts, and each probe restored its file.

| finding | probe after the fix | verdict |
|---|---|---|
| C2-2 | The target-alone grade ignores the lane (`landing.New()` for `mergeOwner`). | `bit`, 4 failed, including both named-lane rows |
| C2-2 | The target-alone grade reads `PreviousTip` instead of `Incoming`. | `bit`, 1 failed: `merge-target-red/named_lane` |
| C2-3 | `path.go` forces the landed route. | `bit`: the assignment stays active after its route |
| C2-4 (merge half) | `mergeRerun` prints the raw values. | `bit`, 1 failed: the new placeholder test |
| C2-4 (reset half) | The reset rerun prints the raw target. | `bit`: the route printed `<rerun>` |
| C2-5 | The lane-read arm returns `LanePass` instead of `Infrastructure`. | `bit`, 1 failed: `target_lane_unreadable` |
| C2-5 | The unreadable-tree arm returns `Green` instead of `Infrastructure`. | `bit`, 1 failed: `target_tree_unreadable` |
| C2-6 | The `case LeaseUnknown:` arm is removed. | `bit`: the release exited 0 |
| C2-P4 | The absent arm of `planLandedAssignment` is disabled. | `bit`: the absent sibling stays active |
| C2-P5 (list) | The help row reads the landed cell again. | `bit`, the leased landed row failed |

## RR-C3 pickup

The RR-C3 review returned 14 findings, and two of them are blocking.
The findings collapse to repair targets in tickets 08 and 09, and to one plan correction.
Each repair session runs the debug step and states the diagnosed cause before its fix.
The ticket 08 repair runs first, because the ticket 09 repair builds on its shared walk.

### Standards

Count: 4. Worst issue: C3-S1, the route-walk harness that the commit tests repeat.

- C3-S1, auto-fix, confidence 7: one shared test package, `internal/refusalroute/routetest`, holds the route walk, and the worktree and commit walks use it. Ticket 08 repair.
- C3-S2, auto-fix, confidence 7: `landingRerunArg` uses `refusalroute.Arg` and keeps no copy of its rule. Ticket 08 repair.
- C3-S3, no-op, confidence 5: the usage grammar and the face route are two facts, and the spec fixes that the face declares no sentence. Flagged for reviewer veto.
- C3-S4, auto-fix, confidence 6: the spec states that the commit prints `next=` on its own stderr line after the refusal sentence. Plan correction.

### Spec

Count: 5. Worst issue: C3-P1, the agent-clearable causes that route to the reviewer.

- C3-P1, auto-fix, confidence 6: each landing cause that an agent clears with its own edits takes an agent face, by the spec authority rule. The handback keeps only causes outside agent authority. Ticket 09 repair.
- C3-P2, no-op, confidence 6: the same item as C3-S3.
- C3-P3, auto-fix, confidence 7: the same target as C3-S1. Ticket 08 repair.
- C3-P4, auto-fix, confidence 6: the same target as C3-S2. Ticket 08 repair.
- C3-P5, auto-fix, confidence 7: hostile-argument rows hold the commit rerun. Ticket 09 repair.

### Coverage

Count: 5. Worst issue: C3-2, the same item as C3-P1.

- C3-1, auto-fix, confidence 7: the same target as C3-P5, with rows for `--dry-run`, `--preflight-build`, a multi-line message, and an unsafe path. Ticket 09 repair.
- C3-2, auto-fix, confidence 5: the same target as C3-P1. Ticket 09 repair.
- C3-3, no-op, confidence 6: ticket 11 widens the commit walk to the commitment faces.
- C3-4, auto-fix, confidence 6: a lane-red fixture and a dry-run red fixture reach `commit-red`. Ticket 09 repair.
- C3-5, auto-fix, confidence 5: a fixture with no owning assignment prints the `<label>` placeholder. Ticket 09 repair.

### RR-C3 repair cycle 1

RR-C3 consumed 1 of its 2 repair cycles. The repair commits are 0bc09add and 12c02202 (ticket 08), and 0d0e0c1f (ticket 09).
The repair sessions ran the debug step and stated these causes before their fixes.

| finding | diagnosed cause | repro |
|---|---|---|
| C3-S1, C3-P3 | The route-walk rules lived only in worktree test files, which no other package can import, so the commit walk wrote them again. | A `bench probe` on the worktree marker check, run against the commit walk, returned `silent`. |
| C3-S2, C3-P4 | `landingRerunArg` kept its own placeholder rule beside `refusalroute.Arg`. | A `bench probe` on `pasteable` in `route.go`, run against the merge placeholder test, returned `silent`. |
| C3-S1, C3-S2 (follow-up) | Three callers outside the first repair fence kept the old copies: a reviewer marker constant, three walk aliases, and a bracket trim in `landingRerunArg`. | `rg` over `internal` for the marker, the aliases, and `landingRerunArg`. |
| C3-P1, C3-2 | The landing package had no typed signal for a cause that the caller clears by correcting its named paths, so `landingFace` sent each such error to the handback. A missing path also printed no `next=` line. | `TestCommitRerunKeepsTheCallersArguments` at the start tip: 6 of 6 rows red, with no `next=`. |
| C3-P5, C3-1 | No row held the argument composition of `refuse`; the composition was correct once a refusal reached it. | One `bench probe` for each argument rule. |
| C3-4 | No fixture drove a lane red or a dry-run red; the code already routed both to `commit-red`. | Two `bench probe` runs on the walk. |
| C3-5 | No fixture covered a commit that no assignment owns; the code already printed `<label>`. | A `bench probe` on the owner label. |

The ticket 09 repair keeps the malformed lane declaration in `commit-handback`, because an agent edit of a gate check is the reviewer's call.

### RR-C3 confirming round

The confirming round confirmed every fold for the named-path causes, and it returned one blocking cause on all three axes.

- C3-RC1, C3-RS2, C3-RP1, auto-fix, confidence 6: the compare-and-swap refusal of the landing commit tells the caller to rerun, but it routes to `commit-handback`. A rerun clears it inside the agent's own worktree, so it takes the agent face `commit-tip-moved` with the rerun route, as `checkpoint-tip-moved` does. The landing package gives it a typed signal. Ticket 09 repair, cycle 2.
- C3-RS1, auto-fix, confidence 6: the spec Authority section now names an edit of a gate check, such as the lane declaration, as a reviewer cause. Plan correction, flagged for reviewer veto.
- C3-RP2, auto-fix, confidence 7: RR36 and the ticket 09 row now say one producing fixture for each face and cause. Plan correction.

The face table now lists `commit-named-path` and `commit-tip-moved`.
RR-C3 repair cycle 2 is the last cycle that the bounded repair policy allows without a reviewer extension.

### RR-C3 repair cycle 2

RR-C3 consumed 2 of its 2 repair cycles. The repair commit is 0f2825e1 (ticket 09).
The repair session ran the debug step and stated this cause before its fix.

| finding | diagnosed cause | repro |
|---|---|---|
| C3-RC1, C3-RS2, C3-RP1 | `destinationUpdateFailure` returned the compare-and-swap refusal as an untyped error, so `landingFace` had no type to select and the refusal fell to `commit-handback`. | The `commit-tip-moved` walk fixture printed the reviewer handback route before the fix. |

A `bench probe` that removed the typed wrap after the fix returned `bit`: the fixture printed the reviewer handback route again.

## RR-C4 pickup

The RR-C4 review returned 16 findings, and Coverage marked one of them blocking.
The findings collapse to repair targets in ticket 10 and to two plan corrections.
The repair session runs the debug step and states the diagnosed cause before its fix.

### Standards

Count: 6. Worst issue: C4-S2, the capture faults that route to the reviewer handback.

- C4-S1, no-op, confidence 6: the raising site owns the forced `--fresh` rerun. The registry has no optional word.
- C4-S2, auto-fix, confidence 7: `completionTree` splits capture and read faults from proof faults, and only a proof fault takes `gate-handback`. Ticket 10 repair.
- C4-S3, auto-fix, confidence 6: `routetest` owns the `next=` line reader, and the gate and commit tests use it. Ticket 10 repair.
- C4-S4, auto-fix, confidence 4: no route prints the `<arguments>` placeholder for a rerun that has no arguments. Ticket 10 repair, with C4-2.
- C4-S5, no-op, confidence 6: the staged native-record-operations spec cites a deleted symbol; its own staleness pass corrects the citation. A learning records it.
- C4-S6, auto-fix, confidence 5: the new gate test comments drop the spec row ids. Ticket 10 repair.

### Spec

Count: 4. Worst issue: C4-P1, the same cause as C4-S2.

- C4-P1, auto-fix, confidence 7: the same target as C4-S2. Ticket 10 repair.
- C4-P2, no-op, confidence 6: the engine refusal serves every gate run. The spec won't-handle row for run-state refusals covers it.
- C4-P3, auto-fix, confidence 6: the slug renders through a quoted fact, and RR40 now reads `next=bench preflight review 'example'`. Ticket 10 repair and plan correction.
- C4-P4, auto-fix, confidence 6: the same target as C4-2. Ticket 10 repair.

### Coverage

Count: 6. Worst issue: C4-2, the gate route inside a land or commit authorization run.

- C4-1, auto-fix, confidence 7: the same target as C4-P3, with fixtures for a hostile slug, spec path, and chunk id. Ticket 10 repair.
- C4-2, auto-fix, confidence 6: a funnel refusal inside a land or commit run prints no gate route. The calling verb prints the one route (story 35). A test drives such a run and counts one `next=` line. Ticket 10 repair.
- C4-3, auto-fix, confidence 6: the same target as C4-S2, with a test that injects a capture fault. Ticket 10 repair.
- C4-4, no-op, confidence 5: the same item as C4-P2.
- C4-5, auto-fix, confidence 6: a fixture runs a refused `--fresh` checkpoint and holds `--fresh` in the rerun. Ticket 10 repair.
- C4-6, no-op, confidence 4: the handback and tip-moved walks need no repair step by their nature.

### RR-C4 repair cycle 1

RR-C4 consumed 1 of its 2 repair cycles. The repair commit is 819154c6 (ticket 10).
The repair session ran the debug step and stated these causes before its fixes.

| finding | diagnosed cause | repro |
|---|---|---|
| C4-2, C4-P4, C4-S4 | The funnel served both the gate verb and the authorization entry, and `refuse` printed the gate route in every run. A commit or land run printed a second route, and a tickets-only landing printed `<arguments>`. | `TestAuthorizationRunPrintsTheCallersRoute` printed two `next=` lines before the edit. |
| C4-S2, C4-P1, C4-3 | `applyCheckpoint` wrapped every `completionTree` error as a proof fault, so read and capture faults took `gate-handback`. | `TestCompletionFaultsRouteByTheirCause/unreadable_source` printed the reviewer route before the edit. |
| C4-P3, C4-1 | The gate face declared the slug as a composed value, which prints as written. | `TestCheckpointRouteRendersHostileFacts/spaced_slug` printed the raw slug before the edit. |
| C4-5 | No test graded the `mode == forceRun` term. | A `bench probe` that dropped the term returned `bit` after the new row. |
| C4-S3, C4-S6 | The `next=` reader had three copies, and new comments carried spec row ids. | A read of the gate and commit tests. |

Only the `bench gate` verb prints a gate route now; the land and commit verbs print their own route for a funnel refusal.
