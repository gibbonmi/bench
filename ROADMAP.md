# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

## Features, in priority order

**FT44 — salvage branches sweep automatically once their content lands.**
`bench status` perpetually recommends `bench worktree clean` for unmerged
`worktree-*` salvage branches the sweep deliberately keeps, so the
recommended action provably no-ops. Reviewer direction (2026-07-08): the user
must never have to investigate — the sweep's proof of "landed" widens beyond
ancestry so a salvage branch whose changes are already contained in the
default branch (the common case: a delegate draft superseded by the merged
version) is deleted automatically; only genuinely un-landed content is kept,
and that keep gets an honest status action.
Spec: `specs/defect-batch-ft43-49.md`. Next: `/bench-implement-spec`.

**FT45 — worktree lease reclaim race.** Two crash-recovery reclaimers of one
dead-pid lease can both win the same worktree: `Claim`'s takeover rename is
blind to lease identity, falsifying the code's own "cannot both win" comment;
the concurrent-acquire contract covers fresh-mint only. Identity-verified
takeover plus a two-reclaimer stress case.
Spec: `specs/defect-batch-ft43-49.md`. Next: `/bench-implement-spec`.

**FT46 — ADR 0002 posture 5 amendment.** The record asserts nothing
short-circuits a gate run on a cache hit; `bench commit`'s verdict reuse now
does exactly that (soundly — exact tree-hash key, fresh-only). Amend the
posture to record the decision, and pin the capture-only-allowlist reuse
regression while there. Spec: `specs/defect-batch-ft43-49.md`.
Next: `/bench-implement-spec`.

**FT47 — unlink leaves a dangling CLAUDE.md.** Link-created CLAUDE.md is
written outside the install plan, so it is never manifest-recorded and
`bench unlink` leaves it importing just-deleted files. Record it in the
manifest when link created it; add it to the README leave-behind list.
Spec: `specs/defect-batch-ft43-49.md`. Next: `/bench-implement-spec`.

**FT48 — CHANGELOG append duty has no owner.** Both 2026-07-07
learnings-sourced promotions are missing their mandated entries, and no gate
check or phase step owns the duty — the `/bench-update-kit` baseline drifts
silently. Backfill the two entries and anchor the duty (a `/bench-what-next`
drain-checklist step, or a conformance check if drift recurs).
Spec: `specs/defect-batch-ft43-49.md`. Next: `/bench-implement-spec`.

**FT49 — pre-push guards a fabricated default branch.** `git.DefaultBranch`
falls back to `main` when `origin/HEAD` is unresolvable and link bakes that
answer into the hook, so a `master`/`trunk` repo linked before its remote
exists gets a backstop that never fires. Resolve at push time, or warn at
link and add a doctor row comparing baked vs live default.
Spec: `specs/defect-batch-ft43-49.md`. Next: `/bench-implement-spec`.

**FT38 — dashboard visual identity pass.** `bench dashboard` v1 shipped
data-faithful and visually neutral; the original idea wanted an
ui_examples-inspired rich treatment with animated characters. Taste is a
reviewer call, so the work starts as a grill, not a build. Decision detail is
recoverable via `bench spec history dashboard`. Next: `/bench-shape-idea`.

**FT50 — one-source collapse batch.** Export the pre-push marker const to
`bench guards` instead of its copied literal; a shared const for the
`bench-last-gate` cache filename (writer and reader currently copy it); a
sync test for the git guard's deliberately-inlined wrapper search order;
delete the stale `default_branch` shell-mirror comment. Next:
`/bench-write-spec`.

**FT51 — CLI hygiene batch.** Unknown subcommand prints help at exit 0 (typo
indistinguishable from success — dispatcher is the outlier against the exit-2
norm); `--version`/`--help` fall into the same case; `bench canary` passes
silently; `coverage`/`link` error-message nits; decide and record the
harness-form posture for CLI strings that print `/bench-*` phases. Next:
`/bench-write-spec`.

