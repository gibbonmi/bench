# Rename /bench-what-next to /bench-drain with a one-release alias

Blocked by: 04-status-route-flag.md
Writes: .agents/commands, .agents/skills/bench-drain, .agents/skills/bench-what-next, internal/status, internal/roadmap, internal/learnings, internal/axi/action.go, internal/dashboard, internal/adopt/init.go, internal/anchors/registry_data.go, internal/conformance, tests/canary, .bench/BENCH.md, .bench/BENCH-reference.md, README.md, CONTEXT.md, ROADMAP.md, projects/benchkit.md, .agents/skills/bench-craft-synthesis/SKILL.md, docs/reporesident-distillation.md, docs/greenfield-build-sequence.md

## What to build

`git mv` the phase file to `bench-drain.md` and its Codex pair to `bench-drain/`;
leave `bench-what-next.md` and `bench-what-next/` behind as thin aliases (same
`disable-model-invocation` posture, body: renamed, read and follow `bench-drain.md`).
Rename every production string (drain/roadmap rows, roadmap help row, learnings
HarnessPhase, `validHarnessPhase` accepting both, dashboard prose, init scaffold);
re-point every anchor `File:`/needle/diagnostic; update the path-hardcoded conformance
tests and the `registry_test` fixture names; rename the six `what-next-*` canary dirs to
`drain-*` with their fixture files and EXPECT/BASE/MUTATE.json; update live prose
(README, BENCH.md, BENCH-reference mapping row, CONTEXT, benchkit profile, sibling
command files, craft-synthesis, ROADMAP.md, the two docs guides). Leave dated occurrence
lines, `docs/audits/`, `capture/` untouched.

Covers: R41, R42, R43, R44, R45, R46, R47

## Acceptance

- [ ] `bench anchors .agents/commands/bench-drain.md` lists every row formerly on `bench-what-next.md`; `bench canary` green; conformance green; `bench skills-index --check` green.
- [ ] Every board/roadmap/learnings/init/dashboard string says `/bench-drain`; `validHarnessPhase` accepts both names.
- [ ] `bench handoff --next /bench-what-next` still validates on a fixture (alias present).
- [ ] With the two alias files removed in a throwaway copy, the stale-command sweep names every remaining old reference (record the probe output).
- [ ] Gate green.
