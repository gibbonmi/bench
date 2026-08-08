# FT171 shared fixture staged-binary architecture review

Status: reviewed architecture input for the `shared-fixture-staged-binary` light-path build.

Reviewer: Claude Fable, read-only architecture pass. The reviewer could inspect the
primary checkout but not the isolated candidate worktree, so its code citations are to
the primary checkout and its characterization of the parallel-marker draft uses the
recorded evidence below.

## Recommendation

The repeated disk work is fixable at one bounded fixture seam. The active three-ticket
spec build does not need to be restructured: finish the staged kit-root injection,
fixture-pin retirement, and eligible-test parallelism work, then land fixture binary
reuse as one independently-green light-path ticket after the spec retires.

Build the kit-shaped fixture's staged binary once per test process behind `sync.Once`.
For every independent mutable fixture root, copy that immutable binary into the root's
staging path and invoke the existing per-root publication path. Keep real builds only
where changed source, changed executable bytes, or build authorship is the behavior under
test.

## Evidence that selects this seam

The untouched package median was 150.85 seconds. The paused parallel-marker draft
measured 61.55 seconds at default width, 99.322 seconds with `GOMAXPROCS=2`, and 59.827
seconds under the race detector. It contains only 190 `t.Parallel` insertions across 32
`internal/gate` test files and remains unintegrated.

Resource controls show that concurrency compresses existing work rather than creating
more of it:

| run | wall | maximum RSS | filesystem-output blocks |
|---|---:|---:|---:|
| width 1 | 141.07 s | 225,480 KiB | 1,263,056 |
| width 2 | 100.47 s | 225,828 KiB | 1,262,720 |

Both widths write about 617 MiB. A single focused
`TestKitShapedFixtureCarriesBuildAndCanary` run writes 15,768 output blocks, about 7.7
MiB. `newKitShapedFixture` is called about 51 times, and each call follows
`sealKitShapedBinary` to `buildFixtureBinaryTo`, which executes
`go build -buildvcs=false -o <fresh root>/dist/bench.staged ./cmd/bench`. Roughly ten
additional real build/publication calls remain elsewhere. Shared Go build and module
caches can reuse compilation, but every distinct `-o` path still pays for a fresh link
and artifact write.

The repeated constructor chain is at:

- `internal/gate/kitshaped_fixture_test.go`: `newKitShapedFixture`,
  `writeKitShapedTree`, and `sealKitShapedBinary`.
- `internal/gate/build_attestation_test.go`: `buildFixtureBinaryTo` and
  `attestationFixture.buildAndPublishAt`.
- `internal/freshness/freshness.go`: `Digest`, `Publish`, and `Verify`.

`writeKitShapedTree` is deterministic within one test process. Its inputs are literals
or derivations from the registry, payload rows, race-test synthetic sources, reduced
scope, and input resolvers. The freshness digest is content-based: it hashes
repository-relative paths and bytes, not absolute paths, mtimes, or Git state.
`Publish` computes the seal against the destination root and `Verify` recomputes it
there. A binary built once from the byte-identical template can therefore be copied and
honestly published in every independent root. Existing seal verification turns red if
a future constructor input diverges from the template.

## Architecture

Use lazy `sync.Once` construction, with `TestMain` responsible only for removing the
package-scoped temporary directory. A focused run that creates no kit-shaped fixture
then pays nothing. The once stores both the staged artifact path and any construction
error; every dependent caller reports the same error rather than allowing only the
first test to fail. The existing per-test Go-toolchain capability check remains at the
fixture entry.

The template cannot live in a test's `t.TempDir`, because its lifetime must span
parallel tests. It also cannot be passed directly to `Publish`, because publication
consumes the staged path by rename. Each fixture receives its own ordinary byte copy,
then follows the existing publication path so its seal and mutable artifact remain
root-owned.

Use a portable copy, not a hardlink or reflink:

- A hardlink would share an inode across fixtures. Tests backdate, truncate, replace,
  and rewrite artifact bytes; one parallel test could silently mutate another.
