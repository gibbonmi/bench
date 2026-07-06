# `bench worktree clean` + honest status actions (FT11)

`bench worktree` ignores every argument and drops into a pool subshell
(`main.go` `case "worktree"` → `worktree.Subshell`), so `bench worktree clean`
— the very action `bench status` recommends for a stray worktree — exits 0 and
silently does nothing. FT11 adds a `clean` verb that removes clean out-of-pool
worktrees and prunes, rejects unknown args instead of swallowing them, and
rewrites the status action so each worktree signal class points at the command
that actually clears it.

## #1: What may `clean` remove, and by what mechanism?

Type: Grill

### Answer
**Scope: out-of-pool worktrees + prune only.** `clean` removes registered
worktrees that are neither the repo root nor under the warm pool, then runs
`git worktree prune`. It deliberately leaves the other two classes `status`
counts alone: **leased pool entries** are live/crashed shift or delegate
state that `Acquire` reclaims automatically on the next lease (a manual sweep
would race a live delegate), and **orphaned `worktree-*` scratch branches**
already have their own `delete scratch branch` action.

**Mechanism: plain `git worktree remove`, never `--force`.** Non-force remove
refuses any worktree with uncommitted or untracked files, so "remove only
clean" comes for free and the verb stays off the `block-dangerous-git`
denylist (`gitguard.go:44` denies `git worktree remove --force`). The
"would remove" candidate list is pre-filtered to out-of-pool ∩ clean so it is
accurate; a candidate that turned dirty between listing and removal is still
refused by the non-force remove and reported, not forced.

Rejected: unify all three status classes under one verb (leased entries are
live state, not garbage — reclaim is automatic; a stale-lease heuristic races
a live owner). Rejected: `--force` removal (on the denylist; destroys
uncommitted work). Rejected: a separate `prune` verb (prune is folded into
`clean`).

## #2: What removal posture — just-remove, or human-attended?

Type: Grill

### Answer
**List, confirm on a TTY, refuse non-interactive stdin** — the same posture as
`bench gate pin`. `clean` prints the clean out-of-pool worktrees it would
remove, requires a typed confirmation token on a TTY, and on non-TTY stdin
refuses with a structured error having removed nothing. A headless agent
cannot sweep worktrees; removal is a deliberate human-attended act.

