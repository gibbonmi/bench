# 3. Declare the discarded namespace and sweep it at 30 days

Blocked by: none
Writes: internal/intent/ledger/ledger.go, internal/intent/ledger_aliases.go, internal/worktree/reconcile.go, internal/worktree/reconcile_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RI42, RI43, RI44, RI45, RI46, RI47

## What to build

Chunk: RI-C2a.

Declare `refs/bench/discarded/` as one ledger constant beside the recovery and reset namespaces, with an alias in the intent package.
Add one function that renders a discarded ref path from a UTC instant and a branch ref.
The path is the namespace, the date as `yyyymmdd`, a slash, and the branch path without `refs/heads/`.
The namespace joins neither the lifecycle emptying list nor the reset rule.

Extend the session-start sweep with the discarded rule.
The sweep lists the namespace, parses the first path segment under it as a UTC date, and keeps a ref whose segment does not parse.
It deletes a ref at its listed object when the resume instant is 30 days or more past that date.
The swept count joins the existing total.
A delete that fails because the ref moved reports an error, as the existing lifecycle deletes do.

## Acceptance

- [ ] The path function renders `refs/bench/discarded/20260927/bench/assign/<owner>/<id>` for 2026-09-27 and that branch.
- [ ] A planted ref dated 30 days before the instant is deleted and the swept count is 1.
- [ ] A planted ref dated 29 days before the instant survives and the swept count is 0.
- [ ] A planted ref dated today survives the pass that deletes a planted recovery ref.
- [ ] A planted ref whose date segment is `latest` survives and the swept count is 0.
- [ ] A ref moved between the listing and the delete stays, and the sweep returns an error.
