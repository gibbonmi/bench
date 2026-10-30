# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/markdown-block-reader/spec.md",
  "plan_digest": "sha256:a7ef5bc715c0df86eb39bf91e71bf940aeedc093cbfa00369948467107736c0d",
  "implementation_session": "",
  "chunks": [
    {
      "id": "MB-C1",
      "base": "36a07d83591f99cc079f571f185d25e606db0a2c",
      "tip": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
      "plan_digest": "sha256:27f3d64fb8abd49d83a5d49138560b8d5ba8dcbb419b54cbbd4729efa68278ce",
      "source_digest": "23e61e60d055c6e595d4a32ab85c7df816423ceb",
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
        },
        {
          "id": "r1-markdown",
          "performer": "ft358-repair-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "1efde0ff74f0f0f1591151bfbd6a1a7aa168d711",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-01",
            "digest": "sha256:a5c5cb4e8924755c278a6cb419b8a6ce9b7bf1b4ca0c76304fdce998f10b455b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/markdown/markdown.go,swap,failed,1,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/markdown,^TestReadFenceEdges$,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,356\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/markdown,fail,2,10\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/markdown,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,\"markdown_test.go:68: reading of \\\"````\\\\\\\\n```\\\\\\\\nx\\\\\\\\n````\\\":\",11\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "503aef49a0de6507e0690eb4161aa0e09a71536d",
          "requirement": "markdown",
          "command": "bench test --package ./internal/markdown --run '^TestReadFenceEdges$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/markdown: make the fence closer accept any run of three or more of the opener's character; expect TestReadFenceEdges",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-01",
              "digest": "sha256:a5c5cb4e8924755c278a6cb419b8a6ce9b7bf1b4ca0c76304fdce998f10b455b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/markdown/markdown.go,swap,failed,1,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/markdown,^TestReadFenceEdges$,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,356\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/markdown,fail,2,10\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/markdown,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,\"markdown_test.go:68: reading of \\\"````\\\\\\\\n```\\\\\\\\nx\\\\\\\\n````\\\":\",11\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r1-prose",
          "performer": "ft358-repair-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "1efde0ff74f0f0f1591151bfbd6a1a7aa168d711",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-01",
            "digest": "sha256:edf514ced4bac7a7d5ffa5addb939b565b68891fa0aadba3467e6061e5fe8229",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestFindingsSkipsARejoinedComment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/prose,^TestFindingsSkipsARejoinedComment$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/prose,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/prose,TestFindingsSkipsARejoinedComment,\"prose_test.go:30: Findings() = [{Kind:sentence Line:1 Count:31 Starts:[]}], want none\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "503aef49a0de6507e0690eb4161aa0e09a71536d",
          "requirement": "prose",
          "command": "bench test --package ./internal/prose --run '^TestFindingsSkipsARejoinedComment$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/prose/prepare.go: keep each body line's raw text in place of its comment-free text; expect TestFindingsSkipsARejoinedComment",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-01",
              "digest": "sha256:edf514ced4bac7a7d5ffa5addb939b565b68891fa0aadba3467e6061e5fe8229",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestFindingsSkipsARejoinedComment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/prose,^TestFindingsSkipsARejoinedComment$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/prose,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/prose,TestFindingsSkipsARejoinedComment,\"prose_test.go:30: Findings() = [{Kind:sentence Line:1 Count:31 Starts:[]}], want none\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r1-prose-mechanics",
          "performer": "ft358-repair-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "1efde0ff74f0f0f1591151bfbd6a1a7aa168d711",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-01",
            "digest": "sha256:8c85b607a7f873e23634bb638a21e98c54df40fd9c188320cb63f07fdffcf122",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,prose-mechanics,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,467\ncheck[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,102,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: prose: \\\".agents/skills/bench-craft-tickets/SKILL.md\\\" line 82: sentence of 26 words is over the 25-word bound\",13\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "503aef49a0de6507e0690eb4161aa0e09a71536d",
          "requirement": "prose-mechanics",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/prose/prepare.go: stop blanking fenced lines, so a live fenced line gives a prose finding; expect TestRootConformance",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-01",
              "digest": "sha256:8c85b607a7f873e23634bb638a21e98c54df40fd9c188320cb63f07fdffcf122",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,prose-mechanics,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,467\ncheck[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,102,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: prose: \\\".agents/skills/bench-craft-tickets/SKILL.md\\\" line 82: sentence of 26 words is over the 25-word bound\",13\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r2-markdown",
          "performer": "ft358-repair-01b",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "23e61e60d055c6e595d4a32ab85c7df816423ceb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-01b",
            "digest": "sha256:a5c5cb4e8924755c278a6cb419b8a6ce9b7bf1b4ca0c76304fdce998f10b455b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/markdown/markdown.go,swap,failed,1,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/markdown,^TestReadFenceEdges$,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,356\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/markdown,fail,2,10\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/markdown,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,\"markdown_test.go:68: reading of \\\"````\\\\\\\\n```\\\\\\\\nx\\\\\\\\n````\\\":\",11\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "requirement": "markdown",
          "command": "bench test --package ./internal/markdown --run '^TestReadFenceEdges$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/markdown: make the fence closer accept any run of three or more of the opener's character; expect TestReadFenceEdges",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-01b",
              "digest": "sha256:a5c5cb4e8924755c278a6cb419b8a6ce9b7bf1b4ca0c76304fdce998f10b455b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/markdown/markdown.go,swap,failed,1,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/markdown,^TestReadFenceEdges$,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,356\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/markdown,fail,2,10\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/markdown,TestReadFenceEdges/MB2_a_shorter_run_stays_inside_a_longer_block,\"markdown_test.go:68: reading of \\\"````\\\\\\\\n```\\\\\\\\nx\\\\\\\\n````\\\":\",11\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r2-prose",
          "performer": "ft358-repair-01b",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "1efde0ff74f0f0f1591151bfbd6a1a7aa168d711",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "carried:r1-prose",
            "digest": "sha256:edf514ced4bac7a7d5ffa5addb939b565b68891fa0aadba3467e6061e5fe8229",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestFindingsSkipsARejoinedComment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/prose,^TestFindingsSkipsARejoinedComment$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/prose,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/prose,TestFindingsSkipsARejoinedComment,\"prose_test.go:30: Findings() = [{Kind:sentence Line:1 Count:31 Starts:[]}], want none\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "503aef49a0de6507e0690eb4161aa0e09a71536d",
          "requirement": "prose",
          "command": "bench test --package ./internal/prose --run '^TestFindingsSkipsARejoinedComment$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/prose/prepare.go: keep each body line's raw text in place of its comment-free text; expect TestFindingsSkipsARejoinedComment",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "carried:r1-prose",
              "digest": "sha256:edf514ced4bac7a7d5ffa5addb939b565b68891fa0aadba3467e6061e5fe8229",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestFindingsSkipsARejoinedComment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/prose,^TestFindingsSkipsARejoinedComment$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/prose,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/prose,TestFindingsSkipsARejoinedComment,\"prose_test.go:30: Findings() = [{Kind:sentence Line:1 Count:31 Starts:[]}], want none\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r2-prose-mechanics",
          "performer": "ft358-repair-01b",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "23e61e60d055c6e595d4a32ab85c7df816423ceb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-01b",
            "digest": "sha256:cd8b0baa8b6964391c3df93a377b7fc139309b9293f25d896bacbaa685811869",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,prose-mechanics,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,467\ncheck[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,100,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: prose: \\\".agents/skills/bench-craft-tickets/SKILL.md\\\" line 82: sentence of 26 words is over the 25-word bound\",17\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "requirement": "prose-mechanics",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/prose/prepare.go: stop blanking fenced lines, so a live fenced line gives a prose finding; expect TestRootConformance",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-01b",
              "digest": "sha256:cd8b0baa8b6964391c3df93a377b7fc139309b9293f25d896bacbaa685811869",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/prose/prepare.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,prose-mechanics,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,467\ncheck[1]{name,kind,tests_run,subjects}:\n  prose-mechanics,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,100,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: prose: \\\".agents/skills/bench-craft-tickets/SKILL.md\\\" line 82: sentence of 26 words is over the 25-word bound\",17\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c1-standards",
          "performer": "ft358-review-c1-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "780d3f3f9145f6dc8894c30ca4169925915a6fa6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c1-standards",
            "digest": "sha256:2215cf9738d3f0425a3b07371589740486ff9404151a8c8ea36acf7068b214b9",
            "excerpt": "Standards axis, MB-C1 (base 36a07d83, tip 252408eb). One finding, low severity; no hard standards violation.\nS1 (low, confidence 6, auto-fix): internal/markdown/markdown.go:187-193 Line struct. The Section doc comment also documents InSection, but InSection is declared after Heading, under the Heading comment, so go doc shows it undocumented and misplaced. Repair: declare InSection directly after Section.\nRefuted: literal ban in internal/prose (test fixtures only); coverage seam flips follow the planned-row convention; dependency standard (sort, strings only); file sizes 229, 102, 268; cutComments duplication is closed by ticket 8 per spec.md:316.\n"
          },
          "axis": "Standards",
          "base": "36a07d83591f99cc079f571f185d25e606db0a2c",
          "tip": "252408eb9bce41230a88ccc4ca7a511d4415b848",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "c1-spec",
          "performer": "ft358-review-c1-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "780d3f3f9145f6dc8894c30ca4169925915a6fa6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c1-spec",
            "digest": "sha256:78c288c11fa6fa181ef000bc923b59b729164e82ff4bdc9dfe396a0fb8f7f6c7",
            "excerpt": "Spec axis, MB-C1 (base 36a07d83, tip 252408eb). One finding; every other mapped row (MB1-MB19, MB71-MB79, MB88-MB90, MB92) and acceptance item matches.\nP1 (medium, confidence 8, auto-fix): internal/markdown/comments.go:43-56 joinedOpener. Spec.md:210-211 and Closed decision 10 require one linear pass. For adjacent comments with no kept byte between them (repeated \u003c!---->), the walk-back steps across every earlier cut, so the work is quadratic. MB79 measures allocation only and its comments are 93 bytes apart, so it stays green. Repair: make the skip over a run of adjacent cuts O(1), and add a run of adjacent comments to the comment-dense input under a bound.\nChecked: comment rule matches stripCommentsMapped; Closed decisions 5, 8, 14, 16; fault order frontmatter, comment, fence; exported constants; prose prepare fault mapping; no banned literal in internal/prose non-test files.\n"
          },
          "axis": "Spec",
          "base": "36a07d83591f99cc079f571f185d25e606db0a2c",
          "tip": "252408eb9bce41230a88ccc4ca7a511d4415b848",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "c1-coverage",
          "performer": "ft358-review-c1-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "780d3f3f9145f6dc8894c30ca4169925915a6fa6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c1-coverage",
            "digest": "sha256:754d97ad9b203352ee61c4025e92893d88067a73c7f5c70c026450e7bbf3047c",
            "excerpt": "Coverage axis, MB-C1 (base 36a07d83, tip 252408eb). Three findings, tests only; production passed six added hostile rows. Every probe restored; git status clean.\nC1 (medium, confidence 9, auto-fix): internal/markdown/comments.go:39, :53-56, :76, :85. Four silent mutants (omit top--, pushCut for->if, after = at+4-n, unterminated branch append for pushCut). Repair: add rows X1 chained rejoin, X2 join at p=1, X3 three kept bytes, X4 unterminated rejoin to TestReadBlockRules.\nC2 (low-medium, confidence 9, auto-fix): internal/markdown/markdown.go:249. CutPrefix on line.Raw in place of line.Text is silent; spec.md:225 and :240-241 require the H2 test on comment-free text. Repair: add row X6.\nC3 (low-medium, confidence 8, auto-fix): internal/markdown/markdown.go:244-247. Omitting the frontmatter continue is silent; a fence marker inside frontmatter then fences the body. Repair: add row X5.\n"
          },
          "axis": "Coverage",
          "base": "36a07d83591f99cc079f571f185d25e606db0a2c",
          "tip": "252408eb9bce41230a88ccc4ca7a511d4415b848",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "c1r-standards",
          "performer": "ft358-confirm-c1-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "1efde0ff74f0f0f1591151bfbd6a1a7aa168d711",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c1-standards",
            "digest": "sha256:ce2219863e6f3f8871a2299093455dd83360fde7ae61aaea1e970b5520e8bf57",
            "excerpt": "Standards axis, MB-C1 confirming round (repair delta 252408eb..503aef49). Zero findings.\nS1 closed: internal/markdown/markdown.go:73-79 declares InSection directly after Section under the shared comment, and the Heading comment sits only over Heading.\nRepair shape checked: cuts.push (comments.go:88-102) trims first with spans; joinedOpener is a pointer method because it writes steps; the cuts type comment states the reason for first; the steps counter (comments.go:48) is test-read only and its comment states why; every edit falls under an MB-C1 finding.\n"
          },
          "axis": "Standards",
          "base": "36a07d83591f99cc079f571f185d25e606db0a2c",
          "tip": "503aef49a0de6507e0690eb4161aa0e09a71536d",
          "finding_ids": [],
          "supersedes": [
            "c1-standards"
          ]
        },
        {
          "id": "c1r-spec",
          "performer": "ft358-confirm-c1-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "1efde0ff74f0f0f1591151bfbd6a1a7aa168d711",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c1-spec",
            "digest": "sha256:d693de3cfd6f91398b89222f6872fcd5a885ebe1d2c8f08b1ad50c0cea6dec6b",
            "excerpt": "Spec axis, MB-C1 confirming round (repair delta 252408eb..503aef49). Zero findings.\nP1 closed: internal/markdown/comments.go:54-68 and :90-102. Spans stay ordered; first[i] lets one jump cross a whole adjacent run; after a jump the next byte is kept, so a walk takes at most 6 steps; total work is O(n), under the MB79 bound of 8 steps per comment. The walk visits exactly the kept positions of the old walk, so the result still equals stripCommentsMapped (internal/anchors/locate.go:58-74).\nRows MB10-MB13, MB88, MB89 unchanged; MB79 gains the work bound of Closed decision 10; spec.md:206-211 now holds. No edit outside the acceptance's reach.\n"
          },
          "axis": "Spec",
          "base": "36a07d83591f99cc079f571f185d25e606db0a2c",
          "tip": "503aef49a0de6507e0690eb4161aa0e09a71536d",
          "finding_ids": [],
          "supersedes": [
            "c1-spec"
          ]
        },
        {
          "id": "c1r-coverage",
          "performer": "ft358-confirm-c1-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "1efde0ff74f0f0f1591151bfbd6a1a7aa168d711",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c1-coverage",
            "digest": "sha256:6e4ee5e02f50ed2e1c4038d9c29ab18ff487b56cf321202644190bff54cf7bce",
            "excerpt": "Coverage axis, MB-C1 confirming round (repair delta 252408eb..503aef49). C1, C2, C3 closed; one new finding.\nRound 1 mutants re-run, all red: omit top-- (comments.go:62), for->if in push (:91), at+len(commentOpen)-n (:82), append in the unterminated branch (:33), CutPrefix on line.Raw (markdown.go:136), omit the frontmatter continue (markdown.go:133), two-byte kept window (comments.go:56).\nNew mutants: omit top = c.first[top] red (MB79); first = i red (MB79); End \u003c= s.Start red (MB12).\nN1 (medium, confidence 8, auto-fix): internal/markdown/comments.go:95. Omitting c.first = c.first[:i] is silent. Witness input \u003c!\u003c!---->\u003c!---->---->\u003c!\u003c!---->----> reads comment 0-21 and comment 21-35 with empty text; the mutant reads comments 0-21 and 23-30 with visible \u003c!---->. Repair: add that row to TestReadBlockRules.\n"
          },
          "axis": "Coverage",
          "base": "36a07d83591f99cc079f571f185d25e606db0a2c",
          "tip": "503aef49a0de6507e0690eb4161aa0e09a71536d",
          "finding_ids": [
            "N1"
          ],
          "supersedes": [
            "c1-coverage"
          ]
        }
      ],
      "confirmations": [
        {
          "id": "c1-n1-confirmation",
          "performer": "session_01JK5CPbT18o2hCYe8SVhPcj",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "medium",
          "source_digest": "23e61e60d055c6e595d4a32ab85c7df816423ceb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:session_01JK5CPbT18o2hCYe8SVhPcj:n1-probe",
            "digest": "sha256:8729f6ca95b5b9d966ebceb870e0ea64d9ecde4b80b1a802ecf998e3a8e9116c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/markdown/comments.go,omit,failed,1,TestReadBlockRules/MB12_a_rejoin_after_a_rejoin_of_two_adjacent_removals_is_removed_too,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/markdown,^TestReadBlockRules$,passed,32\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,386\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/markdown,fail,5,32\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/markdown,TestReadBlockRules/MB12_a_rejoin_after_a_rejoin_of_two_adjacent_removals_is_removed_too,\"markdown_test.go:68: reading of \\\"\u003c!\u003c!---->\u003c!---->---->\u003c!\u003c!---->---->\\\":\",8\nskips[0]{package,test,reason}:\n"
          },
          "base": "36a07d83591f99cc079f571f185d25e606db0a2c",
          "reviewed_tip": "503aef49a0de6507e0690eb4161aa0e09a71536d",
          "repaired_tip": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "finding_ids": [
            "N1"
          ],
          "supersedes": [],
          "probe": {
            "mutation": "internal/markdown/comments.go: omit c.first = c.first[:i]",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:session_01JK5CPbT18o2hCYe8SVhPcj:n1-probe",
              "digest": "sha256:8729f6ca95b5b9d966ebceb870e0ea64d9ecde4b80b1a802ecf998e3a8e9116c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/markdown/comments.go,omit,failed,1,TestReadBlockRules/MB12_a_rejoin_after_a_rejoin_of_two_adjacent_removals_is_removed_too,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/markdown,^TestReadBlockRules$,passed,32\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,386\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/markdown,fail,5,32\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/markdown,TestReadBlockRules/MB12_a_rejoin_after_a_rejoin_of_two_adjacent_removals_is_removed_too,\"markdown_test.go:68: reading of \\\"\u003c!\u003c!---->\u003c!---->---->\u003c!\u003c!---->---->\\\":\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ]
    },
    {
      "id": "MB-C2",
      "base": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
      "tip": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
      "plan_digest": "sha256:4b0b90422efdddcf683c1eb5a32d0bf68d153fed48f264671dfebc96e810e3db",
      "source_digest": "ce1c5f482491cc109c94105f85ea9955ad3736ea",
      "acceptance_rows": [
        "MB24",
        "MB25",
        "MB26",
        "MB29",
        "MB30",
        "MB80",
        "MB83",
        "MB31",
        "MB32",
        "MB33",
        "MB81",
        "MB99",
        "MB100",
        "MB103",
        "MB34",
        "MB35",
        "MB36",
        "MB37",
        "MB95",
        "MB102"
      ],
      "verification": [
        {
          "id": "t2-spec",
          "performer": "ft358-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "956cb86cb6530e6d9fd110cc26d195b1ff130db4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-02",
            "digest": "sha256:148ec91d985dbd280b3a968387fa1498e6dd74a7fea97ba07b1cb0b61635bc44",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/spec.go,swap,failed,2,TestMetadataSkipsAFencedStatus/indented_fence,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,^TestMetadataSkipsAFencedStatus$,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,481\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,2,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestMetadataSkipsAFencedStatus/indented_fence,\"header_fence_test.go:15: LiveSpecSlugs = [\\\"hidden\\\"], want none\",2\n  github.com/gibbonmi/bench/internal/spec,TestMetadataSkipsAFencedStatus/tilde_fence,\"header_fence_test.go:15: LiveSpecSlugs = [\\\"hidden\\\"], want none\",2\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "08479e59b53cbac724b91d0a2955782d46d03cc8",
          "requirement": "spec",
          "command": "bench test --package ./internal/spec --run '^TestMetadataSkipsAFencedStatus$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/spec/spec.go: make headerLines ignore the reader's fence class; expect TestMetadataSkipsAFencedStatus",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-02",
              "digest": "sha256:148ec91d985dbd280b3a968387fa1498e6dd74a7fea97ba07b1cb0b61635bc44",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/spec/spec.go,swap,failed,2,TestMetadataSkipsAFencedStatus/indented_fence,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/spec,^TestMetadataSkipsAFencedStatus$,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,481\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,fail,2,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/spec,TestMetadataSkipsAFencedStatus/indented_fence,\"header_fence_test.go:15: LiveSpecSlugs = [\\\"hidden\\\"], want none\",2\n  github.com/gibbonmi/bench/internal/spec,TestMetadataSkipsAFencedStatus/tilde_fence,\"header_fence_test.go:15: LiveSpecSlugs = [\\\"hidden\\\"], want none\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2-coverage",
          "performer": "ft358-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "956cb86cb6530e6d9fd110cc26d195b1ff130db4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-02",
            "digest": "sha256:e89536ec4d293bbf36c83ff2d267dc55904ea9ea10028c4e42f6a9ad7da87ed0",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/coverage/rows.go,swap,failed,1,TestParseSkipsFencedCoverageLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/coverage,^TestParseSkipsFencedCoverageLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,399\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/coverage,TestParseSkipsFencedCoverageLines,\"coverage_test.go:286: Check = [coverage map leaves story 2 unreferenced; add a row or a `Not covered: story 2 — \u003creason>` line], want none\",2\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "08479e59b53cbac724b91d0a2955782d46d03cc8",
          "requirement": "coverage",
          "command": "bench test --package ./internal/coverage --run '^TestParseSkipsFencedCoverageLines$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/coverage/rows.go: drop the fenced-line skip in parse; expect TestParseSkipsFencedCoverageLines",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-02",
              "digest": "sha256:e89536ec4d293bbf36c83ff2d267dc55904ea9ea10028c4e42f6a9ad7da87ed0",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/coverage/rows.go,swap,failed,1,TestParseSkipsFencedCoverageLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/coverage,^TestParseSkipsFencedCoverageLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,399\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/coverage,TestParseSkipsFencedCoverageLines,\"coverage_test.go:286: Check = [coverage map leaves story 2 unreferenced; add a row or a `Not covered: story 2 — \u003creason>` line], want none\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3-roadmap",
          "performer": "ft358-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "5fe25d8a9c17aabd33a4fef9447d69a56b240c87",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-03",
            "digest": "sha256:faac94ce3104d1126decacad6d94cfdbc68ba9881a4d8959f0c5e2af0a2d8ea9",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:165: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", \u003cnil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5c654076827727cd0609f18f78d4d251c65b4759",
          "requirement": "roadmap",
          "command": "bench test --package ./internal/roadmap --run '^TestProjectSequenceSkipsFencedHeadings$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/roadmap/sequence_projection.go: make sectionBounds ignore the reader's fence class; expect TestProjectSequenceSkipsFencedHeadings",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-03",
              "digest": "sha256:faac94ce3104d1126decacad6d94cfdbc68ba9881a4d8959f0c5e2af0a2d8ea9",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:165: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", \u003cnil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t4-maps",
          "performer": "ft358-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "540aea0cabec043307afb28b3e8ad136d026e091",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-04",
            "digest": "sha256:3f4b335f551a418a84db54df36c33ce60a8add52178e98526e7c80edc1941dc3",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/maps/fields.go,swap,failed,1,TestFieldScanMarksTildeFencedLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/maps,^TestFieldScanMarksTildeFencedLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,333\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/maps,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/maps,TestFieldScanMarksTildeFencedLines,\"tickets_test.go:278: Scan diagnostics = []string{\\\"duplicate Writes\\\"}, want none\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "cd86cf1efca9403f1d2908aadd32612e62857e42",
          "requirement": "maps",
          "command": "bench test --package ./internal/maps --run '^TestFieldScanMarksTildeFencedLines$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/maps/fields.go: mark only backtick-fenced lines as fenced; expect TestFieldScanMarksTildeFencedLines",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-04",
              "digest": "sha256:3f4b335f551a418a84db54df36c33ce60a8add52178e98526e7c80edc1941dc3",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/maps/fields.go,swap,failed,1,TestFieldScanMarksTildeFencedLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/maps,^TestFieldScanMarksTildeFencedLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,333\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/maps,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/maps,TestFieldScanMarksTildeFencedLines,\"tickets_test.go:278: Scan diagnostics = []string{\\\"duplicate Writes\\\"}, want none\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t4-tickets",
          "performer": "ft358-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "540aea0cabec043307afb28b3e8ad136d026e091",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-04",
            "digest": "sha256:cbb4803ea5e1e022f51a17aeb9f4ee6baea75ee934ec2dd1287b2274a8185821",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/tickets/tickets.go,swap,failed,3,TestParseTicketReadsTheBlockReaderFault/MB35_an_unterminated_tilde_block,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/tickets,^TestParseTicketReadsTheBlockReaderFault$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,873\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,fail,2,6\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/tickets,TestParseTicketReadsTheBlockReaderFault/MB35_an_unterminated_tilde_block,\"fence_test.go:50: diagnostics = []string(nil), want []string{\\\"unterminated fence\\\"}\",1\n  github.com/gibbonmi/bench/internal/tickets,TestParseTicketReadsTheBlockReaderFault/MB36_a_four-backtick_block_that_holds_a_three-backtick_line,\"fence_test.go:50: diagnostics = []string{\\\"unterminated fence\\\"}, want []string(nil)\",1\n  github.com/gibbonmi/bench/internal/tickets,TestParseTicketReadsTheBlockReaderFault/MB37_a_Covers_line_inside_the_four-backtick_block,\"fence_test.go:50: diagnostics = []string{\\\"unterminated fence\\\", \\\"missing Covers\\\"}, want []string{\\\"missing Covers\\\"}\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "cd86cf1efca9403f1d2908aadd32612e62857e42",
          "requirement": "tickets",
          "command": "bench test --package ./internal/tickets --run '^TestParseTicketReadsTheBlockReaderFault$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/tickets/tickets.go: report an unterminated fence on an odd count of backtick marker lines; expect TestParseTicketReadsTheBlockReaderFault",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-04",
              "digest": "sha256:cbb4803ea5e1e022f51a17aeb9f4ee6baea75ee934ec2dd1287b2274a8185821",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/tickets/tickets.go,swap,failed,3,TestParseTicketReadsTheBlockReaderFault/MB35_an_unterminated_tilde_block,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/tickets,^TestParseTicketReadsTheBlockReaderFault$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,873\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,fail,2,6\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/tickets,TestParseTicketReadsTheBlockReaderFault/MB35_an_unterminated_tilde_block,\"fence_test.go:50: diagnostics = []string(nil), want []string{\\\"unterminated fence\\\"}\",1\n  github.com/gibbonmi/bench/internal/tickets,TestParseTicketReadsTheBlockReaderFault/MB36_a_four-backtick_block_that_holds_a_three-backtick_line,\"fence_test.go:50: diagnostics = []string{\\\"unterminated fence\\\"}, want []string(nil)\",1\n  github.com/gibbonmi/bench/internal/tickets,TestParseTicketReadsTheBlockReaderFault/MB37_a_Covers_line_inside_the_four-backtick_block,\"fence_test.go:50: diagnostics = []string{\\\"unterminated fence\\\", \\\"missing Covers\\\"}, want []string{\\\"missing Covers\\\"}\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r2-roadmap",
          "performer": "ft358-repair-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "86d28a1c0b9b0e25f9b79741a8a669daf3bfbf47",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-repair-02",
            "digest": "sha256:dd395c7ec98a337b102d4f5cdfd3bb22e5fb31fd24ca39a2ed234d35fd09de93",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:170: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", \u003cnil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "592040eab58cc8bdcc388aa0973757a6c86cabd4",
          "requirement": "roadmap",
          "command": "bench test --package ./internal/roadmap --run '^TestProjectSequenceSkipsFencedHeadings$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/roadmap/sequence_projection.go: make sectionBounds ignore the reader's fence class; expect TestProjectSequenceSkipsFencedHeadings",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-repair-02",
              "digest": "sha256:dd395c7ec98a337b102d4f5cdfd3bb22e5fb31fd24ca39a2ed234d35fd09de93",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:170: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", \u003cnil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r3-roadmap",
          "performer": "ft358-repair-02b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ce1c5f482491cc109c94105f85ea9955ad3736ea",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-repair-02b",
            "digest": "sha256:dd395c7ec98a337b102d4f5cdfd3bb22e5fb31fd24ca39a2ed234d35fd09de93",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:170: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", \u003cnil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "requirement": "roadmap",
          "command": "bench test --package ./internal/roadmap --run '^TestProjectSequenceSkipsFencedHeadings$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/roadmap/sequence_projection.go: make sectionBounds ignore the reader's fence class; expect TestProjectSequenceSkipsFencedHeadings",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-repair-02b",
              "digest": "sha256:dd395c7ec98a337b102d4f5cdfd3bb22e5fb31fd24ca39a2ed234d35fd09de93",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:170: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", \u003cnil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c2-standards",
          "performer": "ft358-review-c2-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "540aea0cabec043307afb28b3e8ad136d026e091",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c2-standards",
            "digest": "sha256:30e08c710115aa177b06d0344236373f93731a36ab180e6f48c2593be1587d3a",
            "excerpt": "Standards axis, MB-C2 (base 2c79f46e, tip cd86cf1e). Three findings.\nPer-consumer Class != markdown.Body tests are not duplicated knowledge: the enum owns the block fact, each site states its own grammar. maps.Blocks is the single source of the ReadFragment choice that Scan and ParseTicket share (no-op).\nS1 (medium, high confidence, ask-user): internal/roadmap/closure.go:19 and :27. Close matches roadmapStartRe on every line, fenced or not, while its doc comment says a row heading is what ParseDocument reads as one; ParseDocument now skips fenced rows (tree.go:116). Repair: Close skips non-Body lines; widen ticket 3 Writes to internal/roadmap/closure.go.\nS2 (low, medium confidence, ask-user): internal/tickets/tickets.go:108. ParseTicket tests only FenceFault; a comment fault comes first (markdown.go:142-148), so an unclosed comment hides an unterminated fence that the old count reported.\nS3 (low, high confidence, no-op): tree.go:92, :274, sequence_projection.go:27 read the index three times; the parallel strings.Split matches by construction.\n"
          },
          "axis": "Standards",
          "base": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "tip": "cd86cf1efca9403f1d2908aadd32612e62857e42",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c2-spec",
          "performer": "ft358-review-c2-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "540aea0cabec043307afb28b3e8ad136d026e091",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c2-spec",
            "digest": "sha256:ef1202045ecbc8b470639b25412b49fcfbd5ea897da0809a122dec93f1328216",
            "excerpt": "Spec axis, MB-C2 (base 2c79f46e, tip cd86cf1e). Two findings; rows MB24-26, MB29-37, MB80, MB81, MB83, MB95, MB99, MB100, MB102 hold.\nP1 (severity 5, confidence 10, auto-fix): internal/coverage/rows.go grows from 301 to 303 lines; ticket 2 and spec lines 341-344 require no net growth. Repair: cut two net lines in parse.\nP2 (severity 4, confidence 8, ask-user): internal/roadmap/closure.go:19 and :27. Close counts a fenced row heading that ParseDocument now skips (FA6), so a fenced example beside a live row refuses, and its doc comment is false. Outside every Writes line.\nFlag b: a non-behavioral contradiction; the behavior clauses (lines 288-291, Closed decision 16, MB95) require ReadFragment for tickets, so lines 293-296 and the Closed decision 16 caller list need a prose correction.\nAdvice: flag a is intended by story 9; flag c follows from the MB71 fault order and is not required by any row.\n"
          },
          "axis": "Spec",
          "base": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "tip": "cd86cf1efca9403f1d2908aadd32612e62857e42",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "c2-coverage",
          "performer": "ft358-review-c2-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "540aea0cabec043307afb28b3e8ad136d026e091",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c2-coverage",
            "digest": "sha256:4e679d779a9c92de4278fdae61186b4e1ed86bbc15b72b6efc396a0ba7d0d3ea",
            "excerpt": "Coverage axis, MB-C2 (base 2c79f46e, tip cd86cf1e). Five findings; the six named packages pass; seven probes restored; git status clean.\nC1 (medium, confidence 8, auto-fix): spec.go:67, coverage/rows.go:133, roadmap/tree.go:92 Read->ReadFragment swaps are silent; no test pins the frontmatter rule in these readers. Add a first-line --- ... --- case to each module test.\nC2 (medium, confidence 9, auto-fix): tickets.go:108 fault-kind filter swap to any fault is silent across 11151 tests. Add an unterminated-comment case that expects no unterminated fence.\nC3 (medium, confidence 9, auto-fix): tree.go:276 fence skip inside the live sequence is silent. Add a fenced 1. fake line inside the live sequence.\nC4 (low, confidence 7, auto-fix): sequence_projection.go:32 title compare tests only the accepted side. Add a rejected-side case.\nC5 (low, confidence 7, ask-user): a longer closer inside the planning-facts fence now closes with an empty label (maps/fields.go:92), so metadata.go:63 truncates the payload with no refusal; spec.md:330 says a nested marker still refuses.\n"
          },
          "axis": "Coverage",
          "base": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "tip": "cd86cf1efca9403f1d2908aadd32612e62857e42",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
            "C5"
          ],
          "supersedes": []
        },
        {
          "id": "c2r-standards",
          "performer": "ft358-confirm-c2-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "86d28a1c0b9b0e25f9b79741a8a669daf3bfbf47",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c2-standards",
            "digest": "sha256:34fd988c26ea2020153e934f1e274bb3552f4f4e92d9d4eab5f185a5693ed8e0",
            "excerpt": "Standards axis, MB-C2 confirming round (repair delta cd86cf1e..592040ea). S1 closed; two new findings.\nS1 closed: internal/roadmap/closure.go:31 counts a row heading only on a Body line read through markdown.Read; the doc comment at :21-24 is true.\nS4 (low-medium, high confidence, auto-fix): closure.go:31 and tree.go:116 + :123 state the row-heading rule (Body line plus roadmapStartRe match) twice, against AGENTS.md one source per fact. Repair: one rowHeading helper in internal/roadmap that both call.\nS5 (low, medium confidence, no-op): internal/maps/fields.go openerRun derives the run from the reader's Text and Info and copies no marker rule; a run-length field on markdown.Line would be cleaner but is a reader API change.\n"
          },
          "axis": "Standards",
          "base": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "tip": "592040eab58cc8bdcc388aa0973757a6c86cabd4",
          "finding_ids": [
            "S4"
          ],
          "supersedes": [
            "c2-standards"
          ]
        },
        {
          "id": "c2r-spec",
          "performer": "ft358-confirm-c2-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "86d28a1c0b9b0e25f9b79741a8a669daf3bfbf47",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c2-spec",
            "digest": "sha256:da5cc9cf26a516ff46454b4eb5729579a0e8dc9ee545518a323ae3ec15f81b97",
            "excerpt": "Spec axis, MB-C2 confirming round (repair delta cd86cf1e..592040ea). P1 and P2 closed; one new finding.\nP1 closed: internal/coverage/rows.go is 301 lines, numstat 16/16 from the base. P2 closed: Close matches a row only on a Body line, aligned with ParseDocument. C5 matches the corrected spec text and MB102; MB103 holds; bench coverage --check is valid with 100 rows; FA10 and FA11 hold.\nP3 (low, confidence 8, auto-fix): specs/markdown-block-reader/tickets/3-read-roadmap-blocks.md names MB103 in Covers but not in What to build or Acceptance. Repair: add one sentence and one acceptance line.\n"
          },
          "axis": "Spec",
          "base": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "tip": "592040eab58cc8bdcc388aa0973757a6c86cabd4",
          "finding_ids": [
            "P3"
          ],
          "supersedes": [
            "c2-spec"
          ]
        },
        {
          "id": "c2r-coverage",
          "performer": "ft358-confirm-c2-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "86d28a1c0b9b0e25f9b79741a8a669daf3bfbf47",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c2-coverage",
            "digest": "sha256:9bfa148b2d2a90e112fd456ddbe7535c292166867772663e200f9758f5ac36fe",
            "excerpt": "Coverage axis, MB-C2 confirming round (repair delta cd86cf1e..592040ea). C1-C5 closed (each round 1 mutant now red); two new findings. Every probe restored; git status clean.\nN1 (low, confidence 8, auto-fix): internal/maps/fields.go:96, :98, :145-147. Silent mutants: run = markdown.FenceMarker; strings.CutPrefix for strings.Cut; openerRun without the outer TrimSpace. Repair: add a four-run opener with a four-run closer, an indented closer, and an opener with a space before its info string plus a longer closer to TestMetadataGrammarFenceRoles.\nN2 (low, confidence 7, auto-fix): internal/roadmap/closure.go:22-23 claims frontmatter parity with ParseDocument, but Read->ReadFragment at closure.go:26 is silent. Repair: add a first-line --- block holding **FT1 — x** to TestCloseSkipsAFencedRowHeading.\n"
          },
          "axis": "Coverage",
          "base": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "tip": "592040eab58cc8bdcc388aa0973757a6c86cabd4",
          "finding_ids": [
            "N1",
            "N2"
          ],
          "supersedes": [
            "c2-coverage"
          ]
        },
        {
          "id": "c2r2-standards",
          "performer": "ft358-confirm2-c2-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ce1c5f482491cc109c94105f85ea9955ad3736ea",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm2-c2-standards",
            "digest": "sha256:031a04251ef0fcb4199dcbe74ec43b539718119d8d0af5adfb26b2c23906664d",
            "excerpt": "Standards axis, MB-C2 confirming round 2 (repair delta 592040ea..60f0d9ee). Zero findings.\nS4 closed: rowHeading (internal/roadmap/tree.go:205-212) is the one place that combines the body-class test with roadmapStartRe; ParseDocument (tree.go:116) and Close (closure.go:30) both call it; behavior unchanged; comments meet craft-comments.\nOptional advice (no finding, no repair): premises_test.go:99-101 packs four unnamed closer-label cases on one line, and a named-case table would read better; TestCloseSkipsAFencedRowHeading now also covers a frontmatter example, so a name such as TestCloseSkipsNonBodyRowHeadings would fit both cases.\n"
          },
          "axis": "Standards",
          "base": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "tip": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "finding_ids": [],
          "supersedes": [
            "c2r-standards"
          ]
        },
        {
          "id": "c2r2-spec",
          "performer": "ft358-confirm2-c2-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ce1c5f482491cc109c94105f85ea9955ad3736ea",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm2-c2-spec",
            "digest": "sha256:feb72f0260a7188b38a9b518820b32a075384a3adaee351d13a5a8e61bbbfdb0",
            "excerpt": "Spec axis, MB-C2 confirming round 2 (repair delta 592040ea..60f0d9ee). Zero findings.\nP3 closed: ticket 3 names MB103 in What to build (:27-29) and in Acceptance (:44).\nThe rowHeading refactor keeps ParseDocument (tree.go:114-131) and Close (closure.go:30) behavior exactly; both read through markdown.Read, so FA10 and FA12 hold; MB31, MB81, FA6 hold. tree.go numstat from the chunk base is 36/36, so it does not grow. The premises_test.go cases are in scope.\n"
          },
          "axis": "Spec",
          "base": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "tip": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "finding_ids": [],
          "supersedes": [
            "c2r-spec"
          ]
        },
        {
          "id": "c2r2-coverage",
          "performer": "ft358-confirm2-c2-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ce1c5f482491cc109c94105f85ea9955ad3736ea",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm2-c2-coverage",
            "digest": "sha256:c8046d90e4fcc8bba69a5174da7513c0bbe8bd84d0e37cca68c8d4659ebd3a09",
            "excerpt": "Coverage axis, MB-C2 confirming round 2 (repair delta 592040ea..60f0d9ee). Zero findings. Every probe restored; git status clean.\nN1 closed: run = markdown.FenceMarker, strings.CutPrefix for strings.Cut, and openerRun without the outer TrimSpace each fail their own case at premises_test.go:102.\nN2 closed: markdown.ReadFragment in Close fails at commitment_test.go:186.\nMutants on rowHeading (tree.go:208-209) and both callers (closure.go:30, tree.go:116-117) are all red.\n"
          },
          "axis": "Coverage",
          "base": "2c79f46e4dd6031016a82e0cc335762ae3ffc6b8",
          "tip": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "finding_ids": [],
          "supersedes": [
            "c2r-coverage"
          ]
        }
      ]
    },
    {
      "id": "MB-C3",
      "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
      "tip": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
      "plan_digest": "sha256:c31bcf3103d4d0c4225c463e4b60f241de6c3c6b338663b5093a9bf6b765733a",
      "source_digest": "0f6e4939b773f5cf85c949e9084f7a32f43f7b17",
      "acceptance_rows": [
        "MB38",
        "MB39",
        "MB40",
        "MB82",
        "MB93",
        "MB41",
        "MB42",
        "MB43",
        "MB44",
        "MB45",
        "MB46",
        "MB47",
        "MB84",
        "MB104",
        "MB105",
        "MB106",
        "MB59",
        "MB60"
      ],
      "verification": [
        {
          "id": "t5-handoffdoc",
          "performer": "ft358-author-05",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a80577a8fa012564e7f053eafb36716acdbd45e7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-05",
            "digest": "sha256:1a480f97a9f6ca1423b8c6e34a732dcdf30482dbff1fcdf62c98549aca1d7fe2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/handoffdoc/fence.go,swap,failed,4,TestOpenFenceUsesTheCloseRule/longer_run_closes,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoffdoc,^TestOpenFenceUsesTheCloseRule$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,870\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,fail,2,6\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestOpenFenceUsesTheCloseRule/longer_run_closes,\"document_test.go:269: longer run closes: OpenFence = (3, true), want (0, false)\",1\n  github.com/gibbonmi/bench/internal/handoffdoc,TestOpenFenceUsesTheCloseRule/open_comment,\"document_test.go:269: open comment: OpenFence = (1, true), want (0, false)\",1\n  github.com/gibbonmi/bench/internal/handoffdoc,TestOpenFenceUsesTheCloseRule/other_character,\"document_test.go:269: other character: OpenFence = (0, false), want (1, true)\",1\n  github.com/gibbonmi/bench/internal/handoffdoc,TestOpenFenceUsesTheCloseRule/shorter_run,\"document_test.go:269: shorter run: OpenFence = (0, false), want (1, true)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e2cda7748aa220e019b688c69a6d3063a5214f89",
          "requirement": "handoffdoc",
          "command": "bench test --package ./internal/handoffdoc --run '^TestOpenFenceUsesTheCloseRule$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/handoffdoc/fence.go: close an open block on any marker line; expect TestOpenFenceUsesTheCloseRule",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-05",
              "digest": "sha256:1a480f97a9f6ca1423b8c6e34a732dcdf30482dbff1fcdf62c98549aca1d7fe2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/handoffdoc/fence.go,swap,failed,4,TestOpenFenceUsesTheCloseRule/longer_run_closes,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoffdoc,^TestOpenFenceUsesTheCloseRule$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,870\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,fail,2,6\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestOpenFenceUsesTheCloseRule/longer_run_closes,\"document_test.go:269: longer run closes: OpenFence = (3, true), want (0, false)\",1\n  github.com/gibbonmi/bench/internal/handoffdoc,TestOpenFenceUsesTheCloseRule/open_comment,\"document_test.go:269: open comment: OpenFence = (1, true), want (0, false)\",1\n  github.com/gibbonmi/bench/internal/handoffdoc,TestOpenFenceUsesTheCloseRule/other_character,\"document_test.go:269: other character: OpenFence = (0, false), want (1, true)\",1\n  github.com/gibbonmi/bench/internal/handoffdoc,TestOpenFenceUsesTheCloseRule/shorter_run,\"document_test.go:269: shorter run: OpenFence = (0, false), want (1, true)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t5-handoff",
          "performer": "ft358-author-05",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a80577a8fa012564e7f053eafb36716acdbd45e7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-05",
            "digest": "sha256:45188b7031c63fbc19bf6d387312b9d23a5e332a9f19cc3a4dcaf5d5e14dd2eb",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/handoff/state_file.go,omit,failed,2,TestStateFileRefusesAnUnterminatedComment/open_comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoff,^TestStateFileRefusesAnUnterminatedComment$,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,831\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoff,fail,4,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoff,TestStateFileRefusesAnUnterminatedComment/open_comment,\"state_file_test.go:249: open comment: readStateFile refusal kind \\\"the handoff state file opens a section heading\\\" (hint names the opener: false), want \\\"the handoff state file opens an HTML comment it never closes\\\" naming the opener line\",1\n  github.com/gibbonmi/bench/internal/handoff,TestStateFileRefusesAnUnterminatedComment/open_fence_above_it,\"state_file_test.go:249: open fence above it: readStateFile refusal kind \\\"\\\" (hint names the opener: false), want \\\"the handoff state file opens an HTML comment it never closes\\\" naming the opener line\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e2cda7748aa220e019b688c69a6d3063a5214f89",
          "requirement": "handoff",
          "command": "bench test --package ./internal/handoff --run '^TestStateFileRefusesAnUnterminatedComment$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/handoff/state_file.go: omit the comment-fault refusal in readStateFile; expect TestStateFileRefusesAnUnterminatedComment",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-05",
              "digest": "sha256:45188b7031c63fbc19bf6d387312b9d23a5e332a9f19cc3a4dcaf5d5e14dd2eb",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/handoff/state_file.go,omit,failed,2,TestStateFileRefusesAnUnterminatedComment/open_comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoff,^TestStateFileRefusesAnUnterminatedComment$,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,831\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoff,fail,4,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoff,TestStateFileRefusesAnUnterminatedComment/open_comment,\"state_file_test.go:249: open comment: readStateFile refusal kind \\\"the handoff state file opens a section heading\\\" (hint names the opener: false), want \\\"the handoff state file opens an HTML comment it never closes\\\" naming the opener line\",1\n  github.com/gibbonmi/bench/internal/handoff,TestStateFileRefusesAnUnterminatedComment/open_fence_above_it,\"state_file_test.go:249: open fence above it: readStateFile refusal kind \\\"\\\" (hint names the opener: false), want \\\"the handoff state file opens an HTML comment it never closes\\\" naming the opener line\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t6-learnings",
          "performer": "ft358-author-06",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "35971c128fe7f207e4c16fa1a02f69237a952aae",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-06",
            "digest": "sha256:24af14ac1b97501cebd02531dbe8bda3b52ceae948fe01bdd0fddb4c973da20c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/learnings/learnings.go,swap,failed,3,TestParseReportsAnUnterminatedBlock/comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/learnings,^TestParseReportsAnUnterminatedBlock$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,782\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/learnings,fail,2,4\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/comment,\"learnings_test.go:389: comment: Parse gives 0 entries and 2 malformed records, want 0 and 3 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/fence,\"learnings_test.go:389: fence: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/frontmatter,\"learnings_test.go:389: frontmatter: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ffc1d806ccde0e7e73fe12199dfb574d660911d5",
          "requirement": "learnings",
          "command": "bench test --package ./internal/learnings --run '^TestParseReportsAnUnterminatedBlock$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/learnings: omit the malformed record for the reader fault; expect TestParseReportsAnUnterminatedBlock",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-06",
              "digest": "sha256:24af14ac1b97501cebd02531dbe8bda3b52ceae948fe01bdd0fddb4c973da20c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/learnings/learnings.go,swap,failed,3,TestParseReportsAnUnterminatedBlock/comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/learnings,^TestParseReportsAnUnterminatedBlock$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,782\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/learnings,fail,2,4\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/comment,\"learnings_test.go:389: comment: Parse gives 0 entries and 2 malformed records, want 0 and 3 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/fence,\"learnings_test.go:389: fence: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/frontmatter,\"learnings_test.go:389: frontmatter: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t6-retros",
          "performer": "ft358-author-06",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "35971c128fe7f207e4c16fa1a02f69237a952aae",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-06",
            "digest": "sha256:90153edbc2df8f7083b9258d0af1637d7b2af06454703b168348698ab5815f96",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/retros/recommendations.go,swap,failed,1,TestRecommendationsSkipFencedLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/retros,^TestRecommendationsSkipFencedLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,369\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/retros,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/retros,TestRecommendationsSkipFencedLines,\"recommendations_test.go:33: fenced lines: Recommendations gives 1 units or other units, want the 3 unfenced units\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ffc1d806ccde0e7e73fe12199dfb574d660911d5",
          "requirement": "retros",
          "command": "bench test --package ./internal/retros --run '^TestRecommendationsSkipFencedLines$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/retros/recommendations.go: end the improvements on any H2 line, fenced or not; expect TestRecommendationsSkipFencedLines",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-06",
              "digest": "sha256:90153edbc2df8f7083b9258d0af1637d7b2af06454703b168348698ab5815f96",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/retros/recommendations.go,swap,failed,1,TestRecommendationsSkipFencedLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/retros,^TestRecommendationsSkipFencedLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,369\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/retros,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/retros,TestRecommendationsSkipFencedLines,\"recommendations_test.go:33: fenced lines: Recommendations gives 1 units or other units, want the 3 unfenced units\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t9-skillsindex",
          "performer": "ft358-author-09",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "183a350962746062b7ffb6fc133d42339b3ace33",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-09",
            "digest": "sha256:b4f6e7132ddd0a25db252dcbcb1d8d548b2f3fe8e52c3e8e9a113ac59637361a",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/skillsindex/skillsindex.go,insert,failed,1,TestFrontmatterFieldReadsCRLF/CRLF_block,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/skillsindex,^TestFrontmatterFieldReadsCRLF$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,332\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/skillsindex,fail,2,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/skillsindex,TestFrontmatterFieldReadsCRLF/CRLF_block,\"frontmatter_test.go:23: CRLF block: index = \\\"\\\", want \\\"x\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "251b6e3adf4950681b0227de1fc93449dc8910ef",
          "requirement": "skillsindex",
          "command": "bench test --package ./internal/skillsindex --run '^TestFrontmatterFieldReadsCRLF$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/skillsindex/skillsindex.go: compare the first line to the frontmatter marker before the carriage-return strip; expect TestFrontmatterFieldReadsCRLF",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-09",
              "digest": "sha256:b4f6e7132ddd0a25db252dcbcb1d8d548b2f3fe8e52c3e8e9a113ac59637361a",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/skillsindex/skillsindex.go,insert,failed,1,TestFrontmatterFieldReadsCRLF/CRLF_block,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/skillsindex,^TestFrontmatterFieldReadsCRLF$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,332\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/skillsindex,fail,2,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/skillsindex,TestFrontmatterFieldReadsCRLF/CRLF_block,\"frontmatter_test.go:23: CRLF block: index = \\\"\\\", want \\\"x\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r3-learnings",
          "performer": "ft358-repair-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9eceeda011572dab56efc7b5241150fa3518bbe0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-repair-03",
            "digest": "sha256:2a5cf2c0e2d8f5a815bca77a558bf429b0546b6449c1978227d6cc54c4915cc4",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/learnings/learnings.go,swap,failed,3,TestParseReportsAnUnterminatedBlock/comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/learnings,^TestParseReportsAnUnterminatedBlock$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,782\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/learnings,fail,2,4\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/comment,\"learnings_test.go:396: comment: Parse gives 0 entries and 2 malformed records, want 0 and 3 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/fence,\"learnings_test.go:396: fence: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/frontmatter,\"learnings_test.go:396: frontmatter: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "a73b8668bef51370c7aa42e462ca7afaca153097",
          "requirement": "learnings",
          "command": "bench test --package ./internal/learnings --run '^TestParseReportsAnUnterminatedBlock$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/learnings: omit the malformed record for the reader fault; expect TestParseReportsAnUnterminatedBlock",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-repair-03",
              "digest": "sha256:2a5cf2c0e2d8f5a815bca77a558bf429b0546b6449c1978227d6cc54c4915cc4",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/learnings/learnings.go,swap,failed,3,TestParseReportsAnUnterminatedBlock/comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/learnings,^TestParseReportsAnUnterminatedBlock$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,782\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/learnings,fail,2,4\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/comment,\"learnings_test.go:396: comment: Parse gives 0 entries and 2 malformed records, want 0 and 3 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/fence,\"learnings_test.go:396: fence: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/frontmatter,\"learnings_test.go:396: frontmatter: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r3-retros",
          "performer": "ft358-repair-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9eceeda011572dab56efc7b5241150fa3518bbe0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-repair-03",
            "digest": "sha256:42754338134bf11b5e60cf7530840544ff9dd822a5ac608c0b45aef81bb60b14",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/retros/recommendations.go,swap,failed,1,TestRecommendationsSkipFencedLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/retros,^TestRecommendationsSkipFencedLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,402\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/retros,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/retros,TestRecommendationsSkipFencedLines,\"recommendations_test.go:36: fenced lines: Recommendations gives 1 units or other units, want the 3 units outside every fence and frontmatter block\",2\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "a73b8668bef51370c7aa42e462ca7afaca153097",
          "requirement": "retros",
          "command": "bench test --package ./internal/retros --run '^TestRecommendationsSkipFencedLines$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/retros/recommendations.go: end the improvements on any H2 line, fenced or not; expect TestRecommendationsSkipFencedLines",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-repair-03",
              "digest": "sha256:42754338134bf11b5e60cf7530840544ff9dd822a5ac608c0b45aef81bb60b14",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/retros/recommendations.go,swap,failed,1,TestRecommendationsSkipFencedLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/retros,^TestRecommendationsSkipFencedLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,402\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/retros,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/retros,TestRecommendationsSkipFencedLines,\"recommendations_test.go:36: fenced lines: Recommendations gives 1 units or other units, want the 3 units outside every fence and frontmatter block\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r4-learnings",
          "performer": "ft358-repair-03b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8a003b2caf0b92790af1f599434e4241878fb1d6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-repair-03b",
            "digest": "sha256:2a5cf2c0e2d8f5a815bca77a558bf429b0546b6449c1978227d6cc54c4915cc4",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/learnings/learnings.go,swap,failed,3,TestParseReportsAnUnterminatedBlock/comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/learnings,^TestParseReportsAnUnterminatedBlock$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,782\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/learnings,fail,2,4\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/comment,\"learnings_test.go:396: comment: Parse gives 0 entries and 2 malformed records, want 0 and 3 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/fence,\"learnings_test.go:396: fence: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/frontmatter,\"learnings_test.go:396: frontmatter: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "73bfbb2d56ac678a658b54d24561ba91636525c6",
          "requirement": "learnings",
          "command": "bench test --package ./internal/learnings --run '^TestParseReportsAnUnterminatedBlock$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/learnings: omit the malformed record for the reader fault; expect TestParseReportsAnUnterminatedBlock",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-repair-03b",
              "digest": "sha256:2a5cf2c0e2d8f5a815bca77a558bf429b0546b6449c1978227d6cc54c4915cc4",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/learnings/learnings.go,swap,failed,3,TestParseReportsAnUnterminatedBlock/comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/learnings,^TestParseReportsAnUnterminatedBlock$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,782\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/learnings,fail,2,4\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/comment,\"learnings_test.go:396: comment: Parse gives 0 entries and 2 malformed records, want 0 and 3 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/fence,\"learnings_test.go:396: fence: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/frontmatter,\"learnings_test.go:396: frontmatter: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r4-retros",
          "performer": "ft358-repair-03b",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8a003b2caf0b92790af1f599434e4241878fb1d6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-repair-03b",
            "digest": "sha256:42754338134bf11b5e60cf7530840544ff9dd822a5ac608c0b45aef81bb60b14",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/retros/recommendations.go,swap,failed,1,TestRecommendationsSkipFencedLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/retros,^TestRecommendationsSkipFencedLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,402\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/retros,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/retros,TestRecommendationsSkipFencedLines,\"recommendations_test.go:36: fenced lines: Recommendations gives 1 units or other units, want the 3 units outside every fence and frontmatter block\",2\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "73bfbb2d56ac678a658b54d24561ba91636525c6",
          "requirement": "retros",
          "command": "bench test --package ./internal/retros --run '^TestRecommendationsSkipFencedLines$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/retros/recommendations.go: end the improvements on any H2 line, fenced or not; expect TestRecommendationsSkipFencedLines",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-repair-03b",
              "digest": "sha256:42754338134bf11b5e60cf7530840544ff9dd822a5ac608c0b45aef81bb60b14",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/retros/recommendations.go,swap,failed,1,TestRecommendationsSkipFencedLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/retros,^TestRecommendationsSkipFencedLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,402\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/retros,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/retros,TestRecommendationsSkipFencedLines,\"recommendations_test.go:36: fenced lines: Recommendations gives 1 units or other units, want the 3 units outside every fence and frontmatter block\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r5-learnings",
          "performer": "ft358-repair-03c",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f6e4939b773f5cf85c949e9084f7a32f43f7b17",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-repair-03c",
            "digest": "sha256:2a5cf2c0e2d8f5a815bca77a558bf429b0546b6449c1978227d6cc54c4915cc4",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/learnings/learnings.go,swap,failed,3,TestParseReportsAnUnterminatedBlock/comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/learnings,^TestParseReportsAnUnterminatedBlock$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,782\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/learnings,fail,2,4\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/comment,\"learnings_test.go:396: comment: Parse gives 0 entries and 2 malformed records, want 0 and 3 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/fence,\"learnings_test.go:396: fence: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/frontmatter,\"learnings_test.go:396: frontmatter: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "requirement": "learnings",
          "command": "bench test --package ./internal/learnings --run '^TestParseReportsAnUnterminatedBlock$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/learnings: omit the malformed record for the reader fault; expect TestParseReportsAnUnterminatedBlock",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-repair-03c",
              "digest": "sha256:2a5cf2c0e2d8f5a815bca77a558bf429b0546b6449c1978227d6cc54c4915cc4",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/learnings/learnings.go,swap,failed,3,TestParseReportsAnUnterminatedBlock/comment,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/learnings,^TestParseReportsAnUnterminatedBlock$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,782\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/learnings,fail,2,4\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/comment,\"learnings_test.go:396: comment: Parse gives 0 entries and 2 malformed records, want 0 and 3 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/fence,\"learnings_test.go:396: fence: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\n  github.com/gibbonmi/bench/internal/learnings,TestParseReportsAnUnterminatedBlock/frontmatter,\"learnings_test.go:396: frontmatter: Parse gives 1 entries and 1 malformed records, want 1 and 2 in line order\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r5-retros",
          "performer": "ft358-repair-03c",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f6e4939b773f5cf85c949e9084f7a32f43f7b17",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-repair-03c",
            "digest": "sha256:42754338134bf11b5e60cf7530840544ff9dd822a5ac608c0b45aef81bb60b14",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/retros/recommendations.go,swap,failed,1,TestRecommendationsSkipFencedLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/retros,^TestRecommendationsSkipFencedLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,402\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/retros,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/retros,TestRecommendationsSkipFencedLines,\"recommendations_test.go:36: fenced lines: Recommendations gives 1 units or other units, want the 3 units outside every fence and frontmatter block\",2\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "requirement": "retros",
          "command": "bench test --package ./internal/retros --run '^TestRecommendationsSkipFencedLines$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/retros/recommendations.go: end the improvements on any H2 line, fenced or not; expect TestRecommendationsSkipFencedLines",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-repair-03c",
              "digest": "sha256:42754338134bf11b5e60cf7530840544ff9dd822a5ac608c0b45aef81bb60b14",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/retros/recommendations.go,swap,failed,1,TestRecommendationsSkipFencedLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/retros,^TestRecommendationsSkipFencedLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,402\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/retros,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/retros,TestRecommendationsSkipFencedLines,\"recommendations_test.go:36: fenced lines: Recommendations gives 1 units or other units, want the 3 units outside every fence and frontmatter block\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c3-standards",
          "performer": "ft358-review-c3-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "183a350962746062b7ffb6fc133d42339b3ace33",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c3-standards",
            "digest": "sha256:f0868fc565ae8f24977eb84f73117e9147b700bceab400c682e56d961bf69ed5",
            "excerpt": "Standards axis, MB-C3 (base 60f0d9ee, tip 251b6e3a). Six findings, all low.\nS1 (low, high, auto-fix): internal/handoffdoc/document.go:7 package doc says it imports the standard library and the block reader only; store.go imports internal/bounds and internal/durablefile.\nS2 (low, high, auto-fix): internal/handoff/state_file.go:72 calls markdown.ReadFragment and checks CommentFault itself, while handoffdoc/fence.go owns the State-body read; export one State-body fault accessor from handoffdoc.\nS3 (low, medium-high, auto-fix): internal/learnings/learnings.go:55, :67, :157 keep a parallel strings.Split and strip CR again; use the reader's Line.Raw.\nS4 (low, high, no-op this chunk): internal/handoff/sections.go:7-10 ShapeHeading duplicates handoffdoc shapeHeading; predates the delta and is outside every Writes line; capture as a learning.\nS5 (low, medium, ask-user): learnings.go:194, :209 and handoffdoc/document.go:270 detect a heading on Text but take the title from Raw; a heading behind a leading inline comment may fail the title parse.\nS6 (low, medium, auto-fix): new test doc comments lead with row IDs (learnings_test.go:358, :373, retros_test.go:60, recommendations_test.go:21, handoffdoc/document_test.go:196, :249), which craft-comments treats as provenance.\n"
          },
          "axis": "Standards",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "251b6e3adf4950681b0227de1fc93449dc8910ef",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "S5",
            "S6"
          ],
          "supersedes": []
        },
        {
          "id": "c3-spec",
          "performer": "ft358-review-c3-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "183a350962746062b7ffb6fc133d42339b3ace33",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c3-spec",
            "digest": "sha256:b9a4ffc40f0a48a45e3c1e7a8bba7237ae1608d64f950d2f2ef31599ae814ca4",
            "excerpt": "Spec axis, MB-C3 (base 60f0d9ee, tip 251b6e3a). Four findings; rows MB38-MB47, MB59, MB60, MB82, MB84, MB93 hold; no Writes file grows past its bound.\nP1 (medium, high, auto-fix): internal/handoff/state_file.go:78-79. The guard tests H2Prefix on trimmed UnfencedLines text, while Parse takes headings from the reader; a draft line of ## with trailing spaces passes the guard, Render writes it verbatim, and every later Parse refuses the document. Repair: refuse on line.Heading from the ReadFragment result at :72.\nP2 (low, high, ask-user): spec rule-delta row :737 and the flagged additions do not record that a handoff ## line with trailing whitespace, or with comment-free text that starts with ##, now opens a heading that Parse refuses. Add an FA line.\nP3 (low, medium, ask-user): handoffdoc fence.go:31-43 with document.go:209. An open fence followed by an unterminated comment lost its Parse fence refusal, the handoffdoc twin of FA11. Add an FA line.\nP4 (low, high, ask-user): internal/retros/retros.go:76 matches on line.Text, but ticket 6 and spec :254 say retros reads raw text; a commented H3 entry no longer counts and ## Outcome followed by a comment now counts.\n"
          },
          "axis": "Spec",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "251b6e3adf4950681b0227de1fc93449dc8910ef",
          "finding_ids": [
            "P1",
            "P2",
            "P3",
            "P4"
          ],
          "supersedes": []
        },
        {
          "id": "c3-coverage",
          "performer": "ft358-review-c3-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "183a350962746062b7ffb6fc133d42339b3ace33",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c3-coverage",
            "digest": "sha256:15b2464fd8a5e08bcb2c45c48f3756d05edc8ac5eda410fc08258b7fc90e8a6b",
            "excerpt": "Coverage axis, MB-C3 (base 60f0d9ee, tip 251b6e3a). Five findings. Every probe restored; git status clean.\nC1 (high, confidence 9, ask-user): internal/handoffdoc/document.go:209. Parse reads through Read; an unclosed first-line --- reports a frontmatter fault that masks the fence fault, so Parse accepts an open fence that the base refused (Closed decision 1).\nC2 (medium, confidence 8, ask-user): internal/handoffdoc/document.go:209. Parse does not refuse an unterminated comment, so a later section silently joins the State.\nC3 (medium, confidence 9, auto-fix): internal/handoff/state_file.go:72 ReadFragment->Read is silent. Add an MB82 subtest whose first line is ---.\nC4 (medium, confidence 9, auto-fix): FA2 (a heading inside a closed comment opens nothing) has no test at handoffdoc/document.go:262, learnings/learnings.go:113, retros/recommendations.go:25. Add one commented-heading case per site.\nC5 (low, confidence 7, ask-user): Read->ReadFragment is silent at retros/retros.go:75, retros/recommendations.go:20, learnings.go:277; a closed frontmatter block in a journal hides an entry with no record. Add a first-line --- row per site and decide the learnings frontmatter.\n"
          },
          "axis": "Coverage",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "251b6e3adf4950681b0227de1fc93449dc8910ef",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
            "C5"
          ],
          "supersedes": []
        },
        {
          "id": "c3r-standards",
          "performer": "ft358-confirm-c3-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9eceeda011572dab56efc7b5241150fa3518bbe0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c3-standards",
            "digest": "sha256:d4495394aaf491ad7d2d9ada29f3f34436d8b35a5bca297158c3f8caf06ed055",
            "excerpt": "Standards axis, MB-C3 confirming round (repair delta 251b6e3a..a73b8668). S1, S2, S3, S5, S6 closed; three new low findings.\nS7 (low, high, auto-fix): internal/handoffdoc/fence.go:30-32 OpenFence doc says a writer calls it before the render; no production code calls it now. Rewrite the doc to state what it reports.\nS8 (low, medium, auto-fix): fence.go:52-55 and internal/handoff/state_file.go:74-77 both map a comment fault to the comment repair and any other fault to the fence repair. Export one handoffdoc function that both call.\nS9 (low, high, auto-fix): internal/learnings/entry.go:48 sourceBody doc says it trims the empty lines at each end; it trims only newline bytes.\nReadState, OpenCommentRepair, and sourceBody are sound against one source per fact.\n"
          },
          "axis": "Standards",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "a73b8668bef51370c7aa42e462ca7afaca153097",
          "finding_ids": [
            "S7",
            "S8",
            "S9"
          ],
          "supersedes": [
            "c3-standards"
          ]
        },
        {
          "id": "c3r-spec",
          "performer": "ft358-confirm-c3-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9eceeda011572dab56efc7b5241150fa3518bbe0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c3-spec",
            "digest": "sha256:49b5a145908e9b720123006c977833666d7448d5997973309180af5fa4b97b0f",
            "excerpt": "Spec axis, MB-C3 confirming round (repair delta 251b6e3a..a73b8668). P1, P2, P3 closed; P4 confirmed; FA10, FA15 met; every file within its bound; two new low findings.\nP5 (low, confidence 7, auto-fix decided): internal/handoffdoc/document.go:271 takes block.heading from line.Text, so the existing message at document.go:356 quotes the comment-free heading, against spec.md:871 (existing messages keep their bytes). Repair: quote line.Raw in the message.\nP6 (low, confidence 6, auto-fix, spec text only): spec.md:255 marks internal/handoff comment-free, but the State pin scan reads raw text through UnfencedLines. Narrow the row to the heading grammar.\n"
          },
          "axis": "Spec",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "a73b8668bef51370c7aa42e462ca7afaca153097",
          "finding_ids": [
            "P5",
            "P6"
          ],
          "supersedes": [
            "c3-spec"
          ]
        },
        {
          "id": "c3r-coverage",
          "performer": "ft358-confirm-c3-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9eceeda011572dab56efc7b5241150fa3518bbe0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c3-coverage",
            "digest": "sha256:b8be6f3c0224ce9abefac80bf262f77a922e793ca6d7ac537026ac4fadfbb39a",
            "excerpt": "Coverage axis, MB-C3 confirming round (repair delta 251b6e3a..a73b8668). C1-C5 closed (each round 1 mutant now red); one new finding. Every probe restored; git status clean.\nNew mutants on the repaired code red: Text->Raw title reads, faultRefusal opener and comment branch, readStateFile guard, opener, and comment branch, sourceBody guard and trim.\nC6 (low, high confidence, auto-fix): internal/handoffdoc/fence.go:43 frontmatter-fault narrowing is untested; if true passes 1696 tests. Add a case with a closed first-line --- block that holds a fence marker line, which Parse must accept.\n"
          },
          "axis": "Coverage",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "a73b8668bef51370c7aa42e462ca7afaca153097",
          "finding_ids": [
            "C6"
          ],
          "supersedes": [
            "c3-coverage"
          ]
        },
        {
          "id": "c3r2-standards",
          "performer": "ft358-confirm2-c3-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8a003b2caf0b92790af1f599434e4241878fb1d6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm2-c3-standards",
            "digest": "sha256:fa0df6652ed8b3c0ff4f4b649849e4d5c27ba0652832b561d291d16bb95c7647",
            "excerpt": "Standards axis, MB-C3 confirming round 2 (repair delta a73b8668..73bfbb2d). Zero findings.\nS7 closed: internal/handoffdoc/fence.go:38-42 OpenFence doc states only what it reports. S8 closed: handoffdoc.FaultRepair (fence.go:24-29) is the one fault-to-repair map; faultRefusal and readStateFile both call it. S9 closed: entry.go:48 doc matches the code.\nOptional advice (no finding): stateRepair in internal/handoff/state_scan.go:34 still spells the rerun suffix itself; the legacy next-command refusal at document.go:320 quotes the comment-free heading; the OpenFence doc phrase about a comment coming first is ambiguous; bodyText's doc at fence.go:96 says empty lines where it trims newline bytes.\n"
          },
          "axis": "Standards",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "73bfbb2d56ac678a658b54d24561ba91636525c6",
          "finding_ids": [],
          "supersedes": [
            "c3r-standards"
          ]
        },
        {
          "id": "c3r2-spec",
          "performer": "ft358-confirm2-c3-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8a003b2caf0b92790af1f599434e4241878fb1d6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm2-c3-spec",
            "digest": "sha256:331b2c2e9283e5953c82c1c1a960077a3c9dbafd0e3b76445d33e2530795e89f",
            "excerpt": "Spec axis, MB-C3 confirming round 2 (repair delta a73b8668..73bfbb2d). P5 closed at its cited sites; P6 closed; one new finding. document.go is 399 lines.\nEvery base refusal keeps its bytes except the legacy-mix refusal.\nP7 (low, confidence 8, auto-fix): internal/handoffdoc/document.go:298 sets section = b.heading and :305 quotes it, so a legacy document with a heading followed by a trailing comment quotes the comment-free heading, against spec.md:871. Repair: section = b.raw at :298, and one legacy-mix row in TestParseQuotesTheRawHeading.\nOptional advice: :320 quotes b.heading; quoting b.raw would match the other heading refusals.\n"
          },
          "axis": "Spec",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "73bfbb2d56ac678a658b54d24561ba91636525c6",
          "finding_ids": [
            "P7"
          ],
          "supersedes": [
            "c3r-spec"
          ]
        },
        {
          "id": "c3r2-coverage",
          "performer": "ft358-confirm2-c3-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8a003b2caf0b92790af1f599434e4241878fb1d6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm2-c3-coverage",
            "digest": "sha256:75e58352af0ad65cfdf926a769d220c76f0bb417e275f4028d52c11948a4f11f",
            "excerpt": "Coverage axis, MB-C3 confirming round 2 (repair delta a73b8668..73bfbb2d). Zero findings. Every probe restored; git status clean.\nC6 closed: true in place of the frontmatter-fault narrowing at internal/handoffdoc/fence.go:51 fails the fence marker in closed frontmatter case at document_test.go:259. FaultRepair if false and the faultRefusal comment opener omission both fail tests.\nOptional advice (no finding): FaultRepair if true passes, because no test checks the fence repair text (a gap that predates the repair); the TrimRight on block.raw at document.go:271 has no test.\n"
          },
          "axis": "Coverage",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "73bfbb2d56ac678a658b54d24561ba91636525c6",
          "finding_ids": [],
          "supersedes": [
            "c3r-coverage"
          ]
        },
        {
          "id": "c3r3-spec",
          "performer": "ft358-confirm3-c3-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f6e4939b773f5cf85c949e9084f7a32f43f7b17",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm3-c3-spec",
            "digest": "sha256:046fde85e4e85ab80f3892fd695d46c02c59a6ed1ab8a6a7b58bfc23b1998324",
            "excerpt": "Spec axis, MB-C3 confirming round 3 (repair delta 73bfbb2d..82fe6fdb). Zero findings.\nP7 closed: every heading refusal in internal/handoffdoc/document.go quotes b.raw (:298 and :305, :320, :352, :356), the same bytes that the base quoted; :305 quotes the legacyStateHeading constant as the base did; matching stays on comment-free text (FA15); document.go is 399 lines; the two new rows reach :305 and :320.\nOptional advice: the test assertion is now Contains, which the legacy-mix row needs.\n"
          },
          "axis": "Spec",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "finding_ids": [],
          "supersedes": [
            "c3r2-spec"
          ]
        },
        {
          "id": "c3r3-standards",
          "performer": "ft358-confirm3-c3-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f6e4939b773f5cf85c949e9084f7a32f43f7b17",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm3-c3-standards",
            "digest": "sha256:0d26c69e57a5ae0e53b5e6ad7d977a76b1b5744722e3325eee415cd5ac2cfa94",
            "excerpt": "Standards axis, MB-C3 confirming round 3 (repair delta 73bfbb2d..82fe6fdb), native reaffirmation. Zero findings.\nThe two refusals quote b.raw (internal/handoffdoc/document.go:298, :320), the field that splitSections sets for refusals; the test builds its headings from the package's own constants; the test doc comment states the current state.\nOptional advice: the assertion is now Contains in place of the heading prefix check; :305 still quotes the legacyStateHeading constant, as at the base.\n"
          },
          "axis": "Standards",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "finding_ids": [],
          "supersedes": [
            "c3r2-standards"
          ]
        },
        {
          "id": "c3r3-coverage",
          "performer": "ft358-confirm3-c3-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0f6e4939b773f5cf85c949e9084f7a32f43f7b17",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm3-c3-coverage",
            "digest": "sha256:04535b54fbdbc9c732ad6a6a4375d7e8e8a3cfb71c6d53e023efebe0c9e1103d",
            "excerpt": "Coverage axis, MB-C3 confirming round 3 (repair delta 73bfbb2d..82fe6fdb), native reaffirmation. Zero findings. Every probe restored; git status clean.\nb.raw back to b.heading at internal/handoffdoc/document.go:298 fails only TestParseQuotesTheRawHeading/legacy_mix; the same swap at :320 fails only TestParseQuotesTheRawHeading/legacy_next_command.\n"
          },
          "axis": "Coverage",
          "base": "60f0d9ee1fe632cc0ae916661e01dc2a9c869c80",
          "tip": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "finding_ids": [],
          "supersedes": [
            "c3r2-coverage"
          ]
        }
      ]
    },
    {
      "id": "MB-C4",
      "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
      "tip": "15903e2ce3aab9a0db82da0678bf3a7359ec8d58",
      "plan_digest": "sha256:6c72c6b28d025c0c98a527879153466994ffe1c855fef1ae3dd055075067a043",
      "source_digest": "f39a6cfb42f4d277a91ca1bf629d1acc9bce4d9c",
      "acceptance_rows": [
        "MB48",
        "MB49",
        "MB50",
        "MB51",
        "MB96",
        "MB107",
        "MB52",
        "MB53",
        "MB55",
        "MB56",
        "MB85",
        "MB94",
        "MB108",
        "MB109"
      ],
      "verification": [
        {
          "id": "t7-reviewrecord",
          "performer": "ft358-author-07",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "706396690a176f0a092ee846506269477f83eeb3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-07",
            "digest": "sha256:be991a7d4dc62769a68fd5d82d3ed44c752f0acf215cc76625cb675db3aed865",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/parse.go,swap,failed,1,TestRecordSkipsAQuotedFence,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,^TestRecordSkipsAQuotedFence$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,440\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestRecordSkipsAQuotedFence,\"record_test.go:193: a longer fence that quotes the record opener: render and read = invalid duplicate bench-review-record fence; want the record payload replaced and read in place\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8bc499d87ee64737ed23bd7013abb2b1260be5ae",
          "requirement": "reviewrecord",
          "command": "bench test --package ./internal/reviewrecord --run '^TestRecordSkipsAQuotedFence$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/reviewrecord/parse.go: match a record opener on fenced lines too; expect TestRecordSkipsAQuotedFence",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-07",
              "digest": "sha256:be991a7d4dc62769a68fd5d82d3ed44c752f0acf215cc76625cb675db3aed865",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/parse.go,swap,failed,1,TestRecordSkipsAQuotedFence,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,^TestRecordSkipsAQuotedFence$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,440\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestRecordSkipsAQuotedFence,\"record_test.go:193: a longer fence that quotes the record opener: render and read = invalid duplicate bench-review-record fence; want the record payload replaced and read in place\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t8-anchors",
          "performer": "ft358-author-08",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f4c2bdff3d46f2490d899712e4af6a1dea3b3113",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-author-08",
            "digest": "sha256:95a7b1bbdfa29bdfa0e18ccbd0e843032376d4af4f74d6d304293620f8c71fc7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap+swap,failed,2,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,^TestLocateReadsTheBlockReader$,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,588\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,2,10\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,\"locate_test.go:293: a tilde-fenced heading counts as no section: section count = 2, want 1\",2\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/an_indented_heading_opens_no_section,\"locate_test.go:293: an indented heading opens no section: section count = 1, want 0\",3\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ea07f8994bc0dc3802829aace137e7c5288a1b96",
          "requirement": "anchors",
          "command": "bench test --package ./internal/anchors --run '^TestLocateReadsTheBlockReader$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/anchors/locate.go: open a section on a trimmed heading line; expect TestLocateReadsTheBlockReader",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-author-08",
              "digest": "sha256:95a7b1bbdfa29bdfa0e18ccbd0e843032376d4af4f74d6d304293620f8c71fc7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap+swap,failed,2,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,^TestLocateReadsTheBlockReader$,passed,10\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,588\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,2,10\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,\"locate_test.go:293: a tilde-fenced heading counts as no section: section count = 2, want 1\",2\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/an_indented_heading_opens_no_section,\"locate_test.go:293: an indented heading opens no section: section count = 1, want 0\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t8-docs-currency-workflow",
          "performer": "ft358-author-08",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f4c2bdff3d46f2490d899712e4af6a1dea3b3113",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-author-08",
            "digest": "sha256:8b762d3326ef5bd56aa175716e476851f86f11d44b6f6b2c9671a01622e69cb0",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,476\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3397,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: .agents/skills/bench-craft-spec/SKILL.md Template missing acceptance coverage anchor: why it catches the failure\",9\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "ea07f8994bc0dc3802829aace137e7c5288a1b96",
          "requirement": "docs-currency-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/anchors/locate.go: end a section at any line that starts with `#`; expect TestRootConformance",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-author-08",
              "digest": "sha256:8b762d3326ef5bd56aa175716e476851f86f11d44b6f6b2c9671a01622e69cb0",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,476\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3397,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: .agents/skills/bench-craft-spec/SKILL.md Template missing acceptance coverage anchor: why it catches the failure\",9\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r6-anchors",
          "performer": "ft358-repair-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "46a61529456dfc592f5e2c20c015713aaa8ac2fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04",
            "digest": "sha256:cb0fffe7da3b297f6517717df8b494520634ba06dbd024369dbfa58b41a0c7fa",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap+swap,failed,2,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,^TestLocateReadsTheBlockReader$,passed,11\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,588\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,2,11\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,\"locate_test.go:295: a tilde-fenced heading counts as no section: section count = 2, want 1\",2\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/an_indented_heading_opens_no_section,\"locate_test.go:295: an indented heading opens no section: section count = 1, want 0\",3\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "b5f0b2de83f9885e2d526e1d79b6529e65b6900d",
          "requirement": "anchors",
          "command": "bench test --package ./internal/anchors --run '^TestLocateReadsTheBlockReader$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/anchors/locate.go: open a section on a trimmed heading line; expect TestLocateReadsTheBlockReader",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04",
              "digest": "sha256:cb0fffe7da3b297f6517717df8b494520634ba06dbd024369dbfa58b41a0c7fa",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap+swap,failed,2,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,^TestLocateReadsTheBlockReader$,passed,11\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,588\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,2,11\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,\"locate_test.go:295: a tilde-fenced heading counts as no section: section count = 2, want 1\",2\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/an_indented_heading_opens_no_section,\"locate_test.go:295: an indented heading opens no section: section count = 1, want 0\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r6-docs-currency-workflow",
          "performer": "ft358-repair-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "46a61529456dfc592f5e2c20c015713aaa8ac2fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04",
            "digest": "sha256:ee1edda5bb31f3bafd49b8f48d19b539fb94a4b88cb6c873aee99b428bc7fa34",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,476\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3378,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: .agents/skills/bench-craft-spec/SKILL.md Template missing acceptance coverage anchor: why it catches the failure\",9\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "b5f0b2de83f9885e2d526e1d79b6529e65b6900d",
          "requirement": "docs-currency-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/anchors/locate.go: end a section at any line that starts with `#`; expect TestRootConformance",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04",
              "digest": "sha256:ee1edda5bb31f3bafd49b8f48d19b539fb94a4b88cb6c873aee99b428bc7fa34",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,476\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3378,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: .agents/skills/bench-craft-spec/SKILL.md Template missing acceptance coverage anchor: why it catches the failure\",9\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r7-anchors",
          "performer": "ft358-repair-04b",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "3758c40164a6d4adc2c8d4f2c6178106b4e31cd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04b",
            "digest": "sha256:653b50888e244cabec7979a82a27bd634544e2238170e92e5b7f2348408e886d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap+swap,failed,2,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,^TestLocateReadsTheBlockReader$,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,588\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,3,16\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,\"locate_test.go:301: a tilde-fenced heading counts as no section: section count = 2, want 1\",3\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/an_indented_heading_opens_no_section,\"locate_test.go:301: an indented heading opens no section: section count = 1, want 0\",4\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f04f636e8e8f38beb3159e4984d8b507df83bf50",
          "requirement": "anchors",
          "command": "bench test --package ./internal/anchors --run '^TestLocateReadsTheBlockReader$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/anchors/locate.go: open a section on a trimmed heading line; expect TestLocateReadsTheBlockReader",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04b",
              "digest": "sha256:653b50888e244cabec7979a82a27bd634544e2238170e92e5b7f2348408e886d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap+swap,failed,2,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,^TestLocateReadsTheBlockReader$,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,588\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,3,16\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,\"locate_test.go:301: a tilde-fenced heading counts as no section: section count = 2, want 1\",3\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/an_indented_heading_opens_no_section,\"locate_test.go:301: an indented heading opens no section: section count = 1, want 0\",4\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r7-docs-currency-workflow",
          "performer": "ft358-repair-04b",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "3758c40164a6d4adc2c8d4f2c6178106b4e31cd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04b",
            "digest": "sha256:eea70f2cc14f6c07fbe7d6f9b2cc6f8304a6dff71c35920570a7578e895706d0",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,476\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3422,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: .agents/skills/bench-craft-spec/SKILL.md Template missing acceptance coverage anchor: why it catches the failure\",9\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f04f636e8e8f38beb3159e4984d8b507df83bf50",
          "requirement": "docs-currency-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/anchors/locate.go: end a section at any line that starts with `#`; expect TestRootConformance",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04b",
              "digest": "sha256:eea70f2cc14f6c07fbe7d6f9b2cc6f8304a6dff71c35920570a7578e895706d0",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,476\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3422,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: .agents/skills/bench-craft-spec/SKILL.md Template missing acceptance coverage anchor: why it catches the failure\",9\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r8-anchors",
          "performer": "ft358-repair-04c",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f39a6cfb42f4d277a91ca1bf629d1acc9bce4d9c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04c",
            "digest": "sha256:6fc1da5ec3aa540947e8d6ae67f62036b1ed8758af6fc85b03097f8b6e7510f5",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap+swap,failed,2,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,^TestLocateReadsTheBlockReader$,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,588\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,2,16\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,\"locate_test.go:309: a tilde-fenced heading counts as no section: section count = 2, want 1\",3\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/an_indented_heading_opens_no_section,\"locate_test.go:309: an indented heading opens no section: section count = 1, want 0\",4\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "15903e2ce3aab9a0db82da0678bf3a7359ec8d58",
          "requirement": "anchors",
          "command": "bench test --package ./internal/anchors --run '^TestLocateReadsTheBlockReader$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/anchors/locate.go: open a section on a trimmed heading line; expect TestLocateReadsTheBlockReader",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04c",
              "digest": "sha256:6fc1da5ec3aa540947e8d6ae67f62036b1ed8758af6fc85b03097f8b6e7510f5",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap+swap,failed,2,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/anchors,^TestLocateReadsTheBlockReader$,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,588\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,fail,2,16\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/a_tilde-fenced_heading_counts_as_no_section,\"locate_test.go:309: a tilde-fenced heading counts as no section: section count = 2, want 1\",3\n  github.com/gibbonmi/bench/internal/anchors,TestLocateReadsTheBlockReader/an_indented_heading_opens_no_section,\"locate_test.go:309: an indented heading opens no section: section count = 1, want 0\",4\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "r8-docs-currency-workflow",
          "performer": "ft358-repair-04c",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "f39a6cfb42f4d277a91ca1bf629d1acc9bce4d9c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04c",
            "digest": "sha256:9ae09b0cd9216381e45e585ca6e7405037876326e51a6d81b294c4c351931742",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,476\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3373,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: .agents/skills/bench-craft-spec/SKILL.md Template missing acceptance coverage anchor: why it catches the failure\",9\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "15903e2ce3aab9a0db82da0678bf3a7359ec8d58",
          "requirement": "docs-currency-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/anchors/locate.go: end a section at any line that starts with `#`; expect TestRootConformance",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft358-repair-04c",
              "digest": "sha256:9ae09b0cd9216381e45e585ca6e7405037876326e51a6d81b294c4c351931742",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/anchors/locate.go,swap,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,476\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3373,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: .agents/skills/bench-craft-spec/SKILL.md Template missing acceptance coverage anchor: why it catches the failure\",9\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c4-standards",
          "performer": "ft358-review-c4-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f4c2bdff3d46f2490d899712e4af6a1dea3b3113",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c4-standards",
            "digest": "sha256:71b550334ba4b621b40e4af5873b4f2d02107af3d99ce063efa480cbdd6d22e5",
            "excerpt": "Standards axis, MB-C4 (base 82fe6fdb, tip ea07f899). Two findings.\nS1 (medium, high, auto-fix): internal/anchors/registry.go:217 restates the comment rule in prose (removes complete comments and truncates at an unterminated comment), which is now incomplete because a comment inside valid frontmatter stays. Repair: say the function returns the text without the comments that the block reader removes; widen ticket 8 to registry.go.\nS2 (low, high, auto-fix): internal/anchors/locate.go:72-73 says the walk reads a section body; every caller passes the whole comment-stripped document.\nThe payload escape in internal/reviewrecord/write.go:46-58 is the one encoding site, and its comment is accurate.\nOptional advice: write.go:49 frames the escape as excerpts only; recordtest/fixture.go:86 writes the H2 heading as a literal; the strip reads frontmatter through Read while the walk reads through ReadFragment.\n"
          },
          "axis": "Standards",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "ea07f8994bc0dc3802829aace137e7c5288a1b96",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c4-spec",
          "performer": "ft358-review-c4-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f4c2bdff3d46f2490d899712e4af6a1dea3b3113",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c4-spec",
            "digest": "sha256:dd22ad0a59a81368bbb6b15373a2ff6b67df63433d7c1d78148f7830abdbaf31",
            "excerpt": "Spec axis, MB-C4 (base 82fe6fdb, tip ea07f899). One finding; all 12 rows mapped with tests that match; every line bound and refusal byte holds.\nP1 (medium, confidence 7, auto-fix decided): internal/reviewrecord/parse.go:294-323. A second record opener directly after a record reads as a closer with trailing text, so the base duplicate refusal no longer fires and Render writes a document that ends in an open fence; story 48 keeps the duplicate refusal. Repair: refuse a closer whose label equals the record name with the existing duplicate message, and add a nested case to record_test.go.\nFlags: (a) permitted; (b) required by spec:300-301 and rule 2; (c) required by spec:231-232; (d) permitted.\nOptional advice: the rule-delta rows could name the (b), (c), and closer deltas; ticket 7 and write.go say excerpt where the code escapes every less-than sign.\n"
          },
          "axis": "Spec",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "ea07f8994bc0dc3802829aace137e7c5288a1b96",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "c4-coverage",
          "performer": "ft358-review-c4-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f4c2bdff3d46f2490d899712e4af6a1dea3b3113",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c4-coverage",
            "digest": "sha256:863538561f36b52375b67c141713e006f2265b011ec95492a08821d2ac437e3c",
            "excerpt": "Coverage axis, MB-C4 (base 82fe6fdb, tip ea07f899). One finding; the other hostile-input classes hold (CRLF, no trailing newline, tilde and indented openers, the less-than escape, non-ASCII and invalid UTF-8 mapping, a heading behind a comment). git status clean.\nC1 (medium, confidence 7, ask-user): internal/anchors/locate.go:76 scopeRunesMapped reads the whole comment-stripped document through ReadFragment, while stripCommentsMapped at :52 reads it through Read. On a document with valid frontmatter that holds a comment opener or a tilde marker, the walk loses section S (count 0; the base gave 1 for the tilde case). The walk reads the whole file, not a section body, so spec 298-301, Closed decision 16, and MB94 rest on a wrong premise. Repair: read through Read like the strip, and restate MB94.\nOptional advice: FA7 has no test (the exact-line opener mutant passes 444 tests); raw greater-than is not pinned; a closer with trailing text now closes a record.\n"
          },
          "axis": "Coverage",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "ea07f8994bc0dc3802829aace137e7c5288a1b96",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "c4r-standards",
          "performer": "ft358-confirm-c4-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "46a61529456dfc592f5e2c20c015713aaa8ac2fa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c4-standards",
            "digest": "sha256:17abd62fa80d9efbf0008cc9b6d16d6cecc6a4d781e103a0060225a78a461eab",
            "excerpt": "Standards axis, MB-C4 confirming round (repair delta ea07f899..b5f0b2de). Zero findings.\nS1 closed: internal/anchors/registry.go:217 says the function returns the text without the comments that the block reader removes. S2 closed: internal/anchors/locate.go:72-73 says the walk reads the whole comment-free document.\nThe shared duplicate-refusal variable at internal/reviewrecord/parse.go:299 is one source for both refusal paths, and the new test derives its expected message from a two-record Render. Every doc comment above a changed body is accurate.\nOptional advice: the ReadFragment doc at internal/markdown/markdown.go:96 still names a section body; a duplicateFence type would avoid building the error up front.\n"
          },
          "axis": "Standards",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "b5f0b2de83f9885e2d526e1d79b6529e65b6900d",
          "finding_ids": [],
          "supersedes": [
            "c4-standards"
          ]
        },
        {
          "id": "c4r-spec",
          "performer": "ft358-confirm-c4-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "46a61529456dfc592f5e2c20c015713aaa8ac2fa",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c4-spec",
            "digest": "sha256:1d41463a74344c77467a091fc0d40bf1e57fbb44ddebbe2fa2abb57866f59bea",
            "excerpt": "Spec axis, MB-C4 confirming round (repair delta ea07f899..b5f0b2de). P1 closed; two new findings.\nP1 closed: internal/reviewrecord/parse.go:312-318 refuses the record's first closer with Info equal to the name with the existing duplicate message; every refusal keeps its base bytes; story 84, MB108, MB94, FA16, FA17, the fragment-caller paragraph, and Closed decision 16 match the code; every file is within its bound. The author's reading of MB108 (the record's own closer) is correct.\nP2 (medium, high, auto-fix): specs/markdown-block-reader/tickets/8-read-anchor-blocks.md:13-14 still says the walk uses ReadFragment, and line 42 keeps the old MB94 text.\nP3 (low, medium, decided auto-fix with Coverage C2): internal/anchors/locate.go:76 applies the frontmatter rule to the comment-free text, so a comment followed by --- on line 1 hides section S.\n"
          },
          "axis": "Spec",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "b5f0b2de83f9885e2d526e1d79b6529e65b6900d",
          "finding_ids": [
            "P2",
            "P3"
          ],
          "supersedes": [
            "c4-spec"
          ]
        },
        {
          "id": "c4r-coverage",
          "performer": "ft358-confirm-c4-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "46a61529456dfc592f5e2c20c015713aaa8ac2fa",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm-c4-coverage",
            "digest": "sha256:41806e45cedf671cb37d15676a61017751548b53b31c45c56cfcb1863979338a",
            "excerpt": "Coverage axis, MB-C4 confirming round (repair delta ea07f899..b5f0b2de). C1 closed; one new finding. Every probe restored; git status clean.\nC1 closed: ReadFragment at internal/anchors/locate.go:76 fails both new rows with section count = 0, want 1.\nC2 (medium, high, auto-fix): internal/anchors/locate.go:76 (and the evaluator path through registry.go:185) applies the line-1 frontmatter rule to the comment-free text. For a first line of a comment followed by --- on the same line, Locate returns 0 where the original document gives line 3; the base walk found it. Repair: take the frontmatter state from the original document, and add that row as the witness.\nThe nested-closer refusal and the shared duplicate variable have no silent mutant.\n"
          },
          "axis": "Coverage",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "b5f0b2de83f9885e2d526e1d79b6529e65b6900d",
          "finding_ids": [
            "C2"
          ],
          "supersedes": [
            "c4-coverage"
          ]
        },
        {
          "id": "c4r2-standards",
          "performer": "ft358-confirm2-c4-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "3758c40164a6d4adc2c8d4f2c6178106b4e31cd0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm2-c4-standards",
            "digest": "sha256:d204b77e54537e61e27ae02110203c610826452bc8f79e23be4933a3dc2105cc",
            "excerpt": "Standards axis, MB-C4 confirming round 2 (repair delta b5f0b2de..f04f636e). One finding.\nS1 (low, high, broken mandatory standard craft-comments, auto-fix): internal/anchors/match.go:50-53 says MarkdownH2Sections projects the walk so the evaluator and Locate read the same body; resolveSection and Locate now read through commentFreeSectionMapped. Repair: name commentFreeSectionMapped as their section reader, or drop the clause.\nOptional advice: commentFreeSectionMapped calls markdown.Read twice on the same data; the origin parameter of sectionRunesMapped is now always identityOrigin.\n"
          },
          "axis": "Standards",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "f04f636e8e8f38beb3159e4984d8b507df83bf50",
          "finding_ids": [
            "S1"
          ],
          "supersedes": [
            "c4r-standards"
          ]
        },
        {
          "id": "c4r2-spec",
          "performer": "ft358-confirm2-c4-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "3758c40164a6d4adc2c8d4f2c6178106b4e31cd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm2-c4-spec",
            "digest": "sha256:a3826ddc699f6b93f271bccd74527d69e3d3d128cbb9fed8b9f94b7e8fc7ac39",
            "excerpt": "Spec axis, MB-C4 confirming round 2 (repair delta b5f0b2de..f04f636e). Zero findings.\nP2 closed: ticket 8 says the file read and the section walk both use Read, with the frontmatter state taken from the original document; its acceptance line names an existing row. P3 closed: story 85 and MB109 match the code and the row at locate_test.go:294; the walk prose, Closed decision 16, and the sweep line match commentFreeSectionMapped and its two callers; bench coverage --check passes with 106 rows. MB94, MB52, MB53, MB85, MB55 still hold. The three author flags are accepted.\nOptional advice: the match.go doc is stale; docs_workflow_helpers_test.go:96 passes comment-free text to MarkdownH2Section, harmless today.\n"
          },
          "axis": "Spec",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "f04f636e8e8f38beb3159e4984d8b507df83bf50",
          "finding_ids": [],
          "supersedes": [
            "c4r-spec"
          ]
        },
        {
          "id": "c4r2-coverage",
          "performer": "ft358-confirm2-c4-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "3758c40164a6d4adc2c8d4f2c6178106b4e31cd0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm2-c4-coverage",
            "digest": "sha256:23a489b2ebaf1259c5fd1dc651cb40eb7f26f5d2ea255161998dad815892bf29",
            "excerpt": "Coverage axis, MB-C4 confirming round 2 (repair delta b5f0b2de..f04f636e). C2 closed; one new finding. Every probe restored; git status clean.\nC2 closed: the old walk in commentFreeSectionMapped fails the MB109 row and the closer row; feeding resolveSection the comment-free text fails the evaluator check. The to-search, from, to, and empty-range mutants all fail; ten hostile-input rows (CRLF, non-ASCII, spanning comments, end of file) all pass.\nC3 (medium, high, auto-fix): internal/anchors/registry.go:187 swapping commentFreeSectionMapped for sectionRunesMapped over the raw text is silent across 140 anchors tests and docs-currency-workflow, because the comment-hiding rows assert only Locate. Repair: add the row \"a comment in a section hides its needle from the evaluator\" (## S, a commented needle) to the first table, which also checks the evaluator.\n"
          },
          "axis": "Coverage",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "f04f636e8e8f38beb3159e4984d8b507df83bf50",
          "finding_ids": [
            "C3"
          ],
          "supersedes": [
            "c4r-coverage"
          ]
        },
        {
          "id": "c4r3-standards",
          "performer": "ft358-confirm3-c4-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f39a6cfb42f4d277a91ca1bf629d1acc9bce4d9c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm3-c4-standards",
            "digest": "sha256:6bbbbc4a4fce2b738eb5180d3ec33c9dd26317d5bf4a7e5ddc974a3b34099a0c",
            "excerpt": "Standards axis, MB-C4 confirming round 3 (repair delta f04f636e..15903e2c), reaffirmation. Zero findings.\nS1 closed: internal/anchors/match.go:50-52 says the body comes as stored with comments included, and names commentFreeSectionMapped as the reader of Locate and the evaluator, which matches locate.go:28 and registry.go:187. The evaluatorBody closure (locate_test.go:284-291) puts the write-and-resolve steps in one place.\nOptional advice: match.go:52 repeats the consumer list that the commentFreeSectionMapped doc states.\n"
          },
          "axis": "Standards",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "15903e2ce3aab9a0db82da0678bf3a7359ec8d58",
          "finding_ids": [],
          "supersedes": [
            "c4r2-standards"
          ]
        },
        {
          "id": "c4r3-spec",
          "performer": "ft358-confirm3-c4-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f39a6cfb42f4d277a91ca1bf629d1acc9bce4d9c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm3-c4-spec",
            "digest": "sha256:4c5e24da5393e141d7f9d4fbe6d307553b02c14417ae8b4e895c555fc78f5404",
            "excerpt": "Spec axis, MB-C4 confirming round 3 (repair delta f04f636e..15903e2c), reaffirmation. Zero findings.\nS1 closed: the match.go doc matches the walk prose (spec.md:302-303), Closed decision 16, and the sweep line. C3 closed: locate_test.go:339-344 checks the evaluator body on the section-scoped rows. The delta changes no behavior; story 85, MB109, MB52, MB53, and MB94 still hold.\n"
          },
          "axis": "Spec",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "15903e2ce3aab9a0db82da0678bf3a7359ec8d58",
          "finding_ids": [],
          "supersedes": [
            "c4r2-spec"
          ]
        },
        {
          "id": "c4r3-coverage",
          "performer": "ft358-confirm3-c4-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f39a6cfb42f4d277a91ca1bf629d1acc9bce4d9c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-confirm3-c4-coverage",
            "digest": "sha256:16d736b2251187e842fd3552928f8d33bb211c1bec6c5b98deb89a0044b14245",
            "excerpt": "Coverage axis, MB-C4 confirming round 3 (repair delta f04f636e..15903e2c), reaffirmation. Zero findings. The probe restored; git status clean.\nC3 closed: the raw sectionRunesMapped walk at internal/anchors/registry.go:187 fails only TestLocateReadsTheBlockReader/a_comment_in_a_section_hides_its_needle at locate_test.go:343. The first table keeps every evaluator assertion through the evaluatorBody closure, and the second table checks the evaluator for each section-scoped row.\n"
          },
          "axis": "Coverage",
          "base": "82fe6fdb187d5b59833c8d5bdee55c94492fe8d2",
          "tip": "15903e2ce3aab9a0db82da0678bf3a7359ec8d58",
          "finding_ids": [],
          "supersedes": [
            "c4r2-coverage"
          ]
        }
      ]
    },
    {
      "id": "MB-C5",
      "base": "15903e2ce3aab9a0db82da0678bf3a7359ec8d58",
      "tip": "9d5e681d32acc80cc0654018f783b1e23cb3a454",
      "plan_digest": "sha256:a7ef5bc715c0df86eb39bf91e71bf940aeedc093cbfa00369948467107736c0d",
      "source_digest": "2c8a1051010604cd1d2c0ef9a22f1be66f667273",
      "acceptance_rows": [
        "MB57",
        "MB58",
        "MB61",
        "MB70",
        "MB86",
        "MB87",
        "MB91",
        "MB97",
        "MB62",
        "MB63",
        "MB64",
        "MB65",
        "MB66",
        "MB67",
        "MB68",
        "MB69",
        "MB98",
        "MB101"
      ],
      "verification": [
        {
          "id": "t10-adopt",
          "performer": "ft358-author-10",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3343f290618c7d56f3650ad7e26e04dd2288e565",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-10",
            "digest": "sha256:8720a3bcc4b3aaa48c9c31dc5438b29fff9b3f67396393517a8dcc4f47c2013b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/adopt/marker.go,swap,failed,1,TestMarkersSkipTildeFences,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestMarkersSkipTildeFences$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,378\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestMarkersSkipTildeFences,\"link_plan_test.go:76: MB57 rewrite keeps a tilde-fenced example: err = conflict: AGENTS.md has malformed Bench managed block markers\",3\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "fe458ea701f1b6e1a32919ff8819a4f4eae8debd",
          "requirement": "adopt",
          "command": "bench test --package ./internal/adopt --run '^TestMarkersSkipTildeFences$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/adopt/marker.go: toggle the fence state on backtick marker lines only; expect TestMarkersSkipTildeFences",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-10",
              "digest": "sha256:8720a3bcc4b3aaa48c9c31dc5438b29fff9b3f67396393517a8dcc4f47c2013b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/adopt/marker.go,swap,failed,1,TestMarkersSkipTildeFences,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestMarkersSkipTildeFences$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,378\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestMarkersSkipTildeFences,\"link_plan_test.go:76: MB57 rewrite keeps a tilde-fenced example: err = conflict: AGENTS.md has malformed Bench managed block markers\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t11-markdown-block-owner",
          "performer": "ft358-author-11",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c8a1051010604cd1d2c0ef9a22f1be66f667273",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-11",
            "digest": "sha256:4333f1b0e4418a90a9d0e3a1ceed0e0ea45444fee6e7a5fefc1b22720d54a217",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/markdown_block_owner_test.go,omit,failed,1,TestMarkdownBlockOwnerBites/each_kind_in_internal,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestMarkdownBlockOwnerBites$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,540\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10,6\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestMarkdownBlockOwnerBites/each_kind_in_internal,\"markdown_block_owner_test.go:111: diagnostics = [\\\"internal/example/example.go:4 spells the backtick fence run outside internal/markdown; read or render the block through the block reader\\\" \\\"internal/example/example.go:6 spells the frontmatte… (761 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "9d5e681d32acc80cc0654018f783b1e23cb3a454",
          "requirement": "markdown-block-owner",
          "command": "bench test --package ./internal/conformance --run '^TestMarkdownBlockOwnerBites$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/conformance/markdown_block_owner_test.go: omit the tilde-run kind; expect TestMarkdownBlockOwnerBites",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-11",
              "digest": "sha256:4333f1b0e4418a90a9d0e3a1ceed0e0ea45444fee6e7a5fefc1b22720d54a217",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/markdown_block_owner_test.go,omit,failed,1,TestMarkdownBlockOwnerBites/each_kind_in_internal,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestMarkdownBlockOwnerBites$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,540\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10,6\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestMarkdownBlockOwnerBites/each_kind_in_internal,\"markdown_block_owner_test.go:111: diagnostics = [\\\"internal/example/example.go:4 spells the backtick fence run outside internal/markdown; read or render the block through the block reader\\\" \\\"internal/example/example.go:6 spells the frontmatte… (761 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t11-skill-description-budgets",
          "performer": "ft358-author-11",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c8a1051010604cd1d2c0ef9a22f1be66f667273",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-11",
            "digest": "sha256:8bab6fbc5feaa8f70cd65b7fcf66efd2268943e9107d616b0a10ddb869a2f9c8",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/skill_description_budgets_test.go,insert,failed,1,TestSkillDescriptionBudgetReportsMissingAndFoldedValues,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestSkillDescriptionBudgetReportsMissingAndFoldedValues$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,465\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,11,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestSkillDescriptionBudgetReportsMissingAndFoldedValues,\"skill_description_budgets_test.go:431: want \\\"skill-description folded: .agents/commands/bench-fixture.md continues its description on an indented second line\\\" in one run, got:\",2\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "9d5e681d32acc80cc0654018f783b1e23cb3a454",
          "requirement": "skill-description-budgets",
          "command": "bench test --package ./internal/conformance --run '^TestSkillDescriptionBudgetReportsMissingAndFoldedValues$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/conformance/skill_description_budgets_test.go: make skillDescriptionFolded return false; expect TestSkillDescriptionBudgetReportsMissingAndFoldedValues",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-11",
              "digest": "sha256:8bab6fbc5feaa8f70cd65b7fcf66efd2268943e9107d616b0a10ddb869a2f9c8",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/skill_description_budgets_test.go,insert,failed,1,TestSkillDescriptionBudgetReportsMissingAndFoldedValues,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestSkillDescriptionBudgetReportsMissingAndFoldedValues$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,465\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,11,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestSkillDescriptionBudgetReportsMissingAndFoldedValues,\"skill_description_budgets_test.go:431: want \\\"skill-description folded: .agents/commands/bench-fixture.md continues its description on an indented second line\\\" in one run, got:\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t11-skills-index-command-adapters",
          "performer": "ft358-author-11",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c8a1051010604cd1d2c0ef9a22f1be66f667273",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-11",
            "digest": "sha256:952a78f224326566e450cb63319fb566b6c4835b3491ece3f7932a91fe3e4b1f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/skills_index_checks_test.go,omit,failed,1,TestCommandInvocationPolicyGradesTableCompleteness/body_prose_mention_of_the_key_is_inert,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestCommandInvocationPolicyGradesTableCompleteness$,passed,5\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,563\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,14,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestCommandInvocationPolicyGradesTableCompleteness/body_prose_mention_of_the_key_is_inert,\"skills_index_checks_test.go:537: body-prose mention of disable-model-invocation reddened the check: [\\\"command 'bench-write-spec' frontmatter carries disable-model-invocation though the invocation policy makes it model-invocable on Claude\\\"]\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "9d5e681d32acc80cc0654018f783b1e23cb3a454",
          "requirement": "skills-index-command-adapters",
          "command": "bench test --package ./internal/conformance --run '^TestCommandInvocationPolicyGradesTableCompleteness$'",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/conformance/skills_index_checks_test.go: make frontmatterHasKey read past the closing frontmatter marker; expect TestCommandInvocationPolicyGradesTableCompleteness",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-11",
              "digest": "sha256:952a78f224326566e450cb63319fb566b6c4835b3491ece3f7932a91fe3e4b1f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/skills_index_checks_test.go,omit,failed,1,TestCommandInvocationPolicyGradesTableCompleteness/body_prose_mention_of_the_key_is_inert,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestCommandInvocationPolicyGradesTableCompleteness$,passed,5\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,563\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,14,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestCommandInvocationPolicyGradesTableCompleteness/body_prose_mention_of_the_key_is_inert,\"skills_index_checks_test.go:537: body-prose mention of disable-model-invocation reddened the check: [\\\"command 'bench-write-spec' frontmatter carries disable-model-invocation though the invocation policy makes it model-invocable on Claude\\\"]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t11-conformance-meta",
          "performer": "ft358-author-11",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c8a1051010604cd1d2c0ef9a22f1be66f667273",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft358-author-11",
            "digest": "sha256:7b1fb9a77dd0964497c11c217b962a4fa232dfc095cafbe2bc2ebeed4b507a5c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/tier_live_tree_test.go,omit,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,conformance-meta,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,446\ncheck[1]{name,kind,tests_run,subjects}:\n  conformance-meta,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,31,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: conformance meta unregistered live-tree assertion TestMarkdownBlockOwnerHoldsOnTheLiveTree\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "9d5e681d32acc80cc0654018f783b1e23cb3a454",
          "requirement": "conformance-meta",
          "command": "bench test --check conformance-meta",
          "exit_code": 0,
          "probe": {
            "mutation": "internal/conformance/tier_live_tree_test.go: omit the TestMarkdownBlockOwnerHoldsOnTheLiveTree entry of classifiedLiveTreeTests, so the unregistered live-tree assertion diagnostic fires; expect TestRootConformance",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft358-author-11",
              "digest": "sha256:7b1fb9a77dd0964497c11c217b962a4fa232dfc095cafbe2bc2ebeed4b507a5c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/tier_live_tree_test.go,omit,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,conformance-meta,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,446\ncheck[1]{name,kind,tests_run,subjects}:\n  conformance-meta,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,31,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: conformance meta unregistered live-tree assertion TestMarkdownBlockOwnerHoldsOnTheLiveTree\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c5-standards",
          "performer": "ft358-review-c5-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c8a1051010604cd1d2c0ef9a22f1be66f667273",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c5-standards",
            "digest": "sha256:ec296054431f6bd29cbb104cdf7604364246f7a4e8ad03618ca87f9c9bea8ec2",
            "excerpt": "Standards axis, MB-C5 (base 15903e2c, tip 6ac6dc60). Two findings.\nS1 (low, medium, auto-fix decided): internal/conformance/markdown_block_owner_test.go:387-390 spells the tilde run, the frontmatter marker, and the two comment delimiters itself, while internal/markdown owns them (markdown.go:23, comments.go:7-8, markdown.go:215, :221). Repair: internal/markdown exports the frontmatter marker and the comment delimiters, and the check binds to them; build the tilde kind from the fence-marker length.\nS2 (low, high, auto-fix decided): internal/conformance/skills_index_checks_test.go:95-98 tests the first line against the frontmatter marker with no carriage-return strip, against Closed decision 8, and docs_workflow_checks_test.go:688-693 reads a fenced or commented H2 line as the README's first H2. Repair: move both helpers onto the block reader.\nThe profile row agrees with the registry; blockSpan in internal/adopt/marker.go collapses three split loops into one source.\n"
          },
          "axis": "Standards",
          "base": "15903e2ce3aab9a0db82da0678bf3a7359ec8d58",
          "tip": "9d5e681d32acc80cc0654018f783b1e23cb3a454",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c5-spec",
          "performer": "ft358-review-c5-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c8a1051010604cd1d2c0ef9a22f1be66f667273",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c5-spec",
            "digest": "sha256:9e1e6629f1308c00bf7409c235e6f404a4fa941e2f633749a8caee7e688da372",
            "excerpt": "Spec axis, MB-C5 (base 15903e2c, tip 6ac6dc60). Zero findings.\nbench coverage --check gives ok with 106 rows; every MB-C5 row cites its seam and no planned cell is left. Ticket 10 meets MB57, MB58, MB61, MB70, MB86, MB87, MB91, MB97 (fragment read, Bench-text fault condition, exact message bytes, CR-preserving rebuild); ticket 11 meets MB62-MB69, MB98, MB101 (six-kind bite proof, exemptions, registry and profile row, live-tree entry, the three readers on the block reader). Every line bound holds.\nOptional advice: the check also grades the root-level package, wider than spec.md:381; two block-rule copies stay in conformance test helpers; the check spells four literals itself.\n"
          },
          "axis": "Spec",
          "base": "15903e2ce3aab9a0db82da0678bf3a7359ec8d58",
          "tip": "9d5e681d32acc80cc0654018f783b1e23cb3a454",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c5-coverage",
          "performer": "ft358-review-c5-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c8a1051010604cd1d2c0ef9a22f1be66f667273",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft358-review-c5-coverage",
            "digest": "sha256:e04527d53bb33c4f541dc8efa99b679cd351c8bae6acf7dcad4b9dcb935d676b",
            "excerpt": "Coverage axis, MB-C5 (base 15903e2c, tip 6ac6dc60). Zero findings. Every probe restored; git status clean.\nEvery acceptance row has a test that fails under its named mutation: the Bench-text condition, the strip suffix, the blockSpan end guard, the equal-versus-contains rule, and the frontmatterHasKey class test all bit.\nOptional advice: skillDescriptionFolded has no body-inert case; the owner exemption has no sibling-path case; marker.go:56 and :59 and docs_workflow_checks_test.go:747 have silent mutants that do not change behavior today.\n"
          },
          "axis": "Coverage",
          "base": "15903e2ce3aab9a0db82da0678bf3a7359ec8d58",
          "tip": "9d5e681d32acc80cc0654018f783b1e23cb3a454",
          "finding_ids": [],
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
      "from": "sha256:7c5e675de050f405240c89a679c3214d28f457e88d0a29c7934ebc3cdcedd165",
      "to": "sha256:fcf4869a9e51f9cc4e151c012af5028c10e27b8cc88e6e18f5ffa7fa55d53a41",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ]
      }
    },
    {
      "from": "sha256:fcf4869a9e51f9cc4e151c012af5028c10e27b8cc88e6e18f5ffa7fa55d53a41",
      "to": "sha256:27f3d64fb8abd49d83a5d49138560b8d5ba8dcbb419b54cbbd4729efa68278ce",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ]
      }
    },
    {
      "from": "sha256:27f3d64fb8abd49d83a5d49138560b8d5ba8dcbb419b54cbbd4729efa68278ce",
      "to": "sha256:284407be4ffe0befb25fe8c8ac81e8cc73e5d70e181611d5895d02443cd6b416",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ]
      }
    },
    {
      "from": "sha256:284407be4ffe0befb25fe8c8ac81e8cc73e5d70e181611d5895d02443cd6b416",
      "to": "sha256:db3964c4226ddfa7d839b061b646cbd166c039cb8da1620bac16b6f552cdd706",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ]
      }
    },
    {
      "from": "sha256:db3964c4226ddfa7d839b061b646cbd166c039cb8da1620bac16b6f552cdd706",
      "to": "sha256:5282f157288e552be111dee9dbb280c348f121f3ef44d7fa67b8b3df386f68fe",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ]
      }
    },
    {
      "from": "sha256:5282f157288e552be111dee9dbb280c348f121f3ef44d7fa67b8b3df386f68fe",
      "to": "sha256:a0517fefa850069b8ac3657f4567d1e6ffff412ce1f57ded61459a0e554fb375",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ]
      }
    },
    {
      "from": "sha256:a0517fefa850069b8ac3657f4567d1e6ffff412ce1f57ded61459a0e554fb375",
      "to": "sha256:4b0b90422efdddcf683c1eb5a32d0bf68d153fed48f264671dfebc96e810e3db",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ]
      }
    },
    {
      "from": "sha256:4b0b90422efdddcf683c1eb5a32d0bf68d153fed48f264671dfebc96e810e3db",
      "to": "sha256:b12ab9e5b2d403cce066d1e10815d6524f165eecce1b3013a84eef59a7c3e75c",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ]
      }
    },
    {
      "from": "sha256:b12ab9e5b2d403cce066d1e10815d6524f165eecce1b3013a84eef59a7c3e75c",
      "to": "sha256:85d9574241a6429a3b3866d15ce2f8c10ea7a8832e6c5ba63d9573eb7b0c4aa4",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ]
      }
    },
    {
      "from": "sha256:85d9574241a6429a3b3866d15ce2f8c10ea7a8832e6c5ba63d9573eb7b0c4aa4",
      "to": "sha256:dd882480f65b3003489ea44e0c0ec2e2f4d9f36046718103bbee28583a6bff69",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ]
      }
    },
    {
      "from": "sha256:dd882480f65b3003489ea44e0c0ec2e2f4d9f36046718103bbee28583a6bff69",
      "to": "sha256:91964500432f92c3f973997f471817b3f52158f7bb423122eacd110baaa95596",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ]
      }
    },
    {
      "from": "sha256:91964500432f92c3f973997f471817b3f52158f7bb423122eacd110baaa95596",
      "to": "sha256:5e374001bec3bab9c1d9f5ede599951b1d8344ac5ce244bc1ad38319a78d34ab",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ]
      }
    },
    {
      "from": "sha256:5e374001bec3bab9c1d9f5ede599951b1d8344ac5ce244bc1ad38319a78d34ab",
      "to": "sha256:c31bcf3103d4d0c4225c463e4b60f241de6c3c6b338663b5093a9bf6b765733a",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ]
      }
    },
    {
      "from": "sha256:c31bcf3103d4d0c4225c463e4b60f241de6c3c6b338663b5093a9bf6b765733a",
      "to": "sha256:e5d6e08b78bc3f83cf7b2f61e8ef9579bad2cce480c69fcec4a72a109430a5ee",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ]
      }
    },
    {
      "from": "sha256:e5d6e08b78bc3f83cf7b2f61e8ef9579bad2cce480c69fcec4a72a109430a5ee",
      "to": "sha256:883b6bddb7b21f4ecb92530705af9a0faf9b3a10ccf428aff9805f5cede17939",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ],
        "MB-C4": [
          "MB-C4"
        ]
      }
    },
    {
      "from": "sha256:883b6bddb7b21f4ecb92530705af9a0faf9b3a10ccf428aff9805f5cede17939",
      "to": "sha256:c23ac50b1b8ab72bf7fca4145523c99c6a021058aa10f9b9752578c2e22ebf28",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ],
        "MB-C4": [
          "MB-C4"
        ]
      }
    },
    {
      "from": "sha256:c23ac50b1b8ab72bf7fca4145523c99c6a021058aa10f9b9752578c2e22ebf28",
      "to": "sha256:fc56090ad4c5b17107165bec91dc5d12946d0aacaebc407a376e940d178c334e",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ],
        "MB-C4": [
          "MB-C4"
        ]
      }
    },
    {
      "from": "sha256:fc56090ad4c5b17107165bec91dc5d12946d0aacaebc407a376e940d178c334e",
      "to": "sha256:6c72c6b28d025c0c98a527879153466994ffe1c855fef1ae3dd055075067a043",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ],
        "MB-C4": [
          "MB-C4"
        ]
      }
    },
    {
      "from": "sha256:6c72c6b28d025c0c98a527879153466994ffe1c855fef1ae3dd055075067a043",
      "to": "sha256:b9b3d587f8e5b0a278aeda03fdbf9f50b424e699eb649789d12eabe285e94a7b",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ],
        "MB-C4": [
          "MB-C4"
        ]
      }
    },
    {
      "from": "sha256:b9b3d587f8e5b0a278aeda03fdbf9f50b424e699eb649789d12eabe285e94a7b",
      "to": "sha256:a7ef5bc715c0df86eb39bf91e71bf940aeedc093cbfa00369948467107736c0d",
      "chunk_ids": {
        "MB-C1": [
          "MB-C1"
        ],
        "MB-C2": [
          "MB-C2"
        ],
        "MB-C3": [
          "MB-C3"
        ],
        "MB-C4": [
          "MB-C4"
        ],
        "MB-C5": [
          "MB-C5"
        ]
      }
    }
  ]
}
```

## Author witness reds

### MB-C1, ticket 1 (ft358-author-01, opus / xhigh)

- MB22: the base `prepare` gave `Findings() = [{Kind:sentence Line:1 Count:30}], want none` for a 30-word sentence inside a rejoined comment.
- The author ran 58 `bench probe` mutants, one for each coverage row and each new branch. Each mutant returned `bit` and failed only its own subtest. Two first spellings missed and the author spelled them again: one did not build, and one named a row with no prose.
- MB23: the `prose-mechanics` plan probe made `TestRootConformance` red with `gate: prose: ".agents/skills/bench-craft-tickets/SKILL.md" line 82: sentence of 26 words is over the 25-word bound`.
- The comment rule gave the result of `stripCommentsMapped` on 2,000,000 random inputs. MB79 allocates 2.42 times the input size.
- The `bench gate-prose` output over all 1000 tracked Markdown files is identical at the base and at the tip.

## MB-C1 review round 1

The repair targets are five findings on ticket 1, and one combined repair charge carries them all.

### Standards

Count: 1. Worst issue: S1.

- S1, low, confidence 6, auto-fix: `internal/markdown/markdown.go:187-193` declares `InSection` under the `Heading` comment. Declare it directly after `Section`.

### Spec

Count: 1. Worst issue: P1.

- P1, medium, confidence 8, auto-fix: `internal/markdown/comments.go:43-56` walks back across every earlier cut on a run of adjacent comments. Spec.md:210-211 and Closed decision 10 require one linear pass. Make the skip over a run of adjacent cuts constant time, and bound the work of a run of adjacent comments in a test.

### Coverage

Count: 3. Worst issue: C1.

- C1, medium, confidence 9, auto-fix: four mutants of `internal/markdown/comments.go:39`, `:53-56`, `:76`, and `:85` are silent. Add the chained rejoin, the join at offset 1, the three-kept-byte join, and the unterminated rejoin as rows of `TestReadBlockRules`.
- C2, low-medium, confidence 9, auto-fix: `internal/markdown/markdown.go:249` tests the H2 prefix on raw text without a red. Add a row whose heading follows a comment on the same line and whose commented line holds `## B`.
- C3, low-medium, confidence 8, auto-fix: omitting the frontmatter skip at `internal/markdown/markdown.go:244-247` is silent. Add a row whose frontmatter holds a fence marker.

