# 6. Print the evidence manifest summary by default

Blocked by: none
Writes: internal/preflight/evidencecmd/, internal/chargeevidence/schema.go, internal/chargeevidence/read.go
Covers: BO42, BO43, BO44

## What to build

`bench preflight evidence <id>`, with no cursor and no source, prints one `evidence_summary` block. The block holds the evidence identity, the source count, the page count, the manifest bytes, the total source bytes, and `next`. `next` is the exact command that reads the first manifest page. Register the block in the chargeevidence response schema, so the ADR 0022 encoded-response bound applies to it.

The cursor, source, verify, and check-current forms keep their current output.

## Acceptance

- [ ] `bench preflight evidence <id>` prints one `evidence_summary` block with no content cell.
- [ ] The summary's `next` command returns the first manifest page byte-for-byte as the old default did.
- [ ] The existing cursor, source, verify, and check-current tests pass unchanged.
