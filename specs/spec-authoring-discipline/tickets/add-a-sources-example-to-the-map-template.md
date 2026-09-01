# Add a Sources example to the decision-map template

Blocked by: none
Writes: internal/maps/schema.go, internal/maps/maps_parse_test.go, internal/conformance/decision_map_integrity_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SAD1, SAD2, SAD3, SAD4, SAD7, SAD8

## What to build

`bench maps --template` renders one complete Sources entry under the Sources
terminal heading. The entry starts with a `URL:` locator. A `Path:` locator
routes to `validateSourcePath`, which resolves against the repository root and
finds no placeholder file there. The entry writes the Supports line and the
Drift line as two separate physical lines.

The template states that a resolved decision ticket cannot stay blocked by an
unresolved ticket. The graph walk stays the one enforcement of that rule. Add
no second check. The template keeps exactly one decision ticket.

The rendered template validates with zero diagnostics as a shaping map and as a
ready map. `internal/conformance/decision_map_integrity_test.go` drives the
rendered template through the live validator. Eleven test call sites resolve the
`<answer>` placeholder with a single-count replacement, so a second decision
ticket reds them.

The five `cmd/bench` and `internal/conformance` entries in `Writes:` are the
registry closure of the `internal/maps` package. Edit one only when your change
reaches it.

## Acceptance

- [ ] SAD1 — the rendered template holds one Sources bullet under the Sources heading.
- [ ] SAD2 — the Sources bullet starts with a `URL:` locator.
- [ ] SAD3 — the Supports line and the Drift line are two physical lines.
- [ ] SAD4 — the rendered template validates with zero diagnostics as shaping and as ready.
- [ ] SAD7 — the rendered template holds exactly one decision-ticket heading.
- [ ] SAD8 — the rendered template states the resolved-blocked rule.
- [ ] The eleven single-count answer replacements stay green.

## Delegate charge

You work in the Bench repo on the `spec-authoring-discipline` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/spec-authoring-discipline/spec.md` first. Then read
`canonicalDecisionMapSchema` and `DecisionMapTemplate` in
`internal/maps/schema.go`. Read `sourceDiagnostics` and `validateSourcePath` in
`internal/maps/validation.go`. Read `TestParseDecisionMapSchemaAndTemplate` in
`internal/maps/maps_parse_test.go`. Read
`internal/conformance/decision_map_integrity_test.go`.

Render the Sources example as body text under the Sources terminal heading. Use
a `URL:` locator. Keep each record field on one physical line.

State the resolved-blocked rule as one template line. Do not add a second
validator check for it.

Coverage rows: SAD1, SAD2, SAD3, SAD4, SAD7, SAD8. Show each row red before your edit. Show
each row green after. Return the red-to-green log.

Self-probe with an omission mutation. Join the Supports line and the Drift line
into one line, then report the observed result. If the mutation returns green,
add the missing assertion.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/maps/ ./internal/conformance/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
