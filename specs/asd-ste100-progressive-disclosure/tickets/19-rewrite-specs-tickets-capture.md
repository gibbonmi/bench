# Rewrite the other spec, the light-path tickets, and the capture files

Blocked by: 01c-register-the-check-and-exclusions.md
Writes: specs/inherited-toolchain-environment/, tickets/, capture/agent-performance/, capture/audits/, capture/FIXES.md, capture/parallel-session-friction.md, capture/learnings.md, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The light-path ticket, the agent-performance scorecards, the capture audits and notes, and the learnings preamble read in ASD-STE100. The delegate rewrites prose sentences only. These parts stay intact: headings, tables, fenced blocks, the `Feeds:` and journal markers, every inline-code token, and every coverage-row id. `specs/inherited-toolchain-environment/` is rewritten when its folder exists; when the folder is gone, its row is stale and leaves anyway. The orchestrator runs the token-set and line-class comparisons and reads every sentence.

All rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the batch's rows removed (covers PD27).
- [ ] `bench coverage --check` stays green on every spec present (covers PD29).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in frontmatter, a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
