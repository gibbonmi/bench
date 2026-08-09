# Assert the exact approved AXI disposition

Blocked by: declare-production-axi-registry.md
Ownership fence: `cmd/bench`, `internal/conformance`, `projects/benchkit.md`
Integration surfaces: production declarations→declare-production-axi-registry.md; guidance→document-ten-principle-axi.md
Contracts: member-name set and disposition enum cross `cmd/bench`→`internal/conformance`, membership is six root queries plus worktree list versus explicit exemptions, order is irrelevant set equality, and no disposition may be absent, asserted by DS1
Closure: DS1/roots, DS1/worktree-list, DS1/exemptions

## What to build

the exact AXI set remains six approved roots plus worktree list and every other member is explicitly exempt.

## Acceptance

- [ ] [DS1] (covers CR7) the exact AXI set remains six approved roots plus worktree list and every other member is explicitly exempt.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DS1/roots | mark one approved root non-AXI | independent exact-set test | derive and require the six names |
| DS1/worktree-list | mark worktree list non-AXI | independent exact-set test | derive and require the nested member |
| DS1/exemptions | mark one operational member AXI | independent exact-set test | derive and require exact set equality |

