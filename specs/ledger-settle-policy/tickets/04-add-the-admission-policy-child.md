# Add the admission policy child

Blocked by: 01-forbid-the-new-parent-adapters-in-the-policy-census.md, 03-add-the-ledger-transaction.md
Writes: internal/intent/admissionpolicy/ (new), internal/intent/intent.go, internal/intent/intent_test.go
Covers: LS15, LS16, LS17, LS19

## What to build

Verify the premise first. Read `Upsert` and `Compact` in
internal/intent/intent.go, and read `filterLive` and `claudeCandidates` in the
same file. Read `PutCleanupReceipt`, `PutAssignment`,
`ReauthorizeAssignment`, `compareAndSwapRequestDigest`, `PurgeAssignments`, and
`DeleteAssignment` in internal/intent/assignment.go. Read the transaction
ticket 03-add-the-ledger-transaction.md adds.

Read internal/worktree/landingpolicy/landingpolicy.go and
internal/worktree/landingpolicy/landingpolicy_test.go as the pure-policy shape.
Read `PolicyPackage` in internal/puritycensus/census.go. Read
`TestSnapshotProofLifecycle` and `TestUncorrelatedEntriesUseCandidateSet` in
internal/intent/intent_test.go.

Create the package `internal/intent/admissionpolicy` with the import path
`github.com/gibbonmi/bench/internal/intent/admissionpolicy`. Ticket
05-migrate-the-assignment-mutators.md imports that exact path, so do not rename
it. The package imports `internal/intent/ledger` alone, and it imports neither
`internal/git` nor `internal/intent`.

Hold the admission rule of each of the seven mutators in this package. The
seven mutators are `Upsert`, `Compact`, `PutCleanupReceipt`, `PutAssignment`,
`ReauthorizeAssignment`, `PurgeAssignments`, and `DeleteAssignment`. Each rule
answers the next ledger and the changed report from the supplied ledger and the
mutator's operand alone. Each rule takes no root, no path, and no repository.

Give the compaction rule typed liveness facts. The facts are one
candidate-probe boolean, one worktree-existence answer for each entry, and one
landed answer for each branch. The rule reads no filesystem and no Git.

Migrate `Upsert` and `Compact` onto the transaction and their rules. The intent
package translates `git.Worktrees`, `git.LocalBranches`, `os.Stat`,
`git.ResolvedDefault`, and `git.LandedInDefault` into the compaction facts at
its own boundary. Leave the five assignment mutators on their current envelope,
because ticket 05-migrate-the-assignment-mutators.md moves them. The seven
exported signatures do not change.

Write the census wrapper internal/intent/admissionpolicy/purity_census_test.go.
It runs `puritycensus.Scan` with `PolicyPackage()` and names this package's own
source through `MustHold`.

Write the two table tests in the new package. Name the first
`TestAdmissionPolicyDecidesEveryMutator`. It drives each of the seven rules
with a literal ledger value and observes the returned ledger and the changed
report.

Name the second table test `TestCompactionRuleReadsTypedLivenessFacts`. It
covers an entry whose worktree exists and an entry whose worktree is absent. It
also covers an entry whose branch landed. It also covers an uncorrelated Claude
entry with a true candidate probe. Neither table test builds a repository, and
neither one takes a lock.

## Acceptance

- [ ] The admission policy answers the next ledger and the changed report for each of the seven mutators from a literal ledger value.
- [ ] The policy table tests build no repository and take no lock.
- [ ] The compaction rule keeps an entry whose worktree exists.
- [ ] The compaction rule drops an entry whose worktree is absent and an entry whose branch landed.
- [ ] The compaction rule keeps an uncorrelated Claude entry when the candidate probe is true.
- [ ] `Snapshot` and `Compact` over a real repository return the same live entries as before the change.
- [ ] The policy census reports no import of `internal/git`, `internal/intent`, `internal/worktree`, `internal/bounds`, `os/exec`, or `syscall`, no ambient effect, and no parallel call.
- [ ] The policy census wrapper names this package's own source.
- [ ] The pre-existing `internal/intent` suite passes with its test logic unchanged.
- [ ] Self-probe: ignore the candidate probe in the compaction rule, and report the liveness test red.
