# Landing parse seams (FT169 research)

Source: read-only Terra delegate, 2026-09-01. Coordinator spot checks:
`internal/worktree/land_identity.go:195-203`,
`internal/spec/spec.go:175-179` — both resolve.

## 1. The source tip

`--source-tip` is required and refuses an empty value
(`internal/worktree/land.go:39`, `internal/usage/parse.go:148-150`).
`expandIdentity` expands a 4-to-39-character hex prefix, in either case, to
the full worktree head (`internal/worktree/land_identity.go:195-211`;
tests pin `tip[:4]`, `tip[:12]`, `tip[:39]` in
`internal/worktree/land_reauthorization_test.go:172-176`). A full
40-character lower-case sha passes unchanged.

A symbolic ref (`HEAD`, a branch, a label) is not expanded and refuses as
`worktree source tip mismatch` or
`assignment branch source tip mismatch`
(`internal/worktree/land_identity.go:105`, `:112`). An upper-case full sha
refuses the same way, because expansion skips a 40-character value.

No landing surface resolves a worktree tip for the caller.
`bench worktree list` prints no head column
(`internal/worktree/list.go:15`). Only `bench preflight review` prints a
`source` table whose tip is the checkout head
(`internal/preflight/command.go:78-84`).

## 2. The `--spec` flag

Two branches of the landing normalize differently. The tickets-only branch
calls `spec.LiveSpecSlug`, which accepts a bare slug, `foo.md`,
`specs/foo/spec.md`, or `specs/foo`
(`internal/worktree/land_identity.go:123`,
`internal/spec/spec.go:47-62`). The staged-spec branch passes the raw
argument to `spec.Resolve` (`internal/spec/spec.go:175-179`). That
function tries the raw argument as a filesystem path first, relative to
the process working directory. Only a slash-free slug reaches the
`specs/<slug>/spec.md` fallback rooted at the base.

A path form also
breaks the status lookup, which compares the raw argument to the
enumerated slug by equality (`internal/preflight/gather.go:242-250`).
The landing then refuses with `reviewed source range or ownership fence
is invalid: spec status not readable`
(`internal/worktree/land_identity.go:179-184`).

## 3. The `--base` flag

The landing runs three validations
(`internal/worktree/land.go:241-252`,
`internal/diff/range.go:149-165`):

- The base resolves and is an ancestor of the source tip.
- The base is the assignment's recorded start or a descendant of it.
- The base is an ancestor of the landing destination.

`bench preflight review` validates none of the three landing rules. Its
optional `--base` only checks resolvability and ancestry of the checkout
head (`internal/preflight/gather.go:51-58`,
`internal/preflight/decision.go:221-227`). So a review can pass green on a
base the landing later refuses. This is the exact 2026-08-29
`worktree-merge` occurrence, and it is open.

## Read and not read

The delegate read the land, usage, spec, landing-close, preflight, and
diff-range code. It did not read `land_resume.go` past the flag
declarations, the `landing` package composition code, `internal/intent`,
or the system tests.
