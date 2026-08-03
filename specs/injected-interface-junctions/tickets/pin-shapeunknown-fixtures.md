# Pin the ShapeUnknown classifier sites and their specbuild composition

Blocked by: none
Ownership fence: `internal/worktree/classifier_shape_test.go`, `internal/specbuild/shape_junction_test.go`
Contracts: the ShapeUnknown enum value crosses `internal/worktree`→`internal/specbuild` through the real `ClassifyPathShape` into the precondition's ownership refusal, asserted by SU2 against the real classifier
Assumptions: production is presumed correct — a red against unmodified production is a stop-and-surface finding; fixtures are privilege-free — ENOTDIR via file-as-parent for site one, ELOOP via self-symlink for site two, EACCES via chmod-0 dir (root-guarded, `t.Cleanup` restore) for site three; the consumer test uses the ELOOP self-symlink at the assignment path, not ENOTDIR; the discarded error at the consumer stays discarded (out of scope); claims re-derived from the tree at pickup

## What to build

Deterministic fixtures for all three `ShapeUnknown` return sites in
`ClassifyPathShape`, each asserting the enum plus a non-nil error, and one
specbuild test composing the real classifier's `ShapeUnknown` into the
ownership refusal by replacing the assignment path with a self-symlink.

## Acceptance

- [ ] [SU1] `ClassifyPathShape` returns `ShapeUnknown` with a non-nil error at all three sites: `<regular-file>/child` (ENOTDIR), a self-symlink (ELOOP), and `<chmod-0 dir>/.git` (EACCES, skipped as root with the skip reason stated).
- [ ] [SU2] a specbuild precondition path whose assignment checkout is a self-symlink is refused with the ownership refusal, driven through the real classifier.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SU1 | change one site to return `ShapeAbsent, nil` | the classifier test | apply the mutation, run `go test ./internal/worktree -run <SU1 test>`, expect the site's assertion to fail |
| SU2 | short-circuit the consumer's default branch to nil | the junction test | apply the mutation, run `go test ./internal/specbuild -run <SU2 test>`, expect the missing refusal to fail the test |
