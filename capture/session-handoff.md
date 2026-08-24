# Session handoff

Repository: `b288be00ab6c7185ed1ec38ee56d7a85-253130b56b7c625b08d484274fbd1250` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/b288be00ab6c7185ed1ec38ee56d7a85-253130b56b7c625b08d484274fbd1250`
Branch: worktree `local-handoff` — base `a3e84401` on `main`, two light-path commits staged for one landing
Spec: none staged; the tickets-only folder `specs/local-handoff/` closes at this landing
Gate: green in the worktree before each commit

## State

The capture inboxes (`capture/IDEAS.md`, `capture/learnings.md`) are
git-ignored since `a3e84401`; `bench idea` writes the primary checkout's copy
when the inbox is ignored. This landing carries the same treatment for the
session handoff, in two steps. Step one (this landing) ships the code with the handoff still
tracked. The code adds the shared `git.LocalNoteRoot` helper, the
`bench handoff` redirection, the mtime-based status age, and the
tracked-only Shape check. Step two, a follow-up landing, unstages the file so the `.gitignore`
line takes effect; until then the tracked file wins over the ignore rule.

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
file's last write and reports a `handoff` row once anything has landed since.
That write is the commit that carried the file, or the file's own timestamp
when git ignores it. Where this document and the tree disagree, the tree wins.
