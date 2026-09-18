# Publish bounded build evidence

Blocked by: 1-validate-legacy-prepared-packs.md
Writes: internal/preflight, internal/chargeevidence (new), .agents/skills/bench-craft-delegate/references/charge-evidence-format.md (new), reviews/bounded-charge-evidence.md (new), internal/systemtest, internal/conformance/injected_ports_registry_test.go, internal/git/worktree_admin.go, cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/package_shipped_surface_test.go, internal/conformance/help_inventory_single_source_test.go (new), tests/canary/package-core-guard/unrouted-subcommand
Covers: CE1, CE2, CE3, CE4, CE5, CE6, CE9, CE10, CE11, CE12, CE13, CE14, CE15, CE17, CE39, CE40, CE53, CE54, CE63, CE64, CE74, CE75, CE76, CE77, CE86, CE87, CE88, CE89, CE90, CE91, CE92, CE93, CE94, CE95, CE96, CE97, CE98, CE99, CE118, CE128, CE129, CE130, CE131, CE132, CE138, CE139, CE150, CE151, CE152, CE153, CE154, CE155, CE156, CE157

## What to build

Prepare persistent build evidence and traverse its manifest and bodies through bounded responses.

Use the reviewed pack interface to implement build preparation and default bounded traversal.
This checkpoint adds the final build preparation form and default evidence read form.
Default successors use only the implemented cursor form.
Explicit source selection, verify, and check-current remain unavailable until ticket 3.
The existing build --full route remains for the current canonical phase until ticket 4.

Add repo-common storage, safe temporary publication, quota accounting, and process-safe writer serialization.
Hold the writer lock from capacity calculation through verification and publication.
Do not add a persistent reservation ledger or consumer progress file.
Every default read is already stateless; repeated reads cannot create a reading session.

Drive prepare, manifest traversal, and every required source through the public command.
Use the reviewed format module for corruption and deterministic filesystem failures.
Rerun inherited source-refusal cases and require no published artifact from any refusal.
Expose the command registry projection through cmd/bench/main.go and test the exact root help inventory.
Advertise only this checkpoint's implemented operations.

Generate prepared and page schema documentation from the canonical registry.
Record each schema omission, preparation-next mutation, and quota recovery mutation.
Capacity recovery gives the exact larger-quota option until cleanup exists.
CE157 records this checkpoint; ticket 6 migrates the same oracle under CE175 when cleanup becomes available.
Build guidance keeps the legacy path and its existing action requirements at this checkpoint.

Review chunk: CE-C1B

Start only after the preceding green checkpoint and its independent chunk review close.
Use the shared contract in the spec; do not duplicate its field inventories here.
The successor starts only after this chunk review and its repair coverage close.

## Context price and checks

10 focused files; at most 2,000 source lines. Start from the reviewed pack API, preflight command, Git common-directory owner, and filesystem precedent.
This slice owns 54 predicates, including separately named table cases.
Use one retained context and the existing shared fixture harness; do not copy private helpers across packages.
Do not add code to an over-budget file without moving its responsibility and headroom in this ticket.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./internal/chargeevidence`
- `bench test --package ./cmd/bench`
- `bench test --check system`

Record a biting omission or mutation for every independent expected schema or policy fact.
The completion-plan probe is the minimum named mutation, not a substitute for the acceptance cases.

## Acceptance

- [ ] CE1: A prepared build reconstructs the selected ticket and every required canonical source exactly.
- [ ] CE2: A ticket larger than 48,000 bytes remains fully retrievable.
- [ ] CE3: Preparation reports delivery as unverified.
- [ ] CE4: A page reports only its response completeness and stream state.
- [ ] CE5: Default traversal returns the complete manifest before metadata or bodies.
- [ ] CE6: A manifest scalar larger than the response limit reconstructs exactly.
- [ ] CE9: A source without a final newline reconstructs without an added newline.
- [ ] CE10: Unicode and escapable control bytes reconstruct exactly.
- [ ] CE11: Numeric-looking identifiers decode as strings.
- [ ] CE12: Representable path text remains data throughout retrieval.
- [ ] CE13: Build preparation contains at most 48,000 encoded stdout bytes.
- [ ] CE14: Oversized required metadata reconstructs without omission.
- [ ] CE15: An oversized invalid operand produces a bounded usage refusal.
- [ ] CE17: Sibling assignment and checkout locations do not change logical identity.
- [ ] CE39: A selected-page read refuses a changed page body.
- [ ] CE40: A selected-page read refuses a changed manifest.
- [ ] CE53: One movement discards the first attempt before the retry publishes.
- [ ] CE54: A second movement publishes no artifact from either attempt.
- [ ] CE63: The reader refuses malformed or foreign cursors.
- [ ] CE64: Every nonterminal read supplies the exact successor command.
- [ ] CE74: The default quota includes published artifacts and temporary writes.
- [ ] CE75: An explicit larger quota admits an artifact above the default quota.
- [ ] CE76: Capacity refusal preserves every existing artifact.
- [ ] CE77: Concurrent preparations cannot exceed the selected store quota.
- [ ] CE86: Interruption before source verification exposes no published handle.
- [ ] CE87: Interruption during artifact verification exposes no published handle.
- [ ] CE88: Publication failure at the terminal step exposes no partial artifact.
- [ ] CE89: Identical concurrent writers reuse one fully verified published artifact.
- [ ] CE90: Distinct concurrent writers retain their distinct artifacts.
- [ ] CE91: The artifact reader refuses a symlink object without consuming it.
- [ ] CE92: The artifact reader refuses a FIFO object without consuming it.
- [ ] CE93: The artifact reader refuses a socket object without consuming it.
- [ ] CE94: The artifact reader refuses a device object without consuming it.
- [ ] CE95: The artifact reader refuses a directory object without consuming it.
- [ ] CE96: An invalid artifact identifier refuses before path resolution.
- [ ] CE97: A corrupt existing artifact refuses replacement during preparation.
- [ ] CE98: An unsafe store directory refuses preparation and retrieval.
- [ ] CE99: An absent store creates only during explicit preparation.
- [ ] CE118: An absent store read refuses without creating the store.
- [ ] CE128: A missing required tool preserves a bounded preparation refusal.
- [ ] CE129: A replaced open artifact refuses a changed selected page.
- [ ] CE130: Invalid quota operands refuse before temporary storage changes.
- [ ] CE131: The manifest read response contains at most 48,000 encoded stdout bytes.
- [ ] CE132: The source read response contains at most 48,000 encoded stdout bytes.
- [ ] CE138: The usage refusal response contains at most 48,000 encoded stdout bytes.
- [ ] CE139: The operational refusal response contains at most 48,000 encoded stdout bytes.
- [ ] CE150: Build preparation accepts exactly its declared argument forms.
- [ ] CE151: Default evidence reads accept exactly their declared argument forms.
- [ ] CE152: The public help inventory projects exactly the implemented preparation and read grammar.
- [ ] CE153: Build preparation returns the complete registered prepared response schema.
- [ ] CE154: Manifest fragments return the complete registered page response schema.
- [ ] CE155: Source fragments return the complete registered page response schema.
- [ ] CE156: Build preparation returns the exact manifest-first retrieval command.
- [ ] CE157: Before cleanup lands, capacity refusal returns the exact larger-quota retry instruction.
