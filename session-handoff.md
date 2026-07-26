# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `33724b4`, 8 dirty paths, 2 unpushed commits
Spec: none staged.
Gate: green at `edc15e9` — stale, work tree `5b4a4e6`

## State

- **FT136 is built, reviewed, and retired.** The slice-fence rule lives once in
  `craft-spec` ("Slicing a build for delegates": fence rule,
  shared-primitives-first, tier-independence); `craft-review`'s Standards axis
  hunts fence-boundary duplication; `craft-delegate` carries the content-free
  charge-time pointer; frontmatter/index regenerated; the section heading is
  pinned in the conformance anchor registry (observed red-then-green). Landed
  gate-green as `b381816`, spec retired in `33724b4`.
- **Review disposition that stays closed:** the tier-independence clause and
  the section placement were challenged by the Standards axis and rejected as
  spec-closed; the pair captions were reworded to timeless register; the
  anchor-registry rows upgraded the spec's "not TDD-able" anchor edge into
  gate coverage — flagged for post-hoc veto, two lines in
  `internal/conformance/docs_workflow_helpers_test.go`. Note: the
  craft-delegate pointer sentence is anchor-pinned raw — rewrapping it mid-phrase
  reds the gate by design.
- **Open map tickets, untouched:** `decisions/cost-follows-project-size.md`
  stays (backs FT91/FT101): #2 (time the conformance checks) then #3 (FT91
  go/no-go); #6 waits opportunistically for a genuinely seam-shaped slice.
- **The unpushed commits are ready; the push is the reviewer's.** Nothing in
  the tree waits on it.
- Known ambient facts, unchanged: 17 worktrees remain from earlier sessions and
  were left untouched (the FT136 one was cleaned); the structure budget
  violations and the conformance-phase long pole stand where they were.

## Next command

`git push origin main` — then `/bench-shape-idea`, the board's leading
invocable signal (`decisions`).

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
