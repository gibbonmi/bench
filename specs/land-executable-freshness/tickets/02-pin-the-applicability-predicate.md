# Pin the applicability predicate

Blocked by: 01-refuse-stale-landing-executable
Writes: internal/freshness/freshness.go (advisory: no production change expected),
internal/freshness/freshness_test.go, specs/land-executable-freshness/spec.md

## What to build

Review's Coverage axis found the feature's applicability gate ungraded.
`freshness.DeclaresBuildInputs` is reached only through the land surface, which
exercises a plain absent manifest and a plain present one. The spec's edge
inventory promises more than that — a dangling or live symlink at the manifest
path reports presence and routes to the owner, never an authoritative absence —
and it cites the owner's symlink test for it. That citation does not hold: the
owner's `TestVerifyRefusesLiveAndDanglingSymlinkArtifacts` covers the seal and
the executable, never the manifest path.

Add a direct table test over `DeclaresBuildInputs` and correct the edge
inventory's citation to name it.

## Acceptance

- [ ] A table test grades `DeclaresBuildInputs` over an absent manifest, a
      regular manifest, an empty manifest, a live symlink, a dangling symlink,
      and a directory at the manifest path.
- [ ] Swapping `Lstat` for `Stat` turns the dangling-symlink case red.
- [ ] The spec's edge-inventory symlink line cites the new test rather than the
      owner's seal-and-executable test.
