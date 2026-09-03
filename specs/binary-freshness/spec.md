# Binary freshness

Status: staged

Roadmap: FT177

Decision source: roadmap row FT177, carrying its "Decided 2026-09-02" paragraph from the reviewer's grill of that date

Verification log: 2 iteration(s) to accept — the first round blocked on two missing test owners, a partial fence, two unordered edges, an import cycle, and a publish signature; the second pass found one fence line

## Problem

A stale `dist/bench` yields a false pass wherever a hand run, a hook, the
wrapper, or the landing executes it. The gate never touches that binary, so
nothing in the ordinary loop reds it. `bench doctor` reports healthy beside a
stale binary and a stale broker, because it has no row for either. The
follow-on guard turns a stale binary's `unknown subcommand` answer into a
denial of every Bash call. A hand rebuild replaces the broker bytes without a
manifest digest, and the next landing refuses with exit 127.

The landing tail then names `bench repair`, which refuses in the kit source
checkout because the pin manifest exists only in the release artifact. The
system suite's hand-run route is unnamed in guidance, so a probe fails on
plumbing before it reaches its subject.

## Solution

Every consumer of `dist/bench` grades it through the one existing digest
primitive in `internal/freshness`, never through an mtime. The stamped build
publishes the executable, its seal, and the broker manifest as one outcome. The
land route recovers a `dev` manifest through that build and keeps a digest
mismatch as a tamper refusal. The guard classifies a stale answer and fails
open for a non-Bench call. `bench doctor`, `bench commands --brief`, and the
`--full` preflight each gain the freshness verdict with the one rebuild
sentence. The landing line picks its route by checkout kind.

The implementation of this spec starts after the `roadmap-light-path-fixes-3`
landing, because both edit the profile.

## User stories

### Group A — the freshness verdict reaches its consumers

Line: opus / low. The digest, the decision, and the rebuild sentence exist;
the work wires them into three consumers.

1. As a session, I want `bench commands --brief` to refuse when the binary's seal does not match the sources, so that the liveness probe cannot pass stale.
2. As a session, I want that refusal to print the one rebuild sentence, so that the next action is exact.
3. As an operator, I want `bench doctor` to print a row for the `dist/bench` seal, so that a stale binary is visible before a landing.
4. As an operator, I want `bench doctor` to print a row for the broker manifest, so that a refusal at exit 127 is predicted, not met.
5. As an operator, I want a landing refusal to name the kit root as the rebuild root, so that the printed command runs here.
6. As a maintainer, I want every freshness verdict derived from the source digest, so that a checkout that rewrites mtimes stays green.

### Group B — the inherited binary is verified

Line: opus / low. The verifier exists; one branch and one owner change.

7. As a test author, I want an inherited `BENCH_RUN_BINARY` verified against a named source root, so that a self-consistent stale seal cannot pass.
8. As a test author, I want the system suite owner to verify the binary against `BENCH_KIT`, so that a stale hand run reds at setup.
9. As a gate operator, I want the gate's own private build unchanged, so that the ordinary loop keeps its exact-source binary.

### Group C — the stamped build publishes the broker

Line: opus / medium. The publish transaction is a rollback owner, and the
manifest lands in a second directory.

10. As an operator, I want the subject-mode build to write the broker manifest in the seal's publication, so that no crash parts them.
11. As a release maintainer, I want artifact mode to keep never executing what it built, so that the release path stays inert.
12. As a maintainer, I want the manifest publication to stay inside the two files the topology test allows, so that no third caller appears.
13. As an operator, I want a manifest written by the build to carry the stamped version, not `dev`, so that the land route accepts it.

### Group D — the land route recovers a dev manifest

Line: opus / medium. The route is shell before any repository read, and the
recovery must stay single-pass.

