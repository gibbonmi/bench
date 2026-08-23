# Rewrite root docs batch B: README and the kit profile

Blocked by: 01c-register-the-check-and-exclusions.md
Writes: README.md, projects/benchkit.md, .bench/prose-exclusions (new in ticket 01c)
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

`README.md` and `projects/benchkit.md` read in ASD-STE100. The delegate rewrites prose sentences only and leaves headings, fenced blocks, the layout listing, every table, every inline-code token, and every anchor needle byte-identical. README keeps `## Reviewer quick start` as its first H2 and carries no shared-rule marker. The profile keeps its hostile-input heading, its budget table, its conformance table, and its lines table untouched. The orchestrator runs the token-set and line-class comparisons and reads every rule sentence. The two rows for this batch leave `.bench/prose-exclusions` in this commit.

## Acceptance

- [ ] The live-tree test passes with the two rows removed (covers PD27).
- [ ] The shared-rule, prose-budget, and command-first checks stay green (covers PD29).
- [ ] The inline-code token multiset of each file equals the original's (covers PD30).
- [ ] No changed line lies in frontmatter, a fence, a table, or a marker (covers PD31).
- [ ] No row was added to `.bench/prose-exclusions` (covers PD32).
