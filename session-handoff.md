# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `0ebba2f`, clean except this file
Spec: none staged.
Gate: green at `f28e44d`; only capture-only paths written since (`ROADMAP.md`,
`IDEAS.md`, this file).

## State

- **FT86 is fully closed.** Built, three-axis reviewed (26 findings, all
  resolved across four seam-scoped slices), spec flipped to implemented and
  retired, decision map retired. Nothing FT86-shaped remains open except
  FT140's two residual reviewer calls (the 2 MiB learnings read bound —
  keep or reverse — and `TestAXILearningsWrongType`'s missing coverage
  provenance).
- **The 2026-07-26 drain was a reconcile + restructure pass** under the
  reviewer's batch approval. FT86's row removed as shipped; FT96 absorbed
  FT118, FT103's charge half, FT131's guidance arm, and the worktree-hook
  keying idea; FT98 absorbed FT127 and FT114 as one preserve-then-discard
  primitive; FT107 absorbed FT110 and FT119; FT113 absorbed FT121;
  FT123/FT124/FT125 grouped under a session-tax section; ten ideas drained
  into those merges and new rows FT133–FT140. `ROADMAP.md` parses clean
  (FT86's bold-leading paragraph was the malformed row); the row-grammar
  decision itself is FT132's, still the reviewer's.
- **Calls flagged for post-hoc veto** (all deviations from the pre-approved
  restructure plan, forced by the tree having moved past its pins):
  FT86's row removed rather than kept (the plan assumed the review had not
  run); FT132 narrowed because its conflation face shipped; FT130 excluded
  from the FT107 prose batch (its preferred fix is CLI-mechanical); the
  capability-skip visibility idea merged into FT124's reader rather than a
  new row; the recommended sequence no longer leads with the FT86 review.
- **Both capture sources are at zero.** `IDEAS.md` empty, `.bench/learnings.md`
  holds no open entries.
- Known ambient facts: `bench status` may report the gate stale after this
  handoff write (`session-handoff.md` is not on the capture-only allowlist —
  FT113's gate face); the FT91 conformance-phase long pole and the structure
  budget violations are unchanged.

## Next command

`/bench-write-spec` — FT96, the batched delegation discipline (top of the
roadmap's recommended sequence).

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
