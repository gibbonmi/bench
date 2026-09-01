# Add the buildvcs defense to the citation expansion

Blocked by: none
Writes: internal/coverage/citation_execution.go, internal/coverage/citation_execution_test.go, internal/gate/gate_go.go, internal/gate/lane.go, internal/gate/lane_test.go

## What to build

The coverage citation expansion runs `go list` without the
`-buildvcs=false` defense. Go's VCS discovery treats a linked worktree's
`.git` file as no root. It walks up and adopts a stray `.git` directory
above the temporary checkout. Git refuses that directory, and the
expansion fails with `error obtaining VCS status`. The gate's build and
test commands already carry the defense
(`internal/gate/gate_go.go:27`, `disableBuildVCS`). The comment at
`internal/gate/lane_test.go:176-180` documents the hazard.

Give `selectedPackageDirs` in
`internal/coverage/citation_execution.go` the same defense. Share one
source: export the constant from the gate package and use it in both
call sites. Do not weaken what the expansion proves. A package operand
that does not resolve still fails the expansion.

## Acceptance

- [ ] A stray `.git` directory above the phase checkout no longer fails
      the citation expansion. A new test proves it, on the pattern of
      `TestBenchkitLaneBuildIgnoresAStrayGitDirAboveTheCheckout`
      (`internal/gate/lane_test.go:181`).
- [ ] The new test is red without the flag and green with it. Record the
      red-to-green log.
- [ ] An unresolvable package operand still fails the expansion, and the
      existing `citation_execution_test.go` suite stays green.

## Delegate charge

You work in the Bench repo. Line: opus / medium. Effort: medium, at most
3 iterations.

Read `internal/coverage/citation_execution.go` and its test file in
full. Read `internal/gate/gate_go.go:20-40` and
`internal/gate/lane_test.go:176-210`.

Show the new test red before the fix and green after. Self-probe with a
swap mutation: point the exported constant at `-buildvcs=true` and
report the observed result.

Run `bench worktree exec "<label>" -- go test -parallel 2
./internal/coverage/ ./internal/gate/`. The exec form is the only
command form. Do not use `cd`. Do not commit. Do not edit this ticket.
