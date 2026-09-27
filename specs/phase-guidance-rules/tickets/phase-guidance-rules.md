# State the ticket check floor, two repair causes, the confirming Standards focus, and one retro per spec

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, .agents/commands/bench-final-check.md, .agents/commands/bench-review-implementation.md, AGENTS.md, internal/anchors, tests/canary/workflow-guidance-anchors
Covers: none

## What to build

The phase guidance states four rules, each from one source. The implementation phase states the checks that a ticket author runs before the ticket commit. The final-check retro vocabulary gains the `one-source` and `check-gap` repair causes. The anchor registry needle for that vocabulary changes in place to match. A confirming review round charges the Standards axis to look first for duplicated knowledge in the repair delta.

The working agreement states that only the implementation close writes the retro. A spec-stage close commits its scorecard updates and no retro.

## Acceptance

- [ ] The revalidation paragraph of `bench-implement-spec.md` states the check floor before the ticket commit. The floor is root conformance and the tests of each package that the ticket writes. It adds the whole `cmd/bench` package when the ticket changes a public response or an embed pattern. No other guidance file states this floor.
- [ ] `bench-final-check.md` lists `one-source` and `check-gap` in the repair-cause vocabulary and defines each term in one clause. A probe that omits `check-gap` from the guidance turns `docs-currency-workflow` red.
- [ ] `bench-review-implementation.md` states that a confirming round charges the Standards axis to look first for duplicated knowledge in the repair delta. The pinned confirming-round sentence stays byte for byte.
- [ ] The `**Phase-close handoff.**` paragraph in `AGENTS.md` states that only the implementation close writes the retro, one for each spec. It states that a spec-stage close commits its scorecard updates and no retro.
