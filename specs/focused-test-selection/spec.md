# focused-test-selection

Status: implemented

Roadmap: FT215

Decision source: specs/focused-test-selection/decisions/ft215-focused-iteration.md (compiled ready map, resolved 2026-08-27; focused-test decisions #1–#9 and #11 govern this spec)

Verification log: 1 iteration(s) to accept — one independent `gpt-5.6-terra`
high-effort round found four blockers. The author folded final-row ownership,
the unsafe-path oracle, the executable trust trace, and exact new-file fences.

Reviewer disposition: 2026-08-27 — approve all three tickets after those four
repairs, with no merge or split. Ticket 01 is the narrower capability that can
ship on its own gate; the complete surface remains serial at `command.go`.

## Problem

`bench test` is the only supported focused test renderer, but it accepts only
one optional package expression. An agent that wants one test, touched
packages, or one conformance check must reconstruct the invocation by hand.
Those reconstructions do not share
the coherent diff subject, the selected Bench executable, the conformance
registry, or one refusal posture.

The tempting correction is a scoped gate. That would give a partial run
authority it has not earned. FT215 instead keeps `bench gate` as the
whole-project oracle and gives focused iteration one explicitly
non-authoritative surface.

## Solution

Extend `bench test` with three focused forms while retaining its current
default and positional-package behavior:

- `bench test [--full] [--package <expr> | <legacy-package> | --changed]
  [--run <go-regex>]`
- `bench test [--full] --check <name>`
- `--changed --base <commit> [--source-tip <commit>]` for explicit live and
  immutable subjects; `--source-tip` without `--base` is a grammar error.

Explicit package selection and `--run` remain one fresh `go test` process.
Changed selection asks `internal/diff` for one coherent path subject and maps
Go-relevant paths through one current `go list` graph. It adds transitive
reverse dependents and invokes the selected packages together. A proven non-Go or
empty subject emits the normal three empty result tables and exits zero. A
Go-relevant path that cannot be mapped completely refuses before tests run.

Named checks validate against the dev conformance registry and use the
existing singular-scope runner with the graded root and dev tier pinned.
Every form selects one Bench executable, keeps the current TOON result
renderer, and honors cancellation. It writes no gate verdict, green marker,
lane record, or reusable evidence.

## Terms

- **Focused run** means a `bench test` result with no oracle authority. Avoid
  “scoped gate” and “partial gate.”
- **Changed subject** means the coherent live or immutable path set returned by
  `internal/diff`, not a second Git snapshot assembled in `testreport`.
- **Selected package** means a runnable current-module package mapped directly
  from a Go-relevant path or reached through a production, test, or
  external-test reverse-import edge.
- **Explicit empty** means successful zero-row `packages`, `failures`, and
  `skips` tables after the resolver proves that no changed path is
  Go-relevant. It is not a silent no-op.

## User stories

1. **Select an explicit package without breaking the existing command.** As an
   agent I can use `--package <expr>` or the compatible positional expression,
   while an absent selection still means `./...`. Existing directory
   normalization and pass-through of `../`, absolute, import-path, wildcard,
   and outside-module expressions remain unchanged. Mutually exclusive
   selectors and misplaced subject flags exit 2 with usage.
   Line: `gpt-5.6-terra` / high / ~2 iterations / serial.
2. **Select one Go test pattern.** As an agent I can combine `--run
   <go-regex>` with the default, explicit-package, or changed-package subject.
   The regex is one argv value passed to Go. A completed invocation that emits
   no test run event refuses instead of reporting a misleading package pass.
   Line: `gpt-5.6-terra` / high / ~2 iterations / serial after story 1.
3. **Derive one coherent changed subject.** As an agent I can select the
   default live review subject, an explicit-base live subject, or an immutable
   base-to-source-tip subject. The live forms include committed, staged,
   tracked-worktree, and untracked paths under one movement check. The
   immutable form never reads unrelated checkout state.
   Line: `gpt-5.6-terra` / high / ~3 iterations / serial after stories 1–2.
4. **Select changed packages.** As an agent I receive one deterministic
   package set. It derives from current Go packages, named embed inputs,
   module-wide Go metadata, and all three import edge classes. Mixed paths take
   the union. Non-Go-only and empty
   subjects return explicit empty output. Deleted packages, deleted embed
   inputs, special nodes, symlinks, and any other Go-relevant path the current
   graph cannot map refuse rather than disappear or widen silently.
   Line: `gpt-5.6-terra` / high / ~3 iterations / serial after story 3.
5. **Run one registered conformance check.** As an agent I can name one
   registry check that runs at the dev tier. The command pins the repository
   root, dev tier, selected scope, kit root, and selected executable; unknown
   or ship-only names refuse before execution. No ambient
   conformance control variable can widen or redirect the run.
   Line: `gpt-5.6-terra` / high / ~2 iterations / serial after stories 1–2.
6. **Keep focused evidence honest.** As a reviewer I see the same package,
   failure, and skip tables for every executed form. Non-Go changes have one
   explicit empty shape, cancellation has one posture, and no focused run
   mutates a gate-owned record. Public help names all supported forms without
   describing them as a verdict.
   Line: `gpt-5.6-terra` / high / ~2 iterations / serial after stories 1–5.

## Implementation decisions

- **One parsed request.** `internal/testreport` parses the grammar into one
  typed request before selecting a child process. Package, changed, and check
  selectors are mutually exclusive. `--base` and `--source-tip` are valid
  only with `--changed`; `--check` is exclusive except for `--full`.
- **Preserve the renderer, separate command mechanics.** Move command parsing
  and child-process lifecycle out of the near-budget `testreport.go` into
  `command.go`. Keep JSON decoding and TOON rendering in `testreport.go`.
  Changed selection belongs in `selection.go`; it is not embedded in the
  renderer.
- **Coherent paths stay in `internal/diff`.** Add one exported, root-bound
  subject API in the existing `range.go` seam. It owns default-base
  resolution, explicit base and tip validation, live movement retry, and path
  collection. `testreport` receives paths plus resolved identity and never
  calls raw Git.
- **One current package graph.** One `go list -json ./...` read supplies
  package directory, import path, `Imports`, `TestImports`,
  `XTestImports`, `EmbedFiles`, `TestEmbedFiles`, and
  `XTestEmbedFiles`. The resolver maps ordinary Go files by containing
  package directory. It maps named embed files by exact repository-relative
  path. The files `go.mod`, `go.sum`, `go.work`, and `go.work.sum`
  are module-wide inputs that select every listed package.
- **Reverse closure includes test edges.** Starting from directly touched
  packages, repeatedly add current-module packages whose production, test, or
  external-test imports name a selected package. Sort and deduplicate the final
  import paths, then pass them to one `go test` invocation.
- **Current-graph refusal is deliberate.** A deleted file inside a surviving
  package still maps by directory. A deleted package or embed input absent from
  the current graph cannot be proved complete and refuses. This spec does not
  materialize or query a baseline checkout.
- **Run matching is observed, not guessed.** With `--run`, the JSON decoder
  records Go's test `run` events. Zero events after a complete process is an
  exit-1 structured refusal. Invalid regexes and ordinary test failures retain
  Go's failure result and the current renderer.
- **Singular conformance transport.** Add one registry-owned scope environment
  name. `TestRootConformance` passes it to the existing
  `RunConformance(..., scope)` path. `bench test --check` validates
  `registry.Find(name)` and `RunsAt(Dev)`, then runs only
  `TestRootConformance` with the scope, root, and dev tier pinned.
- **Environment closure.** Every mode removes ambient conformance root, tier,
  scope, selected-set, inherited-set, and capability-log controls. Check mode
  then installs only its exact controls. The selected run-binary and
  `BENCH_KIT` values continue to come from the one `runbinary` selection.
- **Bootstrap authority before execution.** The running Bench process resolves
  the fixed literal `go` name once from the operator-owned `PATH`; no
  request value can choose that executable. That operator environment is the
  independent trust root already used by the gate's Go phases.

  Before Go starts, `runbinary.ReuseOrOwn` canonicalizes the source root.
  An owned selection is built by the fixed builder and freshness-checked
  against that source. An inherited selection must be a clean absolute,
  non-symlink, regular executable whose embedded freshness record verifies.
  A missing or corrupted inherited selection refuses before `go`.

  Go compiles the requested test package from the command-owned root. Check
  mode then pins that canonical Git root as `BENCH_CONFORMANCE_ROOT`, pins the
  validated selection's source as `BENCH_KIT`, and pins the selected binary.
  The resulting conformance test binary can reach only the existing
  registry-bound check runner. Downstream process hops remain authenticated by
  those existing check owners. N02 proves the common chain; K03 adds a corrupt
  inherited-selection refusal before the conformance child starts.

  No new mode shells out from raw parsed flags. Its typed resolver and refusal
  tests land with the first executable path. No temporary allowlist, broad
  fallback, or advertised not-yet-runnable form is permitted.
- **No oracle projection.** The command does not call `internal/gate` and
  does not write the verdict cache, project-green marker, lane record, or
  evidence store. `bench gate` and landing remain unchanged.

## Testing decisions

- Grammar and explicit-package tests attach at `testreport.Command`, using
  temporary Go modules to inspect the child argv through observable test
  events. Existing positional, selected-binary, malformed-output, and
  cancellation tests stay the compatibility oracle.
- Coherent-subject tests attach at the exported `internal/diff` subject seam.
  They prove default and explicit live composition, immutable base-to-tip
  isolation, one retry then drift refusal, rename-as-delete-plus-add, and no
  Git config or ref mutation.
- Package-selection tests attach below process execution at the typed
  `testreport` resolver. Fixture modules cover all three reverse edge classes,
  named embeds, metadata-wide selection, mixed paths, and deterministic order.
  They also cover surviving-directory deletion, deleted package/embed refusal,
  special files, and symlinks.
- Conformance transport tests attach at `TestRootConformance` and the
  registry scope constant. A command-level test runs one cheap registered
  check and proves timing contains that check only. It also proves ambient
  selection controls are scrubbed and an invalid or ship-only name starts no
  child.
- A record-integrity test snapshots the gate-owned record paths before each
  focused form and proves their absence or bytes are unchanged afterward.
- The ordinary `test` phase is the gate seam. No gate phase or shell entry
  changes.

### Seam diagram

    bench test grammar
           |
           v
    [ typed focused request ]
       /        |         \
      v         v          v
 package/run  changed     named check
      |       subject       |
      |       + graph       v
      |         |      conformance registry
      \________|__________/
               v
      [ one Go test process ]
               |
               v
    package / failure / skip TOON
               |
               x  no gate-owned record writes

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| F01 | 1 | default, positional, and `--package` forms select `./...` or the same normalized expression as today | `TestFocusedRequestPackageCompatibility` plus `TestPackagePattern` | A parser rewrite that drops the legacy operand or changes outside-module pass-through is red before execution. |
| F02 | 1 | package, changed, and check selectors are exclusive, and subject flags outside changed mode or source tip without base exit 2 | `TestFocusedRequestGrammarRefusals` | A flag accidentally treated as a positional or silently ignored is red at the command seam. |
| F03 | 2 | `--run` is one unmodified argv value with default and explicit-package subjects | `TestRunPatternReachesGoAsOneArgument` | Shell splitting, option injection, or omission changes the fixture's observed test set. |
| F04 | 2 | a completed `--run` with zero test run events exits 1, while a matched skip or failure remains an observed run | `TestRunPatternRefusesZeroMatches` | Go's package-level exit zero can no longer masquerade as evidence that a named test ran. |
| C01 | 3 | default changed mode returns committed, staged, tracked-worktree, and untracked paths from one movement-checked live subject | `TestChangedSubjectDefaultLiveComposition` | Omitting any producer or rereading it outside the captured identity changes the exact path set. |
| C02 | 3 | explicit base keeps the complete live subject, while base plus source tip returns only the immutable range | `TestChangedSubjectExplicitLiveAndFrozenPair` | Checkout-only paths leaking into the frozen pair, or disappearing from live mode, are red. |
| C03 | 3 | movement gets one retry and a second drift refuses without emitting a partial selection | `TestChangedSubjectRetriesThenRefusesMovement` | A mixed-revision package list cannot reach execution. |
| C04 | 4 | Go files and exact production/test/external-test embed inputs map to their current packages, then all three reverse-import edge classes close transitively | `TestChangedPackageClosureAcrossAllGoEdges` | Dropping a test-only dependent or treating an arbitrary non-Go file as an embed changes the expected set. |
| C05 | 4 | Go metadata selects every current package and mixed known paths select the deterministic union once | `TestChangedPackageSelectionMetadataAndMixedUnion` | A global graph input cannot be narrowed to one directory or produce duplicate invocations. |
| C06 | 4 | empty and proven non-Go-only subjects emit zero-row package, failure, and skip tables and exit 0 | `TestChangedNonGoSubjectRendersExplicitEmpty` | Returning no bytes or invoking `go test ./...` is red. |
| C07 | 4 | control-byte paths, special nodes, symlinks, and unmappable Go-relevant paths refuse without `go list` or `go test`, while a deleted file in a surviving package still maps | `TestChangedPackageSelectionRefusalMatrix` | Silent omission, unsafe reads, raw control output, and broad fallback all fail the expected no-child cases. |
| C08 | 4 | `--run` applies unchanged to the complete changed-package union | `TestChangedPackageRunPattern` | Dropping the filter or applying it while resolving the graph changes the observed selected tests. |
| K01 | 5 | one known dev check runs through the singular scope and no other conformance check appears in timing | `TestNamedCheckRunsOnlyRegisteredDevScope` | Reusing the gate's ordered partition or falling back to the full tier is visible in the timing names. |
| K02 | 5 | unknown and ship-only names, plus `--check` combined with any selector or `--run`, refuse before child start | `TestNamedCheckRefusalMatrix` | Registry drift cannot become an empty green or an unexpectedly broad run. |
| K03 | 5 | root, dev tier, scope, kit, and selected binary are exact, ambient controls are absent, and a corrupt inherited selection refuses before child start | `TestNamedCheckOwnsConformanceEnvironment` plus `TestNamedCheckRefusesCorruptInheritedSelection` | An inherited redirect or self-authenticating executable changes captured state or starts the forbidden child. |
| N01 | 6 | default, explicit-package, and run-filtered forms leave existing or absent gate-owned records unchanged | `TestExplicitFocusedRunsWriteNoGateOwnedRecords` | Any gate persistence added to the first focused forms becomes a byte-level test failure. |
| N02 | 6 | every form uses the fixed Go hop and one independently validated selected binary, cancels its process group, and retains the current renderer | run-binary, corrupt-selection, cancellation, and decoder tests across typed requests | A raw executable value, second selection, leaked child, or separate rendering path breaks the common chain. |
| N03 | 6 | changed-package forms leave existing or absent gate-owned records unchanged | `TestChangedRunsWriteNoGateOwnedRecords` | A changed subject cannot acquire gate authority through persistence. |
| N04 | 6 | named-check forms leave existing or absent gate-owned records unchanged | `TestNamedChecksWriteNoGateOwnedRecords` | Reusing the gate's conformance entry cannot accidentally reuse its verdict writer. |
| H01 | 6 | public help advertises package, run, changed-subject, and check forms as focused evidence, not a verdict | command-registry help inventory tests | A stale public grammar or “gate” authority claim fails the front-door inventory. |

### Edge inventory

- **Paths with spaces or glob characters:** changed paths stay typed argv/data;
  fixture package and embed paths include both. No shell joins them.
- **Git-sourced control bytes:** ESC and BEL in paths refuse before TOON output.
  Tab, newline, and return exercise the encoder-permitted half and must remain
  one escaped cell or refuse before execution; raw line splitting is red.
- **Flag-looking and hostile regex values:** `--` retains the legacy
  positional route. `--run` and `--package` values beginning with `-`,
  containing spaces, brackets, or shell operators remain one argv value.
- **Empty/absent:** no changed paths and non-Go-only paths are explicit empty.
  Missing package terminals after an attempted run remain an error. Zero
  matched tests differs from a package with no tests when no filter was
  requested.
- **Deleted/renamed:** Git's no-renames projection supplies deletion plus
  addition. Both feed classification. Current-directory mapping admits a
  deleted file in a surviving package; absent package/embed ownership refuses.
- **Special files and symlinks:** changed selection uses lstat-style
  classification before any file read. FIFO, socket, device, dangling symlink,
  and live symlink cases start no `go list` or `go test` child.
- **Malformed identity:** missing, ambiguous, or non-ancestor base/tip values
  retain `internal/diff`'s flag-specific errors and write no Git state.
- **Movement and interruption:** live subjects retry once on movement.
  Cancellation covers both `go list` and `go test` process groups and emits
  no partial package tables.
- **Missing tools and invocation location:** missing `go` is a structured
  start refusal. Invocation below the repository root still resolves and runs
  from the root. Existing real-kit, linked-repo, and symlinked launcher routing
  continues through the command registry.
- **Ambient environment:** every conformance control is scrubbed; check mode
  reinstalls exact values. Capability-log control remains scrubbed as today.
- **No attached stateful classes:** the command writes no project file and
  reads no hand-edited record or path pointer. It prompts for no input,
  performs no destructive worktree operation, parses no Markdown, and
  serializes no cross-process state. The profile's rewrite-self-report, whitespace,
  unterminated-delimiter, prompt, lifecycle, and merge-recomposition classes
  therefore do not attach.
- **Won't handle:** resolving a deleted package or embed from a baseline
  checkout. That requires a second graph at a private immutable tree,
  estimated 8–12 production/test edits and 2 full gate runs.
- **Won't handle:** deriving conformance checks from changed paths. Registry
  `Inputs` remain labels rather than selectors; an executable derivation is
  estimated 10–16 registry, selector, fixture, and test edits plus 2–3 full
  gate runs.

## Ownership fences

One writer at a time. Tickets are serial where their grammar and command
mechanics overlap. Reviewer disposition: approve, merge, or split each fence
during sign-off.

- Stories 1–2:
  `internal/testreport/command.go` (new),
  `internal/testreport/testreport.go`,
  `internal/testreport/testreport_test.go`,
  `internal/testreport/runbinary_test.go`,
  `internal/testreport/cancel_test.go`,
  `cmd/bench/main.go`, and
  `cmd/bench/command_registry_test.go`.
- Stories 3–4:
  `internal/diff/range.go`,
  `internal/diff/explicit_base_test.go`,
  `internal/diff/source_tip_pair_test.go`,
  `internal/testreport/command.go`,
  `internal/testreport/selection.go` (new), and
  `internal/testreport/selection_test.go` (new).
- Story 5:
  `internal/conformance/registry/scope.go` (new),
  `internal/conformance/gate_entry_test.go`,
  `internal/conformance/tier_test.go`,
  `internal/testreport/command.go`,
  `internal/testreport/check_test.go` (new),
  `cmd/bench/main.go`, and
  `cmd/bench/command_registry_test.go`.
- Story 6: `internal/testreport/testreport.go`,
  `internal/testreport/testreport_test.go`,
  `internal/testreport/runbinary_test.go`,
  `internal/testreport/cancel_test.go`,
  `internal/testreport/selection_test.go` (new),
  `internal/testreport/check_test.go` (new),
  `cmd/bench/main.go`, and
  `cmd/bench/command_registry_test.go`. No `internal/gate` production file.
- Spec compilation only: `ROADMAP.md`, `roadmap/FT215.md`,
  `decisions/ft215-focused-iteration.md`,
  `decisions/assets/ft215-focused-test-inputs.md`,
  `decisions/assets/ft215-fast-lane-inputs.md`, and
  `specs/focused-test-selection/decisions/`.
- Every ticket: `specs/focused-test-selection/` for ticket checkboxes,
  review-approved amendments, and lifecycle status changes.

## Out of scope

- The independent path-aware `bench commit` fast lane decided by FT215. Its
  decisions remain in the compiled map and must be re-homed before this spec
  retires. Expected size: 12–20 production/test edits and 2–3 full gate runs.
- A partial gate verdict, changed-input green marker, scoped evidence reuse, or
  any change to the whole-project gate at landing.
- Package-scoped vet, build, gofmt, race, system, or release work.
- A general file-to-test build system, test-result cache, daemon, or watch
  mode.
- Multiple explicit `--package` values or a public package-set expression
  beyond the compatible one-expression form.
- The spec-build checkpoint lifecycle in
  `decisions/spec-build-review-gate-cadence.md`.