## MB-C1 repair witness reds (ft358-repair-01, opus / xhigh, repair cycle 1 of 2)

- P1: with the old walk back, the adjacent-comment row gave `cutComments walked 50005000 steps over 10000 adjacent comments, want at most 80000`.
- C1: X1 is red against the omission of `top--` and against `for` to `if` in `push`. X2 is red against `after = at + len(commentOpen) - n`. X4 and X2 are red against `append` in place of `cut.push` in the unterminated branch. X3 is red against a two-byte `kept` window.
- C2: X6 is red against `strings.CutPrefix(line.Raw, H2Prefix)`.
- C3: X5 is red against the omission of the frontmatter `continue`.
- Each red failed only its new row, and each probe restored its subject.

## MB-C1 confirming round 1

Standards and Spec pass: S1 and P1 are closed. Coverage closes C1, C2, and C3, and it reports one new finding. The repair allowance consumed is 1 of 2.

### Standards

Count: 0.

### Spec

Count: 0.

### Coverage

Count: 1. Worst issue: N1.

- N1, medium, confidence 8, auto-fix: omitting `c.first = c.first[:i]` at `internal/markdown/comments.go:95` is silent. Add the row `<!<!----><!---->----><!<!---->---->` to `TestReadBlockRules`, with the spans 0-21 and 21-35 and empty text. The repair is test-only, so a coordinator probe that bites closes it through `bench record confirmation`.

