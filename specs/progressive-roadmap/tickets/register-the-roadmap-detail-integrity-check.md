# Register the roadmap-detail-integrity conformance check with its canary family

Blocked by: split-the-board-parser-and-migration-in-one-green.md
Writes: internal/roadmap, internal/conformance, internal/conformance/registry, internal/conformance/tier_test.go, tests/canary/roadmap-detail-integrity, projects/benchkit.md

## What to build

A `func(root string) []string` in the roadmap package returns the loader's
diagnostics; the conformance table and registry bind it as
`roadmap-detail-integrity` (Dev, SubjectRoot, its own input source naming
`ROADMAP.md` and `roadmap/`); a canary family of the same name carries one fixture
per diagnostic class with `EXPECT` and `MUTATE.json`, the check's own test asserts that fixture inventory and is listed in
`tier_test.go`'s live-tree classification, and `projects/benchkit.md` gains the
check-input row in registry order. Each fixture's `files/roadmap/` carries a file
for every index row it keeps, because restore overlays the live board. Lands after the migration so this repo's board is
green under the check the moment it is registered. Coverage row PR15.

## Acceptance

- [ ] The registry and conformance table both carry `roadmap-detail-integrity` bound to the roadmap validator at Dev tier and SubjectRoot.
- [ ] Every fixture in the family emits exactly its `EXPECT` diagnostic when mutated and no longer emits it after `RestoreMutationFixture`.
- [ ] `projects/benchkit.md` carries the check-input row and `tier_test.go` classifies the live-tree test.
- [ ] Deleting one named fixture from the inventory reds the check's own test.
- [ ] `bench gate` is green on this repo's migrated board.
