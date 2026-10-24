# Run a delta checkpoint for a test-only delta after a whole green

Blocked by: 09-list-the-system-suite-in-the-selector.md
Writes: internal/gate/checkpoint.go, internal/gate/delta_checkpoint.go (new), internal/gate/delta_checkpoint_test.go (new), internal/gate/run_transaction.go, internal/gate/gate.go, internal/gate/chunk_gap.go, internal/gate/engine.go, internal/gate/gate_go.go, internal/gate/phases.go, internal/gate/phases_test.go, internal/gate/tag_census_test.go, internal/gate/report.go, internal/gate/complete_checkpoint_test.go, internal/gate/review_checkpoint_confirmation_test.go, internal/gate/authorization/authorization.go, internal/reviewrecord/confirmation.go, internal/reviewrecord/confirmation_test.go, internal/reviewrecord/write.go, internal/testgap/testgap.go, internal/testgap/testgap_test.go, .bench/BENCH-reference.md, projects/benchkit.md
Covers: RC58, RC59, RC60, RC61, RC62, RC63, RC64, RC65, RC66, RC67, RC71, RC72, RC73

## What to build

Review chunk: RC-C6. RC-C4 and RC-C4B (the subject key and the prepared run) and RC-C5 (the system suite in the selector) land first.

A test-only delta after a whole green at the reviewed tip takes a delta checkpoint: gofmt, vet, `go test` of the changed packages, and root conformance. Any other delta runs the whole gate. A delta green authorizes no landing, no complete checkpoint, and no later whole checkpoint.

The chunk checkpoint decides its run in this order: reuse a whole green of the current subject key (RC-C4); run the delta set when the delta is eligible; run the whole gate. `--fresh` skips the reuse and the delta (`internal/gate/checkpoint.go:65,83-84`). Ticket 11 adds the delta-green reuse step and makes the prepared run share this decision; until then the prepared run keeps its whole run, because it needs a green that a later call can reuse. Plug the decision in before the run seam (`internal/gate/run_transaction.go:43,57,75-95`), which is at 381 of 400 lines. Add only the call there, and put the decision in `delta_checkpoint.go`.

The delta has two operands. The base is the record-stripped tree of the reviewed tip, and the subject is the record-stripped tree that the run grades. Add one `reviewrecord` derivation of the reviewed tip that takes a target commit. It applies the pair checks of `FindReviewedPair` (`internal/reviewrecord/confirmation.go:74-112`): one pair for all three axes, a tip other than the chunk base, a tip that descends from the chunk base, a tip other than the target, and a tip that is an ancestor of the target. The derivation takes the subject's HEAD commit as its target and never reads the chunk tip or the chunk digest of the record. `FindReviewedPair` composes it with the chunk tip as its target, so the confirmation and the delta share one rule. Keep the `FindReviewedPair(root, spec, chunk)` signature, because `internal/reviewrecord/write.go:272` and `confirmation.go:136` call it, and keep the refusal texts that `internal/reviewrecord/confirmation_test.go` and `internal/gate/review_checkpoint_confirmation_test.go` pin.

A delta is eligible when every condition holds:

- The test-gap proof, with the RC-C1 exemption, proves the change from the base to the subject test-only.
- A whole green exists for the base under the prospective base oracle of the base subject, read through the prospective inspection (`internal/gate/engine.go:90-102`) with the non-expiring window (`engine.go:108-141`). Add an unexported read for that; keep the signatures that `internal/gate/authorization/authorization.go:152,187` call.
- No changed test file sits in `gate.SystemSuite` or in a race-table package (`raceTests`, `internal/gate/gate_go.go:32-40`; `internal/gate/tag_census_test.go:159,169-174` pins it).

The delta set runs gofmt and vet with the argv of their gate phases over the whole root (`internal/gate/phases.go:107-110`), `go test` of the directories of the files in the RC-C1 test-gap list (`internal/testgap/testgap.go`), and the root-conformance entry (`phases.go:147-164`). A `_test.go` file always builds with the package of its directory. The gate cannot import `testreport`, so no checkpoint path calls the changed-package selector (RC60). `phases.go` is at 394 of 400 lines; put the delta phase builder in the new `internal/gate/delta_checkpoint.go` and share the argv with the whole phases through one source.

A red delta makes the checkpoint red, with no fallback to the whole gate. A green delta opens the chunk gap as a whole chunk checkpoint does (`internal/gate/gate.go:213-215`). A delta green writes nothing to the evidence store or to `bench-last-gate`, so the landing (`engine.go:70-88`), the complete checkpoint, and a later whole checkpoint find no delta green. The complete form never enters this decision.

The gate output prints the delta phases table and the line `gate: green (delta checkpoint)`. Keep the `aggregateAndReport` signature (`internal/gate/report.go:55`), because it has 35 call sites, and put the delta output tests in `internal/gate/delta_checkpoint_test.go`, not in `report_test.go` (942 lines, over its limit). Whole runs keep their output, so the exact-line pins of `gate: green` stay. `.bench/BENCH-reference.md` Gate output owns the delta output account. `projects/benchkit.md:387-389` states that no stripped-subject phase schedule exists; change it to name the delta schedule as the one reduced schedule. The profile forbids gate-output sentences (`internal/anchors/registry_data.go:159`, `internal/anchors/registry_retained_workflow.go:80,105-110`).

A delta row plants a red in a package outside the delta set. The delta stays green there, and the whole run goes red. Keep `delta_checkpoint_test.go` under 400 lines, because ticket 11 adds rows to it.

## Acceptance

- [ ] RC58: an eligible test-only delta runs gofmt and vet over the whole root, `go test` of the changed packages, and root conformance, and no other phase.
- [ ] RC59: the delta package set equals the directories of the test-gap file list, a path with a space included.
- [ ] RC60: no checkpoint path calls the changed-package selector.
- [ ] RC61: a test-only delta with no whole green at the reviewed tip runs the whole gate.
- [ ] RC62: a whole base green older than the freshness window still makes the delta eligible.
- [ ] RC63: a delta with a non-test file runs the whole gate.
- [ ] RC64: a delta whose only non-test change is an assignment-only plan commit stays eligible.
- [ ] RC65: a delta that touches the system-suite package runs the whole gate.
- [ ] RC66: a delta that touches a race-table package runs the whole gate.
- [ ] RC67: after a green delta, the landing inspection is not green and the evidence store holds no new entry.
- [ ] RC71: `--fresh` on an eligible delta runs the whole gate.
- [ ] RC72: a green delta checkpoint opens the chunk gap.
- [ ] RC73: a red delta run makes the checkpoint red and records no green.
- [ ] The probes "admit a system-suite file to the delta" and "take the delta base from the record chunk tip" each bite in `bench test --package ./internal/gate`.
