# Exercise the prospective bundle root

Blocked by: none
Writes: internal/gate/prospective.go, internal/gate/prospective_owner_test.go
Covers: LF20

## What to build

Remove the test-only artifact-root branch or drive it through production.
Make the own-branch executable test use a nonempty root and prove bundle
confinement inside the graded tree.

## Acceptance

- [ ] Production and tests use one artifact-root path.
- [ ] The executable-authoring test passes a nonempty artifact root.
- [ ] The authored bundle cannot escape the graded tree.

