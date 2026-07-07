# Worktree branch sweep

Status: staged

## Problem

Eight orphaned `worktree-*` delegate scratch branches sit in the repo with no
sanctioned way to remove them. `bench status` surfaces them and tells the reviewer
to "delete scratch branch", but no `bench` subcommand deletes branches, and the
agent-facing git guard denies `git branch -D` — so the one remedy the dashboard
recommends is a dead end. Merged delegate branches accumulate forever, and the
only escape is a manual `git` command the guard blocks for the agent and that the
reviewer must run by hand.

## Solution

`bench worktree clean` gains ownership of the branch sweep. On every invocation it
deletes a delegate scratch branch when — and only when — all three hold: the branch
is scratch-named (`worktree-*`), no live worktree has it checked out, and its tip is
fully merged into the default branch (no unique commits). It prints
`deleted branch <name>` for each one removed. A scratch branch carrying unique
commits is left in place and reported as
`kept branch <name> (unique commits — inspect or delete by hand)`, so no unreviewed
work ever disappears. The CLI is the sanctioned owner of the delete, so the git
guard stays exactly as it is. The `bench status` orphaned-branch row changes its
action from the dead-end `delete scratch branch` to the now-invocable, now-truthful
`bench worktree clean`.

## User stories

1. As the reviewer, I want `bench worktree clean` to delete every fully-merged
   `worktree-*` scratch branch that no live worktree holds, printing
   `deleted branch <name>` for each, so orphaned delegate branches stop
   accumulating and the recommended remedy actually works.
   Line: claude-opus-4-8 / medium. Deciding that a branch is fully merged into the
   default branch is a merge-base ancestry judgment whose correctness is the entire
   safety of an irreversible delete, so it takes the mid tier the profile routes
   oracle-correctness work to.

2. As the reviewer, I want an unmerged scratch orphan left untouched and reported as
   `kept branch <name> (unique commits — inspect or delete by hand)`, so a branch
   holding unreviewed commits is never destroyed by the conservative default.
   Line: claude-opus-4-8 / medium. The keep is the other outcome of the same
   ancestry check, and getting the boundary wrong would either destroy unreviewed
   work or never clean anything, so it shares the mid tier.

3. As the reviewer, I want a scratch branch checked out in a live worktree, and any
   branch that is not scratch-named (a plain branch, or a `bench/shift-*` review
   branch), left completely untouched, so the sweep never races a live delegate or
   deletes a real branch.
   Line: claude-sonnet-5 / low. This safety property is delivered entirely by
   consuming the existing `OrphanedDelegateBranches` filter, so it is wiring at a
   known seam that the cheap tier handles.

4. As the reviewer, I want the command to exit 0 whether it deleted branches, kept
   branches, or found none, so keeping unmerged work is never treated as a failure
   and the command composes cleanly in scripts.
   Line: claude-sonnet-5 / low. Mapping the three outcomes to a single success exit
   code is exit-code plumbing the gate fully observes, so it routes cheap.

5. As the reviewer, I want the branch sweep to run non-interactively, independent of
   the out-of-pool worktree removal's TTY confirmation, so the sweep works in a
   headless invocation and never requires a terminal to remove a merged branch.
   Line: claude-sonnet-5 / low. Making the sweep run outside the confirmation path is
   control-flow wiring at the existing `cleanCommand` seam and is fully
   gate-observable, so it routes cheap.

6. As the reviewer, I want the `bench status` orphaned-branch row to recommend
   `bench worktree clean` instead of the dead-end `delete scratch branch`, so the
   surfaced remedy is invocable and truthful.
   Line: claude-sonnet-5 / low. This is a single action-string change in
   `internal/status` with a status contract already asserting the row, so it routes
   cheap.

7. As the reviewer, I want the sweep to refuse loudly and delete nothing when the
   default branch cannot be resolved, so a repo with no resolvable default branch can
   never have every orphan silently reported clean or wrongly deleted.
   Line: claude-opus-4-8 / medium. The false-empty guard is a correctness decision
   about not reporting an all-clean sweep when mergedness could not be computed, which
   is exactly the oracle-correctness work the profile routes to the mid tier.

## Implementation decisions

