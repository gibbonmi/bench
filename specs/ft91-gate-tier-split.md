# FT91 — two-tier gate: dev tier and `bench prep-release`

Status: implemented

## Problem

Running `bench gate` on this repo costs 400–830 s of wall clock, and a human waits for
all of it. The 2026-07-26 measurement found the cost is not spread across the
conformance checks — `checkPackageCoreAndGuards` is ~99.8% of the phase, and inside it a
single release-evidence probe (a four-platform artifact matrix build, a reproducibility
rebuild, and a real `release-preflight.sh --mode verify`) accounts for ~372 s. Every
small change pays that price to re-prove something only a release can consume.

The same phase carries two latent hangs. The inner `go test` leans on Go's test cache; on
a cache miss `internal/preflight` alone runs 676 s and blows the 600 s default package
timeout, which presents as a gate hang rather than a slow test. And that inner `go test`
includes `internal/conformance` itself — the package whose entry point *is* the outer
run. A live recursion cascade was observed here on 2026-07-26 (children outliving their
600 s package timeout as orphans, killed by hand); the tree has since neutralized the
recursive arm — the inner subprocess env strips `BENCH_CONFORMANCE_ROOT` and
`TestRootConformance` skips without it, under its own regression test — but that
protection is an implicit env contract, not a declared one, and nothing would catch its
regression early: phases carry no individual timeout, so the only backstop is the
45-minute whole-gate deadline, which SIGKILLs the process group and records a `timeout`
verdict long after the machine has become unusable. Story 3 exists to make the
non-recursion contract structural and gate-tested instead of one env-strip regression
away from that failure mode.

There is also no observability: when the gate is slow, nothing in its output says which
check spent the time. The number above came from a throwaway probe file that had to be
written, run, and deleted.

## Solution

Split the oracle into two tiers with a boundary no release can bypass.

The **dev tier** — `bench gate`, the shift loop, `/bench-final-check`, and the pre-push
hook — answers one question: *does the kit work from the tree?* It runs every static
conformance check, and drops the release-evidence probe. Its inner `go test` excludes the
release-only packages the way it already excludes `internal/contract`, and excludes
`internal/conformance` from the unfiltered run — replacing it with a second, filtered
run of that package that skips only its entry-point test, so the package's own suite
stays enforced. Dev green is explicitly the narrower claim.

The **ship tier** — a new `bench prep-release` command — runs the artifact matrix build,
the cross-compile matrix, the release preflight verify, the release-only package tests,
and the canary fixtures that grade ship-tier checks. Exit 0 means ship green with evidence
written to `dist/preflight/release-index.json` and `dist/artifacts`. No check loses
authority, but the hard boundary is not this evidence: the publish path's pre-existing
refusal (`VerifyPublishAuthority`) demands a *publish*-mode index that only the release
workflow's own preflight produces, while `prep-release` emits verify-mode evidence.
`prep-release` is the maintainer's rehearsal that proves the moved checks green before
that boundary is approached; the checks move to a surface that runs once per release
instead of once per commit.

Restaging is not check-weakening. Every check keeps full authority; `bench prep-release`
additionally requires a current dev-green verdict, so a dev-tier failure reds the ship
tier too.

Per-check timing becomes permanent, driver-owned gate output: one stable-format line per
check on every run, in byte-stable order, with values free to vary. The next
wall-clock question is answered by reading the gate, not by writing a probe.

## User stories

1. As a developer running `bench gate` on a small change, I want the release-evidence
   probe to run only at the ship tier, so that I stop paying ~372 s per commit to rebuild
   a four-platform artifact matrix I am not shipping.
   `Line: opus / medium.` This changes what the oracle grades, and the profile routes
   gate and conformance logic to the mid tier because a wrong gate is the worst class of
   bug in a kit whose premise is that the gate is the oracle.

2. As a developer, I want the dev tier's inner `go test` to exclude the release-only
   packages (`internal/preflight`, `internal/releaseevidence`, `internal/publication`),
   so that a cold Go test cache makes the gate slower rather than making it hang past the
   600 s package timeout.
   `Line: opus / medium.` The exclusion predicate decides what the oracle tests, so it
   carries the same correctness weight as the tier membership it sits beside.

3. As a developer, I want the dev tier's inner `go test` to exclude
   `internal/conformance` from the unfiltered run and grade it through an explicit,
   registry-driven filtered invocation instead, so that the recursion the map recorded
   stays structurally impossible — a declared, tested `-skip` contract — rather than one
   implicit env-strip regression away, while the package's ~30 other test files stay
   enforced by the gate. This change is deliberately zero-cost: the filtered set equals
   what the child run already executes today, because the entry point already skips
   under the stripped env.
   `Line: opus / medium.` The inner unfiltered run is the only place the conformance
   package's own suite executes during a gate, so a bare exclusion would silently drop
   it from the oracle — the filtered replacement is what keeps the no-weakening ruling
   intact, and drawing that skip list correctly needs judgment.

