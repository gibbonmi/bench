# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `e3e8c1e`, clean tree, pushed and level with origin.
Spec: none staged. `specs/` is empty.
Gate: green at `16cc9b5`; two doc-only commits have landed since
(`877b65c` drain, `e3e8c1e` map edit), so the pin is stale but nothing
executable changed.

## State

- **FT152 is done, retired, and pushed.** The gate ran green, the spec and its
  decision map are gone, and the durable content was promoted before deletion.
  Ship-tier verification has not run; `bench prep-release` covers that once per
  release.

- **FT154 is unblocked and is the next build.** Its map, `decisions/slice-unit.md`,
  had assumed one staged spec to migrate to the `specs/<slug>/spec.md` layout;
  FT152's retirement emptied `specs/`, so the migration and the build-sequencing
  premises were struck in `e3e8c1e`. Every decision in the map survived — no
  standing dual-form resolution, `bench spec history` still resolving retired
  flat paths, enumerate-every-glob rather than patch-as-found.

- **Read "mirror" as "symlink" anywhere under `.claude/`.** `.claude/commands`
  is one directory-level symlink; `.claude/skills` is a real directory of
  per-skill symlinks. Neither holds copies, so editing `.agents/` is the whole
  edit — but adding a skill still needs its own `.claude/skills/` entry created,
  which FT154's build will hit. The profile's cold-session notes carry this.

- **FT161 (LOW) is open and needs a ruling, not a build.** `bench spec retire`
  tells the retiring session to remove the ROADMAP row;
  `/bench-final-check`'s post-merge tail says to leave roadmap rows to
  `/bench-what-next`; `/bench-what-next`'s own prose calls itself "the backstop
  for anything spec-retire missed", which leans the third way. Three sources,
  one of them wrong. Which yields is the reviewer's call, then a one-line kit
  edit to the loser at `internal/spec/spec.go:336`.

- **The 2026-07-28 drain landed at `877b65c`** — `IDEAS.md` and the journal are
  empty. It created FT159 (model-token sweep), FT160 (a `bench` subcommand for
  the full run), FT161, and a third face on FT156.

## Next command

`/bench-write-spec` for FT154.

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
