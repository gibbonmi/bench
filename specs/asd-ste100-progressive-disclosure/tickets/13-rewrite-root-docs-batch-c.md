# Rewrite root docs batch C: the assessments and the ADRs

Blocked by: 01c-register-the-check-and-exclusions.md
Writes: ASSESSMENT.md, skills-assessment.md, docs/adr/, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The two assessments and the thirteen ADRs read in ASD-STE100. The delegate rewrites prose sentences only and leaves headings, fenced blocks, tables, every inline-code token, and every decision statement's meaning byte-for-meaning intact. An ADR records the decided state; the rewrite adds no history and no path. The orchestrator runs the token-set and line-class comparisons and reads every decision sentence. The three rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the three rows removed (covers PD27).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in frontmatter, a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