4. As a developer diagnosing a slow gate, I want the conformance driver to print one
   timing line per check on every run in a stable order, so that I can see where the wall
   clock went without writing and deleting a probe file.
   `Line: opus / medium.` The output shape is load-bearing because canary EXPECT matching
   is substring-based against inner-gate output, so a careless format can turn an
   unrelated fixture vacuous.

5. As a maintainer preparing a release, I want `bench prep-release` to run the artifact
   matrix build, the cross-compile matrix, the release preflight verify, and the
   release-only package tests, so that everything the dev tier stopped running is proved
   once at the boundary that matters.
   `Line: opus / medium.` It is a thin route over existing scripts, but it is also the
   surface that now carries the authority the dev gate gave up, so its exit code has to
   be right.

6. As a maintainer, I want `bench prep-release` to refuse unless there is a current
   dev-green verdict for this exact tree, so that ship green can never be claimed over a
   tree the dev tier has not passed.
   `Line: opus / medium.` The precondition reuses the gate package's existing
   tree-and-oracle-bound verdict authority, and misreading that contract would silently
   weaken the whole split.

7. As a maintainer, I want the canary fixtures that grade ship-tier checks to run their
   inner gate under `bench prep-release` rather than under the dev gate, so that those
   checks are still proved to bite without dragging the probe back into every commit.
   `Line: opus / medium.` The canary is the gate guarding the gate, and a fixture that
   silently stops biting is exactly the rot the canary exists to catch.

8. As a maintainer, I want `bench prep-release` to be routed through the shared argument
   grammar and reachable from every shipped surface, so that it behaves like every other
   `bench` subcommand instead of hand-rolling its own flag handling.
   `Line: sonnet / low.` This is `bench` CLI shell plumbing at a known seam, which the
   profile routes to the cheap tier, and the conformance routing registry fully observes
   whether it was done right.

9. As a developer finishing a phase, I want `/bench-final-check` on green to print one
   line reminding me that ship-tier verification has not run, so that the narrower
   meaning of dev green is stated at the moment it could be misread.
   `Line: fable / high.` The profile's cached routing sends command and guidance prose to
   the top tier because it compounds through every session that loads it, and this is a
   standing grant rather than a bump.

10. As an agent picking up this repo cold, I want `bench prep-release` documented in the
    CLI inventory and the two tiers described in the project profile's gate section, so
    that a fresh session learns what dev green does and does not claim.
    `Line: fable / high.` Same cached doc-authoring routing as story 9, and the cold-pickup
    inventory is the one surface a session reads before it knows anything else.

## Implementation decisions

**The tier split lands inside the conformance driver, not around it.** `RunConformance`
stays the deep unit: it gains an explicit tier registry — an indexed slice of
`{name, tier, fn}` — and the checks stay pure functions over a read-only tree. Callers ask
for a tier; nothing outside the package learns which check belongs where. This is the
seam the Handoff named, and it keeps tier membership and timing behind one boundary.

**The registry's metadata is promoted into a leaf package
`internal/conformance/registry`; the check functions stay test-side.** Today
`internal/conformance` contains zero non-test `.go` files, so nothing outside its test
binary can read the registry — but the canary tier test and the inner skip list both
need it. The metadata cannot live as a non-test file inside `internal/conformance`
itself: that package's in-package test files already import `internal/canary`, so
`internal/canary` importing `internal/conformance` is an import cycle the test binary
refuses to build. The metadata therefore lands in a leaf package that imports nothing
from `internal/conformance`: an indexed table of (check name, tier), plus the inner
skip list described below. `internal/canary` and the conformance test files both import
the leaf. The check functions stay in test files, bound to metadata entries, with a
completeness assertion that every metadata row has exactly one bound function and vice
versa, so the two halves cannot drift. This is a named scope addition the Handoff did
not anticipate.

