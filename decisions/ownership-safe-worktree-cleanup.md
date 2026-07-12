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

Unattended cleanup may salvage and remove Git-visible dirty state when the
ownership and assignment identities match, no live lease remains, and the dedicated
assignment branch is proven landed. Bench first creates and verifies the synthetic
recovery commit and ref, then unlocks and removes the worktree. Any classification,
preservation, or removal-precondition failure keeps and reports the worktree. Ignored
residuals remain governed by #8.

## #8: What happens to ignored residuals?

Blocked by: #4, #6
Type: Grill

### Question

Ignored paths may be disposable build output or local secrets and are not captured by
the recovery commit. Decide whether cleanup refuses them, archives them, or permits
their disclosed deletion through a separately bound plan action.

### Answer

Ignored residuals always block unattended removal. Explicit cleanup may delete them
only through a separately fingerprinted `--discard-ignored` plan action whose bounded
inventory reports counts, size, and truncated path detail without reading file
contents. Inventory uncertainty or drift refuses removal. Bench never archives
ignored residuals because they may contain secrets or unbounded caches.

## #9: When do assignment branches and records leave?

Blocked by: #3, #4, #7
Type: Grill

### Question

Set the lifecycle for landed assignment branches, completed intent records, stale
markers, and retained recovery refs after the checkout disappears.

### Answer

Successful cleanup removes only the exact assignment branch after landedness proof;
branch-deletion or worktree-removal failure leaves the assignment record in a
cleanup-pending state and reports the residual. Git removes the immutable ownership
marker with the worktree registration. An assignment with no recovery ref completes
and compacts after its worktree and branch are gone. An assignment with recovery refs
transitions to terminal `recovered`, retaining its implementation context until the
last associated ref is explicitly deleted; only then may the record compact.

## #10: What must each cleanup surface report?

Blocked by: #1, #4, #6, #7, #8, #9
Type: Grill

### Question

Define the observable results and fail posture for SessionStart, `resume-clean`, dry
run, apply, partial preservation, and state drift without duplicating classification
logic in their renderers.

### Answer

One typed cleanup result is the single source for classification, preservation,
actions, residuals, and failures. Explicit dry-run and apply emit a minimal stable
TOON plan/outcome schema on stdout, including the plan fingerprint and bounded detail;
their errors use the same schema. SessionStart and `resume-clean` render one compact
operational summary from that result. Policy-preserved and idempotent no-op outcomes
exit 0, failed classification/preservation/removal exits 1, and invalid invocation
exits 2. SessionStart remains non-blocking but does not turn a failure into a false
clean report.

## #11: Does this destination ship as one slice or several?

Blocked by: #7, #8, #9, #10
Type: Grill

### Question

Choose the spec boundary after the ownership, recovery, lifecycle, harness wiring,
and CLI seams are fully visible. Slicing remains the reviewer's call.

### Answer

Use one spec with three ordered stories: (1) the ownership, assignment, landedness,
lock, and recovery core; (2) automatic and explicit cleanup surfaces over that core;
and (3) Claude/multi-harness wiring with preservation contracts and a biting canary.
These are seams of one safety capability rather than independently complete features;
all three must ship before FT77 closes.

## Not yet specified

n/a — exact field names, ref-name encoding, and fingerprint algorithm are
spec-writer choices inside the resolved contracts below, not open product decisions.

## Out of scope

- Identity-safe lease reclamation and private pool-root permissions tracked by FT58.
- Shift-wide failure recovery and result-state semantics tracked by FT79.
- Removing the repository's primary checkout.
- LLM-based ownership, landedness, preservation, or acknowledgement decisions.
- Treating branch-name, directory-name, age, or an observed harness process as proof
  of ownership.

## Handoff

1. **Module boundaries.** `internal/worktree` is the deep lifecycle module: it owns
   Bench-controlled creation, immutable ownership markers, Git locks, classification,
   cleanup planning, recovery commits/refs, apply revalidation, and a typed result.
   `internal/intent` owns the per-assignment persisted schema and terminal states;
   `internal/git` remains the one owner of default resolution and
   `LandedInDefault`. `cmd/bench`, `bin/bench.sh`, SessionStart, and harness hooks are
   thin routers/renderers. The harness-neutral Bench creation surface is primary;
   Claude's lifecycle hooks are one adapter over it. No caller independently infers
   ownership or landedness.
2. **Contracts.** Creation accepts a stable work-item label and produces a dedicated
   assignment branch, worktree, Git-private marker, matching assignment record, and
   Git lock or fails without returning a usable path. Automatic cleanup acts only on
   a matching, non-live,
   landed assignment; foreign, unmerged, ignored-residual, malformed, and uncertain
   states are retained. Dirty Git-visible state is committed through an isolated
   index to a verified recovery ref before removal. `bench worktree clean <path>`
   emits a dry-run plan; `--apply <fingerprint>` revalidates it, and ignored deletion
   additionally requires the fingerprinted `--discard-ignored` action. Explicit
   cleanup may target a foreign exact path. Recovery refs leave only through explicit
   cleanup after landedness proof. Explicit output is stable TOON; exits are 0 for
   success/policy-preserved/idempotent no-op, 1 for failed intent, and 2 for usage.
