# Host-only builds in artifact contract tests

Status: implemented

Map: `decisions/gate-host-only-test-builds.md` — written on the
`/bench-write-spec` same-session path from decisions the reviewer closed in
session; the map itself is flagged for post-hoc veto.

## Problem

`bench gate` costs ~10–15 minutes of wall time. The dominant cost is the
artifact contract suite: each `build-artifacts.sh` invocation cold-builds the
full 4-platform Go matrix in a fresh hermetic per-stage cache (~62s, ~3m45 CPU),
and the suite invokes it repeatedly — plus a reproducibility second build per
invocation. The 4-way matrix proves nothing per-commit that the release
workflow doesn't already prove at release time.

## Solution

Artifact contract tests build from a staged source clone whose committed
`scripts/release-plan.json` is filtered to the host target (probe-verified:
14s vs 62s per invocation). One test keeps a 2-row plan so matrix iteration
stays red-capable. Production scripts, the release workflow, and the hermetic
staging design are untouched. Expected result: artifact suite ~280s → ~90s,
gate wall roughly halved.

## User stories

1. As a developer running `bench gate`, I want every artifact contract test
   that invokes `build-artifacts.sh` to build from a staged clone with a
   host-only committed plan, so the gate stops paying for cross-compiles that
   prove nothing per-commit.
   Line: sonnet (Luna) / medium. This is test-harness work at an existing,
   probe-verified seam whose correctness the gate fully observes.

2. As the release owner, I want exactly one artifact test to stage a 2-row plan
   (host plus the first non-host target), so a builder that ignores plan rows
   beyond the first still turns the gate red through the existing count and
   per-platform-name assertions.
   Line: sonnet (Luna) / medium. The assertions already exist and the staged
   plan is the only new input, so the gate observes the whole behavior.

3. As a developer on a host absent from the release plan, I want the staging
   helper to skip the affected tests with a named reason instead of fabricating
   a target row, so target knowledge stays single-sourced in
   `release-plan.json`.
   Line: sonnet (Luna) / medium. A small, mechanical branch in one helper with
   a direct test.

4. As the release owner, I want the production build surfaces —
   `build-artifacts.sh`, `release-plan.mjs`, `go-build.sh`,
   `build-release-evidence.mjs`, the release workflow — byte-untouched by this
   change, so the shipped artifact path cannot regress.
   Line: sonnet (Luna) / medium. The story is an absence of change, enforced by
   diff scope and the existing conformance checks over the release workflow.

5. As the release owner, I want the reproducibility second build to keep
   running in every test invocation, so hermetic double-build comparison stays
   proven per-commit (now against the staged plan its clone inherits).
   Line: sonnet (Luna) / medium. The builder itself fails when the second build
   diverges, so the gate observes it without new assertions.

## Implementation decisions

