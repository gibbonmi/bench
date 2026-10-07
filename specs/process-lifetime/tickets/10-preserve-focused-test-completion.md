# Preserve focused test decoding through bounded cancellation

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md, 08-preserve-release-step-evidence.md, 09-bound-preflight-and-scanner-shutdown.md
Writes: internal/conformance/injected_ports_registry_test.go, internal/testreport/cancel_test.go, internal/testreport/command.go, internal/testreport/lifetime_test.go (new), internal/testreport/selection.go, internal/testreport/selection_test.go
Covers: PL64, PL65, PL79, PL80, PL122

## What to build

A canceled focused test resolves a descriptor-retaining descendant without falsely completing its JSON report.
Migrate focused execution and package-list execution through Lifetime.
Use INT and the current builder grace for cancellation.
Preserve package selection, list attribution, and complete JSON validation after natural EOF.

Replace StdoutPipe coupling with the established owned file-pipe consumer.
Malformed JSON, forced closure, and missing terminal package events remain failures.
No caller reads mutable output before its worker completes.

Consume selected-binary, private-run, and cache descriptors from their landed owner tickets.
The focused operation consumes Selection.Close and private-run Close errors before returning success.
Preserve child status beside required cleanup failure.
Keep current command inputs and ordinary complete selection cleanup.

Before implementation starts, confirm complete accepted and landed durable-file-replacement delivery.
Its leaf, review-record caller, and native qualification must be complete.
Durable-caller-migration is not a prerequisite.
Use one retained integration source and a fresh author for this ticket.

Keep each caller's selected normal-exit policy.
A raw child status cannot certify required cleanup or publication.
Retain uncertain resources and every unresolved obligation.
Never signal from recovered records or add a blanket normal-exit kill.

This ticket belongs to PL-C4.
Finish the preceding chunk review before this chunk starts.
Its blockers include every prior overlapping Writes owner.
No successor caller must exist for this ticket's owned predicates.

## Acceptance

- [ ] PL64: Focused-test cancellation resolves a descriptor-retaining descendant.
- [ ] PL65: Malformed focused JSON remains a decode failure.
- [ ] PL79: Focused output requires terminal events for every package.
- [ ] PL80: Package-list cancellation keeps current interruption attribution.
- [ ] PL122: The focused run consumes selection Close failure.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/testreport`

Mutation witness: Treat forced EOF as a complete report, wait before decoder completion, or discard Selection.Close failure. The real report and retained executable witnesses must fail.
