# Test-suite structure debt

`bench structure` caps a directory at 12 source files; `internal/conformance` sits
at 13 (all test files, ~2,500 lines) and `internal/contract` at 22 (~4,200 lines,
still growing — the in-flight binary-auto-repair spec adds another contract file).
The two caps squeeze jointly: the 400-line file cap forces splits, the splits feed
the 12-file dir cap. Folded in from the roadmap: the tests/canary rework (56 flat
fixtures) and the tests/ bloat audit — one debt-reduction program, per grill
2026-07-05.

## #1: What direction clears the crowded Go dirs?

Blocked by: —
Type: Grill

### Answer
Subpackages — group into real Go packages, which is what the check's own remedy
("group into modules") means in Go. Resolved 2026-07-05. Rejected: raising or
test-exempting the dir cap (oracle-weakening, and the caps would conflict again as
suites grow); consolidation-first (fights the split-by-responsibility layout that
just landed).

## #2: Where do the contract subpackage boundaries and shared helpers land?

Blocked by: —
Type: Grill

### Question
22 files share unexported helpers (`helper.go`, `runtime_helpers_test.go`,
`axi_helpers_test.go`); subpackages need an exported test-helper package. Which
family boundaries (runtime / axi+guards / doctor+link+package / repair?) and where
do helpers live? Does the parked "unify subprocess-capture seams" idea (conformance
`runProbe`, `Harness.Run`, canary `defaultRunner`) land inside this move or stay a
separate change? Timing: build after binary-auto-repair merges — its contract test
is in flight in the same package.

### Answer
Root lib + three family subpackages, resolved 2026-07-05. The root
`internal/contract` package stays the exported fixture-harness library:
`helper.go`/`helper_test.go` plus the generic assertions promoted out of
`runtime_helpers_test.go` into an exported assert surface. Consumer subpackages:
`runtime/` (status, gate, git, shift, shift adapters/helpers, worktree,
structure), `axi/` (axi, wave2, guards, fail_closed, axi asserts), `surface/`
(doctor, doctor_shim, link, link_marker, package, go_routing, binary_repair).
Family-specific helpers stay unexported in their family package; only generic
ones live in the root lib. The relocation is behavior-preserving — same tests,
same subtest counts — and builds only after binary-auto-repair merges. The
"unify subprocess-capture seams" idea stays a separate parked change (bigger
blast radius: it spans conformance and canary). Rejected: one-dir-per-family
(2-file fragment packages); moving only `runtime/` out (root stays at the cap
and re-triggers as suites grow).

## #3: How does the oracle rewire for subpackages?

Blocked by: #2
Type: Grill

### Question
`.bench/gate.sh` runs `go test ./internal/contract` with `BENCH_CONTRACT_ROOT`,
and `checkGoCore` excludes exactly `internal/contract` (suffix match) from the
conformance unit sweep. Subpackages like `internal/contract/axi` would escape the
exclusion and silently resurrect the duplicate contract run the test-layout change
removed. Prefix-match exclusion vs deriving the exclusion from the gate's phase
list; does the contract phase become `./internal/contract/...`? This edits what
the oracle covers — craft-gate discipline applies.

