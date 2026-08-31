# Register the ticket-grammar conformance sweep

Blocked by: close-the-ownership-closures.md
Writes: internal/conformance/ticket_grammar_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/registry.go, projects/benchkit.md
Covers: TG17, TG18, TG25, TG27, TG28, TG36

## What to build

A registered check named `ticket-grammar` grades ticket files inside the gate.
It runs at the dev tier over the graded root. It sweeps every
`specs/*/tickets/` directory with the parser the preflight rows already use. It
is fail-closed.

The sweep tolerates two tree states. A staged spec with no tickets directory
stays green. A tickets folder with no `spec.md` is graded with the `Covers:`
checks skipped, so a light-path landing needs no declared rows.

The same check grades the binding table. A bound registry file that the tree
does not hold reds the check. An approved query command package with no binding
row also reds the check. The check also proves the dispatcher, renderer, and
terminal-lifecycle owner rows are present.

The profile advertises the new input binding. The registry order and the
profile row order must agree, because the meta check compares them. This ticket
adds no canary family, so the family binding waits for the next ticket.

## Acceptance

- [ ] TG17 — a binding row that names a nonexistent registry file reds the check.
- [ ] TG18 — an approved query package with no binding row reds the check.
- [ ] TG25 — a tickets folder with no `spec.md` is graded without the `Covers:` checks.
- [ ] TG27 — a staged spec whose ticket carries a dangling blocker reds the sweep.
- [ ] TG28 — a staged spec with no tickets directory stays green.
- [ ] TG36 — a deleted owner row reds the check.
- [ ] `bench test --check ticket-grammar` reaches the check.

## Delegate charge

You work in the Bench repo on the `ticket-grammar` spec. Read
`specs/ticket-grammar/spec.md` first. Then read
`internal/conformance/registry/registry.go` and
`internal/conformance/checks_test.go` in full. Read
`internal/conformance/row_next_grammar_test.go` for the check shape. Read
`internal/conformance/docs_workflow_checks_test.go` for the sweep shape.

Create `internal/conformance/ticket_grammar_test.go` and write
`checkTicketGrammar(root string) []string` in it. Sweep every
`specs/*/tickets/` directory through `internal/tickets`. Return one diagnostic
per fault. Keep the check fail-closed on an unreadable directory.

Return no diagnostic for a staged spec with no tickets directory. Skip the
`Covers:` checks for a tickets folder that holds no `spec.md`. Keep every other
grammar in force there.

Grade the binding table in the same check. Red a bound registry file the tree
does not hold. Red an approved query command package that no binding row
covers.

Register the check in `registry.Checks` at the dev tier with the root subject
and the `catch-all` input source. Add no new `InputSource` constant.
Add the dispatch entry in `checks_test.go`. Add the matching row to the
profile's check-input table in `projects/benchkit.md`. Keep both orders in
agreement. Do not add a `familyChecks` entry, because the fixtures land in the
next ticket.

Add `TestTicketGrammarSweepRedsDanglingBlocker`,
`TestTicketGrammarSweepToleratesAbsentTicketsDir`,
`TestTicketGrammarSweepSkipsCoversWithoutSpec`,
`TestTicketGrammarRedsMissingRegistryFile`, and
`TestTicketGrammarRedsUnboundQueryPackage`, and
`TestTicketGrammarRedsMissingOwnerRow` in `internal/conformance`. Build a
synthetic root in each test. Assert the exact diagnostic text.

Run only `bench worktree exec ft174-ticket-grammar -- go test ./internal/conformance/...`.
Do not commit. Do not edit the spec.
