# Show the whole completion plan and its record state

Blocked by: 01-collect-record-findings.md
Writes: internal/reviewrecord/recordcmd/status.go (new), internal/reviewrecord/recordcmd/status_test.go (new), internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/repair.go, internal/reviewrecord/recordtest/status.go (new), cmd/bench/main.go, cmd/bench/axi_query_cases_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, .agents/skills/bench-craft-cli/SKILL.md, projects/benchkit.md
Covers: RP1, RP2, RP3, RP4, RP5, RP6, RP7, RP8, RP9, RP10, RP11, RP12, RP13, RP14, RP28

## What to build

Review chunk: RP-C1.

Add the read-only form `bench record status <slug> [--offset <byte-offset> --source <identity>]` to the `bench record` family. Declare the form in `status.go`, and add one entry to the `forms` list in `command.go`. Put its form-name constant beside the other form names in `internal/reviewrecord/repair.go`. The RR49 route-bypass scan reds a bare `"next"` literal or a `next=` needle in a production `recordcmd` file, so compose no route text by hand.

The form reads the record and the plan at HEAD of the target tree, and it writes no file and no ledger entry. It skips the primary-checkout refusal (`command.go:286-291`), and every other form keeps that refusal. It prints the tables `plan[1]`, `chunks[N]`, `tickets[N]`, `verification[N]`, `final_verification[N]`, and `assignments[N]` with the fields that the spec's output table fixes. Each `entry` and `state` cell comes from the RP-C1 collector. With no record file, `record` is `absent` and every chunk `entry` is `missing`, at exit 0. The form pages its projection through `responsebound.PageRead`, as `bench worktree list --spec` does (`internal/worktree/related.go:85-96`). Every read ends with the `help[N]{cmd,why}:` block of `axi.RenderHelp`.

Approve the child `status` on the `record` row of `cmd/bench/main.go:117`. Add it to `approvedAXIQueries`, to the AXI envelope cases in `cmd/bench/axi_query_cases_test.go`, to the approved-query table of the craft-cli skill, and to the AXI query list in `projects/benchkit.md`. `internal/conformance/axi_query_registry_test.go` is over its line budget, so keep its net line count unchanged. Add the per-file prose budget row `.agents/skills/bench-craft-cli/SKILL.md | 122` to the budget table in `projects/benchkit.md`, by the reviewer decision of 2026-11-08. Add the new help row to the exact inventory in `cmd/bench/help_inventory_test.go`. `internal/reviewrecord/recordtest/fixture.go` has 3 lines of headroom, so put any new fixture helper in the new `recordtest/status.go`.

## Acceptance

- [ ] The view prints one `chunks` row per planned chunk, one `tickets` row per planned ticket, and one `assignments` row per recorded assignment (RP1, RP2, RP3).
- [ ] A stale chunk entry shows `stale`, and an absent one shows `missing` (RP4, RP5).
- [ ] A missing chunk requirement shows `missing`, a stale one shows `stale`, and each final requirement shows its state (RP6, RP7, RP8).
- [ ] A long projection prints its first page and a cursor, and a stale cursor refuses with `source changed since the cursor` (RP9, RP10).
- [ ] With no record file, the view prints `absent` and `missing` at exit 0 (RP11).
- [ ] The view runs in the primary checkout, and `bench record chunk` there keeps its refusal (RP12, RP13).
- [ ] A status read leaves the record bytes, the spec bytes, and `git status` unchanged (RP14).
- [ ] `bench test --check axi-query-registry` passes with the `status` child approved (RP28).

## Verification

Run `bench test --package ./internal/reviewrecord/...` and `bench test --package ./cmd/bench` whole. Run `bench test --check axi-query-registry` and `bench test --check guidance-prose-budgets`. Observe each planned row red before its production edit, and record the route of each red. Self-probe: drop the `state` cells from the `verification` table, and report that RP6 and RP7 go red.
