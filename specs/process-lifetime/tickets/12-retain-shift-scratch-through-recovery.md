# Retain shift scratch through startup cancellation and recovery

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md, 08-preserve-release-step-evidence.md, 09-bound-preflight-and-scanner-shutdown.md, 10-preserve-focused-test-completion.md, 11-preserve-worktree-exec-and-shell.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/shift/lifetime_test.go (new), internal/shift/loop.go, internal/shift/recover.go, internal/shift/session.go, internal/systemtest/process_lifetime_shift_test.go (new), internal/worktree/lifecycle.go, internal/worktree/pool_reclaim.go, internal/worktree/protection.go (new), internal/worktree/protection_test.go (new), internal/worktree/resume.go
Covers: PL83, PL84, PL85, PL86, PL87, PL88, PL117

## What to build

Actual shift startup cancellation reaches a resistant adapter while protected scratch survives incomplete rollback or recovery.
Store the established lifetime handle before Start rather than racing cmd.Process publication.
Signal and deadline paths call Cancel on that same handle.
Use the current signal, accepted missing two-second grace, and shared final cleanup window.

Consume the worktree Descriptor and protected-user predicate from ticket 11.
Refuse rollback, scratch removal, release, reuse, and reclaim while users remain unresolved.
Recovery retains protected scratch even when its ordinary dirty set is empty.
Recheck protection after cleanup planning and before apply under the existing exclusion.

Preserve deadline precedence and the current evidence taxonomy.
Pause the production startup transition in the test, not a modeled command decision.
Run the actual adapter and recovery owners in the second-process witness.

Keep over-budget resume and shift sources from gaining lines.
Move only cohesive protection or shutdown composition into this ticket's co-owned protection, session, or recovery source.

Before implementation starts, confirm complete accepted and landed durable-file-replacement delivery.
Its leaf, review-record caller, and native qualification must be complete.
Durable-caller-migration is not a prerequisite.
Use one retained integration source and a fresh author for this ticket.

Keep each caller's selected normal-exit policy.
A raw child status cannot certify required cleanup or publication.
Retain uncertain resources and every unresolved obligation.
Never signal from recovered records or add a blanket normal-exit kill.

This ticket belongs to PL-C5.
Finish the preceding chunk review before this chunk starts.
Its blockers include every prior overlapping Writes owner.
No successor caller must exist for this ticket's owned predicates.

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.
Do not introduce another subprocess fixture runner or private binary publisher.

## Acceptance

- [ ] PL83: Actual shift startup cancellation reaches a resistant adapter.
- [ ] PL84: Shift deadline keeps precedence over simultaneous interruption.
- [ ] PL85: Recovery retains protected scratch in a clean worktree.
- [ ] PL86: Explicit cleanup rechecks protection before apply.
- [ ] PL87: Pool reuse refuses a protected checkout.
- [ ] PL88: Pool reclaim retains protected member users.
- [ ] PL117: Shift refuses rollback while resource users remain unresolved.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/shift`
- `bench test --package ./internal/worktree`
- `bench test --check system`

Mutation witness: Read cmd.Process before the serialized transition, let an empty dirty set authorize removal, or trust a stale cleanup fingerprint. The actual resistant adapter and protected scratch witnesses must fail.
