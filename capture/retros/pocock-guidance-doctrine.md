# Retro — pocock-guidance-doctrine

Landed 2026-08-12 at `a354cede`: 11 implementation tickets plus one review
round (16 raw findings → 8 repair commits), every landing gate-green through
path-scoped `bench commit`, final `--spec` flip on the last commit.

## What worked

- Serial tickets with per-ticket worktrees and independent coordinator probes
  caught real defects delegates missed: three vacuous sed probes were exposed
  by re-checking that a mutation actually landed before trusting its verdict.
- The pre-prune worktree route satisfied PG15's red-before-prune sequence
  without a spec edit: the checker observed all seven over-budget subjects on
  the real tree, then rebased and landed last.
- The cross-harness reviews earned their cost: sol's breakdown REJECT found
  the out-of-fence canary family and the mirror-classifier refusal before any
  code; terra's falsification pass found the symlinkable skills root and the
  spec-build residue. Terra's REFUTED verdict was partly a spec-time versus
  ticket-time fence conflation — verify a cross-harness verdict's premises
  before acting on its severity.
- The shortfall predicate landed by this very build fired during its own
  review repair: the AF1 delegate measured 190 lines against the 150 budget
  and exited with evidence instead of deleting doctrine; the budget moved to
  175 as a reviewer decision.

## What to change

- Line-count budgets invite wrap-width gaming; the AF1 episode is parked in
  `capture/IDEAS.md` (a word-count or house-wrap-normalized measure would
  close the loophole, and a further BENCH.md trim needs a doctrine decision).
- Guidance prose wraps mid-phrase, so single-line sed probes and raw-literal
  fixtures both need wrap awareness; three probes went vacuous this run
  before verification caught them (learnings entry filed).
- The `BENCH_CONFORMANCE_ROOT` probe surface evaluates families the real gate
  scopes differently and emits ~8 pre-existing artifact diagnostics; every
  later charge had to carry the ignore-list. Worth a dedicated clean probe
  entry point someday.
