# Document the enforced ten-principle contract

Blocked by: validate-axi-registry.md, assert-approved-axi-disposition.md
Ownership fence: `.agents/skills/bench-craft-cli/SKILL.md`, `projects/benchkit.md`, `CHANGELOG.md`, `internal/conformance`
Integration surfaces: validated registry→validate-axi-registry.md; exact set→assert-approved-axi-disposition.md
Contracts: ordered ten-principle names, approved member set, exemption set, and deferred features cross `projects/benchkit.md`→`.agents/skills/bench-craft-cli/SKILL.md`, membership/order match the published contract, and no principle/scope may be absent, asserted by DG1
Closure: DG1/ten, DG1/scope, DG1/set, DG1/exemptions, DG1/deferrals, DG1/changelog

## What to build

canonical guidance states all ten principles, exact scope, exemptions, and current deferrals without widening public behavior.

## Acceptance

- [ ] [DG1] (covers CR8) canonical guidance states all ten principles, exact scope, exemptions, and current deferrals without widening public behavior.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DG1/ten | remove principle 8 | docs conformance test | inspect and require all ten ordered identities |
| DG1/scope | claim whole-binary AXI | docs/registry consistency test | compare to production disposition and require per-surface scope |
| DG1/set | omit one approved member | docs/registry consistency test | require exact approved set |
| DG1/exemptions | omit operational exemptions | docs/registry consistency test | require explicit non-AXI posture |
| DG1/deferrals | advertise fields or help output | docs current-state test | require unsupported features absent |
| DG1/changelog | omit the user-visible entry | changelog conformance test | require one concise typed entry |

