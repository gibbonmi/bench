# Expose spec-build porcelain

Blocked by: Promote reviewed candidates exactly

Ownership fence: `internal/spec`, `cmd/bench`, `internal/contract/runtime`, `internal/gate/authorization`, the `ReleaseOwner` seam in `internal/specbuild`, and provisional cleanup in `internal/worktree`
Assumptions: `internal/specbuild` exposes the complete lifecycle without leaking record fields and already imports `internal/spec`, so parsing belongs in `internal/spec` while runtime composition stays in `cmd/bench`

## What to build

Expose the eight-command `bench spec build` family through grammar-only
`internal/spec` parsing and a thin `cmd/bench` runtime adapter, then drive it
through black-box runtime contracts using the real kit and linked-repo wrappers.
TOON rendering stays at the binary adapter; lifecycle state stays in
`internal/specbuild`; prospective-gate attribution and provisional worktree
cleanup stay with their existing owners.

## Acceptance

- [ ] [R24, R27] Status and every mutator return the post-transition TOON projection with exit 0/1/2, exact no-op/conflict behavior, and no prompts.
- [ ] [R30, R33] Runtime calls prove only promote reaches prospective gate execution and full status retains attributed evidence after squash.
- [ ] [R40] A red/repair/green runtime trace pays no ticket gate and requires a fresh composed review.
- [ ] [R53, R55-R56] Flag values and `--` parse literally; control text cannot split output; nested cwd and both supported wrappers resolve one common-directory run.
- [ ] [R57] Missing Git, gate, or real binary reports one actionable recovery route without unsafe mutation.

## Live-tree seam requirements

- `internal/spec` must not import `internal/specbuild`; that reverses the existing dependency and creates a cycle.
- The gate owner must classify candidate, inherited, and infrastructure outcomes and issue opaque exact-tree evidence; the CLI must not parse gate output or cache records.
- `start` must bootstrap a missing branch-scoped green marker only from reusable evidence for the exact branch and HEAD, using compare-and-swap and refusing conflicting markers.
- Integrated assignment cleanup must use a worktree-owned provisional release contract: exact request/path ownership, a durable checkpoint ref at the expected commit, a clean checkout with the same tree, and refusal before mutation on any mismatch. Ordinary landedness is not sufficient proof.
