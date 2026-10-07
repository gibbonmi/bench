# Grade the published tree at the complete checkpoint

Status: staged

Roadmap: FT392

Decision source: named reviewed artifact `roadmap/FT392.md`, opened by drain d-6259cc0d8421 and placed first in quality-1 by the reviewer on 2026-10-07

Audience: every repository that links the kit

Verification log: 2 iteration(s) to accept — iteration 1 rejected on the fixture git-dir, the witness readers, and the unproven route seam. Iteration 2 accepted with three minor folds.

## Problem

`bench gate --checkpoint <spec> --complete` grades the checkout tree of the reviewed source.
In that tree, the spec is still staged and the delivery closure has not run.
The landing grades a different tree.
The landing flips the spec status, closes the roadmap rows, and deletes their detail files before its gate runs.
A red that the closure creates therefore shows first at the landing gate.
On 2026-10-07 two decision maps cited `roadmap/FT290.md`, which the closure deletes, and only the landing gate went red.

## Solution

The complete checkpoint grades the tree that the landing publishes.
It composes that tree from the committed source tip through the landing's own transform, `published.Tree`.
It grades the composed tree through the landing's own prospective execution, with the same completion obligation.
A red that the closure creates then refuses the checkpoint, with the project gate's own diagnostic.
A checkout with uncommitted changes refuses before the oracle runs, because the landing publishes only committed bytes.
The chunk checkpoint and the landing do not change.

## User stories

Line: `opus` / high.
Implementation-line reason: the hardest chunk changes oracle code on the CLI path and moves shared fixture witnesses. A probe resolved the seams, the spec is precise, and the repro is red-capable.
Harder chunks: CC1.

### The complete checkpoint grades the published tree

1. As an orchestrator, I want the complete checkpoint to grade the tree that the landing publishes, so that a closure red shows before the landing.
2. As an orchestrator, I want the checkpoint to compose that tree through the landing's own transform, so that the two grades cannot drift apart.
3. As an orchestrator, I want the graded tree to carry the implemented spec status, so that a status check grades the published status.
4. As an orchestrator, I want the graded tree to omit each detail file that the closure deletes, so that a dangling citation fails.
5. As an orchestrator, I want a red checkpoint to keep the project gate's exit and diagnostic, so that I can attribute the red without archaeology.
6. As an orchestrator, I want green evidence keyed to the published tree and the completion purpose, so that it names the landing's tree.
7. As an orchestrator, I want a repeated checkpoint on an unchanged source to reuse its green, so that I avoid a second run.
8. As an orchestrator, I want `--fresh` to force a new run on the published tree, so that I keep my escape from a stale green.
9. As an orchestrator, I want the checkpoint to leave my checkout unchanged, so that the grade does not move the source that I land.

### The checkpoint composes from the committed source

10. As an orchestrator, I want the checkpoint to compose from the committed source tip, so that it grades the bytes that the landing consumes.
11. As an orchestrator, I want an uncommitted tracked change to refuse the checkpoint before the oracle runs, so that a green never covers unpublished bytes.
12. As an orchestrator, I want an untracked file to refuse the checkpoint before the oracle runs, so that a stray file cannot enter the grade.
13. As an orchestrator, I want an uncommitted review record to refuse the checkpoint, so that the landing reads the graded evidence.
14. As an orchestrator, I want the clean-checkout refusal to tell me to commit or remove the change, so that I can recover.
15. As an orchestrator, I want the clean-checkout refusal to come before the evidence refusals, so that I repair the cause that blocks every later check.

### Existing refusals and inputs keep their behavior

16. As an orchestrator, I want stale completion evidence to refuse before the oracle runs, so that the transform keeps the evidence checks.
17. As an orchestrator, I want a spec with no staged status line to refuse before the oracle runs, so that the landing's refusal comes first.
18. As an orchestrator, I want a spec path with spaces and glob characters to compose and pass, so that the literal path rule survives.
19. As a linked-repository maintainer, I want a tree with no commitment policy to grade the status transform alone, so that my checkpoint works.
20. As an orchestrator, I want a comment-only gap past the last chunk to keep passing the complete checkpoint, so that the accepted comment rule survives.

### Unchanged scope and its proof

