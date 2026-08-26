# Derive the guard and conformance harness loops

Blocked by: 01-add-the-harness-record-package.md
Writes: internal/guards/guards.go, internal/guards/guards_test.go, internal/conformance/line_routing_exec_test.go

## What to build

Two more readers stop naming their own harness list.

`guards.wiredHarnesses` reads each record row's hook config path instead of
the two paths it names today. A row with no hook config contributes nothing.
`bench guards` therefore reports `claude,codex` for a script both configs
name, and `none` for a script neither config names. A new harness with a hook
config joins the report as one record row.

The conformance `harnessOf` map and the literal harness loops beside it read
the record. `checkAdapterLineGuards` walks the record's rows that name a
headless adapter, so a missing `.bench/adapters/opencode` names that adapter
in a diagnostic. `checkLineHarnessSurfaces` walks the same source.

This ticket reads the exported contract ticket 01 states: `harnesses.Rows`,
the row's `HookConfig`, and the row's `Headless`. It writes no file that
tickets 02 and 03 write, so the three run in parallel. Keep every existing
diagnostic wording, because the record's row order matches the old literal
order.

## Acceptance

- [ ] `bench guards` reports `claude,codex` for a script both configs name and `none` for a script neither names. (covers HC14)
- [ ] `checkAdapterLineGuards` on a root missing `.bench/adapters/opencode` names that adapter. (covers HC15)
- [ ] `internal/guards` names no hook config path as a literal.
- [ ] The conformance harness loops declare no harness name as a literal.
