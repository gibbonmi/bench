# Description inventory of the listed skills and phases

Consumed by: the guidance-weight map (#13) and the FT100 always-loaded spec.
Drift: an edit of a `description:` line in `.agents/skills/` or `.agents/commands/`.
Retire when: the FT100 always-loaded spec lands.

## Recommendation, scope, and evidence status

Recommendation: fix the four description defects now. Trim the content phrases only after an FT231 trial measures trigger accuracy. The content phrases are the only text that the harness can match at selection time.

Scope: the 26 kit descriptions that Claude Code listed on 2026-10-11, 6,188 bytes in total. Slice D1 holds nine craft skills, slice D2 holds eight craft skills and `prototype`, and slice D3 holds eight phases.

Evidence status: all three reader pairs agree on every quoted value, count, and body overlap. The coordinator checked the standard's text at `.agents/skills/bench-craft-skills/SKILL.md:44-45` and the two term defects.

## Facts

| Slice | Description | Characters | Trigger phrase | Content restates the body |
|---|---|---|---|---|
| D1 | craft-adr | 229 | yes | yes |
| D1 | craft-cli | 245 | yes | yes |
| D1 | craft-comments | 225 | yes | yes |
| D1 | craft-delegate | 241 | yes | yes |
| D1 | craft-design-system | 179 | yes | yes |
| D1 | craft-domain | 245 | yes | yes |
| D1 | craft-gate | 240 | yes | yes |
| D1 | craft-grill | 244 | yes | yes |
| D1 | craft-line | 234 | yes | yes |
| D2 | craft-research | 236 | yes | yes |
| D2 | craft-review | 250 | yes | yes |
| D2 | craft-seams | 236 | yes | yes |
| D2 | craft-skills | 184 | yes | yes |
| D2 | craft-spec | 250 | yes | yes |
| D2 | craft-synthesis | 247 | yes | yes |
| D2 | craft-tdd | 250 | yes | yes |
| D2 | craft-tickets | 244 | yes | yes |
| D2 | prototype | 239 | yes | yes |
| D3 | bench | 77 | no | yes |
| D3 | bench-debug | 249 | yes | yes |
| D3 | bench-drain | 250 | no | yes |
| D3 | bench-final-check | 237 | partial | yes |
| D3 | bench-implement-spec | 247 | yes | yes |
| D3 | bench-review-implementation | 156 | yes | yes |
| D3 | bench-shape-idea | 167 | yes | yes |
| D3 | bench-write-spec | 183 | yes | yes |

Every description has content phrases that restate a heading or a line of its own body. The standard says that a model-invoked description sits in the context window on every turn (`.agents/skills/bench-craft-skills/SKILL.md:17-19`).

The triad rule targets a triad that is spelled out three times (`.agents/skills/bench-craft-skills/SKILL.md:44-45`). The research charge said "a triad spelled out in full", so the Sonnet triad findings apply the rule too widely. Two lists meet the literal rule. "Standards, Spec, and Coverage" appears five times, at `.agents/commands/bench-review-implementation.md:2`, `:28`, `:43`, `:67` and `.bench/BENCH.md:120`. The `bench-debug` symptom list appears three times, at `.agents/commands/bench-debug.md:2`, `:9` and `.agents/skills/bench-debug/SKILL.md:3`.

Each Codex adapter except `bench-debug` has its own explicit-invoke description, for example `.agents/skills/bench-drain/SKILL.md:3`. So a command description is the Codex trigger only for `bench-debug`.

## Defects

- `.agents/commands/bench.md:2` and `.agents/commands/bench-drain.md:2` have no trigger phrase.
- `.agents/commands/bench-drain.md:2` ends "Maintenance, not a phase." The body says "This is the single roadmap-maintenance phase" at `:9`. At `:15` it quotes a "not a workflow phase" label, which no line of the file holds.
- `.agents/commands/bench-implement-spec.md:2` and `.agents/commands/bench-write-spec.md:2` name "the lighter-path threshold in .bench/BENCH.md". `.bench/BENCH.md` uses the words "lighter path" and "Right-size", and never "threshold".
- `.bench/BENCH-reference.md:123-125` says that most phases are reviewer-chosen and that a model does not start them unbidden. The policy table at `internal/conformance/skills_index_checks_test.go:216-229` marks all eight listed phases as Claude model-invocable.

## Inferences

A trim of the content phrases changes the only text that the harness can match before a skill loads. That is a behavior change, not a repeated copy, so it fits the FT231 rule of #2. The four defects are corrections of stale or missing text, so they can go now.

A trim to the trigger and one leading-word phrase would remove about 1,700 bytes, or about 430 tokens on every turn.

## Unknowns

- Whether the Claude harness matches on the content phrases or only on the trigger phrases is unknown. No file in the tree states it.

## Verification record

- [x] The output opens with the recommendation, the scope, and the evidence status.
- [x] Facts, defects, and inferences sit in separate sections.
- [x] A table holds one row for each description.
- [x] The unknowns are stated.
- [x] A validation plan ends the output.

## Validation plan

1. After the defect fixes, run `bench test --check skill-description-budgets` and the skills-index conformance check.
2. Before any content trim, run an FT231 trial that compares skill selection with the current descriptions and with the trimmed ones.
