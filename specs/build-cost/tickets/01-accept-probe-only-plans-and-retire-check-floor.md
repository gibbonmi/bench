# Accept probe-only plans and retire the check-floor row

Blocked by: none
Writes: internal/reviewrecord/plan.go, internal/reviewrecord/ticket.go, internal/reviewrecord/plan_probe_test.go (new), internal/reviewrecord/recordcmd/ticket_test.go, internal/preflight/check_floor.go, internal/preflight/check_floor_test.go, internal/preflight/check_floor_command_test.go, internal/preflight/decision.go, internal/preflight/decision_test.go, internal/preflight/plan.go, internal/preflight/charge_cells_test.go, internal/testreport/named_check.go, internal/conformance/ticket_grammar_test.go, internal/reviewrecord/recordcmd/ticket.go, internal/reviewrecord/delegated_test.go, internal/preflight/source_tip_test.go, internal/reviewrecord/recordcmd/ticket_probe_test.go (new), internal/testreport/floor.go, internal/testreport/floor_test.go, internal/testreport/selection_facts.go, internal/testreport/selection_facts_test.go
Covers: BC24, BC25, BC26, BC27, BC28, BC29, BC30, BC31, BC32

## What to build

The plan grade accepts a chunk whose `verification` list is empty or absent. A version 2 plan no longer needs one requirement for each ticket. Each listed chunk requirement of a version 2 plan still names a ticket of its own chunk.
The final verification list must still be non-empty. In `requirementsValid`, the empty check reads `if len(items) == 0 && owned == nil {`, because the probe-01 command of the plan swaps that exact line.

`TicketRequirements` answers an empty list for a ticket that owes nothing, also when the chunk lists no requirement. `bench record ticket` for such a ticket writes the chunk step and no verification entry, and it exits 0.
The chunk checkpoint keeps its evidence pass and its oracle selection. A planned probe requirement without an entry still refuses with the verification form named.

Delete the `check-floor` row: `check_floor.go`, its two test files, its row in `internal/preflight/decision.go`, and its facts in `internal/preflight/plan.go`. Remove the row from each test that lists the rows, and correct each comment that names the row.
The charge's root conformance anchor in `internal/anchors/registry_check_floor.go` stays.

In the same commit, delete every plain chunk requirement from the completion plan of `specs/build-cost/spec.md`, as the spec's self-hosting order states. Keep each probe requirement and the final list.
Record this ticket with the current binary, and commit the record with plain `bench commit`. Run the post-ticket preflight as `bench preflight --in <target> build build-cost --base <chunk-base> --source-tip <tip>` after `bench worktree build <target>`.

## Acceptance

- [ ] A chunk with an empty `verification` list passes the plan grade (BC24).
- [ ] A chunk without a `verification` key passes the plan grade (BC25).
- [ ] A version 2 plan in which one ticket owes no requirement passes, and a requirement that names a ticket outside its chunk still refuses (BC26, BC27).
- [ ] A plain final requirement still passes (BC28).
- [ ] `bench record ticket` for a ticket that owes nothing exits 0, freezes the chunk, and prints an empty `verifications` cell (BC29).
- [ ] A chunk checkpoint whose probe requirement has no entry refuses with the verification form named (BC30).
- [ ] `bench preflight build` prints no `check-floor` row (BC32).
- [ ] The completion plan of `specs/build-cost/spec.md` holds no chunk requirement without `probe`.
