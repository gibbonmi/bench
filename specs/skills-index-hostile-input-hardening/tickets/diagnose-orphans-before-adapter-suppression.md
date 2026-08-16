# Diagnose orphans before adapter suppression

Blocked by: enumerate-skill-directories-literally.md
Writes: internal/skillsindex/skillsindex.go, internal/skillsindex/skillsindex_test.go

## What to build

Turn literal discovery's missing-file state into an attributed orphan diagnostic before
checking for a same-named command adapter. A command-named directory with a valid
`SKILL.md` is classified and then suppressed; its paired missing-file directory is
diagnosed and blocks Write. Empty regular bytes continue through frontmatter grading
and remain distinct from absence.

## Acceptance

- [ ] `(covers HI3)` The paired command-named orphan/valid-skill fixtures prove missing
  classification precedes suppression, and missing versus empty retain distinct
  diagnostics and write behavior.

