# Choose ticket boundaries by tracer outcome

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `ROADMAP.md`
Contracts: none crosses

## What to build

Ticket authors inspect seams without treating them as automatic split points. A consolidated signal/response table keeps one independently-green behavior together across seams, splits independently useful behaviors, prefactors shared primitives, creates junction tickets for contracts neither side proves alone, points disjoint package groups back to `craft-spec`, and preserves expand-migrate-contract for mechanical refactors. The fence-breadth row demands justification rather than an automatic split, and the same edit absorbs `specs/authoring-hardening/tickets/correct-fence-breadth-evidence.md` by removing its false narrow-fences-landed-first-time claim. Record the discovered vacuous `Contracts:` template anchor on FT156 without changing the oracle in this guidance batch.

## Acceptance

- [ ] [TB1] `craft-tickets` preserves verbatim: "A seam is a reason to inspect a ticket boundary, not an automatic ticket boundary. Split only when both resulting tickets remain complete, independently green tracer outcomes."
- [ ] [TB2] one table maps all seven required signals to keep, split, prefactor, junction, spec split, expand-migrate-contract, or justify-without-auto-splitting responses without re-teaching the existing ticket definition.
- [ ] [TB3] the fence-breadth paragraph keeps its threshold, counting convention, and justification price while no longer claiming or implying that a narrow fence predicts a sound ticket.
- [ ] [TB4] FT156 records that deleting the ticket template's `Contracts:` line survives a real graded-root conformance run because no section-scoped anchor requires that line.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TB1 | replace the principle with seam-first grouping | the semantic reviewer | re-read the breakdown method, expect seam crossings to remain inspection signals rather than split commands |
| TB2 | remove the junction or prefactor response | the semantic reviewer | enumerate the seven requested signals against the table, expect one signal to have no response |
| TB3 | restore the narrow-fences-landed-first-time clause | the semantic reviewer | compare the paragraph with `specs/authoring-hardening/tickets/correct-fence-breadth-evidence.md`, expect the one-directory counterexample to falsify it |
| TB4 | omit the roadmap rider | the roadmap review | inspect FT156 after the graded-root mutation probe, expect the concrete false-green face to be absent |
