# FT78 implementation review

## Blocked

1 Standards finding remains. Worst: P2 proof inventories conflict with the approved
specification.

- **P2 — proof IDs are authored twice by explicit requirement.** The proof registries
  and their independently literal expected-ID lists repeat the same inventories in
  `internal/contract/runtime/runtime_gate_action_registry_test.go`,
  `runtime_gate_projection_registry_test.go`, `runtime_gate_proof_test.go`,
  `internal/gate/fault_engine_test.go`, and `internal/gate/story4_proof_test.go`.
  AGENTS.md's **one source per fact** standard requires one derived inventory, while
  the approved spec's Proof discipline requires an **independently literal** expected-ID
  set so removing a registration still emits
  `FT78 proof ledger completeness contract failed`. Deriving the expected set from the
  registry removes that completeness bite; retaining the literal set preserves the
  oracle but duplicates the fact. Resolution requires the reviewer to reopen either
  the repository standard or the approved proof requirement. The implementation must
  not choose between them silently.
