# Rewrite review as re-derive-then-compare with dispositions and committed pickup

Blocked by: 05-slim-delegate-and-line.md
Ownership fence: `.agents/skills/bench-craft-review/SKILL.md`, `.agents/commands/bench-review-implementation.md`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors`, `.bench/BENCH-reference.md`
Integration surfaces: anchors→`internal/anchors/registry_data.go`; anchor canaries→`tests/canary/workflow-guidance-anchors`; index currency→`.bench/BENCH-reference.md`
Contracts: the finding-disposition vocabulary (`no-op`/`auto-fix`/`ask-user`) crossing `.agents/skills/bench-craft-review/SKILL.md`→`.agents/commands/bench-review-implementation.md`, asserted by RD3 against both files' rereads
Closure: RD1/rederive-mandate, RD2/per-axis-derivations, RD3/dispositions, RD4/committed-pickup, RD5/dual-counts, RD6/anchors-green

## What to build

Rewrite `craft-review` (keep ≤120 lines) so every axis derives its facts from
the current primary source before comparing the candidate: Coverage
independently enumerates the producer-derived input family and the
spec-authorized write set from the approved spec; Spec drives the behavior and
quotes the applicable spec line rather than trusting ticket claims; Standards
independently reads the current conventions and keeps the Fowler smell
baseline (`references/smell-baseline.md` stays the owner). A declaration-only
confirmation is incomplete; findings cite their derivation. The three axes run
in parallel fresh contexts so one derivation cannot seed another, and review
asks what authenticates the verifier before any candidate-controlled
execution. Review treats a compiled map's defaulted decisions as authoritative
unless the spec overrides them; a claimed repair is checked against both its
coverage row and the applicable defaulted-decision table. Update
`bench-review-implementation.md`: every finding carries exactly one
disposition — `no-op` (refuted, no repair target), `auto-fix` (deterministic
rule or exact spec predicate inside approved scope), `ask-user` (judgment,
scope, authority, or oracle change) — as repair-routing labels, never edit
permission for the read-only phase; actionable findings are written to
`reviews/<slug>.md` and committed as an ordered step before repair begins,
including findings returned by another harness; a clean review writes no
pickup artifact; the handoff reports raw per-axis counts and the de-duplicated
repair-target count separately. Migrate the 17 command anchors and 2 skill
anchors with their canaries; keep surviving obligations pinned. Won't-handle
(explicit unused disposition): tracker-backed maps and the upstream two-axis
review remain closed non-adoptions — this ticket keeps local decision maps and
all three axes, adding no adoption route.

## Acceptance

- [ ] [RD1] (covers PG10) `craft-review` mandates re-derive-then-compare per axis with derivation-citing findings; declaration-only confirmation named incomplete.
- [ ] [RD2] (covers PG11) the per-axis derivations, parallel fresh contexts, verifier-bootstrap question, and Fowler baseline are all present.
- [ ] [RD3] (covers PG12) both files carry the three-value disposition vocabulary with exactly-one-per-finding.
- [ ] [RD4] (covers local) the ordered write-and-commit-before-repair step covers cross-harness returns; clean review writes nothing.
- [ ] [RD5] (covers local) the handoff states raw axis counts and unique repair targets separately.
- [ ] [RD6] (covers local) `go test ./internal/anchors ./internal/conformance` green after anchor migration; defaulted-decision authority stated.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RD1/rederive-mandate | let Coverage consume ticket claims | semantic review reread | reviewer-graded against PG10's planted-omission scenario |
| RD2/per-axis-derivations | drop the verifier-bootstrap question | anchors check | remove it, run the docs-currency check, expect the owning anchor's red |
| RD3/dispositions | allow a finding with no disposition | semantic review reread | reviewer-graded against PG12 |
| RD4/committed-pickup | move the commit step after repair | anchors check | reorder, run the docs-currency check, expect the ordering anchor's red |
| RD5/dual-counts | merge the two counts into one | semantic review reread | reviewer-graded against PG12 |
| RD6/anchors-green | leave a retired review Needle live | docs-currency-workflow check | run the check, expect the missing-needle diagnostic |
