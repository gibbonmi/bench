# Review — ft91-phase-manifest-dag

Three-axis semantic review of `7eecc1f..bd761db` (slices A–D), run after the build
went gate-green. Advisory: the gate is still the oracle.

## Standards

4 actionable findings. **Worst: the manifest loader re-derives the package's own
strict-JSON decoder and diverges from it, so a duplicate JSON key silently shadows a
phase and the gate reports green.**

- **S1 — `parseManifest` re-derives `strictJSON` and omits its duplicate-key
  rejection.** `internal/gate/manifest.go:91-104` builds its own `json.Decoder` with
  `DisallowUnknownFields()` plus a `decoder.More()` trailing-content check.
  `internal/gate/verdict.go:251` (`strictJSON`, same package, no import needed)
  already owns all of that *and* recursive duplicate-key rejection via
  `rejectDuplicateNames`. AGENTS.md, "Code standard — one source per fact": "two
  derivations of the same fact … must collapse to one source". Confirmed
  behavioral consequence, not just duplication — a manifest carrying
  `{"name":"real","argv":[…],"name":"shadow","argv":[…]}` runs `shadow` and exits
  0. That is the silent-drop class the loader exists to refuse. Route `manifestDoc`
  through `strictJSON` and map its error onto the existing `unknown field` /
  `parse error` classes.
- **S3 — `Phase.Dir` has two anchoring authorities.** `internal/gate/manifest.go:156`
  anchors a manifest `dir` to the graded root and always emits an absolute path;
  `internal/gate/runner.go:347-352` independently re-implements relative joining
  against the *runner's* root (which `PhasesCommand` passes as `kit`); the `Phase`
  doc comment in `phases.go` states a third version of the rule. No production
  producer emits a relative `Dir`, so the runner branch is kept alive only by
  `runner_test.go:263`. Make `Dir` contractually absolute-or-empty, delete the
  runner branch, and let the doc comment carry the one rule. (Same defect as Spec
  finding Sp3.)
- **S2 — two duplicated scheduler tests.**
  `TestRunnerNeededPhaseCompletesBeforeDependentsStart`
  (`runner_serial_test.go:236`) and `TestSchedulerRespectsNeeds` (`:371`) are the
  same test with different marker and phase names; `TestRunnerNeededPhaseRedSkipsDependents`
  (`:260`) is wholly subsumed by `TestSchedulerSkipsDependentsOfRed` (`:414`). The
  AGENTS.md independent-expectation exception does not cover a second copy of the
  same expectation. Drop the renamed-from-`Serial` originals where the scheduler
  test covers them; keep `…NotFirstStillRunsFirst` and the inner-mode pair, which
  carry distinct facts.
- **S5 — shared `needsBuild` slice aliased into four phases.**
  `internal/gate/phases.go:64-95` hands the same backing array to all four
  downstream phases, so mutating one phase's `Needs` mutates the others. Construct
  per phase.

Not actionable, recorded for the reviewer: `dedupe` (`manifest.go:229`) has no
observable effect — `edgeState` and `firstUnsettledNeed` are already
duplicate-tolerant and no diagnostic renders `Needs`. It implements the spec's
veto item (i) literally, so it is defensible dead code rather than a defect.

## Spec

3 findings, 27 of 29 coverage rows realized and 2 partial. **Worst: coverage row 7's
red signal is falsely classified, which makes its mapped test vacuous — a defect in
the spec, not in the build.**

- **Sp1 — row 7's stated red signal cannot occur.** The row says "append semantics
  hand the child two values for one key". `os/exec` dedups `cmd.Env` before exec,
  keeping the last entry, so `TestRunnerPhaseEnvStripsThenSets`
  (`internal/gate/runner_test.go:301`) passes identically under the old
  `append(gateEnv(), phase.Env...)`. Confirmed by standalone probe. Consequence:
  the unmapped `TestMergeEnvStripsThenSets` (`runner_test.go:285`) is the only
  non-vacuous coverage story 7 has — it is not scope creep. **Fixing this means
  editing the spec's coverage map, which is the reviewer's call, not the build's.**