### MB-C2, ticket 2 (ft358-author-02, opus / high)

- MB24: the old code gave `LiveSpecSlugs = ["hidden"], want none` for a tilde-fenced spec path.
- MB25: the old code gave `Facts = [...Status:"implemented"...], want one fact with the status staged` for an indented fence.
- MB26: the old code gave `Implemented = "", spec has more than one Status: staged line`.
- MB80: the old code gave `Implemented = "", spec has no Status: staged line` on a CRLF spec.
- MB29 and MB30: the old code gave `references story 3, which the spec does not declare (has: 1, 2)` and `leaves story 2 unreferenced`.
- MB83: a marker match on comment-free text gave `State = "mapped", want historical`.
- `bench coverage --check` gives the same output and exit code at the base and at the tip on all 14 live specs.

### MB-C2, ticket 3 (ft358-author-03, opus / high)

- MB31 and MB81: the old code gave `failures = [{ROADMAP.md malformed roadmap row  11}]` for a tilde-fenced and a backtick-fenced row example.
- MB32: the old code gave `sequence = [{Rank:1 Text:quoted Command:}]` for a tilde-fenced sequence heading.
- MB33: the old code gave `roadmap/FT1.md: duplicate Next: line at line 4; a row carries one`.
- MB99: the old projection replaced the quoted sequence section.
- MB100: the old closure removed the quoted table row and kept `| FT2 | FT1 |` in the live table.
- `bench roadmap` gives byte-identical output at the base and at the tip.

