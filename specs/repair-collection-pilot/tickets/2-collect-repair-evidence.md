# Collect attributed repair evidence

Blocked by: 1-activate-pilot.md
Writes: internal/repairpilot (new), internal/assessment/command.go, internal/assessment/cost.go, internal/assessment/harness.go, internal/assessment/plan_validate.go, internal/assessment/store.go, internal/assessment/validate.go, internal/bounds/classify.go, specs/repair-collection-pilot, reviews/repair-collection-pilot.md (new), cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: RP13, RP14, RP15, RP16, RP17, RP18, RP19, RP20, RP21, RP22, RP23, RP24, RP25, RP26, RP27, RP28, RP29, RP30, RP32, RP33, RP34, RP35, RP36, RP37, RP45, RP48, RP51, RP52, RP53, RP54, RP60, RP62, RP63, RP64, RP65

## What to build

Introduce `bench repair-pilot record --input <file>` together with its help and production route.
Extend the active pilot with explicitly imported observations, endpoints, proposals, and audits.
Preserve one sequence across assignments, sessions, findings, and review stages when the implementation source, spec, and chunk match.
Retain comparable failures, their ownership, and the failure set's completeness.
Do not infer assignment identity from overlapping timestamps.

Retain explicit activity interval bounds and provenance for the later overlap report.
An initial nonblocking observation cannot start a sequence.

Use RP-C1's document owner and update protocol for every import.
Recheck the stored hostile-state inventory through the new record operation.
Validate every supplied activity timestamp, including both interval bounds, against the active window.

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
- [ ] A first blocking failure starts the sequence without inventing a repair attempt.
- [ ] The source, spec, and chunk identify one sequence across contributor assignments.
- [ ] Explicit interval bounds retain their evidence reference.
- [ ] Every supplied activity timestamp remains inside the active observation window.
- [ ] Record import refuses hostile stored documents without changing prior evidence.
