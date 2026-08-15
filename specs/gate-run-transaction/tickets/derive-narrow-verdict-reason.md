# Derive the narrow-verdict reason from the registry

Blocked by: register-verdict-record-classes.md
Writes: internal/gate/verdict.go, internal/gate (test file)

## What to build

Story 3, consumer half. `narrowVerdictReason` answers from the registry row
the loader selected — per-row attribute or name predicate, discretion — instead
of the ad-hoc `partitions()`/`checkPartitions()` predicates, so the reuse
refusal reads the registry too. The observable reason string stays exactly
`partial verdict` for the three narrow classes and empty for the full class.
Consumes the registry the blocker landed; adds no row.

Return note (not acceptance): the three per-row mutations (each narrow row
made to report the full class's reason) and their observed reds.

## Acceptance

- [ ] Partial, check-partial, and combined-partial records each yield a non-reusable inspection with reason `partial verdict`; a full record is reusable (covers GR2)
- [ ] `narrowVerdictReason` no longer references `partitions()`/`checkPartitions()` directly
