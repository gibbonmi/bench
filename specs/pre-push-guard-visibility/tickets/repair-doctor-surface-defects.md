# Repair doctor surface defects

Blocked by: none
Ownership fence: `internal/adopt/doctor.go`, `internal/contract/surface/doctor_test.go`
Integration surfaces: shim-current PATH notice→`internal/adopt/doctor.go` + DSD1; doctor runtime contract→`internal/contract/surface/doctor_test.go` + DSD1
Contracts: none crosses
Closure: DSD1/current-path-notice, DSD2/stale-remedy-doc, DSD3/stale-classify-comment

## What to build

Restore the PATH guidance `doctorFix` lost on the shim-already-current branch:
when the shim is current but its directory is not on `PATH`, `bench doctor
--fix` must still emit the `doctorPathNotice` line before returning through the
stale-pre-push repair. Extend the doctor runtime contract to exercise the
already-current path, not only the fresh-write path. In the same fence, correct
two stale comments: the `reportPrePush` doc comment still claims doctor never
installs or relabels the hook and points at `bench link`, while the code now
recommends and performs `bench doctor --fix`; and a `doctor_test.go` comment
still names the removed `ClassifyPrePush`.

## Acceptance

- [ ] [DSD1] (covers local) `bench doctor --fix` with a current shim in a directory absent from `PATH` prints the PATH notice.
- [ ] [DSD2] (covers local) The `reportPrePush` doc comment states the current repair remedy semantics and does not claim doctor never repairs the hook.
- [ ] [DSD3] (covers local) No `internal/contract/surface/doctor_test.go` comment names `ClassifyPrePush`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DSD1/current-path-notice | drop the `doctorPathNotice` call from the shim-already-current branch | doctor runtime contract | seed a current shim off `PATH`, run `bench doctor --fix`, expect the PATH-notice assertion red |
| DSD2/stale-remedy-doc | restore the never-repairs claim to the doc comment | Standards review | read `reportPrePush` and its doc, expect the comment-contradicts-code finding |
| DSD3/stale-classify-comment | reintroduce `ClassifyPrePush` into the test comment | Standards review | enumerate `ClassifyPrePush` references under the fence, expect the stale-name finding |
