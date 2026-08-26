# Add the harness record package

Blocked by: none
Writes: internal/harnesses/harnesses.go (new), internal/harnesses/harnesses_test.go (new)

## What to build

A new low package, `internal/harnesses`, owns one versioned record of every
harness Bench knows. The record file imports only the standard library, so
`lines`, `status`, `handoff`, `guards`, and `conformance` can all compose it
without an import cycle. Ticket 06 adds the verb beside the record in the
same package.

This exported contract crosses into tickets 02 through 09:

- `harnesses.Schema` is the integer schema version, and it equals 1.
- `harnesses.Rows` is the ordered slice `codex`, `claude`, `opencode`, `none`.
- `harnesses.Lookup(name string) (Row, bool)` returns one row by harness name.
- A `Row` holds `Harness`, `Providers`, `PhaseForm`, `HookConfig`, `HookEvents`, `DelegationGuard`, `Headless`, and `Mechanics`.
- A `Cell` holds `Value`, `Source`, and `Checked`.

`Providers` holds a closed value: one provider name, `any`, or `none`. The
`codex` row binds `openai`, the `claude` row binds `anthropic`, the
`opencode` row binds `any`, and the `none` row binds `none`. `PhaseForm` is
`/bench-` for `claude`, `$bench-` for `codex`, and empty for the other two.

`HookConfig` names `.claude/settings.json` for `claude` and `.codex/hooks.json`
for `codex`, and `HookEvents` lists the events each config wires. Read both
config files first and record the events the tree wires today. `opencode` and
`none` name no hook config. Only the `claude` row holds
`DelegationGuard: yes`. `Headless` names `.bench/adapters/<harness>` for the
three harness rows, and the `none` row names no adapter.

`Mechanics` holds every one of the twelve mechanics the spec lists. A cell
value is `yes`, `no`, or `unknown`. A `yes` or `no` cell carries a non-empty
source and an ISO date; an `unknown` cell carries neither. Fill a cell only
from a fact the tree records:

- the two hook configs
- the three headless adapters
- the Hook Layers verdict dated 2026-07-11
- the reference's effort rule

Every other cell starts as `unknown`.

The unit test walks every row and every cell. It grades the order, the
providers, the phase forms, the hook configs, the guard verdict, the adapter
paths, and the enum. A mutated copy of the record proves the walk bites: an
empty cell, a value outside the enum, and a future `checked` date each fail.

## Acceptance

- [ ] The record's rows are exactly `codex`, `claude`, `opencode`, and `none` in that order. (covers HC01)
- [ ] The `codex` row binds `openai`, `claude` binds `anthropic`, `opencode` binds `any`, and `none` binds `none`. (covers HC02)
- [ ] The `claude` phase form is `/bench-`, the `codex` form is `$bench-`, and the other two forms are empty. (covers HC03)
- [ ] The `claude` row names `.claude/settings.json` with six events, and the `codex` row names `.codex/hooks.json` with three. (covers HC04)
- [ ] Only the `claude` row holds `delegation_guard: yes`. (covers HC05)
- [ ] Each of the three harness rows names `.bench/adapters/<harness>`, and the `none` row names no adapter. (covers HC06)
- [ ] Every row holds every one of the twelve mechanics with a value inside the closed enum. (covers HC07)
- [ ] Every `yes` or `no` cell carries a non-empty source and a date that parses as ISO. (covers HC08)
- [ ] The `none` row has an empty phase form, no hook config, no adapter, and `headless execution: no`. (covers HC09)
- [ ] `harnesses.Schema` equals 1. (covers HC10)
- [ ] The record unit test fails on a row with an empty cell or a value outside the enum. (covers HC16)
- [ ] A cell whose `checked` date is later than today fails the record unit test.
- [ ] The record file imports only the standard library.
