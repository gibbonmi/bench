# 5. Summarize the green preflight checks on one line

Blocked by: none
Writes: internal/preflight/command.go, internal/preflight/command_review_test.go, internal/preflight/source_tip_test.go, internal/preflight/verdict_summary_test.go (new)
Covers: BO37, BO38, BO39, BO40

## What to build

`bench preflight build <slug>` and `bench preflight review <slug>` keep the `phase`, `spec`, and `source` lines. The check table becomes one line: `checks{green=<n>,not_applicable=<n>,red=<n>}`. When one or more checks are red, a `checks[<n>]{check,verdict,detail,next}` table follows with the red rows only. The exit code stays 1 for a red verdict and 0 otherwise. The charge forms keep their complete check table in the evidence artifact.

## Acceptance

- [ ] An all-green build preflight prints `checks{green=13,not_applicable=2,red=0}` and no `checks[` table.
- [ ] An all-green review preflight prints the same summary line and no `checks[` table.
- [ ] A preflight with one red check prints the summary line and `checks[1]{check,verdict,detail,next}` with the red row only.
- [ ] A red preflight exits 1.
