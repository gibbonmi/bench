# Fix the label-line rule in the prose grader

Blocked by: none
Writes: internal/prose/parse.go, internal/prose/parse_test.go
Covers: LP1, LP2, LP3, LP5

## What to build

Verify the premise first: `isLabelLine` in internal/prose/parse.go tests only the
prefix before the first colon, and `gradeBlocks` splits the paragraph before it
tests the line for a terminator. Then change the rule. A line is a label when
its one-to-four-word prefix is a template field name, or when the whole line
carries no sentence terminator. Keep the field-name list as a closed
constant beside the abbreviation list. A label-shaped prose line with a
terminator stays inside its paragraph, so the six-sentence bound counts it.

Add the table cases: a label-shaped line with a mid-line terminator and no
trailing terminator inside an eight-sentence run reds; twenty `Occurrence:`
lines that end with a period stay green; a `Writes:` line stays green. Run the
live-tree prose check and report it green.

## Acceptance

- [ ] `TestFindings` reds an eight-sentence run that holds `Run the real path: gate it. Then land` on one line.
- [ ] `TestFindings` reports no finding for twenty `Occurrence:` lines that end with a period.
- [ ] `TestFindings` reports no finding for a `Writes:` line with no terminator.
- [ ] `TestProseMechanicsHoldsOnTheLiveTree` passes on the worktree.
- [ ] Self-probe: change the rule to test only the last token, and report the LP1 case green as the observed wrong result.
