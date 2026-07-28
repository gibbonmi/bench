# Review — ft91-canary-compiled-bites

Reviewed commit `a92cae5` ("canary: grade behavior-owned fixtures by compiled contract
test binaries") against `specs/ft91-canary-compiled-bites.md`. Base `e079dcf`. The build
sits on its assignment branch and has not landed on the default branch.

Three axes, kept separate. Findings marked **blocking** should be resolved before the
merge; the rest are the reviewer's call.

## Standards

**9 findings — 4 hard violations, 5 judgment calls. Worst: the `./internal/contract/`
prefix is now derived three ways.**

### Hard violations

1. **One-source-per-fact regression on the contract import prefix.** `AGENTS.md`, "Code
   standard — one source per fact": *"two derivations of the same fact … must collapse to
   one source."* The path is now spelled independently in three places:
   `internal/canary/canary.go:1032` (`const contractPackagePrefix`),
   `internal/canary/canary.go:487` (`filepath.Join(root, "internal", "contract", …)`), and
   `internal/gate/phases.go:35` (`const contractSubtree = "./internal/contract/..."`, an
   inlined literal where the pre-change file composed it from the prefix const).
   Pre-change the tree carried two derivations, so this is a net +1 and an un-composing of
   the gate side. The comment at `canary.go:1030-1031` asserting the two must agree is the
   tell. `internal/gate` already imports `internal/canary`, so one exported prefix collapses
   all three.

2. **`projects/benchkit.md` left describing the retired run shape.** `.bench/BENCH.md`
   invariant 3: docs *"describe the current decided state."* `projects/benchkit.md:178-179`
   — *"a nested fixture exercises only its owning phase (conformance for conformance
   families, contract for `behavior-owned`)"*; `:214-215` — *"Each nested fixture keeps the
   real gate entry path but runs only the phase that owns its failure"*; `:219-220` —
   *"every inner gate is invoked with an explicit `GOMAXPROCS`"* (the pin now also rides
   compiled bites, `canary.go:956`). All false for the behavior-owned family after this
   diff. Not gate-enforced: the docs-currency check grades only named anchors, none of them
   this prose. Same object as Spec S4.

3. **Comment register — change-narration and PR-talk that survives the merge.**
   `craft-comments`, The register: *"no narration … no argument for its own correctness.
   That prose addresses the reviewer of the diff."* Instances:
   `internal/canary/compiled_bite_test.go:16, 75, 92-94, 172, 195`;
   `internal/canary/contract_scope_test.go:260`;
   `internal/gate/contract_phase_argv_test.go:29-33`; and, weakest,
   `internal/canary/canary.go:1009-1011` (first clause is a legitimate why, the rest argues
   suite adequacy). "the slice", "the migration", "the existing test" resolve only against
   this diff.

4. **Comments left describing the code that was, in files the diff edited.**
   `craft-comments`, Aging: *"read the comment as part of the diff: update it or delete
   it."* `internal/canary/canary.go:29-30` (`PhaseEnv`: *"An absent value means the fixture
   must exercise the full inner gate"*); `:131-133` (`Fixture`: *"the family that routes its
   inner gate"*); `:351-353` (`selected`: *"everything its inner run is scoped by"* — `pkg`
   no longer scopes a run, it selects a run kind).

### Judgment calls

5. **`FixturePhase` is a middle man returning a phase nothing pins.** `canary.go:43-53`
   returns `PhaseContract` for `behavior-owned`, but `subjectCall` (`:919`) returns before
   the phase pin whenever `fixture.pkg != ""`. The value now survives only as an equality
   token for "is this the behavior-owned family" (`:466`, `:758`).
   `internal/gate/phases_command_test.go:218` still asserts a combination nothing produces.

6. **Duplicated test, inert by its own admission.**
   `internal/gate/contract_phase_argv_test.go:34-43` is the body of `:19-27` plus a
   `t.Setenv` of a name no code reads; its comment concedes the export is *"inert by
   construction"*, so the added line cannot change the outcome.

