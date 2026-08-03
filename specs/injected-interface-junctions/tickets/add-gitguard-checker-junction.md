# Add the gitguard real-checker junction tests

Blocked by: none
Ownership fence: `internal/gitguard/checker_junction_test.go`
Contracts: the two probe booleans with opposite fail-safe polarities cross `internal/gitguard`→`internal/git` through the real `git.RefResolves`/`git.BranchExists`, asserted by GC1 and GC2 against the real probes
Assumptions: production is presumed correct — a red against unmodified production is a stop-and-surface finding; the probes resolve in the process working directory so the tests use `t.Chdir` and must not run parallel; the timeout test plants a PATH-front stub `git` sleeping past the 2s bound and tolerates its wall-clock cost; claims re-derived from the tree at pickup

## What to build

`Classify` composed with the real `Checker{git.RefResolves, git.BranchExists}`:
resolved answers over a temp repo (block checkout of an unresolvable ref,
permit a resolvable one, and the forced-creation clobber verdicts), plus the
timeout composition where a hung git makes both polarities land on "block".

## Acceptance

- [ ] [GC1] `Classify` with the real probes over a temp repo (one commit, one branch) blocks checkout of an unresolvable ref, permits a resolvable one, and blocks forced creation clobbering the existing branch.
- [ ] [GC2] `Classify` with a PATH-front stub `git` sleeping past the probe bound blocks both the checkout verdict (ref presumed unresolvable) and the forced-creation verdict (branch presumed present).

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GC1 | point the fixture's resolvable-ref case at a ref the temp repo does not have | the junction test | apply the mutation, run `go test ./internal/gitguard -run <GC1 test>`, expect the permit assertion to fail |
| GC2 | shorten the stub's sleep below the probe bound | the junction test | apply the mutation, run the GC2 test, expect the timeout-polarity assertions to fail because the probes now resolve |
