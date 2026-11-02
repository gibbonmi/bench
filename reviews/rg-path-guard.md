# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/rg-path-guard/spec.md",
  "plan_digest": "sha256:6c2edf10a2f6a970760dd07cf4a7b57136b06fac0fc2d8cf365022c799e8afb4",
  "implementation_session": "",
  "chunks": [
    {
      "id": "C1",
      "base": "7cc35cbbf401ead2a5b374454e7441a81ec884d0",
      "tip": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
      "plan_digest": "sha256:61fa77cebbc88c0140fa3687c5d113a593a9c56d22711ec3960021034d1a5635",
      "source_digest": "b69a7389267bef67316b749a3f76c7811a311420",
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
        },
        {
          "id": "c1-t01-rep1-mark-tests",
          "performer": "ft366-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "b69a7389267bef67316b749a3f76c7811a311420",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft366-author-01-rep1",
            "digest": "sha256:930b81168d890dc408dfeb29d8f70bb416ddedaa4286d23f6e0bf09f4bdb4926",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/lex.go,omit,failed,2,TestParseMarksLiveSubstitution/R01_double-quoted_dollar_paren,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/shellcommand,all,passed,63\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,687\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,3,63\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R01_double-quoted_dollar_paren,\"substitution_test.go:43: marked tokens of Parse(\\\"rg \\\\\\\\\\\"$(x)\\\\\\\\\\\" .\\\") = [], want [{Index:1 Token:{Kind:0 Text:$(x) Substitution:true}}]\",1\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R02_double-quoted_backtick,\"substitution_test.go:43: marked tokens of Parse(\\\"rg \\\\\\\\\\\"`x`\\\\\\\\\\\" .\\\") = [], want [{Index:1 Token:{Kind:0 Text:`x` Substitution:true}}]\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "requirement": "mark-tests",
          "command": "bench test --package ./internal/shellcommand",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the double-quote branch of the mark and record the R01 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft366-author-01-rep1",
              "digest": "sha256:930b81168d890dc408dfeb29d8f70bb416ddedaa4286d23f6e0bf09f4bdb4926",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/lex.go,omit,failed,2,TestParseMarksLiveSubstitution/R01_double-quoted_dollar_paren,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/shellcommand,all,passed,63\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,687\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,3,63\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R01_double-quoted_dollar_paren,\"substitution_test.go:43: marked tokens of Parse(\\\"rg \\\\\\\\\\\"$(x)\\\\\\\\\\\" .\\\") = [], want [{Index:1 Token:{Kind:0 Text:$(x) Substitution:true}}]\",1\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R02_double-quoted_backtick,\"substitution_test.go:43: marked tokens of Parse(\\\"rg \\\\\\\\\\\"`x`\\\\\\\\\\\" .\\\") = [], want [{Index:1 Token:{Kind:0 Text:`x` Substitution:true}}]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c1-t01-rep1-guard-preservation",
          "performer": "ft366-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "0f739eaa220393c4db16c40462326410c8702b5d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "carried:c1-t01-guard-preservation",
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
              "ref": "carried:c1-t01-guard-preservation",
              "digest": "sha256:4cba389c2f46ffe353dae9aa9bb3b1fc8a53c14f569dd93b355c86f2af8b1326",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/lex.go,swap,failed,2,TestParseMarksLiveSubstitution/R06_spaced_dollar,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gitguard ./internal/benchguard ./internal/census ./internal/shellcommand,all,passed,421\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,804\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/benchguard,pass,3,34\n  github.com/gibbonmi/bench/internal/census,pass,119,75\n  github.com/gibbonmi/bench/internal/gitguard,pass,15196,250\n  github.com/gibbonmi/bench/internal/shellcommand,fail,3,62\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R06_spaced_dollar,\"substitution_test.go:42: marked tokens of Parse(\\\"a $ (b)\\\") = [{Index:2 Token:{Kind:1 Text:( Substitution:true}}], want []\",1\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R06_subshell,\"substitution_test.go:42: marked tokens of Parse(\\\"(rg x .)\\\") = [{Index:0 Token:{Kind:1 Text:( Substitution:true}}], want []\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c1-t01-rep1-guard-golden",
          "performer": "ft366-author-01-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "0f739eaa220393c4db16c40462326410c8702b5d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "carried:c1-t01-guard-golden",
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
              "ref": "carried:c1-t01-guard-golden",
              "digest": "sha256:ac93d42045f2b0f1c2506b94505ca51bc8b2b2b9c979b70ddc9e97045611d034",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/shellcommand.go,swap,failed,1,TestGuardProjectionDifferentialFamilies,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,TestGuardProjectionDifferentialFamilies,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,598\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,13672,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestGuardProjectionDifferentialFamilies,\"guard_projection_differential_test.go:161: GP46 members differ from the golden: P \\\"env -u X bench help | cat\\\"; P \\\"timeout -s KILL -k 1 5 bench help | cat\\\"; P \\\"xargs -n 1 bench help | cat\\\"; W \\\"bash -lc 'bench help | cat'\\\"; W4 \\\"bash -c 'echo … (296 bytes)\",13\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c1-r1-standards",
          "performer": "ft366-review-c1-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0f739eaa220393c4db16c40462326410c8702b5d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c1-r1",
            "digest": "sha256:a3f18338c99513d5c2ebc68b19afee75da82e2ea514bac6023250950bbda0b63",
            "excerpt": "## Standards\nFindings: 0. Worst issue: none.\ndollarParen is the one $( test; shellcommand.go drops to about 337 lines; the canary anchor at shellcommand.go:18 does not move; comments state current behaviour; the R01-R06 cell rewrite is required by bench-implement-spec.\n"
          },
          "axis": "Standards",
          "base": "7cc35cbbf401ead2a5b374454e7441a81ec884d0",
          "tip": "eceaeaceb006462864ea135e751b1698768f6cfc",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1-r1-spec",
          "performer": "ft366-review-c1-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0f739eaa220393c4db16c40462326410c8702b5d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c1-r1",
            "digest": "sha256:66484bf9e4dd0a6c2e680e37ec3cb8f03caebb3964f69229cf828b3199e6ddb2",
            "excerpt": "## Spec\nFindings: 0. Worst issue: none.\nR01-R06 each have a matching case in substitution_test.go; R07 literals keyed with no text, kind, or span change; R55 golden probe bit; hand traces of lex agree with POSIX and spec:129; plan note (guard-preservation package) is no-op, disclosed in the review record.\n"
          },
          "axis": "Spec",
          "base": "7cc35cbbf401ead2a5b374454e7441a81ec884d0",
          "tip": "eceaeaceb006462864ea135e751b1698768f6cfc",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1-r1-coverage",
          "performer": "ft366-review-c1-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0f739eaa220393c4db16c40462326410c8702b5d",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c1-r1",
            "digest": "sha256:bd11dede69134a6f99226bb0d6de40155a49c99358b9539b208eaa491b6251d1",
            "excerpt": "## Coverage\nFindings: 1. Worst issue: C1.\nC1: internal/shellcommand/lex.go:34, confidence 7, auto-fix. No test grades the opener reset after an operator; bench probe --omit 'opener = false' was silent. Fix: add an R03 case such as rg $('x') . that expects only the ( operator marked.\nProbe 2: swap of the unlexed fallback marks bit (TestFallbackSplitsShellOperators).\n"
          },
          "axis": "Coverage",
          "base": "7cc35cbbf401ead2a5b374454e7441a81ec884d0",
          "tip": "eceaeaceb006462864ea135e751b1698768f6cfc",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        }
      ],
      "confirmations": [
        {
          "id": "c1-r1-confirm",
          "performer": "session_01JK5CPdJaS9Xn1VAmFL1sNN",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b69a7389267bef67316b749a3f76c7811a311420",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:session_01JK5CPdJaS9Xn1VAmFL1sNN:c1-confirm",
            "digest": "sha256:8ff3049f579ef30ec19b9a09bd76ce32eaacee5482f839f4ce7e66ba19c6858f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/lex.go,swap,failed,1,TestParseMarksLiveSubstitution/R03_quoted_word_in_bare_dollar_paren,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/shellcommand,^TestParseMarksLiveSubstitution$,passed,13\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,528\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,2,13\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R03_quoted_word_in_bare_dollar_paren,\"substitution_test.go:43: marked tokens of Parse(\\\"rg $('x') .\\\") = [{Index:2 Token:{Kind:1 Text:( Substitution:true}} {Index:4 Token:{Kind:1 Text:) Substitution:true}}], want [{Index:2 Token:{Kind:1 Text:( Substitution:true}}]\",1\nskips[0]{package,test,reason}:\n"
          },
          "base": "7cc35cbbf401ead2a5b374454e7441a81ec884d0",
          "reviewed_tip": "eceaeaceb006462864ea135e751b1698768f6cfc",
          "repaired_tip": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "finding_ids": [
            "C1"
          ],
          "supersedes": [],
          "probe": {
            "mutation": "In lex, the single-quote loop sets opener true, so a punctuation run after a quoted word keeps a stale opener; TestParseMarksLiveSubstitution R03 quoted word in bare dollar paren fails",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:session_01JK5CPdJaS9Xn1VAmFL1sNN:c1-confirm",
              "digest": "sha256:8ff3049f579ef30ec19b9a09bd76ce32eaacee5482f839f4ce7e66ba19c6858f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/shellcommand/lex.go,swap,failed,1,TestParseMarksLiveSubstitution/R03_quoted_word_in_bare_dollar_paren,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/shellcommand,^TestParseMarksLiveSubstitution$,passed,13\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,528\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/shellcommand,fail,2,13\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/shellcommand,TestParseMarksLiveSubstitution/R03_quoted_word_in_bare_dollar_paren,\"substitution_test.go:43: marked tokens of Parse(\\\"rg $('x') .\\\") = [{Index:2 Token:{Kind:1 Text:( Substitution:true}} {Index:4 Token:{Kind:1 Text:) Substitution:true}}], want [{Index:2 Token:{Kind:1 Text:( Substitution:true}}]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ]
    },
    {
      "id": "C2",
      "base": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
      "tip": "eda20a987adc8e38dced29a3fda32b8b07898b04",
      "plan_digest": "sha256:6c2edf10a2f6a970760dd07cf4a7b57136b06fac0fc2d8cf365022c799e8afb4",
      "source_digest": "aacdc93af5b100c577ff86add639ceb90937cdca",
      "acceptance_rows": [
        "R08",
        "R09",
        "R10",
        "R11",
        "R12",
        "R13",
        "R14",
        "R15",
        "R16",
        "R17",
        "R18",
        "R19",
        "R20",
        "R21",
        "R22",
        "R23",
        "R24",
        "R25",
        "R26",
        "R27",
        "R28",
        "R29",
        "R30",
        "R31",
        "R32",
        "R33",
        "R34",
        "R35",
        "R36",
        "R37",
        "R38",
        "R58"
      ],
      "verification": [
        {
          "id": "c2-t02-verdict-tests",
          "performer": "ft366-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aacdc93af5b100c577ff86add639ceb90937cdca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft366-author-02",
            "digest": "sha256:d379b9892b8d0c0c7bea06d815f29759cb179e23c359fdf36e0545a805543fb2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/options.go,omit,failed,1,TestClassifyRgCalls,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/rgguard,^TestClassifyRgCalls$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,377\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/rgguard,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/rgguard,TestClassifyRgCalls,\"rgguard_test.go:81: R13: Classify(\\\"rg -g '*.go' foo\\\") = {Blocked:false Segment:[] reason:0}, want {Blocked:true Segment:[] reason:0}\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "eda20a987adc8e38dced29a3fda32b8b07898b04",
          "requirement": "verdict-tests",
          "command": "bench test --package ./internal/rgguard",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove one value letter from the option table and record the R13 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft366-author-02",
              "digest": "sha256:d379b9892b8d0c0c7bea06d815f29759cb179e23c359fdf36e0545a805543fb2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/options.go,omit,failed,1,TestClassifyRgCalls,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/rgguard,^TestClassifyRgCalls$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,377\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/rgguard,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/rgguard,TestClassifyRgCalls,\"rgguard_test.go:81: R13: Classify(\\\"rg -g '*.go' foo\\\") = {Blocked:false Segment:[] reason:0}, want {Blocked:true Segment:[] reason:0}\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c2-t02-projection-owner",
          "performer": "ft366-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aacdc93af5b100c577ff86add639ceb90937cdca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft366-author-02",
            "digest": "sha256:f9b386aca52e1a5bf6afc69a4e0e63db86a361846aafda9a56288674d6d8cf5d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/options.go,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,package-core-guard,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,441\ncheck[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,2723,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: guard projection owner internal/rgguard/options.go:6: hook envelope JSON decoding\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "eda20a987adc8e38dced29a3fda32b8b07898b04",
          "requirement": "projection-owner",
          "command": "bench test --check package-core-guard",
          "exit_code": 0,
          "probe": {
            "mutation": "Plant a private envelope decode in internal/rgguard and record the R37 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft366-author-02",
              "digest": "sha256:f9b386aca52e1a5bf6afc69a4e0e63db86a361846aafda9a56288674d6d8cf5d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/options.go,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,package-core-guard,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,441\ncheck[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,2723,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: guard projection owner internal/rgguard/options.go:6: hook envelope JSON decoding\",1\nskips[0]{package,test,reason}:\n"
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
      "from": "sha256:661349ba0ae7d631cf0c2b43e02bb65b60aefca056311d23ffec0b8190db0fc9",
      "to": "sha256:61fa77cebbc88c0140fa3687c5d113a593a9c56d22711ec3960021034d1a5635",
      "chunk_ids": {
        "C1": [
          "C1"
        ]
      }
    },
    {
      "from": "sha256:61fa77cebbc88c0140fa3687c5d113a593a9c56d22711ec3960021034d1a5635",
      "to": "sha256:6c2edf10a2f6a970760dd07cf4a7b57136b06fac0fc2d8cf365022c799e8afb4",
      "chunk_ids": {
        "C1": [
          "C1"
        ]
      }
    }
  ]
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

## C1 review round 1

One unified session (`execution.review_mode` is `unified`) ran the three axes separately on opus at medium effort over `7cc35cbb..eceaeace`.
The orchestrator accepts C1 under `--auto-approve`. The round has one repair target on ticket 01, and the repair adds only a test case.
A coordinator probe that bites then closes it through `bench record confirmation`.

## Standards

Findings: 0. Worst issue: none.

## Spec

Findings: 0. Worst issue: none. The reviewer marks the `guard-preservation` plan text as a no-op, because the record above discloses the added package.

## Coverage

Findings: 1 accepted. Worst issue: C1.

- C1 (`c1-r1-coverage`), auto-fix, confidence 7: no test grades the reset of the opener flag at `internal/shellcommand/lex.go:34`. A probe that omits `opener = false` stays silent, and then `rg $('x') .` also marks the closing `)`. Add an R03 case such as `rg $('x') .` to `TestParseMarksLiveSubstitution` that expects only the `(` operator marked.

## C1 repair cycle 1

The repair session `ft366-author-01-rep1` (opus / xhigh) added the R03 case `rg $('x') .` at `c8519a1b0` and changed no production code. It refroze C1 at `fb1c15f6c`.
Its witness probe omitted `opener = false` and failed the new case, and its rerun of the `mark-tests` plan probe answered `bit`.
The coordinator probe made the single-quote loop of `lex` set the opener, and the new case failed. The confirmation `c1-r1-confirm` closes C1.
C1 used one of its two repair cycles.

## C2 ticket 02 author witness reds

The ticket 02 author `ft366-author-02` (opus / xhigh) committed ticket 02 at `eda20a987` and froze C2 at `be4ade1fc`.
Before the first production edit, a `Classify` stub that permits every call made the refusal rows red.
These rows are R08, R11, R13–R15, R17, R24, R25, R27–R31, R33, and R58.
A stub `IsBenchName` that answers false made `TestIsBenchName` red, and R37 was red before `internal/rgguard` joined `guardProjectionGoverned`.
Each permit row went red under a `bench probe` mutation before the commit, and each probe restored the file.

| row | mutation in `internal/rgguard` |
| --- | --- |
| R09 | the one-positional test answers true |
| R10 | the `\|` lead operator removed |
| R12 | the input-redirection test answers false |
| R16 | the mode term of the permit removed |
| R18 | an unknown option word reads as a flag |
| R19 | an unlexed stream refuses |
| R20 | the `ViaXargs` test removed |
| R21 | the process-substitution test removed |
| R22 | the inline config assignment test removed |
| R23 | the one-positional test reads at most one |
| R26 | the mark test reads raw text for `$` and a backtick |
| R32 | the fed flag of the exec child cleared |
| R34 | the `Executes` test removed |
| R35 | the stdin request cleared |
| R36 | the one-level wrapper bound removed |
| R58 | the marks of the outer exec span cleared |

The plan probes `verdict-tests` and `projection-owner` answered `bit` with `restored` `yes`, and the record holds them.
The first R16 probe was silent, because each R16 call has no pattern and the pathless step already permits it. The author added `rg --files src`, and the R16 probe then failed. For the same reason, R20 gained `xargs rg foo` and R34 gained `command -v rg foo`.

The coordinator probe made the `-e` and `-f` branch of `pathless` count one positional as the pattern. `TestClassifyRgCalls` failed on R15, and the probe restored the file.
