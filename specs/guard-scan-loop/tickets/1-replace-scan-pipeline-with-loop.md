# Replace the scan pipeline with a synchronous loop

Blocked by: none
Writes: internal/guards/guards.go, internal/guards/scan_loop_test.go (new), internal/guards/testdata/scan-family-s.golden (new)
Covers: GS01, GS02, GS03, GS04, GS05, GS06, GS07, GS08, GS09, GS10, GS11, GS12, GS13, GS14, GS15, GS16, GS17, GS18, GS19, GS20, GS21, GS22, GS27, GS28, GS29, GS30

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

Keep guards_test.go, guards_cleanup_test.go, and the four testdata stdout fixtures unchanged.
Keep the test names `TestScanTimeoutPreservesPartialRowsAndHonestCounts` and `TestScanEnumerationTimeoutUsesUnknownCounts`, because the race phase names them (`internal/gate/gate_go.go:36-37`).
guards.go stays within 400 lines, and each new file stays within 400 lines. Ticket 2 adds the ownership check after this loop lands.

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
- [ ] Family S candidate tuples equal `scan-family-s.golden`, and reversing the candidate order makes GS27 red (GS27).
- [ ] guards.go keeps the text that bounds-policy pins (GS29).
- [ ] Gating the post-enumeration check on the error alone makes `TestScanLoopEnumerationReturnsAfterCancellation` red, and removing the post-inspection check makes `TestScanLoopPostInspectionCancellation` red; restoration returns the same green.