3. **Deep vs thin.** The worktree lifecycle hides Git admin addressing, identity
   validation, alternate-index salvage, ref ordering, lock ordering, plan hashing,
   and rollback behind create/plan/apply/automatic-clean interfaces. Intent persistence
   and Git landedness stay deep existing collaborators. Shell hooks, routing, and
   status text pass inputs and render the typed result; they have no cleanup policy.
4. **Black-box assertables.** Ordinary branch and detached-unique foreign worktrees
   survive SessionStart, `resume-clean`, and default explicit dry-run. A marker/path
   mismatch, reused path, live lease, lock, unmerged branch, unknown default, or Git
   query failure also preserves the target. Dirty submodules or nested repositories
   whose contents cannot be represented by the parent recovery commit are likewise
   preserved. A matching landed clean target is removed;
   a matching landed dirty target is removed only after a ref-reachable recovery
   commit contains its Git-visible changes without moving its assignment branch.
   Ignored residuals block unattended removal; explicit apply without the separately
   bound discard action refuses them, and inventory drift invalidates the plan. An
   exact foreign plan/apply removes only its target. Recovery-ref and assignment-state
   lifecycles remain discoverable until explicit retirement. Claude creation returns
   a marked, assigned, locked worktree; safe removal unlocks it, while injected hook
   failure leaves the lock and checkout intact. Every outcome has the agreed output
   shape and exit code.
5. **Gate attachment.** Unit tests attach at the worktree lifecycle interface for
   marker parsing, assignment matching, fingerprint stability, state drift, alternate
   index salvage, ref-before-unlock ordering, and injected filesystem/Git failures.
   Built-CLI runtime contracts exercise automatic and explicit paths in throwaway
   repositories, including the shipped SessionStart and Claude hook adapters. A
   behavior-owned canary must go red when ownership matching, recovery-before-removal,
   or the fail-closed lock is bypassed. A fresh Claude Code worktree run verifies the
   real `WorktreeCreate`/`WorktreeRemove` cadence that the gate can only simulate.
6. **Hostile-input owners.** The CLI parser and canonical-path classifier own spaces,
   glob characters, leading dashes, exact-path grammar, symlinks, deep cwd, and
   traversal. Marker/assignment parsers own absent, empty, malformed, wrong-type,
   no-trailing-newline, schema, mode, and identity-mismatch states. The TOON renderer
   owns control bytes and bounded/truncated ignored inventories. The local CLI
   resolver owns no-global-`bench` hook invocation; Git failures remain structured
   lifecycle errors. The transaction owns SIGINT/partial creation or apply, concurrent
   cleanup, state drift, nested-repository/submodule detection, and recovery-ref
   durability. Re-run contracts own idempotent create/plan/apply/remove behavior.
   Cross-surface contracts prove the real CLI,
   linked by-path CLI, SessionStart, and Claude hooks reach the same core.
7. **Uncertainty flags.** None require a higher line. Claude's real lifecycle cadence
   is a dogfood obligation, not an unresolved design choice; Git-lock refusal was
   locally probed, and the hook roles are documented by Claude's official contract.
8. **Rejected alternatives.** Rejected automatic foreign cleanup; post-hoc adoption;
   a drifting central ownership registry; normal detached assignments; ancestry-only
   landedness; generic `--yes`; salvage commits on the user's branch; archiving ignored
   data; unattended recovery-ref deletion; LLM classification or marker creation;
   relying on a non-blockable harness removal hook without a Git lock; and separate
   specs for the core, surfaces, and harness wiring.
9. **Domain watch-outs.** A Git lock is the fail-closed substrate because Claude's
   removal hook cannot block; Bench removes the lock only inside the proven cleanup
   transaction. Ignored does not mean disposable and may mean secret or enormous.
   A parent Git commit cannot preserve dirty content inside a nested repository or
   submodule, so that state is a preservation failure rather than salvage success.
   The Git-private ownership marker disappears with registration, so assignment state
   transitions before removal. A recovery ref must be durably visible before unlock
   or removal. Patch-equivalent non-merge work counts as landed, while squash or merge
   ambiguity deliberately retains work. FT58 owns lease-reclamation correctness and
   pool permissions; FT79 owns shift-wide post-mutation recovery and result states.

Dependency order: one spec, implemented as (1) lifecycle core, (2) automatic and
explicit surfaces, then (3) multi-harness wiring, black-box preservation contracts,
and canary. No story alone closes FT77.
