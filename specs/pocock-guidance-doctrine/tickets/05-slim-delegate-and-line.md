# Slim craft-delegate to its safety core and teach craft-line owned-red convergence

Blocked by: 04-slim-tickets-and-implement-spec.md
Ownership fence: `.agents/skills/bench-craft-delegate/SKILL.md`, `.agents/skills/bench-craft-line/SKILL.md`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors`, `.bench/BENCH-reference.md`
Integration surfaces: anchors→`internal/anchors/registry_data.go`; anchor canaries→`tests/canary/workflow-guidance-anchors`; line-routing conformance table→`.agents/skills/bench-craft-line/SKILL.md` stage-default rows; index currency→`.bench/BENCH-reference.md`
Contracts: the three stage-default table rows crossing `.agents/skills/bench-craft-line/SKILL.md`→`internal/anchors/registry_data.go` (AfterImplementSpec anchors), asserted by DL4 against the real docs-currency check
Closure: DL1/delegate-core-kept, DL2/prohibited-posture, DL3/line-owned-reds, DL4/anchors-and-conformance-green

## What to build

Rewrite `craft-delegate` to ≤120 lines keeping: fresh context with an explicit
line and bounded charge, worktree isolation for writes, read-only review
without a worktree, independent verification of every done-claim against the
tree and gate, and the capability-aware stop. Extend the stop to cover
reviewer-prohibited delegation: when delegation is unavailable *or* prohibited,
stop with one executable resume handoff (a spec-doc-only correction is not a
silent exception), and surface before spawn any delegation that changes who
performs requested work. Shed receipts, lifecycle plumbing, duplicated
mutation-probe ceremony, and charge boilerplate. Rewrite `craft-line` to ≤120
lines keeping the harness-local tier binding, the three starting signals plus
leverage override, the declaration, the cap, and the ladder — and add owned-red
classification: before any retry or escalation, classify reds into
diff-owned versus inherited-baseline versus spec-predicted against the pinned
inherited baseline; only diff-owned reds count, and a non-shrinking owned-red
set stops and surfaces a likely seam/spec contradiction instead of buying a
more expensive attempt. Preserve the three stage-default rows the
AfterImplementSpec anchors pin (`| Orchestration | mid + medium |`,
`| Ticket implementation | cheap + low |`, `| Review (axis or falsification) |
mid + high |`). Migrate the 23 delegate anchors and 3 line anchors with their
canary fixtures, retiring ceremony-only pins.

## Acceptance

- [ ] [DL1] (covers PG17) `craft-delegate` ≤120 lines keeps isolation, fresh charges, read-only review, and independent verification; receipts and lifecycle ceremony are gone.
- [ ] [DL2] (covers local) unavailable and prohibited postures both stop with one executable handoff; who-performs changes are surfaced before spawn.
- [ ] [DL3] (covers PG18) `craft-line` ≤120 lines classifies owned reds against the pinned baseline before retry/escalation and stops on a non-shrinking owned-red set.
- [ ] [DL4] (covers local) `go test ./internal/anchors ./internal/conformance` green, including the line-routing check over the preserved stage-default rows.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DL1/delegate-core-kept | delete the worktree-isolation rule | anchors check | remove it, run the docs-currency check, expect the owning anchor's red |
| DL2/prohibited-posture | drop the prohibited branch, keep only incapable | semantic review reread | reviewer-graded against PG17's three postures |
| DL3/line-owned-reds | escalate on any second red without classification | semantic review reread | reviewer-graded against PG18 |
| DL4/anchors-and-conformance-green | remove one stage-default table row | line-routing/anchors check | delete the row, run `go test ./internal/conformance`, expect the pinned-row red |