### MB-C2, ticket 4 (ft358-author-04, opus / high)

- MB34: the old code gave `Scan diagnostics = []string{"duplicate Writes"}, want none`.
- MB35: the old code gave `diagnostics = nil, want ["unterminated fence"]`.
- MB36 and MB37: the old code gave `["unterminated fence", "missing Acceptance section"]` for a balanced four-backtick block.
- MB95: a `Blocks` that returns `markdown.Read` lost the fence fault after an unclosed `---` first line.
- MB102: a `Scan` that drops `FenceLabel` failed `TestMetadataGrammarFenceRoles` at `premises_test.go:89`.
- `ParseTicket` gives byte-identical diagnostics, `Covers`, and `Writes` at the base and at the tip over all 1001 tracked Markdown files.

## MB-C2 review round 1

The orchestrator decides each finding under the reviewer's advance approval of the `--full` run. One combined repair charge carries every repair target, and it routes to ticket 3, which holds most of the paths. Repair allowance consumed for MB-C2 before this repair: 0 of 2.

### Standards

Count: 3. Worst issue: S1.

- S1, medium, high confidence, auto-fix (decided): `internal/roadmap/closure.go:19` and `:27`. `Close` counts a fenced row heading that `ParseDocument` skips. Make `Close` skip non-body lines. New row MB103 covers it.
- S2, low, medium confidence, no-op (decided): `internal/tickets/tickets.go:108`. The reader reports the comment fault first (MB71), so FA11 records the result. The C2 test pins it.
- S3, low, high confidence, no-op: the triple index read in `internal/roadmap` matches by construction.

