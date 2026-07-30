# Bind decision-map integrity to the gate

Blocked by: Project the active-map model

## What to build

One registered `decision-map-integrity` conformance check validates active and
compiled candidates through `internal/maps`, and a dedicated canary family keeps
every refusal live.

## Acceptance

- [ ] The focused check accepts valid active and compiled maps and specs with no
  decisions directory.
- [ ] An invalid compiled map makes the focused check and the real gate red with
  its compiled path and targeted diagnostic.
- [ ] Schema, graph, readiness, and source refusal classes each have a canary
  expectation matching their own diagnostic.
- [ ] The check records independent candidate diagnostics instead of stopping at
  the first failure.
- [ ] Compiled failures stay outside `bench maps` rows and ambient count.
