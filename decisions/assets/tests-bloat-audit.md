# tests/ bloat audit — evidence for test-suite-structure #5 and #6

Audited 2026-07-05 (repo age 8 days, 245 commits). Counts at audit time: canary
59 fixtures; conformance 2,474 test lines / 13 files; contract 4,537 lines / 24
files.

## A) Canary fixture anatomy (feeds #5)

**Shape.** A fixture is a directory with exactly two things, enforced by the
runner (`internal/canary/canary.go:126-137`): `EXPECT` (one substring the gate
output must contain) and `files/` (the materialized repo tree; dot-dirs stored
as `dot-` prefixes).

| metric | value |
|---|---|
| fixtures | 59 |
| total files under `tests/canary` | 139 (avg ~2.4/fixture; most are 2) |
| `EXPECT` bytes (all 59) | 2,911 (~49 B each) |
| `files/` bytes (all 59) | 96,926 (~1.6 KB each) / 2,037 lines |

**Boilerplate vs varies: honest-repetition, already lean.** `files/` trees carry
only the artifact under test. Recurring filenames are check inputs, not
boilerplate: `bin/bench.sh` (18×, all 18 distinct hand-written stubs of 9-23
lines — not copies of the 270-line root gate), `dot-bench/BENCH.md` (8×),
`dot-bench/BENCH-reference.md` (7×). Byte-identical duplication is marginal
(~2 identical pairs each of BENCH.md / BENCH-reference.md).

**Runner.** `bench canary [root]` → `Sweep` (canary.go:67) → parallel
`runFixtures` (90, NumCPU workers) → `runFixture` (126). Each fixture: reads
`EXPECT`, rejects a vacuous expectation that also matches an empty-repo
baseline (139-141), materializes `files/` into a temp git repo, runs the real
`.bench/gate.sh`, requires exit≠0 AND output contains EXPECT (153).

**Registry / second-source question.** Two registries exist and are distinct:

1. The check registry = `RunConformance`'s 5 check functions
   (`internal/conformance/checks_test.go:13-21`) — computes diagnostics.
2. A fixture-classification map, `canaryFixtureRegistry`
   (`internal/conformance/registry_test.go:24-91`) — re-declares every fixture
   → owning family → retired shell source; bidirectionally enforced
   (`TestCanaryFixtureRegistryClassifiesEveryFixture`, registry_test.go:101-145).

Generating `files/` from the check registry WOULD create a second source: the
fixture `files/` tree is today the only place the biting input is expressed;
the Go check is the consumer, not a generator. The one honest cross-derivation:
`EXPECT` is a substring of the check's emitted message, bridged (not generated)
by `fixture_bite_test.go:265-266`.

**One-canary-per-check: NOT enforced — manual.** The registry permits many
fixtures per family (`workflow-guidance-anchors` owns ~12,
registry_test.go:52-63). No test asserts a family has ≥1 canary.

## B) Suite-wide leanness (feeds #6)

### Three subprocess-capture seams

| seam | location | lines | output | exit-code extraction |
|---|---|---|---|---|
| conformance `runProbe` | `checks_test.go:114-127` | ~14 | split stdout/stderr | `ProcessState.ExitCode()` |
| contract `runFixtureCommandSpec` | `helper.go:114-135` | ~22 | split stdout/stderr | `errors.As(*ExitError)` |
| canary `defaultRunner` | `canary.go:256-268` | ~13 | merged `CombinedOutput` | `ProcessState.ExitCode()` |

Partly duplicated-knowledge (the "run err → exit code" fallback re-derived 3×
with two idioms for the same fact), partly honest mode difference (canary
deliberately merges output for substring matching; conformance/contract split
to parse stdout in isolation). The mode split is defensible; the two divergent
exit-code idioms are not.

### Repo-scaffolding helpers

At least 3-4 independent "temp dir + `git init -q`" implementations: contract
`NewFixture` + `isolatedEnv` (`helper.go:65-84, 302-325`, full HOME/XDG/
BENCH_HOME sandbox — the only rich version, already single-sourced);
conformance ad-hoc (`checks_test.go:214-221`, `fixture_bite_test.go:264-272`);
canary `gitInit` + `materialize` (canary.go:187-192, 270-274). The git-init
step is duplicated-knowledge, mechanical and low-risk.

### Other N-times-derived facts

- Retired gate-fragment path list in two files: `gate_entry_test.go:22-35`
  (asserted absent from gate.sh) and `registry_test.go` `ShellSources`
  (asserted not to leak EXPECT text). Borderline duplicated-knowledge; drift
  possible.
- Check-family names: concentrated in `registry_test.go` only — single-source.
- Live gate-entry invocations asserted once (`gate_entry_test.go:12-15`) —
  single-source.

### Growth trend

Test-estate line insertions per day: 06-28: 28 · 06-30: 53 · 07-01: 260 ·
07-02: 336 · 07-03: 1,387 · 07-04: 10,696 (contract-subpackage split landing) ·
07-05: 380. Canary fixtures arrived in bursts (5 on 06-30, 11 on 07-02, 7 on
07-03). The test estate is the fastest-growing part of the repo.

## Candidate leanness disciplines (options, not a chosen recommendation)

1. **Unify the subprocess seam** — one `Probe` runner with a merged/split mode.
   Tradeoff: kills the divergent exit-code idiom, but adds a cross-package
   dependency across three currently-independent test packages.
2. **One temp-git-repo helper** — exported `NewRepo(t)` consumed by all suites.
   Tradeoff: removes 3-4 re-derivations, but the suites want different
   env-sandbox depths, so mode flags may reintroduce the branching.
3. **One-canary-per-check meta-check** — assert every check family has ≥1
   registered canary. Tradeoff: converts manual bookkeeping into an enforced
   invariant with zero new source; couples the fixture registry to family names.
4. **Keep canary fixtures hand-tended; do NOT generate from the registry** —
   fixtures are already lean and `files/` is the sole source for biting input.
   Tradeoff: keeps 59 hand-tended dirs and the manual EXPECT↔check-message
   correspondence (bridged by `fixture_bite_test.go`).
