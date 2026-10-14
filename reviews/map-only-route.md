# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/map-only-route/spec.md",
  "plan_digest": "sha256:c1c0aa81ad757b6ca289c99ba2d3700ba0fb76f4d3b9e49c12df4aef0a0db983",
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
