# Rewrite the core and the reference in ASD-STE100

Blocked by: 02-split-bench-core-from-reference.md
Writes: .bench/BENCH.md, .bench/BENCH-reference.md, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

Both files read in ASD-STE100 and pass the prose mechanics check. The kept sections of the core carry sentences over 25 words. The invariant items and the predicate paragraph carry seven and eight sentences. The reference carries long sentences and two long paragraphs.

The delegate rewrites prose sentences only and splits the long items into paragraphs of six sentences or fewer. These parts stay byte-identical: headings, tables, fenced blocks, the skills-index block between its markers, every inline-code token, every shared-rule marker, and every anchor needle. The structured-phase contract block keeps its declaration line and its four clause lines. The orchestrator runs the token-set and line-class comparisons and reads every rule sentence.

The two `.bench/` rows leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] Each invariant item and the predicate paragraph hold six sentences or fewer (covers PD8).
- [ ] The live-tree test passes with the two `.bench/` rows removed (covers PD27).
- [ ] Every anchor needle on both files survives (covers PD28).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
