# Grade a chained commit before it publishes

Status: staged

Roadmap: FT380

Decision source: named reviewed artifact `roadmap/FT380.md` (the drained FT380 row) and its approved commitment criterion `FT380-residual.before-publication`, with reviewer decision D1 closed in conversation on 2026-11-02: drop the `check-floor` clause, because build-cost retired that row.

Verification log: 0 iteration(s) to accept — pending review.

## Problem

An orchestrator commits a plan commit that no dispatch follows with `bench commit --preflight-build <slug>` (`.agents/commands/bench-implement-spec.md:62`).
That chain runs the build preflight only after the ref moves (`cmd/bench/commit_chain.go:44-63`).
A red preflight then exits 4, and the red commit is already on the integration source (`cmd/bench/commit_chain.go:26-30`, `:67-69`).
The roadmap row records several plan commits that published a red plan and needed a second commit to repair it (`roadmap/FT380.md:22-36`).

The same chain grades each ticket record commit, because `bench record ticket` prints a record commit with `--preflight-build <slug> --base <base>` (`internal/reviewrecord/recordcmd/ticket.go:101-102`).

A pre-publication grade cannot simply run the preflight before the commit verb. The build preflight reads the spec and the tickets from the checkout (`internal/preflight/gather.go:87-97`).
It reads the completion plan, the decision maps, and the review record from the source tip commit (`internal/preflight/plan.go:16-29`, `internal/preflight/map_sources.go:15-28`, `internal/preflight/gather_inputs.go:105-106`).
Before publication, that source tip is the old HEAD (`internal/preflight/gather.go:77`, `:189-195`). So a run before the commit grades the old completion plan, not the candidate's plan.

Two faces of the row are already closed:

- At the kit root, the Markdown lane class runs `docs-currency-workflow` on a spec or ticket commit, and that check runs `bench coverage --check` on every spec before publication (`internal/gate/lane_select.go:147-151`, `internal/conformance/docs_workflow_checks_test.go:333-358`).
- Build-cost ticket 01 retired the `check-floor` row (commit `f0d811335`, `CHANGELOG.md:67`), and a test keeps it out (`internal/preflight/charge_cells_test.go:51-58`).

## Solution

`bench commit --preflight-build <slug> [--base <commit>]` runs its chain steps before publication.
The landing writes the candidate commit object, and then it runs the chain check on that candidate. The check runs the worktree build and then the build preflight at the candidate commit.
If a step is red, the landing moves no ref, and the commit exits 1 with nothing committed. If every step is green, the landing moves the ref as it does now.

The build preflight gains an in-process grade at an unpublished candidate commit. That grade reads every source-tip input at the candidate and prints the same verdict text that the published grade prints for the same commit.
The build preflight runs the coverage parser as its bootstrap step (`internal/preflight/gather.go:54-63`), so one grade covers the coverage parser and the build preflight in every repository.

A plan commit that a dispatch follows still uses plain `bench commit`, because the charge preparation is the one build preflight of that dispatch (`.agents/commands/bench-implement-spec.md:21`, `:62`). This spec keeps that route and pins it with an anchor.

## User stories

Line: opus / high.
Implementation-line reason: C1 substitutes the candidate for HEAD across several gather reads, and a missed read is silent unless a differential row varies it, so C1 runs at high. C2 composes existing seams with exact stub rows and one system row, so it runs at medium.
Harder chunks: C1.

### The candidate grade

1. As an orchestrator, I want the build preflight to grade an unpublished candidate commit, so that a chained commit is graded before its ref moves.
2. As an orchestrator, I want the candidate grade to read the completion plan at the candidate, so that a broken plan fence in the plan commit goes red.
3. As an orchestrator, I want the candidate grade to read the decision maps at the candidate, so that a map source that the plan commit deletes goes red.
4. As an orchestrator, I want the candidate verdict text to equal the published verdict text for the same commit, so that one verdict shape serves both routes.
5. As an orchestrator, I want the candidate grade to run the coverage parser as its bootstrap step, so that an invalid coverage map refuses before publication.
6. As an orchestrator, I want the candidate grade to take an explicit base, so that a record commit inside a chunk grades the range of its chunk.
7. As an orchestrator, I want a candidate that names no commit refused, so that a bad candidate never reads as green.
8. As an orchestrator, I want a base that is not an ancestor of the candidate refused, so that a wrong range never grades.
9. As a CLI user, I want `bench preflight build <slug> --source-tip <commit>` to keep its HEAD comparison, so that a manual pin still reds a mismatch.

