# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `2f78e9a`, clean tree, 8 unpushed commits
Spec: `specs/ft126-recurrence-tallying/spec.md` (Status: staged), `specs/ft128-agent-line-binding/spec.md` (Status: staged)
Gate: green at `32285dd` — stale, work tree `dcbbc5b`

## State

- **FT126 implementation and semantic review are complete.** All 33 mapped rows
  landed across `662650d`, `592e253`, `1a02ac4`, and `8a40ed3`.
- **Fresh review `/root/ft126_full_review` found two concrete issues.** The exact
  count-mutation proof and fail-closed degraded-source trust repair landed green at
  `2f78e9a`; its pickup file was deleted in that same commit.
- **The reviewer chose the mid binding and skipped optional falsification.** Top
  escalation and the separate top-binding Codex CLI pass were both declined.
- Final-check is next. The `Status: implemented` transition, default-branch
  retirement, implementation retro, push, and ship-tier verification have not run.

## Next command

`$bench-implement-spec ft126 --full`

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
