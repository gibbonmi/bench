# Session handoff

Repository: `dcb9a3236bed02ce621939bf4d8821b5-fa4b747ff712a5b67f4f226c573d84e5` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/dcb9a3236bed02ce621939bf4d8821b5-fa4b747ff712a5b67f4f226c573d84e5`
Branch: `bench/assign/dcb9a3236bed02ce621939bf4d8821b5/fa4b747ff712a5b67f4f226c573d84e5` — HEAD `5f4796d`, 4 dirty paths, 22 unpushed commits
Spec: `specs/worktree-test-latency/spec.md` (Status: staged)
Gate: no gate has run.

## State

The ready decision map and its invocation census moved into the staged spec.
One read-only `gpt-5.6-terra` high review accepted the spec and seven-ticket
graph after two iterations.

The first spec removes deterministic demand. It selects one Bench binary,
introduces explicit effect inputs, extracts landing, lifecycle, and reclaim
policy owners, and contracts tests to representative serial journeys.

The first spec adds no scheduler or `t.Parallel`. It preserves `-count=1`,
records before-and-after demand, and leaves measured pure-test parallelism plus
the slow-package budget to the second spec.

Publication's uncontrolled connection wait and the WSL Go bootstrap defect
remain separate debug tracks. Neither may count as worktree latency improvement.

Tickets 01 and 02 are the initial frontier. Their writes overlap, so work them
serially on one retained integration source even though neither semantically
blocks the other.

## Next command

`$bench-implement-spec --full worktree-test-latency`

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