7. **`SubjectRootEnv`'s stated placement rationale is false, and the placement has an
   unstated cost.** `canary.go:34-36` claims *"this package is the only one all three can
   read."* `internal/bounds` is stdlib-only in production and already imported by
   `internal/gate`, `internal/contract`, and `internal/canary`; `internal/capability` owns
   an exactly analogous env-name constant (`LogEnv`). Verified cost:
   `go list -deps -test ./internal/contract/axi` now returns `internal/canary` (plus `git`,
   `subprocess`, `toon`) and did **not** pre-change — every contract test binary the sweep
   compiles now links the package that compiles it. Bears on the story-8 named deviation.

8. **Implicit zero value for the gate call kind.** `canary.go:931` builds the gate call with
   no `Kind`, relying on `RunGate` being `iota`'s zero, and `runnerCommand` (`:1022`)
   handles it via `default:` rather than `case RunGate:`. The test spells it explicitly
   (`compiled_bite_test.go:45`). A fourth kind added later falls through to
   `exec.Command("bash", call.Gate)` with an empty `Gate` — fail-open dispatch.

9. **Half the strip list is literals, half constants.** `canary.go:970-980` comments that
   the strip and the set must agree, then spells `"GOMAXPROCS"` (`:978`) against
   `innerWidthPin`'s `fmt.Sprintf("GOMAXPROCS=%d", …)` (`:963`), and `"BENCH_CANARY_INNER"`
   (`:974`) against `gateEnv`'s literal (`:1002`) and `internal/gate/phases.go:220`. The
   other four entries in the same list are single-sourced. Minor, folded here: the renamed
   `canary.gateEnv` (`:1000`) now collides by name with `gate.gateEnv`
   (`internal/gate/gate.go:147`), which performs a different strip.

## Spec

**5 findings — 1 blocking, 1 spec-authoring defect, 2 documentation/classification, 1
minor. Worst: S1.** 28 coverage rows audited; 24 present, verified against their objects.

1. **S1 — blocking. Vacuity baselines for every *non*-contract family were silently
   narrowed to a single phase.** Spec, Implementation decisions: *"Every other family's run
   shape is untouched"* and *"Bite and vacuity semantics are unchanged."* The diff routes
   baseline calls through `subjectCall` (`canary.go:610` → `:916-932`), whose non-contract
   branch appends the phase pin whenever `FixturePhase(family) != ""`. The deleted
   `baselineEnv` added that pin **only** when `group.pkg != ""` — contract groups — so every
   other group's baseline ran the full inner gate. Verified against `e079dcf`.
   Consequence on this kit: non-contract groups key on `scope` (`canary.go` `group()`), so
   every legacy-flat and phase-named fixture shares the one unscoped group, whose
   representative is the first in directory order — `tests/canary/build/…`. That shared
   baseline now runs the **build phase alone** over the empty tree instead of a full gate,
   so any EXPECT a full empty-tree gate would emit is no longer flagged vacuous. Conformance
   groups are narrowed the same way (`line-routing`: `phases=[]` → `phases=[conformance]`).
   It is also order-dependent: adding a family that sorts before `build` silently changes
   which phase the shared baseline runs. This is exactly the failure story 4 exists to
   prevent, newly introduced for the families the spec declared out of scope, and the real
   sweep stays green throughout.
   Knock-on: the ~11 s gate improvement in `decisions/assets/gate-critical-path-timeline.md`
   is partly attributable to this unrequested narrowing, not solely to compiled bites.

2. **S2 — coverage row 7c cannot go green as written.** The row's own command,
   `rg -n -e narrowContractScope -e ContractPackageEnv -e BENCH_CANARY_CONTRACT_PACKAGE
   --type go`, returns one hit: `internal/gate/contract_phase_argv_test.go:37`. It is
   self-contradictory with row 7b, which *requires* a test exporting the retired name
   literally, because `--type go` includes `_test.go`. The mechanism is fully deleted
   (`narrowContractScope`, `checkContractPackage`, `ContractPackageEnv`, the gate env scrub,
   `internal/gate/contract_scope_test.go`), so this is a spec-authoring defect, not a
   partial deletion. Re-word the row (e.g. add `--glob '!*_test.go'`, or scope the env-name
   grep separately from the two production symbols).