21. As an author, I want the chunk checkpoint to keep grading the checkout tree, so that in-progress review evidence does not change.
22. As a reviewer, I want the landing's own completion proof unchanged, so that the landing stays the authority that publishes.
23. As a maintainer, I want the reference guide to state what the complete checkpoint grades, so that a new teammate reads the decided state.
24. As a maintainer, I want the collision repro promoted into an ordinary test, so that the gate keeps the regression without a build tag.
25. As a maintainer, I want the fixture run counter to observe a prospective run, so that a no-oracle-run assertion cannot pass by construction.
26. As a maintainer, I want one closure derivation, so that the checkpoint and the landing cannot compose different trees.

## Implementation decisions

The CLI gate path owns the change.
When the parsed checkpoint selects `--complete`, the run takes the completion route.
Every other gate run keeps its current evaluation.

The completion route does four steps in this order:

1. It refuses a checkout whose working-tree hash differs from the tree of `HEAD`.
2. It resolves `HEAD` as the source tip.
3. It composes the published tree with `published.Tree` from the tip's tree, the spec path, and the tip.
4. It grades that tree through the prospective execution that `ExecuteTree` uses, with `WithCompletion` for the spec and the tip.

The route keeps the caller's run mode, so `--fresh` still forces execution.
It keeps the signal arm that the CLI path installs today.
A compose error refuses with exit 1 and its own reason, before the oracle runs.

Reviewer decision (2026-10-07): a dirty checkout refuses the complete checkpoint before the oracle runs.
The review record has no exemption from this rule.
The working-tree hash is the existing ordinary subject tree, which includes untracked files and excludes ignored files.
The refusal text contains `complete checkpoint requires a clean checkout: commit or remove each uncommitted change`.
FT393 owns the typed recovery route that a later refusal registry prints.

The gate package imports `internal/landing/published`.
`published` does not import the gate package, so the edge adds no cycle.
`go list -deps ./internal/landing/published` lists only `internal/gate/greenmarker` from the gate tree.

The existing completion proof in the prospective evaluation stays the authority over the graded tree.
That proof reads the reviewed source tree for the plan and the source digest, and it reads the record from the graded tree.
The proof requires every path outside the transform to keep its source bytes and mode.
So any landing that the gate accepts publishes exactly the tree that the checkpoint composes from the same tip.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| CC1 / 01-move-fixture-witnesses-to-the-common-directory.md, 02-grade-the-published-tree-at-the-complete-checkpoint.md | The complete checkpoint grades the published tree of the committed source, and the reference guide states it | CC01-CC32 | `internal/gate` package tests, the `cmd/bench` gate route test, and the `internal/worktree` landing journey | yes |

This is one behavior outcome with one review checkpoint.
Ticket 1 moves the fixture witnesses to the Git common directory.
CC25 and CC32 then keep the refusal assertions honest when ticket 2 changes the route.

### Completion plan

The version 2 plan records future implementation evidence.
It claims no current implementation pass, red, or probe result.
Each chunk row names the ticket that owes it.
The orchestrator adds each author session to the execution block before dispatch.

