# Retro: progressive-roadmap (FT198)

## Outcome

Landed via `bench worktree land` at `5ee78026`, spec-retired at `595763a1`.
`ROADMAP.md` is now a body-less index (67 rows, one physical heading line
each); each row's body, occurrence ledger, and `Sources:` line live in
`roadmap/FT<n>.md`. The `roadmap-detail-integrity` conformance check (Dev
tier, 9 canary fixtures) grades the split shape; `bench roadmap`, `--context`,
`bench status`, the dashboard, and `bench idea --owner` all read through one
loader. Migration differential (base vs. migrated `roadmap_rows`/`sequence`
blocks) was empty, independently re-derived by the coordinator rather than
trusted from the delegate's report. Nine tickets total: five build tickets,
four repair tickets from a three-axis review (19 raw findings, 12
de-duplicated repair targets, all accepted findings closed).

## Gate-stage timings

No single wall-clock total was captured cleanly — every `bench gate --fresh`
run in this build was invoked through the worktree-vs-kit binary-trust
workaround (see Agent-experience improvements), which adds its own variable
overhead. Sub-phase timings were stable across the ~15 full gate runs this
build required: `internal/worktree`'s test package dominates every run
(51–56s), `internal/publication` (~30s), `internal/diff` (~8.5s), and
`internal/intent`/`internal/gitguard` (~4–5s each) are the next-largest
fixed costs regardless of what changed. A change confined to
`internal/roadmap`/`internal/conformance` still pays this fixed floor because
`bench gate` runs the whole suite — no incremental/changed-package-only gate
path exists yet.

## Ticket-versus-spec-slice and delegate performance

All nine charges were ticket-sized (one ticket, one delegate, one worktree);
none were handed a whole spec slice. Ticket sizes varied widely in practice
despite uniform charging: ticket 1 (parser + this repo's own 67-row
migration) ran 83 tool calls / ~28 minutes wall-clock, the largest single
charge in this build and roughly 3x the next-largest (repair D, 70 tool
calls / ~16 minutes). The smaller CLI-rendering and status-wiring tickets (3,
4) ran 60–80 tool calls each at Sonnet/low despite being scoped tighter than
ticket 1 — Sonnet's charges took comparably long wall-clock to Opus's on
similar-complexity work, but every one first-pass accepted, so the tier
choice held.

## Coordinator catches

- **Ticket 1's delegate correctly overrode the coordinator's charge.** The
  charge paraphrased the per-class fault disposition wrong (said inline-body
  and heading-mismatch drop the row); the delegate re-read the spec, found
  the charge and spec disagreed, and implemented the spec's actual rule —
  flagging the discrepancy rather than silently picking one. Caught by the
  coordinator only in the sense of verifying the delegate's citation against
  `spec.md:125` before accepting it; the delegate did the real catching.
- **An under-declared ownership fence stopped `bench preflight review`
  cold.** Ticket 1's migration necessarily touched
  `internal/conformance/recurrence_maintenance_contract_test.go` (its ledger
  fixtures parsed `ROADMAP.md` bodies directly, pre-migration) — a file
  neither the ticket's own `Writes:` line nor the spec's Ownership-fences
  list named. Two small administrative commits (inline-allowance edits, not
  delegated) corrected both before review could proceed.
- **Porting repair ticket A onto the retained source produced a real merge
  conflict**, not a clean patch apply: `internal/status/status_test.go`
  conflicted with repair C's already-landed removal of a test ticket A's base
  predated, and `internal/roadmap/context_test.go` carried a `[2]string`
  literal for `writeBoard` that repair B's later-landed typed-diagnostic test
  addition didn't know had become a `Row{}` struct. Both were coordinator
  merge-conflict resolutions, not delegate defects — a direct consequence of
  running repair tickets A/B/C in parallel from a shared pre-repair base
  rather than serially.
- **Repair D's delegate deviated from the ticket's suggested implementation**
  for the wrapped-heading double-diagnostic fix (C3): the ticket suggested
  marking the shared `indexed` map before the wrapped-heading `continue`;
  the delegate used two separate maps (`rowed` for duplicate detection,
  `claimed` for the orphan check) instead, because the ticket's suggested
  fix would have silently changed duplicate-ID behavior too. Verified
  independently with a fresh mutation probe (a different site than the
  delegate's own) before accepting.

## Repair attribution

| ticket | repair rounds | cause |
| --- | --- | --- |
| split-the-board-parser-and-migration-in-one-green | 0 | none |
| register-the-roadmap-detail-integrity-check | 0 | none |
| render-the-split-tree-through-the-roadmap-commands | 0 | none |
| read-row-files-from-status-dashboard-and-the-stripped-journey | 0 | none |
| document-the-split-board | 0 | none |
| repair-fixture-harness-dedup-and-comments | 1 | tree-drift |
| repair-typed-diagnostic | 0 | none |
| repair-dashboard-pin-location | 0 | none |
| repair-roadmap-detail-integrity-correctness | 0 | none |

## Agent-experience improvements

### Bench CLI

- **The dev-vs-installed-binary trust mismatch cost real time across this
  entire build.** The globally installed `bench` (0.2.0) predates this
  repo's dev source; every gate/commit/preflight call needed a hand-built,
  hand-sealed `dist/bench` plus `BENCH_RUN_BINARY`, and `bench worktree
  land`'s internal gate call additionally refuses an *inherited* binary for
  a prospective (never-before-graded) composed tree — discovered only by
  reading `internal/runbinary`/`internal/freshness` source directly, since
  the error text ("seal source digest does not match current build inputs")
  names the wrong root (the kit tree, not the worktree under grade) when the
  binary was built at the worktree instead of the kit root. A `bench doctor`
  or `bench repair` diagnostic that names this exact mismatch and its fix
  (`bash scripts/go-build.sh <kit-root> dist/bench`, no `BENCH_RUN_BINARY`
  needed for `land`) would save every future dev-build session the same
  investigation.
- **No changed-package-scoped gate path.** A one-file diff still pays the
  full ~52s `internal/worktree` test floor. Not this build's fix, but the
  fixed cost was felt on every one of the ~15 full gate runs this build
  required.

### Skills

- `bench-craft-tickets`' `Writes:` fence should be checked against the
  gate's own docs-currency/ledger-migration test fixtures when a ticket's
  work is a migration — this build's ticket 1 needed a collateral fixture
  fix in a file no story or ownership-fence line named, caught only at
  `bench preflight review` rather than at ticket-slicing time.

### Process

- Running repair tickets in parallel from one shared base, then porting
  each onto the retained source serially, is real leverage (three tickets
  landed in the time of roughly one and a half sequential ones) but shifts
  merge-conflict cost onto the coordinator's port step instead of a
  delegate's own worktree. Worth naming explicitly in `craft-delegate`: a
  parallel repair batch's tickets should declare which files they expect to
  touch in common, so the coordinator can order the ports to minimize
  conflict surface rather than discovering it at apply time.
