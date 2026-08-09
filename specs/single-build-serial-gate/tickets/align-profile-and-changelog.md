# Align the profile and changelog

Blocked by: contract-ordinary-build-census.md, contract-run-directory-lifecycle.md
Ownership fence: `projects/benchkit.md`, `CHANGELOG.md`
Integration surfaces: closed build census→contract-ordinary-build-census.md; terminal run lifecycle→contract-run-directory-lifecycle.md; current gate execution profile→`projects/benchkit.md`; user-visible kit change→`CHANGELOG.md`
Contracts: final run-scoped build count, serial phase schedule, and independent exception policy cross implementation→current-state descriptions in `projects/benchkit.md` and `CHANGELOG.md`, membership is one profile statement and one changelog entry, ordering is document only the contracted behavior after its tests land, asserted by DC1 against the final source and census names
Closure: DC1/profile-build-count, DC1/profile-serial-schedule, DC1/profile-exception-boundary, DC1/changelog-run-scope

## What to build

Replace the profile's concurrent-phase and ordinary-build descriptions with the final current state and add the kit-visible change to the changelog. State the resulting behavior, not the migration history or the superseded durable-store design.

## Acceptance

- [ ] [DC1] (covers local) the Bench kit profile and changelog describe one private exact-snapshot build per top-level run, one serial phase schedule, cleanup with no later reuse, and the closed independent-proof boundary without advertising excluded capabilities.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DC1/profile-build-count | leave the old phase-owned or concurrent build description | profile conformance/read-through test | compare the current phase table and owner census to the profile and require one run-scoped build statement |
| DC1/profile-serial-schedule | advertise four concurrent phase processes | profile conformance/read-through test | compare scheduler max-active policy to the profile and require serial wording |
| DC1/profile-exception-boundary | omit or broaden the independent proof classes | profile/census consistency test | compare named profile classes to the exact census and require parity |
| DC1/changelog-run-scope | describe durable or cross-run reuse | changelog review check | read the entry against the no-survivor lifecycle test and require run-scoped wording |
