# Read the known issues at the spec phase entry

Blocked by: none
Writes: .agents/commands/bench-write-spec.md, projects/benchkit.md, CHANGELOG.md
Covers: none

## What to build

At the entry of `/bench-write-spec`, the invoking session reads `docs/known-issues.md` before the research fan-out. The session names each row whose issue can affect a step of the phase, from the research to the stage close. Each charge of an affected step names that path and each affected row by its issue, and it copies no row. The invoking session applies the workaround of each affected row at its own steps.

The command file states this rule once. The research question that read the file and the probe input of the author step go away, because the new step replaces them. The prose budget of the command file grows by the net count of the added lines.

## Acceptance

- [ ] The `Who runs this phase` section of `.agents/commands/bench-write-spec.md` has one step that reads `docs/known-issues.md` before the research fan-out.
- [ ] The command file names `docs/known-issues.md` only in that step.
- [ ] The step copies no row of `docs/known-issues.md`.
- [ ] The `guidance-prose-budgets` check and the `prose-mechanics` check pass on the changed tree.
