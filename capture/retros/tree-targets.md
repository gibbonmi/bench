## Outcome

Spec A of FT341 landed at `27e67a812dd2e5237d15ed3c825c752f93b4e999` from the source tip `59e05dc1` on the base `f981cd3d`. The build had five tickets in four chunks. Each chunk took three review axes and a green checkpoint.

- Ticket 1 declares one scope on each public command.
- Ticket 2 derives the kit source path through `canonicalpath` for a relative kit spelling.
- Ticket 3 leads each tree-scoped response with the `tree[1]{target,head,dirty}` row.
- Ticket 4 runs a tree-scoped verb in the tree target that `--in <label|primary>` names.
- Ticket 5 runs a kit worktree target on its own current build, and it refuses a missing or stale build.

The first landing attempt refused after every gate phase passed. The kit test run could not remove its temporary directory, because a crash system test left a subtree there. The retry landed with a green gate, and both landing effects completed.

## Gate-stage timings

- landing: commit 27e67a812dd2e5237d15ed3c825c752f93b4e999, trace 7ec202c208baba0cc28bda73246e8c7d
- gofmt: 156 ms
- vet: 1368 ms
- test: 155710 ms
- race: 3164 ms
- system: 48317 ms
- shellcheck: 597 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| TT-C4 R1 | one helper builds the ambiguity id list | verified | 9 | true | opus / xhigh / repair |
| TT-C4 R2 | `internal/env` names the wrapper variable once | verified | 9 | true | opus / xhigh / repair |
| TT-C4 R3 | the missing-tree and bundle refusals are pinned at `--in` | verified | 9 | true | opus / xhigh / repair |
| TT-C4 R4 | the build refusal prints through the shared printer | verified | 9 | true | opus / xhigh / repair |
| TT-C4 R5 | `kittest` writes its manifest through `internal/freshness` | verified | 9 | true | opus / xhigh / repair |
| TT-C4 R6 | the TT48 test pins the child home | verified | 9 | true | opus / xhigh / repair |
| TT-C4 R7 | the shared printer escapes only a line that fails `LineSafe` | verified | 9 | true | opus / xhigh / repair |
| TT-C4 R8 | one function fills the build grammar | verified | 9 | true | opus / xhigh / repair |
| C4-S1 | the ambiguity id list has two sources | finding | 7 | true | fable / high / reviewer |
| C4-S2 | the build refusal restates the shared printer | finding | 6 | true | fable / high / reviewer |
| C4-S3 | `WrapperEnv` restates `env.WrapperRouting` | finding | 6 | true | fable / high / reviewer |
| C4-S4 | `kittest` restates the manifest grammar | finding | 5 | true | fable / high / reviewer |
| C4-P1 | the spec names the wrong quoting function | finding | 6 | true | fable / high / reviewer |
| C4-P2 | the TT34 row text is a paraphrase | no-op | 6 | unknown | fable / high / reviewer |
| C4-C1 | no test pins the missing-tree or bundle refusal | finding | 8 | true | fable / high / reviewer |
| C4-C2 | no test pins the child home | finding | 7 | true | fable / high / reviewer |
| C4-C3 | a terminal read can hang a tree child | finding | 5 | true | fable / high / reviewer |
| C4-S5 | the build grammar is filled at two sites | finding | 4 | true | fable / high / reviewer |
| C4-P3 | ticket 5 names the wrong quoting function | finding | 8 | true | fable / high / reviewer |
| C4-P4 | the fence disposition contradicts the ownership list | finding | 7 | true | fable / high / reviewer |
| C4-P5 | a backslash label prints a wrong repair command | finding | 4 | true | fable / high / reviewer |

The Brier mean is 0.107 over 20 labeled pairs, with 1 abstention. C4-P2 abstains, because its veto flag is still open. This table holds only the claims of the 2026-09-30 session. The review record holds the claims of the earlier sessions: TT-C1 to TT-C3 and the first ticket 4 and ticket 5 authors. This retro does not label them.

## Coordinator catches

- The handoff was behind the tree. Ticket 5 and the TT-C4 round 1 review had committed after the last handoff write, so the session resumed from the tree.
- Two fence expansions went red at the preflight: `registry-closure` for `internal/worktree` and `anchor-closure` for `internal/usage`. Each closure fix went into its plan commit before dispatch.
- The ticket 5 repair moved the label quoting into `internal/worktree`. So the import-edge sentence that the C4-P1 fold had added became false, and a plan commit corrected it.
- The first TT-C4 checkpoint refused with no retained verification, because the repair charges told each session not to edit the record JSON. Each repair session then ran its ticket checks again at the final source and retained them.
- The orchestrator reproduced the plan and source digests from recorded values before each record update.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-declare-command-scope.md | 1 | one-source |
| 2-derive-kit-source-path.md | 1 | other |
| 3-name-the-graded-tree.md | 2 | spec-row, spec-row |
| 4-run-verbs-in-tree-target.md | 1 | one-source |
| 5-run-kit-worktree-build.md | 2 | one-source, one-source |

## Agent-experience improvements

### Bench CLI

- The landing census counted 15 raw calls, and the learning "tt-integration landing census: 15 raw calls" holds the count for each verb head.
  Feeds: new
- Add a `bench review-record` verb that appends a review or verification entry and pins the plan and source digests.
  Feeds: new
- Let `bench probe` take a repeated `--package`, so one mutation grades two readers in one private build.
  Feeds: new
- Make the inherited-red landing refusal name a baseline that can pass when the primary checkout holds git-ignored local files.
  Feeds: new

### Skills

- Make the repair charge in `bench-implement-spec` require the retained verification entries at the final source before the checkpoint.
  Feeds: new
- Make a plan commit that renames a spec symbol search the tickets for the old symbol.
  Feeds: new
- Make a fence expansion read the fence disposition text and correct each sentence that names the added package.
  Feeds: new

### Process

- Make `TestOtelCrashKeepsStartedPhaseLine` wait for its crashed child and remove its own home, so the kit test run can remove its temporary directory.
  Feeds: new