- The filter lives in the staged fixture's committed `release-plan.json`
  `targets` array, never in production code or an env knob — every consumer
  (plan queries, count checks, evidence, offline archives, the reproducibility
  clone, and Go test assertions) derives from that one committed plan (map #2).
- The artifact package's staging logic stays a single-sourced helper (code
  standard): stage clone → filter `targets` to the host row (or host + first
  non-host row for the breadth-keeper) → commit → return the staged root, or
  skip with a named reason when no plan row matches `go env GOOS/GOARCH`
  (map #4).
- Test assertions that need the matrix read the staged clone's plan, never the
  repo root's.
- The breadth-keeper is the distributable-artifact acceptance test, since that
  seam names the exact tarball set as its subject (map #3).
- No new gate check and no gate wiring change; the speedup lands inside the
  existing contract phase.

## Testing decisions

- A good test here invokes `bash scripts/build-artifacts.sh <staged-root> <out>`
  exactly as production does and asserts on its observable outputs: exit code,
  the artifact directory against `release-plan.mjs artifact-names` for the
  staged plan, and promoted `reproducibility.json`. No test reaches into
  staging internals.
- Seam: the `build-artifacts.sh` CLI — the profile's distributable-artifact
  contract seam. Prior art: every test in `internal/contract/surface/artifact`
  already attaches here; `artifact_source_state_test.go` already stages
  committed clones with git fixtures.
- Gate command: `.bench/gate.sh` (the project gate).

### Seam diagram

    trigger: go test (contract phase) — each artifact test
        │
        ▼
    staged clone (committed          ┌──────────────────────────┐
    host-filtered release-plan) ──▶  │  build-artifacts.sh      │ ──▶ artifact dir
    output dir                  ──▶  │  (production script,     │     (tarballs, archives,
                                     │   untouched)             │      reproducibility.json)
                                     └──────────────────────────┘
                      ◀ tests attach here: run the script, then assert the
                        artifact set against the staged clone's own plan

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | staging helper produces a committed clone whose plan holds exactly the host row, and the builder goes green on it | build-artifacts.sh CLI | new fixture test written first: asserts the staged plan's targets equal the host row and the builder exits 0 with the 1-row artifact set; red until the helper exists | a helper that stages the unfiltered plan fails the exact-row assertion; a broken staging commit makes the builder's clean-state check fail |
| 1 | artifact set matches `artifact-names` for the staged (not root) plan | build-artifacts.sh CLI | existing count/name assertions re-pointed at the staged clone's plan; red if any assertion still reads the root plan (4-row expectation vs 1-row output) | a test still reading the root plan expects 4 platforms and fails against host-only output |
| 1 | suite wall time drops | n/a | not TDD-able — performance; verified by measured before/after suite times reported at final check | timing is environment-dependent; an assertion would be flaky, measurement is honest |
| 2 | a 2-row staged plan yields 2 platform tarballs + 2 offline archives + wrapper | build-artifacts.sh CLI | existing count and per-platform-name assertions running on the 2-row fixture; demonstrated red once against a first-row-only builder mutation during the TDD pass | a builder that ignores rows beyond the first emits 1 platform set; the count assertion fails |
| 3 | no plan row matches the host → named skip | staging helper (in-package) | new test written first: helper given a plan filtered to a non-host target reports skip with the named reason; red until the branch exists | without the branch the helper stages a hostless plan and the builder's own artifact-count check fails confusingly instead of skipping |
| 4 | production surfaces byte-untouched | review / diff scope | not TDD-able — absence of change; enforced by diff scope at review and the existing release-workflow conformance checks | any diff to the named scripts or workflow is visible in the review diff and graded by existing conformance |
| 5 | reproducibility second build still runs and promotes reproducibility.json | build-artifacts.sh CLI | already covered — the builder exits nonzero when the second build fails or diverges, and existing tests assert the promoted file | the second build clones the staged source, so a filter that broke it would fail the builder itself |

### Edge inventory

- host row absent from the plan → coverage row (story 3).
- dirty staged clone → already covered: `artifact_source_state_test.go` asserts
  the builder's clean-state refusal; unchanged by this spec.
- paths with spaces/glob characters in output dirs → already covered: the
  existing `[hostile]`-named output fixtures are retained on the 1-row plan.
- interrupt (SIGINT/stage interruption) mid-build → already covered: the
  existing interruption tests run unchanged on the 1-row plan.
- absent vs empty `targets` array → **Won't handle** — `release-plan.mjs`
  already fail-closes on empty-target cardinality in production; this spec
  never produces an empty filter result (story 3 skips first).
- required tool missing (rsync/node/npm) → **Won't handle** — existing skip and
  hard-requirement behavior is unchanged; no new tool is introduced.
- invocation through a symlink / non-root cwd → **Won't handle** — the builder's
  own path normalization is production behavior this spec does not touch.

## Out of scope

- **Conformance-phase cost (~250s)** — a different phase with different
  mechanisms (clone `go build ./...`, `npm pack`); separate capability —
  6 edits, 3 gate runs.
- **Parallelizing the contract phase's serial package list** — gate-runner
  change, interacts with worktree-serialization learning; separate capability —
  3 edits, 2 gate runs.
- **`surface` package cost (199s)** — not a `build-artifacts.sh` consumer;
  needs its own diagnosis first — 1 debug session to size.
