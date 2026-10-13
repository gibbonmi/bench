# Migrate prospective owner-record replacement

Blocked by: none
Writes: internal/gate/prospectiveartifact/prospectiveartifact.go, internal/gate/prospectiveartifact/replacement_test.go (new), internal/durablefile/durablefiletest/assert.go (new), internal/gate/durable_replacement_test.go
Covers: D18

## What to build

Use the accepted leaf for the prospective owner record.
Keep RecordMode, encoded newline, destination authorization, and existing refusal policy.
The caller test also proves one pre-rename failure, so the pre-rename refusal evidence of D36 exists for this caller.
One exported assertion in durablefiletest owns the published-failure and the record-file checks, and the gate caller test uses it too.
Keep the exported Publish signature; it supplies the native replacer to an unexported form in the same file, and the test injects through that form.

Review chunk: D-B2.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Publish the real record with a directory-sync fault after rename.
The caller returns classified failure and read-back observes the complete new record.

- [x] The prospective owner record caller exposes an injected directory-sync failure. (D18).

## Checkpoint verification

- `bench test --package ./internal/gate/prospectiveartifact`

## Required omission evidence

Restore the former local replacement sequence.
The real owner-record junction test must detect the missing shared failure.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
