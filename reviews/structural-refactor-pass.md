# Review pickup: structural-refactor-pass

Frozen base `8eea2d15d661dd590e0f36958c27fa0c9f3f69cd`, reviewed tip
`cace3cafed2620f3595288e078db27d7bba7789b`, 29 commits, 60 files. Three axes ran
at `opus` / medium on 2026-09-05. Raw findings: 13. De-duplicated repair
targets: 7. One finding blocks. The reviewer accepted every `auto-fix` and
`ask-user` disposition below on 2026-09-05.

## Standards

Count: 4. Worst: S1, the growth mode re-derives the over-budget rule the all
scan owns.

| id | finding | citation | disposition | blocks |
|---|---|---|---|---|
| S1 | `Growth` restates four facts `evaluate` owns: the newline count of a regular file, the over-budget compare, the accept exemption, and the loud unreadable-accept line | `internal/structure/structure.go` :95-106 and :72 against :284, :305-313, :370-382; AGENTS.md "one source per fact" | auto-fix | no |
| S2 | the usage line spells `--since` and `--growth` as an alternation, and `Command` accepts both and drops `--growth` in silence | `internal/structure/structure.go:33` and :209-216 | ask-user, accepted: refuse the pair | no |
| S3 | a test comment narrates the change ("a caller used to state by hand") | `internal/landing/lane_test.go:15`; craft-comments register | auto-fix | no |
| S4 | the package doc line "the two-derivations bug class this slice ends" was reflowed and kept | `internal/structure/structure.go:7` | no-op, folds into S1 | no |

## Spec

Count: 5. Worst: P1, one over-budget file has no backlog row.

| id | finding | citation | disposition | blocks |
|---|---|---|---|---|
| P1 | SR58 is partial: `internal/guards/guards_test.go` (427 lines) is over budget at the base and absent from the backlog | `capture/restructure-backlog.md` partition 1 | auto-fix | no |
| P2 | ticket 8 names `TestLaneClassesNameOnlyDeclaredChecks` as a changed pin; decision (n) and the tree say it needed no edit | `specs/structural-refactor-pass/tickets/run-the-growth-check-in-the-fast-lane.md:22-24` | no-op, the ticket text follows the tree | no |
| P3 | three added tests sit outside the "Flagged additions" list: the kit-source symlink test, the recorded-worktree symlink test, and the nil-lane owner test | `internal/gate/kit_source_test.go:14`, `internal/intent/worktree_owner_test.go:47`, `internal/landing/lane_test.go:36` | no-op, the list gains them | no |
| P4 | SR14, SR56, and SR57 name a build-recorded artifact that does not exist; the axis ran both and both hold | spec rows SR14, SR56, SR57 | no-op, the record is this file | no |
| P5 | the resume's spec-file path now derives through the closed-folder path; equivalent for every spelling the grammar admits | `internal/worktree/land_resume.go:284-291` | no-op | no |

## Coverage

Count: 4. Worst: C1, the growth check never fires in the fast lane.

| id | finding | citation | disposition | blocks |
|---|---|---|---|---|
| C1 | the lane runs the check in a detached checkout whose HEAD equals the base, and the growth query diffs `base..HEAD`, so the change list is always empty; a planted growth passed `bench commit --dry-run` with `check=structure` green | `internal/structure/structure.go` `growthChanges`; `internal/gate/prospectiveartifact/prospectiveartifact.go:103-106`; `internal/commit/commit.go:112,153` | ask-user, accepted: compare the base against the working tree, add a detached-checkout test and a live dry-run dogfood | yes |
| C2 | `--growth ""` resolves to `HEAD..HEAD`, passes, and prints a dangling `since ` | `internal/structure/structure.go` `growthChanges` | ask-user, accepted: refuse an empty value | no |
| C3 | the git-object reader answers "directory" for a committed blob at the folder path; the working-tree reader answers false | `internal/landing/close.go` `commitTreeReader.FolderIsDirectory` | auto-fix: require a tree object, add the fifth SR31 case | no |
| C4 | the exported match absolutizes an empty `Worktree` row to the cwd where the old compare did not; unreachable from the one caller | `internal/worktree/landed.go:37-41` | no-op | no |

## Repair targets

| target | findings | files |
|---|---|---|
| R1 growth query against the working tree, empty base refused, detached-checkout test, dry-run dogfood | C1, C2 | `internal/structure/structure.go`, `internal/structure/growth_test.go` |
| R2 one over-budget predicate shared by `evaluate` and `Growth`, package doc corrected | S1, S4 | `internal/structure/structure.go` |
| R3 `--since` with `--growth` refuses at exit 2 | S2 | `internal/structure/structure.go`, `internal/structure/growth_test.go` |
| R4 timeless comment | S3 | `internal/landing/lane_test.go` |
| R5 tree-object check in the git reader, fifth SR31 case | C3 | `internal/landing/close.go`, `internal/landing/close_test.go` |
| R6 backlog row for the guards test file | P1 | `capture/restructure-backlog.md` |
| R7 spec and ticket text: rows SR62 and SR63, decisions (t), (v), (w), the flagged-additions list, ticket 8's pin sentence | P2, P3, C1 | `specs/structural-refactor-pass/` |
