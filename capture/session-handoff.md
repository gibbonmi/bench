# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — 22 unpushed commits, clean tree.
Gate: green.

## State

FT229's spec is staged at `specs/ft229-hygiene-batch/spec.md` — 46 stories in
seven groups, 33 coverage rows, one named `Not covered` exception, eleven
tickets. Reviewer approved the spec-and-tickets pair on 2026-08-19. It is the
last actionable row of the 2026-08 capability audit portfolio; A1, A2, A3, A6,
and FT228 already landed.

Two forks are closed and stay closed. Every tickets-only spec folder is deleted
in this build, and `bench preflight` gains `--source-tip` on both `review` and
`build` rather than the review phase dropping the word. The reviewer also
confirmed the `build` extension and the retention count of 20 gate run records.

L-30 is out of scope: FT228 already removed the inert Codex invocation key and
added the Claude-side parity grading. The residue count is 37 folders, not the
34 the audit recorded or the 30 quoted during shaping — re-count at ticket entry,
because a light-path landing adds one.

The ticket graph has three edges. `count-tickets-only-folders-in-status` blocks
on `close-the-light-path-ticket-on-landing` for the tickets-only predicate;
`delete-the-landed-tickets-only-residue` blocks on both; and
`name-the-build-command-on-a-cold-session` blocks on
`narrow-the-degraded-guard-rim` only because both write `internal/systemtest`.
The other seven tickets are frontier.

`capture/learnings.md` carries one open entry from this phase's review round.

Nothing is pushed. The push is the reviewer's.

## Next command

`/bench-implement-spec specs/ft229-hygiene-batch/spec.md` — a fresh mid-tier
session on one retained integration source.

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
