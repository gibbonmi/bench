# Tell the two merge states apart

Blocked by: declare-the-landing-refusal-registry.md
Writes: internal/worktree/land.go, internal/worktree/land_refusal.go, internal/worktree/land_journey_test.go, internal/worktree/land_surface_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LRS6, LRS7, LRS22

## What to build

The landing reads `MERGE_HEAD` in the source worktree. The read decides which
route the conflict refusal names. A source worktree that holds `MERGE_HEAD`
names `git merge --continue`. A source worktree without `MERGE_HEAD` names the
commit-and-review route, and that route omits `git merge --continue`.

An unreadable source Git directory leaves the state undecided. The route then
falls back to the commit-and-review form. The `bench commit` contract does not
change, and FT258 still owns it.

The conflict face constructs through the registry constructor that the blocker
landed in `internal/worktree/land_refusal.go`. That constructor takes the route
as a required argument, so this ticket composes no route of its own. Every
ticket that writes `internal/worktree/land_refusal.go` builds serially after
the registry ticket, one at a time.

## Acceptance

- [ ] LRS6 — a source worktree that holds `MERGE_HEAD` produces a `next=` field
      that names `git merge --continue`.
- [ ] LRS7 — a source worktree without `MERGE_HEAD` produces a `next=` field
      that names the commit-and-review route and omits `git merge --continue`.
- [ ] LRS22 — an unreadable source Git directory produces the commit-and-review
      route.
- [ ] The conflict face constructs through the registry constructor in
      `internal/worktree/land_refusal.go`.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read
`conflictRepairPrefix` and `landingConflictNext` in
`internal/worktree/land_refusal.go`. Read the `landing.ConflictError` branch in
`internal/worktree/land.go`. Read
`TestLandCommandPublicConflictRepairRequiresNewReviewedTip` in
`internal/worktree/land_journey_test.go`. Read
`TestLandCommandConflictRefusalNamesTheSourceRepair` in
`internal/worktree/land_surface_test.go`.

Coverage rows: LRS6 and LRS22 ride in the journey test, and LRS7 rides in the
surface test. Show LRS6 and LRS22 red before your edit. Show both rows green
after. Return the red-to-green log per row.

LRS7 is a guard row, and it is green before your edit by design. Give LRS7 no
red-before log. Probe LRS7 instead: append `git merge --continue` to every
conflict route and report the observed result. If LRS7 stays green under that
mutation, add the missing assertion.

Branch the existing route builder in `internal/worktree/land_refusal.go`; do
not add a second builder. The five `cmd/bench` and `internal/conformance`
entries are the registry closure for the `internal/worktree` package. Edit them
only if your change reaches them.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/worktree/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
