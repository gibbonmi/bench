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
— (open)

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
— (open)

## #6: What does the tests/ bloat audit actually show?

Blocked by: —
Type: Research

### Question
Audit suite growth and duplication (pasted fixture harnesses, hand-rolled
subprocess capture) and propose a leanness discipline. Produces a short asset in
`decisions/assets/`; its evidence feeds the boundary choices in #2 and the
fixture scheme in #5.

### Answer
— (open)
