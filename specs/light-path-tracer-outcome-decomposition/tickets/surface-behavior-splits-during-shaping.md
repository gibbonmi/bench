# Surface behavior splits during shaping

Blocked by: none
Ownership fence: `.agents/commands/bench-shape-idea.md`
Contracts: none crosses

## What to build

Shaping resolves the scope into deliverable outcomes instead of producing a seam inventory. When the shaped scope contains two independently useful behaviors, shaping makes that split visible as a reviewer decision before a spec can bundle them by default.

## Acceptance

- [ ] [BS1] shaping names deliverable outcomes as its scope result and does not inventory engineering seams as decomposition units.
- [ ] [BS2] two independently useful behaviors surface a split decision during shaping rather than surviving silently into one bundled spec.
- [ ] [BS3] spec authoring retains ownership of engineering seams unless the reviewer explicitly chooses one while shaping.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BS1 | direct shaping to inventory seams | the semantic reviewer | re-read the ownership and exit sections, expect engineering decomposition to remain outside shaping |
| BS2 | defer the first visible behavior split to spec authoring | the semantic reviewer | shape two independently useful behaviors, expect a decision ticket before readiness |
| BS3 | move ordinary seam choice into shaping | the consistency review | compare Ownership with `/bench-write-spec`, expect spec authoring to retain the engineering seam decision |
