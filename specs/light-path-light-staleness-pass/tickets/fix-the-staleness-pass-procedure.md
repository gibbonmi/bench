# Fix the staleness pass procedure

Blocked by: none
Writes: .agents/skills/bench-implement-spec/references/staleness-pass.md, .agents/commands/bench-implement-spec.md, CHANGELOG.md
Covers: none

## What to build

The reviewer decided that the spec staleness pass of `/bench-implement-spec`
is a light assessment with a fixed short procedure. The implementation
command names only the mid tier, so each orchestrator writes its own charge.
One such charge took 98 tool calls, and its findings were non-behavioral.

One reference file under the phase adapter skill owns the procedure. The
charge is the mid tier, medium effort, one iteration, and 15 tool calls. The
procedure lists the drift since the reviewed graph commit and reads only the
drifted files. It then classifies each contradiction and returns blocking
contradictions only. The implementation command points to the file in one
sentence and restates nothing.

New build preflight rows are out of scope.

## Acceptance

- [ ] One reference file states the charge, the drift procedure, and the blocking-only return.
- [ ] The implementation command charges the pass through one sentence that names the reference file.
- [ ] No other guidance file restates the charge line, the budget, or the procedure steps.
- [ ] The prose mechanics check, the prose budget check, and the root conformance pass stay green.
