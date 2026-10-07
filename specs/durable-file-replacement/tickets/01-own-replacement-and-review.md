# Own durable replacement through review records

Blocked by: none
Writes: internal/durablefile/ (new), internal/reviewrecord/write.go, internal/reviewrecord/replacement_test.go (new)
Covers: D01, D02, D03, D04, D05, D06, D07, D08, D09, D10, D11, D12, D13, D14, D15, D19, D32

## What to build

Implement the accepted durablefile API and its classified error, then migrate the real review-record writer in this checkpoint.
The leaf owns one adjacent regular-file replacement and imports no caller policy.
Review-record mode, rendering, validation, and destination authority remain with that caller.
Its added directory synchronization is an explicit behavior change.
Keep every accepted operation stage, error cause, temporary cleanup, and post-rename visibility requirement.

Review chunk: D-A.
This complete prerequisite has no caller-migration predecessor.
Its accepted and landed spec supplies the leaf to process-lifetime PL-C2 and the separate migration successor.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Replace a real review record, then fail directory synchronization after rename.
The caller returns *durablefile.Error with Published true and the uncertainty phrase, while the complete new record remains visible.
Qualify file and directory synchronization natively on both Linux and macOS before D-A acceptance.

- [ ] Replacement of old bytes installs the exact supplied bytes. (D01).
- [ ] An empty payload creates a present regular file. (D02).
- [ ] The installed file has the requested permission bits. (D03).
- [ ] A concurrent reader observes only the old or new complete payload. (D04).
- [ ] The temporary regular file belongs to the destination parent. (D05).
- [ ] Rename follows the successful mode, write, file-sync, and file-close operations. (D06).
- [ ] Success follows parent-directory synchronization and close. (D07).
- [ ] A short write returns an error before rename. (D08).
- [ ] Each pre-rename failure leaves the existing destination bytes unchanged. (D09).
- [ ] Each pre-rename failure leaves an absent destination absent. (D10).
- [ ] Each post-rename failure reports that new bytes may already be visible. (D11).
- [ ] A post-rename sync failure leaves the new destination bytes visible. (D12).
- [ ] Replacement errors preserve errors.Is access to the injected cause. (D13).
- [ ] A pre-rename failure removes its temporary file when removal succeeds. (D14).
- [ ] A failed temporary cleanup preserves the primary error and the cleanup cause. (D15).
- [ ] The review record caller exposes an injected directory-sync failure. (D19).
- [ ] Native Linux and macOS qualification establishes the declared synchronization result. (D32).

## Checkpoint verification

- `bench test --package ./internal/durablefile`
- `bench test --package ./internal/reviewrecord`
- `bench test --package ./internal/durablefile --run 'TestNativeReplacement'`
- `bench test --package ./internal/durablefile --run 'TestNativeReplacement'`

## Required omission evidence

Omit file sync, omit directory sync, rename before file close, accept a short write, or restore the former review writer.
Run owner omissions with the durablefile package command.
Run the review-writer bypass with the reviewrecord package command.
Each targeted assertion must report its named red, then pass after restoration.
Owner preservation and real review-record preservation have separate package commands and independent witnesses.
Record the host and supported filesystem for each native requirement.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
