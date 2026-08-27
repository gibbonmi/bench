# Repair selection and environment closure

Blocked by: 03-run-one-named-conformance-check.md
Writes: internal/testreport/selection.go, internal/testreport/selection_test.go, internal/testreport/command.go, internal/testreport/testreport_test.go

## What to build

Return an explicit empty report for a deleted non-Go-only subject. Keep the refusal for absent Go-relevant inputs.

Remove all ambient conformance controls from every focused mode. Keep only the exact controls for a named check.

Keep one source for the injected changed-package graph. Use the fixture module only for filesystem facts.

## Acceptance

- [ ] C06 — a deleted non-Go-only subject returns the three zero-row tables and exit 0.
- [ ] C07 — a deleted Go-relevant input keeps its current mapping or refusal posture.
- [ ] K03 — every focused mode removes all ambient conformance controls.
- [ ] K03 — a named check installs only its exact conformance controls.
- [ ] C09 — the changed-selection fixture has one source for package-graph facts.