### Spec

Count: 2. Worst issue: P1.

- P1, severity 5, confidence 10, auto-fix: `internal/coverage/rows.go` grows from 301 to 303 lines. Cut two net lines.
- P2, severity 4, confidence 8, auto-fix (decided): the same defect as S1.
- The non-behavioral contradiction about fragment callers is corrected in the spec prose, in Closed decision 16, and in FA10 to FA12.

### Coverage

Count: 5. Worst issue: C1.

- C1, medium, confidence 8, auto-fix: add a first-line `---` to `---` frontmatter case to `TestMetadataSkipsAFencedStatus`, `TestParseSkipsFencedCoverageLines`, and `TestParseDocumentSkipsTildeFences`.
- C2, medium, confidence 9, auto-fix: add an unterminated-comment case to `TestParseTicketReadsTheBlockReaderFault` that expects no `unterminated fence`.
- C3, medium, confidence 9, auto-fix: add a fenced `1. fake` line inside the live sequence to `TestParseDocumentSkipsTildeFences`.
- C4, low, confidence 7, auto-fix: add a rejected-side title case for `internal/roadmap/sequence_projection.go:32`.
- C5, low, confidence 7, auto-fix (decided): `internal/maps/fields.go:92`. A closer's label is the text after a run of the opener's length. Thus a longer nested run inside the planning-facts fence still refuses (story 78). Add that case to `TestMetadataGrammarFenceRoles`.

