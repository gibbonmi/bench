# Document repair-ticket reslicing

Blocked by: none
Ownership fence: `CHANGELOG.md`, `package.json`
Integration surfaces: implemented repair-reslicing behavior→`CHANGELOG.md` + CH3 + CH4 + CH5 + CH6; unchanged canonical policy owner→`.agents/skills/bench-craft-tickets/SKILL.md` + CH3 + CH4 + CH5; unchanged lifecycle orchestrator→`.agents/commands/bench-implement-spec.md` + CH6; unchanged shipped-package manifest→`package.json` + CH7
Contracts: user-visible repair-reslicing summary (typed Changed-entry prose, exact receipt-fence envelope plus one-ticket-or-ordered-chain and terminal-refresh domain, placement under `[Unreleased]` then `Changed` ordering, absent or contradictory entry invalid) crossing the implemented command and skill behavior→`CHANGELOG.md`, asserted by CH1 through CH6 against the exact integrated candidate; shipped changelog membership (package-file-list type, exact `CHANGELOG.md` domain, manifest selection before package materialization ordering, missing membership invalid) crossing `package.json`→`CHANGELOG.md`, asserted by CH7 through package materialization
Closure: CH1/entry-present, CH2/unreleased-changed-placement, CH3/receipt-envelope-summary, CH4/permitted-result-summary, CH5/contained-union-summary, CH6/terminal-refresh-summary, CH7/changelog-shipped

## What to build

Close review finding `S1-missing-repair-reslicing-changelog` with one concise `Changed` entry under `[Unreleased]`. State that a validated receipt's fence is the maximum repair envelope, that independently-green reslicing may yield one repair ticket or an ordered reciprocal chain whose union remains inside that envelope, and that the original blocked assignment refreshes only after the terminal repair ticket lands. Do not restate the slicing algorithm or add another policy owner.

## Acceptance

- [ ] [CH1] (covers local) One concise typed repair-reslicing entry exists.
- [ ] [CH2] (covers local) The entry appears under `[Unreleased]` `Changed`.
- [ ] [CH3] (covers local) The entry identifies the validated receipt fence as the maximum repair envelope.
- [ ] [CH4] (covers local) The entry permits one repair ticket or an ordered repair chain.
- [ ] [CH5] (covers local) The entry keeps the union of repair-ticket fences inside the receipt envelope.
- [ ] [CH6] (covers local) The entry permits original-assignment refresh only after the terminal repair ticket lands.
- [ ] [CH7] (covers local) The unchanged package manifest materializes `CHANGELOG.md` on the shipped surface.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CH1/entry-present | delete the entry | changelog semantic review (no executable changelog-content owner exists) | compare the exact candidate's canonical behavior to `CHANGELOG.md`, require the missing typed entry as an attributable finding, restore the entry, and rerun the same review green |
| CH2/unreleased-changed-placement | move the entry outside `[Unreleased]` `Changed` | changelog semantic review (no executable placement owner exists) | compare the exact candidate's `CHANGELOG.md` section structure, require the misplaced entry as an attributable finding, restore the entry, and rerun the same review green |
| CH3/receipt-envelope-summary | replace maximum-envelope wording with an indivisible one-ticket fence | changelog semantic review (no executable prose-truth owner exists) | compare the mutated entry to the exact candidate's canonical skill behavior, require the contradiction as an attributable finding, restore the entry, and rerun the same review green |
| CH4/permitted-result-summary | replace one-ticket-or-chain wording with a chain-only result | changelog semantic review (no executable prose-truth owner exists) | compare the replacement independently to the exact candidate's canonical skill behavior, require the contradiction as an attributable finding, restore the entry, and rerun the same review green |
| CH4/permitted-result-summary | replace one-ticket-or-chain wording with a one-ticket-only result | changelog semantic review (no executable prose-truth owner exists) | compare the replacement independently to the exact candidate's canonical skill behavior, require the contradiction as an attributable finding, restore the entry, and rerun the same review green |
| CH5/contained-union-summary | permit one chain ticket to escape the receipt envelope | changelog semantic review (no executable prose-truth owner exists) | compare the mutated entry to the exact candidate's canonical skill behavior, require the contradiction as an attributable finding, restore the entry, and rerun the same review green |
| CH6/terminal-refresh-summary | replace terminal-refresh wording with permission to refresh after any non-terminal repair ticket lands | changelog semantic review (no executable prose-truth owner exists) | compare the mutated entry to the exact candidate's canonical command behavior, require the contradiction as an attributable finding, restore the entry, and rerun the same review green |
| CH7/changelog-shipped | remove `CHANGELOG.md` from the package file list | npm package materialization | mutate the real manifest, run `npm pack --dry-run --json`, require the materialized file list to omit `CHANGELOG.md`, restore the manifest, and rerun the same command requiring `CHANGELOG.md` present |
