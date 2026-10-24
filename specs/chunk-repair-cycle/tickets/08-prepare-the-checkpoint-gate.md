# Prepare the checkpoint gate beside the review and the confirmation

Blocked by: 07-key-the-chunk-checkpoint-without-the-record.md
Writes: internal/gate/checkpoint.go, internal/gate/gate.go, internal/gate/chunk_gap.go, internal/gate/run_transaction.go, internal/gate/evaluation.go, internal/gate/checkpoint_prepare_test.go (new), internal/gate/chunk_gap_test.go, internal/gate/complete_checkpoint_test.go, internal/gate/review_checkpoint_test.go, internal/gate/checkpoint_persisted_test.go, internal/gate/refusal_route_test.go, internal/gate/command_test.go, internal/refusalroute/route_test.go, cmd/bench/help_inventory_test.go, cmd/bench/gate_route_test.go, cmd/bench/tree_scope_test.go, cmd/bench/main.go, .agents/commands/bench-implement-spec.md, .bench/BENCH-reference.md, projects/benchkit.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, docs/adr/0002-accepted-enforcement-postures.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go
Covers: RC34, RC35, RC36, RC37, RC38, RC42, RC43, RC45, RC47, RC49, RC50, RC51, RC52

## What to build

Review chunk: RC-C4.

The orchestrator starts the checkpoint gate when a chunk freezes or a repair commits, with `bench gate --checkpoint <spec-path> --chunk <id> --prepare`. The run grades the subject that the chunk checkpoint grades (ticket 07) and writes its green to the evidence store under that subject key. The chunk checkpoint after the record commit reuses that green.

The prepared run skips the evidence pass, the record prose grade, the tip-moved guard, and the chunk gap. The key holds only the stripped tree that the run captured at its start, so a moved HEAD cannot misattribute the green. A record write or a record commit during the run leaves its verdict valid, because the post-run capture reads the private checkout (`internal/gate/evaluation.go:53-62`) and the subject excludes the record. Skip the chunk-gap write at `internal/gate/gate.go:213-215` (`openChunkGap`, `internal/gate/chunk_gap.go:14`) for a prepared run.

Add the flag to the gate flag constants, the parser, and the validator (`internal/gate/checkpoint.go:47-104`). `--prepare` with `--complete`, or with no `--checkpoint`, is a usage error at exit 2. If the checkpoint value gains a field, keep the complete form's purpose JSON unchanged (`checkpoint.go:277` marshals the checkpoint value). The usage line is `CommandUsage` (`internal/gate/gate.go:220`), which `internal/gate/command_test.go:17-20`, `cmd/bench/gate_route_test.go:78-80,116`, `cmd/bench/tree_scope_test.go:120-121`, and `cmd/bench/main.go:149` read. `cmd/bench/help_inventory_test.go:128` holds a hand-typed copy of the gate row; add `--prepare` there. `cmd/bench/main.go` is at 429 lines over its limit, so keep it line-neutral. The RC50 and RC51 cases go beside the usage-error table at `cmd/bench/gate_route_test.go:78`. The refusal routes reuse the existing gate faces. `rerunArguments` (`checkpoint.go:221-236`) rebuilds the rerun of a refused call; keep `--prepare` in it so the rerun repeats the refused command, and move the rerun pins at `internal/gate/review_checkpoint_test.go:279`, `internal/gate/checkpoint_persisted_test.go:114,170`, `internal/gate/refusal_route_test.go:227-332`, and `internal/refusalroute/route_test.go:21` only where their text changes.

A chunk checkpoint while a prepared run holds the gate lock refuses with the existing message (`internal/gate/run_transaction.go:90-96`). `bench probe` refuses on the gate lock (`internal/probe/set.go:223`).

Update the guidance (RC43, RC47). In `.agents/commands/bench-implement-spec.md`, the orchestrator starts a prepared run when a chunk freezes, beside the review axes, and when a repair commits, beside the confirmation. It runs the coordinator probe before the prepared run, and it waits on the run's PID before the checkpoint. At a freeze, the prepared run starts only when each axis probes in its own worktree. The file is at 95 of 95 lines and each paragraph is one line, so keep the paragraph count. Keep the line 68 sentence "Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`." verbatim (`internal/anchors/registry_chunk_chain.go:23`), and avoid the Forbid needle "--chunk <id>` before the successor" (`registry_chunk_chain.go:24`).

State the chunk-checkpoint key and the prepared run in `.bench/BENCH-reference.md` (lines 339, 379-383 for the chunk gap, 547-559, and Gate output), `projects/benchkit.md:469-473`, the bounded repair policy (line 62, "It does not edit the gate subject while the gate runs", and line 47-48), and ADR 0002 posture 5. The profile forbids gate-output sentences (`bench anchors projects/benchkit.md`), so the output account goes to the reference only.

`internal/gate/checkpoint_prepare_test.go` holds the rows of tickets 06, 07, and 08. A reuse row runs the prepared gate and the checkpoint as two process calls, so the second call reloads the evidence store. Keep the file under 400 lines; put shared helpers in an existing gate test file.

## Acceptance

- [ ] RC34: `--prepare` runs the whole gate on the record-stripped HEAD tree and writes a green to the evidence store.
- [ ] RC35: `--prepare` leaves the chunk gap closed.
- [ ] RC36: a chunk checkpoint in a second process after a record-only commit reuses the prepared green and runs no phase.
- [ ] RC37: a spec change after the prepared run makes the chunk checkpoint run the gate.
- [ ] RC38: a source change other than the record after the prepared run makes the chunk checkpoint run the gate.
- [ ] RC42: a record write and a record commit in the worktree during the prepared run leave its green.
- [ ] RC43: the implement phase runs the coordinator probe before the prepared run and waits on the run's PID before the checkpoint.
- [ ] RC45: a complete checkpoint after a prepared chunk green runs the gate on its published tree.
- [ ] RC47: the reference, the profile, the bounded repair policy, and ADR 0002 state the chunk-checkpoint key and the prepared run.
- [ ] RC49: `bench gate --help` names `--prepare` in its usage line.
- [ ] RC50: `--prepare` with `--complete` is a usage error at exit 2.
- [ ] RC51: `--prepare` without `--checkpoint` is a usage error at exit 2.
- [ ] RC52: a prepared green past the freshness window makes the chunk checkpoint run the gate.
- [ ] The probe "keep the tip-moved guard in the prepared run" bites in `bench test --package ./internal/gate`.
