# Allow a verified clean provisional release

Blocked by: allow-already-covered-clean-integration.md
Ownership fence: `internal/worktree/path.go`, `internal/worktree/orphan_test.go`
Integration surfaces: exact provisional evidence -> cleanup plan; no-op checkpoint -> clean assignment checkout; terminal release -> assignment compaction
Contracts: `internal/worktree/path.go` releases a clean checkout at its assignment base only after the retained checkpoint, integrated candidate, index, and live tree all validate exactly; `internal/worktree/orphan_test.go` proves the no-op lifecycle shape releases without recovery while preserving every evidence-drift refusal

## What to build

Permit release of a verified provisional assignment whose no-op checkpoint leaves
the checkout clean at its base. The exact retained checkpoint and integrated
candidate evidence must remain mandatory, and the live tree and index must still
match the base/checkpoint tree. Release must use the terminal provisional-release
action, avoid a redundant recovery ref, compact the assignment, and remain
idempotent on replay. Do not admit an arbitrary clean checkout or weaken any dirty
payload, ignored-output, unsafe-path, or evidence-drift refusal.

## Acceptance

- [ ] [VR1] A no-op checkpoint whose checkout remains clean at the assignment base releases, creates no recovery ref, compacts the assignment, and succeeds on replay.
- [ ] [VR2] Clean release still refuses incomplete, retargeted, deleted, mismatched-tree, unsafe-path, and unretained checkpoint or integration evidence; dirty payload release keeps its existing exactness rules.

## Red mutations

- [ ] [MVR1] Requiring every non-legacy provisional release checkout to be dirty makes the honest no-op release case red.
- [ ] [MVR2] Accepting a clean checkout before exact evidence validation makes an existing evidence-drift case green.
