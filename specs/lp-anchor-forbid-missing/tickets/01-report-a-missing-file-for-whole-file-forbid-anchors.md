# Report a missing file for whole-file forbid anchors

Blocked by: none
Writes: internal/anchors/registry.go, internal/anchors/anchor_harness_diagnostics_test.go, cmd/bench/anchor_help_test.go
Covers: none

## What to build

Source: drain d-359accf567c1, learning L62.

The anchor registry evaluator reports `acceptance coverage anchor file missing: <file>` for a `Require` row whose file is absent.
A `Forbid` row or a `ForbidCaseFoldedEmphasis` row on an absent file reads the absent file as clean text, so the row is green.
Thus a rename of a guarded file silently disables each forbid on that file.

The reviewer decision: a whole-file forbid row on an absent file reports the same missing-file diagnostic as a `Require` row.
The evaluator uses one whole-file branch for every kind that is not section-scoped.
That branch reports the existing diagnostic once for each row, as the `Require` branch does now.
Section-scoped kinds keep their own `section-scoped anchor file missing: <file>` diagnostic.

The `bench anchors AGENTS.md` absent-file expectation counts one missing-file diagnostic for each `AGENTS.md` row, because that file has a whole-file `Forbid` row.

## Acceptance

- [ ] On an absent file, each `Require`, `Forbid`, and `ForbidCaseFoldedEmphasis` row reports `acceptance coverage anchor file missing: <file>` once, and no row reports its needle diagnostic.
- [ ] On a present file, a whole-file forbid row reports its needle diagnostic only when the file carries the needle.
- [ ] `bench anchors AGENTS.md` on an absent file prints one missing-file error for each `AGENTS.md` row.
- [ ] The root conformance check stays green on the kit tree, so no registry row names an absent file.
