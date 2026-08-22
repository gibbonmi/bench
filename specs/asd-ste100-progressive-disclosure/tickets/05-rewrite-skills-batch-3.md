# Rewrite skill batch 3: craft-tickets, craft-tdd, craft-delegate, craft-grill, craft-adr

Blocked by: 04-rewrite-skills-batch-2.md
Writes: .agents/skills/bench-craft-tickets/, .agents/skills/bench-craft-tdd/, .agents/skills/bench-craft-delegate/, .agents/skills/bench-craft-grill/, .agents/skills/bench-craft-adr/, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The five skills and their references read in ASD-STE100. The delegate rewrites prose sentences only and leaves frontmatter, fenced examples, the ticket-example block, tables, headings, every inline-code token, and every anchor needle byte-identical. `craft-tickets` stays within 100 lines and each other skill within 120 lines. The cross-harness reviewer recipes keep their pinned command lines. The orchestrator runs the token-set and line-class comparisons and reads every rule sentence. The five rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the five rows removed (covers PD27).
- [ ] Every anchor needle on the batch's files survives (covers PD28).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in frontmatter, a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
