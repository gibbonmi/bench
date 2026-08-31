# Document gate API contracts

Blocked by: none
Writes: internal/gate/gate.go, internal/gate/decision.go, internal/gate/engine.go, internal/gate/phases.go
Covers: LF18

## What to build

Document Resolve, Run, RunContext, RunAndRecord, RunAndRecordContext,
RunCommand, Command, Execute, ExecuteReusingFreshGreen, Decide, InspectTree,
ExecuteTree, ValidateProjectGreen, BenchkitPhases, SystemSuite,
SystemSuiteRuns, BaseTestArgv, PhasesCommand, and PinCommand. State selection,
execution, cancellation, and evidence contracts. Do not edit FT281 files.

## Acceptance

- [ ] All nineteen named functions state their caller-facing contracts.
- [ ] Selection comments remain distinct from execution comments.
- [ ] Cancellation and evidence ownership are explicit.
- [ ] No tag-census file changes.
