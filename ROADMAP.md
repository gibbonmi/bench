# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

## Features, in priority order

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

**FT59 — sentinel precondition for fix-pass delegations.** A fix-pass
delegate building on a repo snapshot verifies a commit-specific sentinel (a
function or test the commit under fix added) before working, so a stale
snapshot fails fast instead of producing a divergent rebuild. One clause in
`craft-delegate`'s charge template. XS kit edit under `craft-synthesis`. Next:
direct kit edit, gated as usual.

**FT62 — structure debt: `internal/contract/helper.go` at 441/400.** Grown by
FT53's RunAtWithTimeout helper; the one open `bench structure` issue. Split
along responsibility or propose a reviewer grant — check the dir's file-count
headroom before choosing (the FT55 rule; `internal/contract/runtime/` already
holds a dir grant, so a split into that dir deepens dir debt). Next:
`/bench-implement-spec` (lighter path).

**FT63 — staged-spec posture for new phase commands.** The stale-command-
reference sweep reds a `Status: staged` spec that names its own not-yet-built
`/bench-<new>` command (hit during FT54; worked around by landing spec and
command in one diff). Record the posture in `/bench-write-spec`: a spec whose
deliverable is a new phase command lands in the same diff as the command —
recommended over teaching the sweep a staged-spec exemption, which would weaken
the gate for every other stale ref. XS kit edit under `craft-synthesis`;
alternative (sweep exemption) is the veto surface. Next: direct kit edit,
gated as usual.

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check, and a per-anchor
bite-proof meta-test (canaries prove one needle per family today; graduate on
observed anchor rot). `bench symbols` is not carried; restore only if agents
demonstrably burn turns on symbol search.

**FT58 (parked pending evidence) — reclaim-lock protocol for lease Claim.**
Serialize dead-pid takeovers via a crash-safe lock file to close the residual
B-vs-C double-use window the FT45 no-clobber restore leaves open. Graduate on
an observed double-win (the two-reclaimer stress case red, or a real
double-use in the pool); until then the shipped identity-verify plus restore
is the accepted posture.

**FT24 (parked pending upstream) — Codex agent-line guard parity.** Researched
2026-07-07: not implementable on current Codex — delegation never surfaces as a
matchable `tool_name` on a deny-capable hook event, and `SubagentStart` neither
carries the delegate's resolved model nor honors a deny (verdict recorded in
`.bench/BENCH-reference.md` Hook Layers). Graduate only when the Codex
changelog adds a spawn tool name or a deny-capable SubagentStart.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

**FT38 (tabled, revisit on or after 2026-08-09) — dashboard visual identity
pass.** `bench dashboard` v1 shipped data-faithful and visually neutral; the
original idea wanted a rich treatment with animated characters, reference
saved at `ui_example/` (Gather-style pixel office with activity feed).
Reviewer tabled it 2026-07-09 for at least a month. When it revives, the work
starts as a grill (`/bench-shape-idea`); decision detail recoverable via
`bench spec history dashboard`.

## Recommended sequence

1. `/bench-implement-spec` — the XS kit-edit batch (FT55, FT56, FT57, FT59,
   FT63): five prose-only rule edits under `craft-synthesis`, one gated pass.
2. `/bench-implement-spec` — FT62 helper.go split-or-grant, applying the FT55
   rule the batch just landed.
