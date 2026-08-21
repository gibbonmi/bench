# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `a1cb0ed`, 2 dirty paths, 2 unpushed commits
Spec: `specs/learnings-dated-line-visibility/spec.md` (Status: staged)
Gate: green at `602a9d0` — stale, work tree `e1e379d`

## State

FT243 is implemented, landed, reviewed, retired, and clean. Nothing is in
flight; the integration worktree is released and the pool is empty.

`capture/learnings.md` no longer reports a false zero. `learnings.Parse` reports
a dated line that is not a heading, and content below the exported
`JournalEntriesMarker` that belongs to no entry. `bench learnings`,
`bench roadmap --context`, and `bench status` all flipped without a production
edit of their own; production changed in two files only. The observed red was
flipped through the built binary, not just in tests.

The spec is retired. Its durable hostile-input edge — name the exact whitespace
predicate and assert both sides of it, because `White_Space` and zero-width are
different classes — is promoted into `projects/benchkit.md`'s checklist. That
one-bullet profile edit was written inline rather than at the top tier the
leverage override asks for; flagged for veto.

FT243's `ROADMAP.md` row is deliberately left standing. `bench spec retire`'s
own `next:` line says to remove it, and `/bench-final-check`'s post-merge tail
says roadmap rows belong to `/bench-drain`. The phase text won; the
contradiction is recorded in the retro.

Uncommitted and waiting for the drain: `capture/retros/learnings-dated-line-visibility.md`
and the refreshed `capture/agent-performance/claude-models.md`. The OpenAI
scorecard is untouched, because no OpenAI model served on this landing.

Ten commits are unpushed. Pushing is the reviewer's call.

## Next command

`/bench-drain`

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
