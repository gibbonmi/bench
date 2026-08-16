# Require exactly one marker span

Blocked by: harden-every-reference-file-reader.md
Writes: internal/skillsindex/skillsindex.go, internal/skillsindex/skillsindex_test.go

## What to build

Have the skills-index marker parser scan the complete classified reference and return
a span only for exactly one start marker followed by exactly one end marker. Check and
Write consume the same result, and Write preserves bytes for zero, reversed, unclosed,
duplicate-start, duplicate-end, or two-complete-span documents.

The cardinality table and both consumers form one smallest tracer: landing only the
parser or only one consumer leaves the other able to select the first plausible span,
which strands HI7 red at the public check/write seam.

## Acceptance

- [ ] `(covers HI7)` The full zero/reversed/unclosed/duplicated/two-span table refuses,
  while exactly one ordered span succeeds.
- [ ] `(covers HI7)` Check cannot pass by selecting the first span and Write changes no
  malformed fixture bytes.

