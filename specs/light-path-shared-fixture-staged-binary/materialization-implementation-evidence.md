# Gate fixture materialization implementation evidence

## Census

The fresh width-one JSON census was green in 140.62 wall seconds with 235,544
KiB maximum RSS and 1,384,784 filesystem-output blocks. Earlier env-gated
operation logging, removed before implementation, counted 101 initial fixture
materializations and 23 successful real builds. The largest isolated families
were:

| family | wall | output blocks |
|---|---:|---:|
| exact unpublished prospective execution | 6.47s | 358,160 |
| FT78 story-four proof ledger | 6.37s | 106,752 |
| public-document check partition | 18.11s | 29,976 |
| declared-document invalidation | 4.84s | 27,192 |

The census selected four bounded reductions: share Go's ordinary default cache
with nested temporary roots, hardlink immutable initial executables and detach
before in-place mutation, replace twelve setup-only links with freshness or
store-specific setup, and shorten only the cancellation tests' process-group
grace through their contexts.

## Resource receipts

The exact width-one acceptance command was green in 111.887 package seconds and
112.23 wall seconds, with 225,612 KiB maximum RSS and 682,888 output blocks. It
meets the 112.85-second and 757,833-block limits. Against the fresh census, wall
fell 20.2 percent and output fell 50.7 percent.

The exact width-two command was green in 78.305 package seconds and 78.69 wall
seconds, with 224,528 KiB maximum RSS and 688,496 output blocks. Compared with
width one, wall fell 29.9 percent, maximum RSS fell 0.5 percent, and output rose
0.8 percent; modest two-core overlap adds no material resource contention.

The combined hostile-kit/race command was green in 86.386 package seconds and
88.74 wall seconds, with 226,528 KiB maximum RSS and 726,056 output blocks:
`BENCH_KIT=/nonexistent GOMAXPROCS=2 go test -race -p 1 -parallel 2 -count=1 -timeout 600s ./internal/gate`.

The three-fixture hardlink witness wrote 20,768 blocks, versus 30,376 for the
copy-based intermediate. The exact prospective test now writes 223,248 blocks
with the package-resolved default Go cache and deterministic ordinary-route
control, versus 358,160 in the census. The focused cancellation family is green
in 5.181 package seconds with its 100-millisecond context grace.

## Behavioral boundaries

- Every fixture still calls `git init` and owns its Git directory, index, refs, freshness seal, and Bench evidence paths. No test fixture is a linked Git worktree.
- Initial executable paths share only the immutable template inode. `makeFixtureBinaryPrivate` atomically replaces the path before truncation or backdating; ordinary removal and freshness publication replace directory entries and are already root-local.
- The check-slot edits replaced by `requireKitShapedBinaryFresh` are outside the synthetic `cmd/bench` build closure. The assertion fails if a future edit makes that assumption false.
- Attestation parser and atomic-store tests seed a valid record from already published bytes. Changed source, planted bytes, alternate artifacts, prospective execution, executable-digest authorship, and every compiler-observing assertion retain real builds.
- An unmarked context retains the two-second production cancellation grace. Only the signal/cancellation tests carry the 100-millisecond value, and `TestProcessGroupCancelGraceProductionDefault` independently pins both resolutions.

## Focused proof

The link-unavailable and existing-destination fallback witness was green together
with the detach, backdate, and production-grace seams in 0.979 seconds. The six
affected check-slot tests were green together in 23.008 seconds. The
five attestation tests were green together in 2.437 seconds. The cancellation
and signal family was green together in 5.181 seconds. `git diff --check` was
green before resource validation.

Removing the private detach from the hardlink witness failed in 0.316 seconds
because the first executable still shared the template inode. Exact restoration,
together with the backdated-artifact and production-grace pins, was green in
0.981 seconds.

## Review

Local exact-diff Standards, Spec, and Coverage review found and repaired the
unwitnessed copy fallback and its existing-destination truncation risk. Bounded
Opus reviews of the dependent full diff exhausted their harness turn or time
limits without returning a verdict; that is recorded as unavailable, not clean.
The earlier Opus review of the shared-template seam returned concrete isolation,
error, and cleanup findings, and its repair pass verified those behavior fixes.
