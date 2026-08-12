# Slim craft-tickets to the tracer contract and implement-spec to the orchestration pointer

Blocked by: 03-frontier-grill-and-tdd-seam-gate.md
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `.agents/commands/bench-implement-spec.md`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors`, `.bench/BENCH-reference.md`, `internal/conformance/fixture_bite_test.go`, `internal/conformance/registry_test.go`
Integration surfaces: ticket-shape consumers→`.agents/commands/bench-implement-spec.md`; anchors→`internal/anchors/registry_data.go`; anchor canaries→`tests/canary/workflow-guidance-anchors`; index currency→`.bench/BENCH-reference.md`
Contracts: the slim ticket field set (title, `Blocked by:`, What to build, Acceptance, advisory `Writes:`) crossing `.agents/skills/bench-craft-tickets/SKILL.md`→`.agents/commands/bench-implement-spec.md`, asserted by ST4 against the rewritten command's own reread
Closure: ST1/slim-contract, ST2/retired-fields-gone, ST3/reviewer-approved-breakdown, ST4/implement-spec-pointer, ST5/anchors-migrated

## What to build

Rewrite `craft-tickets` to ≤100 lines holding only the durable contract: the
independently-green tracer rule, `Blocked by:` frontier order, `What to build`,
`Acceptance`, the advisory `Writes:` note used only to judge parallel
disjointness, serial commit-on-green through path-scoped `bench commit`, and
the reviewer-approved breakdown: before assignment the coordinator presents a
numbered list of title, `Blocked by:`, and delivered outcome, iterates it with
the reviewer, and records approval — the batch-approval AFK carve-out stays the
only no-round-trip route. Remove (not summarize) Contracts, Integration
surfaces, Closure, covers annotations, red-mutation tables, handoff ledgers,
fence enforcement, and the read-only breakdown-review delegate; a meaningful
contract is stated in What to build and Acceptance and re-derived from the tree
by review. Rewrite `bench-implement-spec.md` to ≤60 lines as the orchestration
pointer: preflight entry, line declaration, ticket derivation per the slim
shape, reviewer breakdown approval, serial commit-on-green landing, TDD at
agreed seams, review before final landing, `--spec` final flip, `--full`
phase-boundary handoff, and stop-short/blocked routes in compact form. Migrate
the 45 craft-tickets and 54 implement-spec anchors: retire pins whose only
subject is deleted ceremony together with their canary fixtures under
`tests/canary/workflow-guidance-anchors`; keep one owning anchor per surviving
obligation (independently-green rule, `Blocked by:`, What to build, Acceptance,
reviewer approval, commit-on-green). Never copy anchor counts into prose.

## Acceptance

- [ ] [ST1] (covers PG16) `craft-tickets` is ≤100 lines and specifies exactly the slim field set with frontier order and serial commit-on-green.
- [ ] [ST2] (covers local) `rg -i "Integration surfaces|Closure:|red mutation|handoff ledger|covers annotation|ownership fence" .agents/skills/bench-craft-tickets/SKILL.md` returns nothing.
- [ ] [ST3] (covers PG9) both files route breakdown approval to the reviewer as an iterated numbered title/blocked-by/outcome list; no breakdown-review delegate remains; the AFK batch carve-out survives.
- [ ] [ST4] (covers local) `bench-implement-spec.md` is ≤60 lines and consumes the slim ticket shape without restating its schema.
- [ ] [ST5] (covers local) `go test ./internal/anchors ./internal/conformance` green after the anchor migration.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| ST1/slim-contract | delete the `Writes:` note rule | anchors check | remove it, run the docs-currency check, expect the owning anchor's red |
| ST2/retired-fields-gone | reintroduce a `Contracts:` schema line | semantic review reread | the ST2 rg probe returns a hit; review cites PG16 |
| ST3/reviewer-approved-breakdown | route breakdown review back to a delegate | anchors check | swap the sentence, run the docs-currency check, expect the reviewer-approval anchor's red |
| ST4/implement-spec-pointer | paste the ticket schema into the command | semantic review reread | reviewer-graded: duplicated knowledge and budget breach |
| ST5/anchors-migrated | leave one ceremony Needle live | docs-currency-workflow check | run the check, expect the missing-needle diagnostic |
