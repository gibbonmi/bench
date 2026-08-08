# Reduce gate fixture materialization

Blocked by: share-the-kit-shaped-staged-binary.md
Ownership fence: `internal/gate/kitshaped_fixture_test.go`, `internal/gate/prospective_test.go`, `internal/gate/check_slots_test.go`, `internal/gate/build_attestation_test.go`, `internal/gate/build_skip_test.go`, `internal/gate/component_decision_test.go`, `internal/gate/runner.go`, `internal/gate/phases.go`, `internal/gate/runner_test.go`, `internal/gate/runner_serial_test.go`, `internal/gate/resource_bounds_test.go`, `internal/gate/phases_command_test.go`, `internal/gate/story4_proof_test.go`
Integration surfaces: process-scoped staged template→per-root hardlink with portable copy fallback + GFM2-GFM3; in-place executable metadata mutation→private detach and byte replacement→atomic root-local rename + GFM3; package `TestMain`→ordinary resolved Go cache for nested temporary roots + GFM1/GFM6; fixture-only check-slot edits→freshness assertion, store-only attestation setup and decision-site faults→direct validated evidence + GFM4-GFM6; cancellation context→test-scoped grace while the production default remains two seconds + GFM5/GFM9; package resource budget→serial width-one and width-two commands + GFM5-GFM7
Contracts: every test retains a private mutable root, Git directory, freshness seal, and Bench evidence store; only the immutable initial executable inode crosses roots, and any in-place byte or metadata mutation detaches first while unlink/rename replacement remains root-local; linked worktrees remain excluded because their common Git directory would share Bench evidence; nested Go commands inherit the caller's explicit cache or Go's ordinary resolved default rather than manufacturing one cache per temporary root; changed-source, alternate-artifact, planted-binary, prospective-execution, and build-authorship assertions retain real builds; production process cancellation retains its two-second courtesy while only context-marked tests use the shorter grace
Closure: GFM1/runtime-census, GFM1/shared-go-cache, GFM2/private-git-dir, GFM2/private-evidence, GFM3/immutable-hardlink, GFM3/private-detach, GFM3/copy-fallback, GFM4/real-build-retention, GFM4/freshness-assertion, GFM4/store-only-attestation, GFM4/decision-only-evidence, GFM5/width-one-wall, GFM6/width-one-output, GFM7/two-core-contention, GFM8/full-behavior, GFM9/production-grace, GFM9/test-grace

## What to build

Use the fresh census to remove four kinds of repeated setup. Resolve Go's ordinary
default build cache once before package tests so nested prospective roots do not each
create a cold `dist/.freshness-go-cache`. Hardlink the immutable process template into
each fixture root, falling back to an ordinary copy where the filesystem refuses links,
and atomically detach before in-place metadata mutation or rename a private temporary file
over the fixture path for byte replacement. Replace check-slot re-links
whose edits are outside the synthetic binary's build closure with an explicit freshness
assertion, and seed attestation parser/store tests from already published bytes while
retaining real builds in every behavior-owning class. Seed decision-site refusal tests
through the same valid slot and attestation writers, independently read that evidence,
and call the production scoping decision without executing behavior-unrelated marker phases.

Keep every working tree, Git directory, seal, and Bench evidence store private; ordinary
`git worktree` children remain unsafe because their common Git directory would collide.
Finally, pass a shorter process-group cancellation grace only through test contexts. The
production two-second grace remains independently pinned, while serial tests no longer
pay it repeatedly merely to prove the same TERM/INT-before-KILL cascade.

## Acceptance

