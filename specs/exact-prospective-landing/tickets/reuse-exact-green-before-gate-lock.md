# Reuse exact green before the gate lock

Blocked by: none
Ownership fence: `internal/gate`
Integration surfaces: exact-subject evaluation plan→existing `internal/gate/gate.go` asserted by RL1-RL3; prospective-tree execution→existing `internal/gate/engine.go` plus RL1; commit adapter consumer→adopt-exact-landing-in-commit.md
Contracts: immutable tree and oracle identity cross prospective evaluation→retained verdict inspection in `internal/gate`, asserted by RL1-RL3 against the real verdict record; non-reusable evidence crosses retained inspection→locked gate execution in `internal/gate`, asserted by RL2 against the real lock

## What to build

Answer an exact fresh-green subject from its retained verdict before attempting the repository gate lock, while preserving the locked validation and execution path for stale, mismatched, unavailable, or red evidence.

## Acceptance

- [ ] [RL1] An exact reusable green subject returns the canonical reuse result without opening or acquiring the gate lock, including prospective-tree execution while another process holds that lock.
- [ ] [RL2] Stale, tree-mismatched, oracle-mismatched, unavailable, and red evidence never bypasses lock acquisition or gate execution through the pre-lock check.
- [ ] [RL3] A reusable green that becomes available only after the initial check can still be reused under the lock from the same single-sourced evidence predicate without rewriting its verdict timestamp.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RL1 | move the first reusable-evidence check below lock acquisition | gate reuse-under-contention test | seed exact green, hold `bench-gate.lock`, execute the same prospective tree, expect canonical reuse and byte-identical verdict record |
| RL2 | treat any ready green record as reusable without the current plan's tree and oracle identities | hostile retained-evidence table | seed each stale or mismatched record, hold or observe the lock path, expect no reuse result |
| RL3 | delete the post-acquire reuse check | injected timing test | make exact green appear between the pre-lock check and lock acquisition, acquire the lock, expect reuse without oracle execution or verdict rewrite |