### The check before publication

10. As a commit caller, I want the landing to run a supplied check after it writes the commit object and before it moves the ref, so that the check grades the exact candidate.
11. As a commit caller, I want the candidate tree to equal the composed tree and its first parent to equal the expected base, so that the check grades what would publish.
12. As a commit caller, I want a red check to move no ref and to leave the checkout unchanged, so that nothing publishes.
13. As a commit caller, I want a red check to exit 1 with a `commit-chain-red` refusal that reruns my exact commit call, so that I repair and retry with one paste.
14. As a commit caller, I want a lane red to stop the commit before the check runs, so that the cheaper refusal wins.
15. As a commit caller, I want a commit with no check to behave as it does now, so that plain `bench commit` keeps its contract.
16. As a commit caller, I want a destination ref that moves during the check to refuse at the compare-and-swap, so that a stale candidate never publishes.

### The chain

17. As an orchestrator, I want `bench commit --preflight-build <slug>` to run the worktree build and then the build preflight before publication, so that a red chained commit never publishes.
18. As an orchestrator, I want the chain to give the candidate commit to the build preflight, so that the preflight grades the new commit and not the old tip.
19. As an orchestrator, I want the chain to give `--base` to the candidate grade, so that a record commit keeps the chunk range.
20. As an orchestrator, I want a red build to stop the chain before the preflight and to publish nothing, so that a stale build never grades.
21. As an orchestrator, I want a red preflight to publish nothing and to exit 1, so that exit 1 keeps its meaning of nothing committed.
22. As an orchestrator, I want a lane refusal to print a chain line with every step skipped, so that I see that the chain never started.
23. As an orchestrator, I want a green chain to publish and to name the commit in the chain line, so that the record names the published commit.
24. As an orchestrator, I want a publication that does not reconcile to keep exit 3, so that the remainder route stays.
25. As an orchestrator, I want exit 4 retired, so that no exit means published and red.
26. As an agent, I want the commit help to state that the chain steps run before publication, so that I learn the order without a trial.
27. As a ticket author, I want the record commit that `bench record ticket` prints to keep the chain flags, so that a red record commit never publishes.
28. As an operator, I want the real `bench commit` verb to refuse a red candidate and publish a green one, so that the composed route works end to end.

### The plan-commit route

29. As an orchestrator, I want the implement guidance to keep routing a plan commit that a dispatch follows to plain `bench commit`, so that the charge preparation stays the one preflight of that dispatch.
30. As an orchestrator, I want the implement guidance to keep routing a plan commit that no dispatch follows to the chain, so that such a plan commit is graded before it publishes.
31. As a teammate, I want the changelog to state the new step order and the retired exit, so that a reader of the release sees the contract change.

### Reviewed exclusions

32. As a reviewer, I want no `check-floor` grade restored, so that the build-cost retirement stays closed.
33. As a maintainer, I want no lane change without a `bench commit` repro of a skipped parser, so that wider lane enforcement waits for evidence.
34. As a maintainer, I want the comment-only correction after a chunk freeze left to FT439, so that this spec keeps one outcome.

## Implementation decisions

**The candidate grade.** The preflight package exports one in-process build grade at a candidate commit: `BuildVerdictAt(root, slug, base, candidate string) (string, int)`.
The grade resolves the candidate as a commit and gathers the build facts with that commit as the source tip, in place of HEAD.
The source-tip reads are the default-branch ancestry of `base-current`, the changed paths, the default-tip paths, the build deletions, the decision maps, the completion plan, and the fold resolutions (`internal/preflight/gather.go:76-150`, `internal/preflight/plan.go:16-29`, `internal/preflight/map_sources.go:15-28`, `internal/preflight/gather_inputs.go:105-106`).
The grade pins the candidate as the source-tip pin, so its verdict holds the `tip-current` row, as the published grade does when the chain passes `--source-tip` (`cmd/bench/commit_chain.go:58`).
The grade renders through the same verdict renderer as `bench preflight build` (`internal/preflight/command.go:153-213`). It appends the `commitment` row under the same rule.
With an explicit base, the grade resolves the range from that base to the candidate through `diff.ResolveSourceRange(root, base, candidate)` (`internal/diff/range.go:150-162`).
The public `--source-tip` flag keeps its meaning: a pin that compares with HEAD (`internal/preflight/source_tip_test.go:75`).

