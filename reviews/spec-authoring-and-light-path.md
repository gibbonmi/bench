## Standards

Finding count: 2. Worst issue: duplicated verification-loop policy across the command and project profile.

- `auto-fix` — `AGENTS.md`'s one-source-per-fact rule is violated by `projects/benchkit.md` Lines repeating the loop inputs, retry policy, reporting shape, advisory status, and sign-off stop owned by `.agents/commands/bench-write-spec.md`. Keep cached routing and a lean stage identity in the profile; point to the command for process.
- `auto-fix` — `.agents/skills/bench-craft-delegate/SKILL.md` introduces the otherwise undefined phrase `bound alias`, while `craft-line` uses `resolved model id`. Replace it with `resolved bound model id` so tier-versus-id routing remains unambiguous.

## Spec

Finding count: 1. Worst issue: WF4 leaves one ticketless state on an unexplained red path.

- `auto-fix` — `specs/spec-authoring-and-light-path/spec.md` WF4 requires ticketless spec-backed runs to return to `/bench-write-spec`, but `.agents/commands/bench-implement-spec.md` checks only for an absent directory. `internal/preflight/gather.go` distinguishes absent from present-empty, so a present directory with no ticket files reaches red preflight. Cover absent or empty/no-ticket-file and update the WF4 anchor fixture.

## Coverage

Finding count: 1. Worst issue: WF10's two-loop profile routing can silently rot.

- `auto-fix` — `projects/benchkit.md` can lose loop-1 spec-before-slicing or loop-2 breakdown-after-slicing semantics while retaining only the two currently required headings. In the same lean profile repair as the Standards finding, add exact Require anchors and mutations for the two cached-routing bullets without copying the command's retry or stop policy.

Raw findings: Standards 2, Spec 1, Coverage 1. De-duplicated repair targets: 3.
