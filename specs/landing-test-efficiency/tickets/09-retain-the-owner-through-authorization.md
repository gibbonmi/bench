# 09 Retain the prospective owner through durable authorization

Blocked by: 01-build-isolated-repair-kits.md, 02-shrink-local-adoption-fixtures.md, 03-consolidate-conformance-executions.md, 07-share-skip-scans-and-prove-dispatch-reuse.md, 08-materialize-only-the-requested-tree.md
Writes: internal/gate/engine.go, internal/gate/execution_probe_test.go, internal/gate/authorization/authorization.go, internal/gate/authorization/authorization_test.go
Covers: LTE38, LTE39, LTE40, LTE41, LTE42, LTE43, LTE44, LTE45

## What to build

Review chunk: LTE-C8.

Add a gate-owned execution-and-inspection operation and route authorization through it. Retain one artifact owner until execution and independent durable inspection finish.
Construct a fresh subject evaluation after execution. Reload the existing durable evidence record rather than accepting the execution result's cached inspection.
Keep standalone ExecuteTree and InspectTreeContext contracts. Keep attribution and opaque evidence-token production under their current owner.

Preserve exact tree, baseline runner, freshness, and baseline phase-schedule checks. A successful child process alone cannot authorize publication.
Keep cancellation, timeout, build refusal, inherited binary retention, cleanup ordering, and concurrent-owner recovery behavior.
Place focused lifecycle tests in the existing execution test file to avoid increasing the oversized gate directory's file count.

Reconcile the complete build's quality and cost evidence after the prior outcomes are integrated. This dependency gathers the final measurement report.
Report required-run timings, observation counts, named reds, restoration, and remaining uncertainty. Do not add a standalone full benchmark for the report.

## Acceptance

- [ ] One authorization constructs one prospective artifact owner through its final evidence inspection.
- [ ] Successful execution with missing or corrupted durable evidence does not return Green.
- [ ] Tree, baseline identity, and freshness mismatches still refuse prior evidence.
- [ ] Cancelled, timed-out, failed-build, and normal exits close the owned bundle.
- [ ] Standalone inspection still publishes an independently recoverable owner record.
- [ ] Concurrent authorization retains the live owner's artifacts, and killed landing recovery removes the dead owner's artifacts.
- [ ] The final report states measured costs beside retained behavioral and omission evidence.

## Verification

Run gate and authorization focused tests, including every existing prospective-owner failure path. Demonstrate a durable-reload omission red. BENCH_KIT must identify the candidate kit for system verification. Run the system suite through `bench test --check system` during required integration verification. Final landing uses the existing complete gate.

