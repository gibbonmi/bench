# Review pickup — stable-owner-landing

Base `7eca97bf`, reviewed tip `4d6c9c55`. Raw findings: 9. Repair targets: 8
(the Spec finding and Coverage finding 4 share one site).

## Standards

Count: 2. Worst: the platform-suffix duplication.

- S1 — auto-fix. `CHANGELOG.md` new entry holds a 31-word sentence; the STE
  limit for a descriptive sentence is 25 words. Split it.
- S2 — accepted (reviewer): single-source the platform fact via the manifest.
  The Go installer writes the platform value; the wrapper drops its own
  `platform_suffix()` derivation for the broker check. The digest comparison
  remains the binding check. Cited: `internal/adopt/broker.go:290-298`,
  `bin/bench.sh` land route platform comparison.

## Spec

Count: 1. Worst: the SOL06 base binding.

- P1 — accepted (reviewer): bind `--base` to the assignment's recorded
  `Start`, keep the destination-ancestry check as a second guard, and add the
  mutation test that pins the exact refusal for a different-valid-ancestor
  base. Cited: `internal/worktree/land.go:99-102`,
  `internal/worktree/reauthorize.go:86-92` (the recorded-start precedent).
  This target also closes Coverage finding C4.

## Coverage

Count: 6. Worst: the untested manifest refusal branches.

- C1 — auto-fix. Test each fail-closed branch of the wrapper's land route:
  missing manifest, incomplete manifest, platform or version mismatch,
  non-regular or empty broker, digest mismatch. Site: `bin/bench.sh` land
  route; extend `internal/systemtest/land_route_test.go`.
- C2 — auto-fix. Add the empty-but-set case (`BENCH_KIT=""` and the other two
  variables) to the inherited-routing refusal test.
- C3 — parked (reviewer): crash-safe cleanup is a new seam. Parked to
  `capture/IDEAS.md`; not a repair in this build.
- C4 — folded into P1.
- C5 — auto-fix. Drive `prospectiveRunBinaryOwner` through a `Build` failure
  and assert a clean error with no residue. Site:
  `internal/gate/prospective.go:597`, `prospective_owner_test.go`.
- C6 — accepted (reviewer): amend the spec with a Won't-handle line for
  newline and control-byte roots. Spaces and glob characters stay covered by
  the existing fixtures.
