# FT78 implementation review

## Standards

3 findings. Worst: P2 duplicate gate-execution ownership.

- **P2 — duplicate gate-execution ownership.** `internal/gate/gate.go:151` and
  `internal/gate/gate.go:176` retain the exported `Run` and `RunContext` execution
  paths while `internal/gate/gate.go:308` adds `runCaptured` with a different
  environment and output policy. Repository search finds no callers of `Run` or
  `RunContext`; live callers route through `Execute`. This violates AGENTS.md's
  **one source per fact** rule and the Duplicated Code baseline. Consolidate the
  runner ownership without changing a public interface unless the reviewer reopens
  that out-of-scope spec decision.

- **P2 — proof IDs are authored twice.** The proof registries and their independently
  literal expected-ID lists repeat the same inventories, for example
  `internal/contract/runtime/runtime_gate_action_registry_test.go:30` and `:65`;
  the pattern recurs in `runtime_gate_projection_registry_test.go`,
  `runtime_gate_proof_test.go`, `internal/gate/fault_engine_test.go`, and
  `internal/gate/story4_proof_test.go`. This violates AGENTS.md's **one source per
  fact** rule and is Duplicated Code/Shotgun Surgery. The approved spec explicitly
  requires the independent literal lists, so resolving this finding requires a
  reviewer decision on that spec/standard conflict rather than deriving away the
  completeness check silently.

- **P3 — speculative projection field.** `internal/roadmap/context_types.go:50`
  adds `GateCacheFact.Reason`, and `internal/status/status.go:198` plus
  `cmd/bench/main.go:86` thread it through, but no roadmap, status, dashboard, or AXI
  renderer reads it. This is Speculative Generality. Remove it until a surface owns
  it, or add the decided consumer if its omission is unintended.

## Spec

1 finding. Worst: High — R17 retains old green evidence after lock faults.

- **High — R17 allows old green evidence to become reusable after lock faults.**
  `specs/oracle-bound-gate-verdicts.md:132` requires that for every R21 operation
  fault, **"old green never returns."** `internal/gate/gate.go:249` returns on lock
  open or acquisition failure without invalidating the prior durable verdict, and
  `internal/gate/runner_test.go:332` explicitly accepts reusable ready-green durable
  state for those two faults. Independent focused runs reported
  `durable=ready-green` for both `lock-open` and `lock-acquisition`. Persist a
  non-reusable state before returning and make the focused bridge reject old green.

## Coverage

1 finding. Worst: High — stale green reuse after a declared symlink target changes.

- **High — declared symlink target content can change without invalidating green.**
  For an ignored declared chain `inputs/link-a -> link-b -> target`,
  `internal/gate/subject.go:224` hashes each symlink's mode and link text but not the
  resolved target bytes. `internal/gate/story3_proof_helpers_test.go:347` proves only
  initial acceptance. In an independent real-wrapper repro, the target changed from
  green to red after a green verdict; `bench commit` still committed and the gate run
  count remained `1`. Add a real-wrapper case that mutates the in-repo target and
  asserts changed oracle identity, a second gate run, and commit refusal.
