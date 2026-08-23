# Rewrite the thin bench-* adapter skills

Blocked by: 09-rewrite-commands-batch-b.md
Writes: .agents/skills/bench/, .agents/skills/bench-debug/, .agents/skills/bench-assess/, .agents/skills/bench-deepen/, .agents/skills/bench-drain/, .agents/skills/bench-final-check/, .agents/skills/bench-implement-spec/, .agents/skills/bench-review-implementation/, .agents/skills/bench-setup-repo/, .agents/skills/bench-shape-idea/, .agents/skills/bench-update-kit/, .agents/skills/bench-what-next/, .agents/skills/bench-write-spec/, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The thirteen adapter skill bodies, the `bench-debug` loop-construction reference, and the `bench-deepen` report intro read in ASD-STE100. The delegate rewrites prose sentences only and leaves frontmatter, `agents/openai.yaml`, fenced blocks, the HTML scaffold, every inline-code token, and every anchor needle byte-identical. The orchestrator runs the token-set and line-class comparisons and reads every rule sentence. The thirteen rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the thirteen rows removed (covers PD27).
- [ ] Every anchor needle on the batch's files survives (covers PD28).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in frontmatter, a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
