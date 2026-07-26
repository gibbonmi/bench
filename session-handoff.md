# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `4085b4d`, main checkout clean, 13 unpushed commits
Spec: `specs/ft136-slice-fence-rule.md` — staged, approved 2026-07-26
Gate: green at `edc15e9`, stale by doc-only commits (map, handoff, spec).

## State

- **The FT136 spec is staged and approved.** Compiled off
  `decisions/cost-follows-project-size.md` (tickets #4/#5 closed): the fence
  rule and shared-primitives-first land once in `craft-spec`, `craft-review`'s
  Standards axis gains the fence-boundary duplication hunt, `craft-delegate`
  gets a content-free charge-time pointer. The top-tier falsification pass
  (reviewer-directed) returned RECOMMEND; its two findings are fixed in the
  committed draft.
- **Decisions that stay closed:** rule ships now, tier-independent — the
  cheap-tier retest (ticket #6) gates only mid-tier-by-default, never the rule;
  three surfaces, one source — do not re-litigate a single-skill home; the
  build runs as one write-delegation in one worktree at `gpt-5.6-sol` / high,
  and every build delegate for this spec is mid tier or higher (reviewer
  directive, 2026-07-26).
- **Open map tickets, untouched by this spec:** #2 (time the fifteen
  conformance checks) then #3 (FT91 go/no-go) for the gate wall-clock arm;
  #6 waits opportunistically for a genuinely seam-shaped slice.
- **The 13 commits are unpushed, and the push is the reviewer's.** Nothing in
  the tree waits on it.
- Known ambient facts, unchanged: 17 worktrees remain from earlier sessions and
  were left untouched; the structure budget violations and the conformance-phase
  long pole stand where they were.

## Next command

`/bench-implement-spec specs/ft136-slice-fence-rule.md` — in a fresh mid-tier
session, per the spec-phase exit contract: orchestration stays off top-tier
context while the spec's single write-delegate runs at its declared
`gpt-5.6-sol` / high line.

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
