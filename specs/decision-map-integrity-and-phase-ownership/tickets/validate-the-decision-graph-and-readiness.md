# Validate the decision graph and readiness

Blocked by: Expand the decision-map schema

## What to build

The new `internal/maps` model resolves dependency graphs, readiness, terminal
section shapes, and offline structured sources so every future consumer receives
one complete validation result.

## Acceptance

- [ ] Duplicate IDs and blockers, dangling blockers, self-edges, cycles, and
  resolved-on-unresolved blockers each earn a targeted diagnostic.
- [ ] Shaping maps may carry honest unresolved tickets or fog; ready maps reject
  every unresolved/deferred marker and non-empty fog.
- [ ] Path sources resolve only to in-root regular files, and URL sources accept
  only absolute HTTP(S) locators with hosts without network access.
- [ ] Empty Sources is valid while unstructured sentinels and incomplete source
  entries are rejected.
- [ ] Non-empty fog, discretion, and out-of-scope bodies require Markdown bullet
  lists.
- [ ] Validation continues across independent failures and names the map,
  decision ticket, graph handle, and bad edge where applicable.
