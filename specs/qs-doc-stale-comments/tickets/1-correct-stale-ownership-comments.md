# Correct the stale ownership comments

Blocked by: none
Writes: internal/subprocess/subprocess.go, internal/axi/action.go, internal/consumers/citation.go, internal/publication/fixture_registry.go, internal/publication/registry.go, internal/publication/npm_registry.go, internal/publication/plan.go, internal/roadmapflow/flow.go
Covers: none

## What to build

The quality survey of 2026-09-29 found four comments that state a false fact about the tree. No anchor entry, conformance check, or test pins the text of these comments. Change only the comments, and change no code.

- The `internal/subprocess` package comment calls the package the one seam for running an external command. Only two conformance test files call `Capture`, and many production owners start, cancel, and kill their own process groups. State what the package holds, and say that it does not own the process group.
- The `axi.ShellQuote` comment calls the function the one derivation of the kit's shell quoting. `sanitize.ShellQuote` and a private quoter in package `freshness` also quote shell words. The `citation.cmd` comment in `internal/consumers` repeats the claim. Remove the claim from both comments.
- The `FixtureRegistry` comment calls it the adapter that the gate exercises. No gate path starts the `offline-registry.mjs` fixture. The gate runs `FixtureRegistry` only against an unreachable base URL, and it runs `NPMCLIRegistry` against an `npm` stub. Three more comments in `internal/publication` repeat a false gate claim. Keep the gate-coverage fact once, in the `Registry` port comment, and remove it from the adapter and package comments.
- The `internal/roadmapflow` package comment says that package `roadmap` sits at 11 source files against a reviewer-owned directory budget of 12. Package `roadmap` has 20 source files, and `.bench/structure.budgets` grants it no budget. Remove the budget sentence and the placement claim that depends on it.

## Acceptance

- [ ] No comment calls package `subprocess` the one seam for running an external command.
- [ ] No comment calls `axi.ShellQuote` the one derivation or the one owner of the kit's shell quoting.
- [ ] Only the `Registry` port comment states the gate coverage of each adapter, and it says that no gate path starts a registry.
- [ ] No comment in `internal/roadmapflow` cites a directory budget for package `roadmap`.
- [ ] The diff changes comment lines only, and `go vet ./...` and `bench test --changed` pass.
