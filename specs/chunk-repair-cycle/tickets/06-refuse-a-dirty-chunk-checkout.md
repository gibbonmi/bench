# Refuse a dirty checkout at the chunk checkpoint

Blocked by: 05-one-plan-commit-and-preflight-per-dispatch.md
Writes: internal/gate/gate.go, internal/gate/checkpoint.go, internal/gate/complete_checkpoint.go, internal/refusalroute/faces_gate.go, internal/gate/checkpoint_prepare_test.go (new), internal/gate/complete_checkpoint_test.go, internal/gate/review_checkpoint_test.go, internal/gate/refusal_route_test.go, internal/gate/review_checkpoint_confirmation_test.go, internal/gate/record_repair_test.go, internal/gate/chunk_gap_test.go, internal/gate/checkpoint_persisted_test.go, internal/gate/review_checkpoint_commits_test.go, internal/gate/delegated_checkpoint_test.go, internal/gate/commitment_completion_test.go, internal/reviewrecord/recordtest/fixture.go, cmd/bench/gate_route_test.go, internal/worktree/delegated_integration_test.go
Covers: RC44

## What to build

Review chunk: RC-C4. This ticket is the prefactor of the chunk-checkpoint key (ticket 07).

Before the chunk checkpoint runs, it refuses a working tree that differs from HEAD, with the existing `checkpoint-dirty-checkout` face. The complete form already does this (`internal/gate/complete_checkpoint.go:31-37`), so the chunk form composes the same comparison and the same face. The review record has no exemption. If the chunk form shares the refusal text with the complete form, keep one source for it and move the pins at `internal/gate/complete_checkpoint_test.go:207,242-243`. The comment on the face (`internal/refusalroute/faces_gate.go:55-62`) names the landing only; extend it to the chunk checkpoint. Keep the candidate check first (`internal/gate/gate.go:201-203`).

This is a posture change, and it reds each test that saves the record and then runs a chunk checkpoint with no commit. Move each such call site to a committed fixture in this ticket:

- `internal/gate/review_checkpoint_test.go:100,137,170,189,192,307` (`runCheckpoint` at line 67 is the shared helper).
- `internal/gate/refusal_route_test.go:298,320` (`evidenceFixture`, `repairFixture`). The file is at 379 of 400 lines.
- `internal/gate/review_checkpoint_confirmation_test.go:274`, `internal/gate/record_repair_test.go:33-56`, `internal/gate/chunk_gap_test.go:124`, `internal/gate/checkpoint_persisted_test.go:148`, and the chunk runs of `internal/gate/review_checkpoint_commits_test.go` and `internal/gate/delegated_checkpoint_test.go`.
- `cmd/bench/gate_route_test.go:105-107` and `internal/worktree/delegated_integration_test.go:150`.

Each refusal case keeps its own reason. A commit in the shared helper, or a commit step in the fixture (`internal/reviewrecord/recordtest/fixture.go`, at 397 of 400 lines), is one way to do it. `TestChunkCheckpointGradesTheCheckoutTree` stays as it is, because its fixture commits.

Put the RC44 test in the new file `internal/gate/checkpoint_prepare_test.go`. Tickets 07 and 08 add their rows to the same file, and a new file over 400 lines reds the growth check, so keep its fixture helpers in an existing test file.

## Acceptance

- [ ] RC44: the chunk checkpoint refuses a working tree that differs from HEAD with the `checkpoint-dirty-checkout` face, before the evidence pass and before the oracle.
- [ ] A dirty review record alone also refuses.
- [ ] Each existing chunk-checkpoint refusal case keeps its own reason over a committed fixture.
- [ ] The complete checkpoint keeps its clean-checkout refusal.