**FT52 — docs batch.** CONTEXT.md names 6 of the 10 live status signals and
has no canonical term for the `bench dashboard` artifact (feeds the FT38
grill); README's `internal/` layout omits `dashboard/` and `outline/`; the
lighter-path threshold is worded three ways across BENCH.md and two commands;
decide whether `projects/benchkit.md` (the internal dogfood profile) keeps
shipping; disposition `research/unit_testing.pdf`. Next: `/bench-write-spec`.

**FT53 — test/hardening batch.** Pre-push read-loop newline-tail guard;
symlink-loop cap in the wrapper's path resolver; propagate the
`RegisteredWorktrees` classify error (last false-empty instance); fix the
concurrent-acquire overlap comment overstating its barrier; a status signal
for a missing pre-push hook. Next: `/bench-write-spec`.

**FT54 — assessment owner.** Third consecutive platform assessment run with
no owner: each run re-derives the method (verify last drain, area sweeps,
adversarial synthesis, ranked backlog) from the prior file. A `craft-assess`
skill or `/bench-assess` phase codifies the drill. Next: `/bench-write-spec`.

**FT55 — split-vs-grant rule for structure budgets.** A file-length split is
only free when the dir has file-count headroom; the conformance split traded
`FILE TOO LONG` for `DIR CROWDED` (resolved by a reviewer grant). Add one line
to `craft-seams` and the implement-spec structure-housekeeping note: check both
budgets before choosing split-vs-grant. XS kit edit under `craft-synthesis`.
Next: direct kit edit, gated as usual.

**FT56 — sanctioned-command contracts visible at the point of use.** Three
contract-guesses in one day: a pathless `bench commit -m` (exit 2), `--spec`
misread as mere association when it flips the spec to `Status: implemented`,
and `bench spec history` not recalled during a shipped-row check. Amend the
BENCH.md CLI-inventory work-execution line with the one-clause commit contract
(`bench commit -m <msg> <path>...`, path-scoped, stages its own paths; `--spec
<slug>` marks the spec implemented — implementation green commit only), extend
the `--spec` usage string to state the flip, and name `bench spec history` as
the tool in the what-next reconcile step. XS kit edit under `craft-synthesis`.
Next: direct kit edit, gated as usual.

**FT57 — shared-worktree path pinning in craft-delegate.** During the FT43
stacked fix both sides missed that `cd` governs Bash CWD only: the delegate's
file tools used repo-root absolute paths and wrote into the main tree, and the
orchestrator's stale worktree CWD made a correct `bench commit` "nothing to
commit" get misfiled as a tool bug. Add two clauses to `craft-delegate`'s
Isolation section: a charge that shares an existing worktree pins all
file-tool paths to the worktree root, and a "nothing to commit" against a
visibly-modified file reads as a CWD/tree mismatch before it reads as a
defect. XS kit edit under `craft-synthesis`. Next: direct kit edit, gated as
usual.

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check, and a per-anchor
bite-proof meta-test (canaries prove one needle per family today; graduate on
observed anchor rot). `bench symbols` is not carried; restore only if agents
demonstrably burn turns on symbol search.

**FT24 (parked pending upstream) — Codex agent-line guard parity.** Researched
2026-07-07: not implementable on current Codex — delegation never surfaces as a
matchable `tool_name` on a deny-capable hook event, and `SubagentStart` neither
carries the delegate's resolved model nor honors a deny (verdict recorded in
`.bench/BENCH-reference.md` Hook Layers). Graduate only when the Codex
changelog adds a spawn tool name or a deny-capable SubagentStart.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

## Recommended sequence

1. `/bench-implement-spec` — FT44–FT49 defect batch
   (`specs/defect-batch-ft43-49.md`, staged; FT43 shipped, FT44 next).
   Fresh mid-tier session.
2. `/bench-shape-idea` — FT38 dashboard visual identity: pure reviewer taste,
   grill before build.
3. `/bench-write-spec` — FT50 one-source collapse batch: mechanical,
   well-scoped, no open decisions.
