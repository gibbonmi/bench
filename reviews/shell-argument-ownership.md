# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/shell-argument-ownership/spec.md",
  "plan_digest": "sha256:c20f6c65e343e0a2334a915e1df675c951ffcad45fc3ce413f0cd12dd2a0d76e",
  "implementation_session": "",
  "chunks": [
    {
      "id": "S-A",
      "base": "1586949596564559881814a1695df97030c333ec",
      "tip": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
      "plan_digest": "sha256:dd10c1e33a504fed07206dd2621197369959e9b0cc4f50eaa6fef622276f7818",
      "source_digest": "3af663fa6facc743e27596274e13d17e252df1cc",
      "acceptance_rows": [
        "S01",
        "S02",
        "S03",
        "S04",
        "S15"
      ],
      "verification": [
        {
          "id": "S-A-T01-owner-omission-proof",
          "performer": "ft354-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3af663fa6facc743e27596274e13d17e252df1cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-01",
            "digest": "sha256:06406b1862d571c5906704afe5bd3cc34d1a57f3e569ff00034ef3124ee8e8b4",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,3,TestShellQuoteAlwaysQuoted/a-b_c.d/e,24,1,0,null,yes\n  2,bit,swap,failed,3,TestShellQuoteRoundTrip/a_leading_quote,24,1,0,null,yes\n  3,bit,swap,failed,1,TestShellQuoteRoundTrip/empty,24,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/sanitize,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
          "requirement": "owner-omission-proof",
          "command": "bench test --package ./internal/sanitize",
          "exit_code": 0,
          "probe": {
            "mutation": "Return bare alpha, split an embedded quote, or drop the empty owner argument. The exact quoted-byte and NUL-framed recovered-argv owner assertions must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-01",
              "digest": "sha256:06406b1862d571c5906704afe5bd3cc34d1a57f3e569ff00034ef3124ee8e8b4",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,insert,failed,3,TestShellQuoteAlwaysQuoted/a-b_c.d/e,24,1,0,null,yes\n  2,bit,swap,failed,3,TestShellQuoteRoundTrip/a_leading_quote,24,1,0,null,yes\n  3,bit,swap,failed,1,TestShellQuoteRoundTrip/empty,24,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/sanitize,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "S-A-T01-owner-preservation",
          "performer": "ft354-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3af663fa6facc743e27596274e13d17e252df1cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-01",
            "digest": "sha256:994659707a1f3cab7c1f3a6ced0456007d5473239a7ba597e0495b1440177c8b",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,8,TestShellQuoteRoundTrip/Unicode,24,1,0,null,yes\n  2,bit,swap,failed,1,TestShellQuoteRoundTrip/Unicode,24,1,0,null,yes\n  3,bit,swap,failed,3,TestShellQuoteRoundTrip/a_leading_quote,24,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/sanitize,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
          "requirement": "owner-preservation",
          "command": "bench test --package ./internal/sanitize",
          "exit_code": 0,
          "probe": {
            "mutation": "Alter the existing always-quoted owner result for spaces, Unicode, or embedded quotes. The independent input-to-argv witness must fail, then pass after restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-01",
              "digest": "sha256:994659707a1f3cab7c1f3a6ced0456007d5473239a7ba597e0495b1440177c8b",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,8,TestShellQuoteRoundTrip/Unicode,24,1,0,null,yes\n  2,bit,swap,failed,1,TestShellQuoteRoundTrip/Unicode,24,1,0,null,yes\n  3,bit,swap,failed,3,TestShellQuoteRoundTrip/a_leading_quote,24,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/sanitize,all,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "S-A-T01-caller-bypass-proof",
          "performer": "ft354-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3af663fa6facc743e27596274e13d17e252df1cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-01",
            "digest": "sha256:ccbdf5993d4f370d001cd49b8ae806f17bed596b70619b15e4f621232f643e4f",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,10,TestTargetRefusalEscapesOnlyAnUnsafeLine/a_control_byte_in_the_detail,12,1,0,null,yes\n  2,bit,swap+swap,failed,2,TestTargetRefusalEscapesOnlyAnUnsafeLine/a_plain_label,12,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,TreeBuildArgument|TargetRefusal,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
          "requirement": "caller-bypass-proof",
          "command": "bench test --package ./internal/worktree --run 'TreeBuildArgument|TargetRefusal'",
          "exit_code": 0,
          "probe": {
            "mutation": "Bypass sanitize.ShellQuote in PrintTreeBuildRefusal or restore its former bare alpha spelling. The actual tree-build producer fixture and recovered argv assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-01",
              "digest": "sha256:ccbdf5993d4f370d001cd49b8ae806f17bed596b70619b15e4f621232f643e4f",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,10,TestTargetRefusalEscapesOnlyAnUnsafeLine/a_control_byte_in_the_detail,12,1,0,null,yes\n  2,bit,swap+swap,failed,2,TestTargetRefusalEscapesOnlyAnUnsafeLine/a_plain_label,12,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,TreeBuildArgument|TargetRefusal,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "S-A-T01-caller-preservation",
          "performer": "ft354-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3af663fa6facc743e27596274e13d17e252df1cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-01",
            "digest": "sha256:92abc0b0bbc763578af04a2948ff056a20d73a0f73318594c959135d3a1db573",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,10,TestTargetRefusalEscapesOnlyAnUnsafeLine/a_control_byte_in_the_detail,12,1,0,null,yes\n  2,bit,swap+swap,failed,1,TestTreeBuildArgument/an_embedded_quote,12,1,0,null,yes\n  3,bit,swap,failed,3,TestTargetRefusalEscapesOnlyAnUnsafeLine/a_control_byte_in_the_detail,12,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,TreeBuildArgument|TargetRefusal,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
          "requirement": "caller-preservation",
          "command": "bench test --package ./internal/worktree --run 'TreeBuildArgument|TargetRefusal'",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the tree-build operand, change its exact quoted spelling, or bypass retained whole-line control escaping. The real producer or refusal assertion must fail, then pass after restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-01",
              "digest": "sha256:92abc0b0bbc763578af04a2948ff056a20d73a0f73318594c959135d3a1db573",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap,failed,10,TestTargetRefusalEscapesOnlyAnUnsafeLine/a_control_byte_in_the_detail,12,1,0,null,yes\n  2,bit,swap+swap,failed,1,TestTreeBuildArgument/an_embedded_quote,12,1,0,null,yes\n  3,bit,swap,failed,3,TestTargetRefusalEscapesOnlyAnUnsafeLine/a_control_byte_in_the_detail,12,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/worktree,TreeBuildArgument|TargetRefusal,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "S-A-R1-standards",
          "performer": "ft354-review-S-A",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3af663fa6facc743e27596274e13d17e252df1cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-A",
            "digest": "sha256:f398efeaf80bd5fd65ed8f29a5c0df5292cec2758a5aa8d3156e346b1a822bdc",
            "excerpt": "## Standards\nCount: 2, both no-op. Worst: the shell-argv-plus-sentinel test pattern is hand-written per package.\n- internal/worktree/shell_arguments_test.go:161-184 and internal/sanitize/shellquote_test.go:25-35: per-package sentinel setup around axitest.HelpAction.Argv; a shared axitest helper would prevent a pasted harness in tickets 04 and 07. Confidence 4. no-op for S-A.\n- internal/worktree/shell_arguments_test.go:174 repeats the refusal template of internal/treetarget/build_test.go:76; independent expectations with recorded reds meet the AGENTS.md exception. Confidence 3. no-op.\n"
          },
          "axis": "Standards",
          "base": "1586949596564559881814a1695df97030c333ec",
          "tip": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "S-A-R1-spec",
          "performer": "ft354-review-S-A",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3af663fa6facc743e27596274e13d17e252df1cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-A",
            "digest": "sha256:89f10fc7f826a248d6b94052df797739618dea557a98656be846cecca42bb6bd",
            "excerpt": "## Spec\nCount: 1, no-op. Worst: the record lists only two of six ticket checkpoint commands.\n- reviews/shell-argument-ownership.md:251-371 verification covers sanitize and worktree only; the reviewer re-ran treetarget (54 pass) and the system tree-target test (pass). Confidence 5. no-op.\nRows S01-S04 and S15 are met; readers quoted at build_test.go:117,118,120,214, target_refusal_test.go:19, systemtest/tree_target_test.go:246.\n"
          },
          "axis": "Spec",
          "base": "1586949596564559881814a1695df97030c333ec",
          "tip": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "S-A-R1-coverage",
          "performer": "ft354-review-S-A",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3af663fa6facc743e27596274e13d17e252df1cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-A",
            "digest": "sha256:d3a9cbe3463f92fdcc9a19e1665fae00b3604330d180b07fd39758bf4aecb5a0",
            "excerpt": "## Coverage\nCount: 0.\n- Caller mutation: naive \"'\"+label+\"'\" in PrintTreeBuildRefusal bit TestTreeBuildArgument/an_embedded_quote, restored=yes.\n- Owner mutation: escape changed to the argv-equivalent '\"'\"' form bit 3 byte assertions, restored=yes.\n- The touched=false consumer internal/treetarget/run.go:65 passes the label unchanged.\n"
          },
          "axis": "Coverage",
          "base": "1586949596564559881814a1695df97030c333ec",
          "tip": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "S-B1",
      "base": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
      "tip": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
      "plan_digest": "sha256:74be021ade7410316adf7b5b52ba3e17eb821480cdd19ec6aad7a4f0bfbf7ad9",
      "source_digest": "fa087909f4ea9649d0094baa98f3ca249d38002a",
      "acceptance_rows": [
        "S13",
        "S40"
      ],
      "verification": [
        {
          "id": "S-B1-T02-omission-proof",
          "performer": "ft354-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fa087909f4ea9649d0094baa98f3ca249d38002a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-02",
            "digest": "sha256:9a6202ae101cdbf978c83efaf9b3bc2aff4a99c6db4ddc5049f43e4831be3795",
            "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap+swap,failed,2,TestSpecWidenRefusesAnOutOfScopeEntry,196,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/spec,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0,
          "probe": {
            "mutation": "Bypass the owner with former bare-token output at historyDetail or at the widen rerun line. The independent exact fixture must fail even if argv still round-trips.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-02",
              "digest": "sha256:9a6202ae101cdbf978c83efaf9b3bc2aff4a99c6db4ddc5049f43e4831be3795",
              "excerpt": "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap+swap,failed,2,TestSpecWidenRefusesAnOutOfScopeEntry,196,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/spec,all,passed,2\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "S-B1-T02-preservation",
          "performer": "ft354-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fa087909f4ea9649d0094baa98f3ca249d38002a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-02",
            "digest": "sha256:fc032360749126f97337e49a5841d828959da8cc32f781c84e8a1d3ac7f7f7b7",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestSelectedHistoryHostileTarget,196,1,0,null,yes\n  2,bit,swap,failed,2,TestSelectedSpecDetailRoute/--limit,196,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/spec,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
          "requirement": "preservation",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-02",
              "digest": "sha256:fc032360749126f97337e49a5841d828959da8cc32f781c84e8a1d3ac7f7f7b7",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestSelectedHistoryHostileTarget,196,1,0,null,yes\n  2,bit,swap,failed,2,TestSelectedSpecDetailRoute/--limit,196,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/spec,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "S-B1-R1-standards",
          "performer": "ft354-review-S-B1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fa087909f4ea9649d0094baa98f3ca249d38002a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B1",
            "digest": "sha256:b2e93aa1d837a5360b8e6bfa71f2518cd4eeee47a20c3c7d7993376b840a087d",
            "excerpt": "## Standards\nCount: 1, no-op. Worst: overlapping widen rerun expectations across two files.\n- internal/spec/widen_test.go:142-147 and internal/spec/widen_rerun_test.go:117 build the same fixture and assert the same rerun prefix; the widen_test.go block predates this change and the diff only updated its literal. Confidence 4. no-op.\n"
          },
          "axis": "Standards",
          "base": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
          "tip": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "S-B1-R1-spec",
          "performer": "ft354-review-S-B1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fa087909f4ea9649d0094baa98f3ca249d38002a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B1",
            "digest": "sha256:e20e76d6959c0c9a720e9d65d31f2ab488879794d221da91653a8eb1fbaa3eba",
            "excerpt": "## Spec\nCount: 0.\n- S13: history_selected.go:120 uses sanitize.ShellQuote; the -- marker (:117-118) and LineSafe (:64) remain.\n- S40: widen.go:129-131 quotes slug, basename, and entries; --ticket, --decision, and \u003creference> stay bare.\n"
          },
          "axis": "Spec",
          "base": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
          "tip": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "S-B1-R1-coverage",
          "performer": "ft354-review-S-B1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "fa087909f4ea9649d0094baa98f3ca249d38002a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B1",
            "digest": "sha256:b19a78f23baf3cbf66c78eefc5d1eab3aa341339b78c2adf9220f6bb32ec2f42",
            "excerpt": "## Coverage\nCount: 0.\n- Argv-equivalent escape at historyDetail bit TestSelectedSpecDetailRoute/quote'_name; restored=yes.\n- Argv-equivalent escape at the widen entries bit TestWidenRerunArguments; restored=yes.\n- Quoting \u003creference> and dropping the last entry each bit both widen tests; restored=yes.\n## Advice\n- The r1 placeholder substitution at widen_rerun_test.go:125 runs only after the exact-bytes assertion; the exact fixture carries the S40 detection.\n"
          },
          "axis": "Coverage",
          "base": "e9c03d6b5704fb9c4fb24c2cf1923effa6eaac8c",
          "tip": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "S-B2",
      "base": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
      "tip": "03bd8984525222136d809db96fa9038cf2524922",
      "plan_digest": "sha256:07216e7d8d911629ed14a562af6382f4a3a8345713d694278b971ea038f8134e",
      "source_digest": "d47637eae9b735d54f560c950e18a345db627151",
      "acceptance_rows": [
        "S14"
      ],
      "verification": [
        {
          "id": "S-B2-T03-omission-proof",
          "performer": "ft354-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ad5b5a96cd34ff88839539d6965cd29321df8534",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-03",
            "digest": "sha256:0ecade0736af8589b29af7280786872112f535948df2ae3013165a8be412914f",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  2,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  3,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  4,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/consumers,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "fb2f38e9c09ce5f5cc61b8e7cf691285a3b02001",
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0,
          "probe": {
            "mutation": "Return the former bare spelling at citation.cmd or omit a supplied flag. The exact fixture or recovered argv must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-03",
              "digest": "sha256:0ecade0736af8589b29af7280786872112f535948df2ae3013165a8be412914f",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  2,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  3,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  4,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/consumers,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "S-B2-T03-preservation",
          "performer": "ft354-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ad5b5a96cd34ff88839539d6965cd29321df8534",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-03",
            "digest": "sha256:0ecade0736af8589b29af7280786872112f535948df2ae3013165a8be412914f",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  2,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  3,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  4,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/consumers,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "fb2f38e9c09ce5f5cc61b8e7cf691285a3b02001",
          "requirement": "preservation",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-03",
              "digest": "sha256:0ecade0736af8589b29af7280786872112f535948df2ae3013165a8be412914f",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  2,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  3,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\n  4,bit,swap,failed,2,TestCitationCmdRoundTripsThroughAPOSIXShellSplit,66,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/consumers,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "S-B2-T03-R1-omission-proof",
          "performer": "ft354-repair-03-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d47637eae9b735d54f560c950e18a345db627151",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-repair-03-r1",
            "digest": "sha256:8260701e46c5ccdac6674ec1068e99d808768bf8b5954a0358c964dfe6ee632d",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+omit,failed,3,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/metacharacter_and_quote,68,1,0,null,yes\n  2,bit,swap,failed,3,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/metacharacter_and_quote,68,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/consumers,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "03bd8984525222136d809db96fa9038cf2524922",
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0,
          "probe": {
            "mutation": "Return the former bare spelling at citation.cmd or omit a supplied flag. The exact fixture or recovered argv must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-repair-03-r1",
              "digest": "sha256:8260701e46c5ccdac6674ec1068e99d808768bf8b5954a0358c964dfe6ee632d",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+omit,failed,3,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/metacharacter_and_quote,68,1,0,null,yes\n  2,bit,swap,failed,3,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/metacharacter_and_quote,68,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/consumers,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "S-B2-T03-R1-preservation",
          "performer": "ft354-repair-03-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d47637eae9b735d54f560c950e18a345db627151",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-repair-03-r1",
            "digest": "sha256:8260701e46c5ccdac6674ec1068e99d808768bf8b5954a0358c964dfe6ee632d",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+omit,failed,3,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/metacharacter_and_quote,68,1,0,null,yes\n  2,bit,swap,failed,3,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/metacharacter_and_quote,68,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/consumers,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "03bd8984525222136d809db96fa9038cf2524922",
          "requirement": "preservation",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-repair-03-r1",
              "digest": "sha256:8260701e46c5ccdac6674ec1068e99d808768bf8b5954a0358c964dfe6ee632d",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+omit,failed,3,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/metacharacter_and_quote,68,1,0,null,yes\n  2,bit,swap,failed,3,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/metacharacter_and_quote,68,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/consumers,all,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "S-B2-R1-standards",
          "performer": "ft354-review-S-B2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ad5b5a96cd34ff88839539d6965cd29321df8534",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B2",
            "digest": "sha256:97bfa5674035ff501e562c25969a02aef04c093dfd7c818768353aec757363b2",
            "excerpt": "## Standards\nCount: 0.\n- citation.go:58 uses sanitize.ShellQuote; the bench consumers prefix at :56 stays bare.\n- The removed shellSplit duplicated axitest.shellArgv (help.go:50); the fold removes a duplicated harness.\n## Advice\n- citationCells (citation_test.go:37) is a hand-written TOON row splitter beside axitest.DecodeDocument; it predates the diff.\n"
          },
          "axis": "Standards",
          "base": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
          "tip": "fb2f38e9c09ce5f5cc61b8e7cf691285a3b02001",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "S-B2-R1-spec",
          "performer": "ft354-review-S-B2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ad5b5a96cd34ff88839539d6965cd29321df8534",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B2",
            "digest": "sha256:b9202a81291b3d450d39e1b39080a9079f2cb305c077b8129555c2db0456844d",
            "excerpt": "## Spec\nCount: 1. Worst: the acceptance input class \"spaces\" was dropped on a premise that does not hold.\n- S-B2-SPEC-1: ticket 03 acceptance names spaces; citation_test.go:181 has no spaced operand. --base and --source-tip take git revisions (internal/git/git.go:259-261), and the reviewer ran bench consumers --changed --base \"HEAD^{/assign ticket 03}\" --source-tip HEAD at exit 0. Confidence 8. auto-fix (test-only, inside Writes). Implementation command contributed: an author who omits a named acceptance input class must show the producer refuses it.\n"
          },
          "axis": "Spec",
          "base": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
          "tip": "fb2f38e9c09ce5f5cc61b8e7cf691285a3b02001",
          "finding_ids": [
            "S-B2-SPEC-1"
          ],
          "supersedes": []
        },
        {
          "id": "S-B2-R1-coverage",
          "performer": "ft354-review-S-B2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ad5b5a96cd34ff88839539d6965cd29321df8534",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B2",
            "digest": "sha256:8696caaefbf9c6f519968b89f53c5307f9d02eab6d2abee63a9d25873bd53ef9",
            "excerpt": "## Coverage\nCount: 1. Worst: a space-splitting caller bypass survives.\n- S-B2-COV-1: at citation.go:58, strings.ReplaceAll(sanitize.ShellQuote(arg), \" \", \"' '\") reads silent (66 tests pass, restored=yes). Missing: a spaced-revision case with an exact literal and recovered argv in TestCitationCmdRoundTripsThroughAPOSIXShellSplit. Confidence 8. auto-fix.\n"
          },
          "axis": "Coverage",
          "base": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
          "tip": "fb2f38e9c09ce5f5cc61b8e7cf691285a3b02001",
          "finding_ids": [
            "S-B2-COV-1"
          ],
          "supersedes": []
        }
      ],
      "confirmations": [
        {
          "id": "S-B2-R1-confirmation",
          "performer": "ft354-orchestrator",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d47637eae9b735d54f560c950e18a345db627151",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:ft354-orchestrator",
            "digest": "sha256:5d43f9f70bc1c8028aee8c6a4af402570fb5b550a2b639f4fbbe000d419afa9e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/consumers/citation.go,insert,failed,1,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/space,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/consumers,all,passed,68\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,553\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,fail,1969,68\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/consumers,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/space,\"citation_test.go:208: cmd = \\\"bench consumers '--changed' '--base' 'feature;echo'\\\\\\\\\\\\\\\\''q^{/base|no' 'such subject}' '--source-tip' '410c33073fcb02ecbaaa118d25a304dc42e993e2' '--full'\\\", want \\\"bench consumers '--changed' '--base' 'feature;echo'\\\\\\\\… (337 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "base": "7cc8d6d041fa0dbc117d77c8c74beb51a2555bab",
          "reviewed_tip": "fb2f38e9c09ce5f5cc61b8e7cf691285a3b02001",
          "repaired_tip": "03bd8984525222136d809db96fa9038cf2524922",
          "finding_ids": [
            "S-B2-SPEC-1",
            "S-B2-COV-1"
          ],
          "supersedes": [],
          "probe": {
            "mutation": "insert before the citation.cmd Join return a split of the first space of each quoted token",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:ft354-orchestrator",
              "digest": "sha256:5d43f9f70bc1c8028aee8c6a4af402570fb5b550a2b639f4fbbe000d419afa9e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/consumers/citation.go,insert,failed,1,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/space,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/consumers,all,passed,68\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,553\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,fail,1969,68\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/consumers,TestCitationCmdRoundTripsThroughAPOSIXShellSplit/space,\"citation_test.go:208: cmd = \\\"bench consumers '--changed' '--base' 'feature;echo'\\\\\\\\\\\\\\\\''q^{/base|no' 'such subject}' '--source-tip' '410c33073fcb02ecbaaa118d25a304dc42e993e2' '--full'\\\", want \\\"bench consumers '--changed' '--base' 'feature;echo'\\\\\\\\… (337 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ]
    },
    {
      "id": "S-B4",
      "base": "03bd8984525222136d809db96fa9038cf2524922",
      "tip": "18517a9a7824394d85a6d262579584f76d60cf2d",
      "plan_digest": "sha256:cb192f02c39a015884882b7ad75e3747b70e693f02ac4763a1ef507632a07852",
      "source_digest": "88c0d3ced019a22c65fee32bb46078550dbdd808",
      "acceptance_rows": [
        "S17"
      ],
      "verification": [
        {
          "id": "S-B4-T05-omission-proof",
          "performer": "ft354-author-05",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de67355581646dd0c88b00cffd79357d0a5d4213",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-05",
            "digest": "sha256:3446816b9b58afa47933d0c5ff02a02cd4ba144f9cb95d0fc318ed5fe2117d66",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestEvidenceArgumentCommands/source_page,239,1,0,null,yes\n  2,bit,swap,failed,6,TestEvidenceArgumentCommands/preparation,239,1,0,null,yes\n  3,bit,omit,failed,4,TestEvidenceArgumentCommands/source_page,239,1,0,null,yes\n  4,bit,swap,failed,6,TestEvidenceArgumentCommands/preparation,239,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/evidencecmd,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "e93f34e3705456c6cd8fc9bdfad72c9b4218b733",
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/preflight/evidencecmd",
          "exit_code": 0,
          "probe": {
            "mutation": "Return a bare cursor or bypass the owner at evidenceInvocation. The exact producer fixture and argv oracle must detect each omission.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-05",
              "digest": "sha256:3446816b9b58afa47933d0c5ff02a02cd4ba144f9cb95d0fc318ed5fe2117d66",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestEvidenceArgumentCommands/source_page,239,1,0,null,yes\n  2,bit,swap,failed,6,TestEvidenceArgumentCommands/preparation,239,1,0,null,yes\n  3,bit,omit,failed,4,TestEvidenceArgumentCommands/source_page,239,1,0,null,yes\n  4,bit,swap,failed,6,TestEvidenceArgumentCommands/preparation,239,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/evidencecmd,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "S-B4-T05-preservation",
          "performer": "ft354-author-05",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de67355581646dd0c88b00cffd79357d0a5d4213",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-05",
            "digest": "sha256:3446816b9b58afa47933d0c5ff02a02cd4ba144f9cb95d0fc318ed5fe2117d66",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestEvidenceArgumentCommands/source_page,239,1,0,null,yes\n  2,bit,swap,failed,6,TestEvidenceArgumentCommands/preparation,239,1,0,null,yes\n  3,bit,omit,failed,4,TestEvidenceArgumentCommands/source_page,239,1,0,null,yes\n  4,bit,swap,failed,6,TestEvidenceArgumentCommands/preparation,239,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/evidencecmd,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "e93f34e3705456c6cd8fc9bdfad72c9b4218b733",
          "requirement": "preservation",
          "command": "bench test --package ./internal/preflight/evidencecmd",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-05",
              "digest": "sha256:3446816b9b58afa47933d0c5ff02a02cd4ba144f9cb95d0fc318ed5fe2117d66",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestEvidenceArgumentCommands/source_page,239,1,0,null,yes\n  2,bit,swap,failed,6,TestEvidenceArgumentCommands/preparation,239,1,0,null,yes\n  3,bit,omit,failed,4,TestEvidenceArgumentCommands/source_page,239,1,0,null,yes\n  4,bit,swap,failed,6,TestEvidenceArgumentCommands/preparation,239,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/evidencecmd,all,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "S-B4-T05-R1-omission-proof",
          "performer": "ft354-repair-05-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "88c0d3ced019a22c65fee32bb46078550dbdd808",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-repair-05-r1",
            "digest": "sha256:746a67ef1bb9edcadee2ca16f273ba7cf750dc90ca32e3978a5d1901b6b200dd",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestEvidenceArgumentCommands/source_page,7,1,0,null,yes\n  2,bit,swap,failed,3,TestEvidenceArgumentCommands/preparation,7,1,0,null,yes\n  3,bit,omit,failed,1,TestEvidenceArgumentCommands/source_page,7,1,0,null,yes\n  4,bit,swap,failed,3,TestEvidenceArgumentCommands/preparation,7,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/evidencecmd,^TestEvidenceArgumentCommands$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "18517a9a7824394d85a6d262579584f76d60cf2d",
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/preflight/evidencecmd",
          "exit_code": 0,
          "probe": {
            "mutation": "Return a bare cursor or bypass the owner at evidenceInvocation. The exact producer fixture and argv oracle must detect each omission.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-repair-05-r1",
              "digest": "sha256:746a67ef1bb9edcadee2ca16f273ba7cf750dc90ca32e3978a5d1901b6b200dd",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestEvidenceArgumentCommands/source_page,7,1,0,null,yes\n  2,bit,swap,failed,3,TestEvidenceArgumentCommands/preparation,7,1,0,null,yes\n  3,bit,omit,failed,1,TestEvidenceArgumentCommands/source_page,7,1,0,null,yes\n  4,bit,swap,failed,3,TestEvidenceArgumentCommands/preparation,7,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/evidencecmd,^TestEvidenceArgumentCommands$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "S-B4-T05-R1-preservation",
          "performer": "ft354-repair-05-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "88c0d3ced019a22c65fee32bb46078550dbdd808",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-repair-05-r1",
            "digest": "sha256:746a67ef1bb9edcadee2ca16f273ba7cf750dc90ca32e3978a5d1901b6b200dd",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestEvidenceArgumentCommands/source_page,7,1,0,null,yes\n  2,bit,swap,failed,3,TestEvidenceArgumentCommands/preparation,7,1,0,null,yes\n  3,bit,omit,failed,1,TestEvidenceArgumentCommands/source_page,7,1,0,null,yes\n  4,bit,swap,failed,3,TestEvidenceArgumentCommands/preparation,7,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/evidencecmd,^TestEvidenceArgumentCommands$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "18517a9a7824394d85a6d262579584f76d60cf2d",
          "requirement": "preservation",
          "command": "bench test --package ./internal/preflight/evidencecmd",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-repair-05-r1",
              "digest": "sha256:746a67ef1bb9edcadee2ca16f273ba7cf750dc90ca32e3978a5d1901b6b200dd",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestEvidenceArgumentCommands/source_page,7,1,0,null,yes\n  2,bit,swap,failed,3,TestEvidenceArgumentCommands/preparation,7,1,0,null,yes\n  3,bit,omit,failed,1,TestEvidenceArgumentCommands/source_page,7,1,0,null,yes\n  4,bit,swap,failed,3,TestEvidenceArgumentCommands/preparation,7,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/evidencecmd,^TestEvidenceArgumentCommands$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "S-B4-R1-standards",
          "performer": "ft354-review-S-B4",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de67355581646dd0c88b00cffd79357d0a5d4213",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B4",
            "digest": "sha256:bf1562047cb903f9ebe41bdec0c968d817a5a9fc69ec9352f9f1456acb8ac7ea",
            "excerpt": "## Standards\nCount: 1. Worst: S-B4-STD-1.\n- S-B4-STD-1: internal/preflight/evidencecmd/shell_arguments_test.go:20-21 claims that a lost quote or a split reds the recovered argv, but the producer renders only safe-alphabet values (chargeevidence/manifest.go:63,78), so only the exact-byte check at :52 can red. Confidence 8. auto-fix.\n"
          },
          "axis": "Standards",
          "base": "03bd8984525222136d809db96fa9038cf2524922",
          "tip": "e93f34e3705456c6cd8fc9bdfad72c9b4218b733",
          "finding_ids": [
            "S-B4-STD-1"
          ],
          "supersedes": []
        },
        {
          "id": "S-B4-R1-spec",
          "performer": "ft354-review-S-B4",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de67355581646dd0c88b00cffd79357d0a5d4213",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B4",
            "digest": "sha256:1574ddf751776774ae64e09af95f66b7d1eb350ae7761eba46f223fca96763f6",
            "excerpt": "## Spec\nCount: 1. Worst: S-B4-SPEC-1.\n- S-B4-SPEC-1: spec.md:125 \"Each rendered argv value becomes always quoted\" disagrees with the values-only spelling. Reviewer decision 2026-11-03: values only; the row is reworded. Confidence 5. ask-user, resolved by the reviewer.\nThe hostile-value refusal before render is verified at internal/preflight/command.go:65-82 and evidence.go:54-55,115-132.\n"
          },
          "axis": "Spec",
          "base": "03bd8984525222136d809db96fa9038cf2524922",
          "tip": "e93f34e3705456c6cd8fc9bdfad72c9b4218b733",
          "finding_ids": [
            "S-B4-SPEC-1"
          ],
          "supersedes": []
        },
        {
          "id": "S-B4-R1-coverage",
          "performer": "ft354-review-S-B4",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de67355581646dd0c88b00cffd79357d0a5d4213",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B4",
            "digest": "sha256:b645bea5b91102275db336c3fc0de30a985300bef87e55bd7694769f0dd2bb38",
            "excerpt": "## Coverage\nCount: 1. Worst: S-B4-COV-1.\n- S-B4-COV-1: the sentinel check at shell_arguments_test.go:76 cannot red, because the hostile loop (:61-78) passes no output to a shell; the comment at :22-23 claims it as a guard. Confidence 8. auto-fix.\n## Advice\n- nextArgs and commandArgs are thin wrappers over axitest.HelpAction.Argv in two packages; incidental repetition.\n"
          },
          "axis": "Coverage",
          "base": "03bd8984525222136d809db96fa9038cf2524922",
          "tip": "e93f34e3705456c6cd8fc9bdfad72c9b4218b733",
          "finding_ids": [
            "S-B4-COV-1"
          ],
          "supersedes": []
        },
        {
          "id": "S-B4-R2-standards",
          "performer": "ft354-review-S-B4-c1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "88c0d3ced019a22c65fee32bb46078550dbdd808",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B4-c1",
            "digest": "sha256:e7032b8ae3ccaf815e94f3a18560894567cea2e5caf7a94415707a29be309594",
            "excerpt": "## Standards\nCount: 0 open.\n- S-B4-STD-1 closed: shell_arguments_test.go:15-17 credits the exact bytes with the lost-quote catch (:48) and the shell read-back with replay (:51-53). Confidence 9.\n"
          },
          "axis": "Standards",
          "base": "03bd8984525222136d809db96fa9038cf2524922",
          "tip": "18517a9a7824394d85a6d262579584f76d60cf2d",
          "finding_ids": [],
          "supersedes": [
            "S-B4-R1-standards"
          ]
        },
        {
          "id": "S-B4-R2-spec",
          "performer": "ft354-review-S-B4-c1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "88c0d3ced019a22c65fee32bb46078550dbdd808",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B4-c1",
            "digest": "sha256:f2570b4baea561e9b1172210c4076242c5d9fe074802557a4b09d717a3b69207",
            "excerpt": "## Spec\nCount: 0 open.\n- S-B4-SPEC-1 closed: spec.md:125 names the identity, source, and cursor values and the reviewer decision; the fixtures at shell_arguments_test.go:38-45 match. Confidence 9.\n"
          },
          "axis": "Spec",
          "base": "03bd8984525222136d809db96fa9038cf2524922",
          "tip": "18517a9a7824394d85a6d262579584f76d60cf2d",
          "finding_ids": [],
          "supersedes": [
            "S-B4-R1-spec"
          ]
        },
        {
          "id": "S-B4-R2-coverage",
          "performer": "ft354-review-S-B4-c1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "88c0d3ced019a22c65fee32bb46078550dbdd808",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-B4-c1",
            "digest": "sha256:f0481cd0c2653781fecce8c16d3c13a5d1be28dbdcfa6f2ef758dbe029377b07",
            "excerpt": "## Coverage\nCount: 0 open.\n- S-B4-COV-1 closed: the dead sentinel and its imports are removed; the hostile loop (:57-70) keeps its exit-code and no-next assertions. Confidence 9.\n"
          },
          "axis": "Coverage",
          "base": "03bd8984525222136d809db96fa9038cf2524922",
          "tip": "18517a9a7824394d85a6d262579584f76d60cf2d",
          "finding_ids": [],
          "supersedes": [
            "S-B4-R1-coverage"
          ]
        }
      ]
    },
    {
      "id": "S-D1",
      "base": "4ac0f5b42b6d39a58999544e3b7835a3d40b3009",
      "tip": "79c12fb863f1429d46033be7f4c7c86a2c8b034f",
      "plan_digest": "sha256:ba5cd9dee92e0db74c09fc55407444b2d0c5e1f094ca6df0cd67710c797db240",
      "source_digest": "dda393fd3e2d4e108d4f826ddfd6ad9a7f176e31",
      "acceptance_rows": [
        "S18",
        "S37"
      ],
      "verification": [
        {
          "id": "S-D1-T06-preservation",
          "performer": "ft354-author-06",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dda393fd3e2d4e108d4f826ddfd6ad9a7f176e31",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-06",
            "digest": "sha256:d55a1b07ea3065d2a7faab7cf7ad91b797dd1ebb321a8be6aaaa99370491a190",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,all,passed,95\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,546\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,10213,95\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths,\"freshness_verify_test.go:267: rebuild action args = \\\"/dev/shm/5O3DCD/t/TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths1044644491/001/root [*] with ' quote\\\\\\\\n/dev/shm/5O3DCD/t/TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths10446… (633 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "79c12fb863f1429d46033be7f4c7c86a2c8b034f",
          "requirement": "preservation",
          "command": "bench test --package ./internal/freshness",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-06",
              "digest": "sha256:d55a1b07ea3065d2a7faab7cf7ad91b797dd1ebb321a8be6aaaa99370491a190",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,all,passed,95\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,546\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,10213,95\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths,\"freshness_verify_test.go:267: rebuild action args = \\\"/dev/shm/5O3DCD/t/TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths1044644491/001/root [*] with ' quote\\\\\\\\n/dev/shm/5O3DCD/t/TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths10446… (633 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "S-D1-T06-rebuild-argument-proof",
          "performer": "ft354-author-06",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dda393fd3e2d4e108d4f826ddfd6ad9a7f176e31",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-06",
            "digest": "sha256:ce6e4b5a586e2aefbfef36dad4d253935fe7b893146370bba2108fdcbdea5638",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,all,passed,95\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,546\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,10780,95\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths,\"freshness_verify_test.go:267: rebuild action args = \\\"/dev/shm/3CIPD5/t/TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths3421749207/001/root [*] with ' quote\\\\\\\\n/dev/shm/3CIPD5/t/TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths34217… (723 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "79c12fb863f1429d46033be7f4c7c86a2c8b034f",
          "requirement": "rebuild-argument-proof",
          "command": "bench test --package ./internal/freshness",
          "exit_code": 0,
          "probe": {
            "mutation": "Emit a malformed quoted rebuild operand or drop a supplied rebuild operand. Its real producer fixture and independent shell argv assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-06",
              "digest": "sha256:ce6e4b5a586e2aefbfef36dad4d253935fe7b893146370bba2108fdcbdea5638",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/freshness/freshness_verify.go,swap,failed,1,TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/freshness,all,passed,95\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,546\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/freshness,fail,10780,95\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/freshness,TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths,\"freshness_verify_test.go:267: rebuild action args = \\\"/dev/shm/3CIPD5/t/TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths3421749207/001/root [*] with ' quote\\\\\\\\n/dev/shm/3CIPD5/t/TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths34217… (723 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "S-D1-T06-diff-fixture-argument-proof",
          "performer": "ft354-author-06",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dda393fd3e2d4e108d4f826ddfd6ad9a7f176e31",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-06",
            "digest": "sha256:2eb0921eda1b8e7fda4b9b6a7834a44cace79afc0ac263f864b55ec5e5ca68b3",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/diff/identity_test.go,swap,failed,1,TestCommandResolvesRecordedBaseOncePerAttempt,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/diff,all,passed,100\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,395\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,fail,5645,100\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/diff,TestCommandResolvesRecordedBaseOncePerAttempt,\"identity_test.go:210: recorded-base reads = 0, want one resolution plus one after-read capture per each of two attempts; trace:\",2\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "79c12fb863f1429d46033be7f4c7c86a2c8b034f",
          "requirement": "diff-fixture-argument-proof",
          "command": "bench test --package ./internal/diff",
          "exit_code": 0,
          "probe": {
            "mutation": "Corrupt quote escaping for a hostile diff fixture operand. Its independent observed argv assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-06",
              "digest": "sha256:2eb0921eda1b8e7fda4b9b6a7834a44cace79afc0ac263f864b55ec5e5ca68b3",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/diff/identity_test.go,swap,failed,1,TestCommandResolvesRecordedBaseOncePerAttempt,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/diff,all,passed,100\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,395\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,fail,5645,100\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/diff,TestCommandResolvesRecordedBaseOncePerAttempt,\"identity_test.go:210: recorded-base reads = 0, want one resolution plus one after-read capture per each of two attempts; trace:\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "S-D1-T06-release-fixture-argument-proof",
          "performer": "ft354-author-06",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dda393fd3e2d4e108d4f826ddfd6ad9a7f176e31",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-06",
            "digest": "sha256:63ba74a8f3a593d0a119c26ac4bd57e328c0edf623bcdca5408ed6ce4b2b15f7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preprelease/preprelease_test.go,swap,failed,1,TestCanaryCompatibilityShimAttributesInventoryFailures/resolved_Bench,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preprelease,all,passed,15\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,401\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preprelease,fail,2030,15\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preprelease,TestCanaryCompatibilityShimAttributesInventoryFailures/resolved_Bench,\"preprelease_test.go:145: resolved Bench argv = \\\"canary \\\\\\\\n\\\", \u003cnil>; want canary fixture root\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "79c12fb863f1429d46033be7f4c7c86a2c8b034f",
          "requirement": "release-fixture-argument-proof",
          "command": "bench test --package ./internal/preprelease",
          "exit_code": 0,
          "probe": {
            "mutation": "Corrupt quote escaping for a hostile preprelease fixture operand. Its independent observed argv assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-06",
              "digest": "sha256:63ba74a8f3a593d0a119c26ac4bd57e328c0edf623bcdca5408ed6ce4b2b15f7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preprelease/preprelease_test.go,swap,failed,1,TestCanaryCompatibilityShimAttributesInventoryFailures/resolved_Bench,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preprelease,all,passed,15\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,401\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preprelease,fail,2030,15\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preprelease,TestCanaryCompatibilityShimAttributesInventoryFailures/resolved_Bench,\"preprelease_test.go:145: resolved Bench argv = \\\"canary \\\\\\\\n\\\", \u003cnil>; want canary fixture root\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "S-D1-R1-standards",
          "performer": "ft354-review-S-D1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dda393fd3e2d4e108d4f826ddfd6ad9a7f176e31",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-D1",
            "digest": "sha256:dc2bd303d3904cbea7e582d4508abf1f72e80917130e5e7b92703c9ad93702f6",
            "excerpt": "## Standards\nCount: 0. Worst issue: none.\n- The diff removes freshness.shellQuote, freshness quoteForShell, the diff test shellQuote, and the preprelease shellQuote; it adds no duplicated knowledge.\n"
          },
          "axis": "Standards",
          "base": "4ac0f5b42b6d39a58999544e3b7835a3d40b3009",
          "tip": "79c12fb863f1429d46033be7f4c7c86a2c8b034f",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "S-D1-R1-spec",
          "performer": "ft354-review-S-D1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dda393fd3e2d4e108d4f826ddfd6ad9a7f176e31",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-D1",
            "digest": "sha256:01907644dd38cb13324c1a91d7a85fb0cb293dcb94249d2e61891b530f6473b9",
            "excerpt": "## Spec\nCount: 1. Worst issue: S-D1-SPEC-1.\n- S-D1-SPEC-1 (confidence 8, no-op): spec.md:530 and ticket 06 line 22 predict a preprelease script-text change, but no fixture input holds a quote (preprelease_test.go:162-163), so the escape never runs. Non-behavioral premise error; divergence 2 accepted.\n"
          },
          "axis": "Spec",
          "base": "4ac0f5b42b6d39a58999544e3b7835a3d40b3009",
          "tip": "79c12fb863f1429d46033be7f4c7c86a2c8b034f",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "S-D1-R1-coverage",
          "performer": "ft354-review-S-D1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "dda393fd3e2d4e108d4f826ddfd6ad9a7f176e31",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-review-S-D1",
            "digest": "sha256:4475439bbc3a2ca506cb36b05322792b2a52178a06a4503b93ae60f61e9f9c07",
            "excerpt": "## Coverage\nCount: 1. Worst issue: S-D1-COV-1.\n- S-D1-COV-1 (confidence 6, no-op): the preprelease fixture never exercises quote escaping; a no-escape wrap of the shim stays green. TestShellQuoteRoundTrip and the diff and freshness hostile fixtures prove the escape, so no repair. Divergence 1 justified.\n"
          },
          "axis": "Coverage",
          "base": "4ac0f5b42b6d39a58999544e3b7835a3d40b3009",
          "tip": "79c12fb863f1429d46033be7f4c7c86a2c8b034f",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "S-E",
      "base": "79c12fb863f1429d46033be7f4c7c86a2c8b034f",
      "tip": "80356fc1252d49563f727e91b417482f28320e55",
      "plan_digest": "sha256:c20f6c65e343e0a2334a915e1df675c951ffcad45fc3ce413f0cd12dd2a0d76e",
      "source_digest": "097cae6c5d08f3b33c93e89abdf49865d12bf8b7",
      "acceptance_rows": [
        "S16",
        "S05",
        "S06",
        "S07",
        "S08",
        "S09",
        "S10",
        "S11",
        "S12",
        "S19",
        "S20",
        "S21",
        "S22",
        "S23",
        "S24",
        "S25",
        "S26",
        "S27",
        "S28",
        "S29",
        "S30",
        "S31",
        "S32",
        "S33",
        "S34",
        "S35",
        "S36",
        "S38",
        "S39",
        "S41",
        "S42",
        "S43",
        "S44"
      ],
      "verification": [
        {
          "id": "S-E-T04-recovery-line-agreement",
          "performer": "ft354-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "097cae6c5d08f3b33c93e89abdf49865d12bf8b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft354-author-04",
            "digest": "sha256:41b5dc50c874a42ebfb2ace48a41dea20475b0df86ab71da8049ec002a238bcb",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/list.go,omit,failed,7,TestRecoveryRouteArgument/a_command_list/discard,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,RecoveryRouteArgument|MissingTreeRoute|UnlandedRelease|IdentifierOperand,passed,20\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2668\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,193,20\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRecoveryRouteArgument/a_command_list/discard,\"shell_arguments_test.go:89: recovery line = \\\"bench worktree clean '/gone/a; touch /dev/shm/SL67ZH/t/TestRecoveryRouteArgument1822496834/002/executed'\\\" and help commands = [\\\"bench worktree clean '/gone/a; touch /dev/shm/SL67ZH/t/TestRecovery… (539 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRecoveryRouteArgument/a_command_substitution/discard,\"shell_arguments_test.go:89: recovery line = \\\"bench worktree clean '/gone/$(touch /dev/shm/SL67ZH/t/TestRecoveryRouteArgument1822496834/002/executed)'\\\" and help commands = [\\\"bench worktree clean '/gone/$(touch /dev/shm/SL67ZH/t/TestRecoveryR… (539 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRecoveryRouteArgument/an_angle_bracket/discard,\"shell_arguments_test.go:89: recovery line = \\\"bench worktree clean '/gone/a\u003cb>'\\\" and help commands = [], want \\\"bench worktree clean --discard-branch '/gone/a\u003cb>'\\\" and []\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRecoveryRouteArgument/an_embedded_quote/discard,\"shell_arguments_test.go:89: recovery line = \\\"bench worktree clean '/gone/it'\\\\\\\\\\\"'\\\\\\\\\\\"'s'\\\" and help commands = [\\\"bench worktree clean '/gone/it'\\\\\\\\\\\"'\\\\\\\\\\\"'s'\\\"], want \\\"bench worktree clean --discard-branch '/gone/it'\\\\\\\\\\\"'\\\\\\\\\\\"'s'\\\" and [\\\"bench worktree clea… (279 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRecoveryRouteArgument/spaces/discard,\"shell_arguments_test.go:89: recovery line = \\\"bench worktree clean '/gone/my  tree '\\\" and help commands = [\\\"bench worktree clean '/gone/my  tree '\\\"], want \\\"bench worktree clean --discard-branch '/gone/my  tree '\\\" and [\\\"bench worktree clean -… (275 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnlandedReleaseRefusalNamesTheDiscardClean/no_own_commit_with_output,\"unlanded_route_test.go:70: release refusal = \\\"bench worktree release: worktree retained (unmerged): assignment branch has not landed; next=bench worktree clean /dev/shm/SL67ZH/t/TestUnlandedReleaseRefusalNamesTheDiscardCleanno_own_commit_wi… (611 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnlandedReleaseRefusalNamesTheDiscardClean/unlanded_commit,\"unlanded_route_test.go:70: release refusal = \\\"bench worktree release: worktree retained (unmerged): assignment branch has not landed; next=bench worktree clean /dev/shm/SL67ZH/t/TestUnlandedReleaseRefusalNamesTheDiscardCleanunlanded_commit4… (607 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "80356fc1252d49563f727e91b417482f28320e55",
          "requirement": "recovery-line-agreement",
          "command": "bench test --package ./internal/worktree --run 'RecoveryRouteArgument|MissingTreeRoute|UnlandedRelease|IdentifierOperand'",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the path from the shared recovery argument list, or drop a route word. The real producer fixture and the line-and-help-row agreement assertion must fail, then pass after restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft354-author-04",
              "digest": "sha256:41b5dc50c874a42ebfb2ace48a41dea20475b0df86ab71da8049ec002a238bcb",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/worktree/list.go,omit,failed,7,TestRecoveryRouteArgument/a_command_list/discard,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,RecoveryRouteArgument|MissingTreeRoute|UnlandedRelease|IdentifierOperand,passed,20\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2668\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,193,20\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestRecoveryRouteArgument/a_command_list/discard,\"shell_arguments_test.go:89: recovery line = \\\"bench worktree clean '/gone/a; touch /dev/shm/SL67ZH/t/TestRecoveryRouteArgument1822496834/002/executed'\\\" and help commands = [\\\"bench worktree clean '/gone/a; touch /dev/shm/SL67ZH/t/TestRecovery… (539 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRecoveryRouteArgument/a_command_substitution/discard,\"shell_arguments_test.go:89: recovery line = \\\"bench worktree clean '/gone/$(touch /dev/shm/SL67ZH/t/TestRecoveryRouteArgument1822496834/002/executed)'\\\" and help commands = [\\\"bench worktree clean '/gone/$(touch /dev/shm/SL67ZH/t/TestRecoveryR… (539 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRecoveryRouteArgument/an_angle_bracket/discard,\"shell_arguments_test.go:89: recovery line = \\\"bench worktree clean '/gone/a\u003cb>'\\\" and help commands = [], want \\\"bench worktree clean --discard-branch '/gone/a\u003cb>'\\\" and []\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRecoveryRouteArgument/an_embedded_quote/discard,\"shell_arguments_test.go:89: recovery line = \\\"bench worktree clean '/gone/it'\\\\\\\\\\\"'\\\\\\\\\\\"'s'\\\" and help commands = [\\\"bench worktree clean '/gone/it'\\\\\\\\\\\"'\\\\\\\\\\\"'s'\\\"], want \\\"bench worktree clean --discard-branch '/gone/it'\\\\\\\\\\\"'\\\\\\\\\\\"'s'\\\" and [\\\"bench worktree clea… (279 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestRecoveryRouteArgument/spaces/discard,\"shell_arguments_test.go:89: recovery line = \\\"bench worktree clean '/gone/my  tree '\\\" and help commands = [\\\"bench worktree clean '/gone/my  tree '\\\"], want \\\"bench worktree clean --discard-branch '/gone/my  tree '\\\" and [\\\"bench worktree clean -… (275 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnlandedReleaseRefusalNamesTheDiscardClean/no_own_commit_with_output,\"unlanded_route_test.go:70: release refusal = \\\"bench worktree release: worktree retained (unmerged): assignment branch has not landed; next=bench worktree clean /dev/shm/SL67ZH/t/TestUnlandedReleaseRefusalNamesTheDiscardCleanno_own_commit_wi… (611 bytes)\",1\n  github.com/gibbonmi/bench/internal/worktree,TestUnlandedReleaseRefusalNamesTheDiscardClean/unlanded_commit,\"unlanded_route_test.go:70: release refusal = \\\"bench worktree release: worktree retained (unmerged): assignment branch has not landed; next=bench worktree clean /dev/shm/SL67ZH/t/TestUnlandedReleaseRefusalNamesTheDiscardCleanunlanded_commit4… (607 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
      "from": "sha256:dd10c1e33a504fed07206dd2621197369959e9b0cc4f50eaa6fef622276f7818",
      "to": "sha256:74be021ade7410316adf7b5b52ba3e17eb821480cdd19ec6aad7a4f0bfbf7ad9",
      "chunk_ids": {
        "S-A": [
          "S-A"
        ]
      }
    },
    {
      "from": "sha256:74be021ade7410316adf7b5b52ba3e17eb821480cdd19ec6aad7a4f0bfbf7ad9",
      "to": "sha256:1a12724ba4cfce35272a0d678e74da28c1130d4ade56e7dd1a8f2c0c25fdf235",
      "chunk_ids": {
        "S-A": [
          "S-A"
        ],
        "S-B1": [
          "S-B1"
        ]
      }
    },
    {
      "from": "sha256:1a12724ba4cfce35272a0d678e74da28c1130d4ade56e7dd1a8f2c0c25fdf235",
      "to": "sha256:07216e7d8d911629ed14a562af6382f4a3a8345713d694278b971ea038f8134e",
      "chunk_ids": {
        "S-A": [
          "S-A"
        ],
        "S-B1": [
          "S-B1"
        ],
        "S-B2": [
          "S-B2"
        ]
      }
    },
    {
      "from": "sha256:07216e7d8d911629ed14a562af6382f4a3a8345713d694278b971ea038f8134e",
      "to": "sha256:13292e498640b06e3de0bd8d00eb08d63e38c71a9b09efcbd62ea4757f71ec98",
      "chunk_ids": {
        "S-A": [
          "S-A"
        ],
        "S-B1": [
          "S-B1"
        ],
        "S-B2": [
          "S-B2"
        ]
      }
    },
    {
      "from": "sha256:13292e498640b06e3de0bd8d00eb08d63e38c71a9b09efcbd62ea4757f71ec98",
      "to": "sha256:cb192f02c39a015884882b7ad75e3747b70e693f02ac4763a1ef507632a07852",
      "chunk_ids": {
        "S-A": [
          "S-A"
        ],
        "S-B1": [
          "S-B1"
        ],
        "S-B2": [
          "S-B2"
        ],
        "S-B4": [
          "S-B4"
        ]
      }
    },
    {
      "from": "sha256:cb192f02c39a015884882b7ad75e3747b70e693f02ac4763a1ef507632a07852",
      "to": "sha256:ba5cd9dee92e0db74c09fc55407444b2d0c5e1f094ca6df0cd67710c797db240",
      "chunk_ids": {
        "S-A": [
          "S-A"
        ],
        "S-B1": [
          "S-B1"
        ],
        "S-B2": [
          "S-B2"
        ],
        "S-B4": [
          "S-B4"
        ]
      }
    },
    {
      "from": "sha256:ba5cd9dee92e0db74c09fc55407444b2d0c5e1f094ca6df0cd67710c797db240",
      "to": "sha256:c20f6c65e343e0a2334a915e1df675c951ffcad45fc3ce413f0cd12dd2a0d76e",
      "chunk_ids": {
        "S-A": [
          "S-A"
        ],
        "S-B1": [
          "S-B1"
        ],
        "S-B2": [
          "S-B2"
        ],
        "S-B4": [
          "S-B4"
        ],
        "S-D1": [
          "S-D1"
        ]
      }
    }
  ]
}
```

## S-A author witness reds

The ticket 01 author reported these witness reds. The S01 to S04 reds came from probe mutations of the owner, and the S15 red came before the production edit. Each mutation was restored, and the packages then passed.

| Row | Mutation | Red test |
|---|---|---|
| S01 | The empty quoted pair is returned as an empty string. | `TestShellQuoteRoundTrip/empty` |
| S02 | The embedded quote escape is replaced with a split. | `TestShellQuoteRoundTrip/a_leading_quote` |
| S03 | The owner returns the bare value. | `TestShellQuoteRoundTrip`, 13 subtest failures |
| S04 | A safe-token shortcut returns a safe value bare. | `TestShellQuoteAlwaysQuoted/a-b_c.d/e` |
| S15 | The producer still calls `axi.ShellQuote` before the production edit. | `TestTreeBuildArgument/an_embedded_quote`, `TestTargetRefusalEscapesOnlyAnUnsafeLine/a_plain_label` |

The author reported one limit: the S03 probe did not isolate the sentinel-file assertion from the argv and byte failures.

## S-B1 author witness reds

The ticket 02 author reported these witness reds before the first production edit. The packages passed after the migration.

| Row | Mutation | Red test |
|---|---|---|
| S13 | The history detail operand still renders through the bare-token policy. | `TestSelectedSpecDetailRoute` at `history_selected_test.go:140`, for `--limit`, `x.md.md`, and `quote' name` |
| S40 | The widen rerun line still renders through the bare-token policy. | `TestWidenRerunArguments`, `TestSpecWidenRefusesAnOutOfScopeEntry` |

The S40 argv recovery replaces the `<reference>` placeholder with `r1` in its recovery input only, because a shell reads a bare `<reference>` as a redirection.

## S-B2 author witness reds

The ticket 03 author reported this witness red before the production edit. The package passed after the migration.

| Row | Mutation | Red test |
|---|---|---|
| S14 | The citation cmd cell still renders through the bare-token policy. | `TestCitationCmdRoundTripsThroughAPOSIXShellSplit` at `citation_test.go:200`, `TestCitationNamesHeadVersionAndReplaySpelling` at `citation_test.go:166` |

The author folded the package-local `shellSplit` harness into `axitest.HelpAction.Argv`. A git ref cannot hold a space, so the S14 fixture uses an embedded quote and a metacharacter, and the owner tests cover spaces.

## S-B2 review round 1 pickup

Repair allowance: S-B2 consumed 0 of 2 repair cycles before this round.

### Standards

Count: 0. Worst issue: none.

### Spec

Count: 1. Worst issue: S-B2-SPEC-1.

- S-B2-SPEC-1 (confidence 8, auto-fix): the ticket 03 acceptance names spaces, but the S14 fixture has no spaced operand. A `--base` or `--source-tip` revision can hold a space, for example `HEAD^{/assign ticket 03}`. Citation: `internal/consumers/citation_test.go:181`, `internal/git/git.go:259`.

### Coverage

Count: 1. Worst issue: S-B2-COV-1.

- S-B2-COV-1 (confidence 8, auto-fix): a caller bypass that splits a spaced argument into two shell words stays green. Citation: `internal/consumers/citation.go:58`.

Both findings name one test-only repair on ticket 03. A fresh repair session takes it.

## S-B2 repair 1 witness reds

The repair session ft354-repair-03-r1 reported this witness red. The coordinator probe also bit at a different site and with a different mutation kind.

| Source | Mutation | Red test |
|---|---|---|
| Repair witness | At `citation.go:58`, each quoted argument splits on every space. | `TestCitationCmdRoundTripsThroughAPOSIXShellSplit/space` |
| Coordinator probe | Before the `Join` return, the first space of each token splits. | `TestCitationCmdRoundTripsThroughAPOSIXShellSplit/space` |

Repair allowance: S-B2 consumed 1 of 2 repair cycles.

## S-B4 author witness reds

The ticket 05 author reported this witness red before the production edit. The package passed after the migration.

| Row | Mutation | Red test |
|---|---|---|
| S17 | The evidence next command still renders through the bare-token policy. | `TestEvidenceArgumentCommands/preparation`, `/summary`, and `/source_page` at `shell_arguments_test.go:53` |

The producer refuses a hostile identity, source, or cursor before any command renders. Thus the S17 test shows those refusals.
The orchestrator decision, flagged for reviewer veto: only the identity, source, and cursor values take `sanitize.ShellQuote`, and the fixed tokens stay bare.

## S-B4 review round 1 pickup

Repair allowance: S-B4 consumed 0 of 2 repair cycles before this round.

### Standards

Count: 1. Worst issue: S-B4-STD-1.

- S-B4-STD-1 (confidence 8, auto-fix): the test comment claims an argv oracle that cannot red for safe-alphabet values. Citation: `internal/preflight/evidencecmd/shell_arguments_test.go:20`.

### Spec

Count: 1. Worst issue: S-B4-SPEC-1.

- S-B4-SPEC-1 (confidence 5, ask-user): the direct inventory row read as quoting every word. The reviewer decided values only on 2026-11-03, and the row now names the identity, source, and cursor values. Citation: `specs/shell-argument-ownership/spec.md:125`.

### Coverage

Count: 1. Worst issue: S-B4-COV-1.

- S-B4-COV-1 (confidence 8, auto-fix): the sentinel check of the hostile loop cannot red, because no output reaches a shell there. Citation: `internal/preflight/evidencecmd/shell_arguments_test.go:76`.

S-B4-STD-1 and S-B4-COV-1 name one test-only repair on ticket 05. A fresh repair session takes it.

## S-B4 repair 1 confirmation

The repair session ft354-repair-05-r1 added no expectation. It rewrote the oracle comment and removed the sentinel check that could not red.
The coordinator probe omitted the source validation in `operandRefusal`, and `TestEvidenceArgumentCommands/source_with_an_embedded_quote` went red. The hostile refusal rows therefore still bite.
The spec reword makes the delta more than test-only, so a three-axis confirming round closes the findings.

Repair allowance: S-B4 consumed 1 of 2 repair cycles.

## S-D1 author witness reds

The ticket 06 author reported no witness red apart from the plan probes. The fenced record holds those four probes, and each probe went red and then passed after restoration.

| Row | Mutation | Red test |
|---|---|---|
| S18 | The second rebuild operand takes a malformed quote around its embedded quote. | `TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths` at `freshness_verify_test.go:267` |
| S37 | The diff fixture git path takes a malformed quote. | `TestCommandResolvesRecordedBaseOncePerAttempt` at `identity_test.go:210` |
| S37 | The preprelease fixture root operand loses its quotes. | `TestCanaryCompatibilityShimAttributesInventoryFailures/resolved_Bench` at `preprelease_test.go:145` |

The author added hostile fixture inputs so that the omission evidence can red. The diff wrapper runs git through a path with a space, a glob, a dollar sign, and a quote. The preprelease script builds its root operand through `sanitize.ShellQuote`.

Two divergences are flagged for reviewer veto:

- The preprelease operand contains no embedded quote, so its probe removes the quotes and does not corrupt the escape. An embedded quote would change the observed argv bytes, and the acceptance forbids that change.
- The preprelease script text changes only where the owner now builds the root operand. The removed helper's quote escape never reached the shim path, so those bytes stay the same.

## S-D1 review round 1 pickup

Repair allowance: S-D1 consumed 0 of 2 repair cycles. The round has no repair target, so S-D1 closes at its frozen tip.
The frozen pair starts at the S-B4 source tip `18517a9a7`. The chunk entry names the record commit `4ac0f5b42`, and the two bases give the same code delta.

### Standards

Count: 0. Worst issue: none.

### Spec

Count: 1. Worst issue: S-D1-SPEC-1.

- S-D1-SPEC-1 (confidence 8, no-op): the spec and ticket 06 predict a change of the preprelease script text, but no fixture input holds a quote. The orchestrator accepts divergence 2, and the next plan commit corrects the prose. Citation: `specs/shell-argument-ownership/spec.md:530`.

### Coverage

Count: 1. Worst issue: S-D1-COV-1.

- S-D1-COV-1 (confidence 6, no-op): the preprelease fixture never exercises quote escaping. The owner round-trip test and the hostile diff and freshness fixtures prove the escape, so no repair is necessary. The orchestrator accepts divergence 1. Citation: `internal/preprelease/preprelease_test.go:162`.

The definition census found two surviving quote algorithms: the owner and `axi.ShellQuote`, which ticket 07 owns in S-E.
