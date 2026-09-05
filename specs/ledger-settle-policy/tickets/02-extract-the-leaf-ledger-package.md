# Extract the leaf ledger package

Blocked by: none
Writes: internal/intent/ledger/ (new), internal/intent/ledger_aliases.go (new), internal/intent/ledger_aliases_test.go (new), internal/intent/intent.go, internal/intent/assignment.go, internal/intent/validate.go, internal/intent/branches.go, internal/intent/intent_test.go, internal/worktree/lifecyclepolicy/lifecyclepolicy.go, internal/worktree/lifecyclepolicy/lifecyclepolicy_test.go
Covers: LS11, LS12, LS13, LS14

## What to build

Verify the premise first. Read `Schema`, `LegacySchema`, `Kind`, `Entry`,
`Ledger`, `Address`, `Read`, `readPath`, `NewEntry`, and `writePath` in
internal/intent/intent.go. Read `AssignmentRecordSchema`,
`assignmentBranchNamespace`, `AssignmentBranchPrefix`, `RecoveryRefNamespace`,
`AssignmentBranchRef`, `RecoveryRefPrefix`, `validAssignmentBranchRef`,
`AssignmentState`, `Recovery`, `Assignment`, `CleanupReceipt`,
`validateCleanupReceipts`, `idPattern`, `ValidIdentity`, `ValidateAssignment`,
and `RequestDigest` in internal/intent/assignment.go. Read `validEntry` in
internal/intent/validate.go. Read `LeafPackage`, `Scan`, and `MustHold` in
internal/puritycensus/census.go. Read
internal/worktree/landingpolicy/purity_census_test.go as the wrapper shape.

Create the package `internal/intent/ledger` with the import path
`github.com/gibbonmi/bench/internal/intent/ledger`. Every sibling ticket
imports that exact path, so do not rename it. Move the schema constants, the
seven types, their constants, the identity patterns, the branch-ref helpers,
`RequestDigest`, `ValidIdentity`, `ValidateAssignment`, `validEntry`, and
`validateCleanupReceipts` into it. Export `validEntry` and
`validateCleanupReceipts` under names the intent package calls, because a
sibling package cannot reach an unexported name. The package imports nothing
under `internal/`, so leave `Address`, `Read`, `readPath`, `writePath`, and
`NewEntry` in the intent package.

Re-export every moved name from the intent package. A type becomes a type
alias, and a function and a constant become a declared value. Every importing
file then keeps its `intent.` spelling and its import line. The tree holds 68
Go files that import `internal/intent`. Change none of them. The enumeration
command is `rg -l 'bench/internal/intent"' --type go`.

Shrink `internal/intent/assignment.go` under the 400-line budget. The file
holds 584 lines today, and the moved schema is the reduction this ticket owns.
Report to the coordinator when the move leaves the file over 400 lines.

Write the census wrapper internal/intent/ledger/purity_census_test.go. It runs
`puritycensus.Scan` with `LeafPackage()` and names this package's own source
through `MustHold`.

Write two new tests in internal/intent/intent_test.go. Name the first
`TestIntentReExportsEveryMovedLedgerName`. It assigns each alias to the leaf
declaration, so a dropped alias reds compilation and a renamed type reds the
build. Name the second `TestOnlyTheIntentOwnerImportsTheLeafLedger`. It walks
the module sources and refuses an import of the leaf path outside
`internal/intent`, `internal/intent/admissionpolicy`, and their test files.

## Acceptance

- [ ] Every moved name resolves through its `intent.` spelling and equals the leaf declaration.
- [ ] No Go file outside `internal/intent` and `internal/intent/admissionpolicy` imports the leaf package.
- [ ] The leaf census reports no import under `internal/` other than `internal/puritycensus`, no ambient effect, and no parallel call.
- [ ] The leaf census wrapper names this package's own source.
- [ ] `bench structure` reports `internal/intent/assignment.go` under the 400-line budget.
- [ ] The pre-existing `internal/intent` suite passes with its test logic unchanged.
- [ ] Self-probe: drop one alias, and report the re-export test red.
