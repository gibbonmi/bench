# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `a0c9d26e`, working tree clean apart from this file, 27 unpushed commits
Spec: none active — FT198 shaped, awaiting `/bench-write-spec`
Gate: green at `a0c9d26e`

## State

`/bench-shape-idea 198` ran a single-session grill (rounds 1–3, all answered,
reviewer confirmed shared understanding). No decision map was created: the tree had
no multi-session dependency. The reviewed decisions, which stay closed:

- `ROADMAP.md` becomes a hand-ordered index: section headings, board prose
  (release-readiness, dependencies, recommended sequence, section intros) and, per
  row, exactly today's heading line `**FT<n> (<priority>[, decision required]) — <title>.**`
  with no body.
- Each row's body — prose, `Occurrence:` ledger, `Sources:` — lives in
  `roadmap/FT<n>.md`, which opens with the same heading line repeated verbatim.
- The CLI parses only that split shape (no legacy inline-body compatibility: this
  is the only repo with a roadmap). It exits non-zero on any of: an index row with no
  file, a file with no index row, an inline body under an index row, a file heading
  that differs from its index line. `.bench/gate.sh` invokes that check so a lost
  detail owner is a red.
- Migration is a one-shot script, run once in the build then deleted; acceptance is
  `bench roadmap --context --full` rows, ledgers, and sequence identical before and
  after.
- One spec covers layout, parser, gate check, and this board's migration.

Spec-writer discretion (bounded, non-behavioral): whether the check rides the
schema-4 parse or a `--check` flag; which `ROADMAP.md` readers move to the row-file
parse (`bench idea --owner` validation, status spec-path reconcile, `--context --row`
fetch); the migration script's language.

## Next command

`/bench-write-spec` FT198 (progressive roadmap) — from this session's decisions above.

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
