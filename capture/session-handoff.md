# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main`

## State

`worktree-exec-run-binary` is landed, published, and retired. Nothing about it
is open.

`bench worktree exec <target> -- ./dist/bench gate` now runs. It refused every
time before: `craft-delegate` addresses a worktree through `exec`, the profile
says to exercise a worktree artifact through that worktree's own `./dist/bench`,
and composed they produced a command the gate entry always refused. The child is
now marked wrapper-rooted so its gate owns a private exact-source build, and the
gate entry's refusal names the wrapper invocation instead of a variable no
operator sets.

The reviewer reopened the resolution fork once, on a `fable`/high finding:
pointing the child at the worktree's `dist/bench` would make it *inherit*, and an
inherited selection is verified against its own seal rather than its source, so a
stale artifact could grade the tree.

**Everything now open is capture, and all of it belongs to `/bench-drain`:**

- 6 open learnings — three from FT229, three from this build. The load-bearing
  one is the landing's undocumented ordering rule: a mid-build spec amendment
  must become the source's base, never a commit inside the reviewed range.
- 4 parked ideas, including the `bench worktree path` tilde form that its own
  sibling verbs reject, and a CONTEXT.md glossary term for the executable the
  gate authorizes and runs.
- 2 pending retros: `ft229-hygiene-batch` and `worktree-exec-run-binary`.
- `capture/agent-performance/claude-models.md`, refreshed and uncommitted.
- `ROADMAP.md` still carries three FT229 references including the FT174
  dependency row FT229 unblocks. FT223 has a fourth occurrence to record.

The retro and scorecard are uncommitted by design — `/bench-final-check` is
forbidden from committing them, and the drain owns their capture commit.

Eight landed worktrees remain, all from FT229. `bench worktree clean --landed`
plans `retain` for every one: each carries uncommitted tracked changes, so the
sweep refuses to remove work. They need per-path resolution by whoever owns them.

## Next command

`/bench-drain`

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
