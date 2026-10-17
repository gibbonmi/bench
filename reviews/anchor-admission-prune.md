# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/anchor-admission-prune/spec.md",
  "plan_digest": "sha256:9ae1c28cf65da60d226e873b3fb29a0ac2ef2a7eb9b940da5a78c04c69067181",
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
    },
    {
      "id": "AP-C1",
      "base": "bba9f7e68856f827f535943983102164c6a79061",
      "tip": "8dee1a027b4b1a32731bffa71f620d3c31a10413",
      "plan_digest": "sha256:9ae1c28cf65da60d226e873b3fb29a0ac2ef2a7eb9b940da5a78c04c69067181",
      "source_digest": "456c681ccc9744bd468cd19e7fe78d75aaee79ec",
      "acceptance_rows": [
        "AP02",
        "AP03",
        "AP04",
        "AP05",
        "AP06",
        "AP07",
        "AP08",
        "AP10",
        "AP18",
        "AP20",
        "AP21",
        "AP22",
        "AP23",
        "AP24",
        "AP46",
        "AP47"
      ],
      "verification": [
        {
          "id": "t2-t2-anchors",
          "performer": "claude:aap_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "456c681ccc9744bd468cd19e7fe78d75aaee79ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t02",
            "digest": "sha256:08cd94211f1a3f72d93cbb37b5ced41e131e73c77bedecb52331f3aced4321cb",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,1308,165\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t2-t2-conformance",
          "performer": "claude:aap_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "456c681ccc9744bd468cd19e7fe78d75aaee79ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t02",
            "digest": "sha256:d9b065011e4bea383207cddba3500531f2d8c08570183303036e257be725691a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,41124,1377\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t2-t2-docs-currency-workflow",
          "performer": "claude:aap_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "456c681ccc9744bd468cd19e7fe78d75aaee79ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t02",
            "digest": "sha256:44b449b364b256422bb4df596dfce6956f10f18bb90a2a8d8ce663d444bb0b61",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4483,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-docs-currency-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t2-t2-conformance-canary-families",
          "performer": "claude:aap_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "456c681ccc9744bd468cd19e7fe78d75aaee79ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t02",
            "digest": "sha256:f0c7d19ef114ea1891fb63e4ce57b36704bb46c5f96c58cb1d59c40d7f775b06",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  conformance-canary-families,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,10,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-conformance-canary-families",
          "command": "bench test --check conformance-canary-families",
          "exit_code": 0
        },
        {
          "id": "t2-t2-cmd-bench",
          "performer": "claude:aap_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "456c681ccc9744bd468cd19e7fe78d75aaee79ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t02",
            "digest": "sha256:c7d65416ca8922a62434f7f7b6dbb59a4e7bdb2d8d331df4bac792240fe3d663",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14193,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t2-t2-root-conformance",
          "performer": "claude:aap_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "456c681ccc9744bd468cd19e7fe78d75aaee79ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:aap_t02",
            "digest": "sha256:cc44ee14699f14e137b11b62e60e5737df42460d5588d85ffe350505b6ae0f18",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14490,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "apc1-standards-r1",
          "performer": "claude:aap_apc1_standards_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "456c681ccc9744bd468cd19e7fe78d75aaee79ec",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:aap_apc1_standards_r1",
            "digest": "sha256:f42ba4eb736183c1c3b1aa36ee883a59878b55f270b8cf44b909da42357875fb",
            "excerpt": "Standards axis, AP-C1, pair bba9f7e6..18f6c6dd; check-current current=true, dirty_record=false.\nS1 | low | internal/anchors/registry_data_test.go:352 | TestDelegateExecOnlyAndCapChangeAnchorsRedOnRemoval still names the cap-change rule that the prune removed; the function holds only the exec-only rule | craft-comments SKILL.md:16-17 and :68-70 | ask-user | confidence 6\nClean: folded marker comment (registry_data.go:5-6), staleLandingPathAnchors single source (:26-59), drain family comment (:130), ft311_preparation_test.go removals.\nCount 1. Worst: S1.\n"
          },
          "axis": "Standards",
          "base": "bba9f7e68856f827f535943983102164c6a79061",
          "tip": "8dee1a027b4b1a32731bffa71f620d3c31a10413",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "apc1-spec-r1",
          "performer": "claude:aap_apc1_spec_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "456c681ccc9744bd468cd19e7fe78d75aaee79ec",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:aap_apc1_spec_r1",
            "digest": "sha256:d8fa3e2cc50c1fc03deb4c43cbd36c1c892eeffa7e9cd83eeaa22b0f87221c82",
            "excerpt": "Spec axis, AP-C1, base bba9f7e6, tip 18f6c6dd; check-current current=true. All covered rows audited; none of the eight author calls contradicts the spec or the ticket.\nP1 | low | ticket commit 8dee1a02 edited specs/anchor-admission-prune/spec.md seam cells AP06-AP08 (lines 214-216) and planning seam-8 (line 410), outside the ticket Writes line | ask-user | confidence 7\nP2 | low | internal/anchors/registry_data.go:92 row 133 diagnostic still names \"does not create or retain the integration worktree\", which the re-cut needle no longer holds | auto-fix with tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight/EXPECT | confidence 6\nCensus (full enumeration): 333 DROP rows and the Roles row removed; no FORBID, LABEL, ROUTE, or COMMAND row left; 155 fixtures deleted (154 listed plus implement-spec-entry-validation); 9 integration-source files covered; no kept Require needle over 25 words.\nCount 2. Worst: P2.\n"
          },
          "axis": "Spec",
          "base": "bba9f7e68856f827f535943983102164c6a79061",
          "tip": "8dee1a027b4b1a32731bffa71f620d3c31a10413",
          "finding_ids": [
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "apc1-coverage-r1",
          "performer": "claude:aap_apc1_coverage_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "456c681ccc9744bd468cd19e7fe78d75aaee79ec",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:aap_apc1_coverage_r1",
            "digest": "sha256:69eaf1da96ff42e3f884ccf567305f1c5eee3219f5f6245b7c255180cb46b60e",
            "excerpt": "Coverage axis, AP-C1, base bba9f7e6, tip 18f6c6dd; evidence current=true.\nC1 | medium | internal/anchors/locate.go:199-208 with registry_data.go:41-48 | ForbidCaseFoldedEmphasis strips only emphasis; the old check stripped every backtick and <code> tag, so four partial-span spellings now pass | ask-user | confidence 8\nC2 | low | internal/anchors/registry_data.go:36 | docs/greenfield-build-sequence.md keeps only Forbid rows, so a rename leaves them silently green | ask-user | confidence 7\nC3 | low | internal/anchors/registry_data.go:110 | the re-cut needle's diagnostic still names the create-or-retain clause | ask-user | confidence 7\nC4 | low | cmd/bench/anchor_help_test.go:99-116 | the sibling-close-up branch no longer runs with one Require row | no-op | confidence 6\nClean: AP06-AP08 each red on clause removal; wrapped needle matches; implement-spec-entry-validation left with its Require owner.\nTests: bench test --package ./internal/anchors pass, 165 tests; bench anchors .bench/BENCH.md 63 rows, no diagnostics.\nCount 4. Worst: C1.\n"
          },
          "axis": "Coverage",
          "base": "bba9f7e68856f827f535943983102164c6a79061",
          "tip": "8dee1a027b4b1a32731bffa71f620d3c31a10413",
          "finding_ids": [
            "C3"
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
      "from": "sha256:294bb59a470b025b0faa1f465f8561a826966e0879b0b69714709866831db182",
      "to": "sha256:9ae1c28cf65da60d226e873b3fb29a0ac2ef2a7eb9b940da5a78c04c69067181",
      "chunk_ids": {
        "AP-C2": [
          "AP-C2"
        ]
      }
    }
  ],
  "folds": [
    {
      "id": "F1",
      "after": "AP-C2",
      "from": "ebef308d813d87763b2482ada126ffe8103c648b",
      "commit": "bba9f7e68856f827f535943983102164c6a79061",
      "incoming": "d568cb05bc180bcd39ed95041edb265393a5fbb2",
      "source_digest": "6fb54f2d9c6369b4a1672a657335fb63c7596ada",
      "resolution": [],
      "resolution_digest": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "test_packages": [],
      "reviews": []
    }
  ]
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