```bench-completion-plan
{
  "version": 2,
  "chunks": [
    {
      "id": "CC1",
      "tickets": [
        "01-move-fixture-witnesses-to-the-common-directory.md",
        "02-grade-the-published-tree-at-the-complete-checkpoint.md"
      ],
      "verification": [
        {
          "id": "t1-gate-package",
          "command": "bench test --package ./internal/gate",
          "ticket": "01-move-fixture-witnesses-to-the-common-directory.md"
        },
        {
          "id": "t2-gate-package",
          "command": "bench test --package ./internal/gate",
          "ticket": "02-grade-the-published-tree-at-the-complete-checkpoint.md"
        },
        {
          "id": "t2-complete-checkpoint",
          "command": "bench test --package ./internal/gate --run 'TestCompleteCheckpoint|TestChunkCheckpointGradesTheCheckoutTree|TestReviewCheckpoint|TestCommitmentExactTransform'",
          "ticket": "02-grade-the-published-tree-at-the-complete-checkpoint.md"
        },
        {
          "id": "t2-public-route",
          "command": "bench test --package ./cmd/bench --run 'TestGateCheckpointRoute'",
          "ticket": "02-grade-the-published-tree-at-the-complete-checkpoint.md"
        },
        {
          "id": "t2-landing-journey",
          "command": "bench test --package ./internal/worktree --run 'TestLandCommandPublicRealGitJourney'",
          "ticket": "02-grade-the-published-tree-at-the-complete-checkpoint.md"
        },
        {
          "id": "t2-route-omission-proof",
          "command": "bench test --package ./internal/gate --run 'TestCompleteCheckpoint'",
          "probe": "Route --complete through the ordinary checkout evaluation. TestCompleteCheckpointGradesTheClosedTree must fail, then pass after source restoration.",
          "ticket": "02-grade-the-published-tree-at-the-complete-checkpoint.md"
        },
        {
          "id": "t2-clean-checkout-proof",
          "command": "bench test --package ./internal/gate --run 'TestCompleteCheckpointRefusesADirtyCheckout'",
          "probe": "Remove the clean-checkout refusal, or exempt the review record from it. Each named refusal assertion must fail, then pass after source restoration.",
          "ticket": "02-grade-the-published-tree-at-the-complete-checkpoint.md"
        },
        {
          "id": "t2-witness-proof",
          "command": "bench test --package ./internal/gate --run 'TestReviewCheckpointReuse|TestCompleteCheckpointEvidenceNamesThePublishedTree'",
          "probe": "Point the witness helper and the fixture witnesses at the checkout. The CC25 and CC32 assertions must fail, then pass after source restoration.",
          "ticket": "02-grade-the-published-tree-at-the-complete-checkpoint.md"
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "ordinary-integration",
      "command": "bench test --package ./..."
    },
    {
      "id": "coverage",
      "command": "bench coverage --check complete-checkpoint-closure"
    }
  ],
  "execution": {
    "mode": "delegate",
    "run_id": "ft392-build-20261007",
    "orchestrator_session": "claude:ft392-orchestrator-20261007",
    "author_limit": 1,
    "assignments": {
      "01-move-fixture-witnesses-to-the-common-directory.md": [],
      "02-grade-the-published-tree-at-the-complete-checkpoint.md": []
    }
  }
}
```

## Testing decisions

The external behavior is the exit code, the diagnostic, and the retained verdict evidence of `gate.RunCommand` with `--checkpoint <spec> --complete`.
The package tests drive `RunCommand` in process.
The public route test drives a subprocess through the shell wrapper.

Promote `TestCollisionCompletionCheckpointGradesTheClosedTree` from the `landing-collisions` worktree into `internal/gate/complete_checkpoint_test.go` with no build tag.
Its fixture is `outcomeFixture`, `recordtest.Attach`, and `commitmenttest.SeedClosure`.
Its project gate exits 9 when `decisions/x.md` exists and `roadmap/FT1.md` does not.
Its witness reads the landing's closure of the same source through `published.Tree` and grades it red.

The prospective run executes in a private linked checkout that the run removes.
The gate writes `bench-last-gate` and its verdict evidence in the Git common directory.
A linked checkout's own Git directory does not hold them.

The `outcomeFixture` script resolves `git rev-parse --path-format=absolute --git-common-dir` once.
It writes `.gate-run-count` and `.gate-record-during` there, and it copies `bench-last-gate` from there.
The `--path-format=absolute` flag requires Git 2.31 or later.
The kit declares no minimum Git version, although production Git administration code already uses this flag.
The build's `CHANGELOG.md` entry records Git 2.31 as the minimum Git version.

One test helper owns the witness path for every reader.
These readers move to it:

- `outcomeRuns` in `internal/gate/run_outcomes_test.go`
- the `.gate-record-during` read in `internal/gate/run_outcomes_test.go`
- the no-run assertion in `internal/gate/review_checkpoint_test.go`
- the two no-run assertions in `internal/gate/review_checkpoint_commits_test.go`
- the no-run and no-record assertions and the record wait in `internal/gate/run_failure_outcomes_test.go`

A reader that keeps the checkout path passes by construction on every prospective run.
CC25 and CC32 keep the helper honest.

The gate calls the run-binary owner only for a gate script that hands off to `gate-phases`.
The route fixture's `.bench/gate.sh` is `exit 0`, so the public route test calls no owner and needs no fixture change.

