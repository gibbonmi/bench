# Repair dangling cross-references left by the doctrine slimming

Blocked by: 09-prose-budget-conformance.md
Writes: `.agents/skills/bench-craft-spec/SKILL.md`, `.agents/skills/bench-craft-review/SKILL.md`, `.agents/commands/bench-debug.md`, `.agents/commands/bench-shape-idea.md`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors`

## What to build

Ticket 10's composed sweep found four live references to prose the slimming
removed. Repair each at its own site; no budget subject may exceed its limit
(craft-spec and craft-review sit exactly at 120 — replacements, not additions).

1. `craft-spec` points to `` `Discover the contracts before writing files` ``
   in `craft-tickets` — a section ticket 04 removed. The pointer is gate-pinned
   by a `RequireInSection` anchor. Repair the sentence to state the surviving
   doctrine (a contract between tickets is stated in What to build and
   Acceptance and re-derived from the tree by review) and migrate the anchor's
   Needle to the new wording in the same edit; update any canary fixture
   naming the old text.
2. `craft-review` cites "the prefactor rule lives in `craft-tickets`" and
   "`craft-tickets`' integration-surface discovery" — neither exists. Reword
   to cite what the slim skill actually owns (independently-green tracer
   grouping and the advisory `Writes:` disjointness note), or drop the
   citation where the sentence stands without it.
3. `bench-debug.md` still describes the retired debug receipt, assignment ID,
   and an `assign --refresh` verb, and points at implement-spec's removed
   blocked-outside-fence section. Rewrite that passage to the surviving
   contract: a blocked write-delegate stops at its boundary and returns a
   bounded blocked report (repro command, red output digest, failing surface,
   in-scope dirty paths); the coordinator validates it and reslices repair
   tickets; no receipt artifact, no assign verbs.
4. `bench-shape-idea.md` names the Prototype decision-ticket type but never
   points to the `prototype` skill. Add the one-line pointer at the ticket
   type's definition.

## Acceptance

- [ ] [CR1] no file under `.agents/` or `.bench/` references a craft-tickets section, receipt artifact, or assign verb that no longer exists (rg for the four quoted strings above returns nothing outside specs/, decisions/, capture/, CHANGELOG history).
- [ ] [CR2] the migrated craft-spec anchor bites: mutating the new sentence reds docs-currency-workflow, restoring greens it.
- [ ] [CR3] `wc -l` of craft-spec and craft-review stays ≤120; `go test ./internal/anchors ./internal/conformance` green; `.bench/skills-index.sh --check` green.
- [ ] [CR4] shape-idea's Prototype ticket type names the `prototype` skill.
