# Add red-capable tests for the release identity checks and two hook verdicts

Blocked by: none
Writes: internal/releasepreflight/identity_test.go (new), internal/harness/worktree_test.go, internal/gittest/gittest.go, internal/gittest/gittest_test.go
Covers: none

## What to build

The release identity gate refuses a release whose identity disagrees. Tests in
`internal/releasepreflight` name `checkIdentity`, `checkAncestry`, and
`checkChangelog`. A fixture repository holds a tagged commit, a package version,
and a `go.mod` toolchain line. Each case changes one fact from a green baseline,
so only that fact can cause the refusal.

`WorktreeCommand` has two more verdict tests. One test makes `worktree.Create`
fail after the event passes validation: a repository with no commit gives no
start for the worktree. One test makes the worktree registration unavailable for
a remove event: a symlinked `worktrees` admin directory makes the registration
scan refuse.

Each new test names, in its comment, the production mutation that turns it red.

The identity fixture runs git through `gittest.Output`, which the `gittest`
package exports for this ticket. The remove verdict test takes the expected
refusal text from `git.ScanWorktreeAdmin`, so the test holds no second copy of
that text.

## Acceptance

- [x] `checkIdentity` accepts the baseline. It refuses a ref that is not an
      exact tag and a tag that does not resolve to `HEAD`. It also refuses a
      package version or a binary version that disagrees.
- [x] `checkAncestry` accepts a `HEAD` that `origin/main` contains and refuses
      a `HEAD` that it does not contain.
- [x] `checkChangelog` accepts one dated release heading. It refuses a
      changelog with no matching heading or with two matching headings.
- [x] A `WorktreeCommand` create event in a repository with no commit exits 1
      with the start-resolution verdict and an empty stdout.
- [x] A `WorktreeCommand` remove event in a repository whose worktree admin
      directory is a symlink exits 1 with the registration-unavailable verdict.
- [x] Each named mutation turns its test red; the ticket author records each red.

## Verification

Each probe applied one mutation to the production file, ran only the named test,
and restored the file. Each probe was red.

| Mutation | Test | Result |
|---|---|---|
| Delete the `exactTag` refusal | `TestCheckIdentityRefusesADisagreement/ref_is_not_an_exact_tag` | red |
| Drop the `tagCommit != commit` comparison | `TestCheckIdentityRefusesADisagreement/tag_does_not_resolve_to_HEAD` | red |
| Drop the `pkg != version` comparison | `TestCheckIdentityRefusesADisagreement/package_version_disagrees` | red |
| Drop the `r.binaryVersion != version` comparison | `TestCheckIdentityRefusesADisagreement/binary_version_disagrees` | red |
| Swap the merge-base operands | `TestCheckAncestryRefusesAHeadThatOriginMainDoesNotContain` | red |
| Delete the `len(matches) != 1` refusal | `TestCheckChangelogRefusesAMissingReleaseHeading` | red |
| Weaken `len(matches) != 1` to `len(matches) == 0` | `TestCheckChangelogRefusesADuplicateReleaseHeading` | red |
| Delete the `worktree.Create` error return | `TestWorktreeCommandCreateRefusesACreateError` | red |
| Delete the registration guard | `TestWorktreeCommandRemoveRefusesAnUnavailableRegistration` | red |
