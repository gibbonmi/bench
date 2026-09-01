# Name the proofs a group did not run

Blocked by: declare-the-landing-refusal-registry.md
Writes: internal/worktree/land.go, internal/worktree/land_refusal.go, internal/worktree/land_surface_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LRS9

## What to build

The landing preflight groups its proofs, and a fault short-circuits the group.
The refusal states its own fault today, and it says nothing about the proofs
behind it. The refusal gains one sentence that states that the group's later
proofs did not run. The `next=` field names `later proofs in this group did not
run`. The operator then expects another refusal after the repair.

The short-circuit face constructs through the registry constructor that the
blocker landed in `internal/worktree/land_refusal.go`. That constructor takes
the route as a required argument, so this ticket composes no route of its own.
Every ticket that writes `internal/worktree/land_refusal.go` builds serially
after the registry ticket, one at a time.

## Acceptance

- [ ] LRS9 — a refusal from a short-circuited proof group names `later proofs
      in this group did not run` in its `next=` field.
- [ ] A refusal from a complete proof group carries no such sentence.
- [ ] The short-circuit face constructs through the registry constructor in
      `internal/worktree/land_refusal.go`.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read the preflight proof
groups in `internal/worktree/land.go`. Read
`TestLandCommandReportsIdentityAndDestinationInOnePreflight` in
`internal/worktree/land_surface_test.go`.

Your blocker landed the registry constructor in
`internal/worktree/land_refusal.go`. Carry the sentence through that
constructor; do not compose it at the call site.

The five `cmd/bench` and `internal/conformance` entries are the registry
closure for the `internal/worktree` package. Edit them only if your change
reaches them.

Coverage rows: LRS9. Show LRS9 red before your edit. Show LRS9 green after.
Return the red-to-green log.

Self-probe with an omission mutation. Drop the sentence from the
short-circuited group and report the observed result. If the mutation returns
green, add the missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/worktree/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
