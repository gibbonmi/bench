# Record the drain's user-facing fixes in the changelog

Blocked by: none
Writes: CHANGELOG.md
Covers: none

## What to build

The 2026-10-08 drain d-4df49992a158 landed four light-path fixes. Two of them change what a
user sees: `bench gate-prose` marks an excluded path `excluded`, and `bench release --help`
and `bench upgrade --help` print their usage lines with exit 0. Those fixes landed without a
changelog entry, so that four tickets did not collide on one file.

Add one typed entry under the Unreleased section of `CHANGELOG.md` that states both
changes, including the repeated-flag refusal that `bench upgrade` now makes.

## Acceptance

- [ ] `CHANGELOG.md` holds one Unreleased entry that names the `excluded` verdict of `bench gate-prose`.
- [ ] The same entry names the `--help` behavior of `bench release` and `bench upgrade`, and the repeated-flag refusal of `bench upgrade`.
