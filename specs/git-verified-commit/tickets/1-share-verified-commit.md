# Share the verified commit query

Blocked by: none
Writes: internal/diff/range.go, internal/gate/authorization/authorization.go, internal/gate/checkpoint.go, internal/gate/greenmarker/greenmarker.go, internal/git/git.go, internal/landing/attribution.go, internal/landing/composition.go, internal/landing/gitexec.go, internal/landing/landing.go, internal/landing/merge.go, internal/preflight/gather.go, internal/worktree/clean.go, internal/worktree/clean_classes.go, internal/worktree/clean_landed.go, internal/worktree/land_identity.go, internal/worktree/land_resume.go, internal/worktree/merge.go, internal/worktree/ownership.go, internal/worktree/reauthorize.go, internal/worktree/reset.go, internal/worktree/reset_envelope.go, internal/worktree/reset_restore.go, internal/worktree/subshell.go, internal/git/refs_test.go
Covers: none

## What to build

Promote the repeated verified commit query into the existing Git owner for FT302.
Each caller keeps its revision, flags, output, error handling, and side effects.
The helper owns the commit peel and verification arguments.
Other Git operations retain their existing owners.

## Acceptance

- [x] Each enumerated Go query uses one Git helper.
- [x] Branches, commit IDs, and annotated tags produce the same commit IDs as the original queries.
- [x] Missing refs, non-commit objects, option-shaped revisions, and invalid roots preserve the original outputs and errors.
- [x] Quiet and end-of-options flags retain their original positions before the revision.
- [x] Existing caller tests pass without changed expectations.

## Verification

Run affected packages and root conformance.
Compare the helper against the original Git queries across the enumerated revision and flag families.
Omit the commit peel and the optional flags in separate probes.
Require each probe to fail the differential test.

The affected selection passed all 52 packages, including root conformance.
The peel omission probe returned `bit` with nine failed cases.
The flag omission probe returned `bit` with eight failed cases.
Both probes restored the source.
The differential test retains the original queries because these two demonstrated mutations require an independent expectation.
