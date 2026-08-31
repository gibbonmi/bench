# Add a negative ReadPublished caller

Blocked by: reject-invalid-prospective-temp-roots.md
Writes: internal/gate/prospectiveartifact/prospectiveartifact_test.go
Covers: LF23

## What to build

Add a direct negative caller for ReadPublished. Use a malformed or absent
publication that reaches its owned validation seam.

## Acceptance

- [ ] ReadPublished rejects the named invalid publication.
- [ ] The test asserts the stable failure class.
- [ ] No higher-level helper replaces the direct caller.

