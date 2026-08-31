# Derive the executed tag census from the phase table

Blocked by: none
Writes: internal/gate/tag_census.go (new), internal/gate/tag_census_test.go (new)
Covers: CE3

## What to build

The gate package exports one derivation of the build tags the gate executes on a
host. The derivation reads the resolved phase table: the root's manifest when the
root declares one, else the built-in kit table. The runner executes the same
resolution.

The derivation parses
each test phase's argv. It collects the `-tags` set of each such argv. It also
collects the untagged default set, because the ordinary test phase carries no
`-tags` flag. The result is the census that a later ticket consumes.

The census has one source. Do not copy the tag list into a constant. Derive the list from the resolved phase table, so the census and the oracle
cannot disagree.

A root with no test phase yields an empty census. The caller reads that empty
result as "inapplicable", so a non-Go root keeps its current behavior.

## Acceptance

- [ ] CE3 — the census on the kit root holds the untagged default set and the
      `system` set.
- [ ] a root whose phase manifest declares a custom `-tags` set yields that
      set in the census.
- [ ] CE3 — the census is empty for a root that declares no test phase.
- [ ] `bench gate` stays green.

## Delegate charge

You work in the Bench repo on the `citation-execution-proof` spec. Read
`specs/citation-execution-proof/spec.md` first. Then read
`internal/gate/phases.go` in full.

Add `internal/gate/tag_census.go`. Declare one exported function there. It
returns the build-tag sets the gate executes for a given root and kit.

Derive
the sets from the same phase-table resolution the runner executes: the root's
manifest when one is declared, else the built-in kit table. The gate package
already holds that unexported resolver. Call it. Do not re-derive the choice.

Select the phases whose argv is a `go test` invocation. Parse the `-tags=`
operand of each such argv. Add the empty set for a test phase with no `-tags`
operand. Return an empty result when the table declares no test phase. Do not add
a second copy of the tag list anywhere.

Add `internal/gate/tag_census_test.go` with `TestExecutedTagCensus`. Assert the untagged and `system` sets for the kit root. Assert a
manifest-declared custom tag appears in that root's census. Also assert the
empty result for a root with no test phase. Write the Go doc comments in the register of
`internal/gate/phases.go`.

Run only `bench worktree exec <label> -- go test ./internal/gate/`. Do not
commit. Do not edit the spec. Report the exported name and its signature, because
a sibling ticket consumes it.
