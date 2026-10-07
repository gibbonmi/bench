# Own bounded lifetime through the bounds runner

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/bounds/bounds.go, internal/bounds/bounds_test.go, internal/bounds/lifetime_test.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/subprocess, projects/benchkit.md, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/package-core-guard/bounds-classify-limit-restated, tests/canary/package-core-guard/bounds-discovery-window-unwrapped, tests/canary/package-core-guard/bounds-dot-import-package-alias, tests/canary/package-core-guard/bounds-dot-import-wait, tests/canary/package-core-guard/bounds-duplicate-owner, tests/canary/package-core-guard/bounds-intent-window-fixed, tests/canary/package-core-guard/bounds-multiple-dot-import-wait, tests/canary/package-core-guard/bounds-parenthesized-wait, tests/canary/package-core-guard/bounds-raw-elapsed-wait, tests/canary/package-core-guard/bounds-raw-injected-wait, tests/canary/package-core-guard/bounds-raw-wait-deadline, tests/canary/package-core-guard/bounds-raw-wait-duration, tests/canary/package-core-guard/bounds-read-limit-restated, tests/canary/package-core-guard/bounds-reassigned-wait-duration, tests/canary/package-core-guard/bounds-redeclared-wait-duration, tests/canary/package-core-guard/bounds-worktree-window-unwrapped, tests/canary/skill-description-budgets/budget-table-missing, tests/canary/skill-description-budgets/description-folded, tests/canary/skill-description-budgets/description-missing, tests/canary/skill-description-budgets/over-budget-command, tests/canary/skill-description-budgets/over-budget-description, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route
Covers: PL2, PL3, PL4, PL5, PL6, PL7, PL8, PL9, PL10, PL11, PL12, PL13, PL14, PL15, PL16, PL17, PL18, PL19, PL28, PL66, PL67, PL68, PL69, PL124, PL125

## What to build

A canceled bounds run returns its current classification after bounded teardown of a resistant local child.
Implement subprocess Lifetime and migrate the real bounds runner through it.
Keep the other ten caller implementations operational until their migration tickets.
This checkpoint owns no global caller census assertion.

Expose NewLifetime, Lifetime.Start, Lifetime.Cancel, and Lifetime.Wait to later callers.
Start queues exactly one serialized startup worker; callers do not infer successful launch from request acceptance.
Wait returns immutable Outcome and an eventual Completion handle.
Outcome separates process, cancellation, stream, cleanup, and protection facts from the accepted spec.
Completion publishes a new snapshot after unfinished work completes.

Policy takes caller-selected normal mode, signal, grace, final window, and normal-output delay.
Keep command platform attributes and caller files intact.
Owner-created file pipes separate process wait from independent consumers.
Bounds remains the duration and limit owner; subprocess imports neither bounds nor resource packages.

Use immediate KILL for bounds cancellation and retain its current status projection.
Bounds supplies the accepted two-second missing grace and three-second final window to later catchable callers.
Use FixedWindow for teardown under the unbounded-verdict switch.
Leave a registration adapter point for ticket 2 without publishing invented resources.

Correct the profile owner claim when this first shared owner appears.
Keep its other accepted guidance and anchor policy intact.
Register each new fault port with its real producer at this checkpoint.

Before implementation starts, confirm complete accepted and landed durable-file-replacement delivery.
Its leaf, review-record caller, and native qualification must be complete.
Durable-caller-migration is not a prerequisite.
Use one retained integration source and a fresh author for this ticket.

Keep each caller's selected normal-exit policy.
A raw child status cannot certify required cleanup or publication.
Retain uncertain resources and every unresolved obligation.
Never signal from recovered records or add a blanket normal-exit kill.

This ticket belongs to PL-C1.
Finish the preceding chunk review before this chunk starts.
Its blockers include every prior overlapping Writes owner.
No successor caller must exist for this ticket's owned predicates.

## Acceptance

- [ ] PL2: A pre-canceled lifetime starts no child.
- [ ] PL3: Cancellation during Start reaches the eventual resistant child.
- [ ] PL4: A failed Start creates no waiter.
- [ ] PL5: Concurrent cancellation and completion create exactly one waiter.
- [ ] PL6: Exit seven retains raw status seven and complete output.
- [ ] PL7: Cleanup failure retains an already observed child exit.
- [ ] PL8: Captured output remains unread until its worker completes.
- [ ] PL9: A direct caller file remains open after completion.
- [ ] PL10: Forced owned-pipe closure returns incomplete output.
- [ ] PL11: A blocked caller writer returns a pending stream snapshot.
- [ ] PL12: A TERM-resistant group receives KILL after its grace.
- [ ] PL13: An unresolved waiter returns incomplete by the final deadline.
- [ ] PL14: A descriptor-retaining descendant cannot hold canceled decoding forever.
- [ ] PL15: Current ESRCH marks the group absent.
- [ ] PL16: An answering probe marks the group present.
- [ ] PL17: EPERM retains uncertain group state.
- [ ] PL18: An unknown probe error retains uncertainty.
- [ ] PL19: Preserve mode leaves a normal surviving group unsignaled.
- [ ] PL28: Bounds preserves immediate-KILL cancellation classification.
- [ ] PL66: Wait and stream obligations share one final cleanup window.
- [ ] PL67: The unbounded-verdict switch keeps teardown finite.
- [ ] PL68: An incomplete worker holds only required resources.
- [ ] PL69: Blocked caller stdin yields a pending local worker.
- [ ] PL124: Pending Start returns incomplete within the cancellation deadline.
- [ ] PL125: Late successful Start receives shutdown and one waiter.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/subprocess`
- `bench test --package ./internal/bounds`

Mutation witness: Omit remembered cancellation for a held Start, start a second waiter, or wait unconditionally on a blocked consumer. The corresponding real-child or controlled-start assertion must fail.
