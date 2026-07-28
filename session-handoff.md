# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `16cc9b5`, clean tree, 4 unpushed commits
Spec: none staged.
Gate: green at `16cc9b5`.

## State

- **FT152 is done and retired.** The gate ran green on `main`; the build was
  already landed and status-flipped, so `/bench-final-check` was an honest
  no-op. `bench spec retire implement-spec-full-run` then removed the spec and
  its decision map in `16cc9b5`, gate-green. Ship-tier verification has not run;
  `bench prep-release` covers that once per release.

- **What the retirement promoted before deleting.** One cold-session note in
  `projects/benchkit.md`: `.claude/commands` is a git-tracked symlink to
  `../.agents/commands`, not a copied tree — two FT152 artifacts assumed a copy
  and specified mirror work that did not exist. Everything else the spec decided
  is enforced in the tree by its own anchors, so it needed no second home.

- **Five entries are parked in `IDEAS.md` for the next drain.** Three are FT152
  deferrals that would have died with the spec: the section-scoped-anchor
  fixture gap (a clause for FT156), the model-token sweep over shipped command
  prose, and a `bench` subcommand for the full run. Two are findings the
  retirement surfaced — see below.

- **Retirement broke a pointer in FT154's map, and that is the one thing worth
  reading before FT154 is specced.** `decisions/slice-unit.md` names
  `implement-spec-full-run` as the single staged spec to migrate to the
  `specs/<slug>/spec.md` layout, and as the next queued build. There are now zero
  flat specs, so that migration story is a no-op or changes shape. Parked, not
  fixed — it is spec-level design, so it is the reviewer's call.

- **A second finding: two sources give the retiring session opposite orders.**
  `bench spec retire` prints "remove the ROADMAP row"; `/bench-final-check`'s
  post-merge tail says to leave roadmap rows to `/bench-what-next`. I followed
  the command file, and `bench status` agrees — it now flags `1 row names a
  retired spec → /bench-what-next`. Parked for the drain to verdict.

## Next command

`/bench-what-next` — the drain owns the stale FT152 roadmap row and the five
parked ideas, and its verdict on the FT154 map pointer is what unblocks the next
build.

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
