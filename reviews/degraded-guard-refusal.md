# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/degraded-guard-refusal/spec.md",
  "plan_digest": "sha256:a4fb59f184d65ef47e3267596181bcbb65d8e91935385d4eba7a3659b62497e5",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GR1",
      "base": "8ce8934d629efacb52a3b0c714afbebacf5f3435",
      "tip": "a9324deb788d6de6d04313a79aad653246df064c",
      "plan_digest": "sha256:a4fb59f184d65ef47e3267596181bcbb65d8e91935385d4eba7a3659b62497e5",
      "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
      "acceptance_rows": [
        "GR01",
        "GR02",
        "GR03",
        "GR04",
        "GR05",
        "GR06",
        "GR07",
        "GR08",
        "GR09",
        "GR10",
        "GR11",
        "GR12",
        "GR13",
        "GR14",
        "GR15",
        "GR16",
        "GR17",
        "GR18",
        "GR19",
        "GR20",
        "GR21",
        "GR22",
        "GR23",
        "GR24",
        "GR25",
        "GR26",
        "GR27",
        "GR28",
        "GR29",
        "GR30",
        "GR31",
        "GR32",
        "GR33",
        "GR34",
        "GR35",
        "GR36",
        "GR37",
        "GR38",
        "GR39",
        "GR40",
        "GR41"
      ],
      "verification": [
        {
          "id": "t1-t1-gitguard",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:794e1801dee75687a33d3269d2e6dea21945b0c79d22174cae45ee3438087d2b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gitguard,pass,15172,230\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-gitguard",
          "command": "bench test --package ./internal/gitguard",
          "exit_code": 0
        },
        {
          "id": "t1-t1-core-entry",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:99d2631b484c34169f0ef841b87ef22f36f9348f35febe4d19f30b16f24ecedd",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18034,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-core-entry",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1-t1-hook-entries",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:c74db1d570b9c787307a123e64aaf9d33a69829c1bc9972922440cbf61fd122e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,175,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,79095,175\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-hook-entries",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t1-t1-hook-contracts",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:31993f343711a7332b458fd41dcd6d63a6725323a4e441a6e825f5ae9b0a75c7",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  package-core-guard,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,2775,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-hook-contracts",
          "command": "bench test --check package-core-guard",
          "exit_code": 0
        },
        {
          "id": "t1-t1-author-parse-success-omission",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:55d3864749e19a9d47d2744a9282cb03b5dd5b8198131f9d6fadbfcfed6391fc",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.bench/hooks/block-dangerous-git.sh,swap,failed,1,TestDegradedGuardRefusesMalformedCommands,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,TestDegradedGuardRefusesMalformedCommands,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,438\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,381,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestDegradedGuardRefusesMalformedCommands,\"guard_refusal_test.go:46: GR03 trailing escape after an attached separator exit = 0 (stderr \\\"\\\"), want 2\",4\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-author-parse-success-omission",
          "command": "bench test --check system",
          "exit_code": 0,
          "probe": {
            "mutation": "restore permissive malformed success at the trailing-escape branch; GR03 must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
              "digest": "sha256:55d3864749e19a9d47d2744a9282cb03b5dd5b8198131f9d6fadbfcfed6391fc",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,.bench/hooks/block-dangerous-git.sh,swap,failed,1,TestDegradedGuardRefusesMalformedCommands,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,system,TestDegradedGuardRefusesMalformedCommands,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,438\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,381,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/systemtest,TestDegradedGuardRefusesMalformedCommands,\"guard_refusal_test.go:46: GR03 trailing escape after an attached separator exit = 0 (stderr \\\"\\\"), want 2\",4\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-t1-canonical-path-policy",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:154f593d7d73b3e990af6a23d8a1cc5ef43f8ab8ad56d6919fc97a013bcbbdd3",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  canonical-path-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,115,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-canonical-path-policy",
          "command": "bench test --check canonical-path-owner",
          "exit_code": 0
        },
        {
          "id": "t1-t1-injected-port-policy",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:9a49626b0abd49fcae76c5e872ceff75e9558cbd53a3c2777967bc1f92d1f32f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  injected-port-registry,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,25,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-injected-port-policy",
          "command": "bench test --check injected-port-registry",
          "exit_code": 0
        },
        {
          "id": "t1-t1-guard-anchors",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:c6b23801a831b75949f1915363bdd979641a7199ea749c96f92a6fd5c1764263",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4975,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-guard-anchors",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t1-t1-axi-registry",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:9bb0abdb3205d48eb90e52b735058f8f1eb39c284c726525d8897e6fe96907e6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  axi-query-registry,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,12,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-axi-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "t1-t1-command-routing",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:5f31975f6976fd1476a5cae53a7c426194874510d089b2fb2bfcc648371acd30",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  subcommand-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,28,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-command-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "t1-t1-root-conformance",
          "performer": "claude:ft366-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9f6b558c91eb76ac2970a4edb5899bad31ab26ca",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft366-t01-author-20261019",
            "digest": "sha256:8866917ee404d739a2f91712b4c4e856906d64986d2456ea5dcccaf6320c789e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16401,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
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
