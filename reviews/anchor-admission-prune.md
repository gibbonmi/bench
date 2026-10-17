# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/anchor-admission-prune/spec.md",
  "plan_digest": "sha256:294bb59a470b025b0faa1f465f8561a826966e0879b0b69714709866831db182",
  "implementation_session": "",
  "chunks": [
    {
      "id": "AP-C2",
      "base": "ca0e6ec76eac0a2ef65d22c1c5405b948bb15dcf",
      "tip": "1e169d060fc785115f5b6f42eff6a83761bae863",
      "plan_digest": "sha256:294bb59a470b025b0faa1f465f8561a826966e0879b0b69714709866831db182",
      "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
      "acceptance_rows": [
        "AP02",
        "AP03",
        "AP04",
        "AP10",
        "AP25",
        "AP26",
        "AP27",
        "AP31",
        "AP46"
      ],
      "verification": [
        {
          "id": "t1-t1-anchors",
          "performer": "claude:aap_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t01",
            "digest": "sha256:bc1cf1f92ded0de4ca5c3c36229d3c25cdecbb3a885206be420b6f556112370d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,1723,181\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t1-t1-docs-currency-workflow",
          "performer": "claude:aap_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t01",
            "digest": "sha256:3edeee72ada88b9370cf9fdbb32754a73daff4a95a2e70e45ab119c6b780b173",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4600,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-docs-currency-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t1-t1-conformance",
          "performer": "claude:aap_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t01",
            "digest": "sha256:56e103fd930c1fdf24313d2c68701c5fcebf4aa6138e998f454516c68cdfef31",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,51490,1561\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-t1-cmd-bench",
          "performer": "claude:aap_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t01",
            "digest": "sha256:f61428016e9012f55ecea1e9090d037635d28457632dbc4c78e2dc1a6fc59366",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14734,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1-t1-root-conformance",
          "performer": "claude:aap_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t01",
            "digest": "sha256:fa8c8775eb396a0cee658a7df678d7e7f43fd85231509c8fc6e09314f6c8009d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14785,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "apc2-standards-r1",
          "performer": "claude:aap_apc2_standards_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:aap_apc2_standards_r1",
            "digest": "sha256:2b738dc2d2c67af8808da5a3a392bbbbd442541779eb5b2184eada61cfb8ccab",
            "excerpt": "Standards axis, AP-C2, frozen pair ca0e6ec7..32923169; check-current current=true, dirty_record=false.\nS1 | low | internal/anchors/registry_retained_workflow.go:9-10 | the family header still calls defaultWorkflowAnchors \"the default contract below\" and delegatedWorkflowAnchors \"the opt-in delegated exception\"; both groups now hold only Forbid rows for retired text | craft-comments Aging | auto-fix | confidence 7\nS2 | low | internal/anchors/registry_retained_workflow.go:153-155 | the laneAndLandingAnchors header says the family pins \"`bench worktree land` runs the whole-project gate\"; the delta removed the only row that pinned it | craft-comments Aging | auto-fix | confidence 7\nS3 | low | internal/conformance/retained_workflow_test.go:44-58 | checkRetainedWorkflow writes its own EvaluateGroup loop and prefix filter beside checkWorkflowAnchors | BENCH.md invariant 4, AGENTS.md one source | ask-user | confidence 4\nCount 3, all low. Worst: S2.\n"
          },
          "axis": "Standards",
          "base": "ca0e6ec76eac0a2ef65d22c1c5405b948bb15dcf",
          "tip": "1e169d060fc785115f5b6f42eff6a83761bae863",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "apc2-spec-r1",
          "performer": "claude:aap_apc2_spec_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_apc2_spec_r1",
            "digest": "sha256:08e07cf3bff63d8b2037343ab9988223bc6b4b418efdab9e4a505869d0ef51d2",
            "excerpt": "Spec axis, AP-C2, tip 32923169; check-current current=true, dirty_record=false. All nine covered rows hold (AP02, AP03, AP04, AP10, AP25, AP26, AP27, AP31, AP46). Census sample of nine removed rows: each DROP in classified.md; 15 stale phrases each have a Forbid row at registry_retained_workflow.go:21-35; 25 deleted fixtures match fixture-census.md:8-32.\nP1 | low | the moved stale phrases now match through the registry matcher, which strips HTML comments and collapses whitespace; spec.md:139 requires only a Forbid row on the same file | no-op, flagged for reviewer veto | confidence 7\nP2 | low | removal of checkRetainedWorkflow from the aggregate (docs_workflow_checks_test.go:30) keeps the document reads; checkWorkflowAnchors already evaluates the group | no-op | confidence 8\nP3 | low | registry_retained_workflow.go:164 needle holds two sentences at 23 words; story 17 intent versus the AP17 word bound | ask-user | confidence 5\nCount 3, all low. Worst: P1.\n"
          },
          "axis": "Spec",
          "base": "ca0e6ec76eac0a2ef65d22c1c5405b948bb15dcf",
          "tip": "1e169d060fc785115f5b6f42eff6a83761bae863",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "apc2-coverage-r1",
          "performer": "claude:aap_apc2_coverage_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_apc2_coverage_r1",
            "digest": "sha256:e86fa20b61e715e2f8812d195fb9b28ef9e45d37e3b8865e5c5078d5ce116cab",
            "excerpt": "Coverage axis, AP-C2, tree at 32923169; check-current current=true, dirty_record=false. 2 low findings, no blocker.\nC1 | low | internal/anchors/registry.go:234 with registry_retained_workflow.go:33-34 | the 15 stale phrases now pass the HTML-comment strip, so a stale instruction inside an HTML comment no longer reds | ask-user | confidence 7\nC2 | low | internal/conformance/retained_workflow_test.go:73 | the no-registry-literal subtest parses one hard-coded file name and matches exact literals only | no-op for this chunk | confidence 5\nVerified: each new Forbid row bites through runAnchorBites (docs_workflow_helpers_test.go:449); checkWorkflowAnchors still evaluates every row (docs_workflow_helpers_test.go:36,42); 25 deleted fixtures name removed rows.\nTests: bench test --package ./internal/conformance --run \"TestRetainedWorkflow|TestEveryRetainedFixtureBitesThroughRegisteredOwner|TestFreshTicketAuthors|TestImplementationContinuation|TestLineBinding\" pass, 707 tests; bench test --package ./cmd/bench --run TestAnchorsReportsAbsentNeedles pass, 1 test.\nCount 2. Worst: C1.\n"
          },
          "axis": "Coverage",
          "base": "ca0e6ec76eac0a2ef65d22c1c5405b948bb15dcf",
          "tip": "1e169d060fc785115f5b6f42eff6a83761bae863",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "apc2-standards-r2",
          "performer": "claude:aap_apc2_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e9b4216c4d59e10744f9ff0b06de7a959f815caa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_apc2_standards_r2",
            "digest": "sha256:a49549bcfd1a8b603d502eccde24ac5ee3aead9aa8f9f9067ea48152fe8d69d4",
            "excerpt": "Standards confirming round, AP-C2, correction delta c532d30d..ae3f7254 (comment-only; the file shrank by two lines).\nS1 closed, confidence 9: internal/anchors/registry_retained_workflow.go:9-10 reads \"the retired default and delegated workflow forms\", matching the Forbid-only groups at lines 181-195.\nS2 closed, confidence 9: registry_retained_workflow.go:153-159 no longer claims a landing-gate pin; each claim matches a row at lines 162-178.\nS3 not re-graded (no-op).\nNew findings: none. Verdict: pass.\n"
          },
          "axis": "Standards",
          "base": "ca0e6ec76eac0a2ef65d22c1c5405b948bb15dcf",
          "tip": "1e169d060fc785115f5b6f42eff6a83761bae863",
          "finding_ids": [],
          "supersedes": [
            "apc2-standards-r1"
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
  }
}
```

## Repair state

AP-C2: no repair cycle is consumed. The S1 and S2 correction is comment-only, so it is evidence-only under the bounded repair policy. No reviewer extension exists.

## Standards

AP-C2: 3 findings. The worst issue is S2.

- S1 (auto-fix, confidence 7): `internal/anchors/registry_retained_workflow.go:9-10` still calls the default and delegated groups a contract and an exception, and both groups now hold only Forbid rows. The orchestrator corrects the comment.
- S2 (auto-fix, confidence 7): `internal/anchors/registry_retained_workflow.go:153-155` says the family pins that `bench worktree land` runs the whole-project gate, and no remaining row pins it. The orchestrator corrects the comment.
- S3 (no-op, confidence 4): `internal/conformance/retained_workflow_test.go:44-58` keeps its own group loop. The spec says the check keeps its document reads, and the group set derives from the rows. Flagged for reviewer veto.

## Spec

AP-C2: 0 actionable findings. The worst issue is P1.

- P1 (no-op, confidence 7): the moved stale phrases match through the registry matcher, which strips HTML comments. The spec requires only a Forbid row on the same file. Flagged for reviewer veto.
- P2 (no-op, confidence 8): the aggregate check no longer calls `checkRetainedWorkflow`, because `checkWorkflowAnchors` already evaluates the group.
- P3 (no-op, confidence 5): `internal/anchors/registry_retained_workflow.go:164` holds two sentences in 23 words. AP17 bounds words, and the AP-C5 admission test grades it. Flagged for reviewer veto.

## Coverage

AP-C2: 0 actionable findings. The worst issue is C1.

- C1 (no-op, confidence 7): the same HTML-comment narrowing as P1. Merged with P1.
- C2 (no-op, confidence 5): `internal/conformance/retained_workflow_test.go:73` parses one fixed file name. AP25 holds as written.
