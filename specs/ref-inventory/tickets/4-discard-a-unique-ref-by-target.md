# 4. Discard one unique unrecorded ref by target with a discarded ref first

Blocked by: 1-classify-unclaimed-refs.md, 3-sweep-discarded-refs.md
Writes: internal/worktree/clean_discard.go (new), internal/worktree/clean_discard_test.go (new), internal/worktree/clean_discard_transaction_test.go (new), internal/worktree/clean_classes.go, internal/worktree/clean_set.go, internal/worktree/clean_set_apply.go, internal/worktree/clean_set_command_test.go, internal/worktree/clean_unclaimed.go, internal/worktree/clean_unclaimed_test.go, internal/worktree/joins.go, internal/worktree/path.go, internal/worktree/worktree.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RI16, RI29, RI30, RI31, RI32, RI33, RI34, RI35, RI36, RI37, RI38, RI39, RI40, RI41, RI58, RI64, RI65, RI67, RI68, RI74, RI75, RI76, RI78, RI96, RI98, RI99, RI101

## What to build

Chunk: RI-C2b.

Route an operand to the unrecorded-branch fallback in `planExplicitSet` before the relative-path check.
The fallback takes an operand that starts with `refs/heads/bench/assign/`, `bench/assign/`, `refs/heads/bench/shift-`, or `bench/shift-`.
It also takes a 32-character hexadecimal operand that resolves no record, and matches it to the last segment of exactly one unclaimed assignment branch.
Two matches refuse with an error row that names both refs, and no match keeps the unassigned refusal.
Every other operand keeps today's path check, so `archive/x` still refuses with `relative path targets are unsupported`.

The fallback's candidate set is the unclaimed set that the bulk planner selects, not every unrecorded branch.
A prefix operand that names a checked-out branch refuses with an error row that names the checkout path.
A prefix operand that names the default branch refuses with an error row.
A prefix operand that names a recorded branch plans through the record and prints no `class=` prefix.
The set collapses `--target <id>` and `--target bench/assign/<owner>/<id>` into one row by branch ref.
Recorded rows sort by assignment id first, then unrecorded rows sort by branch ref, and the explicit set fingerprint version becomes `bench-explicit-set/v2`.

An unrecorded target row classifies through the class function of ticket 1.
Without `--discard-branch` the row carries `retain`, its detail names `--discard-branch`, and the apply removes nothing.
A unique row plans the discarded ref of ticket 3 in its recovery cell, dated by the clock join added to the package's join set.
Its detail starts with `class=unique`.
A landed or subsumed row plans recovery `none`.

The explicit set fingerprint binds each unrecorded row's ref, tip, class, holder, and planned discarded ref.
The apply needs the set fingerprint, and `--apply-current` stays refused outside the unclaimed mode.
Change the unique row's detail suffix in the unclaimed plan to `bench worktree clean --discard-branch --target <branch path>`, the row's ref without `refs/heads/`.
A unique shift row ends with `bench worktree clean --discard-branch --target bench/shift-<stamp>`, and that command plans one row for it.
The class file owns that suffix in its detail rendering, and the retained text of ticket 1 goes away with its constant.
The printed route and the apply action derive their selector from one function.

The apply of a unique row runs in this order inside the set apply.
When the planned ref exists at the row's tip as a direct ref, the write is skipped.
When the planned path is a symbolic ref, the apply prints an error row, keeps both refs, and writes nothing.
When it exists at another tip, the apply refuses and keeps the branch.
Otherwise the write uses the zero old value.

The delete uses the exact tip.
A branch moved after the write survives with an error row, and the discarded ref stays at the old tip.
One fault boundary step precedes the read of the planned path, one sits between that read and the write, and one follows the write.
A target that the class function faults, such as a resolving symref, prints an error row with no fingerprint, and both apply forms refuse.
The outcome row's recovery cell names the discarded ref.

## Acceptance

- [ ] `--target` with the full branch path, the path without `refs/heads/`, a shift branch path, or the id segment plans one row for the unrecorded branch.
- [ ] `--target <id>` and `--target bench/assign/<owner>/<id>` in one call plan one row.
- [ ] An id segment shared by two unrecorded branches prints an error row that names both refs and no fingerprint.
- [ ] `--target archive/x` prints the refusal `relative path targets are unsupported`.
- [ ] A prefix operand for a checked-out shift branch prints an error row that names the checkout path.
- [ ] A prefix operand for the default branch prints an error row.
- [ ] A prefix operand for a recorded active branch plans through the record with no `class=` prefix.
- [ ] `--target <unrecorded id>` without `--discard-branch` plans `retain` with `--discard-branch` in its detail, and the apply leaves the branch.
- [ ] The plan row of a unique target shows `class=unique` and the recovery `refs/bench/discarded/<yyyymmdd>/bench/assign/<owner>/<id>` under the fixed clock.
- [ ] After the apply, the discarded ref resolves to the old tip, the branch is gone, and the outcome row's recovery cell equals the discarded ref.
- [ ] A fault before the write leaves the branch and no ref under `refs/bench/discarded/`.
- [ ] A fault after the write leaves both refs, and a re-plan plus second apply removes the branch.
- [ ] A branch moved at the after-write step survives, the discarded ref stays at the old tip, and the row reads action `error`.
- [ ] A planted ref at the planned path at another commit makes the apply refuse and keep the branch.
- [ ] An explicit plan, then a recorded active branch at a descendant of the target, then the old fingerprint refuses as stale.
- [ ] After that stale refusal, `for-each-ref refs/bench/discarded/` prints nothing.
- [ ] A landed unrecorded target and a subsumed unrecorded target each apply with recovery `none` and write no discarded ref.
- [ ] A recorded active assignment by label plans through its record with no `class=` prefix.
- [ ] `--target <id> --apply-current` exits 2 with the usage line.
- [ ] A unique row in the unclaimed plan ends with `bench worktree clean --discard-branch --target <branch path>`.
- [ ] A symref at the planned discarded path makes the apply print an error row, keep both refs, and write nothing.
- [ ] An explicit target of a resolving Bench symref prints an error row naming `symref`, no fingerprint, and both apply forms refuse.
- [ ] After an explicit set removes a recorded holder, its formerly subsumed unrecorded member refuses as stale and survives.
- [ ] A direct ref planted at the planned path between the read and the write makes the apply refuse and keep both refs.
- [ ] A unique shift row ends with `bench worktree clean --discard-branch --target bench/shift-<stamp>`, and that command plans one row.
