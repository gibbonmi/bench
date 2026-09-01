# Route the undeclared residue refusal

Blocked by: declare-the-landing-refusal-registry.md
Writes: internal/worktree/land_identity.go, internal/worktree/land_refusal.go, internal/worktree/land_release_refusal_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LRS10, LRS11

## What to build

The residue refusal names its paths today, and it names no repair. The refusal
gains a route with two parts. The route names the declaration file
`.bench/build-outputs.json`. The route also names an exact removal command for
the undeclared path, so the operator chooses between the two repairs.

The residue face constructs through the registry constructor that the blocker
landed in `internal/worktree/land_refusal.go`. That constructor takes the route
as a required argument, so this ticket composes no route of its own. Every
ticket that writes `internal/worktree/land_refusal.go` builds serially after
the registry ticket, one at a time.

## Acceptance

- [ ] LRS10 — the residue refusal's `next=` field names
      `.bench/build-outputs.json`.
- [ ] LRS11 — the residue refusal's `next=` field names an exact removal
      command for the undeclared path.
- [ ] The residue face constructs through the registry constructor in
      `internal/worktree/land_refusal.go`.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read `landingDestination`
and `undeclaredLandingIgnoredPaths` in `internal/worktree/land_identity.go`.
Read `releaseNext` in `internal/worktree/ownership.go` as the prior art for a
residue route. Read `TestLandCommandRefusalListsIgnoredPaths` in
`internal/worktree/land_release_refusal_test.go`.

Coverage rows: LRS10, LRS11. Both rows ride in
`TestLandCommandRefusalListsIgnoredPaths`. Show each row red before your edit.
Show each row green after. Return the red-to-green log per row.

Self-probe with an omission mutation. Drop the declaration file from the route
and report the observed result. If the mutation returns green, add the missing
row.

The five `cmd/bench` and `internal/conformance` entries are the registry
closure for the `internal/worktree` package. Edit them only if your change
reaches them.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/worktree/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
