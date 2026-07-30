# Migrate the newly tracked diff-visual map

Blocked by: Migrate the tracked decision-map corpus

## What to build

The newly tracked active `decisions/diff-visual.md` map adopts the canonical
decision-map schema without changing its twelve decision questions or answers,
exclusions, fog, reviewer decisions, or research facts.

## Acceptance

- [x] All twelve decision questions and answers, the #11 → #10 blocker, current
  fog, exclusions, reviewer decisions, and research facts are preserved.
- [x] The active map is honestly `shaping`, has canonical statuses and blockers,
  has no Handoff section, and leaves spec-writer discretion empty.
- [x] Every research citation uses a structured Path or absolute HTTP(S) URL
  record with non-empty Supports and Drift, including both tracked HTML assets.
- [x] `ValidateDecisionMap` rejects the legacy form and accepts the migrated map;
  focused maps tests, vet, and a second semantic migration pass are green and
  leave no additional migration diff.
