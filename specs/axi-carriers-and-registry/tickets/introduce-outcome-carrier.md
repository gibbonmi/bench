# Introduce the outcome carrier

Blocked by: none
Ownership fence: `internal/axi/outcome.go`, `internal/axi/outcome_test.go`
Integration surfaces: outcome policy API→`internal/axi/outcome.go`; registry declarations→declare-axi-query-root-metadata.md; compatibility oracle package root→`internal/axi/compatibility` exercised unchanged (the sibling spec owns it; this ticket adds only root-package files beside it)
Contracts: domain kind string, exit-policy map, typed payload, aggregate list, and action list cross caller→`internal/axi/outcome.go`, membership is owner-declared, ordering is payload-aggregates-actions, and absence is an explicit empty slice, asserted by OC1 against a real domain policy value rather than a stub validator
Closure: OC1/kind, OC1/exit-policy, OC1/exit-preimage, OC1/payload, OC1/order, OC1/absence

## What to build

`internal/axi` gains a generic outcome carrier: the domain supplies its own kind
set, its own kind-to-exit map, a typed payload, an ordered aggregate list, and an
ordered action list, and the carrier validates and returns exactly those facts.
The carrier derives no meaning of its own — it never invents a kind, never falls
back to a universal 0/1/2 exit table, never re-encodes a payload, and never turns
an owner-declared empty list into a default row.

Tree condition at refresh time: this spec follows `axi-compatibility-oracle`, so
`internal/axi/compatibility` already exists as a sibling package directory. This
ticket writes only the two root-package files on its fence line and must not
touch anything under `internal/axi/compatibility`. If `internal/axi` does not yet
exist as a directory, this ticket creates it together with `outcome.go`.

## Acceptance

- [ ] [OC1] (covers CR1) shared outcomes preserve owner-declared kinds, typed payloads, component order, explicit absence, and the domain's exact exit policy.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OC1/kind | delete the domain kind-set membership lookup in `Outcome.Validate` so any kind string is accepted | `TestOutcomeRejectsKindOutsideDomainPolicy` in `internal/axi` (this ticket authors it) | run `go test ./internal/axi -run TestOutcomeRejectsKindOutsideDomainPolicy -timeout 60s`; expect the `errors.Is(err, axi.ErrUnknownKind)` assertion to report a nil error for kind `"invented"`; bound is the `-timeout 60s` binary deadline, and the carrier does no I/O and starts no goroutine so nothing else can stall |
| OC1/exit-policy | replace the domain-supplied exit map with a universal `ok=0 / error=1 / usage=2` table | `TestOutcomeExitCodeComesFromDomainPolicy` in `internal/axi` | run `go test ./internal/axi -run TestOutcomeExitCodeComesFromDomainPolicy -timeout 60s`; expect the equality assertion `ExitCode() = 1, want 3` for the test policy's `stale` kind, which declares exit 3; bound is the `-timeout 60s` binary deadline over an in-memory table |
| OC1/exit-preimage | drop the kind field from the key `exitFor` looks up, so every kind resolves to the policy's first entry (input mutation of the value that authorizes the process exit) | `TestOutcomeExitCodeComesFromDomainPolicy` in `internal/axi` | run `go test ./internal/axi -run TestOutcomeExitCodeComesFromDomainPolicy -timeout 60s`; expect the equality assertion `ExitCode() = 0, want 3` for kind `stale` while kind `ok` still reports 0, so a control-flow-only fix cannot make it green; bound is the `-timeout 60s` binary deadline |
| OC1/payload | return the payload re-encoded as its `fmt.Sprint` string instead of the owner's typed value | `TestOutcomePreservesTypedPayload` in `internal/axi` | run `go test ./internal/axi -run TestOutcomePreservesTypedPayload -timeout 60s`; expect the type assertion `Payload().(*policyPayload)` to fail with `payload has dynamic type string`; bound is the `-timeout 60s` binary deadline |
| OC1/order | emit the action list before the aggregate list in `Outcome.Components` | `TestOutcomePreservesPayloadAggregateActionOrder` in `internal/axi` | run `go test ./internal/axi -run TestOutcomePreservesPayloadAggregateActionOrder -timeout 60s`; expect the index assertion `Components()[1].Kind = "action", want "aggregate"`; bound is the `-timeout 60s` binary deadline |
| OC1/absence | substitute a one-row default aggregate list when the owner supplies an explicitly empty one | `TestOutcomeKeepsDeclaredAbsenceDistinct` in `internal/axi` | run `go test ./internal/axi -run TestOutcomeKeepsDeclaredAbsenceDistinct -timeout 60s`; expect the assertion `len(Aggregates()) = 1, want 0` on the outcome constructed with an explicit empty slice; bound is the `-timeout 60s` binary deadline |
