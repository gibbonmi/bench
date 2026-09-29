# Tree targets review record

## TT-C1 author evidence

Ticket 1 had a fresh `bench-writer` author, `tt-t1-author`, on opus at high effort, with a cap of 3 attempts. The author started at `1634be1e` and committed `166a3a8f` on a lane pass in the first attempt. The author then committed this record in a second commit.

The chunk pair is `f981cd3d..166a3a8f`. The base is the `main` tip, because a plan commit is never a chunk base. The payload names the ticket commit as the tip, because a record cannot name its own commit. The coordinator moves the tip to the record commit. The source digest does not change, because the record file is outside the graded source.

The author wrote each test before the code that it grades. The red and green log for each row follows:

- TT6: before the first production edit, `TestRepositoryVerbRefusesTreeTarget/version` failed. The verb printed `bench dev (linux/amd64)` at exit 0. After the refusal in `cmd/bench/tree_scope.go`, the subtest passed.
- TT7: before the change, the `idea` subtest passed, because the idea grammar refuses `--in` itself. A probe made the idea parser skip its first two arguments, and the subtest failed with `parked: x` at exit 0. After the change, the same probe was silent, because the dispatcher refuses before the verb runs.
- TT2 to TT5: with an empty check, `TestCommandScopeCheckBites` failed with no scope diagnostic. After the check, the test passed.
- TT1: after the check and before the registry declarations, `TestRootConformance` failed with `gate: command "anchors" declares no scope`. The named check with the worktree build gave the same red. After the declarations, both passed.

The author reported these deviations from the ticket. The Spec axis grades each one.

- The scope check reads a classified entry that is not internal as public. The help-inventory check refuses any third classification. So the scope check reuses `parityInternalCommand` and adds no second test for `publicInventory`. This change follows the duplicated-facts sweep.
- The check reads the presence of a `Scope` field and does not grade its value. The scope type has no named zero value, so an undeclared value needs an explicit conversion.
- A tree without the leaf file grades the registry only. The routing check applies the same rule to an absent package, and the canary fixture copies only `cmd/bench/main.go`.
- Each registry line takes its scope after its bound disposition. So the AXI mutation anchors still match, and `internal/conformance/axi_query_registry_test.go` stays unchanged.
- No check required an edit to these `Writes:` paths, so the diff leaves them unchanged:
  - `cmd/bench/command_registry_test.go` and `cmd/bench/help_inventory_test.go`;
  - `internal/conformance/axi_query_registry_test.go` and `internal/conformance/subcommand_routing_table_test.go`.
- The canary fixture change is in `MUTATE.json`, because that file quotes the `learnings` registry line byte for byte.
- The TT7 red probe changed `internal/roadmap/roadmap.go`, which is outside the fence. The probe restored the file, and the diff does not touch it.

The author computed the plan digest from the plan rule. The input is the JSON array of the spec bytes and the five ticket bytes, in plan order. The same rule gives the digest `48ff3619` for the session context queries record at `a5a527c8`, which confirms the method. The source digest is the tree of `166a3a8f` without this record file.

### Probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`. The first row is the plan probe `1-scope-probe`. The JSON payload holds the exact command and output of the first three rows. The last two rows ran on the source before the ticket commit or record the silent result, so only this table holds them.

| File | Mutation | Test | Row | Verdict |
|---|---|---|---|---|
| `cmd/bench/tree_scope.go` | swap: `"--in"` to `"--in-x"` | TestRepositoryVerbRefusesTreeTarget | TT6 | bit |
| `internal/conformance/command_scope_test.go` | swap: `if len(leaf.fields["Scope"]) == 0 {` to `if len(leaf.fields["Scope"]) < 0 {` | TestCommandScopeCheckBites | TT3 | bit |
| `internal/conformance/command_scope_test.go` | swap: `case len(entry.fields["Inventory"]) == 1 && !declared:` to `case len(entry.fields["Inventory"]) == 2 && !declared:` | TestCommandScopeCheckBites | TT2 | bit |
| `internal/roadmap/roadmap.go` | swap: `usage.Parse(ideaGrammar, args)` to `usage.Parse(ideaGrammar, args[min(2, len(args)):])`, before the change | TestRepositoryVerbRefusesTreeTarget/idea | TT7 | bit |
| `internal/roadmap/roadmap.go` | the same swap, after the change | TestRepositoryVerbRefusesTreeTarget/idea | TT7 | silent |

### Verification

