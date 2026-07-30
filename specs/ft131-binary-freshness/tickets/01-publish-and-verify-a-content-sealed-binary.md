# Publish and verify a content-sealed binary

Blocked by: none

## What to build

The Go freshness owner discovers the Bench command's complete local build-input
set, validates the executable and companion seal without following untrusted
paths, and gives the build owner an atomic publication path. The built CLI
exposes that owner through `freshness-check`.

## Acceptance

- [x] Successful builds publish an executable and matching content seal; failed
  or interrupted publication never produces a trusted pair.
- [x] Missing, unreadable, malformed, partial, non-regular, and symlinked
  executable or seal artifacts fail closed with the prescribed rebuild action.
- [x] Equal, newer, and older mtimes do not affect a content-seal verdict.
- [x] The discovered input set includes every regular local Go build input plus
  build-script/version/flag inputs, including resolved generated or untracked
  files, while unrelated files do not change the seal.
- [x] Root and nested-cwd checks agree, repeated unchanged checks stay green,
  and checking does not mutate tracked configuration or self-invalidate.
