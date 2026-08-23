# Make the guidance name one flip author

Blocked by: retire-bench-spec-implemented.md, commit-exit-3-names-the-remainder.md
Writes: .agents/skills/bench-craft-spec/SKILL.md, .agents/commands/bench-final-check.md, projects/benchkit.md

## What to build

The spec template in `craft-spec` gains the optional `Roadmap: FT<n>` line
between `Status:` and `Decision source:`. A one-clause note says that `bench
spec retire` names that row's detail file. The skill keeps its prose budget:
trim one line elsewhere in the skill, or raise the skill's budget row in
`projects/benchkit.md` by one. `/bench-final-check` drops the `bench spec
implemented` sentence and states that the landing verb is the only flip
author. It names `bench commit`'s exit 3 as published, with the record's
`next=` restore as the repair. Every sentence is ASD-STE100 per
`references/ste-prose.md`.

## Acceptance

- [ ] The template block in `craft-spec` shows `Roadmap: <optional — the FT<n> row this spec settles>` between `Status:` and `Decision source:`.
- [ ] `.agents/skills/bench-craft-spec/SKILL.md` is at most its budgeted line count in `projects/benchkit.md`.
- [ ] `.agents/commands/bench-final-check.md` contains no `bench spec implemented` and names `bench worktree land` as the only `Status: implemented` author.
- [ ] `.agents/commands/bench-final-check.md` names `bench commit` exit 3 as published, with the record's `next=` restore as the repair.
- [ ] `.agents/commands/bench-final-check.md` still contains `bench commit -m`.