The author ran each TT-C1 verification on the source of `166a3a8f`, and each passed. The JSON payload holds each result. The conformance excerpt omits its three skip rows. Each skip is an environment capability skip for unix sockets or device nodes. The author also ran these checks, and each passed: `bench structure --growth f981cd3d`, `bench canary`, and the named checks `canary-fixture-compliance`, `package-core-guard`, `conformance-canary-families`, and `axi-query-registry`. After the commit, `bench preflight build tree-targets` reported 14 green checks, 1 check that does not apply, and 0 red checks.

```bench-review-record
{
  "version": 2,
  "spec": "specs/tree-targets/spec.md",
  "plan_digest": "sha256:7355203a687409e19b68ed6613479e382ad1682b9ccc27fc2000604a45e2920c",
  "implementation_session": "",
  "chunks": [
    {
      "id": "TT-C1",
      "base": "f981cd3db1a4f27eddba5feede40a80e205bf1c1",
      "tip": "166a3a8f6cb850de7a5e15b5b273dd415b46ff17",
      "plan_digest": "sha256:7355203a687409e19b68ed6613479e382ad1682b9ccc27fc2000604a45e2920c",
      "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
      "acceptance_rows": [
        "TT1",
        "TT2",
        "TT3",
        "TT4",
        "TT5",
        "TT6",
        "TT7"
      ],
      "verification": [
        {
          "id": "tt-c1-1-cmd-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/1-cmd@166a3a8f",
            "digest": "sha256:b8706b37336cad03b084fddd17aff5f746edb6b943cfc5c83ebfe709a7766696",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16090\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "tt-c1-1-routing-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/1-routing@166a3a8f",
            "digest": "sha256:915a63e3adf95500b076f300919bf9a71615b398618b71408f9c5bc26dab14d8",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,20\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "tt-c1-1-conformance-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/1-conformance@166a3a8f",
            "digest": "sha256:2ccce730465fdae97c867b5b5ecec44425db5d7d0e68002847654e251e563f68",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,51673\nfailures[0]{package,test,line}:\nskips[3]{package,test,reason}:"
          },
          "requirement": "1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "tt-c1-1-scope-probe-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/1-scope-probe@166a3a8f",
            "digest": "sha256:ca4cb4929604998da4cf09f2ae158ca735484565de97eef0e59bd357d08297ab",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestRepositoryVerbRefusesTreeTarget,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,55"
          },
          "requirement": "1-scope-probe",
          "command": "bench probe cmd/bench/tree_scope.go --swap '\"--in\"' --with '\"--in-x\"' --package ./cmd/bench --run TestRepositoryVerbRefusesTreeTarget",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t1-author-20260929/1-scope-probe@166a3a8f",
              "digest": "sha256:ca4cb4929604998da4cf09f2ae158ca735484565de97eef0e59bd357d08297ab",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestRepositoryVerbRefusesTreeTarget,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,55"
            }
          }
        },
        {
          "id": "tt-c1-author-probe-tt3-leaf-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/tt3-leaf-probe@166a3a8f",
            "digest": "sha256:d52c3fad433dc890602431309f3217fb4f9944bad20b8cb108351e881f4f6f51",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/conformance/command_scope_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommandScopeCheckBites,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5"
          },
          "requirement": "author-probe-TT3-leaf",
          "command": "bench probe internal/conformance/command_scope_test.go --swap 'if len(leaf.fields[\"Scope\"]) == 0 {' --with 'if len(leaf.fields[\"Scope\"]) < 0 {' --package ./internal/conformance --run TestCommandScopeCheckBites",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t1-author-20260929/tt3-leaf-probe@166a3a8f",
              "digest": "sha256:d52c3fad433dc890602431309f3217fb4f9944bad20b8cb108351e881f4f6f51",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/conformance/command_scope_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommandScopeCheckBites,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5"
            }
          }
        },
        {
          "id": "tt-c1-author-probe-tt2-public-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/tt2-public-probe@166a3a8f",
            "digest": "sha256:3064fcd3c3bc49a39214018ad6c41a263ea660bb39faaaab910d795894018d0a",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/conformance/command_scope_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommandScopeCheckBites,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10"
          },
          "requirement": "author-probe-TT2-public",
          "command": "bench probe internal/conformance/command_scope_test.go --swap 'case len(entry.fields[\"Inventory\"]) == 1 && !declared:' --with 'case len(entry.fields[\"Inventory\"]) == 2 && !declared:' --package ./internal/conformance --run TestCommandScopeCheckBites",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t1-author-20260929/tt2-public-probe@166a3a8f",
              "digest": "sha256:3064fcd3c3bc49a39214018ad6c41a263ea660bb39faaaab910d795894018d0a",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/conformance/command_scope_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommandScopeCheckBites,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10"
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