14. As an operator, I want a manifest whose version reads `dev` to trigger one stamped rebuild and one re-read, so that the landing continues.
15. As an operator, I want a digest mismatch to keep its exit 127, so that a tampered broker never recovers.
16. As an operator, I want the recovery to read no repository and honor no inherited override, so that the route keeps its trust boundary.
17. As an operator, I want a second `dev` after the rebuild to refuse, so that the recovery cannot loop.

### Group E — the landing line and the guidance

Line: opus / medium. Guidance prose steers later sessions, and the landing
line is one composed string.

18. As an operator in the kit checkout, I want the landing line to name the stamped rebuild and `bench doctor --fix`, so that its route works here.
19. As an operator on an installed kit, I want the landing line to name `bench repair`, so that the installed route stays advertised where it works.
20. As a maintainer, I want the landing line composed from the checkout predicate and the rebuild sentence, so that no second copy of either exists.
21. As a session, I want the working agreement and the profile to name `bench test --check system` as the hand-run route, so that a probe reaches its subject.

### Group F — the guard on a stale answer

Line: opus / medium. The shim must classify without the binary, and the
one-source hazard is real.

22. As a session, I want a non-Bench call to pass with a warning on a stale answer, so that the shell stays usable.
23. As a session, I want a Bench call to refuse with the rebuild sentence on that answer, so that no Bench verb runs stale.
24. As a session, I want a genuine refusal, with its `BLOCKED:` line, to keep exit 2, so that the pool denial does not fail open.
25. As a maintainer, I want the shim's Bench-call word test pinned against the Go classifier's fixture table, so that the two cannot drift silently.

### Group G — the preflight seal row

Line: opus / low. One new check row in an existing table.

26. As a coordinator, I want `bench preflight build` to red on a mismatched destination seal, so that the `--full` path names the rebuild before the transaction.
27. As a coordinator, I want a linked repo with no `dist/bench` to report that row not applicable, so that a consumer build is not refused.

## Implementation decisions

- Staleness is `freshness.Select` over the seal's source digest and `freshness.Digest` of the current build inputs. No consumer reads an mtime, and no consumer re-derives the digest.
- `freshness.RebuildAction` is the one rebuild sentence. Every new refusal or warning calls it; none inlines the text.
- `bench commands --brief` becomes a root-taking handler. It resolves the running executable, verifies it against the kit root, and refuses with the sentence on a mismatch. Outside a repository it keeps its three-line answer.
- `bench doctor` gains two rows: the `dist/bench` seal verdict, and the broker manifest graded by the same five predicates the land route applies. The doctor row and the land route are two derivations by necessity, because the route runs before any binary is trusted. One authored conformance expectation pins that both enumerate the same five reasons.
- The broker manifest reader and writer move to a new leaf package `internal/brokermanifest`, which `internal/adopt` and `internal/freshness` both import. Without it the doctor row and the publish transaction form an import cycle.
- The landing's gate refusal takes a repair root distinct from the digest root, threaded through a new `Factory` field. The prospective owner passes the kit checkout as the repair root.
- `kitSourceCheckout` is exported from `internal/adopt`, so the landing line can call it.
- `Factory.validate` verifies an inherited executable against its source root when the caller names one. The system suite owner names `BENCH_KIT`. The gate's private build path is unchanged.
- The subject-mode build's `freshness-publish` step writes the broker manifest inside the same publication transaction, with the stamped version. The verb gains two arguments, the published path and the version, so the manifest binds the published executable and not the staged one. Artifact mode stays excluded. The manifest lands beside the resolved wrapper, as today.
- The shim's shared classifier table holds resolver-independent rows only, because `benchguard.InvokesBench` takes a resolver the shell test lacks.
- The land route splits refusal three from refusal five. The version branch runs the stamped build at the install root once, re-reads the manifest once, and refuses if the version still reads `dev`. The digest branch keeps its unconditional exit 127. The route reads no repository and honors no inherited override.
- The landing line composes `kitSourceCheckout` with `RebuildAction` for the source checkout, and names `bench repair` for an installed kit.
- The shim captures the child's stderr. When the child exits 2 and its stderr holds `unknown subcommand`, the shim runs a shell word test for a Bench call. A Bench call refuses with the sentence; any other call passes with a warning that carries the sentence. A genuine refusal keeps exit 2. The shell word test is a second derivation of `benchguard.InvokesBench`, pinned by one shared fixture table both sides read.
- `bench preflight build` gains a `binary-seal` row. It verifies the destination `dist/bench` when present and reports not applicable when absent.
- No wrapper script for the system suite is added. `bench test --check system` is the one seam, and the guidance names it.

