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

- [x] The after numbers for one `bench gate` on a fresh worktree path are recorded here as a byte total and a file count.
- [x] The after numbers for one `bench gate --fresh` on that same path are recorded here as a byte total and a file count.
- [ ] The reviewer's decision on tickets 04 to 06 is recorded here.

## Measurement, 2026-08-27, source tip `62e0c703`

The Bench build cache is `/home/mgibs/.cache/bench/go-build`. Each figure is
`du -sb` apparent bytes and a regular-file count.

| step | before | after | growth |
|---|---|---|---|
| gate on the integration worktree | 375,986,901 bytes, 5,808 files | 376,606,010 bytes, 5,860 files | 619,109 bytes, 52 files |
| gate on a fresh detached checkout path | 376,606,010 bytes, 5,860 files | 376,864,492 bytes, 6,057 files | 258,482 bytes, 197 files |
| `--fresh` gate on that same path | 376,864,492 bytes, 6,057 files | 376,877,984 bytes, 6,095 files | 13,492 bytes, 38 files |

The spec's before pair on the interim directory is the comparison. The fresh
path grew it by 198,839,178 bytes and 2,621 files. The `--fresh` gate grew it
by 304,494 bytes and 73 files.
The fresh-path growth fell by a factor of about 770.

Delivered outcome: the reviewer reads the measured saving beside the before
numbers and decides whether the bound's tickets build.
