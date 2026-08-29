# Review pickup: ft213-craft-delegate-visit

Frozen base `129fcbc1`, reviewed tip `eba49e0e`. Raw findings: Standards 5,
Spec 3, Coverage 4. Repair targets after de-duplication: 5.

## Standards

Count 5. Worst: the new anchors test comment narrates the change.

- S1 (auto-fix): `internal/anchors/registry_data_test.go` — the doc comment on
  `TestCraftDelegateDisciplineAnchorsRedOnRemoval` names the FT213 visit, which
  `craft-comments` bars as provenance. Reword it to state what the pins hold.
  The FT214 sibling comment has the same shape; reword it in the same pass.
- S2 (no-op): "wrap width" is a per-file value the charge states, not a term.
- S3 (no-op): the go-back-for-reds bullet adds the remedy branch.
- S4 (no-op): the "The charge" paragraph sits at six sentences; pre-existing.
- S5 (no-op): `bench-write-spec.md:66` is 151 columns; no wrap convention is
  documented and the file sits at its budget.

## Spec

Count 3. Worst: the refusal does not resolve symlinks (repaired under C1).

- P1 (no-op): the `committed` line rule names the coordinator, the correct actor.
- P2 (no-op): the coverage-row citation rule applies to every charge.
- P3 (no-op): folded into C1.

## Coverage

Count 4. Worst: a symlinked root spelling escapes the `dst == root` refusal.

- C1 (auto-fix): `internal/canary/mutation.go` compares `filepath.Abs` paths
  only. Resolve both paths with `filepath.EvalSymlinks` when they exist, and add
  a test that spells `root` through a symlink.
- C2 (auto-fix): `TestRestoreMutationFixtureRefusesDestinationEqualToRoot`
  places `added.txt` in `root`, so the removal branch never runs without the
  guard. Keep the overlay-only file out of `root` and assert it survives.
- C3 (auto-fix): nothing grades the existence of
  `references/delegation-discipline.md`. Add one `Require` pin on its lead
  sentence; give `references/map-discipline.md` the same pin.
- C4 (auto-fix, partial): add a `Forbid` pin on `Read-only delegations need no
  worktree` in `SKILL.md` and a `Require` pin on the cap sentence in
  `bench-write-spec.md`. The two probe sentences stay unpinned.
