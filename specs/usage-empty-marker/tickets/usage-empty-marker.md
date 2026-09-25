# Give the empty-value refusal marker one owner

Blocked by: none
Writes: internal/usage/parse.go, internal/usage/parse_test.go, internal/commit/chain_grammar_test.go
Covers: none

## What to build

The usage parser refuses an empty positional and an empty value for a flag that
declares `NoEmptyValue`. Each refusal names the empty token as `""`. The parser writes
that marker inline twice and exports no owner. So the parser tests and the commit chain
test restate the marker.

The usage package exports the marker as one constant and exports one function that
names a flag with an empty value. The parser and the tests read these owners.

## Acceptance

- [ ] Each empty-argument refusal names its token through `usage.EmptyOperand` or `usage.EmptyFlagValue`. No other marker literal remains in `parse.go`.
- [ ] The parser tests and the commit chain test read the owner for the refused token.
- [ ] A probe that swaps the flag refusal back to the bare flag name turns the commit chain test red. The restore returns it to green.
- [ ] `bench test --package ./internal/usage` and `bench test --package ./internal/commit` pass, and the gate is green at the landing.
