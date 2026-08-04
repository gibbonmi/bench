# Converge an interrupted retire or discard on re-run

Blocked by: add-recovery-discard.md
Ownership fence: `internal/worktree`
Contracts: the named lifecycle step reaching `ApplyRecovery` crosses `internal/worktree/ownership.go`→`internal/worktree/resume.go`, asserted by IR1 and IR2 through the package's real fault seam driving real `ApplyRecovery` calls rather than a hand-edited intent record
Assumptions: `add-recovery-discard.md` has landed, so both verbs exist and share the delete-then-close ordering; the fault seam is the package's existing `Fault`/`LifecycleStep` mechanism, extended to reach this verb rather than replaced; claims re-derived from the tree at pickup

## What to build

Both recovery verbs delete the recovery ref before closing the assignment row, so
a process that dies between the two leaves a row naming a ref that no longer
exists. Re-running either verb against that state must close the row rather than
refuse — otherwise the crash window strands a row no other command can reach.

`ApplyRecovery` makes no `hit` call today and takes no fault, unlike
`applyCleanupTransaction` which threads one. Reach the existing named-step seam
into this verb: add the step constant that names the boundary between the ref
delete and the row close, and inject the fault there. The interruption is then
produced by the package's real mechanism, not simulated by hand-editing the
intent record — which is what makes the coverage evidence about production
behavior rather than about the fixture.

This holds for `--apply` and `--discard` alike. The `--apply` window is
pre-existing and was never covered; adding the discard coverage without it would
leave the identical latent defect in the older verb.

## Acceptance

- [ ] [IR1] a discard interrupted at the named step after the ref delete leaves the ref gone and the row open, and re-running the discard closes the row and succeeds.
- [ ] [IR2] an `--apply` interrupted at the same named step leaves the same state, and re-running the apply closes the row and succeeds.
- [ ] [IR3] the interruption is produced by installing the package's existing fault seam against the new step, and an uninstalled fault leaves both verbs' happy paths unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| IR1 | make the re-run refuse when the named ref no longer resolves | the interrupted-discard convergence test | return the missing-ref error before the compaction, run `go test ./internal/worktree -run Interrupted -timeout 120s`, expect the closed-row assertion to fail |
| IR2 | close the row only when the ref delete reported a deletion this invocation | the interrupted-apply convergence test | guard the compaction on the delete result, run `go test ./internal/worktree -run Interrupted -timeout 120s`, expect the closed-row assertion to fail |
| IR3 | drop the `hit` call at the new step so the seam never fires | the fault-reaches-recovery test | remove the call from `ApplyRecovery`, run `go test ./internal/worktree -run Interrupted -timeout 120s`, expect the injected-error assertion to fail |
