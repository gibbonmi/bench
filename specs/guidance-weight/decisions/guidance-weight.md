# Always-loaded guidance weight (FT100)

Status: ready

## Destination

Shrink the always-loaded guidance that every session pays for: `.bench/BENCH.md`, `AGENTS.md`, `CLAUDE.md`, and the skill and phase descriptions that a harness lists on every turn. The outcome moves each phase-only clause to a file that its phases load. It also removes repeated copies, and it replaces the line budget with a byte budget.

## Notes

FT100 is the source row. FT231 owns the measurement that a rule deletion needs, and FT89 owns the guidance contradictions.
A map-owned asset stays in the map's assets folder, specs/guidance-weight/decisions/guidance-weight/assets/.

Terms:

- **Always-loaded guidance** is text that a harness loads into every session or every turn. Avoid: core docs, root docs, system prompt.
- **On-demand guidance** is a phase command, a skill body, a reference file, or `.bench/BENCH-reference.md`.
- **Clause** is one rule paragraph or bullet of always-loaded guidance. Avoid: line, item.
- **Move** puts a clause, word for word and with its anchors, into a file that its acting phases load.
- **Deletion** removes a rule that has no other copy. Avoid: cut and trim, when either word could mean a move.

Invariant 1 applies to every move. An anchor moves with its clause or with the copy that stays, and no anchor goes without a reviewer decision.
Research pairs run one Sonnet/high reader and one Haiku/xhigh reader for each slice, by the reviewer's scorecard decision of 2026-10-10.
Use craft-research for Research tickets and craft-grill for Grill tickets. The spec uses craft-skills and craft-gate.

## Decisions so far

- [Is FT100 one outcome or two?](guidance-weight/tickets/1.md): two; this map shapes the always-loaded set, and the craft-skill library audit is a later outcome.
- [Which rule places and removes an always-loaded clause?](guidance-weight/tickets/2.md): a phase-only clause moves with its anchors; a unique rule goes only after an FT231 result.
- [Which measure does the prose budget use?](guidance-weight/tickets/3.md): bytes, for every row of the guidance prose budget table.
- [Which light-path fixes go to a write delegate?](guidance-weight/tickets/4.md): every learning or drain fix; a fix that the reviewer asks for stays inline.
- [What byte limit does each existing budget row start with?](guidance-weight/tickets/6.md): the file's size at the switch, rounded up to a multiple of 100 bytes.
- [Does AGENTS.md get a byte budget row?](guidance-weight/tickets/7.md): yes, an exact row at its size after this outcome.
- [Does the kit's CLAUDE.md take the managed form?](guidance-weight/tickets/8.md): yes; the repeated sentences go.
- [Where does each always-loaded clause act, and where else does it live?](guidance-weight/tickets/5.md): see `assets/clause-inventory.md`; the build clauses and B43 act only in some phases.
- [Where do the build-orchestration clauses move?](guidance-weight/tickets/9.md): into one build reference beside the build phase, with their pointers and anchors.
- [How do the light-path clauses merge into one block?](guidance-weight/tickets/10.md): one block under the table; the repeated drain sentences go.
- [What byte bound does BENCH.md get?](guidance-weight/tickets/11.md): its size after this outcome, rounded up to a multiple of 100 bytes; about 11,600.
- [Which repeated copies go now?](guidance-weight/tickets/14.md): the B10 exception, B12, and A17; their anchors move to the kept copies.
- [Where does the batch-approval rule go?](guidance-weight/tickets/15.md): beside the sign-off stop in `/bench-write-spec`; the drain line gets its own words.
- [Which part of each listed description does its trigger need?](guidance-weight/tickets/12.md): see `assets/description-inventory.md`; content restates the body, and four defects need a fix.
- [What byte limit does a listed description get?](guidance-weight/tickets/13.md): 250 stays; the defects get fixed, and a trim waits for an FT231 trial.
- [Does this outcome ship as one spec or split?](guidance-weight/tickets/16.md): one spec for the moves and the budget; two light-path tickets for the small fixes.

## Not yet specified

## Spec-writer discretion

- How the byte count treats line endings and the final newline.
- The order of the moves inside the spec.
- Whether a check holds the kit's `CLAUDE.md` to the managed form.
- The name of the build reference file.

## Out of scope

- The craft-skill library audit and the merge of overlapping craft skills: a later FT100 outcome.
- The deletion of a rule that has no other copy: it waits for an FT231 result.
- A trim of the description content phrases: it waits for an FT231 trial of skill selection.
- The admission of this outcome into the delivery commitment: a commitment plan that the reviewer approves decides it.

## Sources

- Path: `specs/guidance-weight/decisions/guidance-weight/assets/clause-inventory.md`
  Supports: #9, #10, #14, and #15.
  Drift: an edit of `.bench/BENCH.md`, `AGENTS.md`, or an anchor row that names either file.
- Path: `specs/guidance-weight/decisions/guidance-weight/assets/description-inventory.md`
  Supports: #13 and #16.
  Drift: an edit of a listed `description:` line.
- Path: `roadmap/FT100.md`
  Supports: the scope, the anchor-first order, and the budget-measure requirement.
  Drift: a drain edit of the row.
- Path: `roadmap/FT231.md`
  Supports: the measurement rule for a rule deletion.
  Drift: an FT231 delivery or a change of its scope.
- Path: `capture/agent-performance/claude-models.md`
  Supports: the models of the research pairs.
  Drift: a change of the research routing decision.
- Path: `internal/conformance/prose_budget_test.go`
  Supports: the line count that the budget check uses today.
  Drift: a change of the budget measure.
- Path: `internal/adopt/marker.go`
  Supports: the managed CLAUDE.md form that `bench link` writes.
  Drift: a change of the managed form.