## AP-C1 repair state

AP-C1: repair cycle 1 of 2 starts for ticket 2, with the targets S1 and C3. P2 merges with C3. No reviewer extension exists.

## AP-C1 Standards

AP-C1: 1 finding. The worst issue is S1.

- S1 (auto-fix, confidence 6): `internal/anchors/registry_data_test.go:352` names a cap-change rule that the prune removed. The repair renames the test to the rule that it holds.

## AP-C1 Spec

AP-C1: 1 actionable finding. The worst issue is P2.

- P1 (no-op, confidence 7): the ticket commit rewrote the AP06 to AP08 seam cells and added a planning-facts seam. The implement-spec phase directs the ticket commit to rewrite each covered planned seam cell into the citation form.
- P2 (auto-fix, confidence 6): `internal/anchors/registry_data.go:92` keeps a diagnostic that names a clause that the re-cut needle no longer holds. Merged with C3.

## AP-C1 Coverage

AP-C1: 1 actionable finding. The worst issue is C1.

- C1 (no-op, confidence 8): the generated case-folded rows do not catch every partial backtick or `<code>` spelling that the removed check stripped. The ticket asks for a row for each spelling that the old normalization matched at the ticket base. The spec says that this prune changes no matcher behavior. The repair records the narrowing in the AP47 census. Flagged for reviewer veto.
- C2 (no-op, confidence 7): `docs/greenfield-build-sequence.md` keeps only Forbid rows, so a rename leaves them silently green. Every Forbid row has this property, so it goes to `bench learning`. Flagged for reviewer veto.
- C3 (auto-fix, confidence 7): `internal/anchors/registry_data.go:110` keeps a diagnostic that names the create-or-retain clause. The repair rewords the diagnostic and the EXPECT of `implement-spec-worktree-before-preflight` together.
- C4 (no-op, confidence 6): `cmd/bench/anchor_help_test.go:99-116` keeps a branch that one Require row cannot reach; the unit tests in `locate_test.go` cover the line shift.
