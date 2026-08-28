# Count release evidence against the proven targets

Blocked by: add-proven-target-field.md
Writes: internal/releaseevidence, scripts/verify-release-artifact.mjs, internal/conformance/native_workflow_test.go

## What to build

Every proof count derives from the proven targets, not from the shipped targets. Three consumers count today, and all three move in this ticket.

`internal/releaseevidence` reads the plan's `native_proof` flag into its target evidence. Proof inspection skips an unproven target, and finalization reports an incomplete proof set only when a proven target has no proof.

`scripts/verify-release-artifact.mjs` compares the release index's `native_proofs` length against the proven-target count. This verifier runs inside the workflow through `smoke-artifacts.sh`, so a shipped-count comparison would redden every evidence and smoke run.

The release evidence probe in `internal/conformance/native_workflow_test.go` synthesizes a proof per proven target, so the probe grades the set a real run produces.

All three keep their per-operating-system clauses unchanged, so a target returning to the proven list needs plan data only.

The contract this ticket reads from `add-proven-target-field.md` is the plan field name, its required-boolean rule, and the Go plan struct field that ticket adds.

## Acceptance

- [ ] Finalization succeeds when proofs exist for the proven targets only (row B9).
- [ ] Finalization fails when one proven target has no proof (row B10).
- [ ] Finalization still fails when a present proof does not match the inspected artifacts.
- [ ] `verify-release-artifact.mjs` accepts an index holding proofs for the proven targets only (row B13).
- [ ] `verify-release-artifact.mjs` still fails an index whose proof is red or mismatched.
- [ ] The release evidence probe writes one proof per proven target (row B14).
