# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/markdown-block-reader/spec.md",
  "plan_digest": "sha256:fcf4869a9e51f9cc4e151c012af5028c10e27b8cc88e6e18f5ffa7fa55d53a41",
  "implementation_session": "",
  "chunks": [
    {
      "id": "MB-C1",
      "base": "36a07d83591f99cc079f571f185d25e606db0a2c",
      "tip": "503aef49a0de6507e0690eb4161aa0e09a71536d",
      "plan_digest": "sha256:fcf4869a9e51f9cc4e151c012af5028c10e27b8cc88e6e18f5ffa7fa55d53a41",
      "source_digest": "1efde0ff74f0f0f1591151bfbd6a1a7aa168d711",
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
