# Derive the binding matrix from the record

Blocked by: 01-add-the-harness-record-package.md
Writes: internal/lines/lines.go, internal/lines/lines_parse_test.go, internal/lines/lines_verdict_test.go

## What to build

`internal/lines` keeps no harness list of its own. `lines.Harnesses` becomes
the record's rows whose providers are not `none`, so it holds `codex`,
`claude`, and `opencode` in the record's order. The `none` row stays out,
because `none` binds no provider and takes no cell.

`lines.CellFault` reads the row's providers instead of a harness literal. A
row with providers `any` demands a provider-qualified model id, so
`CellFault("opencode", "gpt-5")` reports the qualification fault. A row that
binds one provider accepts a bare id, so `CellFault("codex", "gpt-5")`
reports no fault. No harness name stays as a literal inside the rule.

This ticket reads the exported contract ticket 01 states: `harnesses.Rows`
and the row's `Providers` value. It changes no other package, so it runs
beside tickets 03 and 04.

The existing `lines` unit tests extend rather than move. Keep the current
diagnostic wording, because `Harnesses` feeds the parse errors and the
foreign-key report.

## Acceptance

- [ ] `lines.Harnesses` equals `codex`, `claude`, `opencode`. (covers HC11)
- [ ] `lines.CellFault("opencode", "gpt-5")` reports the provider-qualified fault and `CellFault("codex", "gpt-5")` reports none. (covers HC12)
- [ ] A `lines.env` key `BENCH_NONE_MID=` reports as a foreign key.
- [ ] `internal/lines` holds no harness name as a literal.
