# shift loop hardening

## Problem

`bench shift` owns the autonomous loop, so its tail must be conservative: it should
only ask for refactors caused by the current shift, report commits truthfully, clean
its scratch files on interruption, and stop when the agent makes no refactor change.
The current tail can scan unrelated repo history, report a commit when none happened,
leave scratch files after Ctrl-C, and point at an unbundled command.

The pooled worktree cleanup has one more isolation leak: release resets tracked files
and removes untracked non-ignored files, but ignored artifacts survive. A later
interactive `bench worktree` or `bench shift` can observe build/cache state from the
previous lease even though the pool advertises a clean worktree.

## Solution

Keep the implementation loop shape, but harden the post-implementation path:

- Track the files touched by the shift and run the structure refactor phase only on
  that touched scope.
- Clean `.bench-objective` and `.bench-notes.md` on normal exit, interrupt, and
  termination.
- Clean ignored artifacts on worktree release, so reuse starts from a truly clean
  worktree rather than a porcelain-clean one.
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
7. As a user reusing a pooled worktree, I want ignored build/cache artifacts removed
   before the next lease, so hidden state cannot cross shifts.

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
- **Pool cleanliness:** release must clean both untracked non-ignored files and ignored
  artifacts in the pooled worktree. The main checkout is never cleaned; this policy
  applies only to the leased worktree Bench owns.
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
  reset and reused. Include an ignored artifact in that dirty state and prove it is
  gone on the second lease.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Pre-existing structural debt outside the shift does not trigger the refactor phase. | `bench shift` CLI | Already covered by the runtime contract for unrelated over-budget source. | The loop must only react to files touched by the shift. |
| 2, 3 | A no-op refactor reports no staged change, creates no phantom commit, and exits early. | `bench shift` CLI | Already covered by the runtime contract for the touched over-budget no-op refactor case. | It catches the false "committed" output and wasted retry behavior at the CLI seam. |
| 4 | Interrupted shifts remove scratch files and release the worktree lease. | `bench shift` CLI | Already covered by the runtime interrupt contract. | It proves the next shift can start without manual cleanup. |
| 5 | Unresolved structure guidance names only bundled Bench capabilities. | `bench shift` CLI output | Already covered by the runtime contract that rejects `/improve-codebase-architecture` in fallback output. | It catches a handoff that points at an unbundled command as if it ships. |
| 7 | Ignored artifacts created during one pooled worktree lease are absent on the next lease of the same path. | `bench worktree` CLI | Observed red before implementation: a throwaway repo with `ignored/` in `.gitignore`, first lease writes `ignored/leak.txt`, second lease sees that file and exits non-zero. | Current release uses `git clean -fd`, which leaves ignored files behind. |

## Out of scope

- Changing `bench structure` output or default whole-repo behavior. The manual command
  remains a repo-wide structural audit.
- Making `bench shift` create or enter a pooled worktree. That is a separate product
  decision already tracked in the dogfood docs.
- Testing `bench models` against live network APIs.
