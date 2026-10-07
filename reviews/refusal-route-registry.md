# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/refusal-route-registry/spec.md",
  "plan_digest": "sha256:10d938cdde11624c65c40c632beebb04b986d194da031b1c8dd8d470d65aaa05",
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
