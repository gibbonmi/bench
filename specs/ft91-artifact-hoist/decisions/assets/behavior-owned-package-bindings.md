# Behavior-owned fixture → contract package bindings (2026-07-28)

Build input for `specs/ft91-canary-contract-scoping.md` story 7. Method: each
fixture's EXPECT string traced with a fixed-string search to the test file
that emits it; the package column is the file's package directory relative to
`internal/contract`. The build verifies every binding the cheap way — a wrong
one reds as did-not-bite at the first gate.

| fixture | package |
|---|---|
| coverage-extraction-dropped | axi |
| diff-recorded-base-dropped | axi |
| guards-aggregation-dropped | axi |
| learnings-parse-broken | axi |
| roadmap-context-incomplete | axi |
| session-start-resume-cleanup-dropped | axi |
| toon-escaping-dropped | axi |
| gate-verdict-invalidation-bypassed | runtime |
| gate-verdict-oracle-binding-bypassed | runtime |
| intent-common-dir-address-regressed | runtime |
| status-landed-aggregation-regressed | runtime |
| status-regressed | runtime |
| worktree-lifecycle-safety-bypassed | runtime |
| doctor-foreign-clobbered | surface |
| doctor-manager-dir-chosen | surface |
| doctor-stale-silent | surface |
| native-selection-regressed | surface |
| postinstall-guard-bypassed | surface |
| postinstall-nonzero-exit | surface |
| repo-local-forwarding-dropped | surface |
| session-start-advice-dropped | surface |
| unscaffolded-bench-file | surface |
| wrapper-args-dropped | surface |
| native-proof-aggregation-bypassed | surface/artifact |
| native-proof-digest-binding-bypassed | surface/artifact (the fixture injects a build-tagged test into the graded root; `TestAuthoritativeNativeProofBehaviorCanary` drives it as a subprocess and reports its output, so the identically worded string in the release-evidence probe is not the emitter under the gate) |
| offline-archive-digest-binding-omitted | surface/artifact |
| wrapper-contamination-admitted | surface/artifact |
| wrapper-required-surface-dropped | surface/artifact |
| integrity-mismatch-acceptance | surface/publication |
| premature-wrapper-promotion | surface/publication |
| publication-order-bypass | surface/publication |
| publication-unpublish-attempt | surface/publication (emitter is `scripts/offline-registry.mjs`, driven by the publication suite — verify via bite) |
| roadmap-regressed | runtime or axi (EXPECT fragment emitted by both suites — bind to whichever bites; verify via did-not-bite) |

## Not contract-owned — relocate to legacy flat fixtures

One EXPECT is emitted by no phase at all, so that fixture is mis-familied
today and can take no package binding:

| fixture | actual emitter |
|---|---|
| phase-manifest-defect-admitted | the gate's manifest loader (`internal/gate`), which reds before any phase runs |

It moves to `tests/canary/<fixture>/` — the existing legacy flat class: full
inner gate, shared unscoped baseline, no new mechanism. Cost accepted: one
full inner run (~1 worker-minute) against ~1,600 worker-seconds removed from
the other 33.
