# Session handoff

Repository: `0e3b0d6e38fcbb8d24613567d04cf14b-5d31fc65aa60735401d7a9d314f20918` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/0e3b0d6e38fcbb8d24613567d04cf14b-5d31fc65aa60735401d7a9d314f20918`
Branch: `bench/assign/0e3b0d6e38fcbb8d24613567d04cf14b/5d31fc65aa60735401d7a9d314f20918` — HEAD `8252957`, clean tree, 20 unpushed commits
Spec: `specs/handoff-sections/spec.md` (Status: staged)
Gate: no gate has run.

## State

`/bench-implement-spec specs/handoff-sections/spec.md --full` is in the
build phase. The integration worktree is labelled `handoff-sections`, request
id `impl-handoff-sections-20260904`, over `main` at `6fcd0882`. Its tip
`82529573` carries one fence-closure commit: the guidance ticket now names
the `agents-system-suite-route` fixture, which `binary-freshness` landed
after the spec was staged. `bench preflight build handoff-sections` is green
there. This closure edit is flagged for reviewer veto.

Wave 1 is charged: the leaf package ticket runs in `handoff-sections`, and
the ignore ticket runs in the sibling `hs-ignore` (request id
`impl-hs-ignore-20260904`), both from tip `82529573`. Five tickets remain in
`Blocked by:` order. No ticket has committed.

`specs/pin-removal/spec.md` waits on sign-off in the worktree labelled
`pin-removal` (request id `spec-pin-removal-20260904`). Its veto surface is
in that spec. `specs/agent-push-guard/spec.md` is implemented and not
retired. Three learnings are open, and one retro is pending for the drain.

The user owns the push of `main`.

## Closed decisions

- The agent may push any branch but the default. Force, delete, `--all`, `--mirror`, and `--tags` stay denied. `bench.allowProtectedPush` never lifts the guard.
- The drift clause, `bench gate pin`, the pin file, and the `gate unpinned` warning leave the kit. The branch clause and the config lift stay.
- Every subagent runs `opus` at low or medium effort.
- Do not start FT71 before 2026-09-16.

## Next command

`/bench-implement-spec specs/handoff-sections/spec.md --full`

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
