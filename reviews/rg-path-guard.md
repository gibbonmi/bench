# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/rg-path-guard/spec.md",
  "plan_digest": "sha256:661349ba0ae7d631cf0c2b43e02bb65b60aefca056311d23ffec0b8190db0fc9",
  "implementation_session": "",
  "chunks": [
    {
      "id": "C1",
      "base": "7cc35cbbf401ead2a5b374454e7441a81ec884d0",
      "tip": "eceaeaceb006462864ea135e751b1698768f6cfc",
      "plan_digest": "sha256:661349ba0ae7d631cf0c2b43e02bb65b60aefca056311d23ffec0b8190db0fc9",
      "source_digest": "0f739eaa220393c4db16c40462326410c8702b5d",
      "acceptance_rows": [
        "R01",
        "R02",
        "R03",
        "R04",
        "R05",
        "R06",
        "R07",
        "R55"
      ],
      "verification": [
        {
          "id": "c1-t01-mark-tests",
          "performer": "ft366-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "0f739eaa220393c4db16c40462326410c8702b5d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft366-author-01",
            "digest": "sha256:223c7d2e9b798089b676d00358709a2d5ae99d667b6df10b5a99ff40138cf296",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/lex.go,omit,failed,2,TestParseMarksLiveSubstitution/R01_double-quoted_dollar_paren,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/shellcommand,^TestParseMarksLiveSubstitution$,passed,12\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,687\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,2,12\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R01_double-quoted_dollar_paren,\"substitution_test.go:42: marked tokens of Parse(\\\"rg \\\\\\\\\\\"$(x)\\\\\\\\\\\" .\\\") = [], want [{Index:1 Token:{Kind:0 Text:$(x) Substitution:true}}]\",1\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R02_double-quoted_backtick,\"substitution_test.go:42: marked tokens of Parse(\\\"rg \\\\\\\\\\\"`x`\\\\\\\\\\\" .\\\") = [], want [{Index:1 Token:{Kind:0 Text:`x` Substitution:true}}]\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "eceaeaceb006462864ea135e751b1698768f6cfc",
          "requirement": "mark-tests",
          "command": "bench test --package ./internal/shellcommand",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the double-quote branch of the mark and record the R01 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft366-author-01",
              "digest": "sha256:223c7d2e9b798089b676d00358709a2d5ae99d667b6df10b5a99ff40138cf296",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/lex.go,omit,failed,2,TestParseMarksLiveSubstitution/R01_double-quoted_dollar_paren,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/shellcommand,^TestParseMarksLiveSubstitution$,passed,12\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,687\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,2,12\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R01_double-quoted_dollar_paren,\"substitution_test.go:42: marked tokens of Parse(\\\"rg \\\\\\\\\\\"$(x)\\\\\\\\\\\" .\\\") = [], want [{Index:1 Token:{Kind:0 Text:$(x) Substitution:true}}]\",1\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R02_double-quoted_backtick,\"substitution_test.go:42: marked tokens of Parse(\\\"rg \\\\\\\\\\\"`x`\\\\\\\\\\\" .\\\") = [], want [{Index:1 Token:{Kind:0 Text:`x` Substitution:true}}]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c1-t01-guard-preservation",
          "performer": "ft366-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "0f739eaa220393c4db16c40462326410c8702b5d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft366-author-01",
            "digest": "sha256:4cba389c2f46ffe353dae9aa9bb3b1fc8a53c14f569dd93b355c86f2af8b1326",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/lex.go,swap,failed,2,TestParseMarksLiveSubstitution/R06_spaced_dollar,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gitguard ./internal/benchguard ./internal/census ./internal/shellcommand,all,passed,421\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,804\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,pass,3,34\n  github.com/gibbonmi/bench/internal/census,pass,119,75\n  github.com/gibbonmi/bench/internal/gitguard,pass,15196,250\n  github.com/gibbonmi/bench/internal/shellcommand,fail,3,62\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R06_spaced_dollar,\"substitution_test.go:42: marked tokens of Parse(\\\"a $ (b)\\\") = [{Index:2 Token:{Kind:1 Text:( Substitution:true}}], want []\",1\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R06_subshell,\"substitution_test.go:42: marked tokens of Parse(\\\"(rg x .)\\\") = [{Index:0 Token:{Kind:1 Text:( Substitution:true}}], want []\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "eceaeaceb006462864ea135e751b1698768f6cfc",
          "requirement": "guard-preservation",
          "command": "bench test --package ./internal/gitguard --package ./internal/benchguard --package ./internal/census",
          "exit_code": 0,
          "probe": {
            "mutation": "Mark the bare subshell operator and record the R06 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft366-author-01",
              "digest": "sha256:4cba389c2f46ffe353dae9aa9bb3b1fc8a53c14f569dd93b355c86f2af8b1326",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/lex.go,swap,failed,2,TestParseMarksLiveSubstitution/R06_spaced_dollar,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gitguard ./internal/benchguard ./internal/census ./internal/shellcommand,all,passed,421\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,804\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,pass,3,34\n  github.com/gibbonmi/bench/internal/census,pass,119,75\n  github.com/gibbonmi/bench/internal/gitguard,pass,15196,250\n  github.com/gibbonmi/bench/internal/shellcommand,fail,3,62\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R06_spaced_dollar,\"substitution_test.go:42: marked tokens of Parse(\\\"a $ (b)\\\") = [{Index:2 Token:{Kind:1 Text:( Substitution:true}}], want []\",1\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R06_subshell,\"substitution_test.go:42: marked tokens of Parse(\\\"(rg x .)\\\") = [{Index:0 Token:{Kind:1 Text:( Substitution:true}}], want []\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c1-t01-guard-golden",
          "performer": "ft366-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "0f739eaa220393c4db16c40462326410c8702b5d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft366-author-01",
            "digest": "sha256:ac93d42045f2b0f1c2506b94505ca51bc8b2b2b9c979b70ddc9e97045611d034",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/shellcommand.go,swap,failed,1,TestGuardProjectionDifferentialFamilies,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,TestGuardProjectionDifferentialFamilies,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,598\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,13672,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestGuardProjectionDifferentialFamilies,\"guard_projection_differential_test.go:161: GP46 members differ from the golden: P \\\"env -u X bench help | cat\\\"; P \\\"timeout -s KILL -k 1 5 bench help | cat\\\"; P \\\"xargs -n 1 bench help | cat\\\"; W \\\"bash -lc 'bench help | cat'\\\"; W4 \\\"bash -c 'echo … (296 bytes)\",13\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "eceaeaceb006462864ea135e751b1698768f6cfc",
          "requirement": "guard-golden",
          "command": "bench test --check system --run TestGuardProjectionDifferentialFamilies",
          "exit_code": 0,
          "probe": {
            "mutation": "Lex a | operator run as a Word token and record the R55 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft366-author-01",
              "digest": "sha256:ac93d42045f2b0f1c2506b94505ca51bc8b2b2b9c979b70ddc9e97045611d034",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/shellcommand.go,swap,failed,1,TestGuardProjectionDifferentialFamilies,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,TestGuardProjectionDifferentialFamilies,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,598\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,13672,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestGuardProjectionDifferentialFamilies,\"guard_projection_differential_test.go:161: GP46 members differ from the golden: P \\\"env -u X bench help | cat\\\"; P \\\"timeout -s KILL -k 1 5 bench help | cat\\\"; P \\\"xargs -n 1 bench help | cat\\\"; W \\\"bash -lc 'bench help | cat'\\\"; W4 \\\"bash -c 'echo … (296 bytes)\",13\nskips[0]{package,test,reason}:\n"
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

## C1 ticket 01 author witness reds

The ticket 01 author `ft366-author-01` (opus / xhigh) committed ticket 01 at `eceaeaceb` and froze C1 at `ecaf92782`.
The R01, R02, and R03 cases went red before the first production edit, when `Token.Substitution` existed with no logic. Each other witness went red under a `bench probe` mutation after the ticket commit.
The plan probes `mark-tests`, `guard-preservation`, and `guard-golden` answered `bit` with `restored` `yes`, and the record holds them.

| subject | mutation | failing witness | row |
| --- | --- | --- | --- |
| `internal/shellcommand/lex.go` | field with no logic, before the production edit | `TestParseMarksLiveSubstitution` R01, R02, and R03 cases | R01, R02, R03 |
| `internal/shellcommand/lex.go` | the flush marks from token text alone | `TestParseMarksLiveSubstitution` R04 single-quoted backtick case | R04 |
| `internal/shellcommand/lex.go` | the double-quote escape branch for `\$` and a backslash before a backtick removed | `TestParseMarksLiveSubstitution` R05 escaped dollar case | R05 |
| `internal/shellcommand/lex.go` | the flush keeps the `active` flag | `TestParsePublishesTokenClassesAndSimpleCommands` | R07 |

Two plan probe facts differ from the plan text.
The `guard-preservation` plan command names `./internal/gitguard`, `./internal/benchguard`, and `./internal/census`, but the R06 witness is in `./internal/shellcommand`. The author added that package to the same probe call, and the three guard packages stayed green.
The first `guard-golden` probe was silent, because its mutation in `lex` was equivalent: `Parse` classifies a lone `|` as an operator through `allPunct`. The second probe made `Parse` keep a `|` run as a Word token, and the golden failed on the `| cat` follow-on tuples. No row case is missing.

The coordinator probe removed the unquoted backtick mark in `lex` at a site that no author probe reached. `TestParseMarksLiveSubstitution` failed on its R02 bare backtick case, and the probe restored the file.
