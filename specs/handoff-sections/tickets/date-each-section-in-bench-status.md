# Date each section in bench status

Blocked by: add-the-handoff-document-leaf-package.md
Writes: internal/status/handoff.go, internal/status/handoff_test.go, internal/status/status.go, tests/canary/docs-currency-token-diet/signal-vocabulary-drift
Covers: HS21, HS22, HS23

## What to build

Verify the premise first: `appendIgnoredHandoff` in internal/status/handoff.go
dates the whole file by its write time. Then parse the sections through the
leaf package. For each request section, count `git rev-list --count
<section tip>..<assignment branch tip>` and name the section with the largest
count in the row. Keep the file-age rule for `main`. Expand per-section rows
only under `--all`, in the `expandCensusSignals` shape.

## Acceptance

- [ ] With two sections, one three commits behind on its branch, the row names that section with `3 commits behind`.
- [ ] Rewriting the fresh section leaves the behind row's count unchanged.
- [ ] The existing ignored-handoff rows still date `main` by the file write time.
- [ ] Self-probe: keep the file-level clock and label it with the first section, and report the two-section test red.
