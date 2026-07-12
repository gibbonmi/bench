## Destination

Make every destructive worktree cleanup prove Bench ownership and landed or
recoverable implementation state. Keep the lifecycle deterministic and
multi-harness: code creates, classifies, preserves, and removes; no LLM does.

## #1: Which cleanup path may cross the ownership boundary?

Type: Grill

### Question

Separate unattended cleanup authority from an explicitly acknowledged request.

### Answer

SessionStart and `resume-clean` may remove only an unlocked worktree whose Bench
ownership and current assignment both validate. Foreign, missing-record, malformed,
unmerged, and otherwise ambiguous worktrees are preserved and reported. Explicit
cleanup may target a foreign worktree, but only by exact path through dry-run and an
affirmative application of that exact plan.

## #2: What proves that Bench owns a worktree?

Type: Grill

### Question

Choose one ownership source that cannot be confused by path reuse or drift.

### Answer

An immutable marker in the worktree's Git-private admin directory contains a schema,
random identity, and canonical path and is verified against the current Git
registration. Bench writes it only during a Bench-controlled creation operation. A
harness-created worktree counts only when its creation is routed through that same
deterministic operation; observing an existing checkout never adopts it. Existing
unmarked worktrees remain foreign.

## #3: What identifies the implementation and proves it landed?

Type: Grill

### Question

Give cleanup a machine-verifiable work reference without making mutable assignment
state part of physical ownership.

### Answer

The ownership marker remains immutable. A separate per-use lease or intent record
binds the current assignment to its starting commit and an exact dedicated branch;
every Bench-controlled assignment starts on such a branch. Cleanup requires the
ownership and assignment identities to match. It reuses `LandedInDefault`: ancestry
or complete non-merge patch equivalence counts as landed, while squash, merge-only,
missing-default, query-error, and other ambiguity fail closed.

## #4: What preservation precedes removal?

Type: Grill

### Question

Preserve detached commits and Git-visible dirty state without rewriting the user's
branch.

### Answer

A detached HEAD receives a durable Bench recovery ref before removal. Git-visible
dirty state—staged, unstaged, conflicted, deleted, or untracked non-ignored paths—is
snapshotted as a synthetic commit through an isolated index and anchored by a Bench
recovery ref; the assignment branch is not moved. Recovery refs are never deleted by
SessionStart or `resume-clean`. Explicit cleanup may delete one only after showing it
and proving its commit landed in the resolved default branch.

## #5: How does the lifecycle stay fail-closed across harnesses?

Type: Grill

### Question

Compose one deterministic lifecycle with harnesses whose native cleanup cannot be
blocked.

### Answer

One Bench core operation creates the dedicated branch, worktree, ownership marker,
assignment record, and Git lock atomically enough to fail creation on any incomplete
state. Claude Code's `WorktreeCreate` command hook routes creation through it. Its
`WorktreeRemove` hook calls the matching Bench cleanup operation; Bench unlocks and
removes only after its checks succeed. If the hook fails or preservation is not
proven, the lock remains and ordinary Git removal fails. This lock is the substrate
backstop because Claude's removal hook itself cannot block removal. Other harnesses
use the same core through their available Bench worktree surface, not a parallel
policy. The current Claude hook behavior is documented in the
[official hook contract](https://code.claude.com/docs/en/hooks).

## #6: How is explicit destructive acknowledgement bound to inspected state?

Type: Grill

### Question

Prevent a generic confirmation from authorizing a different worktree state.

### Answer

`bench worktree clean <path>` is a non-interactive dry run that emits a deterministic
plan fingerprint. Execution requires `--apply <fingerprint>` and revalidates
ownership, assignment, HEAD, dirty state, ignored residuals, and recovery actions.
Any drift invalidates the plan and emits a new dry run; a generic `--yes` is not an
acknowledgement.

## #7: May unattended cleanup salvage dirty but landed work?

Blocked by: #3, #4, #5
Type: Grill

### Question

Dirty worktrees must not accumulate, but SessionStart would create a synthetic
recovery commit without a reviewer applying a dry-run plan. Decide whether proven
landed assignment state is enough authority for that unattended salvage and removal,
or whether every dirty target requires explicit plan/apply.

### Answer

— (open)

## #8: What happens to ignored residuals?

Blocked by: #4, #6
Type: Grill

### Question

Ignored paths may be disposable build output or local secrets and are not captured by
the recovery commit. Decide whether cleanup refuses them, archives them, or permits
their disclosed deletion through a separately bound plan action.

### Answer

— (open)

## #9: When do assignment branches and records leave?

Blocked by: #3, #4, #7
Type: Grill

### Question

Set the lifecycle for landed assignment branches, completed intent records, stale
markers, and retained recovery refs after the checkout disappears.

### Answer

— (open)

## #10: What must each cleanup surface report?

Blocked by: #1, #4, #6, #7, #8, #9
Type: Grill

### Question

Define the observable results and fail posture for SessionStart, `resume-clean`, dry
run, apply, partial preservation, and state drift without duplicating classification
logic in their renderers.

### Answer

— (open)

## #11: Does this destination ship as one slice or several?

Blocked by: #7, #8, #9, #10
Type: Grill

### Question

Choose the spec boundary after the ownership, recovery, lifecycle, harness wiring,
and CLI seams are fully visible. Slicing remains the reviewer's call.

### Answer

— (open)

## Not yet specified

- Exact marker and assignment schemas, recovery-ref namespace, and plan-fingerprint
  encoding after the remaining lifecycle policies close.
- The smallest shared classifier/result type that keeps all renderers thin.

## Out of scope

- Identity-safe lease reclamation and private pool-root permissions tracked by FT58.
- Shift-wide failure recovery and result-state semantics tracked by FT79.
- Removing the repository's primary checkout.
- LLM-based ownership, landedness, preservation, or acknowledgement decisions.
- Treating branch-name, directory-name, age, or an observed harness process as proof
  of ownership.
