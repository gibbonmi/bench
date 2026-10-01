# Make the worktree test-count pin catch a real removal

Blocked by: none
Writes: internal/worktree/parallel_census_test.go
Covers: none

## What to build

The quality survey of 2026-09-29, card 06, found a weak pin in `internal/worktree/parallel_census_test.go`. `TestPackageTestCountPin` counts the top-level tests in the package test files. It refuses a count below `worktreeTestFloor`, which is 334. The package now declares 662 top-level tests. Thus an author can remove 328 tests before the pin turns red. The comment on the constant says that a removal below the pin turns the gate red, but a removal of one test passes.

The count cannot come from the live tree, because a removal changes the live count and the expectation together. Thus the honest guard is an exact pinned number. Set the pin to the true count, 662, and make the test require an exact match.

A removal or a merge of one test turns the pin red. An addition also turns the pin red, so the author raises the pin on purpose in the same change. Thus the pin does not drift below the true count again.

Keep the count itself as it is: every top-level test in the package test files. Make the two failure messages different, so a red names its cause: a count below the pin, or a count above the pin. Change the comment so that it states what the guard proves.

On 2026-10-01, the reviewer chose the exact pin over a raised floor. A floor at the true count goes silent again after the next addition. The cost of the exact pin is a constant bump for each added test. When two sibling landings both add tests, the second landing must merge the constant again.

The below-pin message names every cause of a lower count. A test was removed, merged, renamed off the `Test` prefix, or moved to a different package.

## Acceptance

- [ ] The pin is 662, and the test requires a count equal to the pin.
- [ ] A count below the pin and a count above the pin each give their own message.
- [ ] The comment on the pin states that one removal and one addition each turn the test red.
- [ ] At the base, a probe that omits `TestPackageClausePin` is silent for `TestPackageTestCountPin`. Observed: the test passed at the base.
- [ ] After the change, the same probe turns `TestPackageTestCountPin` red, and the probe restores the file. Observed: the test failed with "the package declares 661 top-level tests, below the pin of 662".
- [ ] After the change, a probe that adds one test turns `TestPackageTestCountPin` red, and the probe restores the file. Observed: the test failed with "the package declares 663 top-level tests, above the pin of 662".
- [ ] `go vet ./...` passes.
- [ ] `bench test --changed` passes for the changed packages.
