# Cite each existing test seam in the resolvable form

Blocked by: none
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go
Covers: none

## What to build

The ft336 coverage map cited `TestWorktreeExecGrammar` as the existing seam for row BO31. The
tree has no such test. The cell named the test in prose, so `bench coverage --check` did not
read it as a citation, and build preflight did not resolve it. Review found the defect after
the build.

The coverage check already resolves a seam-cell citation of the form
`` `<path>_test.go` (`<Name>`) `` against the tree. It reports a missing file or an undeclared
name as a violation, and build preflight refuses the spec on that violation. A prose cell
holds no citation, so the check ignores it by design.

Widen the existing-test rule in the map discipline, which owns the rules for each row. A seam
cell that names an existing test cites that test in the citation form. Then the check and build
preflight resolve the name before the spec goes to review. The ticket-slicing anchor rows pin the
sentence, so a removal turns the anchor check red.

## Acceptance

- [ ] The map discipline tells the author to cite an existing test seam in the citation form that `bench coverage --check` resolves.
- [ ] No other guidance file states how a seam cell names an existing test.
- [ ] A probe that removes the sentence turns the anchor check red.
