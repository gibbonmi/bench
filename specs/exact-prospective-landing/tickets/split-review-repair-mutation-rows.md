# Split review repair mutation rows

Blocked by: normalize-review-repair-ticket-metadata.md
Ownership fence: `specs/exact-prospective-landing/tickets/preserve-executable-spec-mode.md`, `specs/exact-prospective-landing/tickets/close-adapter-blocker-metadata.md`
Integration surfaces: executable-mode repair mutation table→`specs/exact-prospective-landing/tickets/preserve-executable-spec-mode.md`; blocker-metadata repair mutation table→`specs/exact-prospective-landing/tickets/close-adapter-blocker-metadata.md`; one-row-per-acceptance-ID rule→existing `.agents/skills/bench-craft-tickets/SKILL.md` plus SM1-SM2
Contracts: each repair acceptance ID and its independent subject mutation cross `specs/exact-prospective-landing/tickets/preserve-executable-spec-mode.md` and `specs/exact-prospective-landing/tickets/close-adapter-blocker-metadata.md`→fresh Standards review, asserted by SM1-SM2 against the complete mutation tables

## What to build

Complete the two accepted-review repair records by giving every MR and BR acceptance
ID its own concrete subject mutation, independent owner, and bounded public operation.
Do not change their implementation scope, acceptance criteria, or metadata fields.

## Acceptance

- [ ] [SM1] The executable-mode repair has distinct MR1, MR2, and MR3 mutation rows that respectively catch executable-mode collapse, non-executable-mode drift, and a package-level regression outside the focused transition pair.
- [ ] [SM2] The blocker-metadata repair has distinct BR1 and BR2 mutation rows that respectively catch an omitted producer basename and a malformed/unresolved complete blocker record.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SM1 | omit any one MR mutation row | exact mutation-table shape probe plus fresh Standards review | inspect the executable-mode repair record and expect its acceptance-to-mutation ID sets to differ |
| SM2 | combine BR1 and BR2 into one row again | exact mutation-table shape probe plus fresh Standards review | inspect the blocker-metadata repair record and expect its acceptance-to-mutation ID sets to differ |
