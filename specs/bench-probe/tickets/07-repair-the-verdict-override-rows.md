# Repair the verdict override rows from the review round

Blocked by: 06-state-the-probe-verb-in-the-guidance.md
Writes: internal/probe/probe.go, internal/probe/outcome_test.go, internal/probe/probe_test.go, internal/probe/refusal_test.go, internal/testreport/outcome_test.go
Covers: PB38, PB47

## What to build

Verify the premise first. Read `render` and `verdictFor` in internal/probe/probe.go.
Read `TestVerdictExitCodes` in internal/probe/outcome_test.go and
`TestProbeRestoresBeforeARenderRefusal` and `TestProbeReportsARestoreFailure` in the
probe tests. Read the spec's implementation decision on the render after the restore.

Make `render` obey the `restore-failed` override when the verdict row cannot render.
When the restore failed and the row is unrepresentable, print the render error line,
then the `preserved[1]{path,reason}` row, and exit 2. Write
`TestProbeNamesTheCopyWhenTheRenderAndTheRestoreFail` for row PB47 with the PB41
subject and the PB14 read-only directory.

Extend `TestVerdictExitCodes` so that the table drives `render` with a failed restore
over every outcome kind and asserts `restore-failed` at exit 2 each time.

Fold the three accepted Standards repairs. Add one local helper for the lock path in
internal/probe/refusal_test.go. Remove the unused `awaitFile` parameter. Remove the
provenance and narration comments from internal/testreport/outcome_test.go and
internal/probe/probe_test.go. Keep every golden and every assertion unchanged.

Self-probe: return the render error before the restore check and show
`TestProbeNamesTheCopyWhenTheRenderAndTheRestoreFail` red.

## Acceptance

- [ ] `TestProbeNamesTheCopyWhenTheRenderAndTheRestoreFail` prints the render error line, the preserved row, and exits 2.
- [ ] `TestVerdictExitCodes` answers `restore-failed` at 2 over every kind.
- [ ] `go test ./internal/probe ./internal/testreport` passes with every existing assertion unchanged.
