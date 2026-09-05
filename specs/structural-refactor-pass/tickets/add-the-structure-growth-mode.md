# Add the structure growth mode

Blocked by: declare-the-worktree-leaf-table.md, construct-the-lane-owner-once.md, move-the-kit-source-predicate-into-gate.md, own-the-operand-path-match-in-intent.md, state-the-tickets-only-rule-once.md, move-the-refresh-package-beside-its-consumers.md
Writes: internal/structure/structure.go, internal/structure/budgets.go (new), internal/structure/growth_test.go (new)
Covers: SR42, SR43, SR44, SR45, SR46, SR47, SR48, SR49, SR50, SR51, SR56, SR57

## What to build

The line is opus / medium. Add `bench structure --growth <base>`. The mode
lists the source files that changed between the base and HEAD, with NUL
framing and exact-rename pairing. It reads each tip count from the working
tree and each base count from the base commit's blob. A file absent at the
base counts zero, and an exact rename reads the old path's blob. A file reds
when its tip count exceeds both its limit and its base count.

The limit and the accept list come from the engine the all scan uses, per
decision (m). The limit is the `structure.budgets` row when one exists, else
`BENCH_MAX_LINES` or 400. A structure-accept row exempts the path.

A red prints one `FILE GREW` row per file with the tip count, the base count,
the limit, and the path. One summary line follows, and the command exits 1.
No red prints one ok line at exit 0. An unresolvable base prints `git diff
failed:` on stderr at exit 1, as `--since` does. A present-but-unreadable
accept file prints the loud accept line at exit 1.

Bare `bench structure` and `bench structure --since <base>` keep their bytes,
so `bench status` reads the same violation count. The structure usage line
gains the growth flag, per decision (n).

Split `internal/structure/structure.go`, per decision (o): the budget and
accept loaders move to `internal/structure/budgets.go`, which pairs with
`accept_validation.go`. The moved symbols are `loadBudgets`, `loadAccepts`,
`staleAcceptWarnings`, `budgetFor`, `envInt`, and `allDigits`. Every existing
test function keeps its name. The new growth tests land in the new file
`internal/structure/growth_test.go`, and they set `BENCH_MAX_LINES` through
the environment as the since-mode tests do.

This ticket's diff supplies the contract the lane-check ticket reads: the
growth flag spelling `--growth`, the base operand, and the `FILE GREW` row
label.

## Acceptance

- [ ] A grown over-budget file prints one `FILE GREW` row with the tip count, the base count, the limit, and the path.
- [ ] A red growth run prints one summary line after the rows and exits 1.
- [ ] An over-budget file that lost lines or held its count prints no row, and the command exits 0.
- [ ] A file at or under its limit that gained lines prints no row.
- [ ] A file absent at the base and over its limit at the tip prints a row with a base count of 0.
- [ ] An over-budget file with a structure-accept row that gained lines prints no row, and the command exits 0.
- [ ] A file with a `structure.budgets` row prints no row within that budget and prints a row with that budget as the limit past it.
- [ ] An exact rename of an over-budget file with no content change prints no row.
- [ ] A changed source file whose name carries a byte above ASCII is counted, and its row prints the name's own bytes.
- [ ] `--growth` with a present-but-unreadable accept file prints the loud accept line and exits 1.
- [ ] `--growth` with a base that does not resolve prints `git diff failed:` on stderr and exits 1.
- [ ] Bare `bench structure` and `bench structure --since <base>` print the same bytes as at the base commit over the fixture repository.
- [ ] The build records a red `FILE GREW` row against a planted growth in a throwaway repository, and the ok line after the revert.
- [ ] Self-probe: omit the base-count comparison, and report the observed red.
