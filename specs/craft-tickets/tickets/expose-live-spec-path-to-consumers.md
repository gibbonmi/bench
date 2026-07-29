# Expose live spec path to consumers

Blocked by: Expand spec resolution with folder form

## What to build

Expose each enumerated spec's actual repository-relative path from the
`internal/spec` deep unit. Consumers can then render folder specs accurately
without re-deriving layout rules, while transitional flat specs retain their
real path until the contract ticket removes them.

## Acceptance

- [x] Folder facts carry `specs/<slug>/spec.md`.
- [x] Transitional flat facts carry `specs/<slug>.md`.
- [x] Consumers need no file-existence probe to choose a layout.
