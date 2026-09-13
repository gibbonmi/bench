# Collect attributed repair evidence

Blocked by: 1-activate-pilot.md
Writes: internal/repairpilot (new), specs/repair-collection-pilot, reviews/repair-collection-pilot.md (new)
Covers: RP13, RP14, RP15, RP16, RP17, RP18, RP19, RP20, RP21, RP22, RP23, RP24, RP25, RP26, RP27, RP28, RP29, RP30, RP31, RP32, RP33, RP34, RP35, RP36, RP37, RP48

## What to build

Extend the active pilot with explicitly imported observations, endpoints, proposals, and audits.
Preserve one sequence across assignments, sessions, findings, and review stages when the implementation source and chunk match.
Retain comparable failures, their ownership, and the failure set's completeness.
Do not infer assignment identity from overlapping timestamps.

Use RP-C1's document owner and update protocol for every import.
The command accepts the tenth sequence endpoint, then refuses later observations.
The deadline stops collection even when the sample contains fewer than ten completed sequences.
Audits of existing evidence can continue after cutoff.
A delayed observation cannot alter the frozen sample.

Implement the structural evidence rules and retain unknown or conflicting labels.
The software does not certify the truth of a native reference.
RP-C3 supplies the full evidence report and operator audit procedure.
The stored document remains the shared value contract between collection and reporting.

## Acceptance

- [ ] Every accepted observation retains explicit sequence, assignment, session, and source identity.
- [ ] Several findings and both review stages retain one sequence for the same chunk.
- [ ] Failure completeness and ownership retain all source vocabulary values.
- [ ] Verified closure and reviewer handoff retain distinct endpoints.
- [ ] Ten completed sequences or the original fourteen-day deadline stops new observations.
- [ ] Identical imports count once and conflicting imports preserve the previous evidence.
- [ ] Audits retain proposals and classify unsupported or conflicting progress as unknown.
