# Derive the status harness grammar from the record

Blocked by: 01-add-the-harness-record-package.md
Writes: internal/status/route.go, internal/status/route_test.go, internal/status/status_command_test.go, internal/conformance/handoff_single_source_test.go

## What to build

`internal/status` builds its `harnessPrefix` table from the record's phase
forms at package init. The table therefore holds one entry per record row,
and a formless row holds an empty form. `HarnessChoices` and `ValidHarness`
read the record, so the grammar names all four harnesses with `claude` first
and the rest sorted.

`bench status -h` and `bench handoff -h` advertise the four names, because
both help lines compose `HarnessChoices`. An unknown `--harness` value stays
a usage error at exit 2 and prints the grammar. `--harness` without `--route`
stays a usage error, and `--route --all` stays a usage error.

The `harness-prefix-single-source` check moves its owner constants to the
record's file. Set `prefixTablePkgDir` to `internal/harnesses`,
`prefixTableFile` to the record's file name, and `prefixTableVar` to the
record's exported rows variable. The handoff package still holds no phase
form as a literal, and the check's own bite test still reds on a removed
table.

The check's `prefixTable` reader walks one level of key-value pairs today.
The record's rows are struct literals nested inside a slice literal, so the
reader must collect the string literal behind each `PhaseForm:` key there. An empty form contributes no forbidden string. The bite test's
synthetic table takes the record's nested shape, so the test proves the
reader against the shape the live tree carries.

This ticket owns the table's shape. Ticket 05 owns how `Route` treats a row
with an empty form, so this ticket keeps today's translation behavior.

## Acceptance

- [ ] `status.HarnessChoices()` returns the four record names with `claude` first and the rest sorted. (covers HC13)
- [ ] `bench status -h` advertises `--harness` with the four record names, `claude` first. (covers HC21)
- [ ] `bench status --route --harness cursor` prints the grammar and exits 2. (covers HC23)
- [ ] `internal/status` declares no harness name and no phase form as a literal.
- [ ] `checkHarnessPrefix` grades the record's file and stays green over the live root.
- [ ] The `harness-prefix-single-source` bite test reds on a record file without the rows table.
- [ ] `checkHarnessPrefix` collects `/bench-` and `$bench-` from a synthetic record in the nested row shape.
