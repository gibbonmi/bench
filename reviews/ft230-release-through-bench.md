# Review pickup — ft230-release-through-bench

Base `1766e4d1`, tip `2277209b`. Raw findings: 3 (Standards 0, Spec 0,
Coverage 3). De-duplicated repair targets: 1 — findings 1 and 2 share one fix
(a failing stub-npm path through the command surface).

## Standards

Count: 0. No worst issue.

## Spec

Count: 0. No worst issue.

## Coverage

Count: 3. Worst issue: finding 2 — the resumability guarantee has no
adapter-level proof.

1. **auto-fix** — No test isolates an absent `npm` binary with `--adapter
   npm`. The spec's edge inventory decides the behavior ("the adapter's
   structured `npm ... failed` error surfaces as unsatisfied release intent,
   exit 1") and no test in `internal/publication/command_adapter_test.go` or
   `npm_registry_test.go` asserts it.
2. **auto-fix** — The stub `npm` in `command_adapter_test.go` always exits 0,
   so a non-zero npm exit on an interior package (partial publication) is
   never exercised through the npm adapter: no assertion that the record
   marks the published prefix, surfaces the failure as exit 1, and resumes on
   a re-run. The spec's edge inventory owns interruption at the state-machine
   seam; the adapter-path proof is the gap.
3. **no-op** — `--profile bank --adapter npm` leaves `Access` empty and is
   untested. Refuted as a defect: the spec's implementation decision sets
   `Access: "public"` for the public profile only, and no acceptance row maps
   the bank profile onto the npm adapter.
