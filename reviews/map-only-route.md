# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/map-only-route/spec.md",
  "plan_digest": "sha256:12a8f3f246e1dac6c12f0dec89ac955342f094a8545b16c69e6a4600e2d4fc26",
  "implementation_session": "",
  "chunks": [
    {
      "id": "MR-A",
      "base": "b10edf16c05da80525b0bb0b6ca8341c37a912f3",
      "tip": "e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a",
      "plan_digest": "sha256:12a8f3f246e1dac6c12f0dec89ac955342f094a8545b16c69e6a4600e2d4fc26",
      "source_digest": "40d64a917901a56a853a5b45d1213d5890abc577",
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
        },
        {
          "id": "mr-t01-rep1-guidance",
          "performer": "claude:mr_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
            "digest": "sha256:f954b594bd0e783839745ece25d4ee7f6713a434fcbb97522b6d0512e8eced57",
            "excerpt": "bench test --check docs-currency-workflow\ntree[1]{target,head,dirty}:\n  map-only-route-build,3327a8e356a2b5bb9f444b4e6dac75503c53b01c,true\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5011,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "guidance",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep1-prose",
          "performer": "claude:mr_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
            "digest": "sha256:ab896b680fae4b23a2e0e0bb78f28bc9988804c335ea41eb075084ef46afd0df",
            "excerpt": "bench test --check prose\ntree[1]{target,head,dirty}:\n  map-only-route-build,3327a8e356a2b5bb9f444b4e6dac75503c53b01c,true\ncheck[1]{name,kind,tests_run,subjects}:\n  prose,prose,0,536\n"
          },
          "requirement": "prose",
          "command": "bench test --check prose",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep1-budgets",
          "performer": "claude:mr_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
            "digest": "sha256:0b24779b59384f086999bfb05a949b46d5c9e6e6b807db075841d4bfa8ddfb7a",
            "excerpt": "bench test --check guidance-prose-budgets\ntree[1]{target,head,dirty}:\n  map-only-route-build,3327a8e356a2b5bb9f444b4e6dac75503c53b01c,true\ncheck[1]{name,kind,tests_run,subjects}:\n  guidance-prose-budgets,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,11,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep1-anchors-package",
          "performer": "claude:mr_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
            "digest": "sha256:48ff8cf366a9d4683c11b4bce36050e737fe728d72d2a8e7c615f126811fed7d",
            "excerpt": "bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  map-only-route-build,3327a8e356a2b5bb9f444b4e6dac75503c53b01c,true\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17519,335\n  github.com/gibbonmi/bench/internal/anchors,pass,2415,180\n  github.com/gibbonmi/bench/internal/conformance,pass,61132,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "anchors-package",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep1-command-package",
          "performer": "claude:mr_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
            "digest": "sha256:48ff8cf366a9d4683c11b4bce36050e737fe728d72d2a8e7c615f126811fed7d",
            "excerpt": "bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  map-only-route-build,3327a8e356a2b5bb9f444b4e6dac75503c53b01c,true\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17519,335\n  github.com/gibbonmi/bench/internal/anchors,pass,2415,180\n  github.com/gibbonmi/bench/internal/conformance,pass,61132,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep1-conformance-package",
          "performer": "claude:mr_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
            "digest": "sha256:48ff8cf366a9d4683c11b4bce36050e737fe728d72d2a8e7c615f126811fed7d",
            "excerpt": "bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  map-only-route-build,3327a8e356a2b5bb9f444b4e6dac75503c53b01c,true\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17519,335\n  github.com/gibbonmi/bench/internal/anchors,pass,2415,180\n  github.com/gibbonmi/bench/internal/conformance,pass,61132,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep1-root-conformance",
          "performer": "claude:mr_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
            "digest": "sha256:301a5a2b9723d9cac3a3b949800fcd9fe9acca9a8437e4ac005c7c1344a9889c",
            "excerpt": "bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  map-only-route-build,3327a8e356a2b5bb9f444b4e6dac75503c53b01c,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16063,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep1-central-property",
          "performer": "claude:mr_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
            "digest": "sha256:756b082f1e29bbe04e62056b93a0e2bd6a0b1e2971cc4bbdb3f9cd247d91d4db",
            "excerpt": "bench probe internal/anchors/registry_ticket_passes.go --omit '\t{Group: AfterImplementSpec, File: mapOnlyRoute, Kind: RequireInSection, Section: \"Stop and return\", Needle: \"The slicer stops when a map decision is not an exact predicate.\", Diagnostic: \"map-only route: Stop and return must stop the slicer when a map decision is not an exact predicate\"},' --package ./internal/anchors --run TestWriteSpecMapOnlyRoute --expect TestWriteSpecMapOnlyRoute\ntree[1]{target,head,dirty}:\n  map-only-route-build,3327a8e356a2b5bb9f444b4e6dac75503c53b01c,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/registry_ticket_passes.go,omit,failed,1,TestWriteSpecMapOnlyRoute,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,TestWriteSpecMapOnlyRoute,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,508\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,57,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestWriteSpecMapOnlyRoute,\"registry_ticket_passes_test.go:124: tree broken at \\\"The slicer stops when a map decision is not an exact predicate.\\\" in .agents/skills/bench-craft-spec/references/map-only-route.md = [], want \\\"map-only route: Stop and return must stop the s… (292 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
              "digest": "sha256:756b082f1e29bbe04e62056b93a0e2bd6a0b1e2971cc4bbdb3f9cd247d91d4db",
              "excerpt": "bench probe internal/anchors/registry_ticket_passes.go --omit '\t{Group: AfterImplementSpec, File: mapOnlyRoute, Kind: RequireInSection, Section: \"Stop and return\", Needle: \"The slicer stops when a map decision is not an exact predicate.\", Diagnostic: \"map-only route: Stop and return must stop the slicer when a map decision is not an exact predicate\"},' --package ./internal/anchors --run TestWriteSpecMapOnlyRoute --expect TestWriteSpecMapOnlyRoute\ntree[1]{target,head,dirty}:\n  map-only-route-build,3327a8e356a2b5bb9f444b4e6dac75503c53b01c,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/registry_ticket_passes.go,omit,failed,1,TestWriteSpecMapOnlyRoute,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,TestWriteSpecMapOnlyRoute,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,508\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,57,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestWriteSpecMapOnlyRoute,\"registry_ticket_passes_test.go:124: tree broken at \\\"The slicer stops when a map decision is not an exact predicate.\\\" in .agents/skills/bench-craft-spec/references/map-only-route.md = [], want \\\"map-only route: Stop and return must stop the s… (292 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "mr-t01-rep2-guidance",
          "performer": "claude:mr_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "40d64a917901a56a853a5b45d1213d5890abc577",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a",
            "digest": "sha256:53c2b87f8155e7078a2b19d97f5a84f1e866198d049c61647348028bc7f38cbd",
            "excerpt": "bench worktree exec map-only-route-build -- bench test --check docs-currency-workflow\ntree map-only-route-build,e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a\ncheck docs-currency-workflow,conformance,1,0\ngithub.com/gibbonmi/bench/internal/conformance,pass,4762,1\nfailures[0] skips[0]\n"
          },
          "requirement": "guidance",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep2-prose",
          "performer": "claude:mr_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "40d64a917901a56a853a5b45d1213d5890abc577",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a",
            "digest": "sha256:978611be147591f3757d45fb01c1256263d83c341b3d88aca1202d5871e7e074",
            "excerpt": "bench worktree exec map-only-route-build -- bench test --check prose\ntree map-only-route-build,e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a\ncheck prose,prose,0,536\n"
          },
          "requirement": "prose",
          "command": "bench test --check prose",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep2-budgets",
          "performer": "claude:mr_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "40d64a917901a56a853a5b45d1213d5890abc577",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a",
            "digest": "sha256:4745587ab7dd8785f85c944353c36e4a9b6867ea302fab1a1ac96abc71defc09",
            "excerpt": "bench worktree exec map-only-route-build -- bench test --check guidance-prose-budgets\ntree map-only-route-build,e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a\ncheck guidance-prose-budgets,conformance,1,0\ngithub.com/gibbonmi/bench/internal/conformance,pass,11,1\nfailures[0] skips[0]\n"
          },
          "requirement": "budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep2-anchors-package",
          "performer": "claude:mr_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "40d64a917901a56a853a5b45d1213d5890abc577",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a",
            "digest": "sha256:f9c642a94f8a991028659993b35f93b2642d324cbcf5bdc87846c46df1c9f0b4",
            "excerpt": "bench worktree exec map-only-route-build -- bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\ntree map-only-route-build,e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a\ngithub.com/gibbonmi/bench/cmd/bench,pass,16368,335\ngithub.com/gibbonmi/bench/internal/anchors,pass,2376,180\ngithub.com/gibbonmi/bench/internal/conformance,pass,57220,1690\nfailures[0]\nskips[1] internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,capability: privilege: cannot create a character device: operation not permitted\n"
          },
          "requirement": "anchors-package",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep2-command-package",
          "performer": "claude:mr_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "40d64a917901a56a853a5b45d1213d5890abc577",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a",
            "digest": "sha256:f9c642a94f8a991028659993b35f93b2642d324cbcf5bdc87846c46df1c9f0b4",
            "excerpt": "bench worktree exec map-only-route-build -- bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\ntree map-only-route-build,e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a\ngithub.com/gibbonmi/bench/cmd/bench,pass,16368,335\ngithub.com/gibbonmi/bench/internal/anchors,pass,2376,180\ngithub.com/gibbonmi/bench/internal/conformance,pass,57220,1690\nfailures[0]\nskips[1] internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,capability: privilege: cannot create a character device: operation not permitted\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep2-conformance-package",
          "performer": "claude:mr_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "40d64a917901a56a853a5b45d1213d5890abc577",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a",
            "digest": "sha256:f9c642a94f8a991028659993b35f93b2642d324cbcf5bdc87846c46df1c9f0b4",
            "excerpt": "bench worktree exec map-only-route-build -- bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\ntree map-only-route-build,e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a\ngithub.com/gibbonmi/bench/cmd/bench,pass,16368,335\ngithub.com/gibbonmi/bench/internal/anchors,pass,2376,180\ngithub.com/gibbonmi/bench/internal/conformance,pass,57220,1690\nfailures[0]\nskips[1] internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,capability: privilege: cannot create a character device: operation not permitted\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep2-root-conformance",
          "performer": "claude:mr_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "40d64a917901a56a853a5b45d1213d5890abc577",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a",
            "digest": "sha256:3d7ad0760c56a4294bce8dbf58618a372ed2e4abd9dfbe696a8d7fac66cae0df",
            "excerpt": "bench worktree exec map-only-route-build -- bench test --check root-conformance\ntree map-only-route-build,e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a\ncheck root-conformance,conformance,1,0\ngithub.com/gibbonmi/bench/internal/conformance,pass,15476,1\nfailures[0] skips[0]\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "mr-t01-rep2-central-property",
          "performer": "claude:mr_t01_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "40d64a917901a56a853a5b45d1213d5890abc577",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a",
            "digest": "sha256:98106b98ac27d8b4a4d7a5c8d3c90cc70424e25676cc3d74dfb3dcbbaf0b1e24",
            "excerpt": "bench worktree exec map-only-route-build -- bench probe internal/anchors/registry_ticket_passes.go --omit \"<MR20 anchor line: Needle 'The slicer stops when a map decision is not an exact predicate.'>\" --package ./internal/anchors --run TestWriteSpecMapOnlyRoute --expect TestWriteSpecMapOnlyRoute\ntree map-only-route-build,e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a\nprobe bit,internal/anchors/registry_ticket_passes.go,omit,failed,1,TestWriteSpecMapOnlyRoute,yes\ngithub.com/gibbonmi/bench/internal/anchors,fail,54,1 (mutated)\nregistry_ticket_passes_test.go:124: tree broken at \"The slicer stops when a map decision is not an exact predicate.\" in .agents/skills/bench-craft-spec/references/map-only-route.md = [], want \"map-only route: Stop and return must stop the s... (292 bytes)\"\n"
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
              "ref": "e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a",
              "digest": "sha256:98106b98ac27d8b4a4d7a5c8d3c90cc70424e25676cc3d74dfb3dcbbaf0b1e24",
              "excerpt": "bench worktree exec map-only-route-build -- bench probe internal/anchors/registry_ticket_passes.go --omit \"<MR20 anchor line: Needle 'The slicer stops when a map decision is not an exact predicate.'>\" --package ./internal/anchors --run TestWriteSpecMapOnlyRoute --expect TestWriteSpecMapOnlyRoute\ntree map-only-route-build,e971a5e76146d2f0db342844ff8a0bb6e0b4fb2a\nprobe bit,internal/anchors/registry_ticket_passes.go,omit,failed,1,TestWriteSpecMapOnlyRoute,yes\ngithub.com/gibbonmi/bench/internal/anchors,fail,54,1 (mutated)\nregistry_ticket_passes_test.go:124: tree broken at \"The slicer stops when a map decision is not an exact predicate.\" in .agents/skills/bench-craft-spec/references/map-only-route.md = [], want \"map-only route: Stop and return must stop the s... (292 bytes)\"\n"
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
        },
        {
          "id": "mr-a-standards-r2",
          "performer": "claude:mr_a_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:mr_a_standards_r2",
            "digest": "sha256:e9f131eb84f15ed20021d0f7eb36af1c666cb479da8c2e5ccfd1dad3a84ce01a",
            "excerpt": "## Standards (claude:mr_a_standards_r2, opus/high, confirming round 2, pair b10edf16c..306f1eb9, chunk tip 3327a8e3, evidence sha256:b9db3f727f1778b7625cd781458f2bbd858127493840f473f5db80c5791d4483 current=true)\n\n- S1 closed: map-only-route.md:42-43 lists only the check-read parts, and line 43 says no check reads the gap log (true: gather.go:48-70 reads status, fence, coverage map). Anchor registry_ticket_passes.go:91 and test literal registry_ticket_passes_test.go:121 changed together; no canary pins the old bytes.\n- S2 closed: field-guide.html:1086-1089 and 1329-1331 put the map-only route before the fan-out; sentences 17/18 and 12/21 words; <p> markup unchanged.\n- S4 closed: benchkit.md:590-591 is its own one-sentence paragraph; the first paragraph (581-588) has 6 sentences.\n- S5 closed: const route at test line 100; the test does not use the registry's mapOnlyRoute; all 20 reference rows use route; the pointer row names its own file.\n- S6 closed: the doc comment at test lines 96-98 is active voice.\n- S3/S7: spec line 121 and the Closed-decisions lines untouched.\n- S8 (new, one source per fact): map-only-route.md:37 \"The slicer also stops when the map holds no decision ticket.\" internal/maps/validation.go:21-22 already gives the diagnostic \"missing decision ticket\" for that map, and Eligibility (line 7) refuses any map with a bench maps diagnostic, so the new stop restates the rule and can never fire. MR34's reason at spec.md:228 (\"no check turns red\") is false. Disposition: auto-fix. Confidence 8.\n- S9 (new, drift): D5 at spec.md:297 names two stop conditions while spec.md:131-132 has three. Disposition: auto-fix. Confidence 7.\n\nOpen: 2. Worst: S8.\n\n### Advice\n- Fix S8 by dropping MR34 with its anchor and test row, or by pointing the edge-inventory line (spec.md:239) at Eligibility; if MR34 stays, its reason must cite the bench maps diagnostic.\n- The new ticket acceptance item says \"The Thin spec sentence\" while the reference has two sentences. MR23 claims its test covers \"names the gap log apart\", but no anchor pins line 43.\n\n### Examined\n--check-current inside the worktree (the primary checkout refused with \"assignment required\"); the full diff 0fc73c29..306f1eb9; map-only-route.md; spec.md; field-guide.html 1078-1091 and 1324-1335; benchkit.md 575-599; gather.go 40-79; validation.go 1-45; ste-prose.md; finding-discipline.md; searches over Further notes and gap-log readers, canary pins, and \"decision ticket\".\n"
          },
          "axis": "Standards",
          "base": "b10edf16c05da80525b0bb0b6ca8341c37a912f3",
          "tip": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
          "finding_ids": [
            "S8",
            "S9"
          ],
          "supersedes": [
            "mr-a-standards-r1"
          ]
        },
        {
          "id": "mr-a-spec-r2",
          "performer": "claude:mr_a_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:mr_a_spec_r2",
            "digest": "sha256:db236fa207bdf27c40479b740e72bcb83b9efe7085e41c2c6af9e921306e8087",
            "excerpt": "## Spec (claude:mr_a_spec_r2, opus/high, confirming round 2, pair b10edf16c..306f1eb9, chunk tip 3327a8e3, evidence sha256:b9db3f727f1778b7625cd781458f2bbd858127493840f473f5db80c5791d4483 current=true)\n\n- P1 closed: story 31 \"both write-spec passages name the route, so that the guide stays true\"; field-guide.html:1086-1089 and 1329-1331 now put the map-only route first and the fan-out after \"Otherwise\".\n- P2 closed: the reference limits the read list to six parts at map-only-route.md:42, and line 43 adds the gap log for the sign-off with no check reading it; spec.md:137 says the same.\n- Expansion accepted: the edge inventory promised the stop under story 22 (spec.md:239); the new stop is one sentence in spec.md:132, map-only-route.md:37, the registry, and the test; MR34 seam cell (spec.md:228) in the citation form; MR28, the chunk row, ticket Covers, anchor list, and acceptance extended; learning at capture/learnings.md:27.\n- P3 (new): MR23 now says \"it names the gap log apart from them\" with TestWriteSpecMapOnlyRoute as its seam, and MR28 needs one anchor per MR23 route sentence; map-only-route.md:43 has no anchor and no test row, so deleting it stays green. Disposition: auto-fix. Confidence 6.\n- P4 (new): D5 at spec.md:297 still names two stop triggers; the veto list leaves out the empty-map stop. Disposition: auto-fix. Confidence 5.\n\nRow audit: MR23 partial (P3), MR31 met, MR34 met.\nScope: the diff touches only the reference, the registry, the test, the field guide, benchkit.md, the spec, ticket 1, and reviews/map-only-route.md. The benchkit.md:590-593 paragraph split traces to no Spec finding (it is Standards S4); inside the fence, no behavior change.\nOpen: 2. Worst: P3.\n\n### Advice\nAdd a RequireInSection anchor for the gap-log sentence with its own independent test literal. Add the empty-map stop to the D5 veto line.\n\n### Examined\nThe diff 0fc73c29..306f1eb9; spec.md 129-147, 193-245, 287-301; ticket 1; map-only-route.md in full; field-guide.html 1076-1091 and 1322-1333; benchkit.md 580-593; capture/learnings.md:27-30; bench coverage (34 rows); rg for gap-log readers and anchors.\n"
          },
          "axis": "Spec",
          "base": "b10edf16c05da80525b0bb0b6ca8341c37a912f3",
          "tip": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
          "finding_ids": [
            "P3",
            "P4"
          ],
          "supersedes": [
            "mr-a-spec-r1"
          ]
        },
        {
          "id": "mr-a-coverage-r2",
          "performer": "claude:mr_a_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "26450ed87e87af8771862cc4cf253ee1669672ab",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:mr_a_coverage_r2",
            "digest": "sha256:8973cb26746edea0c2eaa0018a42330507bfcf947c49e752dbb389d05b637308",
            "excerpt": "## Coverage (claude:mr_a_coverage_r2, opus/high, confirming round 2, pair b10edf16c..306f1eb9, chunk tip 3327a8e3, evidence sha256:b9db3f727f1778b7625cd781458f2bbd858127493840f473f5db80c5791d4483 current=true)\n\n- C1 closed: the third stop sentence at map-only-route.md:37 under ## Stop and return, one registry anchor at registry_ticket_passes.go:92, one independent test row at registry_ticket_passes_test.go:119. Registry lines 74-95 hold 22 anchors (4 Eligibility, 12 Route steps, 4 Stop and return, 1 Thin spec, 1 pointer); the test holds 22 rows; one-to-one with MR1-4, MR6-8, MR10-18, MR20-23, MR25, MR34. The S5 const route holds only the file path; every needle and want stays an inline literal; the test reads no registry symbol.\n- MR23 split confirmed: the anchored first sentence is one physical line at reference line 42 under ## Thin spec; registry needle (line 94) and test literal match byte for byte.\n- C2 (new): the second Thin spec sentence at reference line 43 has no anchor. Input: delete line 43. Expected break: the MR23 clause \"names the gap log apart from them\" disappears and nothing turns red (probe 3 silent); MR28 partly unmet. Should exist: a second Thin spec RequireInSection anchor and a matching test row, or a narrowed MR23 behavior cell. Disposition: ask-user. Confidence 6.\n\n| command | verdict | first_failure | restored |\n|---|---|---|---|\n| swap the MR34 registry Diagnostic, --package ./internal/anchors --run TestWriteSpecMapOnlyRoute | bit | TestWriteSpecMapOnlyRoute | yes |\n| swap the MR23 sentence to drop \"ownership fences\", --check docs-currency-workflow | bit | TestRootConformance | yes |\n| omit reference line 43, --check docs-currency-workflow | silent | null | yes |\n\nbench coverage --check: valid, 34 rows. docs-currency-workflow: pass.\nOpen: 1. Worst: C2.\n\n### Advice\nThe MR34 story cell reads 22, the same story as MR22; the Spec axis should confirm that story 22 owns the empty-map stop.\n\n### Examined\nThe repair diff 0fc73c29..306f1eb9; map-only-route.md 1-53; registry_ticket_passes.go 60-96; anchor_harness_test.go 155-224; rg over internal/anchors for registry symbols; the two skill files. Final git status --short: empty.\n"
          },
          "axis": "Coverage",
          "base": "b10edf16c05da80525b0bb0b6ca8341c37a912f3",
          "tip": "3327a8e356a2b5bb9f444b4e6dac75503c53b01c",
          "finding_ids": [
            "C2"
          ],
          "supersedes": [
            "mr-a-coverage-r1"
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
    },
    {
      "from": "sha256:ca1148ab368111df094e860c4d985c2a720ff0921efcfeec3f0fb926a3c8f927",
      "to": "sha256:9a4314eef28fe52d8f9e8b7a0d9b3516e1078a10b294009e817c0c9731fe221e",
      "chunk_ids": {
        "MR-A": [
          "MR-A"
        ]
      }
    },
    {
      "from": "sha256:9a4314eef28fe52d8f9e8b7a0d9b3516e1078a10b294009e817c0c9731fe221e",
      "to": "sha256:12a8f3f246e1dac6c12f0dec89ac955342f094a8545b16c69e6a4600e2d4fc26",
      "chunk_ids": {
        "MR-A": [
          "MR-A"
        ]
      }
    }
  ]
}
```
