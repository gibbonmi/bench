# Name the retire step in its refusal

Blocked by: none
Writes: internal/spec/spec.go, internal/spec/spec_test.go, internal/roadmap/roadmap_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LRS16, LRS18

## What to build

`bench spec retire` refuses the primary checkout with the shared refusal today.
That shared line names the worktree creation, and it names no retire step. The
retire verb composes its own route. The verb appends its own follow-on step, so
the operator reads one route.

`usage.PrimaryCheckoutRefusal` keeps its current text. `bench commit` and the
idea verb stay unchanged, so the shared source stays one source. The idea
verb's refusal therefore names no `bench spec retire` step.

## Acceptance

- [ ] LRS16 — the `bench spec retire` primary-checkout refusal names
      `bench spec retire` as the follow-on step.
- [ ] LRS18 — the idea verb's primary-checkout refusal names no
      `bench spec retire` step.
- [ ] `internal/usage/worktree.go` stays unchanged.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 2 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read
`usage.PrimaryCheckoutRefusal` in `internal/usage/worktree.go` as read-only
context. Read the retire path in `internal/spec/spec.go`. Read
`TestRetireOnPrimaryCheckoutRefusesAndDeletesNothing` in
`internal/spec/spec_test.go`. Read `TestIdeaRefusesPrimaryCheckout` in
`internal/roadmap/roadmap_test.go`.

`internal/usage/worktree.go`, `internal/commit/commit.go`, and
`internal/roadmap/learning.go` are outside your fence. Do not edit them. Append
the retire step at the retire call site.

The five `cmd/bench` and `internal/conformance` entries are the registry
closure for the `internal/roadmap` package. Edit them only if your change
reaches them.

Coverage rows: LRS16, LRS18. Show each row red before your edit. Show each row
green after. Return the red-to-green log per row.

Self-probe with a swap mutation. Move the retire step into the shared function
and report the observed result. If the mutation returns green, add the missing
row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/spec/ ./internal/roadmap/ ./internal/commit/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
