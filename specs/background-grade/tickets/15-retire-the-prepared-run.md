# Retire the prepared run

Blocked by: 10-let-a-chunk-checkpoint-wait-on-a-background-grade.md
Writes: internal/gate/checkpoint.go, internal/gate/gate.go, internal/gate/checkpoint_prepare_test.go, internal/gate/chunk_gap_test.go, internal/gate/review_checkpoint_test.go, internal/gate/lane_record_test.go, internal/gate/delta_checkpoint_test.go, internal/gate/complete_checkpoint_test.go, internal/gate/background_grade_test.go (new), cmd/bench/gate_route_test.go, cmd/bench/help_inventory_test.go, .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .bench/BENCH-reference.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, CHANGELOG.md, internal/gate/baseline_commit_test.go
Covers: BG4, BG86, BG87, BG88, BG89, BG90, BG91, BG92, BG100, BG101, BG102, BG103

## What to build

The background grade now grades each chunk subject ahead of its checkpoint, so the prepared run retires.

- Remove `--prepare` from the gate grammar, the usage line, the refusal reruns, and `gradesSubjectAlone` (`internal/gate/checkpoint.go:26-28`, `internal/gate/checkpoint.go:84-88`, `internal/gate/checkpoint.go:96`, `internal/gate/checkpoint.go:125-126`, `internal/gate/checkpoint.go:322-323`, `internal/gate/gate.go:224`, `internal/gate/gate.go:231`).
- A call with `--prepare` exits 2 with the usage line.
- Move each prepared-run test behavior that still matters onto the background grade, and delete the rest. `TestCompleteCheckpointIgnoresAPreparedChunkGreen` becomes a test over a background chunk green.
- The Land section states that a chunk checkpoint waits for the live background grade of its key.

Remove each sentence that grants the prepared run, and add one Forbid anchor for each:

- The Land section of `/bench-implement-spec`: the start at a chunk freeze, the start after a repair, and the wait on the PID of a prepared run (`.agents/commands/bench-implement-spec.md:64`, `.agents/commands/bench-implement-spec.md:66`).
- The bounded repair policy: the prepared chunk run sentence and the ordering sentence (`.agents/skills/bench-craft-line/references/bounded-repair-policy.md:66`, `.agents/skills/bench-craft-line/references/bounded-repair-policy.md:68`).
- `.bench/BENCH-reference.md`: every sentence that holds `--prepare` or the phrase `prepared run`, and the sentence that lists what the prepared run skips (`.bench/BENCH-reference.md:341`, `.bench/BENCH-reference.md:387`, `.bench/BENCH-reference.md:650-658`, `.bench/BENCH-reference.md:721`, `.bench/BENCH-reference.md:745`).

The staged `repair-carry` spec also edits the Land section. Edit the text that is current when this ticket starts.

## Acceptance

- [ ] `bench gate --checkpoint <spec> --chunk 1 --prepare` exits 2 with the usage line.
- [ ] The gate help row names no `--prepare`.
- [ ] A complete checkpoint after a background chunk green of the same source runs its own oracle.
- [ ] The Land section holds no `--prepare` token, no sentence that starts a prepared run after a repair, and no wait on the PID of a prepared run.
- [ ] The bounded repair policy holds neither of its two prepared-run sentences.
- [ ] `.bench/BENCH-reference.md` holds no `--prepare` token, no phrase `prepared run`, and no sentence that lists what a prepared run skips.
- [ ] Each retired sentence, put back alone, makes `bench test --check docs-currency-workflow` red.
- [ ] The Land section states that a chunk checkpoint waits for the background grade of its key.
