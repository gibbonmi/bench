# Make mutation inventories optional by default

Blocked by: none
Ownership fence: `internal/specbuild/assign.go`, `internal/specbuild/assign_covers_test.go`, `internal/specbuild/refresh.go`, `internal/specbuild/refresh_repin_test.go`, `scripts/gremlins-diff.sh`, `CHANGELOG.md`
Integration surfaces: parsed presence of optional `Closure:` and `## Red mutations` fields→`requireClosure`; optional legacy graph values→existing bidirectional validator; assignment and refresh revalidation→same `requireClosure` owner; permanent behavior→`scripts/gremlins-diff.sh` and `CHANGELOG.md`
Contracts: four parser states cross `internal/specbuild/assign.go`→`internal/specbuild/refresh.go` unchanged—neither field succeeds, Closure-only refuses, heading-only refuses, and both declared invoke the complete existing validator—asserted by PO1 and PO2; the permanent no-flag behavior crosses `internal/specbuild/assign.go`→`scripts/gremlins-diff.sh` and `CHANGELOG.md`, asserted by PO1

## What to build

Retire `BENCH_RED_MUTATIONS_OPTIONAL` as an experiment and make its safe absent-pair
behavior permanent. A modern ticket declaring neither optional inventory assigns
and refreshes by default. Extend the parsed ticket model with explicit
field-presence bits so an empty declared heading cannot masquerade as absence:
neither field succeeds; Closure-only and heading-only refuse; both declared run
the existing validator, where open, malformed, duplicated, or non-reciprocal
graphs still refuse exactly as today. Exercise both initial assignment and
refreshed-assignment revalidation. Remove flag-dependent branches/tests and
update the changed-package referee and changelog to describe the permanent
default. Do not remove parser compatibility for already-landed tickets that
retain the fields.

## Acceptance

- [ ] [PO1] (covers local) modern tickets with neither optional field assign and refresh with no environment flag, ordinary acceptance/fence/contracts/covers validation remains unchanged, and the referee/changelog describe the permanent no-flag default.
- [ ] [PO2] (covers local) explicit parser presence distinguishes all four optional-field states: Closure-only and heading-only refuse, both declared enter the existing complete bidirectional validator, and every malformed, duplicated, non-reciprocal, or open declared graph remains red.
