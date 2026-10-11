# Read fence entries only from the list items of the fence section

Blocked by: none
Writes: internal/spec/fences.go, internal/spec/fences_test.go, internal/preflight/fence_writes_test.go
Covers: none

## What to build

The `## Ownership fences` section of a spec authorizes the paths that a build can write. Build preflight and the review-pickup grade read that section through `FenceTokens` in `internal/spec/fences.go`. The parser reads every backticked token outside parentheses in the section, and it also reads tokens in prose sentences. Thus a prose sentence such as "`.bench/BENCH.md` governs execution-plan changes." authorizes a path that no ticket writes.

The parser reads backticked tokens only from list items of the section and from their indented continuation lines. A list item is a line that starts with `- ` or `* ` after optional indentation. A prose paragraph in the section authorizes nothing. The parenthesis shield stays: a token inside parentheses authorizes nothing, also when the parenthetical spans lines.

This rule supersedes the behavior text of SC20 in the retired slicing-closure spec, which said that preflight reports a backticked prose token as fence-only. The SC20 test keeps its purpose: the fence-writes row of `bench preflight build` reads the same parser as the preflight facts. Its fixture now puts the unwritten token in a list item.

## Acceptance

- [ ] A backticked path in a prose line of the fence section is not a fence entry.
- [ ] A backticked path in a list item of the fence section is a fence entry.
- [ ] A backticked path on an indented continuation line of a list item is a fence entry.
- [ ] A backticked path inside parentheses is not a fence entry, also when the parenthetical spans a list item and its continuation line.
- [ ] A backticked token in a prose line of the fence section produces no fence-only row in `bench preflight build`.
