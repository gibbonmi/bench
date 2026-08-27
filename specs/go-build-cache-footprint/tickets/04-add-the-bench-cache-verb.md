# Add the bench cache verb

Blocked by: 02-own-the-go-build-cache-path.md, 07-measure-the-footprint-after-trimpath.md

Writes: internal/gocache/ (new), cmd/bench/main.go, cmd/bench/command_registry_test.go, internal/conformance/subcommand_routing_test.go, projects/benchkit.md, CHANGELOG.md

## What to build

An operator reads the cache size from one command. Add the read-only
`bench cache` verb. This ticket adds no gate line, no log event, and no
`clean` child.

Extend `internal/gocache/` with the footprint walk, the last trim time, and the
bound. The walk uses `lstat` only, opens no file, and follows no symlink. It
recurses into a `-d` directory, because Go stores a cached executable as a
directory. It counts regular files and sums their sizes. The last trim time
comes from `trim.txt` in the directory as unix seconds, after an `lstat`
regular-file check. The value renders as UTC RFC 3339, and it is empty when the
file is absent, not regular, or unparsable.

The bound is one constant of 10,737,418,240 bytes in the new module, and there
is no knob.

`bench cache` prints one `go_build_cache` table with the columns `dir`,
`bytes`, `files`, `last_trim`, `bound`, and `over_bound`. It reports zeros and
exits 0 on an absent or empty directory. It works outside a git repository. It
refuses with a named reason when the directory path holds a control byte.

Add the registry row in `cmd/bench/main.go` beside the `structure` row at line
89, with the same `axiExempt(axiReasonOperational)` marking. `bench cache` is
an operational report, so it joins no AXI query registry and needs no
agent-facing table edit. Add the routing map entry in
`internal/conformance/subcommand_routing_test.go`, and note the new seam in
`projects/benchkit.md`.

This ticket registers the `cache` parent verb.
`05-hold-the-cache-lock-and-add-bench-cache-clean.md` adds the `clean` child to
that verb, and `06-print-the-gate-line-and-log-the-footprint.md` reuses this
walk at the gate tail.

## Acceptance

- [ ] R01 — `bench cache` prints one `go_build_cache` table with `dir`, `bytes`, `files`, `last_trim`, `bound`, and `over_bound`.
- [ ] R02 — `bytes` equals the size sum of the regular files under the directory, with a `-d` directory recursed.
- [ ] R03 — `bench cache` on an absent directory prints zero bytes and zero files with exit 0.
- [ ] R04 — `bench cache` on an empty directory prints a zero-byte, zero-file row with exit 0.
- [ ] R05 — `bench cache` run in a directory outside any git repository exits 0.
- [ ] R06 — `last_trim` renders the `trim.txt` unix seconds as UTC RFC 3339.
- [ ] R07 — `last_trim` is empty when `trim.txt` is absent, a symlink, or a FIFO.
- [ ] R14 — A directory path with an ESC byte makes `bench cache` refuse with a named reason.
- [ ] R16 — The `cache` verb routes to the new module in the subcommand routing map.

Delivered outcome: an operator reads the Bench build cache size, its file
count, and its last trim from one command, anywhere on the machine.