Removing a clean worktree only drops its checkout directory — commits and
branches survive and the worktree is recoverable by re-adding — so the stakes
are low, but out-of-pool worktrees are more likely a user's own `git worktree
add` checkout than harness scratch (delegates and shifts both use *pool*
worktrees), which is why the confirm step guards against sweeping a checkout
the user still wanted.

Rejected: remove-clean-and-report with no prompt (would drop a user's
deliberate clean checkout without a beat).

## #3: How does the status action become per-class honest?

Type: Grill

### Answer
Today `appendWorktree` sums leased-pool and out-of-pool entries into one count
and emits one action, `resume or clean up (bench worktree)`, which points at a
command that clears neither. With scope narrowed to out-of-pool, the honest
mapping is **one row per class present**:

- **out-of-pool worktree(s)** → `resume or clean up (bench worktree clean)`
  (the verb that actually removes them).
- **leased pool entry(ies)** → an action that does *not* name `clean`
  (e.g. `resume the shift, or leave to reclaim`), because `clean` will not
  touch them.
- **orphaned scratch branch(es)** → `delete scratch branch` (unchanged).

Recommendation, flagged for veto: split the counts by class rather than keep a
single lumped row, so no signal is described by an action that ignores it.
Liveness of a leased entry (dead-owner vs in-flight) stays out of scope —
`status` does not check owner liveness today and the reclaim action is honest
either way.

## Handoff

1. **Module boundaries.**
   - `main.go` `case "worktree"` — arg dispatch: no args → `Subshell`
     (unchanged); `clean` → new `worktree.CleanCommand`; anything else →
     usage error, exit 2. Inside FT11. Outside: the subshell behavior itself.
   - `worktree.CleanCommand(args, stdin, stdout, stderr) int` — new verb:
     enumerate clean out-of-pool worktrees, show them, confirm on TTY / refuse
     non-TTY, remove via non-force `git worktree remove`, then
     `git worktree prune`, report each outcome.
   - A shared worktree-classification seam in `internal/worktree` — returns
     each registered worktree tagged root / pool-warm / pool-leased /
     out-of-pool. **One source per fact:** `status.appendWorktree` and
     `CleanCommand` must derive "out-of-pool" from this one function, not two
     copies of the pool-prefix + lease-file logic.
   - A shared TTY-detect + typed-confirm helper — `stdinIsTerminal` (currently
     unexported in `gate/pin.go`) gets its second caller here; promote it to a
     shared home so the `ModeCharDevice` check has one source. `CleanCommand`
     mirrors `pinCommand`'s injectable `isTerminal func(io.Reader) bool` seam
     for testability.
   - `status.appendWorktree` — emit per-class rows/actions (#3).

2. **Contracts.**
   - `bench worktree` (no args): unchanged subshell.
   - `bench worktree clean`, TTY: prints the clean out-of-pool candidates,
     prompts for a typed token; on confirm removes each (non-force) + prunes,
     reports removed/refused paths, exit 0; on decline removes nothing, exit
     non-zero. Nothing to clean → a "nothing to clean" line, exit 0.
   - `bench worktree clean`, non-TTY stdin: structured error, removes nothing,
     exit non-zero.
   - `bench worktree <unknown>`: usage error naming the valid forms, exit 2 —
     never falls through to the subshell.
   - Removal never uses `--force`; a dirty candidate is refused and reported,
     never destroyed.

3. **Deep vs thin.** `CleanCommand` is deep (enumerate → classify → confirm →
   remove → prune → report). The `main.go` dispatch is thin routing. The
   classification function and the TTY helper are thin shared primitives with
   one source each.

4. **Black-box assertables** (throwaway fixture repo, `<git-dir>`/worktree
   state like the gate-pin verb tests): create a clean out-of-pool worktree →
   `clean` with piped (non-TTY) stdin refuses, worktree still registered, exit
   non-zero; with `isTerminal` forced + confirm token → worktree removed, exit
   0. Dirty out-of-pool worktree → left intact (non-force refuses). Leased pool
   entry present → not offered as a candidate and not removed. `bench worktree
   badverb` → exit 2, error names valid forms. `bench worktree` no args → still
   the subshell (existing test unchanged). Status: out-of-pool present → action
   text names `bench worktree clean`; leased-only → action does not name
   `clean`.

5. **Gate attachment.** `internal/contract/runtime` — beside the gate-cache and
   pin-verb tests that drive `<git-dir>` state; the non-TTY refusal and the
   forced-confirm removal are both gate-drivable via the injected `isTerminal`
   seam. The one gate-blind seam is the literal interactive TTY prompt — same
   as `bench gate pin`: **test the refusal + forced-confirm write paths in the
   gate; the raw interactive prompt is manual-verify, no PTY helper built**
   (FT2's precedent; a PTY helper still has no second caller). Status per-class
   rows are asserted in the existing status runtime contracts. The arg-reject
   and status-action-string behaviors also want canary needles so a rot back to
   swallow-args / stale action text goes red.

6. **Hostile-input owners** (profile checklist):
   - paths/dirs with spaces or globs → the remove loop; Go `exec.Command` args
     don't word-split, but the "would remove" print must show paths intact.
   - unquoted multi-word args → `main.go` dispatch (slice args, no split).
   - tool missing from PATH → n/a, `clean` is Go+git, calls no `bench`.
   - invocation through every surface (kit CLI, linked by-path CLI) reaches one
     impl → shell `worktree)` → `route_porcelain` (unchanged); the status
     action must name the verb that impl actually provides.
   - SIGINT mid-run → confirm precedes any removal, so interrupt-before-confirm
     removes nothing; each removal is atomic and idempotent, and out-of-pool
     worktrees carry no lease to strand.
   - re-run idempotency → second `clean` finds nothing, exit 0.
   - cwd deeper than root → resolve `git.Root` like `Subshell`; a candidate
     that is the current worktree is refused by git, reported and skipped.
   - registered worktree whose dir is gone on disk → `git worktree prune`
     covers it.
   - absent vs empty, trailing-newline → n/a (no file parsing; reads
     `git worktree list --porcelain`).

7. **Uncertainty flags.** None blocking the build. #3's split-into-per-class
   rows is a recommendation flagged for veto, not an open seam — a single
   combined row would still be buildable at the same seam.

8. **Rejected alternatives.** Unify all three status classes under `clean`;
   `--force` removal; a separate `prune` verb; remove-without-confirm; checking
   leased-entry owner liveness. (See tickets #1–#3.)

9. **Domain watch-outs.** Out-of-pool worktrees are more often a human's own
   `git worktree add` checkout than harness scratch — delegates and shifts use
   pool worktrees, so `clean`'s candidate set skews toward user-created work,
   which is what the non-force refusal and the confirm step both guard.
   Non-force `git worktree remove` refuses a worktree with uncommitted *or*
   untracked files and refuses the currently-checked-out worktree — the safety
   floor is git's own, not a bench re-implementation. Removing a worktree drops
   only its checkout directory; commits and branches are untouched.

Dependency order: n/a — single spec.
