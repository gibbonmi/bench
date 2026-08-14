# Document the --reviewer override flag

Blocked by: verification-loops.md, cross-harness-reviewer-recipes.md
Writes: .agents/commands/bench-write-spec.md, internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors/

## What to build

bench-write-spec.md documents `--reviewer <tier-or-model> [effort]` on the
verification loops: a tier resolves through the invoking harness's own
`lines.env` column (`--reviewer mid xhigh` under Codex resolves
`BENCH_CODEX_MID` at xhigh), a model id must already be bound in `lines.env`
and an unbound id is refused, an own-family id runs through the native agent
surface and a cross-family id through the recipes reference, which this
passage cites by path. Blocked by verification-loops.md (the flag rides the
loop clause that ticket installs) and cross-harness-reviewer-recipes.md (the
cited file must exist before the citation lands, or the reference dangles).
One long-needle Require anchor pins the grammar sentence — flag,
tier-resolution, and refusal — with a new fixture biting both halves.
Registry and fixture paths land serially across the whole spec.

## Acceptance

- [ ] the grammar sentence with tier resolution and unbound-id refusal is
      documented, cites the recipes file, and its long-needle fixture bites
      both halves (covers WF8)
