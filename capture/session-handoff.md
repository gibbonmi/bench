# Session handoff

Repository: `b30c67dabf3a05462644ce3a222f1997-e88bc453ebc75b26ffb8060ff52d32fa` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/b30c67dabf3a05462644ce3a222f1997-e88bc453ebc75b26ffb8060ff52d32fa`
Branch: `bench/assign` worktree `spec-amend-r13` — HEAD `a792d5e1`, clean tree, base `888babfe` on `main`
Spec: `specs/module-size-split/spec.md` (Status: staged)
Gate: green at `a792d5e1` — current

## State

Batch 2 of `module-size-split` landed on `main` at `888babfe`; tickets 03 to 12
are done and `bench structure` reports 56 issues. The R13 decision is closed:
on 2026-08-24 the reviewer widened the scope by one file instead of amending
R13 to 56. The amendment at `a792d5e1` rewrites ticket 13 to split
`internal/lines/lines_test.go` (612 lines). It adds story 27 and row R20,
re-fences the ticket to `internal/lines/`, and keeps R13 at 55. One Sonnet review round
returned two blocking consistency findings; both are folded. The gate and
`bench coverage --check` are green on the amendment.

Ticket 13 is the only open ticket. It has no blocker, and its fence is disjoint
from every landed ticket. The amendment sign-off is the pending reviewer action.
The earlier veto flags (the `checks_test.go` cap at 692, the `buildLandingBinary`
name, the ungraded `writeHostileSkillRoot` plant) stay open.

## Next command

`/bench-implement-spec specs/module-size-split/spec.md --full`

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