3. **S3 — coverage row 11 is falsely classified.** The row claims its command *"reds today
   with the `session-handoff.md` reference"*. At `e079dcf`,
   `git grep -l BENCH_CANARY_CONTRACT_PACKAGE` returns only `internal/canary/canary.go` and
   the spec itself — the handoff rewrite had already dropped it. The command was green
   before the diff and after, so the row closed no work.

4. **S4 — story 11 is partial.** Story 11: *"every document that names the retired
   narrowing or the retired environment variable updated in the same diff, so no reference
   outlives the mechanism."* `projects/benchkit.md:178-181` still describes behavior-owned
   fixtures as nested and their EXPECT as matching *"inner-gate output"*. The row's literal
   env-var grep does not reach this prose, so the row is green while the story's intent is
   not met. Same object as Standards finding 2.

5. **S5 — minor. Row 12a covers one of two spellings.**
   `TestContractFixtureDeclaringItsOwnPhaseTableIsSwept`
   (`internal/canary/compiled_bite_test.go:403-412`) exercises only
   `files/dot-bench/phases.json`; the deleted refusal covered the literal
   `files/.bench/phases.json` spelling too. Coverage thinness, not a behavior gap.

**No scope creep found beyond S1.**

### Named deviations — all landed as the spec described (reviewer's call, not defects)

| deviation | evidence |
|---|---|
| story 8 edits `internal/contract/helper.go` | `helper.go:13,49-52`; `rg -c BENCH_CONTRACT_ROOT --type go --glob '!*_test.go'` → one file. See Standards finding 7 for the linking cost. |
| compile source resolved against the swept root | `canary.go:535` (`Cwd: root`), `:1014`, resolver shared with the structural refusal at `:487` |
| `GOMAXPROCS` pin kept on bite invocations | `canary.go:954-957`, asserted `contract_scope_test.go:156` |
| story 12 `.bench/phases.json` refusal removed | `declaresPhaseManifest`, its `assertContractScopes` case, and `dotEncodedPath` all deleted |

## Coverage

**5 findings. Worst: C1 — the vacuity guard for the behavior-owned family grades nothing.**

1. **C1 — blocking. A contract group's baseline exits zero printing `PASS`, which the
   empty-output refusal accepts.** Reproduced against the post-change worktree: compiled
   `axi.test`, invoked over an empty git-init'd tree the way `scopeBaselines` does.
   Under a gate run (`BENCH_SKIP_LOG` set — see C2): **exit 0, stdout `PASS\n`, 5 bytes**,
   34 skip lines diverted to the log. Standalone `bench canary .`: exit 0, 2317 bytes of
   `bench-skip kind=environment` lines. Every `axi` test is guarded by
   `SkipIfSubjectBenchMissing` (`internal/contract/helper.go:246`) and an empty tree has no
   `bin/bench.sh`, so all 34 skip. `scopeBaselines` (`canary.go:608-624`) discards
   `ExitCode` and refuses only `outputs[idx] == ""`, so both shapes are accepted and
   `strings.Contains` can never match an EXPECT. **No behavior-owned fixture can ever be
   flagged vacuous** — a fixture whose EXPECT is a string its contract test prints
   unconditionally reports a bite forever.
   Note on scope: the underlying cause is partly inherited — on an empty tree the contract
   tests skip either way, so a pre-change baseline could not contain a contract failure
   message either. What this diff changes is that it narrows the baseline from gate framing
   plus phase output down to `PASS`, while story 4 claims the comparison *"keeps the meaning
   stage 1 gave it"*. That claim is what is unsupported.
   Row that should exist: "a contract group whose baseline exits zero is a red naming the
   group" — seam A, fake runner returning `RunResult{ExitCode: 0, Output: "PASS\n"}`. The
   existing row's test feeds `ExitCode: 1, Output: ""`, a shape the real baseline never
   produces.

