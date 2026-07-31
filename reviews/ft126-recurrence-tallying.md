## Standards

Count: 1

Worst: The migration-count expectation duplicates the nine-row count fact without
demonstrating the count mutation that makes the duplication independently necessary.

- `internal/conformance/docs_workflow_checks_test.go:41` repeats the nine durable
  counts encoded by `ROADMAP.md`; `AGENTS.md:30-35` permits that independently
  authored expectation only when a named omission or mutation makes the gate red and
  that red is recorded and demonstrated. `TestOccurrenceLedgerMigrationCheckBites`
  (`internal/conformance/docs_workflow_checks_test.go:60-68`) mutates only the legacy
  heading text. No count-mutation test exists (`rg` finds the count diagnostic and
  map only at lines 41 and 51), so changing (for example) FT158 from three keys to
  two has no demonstrated red. Add a focused count mutation that asserts the
  FT158-specific migration diagnostic.

## Spec

Count: 0

Worst: none. All 33 acceptance-coverage rows were audited against the pinned
`8e80941..8a40ed3` diff; the required behavior is present with no scope finding.

## Coverage

Count: 1

Worst: A degraded capture source can leave `sequence_trusted=true` and authorize
maintenance from partial evidence.

- With a valid ledger and a FIFO/wrong-type retrospective (or unreadable/wrong-type
  `IDEAS.md` or `.bench/learnings.md`), `BuildContext` records the source failure
  (`internal/roadmap/context_parse.go:225-251`) but computes trust only from
  occurrence discrepancies (`internal/roadmap/context_parse.go:302-303`). The
  maintenance instruction stops only when that field is false
  (`.agents/commands/bench-what-next.md:23-27`), so it can drain/rewrite while an
  unobserved capture unit remains. This violates the profile's fail-closed
  special-file rule (`projects/benchkit.md:124-126`). Existing AXI coverage proves a
  retrospective FIFO renders and exits zero
  (`internal/contract/axi/axi_roadmap_context_test.go:109-127`), but does not assert
  false trust or a maintenance refusal; structural occurrence failures do assert
  false trust (`internal/roadmap/context_test.go:139-198`). Make degraded capture
  sources fail the recurrence-trust authorization and add the focused assertion.