## Testing decisions

- Consumer verdicts use the existing freshness decision tests plus one behavior test per consumer.
- The inherited-binary branch uses the runbinary bootstrap tests and one system-suite setup test.
- The broker publication uses the publish rollback tests, the broker tests, the topology test, and the land-route system test.
- The land route uses the land-route refusal table with a `dev` row that expects success and a digest row that keeps 127.
- The guard uses the degraded-rim fixture with a fake binary that answers `unknown subcommand`.
- The preflight row uses the preflight decision tests.
- Guidance uses anchor tuples and live-mirror fixtures in the five-part shape.

### Seam diagram

    scripts/go-build.sh (subject mode)
        │
        ▼
    staged binary ──▶ [ freshness-publish ] ──▶ dist/bench + .seal + bench-broker.manifest
                              ◀ tests attach here: publish rollback, broker read, topology

    consumer (commands --brief, doctor, preflight, guard, land route)
        │
        ▼
    executable ──▶ [ freshness.Verify / Select ] ──▶ verdict + RebuildAction
                              ◀ tests attach here: per-consumer behavior tests

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| BF1 | 1, 2 | `bench commands --brief` under a mismatched seal exits non-zero and prints the `RebuildAction` sentence | new `cmd/bench` behavior test | the fixed three-line string passes stale |
| BF2 | 1 | `bench commands --brief` under a matching seal prints the three-line answer | new `cmd/bench` behavior test | a refusal on every run breaks the liveness probe |
| BF3 | 3 | `bench doctor` prints a seal row whose verdict names the source-digest mismatch | new doctor rows test | nine rows report healthy beside a stale binary |
| BF4 | 4 | `bench doctor` prints a broker row that names each of the five land-route reasons | new doctor rows test and the authored reason-list expectation | a row with four reasons predicts only part of the 127 |
| BF5 | 5 | the landing's inherited-binary refusal prints the kit root in its rebuild command | new `TestProspectiveOwnerRefusalNamesTheKitRoot` in the prospective owner tests | a refusal that prints the composed tree names a root nobody can use |
| BF6 | 6 | `freshness.Select` decides the seal by content with equal mtimes, as a regression row | `TestVerifyUsesContentRatherThanMtime` | an mtime comparison reds a fresh checkout |
| BF7 | 7 | an inherited executable with a stale source digest and a named source root refuses in `Factory.validate` | new `TestFactoryValidateRefusesAStaleInheritedSeal` beside the validate tests | the seal-pair check alone passes a self-consistent stale binary |
| BF8 | 8 | the system suite owner refuses a `BENCH_RUN_BINARY` whose seal mismatches `BENCH_KIT` | new out-of-process test that runs the suite binary with a stale seal and reads its exit | `identifyExecutable` trusts any executable file |
| BF9 | 9 | `runbinary.Own` still builds and verifies a private executable, as a regression row | `TestFactoryOwnBuildsOnePrivateAbsoluteSelectionAndCleansIt` | a changed owner path reds the gate loop |
| BF10 | 10 | after a subject-mode build, the manifest digest equals the published executable's digest | new publish test | two commands leave the manifest one crash behind |
| BF11 | 10 | a publication rolled back leaves no manifest change | publish rollback tests | a manifest written before the rename survives the rollback |
| BF12 | 11 | artifact mode writes no manifest and executes nothing | new `TestGoBuildArtifactModeWritesNoManifest` in `cmd/bench` | a manifest write in artifact mode runs the new binary |
| BF13 | 12 | only `cmd/bench/freshness_publish.go`, `cmd/bench/main.go`, and `scripts/go-build.sh` carry the publish token | `TestFreshnessPublicationTopology` | a fourth caller reds the topology |
| BF14 | 13 | the manifest written by the build carries the package version | new broker test | a `dev` manifest re-creates refusal three |
| BF15 | 14 | a `dev` manifest with a sound digest lands after one rebuild and one re-read | `TestWorktreeLandRouteRefusesEveryUnauthenticatedBroker` new row | a route that refuses `dev` keeps the 2026-08-29 occurrence |
| BF16 | 15 | a digest mismatch exits 127 with the repair advice | `TestWorktreeLandRouteRefusesEveryUnauthenticatedBroker` digest row | a route that recovers the digest branch rubber-stamps tampered bytes |
| BF17 | 16 | the recovery with an inherited `BENCH_KIT` exits 1 before any rebuild | the land-route inherited-override rows | a recovery that honors the override reads the wrong tree |
| BF18 | 17 | a manifest still `dev` after the rebuild exits 127 | new land-route row | a recursive retry loops |
| BF19 | 18, 20 | in the kit source checkout the landing line names the `RebuildAction` sentence and `bench doctor --fix` | `TestLandCommandReportsInstallStepForABrokerChangingDiff` | a line that names `bench repair` here advertises a refusing route |
| BF20 | 19 | on an installed kit the landing line names `bench repair` | new land freshness test with the predicate false | one unconditional line loses the installed route |
| BF21 | 21 | the working agreement and the profile name `bench test --check system` | anchor registry test and fixtures `agents-system-suite-route` and `benchkit-system-suite-route` | an unnamed route repeats the plumbing failure |
| BF22 | 22 | a fake binary that answers `unknown subcommand` with exit 2 lets `ls` pass at exit 0 with a warning that holds the sentence | `TestBenchFollowOnHookDegradedRim` new rows | a widened rc case fails open on every refusal |
| BF23 | 23 | the same fake binary refuses `bench gate` at exit 2 with the sentence | `TestBenchFollowOnHookDegradedRim` new rows | a shim that passes every call runs a stale verb |
| BF24 | 24 | a fake binary that prints `BLOCKED:` with exit 2 keeps exit 2 for `ls` | `TestBenchFollowOnHookProcess` refuse row | a shim that treats every 2 as stale fails open on the pool denial |
| BF25 | 25 | the shell word test and `benchguard.InvokesBench` agree on every row of one shared fixture table | new conformance test over the shared table | two classifiers drift with no red |
| BF26 | 26 | `bench preflight build` reds `binary-seal` when the destination seal mismatches | new preflight decision test | a `--full` run enters the transaction on a stale binary |
| BF27 | 27 | a root with no `dist/bench` reports `binary-seal` not applicable | new preflight decision test | a refusal on absence blocks every consumer build |