**The check seam in the landing.** `landing.Request` gains one optional field, `Check func(candidate string) error`.
`Owner.Land` calls a non-nil check after `commit-tree` and before `updateRef` (`internal/landing/landing.go:184-193`). A check error returns `landing.CheckRefusal`, and Land moves no ref and runs no reconcile.
A nil check keeps every current step. `DryRun` never calls the check, because the commit grammar refuses `--dry-run` with the chain (`internal/commit/grammar.go:85-88`).
The candidate commit object stays unreferenced after a refusal.

**The commit entry.** The commit package exports `RunChecked(args, stdout, stderr, check)`, where `check` has the signature `func(root, candidate string) error`. `Run` keeps its signature and calls `RunChecked` with a nil check, so the callers of `Run` do not change.
`landingFace` maps `landing.CheckRefusal` to a new face `commit-chain-red` (`internal/commit/commit.go:209-236`). That face exits 1 through `refuse`, which already repeats the chain flags in the rerun (`internal/commit/commit.go:248-251`).
The refusal-route registry declares the face beside `commit-red` (`internal/refusalroute/faces_commit.go:26-32`). Its route is one instruction to repair each red chain step, then the caller's own commit rerun.

**The chain.** `commitChainCommand` hands a check to `RunChecked` when `--preflight-build` names a slug (`cmd/bench/commit_chain.go:42-73`).
The check runs the worktree build first. If the build is green, it runs `BuildVerdictAt` with the slug, the base, and the candidate. It prints each step's own response and records each step state.
The chain line keeps its fields and their order: `commit-chain{commit=<sha|none>,build=<state>,preflight=<state>}`.
`publishedStepRedExit` and exit 4 retire, because no chain step runs after publication. The commit's exits keep their meanings: 0, 1, 2, and 3.
The commit's own post-publication rebuild stays (`internal/commit/background.go:24-45`).
`chainSteps.Commit` takes the check, and `chainSteps.Preflight` takes the root, the slug, the base, and the candidate, so the stub tests drive each step (`cmd/bench/commit_chain_test.go:33-61`).

**The route anchors.** One new anchor file requires the two halves of the plan-commit route sentence in `.agents/commands/bench-implement-spec.md:62`. No anchor pins that sentence now. The guidance text does not change.

## Implementation chunks

C1 adds the candidate grade and is its own chunk, because its seam is the uncertain one. Its review closes before C2 composes it.
C2 composes the check seam and the chain. Ticket 02 adds the seam with a stub check, and ticket 03 binds the real chain.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| C1 / 01-grade-the-candidate-commit.md | `BuildVerdictAt` grades an unpublished candidate and prints the published verdict text for the same commit | PP01, PP02, PP03, PP04, PP05, PP06, PP07, PP08, PP09, PP32 | preflight candidate tables, the published-verdict differential, existing source-tip tests | yes |
| C2 / 02-check-before-publication.md, 03-chain-before-publication.md | `bench commit --preflight-build` grades the candidate before publication and publishes nothing on red | PP10, PP11, PP12, PP13, PP14, PP15, PP16, PP17, PP18, PP19, PP20, PP21, PP22, PP23, PP24, PP25, PP26, PP27, PP28, PP29, PP30, PP31 | landing check tests, commit check tests, chain stub tests, chain system test, plan-commit anchors | no |

