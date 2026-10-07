# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/refusal-route-registry/spec.md",
  "plan_digest": "sha256:f0639cf27495f50d19843b401bf84eab059aa614e05f76f81c98dc485790924b",
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
