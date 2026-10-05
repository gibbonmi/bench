# Enforce completed fixture ownership

Blocked by: 12-migrate-capability-conformance-callers.md
Writes: internal/conformance/checks_test.go, internal/conformance/git_fixture_owner_test.go (new), internal/conformance/git_plumbing_owner_test.go, internal/conformance/ordinary_build_census_test.go, internal/conformance/registry/checks.go, internal/conformance/skip_ownership_test.go, internal/conformance/skip_reason_owner_test.go (new), internal/conformance/testdata/fixture-skip-classes.json (new), internal/conformance/tier_test.go
Covers: GF10, GF24, GF25, GF26, GF27, GF28, GF31, GF30

## What to build

Extend the existing Git and skip ownership checks after all migrations pass. Reuse moduleGoFiles and existing AST helpers. Preserve the production Git administration rule and its negative controls. Record exact specialized-site exceptions without granting file-wide exemptions. Add the independent caller-class fixture and demonstrate its named class-change red.

Plant one generic Git runner outside gittest. The registered ownership check fails. Restore it, then repeat with a duplicate identity and an unconverted capability caller.

## Acceptance

- [ ] Every former testrepo Git export caller reaches gittest (GF10).
- [ ] A planted generic Git runner outside gittest fails the registered check (GF24).
- [ ] A planted default identity outside gittest fails the registered check (GF25).
- [ ] A planted direct FIFO or symlink skip caller fails the registered check (GF26).
- [ ] Source examples in literals do not trigger ownership diagnostics (GF27).
- [ ] A specialized-site exception does not permit a generic sibling wrapper (GF28).
- [ ] Every migrated capability caller retains its baseline class (GF31).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --check git-plumbing-owner`
- `bench test --check skip-ownership`
- `bench test --check ordinary-build-census`
- `bench test --package ./internal/conformance`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
This ticket owns the final package invariants and all specialized-site exceptions. No generic copy may remain outside its owner.
