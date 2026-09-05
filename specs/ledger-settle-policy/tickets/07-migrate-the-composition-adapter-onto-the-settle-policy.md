# Migrate the composition adapter onto the settle policy

Blocked by: 06-add-the-settle-policy-child.md
Writes: internal/landing/composition.go, internal/landing/composition_test.go, internal/landing/landing_helpers_test.go
Covers: LS22, LS24, LS25, LS27, LS39

## What to build

Verify the premise first. Read `Compose`, `resolveCaptureConflict`,
`unionStages`, `parseConflict`, `contentConflictKind`, `CaptureSide`,
`stageRecord`, `unionStage`, `Conflict`, `CompositionResult`, `mergeTree`,
`mergeTreeResult`, `editTree`, and `indexRun` in
internal/landing/composition.go.

Read `TestComposeClassifiesRealGitConflictsWithoutMutation`,
`TestConflictKindRejectsEmptyMergeTreeOutput`,
`TestComposeSettlesPhaseOwnedConflictsByRule`, and
`TestComposeRefusesConflictsTheRuleTableCannotSettle` in
internal/landing/composition_test.go. Read the settle policy ticket
06-add-the-settle-policy-child.md adds.

Migrate the adapter onto the policy. `parseConflict` stays in the adapter and
turns the `merge-tree --write-tree -z` output into the policy's stage-record
type. The adapter calls the policy, and it applies the verdict with
`update-index`. The union verdict's text merge stays in the adapter, because it
shells to `git merge-file`. Git's output format keeps one reader.

`ConflictError.Error` does not change in this ticket. Ticket
08-surface-the-settle-refusal-reason.md adds the reason field and the new text,
so leave today's `composition conflict: <kind>` whole here. This ticket writes
no new file inside `internal/landing/` itself, because that directory already
holds 16 source files against a 12-file budget.

Shrink `internal/landing/composition_test.go` under the 400-line budget. The
file holds 585 lines today, and the settle partitions that move to the policy
tables are the reduction this ticket owns. Report to the coordinator when the
move leaves the file over 400 lines. LS39 is the row this shrink answers.

Keep the seven real-Git classification journeys. Keep one settle journey for
each verb and one refusal journey. Move every other settle partition to the
policy tables ticket 06-add-the-settle-policy-child.md adds. Delete no
partition; move it.

Add one malformed-record case to
`TestConflictKindRejectsEmptyMergeTreeOutput`. The case drives `parseConflict`
with a record whose header holds the wrong field count.

## Acceptance

- [ ] `parseConflict` turns a `merge-tree -z` output into the policy's stage records.
- [ ] `parseConflict` returns an error for empty output and for a malformed record.
- [ ] Over real Git conflicts, `Compose` answers the kinds `textual`, `modify/delete`, `rename/rename`, `file/directory`, `mode`, `symlink`, and `gitlink` without mutating the repository.
- [ ] The composition keeps one settle journey for each verb and one refusal journey.
- [ ] Every other settle partition is a policy table case.
- [ ] `ConflictError.Error` still reads `composition conflict: <kind>`.
- [ ] `bench structure` reports `internal/landing/composition_test.go` under the 400-line budget.
- [ ] `bench structure` reports no crowded-directory growth for `internal/landing/`.
- [ ] The pre-existing `internal/landing` suite passes with its test logic unchanged, except the moved settle partitions.
- [ ] Self-probe: move the parser into the policy, and report the settle policy census red.
