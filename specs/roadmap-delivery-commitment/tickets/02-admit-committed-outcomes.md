# Admit only eligible committed outcomes

Blocked by: 01-plan-exact-policy-changes.md
Writes: internal/commitment (new), internal/intent, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DC16, DC17, DC18, DC19, DC20, DC21, DC22, DC33, DC48, DC61, DC66

## What to build

Deliver `commitment start`, `block`, and `unblock`. A start atomically binds the current assignment, deliverable, and eligible outcome under the intent transaction.
Ticket 01 supplies validated policy facts and exact approval identity. This ticket adds the one shared admission decision and durable runtime claim.
The command refuses absent policy with the adoption remedy. It distinguishes corruption from missing adoption.

Apply the committed order, dependency edges, blocker state, and exact parallel grants. A blocker preserves the obligation and allows only the next independent committed outcome.
Unblocking A does not displace active B. Process exit or assignment release does not complete an outcome.
A second start can share the same outcome for ticket work, but a different outcome needs the exact approved grant.

This checkpoint proves admission through its own commands. Production consumer integration remains explicit in later tickets.
Do not add a process-liveness scheduler or infer authorization from an old sequence.

Read ticket 01 exports and the intent lock and record lifecycle. Keep this new admission seam in its own reviewed chunk before consumer work.

## Acceptance

- [ ] Uncommitted C and premature B refuse with no claim (DC16, DC17).
- [ ] Blocked A permits independent B, refuses dependent B, and never permits unrelated C (DC18 to DC20).
- [ ] A and B race for one default slot; exactly one succeeds (DC21, DC61).
- [ ] A grant naming A and B permits those outcomes and still refuses C (DC22).
- [ ] Restart retains the claim, and unblocking A leaves active B in place (DC33, DC66).
- [ ] Missing policy refuses start with a concrete initial-planning remedy (DC48).

## Checkpoint verification

Run `bench test --package ./internal/commitment` and `bench test --package ./internal/intent`. Use a controlled transaction race and verify both the reply and persisted claims. Observe a separate-check-and-write mutation turn the race check red.
