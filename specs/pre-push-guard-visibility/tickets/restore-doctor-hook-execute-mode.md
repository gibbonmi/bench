# Restore doctor hook execute mode

Blocked by: none
Ownership fence: `internal/adopt/doctor.go`, `internal/contract/surface/doctor_test.go`
Integration surfaces: stale-hook repair route→`internal/adopt/doctor.go` + DEX1; doctor runtime contract→`internal/contract/surface/doctor_test.go` + DEX1
Contracts: a stale managed hook-health record crosses `internal/adopt/link_hook.go`→`internal/adopt/doctor.go`, asserted by DEX1 against the real doctor repair path
Closure: DEX1/stale-existing-execute-mode

## What to build

After `bench doctor --fix` rewrites an existing stale managed pre-push hook whose mode lacks execute bits, restore executable mode before reporting repair success. Keep the rewrite and its post-write mode restoration in this one vertical cut: a test-only or chmod-only cut strands the doctor runtime contract red because the public repair operation must both render current bytes and leave an executable guard.

## Acceptance

- [ ] [DEX1] (covers local) `bench doctor --fix` repairs a stale managed hook seeded at mode `0644` to current bytes with one or more executable permission bits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DEX1/stale-existing-execute-mode | omit the post-repair chmod of the existing hook | doctor runtime contract | seed stale managed bytes at `0644`, run `bench doctor --fix`, stat the hook, expect the executable-mode assertion red |
