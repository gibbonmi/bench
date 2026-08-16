# Clean up failed replacements

Blocked by: harden-every-reference-file-reader.md
Writes: internal/skillsindex/skillsindex.go, internal/skillsindex/skillsindex_test.go

## What to build

Give replacement cleanup authority from successful sibling temp creation until a
successful rename. Defer removal across write, close, chmod, and rename errors, and
disarm it only after rename. Expose the production replacement-operation seam needed
for an injected post-creation rename failure; the dependent SIGINT ticket extends the
same lifetime with cancellation rather than copying cleanup.

## Acceptance

- [ ] `(covers HI9)` An injected rename failure after temp creation returns the error,
  preserves original bytes, and leaves no `.bench/.skills-index-*` residue.

