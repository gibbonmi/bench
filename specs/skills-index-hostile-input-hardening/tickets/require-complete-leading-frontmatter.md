# Require complete leading frontmatter

Blocked by: harden-every-skill-file-reader.md
Writes: internal/skillsindex/skillsindex.go, internal/skillsindex/skillsindex_test.go

## What to build

Parse only classified SKILL bytes whose frontmatter begins at byte zero with an opener
line and contains a closing `---` line. Preserve the first exact-key value, duplicate
key posture, and no-trailing-newline case inside a valid fence. Publish only a parsed
field value; the dependent sink ticket owns whether that value is safe to render.

## Acceptance

- [ ] `(covers HI5)` Late, unclosed, valid, duplicate-key, and no-trailing-newline
  fixtures pin the complete leading-fence and first-value contract.

