# Enforce one scan loop owner

Blocked by: 1-replace-scan-pipeline-with-loop.md
Writes: internal/conformance/guard_scan_owner_test.go (new), internal/conformance/package_core_checks_test.go, internal/conformance/checks_test.go, internal/conformance/tier_live_tree_test.go
Covers: GS23, GS24, GS25, GS26

## What to build

Make the registered conformance check refuse a second scan implementation in internal/guards.
Ticket 1 supplies a guards.go with no goroutine, select, or channel, so the live tree is clean when this check starts.

Add `checkGuardScanLoopOwner` beneath package-core-guard.
It grades the production Go AST under internal/guards, unused function bodies included, and refuses go statements, select statements, channel construction, and send or receive nodes.
It ignores test code, comments, and fixture strings.
It engages only when the root holds an `internal/guards` directory. A missing or unreadable `guards.go` in that directory gets a named diagnostic; a root without the directory gets none.
No package-core-guard canary tree lists a file under `internal/guards`, so no canary gets a new diagnostic.
Classify its live-tree test in `tier_live_tree_test.go`.

If `checkPackageCoreAndGuards` is still in package_core_checks_test.go, move only that 16-line declaration to checks_test.go, and add the new subcheck there.
If the shared-guard-projection outcome moved it first, extend that sole declaration and keep its subcheck.
Keep the name, signature, binding, registry metadata, and subcheck order. Keep the npm producers and their formatter call in package_core_checks_test.go. package_core_checks_test.go must not grow.
Keep checks_test.go within its 739-line grant after either outcome order, and do not add the banned skills-index literals (`skills_index_checks_test.go:375-387`).
The new file stays within 400 lines. Add no grant.

## Acceptance

- [ ] A renamed unused copy of the old pipeline makes package-core-guard red; a test goroutine does not (GS23).
- [ ] A root without `internal/guards` passes, and a root whose `internal/guards` lacks a readable `guards.go` gets the named diagnostic (GS24, GS25).
- [ ] The live-tree ownership test runs through its registered binding and its classification (GS26).
- [ ] Every retained canary fixture keeps its expected diagnostics.
- [ ] Omitting only the registered `checkGuardScanLoopOwner` call makes `TestGuardScanLoopOwnerBitesOnPrivateCopy` red, and restoration returns the same green.
