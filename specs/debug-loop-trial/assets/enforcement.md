# Enforcement and reader inventory

Baseline: `2f266312`, current main at phase entry, 2026-09-15.
Subject: the loop reference, the debug-phase changes, and the trial flag.

## Existing enforcement

| Read source | Observed contract | Disposition |
| --- | --- | --- |
| `internal/anchors/registry.go` | Ordered anchors of kind Require, Forbid, or RequireInSection, evaluated per file with HTML comments stripped | DL1 to DL20 add rows in a new registry file |
| `internal/anchors/registry_data.go` | The `registry` variable appends every family; line 110 pins `- [ ] **red-capable**` in `bench-debug.md` | DL2 moves that anchor's file to `loop.md` |
| `internal/anchors/registry_ft311_preparation.go` | A family file is one exported-by-append slice with a header comment | The new family mirrors this file |
| `internal/conformance/registry_test.go` | `canaryFixtureFamilyRegistry` lists every anchor source file for `workflow-guidance-anchors` | DL-C1 adds the new registry file |
| `internal/conformance/fixture_bite_test.go` | Every retained fixture bites through its registered owner and restores | DL22 |
| `internal/conformance/prose_budget_test.go` | Line counts against the `projects/benchkit.md` table, exact row before glob row | DL21 |
| `internal/assessment/plan_validate.go` | A plan needs purpose, variable, tolerance, tasks, and conditions that differ only in the variable; kit-causal needs three named conditions | The template uses `descriptive` and `capability` |
| `internal/assessment/comparison.go` | A run binds by condition id, pinned revision, pinned lines per role, plan id, harness, limits, capabilities, and repetition | DL17 names every field |
| `internal/assessment/types.go` | A run carries `trial`, `condition`, `task_id`, `holdout`, `source`, and a `quality` map of measures | DL17 and DL19 |
| `internal/assessment/harness.go` | Only the Codex token-count fragment maps to usage events | DL19 keeps Claude counters unknown |
| `tests/canary/workflow-guidance-anchors/write-spec-conversation-fork` | `BASE` is the path, `EXPECT` is the diagnostic, `MUTATE.json` is one old-to-new edit | Every new fixture mirrors it |

## Line counts at phase entry

| file | lines | budget |
| --- | --- | --- |
| `.agents/commands/bench-implement-spec.md` | 80 | 80 |
| `.agents/commands/bench-write-spec.md` | 73 | 73 |
| `.agents/commands/bench-debug.md` | 169 | 170 |
| `.agents/commands/bench-review-implementation.md` | 194 | none |
| `.agents/skills/bench-craft-tdd/SKILL.md` | 122 | 122 |
| `.agents/skills/bench-craft-delegate/SKILL.md` | 124 | 126 |
| `.agents/skills/bench-craft-review/SKILL.md` | 122 | 122 |

## Readers not changed

- `docs/field-guide.html` describes the debug phase without a delegate route, so its text stays true.
- `.bench/BENCH.md` names `--delegate` because that flag changes authorship policy; `--debug-trial` changes none.
- `.agents/skills/bench-craft-review/references/finding-discipline.md` already requires a real run before a strong finding; the review arm adds the pasted run only under the flag.
