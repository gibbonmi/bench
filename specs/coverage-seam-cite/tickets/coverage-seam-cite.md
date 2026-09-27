# Refuse a seam-cell test name outside the citation form

Blocked by: none
Writes: internal/coverage/citations.go, internal/coverage/citation_form_test.go (new), .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, specs/ft290-test-projection/spec.md, specs/test-determinism/spec.md
Covers: none

## What to build

`bench coverage --check` reads each seam cell of a mapped spec for Go test names.
A test name is a token that matches `Test[A-Z][A-Za-z0-9_]*`.
A test name in the citation form, `` `<path>_test.go` (`<Name>`) ``, keeps the current resolution against the tree.
A test name outside that form is a violation.
The diagnostic names the row, the test name, and the citation form to use.

Examples of a name outside the form are a bare name and a name beside a path with no name list.
A name without backticks in the name list is also outside the form.
A seam cell that holds the word `planned` names a test that does not exist yet.
That cell keeps the current behavior: it gets no violation, and the uncited line names its row.
The uncited line and the historical opt-out do not change.
The map discipline reference states the rule and the `planned` marker.

Seven seam cells in two staged specs name existing tests outside the form.
This ticket rewrites those cells in the form and keeps the other words of each cell.

## Acceptance

- [ ] A seam cell `covered by TestMissing` makes `bench coverage --check` exit 1 with `coverage map row 1 names 'TestMissing' outside the citation form`.
- [ ] A seam cell `` `internal/x/foo_test.go` (TestMissing) `` gives the same violation.
- [ ] A seam cell `` `internal/x/foo_test.go` TestPresent `` gives the mention violation and a violation for `TestPresent`.
- [ ] A seam cell `planned TestMissing in internal/x` gives no violation, and the uncited line names its row.
- [ ] A seam cell `unplanned TestMissing` gives the violation.
- [ ] A test name in the why cell gives no violation.
- [ ] A declared citation, subtest included, gives no violation.
- [ ] `bench coverage --check` exits 0 on each staged spec.
