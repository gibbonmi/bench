# Migrate the tracked decision-map corpus

Blocked by: Validate the decision graph and readiness

## What to build

Every tracked, non-exempt active and compiled candidate present at this ticket's
start moves to the canonical schema without changing its decisions, answers,
exclusions, or research facts.

## Acceptance

- [x] The migration inventory explicitly covers the six tracked active maps and
  three tracked compiled maps present at ticket start; nested map assets remain
  excluded.
- [x] Active maps are honestly `shaping` or `ready`, and every compiled map is
  `ready`.
- [x] Every decision ticket has canonical blockers and type, and no answer or
  graph dependency changes meaning.
- [x] Handoff sections are removed only after their durable facts are retained in
  Destination, terminal sections, or the consuming spec.
- [x] Sources use the structured locator, Supports, and Drift form without losing
  a cited fact.
- [x] The complete tracked corpus validates through `internal/maps`, and a second
  migration pass changes nothing.
