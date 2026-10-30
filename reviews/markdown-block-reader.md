# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/markdown-block-reader/spec.md",
  "plan_digest": "sha256:dd882480f65b3003489ea44e0c0ec2e2f4d9f36046718103bbee28583a6bff69",
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
            "excerpt": "Spec axis, MB-C1 (base 36a07d83, tip 252408eb). One finding; every other mapped row (MB1-MB19, MB71-MB79, MB88-MB90, MB92) and acceptance item matches.\nP1 (medium, confidence 8, auto-fix): internal/markdown/comments.go:43-56 joinedOpener. Spec.md:210-211 and Closed decision 10 require one linear pass. For adjacent comments with no kept byte between them (repeated <!---->), the walk-back steps across every earlier cut, so the work is quadratic. MB79 measures allocation only and its comments are 93 bytes apart, so it stays green. Repair: make the skip over a run of adjacent cuts O(1), and add a run of adjacent comments to the comment-dense input under a bound.\nChecked: comment rule matches stripCommentsMapped; Closed decisions 5, 8, 14, 16; fault order frontmatter, comment, fence; exported constants; prose prepare fault mapping; no banned literal in internal/prose non-test files.\n"
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
            "excerpt": "Coverage axis, MB-C1 confirming round (repair delta 252408eb..503aef49). C1, C2, C3 closed; one new finding.\nRound 1 mutants re-run, all red: omit top-- (comments.go:62), for->if in push (:91), at+len(commentOpen)-n (:82), append in the unterminated branch (:33), CutPrefix on line.Raw (markdown.go:136), omit the frontmatter continue (markdown.go:133), two-byte kept window (comments.go:56).\nNew mutants: omit top = c.first[top] red (MB79); first = i red (MB79); End <= s.Start red (MB12).\nN1 (medium, confidence 8, auto-fix): internal/markdown/comments.go:95. Omitting c.first = c.first[:i] is silent. Witness input <!<!----><!---->----><!<!---->----> reads comment 0-21 and comment 21-35 with empty text; the mutant reads comments 0-21 and 23-30 with visible <!---->. Repair: add that row to TestReadBlockRules.\n"
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
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/markdown/comments.go,omit,failed,1,TestReadBlockRules/MB12_a_rejoin_after_a_rejoin_of_two_adjacent_removals_is_removed_too,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/markdown,^TestReadBlockRules$,passed,32\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,386\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/markdown,fail,5,32\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/markdown,TestReadBlockRules/MB12_a_rejoin_after_a_rejoin_of_two_adjacent_removals_is_removed_too,\"markdown_test.go:68: reading of \\\"<!<!----><!---->----><!<!---->---->\\\":\",8\nskips[0]{package,test,reason}:\n"
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
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/markdown/comments.go,omit,failed,1,TestReadBlockRules/MB12_a_rejoin_after_a_rejoin_of_two_adjacent_removals_is_removed_too,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/markdown,^TestReadBlockRules$,passed,32\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,386\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/markdown,fail,5,32\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/markdown,TestReadBlockRules/MB12_a_rejoin_after_a_rejoin_of_two_adjacent_removals_is_removed_too,\"markdown_test.go:68: reading of \\\"<!<!----><!---->----><!<!---->---->\\\":\",8\nskips[0]{package,test,reason}:\n"
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
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/coverage/rows.go,swap,failed,1,TestParseSkipsFencedCoverageLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/coverage,^TestParseSkipsFencedCoverageLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,399\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/coverage,TestParseSkipsFencedCoverageLines,\"coverage_test.go:286: Check = [coverage map leaves story 2 unreferenced; add a row or a `Not covered: story 2 — <reason>` line], want none\",2\nskips[0]{package,test,reason}:\n"
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
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/coverage/rows.go,swap,failed,1,TestParseSkipsFencedCoverageLines,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/coverage,^TestParseSkipsFencedCoverageLines$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,399\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/coverage,TestParseSkipsFencedCoverageLines,\"coverage_test.go:286: Check = [coverage map leaves story 2 unreferenced; add a row or a `Not covered: story 2 — <reason>` line], want none\",2\nskips[0]{package,test,reason}:\n"
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
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:165: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", <nil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
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
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:165: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", <nil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
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
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:170: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", <nil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
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
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:170: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", <nil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
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
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:170: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", <nil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
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
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/roadmap/sequence_projection.go,swap,failed,1,TestProjectSequenceSkipsFencedHeadings,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/roadmap,^TestProjectSequenceSkipsFencedHeadings$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,566\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/roadmap,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/roadmap,TestProjectSequenceSkipsFencedHeadings,\"commitment_test.go:170: projection = \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n\\\", <nil>; want \\\"# Roadmap\\\\\\\\n\\\\\\\\n~~~\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. quoted\\\\\\\\n~~~\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n"
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
      "tip": "251b6e3adf4950681b0227de1fc93449dc8910ef",
      "plan_digest": "sha256:dd882480f65b3003489ea44e0c0ec2e2f4d9f36046718103bbee28583a6bff69",
      "source_digest": "183a350962746062b7ffb6fc133d42339b3ace33",
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