### MB-C2 repair witness reds (ft358-repair-02, opus / high, repair cycle 1 of 2)

- MB103: the old regex scan in `Close` gave `closure = "", closed row FT1 appears 2 times in ROADMAP.md`.
- C5: the old whole-run label gave an empty `FenceLabel` for a longer nested closer.
- C1: `Read` to `ReadFragment` gave `LiveSpecSlugs = ["hidden"]` in spec, `coverage map missing the canonical header` in coverage, and `malformed roadmap row` in roadmap.
- C2: a filter that accepts every fault kind gave `diagnostics = []string{"unterminated fence"}, want []string(nil)`.
- C3: a disabled fence skip gave `sequence = [{Rank:1 Text:live} {Rank:1 Text:fake}]`.
- C4: the compare `H2Prefix+line.Section` turned both rejected headings red.
- P1: `internal/coverage/rows.go` ends at 301 lines.

## MB-C2 confirming round 1

Standards closes S1, Spec closes P1 and P2, and Coverage closes C1 to C5. The round reports four new findings. The repair allowance consumed for MB-C2 is 1 of 2, so the next repair is cycle 2 of 2.

### Standards

Count: 2. Worst issue: S4.

- S4, low-medium, high confidence, auto-fix: `internal/roadmap/closure.go:31` and `internal/roadmap/tree.go:116` and `:123` state the row-heading rule twice. Add one `internal/roadmap` helper that both call.
- S5, low, medium confidence, no-op: `openerRun` in `internal/maps/fields.go` derives the run from the reader and copies no marker rule.