### Edge inventory

- Error paths: a missing seal, a malformed seal, and a digest mismatch each refuse with the sentence.
- Empty input: no `dist/bench` at all is not applicable for preflight and a plain row for doctor.
- Boundary values: equal mtimes never decide; a `dev` version recovers once and only once.
- Interrupted state: a publication interrupted before the rename leaves the old executable, seal, and manifest.
- Re-run idempotency: a second `bench doctor` after a rebuild prints two green rows.
- Hostile paths: a symlinked shim resolves to the real executable before verification.
- Partial implementation: a consumer that inlines the rebuild text reds the one-source sweep the topology test already runs.

**Won't handle** — a wrapper script for the system suite — `bench test --check system` is the one seam, and a script would be a second producer of its operands.

**Won't handle** — the shell classifier reimplementing symlink resolution — the word test covers the first word of each segment, and the shared table pins its reach.

**Won't handle** — `bench repair` reading a pin manifest in the kit source checkout — the decision keeps `bench repair` as the installed-kit route.

**Won't handle** — a tamper refusal that recovers — the digest branch keeps exit 127 by decision.

**Won't handle** — a sweep for bare `go build` paths — the tree holds `go build` only inside the stamped build script, and BF13 pins the publish token.

## Ownership fences

- `specs/binary-freshness/`
- `reviews/binary-freshness.md`
- `cmd/bench/main.go`
- `cmd/bench/main_test.go`
- `internal/adopt/doctor_rows.go`
- `internal/adopt/doctor.go`
- `internal/adopt/broker.go`
- `internal/adopt/broker_test.go`
- `internal/adopt/setup_test.go`
- `internal/brokermanifest/`
- `cmd/bench/freshness_publish.go`
- `cmd/bench/freshness_publish_test.go`
- `cmd/bench/build_artifact_mode_test.go`
- `cmd/bench/build_subject_mode_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `internal/freshness/`
- `internal/runbinary/`
- `internal/gate/prospective.go`
- `internal/gate/prospective_owner_test.go`
- `internal/systemtest/owner_test.go`
- `internal/systemtest/owner_stale_seal_test.go`
- `internal/systemtest/land_route_test.go`
- `internal/systemtest/owner_artifact_recovery_test.go`
- `internal/systemtest/bench_follow_on_test.go`
- `internal/worktree/land.go`
- `internal/worktree/land_freshness_test.go`
- `internal/preflight/decision.go`
- `internal/preflight/gather.go`
- `internal/preflight/decision_test.go`
- `internal/preflight/source_tip_test.go`
- `internal/benchguard/`
- `internal/conformance/`
- `internal/conformance/registry/registry.go`
- `bin/bench.sh`
- `scripts/go-build.sh`
- `.bench/hooks/block-bench-follow-on.sh`
- `.bench/lib/resolve-bench.sh`
- `AGENTS.md`
- `projects/benchkit.md`
- `internal/anchors/registry_data.go`
- `cmd/bench/testdata/anchors/pre-disclosure-populated.stdout`
- `internal/anchors/registry_data_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `tests/canary/guard-classifier-table/`
- `tests/canary/load-validity-metadata/shared-rule-drift`
- `tests/canary/load-validity-metadata/extensionless-gate-ref`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/docs-currency-token-diet/missing-cli-inventory`
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference`
- `tests/canary/docs-currency-token-diet/stale-skill-cli-reference`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/package-core-guard/bounds-duplicate-owner`
- `tests/canary/package-core-guard/reintroduced-bare-skip`
- `tests/canary/package-core-guard/guard-resolver-order-drift`

The fence is the union of the nine tickets' `Writes:` lines, closed by
`bench preflight build` over the fixture and registry pins. A closure headroom
file creates no blocker edge; only a file a ticket's `What to build` names
does. Four real edges follow. The publish ticket and the landing line follow the
doctor ticket, and commands-brief and the land route follow the publish
ticket.

## Out of scope

- `bench repair` in the kit source checkout: 0 edits by decision.
- A wrapper script for the system suite: 0 edits by decision.
- A second `bench doctor --fix` action that rebuilds `dist/bench`: 2 edits, 1 gate run.

## Further notes

Recorded during the build, open to reviewer veto:

- The guard classifier table gets a new `familyChecks` row in
  `internal/conformance/registry/registry.go`. A canary fixture bites only
  through a registered family owner, and every existing family is bound
  there.
- The guard classifier table's drift check binds to the existing
  `package-core-guard` owner through one call line in
  `internal/conformance/package_core_checks_test.go`. The check parses the
  table's rows from the graded root. It then runs the shell word test over
  each one, so the fixture bites under any graded tree.
- The shell word test lives in `.bench/lib/resolve-bench.sh` as
  `bench_invokes_bench`, not in the shim. The conformance test sources the
  lib, and the shim keeps the posture.

- The doctor broker row names an absent manifest with an `ok:` verdict, not a
  red one. The kit source checkout publishes no broker manifest beside its
  resolved wrapper, and its landings still work. So a red there is a false
  alarm for every kit session. A red also reds four tests outside the fence.
  The row still carries refusal one's wording, so the reason list stays one
  source.
- `internal/conformance/compliance_checks_test.go` calls the new reason-list
  expectation from the registered `kit-compliance` owner. A standalone
  live-tree test reds `TestConformanceMetaBites` as unregistered.
- `cmd/bench/testdata/anchors/pre-disclosure-populated.stdout` joins the
  fence. It is a captured snapshot of the AGENTS.md anchor list, so a new
  anchor moves it.

- The `binary-seal` row is gated to build mode, because story 26 and rows
  BF26 and BF27 name `bench preflight build`. An unconditional row also reds
  `bench preflight review` on a stale `dist/bench`, which no source sentence
  asks for.
- `internal/preflight/source_tip_test.go` joins the fence. It asserts a
  literal preflight row count for both modes, so the new row moves its
  build-mode expectation. `internal/preflight/command_review_test.go` stays
  outside the fence, because the build-mode gate keeps its review-mode count
  at eleven.

Recorded during the build for the publish ticket, open to reviewer veto:

- `freshness-publish` gains one argument, the version, not two. The published
  path is already argument two, so a second copy of it would be a second
  source of one fact. The manifest directory is derived from that path.
- `freshness.Publish` gains a version parameter, so five call sites take a
  one-argument edit. `internal/systemtest/owner_artifact_recovery_test.go`
  joins the fence, because its fake builder calls the verb.
- A coordinator probe found a missing row. The build script's version
  argument was ungraded, because BF14's test calls `freshness.Publish`
  directly. `TestGoBuildSubjectModePublishesTheStampedVersion` now runs a real
  subject-mode build and grades the published manifest.
- One mutation came back silently green: a manifest write that stays in place
  but bypasses the transaction's step lock. It is behavior-equivalent for the
  rollback, and only a signal-timing test in the manifest's own rename window
  could see it. No row was added.

Flagged additions beyond the decision source:

- The `binary-seal` preflight row. The row's text names a `--full` landing preflight check; this is its executable form.
- The shared classifier fixture table between the shim and `benchguard`. The decision names the guard posture; the table is what keeps two classifiers honest.
- The authored reason-list expectation between the doctor broker row and the land route. The route is shell before trust, so the row cannot call it.
- Two guidance sentences that name `bench test --check system`. The row asks for one wrapped invocation; the verb exists, and the gap is the guidance.
- The `internal/brokermanifest` leaf package. The doctor row and the publish transaction would otherwise form an import cycle.
- The BF9 regression row on the gate's private build. The row guards a no-change, and the source names no such sentence.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| staleness is a stamped source-digest mismatch, never an mtime | BF1, BF3, BF6, BF7, BF8 |
| the guard warns and passes a non-Bench call, refuses a Bench call with the rebuild named | BF22, BF23, BF24, BF25 |
| the stamped build publishes the broker atomically; `bench doctor --fix` stays the repair | BF10 to BF14 |
| a `dev` manifest rebuilds and continues; a digest mismatch stays a tamper refusal | BF15 to BF18 |
| the landing line names the stamped rebuild and `bench doctor --fix` for the source checkout; `bench repair` stays the installed-kit route | BF19, BF20 |
| `bench commands --brief` verifies identity or refuses with a repair action | BF1, BF2 |
| `bench doctor` names the mismatch and its fix, including the worktree-built case | BF3, BF4, BF5 |
| the `--full` landing preflight checks the destination binary seal | BF26, BF27 |
| the tagged systemtest suite gets one wrapped invocation | BF21 |
| the gate's private build is unchanged | BF9 |

The 2026-09-02 grill closed every fork on the row. Every subagent runs `opus`
at low or medium effort.
