# Fill the preflight next column

Blocked by: none
Writes: internal/preflight/decision.go, internal/preflight/gather.go, internal/preflight/command.go, internal/preflight/decision_test.go, internal/preflight/command_review_test.go, internal/preflight/gather_test.go, internal/preflight/source_tip_test.go

## What to build

An agent reads a preflight red and gets its remedy in the row. `CheckResult`
gains a `Next` field, and the rendered table becomes
`checks[N]{check,verdict,detail,next}`. `Facts` gains `DefaultBranch` and
`AssignmentTarget`. The gatherer fills `DefaultBranch` from `git.ResolvedDefault`.
It fills `AssignmentTarget` with the id of the active assignment whose canonical
worktree path equals the canonical preflight root.

One red carries a remedy. It is the `base-current` red whose detail reads
`default branch tip is not an ancestor of HEAD`.
Its `Next` is `bench worktree merge --from <DefaultBranch> <AssignmentTarget>`.
The literal `<target>` stands when the assignment id is empty.

This ticket is the only one to write `internal/preflight`, so it carries that
package's gate invariant. It crosses no shared inventory surface.

## Acceptance

- [ ] WF32: `bench preflight review <slug>` on a conformant tree prints `checks[5]{check,verdict,detail,next}`.
- [ ] WF33: `Decide` with a false `DefaultBranchCurrent` and a filled `AssignmentTarget` yields `Next` equal to `bench worktree merge --from main <id>`.
- [ ] WF34: the same red with an empty `AssignmentTarget` yields `Next` equal to `bench worktree merge --from main <target>`.
- [ ] WF35: every green row and every not-applicable row carries an empty `Next`.
- [ ] WF36: every other red `Decide` produces carries an empty `Next`: two `base-current`, two `tip-current`, two `paths-authorized`, `rows-owned`, `rows-membership`, and two `diff-nonempty`.
- [ ] WF37: a preflight behind the default branch prints the row `base-current,red,default branch tip is not an ancestor of HEAD,bench worktree merge --from main <target>` at exit 1.
- [ ] WF38: `Gather` from the worktree path of an active assignment fills `AssignmentTarget` with that assignment's id.
- [ ] WF39: `Gather` from a symlink to that worktree path fills the same id.
- [ ] WF40: `Gather` from the primary checkout leaves `AssignmentTarget` empty.
- [ ] WF41: a preflight whose `--source-tip` is pinned prints `checks[6]{check,verdict,detail,next}`.
- [ ] The gate `test` phase stays green for the whole `internal/preflight` package.