### Answer
Three touchpoints in one change, resolved 2026-07-05. (1) The gate's contract
phase becomes `go test -count=1 ./internal/contract/...` — root lib plus
subpackages, still one phase, still one run. (2) The `goCoreTestPackages`
exclusion widens from exact/suffix match to the whole `internal/contract/`
subtree, so subpackages cannot escape and resurrect the duplicate run. (3) The
`gate_entry` conformance anchor pins the gate line verbatim and updates in the
same change — it is the existing bite-proof for the gate edit. No-weakening:
the contract suite still runs exactly once; the unit sweep loses only packages
the contract phase itself owns. Proof for the widened exclusion: a unit test on
`goCoreTestPackages` asserting contract-subtree packages are excluded and a
sibling (e.g. `internal/conformance`) survives — it catches the too-wide
direction, the only one that loses coverage; too-narrow is cost-only and
self-announces via gate wall time. Rejected: a new canary fixture (the sprawl
#5 is shrinking; the unit test catches the same bug cheaper); relying on the
existing canary alone (silent on too-wide over-filtering).

## #4: Does conformance subpackage too, or fit under the cap another way?

Blocked by: —
Type: Grill

### Question
Only one file over the cap, and the five check families are already named in the
file layout — but `RunConformance` and the gate's
`-run '^TestRootConformance$'` entry live in the single package, so subpackaging
changes the gate's entry contract. Merging two helper files clears the cap today
but leaves the same squeeze to recur. Same-direction split as #1, or the one
place a lighter fix is honest?

### Answer
Lighter fix, resolved 2026-07-05: merge `harness_helpers_test.go` (145 lines)
into `harness_test.go` (88) and fold `root_conformance_test.go` (17) into
`gate_entry_test.go` — 13 files → 11, every merged file under the 400-line cap.
Conformance is one deep module by design (single check registry, single pinned
gate entry `-run '^TestRootConformance$'`), so subpackaging would fragment a
registry whose point is to be one list and rewire the gate for zero coverage
gain. Not a reopen of #1's consolidation-first rejection: these are
helper-beside-single-consumer merges (cohesion), not responsibility
consolidation. `docs_workflow_helpers_test.go` stays separate — two consumers,
and merging it into its checks file would break the line cap. If a genuine new
check family pushes the dir over the cap again, the #1/#2 subpackage direction
applies then; don't pre-build. Rejected: subpackaging now (fragments the
registry, changes the gate entry contract); raising the cap (oracle-weakening,
per #1).

## #5: What structure replaces 56 flat canary fixtures?

Blocked by: —
Type: Grill

### Question
`tests/canary/` is 56 flat fixture directories, hand-tended. Group by check
family, or generate fixtures from the check registry? And does the
one-canary-per-check rule become a meta-check instead of manual bookkeeping?
Registry-generation risks a second source for what a check needs to bite;
grouping alone may just move the sprawl.

### Answer
Group by family with path-as-classification, plus a meta-check; resolved
2026-07-05 on the #6 audit evidence (`decisions/assets/tests-bloat-audit.md`).
Layout becomes `tests/canary/<family>/<fixture>/`; the family column of
`canaryFixtureRegistry` is deleted and derived from the directory, collapsing a
59-entry hand-tended map into the tree (the registry keeps only the
retired-shell-source column). `Sweep` walks two levels. Fixture content stays
hand-tended. One-canary-per-check becomes a meta-check: every check family in
`RunConformance` must own ≥1 fixture directory; the rename-coupling is intended
(a family rename must move its fixtures). Rejected: generating fixtures from
the check registry (`files/` is the sole expression of each biting input —
generation manufactures a second source); flat layout + registry map as-is
(keeps 59 hand-maintained classification entries); grouping without deriving
family from path (moves sprawl, dedups nothing).

## #6: What does the tests/ bloat audit actually show?

Blocked by: —
Type: Research

### Question
Audit suite growth and duplication (pasted fixture harnesses, hand-rolled
subprocess capture) and propose a leanness discipline. Produces a short asset in
`decisions/assets/`; its evidence feeds the boundary choices in #2 and the
fixture scheme in #5.

### Answer
Audited 2026-07-05; asset at `decisions/assets/tests-bloat-audit.md`. Headlines:
canary fixtures are already lean (~1.6 KB each, near-zero identical scaffolding
— honest repetition, not bloat); generating fixtures from the check registry
would create a second source, since `files/` is the sole expression of each
biting input; one-canary-per-check is manual, not enforced. Real duplicated
knowledge found: the exit-code extraction idiom re-derived 3× across the
subprocess seams, 3-4 independent temp-git-repo scaffolds, and the retired
gate-fragment list derived in two files. Growth is front-loaded (test estate is
the repo's fastest-growing part). The asset ends with four candidate leanness
disciplines for the reviewer; adopting any is decided in #5 (fixtures) or a
future map, not here.

## Handoff

1. **Module boundaries.** Three independent slices. (a) Contract: root
   `internal/contract` becomes the exported fixture-harness lib (helper +
   generic asserts promoted from `runtime_helpers_test.go`); consumer
   subpackages `runtime/`, `axi/`, `surface/` per the #2 file assignment;
   family-specific helpers stay unexported in their family. (b) Conformance:
   stays one package; two merges (#4) take it 13 → 11 files. (c) Canary:
   `tests/canary/<family>/<fixture>/`, family derived from path; `Sweep` walks
   two levels; `canaryFixtureRegistry` keeps only the retired-shell-source
   column; new meta-check requires ≥1 fixture dir per `RunConformance` family.
2. **Contracts.** Gate contract phase becomes
   `go test -count=1 ./internal/contract/...` — still one phase, one run.
   `goCoreTestPackages` excludes the whole `internal/contract/` subtree.
   Root contract lib's exported API is the harness + asserts consumed by the
   three subpackages; relocation is behavior-preserving (same tests, same
   subtest counts). `bench canary` CLI behavior (exit codes, diagnostics,
   vacuous-EXPECT rejection, bite semantics) is unchanged; only fixture
   addressing gains a family level. Meta-check emits a conformance diagnostic
   for a family with zero fixtures.
3. **Deep vs thin.** Deep: the contract root lib (hides fixture/env-isolation
   complexity; seam at its exported API) and canary `Sweep` (hides
   walk/materialize/parallelism; seam at `bench canary`). Thin: the family
   subpackages and the conformance merges — pure consumers/relocations, no
   seam of their own.
4. **Black-box assertables.** `bench gate` exit 0 on the restructured tree
   (the whole program is behavior-preserving). Unit test on
   `goCoreTestPackages`: contract-subtree package excluded, an
   `internal/conformance` sibling survives (#3's bite proof). Meta-check bite:
   a family with zero fixture dirs turns conformance red. `bench structure`
   reports zero dir-cap issues after. Subtest-count equality before/after the
   contract relocation.
5. **Gate attachment.** The gate line edit is pinned by the `gate_entry`
   conformance anchor and updates in the same change (#3). The meta-check
   lands inside `RunConformance`, observed via `TestRootConformance`. Not
   gate-visible: the "same subtest counts" relocation claim — capture the
   `go test -v` name inventory before/after during the build (one-time
   verify), it has no permanent seam.
6. **Hostile-input owners.** Mostly n/a — internal Go restructure, no new
   shell surface. Live classes: absent vs empty (family dir with zero
   fixtures → meta-check; fixture missing EXPECT/`files/` → existing
   `runFixture` diagnostics, unchanged); paths with spaces/globs in fixture
   names → `Sweep`'s two-level walk (kit-authored names, but the walk must not
   glob). Remaining checklist classes: n/a — no CLI/link/init surface touched.
7. **Uncertainty flags.** None — no seam needs escalation; spec on the mid
   tier per profile.
8. **Rejected alternatives.** Raising/test-exempting the dir cap;
   consolidation-first (#1). One-dir-per-family contract split; moving only
   `runtime/` (#2). New canary fixture for the exclusion proof; trusting the
   existing canary alone (#3). Subpackaging conformance now (#4). Generating
   fixtures from the check registry; flat layout + hand-tended map; grouping
   without path-derived family (#5). Subprocess-capture-seam unification stays
   a separate parked change (#2/#6 — bigger blast radius).
9. **Domain watch-outs.** Merged Go test files share package scope — helper
   symbol collisions must be resolved at merge time. `runFixture` names a
   fixture by `filepath.Base`, so the family level removes filesystem-enforced
   global uniqueness of fixture names — diagnostics and the registry key may
   need family-qualified names. The contract relocation must land as move-only
   commits, no behavior edits mixed in, so review can grade it as relocation.
   `BENCH_CONTRACT_ROOT` must reach subpackage tests unchanged.

Dependency order: #4 conformance merges first (touches `gate_entry_test.go`
before #3 re-pins the anchor there), then #2+#3 contract split + gate rewire
(one change per #3), then #5 canary regroup (independent; parallel-safe).
