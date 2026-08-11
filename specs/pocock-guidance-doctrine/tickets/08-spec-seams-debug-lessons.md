# Land the spec-reread, read-budget, and debug-isolation lessons; trim craft-spec

Blocked by: 07-bench-and-agents-lessons.md
Ownership fence: `.agents/skills/bench-craft-spec/SKILL.md`, `.agents/skills/bench-craft-seams/SKILL.md`, `.agents/commands/bench-debug.md`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors`, `.bench/BENCH-reference.md`
Integration surfaces: craft-spec section-scoped anchors→`internal/anchors/registry_data.go`; anchor canaries→`tests/canary/workflow-guidance-anchors`; index currency→`.bench/BENCH-reference.md`
Contracts: the section headings craft-spec's RequireInSection anchors resolve (`Slicing a build for delegates`, `The edge inventory`, `The acceptance coverage map`) crossing `.agents/skills/bench-craft-spec/SKILL.md`→`internal/anchors/registry_data.go`, asserted by SD4 against the real docs-currency check
Closure: SD1/spec-reread-and-trim, SD2/seams-read-budget, SD3/debug-entry-isolation, SD4/anchors-green

## What to build

Three lesson landings. (1) `craft-spec`: add the whole-artifact reread rule —
after wide structured-prose edits, reread the complete artifact and reconcile
before handing off — and trim the skill to ≤120 lines while preserving the
row schema, red-signal standard, edge-inventory classes, story sizing, and
delegate-slicing content its 11 section-scoped anchors pin; where a pinned
section must shrink, migrate the anchor to the surviving wording rather than
deleting the obligation. (2) `craft-seams`: add the observable
exploration-read budget — declare a small read budget for exploration; when it
is spent without traction, reroute through `bench outline` and report the
budget spent (PG19's seam half; keep ≤120 lines and the ticket-02 leaf
pointers). (3) `bench-debug.md`: debug work that may write creates or selects
its isolated worktree *before* the first repro artifact exists, so a
clean-at-start checkout cannot become unattributably dirty mid-debug; do not
generalize the one-build Opus delegate-debug override. Migrate the affected
anchors (11 craft-spec, 3 craft-seams, 6 bench-debug) and any canary fixture
naming moved wording.

## Acceptance

- [ ] [SD1] (covers PG21) the complete surviving-lesson family is explicit across the landed tree: review capture, dual counts, verifier bootstrap, and defaulted-decision authority (landed by ticket 06 — assert by reread), the acceptance-shortfall exit (ticket 07 — assert by reread), and this ticket's whole-artifact reread rule in `craft-spec`, which is ≤120 lines with schema, red-signal, edge classes, and slicing content intact.
- [ ] [SD2] (covers PG19) both PG19 halves are explicit: BENCH.md's source-before-call predicate (ticket 07 — assert by reread) and `craft-seams`' declared read budget, `bench outline` reroute, and report-budget-spent obligation.
- [ ] [SD3] (covers PG23) `bench-debug.md` requires isolation before the first repro artifact and keeps the Opus override non-standing.
- [ ] [SD4] (covers local) `go test ./internal/anchors ./internal/conformance` green after anchor migration.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SD1/spec-reread-and-trim | delete the reread rule | anchors check | remove it, run the docs-currency check, expect the new owning anchor's red |
| SD2/seams-read-budget | drop the report-budget-spent clause | semantic review reread | reviewer-graded: the reroute stops being observable against PG19 |
| SD3/debug-entry-isolation | defer isolation to the fix phase again | anchors check | reorder, run the docs-currency check, expect the entry-isolation anchor's red |
| SD4/anchors-green | rename a pinned section heading without migrating | docs-currency-workflow check | run the check, expect resolveSection's missing-section diagnostic |