**The conformance package's own suite stays in the oracle via a second, filtered inner
run.** The inner unfiltered `go test` inside `checkGoCore` is the only place that
package's ~30 test files execute during a gate (the outer phase runs
`^TestRootConformance$` alone), so a bare exclusion would de-enforce the registry tests,
the fixture-bite tests, and four of this spec's own red signals — colliding with the
map's closed no-weakening ruling. Instead: the unfiltered package list excludes
`internal/conformance`, and `checkGoCore` adds one more invocation —
`go test ./internal/conformance -skip <pattern>` — guarded on that package existing in
the graded root. The guard matters: `checkGoCore` grades whatever tree it is given, and
linked repos plus the two canary fixtures that carry a `go.mod` have no
`internal/conformance`, so an unguarded invocation would add a spurious `go test
failed` diagnostic to every such gate. The pattern comes from a **skip list of test
names** held in the leaf registry package. A marker on *checks* cannot produce this
pattern — the skip selects *tests*, and there is no check→test map; the entry-point
test `TestRootConformance` runs every check and so belongs to no single check's marker.
The list is honest data, single-sourced in the leaf package, and today it contains
exactly `TestRootConformance`. It exists as defense-in-depth only, not cost control:
the tree already neutralizes the recursive arm (the inner subprocess env strips
`BENCH_CONFORMANCE_ROOT`, under its own regression test), so the entry point already
skips in the child run and the filtered set equals what runs today — the change costs
nothing and removes nothing from the oracle. Its value is that the non-recursion
contract becomes declared and gate-tested instead of implicit. The map's item-9 claim
that the recursion is live is contradicted by the tree — the tree wins.
**(flagged for reviewer veto — deviates from Handoff item 9's live-recursion premise;
the alternatives are keeping the status quo inclusion, which reverses the map's decided
exclusion, or also skipping the fixture-bite sweep for a real cost cut, which removes
tests from the oracle and reopens the no-weakening ruling.)**
The remaining suite is genuinely affordable: most fixture trees short-circuit
`checkGoCore` at its `go.mod` guard, and the two that do carry a `go.mod` are two-file
modules whose full body runs in seconds. Mechanically, `checkGoCore` builds the
invocation through a deterministic helper beside `goCoreTestPackages` — it takes the
graded root, and its only I/O is the package-presence stat the guard needs — whose
output embeds the leaf registry's pattern, so a test asserts the helper's argv against
the registry directly; argv observation via a recorder is impossible here because `go`
itself cannot be swapped out. Two assertions keep the list from rotting: an anti-drift check
that every listed name matches an existing test function (via `go test -list`, which is
the full inventory — `-list` ignores `-skip`, so it can enumerate but never observe
filtering), and the non-vacuity row below, which compiles the skip pattern in-test and
asserts the named cheap tests do **not** match it. Filtering uses `go test -skip`
(Go ≥1.20; this module is at 1.25), not an in-test guard, because the gate bans bare
`t.Skip` outside the capability helper package. Nothing else leaves the oracle.

**What actually moves is `runReleaseEvidenceProbe`, not all of `checkReleasePreflight`.**
The verification read changed the granularity here, and the change is material enough to
state plainly. `checkReleasePreflight` has two halves: a static half (registry JSON
validity, workflow structure anchors, the exact-toolchain pin, the comparator and
offline-smoke script anchors) that is file reads and regex and costs milliseconds, and
`runReleaseEvidenceProbe`, which materializes an authenticated clone, runs
`scripts/build-artifacts.sh`, inspects every archive, and executes a real
`release-preflight.sh --mode verify`. The ~372 s is entirely the probe. Moving the whole
check would also strip nineteen `package-core-guard` canary fixtures of the diagnostic
they assert, turning the dev gate red on "did not bite" — so the dev tier keeps the static
half and the ship tier takes the probe. **(flagged for reviewer veto — Handoff item 7.)**

**The cross-compile matrix needs no removal, only a caller.** `crossCompileMatrix` is
already a no-op behind a `!stress` build tag, so the dev tier does not run it today. The
work is the other direction: `bench prep-release` must invoke the conformance suite with
`-tags stress` so the four-platform matrix actually runs somewhere.

**Release-only packages are `internal/preflight`, `internal/releaseevidence`, and
`internal/publication`.** The import graph settles this: `internal/releaseevidence` is
imported only by `internal/preflight`; `internal/preflight` and `internal/publication` are
imported only by `cmd/bench/main.go`'s dispatch switch. No ordinary subcommand package
reaches any of them. Their suites are the heavy ones — `internal/preflight` rebuilds the
binary through `scripts/go-build.sh` in seven separate tests, and the other two shell out
to `node`. The exclusion predicate generalizes the existing contract-subtree predicate
rather than growing a second one. **(flagged for reviewer veto — Handoff item 7.)**

