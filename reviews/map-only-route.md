# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/map-only-route/spec.md",
  "plan_digest": "sha256:ca1148ab368111df094e860c4d985c2a720ff0921efcfeec3f0fb926a3c8f927",
  "implementation_session": "",
  "chunks": [
    {
      "id": "MR-A",
      "base": "b10edf16c05da80525b0bb0b6ca8341c37a912f3",
      "tip": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
      "plan_digest": "sha256:c1c0aa81ad757b6ca289c99ba2d3700ba0fb76f4d3b9e49c12df4aef0a0db983",
      "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
      "acceptance_rows": [
        "MR1",
        "MR2",
        "MR3",
        "MR4",
        "MR5",
        "MR6",
        "MR7",
        "MR8",
        "MR9",
        "MR10",
        "MR11",
        "MR12",
        "MR13",
        "MR14",
        "MR15",
        "MR16",
        "MR17",
        "MR18",
        "MR19",
        "MR20",
        "MR21",
        "MR22",
        "MR23",
        "MR24",
        "MR25",
        "MR26",
        "MR27",
        "MR28",
        "MR29",
        "MR30",
        "MR31",
        "MR32",
        "MR33"
      ],
      "verification": [
        {
          "id": "mr-t01-guidance",
          "performer": "agy:mr_t01",
          "role": "author-verification",
          "model": "gemini-4-argon-high",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
            "digest": "sha256:3019b64e33b774c5e58c9d8631876af5185964c9973d9faad7d488db83f5bede",
            "excerpt": "bench test --check docs-currency-workflow\ntree[1]{target,head,dirty}:\n  map-only-route-build,6d0989b8dc4a608a1b8b5721194f03d5846af07b,false\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5727,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "guidance",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "mr-t01-prose",
          "performer": "agy:mr_t01",
          "role": "author-verification",
          "model": "gemini-4-argon-high",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
            "digest": "sha256:e98196db8989c6590d8ddae1acc69dff82514b6178734b7272b3aa7436f7ed05",
            "excerpt": "bench test --check prose\ntree[1]{target,head,dirty}:\n  map-only-route-build,6d0989b8dc4a608a1b8b5721194f03d5846af07b,false\ncheck[1]{name,kind,tests_run,subjects}:\n  prose,prose,0,535\n"
          },
          "requirement": "prose",
          "command": "bench test --check prose",
          "exit_code": 0
        },
        {
          "id": "mr-t01-budgets",
          "performer": "agy:mr_t01",
          "role": "author-verification",
          "model": "gemini-4-argon-high",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
            "digest": "sha256:d8b44a4b6e98816874eaed49e1641b7df6a1bfa4c6724a139907b4d6d8a47907",
            "excerpt": "bench test --check guidance-prose-budgets\ntree[1]{target,head,dirty}:\n  map-only-route-build,6d0989b8dc4a608a1b8b5721194f03d5846af07b,false\ncheck[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "mr-t01-anchors-package",
          "performer": "agy:mr_t01",
          "role": "author-verification",
          "model": "gemini-4-argon-high",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
            "digest": "sha256:1df57300578f6e06adc5e48ecbdf7ed1f8f57adc84fd9623fc01446931a8c25c",
            "excerpt": "bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  map-only-route-build,6d0989b8dc4a608a1b8b5721194f03d5846af07b,false\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17850,335\n  github.com/gibbonmi/bench/internal/anchors,pass,2550,180\n  github.com/gibbonmi/bench/internal/conformance,pass,67516,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "anchors-package",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "mr-t01-command-package",
          "performer": "agy:mr_t01",
          "role": "author-verification",
          "model": "gemini-4-argon-high",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
            "digest": "sha256:1df57300578f6e06adc5e48ecbdf7ed1f8f57adc84fd9623fc01446931a8c25c",
            "excerpt": "bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  map-only-route-build,6d0989b8dc4a608a1b8b5721194f03d5846af07b,false\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17850,335\n  github.com/gibbonmi/bench/internal/anchors,pass,2550,180\n  github.com/gibbonmi/bench/internal/conformance,pass,67516,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "mr-t01-conformance-package",
          "performer": "agy:mr_t01",
          "role": "author-verification",
          "model": "gemini-4-argon-high",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
            "digest": "sha256:1df57300578f6e06adc5e48ecbdf7ed1f8f57adc84fd9623fc01446931a8c25c",
            "excerpt": "bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  map-only-route-build,6d0989b8dc4a608a1b8b5721194f03d5846af07b,false\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17850,335\n  github.com/gibbonmi/bench/internal/anchors,pass,2550,180\n  github.com/gibbonmi/bench/internal/conformance,pass,67516,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "mr-t01-root-conformance",
          "performer": "agy:mr_t01",
          "role": "author-verification",
          "model": "gemini-4-argon-high",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
            "digest": "sha256:4d1f083c381974064f0d1650aa88459062ec2a55f95473835bd0a13930f63962",
            "excerpt": "bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  map-only-route-build,6d0989b8dc4a608a1b8b5721194f03d5846af07b,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16781,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "mr-t01-central-property",
          "performer": "agy:mr_t01",
          "role": "author-verification",
          "model": "gemini-4-argon-high",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
            "digest": "sha256:18ecf290c2706456ffafcc9385f2af2edbb6dc0d7788f7fca41619629db81172",
            "excerpt": "bench probe internal/anchors/registry_ticket_passes.go --omit '\t{Group: AfterImplementSpec, File: mapOnlyRoute, Kind: RequireInSection, Section: \"Stop and return\", Needle: \"The slicer stops when a map decision is not an exact predicate.\", Diagnostic: \"map-only route: Stop and return must stop the slicer when a map decision is not an exact predicate\"},' --package ./internal/anchors --run TestWriteSpecMapOnlyRoute --expect TestWriteSpecMapOnlyRoute\ntree[1]{target,head,dirty}:\n  map-only-route-build,6d0989b8dc4a608a1b8b5721194f03d5846af07b,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/registry_ticket_passes.go,omit,failed,1,TestWriteSpecMapOnlyRoute,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,TestWriteSpecMapOnlyRoute,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,508\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,59,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestWriteSpecMapOnlyRoute,\"registry_ticket_passes_test.go:122: tree broken at \\\"The slicer stops when a map decision is not an exact predicate.\\\" in .agents/skills/bench-craft-spec/references/map-only-route.md = [], want \\\"map-only route: Stop and return must stop the s… (292 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0,
          "probe": {
            "mutation": "Delete the stop-and-return anchor from internal/anchors/registry_ticket_passes.go. TestWriteSpecMapOnlyRoute must turn red.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
              "digest": "sha256:18ecf290c2706456ffafcc9385f2af2edbb6dc0d7788f7fca41619629db81172",
              "excerpt": "bench probe internal/anchors/registry_ticket_passes.go --omit '\t{Group: AfterImplementSpec, File: mapOnlyRoute, Kind: RequireInSection, Section: \"Stop and return\", Needle: \"The slicer stops when a map decision is not an exact predicate.\", Diagnostic: \"map-only route: Stop and return must stop the slicer when a map decision is not an exact predicate\"},' --package ./internal/anchors --run TestWriteSpecMapOnlyRoute --expect TestWriteSpecMapOnlyRoute\ntree[1]{target,head,dirty}:\n  map-only-route-build,6d0989b8dc4a608a1b8b5721194f03d5846af07b,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/registry_ticket_passes.go,omit,failed,1,TestWriteSpecMapOnlyRoute,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,TestWriteSpecMapOnlyRoute,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,508\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,59,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestWriteSpecMapOnlyRoute,\"registry_ticket_passes_test.go:122: tree broken at \\\"The slicer stops when a map decision is not an exact predicate.\\\" in .agents/skills/bench-craft-spec/references/map-only-route.md = [], want \\\"map-only route: Stop and return must stop the s… (292 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "mr-a-standards-r1",
          "performer": "claude:mr_a_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:mr_a_standards",
            "digest": "sha256:4cd8d1ecff64d23501847becb7422f9ca8706287d7041e990ba6efe76448c82f",
            "excerpt": "## Standards (claude:mr_a_standards, opus/high, round 1, pair b10edf16c..fe033446d, evidence sha256:b96e00c889bfe1fe52f9c6e097ba20dbc22a2cdb70943ddbfa6ef0d13257a89f current=true)\n\n- S1. Rule: BENCH.md invariant 3. map-only-route.md:60 says the coverage check and the build preflight read the ownership fences and the gap log; no non-test Go file reads `Further notes` (internal/preflight/gather.go:48-70 reads status, fence, coverage map). Anchored sentence: fix moves registry_ticket_passes.go:93 and registry_ticket_passes_test.go:120. Spec line 136 has the same wording. Disposition: ask-user. Confidence 7.\n- S2. Rule: invariant 3 and STE 25-word bound. docs/field-guide.html:1086-1089 and :1329-1331 put the research fan-out before the map-only route, against reference line 34 (\"no research fan-out\"); the next sentence gives \"That fork\" the research findings; sentences of 29 and 32 words. Not anchored. Disposition: auto-fix. Confidence 8.\n- S3. Rule: AGENTS.md one source per fact. map-only-route.md:24 \"The proposal tolerates only the closure-check reds and the `fence-writes` red.\" repeats proposalToleratedChecks at internal/preflight/proposal.go:14 with no anchor tying them. Spec line 121 requires the sentence. Disposition: ask-user. Confidence 6.\n- S4. Rule: STE one topic per paragraph. projects/benchkit.md:590-592: the new map-only sentence opens the escalation-and-build paragraph, so \"Top + high remains a reviewer-approved escalation before the fork\" seems to cover a route with no fork. A separate paragraph fixes both. Disposition: auto-fix. Confidence 5.\n- S5. Smell: Duplicated Code. registry_ticket_passes_test.go:101-121 repeats the reference path literal 20 times; TestWriteSpecChunkedSlicing names each file once in a local const (:54-57). Disposition: auto-fix. Confidence 5.\n- S6. Rule: STE active voice via craft-comments. registry_ticket_passes_test.go:96-98 \"is written here apart from the registry\" is passive. Disposition: auto-fix. Confidence 4.\n- S7. Rule: one source per fact. map-only-route.md:44-45 repeats the ticket-review cap and the fence/plan-only/approval order from bench-write-spec.md:77-81; spec lines 122-123 require both. Disposition: no-op. Confidence 4.\n\nbench-implement-spec.md did not cause S1-S7; no command change is necessary.\nCount: 7 (6 held, 1 no-op). Worst: S2. Repair targets (de-duplicated): 5.\n\n### Advice\n- Name drift: \"the slicing research\" (pointer), \"paired slicing readers\" (reference), \"two read-only readers for each chunk\" (full route); \"the fork\" vs \"spec fork\".\n- \"thin spec\" is new and undefined in CONTEXT.md; used in step 3 before the Thin spec section.\n- Route steps 1-2 imperative, 3-10 descriptive.\n- map-only-route.md:66 restates the interim fence from step 6.\n- CHANGELOG.md names \"anchors in internal/anchors\" (implementation detail); field-guide.html:1089 wraps short.\n\n### Examined\nDelta b10edf16..fe033446; craft-review SKILL.md and finding-discipline.md; AGENTS.md; BENCH.md; ste-prose.md; bench-craft-comments/SKILL.md; bench-craft-skills/SKILL.md:45-74; the spec and ticket 1; bench-write-spec.md; projects/benchkit.md:570-604; field-guide.html:1078-1099 and :1324-1337; registry_ticket_passes.go:1-75; registry_ticket_passes_test.go:1-95; gather.go:40-79; retained_workflow_test.go:252-265; rg sweeps for Further notes, fence-writes, CONTEXT.md vocabulary, craft-line and slicing checks, field-guide anchors. Not read: the rest of CONTEXT.md and the remaining evidence pages. No probes.\n"
          },
          "axis": "Standards",
          "base": "b10edf16c05da80525b0bb0b6ca8341c37a912f3",
          "tip": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "S5",
            "S6",
            "S7"
          ],
          "supersedes": []
        },
        {
          "id": "mr-a-spec-r1",
          "performer": "claude:mr_a_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:mr_a_spec",
            "digest": "sha256:ef9e69333b3db82c17715cc79ce255fe45261757915a4698839dc1b781a1e074",
            "excerpt": "## Spec (claude:mr_a_spec, opus/high, round 1, pair b10edf16c..fe033446d, evidence sha256:b96e00c889bfe1fe52f9c6e097ba20dbc22a2cdb70943ddbfa6ef0d13257a89f current=true)\n\n- P1. Story 7 (\"omit the research, the spec fork, the spec review, and the slicing readers\") and story 31 (\"so that the guide stays true\"). docs/field-guide.html:1086-1089 and :1329-1331 attach the map-only route after \"fans out read-only research delegates, and then ...\"; read as written, the fan-out still runs before the route. No anchor or canary pins either passage. Disposition: auto-fix (put the route ahead of the fan-out as the other branch). Confidence 7. No command change necessary.\n- P2. D6: the thin spec \"holds only the parts that the coverage check and the build preflight read, plus the gap log.\" map-only-route.md:60 says \"The coverage check and the build preflight read ... ownership fences, and gap log.\" No production code reads Further notes (rg over internal and cmd excluding tests finds only anchor needles; gather.go:48-70 reads status, fence, coverage map). Spec line 136 has the same wording: a contradiction inside the spec. Anchored: a reword moves the registry needle (Thin spec entry) and the TestWriteSpecMapOnlyRoute literal. Disposition: ask-user. Confidence 6. No command change necessary.\n\nRow audit: met MR1-MR30, MR32, MR33; partial MR31 (P1); missing none; MR23 met as worded but wrong against D6 (P2).\nCount: 2. Worst: P1. Repair targets (de-duplicated): 2.\n\nVerified with no finding: ten route steps in spec order (reference 32-48); step 1 sends to #entry-contract and line 50 to #exit-handoff, no copy, links resolve; 21 anchors and 21 test literals in AfterImplementSpec; pointer verbatim on line 49, file at 87 lines; README.md and ADR 0022 untouched; the Spec authoring phrase at benchkit.md:581 unchanged; 22 seam cells in the `path` (`TestName`) form and bench coverage parses them; bench maps, --propose-writes, --plan-only exist; version 2 plan amendment in scope; no scope creep.\nAuthor judgment calls: benchkit.md placement accepted (six-sentence bound not verified against internal/prose); MR20-MR22 as three sentences accepted.\n\n### Advice\n- benchkit.md:591-592 \"Top + high remains a reviewer-approved escalation before the fork\" follows the map-only sentence; \"before the full route's fork\" removes the ambiguity.\n- field-guide.html:1089 short wrapped line; reflow with P1.\n- The writeSpecCommand comment at registry_ticket_passes.go:144-146 is incidental but harmless.\n\n### Examined\n--check-current only, no evidence page; one full git diff plus --stat; spec.md whole; ticket 1; bench-write-spec.md whole; benchkit.md:565-604; field-guide.html:1080-1099 and 1324-1337; preflight/gather.go:40-79; bench coverage (33 rows); bench preflight build --help and bench maps --help; rg over seam cell form, Further notes readers, field-guide and fork anchors; craft-review SKILL.md and finding-discipline.md. No test or probe.\n"
          },
          "axis": "Spec",
          "base": "b10edf16c05da80525b0bb0b6ca8341c37a912f3",
          "tip": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "mr-a-coverage-r1",
          "performer": "claude:mr_a_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f9cec8c8aab62669336c04997d913ec15c9ca761",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:mr_a_coverage",
            "digest": "sha256:9a44fb6d66aa9d6717f414d9f0e4492c8991a7bf941ae3bb4691c7c0cf6a8c48",
            "excerpt": "## Coverage (claude:mr_a_coverage, opus/high, round 1, pair b10edf16c..fe033446d, evidence sha256:b96e00c889bfe1fe52f9c6e097ba20dbc22a2cdb70943ddbfa6ef0d13257a89f current=true)\n\n- C1. Input: a ready compiled map with no decision ticket. Expected break: spec.md:237 (edge inventory) says the slicer \"stops and returns to the reviewer under story 22\", but the reference's Stop and return (map-only-route.md:35-36) and the spec's stop decision (spec.md:131) give only two triggers (non-predicate decision, open code seam). A slicer on an empty map finds no stop, and \"at least one story for each map decision\" (reference 44) is true with zero stories. Should exist: a third stop sentence with an anchor and an independent TestWriteSpecMapOnlyRoute row, or a recorded out-of-scope decision. Uncertain: internal/maps/validation.go:99-123 raises no diagnostic for zero tickets; parser not read. Disposition: ask-user. Confidence 5. No command change necessary.\n\n| probe | verdict | first_failure | restored |\n|---|---|---|---|\n| reference --insert-before \"The slicer also stops when the map leaves a code seam open.\" --with \"## Stop notes\\n\\n\" (moves MR21 and MR22 needles out of their H2), --check docs-currency-workflow | bit | TestRootConformance (Stop and return must stop the slicer when the map leaves a code seam open) | yes |\n| registry --swap of the MR23 diagnostic to \"...list the parts that the checks read\", test unchanged, --package ./internal/anchors --run TestWriteSpecMapOnlyRoute | bit | TestWriteSpecMapOnlyRoute (registry_ticket_passes_test.go:122) | yes |\n| reference --omit of the step 9 line \"These rows are not late clarification questions, so the two-question limit does not apply to them.\", --check docs-currency-workflow | silent | null | yes |\n\nbench coverage --check specs/map-only-route/spec.md: valid, 33 rows (11 uncited rows cite review-owned or gate evidence). bench test --check docs-currency-workflow: pass.\nCount: 1. Worst: C1. Repair targets (de-duplicated): 1.\n\n### Advice\n- MR17's why clause names the two-question-limit failure; the step 9 exemption sentence has no anchor (third probe silent). MR17's approved behavior is anchored (registry_ticket_passes.go:88), so this is hardening.\n- No row covers the step 6 sentence \"The proposal tolerates only the closure-check reds and the `fence-writes` red.\" (reference 24).\n- Removing a test row does not red the gate: the harness (anchor_harness_test.go:161-173) never checks that each registry entry has a test row; every anchor family works this way. Not probed.\n\n### Examined\nEvidence current=true; full base..tip diff (spill 1791936625585150614-cb40aaa388e386ea.out); anchors and test rows enumerated independently: 21 rows in scope, 21 registry entries (registry_ticket_passes.go:74-94), 21 test rows (registry_ticket_passes_test.go:101-121) with plain-string needles and diagnostics (MR29); needle placement: Eligibility lines 5, 6, 7, 9; Route steps 14-29; Stop and return 35-37; Thin spec 41; pointer at bench-write-spec.md:49 under ## Who runs this phase (line 40), coordinator probe already showed it bites; all 5 edges and both won't-handle lines walked; the harness, anchor_harness_diagnostics_test.go, maps/validation.go readiness checks, craft-review SKILL.md, finding-discipline.md, bounded repair policy. Final git status --short: empty.\n"
          },
          "axis": "Coverage",
          "base": "b10edf16c05da80525b0bb0b6ca8341c37a912f3",
          "tip": "6d0989b8dc4a608a1b8b5721194f03d5846af07b",
          "finding_ids": [
            "C1"
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
      "from": "sha256:c1c0aa81ad757b6ca289c99ba2d3700ba0fb76f4d3b9e49c12df4aef0a0db983",
      "to": "sha256:5f596423575e8a3108be9ac5d57b0140ee4d6f78d5fecdb55a0b2b048e8bc03b",
      "chunk_ids": {
        "MR-A": [
          "MR-A"
        ]
      }
    },
    {
      "from": "sha256:5f596423575e8a3108be9ac5d57b0140ee4d6f78d5fecdb55a0b2b048e8bc03b",
      "to": "sha256:ca1148ab368111df094e860c4d985c2a720ff0921efcfeec3f0fb926a3c8f927",
      "chunk_ids": {
        "MR-A": [
          "MR-A"
        ]
      }
    }
  ]
}
```
