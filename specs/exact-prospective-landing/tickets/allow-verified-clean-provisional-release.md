# Allow a verified clean provisional release

Blocked by: allow-already-covered-clean-integration.md
Ownership fence: `internal/worktree/path.go`, `internal/worktree/orphan_test.go`
Integration surfaces: exact provisional evidence, cleanup plan, and assignment compaction→`internal/worktree/path.go`; clean no-op checkpoint producer→allow-already-covered-clean-integration.md; successful release and replay proof→`internal/worktree/orphan_test.go` plus VR1; evidence-drift refusals→`internal/worktree/orphan_test.go` plus VR2
Contracts: retained checkpoint and integrated-candidate identities cross allow-already-covered-clean-integration.md→`internal/worktree/path.go` release validation, asserted by VR1-VR2 against the real provisional release lifecycle in `internal/worktree/orphan_test.go`

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

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| VR1 | require every non-legacy provisional release checkout to be dirty | `TestReleaseProvisionalRemovesVerifiedCleanNoOpCheckpoint` | run `go test ./internal/worktree -run '^TestReleaseProvisionalRemovesVerifiedCleanNoOpCheckpoint$' -count=1`; expect the verified clean release and replay to fail |
| VR2 | accept a clean checkout before exact evidence validation | `TestReleaseProvisionalRefusesVerifiedCleanNoOpEvidenceDrift` and `TestReleaseProvisionalRefusesLiveCheckpointDrift` | run `go test ./internal/worktree -run '^(TestReleaseProvisionalRefusesVerifiedCleanNoOpEvidenceDrift|TestReleaseProvisionalRefusesLiveCheckpointDrift)$' -count=1`; expect an evidence-drift case to stop refusing |
