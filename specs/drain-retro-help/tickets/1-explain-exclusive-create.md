# Explain exclusive retrospective creation

Blocked by: none
Writes: internal/roadmap/retro.go
Covers: none

## What to build

Add one sentence to the detailed `bench retro` help. State that `--body` creates the retrospective only once. State that a refusal preserves an existing retrospective.

## Acceptance

- [ ] `bench retro --help` says that `--body` creates only once and preserves an existing retrospective when it refuses the write.
- [ ] `bench help` keeps its concise `retro` inventory without the detailed creation sentence.
- [ ] A repeated `--body` call refuses and preserves the body from the first call.
