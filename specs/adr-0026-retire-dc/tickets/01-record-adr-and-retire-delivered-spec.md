# Record ADR 0026 and retire the delivered commitment spec

Blocked by: none
Writes: docs/adr/0026-work-class-sets-the-default-delivery-order.md, specs/roadmap-delivery-commitment, reviews/roadmap-delivery-commitment.md
Covers: none

## What to build

The roadmap-delivery-commitment spec is implemented, and quality-1 is adopted. Three of its decisions have no durable home outside the spec: the default class order, its current priority, and the out-of-scope list.
ADR 0026 records those decisions as the current state. Then `bench spec retire roadmap-delivery-commitment` removes the spec folder and its review record.
The adoption proposal leaves with the spec folder, because the adoption consumed it. Git history keeps it.

## Acceptance

- [ ] ADR 0026 states the default order of confirmed defects, then refactoring, then new features, and states that an explicit reviewer order governs.
- [ ] ADR 0026 states that a routine drain never displaces or admits work.
- [ ] ADR 0026 states that an existing sequence never becomes a commitment automatically.
- [ ] The spec folder and its review record are absent, and the landing commit subject ends with `spec-retire: roadmap-delivery-commitment`.
