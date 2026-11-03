# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/shell-argument-ownership/spec.md",
  "plan_digest": "sha256:dd10c1e33a504fed07206dd2621197369959e9b0cc4f50eaa6fef622276f7818",
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
