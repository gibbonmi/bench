# Bounded default output review record

## BO-C1 author evidence

Ticket 1 had a fresh `bench-writer` author on opus at high effort, with a cap of 3 attempts. The author started at `c24229b1` and committed `53dc2f50` and `77870627`. The second commit replaced the BO68 memory measurement, because the wait status of a Go child reports the peak of the parent test process.

The orchestrator then committed `463b4908`. That commit states `BENCH_KIT` on tickets 1 and 10, because review preflight was red on `kit-pin` after the new system test file existed. The author ran the three ticket checks again at `463b4908`, and each check passed. The JSON payload holds each result.

The author reported these deviations. The Spec axis graded each one.

- The wait delay is 3 seconds, because a 2-second constant trips the bounds duplicate-owner check. The interrupt grace moved from 2 to 3 seconds with it. The spec fixes no value.
- Outside a repository, a spill goes to `responses/none/primary/`.
- On the create-failure route, an unterminated final line gets one newline before the `spill-failed` line.

### Probe verdicts

Each probe ran through `bench probe`. Each probe bit, and each restore reads `yes`.

| File | Mutation | Test |
|---|---|---|
| `cmd/bench/command_registry.go` | omission: the owner wrap in `Command.Run` | TestDispatcherBoundsPublicResponse |
| `internal/responsebound/owner.go` | swap: `tailLines` 5 to 4 | TestOwnerProjectsHeadAndTail |
| `internal/responsebound/` | swap: the bound test `>` to `>=` | TestDispatcherPassesBoundaryResponse |
| `internal/responsebound/owner.go` | swap: `Finish` starts a spill | TestOwnerEmptyResponse |
| `cmd/bench/command_registry.go` | swap: the dispatcher returns exit 0 | TestDispatcherKeepsExitCode |
| `internal/responsebound/owner.go` | swap: every write goes to stdout | TestExecGrammarRefusalKeepsUsageLine |

### Verification

The whole-project gate was green at `77870627`. The author's three ticket checks passed at `463b4908`.

## BO-C1 chunk review, round 1

The frozen pair is base `80780c046df2796b635d0c139dae3b505d2891f1` and tip `463b49086757cde37b79d23812289e4df318250e`. The shared evidence is `sha256:e5c3c3d80f7db0b0e8c7c67ec88b218d0331ae1246d7ab426267816ae3fb110b`. Each axis ran in a fresh `bench-reviewer` session on opus at medium effort, by user direction. Only the Coverage axis ran probes, and it left the tree clean.

The consumer table has these rows outside the diff:

- The test callers of `Command.Run` in `cmd/bench`. Every public entry except `worktree exec` is `pending`, so their output does not change.
- Four `runWorktreeChild` calls in `internal/worktree/exec_test.go`.
- The race phase in `internal/gate/gate_go.go`, which reads `racetests.Tests`.

The whole-project gate at `77870627` ran each of these consumers green.

The raw finding count is 6: Standards 2, Spec 1, and Coverage 3. No two findings name the same fix, so 6 repair targets remain.

## Standards

Findings: 2. The worst issue is that the worktree leaf table is stated at two sites.

- `cmd/bench/main.go:128` and `cmd/bench/worktree_leaves.go:69` each name `worktreeLeaves`. The bound check reads the registry field, and the family dispatch reads its own argument. A family that omits `Leaves:` escapes the disposition test. Target R1. `auto-fix`. Confidence 7.
- `cmd/bench/response_bound_test.go:17-18`, `internal/responsebound/owner_test.go:17-19`, and `internal/systemtest/exec_bound_test.go:18-19` restate 10, 4, 5, and `responses/none/primary` apart from their constants. `AGENTS.md` allows an independent expectation only with a recorded red. Target R2. `auto-fix`. Confidence 5.

## Spec

Findings: 1. The worst issue is a coverage-map seam that does not exist.