### Spec

Count: 1. Worst issue: P3.

- P3, low, confidence 8, auto-fix: the ticket 3 body did not name MB103. The orchestrator corrects the ticket prose in this commit.

### Coverage

Count: 2. Worst issue: N1.

- N1, low, confidence 8, auto-fix: three closer-label mutants in `internal/maps/fields.go` are silent. Add three cases to `TestMetadataGrammarFenceRoles`. They are a four-run opener with a four-run closer, an indented closer, and a spaced info string with a longer closer.
- N2, low, confidence 7, auto-fix: the frontmatter parity of `Close` has no test. Add a first-line frontmatter block that holds a row heading to `TestCloseSkipsAFencedRowHeading`.

### MB-C2 repair 2 witness reds (ft358-repair-02b, opus / high, repair cycle 2 of 2)

- S4: `rowHeading` in `internal/roadmap/tree.go` is the one row-heading predicate, and `ParseDocument` and `Close` both call it. The roadmap, roadmapflow, maps, and ownership packages pass.
- N1: `run = markdown.FenceMarker` failed the four-run case, `strings.CutPrefix` failed the indented case, and `openerRun` without its outer `TrimSpace` failed the spaced-info case. Each failed at `premises_test.go:101`.
- N2: `markdown.ReadFragment` in `Close` gave `closure = "", closed row FT1 appears 2 times in ROADMAP.md`.

