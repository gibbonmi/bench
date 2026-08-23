# Review pickup — worktree-test-latency

Frozen base `4a8aa16a`, reviewed tip `e26be7cc`. Initial three-axis review.

## Standards

Count: 1. Worst issue: the pasted purity census harness.

- **Duplicated census harness.** The three `purity_census_test.go` files share
  one scan algorithm; only the package line and one import row differ. The
  one-source-per-fact standard in `AGENTS.md` names a pasted fixture harness
  as a defect. The precedent shape is a shared test-support package, as in
  `internal/gittest` and `internal/axi/axitest`. A shared function keeps the
  per-package fence, because `go test` sets the working directory to the
  package under test. Disposition: **ask-user** — the extraction adds a new
  shared test-support seam, and the reviewer owns new seams.

## Spec

Count: 0. Clean. The per-row audit found all 13 rows covered.

## Coverage

Count: 2. Worst issue: lost precedence coverage in the explicit decision table.

- **Lost precedence coverage.** The deleted `TestExplicitEligibilityOutcomeMatrix`
  pinned cross-block precedence with competing facts (its EX3, EX5, and EX6
  subtests). Every case of `TestExplicitDecisionTable` in
  `internal/worktree/lifecyclepolicy/lifecyclepolicy_test.go` sets one fact,
  so a block reorder in `DecideExplicit` passes every current test. The
  automatic table already carries combined-fact cases; the explicit table
  does not. Disposition: **auto-fix** — add combined-fact cases that
  reproduce the EX3, EX5, and EX6 precedence verdicts.
- **Ledger wording on symlink shapes.** The ledger labels the deleted
  symlink test "adapter", but the adapter test covers only the key-level
  symlink with a real link. The shared `lstatShape` helper covers the other
  two levels through the same code path. Disposition: **no-op**.
