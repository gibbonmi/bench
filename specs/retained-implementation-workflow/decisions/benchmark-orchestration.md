# Benchmark orchestration

Status: ready

## Destination

Decide the orchestration, evidence, and measurement changes proposed by the FT311 benchmark.
Keep the outcomes in one map and deliver independently useful outcomes through separate specs.

## Notes

Use craft-domain, craft-grill, and craft-synthesis. Research questions also use craft-research.
Historical research assets supply benchmark evidence. Resolved decision tickets own the current workflow choices.

The reviewer requires independent review for solo implementation.
Ticket #2 records the reviewer-selected implementation default. Portable workflow changes have not yet been implemented.
The existing roadmap owners remain in place.

The reviewer requires flexible implementation effort.

## Decisions so far

- [How is the proposal set shaped and delivered?](benchmark-orchestration/tickets/1.md): One map coordinates separate deliverable specs.
- [Does retained authorship become the default?](benchmark-orchestration/tickets/2.md): Implementation stays in one model session; reviews use subagents, with brief diagnostic consultation permitted by #6.
- [Does uncapped continuation extend beyond this benchmark?](benchmark-orchestration/tickets/3.md): Support explicit uncapped continuation with progress and stop conditions.
- [Who can change chunk boundaries during implementation?](benchmark-orchestration/tickets/4.md): The implementer can revise chunks while preserving required behavior and review checkpoints.
- [What write authority and gate expansion can the implementation session use?](benchmark-orchestration/tickets/5.md): Expansions have standing approval and enter the drain.
- [What changes after repeated failure to advance?](benchmark-orchestration/tickets/6.md): Reassess after two attempts without progress; use bench-debug or brief advice from a higher-tier model.
- [What assurance must completion evidence establish?](benchmark-orchestration/tickets/7.md): Bind terminal verification and review outcomes to their source and performer.
- [How are the independently useful outcomes split?](benchmark-orchestration/tickets/8.md): Separate specs deliver implementation workflow, continuation, completion evidence, and assessment.
- [What evidence permits a later default change?](benchmark-orchestration/tickets/9.md): Use repeated comparisons with equivalent assurance and complete cost; the reviewer decides adoption.
- [Where must missing or stale assurance block completion?](benchmark-orchestration/tickets/10.md): Missing or stale required evidence blocks chunk advancement and spec completion through the gate-owned check.
- [What assessment record and trial rule support adoption?](benchmark-orchestration/tickets/11.md): Keep local run records and require an approved plan for paid comparison trials.
- [How are implementation chunks defined and reviewed?](benchmark-orchestration/tickets/12.md): The spec author defines chunks. Each chunk receives delegated review before the next starts.
- [Who recommends the implementation line, and what line do reviews use?](benchmark-orchestration/tickets/13.md): The spec author recommends the implementation line. Reviews default to mid/high, with Astra/high for Sol implementations.

## Not yet specified

## Spec-writer discretion

## Out of scope

## Sources

- Path: `specs/retained-implementation-workflow/decisions/benchmark-orchestration/assets/benchmark-recommendations.md`
  Supports: #1 through #9 benchmark proposals and their evidence boundaries.
  Drift: re-read when the benchmark findings or recommendations change.
- Path: `roadmap/FT293.md`
  Supports: #5 prior closure and write-authority questions. The resolved ticket governs this workflow.
  Drift: re-read when the reviewer settles closure overlap.
- Path: `roadmap/FT231.md`
  Supports: #9 existing measurement ownership and comparison requirements.
  Drift: re-read when the measurement contract changes.
- Path: `roadmap/FT200.md`
  Supports: #7 existing preflight enforcement decision and oracle boundary.
  Drift: re-read when the enforcement decision changes.
- Path: `specs/retained-implementation-workflow/decisions/benchmark-orchestration/assets/frontier-reassessment.md`
  Supports: #4 through #11 historical baseline evidence. Resolved tickets replace its proposed policies.
  Drift: re-read if the frozen source identities or their findings change.
