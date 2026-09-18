# Match the spec fence to the ticket Writes union

Blocked by: 1-close-anchor-registry-writes.md
Writes: internal/preflight/fence_writes.go (new), internal/preflight/fence_writes_test.go (new), internal/preflight/decision.go, internal/preflight/gather.go, internal/preflight/proposal.go, internal/preflight/closure.go, internal/preflight/decision_test.go, internal/preflight/command_build_test.go, internal/preflight/proposal_test.go, internal/preflight/anchor_closure_test.go, internal/preflight/charge_evidence_test.go, internal/preflight/charge_test.go, internal/preflight/command_review_test.go, internal/preflight/source_tip_test.go, internal/preflight/command_bootstrap_test.go, internal/preflight/evidencecmd/evidence_budget_test.go, internal/preflight/preflighttest/fixture.go, internal/reviewrecord/recordtest/fixture.go, internal/systemtest/owner_landing_fixture_test.go, cmd/bench/preflight_version_test.go, internal/worktree/land_fixtures_test.go, internal/worktree/land_effects_test.go, internal/worktree/land_journey_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: SC14, SC15, SC16, SC17, SC18, SC19, SC20, SC21, SC22, SC23

## What to build

A spec author runs `bench preflight build <slug>` or `bench preflight review <slug>` and sees a `fence-writes` row after `anchor-closure`. The row compares the fence tokens with the union of owned paths across the parsed tickets. Each side drops a trailing `/` and the `(new)` marker. Each side removes the review pickup and every path that the spec folder or the capture folder authorizes.

When the sets differ, the row reds with `spec fence and ticket Writes: union differ: fence only: <sorted list>; Writes only: <sorted list>`, and it omits an empty side. The gatherer records the review pickup from the review record path owner. This ticket adds `fence-writes` to the proposal-tolerated checks list that ticket 1 created, so `--propose-writes` tolerates the red.

The fixture has two fence tables: the conformant fence and the review fence extras. Each seeded ticket's `Writes:` values derive from the fence table that its seed declares. Every derived entry carries the `(new)` marker, because the ticket builder has no tree access. A writes-aware ticket builder takes the entries for a test that adds fence lines.

Three seed readers move to that builder. The 5000-entry budget seed passes its generated entries. The closed-parenthesis bootstrap test passes `internal/real/` in a net-neutral edit, so `command_bootstrap_test.go` stays at 395 lines. The proposal tests stop rewriting the literal `Writes: specs` text.

By reviewer decision, the shared seeds that the review record fixture, the landing race, and the version test own also become union-exact. The landing race fixture is a system test. Its focused run is `bench test --check system`, which pins `BENCH_KIT` to the graded kit root.

The new row also changes the asserted rows. This ticket adds `fence-writes` to the legacy charge baseline, the review command's row count, and the two `Decide` row-order lists. It makes `baseFacts` union-exact. The edits in `command_review_test.go` and `decision_test.go` are net-neutral in line count.

## Acceptance

- [ ] `Decide` reds `fence-writes` when the fence holds `a/` and the tickets write `b.go`, with the exact two-sided detail.
- [ ] `fence-writes` is green for equal sets and for a pickup on one side only.
- [ ] `fence-writes` is green for `internal/x/` against `internal/x (new)`, and for a `Writes:` path under the spec folder or `capture`.
- [ ] `bench preflight build` reports a backticked prose token in the fence section as fence-only.
- [ ] A `fence-writes` red does not refuse `--propose-writes`.
- [ ] `bench preflight build` over the conformant seed renders `fence-writes,green`.
- [ ] The preflight and evidence command suites stay green, including the review evidence, budget, and bootstrap seeds.
- [ ] In build mode with no tickets directory, the row renders not-applicable after `anchor-closure`.
