# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `2de69aa`, checkout clean, 13 unpushed commits.
`bench status` still reports 15 dirty paths; they are the abandoned worktrees
below, not this checkout.
Spec: `specs/ft86-fail-closed-control-records.md` (Status: staged)
Gate: green at `32fdf17` — current

## State

- **FT86 is BUILT and landed — all 20 stories, all 34 coverage rows green.** Seven
  code commits (`96faeba` classifier · `f5051fd` coverage · `b1e2778` git default
  branch · `72d787c` learnings · `b987e56` maps · `5697ec3` roadmap · `f9cc907`
  status), each on its own green whole-tree gate, plus a confirming gate on the
  integrated tree. The spec's `Status:` line is still `staged` — the transition to
  `implemented` is `/bench-final-check`'s, not the build's.
- **Semantic review has NOT run.** That is the next phase. Nothing has been pushed;
  the reviewer owns the merge.
- **Two abandoned worktrees still hold discarded work.** `main` is clean; two
  round-2 delegates cross-contaminated each other through the stash (see below),
  and their worktrees still hold the discarded mixed diff. They can be cleared
  with `bench worktree clean <path>` — left for the reviewer because discarding
  is hard to reverse.
- **Four calls are open for post-hoc veto.** (1) `bench maps` now exits 1 when
  `decisions/` holds a file with no ticket heading; all three real decision files
  have one, so this repo is unaffected, but a *linked* repo with a `README.md`
  there would get a permanent exit 1 — a pre-existing assertion that such a file
  was silently skipped was updated to match. (2) `learnings`' `unsupported-schema`
  predicate was narrowed to "zero `## ` headings at all" rather than the spec's
  literal "no dated heading", because the literal rule fails-closes every fresh
  repo whose seeded template has a placeholder heading. (3) `bench dashboard` still
  degrades a failed learnings read to `0` — story 17 names `bench status` only, so
  the scope was held, and the choice is recorded in a comment at the call site.
  (4) The three calls the spec itself flagged (outline's exclusion, declining the
  `traversal` state, `git.Facts`'s `DefaultResolved` shape) remain open.
- **The spec's coverage map names story 16's guard `TestAXIOutlineSymlinkSkipped`.**
  The real guard is a subtest, `TestAXIOutlineContracts/AXI_outline_tracked_symlink_skipped`.
  It is green and asserts the `link.go,nonregular` row; only the map's symbol name
  is imprecise.
- **Never run `git stash` in a delegate charge.** The stash ref is repository-global,
  not per-worktree; two concurrent delegates in separate worktrees cross-applied each
  other's changes through it and both diffs became unattributable. Use
  `cp` aside plus `git show HEAD:<path> > <path>` to prove a red. Captured in
  `.bench/learnings.md`.
- **Rebuild `dist/bench` before any suite that drives it.** The AXI and runtime
  contract tests execute the built binary, not the package source; a stale binary
  produced false reds twice this session and can equally make broken code look
  green. Captured in `.bench/learnings.md`.
- **Both capture sources are drained to zero** (`8e4ded4`). The stash hazard was
  merged into FT96 and the stale-binary trap opened FT131; the parked
  empty-vs-absent `ROADMAP.md` conflation opened FT132.
- **Never mutate the repository while a gate runs**, and build `dist/bench` only with
  `scripts/go-build.sh`. `projects/benchkit.md`'s cold-session notes carry these and
  the `internal/canary` nested-run trap.
- `bench structure` reports 18 issues; 17 are pre-existing, and the new one is the
  `internal/conformance/` file added by the git-default-branch slice. That
  directory's structure grant still records "sits one over a full budget", which now
  understates by one file — reviewer-owned prose, deliberately not rewritten.

## Next command

`/bench-review-implementation — specs/ft86-fail-closed-control-records.md`

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
