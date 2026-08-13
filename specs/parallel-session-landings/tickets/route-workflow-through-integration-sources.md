# Route the workflow through reviewed integration sources

Blocked by: land-reviewed-sources-atomically.md, resume-published-landings.md
Writes: .bench/BENCH.md, .bench/BENCH-reference.md, .agents/commands/bench-write-spec.md, .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, .agents/commands/bench-final-check.md, .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-tickets/SKILL.md, projects/benchkit.md, docs/field-guide.html, docs/greenfield-build-sequence.md, CHANGELOG.md, CONTEXT.md, internal/conformance/docs_workflow_helpers_test.go, internal/anchors/registry_data.go, tests/canary

## What to build

Rewrite the twelve named current-state surfaces as one workflow contract:
platform workflow, spec authoring, implementation, semantic review, final check,
delegation, ticketing, profile, reference, field guide, greenfield sequence, and
help. They route tickets into one retained integration source, review the
explicit frozen pair, and hand that source to `bench worktree land`; remove the
stale claim that path-scoped `bench commit` is the sole landing route.

Move the suggested fresh-session boundary to after reviewed ticket slicing, so
map, spec, and ticket planning stay in one decision context. Keep executable
grammar single-sourced in usage. For every edited subject named in the profile's
guidance-prose budget table, make compensating cuts rather than raising its
limit. Update the existing anchors and canaries so reverting any required
current-state sentence makes the gate red.

## Acceptance

- [ ] All twelve enumerated surfaces agree on integration-source review and
      landing, and the spec-authoring surface recommends a fresh session only
      after tickets are sliced (covers PL23).
- [ ] The stale scalar-base and sole-path claims are absent, while command flags
      remain owned by executable help and conformance checks bite on the new
      workflow anchors (covers PL23).
- [ ] Every edited budgeted subject remains at or below its current profile limit
      through demonstrated compensating cuts; unbudgeted subjects are not folded
      into that assertion (covers PL32).
