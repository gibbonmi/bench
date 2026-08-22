# Write the drain rules and bind the token table

Blocked by: 01-report-the-board-flow.md, 02-parse-the-next-token.md, 03-require-the-feeds-marker-on-retro-items.md
Writes: .agents/commands/bench-drain.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/conformance/registry/registry.go, internal/conformance/checks.go, internal/conformance/checks_test.go, internal/conformance/registry_test.go, internal/conformance/row_next_grammar_test.go, tests/canary/row-next-grammar/, projects/benchkit.md

## What to build

`/bench-drain` states the rules the map decided: the exit quotes
`bench roadmap --flow`; an entry feeds a row only when it changes the row's
priority, scope, or `Next:`; an occurrence-only entry is dismissed with one
line of why; a new row needs a `Next:` token and a class; a positive net delta
forces reducing moves in the next batch diff; a light-path item is built in the
session by default. The command carries the token table, and a new Dev check
`row-next-grammar` compares that table against the exported token set. One
anchors-registry entry pins each rule. Spec group D, rows RF16, RF26, RF27.

## Acceptance

- [ ] Removing any one of the six drain rules yields the anchor diagnostic that names the dropped rule.
- [ ] A token table that lacks `kit-edit` yields the drift diagnostic naming that token.
- [ ] The canary fixture for `row-next-grammar` makes the check red, and restoring the token clears it.
