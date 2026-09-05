# Repair the review findings

Blocked by: add-the-structure-growth-mode.md, run-the-growth-check-in-the-fast-lane.md, state-the-tickets-only-rule-once.md, construct-the-lane-owner-once.md
Writes: internal/structure/structure.go, internal/structure/budgets.go, internal/structure/growth_test.go, internal/landing/close.go, internal/landing/close_test.go, internal/landing/lane_test.go
Covers: SR62, SR63

## What to build

The line is opus / medium. Repair the five code targets the review pickup
`reviews/structural-refactor-pass.md` lists as R1 to R5, per decision (w).

R1. The growth query compares the base commit against the working tree, not
against HEAD. The lane's private checkout keeps the repository's HEAD, which both
callers pass as the base, and holds the composed tree in its working tree. A query over `base..HEAD` is empty
there, so the check never fires. Keep the NUL framing and the exact-rename
pairing. An empty `--growth` value refuses at exit 2 through the grammar.

R2. One predicate owns the over-budget rule for the all scan and the growth
mode: the newline count of a regular working-tree file, the compare against
the budget, the accept exemption, and the loud unreadable-accept line. Both
callers call it, and the package doc no longer claims the duplication ended.

R3. A call that names both `--since` and `--growth` refuses at exit 2 with the
usage line and one reason. The shape is the one the commit grammar uses.

R4. The lane test comment in `internal/landing/lane_test.go` states the
contract in timeless form, with no reference to what callers used to do.

R5. The git-object reader answers "directory" only for a tree object at the
folder path, so a committed blob at that path is never tickets-only.

## Acceptance

- [ ] `Growth` over a detached checkout whose HEAD equals the base prints the `FILE GREW` row at exit 1 when the read-in tree grows an over-budget file.
- [ ] A live `bench commit --dry-run` in a worktree with a planted growth in an over-budget Go file fails the lane with `check=structure` and the `FILE GREW` line.
- [ ] `bench structure --since <base> --growth <base>` refuses at exit 2 with the usage line and one reason.
- [ ] `bench structure --growth ""` refuses at exit 2.
- [ ] The all scan and the growth mode call one over-budget predicate, and the ten growth tests and the four pre-existing structure tests pass unchanged.
- [ ] `TicketsOnly` through the git-object reader answers false for a committed blob at the folder path.
- [ ] Self-probe: restore the `..HEAD` suffix in the growth query, and report the observed red.
