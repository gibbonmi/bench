# Print shift help before the adapter check

Blocked by: none
Writes: internal/shift/shift.go, internal/shift/shift_test.go

## What to build

`bench shift --help` prints its public grammar before it validates `BENCH_AGENT`.

## Acceptance

- [ ] `bench shift --help` exits 0 when `BENCH_AGENT` is unset.
- [ ] `bench shift -h` prints the same grammar and exits 0.
- [ ] An ordinary shift still rejects an unset adapter with exit 2.