- [x] [GFM1] the removed diagnostic census records 101 initial materializations, 23 real builds, the isolated prospective test's 358,160 output blocks, and the fixed-grace cancellation family; nested gates inherit one ordinary resolved Go cache.
- [x] [GFM2] every realized fixture has its own Git directory and Bench evidence location; no two live fixtures share lock files, refs, indexes, freshness seals, or gate evidence.
- [x] [GFM3] initial executables share the immutable inode; metadata mutation atomically detaches first, byte replacement atomically renames a private file, and replacing or removing one fixture path cannot change another fixture or the template.
- [x] [GFM4] changed-source, alternate-artifact, planted-binary, prospective-execution, and build-authorship tests continue to perform the real builds that their assertions observe; parser, store, and decision-only tests seed and independently read valid evidence without behavior-unrelated builds or phase executions.
- [x] [GFM5] `/usr/bin/time -v env GOMAXPROCS=2 go test -p 1 -parallel 1 -count=1 -timeout 600s ./internal/gate` completes within 112.85 seconds, at least 20 percent below the restored 141.07-second serial baseline.
- [x] [GFM6] the GFM5 command reports at most 757,833 filesystem-output blocks, at least 40 percent below the restored 1,263,056-block serial baseline.
- [x] [GFM7] `/usr/bin/time -v env GOMAXPROCS=2 go test -p 1 -parallel 2 -count=1 -timeout 600s ./internal/gate` remains green, does not exceed the serial run's filesystem output by more than 5 percent, and does not exceed its maximum RSS by more than 25 percent.
- [x] [GFM8] the focused mutation witnesses and `BENCH_KIT=/nonexistent GOMAXPROCS=2 go test -race -p 1 -parallel 2 -count=1 -timeout 600s ./internal/gate` are green, combining hostile-kit and race coverage without a duplicate full sweep.
- [x] [GFM9] an unmarked production context resolves a two-second process-group cancellation grace while only explicitly marked tests resolve the 100-millisecond grace.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GFM1/runtime-census | omit one selected materialization family from the diagnostic counter | the family total versus the package-level residual | restore the diagnostic locally, omit the family, run width one, and expect the accounted operations or output contribution to fall outside the recorded selection bound |
| GFM2/private-git-dir | realize two fixtures as linked worktrees of one common Git directory | a concurrent two-fixture evidence writer | apply, make both fixtures publish independent evidence, and expect a shared evidence path, lock collision, or cross-fixture observation |
| GFM2/private-evidence | point two private Git directories at one Bench evidence directory | the same concurrent evidence writer | apply, publish distinct evidence, and expect one fixture to observe or overwrite the other's record |
| GFM3/mutable-root-isolation | skip the private detach before backdating or rewrite bytes through the shared inode instead of atomic rename | the two-fixture inode/digest witness, the executable-mismatch case, and `TestAttestedSealSkipsTheBuild` | apply each mutation, then expect the second fixture's freshness result, attestation, or template metadata to move |
| GFM4/real-build-retention | route one behavior-owning build family through immutable setup | the existing source/digest/authorship assertion for that family | apply one family at a time, run its focused test, and expect the observed changed source, digest, prospective execution, or authorship to red |
| GFM4/decision-only-evidence | restore full-green phase execution to `TestDecisionSiteFailsClosed` or stop validating seeded evidence | its three-run elapsed loop and per-fault scoping assertions | apply either mutation, run three focused repetitions, and expect the loop to exceed 12 seconds or a hostile evidence class to stop forcing its component |
| GFM5/width-one-wall | make cancellation tests use the production grace or restore setup-only real links | the exact serial timed command | apply, run width one, and expect wall time to exceed 112.85 seconds |
| GFM6/width-one-output | restore byte copies, per-root freshness caches, or setup-only real links | the exact serial filesystem-output receipt | apply one family, run width one, and expect output to move materially toward or beyond 757833 blocks |
| GFM7/two-core-contention | permit concurrent realization to share a mutable writer or unbounded expensive builder | the paired serial/two-core resource receipts | apply, run both commands, and expect output inflation, RSS inflation, a race, or a timeout beyond the stated bounds |
| GFM8/full-behavior | resolve immutable setup through ambient `BENCH_KIT` or discard synchronization | the hostile-kit package run and race detector | apply each change, run its owner, and expect foreign-kit coupling or a reported race |
| GFM9/production-grace | change the production default or ignore the scoped test value | `TestProcessGroupCancelGraceProductionDefault` | apply either mutation, run the test, and expect the exact default or scoped-duration assertion to red |
