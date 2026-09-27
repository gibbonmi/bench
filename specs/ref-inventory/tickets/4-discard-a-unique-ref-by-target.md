# 4. Discard one unique unrecorded ref by target with a discarded ref first

Blocked by: 1-classify-unclaimed-refs.md, 3-sweep-discarded-refs.md
Writes: internal/worktree/clean_discard.go (new), internal/worktree/clean_discard_test.go (new), internal/worktree/clean_set.go, internal/worktree/path.go, internal/worktree/worktree.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RI29, RI30, RI31, RI32, RI33, RI34, RI35, RI36, RI37, RI38, RI39, RI40, RI41

## What to build

Chunk: RI-C2b.

Add one fallback to `--target` resolution after the assignment resolver refuses with the unassigned error.
The fallback accepts a full branch path with or without `refs/heads/`.
It also accepts a 32-character hexadecimal id equal to the last segment of exactly one unrecorded assignment branch.
Two matches refuse with an error row that names both refs, and no match keeps the unassigned refusal.
A recorded assignment keeps today's route.

An unrecorded target row classifies through the class function of ticket 1.
A unique row plans the discarded ref of ticket 3 in its recovery cell, dated by the current UTC day.
Its detail starts with `class=unique`.
A landed or subsumed row plans recovery `none`.
The explicit set fingerprint binds each unrecorded row's ref, tip, class, holder, and planned discarded ref.
The apply needs the set fingerprint, and `--apply-current` stays refused outside the unclaimed mode.

The apply of a unique row runs in this order inside the set apply.
An existing ref at the planned path must resolve to the row's tip, or the apply refuses and keeps the branch.
The write uses the zero old value, and the delete uses the exact tip.
One fault boundary step precedes the write and one follows it.
The outcome row's recovery cell names the discarded ref.

## Acceptance

- [ ] `--target` with the full branch path, with the path without `refs/heads/`, and with the id segment each plan one row for the unrecorded branch.
- [ ] An id segment shared by two unrecorded branches prints an error row that names both refs and no fingerprint.
- [ ] `--target archive/x` prints the unassigned refusal.
- [ ] The plan row of a unique target shows `class=unique` and the recovery `refs/bench/discarded/<yyyymmdd>/bench/assign/<owner>/<id>`.
- [ ] After the apply, the discarded ref resolves to the old tip, the branch is gone, and the outcome row's recovery cell equals the discarded ref.
- [ ] A fault before the write leaves the branch and no ref under `refs/bench/discarded/`.
- [ ] A fault after the write leaves both refs, and a re-plan plus second apply removes the branch.
- [ ] A planted ref at the planned path at another commit makes the apply refuse and keep the branch.
- [ ] A landed unrecorded target applies with recovery `none` and writes no discarded ref.
- [ ] A recorded active assignment by label plans through its record with no `class=` prefix.
- [ ] `--target <id> --apply-current` exits 2 with the usage line.
