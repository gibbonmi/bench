# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `9616919a`, working tree carries the reviewed `/bench-what-next` capture-drain batch, about to land as one commit
Spec: none staged; `specs/` is empty

## State

`/bench-what-next` reconciled the roadmap against the tree (nothing shipped
since the last drain; `specs/` and `spec_history` are both empty, so no
`bench spec retire` was owed this pass) and drained every capture source to
zero:

- `capture/IDEAS.md` was already empty.
- `capture/learnings.md`'s three open entries are verdicted: one (resolved
  `reviews/<slug>.md` blocking `bench worktree land`) merged as an occurrence
  into `roadmap/FT89.md`; two (refresh a ready `/bench-deepen` map before
  handoff; census shared decision readers before `/bench-write-spec` ticket
  slicing) became new kit-edit rows `roadmap/FT219.md` and
  `roadmap/FT220.md`.
- `capture/retros/worktree-cleanup-eligibility.md` is drained and deleted.
  Its actionable recommendations: three merged as occurrences into existing
  rows (`roadmap/FT113.md` — sharpen the `--spec` flip's `--help`/guard;
  `roadmap/FT169.md` — land's ownership-fence refusal should name the
  offending path; `roadmap/FT177.md` — landing guard should warn before
  removing a load-bearing `dist/bench`), one became a new kit-edit row
  (`roadmap/FT221.md` — promote `craft-delegate`'s cp-aside guidance to a
  named checklist step), and one became a new decision-required row
  (`roadmap/FT222.md` — standing per-repair-class delegate-tier preference in
  `projects/benchkit.md`). Repair-attribution tally from its table: 8 tickets,
  6 one-shots (0 repair rounds), 2 with repairs — cause counts:
  delegate-error 1, spec-row 1, other 2.
- `capture/agent-performance/claude-models.md`'s refresh from the prior
  `/bench-final-check` landing evidence (FT216, `worktree-cleanup-eligibility`)
  is carried through in this same commit.

`## Recommended sequence` rank 1 is now FT100 (`/bench-shape-idea` — cut prose
weight from `AGENTS.md`/`.bench/BENCH.md` and the craft-skill library to
demonstrated-delta clauses), reviewer-prioritized ahead of its
recommended-after-FT89 sequencing. FT207 (`/bench-shape-idea`) and FT213
(`/bench-write-spec`) follow at ranks 2 and 3, unchanged. This run's roadmap
edits are all feature/guidance-shaped except FT222, which is decision-only.

`dist/bench` must exist and be reasonably fresh for local `bench` CLI
resolution to work in this checkout; rebuild with
`go build -o dist/bench ./cmd/bench` if it's ever removed.

## Next command

`/bench-shape-idea`

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
