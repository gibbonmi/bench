# Repair the prose engine findings S1, S3, C1, C2, C3

Blocked by: 29-close-the-chain.md
Writes: internal/prose/, internal/conformance/prose_mechanics_test.go (advisory)
Line: `opus` / medium under Claude Code, `gpt-5.6-terra` / medium under Codex — oracle code at a known seam.

## What to build

S1: `internal/prose` exports one row-reading seam, and `prose_mechanics_test.go` reads the live rows through it. The approved set stays an independent literal. S3: the package doc names the two bounds and drops the numerals. C1: the parser folds code spans before the label test, so a colon inside a span never makes a field line. C2: a row that names a directory without a trailing `/` reds as malformed. C3: a double-backtick span counts as one token.

## Acceptance

- [ ] One exclusion-row parser exists, and the subset test still reds on an added row (covers PD32).
- [ ] A 40-word line with a colon inside a code span returns one finding. A label line inside a paragraph run splits that run (covers PD41).
- [ ] A directory row with no trailing `/` reds with its own message (covers PD19).
- [ ] A double-backtick span with an inner backtick counts as one word (covers PD13).
