# Normalize review repair ticket metadata

Blocked by: preserve-executable-spec-mode.md, close-adapter-blocker-metadata.md
Ownership fence: `specs/exact-prospective-landing/tickets/preserve-executable-spec-mode.md`, `specs/exact-prospective-landing/tickets/close-adapter-blocker-metadata.md`
Integration surfaces: executable-mode repair metadata→`specs/exact-prospective-landing/tickets/preserve-executable-spec-mode.md`; blocker-metadata repair metadata→`specs/exact-prospective-landing/tickets/close-adapter-blocker-metadata.md`; required ticket shape→existing `.agents/skills/bench-craft-tickets/SKILL.md` plus rows NM1 and NM2
Contracts: executable-mode repair metadata crosses `specs/exact-prospective-landing/tickets/preserve-executable-spec-mode.md`→semantic review and lifecycle parsing, asserted by NM1 against the complete ticket record; blocker-metadata repair metadata crosses `specs/exact-prospective-landing/tickets/close-adapter-blocker-metadata.md`→semantic review and lifecycle parsing, asserted by NM2 against the complete ticket record

## What to build

Give both accepted-review repair tickets the complete current ticket shape. Put all
four lifecycle/review metadata fields directly below each title, keep every field on
one line, backtick every ownership-fence path, and preserve the already-integrated
behavioral acceptance and mutation rows unchanged.

## Acceptance

- [ ] [NM1] `preserve-executable-spec-mode.md` declares its landing-owner blocker, exact two-file fence, every real integration surface, and the transformed spec byte/mode contract crossing `internal/spec` into `internal/landing`.
- [ ] [NM2] `close-adapter-blocker-metadata.md` declares the four producer blockers it derives from, its exact one-file fence, the producer/dependent resolution surfaces, and the producer-basename contract its BR rows assert.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| NM1 | omit one required metadata line or move the fence below Acceptance again | fresh Standards review plus exact ticket-shape probe | inspect the executable-mode repair ticket and expect the missing or misplaced field to fail the required shape |
| NM2 | omit one producer basename or one required metadata line | fresh Standards review plus exact ticket-shape probe | inspect the blocker-metadata repair ticket and expect incomplete dependency resolution or shape |
