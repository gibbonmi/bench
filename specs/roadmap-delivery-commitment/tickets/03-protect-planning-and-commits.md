# Protect planning commits and build charges

Blocked by: 02-admit-committed-outcomes.md
Writes: internal/commitment (new), internal/intent, internal/commit, internal/preflight, internal/roadmap, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/systemtest/otel_verbs_test.go, internal/systemtest/charge_evidence_test.go
Covers: DC6, DC10, DC11, DC13, DC15, DC26, DC27, DC31

## What to build

Make the real commit command grade its complete candidate before formatting or lane effects. Add plan-only build preflight and preserve ordinary build admission checks.
Ticket 02 supplies the assignment admission decision. Ticket 01 supplies proposal identity and canonical policy facts.
Create one planning change classifier in the commitment owner. Both commit and the later publication consumer use that exported decision.

Permit scoped maps, staged specs, research, capture, and their required promotion documents. Append uncommitted roadmap rows without changing active obligations.
The roadmap owner separates occurrence evidence from requirement identity. Reject protected-row deletion, requirement replacement, stale approval bytes, and production changes under planning purpose.
Plan-only preflight checks the authored graph without creating a delivery claim. Ordinary preflight requires a bound deliverable.

Keep the existing invocation and identity precedence. The new commitment refusal must happen before Go formatting, lane execution, or a persisted build charge.
Update existing command fixture setup in these packages through their shared helpers.

The separate `newEvidenceJourney` constructor in `internal/systemtest/charge_evidence_test.go` also drives real build charges.
Seed an approved policy in that shared fixture and bind each assignment to its approved outcome before the charge.
Its callers in `charge_evidence_test.go` and `charge_cleanup_test.go` must retain their intended evidence assertions. Do not scatter independent permissive policy fixtures across callers or bypass admission for system tests.

Read commit orchestration, preflight gather and decisions, roadmap context parsing, and the new owner. Existing large files receive an extracted focused sibling file in the same package before added growth. This ticket owns that headroom move and its tests.

## Acceptance

- [ ] A changed proposal refuses commit despite an earlier receipt (DC6).
- [ ] Appending C and adding a valid occurrence to A preserve the current obligation and order (DC10, DC15).
- [ ] Deleting A or replacing its requirement refuses before writing a commit (DC11, DC13).
- [ ] An unbound production candidate refuses before a formatter marker appears (DC26).
- [ ] Ordinary preflight refuses an unbound deliverable; plan-only succeeds without creating a claim (DC27, DC31).
- [ ] The planning classifier has one production definition. Publication can consume that definition without copying its allowlist.

- [ ] The shared evidence journey produces successful real build charges for its admitted assignments. Location independence, concurrency, interruption, and cleanup callers keep testing their original behavior.

## Checkpoint verification

The existing system journey uses BENCH_KIT through `bench test --check system`. Adapt its fixture at this checkpoint when admission changes its route.

Run `bench test --package ./internal/commit`, `bench test --package ./internal/preflight`, `bench test --package ./internal/roadmap`, and `bench test --package ./internal/commitment`. Exercise commit through its actual composed-candidate path.