The reviewer's refusal decision reds the `complete` case of `TestReviewCheckpointCommentOnlyGap`.
That case saves its completion record and does not commit it, so the clean-checkout refusal stops it.
Its edit commits the record before `--complete`, and its assertions stay unchanged.
No other existing complete-checkpoint fixture leaves an uncommitted change.

Each new expectation derives from the fixture's inputs, never from the code under test.
The evidence-key row computes its expected tree through `published.Tree`, the one closure owner, because that equality is the promise.

### Seam diagram

    trigger: bench gate --checkpoint <spec> --complete
        │
        ▼
    RunCommand ──▶ [ completion route ]
                      ├── dirty checkout ──▶ refusal, exit 1, no oracle
                      ├── published.Tree(HEAD^{tree}, spec, HEAD) ──▶ compose error ──▶ refusal, exit 1
                      └── prospective execution + WithCompletion(spec, HEAD)
                              ├── completion evidence refusal ──▶ exit 1, next=bench preflight review <slug>
                              └── project gate on the published tree ──▶ exit and evidence
        ◀ tests attach at RunCommand and at InspectTreeContext for the published tree

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| CC01 | 1, 4 | The complete checkpoint on the repro source exits non-zero | planned TestCompleteCheckpointGradesTheClosedTree in internal/gate/complete_checkpoint_test.go | The current checkpoint grades the unclosed tree and exits 0 |
| CC02 | 5 | The red complete checkpoint on the repro source exits 9 | planned TestCompleteCheckpointGradesTheClosedTree in internal/gate/complete_checkpoint_test.go | A refusal that replaces the oracle verdict with exit 1 hides the gate's own exit |
| CC03 | 5 | The red complete checkpoint's output contains `decisions/x.md cites missing roadmap/FT1.md` | planned TestCompleteCheckpointGradesTheClosedTree in internal/gate/complete_checkpoint_test.go | A route that drops the phase diagnostic loses the attribution |
| CC04 | 4 | An ordinary gate run on the repro source's checkout exits 0 | planned TestCompleteCheckpointGradesTheClosedTree in internal/gate/complete_checkpoint_test.go | A fixture whose source is already red makes CC01 pass for the wrong reason |
| CC05 | 2, 6 | After a green complete checkpoint, `InspectTreeContext` under `WithCompletion` for the `published.Tree` result of `HEAD` reports a reusable green | planned TestCompleteCheckpointEvidenceNamesThePublishedTree in internal/gate/complete_checkpoint_test.go | Evidence keyed to the checkout tree or to a second closure derivation misses this key |
| CC06 | 3, 19 | A complete checkpoint passes when the project gate exits 8 on a `Status: staged` spec in a tree with no commitment policy | planned TestCompleteCheckpointGradesTheImplementedStatus in internal/gate/complete_checkpoint_test.go | A checkpoint that grades the checkout tree reads the staged status and exits 8 |
| CC07 | 7 | A second complete checkpoint on the unchanged source prints `gate: green (fresh verdict reused for this tree)` | planned TestCompleteCheckpointReuseAndFresh in internal/gate/complete_checkpoint_test.go | A route that always executes never prints the reuse line |
| CC08 | 7 | A second complete checkpoint on the unchanged source leaves the common-directory run count unchanged | planned TestCompleteCheckpointReuseAndFresh in internal/gate/complete_checkpoint_test.go | A route that prints reuse but executes increments the count |
| CC09 | 8 | A complete checkpoint with `--fresh` after a green increments the common-directory run count by one | planned TestCompleteCheckpointReuseAndFresh in internal/gate/complete_checkpoint_test.go | A route that drops the run mode reuses the green |
| CC10 | 9 | A complete checkpoint leaves `HEAD` at its prior commit | planned TestCompleteCheckpointLeavesTheCheckout in internal/gate/complete_checkpoint_test.go | A route that commits the closed tree into the checkout moves `HEAD` |
| CC11 | 9 | A complete checkpoint leaves `git status --porcelain` empty in the checkout | planned TestCompleteCheckpointLeavesTheCheckout in internal/gate/complete_checkpoint_test.go | A route that reads the closed tree into the checkout leaves a dirty status |
| CC12 | 11 | A complete checkpoint with an uncommitted tracked edit exits 1 | planned TestCompleteCheckpointRefusesADirtyCheckout in internal/gate/complete_checkpoint_test.go | The current ordinary subject grades the working tree and passes |
| CC13 | 11 | A complete checkpoint with an uncommitted tracked edit leaves the common-directory run count absent | planned TestCompleteCheckpointRefusesADirtyCheckout in internal/gate/complete_checkpoint_test.go | A refusal after the oracle wastes the run and grades bytes that the landing never publishes |
| CC14 | 12 | A complete checkpoint with one untracked file exits 1 | planned TestCompleteCheckpointRefusesADirtyCheckout in internal/gate/complete_checkpoint_test.go | A cleanliness test that reads only `git diff` misses an untracked file |
| CC15 | 13 | A complete checkpoint whose only change is an uncommitted review record prints `complete checkpoint requires a clean checkout: commit or remove each uncommitted change` | planned TestCompleteCheckpointRefusesADirtyCheckout in internal/gate/complete_checkpoint_test.go | A record-path exemption reaches the evidence refusal, which prints a different reason |
| CC16 | 14 | A complete checkpoint with an uncommitted edit to `tracked.txt` prints `complete checkpoint requires a clean checkout: commit or remove each uncommitted change` | planned TestCompleteCheckpointRefusesADirtyCheckout in internal/gate/complete_checkpoint_test.go | A generic subject refusal names no recovery |
| CC17 | 15 | A dirty checkout with stale completion evidence prints the dirty-checkout refusal instead of `completion is incomplete or stale` | planned TestCompleteCheckpointRefusesADirtyCheckout in internal/gate/complete_checkpoint_test.go | A route that grades evidence first names the wrong cause |
| CC18 | 16 | A complete checkpoint with stale completion evidence prints `completion is incomplete or stale` | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | A route that skips the completion proof loses the evidence refusal |
| CC19 | 16 | A complete checkpoint with stale completion evidence leaves the common-directory run count absent | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | A route that checks evidence after the oracle runs the oracle first |
| CC20 | 17 | A complete checkpoint on a spec with no `Status: staged` line prints `spec has no Status: staged line` | planned TestCompleteCheckpointRefusesAnUntransformableSpec in internal/gate/complete_checkpoint_test.go | A route that falls back to the checkout tree passes a spec that the landing refuses |
| CC21 | 17 | A complete checkpoint on a spec with no `Status: staged` line leaves the common-directory run count absent | planned TestCompleteCheckpointRefusesAnUntransformableSpec in internal/gate/complete_checkpoint_test.go | A route that composes after the oracle runs the oracle first |
| CC22 | 18 | The public complete checkpoint on `specs/example [*] space/spec.md` exits 0 | `cmd/bench/gate_route_test.go` (`TestGateCheckpointRoute`) | A pathspec that is not literal mangles the space and glob characters |
| CC23 | 20 | A comment-only gap past the last chunk passes the complete checkpoint after one oracle run, with its completion record committed | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointCommentOnlyGap`) | A route that drops the comment-gap proof refuses the accepted gap |
| CC24 | 21 | A chunk checkpoint on the repro source exits 0 | planned TestChunkCheckpointGradesTheCheckoutTree in internal/gate/complete_checkpoint_test.go | A route that composes the closure for every checkpoint changes the chunk grade |
| CC25 | 25 | The complete step of the reuse test counts its third oracle run | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointReuse`) | A counter left in the private checkout stays at two runs, and a checkout-path stat passes by construction in `TestReviewCheckpoint`, `TestReviewCheckpointCommentOnlyGap`, `TestReviewCheckpointKeepsStrictEvidence`, `TestGateRunRefusesInitialPendingPersistenceWhenCacheIsDirectory`, and `TestGateRunRefusesOwnerPersistenceWhenOwnerPathIsDirectory` |
| CC26 | 22 | The prospective completion proof accepts the exact transform of a closure fixture | `internal/gate/commitment_completion_test.go` (`TestCommitmentExactTransform`) | A change to the completion proof breaks the landing's authority |
| CC27 | 22 | The public landing journey grades its published tree green and lands | `internal/worktree/land_journey_test.go` (`TestLandCommandPublicRealGitJourney`) | A change to the shared prospective execution breaks the landing |
| CC28 | 23 | The reference guide states that `--complete` grades the published tree of the committed source | review-owned | No gate check reads this prose |
| CC29 | 24 | `internal/gate/complete_checkpoint_test.go` carries no build constraint | review-owned | A retained `collisionrepro` tag removes the regression from the gate |
| CC30 | 26 | The completion route composes through `published.Tree` and through no second closure derivation | review-owned | A second derivation can agree on the fixtures and drift later |
| CC31 | 10 | A complete checkpoint after a new commit on the source grades that commit's published tree | planned TestCompleteCheckpointEvidenceNamesThePublishedTree in internal/gate/complete_checkpoint_test.go | A route that caches the first tip grades a source that the landing no longer consumes |
| CC32 | 25 | After a green complete checkpoint, `.gate-record-during` exists in the Git common directory | planned TestCompleteCheckpointEvidenceNamesThePublishedTree in internal/gate/complete_checkpoint_test.go | A fixture that copies the record from the linked checkout's private Git directory exits 1 before the project gate runs |

