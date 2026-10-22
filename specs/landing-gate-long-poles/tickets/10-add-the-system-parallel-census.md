# Add the system parallel census

Blocked by: 09-move-the-exec-child-environment-tests.md
Writes: internal/systemtest/parallel_census_test.go (new), internal/systemtest/serial_table_test.go (new), internal/systemtest/owner_test.go, internal/systemtest/land_route_test.go, internal/testcensus/census.go (new), internal/testcensus/census_test.go (new), internal/worktree/parallel_census_test.go, internal/systemtest/parallel_census_cases_test.go (new)
Covers: LGP25, LGP26, LGP28, LGP29, LGP31, LGP34

## What to build

The system package gains a census that classifies each top-level test as parallel or serial with a reason. It also gains one journey-repository helper. No test moves yet.

Create `internal/systemtest/parallel_census_test.go` with the `//go:build system` tag. The census refuses a test that neither calls `t.Parallel` nor appears in the serial table with a reason. The census's own tests call `t.Parallel`. `TestMain` is not a classified test, so 105 of the 106 top-level functions are classified.

The census walk has one source. Extract the pure walk of the worktree census into a new non-test package, `internal/testcensus`: the test-file parser (`internal/worktree/parallel_census_test.go:52`), `isTopLevelTest` (`:95`), and `callsParallelDirectly` (`:292`). The worktree census imports it in place of its own copies, and the system census imports it too. Keep the worktree serial helpers and `plantTestFiles` (`:455`) where they are, or move `plantTestFiles` as a shared helper. If a conformance check requires a test for the new package, give it one.

The tests are `TestSystemParallelCensus` (LGP25, LGP29), `TestLedgerWritesHoldTheOwnerLock` (LGP28), and `TestSystemSerialTableIsExact` (LGP31).

The census also refuses a parallel test that reaches `owner.repos`, or that names a ledger field outside an owner writer. An owner writer writes at least one ledger field, and every ledger-field reference in it sits between a `Lock` and an `Unlock` of the owner mutex. The owner value is spelled `owner` or `o`: six writers lock through `o.mu` (for example `owner_test.go:158-160`), and two through `owner.mu`.

Match selectors on the owner value, such as `owner.starts` or `o.starts`, not bare words. A local `seen` at `charge_evidence_stateless_test.go:57` and the strings in `TestUnfilteredRunVerifiesLedger` (`owner_selection_test.go:80-86`) are not ledger references.

`TestLedgerWritesHoldTheOwnerLock` refuses an unguarded ledger write by name. It exempts `TestOwnerSelectionChild`, which plants the negative controls in a child process.

Put the serial table in the new `internal/systemtest/serial_table_test.go`, so that neither file passes 400 lines. A new file over its limit reds the growth ratchet, because its base count is zero.

The table holds a temporary reason, "pending migration", for each of the 103 tests that later tickets move. The two subject tests hold the owner-repository reason. Tickets 11 to 15 remove the temporary rows, and ticket 16 forbids the reason.

Prove each refusal on planted fixture sources, as the worktree census tests do. Also run the census on the live package.

Add the journey-repository helper to `owner_test.go`. `newSystemOwner` builds its three repositories through it, and `landRouteRepo` (now in `land_route_fixture_test.go`) wraps it. `newSystemOwner` returns errors and runs before the global `owner` is set (`owner_test.go:56-61`), so the helper returns `(string, error)`, and a test wrapper such as `landRouteRepo` turns the error into `t.Fatal`.

The helper must not name `repos`. `newSystemOwner` keeps the append (`owner_test.go:112`). Otherwise every migrated caller reaches `owner.repos` through the helper, and the census refuses it.

Keep the `for range 3` loop, `len(o.repos) != 3`, and both `strippedJourneyMarker` occurrences in `owner_test.go`. `TestBranchNativeArchitectureCensus` (`internal/conformance/ordinary_build_census_test.go:88-99`, `:304`) pins them, and it is the one source of LGP26. The new census does not parse those facts again.

`owner_test.go` has 485 lines against its grant of 520. `land_route_test.go` has 508 lines and no grant, so it may not gain net lines. `internal/worktree/parallel_census_test.go` has no grant, and the extraction only removes lines from it.

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.

## Acceptance

- [ ] A planted test with no `t.Parallel` call and no serial row reds by name (LGP25).
- [ ] `TestBranchNativeArchitectureCensus` passes with the literal loop and the repository check kept (LGP26).
- [ ] The worktree census passes on its live tree through the shared `internal/testcensus` walk.
- [ ] A planted unguarded ledger write reds by name, and the live package passes (LGP28).
- [ ] A planted parallel test that reaches `owner.repos` or a ledger field reds (LGP29).
- [ ] A serial row for a missing test, or a row with no reason, reds (LGP31).
- [ ] The probe "drop the serial reason of an owner-repository mutator" bites.
- [ ] The verification record holds the system phase time before and after the change (LGP34).
