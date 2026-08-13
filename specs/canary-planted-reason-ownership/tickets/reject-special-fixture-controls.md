# Reject special fixture controls before reading

Blocked by: none
Writes: `internal/canary`

## What to build

Classify fixture control records by metadata before reading them. Add bounded synthetic discovery coverage for missing files, dangling symlinks, and special files so a FIFO `EXPECT` refuses immediately instead of blocking or entering the accepted inventory.

## Acceptance

- [ ] (covers CI4) Missing, dangling, and special fixture control records fail closed by metadata before content is read; a FIFO `EXPECT` neither blocks nor counts.
