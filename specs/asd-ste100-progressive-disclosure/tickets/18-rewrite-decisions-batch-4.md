# Rewrite decisions batch 4: the remaining working maps and assets

Blocked by: 01c-register-the-check-and-exclusions.md
Writes: decisions/craft-research.md, decisions/assets/craft-research-research.md, decisions/cost-follows-project-size.md, decisions/diff-visual.md, decisions/spec-build-review-gate-cadence.md, decisions/worktree-orphan-retirement.md, decisions/ft144-post-approval-edits.md, decisions/assets/ft171-shared-fixture-staged-binary.md, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

The eight files read in ASD-STE100. The delegate rewrites prose sentences only. These parts stay intact: the decision-ticket grammar, `Blocked by` lines, state labels, headings, tables, fenced blocks, every inline-code token, and every recorded reviewer choice. The orchestrator runs the token-set and line-class comparisons and reads every decision sentence.

The eight rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the eight rows removed (covers PD27).
- [ ] The decision-map integrity check stays green (covers PD29).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
