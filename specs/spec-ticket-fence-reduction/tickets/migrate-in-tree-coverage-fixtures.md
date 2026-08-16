# Migrate in-tree coverage fixtures

Blocked by: accept-the-reduced-coverage-header.md, project-one-row-shape-across-schemas.md
Writes: internal/preflight/command_bootstrap_test.go, internal/preflight/command_review_test.go, internal/preflight/command_harness_test.go, internal/worktree/land_test.go, internal/systemtest/owner_test.go, cmd/bench/command_registry_test.go

## What to build

Every in-tree test fixture that embeds an acceptance coverage map uses the
reduced header, so the reduced schema is exercised through its real consumers —
preflight's spec gather, the worktree landing path, the system test's ownership
fences, and the command registry — rather than only through the coverage
package's own tests. This also shrinks the eventual contract to a parser-and-two-specs
edit instead of a fixture sweep.

Fixtures keep asserting whatever they assert today; only their embedded header
and row widths change. A fixture whose assertion depends on the red-signal cell's
content is a finding to surface, not to paper over.

## Acceptance

- [ ] None of the six named Go test files embeds a six-column coverage map.
- [ ] `tests/canary/coverage-map-validation/broken-coverage-map` is left on its
      legacy header on purpose — it is what proves the legacy branch still
      validates — and that intent is stated where the fixture lives.
- [ ] The preflight, worktree, systemtest, and command-registry suites pass with
      their fixtures migrated, asserting the same outcomes as before.
- [ ] No fixture assertion was weakened or deleted to accommodate the new header.
