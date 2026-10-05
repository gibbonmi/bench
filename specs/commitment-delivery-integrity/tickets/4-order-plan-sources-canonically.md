# 4. Order plan sources in one canonical order

Blocked by: 3-bind-only-open-plan-sources.md
Writes: internal/commitment/authority.go, internal/commitment/authority_test.go, internal/commitment/authority_internal_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: FD20, FD21, FD22

## What to build

This ticket is the second ticket of review chunk FD-C2. It orders the source
set that ticket 3 defines, and it delivers a plan identity that no traversal
change can move.

Add one unexported function, `boundSources`, to `authority.go`. It sorts a copy
of the plan sources once and returns the sorted list and the identity of its
encoding. The order compares the identifier, then the path, then the identity,
each as a Go string, which is byte order. `BuildPlan` stores the sorted list in
`Plan.Sources` and puts the returned identity into the bound plan input. The
inline encoding of the sources in `BuildPlan` goes away, so one function owns
the order and the encoding.

The readers of `Plan.Sources` and of the plan identity make no order
assumption, so none needs an edit. A receipt approved before this change can
fail to match at the commit. The commit then refuses with its existing plan
guidance. Each fixture binding in an outcome with sources lists an obligation,
because ticket 1 refuses a new obligation-free binding.

## Acceptance

- [ ] `TestBoundSourcesIgnoresInputOrder` in `internal/commitment/authority_internal_test.go` shows that `boundSources` returns one sorted list and one identity for the inputs `FT9`, `FT1` and `FT1`, `FT9` (FD20).
- [ ] `TestCommitmentPlanSourcesAreCanonical` in `internal/commitment/authority_test.go` shows that `BuildPlan` lists `FT10` before `FT9` when outcome `A` owns `FT9` and the later outcome `B` owns `FT10` (FD21).
- [ ] Outcome `A` approves the deliverable `spec` at `specs/b/spec.md`, and outcome `B` approves `spec` at `specs/a/spec.md`. The same test shows that `BuildPlan` lists the `a` path first (FD22).
- [ ] `bench test --package ./internal/commitment` and `bench test --package ./internal/commitment/repository` pass.
