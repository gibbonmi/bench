# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` at `d5915574`
Spec: none staged; `specs/` is empty
Gate: green on `d5915574` (2026-08-23)

## State

`/bench-drain` reconciled the board against the tree. FT113 shipped in full.
Its path counting, its one author, and its named remainder all landed under
`148f3a68`/`bc3e9377`/`d5915574`. Its `ROADMAP.md` row and `roadmap/FT113.md`
are retired.

The one parked idea and the `landing-authors-the-flip` retro are both
drained. The idea named a `sanitize.Controls`/`ShellQuote` ordering defect in
`worktree.landedIncomplete`'s `next=` field. It merged onto FT169.

The retro's recommendations merged as `Occurrence:` evidence on four rows.
FT169 gained the paste defect and an unnamed moved-destination refusal. FT238
gained the missing worktree-rebase verb. FT177 gained the post-landing
stale-`dist/bench` gap. FT89 gained `bench idea` writing to `main`'s checkout
when run outside a worktree.

One recommendation — the `craft-delegate` worktree-exec rule — was already
documented, so it was dismissed rather than merged. The inbox and the retro
directory are both empty. `capture/learnings.md` had no open entries.

The `## Recommended sequence` is refreshed to drop the shipped FT113 line and
add FT224 third. `FT166`'s recommended-dependency row dropped its shipped
`FT113` half.

Two files stay uncommitted outside this drain's scope: `decisions/worktree-test-latency.md`
and `decisions/assets/worktree-test-invocation-census.md`, from an unrelated,
still-open `/bench-shape-idea` decision map. This drain did not touch them.

## Next command

`/bench-implement-spec`

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
