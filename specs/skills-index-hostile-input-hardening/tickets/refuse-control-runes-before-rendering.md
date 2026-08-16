# Refuse control runes before rendering

Blocked by: require-complete-leading-frontmatter.md
Writes: internal/skillsindex/skillsindex.go, internal/skillsindex/skillsindex_test.go

## What to build

Before a field published by the complete-frontmatter parser reaches `Entry.Line`,
reject every Unicode control rune in both `index:` and `index-note:` while preserving
ordinary graphic Unicode. The observable contract is the one-line sink, not a helper
counter or an ASCII-only byte list.

## Acceptance

- [ ] `(covers HI6)` Tab, CR, ESC, BEL, NUL, DEL, and every other control rune are
  refused before rendering; CR cannot forge a second line and graphic Unicode remains
  one rendered line.

