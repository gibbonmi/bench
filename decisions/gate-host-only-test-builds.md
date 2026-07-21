# Host-only builds in artifact contract tests

Reviewer-closed in session (2026-07-20): the gate-speed diagnosis ended with the
reviewer approving option 1 — host-only test builds, full matrix stays with the
release workflow. This map records those settled decisions per the
`/bench-write-spec` same-session path; it is flagged in the spec for post-hoc
veto.

## Destination

`bench gate` stops paying for 4-platform cold cross-compiles it doesn't need:
artifact contract tests build only the host target (one keeps a 2-row plan so
matrix iteration stays red-capable), the release workflow's full-matrix path is
untouched, and the hermetic per-stage cache design is preserved exactly.

## #1: Which of the three speed options ships first?

Type: Grill (closed in session)

### Question
Gate wall is ~10–15 min. Host-only test builds, a shared persistent hermetic
cache, or parallel phases?

### Answer
**Host-only test builds.** Smallest diff, no oracle weakening: hermetic staging
(`build-artifacts.sh`'s per-stage `GOCACHE`/`GOMODCACHE`/`HOME`) stays exactly as
is; tests just build 1 target instead of 4, and the 4-way matrix is still proven
where it matters — the release workflow. Rejected for now: shared cache (trades
away the hermeticity the per-stage cache exists to guarantee), parallel phases
(phases already run concurrently; the serial cost is inside the contract phase,
and parallelism interacts with the worktree-serialization learning).

## #2: Where does the target filter live?

Type: Grill (closed in session; mechanism probe-verified)

### Answer
**In the staged fixture's `scripts/release-plan.json`, not in production code.**
Tests already build from committed clones (the builder refuses dirty trees);
they commit a plan whose `targets` is filtered to the host row. Every consumer —
`release-plan.mjs` (`targets`, `artifact-names`, `archive-inventory`), the
builder's tarball-count check, the evidence builder, offline archives, the
reproducibility second build (clones the staged source), and the Go tests' own
`plan.Targets` read — derives from that one committed plan. One source per fact;
zero production-script or env-knob changes. Probe (2026-07-20): a 1-row
committed plan builds end-to-end green in 14s vs 62s for the 4-row plan.
Rejected: an env-var filter in `release-plan.mjs` (a production knob plus a
duplicate filter in Go test assertions), filtering in `build-artifacts.sh` alone
(trips its own count check against `artifact-names`).

## #3: What keeps matrix iteration red-capable?

Type: Grill (closed in session)

### Answer
**Exactly one test keeps a 2-row plan** — host plus the first non-host target —
so a builder that ignores plan rows beyond the first still turns the gate red via
the existing count and per-platform-name assertions. Everything else runs the
1-row plan. The full 4-way matrix remains release-workflow-owned; the
publication verifier still requires fully-approved evidence at release time.

## #4: Host row absent from the plan?

Type: Grill (closed in session)

### Answer
**Skip, don't fabricate.** The filter selects the plan row matching
`go env GOOS/GOARCH`; on a host with no matching row the affected tests skip
with a named reason (precedent: the rsync-absent reproducibility skip). Target
knowledge stays in `release-plan.json`; tests never invent os/arch/runner
strings.

## Out of scope

- Conformance-phase cost (~250s: clone builds, `npm pack`) — separate phase,
  separate capability.
- Parallelizing the contract phase's serial package list — separate capability.
- The `surface` package's 199s (not a `build-artifacts.sh` consumer) — unmeasured,
  separate diagnosis.

## Handoff

1. **Module boundaries.** One module: the artifact contract test package
   (`internal/contract/surface/artifact`), whose single staging helper grows a
   target-filter step (stage clone → filter+commit plan → build). Production
   scripts (`build-artifacts.sh`, `release-plan.mjs`, `go-build.sh`,
   `build-release-evidence.mjs`), the release workflow, and the gate phase
   wiring are untouched (thin: no diff at all).

2. **Contracts.** The staged source is a committed clone whose
   `scripts/release-plan.json` `targets` array is filtered to the host row
   (default) or host + first non-host row (the one breadth-keeper test);
   everything else in the plan is unchanged. `build-artifacts.sh` is invoked
   exactly as production invokes it. Test assertions that need the matrix read
   the staged clone's plan, never the repo root's.

3. **Deep vs thin.** The staging helper is the one deep spot (clone + filter +
   commit + skip-if-no-host-row); every test call site is a thin consumer.
   The fixture harness stays single-sourced per the code standard.

4. **Black-box assertables.** Builder exit 0 on the staged plan; artifact dir
   contains exactly `artifact-names`-derived files for the filtered plan
   (wrapper + per-target platform tarball + per-target offline archive);
   reproducibility.json promoted; 2-row test proves N-row iteration (2 targets →
   2 platform tarballs + 2 archives); existing interruption/concurrency/
   source-state behaviors unchanged under the 1-row plan.

5. **Gate attachment.** Entirely inside the existing contract-phase artifact
   package — the gate command is unchanged (`.bench/gate.sh`). No new gate
   check: the speedup is observable as suite wall time; correctness rides the
   existing assertions re-derived from the staged plan.

6. **Hostile-input owners.** Paths with spaces (t.TempDir staging — existing
   coverage retained, e.g. the `[hostile]` output dirs); host row absent →
   helper's named skip; dirty staged clone → builder's own clean-state refusal
   (already asserted by the source-state test, unchanged); interrupt mid-build →
   existing interruption tests now on the 1-row plan.

7. **Uncertainty flags.** None. The mechanism risk (evidence validation or
   count checks rejecting a filtered plan) was probed green in session.
