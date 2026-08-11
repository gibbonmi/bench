# Read each retained run from one validated object

Blocked by: none
Ownership fence: `internal/specbuild/state.go`, `internal/specbuild/runs_test.go`
Integration surfaces: retained-state directory enumeration and record reading→`internal/specbuild/state.go`; deterministic swap fixtures and bounded blocking guard→`internal/specbuild/runs_test.go`; family-home rendering→existing `cmd/bench/specbuild.go` exercised unchanged through the returned runs and diagnostics
Contracts: one retained record's validated identity and bytes cross the filesystem read inside `internal/specbuild/state.go`→the record decoder in the same owner, asserted by RS1 and RS2 against the real retained-state directory; diagnostic state crosses `internal/specbuild/state.go`→the existing family-home renderer, asserted by RS3 through the returned public projection
Closure: RS1/opened-object-identity, RS2/no-path-reread, RS3/symlink-cannot-redirect-open-object, RS4/fifo-cannot-block-open-object

## What to build

Close the accepted round-8 Spec and Coverage findings P2/C2. `Service.Runs`
must validate and decode one opened retained-state object rather than validate a
path and then read that path again. Once the validated object is open, a
deterministic path replacement cannot redirect the read into a symlink target,
another regular record, or a blocking special file.

## Acceptance

- [ ] [RS1] (covers local) (P2, C2) each retained record is decoded only from an opened regular object whose identity is the identity `Service.Runs` validated.
- [ ] [RS2] (covers local) (P2, C2) replacing the directory path with a second regular record after the validated object is open cannot change the bytes decoded from that object or substitute the second run.
- [ ] [RS3] (covers local) (P2, C2) replacing the directory path with a symlink after the validated object is open cannot redirect its read to the symlink target.
- [ ] [RS4] (covers local) (P2, C2) replacing the directory path with a FIFO after the validated object is open cannot block its read.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RS1/opened-object-identity | restore separate path validation and `os.ReadFile(path)` calls | the retained-run object-snapshot test | replace the path after validation, run the focused `Service.Runs` test, and require the decoded bytes and validated identity to come from one opened object |
| RS2/no-path-reread | reopen the pathname after the validated object is available | the regular-record replacement test | replace the path with a second valid record at the deterministic fault point, call `Service.Runs`, and require that the replacement cannot be reported as the original entry |
| RS3/symlink-cannot-redirect-open-object | reopen the path after a symlink replaces it | the opened-object symlink-swap test | open and validate the original record, replace its path with a symlink to a different valid record, then require `Service.Runs` to decode the original object rather than the symlink target |
| RS4/fifo-cannot-block-open-object | reopen the path after a FIFO replaces it | the bounded opened-object FIFO-swap test | open and validate the original record, replace its path with a FIFO, then require `Service.Runs` to return the original record within the existing one-second bound |
