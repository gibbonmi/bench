# Gate-pipeline fixture inventory (asset for decisions/gate-pipeline.md #6)

Research output, 2026-07-26: every claim cited to the code that emits the
diagnostic the fixture's EXPECT matches. Composite entry:
`checkPackageCoreAndGuards` (`internal/conformance/package_core_checks_test.go:20`);
`checkGoCore` (`:180`).

## Part A — package-core-guard family (39 fixtures), post-split buckets

Bucket totals: build-phase 1, test-phase 1, residual-structural 28, ship 2,
stray (other check) 7. gofmt/vet/race/cross-compile: **zero fixtures today**.

| fixture | EXPECT stem | emitter | bucket |
|---|---|---|---|
| go-build-broken | `go build failed` | package_core_checks_test.go:205 | build-phase ⚠ |
| go-test-failing | `go test failed` | package_core_checks_test.go:216 | test-phase ⚠ |
| missing-files-entry | `files[] missing` | package_core_checks_test.go:52 | residual |
| kit-only-allowlist-emptied | `consumer payload allowlist declares no kit-only rows` | package_core_checks_test.go:105 | residual |
| kit-only-asset-admitted | `npm package includes kit-only allowlist asset` | package_core_checks_test.go:118 | residual |
| guard-describe-boundary-dropped | `manifest missing boundary` | package_core_checks_test.go:488 | residual |
| guard-resolver-order-drift | `git guard inlined resolver order drifts…` | guard_resolver_drift_test.go:218 | residual |
| native-trigger-comment-spoof | `native verification workflow does not include pull requests` | workflow_checks_test.go:53 | residual |
| native-reproducibility-handoff-dropped | `…does not hand reproducibility records…` | workflow_checks_test.go:82 | residual |
| native-smoke-workflow-dropped | `…does not run smoke from finalized evidence` | workflow_checks_test.go:87 | residual |
| preflight-verify-{analysis,artifact,gate,smoke,vulnerability}-omitted | `release preflight verify registry omits or reorders…` | native_workflow_test.go:123 | residual (5) |
| preflight-publish-{ancestry,changelog,identity}-omitted | `release preflight publish registry omits or reorders…` | native_workflow_test.go:126 | residual (3) |
| release-future-owner-omitted | `release requirement registry omits public.ft88.data_handling` | native_workflow_test.go:118 | residual |
| preflight-native-{call,upload}-bypassed | `native verification does not finalize full release evidence…` | native_workflow_test.go:139 | residual (2) |
| preflight-release-call-bypassed, release-public-profile-omitted | `tag publication bypasses full release preflight` | native_workflow_test.go:146–148 | residual (2) |
| preflight-publish-needs-bypassed | `publication does not wait for finalized evidence…` | native_workflow_test.go:155 | residual |
| reproducibility-byte-compare-bypassed | `reproducibility comparator does not require exact byte equality` | native_workflow_test.go:169 | residual |
| offline-network-repair-allowed | `offline smoke permits repair or network fallback` | native_workflow_test.go:173 | residual |
| offline-stage-interruption-ignored | `offline smoke omits stage interruption recovery` | native_workflow_test.go:176 | residual |
| offline-slice1-operation-omitted | `offline smoke omits slice-1 suppressed operation proof` | native_workflow_test.go:179 | residual |
| offline-registry-fallback-allowed | `offline registry smoke does not fail closed` | native_workflow_test.go:182 | residual |
| release-digest-binding-omitted (CHECK) | `release index does not bind component manifest digests` | release_probe_fixture_test.go:62 | ship |
| release-package-evidence-omitted (CHECK) | `package evidence is missing or unsafe: LICENSE` | scripts/build-release-evidence.mjs:159 via native_workflow_test.go:247 | ship |
| bounds-canary-width-unconsumed | `…does not consume bounds.CanaryInnerWidth` | bounds_policy_test.go:45 | stray → bounds-policy |
| bounds-duplicate-canary-width | `…redeclares CanaryInnerWidth policy value` | bounds_policy_test.go:115 | stray → bounds-policy |
| bounds-duplicate-owner | `internal/models/models.go redeclares ProviderTimeout…` | bounds_policy_test.go:115 | stray → bounds-policy |
| default-branch-refabricated | `declares or calls a DefaultBranch function…` | git_facts_checks_test.go:40 | stray → default-branch-single-source |
| marker-wait-literal-deadline | `…duration literal 60 * time.Second…` | marker_wait_deadline_test.go:68 | stray → marker-wait-deadlines |
| reintroduced-bare-skip | `…calls t.Skip outside internal/capability` | skip_ownership_test.go:86 | stray → skip-ownership |
| unrouted-subcommand | `…dispatches "unregistered" with no entry…` | subcommand_routing_test.go:126 | stray → subcommand-routing |

### No seam will emit after the split

- **go-build-broken** — `go build failed` is `formatProbeFailure` framing
  (package_core_diagnostics_test.go:10), not compiler output; the serial build
  phase streams raw child output plus `phase build: red (exit N)`. Compounding:
  the build phase materializes only when both `scripts/go-build.sh` and
  `go.mod` exist (gate/phases.go:61–67) and the fixture tree ships only
  `go.mod` + `cmd/` — unchanged, it goes green after the collapse.
- **go-test-failing** — same framing; a standalone test phase emits Go's own
  `--- FAIL` lines, never the literal.
- **Routing gap:** `FixturePhase` (canary/canary.go:32–41) pins every
  conformance family to the `conformance` phase and `phasesForMode`
  (gate/phases.go:271–290) honours only `conformance`/`contract` — new phases
  are invisible to fixtures until both widen; an unknown family name falls
  through to owner `""` (runs every non-canary phase).
- **Uncanaried seams the split creates:** `gofmt: unformatted Go files`
  (:189), `go vet failed` (:208), race-test failures (:227, :229), `go list
  failed` (:212), `go build setup failed` (:197), `cross-compile failed`
  (cross_compile_stress_test.go:35) — no fixture today.
- **Coupled assertion:** `TestCoreSubprocessFailuresUseProbeFormatter`
  (package_core_diagnostics_test.go:99–128) text-scans the checks file for the
  exact label inventory; every label that leaves `checkGoCore` breaks it —
  mechanical build casualty, update in the migration.

## Part B — family→check mapping (verified by emitter, not name)

| family | check | fixtures | exceptions |
|---|---|---|---|
| package-core-guard | package-core-guard | 39 | 9 strays above (7 gain CHECK; 2 ship already carry it) |
| line-routing | line-routing | 5 | none |
| load-validity-metadata | load-validity-metadata | 10 | none (two skills-file emitters bound through checkLoadValidityMetadata) |
| skills-index-command-adapters | skills-index-command-adapters | 5 | none |
| data-handling-derivation | data-handling-derivation | 1 | none |
| docs-currency-token-diet | docs-currency-workflow | 14 | none |
| workflow-guidance-anchors | docs-currency-workflow | 41 | none — shares the check with two other families |
| coverage-map-validation | docs-currency-workflow | 2 | none — emitters live in internal/coverage, shelled through checkCoverageMaps |
| compliance-hardening | canary-inner-compliance | 2 | none — fires via checkCanaryInnerCompliance grading root under BENCH_CANARY_INNER=1 + marker file, NOT the same-named kit-compliance (which grades kitRoot) |
