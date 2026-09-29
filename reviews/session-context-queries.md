# Session context queries review record

## QU-C1 author evidence

Ticket 1 had a fresh `bench-writer` author, `scq-t1-author`, on opus at high effort, with a cap of 3 attempts. The author started at `3cac9fca` and committed `7c258f25` on a lane pass in the first attempt. The chunk pair is `3cac9fca..7c258f25`, because `3cac9fca` holds the version 2 plan and precedes the ticket commit.

The author wrote the tests before the implementation. At that point, every selected worktree test in `internal/worktree` failed with `unknown argument: --view`, and the three command tests in `cmd/bench` failed too. `TestSelectedWorktreesPreserveDefault` passed against the unchanged bare path, which is its purpose. The author then added the selected view, and each test passed.

The author reported these deviations from the ticket and the stream source. The Spec axis grades each one.

- The bare baseline fixture holds the facts that the unchanged `ListCommand` printed at `3cac9fca`, not the stream fixture. On `main`, the active help rows now read `bench worktree path <target>` and `inspect an active worktree by its id`. The fixture also adds a `recovered` case with valid recovery metadata.
- The command test file keeps only the worktree tests of the stream file. The stream history tests belong to ticket 2. The file adds `TestSelectedWorktreeWithinResponseBound` for QU26, and one helper records each synthetic assignment.
- The help expectations read `usage.WorktreeListPaths` and the selected grammar help, not a copied literal. This change follows the duplicated-facts sweep.
- A named predicate, `selectsWorktrees`, holds the route check that the stream wrote inline in `ListCommand`.
- The CHANGELOG entry names only the worktree view. Ticket 2 owns the history entry.
- No check required an edit to these `Writes:` paths, so the diff leaves them unchanged:
  - `internal/worktree/path.go`, `internal/worktree/list_actions_test.go`, and `internal/worktree/unlanded_route_test.go`;
  - `cmd/bench/main.go`, `cmd/bench/worktree_leaves.go`, and `cmd/bench/command_registry.go`;
  - `cmd/bench/command_registry_test.go` and `cmd/bench/help_inventory_test.go`;
  - the two conformance tests, the anchor registry files, and the three canaries.
- The diff adds no `bench help` inventory row for the selected view. The stream added none, and `bench worktree --help` shows the grammar. This call is open to reviewer veto.

The author computed the plan digest from the plan rule. The input is the JSON array of the spec bytes and the three ticket bytes, in plan order. No Bench verb printed the value, so the review preflight must confirm it.

### Probe verdicts

Each probe ran through `bench probe`. Each probe bit, and each restore reads `yes`. The first row is the plan probe `1-state-probe`. The other rows are author probes for rows with no useful red before the change. One earlier QU9 probe did not compile and is not in the table.

| File | Mutation | Test | Row |
|---|---|---|---|
| `internal/worktree/list_selected.go` | swap: `string(selected.State)` to `string(intent.StateActive)` | TestSelectedWorktreeFacts | QU1 |
| `internal/worktree/list.go` | swap: `return out + help, 0` to `return help + out, 0` | TestSelectedWorktreesPreserveDefault | QU9 |
| `internal/worktree/list_selected.go` | swap: `if seenIDs[selected.ID] {` to `if false {` | TestSelectedWorktreeHostilePath, TestSelectedWorktreeAliases | QU3 |
| `internal/worktree/list_selected.go` | swap: one added help action | TestSelectedWorktreeWithinResponseBound | QU26 |
| `internal/usage/worktree.go` | omission: `WorktreeListPaths,` in the usage list | TestSelectedWorktreeHelpDiscovery | QU18 |

### Verification

The author ran each QU-C1 verification at `7c258f25`, and each passed. The JSON payload holds each result. The bare-matrix excerpt omits its two skip rows. Each skip is an environment capability skip for unix sockets. The author also ran these checks at `7c258f25`, and each passed: `bench test --package ./internal/conformance`, `bench test --package ./internal/worktree`, `bench test --package ./internal/usage`, and `bench test --package ./cmd/bench`. After the commit, `bench preflight build session-context-queries` reported 13 green checks and 0 red checks.

