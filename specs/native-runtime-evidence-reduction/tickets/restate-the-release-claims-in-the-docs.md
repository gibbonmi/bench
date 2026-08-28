# Restate the release claims in the docs

Blocked by: collapse-the-artifacts-job-to-one-generation.md
Writes: docs/release-runbook.md, docs/field-guide.html, projects/benchkit.md, ROADMAP.md

## What to build

The docs state what the release now proves, so a reader with no memory of this change reads no retired promise.

Two claims are retired, and the docs name both. Bench no longer proves that two independent checkouts finalize identical release evidence. Bench no longer proves the shipped macOS binaries through `native-proof`; the `smoke` job is the remaining macOS execution evidence.

The docs record the resulting state, not the change history. They carry no file paths and no snippets where an ADR rule forbids them.

## Acceptance

- [ ] The release runbook states the current native-proof scope and names the proven targets as plan data (row C4).
- [ ] The release runbook and the field guide state the surviving reproducibility claim, and neither promises the cross-checkout comparison (row C4).
- [ ] `projects/benchkit.md` and `ROADMAP.md` carry no retired claim (row C4).
