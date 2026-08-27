# Revert the interim machine setting

Blocked by: 03-add-trimpath-to-every-bench-go-argv.md, 05-hold-the-cache-lock-and-add-bench-cache-clean.md, 06-print-the-gate-line-and-log-the-footprint.md, 07-measure-the-footprint-after-trimpath.md

Writes: specs/go-build-cache-footprint/tickets/08-revert-the-interim-machine-setting.md

## What to build

This ticket carries no code and no coverage row. The coordinating session runs
the revert itself and writes the results into the checklist below. The spec
marks stories 34 and 35 as not covered by a test, and this checklist is their
record.

An interim `go env -w GOCACHE=/home/mgibs/.cache/bench-go-build` was set on this
box on 2026-08-27. It relocates the writes but bounds nothing, and it applies to
every Go workload, so Bench does not own that directory.

The revert runs once the code tickets are green on the integration source, and
before the landing. The landing's private exact-source binary then already
carries the Bench-owned path. Run `go env -u GOCACHE` on this box. Confirm that
`go env GOCACHE` prints `/home/mgibs/.cache/go-build`. Then delete the orphaned
directory `/home/mgibs/.cache/bench-go-build`.

## Acceptance

- [x] `go env -u GOCACHE` has run on this box, and the result is recorded here: it ran on 2026-08-27 with no output and exit 0.
- [x] `go env GOCACHE` prints `/home/mgibs/.cache/go-build`.
- [x] The directory `/home/mgibs/.cache/bench-go-build` no longer exists. It held 1.8 GB before the delete.

Delivered outcome: this box carries no unowned directory that looks
Bench-owned, and the handoff reader sees the revert done.
