# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/rg-path-guard/spec.md",
  "plan_digest": "sha256:d96ea735f3a085d172be23c274b2a29171ca1dc2bd702c49dd18f7af18b9178e",
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
      "tip": "92dda48de69eed72ed8c3de6a1c7b3894ca81f52",
      "plan_digest": "sha256:079fe43eda2b77a7eb921c044d8bc4efb857f81a6d1f197acaae4989661a3ce7",
      "source_digest": "38f2fe6551de667e4f1fd73cd9b8fadaa33b2044",
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
        "R58",
        "R59"
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
        },
        {
          "id": "c2-t02-rep1-verdict-tests",
          "performer": "ft366-author-02-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "2b9c376541b7dda0724520bd0ae65fe00dc4cf75",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft366-author-02-rep1",
            "digest": "sha256:d379b9892b8d0c0c7bea06d815f29759cb179e23c359fdf36e0545a805543fb2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/options.go,omit,failed,1,TestClassifyRgCalls,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/rgguard,^TestClassifyRgCalls$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,377\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/rgguard,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/rgguard,TestClassifyRgCalls,\"rgguard_test.go:81: R13: Classify(\\\"rg -g '*.go' foo\\\") = {Blocked:false Segment:[] reason:0}, want {Blocked:true Segment:[] reason:0}\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e11d6e68f84315580fc6e001e296a0ea1dd2be12",
          "requirement": "verdict-tests",
          "command": "bench test --package ./internal/rgguard",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove one value letter from the option table and record the R13 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft366-author-02-rep1",
              "digest": "sha256:d379b9892b8d0c0c7bea06d815f29759cb179e23c359fdf36e0545a805543fb2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/options.go,omit,failed,1,TestClassifyRgCalls,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/rgguard,^TestClassifyRgCalls$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,377\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/rgguard,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/rgguard,TestClassifyRgCalls,\"rgguard_test.go:81: R13: Classify(\\\"rg -g '*.go' foo\\\") = {Blocked:false Segment:[] reason:0}, want {Blocked:true Segment:[] reason:0}\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c2-t02-rep1-projection-owner",
          "performer": "ft366-author-02-rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "2b9c376541b7dda0724520bd0ae65fe00dc4cf75",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft366-author-02-rep1",
            "digest": "sha256:21247a13ddfa32339ebc29c9abf450e4d18069d7321a5b96117caefa051504db",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/rgguard.go,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,package-core-guard,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,442\ncheck[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,2694,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: guard projection owner internal/rgguard/rgguard.go:18: hook envelope JSON decoding\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e11d6e68f84315580fc6e001e296a0ea1dd2be12",
          "requirement": "projection-owner",
          "command": "bench test --check package-core-guard",
          "exit_code": 0,
          "probe": {
            "mutation": "Plant a private envelope decode in internal/rgguard and record the R37 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft366-author-02-rep1",
              "digest": "sha256:21247a13ddfa32339ebc29c9abf450e4d18069d7321a5b96117caefa051504db",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/rgguard.go,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,package-core-guard,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,442\ncheck[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,2694,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: guard projection owner internal/rgguard/rgguard.go:18: hook envelope JSON decoding\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c2-t02-rep2-verdict-tests",
          "performer": "ft366-author-02-rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "38f2fe6551de667e4f1fd73cd9b8fadaa33b2044",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft366-author-02-rep2",
            "digest": "sha256:b3e9fe13dad17ef55a2b3c744b1ae255f12a6732a900f1b4a4b7001a4bd1c01c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/options.go,omit,failed,1,TestClassifyRgCalls,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/rgguard,^TestClassifyRgCalls$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,377\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/rgguard,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/rgguard,TestClassifyRgCalls,\"rgguard_test.go:82: R13: Classify(\\\"rg -g '*.go' foo\\\") = {Blocked:false Segment:[] reason:0}, want {Blocked:true Segment:[] reason:0}\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "92dda48de69eed72ed8c3de6a1c7b3894ca81f52",
          "requirement": "verdict-tests",
          "command": "bench test --package ./internal/rgguard",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove one value letter from the option table and record the R13 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft366-author-02-rep2",
              "digest": "sha256:b3e9fe13dad17ef55a2b3c744b1ae255f12a6732a900f1b4a4b7001a4bd1c01c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/options.go,omit,failed,1,TestClassifyRgCalls,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/rgguard,^TestClassifyRgCalls$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,377\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/rgguard,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/rgguard,TestClassifyRgCalls,\"rgguard_test.go:82: R13: Classify(\\\"rg -g '*.go' foo\\\") = {Blocked:false Segment:[] reason:0}, want {Blocked:true Segment:[] reason:0}\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c2-t02-rep2-projection-owner",
          "performer": "ft366-author-02-rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "38f2fe6551de667e4f1fd73cd9b8fadaa33b2044",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:ft366-author-02-rep2",
            "digest": "sha256:4dc79f1d53f618641e836d2d8ba2b53335dbff584967f9c7069ba996a12d9f58",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/rgguard.go,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,package-core-guard,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,442\ncheck[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,2546,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: guard projection owner internal/rgguard/rgguard.go:17: hook envelope JSON decoding\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "92dda48de69eed72ed8c3de6a1c7b3894ca81f52",
          "requirement": "projection-owner",
          "command": "bench test --check package-core-guard",
          "exit_code": 0,
          "probe": {
            "mutation": "Plant a private envelope decode in internal/rgguard and record the R37 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:ft366-author-02-rep2",
              "digest": "sha256:4dc79f1d53f618641e836d2d8ba2b53335dbff584967f9c7069ba996a12d9f58",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/rgguard/rgguard.go,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,package-core-guard,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,442\ncheck[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,2546,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: guard projection owner internal/rgguard/rgguard.go:17: hook envelope JSON decoding\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c2-r1-standards",
          "performer": "ft366-review-c2-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "aacdc93af5b100c577ff86add639ceb90937cdca",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c2-r1",
            "digest": "sha256:47a9be35b0fda166e53caf151aab459b77b99d45b2ac0573f2b5fc94d4ea5971",
            "excerpt": "## Standards\nFindings: 1. Worst issue: S1.\nS1: internal/rgguard/rgguard_test.go:19-21, confidence 5, auto-fix. The comment tells the history of the added inputs and why they turn red, which the review record owns. Cut it to the expectation-source sentence.\nIsBenchName is the one Bench-name test; R38 sweep finds no second option list; no new dependency.\n"
          },
          "axis": "Standards",
          "base": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "tip": "eda20a987adc8e38dced29a3fda32b8b07898b04",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "c2-r1-spec",
          "performer": "ft366-review-c2-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "aacdc93af5b100c577ff86add639ceb90937cdca",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c2-r1",
            "digest": "sha256:14b5ded9adb91fa127a98f3b5033c1e2d2eb9c029af700ee1779cd859f61fa85",
            "excerpt": "## Spec\nFindings: 1. Worst issue: P1.\nP1: R16 row, options.go:163 and :230-231, confidence 9, auto-fix. rg --version and rg -nh have no pattern, so the pathless step permits them without the mode list; a probe of -h mode to flag was silent. Add rg -nh foo and rg --version foo to R16 and fix the row text.\nAuthor divergences hold; option table matches rg 15.1.0 help; Classify starts no process; verdict order matches steps 2-6.\n"
          },
          "axis": "Spec",
          "base": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "tip": "eda20a987adc8e38dced29a3fda32b8b07898b04",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "c2-r1-coverage",
          "performer": "ft366-review-c2-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "aacdc93af5b100c577ff86add639ceb90937cdca",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c2-r1",
            "digest": "sha256:b06f0303416b329866d0ca9127926eae7231ec98f798f9cd27f1845bac5cee6b",
            "excerpt": "## Coverage\nFindings: 3. Worst issue: C1.\nC1: same gap as P1, confidence 9, auto-fix.\nC2: internal/rgguard/rgguard.go:78, confidence 8, auto-fix. Omitting fed = false is silent; add cat f | bash -c 'rg x .; rg foo' (refuse).\nC3: internal/rgguard/rgguard.go:145, confidence 8, ask-user. cat f | (rg foo) is refused although the pipe is the stdin of rg; decide over-refusal or treat a lead ( after a pipe as fed.\n"
          },
          "axis": "Coverage",
          "base": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "tip": "eda20a987adc8e38dced29a3fda32b8b07898b04",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "c2-r2-standards",
          "performer": "ft366-review-c2-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2b9c376541b7dda0724520bd0ae65fe00dc4cf75",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c2-r2",
            "digest": "sha256:b9b8c097c6d98f856113c07eb9bd3287314f567cf4d4b9081c7cda7b5c5c4bfa",
            "excerpt": "## Standards\nS1 closed. One-source check passes: supplied is the one lead-operator read in rgguard.\nNew findings: 1. S-n1: internal/rgguard/rgguard.go:146, confidence 5, auto-fix. The !Substitution term on the group-lead ( is dead code (lex.go:79 and :33 keep the $ in the word before a marked (); deleting it is silent. Delete the term and its comment sentence.\n"
          },
          "axis": "Standards",
          "base": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "tip": "e11d6e68f84315580fc6e001e296a0ea1dd2be12",
          "finding_ids": [
            "S2"
          ],
          "supersedes": [
            "c2-r1-standards"
          ]
        },
        {
          "id": "c2-r2-spec",
          "performer": "ft366-review-c2-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2b9c376541b7dda0724520bd0ae65fe00dc4cf75",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c2-r2",
            "digest": "sha256:59ba04286f56312be1f10ee3c85640585653abb4fbb1291d80cabac0a7751c84",
            "excerpt": "## Spec\nP1 closed. C3/R59 partly closed.\nNew findings: 1. P-n1: spec.md:92 story 59, rgguard.go:146-152, confidence 8, auto-fix. cat f |(rg foo) is refused: the lexer joins |( into one operator token (lex.go:26-33). Treat a pipe operator followed by ( as a supplier and add the input to R59.\n"
          },
          "axis": "Spec",
          "base": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "tip": "e11d6e68f84315580fc6e001e296a0ea1dd2be12",
          "finding_ids": [
            "P2"
          ],
          "supersedes": [
            "c2-r1-spec"
          ]
        },
        {
          "id": "c2-r2-coverage",
          "performer": "ft366-review-c2-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2b9c376541b7dda0724520bd0ae65fe00dc4cf75",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c2-r2",
            "digest": "sha256:ad8e0e7fecd01977f03548db6a62d00893c75e31c39a012bb458c23eba7976dd",
            "excerpt": "## Coverage\nC1 closed. C2 closed. C3 closed for the spaced form only; the no-space form is P-n1.\nNew findings: 0.\n"
          },
          "axis": "Coverage",
          "base": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "tip": "e11d6e68f84315580fc6e001e296a0ea1dd2be12",
          "finding_ids": [],
          "supersedes": [
            "c2-r1-coverage"
          ]
        },
        {
          "id": "c2-r3-standards",
          "performer": "ft366-review-c2-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38f2fe6551de667e4f1fd73cd9b8fadaa33b2044",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c2-r3",
            "digest": "sha256:63ab90b627f1d2a46ad3c75d494fe25195de079fb7927d7738418952891f63a5",
            "excerpt": "## Standards\nS2 closed: the Substitution term and its comment sentence are gone; one list of pipe texts at rgguard.go:154.\nNew findings: 0.\n"
          },
          "axis": "Standards",
          "base": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "tip": "92dda48de69eed72ed8c3de6a1c7b3894ca81f52",
          "finding_ids": [],
          "supersedes": [
            "c2-r2-standards"
          ]
        },
        {
          "id": "c2-r3-spec",
          "performer": "ft366-review-c2-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38f2fe6551de667e4f1fd73cd9b8fadaa33b2044",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c2-r3",
            "digest": "sha256:24942af5182f39590bb9ba15862de4ca66d37e7b8a2226e004199885addf579c",
            "excerpt": "## Spec\nP2 closed: cat f |(rg foo) permits; the comment at rgguard.go:142-144 and spec line 149 match the code. Concern (a) probed over 14 inputs: no earlier verdict changed except |&( which P2 covers.\nNew findings: 0.\n"
          },
          "axis": "Spec",
          "base": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "tip": "92dda48de69eed72ed8c3de6a1c7b3894ca81f52",
          "finding_ids": [],
          "supersedes": [
            "c2-r2-spec"
          ]
        },
        {
          "id": "c2-r3-coverage",
          "performer": "ft366-review-c2-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "38f2fe6551de667e4f1fd73cd9b8fadaa33b2044",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft366-review-c2-r3",
            "digest": "sha256:f3f8e0199a369d9604d9c6f1a942b442a3c08dc04c13b3cc43214a2f8c2267e9",
            "excerpt": "## Coverage\nR59 has both inputs: removing the join fails the spaced input, removing the trim fails the no-space input.\nNew findings: 0.\n"
          },
          "axis": "Coverage",
          "base": "c8519a1b0514ac355c11bc7b1a65575b3cf4a9eb",
          "tip": "92dda48de69eed72ed8c3de6a1c7b3894ca81f52",
          "finding_ids": [],
          "supersedes": [
            "c2-r2-coverage"
          ]
        }
      ]
    },
    {
      "id": "C3",
      "base": "92dda48de69eed72ed8c3de6a1c7b3894ca81f52",
      "tip": "eb315403e4cb3bb95148af1d60e4c57b30841665",
      "plan_digest": "sha256:d96ea735f3a085d172be23c274b2a29171ca1dc2bd702c49dd18f7af18b9178e",
      "source_digest": "f379142159ce57929cd7c42159b5d82cf80d5460",
      "acceptance_rows": [
        "R41",
        "R42",
        "R43",
        "R44",
        "R45",
        "R46",
        "R50",
        "R56",
        "R57",
        "R39",
        "R40",
        "R47",
        "R48",
        "R49"
      ],
      "verification": [
        {
          "id": "c3-t03-registrations",
          "performer": "ft366-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c59bed564a4745868b4ae7f35d2275df9b280df",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft366-author-03",
            "digest": "sha256:9b99886972291d95fb8bb018de1da3f0ee31e98dbb0b96600d3eaab92c724203",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/entry_point_parity_rows_test.go,omit,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,entry-point-parity,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,484\ncheck[1]{name,kind,tests_run,subjects}:\n  entry-point-parity,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,771,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: entry-point-parity: .bench/hooks/block-unsafe-rg.sh reaches no parity row; give it a row naming the registry command it calls\",2\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f0c28331b775b7e06b736991368a4f3bde8905d6",
          "requirement": "registrations",
          "command": "bench test --check entry-point-parity",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the guard-rg parity row and record the R56 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft366-author-03",
              "digest": "sha256:9b99886972291d95fb8bb018de1da3f0ee31e98dbb0b96600d3eaab92c724203",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/entry_point_parity_rows_test.go,omit,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,entry-point-parity,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,484\ncheck[1]{name,kind,tests_run,subjects}:\n  entry-point-parity,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,771,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: entry-point-parity: .bench/hooks/block-unsafe-rg.sh reaches no parity row; give it a row naming the registry command it calls\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c3-t03-command-contract",
          "performer": "ft366-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c59bed564a4745868b4ae7f35d2275df9b280df",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft366-author-03",
            "digest": "sha256:1e23a014e4d51d49b89ddce93f866044eb0d3f5f74f268f13c2bb9fb425e0e93",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/otelrecord/registry.go,omit,failed,1,TestHookSeamsMatchTheHookDispatchRows,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,HookSeams|GuardRg|CommandRegistry,passed,7\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,502\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,7,7\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestHookSeamsMatchTheHookDispatchRows,\"otel_hook_seams_test.go:40: hook seams: dispatch [hook.check-agent-line hook.guard-bench-follow-on hook.guard-file-write hook.guard-git hook.guard-rg hook.session-inspect hook.stop-verdict hook.worktree-hook], registry [hook.check-agent-lin… (364 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f0c28331b775b7e06b736991368a4f3bde8905d6",
          "requirement": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'HookSeams|GuardRg|CommandRegistry'",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the hook.guard-rg seam row and record the R46 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft366-author-03",
              "digest": "sha256:1e23a014e4d51d49b89ddce93f866044eb0d3f5f74f268f13c2bb9fb425e0e93",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/otelrecord/registry.go,omit,failed,1,TestHookSeamsMatchTheHookDispatchRows,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,HookSeams|GuardRg|CommandRegistry,passed,7\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,502\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,7,7\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestHookSeamsMatchTheHookDispatchRows,\"otel_hook_seams_test.go:40: hook seams: dispatch [hook.check-agent-line hook.guard-bench-follow-on hook.guard-file-write hook.guard-git hook.guard-rg hook.session-inspect hook.stop-verdict hook.worktree-hook], registry [hook.check-agent-lin… (364 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c3-t03-routing",
          "performer": "ft366-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c59bed564a4745868b4ae7f35d2275df9b280df",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft366-author-03",
            "digest": "sha256:a3caa09cc42439a4f344eac717fe971e53b5af66c2b36b3c9d04ecc09e0169d9",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/subcommand_routing_table_test.go,omit,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,subcommand-routing,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,533\ncheck[1]{name,kind,tests_run,subjects}:\n  subcommand-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,25,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: cmd/bench/main.go dispatches \\\"guard-rg\\\" with no entry in the subcommand argument-routing registry; record it as routed through usage.Parse or as an exemption with its reason\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "f0c28331b775b7e06b736991368a4f3bde8905d6",
          "requirement": "routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the guard-rg subcommandRouting row and record the R57 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft366-author-03",
              "digest": "sha256:a3caa09cc42439a4f344eac717fe971e53b5af66c2b36b3c9d04ecc09e0169d9",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/subcommand_routing_table_test.go,omit,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,subcommand-routing,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,533\ncheck[1]{name,kind,tests_run,subjects}:\n  subcommand-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,25,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: cmd/bench/main.go dispatches \\\"guard-rg\\\" with no entry in the subcommand argument-routing registry; record it as routed through usage.Parse or as an exemption with its reason\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "c3-t04-wiring",
          "performer": "ft366-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f379142159ce57929cd7c42159b5d82cf80d5460",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft366-author-04",
            "digest": "sha256:124a6dc693cb41eb33d0a060a2048cd8cb2845a99e1a11a4ad7d128ceb57a48c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.claude/settings.json,omit,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,line-routing,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,453\ncheck[1]{name,kind,tests_run,subjects}:\n  line-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,860,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: claude settings.json PreToolUse Bash matcher missing or does not run .bench/hooks/block-unsafe-rg.sh\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "eb315403e4cb3bb95148af1d60e4c57b30841665",
          "requirement": "wiring",
          "command": "bench test --check line-routing",
          "exit_code": 0,
          "probe": {
            "mutation": "Remove the shim from the Claude Bash group and record the R39 red, the restore, and the green.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft366-author-04",
              "digest": "sha256:124a6dc693cb41eb33d0a060a2048cd8cb2845a99e1a11a4ad7d128ceb57a48c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.claude/settings.json,omit,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,line-routing,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,453\ncheck[1]{name,kind,tests_run,subjects}:\n  line-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,860,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: claude settings.json PreToolUse Bash matcher missing or does not run .bench/hooks/block-unsafe-rg.sh\",1\nskips[0]{package,test,reason}:\n"
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
    },
    {
      "from": "sha256:6c2edf10a2f6a970760dd07cf4a7b57136b06fac0fc2d8cf365022c799e8afb4",
      "to": "sha256:cc4f36a7ad3e0e75becff9ce97c4d929c24cdf9b9b96941687316c202fb103ca",
      "chunk_ids": {
        "C1": [
          "C1"
        ],
        "C2": [
          "C2"
        ]
      }
    },
    {
      "from": "sha256:cc4f36a7ad3e0e75becff9ce97c4d929c24cdf9b9b96941687316c202fb103ca",
      "to": "sha256:079fe43eda2b77a7eb921c044d8bc4efb857f81a6d1f197acaae4989661a3ce7",
      "chunk_ids": {
        "C1": [
          "C1"
        ],
        "C2": [
          "C2"
        ]
      }
    },
    {
      "from": "sha256:079fe43eda2b77a7eb921c044d8bc4efb857f81a6d1f197acaae4989661a3ce7",
      "to": "sha256:e083920ee6f761c9c2504d1360f5c4b9d72bf8cc552eb2fb3781c35a976a926b",
      "chunk_ids": {
        "C1": [
          "C1"
        ],
        "C2": [
          "C2"
        ]
      }
    },
    {
      "from": "sha256:e083920ee6f761c9c2504d1360f5c4b9d72bf8cc552eb2fb3781c35a976a926b",
      "to": "sha256:d96ea735f3a085d172be23c274b2a29171ca1dc2bd702c49dd18f7af18b9178e",
      "chunk_ids": {
        "C1": [
          "C1"
        ],
        "C2": [
          "C2"
        ],
        "C3": [
          "C3"
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

## C2 review round 1

One unified session ran the three axes separately on opus at medium effort over `c8519a1b..eda20a98`.
The orchestrator accepts S1, P1, C1, C2, and C3 under `--auto-approve`. P1 and C1 name one repair target, so the round has four repair targets, and all of them are on ticket 02.
C3 changes production code, so C2 takes a confirming round after the repair.
The spec amendment corrects the R16 inputs and row text, adds a refusal case to R32, and adds story 59 and row R59 for C3.

The reviewer proposed `ask-user` for C3. The orchestrator decides it. The decision source permits a call when a pipe in the same command supplies stdin, and in `cat f | (rg foo)` the pipe does.
The refusal advice to add `.` would also change what that call searches. The reviewer can veto R59.

## Standards

Findings: 1 accepted. Worst issue: S1.

- S1 (`c2-r1-standards`), auto-fix, confidence 5: the comment at `internal/rgguard/rgguard_test.go:19-21` tells the history of the added inputs and why they turn red. The review record owns that history, so cut the comment to its expectation-source sentence.

## Spec

Findings: 1 accepted. Worst issue: P1.

- P1 (`c2-r1-spec`), auto-fix, confidence 9: R16 has two inputs that no test can fail. `rg --version` and `rg -nh` name no pattern, so the pathless step permits them without the mode list. Add `rg -nh foo` and `rg --version foo` to R16.

## Coverage

Findings: 3 accepted. Worst issue: C1.

- C1 (`c2-r1-coverage`), auto-fix, confidence 9: the same gap as P1. A probe that reads `-h` as a flag at `internal/rgguard/options.go:163` stays silent.
- C2 (`c2-r1-coverage`), auto-fix, confidence 8: no row grades that only the first `rg` span of a fed `-c` child is fed. A probe that omits `fed = false` at `internal/rgguard/rgguard.go:78` stays silent. Add `cat f | bash -c 'rg x .; rg foo'` as a refusal case of R32.
- C3 (`c2-r1-coverage`), auto-fix, confidence 8: `cat f | (rg foo)` is refused, because `internal/rgguard/rgguard.go:145` reads only the token directly before the span. A pipe before a lead `(` feeds the first `rg` span of the group (R59).

## Advice

- The zero `Verdict` holds the pathless reason, so `Message()` on a permit verdict prints the pathless sentence.
- No row grades the attached value branch, such as `-A=3` or `--hidden=x`.

## C2 repair cycle 1

The repair session `ft366-author-02-rep1` (opus / xhigh) closed S1, P1, C1, C2, and C3 at `e11d6e68f` and refroze C2 at `0b9571087`.
Before the production edit, only the R59 permit case was red. The new R16 inputs, the R32 refusal, and the R59 refusal were green on the old code.
Its three witness probes answered `bit` with `restored` `yes`.
`-h` read as a flag failed R16, the omitted `fed = false` failed R32, and the removed group lead failed R59.
It reran the plan probes `verdict-tests` and `projection-owner`, and both answered `bit`.

The repair session reports that the `Substitution` test on the group-lead `(` is unreachable. The lexer marks a `(` only directly after a `$`, so a pipe never comes directly before a marked `(`.
The spec text of verdict step 4 still names only the token directly before the span. The next plan commit corrects that prose.
The coordinator probe added `;` to the feeding operators, and `TestClassifyRgCalls` failed on R32 and R59.
C2 used one of its two repair cycles, and a confirming round of the three axes follows.

## C2 confirming round

One unified session confirmed the repair delta `cbe816fa..e11d6e68` on opus at medium effort. S1, P1, C1, and C2 are closed, and C3 is closed for the spaced form only.
The orchestrator accepts the two new findings under `--auto-approve`, and both are on ticket 02. They take C2 repair cycle 2, the last cycle of the allowance.
The spec amendment adds the no-space input to R59 and states the group lead after verdict step 4.

## Standards

Findings: 1 accepted. Worst issue: S2.

- S2 (`c2-r2-standards`), auto-fix, confidence 5: the `Substitution` term of the group-lead test at `internal/rgguard/rgguard.go:146` is dead code. The lexer keeps the `$` in the word before a marked `(`. Delete the term and its comment sentence.

## Spec

Findings: 1 accepted. Worst issue: P2.

- P2 (`c2-r2-spec`), auto-fix, confidence 8: `cat f |(rg foo)` is refused, because the lexer joins `|(` into one operator token. Treat a pipe operator that ends in `(` as a supplier, and add the input to R59.

## Coverage

Findings: 0. Worst issue: none.

## C2 repair cycle 2

The repair session `ft366-author-02-rep2` (opus / xhigh) closed S2 and P2 at `92dda48de` and refroze C2 at `d82475406`.
The `supplied` read joins a lone `(` to the token before it and drops one trailing `(`.
It keeps the one list of pipe texts.

The new R59 input `cat f |(rg foo)` was red before the edit. A witness probe that restored the old one-step lead failed only that input.
The plan probes `verdict-tests` and `projection-owner` answered `bit`. C2 used both of its repair cycles, and a confirming round follows.

## C2 second confirming round

One unified session confirmed the repair delta `7cbb9f61..92dda48d` on opus at medium effort. S2 and P2 are closed, and no axis has a new finding.
The reviewer probed 14 group and substitution inputs, and no earlier verdict changed except `|&(`, which P2 covers.

## Advice

- No test holds a group input that starts with `|&`, such as `cat f |&(rg foo)`.
- `cat f | $(rg foo)` refuses as pathless, although the substitution can read the pipe. The code before this chunk refused it the same way.

## C3 ticket 03 author witness reds

The ticket 03 author `ft366-author-03` (opus / medium) committed ticket 03 at `f0c28331b` and recorded it at `df454f276`.
Each witness went red under a `bench probe` mutation after the ticket commit, and each probe restored the file.
The plan probes `registrations`, `command-contract`, and `routing` answered `bit`, and the record holds them.

| subject | mutation | failing witness | row |
| --- | --- | --- | --- |
| `internal/conformance/guard_header_manifests_test.go` | the shim entry omitted from the manifest list | `TestGuardHeaderManifestsCoverTheRgGuard` | R41 |
| `.bench/hooks/block-unsafe-rg.sh` | a missing resolver or core exits 2 | `TestRgGuardHookFailsOpen` | R42 |
| `.bench/hooks/block-unsafe-rg.sh` | a stale core exits 2 | `TestRgGuardHookFailsOpen` | R43 |
| `cmd/bench/guards.go` | an unreadable envelope returns 2 | `TestGuardRgVerb` | R44 |
| `.bench/hooks/block-unsafe-rg.sh` | the core exit 2 maps to 0 | `TestRgGuardHookRefuses` | R45 |
| `cmd/bench/guards.go` | a blocked verdict returns 0 | `TestGuardRgVerb` | R45 |
| `.bench/hooks/block-unsafe-rg.sh` | a core error exits 2 | `TestRgGuardHookFailsOpen` | R42 |

The author traced `guard-bench-follow-on` through the tree and found no registry that ticket 03 or ticket 04 does not name.
The shim warns and exits 0 for a stale core, as `block-primary-file-write.sh` does, because R43 needs each degraded case to fail open.
