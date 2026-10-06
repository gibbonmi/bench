# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/spec-stage-grader-trace/spec.md",
  "plan_digest": "sha256:1b0cc417b91dfb0abeb241102024cfed076c5e5d4ead195f5d96e11b1f206846",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GT-C1",
      "base": "f7ef3cee4ed28920a16168ea6e900ce7c6e71af7",
      "tip": "102144309646bffeba2b48b1168b4ce83341df76",
      "plan_digest": "sha256:1b0cc417b91dfb0abeb241102024cfed076c5e5d4ead195f5d96e11b1f206846",
      "source_digest": "034738cb6dce355d365eeb057f37f3bb1ebc7bca",
      "acceptance_rows": [
        "GT1",
        "GT2",
        "GT3",
        "GT4",
        "GT5",
        "GT6",
        "GT7",
        "GT8",
        "GT9",
        "GT10",
        "GT11",
        "GT12",
        "GT13",
        "GT14",
        "GT15",
        "GT16",
        "GT17",
        "GT18",
        "GT19"
      ],
      "verification": [
        {
          "id": "t1-anchors-v1",
          "performer": "claude:ft376_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "54595b43c8eaf671420d796061bee64f7c129076",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft376_t1",
            "digest": "sha256:8f80beeca4934e65fed60b0da85aeb5afcfb253015d37af475ebe38beab8d94e",
            "excerpt": "$ bench test --package ./internal/anchors\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/anchors,pass,1195\nfailures[0]{package,test,line}:\n"
          },
          "requirement": "t1-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t1-fixture-bites-v1",
          "performer": "claude:ft376_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "54595b43c8eaf671420d796061bee64f7c129076",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft376_t1",
            "digest": "sha256:003af642c4278fd2f518a21fff6fd52b1144b5ca2ceedbebfb8dbe6b84436584",
            "excerpt": "$ bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,19411\nfailures[0]{package,test,line}:\n$ bench probe internal/anchors/registry_spec_trace.go --swap (N1 full needle) --with '`Pin operators`' --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/anchors/registry_spec_trace.go,swap,failed,2,yes\n  TestEveryRetainedFixtureBitesThroughRegisteredOwner/map-discipline-pin-operator-trace did not bite through owner docs-currency-workflow\n"
          },
          "requirement": "t1-fixture-bites",
          "command": "bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner",
          "exit_code": 0,
          "probe": {
            "mutation": "Run the seven needle probes that the ticket's Acceptance names, one at a time: shorten that needle in specTraceAnchors as stated. TestEveryRetainedFixtureBitesThroughRegisteredOwner must fail on the named canary, TestSpecGraderTraceAnchors must stay green, and each restore must be exact.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft376_t1",
              "digest": "sha256:003af642c4278fd2f518a21fff6fd52b1144b5ca2ceedbebfb8dbe6b84436584",
              "excerpt": "$ bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,19411\nfailures[0]{package,test,line}:\n$ bench probe internal/anchors/registry_spec_trace.go --swap (N1 full needle) --with '`Pin operators`' --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/anchors/registry_spec_trace.go,swap,failed,2,yes\n  TestEveryRetainedFixtureBitesThroughRegisteredOwner/map-discipline-pin-operator-trace did not bite through owner docs-currency-workflow\n"
            }
          }
        },
        {
          "id": "t1-conformance-v1",
          "performer": "claude:ft376_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "54595b43c8eaf671420d796061bee64f7c129076",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft376_t1",
            "digest": "sha256:047e89c9d0ad19c48964f0aa2ea627518d05935068685b462ae21e7e8f4dca12",
            "excerpt": "$ bench test --package ./internal/conformance\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,56411\nfailures[0]{package,test,line}:\nskips[3]: environment capability skips (unix sockets, character device)\n"
          },
          "requirement": "t1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-canary-fixture-compliance-v1",
          "performer": "claude:ft376_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "54595b43c8eaf671420d796061bee64f7c129076",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft376_t1",
            "digest": "sha256:9c416a90716bfdaf9e1ebb3c9686e3376b771796e4152319573b5f9a8056f21f",
            "excerpt": "$ bench test --check canary-fixture-compliance\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4\nfailures[0]{package,test,line}:\n"
          },
          "requirement": "t1-canary-fixture-compliance",
          "command": "bench test --check canary-fixture-compliance",
          "exit_code": 0
        },
        {
          "id": "t1-docs-currency-workflow-v1",
          "performer": "claude:ft376_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "54595b43c8eaf671420d796061bee64f7c129076",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft376_t1",
            "digest": "sha256:9c927f41e95ead144d3f334c318189bc5bb6aadfa5ebb20e8201d738c4b81d14",
            "excerpt": "$ bench test --check docs-currency-workflow\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,1550\nfailures[0]{package,test,line}:\n"
          },
          "requirement": "t1-docs-currency-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t1-anchors-v3",
          "performer": "claude:ft376_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "034738cb6dce355d365eeb057f37f3bb1ebc7bca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft376_t1",
            "digest": "sha256:2d3a13df32cc58aa496628ea522c5bdd4663495cce5abc3ce466d71bf6dfed8d",
            "excerpt": "$ bench test --package ./internal/anchors\ntree: FT376-build,102144309646bffeba2b48b1168b4ce83341df76\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/anchors,pass,1125\nfailures[0]{package,test,line}:\nexit code 0\n"
          },
          "requirement": "t1-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "t1-fixture-bites-v3",
          "performer": "claude:ft376_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "034738cb6dce355d365eeb057f37f3bb1ebc7bca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft376_t1",
            "digest": "sha256:c95d26dd0e0e54ef30726e41e0e6cd3cd97e40739be8f3481232491972183900",
            "excerpt": "$ bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner\ntree: FT376-build,102144309646bffeba2b48b1168b4ce83341df76\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16148\nfailures[0]{package,test,line}:\nexit code 0\n$ bench probe internal/anchors/registry_spec_trace.go --swap (N1 full needle) --with '`Pin operators`' --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/anchors/registry_spec_trace.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/conformance,fail,16297\n  TestEveryRetainedFixtureBitesThroughRegisteredOwner/map-discipline-pin-operator-trace did not bite through owner docs-currency-workflow\nprobed test exit code 1 (observed under the same N1 mutation by bench test: Exit code 1)\n"
          },
          "requirement": "t1-fixture-bites",
          "command": "bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner",
          "exit_code": 0,
          "probe": {
            "mutation": "Run the seven needle probes that the ticket's Acceptance names, one at a time: shorten that needle in specTraceAnchors as stated. TestEveryRetainedFixtureBitesThroughRegisteredOwner must fail on the named canary, TestSpecGraderTraceAnchors must stay green, and each restore must be exact.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft376_t1",
              "digest": "sha256:c95d26dd0e0e54ef30726e41e0e6cd3cd97e40739be8f3481232491972183900",
              "excerpt": "$ bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner\ntree: FT376-build,102144309646bffeba2b48b1168b4ce83341df76\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16148\nfailures[0]{package,test,line}:\nexit code 0\n$ bench probe internal/anchors/registry_spec_trace.go --swap (N1 full needle) --with '`Pin operators`' --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/anchors/registry_spec_trace.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/conformance,fail,16297\n  TestEveryRetainedFixtureBitesThroughRegisteredOwner/map-discipline-pin-operator-trace did not bite through owner docs-currency-workflow\nprobed test exit code 1 (observed under the same N1 mutation by bench test: Exit code 1)\n"
            }
          }
        },
        {
          "id": "t1-conformance-v3",
          "performer": "claude:ft376_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "034738cb6dce355d365eeb057f37f3bb1ebc7bca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft376_t1",
            "digest": "sha256:b8236abec38b4941162c49bed9251d8dc0fb8268165407d54e2e76874c46ee86",
            "excerpt": "$ bench test --package ./internal/conformance\ntree: FT376-build,102144309646bffeba2b48b1168b4ce83341df76\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,41099\nfailures[0]{package,test,line}:\nskips[3]: environment capability skips (unix sockets, character device)\nexit code 0\n"
          },
          "requirement": "t1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-canary-fixture-compliance-v3",
          "performer": "claude:ft376_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "034738cb6dce355d365eeb057f37f3bb1ebc7bca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft376_t1",
            "digest": "sha256:514eb0f4be044a9b7b74cbf8f743fd0fdddb8e71e2681b8cbe39e9f57b58f1f6",
            "excerpt": "$ bench test --check canary-fixture-compliance\ntree: FT376-build,102144309646bffeba2b48b1168b4ce83341df76\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5\nfailures[0]{package,test,line}:\nexit code 0\n"
          },
          "requirement": "t1-canary-fixture-compliance",
          "command": "bench test --check canary-fixture-compliance",
          "exit_code": 0
        },
        {
          "id": "t1-docs-currency-workflow-v3",
          "performer": "claude:ft376_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "034738cb6dce355d365eeb057f37f3bb1ebc7bca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft376_t1",
            "digest": "sha256:ae37369aed739c1281b0a0feab6e811de11f3aa25948f8b6780c5dffc64985a9",
            "excerpt": "$ bench test --check docs-currency-workflow\ntree: FT376-build,102144309646bffeba2b48b1168b4ce83341df76\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,1459\nfailures[0]{package,test,line}:\nexit code 0\n"
          },
          "requirement": "t1-docs-currency-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
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
