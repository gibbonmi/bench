# Place the fix at the shared owner

Blocked by: none
Ownership fence: `.agents/commands/bench-debug.md`, `CHANGELOG.md`, `specs/light-path-debug-edit-owner/tickets/place-the-fix-at-the-shared-owner.md`
Integration surfaces: canonical debug guidance→`.agents/commands/bench-debug.md`; Claude command consumer→`.claude/commands/bench-debug.md` + EO3; Codex phase adapter→`.agents/skills/bench-debug/SKILL.md` + EO3; user-facing change→`CHANGELOG.md`
Contracts: the canonical debug phase crosses `.agents/commands/bench-debug.md`→Claude and Codex phase surfaces through the tracked command symlink and explicit adapter, asserted by EO3; absence means a harness can load stale debug guidance
Closure: EO1/edit-owner, EO2/placement-separation, EO3/portable-consumers

## What to build

The debug phase distinguishes where a regression test attaches from where a repair lands. Before editing, it enumerates the relevant callers and places a uniform invariant once at their narrowest honest shared owner. The regression remains at the highest seam where the reported failure is observable; caller-specific policy stays with its actual owner, and paths with no shared owner become the Phase 6 architecture finding rather than repeated caller patches.

## Acceptance

- [ ] [EO1] Phase 5 enumerates the relevant callers before editing and places a uniform invariant once at their narrowest shared function or module.
- [ ] [EO2] Phase 5 keeps the regression test at the highest observable seam, refuses to force differing caller policy into a shared helper, and routes paths with no honest shared owner to the Phase 6 architecture finding.
- [ ] [EO3] the Claude command symlink and Codex phase adapter consume the edited canonical command without a copied rule, and `CHANGELOG.md` records the user-facing change.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| EO1 | remove caller enumeration or place the repair at the reported call site without tracing sibling paths | the semantic reviewer | read Phase 5 against two affected callers converging on one invariant owner, expect one repair at that owner rather than one patch per caller |
| EO2 | collapse the test seam and edit owner into one placement or force caller-specific policy into their shared helper | the semantic reviewer | trace callers with different policy through one helper, expect the test at the highest observable seam and each policy at its honest owner |
| EO3 | copy the rule into a harness adapter or omit the changelog entry | the conformance gate and consistency reviewer | resolve `.claude/commands/bench-debug.md` and `.agents/skills/bench-debug/SKILL.md` to the canonical command, then run the path-scoped `bench commit` and inspect `CHANGELOG.md` |
