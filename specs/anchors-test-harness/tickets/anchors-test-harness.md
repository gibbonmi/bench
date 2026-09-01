# Collapse the repeated evaluate closures in the anchors registry test into one harness

Blocked by: none
Writes: internal/anchors/registry_data_test.go (and at most one new sibling test helper file in internal/anchors)

## What to build

`internal/anchors/registry_data_test.go` holds 18 hand-rolled
`evaluate := func(t *testing.T, broken int) []string` closures. Each one writes
every rule's file, section, and needle into a `t.TempDir()`, skips the needle at
index `broken`, and returns `EvaluateGroup(root, <group>)`. Each test then runs
the conformant pass, a per-rule red, and a cross-talk check. This is a fixture
harness pasted N times, which the repo's one-source-per-fact rule names as
duplicated knowledge.

Replace the closures with one shared test harness, parameterized by the rule set
and the group. Two groups are used today; one variant handles a `forbidden`
needle. Keep every per-test expectation independent: the needle, the wanted
diagnostic, the file, and the section stay in each test's own data. Only the
loop mechanism moves. Do not change `registry_data.go`, `match.go`, or any
canary fixture. Do not weaken any assertion: every rule that was checked red
before is checked red after, and the cross-talk check stays.

## Acceptance

- [ ] `internal/anchors/registry_data_test.go` contains no more than one `EvaluateGroup` call site inside a harness, and no per-test `evaluate := func` closure remains.
- [ ] `go test ./internal/anchors/ -parallel 2` is green, and the test count (from `go test -v` names) is not lower than before the change.
- [ ] Deletion of one registry needle from `registry_data.go` turns the anchors test red at the same test name as before the change. Restore the needle afterwards from `HEAD`.
- [ ] `go vet` and `gofmt` are clean on the edited files.