- **Module boundaries (off the map's Handoff).** `internal/worktree` owns the sweep,
  fed by its existing `OrphanedDelegateBranches`, which already returns only
  `worktree-*` refs that no registered worktree holds — so the "checked-out survives"
  and "non-scratch untouched" safety comes for free from the detection seam, not new
  logic. `internal/status` changes only the orphaned-branch row's action string.
  `bin/bench.sh` routing is unchanged: `worktree clean` already routes through
  `route_porcelain`.
- **Deep vs thin.** The one deep decision is *merged into the default branch* — a
  merge-base ancestry check (`git merge-base --is-ancestor <branch> <default>`),
  never a string or name compare. Everything else is thin: the detection filter
  already exists, and the delete is a single git call once safety holds.
- **The delete.** After the ancestry check confirms merged-into-default, the branch
  is removed with a forced git branch delete. The ancestry check is the sole safety
  gate; the force flag makes the delete deterministic even when the repo root's HEAD
  is detached at the default tip or sits on another branch, because git's own
  merged-check for a non-forced delete is HEAD-relative and would falsely refuse a
  branch that is merged into the default branch but not into HEAD. This reads the
  map's "`git branch -d`-equivalent" as "a merged-gated safe delete", with the merged
  gate under our control against the correct ref. *Reviewer veto point: this is the
  one place the map named the flag rather than the guarantee.*
- **False-empty guard.** The default branch is resolved through `git.DefaultBranch`,
  the single source `diff` and `status` already share. Before any branch is
  classified, the sweep verifies the resolved default resolves to a commit; if it
  does not, it prints a loud refusal to stderr, deletes nothing, and the sweep
  contributes a non-zero (exit 1) status — never a silent all-clean report, which is
  the FT29 false-empty failure applied here from day one. *Reviewer veto point: the
  map pinned "report loudly, delete nothing" but left the exit code open; exit 1 is
  chosen so a scripted caller notices the sweep could not run.*
- **Integration into `cleanCommand`.** The sweep runs unconditionally after the
  in-repo guard, as an independent phase that needs no TTY and no confirmation, so it
  is observable in a non-interactive invocation with no out-of-pool worktrees to
  remove. The command's exit code is the higher severity of the sweep and the
  existing worktree-removal phase. The sweep reads the orphan set once; a branch
  freed by that same run's worktree removal is swept on the next `bench worktree
  clean` (see Out of scope), which re-run convergence already covers.
- **Output streams.** `deleted branch <name>` and
  `kept branch <name> (unique commits — inspect or delete by hand)` are results and
  print to stdout, matching the existing `removed`/`refused` lines. The unresolvable-
  default refusal is an error and prints to stderr.

## Testing decisions

- **What a good test is here.** Drive the built `bench worktree clean` (and `bench
  status`) in a throwaway fixture repo and observe external behavior: which branches
  survive `git branch`, the exact printed lines, and the exit code — never a reading
  of the diff or an internal call.
- **Seams and prior art.** The primary seam is the runtime worktree-clean contract
  family (`TestRuntimeWorktreeContracts`), which already exercises confirmed cleanup,
  non-TTY refusal, dirty refusal, stale-registration prune, and pool-cwd
  classification against the built binary — the sweep behaviors attach as new cases
  there. The status action-string change attaches to the runtime status contract
  (which already stands up an orphaned scratch branch) and the `status_test.go` unit.
  The Go unit `branches_test.go` already pins `OrphanedDelegateBranches` detection; the
  sweep's classification is observed end-to-end through the CLI, so it needs no second
  seam.
- **Gate.** `.bench/gate.sh` (the project gate). The runtime contracts run under its
  runtime-and-behavior-contracts phase against `dist/bench`, so the binary is rebuilt
  before the contracts run.

### Seam diagram

Seam 1 — the branch sweep behind `bench worktree clean`:

    trigger: reviewer runs `bench worktree clean` (the action bench status now recommends)
        │
        ▼
    repo root      ──▶  [ cleanCommand ─▶ branch sweep                       ]
    default branch ──▶  [   OrphanedDelegateBranches(root): worktree-*,      ] ──▶ stdout: deleted branch <name>
                        [     not checked out in any live worktree           ] ──▶ stdout: kept branch <name> (unique commits — inspect or delete by hand)
                        [   merge-base --is-ancestor gate ─▶ forced delete   ] ──▶ stderr + exit 1: loud refusal when default branch unresolvable
                             ◀ tests attach here: the runtime worktree-clean contract runs the built bench in a
                               throwaway fixture repo and asserts branch presence/absence, the printed lines, and exit code

Seam 2 — the status orphaned-branch action string:

    trigger: SessionStart hook or reviewer runs `bench status`
        │
        ▼
    orphan count > 0 ──▶ [ appendWorktree row builder ] ──▶ stdout row action: bench worktree clean
                              ◀ tests attach here: status_test.go (unit) and the runtime status contract assert
                                the orphaned-branch row recommends `bench worktree clean`

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | a merged worktree-agent-* orphan is deleted and stdout prints deleted branch name | runtime worktree-clean contract (built bench, fixture repo) | new test commits a scratch branch merged into the default branch, runs clean, asserts the branch is absent from git branch and stdout contains deleted branch; RED on the current binary because no sweep exists so the branch survives and the line is absent | an always-green stub that skips the sweep leaves the branch present and omits the line, failing both assertions |
| 2 | an unmerged scratch orphan survives and stdout prints the kept branch line verbatim | runtime worktree-clean contract | new test gives the scratch branch a commit not in the default branch, asserts the branch present and stdout contains kept branch name (unique commits — inspect or delete by hand); RED on the current binary (no sweep, no line) | a delete-all-by-name impl removes the branch and fails the present assertion; a no-op impl omits the kept line |
| 3 | a scratch branch checked out in a live worktree, and a non-scratch branch (main and a bench/shift-* branch), are never deleted | runtime worktree-clean contract | new test adds a checked-out worktree-agent-* branch plus a bench/shift-* branch and a plain branch, runs clean, asserts all survive; RED on a delete-all-by-name impl | catches the degenerate that deletes by prefix alone or ignores the active-worktree filter |
| 1 | edge: a scratch name with a slash and a non-ASCII character inside the worktree- prefix, when merged, is deleted | runtime worktree-clean contract | new test adds a merged branch named worktree-agent-cafe/x (unicode variant), asserts it is gone after clean; RED if detection or the delete mangles the name so the branch survives | a string-trimming or shell-unsafe interpolation of the ref name leaves the branch present |
| 4 | the command exits 0 when the outcome is a kept (unmerged) branch | runtime worktree-clean contract | new test asserts RequireExit(0) on the unmerged-only scenario; RED on a keeping-is-error impl that exits non-zero | a keeping-equals-failure impl exits non-zero and fails this row |
| 5 | a non-interactive clean (no TTY) with only orphan branches sweeps merged ones and exits 0 | runtime worktree-clean contract | new test runs f.Bench(worktree, clean) with a merged orphan and no out-of-pool worktrees, asserts deletion and RequireExit(0); RED on the current binary and on any impl that gates the sweep behind the TTY confirmation | proves the sweep runs without a terminal and is not coupled to worktree-removal confirmation |
| 6 | the status orphaned-branch row recommends bench worktree clean | runtime status contract plus status_test.go unit | flip the existing assertion from delete scratch branch to bench worktree clean; RED on the current binary which still emits the old string | the new substring assertion fails until the action string is changed |
| 7 | an unresolvable default branch yields a loud stderr refusal, deletes no branch, and exits 1 | runtime worktree-clean contract | new test in a fixture whose resolved default branch ref does not resolve to a commit, with a scratch branch present, asserts nothing deleted, a loud substring on stderr, and RequireExit(1); RED on an impl that runs the ancestry check anyway (branch silently kept, no loud line) or one that treats the error as merged and deletes | pins the false-empty guard: no silent all-clean report and no false deletion when mergedness cannot be computed |

Cheapest wrong implementations checked against the map: an always-green stub fails
rows 1 and 5; a delete-all-by-name impl fails rows 2 and 3; a keeping-is-error impl
fails row 4; leaving the status string unchanged fails row 6; swallowing the
default-branch failure fails row 7.

### Edge inventory

Walked per behavior; each resolves to a coverage row above or a **Won't handle** line
here.

- Error path (git delete fails after classification) — **Won't handle** as a dedicated
  assertion: classify-before-attempt removes the only common cause (active branches are
  already excluded), a delete that still errors is reported to stderr rather than
  swallowed, and the concurrent-checkout race is non-deterministic to force under git's
  ref lock.
- Boundary (scratch tip equals the default tip, zero unique commits) — covered by row 1;
  the ancestry gate treats a tip-equal branch as merged and deletes it.
- Malformed input (slash and non-ASCII in the ref name) — covered by the row-4 edge row.
- Empty/absent input (no orphan branches at all) — covered: the non-interactive clean
  prints no branch lines and exits 0, a subset of the row-5 scenario and the existing
  nothing-to-clean path.
- Interrupted/partial (SIGINT mid-sweep) — **Won't handle**: per-branch git deletes are
  atomic and idempotent, so an interrupt leaves a consistent partial set with no
  partial-branch state to repair, and the next run converges.
- Re-run idempotency — covered: a second clean finds the deleted branches already gone
  and prints nothing. **Won't handle** in one pass: a branch freed by the same run's
  out-of-pool worktree removal is swept on the next run, not the same run — the sweep
  reads the orphan set once and re-run converges (see Out of scope).
- Hostile environment (control bytes in git-sourced text, toon refusal) — **Won't
  handle / N-A**: git refname rules forbid control bytes and spaces in branch names,
  and the sweep prints plain lines rather than a toon.Table, so the control-byte
  refusal path does not apply to this surface.
- cwd deeper than the repo root — covered: cleanCommand resolves the root via
  git.Root(), and the existing pool-cwd contract already drives clean from a non-root
  cwd.
- Invocation surface — covered: the command routes through route_porcelain (unchanged),
  and the runtime contract drives the built binary the linked-repo by-path CLI resolves.

## Out of scope

- **Force-deleting unmerged scratch branches** (an opt-in mode that removes a scratch
  orphan carrying unique commits): a distinct capability with its own safety decision,
  deliberately excluded by the conservative merged-only default the decision map chose;
  ~2 edits, 2 gate runs (a flag plus a contract case).
- **Same-run sweep of a branch freed by that run's out-of-pool worktree removal** (the
  coherence variant where removing a delegate worktree in the same invocation makes its
  branch immediately sweepable): a refinement only — the sweep reads the orphan set once
  and re-running clean already converges; ~2 edits, 1 gate run.

The decision map's other alternatives — loosening the git guard so the agent can delete
branches, and a separate `bench branch` subcommand — remain **rejected**, not deferred,
and are not parked here.
