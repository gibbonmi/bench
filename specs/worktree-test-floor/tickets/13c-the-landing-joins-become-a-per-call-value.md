# The landing joins become a per-call value

Line: opus / high.

The change moves six landing joins, the lifecycle gaps, and the live-binary
resolvers from package-level variables into a value the verb builds at its
boundary. It is a foundational Go-seam rewrite of the landing transaction, so
the top of the mid tier applies.

Blocked by: 13b-a-stubbing-test-is-serial.md
Writes: internal/worktree/land.go, internal/worktree/land_resume.go, internal/worktree/lifecycle.go, internal/worktree/ownership.go, internal/worktree/live_binary.go, internal/worktree/parallel_census_test.go, and every `internal/worktree/*_test.go` that assigns a package-level variable

## What to build

Compose the `Fault` precedent in `ownership.go`. The landing's joins —
`landReviewed`, `advanceLandingMarker`, `reconcileLanding`,
`releaseLandingAssignment`, `authorizeLandingSource`, and
`cleanupTransactionBoundary` — move into one value with a default instance.
`LandCommand` and `ResumeLandCommand` build it at their boundary, and an
internal form takes it. The lifecycle gaps and the live-binary resolvers
follow the same shape.

A test passes its own value to the internal form. It assigns no package-level
variable. Every such test then calls `t.Parallel()`. The census rule of ticket
13b turns from a serial edge into a refusal. A test that assigns a
package-level identifier of a non-test file is reported.

## Acceptance

- [ ] No `internal/worktree` test assigns a package-level identifier of a non-test file, and the census refuses one that does.
- [ ] WF01 stays green on the live tree, and the tests ticket 13b left serial call `t.Parallel()`.
- [ ] `go test -race -count=1 -parallel 2 ./internal/worktree` passes.
- [ ] `bench worktree land --help` and the landing refusal bytes do not change.
