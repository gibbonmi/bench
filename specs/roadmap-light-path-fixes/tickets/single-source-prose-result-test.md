# Single-source the prose result test

Blocked by: structure-gate-prose-findings.md
Writes: internal/prose/prose_test.go
Covers: none

## What to build

Keep the structured prose-result assertion without independently reproducing
the full rendered diagnostic protocol already covered at its public gate seam.

## Acceptance

- [ ] The prose test independently checks structured path, line, rule, count,
      and sentence fields.
- [ ] The prose test does not reconstruct the rendered diagnostic string.
- [ ] Existing gate-prose public-diagnostic coverage remains green.
