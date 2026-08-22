# Rewrite skill batch 4: craft-review, craft-synthesis

Blocked by: 05-rewrite-skills-batch-3.md
Writes: .agents/skills/bench-craft-review/, .agents/skills/bench-craft-synthesis/, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The two skills and the smell-baseline reference read in ASD-STE100. The delegate rewrites prose sentences only and leaves frontmatter, fenced examples, tables, headings, every inline-code token, and every anchor needle byte-identical. Each skill stays within 120 lines. The orchestrator runs the token-set and line-class comparisons and reads every rule sentence. The two rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the two rows removed (covers PD27).
- [ ] Every anchor needle on the batch's files survives (covers PD28).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in frontmatter, a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
