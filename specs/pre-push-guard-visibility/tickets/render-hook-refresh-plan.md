# Render the hook refresh plan

Blocked by: none
Ownership fence: `internal/adopt/upgrade.go`, `internal/contract/surface/upgrade_test.go`
Integration surfaces: prospective hook refresh→`internal/adopt/upgrade.go` + U1/U2/U3; upgrade plan-mode contract→`internal/contract/surface/upgrade_test.go` + U1/U2/U3
Contracts: prospective hook bytes cross the upgrade transaction→`internal/adopt/upgrade.go`, asserted by U1 against the real upgrade plan
Closure: U1/different-hook-count, U1/current-hook-count, U2/equal-version-no-count, U3/check-no-write

## What to build

Render an upgrade plan that counts a real differing hook refresh, omits a current one, and never promises work on the equal-version early return or writes in `--check` mode. Keep the count and no-write cases together: a count-only cut strands the upgrade plan-mode contract red when its prospective hook inspection stages or writes the hook during `--check`.

## Acceptance

- [ ] [U1] `bench upgrade --check` counts a differing hook refresh and omits the count when the prospective hook already matches.
- [ ] [U2] Equal installed and linked versions count no hook refresh because that path performs none.
- [ ] [U3] `bench upgrade --check` performs no write.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| U1/different-hook-count | derive changed count only from plan entries | upgrade surface contract | use differing versions and stale hook bytes, run `--check`, expect the added hook change |
| U1/current-hook-count | count the hook unconditionally on a real upgrade | upgrade surface contract | repeat with current hook bytes, expect one fewer change |
| U2/equal-version-no-count | count hook drift before the equal-version return | upgrade surface contract | use equal versions with a stale hook, run `--check`, expect no hook refresh count |
| U3/check-no-write | stage or promote the hook during `--check` | existing upgrade plan-mode contract | snapshot the fixture, run `--check`, expect no filesystem change |
