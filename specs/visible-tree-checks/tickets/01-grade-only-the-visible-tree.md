# Grade only the visible tree in the gofmt step and in-place checks

Blocked by: none
Writes: CHANGELOG.md, internal/visibletree (new), internal/gate/gate_go.go, internal/gate/gate_go_test.go (new), internal/gittest/gittest.go, internal/testrepo/working_tree.go, internal/intent/ledger_aliases_test.go, internal/freshness/publication_topology_test.go, internal/conformance/data_handling_test.go, internal/conformance/prose_budget_test.go, internal/conformance/skill_description_budgets_test.go, internal/conformance/visible_tree_test.go (new), internal/skillsindex/skillsindex.go, internal/skillsindex/ignored.go (new), internal/skillsindex/ignored_test.go (new)
Covers: none

## What to build

An ordinary `bench gate` runs its phases in the live checkout. Ignored files stay in that checkout, but they are not part of the subject tree that the gate grades. The gofmt step and five in-place checks must read only the visible tree: tracked files that exist, and untracked files that git does not ignore. Outside a repository, git ignores nothing, so the checks read every file.

The package `internal/visibletree` is the one owner of this definition. These checks use it:

- the gofmt step
- the leaf-ledger import census
- the freshness publication topology
- the control-rune escaper census
- the shared skills listing of the two budget checks
- the skills index

The test scaffold `testrepo.CommitWorkingTree` also uses it, instead of its own listing.

The reviewer approved this fix on 2026-10-16, outside the active commitment.

## Acceptance

- [ ] An ordinary `bench gate` stays green when ignored copies of Go source and an ignored local skill are in the checkout.
- [ ] Each changed check ignores an offender at an ignored path, and reports the same offender at a visible untracked path.
- [ ] A mutation that makes `internal/visibletree` ignore nothing turns each new regression test red.
