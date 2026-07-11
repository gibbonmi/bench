# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

## Features, in priority order

**FT68 — structured Bench phase conversation.** Build the staged spec
`specs/structured-phase-conversation.md`: the two proportional conversation
patterns (**Status:**/**Next:** for in-progress, `## Result`/`## Details`/`##
Next` for completed phases) land once in the shared communication rules, with
conformance and canary coverage. Next: `/bench-implement-spec`.

**FT62 — structure debt: `internal/contract/helper.go` at 441/400.** Grown by
FT53's RunAtWithTimeout helper. Split along responsibility or propose a
reviewer grant — check the dir's file-count headroom before choosing
(`internal/contract/runtime/` already holds a dir grant, so a split into that
dir deepens dir debt). Next: `/bench-implement-spec` (lighter path).

**FT64 — structure debt: `internal/canary/canary_test.go` at 407/400.** The
targeted canary-phase work grew this past the file budget after FT62 was
recorded. Split along responsibility or propose a reviewer grant, checking
both file and directory budgets before choosing. Next: `/bench-implement-spec`
(lighter path).

**FT69 — structure debt: `internal/contract/surface/binary_repair_test.go` at
416/400.** Split along responsibility or propose a reviewer grant, checking
both file and directory budgets before choosing. Next: `/bench-implement-spec`
(lighter path).

**FT70 — structure debt: `internal/worktree/lifecycle_test.go` at 405/400.**
Split along responsibility or propose a reviewer grant, checking both file and
directory budgets before choosing. Next: `/bench-implement-spec` (lighter
path).

**FT72 — structure debt: `internal/git/git.go` at 403/400.** Split along
responsibility or propose a reviewer grant, checking both file and directory
budgets before choosing. Next: `/bench-implement-spec` (lighter path).

**FT73 — structure debt: `internal/worktree/worktree_test.go` at 489/400.**
Split along responsibility or propose a reviewer grant, checking both file and
directory budgets before choosing. Next: `/bench-implement-spec` (lighter
path).

**FT74 — structure debt: `internal/intent/intent.go` at 406/400.** Split along
responsibility or propose a reviewer grant, checking both file and directory
budgets before choosing. Next: `/bench-implement-spec` (lighter path).

**FT71 (parked pending evidence) — kit-sized shift-session log.** Durable
local evidence per shift run: agent identity, gate verdicts, commits. Cousin
of the compliance-assessment H-03 audit-trail finding, which stays out of kit
scope. Graduate on an observed need for durable run evidence — a disputed gate
verdict or an audit ask against a real shift.

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
2026-07-11: still not implementable on current Codex — delegation has no
matchable `tool_name` on a deny-capable hook event, and `SubagentStart` cannot
stop the subagent. The current surface verdict is canonical in
`.bench/BENCH-reference.md` Hook Layers. Graduate only when the Codex changelog
adds a spawn tool name or a deny-capable SubagentStart.

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

1. `/bench-implement-spec` — FT68 structured phase conversation
   (`specs/structured-phase-conversation.md`).
2. `/bench-implement-spec` — FT62/FT64/FT69/FT70/FT72/FT73/FT74 structure
   split-or-grant passes (lighter path, one file per pass).