### Completion plan

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "C1",
      "tickets": ["01-grade-the-candidate-commit.md"],
      "verification": [
        {"id": "candidate-tests", "command": "bench test --package ./internal/preflight", "probe": "Read the completion plan at HEAD in the candidate grade and record the PP02 red, the restore, and the green."}
      ]
    },
    {
      "id": "C2",
      "tickets": ["02-check-before-publication.md", "03-chain-before-publication.md"],
      "verification": [
        {"id": "check-seam", "command": "bench test --package ./internal/landing --package ./internal/commit --package ./internal/refusalroute", "probe": "Call the check after the ref update and record the PP10 red, the restore, and the green."},
        {"id": "chain-stubs", "command": "bench test --package ./cmd/bench --run 'CommitChain|CommitHelp'", "probe": "Give the published HEAD in place of the candidate to the preflight step and record the PP18 red, the restore, and the green."},
        {"id": "chain-system", "command": "bench test --check system --run TestCommitChainRefusesRedCandidate", "probe": "Drop the check from the chain commit call and record the PP28 red, the restore, and the green."},
        {"id": "route-anchors", "command": "bench test --package ./internal/anchors --run PlanCommit", "probe": "Delete the dispatch half of the plan-commit route sentence and record the PP29 red, the restore, and the green."}
      ]
    }
  ],
  "final_verification": [
    {"id": "system", "command": "bench test --check system"},
    {"id": "coverage", "command": "bench coverage --check prepublication-preflight"}
  ]
}
```

## Testing decisions

A good test drives one public entry with one exact input and observes the verdict text, the exit, and the ref state.
The source-tip tests are the prior art for a verdict over a fixture repository (`internal/preflight/source_tip_test.go:23-105`). They use the exported `preflighttest` fixtures (`internal/preflight/preflighttest/fixture.go:161-209`).
The landing tests are the prior art for a ref-state assertion after a refusal (`internal/landing/landing_test.go:259`, `TestLandRefusesEveryNonGreenAuthorizationWithoutMutation`).
The chain stub harness `runCommitChain` is the prior art for step order and the chain line (`cmd/bench/commit_chain_test.go:33-61`). C2 changes its stub commit so that it calls the check with a fixed candidate.

The differential row PP04 compares two runs of the same code over one commit: the candidate grade before `update-ref`, and `Command` with `--source-tip` after it. The fixture family varies each source-tip input between HEAD and the candidate.
The system row PP28 drives the real `bench commit` verb in a linked-project worktree, after the precedent of `internal/systemtest/otel_verbs_test.go:40`.

### Seam diagram

    trigger: bench commit --preflight-build <slug> [--base <commit>] -m <msg> -- <path>...
        │
        ▼
    named paths ──▶ [ compose ──▶ lane ──▶ commit-tree ──▶ check(candidate) ──▶ update-ref ]
                                                  │
                                                  ▼
                              [ worktree build ──▶ preflight.BuildVerdictAt(candidate) ]
                                                  │
                     red: no ref moves, exit 1, commit-chain{commit=none,...}
                     green: the ref moves, exit 0, commit-chain{commit=<sha>,...}
            ◀ tests attach here: BuildVerdictAt tables and the differential; Land with a stub check;
              RunChecked with a stub check; runCommitChain stubs; one system test through the real verb

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PP01 | 1 | `BuildVerdictAt` on a green candidate that HEAD does not hold exits 0 with a green `tip-current` row | planned TestBuildVerdictAtGradesTheCandidate in internal/preflight/candidate_test.go | A grade that pins the candidate and keeps HEAD as the source tip reds `tip-current` |
| PP02 | 2 | A candidate with an unreadable completion-plan fence over a HEAD with a valid fence exits 1 with `completion-plan` red | planned TestBuildVerdictAtGradesTheCandidate in internal/preflight/candidate_test.go | A grade that reads the plan at HEAD stays green |
| PP03 | 3 | A candidate that deletes a path that a decision-map source cites exits 1 with `map-sources` red | planned TestBuildVerdictAtGradesTheCandidate in internal/preflight/candidate_test.go | A grade that lists the maps at HEAD still finds the path |
| PP04 | 4 | For the PP01, PP02, and PP03 fixtures, the candidate verdict text equals the `--source-tip` verdict text after `update-ref` moves HEAD to that candidate | planned TestBuildVerdictAtMatchesThePublishedVerdict in internal/preflight/candidate_test.go | A grade that keeps any HEAD read differs on the fixture that varies that read |
| PP05 | 5 | A candidate whose coverage map holds a two-predicate row exits 1 with the `coverage map invalid` refusal | planned TestBuildVerdictAtRefusesBadInput in internal/preflight/candidate_test.go | A grade that skips the gather bootstrap prints a verdict |
| PP06 | 6 | With an explicit base, the verdict prints a `source` table whose tip is the candidate | planned TestBuildVerdictAtGradesTheCandidate in internal/preflight/candidate_test.go | A grade that resolves the range at HEAD prints the old tip |
| PP07 | 7 | A candidate that names no commit exits 1 with a refusal that names the candidate | planned TestBuildVerdictAtRefusesBadInput in internal/preflight/candidate_test.go | A grade that ignores a failed resolution falls back to HEAD and exits 0 |
| PP08 | 8 | A base that is not an ancestor of the candidate exits 1 with the `--base is not an ancestor` refusal | planned TestBuildVerdictAtRefusesBadInput in internal/preflight/candidate_test.go | A grade that skips the ancestry test grades an unrelated range |
| PP09 | 9 | `bench preflight build <slug> --source-tip <other>` still prints `tip-current` red | `internal/preflight/source_tip_test.go` (`TestSourceTipMismatchRendersRedRow`) | A change that makes the public pin the source tip greens the mismatch |
| PP10 | 10 | Inside the check, the destination ref still names the expected base and the candidate object exists | planned TestLandRunsTheCheckBeforePublication in internal/landing/check_test.go | A check call after `updateRef` sees the moved ref |
| PP11 | 11 | The candidate tree equals the composed tree, and its first parent equals the expected base | planned TestLandRunsTheCheckBeforePublication in internal/landing/check_test.go | A check that receives HEAD or the base reads the wrong tree |
| PP12 | 12 | A check error returns `landing.CheckRefusal`, and the destination ref and the checkout fingerprint stay unchanged | planned TestLandCheckRefusalPublishesNothing in internal/landing/check_test.go | A Land that logs the error and publishes moves the ref |
| PP13 | 13 | `RunChecked` with a red check exits 1, prints `next=` with the caller's `--preflight-build` rerun, and leaves HEAD unchanged | planned TestRunCheckedRefusesARedCheck in internal/commit/check_test.go | A commit that maps the refusal to the handback face prints the reviewer route |
| PP14 | 14 | A lane red with a check prints the lane refusal, and the check never runs | planned TestRunCheckedLaneRedWinsOverTheCheck in internal/commit/check_test.go | A check that runs before the lane runs on a tree the lane refuses |
| PP15 | 15 | Land with no check composes only the named paths and moves the ref by compare-and-swap | `internal/landing/landing_test.go` (`TestLandComposesOnlyNamedPathsAndCASesExpectedBase`) | A nil-check path that changes the landing breaks this existing test |
| PP16 | 16 | A check that moves the destination ref leaves the winner in place and returns the tip-moved refusal | planned TestLandCheckRefusalPublishesNothing in internal/landing/check_test.go | A Land that updates the ref without the expected value overwrites the winner |
| PP17 | 17 | A green chain prints the build, the preflight, and the commit lines in that order, then the chain line, at exit 0 | `cmd/bench/commit_chain_test.go` (`TestCommitChainRunsThreeSteps`) | A chain that keeps the steps after the commit prints the commit line first |
| PP18 | 18 | The preflight step receives the stub candidate, not the published commit | `cmd/bench/commit_chain_test.go` (`TestCommitChainRunsThreeSteps`) | A chain that passes the published commit or HEAD fails the exact argument comparison |
| PP19 | 19 | A chain that names a base gives that base to the preflight step | `cmd/bench/commit_chain_test.go` (`TestCommitChainPassesTheBase`) | A chain that drops the base fails the exact argument comparison |
| PP20 | 20 | A red build exits 1, calls no preflight, and prints `commit-chain{commit=none,build=red,preflight=skipped}` | `cmd/bench/commit_chain_test.go` (`TestCommitChainStopsAtRedBuild`) | A chain that grades after a red build records a preflight call |
| PP21 | 21 | A red preflight at exit 1, 2, or 3 makes the chain exit 1 and print `commit-chain{commit=none,build=green,preflight=red}` | `cmd/bench/commit_chain_test.go` (`TestCommitChainReportsRedPreflight`) | A chain that keeps exit 4 or publishes fails the exit and the line |
| PP22 | 22 | A lane refusal before the check exits 1 and prints `commit-chain{commit=none,build=skipped,preflight=skipped}` | `cmd/bench/commit_chain_test.go` (`TestCommitChainSkipsAfterRefusal`) | A chain that runs its steps outside the commit runs them after a lane red |
| PP23 | 23 | A green chain line names the published commit | `cmd/bench/commit_chain_test.go` (`TestCommitChainRunsThreeSteps`) | A chain that prints `none` after a publication fails the line comparison |
| PP24 | 24 | A publication that does not reconcile exits 3 and names the commit with both steps green | `cmd/bench/commit_chain_test.go` (`TestCommitChainStopsAtRemainder`) | A chain that maps exit 3 to another exit fails the exit comparison |
| PP25 | 25 | The commit help holds no `exit 4` text | `cmd/bench/commit_chain_test.go` (`TestCommitHelpEndsWithChainHelp`) | A help that keeps the retired exit fails the absence test |
| PP26 | 26 | The chain help line starts with the `--preflight-build` group and states that the steps run before publication | `cmd/bench/commit_chain_test.go` (`TestCommitHelpEndsWithChainHelp`) | A help that keeps the after-publication text fails the containment test |
| PP27 | 27 | The `next=` route of `bench record ticket` names `--preflight-build <slug> --base <base>` | `internal/reviewrecord/recordcmd/ticket_probe_test.go` (`TestTicketStagesTheRetainedProbe`) | A change that drops the chain flags from the record route fails the literal |
| PP28 | 28 | In a linked-project worktree, a red plan commit exits 1 with the branch ref unchanged, and its repaired retry exits 0 with the chain line naming the new branch tip | planned TestCommitChainRefusesRedCandidate in internal/systemtest/commit_chain_test.go | A real verb that drops the check publishes the red plan commit |
| PP29 | 29 | The implement guidance holds the sentence half "A plan commit that a dispatch follows uses plain `bench commit`" | planned TestPlanCommitAnchors in internal/anchors/registry_plan_commit_test.go | A guidance edit that drops the dispatch route reds the Require anchor |
| PP30 | 30 | The implement guidance holds the sentence half "a plan commit that no dispatch follows uses `bench commit --preflight-build <slug>`" | planned TestPlanCommitAnchors in internal/anchors/registry_plan_commit_test.go | A guidance edit that drops the chain route reds the Require anchor |
| PP31 | 31 | The changelog entry states that the chain steps run before publication and that exit 4 retired | review-owned in C2 | A changelog that omits the contract change fails review |
| PP32 | 32 | The build verdict holds no `check-floor` row | `internal/preflight/charge_cells_test.go` (`TestDecideHasNoCheckFloorRow`) | A grade that restores the row fails this existing test |

Not covered: story 33 — no recorded `bench commit` repro of a skipped parser exists, so the lane stays as it is.
Not covered: story 34 — FT439 owns the comment-only correction after a chunk freeze.

### Edge inventory

This kit change serves every repository that links the kit, because `bench commit` and the build preflight ship in the binary.
The check holds no lock and writes no file. It writes one unreferenced commit object, and the landing's compare-and-swap is the only publication.
The hostile-input walk covers a slug or base with spaces or glob characters through the existing `refusalroute.Arg` quoting (`internal/commit/grammar.go:55-61`), a moved destination in PP16, and a non-commit candidate in PP07.
A transaction has three failure points. Before the check runs, nothing persists but the candidate object (PP12). Inside the check, a red or a bootstrap failure moves no ref (PP12, PP05). At the terminal step, the compare-and-swap refuses a moved ref (PP16).

- **Won't handle:** a chained commit that changes a Bench build input — the build step rebuilds the working tree before the preflight, and PP20 covers a red build; the guidance routes only plan and record commits through the chain.
- **Won't handle:** an unreferenced candidate object after a refusal — Git's own garbage collection removes it, and PP12 covers the unchanged ref.
- **Won't handle:** an unnamed dirty spec or ticket file in the checkout during the check — the build preflight reads the checkout tickets today (`internal/preflight/gather.go:87-89`), and the post-publication chain read the same checkout.
- **Won't handle:** `--dry-run` with the chain — the grammar refuses the pair (`internal/commit/grammar.go:85-88`), and PP15 covers the plain landing.

## Out of scope

- A lane class that runs the build preflight on every spec commit. FT380 graduates wider lane enforcement only on a `bench commit` repro of a skipped parser. Estimate: 6 edits, 2 gate runs.
- The comment-only correction after a chunk freeze, which FT439 owns. Estimate: 5 edits, 2 gate runs.
- A `--dry-run` form of the chain. Estimate: 3 edits, 1 gate run.

## Further notes

### Flagged additions

Each item below goes beyond the literal decision source. The reviewer can veto each one.

1. The record commit takes the pre-publication grade (story 27). The chain is one mechanism, and `bench record ticket` composes it. A red record commit now publishes nothing, where it published with exit 4 before.
2. Exit 4 retires (story 25). No chain step runs after publication, so no exit means published and red.
3. The build step moves before publication with the preflight (story 20). The binary-seal row needs a fresh build, as the chain needs now.
4. The two plan-commit route anchors (PP29, PP30). The source's second sentence is a guidance fact, and no anchor pins it now.
5. The new refusal face `commit-chain-red` (story 13). The existing faces name a lane red or a handback, and neither names the chain step.

### Seams beyond the decision source

- `preflight.BuildVerdictAt`: the gather reads the source tip at `internal/preflight/gather.go:77`, `internal/preflight/plan.go:17`, `internal/preflight/map_sources.go:19`, and `internal/preflight/gather_inputs.go:106`. The coordinator's own read found these reads. No research reader found them, and every reader proposed a checkout grade before `Land`.
- The `landing.Request.Check` field: `internal/landing/landing.go:184-193`. The enforcement readers found that no step runs between `commit-tree` and `updateRef`.
- `commit.RunChecked` and the `commit-chain-red` face: `internal/commit/commit.go:82-116`, `:209-236`, and `internal/refusalroute/faces_commit.go:7-72`. The reader sweep found every caller of `commit.Run` and the face registry.

### Source-sentence-to-row table

| Source clause | Coverage |
|---|---|
| A plan commit publishes before any check runs. | PP10, PP12, PP17, PP21, PP28 |
| A check before publication runs the coverage parser and the build preflight first. | PP01, PP02, PP03, PP04, PP05, PP18 |
| The commit lane does not run the `check-floor` row, so a plan commit can publish a red floor. | D1: dropped; PP32 keeps the row retired |
| The Markdown lane runs the coverage parser through `docs-currency-workflow`. | Already true at the kit root (`internal/gate/lane_select.go:147-151`); PP05 covers every root |
| Wider lane enforcement graduates only on a `bench commit` repro of a skipped parser. | Not covered: story 33 |
| FT439 owns the comment-only correction after a chunk freeze. | Not covered: story 34 |
| Criterion: a plan commit that no dispatch follows runs the coverage parser and the build preflight before it publishes. | PP05, PP21, PP28, PP30 |
| Criterion: a plan commit that a dispatch follows takes its build preflight from the charge preparation (FT445). | PP29 |
| Occurrence 2026-10-04 DC-C8: the coverage parser refused rows only after the commit. | PP05, PP28 |
| Occurrence 2026-10-04 DC-C8 seam correction: a plan commit after the freeze forced a re-freeze. | Guidance: `.agents/commands/bench-implement-spec.md:60`; no row |
| Occurrence 2026-10-04 DC-C9: the `Writes:` change needed the fence union and the registry closure. | PP02, PP04 |
| Occurrence 2026-10-03 FT370: citations changed after the final freeze. | Guidance: `.agents/commands/bench-implement-spec.md:60`; no row |
| Occurrence 2026-10-06 FT376: the amendment left a row out of the chunk table. | PP04 (`chunk-rows` in the same verdict) |
| Occurrence 2026-10-08 FT349: a ticket sentence after the freeze staled the source. | Guidance: `.agents/commands/bench-implement-spec.md:60`; no row |
| Occurrence 2026-10-10 BD-C4 and the build-directions plan commit. | PP21 |
| Occurrence 2026-10-10 BD-C3: the record-commit base was refused. | Guidance: `.agents/commands/bench-review-implementation.md:99`; no row |
| Occurrence 2026-10-11 SB1: a prose-only plan commit staled the chunk review. | Guidance: `.agents/commands/bench-implement-spec.md:60`; no row |
| Occurrence 2026-10-12: a chunk reorder left D-C1 graded inside D-B6. | Guidance: `.agents/commands/bench-implement-spec.md:60`; no row |
| Occurrence 2026-10-13 PS-D: the expansion omitted a floor check. | D1: the row retired; PP32 |
| Occurrence 2026-10-14 FT405: a plan commit after the chunk tip stalled the checkpoint. | Guidance: `.agents/commands/bench-implement-spec.md:60`; no row |
| Occurrence 2026-10-14 FT379: the `check-floor` row went red three times. | D1: the row retired; PP32 |
| Occurrence 2026-10-18 SS-C2 and SS-C3: the completion entry refused the chain gap. | Guidance: `.agents/commands/bench-review-implementation.md:99`; no row |

### Pre-review proof checklist

- Cited symbols: `commitChainCommand`, `chainSteps`, `publishedStepRedExit`, `chainHelp`, `Run`, `commitAttributed`, `refuse`, `landingFace`, `Owner.Land`, `composeAuthorized`, `gather`, `headTip`, `withCompletionPlan`, `mapSourceDiagnostics`, `withFoldResolution`, `verdictCommand`, `Decide`, `ResolveSourceRange`, `ChainArgs`, `chainCommit`, and `refreshWorktreeBuild` were opened at their cited lines.
- Import edges: `cmd/bench` already imports `internal/commit`, `internal/preflight`, and `internal/worktree` (`cmd/bench/commit_chain.go:9-12`). `internal/commit` already imports `internal/landing`. `go list -deps ./internal/preflight` and `go list -deps ./internal/landing` name neither `internal/commit` nor each other's importer, so no cycle forms. `internal/systemtest` already imports `preflighttest`.
- Source-row clauses and occurrences: each clause and each occurrence is in the source table.
- Promised field labels: the function `BuildVerdictAt`, the field `landing.Request.Check`, the error `landing.CheckRefusal`, the function `commit.RunChecked`, the face `commit-chain-red`, and the chain line `commit-chain{commit,build,preflight}`.
- Changed-function callers: `commit.Run` has the chain and `internal/gate/refusal_route_test.go:49`, `:182`; `Run` keeps its signature. `Owner.Land` callers keep a nil check. `chainSteps` has one test harness, `runCommitChain`.
- Copy survival: none.
- Rendered-shape readers: the chain help text and the exit-4 text have two readers, `cmd/bench/commit_chain.go:33-34` and `cmd/bench/commit_chain_test.go:185`. The chain line literal has one test reader, `cmd/bench/commit_chain_test.go:66`. The `bench help` row (`cmd/bench/help_inventory_test.go:153`) reads `HelpRowSuffix`, which does not change.
- Pin operators: PP04 uses byte equality; PP17, PP20, PP21, PP22, and PP24 use exact string equality on stdout or the last line; PP25 uses absence; PP26 uses a prefix and containment; PP29 and PP30 use the anchor Require containment.
- Entry reads: `BuildVerdictAt` reads the repository at `root`, and its candidate and base are values. `RunChecked` reads the Git root from the working directory, as `Run` does.
- Derived expectations: the PP04 expectation comes from the published `Command` run, not from the candidate grade; every other expectation is an authored literal.
- Consolidated rules: none.
- Quantified obligations: "every source-tip read" holds for each read that the candidate grade lists, and PP04 varies the completion plan, the decision maps, and the green case.
- Workflow-step writes: the check writes no file. It reads HEAD as the expected base, the review record at the candidate, and the checkout spec and tickets. The landing's compare-and-swap reads the destination ref after the check.

### Research comparison

The coordinator read the enforcement surface itself and charged six read-only readers in three pairs. The agreed set matched the coordinator's reads: the post-publication chain, the retired `check-floor` row, the Markdown lane's coverage parser, and the FT445 charge preparation.
The readers found the record-commit route at `internal/reviewrecord/recordcmd/ticket.go:101`, which the coordinator had not read.
The coordinator found the source-tip reads of the build preflight, which no reader found. Two readers proposed a checkout grade before `Land`, and that grade is the cheapest wrong implementation that PP02 and PP04 refuse.
