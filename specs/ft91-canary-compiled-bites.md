# ft91-canary-compiled-bites

Status: implemented

Compiled from `decisions/gate-critical-path.md` (#7, resolved 2026-07-28) — FT91's
eighth arm, stage 2. Stage 1 (`ft91-canary-contract-scoping`) shipped and is retired;
its subfamily binding is this slice's input. The map's #2/#3 stay open and gate only
the artifact-hoist slice that follows this one, not this spec.

## Problem

The dev gate is canary-bound end to end: gate wall 267 s against a solo canary of
250 s. The 33 `behavior-owned` fixtures are the bulk of it. Each one spawns a nested
`bench gate` over its mutated tree, and even scoped by stage 1 to the single contract
package that owns its EXPECT, that nested gate re-pays process startup, phase-table
construction, and a fresh `go test` compile-and-link of that package — per fixture,
at an inner width of 2.

Nothing about the nested gate is load-bearing for what the fixture proves. A
behavior-owned fixture's EXPECT is a contract test's own failure message; the gate
around it only carries that message upward. The sweep is paying for a process tree to
observe a string that the test binary already prints.

## Solution

Stop spawning nested gates for the `behavior-owned` family. Compile the owning
contract package's test binary once per package group with `go test -c`, then invoke
that binary once per fixture root with the contract subject root pointed at the
fixture's materialized tree. The canary sweep still owns bite, did-not-bite, and the
vacuity baseline, so `bench canary .` stays standalone-provable and vacuity stays a
canary concept.

Because no behavior-owned fixture spawns an inner gate any more, stage 1's inner-mode
contract-phase narrowing has no caller. It goes, with its environment variable and its
now-unreachable tests.

Measured feasibility (prototype, 2026-07-28, this tree): `go test -c ./internal/contract/axi`
compiles in 0.5 s warm, and the resulting binary run against the materialized
`roadmap-regressed` fixture root reds in 0.49 s with its EXPECT
(`missing "no ROADMAP.md"`) present in test-level output.

## User stories

1. As the reviewer waiting on the gate, I want a `behavior-owned` fixture's bite proven
   by invoking its owning contract package's compiled test binary rather than by
   spawning a nested gate, so the sweep stops paying a process tree per fixture.
   Line: gpt-5.6-terra / medium. This is the oracle's own run shape, and the profile's
   cached routing sends gate and conformance logic to the mid tier because a wrong gate
   is the worst class of bug in this kit.

2. As a fixture author, I want each contract package group's test binary compiled
   exactly once per sweep and reused by every fixture in that group, so compile and link
   are paid once per package instead of once per fixture.
   Line: gpt-5.6-terra / medium. The reuse is the whole win of the slice and it sits in
   the same sweep logic as story 1.

3. As the reviewer, I want every bite invocation to carry its own fixture's materialized
   tree as the contract subject root, so one fixture's mutation can never be graded
   against another fixture's tree.
   Line: gpt-5.6-terra / medium. A subject root that leaked between concurrent
   invocations would produce a green that means nothing, which is oracle logic.

4. As the reviewer, I want each contract group's baseline to run the identical
   compile-and-invoke shape against an empty tree, so an EXPECT that collides with the
   infrastructure noise a test binary emits over a tree it cannot grade is screened out.
   That screen is all the empty-tree comparison establishes for this family — vacuity
   checking is not live here, and was not live before either, when the same empty tree
   sat under a phase-narrowed gate — because an empty tree is a degenerate tree, not an
   unmutated twin of any fixture's. What the comparison does and does not establish is
   stated on the baseline comparison in `internal/canary`.
   Line: gpt-5.6-terra / medium. A baseline whose shape differs from its group's fixtures
   emits output no fixture in the group can produce, so the collision screen compares
   against the wrong noise in both directions — matching what the group never prints,
   missing what it does — and a silently wrong screen is oracle logic.

5. As the reviewer, I want a package group whose test binary fails to compile — or whose
   compile produces no binary at all — to red naming that package, never to be skipped.
   Line: gpt-5.6-terra / medium. A skipped group is a family of fixtures that stops
   grading anything while the sweep still reports green.

6. As a session running the sweep, I want every compiled test binary to live under one
   sweep-owned temporary directory that is removed on every exit path, so a sweep that
   errors out leaves no artifacts behind.
   Line: gpt-5.6-terra / medium. This is the map's one open spec-time detail, and cleanup
   on the error path is where the obvious implementation is wrong.

7. As a maintainer, I want the gate's inner-mode contract-phase narrowing, its
   environment variable, and their now-unreachable tests deleted, so no callerless
   scoping path survives to be maintained or re-entered.
   Line: gpt-5.6-terra / medium. Deleting a scoping path from the phase table is gate
   logic even when the edit is mostly subtraction.

8. As a maintainer, I want the contract subject-root environment variable named in
   exactly one place and consumed by the gate phase table, the canary sweep, and the
   contract test helper, so a rename cannot leave one writer disagreeing with the reader.
   Line: gpt-5.6-terra / medium. The sweep becomes a third writer of a name that is
   currently a repeated literal, which the code standard's one-source-per-fact rule
   forbids. **This story edits `internal/contract/helper.go`, which the Handoff's item 1
   says stays untouched except where an EXPECT proves unobservable — a named deviation,
   flagged for reviewer veto.** The fallback if vetoed is that the sweep carries a third
   literal of the same name, which the code standard forbids; that is the trade the veto
   makes.

9. As a fixture author, I want any EXPECT that was observable only in gate-level framing
   fixed by making the owning contract test's own failure message carry the fact, so no
   fixture is retired or exempted to make the migration land.
   Line: gpt-5.6-terra / medium. Judging whether a failure message states the fact well
   is prose the gate can only grade as red or green, not as good.

10. As the reviewer, I want the before and after canary and gate wall-clock recorded in
    `decisions/assets/gate-critical-path-timeline.md`, so the slice ships with evidence
    even though no wall-clock threshold gates it.
    Line: gpt-5.6-terra / medium. This deviates from the profile's doc-authoring row,
    which sends writing to the top tier because guidance prose compounds through every
    session that loads it; a recorded measurement is a fact about one host at one moment
    and compounds through nothing, so the leverage override does not apply.

11. As the next cold session, I want every document that names the retired narrowing or
    the retired environment variable updated in the same diff, so no reference outlives
    the mechanism.
    Line: gpt-5.6-luna / low. This deviates from the profile's doc-authoring row for the
    same reason as story 10 and one more: deleting a reference to a mechanism that no
    longer exists authors no prose at all, and the docs-currency conformance check grades
    the result.

12. As a fixture author, I want the refusal that rejects a `behavior-owned` fixture for
    declaring its own `.bench/phases.json` removed, so the harness stops enforcing a
    constraint whose mechanism no longer exists.
    Line: gpt-5.6-terra / medium. Removing an enforcement is exactly the kind of edit
    that needs a stated rationale rather than a quiet deletion.

## Implementation decisions

**The run shape.** Per contract package group present in the swept tier: one
`go -C <root> test -c -o <binary> ./internal/contract/<package>`, then one invocation
of `<binary>` per fixture root in that group, plus one invocation for the group's
vacuity baseline. No `bench gate` process is spawned for the `behavior-owned` family.
Every other family's run shape is untouched.

**Why compile-once rather than one process for all roots.** `SubjectRoot` reads the
subject-root environment variable per call, so one process cannot serve two roots; and
Go cannot re-enter a package's `Test*` functions as parameterized subtests without a
per-package registry that would duplicate the package's own test list and collide with
the `t.Parallel` those tests already call. Compile-once delivers the same intent — no
nested gate, no process tree, compile paid once per group — and each fixture keeps the
clean process the root swap needs anyway.

**The compile source is the swept root.** The compile runs `go -C <swept root>`, and
the package is resolved against `<swept root>/internal/contract/`, matching the
resolution the sweep's existing structural refusal already performs. The Handoff's item
6 says "a bound package the **kit tree** lacks"; the swept root and the kit tree
coincide for every tree that can carry a `behavior-owned` fixture, because the surviving
structural refusal already resolves the binding against the swept root and reds a linked
repo that has no `internal/contract/` at all. Resolving against the swept root keeps one
resolver rather than two that can disagree. **Named deviation from the Handoff's wording,
flagged for reviewer veto.**

**The `Runner` seam grows a call kind.** `RunCall` gains an explicit kind
discriminator with three values — spawn an inner gate, compile a package's test binary,
invoke a compiled binary over a subject root — plus the fields the two new kinds need
(the package path, the binary path). Exactly one kind's fields are populated per call.
One injected seam rather than a second injected compiler keeps every sweep assertion
readable from one recorded call list, which is how the existing sweep tests are written.

**Invocation working directory.** A bite invocation runs with its working directory set
to the compiled package's source directory, which is where `go test` runs a test binary.
The subject root travels in the environment, not the working directory.

**Environment.** One strip list, two setters. The list of sweep-controlled variables —
the wrapper routing pair, the canary inner-gate marker, the phase pin, the conformance
tier and check, `GOMAXPROCS`, and the contract subject root — is declared once and
stripped from the inherited environment for every call the sweep makes, keeping the
strip-then-set discipline Go's duplicate-key-free exec environment requires. What is
then set differs by call kind:

- A **gate spawn** gets what it gets today: the inner-gate marker, the tier pin, and the
  `GOMAXPROCS` pin, plus the phase and check scope its caller knows.
- A **bite invocation** gets the subject root and the `GOMAXPROCS` pin, and nothing else.
  It carries no inner-gate marker, no phase pin, and no conformance tier or check —
  there is no gate to read any of them, and the Handoff's item 4 makes "no call carrying
  the canary inner-gate mode for this family" an assertable.

The `GOMAXPROCS` pin is the one thing a bite invocation keeps from the inner-gate
environment. The map says the inner-width pin "becomes moot for this family", which is
true of the *nested gate* it was written for and not of the test binary underneath: the
sweep's worker budget still divides the machine by `bounds.CanaryInnerWidth`, so
unpinned binaries running at full width across every worker would restore exactly the
oversubscription the map's #1 measured. Keeping the pin also keeps the budget arithmetic
single-sourced. **Named deviation from the map's wording, flagged for reviewer veto**;
the alternative is dropping the pin and re-deriving the worker budget for a mixed sweep,
which is a bigger change than this slice.

The retired contract-package variable leaves both the strip list and the gate's own
environment scrub, because with no narrowing code left an ambient export is inert by
construction rather than by suppression.

**The real runner is part of the change, not only the call shape.** The sweep's default
runner today builds one command — `bash <gate>` — for every call. It gains a branch per
call kind: compile execs the Go toolchain, and a bite invocation execs the compiled
binary directly. This is called out because every sweep assertion below runs through an
injected fake, and a change that relabels the calls while the real runner still spawns
gates would satisfy all of them.

**Binary location and cleanup.** One temporary directory per sweep holds every compiled
binary, created before the first compile and removed on every exit path including the
error returns, matching how the baseline scratch directories are already handled. Each
binary's file name is derived from its package path by a substitution that cannot make
two distinct package paths collide.

**What the gate loses.** `narrowContractScope`, its package-value validation, the
contract-package environment constant, the gate's environment scrub of that variable,
and the gate-side tests covering all of it. The contract phase's argv becomes the
unnarrowed subtree in every mode. The phase-name constants and the family-to-phase
router stay: the router is still what identifies the `behavior-owned` family as the
contract-bite family, and the gate's phase table still reads its phase names from the
canary package.

**What the sweep keeps.** The walk, the subfamily package binding under
`tests/canary/behavior-owned/<package path>/<fixture>/`, per-package vacuity grouping,
the structural refusals for a fixture bound to no package and for a fixture bound to a
package the tree lacks, global base-name uniqueness, and the shared worker budget.

**What the sweep drops.** The refusal that rejects a `behavior-owned` fixture declaring
its own `.bench/phases.json`. That refusal existed because a declared phase table
replaces the built-in one the narrowing lived in; with no gate spawned for the family,
such a file is an inert artifact in a subject tree. Removing it is not a weakening —
there is no longer a mechanism for it to protect. **Flagged for reviewer veto.**

**Bite and vacuity semantics are unchanged.** A fixture bites when its invocation exits
nonzero and its output contains the EXPECT; an EXPECT already present in its group's
baseline output is vacuous; anything else is a did-not-bite red. Canary and gate exit
codes are unchanged.

## Testing decisions

A good test here drives the sweep through its injected `Runner` and asserts the call
list — kinds, counts, per-call package, binary, subject root, and environment — rather
than reaching into how the sweep sequences its stages. The prior art is exactly this:
`internal/canary/contract_scope_test.go`, `phase_family_test.go`, and `scope_test.go`
all build a fixture tree in a temp dir, sweep it with a recording runner, and assert
over the recorded calls. On the gate side the prior art is
`internal/gate/phases_test.go` and the argv assertions in the file this slice deletes.

The end-to-end proof stays where it already is: the gate's canary phase runs the real
sweep over the kit's own 33 behavior-owned fixtures, so a migration that broke bite or
vacuity for any of them reds the gate. The red-test-to-red-phase-to-red-gate path is
asserted once by the gate's own phase tests, not 33 times.

Gate command: the project gate, `bench gate`.

### Seam diagram

**Seam A — the canary sweep behind its injected `Runner`.**

    trigger: `bench canary [root]`, and the gate's canary phase
        │
        ▼
    fixture tree ──▶ [ SweepTier                              ] ──▶ error or nil
    tier          ──▶ [   walk → group → compile → baseline    ]
    Runner        ──▶ [   → per-fixture invoke → bite/vacuity  ] ──▶ RunCall stream
                          ◀ tests attach here: inject a recording Runner over a
                            temp fixture tree; assert the recorded call kinds,
                            counts, packages, binaries, subject roots, and env

**Seam B — the gate's phase table.**

    trigger: gate run (outer), and any surviving inner-gate spawn
        │
        ▼
    root, kit ──▶ [ BenchkitPhases → phasesForMode ] ──▶ []Phase
    environment ▶ [                               ]
                      ◀ tests attach here: build the table in each mode, with and
                        without an ambient export of the retired variable, and
                        assert the contract phase's argv

**Seam C — the real sweep over the kit's own fixtures.**

    trigger: the gate's canary phase
        │
        ▼
    kit tree ──▶ [ real Runner: compile + invoke per fixture ] ──▶ exit 0 / red
                      ◀ tests attach here: no new test — the gate is the assertion;
                        every one of the 33 fixtures must still bite non-vacuously

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | no behavior-owned fixture produces an inner-gate call | seam A, recording runner | new test asserting the recorded call kinds for a behavior-owned tree contain no gate-spawn kind; reds against today's sweep, which spawns one per fixture | the cheapest wrong implementation adds compile-and-invoke beside the existing spawn and reports the same green while paying more, not less |
| 1 | the real runner execs the compiled binary for a bite call and the gate script only for a gate call | seam A, default runner unit test | new test building each call kind and asserting the command the default runner constructs; reds today, where every call builds `bash <gate>` | this is the row that stops the whole map being satisfied by relabelling: every other sweep assertion runs through an injected fake, so a change that emits correct call metadata and still dispatches a nested gate passes all of them and buys nothing |
| 1 | a behavior-owned fixture produces exactly one bite invocation | seam A, recording runner | new test counting invocations per fixture directory; reds at zero and at more than one | a fixture that produces no invocation grades nothing and still reports green, which is the silent-unswept failure the canary exists to prevent |
| 1 | a bite invocation carries no canary inner-gate marker and no phase pin | seam A, recording runner | new test asserting neither variable appears in any bite call's environment; reds if the bite path reuses the inner-gate environment builder wholesale | the Handoff makes this an assertable, and reusing the inner-gate builder is the obvious shortcut — it leaves each bite claiming to be a nested gate to anything that later reads those markers |
| 1 | every other family still spawns its inner gate unchanged | seam A, recording runner | new test asserting a phase-named family's fixture still produces a gate-spawn call; reds if the migration converts every family to compile-and-invoke | the existing phase-family test asserts only the phase pin's value, so a wholesale conversion that carried the pin on a compile call could pass it; the spawn kind has to be asserted directly |
| 2 | one compile call per package group, regardless of how many fixtures the group holds | seam A, recording runner | new test with two fixtures under one package and one under another; asserts exactly two compile calls, naming both packages; reds at one-per-fixture | compiling per fixture is the implementation that passes every bite assertion while buying nothing, which is precisely the failure this slice exists to remove |
| 2 | every fixture in a group is invoked with that group's binary | seam A, recording runner | new test asserting each invocation's binary path equals its group's compile output; reds if the paths are per-fixture or crossed | a per-group compile whose output nobody reuses is the same waste wearing the right call count |
| 3 | each invocation carries its own fixture's materialized tree as the subject root | seam A, recording runner | new test asserting the subject-root value per invocation is that fixture's own work directory and distinct across fixtures; reds on a shared or empty value | a shared root would grade every fixture against one tree, so 32 of 33 would report did-not-bite — or worse, all 33 would bite for one fixture's reason |
| 3 | an ambient export of the subject-root variable does not reach a bite invocation | seam A, recording runner | new test exporting the variable in the test process and asserting exactly one occurrence per invocation env, carrying the sweep's value | Go's exec environment has no duplicate-key precedence, so a set without its matching strip hands an ambient export control of what every fixture grades |
| 4 | one baseline invocation per package group, in the same shape as the group's fixtures | seam A, recording runner | new test asserting baseline invocation count equals group count and each baseline carries its group's binary; reds on a single shared baseline or a gate-spawn baseline | a baseline running a different shape emits output no fixture in the group can produce, so the collision screen compares against the wrong noise in both directions — matching what the group never prints, missing what it does |
| 4 | contract groups do not fold into the unscoped baseline legacy flat fixtures share | seam A, recording runner | new test with a behavior-owned group and a legacy flat fixture in one tree, asserting the unscoped baseline count stays at one and the contract group gets its own | the existing unscoped-baseline test builds no behavior-owned group at all, so it cannot see this folding; a map claiming it as covered would grade nothing |
| 4 | a contract group whose baseline produces no output is a red naming the group | seam A, fake runner returning empty baseline output | new test driving the empty output through a behavior-owned group under the compiled shape; reds if the refusal is bypassed on the new path | the existing empty-baseline test exercises conformance groups through the gate-spawn path, so the compiled path could route around the refusal and still pass it — and an empty baseline contains no EXPECT, so every fixture in that group would clear the collision screen unguarded |
| 5 | a compile that exits nonzero reds naming the package, and no fixture of that group reports green | seam A, fake runner returning a compile failure | new test asserting the returned error names the package and that the group's fixtures do not report success; reds if the failure is swallowed or reported as did-not-bite | a swallowed compile failure turns a broken package into a silently unswept family, which is the same class as an unbound family stage 1 refused |
| 5 | a compile that exits zero but writes no binary reds naming the package | seam A, fake runner returning success with no binary written | new test asserting the same red; reds if the sweep proceeds to invoke a path that does not exist | `go test -c` on a package with no test files exits zero and writes nothing, so an exit-code-only check would invoke a missing file and surface an exec error naming nothing the author can act on |
| 6 | every compiled binary lands under one sweep-owned temporary directory | seam A, recording runner | new test asserting all compile output paths share one parent and that the parent is not inside the kit tree | binaries written beside the source would dirty the tree the gate grades, turning a sweep into a git-status change |
| 6 | the directory is removed when the sweep returns an error, not only on success | seam A, fake runner forcing a mid-sweep error | new test asserting the recorded parent directory is gone after the failing sweep; reds against a success-path-only cleanup | the obvious implementation defers cleanup after the happy path or forgets the early error returns, and the leak is invisible until a disk fills |
| 6 | two package paths cannot produce the same binary file name | seam A, recording runner | new test with fixtures under two packages whose paths differ only in a separator position; asserts distinct compile output paths | a name derived by flattening separators lets one group's binary overwrite another's mid-sweep, and the loser grades the wrong package's tests |
| 7 | the contract phase's argv is the unnarrowed subtree in every mode | seam B, gate phase table | new test replacing the deleted narrowing tests: build the table in outer and inner mode and assert the subtree argv; reds if any narrowing survives | a narrowing left behind with no writer is dead code that a later edit can reanimate against a variable nothing sets |
| 7 | an ambient export of the retired variable changes nothing | seam B, gate phase table | new test exporting the retired name literally and asserting the same unnarrowed argv | the retired name is the one input that could still be sitting in an operator's shell or a stale adapter, and inertness is the property the deletion claims |
| 7 | no reference to the retired narrowing or its variable survives in Go source | seam B, repository sweep | `rg -n -e narrowContractScope -e ContractPackageEnv -e BENCH_CANARY_CONTRACT_PACKAGE --type go` returns nothing; reds today across the gate and canary packages, including the comment reference in `internal/gate/manifest.go` | a partial deletion that leaves the constant exported invites a second caller, which is how a retired mechanism comes back |
| 8 | the subject-root variable is declared once and consumed by every reader | seam B, repository sweep | `rg -c BENCH_CONTRACT_ROOT --type go --glob '!*_test.go'` names exactly one file; reds today, where the gate and the contract helper each carry the literal | two literals for one name is the duplicated knowledge the code standard forbids, and a rename that misses one leaves the sweep setting a variable nobody reads |
| 9 | every one of the 33 behavior-owned fixtures still bites, non-vacuously, under the new run shape | seam C, the real gate | `bench canary .` reds naming any fixture whose EXPECT is not observable in test-level output | this is the migration-casualty class the map named; the prototype found the first fixture's EXPECT present, but the remaining 32 are unverified until the sweep runs |
| 9 | a bite invocation that exits nonzero with the EXPECT absent from its output reports did-not-bite | seam A, fake runner returning nonzero and non-matching output | new test on that branch; the existing did-not-bite test covers only the exit-zero-with-matching-output branch, so this direction is unasserted today | the migration's whole casualty class is an EXPECT that stops appearing while the test still fails — exactly nonzero-with-absent-output, which is the branch nothing currently grades |
| 9 (edge of 6) | a fixture root that fails to materialize is a red, never a silent skip | seam A, real materialization | already covered — the broken-symlink case in `canary_concurrency_test.go` asserts the `setup failed:` diagnostic, and materialization is untouched by this slice | the Handoff names this hostile input; the row claims only what the existing test proves, which is that the failure surfaces rather than skipping the fixture |
| 10 | the timeline asset records the post-change canary and gate wall-clock beside the pre-change figures | manual, reviewer-graded | not TDD-able — a measurement is a fact about one host at one moment, and no assertion can distinguish a real figure from an invented one | recorded as a stated limitation: the map made this ship evidence rather than a threshold, so the reviewer reading the asset is the check |
| 11 | no document names the retired narrowing or its variable | seam B, repository sweep | `rg -n BENCH_CANARY_CONTRACT_PACKAGE --glob '!*.go' --glob '!specs/**' --glob '!decisions/**'` returns nothing; reds today with the `session-handoff.md` reference | a doc naming a retired mechanism sends the next cold session to a surface that no longer exists; the spec and the map are excluded because describing the retirement requires naming what was retired |
| 12 | a behavior-owned fixture declaring its own phase table is swept normally | seam A, recording runner | new test building such a fixture and asserting it produces a bite invocation; reds today with the structural refusal | the removal has to be observable, otherwise the refusal could be left in place and the story silently unbuilt |
| 12 | the two surviving structural refusals still red | seam A, recording runner | already covered — the existing no-package-directory and package-absent-from-the-tree cases in `TestSweepRejectsStructurallyUnscopableContractFixtures` | removing one case of a table-driven refusal by deleting the whole test is the cheap wrong edit; the other two cases surviving is what proves the scope |

### Edge inventory

- **absent vs present-but-empty** — a fixture bound to no package and a package absent
  from the tree are covered by the surviving structural refusals (story 12 row); a
  package present but holding no test files is covered by the no-binary-written row
  (story 5).
- **ambient environment export** — covered twice: the subject-root variable (story 3)
  and the retired contract-package variable (story 7).
- **required tool missing from PATH** — `go` absent: the compile call exits nonzero and
  lands on story 5's compile-failure red, which names the package. No separate row; the
  diagnostic is the same one and adding a probe would duplicate the failure path.
- **paths containing spaces** — the sweep's temporary directories and the compiled
  binary paths are passed through `exec.Command` argv, never a shell, and the existing
  sweep already materializes fixtures under generated temp paths. **Won't handle** as a
  new row — no shell interpolation exists on this path to break.
- **re-run idempotency** — a second sweep in the same process creates a fresh temporary
  directory and recompiles; story 6's cleanup rows cover the residue. **Won't handle**
  as a separate row — a second sweep is the first sweep's assertions run twice.
- **interrupt mid-sweep** — SIGINT leaves the temporary binary directory behind exactly
  as it already leaves fixture work directories behind. **Won't handle** — unchanged
  from today's behavior, and fixing it is a signal-handling capability the sweep has
  never had.
- **special files and dangling symlinks in a fixture tree** — materialization is
  untouched by this slice; a root that fails to materialize is covered by the story 9
  edge row, citing the existing broken-symlink assertion.
- **concurrent workers** — two workers writing one binary path is covered by story 6's
  collision row; two workers sharing an environment backing array is covered by the
  existing per-call environment copy, whose test `canary_concurrency_test.go` already
  owns.
- **amputated caller** — the gate's contract phase after the narrowing is deleted is
  covered by story 7's argv rows in both modes.

## Out of scope

- **The prepared-artifact hoist** (map #2/#3). The five `surface/artifact` fixtures each
  pay a full artifact contract suite, and that multiplier survives this slice by design —
  it is the next slice's target and is blocked on an open test-independence ruling.
  Estimate: 8 edits, 5 gate runs, plus the ruling.
- **Reviving the outer conformance and contract width cap.** The CPU oversubscription
  measured in #1 may shrink to nothing once nesting is gone, so the re-check is
  deliberately sequenced after this slice's measurement. Estimate: 3 edits, 3 gate runs.
- **Oracle-semantics levers** — gate-verdict caching keyed on the pinned subject, and
  `-count=1` freshness. These graduate into tickets only if the post-stage-2
  re-measurement leaves the gate above the 60 s stop rule. Estimate: unknown until the
  measurement lands, which is why they are parked rather than sized.
