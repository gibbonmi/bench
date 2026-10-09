# Add the probe cost row to the known issues

Blocked by: none
Writes: docs/known-issues.md
Covers: none

## What to build

A coordinator who writes an author charge reads one row that names the cost shape of `bench probe` and its workaround.
`bench probe` links the package test binary twice for each call, once for the baseline run and once for the mutant.
In a large package such as `internal/adopt`, each call costs about one minute.
A charge that asks for one probe per independent expectation therefore multiplies that minute by the expectation count.
The FT217 ticket 01 author ran 54 serial probes after it finished its code in 13 minutes.

Add one row to the table in `docs/known-issues.md`.
The issue cell states the double link and the per-expectation multiplication.
The retire condition cell names FT168: one probe takes several mutations in one private build.

The workaround cell states three things.
First, an author charge caps witnesses at one per coverage row plus one per production branch the ticket adds.
Second, the author writes fewer independent expectations, for example one table-driven test whose rows share one derivation.
The one-source rule in `AGENTS.md` still owes each independent expectation a demonstrated red.
Third, the author runs the probes after the floor checks, so an invalid mutant does not block the commit.

Keep the existing rows unchanged.
Use ASD-STE100 prose and keep each sentence inside the 25-word bound that `bench gate-prose` grades.

## Acceptance

- [ ] `docs/known-issues.md` holds one new row whose issue cell names the double link of `bench probe` and the per-expectation multiplication.
- [ ] The workaround cell names the witness cap, the fewer-expectations condition under the one-source rule, and the floor-checks-first order.
- [ ] The retire condition cell names FT168 and its multi-mutation form.
- [ ] `bench gate-prose . -- docs/known-issues.md` passes, and `bench test --check docs-currency-workflow` passes.
