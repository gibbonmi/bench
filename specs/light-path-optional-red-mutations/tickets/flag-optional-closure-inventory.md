# Flag the closure inventory optional for red-mutation-free runs

Blocked by: none
Ownership fence: `internal/specbuild/assign.go`, `internal/specbuild/assign_covers_test.go`, `CHANGELOG.md`
Integration surfaces: closure refusal→internal/specbuild/assign.go; refresh re-validation→existing requireClosure call + FL1; flag advertisement→CHANGELOG.md
Contracts: none crosses
Closure: FL1/absent-inventory-assigns, FL2/declared-graph-still-validates, FL3/default-still-refuses

## What to build

The reviewer can set `BENCH_RED_MUTATIONS_OPTIONAL=1` to run experimental builds
whose tickets omit the `Closure:` line and the `## Red mutations` section. With
the flag set, a modern ticket that declares neither assigns; a ticket that
declares either still validates the full closed graph, so the flag makes the
inventory optional, never invalid. Unset, enforcement is byte-identical to
today: fail closed is the default posture, and turning it off is a visible
per-run act, not a standing edit.

## Acceptance

- [ ] [FL1] with the flag set, assign leases a modern ticket that declares no Closure inventory and no Red mutations rows.
- [ ] [FL2] with the flag set, assign still refuses a declared closure graph that is open.
- [ ] [FL3] with the flag unset, assign refuses a modern ticket without a Closure inventory, unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FL1/absent-inventory-assigns | ignore the env flag and always enforce | the flag-on assign test | set the flag in-test, assign the inventory-free ticket, expect a lease not a refusal |
| FL2/declared-graph-still-validates | return nil from requireClosure whenever the flag is set | the flag-on open-graph test | set the flag in-test, assign a ticket whose fact has no mutation row, expect the existing refusal message |
| FL3/default-still-refuses | invert the flag comparison so unset skips enforcement | the existing missing-closure refusal test | run the closure refusal suite with no env set, expect the atomic-inventory refusal |
