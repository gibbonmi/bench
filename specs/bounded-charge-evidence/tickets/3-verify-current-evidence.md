# Verify current evidence

Blocked by: 2-publish-bounded-build-evidence.md
Writes: internal/preflight, internal/chargeevidence (new), .agents/skills/bench-craft-delegate/references/charge-evidence-format.md (new), reviews/bounded-charge-evidence.md (new), internal/systemtest, internal/conformance/injected_ports_registry_test.go, cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: CE7, CE8, CE55, CE56, CE57, CE58, CE59, CE60, CE61, CE62, CE65, CE106, CE120, CE123, CE124, CE133, CE134, CE158, CE159, CE160, CE161, CE162, CE174

## What to build

Verify all bytes, bind current action, and navigate sources independently across released assignments.

Add the approved source selector, full verification, and current-action modes to the working evidence path.
An independent consumer reconstructs the manifest and checks actual membership and delivery coverage.
Check-current validates the current active assignment and required source state.
Origin assignment identifiers never become logical identity or surviving action authority.

Exercise separate processes with repeated and out-of-order reads.
Compare state before and after, excluding only named lock files and existing ordinary telemetry.
A disk-backed cursor or consumer reading log must make CE65 red.
Prove the same artifact remains readable after production Bench assignment release.

This ticket alone owns the assignment-release proof; later cleanup work does not repeat it.
Activate only the new implemented grammar in public help and generated response documentation.
Record every mode-exclusion and response-field mutation at the public command seam.
CE174 owns source-selection grammar, including valid cursor pairing and each refused combination.
Keep canonical build guidance on its legacy route until ticket 4 activates the complete action path.

Review chunk: CE-C1C

Start only after the preceding green checkpoint and its independent chunk review close.
Use the shared contract in the spec; do not duplicate its field inventories here.
The successor starts only after this chunk review and its repair coverage close.

## Context price and checks

8 focused files; at most 1,500 source lines. Read the reviewed store API, assignment checks, and the production system harness.
This slice owns 23 predicates, including separately named table cases.
Use one retained context and the existing shared fixture harness; do not copy private helpers across packages.
Do not add code to an over-budget file without moving its responsibility and headroom in this ticket.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./internal/chargeevidence`
- `bench test --package ./cmd/bench`
- `bench test --check system`

Record a biting omission or mutation for every independent expected schema or policy fact.
The completion-plan probe is the minimum named mutation, not a substitute for the acceptance cases.

## Acceptance

- [ ] CE7: The independent reader rejects a manifest whose hash differs from the trusted identity.
- [ ] CE8: The independent reader rejects a source page outside verified manifest membership.
- [ ] CE55: A sibling worktree reads the same prepared bytes.
- [ ] CE56: Production assignment release preserves historical evidence reads.
- [ ] CE57: Check-current refuses a moved source tip.
- [ ] CE58: Check-current refuses a released assignment binding.
- [ ] CE59: Check-current returns the current verified assignment binding for unchanged pins.
- [ ] CE60: Full verification refuses corruption in an unread page.
- [ ] CE61: Full verification refuses a source digest mismatch.
- [ ] CE62: An explicit source selector returns only that declared source stream.
- [ ] CE65: Repeated and out-of-order reads create or change no consumer-progress state.
- [ ] CE106: The independent delivery check rejects missing, duplicate, and final-only page sets.
- [ ] CE120: A nested working directory retrieves the same artifact and source bytes.
- [ ] CE123: Check-current refuses a dirty assignment checkout.
- [ ] CE124: Check-current refuses changed required-source bytes.
- [ ] CE133: The verify response contains at most 48,000 encoded stdout bytes.
- [ ] CE134: The check-current response contains at most 48,000 encoded stdout bytes.
- [ ] CE158: Full verification accepts only its exclusive declared argument form.
- [ ] CE159: Current checking accepts only its exclusive declared argument form.
- [ ] CE160: Full verification returns the complete registered verified response schema.
- [ ] CE161: Current checking returns the complete registered current response schema.
- [ ] CE162: The public help inventory projects the implemented source, verify, and current grammar.
- [ ] CE174: Source-selected reads accept exactly their declared argument forms.
