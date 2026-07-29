# Build the package-scoped shared artifact set

Blocked by: none

## What to build

The artifact contract package lazily builds one package-owned, host-narrowed
artifact set; protects it with read-only and digest belts; fails closed after
staging failure; cleans it up from `TestMain`; and routes all six ruled sharers
through it without changing their existing assertions. The contract helper
stages a release plan under a caller-owned directory through one shared
implementation.

## Acceptance

- [ ] All six ruled sharers pass against one shared generation and assert the build count is exactly one.
- [ ] Writes to one promoted tarball and one staged-source worktree file fail with a permission error.
- [ ] Digest verification attributes mutation as `shared artifact set mutated`.
- [ ] A fabricated-fixture-only `go test -run` selection performs no artifact build.
- [ ] A requester after failed staging reports `shared artifact set staging failed in an earlier test`.
- [ ] Existing `NarrowReleasePlan` callers retain their behavior through the caller-owned-directory variant.
- [ ] The package-owned read-only directory is made writable and removed after the package run, including on test failure.
