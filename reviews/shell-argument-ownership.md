# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/shell-argument-ownership/spec.md",
  "plan_digest": "sha256:74be021ade7410316adf7b5b52ba3e17eb821480cdd19ec6aad7a4f0bfbf7ad9",
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
