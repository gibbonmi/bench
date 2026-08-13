# Reconcile canary inventory consumers

Blocked by: report-truthful-canary-inventory.md, remove-production-canary-dispatch.md
Writes: `internal/preprelease`, `.bench/lib/canary-run.sh`

## What to build

Reconcile the ship unit seam and both branches of the shipped compatibility shim with the landed inventory decision and exact success vocabulary, and verify the selected-executable journey already updated with CI1 reaches that same route. Keep this consumer family together because CI5 quantifies over the shared production route: omitting any one branch leaves the cross-consumer row red even though the command itself is correct.

## Acceptance

- [ ] (covers CI5) The ship step, both compatibility-shim branches, unit tests, and selected-executable journey use the same inventory decision and inventory-only diagnostics without starting an owner, gate, wrapper, `go test`, or `go run`.