### Edge inventory

The hostile-input checklist applies to the spec path.
CC22 keeps a spec path with a space and glob characters on the public route.
`implementedSpec` passes the spec path to `ls-tree` with no literal pathspec magic.
CC22 proves that the space and glob path composes and passes; the author observed it pass on the probe route.
The composed tree reaches no TOON cell or patch header.

The inventory walks the absent-versus-empty pair for the commitment policy.
An absent policy gives no closure edits, and CC06 grades that case.
A policy with no delivery fact for the spec gives no closure edits through the existing `Closure` rule.

The fail-closed refusals are the dirty checkout, the compose error, and the completion evidence.
CC12-CC21 give each refusal one red-capable row and one no-oracle row where the order matters.
No test swaps a package variable across a subprocess.
The route test needs no swap, because its gate script never selects a run binary.

- **Won't handle**: a moved `main` — the landing's completion proof refuses any destination change outside the source, and FT393 with FT342 owns that recovery route.
- **Won't handle**: verdict reuse between the checkpoint and the landing — the oracle frames the identity root, and the landing grades from the primary checkout.
- **Won't handle**: the baseline runner of the landing — the checkpoint runs under its own checkout's kit, and the landing's gate stays the publishing authority.
- **Won't handle**: a tickets-only folder — the checkpoint grammar requires a spec record path, and the landing still grades its closed folder.
- **Won't handle**: the chunk checkpoint — it grades in-progress work that has no delivery, and CC24 keeps it on the checkout tree.
- **Won't handle**: the typed recovery route of the new refusal — FT393 owns the refusal-route registry.

