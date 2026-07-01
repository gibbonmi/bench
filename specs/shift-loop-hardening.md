# shift loop hardening

## Problem

`bench shift` owns the autonomous loop, so its tail must be conservative: it should
only ask for refactors caused by the current shift, report commits truthfully, clean
its scratch files on interruption, and stop when the agent makes no refactor change.
The current tail can scan unrelated repo history, report a commit when none happened,
leave scratch files after Ctrl-C, and point at an unbundled command.

## Solution

Keep the implementation loop shape, but harden the post-implementation path:

- Track the files touched by the shift and run the structure refactor phase only on
  that touched scope.
- Clean `.bench-objective` and `.bench-notes.md` on normal exit, interrupt, and
  termination.
- Keep those scratch files from tripping the clean-tree precondition if a prior run
  was interrupted.
- In the refactor loop, commit only when there is a staged change, say "no change"
  otherwise, and break on that no-op.
- When structure remains over budget, describe the self-contained Bench deep pass
  instead of naming the external improve-codebase-architecture command as if it ships
  with the kit.

## User stories

1. As a user running a shift in a repo with pre-existing structural debt, I want the
   refactor phase to ignore files this shift did not touch, so the loop does not drift
   into unrelated work.
2. As a user reading shift output, I want "committed" to mean a git commit was created.
3. As a user whose agent declines a bad split, I want the refactor loop to stop after a
   no-op pass instead of burning every retry.
4. As a user pressing Ctrl-C to pull the line, I want the next shift to start without
   manual cleanup.
5. As a user seeing unresolved structure debt, I want guidance that names only bundled
   Bench capabilities unless it clearly labels an external optional skill.
6. As a kit developer, I want the gate to exercise these behaviors through the real
   CLI in throwaway repos.

## Implementation decisions

- **Seam:** the `bench shift` CLI contract. Tests run the real `bin/bench.sh shift`
  with a throwaway repo, controlled gate, and tiny shell agent. They assert external
  behavior: commits, output, scratch-file cleanup, and final repo state.
- **Touched scope:** derive touched files from the branch diff against the parent
  commit captured before the shift branch is created. Pass that file list into a
  structure helper rather than changing the default `bench structure` command.
- **Directory crowding:** when checking touched scope, include only directories that
  contain a touched source file. This keeps unrelated crowded directories from pulling
  the refactor phase into unrelated work.
- **Scratch cleanup:** install a trap after scratch files are created. The trap removes
  only Bench scratch files and otherwise leaves the user's working tree alone.
- **Clean-tree precondition:** ignore reserved Bench scratch files in the status check.
  This is defense-in-depth for previously interrupted runs.

## Testing decisions

- Add focused gate checks under the existing `bench shift` block:
  - pre-existing over-budget source outside the shift does not trigger the refactor
    phase;
  - touched over-budget source triggers refactor, a no-op refactor reports no change,
    creates no phantom commit, and exits early;
  - an interrupted shift removes scratch files and a follow-up no-op shift is not
    blocked by leftover scratch state;
  - `.bench/done.sh` can end the loop before the iteration cap;
  - the unresolved-structure fallback does not mention `/improve-codebase-architecture`
    as a bundled command.
- Add a separate `bench worktree` contract: create a pooled worktree, dirty it during
  the subshell, release it, then run `bench worktree` again and prove the same path is
  reset and reused.

## Out of scope

- Changing `bench structure` output or default whole-repo behavior. The manual command
  remains a repo-wide structural audit.
- Making `bench shift` create or enter a pooled worktree. That is a separate product
  decision already tracked in the dogfood docs.
- Testing `bench models` against live network APIs.
