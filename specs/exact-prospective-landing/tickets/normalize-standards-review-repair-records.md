# Normalize Standards review repair records

Blocked by: allow-clean-assignment-refresh.md, allow-already-covered-clean-checkpoint.md, allow-already-covered-clean-integration.md, allow-verified-clean-provisional-release.md, extend-bounded-operation-journal.md, remove-circular-review-meta-tickets.md, commit-r2-oracle-fixture-base.md, realign-runtime-commit-fixtures-with-prospective-authorization.md, reconcile-authorized-snapshot-after-publication.md
Ownership fence: `specs/exact-prospective-landing/tickets/allow-clean-assignment-refresh.md`, `specs/exact-prospective-landing/tickets/allow-already-covered-clean-checkpoint.md`, `specs/exact-prospective-landing/tickets/allow-already-covered-clean-integration.md`, `specs/exact-prospective-landing/tickets/allow-verified-clean-provisional-release.md`, `specs/exact-prospective-landing/tickets/extend-bounded-operation-journal.md`, `specs/exact-prospective-landing/tickets/remove-circular-review-meta-tickets.md`, `specs/exact-prospective-landing/tickets/commit-r2-oracle-fixture-base.md`, `specs/exact-prospective-landing/tickets/realign-runtime-commit-fixtures-with-prospective-authorization.md`, `specs/exact-prospective-landing/tickets/reconcile-authorized-snapshot-after-publication.md`
Integration surfaces: clean-refresh record→`specs/exact-prospective-landing/tickets/allow-clean-assignment-refresh.md`; clean-checkpoint record→`specs/exact-prospective-landing/tickets/allow-already-covered-clean-checkpoint.md`; clean-integration record→`specs/exact-prospective-landing/tickets/allow-already-covered-clean-integration.md`; clean-release record→`specs/exact-prospective-landing/tickets/allow-verified-clean-provisional-release.md`; operation-journal record→`specs/exact-prospective-landing/tickets/extend-bounded-operation-journal.md` plus existing consumer `internal/specbuild/state.go`; circular-meta cleanup record→`specs/exact-prospective-landing/tickets/remove-circular-review-meta-tickets.md`; R2 fixture record→`specs/exact-prospective-landing/tickets/commit-r2-oracle-fixture-base.md`; runtime-fixture record→`specs/exact-prospective-landing/tickets/realign-runtime-commit-fixtures-with-prospective-authorization.md`; authorized-snapshot record→`specs/exact-prospective-landing/tickets/reconcile-authorized-snapshot-after-publication.md`; required ticket shape→existing `.agents/skills/bench-craft-tickets/SKILL.md` plus rows NS1 through NS4
Contracts: each repaired ticket remains the single lifecycle record for its already-integrated behavioral rows; its four discovery fields name only resolvable blockers and map every contract crossing to a fence path, existing path plus row, blocker, or dependent, while each red-mutation row names the acceptance criterion, concrete mutation, independent existing test or Standards owner, and public operation that exposes it

## What to build

Normalize the nine repair-ticket records identified by the exact-candidate Standards
review. Add missing canonical discovery fields, replace acceptance-style mutation
checkboxes with the required four-column tables, map conceptual integration surfaces
to concrete paths or lifecycle relationships, name the unchanged operation-journal
state consumer, and remove deleted meta-ticket basenames from the cleanup ticket's
blocker list. Preserve all already-integrated behavioral claims and their existing
focused commands. Do not change production code, tests, the spec, or gate policy.

## Acceptance

- [ ] [NS1] The R2 and runtime-fixture repair tickets place all four canonical discovery fields directly below their titles with exact fences and concrete integration mappings.
- [ ] [NS2] The refresh, clean-checkpoint, clean-integration, clean-release, and operation-journal tickets map each integration surface to a fence path, existing path plus row, blocker, or dependent; the journal record names `internal/specbuild/state.go` as its unchanged bound consumer.
- [ ] [NS3] Every affected repair ticket's mutation section is a criterion/mutation/owner/operation table whose owner is independent of the implementation and whose command or lifecycle operation is concrete and public.
- [ ] [NS4] The circular-meta cleanup ticket names only existing sibling blockers and preserves its substantive producer-set and executable-mode evidence.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| NS1 | omit or misplace one canonical discovery field in either prerequisite-fixture record | fresh Standards axis against the current ticket-shape standard | inspect both exact-candidate records against `.agents/skills/bench-craft-tickets/SKILL.md`; the missing discovery field is a hard finding |
| NS2 | leave one conceptual surface unmapped or omit the unchanged journal state consumer | fresh Standards axis against the current integration-surface standard | trace every declared crossing to its fence path, existing path plus row, blocker, or dependent; the unresolved crossing is a hard finding |
| NS3 | retain one acceptance-style mutation checkbox or omit a table column, independent owner, or public operation | fresh Standards axis against the current mutation-table standard | inspect every affected `## Red mutations` section; the malformed row is a hard finding |
| NS4 | retain either deleted meta-ticket basename in `Blocked by:` | sibling-ticket blocker resolution | resolve every cleanup-ticket blocker basename in the exact candidate; the deleted sibling is absent |
