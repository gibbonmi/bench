# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean, one commit ahead of `origin/main` before this drain's own commit.
Spec: none staged.
Gate: green at `f63384b3`.

## State

The drain is closed. `capture/IDEAS.md`, `capture/learnings.md`, and
`capture/retros/` are all empty, and the board carries the evidence instead. Four
rows took occurrences from the `light-path-wrapper-home-guard` retro and its two
journal entries: FT113, FT243, FT215, and FT106.

FT113 rose from LOW to MEDIUM. `bench commit --spec` exits 1 on
`landed-but-checkout-incomplete` for a tickets-only light-path spec, after a
fully correct green commit, and leaves `specs/<slug>/` as untracked residue. The
drain reproduced that exit through the accused command on its own implement-now
commit, so the row now has two same-day occurrences of the face. The fix still
needs the reviewer's call on retry-versus-report. Expect the residue and remove
it by hand after any `--spec` landing until the row ships.

The retro's `craft-seams` recommendation shipped in this session as an
implement-now light-path item at `f63384b3`: a registry-owning package carries an
admission rule that a test's author reads before placing a test there. It was
implemented inline rather than through a write-delegate, which follows the
right-size table and departs from this phase's own delegate wording — flagged for
veto.

## Next command

`/bench-write-spec` — FT243: a capture entry the parser cannot see is reported as
zero, not as a failure. Top-ranked for a second drain running, and the face
repeated during this one.

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
