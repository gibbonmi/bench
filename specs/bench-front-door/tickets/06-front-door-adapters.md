# Ship the /bench and $bench front-door adapters

Blocked by: 05-rename-what-next-to-drain.md
Writes: .agents/commands/bench.md, .agents/skills/bench, .bench/BENCH.md, .bench/BENCH-reference.md

## What to build

`.agents/commands/bench.md`: model-invocable (no `disable-model-invocation` key),
`description` says route from observed state; body: run `bench status --route`, take
the one row; if `command` opens `/bench-` read `.agents/commands/<phase>.md` completely
and follow it as the active phase; otherwise run the command exactly and report; load
nothing else; on an empty command report the state and stop. `.agents/skills/bench/`:
SKILL.md (`name: bench`, references the command file, same adapter shape as the others)
+ `agents/openai.yaml` with `allow_implicit_invocation: false`, running
`bench status --route --harness codex`. Name `/bench` and `$bench` in `.bench/BENCH.md`
(commands bullet + workflow) and the BENCH-reference mapping.

Covers: R33, R34, R35

## Acceptance

- [ ] Conformance `checkCommandGuideReferences`, `checkCodexCommandAdapters`, `checkClaudeSkillMirror` green with the new pair; no `.claude/skills/bench` entry.
- [ ] `bench.md` has no `disable-model-invocation` key and names `bench status --route`.
- [ ] Gate green.