## Ownership fences

The prospective build owns these exact paths:

- `internal/gate/gate.go`
- `internal/gate/checkpoint.go`
- `internal/gate/engine.go`
- `internal/gate/complete_checkpoint.go`
- `internal/gate/complete_checkpoint_test.go`
- `internal/gate/run_outcomes_test.go`
- `internal/gate/review_checkpoint_test.go`
- `internal/gate/review_checkpoint_commits_test.go`
- `internal/gate/prospective_owner_test.go`
- `internal/gate/run_failure_outcomes_test.go`
- `.bench/BENCH-reference.md`
- `CHANGELOG.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_retained_workflow.go`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/agents-handoff-section-rule`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`
- `reviews/complete-checkpoint-closure.md`

The ticket slice closed this list: it equals the union of the ticket write lines and the review pickup.
The coordinator owns the review pickup.
The planning author owns this spec folder during the spec phase.

## Out of scope

- The typed recovery route for checkpoint refusals: FT393, 9 edits, 3 gate runs.
- A checkpoint that composes onto the moved destination: FT342 route decision first, then 6 edits, 2 gate runs.

These are planning estimates, not approved priorities.

## Further notes

Evidence status: source-backed proposal with one throwaway probe.
The probe routed `--complete` through `published.Tree` and the prospective execution in the phase worktree.
The author removed the probe, and the phase commits no code.

### Flagged additions

- The clean-checkout refusal, CC12-CC17. The reviewer closed this fork on 2026-10-07: a dirty checkout refuses before the oracle runs, with no review-record exception.
- The fixture witness move, CC25 and CC32. It adds no behavior; it keeps the existing no-oracle assertions red-capable on the prospective route.

### Posture change

