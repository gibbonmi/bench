# Audit — injected-interface composition (2026-08-03)

Decision source for the junction-test spec. Sweep charge: every injected
interface on a Service-like type repo-wide; for each, the production wiring,
whether any test drives the real implementation through the consuming surface,
and any written exemption. Delegate findings spot-verified by the coordinator
at the four load-bearing citations.

## Inventory

| Interface (pkg) | Production wiring | Real-producer test through the consuming surface | Verdict |
|---|---|---|---|
| `GateOwner` (specbuild) | `productionGateOwner{}` — `cmd/bench/specbuild.go:30` | Binary-level only: `internal/contract/runtime/runtime_spec_build_refusal_test.go` `TestStartMarkerAncestryEndToEnd`; no Go-level Service composition | priced-out — contract phase covers it |
| `PromotionGateOwner` (specbuild) | same wiring; acquired via type assertion at `internal/specbuild/assign.go:300` | Binary-level: `runtime_spec_build_test.go` `TestRuntimeSpecBuildRedRepairGreenRetainsComposedEvidence` | priced-out, plus one compile-time pin (see recommendations) |
| `WorktreeOwner.Create` (specbuild) | `productionWorktreeOwner{}` — `cmd/bench/specbuild.go:30` | YES — `realOwner` delegates to real `worktree.Create` (`internal/specbuild/lifecycle_test.go:321`) | covered |
| `ReleaseOwner` (specbuild) | `cmd/bench/specbuild.go:188-192` → `worktree.ReleaseProvisional` | **NONE** — `realOwner.Release` is a bare `return nil` (`lifecycle_test.go:330`); real refusal surface pinned only in `internal/worktree/orphan_test.go` | **P1 — junction test** |
| `AbandonOwner` (specbuild) | `cmd/bench/specbuild.go:194-201` → real planner | Mixed — real planner for the plan/apply family (`abandon_test.go:774`), **fake `decayedOwner` for the decayed/husk/unreadable family** (`abandon_test.go:805-817`) | **P1 — junction test** (FT181's exact residual surface) |
| `Runner` (specbuild) | `processRunner{}` default (`lifecycle.go:89`) | YES — real in nearly every test; `countingRunner` is a pass-through decorator | covered |
| `gateEngine` / `gateFile` / `FS` (internal/gate) | `productionGateEngine{}`, `RealFS()` in-package | YES — real engine driven in `check_slots_test.go`, `verdict_reuse_test.go`, `component_inputs_test.go` | covered |
| `Checker` (internal/gitguard) | `cmd/bench/main.go:363` wiring `git.RefResolves` / `git.BranchExists` | **NONE** — all `Classify` tests use constant `refYes`/`refNo` fakes; the two real probes have opposite fail-safe polarities (`internal/git/git.go:29-46`) never composed | **P1 — junction test** |
| `Runner` (internal/canary) | `defaultRunner` (`canary.go:259`, `:278`; prep-release via `internal/preprelease/preprelease.go:121`) | **NONE** — `resolvingRunner` (`canary_path_test.go:102`) reimplements `defaultRunner`'s resolution semantics by its own comment | **P1 — junction test** |
| `Registry` (internal/publication) | `NewFixtureRegistry` — `command.go:173` | YES for `FixtureRegistry` (contract surface suite); `NPMCLIRegistry` exempt | covered / exempt — written exemption at `internal/publication/registry.go:5-11` (no egress, no credential) |
| `buildService` (cmd/bench) | real `*specbuild.Service` at `specbuild.go:30` | Binary-level runtime contract suite; `executeBuild` is pure dispatch/render | priced-out — exemption line |
| `ShapeUnknown` Lstat branches (worktree classifier) | `classifier.go:138/:145/:155`; consumed at `internal/specbuild/precondition.go:252` (error discarded) | **NONE** — `ShapeUnknown` appears in zero test files | **P1 — fixture exists** (see below) |

Service-like types confirmed to carry no injected port (sweep completeness):
internal/worktree (free functions — the producer side, no seam by design),
internal/shift, internal/adopt, internal/preflight, internal/preprelease
(static step table), internal/contract harness types, and the data/pure
packages. `FixtureRegistry.Client *http.Client` is concrete with a default,
never overridden.

## Triage

**P1-class — production composition can diverge silently; junction test now:**

1. `ReleaseOwner`: structurally identical to FT181's abandon gap. Every
   `Integrate` test asserts against a release that cannot fail while
   `worktree.ReleaseProvisional` refuses drifted/unretained checkpoints.
2. `AbandonOwner` decayed family: the husk/decayed/unreadable shapes — the
   states FT181 found broken — still graded against a synthetic fingerprint.
   The in-code note ("internal/worktree's contract") is an existing-control
   claim that names no control reaching this surface; under craft-spec's
   existing-control rule it becomes a row.
3. `canary.Runner`: the fake is a second derivation of the production
   runner's path resolution — exactly what craft-gate's run-the-real-path
   rule forbids in a check; a cwd/arg change in `defaultRunner` stays green.
4. `gitguard.Checker`: two real probes with opposite fail-safe polarities,
   composed with `Classify`'s clobber logic by nobody.
5. `ShapeUnknown`: privilege-free deterministic fixture exists — `os.Lstat`
   of `<regular-file>/child` fails `ENOTDIR`, which is not `ErrNotExist`, so
   `classifier.go:138` returns `ShapeUnknown` with no chmod or root needed.
   In sweep scope per the charge. The discarded error at
   `precondition.go:252` is noted; behavior change is out of scope.

**Priced-out — leave with reason (registry exemption line):**

- `GateOwner` / `PromotionGateOwner` / `buildService`: composition covered by
  the binary-level runtime contract suite, which drives the real wiring end
  to end — the fake is below the seam the contract phase tests. One cheap
  addition: a compile-time `var _ specbuild.PromotionGateOwner =
  productionGateOwner{}` pin in cmd/bench so the `assign.go:300` type
  assertion can never silently downgrade.
- `NPMCLIRegistry`: written exemption already exists and is sound.
- `WorktreeOwner.Create`, specbuild `Runner`, internal/gate ports: covered.

**Gate check (parked idea, recommended in scope):** a conformance check that
every injected interface names its real-producer test or its exemption —
registry single-sourced, existence-verified, canary-proven per craft-gate —
so this audit never re-runs by hand.
