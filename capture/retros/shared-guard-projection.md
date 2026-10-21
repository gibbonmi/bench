## Outcome

The `shared-guard-projection` spec landed at `0bcfa8363c04fe84705dd5b28aacd3450fb53ec5` with a green whole-project gate. The `shellcommand` package now owns the wrapper and worktree projections and the hook envelope decoding. The Git, Bench, and Write guards keep their own policy on top of it. All 74 acceptance rows reconciled, and the 11 final requirements passed at `4e8cf2af`.

The run was a `--full --auto-approve` build with an author limit of 1 and three chunks. Every author and repair ran on Opus/medium by reviewer direction, and Sonnet/high ran every review axis. An Opus/medium and Sonnet/medium consultant pair decided each judgment call, and a Fable/medium tie-break decided one split. The GP46 golden came from a reference worktree at the build-entry commit, and the candidate matched it on all 111 members.

## Gate-stage timings

- landing: commit 0bcfa8363c04fe84705dd5b28aacd3450fb53ec5, trace b97a818fa2f2de078130d882d0caa78c
- gofmt: 163 ms
- vet: 1433 ms
- test: 175770 ms
- race: 2376 ms
- system: 101794 ms
- shellcheck: 556 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 1 | the floor checks pass and the Git verdicts are unchanged | verified | 9 | true | opus / medium / implementer |
| ticket 2 | the Git decoder-form plan probe bites and restores | verified | 10 | true | opus / medium / implementer |
| ticket 3 | both readers keep their error text and failure posture | verified | 9 | true | opus / medium / implementer |
| ticket 4 | no private wrapper or exec grammar remains in benchguard | verified | 9 | true | opus / medium / implementer |
| ticket 5 | renamed private copies red the registered check | verified | 9 | true | opus / medium / implementer |
| ticket 6 | GP45 leaves the sentinel absent through every actual entry | verified | 8 | true | opus / medium / implementer |
| ticket 7 | the golden came only from the reference pair | verified | 9 | true | opus / medium / implementer |
| GP1 | expected repair rounds: 1 | declared | 5 | true | opus / medium / orchestrator |
| GP2 | expected repair rounds: 2 | declared | 5 | false | opus / medium / orchestrator |
| GP3 | expected repair rounds: 1 | declared | 5 | true | opus / medium / orchestrator |

Brier mean: 0.084 over 10 pairs. Abstentions: 0.

Each ticket fit its slice. The ticket 5 author peaked at about 218k tokens against the 200k estimate, and the other authors peaked between 75k and 160k. Every author used 1 attempt, except the ticket 7 author, which used 3 of 4 to fit its golden within 400 lines. The ticket 3 author first gave its confidences on a 0 to 100 scale.

## Coordinator catches

- At build entry, the coordinator removed five `(new)` markers from shared test files. The `writes-resolve` check refused that, and a second plan commit restored them.
- The ticket 1 author recorded no excerpt for its probes, because each probe was not the call just before `bench record excerpt`. The author reran both probes and then recorded.
- A coordinator probe on `spanInvokesBench` read silent and exposed the two-level depth gap that the GP2 Coverage axis then raised as C1.
- A coordinator respell probe closed the GP2 Standards finding F1 as evidence only.
- A coordinator rerun of the hook `eval` mutant closed the one GP45 entry that the ticket 6 author had not seen red.
- After a comment-only fix in GP2, the coordinator refroze the chunk, and that refreeze would orphan 20 verifications. A `bench worktree reset` removed it before any commit.
- The complete checkpoint then refused a stale review chain gap at the same comment commit. The coordinator moved the GP3 base to the GP2 tip and ran one more round of three axes.
- The coordinator did not run the `bench worktree build` and `doctor` rehearsal before the first landing, although the spec changed Bench build inputs. The landing refresh completed.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-share-git-command-projections.md | 1 | delegate-error |
| 2-add-envelope-owner-for-git.md | 0 | none |
| 3-migrate-bench-and-write-readers.md | 1 | delegate-error |
| 4-migrate-bench-projection-consumers.md | 1 | delegate-error |
| 5-enforce-guard-projection-ownership.md | 1 | one-source |
| 6-prove-guard-entry-regression.md | 0 | none |
| 7-prove-guard-migration-differential.md | 0 | none |

The ticket 1 repair added two negative projection cases and split one comment. The ticket 3 and ticket 4 repairs added caller cases for cwd independence, wrapper depth, and the outer-only exec exception. The ticket 5 repair moved the walk and filter that two owner checks shared into one helper file.

## Agent-experience improvements

### Bench CLI

- Census account: 33 raw calls. The 15 rg calls, 4 cat calls, 2 sed calls, and 1 ls call are sanctioned reads and spill copies. The 7 python3 calls rewrote the plan fence or wrote mutant and assessment JSON, so they are an owned gap. The 4 cd calls tried the pool path, and the hook refused each one.
  Feeds: none
- Add a `bench spec` form that writes the version 2 delegate amendment and its floor rows, so that no session edits the plan fence with python3.
  Feeds: new
- Add an `--out <path>` option to `bench probe`, because four authors copied spilled answers by hand after the hook refused a redirect.
  Feeds: new
- Let `bench record ticket` replace a ticket's rows after a refreeze, instead of a refusal on an existing evidence ID.
  Feeds: new
- Make the pool-path hook refusal name `bench worktree path` for a file read, so that a delegate does not retry `cd`.
  Feeds: new
- Make the review-chain check in `internal/reviewrecord/coverage.go` prove a comment-only gap, as the chunk checkpoint already does.
  Feeds: new

### Skills

- Align the bounded repair policy with the grader, so that a comment-only correction after a chunk freeze takes no refreeze.
  Feeds: new
- In the author charge template, say that each plan probe runs immediately before its `bench record excerpt` call.
  Feeds: new

### Process

- Run the `bench worktree build` and `doctor` rehearsal before the landing of every spec that changes Go build inputs.
  Feeds: none
- After a refreeze in a multi-ticket chunk, resume every author of that chunk to record again, and plan that time in the repair cycle.
  Feeds: none