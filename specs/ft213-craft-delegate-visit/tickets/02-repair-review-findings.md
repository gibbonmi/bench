# Repair the review findings

Blocked by: 01-fold-ft213-into-craft-delegate.md
Writes: internal/canary/mutation.go, internal/canary/mutation_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, reviews/ft213-craft-delegate-visit.md (delete)

## What to build

The five repairs `reviews/ft213-craft-delegate-visit.md` names: S1, C1, C2,
C3, and C4, each red-first where a test can go red. The same commit deletes
`reviews/ft213-craft-delegate-visit.md`.

## Acceptance

- [ ] `RestoreMutationFixture` refuses a `root` spelled through a symbolic link when `dst` is the real directory, and the test shows that row red before the edit.
- [ ] With the guard removed, the refusal test fails because an overlay-only file is deleted from `root`.
- [ ] A `Require` pin grades the existence of `references/delegation-discipline.md` and of `references/map-discipline.md`, and the independent test shows each red on removal.
- [ ] A `Forbid` pin rejects `Read-only delegations need no worktree` in `SKILL.md`, and a `Require` pin holds the cap sentence in `bench-write-spec.md`; the independent test shows each red.
- [ ] No comment in `internal/anchors/registry_data_test.go` names a roadmap item or a visit.
