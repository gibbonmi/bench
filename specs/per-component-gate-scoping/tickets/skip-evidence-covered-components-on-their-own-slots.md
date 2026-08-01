# Skip evidence-covered components on their own slots

Blocked by: Build the kit-shaped gate fixture root; Declare shellcheck's and
canary's input sets; Identify a component by its inputs and execution closure;
Record a component ancestor slot as its own class; Carry the partition in the
verdict and refuse its reuse; Render a partial verdict in bench status; Refuse a
partial verdict at prep-release
Ownership fence: new `internal/gate/component_decision.go`, the reduction call
sites in `internal/gate/gate.go`, `internal/gate/reduced_run_test.go`
Assumptions: `reducedInheritance` is today's whole-changeset decision and this
ticket replaces it as the reduction site; `ReducedScope()` stays and the
stripped-worktree construction keeps consuming it; the whole-tree fresh-green
reuse still answers first and `--fresh` never consults the reduction. Re-derive
from the tree at pickup.

## What to build

The feature's centre: one decision function, computed inside the gate run,
answering per component whether its inputs moved. Each of `gofmt`, `vet`,
`test`, `race`, `contract`, `shellcheck`, and `canary` resolves its own ancestor
slot by its own identity and skips only on its own evidence. `conformance` and
`conformance-suite` are unconditional — they enforce the bindings this feature
adds, so they can never be skipped by the surfaces they enforce. `build` still
runs every time in this ticket; its skip arrives next.

Selection is the gate's own. No caller, flag, agent, or session names components
to run or skip: `bench commit` keeps asking unconditionally, and the only
operator lever stays `bench gate --fresh`, which forces more grading and never
less.

Every error return of this one function — slot unreadable, seal unreadable,
derivation failure, identity computation failure, domain mismatch — answers
run-the-component, so fail-closed is a property of one site rather than seven.
A non-kit root and a root with no Go module never scope at all; the declarations
are the kit's own, matched by directory identity rather than spelling.

On green, the run authors a slot for each component it executed and leaves every
skipped component's slot byte-identical. On red, the red component's slot is
invalidated and no other's is touched. Every skip is announced on its own line
with its evidence, and the verdict records the partition.

**Evidence authorship.** `bench gate` authors slots for the components it
executed green; `bench gate --fresh` executes everything and re-authors every
slot; `gate-phases` records no verdict and authors nothing.

## Acceptance

- [ ] PC1 — a capture-only changeset executes `conformance` and `conformance-suite` and skips `gofmt`, `vet`, `test`, `race`, `contract`, `shellcheck`, and `canary` on their own evidence.
- [ ] PC2c — every skipped component is announced on its own line naming its ancestor identity and recorded time, and the recorded verdict carries one entry per skip.
- [ ] PC11 — forcing one component red invalidates that component's slot and leaves every other component's slot bytes unchanged.
- [ ] PC12 — editing only a file under `tests/canary/` runs `canary` and skips the toolchain components.
- [ ] PC13 — editing `bin/bench.sh` runs `canary`.
- [ ] PC14 — editing an ordinary Go source skips `canary` and runs the toolchain components.
- [ ] PC19 — each of slot unreadable, seal unreadable, derivation failure, identity failure, and domain mismatch causes its component to run, one subtest per class.
- [ ] PS29 — `bench gate --fresh` executes every component and re-authors every slot, and the whole-tree fresh-green reuse still answers ahead of the decision.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC1 | answer the partition from `ReducedScope().Confines` for the whole changeset | `TestCaptureOnlyChangesetExecutesConformanceOnly` | seed a green kit-shaped fixture, edit `ROADMAP.md`, execute, compare `.git/phase-runs` against the exact expected executed set |
| PC2c | announce one summary line for all skips | `TestEverySkipIsAnnouncedWithItsEvidence` | execute a partial run, assert one stdout line per skipped component naming its ancestor, and one record entry per skip |
| PC11 | invalidate every slot on any red | `TestRedComponentInvalidatesOnlyItsOwnSlot` | force `vet` red, execute, assert `vet`'s slot absent and `test`'s slot bytes unchanged |
| PC12 | leave `tests/canary/` out of canary's resolved inputs | `TestCanaryFixtureEditRunsCanary` | edit a file under `tests/canary/`, execute, assert `canary` in the executed set and the toolchain components skipped |
| PC13 | declare canary over its two directories only | `TestWrapperEditRunsCanary` | edit `bin/bench.sh`, execute, assert `canary` executed |
| PC14 | add the binary digest to canary's inputs | `TestOrdinaryGoEditSkipsCanary` | edit a Go source, execute, assert `canary` skipped and `gofmt`/`vet`/`test` executed |
| PC19 | return "skip" on any decision-site error | `TestDecisionSiteFailsClosed` | inject each of the five error classes in its own subtest, execute, assert the affected component ran |
| PS29 | consult the partition under `forceRun` | `TestFreshExecutesEveryComponent` | seed a partial green, run with `--fresh`, assert every component executed and every slot re-authored |
