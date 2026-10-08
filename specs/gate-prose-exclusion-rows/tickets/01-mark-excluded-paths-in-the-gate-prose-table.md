# Mark an excluded path as excluded in the gate-prose table

Blocked by: none
Writes: internal/gate/gate_prose.go, internal/gate/gate_prose_test.go, internal/prose/subject.go, internal/prose/prose_test.go
Covers: none

## What to build

The `bench gate-prose` verb prints a `prose[N]{path,verdict}` table when the graded list is clean. The exclusion list in `.bench/prose-exclusions` removes a path from the grade. At this time, the table gives a `pass` row to an excluded path, and a reader can think that the gate graded that path.

The verdict of an excluded path must be `excluded`. The verdict of a graded clean path stays `pass`. A graded path with a finding keeps its diagnostic line. The named form and the staged form both use this rule.

The prose package makes the decision to grade or to exclude a subject. That one decision must supply the verdict for each row. The renderer must not do a second exclusion check.

## Acceptance

- [ ] `bench gate-prose <root> -- <path>` prints the row `<path>,excluded` for a path that the exclusion list excludes. It exits 0 when every named path is excluded.
- [ ] `bench gate-prose <root> -- <excluded path> <clean path>` prints `excluded` for the first path and `pass` for the second path, and it exits 0.
- [ ] `bench gate-prose <root> --staged` prints `excluded` for a staged path that the staged exclusion list excludes.
- [ ] A clean graded path keeps the `pass` row, and a path with a finding keeps its diagnostic line and exit 1.
- [ ] A recorded probe that changes the excluded verdict back to `pass` makes the focused test fail.
