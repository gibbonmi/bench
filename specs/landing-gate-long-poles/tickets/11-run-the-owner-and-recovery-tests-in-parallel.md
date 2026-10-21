# Run the owner and recovery tests in parallel

Blocked by: 10-add-the-system-parallel-census.md
Writes: internal/systemtest/owner_artifact_recovery_test.go, internal/systemtest/owner_artifact_children_test.go (new), internal/systemtest/owner_test.go, internal/systemtest/owner_selection_test.go, internal/systemtest/owner_teardown_test.go, internal/systemtest/owner_land_race_test.go, internal/systemtest/owner_land_lease_test.go, internal/systemtest/owner_stale_seal_test.go, internal/systemtest/parallel_census_test.go (new), internal/systemtest/serial_table_test.go (new)
Covers: LGP27, LGP30, LGP34

## What to build

The owner and recovery tests run in parallel. The three recovery tests stop binding TMPDIR in the process.

Remove the three `t.Setenv("TMPDIR", ...)` calls (`internal/systemtest/owner_artifact_recovery_test.go:24`, `:280`, `:301`). Pass the private directory as an explicit child value instead.

`artifactLandEnv` (`:161-163`) and `startArtifactLand` (`:165-191`) take the private directory as an argument. The review children use `systemLandEnv` at `:29`, `:59`, and `:398`, and they get it too. Every child that writes a bundle gets it: the calls at `:29`, `:33`, `:59`, `:62`, `:93`, `:284`, `:305`, `:318`, `:398`, and `:401`. A missed child writes to the shared temporary root, and its bundle lookup fails.

The census gains the LGP30 rule: no `t.Setenv("TMPDIR"` call remains in the package.

These tests move to a private journey repository: `TestStrippedDistributionJourney` (`owner_test.go:449`), `TestRedTimeoutAndDescendantTeardown` (`owner_teardown_test.go:15`), and the tests at `owner_selection_test.go:227` and `:245`.

`TestSelectedExecutableComposition` and `TestOwnerSelectionChild` keep `owner.repos` and stay serial. The other six `owner_selection_test.go` tests, and the tests in `owner_land_race_test.go`, `owner_land_lease_test.go`, and `owner_stale_seal_test.go`, call `t.Parallel`.

Remove each moved test's "pending migration" row from the serial table. Run `bench test --check system` twice before you call a test safe in parallel. A test that fails under parallel load keeps a serial row, and its reason quotes the failing run.

`owner_artifact_recovery_test.go` has 430 lines and no grant, so it may not gain net lines. If it must grow, move its child-environment helpers to the new `owner_artifact_children_test.go`. Keep both `strippedJourneyMarker` occurrences in `owner_test.go`.

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.

## Acceptance

- [ ] No `t.Setenv("TMPDIR"` call remains, and a planted one reds the census by test name (LGP30).
- [ ] The unfiltered run still names each missing terminal outcome (LGP27).
- [ ] The three recovery tests pass in parallel in two runs, with the private directory on every bundle child.
- [ ] The probe "drop TMPDIR from the startArtifactLand child environment" bites.
- [ ] The verification record holds the system phase time before and after the change (LGP34).