**The dev-green precondition reads the gate verdict cache, not `bench gate pin`.** The
Handoff recommended the pin machinery; the verification read shows `bench gate pin` is a
different mechanism — it records HEAD's committed `.bench` tree hash for the managed
pre-push hook to verify, and carries no verdict at all. The facility that answers "is
there a current green verdict for this exact tree" is the gate verdict cache:
`gate.Inspect(root)` returns an `Inspection` whose `ReusableGreen` field is true only when
the record is green, its tree hash matches the current subject, its oracle hash matches,
the subject is closed, and the verdict is under the ten-minute freshness window.
`bench prep-release` calls that and refuses on anything else, naming `bench gate` as the
remedy. The ten-minute window is the real ergonomic cost and is stated so the reviewer can
veto it: the dev gate must have finished within ten minutes of starting `prep-release`.
This is an **observer**, not a second authority grant. ADR 0002 posture 5 confines verdict
*reuse* — skipping a run because a cached green exists — to `internal/commit`, which is the
only caller that reads `ReusableGreen` for authority today. `prep-release` never skips
anything on the strength of the cache: a green observation lets it *start*, and every
ship-tier check still runs in full. Refusing on a non-green observation is strictly more
conservative than the status quo, so the posture is unchanged and no ADR edit is needed.
Prior art for the call shape is `internal/commit/commit.go:118`;
`internal/status/status.go:209` is the existing observe-only precedent.
**(flagged for reviewer veto — Handoff item 7.)**

**Timing travels through a file under the git dir, because no in-test write survives a
passing `go test`.** The conformance phase runs under non-verbose `go test`, which
suppresses the test binary's stdout and `t.Log` alike when the package passes — the one
case where timing is most wanted, a slow green gate, is exactly the case a stdout
mechanism goes silent (verified empirically during falsification). So the driver writes
the timing lines to a file beside the gate verdict cache (`<git-dir>/bench-conformance-timing`,
truncated per run), and the gate runner prints that file's lines into gate output after
the conformance phase completes — green or red, and in **both outer and inner mode**,
because the canary vacuity guard grades inner-mode output and an outer-only print would
leave that guard permanently blind to the timing format. The file's git dir is derived
from the runner's **graded root** by resolving `<root>/.git` directly — never from cwd,
and never through `git rev-parse --absolute-git-dir`, which ascends: inner gates run
with cwd at the host kit checkout while grading per-fixture temp roots, canary fixtures
run concurrently, and an ascending or cwd-derived path would have every concurrent run
truncating one shared host file (or, for a temp root that happens to sit inside a git
repo, silently rebinding to that repo). When the root has no `.git` or no timing file,
the print **fails open and silent** — that stated fail-open, plus non-ascending
derivation, is what keeps `TestRunnerInnerModeByteShape` green against its `t.TempDir`
root. Driver-side tests that need timing observe it in git-inited fixture roots (the
existing `materializeConformanceFixture` prior art), since a bare temp dir gives the
file nowhere to live.
The git-dir location itself is sound: the gate's tree hash and oracle hash both exclude
`.git`, and the verdict cache already lives there, so a mid-run write cannot trip the
subject-drift rejection. Format is one line per check carrying a zero-padded index, the
registry check name verbatim, and a duration; the index is what makes ordering
byte-stable while durations vary, and the name being registry-sourced keeps the
enumeration single-sourced with tier membership.

**`prep-release` is a thin route in its own small package.** It invents no machinery: it
calls `gate.Inspect`, then the existing `scripts/build-artifacts.sh`, the conformance
suite at the ship tier with `-tags stress`, `go test` over the three release-only
packages, `scripts/release-preflight.sh --mode verify`, and the ship-tier canary set. It
takes a flat argv with no subcommand tree, so it is recorded in the conformance routing
registry as `routed` and reaches the shared `usage.Parse` grammar — not as a `whyNested`
exemption like `release` and `release-preflight`.

**The release path's refusal needs no new code, and `prep-release` evidence does not
satisfy it.** `publication.VerifyPublishAuthority` already reads
`dist/preflight/release-index.json` and requires mode `publish`, non-focused scope, green
status, and a matching profile, and `VerifyApprovedSet` cross-checks per-artifact digests
against `dist/artifacts`. `prep-release` runs `release-preflight.sh --mode verify`, whose
index is verify-mode — deliberately insufficient for publish authority, so a rehearsal
can never be mistaken for the release's own preflight. This spec asserts the existing
refusal rather than building it.

**Canary inner gates run the dev tier; the two probe-derived fixtures move to the ship
tier.** The Handoff's item-7 recommendation was "dev" without qualification; moving
`release-package-evidence-omitted` and `release-digest-binding-omitted` to
`prep-release`'s ship-canary step is the necessary consequence of the probe moving —
under a dev inner gate their checks no longer run, so they would report "did not bite"
on every commit. **(flagged for reviewer veto — Handoff item 7.)**

## Testing decisions

