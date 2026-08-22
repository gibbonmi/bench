# Migrate the board under the Next grammar

Blocked by: 02-parse-the-next-token.md, 04-write-the-drain-rules-and-bind-the-token-table.md
Writes: ROADMAP.md, roadmap/, internal/roadmap/tree_validation.go, internal/conformance/roadmap_detail_integrity_test.go, tests/canary/roadmap-detail-integrity/, CHANGELOG.md, capture/session-handoff.md

## What to build

One reviewed `--restructure` pass proposes a `Next:` token for every existing
row as one batch diff, and moves a row with no honest next action to the
parked section. The same commit turns the missing-line class into a gate fault
in `roadmap-detail-integrity` and plants its canary fixture, so the tree is
never red between the check and the board. The reviewer approves the batch
diff before the commit. Spec group E, rows RF14 and RF29.

## Acceptance

- [ ] After the commit, `roadmap-detail-integrity` returns no `Next:` diagnostic over the live tree.
- [ ] The canary fixture with a deleted `Next:` line makes the owner check red, and restoring the line clears it.
- [ ] The batch diff changes no row's priority order and no body prose beyond the added marker or the move to the parked section.
- [ ] The reviewer's approval of the batch diff is recorded in the commit message.
