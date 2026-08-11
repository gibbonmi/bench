# Export the diff package's complete review-base resolution

Blocked by: none
Ownership fence: `internal/diff/`
Integration surfaces: recorded-key half→existing `resolveBase` in `internal/diff/diff.go`; fallback + error half→existing `resolveBranchRange` in `internal/diff/diff.go`; default branch→existing `internal/git/default_branch.go` `ResolvedDefault`; future consumer→implement-preflight-review.md
Contracts: the complete resolved review base — exported as `func ResolveReviewBase(root string) (base, method, errKind, errHint string)` in `internal/diff` (type: base commit-ish + method string + structured-error pair; method domain: `recorded`, `merge-base`, `merge-base (recorded sha unreachable)`, `merge-base (recorded sha not an ancestor)`; absence: a non-empty `errKind`/`errHint` with empty base when no default branch or merge-base resolves, never an empty base with all-empty returns) — will cross `internal/diff`→`internal/preflight` in implement-preflight-review.md, whose C8 row asserts consumption against this real producer; inside this ticket nothing crosses the fence yet
Closure: DX1/recorded-precedence, DX2/fallback-merge-base, DX3/unreachable-label, DX3/non-ancestor-label, DX4/unresolvable-error, DX5/branch-range-consumes-export

## What to build

Today no complete review-base resolution exists behind one symbol: unexported
`resolveBase` (`internal/diff/diff.go:267`) handles only the recorded
`branch.<name>.benchBase` key, while the default-branch merge-base fallback and
the cannot-resolve error semantics live inline in `resolveBranchRange`
(`diff.go:175`). Extract and export the *complete* resolution as
`func ResolveReviewBase(root string) (base, method, errKind, errHint string)` —
recorded key first, `git.ResolvedDefault` + merge-base fallback, the existing
cannot-resolve error kind/hint when neither answers — and make
`resolveBranchRange` consume the exported function, so a second consumer
cannot disagree with `bench diff` about the base, the method label, or the
failure shape. Pure expand: no behavior change, bare `bench diff` output stays
byte-identical, existing method-label spellings preserved exactly. Add
regression tests driving the exported function in a seeded throwaway repo
(mirror `internal/coverage/coverage_test.go` `TestCommand`'s
`t.Chdir(t.TempDir())` seeding shape) covering all four method labels and the
error path.

## Acceptance

- [ ] [DX1] (covers PF25) with a recorded `benchBase` key naming a reachable ancestor, the exported resolution returns that base with method `recorded`, in a seeded repo where merge-base would answer differently.
- [ ] [DX2] (covers PF25) with no recorded key, it returns the default-branch merge-base with method `merge-base`.
- [ ] [DX3] (covers PF25) an unreachable recorded sha and a non-ancestor recorded sha each fall back to merge-base carrying their exact loud method labels.
- [ ] [DX4] (covers PF25) with no resolvable default branch, it returns the existing cannot-resolve error kind and hint, never an empty base with a green return.
- [ ] [DX5] (covers PF25) `resolveBranchRange` consumes the exported function, every existing `internal/diff` test passes unchanged, and bare `bench diff` output is byte-identical before and after the export.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DX1/recorded-precedence | skip the recorded-key lookup and always merge-base | the recorded-key regression test | seed a benchBase key past a divergence, `go test ./internal/diff`, expect the precedence failure |
| DX2/fallback-merge-base | return an empty base instead of falling back when no key is set | the fallback regression test | seed no key, run, expect the missing-fallback failure |
| DX3/unreachable-label | collapse the unreachable-sha label into plain `merge-base` | the unreachable-sha label regression test | seed a recorded sha absent from the repository, run, expect the exact-label failure |
| DX3/non-ancestor-label | collapse the non-ancestor label into plain `merge-base` | the non-ancestor label regression test | seed a recorded sha that is reachable but not an ancestor of HEAD, run, expect the exact-label failure |
| DX4/unresolvable-error | return an empty base with empty error kind when no default branch resolves | the unresolvable regression test | seed a repo with no default branch, run, expect the missing-error failure |
| DX5/branch-range-consumes-export | leave `resolveBranchRange`'s inline fallback in place, bypassing the export | the fallback regression test asserted through `bench diff`'s own output plus a grep asserting the inline merge-base derivation is gone from `resolveBranchRange` | run the existing diff tests and the grep, expect the surviving-inline-derivation red |
