# landing-refusal-diagnostics review pickup

Frozen base: `850fd2677c4fd56ff2a0e8f2f1c5ef1698ba1148`

Reviewed tip: `03612a7b93c19c3efb71d02963304af85d2ae90f`

## Standards

Finding count: 1. Worst issue: P3 long parameter list/data clump.

- **ask-user — P3 long parameter list/data clump.**
  `assignmentForRequest` takes six bare strings, while the release caller passes
  empty identity placeholders and the land caller passes the base/tip pair. The
  optional recovery context is therefore implicit at both call sites. The
  review smell baseline says multiplying arguments should become a type before
  they travel together (`.agents/skills/bench-craft-review/references/smell-baseline.md`,
  Long Parameter List and Data Clumps); the candidate is at
  `internal/worktree/worktree.go:432`, with callers in
  `internal/worktree/ownership.go:354` and `internal/worktree/land.go:408`.
  Whether a small recovery-context type improves this two-caller helper enough
  to justify another abstraction is reviewer judgment.

## Spec

Finding count: 1. Worst issue: caller-token disclosure in a `next=` value.

- **auto-fix — retained-release recovery echoes the caller token.**
  The approved `next=` contract says, “A caller token is never echoed or
  persisted” (`specs/landing-refusal-diagnostics/spec.md:163`).
  `releaseNext` instead shell-quotes every line-safe caller token into the
  emitted command (`internal/worktree/ownership.go:431`), and
  `TestReleaseCommandRefusalListsBoundedIgnoredPathsWithTrueTotal` pins the
  disclosed token (`internal/worktree/land_test.go:1421`). Always render the
  request placeholder in this continuation and update the exact-output tests;
  preserve the unsafe-path assignment-pointer form.

## Coverage

Finding count: 0. Worst issue: none.

The review enumerated identity boundaries, all first-run and resume follow-up
states, runtime and foreign residue, bounded tables, output channels, hostile
paths and control bytes, and zero/one/multiple-assignment recovery states.
Existing tests cover every mapped LR1–LR23 behavior and the specified edge
inventory. Coordinator-focused tests for the representative land, release,
runtime-fingerprint, and recovery paths passed.
