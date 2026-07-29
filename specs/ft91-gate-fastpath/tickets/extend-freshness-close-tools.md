# Extend freshness to 60 minutes and declare the toolchain tools

Blocked by: none

## What to build

Stories 12 and 13 of `specs/ft91-gate-fastpath/spec.md`: `freshness` in
`internal/gate/verdict.go` moves 10 → 60 minutes in the one shared constant —
no second constant, no `policyVersion` bump — replacing the existing 10-minute
boundary pin in `gate_test.go` in the same change. `.bench/gate-inputs.json`
gains `"go", "node", "npm"` in `tools`; `shellcheck` stays undeclared.

## Acceptance

- [ ] Freshness is 60 minutes, pinned at both edges: fake-`Now()` inspect-path
      cases at 59 minutes (reusable) and 61 minutes (not) replace the
      10-minute pin.
- [ ] Exactly one freshness constant exists; `policyVersion` is unbumped.
- [ ] `.bench/gate-inputs.json` declares `go`, `node`, `npm` in `tools` and
      not `shellcheck`; the `R1/manifest-missing-tool` contract still proves a
      declared-but-absent tool opens the subject, exercised against a
      declared-tool name rather than a fixture-only manifest.
