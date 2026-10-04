# Authorize publication against current commitment

Blocked by: 04-bind-worktrees-and-shifts.md
Writes: internal/commitment (new), internal/intent, internal/worktree, internal/landing, internal/roadmap, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/systemtest/otel_verbs_test.go, internal/systemtest/adoption_test.go, internal/systemtest/owner_landing_fixture_test.go
Covers: DC12, DC14, DC28, DC29, DC30, DC32, DC49, DC72

## What to build

Use the installed publication broker to admit the exact candidate against current default-branch authority and local receipts.
Ticket 04 supplies assignment and continuation identity. Ticket 03 supplies the shared planning classifier.
Planning can commit and land a staged spec before adoption, but cannot publish production files or rewrite a protected obligation.

Recheck authority after the prospective gate. Hold the existing intent transaction across the final admission decision and ref publication.
Do not hold that lock throughout the gate. A policy replacement or competing blocker must not win between the final decision and publication.
Preserve current compare-and-swap and publication recovery behavior.

Honor only explicitly listed legacy continuations at adoption, within their approved scope. Do not treat all pre-existing assignments as authorized.
The installed broker remains the independent publication owner. No new wrapper authentication system or startup mechanism is required.

Read worktree land admission and resume, landing LandReviewed and ref publication, and the intent transaction. Reuse the current landing fixture owner. The new publication decision is reviewed before closure extends its candidate transform.

## Acceptance

- [ ] Planning lands a map and staged spec before adoption while refusing production files (DC30, DC32).
- [ ] Protected-row rename and sequence reset refuse publication (DC12, DC14).
- [ ] A displaced assignment and a policy change during the gate cannot publish (DC28, DC29).
- [ ] A listed legacy run can finish only its named existing scope (DC49).
- [ ] A coordinated blocker race cannot alter admission between the final check and ref update (DC72).

## Checkpoint verification

The existing system journey uses BENCH_KIT through `bench test --check system`. Adapt its fixture at this checkpoint when admission changes its route.

Run `bench test --package ./internal/worktree`, `bench test --package ./internal/landing`, and `bench test --package ./internal/commitment`. Observe the gate-time policy race through the real broker orchestration seam.
