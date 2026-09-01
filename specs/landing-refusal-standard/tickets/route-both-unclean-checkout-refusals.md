# Route both unclean checkout refusals

Blocked by: declare-the-landing-refusal-registry.md
Writes: internal/worktree/land_identity.go, internal/worktree/land_refusal.go, internal/worktree/land_release_refusal_test.go, internal/worktree/land_identity_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LRS12, LRS13

## What to build

Both callers of `checkoutClean` pass an empty route today. Each caller gains
its own route, because the fact reads the same and the repair does not.

The dirty-destination refusal prints a route beside its moved-path table, so
the operator cleans the destination. The dirty-source refusal prints a route on
one line and no `refusal_paths` table, so the hostile-source surface stays
bounded.

Both faces construct through the registry constructor that the blocker landed
in `internal/worktree/land_refusal.go`. That constructor takes the route as a
required argument, so this ticket composes no route of its own. Every ticket
that writes `internal/worktree/land_refusal.go` builds serially after the
registry ticket, one at a time.

## Acceptance

- [ ] LRS12 — the dirty-destination refusal prints a `next=` field beside its
      moved-path table.
- [ ] LRS13 — the dirty-source refusal prints a `next=` field and no
      `refusal_paths` table.
- [ ] Both unclean checkout faces construct through the registry constructor in
      `internal/worktree/land_refusal.go`.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read `checkoutClean` and both
of its landing callers in `internal/worktree/land_identity.go`. Read the
`refusal` struct and its `table` method in `internal/worktree/classifier.go`.
Read `TestLandCommandRefusalListsDestinationPaths` in
`internal/worktree/land_release_refusal_test.go`. Read
`TestLandCommandInvalidatesAChangedSourceFingerprintBeforeTheGate` in
`internal/worktree/land_identity_test.go`.

Your blocker landed the registry constructor in
`internal/worktree/land_refusal.go`. Build both faces through that constructor.
Keep the dirty-source face off the path-carrying form.

The five `cmd/bench` and `internal/conformance` entries are the registry
closure for the `internal/worktree` package. Edit them only if your change
reaches them.

Coverage rows: LRS12, LRS13. Show each row red before your edit. Show each row
green after. Return the red-to-green log per row.

Self-probe with a swap mutation. Route the dirty-source face through the
path-carrying form and report the observed result. If the mutation returns
green, add the missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/worktree/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
