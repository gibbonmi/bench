# Retire the Assumptions field from the spec-build machinery

Blocked by: none
Ownership fence: `internal/specbuild`, `internal/contract/runtime/runtime_spec_build_test.go`
Contracts: the checkpoint receipt shape crosses `internal/specbuild`→the runtime fixture in `internal/contract/runtime/runtime_spec_build_test.go`, asserted by RA3 against receipts the real built binary emits, whose expectations move in this same change

## What to build

`ParseTicket` skips a legacy `Assumptions:` line instead of parsing it; the
ticket, assignment, checkpoint, and integration structures drop the field, its
digests (`assumptionDigests`), and its comparisons in the checkpoint and
integrate paths. The digest duplicated a fact `TicketDigest` already seals —
the comma-splitting parser never produced the per-clause structure the field
advertised — so this is the removal of a second derivation, not of a control.

Legacy artifacts stay readable: a staged ticket carrying the old line still
parses, assigns, checkpoints, and integrates with the line ignored, and a
persisted run record or receipt carrying assumption data still loads because
unknown JSON fields are ignored. No migration pass, no version bump. The
runtime fixture's receipt constructions and expectations move in this same
change, because a receipt asserted against removed machinery is unlandable.

## Acceptance

- [ ] [RA1] a ticket file carrying a legacy `Assumptions:` line parses, assigns, checkpoints, and integrates with the line ignored end to end.
- [ ] [RA2] a pre-retirement persisted record carrying assumption data reloads in a fresh process and its lifecycle continues.
- [ ] [RA3] the runtime spec-build receipts construct no assumption digests, and the fixture's receipt expectations agree.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RA1 | make `ParseTicket` refuse a line it does not recognize | the legacy-tolerance test | refuse on the legacy prefix, run `go test ./internal/specbuild -run Legacy -timeout 180s`, expect a nonzero matched-test count and the parses-clean assertion to fail |
| RA2 | reject unknown JSON fields when loading records | the fresh-process reload test | switch the decoder to refuse unknown fields, run `go test ./internal/specbuild -run Reload -timeout 180s`, expect a nonzero matched-test count and the lifecycle-continues assertion to fail |
| RA3 | leave the assumption digests in the checkpoint receipt | the runtime receipt fixture | restore the digest construction, run `go test ./internal/contract/runtime -run TestRuntimeSpecBuild -timeout 300s`, expect a nonzero matched-test count and the receipt-shape assertion to fail |
