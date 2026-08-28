# Single-source test knowledge

Blocked by: 03-confine-run-binaries-and-resolved-bases.md
Writes: internal/gate/prospectiveartifact/prospectiveartifact.go, internal/gate/prospectiveartifact/prospectiveartifact_test.go, internal/gate/prospective_owner_test.go, internal/gate/lane_test.go, internal/systemtest/owner_artifact_recovery_test.go, internal/worktree/land_freshness_test.go

## What to build

Repair review findings S2, S3, S5, S7, and P3 from `reviews/prospective-artifact-recovery.md`.

Export the owner-record shape and the bundle prefix from the owner module and read both through it in every test.
Drop the dead `bench-gate-subject-` arm and the change-narrating comment.
Make the PAR04 second sweep prove that no path changed.

## Acceptance

- [ ] S2: one exported record shape serves every test that reads or plants a record.
- [ ] S3: the bundle prefix has one source.
- [ ] S5: the residue matcher names only prefixes a producer creates.
- [ ] S7: the lane test comment states the mutation only.
- [ ] PAR04: the second fresh sweep changes no path in the temporary root or the Git worktree list.
