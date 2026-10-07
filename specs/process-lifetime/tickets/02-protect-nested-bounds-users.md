# Protect cache users through nested bounds launches

Blocked by: 01-own-lifetime-through-bounds.md
Writes: DATA_HANDLING.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/process_resources.go (new), internal/anchors/registry_commitment.go, internal/bounds/bounds.go, internal/bounds/lifetime_test.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gocache/clean.go, internal/gocache/lock.go, internal/gocache/protection.go (new), internal/gocache/protection_test.go (new), internal/subprocess, internal/systemtest/process_lifetime_test.go (new), tests/canary/data-handling-derivation/undocumented-passlist-var, tests/canary/package-core-guard/bounds-classify-limit-restated, tests/canary/package-core-guard/bounds-discovery-window-unwrapped, tests/canary/package-core-guard/bounds-dot-import-package-alias, tests/canary/package-core-guard/bounds-dot-import-wait, tests/canary/package-core-guard/bounds-duplicate-owner, tests/canary/package-core-guard/bounds-intent-window-fixed, tests/canary/package-core-guard/bounds-multiple-dot-import-wait, tests/canary/package-core-guard/bounds-parenthesized-wait, tests/canary/package-core-guard/bounds-raw-elapsed-wait, tests/canary/package-core-guard/bounds-raw-injected-wait, tests/canary/package-core-guard/bounds-raw-wait-deadline, tests/canary/package-core-guard/bounds-raw-wait-duration, tests/canary/package-core-guard/bounds-read-limit-restated, tests/canary/package-core-guard/bounds-reassigned-wait-duration, tests/canary/package-core-guard/bounds-redeclared-wait-duration, tests/canary/package-core-guard/bounds-worktree-window-unwrapped
Covers: PL20, PL21, PL22, PL23, PL24, PL25, PL26, PL27, PL39, PL40, PL41, PL42, PL43, PL44, PL45, PL46, PL47, PL48, PL49, PL50, PL51, PL52, PL53, PL54, PL70, PL71, PL72, PL73, PL74, PL104, PL105, PL106, PL107, PL108, PL109, PL110, PL111, PL112, PL113, PL120

## What to build

A real nested Bench provider probe registers its separate group before it uses an inherited cache resource.
Use the existing bench models route through bounds.Run and a local fake codex provider.
Keep network credentials absent and discovery online so that the actual bounds producer runs.
Do not add a public command or change models policy.

Consume ticket 1's Lifetime control and immutable outcome.
Implement OpenResource, Prepare, Started, Finish, SealAndObserve, and Recover in subprocess/protection.
Use the accepted resource and user schema, exact status values, and complete inherited/local union.
Expose Descriptor, Registration, and Disposition to subsequent resource owners.
Safe, protected, and uncertain disposition retain the exact identifiers needed by their caller.

Prepare persists pending users before Start, in canonical store-lock order after owner cleanup exclusions.
Started publishes actual identity; Finish publishes the separate immutable completion facts.
Sealing refuses new users while permitting updates for already admitted pending users.
Partial persistence prevents launch or retains uncertainty as specified.
Use landed durablefile.Replace and its classified error; implement no second generic replacement algorithm.

Install the CLI composition-root validator table with the cache identity adapter.
Later resource tickets extend that same table before publishing descriptors of their kinds.
Unknown or unavailable kind validators refuse admission; they never authorize arbitrary paths.
Hold the cache shared lock before protection and observe protection before exclusive clean.
Keep same-process descriptors shared under the existing record-lock owner.

Publish the reserved BENCH_PROCESS_RESOURCES vector through the real bounds launch adapter.
The nested CLI installs the same available validators and registers its provider group.
Prove two nested users, partial records, sealed pending admission, and retained disposition after CLI exit.
This checkpoint does not depend on later gate, focused-test, worktree, or shift migrations.

Before implementation starts, confirm complete accepted and landed durable-file-replacement delivery.
Its leaf, review-record caller, and native qualification must be complete.
Durable-caller-migration is not a prerequisite.
Use one retained integration source and a fresh author for this ticket.

Keep each caller's selected normal-exit policy.
A raw child status cannot certify required cleanup or publication.
Retain uncertain resources and every unresolved obligation.
Never signal from recovered records or add a blanket normal-exit kill.

This ticket belongs to PL-C2.
Finish the preceding chunk review before this chunk starts.
Its blockers include every prior overlapping Writes owner.
No successor caller must exist for this ticket's owned predicates.

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.
Do not introduce another subprocess fixture runner or private binary publisher.

## Acceptance

- [ ] PL20: Protection persists before the child's resource-use marker.
- [ ] PL21: First user-publication failure prevents Start.
- [ ] PL22: Later resource-publication failure prevents Start.
- [ ] PL23: Failed rollback retains its partial pending record.
- [ ] PL24: Partial known-identity publication retains the launch disposition.
- [ ] PL25: Actor death retains a pending startup identity.
- [ ] PL26: Two concurrent users require two complete group-absence proofs.
- [ ] PL27: Sealed admission refuses new registration before Start.
- [ ] PL39: A foreign binding prevents registration.
- [ ] PL40: Malformed protection retains the resource.
- [ ] PL41: A protection symlink prevents admission.
- [ ] PL42: Trailing JSON prevents admission.
- [ ] PL43: Unknown record fields prevent admission.
- [ ] PL44: Wrong record mode retains the resource.
- [ ] PL45: Duplicate user identity makes the set uncertain.
- [ ] PL46: An admitted pending user blocks removal after sealing.
- [ ] PL47: An admitted user can finish after sealing.
- [ ] PL48: Reversed resource vectors acquire locks without deadlock.
- [ ] PL49: A nested CLI registers its separate group before use.
- [ ] PL50: Overlapping nested users retain independent records.
- [ ] PL51: Missing nested completion prevents outer safe release.
- [ ] PL52: Recovery requires current absence for every known group.
- [ ] PL53: Recovery sends no signal to recorded identities.
- [ ] PL54: Record age grants no removal authority.
- [ ] PL70: A proved no-start removes its protection durably.
- [ ] PL71: Resource-record sync failure prevents first launch.
- [ ] PL72: Deletion-sync failure returns uncertain cleanup.
- [ ] PL73: Actor absence resolves workers only after known groups are absent.
- [ ] PL74: A missing inherited resource record prevents Start.
- [ ] PL104: A FIFO protection record prevents admission without blocking.
- [ ] PL105: A directory at a record path prevents admission.
- [ ] PL106: Unsupported schema prevents admission.
- [ ] PL107: Duplicate JSON fields prevent admission.
- [ ] PL108: Oversized protection data prevents admission.
- [ ] PL109: Conflicting generations at one store prevent registration.
- [ ] PL110: Control-byte descriptor paths prevent launch.
- [ ] PL111: Space and glob paths retain exact resource identity.
- [ ] PL112: Repeated identical descriptors create one user per resource.
- [ ] PL113: A failed KILL remains unresolved until absence is proved.
- [ ] PL120: OpenResource preserves users when a second owner joins.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/subprocess/protection`
- `bench test --package ./internal/bounds`
- `bench test --package ./internal/gocache`
- `bench test --check system`

Mutation witness: Drop inherited descriptors from the real bounds launch. The nested-provider witness must fail through bench test --check system.
