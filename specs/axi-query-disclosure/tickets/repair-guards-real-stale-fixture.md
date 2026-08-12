# Repair guards real stale fixture

Blocked by: none
Writes: `internal/guards/`

## What to build

Give the guards stale disclosure a real fixture (review finding R7): the
stale/unwired QD6 evidence is injected through the `enumerateGuards`/
`inspectGuard` package variables with test-authored cells, against the spec's
testing decision "real fixtures (a stale guard manifest, …)". Add one
public-command test that installs a genuinely stale managed pre-push hook in a
fixture repository and observes the `bench link` repair action end to end. The
existing stub-driven tests stay.

## Acceptance

- [ ] [GS1] (covers QD1) a public `guards` command test installs a real
  managed pre-push hook whose manifest the production scan classifies stale,
  and asserts the rendered row plus the appended `bench link` action from
  production derivation alone — no `enumerateGuards`/`inspectGuard` stubs.
- [ ] [GS2] (covers QD1) the test demonstrates its red by the row-shape route
  the stubs cannot catch (a reordered or mislabeled produced cell fails it),
  recorded in the ticket evidence, then green.
