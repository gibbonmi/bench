# Add the prose engine: parser, walk, exclusion grammar, classification

Blocked by: none
Writes: internal/prose/ (new)
Line: `opus` / medium under Claude Code, `gpt-5.6-terra` / medium under Codex — the oracle's correctness outranks speed, and the parser follows the prose-budget precedent.

## What to build

`internal/prose` is a new deep package. It owns the document parser, the tree walk, the exclusion-file grammar, and the subject classification through `bounds.ClassifyNoFollow`. It returns diagnostics for one root. The parser order, the boundary rules, the field-line rule, the skip list, the no-subject guard, and the fail-closed states are the ones the spec fixes. Nothing registers the package yet, so the tree stays green.

Its own comments are authored in ASD-STE100 under `craft-comments`. Table-driven unit tests feed one document or one temporary root and assert the findings.

## Acceptance

- [ ] A 26-word sentence returns one finding with its line and count, and a 25-word sentence returns none (covers PD9, PD10).
- [ ] A seven-sentence paragraph returns one finding, and a six-sentence paragraph returns none (covers PD11, PD12).
- [ ] Code spans, fenced blocks, table rows, headings, frontmatter, HTML comments, and link targets return no finding (covers PD13).
- [ ] An indented block after a blank line is skipped, and a list continuation after a non-blank line is graded (covers PD14).
- [ ] A seven-sentence list item reds, two adjacent four-sentence items do not, and a 25-word ordered-list item does not (covers PD15).
- [ ] In-token periods, all five abbreviations, and an ellipsis do not split, and `word. Next` does (covers PD16).
- [ ] A no-break space splits a word and a zero-width space does not (covers PD17).
- [ ] A link, a FIFO, a non-UTF-8 file, and an over-bound file each red, and a file at the bound is graded (covers PD20).
- [ ] An unterminated fence, HTML comment, or frontmatter block reds naming the opening line (covers PD21).
- [ ] Files under `tests/canary/`, `node_modules/`, `dist/`, `.git/`, and `testdata/` are not graded (covers PD22).
- [ ] A `.go` file and a `.sh` file with a long comment sentence are not graded (covers PD43).
- [ ] The four named field lines and a four-word label line are not graded, and an `Occurrence:` line and a five-word label line are (covers PD41).
- [ ] Twenty contiguous one-sentence `Occurrence:` lines form twenty paragraphs and return no finding (covers PD41).
