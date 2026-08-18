# Proposed roadmap — active work only (twelve items, same portfolio as `action-items.yaml`)

Target architecture: the current Bench, incrementally — the oracle regains sight of the
kit's own contracts, one verdict reader, one front door projected from existing status
signals, the debug discipline restored and reachable on both harnesses, adoption that
goes green, and a measurement harness before any prose is cut. No spine, no new state
layer, no compiler, no claim graph (A12 records that decision). Demonstrated behavior
preserved: tree-keyed gate evidence, prospective commits, ownership fences and
preflight, coverage `--check`, guards, `bench setup`, pooled worktrees, `/bench-debug`,
`craft-review`'s citation standard, `craft-delegate`'s probe rule, the retro's repair
table.

Ordered by dependency and leverage. **Commitments** are items whose acceptance criteria
are finite and mechanically checkable; **experiments** produce a measurement and a
decision, not a feature. This does not modify the live `ROADMAP.md`; it is what the next
`/bench-what-next` drain should turn the board into. Existing rows fold in as noted
(full mapping in `roadmap-dispositions.yaml`).

## Immediate next ticket

**A1 — the dev gate grades the kit's own contracts** (`next-ticket.md`). Blocked by
nothing; unblocks A3, A5, A8, A9, A10.

## Active items

| Rank | Id | Priority | Kind | Title | Depends on | Folds in |
|---|---|---|---|---|---|---|
| 1 | A1 | P0 | commitment | Live-root conformance runs in the dev gate; environment skip inside the oracle is red; skips named | — | FT133 (skip occurrences), FT120 (skip occurrence), FT213 (skip occurrence), FT89 (doc-currency half), the invalid map |
| 2 | A2 | P0 | commitment | One staleness rule for verdicts; status/handoff agree with `bench gate` | — | FT162 (handoff-state half) |
| 3 | A6 | P1 | commitment | Kit tests stop writing into the operator's real `BENCH_HOME` | — | (Sol E006) |
| 4 | A4 | P1 | commitment | Adoption smoke — scaffolded gate goes green | — | — |
| 5 | A3 | P1 | commitment | `/bench` front door over `bench status --route`; staged/setup signals; invocable actions; what-next → drain | A1, A2 | FT89 (root/inventory generation), FT180 (route decision surface) |
| 6 | A5 | P1 | commitment | Restore `/bench-debug` constraints; Codex trigger decision; Claude parity check; drop inert key | A1 | FT112, FT24 (decision record) |
| 7 | A7 | P1 | commitment | Release workflow publishes only through `bench release submit/promote` | — (before next release) | ASSESSMENT.md §6; release-readiness item 4 |
| 8 | A8 | P2 | commitment | Minimal work-state extension: constrained `## State`, `Repro:` line, Next from the router | A2, A3 | FT162 (handoff half) |
| 9 | A9 | P2 | commitment | Hygiene batch: light-path close step + residue, stale prose/accept rows, outline bare form, guard rim + gaps, session-start hint, log pruning, dist/bench landing warning | A1, A3 | FT177 (cold-start half), FT180 (residue), FT201 |
| 10 | A10 | EXPERIMENT | experiment | Repair-loop tripwire — advisory row from gate records; measured against the retro repair table | A1, A5 | — |
| 11 | A11 | EXPERIMENT | experiment | Measurement harness (arms A/B/C, D/E/F, G) + instrumentation; pre-registered criteria | A1–A3 for arm C | FT170, FT138; **blocks FT100** |
| 12 | A12 | DELETE | decision | Do not build a work-state store, context compiler, claim graph, review runner, or strangler spine; revisit only on A11 evidence | A11 | (Sol T.1–T.3, R) |

## What leaves the active board

- **Done / archive:** FT212 (`--landed` works), FT190 (conformance check exists), FT165
  (craft-domain wired).
- **Not the next thing:** FT100 stays but moves from rank 1 to *after A11* — the row's own
  text and both audits agree cuts need a measurement first. FT207 and FT213 (ranks 2–3)
  are cheap decisions/edits with no fresh occurrence; they wait behind the P0/P1 pair.
- **Rewrite before working:** FT98, FT169, FT164, FT158, FT144, FT199/FT206, FT197,
  FT174, FT130, FT172, FT141 (name collision with the existing `bench gate pin`).
- **Deferred backlog (kept, not active):** everything else in `roadmap-dispositions.yaml`
  marked DEFER or FIX-WHEN-TOUCHED — batched by file (craft-delegate: FT205/FT213/FT221;
  craft-spec: FT192/FT209/FT214; comments: FT111/FT179; standards-debt: FT94/FT117) so
  each lands under one gate when its file is next opened.
- **Parked as before:** FT6, FT8, FT24 (Codex spawn hook), FT38.

## Ordering rationale (one line each)

1. A1 first because it is the only item that would have caught the others, it is small,
   and every later doc/adapter/phase change (A3, A5, A8, A9) needs the contracts running.
2. A2 next because a router (A3) and a resume artifact (A8) must not read a phantom red.
3. A6 and A4 are tiny, independent, and stop the kit from harming operators and adopters
   on every run.
4. A3 is the user-visible payoff and depends on A1/A2 only.
5. A5 protects the highest-evidence behavior and is cheap; A7 is release-boundary work
   that has no deadline until a release is attempted.
6. Experiments (A10, A11) run once the tree they measure is the post-A1..A9 tree for arm C;
   arms A/B/D/E/F/G can start immediately.
