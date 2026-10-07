# Relocate the existing audit through its real check

Blocked by: none
Writes: internal/conformance/injected_ports_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/injectedports/ (new), internal/conformance/checks_test.go, tests/canary/injected-ports/
Covers: PS01, PS02, PS03, PS12, PS13, PS14, PS41, PS42, PS43, PS50, PS51, PS52, PS53, PS54, PS55, PS61

## What to build

Move the existing derivation into `internal/conformance/injectedports/` and make the actual `checkInjectedPortRegistry` wrapper call `Check(root, Policy)` once.
Preserve the seven named-port rows, five required-package expectations, and exact legacy diagnostic projections.
Keep existing assertions in `TestInjectedPortRegistryCheckBites` and `TestInjectedPortDerivationSeesEveryPortShape` intact.
Derive named shapes across eligible production packages, with the legacy partial-root contract retained.
Do not enable variable, setter, runtime, or environment enforcement in this ticket.

Carry the production adapter, owner tests, real check binding, and existing injected-ports canary in the same green checkpoint.
Remove each moved walker from the old file immediately; TS4's final whole-tree census cannot authorize a duplicate owner here.
Keep the shared conformance fixture builder at its existing owner, and give direct owner fixtures one local constructor.
Read declaration and injection-site owners, their actual consumers, and the legacy registry before moving symbols.
No unused public boundary or second AST scanner may remain.

## Acceptance

- [ ] Nonempty interfaces, named function types, and all-function structs retain their unregistered findings, including pointer parameters, assertions, and switches.
- [ ] Empty interfaces, data-bearing structs, and unused port shapes stay negative.
- [ ] Parse errors, present required zero inventory, orphan rows, and blank external exemptions retain their existing refusal behavior.
- [ ] Absent fixture packages stay benign, unreadable present source refuses, and hostile root paths never read ambient checkout evidence.
- [ ] `checkInjectedPortRegistry` actually reaches the relocated owner; bypassing that call makes the existing canary assertion fail with its exact expected diagnostic.
- [ ] Existing named-port assertions and external publication limitations remain unchanged; new fault-evidence contracts do not rewrite legacy expectations.
- [ ] The old file loses relocated derivation and fits 400 lines. The new directory stays within 12 files and 400 lines per file.
- [ ] The focused owner and actual conformance checks pass with every successor still unbuilt, and independent TS1 review closes before TS2 starts.

Integration checks: the TS1 completion-plan commands exercise owner fixtures, the existing registered check, and its canary.
Record moved definitions and callers so the real wrapper is distinguishable from a helper-only proof.
