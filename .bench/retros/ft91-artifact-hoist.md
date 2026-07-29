# FT91 artifact hoist retrospective

## Outcome

Six ruled artifact-contract consumers now share one lazy, package-scoped,
host-narrowed artifact set. Read-only and digest belts protect the shared
fixture, later consumers fail closed after an earlier staging failure, and
package cleanup restores permissions before removal. `NarrowReleasePlanIn`
keeps the caller-owned staging variant single-sourced with the original helper.
The implemented-status transition landed green at `235ca1d`.

The artifact package improved from about 141 seconds before the hoist to
108.66 seconds in an isolated post-change run and 130–134 seconds under full
gate load. The full gate remains above FT91's 60-second stop condition because
contract and canary processes still rerun the serial artifact package.

## Gate-stage timings

- Focused delegate verification: artifact package green in 89.624 seconds.
- Recorded isolated post-hoist run: 108.66 seconds wall; package reported
  103.482 seconds.
- Ticket landing gates: artifact contract reported 123.289–134.759 seconds;
  complete gate wall was about 135–142 seconds on the less-contended runs.
- Final implemented-status gate: artifact contract 130.332 seconds; runtime
  49.411 seconds; surface 38.232 seconds; publication 20.049 seconds; axi
  4.451 seconds; canary green.
- One inherited runtime cleanup race produced the only red; the same-tier
  retry was green without an implementation change.

## Ticket-versus-spec-slice and delegate performance

The build split cleanly into one implementation ticket, one evidence ticket,
and one semantic-review repair ticket. The Sol/high implementation delegate
completed stories 1–9 with the required red-to-green probes and a focused
package run. The Sol/low evidence delegate independently measured the package
and verified that moved decision paths had no stale references. The required
fresh Sol/high/fast/yolo review found no Spec issues and identified two concrete
repairs; the original Sol/high author closed both in one focused continuation.

The slices stayed independently green, but the spec's 60-second story-11 stop
rule was predictably a measurement rather than an implementation boundary. Its
miss correctly graduated the next gate work instead of expanding this build.

## Coordinator catches

- Independently verified that two shared consumers produced exactly one
  build-log entry.
- Rejected a review suggestion to fingerprint the entire staged source because
  the approved design assigns source integrity to the write-time belt and
  fingerprints promoted artifacts; widening it would add unpriced work.
- Accepted and repaired the duplicated entry-count derivation.
- Accepted and repaired the UID-0 false belt by using the repository's
  privilege capability posture before singleton construction.
- Retained the FIFO build-log-path observation as a contestable reviewer
  judgment rather than changing a test-only diagnostic contract without a
  demonstrated failure.
- Preserved the independently produced gate assessment while keeping it out of
  FT91's implementation-status commit.

## Agent-experience improvements

### Bench CLI

Expose phase and package timings in the normal gate summary so identifying the
critical path does not require a separate instrumented run. A retirement
command could also offer a planned promotion/removal manifest before deleting
spec-local evidence.

### Skills

The full-build workflow should budget the default-branch retirement commit
explicitly when implementation lands directly on the default branch. The
current sequence can consume the declared final gate and only then discover
that retirement requires another gated documentation commit.

### Process

Treat performance slices as Amdahl-law exercises at spec time: record both the
share being removed and the process boundaries that prevent reuse. For this
repo, the next work should target per-test canary bites, exact-subject verdict
reuse, and artifact-suite scheduling rather than another local fixture hoist.
