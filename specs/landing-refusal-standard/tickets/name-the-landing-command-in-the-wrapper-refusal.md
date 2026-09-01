# Name the landing command in the wrapper refusal

Blocked by: none
Writes: bin/bench.sh, internal/systemtest/land_route_test.go, tests/canary/package-core-guard/unrouted-subcommand, tests/canary/package-core-guard/reintroduced-bare-skip, tests/canary/package-core-guard/bounds-duplicate-owner, tests/canary/load-validity-metadata/extensionless-gate-ref, tests/canary/docs-currency-token-diet/stale-skill-cli-reference, tests/canary/docs-currency-token-diet/stale-cli-doc-reference, tests/canary/docs-currency-token-diet/missing-cli-inventory
Covers: LRS14, LRS15

## What to build

The wrapper refuses an inherited exec route before the first repository read.
The refusal names the environment variable today, and it names no command. The
refusal text gains the exact `bench worktree land` command. The text also gains
a sentence that states that the landing runs outside `bench worktree exec`.

The refusal reads no repository. The refusal still fires before the first
repository read, which the existing route test proves with an absent marker
file.

## Acceptance

- [ ] LRS14 — the wrapper's exec-route refusal prints the exact
      `bench worktree land` command.
- [ ] LRS15 — the wrapper's exec-route refusal states that the landing runs
      outside `bench worktree exec`.
- [ ] The refusal still fires with an absent marker file.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 2 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read the exec-route refusal in
`bin/bench.sh`. Read
`TestWorktreeLandRouteRefusesInheritedRoutingBeforeRepositoryReads` in
`internal/systemtest/land_route_test.go`.

`internal/systemtest/land_route_test.go` carries the `system` build tag. Pin
`BENCH_KIT` to your worktree root for every run of that suite. An ambient
`BENCH_KIT` flips the fixture verdict under composition.

Seven canary fixtures pin `bin/bench.sh`, and your `Writes` names each one. Run
those fixture bites beside your focused checks.

Coverage rows: LRS14, LRS15. Show each row red before your edit. Show each row
green after. Return the red-to-green log per row.

Self-probe with an omission mutation. Delete the venue sentence and report the
observed result. If the mutation returns green, add the missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 -tags=system ./internal/systemtest/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
