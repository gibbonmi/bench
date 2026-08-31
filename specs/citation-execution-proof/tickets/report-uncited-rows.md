# Report the uncited rows behind the review-owned marker

Blocked by: tighten-citation-grammar.md, refuse-mixed-tag-row-ids.md
Writes: internal/coverage/citations.go, internal/coverage/coverage.go, internal/coverage/coverage_command_test.go
Covers: CE23, CE24, CE25, CE26, CE27

## What to build

`bench coverage --check` names each mapped row that no citation backs. The
report is informational. It joins no violation list, and the exit code stays 0.

The row classification extends the one citation grammar in
`internal/coverage/citations.go`. That file already reads the seam cell through
`citationRe`. A row is cited when its seam cell holds at least one citation. A
row is exempt when its trimmed seam cell starts with the `review-owned:`
prefix. Every other mapped row is uncited.

The rendering lives in the `--check` arm of `Command`, in
`internal/coverage/coverage.go`. The report prints beside the pass line, when
the spec is mapped and the violation list is empty. The report is one bounded
line: the row count plus the row names. An opted-in map names each
row by its row ID. A non-opt-in map names each row by its row number. A
historical spec emits no report. A fully cited map emits no report.

Read the pass-line comment in `Command` before you write this arm. That comment
pins the pass line's contract. A pass is a definitive one-line result, and a
new state beside it needs its own line. Keep the pass line unchanged, and add
the report as a second output state beside it.

The Blocked-by edge serializes writes to `internal/coverage/citations.go` and
`internal/coverage/coverage.go`. It is not a semantic prerequisite.

## Acceptance

- [ ] CE23 — `--check` on a green non-opt-in map names each uncited row by its
      row number.
- [ ] CE24 — a seam cell with the `review-owned:` prefix is absent from the
      uncited report.
- [ ] CE25 — the uncited report leaves the check's exit code at zero.
- [ ] CE26 — an opted-in map's uncited report names the row by its row ID.
- [ ] CE27 — a row whose seam cell holds a citation is absent from the uncited
      report.
- [ ] a historical spec emits no uncited report.
- [ ] the report is one line that holds the row count and the row names.
- [ ] `bench gate` stays green.

## Delegate charge

You work in the Bench repo on the `citation-execution-proof` spec. Read
`specs/citation-execution-proof/spec.md` first. Then read
`internal/coverage/citations.go` and `internal/coverage/coverage.go` in full.
Also read `internal/coverage/citations_test.go` and
`internal/coverage/coverage_command_test.go` in full. Include the grammar arm
and the row-ID arm that the sibling tickets already landed.

Add the row classification beside the citation grammar in `citations.go`. Reuse
`citationRe`, `p.dataRows`, and the seam-cell projection. Do not add a second
map parser. Do not re-derive the map structure elsewhere.

Classify each mapped row. A row is cited when its seam cell holds at least one
citation. A row is exempt when its trimmed seam cell starts with
`review-owned:`. Every other row is uncited.

Render the report in the `--check` arm of `Command`, beside the pass line. Keep
the pass line's text and the exit code unchanged. Emit one bounded line: the
row count plus the row names.

Emit no report for a
historical spec. Emit no report when every mapped row is cited or exempt. Name
an opted-in map's rows by their row IDs. Name a non-opt-in map's rows by their
row numbers.

Do not change the message text of a delivered check. Do not add the report to a
violation list.

Add `TestUncitedRowReport` in `internal/coverage/coverage_command_test.go`.
Assert the exact report text for a non-opt-in map and for an opted-in map.
Assert that the exit code stays 0. Assert that a `review-owned:` row is absent
from the report. Give one fixture row a citation to a declared test file, and
assert that row is absent from the report. Drive a historical spec, and assert
no report appears.

Run only `bench worktree exec <label> -- go test ./internal/coverage/`. Do not
commit. Do not edit the spec.
