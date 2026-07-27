# Review — ft91-gate-phase-split

Three-axis semantic review of the branch diff (57 files, base `9399ba4`). Advisory:
the gate was green when this was written. Each axis ran as its own read-only
delegate so none polluted another's context.

Findings accepted for repair are marked **fix**; findings the reviewer is asked to
accept as residual risk are marked **accept**.

## Standards — 6 findings (4 hard, 2 judgment)

Worst: the phase-name inventory is derived twice, which is the one fact this split
most needs single-sourced, and a comment in the diff claims it already is.

- **fix** — `internal/canary/canary.go`, `phaseFamilies`: a hand-typed copy of the
  phase names `internal/gate/phases.go` produces. `canary` cannot import `gate`
  (the dependency runs the other way), so the two can only be kept in sync by hand.
  AGENTS.md: "Knowledge duplication is a defect: two derivations of the same fact …
  must collapse to one source." The comment in `internal/conformance/registry_test.go`
  asserting "Asking the router keeps the phase names in one place instead of listing
  them again here" is false while the router is itself the second place.
- **fix** — comments narrate the change rather than the code, in eight places:
  `internal/conformance/package_core_checks_test.go` (the residual's doc, and
  `TestResidualCheckBuildsNothing`), `package_core_diagnostics_test.go`,
  `fixture_bite_test.go`, `internal/gate/gate_go.go` (`gofmtStep`),
  `internal/preprelease/preprelease.go` and `preprelease_test.go`, and
  `internal/conformance/registry_test.go`. craft-comments: "the comment never
  mentions [the change] … What survives merge is only what's true of the code as it
  stands." One of them also names a symbol the tree no longer has.
- **fix** — `internal/gate/gate_go.go`, `GateGoCommand` doc: "Exit 0 is a green step,
  1 a red one, and 2 a usage error" omits the 3 the function returns for a
  not-in-repo root.
- **fix** — `internal/gate/gate_go.go`, `gofmtStep` writes its red verdict to stdout
  while every other red in the file writes to stderr, with no stated reason.
- **fix** — `internal/conformance/gate_entry_test.go`: `var entryTier = registry.TierFor`
  is a second name for one function and nothing else.

## Spec — 3 findings

Worst: story 5's substitute probe reintroduces the consumer-facing overreach the
story exists to remove.

- **fix** — `internal/gate/phases.go`: the `conformance-suite` probe tests
  `isDir(root/internal/conformance)` — a directory *name*, which is the kit-specific
  thing itself. Any linked Go repo with a package at that path materializes the phase
  and gets the filtered suite run against it. Story 4's twin probe is honest by
  contrast: it scans for a test-name declaration no consumer will collide with.
  Story 5 is narrowed, not closed.
- **fix** — `internal/gate/phases.go`, `BenchkitPhases` doc still says "Conformance's
  own build check goes to a throwaway path for the same reason", describing the path
  story 6 deleted.
- **accept** — `internal/conformance/cross_compile_stress_test.go` is `//go:build
  stress`, and the only `-tags stress` run over that package is prep-release's
  `conformance-ship`, which filters to `^TestRootConformance$`. So
  `TestResidualCheckKeepsCrossCompile` runs only by hand. The coverage row's red
  signal literally says "observed red … before build", so the row is satisfied as
  written, and cross-compile itself is still graded at ship through
  `conformance-ship` → `checkGoToolchain` → `crossCompileMatrix`. Recorded as
  residual risk rather than repaired, because making it run in the dev tier means
  paying the stress matrix on every commit.

**37-row audit: 34 genuinely delivered.** Row 9 and row 15's fixture-local-manifest
half are unsatisfiable under the reviewer-approved drop of story 9 and are superseded
by the probes; row 22 is the non-gate-assertable measurement, reported separately.
No "Won't handle" or "Out of scope" item was quietly built.

## Coverage — 7 findings

Worst: the `race` fixture's EXPECT is emitted by the `test` phase too, so deleting the
`race` phase would leave the fixture still red — the one failure mode canaries exist
to catch.

- **fix** — `tests/canary/race/race-cleanup-test-failing/EXPECT` is
  `--- FAIL: TestConcurrentCleanupRecordsOneTransaction`. Remove the `race` entry from
  the phase table and the owner is absent, so `phasesForMode` falls back to the full
  inner gate and the `test` phase's plain `go test` emits that exact string. The build
  fixture was deliberately designed against this; the race fixture was not. The tree
  should carry a genuine data race that only `-race` detects.
- **fix** — `internal/gate/phases.go`, `declaresCleanupRaceTest` scans bytes, not
  syntax: a commented-out or string-literal occurrence of the test name materializes a
  phase whose `-run` then matches nothing, and the did-it-run guard reds a repo that
  declares no such test — the precise harm the probe exists to prevent.
- **fix** — `internal/gate/gate_go.go`, `runStep` discards spawn errors. With `go` off
  PATH, `gate-go race` returns 1 having written nothing to either stream: a red phase
  with no diagnostic. `coreTestStep` prints one; `race` and `conformance-suite` do not.
- **fix** — `internal/gate/gate_go.go` exit 3 and `toon.NotInRepo()` have no test
  anywhere; neither the absent-root nor the empty-root argument form is exercised.
- **fix** — `internal/gate/gate_go.go` parses `gofmt -l` output with `strings.Fields`,
  so an unformatted file named `bad file.go` is reported as two files that do not
  exist. Line-based parsing fixes this and still covers the spec's recorded
  won't-handle (a last line lacking a trailing newline).
- **fix** — `internal/canary/canary.go`, `assertFamilyBinding` fails open when
  `BENCH_KIT` is unset, empty, or relative: the kit then skips its own binding check.
  The tests always set it, so the case where the assertion silently does not run is
  the untested one.
- **accept** — the `build` and `vet` fixture EXPECTs quote Go toolchain wording, so a
  Go release that rewords them breaks the fixtures. It breaks loudly ("did not bite"),
  which is the failure mode a canary is built to survive.
