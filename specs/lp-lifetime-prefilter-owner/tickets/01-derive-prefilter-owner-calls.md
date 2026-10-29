# The lifetime closure prefilter derives its owner calls from the owner-call list

Blocked by: none
Writes: internal/conformance/process_lifetime_closure_test.go
Covers: none

## What to build

The process-lifetime closure reads a production file only when the file spells a lifetime mechanic or an owner call.
The prefilter in `internal/conformance/process_lifetime_closure_test.go` types the owner calls `NewLifetime` and `RunInherited` a second time.
The owner-call list `processLifetimeOwnerCalls` in `internal/conformance/process_lifetime_test.go` is the source of these names.
Thus, an owner call that a change adds to that list can escape the prefilter.

Make the prefilter the concatenation of the process mechanics `Setpgid`, `Kill`, `Signal`, and `Wait` and the list `processLifetimeOwnerCalls`.
Do not change the closure diagnostics.

## Acceptance

- [ ] The prefilter in `internal/conformance/process_lifetime_closure_test.go` contains no owner-call name of its own and reads `processLifetimeOwnerCalls`.
- [ ] When a probe adds a new owner call to `processLifetimeOwnerCalls` and to the fixture, `TestProcessLifetimeOwnerEntries` stays green.
- [ ] When a probe removes `processLifetimeOwnerCalls` from the prefilter, `TestProcessLifetimeOwnerEntries` turns red.
- [ ] The `bench commit` lane passes.
