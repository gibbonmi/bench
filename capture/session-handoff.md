# Session handoff

Repository: `6dd12e02b9a9ed79f2e18b6faaaa9570-3a9293717c5b8eab09f3bd47bc5205a2` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/6dd12e02b9a9ed79f2e18b6faaaa9570-3a9293717c5b8eab09f3bd47bc5205a2`
Branch: `bench/assign/6dd12e02b9a9ed79f2e18b6faaaa9570/3a9293717c5b8eab09f3bd47bc5205a2` — HEAD `5a9c452`, 7 dirty paths, 6 unpushed commits
Spec: `specs/asd-ste100-progressive-disclosure/spec.md` (Status: staged)
Gate: green at `67b2065` — stale, work tree `7f593fd`

## State

The `asd-ste100-progressive-disclosure` build runs on this worktree, the retained
integration source. Its label is `ste-prose-progressive-loading`; run every command as
`bench worktree exec ste-prose-progressive-loading -- <command>`. Landed green, one commit
per ticket: 01, 01c, 01b, 02, 02b, 03, 04, 11, 12, 13, 14, 19, and 20a–28b. Tickets 15–18
landed as one commit.

Remaining: 05 to 10 (the serial skill and command chain; 05 runs in worktree `ste-05`),
then 29. Then run `/bench-review-implementation` over base `b1abde47` and the source tip,
`bench worktree land`, and `/bench-final-check`.

Orchestrator judgment calls flagged for post-hoc veto: two long needles were re-cut into
two rows each, with fixtures (`coverage row` in CONTEXT.md; five craft-spec needles). One
profile bullet's fixture string moved in `fixture_bite_test.go`. One learnings test line
number moved with the longer preamble.

Three open learnings wait for the drain. A piped `bench commit` hid two red gates. A
`Supports:` field must stay on one physical line. The shared-rule check is a raw
substring. The delegate patches and the batch-check script live under this session's
scratchpad.

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
