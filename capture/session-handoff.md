# Session handoff

Repository: `6dd12e02b9a9ed79f2e18b6faaaa9570-3a9293717c5b8eab09f3bd47bc5205a2` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/6dd12e02b9a9ed79f2e18b6faaaa9570-3a9293717c5b8eab09f3bd47bc5205a2`
Branch: `bench/assign/6dd12e02b9a9ed79f2e18b6faaaa9570/3a9293717c5b8eab09f3bd47bc5205a2` — HEAD `76cff41`, 4 dirty paths, 6 unpushed commits
Spec: `specs/asd-ste100-progressive-disclosure/spec.md` (Status: staged)
Gate: green at `5c776b5` — stale, work tree `c4b2a89`

## State

Implementation of `asd-ste100-progressive-disclosure` is complete on this worktree, the
retained integration source. Its label is `ste-prose-progressive-loading`; run every
command as `bench worktree exec ste-prose-progressive-loading -- <command>`. All 41
tickets landed green, one commit per ticket, except tickets 15 to 18, which landed as one
commit. The review base is `b1abde47` (the merge of `main` into this branch before ticket
01). The source tip is the commit that closes the chain (ticket 29).

The next phase is the three-axis review over that base and tip, then the repair-scoped
re-review, then `bench worktree land` from the destination, then `/bench-final-check`.

Orchestrator judgment calls are flagged for post-hoc veto. Long anchor needles were
re-cut into two or more rows each, with fixtures. The files are CONTEXT.md, craft-spec,
craft-delegate, craft-tdd, craft-tickets, bench-write-spec, bench-review-implementation,
bench-final-check, and bench-drain.

One profile bullet's fixture string moved in `fixture_bite_test.go`. One
learnings testdata line number moved with the longer preamble. One row in
`registry_data_test.go` was re-pinned. Tickets 07 to 10 were authored in parallel and
landed in chain order.

Four open learnings wait for the drain. A piped `bench commit` hid two red gates. A
`Supports:` field must stay on one physical line. The shared-rule check is a raw
substring. A short label line with no terminator is a field line and is not graded. The
delegate patches and the batch-check script live under this session's scratchpad.

## Next command

`/bench-implement-spec specs/asd-ste100-progressive-disclosure/spec.md --full --reviewer opus`

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
