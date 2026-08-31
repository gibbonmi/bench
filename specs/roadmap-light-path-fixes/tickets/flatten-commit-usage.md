# Flatten commit usage errors

Blocked by: route-doctor-and-grade-nested-dispatch.md
Writes: internal/commit/commit.go, internal/commit/commit_test.go, internal/conformance/subcommand_routing_test.go
Covers: LF8

## What to build

Make bench commit usage failures print one flat contract. Keep the top-level
worktree dispatcher exempt because supported leaves own their grammars.

## Acceptance

- [ ] Commit grammar errors print one usage line.
- [ ] Nested usage text no longer appears.
- [ ] The routing census retains the worktree dispatcher exemption.