A good test here drives the real driver or the real command and observes what a reviewer
would observe: which checks ran, what got written to disk, what the exit code was. Tier
membership is a fact about the registry and is asserted directly; everything else is
asserted black-box.

Prior art the new tests follow: `TestGoCoreTestPackagesExcludesContractSubtreeOnly` in
`internal/conformance/package_core_checks_test.go` is the existing shape for asserting a
package-list exclusion, and `TestCheckGoCoreDoesNotWriteRootDistBench` in the same file is
the existing shape for asserting that a check writes nothing into the graded root. The
canary families in `tests/canary/` are the prior art for proving a check still bites.

**The gate command that defines done: `.bench/gate.sh`.**

Two gate-shape interactions have to be handled by this build rather than discovered by it.
First, `runFixture` in `internal/canary/canary.go` rejects a fixture as vacuous when its
EXPECT string is a substring of the *baseline* (empty-fixture) inner-gate output. Timing
lines reach gate output through the runner's file-print step, so they widen that
baseline on every run, and a timing format that embeds free prose could
newly make an unrelated fixture vacuous and turn the gate red. The format is therefore
constrained to an index, the registry check name, and a duration. One nearby test
grades runner output byte-for-byte and is *not* insulated from this change:
`TestRunnerInnerModeByteShape` (`internal/gate/runner_test.go:152`) runs a fake phase
named literally `conformance` and asserts exact inner-mode stdout — a runner that
prints the timing file after the phase named `conformance` is graded by it. It stays
green only because the file path derives non-ascendingly from the runner's graded root
(a `t.TempDir` with no `.git`, where the print fails open and silent), which is why
root-derivation and the fail-open are stated decisions above, not implementation
details. Beyond that test, there is no byte-exact or line-count assertion
against real conformance output — the vacuity guard is the only coupling, and it is a
runtime guard, so `bench canary` has to actually run after the format lands rather than
being reasoned about from source. Second, EXPECT matching is plain `strings.Contains` against inner-gate
output, so the nineteen `package-core-guard` fixtures asserting static release-preflight
diagnostics keep biting only because the static half stays dev tier; the two probe-derived
fixtures move with the probe.

### Seam diagram

**Seam 1 — the conformance driver (tier membership + timing).**

    trigger: gate conformance phase (dev) │ bench prep-release (ship)
        │
        ▼
    root, kitRoot  ──▶  [ RunConformance(tier)            ]  ──▶  []diagnostic
    tier            ──▶  [   indexed registry of checks   ]  ──▶  timing lines (git-dir file,
                                                                  printed by the gate runner)
                        [   pure check fns unchanged      ]
                              ◀ tests attach here: call RunConformance with each tier
                                against a fixture tree; assert which diagnostics appear,
                                that dist/artifacts is untouched at the dev tier, and
                                that timing line order is identical across two runs.

