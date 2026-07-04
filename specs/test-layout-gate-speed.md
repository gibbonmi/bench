# test-layout-gate-speed

Status: staged

## Problem

Seven test files exceed the 400-line structure limit, so `bench structure` is red
and the gate's structural-debt signal is noise. Separately, the gate takes ~60s
wall because the contract suite runs twice per gate (once inside the conformance
sweep's `go test ./...`, once as the gate's explicit contract phase) and the
contract suite itself serializes ~19s of subprocess waits across tests that each
own an isolated fixture.

## Solution

Split the over-limit files along the responsibility families their subtests
already name, remove the duplicated contract run from the conformance sweep, and
let the contract tests run in parallel. Same tests, same oracle authority, about
half the wall time (~60s → ~25s), structure check green.

## User stories

1. As the reviewer, I want `internal/conformance/checks_test.go` (1643 lines)
   split by the five check families `RunConformance` already names — validity,
   docs/workflow, skills/index, line-routing, package/go-core — with
   `RunConformance` and shared helpers in a small core file, so each family reads
   as one unit and the structure check passes.
   Line: claude-sonnet-5 / low. This is a mechanical relocation within one
   package and the gate fully observes the result.
2. As the reviewer, I want `internal/contract/runtime_test.go` (710 lines) split
   by command area — status, gate/stop-hook, worktree, structure — one top-level
   test function per file following the existing `runtime_git_test.go` precedent,
   so runtime contracts are navigable by command.
   Line: claude-sonnet-5 / low. The subtests are already named by area and the
   gate fully observes the result.
3. As the reviewer, I want `runtime_shift_test.go` (510), `axi_test.go` (471),
   `doctor_test.go` (450), and `link_test.go` (427) each split into two files
   along their scenario families (shift loop vs shift adapters; guard fail-closed
   vs CLI surface contracts; doctor report/fix vs shim/postinstall; link
   marker/fence handling vs link environment/init), matching the existing
   `axi_wave2`/`axi_guards` precedent.
   Line: claude-sonnet-5 / low. Same mechanical pattern as story 2, gate-observed.
4. As the reviewer, I want `internal/canary/canary_test.go` (424) split two ways —
   sweep validation/bite behavior vs concurrency/cleanup — so the last over-limit
   file clears the check without fragmenting.
   Line: claude-sonnet-5 / low. A 24-line overage with an obvious two-family
   boundary, gate-observed.
