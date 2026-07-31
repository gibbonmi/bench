# Expose spec-build porcelain

Blocked by: Promote reviewed candidates exactly

Ownership fence: `internal/spec`, `cmd/bench`, `internal/contract/runtime`
Assumptions: internal/specbuild exposes the complete lifecycle without leaking record fields

## What to build

Expose the eight-command `bench spec build` family through thin `internal/spec`
and binary dispatch adapters, then drive it through black-box runtime contracts
using the real kit and linked-repo wrappers. Parsing and TOON rendering stay at
the adapter; lifecycle state and Git mutation stay in `internal/specbuild`.

## Acceptance

- [ ] [R24, R27] Status and every mutator return the post-transition TOON projection with exit 0/1/2, exact no-op/conflict behavior, and no prompts.
- [ ] [R30, R33] Runtime calls prove only promote reaches prospective gate execution and full status retains attributed evidence after squash.
- [ ] [R40] A red/repair/green runtime trace pays no ticket gate and requires a fresh composed review.
- [ ] [R53, R55-R56] Flag values and `--` parse literally; control text cannot split output; nested cwd and both supported wrappers resolve one common-directory run.
- [ ] [R57] Missing Git, gate, or real binary reports one actionable recovery route without unsafe mutation.
