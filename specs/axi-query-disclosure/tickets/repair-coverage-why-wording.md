# Repair coverage why wording

Blocked by: none
Writes: `internal/coverage/coverage.go`, `internal/coverage/coverage_test.go`

## What to build

Fix the coverage disclosure `why` cell (review finding R8): `coverage.go:444`
builds `"check coverage row "+row[0]` where `row[0]` is the STORY column, so
the emitted text labels a story list as a row (traced: `check coverage row
1,2` for map row QD6). Reword the `why` so it no longer mislabels the story
cell; the `cmd` cell and all pre-disclosure bytes are untouched. Candidate-side
test expectations in `coverage_test.go` are re-cut to the new wording; no
testdata changes (the `pre-disclosure-*.stdout` fixtures carry no help block).

## Acceptance

- [ ] [CW1] (covers local) the mapped-rows `why` cell no longer contains the
  word "row" applied to the story value; the new wording names what the value
  is, and duplicate templates still dedupe to one action per identical
  cmd+why.
- [ ] [CW2] (covers QD6) every pre-disclosure fixture byte, stream, and exit
  is unchanged; only candidate-side help-block expectations differ.
