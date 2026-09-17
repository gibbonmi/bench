# Clean evidence explicitly

Blocked by: 5-freeze-bounded-review-evidence.md
Writes: internal/preflight, internal/chargeevidence (new), internal/systemtest, internal/conformance/injected_ports_registry_test.go, reviews/bounded-charge-evidence.md (new), .agents/skills/bench-craft-delegate/references/charge-evidence-format.md (new), cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: CE78, CE79, CE80, CE81, CE82, CE83, CE84, CE85, CE100, CE117, CE119, CE121, CE122, CE136, CE137, CE168, CE169, CE170, CE171, CE172, CE175

## What to build

Deliver exact cleanup plans and safe apply behavior on supported native platforms.

Add the approved cleanup plan and apply forms to the existing store.
Return every exact target through bounded pages and bind the complete plan fingerprint.
Revalidate the fingerprint and target file identities before deletion.
Refuse active readers or writers without deleting any target.

Keep partial deletion and interruption outcomes honest and recoverable with a fresh plan.
Change capacity refusal recovery to the exact cleanup command now that its producer exists.
CE175 updates TestEvidenceCapacityRecovery; replace its earlier quota-first expectation instead of preserving an obsolete branch.
Retain an explicit larger-quota retry as the documented alternative.
Activate cleanup in public help and generate its response schema reference from the canonical registry.

Record cleanup grammar, exact-schema, inventory, stale-plan, and lock-exclusion mutations.
Run production storage checks on Linux, macOS, and Windows; record unavailable native evidence as pending.
Reconcile the complete chargeevidence package and every response-budget operation.

Review chunk: CE-C3

Start only after the preceding green checkpoint and its independent chunk review close.
Use the shared contract in the spec; do not duplicate its field inventories here.
The successor starts only after this chunk review and its repair coverage close.

## Context price and checks

8 focused files; at most 1,600 source lines. Read the store locks, cleanup parser, public inventory, and native filesystem test owner.
This slice owns 21 predicates, including separately named table cases.
Use one retained context and the existing shared fixture harness; do not copy private helpers across packages.
Do not add code to an over-budget file without moving its responsibility and headroom in this ticket.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./internal/chargeevidence`
- `bench test --package ./cmd/bench`
- `bench test --check system`

Record a biting omission or mutation for every independent expected schema or policy fact.
The completion-plan probe is the minimum named mutation, not a substitute for the acceptance cases.

## Acceptance

- [ ] CE78: Cleanup returns a bounded plan of exact deletion targets.
- [ ] CE79: Cleanup apply refuses a changed plan fingerprint before deletion.
- [ ] CE80: Cleanup without apply deletes nothing.
- [ ] CE81: An idle handle removed by explicit cleanup subsequently refuses retrieval.
- [ ] CE82: A cleanup deletion failure returns the exact unfinished disposition.
- [ ] CE83: Cleanup refuses while a reader holds the operation lock.
- [ ] CE84: Cleanup refuses while a writer holds the operation lock.
- [ ] CE85: Cleanup includes an orphan temporary file only after writer exclusion succeeds.
- [ ] CE100: An empty existing store reports an empty cleanup plan.
- [ ] CE117: Production storage checks pass on each supported native platform.
- [ ] CE119: An absent store returns an empty cleanup plan without creation.
- [ ] CE121: Cleanup refuses an unsafe target before deleting any target.
- [ ] CE122: Interrupted cleanup requires a fresh plan for the remaining targets.
- [ ] CE136: The cleanup plan response contains at most 48,000 encoded stdout bytes.
- [ ] CE137: The cleanup apply response contains at most 48,000 encoded stdout bytes.
- [ ] CE168: Cleanup accepts exactly its declared argument forms.
- [ ] CE169: Cleanup plans return the complete registered cleanup response schema.
- [ ] CE170: Cleanup target pages return the complete registered targets response schema.
- [ ] CE171: Cleanup apply returns the complete registered applied response schema.
- [ ] CE172: The public help inventory projects exactly the implemented cleanup grammar.
- [ ] CE175: After cleanup lands, capacity refusal returns `bench preflight evidence-clean` first, with the explicit larger-quota alternative.
