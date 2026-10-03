# Read skill frontmatter from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/skillsindex/skillsindex.go, internal/skillsindex/frontmatter_test.go (new)
Covers: MB59, MB60

## What to build

Make `FrontmatterField` in `internal/skillsindex` read the block reader's
frontmatter lines. It returns the first value of the key among those lines. A
frontmatter fault gives an empty value, as an unclosed opener does today. The
classifier read through `bounds.ClassifyNoFollow` stays first.

`internal/skillsindex/skillsindex.go` is over its line budget, so it does not
grow.

## Acceptance

- [ ] A CRLF skill file `---`, `index: x`, `---` gives the value `x`.
- [ ] `TestFrontmatterFieldRequiresCompleteLeadingFence` and `TestFrontmatterFieldReadsOnlyTheLeadingFence` stay green without an edit.
- [ ] `bench skills-index --check` is green on the live tree.
