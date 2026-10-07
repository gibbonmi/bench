# Check source sizing premises through real proposals

Blocked by: 03-propose-helper-and-owner-closure.md
Writes: internal/preflight/closure.go, internal/preflight/decision.go, internal/preflight/gather.go, internal/preflight/ownership (new), internal/preflight/preflighttest/fixture.go, internal/preflight/proposal.go, internal/structure/facts.go, internal/structure/structure.go, internal/structure/structure_test.go
Covers: OC32, OC33, OC34, OC35, OC36, OC37, OC38

## What to build

The real proposal reports changed line counts, complete tracked glob members, and canonical current headroom from one explicit source pair.
Consume the reviewed OC-C2 loader and write-entry contract.
Add sizing.go and TestSizingPremises in premises_test.go.
Track dot members and literal metacharacter paths through filepath.Match's documented grammar.
Unsupported recursive syntax, escaping scopes, or incomplete member lists remain diagnostic.

Expose sizing through existing structure/facts.go and its shared helpers.
Do not derive another limit, grant, newline count, Growth rule, or directory-count policy.
Historical counts name their actual source; current headroom requires tip.
Exact declared deletions affect net files; a future net-line estimate remains unobserved.
File Growth and whole-tree crowding remain distinct.

Put the projection in the 34-line facts.go rather than growing the 396-line structure.go.
No successor pays this checkpoint's headroom debt or supplies its sizing evidence.
Global charge certification remains inactive until all collectors exist.

This ticket belongs to OC-C3.
The explicit predecessor edge orders every overlapping write.
No unimplemented successor is required for this ticket's owned acceptance rows.

## Acceptance

- [ ] OC32: A line-count premise reports its changed current value.
- [ ] OC33: A glob-member premise reports its omitted tracked dot member.
- [ ] OC34: A glob premise retains literal metacharacter path identity.
- [ ] OC35: An unsupported recursive glob remains diagnostic.
- [ ] OC36: A current file-headroom premise uses the canonical limit.
- [ ] OC37: A net file plan counts its exact declared deletion.
- [ ] OC38: Future net-line estimates remain unobserved estimates.

## Verification

Run these future checks through their actual production entries.
Use a few sufficient multi-case witnesses; do not create a redundant test per row.
Retain the named omission mutation, its behavioral failure, restoration, and passing command result.
Compilation failure supplies no behavioral evidence.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./internal/structure`

Focused producer command: `bench test --package ./internal/preflight/... --run TestSizingPremises`.
Named mutation identity: `oc-04-omission`.
This obligation uses the named production seam and observable omission in this ticket, not a compile-failure substitute.

Mutation obligation: Ignore a tracked dot member, copy a stale limit, or omit a deletion from net files. The real source-sizing sentinel must fail.

The completion plan names this ticket's focused command and mutation identity.
Execute bench probe against the specified production seam after its code exists, then restore and rerun that exact focused command.
If the planned mutation cannot reach that producer, amend the enabling plan before verification; no substitute earns completion.

Read the canonical structure report before dispatch.
Remeasure Growth from the actual preceding tip and whole-tree crowding.
Keep cohesive extraction and necessary pin relocation in this same ticket.
