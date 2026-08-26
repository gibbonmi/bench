# Recover the landing broker toolchain

Blocked by: none
Writes: bin/bench.sh, internal/systemtest/land_route_test.go, CHANGELOG.md, ROADMAP.md, roadmap/FT197.md

## What to build

`bench worktree land` recovers Go from the installed package before it starts
the authenticated promotion broker. Users do not add a toolchain directory to
`PATH` when the harness supplies a partial environment.

## Acceptance

- [x] The installed wrapper gives its recovered Go path to the authenticated promotion broker.
- [x] The wrapper authenticates the broker before it runs the clean-login probe.
- [x] The landing route does not read the repository before it selects the broker.
