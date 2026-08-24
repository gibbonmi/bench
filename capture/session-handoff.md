# Session handoff

Repository: `ee9e5c3e35405f8f069b565a5bf54054-e6c6a7f19db9134730d8cd924a19bfae` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/ee9e5c3e35405f8f069b565a5bf54054-e6c6a7f19db9134730d8cd924a19bfae`
Branch: worktree `spec-retire-mss` — base `74ef1125` on `main`, retire and capture edits staged for one landing
Spec: `specs/module-size-split/` retired (was `Status: implemented` at `74ef1125`)
Gate: green at `74ef1125` — the landing's own run

## State

The `module-size-split` spec is complete. The final landing `74ef1125` amended
R13's scope and split `internal/lines/lines_test.go` (ticket 13, Opus/low,
first-pass). `bench structure` reports 55 issues, within R13's bound. The three
review axes returned zero findings.

Final-check ran the post-merge tail. The spec folder is retired, and the
roadmap sequence list no longer names the spec. The retro
`capture/retros/module-size-split.md` is rewritten, the Claude scorecard is
refreshed, and three learnings entries are appended. This worktree carries all
of those edits in one retire landing.

A `git stash` on the primary checkout holds the first, misplaced retire bytes;
it is disposable once this landing reaches `main`. The 400-to-700-line
remainder stays a decided out-of-scope follow-up. The reviewer pre-approved the
drain recommendations on 2026-08-24.

## Next command

`/bench-drain`

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
