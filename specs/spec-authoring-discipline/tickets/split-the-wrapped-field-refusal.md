# Split the wrapped-field refusal from the unknown-field refusal

Blocked by: add-a-sources-example-to-the-map-template.md
Writes: internal/maps/validation.go, internal/maps/maps_graph_test.go, tests/canary/decision-map-integrity, tests/canary/decision-map-integrity/source-wrapped-field (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SAD5, SAD6, SAD9

## What to build

A Sources record line that holds no field separator is a wrapped continuation.
That case reds with a new message, and the message names the one-physical-line
rule. The message speaks about a Sources record alone. A wrapped bullet stays
legal in a terminal list, so the message must not claim that wrapping is illegal
everywhere.

Every other record line keeps the current message word for word. `tests/canary/decision-map-integrity/source-second-locator`
pins those bytes, and `internal/maps/maps_graph_test.go` pins them again at lines
270 to 276. A blanket rewrite of the shared message reds both readers.

The graph walk keeps the resolved-blocked edge rule.
`tests/canary/decision-map-integrity/graph-resolved-on-unresolved` still bites
when a resolved ticket stays blocked by an unresolved ticket. The sibling ticket
stated that rule in the template, so this ticket proves the walk still enforces
it.

This ticket is the last one that writes `tests/canary/decision-map-integrity`,
so it carries the invariant of that whole fixture family. Every fixture in the
family must still bite through `internal/conformance/fixture_bite_test.go`.

The five `cmd/bench` and `internal/conformance` entries in `Writes:` are the
registry closure of the `internal/maps` package. Edit one only when your change
reaches it.

## Acceptance

- [ ] SAD5 — a Sources record line with no field separator reds with the one-physical-line rule.
- [ ] SAD6 — a second `URL:` locator line reds with the exact current message.
- [ ] SAD9 — a resolved ticket blocked by an unresolved ticket still reds through the graph walk.
- [ ] The new `source-wrapped-field` fixture reds, and its restore returns green.
- [ ] A wrapped bullet in a terminal list stays legal.
- [ ] Every fixture under `tests/canary/decision-map-integrity` still bites.

## Delegate charge

You work in the Bench repo on the `spec-authoring-discipline` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/spec-authoring-discipline/spec.md` first. Then read
`sourceDiagnostics` at lines 134 to 202 in `internal/maps/validation.go`. Read the
unknown-field message at lines 154 to 156. Read
`TestMapSourcesRequireExactRecordShape` at line 264 in
`internal/maps/maps_graph_test.go`. Read
`TestMapTerminalContinuationAndEmptyAnswer` at lines 227 to 232. Read
`tests/canary/decision-map-integrity/source-second-locator/` as the fixture prior
art.

Your blocker writes `internal/maps` and runs that test package, and this ticket
does the same. At most two delegates run tests at once. Start only after your
blocker releases the frontier.

Split the one code path into two messages. Give the wrapped continuation the new
message. Leave every other line on the current message, byte for byte.

Add `tests/canary/decision-map-integrity/source-wrapped-field/` with `BASE`,
`EXPECT`, and `MUTATE.json`. Copy the exact diagnostic text the validator emits.
Do not paraphrase it.

Coverage rows: SAD5, SAD6, SAD9. Show each row red before your edit. Show each row
green after. Return the red-to-green log.

Self-probe with an omission mutation. Return the wrapped line to the shared
message and report the observed result. If the mutation returns green, add the
missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/maps/ ./internal/conformance/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
