# Preserve worktree command policy with durable protection

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md, 08-preserve-release-step-evidence.md, 09-bound-preflight-and-scanner-shutdown.md, 10-preserve-focused-test-completion.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/process_resources.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/exec.go, internal/worktree/exec_lifetime_test.go (new), internal/worktree/exec_test.go, internal/worktree/lifecycle.go, internal/worktree/pool_reclaim.go, internal/worktree/protection.go (new), internal/worktree/protection_test.go (new), internal/worktree/resume.go, internal/worktree/subshell.go, internal/worktree/subshell_lifetime_test.go (new), internal/worktree/subshell_test.go
Covers: PL81, PL82, PL114, PL115

## What to build

A real worktree exec keeps its normal output delay without killing a normally surviving group.
Migrate exec and shell through Lifetime while retaining their distinct policies.
Exec keeps INT/130 and its current three-second grace.
Shell forwards the actual signal, keeps its current five-second grace, and retains the interrupted lease.

Store checkout protection in its Git administrative directory outside reset and clean targets.
Expose the worktree Descriptor and install its existing assignment-identity validator in the CLI table.
Keep the reserved inherited/local vector through current environment routing; --env cannot replace it.
Preserve command inputs, caller streams, assignment behavior, and the separate ExecWaitDelay.

Exec can return child zero only after its own required stream and cleanup obligations complete.
It releases no assignment on that return; durable users still prevent reuse and removal.
Required enclosing owners consume unresolved disposition across that CLI boundary.

Introduce the single protected-user predicate into release, cleanup planning/apply, reuse, and reclaim composition.
Fingerprint protection evidence and recheck under exclusion before apply.
Do not add another pool cleanup rule.
Ticket 12 completes and witnesses shift-specific rollback and recovery consumption.
PL114 and PL115 stay with this actual environment-routing consumer.

Keep over-budget resume from gaining lines.
Move only cohesive protection composition into this ticket's co-owned worktree protection source.

Before implementation starts, confirm complete accepted and landed durable-file-replacement delivery.
Its leaf, review-record caller, and native qualification must be complete.
Durable-caller-migration is not a prerequisite.
Use one retained integration source and a fresh author for this ticket.

Keep each caller's selected normal-exit policy.
A raw child status cannot certify required cleanup or publication.
Retain uncertain resources and every unresolved obligation.
Never signal from recovered records or add a blanket normal-exit kill.

This ticket belongs to PL-C4.
Finish the preceding chunk review before this chunk starts.
Its blockers include every prior overlapping Writes owner.
No successor caller must exist for this ticket's owned predicates.

## Acceptance

- [ ] PL81: Exec normal-output delay adds no normal group kill.
- [ ] PL82: Shell interruption retains its lease after complete shutdown.
- [ ] PL114: Worktree environment transforms retain the resource vector.
- [ ] PL115: Worktree --env cannot override the resource vector.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/worktree`

Mutation witness: Add a normal-exit group kill, strip the descriptor vector, or release the interrupted shell lease. The real command, environment, and lease assertions must fail.
