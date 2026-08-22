# `bench worktree create` advertises the label-addressed next action

Blocked by: none
Writes: internal/worktree/worktree.go, internal/worktree/*_test.go, .agents/skills/bench-craft-delegate/SKILL.md

## What to build

After the `worktree_create` table, `bench worktree create` prints a `next[2]:`
block. The block names the label-addressed verbs: `bench worktree exec "<label>" -- <command>`
and `bench worktree path "<label>"`. A caller never caches the path in a
scratch variable. The `bench-craft-delegate` skill's Isolation passage states
the same rule. A caller addresses a created worktree by label
through `bench worktree exec`/`path`, never through a cached `$WT`.

## Acceptance

- [ ] `bench worktree create --request r --label "x"` output ends with a `next[2]:`
      block whose lines contain `bench worktree exec "x" -- <command>` and
      `bench worktree path "x"`; a test asserts it.
- [ ] Delegate skill Isolation passage carries the sentence.
- [ ] `bench gate` green.
