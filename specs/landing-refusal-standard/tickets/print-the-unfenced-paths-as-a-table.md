# Print the unfenced paths as a table

Blocked by: declare-the-landing-refusal-registry.md
Writes: internal/preflight/gather.go, internal/preflight/decision.go, internal/preflight/decision_test.go, internal/worktree/land_identity.go, internal/worktree/land_surface_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LRS4

## What to build

The unauthorized paths travel typed out of the preflight package.
`preflight.AuthorizeReviewedSource` returns a typed error that carries the
unauthorized path list. `landingSourceRange` reads that list into the refusal's
`paths` field. The fence refusal then prints `paths_total=` and one
`refusal_paths` row for each unfenced path.

The joined detail string stays in `pathsAuthorizedCheck`, because the preflight
report reads it. The value contract crosses the package boundary, so the typed
error and its landing reader land together.

The fence face constructs through the registry constructor that the blocker
landed in `internal/worktree/land_refusal.go`. That constructor takes the route
as a required argument, so this ticket composes no route of its own. Every
ticket that writes `internal/worktree/land_refusal.go` builds serially after
the registry ticket, one at a time.

## Acceptance

- [ ] LRS4 — the fence refusal prints `paths_total=` and a `refusal_paths` row
      that names the unfenced path.
- [ ] `preflight.AuthorizeReviewedSource` returns the unauthorized paths as a
      typed value, not as one sentence.
- [ ] The `paths-authorized` preflight row keeps its joined detail string.
- [ ] The fence face constructs through the registry constructor in
      `internal/worktree/land_refusal.go`.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read
`AuthorizeReviewedSource` in `internal/preflight/gather.go` and
`pathsAuthorizedCheck` in `internal/preflight/decision.go`. Read
`landingSourceRange` and `identityRefusal` in
`internal/worktree/land_identity.go`. Read the `refusal` struct and its
`table` method in `internal/worktree/classifier.go`.

Your blocker landed the registry constructor in
`internal/worktree/land_refusal.go`. Build the fence face through that
constructor. LRS4 rides in `TestLandCommandFenceRefusalNamesThePath`.

The five `cmd/bench` and `internal/conformance` entries are the registry
closure for the `internal/worktree` package. Edit them only if your change
reaches them.

Coverage rows: LRS4. Show LRS4 red before your edit. Show LRS4 green after.
Return the red-to-green log.

Self-probe with a swap mutation. Move the path list back into the detail
sentence and report the observed result. If the mutation returns green, add the
missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/preflight/ ./internal/worktree/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
