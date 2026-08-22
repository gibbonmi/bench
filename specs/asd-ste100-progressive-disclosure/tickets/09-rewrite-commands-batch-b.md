# Rewrite command batch B: the maintenance and shaping phases

Blocked by: 08-rewrite-commands-batch-a.md
Writes: .agents/commands/bench-drain.md, .agents/commands/bench-setup-repo.md, .agents/commands/bench-shape-idea.md, .agents/commands/bench-deepen.md, .agents/commands/bench-assess.md, .agents/commands/bench-update-kit.md, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The six command files read in ASD-STE100. The delegate rewrites prose sentences only. These parts stay byte-identical: frontmatter, fenced blocks, the `Next:` token table, the `Entry orientation` and `Exit handoff` headings, every inline-code token, and every anchor needle. The roadmap-context query sentences in `bench-drain.md` survive verbatim. The orchestrator runs the token-set and line-class comparisons and reads every rule sentence.

The six rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the six rows removed (covers PD27).
- [ ] Every anchor needle on the batch's files survives (covers PD28).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in frontmatter, a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
