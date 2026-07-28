# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `9977c27`, clean tree, 9 unpushed commits
Spec: none staged.
Gate: green at `0eaef37` — stale, work tree `ab832d6`

## State

- **FT152 is built, gate-green, and on `main` — unpushed.** `specs/implement-spec-full-run.md`
  is `Status: implemented`. The build is five commits: `3eb1c9a` (stories 1-9,
  11-15), `8218f92` (review repairs), `7f3f25b` (merge of FT155), `dbdcbda`
  (fixture refresh), `fa83f64` (status flip), plus `d09d081` for the captures.
  Gate green on `main`. The push is the reviewer's.

- **What shipped.** Two shared rules in `.bench/BENCH.md` — fix-don't-park in
  Workflow, the outside-source warrant rule in How to talk to me — each pinned by
  a marker in `checkSharedRuleSingleSource` plus a section-scoped anchor, so the
  anchor means placement, not presence. `NEVER assume, always verify` is now
  anchored too; it was unguarded before. Three point-of-use verify hooks in
  `bench-shape-idea`, `bench-implement-spec`, `bench-review-implementation`. One
  bounded `## The `--full` run` section in `bench-implement-spec.md` with 20
  section-scoped requires and 10 paired forbids. Six canary fixtures.

- **The falsification pass changed the outcome, and that is the headline.** A
  three-axis review returned 29 findings and cleared the scope fence. A Codex
  `gpt-5.6-sol` xhigh pass charged to refute found that the fence re-derived
  fix-don't-park as path-width-based where `.bench/BENCH.md` sets it
  decision-based — the two gave opposite answers for a small out-of-story bug.
  Fixed in `8218f92`. Second such result after FT91; journaled, not yet a rule.

- **Spec retirement is deliberately NOT done.** `bench status` will show
  `1 merged spec awaiting retirement`. Retiring deletes the spec and its decision
  map, and the spec header carries the reviewer's open veto list — including
  calls I made on their blanket authorization (the `.bench/structure.budgets`
  grant for `docs_workflow_helpers_test.go 660`, and correcting the spec's
  `.claude/commands/` mirror decisions, which assumed a copied tree where the
  repo has a symlink). Retire after the reviewer has read that list, not before.

- **The 2026-07-28 drain landed at `9977c27`** — inbox and journal are empty.
  FT152's follow-up findings now live on the roadmap: the coverage row-identity
  gap as a clause on FT133, the anchor-mechanism weakness as FT156 (grill), the
  falsification-trigger question parked as FT158, the registration-seam charge
  rule as FT157. FT154 is marked unblocked and specs next after FT152's
  close-out.

## Next command

`/bench-final-check`

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
