# Rewrite root docs batch A: the working agreement, glossary, and small docs

Blocked by: 01c-register-the-check-and-exclusions.md
Writes: AGENTS.md, CLAUDE.md, CONTEXT.md, .claude/README.md, .claude/output-styles/, DATA_HANDLING.md, SECURITY.md, projects/gl-axi.md, docs/greenfield-build-sequence.md, docs/release-runbook.md, docs/reporesident-distillation.md, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The ten files and the output-style file read in ASD-STE100. The delegate rewrites prose sentences only. These parts stay byte-identical:

- headings, fenced blocks, and tables
- the `DATA_HANDLING.md` passlist block and the `CLAUDE.md` import lines
- the `CONTEXT.md` term names and Avoid lists
- every inline-code token and every anchor needle

AGENTS.md keeps its pointer to the canonical `.bench/BENCH.md` and carries no shared-rule marker. The orchestrator runs the token-set and line-class comparisons and reads every rule sentence.

The ten rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the ten rows removed (covers PD27).
- [ ] The shared-rule and data-handling checks stay green (covers PD29).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in frontmatter, a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
