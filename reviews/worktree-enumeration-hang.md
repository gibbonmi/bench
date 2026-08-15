## Standards

Finding count: 1. Worst: duplicated lifecycle control-record name.

- `auto-fix` — `internal/worktree/worktree.go` and `internal/git/git.go` each
  define `bench-lease`. Move the production filename to one dependency-safe
  owner and consume it from both sites; retain only an independent test literal.

## Spec

Finding count: 1. Worst: bounded rev-parse failures omit Git's diagnostic.

- `auto-fix` — The typed resolution refusal must retain rev-parse failure text.
  `bounds.RunOutput` captures stderr, but `boundedGit` drops it for an exit
  failure. Preserve that diagnostic in the typed error and add a non-silent
  failure fixture.

## Coverage

Finding count: 1. Worst: successful corrupt common-dir variants are unproven.

- `auto-fix` — Extend the existing WE19 resolution seam for empty stdout and a
  symlink-to-directory common-dir result. Both must refuse before porcelain,
  proving `validateCommonDir` does not admit a corrupt successful resolution.
