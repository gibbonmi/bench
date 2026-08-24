# Session handoff

Repository: `6351a9e422071a56729e899496249170-6351b5e29c34ff7e8496da98891c5834` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/6351a9e422071a56729e899496249170-6351b5e29c34ff7e8496da98891c5834`
Branch: worktree `ignore-capture-inboxes` — base `26fac60f` on `main`, one light-path commit staged for one landing
Spec: none staged; the tickets-only folder `specs/ignore-capture-inboxes/` closes at this landing
Gate: green in the worktree before the commit

## State

Git now ignores `capture/IDEAS.md` and `capture/learnings.md`, and both left
the index. The files stay on disk as local working notes, and the drain still
reads them. When the ideas inbox is ignored, `bench idea` writes the primary
checkout's copy from any checkout. When the inbox is tracked, the verb keeps
the primary-checkout refusal and the worktree-local write, so a linked repo
keeps the old behavior.

The primary checkout's inbox copies were preserved before the landing and
restored after it, because the merge deletes tracked copies. Linked repos do
not receive this ignore rule yet; that decision is open.

## Next command

`/bench-implement-spec` — FT238, per the recommended sequence.

## Shape

Rewritten in full at every phase close, pruned rather than accreted. A fresh
session pays for every line it reads cold; drop anything it would not act on.

Operational gotchas are placed by lifetime, not copied here. One that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build
belongs instead in that spec's coverage rows.

This file names at most when you'll hit one, never the command — a second copy
drifts from the source.

Keep the three sections above. **State** holds what is true now, including anything
uncommitted. **Next command** holds the exact harness-native invocation, not a
description of it. This section is the third.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
