# Rewrite command batch A: the workflow phases

Blocked by: 07-rewrite-skills-batch-5.md
Writes: .agents/commands/bench.md, .agents/commands/bench-write-spec.md, .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, .agents/commands/bench-final-check.md, .agents/commands/bench-debug.md, .agents/commands/bench-what-next.md, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The seven command files read in ASD-STE100. The delegate rewrites prose sentences only. These parts stay byte-identical: frontmatter, fenced blocks, tables, the `Entry orientation` and `Exit handoff` headings, the retro template, every inline-code token, and every anchor needle. `bench-write-spec.md` stays within 73 lines, `bench-implement-spec.md` within 75, and `bench-debug.md` within 170; a rewrite that cannot exits and reports. The structured-phase sentences, the review convergence sentences, and the two `Bootstrap authority before execution` mentions survive verbatim. The orchestrator runs the token-set and line-class comparisons and reads every rule sentence.

The seven rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the seven rows removed (covers PD27).
- [ ] Every anchor needle on the batch's files survives (covers PD28).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in frontmatter, a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