2. **C2 — a bite invocation inherits `BENCH_SKIP_LOG`, so fixture-tree skips append to the
   outer gate run's capability tally.** `capability.LogEnv` is absent from `sweepEnvKeys`
   (`canary.go:971-980`) and unreferenced anywhere in `canary.go`; `withSkipLog`
   (`internal/gate/capability_skips.go:56-63`) puts it on every phase including canary.
   Measured: one baseline invocation appended **34 `kind=environment` lines**; the family
   makes 38 invocations (33 fixtures + 5 baselines). On a host where a `capability`-kind
   skip fires (`surface`'s `WriteUnreadable` → `capability.Privilege` under root,
   `WriteFifo` → `capability.Fifo`), those reach `strictFailure` and red a
   `BENCH_REQUIRE_CAPABILITIES=1` run. The pre-change path went through
   `bash gate.sh` → `gateEnv()`, which strips `capability.LogEnv` precisely for this reason,
   with a comment naming the hazard; the compiled path routes around it. Unmentioned in the
   spec and in `decisions/gate-critical-path.md`.
   Fix: add `capability.LogEnv` to `sweepEnvKeys` and to the absent-key list in
   `TestContractFixtureBiteCarriesOnlyItsSubjectRoot` (`contract_scope_test.go`).

3. **C3 — the compile call's environment is decided by nobody and asserted by nothing.**
   `compileContractPackages` passes `Env: base` (`canary.go:535`); `base` has `GOMAXPROCS`
   stripped and not re-set, so each `go test -c` runs at full machine width, and
   `eachIndex` gives 5 concurrent compiles on a 16-core box — 5 × 16-way against the 8 × 2-way
   budget every other call kind is held to. The spec's **Environment** section enumerates
   what a gate spawn and a bite invocation get and never says what a compile gets, so this is
   undecided rather than chosen. No test touches `RunCompile` `Env`;
   `TestDefaultRunnerDispatchesOnCallKind` asserts `cmd.Args` and `cmd.Dir` only.

4. **C4 — a bite invocation has no test timeout; the path it replaced had ten minutes.**
   `runnerCommand` builds `exec.Command(call.Binary)` with no `-test.timeout`
   (`canary.go:1018`); a compiled test binary run directly defaults to `0` — no timeout. The
   replaced path ran `go test` (`internal/gate/phases.go:191`), where `cmd/go` injects
   `-test.timeout=10m0s`. A contract test that deadlocks against a mutated fixture tree now
   hangs the worker; standalone `bench canary .` has no outer bound at all, and inside a gate
   only `bounds.GateTimeout` (45 min) catches it, killing the run rather than reporting the
   fixture red with a goroutine dump. Partially mitigated by `bounds.TestDeadline` on most
   contract waits, but the backstop is gone.

5. **C5 — `os.MkdirTemp("", …)` honors an ambient `TMPDIR` the sweep neither strips nor
   pins.** With `TMPDIR` inside the swept tree, five ~5 MB test binaries land in the tree the
   gate grades for cleanliness — the "sweep becomes a git-status change" failure story 6's
   row names. `TestCompiledBinariesShareOneSweepOwnedDirectory` asserts the parent is not
   under `root` but inherits the ambient `TMPDIR`, so it cannot see this. Fix: `t.Setenv`
   the variable in that test, or add `TMPDIR` to `sweepEnvKeys`.

### Classes walked with no finding

`go` missing from PATH (lands on story 5's compile red as the spec says); dangling symlink
or special file at `internal/contract/<pkg>`; stale binary from a prior sweep; a group with
zero fixtures (unreachable); binary removed between compile and invoke; concurrent workers
sharing the binaries map and the env backing array; binary-name collision (injective, and
`TestContractBinaryNameIsInjective` covers it); cleanup on the compile-error and
did-not-bite returns; paths with spaces, re-run idempotency, interrupt mid-sweep (spec marks
won't-handle).

## Unrelated open item

`decisions/assets/gate-critical-path-timeline.md:42-49` records **one unreproduced red** on
the `test` phase from the first full gate run against this build, output not captured, three
subsequent runs green. Not attributed to this diff and not ruled out. Flagged for the
reviewer; not a finding on any axis.
