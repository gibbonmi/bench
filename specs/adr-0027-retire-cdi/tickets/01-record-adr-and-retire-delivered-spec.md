# Record ADR 0027 and retire the delivered FT390 spec

Blocked by: none
Writes: docs/adr/0027-a-commitment-binding-names-its-obligations.md, specs/commitment-delivery-integrity, reviews/commitment-delivery-integrity.md
Covers: none

## What to build

The commitment-delivery-integrity spec is implemented, and FT390 is delivered. Three of its decisions have no durable home outside the spec: the obligation rule, the open-content rule, and the source order.
ADR 0027 records those decisions as the current state. Then `bench spec retire commitment-delivery-integrity` removes the spec folder and its review record.

## Acceptance

- [ ] ADR 0027 states that a binding names its obligations, and that a plan and the completion landing refuse a binding that names none.
- [ ] ADR 0027 states that a plan binds only the content that it keeps open, and that a settled deliverable stays unbound.
- [ ] ADR 0027 states that the plan identity sorts its sources once, in one canonical byte order.
- [ ] The spec folder and its review record are absent, and the landing commit subject ends with `spec-retire: commitment-delivery-integrity`.
