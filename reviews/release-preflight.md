# Release preflight implementation review

## Standards

2 findings. Worst: major — split executable phase registry.

- Major: `internal/preflight/registry.json` owns ordered phase names, while `internal/preflight/command.go` repeats phase dispatch, required inputs/tools, and executable mapping. Collapse production policy into one registry-backed owner per the repository's one-source-per-fact rule.
- Minor: `.github/workflows/native-runtime.yml` and `.github/workflows/release.yml` independently pin `govulncheck@v1.6.0`. Move the shared tool version/setup fact behind one repository-owned source consumed by both workflows.

## Spec

4 findings. Worst: critical — fresh workflow checkouts cannot execute preflight.

- Critical: both workflows call `scripts/release-preflight.sh` before building, installing, or downloading the compiled Bench binary; `bin/bench.sh` therefore exits 127 on a clean runner. Full and focused matrix calls must bootstrap the canonical compiled command.
- High: `internal/preflight/evidence.go` renames the prior verdict away before promoting staging, so SIGKILL between renames leaves no authoritative `dist/preflight` path. Preserve an old-or-new verdict at the canonical path across the specified fault model.
- Medium: command initialization failures return before writing complete terminal records and a manifest, contrary to the every-handled-terminal-state evidence contract.
- Medium: runtime policy tests call private helpers and omit built-command failures for tag/version mismatch, stranded changelog content, unrelated or shallow ancestry, malformed/expired exceptions, and partial-promotion faults.

## Coverage

4 findings. Worst: critical — native workflow bootstrap is not exercised.

- Critical: add a clean-checkout workflow/bootstrap proof that executes the real full and focused command paths, not only YAML string assertions.
- High: add fault injection for SIGKILL at evidence-promotion boundaries and prove the prior complete verdict remains authoritative.
- High: vulnerability scanning uses bare `exec.CommandContext`; prove and implement child-process-group cancellation for a scanner that spawns descendants.
- Medium: regular-file and FIFO evidence targets are currently renamed and overwritten. Reject non-directory targets before blocking or replacement, with hostile-target tests.