- **Sp3 — row 6's mapped test grades the dead branch.**
  `TestRunnerPhaseDirIsRelativeToRoot` exercises the runner's relative-join branch,
  which production never reaches. The real semantic — graded-root anchoring with
  `BENCH_KIT` differing, and the no-`dir` default — is covered only by the unmapped
  `TestManifestDirResolvesAgainstGradedRoot` (`manifest_test.go:311`). Behavior is
  correct; the map points at the wrong test. Resolving S3 should carry the map row
  with it.
- **Sp4 (low) — a comment asserts a guarantee the code does not give.**
  `internal/gate/runner.go:258` claims the post-loop naming "keeps a never-launched
  phase from reporting the zero value"; when the blocking need settled 130,
  `firstUnsettledNeed` returns `""` and the result reads green. Closed by Coverage
  finding C1/C4.

Judged clean: story 12 added no second timer and no bounds constant, exit 124 is
unchanged, `internal/bounds` is untouched; nothing from the "Out of scope" section
was built; `Serial` and `splitSerialPhases` are gone repo-wide. The
`decoder.More()` trailing-content guard is defensible under the spec's own
parse-error class — and S1 supersedes it by routing through `strictJSON`, which
already rejects trailing content.

Story 11's green-path summary reordering (`phase build: green` now prints after the
downstream phases' output rather than before) is authorized by the spec's veto item
(e), which unifies summary printing to "declaration order after all phases settle".
Recorded, not a finding.

## Coverage

4 findings, all confirmed by running the code. **Worst: a phase graph that executes
nothing at all reports `gate: green`.**

- **C1 — a cycle reaching the scheduler produces `gate: green`, exit 0.** With
  `a needs b` and `b needs a`, both phases settle as `SkippedBy`, whose `Code` is 0,
  so neither `runner.go:126` nor `runner.go:152` sees red: outer mode prints two
  `skipped (needs …)` lines and then `gate: green`; inner mode prints a bare
  `gate: green` having run nothing. The loader refuses cycles, but
  `benchkitPhasesForCommand` and any injected table bypass it. Nothing asserts that
  a table can never report green with zero phases executed. Nearest test,
  `TestSchedulerSkipsDependentsOfRed` (`runner_serial_test.go:414`), only covers a
  skip caused by a red need, which does redden the run. Fail closed instead.
- **C2 — a phase that exits 130 on its own is misread as cancellation.**
  `runner.go:252` sets `cancelled` on any result code 130. Probed: a manifest phase
  running `exit 130` alongside a green phase yields rc 130, `gate: cancelled; still
  running: self130`, no summaries and no `gate: red`. Per `story4_proof_test.go:317`
  a 130 records the run as `Pending`, so a genuinely red gate leaves no red verdict.
  This is pre-existing behavior that the rewrite carried forward, now reachable from
  a project-authored manifest. Gate `cancelled` on `ctx.Err() != nil` rather than on
  the exit code.
- **C3 — an optional phase whose `dir` does not exist is silently skipped.** chdir
  fails ENOENT and `runPhase`'s `errors.Is(run.StartErr, os.ErrNotExist)`
  (`runner.go:362`) cannot tell a missing binary from a missing directory, so the
  summary reads `skipped (not installed)` and the gate stays green. Probed. A
  manifest typo therefore removes a check with no red. `containedDir`
  (`manifest.go:193`) validates `dir` only lexically, so nothing else catches it.
  Nearest test, `TestRunnerPhaseDirIsRelativeToRoot` (`runner_test.go:255`), uses an
  existing subdir only.
- **C4 — an unlaunched phase can print as green.** `firstUnsettledNeed`
  (`runner.go:289`) returns `""` when the blocking need has settled non-green,
  yielding a zero-value `phaseResult` that renders `phase <name>: green` for a phase
  that never ran. Not reachable through `PhasesCommand` today (its context carries no
  `errGateTimeout` cause), but it is the same fail-open shape as C1 and closes with
  it.

Probed and found safe: no deadlock or hot spin on a cycle, and `results[i]` writes
are correctly published — `settled[i]` is set only after the `<-done` receive.
