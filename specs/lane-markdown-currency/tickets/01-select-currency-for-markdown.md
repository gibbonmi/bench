# Select the docs-currency check for every Markdown change

Blocked by: none
Writes: internal/gate/lane_select.go, internal/gate/lane_select_test.go, internal/gate/lane_test.go, projects/benchkit.md
Covers: none

## What to build

The commit lane selects `docs-currency-workflow` only for a registered anchor file. The check reads Markdown across the tree: the root documents, the command files, and every Markdown file under `specs/`, `decisions/`, and `.agents/`. A Markdown change outside the anchor files therefore passes the lane with `prose` alone. The whole-project gate then fails it at the landing. Astra's planning branch committed eleven such reds on 2026-10-07.

Make the `markdown` lane class select `docs-currency-workflow` beside `prose`. Every input of the check is Markdown, so this rule needs no copy of the check's file list. The anchor class keeps its selection for an anchor file that is not Markdown. The profile's lane table renders the new selection.

## Acceptance

- [ ] A commit of one Markdown file runs `prose` and `docs-currency-workflow` in the lane.
- [ ] A Markdown file that cites a stale `$bench-` command fails the lane with the check's own diagnostic, so `bench commit` refuses it.
- [ ] A Go-only commit does not run `docs-currency-workflow`.
- [ ] The profile's lane table names the selecting classes of `docs-currency-workflow`, and the `profile-lane-table` check passes.
- [ ] Dogfood: with the candidate binary, the collision 7 repro reports that the lane refused the change.
