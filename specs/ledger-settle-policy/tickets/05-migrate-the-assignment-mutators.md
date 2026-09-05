# Migrate the assignment mutators onto the transaction

Blocked by: 04-add-the-admission-policy-child.md
Writes: internal/intent/assignment.go, internal/intent/intent.go, internal/intent/intent_test.go, internal/intent/purge_test.go, internal/intent/transaction.go, internal/intent/transaction_test.go (new), internal/intent/admissionpolicy/admissionpolicy_test.go, CONTEXT.md, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary
Covers: LS1, LS2, LS5, LS6, LS7, LS8, LS10, LS18, LS37

## What to build

Verify the premise first. Read `PutCleanupReceipt`, `PutAssignment`,
`ReauthorizeAssignment`, `compareAndSwapRequestDigest`, `PurgeAssignments`, and
`DeleteAssignment` in internal/intent/assignment.go. Read the transaction
ticket 03-add-the-ledger-transaction.md adds, and read the five assignment
rules ticket 04-add-the-admission-policy-child.md adds. Read `Read` and
`readPath` in internal/intent/intent.go.

Read
`TestPurgeAssignmentsDropsRecordsReadRefuses`,
`TestPurgeAssignmentsDropsEveryLegacySchemaRecord`, and
`TestPurgeAssignmentsIsANoOpWithoutALedger` in internal/intent/purge_test.go.
Read `TestReadEvidenceStates`, `TestIdenticalStampedAssignmentWritePreservesBytes`,
`TestCleanupReceiptWindowKeepsExactlyLast256Completions`, and
`TestCompareAndSwapRequestDigestRefusesConcurrentMovementAndPreservesOtherFields`
in internal/intent/intent_test.go. Read the glossary shape of an existing term
in CONTEXT.md.

Migrate the five assignment mutators onto the transaction and their rules. Each
mutator supplies the root, the read mode, and its decision closure. The five
exported signatures and the five external behaviors do not change. This ticket
is the last one that touches `internal/intent`, so it carries that package's
whole-package invariant.

`PurgeAssignments` is the one tolerant caller. It keeps its pre-lock
absent-file answer, so a converged re-run creates no lock file. Every other
mutator reads in the strict mode, which keeps refusing a ledger it cannot
account for.

`ReauthorizeAssignment` is the one caller that supplies a compensation step. It
supplies the rollback its `transition` returns. The step runs on the write
failure alone. Its closure's own error path runs the rollback for the
compare-and-swap failure, which keeps today's behavior. Every other mutator
supplies no step.

Keep exactly one real-filesystem journey for each of the seven mutators, and
move every other admission partition to the policy tables ticket
04-add-the-admission-policy-child.md adds. Delete no partition; move it.

Add the term **ledger transaction** to CONTEXT.md. Define it as the one
envelope over the intent ledger. State that it locks, reads in a named mode,
runs one decision closure, and writes a reported change. Give it the Avoid list
`lock helper` and `mutator envelope`. `bench anchors CONTEXT.md` names seven
required needles, and this entry touches none of them.

The eight fixture directories in the `Writes:` line are closure headroom. They
pin CONTEXT.md, and this ticket edits none of them.

Write the new tests in internal/intent/transaction_test.go. Name the first
`TestLedgerTransactionPreservesEveryMutatorOutcome`. It runs the seven mutators
over one fixture ledger and compares each result with the recorded pre-change
bytes.

Name the second new test
`TestReauthorizeCompensatesOnTheTerminalWriteFailure`. Its first case makes the
ledger directory read-only with a transition that returns a rollback. Its
second case drives `PutAssignment` over the same read-only directory. Add one
`Upsert` case to
`TestIdenticalStampedAssignmentWritePreservesBytes`.

## Acceptance

- [ ] Each of the seven mutators produces the same ledger file bytes over its own happy-path operand as the pre-transaction implementation produces.
- [ ] An identical `Upsert` and an identical `PutAssignment` each leave the ledger bytes and the file modification time unchanged.
- [ ] With one valid and one undecodable assignment record, `PurgeAssignments` drops the undecodable record and keeps the valid one.
- [ ] With an unsupported schema number, `Read` returns an error that names the unsupported schema.
- [ ] With a rollback transition and a read-only ledger directory, `ReauthorizeAssignment` returns an error and the rollback ran exactly once.
- [ ] With a read-only ledger directory, `PutAssignment` returns an error and no compensation step runs.
- [ ] With no ledger file present, `PurgeAssignments` returns zero and creates no lock file.
- [ ] CONTEXT.md holds the term `ledger transaction` with its Avoid list.
- [ ] Each of the seven mutators keeps exactly one real-filesystem journey.
- [ ] The pre-existing `internal/intent` suite passes with its test logic unchanged, except the moved admission partitions.
- [ ] Self-probe: drop the compensation step from the transaction, and report the reauthorization test red.
