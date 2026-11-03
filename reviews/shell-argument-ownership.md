# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/shell-argument-ownership/spec.md",
  "plan_digest": "sha256:13292e498640b06e3de0bd8d00eb08d63e38c71a9b09efcbd62ea4757f71ec98",
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
      "tip": "e93f34e3705456c6cd8fc9bdfad72c9b4218b733",
      "plan_digest": "sha256:13292e498640b06e3de0bd8d00eb08d63e38c71a9b09efcbd62ea4757f71ec98",
      "source_digest": "de67355581646dd0c88b00cffd79357d0a5d4213",
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

The producer refuses a hostile identity, source, or cursor before any command renders. Thus the S17 test shows those refusals and a sentinel absence.
The orchestrator decision, flagged for reviewer veto: only the identity, source, and cursor values take `sanitize.ShellQuote`, and the fixed tokens stay bare.