## MB-C2 confirming round 2

All three axes pass. S4, P3, N1, and N2 are closed. The repair allowance consumed for MB-C2 is 2 of 2.

### Standards

Count: 0. Optional advice: a named-case table in `premises_test.go`, and a test name that also covers the frontmatter case.

### Spec

Count: 0.

### Coverage

Count: 0.

### MB-C3, ticket 5 (ft358-author-05, opus / high)

- MB39 and MB40: the old toggle in `OpenFence` and `UnfencedLines` gave the wrong fence on three shapes. They are a longer closer, a closer of the other character, and a shorter run.
- MB93: `OpenFence` through `Read` gave `(0, false), want (2, true)` for a State body whose first line is `---`.
- MB38: the old toggle in `splitSections` split a section at a fenced `## X` line.
- MB82: without the comment-fault refusal, `readStateFile` gave no refusal or the section-heading refusal.
- `handoffdoc.Parse` gives byte-identical output at the base and at the tip on the live session handoff.

### MB-C3, ticket 6 (ft358-author-06, opus / high)

- MB41: a body loop on the raw H2 prefix split a second entry at a fenced dated heading.
- MB42 to MB44: a parse that drops the fault record failed all three fault subtests.
- MB45: a retros parse without the body-class test accepted a fenced `## Outcome`.
- MB46 and MB47: a raw H2 branch and a unit scan without the body class failed their rows.
- The `ResolvedBy` scan without the fence condition failed its case, and each fault reason string failed on its own.
- `bench learnings`, the drain journal read, and `retros.Parse` on the live retro give identical output at the base and at the tip.

### MB-C3, ticket 9 (ft358-author-09, opus / high)

- MB59: the old split and exact compare gave `index = "", want "x"` for a CRLF skill file.
- MB60: a reader fault read as frontmatter gave `index = "smuggled", want ""` for an unclosed opener.
- A disabled stop at the first body line failed eight cases.
- `bench skills-index --check` gives no diagnostic, and `--write` changes no file at the tip.

## MB-C3 review round 1

The orchestrator decides each finding under the reviewer's advance approval of the `--full` run. One combined repair charge carries every repair target, and it routes to ticket 6, which holds most of the paths. Ticket 6 widens to the handoff paths. The repair allowance consumed for MB-C3 before this repair is 0 of 2.

### Standards

Count: 6. Worst issue: S2.

- S1, low, auto-fix: correct the package doc at `internal/handoffdoc/document.go:7`, so that it names the real imports.
- S2, low, auto-fix: `handoffdoc` exports one State-body read, and `readStateFile` uses it (with P1 and MB106).
- S3, low, auto-fix: `learnings.Parse` reads the reader's lines and drops its parallel split.
- S4, low, no-op: the duplicate `ShapeHeading` predates the build and sits outside every `Writes:` line. A learning records it.
- S5, low, auto-fix (decided): each heading title comes from the comment-free text (FA15).
- S6, low, auto-fix: remove the row-ID provenance from the new test doc comments.

### Spec

Count: 4. Worst issue: P1.

- P1, medium, auto-fix: the State guard refuses each line that the reader reads as an H2 heading (MB106).
- P2, low, auto-fix (decided): FA13 records the trailing-space heading.
- P3, low, auto-fix (decided): `handoffdoc.Parse` refuses an unterminated comment (MB105, FA14).
- P4, low, no-op (decided): `retros` keeps comment-free matching, and the grammar-text table now says so.

### Coverage

Count: 5. Worst issue: C1.

- C1, high, auto-fix (decided): `handoffdoc.Parse` refuses an open fence behind an unclosed first-line `---` (MB104).
- C2, medium, auto-fix (decided): the same repair as P3.
- C3, medium, auto-fix: add an MB82 case whose first line is `---`.
- C4, medium, auto-fix: add a commented-heading case at `handoffdoc/document.go:262`, `learnings.go:113`, and `recommendations.go:25`.
- C5, low, auto-fix (decided): add a first-line frontmatter case at `retros.go:75`, `recommendations.go:20`, and `learnings.go:277`. FA10 now names the journal and the retrospectives.

### MB-C3 repair witness reds (ft358-repair-03, opus / high, repair cycle 1 of 2)

- MB104: the old code gave `Parse refusal at line 0, want the fence refusal at line 8`.
- MB105: the old code gave `Parse refusal at line 0 (a refusal: false), want the comment refusal at line 7`.
- MB106: the old guard gave no refusal kind for a trailing-space heading and for a heading after a comment.
- FA15: a title read from raw text refused a heading behind a leading comment.
- C3, C4, and C5: each named mutant (`Read` for `ReadState`, a raw-prefix heading test, `ReadFragment` for `Read`) failed its new case.
- S3: the CR strip, the end bound, and the guard of `sourceBody` each failed a test when omitted.

## MB-C3 confirming round 1

Standards closes S1, S2, S3, S5, and S6. Spec closes P1 to P3, and Coverage closes C1 to C5. The round reports six new findings. The repair allowance consumed for MB-C3 is 1 of 2, so the next repair is cycle 2 of 2.

### Standards

Count: 3. Worst issue: S8.

- S7, low, auto-fix: rewrite the `OpenFence` doc in `internal/handoffdoc/fence.go`, so that it states only what the function reports.
- S8, low, auto-fix: export one `handoffdoc` function that maps a fault kind to its repair text, and call it from `fence.go` and `internal/handoff/state_file.go`.
- S9, low, auto-fix: the `sourceBody` doc in `internal/learnings/entry.go` says that it trims the newline bytes at each end.

### Spec

Count: 2. Worst issue: P5.

- P5, low, auto-fix (decided): the heading refusal at `internal/handoffdoc/document.go:356` quotes the raw line, so the existing message keeps its bytes.
- P6, low, auto-fix: the grammar-text row for the handoff packages now says that only the heading grammar is comment-free. The orchestrator corrects it in this commit.

### Coverage

Count: 1. Worst issue: C6.

- C6, low, auto-fix: add a case beside MB104 whose closed first-line frontmatter block holds a fence marker line, which `Parse` accepts.

### MB-C3 repair 2 witness reds (ft358-repair-03b, opus / high, repair cycle 2 of 2)

- S8: a shared `FaultRepair` that maps the comment fault to the fence repair failed three cases, for example the hint-names-the-line case at `state_file_test.go:256`.
- P5: `TestParseQuotesTheRawHeading` with the raw quote fed from comment-free text failed all three rows.
- C6: `if true` in place of the frontmatter-fault narrowing gave `Parse refusal at line 2, want ... line 0`.
- S7 and S9 change comments only.

## MB-C3 confirming round 2

Standards closes S7, S8, and S9, and Coverage closes C6. Spec closes P5 and P6 and reports P7. The repair allowance for MB-C3 is spent (2 of 2). The reviewer's advance approval of the `--full` run covers repair rounds past the allowance, so repair cycle 3 takes P7 as an overrun.

### Standards

Count: 0. The optional advice is in the excerpt, and a learning records the two handoff items.

### Spec

Count: 1. Worst issue: P7.

- P7, low, confidence 8, auto-fix: `internal/handoffdoc/document.go:298` sets the legacy-mix section from the comment-free heading, so the refusal at `:305` changes its bytes. Quote `b.raw` at `:298` and at `:320`, and add a legacy-mix row to `TestParseQuotesTheRawHeading`.

### Coverage

Count: 0. The optional advice is in the excerpt.

### MB-C3 repair 3 witness reds (ft358-repair-03c, opus / high, repair cycle 3, an overrun)

- P7: the old code failed the legacy-mix row and the legacy next-command row of `TestParseQuotesTheRawHeading`, because each refusal quoted the comment-free heading.

## MB-C3 confirming round 3

All three axes pass, and P7 is closed. MB-C3 took three repair cycles: two from the allowance and one overrun under the reviewer's advance approval.

### Standards

Count: 0.

### Spec

Count: 0.

### Coverage

Count: 0.

### MB-C4, ticket 7 (ft358-author-07, opus / high)

- MB48: the old exact-line opener rule gave a duplicate-fence refusal for a quoted record fence.
- MB96: a reader that refuses any fault refused a record whose first line is an unclosed `---`.
- MB49 to MB51: an accepted unterminated fence, a first-opener rule, and a shifted span end each failed their Render test.
- MB107: a writer without the encoding gave an unterminated-fence refusal after a render of an excerpt that holds a comment opener.
- The record payloads of all 129 review and spec files decode the same at the base and at the tip.

### MB-C4, ticket 8 (ft358-author-08, opus / xhigh)

- MB52, MB53, and MB85: the base `locate.go` failed the tilde-fenced heading, the fenced duplicate heading, and the indented heading cases of `TestLocateReadsTheBlockReader`.
- MB94: `Read` in place of `ReadFragment` for the section body gave a count of 0.
- The line-start, origin, span, and strip mutants each failed their own case, and an origin mutant failed `TestLocateStripsRejoinedComments`.
- MB56: `bench anchors` gives identical verdicts, and all 435 live anchors keep their located lines at the base and at the tip.

## MB-C4 review round 1

The orchestrator decides each finding under the reviewer's advance approval of the `--full` run. One combined repair charge carries every repair target, and it routes to ticket 8, which widens to the review-record repair paths. The repair allowance consumed for MB-C4 before this repair is 0 of 2.

### Standards

Count: 2. Worst issue: S1.

- S1, medium, auto-fix: the `StripHTMLComments` doc in `internal/anchors/registry.go` says the function returns the text without the comments that the block reader removes. It restates no comment rule.
- S2, low, auto-fix: the walk doc in `internal/anchors/locate.go` says that the walk reads the whole document.

### Spec

Count: 1. Worst issue: P1.

- P1, medium, auto-fix (decided): a record closer whose label is the record name refuses with the existing duplicate message. Story 84 and row MB108 cover it.

### Coverage

Count: 1. Worst issue: C1.

- C1, medium, auto-fix (decided): the anchor section walk reads the whole comment-free document through `Read`, as the comment strip does. MB94 now states the frontmatter case, and Closed decision 16 drops the walk from the fragment callers. The repair also adds an FA7 opener test as hardening.

### MB-C4 repair witness reds (ft358-repair-04, opus / xhigh, repair cycle 1 of 2)

- MB94: `ReadFragment` at the section walk gave `section count = 0, want 1` for leading frontmatter that holds a comment opener, and the same for a tilde marker.
- MB108: the old reader gave no refusal for a record fence followed at once by a second record opener.
- FA7: the exact-line opener failed the indented and tilde cases, and a prefix test failed all three cases. An empty-info closer test failed the trailing-text closer case.
- S1 and S2 change comments only. `bench anchors .` gives the same verdicts for all 42 files.

## MB-C4 confirming round 1

Standards closes S1 and S2. Spec closes P1, and Coverage closes C1. The round reports three new findings, and two of them name one defect. The repair allowance consumed for MB-C4 is 1 of 2, so the next repair is cycle 2 of 2.

### Standards

Count: 0. The optional advice is in the excerpt.

### Spec

Count: 2. Worst issue: P2.

- P2, medium, auto-fix: the ticket 8 body said that the walk uses `ReadFragment`. The orchestrator corrects the ticket prose and the FA17 wording in this commit.
- P3, low, auto-fix (decided): the same defect as C2.

### Coverage

Count: 1. Worst issue: C2.

- C2, medium, auto-fix: the anchor section walk takes the line-1 frontmatter rule from the comment-free text. Thus a comment followed by `---` on line 1 hides section `S`. Take the frontmatter state from the original document. Story 85 and row MB109 cover it.

### MB-C4 repair 2 witness reds (ft358-repair-04b, opus / xhigh, repair cycle 2 of 2)

- MB109: the old walk gave `Locate = 0, want line 3` for a first line of a comment followed by `---`.
- A walk over the comment-free text failed the MB109 row and a closer row, and an evaluator fed the comment-free text failed its section-body case.
- The empty-range guard, the range offsets, and a body that keeps its comments each failed their own row.
- A binary built from the chunk base and the tip binary give byte-identical anchor results, and all 435 anchors keep their located lines.

## MB-C4 confirming round 2

Spec closes P2 and P3, and Coverage closes C2. The round reports two new findings. The repair allowance for MB-C4 is spent (2 of 2). The reviewer's advance approval of the `--full` run covers this overrun, which is test-only plus one comment correction.

### Standards

Count: 1. Worst issue: S1.

- S1, low, auto-fix: the `MarkdownH2Sections` doc in `internal/anchors/match.go` names `commentFreeSectionMapped` as the section reader of `Locate` and the evaluator. This change is comment-only and evidence-only.

### Spec

Count: 0.

### Coverage

Count: 1. Worst issue: C3.

- C3, medium, auto-fix: add an evaluator row in `TestLocateReadsTheBlockReader` for a needle inside a comment in a section. The repair is test-only, so a coordinator probe that bites closes it through `bench record confirmation`.

## MB-C4 confirming round 3

All three axes pass, and S1 and C3 are closed. MB-C4 took three repair cycles: two from the allowance and one test-and-comment overrun under the reviewer's advance approval. The coordinator probe on `internal/anchors/registry.go:187` bit on the evaluator row.

### Standards

Count: 0.

### Spec

Count: 0.

### Coverage

Count: 0.

### MB-C5, ticket 10 (ft358-author-10, opus / high)

- MB57 and MB58: the old marker scan refused a tilde-fenced marker example as malformed.
- MB61 and MB70: the old marker scan gave no refusal for an unterminated tilde block and for a comment fault.
- MB86, MB97, MB87, and MB91: the Bench-text condition, the fragment read, and the CR-preserving rebuild each failed their row when omitted.
- The live AGENTS.md rewrites and strips to the same bytes at the base and at the tip.

### MB-C5, ticket 11 (ft358-author-11, opus / high)

- MB62 to MB67: a stubbed check gave no diagnostic for each of the six literal kinds.
- MB68: each exemption, removed alone, turned its side of the exemption case red.
- MB69: before the two fixtures moved to the fence constant, the live-tree test gave four diagnostics.
- MB101: without its entry, `conformance-meta` reported the unregistered live-tree assertion.
- MB98: the old `skillDescriptionFolded`, `frontmatterHasKey`, and `h2Headings` failed `TestConformanceReadersUseTheBlockReader` on CRLF frontmatter and on a fenced or commented heading.

## MB-C5 review round 1

The orchestrator decides each finding under the reviewer's advance approval of the `--full` run. One repair charge carries both findings, and it routes to ticket 11, which widens to `internal/markdown`. The repair allowance consumed for MB-C5 before this repair is 0 of 2.

### Standards

Count: 2. Worst issue: S1.

- S1, low, auto-fix (decided): `internal/markdown` exports the frontmatter marker and the two comment delimiters, and `markdown-block-owner` binds each literal kind to an exported constant.
- S2, low, auto-fix (decided): the skill frontmatter test of `checkSkillFrontmatter` and the README first-H2 read move onto the block reader. Closed decision 9 now names them.

### Spec

Count: 0. The spec text now says that the check grades every non-test Go file of the module.

### Coverage

Count: 0. The repair also adds two hardening rows from the advice: a body-inert case for `skillDescriptionFolded`, and a sibling-path case for the owner exemption.