- `specs/ft336-bounded-output/spec.md:306` cites "existing TestWorktreeExecGrammar tests in internal/worktree". No such test exists. The behavior holds through `TestExecGrammarRefusalKeepsUsageLine` at `cmd/bench/response_bound_test.go:98`. Target R3: the orchestrator amends the seam cell under the plan-expansion policy. `auto-fix`. Confidence 8.

The Spec axis held the other 23 rows. BO20 lists only `internal/responsebound/owner.go:20,102`. BO68 holds: after the spill starts, the owner keeps only the head lines and a 5-slot ring.

## Coverage

Findings: 3. The worst issue is a merged line on the create-failure route.

- After a create failure, stdout ends without a newline and stderr writes a line. The `spill-failed` text then joins the stdout line, because the separator in `internal/responsebound/owner.go:168-173` reads the combined line state. An injected test failed with `partialspill-failed{reason=...}`. Target R4. `auto-fix`. Confidence 7.
- A symlink at `responses/` or at `responses/<repo-key>/` has no test. `TestSpillStoreRefusesSymlink` plants it only at the scope level, and a last-level-only check stays silent. `spec.md:352` says each symlink edge has a row. Target R5. `auto-fix`. Confidence 8.
- No owner test splits one line across two writes or sends several lines in one write. Two mutations of `lines.go` stayed silent. The production code is correct. Target R6. `auto-fix`. Confidence 7.

## Advice

- `Command.Run` calls `owner.Finish()` without `defer`, so a panic in a bounded verb loses the held output.
- `Finish` ignores the error from `file.Close()`.
- The spec calls `none` a repo key at line 163 and a scope at line 165. Ticket 3 prunes both scopes.
- BO69 allows up to 24.5 seconds, not the 3-second wait delay.
- The owner package comment says that the dispatcher and exec "cannot drift apart", but exec reaches the owner only through the dispatcher.
- The BO26 and BO68 system tests check only the `spilled{` prefix, not the `lines=` value.

## BO-C1 repair routing

Each repair goes to one fresh `bench-writer` repair session for ticket 1, whose `Writes:` line holds every path. This is cycle 1 of the two repair cycles for chunk BO-C1. Targets R5 and R6 are the chunk's one hardening cycle.

| Target | Ticket | Repair |
|---|---|---|
| R1 | 1 | Give the leaf family table one source that both the bound check and the dispatch read. |
| R2 | 1 | Derive each expectation from its constant, or record a demonstrated red for each independent value. |
| R3 | orchestrator | Amend the BO31 seam cell to the dispatcher test and the existing exec grammar tests. |
| R4 | 1 | Track the stdout line state for the create-failure separator, and add the row test. |
| R5 | 1 | Test a symlink and a regular file at each store level. |
| R6 | 1 | Add owner tests for a split line and for one multi-line write, with the full spill line. |

