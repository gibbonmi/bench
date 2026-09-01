# Route every resume refusal

Blocked by: declare-the-landing-refusal-registry.md
Writes: internal/worktree/land_resume.go, internal/worktree/land_refusal.go, internal/worktree/land_resume_refusal_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LRS5

## What to build

The resume path prints its own refusals, and some of them carry a bare detail.
Each resume refusal prints a `next=` field, so the route survives a resume. The
resume refusal prints the first run's refusal shape.

The resume faces join the same registry the first run uses. The resume route
keeps the existing `landingResumeNext` command, because that command already
repeats the caller's flags.

Each resume face constructs through the registry constructor that the blocker
landed in `internal/worktree/land_refusal.go`. That constructor takes the route
as a required argument, so this ticket composes no route of its own. Every
ticket that writes `internal/worktree/land_refusal.go` builds serially after
the registry ticket, one at a time.

## Acceptance

- [ ] LRS5 — each resume refusal prints a `next=` field that names the
      `bench worktree land --resume` continuation.
- [ ] The resume route names the `--resume` re-run with the published commit.
- [ ] Each resume face constructs through the registry constructor in
      `internal/worktree/land_refusal.go`.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read
`internal/worktree/land_resume.go` in full. Read `landingResumeNext` in
`internal/worktree/land_refusal.go`. Read
`TestResumeLandCommandPublicRefusesDestructiveDestinationState` in
`internal/worktree/land_resume_refusal_test.go`.

Your blocker landed the registry constructor in
`internal/worktree/land_refusal.go`. Declare each resume face in that same
registry, so the registry walk drives it too.

The five `cmd/bench` and `internal/conformance` entries are the registry
closure for the `internal/worktree` package. Edit them only if your change
reaches them.

Coverage rows: LRS5. Show LRS5 red before your edit. Show LRS5 green after.
Return the red-to-green log.

Self-probe with an omission mutation. Drop the route from one resume face and
report the observed result. If the mutation returns green, add the missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/worktree/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
