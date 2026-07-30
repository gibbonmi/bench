# Resolve owned worktree targets

Blocked by: none

## What to build

Sessions can resolve one active Bench-owned assignment by exact label, absolute
path, or portable home-relative path and print its safe portable path without
mutating assignment, Git, or filesystem state.

## Acceptance

- [x] An exact unique active label resolves its owned present tree.
- [x] Labels, absolute paths, `~`, and `~/...` select the same assignment, and
  path-component-aware home compaction preserves outside-home prefix siblings.
- [x] Spaces and glob characters remain literal while control-bearing output,
  whitespace-only targets, `~user`, and relative paths are rejected.
- [x] Zero, duplicate, inactive, unregistered, foreign, and ownership-mismatched
  targets fail closed without changing the ledger, Git registrations, or files.
