# Publish the complete AXI guidance contract

Blocked by: grade-the-approved-axi-query-set.md
Writes: `.agents/skills/bench-craft-cli/SKILL.md`, `CHANGELOG.md`, `internal/conformance/`

## What to build

Rewrite `craft-cli` as the current ten-principle contract, including content-first output, contextual disclosure, the exact `help[N]{cmd,why}` and honest zero-row forms, and exactly the approved six root queries plus `worktree list`. Add a conformance cross-check against the production registry and a concise typed changelog entry; do not advertise `--fields` or any operational family.

## Acceptance

- [ ] [GC1] (covers QD4) `craft-cli` states exactly ten named principles, explains their per-surface application and help-block contract, and advertises exactly the registry-declared approved AXI set.
- [ ] [GC2] (covers QD4) conformance turns red when guidance removes a principle, changes the count, omits an approved member, adds an operational member, or advertises unshipped `--fields` behavior.
- [ ] [GC3] (covers QD4) the changelog records the user-visible guidance and query-disclosure completion without duplicating the production registry.
