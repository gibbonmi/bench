# The complete checkpoint refuses a declared-builds overlap

Blocked by: none
Writes: internal/gate/checkpoint_persisted_test.go, internal/gate/refusal_route_test.go
Covers: none

## What to build

`.bench/BENCH-reference.md` states that the complete checkpoint refuses with `declared-builds-overlap` while two declared builds are open.
The gate runs the candidate check for each checkpoint of a spec, so the complete checkpoint reaches this refusal.
But the overlap fixture of the gate face walk runs only the chunk checkpoint.
No test runs the complete checkpoint against two open declared builds.

The overlap fixture now takes the checkpoint that it runs.
For the complete checkpoint, the fixture records one chunk and keeps the completion evidence, so the checkpoint can pass after the route.
A new test walks the overlap face under the complete checkpoint through the shared face walk.
The face walk keeps the chunk case.
This change adds only tests, so it changes no production code and adds no `CHANGELOG.md` entry.

## Acceptance

- [ ] A test runs `bench gate --checkpoint <spec> --complete` against the overlap fixture through the shared face walk.
- [ ] That test requires a refusal whose printed route is the `declared-builds-overlap` route, and the route clears the refusal.
- [ ] The chunk case and the complete case use one overlap fixture, and no copy of it.
- [ ] A probe that removes the candidate check from the complete checkpoint makes the new test fail, and the chunk case stays green.
- [ ] `bench test` for `./internal/gate` passes, and the lane passes at `bench commit`.
