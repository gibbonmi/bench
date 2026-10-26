# Executable roles in a checkpoint, a merge, and a landing

This asset answers ticket #2. It records the state of `main` at commit `cfbd1d360` (2026-10-26).
A read-only research delegate read the sources, and the shaping session verified the claims marked "verified".

## Kit repository

| run | running executable | run binary | identity root | phase table |
|---|---|---|---|---|
| chunk checkpoint, with or without `--prepare` | wrapper selection from the worktree | built from the private checkout of the HEAD tree without the review record | the worktree | the built-in table, or the delta set |
| complete checkpoint | wrapper selection from the worktree | built from the private checkout of the published tree | the worktree | the built-in table, whole schedule |
| merge to a lane target | wrapper selection | built from the checkout when no run binary is inherited | none, because a lane writes no evidence | `BenchkitLane` |
| landing | the authenticated promotion broker | built from the private checkout of the composed published tree | the primary checkout | the built-in table, or the reader set |

- The kit has no `.bench/phases.json`, so the phase table is the built-in table (`internal/gate/manifest.go:171`).
- In the kit, `runBinarySource` answers the private checkout, so the gate builds the run binary from the graded tree (`internal/gate/gate.go:305-311`, verified).
- The build uses `-trimpath`, `-buildvcs=false`, and the pinned toolchain (`scripts/go-build.sh:7-8`, verified). So equal trees give equal run binary bytes by design.
- The running executable computes the subject and writes the evidence record in its own process (`internal/gate/run_transaction.go:239`).

## Running executable selection

- In a worktree, `kit_dir` answers the worktree top level (`bin/bench.sh:82-95`, verified).
- `bench_binary_path` tries the platform package, then the version cache, then `<worktree>/dist/bench`, then `<primary>/dist/bench` (`bin/bench.sh:225-246`, verified).
- So after `bench worktree build`, a checkpoint in that worktree runs under the worktree's own candidate build.
- An inherited `BENCH_RUN_BINARY` or `BENCH_KIT` also changes the selection (`bin/bench.sh:295-313`).
- The landing route refuses `BENCH_KIT`, `BENCH_RUN_BINARY`, and `BENCH_WRAPPER`. It runs only the broker whose digest equals `bin/bench-broker.manifest` (`bin/bench.sh:406-458`, verified).

## Key members today

- `buildSubjectForTree` frames the identity root path into the oracle digest and the identity digest (`internal/gate/subject.go:77-78`, verified).
- `baselineRunnerIdentity` digests the declared build inputs in the identity root's working tree (`internal/gate/prospective.go:57-66`, verified).
- So a worktree checkpoint and its landing always differ in two members: the identity root path, and the baseline runner identity.

## Linked repository

- The kit's `gate-prospective.sh` does not ship, so no baseline runner frame exists (`internal/gate/prospective.go:82-85`).
- The run binary source is the kit, selected by `ReuseOrOwn` (`internal/gate/prospective.go:39`).
- The run binary runs only when the repository's `gate.sh` calls `gate-phases`. The setup templates do not, so a run usually selects no run binary.

## Install exception

- `bench worktree build --install` publishes the candidate as `<primary>/dist/bench` and writes the broker manifest (`internal/worktree/build.go:263-280`).
- A declared-build checkpoint then refuses unless the running executable's digest equals that manifest (`internal/commitment/repository/window.go:78-99`, verified).
- So under the exception, the checkpoint and the landing share one running executable by construction.

## Not verified

- A byte comparison of two run binaries that two checkouts built from one tree.
- Whether a stale worktree `dist/bench` refuses before the gate verb runs.
- The reader-set eligibility rules after `internal/gate/reader_plan.go:99`.
