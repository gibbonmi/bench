# Resume an active spec stage before a new worktree starts

Blocked by: none
Writes: .agents/commands/bench-write-spec.md
Covers: none

## What to build

On 2026-10-19, a `/bench-write-spec` run started a second FT380 spec stage. A signed-off FT380
spec stage sat unlanded in another worktree at that time. The entry of
`.agents/commands/bench-write-spec.md` has no step that looks for an active assignment for the
same roadmap row before the phase creates its worktree. No other kit file holds a general resume
rule for this case.

Add one rule to the entry orientation of `.agents/commands/bench-write-spec.md`. Before the phase
creates its worktree, it runs `bench worktree list`. When an active assignment's label or request
names the same roadmap row or outcome, the phase resumes that assignment or asks the reviewer.
Keep the file inside its guidance prose budget, and keep every anchored needle.

## Acceptance

- [ ] The entry orientation of `.agents/commands/bench-write-spec.md` tells the phase to run `bench worktree list` before it creates its worktree.
- [ ] The same entry tells the phase to resume a matching active assignment or to ask the reviewer.
- [ ] `bench test --check guidance-prose-budgets` and `bench test --check prose-mechanics` pass.
- [ ] Every anchored needle that `bench anchors` lists for the file still holds.
