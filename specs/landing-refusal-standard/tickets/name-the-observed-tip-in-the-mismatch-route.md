# Name the observed tip in the mismatch route

Blocked by: declare-the-landing-refusal-registry.md
Writes: internal/worktree/land_identity.go, internal/worktree/land_refusal.go, internal/worktree/land_journey_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LRS17

## What to build

The source-tip mismatch refusal gains a route. The route names the re-run with
the observed tip, so a moved tip has one exact next command. The refusal keeps
its `observed` and `wanted` fields, because the operator reads both values.

The route ends with the caller's own re-run. This ticket adds no accepted-forms
clause; the landing interface spec owns that clause.

The mismatch face constructs through the registry constructor that the blocker
landed in `internal/worktree/land_refusal.go`. That constructor takes the route
as a required argument, so this ticket composes no route of its own. Every
ticket that writes `internal/worktree/land_refusal.go` builds serially after
the registry ticket, one at a time.

## Acceptance

- [ ] LRS17 — the source-tip mismatch refusal's `next=` field names the
      observed tip.
- [ ] The refusal still prints its `observed` and `wanted` fields.
- [ ] The mismatch face constructs through the registry constructor in
      `internal/worktree/land_refusal.go`.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 2 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read `identityRefusal` and its
caller in `internal/worktree/land_identity.go`. Read the route builders in
`internal/worktree/land_refusal.go`. Read
`internal/worktree/land_journey_test.go` for the journey fixture shape.

Your blocker landed the registry constructor in
`internal/worktree/land_refusal.go`. Build the mismatch face through that
constructor. LRS17 rides in
`TestLandCommandPublicConflictRepairRequiresNewReviewedTip`.

The five `cmd/bench` and `internal/conformance` entries are the registry
closure for the `internal/worktree` package. Edit them only if your change
reaches them.

Coverage rows: LRS17. Show LRS17 red before your edit. Show LRS17 green after.
Return the red-to-green log.

Self-probe with a swap mutation. Put a placeholder tip in the route instead of
the observed value and report the observed result. If the mutation returns
green, add the missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/worktree/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