- Reflinks are not portable across ordinary developer filesystems and are unavailable
  on relevant ext4/WSL2 setups.
- A plain copy retains isolation and reduces each constructor from a relink to roughly
  one artifact-sized write.

Do not introduce a template Git repository or nested worktree yet. That adds a
recursive-copy contract, risks bypassing the single publication writer, and saves only
small Git initialization and metadata writes after the link cost is removed. Measure
the bounded binary-sharing fix first.

Do not add a resource semaphore or a second eligibility taxonomy now. The current
evidence shows flat total writes and flat peak RSS across widths; the large cost is
intrinsic repeated linking. A semaphore would add cross-test coordination while
leaving the duplicated work intact. The active spec deliberately leaves scheduling
width to the gate-budget work.

## Real builds that remain

Keep real builds in these behavior-owning classes:

1. Prospective execution that compiles genuinely different unpublished sources. The
   source difference and exact candidate binary are the subject.
2. Build-attestation cases where a distinct executable digest or distinct authorship
   is the planting or publication condition under test.
3. Any future test whose assertion explicitly observes compiler/linker behavior.

The repeated initial build inside `newKitShapedFixture` is not in those classes; it is
byte-identical setup and is the reusable target.

## Delivery slices

Keep the active three tickets unchanged.

After the spec retires, land one light-path ticket, `shared-fixture-staged-binary`, with
this independently-green outcome:

- Construct the deterministic `benchfixture` staged binary once through the existing
  `writeKitShapedTree` and `buildFixtureBinaryTo` sources of truth.
- Copy it into every fresh kit-shaped fixture root.
- Publish and seal independently in each root.
- Let `TestMain` remove only the package-scoped template directory.
- Retain real builds in the behavior-owning cases above.

Only consider a second ticket for default `buildAndPublishAt` reuse if the first
ticket's measurements leave material residual cost.

## Acceptance evidence

Use the same host and commands as the recorded baseline:

- Three `go test -count=1 ./internal/gate` repetitions before and after. The after
  median must be below the parallel draft's 61.55 seconds.
- One width-1 run must be materially below 141.07 seconds, proving intrinsic work was
  removed rather than merely overlapped.
- `/usr/bin/time -v` filesystem-output blocks at width 1 must fall materially from
  1,263,056. The expected order of improvement is roughly two thirds because about 51
  relinks are replaced by one link and per-root artifact copies.
- The focused fixture witness must fall materially below 15,768 output blocks.
- `go test -race -count=1 ./internal/gate` must stay green.
- `GOMAXPROCS=2 go test -count=1 ./internal/gate` must stay green.
- A hostile `BENCH_KIT=<foreign> go test -count=1 ./internal/gate` must stay green so
  the kit-root injection and pin-retirement guarantees survive the seam change.

## Named risks

- Constructor determinism: a future per-test build input must not reuse the shared
  binary. Existing per-root seal verification owns this red signal.
- Publication consumes staging: the package-scoped template must never be the staged
  argument passed to `Publish`.
- Once failure propagation: toolchain or build failure must fail every dependent test
  with the recorded construction error.
- Toolchain identity: reuse is per test process only. Each worktree and Go test process
  builds its own template, preventing cross-worktree staleness.

## Measured follow-up

The bounded copy implementation proved correct but a fresh runtime census found 101
constructor materializations rather than the static estimate of about 51. Its serial
package run still wrote roughly 1.3 million output blocks. The dependent
`reduce-gate-fixture-materialization.md` ticket therefore supersedes only the ordinary-copy
recommendation: initial executables hardlink the immutable template, in-place metadata
mutation atomically detaches first, byte replacement uses an atomic root-local rename, and every Git directory, seal, and Bench evidence
store remains private. This preserves the review's rejected-worktree conclusion while
addressing the measured copy residual. Decision-only refusal tests seed private valid
evidence through its production writers and call the scoping decision directly; they do not
execute a full marker-phase run merely to construct their baseline.
