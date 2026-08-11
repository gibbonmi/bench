# Rewrite grilling to frontier rounds and give light-path TDD a seam gate

Blocked by: 02-add-prototype-and-craft-leaves.md
Ownership fence: `.agents/skills/bench-craft-grill/SKILL.md`, `.agents/commands/bench-shape-idea.md`, `.agents/skills/bench-craft-tdd/SKILL.md`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors`, `.bench/BENCH-reference.md`
Integration surfaces: grill anchors→`internal/anchors/registry_data.go`; anchor canaries→`tests/canary/workflow-guidance-anchors`; frontier vocabulary→`.agents/commands/bench-shape-idea.md`; index currency→`.bench/BENCH-reference.md`
Contracts: each anchor Needle crossing `.agents/skills/*/SKILL.md`→`internal/anchors/registry_data.go`, asserted by FG4 against the real docs-currency-workflow check via `go test`
Closure: FG1/frontier-rounds, FG2/light-path-seam-stop, FG3/spec-backed-no-second-gate, FG4/anchors-migrated

## What to build

Replace `craft-grill`'s one-question-at-a-time cadence with numbered frontier
rounds: every question whose prerequisites are settled appears in the same
round, each with a recommendation; the skill waits, recomputes the frontier
from the answers, and still closes through shared-understanding confirmation.
Facts stay the agent's job, decisions the reviewer's. Update the frontmatter
`description`/`index:` (they currently say one-question-at-a-time) and keep the
close-each-decision-as-exact-predicate obligation. `bench-shape-idea.md` adopts
the same frontier vocabulary. In `craft-tdd`, add the light-path seam gate:
work with no spec sign-off stops before its first TDD test and presents its
seam for reviewer confirmation, while spec-backed work consumes the signed-off
seam without a second gate; keep SKILL.md ≤120 lines and its ticket-02 leaf
pointers intact. Migrate anchors: rewrite or retire the pins on `craft-grill`
("close each decision by restating the answer as the exact predicate it
fixes", "Grill is a decision-ticket type") and `craft-tdd` (6 anchors) so
surviving obligations keep one owning anchor; update any
`tests/canary/workflow-guidance-anchors` fixture whose EXPECT/MUTATE text names
retired prose. Regenerate the skills index if frontmatter changed.

## Acceptance

- [ ] [FG1] (covers PG7) `craft-grill` mandates whole-frontier numbered rounds with recompute and final confirmation; no one-question-per-turn rule survives anywhere in the skill.
- [ ] [FG2] (covers PG8) `craft-tdd` states the light-path stop-for-seam-confirmation rule.
- [ ] [FG3] (covers local) `craft-tdd` states spec-backed work takes the signed-off seam without a second reviewer gate.
- [ ] [FG4] (covers local) `go test ./internal/anchors ./internal/conformance` is green: every surviving anchor resolves, and retired pins are removed with their canary subjects.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FG1/frontier-rounds | restore a "one question per turn" sentence | semantic review reread | reviewer-graded contradiction against PG7 |
| FG2/light-path-seam-stop | delete the light-path stop clause | anchors check | remove the clause, run `go test ./internal/conformance -run Anchors`, expect the owning anchor's red |
| FG3/spec-backed-no-second-gate | make spec-backed TDD also stop for confirmation | semantic review reread | reviewer-graded: duplicate human gate against PG8 |
| FG4/anchors-migrated | leave one retired Needle in registry_data.go | docs-currency-workflow check | run the check via `go test`, expect the missing-needle diagnostic |