```bench-review-record
{
  "version": 2,
  "spec": "specs/session-context-queries/spec.md",
  "plan_digest": "sha256:48ff3619814e6377b7425a7faa8366aa66f4387386f07421a97bed71e08742b2",
  "implementation_session": "",
  "chunks": [
    {
      "id": "QU-C1",
      "base": "3cac9fcacc2a0226c4bc5264ed7bc683076e951e",
      "tip": "7c258f25d6b3f96c024acd1106ebfb9db1058b4d",
      "plan_digest": "sha256:48ff3619814e6377b7425a7faa8366aa66f4387386f07421a97bed71e08742b2",
      "source_digest": "38ce13f3f0ef7c3764e0ce9b7048bcaea8436850",
      "acceptance_rows": [
        "QU1",
        "QU2",
        "QU3",
        "QU9",
        "QU10",
        "QU16",
        "QU17",
        "QU18",
        "QU26"
      ],
      "verification": [
        {
          "id": "qu-c1-1-worktree-r1",
          "performer": "claude:bench-writer/scq-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "38ce13f3f0ef7c3764e0ce9b7048bcaea8436850",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/scq-t1-author-20260928/1-worktree@7c258f25",
            "digest": "sha256:bacc197901752f279722655a014e39b0de43d84e01c1ddbfc475fe7e1101bad4",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1049\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-worktree",
          "command": "bench test --package ./internal/worktree --run TestSelected",
          "exit_code": 0
        },
        {
          "id": "qu-c1-1-bare-matrix-r1",
          "performer": "claude:bench-writer/scq-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "38ce13f3f0ef7c3764e0ce9b7048bcaea8436850",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/scq-t1-author-20260928/1-bare-matrix@7c258f25",
            "digest": "sha256:74af00cabcc539c14a57fed97b00e4f98c66703baa06530c9cf94f9d67189e46",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,5261\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "1-bare-matrix",
          "command": "bench test --package ./internal/worktree --run 'TestList|TestPath|TestCleanLanded|TestLanded|TestUnlanded|TestParallelCensusOnTheLiveTree|TestSerialSetStaysBelowTheCeiling|TestPackage'",
          "exit_code": 0
        },
        {
          "id": "qu-c1-1-command-route-r1",
          "performer": "claude:bench-writer/scq-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "38ce13f3f0ef7c3764e0ce9b7048bcaea8436850",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/scq-t1-author-20260928/1-command-route@7c258f25",
            "digest": "sha256:c82aacba46695ccfd0e851b72e5e31dc24a9e77a6014d34cfce976b509399a83",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,2243\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-command-route",
          "command": "bench test --package ./cmd/bench --run 'TestSelected|TestCommandRegistryAXI|TestAXIRegistry|TestHelp|TestKept|TestWorktree'",
          "exit_code": 0
        },
        {
          "id": "qu-c1-1-state-probe-r1",
          "performer": "claude:bench-writer/scq-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "38ce13f3f0ef7c3764e0ce9b7048bcaea8436850",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/scq-t1-author-20260928/1-state-probe@7c258f25",
            "digest": "sha256:c559177136fb05e5aa5d1d92ffb2e9eccde2316a06566cf1f1694d306cbbce3b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/list_selected.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestSelectedWorktreeFacts,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,168"
          },
          "requirement": "1-state-probe",
          "command": "bench probe internal/worktree/list_selected.go --swap 'string(selected.State)' --with 'string(intent.StateActive)' --package ./internal/worktree --run TestSelectedWorktreeFacts",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/scq-t1-author-20260928/1-state-probe@7c258f25",
              "digest": "sha256:c559177136fb05e5aa5d1d92ffb2e9eccde2316a06566cf1f1694d306cbbce3b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/list_selected.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestSelectedWorktreeFacts,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,168"
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
