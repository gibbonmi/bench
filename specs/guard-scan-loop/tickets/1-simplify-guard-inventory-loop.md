# Simplify guard inventory scanning

Blocked by: none
Writes: internal/guards/guards.go, internal/guards/scan_loop_test.go (new), internal/guards/testdata/scan-family-s.golden (new), internal/conformance/guard_scan_owner_test.go (new), internal/conformance/package_core_checks_test.go, internal/conformance/checks_test.go, internal/conformance/tier_live_tree_test.go
Covers: GS01, GS02, GS03, GS04, GS05, GS06, GS07, GS08, GS09, GS10, GS11, GS12, GS13, GS14, GS15, GS16, GS17, GS18, GS19, GS20, GS21, GS22, GS23, GS24, GS25, GS26, GS27, GS28, GS29, GS30

## What to build

Replace the worker pipeline in Scan with one synchronous loop on the caller goroutine, and keep every observable result.
Keep Scan, Rows, Command, ScanResult, and their signatures unchanged. Rows returns `[][]string`; add no row type.
Keep `guardScanTimeout`, `enumerateGuards`, and `inspectGuard` as the package-variable seams, with their names and types.
Keep the policy line `bounds.VerdictWindow(bounds.GuardScanTimeout)` at its current text, because bounds-policy pins it with `strings.Contains`.

Enumerate once, then inspect the candidates in their current order.
Do a context check before enumeration, after enumeration, before each inspection, and after each inspection.
A pre-cancelled context starts no operation.
If the check after enumeration sees cancellation, return unknown Total and Omitted and start no inspection.
This rule applies to an error return and to a candidate slice with a nil error.
Count an inspection only when its post-operation check sees an active context. An empty-row inspection still counts.
A cancelled inspection adds no rows and no count. Completed rows stay in order.
After the final completion check, a later cancellation does not change the complete result.
Scan returns only after the current operation returns. Do not promise a hard filesystem deadline.
Keep the active-context enumeration-error result, the static header, wiring, no-manifest, symlink, and FIFO behavior, and all command output, help, and exit codes.

Before you replace the loop, build the contract driver for family S and run it on the pre-change code.
Store that reference capture as `internal/guards/testdata/scan-family-s.golden`, and have it reviewed apart from the candidate code.
GS27 then compares the candidate Command tuples with the golden over the same controlled events, without elapsed-time equality.

Add `checkGuardScanLoopOwner` beneath package-core-guard.
It grades the production Go AST under internal/guards, unused function bodies included, and refuses go statements, select statements, channel construction, and send or receive nodes.
It engages only when the root holds an `internal/guards` directory. A missing or unreadable `guards.go` in that directory gets a named diagnostic; a root without the directory gets none.
Classify its live-tree test in `tier_live_tree_test.go`.

If `checkPackageCoreAndGuards` is still in package_core_checks_test.go, move only that 16-line declaration to checks_test.go, and add the new subcheck there.
If the shared-guard-projection outcome moved it first, extend that sole declaration and keep its subcheck.
Keep the name, signature, binding, registry metadata, and subcheck order. Keep the npm producers and their formatter call in package_core_checks_test.go.

Keep guards_test.go, guards_cleanup_test.go, and the four testdata stdout fixtures unchanged.
Keep the test names `TestScanTimeoutPreservesPartialRowsAndHonestCounts` and `TestScanEnumerationTimeoutUsesUnknownCounts`, because the race phase names them (`internal/gate/gate_go.go:36-37`).
Do not add the banned skills-index literals to checks_test.go (`skills_index_checks_test.go:375-387`).
Keep checks_test.go within its 739-line grant after either outcome order. package_core_checks_test.go must not grow. Each new file stays within 400 lines. Add no grant.

## Acceptance

- [ ] A complete scan enumerates once and keeps candidate and row order, and Rows returns `Scan(background).Rows` (GS01, GS02, GS22).
- [ ] A pre-cancelled scan starts no operation (GS03).
- [ ] Cancellation seen after enumeration keeps Total and Omitted unknown for an error return and for a nil-error candidate slice, and starts no inspection (GS04, GS30).
- [ ] Cancellation between candidates keeps the completed prefix and Omitted = Total − Inspected (GS05, GS06).
- [ ] An empty-row inspection counts, and a cancelled inspection adds no rows or count (GS07, GS09).
- [ ] TestScanWaitsForCancelledWorkerCleanup stays unchanged and passes, and a context-ignoring operation keeps Scan pending until its release (GS08, GS11).
- [ ] All ScanResult fields keep their values, and a late cancellation leaves a complete result unchanged (GS10, GS28).
- [ ] An active-context enumeration error keeps the current complete result (GS12).
- [ ] Absent and empty hook directories, static headers, symlinks, FIFOs, and wiring keep their current rows (GS13–GS16).
- [ ] The four checked-in stdout fixtures, help order, help suppression, and command errors stay unchanged (GS17–GS21).
- [ ] A renamed unused copy of the old pipeline makes package-core-guard red; a test goroutine does not (GS23).
- [ ] A root without `internal/guards` passes, and a root whose `internal/guards` lacks a readable `guards.go` gets the named diagnostic (GS24, GS25).
- [ ] The live-tree ownership test runs through its registered binding and its classification (GS26).
- [ ] Family S candidate tuples equal `scan-family-s.golden`, and reversing the candidate order makes GS27 red (GS27).
- [ ] guards.go keeps the text that bounds-policy pins (GS29).
- [ ] Each completion-plan probe compiles, turns its named witness red, and returns to the same green after byte-identical restoration.
