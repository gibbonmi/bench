# Report malformed admin entries in bench doctor

Blocked by: refuse-malformed-admin-entries.md
Writes: internal/adopt

## What to build

On an adopted repo whose admin dir holds a malformed entry, `bench doctor`
prints a red row naming the entry; a healthy repo prints no new row, and
nothing under `<git-common-dir>/worktrees/` is ever touched. One new
`doctorRows` entry plus one evaluator that delegates to the predecessor's
exported `git.ScanWorktreeAdmin` (resolving the admin dir through
`git.CommonDir`) rather than re-deriving the shape predicate — the exported
scanner contract is the crossing this ticket consumes. The fixture must
chdir into an adopted repo (`reportDoctorRows` resolves `git.Root()` and
requires `.bench/`), and lives in `adopt_test.go`-style files — never
`internal/adopt/decision_test.go`, whose census entry forbids repository or
process constructors.

## Acceptance

- [ ] With a FIFO admin entry planted, `bench doctor` prints a red row containing `worktrees/<id>/gitdir` (covers WE10)
- [ ] With healthy worktrees, `bench doctor` prints no admin-entry row (covers WE11)
- [ ] After the doctor run over the FIFO fixture, the FIFO still exists by lstat (covers WE16)
