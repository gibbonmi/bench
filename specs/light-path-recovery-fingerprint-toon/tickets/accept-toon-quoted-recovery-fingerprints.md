# Accept TOON-quoted recovery fingerprints

Blocked by: none
Ownership fence: `internal/worktree/clean.go`, `internal/worktree/recovery_plan_test.go`, `internal/contract/runtime/runtime_worktree_test.go`, `CHANGELOG.md`
Contracts: the fingerprint crosses the `recovery_cleanup` TOON receipt -> shell field extraction -> `internal/worktree/clean.go` argument parser, asserted by RF1 against a ref produced by real create/release commands
Assumptions: the fingerprint value remains exactly 64 lowercase hexadecimal characters; bare fingerprints remain accepted; planning and fingerprint derivation stay unchanged

## What to build

`bench worktree recovery` can quote a fingerprint when the TOON encoder classifies the
hex string as numeric-looking. A shell cleanup loop that extracts the seventh receipt
field passes that representation intact, and the command rejects it before reading the
intent ledger. Normalize exactly one TOON quote layer before applying the existing
fingerprint validation, without widening the fingerprint value domain or requiring a
second plan.

## Acceptance

- [ ] [RF1] A ref produced by real `bench worktree create` and `bench worktree release` can be discarded with the plan's valid fingerprint carrying exactly one surrounding TOON quote layer, without replanning.
- [ ] [RF2] The same valid fingerprint remains accepted bare, while short, uppercase, non-hex, mismatched-quote, and extra-quote forms still refuse as invalid invocation before changing the ref or ledger.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RF1 | pass the quoted token directly to hex decoding | the real create/release runtime recovery contract | remove quote normalization, run the focused runtime recovery contract, expect exit 2 and the fresh ref to survive |
| RF2 | skip the existing lowercase 64-hex validation after normalization | the recovery argument parser contract | accept a short or uppercase quoted token, run the focused worktree parser test, expect the malformed-input refusal to fail |
