# Retained implementation workflow review pickup

Frozen delta: `6cb6314ad20c953cf8ae5f5d41e5deeba2f9dc1f..53fc8a342093111d08077fc1b239d44c86e8ea45`.

## Standards

Finding count: 1. Worst issue: P2.

1. [P2] `auto-fix` — Reconcile the retrospective heading with its canonical schema. `.agents/commands/bench-final-check.md:87` requires a heading that `internal/retros/retros.go:33` rejects. Preserve historical retrospective parsing and align the guidance with the existing schema owner.

## Spec

Finding count: 3. Worst issue: P1.

1. [P1] `auto-fix` — Remove the retired spec from the permanent docs gate. `internal/conformance/retained_workflow_test.go:34` unconditionally requires a spec that final-check retires. Validate the durable contract and use fixture tables for the staged spec's structural mutations.
2. [P2] `auto-fix` — Reconcile the retrospective heading with its canonical schema. This duplicates the Standards finding and has one repair target.
3. [P2] `auto-fix` — Replace stale fresh-session handoffs in `.agents/commands/bench-write-spec.md` and `projects/benchkit.md`, update their anchor, and add negative coverage.

## Coverage

Finding count: 1. Worst issue: P1.

1. [P1] `auto-fix` — Cover the retired-spec state before making the spec a permanent gate input. This duplicates Spec finding 1 and has one repair target.