5. As the reviewer, I want the conformance `checkGoCore` unit sweep to exclude
   `internal/contract`, because the gate already runs that suite as its own
   explicit phase with `BENCH_CONTRACT_ROOT` set, so the gate stops paying ~19s
   for a duplicate run.
   Line: claude-opus-4-8 / medium. This edits what the oracle covers, and the
   no-weakening argument (the suite still runs, once, in the gate's own phase)
   deserves a tier that can hold that reasoning.
6. As the reviewer, I want the contract tests to declare `t.Parallel()` (top-level
   functions and subtests), with the one `t.Setenv`-using test left serial, so the
   ~19s of serialized subprocess waits collapse to roughly the suite's CPU time.
   Line: claude-opus-4-8 / medium. The change is mechanical but parallel-safety
   is gate-invisible unless raced, so the audit needs judgment.
7. As the reviewer, I want a before/after measurement of the gate's wall time and
   of contract-suite subtest count reported at final check, so the speedup and
   the no-test-lost claim are evidence, not belief.
   Line: claude-sonnet-5 / low. Pure measurement and reporting.

## Implementation decisions

- Splits stay within their existing packages; no new packages, no new seams, no
  test rewritten — only relocated. Each new file carries one top-level test
  function (contract) or one check family plus its private helpers (conformance).
- The conformance core file keeps `RunConformance`, `containsDiagnostic`, and the
  cross-family helpers; each family file owns only helpers no other family uses.
- `checkGoCore` derives its test-package list from `go list ./...` filtered in Go
  code (not shell), so graded trees that lack `internal/contract` (canary
  fixtures) are unaffected. The exclusion carries a comment naming the gate phase
  that owns the contract run — that comment is the single cross-reference for the
  "contract is a separate phase" fact.
- Parallel safety rests on the existing fixture design: per-test `t.TempDir`
  roots and per-command env injection. `TestRunCapturesOutputExitAndEnvironment`
  (the only `t.Setenv` user) stays serial. A one-off
  `go test -race ./internal/contract` validates the change; `-race` is not added
  to the gate.
- Gate-phase-level concurrency in `gate.sh` is deliberately not part of this
  change (see Out of scope) — remeasure after these land.

## Testing decisions

- This change is behavior-preserving for every existing test; the oracle for
  "nothing lost, nothing weakened" is the existing gate plus a subtest-count
  parity check, not new tests.
- Gate command: `bench gate` (root conformance + contract phase + canary sweep).
- Prior art: `runtime_git_test.go` and `axi_guards_test.go` are the file-split
  pattern; the `go-test-failing` canary fixture is the existing bite-proof for
  `checkGoCore`'s unit sweep.

### Seam diagram

Seam 1 — the structure check (grades the splits):

    trigger: bench gate / bench structure
        │
        ▼
    repo tree  ──▶  [ structure line-budget check ]  ──▶  issue list (7 today)
                        ◀ tests attach here: `bench structure` exits 0 after the splits

Seam 2 — the contract suite (grades relocation + parallelism):

    trigger: gate phase 2 (BENCH_CONTRACT_ROOT=<root> go test ./internal/contract)
        │
        ▼
    kit tree ──▶  [ contract fixtures: temp-dir repos exec bin/bench.sh ]  ──▶ pass/fail + wall time
                        ◀ tests attach here: suite green, subtest count unchanged,
                          wall time drops from ~19s to roughly CPU-bound

Seam 3 — the conformance sweep + canary bite-proof (grades the dedup):

    trigger: gate phase 1 (TestRootConformance → checkGoCore)
        │
        ▼
    graded root ──▶  [ go build / vet / go test <all pkgs minus internal/contract> ]  ──▶ diagnostics
                        ◀ tests attach here: `go-test-failing` canary fixture still bites;
                          conformance phase no longer executes contract tests

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1–4 | all seven files at or under 400 lines, split by responsibility | structure check | already red: `bench structure` exits 1 with 7 FILE TOO LONG issues (observed 2026-07-04) | the check greens only when every file clears the budget |
| 1–4 | no test lost in relocation | contract suite + `go test ./...` | not TDD-able as a new red test; parity check: `go test -count=1 -v` RUN-line counts equal before/after, per package touched | a dropped or renamed-away test changes the count; a broken relocation fails compilation or the suite |
| 5 | conformance sweep no longer runs `internal/contract` | conformance sweep | not TDD-able (removal of duplicate work); measured: gate phase-1 wall time drops ~19s | the only ~19s item in phase 1 is the duplicate contract run |
| 5 | unit-sweep bite preserved for every non-contract package | canary bite-proof | already covered: `go-test-failing` fixture expects `checkGoCore` to bite | if the exclusion over-filters, the fixture's failing test goes unseen and the canary reports did-not-bite |
| 6 | contract suite green under parallelism, no data races | contract suite | not TDD-able (performance/concurrency); one-off `go test -race -count=1 ./internal/contract` green + wall time ≤ ~8s | `-race` surfaces shared-state violations the plain gate cannot see |
| 7 | speedup and parity are reported with numbers | final check report | n/a — reporting obligation, not a test | the reviewer accepts on evidence, not on a claim |

### Edge inventory

- Split files: helper visibility — a family file using another family's private
  helper fails compilation; covered by the gate's build. No coverage row needed.
- Split files: re-run idempotency / partial state — n/a, pure source relocation.
- Dedup: graded root without `go.mod` — `checkGoCore` already returns early;
  unaffected. Already covered by existing behavior.
- Dedup: graded tree lacking `internal/contract` (all canary fixtures) — the
  in-Go filter skips a package that isn't listed; exercised by every canary run.
- Dedup: contract package later renamed — the exclusion goes stale and the suite
  silently runs twice again (a cost, not a coverage loss). **Won't handle** —
  self-announcing via gate wall time, and the cross-reference comment marks it.
- Parallel: `t.Setenv` + `t.Parallel` — Go panics on the combination, so a missed
  quarantine is loud, not silent. Covered by the suite itself.
- Parallel: CPU oversubscription against the canary sweep — gate phases run
  sequentially, so the suites never overlap. **Won't handle** — becomes relevant
  only if phase-level concurrency (out of scope) ever lands.
- Parallel: tests exercising concurrency internally (worktree concurrent-acquire)
  — already self-contained per fixture; the `-race` pass is the check.

## Out of scope

- **Phase-level concurrency in `gate.sh`** — a separate capability (gate
  orchestration, not test layout); worth remeasuring only after this spec lands,
  since the sequential gate is already ~25s then. Estimate: ~10 edits, ~6 gate
  runs. Parked on the roadmap.
