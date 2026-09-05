# Migrate the worktree directory sites to the reader

Blocked by: add-the-named-git-admin-readers.md
Writes: internal/worktree/classifier.go, internal/worktree/lifecycle.go, internal/worktree/ownership.go, internal/worktree/subshell.go, internal/worktree/clean.go, internal/worktree/land_refusal.go, internal/worktree/admin_readers_test.go (new), internal/worktree/parallel_census_test.go, CONTEXT.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary
Covers: GR11, GR12, GR19, GR34, GR35

## What to build

Verify the premise first. Read `AdminDir` in internal/git/worktree_admin.go, which
the blocker ticket adds. Read `validateOwnerMarker` at internal/worktree/classifier.go
line 443. Read `lockCleanupRegistration` and `cleanupLockPath` at
internal/worktree/lifecycle.go line 300. Read `markerPath` at
internal/worktree/ownership.go line 284. Read `planExplicitWith` at
internal/worktree/subshell.go line 142.

Read `recoverAssignmentWithFault` at internal/worktree/clean.go line 32. Read
`sourceMergePending` at internal/worktree/land_refusal.go line 93. Read the
glossary shape of an existing term in CONTEXT.md. Each of the six sites spells
`rev-parse --path-format=absolute --git-dir` or `rev-parse --absolute-git-dir`
today.

Replace the Git call at each of the six sites with `git.AdminDir`. Keep the
current refusal outcome at four sites. Keep the repository-level lock file that
`lockCleanupRegistration` opens when the reader fails. Keep the false answer that
`sourceMergePending` returns when the reader fails. Drop the hand-written
absolute-path guards the readers now own.

Add the term **checkout administration directory** to CONTEXT.md. Define it as
the directory `rev-parse --git-dir` names. State that it is the common directory
for the primary checkout. State that it is the worktree admin entry directory for
a linked worktree. Give it the Avoid list `git dir`, `gitdir`, and `private git
dir`. Change the classifier refusal text and the ownership refusal text to use
the term.

Five of the files this ticket writes sit over the structure budget. The
classifier, the lifecycle, the ownership, the subshell, and the cleanup file must
each end at or under the base line count. The fast lane reds the structure growth
check on a worktree commit that adds a line to one of them.

The five command registries and the eight fixture directories in the `Writes:`
line are closure headroom. They close the worktree binding and the CONTEXT.md
pins, and this ticket edits none of them.

Write the three new tests in internal/worktree/admin_readers_test.go. Name them
`TestLockCleanupRegistrationFallsBackForNonWorktree`,
`TestSourceMergePendingIsFalseWhenUndecided`, and
`TestOwnershipRefusalTextsUseTheGlossaryTerm`. Point the first two at a directory
that is no repository, so the real `git` fails and no stub is needed. Drive the
third over a real repository with the `fail-git-dir` stub on `PATH`. Assert that
the classifier refusal and the ownership refusal each hold the term.

## Acceptance

- [ ] The six migrated sites call `git.AdminDir` and spell no Git flag.
- [ ] The cleanup registration lock returns a release function over a target that is no repository.
- [ ] The cleanup registration lock holds the repository-level lock file in that case.
- [ ] The merge-pending probe answers false over a source that is no repository.
- [ ] `bench structure --growth <base>` reports no over-budget file that gained a line.
- [ ] CONTEXT.md holds the term `checkout administration directory` with its Avoid list.
- [ ] Under the `fail-git-dir` stub, the classifier refusal text and the ownership refusal text each hold `checkout administration directory`.
- [ ] The pre-existing `internal/worktree` suite passes with its test logic unchanged.
- [ ] Self-probe: propagate the typed failure from the cleanup registration lock, and report its new test red.
