# Rewrite the roadmap bodies

Blocked by: 01c-register-the-check-and-exclusions.md
Writes: ROADMAP.md, roadmap/, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

`ROADMAP.md` and every `roadmap/FT*.md` read in ASD-STE100. The delegate rewrites prose sentences only, file by file. These parts stay byte-identical:

- each row's heading line, its `Next:` line, and any `Occurrences:` ledger line
- the dependency tables and the `## Recommended sequence` lines
- every inline-code token and every roadmap id

An `Occurrence:` line may split into two sentences on one physical line. The FT100 and FT179 remaining-scope paragraphs belong to ticket 29. The orchestrator runs the token-set and line-class comparisons and reads every row. The two rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the two rows removed (covers PD27).
- [ ] The row-grammar and roadmap-detail checks stay green (covers PD29).
- [ ] Every row keeps its heading line, `Next:` line, and ledger byte-identical (covers PD39).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in a table or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
