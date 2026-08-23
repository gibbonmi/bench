# Rewrite skill batch 1: ste-prose, craft-skills, craft-comments

Blocked by: 01b-pin-thresholds-terms-and-profile.md
Writes: .agents/skills/bench-craft-spec/references/ste-prose.md, .agents/skills/bench-craft-skills/SKILL.md, .agents/skills/bench-craft-comments/SKILL.md, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The three files read in ASD-STE100 per `ste-prose.md`. The delegate rewrites prose sentences only. Frontmatter, fenced examples, tables, headings, the contrastive pairs, every inline-code token, and every anchor needle stay byte-identical. `ste-prose.md` keeps the two threshold sentences the registry rows pin. The orchestrator runs the token-set and line-class comparisons and reads every rule sentence before the batch lands. The three rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the three rows removed (covers PD27).
- [ ] Every anchor needle on the three files survives (covers PD28).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in frontmatter, a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
