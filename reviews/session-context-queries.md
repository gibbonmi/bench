# Session context queries review record

## QU-C1 author evidence

Ticket 1 had a fresh `bench-writer` author, `scq-t1-author`, on opus at high effort, with a cap of 3 attempts. The author started at `3cac9fca` and committed `7c258f25` on a lane pass in the first attempt. The author then committed this record at `a5a527c8`.

The chunk pair is `12c7d857..a5a527c8`. The base is the `main` tip of the integration source, because a plan commit is never a chunk base. The coordinator corrected the author's first pair, `3cac9fca..7c258f25`, to this pair. The source digest does not change, because the record file is outside the graded source.

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

## QU-C1 chunk review, round 1

The frozen pair is base `12c7d857c6fe14e03c618c1584d3374371068277` and tip `a5a527c8425cafe6eac17c7515bd328e6ba7d65d`. The shared evidence is `sha256:c03b0356cdb7638e79f9793a96932a9d3d457204a77ffdc8b58e11b0e238243e`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort, on the conditional review line. Only the Coverage axis ran probes, and it left the tree clean.

Every consumer of `ListCommand` outside the diff passes `nil` or a bare argument list. These consumers are `cmd/bench/worktree_leaves.go:37`, `unlanded_route_test.go:53`, and the landed, path identifier, request token, and hostile clean tests. The two matrix callers at `list_actions_test.go:98` and `:120` also stay bare. The `1-bare-matrix` verification ran each of them green.

The raw finding count is 6: Standards 3, Spec 0, and Coverage 3. The two Coverage record findings name one fix, so 5 repair targets remain.

## Standards

Findings: 3. The worst issue is that the selected view restates refusal text that the package already owns.

- `internal/worktree/list_selected.go:47` repeats the assignment-read refusal pair of `internal/worktree/list.go:52`. `list_selected.go:60` adds a copy of `target contains control characters` from `internal/worktree/path.go:68`. `AGENTS.md` requires one source for each fact, and the copy at `internal/worktree/merge.go:177` is outside the ticket fence and stays. Target R1. `auto-fix`. Confidence 6.
- `internal/worktree/list_selected_test.go:141`, `:248`, `:290`, and `:309` restate selector error fragments, the bare usage line, the command name, and the path refusal apart from their owners. No red is recorded for them. Target R2. `auto-fix`. Confidence 5.
- The comment at `cmd/bench/selected_queries_test.go:71-73` restates the red record of QU26 from `specs/session-context-queries/spec.md:203`. `bench-craft-comments` gives the red record to the spec. Target R3. `auto-fix`. Confidence 4.

## Spec

Findings: 0. All nine rows of ticket 1 are met: QU1, QU2, QU3, QU9, QU10, QU16, QU17, QU18, and QU26. All seven author deviations agree with the spec and the ticket. The missing `bench help` row for the selected view follows the tree convention, because the other worktree subcommands appear only through the `worktree --help` row.

## Coverage

Findings: 3. The worst issue is that no test grades a `--target`-only request.

- Spec line 59 starts the selected route on `--view` or `--target`. The `--target` cases at `internal/worktree/list_selected_test.go:234` and `:242` assert only exit 2 and a `usage:` prefix, which the bare grammar refusal also prints. A probe that drops the `--target` half of `selectsWorktrees` at `internal/worktree/list_selected.go:26` stayed silent in both packages. Target R4. `auto-fix`. Confidence 9.
- The probe row labelled QU18 omits `WorktreeListPaths` and grades help discovery, not the detail route. That omission stays silent against `TestSelectedWorktreeDetailRoute`, and a detail-route swap bit. Target R5. `auto-fix`. Confidence 9.
- The QU26 probe row gives no old or new text, and the four author probe rows have no stored command or output in the payload. Nobody can rerun them. Target R5. `auto-fix`. Confidence 7.

## QU-C1 repair routing

Each repair goes to one fresh `bench-writer` repair session for ticket 1, whose `Writes:` line holds every path. This is cycle 1 of the two repair cycles for chunk QU-C1.

| Target | Ticket | Repair |
|---|---|---|
| R1 | 1 | Give the selected view the refusal owners of `list.go` and `path.go` instead of copies. |
| R2 | 1 | Derive each test expectation from its owner, or record a demonstrated red for each independent value. |
| R3 | 1 | Keep the derivation of 6 lines in the comment, and remove the red argument. |
| R4 | 1 | Assert the exact selected grammar usage for each `--target`-only case, and show that the dropped half now bites. |
| R5 | 1 | Record each author probe as a verification entry with its exact `bench probe` command and output, and relabel the QU18 row with the detail-route mutation. |

```bench-review-record
{
  "version": 2,
  "spec": "specs/session-context-queries/spec.md",
  "plan_digest": "sha256:48ff3619814e6377b7425a7faa8366aa66f4387386f07421a97bed71e08742b2",
  "implementation_session": "",
  "chunks": [
    {
      "id": "QU-C1",
      "base": "12c7d857c6fe14e03c618c1584d3374371068277",
      "tip": "a5a527c8425cafe6eac17c7515bd328e6ba7d65d",
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
      "reviews": [
        {
          "id": "qu-c1-standards-r1",
          "performer": "claude:bench-reviewer/qu-c1-standards-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "38ce13f3f0ef7c3764e0ce9b7048bcaea8436850",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/qu-c1-standards-r1@a5a527c8",
            "digest": "sha256:4ff2c584a52b6d2d2998ea447e92b85682936583a6a5b8602b03bb842d7954c6",
            "excerpt": "Standards: 3 findings. Worst: the selected view restates refusal text that list.go and path.go already own."
          },
          "axis": "Standards",
          "base": "12c7d857c6fe14e03c618c1584d3374371068277",
          "tip": "a5a527c8425cafe6eac17c7515bd328e6ba7d65d",
          "finding_ids": [
            "R1",
            "R2",
            "R3"
          ],
          "supersedes": []
        },
        {
          "id": "qu-c1-spec-r1",
          "performer": "claude:bench-reviewer/qu-c1-spec-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "38ce13f3f0ef7c3764e0ce9b7048bcaea8436850",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/qu-c1-spec-r1@a5a527c8",
            "digest": "sha256:d11a6a58a1f35572c613c75ae3269aeb30b70dd8bd33238c65209ca5074bf475",
            "excerpt": "Spec: 0 findings. All nine ticket 1 rows are met, and all seven author deviations agree with the spec."
          },
          "axis": "Spec",
          "base": "12c7d857c6fe14e03c618c1584d3374371068277",
          "tip": "a5a527c8425cafe6eac17c7515bd328e6ba7d65d",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "qu-c1-coverage-r1",
          "performer": "claude:bench-reviewer/qu-c1-coverage-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "38ce13f3f0ef7c3764e0ce9b7048bcaea8436850",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/qu-c1-coverage-r1@a5a527c8",
            "digest": "sha256:e9e280188dab2fea2f29132e8d938b8016656a614d79b15c24fb43056fb67f35",
            "excerpt": "Coverage: 3 findings. Worst: no test grades that --target alone starts the selected route, so a mutation that drops that half of selectsWorktrees stays silent."
          },
          "axis": "Coverage",
          "base": "12c7d857c6fe14e03c618c1584d3374371068277",
          "tip": "a5a527c8425cafe6eac17c7515bd328e6ba7d65d",
          "finding_ids": [
            "R4",
            "R5"
          ],
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
  }
}
```
