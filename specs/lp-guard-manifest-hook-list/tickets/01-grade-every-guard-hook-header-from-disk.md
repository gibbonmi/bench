# Grade every guard hook header from disk

Blocked by: none
Writes: internal/conformance/guard_header_manifests_test.go, internal/guards/guards.go
Covers: none

## What to build

The guard header manifest check grades each hook script that guard discovery finds in `.bench/hooks`. It does not use a hand-written list of hooks. Guard discovery and the check use one exported selector in `internal/guards`, so a new hook cannot escape the check. The check still grades the pre-push hook source. The `bench guards` report does not change.

## Acceptance

- [ ] A hook script in `.bench/hooks` that no list names and that has no manifest header causes the check to report a missing key.
- [ ] When the `# denies:` line is removed from `block-primary-file-write.sh`, the live-tree check fails.
- [ ] The `bench guards` rows do not change for the current tree.