**Seam 2 — the `prep-release` route.**

    trigger: maintainer runs `bench prep-release`
        │
        ▼
    repo root  ──▶  [ preprelease.Command      ]  ──▶  exit 0 | nonzero + diagnostics
                    [  gate.Inspect precheck   ]  ──▶  dist/artifacts/*
                    [  build-artifacts.sh      ]  ──▶  dist/preflight/release-index.json
                    [  conformance ship tier   ]
                    [  release-only go test    ]
                    [  release-preflight.sh    ]
                    [  ship-tier canary set    ]
                          ◀ tests attach here: run the built binary in a throwaway
                            fixture repo with the verdict cache seeded green, absent,
                            stale, and tree-mismatched; observe exit code and what
                            appeared under dist/.

**Seam 3 — the canary fixture tier.**

    trigger: dev gate canary phase │ prep-release ship-canary step
        │
        ▼
    fixture dir  ──▶  [ canary sweep      ]  ──▶  "" | "did not bite" | "vacuous"
    tier         ──▶  [  fixture tier tag ]
                          ◀ tests attach here: assert each fixture's declared tier
                            matches the tier of the check its EXPECT names, and that
                            the dev sweep selects only dev-tier fixtures.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | The release-evidence probe is absent from the dev tier and present in the ship tier | conformance driver | `go test ./internal/conformance -run TestTierMembership` — fails to compile today, no tier registry exists | Asserting both directions rejects the degenerate registry that marks every check dev-tier and the one that marks every check ship-tier |
| 1 | A dev-tier run's timing lines name exactly the dev-tier registry entries — the probe's check never appears | conformance driver | `go test ./internal/conformance -run TestDevTierExecutesExactlyDevChecks` — red, no registry and no timing exist today | Timing lines are written per executed check, so comparing their names against the dev-tier metadata catches a driver that ignores the registry and runs ship checks anyway — the observation a probe-fault seed cannot provide, since the probe only activates on a real tree where the test would cost the full 400 s+ path |
| 1, 7 | A seeded ship-tier-only failure reds `prep-release` while `bench gate` stays green | conformance driver + prep-release | `go test ./internal/contract/... -run TestShipTierFailureIsolation` — red, no ship tier to seed against; the seed is a failing test in a release-only package, the cheapest surface only the ship tier runs | A single-tier implementation cannot produce two different verdicts for one seeded fault, so it fails this row by construction; seeding a release-only package test rather than a probe fault keeps the fixture cheap, since the probe activates only on a real tree where the run costs the full 400 s+ path |
| 2 | `internal/preflight`, `internal/releaseevidence`, and `internal/publication` are absent from the dev inner `go test` package list and present in the ship-tier list | conformance driver | `go test ./internal/conformance -run TestGoCoreTestPackagesExcludesReleaseOnly` — red, all three are currently included | Enumerating all three by name, in both directions, blocks an implementation that excludes only `internal/preflight` and leaves the other two to blow the timeout |
| 3 | `internal/conformance` is absent from the dev inner unfiltered `go test` package list | conformance driver | `go test ./internal/conformance -run TestGoCoreTestPackagesExcludesContractSubtreeOnly` — red once the assertion at `package_core_checks_test.go:300` is inverted (and the test renamed to match its new meaning) | The existing test asserts the package *is* included, so the inverted expectation fails against current behavior and can only go green when the recursion is actually cut |
| 3 | `checkGoCore` runs the second filtered conformance invocation with the skip pattern built from the leaf registry's skip list | conformance driver | `go test ./internal/conformance -run TestInnerSuiteRunsConformanceFiltered` — red, no second invocation exists | An implementation that only deletes the package from the list passes the exclusion row above but fails this one, which is exactly the silent de-enforcement the no-weakening ruling forbids |
| 3 | The filtered run is not vacuous: every skip-list name matches an existing test function (`go test -list` as the inventory), and the compiled skip pattern does not match the named cheap tests (the registry completeness assertion, the fixture-bite suite) | conformance driver | `go test ./internal/conformance -run TestFilteredRunSelectsRealTests` — red, no skip list exists | `go test -list` ignores `-skip`, so listing can never observe filtering — compiling the pattern in-test and asserting non-matches is the only mechanism that makes an over-broad or rotted pattern red |
| 3 | The filtered invocation's argv embeds the leaf registry's pattern and is absent for graded roots without `internal/conformance` | conformance driver | `go test ./internal/conformance -run TestInnerConformanceArgs` — red, the argv helper does not exist | Asserting the pure helper's output against `registry` in both directions rejects a hardcoded skip literal that drifts from the single source, and rejects the unguarded invocation that would add a spurious diagnostic to every linked Go repo's gate |
| 4 | Every check in the registry writes exactly one timing line to the timing file per run | conformance driver | `go test ./internal/conformance -run TestTimingLinePerCheck` — red, nothing is written today | Deriving the expected line count from the same registry the driver iterates means a check added without timing fails immediately, rather than passing against a hardcoded count |
| 4 | Timing line order in the file is byte-identical across two runs of the same tree while durations differ | conformance driver | `go test ./internal/conformance -run TestTimingOrderStable` — red, no file to compare | Comparing only the ordered name sequence across runs catches a map-iteration implementation, which is the natural wrong way to build a check registry |
| 4 | The gate runner prints the timing file's lines into gate output after the conformance phase — on green and red alike, in outer and inner mode alike, from the graded root's git dir resolved non-ascendingly | gate runner | `go test ./internal/gate -run TestRunnerPrintsConformanceTiming` — red, no such plumbing exists; the test seeds a timing file in a fake root's git dir and asserts the lines appear in both modes, and drives a second root that has no `.git` but sits under a git-bearing ancestor holding a decoy timing file, asserting nothing prints | Non-verbose `go test` swallows everything a passing test binary prints, so the runner-side print is the only link that puts timing into gate output; the both-modes case blocks the outer-only implementation that would leave the canary vacuity guard blind, the seeded-root case blocks a cwd-derived path, and the decoy-ancestor case is the only signal that makes an ascending `git rev-parse` implementation — which would let concurrent canary fixtures corrupt one shared host file — go red |
| 4 | No timing line makes an existing canary EXPECT vacuous | canary sweep | `bash .bench/gate.sh` — the canary phase reports `EXPECT is vacuous` if a timing line contains a fixture's expectation substring | The vacuity check compares against baseline inner-gate output, which timing lines now widen, so the whole existing fixture set grades the new format automatically |
| 5 | `bench prep-release` exits 0 and writes the artifact set plus `dist/preflight/release-index.json` | prep-release | `go test ./internal/contract/... -run TestPrepReleaseWritesEvidence` — red, the command does not exist | Asserting the named evidence file and the artifact count from `release-plan.mjs` rejects a stub that exits 0 without building anything |
| 5 | `bench prep-release` runs the cross-compile matrix | prep-release | `go test ./internal/contract/... -run TestPrepReleaseRunsStressTags` — red, nothing invokes the stress-tagged suite | `crossCompileMatrix` is a no-op without `-tags stress`, so an implementation that forgets the tag runs a check that silently returns nil |
| 6 | `bench prep-release` refuses with a nonzero exit and names `bench gate` when the verdict cache is absent, stale, red, or bound to a different tree | prep-release | `go test ./internal/contract/... -run TestPrepReleaseRequiresDevGreen` — red, no precondition exists | Enumerating all four cache states rejects an implementation that only checks for file existence, which would accept a red or tree-mismatched record |
| 6 | A seeded dev-tier failure reds both `bench gate` and `bench prep-release` | prep-release | `go test ./internal/contract/... -run TestDevTierFailureRedsBothTiers` — red, `prep-release` does not exist | This is the row that proves the split is a restaging and not a weakening; without the precondition, ship green could be claimed over a red tree |
| 7 | The dev canary sweep selects only dev-tier fixtures, and every fixture's declared tier matches the tier of the check its EXPECT names | canary sweep | `go test ./internal/canary -run TestFixtureTierMatchesCheckTier` — red, fixtures carry no tier | Deriving the expected tier from the promoted registry metadata rather than a second hand-written list means a fixture and its check cannot drift apart; the metadata promotion is what makes the import possible at all, since the registry currently lives in test files `internal/canary` cannot reach |
| 7 | The two probe-derived fixtures still bite under `prep-release` | canary sweep | `bash bin/bench.sh prep-release` — reports `did not bite` if the ship-tier canary step is skipped | A `prep-release` that runs the probe but never grades the fixtures leaves the checks unproven, which is the exact rot the canary exists to catch |
| 8 | `bench prep-release` is recorded in the subcommand routing registry and reaches `usage.Parse` | prep-release | `bash .bench/gate.sh` — `checkSubcommandRouting` fires `dispatches "prep-release" with no entry in the subcommand argument-routing registry`; canary fixture `unrouted-subcommand` already proves this check bites | The check reads dispatch names from `cmd/bench/main.go` itself rather than a fixed list, so it fails closed on the subcommand just added |
| 8 | `bench prep-release` reaches the same implementation through the kit CLI, the linked-repo by-path CLI, and hooks | prep-release | `go test ./internal/contract/... -run TestPrepReleaseAllSurfaces` — red, the route does not exist | The profile's hostile-input checklist names surface divergence explicitly, and a route added to `bin/bench.sh` alone would pass a single-surface test |
| 9 | `/bench-final-check` states the one-line ship-tier reminder on green | kit content surface | `bash .bench/gate.sh` — the workflow-guidance anchor check fires on the missing anchor, in the family the `acceptance-coverage-anchor` fixture already proves bites | The anchor sweep grades the shipped command file, so a reminder added only to a session's behavior and not to the command prose fails |
| 10 | `.bench/BENCH.md` lists `bench prep-release` | kit content surface | `bash .bench/gate.sh` — `checkColdPickupCLILists` fires `.bench/BENCH.md or .bench/BENCH-reference.md does not list CLI command 'bench prep-release'`; canary fixture `missing-cli-inventory` proves it bites | The check derives command names from `bin/bench.sh` routes, so adding the route without the doc row is red the moment the route lands |
| 10 | The project profile's gate section describes both tiers and what dev green claims | kit content surface | already covered — the profile is prose with no parser; graded by the `/bench-review-implementation` Standards axis | No honest gate signal exists for prose accuracy, and inventing an anchor that matches a keyword would grade spelling rather than correctness |
| 8 | The release path exits nonzero without current ship evidence | prep-release | already covered — `publication.VerifyPublishAuthority` enforces mode, scope, status, and profile, and canary fixture `behavior-owned/integrity-mismatch-acceptance` proves the digest cross-check bites | This spec builds nothing here; the assertable is listed so its coverage is a decision on the page rather than an assumed one |
| Edge of 5 | A required tool missing from PATH (`go`, `node`) fails `prep-release` closed, naming the tool | prep-release | `go test ./internal/contract/... -run TestPrepReleaseMissingTool` — red, no command | A thin route that shells out without checking will surface an opaque `exec: not found` rather than an actionable diagnostic |
| Edge of 5 | A repo root path containing spaces or glob characters runs `prep-release` correctly | prep-release | `go test ./internal/contract/... -run TestPrepReleaseHostilePath` — red, no command | The route passes the root to four shell scripts, and an unquoted expansion is the checklist's first-named failure for this domain |
| Edge of 5 | SIGINT partway through `prep-release` leaves no partially-written `release-index.json` | prep-release | `go test ./internal/contract/... -run TestPrepReleaseInterrupt` — red, no command | Evidence promotion is an atomic directory exchange, so this asserts the route does not bypass it with a direct write |
| Edge of 5 | A second `prep-release` on an unchanged tree succeeds and produces the same evidence | prep-release | `go test ./internal/contract/... -run TestPrepReleaseIdempotent` — red, no command | Re-run idempotency is on the profile checklist, and a route that appends to `dist/` rather than promoting into it fails on the second run |
| Edge of 5 | `prep-release` invoked from a directory deeper than the repo root resolves the root correctly | prep-release | `go test ./internal/contract/... -run TestPrepReleaseDeepCwd` — red, no command | The command assumes root for every script argument, and the checklist names deep-cwd invocation as a recurring defect here |
| Edge of 6 | `prep-release` writing `dist/` does not invalidate the dev-green verdict it read | prep-release | `go test ./internal/contract/... -run TestPrepReleaseDoesNotFalsifyItsOwnPrecheck` — red, no command | The verdict is tree-hash-bound and the checklist names self-falsifying writes explicitly; this asserts the write target stays gitignored so the subject is unchanged |
| Edge of 6 | `release-index.json` absent versus present-but-empty are distinguished | prep-release | already covered — `VerifyPublishAuthority` returns distinct errors for a missing index and an unparseable one | Both states already produce named diagnostics, and duplicating the assertion would be a second derivation of the same fact |
| Edge of 2 | A cold Go test cache does not make the dev gate exceed the 600 s package timeout | conformance driver | already covered by the story 2 row — the exclusion is what removes the failure mode, and there is no honest way to force a cold cache inside the gate without defeating it for every package | Forcing the cache cold to prove the exclusion would reintroduce the exact 10+ minute cost the exclusion exists to remove |

**Won't handle**

- Control bytes (ESC, BEL) in `prep-release` output — the command emits no git-sourced or user-supplied text, only fixed diagnostics and paths already constrained by the path row above.
- Non-TTY stdin on `prep-release` — the command never prompts by design (an approval prompt is a rejected alternative), so there is no interactive path to fail closed on.
- A FIFO, device, or dangling symlink at `dist/preflight/release-index.json` — `dist/` is gitignored build output written and promoted solely by the kit's own atomic exchange, so a hostile node there means the tree is already compromised beyond what this command can adjudicate.
- Invocation through a symlink rather than the real path — `scripts/release-preflight.sh` already resolves the kit root through symlinks, and `prep-release` reuses that resolution rather than adding a second one.
- WSL2 host-backed `fsync` stalls during artifact promotion — ambient host behavior with no assertable boundary; it makes `prep-release` slow, not wrong.
- Destructive worktree state during `prep-release` — the command neither creates nor registers worktrees, so the worktree lifecycle's fail-closed handling is untouched.

## Out of scope

- **Fifteen-check fan-out across the conformance driver** — killed by the 2026-07-26 data rather than deferred: one composite check is 99.8% of the phase, so a scheduler would recover nothing (3 edits, 2 gate runs, and negative expected value).
- **Hermetic build cache and verdicts keyed on the pinned gate subject** — a separate cache-infrastructure capability with its own oracle-semantics questions; revived only if re-measurement after this split still hurts (8 edits, 4 gate runs).
- **The `-count=1` freshness semantics, both levers** — removing the hardcoded `-count=1` on the *outer* gate phases, and adding one to the *inner* suite inside `checkGoCore` (the cache-leaning run). A reviewer-led oracle decision about what freshness means, not the rest of this feature; the measured 10+ minute price attaches to the inner lever, for which this split is the precondition (2 edits, 3 gate runs).
- **Running the ship tier on the pre-push hook** — a different enforcement boundary with its own latency budget; the hook stays dev tier (2 edits, 2 gate runs).
- **A build-approval prompt after a green final-check** — rejected on its merits rather than deferred: it nags at commit cadence, and story 9's one-line reminder is the decided alternative.
- **FT101's docs half and profile half (scoped ambient surfaces, `CONTEXT-MAP.md` layout)** — deferred undesigned behind its own revive trigger, a linked repo with more than one bounded context (unestimated; needs shaping first).
