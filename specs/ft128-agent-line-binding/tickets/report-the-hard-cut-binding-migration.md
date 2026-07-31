# Report the hard-cut binding migration

Blocked by: Migrate line binding through the harness matrix

## What to build

`bench doctor` identifies retired tier and alias keys, prints their exact
harness-matrix rewrites, and leaves the reviewer-owned binding file unchanged.

## Acceptance

- [ ] Doctor names every retired key and its replacement.
- [ ] Re-running doctor produces the same report without changing `lines.env`.
