# Reduce gate fixture materialization

Blocked by: share-the-kit-shaped-staged-binary.md
Ownership fence: `internal/gate/kitshaped_fixture_test.go`, the exact fixture-materialization helpers and tests named by the fresh runtime census
Integration surfaces: process-scoped staged template and per-root publication→`share-the-kit-shaped-staged-binary.md` + SFB5-SFB9; runtime fixture and real-build census→the exact repeated filesystem authors named by GFM1; independent mutable fixture roots→the selected materialization primitive + GFM2-GFM4; package resource budget→serial width-one and two-core commands + GFM5-GFM7
Contracts: every test retains a private mutable root and private freshness evidence while immutable setup may be shared only through an operation whose write and mutation semantics are proved on the supported local filesystems; Bench evidence must remain in an independent Git directory per fixture, so linked worktrees that share a common Git directory are excluded unless the evidence store is first made fixture-local; resource acceptance is priced from restored, serial same-host baselines and never inferred from default-width overlap
Closure: GFM1/runtime-census, GFM2/private-git-dir, GFM2/private-evidence, GFM3/mutable-root-isolation, GFM4/real-build-retention, GFM5/width-one-wall, GFM6/width-one-output, GFM7/two-core-contention, GFM8/full-behavior

## What to build

Start from a fresh runtime operation census, not the earlier static call-site
estimate. The shared-binary repair observed 101 constructor copies, 23 successful
behavior-owning real builds, and roughly 722,000 filesystem-output blocks outside
that helper in one serial package run. Name the repeated repository and file-tree
materializers that own that residual, then replace only byte-identical setup with
one immutable process-scoped source and a cheap per-test realization.

Preserve a private mutable working tree, Git directory, Bench evidence store, and
published executable for every fixture. Ordinary `git worktree` children are not
an acceptable shortcut while Bench evidence lives in the common Git directory.
Prefer a standard Go-test arrangement: package-scoped immutable inputs, `t.TempDir`
for test-owned mutable state, the ordinary shared Go build cache, and explicit real
builds only where compiler output or changed source is the behavior under test.
Choose hardlink, clone-reference, archive extraction, or another realization only
after a focused mutation proves that one fixture cannot alter another through that
mechanism on the supported local-development filesystems.

## Acceptance

- [ ] [GFM1] an env-gated diagnostic census, removed before commit, names the runtime count and filesystem-output contribution of every repeated materialization family that together explains the package residual well enough to select the implementation fence.
- [ ] [GFM2] every realized fixture has its own Git directory and Bench evidence location; no two live fixtures share lock files, refs, indexes, freshness seals, or gate evidence.
- [ ] [GFM3] mutating, replacing, truncating, backdating, or removing one fixture's working files and executable cannot change any other fixture or the immutable source.
- [ ] [GFM4] changed-source, alternate-artifact, planted-binary, prospective-execution, and build-authorship tests continue to perform the real builds that their assertions observe.
- [ ] [GFM5] `/usr/bin/time -v env GOMAXPROCS=2 go test -p 1 -parallel 1 -count=1 -timeout 600s ./internal/gate` completes within 112.85 seconds, at least 20 percent below the restored 141.07-second serial baseline.
- [ ] [GFM6] the GFM5 command reports at most 757,833 filesystem-output blocks, at least 40 percent below the restored 1,263,056-block serial baseline.
- [ ] [GFM7] `/usr/bin/time -v env GOMAXPROCS=2 go test -p 1 -parallel 2 -count=1 -timeout 600s ./internal/gate` remains green, does not exceed the serial run's filesystem output by more than 5 percent, and does not exceed its maximum RSS by more than 25 percent.
- [ ] [GFM8] the focused mutation witnesses, `go test -race -count=1 ./internal/gate`, and `BENCH_KIT=/nonexistent GOMAXPROCS=2 go test -count=1 ./internal/gate` are green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GFM1/runtime-census | omit one selected materialization family from the diagnostic counter | the family total versus the package-level residual | restore the diagnostic locally, omit the family, run width one, and expect the accounted operations or output contribution to fall outside the recorded selection bound |
| GFM2/private-git-dir | realize two fixtures as linked worktrees of one common Git directory | a concurrent two-fixture evidence writer | apply, make both fixtures publish independent evidence, and expect a shared evidence path, lock collision, or cross-fixture observation |
| GFM2/private-evidence | point two private Git directories at one Bench evidence directory | the same concurrent evidence writer | apply, publish distinct evidence, and expect one fixture to observe or overwrite the other's record |
| GFM3/mutable-root-isolation | skip the selected detach/copy boundary before an in-place mutation | the two-fixture mutation matrix | apply each truncate, rewrite, backdate, replace, and remove operation to the first fixture, then expect the second fixture's digest, bytes, or freshness result to change |
| GFM4/real-build-retention | route one behavior-owning build family through immutable setup | the existing source/digest/authorship assertion for that family | apply one family at a time, run its focused test, and expect the observed changed source, digest, prospective execution, or authorship to red |
| GFM5/width-one-wall | restore the replaced repeated materializer | the exact serial timed command | apply, run width one, and expect wall time to exceed 112.85 seconds |
| GFM6/width-one-output | restore the replaced repeated materializer | the exact serial filesystem-output receipt | apply, run width one, and expect more than 757833 output blocks |
| GFM7/two-core-contention | permit concurrent realization to share a mutable writer or unbounded expensive builder | the paired serial/two-core resource receipts | apply, run both commands, and expect output inflation, RSS inflation, a race, or a timeout beyond the stated bounds |
| GFM8/full-behavior | resolve immutable setup through ambient `BENCH_KIT` or discard synchronization | the hostile-kit package run and race detector | apply each change, run its owner, and expect foreign-kit coupling or a reported race |
