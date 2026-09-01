# Single-source the prospective checkout name

Blocked by: none
Writes: internal/gate/prospectiveartifact/prospectiveartifact.go, internal/gate/prospectiveartifact/prospectiveartifact_test.go, internal/systemtest/owner_artifact_recovery_test.go
Covers: LF21

## What to build

Expose one checkout child-name owner from prospectiveartifact. Make the system
journey and package tests consume it instead of retyping the value.
The system journey preserves the current BENCH_KIT pin for its fixture gate.

## Acceptance

- [ ] Production, package tests, and the system journey use one name owner.
- [ ] The system journey no longer carries three copied literals.
- [ ] Existing checkout layout remains compatible.

## Red evidence

Mutation kind and site: replace `CheckoutName = "checkout"` with
`CheckoutName = "mutated-checkout"` in
`internal/gate/prospectiveartifact/prospectiveartifact.go`.

Command: `GOCACHE=/tmp/bench-checkout-mutation-gocache /home/mgibs/.local/opt/go-v1.25.0/bin/go test ./internal/gate/prospectiveartifact -run '^TestCheckoutNameKeepsTheCheckoutLayout$' -count=1`

The mutation failed with `checkout name = "mutated-checkout", want checkout`.
After source restoration, the same command passed.
