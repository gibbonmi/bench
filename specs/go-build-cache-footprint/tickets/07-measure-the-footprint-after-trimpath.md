# Measure the footprint after -trimpath

Blocked by: 02-own-the-go-build-cache-path.md, 03-add-trimpath-to-every-bench-go-argv.md

Writes: specs/go-build-cache-footprint/tickets/07-measure-the-footprint-after-trimpath.md

## What to build

This ticket carries no code and no coverage row. The coordinating session runs
the measurement itself and writes the results into the checklist below. The
spec marks story 33 as not covered by a test, and this checklist is its record.

The measurement repeats the before-and-after pair against the Bench build
cache. Measure the directory with `du -sb` before and after one `bench gate` on
a fresh worktree path. Then measure it before and after one `bench gate --fresh`
on that same path. The `du -sb` apparent-byte figure is the acceptance number,
and the `bench cache` `bytes` column measures the same quantity.

The before numbers stand in the spec's Further notes, measured on 2026-08-27
against the interim directory. The cache held 368,803,070 bytes in 6,884 files.
One gate on a fresh worktree path grew it to 567,642,248 bytes and 9,505 files.
One `--fresh` gate on the same path grew it to 567,946,742 bytes and 9,578
files.

This ticket is the build checkpoint. Tickets 01 to 03 carry the measured saving
on their own. After this measurement the reviewer decides whether tickets 04 to
06 build, because requirement 4 prefers the measurement before the bound's
shape is fixed.

## Acceptance

- [ ] The after numbers for one `bench gate` on a fresh worktree path are recorded here as a byte total and a file count.
- [ ] The after numbers for one `bench gate --fresh` on that same path are recorded here as a byte total and a file count.
- [ ] The reviewer's decision on tickets 04 to 06 is recorded here.

Delivered outcome: the reviewer reads the measured saving beside the before
numbers and decides whether the bound's tickets build.
