# Gate store records share one schema and time format

Blocked by: none
Writes: internal/gate/verdict.go, internal/gate/run_durability.go, internal/gate/lane.go, internal/gate/delta_checkpoint.go, internal/gate/reader_set.go, internal/gate/record_time_skew_test.go
Covers: none

## What to build

The gate store has five record constructors: the ready verdict, the interrupted record, the lane record, the delta record, and the reader record.
Each constructor sets the schema and formats its record time by itself.
The interrupted record writes the schema as the literal `1`, beside the `verdictSchema` constant that the reader grades.
Each constructor also repeats the UTC, whole-second, RFC 3339 time format.

One constructor in `verdict.go` now sets the fields that every store record shares: the schema and the tree.
One helper in `verdict.go` now formats a record time.
Each of the five record constructors starts from the shared constructor and calls the time helper.
The strict record time reader compares a parsed time against the same helper, so the writer and the reader share one time format.
The bytes of each record do not change, so this change has no `CHANGELOG.md` entry.

No test pinned the UTC conversion before this change, because each existing caller passes a UTC time.
A test in `record_time_skew_test.go` now gives a ready verdict a time in a non-UTC zone and requires the exact UTC form.
The same test requires the strict reader to refuse a value with an offset and a value with a fraction of a second.

## Acceptance

- [ ] No record constructor in the `internal/gate` production files sets the schema with a literal or with the `verdictSchema` constant directly; each constructor starts from the shared constructor.
- [ ] One helper formats each record time field, such as `recorded_at` and `started_at`, and each record constructor calls it.
- [ ] The strict record time reader refuses each value that the helper does not produce.
- [ ] The bytes of each record stay the same as before the change.
- [ ] A removal of the UTC conversion from the time helper turns `TestRecordTimeIsTheUTCWholeSecondForm` red.
- [ ] `bench test` for `./internal/gate` passes, and the lane passes at `bench commit`.
