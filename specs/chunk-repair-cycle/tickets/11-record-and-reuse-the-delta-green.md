# Record the delta green in its own class and reuse it only at the chunk checkpoint

Blocked by: 10-run-the-delta-checkpoint.md
Writes: internal/gate/delta_checkpoint.go (new), internal/gate/delta_checkpoint_test.go (new), internal/gate/checkpoint.go, internal/gate/run_transaction.go, internal/gate/gate.go, internal/gate/verdict.go, internal/gate/verdict_registry.go, internal/gate/verdict_registry_guard_test.go, internal/gate/composed_green.go, internal/gate/lane.go, internal/gate/lane_record_test.go, internal/gate/durable_replacement_test.go, internal/gate/complete_checkpoint_test.go, internal/git/git.go, internal/testreport/testreport_test.go, internal/testreport/check_test.go, internal/status/status_gatecache_test.go, .bench/BENCH-reference.md, projects/benchkit.md, CHANGELOG.md
Covers: RC68, RC69, RC70, RC74, RC75, RC82, RC83

## What to build

Review chunk: RC-C6B. The RC-C6 review closes before this ticket starts, because this ticket extends the run decision and reuses the reviewed-tip derivation of ticket 10.

A green delta writes a delta record, and only the chunk checkpoint reads it back. The prepared run at a repair commit shares the run decision of ticket 10, so a test-only repair's prepared run writes a delta green that the checkpoint after the record commit reuses.

The run decision gains its second step: reuse a delta green of the same subject and base, when the delta is still eligible. The order becomes: reuse a whole green of the current subject key; reuse a delta green; run the delta set; run the whole gate. `--fresh` still skips the first three steps. Step 2 reuses a delta green only when its `base_tree` and `tree` equal the current operands, and it checks the base-green condition again. Put the reuse step in `delta_checkpoint.go`, and add only calls to `run_transaction.go`, which is near its 400-line limit.

The delta record is a new verdict record class shaped like the lane record (`laneRecordFile` and `writeLaneRecord`, `internal/gate/lane.go:52-55,329-345`). Its fields are `schema`, `delta`, `base_tree`, `tree`, `recorded_at`, and `run_binary`, with no `oracle` and no `status`. Add the class row with a non-empty reuse refusal to the class table (`verdictRecordClasses`, `internal/gate/verdict_registry.go:28-45`) and to the pinned expectation of `TestVerdictRecordClassRegistryMatchesExpectation` (`internal/gate/verdict_registry_guard_test.go:15-46`). The record struct (`internal/gate/verdict.go:55-71`, at 397 of a 457-line grant) gains the fields. Each general reuse reader refuses the class (`internal/gate/composed_green.go:33`, `verdict.go:120-129`). Only step 2 reads the record, so a complete checkpoint after a delta-only chunk runs the whole gate (RC68).

Write the record to its own file `bench-last-delta` in the worktree git dir, never to the evidence store or to `bench-last-gate`. Name the file beside `GateCacheFile` (`internal/git/git.go:172-175`), so its name has one source. Add the delta writer to both `gateWriter` loops of `internal/gate/durable_replacement_test.go:33-63,80`. Add `bench-last-delta` to the gate-owned git-dir record lists that named checks must leave untouched: `internal/testreport/testreport_test.go:267-274` and `internal/testreport/check_test.go:267-269`. `check_test.go` is at 481 lines over its limit, so keep it line-neutral. `internal/status/status_gatecache_test.go:48-63` reads a legacy reduced class as an invalid cache; confirm that it still passes.

The prepared run (ticket 08) now takes the delta. At a repair commit, before the record commit refreezes the chunk, the record still names the earlier freeze tip. The derivation of ticket 10 takes the subject's HEAD as its target, so the prepared run finds the reviewed pair and writes a delta green for an eligible delta (RC74). A repair whose change from the reviewed tip holds a non-test file runs the whole gate and writes no delta green (RC82).

Update the documents (RC75). `.bench/BENCH-reference.md` names the delta class. `projects/benchkit.md:469-472` states that reduced-scope records are never authored; change it to name the delta class as the one reduced record. `projects/benchkit.md:118` names `bench-last-gate` as the gate cache; add the delta file. Do not add a second table with the `authoritative argv` header to the profile (`internal/conformance/profile_lane_table_test.go:104,165` parses it), and avoid the profile's Forbid needles on gate-output sentences.

## Acceptance

- [ ] RC68: a complete checkpoint after a delta-only chunk runs the whole gate.
- [ ] RC69: a green delta leaves `bench-last-gate` unchanged, and a later whole chunk checkpoint does not reuse the delta.
- [ ] RC70: the class table holds the delta class with a non-empty reuse refusal and no oracle field (`TestVerdictRecordClassRegistryMatchesExpectation`).
- [ ] RC74: a prepared run at a repair commit, before the record commit refreezes the chunk, writes a delta green for an eligible delta, and the checkpoint after the record commit reuses it.
- [ ] RC75: the reference and the profile name the delta schedule and the delta class as the one reduced schedule and record.
- [ ] RC82: a prepared run at a repair commit whose change from the reviewed tip holds a non-test file runs the whole gate and writes no delta green.
- [ ] RC83: after a new review round moves the reviewed tip, the chunk checkpoint does not reuse a delta green whose `base_tree` is the old base.
- [ ] The probe "write the delta green to the evidence store" bites in `bench test --package ./internal/gate`.
