# Refuse a non-regular graded file before the conformance read opens it

Blocked by: none
Writes: specs/lp-conformance-fifo/tickets/01-refuse-a-non-regular-graded-file.md, internal/conformance/checks_test.go, internal/conformance/read_if_exists_test.go
Covers: none

## What to build

The conformance helper `readIfExists` sends a producer path through the no-follow
classifier, which refuses a FIFO before it opens it. Every other path goes to a plain
`os.ReadFile`. The `package-core-guard` check sweeps `.github`, `scripts`, `internal`,
`tests`, and `docs` through that helper, and its walk skips only directories. A FIFO
in a swept tree, for example `docs/x.md`, blocks the read in open(2). The check then
hangs until the test deadline kills it.

Make the plain-read branch of `readIfExists` refuse a path that is not a regular file
before it opens the path. Keep the branch's link behavior: it follows a link to a
regular file. Add a focused regression test that sweeps a temporary tree with a FIFO
ahead of a stale reference, with a derived deadline, so the test cannot hang the suite.

## Acceptance

- [ ] `bench test --check package-core-guard` passes in a tree that holds a FIFO at `docs/x.md`.
- [ ] `TestReadIfExistsPassesAFIFO` fails within its deadline when the regular-file test in `readIfExists` is weakened to a directory test.
- [ ] `TestReadIfExistsPassesAFIFO` passes, and the sweep still reports the stale reference that sorts after the FIFO.
- [ ] The producer branch of `readIfExists` is unchanged.
