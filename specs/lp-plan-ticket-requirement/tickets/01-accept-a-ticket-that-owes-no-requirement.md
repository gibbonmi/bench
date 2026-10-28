# Accept a ticket that owes no verification requirement

Blocked by: none
Writes: internal/reviewrecord/plan.go, internal/reviewrecord/delegated_test.go
Covers: none

## What to build

The plan grade of a version 2 completion plan accepts a chunk ticket that owes no verification requirement. Before, the grade refused that plan with "every ticket owes one". The promotion broker grades the completion evidence of each landing, so the installed broker must hold this rule before the FT463 build-cost landing. That build retires the rule in its ticket 01, and one of its tickets owes no requirement.

This change takes only the code lines of that retirement, with the same bytes, so the later fold of the default branch into the build merges with no conflict. The doc comment of `requirementsValid` stays unchanged, because the build rewrites those lines, and the build-cost landing replaces it. The build-cost landing also carries the `CHANGELOG.md` entry and the plan test of this behavior (BC26 in `TestPlanProbeOnlyChunks`). After the fold, the build tree is unchanged, so its recorded final evidence stays current.

A chunk requirement still names a ticket of its own chunk, and the final list must still be non-empty.

## Acceptance

- [ ] A version 2 plan whose chunk holds a ticket with no verification requirement passes the plan grade.
- [ ] A chunk requirement that names a ticket outside its chunk still refuses with "outside this chunk".