```bench-review-record
{
  "version": 2,
  "spec": "specs/ft336-bounded-output/spec.md",
  "plan_digest": "sha256:0d97c1e30b257bb8ff01f4de56e407176319611d19d78351521263313fc176c4",
  "implementation_session": "",
  "chunks": [
    {
      "id": "BO-C1",
      "base": "80780c046df2796b635d0c139dae3b505d2891f1",
      "tip": "463b49086757cde37b79d23812289e4df318250e",
      "plan_digest": "sha256:0d97c1e30b257bb8ff01f4de56e407176319611d19d78351521263313fc176c4",
      "source_digest": "9d651c1386edaf2c287c2bca4c55d1af34bd457c",
      "acceptance_rows": [
        "BO1",
        "BO2",
        "BO3",
        "BO4",
        "BO5",
        "BO6",
        "BO7",
        "BO14",
        "BO15",
        "BO18",
        "BO19",
        "BO20",
        "BO21",
        "BO22",
        "BO23",
        "BO24",
        "BO25",
        "BO26",
        "BO27",
        "BO28",
        "BO29",
        "BO31",
        "BO68",
        "BO69"
      ],
      "verification": [
        {
          "id": "bo-c1-1-owner-r1",
          "performer": "claude:bench-writer/bo-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9d651c1386edaf2c287c2bca4c55d1af34bd457c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/bo-t1-author-20260924/1-owner@463b4908",
            "digest": "sha256:7480c8dfb0821726a5aaf1092d4aceb1d927c6f199d74443be20caf83ee4394a",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/responsebound,pass,12\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-owner",
          "command": "bench test --package ./internal/responsebound",
          "exit_code": 0
        },
        {
          "id": "bo-c1-1-cmd-r1",
          "performer": "claude:bench-writer/bo-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9d651c1386edaf2c287c2bca4c55d1af34bd457c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/bo-t1-author-20260924/1-cmd@463b4908",
            "digest": "sha256:402d90541cf4d165740601ad7030226fa0c8da301ee1811331294640dfdc4458",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14774\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "bo-c1-1-system-r1",
          "performer": "claude:bench-writer/bo-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9d651c1386edaf2c287c2bca4c55d1af34bd457c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/bo-t1-author-20260924/1-system@463b4908",
            "digest": "sha256:176f078bbe982576ae3cfd5e59c2f77293863920c9e2d78e47ade8a98a5ba2a7",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,59534\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-system",
          "command": "bench test --check system",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "bo-c1-r1-standards",
          "performer": "claude:bench-reviewer/bo-c1-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9d651c1386edaf2c287c2bca4c55d1af34bd457c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/bo-c1-standards@463b4908",
            "digest": "sha256:325f60e100fa115ed43e5ac78976529c4e5915897b54ec4b95b50a6c3f5d40e8",
            "excerpt": "Standards: 2 findings. Worst: the worktree leaf table is stated at two sites, so the bound check and the dispatcher can read different tables."
          },
          "axis": "Standards",
          "base": "80780c046df2796b635d0c139dae3b505d2891f1",
          "tip": "463b49086757cde37b79d23812289e4df318250e",
          "finding_ids": [
            "R1",
            "R2"
          ],
          "supersedes": []
        },
        {
          "id": "bo-c1-r1-spec",
          "performer": "claude:bench-reviewer/bo-c1-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9d651c1386edaf2c287c2bca4c55d1af34bd457c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/bo-c1-spec@463b4908",
            "digest": "sha256:36308cbbc9e18149a678953d031fddc142bd5db61fe905c19d7206b0ebe13b7d",
            "excerpt": "Spec: 1 finding. Worst: the BO31 seam cell names TestWorktreeExecGrammar, which does not exist; the behavior holds through TestExecGrammarRefusalKeepsUsageLine."
          },
          "axis": "Spec",
          "base": "80780c046df2796b635d0c139dae3b505d2891f1",
          "tip": "463b49086757cde37b79d23812289e4df318250e",
          "finding_ids": [
            "R3"
          ],
          "supersedes": []
        },
        {
          "id": "bo-c1-r1-coverage",
          "performer": "claude:bench-reviewer/bo-c1-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9d651c1386edaf2c287c2bca4c55d1af34bd457c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/bo-c1-coverage@463b4908",
            "digest": "sha256:6b6e7b0f46291ad514559b627698c4324ede7829f46cf48cbff460c79062fe68",
            "excerpt": "Coverage: 3 findings. Worst: on the create-failure route, an unterminated stdout line followed by a stderr write merges the spill-failed text into that stdout line."
          },
          "axis": "Coverage",
          "base": "80780c046df2796b635d0c139dae3b505d2891f1",
          "tip": "463b49086757cde37b79d23812289e4df318250e",
          "finding_ids": [
            "R4",
            "R5",
            "R6"
          ],
          "supersedes": []
        }
      ]
    }
  ],
  "completion": {
    "state": "pending"
  }
}
```
