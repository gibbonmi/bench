# Name a Writes entry that landed work deleted in its own preflight row

Blocked by: none
Writes: internal/preflight/decision.go, internal/preflight/fence_writes.go, internal/preflight/gather.go, internal/preflight/gather_inputs.go, internal/preflight/default_branch.go, internal/preflight/tree_probe.go, internal/preflight/writes_landed.go (new), internal/preflight/writes_landed_test.go (new), internal/preflight/decision_test.go, internal/preflight/source_tip_test.go, internal/preflight/charge_cells_test.go, internal/preflight/preflighttest, internal/git/default_branch.go, internal/spec/widen.go, internal/preflight/staleness/drift_test.go, internal/preflight/default_branch_test.go (new), internal/spec/widen_signedoff_test.go
Covers: ES24, ES25, ES26, ES27, ES28, ES29, ES30, ES31, ES32, ES33, ES34, ES35, ES36, ES37, ES38, ES39, ES40, ES45

## What to build

A bare `bench preflight build <slug>` at a fresh build entry tells a `Writes:` entry that landed work deleted from a typing error.
Read these spec sections before you start: "The `writes-landed` fact", "The `writes-resolve` and `writes-landed` precedence", "Row order and applicability", the acceptance coverage map, and the edge inventory.

Make these changes:

- Add the build-mode fact `LandedDeletions` to `Facts`. It maps a `Writes:` entry to the full identity of its deleting commit.
- Fill the fact in build mode only. Consider a plain entry that `tickets.ParseWrite` accepts and that the checkout does not hold. Do not read `DeletedPaths` or `BuildDeletedPaths` in the gatherer, because the check alone owns the precedence.
- Resolve the default branch through `git.ResolvedDefault` and its tip through `git.ResolveCommit`. Run the tip lookup and the deleting-commit query exactly as the spec states them, with the `:(literal)` pathspec. A failed query, an unresolved default branch, or an empty answer records no entry.
- Put the query code in a new file, `internal/preflight/writes_landed.go`, or in `default_branch.go` or `tree_probe.go`. `gather.go` has 7 lines of headroom under its 400-line budget, so it gets only the fact assignment.
- Add the row `writes-landed` through `ticketRow` directly after `writes-resolve`, in build mode only. Give it no `exempt` flag. Update the `Decide` doc comment.
- Change `writesResolveCheck` and add the `writes-landed` check so that both walk the entries in the five-step precedence of the spec. Step 3 wins over step 4. Keep the `writes-resolve` detail text byte for byte, and use the `writes-landed` detail text of the spec.
- Update the row-list and row-count readers: `bareRowCount` in build mode becomes 13, `TestDecideChargeRows` gets `writes-landed` after `writes-resolve`, `TestTicketGatedRowsNotApplicableWithoutTickets` gets `writes-landed` after `writes-resolve`, and the not-applicable count of `TestDecideAllGreen` at `decision_test.go:100` becomes 10.
- `decision_test.go` holds 557 lines over its 400-line budget. Make its edits in place, so its net growth is zero lines. Do not put a repository fixture in it, because `internal/conformance/ordinary_build_census_test.go:158` forbids one.
- Put every new command test in `writes_landed_test.go`. Make each fixture history with `preflighttest.RunGit`. If a shared history helper is needed, put it in a new file under `internal/preflight/preflighttest`, because `fixture.go` has 6 lines of headroom. `FenceWrites` marks every entry `(new)`, so write a plain-entry ticket body.
- Take each expected commit identity from `git rev-parse` of the commit that the fixture made. Write each expected detail as a literal.

Four existing test files keep their answers, so they are not on the `Writes:` line:

- `fence_writes_test.go` seeds through `seedDeletionReview` (`internal/preflight/fence_writes_test.go:34`), which deletes on a feature branch (`internal/preflight/deletion_preflight_test.go:88-89`), so the default-branch tip still holds the path.
- `command_build_test.go` names a never-tracked path (`internal/preflight/command_build_test.go:239`).
- `deletion_preflight_test.go` and `command_review_test.go` run review mode, and the fact is filled in build mode only.

The ES-C1 checkpoint command runs `TestWritesResolveEarlierChunkDeletion`, the committed-deletion tests, and `TestCommandBuildRendersFourGrammarRows`. If one goes red, stop and report it.

## Acceptance

- [ ] A bare build preflight with no recorded chunk reds `writes-landed` and keeps `writes-resolve` green for a plain entry that a default-branch commit deleted, with the full deleting commit in the detail (ES24, ES25, ES37).
- [ ] `writes-landed` reds a merge-resolution deletion, a side-branch deletion before its merge, and a rename (ES26, ES27, ES28).
- [ ] `writes-resolve` keeps a never-tracked entry, a restored entry, and a never-tracked glob entry beside a deleted sibling (ES29, ES31, ES35).
- [ ] A mixed set reds both rows, an earlier-chunk deletion and a folded `--base` deletion are green in both rows, and an unresolved default branch keeps the entry in `writes-resolve` (ES30, ES32, ES33, ES34).
- [ ] `writes-resolve` keeps its detail text, review mode prints no `writes-landed` row, and the row is not applicable without a tickets folder (ES36, ES38, ES39, ES40).
- [ ] `--audit-spec` reports `manual_required` true while `writes-landed` is red (ES45).
- [ ] The five row-list and row-count tests pass with the new row (`TestDecideAllGreen`, `TestDecideChargeRows`, `TestSourceTipOmittedKeepsTodaysVerdict`, `TestSourceTipAcceptedByBothModes`, `TestTicketGatedRowsNotApplicableWithoutTickets`).
