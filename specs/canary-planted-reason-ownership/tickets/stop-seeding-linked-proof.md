# Stop seeding false linked-repo proof

Blocked by: none
Writes: `internal/adopt`

## What to build

Remove the example fixture, marker, example check, and planted-reason claims from the real init scaffold. Keep the configuration sentinel and repo-local-before-PATH resolver, call the public command only for inventory validation, leave a new linked repo red until its project supplies native checks and bindings, and preserve project-owned checks and fixtures across re-entry.

This ticket also removes the scaffold test's use of the provisional production dispatch API. That decoupling must land before `remove-production-canary-dispatch.md`; removing the API first leaves `internal/adopt` unable to compile.

## Acceptance

- [ ] (covers LR1) Real init writes no example fixture, `EXPECT`, `DO-NOT-SHIP` marker, example check, or output claiming Bench proved a project check bites.
- [ ] (covers LR2) The generated gate keeps its fail-closed configuration sentinel and selected repo-local Bench resolution, invokes `"$bench" canary "$root"` only for inventory validation, stays red for absent or empty inventory, and accepts a project two-level family only as a non-executable `project:<family>` binding.
- [ ] (covers LR3) A second init process in a path containing spaces and glob characters does not recreate the retired seed and preserves project-owned checks and fixture bindings.
