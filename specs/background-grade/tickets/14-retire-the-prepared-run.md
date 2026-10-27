# Retire the prepared run

Blocked by: 09-let-a-chunk-checkpoint-wait-on-a-background-grade.md
Writes: internal/gate/checkpoint.go, internal/gate/gate.go, internal/gate/checkpoint_prepare_test.go, internal/gate/chunk_gap_test.go, internal/gate/review_checkpoint_test.go, internal/gate/lane_record_test.go, internal/gate/delta_checkpoint_test.go, internal/gate/complete_checkpoint_test.go, internal/gate/background_grade_test.go (new), cmd/bench/gate_route_test.go, cmd/bench/help_inventory_test.go, .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .bench/BENCH-reference.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, CHANGELOG.md
Covers: BG4, BG86, BG87, BG88, BG89, BG90, BG91, BG92

## What to build

The background grade now grades each chunk subject ahead of its checkpoint, so the prepared run retires.

- Remove `--prepare` from the gate grammar, the usage line, the refusal reruns, and `gradesSubjectAlone` (`internal/gate/checkpoint.go:26-28`, `internal/gate/checkpoint.go:84-88`, `internal/gate/checkpoint.go:96`, `internal/gate/checkpoint.go:125-126`, `internal/gate/checkpoint.go:322-323`, `internal/gate/gate.go:224`, `internal/gate/gate.go:231`).
- A call with `--prepare` exits 2 with the usage line.
- Move each prepared-run test behavior that still matters onto the background grade, and delete the rest. `TestCompleteCheckpointIgnoresAPreparedChunkGreen` becomes a test over a background chunk green.
- Remove the prepared-run sentences from the Land section of `/bench-implement-spec`, from the bounded repair policy, and from `.bench/BENCH-reference.md`.
- The Land section states that a chunk checkpoint waits for the live background grade of its key.
- Add one Forbid anchor for each retired sentence form, so a sentence that grants the prepared run cannot return.

The staged `repair-carry` spec also edits the Land section. Edit the text that is current when this ticket starts.

## Acceptance

- [ ] `bench gate --checkpoint <spec> --chunk 1 --prepare` exits 2 with the usage line.
- [ ] The gate help row names no `--prepare`.
- [ ] A complete checkpoint after a background chunk green of the same source runs its own oracle.
- [ ] The Land section holds no sentence that starts a `--prepare` run or waits on the PID of a prepared run.
- [ ] The bounded repair policy holds no sentence about a prepared chunk run.
- [ ] `.bench/BENCH-reference.md` holds no `--prepare` form.
- [ ] `bench test --check docs-currency-workflow` reds when one retired sentence returns.
- [ ] The Land section states that a chunk checkpoint waits for the background grade of its key.
