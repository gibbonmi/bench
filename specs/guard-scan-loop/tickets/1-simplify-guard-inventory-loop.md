# Simplify guard inventory scanning

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/checks_test.go, internal/conformance/guard_scan_owner_test.go (new), internal/conformance/package_core_checks_test.go, internal/conformance/registry/checks.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/tier_live_tree_test.go, internal/guards/guards.go, internal/guards/scan_loop_test.go (new)
Covers: GS01, GS02, GS03, GS04, GS05, GS06, GS07, GS08, GS09, GS10, GS11, GS12, GS13, GS14, GS15, GS16, GS17, GS18, GS19, GS20, GS21, GS22, GS23, GS24, GS25, GS26, GS27, GS28

## What to build

Deliver the complete synchronous Scan path, unchanged command formatting, and enforced removal of private worker pipelines in one green checkpoint.
Keep Scan, Rows, Command, ScanResult, Row, enumerateGuards, and inspectGuard signatures and existing injection ownership.
Enumerate once, then inspect in current order directly on the caller goroutine.
Check context before and after each operation and publish only candidates whose post-operation checkpoint is active.
An active empty-row inspection still increments Inspected; cancelled inspection contributes neither rows nor an inspected count.

Preserve the completed prefix, unknown enumeration counts, known omitted arithmetic, and current timeout reason.
After final completion, later cancellation does not reopen the complete result.

Scan waits for the current operation's cleanup even when it ignores context.
Do not promise a hard filesystem deadline or add a worker, parallel inspector, new interface, or new enumeration-error diagnostic.
Keep the active-context enumeration-error result and all static header, wiring, no-manifest, symlink, and FIFO dispositions.
Preserve existing full/brief stdout fixtures, complete action help order and deduplication, incomplete help suppression, and all return codes.

Prepare the canned contract driver before replacing Scan and retain its independently reviewed pre-change owner-package test executable.
Run reference and candidate against the same controlled family S events in their separate test processes through the real Command formatter.
Compare exact ScanResult and full/brief exit/stdout/stderr tuples; elapsed time and nondeterministic select ties are not equality predicates.
Use existing package-variable substitutions with cleanup, no parallel tests, and controlled release barriers for cancellation and blocking operations.
Keep all assertions in guards_test.go, guards_cleanup_test.go, and checked-in stdout fixtures unchanged.

Add checkGuardScanLoopOwner beneath package-core-guard, grading production AST nodes including unused function bodies.
Reject goroutine, select, channel construction, and send/receive pipeline nodes; do not reject channels in test fixtures.
Keep optional linked-source and named missing/unreadable kit-owner dispositions and explicit live-tree test classification.

If checkPackageCoreAndGuards still resides in package_core_checks_test.go, move only that existing 16-line declaration to checks_test.go in this checkpoint.
Keep its name, signature, executable binding, registry metadata, and existing subcheck order unchanged.
If the other guard outcome already moved it, extend that sole destination declaration and retain its landed ownership subcheck.
The two outcomes need no unlanded interface from each other, but overlapping conformance and command closure writes cannot run concurrently.
The coordinator serializes those charges against a fresh integrated source and rechecks headroom after the earlier outcome.

Keep npm producers and their one formatter call in package_core_checks_test.go.
Preserve TestCoreSubprocessFailuresUseProbeFormatter's independently authored path and count assertions.
Do not move the other producers, their tests, or their policy literals into the destination.
checks_test.go must remain within its existing 709-line grant after either or both outcomes.
The source file shrinks, each new file stays within 400 lines, and no budget or accepted debt grows.

Keep all existing assertions and independently authored expectations effective.
Do not rewrite old tests to match the new owner or derive expected results from the candidate producer.
For each required omission, pin the landed source and exact diagnostic before running it.
Accept only a compiling behavioral red, byte-identical restoration, and the same focused green afterward.
Invalid, compilation-only, skipped, or restore-failed proof closes no obligation.

Every first-use caller, fixture, registry, and headroom obligation closes in its introducing checkpoint.
A final family audit cannot supply a missing earlier test or repair.
Co-owned registry files may stay unchanged when their existing bindings and assertions suffice.
No new scanner, injected port, count expectation, structure grant, or policy authority enters this slice.

## Acceptance

- [ ] Complete scans enumerate once and retain candidate and row order; Rows still returns Scan(background).Rows.
- [ ] Pre-cancellation starts no operation, and cancellation during enumeration keeps Total/Omitted unknown without publishing partial candidates.
- [ ] Cancellation between candidates keeps the completed prefix and known Total-Inspected arithmetic without starting the next candidate.
- [ ] A cancelled inspection contributes no rows or inspected count; a completed empty inspection still increments Inspected.
- [ ] Explicit post-inspection cancellation withholds that candidate, including simultaneous completion/cancellation at the accepted checkpoint.
- [ ] Existing TestScanWaitsForCancelledWorkerCleanup stays unchanged, and context-ignoring work keeps Scan pending until its controlled release.
- [ ] Late cancellation after final completion leaves the complete result unchanged; fixture waits assert early-return/deadlock rather than a product deadline.
- [ ] Active-context enumeration error retains the current complete result without a new error state.
- [ ] Absent/empty hook directories retain the pre-push candidate; static regular/symlink/missing/incomplete headers and FIFO rejection retain their rows.
- [ ] Claude/Codex wiring, all ScanResult fields, full/brief checked-in responses, action-help order, suppression, command errors, and incomplete exit 0 remain unchanged.
- [ ] Family S compares independently retained reference and candidate Command tuples over identical controlled events without elapsed-time equality.
- [ ] A renamed unused old pipeline makes the registered package-core-guard diagnostic red; a fixture-string or test goroutine does not.
- [ ] Optional linked surfaces pass; a present kit surface with a missing or unreadable owner receives the named diagnostic.
- [ ] The live-tree ownership witness executes through its registered binding and explicit classification.
- [ ] Omitting only the registered scan subcheck call compiles and reds TestGuardScanLoopOwnerBitesOnPrivateCopy; restoration returns the same focused green.
- [ ] Removing the post-inspection context checkpoint compiles and reds TestScanLoopPostInspectionCancellation before any interrupted rows can publish.
- [ ] Every required mutation pins source and diagnostic and proves byte-identical restoration plus the same focused green.
- [ ] The sole relocated orchestrator preserves bindings/order and any already landed projection subcheck; npm producer and independent formatter count stay in their original file.
- [ ] checks_test.go stays within 709 lines after either outcome order; new files stay within 400 lines and no grant or accepted debt grows.
- [ ] All 28 predicates and original assertions close with no unlanded projection API dependency, new hook execution, parallel scan, or hard filesystem deadline.
