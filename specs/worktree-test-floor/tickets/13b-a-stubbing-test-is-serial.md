# A stubbing test is serial

Line: opus / medium.

A census rule is gate logic, so the mid tier applies.

Blocked by: 13-fixtures-pass-a-private-home.md
Writes: internal/worktree/parallel_census_test.go, and every `internal/worktree/*_test.go` the rule names (marks only)

## What to build

The census derives one more serial edge from the AST. A test is serial when
its body assigns a package-level identifier of a non-test file. A test-file
function it reaches counts the same way. The census reads the
non-test files for those declarations and never builds. The rule replaces
`bindGlobal` as the edge for a package-level swap; `bindGlobal` stays as the
effect record.

Every test the rule names drops `t.Parallel()`, so the live census turns
green. A swap of a package-level function variable then never overlaps a
parallel reader, because Go runs every serial test before the first parallel
one resumes.

## Acceptance

- [ ] A synthetic test that assigns a package-level identifier of a non-test file is serial under the census. One that also calls `t.Parallel()` is reported as the pair.
- [ ] WF01 stays green on the live tree.
- [ ] `go test -race -count=1 -parallel 2 ./internal/worktree` passes.