The clean-checkout refusal reds the `complete` case of `TestReviewCheckpointCommentOnlyGap`.
Its edit commits the completion record before `--complete`.
The witness move changes the path of every `.gate-run-count` and `.gate-record-during` reader, and no assertion changes its meaning.

### Source-sentence-to-row table

| FT392 sentence | rows |
|---|---|
| The complete checkpoint grades the tree that the landing will publish. | CC01, CC05, CC06, CC31 |
| It composes that tree through the same closure transform that the landing uses. | CC05, CC30 |
| The two grades cannot disagree. | CC02, CC03, CC26, CC27 |
| The repro is red today. | CC01, CC04 |
| The commitment criterion `FT392.closure`: a red that the closure creates refuses the checkpoint. | CC01, CC02 |

### Pre-review proof checklist

- `Cited symbols`: `gate.RunCommand`, `gate.ExecuteTree`, `gate.WithCompletion`, `gate.InspectTreeContext`, `published.Tree`, `reviewrecord.CheckTrees`, `commitrepo.Store.Closure`, `spec.Implemented`, and `outcomeRuns` resolve in the tree at 23c573a8.
- `Import edges`: gate to `internal/landing/published` is new; `go list -deps` shows no gate import in `published`.
- `Source-row clauses and occurrences`: one occurrence, the 2026-10-07 FT290 landing.
- `Promised field labels`: none.
- `Changed-function callers`: `RunCommand` has `cmd/bench` `gate-run` and `Command`; `executeTreeWithOwner` has `ExecuteTree` and four call sites in `prospective_owner_test.go`.
- `Copy survival`: CC30, review-owned.
- `Rendered-shape readers`: the new refusal text has no prior reader.
- `Pin operators`: CC02 compares the exit with `==`; CC03, CC15, CC16, CC18, and CC20 use substring containment.
- `Entry reads`: the route reads `HEAD` and the working tree from the root that `RunCommand` receives.
- `Derived expectations`: CC05 and CC31 derive the tree through `published.Tree`; the other expectations come from fixture inputs.
- `Consolidated rules`: none.
- `Quantified obligations`: every complete checkpoint takes the route; CC24 shows that no chunk checkpoint does.
- `Workflow-step writes`: none; the route writes only verdict evidence, keyed as the landing keys it.

### Reader inventory

- `cmd/bench/main.go:147-149`: the `gate` and `gate-run` rows call `Command` and `RunCommand`.
- `cmd/bench/help_inventory_test.go:122`: pins the `gate` help row; the grammar does not change.
- `.agents/commands/bench-final-check.md:19-26`: tells the orchestrator to obtain the complete checkpoint on the source; it stays true.
- `.bench/BENCH-reference.md:397-399`: names `--complete`; CC28 adds the graded-tree sentence there as the one owner.
- `internal/preflight/review.go:293-326`: reads the completion record at the source tip; it does not read gate evidence.
- `internal/landing/landing.go:238-251`: the landing composition and completion obligation; unchanged.
- No anchor needle pins the complete-checkpoint prose; `bench anchors` on both guidance files shows none near those lines.

### Primary sources

- `internal/gate/checkpoint.go:111-158`: the ordinary checkpoint grades the generation tree and frames the purpose and tip.
- `internal/gate/completion.go:23-73`: the completion obligation and the exact-transform proof.
- `internal/gate/engine.go:43-65`: the prospective execution owner.
- `internal/gate/evaluation.go:34-62`: the ordinary and prospective evaluations.
- `internal/landing/published/published.go:25-38`: the one closure transform.
- `internal/landing/landing.go:238-251`: the landing composes, transforms, and grades.
- `internal/commitment/repository/closure.go:58-94`: the closure reads the policy from the tree and the record from the source commit.
- `internal/gate/run_outcomes_test.go:225-250`: the fixture counter and record copy in the graded checkout.
- `internal/gate/run_transaction.go:141-144`: the gate selects a run binary only for a gate script that hands off to `gate-phases`.
- `internal/gate/run_transaction.go:263-273`: verdict evidence lives in the Git common directory.
- `internal/gate/prospectiveartifact/prospectiveartifact.go:102-110`: the private linked checkout.
- The repro `internal/gate/collision_repro_test.go` is untracked in the `landing-collisions` worktree.

Cheap-tier research delegates read these sources and reported citations.
The author re-read each definition above before the rows locked.
