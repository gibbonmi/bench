# Make the author line defaults

Blocked by: none
Writes: .agents/skills/bench-craft-line/SKILL.md, .agents/commands/bench-implement-spec.md, projects/benchkit.md
Covers: none

## What to build

A spec-backed build sets the effort of each ticket author from its chunk, and it scans the fence of each ticket before the author charge. `craft-line` states the effort rule in its ticketed-build stage defaults. `/bench-implement-spec` states the fence scan before the charge preparation. The fence scan widens `Writes:` under the approved plan-expansion policy of `.bench/BENCH.md`. The fence scan takes its own paragraph, so the prose budget of `/bench-implement-spec` moves by two lines. Each rule has one owner, and no other file restates it.

## Acceptance

- [ ] `craft-line` sets the ticket author effort from the harder-chunk mark: medium for an unmarked chunk, and the declared effort for a harder chunk.
- [ ] `/bench-implement-spec` has the orchestrator run `bench consumers` before each author charge and widen `Writes:` for each reached caller, test, or fixture.
