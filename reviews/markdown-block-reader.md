# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/markdown-block-reader/spec.md",
  "plan_digest": "sha256:7c5e675de050f405240c89a679c3214d28f457e88d0a29c7934ebc3cdcedd165",
  "implementation_session": "",
  "chunks": [
    {
      "id": "MB-C1",
      "base": "36a07d83591f99cc079f571f185d25e606db0a2c",
      "tip": "252408eb9bce41230a88ccc4ca7a511d4415b848",
      "plan_digest": "sha256:7c5e675de050f405240c89a679c3214d28f457e88d0a29c7934ebc3cdcedd165",
      "source_digest": "780d3f3f9145f6dc8894c30ca4169925915a6fa6",
      "acceptance_rows": [
        "MB1",
        "MB2",
        "MB3",
        "MB4",
        "MB5",
        "MB6",
        "MB7",
        "MB8",
        "MB9",
        "MB10",
        "MB11",
        "MB12",
        "MB13",
        "MB14",
        "MB15",
        "MB16",
        "MB17",
        "MB18",
        "MB19",
        "MB20",
        "MB21",
        "MB22",
        "MB23",
        "MB71",
        "MB72",
        "MB73",
        "MB74",
        "MB75",
        "MB76",
        "MB77",
        "MB78",
        "MB79",
        "MB88",
        "MB89",
        "MB90",
        "MB92"
      ],
      "verification": [
        {
          "id": "t1-markdown",
          "performer": "ft358-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "780d3f3f9145f6dc8894c30ca4169925915a6fa6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-author-01",
            "digest": "sha256:a5c5cb4e8924755c278a6cb419b8a6ce9b7bf1b4ca0c76304fdce998f10b455b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/markdown/markdown.go,swap,failed,1,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/markdown,^TestReadFenceEdges$,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,356\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/markdown,fail,2,10\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/markdown,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,\"markdown_test.go:68: reading of \\\"````\\\\\\\\n```\\\\\\\\nx\\\\\\\\n````\\\":\",11\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "252408eb9bce41230a88ccc4ca7a511d4415b848",
          "requirement": "markdown",
          "command": "bench test --package ./internal/markdown --run '^TestReadFenceEdges$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/markdown: make the fence closer accept any run of three or more of the opener's character; expect TestReadFenceEdges",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-author-01",
              "digest": "sha256:a5c5cb4e8924755c278a6cb419b8a6ce9b7bf1b4ca0c76304fdce998f10b455b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/markdown/markdown.go,swap,failed,1,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/markdown,^TestReadFenceEdges$,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,356\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/markdown,fail,2,10\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/markdown,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,\"markdown_test.go:68: reading of \\\"````\\\\\\\\n```\\\\\\\\nx\\\\\\\\n````\\\":\",11\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-prose",
          "performer": "ft358-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "780d3f3f9145f6dc8894c30ca4169925915a6fa6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-author-01",
            "digest": "sha256:edf514ced4bac7a7d5ffa5addb939b565b68891fa0aadba3467e6061e5fe8229",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestFindingsSkipsARejoinedComment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/prose,^TestFindingsSkipsARejoinedComment$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/prose,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/prose,TestFindingsSkipsARejoinedComment,\"prose_test.go:30: Findings() = [{Kind:sentence Line:1 Count:31 Starts:[]}], want none\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "252408eb9bce41230a88ccc4ca7a511d4415b848",
          "requirement": "prose",
          "command": "bench test --package ./internal/prose --run '^TestFindingsSkipsARejoinedComment$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/prose/prepare.go: keep each body line's raw text in place of its comment-free text; expect TestFindingsSkipsARejoinedComment",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-author-01",
              "digest": "sha256:edf514ced4bac7a7d5ffa5addb939b565b68891fa0aadba3467e6061e5fe8229",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestFindingsSkipsARejoinedComment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/prose,^TestFindingsSkipsARejoinedComment$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/prose,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/prose,TestFindingsSkipsARejoinedComment,\"prose_test.go:30: Findings() = [{Kind:sentence Line:1 Count:31 Starts:[]}], want none\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-prose-mechanics",
          "performer": "ft358-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "780d3f3f9145f6dc8894c30ca4169925915a6fa6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-author-01",
            "digest": "sha256:76b4ab48580e4cc47294e135b31ab44d16447c311c7d06989679b193d838b42d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,prose-mechanics,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,465\ncheck[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,98,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: prose: \\\".agents/skills/bench-craft-tickets/SKILL.md\\\" line 82: sentence of 26 words is over the 25-word bound\",9\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "252408eb9bce41230a88ccc4ca7a511d4415b848",
          "requirement": "prose-mechanics",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/prose/prepare.go: stop blanking fenced lines, so a live fenced line gives a prose finding; expect TestRootConformance",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-author-01",
              "digest": "sha256:76b4ab48580e4cc47294e135b31ab44d16447c311c7d06989679b193d838b42d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,prose-mechanics,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,465\ncheck[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,98,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: prose: \\\".agents/skills/bench-craft-tickets/SKILL.md\\\" line 82: sentence of 26 words is over the 25-word bound\",9\nskips[0]{package,test,reason}:\n"
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
