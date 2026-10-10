# Give the explicit plan the one landedness verdict

Blocked by: none
Writes: internal/worktree/lifecyclepolicy/lifecyclepolicy.go, internal/worktree/lifecyclepolicy/lifecyclepolicy_test.go, internal/worktree/subshell.go, internal/worktree/landed.go, internal/worktree/classifier.go, internal/worktree/clean_landed.go, internal/worktree/list.go, internal/worktree/ownership.go, internal/worktree/eligibility_test.go, internal/worktree/lifecycle_facts_test.go, internal/worktree/unlanded_route_test.go, CHANGELOG.md
Covers: none

## What to build

The explicit planner gives each cleanup plan a landedness verdict from the branch ancestry alone. The automatic classifier then derives a second verdict through `assignmentLandedness`, which applies the own-commit rule of FT408. Thus `plan.landedTyped` reads `true:ancestry` for an assignment with no own commit and uncommitted output. The `bench worktree list` cell and the release read the same assignment as unlanded. A new reader of `plan.landedTyped` gets the ancestry verdict and can bring back the release-route defect.

One rule gives the landedness of an assignment. The rule returns the ancestry proof, except that an assignment that holds unlanded output has not landed. A failed output-fact query makes the landedness `unknown:<err>`. The `lifecyclepolicy` package owns this rule as one pure function, and `DecideExplicit` applies it to the facts of each explicit plan. `assignmentHoldsUnlandedOutput` stays the one source of the output fact.

The explicit planner gathers the output fact for the assignment that the plan names. That assignment is the matched owned assignment, else the foreign recovery assignment. The planner gathers the fact only when the ancestry proves the branch landed. Thus `plan.landedTyped`, the `landed` fingerprint evidence, and the automatic outcome carry the one verdict.

Each reader reads that verdict and derives no second verdict:

- The automatic classifier reads the landedness of `explicitOutcome`, and it does not call `assignmentLandedness`.
- `assignmentLanded` reads `plan.landedTyped` directly.
- The landed selector proof has no explicit plan, so it applies the policy rule to its ancestry proof and the output fact.
- The landed cell of `bench worktree list` applies the same policy rule.

This ticket deletes `assignmentLandedness`. The release, the list cell, the landed sweep, and the session-start landed count keep their current behavior.

Branch-deletion authority stays on the ancestry proof, which is the recommended option in the decision below. That authority answers a different question: does `main` hold every commit of the branch. The comments that call such a branch "landed" in `deletePlannedBranch` and `cleanupRow` then name the ancestry proof.

`bench worktree clean <path>` without `--discard-branch` gives these plans:

| Assignment shape | Plan action | `landed` evidence before | `landed` evidence after | Branch after the apply |
|---|---|---|---|---|
| No own commit, clean tree | `remove` | `true:ancestry` | `true:ancestry` | deleted, unchanged |
| No own commit, uncommitted output | `recover-remove` | `true:ancestry` | `false:ancestry` | deleted, unchanged |
| Own commit that `main` holds, clean tree | `remove` | `true:ancestry` or `true:patch` | unchanged | deleted, unchanged |
| Own commit that `main` holds, uncommitted output | `recover-remove` | `true:ancestry` or `true:patch` | unchanged | deleted, unchanged |
| Own commit that `main` lacks | `remove` or `recover-remove` | `false:ancestry` | unchanged | kept, unchanged |

In each shape, the `action`, `tracked`, `recovery`, and `detail` cells and the exit code stay the same. Only the second shape gets a new `fingerprint` cell, because its `landed` evidence changes. The apply refuses a fingerprint from a plan made before this change as stale. With `--discard-branch`, each shape plans as before, and the detail names `discards branch <ref>`.

`subshell.go`, `classifier.go`, `lifecyclepolicy.go`, and `lifecyclepolicy_test.go` are over their line budgets. Do not grow these files net. The `internal/worktree/` directory is over its file budget, so add no file there.

The central-property mutation for a reviewer probe removes the own-output clause from the policy rule. The rule then returns the ancestry proof for every assignment. A second probe deletes the output-fact query from the explicit planner. At base `63e28eac`, that second mutation reds nothing, because the classifier derives the verdict again.

## Decision required

The choice: does a bare `bench worktree clean <path>` delete the branch of an assignment with no own commit and uncommitted output?

- Recommended: yes, as today. The branch holds no commit that `main` lacks, and the recovery ref keeps the output. After the apply, the assignment is `recovered` and its registration is gone. No Bench verb can then delete a kept branch.
- Alternative: no. Branch deletion follows the landedness verdict, so this shape needs `--discard-branch` to delete the branch. The list row and the release refusal already name that route. The plan fingerprint then also changes its branch-deletion fields, and the receipt records no branch. The second acceptance row changes to "the apply keeps the branch".

## Acceptance

- [ ] The explicit plan of an assignment with no own commit and an untracked file plans `recover-remove` and carries the landedness `false:ancestry`.
- [ ] The apply of that plan without `--discard-branch` removes the tree, keeps the untracked file in the recovery ref, and deletes the branch.
- [ ] The release of that assignment still refuses as `unmerged` and names `bench worktree clean --discard-branch <path>`.
- [ ] The explicit plan of an assignment with no own commit and a clean tree carries `true:ancestry`, and its apply deletes the branch.
- [ ] The explicit plan of an assignment with a landed own commit and an untracked file carries `true:ancestry`, and its release still refuses as `dirty`.
- [ ] The independent fact gathering of the eligibility test gathers the output fact, and its dirty subtest matches the plan landedness `false:ancestry`.
- [ ] Without the own-output clause in the policy rule, the plan landedness row, the release row, and the list cell row go red.
- [ ] Without the output-fact query in the explicit planner, the release row goes red.
- [ ] `CHANGELOG.md` records the change.
