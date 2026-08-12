# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `2747e38b`, clean tree, unpushed backlog (~65 commits; push remains reviewer-owned)
Spec: `specs/axi-query-disclosure/spec.md` (Status: implemented at `2747e38b`); `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `2747e38b` (fresh verdict on the landing)

## State

axi-query-disclosure is closed: review pickup committed (`04ddded7`), the
Opus-accepted nine-ticket repair frontier fully landed (`30c50b50`, `6a38fc6d`,
`d1b2ef19`, `69492c70`, `93fbe8e9`, `5d71b964`, `b18fa260`, fence amendment
`31d63b5f`, final `2747e38b` with the pickup file deleted), and the retro is at
`capture/retros/axi-query-disclosure.md`. Reviewer decisions taken this
session and staying closed: `internal/gittest/` fence amendment; guards
real-stale fixture added; coverage `why` reworded. Uncommitted right now:
`capture/` updates (IDEAS lines, two open learnings entries, the retro, this
file) staged for the close-out commit that immediately follows this rewrite —
if that commit is absent, land those `capture/` paths first.

Parked for the reviewer (in `capture/IDEAS.md`): R11 active-assignment/deleted-
tree spec amendment; census scope for process-backed fixtures (coupled with
gittest census visibility); standing test-support commons fence; pre-existing
`TestListCommandPublicRowsAndDisclosure` flake. Assignment
`4fb88c5a8a96bceffe565d3a540018a1` and all foreign worktrees preserved.

## Next command

`$bench-what-next`

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
